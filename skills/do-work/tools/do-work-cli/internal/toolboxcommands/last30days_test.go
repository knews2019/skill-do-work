package toolboxcommands

import (
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
