package finalization

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestReviewPrimaryCommitFailurePreservesCommittedRisk(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires a Unix executable Git hook")
	}
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
	writeFinalizationFile(t, repositoryRoot, "foreign.txt", "foreign hook addition\n")
	installReviewPrimaryHook(t, repositoryRoot, "git add -- foreign.txt\n")
	before := currentHead(repositoryRoot)
	result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	committed := currentHead(repositoryRoot)
	if committed == before {
		t.Fatal("fixture did not create a commit before exact-path verification failed")
	}
	if got := runFinalizationGit(t, repositoryRoot, "show", committed+":foreign.txt"); got != "foreign hook addition\n" {
		t.Fatalf("hook addition was not committed: %q", got)
	}
	assertReviewCommittedRisk(t, repositoryRoot, committed, result)
}

func TestReviewPrimaryCommitRejectsHookContentChanges(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires a Unix executable Git hook")
	}
	for _, mutation := range []struct{ name, script string }{
		{"rewrite", "printf 'hook rewrite\\n' > implementation.txt\ngit add -- implementation.txt\n"},
		{"delete", "git rm -f -- implementation.txt\n"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
			installReviewPrimaryHook(t, repositoryRoot, mutation.script)
			before := currentHead(repositoryRoot)
			result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
			committed := currentHead(repositoryRoot)
			if committed == before {
				t.Fatal("fixture did not create the hook-mutated commit")
			}
			assertReviewCommittedRisk(t, repositoryRoot, committed, result)
		})
	}
}

func installReviewPrimaryHook(t *testing.T, repositoryRoot, script string) {
	t.Helper()
	hookPath := filepath.Join(repositoryRoot, ".git", "hooks", "pre-commit")
	if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
		t.Fatal(err)
	}
	// Remove the hook before it runs so a later provenance commit cannot be
	// responsible for the failure this fixture attributes to the primary commit.
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\nset -eu\nrm -- \"$0\"\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
}

func assertReviewCommittedRisk(t *testing.T, repositoryRoot, committed string, result resultmodel.CommandResult) {
	t.Helper()
	if result.Outcome == resultmodel.OutcomeSuccess || result.Outcome == resultmodel.OutcomeRolledBack || result.Rollback.Status == resultmodel.RollbackSucceeded {
		t.Fatalf("committed verification failure must remain a failure without pre-primary rollback: outcome=%s rollback=%s", result.Outcome, result.Rollback.Status)
	}
	if result.Finalization == nil || result.Finalization.PrimaryCommit != committed || result.Finalization.CreatedPrimaryCommit != committed || result.Finalization.Phase != string(PhaseReleaseApplied) {
		t.Fatalf("failed primary commit evidence = %#v, want SHA %s in release_applied", result.Finalization, committed)
	}
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-720")
	if err != nil {
		t.Fatal(err)
	}
	journal, err := readJournal(repositoryRoot, journalPath)
	if err != nil || journal.PrimaryCommit != committed || journal.Phase != PhaseReleaseApplied {
		t.Fatalf("persisted committed-risk evidence: journal=%#v error=%v", journal, err)
	}
	if status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("failed verification rolled committed lifecycle bytes back: %q", status)
	}
	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if len(recovered.Finalizations) != 1 || recovered.Finalizations[0].Phase != string(PhaseReleaseApplied) || recovered.Finalizations[0].PrimaryCommit != committed || len(recovered.Finalizations[0].ReasonCodes) == 0 || currentHead(repositoryRoot) != committed {
		t.Fatalf("recovery silently accepted or recommitted the failed primary: outcome=%s HEAD=%s", recovered.Outcome, currentHead(repositoryRoot))
	}
	if status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("recovery rolled back committed lifecycle bytes: %q", status)
	}
}

func TestReviewPreparedIdentityIncludesAdditionsAndDeletionsWithoutChangingIndex(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	for _, path := range []string{"edited.txt", "deleted.txt", "foreign.txt"} {
		writeFinalizationFile(t, repositoryRoot, path, "before\n")
	}
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "edited.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, "added.txt", "new untracked content\n")
	writeFinalizationFile(t, repositoryRoot, "foreign.txt", "unrelated staged content\n")
	if err := os.Remove(filepath.Join(repositoryRoot, "deleted.txt")); err != nil {
		t.Fatal(err)
	}
	runFinalizationGit(t, repositoryRoot, "add", "--", "foreign.txt")
	indexPath := filepath.Join(repositoryRoot, ".git", "index")
	indexBefore, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{"edited.txt", "added.txt", "deleted.txt", "absent-optional.txt"}
	head, digest, err := preparedCommitIdentity(repositoryRoot, paths)
	if err != nil {
		t.Fatal(err)
	}
	indexAfter, err := os.ReadFile(indexPath)
	if err != nil || !bytes.Equal(indexBefore, indexAfter) {
		t.Fatalf("preparation changed the real index: %v", err)
	}
	runFinalizationGit(t, repositoryRoot, "reset", "-q", "HEAD", "--", "foreign.txt")
	runFinalizationGit(t, repositoryRoot, append([]string{"add", "-A", "--"}, paths[:3]...)...)
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "prepared implementation")
	journal := &Journal{PreparedHead: head, PreparedDiffSHA256: digest, EffectiveCommitPaths: paths}
	if matched, ok := matchingHeadCommit(repositoryRoot, journal); !ok || matched != currentHead(repositoryRoot) {
		t.Fatalf("complete prepared identity failed to match committed edits, addition and deletion: matched=%q ok=%t", matched, ok)
	}
}

func TestReviewRecoveryRecognizesTrackedAndNewFilesBeforePrimaryPhasePersisted(t *testing.T) {
	repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(contents, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest.CommitPaths = append(manifest.CommitPaths, "new-implementation.txt")
	contents, err = json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	writeFinalizationFile(t, repositoryRoot, "new-implementation.txt", "new untracked implementation\n")
	journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := lifecyclePlan(repositoryRoot, journal.Manifest, journal.LifecyclePreimages)
	if err != nil {
		t.Fatal(err)
	}
	if applied := requeststate.ApplyPlan(context.Background(), plan); applied.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("apply lifecycle: %#v", applied)
	}
	journal.Phase = PhaseReleaseApplied
	journal.PreparedHead, journal.PreparedDiffSHA256, err = preparedCommitIdentity(repositoryRoot, journal.EffectiveCommitPaths)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeJournal(journal); err != nil {
		t.Fatal(err)
	}
	runFinalizationGit(t, repositoryRoot, append([]string{"add", "-A", "--"}, journal.EffectiveCommitPaths...)...)
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", journal.Manifest.CommitMessage)
	primary := currentHead(repositoryRoot)
	// The process died after Git committed, while the durable journal still
	// says release_applied and carries no primary SHA.
	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if recovered.Outcome != resultmodel.OutcomeSuccess || len(recovered.Finalizations) != 1 {
		t.Fatalf("recovery of tracked plus untracked implementation failed: %#v", recovered)
	}
	record := recovered.Finalizations[0]
	if record.PrimaryCommit != primary || record.CreatedPrimaryCommit != "" || record.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("recovery failed to reuse original primary: %#v", record)
	}
	if got := runFinalizationGit(t, repositoryRoot, "show", primary+":new-implementation.txt"); got != "new untracked implementation\n" {
		t.Fatalf("new implementation missing from recovered primary: %q", got)
	}
}
