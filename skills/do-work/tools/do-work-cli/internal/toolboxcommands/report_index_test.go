package toolboxcommands

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// reportIndexFixture builds an ai-reports/ folder holding one bundle per naming style seen in
// consumer repositories, entry-file fallbacks, a superseded proposal pair and loose files.
func reportIndexFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"2026-09-01_1200_alpha-report/index.html":                    `<html><head><title>Alpha &amp; Report</title></head><body><h1>Ignored</h1></body></html>`,
		"2026-09-02_req-77-beta-notes/index.md":                      "intro\n# Beta notes for UR-5\n",
		"2026-09-03-gamma-notes/aaa.md":                              "# Not the entry\n",
		"2026-09-03-gamma-notes/prompt.md":                           "# Gamma notes\n",
		"REQ-0412-delta-fix/index.html":                              `<body><h2 class="lead">Delta <em>fix</em></h2></body>`,
		"epsilon-review-2026-09-04/README.md":                        "# Epsilon review\n",
		"2026-09-05_120501_zeta-run/appendix.md":                     "# Appendix\n",
		"2026-09-05_120501_zeta-run/summary.html":                    "<title>Zeta run</title>",
		"2026-09-06_1000_architecture-report/architecture-report.md": "# Architecture map\n",
		"2026-09-07_0900_old-proposal/index.html":                    `<head><meta content="proposal" name="ai-report-kind"><meta name="ai-report-verdict" content="rejected"><title>Old proposal</title></head>`,
		"2026-09-08_0900_new-proposal/index.html":                    `<head><meta name='ai-report-kind' content='proposal'><meta name="ai-report-supersedes" content="ai-reports/2026-09-07_0900_old-proposal/"><title>New proposal</title></head>`,
		"2026-09-09_empty-bundle/picture.png":                        "png",
		".hidden-bundle/index.html":                                  "<title>Hidden</title>",
		"notes.patch":                                                "diff",
		"notes.md":                                                   "# Loose notes\n",
	}
	for relative, content := range files {
		path := filepath.Join(root, "ai-reports", filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("2026-09-01_1200_alpha-report", filepath.Join(root, "ai-reports", "linked-bundle")); err != nil {
		t.Fatal(err)
	}
	return root
}

func runReportIndex(t *testing.T, root string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	handler := Handlers()[CommandAIReportIndex]
	if handler == nil {
		t.Fatalf("%s is not registered", CommandAIReportIndex)
	}
	return handler(commandruntime.ExecutionContext{RepositoryRoot: root}, arguments)
}

func readReportCatalog(t *testing.T, root string) reportCatalog {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "ai-reports", "catalog.json"))
	if err != nil {
		t.Fatal(err)
	}
	var catalog reportCatalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		t.Fatalf("catalog.json does not parse: %v\n%s", err, data)
	}
	return catalog
}

func snapshotReportTree(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			snapshot[path] = "link:" + target
			return err
		}
		data, err := os.ReadFile(path)
		snapshot[path] = string(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

// Pins: a closed list of naming styles would drop or duplicate bundles, a loose file would be
// listed as a bundle, and the verb must change no existing file (the bundles are immutable).
func TestAIReportIndexCatalogsEveryBundleNamingStyleOnce(t *testing.T) {
	root := reportIndexFixture(t)
	before := snapshotReportTree(t, root)
	result := runReportIndex(t, root)
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("index outcome=%s findings=%+v", result.Outcome, result.Findings)
	}
	after := snapshotReportTree(t, root)
	for path, content := range before {
		if after[path] != content {
			t.Errorf("existing file changed: %s", path)
		}
		delete(after, path)
	}
	written := []string{}
	for path := range after {
		written = append(written, filepath.ToSlash(strings.TrimPrefix(path, root+string(filepath.Separator))))
	}
	if len(written) != 2 || !reflect.DeepEqual(map[string]bool{written[0]: true, written[1]: true}, map[string]bool{"ai-reports/catalog.json": true, "ai-reports/index.html": true}) {
		t.Fatalf("new files=%v, want only catalog.json and index.html", written)
	}
	catalog := readReportCatalog(t, root)
	if catalog.Generator != reportIndexGenerator {
		t.Fatalf("generator=%q", catalog.Generator)
	}
	type summary struct{ Path, Entry, Title, Date, Kind, IDs string }
	want := []summary{
		{"ai-reports/2026-09-09_empty-bundle", "", "2026-09-09_empty-bundle", "2026-09-09", "unknown", ""},
		{"ai-reports/2026-09-08_0900_new-proposal", "ai-reports/2026-09-08_0900_new-proposal/index.html", "New proposal", "2026-09-08T09:00", "proposal", ""},
		{"ai-reports/2026-09-07_0900_old-proposal", "ai-reports/2026-09-07_0900_old-proposal/index.html", "Old proposal", "2026-09-07T09:00", "proposal", ""},
		{"ai-reports/2026-09-06_1000_architecture-report", "ai-reports/2026-09-06_1000_architecture-report/architecture-report.md", "Architecture map", "2026-09-06T10:00", "architecture-report", ""},
		{"ai-reports/2026-09-05_120501_zeta-run", "ai-reports/2026-09-05_120501_zeta-run/summary.html", "Zeta run", "2026-09-05T12:05:01", "unknown", ""},
		{"ai-reports/epsilon-review-2026-09-04", "ai-reports/epsilon-review-2026-09-04/README.md", "Epsilon review", "2026-09-04", "unknown", ""},
		{"ai-reports/2026-09-03-gamma-notes", "ai-reports/2026-09-03-gamma-notes/prompt.md", "Gamma notes", "2026-09-03", "unknown", ""},
		{"ai-reports/2026-09-02_req-77-beta-notes", "ai-reports/2026-09-02_req-77-beta-notes/index.md", "Beta notes for UR-5", "2026-09-02", "unknown", "REQ-77,UR-5"},
		{"ai-reports/2026-09-01_1200_alpha-report", "ai-reports/2026-09-01_1200_alpha-report/index.html", "Alpha & Report", "2026-09-01T12:00", "unknown", ""},
		{"ai-reports/REQ-0412-delta-fix", "ai-reports/REQ-0412-delta-fix/index.html", "Delta fix", "", "unknown", "REQ-0412"},
	}
	got := []summary{}
	for _, bundle := range catalog.Bundles {
		if bundle.LinkedIDs == nil {
			t.Errorf("%s linked_ids is null, want []", bundle.Path)
		}
		got = append(got, summary{bundle.Path, bundle.Entry, bundle.Title, bundle.Date, bundle.Kind, strings.Join(bundle.LinkedIDs, ",")})
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("catalog bundles:\n got %+v\nwant %+v", got, want)
	}
}

// Pins: superseded_by must be computed from the successor's supersedes meta, and the page must
// grey the old row and link it to the successor rather than hide it.
func TestAIReportIndexLinksSupersededProposalToSuccessor(t *testing.T) {
	root := reportIndexFixture(t)
	if result := runReportIndex(t, root); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("index outcome=%s findings=%+v", result.Outcome, result.Findings)
	}
	bundles := map[string]reportBundle{}
	for _, bundle := range readReportCatalog(t, root).Bundles {
		bundles[bundle.Path] = bundle
	}
	older, newer := bundles["ai-reports/2026-09-07_0900_old-proposal"], bundles["ai-reports/2026-09-08_0900_new-proposal"]
	if newer.Supersedes != older.Path || older.SupersededBy != newer.Path || newer.SupersededBy != "" || older.Verdict != "rejected" {
		t.Fatalf("older=%+v newer=%+v", older, newer)
	}
	page, err := os.ReadFile(filepath.Join(root, "ai-reports", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(page)
	if !strings.Contains(html, `<meta name="generator" content="do-work-cli ai-report-index">`) {
		t.Fatal("index.html lacks the generator marker")
	}
	start := strings.Index(html, `<tr class="superseded">`)
	if start < 0 {
		t.Fatalf("no greyed row in index.html:\n%s", html)
	}
	row := html[start : start+strings.Index(html[start:], "</tr>")]
	if !strings.Contains(row, `href="2026-09-07_0900_old-proposal/index.html"`) || !strings.Contains(row, `superseded by <a href="2026-09-08_0900_new-proposal/index.html">`) {
		t.Fatalf("greyed row does not link the old report to its successor: %s", row)
	}
	if strings.Count(html, `<tr class="superseded">`) != 1 {
		t.Fatal("only the superseded proposal should be greyed")
	}
}

// Pins: find must list superseded bundles too, newest first, and must not touch the catalog files.
func TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded(t *testing.T) {
	root := reportIndexFixture(t)
	result := runReportIndex(t, root, "--find", "proposal")
	if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
		t.Fatalf("find outcome=%s findings=%+v", result.Outcome, result.Findings)
	}
	want := "ai-reports/2026-09-08_0900_new-proposal\tNew proposal\n" +
		"ai-reports/2026-09-07_0900_old-proposal\tOld proposal (superseded by ai-reports/2026-09-08_0900_new-proposal)\n"
	if *result.ExactTextOutput != want {
		t.Fatalf("find output:\n%q\nwant\n%q", *result.ExactTextOutput, want)
	}
	if _, err := os.Lstat(filepath.Join(root, "ai-reports", "catalog.json")); !os.IsNotExist(err) {
		t.Fatal("find wrote catalog.json")
	}
	if result := runReportIndex(t, root, "--find", "req-412"); result.ExactTextOutput == nil || *result.ExactTextOutput != "ai-reports/REQ-0412-delta-fix\tDelta fix\n" {
		t.Fatalf("id find ignoring leading zeros=%+v", result)
	}
	if result := runReportIndex(t, root, "--find", "no", "such", "topic"); result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || *result.ExactTextOutput != "no report matches no such topic\n" {
		t.Fatalf("no-match find=%+v", result)
	}
}

// Pins: a hand-made catalog or index page must never be overwritten, while the verb's own output
// is regenerated freely and byte-identically.
func TestAIReportIndexRefusesHandMadeCatalog(t *testing.T) {
	root := reportIndexFixture(t)
	catalogPath := filepath.Join(root, "ai-reports", "catalog.json")
	pagePath := filepath.Join(root, "ai-reports", "index.html")
	handMade := []byte(`{"bundles": []}` + "\n")
	if err := os.WriteFile(catalogPath, handMade, 0o644); err != nil {
		t.Fatal(err)
	}
	result := runReportIndex(t, root)
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) != 1 || result.Findings[0].Code != "AI-REPORT-INDEX-HAND-MADE" || result.Findings[0].AffectedPaths[0] != "ai-reports/catalog.json" {
		t.Fatalf("hand-made catalog result=%+v", result)
	}
	if data, _ := os.ReadFile(catalogPath); !bytes.Equal(data, handMade) {
		t.Fatal("hand-made catalog.json was overwritten")
	}
	if _, err := os.Lstat(pagePath); !os.IsNotExist(err) {
		t.Fatal("index.html was written beside a refused catalog")
	}
	if err := os.Remove(catalogPath); err != nil {
		t.Fatal(err)
	}
	if result := runReportIndex(t, root); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("first run=%+v", result)
	}
	firstCatalog, _ := os.ReadFile(catalogPath)
	firstPage, _ := os.ReadFile(pagePath)
	if result := runReportIndex(t, root); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("rerun over own output=%+v", result)
	}
	secondCatalog, _ := os.ReadFile(catalogPath)
	secondPage, _ := os.ReadFile(pagePath)
	if !bytes.Equal(firstCatalog, secondCatalog) || !bytes.Equal(firstPage, secondPage) {
		t.Fatal("rerun on an unchanged tree is not byte-identical")
	}
}
