package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestAdvanceExecutesEstimateGateAtPublicCLISeam(t *testing.T) {
	repositoryRoot := t.TempDir()
	requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-710", "claimed", "route: C\n", "## Triage\n\nRoute C.\n")
	command := exec.Command(
		advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "advance", "REQ-710", "--",
		"--route", "C", "--write-set", "3", "--subsystems", "2", "--acceptance", "1",
	)
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("advance estimate gate failed: %v\n%s", runError, output)
	}
	var result struct {
		Advance *struct {
			RequestID   string `json:"request_id"`
			RequestPath string `json:"request_path"`
			GateRecords []struct {
				GateID      string   `json:"gate_id"`
				State       string   `json:"state"`
				Provenance  string   `json:"provenance"`
				OutputLines []string `json:"output_lines"`
			} `json:"gate_records"`
		} `json:"advance"`
	}
	if decodeError := json.Unmarshal(output, &result); decodeError != nil {
		t.Fatalf("decode: %v\n%s", decodeError, output)
	}
	if result.Advance == nil || result.Advance.RequestID != "REQ-710" || result.Advance.RequestPath != requestPath {
		t.Fatalf("advance identity was not preserved: %#v", result.Advance)
	}
	if len(result.Advance.GateRecords) != 1 {
		t.Fatalf("gate records=%#v", result.Advance.GateRecords)
	}
	gate := result.Advance.GateRecords[0]
	if gate.GateID != "estimate-p50" || gate.State != "satisfied" || gate.Provenance != "advance_executed" || len(gate.OutputLines) == 0 {
		t.Fatalf("estimate gate=%#v", gate)
	}
}

func TestAdvanceEvidenceGatesReturnTypedMissingInputs(t *testing.T) {
	tests := []struct {
		name        string
		requestID   string
		frontmatter string
		body        string
		wantGate    string
	}{
		{name: "estimate", requestID: "REQ-720", frontmatter: "route: C\n", body: "## Triage\n\nRoute C.\n", wantGate: "estimate-p50"},
		{name: "preflight", requestID: "REQ-721", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n", body: routeCBodyThrough("Scope"), wantGate: "preflight"},
		{name: "qualification", requestID: "REQ-722", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n", body: routeCBodyThrough("Implementation Summary"), wantGate: "qualify"},
		{name: "focused test", requestID: "REQ-723", frontmatter: "route: A\nestimate:\n  p50_active_minutes: 5\n", body: "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n\n## Implementation Summary\n\n- `owned.go` (modified)\n\n## Qualification\n\nPassed.\n", wantGate: "run-blocked-check"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			requestPath := writeAdvanceRequest(t, repositoryRoot, "working", test.requestID, "claimed", test.frontmatter, test.body)
			result, status := runAdvanceGateJSON(t, repositoryRoot, test.requestID, "--request-path", requestPath)
			gate := findAdvanceGate(result, test.wantGate)
			if status == 0 || gate == nil || gate.State != resultmodel.AdvanceGateNeedsInput || !gateHasFinding(*gate, "ADVANCE-GATE-INPUT-REQUIRED") {
				t.Fatalf("status=%d gate=%#v result=%#v", status, gate, result)
			}
		})
	}
}

func TestAdvanceExecutesPreflightAndProjectsGreenEvidence(t *testing.T) {
	repositoryRoot := t.TempDir()
	requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-711", "claimed", "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n", routeCBodyThrough("Scope"))
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Gate Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")

	gateArgv := []string{"sh", "-c", "exit 0"}
	recordCommand := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "record-green-gate", "--gate-exit-status", "0", "--", gateArgv[0], gateArgv[1], gateArgv[2])
	if output, err := recordCommand.CombinedOutput(); err != nil {
		t.Fatalf("record gate: %v\n%s", err, output)
	}

	result, status := runAdvanceGateJSON(t, repositoryRoot,
		"REQ-711", "--request-path", requestPath,
		"--gate-arg", gateArgv[0], "--gate-arg", gateArgv[1], "--gate-arg", gateArgv[2],
		"--", "sh", "-c", "exit 0",
	)
	if status != 0 || result.Advance == nil || len(result.Advance.GateRecords) != 2 {
		t.Fatalf("status=%d result=%#v", status, result)
	}
	preflight, greenGate := result.Advance.GateRecords[0], result.Advance.GateRecords[1]
	if preflight.GateID != "preflight" || preflight.State != resultmodel.AdvanceGateSatisfied || len(preflight.Changes) == 0 {
		t.Fatalf("preflight=%#v", preflight)
	}
	if greenGate.GateID != "green-gate" || greenGate.State != resultmodel.AdvanceGateSatisfied || greenGate.Provenance != resultmodel.AdvanceGateExistingEvidence || greenGate.GreenGate == nil || !greenGate.GreenGate.Matches {
		t.Fatalf("green gate=%#v", greenGate)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", "baseline.json")); err != nil {
		t.Fatal(err)
	}
}

func TestAdvanceQualificationUsesExactRangeAndRunsScopeDrift(t *testing.T) {
	repositoryRoot := t.TempDir()
	requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-712", "claimed", "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n", routeCBodyThrough("Implementation Summary"))
	writeAdvanceFile(t, repositoryRoot, "owned.go", "package owned\n\nvar Value = 1\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Gate Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "base")
	baseRevision := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))
	writeAdvanceFile(t, repositoryRoot, "owned.go", "package owned\n\nvar Value = 2\n")
	runAdvanceGit(t, repositoryRoot, "add", "owned.go")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "implementation")
	targetRevision := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))

	result, status := runAdvanceGateJSON(t, repositoryRoot, "REQ-712", "--request-path", requestPath, "--diff-range", baseRevision+".."+targetRevision)
	if status != 0 || result.Advance == nil || len(result.Advance.GateRecords) != 2 {
		t.Fatalf("status=%d result=%#v", status, result)
	}
	if result.Advance.GateRecords[0].GateID != "qualify" || result.Advance.GateRecords[0].Provenance != resultmodel.AdvanceGateMergedRange || result.Advance.GateRecords[1].GateID != "scope-drift" {
		t.Fatalf("gate records=%#v", result.Advance.GateRecords)
	}

	result, status = runAdvanceGateJSON(t, repositoryRoot, "REQ-712", "--diff-range", "missing..also-missing")
	if status == 0 || result.Advance == nil || result.Advance.GateRecords[0].State != resultmodel.AdvanceGateFindings || !gateHasFinding(result.Advance.GateRecords[0], "QUALIFY-DIFF-RANGE-INVALID") {
		t.Fatalf("invalid range status=%d result=%#v", status, result)
	}
}

func TestAdvanceGreenGateMissRequiresDirectRunThenRecordsIt(t *testing.T) {
	repositoryRoot := t.TempDir()
	requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-716", "claimed", "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n", routeCBodyThrough("Scope"))
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Gate Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
	arguments := []string{"REQ-716", "--request-path", requestPath, "--gate-arg", "sh", "--gate-arg", "-c", "--gate-arg", "exit 0", "--", "sh", "-c", "exit 0"}
	result, status := runAdvanceGateJSON(t, repositoryRoot, arguments...)
	greenGate := findAdvanceGate(result, "green-gate")
	if status == 0 || greenGate == nil || greenGate.State != resultmodel.AdvanceGateNeedsInput || len(greenGate.NextArgv) != 3 || !gateHasFinding(*greenGate, "ADVANCE-GREEN-GATE-DIRECT-RUN-REQUIRED") {
		t.Fatalf("miss status=%d gate=%#v result=%#v", status, greenGate, result)
	}

	recordArguments := append([]string{"REQ-716", "--request-path", requestPath, "--gate-exit-status", "0"}, arguments[3:]...)
	result, status = runAdvanceGateJSON(t, repositoryRoot, recordArguments...)
	greenGate = findAdvanceGate(result, "green-gate")
	if status != 0 || greenGate == nil || greenGate.State != resultmodel.AdvanceGateSatisfied || greenGate.GreenGate == nil || greenGate.GreenGate.State != resultmodel.GateEvidenceRecorded || len(greenGate.Changes) != 1 {
		t.Fatalf("record status=%d gate=%#v result=%#v", status, greenGate, result)
	}
}

func TestAdvanceFocusedTestGateClassifiesBaselineStates(t *testing.T) {
	tests := []struct {
		name             string
		probe            string
		baselineJSON     string
		baselineFailures string
		wantState        resultmodel.FocusedTestBaselineState
		wantGateState    resultmodel.AdvanceGateState
	}{
		{name: "green", probe: "exit 0", baselineJSON: `{"test_command":"exit 0","exit_status":0,"launched":true}`, wantState: resultmodel.FocusedBaselineGreen, wantGateState: resultmodel.AdvanceGateSatisfied},
		{name: "matching red", probe: "echo same; exit 17", baselineJSON: `{"test_command":"echo same; exit 17","exit_status":17,"launched":true}`, baselineFailures: "same\n", wantState: resultmodel.FocusedBaselineMatchingRed, wantGateState: resultmodel.AdvanceGateSatisfied},
		{name: "different red status", probe: "echo same; exit 18", baselineJSON: `{"test_command":"echo same; exit 18","exit_status":17,"launched":true}`, baselineFailures: "same\n", wantState: resultmodel.FocusedBaselineNewRed, wantGateState: resultmodel.AdvanceGateFindings},
		{name: "new red", probe: "echo new; exit 18", baselineJSON: `{"test_command":"echo old; exit 17","exit_status":17,"launched":true}`, baselineFailures: "old\n", wantState: resultmodel.FocusedBaselineNewRed, wantGateState: resultmodel.AdvanceGateFindings},
		{name: "unusable baseline", probe: "exit 19", baselineJSON: `{"test_command":"exit 19","exit_status":127,"launched":false}`, wantState: resultmodel.FocusedBaselineUnusable, wantGateState: resultmodel.AdvanceGateFindings},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			body := "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n\n## Implementation Summary\n\n- `owned.go` (modified)\n\n## Qualification\n\nPassed.\n"
			writeAdvanceRequest(t, repositoryRoot, "working", "REQ-713", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", body)
			writeAdvanceFile(t, repositoryRoot, "focused.sh", test.probe)
			writeAdvanceFile(t, repositoryRoot, "do-work/working/baseline.json", test.baselineJSON)
			if test.baselineFailures != "" {
				writeAdvanceFile(t, repositoryRoot, "do-work/working/baseline-failures.txt", test.baselineFailures)
			}
			result, _ := runAdvanceGateJSON(t, repositoryRoot, "REQ-713", "--", "--probe-file", "focused.sh", "--timeout-seconds", "2")
			focusedGate := findAdvanceGate(result, "run-blocked-check")
			if focusedGate == nil || focusedGate.FocusedTest == nil || focusedGate.FocusedTest.BaselineState != test.wantState || focusedGate.State != test.wantGateState {
				t.Fatalf("focused gate=%#v result=%#v", focusedGate, result)
			}
		})
	}
}

func TestAdvanceGateInputsFailClosedAndNeverInterpolateHostileTokens(t *testing.T) {
	repositoryRoot := t.TempDir()
	requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-714", "claimed", "route: C\n", "## Triage\n\nRoute C.\n")
	markerPath := filepath.Join(repositoryRoot, "interpolated")
	result, status := runAdvanceGateJSON(t, repositoryRoot, "REQ-714", "--request-path", "do-work/working/REQ-999-wrong.md", "--", "--route", "C")
	if status == 0 || !hasAdvanceResultFinding(result, "ADVANCE-EVIDENCE-MISMATCH") {
		t.Fatalf("mismatched identity status=%d result=%#v", status, result)
	}
	hostile := "$(touch " + markerPath + ")"
	result, status = runAdvanceGateJSON(t, repositoryRoot, "REQ-714", "--request-path", requestPath, "--", "--route", "C", hostile)
	if status == 0 || result.Advance == nil || result.Advance.GateRecords[0].State != resultmodel.AdvanceGateFailed {
		t.Fatalf("hostile token status=%d result=%#v", status, result)
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("hostile token executed: %v", err)
	}
}

func TestAdvanceFocusedTestGateDistinguishesTimeoutAndLaunchFailure(t *testing.T) {
	for _, test := range []struct {
		name       string
		probe      string
		timeout    string
		emptyPath  bool
		wantStatus int
		wantState  resultmodel.AdvanceGateState
	}{
		{name: "timeout", probe: "sleep 5", timeout: "1", wantStatus: 124, wantState: resultmodel.AdvanceGateFindings},
		{name: "launch failure", probe: "exit 0", timeout: "2", emptyPath: true, wantStatus: 125, wantState: resultmodel.AdvanceGateFailed},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			body := "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n\n## Implementation Summary\n\n- `owned.go` (modified)\n\n## Qualification\n\nPassed.\n"
			writeAdvanceRequest(t, repositoryRoot, "working", "REQ-715", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", body)
			writeAdvanceFile(t, repositoryRoot, "focused.sh", test.probe)
			commandArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "advance", "REQ-715", "--", "--probe-file", "focused.sh", "--timeout-seconds", test.timeout}
			command := exec.Command(advanceCLIBinary(t), commandArguments...)
			if test.emptyPath {
				command.Env = append(os.Environ(), "PATH=")
			}
			output, _ := command.CombinedOutput()
			var result resultmodel.CommandResult
			if err := json.Unmarshal(output, &result); err != nil {
				t.Fatalf("decode: %v\n%s", err, output)
			}
			focusedGate := findAdvanceGate(result, "run-blocked-check")
			if focusedGate == nil || focusedGate.FocusedTest == nil || focusedGate.FocusedTest.ExitStatus != test.wantStatus || focusedGate.State != test.wantState {
				t.Fatalf("focused gate=%#v result=%#v", focusedGate, result)
			}
		})
	}
}

func runAdvanceGateJSON(t *testing.T, repositoryRoot string, arguments ...string) (resultmodel.CommandResult, int) {
	t.Helper()
	commandArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "advance"}
	commandArguments = append(commandArguments, arguments...)
	command := exec.Command(advanceCLIBinary(t), commandArguments...)
	output, runError := command.CombinedOutput()
	status := 0
	if runError != nil {
		if exitError, ok := runError.(*exec.ExitError); ok {
			status = exitError.ExitCode()
		} else {
			t.Fatalf("advance launch: %v", runError)
		}
	}
	var result resultmodel.CommandResult
	if decodeError := json.Unmarshal(output, &result); decodeError != nil {
		t.Fatalf("decode: %v\n%s", decodeError, output)
	}
	return result, status
}

func findAdvanceGate(result resultmodel.CommandResult, gateID string) *resultmodel.AdvanceGateRecord {
	if result.Advance == nil {
		return nil
	}
	for index := range result.Advance.GateRecords {
		if result.Advance.GateRecords[index].GateID == gateID {
			return &result.Advance.GateRecords[index]
		}
	}
	return nil
}

func gateHasFinding(gate resultmodel.AdvanceGateRecord, code string) bool {
	for _, finding := range gate.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

func hasAdvanceResultFinding(result resultmodel.CommandResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
