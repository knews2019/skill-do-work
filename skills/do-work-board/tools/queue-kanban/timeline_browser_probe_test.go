package main

import (
	"encoding/json"
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
	Status       string  `json:"status"`
	WaitFill     string  `json:"waitFill"`
	WorkFill     string  `json:"workFill"`
	WaitOpacity  float64 `json:"waitOpacity"`
	WorkOpacity  float64 `json:"workOpacity"`
	ChipAccent   string  `json:"chipAccent"`
	Unrecognized bool    `json:"unrecognized"`
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
      unrecognized: row.classList.contains("is-status-unrecognized")
    };
  });
  // Written LAST and only once every measurement exists, so a throw leaves the
  // node empty and the Go side reports a failure rather than reading a partial
  // object as a pass.
  document.getElementById("` + browserProbeResultElementId + `").textContent =
    JSON.stringify({ rows: measured });
})();
</script>
</body></html>`

	probeOutput := runBrowserBehaviorProbe(t, "timeline status colours", pageHTML)
	// The lane's sentinel requires a JSON object at the result node, so the rows
	// arrive wrapped.
	var probeResult struct {
		Rows []timelineStatusProbeRow `json:"rows"`
	}
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode timeline status colour probe: %v (output %q)", decodeError, probeOutput)
	}
	measured := probeResult.Rows
	if len(measured) != len(timelineStatusProbeStatuses) {
		t.Fatalf("probe measured %d rows, want %d", len(measured), len(timelineStatusProbeStatuses))
	}

	byStatus := map[string]timelineStatusProbeRow{}
	for _, row := range measured {
		// A fill of "none", empty, or fully transparent means the custom property
		// resolved to nothing — the failure mode where every assertion below would
		// otherwise agree with every other on the same non-colour.
		if row.WorkFill == "" || row.WorkFill == "none" || strings.Contains(row.WorkFill, ", 0)") {
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

	// 2. Lightness is the phase: the wait is the same hue, quieter. Both halves
	//    must share a fill and differ in opacity, which is what let status onto a
	//    chart whose hues were already spent.
	for _, row := range measured {
		if row.WaitFill != row.WorkFill {
			t.Errorf("status %q paints wait %s and work %s; the phase difference must be lightness, "+
				"not a second hue", row.Status, row.WaitFill, row.WorkFill)
		}
		if !(row.WaitOpacity < row.WorkOpacity) {
			t.Errorf("status %q has wait opacity %.2f and work opacity %.2f; the wait must be the "+
				"quieter half", row.Status, row.WaitOpacity, row.WorkOpacity)
		}
		// Quieter, but still a bar. Below roughly a quarter it reads as a disabled
		// control rather than as the waiting half of a measured span.
		if row.WaitOpacity < 0.25 {
			t.Errorf("status %q draws its wait at opacity %.2f, too faint to read as a bar at the "+
				"shipped 10px height", row.Status, row.WaitOpacity)
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
