package toolboxcommands

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	if result := architectureScan(context, "missing-reports"); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("absent first-run report directory was not treated as an empty history: %+v", result)
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

func TestArchitecturePublishStopsOnNonCollisionClaimFailure(t *testing.T) {
	repository := toolboxTestRepository(t)
	draft := filepath.Join(repository, "draft.html")
	if err := os.WriteFile(draft, []byte("<html>complete</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	original := architectureClaimBundle
	architectureClaimBundle = func(_ string, relative string, _ os.FileMode) (os.FileInfo, error) {
		return nil, &os.PathError{Op: "mkdir", Path: relative, Err: fs.ErrPermission}
	}
	t.Cleanup(func() { architectureClaimBundle = original })
	candidate := filepath.Join(repository, "reports", "2026-09-01_1200_architecture-report")
	resultChannel := make(chan resultmodel.CommandResult, 1)
	go func() {
		resultChannel <- handleArchitecture(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--publish", draft, candidate})
	}()
	select {
	case result := <-resultChannel:
		if result.Outcome != resultmodel.OutcomeFindings || len(result.Findings) != 1 || result.Findings[0].Code != "ARCHITECTURE-BUNDLE-CLAIM-FAILED" || !strings.Contains(strings.Join(result.Findings[0].AffectedPaths, " "), "2026-09-01_1200_architecture-report") || !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), "permission denied") {
			t.Fatalf("claim failure result = %#v", result)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("non-collision bundle claim failure did not terminate")
	}
}

func TestArchitecturePublishCommitCommitsPublishedBundle(t *testing.T) {
	repository := toolboxTestRepository(t)
	draft := filepath.Join(repository, "draft.html")
	if err := os.WriteFile(draft, []byte("<html>complete</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "draft.html")
	toolboxTestGit(t, repository, "commit", "-m", "draft")
	candidate := filepath.Join(repository, "reports", "2026-09-01_1200_architecture-report")
	result := handleArchitecture(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--commit", "--publish", draft, candidate})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("commit publication = %#v", result)
	}
	if got := strings.TrimSpace(toolboxTestGit(t, repository, "show", "--name-only", "--format=", "HEAD")); got != "reports/2026-09-01_1200_architecture-report/index.html" {
		t.Fatalf("committed paths = %q", got)
	}
}
