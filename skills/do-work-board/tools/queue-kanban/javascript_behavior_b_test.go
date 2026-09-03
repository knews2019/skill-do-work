package main

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestJavaScriptBehaviorDrawerHeadingDeduplication(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function normalizeHeadingText("),
		sliceBalancedBlockAfter(t, indexHtml, "function linkifyDetailBody("),
	}
	javascriptProbe := `
var NodeFilter = { SHOW_TEXT: 4 };
var document = {
  createTreeWalker: function () {
    return { nextNode: function () { return false; } };
  }
};
function makeDrawerBody(headingText) {
  var bodyRoot = {
    firstElementChild: null,
    querySelectorAll: function () { return []; }
  };
  var heading = {
    tagName: "H1",
    textContent: headingText,
    removed: false,
    remove: function () {
      this.removed = true;
      bodyRoot.firstElementChild = null;
    }
  };
  bodyRoot.firstElementChild = heading;
  return { root: bodyRoot, heading: heading };
}
` + strings.Join(functionBlocks, "\n") + `
var requestMatch = makeDrawerBody("  Compile   assets ");
var requestMismatch = makeDrawerBody("Implementation notes");
var userRequestMatch = makeDrawerBody("Launch plan");
var userRequestMismatch = makeDrawerBody("Background");
linkifyDetailBody(requestMatch.root, "Compile assets");
linkifyDetailBody(requestMismatch.root, "Compile assets");
linkifyDetailBody(userRequestMatch.root, "LAUNCH PLAN");
linkifyDetailBody(userRequestMismatch.root, "Launch plan");
process.stdout.write(JSON.stringify([
  requestMatch.heading.removed,
  requestMismatch.heading.removed,
  userRequestMatch.heading.removed,
  userRequestMismatch.heading.removed
]));`
	probeOutput := runJavaScriptBehaviorProbe(t, "drawer heading", javascriptProbe)
	var removedResults []bool
	if decodeError := json.Unmarshal(probeOutput, &removedResults); decodeError != nil {
		t.Fatalf("decode drawer heading behavior: %v (output %q)", decodeError, probeOutput)
	}
	wantedResults := []bool{true, false, true, false}
	if len(removedResults) != len(wantedResults) {
		t.Fatalf("drawer heading result count = %d, want %d: %#v", len(removedResults), len(wantedResults), removedResults)
	}
	for resultIndex := range wantedResults {
		if removedResults[resultIndex] != wantedResults[resultIndex] {
			t.Fatalf("drawer heading result[%d] = %v, want %v; all results=%#v",
				resultIndex, removedResults[resultIndex], wantedResults[resultIndex], removedResults)
		}
	}
}

func TestJavaScriptBehaviorByUserRequestLensCountsRecentlyDoneAsActive(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, requiredToken := range []string{
		"userRequestHasOpenOrRecentWork(userRequest, recentlyDoneIdSet)",
		"recentlyDoneIds(viewState.windowHours)",
	} {
		if !strings.Contains(indexHtml, requiredToken) {
			t.Fatalf("by-UR recent-work behavior is not wired into the generated asset: %q missing", requiredToken)
		}
	}

	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestHasOpenOrRecentWork("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentlyDoneIds("),
	}
	// The calendar carries EVERY REQ, so the stub holds a claimed-an-hour-ago and
	// a failed-an-hour-ago entry alongside the completions. Both are inside the
	// 24h window and neither may reach Recently done: recentlyDoneIds is gated on
	// terminal-resolved, not on mere presence in the array.
	javascriptProbe := `
Date.now = function () { return Date.parse("2026-08-15T12:00:00Z"); };
var boardData = { calendar: [
  { id: "REQ-claimed", status: "claimed", entryTime: "2026-08-15T11:00:00Z" },
  { id: "REQ-failed", status: "failed", entryTime: "2026-08-15T11:00:00Z" },
  { id: "REQ-queued", status: "pending", entryTime: "" },
  { id: "REQ-recent", status: "completed", entryTime: "2026-08-15T06:00:00Z" },
  { id: "REQ-old", status: "completed", entryTime: "2026-08-13T06:00:00Z" }
] };
var requestsById = {
  "REQ-recent": { status: "completed" },
  "REQ-old": { status: "completed" },
  "REQ-open": { status: "pending" }
};
` + strings.Join(functionBlocks, "\n") + `
var recentIds = recentlyDoneIds(24);
var recentlyDoneIdSet = {};
recentIds.forEach(function (requestId) { recentlyDoneIdSet[requestId] = true; });
process.stdout.write(JSON.stringify({
  recentIds: recentIds,
  recentActive: userRequestHasOpenOrRecentWork({ requestIds: ["REQ-recent"] }, recentlyDoneIdSet),
  oldActive: userRequestHasOpenOrRecentWork({ requestIds: ["REQ-old"] }, recentlyDoneIdSet),
  openActive: userRequestHasOpenOrRecentWork({ requestIds: ["REQ-open"] }, recentlyDoneIdSet)
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR recent-work predicate", javascriptProbe)
	var result struct {
		RecentIds    []string `json:"recentIds"`
		RecentActive bool     `json:"recentActive"`
		OldActive    bool     `json:"oldActive"`
		OpenActive   bool     `json:"openActive"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode by-UR recent-work result: %v (output %q)", decodeError, probeOutput)
	}
	if len(result.RecentIds) != 1 || result.RecentIds[0] != "REQ-recent" {
		t.Fatalf("recentlyDoneIds(24) = %#v, want only REQ-recent — a claimed, failed, or queued calendar "+
			"entry inside the window must never reach the Recently done column", result.RecentIds)
	}
	if !result.RecentActive || result.OldActive || !result.OpenActive {
		t.Fatalf("Active predicate result = recent:%v old:%v open:%v, want true, false, true",
			result.RecentActive, result.OldActive, result.OpenActive)
	}
}

func TestJavaScriptBehaviorByUserRequestLensEmptyState(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function recentWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestLensEmptyText("),
	}
	javascriptProbe := strings.Join(functionBlocks, "\n") + `
const results = [
  recentWindowPhrase(1),
  recentWindowPhrase(168),
  userRequestLensEmptyText(true, 4, 2, "the last 24 hours"),
  userRequestLensEmptyText(true, 4, 0, "the last 24 hours"),
  userRequestLensEmptyText(false, 4, 0, "the last 24 hours"),
  userRequestLensEmptyText(false, 0, 0, "the last 24 hours")
];
process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR empty-state decision", javascriptProbe)
	var results []string
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode assembled-client empty-state results: %v (output %q)", decodeError, probeOutput)
	}
	if len(results) != 6 {
		t.Fatalf("empty-state result count = %d, want 6", len(results))
	}
	if results[0] != "the last 1 hour" || results[1] != "the last 7 days" {
		t.Fatalf("recent-window phrases = %q, %q, want singular hour and seven-day copy", results[0], results[1])
	}
	if !strings.Contains(results[2], "switch URs to All") || !strings.Contains(results[2], "2 resolved matches") {
		t.Fatalf("scope-hidden search result = %q, want an All-scope escape with the match count", results[2])
	}
	if results[3] != "No user requests match the current filters." {
		t.Fatalf("genuine filter miss = %q, want the generic no-match message", results[3])
	}
	if !strings.Contains(results[4], "widen the RECENTLY DONE window") || !strings.Contains(results[4], "switch URs to All") {
		t.Fatalf("scope-only empty state = %q, want both scope escapes", results[4])
	}
	if results[5] != "No user requests in this tree yet." {
		t.Fatalf("empty tree state = %q, want the empty-tree message", results[5])
	}
}

func TestJavaScriptBehaviorTestingStatusUpdateInvalidatesUserRequestLens(t *testing.T) {
	indexHtml := generateLiveSite(t)
	postTestingSource := sliceBalancedBlockAfter(t, indexHtml, "function postTestingStatus(")
	updateCallback := sliceBalancedBlockAfter(t, postTestingSource, ".then(function (payload) {")
	const wiringToken = "applyConfirmedTestingTransition(requestId, testingState, feedbackText, payload)"
	if !strings.Contains(updateCallback, wiringToken) {
		t.Fatalf("testing-status success callback is not wired to its confirmed transition: %q missing", wiringToken)
	}

	transitionFunction := sliceBalancedBlockAfter(t, indexHtml, "function applyConfirmedTestingTransition(")
	javascriptProbe := `
var requestsById = {
  "REQ-1": {
    testingStatus: "",
    testedBy: "",
    testingUpdatedAt: "",
    testingFeedback: "",
    testingStatusUnrecognized: true,
    originalTestingStatus: "invalid"
  }
};
var feedbackFormRequestId = "REQ-1";
var feedbackDraftText = "draft";
var renderedOnce = { userRequestLens: true };
var viewState = { view: "board", lens: "user-request" };
var testingRenderCount = 0;
var columnRenderCount = 0;
var lensRenderCount = 0;
function renderTestingView() { testingRenderCount += 1; }
function renderColumns() { columnRenderCount += 1; }
function renderUserRequestLens() { lensRenderCount += 1; }
` + transitionFunction + `
applyConfirmedTestingTransition("REQ-1", "returned", "needs revision", {
  testingStatus: "returned",
  testedBy: "Alex",
  testingUpdatedAt: "2026-08-15T12:00:00Z"
});
var visibleTransition = {
  request: Object.assign({}, requestsById["REQ-1"]),
  feedbackFormRequestId: feedbackFormRequestId,
  feedbackDraftText: feedbackDraftText,
  testingRenderCount: testingRenderCount,
  columnRenderCount: columnRenderCount,
  lensRenderCount: lensRenderCount,
  lensFresh: renderedOnce.userRequestLens
};
viewState.lens = "columns";
renderedOnce.userRequestLens = true;
applyConfirmedTestingTransition("REQ-1", "tested", "", {
  testingStatus: "tested",
  testedBy: "Alex",
  testingUpdatedAt: "2026-08-15T12:05:00Z"
});
process.stdout.write(JSON.stringify({
  visibleTransition: visibleTransition,
  hiddenLensFresh: renderedOnce.userRequestLens,
  hiddenLensRenderCount: lensRenderCount,
  hiddenTestingRenderCount: testingRenderCount,
  hiddenColumnRenderCount: columnRenderCount
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "confirmed testing transition", javascriptProbe)
	var result struct {
		VisibleTransition struct {
			Request struct {
				TestingStatus             string `json:"testingStatus"`
				TestedBy                  string `json:"testedBy"`
				TestingUpdatedAt          string `json:"testingUpdatedAt"`
				TestingFeedback           string `json:"testingFeedback"`
				TestingStatusUnrecognized bool   `json:"testingStatusUnrecognized"`
				OriginalTestingStatus     string `json:"originalTestingStatus"`
			} `json:"request"`
			FeedbackFormRequestId *string `json:"feedbackFormRequestId"`
			FeedbackDraftText     string  `json:"feedbackDraftText"`
			TestingRenderCount    int     `json:"testingRenderCount"`
			ColumnRenderCount     int     `json:"columnRenderCount"`
			LensRenderCount       int     `json:"lensRenderCount"`
			LensFresh             bool    `json:"lensFresh"`
		} `json:"visibleTransition"`
		HiddenLensFresh          bool `json:"hiddenLensFresh"`
		HiddenLensRenderCount    int  `json:"hiddenLensRenderCount"`
		HiddenTestingRenderCount int  `json:"hiddenTestingRenderCount"`
		HiddenColumnRenderCount  int  `json:"hiddenColumnRenderCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode confirmed testing transition: %v (output %q)", decodeError, probeOutput)
	}
	visibleTransition := result.VisibleTransition
	request := visibleTransition.Request
	if request.TestingStatus != "returned" || request.TestedBy != "Alex" || request.TestingUpdatedAt != "2026-08-15T12:00:00Z" || request.TestingFeedback != "needs revision" || request.TestingStatusUnrecognized || request.OriginalTestingStatus != "returned" {
		t.Fatalf("confirmed testing request state = %#v, want server-confirmed returned state", request)
	}
	if visibleTransition.FeedbackFormRequestId != nil || visibleTransition.FeedbackDraftText != "" {
		t.Fatalf("confirmed testing feedback form state = id:%v draft:%q, want cleared", visibleTransition.FeedbackFormRequestId, visibleTransition.FeedbackDraftText)
	}
	if visibleTransition.TestingRenderCount != 1 || visibleTransition.ColumnRenderCount != 1 || visibleTransition.LensRenderCount != 1 || !visibleTransition.LensFresh {
		t.Fatalf("confirmed testing render state = testing:%d columns:%d lens:%d fresh:%v, want one refresh for each visible surface",
			visibleTransition.TestingRenderCount, visibleTransition.ColumnRenderCount, visibleTransition.LensRenderCount, visibleTransition.LensFresh)
	}
	if result.HiddenLensFresh || result.HiddenLensRenderCount != 1 || result.HiddenTestingRenderCount != 2 || result.HiddenColumnRenderCount != 2 {
		t.Fatalf("hidden by-UR render state = testing:%d columns:%d lens:%d fresh:%v, want lens uncalled and marked stale",
			result.HiddenTestingRenderCount, result.HiddenColumnRenderCount, result.HiddenLensRenderCount, result.HiddenLensFresh)
	}
}

func TestJavaScriptBehaviorDurationsWindowsProjectOneSharedRealTimeDomain(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	completionInstants := []time.Time{
		time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 12, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC),
	}
	fixtureTickets := make([]*RequestTicket, 0, len(completionInstants))
	for completionIndex, completedAt := range completionInstants {
		fixtureTickets = append(fixtureTickets, durationTicket(
			fmt.Sprintf("REQ-%03d", completionIndex+1),
			"B",
			completedAt.Add(-10*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	fixtureBoard := &Board{
		GeneratedAt: time.Date(2030, 1, 1, 12, 0, 0, 0, time.UTC),
		AllRequests: fixtureTickets,
	}
	generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	boardDataJSON, encodeError := json.Marshal(generatedData)
	if encodeError != nil {
		t.Fatalf("encode board payload: %v", encodeError)
	}

	probeDriver := `
function renderedWindow(windowName) {
  durationsStubHosts["durations-chart"] = makeStubNode("div");
  durationsStubHosts["durations-summary"] = makeStubNode("p");
  durationsStubHosts["durations-readout"] = makeStubNode("p");
  durationsStubHosts["durations-table-body"] = makeStubNode("tbody");
  setDurationsWindow(windowName);
  renderDurationsView();
  var markCentres = [], panelBBarCentres = [], panelCBars = 0;
  function walkDrawnNodes(parentNode) {
    (parentNode.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      var nodeClass = String(attributes["class"] || "");
      if (childNode.stubName === "circle" && nodeClass.indexOf("durations-mark") !== -1) {
        markCentres.push(Number(attributes.cx));
      }
      if (childNode.stubName === "rect" && nodeClass === "durations-bar") {
        panelBBarCentres.push(Number(attributes.x) + Number(attributes.width) / 2);
      }
      if (childNode.stubName === "rect" && nodeClass.indexOf("durations-bar-count") !== -1) {
        panelCBars += 1;
      }
      walkDrawnNodes(childNode);
    });
  }
  var svg = durationsStubHosts["durations-chart"].children[0];
  walkDrawnNodes(svg);
  return {
    summary: durationsStubHosts["durations-summary"].textContent,
    ariaLabel: svg.attributes["aria-label"],
    markCentres: markCentres,
    panelBBarCentres: panelBBarCentres,
    panelCBars: panelCBars,
    tableRows: durationsStubHosts["durations-table-body"].children.length
  };
}
process.stdout.write(JSON.stringify({
  last30: renderedWindow("30"),
  last90: renderedWindow("90"),
  all: renderedWindow("all")
}));
`
	javascriptProbe := durationsRenderDomStubPreamble +
		"var boardData = " + string(boardDataJSON) + ";\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "Durations projected windows", javascriptProbe)

	type renderedWindow struct {
		Summary          string    `json:"summary"`
		AriaLabel        string    `json:"ariaLabel"`
		MarkCentres      []float64 `json:"markCentres"`
		PanelBBarCentres []float64 `json:"panelBBarCentres"`
		PanelCBars       int       `json:"panelCBars"`
		TableRows        int       `json:"tableRows"`
	}
	var result struct {
		Last30 renderedWindow `json:"last30"`
		Last90 renderedWindow `json:"last90"`
		All    renderedWindow `json:"all"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode projected Durations windows: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	for _, windowCase := range []struct {
		name          string
		window        renderedWindow
		wantCount     int
		wantStartDate string
		wantEndDate   string
	}{
		{name: "Last 30 days", window: result.Last30, wantCount: 5, wantStartDate: "26 Jul", wantEndDate: "25 Aug"},
		{name: "Last 90 days", window: result.Last90, wantCount: 7, wantStartDate: "27 May", wantEndDate: "25 Aug"},
		{name: "All history", window: result.All, wantCount: 8, wantStartDate: "1 Apr", wantEndDate: "25 Aug"},
	} {
		if len(windowCase.window.MarkCentres) != windowCase.wantCount ||
			windowCase.window.TableRows != windowCase.wantCount ||
			len(windowCase.window.PanelBBarCentres) != windowCase.wantCount ||
			windowCase.window.PanelCBars != windowCase.wantCount {
			t.Errorf("%s counts = marks %d, table %d, Panel B %d, Panel C %d; want %d projected samples/days on every surface",
				windowCase.name, len(windowCase.window.MarkCentres), windowCase.window.TableRows,
				len(windowCase.window.PanelBBarCentres), windowCase.window.PanelCBars, windowCase.wantCount)
		}
		for _, surfaceCopy := range []string{windowCase.window.Summary, windowCase.window.AriaLabel} {
			for _, requiredText := range []string{windowCase.name, windowCase.wantStartDate, windowCase.wantEndDate, "end exclusive", fmt.Sprintf("%d archived REQ", windowCase.wantCount)} {
				if !strings.Contains(surfaceCopy, requiredText) {
					t.Errorf("%s accessibility copy %q is missing %q", windowCase.name, surfaceCopy, requiredText)
				}
			}
		}
	}

	marginLeft := durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")
	plotWidth := durationsRendererConstant(t, "DURATIONS_VIEW_WIDTH") - marginLeft - durationsRendererConstant(t, "DURATIONS_MARGIN_RIGHT")
	firstDayRight := marginLeft + plotWidth/30
	if result.Last30.MarkCentres[0] < marginLeft || result.Last30.MarkCentres[0] > firstDayRight {
		t.Errorf("Last 30 days left-boundary sample x=%.2f, want inside first day slot [%.2f, %.2f]", result.Last30.MarkCentres[0], marginLeft, firstDayRight)
	}
	firstDayGap := result.Last30.PanelBBarCentres[2] - result.Last30.PanelBBarCentres[1]
	secondDayGap := result.Last30.PanelBBarCentres[3] - result.Last30.PanelBBarCentres[2]
	if math.Abs(firstDayGap-secondDayGap) > 0.11 || firstDayGap <= 0 {
		t.Errorf("equal UTC-day gaps draw %.2f and %.2f units apart, want equal positive affine spacing", firstDayGap, secondDayGap)
	}
}

func TestJavaScriptBehaviorDurationsUserRequestLaneNamesEveryUserRequest(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	const fixtureUserRequestCount = 40
	const fixtureNoUserRequestCount = 5
	fixtureTickets := durationsUserRequestLaneFixtureTickets(fixtureUserRequestCount, fixtureNoUserRequestCount)
	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC),
		AllRequests: fixtureTickets,
	}
	generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	// The WHOLE payload, not just the durations slice: the join reads
	// boardData.requests, so a probe handed only the samples would take the
	// unknown-UR path for every row and never exercise the join at all.
	boardDataJson, encodeError := json.Marshal(generatedData)
	if encodeError != nil {
		t.Fatalf("encode board payload: %v", encodeError)
	}

	probeDriver := `
renderDurationsView();
var drawnBrackets = [], drawnTexts = [];
function walkDrawnNodes(parentNode) {
  (parentNode.children || []).forEach(function (childNode) {
    var attributes = childNode.attributes || {};
    var nodeClass = String(attributes["class"] || "");
    if (childNode.stubName === "rect" && nodeClass.indexOf("durations-ur-bracket") !== -1) {
      drawnBrackets.push({ class: nodeClass });
    }
    if (childNode.stubName === "text") {
      drawnTexts.push((childNode.children[0] || {}).textContent || "");
    }
    walkDrawnNodes(childNode);
  });
}
walkDrawnNodes(durationsStubHosts["durations-chart"]);
var userRequestCells = durationsStubHosts["durations-table-body"].children.map(function (tableRow) {
  return tableRow.children[1].textContent;
});
process.stdout.write(JSON.stringify({
  bracketClasses: drawnBrackets.map(function (bracket) { return bracket.class; }),
  drawnTexts: drawnTexts,
  unknownName: DURATIONS_UNKNOWN_USER_REQUEST_NAME,
  remainderSentenceForOne: composeDurationsUserRequestRemainderText(1),
  userRequestCells: userRequestCells
}));
`

	javascriptProbe := durationsRenderDomStubPreamble +
		"var boardData = " + string(boardDataJson) + ";\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "durations UR lane", javascriptProbe)

	var drawn struct {
		BracketClasses          []string `json:"bracketClasses"`
		DrawnTexts              []string `json:"drawnTexts"`
		UnknownName             string   `json:"unknownName"`
		RemainderSentenceForOne string   `json:"remainderSentenceForOne"`
		UserRequestCells        []string `json:"userRequestCells"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
		t.Fatalf("decode drawn UR lane: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	// ---- the fixture does what it claims, checked before anything is read from it.
	userRequestBrackets, unknownBrackets := 0, 0
	for _, bracketClass := range drawn.BracketClasses {
		if strings.Contains(bracketClass, "durations-ur-bracket-unknown") {
			unknownBrackets++
			continue
		}
		userRequestBrackets++
	}
	if userRequestBrackets == 0 {
		t.Fatal("the lane drew no UR brackets at all, so nothing below is a test of the lane")
	}
	if userRequestBrackets >= fixtureUserRequestCount {
		t.Fatalf("all %d fixture URs found a row, so the remainder path was never taken — this fixture no "+
			"longer overflows the lane and the rule it exists to pin is untested", fixtureUserRequestCount)
	}

	// ---- rule one: drawn brackets plus the stated remainder account for every UR.
	statedRemainder := -1
	remainderPattern := regexp.MustCompile(`^\+([0-9]+) URs? with no free row$`)
	for _, drawnText := range drawn.DrawnTexts {
		if remainderPattern.MatchString(drawnText) {
			if _, scanError := fmt.Sscanf(drawnText, "+%d", &statedRemainder); scanError != nil {
				t.Fatalf("read the remainder count out of %q: %v", drawnText, scanError)
			}
		}
	}
	if statedRemainder < 0 {
		t.Fatalf("%d of %d fixture URs found no row and the lane said nothing about it — a reader takes the "+
			"brackets they can see for all of them",
			fixtureUserRequestCount-userRequestBrackets, fixtureUserRequestCount)
	}
	if userRequestBrackets+statedRemainder != fixtureUserRequestCount {
		t.Errorf("the lane drew %d brackets and stated %d more, accounting for %d of the fixture's %d URs — "+
			"every UR is either on a row or in the remainder, and the two must add up at any row count",
			userRequestBrackets, statedRemainder, userRequestBrackets+statedRemainder, fixtureUserRequestCount)
	}

	// ---- rule two: the samples with no UR are named, on every surface, and named apart.
	if unknownBrackets != 1 {
		t.Errorf("the lane drew %d unknown-UR brackets for %d samples carrying no user_request, want exactly 1 — "+
			"the bucket holds one reserved row, it does not compete for one",
			unknownBrackets, fixtureNoUserRequestCount)
	}
	if drawn.UnknownName == "" {
		t.Fatal("the unknown-UR name is empty, so every surface below states nothing")
	}
	if strings.Contains(drawn.RemainderSentenceForOne, drawn.UnknownName) {
		t.Errorf("the remainder sentence %q contains the unknown-UR name %q — a UR that found no row and a "+
			"sample that has no UR at all are different facts and must not read as one",
			drawn.RemainderSentenceForOne, drawn.UnknownName)
	}
	if len(drawn.UserRequestCells) != len(fixtureTickets) {
		t.Fatalf("the table has %d UR cells for %d samples", len(drawn.UserRequestCells), len(fixtureTickets))
	}
	namedUnknownCells := 0
	for cellIndex, cellText := range drawn.UserRequestCells {
		if cellText == "" {
			t.Fatalf("table row %d has a blank UR cell — a sample with no UR must SAY so, and a blank cell "+
				"reads as a rendering fault rather than as a fact about the REQ", cellIndex)
		}
		if cellText == drawn.UnknownName {
			namedUnknownCells++
		}
	}
	if namedUnknownCells != fixtureNoUserRequestCount {
		t.Errorf("%d table cells carry the unknown-UR name for %d samples with no user_request — a sample "+
			"with no UR must never be given one, and one with a UR must never lose it",
			namedUnknownCells, fixtureNoUserRequestCount)
	}
}

func TestJavaScriptBehaviorTimelineVirtualizesRowsAtScale(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_ROW_HEIGHT", "TIMELINE_GROUP_HEADER_HEIGHT", "TIMELINE_OVERSCAN_ROWS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFlattenGroups(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineVisibleDisplayRange(") + `
var groups = [];
var requestIndex = 0;
for (var groupIndex = 0; groupIndex < 80; groupIndex++) {
  var members = [];
  for (var memberIndex = 0; memberIndex < 7; memberIndex++) {
    members.push({ row: { id: "REQ-" + requestIndex }, rowIndex: requestIndex });
    requestIndex++;
  }
  groups.push({ label: "UR-" + groupIndex, members: members });
}
var layout = timelineFlattenGroups(groups);
var displayCount = layout.items.length;
var viewportHeight = 600;
var atTop = timelineVisibleDisplayRange(layout.items, 0, viewportHeight);
var midway = timelineVisibleDisplayRange(layout.items, layout.height / 2, viewportHeight);
var atBottom = timelineVisibleDisplayRange(layout.items, layout.height, viewportHeight);
process.stdout.write(JSON.stringify({
  atTopCount: atTop.lastDisplay - atTop.firstDisplay,
  midwayCount: midway.lastDisplay - midway.firstDisplay,
  midwayCoversScrollPosition:
    layout.items[midway.firstDisplay].topPx <= layout.height / 2 &&
    layout.items[midway.lastDisplay - 1].topPx + layout.items[midway.lastDisplay - 1].height > layout.height / 2,
  atBottomLastRow: atBottom.lastDisplay,
  rowCount: displayCount
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline virtualization", javascriptProbe)
	var sliceResult struct {
		AtTopCount                 int  `json:"atTopCount"`
		MidwayCount                int  `json:"midwayCount"`
		MidwayCoversScrollPosition bool `json:"midwayCoversScrollPosition"`
		AtBottomLastRow            int  `json:"atBottomLastRow"`
		RowCount                   int  `json:"rowCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &sliceResult); decodeError != nil {
		t.Fatalf("decode timeline virtualization behavior: %v (output %q)", decodeError, probeOutput)
	}
	// A 600px viewport holds well under a quarter of the flattened headers and
	// members; the slice is bounded by the VIEWPORT and never by archive size.
	if sliceResult.AtTopCount >= sliceResult.RowCount/4 {
		t.Fatalf("a 600px viewport rendered %d of %d group/member items; the slice must be viewport-bounded",
			sliceResult.AtTopCount, sliceResult.RowCount)
	}
	if sliceResult.MidwayCount >= sliceResult.RowCount/4 {
		t.Fatalf("a midway 600px viewport rendered %d of %d group/member items; mixed fixed heights may vary the count, but the slice must stay viewport-bounded",
			sliceResult.MidwayCount, sliceResult.RowCount)
	}
	if !sliceResult.MidwayCoversScrollPosition {
		t.Fatal("the midway slice does not contain the row at the scroll position")
	}
	if sliceResult.AtBottomLastRow != sliceResult.RowCount {
		t.Fatalf("scrolled past the end the slice reached display item %d, want it clamped to %d",
			sliceResult.AtBottomLastRow, sliceResult.RowCount)
	}
}

func TestJavaScriptBehaviorTimelineTypedDatesMoveTheWindow(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineDateFieldToEpoch(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineStartEpochToDateField(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineEndEpochToDateField(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineUtcDayStart(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTypedWindow(") + `
var boundStart = Date.UTC(2026, 3, 7);
var boundEnd = Date.UTC(2026, 8, 2);
var windowStart = Date.UTC(2026, 5, 1);
var windowEnd = Date.UTC(2026, 5, 8);

function typed(startText, endText) {
  return timelineTypedWindow(startText, endText, windowStart, windowEnd, boundStart, boundEnd);
}
function iso(epochMs) { return new Date(epochMs).toISOString(); }

var bothFields = typed("2026-06-01", "2026-06-15");
var startOnly = typed("2026-06-03", "");
var endOnly = typed("", "2026-06-20");
var sameDay = typed("2026-07-04", "2026-07-04");
var reversed = typed("2026-07-10", "2026-07-01");
var beforeRange = typed("2020-01-01", "2020-01-31");
// A date PAST THE END of the range, typed into From with the end field left
// alone. It has to land on the last day the board HAS, not collapse both
// endpoints onto the bound and leave an empty zoom-floor sliver behind the frame.
var pastRangeStartOnly = typed("2026-09-30", "");
// And the mirror, typed into the end field.
var pastRangeEndOnly = typed("", "2026-09-30");
// A start typed while the end field still holds the board's last day. The
// implied span overruns the ceiling, and a span-preserving settle would pin the
// end to the bound and drag this start backwards to keep the width.
var startAgainstCeiling = timelineTypedWindow(
  "2026-08-01", timelineEndEpochToDateField(boundEnd), windowStart, windowEnd, boundStart, boundEnd);
var neither = typed("", "");
var rubbish = typed("not-a-date", "2026-13-45");
var rolled = timelineDateFieldToEpoch("2026-02-31");

process.stdout.write(JSON.stringify({
  bothStartIso: iso(bothFields.windowStartMs),
  bothEndIso: iso(bothFields.windowEndMs),
  startOnlyStartIso: iso(startOnly.windowStartMs),
  startOnlyKeptEnd: startOnly.windowEndMs === windowEnd,
  endOnlyKeptStart: endOnly.windowStartMs === windowStart,
  endOnlyEndIso: iso(endOnly.windowEndMs),
  sameDaySpanMs: sameDay.windowEndMs - sameDay.windowStartMs,
  reversedOrdered: reversed.windowStartMs < reversed.windowEndMs,
  reversedStartIso: iso(reversed.windowStartMs),
  beforeRangeClampedToBound: beforeRange.windowStartMs >= boundStart,
  beforeRangeStartIso: iso(beforeRange.windowStartMs),
  startAgainstCeilingIso: iso(startAgainstCeiling.windowStartMs),
  pastRangeStartIso: iso(pastRangeStartOnly.windowStartMs),
  pastRangeEndIso: iso(pastRangeStartOnly.windowEndMs),
  pastRangeSpanMs: pastRangeStartOnly.windowEndMs - pastRangeStartOnly.windowStartMs,
  pastRangeEndOnlyEndIso: iso(pastRangeEndOnly.windowEndMs),
  lastDayStartMs: timelineUtcDayStart(boundEnd - 1),
  pastRangeStartMs: pastRangeStartOnly.windowStartMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  // THE ROUND TRIP. Render a real whole-day window into the two fields, parse them
  // straight back, and the same window has to come out — otherwise editing one
  // field re-applies a mangled version of the other. One day and seven days,
  // written as the literal UTC pairs they are now that no calendar helper is left
  // to build them.
  windowRoundTrips: (function () {
    return [
      { name: "one day", startMs: Date.UTC(2026, 6, 15), endMs: Date.UTC(2026, 6, 16) },
      { name: "seven days", startMs: Date.UTC(2026, 6, 13), endMs: Date.UTC(2026, 6, 20) }
    ].map(function (wholeDayWindow) {
      var reparsed = timelineTypedWindow(
        timelineStartEpochToDateField(wholeDayWindow.startMs),
        timelineEndEpochToDateField(wholeDayWindow.endMs),
        0, 0, boundStart, boundEnd);
      return {
        name: wholeDayWindow.name,
        fields: timelineStartEpochToDateField(wholeDayWindow.startMs) + ".." +
          timelineEndEpochToDateField(wholeDayWindow.endMs),
        exact: reparsed.windowStartMs === wholeDayWindow.startMs &&
          reparsed.windowEndMs === wholeDayWindow.endMs
      };
    });
  })(),
  neitherIsNull: neither === null,
  rubbishIsNull: rubbish === null,
  rolledIsNaN: isNaN(rolled),
  roundTrip: timelineStartEpochToDateField(Date.UTC(2026, 5, 9, 13, 45))
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline typed dates", javascriptProbe)
	var typedResult struct {
		BothStartIso              string  `json:"bothStartIso"`
		BothEndIso                string  `json:"bothEndIso"`
		StartOnlyStartIso         string  `json:"startOnlyStartIso"`
		StartOnlyKeptEnd          bool    `json:"startOnlyKeptEnd"`
		EndOnlyKeptStart          bool    `json:"endOnlyKeptStart"`
		EndOnlyEndIso             string  `json:"endOnlyEndIso"`
		SameDaySpanMs             float64 `json:"sameDaySpanMs"`
		ReversedOrdered           bool    `json:"reversedOrdered"`
		ReversedStartIso          string  `json:"reversedStartIso"`
		BeforeRangeClampedToBound bool    `json:"beforeRangeClampedToBound"`
		BeforeRangeStartIso       string  `json:"beforeRangeStartIso"`
		StartAgainstCeilingIso    string  `json:"startAgainstCeilingIso"`
		PastRangeStartIso         string  `json:"pastRangeStartIso"`
		PastRangeEndIso           string  `json:"pastRangeEndIso"`
		PastRangeSpanMs           float64 `json:"pastRangeSpanMs"`
		PastRangeEndOnlyEndIso    string  `json:"pastRangeEndOnlyEndIso"`
		LastDayStartMs            float64 `json:"lastDayStartMs"`
		PastRangeStartMs          float64 `json:"pastRangeStartMs"`
		MinSpanMs                 float64 `json:"minSpanMs"`
		WindowRoundTrips          []struct {
			Name   string `json:"name"`
			Fields string `json:"fields"`
			Exact  bool   `json:"exact"`
		} `json:"windowRoundTrips"`
		NeitherIsNull bool   `json:"neitherIsNull"`
		RubbishIsNull bool   `json:"rubbishIsNull"`
		RolledIsNaN   bool   `json:"rolledIsNaN"`
		RoundTrip     string `json:"roundTrip"`
	}
	if decodeError := json.Unmarshal(probeOutput, &typedResult); decodeError != nil {
		t.Fatalf("decode timeline typed dates behavior: %v (output %q)", decodeError, probeOutput)
	}

	if typedResult.BothStartIso != "2026-06-01T00:00:00.000Z" {
		t.Fatalf("typed start = %s, want the UTC midnight opening that day", typedResult.BothStartIso)
	}
	// The end field names a day to INCLUDE and the window's end is EXCLUSIVE, so
	// the day typed resolves to the FOLLOWING midnight. That is not cosmetic: it
	// is what makes a typed pair produce byte-identical windows to the period
	// chips, which is what the round-trip block below turns into a hard rule.
	// (This assertion previously wanted 23:59:59.999 — the inclusive end that
	// made render and parse non-inverses. Changed deliberately, not bent to fit.)
	if typedResult.BothEndIso != "2026-06-16T00:00:00.000Z" {
		t.Fatalf("typed end = %s, want the midnight following the day typed", typedResult.BothEndIso)
	}
	if !typedResult.StartOnlyKeptEnd || typedResult.StartOnlyStartIso != "2026-06-03T00:00:00.000Z" {
		t.Fatalf("typing only a start moved the end too (start %s, kept end %v); each field must "+
			"resolve against the window already on screen",
			typedResult.StartOnlyStartIso, typedResult.StartOnlyKeptEnd)
	}
	if !typedResult.EndOnlyKeptStart || typedResult.EndOnlyEndIso != "2026-06-21T00:00:00.000Z" {
		t.Fatalf("typing only an end moved the start too (end %s, kept start %v)",
			typedResult.EndOnlyEndIso, typedResult.EndOnlyKeptStart)
	}
	// One date in both fields is the commonest thing a reader will do with a date
	// picker, and it must mean that day rather than an empty window.
	if typedResult.SameDaySpanMs != 86400000 {
		t.Fatalf("the same date in both fields spanned %.0f ms, want exactly one day",
			typedResult.SameDaySpanMs)
	}
	if !typedResult.ReversedOrdered || typedResult.ReversedStartIso != "2026-07-10T00:00:00.000Z" {
		t.Fatalf("an end before the start produced %s and ordered=%v; it must clamp forward from "+
			"the start the reader typed, never silently swap the two",
			typedResult.ReversedStartIso, typedResult.ReversedOrdered)
	}
	// The clamp has to be the shared one. A control with its own bounds is how
	// the reader reaches a window no other control can, and then cannot get back.
	if !typedResult.BeforeRangeClampedToBound {
		t.Fatalf("a date before the board's range escaped the bounds (start %s)",
			typedResult.BeforeRangeStartIso)
	}
	// The defect a browser found and the unit fixture had missed: each endpoint is
	// clamped on its own, because a typed date is a position and the shared settle
	// preserves a span. Pinning the end to the ceiling must never move the start
	// the reader typed.
	if typedResult.StartAgainstCeilingIso != "2026-08-01T00:00:00.000Z" {
		t.Fatalf("a start typed against the range ceiling came back as %s, want the date typed; "+
			"the settle preserves a span and would drag the start back to keep the width",
			typedResult.StartAgainstCeilingIso)
	}
	// One assertion covering the whole class: the fields must be a lossless view of
	// a whole-day window. Rendering an exclusive end instant's own date failed it —
	// every window came back a day long, so editing either field silently moved the
	// other one.
	if len(typedResult.WindowRoundTrips) != 2 {
		t.Fatalf("want a round trip for each of the one-day and seven-day windows; got %d",
			len(typedResult.WindowRoundTrips))
	}
	for _, roundTrip := range typedResult.WindowRoundTrips {
		if !roundTrip.Exact {
			t.Errorf("a %s window rendered to fields %s and did not parse back to itself; "+
				"render and parse must be inverses or editing one field mangles the other",
				roundTrip.Name, roundTrip.Fields)
		}
	}
	if !typedResult.NeitherIsNull {
		t.Fatal("two empty fields must return null — a cleared field is not a request to move")
	}
	if !typedResult.RubbishIsNull {
		t.Fatal("unparseable text in both fields must return null rather than moving the window")
	}
	if !typedResult.RolledIsNaN {
		t.Fatal("2026-02-31 must be rejected; Date.UTC rolls it into March and a rolled date is " +
			"not the one that was typed")
	}
	// A DATE PAST THE END OF THE RANGE lands on the last day the board has. Before
	// this, both endpoints collapsed onto the bound and the settle turned that into
	// an empty one-hour window tucked behind the right edge, while the field went on
	// showing the rejected date.
	if typedResult.PastRangeStartMs != typedResult.LastDayStartMs {
		t.Errorf("a From date past the end of the range put the window start at %s, want the last "+
			"day the board has", typedResult.PastRangeStartIso)
	}
	if typedResult.PastRangeSpanMs <= typedResult.MinSpanMs {
		t.Errorf("a From date past the end of the range produced a %.0f ms window (%s → %s); at the "+
			"zoom floor or below it is the empty sliver this clamp exists to prevent",
			typedResult.PastRangeSpanMs, typedResult.PastRangeStartIso, typedResult.PastRangeEndIso)
	}
	if typedResult.PastRangeEndOnlyEndIso != "2026-09-02T00:00:00.000Z" {
		t.Errorf("a `to` date past the end of the range ended the window at %s, want the range's "+
			"own end", typedResult.PastRangeEndOnlyEndIso)
	}

	if typedResult.RoundTrip != "2026-06-09" {
		t.Fatalf("an instant mid-day rendered into the date field as %q, want its UTC date",
			typedResult.RoundTrip)
	}
}

func TestJavaScriptBehaviorTimelinePanThreshold(t *testing.T) {
	indexHtml := generateLiveSite(t)
	// The distance is READ FROM THE RENDERER, so a probe cannot keep passing
	// against a threshold the board stopped using.
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_PAN_THRESHOLD_PX") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePanEngaged(") + `
var pressX = 500;
process.stdout.write(JSON.stringify({
  threshold: TIMELINE_PAN_THRESHOLD_PX,
  atRest: timelinePanEngaged(false, pressX, pressX),
  justUnder: timelinePanEngaged(false, pressX, pressX + TIMELINE_PAN_THRESHOLD_PX - 0.01),
  exactlyAt: timelinePanEngaged(false, pressX, pressX + TIMELINE_PAN_THRESHOLD_PX),
  // Leftward drags are drags. An unsigned comparison would engage in one
  // direction only, and the bug would look like "panning left is sticky".
  justUnderLeftward: timelinePanEngaged(false, pressX, pressX - TIMELINE_PAN_THRESHOLD_PX + 0.01),
  exactlyAtLeftward: timelinePanEngaged(false, pressX, pressX - TIMELINE_PAN_THRESHOLD_PX),
  // Latched: once engaged, back at the press point is still engaged.
  latchedAtRest: timelinePanEngaged(true, pressX, pressX)
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline pan threshold", javascriptProbe)
	var thresholdResult struct {
		Threshold         float64 `json:"threshold"`
		AtRest            bool    `json:"atRest"`
		JustUnder         bool    `json:"justUnder"`
		ExactlyAt         bool    `json:"exactlyAt"`
		JustUnderLeftward bool    `json:"justUnderLeftward"`
		ExactlyAtLeftward bool    `json:"exactlyAtLeftward"`
		LatchedAtRest     bool    `json:"latchedAtRest"`
	}
	if decodeError := json.Unmarshal(probeOutput, &thresholdResult); decodeError != nil {
		t.Fatalf("decode timeline pan threshold behavior: %v (output %q)", decodeError, probeOutput)
	}

	// A range, not a value. Below about 2px a trackpad tremor still trips it and
	// the click is lost again; much above 8px a deliberate short drag feels stuck.
	if thresholdResult.Threshold < 2 || thresholdResult.Threshold > 8 {
		t.Errorf("the pan threshold is %g px; under 2 a hand tremor trips it and the click is "+
			"lost again, over 8 a short deliberate drag feels stuck", thresholdResult.Threshold)
	}
	if thresholdResult.AtRest {
		t.Error("a press that has not moved at all engaged the pan")
	}
	if thresholdResult.JustUnder || thresholdResult.JustUnderLeftward {
		t.Errorf("a press just under the %g px threshold engaged the pan (rightward %v, "+
			"leftward %v)", thresholdResult.Threshold,
			thresholdResult.JustUnder, thresholdResult.JustUnderLeftward)
	}
	if !thresholdResult.ExactlyAt || !thresholdResult.ExactlyAtLeftward {
		t.Errorf("a press exactly at the %g px threshold did not engage the pan (rightward %v, "+
			"leftward %v); the comparison has to be inclusive and unsigned in distance",
			thresholdResult.Threshold,
			thresholdResult.ExactlyAt, thresholdResult.ExactlyAtLeftward)
	}
	if !thresholdResult.LatchedAtRest {
		t.Error("an already-engaged drag disengaged on returning to the press point; " +
			"engagement latches, or a wandering drag flickers in and out of panning")
	}
}

func TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function timelineKeyboardActivationTarget(") + `
function rowEvent(key, detailId) {
  var trigger = detailId === null ? null : {
    getAttribute: function (name) {
      return name === "data-detail-kind" ? "request" : detailId;
    }
  };
  return { key: key, target: { closest: function () { return trigger; } } };
}
process.stdout.write(JSON.stringify({
  enter: timelineKeyboardActivationTarget(rowEvent("Enter", "REQ-401")),
  space: timelineKeyboardActivationTarget(rowEvent(" ", "REQ-402")),
  legacySpace: timelineKeyboardActivationTarget(rowEvent("Spacebar", "REQ-403")),
  tab: timelineKeyboardActivationTarget(rowEvent("Tab", "REQ-404")),
  arrow: timelineKeyboardActivationTarget(rowEvent("ArrowDown", "REQ-405")),
  offRow: timelineKeyboardActivationTarget(rowEvent("Enter", null))
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline keyboard activation", javascriptProbe)
	var activationResult struct {
		Enter       *struct{ DetailKind, DetailId string } `json:"enter"`
		Space       *struct{ DetailKind, DetailId string } `json:"space"`
		LegacySpace *struct{ DetailKind, DetailId string } `json:"legacySpace"`
		Tab         *struct{ DetailKind, DetailId string } `json:"tab"`
		Arrow       *struct{ DetailKind, DetailId string } `json:"arrow"`
		OffRow      *struct{ DetailKind, DetailId string } `json:"offRow"`
	}
	if decodeError := json.Unmarshal(probeOutput, &activationResult); decodeError != nil {
		t.Fatalf("decode timeline keyboard activation: %v (output %q)", decodeError, probeOutput)
	}
	for _, activated := range []struct {
		keyName string
		result  *struct{ DetailKind, DetailId string }
		wantId  string
	}{
		{"Enter", activationResult.Enter, "REQ-401"},
		{"Space", activationResult.Space, "REQ-402"},
		{"Spacebar (legacy)", activationResult.LegacySpace, "REQ-403"},
	} {
		if activated.result == nil {
			t.Fatalf("%s on a focused row activated nothing; the row advertises role=button", activated.keyName)
		}
		if activated.result.DetailId != activated.wantId || activated.result.DetailKind != "request" {
			t.Fatalf("%s activated %+v, want request/%s", activated.keyName, *activated.result, activated.wantId)
		}
	}
	// Navigation keys and keys pressed off a row must not open anything.
	for _, ignored := range []struct {
		keyName string
		result  *struct{ DetailKind, DetailId string }
	}{
		{"Tab", activationResult.Tab},
		{"ArrowDown", activationResult.Arrow},
		{"Enter off a row", activationResult.OffRow},
	} {
		if ignored.result != nil {
			t.Fatalf("%s activated %+v; it must open nothing", ignored.keyName, *ignored.result)
		}
	}
}

func TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePannedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTrailingWindow(") + `
// A board shaped like this repo's own: five months of archive with the now-line
// well inside it, so a 7-day window is nowhere near either bound.
var boundStart = Date.UTC(2026, 3, 7);           // 7 Apr 2026
var boundEnd = Date.UTC(2026, 8, 2);             // 2 Sep 2026
var nowMs = Date.UTC(2026, 7, 18, 10, 30);       // 18 Aug 2026 10:30 UTC

// And a board a few days old — what anyone who has just installed do-work is
// looking at, and the shape that catches a candidate handed over unclamped.
var shortBoundStart = Date.UTC(2026, 7, 16);     // 16 Aug 2026
var shortBoundEnd = Date.UTC(2026, 7, 19);       // 19 Aug 2026, past now: the forecast
var shortNowMs = Date.UTC(2026, 7, 18, 10, 30);

var lastSeven = timelineTrailingWindow("7", nowMs, boundStart, boundEnd);
// Pressing the same chip again must not move the window it just produced.
var lastSevenAgain = timelineTrailingWindow("7", nowMs, boundStart, boundEnd);
var lastNinetyOnShortArchive = timelineTrailingWindow("90", shortNowMs, shortBoundStart, shortBoundEnd);
var allDays = timelineTrailingWindow("all", nowMs, boundStart, boundEnd);

// One screenful either way, from a window clear of both bounds so the pan is
// free to be an inverse pair rather than clamping at an edge.
var readersWindow = { windowStartMs: Date.UTC(2026, 5, 1), windowEndMs: Date.UTC(2026, 5, 8) };
var steppedForward = timelinePannedWindow(
  readersWindow.windowStartMs, readersWindow.windowEndMs, 1, boundStart, boundEnd);
var steppedBack = timelinePannedWindow(
  steppedForward.windowStartMs, steppedForward.windowEndMs, -1, boundStart, boundEnd);

process.stdout.write(JSON.stringify({
  nowMs: nowMs,
  dayMs: TIMELINE_DAY_MS,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  boundStartMs: boundStart,
  boundEndMs: boundEnd,

  lastSevenStartMs: lastSeven.windowStartMs,
  lastSevenEndMs: lastSeven.windowEndMs,
  lastSevenSpanMs: lastSeven.windowEndMs - lastSeven.windowStartMs,
  lastSevenIsIdempotent:
    lastSevenAgain.windowStartMs === lastSeven.windowStartMs &&
    lastSevenAgain.windowEndMs === lastSeven.windowEndMs,

  shortNowMs: shortNowMs,
  shortBoundStartMs: shortBoundStart,
  shortNinetyStartMs: lastNinetyOnShortArchive.windowStartMs,
  shortNinetyEndMs: lastNinetyOnShortArchive.windowEndMs,

  allStartMs: allDays.windowStartMs,
  allEndMs: allDays.windowEndMs,

  readersStartMs: readersWindow.windowStartMs,
  readersSpanMs: readersWindow.windowEndMs - readersWindow.windowStartMs,
  steppedForwardStartMs: steppedForward.windowStartMs,
  steppedForwardSpanMs: steppedForward.windowEndMs - steppedForward.windowStartMs,
  steppedBackStartMs: steppedBack.windowStartMs,
  steppedBackEndMs: steppedBack.windowEndMs
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline trailing windows", javascriptProbe)
	var trailingResult struct {
		NowMs        float64 `json:"nowMs"`
		DayMs        float64 `json:"dayMs"`
		MinSpanMs    float64 `json:"minSpanMs"`
		BoundStartMs float64 `json:"boundStartMs"`
		BoundEndMs   float64 `json:"boundEndMs"`

		LastSevenStartMs      float64 `json:"lastSevenStartMs"`
		LastSevenEndMs        float64 `json:"lastSevenEndMs"`
		LastSevenSpanMs       float64 `json:"lastSevenSpanMs"`
		LastSevenIsIdempotent bool    `json:"lastSevenIsIdempotent"`

		ShortNowMs         float64 `json:"shortNowMs"`
		ShortBoundStartMs  float64 `json:"shortBoundStartMs"`
		ShortNinetyStartMs float64 `json:"shortNinetyStartMs"`
		ShortNinetyEndMs   float64 `json:"shortNinetyEndMs"`

		AllStartMs float64 `json:"allStartMs"`
		AllEndMs   float64 `json:"allEndMs"`

		ReadersStartMs        float64 `json:"readersStartMs"`
		ReadersSpanMs         float64 `json:"readersSpanMs"`
		SteppedForwardStartMs float64 `json:"steppedForwardStartMs"`
		SteppedForwardSpanMs  float64 `json:"steppedForwardSpanMs"`
		SteppedBackStartMs    float64 `json:"steppedBackStartMs"`
		SteppedBackEndMs      float64 `json:"steppedBackEndMs"`
	}
	if decodeError := json.Unmarshal(probeOutput, &trailingResult); decodeError != nil {
		t.Fatalf("decode timeline trailing-window behavior: %v (output %q)", decodeError, probeOutput)
	}

	utcOf := func(epochMs float64) string {
		return time.UnixMilli(int64(epochMs)).UTC().Format(time.RFC3339)
	}

	// (1) A chip ends the window at NOW and gives the span it names. "Seven days
	// long" alone would pass for a window seven days in the wrong place.
	if trailingResult.LastSevenEndMs != trailingResult.NowMs {
		t.Fatalf("the Last 7 days window ends at %s, want the board's now %s",
			utcOf(trailingResult.LastSevenEndMs), utcOf(trailingResult.NowMs))
	}
	if trailingResult.LastSevenSpanMs != 7*trailingResult.DayMs {
		t.Fatalf("the Last 7 days window spans %.0f ms (%s → %s), want exactly seven days",
			trailingResult.LastSevenSpanMs, utcOf(trailingResult.LastSevenStartMs),
			utcOf(trailingResult.LastSevenEndMs))
	}
	if !trailingResult.LastSevenIsIdempotent {
		t.Errorf("pressing the same chip twice moved the window away from %s → %s",
			utcOf(trailingResult.LastSevenStartMs), utcOf(trailingResult.LastSevenEndMs))
	}

	// (2) THE SLIDE TRAP. A window wider than the archive is CUT SHORT at the range
	// start; handing the candidate to the settle unclamped pins the start to the
	// bound and drags the end forward to keep the width, so the window stops ending
	// at now — the one thing every chip on this toolbar promises.
	if trailingResult.ShortNinetyStartMs != trailingResult.ShortBoundStartMs {
		t.Errorf("Last 90 days on a three-day archive starts at %s, want the range start %s",
			utcOf(trailingResult.ShortNinetyStartMs), utcOf(trailingResult.ShortBoundStartMs))
	}
	if trailingResult.ShortNinetyEndMs != trailingResult.ShortNowMs {
		t.Errorf("Last 90 days on a three-day archive ends at %s, want now %s; the settle preserves "+
			"a width and slides, so the candidate has to be clamped into the bounds before it",
			utcOf(trailingResult.ShortNinetyEndMs), utcOf(trailingResult.ShortNowMs))
	}

	// (3) All days is the whole recorded range, so nothing drawn is out of reach of
	// the button that says it shows all of it.
	if trailingResult.AllStartMs != trailingResult.BoundStartMs ||
		trailingResult.AllEndMs != trailingResult.BoundEndMs {
		t.Errorf("All days spans %s → %s, want the recorded range %s → %s",
			utcOf(trailingResult.AllStartMs), utcOf(trailingResult.AllEndMs),
			utcOf(trailingResult.BoundStartMs), utcOf(trailingResult.BoundEndMs))
	}

	// (4) A pan moves one screenful, and forward-then-back is the reader's undo —
	// checked here away from the bounds, where the pan's clamp cannot fire.
	if trailingResult.SteppedForwardStartMs-trailingResult.ReadersStartMs != trailingResult.ReadersSpanMs {
		t.Errorf("one step forward moved the window %.0f ms, want its own %.0f ms span",
			trailingResult.SteppedForwardStartMs-trailingResult.ReadersStartMs, trailingResult.ReadersSpanMs)
	}
	if trailingResult.SteppedForwardSpanMs != trailingResult.ReadersSpanMs {
		t.Errorf("one step forward resized the window from %.0f ms to %.0f ms; a step moves the "+
			"window, it does not resize it", trailingResult.ReadersSpanMs, trailingResult.SteppedForwardSpanMs)
	}
	if trailingResult.SteppedBackStartMs != trailingResult.ReadersStartMs ||
		trailingResult.SteppedBackEndMs-trailingResult.SteppedBackStartMs != trailingResult.ReadersSpanMs {
		t.Errorf("forward then back landed on %s spanning %.0f ms, want the %s it started from "+
			"spanning %.0f ms", utcOf(trailingResult.SteppedBackStartMs),
			trailingResult.SteppedBackEndMs-trailingResult.SteppedBackStartMs,
			utcOf(trailingResult.ReadersStartMs), trailingResult.ReadersSpanMs)
	}
}

func TestJavaScriptBehaviorTimelineWindowStepArrowsAreInversesAtTheBounds(t *testing.T) {
	indexHtml := generateLiveSite(t)

	arrowStepCallSite := sliceBalancedBlockAfter(t, indexHtml, "function steppedWindowFor(")
	arrowStepCallPattern := regexp.MustCompile(`return\s+(timeline[A-Za-z0-9]*)\(`)
	arrowStepCallMatch := arrowStepCallPattern.FindStringSubmatch(arrowStepCallSite)
	if arrowStepCallMatch == nil {
		t.Fatalf("steppedWindowFor does not return a timeline* window function, so this probe has no "+
			"call site to follow:\n%s", arrowStepCallSite)
	}
	arrowStepFunctionName := arrowStepCallMatch[1]

	// timelinePannedWindow is the clamp every window-moving path shares, so it is
	// in the probe whether or not the arrows call it directly.
	shippedFunctions := sliceBalancedBlockAfter(t, indexHtml, "function timelinePannedWindow(")
	if arrowStepFunctionName != "timelinePannedWindow" {
		shippedFunctions += "\n" +
			sliceBalancedBlockAfter(t, indexHtml, "function "+arrowStepFunctionName+"(")
	}

	javascriptProbe := timelineProbePreamble(t, "TIMELINE_DAY_MS") + shippedFunctions + "\n" +
		"var arrowStep = " + arrowStepFunctionName + ";" + `
// A 90-day board and the reader on a 7-day window: wide enough that a screenful
// is nowhere near the whole range, narrow enough that thirteen of them fit.
var boundStartMs = Date.UTC(2026, 3, 7);
var boundEndMs = boundStartMs + 90 * TIMELINE_DAY_MS;
var windowSpanMs = 7 * TIMELINE_DAY_MS;

function hoursOf(spanMs) { return (spanMs / 3600000).toFixed(2) + "h"; }
function noteAtMost(list, entry) { if (list.length < 8) { list.push(entry); } }

// One press, then the opposite press from wherever it landed.
function stepAndStepBack(windowStartMs, stepCount) {
  var windowEndMs = windowStartMs + windowSpanMs;
  var stepped = arrowStep(windowStartMs, windowEndMs, stepCount, boundStartMs, boundEndMs);
  var back = arrowStep(stepped.windowStartMs, stepped.windowEndMs, -stepCount, boundStartMs, boundEndMs);
  return {
    movedMs: stepped.windowStartMs - windowStartMs,
    steppedStartMs: stepped.windowStartMs,
    steppedEndMs: stepped.windowEndMs,
    steppedSpanMs: stepped.windowEndMs - stepped.windowStartMs,
    driftMs: back.windowStartMs - windowStartMs,
    backSpanMs: back.windowEndMs - back.windowStartMs
  };
}

var partialSteps = [];
var refusedWithRoom = [];
var driftedRoundTrips = [];
var resizedWindows = [];
var escapedTheBounds = [];
var wholeScreenfulStepCount = 0;
var refusedStepCount = 0;

function checkOnePress(label, windowStartMs, stepCount) {
  var roomAheadMs = stepCount > 0
    ? boundEndMs - (windowStartMs + windowSpanMs)
    : windowStartMs - boundStartMs;
  var press = stepAndStepBack(windowStartMs, stepCount);
  var wholeScreenfulMs = windowSpanMs * stepCount;
  var moved = press.movedMs !== 0;
  if (moved && press.movedMs !== wholeScreenfulMs) {
    noteAtMost(partialSteps, label + " (step " + stepCount + ", " + hoursOf(roomAheadMs) +
      " of room): moved " + hoursOf(press.movedMs) + " of the window's own " +
      hoursOf(windowSpanMs));
  }
  if (!moved && roomAheadMs >= windowSpanMs) {
    noteAtMost(refusedWithRoom, label + " (step " + stepCount + "): did not move with " +
      hoursOf(roomAheadMs) + " of room ahead for a " + hoursOf(windowSpanMs) + " window");
  }
  if (moved && press.driftMs !== 0) {
    noteAtMost(driftedRoundTrips, label + " (step " + stepCount + "): press then unpress drifted " +
      hoursOf(press.driftMs));
  }
  if (press.steppedSpanMs !== windowSpanMs || press.backSpanMs !== windowSpanMs) {
    noteAtMost(resizedWindows, label + " (step " + stepCount + "): the window became " +
      hoursOf(press.steppedSpanMs) + " and then " + hoursOf(press.backSpanMs) + ", from " +
      hoursOf(windowSpanMs));
  }
  // The cheapest way to make a step its own inverse everywhere is to stop clamping
  // it, which walks the reader off the end of the board. Refusing the step is the
  // remedy; overshooting the bounds is not.
  if (press.steppedStartMs < boundStartMs || press.steppedEndMs > boundEndMs) {
    noteAtMost(escapedTheBounds, label + " (step " + stepCount + "): landed outside the board's " +
      "own range, " + hoursOf(boundStartMs - press.steppedStartMs) + " before its start and " +
      hoursOf(press.steppedEndMs - boundEndMs) + " past its end");
  }
  if (moved) { wholeScreenfulStepCount++; } else { refusedStepCount++; }
  return press;
}

// The five positions the review measured, kept by name so a failure says which
// one, plus their mirrors at the left bound.
var namedPositions = [
  { label: "mid-range", windowStartMs: boundStartMs + 40 * TIMELINE_DAY_MS },
  { label: "one span from the right bound", windowStartMs: boundEndMs - 2 * windowSpanMs },
  { label: "a partial screenful from the right bound", windowStartMs: boundEndMs - windowSpanMs - 2 * TIMELINE_DAY_MS },
  { label: "flush against the right bound", windowStartMs: boundEndMs - windowSpanMs },
  { label: "one span from the left bound", windowStartMs: boundStartMs + windowSpanMs },
  { label: "a partial screenful from the left bound", windowStartMs: boundStartMs + 2 * TIMELINE_DAY_MS },
  { label: "flush against the left bound", windowStartMs: boundStartMs }
];
var namedTable = [];
namedPositions.forEach(function (position) {
  [1, -1].forEach(function (stepCount) {
    var press = checkOnePress(position.label, position.windowStartMs, stepCount);
    namedTable.push(position.label + "  step " + (stepCount > 0 ? "+1" : "-1") +
      "  moved=" + hoursOf(press.movedMs) + "  drift=" + hoursOf(press.driftMs));
  });
});

// And every window position on the board, half a day apart, in both directions.
var sweptPositionCount = 0;
for (var offsetDays = 0; offsetDays <= 2 * (90 - 7); offsetDays++) {
  var sweptStartMs = boundStartMs + offsetDays * TIMELINE_DAY_MS / 2;
  sweptPositionCount++;
  checkOnePress("window opening at day " + (offsetDays / 2), sweptStartMs, 1);
  checkOnePress("window opening at day " + (offsetDays / 2), sweptStartMs, -1);
}

process.stdout.write(JSON.stringify({
  arrowStepFunctionName: ` + "\"" + arrowStepFunctionName + "\"" + `,
  sweptPositionCount: sweptPositionCount,
  wholeScreenfulStepCount: wholeScreenfulStepCount,
  refusedStepCount: refusedStepCount,
  partialSteps: partialSteps,
  refusedWithRoom: refusedWithRoom,
  driftedRoundTrips: driftedRoundTrips,
  resizedWindows: resizedWindows,
  escapedTheBounds: escapedTheBounds,
  namedTable: namedTable
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline window step arrows", javascriptProbe)
	var stepResult struct {
		ArrowStepFunctionName   string   `json:"arrowStepFunctionName"`
		SweptPositionCount      int      `json:"sweptPositionCount"`
		WholeScreenfulStepCount int      `json:"wholeScreenfulStepCount"`
		RefusedStepCount        int      `json:"refusedStepCount"`
		PartialSteps            []string `json:"partialSteps"`
		RefusedWithRoom         []string `json:"refusedWithRoom"`
		DriftedRoundTrips       []string `json:"driftedRoundTrips"`
		ResizedWindows          []string `json:"resizedWindows"`
		EscapedTheBounds        []string `json:"escapedTheBounds"`
		NamedTable              []string `json:"namedTable"`
	}
	if decodeError := json.Unmarshal(probeOutput, &stepResult); decodeError != nil {
		t.Fatalf("decode the window-step sweep: %v (output %q)", decodeError, probeOutput)
	}

	// The sweep has to have swept, and both outcomes have to occur in it: a run
	// where nothing ever moves satisfies "no partial steps" and "no drift" for free.
	const sweptWindowPositionCount = 167
	if stepResult.SweptPositionCount != sweptWindowPositionCount {
		t.Fatalf("the sweep visited %d window positions, want %d", stepResult.SweptPositionCount,
			sweptWindowPositionCount)
	}
	if stepResult.WholeScreenfulStepCount == 0 || stepResult.RefusedStepCount == 0 {
		t.Fatalf("%s moved the window on %d presses and refused %d; a sweep where one of those is "+
			"never reached cannot tell a step apart from a no-op", stepResult.ArrowStepFunctionName,
			stepResult.WholeScreenfulStepCount, stepResult.RefusedStepCount)
	}

	// (1) A step is a WHOLE screenful. A partial one is the defect itself: it is
	// what the back press then fails to undo.
	if len(stepResult.PartialSteps) > 0 {
		t.Errorf("%s moved the window by part of a screenful:\n\t%s", stepResult.ArrowStepFunctionName,
			strings.Join(stepResult.PartialSteps, "\n\t"))
	}

	// (2) THE PROPERTY THE DELETED TEST HELD. Press and unpress is the reader's undo.
	if len(stepResult.DriftedRoundTrips) > 0 {
		t.Errorf("%s is not its own inverse:\n\t%s", stepResult.ArrowStepFunctionName,
			strings.Join(stepResult.DriftedRoundTrips, "\n\t"))
	}

	// (3) The half that keeps (1) and (2) from being satisfied by refusing every
	// press: wherever a screenful of room exists, the arrow uses it.
	if len(stepResult.RefusedWithRoom) > 0 {
		t.Errorf("%s refused a step that had a whole screenful of room:\n\t%s",
			stepResult.ArrowStepFunctionName, strings.Join(stepResult.RefusedWithRoom, "\n\t"))
	}

	// (4) And a step moves the window rather than resizing it, either way.
	if len(stepResult.ResizedWindows) > 0 {
		t.Errorf("%s resized the window:\n\t%s", stepResult.ArrowStepFunctionName,
			strings.Join(stepResult.ResizedWindows, "\n\t"))
	}

	// (5) And it stays on the board. Not clamping at all would make every step an
	// inverse of itself and put the reader past the end of the data to do it.
	if len(stepResult.EscapedTheBounds) > 0 {
		t.Errorf("%s stepped outside the board's range:\n\t%s", stepResult.ArrowStepFunctionName,
			strings.Join(stepResult.EscapedTheBounds, "\n\t"))
	}

	if t.Failed() {
		t.Logf("the arrows are wired to %s; the seven named positions measured:\n\t%s",
			stepResult.ArrowStepFunctionName, strings.Join(stepResult.NamedTable, "\n\t"))
	}
}

func TestJavaScriptBehaviorSpanFormattersCarryRoundedRemainders(t *testing.T) {
	roundingCases := []spanRoundingCase{
		{119.5, "2h 0m", "2h 0m", "the minute remainder rounds to a full hour"},
		{-119.5, "−2h 0m", "−2h 0m", "the same carry on a reversed span"},
		{59.96, "1h 0m", "1h 0m", "the sub-hour branch rounds up to the hour boundary"},
		{119.4, "1h 59m", "1h 59m", "just under the carry still splits normally"},
		{2879, "47h 59m", "2d 0h", "the hour remainder carries into the day field"},
		{1439.5, "24h 0m", "1d 0h", "a day boundary carries the same way"},
		{7.5, "7.5 min", "8 min", "each formatter keeps its own sub-hour precision"},
	}
	probeValues, encodeError := json.Marshal(roundingCases)
	if encodeError != nil {
		t.Fatalf("encode probe values: %v", encodeError)
	}

	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function formatDurationMinutes(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatSpanMinutes(") + `
process.stdout.write(JSON.stringify(` + string(probeValues) + `.map(function (roundingCase) {
  return {
    durationsText: formatDurationMinutes(roundingCase.minutes),
    timelineText: timelineFormatSpanMinutes(roundingCase.minutes)
  };
})));`

	probeOutput := runJavaScriptBehaviorProbe(t, "span formatter rounding", javascriptProbe)
	var drawnTexts []struct {
		DurationsText string `json:"durationsText"`
		TimelineText  string `json:"timelineText"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawnTexts); decodeError != nil {
		t.Fatalf("decode span formatting: %v (output %q)", decodeError, probeOutput)
	}
	if len(drawnTexts) != len(roundingCases) {
		t.Fatalf("probe returned %d results, want %d", len(drawnTexts), len(roundingCases))
	}
	for caseIndex, roundingCase := range roundingCases {
		if drawnTexts[caseIndex].DurationsText != roundingCase.WantDurationsText {
			t.Errorf("%s: formatDurationMinutes(%.2f) drew %q, want %q",
				roundingCase.Requirement, roundingCase.Minutes,
				drawnTexts[caseIndex].DurationsText, roundingCase.WantDurationsText)
		}
		if drawnTexts[caseIndex].TimelineText != roundingCase.WantTimelineText {
			t.Errorf("%s: timelineFormatSpanMinutes(%.2f) drew %q, want %q",
				roundingCase.Requirement, roundingCase.Minutes,
				drawnTexts[caseIndex].TimelineText, roundingCase.WantTimelineText)
		}
	}
}

func TestJavaScriptBehaviorReversedWaitDrawsAsABreak(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	// claimed_at precedes created_at on REQ-901 only. REQ-902 is an ordinary
	// closed wait; REQ-903 is an unclaimed REQ whose wait runs to the now-line.
	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-901","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T11:30:00Z","waitMinutes":-60,"workMinutes":90,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-902","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-903","createdTime":"2026-08-18T11:00:00Z","claimedTime":null,
	     "completedTime":null,"waitMinutes":60,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ]
	}`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		timelineRenderProbeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline reversed wait", javascriptProbe)

	var drawn struct {
		Rows []struct {
			Id    string `json:"id"`
			Rects []struct {
				Class string  `json:"class"`
				Width float64 `json:"width"`
			} `json:"rects"`
		} `json:"rows"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
		t.Fatalf("decode drawn timeline rows: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}
	if len(drawn.Rows) != 3 {
		t.Fatalf("want one drawn group per fixture row, got %d", len(drawn.Rows))
	}

	rowClasses := map[string][]string{}
	rowWidths := map[string]map[string]float64{}
	for _, drawnRow := range drawn.Rows {
		rowWidths[drawnRow.Id] = map[string]float64{}
		for _, rect := range drawnRow.Rects {
			rowClasses[drawnRow.Id] = append(rowClasses[drawnRow.Id], rect.Class)
			rowWidths[drawnRow.Id][rect.Class] = rect.Width
		}
	}

	rowDrewClassContaining := func(rowId string, classFragment string) bool {
		for _, drawnClass := range rowClasses[rowId] {
			if strings.Contains(drawnClass, classFragment) {
				return true
			}
		}
		return false
	}

	// (1) The reversed wait is a break, and there is no wait bar to misread.
	if !rowDrewClassContaining("REQ-901", "timeline-segment-broken") {
		t.Errorf("a wait whose claim precedes its capture must draw the break marker, got %v", rowClasses["REQ-901"])
	}
	if rowDrewClassContaining("REQ-901", "timeline-segment-wait") {
		t.Errorf("a reversed wait must draw NO wait bar — the table prints the negative value beside it, got %v",
			rowClasses["REQ-901"])
	}
	// Its work span is ordinary and must be untouched by the wait's branch.
	if !rowDrewClassContaining("REQ-901", "timeline-segment-work") {
		t.Errorf("a reversed wait must not suppress the row's ordinary work bar, got %v", rowClasses["REQ-901"])
	}

	// (2) An ordinary closed wait still draws its bar, with real width.
	if !rowDrewClassContaining("REQ-902", "timeline-segment-wait") {
		t.Errorf("an ordinary positive wait must still draw its bar, got %v", rowClasses["REQ-902"])
	}
	if rowDrewClassContaining("REQ-902", "timeline-segment-broken") {
		t.Errorf("an ordinary positive wait must not draw a break marker, got %v", rowClasses["REQ-902"])
	}

	// (3) The open wait keeps its is-open bar: it is measured to the now-line and
	// is never reversed, so the new branch must not reach it.
	if !rowDrewClassContaining("REQ-903", "timeline-segment-wait is-open") {
		t.Errorf("an unclaimed REQ must still draw its open wait bar, got %v", rowClasses["REQ-903"])
	}
	if rowDrewClassContaining("REQ-903", "timeline-segment-broken") {
		t.Errorf("an open wait must not draw a break marker, got %v", rowClasses["REQ-903"])
	}

	// The break marker is a fixed-width mark, not a measured span — a break whose
	// width tracked the reversed magnitude would be the same lie in a new shape.
	if brokenWidth := rowWidths["REQ-901"]["timeline-segment-broken"]; brokenWidth != 6 {
		t.Errorf("the break marker must be the same fixed 6-unit mark the work branch draws, got %v", brokenWidth)
	}
}

func TestJavaScriptBehaviorTimelineSummaryCountsRowsDrawnAsBreaks(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T08:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-911","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:30:00Z",
	     "completedTime":"2026-08-18T10:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":true},
	    {"id":"REQ-912","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T11:30:00Z","waitMinutes":-60,"workMinutes":90,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-913","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T11:00:00Z",
	     "completedTime":"2026-08-18T10:00:00Z","waitMinutes":120,"workMinutes":-60,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-914","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T09:00:00Z","waitMinutes":-60,"workMinutes":-60,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-915","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T10:00:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":60,"workMinutes":60,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-916","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:45:00Z",
	     "waitMinutes":45,"workMinutes":0,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":true}
	  ]
	}`

	probeDriver := `
function timelineSummaryWithFilter(visibleIds) {
  [
    "timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
    "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
  ].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
  timelineStubVisibleIds = visibleIds;
  renderTimelineView();
  return timelineStubHosts["timeline-summary"].textContent || "";
}
process.stdout.write(JSON.stringify({
  unfiltered: timelineSummaryWithFilter(null),
  filtered: timelineSummaryWithFilter(["REQ-912", "REQ-915"]),
  anomalyOnly: timelineSummaryWithFilter(["REQ-911"]),
  unresolvedOnly: timelineSummaryWithFilter(["REQ-916"]),
  unresolvedAndAnomalyOnly: timelineSummaryWithFilter(["REQ-911", "REQ-916"]),
  reversedPair: timelineSummaryWithFilter(["REQ-912", "REQ-913"]),
  reversedWaitOnly: timelineSummaryWithFilter(["REQ-912"]),
  reversedWorkOnly: timelineSummaryWithFilter(["REQ-913"]),
  combinedCausesOnly: timelineSummaryWithFilter(["REQ-914"]),
  healthyOnly: timelineSummaryWithFilter(["REQ-915"])
}));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline summary break count", javascriptProbe)

	var summaries struct {
		Unfiltered               string `json:"unfiltered"`
		Filtered                 string `json:"filtered"`
		AnomalyOnly              string `json:"anomalyOnly"`
		UnresolvedOnly           string `json:"unresolvedOnly"`
		UnresolvedAndAnomalyOnly string `json:"unresolvedAndAnomalyOnly"`
		ReversedPair             string `json:"reversedPair"`
		ReversedWaitOnly         string `json:"reversedWaitOnly"`
		ReversedWorkOnly         string `json:"reversedWorkOnly"`
		CombinedCausesOnly       string `json:"combinedCausesOnly"`
		HealthyOnly              string `json:"healthyOnly"`
	}
	if decodeError := json.Unmarshal(probeOutput, &summaries); decodeError != nil {
		t.Fatalf("decode rendered timeline summaries: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	wantBreakClause := func(caseName string, summary string, wantCount int) {
		t.Helper()
		wantClause := fmt.Sprintf(". %d with broken stamps, drawn as breaks.", wantCount)
		if !strings.Contains(summary, wantClause) {
			t.Errorf("%s summary must contain %q, got %q", caseName, wantClause, summary)
		}
	}

	wantNoBreakClause := func(caseName string, summary string) {
		t.Helper()
		if strings.Contains(summary, "with broken stamps") {
			t.Errorf("%s summary must omit the break clause, got %q", caseName, summary)
		}
	}

	wantBreakClause("unfiltered", summaries.Unfiltered, 4)
	wantBreakClause("filtered", summaries.Filtered, 1)
	wantBreakClause("reversed wait and work", summaries.ReversedPair, 2)
	wantBreakClause("reversed wait only", summaries.ReversedWaitOnly, 1)
	wantBreakClause("reversed work only", summaries.ReversedWorkOnly, 1)
	wantBreakClause("combined causes", summaries.CombinedCausesOnly, 1)
	// A REQ that stopped with no resolvable end instant IS drawn as a break, so it
	// is counted.
	wantBreakClause("unresolved only", summaries.UnresolvedOnly, 1)
	// And a row that is merely flagged anomalous, with every span drawn, is NOT —
	// which is the whole of REQ-328's change to this clause. Both of these fail if
	// `row.anomaly ||` comes back: the first would report 2, the second 1.
	wantBreakClause("unresolved beside an anomalous-but-drawn row", summaries.UnresolvedAndAnomalyOnly, 1)
	wantNoBreakClause("anomaly only", summaries.AnomalyOnly)
	wantNoBreakClause("healthy only", summaries.HealthyOnly)
}

func TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan(t *testing.T) {
	boardData := buildImplementationSpanFixturePayload(t)
	requestsById := boardData.Requests
	payloadJson, encodeError := json.Marshal(requestsById)
	if encodeError != nil {
		t.Fatalf("encode fixture payload: %v", encodeError)
	}

	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function truncateBadgeText("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeBadge("),
		sliceBalancedBlockAfter(t, indexHtml, "function futureStampTooltipText("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeImplementationSpanNode("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeRequestCard("),
		sliceBalancedBlockAfter(t, indexHtml, "function formatElapsedDuration("),
	}
	javascriptProbe := `
var filterState = { searchText: "" };
var boardData = { implementationSpanPausedBadgeText: ` + mustMarshalJSONString(t, boardData.ImplementationSpanPausedBadgeText) + ` };
var requestsById = ` + string(payloadJson) + `;
function makeNode(tagName) {
  var node = {
    tagName: tagName,
    className: "",
    textContent: "",
    title: "",
    childNodes: [],
    dataset: {},
    setAttribute: function () {},
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
  node.classList = { add: function (extraClass) { node.className += (node.className ? " " : "") + extraClass; } };
  return node;
}
var document = {
  createElement: function (tagName) { return makeNode(tagName); },
  createTextNode: function (text) { return { nodeType: "text", text: text, className: "", childNodes: [] }; }
};
var futureStampCauseText = "";
// formatElapsedDuration's clock-skew branch must stay UNREACHABLE from a done
// card: the Go verdict is branched on before the formatter runs, so a reversed
// span never reaches it. These two values make that a check rather than a claim.
// The allowance is deliberately hostile — zero, the most permissive setting —
// so any negative span that did reach the formatter would certainly take the
// branch, and the sentinel below would surface in the rendered text.
var futureInstantSkewAllowanceMs = 0;
var clockSkewMarkerText = "SKEW-BRANCH-REACHED";
function formatShortInstantWithRelative(isoText) { return isoText; }
function activeDependentIds() { return []; }
function isTerminalResolvedStatus() { return true; }
function describeRequestStatus(requestId) { return requestId; }
function stateTimerSpecFor() { return null; }
function makeInstantWithStopwatchNode() { return null; }
function makeInstantWithRelativeNode(isoText) { return document.createTextNode(isoText); }
` + strings.Join(functionBlocks, "\n") + `
function nodeText(node) {
  if (node.nodeType === "text") { return node.text; }
  return (node.textContent || "") + node.childNodes.map(nodeText).join("");
}
var renderedCards = {};
Object.keys(requestsById).forEach(function (requestId) {
  var card = makeRequestCard(requestId, { showCompleted: true });
  var doneLines = card.childNodes.filter(function (childNode) { return childNode.className === "req-card-completed"; });
  var spanNodes = doneLines.length === 1
    ? doneLines[0].childNodes.filter(function (childNode) {
        return (childNode.className || "").split(/\s+/).indexOf("elapsed-duration") !== -1;
      })
    : [];
  var markerNodes = spanNodes.length === 1
    ? spanNodes[0].childNodes.filter(function (childNode) {
        return (childNode.className || "").split(/\s+/).indexOf("status-invalid-flag") !== -1;
      })
    : [];
  renderedCards[requestId] = {
    doneLineCount: doneLines.length,
    doneLineText: doneLines.length === 1 ? nodeText(doneLines[0]) : "",
    spanNodeCount: spanNodes.length,
    spanText: spanNodes.length === 1 ? nodeText(spanNodes[0]) : "",
    // A FINISHED span must never tick. refreshRelativeTimeNodes selects on
    // [data-instant-ms], so carrying that key would have the 1s ticker rewrite
    // every done card's span as elapsed-since-epoch. This is the single property
    // that justified not reusing makeElapsedDurationNode, so it is asserted
    // rather than left to a one-off browser observation.
    spanTickerKeys: spanNodes.length === 1
      ? Object.keys(spanNodes[0].dataset || {}).sort()
      : [],
    markerTitle: markerNodes.length === 1 ? markerNodes[0].title : ""
  };
});
process.stdout.write(JSON.stringify(renderedCards));`

	probeOutput := runJavaScriptBehaviorProbe(t, "done-card implementation span", javascriptProbe)
	var renderedCards map[string]struct {
		DoneLineCount  int      `json:"doneLineCount"`
		DoneLineText   string   `json:"doneLineText"`
		SpanNodeCount  int      `json:"spanNodeCount"`
		SpanText       string   `json:"spanText"`
		SpanTickerKeys []string `json:"spanTickerKeys"`
		MarkerTitle    string   `json:"markerTitle"`
	}
	if decodeError := json.Unmarshal(probeOutput, &renderedCards); decodeError != nil {
		t.Fatalf("decode rendered done lines: %v (output %q)", decodeError, probeOutput)
	}

	renderExpectations := []struct {
		requestId      string
		wantVerb       string
		wantInstantIso string
		wantSpanText   string
		requirement    string
	}{
		{"REQ-901", "done", "2026-08-24T12:45:00Z", "wall time 2h 40m", "an ordinary calibration span is labeled as wall time"},
		{"REQ-902", "done", "2026-08-25T04:05:00Z", "wall time 18h 00m over 4h · assumed pause", "an over-ceiling wall span is marked as a duration-quality assumption, not a workflow state"},
		{"REQ-903", "done", "2026-08-24T10:05:00Z", "reversed stamps", "a reversed span refuses to state a number"},
		{"REQ-904", "done", "2026-08-24T12:45:00Z", "", "no parseable claimed_at leaves the done line exactly as it was"},
		{"REQ-905", "cancelled", "2026-08-24T12:45:00Z", "", "a cancelled card states no duration"},
		{"REQ-906", "done", "2026-08-24T12:45:00Z", "", "a git-dated completion instant states no duration (D-01)"},
		{"REQ-907", "done", "2026-08-24T10:39:00Z", "wall time 34m 00s", "a sub-hour wall span keeps seconds — the chart's \"34.0 min\" is a different vocabulary"},
		{"REQ-908", "done", "2026-08-24T10:05:00Z", "wall time 0s", "a zero-minute wall span states zero, never NaN"},
	}
	if len(renderedCards) != len(renderExpectations) {
		t.Fatalf("probe rendered %d cards, want %d", len(renderedCards), len(renderExpectations))
	}

	sawSpanReading := false
	for _, expectation := range renderExpectations {
		rendered := renderedCards[expectation.requestId]
		if rendered.DoneLineCount != 1 {
			t.Fatalf("%s rendered %d done lines; the card must carry exactly one for this probe to mean anything",
				expectation.requestId, rendered.DoneLineCount)
		}
		if rendered.SpanNodeCount > 1 {
			t.Errorf("%s rendered %d span nodes on one done line", expectation.requestId, rendered.SpanNodeCount)
		}
		if rendered.SpanText != expectation.wantSpanText {
			t.Errorf("%s done line said %q about its span, want %q (%s)",
				expectation.requestId, rendered.SpanText, expectation.wantSpanText, expectation.requirement)
		}
		// Placement: the span rides the done line, after the completion instant.
		wantLineText := expectation.wantVerb + " " + expectation.wantInstantIso + expectation.wantSpanText
		if rendered.DoneLineText != wantLineText {
			t.Errorf("%s done line text = %q, want %q (%s)",
				expectation.requestId, rendered.DoneLineText, wantLineText, expectation.requirement)
		}
		if len(rendered.SpanTickerKeys) != 0 {
			t.Errorf("%s span node carries dataset keys %v; a finished span must carry none — "+
				"refreshRelativeTimeNodes selects [data-instant-ms] and would rewrite it every second as elapsed-since-epoch",
				expectation.requestId, rendered.SpanTickerKeys)
		}
		if strings.Contains(rendered.SpanText, "SKEW-BRANCH-REACHED") {
			t.Errorf("%s reached formatElapsedDuration's clock-skew branch; the Go verdict must be branched on first", expectation.requestId)
		}
		if expectation.requestId == "REQ-902" {
			wantTitle := "Duration-quality marker only: this claim-to-completion wall span is longer than the board's " +
				"single-session ceiling, so it is assumed to include a pause and excluded from duration medians. " +
				"The REQ remains completed."
			if rendered.MarkerTitle != wantTitle {
				t.Errorf("%s marker title = %q, want %q", expectation.requestId, rendered.MarkerTitle, wantTitle)
			}
		}
		if expectation.wantSpanText != "" {
			sawSpanReading = true
		}
	}
	if !sawSpanReading {
		t.Fatalf("no fixture rendered any span reading, so this probe cannot fail on the span text")
	}
}
