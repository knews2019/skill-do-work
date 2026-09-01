package toolboxcommands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestFirstFreePortfolioPathUsesNumericSuffix(t *testing.T) {
	root := t.TempDir()
	candidate := filepath.Join(root, "snapshot.md")
	if err := os.WriteFile(candidate, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	rel, _, err := firstFreePortfolioPath(root, candidate)
	if err != nil || rel != "snapshot-2.md" {
		t.Fatalf("%q %v", rel, err)
	}
}

func TestPortfolioCreatesDistinctCanonicalAndSnapshot(t *testing.T) {
	repository := toolboxTestRepository(t)
	source := filepath.Join(repository, "source.md")
	if err := os.WriteFile(source, []byte("summary\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := handlePortfolio(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--with-snapshot", source, "portfolio/latest.md", "portfolio/snapshot.md"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("result=%+v", result)
	}
	canonical := filepath.Join(repository, "portfolio", "latest.md")
	snapshot := filepath.Join(repository, "portfolio", "snapshot.md")
	canonicalInfo, canonicalErr := os.Stat(canonical)
	snapshotInfo, snapshotErr := os.Stat(snapshot)
	if canonicalErr != nil || snapshotErr != nil {
		t.Fatalf("canonical=%v snapshot=%v", canonicalErr, snapshotErr)
	}
	if os.SameFile(canonicalInfo, snapshotInfo) {
		t.Fatal("canonical and immutable snapshot share an inode")
	}
}

func TestRemediationPortfolioRetainsSnapshotWhenCanonicalIsDirectory(t *testing.T) {
	repository := toolboxTestRepository(t)
	source := filepath.Join(repository, "source.md")
	if err := os.WriteFile(source, []byte("summary\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	canonical := filepath.Join(repository, "portfolio", "latest.md")
	if err := os.MkdirAll(canonical, 0o755); err != nil {
		t.Fatal(err)
	}
	snapshot := filepath.Join(repository, "portfolio", "snapshot.md")
	result := handlePortfolio(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--with-snapshot", source, canonical, snapshot})
	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatal("unsafe canonical directory unexpectedly succeeded")
	}
	if contents, err := os.ReadFile(snapshot); err != nil || string(contents) != "summary\n" {
		t.Fatalf("snapshot not retained: %q %v result=%+v", contents, err, result)
	}
	output, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"PORTFOLIO-CANONICAL-UNSAFE", "PORTFOLIO-SNAPSHOT-RETAINED", snapshot} {
		if !strings.Contains(string(output), required) {
			t.Fatalf("portfolio refusal output omitted %q: %s", required, output)
		}
	}
}
