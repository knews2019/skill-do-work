package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
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

func legacyCheckpointRepository(t *testing.T) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	for _, id := range []string{"REQ-713", "REQ-799"} {
		writeAdvanceRequest(t, root, "working", id, "claimed", "claimed_at: 2026-09-04T12:00:00Z\n", "# Legacy claim\n")
	}
	entries := "- REQ-713: Legacy claim — claimed 2026-09-04T12:00:00Z — writer: other:/checkout\n  Keep detail.\n" +
		"- REQ-000713: Alias — claimed earlier — writer: second:/checkout\n\tAlias detail.\n" +
		"- REQ-713: Unlabelled — claimed earlier\n  Unlabelled detail.\n" +
		"- REQ-713: Duplicate — claimed earlier — writer: other:/checkout\n  Duplicate detail.\n"
	foreign := "- REQ-799: Foreign — claimed earlier — writer: keep:/checkout\n  Preserve foreign detail.\n"
	checkpoint := "# Session Checkpoint\n\n" + entries + foreign
	writeAdvanceFile(t, root, "do-work/CHECKPOINT.md", checkpoint)
	runAdvanceGit(t, root, "init", "-q")
	runAdvanceGit(t, root, "config", "user.name", "Legacy Recovery Test")
	runAdvanceGit(t, root, "config", "user.email", "legacy@example.invalid")
	runAdvanceGit(t, root, "add", ".")
	runAdvanceGit(t, root, "commit", "-qm", "fixture")
	return root, checkpoint, foreign
}

func runCheckpointPublicCommand(t *testing.T, root string, args ...string) resultmodel.CommandResult {
	t.Helper()
	argv := append([]string{"--repo-root", root, "--format", "json"}, args...)
	command := exec.Command(advanceCLIBinary(t), argv...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %v\n%s", args, err, output)
	}
	var result resultmodel.CommandResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode: %v\n%s", err, output)
	}
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("%v: %s", args, output)
	}
	return result
}

func TestRecoverLegacyCheckpointClaimsThroughPublicCommand(t *testing.T) {
	root, checkpoint, foreign := legacyCheckpointRepository(t)
	before := advanceTreeDigest(t, root)
	observation := runCheckpointPublicCommand(t, root, "recover")
	if len(observation.Recovery.Claims) != 2 || len(observation.Recovery.Claims[0].CheckpointEvidence) != 4 || observation.Recovery.Claims[0].CheckpointEvidence[0].Writer != "other:/checkout" {
		t.Fatalf("legacy evidence missing: %#v", observation.Recovery)
	}
	if advanceTreeDigest(t, root) != before {
		t.Fatal("observation changed bytes")
	}
	recovered := runCheckpointPublicCommand(t, root, "recover", "--take-over", "REQ-713")
	if len(recovered.Recovery.Claims) != 1 || !recovered.Recovery.Claims[0].Recovered {
		t.Fatalf("recovery = %#v", recovered.Recovery)
	}
	after, err := os.ReadFile(filepath.Join(root, "do-work/CHECKPOINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != "# Session Checkpoint\n\n"+foreign {
		t.Fatalf("recovery did not remove only matching blocks:\n%s\nOriginal:\n%s", after, checkpoint)
	}
	selection := runCheckpointPublicCommand(t, root, "next", "REQ-713")
	if len(selection.Selected) != 1 || selection.Selected[0].RequestID != "REQ-713" {
		t.Fatalf("recovered request not selectable: %#v", selection)
	}
	runCheckpointPublicCommand(t, root, "claim", "REQ-713", "--request-path", "do-work/queue/REQ-713-fixture.md", "--provenance", "explicit-req", "--writer", "current:/checkout", "--at", "2026-09-04T13:00:00Z", "--commit")
	refreshed := runCheckpointPublicCommand(t, root, "recover")
	if len(refreshed.Recovery.Claims) != 2 {
		t.Fatalf("fresh claims = %#v", refreshed.Recovery)
	}
	for _, claim := range refreshed.Recovery.Claims {
		if len(claim.CheckpointEvidence) != 1 {
			t.Fatalf("fresh claim hid evidence: %#v", claim)
		}
		wantWriter := "current:/checkout"
		if claim.RequestID == "REQ-799" {
			wantWriter = "keep:/checkout"
		}
		if claim.CheckpointEvidence[0].Writer != wantWriter {
			t.Fatalf("writer changed: %#v", claim)
		}
	}
	after, err = os.ReadFile(filepath.Join(root, "do-work/CHECKPOINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), foreign) {
		t.Fatalf("fresh claim changed foreign bytes:\n%s", after)
	}
}

// TestRecoverHoldsAClaimedRequestWithAHeavyVerificationPlanForTheDrain pins
// REQ-570: a claimed request in do-work/working/ carrying a Heavy Verification
// Plan section and a commit that is an ancestor of HEAD is held for the heavy
// lanes, so recovery preserves its claim and routes it to this session's drain.
// A claimed request without the section, and one whose commit never landed on
// this history, both recover to the queue as before.
func TestRecoverHoldsAClaimedRequestWithAHeavyVerificationPlanForTheDrain(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	writeAdvanceFile(t, repositoryRoot, "project.txt", "base\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Held Recovery Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "held@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
	landedCommit := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))

	runAdvanceGit(t, repositoryRoot, "checkout", "-q", "-b", "abandoned")
	writeAdvanceFile(t, repositoryRoot, "abandoned.txt", "abandoned\n")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "abandoned attempt")
	abandonedCommit := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "rev-parse", "HEAD")))
	runAdvanceGit(t, repositoryRoot, "checkout", "-q", "-")

	heavyPlanBody := "## Review\n\nReviewed.\n\n## Heavy Verification Plan\n\n- browser lane\n"
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-740", "claimed",
		"claimed_at: 2026-09-04T12:00:00Z\ncommit: "+landedCommit+"\n", heavyPlanBody)
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-741", "claimed",
		"claimed_at: 2026-09-04T12:00:00Z\ncommit: "+landedCommit+"\n", "## Review\n\nReviewed.\n")
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-742", "claimed",
		"claimed_at: 2026-09-04T12:00:00Z\ncommit: "+abandonedCommit+"\n", heavyPlanBody)
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "claims")

	recovered := handleRecover(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--assume-sole-authority"})
	if recovered.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("recovery did not settle: %#v", recovered)
	}

	heldClaim := recoveryClaimFor(t, recovered, "REQ-740")
	if !heldClaim.HeldForHeavyLanes || heldClaim.Recovered {
		t.Fatalf("held claim was not preserved for the drain: %#v", heldClaim)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", "REQ-740-fixture.md")); err != nil {
		t.Fatalf("held REQ lost its claim: %v", err)
	}
	if !carriesRecoveryFinding(recovered, "RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES", "REQ-740") {
		t.Fatalf("held REQ produced no typed finding: %#v", recovered.Findings)
	}

	for _, requestID := range []string{"REQ-741", "REQ-742"} {
		claim := recoveryClaimFor(t, recovered, requestID)
		if claim.HeldForHeavyLanes || !claim.Recovered {
			t.Fatalf("%s should recover as an ordinary interrupted claim: %#v", requestID, claim)
		}
		if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "queue", requestID+"-fixture.md")); err != nil {
			t.Fatalf("%s was not returned to the queue: %v", requestID, err)
		}
	}
}

func carriesRecoveryFinding(result resultmodel.CommandResult, code, requestID string) bool {
	for _, finding := range result.Findings {
		if finding.Code != code {
			continue
		}
		for _, affected := range finding.AffectedIDs {
			if affected == requestID {
				return true
			}
		}
	}
	return false
}
