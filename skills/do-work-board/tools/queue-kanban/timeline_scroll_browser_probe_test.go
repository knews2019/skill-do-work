package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// The Timeline view's scroll surfaces, measured in a real engine.
//
// WHY A BROWSER. "How many things on this page scroll" is not a fact about any
// string the renderer writes. It is the product of one class's max-height and
// overflow, the board's own overflow-y, and the heights the two boxes end up with
// once real rows are laid out — so the Node behavior lane, which sees the DOM the
// client builds but never a layout, cannot answer it at all. Neither can it answer
// whether the axis pins, where it pins, or whether the virtualized range follows
// a scroll it never observes.
//
// Two nested scroll regions give the reader two scrollbars and a mouse wheel whose
// effect depends on where the pointer happens to sit. REQ-587 left the board as
// the single scroll surface: the chart stops being a scroll box, the time axis
// sticks to the top edge of the board while rows pass underneath it, and every
// virtualization decision reads the board's scroll position minus the chart's
// offset inside it.
//
// SCOPE. "One scroll surface" is asserted about the two elements the REQ's own RED
// expression names, by id, and never about "every element on the page". The table
// under the chart (.timeline-table-scroll) is a third overflow box inside a
// <details> that is closed by default, and it is deliberately still a scroll box.
//
// Following REQ-291's lesson, every number this probe reads is guarded as present
// and positive before any comparison is trusted — a probe whose fixture is too
// short to overflow the fold passes every assertion below while measuring nothing.

// timelineScrollProbeRequestCount decides how tall the chart is. At
// TIMELINE_ROW_HEIGHT = 18 plus one 34px header per user-request group, 200 rows
// in ten groups is 200*18 + 10*34 = 3940px of chart against a board viewport of
// roughly 830px in this probe's 1600x900 window — more than four screenfuls, so
// the chart is unambiguously taller than the fold and the before/after row sets
// below have room to be disjoint.
const timelineScrollProbeRequestCount = 200

// timelineScrollProbeUserRequestCount is why the fixture has groups at all. Each
// REQ joins group (index mod this), and the stamps run steadily backwards, so
// every group holds both recent and old members. That is what makes widening the
// window insert rows ABOVE a mid-list row rather than only appending older ones
// below it — the movement the anchor assertion needs in order to mean anything.
const timelineScrollProbeUserRequestCount = 10

// timelineScrollProbeHoursBetweenRequests spreads the fixture over
// 200 * 36h = 300 days, so the "Last 90 days" chip holds roughly the newest sixty
// REQs and "All days" holds all of them. A fixture inside 90 days would make the
// two windows identical and the anchor assertion vacuous.
const timelineScrollProbeHoursBetweenRequests = 36

// timelineScrollProbeRowsScrollPixels is how far, in ROWS coordinates, the probe
// scrolls the chart before asking what is drawn. It is greater than one board
// viewport (~830px) plus two overscans (TIMELINE_OVERSCAN_ROWS = 4 rows of 18px
// at each end = 144px), which is what makes the rendered-row sets before and
// after provably disjoint rather than merely different.
const timelineScrollProbeRowsScrollPixels = 1500

// timelineScrollProbeAnchorScrollPixels puts the reader several groups down the
// narrow window before the window chips change, so rows really do get inserted
// above the anchor row. Small enough to sit inside the ~1420px the 90-day window
// draws, large enough that the anchor is not the chart's first row.
const timelineScrollProbeAnchorScrollPixels = 500

// timelineScrollSurfaceMeasurement is everything the page reports in one node.
// Real heights ship rather than a pair of booleans, so a failure message says how
// far off the layout was instead of only that it was wrong.
type timelineScrollSurfaceMeasurement struct {
	Browser string `json:"browser"`

	RenderedRowIdsBefore []string `json:"renderedRowIdsBefore"`
	RenderedRowIdsAfter  []string `json:"renderedRowIdsAfter"`

	BoardMainClientHeight float64 `json:"boardMainClientHeight"`
	BoardMainScrollHeight float64 `json:"boardMainScrollHeight"`

	TimelineScrollClientHeight float64 `json:"timelineScrollClientHeight"`
	TimelineScrollScrollHeight float64 `json:"timelineScrollScrollHeight"`

	// TimelineScrollOverflowY is the COMPUTED value, and it is not redundant
	// beside the height comparison above. A leftover `overflow-x: hidden` makes
	// a `visible` on the other axis compute to `auto`, so with the max-height
	// gone the element reports scrollHeight == clientHeight while still being a
	// scroll container. The height comparison cannot see that; this can.
	TimelineScrollOverflowY string `json:"timelineScrollOverflowY"`

	RowsSvgHeight float64 `json:"rowsSvgHeight"`
	AxisHeight    float64 `json:"axisHeight"`
	RowsOffsetPx  float64 `json:"rowsOffsetPx"`

	// RedExpressionResult is the REQ's own RED expression, run on the served
	// board. It returned [true, true] before this change and must return
	// [true, false] after it.
	RedExpressionResult []bool `json:"redExpressionResult"`

	BoardMainPaddingTop string `json:"boardMainPaddingTop"`

	// The board's top padding moves onto whatever the reader sees first, which
	// is not always the view panel. AnomaliesStripVisible says the fixture put a
	// real sibling above the view panel, so the rule is being tested in the
	// arrangement that can break it.
	AnomaliesStripVisible      bool    `json:"anomaliesStripVisible"`
	FirstVisibleChildId        string  `json:"firstVisibleChildId"`
	FirstVisibleChildTopOffset float64 `json:"firstVisibleChildTopOffset"`
	ViewPanelMarginTop         string  `json:"viewPanelMarginTop"`

	// TimelineHeadingTopOffset is the heading's distance from the view panel's
	// own top edge — the view's 16px inset, which the padding move must not eat.
	TimelineHeadingTopOffset float64 `json:"timelineHeadingTopOffset"`

	BoardMainScrollTop float64 `json:"boardMainScrollTop"`
	RowsScrollTop      float64 `json:"rowsScrollTop"`

	// AxisTopOffset is the axis's distance from the board's top inner edge AFTER
	// scrolling. A pinned axis reports ~0; one that scrolled away with its rows
	// reports a large negative number; one pinned to a padded content box
	// reports the padding.
	AxisTopOffset float64 `json:"axisTopOffset"`

	// RowsVisibleAboveAxis counts rows painted in the band between the board's
	// top inner edge and the top of the pinned axis. Anything above zero is a
	// row showing through above the axis, which is what the board's own top
	// padding produced on the Activity view before REQ-585 moved it.
	RowsVisibleAboveAxis int `json:"rowsVisibleAboveAxis"`

	// RenderedRowsInViewportAfter proves the rows moved WITH the scroll rather
	// than being drawn somewhere nobody can see.
	RenderedRowsInViewportAfter int `json:"renderedRowsInViewportAfter"`
}

// timelineAnchorMeasurement covers the GREEN condition no existing probe has an
// equivalent for: the row under the pointer stays put when the window chips
// change. It is the requirement most directly threatened by moving the anchor
// read and the anchor write into the board's coordinate space, so it is pinned
// here rather than checked by hand.
type timelineAnchorMeasurement struct {
	AnchorRowId            string  `json:"anchorRowId"`
	AnchorTopOffsetBefore  float64 `json:"anchorTopOffsetBefore"`
	RowsScrollTopBefore    float64 `json:"rowsScrollTopBefore"`
	BoardScrollTopBefore   float64 `json:"boardScrollTopBefore"`
	DisplayListCountBefore int     `json:"displayListCountBefore"`

	// AnchorInDisplayListAfter reads the "REQs in this window, as a table" body,
	// which lists every row the window admits whatever the virtualizer drew. It
	// separates the two ways an anchor can be missing from the screen: a row the
	// wider window legitimately does not carry (a skipped comparison) from a row
	// the window still carries that the reader was scrolled away from (the
	// failure this assertion exists for).
	AnchorInDisplayListAfter bool    `json:"anchorInDisplayListAfter"`
	AnchorDrawnAfter         bool    `json:"anchorDrawnAfter"`
	AnchorTopOffsetAfter     float64 `json:"anchorTopOffsetAfter"`
	BoardScrollTopAfter      float64 `json:"boardScrollTopAfter"`
	DisplayListCountAfter    int     `json:"displayListCountAfter"`
}

func TestBrowserBehaviorTimelineViewHasOneScrollSurface(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)

	fixtureFiles := make([]verifyFixtureFile, 0, timelineScrollProbeRequestCount+1)
	stampBase := time.Now().UTC().Add(-2 * time.Hour)
	for requestIndex := 0; requestIndex < timelineScrollProbeRequestCount; requestIndex++ {
		requestID := fmt.Sprintf("REQ-%04d", 7000+requestIndex)
		fixtureFiles = append(fixtureFiles, verifyFixtureFile{
			RelativePath: "do-work/archive/" + requestID + "-timeline-scroll.md",
			Content:      timelineScrollProbeFixtureRequest(requestID, stampBase, requestIndex),
		})
	}
	// One terminal REQ with no completion instant anyone can resolve — no
	// completed_at, no commit hash. That is exactly what the completion-anomalies
	// strip exists for, and the strip sits OUTSIDE the view panels, so this is
	// what puts a visible sibling above #view-timeline and makes the moved
	// padding testable in the arrangement that can break it. (#board-findings
	// carries the same class and the same position, so pinning one pins the
	// rule for both; producing a real verify finding costs more and proves
	// nothing extra.)
	//
	// It buys a second sibling for free, and a better one: the same REQ raises a
	// data warning, so board-cards.js inserts the warnings banner as the board's
	// FIRST child. That banner is not in template.html and carries no id, so it
	// is the case a rule that named the strips by id would have missed — the
	// measurement below reports whichever element actually came first.
	fixtureFiles = append(fixtureFiles, verifyFixtureFile{
		RelativePath: "do-work/archive/REQ-7999-timeline-scroll-anomaly.md",
		Content: "---\n" +
			"id: REQ-7999\n" +
			"title: Timeline scroll fixture with no resolvable completion instant\n" +
			"status: completed\n" +
			"user_request: UR-100\n" +
			"created_at: " + stampBase.Add(-30*time.Hour).Format("2006-01-02T15:04:05Z") + "\n" +
			"claimed_at: " + stampBase.Add(-26*time.Hour).Format("2006-01-02T15:04:05Z") + "\n" +
			"---\n",
	})
	repoRoot := writeVerifyFixture(t, fixtureFiles)

	board, buildError := buildBoard(repoRoot, time.Now().UTC(), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatal(buildError)
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, board); generateError != nil {
		t.Fatal(generateError)
	}
	// Feed real finding disclosures through the shipped renderer. Their height
	// changes after Timeline has measured its initial position.
	dataPath := filepath.Join(siteDirectory, "board-data.js")
	dataBytes, err := os.ReadFile(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	dataBytes = append(dataBytes, []byte(`
window.queueKanbanBoardData.verifyFindings = Array.from({length: 20}, function (_, index) {
  return {category: "stale-claim", detail: "Finding " + index,
    remedy: "Review this request and resolve the stale claim. ".repeat(300)};
});
`)...)
	if err := os.WriteFile(dataPath, dataBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatal(readError)
	}

	session := startTrustedInputBrowserSession(
		t, "timeline view scroll surfaces", siteDirectory, string(indexBytes), "--window-size=1600,900")
	defer session.closeBrowserSession()

	session.waitForPageCondition(t, "the Timeline view button",
		`document.querySelector('[data-view-target="timeline"]')`)
	session.evaluateInPage(t,
		`(document.querySelector('[data-view-target="timeline"]').click(), "switched")`)
	// The REQ's RED case names the All days window, so the measurement runs there.
	session.evaluateInPage(t,
		`(document.querySelector('[data-timeline-period="all"]').click(), "all days")`)
	// The summary line, not the rows: the board opens above the chart (the
	// heading, the toolbar and any strip come first), so on a tall enough board
	// the correct number of rows drawn at scroll position 0 is zero. The probe
	// expression below puts the chart's top at the board's top edge before it
	// counts anything.
	session.waitForPageCondition(t, "a rendered timeline",
		`document.getElementById('timeline-summary').textContent.length > 0`)

	var measured timelineScrollSurfaceMeasurement
	session.decodeResult(t, "timeline scroll surfaces",
		session.evaluateInPage(t, timelineScrollSurfaceProbeExpression()), &measured)

	// ---- guards, before any comparison is trusted (REQ-291) ----
	if len(measured.RenderedRowIdsBefore) == 0 {
		t.Fatalf("the Timeline drew no rows — every assertion below would measure a chart nobody sees: %+v", measured)
	}
	if measured.RowsSvgHeight <= measured.BoardMainClientHeight {
		t.Fatalf("the rows SVG is %.1f px against a %.1f px board viewport, so the chart does not reach the fold and this fixture cannot tell one scroll surface from two: %+v",
			measured.RowsSvgHeight, measured.BoardMainClientHeight, measured)
	}
	if measured.BoardMainClientHeight <= 0 || measured.TimelineScrollClientHeight <= 0 || measured.AxisHeight <= 0 {
		t.Fatalf("a measured box has no height: %+v", measured)
	}
	if measured.BoardMainScrollHeight <= measured.BoardMainClientHeight {
		t.Fatalf("the board does not scroll (client %.1f, scroll %.1f) — with nothing scrolling at all, \"one scroll surface\" is indistinguishable from \"none\": %+v",
			measured.BoardMainClientHeight, measured.BoardMainScrollHeight, measured)
	}
	if measured.RowsScrollTop < timelineScrollProbeRowsScrollPixels-1 {
		t.Fatalf("the chart only scrolled %.1f px of the requested %d, so everything measured after the scroll proves nothing: %+v",
			measured.RowsScrollTop, timelineScrollProbeRowsScrollPixels, measured)
	}
	if !measured.AnomaliesStripVisible {
		t.Fatalf("the completion-anomalies strip is not on screen, so the moved top padding is being measured with nothing above the view panel — the arrangement that cannot break: %+v", measured)
	}

	// ---- A1: one scroll surface ----
	if measured.TimelineScrollScrollHeight > measured.TimelineScrollClientHeight+1 {
		t.Errorf("the chart is still its own scroll box (client %.1f, scroll %.1f) — that is the second scrollbar REQ-587 removed: %+v",
			measured.TimelineScrollClientHeight, measured.TimelineScrollScrollHeight, measured)
	}
	if measured.TimelineScrollOverflowY != "visible" {
		t.Errorf("#timeline-scroll computes overflow-y: %q, so it is still a scroll container even though its heights are equal — the half-applied shape a height comparison cannot see: %+v",
			measured.TimelineScrollOverflowY, measured)
	}

	// ---- A2: the REQ's own RED expression ----
	if len(measured.RedExpressionResult) != 2 ||
		!measured.RedExpressionResult[0] || measured.RedExpressionResult[1] {
		t.Errorf("the REQ's RED expression returned %v on the served board, want [true false]: %+v",
			measured.RedExpressionResult, measured)
	}

	// ---- A3: the axis sticks ----
	// One pixel of tolerance, for an offset that rounds off the device pixel
	// ratio rather than landing on an exact zero.
	if measured.AxisTopOffset < -1 || measured.AxisTopOffset > 1 {
		t.Errorf("after scrolling %d px of rows the time axis sits %.1f px from the board's top edge, so it is not pinned to it: %+v",
			timelineScrollProbeRowsScrollPixels, measured.AxisTopOffset, measured)
	}
	if measured.RowsVisibleAboveAxis != 0 {
		t.Errorf("%d rows are drawn above the pinned axis: %+v", measured.RowsVisibleAboveAxis, measured)
	}

	// ---- A4: the virtualized range follows the scroll ----
	if len(measured.RenderedRowIdsAfter) == 0 {
		t.Errorf("the chart drew no rows after scrolling %d px — the virtualized range did not follow the board: %+v",
			timelineScrollProbeRowsScrollPixels, measured)
	}
	if overlap := intersectingRowIds(measured.RenderedRowIdsBefore, measured.RenderedRowIdsAfter); len(overlap) > 0 {
		t.Errorf("%d of the rows drawn before the %d px scroll are still drawn after it (%v) — the range is not following the board's scroll position: %+v",
			len(overlap), timelineScrollProbeRowsScrollPixels, overlap, measured)
	}
	if measured.RenderedRowsInViewportAfter == 0 {
		t.Errorf("no rendered row intersects the board's visible rect after the scroll — the range moved but the rows were drawn where nobody can see them: %+v", measured)
	}

	// ---- A5 / A6: the padding move ----
	if measured.BoardMainPaddingTop != "0px" {
		t.Errorf("the Timeline leaves the board's top padding at %q — a sticky element pins to the CONTENT box, so that padding is a band the rows scroll through above the axis",
			measured.BoardMainPaddingTop)
	}
	if measured.FirstVisibleChildTopOffset < 23 || measured.FirstVisibleChildTopOffset > 25 {
		t.Errorf("the board's first visible child (%s) starts %.1f px below the board's top inner edge, want the 24 px the board's own padding used to supply — it is flush against the top bar: %+v",
			measured.FirstVisibleChildId, measured.FirstVisibleChildTopOffset, measured)
	}
	if measured.ViewPanelMarginTop != "0px" {
		t.Errorf("#view-timeline carries margin-top %q while a visible strip precedes it, so the moved padding is being applied twice: %+v",
			measured.ViewPanelMarginTop, measured)
	}

	// ---- A7: the view's own inset survived ----
	if measured.TimelineHeadingTopOffset <= 0 {
		t.Errorf("the Timeline heading starts %.1f px below the view panel's own top edge, so the view's 16 px inset was eaten by the padding move: %+v",
			measured.TimelineHeadingTopOffset, measured)
	}

	t.Logf("timeline view browser=%s rowsBefore=%d rowsAfter=%d board=%.1f/%.1f chart=%.1f/%.1f overflowY=%s rowsSvg=%.1f axisHeight=%.1f rowsOffset=%.1f red=%v boardScrollTop=%.1f rowsScrollTop=%.1f axisTop=%.1f rowsAboveAxis=%d inViewportAfter=%d boardPaddingTop=%s firstVisibleChild=%s@%.1f viewPanelMarginTop=%s headingTop=%.1f",
		measured.Browser, len(measured.RenderedRowIdsBefore), len(measured.RenderedRowIdsAfter),
		measured.BoardMainClientHeight, measured.BoardMainScrollHeight,
		measured.TimelineScrollClientHeight, measured.TimelineScrollScrollHeight,
		measured.TimelineScrollOverflowY, measured.RowsSvgHeight, measured.AxisHeight,
		measured.RowsOffsetPx, measured.RedExpressionResult,
		measured.BoardMainScrollTop, measured.RowsScrollTop, measured.AxisTopOffset,
		measured.RowsVisibleAboveAxis, measured.RenderedRowsInViewportAfter,
		measured.BoardMainPaddingTop, measured.FirstVisibleChildId,
		measured.FirstVisibleChildTopOffset, measured.ViewPanelMarginTop,
		measured.TimelineHeadingTopOffset)

	assertTimelineAnchorSurvivesAWindowChange(t, session)
	assertTimelineRowsSurviveFindingDisclosures(t, session)

	// The padding move is scoped to this one view, so switching away must restore
	// the board's own top padding. Read rather than restated: the assertion is
	// that the two views disagree and that only the Timeline is zero.
	session.evaluateInPage(t,
		`(document.querySelector('[data-view-target="board"]').click(), "switched")`)
	session.waitForPageCondition(t, "the board view",
		`!document.getElementById('view-board').hidden`)
	var boardViewPaddingTop string
	session.decodeResult(t, "board view padding", session.evaluateInPage(t,
		`getComputedStyle(document.getElementById('board-main')).paddingTop`), &boardViewPaddingTop)
	if boardViewPaddingTop == "0px" {
		t.Errorf("switching back to the Kanban view left the board's top padding at %q, so the Timeline-only rule is not scoped to the Timeline view",
			boardViewPaddingTop)
	}
	t.Logf("board view boardMainPaddingTop=%s", boardViewPaddingTop)
}

func assertTimelineRowsSurviveFindingDisclosures(t *testing.T, session *trustedInputBrowserSession) {
	t.Helper()
	var baselineIDs []string
	var baselineWidth float64
	var previousOffset float64
	for _, step := range []struct {
		name   string
		script string
	}{
		{"baseline", `document.getElementById('board-findings-strip').open = false;
document.querySelector('[data-timeline-period="all"]').click();`},
		{"strip expanded", `document.getElementById('board-findings-strip').open = true;`},
		{"finding expanded", `document.querySelector('#board-findings-cards details').open = true;`},
		{"finding collapsed", `document.querySelector('#board-findings-cards details').open = false;`},
		{"strip collapsed", `document.getElementById('board-findings-strip').open = false;`},
	} {
		var measured struct {
			URL       string   `json:"url"`
			Browser   string   `json:"browser"`
			RowIDs    []string `json:"rowIDs"`
			Offset    float64  `json:"offset"`
			Width     float64  `json:"width"`
			ScrollTop float64  `json:"scrollTop"`
			ToggleIDs []string `json:"toggleIDs"`
			ScrollIDs []string `json:"scrollIDs"`
		}
		session.decodeResult(t, step.name, session.evaluateInPage(t, `(async function () {
  var board = document.getElementById('board-main');
  var chart = document.getElementById('timeline-scroll');
  board.style.overflowAnchor = 'none';
  function rowIDs() { return Array.from(chart.querySelectorAll('[data-detail-id]')).map(function (row) {
    return row.getAttribute('data-detail-id');
  }); }
  function settle() { return new Promise(function (resolve) {
    requestAnimationFrame(function () { requestAnimationFrame(resolve); });
  }); }
  `+step.script+`
  await settle();
  var toggleIDs = rowIDs();
  board.dispatchEvent(new Event('scroll'));
  await settle();
  var scrollIDs = rowIDs();
  var offset = chart.getBoundingClientRect().top - board.getBoundingClientRect().top + board.scrollTop;
  board.scrollTop = offset + 1500;
  // Exercise a scroll even when browser anchoring preserved this position.
  board.dispatchEvent(new Event('scroll'));
  await settle();
  return {url: location.href, browser: navigator.userAgent, offset: offset,
    width: chart.clientWidth, scrollTop: board.scrollTop - offset,
    rowIDs: rowIDs(), toggleIDs: toggleIDs, scrollIDs: scrollIDs};
})()`), &measured)
		if measured.URL == "" || measured.Width <= 0 || measured.ScrollTop < 1499 || len(measured.RowIDs) == 0 {
			t.Fatalf("%s did not measure a scrolled chart: %+v", step.name, measured)
		}
		if baselineIDs == nil {
			baselineIDs = measured.RowIDs
			baselineWidth = measured.Width
		} else {
			if measured.Width != baselineWidth {
				t.Fatalf("%s changed width, which can mask the offset regression: %+v", step.name, measured)
			}
			if !reflect.DeepEqual(measured.ToggleIDs, measured.ScrollIDs) {
				t.Errorf("%s left stale rows until scrolling: toggle=%v, scroll=%v", step.name, measured.ToggleIDs, measured.ScrollIDs)
			}
			if measured.Offset-previousOffset < 100 && previousOffset-measured.Offset < 100 {
				t.Fatalf("%s did not move the chart enough to expose stale geometry: %+v", step.name, measured)
			}
			if !reflect.DeepEqual(measured.RowIDs, baselineIDs) {
				t.Errorf("%s virtualized different rows at the same chart scroll position: got %v, want %v", step.name, measured.RowIDs, baselineIDs)
			}
		}
		previousOffset = measured.Offset
		t.Logf("%s: page=%s browser=%s offset=%.1f width=%.1f scroll=%.1f rows=%d",
			step.name, measured.URL, measured.Browser, measured.Offset, measured.Width, measured.ScrollTop, len(measured.RowIDs))
	}
}

// assertTimelineAnchorSurvivesAWindowChange scrolls into a NARROW window, notes
// the row at the top of the viewport, then widens the window to All days — which
// inserts older members into every group above that row — and asserts the row is
// still in the same place on screen.
//
// Widening rather than narrowing is deliberate: every row the narrow window
// carries is also in the wide one, so the anchor cannot legitimately disappear,
// and the display-list guard below can only fire on a real regression.
func assertTimelineAnchorSurvivesAWindowChange(t *testing.T, session *trustedInputBrowserSession) {
	t.Helper()

	session.evaluateInPage(t,
		`(document.querySelector('[data-view-target="timeline"]').click(), "timeline")`)
	// The table under the chart lists every row the window admits, whatever the
	// virtualizer happens to have drawn. It only rebuilds while it is open.
	session.evaluateInPage(t,
		`(document.querySelector('#view-timeline .timeline-table').open = true, "opened")`)
	session.evaluateInPage(t,
		`(document.querySelector('[data-timeline-period="90"]').click(), "last 90 days")`)
	session.waitForPageCondition(t, "the narrow window's rows and table",
		`document.querySelectorAll('#timeline-scroll [data-detail-id]').length > 0 &&`+
			` document.querySelectorAll('#timeline-table-body [data-timeline-table-request]').length > 0`)

	var anchor timelineAnchorMeasurement
	session.decodeResult(t, "timeline anchor before the window change",
		session.evaluateInPage(t, timelineAnchorBeforeProbeExpression()), &anchor)

	if anchor.AnchorRowId == "" {
		t.Fatalf("no row sits at the top of the viewport in the narrow window, so there is no anchor to follow: %+v", anchor)
	}
	if anchor.RowsScrollTopBefore < timelineScrollProbeAnchorScrollPixels-1 {
		t.Fatalf("the chart only scrolled %.1f px of the requested %d, so the anchor row is near the top of the list and nothing can be inserted above it: %+v",
			anchor.RowsScrollTopBefore, timelineScrollProbeAnchorScrollPixels, anchor)
	}

	session.evaluateInPage(t,
		`(document.querySelector('[data-timeline-period="all"]').click(), "all days")`)
	session.waitForPageCondition(t, "the wider window's table",
		fmt.Sprintf(`document.querySelectorAll('#timeline-table-body [data-timeline-table-request]').length > %d`,
			anchor.DisplayListCountBefore))

	var after timelineAnchorMeasurement
	session.decodeResult(t, "timeline anchor after the window change",
		session.evaluateInPage(t, timelineAnchorAfterProbeExpression(anchor.AnchorRowId)), &after)

	if after.DisplayListCountAfter <= anchor.DisplayListCountBefore {
		t.Fatalf("the window change added no rows (%d then %d), so nothing was inserted above the anchor and this comparison proves nothing: before=%+v after=%+v",
			anchor.DisplayListCountBefore, after.DisplayListCountAfter, anchor, after)
	}
	// The guard that keeps a legitimately-absent anchor from reading as a false
	// pass — and, here, from reading as a failure either.
	if !after.AnchorInDisplayListAfter {
		t.Skipf("%s left the window when it widened, so its screen position is not comparable — this fixture is built so that cannot happen, so treat a skip here as a fixture bug: before=%+v after=%+v",
			anchor.AnchorRowId, anchor, after)
	}
	if !after.AnchorDrawnAfter {
		t.Fatalf("%s is still in the window after the chips changed but is not drawn on screen — the reader was carried away from the row they were reading: before=%+v after=%+v",
			anchor.AnchorRowId, anchor, after)
	}
	// Two pixels of tolerance: the write lands on a whole pixel of scroll and the
	// row's own rect carries the device pixel ratio's rounding.
	positionDrift := after.AnchorTopOffsetAfter - anchor.AnchorTopOffsetBefore
	if positionDrift < -2 || positionDrift > 2 {
		t.Errorf("%s moved %.1f px on screen when the window chips changed (%.1f then %.1f) — the row under the pointer did not stay put: before=%+v after=%+v",
			anchor.AnchorRowId, positionDrift,
			anchor.AnchorTopOffsetBefore, after.AnchorTopOffsetAfter, anchor, after)
	}

	t.Logf("timeline anchor id=%s rows=%d->%d boardScrollTop=%.1f->%.1f anchorTopOffset=%.1f->%.1f drift=%.1f",
		anchor.AnchorRowId, anchor.DisplayListCountBefore, after.DisplayListCountAfter,
		anchor.BoardScrollTopBefore, after.BoardScrollTopAfter,
		anchor.AnchorTopOffsetBefore, after.AnchorTopOffsetAfter, positionDrift)
}

// intersectingRowIds reports the ids present in both sets, so a failure names the
// rows that did not move instead of only counting them.
func intersectingRowIds(firstRowIds []string, secondRowIds []string) []string {
	presentInFirst := make(map[string]bool, len(firstRowIds))
	for _, rowID := range firstRowIds {
		presentInFirst[rowID] = true
	}
	var shared []string
	for _, rowID := range secondRowIds {
		if presentInFirst[rowID] {
			shared = append(shared, rowID)
		}
	}
	return shared
}

// timelineScrollProbeFixtureRequest writes one archived REQ with a full set of
// lifecycle stamps, so it draws a wait segment and a work segment on the chart.
// The user request cycles so every group holds both recent and old members.
func timelineScrollProbeFixtureRequest(requestID string, stampBase time.Time, requestIndex int) string {
	completedAt := stampBase.Add(
		-time.Duration(requestIndex*timelineScrollProbeHoursBetweenRequests) * time.Hour)
	claimedAt := completedAt.Add(-4 * time.Hour)
	createdAt := claimedAt.Add(-8 * time.Hour)
	stamp := func(instant time.Time) string { return instant.Format("2006-01-02T15:04:05Z") }
	return "---\n" +
		"id: " + requestID + "\n" +
		"title: Timeline scroll fixture " + requestID + "\n" +
		"status: completed\n" +
		fmt.Sprintf("user_request: UR-%d\n", 100+requestIndex%timelineScrollProbeUserRequestCount) +
		"created_at: " + stamp(createdAt) + "\n" +
		"claimed_at: " + stamp(claimedAt) + "\n" +
		"completed_at: " + stamp(completedAt) + "\n" +
		"---\n"
}

// timelineScrollSurfaceProbeExpression measures both boxes and the moved padding
// before scrolling, scrolls the chart, and then reports where the axis and the
// rows ended up.
//
// The two animation frames are not a sleep: setting scrollTop and asking for a
// rectangle in the same task can read the sticky axis's pre-scroll position,
// because the compositor has not repositioned it yet.
func timelineScrollSurfaceProbeExpression() string {
	return `(function () {
  var boardMain = document.getElementById('board-main');
  var timelineScroll = document.getElementById('timeline-scroll');
  var timelineAxis = document.getElementById('timeline-axis');
  var viewPanel = document.getElementById('view-timeline');
  var heading = document.querySelector('#view-timeline .timeline-heading');
  var anomaliesStrip = document.getElementById('board-anomalies');
  // Two animation frames after every scrollTop write, and they are not a sleep:
  // setting scrollTop and asking for a rectangle in the same task can read the
  // sticky axis's pre-scroll position, because the compositor has not
  // repositioned it yet.
  function afterTwoFrames() {
    return new Promise(function (resolve) {
      requestAnimationFrame(function () { requestAnimationFrame(resolve); });
    });
  }
  function renderedRowNodes() {
    return Array.prototype.slice.call(
      document.querySelectorAll('#timeline-scroll [data-detail-id]'));
  }
  function renderedRowIds() {
    return renderedRowNodes().map(function (rowNode) {
      return rowNode.getAttribute('data-detail-id');
    });
  }
  function boardInnerTopNow() {
    return boardMain.getBoundingClientRect().top +
      parseFloat(getComputedStyle(boardMain).paddingTop);
  }
  var measurement = { browser: navigator.userAgent };
  // The distance from the board's scroll origin to the top of the rows box. It
  // is a fact about the laid-out page, read here only to decide where to scroll
  // and what to report — every assertion is about positions the page produced,
  // not about this number.
  var rowsOffsetPx = timelineScroll.getBoundingClientRect().top -
    boardMain.getBoundingClientRect().top + boardMain.scrollTop;
  measurement.rowsOffsetPx = rowsOffsetPx;

  // ---- unscrolled: the moved padding and the view's own inset ----
  boardMain.scrollTop = 0;
  return afterTwoFrames().then(function () {
    var boardStyle = getComputedStyle(boardMain);
    var boardInnerTop = boardInnerTopNow();
    var firstVisibleChild = Array.prototype.filter.call(
      boardMain.children, function (childNode) { return !childNode.hidden; })[0];
    measurement.boardMainPaddingTop = boardStyle.paddingTop;
    measurement.anomaliesStripVisible = !anomaliesStrip.hidden;
    measurement.firstVisibleChildId = firstVisibleChild
      ? (firstVisibleChild.id || firstVisibleChild.className || firstVisibleChild.tagName)
      : '';
    measurement.firstVisibleChildTopOffset = firstVisibleChild
      ? firstVisibleChild.getBoundingClientRect().top - boardInnerTop : -1;
    measurement.viewPanelMarginTop = getComputedStyle(viewPanel).marginTop;
    measurement.timelineHeadingTopOffset =
      heading.getBoundingClientRect().top - viewPanel.getBoundingClientRect().top;

    // ---- the chart's top against the board's top edge ----
    // Not scroll position 0: the heading, the toolbar and any strip come first,
    // so on a board this tall the correct number of rows drawn at 0 is zero.
    boardMain.scrollTop = rowsOffsetPx;
    return afterTwoFrames();
  }).then(function () {
    var rowsSvg = timelineScroll.querySelector('.timeline-rows-svg');
    measurement.renderedRowIdsBefore = renderedRowIds();
    measurement.boardMainClientHeight = boardMain.clientHeight;
    measurement.boardMainScrollHeight = boardMain.scrollHeight;
    measurement.timelineScrollClientHeight = timelineScroll.clientHeight;
    measurement.timelineScrollScrollHeight = timelineScroll.scrollHeight;
    measurement.timelineScrollOverflowY = getComputedStyle(timelineScroll).overflowY;
    measurement.rowsSvgHeight = rowsSvg ? rowsSvg.getBoundingClientRect().height : 0;
    measurement.axisHeight = timelineAxis.getBoundingClientRect().height;
    // The REQ's own RED expression, character for character apart from the
    // return the wrapper needs. Running anything else is running a different
    // test than the one this REQ will be judged by.
    measurement.redExpressionResult = (function () {
      const m=document.querySelector('.board-main'), t=document.getElementById('timeline-scroll');
      return [m.scrollHeight>m.clientHeight, t.scrollHeight>t.clientHeight];
    })();

    boardMain.scrollTop = rowsOffsetPx + 1500;
    return afterTwoFrames();
  }).then(function () {
    var boardRect = boardMain.getBoundingClientRect();
    var axisRect = timelineAxis.getBoundingClientRect();
    var boardInnerTop = boardInnerTopNow();
    measurement.boardMainScrollTop = boardMain.scrollTop;
    measurement.rowsScrollTop = boardMain.scrollTop - rowsOffsetPx;
    measurement.axisTopOffset = axisRect.top - boardInnerTop;
    measurement.renderedRowIdsAfter = renderedRowIds();
    var drawnRows = renderedRowNodes();
    measurement.rowsVisibleAboveAxis = drawnRows.filter(function (rowNode) {
      var rowRect = rowNode.getBoundingClientRect();
      return Math.min(rowRect.bottom, axisRect.top) -
        Math.max(rowRect.top, boardInnerTop) > 0.5;
    }).length;
    measurement.renderedRowsInViewportAfter = drawnRows.filter(function (rowNode) {
      var rowRect = rowNode.getBoundingClientRect();
      return Math.min(rowRect.bottom, boardRect.bottom) -
        Math.max(rowRect.top, axisRect.bottom) > 0.5;
    }).length;
    return measurement;
  });
})()`
}

// timelineAnchorBeforeProbeExpression scrolls the narrow window and reports the
// row at the top of the viewport, picked by the same rule the renderer's own
// anchor uses: the first row in draw order whose bottom edge is below the board's
// top inner edge.
func timelineAnchorBeforeProbeExpression() string {
	return `(function () {
  var boardMain = document.getElementById('board-main');
  var timelineScroll = document.getElementById('timeline-scroll');
  var rowsOffsetPx = timelineScroll.getBoundingClientRect().top -
    boardMain.getBoundingClientRect().top + boardMain.scrollTop;
  boardMain.scrollTop = rowsOffsetPx + ` + strconv.Itoa(timelineScrollProbeAnchorScrollPixels) + `;
  return new Promise(function (resolve) {
    requestAnimationFrame(function () {
      requestAnimationFrame(function () {
        var boardRect = boardMain.getBoundingClientRect();
        var boardInnerTop = boardRect.top + parseFloat(getComputedStyle(boardMain).paddingTop);
        var anchorRow = Array.prototype.filter.call(
          document.querySelectorAll('#timeline-scroll [data-detail-id]'),
          function (rowNode) {
            return rowNode.getBoundingClientRect().bottom > boardInnerTop + 0.5;
          })[0];
        resolve({
          anchorRowId: anchorRow ? anchorRow.getAttribute('data-detail-id') : '',
          anchorTopOffsetBefore: anchorRow
            ? anchorRow.getBoundingClientRect().top - boardInnerTop : 0,
          rowsScrollTopBefore: boardMain.scrollTop - rowsOffsetPx,
          boardScrollTopBefore: boardMain.scrollTop,
          displayListCountBefore: document.querySelectorAll(
            '#timeline-table-body [data-timeline-table-request]').length
        });
      });
    });
  });
})()`
}

// timelineAnchorAfterProbeExpression reports where the named row ended up once
// the window widened, and separately whether the window still carries it at all.
func timelineAnchorAfterProbeExpression(anchorRowID string) string {
	quotedAnchorRowID := strconv.Quote(anchorRowID)
	return `(function () {
  var boardMain = document.getElementById('board-main');
  var anchorRowId = ` + quotedAnchorRowID + `;
  return new Promise(function (resolve) {
    requestAnimationFrame(function () {
      requestAnimationFrame(function () {
        var boardRect = boardMain.getBoundingClientRect();
        var boardInnerTop = boardRect.top + parseFloat(getComputedStyle(boardMain).paddingTop);
        var drawnAnchor = document.querySelector(
          '#timeline-scroll [data-detail-id="' + anchorRowId + '"]');
        resolve({
          anchorInDisplayListAfter: !!document.querySelector(
            '#timeline-table-body [data-timeline-table-request="' + anchorRowId + '"]'),
          anchorDrawnAfter: !!drawnAnchor,
          anchorTopOffsetAfter: drawnAnchor
            ? drawnAnchor.getBoundingClientRect().top - boardInnerTop : 0,
          boardScrollTopAfter: boardMain.scrollTop,
          displayListCountAfter: document.querySelectorAll(
            '#timeline-table-body [data-timeline-table-request]').length
        });
      });
    });
  });
})()`
}
