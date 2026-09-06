package main

import (
	"encoding/json"
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
      valueText: value.textContent || "",
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
