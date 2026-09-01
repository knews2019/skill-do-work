package knowledgecommands

import (
	"encoding/json"
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

func TestBKBStatusMasterIndexCountMatrixIsActionableAndReadOnly(t *testing.T) {
	tests := []struct {
		name   string
		master string
		codes  []string
	}{
		{"declared inventory mismatch", "Total articles: 17 | Topic clusters: 3\n", []string{"BKB-STATUS-ARTICLE-COUNT-DISK-MISMATCH", "BKB-STATUS-TOPIC-COUNT-DISK-MISMATCH"}},
		{"missing", "# Master\n", []string{"BKB-STATUS-ARTICLE-COUNT-MISSING", "BKB-STATUS-TOPIC-COUNT-MISSING"}},
		{"malformed", "Total articles: many | Topic clusters: ?\n", []string{"BKB-STATUS-ARTICLE-COUNT-MALFORMED", "BKB-STATUS-TOPIC-COUNT-MALFORMED"}},
		{"duplicate", "Total articles: 0\nTotal articles: 0\nTopic clusters: 0\nTopic clusters: 0\n", []string{"BKB-STATUS-ARTICLE-COUNT-DUPLICATE", "BKB-STATUS-TOPIC-COUNT-DUPLICATE"}},
		{"inconsistent", "Total articles: 1\nTotal articles: 2\nTopic clusters: 3\nTopic clusters: 4\n", []string{"BKB-STATUS-ARTICLE-COUNT-INCONSISTENT", "BKB-STATUS-TOPIC-COUNT-INCONSISTENT"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "kb/wiki/_master_index.md", test.master)
			writeFixture(t, root, "kb/wiki/log.md", "# Log\n")
			if err := os.MkdirAll(filepath.Join(root, "kb", "raw", "inbox"), 0o755); err != nil {
				t.Fatal(err)
			}
			before := treeDigest(t, filepath.Join(root, "kb"))
			result := handleBKBStatus(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb"})
			for _, code := range test.codes {
				if !resultHasFindingCode(result, code) {
					t.Fatalf("missing %s: %#v", code, result.Findings)
				}
			}
			for _, finding := range result.Findings {
				if finding.Severity != resultmodel.SeverityInfo && (len(finding.NextArgv) == 0 || len(finding.VerificationArgv) == 0) {
					t.Fatalf("finding %s is not actionable: %#v", finding.Code, finding)
				}
			}
			textOutput, textError := resultmodel.RenderResult(result, resultmodel.FormatText)
			jsonOutput, jsonError := resultmodel.RenderResult(result, resultmodel.FormatJSON)
			if textError != nil || jsonError != nil {
				t.Fatalf("render text=%v json=%v", textError, jsonError)
			}
			var decoded resultmodel.CommandResult
			if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
				t.Fatal(err)
			}
			for _, code := range test.codes {
				if !strings.Contains(string(textOutput), code) || !resultHasFindingCode(decoded, code) {
					t.Fatalf("text/json parity lost %s: text=%s json=%s", code, textOutput, jsonOutput)
				}
			}
			if after := treeDigest(t, filepath.Join(root, "kb")); after != before {
				t.Fatal("BKB status count matrix changed bytes")
			}
		})
	}
}

func TestBKBLintReportsMalformedQuotedScalarAndMissingReciprocal(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "kb", "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "kb/wiki/_master_index.md", "# Master\n\nTotal articles: 2 | Topic clusters: 1\n- [[_index_core]]\n")
	writeFixture(t, root, "kb/wiki/topics/_index_core.md", "# Core\n\nTotal articles: 2\n- [[alpha]]\n- [[beta]]\n")
	writeFixture(t, root, "kb/wiki/concepts/alpha.md", "---\ntitle: \"Alpha \"quoted\" title\"\ntype: concept\ntopic_cluster: core\nsources: [raw/processed/alpha.md]\nrelated:\n  - page: beta\n    rel: complements\ncreated: 2026-01-01\nupdated: 2026-01-01\nconfidence: medium\n---\n# Alpha\n[[beta]]\n")
	writeFixture(t, root, "kb/wiki/concepts/beta.md", "---\ntitle: Beta\ntype: concept\ntopic_cluster: core\nsources: [raw/processed/beta.md]\nrelated: []\ncreated: 2026-01-01\nupdated: 2026-01-01\nconfidence: medium\n---\n# Beta\n[[alpha]]\n")
	before := treeDigest(t, filepath.Join(root, "kb"))
	result := handleBKBLint(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb"})
	for _, code := range []string{"BKB-FRONTMATTER-MALFORMED", "BKB-RELATIONSHIP-RECIPROCITY"} {
		if !resultHasFindingCode(result, code) {
			t.Fatalf("missing %s: %#v", code, result.Findings)
		}
	}
	for _, finding := range result.Findings {
		switch finding.Code {
		case "BKB-FRONTMATTER-MALFORMED":
			if len(finding.AffectedPaths) != 1 || !strings.HasSuffix(finding.AffectedPaths[0], "wiki/concepts/alpha.md") || len(finding.Evidence) == 0 || !strings.Contains(finding.Evidence[0], "field title") {
				t.Fatalf("malformed finding is not actionable: %#v", finding)
			}
		case "BKB-RELATIONSHIP-RECIPROCITY":
			if len(finding.AffectedPaths) != 2 || !strings.Contains(strings.Join(finding.AffectedPaths, " "), "alpha.md") || !strings.Contains(strings.Join(finding.AffectedPaths, " "), "beta.md") || len(finding.NextArgv) == 0 || len(finding.VerificationArgv) == 0 {
				t.Fatalf("reciprocity finding is not actionable: %#v", finding)
			}
		}
	}
	if after := treeDigest(t, filepath.Join(root, "kb")); after != before {
		t.Fatal("BKB lint changed bytes while reporting malformed YAML or reciprocity")
	}
}

func TestBKBLintAcceptsQuotedScalarsAndReciprocalRelationships(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "kb", "raw"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "kb/wiki/_master_index.md", "# Master\n\nTotal articles: 2 | Topic clusters: 1\n- [[_index_core]]\n")
	writeFixture(t, root, "kb/wiki/topics/_index_core.md", "# Core\n\nTotal articles: 2\n- [[alpha]]\n- [[beta]]\n")
	writeFixture(t, root, "kb/wiki/concepts/alpha.md", "---\ntitle: \"Alpha \\\"quoted\\\" title\" # valid YAML comment\ntype: concept\ntopic_cluster: core\nsources: [raw/processed/alpha.md]\nrelated:\n  - page: beta\n    rel: complements\ncreated: 2026-01-01\nupdated: 2026-01-01\nconfidence: medium\n---\n# Alpha\n[[beta]]\n")
	writeFixture(t, root, "kb/wiki/concepts/beta.md", "---\ntitle: 'Beta''s page' # valid YAML comment\ntype: concept\ntopic_cluster: core\nsources: [raw/processed/beta.md]\nrelated:\n  - page: alpha\n    rel: complements\ncreated: 2026-01-01\nupdated: 2026-01-01\nconfidence: medium\n---\n# Beta\n[[alpha]]\n")
	result := handleBKBLint(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--kb", "kb"})
	for _, code := range []string{"BKB-FRONTMATTER-MALFORMED", "BKB-RELATIONSHIP-RECIPROCITY"} {
		if resultHasFindingCode(result, code) {
			t.Fatalf("unexpected %s: %#v", code, result.Findings)
		}
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
