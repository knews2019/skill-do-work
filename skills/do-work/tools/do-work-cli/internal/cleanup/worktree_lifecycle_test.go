package cleanup

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const lifecycleOperativeName = "worktree-agent-REQ-41-lifecycle-demo"

// lifecycleRepository is a main tree at <tempdir>/repo, so the sibling repo-worktrees/
// directory `worktree new` creates also lands inside the test's temporary directory. It
// holds a queued REQ-41, a link config naming one existing ignored directory and one
// missing path, and one commit.
func lifecycleRepository(t *testing.T) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	temporaryRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mainRoot := filepath.Join(temporaryRoot, "repo")
	runCleanupGit(t, temporaryRoot, "init", "-q", mainRoot)
	runCleanupGit(t, mainRoot, "config", "user.name", "Lifecycle Test")
	runCleanupGit(t, mainRoot, "config", "user.email", "lifecycle@example.invalid")
	writeCleanupFile(t, mainRoot, "do-work/queue/REQ-41-lifecycle-demo.md", cleanupRequest("REQ-41", "pending", ""))
	writeCleanupFile(t, mainRoot, "do-work/worktree-links", "shared-deps\n\nmissing-env.vars\n")
	writeCleanupFile(t, mainRoot, ".gitignore", "shared-deps/\n")
	writeCleanupFile(t, mainRoot, "README.md", "fixture\n")
	commitCleanupFixture(t, mainRoot)
	writeCleanupFile(t, mainRoot, "shared-deps/package.txt", "dependency\n")
	return mainRoot
}

func runWorktreeCommand(t *testing.T, mainRoot string, arguments ...string) (resultmodel.CommandResult, int) {
	t.Helper()
	var output bytes.Buffer
	exitCode := commandruntime.NewRuntime(&output, Handlers()).Run(append([]string{"--repo-root", mainRoot, "--format", "json", "worktree"}, arguments...))
	var result resultmodel.CommandResult
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatalf("worktree %v output is not JSON: %v\n%s", arguments, err, output.String())
	}
	if result.Worktree == nil {
		t.Fatalf("worktree %v carried no typed worktree result: exit %d\n%s", arguments, exitCode, output.String())
	}
	return result, exitCode
}

func lifecycleWorktreePath(mainRoot, operativeName string) string {
	return filepath.Join(filepath.Dir(mainRoot), "repo-worktrees", operativeName)
}

func hasFindingCode(result resultmodel.CommandResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

func commitInWorktree(t *testing.T, worktreePath, relativePath string) {
	t.Helper()
	writeCleanupFile(t, worktreePath, relativePath, "builder work\n")
	runCleanupGit(t, worktreePath, "add", "--", relativePath)
	runCleanupGit(t, worktreePath, "commit", "-q", "-m", "[REQ-41] builder work")
}

// Pins the GREEN condition: new then cleanup leaves no worktree, branch or link behind, and the
// merge in between is a real two-parent [REQ-N] merge whose reported range matches git.
func TestWorktreeLifecycleNewMergeCleanupLeavesNothingBehind(t *testing.T) {
	mainRoot := lifecycleRepository(t)
	worktreePath := lifecycleWorktreePath(mainRoot, lifecycleOperativeName)

	created, exitCode := runWorktreeCommand(t, mainRoot, "new", "REQ-41")
	if exitCode != 0 || created.Worktree.OperativeName != lifecycleOperativeName || created.Worktree.WorktreePath != worktreePath {
		t.Fatalf("new: exit %d result %#v findings %#v", exitCode, created.Worktree, created.Findings)
	}
	if strings.Join(created.Worktree.LinksCreated, ",") != "shared-deps" || strings.Join(created.Worktree.LinksSkipped, ",") != "missing-env.vars" || !hasFindingCode(created, "WORKTREE-LINK-SKIPPED") {
		t.Fatalf("new links: created %v skipped %v findings %#v", created.Worktree.LinksCreated, created.Worktree.LinksSkipped, created.Findings)
	}
	linkPath := filepath.Join(worktreePath, "shared-deps")
	if target, err := os.Readlink(linkPath); err != nil || target != filepath.Join(mainRoot, "shared-deps") {
		t.Fatalf("link target %q, %v", target, err)
	}

	commitInWorktree(t, worktreePath, "feature.txt")
	status, exitCode := runWorktreeCommand(t, mainRoot, "status")
	if exitCode != 0 || len(status.Worktree.StatusRows) != 1 {
		t.Fatalf("status: exit %d rows %#v findings %#v", exitCode, status.Worktree.StatusRows, status.Findings)
	}
	if row := status.Worktree.StatusRows[0]; row.OperativeName != lifecycleOperativeName || row.Ahead != 1 || row.Behind != 0 || row.Dirty || row.LastCommitUnix == 0 {
		t.Fatalf("status row %#v", row)
	}

	merged, exitCode := runWorktreeCommand(t, mainRoot, "merge", "REQ-41")
	if exitCode != 0 || merged.Worktree.Pre == "" || merged.Worktree.MergeHash == "" {
		t.Fatalf("merge: exit %d result %#v findings %#v", exitCode, merged.Worktree, merged.Findings)
	}
	if subject := runCleanupGit(t, mainRoot, "log", "-1", "--format=%s"); subject != "[REQ-41] merge builder branch "+lifecycleOperativeName {
		t.Fatalf("merge subject %q", subject)
	}
	parents := strings.Fields(runCleanupGit(t, mainRoot, "log", "-1", "--format=%p", "--abbrev=40"))
	if len(parents) != 2 || runCleanupGit(t, mainRoot, "rev-parse", "--short", parents[0]) != merged.Worktree.Pre {
		t.Fatalf("merge parents %v, pre %s", parents, merged.Worktree.Pre)
	}
	if head := runCleanupGit(t, mainRoot, "rev-parse", "--short", "HEAD"); head != merged.Worktree.MergeHash {
		t.Fatalf("HEAD %s, merge_hash %s", head, merged.Worktree.MergeHash)
	}

	cleaned, exitCode := runWorktreeCommand(t, mainRoot, "cleanup", "REQ-41")
	if exitCode != 0 {
		t.Fatalf("cleanup: exit %d findings %#v", exitCode, cleaned.Findings)
	}
	if listing := runCleanupGit(t, mainRoot, "worktree", "list", "--porcelain"); strings.Contains(listing, lifecycleOperativeName) {
		t.Fatalf("worktree still registered:\n%s", listing)
	}
	if branches := runCleanupGit(t, mainRoot, "branch", "--list", lifecycleOperativeName); branches != "" {
		t.Fatalf("branch still present: %q", branches)
	}
	if _, err := os.Lstat(worktreePath); !os.IsNotExist(err) {
		t.Fatalf("worktree directory still present: %v", err)
	}
}

// Pins the queue guard: a builder commit under do-work/ must refuse before git merge runs, so
// nothing lands and no merge is left in progress.
func TestWorktreeMergeRefusesBuilderCommitUnderDoWork(t *testing.T) {
	mainRoot := lifecycleRepository(t)
	if _, exitCode := runWorktreeCommand(t, mainRoot, "new", "REQ-41"); exitCode != 0 {
		t.Fatalf("new exit %d", exitCode)
	}
	commitInWorktree(t, lifecycleWorktreePath(mainRoot, lifecycleOperativeName), "do-work/queue/REQ-99-builder-wrote-queue.md")
	headBefore := runCleanupGit(t, mainRoot, "rev-parse", "HEAD")

	refused, exitCode := runWorktreeCommand(t, mainRoot, "merge", "REQ-41")
	if exitCode == 0 || refused.Outcome != resultmodel.OutcomeRefused || !hasFindingCode(refused, "WORKTREE-QUEUE-GUARD") {
		t.Fatalf("merge: exit %d outcome %s findings %#v", exitCode, refused.Outcome, refused.Findings)
	}
	if !strings.Contains(strings.Join(refused.Findings[0].AffectedPaths, ","), "do-work/queue/REQ-99-builder-wrote-queue.md") {
		t.Fatalf("refusal does not name the queue path: %#v", refused.Findings)
	}
	if headAfter := runCleanupGit(t, mainRoot, "rev-parse", "HEAD"); headAfter != headBefore {
		t.Fatalf("HEAD moved from %s to %s", headBefore, headAfter)
	}
	if _, err := os.Stat(filepath.Join(mainRoot, ".git", "MERGE_HEAD")); !os.IsNotExist(err) {
		t.Fatalf("MERGE_HEAD present: %v", err)
	}
}

// Pins the empty hand-back: git merge says "Already up to date." and exits 0, so the command
// must detect the empty branch itself and refuse instead of fabricating a commit.
func TestWorktreeMergeReportsEmptyHandBackWithoutCommit(t *testing.T) {
	mainRoot := lifecycleRepository(t)
	if _, exitCode := runWorktreeCommand(t, mainRoot, "new", "REQ-41"); exitCode != 0 {
		t.Fatalf("new exit %d", exitCode)
	}
	headBefore := runCleanupGit(t, mainRoot, "rev-parse", "HEAD")

	empty, exitCode := runWorktreeCommand(t, mainRoot, "merge", "REQ-41")
	if exitCode == 0 || !empty.Worktree.EmptyHandBack || empty.Worktree.MergeHash != "" {
		t.Fatalf("merge: exit %d result %#v findings %#v", exitCode, empty.Worktree, empty.Findings)
	}
	if headAfter := runCleanupGit(t, mainRoot, "rev-parse", "HEAD"); headAfter != headBefore {
		t.Fatalf("HEAD moved from %s to %s", headBefore, headAfter)
	}
}

// Pins cleanup's preflight: builder dirt refuses before the owned link, the worktree or the
// branch is touched.
func TestWorktreeCleanupRefusesDirtyWorktreeBeforeAnySideEffect(t *testing.T) {
	mainRoot := lifecycleRepository(t)
	if _, exitCode := runWorktreeCommand(t, mainRoot, "new", "REQ-41"); exitCode != 0 {
		t.Fatalf("new exit %d", exitCode)
	}
	worktreePath := lifecycleWorktreePath(mainRoot, lifecycleOperativeName)
	writeCleanupFile(t, worktreePath, "uncommitted.txt", "builder dirt\n")

	refused, exitCode := runWorktreeCommand(t, mainRoot, "cleanup", "REQ-41")
	if exitCode == 0 || refused.Outcome != resultmodel.OutcomeRefused || !hasFindingCode(refused, "WORKTREE-DIRTY") {
		t.Fatalf("cleanup: exit %d outcome %s findings %#v", exitCode, refused.Outcome, refused.Findings)
	}
	if _, err := os.Readlink(filepath.Join(worktreePath, "shared-deps")); err != nil {
		t.Fatalf("owned link removed before refusal: %v", err)
	}
	if listing := runCleanupGit(t, mainRoot, "worktree", "list", "--porcelain"); !strings.Contains(listing, worktreePath) {
		t.Fatalf("worktree unregistered before refusal:\n%s", listing)
	}
	if branches := runCleanupGit(t, mainRoot, "branch", "--list", lifecycleOperativeName); branches == "" {
		t.Fatal("branch deleted before refusal")
	}
}

// Pins the collision rule: an existing derived name gets -2, and the leftover branch is neither
// deleted nor moved.
func TestWorktreeNewAppendsNumericSuffixOnCollision(t *testing.T) {
	mainRoot := lifecycleRepository(t)
	runCleanupGit(t, mainRoot, "branch", lifecycleOperativeName)
	leftoverTip := runCleanupGit(t, mainRoot, "rev-parse", lifecycleOperativeName)

	created, exitCode := runWorktreeCommand(t, mainRoot, "new", "REQ-41")
	variantName := lifecycleOperativeName + "-2"
	if exitCode != 0 || created.Worktree.OperativeName != variantName || created.Worktree.WorktreePath != lifecycleWorktreePath(mainRoot, variantName) {
		t.Fatalf("new: exit %d result %#v findings %#v", exitCode, created.Worktree, created.Findings)
	}
	if !hasFindingCode(created, "WORKTREE-NAME-COLLISION") {
		t.Fatalf("coexistence not reported: %#v", created.Findings)
	}
	if tip := runCleanupGit(t, mainRoot, "rev-parse", lifecycleOperativeName); tip != leftoverTip {
		t.Fatalf("leftover branch moved from %s to %s", leftoverTip, tip)
	}
}
