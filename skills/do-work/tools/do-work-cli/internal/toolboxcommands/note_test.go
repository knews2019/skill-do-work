package toolboxcommands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestNormalizeNote(t *testing.T) {
	cases := map[string]string{" add investigate xyz ": "investigate xyz", ` "keep quotes inside" `: "keep quotes inside", "add ": ""}
	for input, want := range cases {
		if got := normalizeNote(input); got != want {
			t.Fatalf("normalizeNote(%q)=%q want %q", input, got, want)
		}
	}
}

func TestNoteDryRunAndAppendPreserveExistingBytes(t *testing.T) {
	repository := toolboxTestRepository(t)
	doWork := filepath.Join(repository, "do-work")
	if err := os.Mkdir(doWork, 0o755); err != nil {
		t.Fatal(err)
	}
	notes := filepath.Join(doWork, "notes.md")
	if err := os.WriteFile(notes, []byte("existing\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "do-work/notes.md")
	toolboxTestGit(t, repository, "commit", "-m", "baseline")
	originalNow := toolboxNow
	toolboxNow = func() time.Time { return time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { toolboxNow = originalNow })
	context := commandruntime.ExecutionContext{RepositoryRoot: repository}
	if result := handleNote(context, []string{"--dry-run", "add", "new item"}); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("dry-run=%+v", result)
	}
	if contents, _ := os.ReadFile(notes); string(contents) != "existing\n" {
		t.Fatalf("dry-run changed notes: %q", contents)
	}
	if result := handleNote(context, []string{"add", "new item"}); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("append=%+v", result)
	}
	if contents, _ := os.ReadFile(notes); string(contents) != "existing\n- [2026-09-01] new item\n" {
		t.Fatalf("append bytes=%q", contents)
	}
}

func TestNoteDirtyGuardAndExactCommit(t *testing.T) {
	repository := toolboxTestRepository(t)
	if err := os.Mkdir(filepath.Join(repository, "do-work"), 0o755); err != nil {
		t.Fatal(err)
	}
	notes := filepath.Join(repository, "do-work", "notes.md")
	if err := os.WriteFile(notes, []byte("baseline\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "do-work/notes.md")
	toolboxTestGit(t, repository, "commit", "-m", "baseline")
	context := commandruntime.ExecutionContext{RepositoryRoot: repository}
	if err := os.WriteFile(notes, []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if result := handleNote(context, []string{"blocked"}); result.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("dirty target result=%+v", result)
	}
	toolboxTestGit(t, repository, "restore", "do-work/notes.md")
	if err := os.WriteFile(filepath.Join(repository, "unrelated.txt"), []byte("leave me"), 0o644); err != nil {
		t.Fatal(err)
	}
	if result := handleNote(context, []string{"--commit", "committed"}); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("commit result=%+v", result)
	}
	paths := toolboxTestGit(t, repository, "show", "--pretty=format:", "--name-only", "HEAD")
	if strings.TrimSpace(paths) != "do-work/notes.md" {
		t.Fatalf("commit paths=%q", paths)
	}
	if contents, _ := os.ReadFile(filepath.Join(repository, "unrelated.txt")); string(contents) != "leave me" {
		t.Fatal("exact commit disturbed unrelated file")
	}
}
