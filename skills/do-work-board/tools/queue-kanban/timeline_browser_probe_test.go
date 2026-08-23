package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
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
        var row = document.querySelector("#view-timeline .timeline-row");
        var box = row.getBoundingClientRect();
        var startX = Math.round(box.left + box.width / 2);
        var y = Math.round(box.top + box.height / 2);
        var before = windowReadout();
        var during = pressDragRelease(row, startX, y, moveOffsets);
        trials[name] = {
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
