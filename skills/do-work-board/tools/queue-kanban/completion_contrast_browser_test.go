package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Exercise real generated cards, not isolated spans: inherited opacity and the
// card background both affect whether the completion readings remain legible.
func TestBrowserBehaviorCompletionCompanionsKeepReadableContrast(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)
	// The client filters Recently Done against its wall clock, even for a static
	// snapshot, so keep the fixture within that window on every future run.
	moment := time.Now().UTC()
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{{
		"do-work/archive/REQ-901-completion-contrast.md",
		spanFixtureFrontmatter("REQ-901", "Keep completion readings legible", "completed",
			moment.Add(-3*time.Hour).Format(time.RFC3339), moment.Add(-15*time.Minute).Format(time.RFC3339)),
	}, {
		"do-work/queue/REQ-902-pending-contrast.md",
		spanFixtureFrontmatter("REQ-902", "Pending control", "pending", "", "",
			"status_changed_at: "+moment.Add(-time.Hour).Format(time.RFC3339)),
	}, {
		"do-work/working/REQ-903-claimed-contrast.md",
		spanFixtureFrontmatter("REQ-903", "Claimed control", "claimed", moment.Add(-time.Hour).Format(time.RFC3339), ""),
	}})
	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow,
		func(string, string) (time.Time, bool) { return time.Time{}, false })
	if buildError != nil {
		t.Fatal(buildError)
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, board); generateError != nil {
		t.Fatal(generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatal(readError)
	}
	session := startTrustedInputBrowserSession(t, "completion companion contrast", siteDirectory, string(indexBytes))
	session.waitForPageCondition(t, "a completed card with both companion readings",
		`document.querySelector('.req-card[data-status="completed"] .elapsed-duration') !== null`)
	for _, scheme := range []string{"light", "dark"} {
		for _, width := range []int{320, 768, 1280} {
			t.Run(fmt.Sprintf("%s-%d", scheme, width), func(t *testing.T) {
				session.callDevToolsMethod(t, "Emulation.setEmulatedMedia", map[string]any{
					"features": []map[string]string{{"name": "prefers-color-scheme", "value": scheme}},
				}, true)
				session.callDevToolsMethod(t, "Emulation.setDeviceMetricsOverride", map[string]any{
					"width": width, "height": 900, "deviceScaleFactor": 1, "mobile": false,
				}, true)
				var measured struct {
					Href, Browser, Scheme, Body, Card string
					FaintInk                          string
					Title                             completionTextStyle
					Companions                        []completionTextStyle
					Controls                          []completionTextStyle
				}
				session.decodeResult(t, "completion text styles", session.evaluateInPage(t, `(function () {
  var card = document.querySelector('.req-card[data-status="completed"]');
  function textStyle(node) {
    var style = getComputedStyle(node);
    var opacity = 1;
    for (var ancestor = node; ancestor; ancestor = ancestor.parentElement) {
      opacity *= Number(getComputedStyle(ancestor).opacity);
    }
    var bounds = node.getBoundingClientRect();
    return {name: node.className, text: node.textContent, color: style.color,
      opacity: opacity, fontSize: parseFloat(style.fontSize), fontWeight: Number(style.fontWeight),
      width: bounds.width, height: bounds.height};
  }
  return {href: location.href, browser: navigator.userAgent,
    scheme: matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light',
    body: getComputedStyle(document.body).backgroundColor,
    card: getComputedStyle(card).backgroundColor,
    faintInk: getComputedStyle(card).getPropertyValue('--ink-faint'),
    controls: Array.from(document.querySelectorAll('.req-card:is([data-status="pending"], [data-status="claimed"]) .elapsed-duration'), textStyle),
    title: textStyle(card.querySelector('.req-card-title')),
    companions: Array.from(card.querySelectorAll('.req-card-completed .relative-time, .req-card-completed .elapsed-duration'), textStyle)};
})()`), &measured)
				if measured.Scheme != scheme || len(measured.Companions) != 2 {
					t.Fatalf("scheme/readings = %q/%d, want %q/2", measured.Scheme, len(measured.Companions), scheme)
				}
				if len(measured.Controls) != 2 {
					t.Fatalf("nonterminal control readings = %d, want pending and claimed", len(measured.Controls))
				}
				for _, control := range measured.Controls {
					if normalizeCSSColour(control.Color) != normalizeCSSColour(measured.FaintInk) || control.Opacity != 0.85 {
						t.Errorf("completion fix changed a nonterminal timer: %+v", control)
					}
				}
				t.Logf("%s; %s; body=%s card=%s", measured.Href, measured.Browser, measured.Body, measured.Card)
				for _, background := range []string{measured.Body, measured.Card} {
					titleContrast := completionTextContrast(t, measured.Title, background)
					for _, companion := range measured.Companions {
						if companion.Text == "" || companion.Width <= 0 || companion.Height <= 0 || companion.FontSize <= 0 {
							t.Fatalf("completion companion is not rendered: %+v", companion)
						}
						contrast := completionTextContrast(t, companion, background)
						if contrast < 4.5 {
							t.Errorf("%s contrast %.2f:1 against %s, want >=4.5:1", companion.Name, contrast, background)
						}
						// Preserve a meaningful hierarchy, not merely a fractional difference:
						// the companions have at most 75% of the title's contrast and smaller,
						// lighter type. The screenshot review checks the combined appearance.
						if contrast > titleContrast*0.75 || companion.FontSize >= measured.Title.FontSize || companion.FontWeight >= measured.Title.FontWeight {
							t.Errorf("%s no longer reads quieter than the title: contrast %.2f/%.2f, size %.0f/%.0f, weight %.0f/%.0f",
								companion.Name, contrast, titleContrast, companion.FontSize, measured.Title.FontSize, companion.FontWeight, measured.Title.FontWeight)
						}
						t.Logf("%s: color=%s opacity=%.2f contrast=%.2f title=%.2f against %s",
							companion.Name, companion.Color, companion.Opacity, contrast, titleContrast, background)
					}
				}
			})
		}
	}
}

type completionTextStyle struct {
	Name, Text, Color             string
	Opacity, FontSize, FontWeight float64
	Width, Height                 float64
}

func completionTextContrast(t *testing.T, textStyle completionTextStyle, background string) float64 {
	t.Helper()
	red, green, blue, alpha, parsed := parseCSSColourChannels(textStyle.Color)
	backRed, backGreen, backBlue, backAlpha, backParsed := parseCSSColourChannels(background)
	if !parsed || !backParsed || backAlpha != 1 {
		t.Fatalf("cannot measure text %q on opaque background %q", textStyle.Color, background)
	}
	// Composite sRGB channels BEFORE linearizing for WCAG luminance. Blending
	// luminances instead would overstate the faint text's contrast in dark mode.
	alpha *= textStyle.Opacity
	textLuminance := 0.2126*srgbToLinear(red*alpha+backRed*(1-alpha)) +
		0.7152*srgbToLinear(green*alpha+backGreen*(1-alpha)) +
		0.0722*srgbToLinear(blue*alpha+backBlue*(1-alpha))
	backgroundLuminance := 0.2126*srgbToLinear(backRed) + 0.7152*srgbToLinear(backGreen) + 0.0722*srgbToLinear(backBlue)
	return contrastRatio(textLuminance, backgroundLuminance)
}
