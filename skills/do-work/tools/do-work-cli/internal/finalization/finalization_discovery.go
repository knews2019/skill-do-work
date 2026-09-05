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
	"github.com/knews2019/skill-do-work/do-work-cli/internal/releaseownership"
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

var enumerateTrackedReleasePaths = func(repositoryRoot string) ([]string, error) {
	return gitLines(repositoryRoot, "ls-files")
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
		if coherentClaimOnlyTopology(snapshot, dirty, origins, destinations) {
			return nil, nil
		}
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
		if sharedFinalizationPath(path) || releaseownership.IsReleaseMetadataPath(path) {
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
	releaseAmbiguous, releaseFailure := associateReleaseMetadata(repositoryRoot, dirty, candidates)
	if releaseFailure != nil {
		failure := discoveryRefusal(repositoryRoot, releaseFailure.code, releaseFailure.reason, releaseFailure.paths)
		return nil, &failure
	}
	ambiguous = append(ambiguous, releaseAmbiguous...)
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

// coherentClaimOnlyTopology distinguishes a claim transaction that stopped
// before its commit from an unjournalled archive/release tail. Recovery owns
// claims after finalization has explicitly declined this exact shape.
func coherentClaimOnlyTopology(snapshot *repositorymodel.RepositorySnapshot, dirty map[string]bool, origins, destinations map[string]string) bool {
	workingByBase := map[string]*repositorymodel.RequestFile{}
	for _, request := range snapshot.RequestFiles {
		if request.TreeSection != "working" || request.TypedRecord.RequestStatus != "claimed" {
			continue
		}
		workingByBase[filepath.Base(request.RelativePath)] = request
	}
	if len(workingByBase) == 0 || len(dirty) == 0 {
		return false
	}
	for path := range dirty {
		if !sharedFinalizationPath(path) && !releaseownership.IsReleaseMetadataPath(path) {
			continue
		}
		if path == "do-work/CHECKPOINT.md" {
			continue
		}
		base := filepath.Base(path)
		if _, found := workingByBase[base]; found && (strings.HasPrefix(path, "do-work/working/") || strings.HasPrefix(path, "do-work/queue/")) {
			continue
		}
		originBase := filepath.Base(origins[path])
		destinationBase := filepath.Base(destinations[path])
		_, originFound := workingByBase[originBase]
		destinationFound := workingByBase[destinationBase] != nil
		if (originBase != "." && originFound) || (destinationBase != "." && destinationFound) {
			continue
		}
		return false
	}
	return true
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
	appended := current.Bytes[len(before.Bytes):]
	return boundedFollowupFoldProves(appended, requestID)
}

func boundedFollowupFoldProves(appended []byte, requestID string) bool {
	headingPattern := regexp.MustCompile(`^\r?\n## (Review|Recovery) Fold — ` + regexp.QuoteMeta(requestID) + `[^\r\n]*\r?\n`)
	heading := headingPattern.FindSubmatch(appended)
	if len(heading) != 2 {
		return false
	}
	sectionHeadings := regexp.MustCompile(`(?m)^##[ \t]+`).FindAllIndex(appended, -1)
	if len(sectionHeadings) != 1 {
		return false
	}
	kind := strings.ToLower(string(heading[1]))
	markerPrefix := []byte("<!-- do-work:finalization-followup-fold-end ")
	if bytes.Count(appended, markerPrefix) != 1 {
		return false
	}
	marker := []byte("<!-- do-work:finalization-followup-fold-end kind=" + kind + " request=" + requestID + " -->")
	terminalMarker := regexp.MustCompile(`(?:^|\r?\n)` + regexp.QuoteMeta(string(marker)) + `\r?\n\z`)
	return terminalMarker.Match(appended)
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

type releaseDiscoveryFailure struct {
	code   string
	reason string
	paths  []string
}

type workspaceReleaseMember struct {
	relativePath string
	packageName  string
}

type workspaceReleaseMirror struct {
	kind       string
	rootCopies int
	members    []workspaceReleaseMember
}

type releaseVersionChange struct {
	path       string
	oldVersion string
	newVersion string
}

type configuredReleaseSet struct {
	paths            []string
	ownedPaths       map[string]bool
	workspaceMirrors map[string]workspaceReleaseMirror
}

func associateReleaseMetadata(repositoryRoot string, dirty map[string]bool, candidates []*discoveryCandidate) ([]string, *releaseDiscoveryFailure) {
	paths := []string{}
	for path := range dirty {
		if releaseownership.IsReleaseMetadataPath(path) {
			paths = append(paths, path)
		}
	}
	paths = uniqueSorted(paths)
	if len(paths) == 0 {
		return nil, nil
	}
	versions := map[string]struct{}{}
	oldVersions := map[string]struct{}{}
	changedSources := map[string]releaseVersionChange{}
	changelogPaths := []string{}
	var firstChangelogBefore, firstChangelogAfter []byte
	for _, path := range paths {
		before, _ := headFileImage(repositoryRoot, path)
		after, err := currentImage(repositoryRoot, path)
		if err != nil || !before.Exists || !after.Exists {
			return paths, nil
		}
		if strings.HasPrefix(filepath.Base(path), "CHANGELOG") {
			if firstChangelogBefore == nil {
				firstChangelogBefore, firstChangelogAfter = before.Bytes, after.Bytes
			} else if !bytes.Equal(firstChangelogBefore, before.Bytes) || !bytes.Equal(firstChangelogAfter, after.Bytes) {
				return paths, nil
			}
			changelogPaths = append(changelogPaths, path)
			continue
		}
		oldVersion, oldOK := releaseVersion(path, before.Bytes)
		newVersion, newOK := releaseVersion(path, after.Bytes)
		if !oldOK || !newOK || !semverGreater(newVersion, oldVersion) || !semanticVersionReplacement(path, before.Bytes, after.Bytes, oldVersion, newVersion) {
			continue
		}
		oldVersions[oldVersion] = struct{}{}
		versions[newVersion] = struct{}{}
		switch filepath.Base(path) {
		case "package.json", "Cargo.toml", "pyproject.toml":
			changedSources[path] = releaseVersionChange{path: path, oldVersion: oldVersion, newVersion: newVersion}
		}
	}
	if len(changelogPaths) == 0 || len(versions) != 1 || len(oldVersions) != 1 {
		return paths, nil
	}
	oldVersion := ""
	for version := range oldVersions {
		oldVersion = version
	}
	configured, err := configuredReleaseMetadataPaths(repositoryRoot, oldVersion, firstChangelogBefore, changedSources)
	if err != nil {
		return nil, &releaseDiscoveryFailure{code: "FINALIZATION-DISCOVERY-RELEASE-ENUMERATION", reason: "configured release-member enumeration failed closed: " + err.Error(), paths: paths}
	}
	unowned := []string{}
	for _, path := range paths {
		if !configured.ownedPaths[path] {
			unowned = append(unowned, path)
		}
	}
	if len(unowned) > 0 {
		return nil, &releaseDiscoveryFailure{code: "FINALIZATION-DISCOVERY-RELEASE-OWNERSHIP", reason: "release metadata lacks affirmative repository-owned or maintainer-mirror topology evidence", paths: uniqueSorted(unowned)}
	}
	configuredPaths := map[string]bool{}
	for _, path := range configured.paths {
		configuredPaths[path] = true
	}
	unassociated := []string{}
	for _, path := range paths {
		if !configuredPaths[path] {
			unassociated = append(unassociated, path)
		}
	}
	if len(unassociated) > 0 {
		return uniqueSorted(unassociated), nil
	}
	missingRequired := []string{}
	for _, requiredPath := range configured.paths {
		if !dirty[requiredPath] {
			missingRequired = append(missingRequired, requiredPath)
		}
	}
	if len(missingRequired) > 0 {
		return uniqueSorted(missingRequired), nil
	}
	newVersion := ""
	for version := range versions {
		newVersion = version
	}
	for _, path := range paths {
		if strings.HasPrefix(filepath.Base(path), "CHANGELOG") {
			continue
		}
		before, _ := headFileImage(repositoryRoot, path)
		after, _ := currentImage(repositoryRoot, path)
		if mirror, ok := configured.workspaceMirrors[path]; ok {
			if !workspaceMirrorReplacement(path, before.Bytes, after.Bytes, oldVersion, newVersion, mirror) {
				return paths, nil
			}
			continue
		}
		beforeVersion, beforeOK := releaseVersion(path, before.Bytes)
		afterVersion, afterOK := releaseVersion(path, after.Bytes)
		if !beforeOK || !afterOK || beforeVersion != oldVersion || afterVersion != newVersion || !semanticVersionReplacement(path, before.Bytes, after.Bytes, oldVersion, newVersion) {
			return paths, nil
		}
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
		return paths, nil
	}
	for _, path := range paths {
		matches[0].releasePaths = append(matches[0].releasePaths, path)
		matches[0].effectivePaths = append(matches[0].effectivePaths, path)
	}
	return nil, nil
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

func configuredReleaseMetadataPaths(repositoryRoot, oldVersion string, changelogBefore []byte, changedSources map[string]releaseVersionChange) (configuredReleaseSet, error) {
	tracked, err := enumerateTrackedReleasePaths(repositoryRoot)
	if err != nil {
		return configuredReleaseSet{}, err
	}
	trackedSet := map[string]bool{}
	for _, path := range tracked {
		trackedSet[filepath.ToSlash(filepath.Clean(path))] = true
	}
	ownership, err := releaseownership.AffirmativeOwnership(tracked, trackedSet, headReleaseImage(repositoryRoot))
	if err != nil {
		return configuredReleaseSet{}, err
	}
	ownedPaths, ownedManifests := ownership.MetadataPaths, ownership.ProjectManifests
	paths := []string{}
	for _, path := range tracked {
		path = filepath.ToSlash(filepath.Clean(path))
		if !ownedPaths[path] || !releaseownership.IsReleaseMetadataPath(path) {
			continue
		}
		before, err := headFileImage(repositoryRoot, path)
		if err != nil || !before.Exists {
			return configuredReleaseSet{}, fmt.Errorf("read tracked release member %s", path)
		}
		if strings.HasPrefix(filepath.Base(path), "CHANGELOG") {
			if bytes.Equal(before.Bytes, changelogBefore) {
				paths = append(paths, path)
			}
			continue
		}
		if filepath.Base(path) == "VERSION" || path == "skills/do-work/actions/version.md" {
			if version, ok := releaseVersion(path, before.Bytes); ok && version == oldVersion {
				paths = append(paths, path)
			}
			continue
		}
		if _, changed := changedSources[path]; changed {
			paths = append(paths, path)
		}
	}
	mirrors := map[string]workspaceReleaseMirror{}
	manifestPaths := make([]string, 0, len(changedSources))
	for manifestPath := range changedSources {
		manifestPaths = append(manifestPaths, manifestPath)
	}
	for _, manifestPath := range uniqueSorted(manifestPaths) {
		change := changedSources[manifestPath]
		base := filepath.Base(manifestPath)
		if !ownedManifests[manifestPath] || change.oldVersion != oldVersion {
			continue
		}
		manifest, readError := headFileImage(repositoryRoot, manifestPath)
		if readError != nil || !manifest.Exists {
			return configuredReleaseSet{}, fmt.Errorf("read tracked workspace member %s", manifestPath)
		}
		packageName, identityError := releasePackageName(base, manifest.Bytes)
		if identityError != nil {
			return configuredReleaseSet{}, fmt.Errorf("parse tracked workspace member identity %s: %w", manifestPath, identityError)
		}
		lockName, kind := workspaceLockName(base)
		rootLockPath := filepath.ToSlash(filepath.Join(filepath.Dir(manifestPath), lockName))
		if trackedSet[rootLockPath] {
			mirror := mirrors[rootLockPath]
			mirror.kind = kind
			if kind == "npm" {
				before, err := headFileImage(repositoryRoot, rootLockPath)
				if err != nil || !before.Exists {
					return configuredReleaseSet{}, fmt.Errorf("read tracked workspace lock %s", rootLockPath)
				}
				rootCopies, parseError := npmRootVersionCopies(before.Bytes, oldVersion)
				if parseError != nil {
					return configuredReleaseSet{}, fmt.Errorf("parse tracked workspace lock %s: %w", rootLockPath, parseError)
				}
				mirror.rootCopies = rootCopies
			} else {
				mirror.members = appendWorkspaceReleaseMember(mirror.members, workspaceReleaseMember{relativePath: ".", packageName: packageName})
			}
			if mirror.rootCopies > 0 || len(mirror.members) > 0 {
				mirrors[rootLockPath] = mirror
				paths = append(paths, rootLockPath)
			}
		}
		workspaceManifest, memberPath, ok := releaseownership.FindOwnedWorkspaceOwner(trackedSet, ownedManifests, manifestPath, base, headReleaseImage(repositoryRoot))
		if !ok {
			continue
		}
		lockPath := filepath.ToSlash(filepath.Join(filepath.Dir(workspaceManifest), lockName))
		if !trackedSet[lockPath] {
			continue
		}
		mirror := mirrors[lockPath]
		mirror.kind = kind
		mirror.members = appendWorkspaceReleaseMember(mirror.members, workspaceReleaseMember{relativePath: memberPath, packageName: packageName})
		mirrors[lockPath] = mirror
		paths = append(paths, lockPath)
	}
	return configuredReleaseSet{paths: uniqueSorted(paths), ownedPaths: ownedPaths, workspaceMirrors: mirrors}, nil
}

func appendWorkspaceReleaseMember(members []workspaceReleaseMember, member workspaceReleaseMember) []workspaceReleaseMember {
	for _, existing := range members {
		if existing == member {
			return members
		}
	}
	return append(members, member)
}

// releasePackageName returns the one identity a changed workspace source
// releases under. That name is what selects the source's lock entry, so an
// absent, unparseable, duplicated, or competing declaration is an error rather
// than a first-match guess.
func releasePackageName(manifestName string, contents []byte) (string, error) {
	if manifestName == "package.json" {
		var document struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(contents, &document); err != nil {
			return "", fmt.Errorf("manifest is not parseable JSON: %w", err)
		}
		if strings.TrimSpace(document.Name) == "" {
			return "", fmt.Errorf("manifest declares no name")
		}
		return document.Name, nil
	}
	sections := []string{"package"}
	if manifestName == "pyproject.toml" {
		sections = []string{"project", "tool.poetry"}
	}
	return tomlUniqueSectionScalar(contents, sections, "name")
}

func workspaceLockName(manifestName string) (string, string) {
	switch manifestName {
	case "package.json":
		return "package-lock.json", "npm"
	case "Cargo.toml":
		return "Cargo.lock", "cargo"
	default:
		return "uv.lock", "uv"
	}
}

// tomlUniqueSectionScalar returns the single value declared for key across the
// accepted sections. Every declaration counts wherever it appears among them,
// so a repeated key, a repeated section, and two dialect sections that disagree
// are all ambiguous and fail closed instead of resolving to the first match.
func tomlUniqueSectionScalar(contents []byte, sections []string, key string) (string, error) {
	accepted := map[string]bool{}
	for _, section := range sections {
		accepted[section] = true
	}
	currentSection := ""
	declarationPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(key) + `[ \t]*=`)
	valuePattern := regexp.MustCompile(`^` + regexp.QuoteMeta(key) + `[ \t]*=[ \t]*["']([^"']+)["'][ \t]*$`)
	values := []string{}
	for _, line := range strings.Split(string(contents), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			currentSection = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "["), "]"))
			continue
		}
		if !accepted[currentSection] || !declarationPattern.MatchString(trimmed) {
			continue
		}
		match := valuePattern.FindStringSubmatch(trimmed)
		if len(match) != 2 {
			return "", fmt.Errorf("%s declaration %q is not a quoted scalar", key, trimmed)
		}
		values = append(values, match[1])
	}
	switch len(values) {
	case 0:
		return "", fmt.Errorf("no %s declaration in section %s", key, strings.Join(sections, " or "))
	case 1:
		return values[0], nil
	default:
		return "", fmt.Errorf("%d competing %s declarations (%s)", len(values), key, strings.Join(values, ", "))
	}
}

// npmRootVersionCopies counts the root version copies a changed npm root must
// mirror in its lock: the top-level "version" and packages[""].version.
// Structural presence is the obligation, never the value, so a present copy
// holding anything other than the released old version fails closed instead of
// dropping out of the required mirror set and being left stale.
func npmRootVersionCopies(contents []byte, oldVersion string) (int, error) {
	var document struct {
		Version  json.RawMessage            `json:"version"`
		Packages map[string]json.RawMessage `json:"packages"`
	}
	if err := json.Unmarshal(contents, &document); err != nil {
		return 0, err
	}
	var rootPackageVersion json.RawMessage
	if rootPackage, present := document.Packages[""]; present {
		var rootPackageFields map[string]json.RawMessage
		if err := json.Unmarshal(rootPackage, &rootPackageFields); err != nil {
			return 0, fmt.Errorf(`packages[""] is not an object: %w`, err)
		}
		rootPackageVersion = rootPackageFields["version"]
	}
	rootVersionCopies := []struct {
		label string
		raw   json.RawMessage
	}{
		{label: `"version"`, raw: document.Version},
		{label: `packages[""]."version"`, raw: rootPackageVersion},
	}
	copies := 0
	for _, rootCopy := range rootVersionCopies {
		if rootCopy.raw == nil {
			continue
		}
		var value string
		if err := json.Unmarshal(rootCopy.raw, &value); err != nil {
			return 0, fmt.Errorf("root %s is not a string: %w", rootCopy.label, err)
		}
		if value != oldVersion {
			return 0, fmt.Errorf("root %s is %q, not the changed root's released version %q", rootCopy.label, value, oldVersion)
		}
		copies++
	}
	return copies, nil
}

func workspaceMirrorReplacement(path string, before, after []byte, oldVersion, newVersion string, mirror workspaceReleaseMirror) bool {
	expectedCopies := mirror.rootCopies + len(mirror.members)
	if expectedCopies == 0 {
		return false
	}
	switch mirror.kind {
	case "npm":
		var beforeLock struct {
			Version  string `json:"version"`
			Packages map[string]struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"packages"`
		}
		if json.Unmarshal(before, &beforeLock) != nil {
			return false
		}
		replacements := []exactByteReplacement{}
		rootReplacements := 0
		if mirror.rootCopies > 0 {
			if beforeLock.Version == oldVersion {
				location, ok := jsonNestedStringFieldLocation(before, nil, "version", oldVersion)
				if !ok {
					return false
				}
				replacements = append(replacements, exactByteReplacement{start: location[0], end: location[1], value: []byte(strconv.Quote(newVersion))})
				rootReplacements++
			}
			if rootPackage, ok := beforeLock.Packages[""]; ok && rootPackage.Version == oldVersion {
				location, found := jsonNestedStringFieldLocation(before, []string{"packages", ""}, "version", oldVersion)
				if !found {
					return false
				}
				replacements = append(replacements, exactByteReplacement{start: location[0], end: location[1], value: []byte(strconv.Quote(newVersion))})
				rootReplacements++
			}
		}
		if rootReplacements != mirror.rootCopies {
			return false
		}
		for _, member := range mirror.members {
			beforeMember, beforeOK := beforeLock.Packages[member.relativePath]
			if !beforeOK || beforeMember.Name != member.packageName || beforeMember.Version != oldVersion {
				return false
			}
			location, ok := jsonNestedStringFieldLocation(before, []string{"packages", member.relativePath}, "version", oldVersion)
			if !ok {
				return false
			}
			replacements = append(replacements, exactByteReplacement{start: location[0], end: location[1], value: []byte(strconv.Quote(newVersion))})
		}
		expected, ok := applyExactByteReplacements(before, replacements)
		return ok && bytes.Equal(expected, after)
	case "cargo", "uv":
		replacements := []exactByteReplacement{}
		for _, member := range mirror.members {
			location, ok := workspaceLockMemberVersionLocation(filepath.Base(path), before, member, oldVersion)
			if !ok {
				return false
			}
			replacements = append(replacements, exactByteReplacement{start: location[0], end: location[1], value: []byte(newVersion)})
		}
		expected, ok := applyExactByteReplacements(before, replacements)
		return ok && bytes.Equal(expected, after)
	default:
		return false
	}
}

type exactByteReplacement struct {
	start int
	end   int
	value []byte
}

func applyExactByteReplacements(contents []byte, replacements []exactByteReplacement) ([]byte, bool) {
	sort.Slice(replacements, func(left, right int) bool { return replacements[left].start > replacements[right].start })
	result := append([]byte(nil), contents...)
	previousStart := len(contents)
	for _, replacement := range replacements {
		if replacement.start < 0 || replacement.end < replacement.start || replacement.end > previousStart {
			return nil, false
		}
		result = append(append(append([]byte(nil), result[:replacement.start]...), replacement.value...), result[replacement.end:]...)
		previousStart = replacement.start
	}
	return result, len(replacements) > 0
}

func jsonNestedStringFieldLocation(contents []byte, objectPath []string, field, want string) ([2]int, bool) {
	if !json.Valid(contents) {
		return [2]int{}, false
	}
	objectStart, objectEnd := 0, len(contents)
	for _, key := range objectPath {
		location, ok := jsonObjectValueLocation(contents[objectStart:objectEnd], key)
		if !ok {
			return [2]int{}, false
		}
		objectEnd = objectStart + location[1]
		objectStart += location[0]
	}
	location, ok := jsonObjectValueLocation(contents[objectStart:objectEnd], field)
	if !ok {
		return [2]int{}, false
	}
	location[0] += objectStart
	location[1] += objectStart
	var value string
	if json.Unmarshal(contents[location[0]:location[1]], &value) != nil || value != want {
		return [2]int{}, false
	}
	return location, true
}

func jsonObjectValueLocation(object []byte, wantedKey string) ([2]int, bool) {
	index := skipJSONWhitespace(object, 0)
	if index >= len(object) || object[index] != '{' {
		return [2]int{}, false
	}
	index++
	found := [2]int{}
	foundCount := 0
	for {
		index = skipJSONWhitespace(object, index)
		if index >= len(object) {
			return [2]int{}, false
		}
		if object[index] == '}' {
			index = skipJSONWhitespace(object, index+1)
			if index != len(object) || foundCount != 1 {
				return [2]int{}, false
			}
			return found, true
		}
		keyEnd, ok := jsonStringEnd(object, index)
		if !ok {
			return [2]int{}, false
		}
		var key string
		if json.Unmarshal(object[index:keyEnd], &key) != nil {
			return [2]int{}, false
		}
		index = skipJSONWhitespace(object, keyEnd)
		if index >= len(object) || object[index] != ':' {
			return [2]int{}, false
		}
		valueStart := skipJSONWhitespace(object, index+1)
		valueEnd, ok := jsonValueEnd(object, valueStart)
		if !ok {
			return [2]int{}, false
		}
		if key == wantedKey {
			found = [2]int{valueStart, valueEnd}
			foundCount++
		}
		index = skipJSONWhitespace(object, valueEnd)
		if index < len(object) && object[index] == ',' {
			index++
			continue
		}
		if index >= len(object) || object[index] != '}' {
			return [2]int{}, false
		}
	}
}

func skipJSONWhitespace(contents []byte, index int) int {
	for index < len(contents) && (contents[index] == ' ' || contents[index] == '\t' || contents[index] == '\r' || contents[index] == '\n') {
		index++
	}
	return index
}

func jsonStringEnd(contents []byte, start int) (int, bool) {
	if start >= len(contents) || contents[start] != '"' {
		return 0, false
	}
	for index := start + 1; index < len(contents); index++ {
		switch contents[index] {
		case '\\':
			index++
		case '"':
			return index + 1, json.Valid(contents[start : index+1])
		}
	}
	return 0, false
}

func jsonValueEnd(contents []byte, start int) (int, bool) {
	if start >= len(contents) {
		return 0, false
	}
	if contents[start] == '"' {
		return jsonStringEnd(contents, start)
	}
	if contents[start] != '{' && contents[start] != '[' {
		end := start
		for end < len(contents) && contents[end] != ',' && contents[end] != '}' && contents[end] != ']' && contents[end] != ' ' && contents[end] != '\t' && contents[end] != '\r' && contents[end] != '\n' {
			end++
		}
		return end, end > start && json.Valid(contents[start:end])
	}
	stack := []byte{contents[start]}
	for index := start + 1; index < len(contents); index++ {
		if contents[index] == '"' {
			end, ok := jsonStringEnd(contents, index)
			if !ok {
				return 0, false
			}
			index = end - 1
			continue
		}
		switch contents[index] {
		case '{', '[':
			stack = append(stack, contents[index])
		case '}', ']':
			open := stack[len(stack)-1]
			if open == '{' && contents[index] != '}' || open == '[' && contents[index] != ']' {
				return 0, false
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return index + 1, json.Valid(contents[start : index+1])
			}
		}
	}
	return 0, false
}

func workspaceLockMemberVersionLocation(lockName string, contents []byte, member workspaceReleaseMember, version string) ([2]int, bool) {
	locations := [][2]int{}
	namePattern := regexp.MustCompile(`(?m)^name[ \t]*=[ \t]*["']([^"']+)["'][ \t]*\r?$`)
	versionPattern := regexp.MustCompile(`(?m)^version[ \t]*=[ \t]*["']([0-9]+\.[0-9]+\.[0-9]+)["'][ \t]*\r?$`)
	for _, location := range projectLockBlockLocations(contents) {
		block := contents[location[0]:location[1]]
		nameMatch := namePattern.FindSubmatch(block)
		versionMatch := versionPattern.FindSubmatchIndex(block)
		if len(nameMatch) != 2 || len(versionMatch) < 4 || string(nameMatch[1]) != member.packageName {
			continue
		}
		if lockName == "Cargo.lock" {
			if regexp.MustCompile(`(?m)^source[ \t]*=`).Match(block) {
				continue
			}
		} else {
			wantSource := regexp.MustCompile(`(?m)^source[ \t]*=[ \t]*\{[^\r\n}]*(editable|virtual)[ \t]*=[ \t]*["']` + regexp.QuoteMeta(member.relativePath) + `["']`)
			if !wantSource.Match(block) {
				continue
			}
		}
		if string(block[versionMatch[2]:versionMatch[3]]) != version {
			return [2]int{}, false
		}
		locations = append(locations, [2]int{location[0] + versionMatch[2], location[0] + versionMatch[3]})
	}
	if len(locations) != 1 {
		return [2]int{}, false
	}
	return locations[0], true
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
		if sharedFinalizationPath(path) || releaseownership.IsReleaseMetadataPath(path) {
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

// headReleaseImage adapts the committed image to the shared ownership
// reader. Recovery classifies work that is already committed, so HEAD is
// the image its ownership judgment must be made against.
func headReleaseImage(repositoryRoot string) releaseownership.ReadImage {
	return func(path string) ([]byte, bool) {
		image, err := headFileImage(repositoryRoot, path)
		if err != nil || !image.Exists {
			return nil, false
		}
		return image.Bytes, true
	}
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
	record := resultmodel.FinalizationResult{Phase: string(PhaseDiscoveryRefused), BlockedPaths: append([]string(nil), paths...), ReasonCodes: []string{code}, NextArgv: collection, VerificationArgv: verification, CollectionArgv: collection}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, AffectedPaths: paths, Evidence: []string{reason}, Fixability: resultmodel.FixabilityRefused,
		AutomationStopReason: "legacy finalization evidence is ambiguous", NextArgv: collection, VerificationArgv: verification,
	}}, Finalization: &record, Finalizations: []resultmodel.FinalizationResult{record}}
}
