package toolboxcommands

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func runReportOrganize(t *testing.T, repository string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	handler := Handlers()["ai-report-organize"]
	if handler == nil {
		t.Fatal("ai-report-organize is not registered")
	}
	return handler(commandruntime.ExecutionContext{RepositoryRoot: repository}, arguments)
}

func writeOrganizeFiles(t *testing.T, repository string, files map[string]string) {
	t.Helper()
	for relative, content := range files {
		path := filepath.Join(repository, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// snapshotOrganizeTree records every directory, file and symlink outside .git, so a dry run that
// creates an empty month folder also counts as a write. .git is left out because git status may
// refresh its index file.
func snapshotOrganizeTree(t *testing.T, repository string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	err := filepath.WalkDir(repository, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, _ := filepath.Rel(repository, path)
		switch {
		case entry.IsDir() && entry.Name() == ".git":
			return filepath.SkipDir
		case entry.IsDir():
			snapshot[relative] = "dir"
		case entry.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(path)
			snapshot[relative] = "link:" + target
			return err
		default:
			data, err := os.ReadFile(path)
			snapshot[relative] = string(data)
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func freezeOrganizeClock(t *testing.T) {
	t.Helper()
	originalNow := toolboxNow
	toolboxNow = func() time.Time { return time.Date(2026, 11, 2, 12, 0, 0, 0, time.Local) }
	t.Cleanup(func() { toolboxNow = originalNow })
}

// Pins the REQ's proof. Moving a finished month's reports breaks every link that names their old
// folder: the queue file's path mention, the current report's link to a moved one, and a moved
// report's link out of ai-reports/ must all still open, while links that keep working (siblings
// that move together, the report's own image) and lookalikes (-rev1, a binary file) keep their
// bytes. Dry-run must write nothing, committing exactly the printed commit paths must leave a
// clean tree (the list is complete), and a second run must change nothing.
func TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen(t *testing.T) {
	freezeOrganizeClock(t)
	repository := toolboxTestRepository(t)
	alphaPage := `<a href="../2026-09-02_1200_beta/index.html#r1">beta</a> <a href="../../do-work/queue/REQ-001-sample.md">req</a> <img src="shot.png">`
	binaryNote := "ai-reports/2026-09-02_1200_beta/index.html\x00binary"
	writeOrganizeFiles(t, repository, map[string]string{
		"ai-reports/2026-09-01_1200_alpha/index.html":   alphaPage,
		"ai-reports/2026-09-01_1200_alpha/shot.png":     "png",
		"ai-reports/2026-09-02_1200_beta/index.html":    `<h1 id="r1">Beta</h1>`,
		"ai-reports/2026-09-03_1200_gamma/index.html":   `<h1>Gamma</h1>`,
		"ai-reports/2026-11-01_0900_current/index.html": `<a href="../2026-09-03_1200_gamma/index.html">gamma</a>`,
		"do-work/queue/REQ-001-sample.md":               "Report: ai-reports/2026-09-02_1200_beta/index.html\nOther: ai-reports/2026-09-02_1200_beta-rev1/index.html\n",
		"docs/diagram.bin":                              binaryNote,
	})
	toolboxTestGit(t, repository, "add", "-A")
	toolboxTestGit(t, repository, "commit", "-m", "fixture")

	before := snapshotOrganizeTree(t, repository)
	dryRun := runReportOrganize(t, repository, "--dry-run")
	if dryRun.Outcome != resultmodel.OutcomeSuccess || dryRun.ExactTextOutput == nil {
		t.Fatalf("dry run: %+v", dryRun)
	}
	dryText := *dryRun.ExactTextOutput
	if !strings.HasPrefix(dryText, "dry run: nothing written\n") || strings.Count(dryText, "\nmoved ") != 3 ||
		strings.Count(dryText, "\nlinks ") != 3 || strings.Contains(dryText, "commit paths:") {
		t.Fatalf("dry run must list 3 moves and 3 link edits and no commit block:\n%s", dryText)
	}
	if after := snapshotOrganizeTree(t, repository); !reflect.DeepEqual(before, after) {
		t.Fatalf("dry run wrote files:\nbefore %v\nafter  %v", before, after)
	}

	result := runReportOrganize(t, repository)
	if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
		t.Fatalf("organize: %+v", result)
	}
	text := *result.ExactTextOutput
	if !strings.Contains(text, "moved ai-reports/2026-09-01_1200_alpha -> ai-reports/2026-09_3-reports/2026-09-01_1200_alpha\n") {
		t.Fatalf("move line missing:\n%s", text)
	}
	read := func(relative string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(repository, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	for _, gone := range []string{"ai-reports/2026-09-01_1200_alpha", "ai-reports/2026-09-02_1200_beta", "ai-reports/2026-09-03_1200_gamma"} {
		if _, err := os.Lstat(filepath.Join(repository, gone)); !os.IsNotExist(err) {
			t.Errorf("%s was not moved", gone)
		}
	}
	movedAlpha := "ai-reports/2026-09_3-reports/2026-09-01_1200_alpha/index.html"
	expectedPages := map[string]string{
		movedAlpha: `<a href="../2026-09-02_1200_beta/index.html#r1">beta</a> <a href="../../../do-work/queue/REQ-001-sample.md">req</a> <img src="shot.png">`,
		"ai-reports/2026-09_3-reports/2026-09-02_1200_beta/index.html":  `<h1 id="r1">Beta</h1>`,
		"ai-reports/2026-09_3-reports/2026-09-03_1200_gamma/index.html": `<h1>Gamma</h1>`,
		"ai-reports/2026-11-01_0900_current/index.html":                 `<a href="../2026-09_3-reports/2026-09-03_1200_gamma/index.html">gamma</a>`,
		"do-work/queue/REQ-001-sample.md":                               "Report: ai-reports/2026-09_3-reports/2026-09-02_1200_beta/index.html\nOther: ai-reports/2026-09-02_1200_beta-rev1/index.html\n",
		"docs/diagram.bin":                                              binaryNote,
	}
	for relative, expected := range expectedPages {
		if actual := read(relative); actual != expected {
			t.Errorf("%s\n got %q\nwant %q", relative, actual, expected)
		}
	}
	// Every link must open from the file's new place, not just look right.
	for _, link := range []struct{ from, target string }{
		{movedAlpha, "../2026-09-02_1200_beta/index.html"},
		{movedAlpha, "../../../do-work/queue/REQ-001-sample.md"},
		{movedAlpha, "shot.png"},
		{"ai-reports/2026-11-01_0900_current/index.html", "../2026-09_3-reports/2026-09-03_1200_gamma/index.html"},
		{"do-work/queue/REQ-001-sample.md", "../../ai-reports/2026-09_3-reports/2026-09-02_1200_beta/index.html"},
	} {
		if _, err := os.Stat(filepath.Join(repository, filepath.Dir(filepath.FromSlash(link.from)), filepath.FromSlash(link.target))); err != nil {
			t.Errorf("link %s from %s does not open: %v", link.target, link.from, err)
		}
	}
	if catalog := readReportCatalog(t, repository); len(catalog.Bundles) != 4 {
		t.Fatalf("catalog lists %d bundles, want 4: %+v", len(catalog.Bundles), catalog.Bundles)
	}

	_, commitBlock, found := strings.Cut(text, "commit paths:\n")
	if !found {
		t.Fatalf("commit paths block missing:\n%s", text)
	}
	commitPaths := strings.Fields(commitBlock)
	toolboxTestGit(t, repository, append([]string{"add", "-A", "--"}, commitPaths...)...)
	toolboxTestGit(t, repository, append([]string{"commit", "-m", "Organize ai-reports by month", "--"}, commitPaths...)...)
	if status := toolboxTestGit(t, repository, "status", "--porcelain", "--untracked-files=all"); status != "" {
		t.Fatalf("commit paths left changes behind:\n%s\ncommit paths: %v", status, commitPaths)
	}

	organized := snapshotOrganizeTree(t, repository)
	again := runReportOrganize(t, repository)
	if again.Outcome != resultmodel.OutcomeSuccess || again.ExactTextOutput == nil || *again.ExactTextOutput != "nothing to organize\n" {
		t.Fatalf("second run: %+v", again)
	}
	if after := snapshotOrganizeTree(t, repository); !reflect.DeepEqual(organized, after) {
		t.Fatal("second run changed the tree")
	}
}

// Pins the recount and the refusals. A stale count (2026-08_1-report holding two reports after
// the move), a plain 2026-07 folder and the singular name must be fixed by renaming the existing
// folder, and a folder-level mention must follow the rename. A current-month folder must be
// emptied back to the top level and removed. A counted name that exists as a file must skip only
// June: POSIX rename would replace an empty directory silently, so no step may target an existing
// path. Loose files, undated folders and symlinks stay and are listed. An uncommitted file whose
// links would change must refuse the whole run and change nothing, because the action commits by
// exact path and would otherwise sweep that unrelated edit into the organize commit. A hand-made
// catalog page must refuse before any move: the catalog is written last, and refusing there left
// the moves uncommitted and a rerun printed nothing to commit. A current-month and a future-month
// report with the same name must not both plan the same top-level path (the second move stopped
// the run part way). A catalog page Git ignores must stay out of the commit paths, where
// `git add -A` rejects the whole list.
func TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets(t *testing.T) {
	freezeOrganizeClock(t)
	repository := toolboxTestRepository(t)
	writeOrganizeFiles(t, repository, map[string]string{
		"ai-reports/2026-08_1-report/a/index.html": "<h1>A</h1>",
		"ai-reports/2026-08-20_1000_b/index.html":  "<h1>B</h1>",
		"ai-reports/2026-07/c/index.html":          "<h1>C</h1>",
		"ai-reports/2026-11/d/index.html":          "<h1>D</h1>",
		"ai-reports/2026-12/d/index.html":          "<h1>D next month</h1>",
		".gitignore":                               "ai-reports/index.html\n",
		"ai-reports/2026-06-01_1000_e/index.html":  "<h1>E</h1>",
		"ai-reports/2026-06_1-report":              "a file, not a month folder\n",
		"ai-reports/notes.md":                      "# Loose notes\n",
		"ai-reports/undated-folder/index.html":     "<h1>Undated</h1>",
	})
	if err := os.Symlink("undated-folder", filepath.Join(repository, "ai-reports", "linked")); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "-A")
	toolboxTestGit(t, repository, "commit", "-m", "fixture")
	notes := "See ai-reports/2026-08_1-report/a/index.html and the folder ai-reports/2026-07/.\n"
	writeOrganizeFiles(t, repository, map[string]string{"docs/notes.md": notes})

	before := snapshotOrganizeTree(t, repository)
	refused := runReportOrganize(t, repository)
	if refused.Outcome != resultmodel.OutcomeRefused || len(refused.Findings) != 1 || refused.Findings[0].Code != "AI-REPORT-ORGANIZE-DIRTY" ||
		!reflect.DeepEqual(refused.Findings[0].AffectedPaths, []string{"docs/notes.md"}) {
		t.Fatalf("an uncommitted file with links to fix must refuse and name it: %+v", refused)
	}
	if after := snapshotOrganizeTree(t, repository); !reflect.DeepEqual(before, after) {
		t.Fatal("the refused run changed the tree")
	}

	toolboxTestGit(t, repository, "add", "docs/notes.md")
	toolboxTestGit(t, repository, "commit", "-m", "notes")
	writeOrganizeFiles(t, repository, map[string]string{"ai-reports/index.html": "<h1>My own index</h1>\n"})
	before = snapshotOrganizeTree(t, repository)
	if handMade := runReportOrganize(t, repository); handMade.Outcome != resultmodel.OutcomeRefused || len(handMade.Findings) != 1 ||
		handMade.Findings[0].Code != "AI-REPORT-INDEX-HAND-MADE" {
		t.Fatalf("a hand-made catalog page must refuse before any write: %+v", handMade)
	}
	if after := snapshotOrganizeTree(t, repository); !reflect.DeepEqual(before, after) {
		t.Fatal("the hand-made catalog refusal changed the tree")
	}
	if err := os.Remove(filepath.Join(repository, "ai-reports", "index.html")); err != nil {
		t.Fatal(err)
	}
	result := runReportOrganize(t, repository)
	if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
		t.Fatalf("organize: %+v", result)
	}
	text := *result.ExactTextOutput
	for _, line := range []string{
		"renamed ai-reports/2026-08_1-report -> ai-reports/2026-08_2-reports\n",
		"renamed ai-reports/2026-07 -> ai-reports/2026-07_1-report\n",
		"removed ai-reports/2026-11 (empty month folder)\n",
		"links docs/notes.md (2)\n",
		"left ai-reports/2026-06_1-report (",
		"left ai-reports/linked (symbolic link)\n",
		"left ai-reports/notes.md (loose file)\n",
		"left ai-reports/undated-folder (no date prefix)\n",
		"left ai-reports/d (another month's report moves here; month 2026-12 not organized)\n",
	} {
		if !strings.Contains(text, line) {
			t.Errorf("output lacks %q:\n%s", line, text)
		}
	}
	for _, present := range []string{
		"ai-reports/2026-08_2-reports/a/index.html",
		"ai-reports/2026-08_2-reports/2026-08-20_1000_b/index.html", "ai-reports/2026-07_1-report/c/index.html",
		"ai-reports/d/index.html", "ai-reports/2026-06-01_1000_e/index.html", "ai-reports/undated-folder/index.html",
		"ai-reports/2026-12/d/index.html",
	} {
		if _, err := os.Lstat(filepath.Join(repository, filepath.FromSlash(present))); err != nil {
			t.Errorf("%s missing after organize: %v", present, err)
		}
	}
	for _, gone := range []string{"ai-reports/2026-08_1-report", "ai-reports/2026-08-20_1000_b", "ai-reports/2026-07", "ai-reports/2026-11"} {
		if _, err := os.Lstat(filepath.Join(repository, filepath.FromSlash(gone))); !os.IsNotExist(err) {
			t.Errorf("%s still exists after organize", gone)
		}
	}
	after := snapshotOrganizeTree(t, repository)
	for _, untouched := range []string{"ai-reports/2026-06_1-report", "ai-reports/notes.md", "ai-reports/linked"} {
		if after[untouched] != before[untouched] {
			t.Errorf("%s changed: %q -> %q", untouched, before[untouched], after[untouched])
		}
	}
	if got, want := after["docs/notes.md"], "See ai-reports/2026-08_2-reports/a/index.html and the folder ai-reports/2026-07_1-report/.\n"; got != want {
		t.Errorf("docs/notes.md\n got %q\nwant %q", got, want)
	}
	if _, commitBlock, _ := strings.Cut(text, "commit paths:\n"); !strings.Contains("\n"+commitBlock, "\nai-reports/catalog.json\n") ||
		strings.Contains("\n"+commitBlock, "\nai-reports/index.html\n") {
		t.Errorf("commit paths must list the catalog and leave out the ignored catalog page:\n%s", commitBlock)
	}
}

// Pins REQ-708: `ai-report index` and `--check`, which every new report runs, recommend organize
// while a finished month's report is still loose, and only then. A hint that printed on an
// organized tree would teach users to ignore it, and `--check` must print the catalog's exact
// line and write nothing, because it runs after every report.
func TestReportOrganizeHintPrintsOnlyWhenWorkRemains(t *testing.T) {
	freezeOrganizeClock(t)
	indexText := func(repository string) string {
		t.Helper()
		result := runReportIndex(t, repository)
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
			t.Fatalf("ai-report-index: %+v", result)
		}
		return *result.ExactTextOutput
	}
	runOrganizeCheck := func(repository string) string {
		t.Helper()
		result := runReportOrganize(t, repository, "--check")
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
			t.Errorf("ai-report-organize --check must succeed with text output: %+v", result)
			return ""
		}
		return *result.ExactTextOutput
	}

	loose := t.TempDir()
	writeOrganizeFiles(t, loose, map[string]string{"ai-reports/2026-09-27_2200_a/index.html": "<h1>A</h1>"})
	before := snapshotOrganizeTree(t, loose)
	hintLine := runOrganizeCheck(loose)
	if strings.Count(hintLine, "\n") != 1 || !strings.HasSuffix(hintLine, "\n") || !strings.Contains(hintLine, "do-work-toolbox ai-report organize") {
		t.Errorf("--check on a loose finished month must print one line naming do-work-toolbox ai-report organize, got %q", hintLine)
	}
	if after := snapshotOrganizeTree(t, loose); !reflect.DeepEqual(before, after) {
		t.Errorf("--check wrote files:\nbefore %v\nafter  %v", before, after)
	}
	if catalogText := indexText(loose); !strings.Contains(catalogText, "do-work-toolbox ai-report organize") || !strings.HasSuffix(catalogText, hintLine) {
		t.Errorf("ai-report-index output must end with the --check line %q:\n%s", hintLine, catalogText)
	}

	organized := t.TempDir()
	writeOrganizeFiles(t, organized, map[string]string{
		"ai-reports/2026-11-01_0900_current/index.html":            "<h1>Current</h1>",
		"ai-reports/2026-09_1-report/2026-09-27_2200_a/index.html": "<h1>A</h1>",
	})
	if checkText := runOrganizeCheck(organized); checkText != "" {
		t.Errorf("--check on an organized tree must print nothing, got %q", checkText)
	}
	if catalogText := indexText(organized); strings.Contains(catalogText, "ai-report organize") {
		t.Errorf("ai-report-index recommended organize on an organized tree:\n%s", catalogText)
	}
}
