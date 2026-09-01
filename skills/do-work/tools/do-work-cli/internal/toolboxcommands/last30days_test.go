package toolboxcommands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestLast30DaysCompleteNeedsRuntime(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "SKILL.md"), []byte("skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	if last30DaysComplete(root) {
		t.Fatal("SKILL-only tree accepted")
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "last30days.py"), []byte("runtime"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !last30DaysComplete(root) {
		t.Fatal("complete tree rejected")
	}
}

func TestLast30DaysCheckReportsTargetPythonAbsence(t *testing.T) {
	repository := toolboxTestRepository(t)
	target := filepath.Join(repository, ".claude", "skills", "last30days")
	if err := os.MkdirAll(filepath.Join(target, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "SKILL.md"), []byte("skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "scripts", "last30days.py"), []byte("runtime"), 0o644); err != nil {
		t.Fatal(err)
	}
	originalLookup := last30DaysLookPath
	last30DaysLookPath = func(string) (string, error) { return "", errors.New("absent") }
	t.Cleanup(func() { last30DaysLookPath = originalLookup })
	result := handleLast30Days(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"check", repository})
	if result.Outcome != resultmodel.OutcomeFindings || result.ExactTextOutput == nil || !strings.Contains(*result.ExactTextOutput, "python 3.12+: FAILED") {
		t.Fatalf("result=%+v", result)
	}
}

func TestRemediationLast30DaysRefusesTrackedTargetBeforeCloneOrMutation(t *testing.T) {
	repository := toolboxTestRepository(t)
	target := filepath.Join(repository, ".claude", "skills", "last30days")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	skill := filepath.Join(target, "SKILL.md")
	if err := os.WriteFile(skill, []byte("tracked original"), 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", ".claude/skills/last30days/SKILL.md")
	toolboxTestGit(t, repository, "commit", "-m", "tracked target")
	if err := os.WriteFile(skill, []byte("dirty second writer"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := handleLast30Days(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"install", repository, filepath.Join(repository, "source-that-must-not-be-read")})
	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("tracked target accepted: %+v", result)
	}
	if contents, err := os.ReadFile(skill); err != nil || string(contents) != "dirty second writer" {
		t.Fatalf("tracked target changed before refusal: %q %v", contents, err)
	}
}

func TestPublishLast30DaysCopyFailureLeavesNoPublication(t *testing.T) {
	parent, source, target := last30DaysPublicationFixture(t, false)
	withLast30DaysPublicationHooks(t)
	last30DaysCopyTree = func(string, string) error { return errors.New("injected copy failure") }

	if err := publishLast30Days(context.Background(), source, target); err == nil || !strings.Contains(err.Error(), "clone/copy FAILED") {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLast30DaysAbsent(t, target)
	assertNoLast30DaysPrivatePaths(t, parent)
}

func TestPublishLast30DaysPartialPublicationFailureRemovesOwnedTarget(t *testing.T) {
	parent, source, target := last30DaysPublicationFixture(t, false)
	withLast30DaysPublicationHooks(t)
	copyCalls := 0
	last30DaysCopyTree = func(from, to string) error {
		copyCalls++
		if copyCalls == 2 {
			if err := os.WriteFile(filepath.Join(to, "partial"), []byte("partial"), 0o600); err != nil {
				return err
			}
			return errors.New("injected publication failure")
		}
		return copyLast30DaysTree(from, to)
	}

	if err := publishLast30Days(context.Background(), source, target); err == nil || !strings.Contains(err.Error(), "publication copy FAILED") {
		t.Fatalf("unexpected error: %v", err)
	}
	assertLast30DaysAbsent(t, target)
	assertNoLast30DaysPrivatePaths(t, parent)
}

func TestPublishLast30DaysFailureRestoresPreviousTree(t *testing.T) {
	parent, source, target := last30DaysPublicationFixture(t, true)
	withLast30DaysPublicationHooks(t)
	copyCalls := 0
	last30DaysCopyTree = func(from, to string) error {
		copyCalls++
		if copyCalls == 2 {
			return errors.New("injected publication failure")
		}
		return copyLast30DaysTree(from, to)
	}

	if err := publishLast30Days(context.Background(), source, target); err == nil {
		t.Fatal("publication failure was accepted")
	}
	assertLast30DaysPreviousTree(t, target)
	assertNoLast30DaysPrivatePaths(t, parent)
}

func TestPublishLast30DaysCancellationAfterBackupRestoresPreviousTree(t *testing.T) {
	parent, source, target := last30DaysPublicationFixture(t, true)
	withLast30DaysPublicationHooks(t)
	ctx, cancel := context.WithCancel(context.Background())
	last30DaysRename = func(from, to string) error {
		err := os.Rename(from, to)
		if err == nil && filepath.Base(to) == "previous" && strings.Contains(filepath.Dir(to), ".last30days.backup.") {
			cancel()
		}
		return err
	}

	if err := publishLast30Days(ctx, source, target); !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected cancellation result: %v", err)
	}
	assertLast30DaysPreviousTree(t, target)
	assertNoLast30DaysPrivatePaths(t, parent)
}

func TestPublishLast30DaysCollisionPreservesBothWriters(t *testing.T) {
	parent, source, target := last30DaysPublicationFixture(t, true)
	withLast30DaysPublicationHooks(t)
	last30DaysRename = func(from, to string) error {
		err := os.Rename(from, to)
		if err == nil && filepath.Base(to) == "previous" && strings.Contains(filepath.Dir(to), ".last30days.backup.") {
			if mkdirErr := os.Mkdir(target, 0o700); mkdirErr != nil {
				return mkdirErr
			}
			return os.WriteFile(filepath.Join(target, "keep.txt"), []byte("second writer"), 0o600)
		}
		return err
	}

	err := publishLast30Days(context.Background(), source, target)
	if err == nil || !strings.Contains(err.Error(), "reappeared") || !strings.Contains(err.Error(), "prior tree remains") {
		t.Fatalf("unexpected collision result: %v", err)
	}
	entries, readErr := os.ReadDir(target)
	if readErr != nil || len(entries) != 1 || entries[0].Name() != "keep.txt" {
		t.Fatalf("second writer changed: entries=%v err=%v", entries, readErr)
	}
	contents, readErr := os.ReadFile(filepath.Join(target, "keep.txt"))
	if readErr != nil || string(contents) != "second writer" {
		t.Fatalf("second writer bytes changed: %q %v", contents, readErr)
	}
	backups, globErr := filepath.Glob(filepath.Join(parent, ".last30days.backup.*", "previous", "SKILL.md"))
	if globErr != nil || len(backups) != 1 {
		t.Fatalf("prior tree is not recoverable: %v %v", backups, globErr)
	}
	prior, readErr := os.ReadFile(backups[0])
	if readErr != nil || string(prior) != "previous skill" {
		t.Fatalf("prior bytes changed: %q %v", prior, readErr)
	}
	stages, _ := filepath.Glob(filepath.Join(parent, ".last30days.staging.*"))
	if len(stages) != 0 {
		t.Fatalf("staging paths leaked: %v", stages)
	}
}

func last30DaysPublicationFixture(t *testing.T, previous bool) (string, string, string) {
	t.Helper()
	parent := filepath.Join(t.TempDir(), "skills")
	source := filepath.Join(t.TempDir(), "source")
	target := filepath.Join(parent, "last30days")
	for _, directory := range []string{parent, filepath.Join(source, "scripts")} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for path, contents := range map[string]string{
		filepath.Join(source, "SKILL.md"):                 "new skill",
		filepath.Join(source, "scripts", "last30days.py"): "new runtime",
	} {
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if previous {
		if err := os.MkdirAll(filepath.Join(target, "legacy"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(target, "SKILL.md"), []byte("previous skill"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(target, "legacy", "data.bin"), []byte("previous data"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return parent, source, target
}

func withLast30DaysPublicationHooks(t *testing.T) {
	t.Helper()
	originalCopy, originalRename, originalMkdir := last30DaysCopyTree, last30DaysRename, last30DaysMkdir
	t.Cleanup(func() {
		last30DaysCopyTree, last30DaysRename, last30DaysMkdir = originalCopy, originalRename, originalMkdir
	})
}

func assertLast30DaysAbsent(t *testing.T, target string) {
	t.Helper()
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatalf("target exists after failure: %v", err)
	}
}

func assertNoLast30DaysPrivatePaths(t *testing.T, parent string) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(parent, ".last30days.*"))
	if err != nil || len(paths) != 0 {
		t.Fatalf("private paths leaked: %v %v", paths, err)
	}
}

func assertLast30DaysPreviousTree(t *testing.T, target string) {
	t.Helper()
	for path, expected := range map[string]string{
		filepath.Join(target, "SKILL.md"):           "previous skill",
		filepath.Join(target, "legacy", "data.bin"): "previous data",
	} {
		contents, err := os.ReadFile(path)
		if err != nil || string(contents) != expected {
			t.Fatalf("previous tree changed at %s: %q %v", path, contents, err)
		}
	}
}
