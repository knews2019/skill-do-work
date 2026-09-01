package knowledgecommands

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestBKBInitCreatesExactScaffoldAndNeverOverwrites(t *testing.T) {
	root := t.TempDir()
	result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--dry-run"})
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Changes) != len(scaffoldFiles) {
		t.Fatalf("dry-run result = %s, %d changes", result.Outcome, len(result.Changes))
	}
	if _, err := os.Stat(filepath.Join(root, "kb")); !os.IsNotExist(err) {
		t.Fatal("dry-run created the target")
	}

	result = initializeWithoutRepository(root, bkbOptions{target: "kb"}, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("init outcome = %s: %+v", result.Outcome, result.Findings)
	}
	for _, file := range scaffoldFiles {
		data, err := os.ReadFile(filepath.Join(root, "kb", filepath.FromSlash(file.path)))
		if err != nil {
			t.Fatalf("read %s: %v", file.path, err)
		}
		if string(data) != renderScaffold(file.content, "2026-09-01") {
			t.Fatalf("%s bytes differ from canonical scaffold", file.path)
		}
	}

	sentinel := filepath.Join(root, "kb", "wiki", "overview.md")
	if err := os.WriteFile(sentinel, []byte("preserve me\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result = initializeWithoutRepository(root, bkbOptions{target: "kb", fillGaps: true}, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("fill-gaps outcome = %s", result.Outcome)
	}
	data, _ := os.ReadFile(sentinel)
	if string(data) != "preserve me\n" {
		t.Fatal("fill-gaps overwrote an existing file")
	}
}

func TestBKBInitCommitStagesOnlyExactScaffoldFiles(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.invalid")
	runGitFixture(t, root, "config", "user.name", "Fixture")
	if err := os.WriteFile(filepath.Join(root, "unrelated.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitFixture(t, root, "add", "unrelated.txt")
	runGitFixture(t, root, "commit", "-m", "baseline")
	if err := os.WriteFile(filepath.Join(root, "unrelated.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--commit"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("commit outcome = %s: %+v", result.Outcome, result.Findings)
	}
	output := runGitFixture(t, root, "show", "--name-only", "--pretty=format:", "HEAD")
	for _, line := range strings.Fields(output) {
		if !strings.HasPrefix(line, "kb/") {
			t.Errorf("commit included unrelated path %q", line)
		}
	}
	if status := runGitFixture(t, root, "status", "--short", "--", "unrelated.txt"); !strings.Contains(status, "unrelated.txt") {
		t.Fatal("unrelated dirty path was changed or staged")
	}
}

func TestStandaloneBKBInitCommitAndFailedCommitRollback(t *testing.T) {
	t.Run("commit exact scaffold", func(t *testing.T) {
		root := t.TempDir()
		t.Setenv("GIT_AUTHOR_NAME", "Fixture")
		t.Setenv("GIT_AUTHOR_EMAIL", "fixture@example.invalid")
		t.Setenv("GIT_COMMITTER_NAME", "Fixture")
		t.Setenv("GIT_COMMITTER_EMAIL", "fixture@example.invalid")
		result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--commit"})
		if result.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("outcome = %s: %+v", result.Outcome, result.Findings)
		}
		paths := strings.Fields(runGitFixture(t, filepath.Join(root, "kb"), "show", "--name-only", "--pretty=format:", "HEAD"))
		if len(paths) != len(scaffoldFiles) {
			t.Fatalf("committed %d paths, want %d", len(paths), len(scaffoldFiles))
		}
	})

	t.Run("failed commit removes only invocation paths", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "kb"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "kb", "preserve.txt"), []byte("keep\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "missing-gitconfig"))
		t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
		t.Setenv("GIT_CONFIG_COUNT", "1")
		t.Setenv("GIT_CONFIG_KEY_0", "user.useConfigOnly")
		t.Setenv("GIT_CONFIG_VALUE_0", "true")
		for _, name := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL", "EMAIL"} {
			old, existed := os.LookupEnv(name)
			_ = os.Unsetenv(name)
			t.Cleanup(func() {
				if existed {
					_ = os.Setenv(name, old)
				} else {
					_ = os.Unsetenv(name)
				}
			})
		}
		result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--commit"})
		if result.Outcome != resultmodel.OutcomeRolledBack {
			t.Fatalf("outcome = %s: %+v", result.Outcome, result.Findings)
		}
		data, err := os.ReadFile(filepath.Join(root, "kb", "preserve.txt"))
		if err != nil || string(data) != "keep\n" {
			t.Fatal("pre-existing file was not preserved")
		}
		if _, err := os.Stat(filepath.Join(root, "kb", ".git")); !os.IsNotExist(err) {
			t.Fatal("invocation-created Git metadata survived rollback")
		}
	})
}

func runGitFixture(t *testing.T, root string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s: %v", arguments, output, err)
	}
	return string(output)
}

func TestScaffoldBytesMatchHumanReadableReference(t *testing.T) {
	referencePath := filepath.Join("..", "..", "..", "..", "..", "do-work-knowledge", "actions", "bkb-reference.md")
	referenceBytes, err := os.ReadFile(referencePath)
	if err != nil {
		t.Fatal(err)
	}
	reference := string(referenceBytes)
	for _, file := range scaffoldFiles {
		var expected string
		if file.path == "CLAUDE.md" {
			expected = fencedContentAfter(t, reference, "Reference form for the `<path>/CLAUDE.md` schema")
		} else {
			expected = fencedContentAfter(t, reference, "**`"+file.path+"`:**")
		}
		if renderScaffold(file.content, "2026-09-01") != strings.ReplaceAll(expected, "{today}", "2026-09-01") {
			t.Errorf("canonical %s bytes drifted from bkb-reference.md", file.path)
		}
	}
}

func fencedContentAfter(t *testing.T, source, marker string) string {
	t.Helper()
	markerIndex := strings.Index(source, marker)
	if markerIndex < 0 {
		t.Fatalf("reference marker %q not found", marker)
	}
	fenceIndex := strings.Index(source[markerIndex:], "```markdown\n")
	if fenceIndex < 0 {
		t.Fatalf("markdown fence after %q not found", marker)
	}
	contentStart := markerIndex + fenceIndex + len("```markdown\n")
	contentEnd := strings.Index(source[contentStart:], "\n```")
	if contentEnd < 0 {
		t.Fatalf("closing fence after %q not found", marker)
	}
	return source[contentStart:contentStart+contentEnd] + "\n"
}

func TestBKBInitRefusesSymlinkParent(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "kb")); err != nil {
		t.Fatal(err)
	}
	result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--fill-gaps"})
	if result.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("outcome = %s, want refused", result.Outcome)
	}
	entries, _ := os.ReadDir(outside)
	if len(entries) != 0 {
		t.Fatal("symlink target was mutated")
	}
}

func TestBKBInitRefusesNestedScaffoldSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "kb"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "kb", "raw")); err != nil {
		t.Fatal(err)
	}
	result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--fill-gaps"})
	if result.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("outcome = %s, want refused", result.Outcome)
	}
	entries, _ := os.ReadDir(outside)
	if len(entries) != 0 {
		t.Fatal("nested symlink target was mutated")
	}
}

func TestBKBInitDryRunHonorsGitDirtyTargetPreflight(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.invalid")
	runGitFixture(t, root, "config", "user.name", "Fixture")
	writeFixture(t, root, "kb/wiki/overview.md", "tracked\n")
	if err := os.MkdirAll(filepath.Join(root, "kb", "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGitFixture(t, root, "add", "kb/wiki/overview.md")
	runGitFixture(t, root, "commit", "-m", "baseline")
	if err := os.Remove(filepath.Join(root, "kb", "wiki", "overview.md")); err != nil {
		t.Fatal(err)
	}
	result := handleBKBInit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb", "--fill-gaps", "--dry-run"})
	if result.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("outcome = %s, want dirty-target refusal: %+v", result.Outcome, result.Findings)
	}
	if _, err := os.Stat(filepath.Join(root, "kb", "wiki", "_master_index.md")); !os.IsNotExist(err) {
		t.Fatal("refused dry-run wrote scaffold bytes")
	}
}
