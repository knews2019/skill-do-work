package toolboxcommands

import (
	"fmt"
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestReportGalleryStarterRoundTripAcrossSources(t *testing.T) {
	root := t.TempDir()
	const shipped = ".claude/skills/do-work-toolbox/docs/report-gallery/shipped-design/sample.md"
	const added = "docs/report-gallery/added-design/sample.md"
	const generated = "ai-reports/2026-10_1-report/generated-example"
	writeGalleryFixture(t, root, shipped, galleryFixtureSample("Shipped design", "Explain a sequence.", "shipped", "", "2026-11-18"))
	writeGalleryFixture(t, root, added, galleryFixtureSample(`'Cost "$x" & rate'`, "Compare two rates.", "shipped", "", "2026-11-18"))
	writeGalleryFixture(t, root, generated+"/index.html", "<title>Old evidence</title>")
	writeGalleryFixture(t, root, generated+"/gallery-design.md", galleryFixtureDesign("Generated design", "Show a handoff."))

	assertGallery := func(t *testing.T, wantPaths map[string]string) {
		t.Helper()
		if result := runGalleryFixture(t, root); result.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("render gallery: %+v", result)
		}
		assertGalleryStarterRecords(t, root, wantPaths)
	}
	wantPaths := map[string]string{
		"shipped-design":           shipped,
		"added-design":             added,
		"report:generated-example": generated + "/gallery-design.md",
	}
	t.Run("installed", func(t *testing.T) { assertGallery(t, wantPaths) })

	var cards []string
	kinds := []string{"Completed work", "Proposal", "Root cause", "Architecture"}
	for index, suffix := range strings.Fields("one two three four five six seven eight nine ten") {
		name := "report-" + suffix
		report := "ai-reports/" + name + "/index.html"
		image := "docs/report-gallery/previews/" + name + ".png"
		writeGalleryFixture(t, root, report, "<title>Legacy evidence</title><main><figure><figcaption>Sequence diagram</figcaption></figure></main>")
		writeGalleryFixture(t, root, image, "screenshot fixture")
		cards = append(cards, fmt.Sprintf(`{"name":%q,"kind":%q,"subject":"Actual evidence","report":%q,"image":%q,"starter":{"id":"injected-prompt","prompt":"do-work-toolbox ai-report --design \"Added design\""}}`, name, kinds[index%len(kinds)], report, image))
		wantPaths["report:"+name] = report
	}
	writeGalleryFixture(t, root, "docs/report-gallery/preview-manifest.json", `{"width":900,"height":560,"cards":[`+strings.Join(cards, ",")+`]}`)
	t.Run("curated", func(t *testing.T) { assertGallery(t, wantPaths) })

	t.Run("malformed-curated-record", func(t *testing.T) {
		writeGalleryFixture(t, root, "ai-reports/report-one/gallery-design.md", "## Borrow\n- Missing metadata.")
		delete(wantPaths, "report:report-one")
		assertGallery(t, wantPaths)
	})

	listing := runGalleryFixture(t, root, "--list")
	wantListing := "added-design\tCompare two rates.\t" + added + "\nshipped-design\tExplain a sequence.\t" + shipped + "\n"
	if listing.ExactTextOutput == nil || *listing.ExactTextOutput != wantListing {
		t.Errorf("automatic ranking rows: want %q, got %+v", wantListing, listing)
	}
}

func assertGalleryStarterRecords(t *testing.T, root string, wantPaths map[string]string) {
	t.Helper()
	page, err := os.ReadFile(filepath.Join(root, galleryFixturePage))
	if err != nil {
		t.Fatal(err)
	}
	prompts := regexp.MustCompile(`<code class="design-prompt" id="([^"]+)">([^<]+)</code>\s*<button [^>]*data-copy="([^"]+)"`)
	seen := map[string]bool{}
	targetIDs := map[string]bool{}
	for _, match := range prompts.FindAllStringSubmatch(string(page), -1) {
		id, prompt, copyTarget := match[1], html.UnescapeString(match[2]), match[3]
		if targetIDs[id] || id != copyTarget {
			t.Errorf("starter %q has duplicate ID or incorrect copy target %q", id, copyTarget)
		}
		targetIDs[id] = true
		// Let the shell parse the pasted command and inspect its argv without running it.
		parsed, err := exec.Command("sh", "-c", "set -- "+prompt+"\nprintf '%s\\n' \"$@\"").Output()
		if err != nil {
			t.Fatalf("parse starter %q: %v", prompt, err)
		}
		args := strings.Split(strings.TrimSuffix(string(parsed), "\n"), "\n")
		if len(args) != 4 || strings.Join(args[:3], " ") != "do-work-toolbox ai-report --design" {
			t.Errorf("starter %q parsed to %q", prompt, args)
			continue
		}
		result := runGalleryFixture(t, root, "--resolve", args[3])
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil {
			t.Errorf("resolve starter %q: %+v", prompt, result)
			continue
		}
		row := strings.Split(strings.TrimSuffix(*result.ExactTextOutput, "\n"), "\t")
		if len(row) != 3 || row[1] == "" || wantPaths[row[0]] != row[2] {
			t.Errorf("starter %q resolved to unexpected row %q", prompt, row)
			continue
		}
		if seen[row[0]] {
			t.Errorf("design %q appears twice", row[0])
		}
		seen[row[0]] = true
	}
	for name := range wantPaths {
		if !seen[name] {
			t.Errorf("no usable starter for %q", name)
		}
	}
}

func galleryFixtureDesign(title, fit string) string {
	return fmt.Sprintf("---\ntitle: %s\nuse_when: %s\n---\n\n## Borrow\n\n- Trace the handoff in a diagram.\n- Put the measure beside its receipt.\n", title, fit)
}

func TestReportGalleryGeneratedResolveQualifiesCollisions(t *testing.T) {
	root := t.TempDir()
	const loose = "ai-reports/same-design/"
	const organized = "ai-reports/2026-09_1-report/same-design/"
	const sample = "docs/report-gallery/project-design/sample.md"
	writeGalleryFixture(t, root, sample, galleryFixtureSample("report:same-design", "Explain a checklist.", "shipped", "", "2026-11-18"))
	for record, fit := range map[string]string{loose: "Show a decision.", organized: "Show a handoff."} {
		writeGalleryFixture(t, root, record+"index.html", "<title>Earlier report</title>")
		writeGalleryFixture(t, root, record+"gallery-design.md", galleryFixtureDesign("Same title", fit))
	}
	const orphan = "ai-reports/orphan-design/gallery-design.md"
	writeGalleryFixture(t, root, orphan, galleryFixtureDesign("Orphan", "Show a change."))
	const malformed = "ai-reports/broken-design/"
	writeGalleryFixture(t, root, malformed+"index.html", "<title>Earlier report</title>")
	writeGalleryFixture(t, root, malformed+"gallery-design.md", "## Borrow\n- Missing metadata.")
	for _, name := range []string{"report:same-design", "report:unknown", "report:orphan-design", "report:broken-design"} {
		result := runGalleryFixture(t, root, "--resolve", name)
		if result.Outcome != resultmodel.OutcomeFailure || result.ExactTextOutput != nil {
			t.Errorf("resolve %q: expected failure, got %+v", name, result)
		}
	}
	if result := runGalleryFixture(t, root); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("render gallery: %+v", result)
	}
	assertGalleryStarterRecords(t, root, map[string]string{
		"project-design":      sample,
		"report:" + loose:     loose + "gallery-design.md",
		"report:" + organized: organized + "gallery-design.md",
	})
}

func TestReportGalleryResolveUsesDiscoveredNames(t *testing.T) {
	root := t.TempDir()
	for record, title := range map[string]string{
		"skills/do-work-toolbox/docs/report-gallery/local-design/sample.md": "Shipped title",
		"docs/report-gallery/local-design/sample.md":                        "Local title",
		"docs/report-gallery/other-design/sample.md":                        "local-design",
		"docs/report-gallery/first-match/sample.md":                         "Same title",
		"docs/report-gallery/second-match/sample.md":                        "Same Title",
	} {
		writeGalleryFixture(t, root, record, galleryFixtureSample(title, "Fit.", "source", "", "2026-11-18"))
	}
	writeGalleryFixture(t, root, galleryFixturePage, "Hand-made page.")
	before := snapshotReportTree(t, root)
	for _, name := range []string{"  LOCAL TITLE ", "local-design"} {
		result := runGalleryFixture(t, root, "--resolve", name)
		want := "local-design\tFit.\tdocs/report-gallery/local-design/sample.md\n"
		if result.Outcome != resultmodel.OutcomeSuccess || result.ExactTextOutput == nil || *result.ExactTextOutput != want {
			t.Errorf("resolve %q: %+v", name, result)
		}
	}
	for name, guidance := range map[string]string{
		"Shipped title": "--list",
		"Local":         "--list",
		"Same Title":    "first-match, second-match",
	} {
		result := runGalleryFixture(t, root, "--resolve", name)
		if result.Outcome != resultmodel.OutcomeFailure || result.ExactTextOutput != nil ||
			!strings.Contains(fmt.Sprint(result.Findings), guidance) {
			t.Errorf("resolve %q: want failure with %q, got %+v", name, guidance, result)
		}
	}
	if after := snapshotReportTree(t, root); !reflect.DeepEqual(before, after) {
		t.Fatal("--resolve changed fixture files")
	}
}

func TestReportGalleryCopyNamesRoundTrip(t *testing.T) {
	samples := []reportGallerySample{
		{Name: "unique-design", Title: "A Reader’s Guide"},
		{Name: "first-match", Title: "Same Title"},
		{Name: "second-match", Title: "Same title"},
		{Name: "other-design", Title: "unique-design"},
	}
	names := reportGalleryCopyNames(samples)
	if names["unique-design"] != "A Reader’s Guide" {
		t.Errorf("unique title was changed: %q", names["unique-design"])
	}
	for _, sample := range samples {
		name := names[sample.Name]
		matches := matchReportGalleryName(samples, name)
		if len(matches) != 1 || matches[0].Name != sample.Name {
			t.Errorf("copy name %q for %s resolves to %v", name, sample.Name, matches)
		}
	}
}
