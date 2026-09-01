package toolboxcommands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestArchitectureNameParserOrdersNumericSuffix(t *testing.T) {
	if !architectureName.MatchString("2026-01-01_0900_architecture-report-10") || architectureName.MatchString("notadate_architecture-report") {
		t.Fatal("architecture name parser regressed")
	}
}

func TestRemediationArchitectureAbsoluteScanAndReadDirFailure(t *testing.T) {
	repository := toolboxTestRepository(t)
	if err := os.WriteFile(filepath.Join(repository, "baseline"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "baseline")
	toolboxTestGit(t, repository, "commit", "-m", "baseline")
	head := strings.TrimSpace(toolboxTestGit(t, repository, "rev-parse", "--short", "HEAD"))
	prior := filepath.Join(repository, "reports", "2026-08-31_1200_architecture-report")
	if err := os.MkdirAll(prior, 0o755); err != nil {
		t.Fatal(err)
	}
	html := `<meta name="architecture-report-verified-at" content="` + head + `">` + "\n"
	if err := os.WriteFile(filepath.Join(prior, "index.html"), []byte(html), 0o644); err != nil {
		t.Fatal(err)
	}
	context := commandruntime.ExecutionContext{RepositoryRoot: repository}
	for _, reports := range []string{"reports", filepath.Join(repository, "reports")} {
		result := architectureScan(context, reports)
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || !strings.Contains(*result.ExactTextOutput, "prior_hash="+head+"\n") {
			t.Fatalf("scan %q=%+v", reports, result)
		}
	}
	if result := architectureScan(context, "missing-reports"); result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("ReadDir failure collapsed to clean result: %+v", result)
	}
}

func TestRemediationArchitectureFailedClaimRemainsOccupied(t *testing.T) {
	repository := toolboxTestRepository(t)
	draft := filepath.Join(repository, "draft.html")
	if err := os.WriteFile(draft, []byte("complete"), 0o644); err != nil {
		t.Fatal(err)
	}
	original := architectureAfterClaim
	architectureAfterClaim = func(bundle string) {
		_ = os.WriteFile(filepath.Join(bundle, "index.html"), []byte("competitor"), 0o644)
	}
	t.Cleanup(func() { architectureAfterClaim = original })
	candidate := filepath.Join(repository, "reports", "2026-09-01_1200_architecture-report")
	result := handleArchitecture(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--publish", draft, candidate})
	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatal("competing publication unexpectedly succeeded")
	}
	if info, err := os.Stat(candidate); err != nil || !info.IsDir() {
		t.Fatalf("failed claimed bundle was reusable: %v", err)
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
