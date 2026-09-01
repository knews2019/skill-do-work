package toolboxcommands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestArchitectureNameParserOrdersNumericSuffix(t *testing.T) {
	if !architectureName.MatchString("2026-01-01_0900_architecture-report-10") || architectureName.MatchString("notadate_architecture-report") {
		t.Fatal("architecture name parser regressed")
	}
}

func TestArchitecturePublishUsesFirstFreeNumericBundle(t *testing.T) {
	repository := toolboxTestRepository(t)
	draft := filepath.Join(repository, "draft.html")
	if err := os.WriteFile(draft, []byte("<html>complete</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	candidate := filepath.Join(repository, "reports", "2026-09-01_1200_architecture-report")
	if err := os.MkdirAll(candidate, 0o755); err != nil {
		t.Fatal(err)
	}
	result := handleArchitecture(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--publish", draft, candidate})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("result=%+v", result)
	}
	contents, err := os.ReadFile(candidate + "-2/index.html")
	if err != nil || string(contents) != "<html>complete</html>" {
		t.Fatalf("published=%q err=%v", contents, err)
	}
}
