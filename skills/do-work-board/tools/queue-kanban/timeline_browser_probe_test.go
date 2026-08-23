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
// The window, in milliseconds, read back off the readout the reader sees rather
// than out of any internal state.
function windowSpanMs() {
  var text = document.getElementById("timeline-range-readout").textContent || "";
  var match = text.match(/^(\S+ \S+) UTC → (\S+ \S+) UTC$/);
  if (!match) { return null; }
  return Date.parse(match[2].replace(" ", "T") + "Z") - Date.parse(match[1].replace(" ", "T") + "Z");
}
function toolbarState(label) {
  var disabled = {};
  ["timeline-zoom-in", "timeline-zoom-out", "timeline-zoom-fit",
   "timeline-period-prev", "timeline-period-next"].forEach(function (buttonId) {
    disabled[buttonId] = !!document.getElementById(buttonId).disabled;
  });
  var rowsHost = document.querySelector("#view-timeline .timeline-scroll");
  return {
    label: label,
    href: location.href,
    readout: document.getElementById("timeline-range-readout").textContent,
    spanMs: windowSpanMs(),
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
      // The current week, then one step forward — which used to land in the
      // cosmetic bound padding with nothing drawn in it.
      document.querySelector('[data-timeline-period="week"]').click();
      probe.currentWeek = toolbarState("currentWeek");
      document.getElementById("timeline-period-next").click();
      probe.afterStepPastTheData = toolbarState("afterStepPastTheData");
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
		SpanMs        *float64        `json:"spanMs"`
		NowRuleDrawn  bool            `json:"nowRuleDrawn"`
		DrawnSegments int             `json:"drawnSegments"`
		Disabled      map[string]bool `json:"disabled"`
	}
	var landingResult struct {
		Fitted               toolbarState `json:"fitted"`
		AfterNow             toolbarState `json:"afterNow"`
		AfterNowThenZoomIn   toolbarState `json:"afterNowThenZoomIn"`
		CurrentWeek          toolbarState `json:"currentWeek"`
		AfterStepPastTheData toolbarState `json:"afterStepPastTheData"`
		FilteredFit          toolbarState `json:"filteredFit"`
		FilteredSummary      string       `json:"filteredSummary"`
	}
	if decodeError := json.Unmarshal(probeOutput, &landingResult); decodeError != nil {
		t.Fatalf("decode timeline landing behavior: %v (output %q)", decodeError, probeOutput)
	}

	states := []toolbarState{
		landingResult.Fitted, landingResult.AfterNow, landingResult.AfterNowThenZoomIn,
		landingResult.CurrentWeek, landingResult.AfterStepPastTheData, landingResult.FilteredFit,
	}
	for _, state := range states {
		if !strings.HasSuffix(state.Href, "/probe.html") {
			t.Fatalf("the %s state was measured on %q, not the probe page", state.Label, state.Href)
		}
		if state.SpanMs == nil {
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
	if !landingResult.CurrentWeek.Disabled["timeline-period-next"] {
		t.Errorf("the step-forward arrow is enabled on the current week (%s), whose next period "+
			"exists only inside the cosmetic bound padding", landingResult.CurrentWeek.Readout)
	}
	if landingResult.AfterStepPastTheData.Readout != landingResult.CurrentWeek.Readout {
		t.Errorf("pressing the step-forward arrow moved the window from %s to %s, past everything "+
			"drawn", landingResult.CurrentWeek.Readout, landingResult.AfterStepPastTheData.Readout)
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

// pointerCapturingFunctionNames lists every named function in the generated page whose
// OWN body requests pointer capture. It is derived from the page rather than hand-listed
// because the point of the check below is to survive a regression that routes the
// request through a function nobody has written yet.
//
// The trailing paren in the match is load-bearing, and it is the trap REQ-333 fell into:
// a search for the bare name also matches the `typeof scrollHost.setPointerCapture ===
// "function"` feature detect, so it passes with the call itself deleted.
func pointerCapturingFunctionNames(t *testing.T, pageSource string) []string {
	t.Helper()
	const declarationToken = "function "
	var capturingNames []string
	searchOffset := 0
	for {
		relativeIndex := strings.Index(pageSource[searchOffset:], declarationToken)
		if relativeIndex == -1 {
			break
		}
		declarationIndex := searchOffset + relativeIndex
		searchOffset = declarationIndex + len(declarationToken)
		nameEnd := strings.IndexByte(pageSource[searchOffset:], '(')
		if nameEnd == -1 {
			break
		}
		functionName := strings.TrimSpace(pageSource[searchOffset : searchOffset+nameEnd])
		// Anonymous function expressions (`function (event) {`) have no name to call, so
		// nothing can route a request through them from another handler.
		if functionName == "" || strings.ContainsAny(functionName, " \t\n(){}") {
			continue
		}
		functionBody := sliceBalancedBlockAfter(t, pageSource[declarationIndex:], declarationToken+functionName)
		if strings.Contains(functionBody, "setPointerCapture(") {
			capturingNames = append(capturingNames, functionName)
		}
	}
	return capturingNames
}

// TestTimelinePointerCaptureWaitsForThePanEngage is the check REQ-337 exists to add, and
// it pins the regression REQ-336 fixed: pointer capture taken on pointerdown.
//
// Why this is not covered by the probe above. Pointer capture retargets the pointer
// events AND the click the engine synthesizes to the capturing element, so a capture
// taken on every press leaves the delegated handler in board-controls.js with no
// [data-detail-kind] ancestor to find and the detail drawer never opens for any click in
// the chart. TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving passed through
// that entire regression, because it dispatches synthetic PointerEvents whose pointerId
// the engine does not know: setPointerCapture throws on them, capture is never
// established, and the failing path is invisible to the lane. This lane runs under
// --dump-dom with no protocol channel, so it cannot dispatch trusted input at all; the
// REQ therefore allows asserting the structural property instead, and this does.
//
// Three assertions, because "capture is not taken on pointerdown" is trivially satisfied
// by taking it nowhere — which is the OTHER bug, the one REQ-333 fixed, where a drag
// released outside the chart never told the host it had ended:
//
//	(a) the pointerdown handler requests no capture itself,
//	(b) it calls nothing that requests capture, resolved from the page rather than from a
//	    list, so a fresh wrapper called from pointerdown fails here too, and
//	(c) the pointermove handler — the engage path — DOES reach a request.
//
// Known residual: this reads text, so a request routed through a variable, a method
// lookup, or an eval would pass. The mutation this REQ pins (reintroducing capture on
// pointerdown, by either spelling) does not.
func TestTimelinePointerCaptureWaitsForThePanEngage(t *testing.T) {
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read generated index.html: %v", readError)
	}
	indexHTML := string(indexBytes)

	capturingNames := pointerCapturingFunctionNames(t, indexHTML)
	// Vacuity guard: with no capturing function found, (b) below asserts nothing and (c)
	// would be measuring the wrong thing. A page that requests capture nowhere is a
	// failure, not a pass.
	if len(capturingNames) == 0 {
		t.Fatal("no named function in the generated page requests pointer capture, so this check " +
			"cannot tell capture-at-the-engage from capture nowhere at all")
	}

	pointerDownBody := sliceBalancedBlockAfter(t, indexHTML,
		"addTimelineListener(scrollHost, \"pointerdown\"")
	pointerMoveBody := sliceBalancedBlockAfter(t, indexHTML,
		"addTimelineListener(scrollHost, \"pointermove\"")

	// (a) Nothing in the press handler requests capture directly.
	if strings.Contains(pointerDownBody, "setPointerCapture(") {
		t.Error("the Timeline pointerdown handler requests pointer capture; capture retargets the " +
			"synthesized click to the capturing element, so every mouse click in the chart loses " +
			"its [data-detail-kind] target and the detail drawer stops opening")
	}
	// (b) And it reaches no function that requests capture on its behalf.
	for _, capturingName := range capturingNames {
		if strings.Contains(pointerDownBody, capturingName+"(") {
			t.Errorf("the Timeline pointerdown handler calls %s, which requests pointer capture; "+
				"capture must wait for the pan to engage or every click in the chart is retargeted "+
				"to the capturing element", capturingName)
		}
	}
	// (c) The engage path still requests it, so (a) and (b) cannot be satisfied by
	// removing capture altogether and reopening REQ-333's never-released drag.
	engageReachesCapture := strings.Contains(pointerMoveBody, "setPointerCapture(")
	for _, capturingName := range capturingNames {
		if strings.Contains(pointerMoveBody, capturingName+"(") {
			engageReachesCapture = true
		}
	}
	if !engageReachesCapture {
		t.Errorf("the Timeline pointermove handler reaches no pointer-capture request (capturing "+
			"functions in the page: %v); a drag released outside the chart then has no guaranteed "+
			"path back to the host that armed it", capturingNames)
	}
}
