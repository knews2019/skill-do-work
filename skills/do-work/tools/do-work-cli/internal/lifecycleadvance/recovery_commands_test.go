package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecoverPublicCommandRunsFinalizationThenRecoversEveryClaim(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-711", "pending", "", "# First\n")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-712", "pending", "", "# Second\n")
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	writeAdvanceFile(t, repositoryRoot, "project.txt", "base\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Recovery Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "recovery@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")

	for claimIndex, requestID := range []string{"REQ-711", "REQ-712"} {
		claimArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "claim", requestID,
			"--request-path", "do-work/queue/" + requestID + "-fixture.md", "--provenance", "explicit-req",
			"--writer", "hostile '$(touch should-not-exist)' ; writer", "--at", "2026-09-04T12:00:00Z"}
		if claimIndex == 0 {
			claimArguments = append(claimArguments, "--commit")
		}
		command := exec.Command(advanceCLIBinary(t), claimArguments...)
		command.Dir = repositoryRoot
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("claim %s: %v\n%s", requestID, err, output)
		}
	}
	checkpointPath := filepath.Join(repositoryRoot, "do-work", "CHECKPOINT.md")
	checkpoint, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatal(err)
	}
	checkpoint = []byte(strings.Replace(string(checkpoint),
		"- REQ-712: Fixture REQ-712 — claimed 2026-09-04T12:00:00Z — writer: hostile '$(touch should-not-exist)' ; writer",
		"- REQ-712: Fixture REQ-712 — claimed 2026-09-04T12:00:00Z", 1) +
		"- REQ-711: duplicate — claimed later — writer: another:/checkout\n" +
		"- REQ-999: unrelated — claimed earlier — writer: keep:/checkout\n")
	if err := os.WriteFile(checkpointPath, checkpoint, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repositoryRoot, "project.txt"), []byte("base\nunrelated dirt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "recover", "--assume-sole-authority")
	command.Dir = repositoryRoot
	output, commandError := command.CombinedOutput()
	if commandError != nil {
		t.Fatalf("recover: %v\n%s", commandError, output)
	}
	var result struct {
		Outcome  string `json:"outcome"`
		Recovery struct {
			FinalizationPassed bool `json:"finalization_passed"`
			Claims             []struct {
				RequestID string `json:"request_id"`
				Recovered bool   `json:"recovered"`
			} `json:"claims"`
		} `json:"recovery"`
		Findings []struct {
			Code        string   `json:"code"`
			AffectedIDs []string `json:"affected_ids"`
		} `json:"findings"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode recovery: %v\n%s", err, output)
	}
	if result.Outcome != "success" || !result.Recovery.FinalizationPassed || len(result.Recovery.Claims) != 2 || !result.Recovery.Claims[0].Recovered || !result.Recovery.Claims[1].Recovered {
		t.Fatalf("recovery result = %#v", result)
	}
	checkpointAfter, err := os.ReadFile(checkpointPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(checkpointAfter), "REQ-711") || strings.Contains(string(checkpointAfter), "REQ-712") || !strings.Contains(string(checkpointAfter), "REQ-999") {
		t.Fatalf("checkpoint recovery was incomplete:\n%s", checkpointAfter)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "should-not-exist")); !os.IsNotExist(err) {
		t.Fatalf("hostile writer label was evaluated: %v", err)
	}
	projectBytes, _ := os.ReadFile(filepath.Join(repositoryRoot, "project.txt"))
	if string(projectBytes) != "base\nunrelated dirt\n" {
		t.Fatalf("unrelated project dirt changed: %q", projectBytes)
	}

	for _, requestID := range []string{"REQ-711", "REQ-712"} {
		if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "queue", requestID+"-fixture.md")); err != nil {
			t.Fatalf("%s was not returned to the queue: %v", requestID, err)
		}
	}
	selection := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "next", "REQ-711")
	selectionOutput, err := selection.CombinedOutput()
	if err != nil || !strings.Contains(string(selectionOutput), `"request_id": "REQ-711"`) {
		t.Fatalf("recovered request was not selectable: %v\n%s", err, selectionOutput)
	}
	claim := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "claim", "REQ-711",
		"--request-path", "do-work/queue/REQ-711-fixture.md", "--provenance", "explicit-req", "--writer", "current:/checkout", "--at", "2026-09-04T13:00:00Z")
	if claimOutput, err := claim.CombinedOutput(); err != nil {
		t.Fatalf("fresh claim: %v\n%s", err, claimOutput)
	}
}

func TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-713", "claimed", "claimed_at: 2026-09-04T12:00:00Z\n", "# Claimed\n")
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-713: Claimed — claimed now — writer: other:/checkout\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Recovery Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "recovery@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
	before := advanceTreeDigest(t, repositoryRoot)
	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "recover")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("recover observation: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), `"code": "RECOVERY-TAKEOVER-AVAILABLE"`) ||
		!strings.Contains(string(output), `"--take-over"`) || !strings.Contains(string(output), `"REQ-713"`) {
		t.Fatalf("typed takeover missing:\n%s", output)
	}
	if after := advanceTreeDigest(t, repositoryRoot); before != after {
		t.Fatalf("authority-free recovery changed bytes")
	}
	takeOver := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "recover", "--take-over", "REQ-713")
	takeOver.Dir = repositoryRoot
	takeOverOutput, err := takeOver.CombinedOutput()
	if err != nil {
		t.Fatalf("authorized take-over: %v\n%s", err, takeOverOutput)
	}
	if !strings.Contains(string(takeOverOutput), `"authority_mode": "take-over"`) || !strings.Contains(string(takeOverOutput), `"recovered": true`) {
		t.Fatalf("typed take-over result missing:\n%s", takeOverOutput)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "queue", "REQ-713-fixture.md")); err != nil {
		t.Fatalf("authorized take-over did not return request to queue: %v", err)
	}
}
