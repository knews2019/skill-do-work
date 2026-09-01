package knowledgecommands

import (
	"os"
	"path/filepath"
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
	after := treeDigest(t, filepath.Join(root, "kb"))
	if before != after {
		t.Fatal("read-only BKB scans changed bytes")
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
