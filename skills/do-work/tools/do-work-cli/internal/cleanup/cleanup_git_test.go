package cleanup

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentationLinkRewritePreservesAnchorAndBareMention(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "docs/prime.md", "[lesson](../do-work/queue/REQ-106-done.md#lessons) and REQ-106-done.md\n")
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-106-done.md", cleanupRequest("REQ-106", "done", ""))
	commitCleanupFixture(t, repositoryRoot)
	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{{Code: "move", Operations: []CleanupOperation{{Kind: OperationMove, SourcePath: "do-work/queue/REQ-106-done.md", DestinationPath: "do-work/archive/REQ-106-done.md"}}}}}
	EnrichDocumentationLinks(context.Background(), &plan)
	if len(plan.Groups[0].Operations) != 2 {
		t.Fatalf("link operation missing: %#v", plan.Groups[0].Operations)
	}
	updated := string(plan.Groups[0].Operations[1].Contents)
	if !strings.Contains(updated, "../do-work/archive/REQ-106-done.md#lessons") || !strings.Contains(updated, "and REQ-106-done.md") {
		t.Fatalf("rewritten doc = %q", updated)
	}
}

func TestBlankedRestorationRequiresExactConsentAndUsesGitHistory(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	relativePath := "do-work/archive/REQ-107-done.md"
	writeCleanupFile(t, repositoryRoot, relativePath, cleanupRequest("REQ-107", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)
	if err := os.WriteFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	commitCleanupFixture(t, repositoryRoot)
	plan := CleanupPlan{RepositoryRoot: repositoryRoot}
	AddBlankedRepairs(context.Background(), &plan, nil)
	if len(plan.Groups) != 0 || len(plan.Findings) != 1 || plan.Findings[0].Code != "BLANKED-RECORD-REQUIRES-CONSENT" {
		t.Fatalf("default blank result = %#v", plan)
	}
	plan = CleanupPlan{RepositoryRoot: repositoryRoot}
	AddBlankedRepairs(context.Background(), &plan, []string{relativePath})
	if len(plan.Groups) != 1 || len(plan.Groups[0].Operations[0].Contents) == 0 {
		t.Fatalf("approved restore plan = %#v", plan)
	}
	if strings.Contains(string(plan.Groups[0].Operations[0].Contents), "commit:") {
		t.Fatal("restore fabricated commit provenance when the blanking commit recorded none")
	}
}

func TestMergedCleanBuilderWorktreeIsAutomaticButUnmergedNeedsExactConsent(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "README.md", "fixture\n")
	commitCleanupFixture(t, repositoryRoot)
	mergedName := "worktree-agent-REQ-111-merged"
	mergedPath := filepath.Join(t.TempDir(), mergedName)
	runCleanupGit(t, repositoryRoot, "worktree", "add", "-q", "-b", mergedName, mergedPath)
	changes, findings := ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{})
	if len(findings) != 0 || len(changes) != 1 {
		t.Fatalf("merged cleanup changes=%#v findings=%#v", changes, findings)
	}
	if _, err := os.Stat(mergedPath); !os.IsNotExist(err) {
		t.Fatalf("merged worktree remains: %v", err)
	}

	unmergedName := "worktree-agent-REQ-112-unmerged"
	unmergedPath := filepath.Join(t.TempDir(), unmergedName)
	runCleanupGit(t, repositoryRoot, "worktree", "add", "-q", "-b", unmergedName, unmergedPath)
	writeCleanupFile(t, unmergedPath, "builder.txt", "unmerged\n")
	commitCleanupFixture(t, unmergedPath)
	changes, findings = ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{})
	if len(changes) != 0 || len(findings) != 1 || findings[0].Code != "WORKTREE-REQUIRES-CONSENT" {
		t.Fatalf("unmerged default changes=%#v findings=%#v", changes, findings)
	}
	if _, err := os.Stat(unmergedPath); err != nil {
		t.Fatalf("unmerged worktree was deleted: %v", err)
	}
	changes, findings = ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{DiscardNames: []string{unmergedName}})
	if len(findings) != 0 || len(changes) != 1 {
		t.Fatalf("explicit discard changes=%#v findings=%#v", changes, findings)
	}
	if _, err := os.Stat(unmergedPath); !os.IsNotExist(err) {
		t.Fatalf("discarded worktree remains: %v", err)
	}
}

func TestDetachedAttributedWorktreeIsStillEnumerated(t *testing.T) {
	parsed := parseWorktrees([]byte("worktree /repo\x00HEAD abc\x00branch refs/heads/main\x00\x00worktree /tmp/worktree-agent-REQ-113-detached\x00HEAD def\x00detached\x00\x00"))
	if len(parsed) != 1 || parsed[0].Name != "worktree-agent-REQ-113-detached" || parsed[0].Path != "/tmp/worktree-agent-REQ-113-detached" {
		t.Fatalf("parsed worktrees = %#v", parsed)
	}
}

func TestBlankedRecoveryUsesRecordedImplementationHashAcrossAllLiveLayouts(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	relativePath := "do-work/queue/REQ-206-blanked.md"
	writeCleanupFile(t, repositoryRoot, relativePath, cleanupRequest("REQ-206", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)
	implementationSHA := runCleanupGit(t, repositoryRoot, "rev-parse", "HEAD")
	writeCleanupFile(t, repositoryRoot, relativePath, strings.ReplaceAll(cleanupRequest("REQ-206", "completed", ""), "Body", "Later parseable body"))
	commitCleanupFixture(t, repositoryRoot)
	recoverySourceSHA := runCleanupGit(t, repositoryRoot, "rev-parse", "HEAD")
	if implementationSHA == recoverySourceSHA {
		t.Fatal("fixture did not separate implementation and recovery-source commits")
	}
	if err := os.WriteFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	runCleanupGit(t, repositoryRoot, "add", relativePath)
	runCleanupGit(t, repositoryRoot, "commit", "-q", "-m", "[REQ-206] record commit hash "+implementationSHA)
	recoveryEvidence, recoveryError := RecoverGitContent(context.Background(), repositoryRoot, relativePath)
	if recoveryError != nil {
		t.Fatalf("RecoverGitContent: %v", recoveryError)
	}
	if recoveryEvidence.SourceCommit != recoverySourceSHA || recoveryEvidence.ImplementationCommit != implementationSHA || !strings.Contains(string(recoveryEvidence.ContentBytes), "Later parseable body") {
		t.Fatalf("reusable recovery evidence = %#v", recoveryEvidence)
	}
	plan := CleanupPlan{RepositoryRoot: repositoryRoot}
	AddBlankedRepairs(context.Background(), &plan, []string{relativePath})
	if len(plan.Groups) != 1 {
		t.Fatalf("live blanked record was not recovered: %#v", plan)
	}
	recovered := string(plan.Groups[0].Operations[0].Contents)
	if !strings.Contains(recovered, "Later parseable body") || !strings.Contains(recovered, "commit: "+implementationSHA[:7]) {
		t.Fatalf("recovered content/provenance = %q", recovered)
	}
}

func TestDocumentationRewritesComposeTwoMovesOnce(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "docs/prime.md", "[one](../do-work/queue/REQ-207.md#anchor) [two](../do-work/queue/REQ-208.md) and REQ-208.md\n")
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-207.md", cleanupRequest("REQ-207", "done", ""))
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-208.md", cleanupRequest("REQ-208", "done", ""))
	commitCleanupFixture(t, repositoryRoot)
	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{
		{Code: "one", Operations: []CleanupOperation{{Kind: OperationMove, SourcePath: "do-work/queue/REQ-207.md", DestinationPath: "do-work/archive/REQ-207.md"}}},
		{Code: "two", Operations: []CleanupOperation{{Kind: OperationMove, SourcePath: "do-work/queue/REQ-208.md", DestinationPath: "do-work/archive/REQ-208.md"}}},
	}}
	EnrichDocumentationLinks(context.Background(), &plan)
	var replacements []string
	for _, group := range plan.Groups {
		for _, operation := range group.Operations {
			if operation.Kind == OperationReplace {
				replacements = append(replacements, string(operation.Contents))
			}
		}
	}
	if len(replacements) != 1 || !strings.Contains(replacements[0], "archive/REQ-207.md#anchor") || !strings.Contains(replacements[0], "archive/REQ-208.md") || !strings.Contains(replacements[0], "and REQ-208.md") {
		t.Fatalf("composed replacements = %#v", replacements)
	}
}

func TestWorktreeEnumerationHandlesNULNewlineDetachedAndAbsentConsent(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "README.md", "fixture\n")
	commitCleanupFixture(t, repositoryRoot)
	detachedPath := filepath.Join(t.TempDir(), "worktree-agent-REQ-209 space\nline")
	runCleanupGit(t, repositoryRoot, "worktree", "add", "-q", "--detach", detachedPath, "HEAD")
	changes, findings := ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{})
	if len(findings) != 0 || len(changes) != 1 {
		t.Fatalf("detached merged worktree changes=%#v findings=%#v", changes, findings)
	}
	if _, err := os.Stat(detachedPath); !os.IsNotExist(err) {
		t.Fatalf("detached merged worktree remains: %v", err)
	}
	changes, findings = ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{DiscardNames: []string{"worktree-agent-REQ-999-absent"}})
	if len(changes) != 0 || len(findings) != 1 || findings[0].Code != "WORKTREE-DISCARD-NOT-FOUND" {
		t.Fatalf("absent discard changes=%#v findings=%#v", changes, findings)
	}
	_, findings = ApplyWorktreeRepairs(context.Background(), t.TempDir(), WorktreeRepairOptions{})
	if len(findings) != 1 || findings[0].Code != "WORKTREE-ENUMERATION-FAILED" {
		t.Fatalf("enumeration errors hidden: %#v", findings)
	}
}
