package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
		// A child may choose the runner's own reserved values for itself; a launched
		// completion at 124 or 125 is ordinary red, not a timeout or a launch failure.
		{name: "ordinary reserved timeout value", probe: "echo same; exit 124", baselineJSON: `{"test_command":"echo same; exit 124","exit_status":124,"launched":true}`, baselineFailures: "same\n", wantState: resultmodel.FocusedBaselineMatchingRed, wantGateState: resultmodel.AdvanceGateSatisfied},
		{name: "ordinary reserved launch value", probe: "echo same; exit 125", baselineJSON: `{"test_command":"echo same; exit 125","exit_status":125,"launched":true}`, baselineFailures: "same\n", wantState: resultmodel.FocusedBaselineMatchingRed, wantGateState: resultmodel.AdvanceGateSatisfied},
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
			// Every row is an ordinary launched completion, whatever integer it exits with.
			if !focusedGate.FocusedTest.Launched || focusedGate.FocusedTest.TimedOut {
				t.Fatalf("execution facts=%#v", focusedGate.FocusedTest)
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

// advanceGateFinding returns the one finding a gate record carries under the
// given code, so a test can assert on its remedy rather than only its presence.
func advanceGateFinding(t *testing.T, gate resultmodel.AdvanceGateRecord, code string) resultmodel.CommandFinding {
	t.Helper()
	for _, finding := range gate.Findings {
		if finding.Code == code {
			return finding
		}
	}
	t.Fatalf("gate %q carries no %s finding: %#v", gate.GateID, code, gate.Findings)
	return resultmodel.CommandFinding{}
}

func hasAdvanceResultFinding(result resultmodel.CommandResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

// focusedGateRouteABody is the shortest claimed Route A body that reaches the
// test-gate phase, matching the literal the existing focused-gate tables use.
const focusedGateRouteABody = "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n\n## Implementation Summary\n\n- `owned.go` (modified)\n\n## Qualification\n\nPassed.\n"

// canonicalGateFixtureBinary stands in for the repository gate: advance never
// executes it, so a single absolute path keeps the recorded argv resolvable
// even when the probe's own PATH is restricted.
const canonicalGateFixtureBinary = "/usr/bin/true"

// TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline pins the
// REQ-506 false success. A saved baseline that happens to carry the same reserved
// exit value must not exclude a current execution that never launched or that the
// timer killed, and advance must not exit successfully for either.
func TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline(t *testing.T) {
	t.Run("current launch failure", func(t *testing.T) {
		repositoryRoot := t.TempDir()
		writeAdvanceRequest(t, repositoryRoot, "working", "REQ-717", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", focusedGateRouteABody)
		writeAdvanceFile(t, repositoryRoot, "focused.sh", "exit 125")
		writeAdvanceFile(t, repositoryRoot, "do-work/working/baseline.json", `{"test_command":"exit 125","exit_status":125,"launched":true}`)
		writeAdvanceFile(t, repositoryRoot, "do-work/working/baseline-failures.txt", "")
		initAdvanceGitFixture(t, repositoryRoot)
		recordAdvanceGreenGate(t, repositoryRoot, canonicalGateFixtureBinary)

		result, status := runAdvanceGateJSONWithPath(t, repositoryRoot, gitOnlyPathDirectory(t),
			"REQ-717", "--gate-arg", canonicalGateFixtureBinary, "--", "--probe-file", "focused.sh", "--timeout-seconds", "2")
		greenGate := findAdvanceGate(result, "green-gate")
		if greenGate == nil || greenGate.State != resultmodel.AdvanceGateSatisfied {
			t.Fatalf("canonical green evidence did not stand: gate=%#v result=%#v", greenGate, result)
		}
		focusedGate := findAdvanceGate(result, "run-blocked-check")
		if focusedGate == nil || focusedGate.FocusedTest == nil || focusedGate.FocusedTest.Launched {
			t.Fatalf("focused gate=%#v result=%#v", focusedGate, result)
		}
		if focusedGate.FocusedTest.BaselineState == resultmodel.FocusedBaselineMatchingRed || focusedGate.State != resultmodel.AdvanceGateFailed {
			t.Fatalf("failed launch cleared the focused boundary: gate=%#v", focusedGate)
		}
		if !gateHasFinding(*focusedGate, "BLOCKED-PROBE-LAUNCH-FAILED") || status == 0 {
			t.Fatalf("status=%d gate=%#v", status, focusedGate)
		}
	})

	t.Run("current timeout", func(t *testing.T) {
		repositoryRoot := t.TempDir()
		writeAdvanceRequest(t, repositoryRoot, "working", "REQ-718", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", focusedGateRouteABody)
		writeAdvanceFile(t, repositoryRoot, "focused.sh", "/bin/sleep 2; exit 124")
		initAdvanceGitFixture(t, repositoryRoot)
		recordAdvanceGreenGate(t, repositoryRoot, canonicalGateFixtureBinary)

		// The saved baseline is the identical command run to completion: it exits 124
		// under its own power, which is the same integer the timeout path reports.
		baseline, _ := runAdvanceGateJSON(t, repositoryRoot, "REQ-718",
			"--gate-arg", canonicalGateFixtureBinary, "--", "--probe-file", "focused.sh", "--timeout-seconds", "10")
		baselineGate := findAdvanceGate(baseline, "run-blocked-check")
		if baselineGate == nil || baselineGate.FocusedTest == nil || baselineGate.FocusedTest.ExitStatus != 124 {
			t.Fatalf("direct baseline run=%#v result=%#v", baselineGate, baseline)
		}
		writeAdvanceFile(t, repositoryRoot, "do-work/working/baseline.json", `{"test_command":"/bin/sleep 2; exit 124","exit_status":124,"launched":true}`)
		writeAdvanceFile(t, repositoryRoot, "do-work/working/baseline-failures.txt", "")

		result, status := runAdvanceGateJSON(t, repositoryRoot, "REQ-718",
			"--gate-arg", canonicalGateFixtureBinary, "--", "--probe-file", "focused.sh", "--timeout-seconds", "1")
		focusedGate := findAdvanceGate(result, "run-blocked-check")
		if focusedGate == nil || focusedGate.FocusedTest == nil || !focusedGate.FocusedTest.TimedOut {
			t.Fatalf("focused gate=%#v result=%#v", focusedGate, result)
		}
		if focusedGate.FocusedTest.BaselineState == resultmodel.FocusedBaselineMatchingRed || focusedGate.State == resultmodel.AdvanceGateSatisfied || status == 0 {
			t.Fatalf("timeout cleared the focused boundary: status=%d gate=%#v", status, focusedGate)
		}
	})
}

// TestAdvanceMissingInputContinuationsPreserveArgumentChannels pins REQ-506's
// second defect: a missing-input continuation must place each judgment-owned
// input in the channel that parses it, so substituting a real value into the
// emitted argv is enough to make progress.
func TestAdvanceMissingInputContinuationsPreserveArgumentChannels(t *testing.T) {
	t.Run("qualification diff range", func(t *testing.T) {
		repositoryRoot := t.TempDir()
		requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-724", "claimed",
			"route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			routeCBodyThrough("Implementation Summary"))
		writeAdvanceFile(t, repositoryRoot, "owned.go", "package owned\n\nvar Value = 1\n")
		initAdvanceGitFixture(t, repositoryRoot)
		baseRevision := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))
		writeAdvanceFile(t, repositoryRoot, "owned.go", "package owned\n\nvar Value = 2\n")
		runAdvanceGit(t, repositoryRoot, "add", "owned.go")
		runAdvanceGit(t, repositoryRoot, "commit", "-qm", "implementation")
		mergeRevision := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))

		result, status := runAdvanceGateJSON(t, repositoryRoot, "REQ-724", "--request-path", requestPath)
		qualify := findAdvanceGate(result, "qualify")
		if status == 0 || qualify == nil || qualify.State != resultmodel.AdvanceGateNeedsInput {
			t.Fatalf("status=%d gate=%#v result=%#v", status, qualify, result)
		}
		if advanceArgvIndex(qualify.NextArgv, "--diff-range") < 0 || advanceArgvIndex(qualify.NextArgv, "--") >= 0 {
			t.Fatalf("qualification continuation is in the wrong channel: %#v", qualify.NextArgv)
		}
		continuation := substituteAdvancePlaceholder(t, qualify.NextArgv, baseRevision+".."+mergeRevision)
		followed, followedStatus := runAdvanceContinuation(t, repositoryRoot, continuation)
		followedQualify := findAdvanceGate(followed, "qualify")
		if followedStatus != 0 || followedQualify == nil || followedQualify.State != resultmodel.AdvanceGateSatisfied {
			t.Fatalf("continuation %v status=%d gate=%#v result=%#v", continuation, followedStatus, followedQualify, followed)
		}
	})

	t.Run("canonical gate tokens", func(t *testing.T) {
		repositoryRoot := t.TempDir()
		requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-725", "claimed",
			"route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			routeCBodyThrough("Scope"))
		initAdvanceGitFixture(t, repositoryRoot)
		recordAdvanceGreenGate(t, repositoryRoot, canonicalGateFixtureBinary)

		result, status := runAdvanceGateJSON(t, repositoryRoot, "REQ-725", "--request-path", requestPath, "--", "sh", "-c", "exit 0")
		greenGate := findAdvanceGate(result, "check-green-gate")
		if status == 0 || greenGate == nil || greenGate.State != resultmodel.AdvanceGateNeedsInput {
			t.Fatalf("status=%d gate=%#v result=%#v", status, greenGate, result)
		}
		separatorIndex, gateArgIndex := advanceArgvIndex(greenGate.NextArgv, "--"), advanceArgvIndex(greenGate.NextArgv, "--gate-arg")
		if gateArgIndex < 0 || separatorIndex < 0 || gateArgIndex > separatorIndex {
			t.Fatalf("canonical gate continuation is in the wrong channel: %#v", greenGate.NextArgv)
		}
		continuation := substituteAdvancePlaceholder(t, greenGate.NextArgv, canonicalGateFixtureBinary)
		followed, followedStatus := runAdvanceContinuation(t, repositoryRoot, continuation)
		followedGreenGate := findAdvanceGate(followed, "green-gate")
		if followedStatus != 0 || followedGreenGate == nil || followedGreenGate.State != resultmodel.AdvanceGateSatisfied {
			t.Fatalf("continuation %v status=%d gate=%#v result=%#v", continuation, followedStatus, followedGreenGate, followed)
		}
	})
}

// TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation pins the remedy
// redirection, which nothing else in this package reads. A subordinate finding
// whose own remedy would send the action back into the evidence helper this
// gate has already run is rewritten to point at the same request-bound advance
// invocation instead. Deleting either `redirectHelperRemedies` call site in
// evidence_gates.go — the one in composeCoreGate or the one in composeGreenGate
// — makes this test fail. It is the only assertion in this file on a *finding's*
// remedy; every other one reads the record-level NextArgv, which is written
// elsewhere and survives both deletions.
//
// The two negative controls prove the rewrite is selective and not blanket: the
// sibling FOCUSED-BASELINE-MISSING finding carries no remedy and must still
// carry none, and the green gate's git remedy diagnoses the failure instead of
// re-entering the helper, so its owner's argv must survive untouched.
func TestAdvanceRedirectsSubordinateRemediesToItsOwnContinuation(t *testing.T) {
	t.Run("core gate call site", func(t *testing.T) {
		repositoryRoot := t.TempDir()
		requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-726", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", focusedGateRouteABody)
		writeAdvanceFile(t, repositoryRoot, "focused.sh", "exit 0")

		// The fixture is deliberately not a Git repository and carries no saved
		// baseline, so one invocation produces all three findings: a redirected probe
		// remedy, a remedy-less baseline finding, and the green gate's git remedy.
		result, _ := runAdvanceGateJSON(t, repositoryRoot, "REQ-726", "--request-path", requestPath,
			"--gate-arg", canonicalGateFixtureBinary,
			"--", "--probe-file", "focused.sh", "--timeout-seconds", "2")

		focusedGate, greenGate := findAdvanceGate(result, "run-blocked-check"), findAdvanceGate(result, "green-gate")
		if focusedGate == nil || greenGate == nil {
			t.Fatalf("focused gate=%#v green gate=%#v result=%#v", focusedGate, greenGate, result)
		}
		wantContinuation := []string{
			"do-work-cli", "--format", "json", "advance", "REQ-726", "--request-path", requestPath,
			"--gate-arg", canonicalGateFixtureBinary, "--", "--probe-file", "focused.sh", "--timeout-seconds", "2",
		}

		probeFinding := advanceGateFinding(t, *focusedGate, "BLOCKED-PROBE-SUCCEEDED")
		if !slices.Equal(probeFinding.NextArgv, wantContinuation) {
			t.Fatalf("the probe finding's remedy was not redirected to this advance invocation:\n got %#v\nwant %#v", probeFinding.NextArgv, wantContinuation)
		}
		if !slices.Equal(probeFinding.VerificationArgv, wantContinuation) {
			t.Fatalf("the probe finding's verification was not redirected to this advance invocation:\n got %#v\nwant %#v", probeFinding.VerificationArgv, wantContinuation)
		}

		baselineFinding := advanceGateFinding(t, *focusedGate, "FOCUSED-BASELINE-MISSING")
		if len(baselineFinding.NextArgv) != 0 || len(baselineFinding.VerificationArgv) != 0 {
			t.Fatalf("the redirection invented a remedy for a finding that has none: next=%#v verification=%#v", baselineFinding.NextArgv, baselineFinding.VerificationArgv)
		}

		gateEvidenceFinding := advanceGateFinding(t, *greenGate, "GATE-EVIDENCE-CHECK-FAILED")
		if len(gateEvidenceFinding.NextArgv) == 0 || gateEvidenceFinding.NextArgv[0] != "git" || len(gateEvidenceFinding.VerificationArgv) == 0 || gateEvidenceFinding.VerificationArgv[0] != "git" {
			t.Fatalf("a diagnostic remedy that does not run do-work-cli was rewritten: next=%#v verification=%#v", gateEvidenceFinding.NextArgv, gateEvidenceFinding.VerificationArgv)
		}
	})

	// The green-gate call site has no natural producer: of the findings the two
	// green-gate helpers can return, the git remedies name another tool and the
	// not-green remedy is the caller's own gate argv, so none of them re-enters
	// the helper that just ran. The condition the redirection keys on is
	// "argv[0] is do-work-cli and its verb is the subordinate command", so this
	// case supplies that condition through the one channel that can carry it —
	// the gate argv, which advance copies into the not-green remedy verbatim and
	// never executes.
	t.Run("green gate call site", func(t *testing.T) {
		repositoryRoot := t.TempDir()
		requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-727", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", focusedGateRouteABody)
		writeAdvanceFile(t, repositoryRoot, "focused.sh", "exit 0")

		result, _ := runAdvanceGateJSON(t, repositoryRoot, "REQ-727", "--request-path", requestPath,
			"--gate-exit-status", "3", "--gate-arg", "do-work-cli", "--gate-arg", "record-green-gate",
			"--", "--probe-file", "focused.sh", "--timeout-seconds", "2")

		greenGate := findAdvanceGate(result, "green-gate")
		if greenGate == nil {
			t.Fatalf("result=%#v", result)
		}
		wantContinuation := []string{
			"do-work-cli", "--format", "json", "advance", "REQ-727", "--request-path", requestPath,
			"--gate-exit-status", "3", "--gate-arg", "do-work-cli", "--gate-arg", "record-green-gate",
			"--", "--probe-file", "focused.sh", "--timeout-seconds", "2",
		}
		notGreenFinding := advanceGateFinding(t, *greenGate, "GATE-EVIDENCE-NOT-GREEN")
		if !slices.Equal(notGreenFinding.NextArgv, wantContinuation) {
			t.Fatalf("the green gate's remedy was not redirected to this advance invocation:\n got %#v\nwant %#v", notGreenFinding.NextArgv, wantContinuation)
		}
		// Same finding, other field: its verification names check-green-gate, which
		// is not the helper this gate ran, so the redirection must leave it alone.
		if slices.Equal(notGreenFinding.VerificationArgv, wantContinuation) {
			t.Fatalf("a remedy naming a different helper was rewritten: %#v", notGreenFinding.VerificationArgv)
		}
	})
}

// TestFocusedGateStateKeepsSubordinateAuthority pins the three-part guard that
// opens focusedGateState, one part per row. Deleting
// `subordinateState == resultmodel.AdvanceGateFailed ||` reddens the two failed
// rows; deleting `|| focusedTest.TimedOut` reddens the two timed-out rows;
// deleting `!focusedTest.Launched ||` reddens the never-launched row. Without
// the guard each of those executions is cleared by a saved baseline it was
// never eligible to be compared against.
//
// This is the only in-process call to product code in this package's tests, and
// it is deliberate. The guard is defence in depth: its one current caller,
// composeCoreGate, receives its FocusedTestResult from run-blocked-check, whose
// own eligibility check (`finishedOnItsOwn` in internal/corehelpers) already
// leaves every failed, unlaunched or timed-out execution at
// FocusedBaselineNotCompared, where the switch below the guard has no case and
// falls through to the same answer. So no argv reaches these rows through the
// CLI, and a public test could only reach them once a second producer of
// FocusedTestResult exists. The rows are therefore keyed on the execution facts
// the guard reads — Launched, TimedOut and BaselineState — never on the exit
// statuses that happen to produce them today.
func TestFocusedGateStateKeepsSubordinateAuthority(t *testing.T) {
	for _, test := range []struct {
		name             string
		subordinateState resultmodel.AdvanceGateState
		launched         bool
		timedOut         bool
		baselineState    resultmodel.FocusedTestBaselineState
		wantState        resultmodel.AdvanceGateState
	}{
		{name: "failed subordinate against a green baseline", subordinateState: resultmodel.AdvanceGateFailed, launched: true, baselineState: resultmodel.FocusedBaselineGreen, wantState: resultmodel.AdvanceGateFailed},
		{name: "failed subordinate against a matching red baseline", subordinateState: resultmodel.AdvanceGateFailed, launched: true, baselineState: resultmodel.FocusedBaselineMatchingRed, wantState: resultmodel.AdvanceGateFailed},
		{name: "timed out against a green baseline", subordinateState: resultmodel.AdvanceGateFindings, launched: true, timedOut: true, baselineState: resultmodel.FocusedBaselineGreen, wantState: resultmodel.AdvanceGateFindings},
		{name: "timed out against a matching red baseline", subordinateState: resultmodel.AdvanceGateFindings, launched: true, timedOut: true, baselineState: resultmodel.FocusedBaselineMatchingRed, wantState: resultmodel.AdvanceGateFindings},
		{name: "never launched against a green baseline", subordinateState: resultmodel.AdvanceGateFindings, baselineState: resultmodel.FocusedBaselineGreen, wantState: resultmodel.AdvanceGateFindings},
		// Negative controls: an execution the guard admits is still classified by
		// its baseline, so the guard is a boundary and not a blanket refusal.
		{name: "eligible execution cleared by a green baseline", subordinateState: resultmodel.AdvanceGateFindings, launched: true, baselineState: resultmodel.FocusedBaselineGreen, wantState: resultmodel.AdvanceGateSatisfied},
		{name: "eligible execution cleared by a matching red baseline", subordinateState: resultmodel.AdvanceGateFindings, launched: true, baselineState: resultmodel.FocusedBaselineMatchingRed, wantState: resultmodel.AdvanceGateSatisfied},
		{name: "eligible execution reddened by a new red baseline", subordinateState: resultmodel.AdvanceGateSatisfied, launched: true, baselineState: resultmodel.FocusedBaselineNewRed, wantState: resultmodel.AdvanceGateFindings},
		{name: "eligible execution that was never compared", subordinateState: resultmodel.AdvanceGateFindings, launched: true, baselineState: resultmodel.FocusedBaselineNotCompared, wantState: resultmodel.AdvanceGateFindings},
	} {
		t.Run(test.name, func(t *testing.T) {
			focusedTest := &resultmodel.FocusedTestResult{Launched: test.launched, TimedOut: test.timedOut, BaselineState: test.baselineState}
			gotState := focusedGateState(test.subordinateState, focusedTest)
			if gotState != test.wantState {
				t.Fatalf("a saved baseline decided the gate for subordinate=%s launched=%t timed_out=%t baseline=%s: got %s, want %s",
					test.subordinateState, test.launched, test.timedOut, test.baselineState, gotState, test.wantState)
			}
		})
	}
}

// TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure pins the finding
// code for the one focused-test outcome nothing in this package asserted: a
// probe that launched, that the runner's timer did not stop, and that the
// runner reported an error for because the run was interrupted. It reports
// BLOCKED-PROBE-FAILED at error severity and fails the gate. That code used to
// be BLOCKED-PROBE-LAUNCH-FAILED, and because both are error severity nothing
// moved when it changed; restoring the old code reddens the interrupted row.
//
// The ordinary red row is what keeps the two conditions apart. Both probes exit
// 143, so the only difference between the rows is whether the runner observed an
// interruption — never the integer the child chose. Collapsing the interrupted
// case into the ordinary non-zero case therefore reddens the interrupted row on
// severity and gate state as well as on the code.
func TestAdvanceFocusedGateReportsAnInterruptedProbeAsAFailure(t *testing.T) {
	for _, test := range []struct {
		name          string
		probe         string
		wantSeverity  resultmodel.FindingSeverity
		wantGateState resultmodel.AdvanceGateState
	}{
		// The probe's parent is the do-work-cli process running this gate, and its
		// signal handling is already armed, so the probe interrupts the run through
		// its own parent. The test process never handles a signal itself.
		{name: "interrupted run", probe: "kill -TERM $PPID; sleep 5", wantSeverity: resultmodel.SeverityError, wantGateState: resultmodel.AdvanceGateFailed},
		{name: "ordinary red run", probe: "echo failing; exit 143", wantSeverity: resultmodel.SeverityWarning, wantGateState: resultmodel.AdvanceGateFindings},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-728", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", focusedGateRouteABody)
			writeAdvanceFile(t, repositoryRoot, "focused.sh", test.probe)

			result, status := runAdvanceGateJSON(t, repositoryRoot, "REQ-728", "--request-path", requestPath,
				"--", "--probe-file", "focused.sh", "--timeout-seconds", "10")
			focusedGate := findAdvanceGate(result, "run-blocked-check")
			if focusedGate == nil || focusedGate.FocusedTest == nil || status == 0 {
				t.Fatalf("status=%d gate=%#v result=%#v", status, focusedGate, result)
			}
			if focusedGate.FocusedTest.ExitStatus != 143 || !focusedGate.FocusedTest.Launched || focusedGate.FocusedTest.TimedOut {
				t.Fatalf("both rows must reach the gate with the same execution facts: %#v", focusedGate.FocusedTest)
			}
			probeFinding := advanceGateFinding(t, *focusedGate, "BLOCKED-PROBE-FAILED")
			if probeFinding.Severity != test.wantSeverity || focusedGate.State != test.wantGateState {
				t.Fatalf("an interrupted probe and a probe that chose 143 for itself were classified alike: severity=%s want %s, gate state=%s want %s",
					probeFinding.Severity, test.wantSeverity, focusedGate.State, test.wantGateState)
			}
		})
	}
}

func initAdvanceGitFixture(t *testing.T, repositoryRoot string) {
	t.Helper()
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Gate Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
}

func recordAdvanceGreenGate(t *testing.T, repositoryRoot string, gateArgv ...string) {
	t.Helper()
	arguments := append([]string{"--repo-root", repositoryRoot, "--format", "json", "record-green-gate", "--gate-exit-status", "0", "--"}, gateArgv...)
	if output, err := exec.Command(advanceCLIBinary(t), arguments...).CombinedOutput(); err != nil {
		t.Fatalf("record green gate: %v\n%s", err, output)
	}
}

// gitOnlyPathDirectory returns a PATH value holding nothing but git, so the real
// focused probe cannot launch its shell while canonical Git evidence still resolves.
func gitOnlyPathDirectory(t *testing.T) string {
	t.Helper()
	gitBinary, err := exec.LookPath("git")
	if err != nil {
		t.Skipf("git is required for canonical gate evidence: %v", err)
	}
	directory := t.TempDir()
	if err := os.Symlink(gitBinary, filepath.Join(directory, "git")); err != nil {
		t.Fatal(err)
	}
	return directory
}

func runAdvanceGateJSONWithPath(t *testing.T, repositoryRoot, pathValue string, arguments ...string) (resultmodel.CommandResult, int) {
	t.Helper()
	commandArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "advance"}
	commandArguments = append(commandArguments, arguments...)
	command := exec.Command(advanceCLIBinary(t), commandArguments...)
	command.Env = append(os.Environ(), "PATH="+pathValue)
	return decodeAdvanceGateRun(t, command)
}

// runAdvanceContinuation executes one emitted advance continuation verbatim,
// adding only the fixture repository root the harness must supply.
func runAdvanceContinuation(t *testing.T, repositoryRoot string, continuation []string) (resultmodel.CommandResult, int) {
	t.Helper()
	if len(continuation) == 0 || continuation[0] != "do-work-cli" {
		t.Fatalf("continuation is not a do-work-cli invocation: %#v", continuation)
	}
	commandArguments := append([]string{"--repo-root", repositoryRoot}, continuation[1:]...)
	return decodeAdvanceGateRun(t, exec.Command(advanceCLIBinary(t), commandArguments...))
}

func decodeAdvanceGateRun(t *testing.T, command *exec.Cmd) (resultmodel.CommandResult, int) {
	t.Helper()
	output, runError := command.CombinedOutput()
	status := 0
	if runError != nil {
		exitError, ok := runError.(*exec.ExitError)
		if !ok {
			t.Fatalf("advance launch: %v", runError)
		}
		status = exitError.ExitCode()
	}
	var result resultmodel.CommandResult
	if decodeError := json.Unmarshal(output, &result); decodeError != nil {
		t.Fatalf("decode: %v\n%s", decodeError, output)
	}
	return result, status
}

// substituteAdvancePlaceholder replaces the single angle-bracket placeholder in an
// emitted continuation with a real value, moving no flag and no separator.
func substituteAdvancePlaceholder(t *testing.T, continuation []string, value string) []string {
	t.Helper()
	substituted := append([]string(nil), continuation...)
	placeholders := 0
	for index, token := range substituted {
		if strings.Contains(token, "<") {
			substituted[index] = value
			placeholders++
		}
	}
	if placeholders != 1 {
		t.Fatalf("expected exactly one placeholder token in %#v", continuation)
	}
	return substituted
}

func advanceArgvIndex(argv []string, wanted string) int {
	for index, token := range argv {
		if token == wanted {
			return index
		}
	}
	return -1
}
