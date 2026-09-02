package finalization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type discoveryCandidate struct {
	request         *repositorymodel.RequestFile
	completedAt     time.Time
	implementation  []string
	lifecyclePaths  []string
	releasePaths    []string
	effectivePaths  []string
	attributedPaths []string
}

// discoverFinalizationJournals freezes one protected inventory before reading
// request bytes. Staged state always refuses. Unstaged X/XD rows are deliberately
// omitted from association, so their contents are never read, staged, or changed.
func discoverFinalizationJournals(repositoryRoot string, assumeSoleReleaser bool) ([]*Journal, *resultmodel.CommandResult) {
	rows, err := corehelpers.ReadProtectedInventory(repositoryRoot)
	if err != nil {
		failure := commandFailure(repositoryRoot, CommandRecoverFinalization, "FINALIZATION-DISCOVERY-INVENTORY", err.Error())
		return nil, &failure
	}
	staged, err := gitLines(repositoryRoot, "diff", "--cached", "--name-only")
	if err != nil {
		failure := commandFailure(repositoryRoot, CommandRecoverFinalization, "FINALIZATION-DISCOVERY-INDEX", err.Error())
		return nil, &failure
	}
	if len(staged) > 0 {
		protected := map[string]bool{}
		for _, row := range rows {
			if row.Classification == "X" || row.Classification == "XD" {
				protected[row.Path] = true
				if row.Origin != "" {
					protected[row.Origin] = true
				}
			}
		}
		protectedStaged := []string{}
		for _, path := range staged {
			if protected[path] {
				protectedStaged = append(protectedStaged, path)
			}
		}
		if len(protectedStaged) > 0 {
			failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-PROTECTED-STAGED", "protected paths are already staged", uniqueSorted(protectedStaged))
			return nil, &failure
		}
		failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-FOREIGN-STAGED", "foreign staged entries prevent exact finalization recovery", uniqueSorted(staged))
		return nil, &failure
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		failure := commandFailure(repositoryRoot, CommandRecoverFinalization, "FINALIZATION-DISCOVERY-SNAPSHOT", err.Error())
		return nil, &failure
	}
	dirty := map[string]bool{}
	origins := map[string]string{}
	destinations := map[string]string{}
	for _, row := range rows {
		if row.Classification == "X" || row.Classification == "XD" {
			continue
		}
		dirty[row.Path] = true
		if row.Origin != "" {
			dirty[row.Origin] = true
			origins[row.Path] = row.Origin
			destinations[row.Origin] = row.Path
		}
	}
	candidates := legacyCandidates(snapshot)
	if len(candidates) == 0 {
		if blocked := ambiguousSharedRemainder(dirty); len(blocked) > 0 {
			failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-AMBIGUOUS", "shared finalization state has no terminal request owner", blocked)
			return nil, &failure
		}
		return nil, nil
	}

	owners := map[string][]*discoveryCandidate{}
	for _, candidate := range candidates {
		for _, path := range candidate.implementation {
			if !sharedFinalizationPath(path) {
				owners[path] = append(owners[path], candidate)
			}
		}
	}
	ambiguous := []string{}
	for path := range dirty {
		if sharedFinalizationPath(path) || releaseMetadataPath(path) {
			continue
		}
		matches := owners[path]
		if len(matches) == 1 && ordinaryWholeDiffProves(repositoryRoot, path) {
			matches[0].effectivePaths = append(matches[0].effectivePaths, path)
		} else if len(matches) > 0 {
			ambiguous = append(ambiguous, path)
		}
	}
	if assumeSoleReleaser {
		active := activeDiscoveryCandidates(candidates)
		if len(active) > 1 {
			paths := []string{}
			for _, candidate := range active {
				paths = append(paths, requestRepositoryPath(candidate.request.RelativePath))
			}
			failure := discoveryRefusal(repositoryRoot, "FINALIZATION-MULTIPLE-TAILS", "sole-releaser attribution requires exactly one legacy finalization tail", uniqueSorted(paths))
			return nil, &failure
		}
	}

	unresolvedShared := []string{}
	for path := range dirty {
		if !sharedFinalizationPath(path) {
			continue
		}
		matches := []*discoveryCandidate{}
		for _, candidate := range candidates {
			if lifecyclePathProves(repositoryRoot, path, origins[path], destinations[path], candidate) {
				matches = append(matches, candidate)
			}
		}
		if len(matches) == 1 {
			matches[0].lifecyclePaths = append(matches[0].lifecyclePaths, path)
			matches[0].effectivePaths = append(matches[0].effectivePaths, path)
		} else {
			unresolvedShared = append(unresolvedShared, path)
		}
	}
	if assumeSoleReleaser {
		active := activeDiscoveryCandidates(candidates)
		if len(active) > 1 {
			paths := []string{}
			for _, candidate := range active {
				paths = append(paths, requestRepositoryPath(candidate.request.RelativePath))
			}
			failure := discoveryRefusal(repositoryRoot, "FINALIZATION-MULTIPLE-TAILS", "sole-releaser attribution requires exactly one legacy finalization tail", uniqueSorted(paths))
			return nil, &failure
		}
		if len(active) == 1 {
			for _, path := range unresolvedShared {
				if soleReleaserSharedPath(path, origins[path], destinations[path], active[0]) {
					active[0].lifecyclePaths = append(active[0].lifecyclePaths, path)
					active[0].effectivePaths = append(active[0].effectivePaths, path)
					active[0].attributedPaths = append(active[0].attributedPaths, path)
				} else {
					ambiguous = append(ambiguous, path)
				}
			}
		} else {
			ambiguous = append(ambiguous, unresolvedShared...)
		}
	} else {
		ambiguous = append(ambiguous, unresolvedShared...)
	}
	ambiguous = append(ambiguous, associateReleaseMetadata(repositoryRoot, dirty, candidates)...)
	if len(ambiguous) > 0 {
		failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-AMBIGUOUS", "shared, foreign-hunk, or multiply-owned state cannot be associated exactly", uniqueSorted(ambiguous))
		return nil, &failure
	}

	admitted := make([]*discoveryCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		candidate.effectivePaths = uniqueSorted(candidate.effectivePaths)
		candidate.attributedPaths = uniqueSorted(candidate.attributedPaths)
		if len(candidate.effectivePaths) > 0 {
			admitted = append(admitted, candidate)
		}
	}
	if assumeSoleReleaser && len(admitted) > 1 {
		paths := []string{}
		for _, candidate := range admitted {
			paths = append(paths, requestRepositoryPath(candidate.request.RelativePath))
		}
		failure := discoveryRefusal(repositoryRoot, "FINALIZATION-MULTIPLE-TAILS", "sole-releaser attribution requires exactly one legacy finalization tail", uniqueSorted(paths))
		return nil, &failure
	}
	sort.Slice(admitted, func(left, right int) bool {
		if admitted[left].completedAt.Equal(admitted[right].completedAt) {
			return admitted[left].request.TypedRecord.RequestID < admitted[right].request.TypedRecord.RequestID
		}
		return admitted[left].completedAt.Before(admitted[right].completedAt)
	})
	journals := make([]*Journal, 0, len(admitted))
	for _, candidate := range admitted {
		journal, buildError := discoveredJournal(repositoryRoot, candidate)
		if buildError != nil {
			failure := discoveryRefusal(repositoryRoot, "FINALIZATION-DISCOVERY-INVALID", buildError.Error(), candidate.effectivePaths)
			return nil, &failure
		}
		journals = append(journals, journal)
	}
	return journals, nil
}

func activeDiscoveryCandidates(candidates []*discoveryCandidate) []*discoveryCandidate {
	active := []*discoveryCandidate{}
	for _, candidate := range candidates {
		if len(candidate.effectivePaths) > 0 {
			active = append(active, candidate)
		}
	}
	return active
}

func soleReleaserSharedPath(path, origin, destination string, candidate *discoveryCandidate) bool {
	if path == "do-work/CHECKPOINT.md" || path == "do-work/calibration-log.tsv" {
		return true
	}
	userRequestID := candidate.request.TypedRecord.UserRequestID
	if userRequestID == "" || origin == "" && destination == "" {
		return false
	}
	activePrefix := "do-work/user-requests/" + userRequestID + "/"
	archivePrefix := "do-work/archive/" + userRequestID + "/"
	for _, candidatePath := range []string{path, origin, destination} {
		if strings.HasPrefix(candidatePath, activePrefix) || strings.HasPrefix(candidatePath, archivePrefix) {
			return true
		}
	}
	return false
}

func ordinaryWholeDiffProves(repositoryRoot, path string) bool {
	output, err := exec.Command("git", "-C", repositoryRoot, "diff", "--no-ext-diff", "--no-textconv", "--unified=0", "HEAD", "--", path).Output()
	if err != nil || len(output) == 0 {
		return false
	}
	// A legacy summary can prove whole-path ownership but carries no hunk
	// manifest. Admit one indivisible hunk (or one binary whole-file delta) and
	// refuse a multi-hunk file because its independent regions cannot be
	// attributed without guessing.
	return bytes.Count(output, []byte("\n@@")) <= 1
}

func legacyCandidates(snapshot *repositorymodel.RepositorySnapshot) []*discoveryCandidate {
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
		if parseError != nil {
			continue
		}
		candidates = append(candidates, &discoveryCandidate{request: request, completedAt: completedAt, implementation: implementation})
	}
	return candidates
}

func lifecyclePathProves(repositoryRoot, path, origin, destination string, candidate *discoveryCandidate) bool {
	record := candidate.request.TypedRecord
	requestID := record.RequestID
	requestPath := requestRepositoryPath(candidate.request.RelativePath)
	if path == requestPath || path == origin || path == destination || strings.HasPrefix(path, "do-work/working/") && filepath.Base(path) == filepath.Base(requestPath) {
		requestOrigin := origin
		if path != requestPath && (destination != "" || strings.HasPrefix(path, "do-work/working/")) {
			requestOrigin = path
		}
		if lifecycleRequestMoveProves(repositoryRoot, requestPath, requestOrigin, destination, requestID) {
			return true
		}
	}
	if path == "do-work/CHECKPOINT.md" {
		return checkpointRemovalProves(repositoryRoot, path, requestID)
	}
	if path == "do-work/calibration-log.tsv" {
		return calibrationAppendProves(repositoryRoot, path, record)
	}
	if !strings.HasPrefix(path, "do-work/") {
		return false
	}
	if followupPathProves(repositoryRoot, path, requestID) {
		return true
	}
	return userRequestMoveProves(repositoryRoot, path, origin, destination, record.UserRequestID, requestID)
}

func lifecycleRequestMoveProves(repositoryRoot, requestPath, origin, destination, requestID string) bool {
	after, err := currentImage(repositoryRoot, requestPath)
	if err != nil || !after.Exists {
		return false
	}
	afterDocument, err := requestmodel.ParseDocument(after.Bytes)
	if err != nil || afterDocument.TypedRecord().RequestID != requestID || !legacyTerminalSuccess(afterDocument.TypedRecord().RequestStatus) || afterDocument.TypedRecord().FieldEvidenceByName["commit"].ScalarValue != "" {
		return false
	}
	beforePath := origin
	if beforePath == "" && destination != "" {
		beforePath = requestPath
	}
	if beforePath == "" {
		beforePath = "do-work/working/" + filepath.Base(requestPath)
	}
	before, err := headFileImage(repositoryRoot, beforePath)
	if err != nil || !before.Exists {
		before, err = headFileImage(repositoryRoot, requestPath)
	}
	if err != nil || !before.Exists {
		return false
	}
	beforeDocument, err := requestmodel.ParseDocument(before.Bytes)
	if err != nil || beforeDocument.TypedRecord().RequestID != requestID {
		return false
	}
	for _, field := range []string{"status", "completed_at", "release_at", "commit"} {
		_ = beforeDocument.DeleteField(field)
		_ = afterDocument.DeleteField(field)
	}
	return bytes.Equal(beforeDocument.DocumentBytes(), afterDocument.DocumentBytes())
}

func checkpointRemovalProves(repositoryRoot, checkpointPath, requestID string) bool {
	before, err := exec.Command("git", "-C", repositoryRoot, "show", "HEAD:"+checkpointPath).Output()
	if err != nil {
		return false
	}
	after, err := currentImage(repositoryRoot, checkpointPath)
	if err != nil || !after.Exists {
		return false
	}
	beforeLines := strings.Split(string(before), "\n")
	afterLines := strings.Split(string(after.Bytes), "\n")
	prefix := 0
	for prefix < len(beforeLines) && prefix < len(afterLines) && beforeLines[prefix] == afterLines[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(beforeLines)-prefix && suffix < len(afterLines)-prefix && beforeLines[len(beforeLines)-1-suffix] == afterLines[len(afterLines)-1-suffix] {
		suffix++
	}
	if prefix+suffix != len(afterLines) || prefix >= len(beforeLines)-suffix {
		return false
	}
	removed := beforeLines[prefix : len(beforeLines)-suffix]
	requestPattern := regexp.MustCompile(`(^|[^0-9])` + regexp.QuoteMeta(requestID) + `([^0-9]|$)`)
	if len(removed) == 0 || !strings.HasPrefix(strings.TrimSpace(removed[0]), "-") || !requestPattern.MatchString(removed[0]) {
		return false
	}
	for _, line := range removed[1:] {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			return false
		}
	}
	return true
}

func calibrationAppendProves(repositoryRoot, path string, record requestmodel.RequestRecord) bool {
	before, err := exec.Command("git", "-C", repositoryRoot, "show", "HEAD:"+path).Output()
	if err != nil {
		before = nil
	}
	after, err := currentImage(repositoryRoot, path)
	if err != nil || !after.Exists || !bytes.HasPrefix(after.Bytes, before) {
		return false
	}
	estimate := record.FieldEvidenceByName["estimate"].NestedValues["p50_active_minutes"]
	claimedAt, claimedError := requestmodel.ParseTimestamp(record.ClaimedAt)
	completedAt, completedError := requestmodel.ParseTimestamp(record.CompletedAt)
	if estimate == "" || claimedError != nil || completedError != nil {
		return false
	}
	route := record.RouteValue
	if route == "" {
		route = "-"
	}
	want := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\n", record.RequestID, route, estimate, int(completedAt.Sub(claimedAt).Minutes()), requestmodel.CanonicalTimestamp(completedAt))
	return string(after.Bytes[len(before):]) == want
}

func followupPathProves(repositoryRoot, path, requestID string) bool {
	current, err := currentImage(repositoryRoot, path)
	if err != nil || !current.Exists {
		return false
	}
	document, err := requestmodel.ParseDocument(current.Bytes)
	if err != nil || document.TypedRecord().AddendumTo != requestID {
		return false
	}
	before, err := headFileImage(repositoryRoot, path)
	if err != nil {
		return false
	}
	if !before.Exists {
		return true
	}
	beforeDocument, err := requestmodel.ParseDocument(before.Bytes)
	if err != nil || beforeDocument.TypedRecord().RequestID != document.TypedRecord().RequestID || beforeDocument.TypedRecord().AddendumTo != requestID {
		return false
	}
	if !bytes.HasPrefix(current.Bytes, before.Bytes) {
		return false
	}
	appended := strings.TrimSpace(string(current.Bytes[len(before.Bytes):]))
	return (strings.HasPrefix(appended, "## Review Fold") || strings.HasPrefix(appended, "## Recovery Fold")) && strings.Contains(appended, requestID)
}

func userRequestMoveProves(repositoryRoot, path, origin, destination, userRequestID, requestID string) bool {
	if userRequestID == "" {
		return false
	}
	sourcePath, destinationPath := origin, path
	if destination != "" {
		sourcePath, destinationPath = path, destination
	}
	if sourcePath == "" {
		archivePrefix := "do-work/archive/" + userRequestID + "/"
		switch {
		case strings.HasPrefix(path, archivePrefix):
			base := filepath.Base(path)
			if base == "input.md" {
				sourcePath = "do-work/user-requests/" + userRequestID + "/input.md"
			} else {
				sourcePath = "do-work/archive/" + base
			}
		case strings.HasPrefix(path, "do-work/user-requests/"+userRequestID+"/"):
			sourcePath, destinationPath = path, archivePrefix+filepath.Base(path)
		case strings.HasPrefix(path, "do-work/archive/") && !strings.Contains(strings.TrimPrefix(path, "do-work/archive/"), "/"):
			sourcePath, destinationPath = path, archivePrefix+filepath.Base(path)
		default:
			return false
		}
	}
	before, beforeError := headFileImage(repositoryRoot, sourcePath)
	after, afterError := currentImage(repositoryRoot, destinationPath)
	if beforeError != nil || afterError != nil || !before.Exists || !after.Exists || !bytes.Equal(before.Bytes, after.Bytes) {
		return false
	}
	document, err := requestmodel.ParseDocument(after.Bytes)
	if err != nil {
		return false
	}
	record := document.TypedRecord()
	if record.RequestID == userRequestID {
		for _, member := range record.FieldEvidenceByName["requests"].ListValues {
			if member == requestID {
				return true
			}
		}
		return false
	}
	return record.UserRequestID == userRequestID
}

func associateReleaseMetadata(repositoryRoot string, dirty map[string]bool, candidates []*discoveryCandidate) []string {
	paths := []string{}
	for path := range dirty {
		if releaseMetadataPath(path) {
			paths = append(paths, path)
		}
	}
	paths = uniqueSorted(paths)
	if len(paths) == 0 {
		return nil
	}
	versions := map[string]struct{}{}
	oldVersions := map[string]struct{}{}
	changelogPaths := []string{}
	var firstChangelogBefore, firstChangelogAfter []byte
	for _, path := range paths {
		before, _ := headFileImage(repositoryRoot, path)
		after, err := currentImage(repositoryRoot, path)
		if err != nil || !before.Exists || !after.Exists {
			return paths
		}
		if strings.HasPrefix(filepath.Base(path), "CHANGELOG") {
			if firstChangelogBefore == nil {
				firstChangelogBefore, firstChangelogAfter = before.Bytes, after.Bytes
			} else if !bytes.Equal(firstChangelogBefore, before.Bytes) || !bytes.Equal(firstChangelogAfter, after.Bytes) {
				return paths
			}
			changelogPaths = append(changelogPaths, path)
			continue
		}
		oldVersion, oldOK := releaseVersion(path, before.Bytes)
		newVersion, newOK := releaseVersion(path, after.Bytes)
		if !oldOK || !newOK || !semverGreater(newVersion, oldVersion) || !semanticVersionReplacement(path, before.Bytes, after.Bytes, oldVersion, newVersion) {
			return paths
		}
		oldVersions[oldVersion] = struct{}{}
		versions[newVersion] = struct{}{}
	}
	if len(changelogPaths) == 0 || len(versions) != 1 || len(oldVersions) != 1 {
		return paths
	}
	oldVersion := ""
	for version := range oldVersions {
		oldVersion = version
	}
	requiredPaths := configuredReleaseMetadataPaths(repositoryRoot, oldVersion, firstChangelogBefore)
	missingRequired := []string{}
	for _, requiredPath := range requiredPaths {
		if !dirty[requiredPath] {
			missingRequired = append(missingRequired, requiredPath)
		}
	}
	if len(missingRequired) > 0 {
		return uniqueSorted(missingRequired)
	}
	newVersion := ""
	for version := range versions {
		newVersion = version
	}
	matches := []*discoveryCandidate{}
	for _, candidate := range candidates {
		if candidate.request.TypedRecord.FieldEvidenceByName["release_at"].ScalarValue == "" {
			continue
		}
		valid := true
		for _, path := range changelogPaths {
			before, _ := headFileImage(repositoryRoot, path)
			after, _ := currentImage(repositoryRoot, path)
			inserted, ok := singleInsertion(before.Bytes, after.Bytes)
			if !ok || !bytes.Contains(inserted, []byte(candidate.request.TypedRecord.RequestID)) || !bytes.Contains(inserted, []byte(newVersion)) {
				valid = false
				break
			}
		}
		if valid {
			matches = append(matches, candidate)
		}
	}
	if len(matches) != 1 {
		return paths
	}
	for _, path := range paths {
		matches[0].releasePaths = append(matches[0].releasePaths, path)
		matches[0].effectivePaths = append(matches[0].effectivePaths, path)
	}
	return nil
}

func semverGreater(newVersion, oldVersion string) bool {
	newParts, oldParts := strings.Split(newVersion, "."), strings.Split(oldVersion, ".")
	for index := 0; index < 3; index++ {
		newPart, _ := strconv.Atoi(newParts[index])
		oldPart, _ := strconv.Atoi(oldParts[index])
		if newPart != oldPart {
			return newPart > oldPart
		}
	}
	return false
}

func releaseMetadataPath(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, "CHANGELOG") || base == "VERSION" || base == "package.json" || base == "package-lock.json" || base == "Cargo.toml" || base == "Cargo.lock" || base == "pyproject.toml" || base == "uv.lock" || path == "skills/do-work/actions/version.md"
}

func releaseVersion(path string, contents []byte) (string, bool) {
	value := strings.TrimSpace(string(contents))
	switch filepath.Base(path) {
	case "package.json":
		var document struct {
			Version string `json:"version"`
		}
		if json.Unmarshal(contents, &document) != nil {
			return "", false
		}
		value = document.Version
	case "package-lock.json":
		var document struct {
			Version  string `json:"version"`
			Packages map[string]struct {
				Version string `json:"version"`
			} `json:"packages"`
		}
		if json.Unmarshal(contents, &document) != nil || document.Version == "" {
			return "", false
		}
		if rootPackage, ok := document.Packages[""]; ok && rootPackage.Version != "" && rootPackage.Version != document.Version {
			return "", false
		}
		value = document.Version
	case "Cargo.toml", "pyproject.toml":
		sections := []string{"package"}
		if filepath.Base(path) == "pyproject.toml" {
			sections = []string{"project", "tool.poetry"}
		}
		sectionVersion, ok := tomlSectionVersion(contents, sections)
		if !ok {
			return "", false
		}
		value = sectionVersion
	case "Cargo.lock", "uv.lock":
		lockVersion, ok := projectLockVersion(filepath.Base(path), contents)
		if !ok {
			return "", false
		}
		value = lockVersion
	case "VERSION":
		// The whole trimmed file is the version.
	default:
		match := regexp.MustCompile(`(?m)^\*\*Current version\*\*:[ \t]*([0-9]+\.[0-9]+\.[0-9]+)[ \t]*$`).FindSubmatch(contents)
		if len(match) != 2 {
			return "", false
		}
		value = string(match[1])
	}
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return "", false
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return "", false
		}
	}
	return value, true
}

func tomlSectionVersion(contents []byte, acceptedSections []string) (string, bool) {
	accepted := map[string]bool{}
	for _, section := range acceptedSections {
		accepted[section] = true
	}
	currentSection := ""
	for _, line := range strings.Split(string(contents), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			currentSection = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "["), "]"))
			continue
		}
		if !accepted[currentSection] {
			continue
		}
		match := regexp.MustCompile(`^version[ \t]*=[ \t]*["']([0-9]+\.[0-9]+\.[0-9]+)["'][ \t]*$`).FindStringSubmatch(trimmed)
		if len(match) == 2 {
			return match[1], true
		}
	}
	return "", false
}

func configuredReleaseMetadataPaths(repositoryRoot, oldVersion string, changelogBefore []byte) []string {
	tracked, err := gitLines(repositoryRoot, "ls-files")
	if err != nil {
		return nil
	}
	paths := []string{}
	for _, path := range tracked {
		if strings.HasPrefix(path, "do-work/") || installedReleasePathForDiscovery(repositoryRoot, path) || !releaseMetadataPath(path) {
			continue
		}
		before, err := headFileImage(repositoryRoot, path)
		if err != nil || !before.Exists {
			continue
		}
		if strings.HasPrefix(filepath.Base(path), "CHANGELOG") {
			if bytes.Equal(before.Bytes, changelogBefore) {
				paths = append(paths, path)
			}
			continue
		}
		if lockMirrorConfigured(repositoryRoot, path, oldVersion) {
			paths = append(paths, path)
			continue
		}
		if version, ok := releaseVersion(path, before.Bytes); ok && version == oldVersion {
			paths = append(paths, path)
		}
	}
	return uniqueSorted(paths)
}

func lockMirrorConfigured(repositoryRoot, path, oldVersion string) bool {
	manifestName := ""
	switch filepath.Base(path) {
	case "Cargo.lock":
		manifestName = "Cargo.toml"
	case "uv.lock":
		manifestName = "pyproject.toml"
	default:
		return false
	}
	manifestPath := filepath.ToSlash(filepath.Join(filepath.Dir(path), manifestName))
	manifest, err := headFileImage(repositoryRoot, manifestPath)
	if err != nil || !manifest.Exists {
		return false
	}
	version, ok := releaseVersion(manifestPath, manifest.Bytes)
	return ok && version == oldVersion
}

func installedReleasePathForDiscovery(repositoryRoot, path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	for _, prefix := range []string{".claude/skills/", ".codex/skills/", "node_modules/", "vendor/", "vendored/", "generated/", ".generated/"} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	if strings.HasPrefix(lower, "skills/do-work/") {
		manifest, _ := headFileImage(repositoryRoot, "suite/modules.tsv")
		return !manifest.Exists
	}
	return false
}

func singleVersionReplacement(before, after []byte, oldVersion, newVersion string) bool {
	return bytes.Equal(bytes.Replace(before, []byte(oldVersion), []byte(newVersion), 1), after) && bytes.Count(before, []byte(oldVersion)) == 1
}

func semanticVersionReplacement(path string, before, after []byte, oldVersion, newVersion string) bool {
	switch filepath.Base(path) {
	case "package.json":
		return bytes.Equal(replaceJSONVersionValues(before, oldVersion, newVersion, 1), after)
	case "package-lock.json":
		var document struct {
			Packages map[string]json.RawMessage `json:"packages"`
		}
		if json.Unmarshal(before, &document) != nil {
			return false
		}
		count := 1
		if _, ok := document.Packages[""]; ok {
			count = 2
		}
		return bytes.Equal(replaceJSONVersionValues(before, oldVersion, newVersion, count), after)
	case "Cargo.lock", "uv.lock":
		return bytes.Equal(replaceProjectLockVersion(filepath.Base(path), before, oldVersion, newVersion), after)
	default:
		return singleVersionReplacement(before, after, oldVersion, newVersion)
	}
}

func projectLockVersion(lockName string, contents []byte) (string, bool) {
	versions := []string{}
	for _, blockLocation := range projectLockBlockLocations(contents) {
		block := contents[blockLocation[0]:blockLocation[1]]
		owned := lockName == "Cargo.lock" && !regexp.MustCompile(`(?m)^source[ \t]*=`).Match(block)
		if lockName == "uv.lock" {
			owned = regexp.MustCompile(`(?m)^source[ \t]*=[ \t]*\{[^{\n]*(editable|virtual)[ \t]*=[ \t]*"\."`).Match(block)
		}
		if !owned {
			continue
		}
		match := regexp.MustCompile(`(?m)^version[ \t]*=[ \t]*["']([0-9]+\.[0-9]+\.[0-9]+)["'][ \t]*$`).FindSubmatch(block)
		if len(match) == 2 {
			versions = append(versions, string(match[1]))
		}
	}
	return exactSingleVersion(versions)
}

func exactSingleVersion(versions []string) (string, bool) {
	if len(versions) != 1 {
		return "", false
	}
	return versions[0], true
}

func replaceProjectLockVersion(lockName string, contents []byte, oldVersion, newVersion string) []byte {
	blocks := projectLockBlockLocations(contents)
	result := append([]byte(nil), contents...)
	offset := 0
	replaced := false
	for _, blockLocation := range blocks {
		start, end := blockLocation[0]+offset, blockLocation[1]+offset
		block := result[start:end]
		owned := lockName == "Cargo.lock" && !regexp.MustCompile(`(?m)^source[ \t]*=`).Match(block)
		if lockName == "uv.lock" {
			owned = regexp.MustCompile(`(?m)^source[ \t]*=[ \t]*\{[^{\n]*(editable|virtual)[ \t]*=[ \t]*"\."`).Match(block)
		}
		if !owned {
			continue
		}
		updated := regexp.MustCompile(`(?m)(^version[ \t]*=[ \t]*["'])`+regexp.QuoteMeta(oldVersion)+`(["'][ \t]*$)`).ReplaceAll(block, []byte("${1}"+newVersion+"${2}"))
		if bytes.Equal(updated, block) {
			continue
		}
		result = append(append(append([]byte(nil), result[:start]...), updated...), result[end:]...)
		offset += len(updated) - len(block)
		replaced = true
	}
	if !replaced {
		return nil
	}
	return result
}

func projectLockBlockLocations(contents []byte) [][2]int {
	headers := regexp.MustCompile(`(?m)^\[\[package\]\][ \t]*\r?$`).FindAllIndex(contents, -1)
	blocks := make([][2]int, 0, len(headers))
	for index, header := range headers {
		end := len(contents)
		if index+1 < len(headers) {
			end = headers[index+1][0]
		}
		blocks = append(blocks, [2]int{header[0], end})
	}
	return blocks
}

func replaceJSONVersionValues(contents []byte, oldVersion, newVersion string, count int) []byte {
	pattern := regexp.MustCompile(`("version"[ \t\r\n]*:[ \t\r\n]*")` + regexp.QuoteMeta(oldVersion) + `(")`)
	replaced := append([]byte(nil), contents...)
	for index := 0; index < count; index++ {
		location := pattern.FindSubmatchIndex(replaced)
		if len(location) == 0 {
			return nil
		}
		updated := make([]byte, 0, len(replaced)-len(oldVersion)+len(newVersion))
		updated = append(updated, replaced[:location[3]]...)
		updated = append(updated, newVersion...)
		updated = append(updated, replaced[location[4]:]...)
		replaced = updated
	}
	return replaced
}

func singleInsertion(before, after []byte) ([]byte, bool) {
	heading := bytes.Index(before, []byte("## "))
	if heading < 0 || heading > len(after) || !bytes.Equal(before[:heading], after[:heading]) || !bytes.HasSuffix(after[heading:], before[heading:]) {
		return nil, false
	}
	insertedEnd := len(after) - len(before[heading:])
	if insertedEnd <= heading {
		return nil, false
	}
	return after[heading:insertedEnd], true
}

func sharedFinalizationPath(path string) bool {
	return path == "do-work" || strings.HasPrefix(path, "do-work/")
}

func ambiguousSharedRemainder(dirty map[string]bool) []string {
	blocked := []string{}
	for path := range dirty {
		if sharedFinalizationPath(path) || releaseMetadataPath(path) {
			blocked = append(blocked, path)
		}
	}
	return uniqueSorted(blocked)
}

func discoveredJournal(repositoryRoot string, candidate *discoveryCandidate) (*Journal, error) {
	requestID := candidate.request.TypedRecord.RequestID
	requestPath := requestRepositoryPath(candidate.request.RelativePath)
	journalPath, _, err := journalLocations(repositoryRoot, requestID)
	if err != nil {
		return nil, err
	}
	paths := uniqueSorted(candidate.effectivePaths)
	lifecyclePreimages, lifecyclePostimages, err := discoveredImages(repositoryRoot, candidate.lifecyclePaths)
	if err != nil {
		return nil, err
	}
	releasePreimages, releasePostimages, err := discoveredImages(repositoryRoot, candidate.releasePaths)
	if err != nil {
		return nil, err
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
		LifecyclePreimages: lifecyclePreimages, LifecyclePostimages: lifecyclePostimages,
		ReleasePreimages: releasePreimages, ReleasePostimages: releasePostimages,
		EffectiveCommitPaths: paths, PreparedHead: preparedHead, PreparedDiffSHA256: preparedDiff, Discovered: true,
		SoleReleaserAttributed: append([]string(nil), candidate.attributedPaths...),
	}
	if err := writeJournal(journal); err != nil {
		return nil, err
	}
	return journal, nil
}

func discoveredImages(repositoryRoot string, paths []string) ([]FileImage, []FileImage, error) {
	preimages := make([]FileImage, 0, len(paths))
	postimages := make([]FileImage, 0, len(paths))
	for _, path := range uniqueSorted(paths) {
		before, err := headFileImage(repositoryRoot, path)
		if err != nil {
			return nil, nil, err
		}
		after, err := currentImage(repositoryRoot, path)
		if err != nil {
			return nil, nil, err
		}
		preimages = append(preimages, before)
		postimages = append(postimages, after)
	}
	return preimages, postimages, nil
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
	collection := []string{"do-work-cli", "--format", "json", "uncommitted-inventory"}
	record := resultmodel.FinalizationResult{Phase: string(PhaseDiscoveryRefused), BlockedPaths: append([]string(nil), paths...), ReasonCodes: []string{code}, NextArgv: verification, VerificationArgv: verification, CollectionArgv: collection}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, AffectedPaths: paths, Evidence: []string{reason}, Fixability: resultmodel.FixabilityRefused,
		AutomationStopReason: "legacy finalization evidence is ambiguous", NextArgv: verification, VerificationArgv: verification,
	}}, Finalization: &record, Finalizations: []resultmodel.FinalizationResult{record}}
}
