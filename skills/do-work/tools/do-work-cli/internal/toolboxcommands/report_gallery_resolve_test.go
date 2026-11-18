package toolboxcommands

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

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
