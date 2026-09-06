// Package repairvalidation owns the one executable decision authority for an
// already-green repository-gate repair. Work and review consume two projections
// from the same evidence instead of rebuilding the predicate in prose.
package repairvalidation

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gateevidence"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/sharedprimitives"
)

const CommandValidateAlreadyGreenRepair = "validate-already-green-repair"

var workingRequestPathPattern = regexp.MustCompile(`^do-work/working/REQ-[0-9]+(?:-[^/]+)?\.md$`)

type Options struct {
	RequestPath string
	WriterLabel string
	Now         time.Time
}

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{CommandValidateAlreadyGreenRepair: handleValidateAlreadyGreenRepair}
}

func handleValidateAlreadyGreenRepair(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseOptions(arguments)
	if err != nil {
		return failure(executionContext.RepositoryRoot, "REPAIR-VALIDATION-USAGE", err)
	}
	validation, evidence, err := Validate(executionContext.RepositoryRoot, options)
	if err != nil {
		result := failure(executionContext.RepositoryRoot, "REPAIR-VALIDATION-FAILED", err)
		result.AlreadyGreenRepair = &validation
		result.GateEvidence = evidence
		return result
	}
	return resultmodel.CommandResult{
		Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: executionContext.RepositoryRoot,
		AlreadyGreenRepair: &validation, GateEvidence: evidence,
	}
}

func parseOptions(arguments []string) (Options, error) {
	var options Options
	seen := map[string]bool{}
	for index := 0; index < len(arguments); index++ {
		name := arguments[index]
		if name != "--request-path" && name != "--writer" && name != "--at" {
			return options, fmt.Errorf("unknown %s option %q", CommandValidateAlreadyGreenRepair, name)
		}
		if seen[name] {
			return options, fmt.Errorf("%s may be supplied only once", name)
		}
		seen[name] = true
		index++
		if index >= len(arguments) || strings.TrimSpace(arguments[index]) == "" {
			return options, fmt.Errorf("%s requires a value", name)
		}
		switch name {
		case "--request-path":
			options.RequestPath = arguments[index]
		case "--writer":
			options.WriterLabel = arguments[index]
		case "--at":
			parsed, err := time.Parse(time.RFC3339, arguments[index])
			if err != nil {
				return options, fmt.Errorf("--at requires RFC3339: %w", err)
			}
			options.Now = parsed
		}
	}
	if options.RequestPath == "" || options.WriterLabel == "" {
		return options, fmt.Errorf("usage: %s --request-path <do-work/working/REQ-NNN-...md> --writer <label> [--at <RFC3339>]", CommandValidateAlreadyGreenRepair)
	}
	return options, nil
}

// Validate observes request evidence, recorded gate evidence, Git state, and a
// canonical completion dry run without mutating the repository.
func Validate(repositoryRoot string, options Options) (resultmodel.AlreadyGreenRepairValidation, *resultmodel.GateEvidenceResult, error) {
	if strings.TrimSpace(options.WriterLabel) == "" {
		return resultmodel.AlreadyGreenRepairValidation{}, nil, fmt.Errorf("completion writer is required")
	}
	if options.Now.IsZero() {
		options.Now = time.Now().UTC()
	}
	options.Now = options.Now.UTC().Truncate(time.Second)
	validation := resultmodel.AlreadyGreenRepairValidation{
		RequestPath: options.RequestPath, Writer: options.WriterLabel,
		PlannedAt: options.Now.Format(time.RFC3339),
	}
	if err := validateRequestPath(options.RequestPath); err != nil {
		return validation, nil, err
	}
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		return validation, nil, fmt.Errorf("discover repository: %w", err)
	}
	target := requestAtPath(snapshot, options.RequestPath)
	if target == nil {
		return validation, nil, fmt.Errorf("request path is not one discovered REQ: %s", options.RequestPath)
	}
	validation.RequestID = target.TypedRecord.RequestID
	coreReasons := map[string]bool{}
	addCore := func(code string) { coreReasons[code] = true }

	if target.TreeSection != "working" || target.ParseFailure != "" || target.FilenameID == "" || target.FilenameID != target.TypedRecord.RequestID || !requestIdentityIsUnique(snapshot, target.TypedRecord.RequestID) {
		addCore("REPAIR-REQUEST-IDENTITY")
	}
	if !exactScalar(target, "status", "claimed") {
		addCore("REPAIR-STATUS-NOT-CLAIMED")
	}
	if !exactScalar(target, "repository_gate_repair", "true") {
		addCore("REPAIR-MARKER-NOT-EXACT")
	}
	if !exactScalar(target, "tdd", "true") {
		addCore("REPAIR-TDD-NOT-ENABLED")
	}
	if _, found := target.TypedRecord.FieldEvidenceByName["release_at"]; found {
		addCore("REPAIR-RELEASE-METADATA-PRESENT")
	}

	body := ""
	if target.ParsedDocument != nil {
		body = string(target.ParsedDocument.BodyBytes())
	}
	sections := markdownSections(body)
	if hasPrefixedSection(sections, "Repository Gate Repair No-Op") {
		addCore("REPAIR-NO-OP-SHAPE")
	}
	intakes := sections["Repository Gate Repair Intake"]
	if len(intakes) == 0 {
		addCore("REPAIR-INTAKE-MISSING")
	}
	var intakeArgv []string
	for _, intake := range intakes {
		fingerprint, fingerprintOK := uniqueLabel(intake, "Diagnostic fingerprint")
		commandJSON, commandOK := uniqueLabel(intake, "Gate command (argv JSON)")
		directStatus, statusOK := uniqueLabel(intake, "Direct exit status")
		parent, parentOK := uniqueLabel(intake, "Parent")
		dependency, dependencyOK := uniqueLabel(intake, "Repair dependency")
		argv, argvOK := parseArgv(commandJSON)
		status, statusErr := strconv.Atoi(directStatus)
		if !fingerprintOK || strings.TrimSpace(fingerprint) == "" || !commandOK || !argvOK || !statusOK || statusErr != nil || status == 0 || !parentOK || strings.TrimSpace(parent) == "" || !dependencyOK || dependency != validation.RequestID {
			addCore("REPAIR-INTAKE-MALFORMED")
			continue
		}
		if validation.IntakeFingerprint == "" {
			validation.IntakeFingerprint = fingerprint
			intakeArgv = argv
		} else if validation.IntakeFingerprint != fingerprint || !equalStrings(intakeArgv, argv) {
			addCore("REPAIR-INTAKE-CONFLICT")
		}
	}

	noOps := sections["Repository Gate Repair No-Op"]
	if len(noOps) != 1 {
		addCore("REPAIR-NO-OP-SHAPE")
	} else {
		noOp := noOps[0]
		if !exactLabeledBlock(noOp, []string{"Expected diagnostic fingerprint", "Gate command", "Direct exit status", "Recorded green revision", "Observed result", "Verified at"}) {
			addCore("REPAIR-NO-OP-SHAPE")
		}
		expected, expectedOK := uniqueLabel(noOp, "Expected diagnostic fingerprint")
		commandJSON, commandOK := uniqueLabel(noOp, "Gate command")
		directStatus, statusOK := uniqueLabel(noOp, "Direct exit status")
		recordedRevision, revisionOK := uniqueLabel(noOp, "Recorded green revision")
		observed, observedOK := uniqueLabel(noOp, "Observed result")
		verifiedAt, verifiedOK := uniqueLabel(noOp, "Verified at")
		argv, argvOK := parseArgv(commandJSON)
		verifiedTime, timeErr := time.Parse(time.RFC3339, verifiedAt)
		verifiedCanonical := timeErr == nil && verifiedTime.UTC().Format(time.RFC3339) == verifiedAt
		validation.ExpectedFingerprint = expected
		validation.GateCommand = argv
		validation.RecordedRevision = recordedRevision
		if !expectedOK || strings.TrimSpace(expected) == "" || !commandOK || !argvOK || !statusOK || directStatus != "0" || !revisionOK || strings.TrimSpace(recordedRevision) == "" || !observedOK || observed != "green before implementation; repair already satisfied" || !verifiedOK || !verifiedCanonical {
			addCore("REPAIR-NO-OP-SHAPE")
		}
	}
	if validation.IntakeFingerprint == "" || validation.ExpectedFingerprint != validation.IntakeFingerprint {
		addCore("REPAIR-FINGERPRINT-MISMATCH")
	}
	if len(intakeArgv) == 0 || !equalStrings(intakeArgv, validation.GateCommand) {
		addCore("REPAIR-GATE-ARGV-MISMATCH")
	}

	const summary = "**Files changed:** None — verified repository-gate repair no-op.\n\n**What was done:** Re-ran the repair's recorded canonical repository gate before source edits and confirmed it is already green; no implementation changes were necessary."
	if values := sections["Implementation Summary"]; len(values) != 1 || strings.TrimSpace(values[0]) != summary {
		addCore("REPAIR-SUMMARY-SHAPE")
	}
	const qualification = "Passed — repository-gate repair no-op; durable gate evidence verified and project diff empty."
	if values := sections["Qualification"]; len(values) != 1 || strings.TrimSpace(values[0]) != qualification {
		addCore("REPAIR-QUALIFICATION-SHAPE")
	}

	var gateResult *resultmodel.GateEvidenceResult
	gateEvidenceMatches := false
	if len(validation.GateCommand) > 0 && validation.RecordedRevision != "" {
		checked, checkErr := gateevidence.CheckGreenGateAtRevision(repositoryRoot, validation.GateCommand, validation.RecordedRevision)
		gateResult = &checked
		if checkErr != nil || !checked.Matches || checked.RecordedRevision != validation.RecordedRevision {
			addCore("REPAIR-GATE-EVIDENCE-UNVERIFIABLE")
		} else {
			gateEvidenceMatches = true
		}
	} else {
		addCore("REPAIR-GATE-EVIDENCE-UNVERIFIABLE")
	}
	if gateEvidenceMatches {
		durableBytes, durableErr := gitFileAtRevision(repositoryRoot, gateResult.TargetRevision, options.RequestPath)
		if durableErr != nil {
			addCore("REPAIR-INTAKE-NOT-DURABLE")
		} else {
			durableSections := markdownSections(string(durableBytes))
			if !equalStrings(durableSections["Repository Gate Repair Intake"], intakes) || len(durableSections["Repository Gate Repair No-Op"]) != 0 {
				addCore("REPAIR-INTAKE-NOT-DURABLE")
			}
		}
	}

	changedPaths, err := gitStatusPaths(repositoryRoot)
	if err != nil {
		return validation, gateResult, fmt.Errorf("observe project status: %w", err)
	}
	for _, path := range changedPaths {
		if !strings.HasPrefix(path, "do-work/") {
			validation.ProjectChangedPaths = append(validation.ProjectChangedPaths, path)
			validation.OffendingPaths = append(validation.OffendingPaths, path)
		}
	}
	if len(validation.ProjectChangedPaths) > 0 {
		addCore("REPAIR-PROJECT-DIFF-NONEMPTY")
	}

	validation.TDDAllowed = len(coreReasons) == 0
	for code := range coreReasons {
		validation.ReasonCodes = append(validation.ReasonCodes, code)
	}

	plan := requeststate.BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), requeststate.StateOptions{
		Transition: requeststate.TransitionComplete, RequestID: validation.RequestID, RequestPath: options.RequestPath,
		TerminalStatus: "completed", WriterLabel: options.WriterLabel, Now: options.Now, DryRun: true,
	})
	dryRun := requeststate.ApplyPlan(context.Background(), plan)
	if dryRun.Outcome != resultmodel.OutcomeSuccess {
		validation.ReasonCodes = append(validation.ReasonCodes, "REPAIR-CANONICAL-COMPLETION-REFUSED")
	} else {
		for _, change := range dryRun.Changes {
			validation.CanonicalCompletionPaths = append(validation.CanonicalCompletionPaths, change.Path)
		}
		validation.CanonicalCompletionPaths = sharedprimitives.UniqueSortedStrings(validation.CanonicalCompletionPaths)
	}

	validation.StagedPaths, err = gitPaths(repositoryRoot, "diff", "--cached", "--name-only", "-z", "--no-renames")
	if err != nil {
		return validation, gateResult, fmt.Errorf("observe staged paths: %w", err)
	}
	allowed := make(map[string]bool, len(validation.CanonicalCompletionPaths))
	for _, path := range validation.CanonicalCompletionPaths {
		allowed[path] = true
	}
	for _, path := range validation.StagedPaths {
		if !allowed[path] {
			validation.OffendingPaths = append(validation.OffendingPaths, path)
			validation.ReasonCodes = append(validation.ReasonCodes, "REPAIR-STAGED-PATH-NOT-CANONICAL")
		}
	}
	validation.ReasonCodes = sharedprimitives.UniqueSortedStrings(validation.ReasonCodes)
	validation.OffendingPaths = sharedprimitives.UniqueSortedStrings(validation.OffendingPaths)
	validation.ReviewAllowed = validation.TDDAllowed && !containsString(validation.ReasonCodes, "REPAIR-CANONICAL-COMPLETION-REFUSED") && !containsString(validation.ReasonCodes, "REPAIR-STAGED-PATH-NOT-CANONICAL")
	return validation, gateResult, nil
}

func validateRequestPath(path string) error {
	cleaned := filepath.ToSlash(filepath.Clean(path))
	if filepath.IsAbs(path) || cleaned != path || !workingRequestPathPattern.MatchString(path) {
		return fmt.Errorf("request path must be one exact contained do-work/working/REQ-NNN-...md path")
	}
	return nil
}

func requestAtPath(snapshot *repositorymodel.RepositorySnapshot, requestPath string) *repositorymodel.RequestFile {
	for _, request := range snapshot.RequestFiles {
		if filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(request.RelativePath))) == requestPath {
			return request
		}
	}
	return nil
}

func requestIdentityIsUnique(snapshot *repositorymodel.RepositorySnapshot, requestID string) bool {
	if requestID == "" || len(snapshot.RequestsByID[requestID]) != 1 {
		return false
	}
	for _, collision := range snapshot.CollisionEntries {
		if collision.RequestID == requestID {
			return false
		}
	}
	return true
}

func exactScalar(request *repositorymodel.RequestFile, field, expected string) bool {
	evidence, found := request.TypedRecord.FieldEvidenceByName[field]
	return found && evidence.DuplicateCount == 0 && evidence.ScalarValue == expected
}

func markdownSections(body string) map[string][]string {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	lines := strings.Split(body, "\n")
	sections := map[string][]string{}
	name, start := "", 0
	flush := func(end int) {
		if name != "" {
			sections[name] = append(sections[name], strings.TrimSpace(strings.Join(lines[start:end], "\n")))
		}
	}
	for index, line := range lines {
		if strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "### ") {
			flush(index)
			name = strings.TrimPrefix(line, "## ")
			start = index + 1
		}
	}
	flush(len(lines))
	return sections
}

func hasPrefixedSection(sections map[string][]string, canonical string) bool {
	for name := range sections {
		if name != canonical && strings.HasPrefix(name, canonical) {
			return true
		}
	}
	return false
}

func exactLabeledBlock(section string, labels []string) bool {
	lines := strings.Split(section, "\n")
	if len(lines) != len(labels) {
		return false
	}
	for index, label := range labels {
		prefix := "- **" + label + ":** "
		if !strings.HasPrefix(lines[index], prefix) || strings.TrimSpace(strings.TrimPrefix(lines[index], prefix)) == "" || lines[index] != strings.TrimSpace(lines[index]) {
			return false
		}
	}
	return true
}

func uniqueLabel(section, label string) (string, bool) {
	prefix := "- **" + label + ":** "
	var values []string
	for _, line := range strings.Split(section, "\n") {
		if strings.HasPrefix(line, prefix) {
			values = append(values, strings.TrimPrefix(line, prefix))
		}
	}
	if len(values) != 1 {
		return "", false
	}
	if values[0] == "" || values[0] != strings.TrimSpace(values[0]) {
		return "", false
	}
	return values[0], true
}

func parseArgv(raw string) ([]string, bool) {
	var argv []string
	if err := json.Unmarshal([]byte(raw), &argv); err != nil || len(argv) == 0 {
		return nil, false
	}
	for _, argument := range argv {
		if argument == "" {
			return nil, false
		}
	}
	return argv, true
}

func gitStatusPaths(repositoryRoot string) ([]string, error) {
	output, err := exec.Command("git", "-C", repositoryRoot, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--no-renames").Output()
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, record := range strings.Split(string(output), "\x00") {
		if len(record) >= 4 {
			paths = append(paths, record[3:])
		}
	}
	return sharedprimitives.UniqueSortedStrings(paths), nil
}

func gitPaths(repositoryRoot string, arguments ...string) ([]string, error) {
	output, err := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...).Output()
	if err != nil {
		return nil, err
	}
	return sharedprimitives.UniqueSortedStrings(strings.FieldsFunc(string(output), func(r rune) bool { return r == 0 })), nil
}

func gitFileAtRevision(repositoryRoot, revision, path string) ([]byte, error) {
	if revision == "" {
		return nil, fmt.Errorf("recorded revision is empty")
	}
	output, err := exec.Command("git", "-C", repositoryRoot, "show", revision+":"+path).Output()
	if err != nil {
		return nil, fmt.Errorf("read durable repair intake at %s: %w", revision, err)
	}
	return output, nil
}

func equalStrings(left, right []string) bool {
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

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func failure(repositoryRoot, code string, err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{err.Error()}, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: "already-green repair eligibility is unverifiable",
	}}}
}
