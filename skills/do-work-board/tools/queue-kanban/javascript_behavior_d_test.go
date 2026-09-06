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
