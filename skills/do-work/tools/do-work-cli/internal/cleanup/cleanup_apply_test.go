package cleanup

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDryRunReportsExactArchiveMoveAndApplyPerformsIt(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-103-finished.md", cleanupRequest("REQ-103", "finished", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plan := BuildPlan(snapshot)
	dryResult := ApplyPlan(context.Background(), plan, ApplyOptions{DryRun: true})
	if dryResult.Outcome != resultmodel.OutcomeSuccess || len(dryResult.Changes) != 1 || !strings.Contains(dryResult.Changes[0].Detail, "do-work/queue/REQ-103-finished.md") {
		t.Fatalf("dry-run result = %#v", dryResult)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-103-finished.md")); err != nil {
		t.Fatalf("dry-run mutated source: %v", err)
	}
	applyResult := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if applyResult.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("apply result = %#v", applyResult)
	}
	archivedPath := filepath.Join(repositoryRoot, "do-work/archive/REQ-103-finished.md")
	contents, err := os.ReadFile(archivedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "status: completed") {
		t.Fatalf("terminal alias not normalized: %s", contents)
	}
}

func TestDirtyGroupIsRefusedWithoutBlockingIndependentSafeGroup(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-104-done.md", cleanupRequest("REQ-104", "done", ""))
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-105-done.md", cleanupRequest("REQ-105", "done", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plan := BuildPlan(snapshot)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-104-done.md", cleanupRequest("REQ-104", "done", "")+"user edit\n")
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/archive/REQ-105-done.md")); err != nil {
		t.Fatalf("safe group did not apply: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-104-done.md")); err != nil {
		t.Fatalf("dirty group moved: %v", err)
	}
}

func TestCleanupCommitContainsOnlyExactTouchedPaths(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-110-done.md", cleanupRequest("REQ-110", "completed", ""))
	writeCleanupFile(t, repositoryRoot, "unrelated.txt", "initial\n")
	commitCleanupFixture(t, repositoryRoot)
	writeCleanupFile(t, repositoryRoot, "unrelated.txt", "user dirt\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{Commit: true, CommitMessage: "cleanup exact move"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("commit result = %#v", result)
	}
	paths := strings.Fields(runCleanupGit(t, repositoryRoot, "show", "--no-renames", "--pretty=format:", "--name-only", "HEAD"))
	sort.Strings(paths)
	want := []string{"do-work/archive/REQ-110-done.md", "do-work/queue/REQ-110-done.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("commit paths = %#v, want %#v", paths, want)
	}
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, "unrelated.txt"))
	if err != nil || string(contents) != "user dirt\n" {
		t.Fatalf("unrelated dirt changed: %q, %v", contents, err)
	}
}

func TestMovePublicationNeverOverwritesAnExistingDestination(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCleanupFile(t, repositoryRoot, "source.txt", "source\n")
	writeCleanupFile(t, repositoryRoot, "archive/destination.txt", "user content\n")
	if err := moveWithoutOverwrite(repositoryRoot, "source.txt", "archive/destination.txt"); err == nil {
		t.Fatal("move overwrote an existing destination")
	}
	destinationContents, err := os.ReadFile(filepath.Join(repositoryRoot, "archive", "destination.txt"))
	if err != nil || string(destinationContents) != "user content\n" {
		t.Fatalf("destination = %q, %v", destinationContents, err)
	}
	sourceContents, err := os.ReadFile(filepath.Join(repositoryRoot, "source.txt"))
	if err != nil || string(sourceContents) != "source\n" {
		t.Fatalf("source = %q, %v", sourceContents, err)
	}
}

func TestRollbackDoesNotClaimAppliedChangesAndPreservesCanonicalRoots(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-205-done.md", cleanupRequest("REQ-205", "completed", ""))
	for _, rootPath := range []string{"do-work/archive", "do-work/working", "do-work/user-requests", "do-work/runs"} {
		if err := os.MkdirAll(filepath.Join(repositoryRoot, filepath.FromSlash(rootPath)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	commitCleanupFixture(t, repositoryRoot)
	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{
		{Code: "forced-rollback", Operations: []CleanupOperation{
			{Kind: OperationMove, SourcePath: "do-work/queue/REQ-205-done.md", DestinationPath: "do-work/archive/REQ-205-done.md"},
			{Kind: OperationDelete, SourcePath: "do-work/queue/missing-after-preflight.md"},
		}},
	}}
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeRolledBack {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	for _, change := range result.Changes {
		if strings.HasPrefix(change.Detail, "applied ") {
			t.Fatalf("rollback claimed applied change: %#v", result.Changes)
		}
	}
	for _, rootPath := range []string{"do-work", "do-work/queue", "do-work/working", "do-work/user-requests", "do-work/archive", "do-work/runs"} {
		if info, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(rootPath))); err != nil || !info.IsDir() {
			t.Fatalf("canonical root %s removed: %v", rootPath, err)
		}
	}
}

func TestConsumedUntrackedRunIsDeletedWithTruthfulNonRollbackEvidence(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/runs/spent/manifest.md", "Status: consumed\n")
	writeCleanupFile(t, repositoryRoot, "do-work/runs/spent/output.txt", "spent\n")
	writeCleanupFile(t, repositoryRoot, "README.md", "tracked\n")
	runCleanupGit(t, repositoryRoot, "add", "README.md")
	runCleanupGit(t, repositoryRoot, "commit", "-q", "-m", "fixture")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("untracked consumed cleanup = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/runs/spent")); !os.IsNotExist(err) {
		t.Fatalf("spent run remains: %v", err)
	}
	foundEvidence := false
	for _, change := range result.Changes {
		if strings.Contains(change.Detail, "non-rollback spent-scratch") {
			foundEvidence = true
		}
	}
	if !foundEvidence {
		t.Fatalf("non-rollback boundary not reported: %#v", result.Changes)
	}
}

func TestCommittedRiskPreservesExactRevertArgvAndDoesNotClaimApplied(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-211-done.md", cleanupRequest("REQ-211", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{Commit: true, CommitMessage: "cleanup", PostCommitVerify: func(context.Context, string) error {
		return errors.New("forced post-commit risk")
	}})
	if result.Outcome != resultmodel.OutcomeRisk {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	if len(result.Findings) == 0 || len(result.Findings[len(result.Findings)-1].NextArgv) != 3 || result.Findings[len(result.Findings)-1].NextArgv[0] != "git" || result.Findings[len(result.Findings)-1].NextArgv[1] != "revert" {
		t.Fatalf("revert evidence = %#v", result.Findings)
	}
	for _, change := range result.Changes {
		if strings.HasPrefix(change.Detail, "applied ") {
			t.Fatalf("risk claimed applied: %#v", result.Changes)
		}
	}
}
