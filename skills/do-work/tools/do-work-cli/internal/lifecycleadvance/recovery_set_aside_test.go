package lifecycleadvance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/finalization"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// REQ-515: the same `recover --assume-sole-authority` call that sets a REQ
// aside must not then release that REQ's claim. A released claim returns the
// REQ to `do-work/queue/`, the same run selects it again, a builder redoes the
// work, and the finalize tail refuses on the unfinished journal still on disk.
// The unrelated claim in the same run is still recovered, so the exclusion
// stays one REQ wide.
func TestRecoverPreservesTheClaimOfARequestFinalizationSetAside(t *testing.T) {
	repositoryRoot := seedSetAsideRecoveryFixture(t)
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}

	recovered := handleRecover(executionContext, []string{"--assume-sole-authority"})
	if recovered.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("recovery did not settle: %#v", recovered)
	}
	if !carriesSetAsideRecord(recovered.Finalizations, "REQ-730") {
		t.Fatalf("the fixture must produce a set-aside record for REQ-730: %#v", recovered.Finalizations)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", "REQ-730.md")); err != nil {
		t.Fatalf("the set-aside REQ lost its claim: %v", err)
	}
	queued, _ := filepath.Glob(filepath.Join(repositoryRoot, "do-work", "queue", "REQ-730*"))
	if len(queued) != 0 {
		t.Fatalf("the set-aside REQ was released back into the queue: %v", queued)
	}
	setAsideClaim, unrelatedClaim := recoveryClaimFor(t, recovered, "REQ-730"), recoveryClaimFor(t, recovered, "REQ-731")
	if setAsideClaim.Recovered {
		t.Fatalf("claim recovery ran on a REQ this run excluded: %#v", setAsideClaim)
	}
	if !unrelatedClaim.Recovered {
		t.Fatalf("one REQ's exclusion stopped an unrelated claim from recovering: %#v", unrelatedClaim)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", "REQ-731.md")); !os.IsNotExist(err) {
		t.Fatalf("the unrelated claim was not released: %v", err)
	}
}

func carriesSetAsideRecord(records []resultmodel.FinalizationResult, requestID string) bool {
	for _, record := range records {
		if record.RequestID != requestID {
			continue
		}
		for _, reasonCode := range record.ReasonCodes {
			if reasonCode == finalization.SetAsideReasonCode {
				return true
			}
		}
	}
	return false
}

func recoveryClaimFor(t *testing.T, result resultmodel.CommandResult, requestID string) resultmodel.RecoveryClaimResult {
	t.Helper()
	if result.Recovery == nil {
		t.Fatalf("recovery result is missing: %#v", result)
	}
	for _, claim := range result.Recovery.Claims {
		if claim.RequestID == requestID {
			return claim
		}
	}
	t.Fatalf("no claim record for %s: %#v", requestID, result.Recovery.Claims)
	return resultmodel.RecoveryClaimResult{}
}

// seedSetAsideRecoveryFixture leaves REQ-730 with an unfinished finalization
// journal whose primary commit keeps refusing, and REQ-731 claimed with nothing
// wrong. A `commit-msg` hook rejects only the finalization commit message, so
// the refusal is REQ-scoped, repeatable, and leaves the index clean.
func seedSetAsideRecoveryFixture(t *testing.T) string {
	t.Helper()
	repositoryRoot := t.TempDir()
	checkpointPath := "do-work/CHECKPOINT.md"
	writeAdvanceFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n"+
		"- REQ-730: Set-aside fixture — claimed now — writer: host:/repo\n"+
		"- REQ-731: Unrelated fixture — claimed now — writer: host:/repo\n")
	for _, requestID := range []string{"REQ-730", "REQ-731"} {
		writeAdvanceFile(t, repositoryRoot, "do-work/working/"+requestID+".md",
			"---\nid: "+requestID+"\ntitle: Fixture "+requestID+"\nstatus: claimed\nclaimed_at: 2026-09-04T12:00:00Z\ncommit:\n---\n\n"+
				"## Implementation Summary\n- `req-730.txt` (modified)\n")
	}
	writeAdvanceFile(t, repositoryRoot, "req-730.txt", "before\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Recovery Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "recovery@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")

	hookPath := filepath.Join(repositoryRoot, ".git", "hooks", "commit-msg")
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\ngrep -q 'finalize set-aside fixture' \"$1\" && exit 1\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeAdvanceFile(t, repositoryRoot, "req-730.txt", "after\n")

	requestPath := "do-work/working/REQ-730.md"
	manifest := map[string]any{
		"request_id": "REQ-730", "request_path": requestPath, "writer_label": "host:/repo",
		"transition": "complete", "terminal_status": "completed", "completed_at": "2026-09-04T13:00:00Z",
		"expected_request_sha256":    advanceFileDigest(t, repositoryRoot, requestPath),
		"expected_checkpoint_sha256": advanceFileDigest(t, repositoryRoot, checkpointPath),
		"commit_paths":               []string{requestPath, "do-work/archive/REQ-730.md", checkpointPath, "req-730.txt"},
		"commit_message":             "[REQ-730] finalize set-aside fixture", "provenance_mode": "primary_commit",
	}
	manifestPath := filepath.Join(t.TempDir(), "REQ-730-manifest.json")
	contents, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}
	finalized := finalization.Handlers()[finalization.CommandFinalize](executionContext, []string{"--manifest", manifestPath})
	if finalized.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("the fixture finalization was supposed to refuse: %#v", finalized)
	}
	journalPath := filepath.Join(repositoryRoot, ".git", "do-work-finalization", "REQ-730.json")
	if _, statError := os.Stat(journalPath); statError != nil {
		t.Fatalf("the fixture needs an unfinished journal on disk: %v", statError)
	}
	return repositoryRoot
}

func advanceFileDigest(t *testing.T, repositoryRoot, relativePath string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(contents)
	return hex.EncodeToString(digest[:])
}
