package heavyverification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const heavyManifestSchemaVersion = 1

type laneManifest struct {
	SchemaVersion    int            `json:"schema_version"`
	Lanes            []manifestLane `json:"lanes"`
	NonHeavyCoverage []coverageRule `json:"non_heavy_coverage"`
}

type manifestLane struct {
	ID       string         `json:"id"`
	Argv     []string       `json:"argv"`
	Coverage []coverageRule `json:"coverage"`
	// Fingerprint is optional and governs evidence reuse only. Lane selection
	// never reads it; a lane that omits it simply always executes. Keeping it
	// optional is what lets the manifest committed at an older revision, which
	// the revalidation planner still reads, decode unchanged.
	Fingerprint *laneFingerprintInputs `json:"fingerprint,omitempty"`
}

type coverageRule struct {
	Kind   string `json:"kind"`
	Path   string `json:"path,omitempty"`
	Root   string `json:"root,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

// Plan derives the heavy lanes required by the exact Git range. Selection is
// conservative: one uncovered path selects every declared lane.
func Plan(repositoryRoot, manifestPath, baseRevision, targetRevision string, forceAll bool) (resultmodel.HeavyVerificationPlan, error) {
	_, manifestRelativePath, err := resolveManifestPath(repositoryRoot, manifestPath)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, err
	}
	baseCommit, err := resolveRevision(repositoryRoot, baseRevision)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("resolve base revision %q: %w", baseRevision, err)
	}
	targetCommit, err := resolveRevision(repositoryRoot, targetRevision)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("resolve target revision %q: %w", targetRevision, err)
	}
	manifest, err := readManifestAtRevision(repositoryRoot, manifestRelativePath, targetCommit)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, err
	}
	changedPaths, err := diffChangedPaths(repositoryRoot, baseCommit, targetCommit)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, err
	}
	plan := resultmodel.HeavyVerificationPlan{
		ManifestPath: manifestRelativePath, BaseRevision: baseCommit, TargetRevision: targetCommit,
		ForcedAll: forceAll, ChangedPaths: changedPaths, UncoveredPaths: []string{}, SelectedLanes: []resultmodel.HeavyLaneSelection{},
	}
	return selectHeavyLanes(plan, manifest), nil
}

// PlanRevalidation maps one or more already-landed historical implementation ranges
// using the manifest committed at the exact revision where the selected lanes will run.
func PlanRevalidation(repositoryRoot, manifestPath string, sourceRanges []resultmodel.HeavySourceRange, executionRevision string, forceAll bool) (resultmodel.HeavyVerificationPlan, error) {
	if len(sourceRanges) == 0 {
		return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("at least one source range is required")
	}
	_, manifestRelativePath, err := resolveManifestPath(repositoryRoot, manifestPath)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, err
	}
	executionCommit, err := resolveRevision(repositoryRoot, executionRevision)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("resolve execution revision %q: %w", executionRevision, err)
	}
	manifest, err := readManifestAtRevision(repositoryRoot, manifestRelativePath, executionCommit)
	if err != nil {
		return resultmodel.HeavyVerificationPlan{}, err
	}
	resolvedRanges := make([]resultmodel.HeavySourceRange, 0, len(sourceRanges))
	changedPathSet := map[string]bool{}
	for rangeIndex, sourceRange := range sourceRanges {
		baseCommit, resolveError := resolveRevision(repositoryRoot, sourceRange.BaseRevision)
		if resolveError != nil {
			return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("resolve source range %d base revision %q: %w", rangeIndex+1, sourceRange.BaseRevision, resolveError)
		}
		targetCommit, resolveError := resolveRevision(repositoryRoot, sourceRange.TargetRevision)
		if resolveError != nil {
			return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("resolve source range %d target revision %q: %w", rangeIndex+1, sourceRange.TargetRevision, resolveError)
		}
		if ancestorError := requireAncestor(repositoryRoot, baseCommit, targetCommit); ancestorError != nil {
			return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("source range %d base must be an ancestor of its target: %w", rangeIndex+1, ancestorError)
		}
		if ancestorError := requireAncestor(repositoryRoot, targetCommit, executionCommit); ancestorError != nil {
			return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("source range %d target must be an ancestor of execution revision: %w", rangeIndex+1, ancestorError)
		}
		changedPaths, diffError := diffChangedPaths(repositoryRoot, baseCommit, targetCommit)
		if diffError != nil {
			return resultmodel.HeavyVerificationPlan{}, fmt.Errorf("source range %d: %w", rangeIndex+1, diffError)
		}
		for _, changedPath := range changedPaths {
			changedPathSet[changedPath] = true
		}
		resolvedRanges = append(resolvedRanges, resultmodel.HeavySourceRange{BaseRevision: baseCommit, TargetRevision: targetCommit})
	}
	changedPaths := make([]string, 0, len(changedPathSet))
	for changedPath := range changedPathSet {
		changedPaths = append(changedPaths, changedPath)
	}
	sort.Strings(changedPaths)
	plan := resultmodel.HeavyVerificationPlan{
		Mode: "historical-revalidation", SourceRanges: resolvedRanges,
		ManifestRevision: executionCommit, ExecutionRevision: executionCommit,
		ManifestPath: manifestRelativePath, ForcedAll: forceAll, ChangedPaths: changedPaths,
		UncoveredPaths: []string{}, SelectedLanes: []resultmodel.HeavyLaneSelection{},
	}
	return selectHeavyLanes(plan, manifest), nil
}

func selectHeavyLanes(plan resultmodel.HeavyVerificationPlan, manifest laneManifest) resultmodel.HeavyVerificationPlan {
	if plan.ForcedAll {
		plan.SelectedLanes = selectAllLanes(manifest.Lanes, "explicit force-all selected every declared heavy lane")
		return plan
	}

	selectedReasons := make(map[string][]string)
	for _, changedPath := range plan.ChangedPaths {
		covered := false
		for _, lane := range manifest.Lanes {
			for _, rule := range lane.Coverage {
				if rule.matches(changedPath) {
					covered = true
					selectedReasons[lane.ID] = append(selectedReasons[lane.ID], rule.reason(changedPath))
				}
			}
		}
		for _, rule := range manifest.NonHeavyCoverage {
			if rule.matches(changedPath) {
				covered = true
			}
		}
		if !covered {
			plan.UncoveredPaths = append(plan.UncoveredPaths, changedPath)
		}
	}
	if len(plan.UncoveredPaths) > 0 {
		plan.Uncertain = true
		reason := fmt.Sprintf("coverage is uncertain for: %s", strings.Join(plan.UncoveredPaths, ", "))
		plan.SelectedLanes = selectAllLanes(manifest.Lanes, reason)
		return plan
	}
	for _, lane := range manifest.Lanes {
		reasons := selectedReasons[lane.ID]
		if len(reasons) == 0 {
			continue
		}
		plan.SelectedLanes = append(plan.SelectedLanes, resultmodel.HeavyLaneSelection{
			LaneID: lane.ID, CommandArgv: append([]string(nil), lane.Argv...), Reasons: reasons,
		})
	}
	return plan
}

func requireAncestor(repositoryRoot, ancestorRevision, descendantRevision string) error {
	command := exec.Command("git", "-C", repositoryRoot, "merge-base", "--is-ancestor", ancestorRevision, descendantRevision)
	if output, err := command.CombinedOutput(); err != nil {
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return fmt.Errorf("%s is not an ancestor of %s: %s", ancestorRevision, descendantRevision, detail)
	}
	return nil
}

func resolveManifestPath(repositoryRoot, manifestPath string) (string, string, error) {
	if strings.TrimSpace(manifestPath) == "" {
		return "", "", fmt.Errorf("manifest path must not be empty")
	}
	repositoryAbsolutePath, err := filepath.Abs(repositoryRoot)
	if err != nil {
		return "", "", fmt.Errorf("resolve repository root: %w", err)
	}
	manifestAbsolutePath := manifestPath
	if !filepath.IsAbs(manifestAbsolutePath) {
		manifestAbsolutePath = filepath.Join(repositoryAbsolutePath, manifestAbsolutePath)
	}
	manifestAbsolutePath = filepath.Clean(manifestAbsolutePath)
	relativePath, err := filepath.Rel(repositoryAbsolutePath, manifestAbsolutePath)
	if err != nil || relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("manifest path must remain inside the repository")
	}
	return manifestAbsolutePath, filepath.ToSlash(relativePath), nil
}

func readManifestAtRevision(repositoryRoot, manifestRelativePath, targetRevision string) (laneManifest, error) {
	treeEntry, err := runGit(repositoryRoot, "ls-tree", "-z", targetRevision, "--", manifestRelativePath)
	if err != nil {
		return laneManifest{}, fmt.Errorf("inspect heavy-lane manifest at target revision: %w", err)
	}
	entryPrefix, entryPath, found := bytes.Cut(bytes.TrimSuffix(treeEntry, []byte{0}), []byte{'\t'})
	entryFields := bytes.Fields(entryPrefix)
	if !found || len(entryFields) != 3 || string(entryPath) != manifestRelativePath {
		return laneManifest{}, fmt.Errorf("heavy-lane manifest is missing or ambiguous at target revision")
	}
	manifestMode := string(entryFields[0])
	manifestType := string(entryFields[1])
	if manifestType != "blob" || (manifestMode != "100644" && manifestMode != "100755") {
		return laneManifest{}, fmt.Errorf("heavy-lane manifest must be a regular file at target revision")
	}
	contents, err := runGit(repositoryRoot, "cat-file", "blob", string(entryFields[2]))
	if err != nil {
		return laneManifest{}, fmt.Errorf("read heavy-lane manifest from target revision: %w", err)
	}
	return decodeManifest(contents)
}

func decodeManifest(contents []byte) (laneManifest, error) {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	var manifest laneManifest
	if err := decoder.Decode(&manifest); err != nil {
		return laneManifest{}, fmt.Errorf("decode heavy-lane manifest: %w", err)
	}
	var trailingValue any
	trailingError := decoder.Decode(&trailingValue)
	if trailingError == nil {
		return laneManifest{}, fmt.Errorf("decode heavy-lane manifest: trailing JSON value")
	}
	if trailingError != io.EOF {
		return laneManifest{}, fmt.Errorf("decode heavy-lane manifest: %w", trailingError)
	}
	if manifest.SchemaVersion != heavyManifestSchemaVersion {
		return laneManifest{}, fmt.Errorf("heavy-lane manifest schema_version must be %d", heavyManifestSchemaVersion)
	}
	if len(manifest.Lanes) == 0 {
		return laneManifest{}, fmt.Errorf("heavy-lane manifest must declare at least one lane")
	}
	seenLaneIDs := make(map[string]bool)
	for laneIndex, lane := range manifest.Lanes {
		if strings.TrimSpace(lane.ID) == "" {
			return laneManifest{}, fmt.Errorf("lane %d has an empty id", laneIndex)
		}
		if seenLaneIDs[lane.ID] {
			return laneManifest{}, fmt.Errorf("lane id %q is duplicated", lane.ID)
		}
		seenLaneIDs[lane.ID] = true
		if len(lane.Argv) == 0 {
			return laneManifest{}, fmt.Errorf("lane %q has empty argv", lane.ID)
		}
		for _, argument := range lane.Argv {
			if argument == "" {
				return laneManifest{}, fmt.Errorf("lane %q argv contains an empty token", lane.ID)
			}
		}
		if len(lane.Coverage) == 0 {
			return laneManifest{}, fmt.Errorf("lane %q has no coverage rules", lane.ID)
		}
		for _, rule := range lane.Coverage {
			if err := rule.validate(); err != nil {
				return laneManifest{}, fmt.Errorf("lane %q coverage: %w", lane.ID, err)
			}
		}
		if lane.Fingerprint != nil {
			if err := lane.Fingerprint.validate(); err != nil {
				return laneManifest{}, fmt.Errorf("lane %q: %w", lane.ID, err)
			}
		}
	}
	for _, rule := range manifest.NonHeavyCoverage {
		if err := rule.validate(); err != nil {
			return laneManifest{}, fmt.Errorf("non-heavy coverage: %w", err)
		}
	}
	return manifest, nil
}

func (rule coverageRule) validate() error {
	switch rule.Kind {
	case "exact", "subtree":
		if !validRepositoryPath(rule.Path) || rule.Root != "" || rule.Suffix != "" {
			return fmt.Errorf("%s rule requires only a repository-relative path", rule.Kind)
		}
	case "suffix-under":
		if (rule.Root != "." && !validRepositoryPath(rule.Root)) || rule.Path != "" || rule.Suffix == "" || strings.Contains(rule.Suffix, "/") {
			return fmt.Errorf("suffix-under rule requires only a repository-relative root and filename suffix")
		}
	default:
		return fmt.Errorf("unsupported rule kind %q", rule.Kind)
	}
	return nil
}

func validRepositoryPath(path string) bool {
	if path == "" || strings.Contains(path, "\\") || strings.HasPrefix(path, "/") {
		return false
	}
	return filepath.ToSlash(filepath.Clean(path)) == path && path != "." && path != ".." && !strings.HasPrefix(path, "../")
}

func (rule coverageRule) matches(path string) bool {
	switch rule.Kind {
	case "exact":
		return path == rule.Path
	case "subtree":
		return path == rule.Path || strings.HasPrefix(path, rule.Path+"/")
	case "suffix-under":
		underRoot := rule.Root == "." || path == rule.Root || strings.HasPrefix(path, rule.Root+"/")
		return underRoot && strings.HasSuffix(path, rule.Suffix)
	default:
		return false
	}
}

func (rule coverageRule) reason(path string) string {
	switch rule.Kind {
	case "exact":
		return fmt.Sprintf("%s matched exact path %s", path, rule.Path)
	case "subtree":
		return fmt.Sprintf("%s matched subtree %s", path, rule.Path)
	case "suffix-under":
		return fmt.Sprintf("%s matched suffix %s under %s", path, rule.Suffix, rule.Root)
	default:
		return path
	}
}

func resolveRevision(repositoryRoot, revision string) (string, error) {
	if strings.TrimSpace(revision) == "" {
		return "", fmt.Errorf("revision must not be empty")
	}
	output, err := runGit(repositoryRoot, "rev-parse", "--verify", revision+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func diffChangedPaths(repositoryRoot, baseRevision, targetRevision string) ([]string, error) {
	output, err := runGit(repositoryRoot, "diff", "--name-status", "-z", "--find-renames", baseRevision, targetRevision, "--")
	if err != nil {
		return nil, fmt.Errorf("read changed paths: %w", err)
	}
	tokens := bytes.Split(output, []byte{0})
	if len(tokens) > 0 && len(tokens[len(tokens)-1]) == 0 {
		tokens = tokens[:len(tokens)-1]
	}
	paths := make(map[string]bool)
	for tokenIndex := 0; tokenIndex < len(tokens); {
		status := string(tokens[tokenIndex])
		tokenIndex++
		if status == "" || tokenIndex >= len(tokens) {
			return nil, fmt.Errorf("read changed paths: malformed name-status record")
		}
		pathCount := 1
		if status[0] == 'R' || status[0] == 'C' {
			pathCount = 2
		}
		if tokenIndex+pathCount > len(tokens) {
			return nil, fmt.Errorf("read changed paths: incomplete %s record", status)
		}
		for pathOffset := 0; pathOffset < pathCount; pathOffset++ {
			path := string(tokens[tokenIndex+pathOffset])
			if path == "" {
				return nil, fmt.Errorf("read changed paths: empty path in %s record", status)
			}
			paths[path] = true
		}
		tokenIndex += pathCount
	}
	changedPaths := make([]string, 0, len(paths))
	for path := range paths {
		changedPaths = append(changedPaths, path)
	}
	sort.Strings(changedPaths)
	return changedPaths, nil
}

func selectAllLanes(lanes []manifestLane, reason string) []resultmodel.HeavyLaneSelection {
	selected := make([]resultmodel.HeavyLaneSelection, 0, len(lanes))
	for _, lane := range lanes {
		selected = append(selected, resultmodel.HeavyLaneSelection{
			LaneID: lane.ID, CommandArgv: append([]string(nil), lane.Argv...), Reasons: []string{reason},
		})
	}
	return selected
}

func runGit(repositoryRoot string, arguments ...string) ([]byte, error) {
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.Output()
	if err == nil {
		return output, nil
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return nil, fmt.Errorf("git %s: %s", strings.Join(arguments, " "), strings.TrimSpace(string(exitError.Stderr)))
	}
	return nil, fmt.Errorf("git %s: %w", strings.Join(arguments, " "), err)
}
