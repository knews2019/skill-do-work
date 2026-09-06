package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// The DOM stub these probes share. It is the makeNode() idiom the other
// behavior probes use, plus the three things a fold needs: hidden, an id, and
// a textContent that really empties the node — the drawer clears its meta list
// with drawerMeta.textContent = "" on every open, and a stub that treats that
// as a plain property assignment would report rows from a previous open.
const foldProbeDomStub = `
function makeNode() {
  var node = {
    childNodes: [],
    dataset: {},
    attributes: {},
    listeners: {},
    hidden: false,
    id: "",
    className: "",
    type: "",
    classList: { add: function () {}, remove: function () {} },
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; },
    removeChild: function (childNode) {
      var childIndex = this.childNodes.indexOf(childNode);
      if (childIndex !== -1) { this.childNodes.splice(childIndex, 1); }
      return childNode;
    },
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = attributeValue; },
    getAttribute: function (attributeName) {
      return Object.prototype.hasOwnProperty.call(this.attributes, attributeName)
        ? this.attributes[attributeName]
        : null;
    },
    addEventListener: function (eventName, handler) {
      this.listeners[eventName] = (this.listeners[eventName] || []).concat([handler]);
    },
    dispatch: function (eventName) {
      (this.listeners[eventName] || []).forEach(function (handler) { handler(); });
    }
  };
  var ownTextContent = "";
  Object.defineProperty(node, "textContent", {
    get: function () { return ownTextContent; },
    set: function (nextTextContent) {
      ownTextContent = String(nextTextContent);
      if (ownTextContent === "") { this.childNodes = []; }
    }
  });
  return node;
}
function collectByClassName(node, wantedClassName, found) {
  found = found || [];
  if (node.className === wantedClassName) { found.push(node); }
  (node.childNodes || []).forEach(function (child) { collectByClassName(child, wantedClassName, found); });
  return found;
}
`

// REQ-486: the by-UR lens folds too. It starts open — that is the whole point
// of the reading — but a UR with dozens of REQ cards is a wall the reader
// cannot get past, so activating the head must put the cards away and bring
// exactly the same cards back. The sibling group must not move when it does.
func TestJavaScriptBehaviorByUserRequestLensFoldsAndRestoresItsCards(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function citationMatchedTicketId("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesUserRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function requestMatchesFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestHasOpenOrRecentWork("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestLensEmptyText("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentlyDoneIds("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderUserRequestLens("),
	}
	functionBlocks = append(userRequestSummaryCallSiteBlocks(t, indexHtml), functionBlocks...)
	javascriptProbe := `
Date.now = function () { return Date.parse("2026-08-15T12:00:00Z"); };
var boardData = {
  requests: {
    "REQ-601": { status: "pending", title: "alpha open", domain: "general" },
    "REQ-602": { status: "completed", title: "alpha shipped", domain: "general" },
    "REQ-603": { status: "claimed", title: "beta running", domain: "general" }
  },
  userRequests: {
    "UR-401": { requestIds: ["REQ-601", "REQ-602"], title: "alpha request", inputFilePresent: true },
    "UR-402": { requestIds: ["REQ-603"], title: "beta request", inputFilePresent: true }
  },
  userRequestOrder: ["UR-401", "UR-402"],
  calendar: []
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { view: "board", lens: "user-request", windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "all" };
// The always-open reading. Its fold starts open; URs only starts collapsed.
var userRequestCardsFolded = false;
` + foldProbeDomStub + `
var userRequestLensNode = makeNode();
var document = {
  getElementById: function (nodeId) { return nodeId === "user-request-lens" ? userRequestLensNode : null; },
  createElement: function () { return makeNode(); }
};
function makeRequestCard(requestId) { return { className: "req-card", requestId: requestId }; }
` + strings.Join(functionBlocks, "\n") + `
function collectDrawerTriggers(node, found) {
  found = found || [];
  if (node.dataset && node.dataset.detailKind === "ur") { found.push(node.dataset.detailId); }
  (node.childNodes || []).forEach(function (child) { collectDrawerTriggers(child, found); });
  return found;
}
function headOf(group) { return collectByClassName(group, "ur-group-head")[0]; }
function describeGroups(groups) {
  return groups.map(function (group) {
    var cardIds = [];
    collectByClassName(group, "ur-group-cards").forEach(function (cardsNode) {
      cardsNode.childNodes.forEach(function (card) { cardIds.push(card.requestId); });
    });
    var detailButton = collectByClassName(group, "ur-group-detail")[0];
    return {
      userRequestId: collectByClassName(group, "ur-id")[0].textContent,
      expanded: headOf(group).getAttribute("aria-expanded") || "",
      cardIds: cardIds,
      drawerTriggers: collectDrawerTriggers(group),
      detailButtonId: (detailButton && detailButton.dataset.detailId) || ""
    };
  });
}

userRequestLensNode = makeNode();
renderUserRequestLens();
var groups = userRequestLensNode.childNodes.filter(function (node) { return node.className === "ur-group"; });
var initial = describeGroups(groups);
headOf(groups[0]).dispatch("click");
var afterCollapse = describeGroups(groups);
headOf(groups[0]).dispatch("click");
var afterReopen = describeGroups(groups);

process.stdout.write(JSON.stringify({
  initial: initial,
  afterCollapse: afterCollapse,
  afterReopen: afterReopen
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR fold", javascriptProbe)

	var result struct {
		Initial       []renderedUserRequestRow `json:"initial"`
		AfterCollapse []renderedUserRequestRow `json:"afterCollapse"`
		AfterReopen   []renderedUserRequestRow `json:"afterReopen"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode by-UR fold output: %v (output %q)", decodeError, probeOutput)
	}
	if len(result.Initial) != 2 {
		t.Fatalf("by-UR lens rendered %d groups, want 2: %#v", len(result.Initial), result.Initial)
	}

	// Open on arrival, with the cards actually in the DOM and the head saying so.
	for _, row := range result.Initial {
		if row.Expanded != "true" || len(row.CardIds) == 0 {
			t.Fatalf("by-UR group %s starts aria-expanded=%q with cards %#v; want \"true\" and its cards present",
				row.UserRequestId, row.Expanded, row.CardIds)
		}
	}

	// Collapsing the first group empties it and leaves the second one alone.
	if result.AfterCollapse[0].Expanded != "false" || len(result.AfterCollapse[0].CardIds) != 0 {
		t.Fatalf("collapsed by-UR group = %#v, want aria-expanded \"false\" and no cards", result.AfterCollapse[0])
	}
	if strings.Join(result.AfterCollapse[1].CardIds, ",") != "REQ-603" {
		t.Fatalf("collapsing UR-401 changed UR-402's cards to %#v; each group's fold is its own",
			result.AfterCollapse[1].CardIds)
	}

	// Reopening restores the same cards, in the same order. A fold that rebuilds
	// a different list is a different bug from a fold that does not reopen.
	if result.AfterReopen[0].Expanded != "true" ||
		strings.Join(result.AfterReopen[0].CardIds, ",") != strings.Join(result.Initial[0].CardIds, ",") {
		t.Fatalf("reopened by-UR group = %#v, want aria-expanded \"true\" and the first render's cards %#v",
			result.AfterReopen[0], result.Initial[0].CardIds)
	}

	// The head folds, the sibling button opens the drawer, and neither does both.
	for _, row := range result.Initial {
		if row.DetailButtonId != row.UserRequestId || len(row.DrawerTriggers) != 1 {
			t.Fatalf("by-UR group %s: Details button names %q, drawer triggers %#v; want one trigger on the sibling button",
				row.UserRequestId, row.DetailButtonId, row.DrawerTriggers)
		}
	}
}

// REQ-486: a UR with dozens of grouped REQs turned the drawer's "REQ ids" row
// into a wall that pushed input.md and the body out of reach. The row folds
// now. This is the first probe to drive openUserRequestDetail at all, so it
// also pins the two rows either side of the fold staying put, and it writes
// the open record through the shipped setDetailTarget rather than assigning
// currentDetailKind / currentDetailId by hand.
func TestJavaScriptBehaviorUserRequestDrawerFoldsItsRequestIdList(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeTicketLink("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeTicketLinkList("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendMetaRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendFoldableMetaRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function clearDetailGlossary("),
		sliceBalancedBlockAfter(t, indexHtml, "function setDetailTarget("),
		sliceBalancedBlockAfter(t, indexHtml, "function openUserRequestDetail("),
	}
	functionBlocks = append(userRequestSummaryCallSiteBlocks(t, indexHtml), functionBlocks...)
	javascriptProbe := `
var boardData = {
  requests: {
    "REQ-701": { status: "pending", title: "first grouped request" },
    "REQ-702": { status: "completed", title: "second grouped request" },
    "REQ-703": { status: "claimed", title: "third grouped request" }
  },
  userRequests: {
    "UR-501": {
      requestIds: ["REQ-701", "REQ-702", "REQ-703"],
      title: "long request",
      inputFilePresent: true,
      bodyHtml: "<p>the input.md body</p>"
    }
  }
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var inlineTicketTitleMaxLength = 60;
` + foldProbeDomStub + `
// The ten ids board-detail.js reads at module scope, answered by the stub in
// the same shape the fragment expects.
var drawerNodesById = {};
["detail-resizer", "detail-drawer", "detail-kind", "detail-id", "detail-drawer-title",
 "detail-meta", "detail-body", "detail-glossary", "detail-copy", "detail-copy-all"
].forEach(function (nodeId) {
  var node = makeNode();
  node.id = nodeId;
  drawerNodesById[nodeId] = node;
});
var document = {
  getElementById: function (nodeId) { return drawerNodesById[nodeId] || null; },
  // tagName is recorded because "is the fold control a real button" is one of
  // the claims here: a clickable <dt> is not keyboard-reachable.
  createElement: function (tagName) { var node = makeNode(); node.tagName = String(tagName).toUpperCase(); return node; },
  createTextNode: function (text) { var node = makeNode(); node.textContent = text; return node; },
  addEventListener: function () {},
  activeElement: null
};
var detailResizer = document.getElementById("detail-resizer");
var drawer = document.getElementById("detail-drawer");
var drawerKind = document.getElementById("detail-kind");
var drawerId = document.getElementById("detail-id");
var drawerTitle = document.getElementById("detail-drawer-title");
var drawerMeta = document.getElementById("detail-meta");
var drawerBody = document.getElementById("detail-body");
var drawerGlossary = document.getElementById("detail-glossary");
var currentDetailKind = "";
var currentDetailId = "";
// Out of this probe's scope: the drawer's chrome, its glossary and the
// Activity view's row highlight. setDetailTarget itself is the shipped one.
function syncActivitySelectionToDrawer() {}
function linkifyDetailBody() { return []; }
function renderDetailGlossary() {}
var openDrawerCount = 0;
function showDrawer() { openDrawerCount += 1; }
` + strings.Join(functionBlocks, "\n") + `
// The stub's textContent is a node's OWN text, so a row whose value is an
// element (since REQ-486 the summary rows are spans) needs its subtree read.
function subtreeText(node) {
  var text = node.textContent || "";
  (node.childNodes || []).forEach(function (child) { text += subtreeText(child); });
  return text;
}
function describeMetaRows() {
  var rows = [];
  for (var childIndex = 0; childIndex + 1 < drawerMeta.childNodes.length; childIndex += 2) {
    var term = drawerMeta.childNodes[childIndex];
    var value = drawerMeta.childNodes[childIndex + 1];
    var foldButton = collectByClassName(term, "detail-fold")[0];
    var ticketIds = [];
    collectByClassName(value, "detail-dep-list").forEach(function (listNode) {
      listNode.childNodes.forEach(function (rowNode) {
        ticketIds.push(rowNode.childNodes[0].dataset.detailId || rowNode.childNodes[0].textContent);
      });
    });
    rows.push({
      label: foldButton ? foldButton.textContent : term.textContent,
      foldable: Boolean(foldButton),
      expanded: foldButton ? foldButton.getAttribute("aria-expanded") || "" : "",
      buttonElement: foldButton ? foldButton.tagName : "",
      controls: foldButton ? foldButton.getAttribute("aria-controls") || "" : "",
      valueId: value.id || "",
      valueHidden: Boolean(value.hidden),
      valueText: subtreeText(value),
      ticketIds: ticketIds
    });
  }
  return rows;
}
function foldControl() {
  return collectByClassName(drawerMeta, "detail-fold")[0];
}

openUserRequestDetail("UR-501");
var onOpen = describeMetaRows();
foldControl().dispatch("click");
var afterFold = describeMetaRows();
var bodyAfterFold = drawerBody.innerHTML;
foldControl().dispatch("click");
var afterUnfold = describeMetaRows();
openUserRequestDetail("UR-501");
var onReopen = describeMetaRows();

process.stdout.write(JSON.stringify({
  onOpen: onOpen,
  afterFold: afterFold,
  afterUnfold: afterUnfold,
  onReopen: onReopen,
  bodyAfterFold: bodyAfterFold,
  openDrawerCount: openDrawerCount,
  detailKind: currentDetailKind,
  detailId: currentDetailId
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "UR drawer id-list fold", javascriptProbe)

	type renderedDrawerMetaRow struct {
		Label         string   `json:"label"`
		Foldable      bool     `json:"foldable"`
		Expanded      string   `json:"expanded"`
		ButtonElement string   `json:"buttonElement"`
		Controls      string   `json:"controls"`
		ValueId       string   `json:"valueId"`
		ValueHidden   bool     `json:"valueHidden"`
		ValueText     string   `json:"valueText"`
		TicketIds     []string `json:"ticketIds"`
	}
	var result struct {
		OnOpen          []renderedDrawerMetaRow `json:"onOpen"`
		AfterFold       []renderedDrawerMetaRow `json:"afterFold"`
		AfterUnfold     []renderedDrawerMetaRow `json:"afterUnfold"`
		OnReopen        []renderedDrawerMetaRow `json:"onReopen"`
		BodyAfterFold   string                  `json:"bodyAfterFold"`
		OpenDrawerCount int                     `json:"openDrawerCount"`
		DetailKind      string                  `json:"detailKind"`
		DetailId        string                  `json:"detailId"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode UR drawer fold output: %v (output %q)", decodeError, probeOutput)
	}

	findRow := func(rows []renderedDrawerMetaRow, wantedLabel string) renderedDrawerMetaRow {
		t.Helper()
		for _, row := range rows {
			if strings.Contains(row.Label, wantedLabel) {
				return row
			}
		}
		t.Fatalf("the UR drawer rendered no %q row: %#v", wantedLabel, rows)
		return renderedDrawerMetaRow{}
	}

	// The drawer opened through the shipped writer, not by hand.
	if result.DetailKind != "ur" || result.DetailId != "UR-501" || result.OpenDrawerCount != 2 {
		t.Fatalf("drawer target = %q/%q after %d opens; want ur/UR-501 after 2",
			result.DetailKind, result.DetailId, result.OpenDrawerCount)
	}

	// Open, and open by default: the ids are readable without a click.
	idRowOnOpen := findRow(result.OnOpen, "REQ ids")
	if !idRowOnOpen.Foldable || idRowOnOpen.Expanded != "true" || idRowOnOpen.ValueHidden {
		t.Fatalf("REQ ids row on open = %#v; want a fold control reading aria-expanded \"true\" over a visible list", idRowOnOpen)
	}
	if strings.Join(idRowOnOpen.TicketIds, ",") != "REQ-701,REQ-702,REQ-703" {
		t.Fatalf("REQ ids row listed %#v, want all three grouped REQs", idRowOnOpen.TicketIds)
	}
	// A real button, and one that says which node it hides: a <dt> made
	// clickable is not reachable by keyboard and announces nothing.
	if idRowOnOpen.ButtonElement != "BUTTON" || idRowOnOpen.Controls == "" || idRowOnOpen.Controls != idRowOnOpen.ValueId {
		t.Fatalf("REQ ids fold control = %#v; want a <button> whose aria-controls names the value it hides", idRowOnOpen)
	}

	// Only this row folds. Its neighbours are the rows the wall was hiding.
	for _, plainLabel := range []string{"Grouped REQs", "input.md"} {
		if row := findRow(result.OnOpen, plainLabel); row.Foldable {
			t.Fatalf("the %q row grew a fold control: %#v; REQ-486 folds the id list alone", plainLabel, row)
		}
	}

	// Folding hides the list and touches nothing else in the same pass.
	idRowAfterFold := findRow(result.AfterFold, "REQ ids")
	if idRowAfterFold.Expanded != "false" || !idRowAfterFold.ValueHidden {
		t.Fatalf("folded REQ ids row = %#v; want aria-expanded \"false\" over a hidden list", idRowAfterFold)
	}
	for _, plainLabel := range []string{"Grouped REQs", "input.md"} {
		row := findRow(result.AfterFold, plainLabel)
		if row.ValueHidden || row.ValueText == "" {
			t.Fatalf("the %q row went with the fold: %#v; only the id list may be hidden", plainLabel, row)
		}
	}
	if !strings.Contains(result.BodyAfterFold, "the input.md body") {
		t.Fatalf("drawer body after folding = %q; the fold must not disturb the rendered body", result.BodyAfterFold)
	}

	// Unfolding brings the same ids back, without a rebuild.
	idRowAfterUnfold := findRow(result.AfterUnfold, "REQ ids")
	if idRowAfterUnfold.Expanded != "true" || idRowAfterUnfold.ValueHidden ||
		strings.Join(idRowAfterUnfold.TicketIds, ",") != strings.Join(idRowOnOpen.TicketIds, ",") {
		t.Fatalf("unfolded REQ ids row = %#v, want the first open's list %#v back and visible",
			idRowAfterUnfold, idRowOnOpen.TicketIds)
	}

	// The drawer rebuilds its meta list on every open, so the fold is
	// deliberately per-open state: reopening starts expanded again.
	idRowOnReopen := findRow(result.OnReopen, "REQ ids")
	if idRowOnReopen.Expanded != "true" || idRowOnReopen.ValueHidden {
		t.Fatalf("reopened drawer's REQ ids row = %#v; the meta list is rebuilt per open, so it must start expanded", idRowOnReopen)
	}
}

// The rollup and its formatter, sliced out of the shipped page. Everything here
// is pure: no DOM, no wall clock. The probes below install a Date.now that
// THROWS before calling in, because "nowMs is a parameter and never a Date.now()
// call inside" is the property that lets one tick state one instant across the
// header, the drawer and the cards below them.
func userRequestSummaryRollupBlocks(t *testing.T, indexHtml string) []string {
	t.Helper()
	return []string{
		sliceDeclarationAfter(t, indexHtml, "var futureInstantSkewAllowanceMs"),
		sliceDeclarationAfter(t, indexHtml, "var clockSkewMarkerText"),
		sliceBalancedBlockAfter(t, indexHtml, "function isCompletedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatSpanMinutes("),
		sliceBalancedBlockAfter(t, indexHtml, "function liveClaimElapsedMinutes("),
		sliceBalancedBlockAfter(t, indexHtml, "function summarizeUserRequestProgress("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryDurationText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryPercentageText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryActiveText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryRemainingText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryMetrics("),
	}
}

// The summary functions the two SHIPPED call sites reach into, in dependency
// order. Any probe that drives renderUserRequestLens or openUserRequestDetail
// needs them: since REQ-486 the By UR header builds a progress strip and the UR
// drawer builds progress meta rows, so a probe slicing only the renderer dies
// with a ReferenceError instead of quietly measuring the markup that used to be
// there. Re-declaring a function block a probe already sliced is harmless.
func userRequestSummaryCallSiteBlocks(t *testing.T, indexHtml string) []string {
	t.Helper()
	return append(userRequestSummaryRollupBlocks(t, indexHtml),
		sliceBalancedBlockAfter(t, indexHtml, "function markUserRequestSummaryValueNode("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeUserRequestSummaryStrip("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendUserRequestSummaryMetaRows("),
	)
}

type userRequestSummaryMetric struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
}

func metricValue(t *testing.T, metrics []userRequestSummaryMetric, key string) string {
	t.Helper()
	for _, metric := range metrics {
		if metric.Key == key {
			return metric.Value
		}
	}
	t.Fatalf("no %q metric in %#v", key, metrics)
	return ""
}

// REQ-486: the summary answers for the UR's WHOLE membership. Filters change
// which cards are on screen; they must never move the numerator or the
// denominator, so the rollup reads userRequestsById[id].requestIds and never the
// filtered list. Successful is completed plus completed-with-issues, resolved
// adds cancelled, and failed counts toward neither — the split a reader uses to
// tell "shipped" from "finished with".
func TestJavaScriptBehaviorUserRequestProgressCountsWholeMembership(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := `
var boardData = {
  requests: {
    "REQ-601": { status: "completed" },
    "REQ-602": { status: "completed-with-issues" },
    "REQ-603": { status: "cancelled" },
    "REQ-604": { status: "failed" },
    "REQ-605": { status: "pending" }
  },
  userRequests: {
    "UR-601": { requestIds: ["REQ-601", "REQ-602", "REQ-603", "REQ-604", "REQ-605"] },
    "UR-602": { requestIds: [] }
  },
  timeline: { projection: { confident: false, declinedReason: "not enough completed work to forecast from" } }
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var nowMs = Date.parse("2026-08-15T12:00:00Z");
` + strings.Join(userRequestSummaryRollupBlocks(t, indexHtml), "\n") + `
// The rollup takes its instant as an argument. A wall-clock read inside it
// would give each of a tick's several surfaces a different "now".
Date.now = function () { throw new Error("summarizeUserRequestProgress read the wall clock"); };

var populated = summarizeUserRequestProgress("UR-601", nowMs);
var empty = summarizeUserRequestProgress("UR-602", nowMs);
var missing = summarizeUserRequestProgress("UR-999", nowMs);

process.stdout.write(JSON.stringify({
  populated: populated,
  empty: empty,
  missing: missing,
  populatedMetrics: userRequestSummaryMetrics(populated),
  emptyMetrics: userRequestSummaryMetrics(empty)
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "UR progress membership", javascriptProbe)

	type rollup struct {
		TotalCount      int `json:"totalCount"`
		SuccessfulCount int `json:"successfulCount"`
		ResolvedCount   int `json:"resolvedCount"`
	}
	var result struct {
		Populated        rollup                     `json:"populated"`
		Empty            rollup                     `json:"empty"`
		Missing          rollup                     `json:"missing"`
		PopulatedMetrics []userRequestSummaryMetric `json:"populatedMetrics"`
		EmptyMetrics     []userRequestSummaryMetric `json:"emptyMetrics"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode UR progress membership output: %v (output %q)", decodeError, probeOutput)
	}

	if result.Populated.TotalCount != 5 || result.Populated.SuccessfulCount != 2 || result.Populated.ResolvedCount != 3 {
		t.Fatalf("status split = %+v; want 5 members, 2 successful (completed + completed-with-issues), 3 resolved (adds cancelled), failed in neither",
			result.Populated)
	}
	if got := metricValue(t, result.PopulatedMetrics, "successful"); got != "2/5 (40%)" {
		t.Fatalf("successful metric = %q, want %q — the count and the total have to be visible beside the percentage", got, "2/5 (40%)")
	}
	if got := metricValue(t, result.PopulatedMetrics, "resolved"); got != "3/5 (60%)" {
		t.Fatalf("resolved metric = %q, want %q", got, "3/5 (60%)")
	}

	// Zero members: an explicit token, and no division at all.
	if result.Empty.TotalCount != 0 {
		t.Fatalf("empty UR reported %d members", result.Empty.TotalCount)
	}
	for _, percentageKey := range []string{"successful", "resolved"} {
		if got := metricValue(t, result.EmptyMetrics, percentageKey); got != "unavailable" {
			t.Fatalf("empty UR's %s metric = %q, want the literal unavailable token rather than a division by zero", percentageKey, got)
		}
	}

	// An id the payload does not carry must produce a well-formed result. This
	// is the totality the ticker depends on: the rollup runs inside the board's
	// only interval callback, and a reach into undefined there freezes every
	// stopwatch on the page.
	if result.Missing.TotalCount != 0 || result.Missing.SuccessfulCount != 0 || result.Missing.ResolvedCount != 0 {
		t.Fatalf("unknown UR id produced %+v, want a well-formed empty rollup", result.Missing)
	}
}

// REQ-486: active time is the sum of the spans the Go side already accepted,
// plus live elapsed for currently claimed members — never a span the browser
// measured for itself. Rejected and unmeasured members are disclosed, never
// counted as zero, and a claim stamped ahead of the viewer's clock is disclosed
// rather than clamped.
func TestJavaScriptBehaviorUserRequestProgressSumsAcceptedSpansAndLiveClaims(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := `
var nowMs = Date.parse("2026-08-15T12:00:00Z");
function minutesAgo(minuteCount) { return new Date(nowMs - minuteCount * 60000).toISOString(); }
var boardData = {
  requests: {
    // Accepted, rejected, and unmeasured — the three verdicts the Go side ships.
    "REQ-701": { status: "completed", hasImplementationSpan: true, implementationSpanMinutes: 45, implementationSpanReason: "" },
    "REQ-702": { status: "completed", hasImplementationSpan: true, implementationSpanMinutes: 300, implementationSpanReason: "paused" },
    "REQ-703": { status: "cancelled", implementationSpanMinutes: 0 },
    // Two live claims, 30 and 90 minutes old.
    "REQ-711": { status: "claimed", claimedAt: minutesAgo(30) },
    "REQ-712": { status: "claimed", claimedAt: minutesAgo(90) },
    // A claim stamped ten minutes into the future.
    "REQ-721": { status: "claimed", claimedAt: new Date(nowMs + 10 * 60000).toISOString() }
  },
  userRequests: {
    "UR-701": { requestIds: ["REQ-701", "REQ-702", "REQ-703"] },
    "UR-702": { requestIds: ["REQ-711", "REQ-712"] },
    "UR-703": { requestIds: ["REQ-721"] }
  },
  timeline: { projection: { confident: false, declinedReason: "not enough completed work to forecast from" } }
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
` + strings.Join(userRequestSummaryRollupBlocks(t, indexHtml), "\n") + `
Date.now = function () { throw new Error("summarizeUserRequestProgress read the wall clock"); };

var spans = summarizeUserRequestProgress("UR-701", nowMs);
var live = summarizeUserRequestProgress("UR-702", nowMs);
var liveOneMinuteLater = summarizeUserRequestProgress("UR-702", nowMs + 60000);
var skewed = summarizeUserRequestProgress("UR-703", nowMs);

process.stdout.write(JSON.stringify({
  spans: spans,
  spansMetrics: userRequestSummaryMetrics(spans),
  live: live,
  liveOneMinuteLater: liveOneMinuteLater,
  skewed: skewed,
  skewedMetrics: userRequestSummaryMetrics(skewed)
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "UR progress active time", javascriptProbe)

	type activeRollup struct {
		ActiveMinutes     float64 `json:"activeMinutes"`
		ExcludedSpanCount int     `json:"excludedSpanCount"`
		UnmeasuredCount   int     `json:"unmeasuredCount"`
		SkewedClaimCount  int     `json:"skewedClaimCount"`
		LiveClaimCount    int     `json:"liveClaimCount"`
	}
	var result struct {
		Spans              activeRollup               `json:"spans"`
		SpansMetrics       []userRequestSummaryMetric `json:"spansMetrics"`
		Live               activeRollup               `json:"live"`
		LiveOneMinuteLater activeRollup               `json:"liveOneMinuteLater"`
		Skewed             activeRollup               `json:"skewed"`
		SkewedMetrics      []userRequestSummaryMetric `json:"skewedMetrics"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode UR progress active-time output: %v (output %q)", decodeError, probeOutput)
	}

	// Only the accepted span is summed. The 300-minute paused span is disclosed
	// as excluded, never added; the cancelled member with no span is unmeasured.
	if result.Spans.ActiveMinutes != 45 || result.Spans.ExcludedSpanCount != 1 || result.Spans.UnmeasuredCount != 1 {
		t.Fatalf("span rollup = %+v; want 45 accepted minutes with 1 excluded and 1 unmeasured", result.Spans)
	}
	activeText := metricValue(t, result.SpansMetrics, "active")
	for _, requiredFragment := range []string{"at least", "45 min", "1 excluded", "1 unmeasured"} {
		if !strings.Contains(activeText, requiredFragment) {
			t.Fatalf("active metric = %q, want it to carry %q — a known-partial sum stated as complete is the failure this REQ names",
				activeText, requiredFragment)
		}
	}

	// Live claims: exactly the elapsed minutes, and exactly the clock delta one
	// minute later. A drifting or re-derived figure fails on the second number.
	if result.Live.ActiveMinutes != 120 || result.Live.LiveClaimCount != 2 {
		t.Fatalf("live rollup = %+v; want 120 minutes from two claims 30 and 90 minutes old", result.Live)
	}
	if result.LiveOneMinuteLater.ActiveMinutes != 122 {
		t.Fatalf("live rollup one minute later = %v minutes, want exactly 122 — the tick advances the sum by the clock delta and nothing else",
			result.LiveOneMinuteLater.ActiveMinutes)
	}

	// A future claim stamp is disclosed, not clamped to a silent zero.
	if result.Skewed.SkewedClaimCount != 1 || result.Skewed.LiveClaimCount != 0 || result.Skewed.ActiveMinutes != 0 {
		t.Fatalf("skewed rollup = %+v; want the claim disclosed as skewed and contributing nothing", result.Skewed)
	}
	if skewText := metricValue(t, result.SkewedMetrics, "active"); !strings.Contains(skewText, "clock skew") {
		t.Fatalf("active metric for a future-stamped claim = %q, want the clock-skew qualifier present", skewText)
	}
}

// REQ-486: remaining active time has three arms and each has to be able to fail
// in both directions — a saved p50_active_minutes, else the Timeline's median
// for that member's effort class but only while the Timeline calls itself
// confident, else unknown. A claimed member spends its estimate down and floors
// at zero rather than going negative, and a failed member is unknown even when
// it carries a saved figure.
func TestJavaScriptBehaviorUserRequestProgressForecastsRemainingTime(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := `
var nowMs = Date.parse("2026-08-15T12:00:00Z");
function minutesAgo(minuteCount) { return new Date(nowMs - minuteCount * 60000).toISOString(); }
var boardData = {
  requests: {
    "REQ-801": { status: "pending", hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 75 },
    "REQ-802": { status: "pending", estimateP50ActiveMinutes: 0, effortEstimate: "effort-substantive" },
    "REQ-803": { status: "blocked", estimateP50ActiveMinutes: 0, effortEstimate: "effort-mechanical" },
    "REQ-804": { status: "failed", hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 200 },
    "REQ-805": { status: "claimed", claimedAt: minutesAgo(100), hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 90 },
    "REQ-806": { status: "pending-answers", estimateP50ActiveMinutes: 0, effortEstimate: "effort-substantive" },
    "REQ-807": { status: "claimed", claimedAt: minutesAgo(600), hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 20 },
    "REQ-808": { status: "claimed", claimedAt: minutesAgo(240), hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 45 },
    "REQ-809": { status: "claimed", claimedAt: "2099-01-01T00:00:00Z", hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 130 }
  },
  userRequests: {
    "UR-801": { requestIds: ["REQ-801", "REQ-802", "REQ-803", "REQ-804", "REQ-805"] },
    "UR-802": { requestIds: ["REQ-806"] },
    "UR-803": { requestIds: ["REQ-807", "REQ-808"] },
    "UR-804": { requestIds: ["REQ-809"] }
  },
  timeline: { projection: { confident: true, normalMinutes: 60, trivialMinutes: 15 } }
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
` + strings.Join(userRequestSummaryRollupBlocks(t, indexHtml), "\n") + `
Date.now = function () { throw new Error("summarizeUserRequestProgress read the wall clock"); };

var confident = summarizeUserRequestProgress("UR-801", nowMs);
var confidentSingle = summarizeUserRequestProgress("UR-802", nowMs);
// The same member, with the Timeline no longer able to call its own median.
boardData.timeline.projection = { confident: false, declinedReason: "not enough completed work to forecast from" };
var declined = summarizeUserRequestProgress("UR-802", nowMs);
// Every known member has already outrun its estimate, so every contribution
// floors at zero and the sum is a clean 0.
var allFloored = summarizeUserRequestProgress("UR-803", nowMs);
// One claimed member whose stamp the board refuses.
var skewed = summarizeUserRequestProgress("UR-804", nowMs);

process.stdout.write(JSON.stringify({
  confident: confident,
  confidentMetrics: userRequestSummaryMetrics(confident),
  confidentSingle: confidentSingle,
  declined: declined,
  declinedMetrics: userRequestSummaryMetrics(declined),
  allFloored: allFloored,
  allFlooredMetrics: userRequestSummaryMetrics(allFloored),
  skewed: skewed,
  skewedMetrics: userRequestSummaryMetrics(skewed)
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "UR progress remaining time", javascriptProbe)

	type forecastRollup struct {
		RemainingMinutes     float64 `json:"remainingMinutes"`
		KnownForecastCount   int     `json:"knownForecastCount"`
		UnknownForecastCount int     `json:"unknownForecastCount"`
		OverrunForecastCount int     `json:"overrunForecastCount"`
		RemainingIsPartial   bool    `json:"remainingIsPartial"`
	}
	var result struct {
		Confident         forecastRollup             `json:"confident"`
		ConfidentMetrics  []userRequestSummaryMetric `json:"confidentMetrics"`
		ConfidentSingle   forecastRollup             `json:"confidentSingle"`
		Declined          forecastRollup             `json:"declined"`
		DeclinedMetrics   []userRequestSummaryMetric `json:"declinedMetrics"`
		AllFloored        forecastRollup             `json:"allFloored"`
		AllFlooredMetrics []userRequestSummaryMetric `json:"allFlooredMetrics"`
		Skewed            forecastRollup             `json:"skewed"`
		SkewedMetrics     []userRequestSummaryMetric `json:"skewedMetrics"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode UR progress forecast output: %v (output %q)", decodeError, probeOutput)
	}

	// 75 saved + 60 confident-normal + 15 confident-mechanical + 0 from the
	// over-running claim (90 estimated against 100 elapsed, floored) = 150. The
	// failed member is unknown despite carrying a saved 200.
	if result.Confident.RemainingMinutes != 150 {
		t.Fatalf("remaining minutes = %v, want 150 (75 saved + 60 normal median + 15 mechanical median + a floored-at-zero over-running claim)",
			result.Confident.RemainingMinutes)
	}
	if result.Confident.UnknownForecastCount != 1 || !result.Confident.RemainingIsPartial ||
		result.Confident.OverrunForecastCount != 1 {
		t.Fatalf("forecast rollup = %+v; the failed member must be unknown, the over-running claim must be counted, "+
			"and the known sum must say it is partial", result.Confident)
	}
	remainingText := metricValue(t, result.ConfidentMetrics, "remaining")
	for _, requiredFragment := range []string{"~", "at least", "1 unknown", "1 over estimate"} {
		if !strings.Contains(remainingText, requiredFragment) {
			t.Fatalf("remaining metric = %q, want it to carry %q — a forecast has to read as approximate and a partial one has to say so",
				remainingText, requiredFragment)
		}
	}

	// The fallback arm alone: one pending member, no saved figure, a confident
	// Timeline. It reads the median; with the Timeline declined it reads unknown.
	if result.ConfidentSingle.RemainingMinutes != 60 || result.ConfidentSingle.UnknownForecastCount != 0 {
		t.Fatalf("confident fallback = %+v, want the Timeline's 60-minute normal median and nothing unknown", result.ConfidentSingle)
	}
	if result.Declined.RemainingMinutes != 0 || result.Declined.UnknownForecastCount != 1 || !result.Declined.RemainingIsPartial {
		t.Fatalf("declined fallback = %+v; a Timeline that will not call itself confident must produce unknown, never a borrowed median",
			result.Declined)
	}
	if declinedText := metricValue(t, result.DeclinedMetrics, "remaining"); !strings.Contains(declinedText, "unknown") {
		t.Fatalf("remaining metric with no usable forecast = %q, want the unknown token rather than a zero", declinedText)
	}

	// The floor is counted, not only applied. UR-803's two claimed members have
	// each run past their saved estimate (600 minutes against 20, 240 against
	// 45), so every known contribution is zero and the sum is a true zero. The
	// defect this pins: without the count the strip reads a clean "~0 min",
	// which a reader takes for nearly done and which means the exact opposite.
	// This is the shape four of five live-claim user requests had on the real
	// board.
	if result.AllFloored.RemainingMinutes != 0 || result.AllFloored.KnownForecastCount != 2 ||
		result.AllFloored.UnknownForecastCount != 0 || result.AllFloored.OverrunForecastCount != 2 {
		t.Fatalf("all-floored rollup = %+v; want a zero sum from two known members, both counted as over estimate",
			result.AllFloored)
	}
	allFlooredText := metricValue(t, result.AllFlooredMetrics, "remaining")
	if !strings.Contains(allFlooredText, "2 over estimate") {
		t.Fatalf("remaining metric with every known member past its estimate = %q; want it to carry \"2 over estimate\". "+
			"A bare zero cannot be told apart from nearly finished work", allFlooredText)
	}

	// A claim stamped ahead of the viewer's clock: the Active figure already
	// refuses that stamp, so how much of the estimate has been spent is
	// unreadable and the member is unknown here too. Charging it the full 130
	// minutes as a KNOWN forecast would build a clean-looking figure on a stamp
	// the same rollup just rejected.
	if result.Skewed.RemainingMinutes != 0 || result.Skewed.KnownForecastCount != 0 ||
		result.Skewed.UnknownForecastCount != 1 || !result.Skewed.RemainingIsPartial {
		t.Fatalf("skewed-claim rollup = %+v; want the member unknown rather than charged its full saved estimate",
			result.Skewed)
	}
	if skewedText := metricValue(t, result.SkewedMetrics, "remaining"); skewedText != "unknown" {
		t.Fatalf("remaining metric for a UR whose only member has an unreadable claim stamp = %q, want %q",
			skewedText, "unknown")
	}
}

// The DOM stub the two ticking probes share, on top of foldProbeDomStub. The
// board's tick is an attribute pass — refreshRelativeTimeNodes walks
// [data-instant-ms] and the UR summary pass walks [data-ur-summary-id] — so the
// stub needs a real document tree and a querySelectorAll that only finds what is
// still in it. That is the whole reason the design has no subscriber registry:
// a node the renderer dropped is simply not selected.
const tickProbeDocumentStub = `
var documentRoot = makeNode();
var nodesById = {};
function registerNodeById(nodeId) {
  var node = makeNode();
  node.id = nodeId;
  nodesById[nodeId] = node;
  documentRoot.appendChild(node);
  return node;
}
function datasetKeyFromAttributeSelector(selector) {
  return selector.replace("[", "").replace("]", "").replace(/^data-/, "")
    .replace(/-([a-z])/g, function (wholeMatch, letter) { return letter.toUpperCase(); });
}
function collectByDatasetKey(node, datasetKey, found) {
  found = found || [];
  if (node.dataset && node.dataset[datasetKey] !== undefined) { found.push(node); }
  (node.childNodes || []).forEach(function (child) { collectByDatasetKey(child, datasetKey, found); });
  return found;
}
var wallClockReadCount = 0;
var stubbedNowMs = Date.parse("2026-08-15T12:00:00Z");
Date.now = function () { wallClockReadCount += 1; return stubbedNowMs; };
var document = {
  getElementById: function (nodeId) { return nodesById[nodeId] || null; },
  createElement: function (tagName) { var node = makeNode(); node.tagName = String(tagName).toUpperCase(); return node; },
  createTextNode: function (text) { var node = makeNode(); node.textContent = text; return node; },
  querySelectorAll: function (selector) {
    return collectByDatasetKey(documentRoot, datasetKeyFromAttributeSelector(selector));
  },
  addEventListener: function () {},
  activeElement: null
};
`

// The clock fragment plus the whole summary fragment, sliced out of the shipped
// page in dependency order.
func userRequestSummaryRenderBlocks(t *testing.T, indexHtml string) []string {
	t.Helper()
	return []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceDeclarationAfter(t, indexHtml, "var futureInstantSkewAllowanceMs"),
		sliceDeclarationAfter(t, indexHtml, "var clockSkewMarkerText"),
		sliceDeclarationAfter(t, indexHtml, "var futureStampCauseText"),
		sliceDeclarationAfter(t, indexHtml, "var clockSkewExplanationText"),
		sliceBalancedBlockAfter(t, indexHtml, "function formatRelativeTime("),
		sliceBalancedBlockAfter(t, indexHtml, "function formatElapsedDuration("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeElapsedDurationNode("),
		sliceBalancedBlockAfter(t, indexHtml, "function syncClockSkewTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function refreshRelativeTimeNodes("),
		sliceBalancedBlockAfter(t, indexHtml, "function refreshTickingSurfaces("),
		sliceBalancedBlockAfter(t, indexHtml, "function isCompletedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatSpanMinutes("),
		sliceBalancedBlockAfter(t, indexHtml, "function liveClaimElapsedMinutes("),
		sliceBalancedBlockAfter(t, indexHtml, "function summarizeUserRequestProgress("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryDurationText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryPercentageText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryActiveText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryRemainingText("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryMetrics("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestSummaryMetricsByKey("),
		sliceBalancedBlockAfter(t, indexHtml, "function markUserRequestSummaryValueNode("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeUserRequestSummaryStrip("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendUserRequestSummaryMetaRows("),
		sliceBalancedBlockAfter(t, indexHtml, "function refreshUserRequestSummaryNodes("),
	}
}

// REQ-486's central claim, and the one a probe calling summarizeUserRequestProgress
// twice could never hold: the By UR header and the UR drawer must state the SAME
// five figures, because they are two renderings of one rollup. This drives both
// SHIPPED call sites in one run under one stubbed clock and compares what each
// one actually put on screen. It then applies a status filter that hides members
// and re-reads: filters move cards, never the summary or its denominator. Last it
// advances the clock one minute through refreshTickingSurfaces and checks the
// header figure, the drawer figure and a claimed card's stopwatch all moved.
func TestJavaScriptBehaviorUserRequestSummaryAgreesOnBothSurfaces(t *testing.T) {
	indexHtml := generateLiveSite(t)
	lensBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function citationMatchedTicketId("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function searchMatchesUserRequest("),
		sliceBalancedBlockAfter(t, indexHtml, "function requestMatchesFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestHasOpenOrRecentWork("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function userRequestLensEmptyText("),
		sliceBalancedBlockAfter(t, indexHtml, "function recentlyDoneIds("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderUserRequestLens("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeTicketLink("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeTicketLinkList("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendMetaRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendFoldableMetaRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function clearDetailGlossary("),
		sliceBalancedBlockAfter(t, indexHtml, "function setDetailTarget("),
		sliceBalancedBlockAfter(t, indexHtml, "function openUserRequestDetail("),
	}
	javascriptProbe := `
var claimInstantMs = Date.parse("2026-08-15T11:30:00Z");
var boardData = {
  requests: {
    "REQ-901": { status: "completed", title: "shipped", domain: "general",
      hasImplementationSpan: true, implementationSpanMinutes: 45, implementationSpanReason: "" },
    "REQ-902": { status: "claimed", title: "running", domain: "general",
      claimedAt: new Date(claimInstantMs).toISOString() },
    "REQ-903": { status: "pending", title: "queued", domain: "general",
      hasEstimateP50ActiveMinutes: true, estimateP50ActiveMinutes: 75 }
  },
  userRequests: {
    "UR-901": { requestIds: ["REQ-901", "REQ-902", "REQ-903"], title: "alpha request", inputFilePresent: true,
      bodyHtml: "<p>the input.md body</p>" }
  },
  userRequestOrder: ["UR-901"],
  calendar: [],
  timeline: { projection: { confident: true, normalMinutes: 60, trivialMinutes: 15 } }
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { view: "board", lens: "user-request", windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "all" };
var userRequestCardsFolded = false;
var inlineTicketTitleMaxLength = 60;
` + foldProbeDomStub + tickProbeDocumentStub + `
var userRequestLensNode = registerNodeById("user-request-lens");
["detail-resizer", "detail-drawer", "detail-kind", "detail-id", "detail-drawer-title",
 "detail-meta", "detail-body", "detail-glossary", "detail-copy", "detail-copy-all"
].forEach(registerNodeById);
var detailResizer = document.getElementById("detail-resizer");
var drawer = document.getElementById("detail-drawer");
var drawerKind = document.getElementById("detail-kind");
var drawerId = document.getElementById("detail-id");
var drawerTitle = document.getElementById("detail-drawer-title");
var drawerMeta = document.getElementById("detail-meta");
var drawerBody = document.getElementById("detail-body");
var drawerGlossary = document.getElementById("detail-glossary");
var currentDetailKind = "";
var currentDetailId = "";
function syncActivitySelectionToDrawer() {}
function linkifyDetailBody() { return []; }
function renderDetailGlossary() {}
function showDrawer() {}
// The card is stubbed, but its stopwatch is the SHIPPED one: this probe's last
// claim is that the header, the drawer and a claim stopwatch move together, and
// a hand-rolled stopwatch node would make that claim about the probe.
function makeRequestCard(requestId) {
  var card = makeNode();
  card.className = "req-card";
  card.requestId = requestId;
  var request = requestsById[requestId];
  if (request && request.status === "claimed") {
    var stopwatch = makeElapsedDurationNode(request.claimedAt);
    if (stopwatch) { card.appendChild(stopwatch); }
  }
  return card;
}
` + strings.Join(append(userRequestSummaryRenderBlocks(t, indexHtml), lensBlocks...), "\n") + `
function readSummaryMetrics(node) {
  return collectByDatasetKey(node, "urSummaryMetric").map(function (metricNode) {
    return metricNode.dataset.urSummaryMetric + "=" + metricNode.textContent;
  });
}
function headerMetrics() { return readSummaryMetrics(userRequestLensNode); }
function drawerMetrics() { return readSummaryMetrics(drawerMeta); }
function cardStopwatchLabel() {
  var found = collectByDatasetKey(userRequestLensNode, "instantMs");
  return found.length > 0 ? found[0].textContent : "";
}
function groupCountText() {
  var counts = collectByClassName(userRequestLensNode, "ur-count");
  return counts.length > 0 ? counts[0].textContent : "";
}

renderUserRequestLens();
openUserRequestDetail("UR-901");
var headerOnLoad = headerMetrics();
var drawerOnLoad = drawerMetrics();
var stopwatchOnLoad = cardStopwatchLabel();
var unfilteredCountText = groupCountText();

// The tick: one wall-clock read for the whole pass, every surface from it.
wallClockReadCount = 0;
refreshTickingSurfaces();
var wallClockReadsPerTick = wallClockReadCount;

stubbedNowMs += 60000;
refreshTickingSurfaces();
var headerAfterTick = headerMetrics();
var drawerAfterTick = drawerMetrics();
var stopwatchAfterTick = cardStopwatchLabel();

// A filter that hides two of the three members. Cards move; the summary must not.
filterState.status = "pending";
renderUserRequestLens();
openUserRequestDetail("UR-901");
var headerFiltered = headerMetrics();
var drawerFiltered = drawerMetrics();
var filteredCountText = groupCountText();
var filteredCardCount = collectByClassName(userRequestLensNode, "req-card").length;

process.stdout.write(JSON.stringify({
  headerOnLoad: headerOnLoad,
  drawerOnLoad: drawerOnLoad,
  stopwatchOnLoad: stopwatchOnLoad,
  wallClockReadsPerTick: wallClockReadsPerTick,
  headerAfterTick: headerAfterTick,
  drawerAfterTick: drawerAfterTick,
  stopwatchAfterTick: stopwatchAfterTick,
  headerFiltered: headerFiltered,
  drawerFiltered: drawerFiltered,
  unfilteredCountText: unfilteredCountText,
  filteredCountText: filteredCountText,
  filteredCardCount: filteredCardCount
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "UR summary on both surfaces", javascriptProbe)

	var result struct {
		HeaderOnLoad          []string `json:"headerOnLoad"`
		DrawerOnLoad          []string `json:"drawerOnLoad"`
		StopwatchOnLoad       string   `json:"stopwatchOnLoad"`
		WallClockReadsPerTick int      `json:"wallClockReadsPerTick"`
		HeaderAfterTick       []string `json:"headerAfterTick"`
		DrawerAfterTick       []string `json:"drawerAfterTick"`
		StopwatchAfterTick    string   `json:"stopwatchAfterTick"`
		HeaderFiltered        []string `json:"headerFiltered"`
		DrawerFiltered        []string `json:"drawerFiltered"`
		UnfilteredCountText   string   `json:"unfilteredCountText"`
		FilteredCountText     string   `json:"filteredCountText"`
		FilteredCardCount     int      `json:"filteredCardCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode UR summary agreement output: %v (output %q)", decodeError, probeOutput)
	}

	if len(result.HeaderOnLoad) != 5 {
		t.Fatalf("the By UR header rendered %d metrics %#v, want the five the request names", len(result.HeaderOnLoad), result.HeaderOnLoad)
	}
	if !reflect.DeepEqual(result.HeaderOnLoad, result.DrawerOnLoad) {
		t.Fatalf("the two surfaces disagree:\n header %#v\n drawer %#v", result.HeaderOnLoad, result.DrawerOnLoad)
	}
	// 45 accepted minutes plus a claim 30 minutes old, and 75 saved plus the
	// Timeline's 60-minute median spent 30 minutes down.
	wantOnLoad := []string{
		"total=3",
		"active=1h 15m",
		"remaining=~1h 45m",
		"successful=1/3 (33%)",
		"resolved=1/3 (33%)",
	}
	if !reflect.DeepEqual(result.HeaderOnLoad, wantOnLoad) {
		t.Fatalf("rendered metrics = %#v, want %#v", result.HeaderOnLoad, wantOnLoad)
	}
	if result.StopwatchOnLoad != "30m 00s" {
		t.Fatalf("claimed card stopwatch = %q, want 30m 00s", result.StopwatchOnLoad)
	}

	// One instant per tick. Two reads would let the header and the card below it
	// straddle a second boundary and disagree by a minute at every rollover.
	if result.WallClockReadsPerTick != 1 {
		t.Fatalf("one tick read the wall clock %d times, want exactly 1", result.WallClockReadsPerTick)
	}

	// A minute later every live surface moved, and by the same minute.
	wantAfterTick := []string{
		"total=3",
		"active=1h 16m",
		"remaining=~1h 44m",
		"successful=1/3 (33%)",
		"resolved=1/3 (33%)",
	}
	if !reflect.DeepEqual(result.HeaderAfterTick, wantAfterTick) {
		t.Fatalf("header after one minute = %#v, want %#v", result.HeaderAfterTick, wantAfterTick)
	}
	if !reflect.DeepEqual(result.DrawerAfterTick, wantAfterTick) {
		t.Fatalf("drawer after one minute = %#v, want %#v — the drawer rides the same tick as the header", result.DrawerAfterTick, wantAfterTick)
	}
	if result.StopwatchAfterTick != "31m 00s" {
		t.Fatalf("claimed card stopwatch after one minute = %q, want 31m 00s", result.StopwatchAfterTick)
	}

	// Filters move cards and nothing else. The metrics compare against the
	// post-tick values because the clock advanced between the two renders.
	if result.FilteredCardCount != 1 {
		t.Fatalf("the status filter left %d cards on screen, want 1 — the fixture has to actually be filtered for the next claim to mean anything",
			result.FilteredCardCount)
	}
	if !reflect.DeepEqual(result.HeaderFiltered, wantAfterTick) || !reflect.DeepEqual(result.DrawerFiltered, wantAfterTick) {
		t.Fatalf("a status filter moved the summary:\n header %#v\n drawer %#v\n want   %#v",
			result.HeaderFiltered, result.DrawerFiltered, wantAfterTick)
	}
	// The strip owns the total, so the group count becomes filter-only: silent
	// when everything is shown, "1 of 3 shown" when a filter hides members.
	if result.UnfilteredCountText != "" {
		t.Fatalf("unfiltered group count = %q, want nothing — the summary strip states the total", result.UnfilteredCountText)
	}
	if result.FilteredCountText != "1 of 3 shown" {
		t.Fatalf("filtered group count = %q, want %q", result.FilteredCountText, "1 of 3 shown")
	}
}

// The single biggest risk in REQ-486, asserted positively. web/board.js's
// interval is the board's ONLY ticker: every claim stopwatch, every relative
// time, every state timer and the clock-skew tooltip hang off it. Pointing it at
// a function that runs a second pass means an unguarded throw in that pass kills
// the callback — the board paints correctly once and then never updates again,
// which reads as a queue full of very young claims and passes every other test
// in the suite.
//
// Two things contain it. THIS probe holds one of them: the rollup is total by
// narrowing rather than by a try/catch that would hide the bug instead of the
// freeze. The payload below is missing everything the summary pass wants — no
// userRequests entry, no timeline block at all — and a stale summary node still
// names a UR that is gone. The other containment, that the existing refresh runs
// FIRST with the captured instant, is held by
// TestJavaScriptBehaviorTickRefreshesExistingSurfacesBeforeTheSummaryPass below:
// this payload never throws, so the stopwatch here advances whichever pass ran
// first and says nothing about the order.
func TestJavaScriptBehaviorTickSurvivesAnIncompleteUserRequestPayload(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := `
var boardData = { requests: {} };
var requestsById = boardData.requests;
var userRequestsById = {};
` + foldProbeDomStub + tickProbeDocumentStub + `
` + strings.Join(userRequestSummaryRenderBlocks(t, indexHtml), "\n") + `
// An existing ticking surface: a claim stopwatch stamped 90 seconds ago.
var stopwatch = makeElapsedDurationNode(new Date(stubbedNowMs - 90000).toISOString());
documentRoot.appendChild(stopwatch);
// A summary node left over from a render whose UR the payload no longer carries.
var strandedSummary = makeNode();
strandedSummary.dataset.urSummaryId = "UR-404";
strandedSummary.dataset.urSummaryMetric = "active";
strandedSummary.textContent = "stale";
documentRoot.appendChild(strandedSummary);

var labelBeforeTick = stopwatch.textContent;
stubbedNowMs += 60000;
var tickThrew = "";
try {
  refreshTickingSurfaces();
} catch (tickError) {
  tickThrew = String((tickError && tickError.message) || tickError);
}

process.stdout.write(JSON.stringify({
  labelBeforeTick: labelBeforeTick,
  labelAfterTick: stopwatch.textContent,
  strandedSummaryText: strandedSummary.textContent,
  tickThrew: tickThrew
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "tick freeze guard", javascriptProbe)

	var result struct {
		LabelBeforeTick     string `json:"labelBeforeTick"`
		LabelAfterTick      string `json:"labelAfterTick"`
		StrandedSummaryText string `json:"strandedSummaryText"`
		TickThrew           string `json:"tickThrew"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode freeze-guard output: %v (output %q)", decodeError, probeOutput)
	}

	// The positive assertion: the existing surface shows its NEW label. A frozen
	// board is one where this still reads 1m 30s.
	if result.LabelBeforeTick != "1m 30s" || result.LabelAfterTick != "2m 30s" {
		t.Fatalf("claim stopwatch went %q -> %q across one tick against an incomplete payload; want 1m 30s -> 2m 30s. "+
			"A stopwatch that did not move is the board's only ticker having died inside the UR summary pass.",
			result.LabelBeforeTick, result.LabelAfterTick)
	}
	// And the pass itself finished. The catch above exists to REPORT a throw, not
	// to tolerate one: production has no try/catch here, so a throw would take
	// the interval callback with it.
	if result.TickThrew != "" {
		t.Fatalf("refreshTickingSurfaces threw %q against an incomplete UR payload; the rollup has to be total by narrowing", result.TickThrew)
	}
	// The stranded node names a UR the payload no longer has. It must render the
	// well-formed empty rollup rather than reaching into undefined.
	if result.StrandedSummaryText == "stale" || result.StrandedSummaryText == "" {
		t.Fatalf("a summary node for an unknown UR read %q after the tick; want the rollup's own empty-membership text",
			result.StrandedSummaryText)
	}
}

// REQ-486 remediation: the ORDER inside the tick, made observable.
//
// refreshTickingSurfaces runs refreshRelativeTimeNodes first and the UR summary
// pass second, and board-core.js calls that order load-bearing. The freeze-guard
// probe above cannot see it: its payload never makes the summary pass throw, so
// the stopwatch advances whichever pass ran first and the assertion holds either
// way. That is a comment claiming a guarantee no test holds.
//
// This probe makes the summary pass throw on purpose — the shipped code walks
// [data-ur-summary-id] through document.querySelectorAll, so a stub that refuses
// exactly that selector fails the second pass and nothing else — and then checks
// the existing claim stopwatch still shows its NEW label. It can only show the
// new label if the older pass had already run. Swap the two calls in
// board-core.js and this goes red on the stopwatch; nothing else in the suite
// notices the swap.
func TestJavaScriptBehaviorTickRefreshesExistingSurfacesBeforeTheSummaryPass(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := `
var boardData = { requests: {} };
var requestsById = boardData.requests;
var userRequestsById = {};
` + foldProbeDomStub + tickProbeDocumentStub + `
` + strings.Join(userRequestSummaryRenderBlocks(t, indexHtml), "\n") + `
var stopwatch = makeElapsedDurationNode(new Date(stubbedNowMs - 90000).toISOString());
documentRoot.appendChild(stopwatch);

// The injected failure. It is scoped to the summary pass's own selector so the
// relative-time pass is untouched — this probe is about which one ran, not about
// a document that is broken for everyone.
var summarySelectorReads = 0;
var querySelectorAllWithoutTheSummaryPass = document.querySelectorAll;
document.querySelectorAll = function (selector) {
  if (selector === "[data-ur-summary-id]") {
    summarySelectorReads += 1;
    throw new Error("injected: the UR summary pass failed");
  }
  return querySelectorAllWithoutTheSummaryPass.call(document, selector);
};

var labelBeforeTick = stopwatch.textContent;
stubbedNowMs += 60000;
var tickThrew = "";
try {
  refreshTickingSurfaces();
} catch (tickError) {
  tickThrew = String((tickError && tickError.message) || tickError);
}

process.stdout.write(JSON.stringify({
  labelBeforeTick: labelBeforeTick,
  labelAfterTick: stopwatch.textContent,
  summarySelectorReads: summarySelectorReads,
  tickThrew: tickThrew
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "tick ordering", javascriptProbe)

	var result struct {
		LabelBeforeTick      string `json:"labelBeforeTick"`
		LabelAfterTick       string `json:"labelAfterTick"`
		SummarySelectorReads int    `json:"summarySelectorReads"`
		TickThrew            string `json:"tickThrew"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode tick ordering output: %v (output %q)", decodeError, probeOutput)
	}

	// The failure has to have actually reached the summary pass. Without this the
	// probe would pass on a tick that never ran that pass at all, which proves
	// nothing about the order.
	if result.SummarySelectorReads != 1 || result.TickThrew != "injected: the UR summary pass failed" {
		t.Fatalf("the injected summary-pass failure fired %d time(s) and the tick reported %q; "+
			"want exactly one read of [data-ur-summary-id] and that failure escaping the tick",
			result.SummarySelectorReads, result.TickThrew)
	}
	// The point of the ordering: the older surfaces already have their new
	// labels when the newer pass dies. If the summary pass ran first, the
	// stopwatch never gets refreshed and still reads its OLD label.
	if result.LabelBeforeTick != "1m 30s" || result.LabelAfterTick != "2m 30s" {
		t.Fatalf("claim stopwatch went %q -> %q across a tick whose UR summary pass threw; want 1m 30s -> 2m 30s. "+
			"A stopwatch left at its old label means refreshUserRequestSummaryNodes runs BEFORE "+
			"refreshRelativeTimeNodes, and the order board-core.js calls load-bearing is gone.",
			result.LabelBeforeTick, result.LabelAfterTick)
	}
}

// REQ-486 remediation: the two structural claims, as assertions instead of a
// one-time grep in the REQ record.
//
// Both are design rules with a real failure behind them. A try/catch anywhere in
// the summary path would swallow the exception that a broken rollup throws,
// which hides the bug and leaves the board silently wrong instead of loudly
// frozen — the freeze-guard probe would still pass. Reading `completedAt` here
// would let the browser measure a finished span for itself, with no outlier rule
// and no origin correction, and contradict the `took …` badge on the card below.
//
// Scoped to the summary path sliced out of the SHIPPED page, and matched on
// syntax rather than on the bare words: the module's own comments state both
// rules in prose, and "registry" contains "try".
func TestJavaScriptBehaviorUserRequestSummaryPathCarriesNoCatchAndNoCompletedAt(t *testing.T) {
	indexHtml := generateLiveSite(t)
	summaryPathSource := strings.Join(
		append(userRequestSummaryCallSiteBlocks(t, indexHtml),
			sliceBalancedBlockAfter(t, indexHtml, "function refreshUserRequestSummaryNodes(")),
		"\n")

	for _, forbidden := range []struct {
		token  string
		reason string
	}{
		{token: "try {", reason: "the rollup is total by NARROWING; a catch here hides the bug instead of the freeze it causes"},
		{token: "try{", reason: "the rollup is total by NARROWING; a catch here hides the bug instead of the freeze it causes"},
		{token: "catch (", reason: "the rollup is total by NARROWING; a catch here hides the bug instead of the freeze it causes"},
		{token: "catch(", reason: "the rollup is total by NARROWING; a catch here hides the bug instead of the freeze it causes"},
		{token: "completedAt", reason: "the Go side owns the span verdict and its origin rule; measuring one here would contradict the card's own `took …` badge"},
	} {
		if strings.Contains(summaryPathSource, forbidden.token) {
			t.Errorf("the shipped UR summary path contains %q — %s", forbidden.token, forbidden.reason)
		}
	}
}
