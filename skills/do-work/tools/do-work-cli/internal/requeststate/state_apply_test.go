package requeststate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestLifecycleApplySynchronizesClaimUnblockCompleteFailAndCancel(t *testing.T) {
	now := "2026-08-31T21:00:00Z"
	t.Run("claim", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/queue/REQ-301.md", "REQ-301", "pending", "assigned_to: other\n")
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionClaim, []string{"REQ-301", "--request-path", "do-work/queue/REQ-301.md", "--provenance", "explicit-req", "--writer", "host:/repo", "--at", now})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/working/REQ-301.md")
		for _, token := range []string{"status: claimed", "claimed_at: 2026-08-31T21:00:00Z"} {
			if !strings.Contains(contents, token) {
				t.Errorf("claim missing %q:\n%s", token, contents)
			}
		}
		if strings.Contains(contents, "assigned_to:") {
			t.Fatal("explicit claim retained assignment")
		}
		if checkpoint := readStateFile(t, root, "do-work/CHECKPOINT.md"); !strings.Contains(checkpoint, "REQ-301") || !strings.Contains(checkpoint, "writer: host:/repo") {
			t.Fatalf("checkpoint not synchronized:\n%s", checkpoint)
		}
	})
	t.Run("unblock", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/queue/REQ-302.md", "REQ-302", "blocked", "blocked_by: service ready\nblocked_at: 2026-08-31T19:00:00Z\nblocked_check: test -e ready\n")
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionUnblock, []string{"REQ-302", "--request-path", "do-work/queue/REQ-302.md", "--original-status", "blocked", "--probe-status", "succeeded", "--unblock-required", "--source", "successful-probe", "--at", now})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/queue/REQ-302.md")
		for _, token := range []string{"status: pending", "blocked_check:", "cleared by probe"} {
			if !strings.Contains(contents, token) {
				t.Errorf("unblock missing %q:\n%s", token, contents)
			}
		}
		if strings.Contains(contents, "blocked_by:") || strings.Contains(contents, "blocked_at:") {
			t.Fatal("unblock retained blocked metadata")
		}
	})
	t.Run("complete closes UR and calibrates", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/working/REQ-303.md", "REQ-303", "claimed", "user_request: UR-081\nroute: C\nclaimed_at: 2026-08-31T20:00:00Z\nestimate:\n  p50_active_minutes: 60\n")
		writeStateRequest(t, root, "do-work/archive/REQ-300.md", "REQ-300", "completed", "user_request: UR-081\n")
		writeStateUserRequest(t, root, "UR-081")
		writeStateCheckpoint(t, root, "- REQ-303: Fixture — claimed 2026-08-31T20:00:00Z — writer: host:/repo\n- REQ-999: Foreign — claimed earlier — writer: other:/repo\n")
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionComplete, []string{"REQ-303", "--terminal-status", "completed", "--implementation-hash", "abcdef0", "--writer", "host:/repo", "--at", now})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/archive/UR-081/REQ-303.md")
		for _, token := range []string{"status: completed", "completed_at: 2026-08-31T21:00:00Z", "commit: abcdef0"} {
			if !strings.Contains(contents, token) {
				t.Errorf("complete missing %q:\n%s", token, contents)
			}
		}
		_ = readStateFile(t, root, "do-work/archive/UR-081/REQ-300.md")
		_ = readStateFile(t, root, "do-work/archive/UR-081/input.md")
		if _, err := os.Stat(filepath.Join(root, "do-work", "user-requests", "UR-081")); !os.IsNotExist(err) {
			t.Fatalf("active UR directory remains after closure: %v", err)
		}
		if checkpoint := readStateFile(t, root, "do-work/CHECKPOINT.md"); strings.Contains(checkpoint, "REQ-303") || !strings.Contains(checkpoint, "REQ-999") {
			t.Fatalf("checkpoint removal damaged foreign entry:\n%s", checkpoint)
		}
		if calibration := readStateFile(t, root, "do-work/calibration-log.tsv"); !strings.Contains(calibration, "REQ-303\tC\t60\t60\t2026-08-31T21:00:00Z") {
			t.Fatalf("calibration missing:\n%s", calibration)
		}
	})
	t.Run("fail keeps UR open", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/working/REQ-304.md", "REQ-304", "claimed", "user_request: UR-081\nclaimed_at: 2026-08-31T20:00:00Z\n")
		writeStateUserRequest(t, root, "UR-081")
		writeStateCheckpoint(t, root, "- REQ-304: Fixture — claimed now — writer: host:/repo\n")
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionFail, []string{"REQ-304", "--error", "tests failed", "--error-type", "code", "--writer", "host:/repo", "--at", now})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/archive/REQ-304.md")
		if !strings.Contains(contents, "status: failed") || !strings.Contains(contents, "error_type: code") {
			t.Fatalf("failure fields missing:\n%s", contents)
		}
		_ = readStateFile(t, root, "do-work/user-requests/UR-081/input.md")
	})
	t.Run("cancel archived failure in place", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/archive/REQ-305.md", "REQ-305", "failed", "completed_at: 2026-08-31T20:00:00Z\nerror: unavailable\nerror_type: environment\n")
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionCancel, []string{"REQ-305", "--confirmed", "--dependent-disposition", "leave", "--reason", "not retrying", "--at", now})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/archive/REQ-305.md")
		for _, token := range []string{"status: cancelled", "error_type: environment", "## Cancelled", "not retrying", "- **Previously:** failed (`error_type: environment`) — failed at 2026-08-31T20:00:00Z — resolved by decision not to retry"} {
			if !strings.Contains(contents, token) {
				t.Errorf("cancel missing %q:\n%s", token, contents)
			}
		}
	})
}

func TestCheckpointClaimRemovalIncludesIndentedContinuationLines(t *testing.T) {
	existing := []byte("# Session Checkpoint\n\n## In Progress (interrupted)\n\n" +
		"- REQ-489: Own claim — claimed now — writer: host:/repo\n" +
		"  - **Last known state:** implementation in progress\n" +
		"  - **Key files being modified:** `state_apply.go`\n" +
		"  - **Known issues:** removal leaves orphaned lines\n" +
		"- REQ-999: Foreign claim — claimed earlier — writer: other:/repo\n" +
		"  - **Last known state:** must remain\n")
	want := "# Session Checkpoint\n\n## In Progress (interrupted)\n\n" +
		"- REQ-999: Foreign claim — claimed earlier — writer: other:/repo\n" +
		"  - **Last known state:** must remain\n"

	if got := string(checkpointWithoutClaim(existing, "REQ-489", "host:/repo")); got != want {
		t.Fatalf("claim removal left continuation lines or damaged the foreign entry:\n%s", got)
	}
}

func TestCheckpointClaimUsesWholeInProgressHeadingLine(t *testing.T) {
	existing := []byte("# Session Checkpoint\n\n## Session Notes\n\n" +
		"- REQ-489: Historical note mentioning `## In Progress (interrupted)` — writer: host:/repo\n" +
		"  - This continuation must remain.\n\n" +
		"## In Progress (interrupted)\n\n" +
		"- REQ-999: Foreign claim — claimed earlier — writer: other:/repo\n" +
		"  - Foreign continuation must remain.\n")
	wantClaimed := string(existing) + "\n- REQ-489: Whole-heading fix — claimed 2026-09-02T10:00:00Z — writer: host:/repo\n"

	claimed := checkpointWithClaim(existing, "REQ-489", "Whole-heading fix", "2026-09-02T10:00:00Z", "host:/repo")
	if got := string(claimed); got != wantClaimed {
		t.Fatalf("claim was not inserted under the whole In Progress heading:\n%s", got)
	}
	departed := string(checkpointWithoutClaim(claimed, "REQ-489", "host:/repo"))
	sectionParts := strings.SplitN(departed, "\n## In Progress (interrupted)\n", 2)
	if len(sectionParts) != 2 || !strings.Contains(sectionParts[0], "- REQ-489: Historical note mentioning `## In Progress (interrupted)` — writer: host:/repo\n  - This continuation must remain.") {
		t.Fatalf("departure removed content outside the whole In Progress section:\n%s", departed)
	}
	if strings.Contains(sectionParts[1], "REQ-489") || !strings.Contains(sectionParts[1], "  - Foreign continuation must remain.") {
		t.Fatalf("departure did not remove only the real claim entry:\n%s", departed)
	}
}

func TestCancelContainsMultilineReasonAndPreservesOptionalFailureHistory(t *testing.T) {
	t.Run("multiline reason uses summary and byte-identical fenced blockquote", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/queue/REQ-306.md", "REQ-306", "pending", "")
		reason := "first line\n## injected\n````\n- [ ] fake"
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionCancel,
			[]string{"REQ-306", "--confirmed", "--dependent-disposition", "leave", "--reason-summary", "superseded after review", "--reason", reason, "--at", "2026-08-31T21:00:00Z"})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/archive/REQ-306.md")
		contained := "- **Why:** superseded after review\n> `````\n> first line\n> ## injected\n> ````\n> - [ ] fake\n> `````\n- **Decided by:** user, via `do-work abandon`"
		if !strings.Contains(contents, contained) {
			t.Fatalf("multiline reason was not contained byte-identically:\n%s", contents)
		}
	})
	t.Run("legacy failed record omits absent type and names absent instant", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/archive/REQ-307.md", "REQ-307", "failed", "error: legacy failure\n")
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionCancel,
			[]string{"REQ-307", "--confirmed", "--dependent-disposition", "leave", "--reason", "not retrying", "--at", "2026-08-31T21:00:00Z"})
		assertStateSuccess(t, result)
		contents := readStateFile(t, root, "do-work/archive/REQ-307.md")
		want := "- **Previously:** failed — failure instant unrecorded — resolved by decision not to retry"
		if !strings.Contains(contents, want) || strings.Contains(contents, "`error_type:") {
			t.Fatalf("optional failed history is wrong:\n%s", contents)
		}
	})
}

func TestLifecycleDryRunIsByteIdentical(t *testing.T) {
	root := newStateRepository(t)
	writeStateRequest(t, root, "do-work/queue/REQ-310.md", "REQ-310", "pending", "")
	original := readStateFile(t, root, "do-work/queue/REQ-310.md")
	result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionClaim, []string{"REQ-310", "--dry-run", "--writer", "host:/repo", "--at", "2026-08-31T21:00:00Z"})
	assertStateSuccess(t, result)
	if got := readStateFile(t, root, "do-work/queue/REQ-310.md"); got != original {
		t.Fatal("dry-run changed request bytes")
	}
	if _, err := os.Stat(filepath.Join(root, "do-work", "working", "REQ-310.md")); !os.IsNotExist(err) {
		t.Fatalf("dry-run created destination: %v", err)
	}
}

func TestLifecycleStalePreimageRollsBackAndDirtyTargetRefuses(t *testing.T) {
	t.Run("stale untracked preimage rolls back", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/queue/REQ-311.md", "REQ-311", "pending", "")
		snapshot, err := discoverRepository(root)
		if err != nil {
			t.Fatal(err)
		}
		plan := BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), StateOptions{Transition: TransitionClaim, RequestID: "REQ-311", WriterLabel: "host:/repo", Now: time.Date(2026, 8, 31, 21, 0, 0, 0, time.UTC)})
		changed := readStateFile(t, root, "do-work/queue/REQ-311.md") + "\nchanged after planning\n"
		if err := os.WriteFile(filepath.Join(root, "do-work/queue/REQ-311.md"), []byte(changed), 0o644); err != nil {
			t.Fatal(err)
		}
		result := ApplyPlan(context.Background(), plan)
		if result.Outcome != resultmodel.OutcomeRolledBack || result.Rollback.Status != resultmodel.RollbackSucceeded {
			t.Fatalf("stale apply = %#v", result)
		}
		if got := readStateFile(t, root, "do-work/queue/REQ-311.md"); got != changed {
			t.Fatal("rollback did not restore the exact pre-transaction bytes")
		}
	})
	t.Run("tracked dirty target refuses", func(t *testing.T) {
		root := newStateRepository(t)
		writeStateRequest(t, root, "do-work/queue/REQ-312.md", "REQ-312", "pending", "")
		configureStateGit(t, root)
		runStateGit(t, root, "add", "do-work/queue/REQ-312.md")
		runStateGit(t, root, "commit", "-q", "-m", "seed")
		path := filepath.Join(root, "do-work/queue/REQ-312.md")
		if err := os.WriteFile(path, []byte(readStateFile(t, root, "do-work/queue/REQ-312.md")+"dirty\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionClaim, []string{"REQ-312", "--writer", "host:/repo", "--at", "2026-08-31T21:00:00Z"})
		if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) == 0 || result.Findings[0].Code != "GIT-DIRTY-TARGET" {
			t.Fatalf("dirty target = %#v", result)
		}
	})
}

func TestSerialCompleteCreatesLifecycleAndMetadataCommits(t *testing.T) {
	root := newStateRepository(t)
	configureStateGit(t, root)
	writeStateRequest(t, root, "do-work/working/REQ-313.md", "REQ-313", "claimed", "claimed_at: 2026-08-31T20:00:00Z\n")
	result := handleStateCommand(commandruntime.ExecutionContext{RepositoryRoot: root}, TransitionComplete,
		[]string{"REQ-313", "--request-path", "do-work/working/REQ-313.md", "--terminal-status", "completed", "--commit", "--writer", "host:/repo", "--at", "2026-08-31T21:00:00Z"})
	assertStateSuccess(t, result)
	logLines := strings.Split(strings.TrimSpace(runStateGit(t, root, "log", "--format=%s")), "\n")
	if len(logLines) != 2 || !strings.Contains(logLines[0], "record commit hash") || !strings.Contains(logLines[1], "complete request lifecycle") {
		t.Fatalf("serial history = %q", logLines)
	}
	lifecycleHash := strings.TrimSpace(runStateGit(t, root, "rev-parse", "HEAD^"))
	contents := readStateFile(t, root, "do-work/archive/REQ-313.md")
	if !strings.Contains(contents, "commit: "+lifecycleHash) {
		t.Fatalf("archived REQ does not record lifecycle hash %s:\n%s", lifecycleHash, contents)
	}
	if status := runStateGit(t, root, "status", "--porcelain"); status != "" {
		t.Fatalf("serial completion left a dirty tree: %q", status)
	}
}

func TestRecordCommitProvenanceChangesOnlyTopLevelScalar(t *testing.T) {
	root := newStateRepository(t)
	configureStateGit(t, root)
	relativePath := "do-work/archive/REQ-314.md"
	writeStateRequest(t, root, relativePath, "REQ-314", "completed", "completed_at: 2026-08-31T21:00:00Z\ncommit: oldhash0\n")
	runStateGit(t, root, "add", relativePath)
	runStateGit(t, root, "commit", "-qm", "fixture")
	hash := strings.TrimSpace(runStateGit(t, root, "rev-parse", "HEAD"))
	before := readStateFile(t, root, relativePath)
	result := RecordCommitProvenance(context.Background(), root, relativePath, hash, false, false)
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("result=%#v", result)
	}
	after := readStateFile(t, root, relativePath)
	if !strings.Contains(after, "commit: "+hash) || strings.Replace(before, "commit: oldhash0", "commit: "+hash, 1) != after {
		t.Fatalf("unexpected rewrite:\n%s", after)
	}
}

func TestRecordCommitProvenanceVerifyProvesExactMetadataPatch(t *testing.T) {
	root := newStateRepository(t)
	configureStateGit(t, root)
	relativePath := "do-work/archive/REQ-315.md"
	writeStateRequest(t, root, relativePath, "REQ-315", "completed", "completed_at: 2026-08-31T21:00:00Z\ncommit: oldhash0\n")
	runStateGit(t, root, "add", relativePath)
	runStateGit(t, root, "commit", "-qm", "implementation")
	hash := strings.TrimSpace(runStateGit(t, root, "rev-parse", "HEAD"))
	writeResult := RecordCommitProvenance(context.Background(), root, relativePath, hash, false, false)
	if writeResult.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("write result=%#v", writeResult)
	}
	runStateGit(t, root, "add", relativePath)
	runStateGit(t, root, "commit", "-qm", "metadata")
	verifyResult := RecordCommitProvenance(context.Background(), root, relativePath, hash, true, false)
	if verifyResult.Outcome != resultmodel.OutcomeSuccess || !hasStateFinding(verifyResult, "PROVENANCE-VERIFIED") {
		t.Fatalf("verify result=%#v", verifyResult)
	}
}

func TestRecordCommitProvenanceVerifyRejectsRestagedBodyRewrite(t *testing.T) {
	root := newStateRepository(t)
	configureStateGit(t, root)
	relativePath := "do-work/archive/REQ-316.md"
	writeStateRequest(t, root, relativePath, "REQ-316", "completed", "completed_at: 2026-08-31T21:00:00Z\ncommit: oldhash0\n")
	runStateGit(t, root, "add", relativePath)
	runStateGit(t, root, "commit", "-qm", "implementation")
	hash := strings.TrimSpace(runStateGit(t, root, "rev-parse", "HEAD"))
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	corrupted := strings.Replace(readStateFile(t, root, relativePath), "commit: oldhash0", "commit: "+hash, 1)
	corrupted = strings.Replace(corrupted, "Body\n", "Body rewritten by hook\n", 1)
	if err := os.WriteFile(path, []byte(corrupted), 0o644); err != nil {
		t.Fatal(err)
	}
	runStateGit(t, root, "add", relativePath)
	runStateGit(t, root, "commit", "-qm", "metadata with restaged rewrite")
	result := RecordCommitProvenance(context.Background(), root, relativePath, hash, true, false)
	if result.Outcome == resultmodel.OutcomeSuccess || !hasStateFinding(result, "PROVENANCE-VERIFY-PATCH") {
		t.Fatalf("unsafe metadata commit verified: %#v", result)
	}
}

func TestRecordCommitProvenanceDryRunLeavesBytesUnchanged(t *testing.T) {
	root := newStateRepository(t)
	configureStateGit(t, root)
	relativePath := "do-work/archive/REQ-317.md"
	writeStateRequest(t, root, relativePath, "REQ-317", "completed", "completed_at: 2026-08-31T21:00:00Z\ncommit: oldhash0\n")
	runStateGit(t, root, "add", relativePath)
	runStateGit(t, root, "commit", "-qm", "fixture")
	hash := strings.TrimSpace(runStateGit(t, root, "rev-parse", "HEAD"))
	before := readStateFile(t, root, relativePath)
	result := RecordCommitProvenance(context.Background(), root, relativePath, hash, false, true)
	if result.Outcome != resultmodel.OutcomeSuccess || !hasStateFinding(result, "PROVENANCE-DRY-RUN") {
		t.Fatalf("result=%#v", result)
	}
	if after := readStateFile(t, root, relativePath); after != before {
		t.Fatalf("dry-run changed bytes:\n%s", after)
	}
}

func TestRecordCommitProvenanceFailsClosedWhenGitCannotAnswerGuardQueries(t *testing.T) {
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name      string
		condition string
	}{
		{"unreadable HEAD blob size", `[ "$3" = "cat-file" ] && [ "$4" = "-s" ]`},
		{"failed pending numstat", `[ "$3" = "diff" ] && [ "$4" = "--numstat" ]`},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := newStateRepository(t)
			configureStateGit(t, root)
			relativePath := "do-work/archive/REQ-318.md"
			writeStateRequest(t, root, relativePath, "REQ-318", "completed", "completed_at: 2026-08-31T21:00:00Z\ncommit: oldhash0\n")
			runStateGit(t, root, "add", relativePath)
			runStateGit(t, root, "commit", "-qm", "fixture")
			hash := strings.TrimSpace(runStateGit(t, root, "rev-parse", "HEAD"))
			before := readStateFile(t, root, relativePath)
			stubDirectory := t.TempDir()
			stub := fmt.Sprintf("#!/bin/sh\nif %s; then exit 73; fi\nexec \"%s\" \"$@\"\n", test.condition, realGit)
			if err := os.WriteFile(filepath.Join(stubDirectory, "git"), []byte(stub), 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", stubDirectory+string(os.PathListSeparator)+os.Getenv("PATH"))
			result := RecordCommitProvenance(context.Background(), root, relativePath, hash, false, false)
			if result.Outcome == resultmodel.OutcomeSuccess || !hasStateFinding(result, "PROVENANCE-PREIMAGE-GUARD") {
				t.Fatalf("Git guard failure was accepted: %#v", result)
			}
			if after := readStateFile(t, root, relativePath); after != before {
				t.Fatalf("Git guard failure changed bytes:\n%s", after)
			}
		})
	}
}

func hasStateFinding(result resultmodel.CommandResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

func newStateRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	command := exec.Command("git", "init", "-q", root)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	return root
}

func configureStateGit(t *testing.T, root string) {
	t.Helper()
	runStateGit(t, root, "config", "user.name", "Request State Fixture")
	runStateGit(t, root, "config", "user.email", "request-state@example.invalid")
}

func runStateGit(t *testing.T, root string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
	return string(output)
}
func writeStateUserRequest(t *testing.T, root, id string) {
	t.Helper()
	path := filepath.Join(root, "do-work", "user-requests", id, "input.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("---\nid: "+id+"\ntitle: Fixture\nstatus: active\n---\nBody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
func writeStateCheckpoint(t *testing.T, root, entries string) {
	t.Helper()
	path := filepath.Join(root, "do-work", "CHECKPOINT.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# Session Checkpoint\n\n## In Progress (interrupted)\n\n"+entries), 0o644); err != nil {
		t.Fatal(err)
	}
}
func readStateFile(t *testing.T, root, path string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}
func assertStateSuccess(t *testing.T, result resultmodel.CommandResult) {
	t.Helper()
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("result = %#v", result)
	}
}

var _ = time.UTC
