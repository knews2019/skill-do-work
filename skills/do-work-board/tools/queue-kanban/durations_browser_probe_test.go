package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

// sliceGeneratedStyleBlock lifts the page's inlined stylesheet for focused
// browser probes that need the shipped board styles without the full page.
func sliceGeneratedStyleBlock(t *testing.T, indexHtml string) string {
	t.Helper()
	styleStart := strings.Index(indexHtml, "<style>")
	styleEnd := strings.Index(indexHtml, "</style>")
	if styleStart < 0 || styleEnd < styleStart {
		t.Fatal("generated page carries no inlined <style> block")
	}
	return indexHtml[styleStart+len("<style>") : styleEnd]
}

type durationsDensePanelProbeResult struct {
	LocationHref          string  `json:"locationHref"`
	RenderedSampleCount   int     `json:"renderedSampleCount"`
	PayloadSampleCount    int     `json:"payloadSampleCount"`
	ActiveDayCount        int     `json:"activeDayCount"`
	DaySlotWidth          float64 `json:"daySlotWidth"`
	BusyDaySpread         float64 `json:"busyDaySpread"`
	Deterministic         bool    `json:"deterministic"`
	EveryMarkInOwnDay     bool    `json:"everyMarkInOwnDay"`
	HoveredRequestId      string  `json:"hoveredRequestId"`
	HoverReadout          string  `json:"hoverReadout"`
	RibbonFiniteBounded   bool    `json:"ribbonFiniteBounded"`
	MedianFiniteBounded   bool    `json:"medianFiniteBounded"`
	RibbonOpacity         float64 `json:"ribbonOpacity"`
	MarkOpacity           float64 `json:"markOpacity"`
	BodyBackground        string  `json:"bodyBackground"`
	LongestSpanCount      int     `json:"longestSpanCount"`
	ExpectedSpanCount     int     `json:"expectedSpanCount"`
	LongestSpanOrder      bool    `json:"longestSpanOrder"`
	EveryListField        bool    `json:"everyListField"`
	CountSentence         string  `json:"countSentence"`
	ListOutsideSVG        bool    `json:"listOutsideSvg"`
	SharedWrapper         bool    `json:"sharedWrapper"`
	SVGRequestLabels      int     `json:"svgRequestLabels"`
	SVGLeaderLines        int     `json:"svgLeaderLines"`
	SVGMoreSentences      int     `json:"svgMoreSentences"`
	OverflowHoverId       string  `json:"overflowHoverId"`
	OverflowHoverDuration string  `json:"overflowHoverDuration"`
	OverflowHoverReadout  string  `json:"overflowHoverReadout"`
	ViewportWidth         float64 `json:"viewportWidth"`
	DocumentWidth         float64 `json:"documentWidth"`
	WrapperRight          float64 `json:"wrapperRight"`
	AsideRight            float64 `json:"asideRight"`
	WrapperColumns        string  `json:"wrapperColumns"`
	AsideOverflowY        string  `json:"asideOverflowY"`
}

type durationsHeadlineBrowserProbeResult struct {
	LocationHref          string     `json:"locationHref"`
	ViewportWidth         float64    `json:"viewportWidth"`
	WindowStats           [][]string `json:"windowStats"`
	StatItemCount         int        `json:"statItemCount"`
	DefinitionTermCount   int        `json:"definitionTermCount"`
	DefinitionValueCount  int        `json:"definitionValueCount"`
	DefinitionListTag     string     `json:"definitionListTag"`
	DefinitionListTabStop bool       `json:"definitionListTabStop"`
	NativeWindowButton    bool       `json:"nativeWindowButton"`
	WindowButtonFocused   bool       `json:"windowButtonFocused"`
	StatRowCount          int        `json:"statRowCount"`
	StatsClearChart       bool       `json:"statsClearChart"`
	StatTilesOverlap      bool       `json:"statTilesOverlap"`
	RollingFinite         bool       `json:"rollingFinite"`
	RollingMarkerCount    int        `json:"rollingMarkerCount"`
	RollingLineContrast   float64    `json:"rollingLineContrast"`
	RollingMarkContrast   float64    `json:"rollingMarkContrast"`
	BodyBackground        string     `json:"bodyBackground"`
	CountTicksSeparate    bool       `json:"countTicksSeparate"`
	CountTickTexts        []string   `json:"countTickTexts"`
	ConsoleErrors         []string   `json:"consoleErrors"`
}

func durationsMarkActivationFixtureTickets() []*RequestTicket {
	fixtureDay := time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
	// Completion order and REQ-id order deliberately disagree. REQ-349's
	// within-day jitter ranks by id, while a pre-jitter projection ranks by the
	// completion instant. A trusted click at REQ-500's rendered centre therefore
	// distinguishes the two targeting models instead of merely showing that some
	// mark opens some drawer.
	busyDaySamples := []struct {
		requestId      string
		completionHour int
	}{
		{requestId: "REQ-500", completionHour: 8},
		{requestId: "REQ-100", completionHour: 9},
		{requestId: "REQ-400", completionHour: 10},
		{requestId: "REQ-200", completionHour: 11},
		{requestId: "REQ-300", completionHour: 12},
	}
	tickets := make([]*RequestTicket, 0, len(busyDaySamples)+2)
	for _, sample := range busyDaySamples {
		completedAt := fixtureDay.Add(time.Duration(sample.completionHour) * time.Hour)
		tickets = append(tickets, durationTicket(
			sample.requestId, "B",
			completedAt.Add(-10*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	for anchorIndex, dayOffset := range []int{-2, 2} {
		completedAt := fixtureDay.AddDate(0, 0, dayOffset).Add(10 * time.Hour)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-60%d", anchorIndex), "A",
			completedAt.Add(-20*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	return tickets
}

func generateDurationsMarkActivationSite(t *testing.T) string {
	t.Helper()
	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC),
		ProjectName: "REQ-354 Durations mark activation probe",
		AllRequests: durationsMarkActivationFixtureTickets(),
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate Durations mark activation fixture: %v", generateError)
	}
	return siteDirectory
}

// Every projected mark is one semantic control, but the chart costs only one Tab
// press: Left/Right move the sole roving stop and Enter/Space activate the focused
// REQ. The walk is exhaustive so a renderer that makes only the first few samples
// keyboard-reachable cannot satisfy the contract by accident.
func TestBrowserBehaviorDurationsMarksAreOneRovingTabStop(t *testing.T) {
	siteDirectory := generateDurationsMarkActivationSite(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read Durations activation fixture: %v", readError)
	}
	probeScript := `
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
function durationMarks() {
  return Array.prototype.slice.call(document.querySelectorAll("#durations-chart circle.durations-mark"));
}
function pressDurationKey(keyName) {
  document.activeElement.dispatchEvent(new KeyboardEvent("keydown", {
    key: keyName, bubbles: true, cancelable: true, composed: true
  }));
}
function durationMarkState() {
  var marks = durationMarks();
  var active = document.activeElement;
  return {
    ids: marks.map(function (mark) { return mark.getAttribute("data-detail-id") || ""; }),
    roles: marks.map(function (mark) { return mark.getAttribute("role") || ""; }),
    labels: marks.map(function (mark) { return mark.getAttribute("aria-label") || ""; }),
    tabbableCount: marks.filter(function (mark) { return mark.getAttribute("tabindex") === "0"; }).length,
    skippedCount: marks.filter(function (mark) { return mark.getAttribute("tabindex") === "-1"; }).length,
    focusedId: active && active.getAttribute ? (active.getAttribute("data-detail-id") || "") : ""
  };
}
window.addEventListener("load", function () {
  setTimeout(function () {
    document.querySelector('[data-view-target="durations"]').click();
    document.querySelector('[data-durations-window="all"]').click();
    var marks = durationMarks();
    var initial = durationMarkState();
    var forwardIds = [];
    var forwardTabCounts = [];
    marks[0].focus();
    forwardIds.push(durationMarkState().focusedId);
    forwardTabCounts.push(durationMarkState().tabbableCount);
    for (var markIndex = 1; markIndex < marks.length; markIndex += 1) {
      pressDurationKey("ArrowRight");
      forwardIds.push(durationMarkState().focusedId);
      forwardTabCounts.push(durationMarkState().tabbableCount);
    }
    var lastId = durationMarkState().focusedId;
    pressDurationKey("Enter");
    var enterDrawerId = document.getElementById("detail-id").textContent.trim();
    var enterOpened = document.getElementById("detail-drawer").hidden === false;
    document.getElementById("detail-close").click();
    pressDurationKey(" ");
    var spaceDrawerId = document.getElementById("detail-id").textContent.trim();
    var spaceOpened = document.getElementById("detail-drawer").hidden === false;
    document.getElementById("detail-close").click();
    var reverseIds = [durationMarkState().focusedId];
    for (var reverseIndex = 1; reverseIndex < marks.length; reverseIndex += 1) {
      pressDurationKey("ArrowLeft");
      reverseIds.push(durationMarkState().focusedId);
    }
    document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
      href: location.href,
      initial: initial,
      forwardIds: forwardIds,
      forwardTabCounts: forwardTabCounts,
      reverseIds: reverseIds,
      lastId: lastId,
      enterOpened: enterOpened,
      enterDrawerId: enterDrawerId,
      spaceOpened: spaceOpened,
      spaceDrawerId: spaceDrawerId,
      final: durationMarkState(),
      hoverSurfaceTabIndex: document.querySelector(".durations-hover-surface").getAttribute("tabindex")
    });
  }, 100);
});
</script>
</body>`
	probePage := strings.Replace(string(indexBytes), "</body>", probeScript, 1)
	if probePage == string(indexBytes) {
		t.Fatal("generated Durations activation page has no </body> injection point")
	}
	resultJSON := runBrowserBehaviorProbeInDirectory(
		t, "Durations roving mark controls", siteDirectory, probePage,
		"--window-size=1280,1000", "--virtual-time-budget=30000",
	)
	type markState struct {
		IDs           []string `json:"ids"`
		Roles         []string `json:"roles"`
		Labels        []string `json:"labels"`
		TabbableCount int      `json:"tabbableCount"`
		SkippedCount  int      `json:"skippedCount"`
		FocusedID     string   `json:"focusedId"`
	}
	var result struct {
		Href                 string    `json:"href"`
		Initial              markState `json:"initial"`
		ForwardIDs           []string  `json:"forwardIds"`
		ForwardTabCounts     []int     `json:"forwardTabCounts"`
		ReverseIDs           []string  `json:"reverseIds"`
		LastID               string    `json:"lastId"`
		EnterOpened          bool      `json:"enterOpened"`
		EnterDrawerID        string    `json:"enterDrawerId"`
		SpaceOpened          bool      `json:"spaceOpened"`
		SpaceDrawerID        string    `json:"spaceDrawerId"`
		Final                markState `json:"final"`
		HoverSurfaceTabIndex *string   `json:"hoverSurfaceTabIndex"`
	}
	if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
		t.Fatalf("decode Durations roving probe: %v\n%s", decodeError, resultJSON)
	}
	if !strings.HasSuffix(result.Href, "/"+browserProbePageFileName) {
		t.Fatalf("Durations roving probe measured %q, not its probe page", result.Href)
	}
	markCount := len(result.Initial.IDs)
	if markCount < 7 {
		t.Fatalf("Durations roving fixture rendered %d marks, want at least seven", markCount)
	}
	if result.Initial.TabbableCount != 1 || result.Initial.SkippedCount != markCount-1 || result.HoverSurfaceTabIndex != nil {
		t.Errorf("initial Tab stops: marks 0/-1=%d/%d of %d, hover tabindex=%v; want 1/%d and no overlay stop",
			result.Initial.TabbableCount, result.Initial.SkippedCount, markCount,
			result.HoverSurfaceTabIndex, markCount-1)
	}
	for markIndex := range result.Initial.IDs {
		if result.Initial.IDs[markIndex] == "" || result.Initial.Roles[markIndex] != "button" ||
			!strings.Contains(result.Initial.Labels[markIndex], result.Initial.IDs[markIndex]) {
			t.Errorf("mark %d semantics id/role/label = %q/%q/%q", markIndex,
				result.Initial.IDs[markIndex], result.Initial.Roles[markIndex], result.Initial.Labels[markIndex])
		}
	}
	if !reflect.DeepEqual(result.ForwardIDs, result.Initial.IDs) {
		t.Errorf("ArrowRight walk = %q, want every projected mark in DOM order %q", result.ForwardIDs, result.Initial.IDs)
	}
	wantReverseIDs := append([]string(nil), result.Initial.IDs...)
	for leftIndex, rightIndex := 0, len(wantReverseIDs)-1; leftIndex < rightIndex; leftIndex, rightIndex = leftIndex+1, rightIndex-1 {
		wantReverseIDs[leftIndex], wantReverseIDs[rightIndex] = wantReverseIDs[rightIndex], wantReverseIDs[leftIndex]
	}
	if !reflect.DeepEqual(result.ReverseIDs, wantReverseIDs) {
		t.Errorf("ArrowLeft walk = %q, want reverse projected order %q", result.ReverseIDs, wantReverseIDs)
	}
	for stepIndex, tabStopCount := range result.ForwardTabCounts {
		if tabStopCount != 1 {
			t.Errorf("ArrowRight step %d left %d tabbable marks, want one", stepIndex, tabStopCount)
		}
	}
	if !result.EnterOpened || result.EnterDrawerID != result.LastID ||
		!result.SpaceOpened || result.SpaceDrawerID != result.LastID {
		t.Errorf("activation at %q: Enter open/id=%v/%q, Space open/id=%v/%q",
			result.LastID, result.EnterOpened, result.EnterDrawerID, result.SpaceOpened, result.SpaceDrawerID)
	}
	if result.Final.FocusedID != result.Initial.IDs[0] || result.Final.TabbableCount != 1 {
		t.Errorf("reverse walk ended focused/tabbable=%q/%d, want first mark %q and one stop",
			result.Final.FocusedID, result.Final.TabbableCount, result.Initial.IDs[0])
	}
}

// The pointer trial uses Chromium's real input path. It aims through the transparent
// overlay at a busy-day circle whose REQ-349 jitter position names a different REQ
// than the old completion-time position would. That mutation guard keeps "the drawer
// opened" from passing when click targeting quietly drifts away from rendered marks.
func TestBrowserBehaviorDurationsTrustedClickOpensJitteredMark(t *testing.T) {
	siteDirectory := generateDurationsMarkActivationSite(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read Durations trusted-click fixture: %v", readError)
	}
	session := startTrustedInputBrowserSession(
		t, "Durations trusted jittered mark click", siteDirectory, string(indexBytes),
		"--window-size=1280,1000",
	)
	session.evaluateInPage(t, `(function () {
  document.querySelector('[data-view-target="durations"]').click();
  document.querySelector('[data-durations-window="all"]').click();
  window.__durationTrustedClickTarget = "";
  document.addEventListener("click", function (event) {
    window.__durationTrustedClickTarget = event.target.getAttribute("class") || event.target.tagName;
  }, true);
  return true;
})()`)
	session.waitForPageCondition(t, "the all-history Durations marks to render",
		`document.querySelectorAll("#durations-chart circle.durations-mark").length >= 7`)

	type clickAim struct {
		RequestID    string  `json:"requestId"`
		RawNearestID string  `json:"rawNearestId"`
		JitteredX    float64 `json:"jitteredX"`
		RawX         float64 `json:"rawX"`
		ViewportX    float64 `json:"viewportX"`
		ViewportY    float64 `json:"viewportY"`
	}
	var aim clickAim
	session.decodeResult(t, "Durations trusted-click aim", session.evaluateInPage(t, `(function () {
  var samples = window.queueKanbanBoardData.durations.samples;
  var targetIndex = samples.findIndex(function (sample) { return sample.id === "REQ-500"; });
  var marks = Array.from(document.querySelectorAll("#durations-chart circle.durations-mark"));
  var target = marks[targetIndex];
  target.scrollIntoView({block: "center", inline: "center"});
  var firstSampleMs = Date.parse(samples[0].completionTime);
  var lastSampleMs = Date.parse(samples[samples.length - 1].completionTime);
  var dayMs = 86400000;
  var timeStart = Math.floor(firstSampleMs / dayMs) * dayMs;
  var timeEnd = Math.floor(lastSampleMs / dayMs) * dayMs + dayMs;
  var hoverSurface = document.querySelector("#durations-chart .durations-hover-surface");
  var plotLeft = Number(hoverSurface.getAttribute("x"));
  var plotWidth = Number(hoverSurface.getAttribute("width"));
  var rawX = function (sample) {
    return plotLeft + ((Date.parse(sample.completionTime) - timeStart) / (timeEnd - timeStart)) * plotWidth;
  };
  var targetX = Number(target.getAttribute("cx"));
  var targetY = Number(target.getAttribute("cy"));
  var rawNearest = samples.reduce(function (nearest, sample, sampleIndex) {
    var markY = Number(marks[sampleIndex].getAttribute("cy"));
    var distance = Math.abs(rawX(sample) - targetX) + Math.abs(markY - targetY) * 0.35;
    return !nearest || distance < nearest.distance ? {id: sample.id, distance: distance} : nearest;
  }, null);
  var rect = target.getBoundingClientRect();
  return {
    requestId: samples[targetIndex].id,
    rawNearestId: rawNearest.id,
    jitteredX: targetX,
    rawX: rawX(samples[targetIndex]),
    viewportX: rect.left + rect.width / 2,
    viewportY: rect.top + rect.height / 2
  };
})()`), &aim)
	if aim.RequestID == "" || aim.RawNearestID == "" || aim.RawNearestID == aim.RequestID ||
		math.Abs(aim.JitteredX-aim.RawX) < 1 {
		t.Fatalf("trusted-click fixture is not mutation-sensitive: target=%q raw-nearest=%q jitter/raw x=%.2f/%.2f",
			aim.RequestID, aim.RawNearestID, aim.JitteredX, aim.RawX)
	}
	session.dispatchTrustedMouseEvent(t, "mouseMoved", aim.ViewportX, aim.ViewportY, "none", 0)
	hoverReadoutArrived := session.pageConditionHoldsWithin(t, "the jittered Durations hover readout",
		`document.getElementById("durations-readout").textContent.indexOf("`+aim.RequestID+` ·") === 0`,
		browserProbeGestureSettleDeadline)
	var hoverReadout string
	session.decodeResult(t, "Durations trusted-hover outcome", session.evaluateInPage(t,
		`document.getElementById("durations-readout").textContent`), &hoverReadout)
	session.pressTrustedMouseAt(t, aim.ViewportX, aim.ViewportY)
	session.releaseTrustedMouseAt(t, aim.ViewportX, aim.ViewportY)
	drawerOpened := session.pageConditionHoldsWithin(t, "the clicked Durations mark drawer to open",
		`document.getElementById("detail-drawer").hidden === false`, browserProbeGestureSettleDeadline)
	var outcome struct {
		ClickTarget string `json:"clickTarget"`
		DrawerID    string `json:"drawerId"`
	}
	session.decodeResult(t, "Durations trusted-click outcome", session.evaluateInPage(t, `({
  clickTarget: window.__durationTrustedClickTarget,
  drawerId: document.getElementById("detail-id").textContent.trim()
})`), &outcome)
	if !hoverReadoutArrived || !strings.HasPrefix(hoverReadout, aim.RequestID+" ·") {
		t.Errorf("trusted hover at jittered %s read %q; raw targeting would choose %s",
			aim.RequestID, hoverReadout, aim.RawNearestID)
	}
	if !drawerOpened {
		t.Errorf("trusted click at jittered %s did not open the detail drawer (target %q, shown %q)",
			aim.RequestID, outcome.ClickTarget, outcome.DrawerID)
	}
	if outcome.ClickTarget != "durations-hover-surface" {
		t.Errorf("trusted click landed on %q, want the overlay that owns nearest-mark targeting", outcome.ClickTarget)
	}
	if outcome.DrawerID != aim.RequestID {
		t.Errorf("trusted click at jittered %s opened %q; raw targeting would choose %s",
			aim.RequestID, outcome.DrawerID, aim.RawNearestID)
	}
}

func durationsHeadlineBrowserFixtureTickets() []*RequestTicket {
	eligibleDays := []struct {
		completed time.Time
		minutes   time.Duration
	}{
		{completed: time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC), minutes: 10 * time.Minute},
		{completed: time.Date(2026, 7, 2, 10, 0, 0, 0, time.UTC), minutes: 70 * time.Minute},
		{completed: time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC), minutes: 20 * time.Minute},
		{completed: time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC), minutes: 60 * time.Minute},
		{completed: time.Date(2026, 7, 11, 10, 0, 0, 0, time.UTC), minutes: 30 * time.Minute},
		{completed: time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC), minutes: 50 * time.Minute},
		{completed: time.Date(2026, 8, 10, 10, 0, 0, 0, time.UTC), minutes: 40 * time.Minute},
		{completed: time.Date(2026, 9, 20, 10, 0, 0, 0, time.UTC), minutes: 80 * time.Minute},
	}
	tickets := make([]*RequestTicket, 0, len(eligibleDays)+9)
	for dayIndex, eligibleDay := range eligibleDays {
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 880+dayIndex), "B",
			eligibleDay.completed.Add(-eligibleDay.minutes).Format(time.RFC3339),
			eligibleDay.completed.Format(time.RFC3339),
		))
	}
	for pausedIndex := 0; pausedIndex < 5; pausedIndex++ {
		completed := time.Date(2026, 7, 7, 10+pausedIndex, 0, 0, 0, time.UTC)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 900+pausedIndex), "C",
			completed.Add(-8*time.Hour).Format(time.RFC3339), completed.Format(time.RFC3339),
		))
	}
	// Four all-history-only paused spans make every headline tile visibly change
	// between 90 days and all history, including p90 and the rounded cadence.
	for oldIndex, oldMinutes := range []time.Duration{10 * time.Hour, 11 * time.Hour, 12 * time.Hour, 15 * time.Hour} {
		completed := time.Date(2026, time.Month(4+oldIndex/2), 1+oldIndex%2, 10, 0, 0, 0, time.UTC)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 920+oldIndex), "A",
			completed.Add(-oldMinutes).Format(time.RFC3339), completed.Format(time.RFC3339),
		))
	}
	return tickets
}

// The complete generated board proves the headline and rolling surfaces in the
// browser. Every measurement returns location.href in the same result, measures
// the rolling ink against the transparent SVG's real body background, and checks
// the responsive stat grid at the three maintained viewport widths.
func TestBrowserBehaviorDurationsHeadlineAndRollingSeries(t *testing.T) {
	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC),
		ProjectName: "REQ-352 Durations headline probe",
		AllRequests: durationsHeadlineBrowserFixtureTickets(),
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate headline fixture board: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read headline fixture index: %v", readError)
	}

	probePage := strings.Replace(string(indexBytes), "<head>", `<head><script>
window.__durationsProbeErrors = [];
window.addEventListener("error", function (event) { window.__durationsProbeErrors.push(String(event.message)); });
(function () {
  var originalConsoleError = console.error;
  console.error = function () {
    window.__durationsProbeErrors.push(Array.prototype.join.call(arguments, " "));
    originalConsoleError.apply(console, arguments);
  };
})();
</script>`, 1)
	probeScript := `
  (function () {
  viewState.view = "durations";
  applyView();

  function statValues() {
    return ["median", "p90", "active-days", "reqs-per-day"].map(function (statName) {
      return document.getElementById("durations-stat-" + statName).textContent;
    });
  }
  var windowStats = ["30", "90", "all"].map(function (windowName) {
    applyDurationsWindowSelection(windowName);
    return statValues();
  });

  var definitionList = document.getElementById("durations-stats");
  var statTiles = Array.from(definitionList.children);
  var statTileRects = statTiles.map(function (tile) { return tile.getBoundingClientRect(); });
  var statRowCount = new Set(statTileRects.map(function (rect) { return rect.top.toFixed(2); })).size;
  var statTilesOverlap = statTileRects.some(function (first, firstIndex) {
    return statTileRects.some(function (second, secondIndex) {
      return secondIndex > firstIndex && first.left < second.right && second.left < first.right &&
        first.top < second.bottom && second.top < first.bottom;
    });
  });
  var chartRect = document.getElementById("durations-chart").getBoundingClientRect();
  var statsRect = definitionList.getBoundingClientRect();

  var rollingPath = document.querySelector("#durations-chart .durations-rolling-line");
  var rollingMarkers = Array.from(document.querySelectorAll("#durations-chart .durations-rolling-marker"));
  var rollingBox = rollingPath ? rollingPath.getBBox() : { x: NaN, y: NaN, width: NaN, height: NaN };
  var rollingFinite = !!rollingPath && !!rollingPath.getAttribute("d") &&
    !/NaN|Infinity/.test(rollingPath.getAttribute("d")) &&
    [rollingBox.x, rollingBox.y, rollingBox.width, rollingBox.height, rollingPath.getTotalLength()].every(Number.isFinite) &&
    rollingPath.getTotalLength() > 0;

  function rgbChannels(colour) {
    return (colour.match(/[\d.]+/g) || []).slice(0, 3).map(Number);
  }
  function relativeLuminance(colour) {
    return rgbChannels(colour).map(function (channel) {
      var normalized = channel / 255;
      return normalized <= 0.03928 ? normalized / 12.92 : Math.pow((normalized + 0.055) / 1.055, 2.4);
    }).reduce(function (sum, channel, index) {
      return sum + channel * [0.2126, 0.7152, 0.0722][index];
    }, 0);
  }
  function contrastRatio(first, second) {
    var firstLuminance = relativeLuminance(first);
    var secondLuminance = relativeLuminance(second);
    return (Math.max(firstLuminance, secondLuminance) + 0.05) /
      (Math.min(firstLuminance, secondLuminance) + 0.05);
  }
  var bodyBackground = getComputedStyle(document.body).backgroundColor;
  var rollingStyle = rollingPath ? getComputedStyle(rollingPath) : { stroke: "rgb(0, 0, 0)" };
  var markerStyle = rollingMarkers.length ? getComputedStyle(rollingMarkers[0]) : { fill: "rgb(0, 0, 0)" };

  var countTicks = Array.from(document.querySelectorAll('[data-durations-count-tick="true"]'));
  var countTickRects = countTicks.map(function (tick) { return tick.getBoundingClientRect(); });
  var countTicksSeparate = countTickRects.every(function (first, firstIndex) {
    return countTickRects.every(function (second, secondIndex) {
      return secondIndex <= firstIndex || first.bottom <= second.top || second.bottom <= first.top;
    });
  });

  var ninetyDayButton = document.querySelector('[data-durations-window="90"]');
  ninetyDayButton.focus();
  document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
    locationHref: location.href,
    viewportWidth: innerWidth,
    windowStats: windowStats,
    statItemCount: statTiles.length,
    definitionTermCount: definitionList.querySelectorAll("dt").length,
    definitionValueCount: definitionList.querySelectorAll("dd").length,
    definitionListTag: definitionList.tagName,
    definitionListTabStop: definitionList.tabIndex >= 0 || !!definitionList.querySelector("[tabindex]"),
    nativeWindowButton: ninetyDayButton.tagName === "BUTTON",
    windowButtonFocused: document.activeElement === ninetyDayButton,
    statRowCount: statRowCount,
    statsClearChart: statsRect.bottom <= chartRect.top,
    statTilesOverlap: statTilesOverlap,
    rollingFinite: rollingFinite,
    rollingMarkerCount: rollingMarkers.length,
    rollingLineContrast: contrastRatio(rollingStyle.stroke, bodyBackground),
    rollingMarkContrast: contrastRatio(markerStyle.fill, bodyBackground),
    bodyBackground: bodyBackground,
    countTicksSeparate: countTicksSeparate,
    countTickTexts: countTicks.map(function (tick) { return tick.textContent; }),
    consoleErrors: window.__durationsProbeErrors
  });
  })();
`
	clientClose := strings.LastIndex(probePage, "})();")
	if clientClose < 0 {
		t.Fatal("generated page has no client IIFE close for headline probe")
	}
	clientScriptStart := strings.LastIndex(probePage[:clientClose], "<script>")
	if clientScriptStart < 0 {
		t.Fatal("generated page has no inline client script for headline probe")
	}
	resultNode := `<pre id="` + browserProbeResultElementId + `" hidden></pre>`
	probePage = probePage[:clientScriptStart] + resultNode + probePage[clientScriptStart:]
	clientClose += len(resultNode)
	probePage = probePage[:clientClose] + probeScript + probePage[clientClose:]

	probeCases := []struct {
		name           string
		width          int
		colourFlag     string
		wantStatRowsAt int
	}{
		{name: "320-light", width: 320, colourFlag: "--force-light-mode", wantStatRowsAt: 2},
		{name: "768-light", width: 768, colourFlag: "--force-light-mode", wantStatRowsAt: 2},
		{name: "1280-light", width: 1280, colourFlag: "--force-light-mode", wantStatRowsAt: 1},
		{name: "1280-dark", width: 1280, colourFlag: "--blink-settings=preferredColorScheme=2", wantStatRowsAt: 1},
	}
	for _, probeCase := range probeCases {
		probeCase := probeCase
		t.Run(probeCase.name, func(t *testing.T) {
			resultJSON := runBrowserBehaviorProbeInDirectory(
				t, "REQ-352 headline "+probeCase.name, siteDirectory, probePage,
				"--headless=new", fmt.Sprintf("--window-size=%d,1100", probeCase.width), probeCase.colourFlag,
			)
			var result durationsHeadlineBrowserProbeResult
			if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
				t.Fatalf("decode headline browser result: %v\n%s", decodeError, resultJSON)
			}
			if result.LocationHref == "" || !strings.Contains(result.LocationHref, browserProbePageFileName) {
				t.Errorf("probe measured unnamed page %q", result.LocationHref)
			}
			if result.StatItemCount != 4 || result.DefinitionTermCount != 4 || result.DefinitionValueCount != 4 || result.DefinitionListTag != "DL" {
				t.Errorf("semantic stats = tag %q, items/terms/values %d/%d/%d; want DL 4/4/4",
					result.DefinitionListTag, result.StatItemCount, result.DefinitionTermCount, result.DefinitionValueCount)
			}
			if result.DefinitionListTabStop || !result.NativeWindowButton || !result.WindowButtonFocused {
				t.Errorf("keyboard semantics = dl tab stop %v, native button %v, focus %v",
					result.DefinitionListTabStop, result.NativeWindowButton, result.WindowButtonFocused)
			}
			if result.StatTilesOverlap || !result.StatsClearChart {
				t.Errorf("stat layout overlaps: tiles=%v chart=%v", result.StatTilesOverlap, !result.StatsClearChart)
			}
			if probeCase.wantStatRowsAt == 1 && result.StatRowCount != 1 {
				t.Errorf("%dpx viewport made %d stat rows, want one", probeCase.width, result.StatRowCount)
			}
			if probeCase.wantStatRowsAt > 1 && result.StatRowCount < probeCase.wantStatRowsAt {
				t.Errorf("%dpx viewport made %d stat rows, want at least %d", probeCase.width, result.StatRowCount, probeCase.wantStatRowsAt)
			}
			if len(result.WindowStats) != 3 {
				t.Fatalf("window update captured %d states, want 30/90/all", len(result.WindowStats))
			}
			for statIndex := 0; statIndex < 4; statIndex++ {
				distinctValues := map[string]bool{}
				for _, windowStats := range result.WindowStats {
					if len(windowStats) != 4 {
						t.Fatalf("window stats = %#v, want four values", result.WindowStats)
					}
					distinctValues[windowStats[statIndex]] = true
				}
				if len(distinctValues) != 3 {
					t.Errorf("headline stat %d did not change across 30/90/all: %#v", statIndex, result.WindowStats)
				}
			}
			if !result.RollingFinite || result.RollingMarkerCount != 2 {
				t.Errorf("rolling geometry finite=%v markers=%d, want a finite line and two points",
					result.RollingFinite, result.RollingMarkerCount)
			}
			if result.RollingLineContrast < 3 || result.RollingMarkContrast < 3 || result.BodyBackground == "" {
				t.Errorf("rolling contrast line/marker %.2f/%.2f against body %q, want both >= 3:1",
					result.RollingLineContrast, result.RollingMarkContrast, result.BodyBackground)
			}
			if !result.CountTicksSeparate || !reflect.DeepEqual(result.CountTickTexts, []string{"5", "2.5", "0"}) {
				t.Errorf("Panel C ticks separate=%v texts=%q, want separated 5/2.5/0", result.CountTicksSeparate, result.CountTickTexts)
			}
			if len(result.ConsoleErrors) != 0 {
				t.Errorf("browser console errors: %q", result.ConsoleErrors)
			}
			t.Logf("%s %.0fpx: stats rows %d, rolling contrast %.2f/%.2f against %s",
				result.LocationHref, result.ViewportWidth, result.StatRowCount,
				result.RollingLineContrast, result.RollingMarkContrast, result.BodyBackground)
		})
	}
}

// REQ-349's target board is materially denser than this repository. This probe
// renders 705 real samples on 47 active days spread across 139 UTC days, making
// each day slot about eight SVG units wide. It measures the complete generated
// page in both colour schemes, including the hover seam that consumes markIndex.
func TestBrowserBehaviorDurationsDensePanelASpreadStaysBoundedAndInteractive(t *testing.T) {
	const samplesPerDay = 15
	const activeDayCount = 47
	fixtureStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fixtureTickets := make([]*RequestTicket, 0, samplesPerDay*activeDayCount)
	for activeDayIndex := 0; activeDayIndex < activeDayCount; activeDayIndex++ {
		dayStart := fixtureStart.AddDate(0, 0, activeDayIndex*3)
		for sampleIndex := 0; sampleIndex < samplesPerDay; sampleIndex++ {
			completedAt := dayStart.Add(time.Duration(8*60+sampleIndex*37) * time.Minute)
			minutes := time.Duration(2+(sampleIndex*7)%57) * time.Minute
			if sampleIndex < 2 {
				// Ninety-four positive overflow samples force the complete list well
				// past one viewport. Equal values also exercise its REQ-id tie-break.
				minutes = time.Duration(90+(activeDayIndex%5)*15) * time.Minute
			}
			ticket := durationTicket(
				fmt.Sprintf("REQ-%04d", activeDayIndex*samplesPerDay+sampleIndex+1),
				[]string{"A", "B", "C"}[sampleIndex%3],
				completedAt.Add(-minutes).Format(time.RFC3339),
				completedAt.Format(time.RFC3339),
			)
			ticket.UserRequestId = fmt.Sprintf("UR-%03d", activeDayIndex+1)
			ticket.Domain = []string{"frontend", "backend", "testing"}[sampleIndex%3]
			ticket.Title = fmt.Sprintf("Dense sample %d with a wrapping title", activeDayIndex*samplesPerDay+sampleIndex+1)
			fixtureTickets = append(fixtureTickets, ticket)
		}
	}

	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC),
		ProjectName: "REQ-349 dense durations probe",
		AllRequests: fixtureTickets,
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate dense fixture board: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read dense fixture index: %v", readError)
	}

	probeScript := `
  (function () {
  var resultNode = document.createElement("pre");
  resultNode.id = "` + browserProbeResultElementId + `";
  document.body.appendChild(resultNode);
  var durationsPanel = document.getElementById("view-durations");
  durationsPanel.hidden = false;
  durationsPanel.style.display = "block";

  function captureMarks() {
    return Array.from(document.querySelectorAll("#durations-chart circle.durations-mark")).map(function (circle) {
      return { x: Number(circle.getAttribute("cx")), y: Number(circle.getAttribute("cy")) };
    });
  }

  setDurationsWindow("all");
  renderDurationsView();
  var firstMarks = captureMarks();
  setDurationsWindow("all");
  renderDurationsView();
  var secondMarks = captureMarks();
  var samples = boardData.durations.samples;
  var firstSampleMs = Date.parse(samples[0].completionTime);
  var lastSampleMs = Date.parse(samples[samples.length - 1].completionTime);
  var timeStart = Math.floor(firstSampleMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS;
  var timeEnd = Math.floor(lastSampleMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS + DURATIONS_DAY_MS;
  var timeSpan = timeEnd - timeStart;
  var xOfEpoch = function (epochMs) {
    return DURATIONS_MARGIN_LEFT + ((epochMs - timeStart) / timeSpan) * DURATIONS_PLOT_WIDTH;
  };
  var everyMarkInOwnDay = secondMarks.every(function (mark, sampleIndex) {
    var sampleMs = Date.parse(samples[sampleIndex].completionTime);
    var dayStartMs = Math.floor(sampleMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS;
    return mark.x >= xOfEpoch(dayStartMs) - 0.05 && mark.x <= xOfEpoch(dayStartMs + DURATIONS_DAY_MS) + 0.05;
  });
  var firstDayMarks = secondMarks.slice(0, ` + fmt.Sprintf("%d", samplesPerDay) + `);
  var busyDaySpread = Math.max.apply(null, firstDayMarks.map(function (mark) { return mark.x; })) -
    Math.min.apply(null, firstDayMarks.map(function (mark) { return mark.x; }));
  var deterministic = firstMarks.length === secondMarks.length && firstMarks.every(function (mark, markIndex) {
    return mark.x === secondMarks[markIndex].x && mark.y === secondMarks[markIndex].y;
  });

  var ribbon = document.querySelector("#durations-chart .durations-quantile-ribbon");
  var median = document.querySelector("#durations-chart .durations-quantile-median");
  function finiteBounded(pathNode) {
    if (!pathNode || !pathNode.getAttribute("d") || /NaN|Infinity/.test(pathNode.getAttribute("d"))) { return false; }
    var box = pathNode.getBBox();
    return [box.x, box.y, box.width, box.height, pathNode.getTotalLength()].every(Number.isFinite) &&
      box.x >= DURATIONS_MARGIN_LEFT - 0.1 && box.x + box.width <= DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT + 0.1 &&
      box.y >= DURATIONS_MAIN_TOP - 0.1 && box.y + box.height <= DURATIONS_MAIN_BOTTOM + 0.1;
  }

  var hoverMarkIndex = 7;
  var hoveredMark = secondMarks[hoverMarkIndex];
  var svg = document.querySelector("#durations-chart svg");
  var hoverSurface = document.querySelector("#durations-chart .durations-hover-surface");
  var bounds = svg.getBoundingClientRect();
  hoverSurface.dispatchEvent(new MouseEvent("mousemove", {
    bubbles: true,
    clientX: bounds.left + hoveredMark.x * bounds.width / DURATIONS_VIEW_WIDTH,
    clientY: bounds.top + hoveredMark.y * bounds.height / DURATIONS_VIEW_HEIGHT
  }));
  var hoveredRequestId = samples[hoverMarkIndex].id;
  var hoverReadout = document.getElementById("durations-readout").textContent;

  var expectedLongestSpans = samples.filter(function (sample) {
    return sample.wallMinutes > DURATIONS_CEILING_MINUTES;
  }).sort(function (first, second) {
    if (first.wallMinutes !== second.wallMinutes) { return second.wallMinutes - first.wallMinutes; }
    if (first.id !== second.id) { return first.id < second.id ? -1 : 1; }
    return first.completionTime < second.completionTime ? -1 : 1;
  });
  var longestSpanRows = Array.from(document.querySelectorAll("#durations-longest-list > li"));
  var longestSpanOrder = longestSpanRows.every(function (row, rowIndex) {
    return expectedLongestSpans[rowIndex] && row.getAttribute("data-request-id") === expectedLongestSpans[rowIndex].id &&
      Number(row.getAttribute("data-wall-minutes")) === expectedLongestSpans[rowIndex].wallMinutes;
  });
  var everyListField = longestSpanRows.every(function (row, rowIndex) {
    var sample = expectedLongestSpans[rowIndex];
    var request = boardData.requests[sample.id];
    return row.querySelector(".durations-longest-spans-request").textContent === sample.id &&
      row.querySelector(".durations-longest-spans-user-request").textContent === request.userRequestId &&
      row.querySelector(".durations-longest-spans-duration").textContent === formatDurationMinutes(sample.wallMinutes) &&
      row.querySelector(".durations-longest-spans-route").textContent === durationRouteName(sample.route) &&
      row.querySelector(".durations-longest-spans-title").textContent === request.title;
  });
  var overflowSample = expectedLongestSpans[0];
  var overflowSampleIndex = samples.indexOf(overflowSample);
  var overflowMark = secondMarks[overflowSampleIndex];
  hoverSurface.dispatchEvent(new MouseEvent("mousemove", {
    bubbles: true,
    clientX: bounds.left + overflowMark.x * bounds.width / DURATIONS_VIEW_WIDTH,
    clientY: bounds.top + overflowMark.y * bounds.height / DURATIONS_VIEW_HEIGHT
  }));
  var overflowHoverReadout = document.getElementById("durations-readout").textContent;
  var chart = document.getElementById("durations-chart");
  var list = document.getElementById("durations-longest-list");
  var aside = list.closest(".durations-longest-spans");
  var wrapper = chart.parentElement;
  var svgRequestLabels = Array.from(svg.querySelectorAll("text")).filter(function (textNode) {
    return /^REQ-[0-9]+(?:\s|$)/.test(textNode.textContent);
  });
  var svgMoreSentences = Array.from(svg.querySelectorAll("text")).filter(function (textNode) {
    return /^\+[0-9]+ more\b/.test(textNode.textContent);
  });
  var wrapperBounds = wrapper.getBoundingClientRect();
  var asideBounds = aside.getBoundingClientRect();

  document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
    locationHref: location.href,
    renderedSampleCount: secondMarks.length,
    payloadSampleCount: samples.length,
    activeDayCount: boardData.durations.days.length,
    daySlotWidth: DURATIONS_PLOT_WIDTH * DURATIONS_DAY_MS / timeSpan,
    busyDaySpread: busyDaySpread,
    deterministic: deterministic,
    everyMarkInOwnDay: everyMarkInOwnDay,
    hoveredRequestId: hoveredRequestId,
    hoverReadout: hoverReadout,
    ribbonFiniteBounded: finiteBounded(ribbon),
    medianFiniteBounded: finiteBounded(median),
    ribbonOpacity: ribbon ? Number(getComputedStyle(ribbon).opacity) : 0,
    markOpacity: Number(getComputedStyle(document.querySelector("#durations-chart circle.durations-mark:not(.durations-mark-critical):not(.durations-mark-unknown)")).opacity),
    bodyBackground: getComputedStyle(document.body).backgroundColor,
    longestSpanCount: longestSpanRows.length,
    expectedSpanCount: expectedLongestSpans.length,
    longestSpanOrder: longestSpanOrder,
    everyListField: everyListField,
    countSentence: document.getElementById("durations-longest-count").textContent,
    listOutsideSvg: !svg.contains(list),
    sharedWrapper: wrapper === aside.parentElement,
    svgRequestLabels: svgRequestLabels.length,
    svgLeaderLines: svg.querySelectorAll(".durations-label-leader").length,
    svgMoreSentences: svgMoreSentences.length,
    overflowHoverId: overflowSample.id,
    overflowHoverDuration: formatDurationMinutes(overflowSample.wallMinutes),
    overflowHoverReadout: overflowHoverReadout,
    viewportWidth: window.innerWidth,
    documentWidth: document.documentElement.scrollWidth,
    wrapperRight: wrapperBounds.right,
    asideRight: asideBounds.right,
    wrapperColumns: getComputedStyle(wrapper).gridTemplateColumns,
    asideOverflowY: getComputedStyle(aside).overflowY
  });
  })();
`
	// The generated client is one IIFE, so the probe has to run before its final
	// close to exercise the private renderer and payload rather than a copied
	// helper. The result node itself remains ordinary HTML outside the script.
	probePage := string(indexBytes)
	clientClose := strings.LastIndex(probePage, "})();")
	if clientClose < 0 {
		t.Fatal("generated page has no client IIFE close for dense durations probe")
	}
	probePage = probePage[:clientClose] + probeScript + probePage[clientClose:]

	for _, viewport := range []struct {
		name          string
		colourFlag    string
		width         int
		stackedLayout bool
	}{
		{name: "light-320", colourFlag: "--force-light-mode", width: 320, stackedLayout: true},
		{name: "dark-320", colourFlag: "--force-dark-mode", width: 320, stackedLayout: true},
		{name: "light-768", colourFlag: "--force-light-mode", width: 768, stackedLayout: true},
		{name: "dark-768", colourFlag: "--force-dark-mode", width: 768, stackedLayout: true},
		{name: "light-1280", colourFlag: "--force-light-mode", width: 1280},
		{name: "dark-1280", colourFlag: "--force-dark-mode", width: 1280},
	} {
		viewport := viewport
		t.Run(viewport.name, func(t *testing.T) {
			resultJSON := runBrowserBehaviorProbeInDirectory(
				t, "dense durations "+viewport.name, siteDirectory, probePage,
				fmt.Sprintf("--window-size=%d,1200", viewport.width), viewport.colourFlag,
			)
			var result durationsDensePanelProbeResult
			if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
				t.Fatalf("decode dense durations probe: %v\n%s", decodeError, resultJSON)
			}
			if result.PayloadSampleCount != len(fixtureTickets) {
				t.Errorf("payload carries %d samples, want the %d-sample fixture",
					result.PayloadSampleCount, len(fixtureTickets))
			}
			if result.RenderedSampleCount != result.PayloadSampleCount {
				t.Errorf("all-history render drew %d .durations-mark circles for %d payload samples",
					result.RenderedSampleCount, result.PayloadSampleCount)
			}
			if result.RenderedSampleCount < 700 || result.ActiveDayCount != activeDayCount {
				t.Errorf("rendered %d samples across %d active days, want at least 700 across %d",
					result.RenderedSampleCount, result.ActiveDayCount, activeDayCount)
			}
			if result.DaySlotWidth < 7.5 || result.DaySlotWidth > 8.5 {
				t.Errorf("day slot width %.2f, want roughly 8 SVG units", result.DaySlotWidth)
			}
			if !result.EveryMarkInOwnDay || result.BusyDaySpread < 5 {
				t.Errorf("own-day=%v, busy-day spread=%.2f; want bounded useful spread", result.EveryMarkInOwnDay, result.BusyDaySpread)
			}
			if !result.Deterministic {
				t.Error("identical payload moved marks across consecutive renders")
			}
			if !strings.HasPrefix(result.HoverReadout, result.HoveredRequestId+" ·") {
				t.Errorf("hover at %s's jittered centre read %q", result.HoveredRequestId, result.HoverReadout)
			}
			if !result.RibbonFiniteBounded || !result.MedianFiniteBounded {
				t.Errorf("distribution geometry bounded: ribbon=%v median=%v", result.RibbonFiniteBounded, result.MedianFiniteBounded)
			}
			if result.RibbonOpacity <= 0 || result.RibbonOpacity >= result.MarkOpacity {
				t.Errorf("ribbon opacity %.2f is not subordinate to mark opacity %.2f", result.RibbonOpacity, result.MarkOpacity)
			}
			if result.LocationHref == "" || !strings.Contains(result.LocationHref, browserProbePageFileName) {
				t.Errorf("probe measured unnamed page %q", result.LocationHref)
			}
			if result.BodyBackground == "" || result.BodyBackground == "rgba(0, 0, 0, 0)" {
				t.Errorf("%s board has no resolved body background: %q", viewport.name, result.BodyBackground)
			}
			if result.ExpectedSpanCount < 60 || result.LongestSpanCount != result.ExpectedSpanCount {
				t.Errorf("longest-spans list rendered %d of %d overflow samples; fixture must carry at least 60",
					result.LongestSpanCount, result.ExpectedSpanCount)
			}
			if !result.LongestSpanOrder || !result.EveryListField {
				t.Errorf("complete list order=%v fields=%v; want descending/tied order with all five fields",
					result.LongestSpanOrder, result.EveryListField)
			}
			wantCountSentence := fmt.Sprintf("%d spans over 60 minutes in this window; all are listed.", result.ExpectedSpanCount)
			if result.CountSentence != wantCountSentence {
				t.Errorf("count sentence = %q, want %q", result.CountSentence, wantCountSentence)
			}
			if !result.ListOutsideSVG || !result.SharedWrapper {
				t.Errorf("list outside SVG=%v shared wrapper=%v, want adjacent HTML", result.ListOutsideSVG, result.SharedWrapper)
			}
			if result.SVGRequestLabels != 0 || result.SVGLeaderLines != 0 || result.SVGMoreSentences != 0 {
				t.Errorf("SVG still carries %d REQ labels, %d leaders, and %d +N-more sentences",
					result.SVGRequestLabels, result.SVGLeaderLines, result.SVGMoreSentences)
			}
			if !strings.HasPrefix(result.OverflowHoverReadout, result.OverflowHoverId+" ·") ||
				!strings.Contains(result.OverflowHoverReadout, " · "+result.OverflowHoverDuration+" ·") {
				t.Errorf("overflow hover for %s read %q", result.OverflowHoverId, result.OverflowHoverReadout)
			}
			if result.AsideOverflowY != "auto" {
				t.Errorf("longest-spans aside overflow-y = %q, want auto", result.AsideOverflowY)
			}
			if result.WrapperRight > result.ViewportWidth+1 || result.AsideRight > result.ViewportWidth+1 {
				t.Errorf("horizontal clipping at viewport %.0f: document %.0f wrapper-right %.0f aside-right %.0f",
					result.ViewportWidth, result.DocumentWidth, result.WrapperRight, result.AsideRight)
			}
			columnCount := len(strings.Fields(result.WrapperColumns))
			if viewport.stackedLayout && columnCount != 1 {
				t.Errorf("%dpx layout has grid columns %q, want one stacked column", viewport.width, result.WrapperColumns)
			}
			if !viewport.stackedLayout && columnCount != 2 {
				t.Errorf("%dpx layout has grid columns %q, want chart plus aside", viewport.width, result.WrapperColumns)
			}
			t.Logf("%s %s: viewport %.0f, %d/%d longest spans, columns %s, body %s",
				viewport.name, result.LocationHref, result.ViewportWidth, result.LongestSpanCount,
				result.ExpectedSpanCount, result.WrapperColumns, result.BodyBackground)
		})
	}
}

// REQ-266's mechanism, chosen as a check rather than a review convention.
//
// The rule REQ-252 established on the Go side: a browser-measured number must
// name the build it was measured on, because a face is per-browser and an undated
// number reads as timeless fact. `go/parser` cannot reach the JS surface, which is
// why the JS comments drifted into carrying three such numbers with no build
// between them.
//
// THE DISCRIMINATOR IS WHAT THE NUMBER CLAIMS, not whether it is a measurement.
// A measured number cited as EVIDENCE FOR A PAST DECISION is already dated — by
// the REQ it names, which is a stronger anchor than a build string, since the REQ
// carries the whole argument. A measured number presented as CURRENT FACT about
// the face in use is the one that needs a build, and after REQ-292 the honest
// answer for those is not to date them but to delete them: the engine answers at
// test time instead.
//
// So this check enforces the rule in the direction that survives: a measured
// number in a JS comment must sit beside either a REQ reference (evidence) or a
// build name (dated fact). A bare one is neither, and is what this catches.
func TestDurationsJavaScriptCommentsDateTheirMeasurements(t *testing.T) {
	rendererBytes, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	// A measurement claim: a decimal with two or more places, next to a word that
	// makes it a statement about drawn text. Plain geometry constants (13, 11, 9)
	// are declared numbers, not measurements, and are deliberately not matched.
	measurementClaim := regexp.MustCompile(
		`(?i)[0-9]+\.[0-9]{2,}[- ]?(?:unit|px|em)?s?\b[^\n]*?\b(?:measur|ascent|descent|line box|per character|draws)`)
	alternate := regexp.MustCompile(
		`(?i)\b(?:measur|ascent|descent|line box|per character|draws)\b[^\n]*?[0-9]+\.[0-9]{2,}`)
	// Either anchor dates the claim: the REQ that measured it, or the build.
	datingAnchor := regexp.MustCompile(`(?i)REQ-[0-9]{3}|chromium|playwright|firefox|webkit|safari`)

	for lineNumber, line := range strings.Split(string(rendererBytes), "\n") {
		commentStart := strings.Index(line, "//")
		if commentStart < 0 {
			continue
		}
		commentText := line[commentStart:]
		if !measurementClaim.MatchString(commentText) && !alternate.MatchString(commentText) {
			continue
		}
		// The dating anchor may sit anywhere in the surrounding comment block, not
		// only on the claim's own line — a paragraph names its REQ once.
		blockText := durationsSurroundingCommentBlock(string(rendererBytes), lineNumber)
		if datingAnchor.MatchString(blockText) {
			continue
		}
		t.Errorf("web/board-durations.js:%d carries a measured number with no REQ or build to date it:\n  %s\n"+
			"A face is per-browser, so an undated measurement reads as timeless fact. Cite the REQ that "+
			"measured it, name the build, or delete the number and let a browser probe answer at test time.",
			lineNumber+1, strings.TrimSpace(commentText))
	}
}

// durationsSurroundingCommentBlock returns the contiguous run of comment lines the
// given line sits in, so a paragraph that names its REQ once dates every claim in it.
func durationsSurroundingCommentBlock(sourceText string, lineNumber int) string {
	lines := strings.Split(sourceText, "\n")
	isComment := func(index int) bool {
		return index >= 0 && index < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[index]), "//")
	}
	blockStart := lineNumber
	for isComment(blockStart - 1) {
		blockStart--
	}
	blockEnd := lineNumber
	for isComment(blockEnd + 1) {
		blockEnd++
	}
	return strings.Join(lines[blockStart:blockEnd+1], "\n")
}
