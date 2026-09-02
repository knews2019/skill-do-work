package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGreenGateEvidenceLifecycle(t *testing.T) {
	binaryPath := filepath.Join(t.TempDir(), "do-work-cli")
	buildCommand := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := buildCommand.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}

	repositoryRoot := t.TempDir()
	runGateEvidenceGit(t, repositoryRoot, "init", "-q")
	runGateEvidenceGit(t, repositoryRoot, "config", "user.email", "gate-evidence@example.invalid")
	runGateEvidenceGit(t, repositoryRoot, "config", "user.name", "Gate Evidence")
	writeGateEvidenceFixture(t, repositoryRoot, "project.txt", "initial\n")
	runGateEvidenceGit(t, repositoryRoot, "add", "project.txt")
	runGateEvidenceGit(t, repositoryRoot, "commit", "-qm", "initial")

	gateCommand := []string{"bash", "_dev/tests/maintainer-verify.sh"}
	recorded := runGateEvidenceCLI(t, binaryPath, repositoryRoot, "record-green-gate", append([]string{"--gate-exit-status", "0", "--"}, gateCommand...)...)
	if recorded.Outcome != "success" || recorded.GateEvidence.State != "recorded" {
		t.Fatalf("record result = %#v", recorded)
	}

	exact := runGateEvidenceCLI(t, binaryPath, repositoryRoot, "check-green-gate", append([]string{"--"}, gateCommand...)...)
	if exact.Outcome != "success" || !exact.GateEvidence.Matches || exact.GateEvidence.MatchBasis != "exact_revision" || exact.GateEvidence.BaselineRevision == "" {
		t.Fatalf("exact check = %#v", exact)
	}

	writeGateEvidenceFixture(t, repositoryRoot, "project.txt", "changed\n")
	runGateEvidenceGit(t, repositoryRoot, "commit", "-qam", "project change")
	changed := runGateEvidenceCLI(t, binaryPath, repositoryRoot, "check-green-gate", append([]string{"--"}, gateCommand...)...)
	if changed.Outcome != "success" || changed.GateEvidence.Matches || changed.GateEvidence.State != "invalidated_by_non_gate_log_commit" {
		t.Fatalf("project-change check = %#v", changed)
	}

	runGateEvidenceCLI(t, binaryPath, repositoryRoot, "record-green-gate", append([]string{"--gate-exit-status", "0", "--"}, gateCommand...)...)
	writeGateEvidenceFixture(t, repositoryRoot, "_dev/gate-runs/2026-09-03-full.log", "passed\n")
	runGateEvidenceGit(t, repositoryRoot, "add", "_dev/gate-runs/2026-09-03-full.log")
	runGateEvidenceGit(t, repositoryRoot, "commit", "-qm", "gate log")
	logOnly := runGateEvidenceCLI(t, binaryPath, repositoryRoot, "check-green-gate", append([]string{"--"}, gateCommand...)...)
	if logOnly.Outcome != "success" || !logOnly.GateEvidence.Matches || logOnly.GateEvidence.MatchBasis != "gate_log_only_descendant" || logOnly.GateEvidence.BaselineRevision != logOnly.GateEvidence.HeadRevision {
		t.Fatalf("log-only check = %#v", logOnly)
	}
}

type gateEvidenceCLIResult struct {
	Outcome      string `json:"outcome"`
	GateEvidence struct {
		State            string `json:"state"`
		Matches          bool   `json:"matches"`
		MatchBasis       string `json:"match_basis"`
		BaselineRevision string `json:"baseline_revision"`
		HeadRevision     string `json:"head_revision"`
	} `json:"gate_evidence"`
}

func runGateEvidenceCLI(t *testing.T, binaryPath, repositoryRoot, commandName string, commandArguments ...string) gateEvidenceCLIResult {
	t.Helper()
	arguments := []string{"--repo-root", repositoryRoot, "--format", "json", commandName}
	arguments = append(arguments, commandArguments...)
	command := exec.Command(binaryPath, arguments...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", commandName, err, output)
	}
	var result gateEvidenceCLIResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode %s result: %v\n%s", commandName, err, output)
	}
	return result
}

func runGateEvidenceGit(t *testing.T, repositoryRoot string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
}

func writeGateEvidenceFixture(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
