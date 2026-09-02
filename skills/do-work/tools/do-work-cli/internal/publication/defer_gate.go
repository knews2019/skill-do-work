package publication

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var (
	diagnosticFingerprintPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._:-]*$`)
	sweepKeyPattern              = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	checkpointClaimPattern       = regexp.MustCompile(`^\s*-\s+(REQ-0*[0-9]+)\s*:.*\s+—\s+writer:\s*(\S(?:.*\S)?)\s*$`)
)

func BuildDeferGatePlan(repositoryRoot string, manifest Manifest) PublicationPlan {
	plan := PublicationPlan{Operation: OperationDeferGate, RepositoryRoot: repositoryRoot, CommitMessage: manifest.CommitMessage}
	gate := manifest.DeferGate
	if gate == nil {
		return refusedPlan(plan, "DEFER-GATE-MANIFEST-MISSING", "defer_gate body is required", nil)
	}
	if refusal := validateDeferGateManifest(plan, gate); refusal != nil {
		return *refusal
	}
	parentPath, _ := containedPath(gate.ParentPath)
	checkpointPath, _ := containedPath(gate.CheckpointPath)
	repairPath, _ := containedPath(gate.RepairPath)
	reservationPath, _ := containedPath(gate.ReservationPath)
	queueParentPath := "do-work/queue/" + filepath.Base(parentPath)
	if _, statError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(queueParentPath))); statError == nil {
		return refusedPlan(plan, "DEFER-GATE-PARENT-DESTINATION-COLLISION", "parent queue destination already exists", []string{gate.ParentID}, queueParentPath)
	} else if !os.IsNotExist(statError) {
		return refusedPlan(plan, "DEFER-GATE-INSPECTION-FAILED", statError.Error(), []string{gate.ParentID}, queueParentPath)
	}

	parentBytes, _, payloadError := readPayload(repositoryRoot, gate.ExpectedParent)
	if payloadError != nil {
		return refusedPlan(plan, "DEFER-GATE-PARENT-PREIMAGE-INVALID", payloadError.Error(), []string{gate.ParentID}, gate.ExpectedParent.SourcePath)
	}
	if staleReason := exactCurrentBytes(repositoryRoot, parentPath, parentBytes); staleReason != "" {
		return refusedPlan(plan, "DEFER-GATE-PARENT-STALE", staleReason, []string{gate.ParentID}, parentPath)
	}
	parentDocument, parseError := requestmodel.ParseDocument(parentBytes)
	if parseError != nil {
		return refusedPlan(plan, "DEFER-GATE-PARENT-INVALID", parseError.Error(), []string{gate.ParentID}, parentPath)
	}
	parentRecord := parentDocument.TypedRecord()
	if parentRecord.RequestID != gate.ParentID || parentRecord.RequestStatus != gate.ExpectedStatus || !userRequestPattern.MatchString(parentRecord.UserRequestID) {
		return refusedPlan(plan, "DEFER-GATE-PARENT-STATE-MISMATCH", "parent id, status, or user_request does not match the manifest", []string{gate.ParentID}, parentPath)
	}
	checkpointBytes, _, checkpointPayloadError := readPayload(repositoryRoot, gate.ExpectedCheckpoint)
	if checkpointPayloadError != nil {
		return refusedPlan(plan, "DEFER-GATE-CHECKPOINT-PREIMAGE-INVALID", checkpointPayloadError.Error(), []string{gate.ParentID}, gate.ExpectedCheckpoint.SourcePath)
	}
	if staleReason := exactCurrentBytes(repositoryRoot, checkpointPath, checkpointBytes); staleReason != "" {
		return refusedPlan(plan, "DEFER-GATE-CHECKPOINT-STALE", staleReason, []string{gate.ParentID}, checkpointPath)
	}
	newCheckpointBytes, claimCount := checkpointWithoutOwnedClaim(checkpointBytes, gate.ParentID, gate.WriterLabel)
	if claimCount != 1 {
		return refusedPlan(plan, "DEFER-GATE-CHECKPOINT-CLAIM-MISMATCH", fmt.Sprintf("expected one exact writer claim, found %d", claimCount), []string{gate.ParentID}, checkpointPath)
	}

	baseCommit, mergeCommit, mergeError := validateDeferredMergeRange(repositoryRoot, gate.DeferredImplementationBase, gate.DeferredImplementationMerge)
	if mergeError != nil {
		return refusedPlan(plan, "DEFER-GATE-MERGE-EVIDENCE-INVALID", mergeError.Error(), []string{gate.ParentID}, parentPath)
	}

	snapshot, discoveryError := repositorymodel.DiscoverRepository(repositoryRoot)
	if discoveryError != nil {
		return refusedPlan(plan, "DEFER-GATE-DISCOVERY-FAILED", discoveryError.Error(), []string{gate.ParentID}, "do-work")
	}
	matchingRepairs := matchingGateRepairs(snapshot, gate.SweepKey, gate.DiagnosticFingerprint)
	if len(matchingRepairs) > 1 {
		paths := make([]string, 0, len(matchingRepairs))
		for _, repair := range matchingRepairs {
			paths = append(paths, "do-work/"+repair.RelativePath)
		}
		return refusedPlan(plan, "DEFER-GATE-REPAIR-AMBIGUOUS", "multiple repair requests carry the same sweep identity and diagnostic fingerprint", []string{gate.RepairID}, paths...)
	}

	repairOutcome := "created"
	var repairMutation PlannedMutation
	if len(matchingRepairs) == 1 {
		repairOutcome = "folded"
		repair := matchingRepairs[0]
		exactRepairPath := "do-work/" + filepath.ToSlash(repair.RelativePath)
		if requestID(repair) != gate.RepairID || exactRepairPath != repairPath || repair.TreeSection != "queue" || repair.TypedRecord.RequestStatus != "pending" {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-MISMATCH", "the unique fold candidate does not match the pending queue repair named by the manifest", []string{gate.RepairID}, exactRepairPath, repairPath)
		}
		if gate.ExpectedRepair == nil {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-PREIMAGE-MISSING", "fold mode requires expected_repair", []string{gate.RepairID}, repairPath)
		}
		expectedRepairBytes, _, expectedRepairError := readPayload(repositoryRoot, *gate.ExpectedRepair)
		if expectedRepairError != nil || !bytes.Equal(expectedRepairBytes, repair.ContentBytes) || exactCurrentBytes(repositoryRoot, repairPath, expectedRepairBytes) != "" {
			reason := "repair preimage is stale"
			if expectedRepairError != nil {
				reason = expectedRepairError.Error()
			}
			return refusedPlan(plan, "DEFER-GATE-REPAIR-STALE", reason, []string{gate.RepairID}, repairPath)
		}
		repairDocument, _ := requestmodel.ParseDocument(expectedRepairBytes)
		related := appendUnique(repairDocument.TypedRecord().RelatedIDs, gate.ParentID)
		if listError := repairDocument.SetList("related", related); listError != nil {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-INVALID", listError.Error(), []string{gate.RepairID}, repairPath)
		}
		if instanceError := appendSweepInstance(repairDocument, repairInstanceLine(gate, parentRecord.UserRequestID)); instanceError != nil {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-INVALID", instanceError.Error(), []string{gate.RepairID}, repairPath)
		}
		appendHistory(repairDocument, repairFoldHistory(gate, baseCommit, mergeCommit))
		repairMutation = PlannedMutation{Kind: MutationReplace, Path: repairPath, ExpectedBytes: expectedRepairBytes, Contents: repairDocument.DocumentBytes(), AllowUntracked: !gitPathTracked(repositoryRoot, repairPath)}
		reservationBytes, reservationError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(reservationPath)))
		if reservationError == nil && !bytes.Equal(reservationBytes, []byte(gate.RepairID+"\n")) {
			return refusedPlan(plan, "DEFER-GATE-RESERVATION-STALE", "current bytes do not match the exact reservation", []string{gate.RepairID}, reservationPath)
		}
		if os.IsNotExist(reservationError) && !gitPathCleanAgainstHead(repositoryRoot, repairPath) {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-COMMIT-AUTHORITY-MISSING", "an absent reservation requires a committed clean repair preimage", []string{gate.RepairID}, repairPath, reservationPath)
		}
		if reservationError != nil && !os.IsNotExist(reservationError) {
			return refusedPlan(plan, "DEFER-GATE-RESERVATION-STALE", reservationError.Error(), []string{gate.RepairID}, reservationPath)
		}
	} else {
		if gate.ExpectedRepair != nil {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-PREIMAGE-UNEXPECTED", "create mode must not supply expected_repair", []string{gate.RepairID}, repairPath)
		}
		for _, requestFile := range snapshot.RequestsByID[gate.RepairID] {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-COLLISION", "repair id is already claimed by another request", []string{gate.RepairID}, "do-work/"+requestFile.RelativePath)
		}
		for _, targetPath := range []string{repairPath, reservationPath} {
			if _, statError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(targetPath))); statError == nil {
				return refusedPlan(plan, "DEFER-GATE-REPAIR-COLLISION", "repair or reservation destination already exists", []string{gate.RepairID}, targetPath)
			} else if !os.IsNotExist(statError) {
				return refusedPlan(plan, "DEFER-GATE-INSPECTION-FAILED", statError.Error(), []string{gate.RepairID}, targetPath)
			}
		}
		repairBytes, authorError := authoredRepairBytes(gate, parentRecord.UserRequestID, baseCommit, mergeCommit)
		if authorError != nil {
			return refusedPlan(plan, "DEFER-GATE-REPAIR-INVALID", authorError.Error(), []string{gate.RepairID}, repairPath)
		}
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationCreate, Path: reservationPath, Contents: []byte(gate.RepairID + "\n"), Mode: 0o644})
		repairMutation = PlannedMutation{Kind: MutationCreate, Path: repairPath, Contents: repairBytes, Mode: 0o644}
	}

	parentDependencies := appendUnique(parentRecord.DependsOn, gate.RepairID)
	if editError := parentDocument.SetScalar("status", "pending"); editError != nil {
		return refusedPlan(plan, "DEFER-GATE-PARENT-INVALID", editError.Error(), []string{gate.ParentID}, parentPath)
	}
	_ = parentDocument.DeleteField("claimed_at")
	_ = parentDocument.SetScalar("gate_deferred", "true")
	_ = parentDocument.SetList("depends_on", parentDependencies)
	if baseCommit == "" {
		_ = parentDocument.DeleteField("deferred_implementation_base")
		_ = parentDocument.DeleteField("deferred_implementation_merge")
	} else {
		_ = parentDocument.SetScalar("deferred_implementation_base", baseCommit)
		_ = parentDocument.SetScalar("deferred_implementation_merge", mergeCommit)
	}
	appendHistory(parentDocument, parentDeferralHistory(gate, baseCommit, mergeCommit))

	plan.Mutations = append(plan.Mutations,
		repairMutation,
		PlannedMutation{Kind: MutationMove, Path: parentPath, DestinationPath: queueParentPath, ExpectedBytes: parentBytes, Contents: parentDocument.DocumentBytes()},
		PlannedMutation{Kind: MutationReplace, Path: checkpointPath, ExpectedBytes: checkpointBytes, Contents: newCheckpointBytes},
	)
	plan.GateDeferral = &resultmodel.GateDeferralResult{
		ParentID: gate.ParentID, ParentPath: queueParentPath, RepairID: gate.RepairID, RepairPath: repairPath,
		CheckpointPath: checkpointPath, RepairOutcome: repairOutcome, RepairDependency: gate.RepairID,
		DiagnosticFingerprint: gate.DiagnosticFingerprint, SweepKey: gate.SweepKey,
		GateCommand: append([]string(nil), gate.GateCommand...), GateExitStatus: gate.GateExitStatus,
		DeferredImplementationBase: baseCommit, DeferredImplementationMerge: mergeCommit,
	}
	plan = finalizePlan(plan)
	if plan.Refusal != nil {
		return plan
	}
	preimagePaths := []string{parentPath, checkpointPath}
	if repairOutcome == "folded" {
		preimagePaths = append(preimagePaths, repairPath)
	}
	dirtyInputs, untrackedInputs, classificationError := classifyGatePreimages(repositoryRoot, preimagePaths...)
	if classificationError != nil {
		return refusedPlan(plan, "DEFER-GATE-PREIMAGE-CLASSIFICATION-FAILED", classificationError.Error(), []string{gate.ParentID}, preimagePaths...)
	}
	plan.ExistingDirtyTargetPaths = dirtyInputs
	plan.ExistingUntrackedTargetPaths = appendUniqueSorted(plan.ExistingUntrackedTargetPaths, untrackedInputs...)
	directories, topologyError := planCreatedDirectories(repositoryRoot, plan.TargetPaths)
	if topologyError != nil {
		return refusedPlan(plan, "DEFER-GATE-TOPOLOGY-UNSAFE", topologyError.Error(), []string{gate.ParentID, gate.RepairID}, "do-work")
	}
	plan.CreatedDirectoryPaths = directories
	return plan
}

func validateDeferGateManifest(plan PublicationPlan, gate *DeferGateManifest) *PublicationPlan {
	refuse := func(code, reason string, ids []string, paths ...string) *PublicationPlan {
		refused := refusedPlan(plan, code, reason, ids, paths...)
		return &refused
	}
	if !requestPattern.MatchString(gate.ParentID) || !requestPattern.MatchString(gate.RepairID) || gate.ParentID == gate.RepairID {
		return refuse("DEFER-GATE-ID-INVALID", "parent and repair ids must be distinct canonical REQ ids", []string{gate.ParentID, gate.RepairID})
	}
	parentPath, parentError := containedPath(gate.ParentPath)
	repairPath, repairError := containedPath(gate.RepairPath)
	checkpointPath, checkpointError := containedPath(gate.CheckpointPath)
	reservationPath, reservationError := containedPath(gate.ReservationPath)
	if parentError != nil || filepath.ToSlash(filepath.Dir(parentPath)) != "do-work/working" || !validRequestBasename(parentPath, gate.ParentID) {
		return refuse("DEFER-GATE-PARENT-PATH-INVALID", "parent path must be its exact slugged file in do-work/working", []string{gate.ParentID}, gate.ParentPath)
	}
	if repairError != nil || filepath.ToSlash(filepath.Dir(repairPath)) != "do-work/queue" || !validRequestBasename(repairPath, gate.RepairID) {
		return refuse("DEFER-GATE-REPAIR-PATH-INVALID", "repair path must be its exact slugged file in do-work/queue", []string{gate.RepairID}, gate.RepairPath)
	}
	if checkpointError != nil || checkpointPath != "do-work/CHECKPOINT.md" {
		return refuse("DEFER-GATE-CHECKPOINT-PATH-INVALID", "checkpoint path must be do-work/CHECKPOINT.md", []string{gate.ParentID}, gate.CheckpointPath)
	}
	if reservationError != nil || reservationPath != "do-work/.req-reservations/"+gate.RepairID {
		return refuse("DEFER-GATE-RESERVATION-PATH-INVALID", "reservation path must exactly match the repair id", []string{gate.RepairID}, gate.ReservationPath)
	}
	if gate.ExpectedStatus != "claimed" {
		return refuse("DEFER-GATE-STATUS-INVALID", "expected_status must be claimed", []string{gate.ParentID}, gate.ParentPath)
	}
	if gate.ExpectedParent.SourcePath == "" || gate.ExpectedCheckpoint.SourcePath == "" || strings.TrimSpace(gate.WriterLabel) == "" || gate.GateExitStatus < 1 || gate.GateExitStatus > 255 || len(gate.GateCommand) == 0 || len(gate.DiagnosticEvidence) == 0 {
		return refuse("DEFER-GATE-EVIDENCE-MISSING", "exact preimages, writer, non-zero gate status, and command argv are required", []string{gate.ParentID}, gate.ParentPath)
	}
	textEvidence := []string{gate.WriterLabel, gate.DiagnosticFingerprint, gate.SweepKey, gate.RepairTitle, gate.RepairCreatedAt}
	textEvidence = append(textEvidence, gate.GateCommand...)
	textEvidence = append(textEvidence, gate.DiagnosticEvidence...)
	for _, value := range textEvidence {
		if strings.TrimSpace(value) == "" || strings.ContainsAny(value, "\x00\r\n") {
			return refuse("DEFER-GATE-EVIDENCE-INVALID", "manifest text evidence must be non-empty single-line values", []string{gate.ParentID, gate.RepairID})
		}
	}
	if !diagnosticFingerprintPattern.MatchString(gate.DiagnosticFingerprint) || !sweepKeyPattern.MatchString(gate.SweepKey) {
		return refuse("DEFER-GATE-IDENTITY-INVALID", "diagnostic_fingerprint and sweep_key must be stable machine identifiers", []string{gate.ParentID, gate.RepairID})
	}
	createdAt, createdError := requestmodel.ParseTimestamp(gate.RepairCreatedAt)
	if createdError != nil || requestmodel.CanonicalTimestamp(createdAt) != gate.RepairCreatedAt {
		return refuse("DEFER-GATE-CREATED-AT-INVALID", "repair_created_at must use canonical UTC RFC3339", []string{gate.RepairID})
	}
	if (gate.DeferredImplementationBase == "") != (gate.DeferredImplementationMerge == "") {
		return refuse("DEFER-GATE-MERGE-EVIDENCE-UNPAIRED", "deferred implementation base and merge must be supplied together", []string{gate.ParentID})
	}
	return nil
}

func validRequestBasename(path, requestID string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, requestID+"-") && strings.HasSuffix(base, ".md") && base != requestID+"-.md"
}

func exactCurrentBytes(repositoryRoot, path string, expected []byte) string {
	current, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
	if readError != nil {
		return readError.Error()
	}
	if !bytes.Equal(current, expected) {
		return "current bytes do not match the exact preimage"
	}
	return ""
}

func matchingGateRepairs(snapshot *repositorymodel.RepositorySnapshot, sweepKey, fingerprint string) []*repositorymodel.RequestFile {
	matching := []*repositorymodel.RequestFile{}
	for _, requestFile := range snapshot.RequestFiles {
		record := requestFile.TypedRecord
		if record.RepositoryGateRepairValue != "true" || record.FieldEvidenceByName["sweep"].ScalarValue != "true" || record.FieldEvidenceByName["sweep_key"].ScalarValue != sweepKey || (record.RequestStatus != "pending" && record.RequestStatus != "claimed") {
			continue
		}
		if requestFile.ParsedDocument != nil && bodyHasExactLabeledValue(requestFile.ParsedDocument.BodyBytes(), "Diagnostic fingerprint", fingerprint) {
			matching = append(matching, requestFile)
		}
	}
	return matching
}

func requestID(requestFile *repositorymodel.RequestFile) string {
	if requestFile.TypedRecord.RequestID != "" {
		return requestFile.TypedRecord.RequestID
	}
	return requestFile.FilenameID
}

func authoredRepairBytes(gate *DeferGateManifest, userRequestID, baseCommit, mergeCommit string) ([]byte, error) {
	document, parseError := requestmodel.ParseDocument([]byte("---\n---\n"))
	if parseError != nil {
		return nil, parseError
	}
	for _, field := range []struct{ name, value string }{
		{"id", gate.RepairID}, {"title", gate.RepairTitle}, {"status", "pending"}, {"route", "C"},
		{"created_at", gate.RepairCreatedAt}, {"user_request", userRequestID}, {"domain", "backend"},
		{"tdd", "true"}, {"maintenance", "false"}, {"impact", "impact-critical"},
		{"effort_estimate", "effort-substantive"}, {"repository_gate_repair", "true"}, {"sweep", "true"}, {"sweep_key", gate.SweepKey},
	} {
		if setError := document.SetScalar(field.name, field.value); setError != nil {
			return nil, setError
		}
	}
	if listError := document.SetList("depends_on", []string{}); listError != nil {
		return nil, listError
	}
	if listError := document.SetList("related", []string{gate.ParentID}); listError != nil {
		return nil, listError
	}
	body := "\n# " + gate.RepairTitle + "\n\n## What\n\nRepair the repository-gate failure recorded below so dependency-gated requests can resume.\n\n## Instances\n\n" + repairInstanceLine(gate, userRequestID) + "\n" + repairFoldHistory(gate, baseCommit, mergeCommit)
	if replaceError := document.ReplaceBodySpan(0, len(document.BodyBytes()), []byte(body)); replaceError != nil {
		return nil, replaceError
	}
	return document.DocumentBytes(), nil
}

func appendHistory(document *requestmodel.RequestDocument, history string) {
	body := document.BodyBytes()
	prefix := "\n"
	if len(body) == 0 || body[len(body)-1] == '\n' {
		prefix = ""
	}
	_ = document.ReplaceBodySpan(len(body), len(body), []byte(prefix+history))
}

func repairInstanceLine(gate *DeferGateManifest, userRequestID string) string {
	return "- [ ] repository gate: " + gate.DiagnosticFingerprint + " affecting " + gate.ParentID + " (found by " + gate.ParentID + " / " + userRequestID + ")"
}

func appendSweepInstance(document *requestmodel.RequestDocument, instance string) error {
	body := document.BodyBytes()
	lineEnding := []byte("\n")
	if bytes.Contains(body, []byte("\r\n")) {
		lineEnding = []byte("\r\n")
	}
	lines := bytes.Split(body, lineEnding)
	sectionStart := -1
	sectionEnd := len(lines)
	for index, rawLine := range lines {
		line := strings.TrimSpace(string(rawLine))
		if strings.HasPrefix(line, "## ") {
			if sectionStart >= 0 {
				sectionEnd = index
				break
			}
			if line == "## Instances" {
				sectionStart = index
			}
		}
	}
	if sectionStart < 0 {
		return errors.New("repair sweep is missing its canonical ## Instances section")
	}
	for _, rawLine := range lines[sectionStart+1 : sectionEnd] {
		if strings.TrimSpace(string(rawLine)) == instance {
			return nil
		}
	}
	for sectionEnd > sectionStart+1 && strings.TrimSpace(string(lines[sectionEnd-1])) == "" {
		sectionEnd--
	}
	insertAt := 0
	for index := 0; index < sectionEnd; index++ {
		insertAt += len(lines[index])
		if index < len(lines)-1 {
			insertAt += len(lineEnding)
		}
	}
	addition := append([]byte(instance), lineEnding...)
	return document.ReplaceBodySpan(insertAt, insertAt, addition)
}

func bodyHasExactLabeledValue(body []byte, label, expected string) bool {
	prefix := "- **" + label + ":** "
	for _, rawLine := range strings.Split(strings.ReplaceAll(string(body), "\r\n", "\n"), "\n") {
		if strings.HasPrefix(rawLine, prefix) && strings.TrimSpace(strings.TrimPrefix(rawLine, prefix)) == expected {
			return true
		}
	}
	return false
}

func parentDeferralHistory(gate *DeferGateManifest, baseCommit, mergeCommit string) string {
	return "\n## Repository Gate Deferral\n\n" + gateEvidenceLines(gate, baseCommit, mergeCommit)
}

func repairFoldHistory(gate *DeferGateManifest, baseCommit, mergeCommit string) string {
	return "\n## Repository Gate Repair Intake\n\n- **Parent:** " + gate.ParentID + "\n" + gateEvidenceLines(gate, baseCommit, mergeCommit)
}

func gateEvidenceLines(gate *DeferGateManifest, baseCommit, mergeCommit string) string {
	commandBytes, _ := json.Marshal(gate.GateCommand)
	var result strings.Builder
	result.WriteString("- **Gate command (argv JSON):** ")
	result.Write(commandBytes)
	result.WriteString("\n- **Direct exit status:** ")
	result.WriteString(strconv.Itoa(gate.GateExitStatus))
	result.WriteString("\n- **Diagnostic fingerprint:** ")
	result.WriteString(gate.DiagnosticFingerprint)
	result.WriteString("\n- **Repair dependency:** ")
	result.WriteString(gate.RepairID)
	result.WriteString("\n")
	for _, evidence := range gate.DiagnosticEvidence {
		result.WriteString("- **Diagnostic evidence:** ")
		encoded, _ := json.Marshal(evidence)
		result.Write(encoded)
		result.WriteString("\n")
	}
	if baseCommit != "" {
		result.WriteString("- **Implementation base:** ")
		result.WriteString(baseCommit)
		result.WriteString("\n- **Implementation merge:** ")
		result.WriteString(mergeCommit)
		result.WriteString("\n")
	}
	return result.String()
}

func appendUnique(values []string, value string) []string {
	result := append([]string(nil), values...)
	for _, existing := range result {
		if existing == value {
			return result
		}
	}
	return append(result, value)
}

func checkpointWithoutOwnedClaim(contents []byte, requestID, writer string) ([]byte, int) {
	lineEnding := "\n"
	if bytes.Contains(contents, []byte("\r\n")) {
		lineEnding = "\r\n"
	}
	hadTerminal := bytes.HasSuffix(contents, []byte(lineEnding))
	lines := strings.Split(strings.TrimSuffix(string(contents), lineEnding), lineEnding)
	filtered := make([]string, 0, len(lines))
	removed := 0
	for lineIndex := 0; lineIndex < len(lines); lineIndex++ {
		line := lines[lineIndex]
		claimMatch := checkpointClaimPattern.FindStringSubmatch(line)
		if claimMatch != nil && claimMatch[1] == requestID && strings.TrimSpace(claimMatch[2]) == writer {
			removed++
			for lineIndex+1 < len(lines) && (strings.HasPrefix(lines[lineIndex+1], " ") || strings.HasPrefix(lines[lineIndex+1], "\t")) {
				lineIndex++
			}
			continue
		}
		filtered = append(filtered, line)
	}
	result := strings.Join(filtered, lineEnding)
	if hadTerminal {
		result += lineEnding
	}
	return []byte(result), removed
}

func validateDeferredMergeRange(repositoryRoot, base, merge string) (string, string, error) {
	if base == "" {
		return "", "", nil
	}
	canonical := make([]string, 2)
	for index, value := range []string{base, merge} {
		command := exec.Command("git", "-C", repositoryRoot, "rev-parse", "--verify", value+"^{commit}")
		output, commandError := command.Output()
		if commandError != nil {
			return "", "", fmt.Errorf("commit evidence %q does not resolve", value)
		}
		canonical[index] = strings.TrimSpace(string(output))
	}
	ancestor := exec.Command("git", "-C", repositoryRoot, "merge-base", "--is-ancestor", canonical[0], canonical[1])
	if ancestor.Run() != nil || canonical[0] == canonical[1] {
		return "", "", fmt.Errorf("implementation merge must be a non-empty descendant of its base")
	}
	mergedToHead := exec.Command("git", "-C", repositoryRoot, "merge-base", "--is-ancestor", canonical[1], "HEAD")
	if mergedToHead.Run() != nil {
		return "", "", fmt.Errorf("implementation merge must be an ancestor of current HEAD")
	}
	return canonical[0], canonical[1], nil
}

func classifyGatePreimages(repositoryRoot string, paths ...string) ([]string, []string, error) {
	dirty := []string{}
	untracked := []string{}
	for _, path := range paths {
		if !gitPathTracked(repositoryRoot, path) {
			untracked = append(untracked, path)
			continue
		}
		command := exec.Command("git", "-C", repositoryRoot, "diff", "--quiet", "HEAD", "--", path)
		commandError := command.Run()
		if commandError == nil {
			continue
		}
		var exitError *exec.ExitError
		if errors.As(commandError, &exitError) && exitError.ExitCode() == 1 {
			dirty = append(dirty, path)
			continue
		}
		return nil, nil, fmt.Errorf("classify preimage %s: %w", path, commandError)
	}
	sort.Strings(dirty)
	sort.Strings(untracked)
	return dirty, untracked, nil
}

func appendUniqueSorted(existing []string, additions ...string) []string {
	seen := make(map[string]bool, len(existing)+len(additions))
	result := make([]string, 0, len(existing)+len(additions))
	for _, path := range append(append([]string(nil), existing...), additions...) {
		if !seen[path] {
			seen[path] = true
			result = append(result, path)
		}
	}
	sort.Strings(result)
	return result
}

func gitPathTracked(repositoryRoot, path string) bool {
	command := exec.Command("git", "-C", repositoryRoot, "ls-files", "--error-unmatch", "--", path)
	return command.Run() == nil
}

func gitPathCleanAgainstHead(repositoryRoot, path string) bool {
	if !gitPathTracked(repositoryRoot, path) {
		return false
	}
	command := exec.Command("git", "-C", repositoryRoot, "diff", "--quiet", "HEAD", "--", path)
	return command.Run() == nil
}
