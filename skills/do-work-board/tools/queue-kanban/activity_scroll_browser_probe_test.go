package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// The Activity view's scroll surfaces, measured in a real engine.
//
// WHY A BROWSER. "How many things on this page scroll" is not a fact about any
// string the renderer writes. It is the product of one class's max-height and
// overflow, the board's own overflow-y, and the heights the two boxes end up
// with once real rows are laid out — so the Node behavior lane, which sees the
// DOM the client builds but never a layout, cannot answer it at all. The
// question is settled by comparing clientHeight with scrollHeight on both
// boxes, which only an engine that laid the page out can report.
//
// Two nested scroll regions give the reader two scrollbars and a mouse wheel
// whose effect depends on where the pointer happens to sit. REQ-585 left the
// board as the single scroll surface: the transitions table stops being a
// scroll box, and its column header sticks to the top edge of the board while
// rows pass underneath it.
//
// Following REQ-291's lesson, every number this probe reads is asserted to be
// present and positive before any comparison is trusted — a probe that measured
// an unrendered box would otherwise report a table that "does not scroll"
// because it has no height at all.

// activityScrollProbeRequestCount and the stamps below decide how tall the
// transitions table is. Three parseable lifecycle stamps per REQ means the row
// count is three times this, which at the probe's 900 px window overflows both
// the board and the pre-REQ-585 70vh table box — the shape the measurement
// needs in order to distinguish them.
const activityScrollProbeRequestCount = 40

// activityScrollProbeBoardScrollPixels is the distance the probe scrolls the
// board before asking where the column header ended up. It is far enough that
// the header has left its unscrolled position and many rows have passed under
// it, which is the state the REQ's GREEN condition describes.
const activityScrollProbeBoardScrollPixels = 700

// activityScrollSurfaceMeasurement is everything the page reports back in one
// node. Both boxes' heights ship rather than a pair of booleans, so a failure
// message says how far off the layout was instead of only that it was wrong.
type activityScrollSurfaceMeasurement struct {
	Browser string `json:"browser"`

	RenderedRowCount int `json:"renderedRowCount"`

	BoardMainClientHeight float64 `json:"boardMainClientHeight"`
	BoardMainScrollHeight float64 `json:"boardMainScrollHeight"`
	TableClientHeight     float64 `json:"tableClientHeight"`
	TableScrollHeight     float64 `json:"tableScrollHeight"`

	// SummaryTextTopOffset is where the summary line's TEXT starts, measured
	// from the board's top inner edge before any scrolling. The padding this
	// REQ moved off the scroll container and onto the summary is supposed to
	// be invisible, so this is the number that says whether the reader can
	// tell. It is the summary's content-box top, not its border-box top: the
	// moved padding lives inside the summary's own box, so a border-box read
	// reports zero and looks like the padding vanished.
	SummaryTextTopOffset float64 `json:"summaryTextTopOffset"`

	BoardMainScrollTop float64 `json:"boardMainScrollTop"`

	// HeaderTopOffset is the column header's distance from the board's top
	// inner edge AFTER scrolling. A stuck header reports ~0; a header that
	// scrolled away with its rows reports a negative number.
	HeaderTopOffset float64 `json:"headerTopOffset"`
	HeaderHeight    float64 `json:"headerHeight"`

	// RowsVisibleAboveHeader counts body rows that paint inside the band
	// between the board's top inner edge and the top of the stuck header.
	// Anything above zero is a row showing through above the header. A row
	// that merely straddles the board's top edge does NOT count: the header is
	// opaque and covers it, which is the whole point of the sticky rule.
	RowsVisibleAboveHeader int `json:"rowsVisibleAboveHeader"`

	BoardMainPaddingTop string `json:"boardMainPaddingTop"`
	ActivityPaddingTop  string `json:"activityPaddingTop"`
}

func TestBrowserBehaviorActivityViewHasOneScrollSurface(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)

	fixtureFiles := make([]verifyFixtureFile, 0, activityScrollProbeRequestCount)
	stampBase := time.Now().UTC().Add(-2 * time.Hour)
	for requestIndex := 0; requestIndex < activityScrollProbeRequestCount; requestIndex++ {
		requestID := fmt.Sprintf("REQ-%03d", 600+requestIndex)
		fixtureFiles = append(fixtureFiles, verifyFixtureFile{
			RelativePath: "do-work/archive/" + requestID + "-activity-scroll.md",
			Content:      activityScrollProbeFixtureRequest(requestID, stampBase, requestIndex),
		})
	}
	repoRoot := writeVerifyFixture(t, fixtureFiles)

	board, buildError := buildBoard(repoRoot, time.Now().UTC(), defaultRecentWindow, nil)
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

	session := startTrustedInputBrowserSession(
		t, "activity view scroll surfaces", siteDirectory, string(indexBytes), "--window-size=1600,900")
	defer session.closeBrowserSession()

	session.waitForPageCondition(t, "the Activity view button",
		`document.querySelector('[data-view-target="activity"]')`)
	session.evaluateInPage(t,
		`(document.querySelector('[data-view-target="activity"]').click(), "switched")`)

	wantedRowCount := activityScrollProbeRequestCount * 3
	session.waitForPageCondition(t, "every transition row",
		fmt.Sprintf(`document.querySelectorAll('.activity-table-scroll tbody tr').length === %d`,
			wantedRowCount))

	var measured activityScrollSurfaceMeasurement
	session.decodeResult(t, "activity scroll surfaces",
		session.evaluateInPage(t, activityScrollSurfaceProbeExpression()), &measured)

	// REQ-291's lesson: refuse to reason about a box that was never laid out.
	// Without this, an empty or display:none table reports 0 === 0 and reads as
	// the fix working.
	if measured.RenderedRowCount != wantedRowCount {
		t.Fatalf("the Activity view drew %d rows, want %d — the rest of this probe would measure a table nobody sees",
			measured.RenderedRowCount, wantedRowCount)
	}
	if measured.TableClientHeight <= 0 || measured.BoardMainClientHeight <= 0 || measured.HeaderHeight <= 0 {
		t.Fatalf("a measured box has no height: %+v", measured)
	}

	if measured.BoardMainScrollHeight <= measured.BoardMainClientHeight {
		t.Fatalf("the board does not scroll (client %.1f, scroll %.1f) — the fixture is too short to tell one scroll surface from two: %+v",
			measured.BoardMainClientHeight, measured.BoardMainScrollHeight, measured)
	}
	if measured.TableScrollHeight > measured.TableClientHeight {
		t.Errorf("the transitions table is still its own scroll box (client %.1f, scroll %.1f) — that is the second scrollbar REQ-585 removed: %+v",
			measured.TableClientHeight, measured.TableScrollHeight, measured)
	}

	if measured.BoardMainScrollTop < activityScrollProbeBoardScrollPixels {
		t.Fatalf("the board only scrolled to %.1f of the requested %d px, so the header position below proves nothing: %+v",
			measured.BoardMainScrollTop, activityScrollProbeBoardScrollPixels, measured)
	}
	// One pixel of tolerance, for a stuck header whose offset rounds off the
	// device pixel ratio rather than landing on an exact zero.
	if measured.HeaderTopOffset < -1 || measured.HeaderTopOffset > 1 {
		t.Errorf("after scrolling %d px the column header sits %.1f px from the board's top edge, so it is not stuck to it: %+v",
			activityScrollProbeBoardScrollPixels, measured.HeaderTopOffset, measured)
	}
	if measured.RowsVisibleAboveHeader != 0 {
		t.Errorf("%d transition rows are drawn above the stuck column header: %+v",
			measured.RowsVisibleAboveHeader, measured)
	}
	if measured.SummaryTextTopOffset <= 0 {
		t.Errorf("the summary line's text starts %.1f px from the board's top edge, so the padding this REQ moved onto it did not arrive and the view now opens flush against the top bar: %+v",
			measured.SummaryTextTopOffset, measured)
	}

	t.Logf("activity view browser=%s rows=%d board=%.1f/%.1f table=%.1f/%.1f scrollTop=%.1f headerTop=%.1f headerHeight=%.1f rowsAboveHeader=%d summaryTextTop=%.1f boardMainPaddingTop=%s activityPaddingTop=%s",
		measured.Browser, measured.RenderedRowCount,
		measured.BoardMainClientHeight, measured.BoardMainScrollHeight,
		measured.TableClientHeight, measured.TableScrollHeight,
		measured.BoardMainScrollTop, measured.HeaderTopOffset, measured.HeaderHeight,
		measured.RowsVisibleAboveHeader, measured.SummaryTextTopOffset,
		measured.BoardMainPaddingTop, measured.ActivityPaddingTop)

	// The padding move is scoped to this one view, so switching back must
	// restore the board's own top padding. Read rather than restated: the
	// assertion is that the two views disagree and that only Activity is zero.
	if measured.BoardMainPaddingTop != "0px" {
		t.Errorf("the Activity view leaves the board's top padding at %q — that padding is scrollable content and rows pass through it above the stuck header",
			measured.BoardMainPaddingTop)
	}
	session.evaluateInPage(t,
		`(document.querySelector('[data-view-target="board"]').click(), "switched")`)
	session.waitForPageCondition(t, "the board view",
		`!document.getElementById('view-board').hidden`)
	var boardViewPaddingTop string
	session.decodeResult(t, "board view padding", session.evaluateInPage(t,
		`getComputedStyle(document.getElementById('board-main')).paddingTop`), &boardViewPaddingTop)
	if boardViewPaddingTop == "0px" {
		t.Errorf("switching back to the Kanban view left the board's top padding at %q, so the Activity-only rule is not scoped to the Activity view",
			boardViewPaddingTop)
	}
	t.Logf("board view boardMainPaddingTop=%s", boardViewPaddingTop)
}

// activityScrollProbeFixtureRequest writes one archived REQ carrying three
// parseable lifecycle stamps inside the Activity view's default 24-hour window,
// so it contributes three transition rows.
func activityScrollProbeFixtureRequest(requestID string, stampBase time.Time, requestIndex int) string {
	stampAt := func(minutesEarlier int) string {
		return stampBase.Add(-time.Duration(requestIndex*3+minutesEarlier) * time.Minute).
			Format("2006-01-02T15:04:05Z")
	}
	return "---\n" +
		"id: " + requestID + "\n" +
		"title: Activity scroll fixture " + requestID + "\n" +
		"status: completed\n" +
		"created_at: " + stampAt(2) + "\n" +
		"claimed_at: " + stampAt(1) + "\n" +
		"completed_at: " + stampAt(0) + "\n" +
		"---\n"
}

// activityScrollSurfaceProbeExpression measures both boxes before scrolling,
// scrolls the board, and then reports where the sticky column header landed.
//
// The two animation frames are not a sleep: setting scrollTop and asking for a
// rectangle in the same task can read the sticky header's pre-scroll position,
// because the compositor has not repositioned it yet.
func activityScrollSurfaceProbeExpression() string {
	return `(function () {
  var boardMain = document.getElementById('board-main');
  var tableScroll = document.querySelector('.activity-table-scroll');
  var headerCell = document.querySelector('.activity-table-scroll thead th');
  var summaryLine = document.getElementById('activity-summary');
  var bodyRows = Array.prototype.slice.call(
    document.querySelectorAll('.activity-table-scroll tbody tr'));
  var unscrolledBoardRect = boardMain.getBoundingClientRect();
  var measurement = {
    browser: navigator.userAgent,
    renderedRowCount: bodyRows.length,
    boardMainClientHeight: boardMain.clientHeight,
    boardMainScrollHeight: boardMain.scrollHeight,
    tableClientHeight: tableScroll.clientHeight,
    tableScrollHeight: tableScroll.scrollHeight,
    summaryTextTopOffset: summaryLine.getBoundingClientRect().top - unscrolledBoardRect.top +
      parseFloat(getComputedStyle(summaryLine).paddingTop),
    boardMainPaddingTop: getComputedStyle(boardMain).paddingTop,
    activityPaddingTop: getComputedStyle(document.getElementById('view-activity')).paddingTop
  };
  boardMain.scrollTop = ` + strconv.Itoa(activityScrollProbeBoardScrollPixels) + `;
  return new Promise(function (resolve) {
    requestAnimationFrame(function () {
      requestAnimationFrame(function () {
        var boardRect = boardMain.getBoundingClientRect();
        var headerRect = headerCell.getBoundingClientRect();
        measurement.boardMainScrollTop = boardMain.scrollTop;
        measurement.headerTopOffset = headerRect.top - boardRect.top;
        measurement.headerHeight = headerRect.height;
        measurement.rowsVisibleAboveHeader = bodyRows.filter(function (bodyRow) {
          var rowRect = bodyRow.getBoundingClientRect();
          var paintedInBand = Math.min(rowRect.bottom, headerRect.top) -
            Math.max(rowRect.top, boardRect.top);
          return paintedInBand > 0.5;
        }).length;
        resolve(measurement);
      });
    });
  });
})()`
}
