package toolboxcommands

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const galleryFixturePage = "docs/report-gallery/index.html"

func TestAIReportGalleryListUsesProjectOverrideAndSkipsBrokenRecords(t *testing.T) {
	for _, toolbox := range []string{"skills/do-work-toolbox", ".claude/skills/do-work-toolbox"} {
		t.Run(toolbox, func(t *testing.T) {
			root := t.TempDir()
			shipped := toolbox + "/docs/report-gallery/"
			project := "docs/report-gallery/"
			writeGalleryFixture(t, root, shipped+"two-measures/sample.md", galleryFixtureSample("Two measures", "Compare a rate and a level.", "shipped", "", "2026-11-18"))
			writeGalleryFixture(t, root, shipped+"verdict-checks/sample.md", galleryFixtureSample("Verdict checks", "Show why the judgment holds.", "shipped", "", "2026-11-18"))
			writeGalleryFixture(t, root, project+"two-measures/sample.md", galleryFixtureSample("Local comparison", "Compare two project measures.", "https://example.test/report", "", "2026-11-18"))
			writeGalleryFixture(t, root, project+"added-sample/sample.md", "\ufeff"+strings.ReplaceAll(galleryFixtureSample("'An example: project evidence'", `"Does it look the same? Check the #details."`, "evidence.html", "", "2026-11-18"), "\n", "\r\n"))
			writeGalleryFixture(t, root, project+"broken-record/sample.md", galleryFixtureSample("Broken", "Compare numbers.", "shipped", "", "2026-02-30"))
			writeGalleryFixture(t, root, project+"broken-borrow/sample.md", strings.Replace(galleryFixtureSample("Broken", "Use for compare.", "shipped", "", "2026-11-18"), "- Name the unit.\n", "", 1))
			writeGalleryFixture(t, root, project+"broken-scalar/sample.md", galleryFixtureSample("Broken", `"Missing a closing quote.`, "shipped", "", "2026-11-18"))
			writeGalleryFixture(t, root, project+"verdict-checks/sample.md", "---\ntitle: Incomplete project override\n")
			writeGalleryFixture(t, root, project+"missing-record/notes.md", "No sample record.")
			writeGalleryFixture(t, root, project+"tmp/notes.md", "One-word directory.")
			writeGalleryFixture(t, root, "ai-reports", "Listing must not read the catalog.")
			writeGalleryFixture(t, root, galleryFixturePage, "A hand-made page must not affect listing.")
			before := snapshotReportTree(t, root)

			result := runGalleryFixture(t, root, "--list")
			want := "added-sample\tDoes it look the same? Check the #details.\tdocs/report-gallery/added-sample/sample.md\n" +
				"two-measures\tCompare two project measures.\tdocs/report-gallery/two-measures/sample.md\n" +
				"verdict-checks\tShow why the judgment holds.\t" + shipped + "verdict-checks/sample.md\n"
			if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || *result.ExactTextOutput != want {
				if result.ExactTextOutput == nil {
					t.Fatalf("missing list output: %+v", result)
				}
				t.Fatalf("want %q, got %q (%s, %v)", want, *result.ExactTextOutput, result.Outcome, result.Findings)
			}
			var records []string
			for _, finding := range result.Findings {
				if finding.Severity != resultmodel.SeverityWarning || len(finding.AffectedPaths) != 1 {
					t.Fatalf("expected a warning naming one sample, got %+v", finding)
				}
				records = append(records, finding.AffectedPaths[0])
			}
			wantWarnings := []string{project + "broken-borrow/sample.md", project + "broken-record/sample.md", project + "broken-scalar/sample.md", project + "missing-record/sample.md", project + "verdict-checks/sample.md"}
			if !reflect.DeepEqual(records, wantWarnings) {
				t.Errorf("warning records = %v, want %v", records, wantWarnings)
			}
			if after := snapshotReportTree(t, root); !reflect.DeepEqual(before, after) {
				t.Fatal("--list changed fixture files")
			}
		})
	}
}

func TestAIReportGalleryPageLinksLatestRealReportsAndAssets(t *testing.T) {
	root := t.TempDir()
	const toolbox = ".claude/skills/do-work-toolbox/docs/"
	const project = "docs/report-gallery/"
	writeGalleryFixture(t, root, toolbox+"ai-report-kit/presentation.css", ".band { display: block; }\n")
	writeGalleryFixture(t, root, toolbox+"second-opinion-presets/first-principles.md", "Work from the underlying constraints.\n")
	writeGalleryFixture(t, root, toolbox+"second-opinion-presets/not-a-preset.txt", "Ignore extra formats.")
	writeGalleryFixture(t, root, toolbox+"report-gallery/verdict-checks/sample.md", galleryFixtureSample("Verdict checks", "Show why the judgment holds.", "shipped", "", "2026-11-18"))
	writeGalleryFixture(t, root, toolbox+"report-gallery/verdict-checks/example.html", "<title>Fixture worked example</title>")
	writeGalleryFixture(t, root, toolbox+"report-gallery/two-measures/sample.md", galleryFixtureSample("Replaced shipped sample", "Compare a rate and a level.", "shipped", "", "2026-11-18"))
	writeGalleryFixture(t, root, project+"two-measures/sample.md", galleryFixtureSample(`"<script>alert(1)</script> & evidence"`, "Compare two project measures.", "shipped", `"images/rate #1.png"`, "2026-11-18"))
	writeGalleryFixture(t, root, project+"two-measures/images/rate #1.png", "fixture image")
	writeGalleryFixture(t, root, project+"linked-example/sample.md", galleryFixtureSample("External reference", "Answer what changed.", "https://example.test/view?a=1&b=2", "", "2026-11-18"))
	writeGalleryFixture(t, root, project+"local-example/sample.md", galleryFixtureSample("Local reference", "Answer what to keep.", `"notes/local #1.html"`, "", "2026-11-18"))
	writeGalleryFixture(t, root, "notes/local #1.html", "<title>Fixture local reference</title>")
	external := t.TempDir()
	writeGalleryFixture(t, external, "other reference.html", "<title>Fixture reference outside the project</title>")
	externalFile := filepath.Join(external, "other reference.html")
	externalRelative, err := filepath.Rel(filepath.Join(root, filepath.Dir(galleryFixturePage)), externalFile)
	if err != nil {
		t.Fatal(err)
	}
	externalLink := (&url.URL{Path: filepath.ToSlash(externalRelative)}).String()
	writeGalleryFixture(t, root, project+"absolute-example/sample.md", galleryFixtureSample("Absolute reference", "See the original report.", fmt.Sprintf("%q", externalFile), "", "2026-11-18"))
	writeGalleryFixture(t, root, project+"unsafe-source/sample.md", galleryFixtureSample("Unsafe reference", "Demonstrate the text fallback.", `"javascript:alert(1)"`, "", "2026-11-18"))
	writeGalleryFixture(t, root, project+"unsafe-image/sample.md", galleryFixtureSample("Unsafe image", "Demonstrate exclusion.", "shipped", "../two-measures/images/rate #1.png", "2026-11-18"))

	writeGalleryFixture(t, root, "ai-reports/catalog.json", "stale and hand-made")
	writeGalleryFixture(t, root, "ai-reports/2026-11-16_0900_old-proposal/index.html", `<title>Old proposal</title><meta name="ai-report-kind" content="proposal">`)
	writeGalleryFixture(t, root, "ai-reports/2026-11_1-report/2026-11-17_0900_new proposal/report #1.html", `<title>Newest &amp; best proposal</title><meta name="ai-report-kind" content="proposal">`)
	writeGalleryFixture(t, root, "ai-reports/2026-11-18_0900_empty-proposal/attachment.png", "no report")
	writeGalleryFixture(t, root, "ai-reports/2026-11-17_1000_REQ-201-completed/index.html", "<title>Completed work report</title>")
	writeGalleryFixture(t, root, "ai-reports/2026-11-17_1100_architecture-report/index.html", "<title>Architecture map</title>")
	writeGalleryFixture(t, root, "ai-reports/2026-11-17_1200_custom/index.html", `<title>Custom example</title><meta name="ai-report-kind" content="custom-view">`)
	before := snapshotReportTree(t, root)

	result := runGalleryFixture(t, root)
	if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || *result.ExactTextOutput != galleryFixturePage+"\n" {
		t.Fatalf("generation failed: %s (%v)", result.Outcome, result.Findings)
	}
	if len(result.Changes) != 1 || result.Changes[0].Path != galleryFixturePage || result.Changes[0].Kind != "created" {
		t.Errorf("expected only a gallery page creation, got %+v", result.Changes)
	}
	after := snapshotReportTree(t, root)
	page := after[filepath.Join(root, filepath.FromSlash(galleryFixturePage))]
	delete(after, filepath.Join(root, filepath.FromSlash(galleryFixturePage)))
	if !reflect.DeepEqual(before, after) {
		t.Fatal("generation wrote outside the gallery index")
	}
	for _, want := range []string{
		`name="generator" content="do-work-cli ai-report-gallery"`,
		`href="../../.claude/skills/do-work-toolbox/docs/ai-report-kit/presentation.css"`,
		`href="../../.claude/skills/do-work-toolbox/docs/report-gallery/verdict-checks/example.html"`,
		`href="two-measures/images/rate%20%231.png"`,
		`href="https://example.test/view?a=1&amp;b=2"`,
		`href="../../notes/local%20%231.html"`,
		`href="` + externalLink + `"`,
		`href="../../ai-reports/2026-11_1-report/2026-11-17_0900_new%20proposal/report%20%231.html"`,
		`href="../../ai-reports/2026-11-17_1000_REQ-201-completed/index.html"`,
		`href="../../ai-reports/2026-11-17_1100_architecture-report/index.html"`,
		`href="../../ai-reports/2026-11-17_1200_custom/index.html"`,
		`href="../../.claude/skills/do-work-toolbox/docs/second-opinion-presets/first-principles.md"`,
		"Newest &amp; best proposal", "Compare two project measures.", "Borrow", "Name the unit.",
		"ai-report --kind root-cause", "ai-report gallery add", "Copy", "clipboard.writeText", "execCommand('copy')",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("page missing %q", want)
		}
	}
	for _, absent := range []string{"Old proposal", "empty-proposal", "Replaced shipped sample", "not-a-preset", "<script>alert(1)", "href=\"javascript:", "Unsafe image"} {
		if strings.Contains(page, absent) {
			t.Errorf("page unexpectedly contains %q", absent)
		}
	}
	for _, match := range regexp.MustCompile(`href="([^"#]+)"`).FindAllStringSubmatch(page, -1) {
		link, err := url.Parse(match[1])
		if err != nil {
			t.Fatal(err)
		}
		if link.IsAbs() {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, filepath.Dir(galleryFixturePage), filepath.FromSlash(link.Path))); err != nil {
			t.Errorf("unresolved generated link %q: %v", match[1], err)
		}
	}

	changed := runGalleryFixture(t, root)
	if changed.Outcome != resultmodel.OutcomeSuccess || len(changed.Changes) != 1 || changed.Changes[0].Kind != "modified" {
		t.Fatalf("repeat generation did not accept its own marker: %+v", changed)
	}
	if got := snapshotReportTree(t, root)[filepath.Join(root, filepath.FromSlash(galleryFixturePage))]; got != page {
		t.Error("same inputs produced different page bytes")
	}
}

func TestAIReportGalleryEmptyListWritesNothingAndSourceCanCreatePage(t *testing.T) {
	root := t.TempDir()
	writeGalleryFixture(t, root, "skills/do-work-toolbox/docs/ai-report-kit/presentation.css", "")
	writeGalleryFixture(t, root, "skills/do-work-toolbox/docs/second-opinion-presets/feynman.md", "Explain it simply.")
	var output bytes.Buffer
	runtime := commandruntime.NewRuntime(&output, Handlers())
	if code := runtime.Run([]string{"--repo-root", root, "ai-report-gallery", "--list"}); code != 0 || output.Len() != 0 {
		t.Fatalf("empty list: exit %d, stdout %q", code, output.String())
	}
	if _, err := os.Stat(filepath.Join(root, "docs")); !os.IsNotExist(err) {
		t.Fatalf("--list created the output directory: %v", err)
	}
	if code := runtime.Run([]string{"--repo-root", root, "--format", "text", "ai-report-gallery"}); code != 0 || output.String() != galleryFixturePage+"\n" {
		t.Fatalf("updater invocation: exit %d, stdout %q", code, output.String())
	}
	page, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(galleryFixturePage)))
	if err != nil || !bytes.Contains(page, []byte(`href="../../skills/do-work-toolbox/docs/second-opinion-presets/feynman.md"`)) {
		t.Fatalf("source layout preset link missing: %v\n%s", err, page)
	}
}

func TestAIReportGalleryRefusesHandMadeOrLinkedOutput(t *testing.T) {
	for _, shape := range []string{"hand-made", "marker in comment", "linked file", "linked parent", "directory as page"} {
		t.Run(shape, func(t *testing.T) {
			root := t.TempDir()
			target := filepath.Join(root, filepath.FromSlash(galleryFixturePage))
			outside := t.TempDir()
			writeGalleryFixture(t, outside, "index.html", `<meta name="generator" content="do-work-cli ai-report-gallery">`)
			switch shape {
			case "hand-made":
				writeGalleryFixture(t, root, galleryFixturePage, "Hand-made page")
			case "marker in comment":
				writeGalleryFixture(t, root, galleryFixturePage, `<html><!-- <meta name="generator" content="do-work-cli ai-report-gallery"> --></html>`)
			case "directory as page":
				if err := os.MkdirAll(target, 0o755); err != nil {
					t.Fatal(err)
				}
			case "linked file", "linked parent":
				destination := filepath.Join(outside, "index.html")
				if shape == "linked parent" {
					target, destination = filepath.Dir(target), outside
				}
				if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(destination, target); err != nil {
					t.Fatal(err)
				}
			}
			before, outsideBefore := snapshotReportTree(t, root), snapshotReportTree(t, outside)
			result := runGalleryFixture(t, root)
			if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) != 1 || !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), galleryFixturePage) {
				t.Fatalf("unsafe output was not refused with its path: %+v", result)
			}
			if !reflect.DeepEqual(before, snapshotReportTree(t, root)) || !reflect.DeepEqual(outsideBefore, snapshotReportTree(t, outside)) {
				t.Fatal("refused write changed files")
			}
		})
	}
}

func TestAIReportGalleryPublicationPreservesAnEditAfterOwnershipCheck(t *testing.T) {
	root := t.TempDir()
	if result := runGalleryFixture(t, root); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("seed generated page: %+v", result)
	}
	const handMade = "Hand-made edit during regeneration."
	filename := filepath.Join(root, filepath.FromSlash(galleryFixturePage))
	err := rootedPublishFileChecked(root, galleryFixturePage, []byte("replacement"), 0o644, true, func(original []byte) error {
		if !bytes.Contains(original, []byte(`name="generator" content="do-work-cli ai-report-gallery"`)) {
			t.Fatal("ownership check did not receive the original generated page")
		}
		return os.WriteFile(filename, []byte(handMade), 0o644)
	})
	if err == nil || !strings.Contains(err.Error(), "target changed") {
		t.Fatalf("expected a target-change refusal, got %v", err)
	}
	if content, err := os.ReadFile(filename); err != nil || string(content) != handMade {
		t.Fatalf("concurrent edit was lost: %q (%v)", content, err)
	}
	entries, err := os.ReadDir(filepath.Dir(filename))
	if err != nil || len(entries) != 1 {
		t.Fatalf("publication left a staged file: %v (%v)", entries, err)
	}
}

func galleryFixtureSample(title, useWhen, source, image, added string) string {
	fields := fmt.Sprintf("---\ntitle: %s\nuse_when: %s\nsource: %s\nadded: %s\n", title, useWhen, source, added)
	if image != "" {
		fields += "image: " + image + "\n"
	}
	return fields + "---\n\n## Borrow\n\n- State the finding.\n- Name the unit.\n"
}

func writeGalleryFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	filename := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runGalleryFixture(t *testing.T, root string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	handler := Handlers()["ai-report-gallery"]
	if handler == nil {
		t.Fatal("ai-report-gallery is not registered")
	}
	return handler(commandruntime.ExecutionContext{RepositoryRoot: root, Format: resultmodel.FormatText}, arguments)
}
