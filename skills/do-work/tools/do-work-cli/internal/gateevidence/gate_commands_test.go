package gateevidence

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestGateEvidenceHandlersAreComplete(t *testing.T) {
	handlers := Handlers()
	if len(handlers) != 2 || handlers[CommandCheckGreenGate] == nil || handlers[CommandRecordGreenGate] == nil {
		t.Fatalf("handlers = %#v", handlers)
	}
}

func TestRecordGreenGateRefusesNonzeroWithoutSelfReferentialRecovery(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	gateCommand := []string{"bash", "verify.sh"}
	result := handleRecordGreenGate(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--gate-exit-status", "17", "--", "bash", "verify.sh"})
	if result.Outcome != resultmodel.OutcomeRefused || result.GateEvidence == nil || result.GateEvidence.State != resultmodel.GateEvidenceNotGreen || len(result.Changes) != 0 || len(result.Findings) != 1 {
		t.Fatalf("result = %#v", result)
	}
	finding := result.Findings[0]
	if !equalArgv(finding.NextArgv, gateCommand) || len(finding.NextArgv) == 0 || finding.NextArgv[0] == CommandRecordGreenGate {
		t.Fatalf("self-referential or inexact next argv: %#v", finding)
	}
	if len(finding.VerificationArgv) < 4 || finding.VerificationArgv[3] != CommandCheckGreenGate {
		t.Fatalf("verification argv = %#v", finding.VerificationArgv)
	}
	commonDirectory := strings.TrimSpace(runGateEvidenceTestGit(t, repositoryRoot, "rev-parse", "--git-common-dir"))
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(repositoryRoot, commonDirectory)
	}
	if _, err := os.Stat(filepath.Join(commonDirectory, gateEvidenceDirectoryName)); !os.IsNotExist(err) {
		t.Fatalf("non-green refusal created evidence directory: %v", err)
	}
}

func TestCheckGreenGateHandlerHonorsExplicitTargetRevision(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	gateCommand := []string{"bash", "verify.sh"}
	recorded, err := RecordGreenGate(repositoryRoot, gateCommand)
	if err != nil {
		t.Fatal(err)
	}
	writeGateEvidenceTestFile(t, repositoryRoot, "peer.txt", "unrelated\n")
	runGateEvidenceTestGit(t, repositoryRoot, "add", "peer.txt")
	runGateEvidenceTestGit(t, repositoryRoot, "commit", "-qm", "unrelated peer commit")

	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}
	headBound := handleCheckGreenGate(executionContext, []string{"--", "bash", "verify.sh"})
	if headBound.Outcome != resultmodel.OutcomeSuccess || headBound.GateEvidence.Matches {
		t.Fatalf("HEAD-bound handler result = %#v", headBound.GateEvidence)
	}
	targeted := handleCheckGreenGate(executionContext, []string{"--at-revision", recorded.RecordedRevision, "--", "bash", "verify.sh"})
	if targeted.Outcome != resultmodel.OutcomeSuccess || !targeted.GateEvidence.Matches || targeted.GateEvidence.TargetRevision != recorded.RecordedRevision {
		t.Fatalf("targeted handler result = %#v", targeted.GateEvidence)
	}
	usage := handleCheckGreenGate(executionContext, []string{"--at-revision", "--", "bash"})
	if usage.Outcome != resultmodel.OutcomeFailure || len(usage.Findings) != 1 || usage.Findings[0].Code != "GATE-EVIDENCE-USAGE" {
		t.Fatalf("usage result = %#v", usage)
	}
}

func TestGateEvidenceCommandsKeepTextAndJSONInParity(t *testing.T) {
	repositoryRoot := newGateEvidenceRepository(t)
	handlers := Handlers()
	gateArguments := []string{"--", "bash", "verify.sh"}

	var jsonOutput bytes.Buffer
	jsonStatus := commandruntime.NewRuntime(&jsonOutput, handlers).Run([]string{"--repo-root", repositoryRoot, "--format", "json", CommandCheckGreenGate, "--", "bash", "verify.sh"})
	if jsonStatus != 0 {
		t.Fatalf("JSON status=%d output=%s", jsonStatus, jsonOutput.String())
	}
	var jsonResult resultmodel.CommandResult
	if err := json.Unmarshal(jsonOutput.Bytes(), &jsonResult); err != nil {
		t.Fatal(err)
	}
	if jsonResult.GateEvidence == nil || jsonResult.GateEvidence.State != resultmodel.GateEvidenceMissing || jsonResult.GateEvidence.GateCommand == nil {
		t.Fatalf("JSON result = %#v", jsonResult)
	}

	var textOutput bytes.Buffer
	textStatus := commandruntime.NewRuntime(&textOutput, handlers).Run(append([]string{"--repo-root", repositoryRoot, CommandCheckGreenGate}, gateArguments...))
	if textStatus != jsonStatus {
		t.Fatalf("statuses text=%d JSON=%d", textStatus, jsonStatus)
	}
	for _, expected := range []string{"check-green-gate: success", "state=missing matches=false basis=none", "gate command: bash verify.sh", jsonResult.GateEvidence.RecordPath} {
		if !strings.Contains(textOutput.String(), expected) {
			t.Errorf("text output missing %q:\n%s", expected, textOutput.String())
		}
	}
}

func TestGateEvidenceArgumentParsersPreserveExactArgv(t *testing.T) {
	targetRevision, check, err := parseCheckArguments([]string{"--", "sh", "-c", "printf '%s' one"})
	if err != nil || targetRevision != "" || !equalArgv(check, []string{"sh", "-c", "printf '%s' one"}) {
		t.Fatalf("check=%#v target=%q err=%v", check, targetRevision, err)
	}
	for _, arguments := range [][]string{
		{"--at-revision", "abc123", "--", "bash", "verify.sh"},
		{"--at-revision=abc123", "--", "bash", "verify.sh"},
	} {
		spec, argv, parseError := parseCheckArguments(arguments)
		if parseError != nil || spec != "abc123" || !equalArgv(argv, []string{"bash", "verify.sh"}) {
			t.Errorf("check %#v -> spec=%q argv=%#v err=%v", arguments, spec, argv, parseError)
		}
	}
	status, record, err := parseRecordArguments([]string{"--gate-exit-status=0", "--", "sh", "-c", "printf '%s' one"})
	if err != nil || status != 0 || !equalArgv(record, check) {
		t.Fatalf("status=%d record=%#v err=%v", status, record, err)
	}
	for _, arguments := range [][]string{
		{},
		{"bash", "verify.sh"},
		{"--"},
		{"--gate-exit-status", "x", "--", "bash"},
		{"--at-revision", "--", "bash"},
		{"--at-revision=", "--", "bash"},
		{"--at-revision", "a", "--at-revision", "b", "--", "bash"},
		{"--at-revision", "a"},
	} {
		if _, _, err := parseCheckArguments(arguments); err == nil {
			t.Errorf("check accepted %#v", arguments)
		}
	}
	for _, arguments := range [][]string{
		{},
		{"--gate-exit-status", "x", "--", "bash"},
		{"--gate-exit-status", "0", "--gate-exit-status", "0", "--", "bash"},
		{"--gate-exit-status", "0"},
		{"--gate-exit-status", "0", "--"},
		{"--unknown", "0", "--", "bash"},
	} {
		if _, _, err := parseRecordArguments(arguments); err == nil {
			t.Errorf("record accepted %#v", arguments)
		}
	}
}
