package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAdvanceCheckpointRemovesStaleSummariesAndPreservesLiveEntries(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-714", "pending", "", "")
	// Completion refreshed the count to seven but left a 28-request summary and
	// an obsolete next request. Keep live claims and authored notes, not that cache.
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "---\nlast_completed: REQ-483\nreqs_processed_this_session: 1\nsession_depth: light\nqueue_state: [28 pending]\ncustom_note: keep me\n---\n\n# Session Checkpoint\n\n## Completed This Session\n\n- REQ-483: old completion\n\n## Still Queued\n\n- 28 pending requests remain; REQ-485 next.\n\n## Operator Notes\n\nKeep this context.\n\n## In Progress (interrupted)\n\n- REQ-800: foreign — writer: other:/checkout\n  keep foreign detail\n- REQ-801: unknown owner\n  keep unknown detail\n\n## Session Notes\n\nOld note.\n")
	writeAdvanceFile(t, repositoryRoot, "project.txt", "base\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Checkpoint Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "checkpoint@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
	if err := os.WriteFile(filepath.Join(repositoryRoot, "project.txt"), []byte("base\nforeign dirt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "advance", "--checkpoint")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("advance checkpoint: %v\n%s", err, output)
	}
	var result struct {
		Checkpoint struct {
			CheckpointPath  string `json:"checkpoint_path"`
			PreservedClaims int    `json:"preserved_claims"`
		} `json:"checkpoint"`
		Changes []struct {
			Path string `json:"path"`
		} `json:"changes"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode checkpoint result: %v\n%s", err, output)
	}
	if result.Checkpoint.CheckpointPath != "do-work/CHECKPOINT.md" || result.Checkpoint.PreservedClaims != 2 || len(result.Changes) != 1 || result.Changes[0].Path != "do-work/CHECKPOINT.md" {
		t.Fatalf("checkpoint result did not expose the exact mutation: %#v", result)
	}
	checkpoint, _ := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "CHECKPOINT.md"))
	for _, exact := range []string{"- REQ-800: foreign — writer: other:/checkout\n  keep foreign detail", "- REQ-801: unknown owner\n  keep unknown detail", "custom_note: keep me", "## Operator Notes\n\nKeep this context.", "## Session Notes\n\nOld note.", "queue_state: [1 pending,"} {
		if !strings.Contains(string(checkpoint), exact) {
			t.Fatalf("checkpoint lost live record %q:\n%s", exact, checkpoint)
		}
	}
	for _, stale := range []string{"last_completed:", "reqs_processed_this_session:", "session_depth:", "## Completed This Session", "## Still Queued", "28 pending", "REQ-483", "REQ-485 next"} {
		if strings.Contains(string(checkpoint), stale) {
			t.Fatalf("refresh retained stale summary %q:\n%s", stale, checkpoint)
		}
	}
	project, _ := os.ReadFile(filepath.Join(repositoryRoot, "project.txt"))
	if string(project) != "base\nforeign dirt\n" {
		t.Fatalf("project dirt changed: %q", project)
	}
	status := string(runAdvanceGit(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=all"))
	if status != " M do-work/CHECKPOINT.md\n M project.txt\n" {
		t.Fatalf("unexpected changed paths:\n%s", status)
	}
}

func TestWorkingAdvanceRemainsReadOnlyAfterCheckpointMode(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-715", "claimed", "", "")
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	before := advanceTreeDigest(t, repositoryRoot)
	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "advance", "REQ-715")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("ordinary advance: %v\n%s", err, output)
	}
	if after := advanceTreeDigest(t, repositoryRoot); before != after {
		t.Fatalf("ordinary advance changed bytes")
	}
}

func TestAdvanceCheckpointPreservesLegacyClaimDiscovery(t *testing.T) {
	root, checkpoint, _ := legacyCheckpointRepository(t)
	before := runCheckpointPublicCommand(t, root, "recover")
	refreshed := runCheckpointPublicCommand(t, root, "advance", "--checkpoint")
	after := runCheckpointPublicCommand(t, root, "recover")
	count := 0
	for i := range after.Recovery.Claims {
		for j := range after.Recovery.Claims[i].CheckpointEvidence {
			after.Recovery.Claims[i].CheckpointEvidence[j].SourceLine = 0
			count++
		}
	}
	for i := range before.Recovery.Claims {
		for j := range before.Recovery.Claims[i].CheckpointEvidence {
			before.Recovery.Claims[i].CheckpointEvidence[j].SourceLine = 0
		}
	}
	if !reflect.DeepEqual(before.Recovery.Claims, after.Recovery.Claims) || refreshed.Checkpoint.PreservedClaims != count || count != 5 {
		t.Fatalf("refresh hid legacy evidence: before=%#v after=%#v preserved=%d observed=%d", before.Recovery.Claims, after.Recovery.Claims, refreshed.Checkpoint.PreservedClaims, count)
	}
	contents, err := os.ReadFile(filepath.Join(root, "do-work/CHECKPOINT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(contents), checkpoint) {
		t.Fatalf("refresh changed legacy body:\n%s", contents)
	}
}

func TestCheckpointSummaryRemovalDoesNotDependOnSectionOrder(t *testing.T) {
	claims := "## In Progress (interrupted)\n\n- REQ-800: foreign — writer: other:/checkout\n  keep foreign detail\n\n"
	notes := "## Session Notes\n\nKeep this authored note.\n"
	summaries := "## Completed This Session\n\n- REQ-483: old completion\n\n## Still Queued\n\n- 28 pending requests remain; REQ-485 next.\n\n"
	for _, order := range []string{"before claims", "after claims"} {
		t.Run(order, func(t *testing.T) {
			body := summaries + claims + notes
			if order == "after claims" {
				body = claims + summaries + notes
			}
			refreshed := string(checkpointSessionBytes([]byte("# Session Checkpoint\n\n"+body), "2026-09-13T12:00:00Z", "[1 pending]"))
			if !strings.Contains(refreshed, claims) || !strings.Contains(refreshed, notes) || !strings.Contains(refreshed, "queue_state: [1 pending]") {
				t.Fatalf("refresh lost claims or notes:\n%s", refreshed)
			}
			for _, stale := range []string{"## Completed This Session", "## Still Queued", "REQ-483", "28 pending"} {
				if strings.Contains(refreshed, stale) {
					t.Errorf("refresh retained %q:\n%s", stale, refreshed)
				}
			}
		})
	}
}
