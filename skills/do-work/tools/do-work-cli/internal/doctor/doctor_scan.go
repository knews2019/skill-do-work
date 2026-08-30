package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/cleanup"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

type ScanOptions struct {
	Now time.Time
}

var commitHashPattern = regexp.MustCompile(`^[0-9a-fA-F]{7,40}$`)

func ScanRepository(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, options ScanOptions) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	if snapshot == nil {
		result.Outcome = resultmodel.OutcomeFailure
		result.Findings = []resultmodel.CommandFinding{doctorFinding("DOCTOR-SNAPSHOT-MISSING", resultmodel.SeverityError, nil, nil,
			"repository snapshot is required", resultmodel.FixabilityManual, "no repository evidence was available", doctorArgv(), doctorJSONArgv())}
		return result
	}
	result.RepositoryRoot = snapshot.RepositoryRoot
	now := options.Now.UTC().Truncate(time.Second)
	if now.IsZero() {
		now = time.Now().UTC().Truncate(time.Second)
	}
	for _, collision := range snapshot.CollisionEntries {
		paths := make([]string, 0, len(collision.ClaimPaths))
		for _, path := range collision.ClaimPaths {
			paths = append(paths, repositoryRelative(snapshot.RepositoryRoot, path))
		}
		result.Findings = append(result.Findings, doctorFinding("REQUEST-ID-COLLISION", resultmodel.SeverityError, []string{collision.RequestID}, paths,
			fmt.Sprintf("%s is claimed by %d exact paths", collision.RequestID, len(paths)), resultmodel.FixabilityManual,
			"an ID collision cannot be resolved without choosing the authoritative record", doctorArgv(), doctorJSONArgv()))
	}
	for _, damaged := range snapshot.DamagedRecords {
		result.Findings = append(result.Findings, damagedRecordFinding(ctx, snapshot, damaged))
	}
	addCoreRequestFindings(&result, snapshot, now)
	_, timestampFindings := BuildTimestampPlan(ctx, snapshot, now)
	result.Findings = append(result.Findings, timestampFindings...)
	damagePaths := map[string]bool{}
	for _, damaged := range snapshot.DamagedRecords {
		damagePaths[damaged.AbsolutePath] = true
	}
	for _, warning := range snapshot.WarningMessages {
		if warningNamesDamagedPath(warning, damagePaths) || warningDuplicatesStatusFinding(snapshot, warning) {
			continue
		}
		result.Findings = append(result.Findings, doctorFinding("DOCTOR-INSPECTION-WARNING", resultmodel.SeverityWarning, nil, warningAffectedPaths(snapshot, warning), warning,
			resultmodel.FixabilityManual, "part of the repository could not be inspected completely", doctorArgv(), doctorJSONArgv()))
	}
	for _, strayPath := range snapshot.StrayRequestPaths {
		result.Findings = append(result.Findings, doctorFinding("STRAY-REQUEST", resultmodel.SeverityWarning, nil,
			[]string{repositoryRelative(snapshot.RepositoryRoot, strayPath)}, "REQ file is outside queue, working, and archive",
			resultmodel.FixabilityManual, "its intended pipeline location requires judgment", doctorArgv(), doctorJSONArgv()))
	}
	if !isGitRepository(ctx, snapshot.RepositoryRoot) {
		result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "GIT-DIVERGENCE-NOT-APPLICABLE", Reason: "repository has no Git history; read-only filesystem diagnosis completed"})
	} else {
		addGitDivergenceFindings(ctx, &result, snapshot)
	}
	sortFindings(result.Findings)
	if len(result.Findings) > 0 {
		result.Outcome = resultmodel.OutcomeFindings
	}
	return result
}

func addCoreRequestFindings(result *resultmodel.CommandResult, snapshot *repositorymodel.RepositorySnapshot, now time.Time) {
	followUps := map[string]bool{}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TypedRecord.AddendumTo != "" {
			followUps[requestFile.TypedRecord.AddendumTo] = true
		}
	}
	implementationOwners := map[string][]*repositorymodel.RequestFile{}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.ParsedDocument == nil {
			continue
		}
		record := requestFile.TypedRecord
		requestID := record.RequestID
		if requestID == "" {
			requestID = requestFile.FilenameID
		}
		body := string(requestFile.ParsedDocument.BodyBytes())
		if requestFile.TreeSection == "working" {
			addStuckWorkFinding(result, requestFile, requestID, body, now)
		}
		if requestFile.TreeSection == "archive" && schemanormalization.IsTerminalSuccess(record.RequestStatus) {
			if !hasImplementationEvidence(body) && record.BuilderDecidedValue != "true" && !strings.Contains(strings.ToLower(body), "no changes needed") {
				result.Findings = append(result.Findings, doctorFinding("HOLLOW-COMPLETION", resultmodel.SeverityError, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
					"terminal-success REQ has no implementation summary or explicit no-change evidence", resultmodel.FixabilityManual,
					"implementation evidence cannot be reconstructed mechanically", doctorArgv(), doctorJSONArgv()))
			}
			if (strings.Contains(body, "## Scope") || strings.Contains(body, "## Pre-Flight")) && !strings.Contains(body, "## Qualification") {
				result.Findings = append(result.Findings, doctorFinding("MISSING-QUALIFICATION", resultmodel.SeverityWarning, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
					"planned REQ has no Qualification section", resultmodel.FixabilityManual, "qualification evidence may have been skipped", doctorArgv(), doctorJSONArgv()))
			}
			for _, changedPath := range implementationPaths(body) {
				if !strings.HasPrefix(changedPath, "do-work/") {
					implementationOwners[changedPath] = append(implementationOwners[changedPath], requestFile)
				}
			}
		}
		if record.RequestStatus == "failed" {
			if record.ErrorTypeEvidence.OriginalValue == "" {
				result.Findings = append(result.Findings, doctorFinding("FAILED-UNCLASSIFIED", resultmodel.SeverityWarning, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
					"failed REQ has no explicit error_type", resultmodel.FixabilityManual, "failure classification requires review", doctorArgv(), doctorJSONArgv()))
			}
			if (record.ErrorTypeValue == "intent" || record.ErrorTypeValue == "spec" || record.ErrorTypeValue == "code") && !followUps[requestID] {
				result.Findings = append(result.Findings, doctorFinding("FAILED-WITHOUT-FOLLOW-UP", resultmodel.SeverityWarning, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
					"failed REQ has no addendum_to recovery path", resultmodel.FixabilityManual, "recovery scope and cancellation require judgment", doctorArgv(), doctorJSONArgv()))
			}
		}
		if record.OriginalStatus != "" && !record.StatusEvidence.IsRecognized {
			result.Findings = append(result.Findings, doctorFinding("INVALID-STATUS", resultmodel.SeverityWarning, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
				fmt.Sprintf("status %q is outside the schema vocabulary", record.OriginalStatus), resultmodel.FixabilityManual,
				"the intended pipeline state cannot be inferred safely", doctorArgv(), doctorJSONArgv()))
		}
		if requestFile.TreeSection == "queue" {
			if schemanormalization.IsStopped(record.RequestStatus) {
				result.Findings = append(result.Findings, doctorFinding("STRANDED-TERMINAL-REQUEST", resultmodel.SeverityWarning, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
					"terminal REQ remains in queue", resultmodel.FixabilityAutomatic, "diagnosis is read-only",
					[]string{"do-work-cli", "cleanup"}, []string{"do-work-cli", "cleanup", "--dry-run"}))
			}
			addStaleQueueFinding(result, requestFile, now)
		}
		if requestFile.TreeSection == "working" && schemanormalization.IsStopped(record.RequestStatus) {
			result.Findings = append(result.Findings, doctorFinding("STRANDED-TERMINAL-REQUEST", resultmodel.SeverityWarning, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
				"terminal REQ remains in working", resultmodel.FixabilityAutomatic, "diagnosis is read-only",
				[]string{"do-work-cli", "cleanup"}, []string{"do-work-cli", "cleanup", "--dry-run"}))
		}
	}
	addOrphanedUserRequestFindings(result, snapshot)
	for path, owners := range implementationOwners {
		requestIDs := []string{}
		userRequests := map[string]bool{}
		addendumLinked := false
		for _, owner := range owners {
			requestIDs = append(requestIDs, owner.TypedRecord.RequestID)
			userRequests[owner.TypedRecord.UserRequestID] = true
			if owner.TypedRecord.AddendumTo != "" {
				addendumLinked = true
			}
		}
		if len(userRequests) >= 3 {
			result.Findings = append(result.Findings, doctorFinding("SCOPE-HOTSPOT", resultmodel.SeverityWarning, requestIDs, []string{path},
				"implementation path is shared by three or more unrelated user requests", resultmodel.FixabilityManual,
				"architectural ownership requires judgment", doctorArgv(), doctorJSONArgv()))
		} else if len(owners) >= 2 && len(userRequests) == 1 && !addendumLinked {
			result.Findings = append(result.Findings, doctorFinding("SCOPE-OVERLAP", resultmodel.SeverityInfo, requestIDs, []string{path},
				"implementation path is shared by multiple same-UR REQs without an addendum relationship", resultmodel.FixabilityManual,
				"requirement decomposition needs review", doctorArgv(), doctorJSONArgv()))
		}
	}
}

func addStaleQueueFinding(result *resultmodel.CommandResult, requestFile *repositorymodel.RequestFile, now time.Time) {
	record := requestFile.TypedRecord
	infoThreshold := time.Duration(0)
	warningThreshold := time.Duration(0)
	baseText := record.CreatedAt
	code := ""
	if record.RequestStatus == "pending-answers" {
		infoThreshold, warningThreshold, code = 3*24*time.Hour, 7*24*time.Hour, "STALE-PENDING-ANSWERS"
	} else if record.RequestStatus == "blocked" {
		infoThreshold, warningThreshold, code = 7*24*time.Hour, 14*24*time.Hour, "STALE-BLOCKED"
		if blocked, found := record.FieldEvidenceByName["blocked_at"]; found {
			baseText = blocked.ScalarValue
		}
	}
	if infoThreshold == 0 {
		return
	}
	base, parseError := requestmodel.ParseTimestamp(baseText)
	if parseError == nil && now.Sub(base) >= infoThreshold {
		severity := resultmodel.SeverityInfo
		if now.Sub(base) > warningThreshold {
			severity = resultmodel.SeverityWarning
		}
		result.Findings = append(result.Findings, doctorFinding(code, severity, nonEmptyString(record.RequestID), []string{requestRepositoryPath(requestFile)},
			fmt.Sprintf("queue state has remained unchanged for %s", now.Sub(base).Round(time.Hour)), resultmodel.FixabilityManual,
			"the waiting condition needs confirmation", doctorArgv(), doctorJSONArgv()))
	}
}

func addStuckWorkFinding(result *resultmodel.CommandResult, requestFile *repositorymodel.RequestFile, requestID, body string, now time.Time) {
	record := requestFile.TypedRecord
	severity := resultmodel.SeverityWarning
	ageEvidence := "claimed_at is absent"
	shouldReport := record.ClaimedAt == ""
	if record.ClaimedAt != "" {
		claimedAt, parseError := requestmodel.ParseTimestamp(record.ClaimedAt)
		if parseError != nil {
			shouldReport = true
			ageEvidence = fmt.Sprintf("claimed_at=%q is unparseable", record.ClaimedAt)
		} else if age := now.Sub(claimedAt); age > time.Hour {
			shouldReport = true
			ageEvidence = fmt.Sprintf("claimed age=%s", age.Round(time.Minute))
			if age > 24*time.Hour {
				severity = resultmodel.SeverityError
			}
		}
	}
	if !shouldReport {
		return
	}
	evidence := fmt.Sprintf("%s; title=%s; route=%s; last_phase=%s", ageEvidence,
		fallbackDash(record.RequestTitle), fallbackDash(record.RouteValue), fallbackDash(lastMarkdownPhase(body)))
	result.Findings = append(result.Findings, doctorFinding("STUCK-WORK", severity, nonEmptyString(requestID), []string{requestRepositoryPath(requestFile)},
		evidence, resultmodel.FixabilityManual, "claimed work needs an ownership decision before reset",
		[]string{"git", "log", "--full-history", "--", requestRepositoryPath(requestFile)}, doctorJSONArgv()))
}

func lastMarkdownPhase(body string) string {
	last := ""
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "## ") {
			last = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "## "))
		}
	}
	return last
}

func hasImplementationEvidence(body string) bool {
	if !strings.Contains(body, "## Implementation Summary") {
		return false
	}
	for _, path := range implementationPaths(body) {
		if !strings.HasPrefix(path, "do-work/") {
			return true
		}
	}
	return false
}

func addOrphanedUserRequestFindings(result *resultmodel.CommandResult, snapshot *repositorymodel.RepositorySnapshot) {
	for _, userRequest := range snapshot.UserRequestFiles {
		if userRequest.TreeSection != "user-requests" || userRequest.ParsedDocument == nil {
			continue
		}
		userRequestID := userRequest.TypedRecord.RequestID
		allResolved := true
		for _, requestFile := range snapshot.RequestFiles {
			if requestFile.TypedRecord.UserRequestID == userRequestID && !schemanormalization.IsTerminalResolved(requestFile.TypedRecord.RequestStatus) {
				allResolved = false
				break
			}
		}
		if allResolved {
			result.Findings = append(result.Findings, doctorFinding("ORPHANED-USER-REQUEST", resultmodel.SeverityWarning, nonEmptyString(userRequestID), []string{"do-work/" + userRequest.RelativePath},
				"active UR has no unresolved member REQ", resultmodel.FixabilityAutomatic, "diagnosis is read-only",
				[]string{"do-work-cli", "cleanup"}, []string{"do-work-cli", "cleanup", "--dry-run"}))
		}
	}
}

func damagedRecordFinding(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, damaged repositorymodel.DamagedRecordFile) resultmodel.CommandFinding {
	evidence := []string{fmt.Sprintf("%s record has %d bytes and cannot parse: %s", damaged.RecordKind, len(damaged.ContentBytes), damaged.ParseFailure)}
	nextArgv := doctorArgv()
	fixability := resultmodel.FixabilityManual
	stopReason := "the record body must be recovered before field-level diagnosis"
	repositoryPath := "do-work/" + damaged.RelativePath
	if recovery, recoveryError := cleanup.RecoverGitContent(ctx, snapshot.RepositoryRoot, repositoryPath); recoveryError == nil {
		evidence = append(evidence, fmt.Sprintf("recoverable bytes=%d source_commit=%s implementation_commit=%s",
			len(recovery.ContentBytes), recovery.SourceCommit, fallbackDash(recovery.ImplementationCommit)))
		nextArgv = []string{"do-work-cli", "cleanup", "--restore-blanked", repositoryPath}
		fixability = resultmodel.FixabilityRefused
		stopReason = "restoration is destructive recovery and requires exact cleanup consent"
	} else {
		evidence = append(evidence, "full-history recovery unavailable: "+recoveryError.Error())
	}
	return doctorFinding("BLANKED-RECORD", resultmodel.SeverityError, nil, []string{repositoryPath}, strings.Join(evidence, "; "),
		fixability, stopReason, nextArgv, doctorJSONArgv())
}

func addGitDivergenceFindings(ctx context.Context, result *resultmodel.CommandResult, snapshot *repositorymodel.RepositorySnapshot) {
	checked := 0
	for requestIndex := len(snapshot.RequestFiles) - 1; requestIndex >= 0 && checked < 10; requestIndex-- {
		requestFile := snapshot.RequestFiles[requestIndex]
		if requestFile.TreeSection != "archive" || requestFile.ParsedDocument == nil {
			continue
		}
		commitField, found := requestFile.TypedRecord.FieldEvidenceByName["commit"]
		if !found || strings.TrimSpace(commitField.ScalarValue) == "" {
			continue
		}
		checked++
		if !commitHashPattern.MatchString(commitField.ScalarValue) {
			result.Findings = append(result.Findings, doctorFinding("GIT-DIVERGENCE-INSPECTION-FAILED", resultmodel.SeverityWarning,
				nonEmptyString(requestFile.TypedRecord.RequestID), []string{requestRepositoryPath(requestFile)}, "recorded implementation commit is not a 7-40 digit hexadecimal object name",
				resultmodel.FixabilityManual, "divergence inspection was incomplete", []string{"git", "log", "--", requestRepositoryPath(requestFile)}, doctorJSONArgv()))
			continue
		}
		for _, changedPath := range implementationPaths(string(requestFile.ParsedDocument.BodyBytes())) {
			if _, statError := os.Stat(filepath.Join(snapshot.RepositoryRoot, filepath.FromSlash(changedPath))); os.IsNotExist(statError) {
				result.Findings = append(result.Findings, doctorFinding("GIT-DIVERGENCE-MISSING-PATH", resultmodel.SeverityWarning,
					nonEmptyString(requestFile.TypedRecord.RequestID), []string{changedPath}, "implementation summary path no longer exists",
					resultmodel.FixabilityManual, "deletion intent cannot be inferred from history alone", []string{"git", "log", "--", changedPath}, doctorJSONArgv()))
				continue
			}
			command := exec.CommandContext(ctx, "git", "-C", snapshot.RepositoryRoot, "--literal-pathspecs", "log", "--format=%H", commitField.ScalarValue+"..HEAD", "--", changedPath)
			output, commandError := command.Output()
			if commandError != nil {
				result.Findings = append(result.Findings, doctorFinding("GIT-DIVERGENCE-INSPECTION-FAILED", resultmodel.SeverityWarning,
					nonEmptyString(requestFile.TypedRecord.RequestID), []string{changedPath}, "Git could not compare the implementation commit with HEAD",
					resultmodel.FixabilityManual, "divergence inspection was incomplete", []string{"git", "log", "--", changedPath}, doctorJSONArgv()))
			} else if strings.TrimSpace(string(output)) != "" {
				result.Findings = append(result.Findings, doctorFinding("GIT-DIVERGENCE-LATER-CHANGE", resultmodel.SeverityInfo,
					nonEmptyString(requestFile.TypedRecord.RequestID), []string{changedPath}, "path was modified after the recorded implementation commit",
					resultmodel.FixabilityManual, "later changes are informational and may be expected", []string{"git", "log", commitField.ScalarValue + "..HEAD", "--", changedPath}, doctorJSONArgv()))
			}
		}
	}
}

func implementationPaths(body string) []string {
	sectionIndex := strings.Index(body, "## Implementation Summary")
	if sectionIndex < 0 {
		return nil
	}
	section := body[sectionIndex+len("## Implementation Summary"):]
	if nextSection := strings.Index(section, "\n## "); nextSection >= 0 {
		section = section[:nextSection]
	}
	paths := []string{}
	seenPaths := map[string]bool{}
	for _, line := range strings.Split(section, "\n") {
		trimmedLine := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(trimmedLine, "- `") {
			continue
		}
		pathRemainder := strings.TrimPrefix(trimmedLine, "- `")
		closingBacktick := strings.IndexByte(pathRemainder, '`')
		if closingBacktick <= 0 {
			continue
		}
		path := filepath.ToSlash(filepath.Clean(pathRemainder[:closingBacktick]))
		if path != "." && !filepath.IsAbs(path) && !strings.HasPrefix(path, "../") && !seenPaths[path] {
			seenPaths[path] = true
			paths = append(paths, path)
		}
	}
	return paths
}

func doctorFinding(code string, severity resultmodel.FindingSeverity, ids, paths []string, evidence string,
	fixability resultmodel.FindingFixability, stopReason string, nextArgv, verifyArgv []string) resultmodel.CommandFinding {
	if fixability == resultmodel.FixabilityManual && len(paths) > 0 && sameArgv(nextArgv, doctorArgv()) {
		nextArgv = append([]string{"git", "log", "--full-history", "--"}, paths...)
	}
	return resultmodel.CommandFinding{Code: code, Severity: severity, AffectedIDs: ids, AffectedPaths: paths,
		Evidence: []string{evidence}, Fixability: fixability, AutomationStopReason: stopReason,
		NextArgv: nextArgv, VerificationArgv: verifyArgv}
}

func sameArgv(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func doctorArgv() []string     { return []string{"do-work-cli", "doctor"} }
func doctorJSONArgv() []string { return []string{"do-work-cli", "--format", "json", "doctor"} }
func repairDoctorArgv(dryRun, commit bool) []string {
	argv := []string{"do-work-cli", "doctor", "--repair-timestamps"}
	if dryRun {
		argv = append(argv, "--dry-run")
	}
	if commit {
		argv = append(argv, "--commit")
	}
	return argv
}

func nonEmptyString(value string) []string {
	if value == "" {
		return []string{}
	}
	return []string{value}
}

func fallbackDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func repositoryRelative(repositoryRoot, path string) string {
	relative, relativeError := filepath.Rel(repositoryRoot, path)
	if relativeError != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}

func requestRepositoryPath(requestFile *repositorymodel.RequestFile) string {
	return "do-work/" + requestFile.RelativePath
}

func warningNamesDamagedPath(warning string, paths map[string]bool) bool {
	for path := range paths {
		if strings.Contains(warning, path) {
			return true
		}
	}
	return false
}

func warningDuplicatesStatusFinding(snapshot *repositorymodel.RepositorySnapshot, warning string) bool {
	for _, requestFile := range snapshot.RequestFiles {
		statusWarning := requestFile.TypedRecord.StatusEvidence.WarningMessage
		if statusWarning != "" && strings.Contains(warning, requestFile.AbsolutePath) && strings.Contains(warning, statusWarning) {
			return true
		}
	}
	return false
}

func warningAffectedPaths(snapshot *repositorymodel.RepositorySnapshot, warning string) []string {
	paths := []string{}
	seen := map[string]bool{}
	addPath := func(absolutePath string) {
		if absolutePath != "" && strings.Contains(warning, absolutePath) {
			relativePath := repositoryRelative(snapshot.RepositoryRoot, absolutePath)
			if !seen[relativePath] {
				seen[relativePath] = true
				paths = append(paths, relativePath)
			}
		}
	}
	for _, requestFile := range snapshot.RequestFiles {
		addPath(requestFile.AbsolutePath)
	}
	for _, userRequest := range snapshot.UserRequestFiles {
		addPath(userRequest.AbsolutePath)
	}
	for _, damaged := range snapshot.DamagedRecords {
		addPath(damaged.AbsolutePath)
	}
	if len(paths) == 0 {
		if pathStart := strings.Index(warning, snapshot.DoWorkRoot); pathStart >= 0 {
			candidate := warning[pathStart:]
			if pathEnd := strings.Index(candidate, ": "); pathEnd >= 0 {
				candidate = candidate[:pathEnd]
			}
			paths = append(paths, repositoryRelative(snapshot.RepositoryRoot, candidate))
		}
	}
	sort.Strings(paths)
	return paths
}

func isGitRepository(ctx context.Context, repositoryRoot string) bool {
	return exec.CommandContext(ctx, "git", "-C", repositoryRoot, "rev-parse", "--git-dir").Run() == nil
}

func sortFindings(findings []resultmodel.CommandFinding) {
	sort.SliceStable(findings, func(leftIndex, rightIndex int) bool {
		leftPath, rightPath := "", ""
		if len(findings[leftIndex].AffectedPaths) > 0 {
			leftPath = findings[leftIndex].AffectedPaths[0]
		}
		if len(findings[rightIndex].AffectedPaths) > 0 {
			rightPath = findings[rightIndex].AffectedPaths[0]
		}
		if leftPath == rightPath {
			return findings[leftIndex].Code < findings[rightIndex].Code
		}
		return leftPath < rightPath
	})
}
