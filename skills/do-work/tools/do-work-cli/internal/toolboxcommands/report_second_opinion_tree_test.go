package toolboxcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// secondOpinionValidChangeList is an editor page whose "Changed in this iteration" list resolves to a badged section.
const secondOpinionValidChangeList = `<!doctype html><html><body><p>Original report: <a href="index.html">index.html</a></p><nav><h2>Changed in this iteration</h2><ul><li><a href="#rev-1">Clearer opening</a></li></ul></nav><section id="rev-1"><span>changed</span><h1>Clear report</h1></section></body></html>`

// secondOpinionFakeTreeEditor installs a fake codex that writes the bytes held in pageSource as its page and logs its
// prompt to promptLog; a test rewrites pageSource to change what the next run writes.
func secondOpinionFakeTreeEditor(t *testing.T, page string) (pageSource, promptLog string) {
	t.Helper()
	scratch := t.TempDir()
	pageSource, promptLog = filepath.Join(scratch, "page.html"), filepath.Join(scratch, "prompt.txt")
	if err := os.WriteFile(pageSource, []byte(page), 0644); err != nil {
		t.Fatal(err)
	}
	secondOpinionFakeCodex(t, `/bin/cat '`+pageSource+`' > "$output"; printf '%s' "$prompt" > '`+promptLog+`'`)
	return pageSource, promptLog
}

func secondOpinionStubJudgePass(t *testing.T) {
	t.Helper()
	previous := secondOpinionJudge
	secondOpinionJudge = func(string, []byte) (string, []resultmodel.AIReportSecondOpinionFinding) {
		return "pass", []resultmodel.AIReportSecondOpinionFinding{}
	}
	t.Cleanup(func() { secondOpinionJudge = previous })
}

// Pins: a page target refused or read as a bundle, a bundle preset lost when the target is a page, directed names
// without a tree id, a reused or wrong child id, a child page without its parent link and tree metas, a prompt that
// never names the base or the change list, a plain run that gains tree marks, an old non-tree page or a missing page
// accepted as a base, and a session finalize that forgets the prepared base.
func TestSecondOpinionTreeTargetsNamesAndParentHeader(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionFixturePresets(t)
	own := filepath.Join(bundle, "second-opinion-presets")
	os.MkdirAll(own, 0755)
	os.WriteFile(filepath.Join(own, "visual.md"), []byte("Bundle visual: one diagram per section.\n"), 0644)
	_, promptLog := secondOpinionFakeTreeEditor(t, secondOpinionValidChangeList)
	relative := "ai-reports/2026-10-15_1200-example"
	read := func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s not written: %v", path, err)
		}
		return string(data)
	}
	for _, id := range []string{"v1c01", "v1c02"} {
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex", "--preset", "bold"))[0]
		if row["page"] != relative+"/index-"+id+"-codex-actual-model-1-bold.html" || row["reply"] != relative+"/"+id+"-codex-actual-model-1-bold-reply.md" || row["tree_id"] != id || row["base_page"] != relative+"/index.html" {
			t.Fatalf("run on index.html, want %s: %v", id, row)
		}
	}
	first := "index-v1c01-codex-actual-model-1-bold.html"
	child := "index-v1c01c01-codex-actual-model-1-visual.html"
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex", "--preset", "visual", filepath.Join(bundle, first)))[0]
	if row["page"] != relative+"/"+child || row["tree_id"] != "v1c01c01" || row["base_page"] != relative+"/"+first || row["preset_file"] != relative+"/second-opinion-presets/visual.md" {
		t.Fatalf("run on the v1c01 page: %v", row)
	}
	prompt := read(promptLog)
	for _, fragment := range []string{"Base: " + first + "\n", "Changed in this iteration", "Direction: Bundle visual: one diagram per section.\n"} {
		if !strings.Contains(prompt, fragment) {
			t.Fatalf("prompt lacks %q: %s", fragment, prompt)
		}
	}
	page := read(filepath.Join(bundle, child))
	for _, mark := range []string{`<a href="` + first + `">` + first + `</a>`, `<meta name="ai-report-second-opinion-id" content="v1c01c01">`, `<meta name="ai-report-second-opinion-base" content="` + first + `">`} {
		if !strings.Contains(page, mark) {
			t.Fatalf("child page lacks %q: %s", mark, page)
		}
	}
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))[0]
	prompt, page = read(promptLog), read(filepath.Join(bundle, "index-codex-actual-model-1.html"))
	if row["page"] != relative+"/index-codex-actual-model-1.html" || row["tree_id"] != nil || row["base_page"] != nil || !strings.Contains(prompt, "Base: index.html\nChange list: None for this pass.\n") || strings.Contains(page, "ai-report-second-opinion-") || strings.Contains(page, "second-opinion-header") {
		t.Fatalf("plain run changed: %v\n%s\n%s", row, prompt, page)
	}
	for _, target := range []string{"index-codex-actual-model-1.html", "index-v1c09-codex-actual-model-1.html"} {
		result := secondOpinionInvoke(t, repo, "--backend", "codex", filepath.Join(bundle, target))
		if want := "base page must be index.html or a tree page"; result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) == 0 || !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), want) {
			t.Fatalf("%s accepted as a base: %+v", target, result.Findings)
		}
	}
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--prepare", "--backend", "claude", "--model", "claude-fable-5-1", filepath.Join(bundle, first)))[0]
	stage, _ := row["staging_page"].(string)
	if prompt, _ := row["prompt"].(string); row["state"] != "prepared" || row["base_page"] != relative+"/"+first || !strings.Contains(prompt, "Base: "+first+"\n") {
		t.Fatalf("prepare on the v1c01 page: %v", row)
	}
	os.WriteFile(stage, []byte(secondOpinionValidChangeList), 0644)
	os.WriteFile(strings.TrimSuffix(stage, ".html")+"-reply.md", []byte("Clarity: simpler.\n"), 0644)
	row = secondOpinionRows(t, secondOpinionInvoke(t, repo, "--finalize", stage))[0]
	if row["state"] != "ran" || row["page"] != relative+"/index-v1c01c02-claude-claude-fable-5-1.html" || row["tree_id"] != "v1c01c02" || row["base_page"] != relative+"/"+first || row["preset"] != nil {
		t.Fatalf("finalize lost the prepared base: %v", row)
	}
}

// Pins: a change list whose link has no target (or only a data-id look-alike) passing, a tree page with no change
// list passing, a failed check that deletes the page or touches index.html, and a valid list that never gets linked.
func TestSecondOpinionChangeListTargetsAreChecked(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionFixturePresets(t)
	secondOpinionStubJudgePass(t)
	pageSource, _ := secondOpinionFakeTreeEditor(t, secondOpinionValidChangeList)
	original, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	cases := []struct {
		id, page, missingTarget string
	}{
		{"v1c01", `<!doctype html><html><body><p>Original report: <a href="index.html">index.html</a></p><ul><li><a href="#rev-1">One</a></li><li><a href="#rev-2">Two</a></li></ul><section id="rev-1">changed</section><section data-id="rev-2">changed</section></body></html>`, "#rev-2"},
		{"v1c02", `<!doctype html><html><body><p>Original report: <a href="index.html">index.html</a></p><h1>No list</h1></body></html>`, ""},
	}
	for _, test := range cases {
		os.WriteFile(pageSource, []byte(test.page), 0644)
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex", "--preset", "bold"))[0]
		findings := []string{}
		for _, finding := range row["judge_findings"].([]any) {
			fields := finding.(map[string]any)
			if fields["kind"] == "AI-REPORT-SECOND-OPINION-CHANGE-LIST" {
				findings = append(findings, fmt.Sprint(fields["target"]))
			}
		}
		want := "[<nil>]"
		if test.missingTarget != "" {
			want = "[" + test.missingTarget + "]"
		}
		if row["judge_verdict"] != "fail" || row["banner_state"] != "not-inserted" || fmt.Sprint(findings) != want {
			t.Fatalf("%s: change list not checked, want findings %s: %v", test.id, want, row)
		}
		for _, name := range []string{"index-" + test.id + "-codex-actual-model-1-bold.html", test.id + "-codex-actual-model-1-bold-reply.md"} {
			if _, err := os.Stat(filepath.Join(bundle, name)); err != nil {
				t.Fatalf("failed check removed %s: %v", name, err)
			}
		}
		if after, _ := os.ReadFile(filepath.Join(bundle, "index.html")); string(after) != string(original) {
			t.Fatalf("%s: failed check changed index.html: %s", test.id, after)
		}
	}
	os.WriteFile(pageSource, []byte(secondOpinionValidChangeList), 0644)
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex", "--preset", "bold"))[0]
	index, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	if row["judge_verdict"] != "pass" || row["banner_state"] != "inserted" || !strings.Contains(string(index), "<!-- second-opinion-tree -->") {
		t.Fatalf("valid change list not linked: %v\n%s", row, index)
	}
	// Judge-failed pages stay in the list, so no later child is orphaned.
	for _, id := range []string{"v1c01", "v1c02", "v1c03"} {
		if !strings.Contains(string(index), `href="index-`+id+`-codex-actual-model-1-bold.html"`) {
			t.Fatalf("tree list lacks %s: %s", id, index)
		}
	}
}

// Pins: a tree list out of pre-order or out of step with a byte sort of the bundle, a second tree block, a per-filename
// banner on a tree page, a plain run that drops the tree block, and a marker in a script string read as the block.
func TestSecondOpinionTreeListOnIndexFollowsTreeOrder(t *testing.T) {
	repo, bundle := secondOpinionFixture(t)
	secondOpinionFixturePresets(t)
	secondOpinionStubJudgePass(t)
	secondOpinionFakeTreeEditor(t, secondOpinionValidChangeList)
	script := `<script>const sample = "<!-- second-opinion-tree --><li>fake</li><!-- /second-opinion-tree -->";</script>`
	os.WriteFile(filepath.Join(bundle, "index.html"), []byte(`<!doctype html><html><head>`+script+`</head><body><h1>Original</h1></body></html>`), 0644)
	first := filepath.Join(bundle, "index-v1c01-codex-actual-model-1-bold.html")
	for _, args := range [][]string{{"--preset", "bold"}, {"--preset", "bold"}, {"--preset", "visual", first}} {
		row := secondOpinionRows(t, secondOpinionInvoke(t, repo, append([]string{"--backend", "codex"}, args...)...))[0]
		if row["judge_verdict"] != "pass" || row["banner_state"] != "inserted" {
			t.Fatalf("%v did not write the tree list: %v", args, row)
		}
	}
	treeBlock := func() string {
		data, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
		index := string(data)
		if !strings.Contains(index, script) || strings.Count(index, "<!-- second-opinion-tree -->") != 2 {
			t.Fatalf("want the script untouched and exactly one real tree block: %s", index)
		}
		body := strings.Index(index, "<body>")
		start, end := strings.Index(index[body:], "<!-- second-opinion-tree -->"), strings.Index(index[body:], "<!-- /second-opinion-tree -->")
		return index[body+start : body+end]
	}
	block := treeBlock()
	listed := []string{}
	for _, match := range regexp.MustCompile(`href="([^"]+)"`).FindAllStringSubmatch(block, -1) {
		listed = append(listed, match[1])
	}
	entries, _ := os.ReadDir(bundle)
	sorted := []string{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "index-v1c") {
			sorted = append(sorted, entry.Name())
		}
	}
	want := "[index-v1c01-codex-actual-model-1-bold.html index-v1c01c01-codex-actual-model-1-visual.html index-v1c02-codex-actual-model-1-bold.html]"
	if fmt.Sprint(listed) != want || fmt.Sprint(sorted) != want {
		t.Fatalf("tree list order: listed %v, sorted %v", listed, sorted)
	}
	for _, field := range []string{"bold", "visual", "codex actual-model-1"} {
		if !strings.Contains(block, field) {
			t.Fatalf("tree list lacks %q: %s", field, block)
		}
	}
	if index, _ := os.ReadFile(filepath.Join(bundle, "index.html")); strings.Contains(string(index), "second-opinion-banner:") {
		t.Fatalf("tree page got a per-filename banner: %s", index)
	}
	row := secondOpinionRows(t, secondOpinionInvoke(t, repo, "--backend", "codex"))[0]
	index, _ := os.ReadFile(filepath.Join(bundle, "index.html"))
	if row["banner_state"] != "inserted" || !strings.Contains(string(index), "<!-- second-opinion-banner:index-codex-actual-model-1.html -->") || treeBlock() != block {
		t.Fatalf("plain run lost its banner or changed the tree block: %v\n%s", row, index)
	}
}
