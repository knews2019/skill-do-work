package main

import (
	"encoding/json"
	"math"
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
