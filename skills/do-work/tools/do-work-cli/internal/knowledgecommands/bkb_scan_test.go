package knowledgecommands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestBKBStatusAndLintAreReadOnlyAndActionable(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "kb", "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "kb/wiki/_master_index.md", "# Master Index\n\nTotal articles: 1 | Topic clusters: 1\n- [[_index_core]]\n")
	writeFixture(t, root, "kb/wiki/topics/_index_core.md", "# Core\n- [[missing]]\n")
	writeFixture(t, root, "kb/wiki/concepts/alpha.md", "---\ntitle: Alpha\ntype: invalid\ntopic_cluster: core\nrelated:\n  - page: missing\n    rel: invalid\ncreated: 2026-01-01\nupdated: 2026-01-01\nconfidence: medium\n---\n# Alpha\n[[missing]]\n")
	before := treeDigest(t, filepath.Join(root, "kb"))
	status := handleBKBStatus(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb"})
	if status.Outcome == resultmodel.OutcomeFailure || len(status.Findings) == 0 {
		t.Fatalf("status = %+v", status)
	}
	lint := handleBKBLint(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb"})
	if lint.Outcome != resultmodel.OutcomeFindings || len(lint.Findings) == 0 {
		t.Fatalf("lint = %+v", lint)
	}
	for _, finding := range lint.Findings {
		if len(finding.NextArgv) == 0 || len(finding.VerificationArgv) == 0 || finding.NextJustRecipe == "" {
			t.Fatalf("finding %s is not actionable", finding.Code)
		}
	}
	foundDangling := false
	for _, finding := range lint.Findings {
		if finding.Code == "BKB-TOPIC-DANGLING-ENTRY" {
			foundDangling = true
			if len(finding.AffectedPaths) != 1 || !strings.HasSuffix(finding.AffectedPaths[0], "wiki/topics/_index_core.md") {
				t.Fatalf("dangling affected=%v", finding.AffectedPaths)
			}
		}
	}
	if !foundDangling {
		t.Fatal("missing actionable dangling topic-index finding")
	}
	after := treeDigest(t, filepath.Join(root, "kb"))
	if before != after {
		t.Fatal("read-only BKB scans changed bytes")
	}
}

func TestBKBScanPreservesAbsoluteAndParentRelativeTargets(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, "kb", "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, outside, "kb/wiki/_master_index.md", "# Master\n")
	absolute := filepath.Join(outside, "kb")
	for _, target := range []string{absolute, filepath.Join("..", filepath.Base(outside), "kb")} {
		result := handleBKBLint(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", target})
		if result.Outcome == resultmodel.OutcomeFailure {
			t.Fatalf("target %q failed: %+v", target, result.Findings)
		}
		for _, finding := range result.Findings {
			if len(finding.NextArgv) < 4 || finding.NextArgv[len(finding.NextArgv)-1] != filepath.Clean(target) {
				t.Fatalf("target %q lost in %v", target, finding.NextArgv)
			}
		}
	}
}

func writeFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
