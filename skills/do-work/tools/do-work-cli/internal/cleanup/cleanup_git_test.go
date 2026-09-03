package cleanup

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDocumentationLinkRewritePreservesAnchorAndBareMention(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "docs/prime.md", "[lesson](../do-work/queue/REQ-106-done.md#lessons) and REQ-106-done.md\n")
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-106-done.md", cleanupRequest("REQ-106", "done", ""))
	commitCleanupFixture(t, repositoryRoot)
	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{{Code: "move", Operations: []CleanupOperation{{Kind: OperationMove, SourcePath: "do-work/queue/REQ-106-done.md", DestinationPath: "do-work/archive/REQ-106-done.md"}}}}}
	EnrichDocumentationLinks(context.Background(), &plan)
	if len(plan.Groups[0].Operations) != 2 || plan.Groups[0].Operations[1].Kind != OperationRewriteLinks {
		t.Fatalf("link operation missing: %#v", plan.Groups[0].Operations)
	}
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %s; findings = %#v", result.Outcome, result.Findings)
	}
	updated, readError := os.ReadFile(filepath.Join(repositoryRoot, "docs/prime.md"))
	if readError != nil {
		t.Fatal(readError)
	}
	if string(updated) != "[lesson](../do-work/archive/REQ-106-done.md#lessons) and REQ-106-done.md\n" {
		t.Fatalf("rewritten doc = %q", updated)
	}
}

func TestDocumentationRewritesComposeTwoMovesOnce(t *testing.T) {
	testCases := []struct {
		name            string
		dirtySource     string
		wantOutcome     resultmodel.CommandOutcome
		wantDocument    string
		wantFirstMoved  bool
		wantSecondMoved bool
	}{
		{
			name:            "first owner refused",
			dirtySource:     "do-work/queue/REQ-207.md",
			wantOutcome:     resultmodel.OutcomeFindings,
			wantDocument:    "[one](../do-work/queue/REQ-207.md#anchor) [two](../do-work/archive/REQ-208.md) and REQ-208.md\n",
			wantSecondMoved: true,
		},
		{
			name:           "second owner refused",
			dirtySource:    "do-work/queue/REQ-208.md",
			wantOutcome:    resultmodel.OutcomeFindings,
			wantDocument:   "[one](../do-work/archive/REQ-207.md#anchor) [two](../do-work/queue/REQ-208.md) and REQ-208.md\n",
			wantFirstMoved: true,
		},
		{
			name:            "all safe rewrites compose",
			wantOutcome:     resultmodel.OutcomeSuccess,
			wantDocument:    "[one](../do-work/archive/REQ-207.md#anchor) [two](../do-work/archive/REQ-208.md) and REQ-208.md\n",
			wantFirstMoved:  true,
			wantSecondMoved: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
			if testCase.dirtySource != "" {
				writeCleanupFile(t, repositoryRoot, testCase.dirtySource, cleanupRequest(filepath.Base(testCase.dirtySource)[:7], "done", "")+"user edit\n")
			}

			result := ApplyPlan(context.Background(), plan, ApplyOptions{})
			if result.Outcome != testCase.wantOutcome {
				t.Fatalf("outcome = %s, want %s; findings = %#v", result.Outcome, testCase.wantOutcome, result.Findings)
			}
			documentBytes, readError := os.ReadFile(filepath.Join(repositoryRoot, "docs/prime.md"))
			if readError != nil {
				t.Fatal(readError)
			}
			if string(documentBytes) != testCase.wantDocument {
				t.Fatalf("document = %q, want %q", documentBytes, testCase.wantDocument)
			}
			assertCleanupMoveState(t, repositoryRoot, "REQ-207.md", testCase.wantFirstMoved)
			assertCleanupMoveState(t, repositoryRoot, "REQ-208.md", testCase.wantSecondMoved)
		})
	}
}

func assertCleanupMoveState(t *testing.T, repositoryRoot, baseName string, wantMoved bool) {
	t.Helper()
	sourcePath := filepath.Join(repositoryRoot, "do-work", "queue", baseName)
	destinationPath := filepath.Join(repositoryRoot, "do-work", "archive", baseName)
	_, sourceError := os.Stat(sourcePath)
	_, destinationError := os.Stat(destinationPath)
	if wantMoved {
		if !os.IsNotExist(sourceError) || destinationError != nil {
			t.Fatalf("%s move state: source error=%v destination error=%v", baseName, sourceError, destinationError)
		}
		return
	}
	if sourceError != nil || !os.IsNotExist(destinationError) {
		t.Fatalf("%s retained state: source error=%v destination error=%v", baseName, sourceError, destinationError)
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
	writeCleanupFile(t, repositoryRoot, "do-work/archive/REQ-111-merged.md", cleanupRequest("REQ-111", "completed", ""))
	writeCleanupFile(t, repositoryRoot, "do-work/archive/REQ-112-unmerged.md", cleanupRequest("REQ-112", "completed", ""))
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

func TestMergedCleanBuilderWorktreeRequiresSettledRequestEvidence(t *testing.T) {
	testCases := []struct {
		name      string
		requestID string
		seed      func(*testing.T, string)
		wantState string
	}{
		{
			name:      "working request remains active",
			requestID: "REQ-114",
			seed: func(t *testing.T, repositoryRoot string) {
				writeCleanupFile(t, repositoryRoot, "do-work/working/REQ-114-active.md", cleanupRequest("REQ-114", "claimed", ""))
			},
			wantState: "working",
		},
		{
			name:      "request identity is absent",
			requestID: "REQ-115",
			seed: func(t *testing.T, repositoryRoot string) {
				writeCleanupFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# no request\n")
			},
			wantState: "absent",
		},
		{
			name:      "request identity is ambiguous",
			requestID: "REQ-116",
			seed: func(t *testing.T, repositoryRoot string) {
				writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-116-one.md", cleanupRequest("REQ-116", "pending", ""))
				writeCleanupFile(t, repositoryRoot, "do-work/archive/REQ-116-two.md", cleanupRequest("REQ-116", "completed", ""))
			},
			wantState: "ambiguous",
		},
		{
			name:      "request identity is malformed",
			requestID: "REQ-117",
			seed: func(t *testing.T, repositoryRoot string) {
				writeCleanupFile(t, repositoryRoot, "do-work/archive/REQ-117-malformed.md", "not frontmatter\n")
			},
			wantState: "malformed",
		},
		{
			name:      "request identity is unreadable",
			requestID: "REQ-118",
			seed: func(t *testing.T, repositoryRoot string) {
				writeCleanupFile(t, repositoryRoot, "do-work/archive/target.md", cleanupRequest("REQ-118", "completed", ""))
				linkPath := filepath.Join(repositoryRoot, "do-work", "archive", "REQ-118-linked.md")
				if err := os.Symlink("target.md", linkPath); err != nil {
					t.Fatal(err)
				}
			},
			wantState: "unreadable",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repositoryRoot := cleanupRepository(t)
			writeCleanupFile(t, repositoryRoot, "README.md", "fixture\n")
			testCase.seed(t, repositoryRoot)
			commitCleanupFixture(t, repositoryRoot)
			worktreeName := "worktree-agent-" + testCase.requestID + "-leftover"
			worktreePath := filepath.Join(t.TempDir(), worktreeName)
			runCleanupGit(t, repositoryRoot, "worktree", "add", "-q", "-b", worktreeName, worktreePath)

			changes, findings := ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{})
			if len(changes) != 0 || len(findings) != 1 || findings[0].Code != "WORKTREE-REQUIRES-CONSENT" {
				t.Fatalf("default cleanup changes=%#v findings=%#v", changes, findings)
			}
			if len(findings[0].AffectedIDs) != 1 || findings[0].AffectedIDs[0] != testCase.requestID {
				t.Fatalf("affected ids = %#v", findings[0].AffectedIDs)
			}
			if len(findings[0].Evidence) != 1 || !strings.Contains(findings[0].Evidence[0], "request_state="+testCase.wantState) {
				t.Fatalf("evidence = %#v", findings[0].Evidence)
			}
			if _, err := os.Stat(worktreePath); err != nil {
				t.Fatalf("worktree was removed without consent: %v", err)
			}
			if !strings.Contains(runCleanupGit(t, repositoryRoot, "branch", "--list", worktreeName), worktreeName) {
				t.Fatal("branch was removed without consent")
			}

			changes, findings = ApplyWorktreeRepairs(context.Background(), repositoryRoot, WorktreeRepairOptions{DiscardNames: []string{worktreeName}})
			if len(findings) != 0 || len(changes) != 1 {
				t.Fatalf("explicit discard changes=%#v findings=%#v", changes, findings)
			}
			if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
				t.Fatalf("explicitly discarded worktree remains: %v", err)
			}
		})
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

func TestWorktreeEnumerationHandlesNULNewlineDetachedAndAbsentConsent(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "README.md", "fixture\n")
	writeCleanupFile(t, repositoryRoot, "do-work/archive/REQ-209-detached.md", cleanupRequest("REQ-209", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)
	detachedPath := filepath.Join(t.TempDir(), "worktree-agent-REQ-209-detached\nline")
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
