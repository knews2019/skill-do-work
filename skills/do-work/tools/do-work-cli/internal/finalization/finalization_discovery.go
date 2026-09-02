package finalization

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type discoveryCandidate struct {
	request        *repositorymodel.RequestFile
	completedAt    time.Time
	implementation []string
	lifecyclePaths []string
	effectivePaths []string
}

// discoverFinalizationJournals freezes one inventory before reading request
// bytes. Secret-classified paths and every pre-existing index entry stop the
// operation before association.
func discoverFinalizationJournals(repositoryRoot string) ([]*Journal, *resultmodel.CommandResult) {
	rows, err := corehelpers.ReadProtectedInventory(repositoryRoot)
	if err != nil {
		failure := commandFailure(repositoryRoot, CommandRecoverFinalization, "FINALIZATION-DISCOVERY-INVENTORY", err.Error())
		return nil, &failure
	}
	blocked := []string{}
	for _, row := range rows {
		if row.Classification == "X" || row.Classification == "XD" {
			blocked = append(blocked, row.Path)
		}
	}
	if len(blocked) > 0 {
		failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-PROTECTED", "protected inventory contains unreadable paths", blocked)
		return nil, &failure
	}
	staged, err := gitLines(repositoryRoot, "diff", "--cached", "--name-only")
	if err != nil {
		failure := commandFailure(repositoryRoot, CommandRecoverFinalization, "FINALIZATION-DISCOVERY-INDEX", err.Error())
		return nil, &failure
	}
	if len(staged) > 0 {
		failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-STAGED", "discovery requires an empty index", staged)
		return nil, &failure
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		failure := commandFailure(repositoryRoot, CommandRecoverFinalization, "FINALIZATION-DISCOVERY-SNAPSHOT", err.Error())
		return nil, &failure
	}
	dirty := map[string]bool{}
	origins := map[string]string{}
	for _, row := range rows {
		dirty[row.Path] = true
		if row.Origin != "" {
			dirty[row.Origin] = true
			origins[row.Path] = row.Origin
		}
	}
	candidates := []*discoveryCandidate{}
	for _, request := range snapshot.RequestFiles {
		record := request.TypedRecord
		if request.TreeSection != "archive" || !legacyTerminalSuccess(record.RequestStatus) || record.FieldEvidenceByName["commit"].ScalarValue != "" {
			continue
		}
		completedAt, parseError := requestmodel.ParseTimestamp(record.CompletedAt)
		if parseError != nil || completedAt.IsZero() {
			continue
		}
		implementation, parseError := implementationSummaryPaths(string(request.ContentBytes))
		if parseError != nil || len(implementation) == 0 {
			continue
		}
		hasDirtyEvidence := dirty[requestRepositoryPath(request.RelativePath)]
		for _, path := range implementation {
			hasDirtyEvidence = hasDirtyEvidence || dirty[path]
		}
		if !hasDirtyEvidence {
			continue
		}
		candidates = append(candidates, &discoveryCandidate{request: request, completedAt: completedAt, implementation: implementation})
	}
	if len(candidates) == 0 {
		return nil, nil
	}

	owners := map[string][]*discoveryCandidate{}
	for _, candidate := range candidates {
		for _, path := range candidate.implementation {
			if path != "do-work" && !strings.HasPrefix(path, "do-work/") {
				owners[path] = append(owners[path], candidate)
			}
		}
	}
	ambiguous := []string{}
	for path := range dirty {
		if path == "do-work" || strings.HasPrefix(path, "do-work/") {
			continue
		}
		matches := owners[path]
		if len(matches) == 1 {
			matches[0].effectivePaths = append(matches[0].effectivePaths, path)
		} else if len(matches) > 1 {
			ambiguous = append(ambiguous, path)
		}
	}
	for path := range dirty {
		if path != "do-work" && !strings.HasPrefix(path, "do-work/") {
			continue
		}
		matches := []*discoveryCandidate{}
		for _, candidate := range candidates {
			requestID := candidate.request.TypedRecord.RequestID
			requestPath := requestRepositoryPath(candidate.request.RelativePath)
			if path == requestPath || path == origins[requestPath] || origins[path] == requestPath ||
				(strings.HasPrefix(path, "do-work/archive/") || strings.HasPrefix(path, "do-work/working/")) && pathCarriesRequestID(repositoryRoot, path, requestID) ||
				path == "do-work/CHECKPOINT.md" && checkpointRemovalProves(repositoryRoot, path, requestID) {
				matches = append(matches, candidate)
			}
		}
		if len(matches) == 1 {
			matches[0].lifecyclePaths = append(matches[0].lifecyclePaths, path)
			matches[0].effectivePaths = append(matches[0].effectivePaths, path)
		} else {
			ambiguous = append(ambiguous, path)
		}
	}
	if len(ambiguous) > 0 {
		failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-AMBIGUOUS", "shared or multiply-owned state cannot be associated exactly", uniqueSorted(ambiguous))
		return nil, &failure
	}

	sort.Slice(candidates, func(left, right int) bool {
		if candidates[left].completedAt.Equal(candidates[right].completedAt) {
			return candidates[left].request.TypedRecord.RequestID < candidates[right].request.TypedRecord.RequestID
		}
		return candidates[left].completedAt.Before(candidates[right].completedAt)
	})
	journals := make([]*Journal, 0, len(candidates))
	for _, candidate := range candidates {
		journal, buildError := discoveredJournal(repositoryRoot, candidate)
		if buildError != nil {
			failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-INVALID", buildError.Error(), candidate.effectivePaths)
			return nil, &failure
		}
		journals = append(journals, journal)
	}
	return journals, nil
}

func pathCarriesRequestID(repositoryRoot, path, requestID string) bool {
	current, err := currentImage(repositoryRoot, path)
	if err == nil && current.Exists {
		if document, parseError := requestmodel.ParseDocument(current.Bytes); parseError == nil && document.TypedRecord().RequestID == requestID {
			return true
		}
	}
	before, err := exec.Command("git", "-C", repositoryRoot, "show", "HEAD:"+path).Output()
	if err != nil {
		return false
	}
	document, err := requestmodel.ParseDocument(before)
	return err == nil && document.TypedRecord().RequestID == requestID
}

func discoveredJournal(repositoryRoot string, candidate *discoveryCandidate) (*Journal, error) {
	requestID := candidate.request.TypedRecord.RequestID
	requestPath := requestRepositoryPath(candidate.request.RelativePath)
	journalPath, _, err := journalLocations(repositoryRoot, requestID)
	if err != nil {
		return nil, err
	}
	paths := uniqueSorted(candidate.effectivePaths)
	preimages := make([]FileImage, 0, len(candidate.lifecyclePaths))
	postimages := make([]FileImage, 0, len(candidate.lifecyclePaths))
	for _, path := range uniqueSorted(candidate.lifecyclePaths) {
		before, imageError := headFileImage(repositoryRoot, path)
		if imageError != nil {
			return nil, imageError
		}
		after, imageError := currentImage(repositoryRoot, path)
		if imageError != nil {
			return nil, imageError
		}
		preimages = append(preimages, before)
		postimages = append(postimages, after)
	}
	manifest := Manifest{
		RequestID: requestID, RequestPath: requestPath, WriterLabel: "legacy-discovery",
		Transition: "complete", TerminalStatus: candidate.request.TypedRecord.RequestStatus,
		CompletedAt: candidate.request.TypedRecord.CompletedAt, ExpectedRequestSHA256: digestBytes(candidate.request.ContentBytes),
		ExpectedCheckpointSHA256: strings.Repeat("0", 64), CommitPaths: paths,
		CommitMessage: "[" + requestID + "] finalize recovered legacy tail", ProvenanceMode: ProvenancePrimaryCommit,
	}
	manifestBytes, _ := json.Marshal(manifest)
	preparedHead, preparedDiff, err := preparedCommitIdentity(repositoryRoot, paths)
	if err != nil {
		return nil, err
	}
	journal := &Journal{
		Version: journalVersion, CreatedAt: candidate.completedAt, UpdatedAt: time.Now().UTC().Truncate(time.Second), Phase: PhaseReleaseApplied,
		Manifest: manifest, ManifestSHA256: digestBytes(manifestBytes), JournalPath: journalPath, ArchivedPath: requestPath,
		LifecyclePreimages: preimages, LifecyclePostimages: postimages, ReleasePreimages: []FileImage{}, ReleasePostimages: []FileImage{},
		EffectiveCommitPaths: paths, PreparedHead: preparedHead, PreparedDiffSHA256: preparedDiff, Discovered: true,
	}
	if err := writeJournal(journal); err != nil {
		return nil, err
	}
	return journal, nil
}

func requestRepositoryPath(relativePath string) string {
	if relativePath == "do-work" || strings.HasPrefix(relativePath, "do-work/") {
		return relativePath
	}
	return "do-work/" + relativePath
}

func implementationSummaryPaths(contents string) ([]string, error) {
	lines := strings.Split(contents, "\n")
	inSection := false
	paths := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if inSection {
				break
			}
			inSection = strings.TrimSpace(line) == "## Implementation Summary"
			continue
		}
		if !inSection {
			continue
		}
		if strings.Count(line, "`")%2 != 0 {
			return nil, fmt.Errorf("unmatched backtick in Implementation Summary")
		}
		parts := strings.Split(line, "`")
		for index := 1; index < len(parts); index += 2 {
			path := filepath.ToSlash(filepath.Clean(parts[index]))
			if path != "." && path != ".." && !strings.HasPrefix(path, "../") && !filepath.IsAbs(path) {
				paths = append(paths, path)
			}
		}
	}
	return uniqueSorted(paths), nil
}

func checkpointRemovalProves(repositoryRoot, checkpointPath, requestID string) bool {
	before, err := exec.Command("git", "-C", repositoryRoot, "show", "HEAD:"+checkpointPath).Output()
	requestPattern := regexp.MustCompile(`(^|[^0-9])` + regexp.QuoteMeta(requestID) + `([^0-9]|$)`)
	if err != nil || !requestPattern.Match(before) {
		return false
	}
	after, err := currentImage(repositoryRoot, checkpointPath)
	return err == nil && after.Exists && !requestPattern.Match(after.Bytes)
}

func headFileImage(repositoryRoot, path string) (FileImage, error) {
	contents, err := exec.Command("git", "-C", repositoryRoot, "show", "HEAD:"+path).Output()
	if err != nil {
		return FileImage{Path: path}, nil
	}
	return FileImage{Path: path, Exists: true, Bytes: contents, Mode: uint32(0o644)}, nil
}

func preparedCommitIdentity(repositoryRoot string, paths []string) (string, string, error) {
	head := currentHead(repositoryRoot)
	arguments := append([]string{"-C", repositoryRoot, "diff", "--binary", head, "--"}, paths...)
	diff, err := exec.Command("git", arguments...).Output()
	if err != nil {
		return "", "", err
	}
	if len(diff) == 0 {
		return "", "", fmt.Errorf("discovered finalization has no exact diff")
	}
	return head, digestBytes(diff), nil
}

func gitLines(repositoryRoot string, arguments ...string) ([]string, error) {
	output, err := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...).Output()
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func legacyTerminalSuccess(status string) bool {
	return status == "completed" || status == "completed-with-issues"
}

func discoveryRefusal(repositoryRoot, code, reason string, paths []string) resultmodel.CommandResult {
	verification := []string{"do-work-cli", "--format", "json", CommandRecoverFinalization, "--discover"}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, AffectedPaths: paths, Evidence: []string{reason}, Fixability: resultmodel.FixabilityRefused,
		AutomationStopReason: "legacy finalization evidence is ambiguous", NextArgv: verification, VerificationArgv: verification,
	}}}
}
