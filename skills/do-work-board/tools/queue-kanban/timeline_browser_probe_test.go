package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// The Timeline's status colours, measured in a real engine.
//
// WHY A BROWSER. What this REQ delivered is a cascade: a status attribute on the
// row group resolves a custom property, the property resolves a token, the token
// resolves per theme, and two segments read that one property at two opacities.
// Every link in that chain is CSS. A Node probe can assert the class names the
// renderer writes and nothing at all about the colours they produce, which is
// the half a reader actually sees — and asserting the attribute while the
// property silently resolved to nothing is exactly the "successful-looking
// measurement of nothing" REQ-291's lesson warns about. So this asks the engine
// for computed fills and refuses to pass on an empty one.
//
// The row markup here is built to the SHIPPED renderer's shape and pinned to it
// by TestTimelineStatusRowMarkupMatchesTheProbe below: if renderVisibleRows stops
// writing data-status, that test fails rather than this one passing against a
// fixture nothing produces any more.

// timelineStatusProbeStatuses is every status the board's own vocabulary can put
// on a row, plus one deliberate typo. The typo is the point: an exact-match rule
// sends it to the unrecognized colour, and a prefix match would colour it as
// real blocked work.
var timelineStatusProbeStatuses = []string{
	"pending",
	"claimed",
	"completed",
	"completed-with-issues",
	"pending-answers",
	"blocked",
	"blocked-archive-collision",
	"blocked-dependency-cycle",
	"failed",
	"cancelled",
	"blockd-dependency-cycle",
}

type timelineStatusProbeRow struct {
	Status          string  `json:"status"`
	WaitFill        string  `json:"waitFill"`
	WorkFill        string  `json:"workFill"`
	WaitOpacity     float64 `json:"waitOpacity"`
	WorkOpacity     float64 `json:"workOpacity"`
	ChipAccent      string  `json:"chipAccent"`
	Unrecognized    bool    `json:"unrecognized"`
	WaitStroke      string  `json:"waitStroke"`
	WaitStrokeWidth float64 `json:"waitStrokeWidth"`
}

func TestBrowserBehaviorTimelineBarsCarryTheirStatusColour(t *testing.T) {
	indexHtml := generateLiveSite(t)
	styleBlock := sliceGeneratedStyleBlock(t, indexHtml)

	rowMarkup := ""
	chipMarkup := ""
	for _, status := range timelineStatusProbeStatuses {
		unrecognizedClass := ""
		if status == "blockd-dependency-cycle" {
			unrecognizedClass = " is-status-unrecognized"
		}
		rowMarkup += `<g class="timeline-row` + unrecognizedClass + `" data-status="` + status + `">` +
			`<rect class="timeline-segment timeline-segment-wait" width="40" height="10"></rect>` +
			`<rect class="timeline-segment timeline-segment-work" x="40" width="40" height="10"></rect>` +
			`</g>`
		// The same status on a calendar chip, so the comparison is against what
		// the board actually paints elsewhere rather than against a token name
		// copied into this file.
		chipMarkup += `<button class="calendar-chip` + unrecognizedClass + `" data-status="` + status + `"></button>`
	}

	pageHTML := `<!doctype html><html><head><meta charset="utf-8"><style>` + styleBlock + `</style></head><body>
<svg class="timeline-rows-svg" width="200" height="200">` + rowMarkup + `</svg>
<div class="calendar-day-entries">` + chipMarkup + `</div>
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
(function () {
  var statuses = ` + mustJSONList(t, timelineStatusProbeStatuses) + `;
  var rows = document.querySelectorAll(".timeline-row");
  var chips = document.querySelectorAll(".calendar-chip");
  var measured = statuses.map(function (status, index) {
    var row = rows[index];
    var waitStyle = getComputedStyle(row.querySelector(".timeline-segment-wait"));
    var workStyle = getComputedStyle(row.querySelector(".timeline-segment-work"));
    return {
      status: status,
      waitFill: waitStyle.fill,
      workFill: workStyle.fill,
      waitOpacity: Number(waitStyle.fillOpacity),
      workOpacity: Number(workStyle.fillOpacity),
      chipAccent: getComputedStyle(chips[index]).getPropertyValue("--chip-accent").trim(),
      unrecognized: row.classList.contains("is-status-unrecognized"),
      waitStroke: waitStyle.stroke,
      waitStrokeWidth: parseFloat(waitStyle.strokeWidth) || 0
    };
  });
  // Written LAST and only once every measurement exists, so a throw leaves the
  // node empty and the Go side reports a failure rather than reading a partial
  // object as a pass.
  document.getElementById("` + browserProbeResultElementId + `").textContent =
    JSON.stringify({
      rows: measured,
      // The scheme the ENGINE resolved, not the one the flag asked for. A flag
      // this build silently ignores would otherwise let one palette be measured
      // twice and reported as two.
      scheme: matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
      surface: getComputedStyle(document.body).backgroundColor
    });
})();
</script>
</body></html>`

	// BOTH schemes. This board is dark-first — :root is the dark palette and
	// @media (prefers-color-scheme: light) overrides it — and Chromium resolves
	// light with no flag, so a single run measures exactly one of the two
	// palettes and the other is checked by nothing.
	for _, scheme := range []struct {
		name string
		flag string
	}{
		{name: "light", flag: "--blink-settings=preferredColorScheme=1"},
		{name: "dark", flag: "--blink-settings=preferredColorScheme=0"},
	} {
		t.Run(scheme.name, func(t *testing.T) {
			assertTimelineStatusColours(t, pageHTML, scheme.name, scheme.flag)
		})
	}
}

// timelineContrastFloor is the board's own floor for a graphical object, stated
// at web/board.css's route-ramp comment and applied there to adjacent ramp steps.
// The two halves of a bar are the same kind of thing: two adjacent shapes a
// reader has to tell apart.
const timelineContrastFloor = 3.0

// timelineSurfaceFloor is deliberately lower. The wait against the page is a
// fill-versus-background question, not two marks side by side, and the wait
// carries a full-strength outline that the fill ratio does not capture. Below
// this it stops reading as a bar at all.
const timelineSurfaceFloor = 1.3

func assertTimelineStatusColours(t *testing.T, pageHTML string, schemeName string, schemeFlag string) {
	t.Helper()
	probeOutput := runBrowserBehaviorProbeWithFlags(
		t, "timeline status colours ("+schemeName+")", pageHTML, schemeFlag)
	// The lane's sentinel requires a JSON object at the result node, so the rows
	// arrive wrapped.
	var probeResult struct {
		Rows    []timelineStatusProbeRow `json:"rows"`
		Scheme  string                   `json:"scheme"`
		Surface string                   `json:"surface"`
	}
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode timeline status colour probe: %v (output %q)", decodeError, probeOutput)
	}
	measured := probeResult.Rows
	if probeResult.Scheme != schemeName {
		t.Fatalf("asked the engine for the %s palette and it resolved %s; this build ignores the "+
			"colour-scheme flag, so one palette would be measured twice and reported as two",
			schemeName, probeResult.Scheme)
	}
	surfaceLuminance, surfaceKnown := relativeLuminanceOfCSSColour(probeResult.Surface)
	if !surfaceKnown {
		t.Fatalf("could not read the page surface colour (%q); every contrast number below is "+
			"measured against it", probeResult.Surface)
	}
	if len(measured) != len(timelineStatusProbeStatuses) {
		t.Fatalf("probe measured %d rows, want %d", len(measured), len(timelineStatusProbeStatuses))
	}

	byStatus := map[string]timelineStatusProbeRow{}
	for _, row := range measured {
		// A fill of "none", empty, or fully transparent means the custom property
		// resolved to nothing — the failure mode where every assertion below would
		// otherwise agree with every other on the same non-colour.
		if row.WorkFill == "" || row.WorkFill == "none" || fullyTransparentCSSColour(row.WorkFill) {
			t.Fatalf("status %q painted its work segment %q — the status accent resolved to nothing",
				row.Status, row.WorkFill)
		}
		byStatus[row.Status] = row
	}

	// 1. Hue is the status, and it is the SAME hue the calendar chip uses. This is
	//    the requirement in one assertion: a REQ must not read as one colour on a
	//    chip and another on a bar.
	for _, row := range measured {
		if row.ChipAccent == "" {
			t.Fatalf("status %q resolved no --chip-accent; the comparison target is missing", row.Status)
		}
		if !sameCSSColour(row.WorkFill, row.ChipAccent) {
			t.Errorf("status %q paints its bar %s but its calendar chip %s; one REQ must be one colour",
				row.Status, row.WorkFill, row.ChipAccent)
		}
	}

	// 2. Lightness is the phase, and it is MEASURED rather than merely ordered.
	//
	//    The first version of this block asserted only waitOpacity < workOpacity
	//    and waitOpacity >= 0.25. That passes at 0.98 and at 0.26 — it pins the
	//    direction of a difference and says nothing about whether a reader can
	//    see it, in a change whose entire risk was exactly that. It also passed
	//    while the light theme sat at 2.47:1, under this board's own 3:1 floor.
	//    So: real contrast, computed from the composited fills, against two
	//    floors that mean different things.
	for _, row := range measured {
		if row.WaitFill != row.WorkFill {
			t.Errorf("status %q paints wait %s and work %s; the phase difference must be lightness, "+
				"not a second hue", row.Status, row.WaitFill, row.WorkFill)
		}
		if !(row.WaitOpacity < row.WorkOpacity) {
			t.Errorf("status %q has wait opacity %.2f and work opacity %.2f; the wait must be the "+
				"quieter half", row.Status, row.WaitOpacity, row.WorkOpacity)
		}

		accentLuminance, accentKnown := relativeLuminanceOfCSSColour(row.WorkFill)
		if !accentKnown {
			t.Fatalf("status %q resolved an unreadable fill %q", row.Status, row.WorkFill)
		}
		// A fill at alpha composites over the page, so the visible wait colour is
		// the blend — not the accent, and not the accent's own luminance.
		waitLuminance := compositeLuminance(accentLuminance, surfaceLuminance, row.WaitOpacity)
		workLuminance := compositeLuminance(accentLuminance, surfaceLuminance, row.WorkOpacity)

		waitAgainstWork := contrastRatio(waitLuminance, workLuminance)
		waitAgainstSurface := contrastRatio(waitLuminance, surfaceLuminance)

		// Wait vs work is two adjacent marks a reader has to tell apart, and there
		// are two ways to be separable: the fills are far enough apart on their
		// own, OR the wait's outline is visible against the wait's own fill and
		// draws the boundary. Either satisfies this; neither being true does not.
		//
		// The second clause is what stops the outline being a licence. An outline
		// is only a channel while it can be SEEN, and it is drawn in the accent at
		// full strength — so as the wait's fill rises towards full strength the
		// outline disappears into it. At alpha 0.98 the fills are one colour and
		// the outline is invisible against them: both clauses fail, which is the
		// mutation the previous version of this test passed.
		outlineAgainstWaitFill := contrastRatio(accentLuminance, waitLuminance)
		fillsAreSeparable := waitAgainstWork >= timelineContrastFloor
		outlineIsVisible := row.WaitStrokeWidth > 0 && outlineAgainstWaitFill >= timelineSurfaceFloor
		if !fillsAreSeparable && !outlineIsVisible {
			t.Errorf("status %q separates its wait from its work by %.2f:1 of fill (floor %.1f:1) "+
				"and its outline reads %.2f:1 against its own fill (floor %.1f:1) at stroke-width "+
				"%.2f; one of the two has to carry the difference",
				row.Status, waitAgainstWork, timelineContrastFloor,
				outlineAgainstWaitFill, timelineSurfaceFloor, row.WaitStrokeWidth)
		}
		if waitAgainstSurface < timelineSurfaceFloor {
			t.Errorf("status %q draws its wait at %.2f:1 against the page; below %.1f:1 it stops "+
				"reading as a bar at the shipped 10px height",
				row.Status, waitAgainstSurface, timelineSurfaceFloor)
		}
		// The outline is the second channel the light palette needs, so it is a
		// requirement rather than a decoration: without it, opacity alone has to
		// clear the floor above and cannot, in this theme, at any single value.
		if row.WaitStrokeWidth <= 0 {
			t.Errorf("status %q draws its wait with no outline (stroke-width %.2f); opacity alone "+
				"cannot separate wait from work and wait from the surface at the same time",
				row.Status, row.WaitStrokeWidth)
		}
		if !sameCSSColour(row.WaitStroke, row.WorkFill) {
			t.Errorf("status %q outlines its wait in %s but paints its work %s; the outline is the "+
				"same accent at full strength or it is a second hue",
				row.Status, row.WaitStroke, row.WorkFill)
		}
	}

	// 3. The statuses a reader must be able to tell apart, do differ. Without this
	//    the two assertions above are satisfied by painting everything one colour.
	distinct := map[string]string{}
	for _, status := range []string{"pending", "claimed", "completed", "blocked", "cancelled"} {
		fill := byStatus[status].WorkFill
		if previous, taken := distinct[fill]; taken {
			t.Errorf("statuses %q and %q both paint %s; they are the five a reader has to tell apart",
				previous, status, fill)
		}
		distinct[fill] = status
	}

	// 4. Grouped statuses share their group's colour, exactly as the calendar does.
	for _, pair := range [][2]string{
		{"completed", "completed-with-issues"},
		{"blocked", "pending-answers"},
		{"blocked", "blocked-archive-collision"},
		{"blocked", "blocked-dependency-cycle"},
		{"blocked", "failed"},
	} {
		if byStatus[pair[0]].WorkFill != byStatus[pair[1]].WorkFill {
			t.Errorf("%q paints %s and %q paints %s; they share a group on every other view",
				pair[0], byStatus[pair[0]].WorkFill, pair[1], byStatus[pair[1]].WorkFill)
		}
	}

	// 5. The typo. An exact-match rule sends an unrecognized status to the
	//    unrecognized colour; a prefix match would colour it as real blocked work
	//    and hide a broken REQ in plain sight. It lands on the same colour as
	//    blocked BY DESIGN — what is asserted is that it got there through the
	//    unrecognized class, not through a prefix match on its name.
	typo := byStatus["blockd-dependency-cycle"]
	if !typo.Unrecognized {
		t.Fatal("the probe's typo row lost its is-status-unrecognized class; the assertion below " +
			"would then prove nothing")
	}
	if typo.WorkFill != byStatus["blocked"].WorkFill {
		t.Errorf("an unrecognized status painted %s; it must take the same accent the board gives "+
			"every unrecognized status, which is %s", typo.WorkFill, byStatus["blocked"].WorkFill)
	}
}

// The probe above builds its own rows. This pins that fixture to the shipped
// renderer, so the probe cannot keep passing against markup the board stopped
// producing — the REQ-305 lesson: a probe that cannot hold its call site is
// testing a copy.
func TestTimelineStatusRowMarkupMatchesTheProbe(t *testing.T) {
	rendererBytes, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}
	rendererSource := string(rendererBytes)
	for _, required := range []string{
		`"data-status": request.status || ""`,
		`request.statusUnrecognized ? " is-status-unrecognized" : ""`,
		`"timeline-row"`,
		`timeline-segment timeline-segment-wait`,
		`timeline-segment timeline-segment-work`,
	} {
		if !strings.Contains(rendererSource, required) {
			t.Errorf("web/board-timeline.js no longer contains %q, so the status-colour browser "+
				"probe is measuring markup the renderer does not produce", required)
		}
	}
}

func mustJSONList(t *testing.T, values []string) string {
	t.Helper()
	encoded, encodeError := json.Marshal(values)
	if encodeError != nil {
		t.Fatalf("encode probe status list: %v", encodeError)
	}
	return string(encoded)
}

// sameCSSColour compares two computed colour strings. Chromium normalises a
// custom property's value differently depending on where it is read — a bar's
// `fill` comes back as `rgb(r, g, b)` while the property itself can come back as
// the authored hex — so the comparison is on the parsed channels, never on the
// text.
func sameCSSColour(left string, right string) bool {
	return normalizeCSSColour(left) == normalizeCSSColour(right)
}

func normalizeCSSColour(colour string) string {
	trimmed := strings.ToLower(strings.TrimSpace(colour))
	if strings.HasPrefix(trimmed, "#") {
		hex := trimmed[1:]
		if len(hex) == 3 {
			hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
		}
		if len(hex) != 6 {
			return trimmed
		}
		var red, green, blue int
		for index := 0; index < 6; index += 2 {
			value := 0
			for _, digit := range hex[index : index+2] {
				value *= 16
				switch {
				case digit >= '0' && digit <= '9':
					value += int(digit - '0')
				case digit >= 'a' && digit <= 'f':
					value += int(digit-'a') + 10
				default:
					return trimmed
				}
			}
			switch index {
			case 0:
				red = value
			case 2:
				green = value
			case 4:
				blue = value
			}
		}
		return formatRGB(red, green, blue)
	}
	if strings.HasPrefix(trimmed, "rgb(") || strings.HasPrefix(trimmed, "rgba(") {
		inner := trimmed[strings.Index(trimmed, "(")+1 : strings.LastIndex(trimmed, ")")]
		parts := strings.Split(inner, ",")
		if len(parts) < 3 {
			return trimmed
		}
		channels := make([]int, 3)
		for index := 0; index < 3; index++ {
			value := 0
			for _, digit := range strings.TrimSpace(parts[index]) {
				if digit < '0' || digit > '9' {
					return trimmed
				}
				value = value*10 + int(digit-'0')
			}
			channels[index] = value
		}
		return formatRGB(channels[0], channels[1], channels[2])
	}
	return trimmed
}

func formatRGB(red int, green int, blue int) string {
	digits := "0123456789abcdef"
	out := make([]byte, 0, 6)
	for _, channel := range []int{red, green, blue} {
		out = append(out, digits[(channel>>4)&0x0f], digits[channel&0x0f])
	}
	return string(out)
}

// ---- colour maths ----------------------------------------------------------
//
// WCAG relative luminance and contrast, plus the alpha composite the fills
// actually undergo. Written here rather than eyeballed because "can a reader
// tell these two apart" is the question this REQ's encoding turns on, and the
// previous version of this file answered it by asserting that one number was
// smaller than another.

func relativeLuminanceOfCSSColour(colour string) (float64, bool) {
	red, green, blue, alpha, parsed := parseCSSColourChannels(colour)
	if !parsed || alpha == 0 {
		return 0, false
	}
	return 0.2126*srgbToLinear(red) + 0.7152*srgbToLinear(green) + 0.0722*srgbToLinear(blue), true
}

func srgbToLinear(channel float64) float64 {
	normalized := channel / 255
	if normalized <= 0.04045 {
		return normalized / 12.92
	}
	return math.Pow((normalized+0.055)/1.055, 2.4)
}

// compositeLuminance is the luminance of a fill drawn at `alpha` over a
// background. Compositing happens in sRGB space, so this blends the linearised
// luminances rather than the channels — close enough for a floor check and
// exact for the alpha=1 case.
func compositeLuminance(fillLuminance float64, backgroundLuminance float64, alpha float64) float64 {
	return fillLuminance*alpha + backgroundLuminance*(1-alpha)
}

func contrastRatio(leftLuminance float64, rightLuminance float64) float64 {
	lighter := math.Max(leftLuminance, rightLuminance)
	darker := math.Min(leftLuminance, rightLuminance)
	return (lighter + 0.05) / (darker + 0.05)
}

func fullyTransparentCSSColour(colour string) bool {
	_, _, _, alpha, parsed := parseCSSColourChannels(colour)
	return parsed && alpha == 0
}

func parseCSSColourChannels(colour string) (float64, float64, float64, float64, bool) {
	trimmed := strings.ToLower(strings.TrimSpace(colour))
	if strings.HasPrefix(trimmed, "#") {
		normalized := normalizeCSSColour(trimmed)
		if len(normalized) != 6 {
			return 0, 0, 0, 0, false
		}
		return hexPairToFloat(normalized[0:2]), hexPairToFloat(normalized[2:4]), hexPairToFloat(normalized[4:6]), 1, true
	}
	if !strings.HasPrefix(trimmed, "rgb") {
		return 0, 0, 0, 0, false
	}
	open := strings.Index(trimmed, "(")
	close := strings.LastIndex(trimmed, ")")
	if open < 0 || close < open {
		return 0, 0, 0, 0, false
	}
	parts := strings.Split(strings.ReplaceAll(trimmed[open+1:close], "/", ","), ",")
	if len(parts) < 3 {
		return 0, 0, 0, 0, false
	}
	channels := [4]float64{0, 0, 0, 1}
	for index := 0; index < len(parts) && index < 4; index++ {
		value, parseError := strconv.ParseFloat(strings.TrimSpace(parts[index]), 64)
		if parseError != nil {
			return 0, 0, 0, 0, false
		}
		channels[index] = value
	}
	return channels[0], channels[1], channels[2], channels[3], true
}

func hexPairToFloat(pair string) float64 {
	value, parseError := strconv.ParseInt(pair, 16, 32)
	if parseError != nil {
		return 0
	}
	return float64(value)
}

// The label measurement, in an engine — and the assertion REQ-292's lesson asks
// for: that the number RESPONDS to the rendered face rather than being a model
// of one. A width model that returns the same value for every face is
// undetectable by any test that only checks the value looks plausible, which is
// how that defect shipped; so this measures the real face, then measures again
// with a deliberately proportional one and requires the code to REFUSE.
func TestBrowserBehaviorTimelineLabelAdvanceTracksTheFace(t *testing.T) {
	indexHtml := generateLiveSite(t)
	styleBlock := sliceGeneratedStyleBlock(t, indexHtml)
	measureSource := sliceBalancedBlockAfter(t, indexHtml, "function timelineMeasureLabelAdvance(")
	budgetSource := sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelCharacterBudget(")

	pageHTML := `<!doctype html><html><head><meta charset="utf-8"><style>` + styleBlock + `
/* The proportional override. Applied to a second SVG only, so the first still
   measures the shipped face. */
.proportional-face .timeline-row-label { font-family: Georgia, "Times New Roman", serif; }
</style></head><body>
<svg id="shipped-face" class="timeline-rows-svg" width="400" height="60"></svg>
<svg id="wide-face" class="timeline-rows-svg proportional-face" width="400" height="60"></svg>
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
var TIMELINE_SVG_NS = "http://www.w3.org/2000/svg";
` + measureSource + `
` + budgetSource + `
(function () {
  var shipped = timelineMeasureLabelAdvance(document.getElementById("shipped-face"));
  var proportional = timelineMeasureLabelAdvance(document.getElementById("wide-face"));
  // What the engine says the two faces actually are, so a failure below can be
  // read without guessing whether the override applied at all.
  function sampleWidth(svgId, sample) {
    var svg = document.getElementById(svgId);
    var node = document.createElementNS(TIMELINE_SVG_NS, "text");
    node.setAttribute("class", "timeline-row-label");
    node.textContent = sample;
    svg.appendChild(node);
    var width = node.getComputedTextLength();
    svg.removeChild(node);
    return width;
  }
  document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
    shippedAdvance: shipped,
    proportionalAdvance: proportional,
    shippedNarrow: sampleWidth("shipped-face", "iiiiiiiiii"),
    shippedWide: sampleWidth("shipped-face", "MMMMMMMMMM"),
    proportionalNarrow: sampleWidth("wide-face", "iiiiiiiiii"),
    proportionalWide: sampleWidth("wide-face", "MMMMMMMMMM"),
    budgetFromShipped: timelineLabelCharacterBudget(shipped, 172)
  });
})();
</script>
</body></html>`

	probeOutput := runBrowserBehaviorProbe(t, "timeline label advance", pageHTML)
	var advanceResult struct {
		ShippedAdvance      float64 `json:"shippedAdvance"`
		ProportionalAdvance float64 `json:"proportionalAdvance"`
		ShippedNarrow       float64 `json:"shippedNarrow"`
		ShippedWide         float64 `json:"shippedWide"`
		ProportionalNarrow  float64 `json:"proportionalNarrow"`
		ProportionalWide    float64 `json:"proportionalWide"`
		BudgetFromShipped   int     `json:"budgetFromShipped"`
	}
	if decodeError := json.Unmarshal(probeOutput, &advanceResult); decodeError != nil {
		t.Fatalf("decode timeline label advance probe: %v (output %q)", decodeError, probeOutput)
	}

	// The fixture has to be a real test of the thing. If the override did not
	// take, both faces are the shipped one and the refusal below proves nothing.
	if advanceResult.ProportionalNarrow == advanceResult.ProportionalWide {
		t.Fatalf("the proportional override did not take — 'iiiiiiiiii' and 'MMMMMMMMMM' both "+
			"measured %.4f, so this test cannot tell a measured advance from a modelled one",
			advanceResult.ProportionalNarrow)
	}
	if advanceResult.ShippedNarrow != advanceResult.ShippedWide {
		t.Fatalf("the shipped label face is not monospace: 'iiiiiiiiii' measured %.4f and "+
			"'MMMMMMMMMM' %.4f. One advance cannot describe it, and the label budget is built on "+
			"the assumption that it can",
			advanceResult.ShippedNarrow, advanceResult.ShippedWide)
	}
	if !(advanceResult.ShippedAdvance > 0) {
		t.Fatalf("the shipped face measured an advance of %.4f; a zero advance disables the title "+
			"column entirely", advanceResult.ShippedAdvance)
	}
	// The refusal. A proportional face must produce no advance at all rather than
	// an average that would mis-cut every title differently.
	if advanceResult.ProportionalAdvance != 0 {
		t.Errorf("a proportional face produced an advance of %.4f; it must refuse, because one "+
			"number cannot describe a face whose glyphs differ in width",
			advanceResult.ProportionalAdvance)
	}
	// And the measured number has to be in the range a 10px face can occupy —
	// wide enough to be real, narrow enough that the column is not one word.
	if advanceResult.ShippedAdvance < 4 || advanceResult.ShippedAdvance > 9 {
		t.Errorf("the shipped 10px label face measured %.4f px per character, outside the 4–9 px "+
			"a 10px monospace face can plausibly occupy — record the browser and build if this "+
			"is a real face difference", advanceResult.ShippedAdvance)
	}
	if advanceResult.BudgetFromShipped < 20 {
		t.Errorf("the shipped face fits only %d characters in a 172px column; the title column "+
			"stops being worth its width below roughly 20", advanceResult.BudgetFromShipped)
	}
	t.Logf("shipped label face: %.4f px/char, %d characters in 172px (record the browser build "+
		"beside this number if it is ever pinned)",
		advanceResult.ShippedAdvance, advanceResult.BudgetFromShipped)
}

// The row's tooltip. A native <title> is the whole feature, so what needs pinning
// is that the renderer still writes one and that it carries the description
// rather than a bare id.
// The hint paragraph says "Click a row for its full detail." Before this it was
// true only of a perfectly still press: the first pointermove armed the pan,
// which re-rendered the rows, and the click the browser then synthesized had no
// surviving [data-detail-kind] ancestor to find. This probe drives the REAL
// board with real pointer events, because that failure lives entirely in the
// interaction between the engine's click synthesis and a DOM the handler
// rebuilt — no sliced function can show it.
func TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function pointerAt(type, node, clientX, clientY) {
  node.dispatchEvent(new PointerEvent(type, {
    bubbles: true, cancelable: true, composed: true,
    clientX: clientX, clientY: clientY, button: 0,
    buttons: type === "pointerup" ? 0 : 1, pointerId: 1, pointerType: "mouse"
  }));
}
function nodeAt(x, y, fallbackNode) {
  return document.elementFromPoint(x, y) || fallbackNode;
}
function isGrabbing() {
  return document.querySelector("#view-timeline .timeline-scroll")
    .classList.contains("is-panning");
}
// The browser synthesizes a click at the nearest common inclusive ancestor of
// the pointerdown and pointerup targets. Dispatching against elementFromPoint at
// release is the closest a script comes to that, and it is what makes a rebuilt
// row observable: the node under the pointer is no longer the node pressed.
//
// Returns what was true DURING the press, because the release clears it.
function pressDragRelease(fallbackNode, startX, y, moveOffsets) {
  var pressTarget = nodeAt(startX, y, fallbackNode);
  pointerAt("pointerdown", pressTarget, startX, y);
  var grabbingAfterPress = isGrabbing();
  var lastX = startX;
  moveOffsets.forEach(function (offset) {
    lastX = startX + offset;
    pointerAt("pointermove", nodeAt(lastX, y, fallbackNode), lastX, y);
  });
  var grabbingDuringDrag = isGrabbing();
  var releaseTarget = nodeAt(lastX, y, fallbackNode);
  pointerAt("pointerup", releaseTarget, lastX, y);
  releaseTarget.dispatchEvent(new MouseEvent("click", {
    bubbles: true, cancelable: true, composed: true, clientX: lastX, clientY: y
  }));
  return { grabbingAfterPress: grabbingAfterPress, grabbingDuringDrag: grabbingDuringDrag };
}
function drawerIsOpen() {
  var drawer = document.getElementById("detail-drawer");
  return !!drawer && !drawer.hidden && !drawer.classList.contains("is-hidden");
}
function closeDrawerIfOpen() {
  var close = document.getElementById("detail-close");
  if (close) { close.click(); }
}
function windowReadout() {
  var node = document.getElementById("timeline-range-readout");
  return node ? node.textContent : "";
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var trials = {};
      // The press point, PROVEN rather than assumed. Earlier this computed a row
      // rect and pressed at its centre; because the trials open and close the
      // detail drawer the document height changes between them, scrollIntoView
      // landed the chart at a different offset each time, and by the fourth trial
      // the point was over the axis strip. Two trials were then measuring a press
      // that reached no handler and reporting it as a result — so the setup is
      // now checked and reported, and the Go side fails on it.
      function pressPointOverFirstRow() {
        var host = document.querySelector("#view-timeline .timeline-scroll");
        var hostBox = host.getBoundingClientRect();
        var rows = document.querySelectorAll("#view-timeline .timeline-row");
        for (var index = 0; index < rows.length; index++) {
          var rowBox = rows[index].getBoundingClientRect();
          var y = Math.round(rowBox.top + rowBox.height / 2);
          // Inside the host's visible band, and far enough from the right edge
          // that a 120px drag stays on the plot.
          if (y <= hostBox.top + 2 || y >= hostBox.bottom - 2) {
            continue;
          }
          var x = Math.round(Math.min(rowBox.left + rowBox.width / 2, hostBox.right - 200));
          var landedOn = document.elementFromPoint(x, y);
          if (landedOn && host.contains(landedOn)) {
            return { node: rows[index], x: x, y: y, ok: true };
          }
        }
        return { node: null, x: 0, y: 0, ok: false };
      }
      function trial(name, moveOffsets) {
        closeDrawerIfOpen();
        // Every trial starts from the SAME window. At Fit all the window sits on
        // both bounds and a pan clamps to the window it started in, so a drag
        // that panned would be indistinguishable from one that did not; two zoom
        // steps put it off the bounds. Resetting each time is what lets two
        // trials' end windows be compared to each other.
        document.getElementById("timeline-zoom-fit").click();
        document.getElementById("timeline-zoom-in").click();
        document.getElementById("timeline-zoom-in").click();
        // elementFromPoint is viewport-relative and returns null below the fold.
        // The chart sits under a warnings banner and an anomalies board.
        document.querySelector("#view-timeline .timeline-chart")
          .scrollIntoView({ block: "start" });
        var press = pressPointOverFirstRow();
        var before = windowReadout();
        var during = press.ok
          ? pressDragRelease(press.node, press.x, press.y, moveOffsets)
          : { grabbingAfterPress: false, grabbingDuringDrag: false };
        trials[name] = {
          pressLandedOnARow: press.ok,
          drawerOpen: drawerIsOpen(),
          windowMoved: windowReadout() !== before,
          windowAfter: windowReadout(),
          grabbingAfterPress: during.grabbingAfterPress,
          grabbingDuringDrag: during.grabbingDuringDrag
        };
      }
      trial("stillPress", []);
      trial("belowThreshold", [2]);
      trial("aboveThreshold", [120]);
      // The same 120 pixels, delivered in two moves with the first one tripping
      // the threshold. Anchoring the shift at the trip point instead of the press
      // point leaves this one short by the threshold, forever.
      trial("aboveThresholdInTwoMoves", [6, 120]);
      // Out past the threshold and back to within it. The reader has visibly
      // panned; releasing there must not be read as a click.
      trial("draggedOutAndBack", [20, 1]);
      closeDrawerIfOpen();
      document.getElementById("` + browserProbeResultElementId + `").textContent =
        JSON.stringify(trials);
      document.title = "READY";
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	// A viewport wide enough for the chart. The default headless window is 800x600,
	// where the plot is a few dozen pixels wide and a 120-pixel drag runs off it.
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline pan threshold", siteDirectory,
		pageHTML, "--window-size=1400,900")

	type pressTrial struct {
		PressLandedOnARow  bool   `json:"pressLandedOnARow"`
		DrawerOpen         bool   `json:"drawerOpen"`
		WindowMoved        bool   `json:"windowMoved"`
		WindowAfter        string `json:"windowAfter"`
		GrabbingAfterPress bool   `json:"grabbingAfterPress"`
		GrabbingDuringDrag bool   `json:"grabbingDuringDrag"`
	}
	var pressResult struct {
		StillPress               pressTrial `json:"stillPress"`
		BelowThreshold           pressTrial `json:"belowThreshold"`
		AboveThreshold           pressTrial `json:"aboveThreshold"`
		AboveThresholdInTwoMoves pressTrial `json:"aboveThresholdInTwoMoves"`
		DraggedOutAndBack        pressTrial `json:"draggedOutAndBack"`
	}
	if decodeError := json.Unmarshal(probeOutput, &pressResult); decodeError != nil {
		t.Fatalf("decode timeline pan threshold behavior: %v (output %q)", decodeError, probeOutput)
	}

	// THE SETUP, BEFORE ANY MEASUREMENT. A press that never reached the scroll
	// host produces "no pan, no drawer" — which satisfies several assertions below
	// while testing nothing at all. This caught exactly that: the trials open and
	// close the drawer, the document height changes, and the later presses were
	// landing on the axis strip above the rows.
	for _, trial := range []struct {
		name   string
		result pressTrial
	}{
		{"stillPress", pressResult.StillPress},
		{"belowThreshold", pressResult.BelowThreshold},
		{"aboveThreshold", pressResult.AboveThreshold},
		{"aboveThresholdInTwoMoves", pressResult.AboveThresholdInTwoMoves},
		{"draggedOutAndBack", pressResult.DraggedOutAndBack},
	} {
		if !trial.result.PressLandedOnARow {
			t.Fatalf("the %s trial could not find a press point inside the scroll host over a "+
				"visible row, so it pressed nothing; every assertion below would then be "+
				"measuring a press that reached no handler", trial.name)
		}
	}

	// THE PAIR. A still press already worked before the threshold existed; the
	// two-pixel one is the regression, and keeping both is what proves the fix is
	// a threshold and not "clicks now always win".
	if !pressResult.StillPress.DrawerOpen {
		t.Error("a perfectly still press did not open the detail drawer; that worked before " +
			"the threshold existed and the threshold must not have cost it")
	}
	if !pressResult.BelowThreshold.DrawerOpen {
		t.Error("a two-pixel press did not open the detail drawer — the press panned and " +
			"re-rendered the rows, so the click the engine synthesized found no surviving " +
			"[data-detail-kind] to open. This is the defect the threshold exists to fix")
	}
	if pressResult.BelowThreshold.WindowMoved {
		t.Error("a two-pixel press moved the time window; nobody asked a click to scroll the " +
			"chart, and a hand tremor is not a pan")
	}
	// Read DURING the press, not after: the release clears the class, so an
	// assertion taken afterwards passes whatever the cursor did.
	if pressResult.BelowThreshold.GrabbingDuringDrag {
		t.Error("a two-pixel press showed the grab cursor; the cursor is a claim about what " +
			"the press is doing, and it was claiming a drag during a click")
	}
	if pressResult.StillPress.GrabbingAfterPress {
		t.Error("the grab cursor appeared on pointerdown alone; a press that has not moved is " +
			"not a drag yet, and saying so is how the reader learns the threshold exists")
	}

	// The other side. A deliberate drag has to pan, and must NOT then open the
	// drawer on release — "clicks always win" would pass every assertion above.
	if !pressResult.AboveThreshold.WindowMoved {
		t.Error("a 120-pixel drag did not move the time window; the threshold is supposed to " +
			"delay the pan, not prevent it")
	}
	if pressResult.AboveThreshold.DrawerOpen {
		t.Error("a 120-pixel drag opened the detail drawer on release; a drag is not a click, " +
			"and every pan across a row would pop the drawer")
	}
	if !pressResult.AboveThreshold.GrabbingDuringDrag {
		t.Error("a 120-pixel drag never showed the grab cursor; the threshold delays the " +
			"cursor change, it does not cancel it")
	}

	// LATCHING. A drag that wanders back inside the threshold is still a drag: the
	// reader has already seen the chart move, and un-engaging would both flicker
	// the cursor and turn the release into a click that opens a drawer they did
	// not ask for.
	if pressResult.DraggedOutAndBack.DrawerOpen {
		t.Error("a drag that went 20px out and came back to 1px opened the detail drawer; " +
			"engagement has to latch, or a wandering drag ends as a click")
	}
	if !pressResult.DraggedOutAndBack.GrabbingDuringDrag {
		t.Error("a drag that went 20px out and came back to 1px dropped the grab cursor; " +
			"engagement has to latch, or a slow drag flickers around the threshold")
	}

	// CONTINUITY FROM THE PRESS POINT. The same 120 pixels in one move and in two
	// have to land on the same window. Re-anchoring at the moment the threshold
	// trips leaves the two-move drag a threshold's worth short — a lag the reader
	// carries for the rest of the drag, and the reason this is measured rather
	// than asserted in a comment.
	if pressResult.AboveThreshold.WindowAfter != pressResult.AboveThresholdInTwoMoves.WindowAfter {
		t.Errorf("120 pixels in one move landed on %q and in two moves on %q; the shift has to "+
			"be measured from the press point, or the chart settles a threshold behind the "+
			"pointer",
			pressResult.AboveThreshold.WindowAfter,
			pressResult.AboveThresholdInTwoMoves.WindowAfter)
	}
}

// The other half of the same handler. Every pointermove used to run a full
// renderAll, and every renderAll read scrollHost.clientWidth once per xOfEpoch
// call — several times per row, per frame of a drag, each one a forced
// synchronous layout. Both are counted here rather than asserted, because
// "should be faster" is not a check and a comment claiming a saving is how the
// last one drifted.
func TestBrowserBehaviorTimelineDragRendersOncePerFrame(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function pointerAt(type, node, clientX, clientY) {
  node.dispatchEvent(new PointerEvent(type, {
    bubbles: true, cancelable: true, composed: true,
    clientX: clientX, clientY: clientY, button: 0,
    buttons: type === "pointerup" ? 0 : 1, pointerId: 1, pointerType: "mouse"
  }));
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var scrollHost = document.querySelector("#view-timeline .timeline-scroll");
      var axisSvg = document.querySelector("#view-timeline .timeline-axis-svg");

      // A render begins by clearing the axis SVG, so one childList record
      // carrying removed nodes is one render. Counting the renderer's own
      // observable side effect beats trusting a counter the renderer would have
      // to be modified to keep.
      var renderCount = 0;
      var renderObserver = new MutationObserver(function (records) {
        records.forEach(function (record) {
          if (record.removedNodes.length > 0) { renderCount++; }
        });
      });
      function countRendersNow() {
        renderObserver.takeRecords().forEach(function (record) {
          if (record.removedNodes.length > 0) { renderCount++; }
        });
        return renderCount;
      }
      renderObserver.observe(axisSvg, { childList: true });

      // clientWidth is the forced-layout read. An own property shadowing the
      // prototype getter counts every call and hands back the real value, so the
      // render under measurement is the real one.
      var widthReads = 0;
      var clientWidthDescriptor =
        Object.getOwnPropertyDescriptor(Element.prototype, "clientWidth");
      Object.defineProperty(scrollHost, "clientWidth", {
        configurable: true,
        get: function () { widthReads++; return clientWidthDescriptor.get.call(this); }
      });

      // ONE render, provoked by a control rather than a drag, so the count is a
      // per-render figure and not a per-frame one.
      renderCount = 0;
      widthReads = 0;
      document.getElementById("timeline-zoom-in").click();
      var widthReadsPerRender = widthReads;
      var rendersPerControlClick = countRendersNow();

      // Now the drag. Five moves dispatched back to back, all inside one task:
      // with coalescing none of them has rendered yet when this line runs.
      var box = scrollHost.getBoundingClientRect();
      var startX = Math.round(box.left + box.width / 2);
      var y = Math.round(box.top + 20);
      renderCount = 0;
      pointerAt("pointerdown", scrollHost, startX, y);
      for (var step = 1; step <= 5; step++) {
        pointerAt("pointermove", scrollHost, startX - step * 20, y);
      }
      var rendersDuringFiveMoves = countRendersNow();
      // The release flushes whatever the frame still owed, so the drag lands on
      // the window it reached rather than the one the last frame drew.
      pointerAt("pointerup", scrollHost, startX - 100, y);
      var rendersAfterRelease = countRendersNow();

      document.getElementById("` + browserProbeResultElementId + `").textContent =
        JSON.stringify({
          widthReadsPerRender: widthReadsPerRender,
          rendersPerControlClick: rendersPerControlClick,
          rendersDuringFiveMoves: rendersDuringFiveMoves,
          rendersAfterRelease: rendersAfterRelease
        });
      document.title = "READY";
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(string(indexBytes), "</body>", probeScript, 1)
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline drag render cost", siteDirectory,
		pageHTML, "--window-size=1400,900")

	var costResult struct {
		WidthReadsPerRender    int `json:"widthReadsPerRender"`
		RendersPerControlClick int `json:"rendersPerControlClick"`
		RendersDuringFiveMoves int `json:"rendersDuringFiveMoves"`
		RendersAfterRelease    int `json:"rendersAfterRelease"`
	}
	if decodeError := json.Unmarshal(probeOutput, &costResult); decodeError != nil {
		t.Fatalf("decode timeline drag render cost: %v (output %q)", decodeError, probeOutput)
	}

	// The control render has to have HAPPENED, or every count below is a count of
	// nothing — REQ-291's lesson, in this module's own prime.
	if costResult.RendersPerControlClick != 1 {
		t.Fatalf("a zoom click produced %d renders, want exactly 1; the observer is watching "+
			"the wrong node or the click did nothing, and the counts below mean nothing",
			costResult.RendersPerControlClick)
	}

	// ONE forced layout read per render. Before the memo, xOfEpoch called
	// plotWidth() several times per row and plotWidth() read clientWidth every
	// time — a 35-row render made hundreds of them. The ceiling is 2 rather than 1
	// so a second legitimate reader is not a failure; it is nowhere near the
	// hundreds the defect produced.
	if costResult.WidthReadsPerRender > 2 {
		t.Errorf("one render read scrollHost.clientWidth %d times; each is a forced synchronous "+
			"layout, and the plot width can only change between renders — measure it once",
			costResult.WidthReadsPerRender)
	}
	if costResult.WidthReadsPerRender == 0 {
		t.Error("one render read scrollHost.clientWidth zero times; the memo is never being " +
			"invalidated, so a resized container would keep drawing at the old width")
	}

	// Coalescing. Five moves in one task must not be five renders; with the frame
	// still pending they are none yet.
	if costResult.RendersDuringFiveMoves != 0 {
		t.Errorf("five pointermoves in one task produced %d renders before the frame ran, "+
			"want 0; a drag delivers moves faster than the compositor draws them and every "+
			"render but the last is thrown away unseen", costResult.RendersDuringFiveMoves)
	}
	// And the flush: the drag still has to land, or coalescing would be achieved
	// by dropping the work rather than by deferring it.
	if costResult.RendersAfterRelease != 1 {
		t.Errorf("five pointermoves and a release produced %d renders in total, want exactly 1 "+
			"— the deferred frame has to be flushed on release, or the chart settles on the "+
			"window the last drawn frame reached instead of the one the drag ended at",
			costResult.RendersAfterRelease)
	}
}

// The chart draws four kinds of vertical mark and their ORDER OF PROMINENCE is
// the whole design: the now-line is the present and has to win, the queue-end
// rule is a forecast and must not outshout it, the gridlines are a backdrop and
// must sit behind the bars they cross without disappearing. None of that can be
// argued from a colour token — it is a rendered result, so it is measured.
func TestBrowserBehaviorTimelineVerticalRulesRankByProminence(t *testing.T) {
	indexHtml := generateLiveSite(t)
	styleBlock := sliceGeneratedStyleBlock(t, indexHtml)

	pageHTML := `<!doctype html><html><head><meta charset="utf-8"><style>` + styleBlock + `
</style></head><body>
<svg class="timeline-axis-svg" width="400" height="26">
  <line class="timeline-axis-tick" x1="10" y1="20" x2="10" y2="26"></line>
  <line class="timeline-now-line" x1="20" y1="0" x2="20" y2="26"></line>
  <line class="timeline-queue-end-line" x1="30" y1="0" x2="30" y2="26"></line>
</svg>
<svg class="timeline-rows-svg" width="400" height="40">
  <line class="timeline-gridline" x1="10" y1="0" x2="10" y2="40"></line>
  <line class="timeline-queue-end-rule" x1="30" y1="0" x2="30" y2="40"></line>
  <line class="timeline-now-rule" x1="20" y1="0" x2="20" y2="40"></line>
</svg>
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
(function () {
  function measureRule(selector) {
    var style = getComputedStyle(document.querySelector(selector));
    return {
      stroke: style.stroke,
      strokeWidth: parseFloat(style.strokeWidth) || 0,
      // The node's own opacity multiplies whatever alpha the stroke colour
      // already carries; reading one and not the other is how a rule that looks
      // faint measures as full strength.
      opacity: Number(style.opacity),
      dashed: (style.strokeDasharray || "none") !== "none",
      // How much of the line's length is actually inked. A dotted rule and a
      // solid one of the same colour and width are not the same mark, and
      // leaving this out of the measurement is leaving out the channel the
      // design uses hardest.
      dashCoverage: (function () {
        var pattern = (style.strokeDasharray || "none").split(/[\s,]+/)
          .map(parseFloat).filter(function (value) { return isFinite(value) && value >= 0; });
        if (pattern.length === 0) {
          return 1;
        }
        // An odd-length pattern repeats to twice its length, alternating roles.
        if (pattern.length % 2 === 1) {
          pattern = pattern.concat(pattern);
        }
        var inked = 0;
        var total = 0;
        for (var index = 0; index < pattern.length; index++) {
          total += pattern[index];
          if (index % 2 === 0) {
            inked += pattern[index];
          }
        }
        return total > 0 ? inked / total : 1;
      })()
    };
  }
  document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
    gridline: measureRule(".timeline-gridline"),
    axisTick: measureRule(".timeline-axis-tick"),
    nowRule: measureRule(".timeline-now-rule"),
    nowLine: measureRule(".timeline-now-line"),
    queueEndRule: measureRule(".timeline-queue-end-rule"),
    queueEndLine: measureRule(".timeline-queue-end-line"),
    scheme: matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
    surface: getComputedStyle(document.body).backgroundColor
  });
})();
</script>
</body></html>`

	for _, scheme := range []struct {
		name string
		flag string
	}{
		{name: "light", flag: "--blink-settings=preferredColorScheme=1"},
		{name: "dark", flag: "--blink-settings=preferredColorScheme=0"},
	} {
		t.Run(scheme.name, func(t *testing.T) {
			assertTimelineRulePromimence(t, pageHTML, scheme.name, scheme.flag)
		})
	}
}

// timelineRuleVisibilityFloor is the least a vertical mark may stand out from the
// page and still be a mark. It is well under the 3:1 the two halves of a BAR owe
// each other, and deliberately so: a gridline that reached 3:1 would be a second
// chart competing with the first.
const timelineRuleVisibilityFloor = 1.03

// timelineForecastPresenceCeiling is how loud the queue-end mark may be relative
// to the now mark beside it. 0.75 is a margin a palette tweak does not erase;
// it is not a perceptual constant, it is a gap wide enough to be a decision.
const timelineForecastPresenceCeiling = 0.75

func assertTimelineRulePromimence(t *testing.T, pageHTML string, schemeName string, schemeFlag string) {
	t.Helper()
	probeOutput := runBrowserBehaviorProbeWithFlags(t,
		"timeline vertical rules ("+schemeName+")", pageHTML, schemeFlag)

	type measuredRule struct {
		Stroke       string  `json:"stroke"`
		StrokeWidth  float64 `json:"strokeWidth"`
		Opacity      float64 `json:"opacity"`
		Dashed       bool    `json:"dashed"`
		DashCoverage float64 `json:"dashCoverage"`
	}
	var ruleResult struct {
		Gridline     measuredRule `json:"gridline"`
		AxisTick     measuredRule `json:"axisTick"`
		NowRule      measuredRule `json:"nowRule"`
		NowLine      measuredRule `json:"nowLine"`
		QueueEndRule measuredRule `json:"queueEndRule"`
		QueueEndLine measuredRule `json:"queueEndLine"`
		Scheme       string       `json:"scheme"`
		Surface      string       `json:"surface"`
	}
	if decodeError := json.Unmarshal(probeOutput, &ruleResult); decodeError != nil {
		t.Fatalf("decode timeline vertical rule prominence: %v (output %q)", decodeError, probeOutput)
	}
	if ruleResult.Scheme != schemeName {
		t.Fatalf("asked the engine for the %s palette and it resolved %s; the flag did not "+
			"take, so one palette would be measured twice and reported as two",
			schemeName, ruleResult.Scheme)
	}

	surfaceLuminance, surfaceReadable := relativeLuminanceOfCSSColour(ruleResult.Surface)
	if !surfaceReadable {
		t.Fatalf("could not read the page surface colour %q", ruleResult.Surface)
	}

	// Presence: how far the composited stroke lands from the page, weighted by how
	// much of the page it actually inks. Three channels — contrast, width, dash
	// duty — because a 1px dotted line at 0.08 alpha and a 1.5px solid one at full
	// strength are not the same mark however close their hues are. Dropping any
	// one of the three is how a rule that plainly looks quieter measures louder.
	presenceOf := func(ruleName string, rule measuredRule) float64 {
		t.Helper()
		strokeLuminance, strokeReadable := relativeLuminanceOfCSSColour(rule.Stroke)
		if !strokeReadable {
			t.Fatalf("could not read the %s stroke colour %q", ruleName, rule.Stroke)
		}
		_, _, _, strokeAlpha, channelsReadable := parseCSSColourChannels(rule.Stroke)
		if !channelsReadable {
			t.Fatalf("could not read the %s stroke alpha from %q", ruleName, rule.Stroke)
		}
		composited := compositeLuminance(strokeLuminance, surfaceLuminance, strokeAlpha*rule.Opacity)
		return (contrastRatio(composited, surfaceLuminance) - 1) * rule.StrokeWidth * rule.DashCoverage
	}

	gridlinePresence := presenceOf("gridline", ruleResult.Gridline)
	axisTickPresence := presenceOf("axis tick", ruleResult.AxisTick)
	nowRulePresence := presenceOf("now rule", ruleResult.NowRule)
	nowLinePresence := presenceOf("now line", ruleResult.NowLine)
	queueEndRulePresence := presenceOf("queue-end rule", ruleResult.QueueEndRule)
	queueEndLinePresence := presenceOf("queue-end line", ruleResult.QueueEndLine)

	gridlineContrast := gridlinePresence/(ruleResult.Gridline.StrokeWidth*ruleResult.Gridline.DashCoverage) + 1
	if gridlineContrast < timelineRuleVisibilityFloor {
		t.Errorf("[%s] the gridline sits %.4f:1 from the page, under the %.2f:1 floor; "+
			"a reference the reader cannot see is a reference they do not have",
			schemeName, gridlineContrast, timelineRuleVisibilityFloor)
	}
	// The axis strip is the measured edge and the plot is the backdrop. A gridline
	// at or above the tick's own presence inverts that and the eye reads the plot
	// as ruled paper with a chart faintly on top.
	if gridlinePresence >= axisTickPresence {
		t.Errorf("[%s] the gridline's presence is %.4f and the axis tick's is %.4f; the "+
			"gridline has to stay behind the tick it descends from",
			schemeName, gridlinePresence, axisTickPresence)
	}

	queueEndContrast := queueEndRulePresence/(ruleResult.QueueEndRule.StrokeWidth*ruleResult.QueueEndRule.DashCoverage) + 1
	if queueEndContrast < timelineRuleVisibilityFloor {
		t.Errorf("[%s] the queue-end rule sits %.4f:1 from the page, under the %.2f:1 floor; "+
			"the forecast paragraph names an instant and this is the only mark for it",
			schemeName, queueEndContrast, timelineRuleVisibilityFloor)
	}
	// THE ORDER THE REQ ASKED FOR. Both halves, because the rule and the label's
	// line are separately styled and only one of them was ever going to be
	// remembered on the next edit.
	// A CLEAR margin, not a tie broken by the fourth decimal. "Quieter by 0.3%" is
	// the same picture as "equally loud", and a bare `>=` would call it a pass.
	if queueEndRulePresence > nowRulePresence*timelineForecastPresenceCeiling {
		t.Errorf("[%s] the queue-end rule's presence is %.4f against the now-rule's %.4f, "+
			"over the %.0f%% ceiling; a forecast that reads as loud as the present moves the "+
			"reader's eye to the wrong mark",
			schemeName, queueEndRulePresence, nowRulePresence,
			timelineForecastPresenceCeiling*100)
	}
	if queueEndLinePresence > nowLinePresence*timelineForecastPresenceCeiling {
		t.Errorf("[%s] in the axis strip the queue-end line's presence is %.4f against the "+
			"now-line's %.4f, over the %.0f%% ceiling; same rule, other half of the chart",
			schemeName, queueEndLinePresence, nowLinePresence,
			timelineForecastPresenceCeiling*100)
	}
	// Presence alone can be matched by two marks of different meaning, so the
	// forecast also carries a shape difference the now-line does not.
	if !ruleResult.QueueEndRule.Dashed || !ruleResult.QueueEndLine.Dashed {
		t.Errorf("[%s] the queue-end mark is not dashed; it is a projection and has to be "+
			"legible as one at a glance, not only by being fainter", schemeName)
	}
	if ruleResult.NowLine.Dashed {
		t.Errorf("[%s] the axis now-line became dashed; it and the queue-end line then differ "+
			"by hue alone", schemeName)
	}
}

// The probe above draws its own three-line SVGs. This is what keeps them the
// lines the renderer actually emits — REQ-305's lesson: a probe that cannot hold
// its call site tests a copy.
func TestTimelineVerticalRuleMarkupMatchesTheProbe(t *testing.T) {
	rendererBytes, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}
	rendererSource := string(rendererBytes)
	for ruleClass, drawnBy := range map[string]string{
		"timeline-gridline":       "drawGridlines",
		"timeline-queue-end-rule": "drawQueueEndRule",
		"timeline-queue-end-line": "renderAxis",
		"timeline-now-rule":       "drawNowRule",
		"timeline-now-line":       "renderAxis",
	} {
		if !strings.Contains(rendererSource, `class: "`+ruleClass+`"`) {
			t.Errorf("the renderer emits no %q line any more (expected from %s); the "+
				"prominence probe is measuring a class nothing draws", ruleClass, drawnBy)
		}
	}
	if !strings.Contains(rendererSource, "drawGridlines();") {
		t.Error("renderVisibleRows no longer calls drawGridlines, so the plot has no vertical " +
			"reference and the prominence probe measures a rule that is never drawn")
	}
	if !strings.Contains(rendererSource, "drawQueueEndRule();") {
		t.Error("renderVisibleRows no longer calls drawQueueEndRule, so the forecast's " +
			"queue-empty instant is named in prose and marked nowhere")
	}
}

func TestTimelineRowTooltipMarkupMatchesTheProbe(t *testing.T) {
	rendererBytes, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}
	rendererSource := string(rendererBytes)
	if !strings.Contains(rendererSource, `makeTimelineSvgNode(rowGroup, "title", {}, timelineRowDescription(row, request));`) {
		t.Error("the row group no longer carries a native <title> built from timelineRowDescription; " +
			"the tooltip at the pointer was the point, and the foot readout is 700px away")
	}
	if !strings.Contains(rendererSource, "timelineRowLabelText(row.id, request.title, labelCharacterBudget)") {
		t.Error("the row label no longer runs through timelineRowLabelText, so the truncation " +
			"assertions in generate_test.go are testing a function nothing calls")
	}
}

// plotEdgeText renders an optional measured edge for a failure message. A *int
// formatted with %v prints its ADDRESS, which is how a real failure came out
// reading "leftmost 0x1b8ca51929b0" — true, useless, and impossible to argue with.
func plotEdgeText(edge *int) string {
	if edge == nil {
		return "none measured"
	}
	return strconv.Itoa(*edge)
}

// Opening the detail drawer narrows the plot, and before this probe nothing
// re-measured it.
//
// WHY A BROWSER, AND WHY THE WHOLE BOARD. The defect is a layout interaction: the
// drawer takes a 620px grid column, the scroll host loses 630px of width, and
// every bar keeps the x it was given against the old plot. Nothing about the
// renderer is wrong in isolation — a sliced function measures whatever width the
// fixture hands it — so the only thing that can show this is a real engine laying
// out the real page, measuring the real bars, and comparing what is inside the
// host's box before and after a real click on a row.
//
// The measurement is getBoundingClientRect() intersection, not node count: the
// bars were all still in the DOM when the chart looked blank. Fifty-five segments,
// zero of them inside the host.
func TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
// Waiting on a CONDITION, never on a duration. Fixed setTimeouts passed alone and
// failed inside the full suite, where a dozen probes share the machine: the two
// 300ms waits after a row click were enough for an idle run and not for a loaded
// one, and the probe then measured a layout that had not settled and reported the
// defect it exists to catch. Bounded by a deadline so a genuine failure still
// fails rather than hanging.
function settleUntil(predicate, thenDo) {
  var attemptsLeft = 60;
  (function attempt() {
    if (predicate() || attemptsLeft-- <= 0) {
      // One more tick after the predicate holds, so a render the change SCHEDULED
      // has run before anything is measured. setTimeout rather than
      // requestAnimationFrame throughout: headless --dump-dom has no compositor
      // driving frames, so an rAF-based poll never resolves and the runner dumps
      // an empty result node.
      setTimeout(thenDo, 50);
      return;
    }
    setTimeout(attempt, 50);
  })();
}
function plotHost() {
  return document.querySelector("#view-timeline .timeline-scroll");
}
// What the reader can actually see: segments whose box overlaps the host's box.
// Counting nodes would have passed straight through the defect.
function plotSnapshot(label) {
  var host = plotHost();
  var hostBox = host.getBoundingClientRect();
  var segments = [].slice.call(host.querySelectorAll("rect.timeline-segment"));
  var inside = 0;
  var rightmost = -Infinity;
  var leftmost = Infinity;
  segments.forEach(function (segment) {
    var box = segment.getBoundingClientRect();
    if (box.right > hostBox.left && box.left < hostBox.right) {
      inside++;
    }
    rightmost = Math.max(rightmost, box.right);
    leftmost = Math.min(leftmost, box.left);
  });
  return {
    label: label,
    href: location.href,
    hostWidth: Math.round(hostBox.width),
    hostRight: Math.round(hostBox.right),
    segments: segments.length,
    inside: inside,
    rightmost: isFinite(rightmost) ? Math.round(rightmost) : null,
    leftmost: isFinite(leftmost) ? Math.round(leftmost) : null
  };
}
function drawerIsOpen() {
  var drawer = document.getElementById("detail-drawer");
  return !!drawer && !drawer.hidden && !drawer.classList.contains("is-hidden");
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var probe = {};
      probe.before = plotSnapshot("before");
      probe.drawerOpenBefore = drawerIsOpen();
      // A row, clicked the way a reader clicks one. The delegated handler reads
      // [data-detail-kind] off the row group.
      var firstRow = plotHost().querySelector('g[data-detail-kind="request"]');
      probe.foundARow = !!firstRow;
      if (firstRow) {
        firstRow.dispatchEvent(new MouseEvent("click", { bubbles: true, cancelable: true, composed: true }));
      }
      // Wait for the RE-LAYOUT, not for the width change and not for a duration.
      //
      // The width narrows synchronously with the click, but the render it triggers is
      // scheduled through requestAnimationFrame — and headless --dump-dom has no
      // compositor driving frames, so under full-suite load that callback lands well
      // after the width has moved. Polling the width therefore measured a settled
      // container around an unsettled plot, and the probe reported the very defect it
      // exists to catch. Twice.
      //
      // Polling the OUTCOME is sound here because the wait is bounded: a genuine
      // regression spins out the attempts and then fails the assertion below, which is
      // exactly what removing the ResizeObserver does.
      var widthBeforeClick = probe.before.hostWidth;
      function someSegmentIsInsideThePlot(wantWidthChanged) {
        var host = plotHost();
        var hostBox = host.getBoundingClientRect();
        var widthChanged = Math.round(hostBox.width) !== widthBeforeClick;
        if (widthChanged !== wantWidthChanged) {
          return false;
        }
        return [].slice.call(host.querySelectorAll("rect.timeline-segment")).some(function (segment) {
          var box = segment.getBoundingClientRect();
          return box.right > hostBox.left && box.left < hostBox.right;
        });
      }
      settleUntil(function () {
        return someSegmentIsInsideThePlot(true);
      }, function () {
        probe.drawerOpenAfter = drawerIsOpen();
        probe.after = plotSnapshot("after");
        var close = document.getElementById("detail-close");
        if (close) { close.click(); }
        settleUntil(function () {
          return someSegmentIsInsideThePlot(false);
        }, function () {
          probe.closed = plotSnapshot("closed");
          probe.drawerOpenAtEnd = drawerIsOpen();
          document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify(probe);
          document.title = "READY";
        });
      });
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	// A LONGER VIRTUAL-TIME BUDGET, because this probe waits on conditions rather
	// than durations and every polled frame advances the virtual clock. The default
	// 5000ms is spent before the poll resolves, and the runner then dumps a DOM
	// whose result node is still empty. The flag repeats deliberately: Chromium
	// takes the last occurrence, and the runner's own default comes first.
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline drawer plot width", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")

	type plotSnapshot struct {
		Label     string `json:"label"`
		Href      string `json:"href"`
		HostWidth int    `json:"hostWidth"`
		HostRight int    `json:"hostRight"`
		Segments  int    `json:"segments"`
		Inside    int    `json:"inside"`
		Rightmost *int   `json:"rightmost"`
		Leftmost  *int   `json:"leftmost"`
	}
	var drawerResult struct {
		Before           plotSnapshot `json:"before"`
		After            plotSnapshot `json:"after"`
		Closed           plotSnapshot `json:"closed"`
		FoundARow        bool         `json:"foundARow"`
		DrawerOpenBefore bool         `json:"drawerOpenBefore"`
		DrawerOpenAfter  bool         `json:"drawerOpenAfter"`
		DrawerOpenAtEnd  bool         `json:"drawerOpenAtEnd"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawerResult); decodeError != nil {
		t.Fatalf("decode timeline drawer behavior: %v (output %q)", decodeError, probeOutput)
	}

	// Every measurement names the page it measured. A confident number about
	// somebody else's board has shipped here before.
	for _, snapshot := range []plotSnapshot{drawerResult.Before, drawerResult.After, drawerResult.Closed} {
		if !strings.HasSuffix(snapshot.Href, "/probe.html") {
			t.Fatalf("the %s snapshot was taken on %q, not the probe page", snapshot.Label, snapshot.Href)
		}
	}

	// SETUP, ASSERTED. Each of these silently turns the test below into a
	// measurement of nothing.
	if !drawerResult.FoundARow {
		t.Fatal("the probe found no [data-detail-kind] row in the plot, so it clicked nothing")
	}
	if drawerResult.DrawerOpenBefore {
		t.Fatal("the detail drawer was already open before the row click")
	}
	if !drawerResult.DrawerOpenAfter {
		t.Fatal("clicking a row did not open the detail drawer, so the plot was never narrowed " +
			"and this probe cannot see the defect it exists for")
	}
	if drawerResult.Before.Segments == 0 {
		t.Fatal("no timeline segments were drawn before the click; the probe measured an empty chart")
	}
	if drawerResult.After.HostWidth >= drawerResult.Before.HostWidth {
		t.Fatalf("the plot host was %dpx before the drawer opened and %dpx after, so the drawer did "+
			"not narrow it and there is nothing here to re-measure",
			drawerResult.Before.HostWidth, drawerResult.After.HostWidth)
	}

	// THE DEFECT. Fifty-five segments, none of them on screen.
	if drawerResult.After.Inside == 0 {
		t.Fatalf("after the drawer opened, %d segments are in the DOM and NONE of them overlap the "+
			"%dpx plot (leftmost %s, host right edge %d) — the chart is blank",
			drawerResult.After.Segments, drawerResult.After.HostWidth,
			plotEdgeText(drawerResult.After.Leftmost), drawerResult.After.HostRight)
	}
	// Not merely non-zero: the bars have to be laid out against the NEW plot, so
	// none of them may sit past its right edge. A stale layout that happened to
	// leave one bar clipped inside the box would satisfy the count alone.
	if drawerResult.After.Rightmost == nil || *drawerResult.After.Rightmost > drawerResult.After.HostRight+2 {
		t.Fatalf("after the drawer opened the rightmost segment ends at %s, past the plot's right "+
			"edge at %d — the bars were not re-laid out against the narrowed plot",
			plotEdgeText(drawerResult.After.Rightmost), drawerResult.After.HostRight)
	}
	// And closing it comes back, without the reader having to move the window.
	if drawerResult.DrawerOpenAtEnd {
		t.Fatal("the probe could not close the detail drawer, so the recovery half is untested")
	}
	if drawerResult.Closed.Inside == 0 {
		t.Fatalf("after the drawer closed, none of the %d segments overlap the %dpx plot",
			drawerResult.Closed.Segments, drawerResult.Closed.HostWidth)
	}
	if drawerResult.Closed.HostWidth != drawerResult.Before.HostWidth {
		t.Fatalf("closing the drawer left the plot %dpx wide, want the %dpx it started at",
			drawerResult.Closed.HostWidth, drawerResult.Before.HostWidth)
	}
	if drawerResult.Closed.Rightmost == nil || *drawerResult.Closed.Rightmost > drawerResult.Closed.HostRight+2 {
		t.Fatalf("after the drawer closed the rightmost segment ends at %s, past the plot's right edge at %d",
			plotEdgeText(drawerResult.Closed.Rightmost), drawerResult.Closed.HostRight)
	}
}

// The From / to fields, driven through the real events a reader generates.
//
// WHY A BROWSER. The defect is entirely in the interaction between focus, the
// engine's own `input` and `change` events on a date input, and the write-back a
// render performs. A sliced probe can call syncRangeField with any arguments it
// likes and learn nothing about which arguments the browser actually produces —
// and the guard that shipped was wrong precisely about that: it decided "the
// reader is mid-edit" by comparing values, and after a commit the field still
// holds the reader's text.
//
// Two states, both reachable in two keystrokes:
//
//	CLEAR the field and commit — the window must stand and the field must come
//	back, in a branch whose own comment says "Restore it unconditionally";
//	type a date the clamp MOVES — the field must show where the chart actually
//	landed, not the date that was rejected.
func TestBrowserBehaviorTimelineRangeFieldsShowTheWindowTheChartIsDrawnAt(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
// Same condition-not-duration rule as the drawer probe above, and the same
// setTimeout-not-rAF reason; see its note.
function settleUntil(predicate, thenDo) {
  var attemptsLeft = 60;
  (function attempt() {
    if (predicate() || attemptsLeft-- <= 0) {
      setTimeout(thenDo, 50);
      return;
    }
    setTimeout(attempt, 50);
  })();
}
function fieldState() {
  return {
    href: location.href,
    from: document.getElementById("timeline-range-start").value,
    to: document.getElementById("timeline-range-end").value,
    readout: document.getElementById("timeline-range-readout").textContent
  };
}
// The events a date input really emits, in the order it emits them, against the
// focused field — which is what makes the mid-edit guard observable at all.
function typeInto(fieldId, text) {
  var field = document.getElementById(fieldId);
  field.focus();
  field.value = text;
  field.dispatchEvent(new Event("input", { bubbles: true }));
  field.dispatchEvent(new Event("change", { bubbles: true }));
  return field;
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var probe = {};
      // Start from a window that is not the bounds, so a clamp is observable.
      document.querySelector('[data-timeline-period="month"]').click();
      probe.beforeClear = fieldState();

      // (a) CLEAR AND COMMIT. The field keeps focus throughout, which is the case
      // the old guard stranded.
      var cleared = typeInto("timeline-range-start", "");
      probe.afterClear = fieldState();
      probe.clearedStillFocused = document.activeElement === cleared;

      // (b) A DATE THE CLAMP MOVES: far past the end of the board's range.
      typeInto("timeline-range-start", "2099-12-31");
      probe.afterOutOfRange = fieldState();

      // (c) And an ordinary in-range date still applies exactly.
      document.querySelector('[data-timeline-period="month"]').click();
      var monthFrom = document.getElementById("timeline-range-start").value;
      typeInto("timeline-range-start", monthFrom);
      probe.afterReapply = fieldState();
      probe.monthFrom = monthFrom;

      // (d) A field left mid-edit and blurred without committing must be restored.
      var abandoned = document.getElementById("timeline-range-end");
      abandoned.focus();
      abandoned.value = "";
      abandoned.dispatchEvent(new Event("input", { bubbles: true }));
      probe.duringAbandon = fieldState();
      abandoned.dispatchEvent(new Event("blur", { bubbles: true }));
      abandoned.blur();
      // Wait for the restore, not for a duration.
      settleUntil(function () {
        return document.getElementById("timeline-range-end").value !== "";
      }, function () {
        probe.afterAbandon = fieldState();
        document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify(probe);
        document.title = "READY";
      });
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	// Same longer budget as the drawer probe, and for the same reason: this one
	// polls for the restore instead of sleeping through it.
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline range fields", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")

	type fieldState struct {
		Href    string `json:"href"`
		From    string `json:"from"`
		To      string `json:"to"`
		Readout string `json:"readout"`
	}
	var fieldResult struct {
		BeforeClear         fieldState `json:"beforeClear"`
		AfterClear          fieldState `json:"afterClear"`
		ClearedStillFocused bool       `json:"clearedStillFocused"`
		AfterOutOfRange     fieldState `json:"afterOutOfRange"`
		AfterReapply        fieldState `json:"afterReapply"`
		MonthFrom           string     `json:"monthFrom"`
		DuringAbandon       fieldState `json:"duringAbandon"`
		AfterAbandon        fieldState `json:"afterAbandon"`
	}
	if decodeError := json.Unmarshal(probeOutput, &fieldResult); decodeError != nil {
		t.Fatalf("decode timeline range-field behavior: %v (output %q)", decodeError, probeOutput)
	}

	for label, state := range map[string]fieldState{
		"before clear": fieldResult.BeforeClear, "after clear": fieldResult.AfterClear,
		"after out of range": fieldResult.AfterOutOfRange, "after reapply": fieldResult.AfterReapply,
		"after abandon": fieldResult.AfterAbandon,
	} {
		if !strings.HasSuffix(state.Href, "/probe.html") {
			t.Fatalf("the %s snapshot was taken on %q, not the probe page", label, state.Href)
		}
	}

	// SETUP, ASSERTED.
	if fieldResult.BeforeClear.From == "" || fieldResult.BeforeClear.Readout == "" {
		t.Fatalf("the fields started empty (%+v), so nothing below is measuring a restore",
			fieldResult.BeforeClear)
	}
	if !fieldResult.ClearedStillFocused {
		t.Fatal("the cleared field lost focus before the assertion, so this probe is not exercising " +
			"the mid-edit guard at all")
	}

	// (a) A cleared field comes back, and the window stands.
	if fieldResult.AfterClear.From == "" {
		t.Errorf("clearing the From field left it empty; it must be restored to the window the "+
			"chart is drawn at (readout %q)", fieldResult.AfterClear.Readout)
	}
	if fieldResult.AfterClear.Readout != fieldResult.BeforeClear.Readout {
		t.Errorf("clearing a field moved the window from %q to %q; an empty field is not a request "+
			"to move", fieldResult.BeforeClear.Readout, fieldResult.AfterClear.Readout)
	}
	if fieldResult.AfterClear.From != fieldResult.BeforeClear.From {
		t.Errorf("clearing the From field restored it to %q, want the %q it held before",
			fieldResult.AfterClear.From, fieldResult.BeforeClear.From)
	}

	// (b) A clamped commit shows where the chart landed, not what was typed.
	if fieldResult.AfterOutOfRange.From == "2099-12-31" {
		t.Errorf("after typing a date past the end of the range the From field still reads %q "+
			"while the chart is drawn at %q; the field has to name the window that exists",
			fieldResult.AfterOutOfRange.From, fieldResult.AfterOutOfRange.Readout)
	}
	if !strings.HasPrefix(fieldResult.AfterOutOfRange.Readout, fieldResult.AfterOutOfRange.From) {
		t.Errorf("after the clamp the From field reads %q and the readout starts %q; the two must "+
			"describe one window", fieldResult.AfterOutOfRange.From, fieldResult.AfterOutOfRange.Readout)
	}

	// (c) An in-range date still applies exactly, so the fix did not buy field
	// honesty by ignoring the reader.
	if fieldResult.AfterReapply.From != fieldResult.MonthFrom {
		t.Errorf("re-typing the month window's own start gave From %q, want %q",
			fieldResult.AfterReapply.From, fieldResult.MonthFrom)
	}

	// (d) Mid-edit is respected while it lasts, and released on blur.
	if fieldResult.DuringAbandon.To != "" {
		t.Errorf("a render overwrote the end field while the reader was part-way through typing "+
			"into it (it reads %q); mid-edit means hands off", fieldResult.DuringAbandon.To)
	}
	if fieldResult.AfterAbandon.To == "" {
		t.Errorf("the end field was left empty after the reader blurred without committing; "+
			"the window is %q and the field has to name it", fieldResult.AfterAbandon.Readout)
	}
}

// Where the Now and Fit all buttons LAND, and what the toolbar says about itself
// once they have.
//
// Now sized its window from the span between the now-line and the forecast's
// queue-empty instant, floored on half the ZOOM FLOOR. On a queue that is nearly
// drained that span is minutes, so Now landed on a one-hour window — exactly the
// floor — and the obvious next move was dead: the + button, ctrl+wheel and the +
// key were all silent no-ops with no disabled state and no message. Fit all had
// the mirror problem: it assigned the payload's whole range, so filtered to one
// domain it left most of the plot blank.
//
// Driven in a real engine because three of the four properties are about the
// toolbar's own rendered state, and the fourth (Fit all under a filter) needs the
// shared filter machinery this view reads but does not own.
func TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function settleUntil(predicate, thenDo) {
  var attemptsLeft = 60;
  (function attempt() {
    if (predicate() || attemptsLeft-- <= 0) { setTimeout(thenDo, 50); return; }
    setTimeout(attempt, 50);
  })();
}
// The window's two ends, in milliseconds, read back off the readout the reader
// sees rather than out of any internal state. One parser, so a span and an
// endpoint can never be read from two different places.
function readoutWindow() {
  var text = document.getElementById("timeline-range-readout").textContent || "";
  var match = text.match(/^(\S+ \S+) UTC → (\S+ \S+) UTC$/);
  if (!match) { return null; }
  var startMs = Date.parse(match[1].replace(" ", "T") + "Z");
  var endMs = Date.parse(match[2].replace(" ", "T") + "Z");
  if (isNaN(startMs) || isNaN(endMs)) { return null; }
  return { startMs: startMs, endMs: endMs };
}
function toolbarState(label) {
  var disabled = {};
  ["timeline-zoom-in", "timeline-zoom-out", "timeline-zoom-fit",
   "timeline-period-prev", "timeline-period-next"].forEach(function (buttonId) {
    disabled[buttonId] = !!document.getElementById(buttonId).disabled;
  });
  var rowsHost = document.querySelector("#view-timeline .timeline-scroll");
  var readWindow = readoutWindow();
  return {
    label: label,
    href: location.href,
    readout: document.getElementById("timeline-range-readout").textContent,
    startMs: readWindow ? readWindow.startMs : null,
    endMs: readWindow ? readWindow.endMs : null,
    spanMs: readWindow ? readWindow.endMs - readWindow.startMs : null,
    nowRuleDrawn: !!rowsHost.querySelector(".timeline-now-rule"),
    drawnSegments: rowsHost.querySelectorAll("rect.timeline-segment").length,
    disabled: disabled
  };
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var probe = {};
      probe.fitted = toolbarState("fitted");
      document.getElementById("timeline-zoom-now").click();
      probe.afterNow = toolbarState("afterNow");
      document.getElementById("timeline-zoom-in").click();
      probe.afterNowThenZoomIn = toolbarState("afterNowThenZoomIn");
      // The current week, then one step forward. Which regime that step is in
      // depends on the live queue — see the Go side's clause (3).
      document.querySelector('[data-timeline-period="week"]').click();
      probe.currentWeek = toolbarState("currentWeek");
      document.getElementById("timeline-period-next").click();
      probe.afterStepFromCurrentWeek = toolbarState("afterStepFromCurrentWeek");
      // Fit all with a filter on. The search box is the shared filter every view
      // reads, so this needs no timeline-specific plumbing.
      var search = document.getElementById("filter-search");
      search.value = "REQ-164";
      search.dispatchEvent(new Event("input", { bubbles: true }));
      settleUntil(function () {
        var summary = document.getElementById("timeline-summary").textContent || "";
        return summary.indexOf("1 REQ in the window") !== -1 ||
          summary.indexOf("Nothing was drawn") !== -1 ||
          summary.indexOf("No REQ matches") !== -1;
      }, function () {
        document.getElementById("timeline-zoom-fit").click();
        settleUntil(function () { return true; }, function () {
          probe.filteredFit = toolbarState("filteredFit");
          probe.filteredSummary = document.getElementById("timeline-summary").textContent;
          // Still filtered to that one archived REQ: the week holding it, then a
          // step forward out of it. This is the refusal case a live queue always
          // supplies no matter what its forecast is doing — the filtered set ends
          // where that REQ ended, years of queue activity later is irrelevant to
          // it. The Go side reads the fitted extent above to confirm that before
          // asserting anything about it.
          document.querySelector('[data-timeline-period="week"]').click();
          probe.filteredWeek = toolbarState("filteredWeek");
          document.getElementById("timeline-period-next").click();
          probe.filteredWeekThenStep = toolbarState("filteredWeekThenStep");
          document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify(probe);
          document.title = "READY";
        });
      });
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline now and fit all", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")

	type toolbarState struct {
		Label         string          `json:"label"`
		Href          string          `json:"href"`
		Readout       string          `json:"readout"`
		StartMs       *float64        `json:"startMs"`
		EndMs         *float64        `json:"endMs"`
		SpanMs        *float64        `json:"spanMs"`
		NowRuleDrawn  bool            `json:"nowRuleDrawn"`
		DrawnSegments int             `json:"drawnSegments"`
		Disabled      map[string]bool `json:"disabled"`
	}
	var landingResult struct {
		Fitted                   toolbarState `json:"fitted"`
		AfterNow                 toolbarState `json:"afterNow"`
		AfterNowThenZoomIn       toolbarState `json:"afterNowThenZoomIn"`
		CurrentWeek              toolbarState `json:"currentWeek"`
		AfterStepFromCurrentWeek toolbarState `json:"afterStepFromCurrentWeek"`
		FilteredFit              toolbarState `json:"filteredFit"`
		FilteredSummary          string       `json:"filteredSummary"`
		FilteredWeek             toolbarState `json:"filteredWeek"`
		FilteredWeekThenStep     toolbarState `json:"filteredWeekThenStep"`
	}
	if decodeError := json.Unmarshal(probeOutput, &landingResult); decodeError != nil {
		t.Fatalf("decode timeline landing behavior: %v (output %q)", decodeError, probeOutput)
	}

	states := []toolbarState{
		landingResult.Fitted, landingResult.AfterNow, landingResult.AfterNowThenZoomIn,
		landingResult.CurrentWeek, landingResult.AfterStepFromCurrentWeek, landingResult.FilteredFit,
		landingResult.FilteredWeek, landingResult.FilteredWeekThenStep,
	}
	for _, state := range states {
		if !strings.HasSuffix(state.Href, "/probe.html") {
			t.Fatalf("the %s state was measured on %q, not the probe page", state.Label, state.Href)
		}
		if state.SpanMs == nil || state.StartMs == nil || state.EndMs == nil {
			t.Fatalf("the %s state's readout %q did not parse as a window", state.Label, state.Readout)
		}
	}

	const oneHourMs = 3600000.0

	// (1) Now lands on a window, not on the zoom floor, and the now-line is in it.
	if *landingResult.AfterNow.SpanMs <= oneHourMs {
		t.Errorf("Now landed on a %.2f-hour window (%s); the zoom floor is one hour, so at or "+
			"below it the reader has nowhere to zoom and no context around the now-line",
			*landingResult.AfterNow.SpanMs/oneHourMs, landingResult.AfterNow.Readout)
	}
	if !landingResult.AfterNow.NowRuleDrawn {
		t.Errorf("Now landed on %s with no now-rule drawn in it", landingResult.AfterNow.Readout)
	}

	// (2) And zoom-in is alive there, in both the state it reports and what it does.
	if landingResult.AfterNow.Disabled["timeline-zoom-in"] {
		t.Errorf("zoom-in reports itself disabled in the window Now lands on (%s)",
			landingResult.AfterNow.Readout)
	}
	if !(*landingResult.AfterNowThenZoomIn.SpanMs < *landingResult.AfterNow.SpanMs) {
		t.Errorf("one zoom-in press after Now left the window at %.2f hours, from %.2f; the press "+
			"has to narrow it", *landingResult.AfterNowThenZoomIn.SpanMs/oneHourMs,
			*landingResult.AfterNow.SpanMs/oneHourMs)
	}

	// (3) A step past everything drawn does not happen, and says so first.
	//
	// WHICH REGIME THE CURRENT WEEK IS IN IS DATA, NOT A CONSTANT, and asserting
	// one of them here is what made capturing three REQs fail this test. On a
	// drained queue the week after this one exists only inside the cosmetic bound
	// padding and the arrow must refuse; while the forecast reaches past this
	// week's end — which one ordinary `capture-request` is enough to do — that
	// next period holds real projected bars and the arrow is right to be enabled.
	// So the arrow's own verdict is READ and the press is checked against it,
	// which is the contract in both regimes: it never lands the reader on an empty
	// chart, and whatever it is about to do it says first. The refusal branch is
	// then exercised deterministically under a filter in (6) below, where the live
	// queue cannot move the answer.
	if landingResult.CurrentWeek.Disabled["timeline-period-next"] {
		if landingResult.AfterStepFromCurrentWeek.Readout != landingResult.CurrentWeek.Readout {
			t.Errorf("the step-forward arrow reported itself disabled on the current week (%s) and "+
				"the press moved the window to %s anyway", landingResult.CurrentWeek.Readout,
				landingResult.AfterStepFromCurrentWeek.Readout)
		}
	} else {
		if landingResult.AfterStepFromCurrentWeek.Readout == landingResult.CurrentWeek.Readout {
			t.Errorf("the step-forward arrow reported itself enabled on the current week (%s) and "+
				"the press did not move the window", landingResult.CurrentWeek.Readout)
		}
		if landingResult.AfterStepFromCurrentWeek.DrawnSegments == 0 {
			t.Errorf("the step-forward arrow was enabled on the current week (%s) and the press "+
				"landed on %s with nothing drawn in it — a step past everything drawn",
				landingResult.CurrentWeek.Readout, landingResult.AfterStepFromCurrentWeek.Readout)
		}
	}

	// (4) Fit all fits WHAT IS ON SCREEN. The fitted window under a one-row filter
	// must be a small fraction of the unfiltered one, and must still draw that row.
	if landingResult.FilteredFit.DrawnSegments == 0 {
		t.Fatalf("Fit all under the filter drew no segments (summary %q), so the span comparison "+
			"below is measuring an empty chart", landingResult.FilteredSummary)
	}
	if *landingResult.FilteredFit.SpanMs >= *landingResult.Fitted.SpanMs/2 {
		t.Errorf("Fit all under a one-row filter produced a %.1f-day window against the unfiltered "+
			"%.1f days; it has to fit the rows on screen, not the payload's whole range",
			*landingResult.FilteredFit.SpanMs/(24*oneHourMs),
			*landingResult.Fitted.SpanMs/(24*oneHourMs))
	}
	// And zoom-out still reaches past the filtered extent, so fitting the filter did
	// not trap the reader inside it.
	if landingResult.FilteredFit.Disabled["timeline-zoom-out"] {
		t.Error("zoom-out is disabled after Fit all under a filter; the clamp bounds must stay the " +
			"payload's so the reader can look outside the filtered extent")
	}

	// (5) And at the full-range window the controls that cannot act say so. Without
	// this the disabled-state assertions above could all pass against code that
	// simply never disables anything.
	if !landingResult.Fitted.Disabled["timeline-zoom-out"] {
		t.Error("zoom-out is enabled at the full-range window, where there is nothing to zoom out to")
	}

	// (6) The refusal itself, on a set the live queue cannot move: one archived
	// REQ, and the week that holds it.
	//
	// THE PREMISE IS READ, NOT RESTATED. Fit all lands on the drawn extent plus
	// breathing room, so the filtered fit window is an outer bound on everything
	// drawn under that filter; if it ends before this week does, the next period
	// holds nothing at all and the arrow owes the reader a refusal. Read it that
	// way round rather than asserting a week number, because a hardcoded premise
	// about what the queue holds is exactly what clause (3) had to stop doing.
	// The readout is minute-truncated, so the fitted end can be up to a minute
	// later than it reads — the allowance keeps the comparison sound rather than
	// merely true.
	const readoutTruncationMs = 60000.0
	if landingResult.FilteredWeek.DrawnSegments == 0 {
		t.Fatalf("the week the chip landed on under the one-REQ filter (%s) draws nothing, so the "+
			"refusal below would be about an empty chart rather than about a step off the end of "+
			"the filtered REQ (summary %q)", landingResult.FilteredWeek.Readout,
			landingResult.FilteredSummary)
	}
	if *landingResult.FilteredFit.EndMs+readoutTruncationMs >= *landingResult.FilteredWeek.EndMs {
		t.Fatalf("everything drawn under the one-REQ filter (fitted to %s) reaches to the end of "+
			"the week the chip landed on (%s), so a step out of that week is not a step past "+
			"everything drawn and this case proves nothing",
			landingResult.FilteredFit.Readout, landingResult.FilteredWeek.Readout)
	}
	if !landingResult.FilteredWeek.Disabled["timeline-period-next"] {
		t.Errorf("the step-forward arrow is enabled on %s, whose next period is past everything "+
			"drawn under the filter (fitted to %s)", landingResult.FilteredWeek.Readout,
			landingResult.FilteredFit.Readout)
	}
	if landingResult.FilteredWeekThenStep.Readout != landingResult.FilteredWeek.Readout {
		t.Errorf("pressing the step-forward arrow moved the window from %s to %s, past everything "+
			"drawn under the filter (fitted to %s), and drew %d segments there",
			landingResult.FilteredWeek.Readout, landingResult.FilteredWeekThenStep.Readout,
			landingResult.FilteredFit.Readout, landingResult.FilteredWeekThenStep.DrawnSegments)
	}
}

// The two sentences this view emits that used to describe a chart nobody was
// looking at, plus the legend entry the plot's third vertical line never had.
//
// The summary always appended "measured to the now-line at <t>". drawNowRule draws
// nothing when now falls outside the window, so on any past week that pointed at a
// rule the reader could not find, beside open bars clipped flush at the frame.
//
// And the forecast contradicted itself inside one sentence pair: "This covers the
// whole queue, not the rows shown." followed by "Nothing left to schedule — every
// remaining REQ is listed below.", above a single row, with the excluded paragraph
// immediately underneath naming a REQ that was not listed anywhere.
func TestBrowserBehaviorTimelineProseDescribesOnlyTheWindowOnScreen(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function settleUntil(predicate, thenDo) {
  var attemptsLeft = 60;
  (function attempt() {
    if (predicate() || attemptsLeft-- <= 0) { setTimeout(thenDo, 50); return; }
    setTimeout(attempt, 50);
  })();
}
function proseState(label) {
  var rowsHost = document.querySelector("#view-timeline .timeline-scroll");
  return {
    label: label,
    href: location.href,
    summary: document.getElementById("timeline-summary").textContent,
    forecast: document.getElementById("timeline-forecast").textContent,
    excluded: document.getElementById("timeline-excluded").textContent,
    readout: document.getElementById("timeline-range-readout").textContent,
    nowRuleDrawn: !!rowsHost.querySelector(".timeline-now-rule")
  };
}
function typeWindow(fromText, toText) {
  [["timeline-range-start", fromText], ["timeline-range-end", toText]].forEach(function (pair) {
    var field = document.getElementById(pair[0]);
    field.focus();
    field.value = pair[1];
    field.dispatchEvent(new Event("input", { bubbles: true }));
    field.dispatchEvent(new Event("change", { bubbles: true }));
    field.blur();
  });
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var probe = {};
      probe.withNow = proseState("withNow");
      // A window well before the now-line that still HAS rows in it. Late July is
      // the busiest stretch of this repo's own archive; an empty past week would
      // take the "Nothing was drawn" branch instead, which is a different sentence
      // and would leave the one under test unexercised.
      typeWindow("2026-07-27", "2026-08-02");
      settleUntil(function () {
        return (document.getElementById("timeline-range-readout").textContent || "")
          .indexOf("2026-07-27") !== -1;
      }, function () {
        probe.withoutNow = proseState("withoutNow");
        // And the forecast under a filter that leaves a subset on screen.
        var search = document.getElementById("filter-search");
        search.value = "REQ-164";
        search.dispatchEvent(new Event("input", { bubbles: true }));
        settleUntil(function () {
          return (document.getElementById("timeline-forecast").textContent || "").length > 0 ||
            (document.getElementById("timeline-summary").textContent || "").indexOf("No REQ matches") !== -1;
        }, function () {
          document.getElementById("timeline-zoom-fit").click();
          settleUntil(function () { return true; }, function () {
            probe.filtered = proseState("filtered");
            probe.legendRules = [].slice.call(
              document.querySelectorAll("#view-timeline .timeline-legend .timeline-swatch[data-rule]")
            ).map(function (swatch) { return swatch.getAttribute("data-rule"); });
            document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify(probe);
            document.title = "READY";
          });
        });
      });
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline prose", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")

	type proseState struct {
		Label        string `json:"label"`
		Href         string `json:"href"`
		Summary      string `json:"summary"`
		Forecast     string `json:"forecast"`
		Excluded     string `json:"excluded"`
		Readout      string `json:"readout"`
		NowRuleDrawn bool   `json:"nowRuleDrawn"`
	}
	var proseResult struct {
		WithNow     proseState `json:"withNow"`
		WithoutNow  proseState `json:"withoutNow"`
		Filtered    proseState `json:"filtered"`
		LegendRules []string   `json:"legendRules"`
	}
	if decodeError := json.Unmarshal(probeOutput, &proseResult); decodeError != nil {
		t.Fatalf("decode timeline prose behavior: %v (output %q)", decodeError, probeOutput)
	}

	for _, state := range []proseState{proseResult.WithNow, proseResult.WithoutNow, proseResult.Filtered} {
		if !strings.HasSuffix(state.Href, "/probe.html") {
			t.Fatalf("the %s state was measured on %q, not the probe page", state.Label, state.Href)
		}
	}

	// SETUP, ASSERTED: the two window states have to differ in whether the rule is
	// drawn, or the wording assertions are comparing one state with itself.
	if !proseResult.WithNow.NowRuleDrawn {
		t.Fatalf("the fitted window drew no now-rule (readout %q), so the case where one IS drawn "+
			"is untested", proseResult.WithNow.Readout)
	}
	if proseResult.WithoutNow.NowRuleDrawn {
		t.Fatalf("the past window still draws a now-rule (readout %q); this probe needs a window "+
			"the now-line falls outside of", proseResult.WithoutNow.Readout)
	}
	// And it has to have ROWS in it, or the render takes the empty-window branch and
	// the sentence under test is never emitted.
	if !strings.Contains(proseResult.WithoutNow.Summary, "still open") {
		t.Fatalf("the past window produced %q rather than a summary of drawn rows; pick a window "+
			"with REQs in it", proseResult.WithoutNow.Summary)
	}

	// (P1) The pointer appears only where the rule does.
	if !strings.Contains(proseResult.WithNow.Summary, "measured to the now-line at") {
		t.Errorf("with the now-rule drawn the summary reads %q; it should name the rule the reader "+
			"can see", proseResult.WithNow.Summary)
	}
	if strings.Contains(proseResult.WithoutNow.Summary, "now-line") {
		t.Errorf("with no now-rule drawn the summary still points at one: %q",
			proseResult.WithoutNow.Summary)
	}
	// The instant itself still has to be stated: it is what the open spans were
	// measured against, and dropping it would trade one defect for another.
	if !strings.Contains(proseResult.WithoutNow.Summary, "measured against") ||
		!strings.Contains(proseResult.WithoutNow.Summary, " UTC") {
		t.Errorf("with no now-rule drawn the summary no longer states the instant open spans were "+
			"measured against: %q", proseResult.WithoutNow.Summary)
	}

	// (P2) The forecast does not both deny and assert that the rows are everything.
	if strings.Contains(proseResult.Filtered.Forecast, "not the rows shown") &&
		strings.Contains(proseResult.Filtered.Forecast, "listed below") {
		t.Errorf("the forecast says both \"not the rows shown\" and \"listed below\" in one "+
			"paragraph: %q", proseResult.Filtered.Forecast)
	}
	if proseResult.Filtered.Excluded != "" && strings.Contains(proseResult.Filtered.Forecast, "listed below") {
		t.Errorf("the forecast claims every remaining REQ is listed below while the excluded "+
			"paragraph under it names one that is not: forecast %q, excluded %q",
			proseResult.Filtered.Forecast, proseResult.Filtered.Excluded)
	}

	// (P4) Every kind of vertical line the plot draws has a key. The renderer draws
	// three: the now rule, the queue-end rule, and one gridline per axis tick.
	wantLegendRules := map[string]bool{"now": false, "queue-end": false, "gridline": false}
	for _, rule := range proseResult.LegendRules {
		if _, isKnown := wantLegendRules[rule]; !isKnown {
			t.Errorf("the legend keys a vertical rule %q the plot does not draw", rule)
			continue
		}
		wantLegendRules[rule] = true
	}
	for rule, isKeyed := range wantLegendRules {
		if !isKeyed {
			t.Errorf("the plot draws a %q vertical line with no legend entry, so a reader has "+
				"nothing to look it up under", rule)
		}
	}
}

// The pointer and keyboard paths, driven with the events a reader really produces.
//
// WHY A BROWSER. Both defects are about what the ENGINE does around the handlers:
// which release events reach a host when a drag ends outside it, and when a scroll
// event is delivered relative to a synchronous focus call. Neither is visible to a
// probe that calls the handlers directly.
//
// Reproduced before fixing, and the reproduction corrected the report on two counts:
// the stuck drag leaves the GRAB CURSOR on but does not go on panning the window,
// and Tab is not trapped in the rows at all — it walks the rendered rows and exits
// after about thirty presses. Only what reproduced is pinned here.
func TestBrowserBehaviorTimelinePointerAndKeyboardPathsStayAlive(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function plotHost() { return document.querySelector("#view-timeline .timeline-scroll"); }
function chartState() {
  var active = document.activeElement;
  return {
    href: location.href,
    readout: document.getElementById("timeline-range-readout").textContent,
    grabbing: plotHost().classList.contains("is-panning"),
    focusedRowId: active && active.getAttribute ? (active.getAttribute("data-detail-id") || "") : "",
    focusIsInsideTheChart: !!(active && plotHost().contains(active)) || active === plotHost()
  };
}
function pointerAt(type, node, clientX, clientY, buttons) {
  node.dispatchEvent(new PointerEvent(type, {
    bubbles: true, cancelable: true, composed: true,
    clientX: clientX, clientY: clientY, button: 0, buttons: buttons,
    pointerId: 1, pointerType: "mouse"
  }));
}
function nodeAt(x, y, fallbackNode) { return document.elementFromPoint(x, y) || fallbackNode; }
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var probe = {};
      // Off the bounds, so a pan has room to move in both directions.
      document.getElementById("timeline-zoom-in").click();
      document.getElementById("timeline-zoom-in").click();
      document.querySelector("#view-timeline .timeline-chart").scrollIntoView({ block: "start" });

      // (a) A DRAG RELEASED OUTSIDE THE CHART. The release is dispatched at a point
      // above the host, which is where the engine used to deliver it — leaving the
      // grab cursor on for the rest of the session.
      var host = plotHost();
      var hostBox = host.getBoundingClientRect();
      var pressY = Math.round(hostBox.top + 20);
      var pressX = Math.round(hostBox.left + hostBox.width * 0.6);
      pointerAt("pointerdown", host, pressX, pressY, 1);
      pointerAt("pointermove", host, pressX - 60, pressY, 1);
      pointerAt("pointermove", host, pressX - 120, pressY, 1);
      probe.duringDrag = chartState();
      // THE RELEASE THE ENGINE SENDS WHEN A CAPTURED POINTER GOES AWAY. A synthetic
      // PointerEvent carries a pointerId the engine does not know, so
      // setPointerCapture on it throws and this lane cannot reproduce a real captured
      // drag end-to-end. What it CAN drive is the event the capture path relies on,
      // and that is the half that was missing: with the release delivered outside the
      // host and the boundary events suppressed while a button is held, nothing
      // reached this host at all and the grab cursor stayed on for the session.
      var outsideY = Math.round(hostBox.top - 200);
      var outsideNode = nodeAt(pressX - 160, outsideY, document.body);
      pointerAt("pointerup", outsideNode, pressX - 160, outsideY, 0);
      probe.afterReleaseOutside = chartState();
      host.dispatchEvent(new PointerEvent("lostpointercapture", {
        bubbles: true, composed: true, pointerId: 1, pointerType: "mouse"
      }));
      probe.afterLostCapture = chartState();
      // And un-buttoned motion back over the chart must not move anything.
      pointerAt("pointermove", host, Math.round(hostBox.left + 40), pressY, 0);
      probe.afterUnbuttonedReentry = chartState();

      // (b) ARROW KEYS WITH A ROW FOCUSED. Three presses: the first worked before
      // this fix and the second did not, so one press proves nothing.
      document.getElementById("timeline-zoom-fit").click();
      document.getElementById("timeline-zoom-in").click();
      document.getElementById("timeline-zoom-in").click();
      setTimeout(function () {
        var firstRow = plotHost().querySelector('g[data-detail-kind="request"]');
        probe.foundARow = !!firstRow;
        if (firstRow) { firstRow.focus(); }
        probe.beforeArrows = chartState();
        var arrowStates = [];
        for (var press = 0; press < 3; press++) {
          plotHost().dispatchEvent(new KeyboardEvent("keydown", {
            key: "ArrowRight", bubbles: true, cancelable: true, composed: true
          }));
          arrowStates.push(chartState());
        }
        probe.arrows = arrowStates;

        // (c) CTRL+WHEEL OVER THE AXIS STRIP, which the hint promises zooms the
        // time axis and which did nothing at all.
        var axisHost = document.getElementById("timeline-axis");
        var axisBox = axisHost.getBoundingClientRect();
        probe.beforeAxisWheel = chartState();
        axisHost.dispatchEvent(new WheelEvent("wheel", {
          bubbles: true, cancelable: true, composed: true, ctrlKey: true, deltaY: -120,
          clientX: Math.round(axisBox.left + axisBox.width * 0.6),
          clientY: Math.round(axisBox.top + axisBox.height / 2)
        }));
        probe.afterAxisWheel = chartState();

        document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify(probe);
        document.title = "READY";
      }, 300);
    }, 500);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline pointer and keyboard", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")

	type chartState struct {
		Href                  string `json:"href"`
		Readout               string `json:"readout"`
		Grabbing              bool   `json:"grabbing"`
		FocusedRowId          string `json:"focusedRowId"`
		FocusIsInsideTheChart bool   `json:"focusIsInsideTheChart"`
	}
	var pathResult struct {
		DuringDrag             chartState   `json:"duringDrag"`
		AfterReleaseOutside    chartState   `json:"afterReleaseOutside"`
		AfterLostCapture       chartState   `json:"afterLostCapture"`
		AfterUnbuttonedReentry chartState   `json:"afterUnbuttonedReentry"`
		FoundARow              bool         `json:"foundARow"`
		BeforeArrows           chartState   `json:"beforeArrows"`
		Arrows                 []chartState `json:"arrows"`
		BeforeAxisWheel        chartState   `json:"beforeAxisWheel"`
		AfterAxisWheel         chartState   `json:"afterAxisWheel"`
	}
	if decodeError := json.Unmarshal(probeOutput, &pathResult); decodeError != nil {
		t.Fatalf("decode timeline pointer/keyboard behavior: %v (output %q)", decodeError, probeOutput)
	}

	// SETUP, ASSERTED.
	if !strings.HasSuffix(pathResult.DuringDrag.Href, "/probe.html") {
		t.Fatalf("measured on %q, not the probe page", pathResult.DuringDrag.Href)
	}
	if !pathResult.DuringDrag.Grabbing {
		t.Fatalf("the drag never engaged (readout %q), so the release below is not ending one",
			pathResult.DuringDrag.Readout)
	}
	if !pathResult.FoundARow {
		t.Fatal("the probe found no row to focus, so the keyboard half measures nothing")
	}
	if pathResult.BeforeArrows.FocusedRowId == "" {
		t.Fatal("the row did not take focus, so the arrow presses below are not testing the " +
			"focus-restore path at all")
	}
	if len(pathResult.Arrows) != 3 {
		t.Fatalf("the probe recorded %d arrow presses, want 3", len(pathResult.Arrows))
	}

	// (a) A drag ENDS when the pointer goes away, however it goes away. The grab
	// cursor was the visible half: it stayed on for the rest of the session.
	if pathResult.AfterLostCapture.Grabbing {
		t.Errorf("the drag survived losing its pointer capture and the grab cursor is still on; " +
			"lostpointercapture is the one release event the engine sends whatever the pointer did")
	}
	if pathResult.AfterUnbuttonedReentry.Readout != pathResult.AfterLostCapture.Readout {
		t.Errorf("moving the pointer back over the chart with no button held moved the window from "+
			"%s to %s", pathResult.AfterLostCapture.Readout,
			pathResult.AfterUnbuttonedReentry.Readout)
	}
	// And the capture is actually requested, which is what makes the release above a
	// fact rather than a hope. Structural, because a synthetic pointerId cannot be
	// captured in this lane.
	//
	// REQ-336 moved the request off pointerdown and onto the pan's ENGAGE: capture also
	// retargets the synthesized click, so taking it on every press cost every mouse
	// click in the chart its [data-detail-kind] target. This asserts the same contract
	// at its new instant, in two halves so deleting either one fails: the capture call
	// lives in capturePanPointer, and the move handler calls it.
	//
	// The CALL, not the feature-detect guard beside it: a check for the bare name
	// matched the `typeof scrollHost.setPointerCapture === "function"` line and passed
	// with the call itself deleted.
	capturePanPointerBody := sliceBalancedBlockAfter(t, indexHTML, "function capturePanPointer(")
	if !strings.Contains(capturePanPointerBody, "scrollHost.setPointerCapture(") {
		t.Error("capturePanPointer does not capture the pointer, so a drag released outside " +
			"the chart has no guaranteed path back to the host that armed it")
	}
	pointerMoveBody := sliceBalancedBlockAfter(t, indexHTML, "addTimelineListener(scrollHost, \"pointermove\"")
	if !strings.Contains(pointerMoveBody, "capturePanPointer()") {
		t.Error("the pointermove handler never calls capturePanPointer, so nothing takes the " +
			"capture when the pan engages and the release is a hope again")
	}

	// (b) EVERY arrow press pans, not just the first. The second press was the dead
	// one, so a single-press test would have passed over this.
	previousReadout := pathResult.BeforeArrows.Readout
	for pressIndex, afterPress := range pathResult.Arrows {
		if afterPress.Readout == previousReadout {
			t.Errorf("arrow press %d did not move the window (still %s); focus was on %q and inside "+
				"the chart: %v", pressIndex+1, afterPress.Readout, afterPress.FocusedRowId,
				afterPress.FocusIsInsideTheChart)
		}
		if !afterPress.FocusIsInsideTheChart {
			t.Errorf("after arrow press %d focus left the chart, so the next key cannot reach the "+
				"handler", pressIndex+1)
		}
		previousReadout = afterPress.Readout
	}

	// (c) The axis strip zooms, which the hint under the chart promises.
	if pathResult.AfterAxisWheel.Readout == pathResult.BeforeAxisWheel.Readout {
		t.Errorf("ctrl+wheel over the axis strip left the window at %s; the hint says holding Ctrl "+
			"and scrolling zooms the time axis, and the axis is where a reader aims",
			pathResult.BeforeAxisWheel.Readout)
	}
}

// timelinePointerCaptureProbeHelpers is the page-side half of the probe below: the
// state readers, the gesture arming, the aim, and the two mutations. It is installed
// once into the real board and then driven from Go over the protocol channel, so
// nothing here dispatches an event — every event in this probe comes from the engine.
const timelinePointerCaptureProbeHelpers = `(function () {
  var probe = {
    clickCount: 0, clickDetailId: "", clickTag: "",
    pointerUpSeen: false, pointerUpTag: "",
    captureEstablished: null, captureListener: null
  };
  window.timelineProbe = probe;
  probe.plotHost = function () { return document.querySelector("#view-timeline .timeline-scroll"); };
  probe.isPanning = function () { return probe.plotHost().classList.contains("is-panning"); };
  probe.drawerIsOpen = function () {
    var drawer = document.getElementById("detail-drawer");
    return !!drawer && !drawer.hidden && !drawer.classList.contains("is-hidden");
  };
  probe.shownDetailId = function () {
    var shown = document.getElementById("detail-id");
    return shown ? shown.textContent.trim() : "";
  };
  probe.closeDrawer = function () {
    var closeButton = document.getElementById("detail-close");
    if (closeButton) { closeButton.click(); }
    return "closed";
  };
  // Every gesture starts from the SAME window, and one that is off both bounds: at
  // Fit all a pan clamps to the window it started in, so a drag that panned would be
  // indistinguishable from one that did not. REQ-324's lesson, met again.
  probe.resetWindow = function () {
    document.getElementById("timeline-zoom-fit").click();
    document.getElementById("timeline-zoom-in").click();
    document.getElementById("timeline-zoom-in").click();
    return "reset";
  };
  // Every watch sits on DOCUMENT in the BUBBLE phase, so they run AFTER the board's
  // own delegated click handler and after the scroll host's release teardown. A
  // capture-phase listener would report the state from before the thing being
  // measured happened.
  document.addEventListener("click", function (clickEvent) {
    probe.clickCount++;
    var clickTarget = clickEvent.target;
    probe.clickTag = clickTarget && clickTarget.tagName ? clickTarget.tagName : "";
    var trigger = clickTarget && clickTarget.closest ? clickTarget.closest("[data-detail-kind]") : null;
    probe.clickDetailId = trigger ? trigger.getAttribute("data-detail-id") : "";
  }, false);
  document.addEventListener("pointerup", function (upEvent) {
    probe.pointerUpSeen = true;
    probe.pointerUpTag = upEvent.target && upEvent.target.tagName ? upEvent.target.tagName : "";
  }, false);
  // THE CLICK-SYNTHESIS PRECONDITION, watched rather than assumed. Chromium creates a
  // click on the nearest common inclusive ancestor of the mousedown and mouseup
  // targets. A mousedown target that has been DETACHED in between has no ancestor in
  // common with anything, so no click is created at all — not a click on the wrong
  // element, none. That is measured, not theorised: forcing this view's own rebuild
  // path between a press and a release produced clickCount 0 three times out of three.
  // These two fields are what turns that from a hang into a sentence.
  document.addEventListener("mousedown", function (downEvent) {
    probe.pressTarget = downEvent.target;
    probe.pressTargetTag = downEvent.target && downEvent.target.tagName ? downEvent.target.tagName : "";
  }, false);
  document.addEventListener("mouseup", function () {
    probe.pressTargetSurvived = probe.pressTarget ? probe.pressTarget.isConnected : null;
  }, false);
  probe.armGesture = function () {
    probe.clickCount = 0; probe.clickDetailId = ""; probe.clickTag = "";
    probe.pointerUpSeen = false; probe.pointerUpTag = "";
    probe.captureEstablished = null;
    probe.pressTarget = null; probe.pressTargetTag = ""; probe.pressTargetSurvived = null;
    return "armed";
  };
  // WHY A GESTURE WAITS FOR TWO FRAMES FIRST. This view rebuilds its virtualized rows
  // from a scroll event (renderVisibleRows) and from a requestAnimationFrame callback
  // (requestFrameRender). Both are asynchronous, so aiming at a row leaves rebuild
  // work in flight — the aim's own scrollIntoView schedules a scroll event, and a
  // press that focuses a partly-clipped row schedules another. Land one of those
  // between the press and the release and the pressed node is detached, no click is
  // synthesized, and a probe waiting for one waits forever.
  //
  // Two frames rather than one, because the rebuild is scheduled from a callback that
  // itself runs in the first.
  probe.settleRendering = function () {
    return new Promise(function (resolve) {
      requestAnimationFrame(function () {
        requestAnimationFrame(function () { resolve("settled"); });
      });
    });
  };
  probe.gestureOutcome = function () {
    return {
      clickCount: probe.clickCount, clickDetailId: probe.clickDetailId, clickTag: probe.clickTag,
      drawerOpen: probe.drawerIsOpen(), shownDetailId: probe.shownDetailId(),
      panning: probe.isPanning(), captureEstablished: probe.captureEstablished,
      pressTargetTag: probe.pressTargetTag, pressTargetSurvived: probe.pressTargetSurvived
    };
  };
  // THE VIEWPORT RULE, learned the hard way by REQ-336's out-of-suite harness: a
  // bounding rect is not a clickable coordinate until the element is on screen. The
  // first bar of a three-hundred-row board sits near y=1538 in a 900px viewport and a
  // press aimed at its rect lands on HTML. So scroll the chart in, then PROVE the
  // point by asking what is actually under it rather than assuming.
  probe.aimAtARow = function () {
    document.querySelector("#view-timeline .timeline-chart").scrollIntoView({ block: "start" });
    var host = probe.plotHost();
    var hostBox = host.getBoundingClientRect();
    var rows = document.querySelectorAll("#view-timeline .timeline-row");
    var aim = {
      aimed: false, x: 0, y: 0, detailId: "", landedOnTag: "", rowCount: rows.length,
      hostTop: hostBox.top, hostLeft: hostBox.left, hostRight: hostBox.right
    };
    // FULLY inside the host's visible band, with a margin — not merely centred in it.
    // A row clipped at either edge gets scrolled into view by the focus that a trusted
    // press performs on a role=button row, and that scroll fires the host's scroll
    // handler, which rebuilds every row out from under the gesture. Centre-only
    // admitted a row three pixels over the top edge, which is one of the two ways this
    // probe's press and release came apart.
    var edgeMargin = 8;
    for (var index = 0; index < rows.length; index++) {
      var rowBox = rows[index].getBoundingClientRect();
      var y = Math.round(rowBox.top + rowBox.height / 2);
      if (rowBox.top < hostBox.top + edgeMargin || rowBox.bottom > hostBox.bottom - edgeMargin) {
        continue;
      }
      // Left of centre, so a rightward drag of 140px stays on the plot.
      var x = Math.round(hostBox.left + hostBox.width * 0.35);
      var landedOn = document.elementFromPoint(x, y);
      var trigger = landedOn && landedOn.closest ? landedOn.closest("[data-detail-kind]") : null;
      aim.landedOnTag = landedOn && landedOn.tagName ? landedOn.tagName : "";
      if (trigger && trigger.getAttribute("data-detail-id") === rows[index].getAttribute("data-detail-id")) {
        aim.aimed = true;
        aim.x = x;
        aim.y = y;
        aim.detailId = trigger.getAttribute("data-detail-id");
        return aim;
      }
    }
    return aim;
  };
  // THE MUTATION REQ-336 FIXED, reintroduced from outside the board's own source so
  // the probe can be shown to catch it. captureEstablished records whether the engine
  // really granted the capture, so "the drawer stayed shut" can be told apart from
  // "the mutation never took" — which is precisely what happens under synthetic
  // events, and precisely why this check could not be behavioural before.
  probe.captureOnPointerdown = function () {
    var host = probe.plotHost();
    probe.captureListener = function (downEvent) {
      try {
        host.setPointerCapture(downEvent.pointerId);
        probe.captureEstablished = !!host.hasPointerCapture(downEvent.pointerId);
      } catch (captureError) {
        probe.captureEstablished = false;
      }
    };
    host.addEventListener("pointerdown", probe.captureListener);
    return "capture-on-pointerdown installed";
  };
  probe.removeCaptureOnPointerdown = function () {
    probe.plotHost().removeEventListener("pointerdown", probe.captureListener);
    probe.captureListener = null;
    return "capture-on-pointerdown removed";
  };
  // The MIRROR mutation: the board asks for capture and gets nothing. It stubs the
  // element's own method, so capturePanPointer's feature detect still passes and its
  // call still runs; the only thing missing is the capture itself.
  probe.swallowPointerCapture = function () {
    probe.plotHost().setPointerCapture = function () { return undefined; };
    return "pointer capture swallowed";
  };
  probe.restorePointerCapture = function () {
    delete probe.plotHost().setPointerCapture;
    return "pointer capture restored";
  };
  return "installed";
})()`

// timelineRowAim is the press point, PROVEN rather than assumed — see aimAtARow.
type timelineRowAim struct {
	Aimed       bool    `json:"aimed"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	DetailId    string  `json:"detailId"`
	LandedOnTag string  `json:"landedOnTag"`
	RowCount    int     `json:"rowCount"`
	HostTop     float64 `json:"hostTop"`
	HostLeft    float64 `json:"hostLeft"`
	HostRight   float64 `json:"hostRight"`
}

// timelineGestureOutcome is everything one gesture is allowed to be judged on.
type timelineGestureOutcome struct {
	ClickCount          int    `json:"clickCount"`
	ClickDetailId       string `json:"clickDetailId"`
	ClickTag            string `json:"clickTag"`
	DrawerOpen          bool   `json:"drawerOpen"`
	ShownDetailId       string `json:"shownDetailId"`
	Panning             bool   `json:"panning"`
	CaptureEstablished  *bool  `json:"captureEstablished"`
	PressTargetTag      string `json:"pressTargetTag"`
	PressTargetSurvived *bool  `json:"pressTargetSurvived"`
}

// aimAtATimelineRow resets the window, scrolls the chart on screen and returns a press
// point that has been checked against elementFromPoint. A press that reaches no row
// produces "no pan, no drawer", which satisfies half the assertions below while
// measuring nothing, so a failed aim is fatal rather than a result.
func aimAtATimelineRow(
	t *testing.T, session *trustedInputBrowserSession, gestureName string,
) timelineRowAim {
	t.Helper()
	session.evaluateInPage(t, `window.timelineProbe.closeDrawer()`)
	session.evaluateInPage(t, `window.timelineProbe.resetWindow()`)
	// Closing the drawer returns its grid column and resetting the window re-renders,
	// so let both land before measuring anything. Aiming into a page that is still
	// rebuilding produces coordinates for rows that are about to be replaced.
	session.evaluateInPage(t, `window.timelineProbe.settleRendering()`)
	session.evaluateInPage(t, `window.timelineProbe.armGesture()`)
	var aim timelineRowAim
	session.decodeResult(t, "aimAtARow", session.evaluateInPage(t, `window.timelineProbe.aimAtARow()`), &aim)
	if !aim.Aimed {
		t.Fatalf("the %s gesture found no press point fully inside the scroll host over a "+
			"visible row (%d rows rendered, the last point tried landed on %s); every assertion "+
			"after this would be measuring a press that reached no handler",
			gestureName, aim.RowCount, aim.LandedOnTag)
	}
	// The aim itself scrolled the chart into view, which schedules one more rebuild.
	// Settle again and CONFIRM the point still resolves to the row it was measured on,
	// so the gesture below is dispatched at a coordinate that is still true.
	session.evaluateInPage(t, `window.timelineProbe.settleRendering()`)
	var confirmedDetailId string
	session.decodeResult(t, "confirm aim", session.evaluateInPage(t,
		`(function () {
                   var landedOn = document.elementFromPoint(`+
			formatViewportCoordinate(aim.X)+`, `+formatViewportCoordinate(aim.Y)+`);
                   var trigger = landedOn && landedOn.closest ? landedOn.closest("[data-detail-kind]") : null;
                   return trigger ? (trigger.getAttribute("data-detail-id") || "") : "";
                 })()`), &confirmedDetailId)
	if confirmedDetailId != aim.DetailId {
		t.Fatalf("the %s gesture aimed at row %s, but after the page settled that point resolves "+
			"to %q; the chart moved between the aim and the gesture, so nothing dispatched there "+
			"would be measuring the row this trial names",
			gestureName, aim.DetailId, confirmedDetailId)
	}
	return aim
}

// formatViewportCoordinate renders a coordinate for embedding in a page expression.
// The aim reports whole pixels, and this keeps them written that way rather than in the
// exponent form a default float format can reach for.
func formatViewportCoordinate(viewportCoordinate float64) string {
	return strconv.FormatFloat(viewportCoordinate, 'f', -1, 64)
}

// readTimelineGestureOutcome waits for the gesture's RELEASE to have been delivered — a
// real condition, never a duration — and then reads what the page made of it.
//
// The release is the right settle for a drag, because what a drag trial asks about is
// what the release did to the pan, and the pan teardown runs on the way to this watch.
// It is the WRONG settle for a click: the engine synthesizes a click in a later task,
// so a click trial reading here would read a drawer that has not been asked to open
// yet. Click trials use readTimelineClickOutcome, which waits for the click itself.
//
// The condition presumes no outcome, which is what lets it serve the drag that must end
// the pan and the drag that must strand it.
func readTimelineGestureOutcome(
	t *testing.T, session *trustedInputBrowserSession, gestureName string,
) timelineGestureOutcome {
	t.Helper()
	session.waitForPageCondition(t, "the "+gestureName+" release reaches the document",
		`window.timelineProbe.pointerUpSeen`)
	var outcome timelineGestureOutcome
	session.decodeResult(t, "gestureOutcome",
		session.evaluateInPage(t, `window.timelineProbe.gestureOutcome()`), &outcome)
	return outcome
}

// readTimelineClickOutcome is the click trials' reader. It settles in two stages, and
// the second stage is the point.
//
// Stage one waits for the RELEASE, which the engine always delivers. Stage two gives
// the engine a short budget to synthesize the click from it — and if no click arrives,
// says why rather than reporting a timeout. There is exactly one way for a press and a
// release on one point to produce no click: the pressed node was detached in between,
// leaving the two with no common ancestor to fire on. The probe watches for precisely
// that, so the failure names the mechanism instead of the wait.
//
// This replaced a single wait on clickCount, which turned a rebuilt row into a 45s
// timeout that read as "the capture-on-pointerdown click gesture is delivered never
// became true" — a sentence about the harness, in a test about pointer capture.
func readTimelineClickOutcome(
	t *testing.T, session *trustedInputBrowserSession, gestureName string, aim timelineRowAim,
) timelineGestureOutcome {
	t.Helper()
	session.waitForPageCondition(t, "the "+gestureName+" release reaches the document",
		`window.timelineProbe.pointerUpSeen`)
	if !session.pageConditionHoldsWithin(t, "the "+gestureName+" click is synthesized",
		`window.timelineProbe.clickCount > 0`, browserProbeGestureSettleDeadline) {
		var starvedOutcome timelineGestureOutcome
		session.decodeResult(t, "gestureOutcome",
			session.evaluateInPage(t, `window.timelineProbe.gestureOutcome()`), &starvedOutcome)
		if starvedOutcome.PressTargetSurvived != nil && !*starvedOutcome.PressTargetSurvived {
			t.Fatalf("the %s press on row %s landed on <%s>, and that node was REPLACED before "+
				"the release: the chart rebuilt its rows mid-gesture, so the engine had no common "+
				"ancestor to synthesize a click on and produced none. The gesture measured nothing "+
				"— this is the probe's own setup failing, not the board",
				gestureName, aim.DetailId, starvedOutcome.PressTargetTag)
		}
		t.Fatalf("the %s press and release on row %s produced no click within %s, and the pressed "+
			"node <%s> was still in the document at the release; a press and a release on one "+
			"point owe a click, so the engine did not synthesize one for a reason this probe "+
			"cannot see", gestureName, aim.DetailId, browserProbeGestureSettleDeadline,
			starvedOutcome.PressTargetTag)
	}
	var outcome timelineGestureOutcome
	session.decodeResult(t, "gestureOutcome",
		session.evaluateInPage(t, `window.timelineProbe.gestureOutcome()`), &outcome)
	return outcome
}

// TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage pins the regression
// REQ-336 fixed — pointer capture taken on pointerdown — and the one REQ-333 fixed, a
// drag released outside the chart that never tells the host it ended.
//
// WHAT REQ-341 CHANGED HERE, AND WHY. This check shipped for REQ-337 as a STRUCTURAL
// one: it read the generated page's text and required that the pointerdown handler
// contain no setPointerCapture( call, that it call no function whose body has one, and
// that the pointermove handler reach one. It read text because the lane could not do
// anything else. Every probe ran under --dump-dom, with no protocol channel, so the only
// input a probe could produce was a synthetic PointerEvent — and a synthetic pointerId
// establishes no capture at all (on Chromium 141.0.7390.37 setPointerCapture does not
// even throw for it; hasPointerCapture simply stays false afterwards). The whole defect
// is a consequence of a REAL capture retargeting the engine's synthesized click, so it
// was unreachable from the lane, and the structural version carried the residual in its
// own comment: a capture routed through a variable, a method lookup or an eval passed it.
//
// REQ-341 gave the lane a trusted-input transport (browser_probe_test.go), so the
// property is now asserted the way a reader meets it. The residual is closed: these
// gestures do not care HOW a capture is requested, only whether one exists when the
// click is synthesized.
//
// Four trials, in one engine on one page, because each pair is what stops the other
// from being trivially satisfied:
//
//	(1) a trusted press and release on a bar opens the detail drawer on that bar,
//	(2) with capture taken on pointerdown from outside the board's source, the same
//	    gesture opens nothing — which is what proves (1) can fail, and proves it against
//	    the exact regression rather than against a mutation chosen to be easy,
//	(3) a trusted drag released OUTSIDE the chart ends the pan, and
//	(4) with the host's setPointerCapture stubbed out, that release strands the pan —
//	    which is what proves (3) can fail and keeps "no capture anywhere" from passing.
//
// WHAT THE GESTURES ARE SENSITIVE TO, stated because it was measured rather than
// assumed. A capture reaches the click whether it is requested directly, through a
// variable, or from a timer or animation frame scheduled by the pointerdown handler:
// all four spellings fail this test. The last two only fail it because the click
// gesture holds the button down for two frames (clickTrustedMouseOnRow). Dispatched
// back to back, a press and a release are one to two milliseconds apart and the
// renderer delivers both before running any timer, so a deferred capture had not yet
// executed and the trial passed a board that was broken. That was a gap in the probe,
// not a boundary of the defect, and the dwell closes it.
//
// It is still a press-and-release on ONE point. A capture requested later than two
// frames after the press — a network round trip, a long task — would be established
// after this release and pass, and so would a capture taken on a gesture this test
// does not perform.
func TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}

	// The REAL page, unedited: this transport drives it from outside, so unlike the
	// dump-dom probes it needs no script injected to carry a result back.
	session := startTrustedInputBrowserSession(t, "timeline pointer capture", siteDirectory,
		string(indexBytes), "--window-size=1600,900")
	session.waitForPageCondition(t, "the board's view switcher renders",
		`document.querySelector('[data-view-target="timeline"]')`)
	session.evaluateInPage(t,
		`(document.querySelector('[data-view-target="timeline"]').click(), "switched")`)
	session.waitForPageCondition(t, "the Timeline draws its rows",
		`document.querySelectorAll("#view-timeline .timeline-row").length > 0`)
	session.evaluateInPage(t, timelinePointerCaptureProbeHelpers)

	// (1) THE CLICK, WITH TRUSTED INPUT. The hint under the chart says "Click a row for
	// its full detail"; this is that sentence, executed.
	cleanAim := aimAtATimelineRow(t, session, "clean click")
	clickTrustedMouseOnRow(t, session, cleanAim)
	cleanClick := readTimelineClickOutcome(t, session, "clean click", cleanAim)
	if !cleanClick.DrawerOpen {
		t.Errorf("a trusted click on row %s did not open the detail drawer; the click was "+
			"delivered to %s and its nearest [data-detail-kind] was %q. Capture taken before the "+
			"pan engages retargets the synthesized click to the capturing element, which is how "+
			"every mouse click in the chart stopped opening anything",
			cleanAim.DetailId, cleanClick.ClickTag, cleanClick.ClickDetailId)
	}
	if cleanClick.ShownDetailId != cleanAim.DetailId {
		t.Errorf("a trusted click on row %s opened the drawer on %q; a click has to open the row "+
			"it landed on", cleanAim.DetailId, cleanClick.ShownDetailId)
	}
	if cleanClick.Panning {
		t.Errorf("a still trusted press on row %s left the chart in the panning state; a press "+
			"that has not moved is not a drag", cleanAim.DetailId)
	}

	// (2) THE MUTATION. Capture on pointerdown, installed from here rather than by
	// editing the board, and the SAME gesture must now open nothing. Without this the
	// assertion above is only a claim that clicking works today.
	session.evaluateInPage(t, `window.timelineProbe.captureOnPointerdown()`)
	capturedAim := aimAtATimelineRow(t, session, "capture-on-pointerdown click")
	clickTrustedMouseOnRow(t, session, capturedAim)
	capturedClick := readTimelineClickOutcome(t, session, "capture-on-pointerdown click", capturedAim)
	session.evaluateInPage(t, `window.timelineProbe.removeCaptureOnPointerdown()`)
	if capturedClick.CaptureEstablished == nil {
		t.Fatal("the capture-on-pointerdown mutation never saw a pointerdown on the scroll host, " +
			"so the trial below proves nothing about the probe's sensitivity")
	}
	if !*capturedClick.CaptureEstablished {
		t.Fatal("the engine refused the capture-on-pointerdown mutation, so this trial is the " +
			"synthetic-event dead end all over again: it cannot tell a probe that catches the " +
			"regression from one that cannot see it")
	}
	if capturedClick.DrawerOpen {
		t.Errorf("with pointer capture taken on pointerdown the drawer still opened on %q; the "+
			"clean trial above therefore does not pin REQ-336's regression, because the probe "+
			"cannot tell a retargeted click from an ordinary one", capturedClick.ShownDetailId)
	}
	if capturedClick.ClickDetailId != "" {
		t.Errorf("with capture on pointerdown the synthesized click still found the trigger %q; "+
			"capture is supposed to retarget it to the capturing element, and if it does not, "+
			"this whole mutation is measuring something else", capturedClick.ClickDetailId)
	}

	// (3) THE RELEASE OUTSIDE THE CHART. Capture at the engage is what guarantees the
	// host hears the end of a drag: Chromium suppresses the boundary events while a
	// button is held, so without it nothing reaches the host and the grab cursor stays
	// on for the rest of the session (REQ-333).
	//
	// Measured from the PAN ENGAGING and from the is-panning state, never from the axis
	// text: a drag that clamps at the window bound leaves every label identical.
	engagedDrag := aimAtATimelineRow(t, session, "release outside")
	releaseY := engagedDrag.HostTop - 40
	if releaseY < 1 {
		t.Fatalf("the scroll host sits at y=%g, so there is no point above it inside the "+
			"viewport to release at; this trial cannot be run from here", engagedDrag.HostTop)
	}
	panningAfterRelease := driveTimelineDragReleasedOutside(t, session, engagedDrag, releaseY)
	if panningAfterRelease {
		t.Error("a trusted drag released above the chart left it in the panning state; capture " +
			"taken when the pan engages is what makes the release a fact rather than a hope, and " +
			"without it the grab cursor stays on for the rest of the session")
	}

	// (4) THE MIRROR MUTATION. (1) and (2) are both satisfied by a board that captures
	// nowhere at all, which is REQ-333's bug. Swallowing the host's own
	// setPointerCapture leaves capturePanPointer's feature detect and call intact and
	// removes only the capture — and the same drag must then strand the pan.
	session.evaluateInPage(t, `window.timelineProbe.swallowPointerCapture()`)
	uncapturedDrag := aimAtATimelineRow(t, session, "release outside without capture")
	strandedPanning := driveTimelineDragReleasedOutside(t, session, uncapturedDrag,
		uncapturedDrag.HostTop-40)
	session.evaluateInPage(t, `window.timelineProbe.restorePointerCapture()`)
	if !strandedPanning {
		t.Error("with the scroll host's setPointerCapture swallowed, a drag released above the " +
			"chart still ended cleanly; trial (3) is then passing for some other reason and does " +
			"not pin the capture at all")
	}

	t.Logf("trusted click on %s: %d click(s), delivered to <%s>, trigger %q, drawer open %v "+
		"showing %q. With capture on pointerdown: %d click(s), trigger %q, drawer open %v. "+
		"Drag released above the host left the chart panning: %v with capture, %v with capture "+
		"swallowed.",
		cleanAim.DetailId, cleanClick.ClickCount, cleanClick.ClickTag, cleanClick.ClickDetailId,
		cleanClick.DrawerOpen, cleanClick.ShownDetailId,
		capturedClick.ClickCount, capturedClick.ClickDetailId, capturedClick.DrawerOpen,
		panningAfterRelease, strandedPanning)
}

// clickTrustedMouseOnRow presses and releases on a row with a real DWELL between the
// two, the way a hand does.
//
// The dwell is not politeness, it is coverage. Dispatched back to back the press and
// the release are one to two milliseconds apart — measured — and the renderer delivers
// both before it runs a single timer callback. A capture requested from the pointerdown
// handler but DEFERRED (a setTimeout, an animation frame) therefore had not executed at
// all when the click was synthesized, so the gesture sailed past a board that would
// break for any human: with two frames of dwell the same deferred capture is granted
// 16ms after the press, the release retargets to the scroll host, and the drawer stops
// opening. That gap was found by measurement and is the reason this is not a bare pair.
//
// Two frames, not a sleep: it is the same wait on the page's own scheduling used before
// the gesture, and no wall-clock duration is assumed to mean anything.
func clickTrustedMouseOnRow(
	t *testing.T, session *trustedInputBrowserSession, aim timelineRowAim,
) {
	t.Helper()
	session.pressTrustedMouseAt(t, aim.X, aim.Y)
	session.evaluateInPage(t, `window.timelineProbe.settleRendering()`)
	session.releaseTrustedMouseAt(t, aim.X, aim.Y)
}

// driveTimelineDragReleasedOutside presses on a row, drags far enough to engage the pan,
// then releases at a point ABOVE the scroll host, and reports whether the chart is still
// in the panning state afterwards.
//
// The engage is asserted, not assumed: a release that ends a pan which never started
// tells nobody anything.
func driveTimelineDragReleasedOutside(
	t *testing.T, session *trustedInputBrowserSession, aim timelineRowAim, releaseY float64,
) bool {
	t.Helper()
	session.pressTrustedMouseAt(t, aim.X, aim.Y)
	// Two moves: the first clears the 4px pan threshold, the second is the drag the
	// reader would see.
	session.dragTrustedMouseTo(t, aim.X+6, aim.Y)
	session.dragTrustedMouseTo(t, aim.X+140, aim.Y)
	session.waitForPageCondition(t, "the drag engages the pan", `window.timelineProbe.isPanning()`)

	releaseX := aim.X + 140
	session.dragTrustedMouseTo(t, releaseX, releaseY)
	session.releaseTrustedMouseAt(t, releaseX, releaseY)
	outcome := readTimelineGestureOutcome(t, session, "drag released outside")
	return outcome.Panning
}

// TestBrowserBehaviorTimelineRowListIsOneTabStop pins REQ-338's keyboard contract: the row
// list is a single Tab stop with arrow-key movement between rows (a roving tabindex),
// rather than one Tab stop per row.
//
// "Tab escapes the list in one press" is asserted as the property that produces it —
// exactly one row is tabbable and every other row is explicitly not — because Tab's
// focus movement is a default action on a TRUSTED key event and this lane dispatches
// synthetic ones. The movement half IS behavioural: the handler calls focus() itself, so a
// synthetic ArrowDown really does move focus, and that is what the probe measures.
//
// It also pins the two contracts a roving index could quietly break: Left/Right must still
// pan the window without moving row focus (REQ-333), and Enter on the focused row must
// still open the detail drawer (REQ-333, restored by REQ-336).
func TestBrowserBehaviorTimelineRowListIsOneTabStop(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function plotHost() { return document.querySelector("#view-timeline .timeline-scroll"); }
function rowNodes() { return Array.prototype.slice.call(plotHost().querySelectorAll("[data-row-index]")); }
function tabStopState() {
  var rows = rowNodes();
  var tabbable = rows.filter(function (row) { return row.getAttribute("tabindex") === "0"; });
  var notTabbable = rows.filter(function (row) { return row.getAttribute("tabindex") === "-1"; });
  var active = document.activeElement;
  var drawer = document.getElementById("detail-drawer");
  var shown = document.getElementById("detail-id");
  return {
    href: location.href,
    rowCount: rows.length,
    tabbableCount: tabbable.length,
    notTabbableCount: notTabbable.length,
    tabbableRowId: tabbable.length === 1 ? (tabbable[0].getAttribute("data-detail-id") || "") : "",
    focusedRowId: active && active.getAttribute ? (active.getAttribute("data-detail-id") || "") : "",
    readout: document.getElementById("timeline-range-readout").textContent,
    drawerOpen: drawer ? drawer.hidden === false : false,
    shownDetailId: shown ? shown.textContent.trim() : ""
  };
}

function pressKey(keyName) {
  var target = document.activeElement || plotHost();
  target.dispatchEvent(new KeyboardEvent("keydown", {
    key: keyName, bubbles: true, cancelable: true, composed: true
  }));
}

window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var probe = {};
      // Off the bounds so ArrowRight has somewhere to pan to.
      document.getElementById("timeline-zoom-in").click();
      var rows = rowNodes();
      probe.initial = tabStopState();
      if (rows.length >= 2) {
        probe.firstRowId = rows[0].getAttribute("data-detail-id") || "";
        probe.secondRowId = rows[1].getAttribute("data-detail-id") || "";
        rows[0].focus();
        probe.afterFocusingFirstRow = tabStopState();
        pressKey("ArrowDown");
        probe.afterArrowDown = tabStopState();
        pressKey("ArrowUp");
        probe.afterArrowUp = tabStopState();
        pressKey("ArrowRight");
        probe.afterArrowRight = tabStopState();
        pressKey("Enter");
        probe.afterEnter = tabStopState();
      }
      document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify(probe);
    }, 700);
  }, 200);
});
</script>
</body>`

	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the probe script before")
	}
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline row tab stops", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")

	type tabStopState struct {
		Href             string `json:"href"`
		RowCount         int    `json:"rowCount"`
		TabbableCount    int    `json:"tabbableCount"`
		NotTabbableCount int    `json:"notTabbableCount"`
		TabbableRowId    string `json:"tabbableRowId"`
		FocusedRowId     string `json:"focusedRowId"`
		Readout          string `json:"readout"`
		DrawerOpen       bool   `json:"drawerOpen"`
		ShownDetailId    string `json:"shownDetailId"`
	}
	var rovingResult struct {
		FirstRowId            string       `json:"firstRowId"`
		SecondRowId           string       `json:"secondRowId"`
		Initial               tabStopState `json:"initial"`
		AfterFocusingFirstRow tabStopState `json:"afterFocusingFirstRow"`
		AfterArrowDown        tabStopState `json:"afterArrowDown"`
		AfterArrowUp          tabStopState `json:"afterArrowUp"`
		AfterArrowRight       tabStopState `json:"afterArrowRight"`
		AfterEnter            tabStopState `json:"afterEnter"`
	}
	if decodeError := json.Unmarshal(probeOutput, &rovingResult); decodeError != nil {
		t.Fatalf("decode timeline row tab-stop behavior: %v (output %q)", decodeError, probeOutput)
	}
	if rovingResult.Initial.Href == "" || !strings.HasSuffix(rovingResult.Initial.Href, "probe.html") {
		t.Fatalf("measured on %q, not the probe page", rovingResult.Initial.Href)
	}
	// Vacuity guard: one row cannot show a roving index, and zero rows cannot show
	// anything. The generated board carries hundreds, so this is a "the probe never
	// reached the chart" failure rather than a real board state.
	if rovingResult.Initial.RowCount < 2 {
		t.Fatalf("the probe rendered %d timeline rows, so nothing below measures a row list",
			rovingResult.Initial.RowCount)
	}
	if rovingResult.FirstRowId == "" || rovingResult.SecondRowId == "" ||
		rovingResult.FirstRowId == rovingResult.SecondRowId {
		t.Fatalf("the probe could not name two distinct rows (%q and %q), so focus movement below "+
			"cannot be told from focus standing still",
			rovingResult.FirstRowId, rovingResult.SecondRowId)
	}

	// (a) ONE Tab stop. Tab's own movement is a trusted-input default action this lane
	// cannot dispatch, so the assertion is the property that produces it: exactly one row
	// is reachable by Tab and every other row is explicitly skipped.
	if rovingResult.Initial.TabbableCount != 1 {
		t.Errorf("the rendered row list has %d rows with tabindex=0 out of %d rows; a keyboard "+
			"reader pays one Tab press per tabbable row to get past the chart",
			rovingResult.Initial.TabbableCount, rovingResult.Initial.RowCount)
	}
	if rovingResult.Initial.NotTabbableCount != rovingResult.Initial.RowCount-1 {
		t.Errorf("%d of %d rendered rows carry tabindex=-1; every row other than the roving one "+
			"must be explicitly skipped, or the browser's own default takes over",
			rovingResult.Initial.NotTabbableCount, rovingResult.Initial.RowCount)
	}

	// (b) The tabbable row FOLLOWS focus, which is what makes Tab return to where the
	// reader left off rather than to the top of the list.
	if rovingResult.AfterFocusingFirstRow.TabbableRowId != rovingResult.FirstRowId {
		t.Errorf("after focusing the first row the tabbable row is %q, want %q",
			rovingResult.AfterFocusingFirstRow.TabbableRowId, rovingResult.FirstRowId)
	}

	// (c) Down and Up MOVE focus between rows, and the roving index goes with it.
	if rovingResult.AfterArrowDown.FocusedRowId != rovingResult.SecondRowId {
		t.Errorf("ArrowDown left focus on %q; it should move to the next row %q",
			rovingResult.AfterArrowDown.FocusedRowId, rovingResult.SecondRowId)
	}
	if rovingResult.AfterArrowDown.TabbableRowId != rovingResult.SecondRowId {
		t.Errorf("after ArrowDown the tabbable row is %q, want the newly focused %q",
			rovingResult.AfterArrowDown.TabbableRowId, rovingResult.SecondRowId)
	}
	if rovingResult.AfterArrowDown.TabbableCount != 1 {
		t.Errorf("after ArrowDown %d rows are tabbable; moving the roving index must not leave two",
			rovingResult.AfterArrowDown.TabbableCount)
	}
	if rovingResult.AfterArrowUp.FocusedRowId != rovingResult.FirstRowId {
		t.Errorf("ArrowUp left focus on %q; it should move back to %q",
			rovingResult.AfterArrowUp.FocusedRowId, rovingResult.FirstRowId)
	}

	// (d) REQ-333's contract: Left/Right still pan, and panning does not move row focus.
	if rovingResult.AfterArrowRight.Readout == rovingResult.AfterArrowUp.Readout {
		t.Errorf("ArrowRight left the window at %s; arrow panning is REQ-333's contract and the "+
			"roving index must not take the horizontal keys",
			rovingResult.AfterArrowRight.Readout)
	}
	if rovingResult.AfterArrowRight.FocusedRowId != rovingResult.FirstRowId {
		t.Errorf("ArrowRight moved row focus to %q; the horizontal keys pan, they do not rove",
			rovingResult.AfterArrowRight.FocusedRowId)
	}

	// (e) And Enter on the focused row still opens the drawer.
	if !rovingResult.AfterEnter.DrawerOpen {
		t.Error("Enter on the focused row did not open the detail drawer; rows advertise " +
			"role=button, so they owe the activation")
	}
	if rovingResult.AfterEnter.ShownDetailId != rovingResult.FirstRowId {
		t.Errorf("Enter opened the drawer on %q, want the focused row %q",
			rovingResult.AfterEnter.ShownDetailId, rovingResult.FirstRowId)
	}
}

// Timeline grouping is rendered by the shipped board, not by a probe fixture.
// The table exposes every window-listed group while the SVG stays virtualized,
// so one run can prove both the complete grouping and bounded viewport nodes.
func TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function timelineGroupingSnapshot() {
  var host = document.getElementById("timeline-scroll");
  var tableBody = document.getElementById("timeline-table-body");
  var payload = window.queueKanbanBoardData || {};
  var requests = payload.requests || {};
  var groups = [];
  var currentGroup = null;
  var columnHeaderIds = Array.prototype.map.call(
    document.querySelectorAll("#view-timeline .timeline-table thead th"),
    function (header) { return header.id; });
  Array.prototype.forEach.call(tableBody.children, function (tableRow) {
    var groupLabel = tableRow.getAttribute("data-timeline-table-group");
    if (groupLabel) {
      var groupHeading = tableRow.querySelector("th");
      currentGroup = {
        label: groupLabel,
        headerId: groupHeading ? groupHeading.id : "",
        statedCount: Number(tableRow.getAttribute("data-group-count")),
        metricText: tableRow.textContent.trim(),
        ids: [],
        joinedUserRequestIds: [],
        cellHeaders: [],
        cellTags: []
      };
      groups.push(currentGroup);
      return;
    }
    var requestId = tableRow.getAttribute("data-timeline-table-request");
    if (currentGroup && requestId) {
      currentGroup.ids.push(requestId);
      currentGroup.joinedUserRequestIds.push((requests[requestId] || {}).userRequestId || "");
      currentGroup.cellHeaders.push(Array.prototype.map.call(tableRow.children, function (cell) {
        return (cell.getAttribute("headers") || "").split(/\s+/).filter(Boolean);
      }));
      currentGroup.cellTags.push(Array.prototype.map.call(tableRow.children, function (cell) {
        return cell.tagName;
      }));
    }
  });
  var headers = Array.prototype.slice.call(host.querySelectorAll("[data-timeline-group]"));
  var rows = Array.prototype.slice.call(host.querySelectorAll("[data-row-index]"));
  return {
    href: location.href,
    columnHeaderIds: columnHeaderIds,
    groups: groups,
    visibleRowCount: rows.length,
    visibleHeaderCount: headers.length,
    headerTabStopCount: headers.filter(function (header) { return header.hasAttribute("tabindex"); }).length,
    tabbableRowCount: rows.filter(function (row) { return row.getAttribute("tabindex") === "0"; }).length,
    renderedNodeCount: headers.length + rows.length,
    svgHeight: Number((host.querySelector("svg") || {}).getAttribute && host.querySelector("svg").getAttribute("height")) || 0,
    viewportHeight: host.clientHeight
  };
}

function afterTimelineFrames(callback) {
  setTimeout(callback, 50);
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var table = document.querySelector("#view-timeline .timeline-table");
      table.open = true;
      table.dispatchEvent(new Event("toggle"));
      document.getElementById("timeline-zoom-fit").click();
      afterTimelineFrames(function () {
        var fitAll = timelineGroupingSnapshot();
        document.querySelector('[data-timeline-period="day"]').click();
        afterTimelineFrames(function () {
          document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
            fitAll: fitAll,
            day: timelineGroupingSnapshot()
          });
        });
      });
    }, 900);
  }, 200);
});
</script>
</body>`
	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the grouping probe before")
	}
	probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline user-request grouping", siteDirectory,
		pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000")
	type browserGroup struct {
		Label                string       `json:"label"`
		HeaderID             string       `json:"headerId"`
		StatedCount          int          `json:"statedCount"`
		MetricText           string       `json:"metricText"`
		Ids                  []string     `json:"ids"`
		JoinedUserRequestIds []string     `json:"joinedUserRequestIds"`
		CellHeaders          [][][]string `json:"cellHeaders"`
		CellTags             [][]string   `json:"cellTags"`
	}
	type groupingSnapshot struct {
		Href               string         `json:"href"`
		ColumnHeaderIDs    []string       `json:"columnHeaderIds"`
		Groups             []browserGroup `json:"groups"`
		VisibleRowCount    int            `json:"visibleRowCount"`
		VisibleHeaderCount int            `json:"visibleHeaderCount"`
		HeaderTabStopCount int            `json:"headerTabStopCount"`
		TabbableRowCount   int            `json:"tabbableRowCount"`
		RenderedNodeCount  int            `json:"renderedNodeCount"`
		SvgHeight          float64        `json:"svgHeight"`
		ViewportHeight     float64        `json:"viewportHeight"`
	}
	var groupingResult struct {
		FitAll groupingSnapshot `json:"fitAll"`
		Day    groupingSnapshot `json:"day"`
	}
	if decodeError := json.Unmarshal(probeOutput, &groupingResult); decodeError != nil {
		t.Fatalf("decode timeline user-request grouping: %v (output %q)", decodeError, probeOutput)
	}
	if groupingResult.FitAll.Href == "" || !strings.HasSuffix(groupingResult.FitAll.Href, "probe.html") {
		t.Fatalf("measured on %q, not the probe page", groupingResult.FitAll.Href)
	}
	if groupingResult.FitAll.VisibleRowCount == 0 {
		t.Fatal("the generated board rendered no Timeline REQ rows, so the grouping assertion is vacuous")
	}
	if len(groupingResult.FitAll.Groups) < 2 || groupingResult.FitAll.VisibleHeaderCount == 0 {
		t.Fatalf("the Timeline rendered %d table groups and %d visible headers; the fixture must exercise grouping",
			len(groupingResult.FitAll.Groups), groupingResult.FitAll.VisibleHeaderCount)
	}
	for _, snapshot := range []groupingSnapshot{groupingResult.FitAll, groupingResult.Day} {
		if len(snapshot.ColumnHeaderIDs) != 6 {
			t.Fatalf("grouped Timeline table has column header ids %v, want six explicit ids", snapshot.ColumnHeaderIDs)
		}
		seenGroupHeaderIDs := map[string]bool{}
		for groupIndex, group := range snapshot.Groups {
			if group.HeaderID == "" || seenGroupHeaderIDs[group.HeaderID] {
				t.Errorf("group %q has missing or duplicate stable header id %q", group.Label, group.HeaderID)
			}
			seenGroupHeaderIDs[group.HeaderID] = true
			if len(group.Ids) != group.StatedCount {
				t.Errorf("%s states %d listed REQs but the grouped table contains %d",
					group.Label, group.StatedCount, len(group.Ids))
			}
			if !strings.Contains(group.MetricText, "Elapsed") ||
				!strings.Contains(group.MetricText, "Accepted work") ||
				!strings.Contains(group.MetricText, "listed REQ") {
				t.Errorf("%s header omits an explicit metric: %q", group.Label, group.MetricText)
			}
			for memberIndex, joinedUserRequestId := range group.JoinedUserRequestIds {
				wantLabel := joinedUserRequestId
				if wantLabel == "" {
					wantLabel = "No UR recorded"
				}
				if group.Label != wantLabel {
					t.Errorf("group %q contains %s joined to %q; the client-side requests join put it under the wrong header",
						group.Label, group.Ids[memberIndex], joinedUserRequestId)
				}
				if len(group.CellHeaders[memberIndex]) != len(snapshot.ColumnHeaderIDs) {
					t.Errorf("%s in %s has %d table cells, want %d", group.Ids[memberIndex], group.Label,
						len(group.CellHeaders[memberIndex]), len(snapshot.ColumnHeaderIDs))
					continue
				}
				for cellIndex, headerTokens := range group.CellHeaders[memberIndex] {
					wantHeaders := []string{group.HeaderID, snapshot.ColumnHeaderIDs[cellIndex]}
					if !reflect.DeepEqual(headerTokens, wantHeaders) {
						t.Errorf("%s cell %d headers = %v, want exactly its own group and column headers %v",
							group.Ids[memberIndex], cellIndex, headerTokens, wantHeaders)
					}
				}
				if len(group.CellTags[memberIndex]) != 6 || group.CellTags[memberIndex][0] != "TH" {
					t.Errorf("%s member cells = %v, want its REQ cell to remain a row header", group.Ids[memberIndex],
						group.CellTags[memberIndex])
				}
			}
			if group.Label == "No UR recorded" && groupIndex != len(snapshot.Groups)-1 {
				t.Errorf("No UR recorded is group %d of %d; it must be last", groupIndex+1, len(snapshot.Groups))
			}
		}
		if snapshot.HeaderTabStopCount != 0 {
			t.Errorf("%d group headers entered the Tab order; Up/Down must move between REQs only",
				snapshot.HeaderTabStopCount)
		}
		if snapshot.TabbableRowCount != 1 {
			t.Errorf("grouped Timeline has %d tabbable REQ rows, want exactly one", snapshot.TabbableRowCount)
		}
		if snapshot.RenderedNodeCount >= 100 {
			t.Errorf("the grouped Timeline rendered %d header/member nodes in one viewport; virtualization is unbounded",
				snapshot.RenderedNodeCount)
		}
	}
	if groupingResult.FitAll.SvgHeight <= groupingResult.FitAll.ViewportHeight {
		t.Errorf("Fit-all SVG height %.0f does not exceed the %.0fpx viewport, so the probe never exercised virtual scrolling",
			groupingResult.FitAll.SvgHeight, groupingResult.FitAll.ViewportHeight)
	}
	fitAllIds := []string{}
	for _, group := range groupingResult.FitAll.Groups {
		fitAllIds = append(fitAllIds, group.Ids...)
	}
	dayIds := []string{}
	for _, group := range groupingResult.Day.Groups {
		dayIds = append(dayIds, group.Ids...)
	}
	if reflect.DeepEqual(fitAllIds, dayIds) {
		t.Errorf("Fit all and Day listed the same %d REQs; the grouping was not rebuilt after the window changed",
			len(fitAllIds))
	}
}

func TestBrowserBehaviorTimelineGroupHeadersReadInBothThemes(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)
	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="timeline"]').click();
    setTimeout(function () {
      var header = document.querySelector(".timeline-group-header");
      var headerRect = header && header.querySelector(".timeline-group-header-fill");
      var label = header && header.querySelector(".timeline-group-label");
      var metrics = header && header.querySelector(".timeline-group-metrics");
      var rectBox = headerRect ? headerRect.getBBox() : {};
      var labelBox = label ? label.getBBox() : {};
      var metricBox = metrics ? metrics.getBBox() : {};
      document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
        href: location.href,
        scheme: matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
        headerFill: headerRect ? getComputedStyle(headerRect).fill : "",
        labelFill: label ? getComputedStyle(label).fill : "",
        metricFill: metrics ? getComputedStyle(metrics).fill : "",
        rectY: rectBox.y || 0,
        rectHeight: rectBox.height || 0,
        labelY: labelBox.y || 0,
        labelHeight: labelBox.height || 0,
        metricY: metricBox.y || 0,
        metricHeight: metricBox.height || 0
      });
    }, 700);
  }, 200);
});
</script>
</body>`
	pageHTML := strings.Replace(indexHTML, "</body>", probeScript, 1)
	if pageHTML == indexHTML {
		t.Fatal("the generated page has no </body> to inject the header probe before")
	}

	for _, scheme := range []struct {
		name string
		flag string
	}{
		{name: "light", flag: "--blink-settings=preferredColorScheme=1"},
		{name: "dark", flag: "--blink-settings=preferredColorScheme=0"},
	} {
		t.Run(scheme.name, func(t *testing.T) {
			probeOutput := runBrowserBehaviorProbeInDirectory(t, "timeline group header "+scheme.name,
				siteDirectory, pageHTML, "--window-size=1600,900", "--virtual-time-budget=30000", scheme.flag)
			var measurement struct {
				Href         string  `json:"href"`
				Scheme       string  `json:"scheme"`
				HeaderFill   string  `json:"headerFill"`
				LabelFill    string  `json:"labelFill"`
				MetricFill   string  `json:"metricFill"`
				RectY        float64 `json:"rectY"`
				RectHeight   float64 `json:"rectHeight"`
				LabelY       float64 `json:"labelY"`
				LabelHeight  float64 `json:"labelHeight"`
				MetricY      float64 `json:"metricY"`
				MetricHeight float64 `json:"metricHeight"`
			}
			if decodeError := json.Unmarshal(probeOutput, &measurement); decodeError != nil {
				t.Fatalf("decode %s group header: %v (output %q)", scheme.name, decodeError, probeOutput)
			}
			if measurement.Href == "" || !strings.HasSuffix(measurement.Href, "probe.html") {
				t.Fatalf("measured on %q, not the probe page", measurement.Href)
			}
			if measurement.Scheme != scheme.name {
				t.Fatalf("asked for %s and the browser resolved %s", scheme.name, measurement.Scheme)
			}
			headerLuminance, headerKnown := relativeLuminanceOfCSSColour(measurement.HeaderFill)
			labelLuminance, labelKnown := relativeLuminanceOfCSSColour(measurement.LabelFill)
			metricLuminance, metricKnown := relativeLuminanceOfCSSColour(measurement.MetricFill)
			if !headerKnown || !labelKnown || !metricKnown {
				t.Fatalf("%s header colours did not resolve: header %q, label %q, metrics %q",
					scheme.name, measurement.HeaderFill, measurement.LabelFill, measurement.MetricFill)
			}
			if ratio := contrastRatio(labelLuminance, headerLuminance); ratio < 4.5 {
				t.Errorf("%s group label contrast %.2f:1, want at least 4.5:1", scheme.name, ratio)
			}
			if ratio := contrastRatio(metricLuminance, headerLuminance); ratio < 4.5 {
				t.Errorf("%s group metrics contrast %.2f:1, want at least 4.5:1", scheme.name, ratio)
			}
			if measurement.RectHeight != 34 || measurement.LabelHeight <= 0 || measurement.MetricHeight <= 0 {
				t.Errorf("%s header geometry rect %.1fpx, label %.1fpx, metrics %.1fpx; expected a 34px fixed header with rendered text",
					scheme.name, measurement.RectHeight, measurement.LabelHeight, measurement.MetricHeight)
			}
			if measurement.LabelY+measurement.LabelHeight > measurement.MetricY ||
				measurement.MetricY+measurement.MetricHeight > measurement.RectY+measurement.RectHeight {
				t.Errorf("%s header text overlaps or clips: rect y %.1f..%.1f, label %.1f..%.1f, metrics %.1f..%.1f",
					scheme.name, measurement.RectY, measurement.RectY+measurement.RectHeight,
					measurement.LabelY, measurement.LabelY+measurement.LabelHeight,
					measurement.MetricY, measurement.MetricY+measurement.MetricHeight)
			}
			t.Logf("%s Timeline group header: label %.2f:1, metrics %.2f:1, 34px fixed height",
				scheme.name,
				contrastRatio(labelLuminance, headerLuminance),
				contrastRatio(metricLuminance, headerLuminance))
		})
	}
}
