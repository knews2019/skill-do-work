package gittransaction

import (
	"context"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestCommitExactPathsLeavesUnrelatedWorktreeChangesUntouched(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "owned.txt", "before\n")
	writeFile(t, repositoryRoot, "unrelated.txt", "before\n")
	commitAll(t, repositoryRoot, "seed")

	writeFile(t, repositoryRoot, "owned.txt", "after\n")
	writeFile(t, repositoryRoot, "unrelated.txt", "user change\n")
	result := CommitExactPaths(context.Background(), repositoryRoot, []string{"owned.txt"}, "exact commit", nil)
	if result.Outcome != resultmodel.OutcomeSuccess || result.CommitSHA == "" {
		t.Fatalf("exact commit = %#v", result)
	}
	if got := runFixtureGit(t, repositoryRoot, "show", "HEAD:owned.txt"); got != "after" {
		t.Fatalf("committed owned bytes = %q", got)
	}
	if got := readFile(t, repositoryRoot, "unrelated.txt"); got != "user change\n" {
		t.Fatalf("unrelated bytes = %q", got)
	}
	if status := runFixtureGit(t, repositoryRoot, "status", "--short"); status != "M unrelated.txt" {
		t.Fatalf("unrelated status = %q", status)
	}
}
