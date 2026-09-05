package main

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestJavaScriptBehaviorTicketMentionTitlesAndGlossary(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeRequestStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildRequestIdByReqSegment("),
		sliceBalancedBlockAfter(t, indexHtml, "function resolveTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function isAmbiguousTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeTicketLink("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeMissingTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeExternalUrlLink("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeRepoFileLink("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildLinkifiedFragment("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderDetailGlossary("),
	}
	declarationBlocks := []string{
		sliceDeclarationAfter(t, indexHtml, "var inlineTicketTitleMaxLength ="),
		sliceDeclarationAfter(t, indexHtml, "var bodyMentionPattern ="),
		sliceDeclarationAfter(t, indexHtml, "var requestIdByReqSegment ="),
	}

	// The 60-character cut is the REQ's number, so the expectations below are
	// written out rather than recomputed from the shipped constant: shrinking
	// the constant must fail this test, not silently move the assertion with it.
	longTitle := "Make every referenced request identifier in a drawer body carry its own title"
	exactlySixtyTitle := "Keep the timeline forecast honest about ordering and timings"
	unbrokenTitle := strings.Repeat("x", 70)

	javascriptProbe := `
function makeStubElement(tagName) {
  return {
    stubTag: tagName,
    className: "",
    dataset: {},
    childNodes: [],
    textContent: "",
    hidden: false,
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
}
var document = {
  createElement: function (tagName) { return makeStubElement(tagName); },
  createTextNode: function (nodeText) { return { stubTag: "#text", textContent: nodeText, childNodes: [] }; },
  createDocumentFragment: function () { return makeStubElement("#fragment"); }
};
var drawerGlossary = makeStubElement("section");
var requestsById = {
  "REQ-1679": { title: ` + mustMarshalJSONString(t, exactlySixtyTitle) + `, status: "completed" },
  "REQ-1108": { title: "Short one", status: "pending" },
  "REQ-1685": { title: ` + mustMarshalJSONString(t, longTitle) + `, status: "claimed" },
  "UR-001-REQ-042": { title: "First half of an ambiguous pair", status: "pending" },
  "UR-002-REQ-042": { title: "Second half of an ambiguous pair", status: "pending" }
};
var userRequestsById = {
  "UR-074": { title: "Ticket ids should carry their titles" },
  // Two titleless shapes the live tree does not currently hold, so they are
  // fixtures rather than sampled data: a UR synthesized because its input.md was
  // never found (linkRequestsToUserRequests leaves Title empty by design), and a
  // UR whose input.md exists but names no title.
  "UR-900": { inputFilePresent: false },
  "UR-901": { title: "", inputFilePresent: true }
};
var repoFileMentionExists = {};
var liveFileApiAvailable = false;
` + strings.Join(functionBlocks, "\n") + "\n" + strings.Join(declarationBlocks, "\n") + `

function collectNodeText(node) {
  if (node.childNodes && node.childNodes.length > 0) {
    return node.childNodes.map(collectNodeText).join("");
  }
  return node.textContent || "";
}
function describeNode(node) {
  return {
    tag: node.stubTag,
    className: node.className || "",
    text: collectNodeText(node),
    title: node.title || "",
    detailKind: (node.dataset && node.dataset.detailKind) || "",
    detailId: (node.dataset && node.dataset.detailId) || "",
    childClassNames: (node.childNodes || []).map(function (childNode) { return childNode.className || ""; })
  };
}
function describeFragment(fragment) {
  return fragment === null ? [] : fragment.childNodes.map(describeNode);
}

var mentionRenderState = { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] };
// A backticked id comes first on purpose: it must not consume the inline
// expansion slot, and it must still earn its glossary line.
var codeSpanFragment = buildLinkifiedFragment("REQ-1108", true, false, mentionRenderState);
var proseFragment = buildLinkifiedFragment(
  "Read REQ-1679 lessons, REQ-1108 again, UR-074 for context, plus REQ-9999 and REQ-042.",
  false,
  false,
  mentionRenderState
);
var repeatFragment = buildLinkifiedFragment("the REQ-1679 note and REQ-1108 once more", false, false, mentionRenderState);
var ambiguousOnlyFragment = buildLinkifiedFragment("see REQ-042 today", false, false, { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] });

// A UR whose input.md was never found is synthesized with no Title
// (linkRequestsToUserRequests). It is a supported board state, so it must not
// fall back to the bare id the whole feature exists to remove.
var synthesizedState = { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] };
var synthesizedFragment = buildLinkifiedFragment("see UR-900 and UR-901 for that", false, false, synthesizedState);
var synthesizedGlossary = synthesizedState.glossaryEntries;

// The two code contexts drive DIFFERENT suppressions, so each is probed alone
// against the same missing id. An inline code span is still a reference and
// flags; a fenced block prints templates and worked examples and must not.
// REQ-042 rides along in both to prove the ambiguity guard is independent of
// context — ambiguous is never flagged, fenced or not.
var inlineCodeMissingFragment = buildLinkifiedFragment(
  "depends_on: [REQ-9999, REQ-042]",
  true,
  false,
  { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] }
);
var fencedMissingFragment = buildLinkifiedFragment(
  "depends_on: [REQ-9999, REQ-042]",
  true,
  true,
  { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] }
);
var proseMissingFragment = buildLinkifiedFragment(
  "see REQ-9999 and REQ-042",
  false,
  false,
  { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] }
);

renderDetailGlossary(mentionRenderState.glossaryEntries);
var glossaryList = drawerGlossary.childNodes.filter(function (childNode) { return childNode.stubTag === "dl"; })[0];
var glossaryRows = [];
if (glossaryList) {
  for (var rowIndex = 0; rowIndex + 1 < glossaryList.childNodes.length; rowIndex += 2) {
    var termNode = glossaryList.childNodes[rowIndex];
    var definitionNode = glossaryList.childNodes[rowIndex + 1];
    glossaryRows.push({
      termTag: termNode.stubTag,
      identifier: collectNodeText(termNode),
      detailKind: termNode.childNodes[0].dataset.detailKind,
      title: collectNodeText(definitionNode.childNodes[0]),
      status: collectNodeText(definitionNode.childNodes[1])
    });
  }
}
var glossaryHidden = drawerGlossary.hidden;

drawerGlossary = makeStubElement("section");
renderDetailGlossary([]);

process.stdout.write(JSON.stringify({
  shortTitles: [
    shortTicketTitle(` + mustMarshalJSONString(t, longTitle) + `),
    shortTicketTitle(` + mustMarshalJSONString(t, exactlySixtyTitle) + `),
    shortTicketTitle(` + mustMarshalJSONString(t, exactlySixtyTitle) + ` + "X"),
    shortTicketTitle(` + mustMarshalJSONString(t, unbrokenTitle) + `),
    shortTicketTitle("")
  ],
  codeSpanFragment: describeFragment(codeSpanFragment),
  inlineCodeMissingFragment: describeFragment(inlineCodeMissingFragment),
  fencedMissingFragment: describeFragment(fencedMissingFragment),
  proseMissingFragment: describeFragment(proseMissingFragment),
  synthesizedFragment: describeFragment(synthesizedFragment),
  synthesizedGlossaryTitles: synthesizedGlossary.map(function (entry) { return { id: entry.id, title: entry.title }; }),
  proseFragment: describeFragment(proseFragment),
  repeatFragment: describeFragment(repeatFragment),
  ambiguousOnlyLinked: ambiguousOnlyFragment !== null,
  metaRowLink: describeNode(makeTicketLink("req", "REQ-1685", null, true)),
  glossary: glossaryRows,
  glossaryHidden: glossaryHidden,
  emptyGlossaryHidden: drawerGlossary.hidden,
  emptyGlossaryChildCount: drawerGlossary.childNodes.length
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "ticket mention titles", javascriptProbe)
	var probeResult ticketMentionProbeResult
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode ticket mention behavior: %v (output %q)", decodeError, probeOutput)
	}

	wantShortTitles := []string{
		"Make every referenced request identifier in a drawer body…",
		exactlySixtyTitle,
		"Keep the timeline forecast honest about ordering and…",
		strings.Repeat("x", 60) + "…",
		"",
	}
	if !reflect.DeepEqual(probeResult.ShortTitles, wantShortTitles) {
		t.Errorf("shortTicketTitle results = %#v, want %#v", probeResult.ShortTitles, wantShortTitles)
	}

	// A backticked id keeps the bare mono link: no title span, no tooltip.
	if len(probeResult.CodeSpanFragment) != 1 {
		t.Fatalf("code-span fragment = %#v, want one node", probeResult.CodeSpanFragment)
	}
	codeSpanLink := probeResult.CodeSpanFragment[0]
	if codeSpanLink.Tag != "a" || codeSpanLink.ClassName != "ticket-link" || codeSpanLink.Text != "REQ-1108" {
		t.Errorf("code-span mention = %#v, want a bare ticket-link reading REQ-1108", codeSpanLink)
	}
	if codeSpanLink.Title != "" || len(codeSpanLink.ChildClassNames) != 0 {
		t.Errorf("code-span mention gained prose: title=%q children=%#v", codeSpanLink.Title, codeSpanLink.ChildClassNames)
	}

	proseLinks := map[string]ticketMentionNodeProbe{}
	for _, proseNode := range probeResult.ProseFragment {
		if proseNode.DetailId != "" {
			proseLinks[proseNode.DetailId] = proseNode
		}
	}
	firstRequestMention, hasFirstRequestMention := proseLinks["REQ-1679"]
	if !hasFirstRequestMention {
		t.Fatalf("prose fragment has no REQ-1679 link: %#v", probeResult.ProseFragment)
	}
	if !reflect.DeepEqual(firstRequestMention.ChildClassNames, []string{"ticket-link-id", "", "ticket-link-title"}) {
		t.Errorf("first REQ-1679 mention children = %#v, want id + separator + title", firstRequestMention.ChildClassNames)
	}
	if firstRequestMention.Text != "REQ-1679 "+exactlySixtyTitle {
		t.Errorf("first REQ-1679 mention text = %q, want the id and its title", firstRequestMention.Text)
	}
	if firstRequestMention.Title != exactlySixtyTitle {
		t.Errorf("first REQ-1679 mention tooltip = %q, want the untruncated title", firstRequestMention.Title)
	}
	// The code span already resolved REQ-1108; the first PROSE mention is still
	// the one that expands.
	firstCodeThenProseMention := proseLinks["REQ-1108"]
	if !reflect.DeepEqual(firstCodeThenProseMention.ChildClassNames, []string{"ticket-link-id", "", "ticket-link-title"}) {
		t.Errorf("REQ-1108 prose mention children = %#v, want the code span not to have spent the expansion",
			firstCodeThenProseMention.ChildClassNames)
	}
	userRequestMention := proseLinks["UR-074"]
	if userRequestMention.DetailKind != "ur" || userRequestMention.Text != "UR-074 Ticket ids should carry their titles" {
		t.Errorf("UR-074 mention = %#v, want an expanded user-request link", userRequestMention)
	}

	var brokenNodes []ticketMentionNodeProbe
	var proseText string
	for _, proseNode := range probeResult.ProseFragment {
		proseText += proseNode.Text
		if proseNode.ClassName == "ticket-missing" {
			brokenNodes = append(brokenNodes, proseNode)
		}
	}
	if len(brokenNodes) != 1 || brokenNodes[0].Tag != "span" || brokenNodes[0].Text != "REQ-9999" {
		t.Errorf("unresolved id nodes = %#v, want one non-link ticket-missing span for REQ-9999", brokenNodes)
	} else if brokenNodes[0].Title != "Not found in this queue" {
		t.Errorf("unresolved id tooltip = %q, want the not-found tooltip", brokenNodes[0].Title)
	}
	// Ambiguous is not missing: the board knows two records and refuses to pick.
	if _, ambiguousWasLinked := proseLinks["UR-001-REQ-042"]; ambiguousWasLinked {
		t.Error("an ambiguous REQ segment was linked — the never-guess rule broke")
	}
	if !strings.Contains(proseText, "REQ-042.") {
		t.Errorf("prose text = %q, want the ambiguous segment left as plain prose", proseText)
	}
	if probeResult.AmbiguousOnlyLinked {
		t.Error("a text run whose only mention is ambiguous was rewritten; it must be left untouched")
	}

	// A titleless record is not a missing one, and both shapes are supported:
	// linkRequestsToUserRequests SYNTHESIZES a UserRequestTicket with no Title
	// whenever a REQ names a UR whose input.md was not found, and a real UR can
	// simply name no title. Before this case the empty title fell through
	// makeTicketLink's !fullTitle branch to the bare id — reintroducing the exact
	// cryptic number the feature exists to remove, on the one kind of record that
	// cannot explain itself. Each now says WHY it has no title, marked a fallback
	// so it renders as a description rather than as the record's own words.
	//
	// Both are fixtures, deliberately: this repo's tree ships zero synthesized URs
	// and one titleless one, so sampling live data would leave the branch untested
	// on the board that matters and silently vacuous on any other.
	synthesizedLinks := map[string]ticketMentionNodeProbe{}
	for _, fragmentNode := range probeResult.SynthesizedFragment {
		if fragmentNode.DetailId != "" {
			synthesizedLinks[fragmentNode.DetailId] = fragmentNode
		}
	}
	for _, titlelessCase := range []struct{ detailId, wantPhrase, why string }{
		{"UR-900", "no input.md", "a UR synthesized from REQ pointers"},
		{"UR-901", "untitled", "a UR that exists but names no title"},
	} {
		titlelessLink, wasLinked := synthesizedLinks[titlelessCase.detailId]
		if !wasLinked {
			t.Errorf("%s (%s) produced no link at all", titlelessCase.detailId, titlelessCase.why)
			continue
		}
		if !reflect.DeepEqual(titlelessLink.ChildClassNames, []string{"ticket-link-id", "", "ticket-link-title is-fallback"}) {
			t.Errorf("%s (%s) children = %#v, want it expanded as a marked fallback rather than left a bare id",
				titlelessCase.detailId, titlelessCase.why, titlelessLink.ChildClassNames)
		}
		if !strings.Contains(titlelessLink.Text, titlelessCase.wantPhrase) {
			t.Errorf("%s (%s) link = %q, want it to say why it has no title (%q)",
				titlelessCase.detailId, titlelessCase.why, titlelessLink.Text, titlelessCase.wantPhrase)
		}
	}
	glossaryFallbacks := map[string]ticketFallbackTitleProbe{}
	for _, glossaryRow := range probeResult.SynthesizedGlossary {
		glossaryFallbacks[glossaryRow.Id] = glossaryRow.Title
	}
	for _, detailId := range []string{"UR-900", "UR-901"} {
		fallbackTitle, wasGlossed := glossaryFallbacks[detailId]
		if !wasGlossed {
			t.Errorf("%s earned no glossary entry", detailId)
			continue
		}
		if !fallbackTitle.IsFallback {
			t.Errorf("%s's substitute title is not marked a fallback — it would render dressed as the record's own title", detailId)
		}
	}

	// Where the broken-reference flag fires, by context. The three cases share one
	// missing id (REQ-9999) and one ambiguous id (REQ-042) so the only variable is
	// the context, and each pins a distinct rule:
	//
	//   prose        → flagged. An id written in prose is a real reference.
	//   inline code  → flagged. A backticked id in prose is still a reference;
	//                  REQ bodies conventionally backtick the ids they cite.
	//   fenced block → NOT flagged. A fence prints templates and worked examples
	//                  ("id: REQ-021"), which point at nothing and must not be
	//                  asserted missing. Without this, 115 of the 397 flags on
	//                  this repo's own board are illustrations, not typos.
	//
	// REQ-042 must be absent from all three: ambiguous is not missing, in any
	// context. Deleting the insideFencedBlock guard makes fencedMissingCount 1;
	// widening it to any code context makes inlineCodeMissingCount 0. Both fail.
	countMissingSpans := func(fragmentNodes []ticketMentionNodeProbe) (missingCount int, sawAmbiguous bool) {
		for _, fragmentNode := range fragmentNodes {
			if fragmentNode.ClassName == "ticket-missing" {
				missingCount++
				if fragmentNode.Text == "REQ-042" {
					sawAmbiguous = true
				}
			}
		}
		return missingCount, sawAmbiguous
	}
	for _, flagCase := range []struct {
		contextName      string
		fragmentNodes    []ticketMentionNodeProbe
		wantMissingCount int
	}{
		{"prose", probeResult.ProseMissing, 1},
		{"inline code span", probeResult.InlineCodeMissing, 1},
		{"fenced code block", probeResult.FencedMissing, 0},
	} {
		gotMissingCount, sawAmbiguous := countMissingSpans(flagCase.fragmentNodes)
		if gotMissingCount != flagCase.wantMissingCount {
			t.Errorf("REQ-9999 in a %s: %d ticket-missing spans, want %d",
				flagCase.contextName, gotMissingCount, flagCase.wantMissingCount)
		}
		if sawAmbiguous {
			t.Errorf("the ambiguous REQ-042 was flagged missing in a %s — ambiguous is not missing",
				flagCase.contextName)
		}
	}

	// A later mention of an already-expanded id stays bare, and so does its tooltip.
	for _, repeatNode := range probeResult.RepeatFragment {
		if repeatNode.DetailId == "" {
			continue
		}
		if len(repeatNode.ChildClassNames) != 0 || repeatNode.Title != "" {
			t.Errorf("repeat mention of %s expanded again: %#v", repeatNode.DetailId, repeatNode)
		}
	}

	// Meta rows are reference lists, not prose: they always carry the title,
	// truncated inline with the full text in the tooltip.
	if !reflect.DeepEqual(probeResult.MetaRowLink.ChildClassNames, []string{"ticket-link-id", "", "ticket-link-title"}) {
		t.Errorf("meta-row link children = %#v, want an always-expanded link", probeResult.MetaRowLink.ChildClassNames)
	}
	if probeResult.MetaRowLink.Title != longTitle {
		t.Errorf("meta-row link tooltip = %q, want the untruncated title", probeResult.MetaRowLink.Title)
	}
	if probeResult.MetaRowLink.Text != "REQ-1685 Make every referenced request identifier in a drawer body…" {
		t.Errorf("meta-row link text = %q, want the id and the truncated title", probeResult.MetaRowLink.Text)
	}

	wantGlossary := []ticketGlossaryRowProbe{
		{TermTag: "dt", Identifier: "REQ-1108", DetailKind: "req", Title: "Short one", Status: "pending"},
		{TermTag: "dt", Identifier: "REQ-1679", DetailKind: "req", Title: exactlySixtyTitle, Status: "completed"},
		{TermTag: "dt", Identifier: "UR-074", DetailKind: "ur", Title: "Ticket ids should carry their titles", Status: "user request"},
	}
	if !reflect.DeepEqual(probeResult.Glossary, wantGlossary) {
		t.Errorf("glossary = %#v, want one line per resolved id in first-mention order %#v", probeResult.Glossary, wantGlossary)
	}
	if probeResult.GlossaryHidden {
		t.Error("the glossary stayed hidden with entries to show")
	}
	if !probeResult.EmptyGlossaryHidden || probeResult.EmptyGlossaryChildCount != 0 {
		t.Errorf("a body that cited nothing left a glossary: hidden=%v children=%d",
			probeResult.EmptyGlossaryHidden, probeResult.EmptyGlossaryChildCount)
	}
}

func TestJavaScriptBehaviorRecentlyDoneWindowRefreshesVisibleLens(t *testing.T) {
	indexHtml := generateLiveSite(t)
	const wiringToken = `applyRecentWindowSelection(parseInt(button.getAttribute("data-window-hours"), 10))`
	if !strings.Contains(indexHtml, wiringToken) {
		t.Fatalf("recent-window click handler is not wired to the transition helper: %q missing", wiringToken)
	}

	recentWindowFunction := sliceBalancedBlockAfter(t, indexHtml, "function applyRecentWindowSelection(")
	javascriptProbe := `
var viewState = { windowHours: 24, view: "board", lens: "user-request" };
var renderedOnce = { userRequestLens: true };
var selectedWindow = "";
var columnRenderCount = 0;
var lensRenderCount = 0;
function setActiveButton(selector, attributeName, attributeValue) { selectedWindow = attributeValue; }
function renderColumns() { columnRenderCount += 1; }
function renderUserRequestLens() { lensRenderCount += 1; }
` + recentWindowFunction + `
applyRecentWindowSelection(168);
var visibleLensState = {
  windowHours: viewState.windowHours,
  selectedWindow: selectedWindow,
  columnRenderCount: columnRenderCount,
  lensRenderCount: lensRenderCount,
  lensFresh: renderedOnce.userRequestLens
};
viewState.lens = "columns";
applyRecentWindowSelection(48);
process.stdout.write(JSON.stringify({
  visibleLensState: visibleLensState,
  hiddenLensState: {
    windowHours: viewState.windowHours,
    columnRenderCount: columnRenderCount,
    lensRenderCount: lensRenderCount,
    lensFresh: renderedOnce.userRequestLens
  }
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "recent-window transition", javascriptProbe)
	var result struct {
		VisibleLensState struct {
			WindowHours       int    `json:"windowHours"`
			SelectedWindow    string `json:"selectedWindow"`
			ColumnRenderCount int    `json:"columnRenderCount"`
			LensRenderCount   int    `json:"lensRenderCount"`
			LensFresh         bool   `json:"lensFresh"`
		} `json:"visibleLensState"`
		HiddenLensState struct {
			WindowHours       int  `json:"windowHours"`
			ColumnRenderCount int  `json:"columnRenderCount"`
			LensRenderCount   int  `json:"lensRenderCount"`
			LensFresh         bool `json:"lensFresh"`
		} `json:"hiddenLensState"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode recent-window transition: %v (output %q)", decodeError, probeOutput)
	}
	visibleState := result.VisibleLensState
	if visibleState.WindowHours != 168 || visibleState.SelectedWindow != "168" || visibleState.ColumnRenderCount != 1 || visibleState.LensRenderCount != 1 || !visibleState.LensFresh {
		t.Fatalf("visible by-UR transition = %#v, want selected window 168 with both lenses refreshed", visibleState)
	}
	hiddenState := result.HiddenLensState
	if hiddenState.WindowHours != 48 || hiddenState.ColumnRenderCount != 2 || hiddenState.LensRenderCount != 1 || hiddenState.LensFresh {
		t.Fatalf("hidden by-UR transition = %#v, want columns refreshed and by-UR marked stale", hiddenState)
	}
}

func TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller(t *testing.T) {
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
    "REQ-501": { status: "completed", title: "old work" },
    "REQ-502": { status: "completed", title: "recent work" }
  },
  userRequests: {
    "UR-301": { requestIds: ["REQ-501"], title: "old request", inputFilePresent: true },
    "UR-302": { requestIds: ["REQ-502"], title: "recent request", inputFilePresent: true }
  },
  userRequestOrder: ["UR-301", "UR-302"],
  calendar: [
    { id: "REQ-502", status: "completed", entryTime: "2026-08-15T06:00:00Z" },
    { id: "REQ-501", status: "completed", entryTime: "2026-08-13T06:00:00Z" }
  ]
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "active" };
// The always-open reading of this lens; the folded one has its own probe.
var userRequestCardsFolded = false;
function makeNode() {
  return {
    childNodes: [],
    dataset: {},
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
}
var userRequestLensNode = makeNode();
var document = {
  getElementById: function (nodeId) { return nodeId === "user-request-lens" ? userRequestLensNode : null; },
  createElement: function () { return makeNode(); }
};
function makeRequestCard(requestId) { return { requestId: requestId }; }
` + strings.Join(functionBlocks, "\n") + `
renderUserRequestLens();
var renderedUserRequestIds = userRequestLensNode.childNodes
  .filter(function (node) { return node.className === "ur-group"; })
  .map(function (groupNode) { return groupNode.childNodes[0].dataset.detailId; });
var scopeNotes = userRequestLensNode.childNodes
  .filter(function (node) { return node.className === "ur-lens-hidden-note"; })
  .map(function (node) { return node.textContent; });
process.stdout.write(JSON.stringify({ renderedUserRequestIds: renderedUserRequestIds, scopeNotes: scopeNotes }));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "by-UR caller", javascriptProbe)
	var result struct {
		RenderedUserRequestIds []string `json:"renderedUserRequestIds"`
		ScopeNotes             []string `json:"scopeNotes"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode assembled-client by-UR caller output: %v (output %q)", decodeError, probeOutput)
	}
	if len(result.RenderedUserRequestIds) != 1 || result.RenderedUserRequestIds[0] != "UR-302" {
		t.Fatalf("Active by-UR caller rendered %#v, want only recent terminal UR-302", result.RenderedUserRequestIds)
	}
	if len(result.ScopeNotes) != 1 || !strings.Contains(result.ScopeNotes[0], "1 UR with no open work or activity in the last 24 hours") {
		t.Fatalf("Active by-UR caller scope notes = %#v, want one old hidden UR and the selected window", result.ScopeNotes)
	}
}

func TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours(t *testing.T) {
	indexHtml := generateLiveSite(t)
	annotationCases := durationsSlowestDayAnnotationCaseList()

	probeCases, encodeError := json.Marshal(durationsSlowestDayAnnotationProbeCases(annotationCases))
	if encodeError != nil {
		t.Fatalf("encode annotation probe cases: %v", encodeError)
	}
	annotationSource := sliceBalancedBlockAfter(t, indexHtml, "function drawDurationsSlowestDayAnnotation(")
	assertDurationsAnnotationBaselineIgnoresItsInputs(t, annotationSource)

	javascriptProbe := fmt.Sprintf("var DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y = %v;\n",
		durationsRendererConstant(t, "DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y")) +
		"var drawnNodes = [];\n" +
		"function makeDurationsSvgNode(svg, name, attributes, textContent) { drawnNodes.push({ name: name, attributes: attributes, text: textContent }); }\n" +
		annotationSource + `
var probeCases = ` + string(probeCases) + `;
probeCases.forEach(function (probeCase) {
  drawDurationsSlowestDayAnnotation(null, { medianMinutes: probeCase.medianMinutes }, probeCase.dayCentreX);
});
process.stdout.write(JSON.stringify(drawnNodes.map(function (node) {
  return { y: Number(node.attributes.y), x: Number(node.attributes.x), anchor: node.attributes["text-anchor"], text: node.text };
})));`

	probeOutput := runJavaScriptBehaviorProbe(t, "durations slowest-day annotation", javascriptProbe)
	var drawnAnnotations []struct {
		Y      float64 `json:"y"`
		X      float64 `json:"x"`
		Anchor string  `json:"anchor"`
		Text   string  `json:"text"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawnAnnotations); decodeError != nil {
		t.Fatalf("decode slowest-day annotation behavior: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}
	if len(drawnAnnotations) != len(annotationCases) {
		t.Fatalf("the renderer drew %d annotations for %d cases — it must draw exactly one per slowest day",
			len(drawnAnnotations), len(annotationCases))
	}

	// The annotation's face is measured independently from the renderer geometry.
	annotationAscent := durationsMeasuredMarkLabelAscentUnits
	annotationDescent := math.Max(durationsLabelTextDescentUnits, durationsMeasuredMarkLabelDescentUnits)
	medianBaseline := durationsRendererConstant(t, "DURATIONS_MEDIAN_BOTTOM")

	// The annotation's neighbours in the strip it now occupies. Panel B's "0"
	// tick is the one a render caught and no arithmetic would have: it sits in
	// the y-axis gutter, so it is only ever in the annotation's way when the
	// slowest day is the leftmost — the same luck-of-x that hid the original
	// defect. The ticks for 15/30/45 need no case of their own: their baselines
	// are above DURATIONS_MEDIAN_BOTTOM, which the baseline check below covers.
	neighbourBoxes := []struct {
		neighbourName string
		baseline      float64
		ascent        float64
		descent       float64
	}{
		{
			"panel B's title",
			durationsRendererConstant(t, "DURATIONS_MEDIAN_TITLE_Y"),
			durationsMeasuredAxisTitleAscentUnits,
			durationsMeasuredAxisTitleDescentUnits,
		},
		{
			"panel B's \"0\" axis tick",
			durationsRendererConstant(t, "DURATIONS_MEDIAN_BOTTOM") + durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP"),
			annotationAscent,
			annotationDescent,
		},
		{
			"panel C's title",
			durationsRendererConstant(t, "DURATIONS_COUNT_TITLE_Y"),
			durationsMeasuredAxisTitleAscentUnits,
			durationsMeasuredAxisTitleDescentUnits,
		},
	}
	for _, neighbour := range neighbourBoxes {
		neighbourTop := neighbour.baseline - neighbour.ascent
		neighbourBottom := neighbour.baseline + neighbour.descent
		for caseIndex, probeCase := range annotationCases {
			drawn := drawnAnnotations[caseIndex]
			annotationTop := drawn.Y - annotationAscent
			annotationBottom := drawn.Y + annotationDescent
			if annotationBottom >= neighbourTop && annotationTop <= neighbourBottom {
				t.Fatalf("%s: the annotation's text box [%.2f, %.2f] intersects %s's box [%.2f, %.2f] — the two overprint wherever their x ranges meet, and x follows whichever day is slowest",
					probeCase.caseName, annotationTop, annotationBottom, neighbour.neighbourName, neighbourTop, neighbourBottom)
			}
		}
	}
	for caseIndex, probeCase := range annotationCases {
		drawn := drawnAnnotations[caseIndex]
		if drawn.Y-annotationAscent <= medianBaseline {
			t.Fatalf("%s: the annotation's text box starts at %.2f, above panel B's baseline at %.2f — inside the plot it overprints the bars, which are 4 units wide and shoulder to shoulder on a dense board",
				probeCase.caseName, drawn.Y-annotationAscent, medianBaseline)
		}
		if drawn.Y != drawnAnnotations[0].Y {
			t.Fatalf("%s: the annotation's baseline is %.2f but %s put it at %.2f — a baseline that moves with the day's position or its bar's height is a clearance that holds only for the days this repository happens to have",
				probeCase.caseName, drawn.Y, annotationCases[0].caseName, drawnAnnotations[0].Y)
		}
		if drawn.X != probeCase.dayCentreX || drawn.Anchor != "middle" {
			t.Fatalf("%s: the annotation was drawn at x=%.2f anchored %q, want x=%.2f anchored \"middle\" — it must stay centred on the day it describes",
				probeCase.caseName, drawn.X, drawn.Anchor, probeCase.dayCentreX)
		}
		if wantText := fmt.Sprintf("%.0f min", math.Round(probeCase.medianMinutes)); drawn.Text != wantText {
			t.Fatalf("%s: the annotation reads %q, want %q — moving it must not cost it the value it exists to state",
				probeCase.caseName, drawn.Text, wantText)
		}
	}

	// The strip's fourth occupant is the month rule, and unlike the other three
	// it cannot be cleared: .durations-month-line spans DURATIONS_MAIN_TOP to
	// DURATIONS_COUNT_BOTTOM, so it crosses EVERY baseline the annotation could
	// legally take, and it crosses panel A's reversed-band labels the same way.
	// The crossing is accepted, not overlooked — on a fixture whose slowest day
	// falls on a month boundary the rule passes between the "9" and the " min".
	// What makes it acceptable is that it is a one-unit soft rule, so that is
	// what gets asserted in place of a clearance: if the month rule ever grows
	// wide or firm, this fires and the acceptance has to be re-argued.
	annotationTop := drawnAnnotations[0].Y - annotationAscent
	annotationBottom := drawnAnnotations[0].Y + annotationDescent
	monthRuleTop := durationsRendererConstant(t, "DURATIONS_MAIN_TOP")
	monthRuleBottom := durationsRendererConstant(t, "DURATIONS_COUNT_BOTTOM")
	if annotationTop <= monthRuleTop || annotationBottom >= monthRuleBottom {
		t.Fatalf("the annotation's text box [%.2f, %.2f] is no longer inside the month rule's span [%.2f, %.2f] — the crossing this test ACCEPTS has become avoidable, so it belongs in the clearance list above instead",
			annotationTop, annotationBottom, monthRuleTop, monthRuleBottom)
	}
	if strokeWidth := durationsStyleDeclaration(t, ".durations-month-line", "stroke-width"); strokeWidth != "1" {
		t.Fatalf("the month rule's stroke-width is %q, not \"1\" — it is allowed to cross the slowest-day annotation only because it is a hairline",
			strokeWidth)
	}
	if strokeColour := durationsStyleDeclaration(t, ".durations-month-line", "stroke"); strokeColour != "var(--line-soft)" {
		t.Fatalf("the month rule is stroked %q, not the soft line token — it is allowed to cross the slowest-day annotation only because it is soft",
			strokeColour)
	}
}

func TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	marginLeft := durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")
	plotRight := durationsRendererConstant(t, "DURATIONS_VIEW_WIDTH") -
		durationsRendererConstant(t, "DURATIONS_MARGIN_RIGHT")

	for _, dayCount := range []int{1, 2, 14, 400} {
		dayCount := dayCount
		t.Run(fmt.Sprintf("%d-active-days", dayCount), func(t *testing.T) {
			fixtureTickets := durationsDayCountFixtureTickets(dayCount)
			aggregate := buildDurationAggregate(fixtureTickets)
			fixtureBoard := &Board{
				GeneratedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC),
				AllRequests: fixtureTickets,
			}
			generatedData, buildError := buildGeneratedBoardData(fixtureBoard)
			if buildError != nil {
				t.Fatalf("buildGeneratedBoardData: %v", buildError)
			}
			durationsJson, encodeError := json.Marshal(generatedData.Durations)
			if encodeError != nil {
				t.Fatalf("encode durations payload: %v", encodeError)
			}

			javascriptProbe := durationsRenderDomStubPreamble +
				"var boardData = { durations: " + string(durationsJson) + " };\n" +
				string(rendererFragment) +
				"setDurationsWindow(\"all\");\n" +
				durationsRenderProbeDriver
			probeOutput := runJavaScriptBehaviorProbe(t,
				fmt.Sprintf("durations day buckets (%d active days)", dayCount), javascriptProbe)

			var drawn struct {
				Bars []struct {
					Class string  `json:"class"`
					X     float64 `json:"x"`
					Width float64 `json:"width"`
				} `json:"bars"`
				AnnotationXs []float64 `json:"annotationXs"`
				MarkCxs      []float64 `json:"markCxs"`
				FirstMarkCxs []float64 `json:"firstMarkCxs"`
			}
			if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
				t.Fatalf("decode drawn durations geometry: %v (output starts %q)",
					decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
			}

			// (1) Every Panel B and C bar inside the plot area. 0.05 covers the
			// renderer's toFixed(1) rounding and nothing else.
			panelBBarCount := 0
			for _, bar := range drawn.Bars {
				if !strings.Contains(bar.Class, "durations-bar-count") &&
					!strings.Contains(bar.Class, "durations-bar-over-ceiling") {
					panelBBarCount++
				}
				if bar.X < marginLeft-0.05 || bar.X+bar.Width > plotRight+0.05 {
					t.Errorf("%q bar spans x %.1f–%.1f, outside the plot area [%.0f, %.0f]",
						bar.Class, bar.X, bar.X+bar.Width, marginLeft, plotRight)
				}
			}

			// (2) One Panel B bar per day with a median — the guard against a
			// render that went off the rails and drew nothing to check.
			medianDayCount := 0
			for _, day := range aggregate.Days {
				if day.HasMedian {
					medianDayCount++
				}
			}
			if panelBBarCount != medianDayCount {
				t.Errorf("drew %d Panel B bars for %d days with a median", panelBBarCount, medianDayCount)
			}

			// (3) Exactly one slowest-day annotation, anchored inside the plot —
			// it exists to state a value a clipped bar cannot, and cannot do
			// that from off-screen.
			if len(drawn.AnnotationXs) != 1 {
				t.Fatalf("drew %d slowest-day annotations, want exactly 1", len(drawn.AnnotationXs))
			}
			if annotationX := drawn.AnnotationXs[0]; annotationX < marginLeft-0.05 || annotationX > plotRight+0.05 {
				t.Errorf("slowest-day annotation anchored at x=%.1f, outside the plot area [%.0f, %.0f]",
					annotationX, marginLeft, plotRight)
			}

			// (4) Every Panel A mark stays inside its own UTC day slot. This pins
			// the shared whole-day axis domain while allowing the intentional
			// deterministic jitter within that slot.
			rangeStart, rangeEnd, hasRange := durationLabelTimeRange(aggregate.Samples)
			if !hasRange {
				t.Fatal("fixture produced no label time range")
			}
			if len(drawn.MarkCxs) != len(aggregate.Samples) {
				t.Fatalf("drew %d marks for %d samples", len(drawn.MarkCxs), len(aggregate.Samples))
			}
			for sampleIndex, sample := range aggregate.Samples {
				dayStart := sample.CompletionTime.UTC().Truncate(24 * time.Hour)
				dayLeft := marginLeft + durationLabelPlotX(dayStart, rangeStart, rangeEnd)
				dayRight := marginLeft + durationLabelPlotX(dayStart.Add(24*time.Hour), rangeStart, rangeEnd)
				if drawn.MarkCxs[sampleIndex] < dayLeft-0.06 || drawn.MarkCxs[sampleIndex] > dayRight+0.06 {
					t.Errorf("%s mark drawn at x=%.2f outside its UTC-day slot [%.2f, %.2f]",
						sample.RequestId, drawn.MarkCxs[sampleIndex], dayLeft, dayRight)
				}
			}

			// (5) The same payload renders to the same jittered coordinates, and
			// a busy day uses more than one x instead of restacking every mark.
			if len(drawn.FirstMarkCxs) != len(drawn.MarkCxs) {
				t.Fatalf("first render drew %d marks and second drew %d", len(drawn.FirstMarkCxs), len(drawn.MarkCxs))
			}
			busyDaySpread := 0.0
			for sampleIndex := range drawn.MarkCxs {
				if math.Abs(drawn.FirstMarkCxs[sampleIndex]-drawn.MarkCxs[sampleIndex]) > 0.001 {
					t.Errorf("mark %d moved from x=%.2f to x=%.2f across identical renders",
						sampleIndex, drawn.FirstMarkCxs[sampleIndex], drawn.MarkCxs[sampleIndex])
				}
				if sampleIndex%2 == 1 {
					busyDaySpread = math.Max(busyDaySpread, math.Abs(drawn.MarkCxs[sampleIndex]-drawn.MarkCxs[sampleIndex-1]))
				}
			}
			if busyDaySpread < 0.1 {
				t.Errorf("busy-day marks have only %.2f units of x spread; jitter is degenerate", busyDaySpread)
			}
		})
	}
}

func TestJavaScriptBehaviorDurationsColourChannelsNameAndRecolourPanelA(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	fixtureTickets := make([]*RequestTicket, 0, 16)
	fixtureStart := time.Date(2026, 8, 10, 9, 0, 0, 0, time.UTC)
	for ticketIndex := 0; ticketIndex < 14; ticketIndex++ {
		completedAt := fixtureStart.Add(time.Duration(ticketIndex) * time.Hour)
		ticket := durationTicket(
			fmt.Sprintf("REQ-%03d", ticketIndex+1),
			"B",
			completedAt.Add(-12*time.Minute).Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		)
		ticket.UserRequestId = fmt.Sprintf("UR-%03d", ticketIndex+1)
		ticket.Domain = []string{"frontend", "backend", "testing"}[ticketIndex%3]
		fixtureTickets = append(fixtureTickets, ticket)
	}
	missingTicket := durationTicket("REQ-901", "", fixtureStart.Add(15*time.Hour).Format(time.RFC3339), fixtureStart.Add(15*time.Hour+12*time.Minute).Format(time.RFC3339))
	fixtureTickets = append(fixtureTickets, missingTicket)
	reversedTicket := durationTicket("REQ-902", "B", fixtureStart.Add(17*time.Hour).Format(time.RFC3339), fixtureStart.Add(16*time.Hour).Format(time.RFC3339))
	reversedTicket.UserRequestId = "UR-015"
	reversedTicket.Domain = "frontend"
	fixtureTickets = append(fixtureTickets, reversedTicket)

	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC),
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
function renderedChannel(colourChannel) {
  durationsStubHosts["durations-chart"] = makeStubNode("div");
  durationsStubHosts["durations-colour-legend"] = makeStubNode("p");
  setDurationsColourChannel(colourChannel);
  renderDurationsView();
  var marks = [];
  function walkDrawnNodes(parentNode) {
    (parentNode.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      if (childNode.stubName === "circle" && String(attributes["class"] || "").indexOf("durations-mark") !== -1) {
        marks.push({ fill: attributes.fill || "", class: attributes["class"] || "" });
      }
      walkDrawnNodes(childNode);
    });
  }
  walkDrawnNodes(durationsStubHosts["durations-chart"]);
  return {
    marks: marks,
    legend: durationsStubHosts["durations-colour-legend"].textContent,
    ariaLabel: (durationsStubHosts["durations-chart"].children[0] || { attributes: {} }).attributes["aria-label"] || ""
  };
}
process.stdout.write(JSON.stringify({
  route: renderedChannel("route"),
  userRequest: renderedChannel("user-request"),
  domain: renderedChannel("domain")
}));
`

	javascriptProbe := durationsRenderDomStubPreamble +
		"var boardData = " + string(boardDataJSON) + ";\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "durations colour channels", javascriptProbe)

	type drawnMark struct {
		Fill  string `json:"fill"`
		Class string `json:"class"`
	}
	type renderedChannel struct {
		Marks     []drawnMark `json:"marks"`
		Legend    string      `json:"legend"`
		AriaLabel string      `json:"ariaLabel"`
	}
	var rendered struct {
		Route       renderedChannel `json:"route"`
		UserRequest renderedChannel `json:"userRequest"`
		Domain      renderedChannel `json:"domain"`
	}
	if decodeError := json.Unmarshal(probeOutput, &rendered); decodeError != nil {
		t.Fatalf("decode rendered colour channels: %v (output starts %q)", decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	if len(rendered.Route.Marks) != len(fixtureTickets) {
		t.Fatalf("route render drew %d marks for %d fixture samples", len(rendered.Route.Marks), len(fixtureTickets))
	}
	if rendered.Route.Marks[0].Fill != "var(--route-b)" {
		t.Errorf("route is not the default channel: first ordinary mark fill = %q, want route B", rendered.Route.Marks[0].Fill)
	}
	if rendered.UserRequest.Marks[0].Fill == rendered.UserRequest.Marks[1].Fill {
		t.Errorf("two named URs share fill %q before the categorical palette is exhausted", rendered.UserRequest.Marks[0].Fill)
	}
	if rendered.UserRequest.Marks[12].Fill != rendered.UserRequest.Marks[13].Fill {
		t.Errorf("the thirteenth and fourteenth named UR fills %q / %q differ — the stated Other URs bucket must be shared", rendered.UserRequest.Marks[12].Fill, rendered.UserRequest.Marks[13].Fill)
	}
	if !strings.Contains(rendered.UserRequest.Legend, "UR") || !strings.Contains(rendered.UserRequest.Legend, "Other URs") || !strings.Contains(rendered.UserRequest.Legend, "No UR recorded") {
		t.Errorf("UR legend %q does not name its active channel, overflow, and missing-value rule", rendered.UserRequest.Legend)
	}
	if !strings.Contains(rendered.Domain.Legend, "Domain") || !strings.Contains(rendered.Domain.Legend, "No domain recorded") {
		t.Errorf("domain legend %q does not name its active channel and missing-value rule", rendered.Domain.Legend)
	}
	if rendered.Domain.Marks[0].Fill == rendered.Domain.Marks[1].Fill {
		t.Errorf("different named domains share fill %q before the palette is exhausted", rendered.Domain.Marks[0].Fill)
	}
	missingMarkIndex := len(fixtureTickets) - 2
	if !strings.Contains(rendered.UserRequest.Marks[missingMarkIndex].Class, "unknown") || !strings.Contains(rendered.Domain.Marks[missingMarkIndex].Class, "unknown") {
		t.Errorf("missing UR/domain sample classes = %q / %q, want the visually distinct unknown mark", rendered.UserRequest.Marks[missingMarkIndex].Class, rendered.Domain.Marks[missingMarkIndex].Class)
	}
	criticalMarkIndex := len(fixtureTickets) - 1
	for channelName, channel := range map[string]renderedChannel{"route": rendered.Route, "user-request": rendered.UserRequest, "domain": rendered.Domain} {
		if channel.Marks[criticalMarkIndex].Fill != "var(--durations-critical)" {
			t.Errorf("%s colour channel changed reversed-stamp fill to %q", channelName, channel.Marks[criticalMarkIndex].Fill)
		}
		if !strings.Contains(channel.AriaLabel, "coloured by ") {
			t.Errorf("%s chart aria label %q does not name the active colour channel", channelName, channel.AriaLabel)
		}
	}
}

func TestJavaScriptBehaviorTimelineRowsFollowTheWindow(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function timelineWaitEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineWorkEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowSegments(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowsInWindow(") + `
var nowMs = Date.UTC(2026, 5, 20, 12, 0);
var windowStartMs = Date.UTC(2026, 5, 10);
var windowEndMs = Date.UTC(2026, 5, 17);

var rows = [
  // Wholly before the window.
  { id: "REQ-001", createdTime: "2026-06-01T09:00:00Z", claimedTime: "2026-06-02T09:00:00Z",
    completedTime: "2026-06-03T09:00:00Z", hasWork: true },
  // Wholly inside it.
  { id: "REQ-002", createdTime: "2026-06-12T09:00:00Z", claimedTime: "2026-06-13T09:00:00Z",
    completedTime: "2026-06-14T09:00:00Z", hasWork: true },
  // Started before, finished after: STRADDLES the window and never sits inside it.
  { id: "REQ-003", createdTime: "2026-06-05T09:00:00Z", claimedTime: "2026-06-06T09:00:00Z",
    completedTime: "2026-06-25T09:00:00Z", hasWork: true },
  // Captured before the window and STILL RUNNING at the now-line beyond it.
  { id: "REQ-004", createdTime: "2026-06-02T09:00:00Z", claimedTime: "2026-06-03T09:00:00Z",
    hasWork: true, workOpen: true },
  // Wholly after the window.
  { id: "REQ-005", createdTime: "2026-06-18T09:00:00Z", claimedTime: "2026-06-19T09:00:00Z",
    completedTime: "2026-06-19T18:00:00Z", hasWork: true }
];
// A reversed stamp: claimed BEFORE created. It carries waitMinutes because that
// is the field BOTH the renderer and the segment list branch on — without it the
// fixture is silently a forward row and tests nothing about reversal.
//
// No segment may come out with its start after its end, or the overlap test
// inverts and the row vanishes from every window. A reversed wait satisfies that
// as a point at the created instant, which is also where the break marker goes.
var reversedRow = { id: "REQ-006", createdTime: "2026-06-14T09:00:00Z",
  claimedTime: "2026-06-12T09:00:00Z", completedTime: "2026-06-15T09:00:00Z",
  hasWork: true, anomaly: true, waitMinutes: -2880, workMinutes: 4320 };
// A pending REQ whose only presence in a FUTURE window is its forecast bar. Its
// measured extent is the open wait, which stops at the now-line; the window sits
// entirely after that, so the row reaches it only if the extent follows the
// projected segment the chart actually draws.
var projectedRow = { id: "REQ-007", createdTime: "2026-06-01T09:00:00Z", waitOpen: true };
var projection = { id: "REQ-007", startTime: "2026-06-21T00:00:00Z", endTime: "2026-06-21T06:00:00Z" };
var futureWindowStartMs = Date.UTC(2026, 5, 20, 18, 0);
var futureWindowEndMs = Date.UTC(2026, 5, 22);

function idsInSpan(rowList, projectedById, spanStartMs, spanEndMs) {
  var extents = rowList.map(function (row) {
    return timelineRowSegments(row, nowMs, (projectedById || {})[row.id]);
  });
  return timelineRowsInWindow(rowList, extents, spanStartMs, spanEndMs)
    .map(function (row) { return row.id; });
}

function idsInWindow(rowList, projectedById) {
  return idsInSpan(rowList, projectedById, windowStartMs, windowEndMs);
}

var reversedSegments = timelineRowSegments(reversedRow, nowMs, undefined);

// THE EXCLUSIVE END. Every window this view can build ends on the NEXT window's
// first instant — the end date field renders windowEndMs - 1 to say so, and
// timelineRowsInWindow tests the two sides asymmetrically for the same reason. A
// segment that begins exactly at
// that instant belongs to the next window, and drawSegment puts its floored
// rectangle at the right edge where it is clipped: listed, nothing drawn.
var edgeWindowStartMs = Date.UTC(2026, 5, 10);
var edgeWindowEndMs = Date.UTC(2026, 5, 17);
var startsAtWindowEndRow = {
  id: "REQ-008",
  createdTime: new Date(edgeWindowEndMs).toISOString(),
  claimedTime: new Date(edgeWindowEndMs + 3600000).toISOString(),
  completedTime: new Date(edgeWindowEndMs + 7200000).toISOString(),
  hasWork: true, waitMinutes: 60, workMinutes: 60
};
// The control, one millisecond earlier: genuinely inside, and must stay listed.
var startsJustInsideRow = {
  id: "REQ-009",
  createdTime: new Date(edgeWindowEndMs - 1).toISOString(),
  claimedTime: new Date(edgeWindowEndMs + 3600000).toISOString(),
  completedTime: new Date(edgeWindowEndMs + 7200000).toISOString(),
  hasWork: true, waitMinutes: 60, workMinutes: 60
};
// The symmetric case at the other end, which must NOT change: windowStartMs is
// inclusive, and a span ending exactly there draws a floored mark at x=0 that
// the reader can see.
var endsAtWindowStartRow = {
  id: "REQ-010",
  createdTime: new Date(edgeWindowStartMs - 7200000).toISOString(),
  claimedTime: new Date(edgeWindowStartMs).toISOString(),
  hasWork: false, waitMinutes: 120
};

// A REVERSED WAIT drawn as what the renderer actually draws. renderVisibleRows
// puts a 6px break marker at the CREATED instant for a reversed wait and at the
// CLAIMED instant for reversed work — it does not draw a bar across the min/max
// interval. Modelling the row as that interval is the forecast-gap defect again
// in another costume: created 14 Jun, claimed 12 Jun, completed 12 Jun 06:00
// puts the hull across 13 June while both drawn marks sit outside it.
var reversedHullRow = {
  id: "REQ-011",
  createdTime: "2026-06-14T12:00:00Z",
  claimedTime: "2026-06-12T00:00:00Z",
  completedTime: "2026-06-12T06:00:00Z",
  hasWork: true, waitMinutes: -3600, workMinutes: 360, anomaly: true
};
var reversedHullGapStartMs = Date.UTC(2026, 5, 13);
var reversedHullGapEndMs = Date.UTC(2026, 5, 13, 23, 59);
// The control: a window over the break marker itself must still list it.
var reversedHullMarkerStartMs = Date.UTC(2026, 5, 14, 6, 0);
var reversedHullMarkerEndMs = Date.UTC(2026, 5, 14, 18, 0);

// THE GAP. A pending REQ draws two disjoint marks: the open wait ending at the
// now-line, and the forecast bar starting after in-flight work finishes. A hull
// over both spans the gap between them, so a window sitting in that gap listed
// the row with nothing drawn on it. Segment-wise overlap is what makes the REQ's
// GREEN — every listed row has something drawn on it — actually true.
var gapWindowStartMs = Date.UTC(2026, 5, 20, 14, 0);   // after the now-line
var gapWindowEndMs = Date.UTC(2026, 5, 20, 18, 0);     // before the forecast bar
var gapProjection = { id: "REQ-007", startTime: "2026-06-21T00:00:00Z", endTime: "2026-06-21T06:00:00Z" };
process.stdout.write(JSON.stringify({
  inWindow: idsInWindow(rows, {}),
  reversedInWindow: idsInWindow([reversedRow], {}),
  reversedExtentOrdered: reversedSegments.every(function (s) { return s.startMs <= s.endMs; }),
  reversedExtentStartIso: new Date(reversedSegments[0].startMs).toISOString(),
  projectedOnlyInWindow:
    idsInSpan([projectedRow], { "REQ-007": projection }, futureWindowStartMs, futureWindowEndMs),
  projectedIgnoredWithoutForecast:
    idsInSpan([projectedRow], {}, futureWindowStartMs, futureWindowEndMs),
  everythingInAWideWindow: idsInSpan(rows, {}, Date.UTC(2026, 0, 1), Date.UTC(2027, 0, 1)),
  inTheForecastGap:
    idsInSpan([projectedRow], { "REQ-007": gapProjection }, gapWindowStartMs, gapWindowEndMs),
  spanningTheForecastGap:
    idsInSpan([projectedRow], { "REQ-007": gapProjection }, nowMs - 3600000, Date.UTC(2026, 5, 21, 3, 0)),

  startsAtWindowEnd: idsInSpan([startsAtWindowEndRow], {}, edgeWindowStartMs, edgeWindowEndMs),
  startsJustInside: idsInSpan([startsJustInsideRow], {}, edgeWindowStartMs, edgeWindowEndMs),
  endsAtWindowStart: idsInSpan([endsAtWindowStartRow], {}, edgeWindowStartMs, edgeWindowEndMs),

  reversedHullInTheGap:
    idsInSpan([reversedHullRow], {}, reversedHullGapStartMs, reversedHullGapEndMs),
  reversedHullAtItsMarker:
    idsInSpan([reversedHullRow], {}, reversedHullMarkerStartMs, reversedHullMarkerEndMs)
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline window rows", javascriptProbe)
	var windowResult struct {
		InWindow                        []string `json:"inWindow"`
		ReversedInWindow                []string `json:"reversedInWindow"`
		ReversedExtentOrdered           bool     `json:"reversedExtentOrdered"`
		ReversedExtentStartIso          string   `json:"reversedExtentStartIso"`
		ProjectedOnlyInWindow           []string `json:"projectedOnlyInWindow"`
		ProjectedIgnoredWithoutForecast []string `json:"projectedIgnoredWithoutForecast"`
		EverythingInAWideWindow         []string `json:"everythingInAWideWindow"`
		InTheForecastGap                []string `json:"inTheForecastGap"`
		SpanningTheForecastGap          []string `json:"spanningTheForecastGap"`

		StartsAtWindowEnd []string `json:"startsAtWindowEnd"`
		StartsJustInside  []string `json:"startsJustInside"`
		EndsAtWindowStart []string `json:"endsAtWindowStart"`

		ReversedHullInTheGap    []string `json:"reversedHullInTheGap"`
		ReversedHullAtItsMarker []string `json:"reversedHullAtItsMarker"`
	}
	if decodeError := json.Unmarshal(probeOutput, &windowResult); decodeError != nil {
		t.Fatalf("decode timeline window rows behavior: %v (output %q)", decodeError, probeOutput)
	}

	wantInWindow := "REQ-002,REQ-003,REQ-004"
	if gotInWindow := strings.Join(windowResult.InWindow, ","); gotInWindow != wantInWindow {
		t.Fatalf("rows in the window = %q, want %q — the straddling and still-running REQs "+
			"overlap it and must be listed; the two outside it must not be",
			gotInWindow, wantInWindow)
	}
	if !windowResult.ReversedExtentOrdered {
		t.Fatalf("a reversed-stamp row produced a segment whose start (%s) is after its end; "+
			"a reversed span is a point at the break marker's own instant, and an inverted "+
			"segment makes the overlap test drop broken rows from every window",
			windowResult.ReversedExtentStartIso)
	}
	if len(windowResult.ReversedInWindow) != 1 {
		t.Fatalf("the reversed-stamp row is inside the window and was not listed (got %v)",
			windowResult.ReversedInWindow)
	}
	if len(windowResult.ProjectedOnlyInWindow) != 1 {
		t.Fatal("a pending REQ whose forecast bar is the only thing it draws in the window " +
			"was not listed; the extent must reach the projected segment")
	}
	// The control for the assertion above: same row, same window, no forecast.
	// Without it the extent stops at the now-line and the row is genuinely absent,
	// which is what proves the projection — not the fixture — put it in the window.
	if len(windowResult.ProjectedIgnoredWithoutForecast) != 0 {
		t.Fatalf("the same pending REQ reached a window beyond the now-line with no forecast "+
			"attached (got %v); its measured extent ends at the now-line",
			windowResult.ProjectedIgnoredWithoutForecast)
	}
	if len(windowResult.EverythingInAWideWindow) != 5 {
		t.Fatalf("a window spanning the whole year listed %d of 5 rows; widening the window "+
			"must never drop a row", len(windowResult.EverythingInAWideWindow))
	}
	// The pair that forces segment-wise overlap rather than a hull. Both windows
	// sit between the same two marks; only the second one touches either.
	if len(windowResult.InTheForecastGap) != 0 {
		t.Fatalf("a window in the gap between a REQ's open wait and its forecast bar listed "+
			"%v; the row draws nothing there, and listing it is what window-scoping exists to "+
			"stop — a hull over the two marks spans the gap between them",
			windowResult.InTheForecastGap)
	}
	if len(windowResult.SpanningTheForecastGap) != 1 {
		t.Fatal("a window reaching across both the open wait and the forecast bar did not list " +
			"the row; the gap rule must not cost a row that genuinely draws inside the window")
	}

	// THE EXCLUSIVE END, and its two controls. windowEndMs is the next window's
	// first instant everywhere else in this module — the end field renders
	// windowEndMs - 1 to say so — so admitting a segment that begins exactly there
	// lists a row whose only mark is clipped at the right edge.
	if len(windowResult.StartsAtWindowEnd) != 0 {
		t.Errorf("a REQ whose span begins exactly at the window's exclusive end was listed "+
			"(got %v); its floored rectangle lands at the clipped right edge, so the row shows "+
			"nothing — and it belongs to the NEXT window", windowResult.StartsAtWindowEnd)
	}
	if len(windowResult.StartsJustInside) != 1 {
		t.Errorf("the same REQ one millisecond earlier was not listed (got %v); the end is "+
			"exclusive by one instant, not by a margin", windowResult.StartsJustInside)
	}
	// Deliberately asymmetric. The START instant IS in the window and a span
	// ending on it draws a visible floored mark at x=0, so this stays inclusive.
	if len(windowResult.EndsAtWindowStart) != 1 {
		t.Errorf("a REQ whose span ends exactly at the window's start was dropped (got %v); "+
			"the start is inclusive and that span draws a mark at the left edge",
			windowResult.EndsAtWindowStart)
	}

	// A REVERSED SPAN IS A POINT, because that is what the renderer draws. The
	// same defect as the forecast gap: a hull over an interval nothing is drawn
	// across lists rows with nothing on them.
	if len(windowResult.ReversedHullInTheGap) != 0 {
		t.Errorf("a window inside a reversed span's min/max interval listed the row (got %v), "+
			"but renderVisibleRows draws only a break marker at the created instant and the "+
			"row's other bar sits elsewhere — nothing is drawn in that window",
			windowResult.ReversedHullInTheGap)
	}
	if len(windowResult.ReversedHullAtItsMarker) != 1 {
		t.Errorf("a window over the reversed span's own break marker did not list the row "+
			"(got %v); the point has to sit where the renderer puts the marker, or a broken "+
			"row becomes unfindable — which is what the min/max existed to prevent",
			windowResult.ReversedHullAtItsMarker)
	}
}

func TestJavaScriptBehaviorTimelineRowLabelTruncation(t *testing.T) {
	indexHtml := generateLiveSite(t)
	// The column width is READ FROM THE RENDERER, not restated here. Restating it
	// is how the first version of this test passed with the constant reverted to
	// its old value: the budget floor below was measuring a column the board had
	// stopped using (REQ-265 — grep the quantity, not the constant name).
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_LABEL_WIDTH") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelCharacterBudget(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelCellCount(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineLabelPrefixWithinCells(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowLabelText(") + `
var longTitle = "Colour timeline bars by REQ status";
// The shipped column, minus the label's own 6px x-offset and a 6px gap before
// the first bar — the same expression the renderer uses.
var shippedColumnCells = timelineLabelCharacterBudget(6.0219, TIMELINE_LABEL_WIDTH - 12);
process.stdout.write(JSON.stringify({
  shippedColumnCells: shippedColumnCells,
  shippedLabelSample: timelineRowLabelText("REQ-042", longTitle, shippedColumnCells),
  taggedLabelSample: timelineRowLabelText(
    "REQ-306", "[impact-rule-change] Judge effort_estimate on review-minted follow-ups too",
    shippedColumnCells),
  cjkCells: timelineLabelCellCount("修复部署管道"),
  asciiCells: timelineLabelCellCount("abcdef"),
  astralCells: timelineLabelCellCount("🚀"),
  cjkLabel: timelineRowLabelText("REQ-042", "修复部署管道的一个严重问题在这里", 28),
  astralBoundary: timelineRowLabelText("REQ-042", "Fix the deploy 🚀 pipeline right now", 26),
  budgetAt6px: timelineLabelCharacterBudget(6, 172),
  budgetWithNoAdvance: timelineLabelCharacterBudget(0, 172),
  budgetWithNegativeAdvance: timelineLabelCharacterBudget(-3, 172),

  roomy: timelineRowLabelText("REQ-042", longTitle, 60),
  cut: timelineRowLabelText("REQ-042", longTitle, 30),
  cutLength: timelineRowLabelText("REQ-042", longTitle, 30).length,
  tight: timelineRowLabelText("REQ-042", longTitle, 14),
  noBudget: timelineRowLabelText("REQ-042", longTitle, 0),
  noTitle: timelineRowLabelText("REQ-042", "", 60),
  exactFit: timelineRowLabelText("REQ-042", "abcdef", 15),
  oneOver: timelineRowLabelText("REQ-042", "abcdefghijkl", 20),
  longId: timelineRowLabelText("REQ-100042", longTitle, 14)
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline row label", javascriptProbe)
	var labelResult struct {
		ShippedColumnCells        int    `json:"shippedColumnCells"`
		ShippedLabelSample        string `json:"shippedLabelSample"`
		TaggedLabelSample         string `json:"taggedLabelSample"`
		CjkCells                  int    `json:"cjkCells"`
		AsciiCells                int    `json:"asciiCells"`
		AstralCells               int    `json:"astralCells"`
		CjkLabel                  string `json:"cjkLabel"`
		AstralBoundary            string `json:"astralBoundary"`
		BudgetAt6px               int    `json:"budgetAt6px"`
		BudgetWithNoAdvance       int    `json:"budgetWithNoAdvance"`
		BudgetWithNegativeAdvance int    `json:"budgetWithNegativeAdvance"`
		Roomy                     string `json:"roomy"`
		Cut                       string `json:"cut"`
		CutLength                 int    `json:"cutLength"`
		Tight                     string `json:"tight"`
		NoBudget                  string `json:"noBudget"`
		NoTitle                   string `json:"noTitle"`
		ExactFit                  string `json:"exactFit"`
		OneOver                   string `json:"oneOver"`
		LongId                    string `json:"longId"`
	}
	if decodeError := json.Unmarshal(probeOutput, &labelResult); decodeError != nil {
		t.Fatalf("decode timeline row label behavior: %v (output %q)", decodeError, probeOutput)
	}

	// THE COLUMN WIDTH, pinned to the shipped constant. Reverting
	// TIMELINE_LABEL_WIDTH to its pre-REQ 104 used to pass this whole file; now it
	// fails here, which is what makes the width a decision the tests hold rather
	// than a number in a comment.
	if labelResult.ShippedColumnCells < 20 {
		t.Errorf("the shipped label column fits %d cells; below about 20 the title is a stub and "+
			"the column stops earning the plot width it costs (label sample %q)",
			labelResult.ShippedColumnCells, labelResult.ShippedLabelSample)
	}
	// A classification tag is metadata for the board's search box, not title text.
	// Unstripped it consumed the entire budget and every review-minted REQ read
	// "[impact-user-visib…".
	if strings.Contains(labelResult.TaggedLabelSample, "[impact-") {
		t.Errorf("a tagged title rendered %q; the leading [impact-…] tag has to come off in the "+
			"label or it is the whole budget", labelResult.TaggedLabelSample)
	}
	if !strings.Contains(labelResult.TaggedLabelSample, "Judge") {
		t.Errorf("a tagged title rendered %q, want the actual title after the tag came off",
			labelResult.TaggedLabelSample)
	}

	// Cells, not characters. The measured advance describes the face's Latin cell,
	// and a face that draws 中 at 10px against a 6.02px cell overruns the column.
	if labelResult.AsciiCells != 6 {
		t.Errorf("six ASCII characters counted %d cells, want 6", labelResult.AsciiCells)
	}
	if labelResult.CjkCells != 12 {
		t.Errorf("six CJK characters counted %d cells, want 12 — one cell each is what let a CJK "+
			"title draw 36px into the plot", labelResult.CjkCells)
	}
	if labelResult.AstralCells != 2 {
		t.Errorf("one astral character counted %d cells, want 2", labelResult.AstralCells)
	}
	// The cut lands on a code point boundary, so an astral character is never
	// split into a lone surrogate that renders as a fallback box.
	for _, character := range labelResult.AstralBoundary {
		if character == '\uFFFD' {
			t.Errorf("a title cut near an astral character produced %q, which contains a "+
				"replacement character — the cut split a surrogate pair",
				labelResult.AstralBoundary)
		}
	}
	if labelResult.CjkLabel == "" || !strings.HasPrefix(labelResult.CjkLabel, "REQ-042") {
		t.Errorf("a CJK title rendered %q, want a label led by its id", labelResult.CjkLabel)
	}

	if labelResult.BudgetAt6px != 28 {
		t.Errorf("a 172px column at a 6px advance fits %d characters, want 28", labelResult.BudgetAt6px)
	}
	// An unmeasurable face must produce NO budget rather than a plausible one.
	// A guessed advance is the REQ-292 defect: a number that looks like a
	// measurement and does not move when the face does.
	if labelResult.BudgetWithNoAdvance != 0 || labelResult.BudgetWithNegativeAdvance != 0 {
		t.Errorf("an unmeasured face produced budgets %d and %d, want 0 for both",
			labelResult.BudgetWithNoAdvance, labelResult.BudgetWithNegativeAdvance)
	}

	if labelResult.Roomy != "REQ-042  Colour timeline bars by REQ status" {
		t.Errorf("a roomy budget rendered %q, want the id and the whole title", labelResult.Roomy)
	}
	if !strings.HasPrefix(labelResult.Cut, "REQ-042  ") || !strings.HasSuffix(labelResult.Cut, "…") {
		t.Errorf("a cut label rendered %q, want the id then a truncated title ending in an ellipsis",
			labelResult.Cut)
	}
	// The ellipsis is inside the budget, not on top of it — otherwise the label
	// the arithmetic says fits is one character wider than the column.
	if labelResult.CutLength > 30 {
		t.Errorf("a 30-character budget produced a %d-character label %q; the ellipsis has to fit "+
			"inside the budget", labelResult.CutLength, labelResult.Cut)
	}
	// Too tight for a useful title: the id survives whole, alone.
	if labelResult.Tight != "REQ-042" {
		t.Errorf("a tight budget rendered %q, want the id alone — a half-drawn id is worse than "+
			"no title", labelResult.Tight)
	}
	if labelResult.NoBudget != "REQ-042" || labelResult.NoTitle != "REQ-042" {
		t.Errorf("no budget rendered %q and no title rendered %q, want the id in both cases",
			labelResult.NoBudget, labelResult.NoTitle)
	}
	if labelResult.ExactFit != "REQ-042  abcdef" {
		t.Errorf("a title that exactly fills the budget rendered %q, want it whole and unmarked",
			labelResult.ExactFit)
	}
	if !strings.HasSuffix(labelResult.OneOver, "…") || len([]rune(labelResult.OneOver)) > 20 {
		t.Errorf("a title one character over the budget rendered %q (%d runes); it must be cut, "+
			"marked, and still inside the budget",
			labelResult.OneOver, len([]rune(labelResult.OneOver)))
	}
	// A longer id eats the same budget. The rule holds by id length, not by a
	// hard-coded seven characters.
	if labelResult.LongId != "REQ-100042" {
		t.Errorf("a ten-character id at a 14-character budget rendered %q, want the id alone",
			labelResult.LongId)
	}
}

func TestJavaScriptBehaviorTimelineGridlinesShareTheAxisTicks(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_AXIS_TICK_COUNT", "TIMELINE_DAY_MS", "TIMELINE_AXIS_TICK_LIMIT") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_WEEK_ALIGNMENT_MS") + "\n" +
		rendererBracketDeclaration(t, "web/board-timeline.js", "TIMELINE_AXIS_TICK_STEPS") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickStepSpanMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTickStep(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickAtOrBefore(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedTick(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTicks(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTickInstants(") + `
var windowStartMs = Date.UTC(2026, 5, 1);
var windowEndMs = Date.UTC(2026, 5, 8);
var instants = timelineAxisTickInstants(windowStartMs, windowEndMs);
process.stdout.write(JSON.stringify({
  tickCount: TIMELINE_AXIS_TICK_COUNT,
  instantCount: instants.length,
  firstIso: new Date(instants[0]).toISOString(),
  lastIso: new Date(instants[instants.length - 1]).toISOString(),
  ascending: instants.every(function (instant, index) {
    return index === 0 || instant > instants[index - 1];
  }),
  // A zero-width window is reachable: the fields accept one day in both boxes
  // before the settle widens it. It must produce a tick list, not NaNs.
  degenerateFinite: timelineAxisTickInstants(windowStartMs, windowStartMs)
    .every(function (instant) { return isFinite(instant); }),
  degenerateCount: timelineAxisTickInstants(windowStartMs, windowStartMs).length
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline axis ticks", javascriptProbe)
	var tickResult struct {
		TickCount        int    `json:"tickCount"`
		InstantCount     int    `json:"instantCount"`
		FirstIso         string `json:"firstIso"`
		LastIso          string `json:"lastIso"`
		Ascending        bool   `json:"ascending"`
		DegenerateFinite bool   `json:"degenerateFinite"`
		DegenerateCount  int    `json:"degenerateCount"`
	}
	if decodeError := json.Unmarshal(probeOutput, &tickResult); decodeError != nil {
		t.Fatalf("decode timeline axis tick behavior: %v (output %q)", decodeError, probeOutput)
	}

	// REQ-327 CHANGED WHAT THIS COUNTS. The axis no longer divides the window into
	// TIMELINE_AXIS_TICK_COUNT equal parts, so the old "one instant per tick plus
	// both window edges" arithmetic no longer describes it — that arithmetic is the
	// defect: on a seven-day window it put ticks 28 hours apart and the formatter
	// labelled each with a bare date. A June week now gets eight midnights, and the
	// count varies with the rung of the ladder the span picks.
	if tickResult.InstantCount != 8 {
		t.Errorf("a seven-day window produced %d ticks, want the 8 midnights that span it "+
			"(1 June through 8 June inclusive)", tickResult.InstantCount)
	}
	if tickResult.FirstIso != "2026-06-01T00:00:00.000Z" || tickResult.LastIso != "2026-06-08T00:00:00.000Z" {
		t.Errorf("the tick list spans %s → %s, want 2026-06-01 → 2026-06-08",
			tickResult.FirstIso, tickResult.LastIso)
	}
	if tickResult.DegenerateCount != 1 {
		t.Errorf("a zero-width window produced %d ticks, want exactly 1 — six copies of one "+
			"instant is what the old equal-parts loop emitted", tickResult.DegenerateCount)
	}
	if !tickResult.Ascending {
		t.Error("the tick instants are not strictly ascending")
	}
	if !tickResult.DegenerateFinite {
		t.Error("a zero-width window produced non-finite tick instants; the fields can reach " +
			"one before the settle widens it, and NaN ticks draw nothing and log nothing")
	}

	// ONE source, and this is the half that bites. The probe above would pass just
	// as well with renderAxis keeping a private copy of the walk, which is exactly
	// the drift the extraction removed — so count the callers instead of trusting
	// the function's existence.
	//
	// The BOUNDARY WALK specifically: the loop that steps from one tick to the next
	// is what places a tick, and timelineSteppedTick is the only thing entitled to
	// do it. Inlining it back into either caller puts a second copy in the page and
	// fails here. (This replaced a count of "tickIndex) / TIMELINE_AXIS_TICK_COUNT",
	// the equal-parts expression REQ-327 deleted.)
	// Counted as "calls from anywhere else", not as a raw string count: the walk
	// legitimately calls timelineSteppedTick twice inside timelineAxisTicks (once to
	// skip a boundary that precedes the window, once per step), so a bare count of 1
	// would be wrong about the healthy shape and could only be satisfied by making
	// the code worse.
	axisTicksBody := sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTicks(")
	steppedTickDefinition := sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedTick(")
	callsEverywhere := strings.Count(indexHtml, "timelineSteppedTick(")
	callsInTheWalk := strings.Count(axisTicksBody, "timelineSteppedTick(")
	callsInItsOwnDefinition := strings.Count(steppedTickDefinition, "timelineSteppedTick(")
	if callsElsewhere := callsEverywhere - callsInTheWalk - callsInItsOwnDefinition; callsElsewhere != 0 {
		t.Errorf("timelineSteppedTick is called from %d place(s) outside timelineAxisTicks; the "+
			"boundary walk lives in one function or the gridlines can start meaning a different "+
			"instant than the ticks above them", callsElsewhere)
	}
	if callsInTheWalk == 0 {
		t.Error("timelineAxisTicks does not call timelineSteppedTick, so the check above is " +
			"comparing two numbers that no longer describe the walk")
	}
	// renderAxis needs the GAP as well as the instants, so it calls timelineAxisTicks
	// directly; drawGridlines needs only the instants. Both bottom out in one walk,
	// and each is named here so neither can quietly grow its own.
	for caller, wantCall := range map[string]string{
		"function renderAxis(":    "timelineAxisTicks(",
		"function drawGridlines(": "timelineAxisTickInstants(",
	} {
		callerBody := sliceBalancedBlockAfter(t, indexHtml, caller)
		if !strings.Contains(callerBody, wantCall) {
			t.Errorf("%s does not read %s; the axis and the gridlines have to come from one "+
				"list or they can disagree", caller, wantCall)
		}
	}
	// And the gap the labels are formatted against comes from the same call that
	// positioned the ticks, rather than being derived a second time.
	renderAxisBody := sliceBalancedBlockAfter(t, indexHtml, "function renderAxis(")
	if !strings.Contains(renderAxisBody, "axisTicks.gapMs") {
		t.Error("renderAxis does not pass the chosen gap to timelineFormatAxisTick; a gap " +
			"derived a second time is how a date-only label came to sit on a 04:00 tick")
	}
}

func TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_ZOOM_STEP", "TIMELINE_PAN_FRACTION") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelinePannedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineKeyboardWindow(") + `
var boundStart = 0;
var boundEnd = 30 * 24 * 3600 * 1000;   // a 30-day board

// Start half zoomed in, so a pan has room to run in both directions before it
// meets a bound.
var halfway = timelineZoomedWindow(boundStart, boundEnd, 2, 0.5, boundStart, boundEnd);
var halfSpanMs = halfway.windowEndMs - halfway.windowStartMs;

function pressKey(currentWindow, keyName) {
  var moved = timelineKeyboardWindow(
    keyName, currentWindow.windowStartMs, currentWindow.windowEndMs, boundStart, boundEnd);
  return moved || currentWindow;
}
function repeatKey(currentWindow, keyName, pressCount) {
  for (var press = 0; press < pressCount; press++) {
    currentWindow = pressKey(currentWindow, keyName);
  }
  return currentWindow;
}

var pannedRight = pressKey(halfway, "ArrowRight");
var pannedLeft = pressKey(halfway, "ArrowLeft");

// Held to the edges: the window must stop AT the bound, keeping its span.
var atRightEdge = repeatKey(halfway, "ArrowRight", 40);
var atLeftEdge = repeatKey(halfway, "ArrowLeft", 40);

// The pointer path's own floor and ceiling, reached through the wheel's
// off-centre anchor rather than the keyboard's centred one.
var pointerFloor = halfway;
var pointerCeiling = halfway;
for (var pointerStep = 0; pointerStep < 40; pointerStep++) {
  pointerFloor = timelineZoomedWindow(
    pointerFloor.windowStartMs, pointerFloor.windowEndMs, TIMELINE_ZOOM_STEP, 0.25, boundStart, boundEnd);
  pointerCeiling = timelineZoomedWindow(
    pointerCeiling.windowStartMs, pointerCeiling.windowEndMs, 1 / TIMELINE_ZOOM_STEP, 0.25, boundStart, boundEnd);
}
var keyboardFloor = repeatKey(halfway, "+", 40);
var keyboardCeiling = repeatKey(halfway, "-", 40);

process.stdout.write(JSON.stringify({
  panStepMs: pannedRight.windowStartMs - halfway.windowStartMs,
  panBackStepMs: halfway.windowStartMs - pannedLeft.windowStartMs,
  wantPanStepMs: halfSpanMs * TIMELINE_PAN_FRACTION,
  windowSpanMs: halfSpanMs,
  panKeepsSpan:
    pannedRight.windowEndMs - pannedRight.windowStartMs === halfSpanMs &&
    pannedLeft.windowEndMs - pannedLeft.windowStartMs === halfSpanMs,
  rightEdgeMs: atRightEdge.windowEndMs,
  leftEdgeMs: atLeftEdge.windowStartMs,
  boundStartMs: boundStart,
  boundEndMs: boundEnd,
  edgesKeepSpan:
    atRightEdge.windowEndMs - atRightEdge.windowStartMs === halfSpanMs &&
    atLeftEdge.windowEndMs - atLeftEdge.windowStartMs === halfSpanMs,
  keyboardFloorSpanMs: keyboardFloor.windowEndMs - keyboardFloor.windowStartMs,
  pointerFloorSpanMs: pointerFloor.windowEndMs - pointerFloor.windowStartMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  keyboardCeilingSpanMs: keyboardCeiling.windowEndMs - keyboardCeiling.windowStartMs,
  pointerCeilingSpanMs: pointerCeiling.windowEndMs - pointerCeiling.windowStartMs,
  boundSpanMs: boundEnd - boundStart,
  unownedKeys: ["Enter", " ", "Spacebar", "Tab", "ArrowUp", "ArrowDown", "a"].map(function (keyName) {
    return timelineKeyboardWindow(keyName, halfway.windowStartMs, halfway.windowEndMs, boundStart, boundEnd);
  })
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline keyboard pan and zoom", javascriptProbe)
	var keyboardResult struct {
		PanStepMs             float64            `json:"panStepMs"`
		PanBackStepMs         float64            `json:"panBackStepMs"`
		WantPanStepMs         float64            `json:"wantPanStepMs"`
		WindowSpanMs          float64            `json:"windowSpanMs"`
		PanKeepsSpan          bool               `json:"panKeepsSpan"`
		RightEdgeMs           float64            `json:"rightEdgeMs"`
		LeftEdgeMs            float64            `json:"leftEdgeMs"`
		BoundStartMs          float64            `json:"boundStartMs"`
		BoundEndMs            float64            `json:"boundEndMs"`
		EdgesKeepSpan         bool               `json:"edgesKeepSpan"`
		KeyboardFloorSpanMs   float64            `json:"keyboardFloorSpanMs"`
		PointerFloorSpanMs    float64            `json:"pointerFloorSpanMs"`
		MinSpanMs             float64            `json:"minSpanMs"`
		KeyboardCeilingSpanMs float64            `json:"keyboardCeilingSpanMs"`
		PointerCeilingSpanMs  float64            `json:"pointerCeilingSpanMs"`
		BoundSpanMs           float64            `json:"boundSpanMs"`
		UnownedKeys           []*json.RawMessage `json:"unownedKeys"`
	}
	if decodeError := json.Unmarshal(probeOutput, &keyboardResult); decodeError != nil {
		t.Fatalf("decode timeline keyboard behavior: %v (output %q)", decodeError, probeOutput)
	}

	// A pan step has to be a fraction of what is on screen: a fixed number of
	// milliseconds is either imperceptible zoomed out or a jump zoomed in.
	if math.Abs(keyboardResult.PanStepMs-keyboardResult.WantPanStepMs) > 1 {
		t.Fatalf("ArrowRight moved the window %.0f ms, want %.0f ms — one step of the visible span",
			keyboardResult.PanStepMs, keyboardResult.WantPanStepMs)
	}
	if math.Abs(keyboardResult.PanBackStepMs-keyboardResult.WantPanStepMs) > 1 {
		t.Fatalf("ArrowLeft moved the window %.0f ms, want %.0f ms back",
			keyboardResult.PanBackStepMs, keyboardResult.WantPanStepMs)
	}
	if keyboardResult.PanStepMs <= 0 || keyboardResult.PanStepMs >= keyboardResult.WindowSpanMs {
		t.Fatalf("a pan step of %.0f ms against a %.0f ms window is not a bounded step; a reader loses their place",
			keyboardResult.PanStepMs, keyboardResult.WindowSpanMs)
	}
	if !keyboardResult.PanKeepsSpan {
		t.Fatal("panning changed the window span; panning moves the window, zooming resizes it")
	}

	// Held down, a pan must stop at the range edge rather than walking the window
	// off the data — the same clamp the drag path applies.
	if math.Abs(keyboardResult.RightEdgeMs-keyboardResult.BoundEndMs) > 1 {
		t.Fatalf("panning right settled with the window ending at %.0f ms, want the range edge %.0f ms",
			keyboardResult.RightEdgeMs, keyboardResult.BoundEndMs)
	}
	if math.Abs(keyboardResult.LeftEdgeMs-keyboardResult.BoundStartMs) > 1 {
		t.Fatalf("panning left settled with the window starting at %.0f ms, want the range edge %.0f ms",
			keyboardResult.LeftEdgeMs, keyboardResult.BoundStartMs)
	}
	if !keyboardResult.EdgesKeepSpan {
		t.Fatal("clamping at a range edge changed the window span; it must only stop the window, not shrink it")
	}

	// The point of routing the keys through timelineZoomedWindow: one floor and
	// one ceiling, whichever driver arrives at them.
	if keyboardResult.KeyboardFloorSpanMs != keyboardResult.PointerFloorSpanMs {
		t.Fatalf("`+` bottomed out at %.0f ms but the pointer path bottoms out at %.0f ms; the two have diverged",
			keyboardResult.KeyboardFloorSpanMs, keyboardResult.PointerFloorSpanMs)
	}
	if keyboardResult.KeyboardFloorSpanMs != keyboardResult.MinSpanMs {
		t.Fatalf("`+` bottomed out at %.0f ms, want the renderer's %.0f ms floor",
			keyboardResult.KeyboardFloorSpanMs, keyboardResult.MinSpanMs)
	}
	if keyboardResult.KeyboardCeilingSpanMs != keyboardResult.PointerCeilingSpanMs {
		t.Fatalf("`-` topped out at %.0f ms but the pointer path tops out at %.0f ms; the two have diverged",
			keyboardResult.KeyboardCeilingSpanMs, keyboardResult.PointerCeilingSpanMs)
	}
	if keyboardResult.KeyboardCeilingSpanMs != keyboardResult.BoundSpanMs {
		t.Fatalf("`-` topped out at %.0f ms, want the full range span %.0f ms",
			keyboardResult.KeyboardCeilingSpanMs, keyboardResult.BoundSpanMs)
	}

	// Enter and Space belong to row activation, and Up/Down to scrolling the
	// queue. Claiming any of them would take a working interaction away.
	unownedKeyNames := []string{"Enter", "Space", "Spacebar", "Tab", "ArrowUp", "ArrowDown", "a"}
	if len(keyboardResult.UnownedKeys) != len(unownedKeyNames) {
		t.Fatalf("probe reported %d unowned-key results, want %d", len(keyboardResult.UnownedKeys), len(unownedKeyNames))
	}
	for keyIndex, unownedKeyName := range unownedKeyNames {
		if keyboardResult.UnownedKeys[keyIndex] != nil {
			t.Fatalf("%s moved the time window to %s; that key belongs to row activation or to scrolling",
				unownedKeyName, string(*keyboardResult.UnownedKeys[keyIndex]))
		}
	}
}

func TestJavaScriptBehaviorTimelineTrailingWindowsSurviveADrainedQueue(t *testing.T) {
	indexHtml := generateLiveSite(t)

	chipValuePattern := regexp.MustCompile(`data-timeline-period="([^"]*)"`)
	var chipValues []string
	for _, chipMatch := range chipValuePattern.FindAllStringSubmatch(indexHtml, -1) {
		chipValues = append(chipValues, chipMatch[1])
	}
	if len(chipValues) < 2 {
		t.Fatalf("the generated page declares %d trailing-window chips (%v); a sweep for distinct "+
			"windows needs the shipped control set to have something to be distinct about",
			len(chipValues), chipValues)
	}
	chipValuesJson, marshalError := json.Marshal(chipValues)
	if marshalError != nil {
		t.Fatalf("marshal the shipped chip values: %v", marshalError)
	}

	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTrailingWindow(") + "\n" +
		"var chipValues = " + string(chipValuesJson) + ";" + `
// A board whose RECORDED RANGE is a fixed 95 days while NOW walks away from it,
// which is exactly what happens as a queue drains: timelineRange stops following
// now the moment the last row closes, and nothing moves the range end after that.
var boundStartMs = Date.UTC(2026, 3, 7);
var boundEndMs = boundStartMs + 95 * TIMELINE_DAY_MS;

var collapsedToFloor = [];
var sharedWindows = [];
var wrongWindow = [];
var wrongClippedFlag = [];
var clippedFlagTrueCount = 0;
var clippedFlagFalseCount = 0;
var sampleCount = 0;
var measuredTable = [];
var measuredAges = { "3": true, "10": true, "40": true, "100": true };

function isoOf(epochMs) { return new Date(epochMs).toISOString().slice(0, 16) + "Z"; }
function hoursOf(spanMs) { return (spanMs / 3600000).toFixed(2) + "h"; }
function noteAtMost(list, entry) { if (list.length < 8) { list.push(entry); } }

// Idle days: negative is a live board, where the payload's range end sits past
// now because a forecast and the cosmetic padding are drawn to the right of it;
// zero is the instant the queue drains; positive is a board nobody has touched
// since. The whole span is swept because any of it is a board somebody has.
for (var idleDays = -20; idleDays <= 120; idleDays++) {
  var nowMs = boundEndMs + idleDays * TIMELINE_DAY_MS;
  var anchorMs = Math.min(Math.max(nowMs, boundStartMs), boundEndMs);
  var windowsThisSample = [];
  sampleCount++;
  for (var chipIndex = 0; chipIndex < chipValues.length; chipIndex++) {
    var chipValue = chipValues[chipIndex];
    var chipWindow = timelineTrailingWindow(chipValue, nowMs, boundStartMs, boundEndMs);
    var chipSpanMs = chipWindow.windowEndMs - chipWindow.windowStartMs;
    var askedDayCount = Number(chipValue);
    var asksForATrailingSpan = isFinite(askedDayCount) && askedDayCount > 0;
    var askedStartMs = anchorMs - askedDayCount * TIMELINE_DAY_MS;
    var wantStartMs = asksForATrailingSpan ? Math.max(askedStartMs, boundStartMs) : boundStartMs;
    var wantEndMs = asksForATrailingSpan ? anchorMs : boundEndMs;
    var wantClipped = asksForATrailingSpan && askedStartMs < boundStartMs;
    if (wantClipped) { clippedFlagTrueCount++; } else { clippedFlagFalseCount++; }

    if (chipSpanMs <= TIMELINE_MIN_SPAN_MS) {
      noteAtMost(collapsedToFloor, idleDays + "d idle: chip " + chipValue + " settled on a " +
        hoursOf(chipSpanMs) + " window at " + isoOf(chipWindow.windowStartMs));
    }
    for (var seenIndex = 0; seenIndex < windowsThisSample.length; seenIndex++) {
      var seenWindow = windowsThisSample[seenIndex];
      if (seenWindow.startMs === chipWindow.windowStartMs &&
        seenWindow.endMs === chipWindow.windowEndMs) {
        noteAtMost(sharedWindows, idleDays + "d idle: chips " + seenWindow.value + " and " +
          chipValue + " share " + isoOf(chipWindow.windowStartMs) + " -> " +
          isoOf(chipWindow.windowEndMs));
      }
    }
    if (chipWindow.windowStartMs !== wantStartMs || chipWindow.windowEndMs !== wantEndMs) {
      noteAtMost(wrongWindow, idleDays + "d idle: chip " + chipValue + " gives " +
        isoOf(chipWindow.windowStartMs) + " -> " + isoOf(chipWindow.windowEndMs) + ", want " +
        isoOf(wantStartMs) + " -> " + isoOf(wantEndMs));
    }
    if (chipWindow.isClippedByBounds !== wantClipped) {
      noteAtMost(wrongClippedFlag, idleDays + "d idle: chip " + chipValue +
        " reports isClippedByBounds=" + chipWindow.isClippedByBounds + ", want " + wantClipped);
    }
    windowsThisSample.push({
      value: chipValue,
      startMs: chipWindow.windowStartMs,
      endMs: chipWindow.windowEndMs
    });
    if (measuredAges[String(idleDays)]) {
      measuredTable.push(idleDays + "d idle  chip " + chipValue + "  " + hoursOf(chipSpanMs) +
        "  " + isoOf(chipWindow.windowStartMs) + " -> " + isoOf(chipWindow.windowEndMs));
    }
  }
}

process.stdout.write(JSON.stringify({
  sampleCount: sampleCount,
  chipCount: chipValues.length,
  collapsedToFloor: collapsedToFloor,
  sharedWindows: sharedWindows,
  wrongWindow: wrongWindow,
  wrongClippedFlag: wrongClippedFlag,
  clippedFlagTrueCount: clippedFlagTrueCount,
  clippedFlagFalseCount: clippedFlagFalseCount,
  measuredTable: measuredTable
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline trailing windows on a drained queue", javascriptProbe)
	var drainedResult struct {
		SampleCount           int      `json:"sampleCount"`
		ChipCount             int      `json:"chipCount"`
		CollapsedToFloor      []string `json:"collapsedToFloor"`
		SharedWindows         []string `json:"sharedWindows"`
		WrongWindow           []string `json:"wrongWindow"`
		WrongClippedFlag      []string `json:"wrongClippedFlag"`
		ClippedFlagTrueCount  int      `json:"clippedFlagTrueCount"`
		ClippedFlagFalseCount int      `json:"clippedFlagFalseCount"`
		MeasuredTable         []string `json:"measuredTable"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drainedResult); decodeError != nil {
		t.Fatalf("decode the drained-queue trailing-window sweep: %v (output %q)", decodeError, probeOutput)
	}

	// The sweep has to have swept. Without this the four assertions below all pass
	// against an empty loop, which is the one mutation none of them can see.
	const sweptBoardAgeCount = 141
	if drainedResult.SampleCount != sweptBoardAgeCount || drainedResult.ChipCount != len(chipValues) {
		t.Fatalf("the sweep visited %d board ages across %d chips, want %d ages across the %d chips "+
			"the page declares", drainedResult.SampleCount, drainedResult.ChipCount,
			sweptBoardAgeCount, len(chipValues))
	}
	// And both verdicts of the clipped flag have to occur in it, or asserting the
	// flag is asserting one constant.
	if drainedResult.ClippedFlagTrueCount == 0 || drainedResult.ClippedFlagFalseCount == 0 {
		t.Fatalf("the sweep expects the clipped verdict to be true %d times and false %d times; "+
			"a sweep that never reaches one of them cannot tell the flag from a constant",
			drainedResult.ClippedFlagTrueCount, drainedResult.ClippedFlagFalseCount)
	}

	// (1) THE REPORTED SYMPTOM. A chip on the zoom floor is a dead window: the
	// board has 95 days to show and the reader asked for a day of them.
	if len(drainedResult.CollapsedToFloor) > 0 {
		t.Errorf("trailing-window chips collapsed onto the one-hour zoom floor on a board with 95 "+
			"days of range:\n\t%s", strings.Join(drainedResult.CollapsedToFloor, "\n\t"))
	}

	// (2) AND ITS CONSEQUENCE. Two chips on one window means the lit chip is
	// whichever comes first in the DOM, not the one the reader pressed.
	if len(drainedResult.SharedWindows) > 0 {
		t.Errorf("two trailing-window chips produced the same window, so pressing the second lights "+
			"the first:\n\t%s", strings.Join(drainedResult.SharedWindows, "\n\t"))
	}

	// (3) The window each chip is FOR: the last N days of the recorded range, cut
	// short at the range start when the board is younger than the span asked for.
	if len(drainedResult.WrongWindow) > 0 {
		t.Errorf("trailing-window chips landed somewhere other than the last N days of the recorded "+
			"range:\n\t%s", strings.Join(drainedResult.WrongWindow, "\n\t"))
	}

	// (4) And the window says whether it got what was asked for, so the toolbar's
	// "part of" readout reads a verdict rather than recomputing the clamp.
	if len(drainedResult.WrongClippedFlag) > 0 {
		t.Errorf("the window's own clipped verdict disagrees with whether the bounds cut it "+
			"short:\n\t%s", strings.Join(drainedResult.WrongClippedFlag, "\n\t"))
	}

	if t.Failed() {
		t.Logf("the four board ages the review measured:\n\t%s",
			strings.Join(drainedResult.MeasuredTable, "\n\t"))
	}
}

func TestJavaScriptBehaviorTimelineNowJumpLandsOnTheOpenWork(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_ROW_HEIGHT",
		"TIMELINE_NOW_JUMP_MARGIN_FRACTION", "TIMELINE_NOW_JUMP_MINIMUM_SPAN_MS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFirstOpenRowIndex(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineNowJump(") + `
// A board whose range spans five months — the shape the user reported, where the
// only sideways movement was a very long drag.
var boundStart = Date.UTC(2026, 3, 7);        // 7 Apr 2026
var boundEnd = Date.UTC(2026, 8, 2);          // 2 Sep 2026
var nowMs = Date.UTC(2026, 7, 18, 10, 30);
var queueEndMs = Date.UTC(2026, 7, 30, 6, 0); // the projection's queue-empty instant

// Closed rows above the still-open ones: the case the row-list jump exists for.
// Under newest-first order (REQ-318) the newest open REQ is usually row 0 and the
// jump is a no-op, so the fixture deliberately puts the open work lower — an old
// REQ still running under newer finished ones is what makes the second movement
// do anything.
var rows = [
  { waitOpen: false, workOpen: false },
  { waitOpen: false, workOpen: false },
  { waitOpen: false, workOpen: false },
  { waitOpen: true, workOpen: false },
  { waitOpen: false, workOpen: true }
];
var scrollHostStub = { scrollTop: 0 };
// The Now button's three steps, in the order the handler runs them: take the
// window, let the rows follow it, then scroll among THOSE rows.
//
// The two row sets below are what makes this an ORDERING assertion rather than a
// restatement. The rows array is what was on screen before the jump;
// rowsAfterJump is what the window the jump chose admits — a narrower set whose
// first open row sits at a different index. Deciding the scroll before the
// refresh, which is what timelineNowJump used to do internally, yields 3;
// deciding it after yields 1. Only the second is right, and only a fixture where
// the two disagree can tell them apart.
var rowsAfterJump = [
  { waitOpen: false, workOpen: false },
  { waitOpen: true, workOpen: false },
  { waitOpen: false, workOpen: true }
];
var scrollTopIfDecidedBeforeRefresh = timelineFirstOpenRowIndex(rows) * TIMELINE_ROW_HEIGHT;
var nowWindow = timelineNowJump(nowMs, queueEndMs, boundStart, boundEnd);
var nowWindowWithNothingScheduled = timelineNowJump(nowMs, nowMs, boundStart, boundEnd);
var nowWindowWithNoForecast = timelineNowJump(nowMs, NaN, boundStart, boundEnd);
var openRowIndex = timelineFirstOpenRowIndex(rowsAfterJump);
if (openRowIndex >= 0) {
  scrollHostStub.scrollTop = openRowIndex * TIMELINE_ROW_HEIGHT;
}

process.stdout.write(JSON.stringify({
  nowWindowIso:
    new Date(nowWindow.windowStartMs).toISOString() + " → " + new Date(nowWindow.windowEndMs).toISOString(),
  nothingScheduledSpanMs:
    nowWindowWithNothingScheduled.windowEndMs - nowWindowWithNothingScheduled.windowStartMs,
  nothingScheduledHoldsNow:
    nowMs >= nowWindowWithNothingScheduled.windowStartMs &&
    nowMs <= nowWindowWithNothingScheduled.windowEndMs,
  noForecastSpanMs: nowWindowWithNoForecast.windowEndMs - nowWindowWithNoForecast.windowStartMs,
  noForecastHoldsNow:
    nowMs >= nowWindowWithNoForecast.windowStartMs && nowMs <= nowWindowWithNoForecast.windowEndMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  nowInsideWindow: nowMs >= nowWindow.windowStartMs && nowMs <= nowWindow.windowEndMs,
  queueEndInsideWindow: queueEndMs >= nowWindow.windowStartMs && queueEndMs <= nowWindow.windowEndMs,
  scrollTop: scrollHostStub.scrollTop,
  wantScrollTop: 1 * TIMELINE_ROW_HEIGHT,
  scrollTopIfDecidedBeforeRefresh: scrollTopIfDecidedBeforeRefresh
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline now jump", javascriptProbe)
	var nowJumpResult struct {
		NowWindowIso                    string  `json:"nowWindowIso"`
		NothingScheduledSpanMs          float64 `json:"nothingScheduledSpanMs"`
		NothingScheduledHoldsNow        bool    `json:"nothingScheduledHoldsNow"`
		NoForecastSpanMs                float64 `json:"noForecastSpanMs"`
		NoForecastHoldsNow              bool    `json:"noForecastHoldsNow"`
		MinSpanMs                       float64 `json:"minSpanMs"`
		NowInsideWindow                 bool    `json:"nowInsideWindow"`
		QueueEndInsideWindow            bool    `json:"queueEndInsideWindow"`
		ScrollTop                       float64 `json:"scrollTop"`
		WantScrollTop                   float64 `json:"wantScrollTop"`
		ScrollTopIfDecidedBeforeRefresh float64 `json:"scrollTopIfDecidedBeforeRefresh"`
	}
	if decodeError := json.Unmarshal(probeOutput, &nowJumpResult); decodeError != nil {
		t.Fatalf("decode timeline now-jump behavior: %v (output %q)", decodeError, probeOutput)
	}

	// Now on a drained queue: a window, not the zoom floor. This is the state the
	// button lands in whenever the forecast has nothing left to schedule, which on a
	// healthy queue is most of the time.
	for _, degenerate := range []struct {
		name     string
		spanMs   float64
		holdsNow bool
	}{
		{"the queue-empty instant equal to now", nowJumpResult.NothingScheduledSpanMs, nowJumpResult.NothingScheduledHoldsNow},
		{"no forecast at all", nowJumpResult.NoForecastSpanMs, nowJumpResult.NoForecastHoldsNow},
	} {
		if degenerate.spanMs <= nowJumpResult.MinSpanMs {
			t.Errorf("with %s, Now lands on a %.0f ms window against a %.0f ms zoom floor; at or "+
				"below the floor there is nowhere left to zoom and no context around the now-line",
				degenerate.name, degenerate.spanMs, nowJumpResult.MinSpanMs)
		}
		if !degenerate.holdsNow {
			t.Errorf("with %s, the window Now lands on does not contain the now-line", degenerate.name)
		}
	}

	if !nowJumpResult.NowInsideWindow || !nowJumpResult.QueueEndInsideWindow {
		t.Fatalf("the Now window %s does not cover both the now-line and the projected queue end", nowJumpResult.NowWindowIso)
	}
	if nowJumpResult.ScrollTop != nowJumpResult.WantScrollTop {
		t.Fatalf("Now left the row list at scrollTop %.0f, want %.0f — the first still-open row",
			nowJumpResult.ScrollTop, nowJumpResult.WantScrollTop)
	}
	// The ordering half, and the reason the fixture carries two row sets. Deciding
	// the scroll from the PRE-jump rows — which is what timelineNowJump did before
	// REQ-319 split it — lands somewhere else entirely. If these two ever agree the
	// fixture has stopped being able to tell the orders apart, and the assertion
	// above has quietly become a restatement.
	if nowJumpResult.ScrollTopIfDecidedBeforeRefresh == nowJumpResult.WantScrollTop {
		t.Fatalf("the fixture's pre-jump and post-jump row sets both give scrollTop %.0f, so this "+
			"test cannot tell a scroll decided before the row refresh from one decided after; "+
			"give the two sets different first-open indices",
			nowJumpResult.WantScrollTop)
	}
}

func TestJavaScriptBehaviorCalendarDayBreakdownGroupsStatuses(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := sliceBalancedBlockAfter(t, indexHtml, "function calendarDayBreakdown(") + `
var entries = [
  { id: "REQ-1", status: "cancelled" },
  { id: "REQ-2", status: "completed" },
  { id: "REQ-3", status: "blocked-archive-collision" },
  { id: "REQ-4", status: "completed" },
  { id: "REQ-5", status: "blocked" },
  { id: "REQ-6", status: "completed-with-issues" },
  { id: "REQ-7", status: "blockd-dependency-cycle" },
  { id: "REQ-8", status: "claimed" },
  { id: "REQ-9" }
];
process.stdout.write(JSON.stringify(calendarDayBreakdown(entries)));`
	probeOutput := runJavaScriptBehaviorProbe(t, "calendar day breakdown", javascriptProbe)

	var breakdown []struct {
		Group string `json:"group"`
		Label string `json:"label"`
		Count int    `json:"count"`
	}
	if decodeError := json.Unmarshal(probeOutput, &breakdown); decodeError != nil {
		t.Fatalf("decode calendar day breakdown: %v (output %q)", decodeError, probeOutput)
	}
	want := []struct {
		group string
		count int
	}{
		{"done", 2},
		{"with-issues", 1},
		{"claimed", 1},
		{"blocked", 2},      // `blocked` + `blocked-archive-collision`, one group
		{"cancelled", 1},    // never folded into done
		{"unrecognized", 2}, // the typo'd status and the one with no status at all
	}
	if len(breakdown) != len(want) {
		t.Fatalf("breakdown = %#v, want %d non-zero groups (empty groups must not render)", breakdown, len(want))
	}
	for index, wantPart := range want {
		if breakdown[index].Group != wantPart.group || breakdown[index].Count != wantPart.count {
			t.Fatalf("breakdown[%d] = %s×%d, want %s×%d (fixed group order, exact status matching)",
				index, breakdown[index].Group, breakdown[index].Count, wantPart.group, wantPart.count)
		}
	}
}

func TestJavaScriptBehaviorTimelineFallbackBoundsSpanTheWholeMatchedSet(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	// rangeStart is deliberately unparseable, which is what takes the fallback. The
	// rows span four hours, and they are in the newest-first order the producer
	// emits, so a fallback anchored on [0] lands at the NEWEST end.
	// REQ-933's WORK RUNS EIGHT HOURS PAST the newest created_at on purpose. An
	// extent taken from created_at alone would end the window at 12:00 and clip that
	// bar off the right edge while still listing the row, so the assertions below
	// would pass; naming the window in the summary is what makes the difference
	// visible.
	brokenRangePayload := `{
	  "now": "2026-08-18T21:00:00Z",
	  "rangeStart": "not-a-timestamp",
	  "rangeEnd": "2026-08-18T21:00:00Z",
	  "rows": [
	    {"id":"REQ-933","createdTime":"2026-08-18T12:00:00Z","claimedTime":"2026-08-18T12:10:00Z",
	     "completedTime":"2026-08-18T20:00:00Z","waitMinutes":10,"workMinutes":470,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-932","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:10:00Z",
	     "completedTime":"2026-08-18T10:30:00Z","waitMinutes":10,"workMinutes":20,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-931","createdTime":"2026-08-18T08:00:00Z","claimedTime":"2026-08-18T08:10:00Z",
	     "completedTime":"2026-08-18T08:30:00Z","waitMinutes":10,"workMinutes":20,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false}
	  ]
	}`

	// A payload whose range is unreadable AND whose only row is still OPEN, with no
	// projection to push the bound past the now-line. timelineRowSegments draws that
	// row's wait to `now`, so an extent built from stored stamps alone stops eight
	// hours short of it and the live part of the bar is unreachable in every window.
	// (Found by Codex on the pull request.)
	openRowPayload := `{
	  "now": "2026-08-18T20:00:00Z",
	  "rangeStart": "not-a-timestamp",
	  "rangeEnd": "not-a-timestamp",
	  "rows": [
	    {"id":"REQ-951","createdTime":"2026-08-18T12:00:00Z","claimedTime":null,
	     "completedTime":null,"waitMinutes":480,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ]
	}`

	// And a payload where no row carries a readable instant at all: the fallback
	// must decline rather than invent a window. Doubly unreachable from the producer
	// — timelineRange always returns real instants, and buildTimelineAggregate drops
	// a ticket whose created_at does not parse — but a payload-integrity failure has
	// to end in a legible message rather than a window built from NaNs, which is the
	// same posture timelineRowSegments takes for its own unreachable case.
	unreadablePayload := `{
	  "now": "2026-08-18T13:00:00Z",
	  "rangeStart": "not-a-timestamp",
	  "rangeEnd": "also-not-a-timestamp",
	  "rows": [
	    {"id":"REQ-941","createdTime":"not-a-timestamp","claimedTime":null,
	     "completedTime":null,"waitMinutes":0,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ]
	}`

	probeDriver := `
function drawnRowIds() {
  var ids = [];
  (function walk(node) {
    (node.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      if (childNode.stubName === "g" && attributes["data-detail-id"]) { ids.push(attributes["data-detail-id"]); return; }
      walk(childNode);
    });
  })(timelineStubHosts["timeline-scroll"]);
  return ids;
}
renderTimelineView();
process.stdout.write(JSON.stringify({
  rowIds: drawnRowIds(),
  summary: timelineStubHosts["timeline-summary"].textContent
}));
`

	renderWith := func(payload string) (drawnIds []string, summary string) {
		t.Helper()
		javascriptProbe := timelineRenderDomStubPreamble +
			"var boardData = { timeline: " + payload + " };\n" +
			string(rendererFragment) +
			probeDriver
		probeOutput := runJavaScriptBehaviorProbe(t, "timeline fallback bounds", javascriptProbe)
		var result struct {
			RowIds  []string `json:"rowIds"`
			Summary string   `json:"summary"`
		}
		if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
			t.Fatalf("decode timeline fallback bounds behavior: %v (output %q)", decodeError, probeOutput)
		}
		return result.RowIds, result.Summary
	}

	drawnIds, summary := renderWith(brokenRangePayload)

	// SETUP, ASSERTED: if the payload did not take the fallback branch, everything
	// below passes for the wrong reason.
	if !strings.Contains(summary, "3 REQs in the window") {
		t.Fatalf("the fallback render summarised %q; want all three rows, which is what says the "+
			"bounds cover the whole matched set rather than one hour around the newest row", summary)
	}
	if len(drawnIds) != 3 {
		t.Fatalf("the fallback render drew %v, want all three rows — the old fallback bounded the "+
			"view to one hour around REQ-933, the newest, leaving the other two unreachable", drawnIds)
	}
	// The oldest row specifically. A fallback anchored on the newest capture leaves
	// exactly this one off the chart.
	oldestIsDrawn := false
	for _, drawnId := range drawnIds {
		if drawnId == "REQ-931" {
			oldestIsDrawn = true
		}
	}
	if !oldestIsDrawn {
		t.Errorf("the oldest row REQ-931 is not on the chart (drawn: %v); the fallback bounds have "+
			"to reach the earliest instant the matched set carries", drawnIds)
	}
	// EVERY instant the rows carry, not just created_at. REQ-933's work ends at
	// 20:00, eight hours past the newest capture; an extent taken from created_at
	// alone ends the window around 12:04 and clips that bar off the frame while
	// still listing its row.
	if !strings.Contains(summary, "→ 2026-08-18 20:") {
		t.Errorf("the fallback window is %q; it has to reach REQ-933's completion at 20:00, so the "+
			"extent must read claimed and completed instants and not created_at alone", summary)
	}

	// AN OPEN ROW'S NOW-LINE IS PART OF THE EXTENT. The window has to reach 20:00,
	// where the bar is drawn to, not stop at the 12:00 the row has stored.
	openRowIds, openRowSummary := renderWith(openRowPayload)
	if len(openRowIds) != 1 {
		t.Fatalf("the open-row fallback render drew %v, want the single fixture row", openRowIds)
	}
	if !strings.Contains(openRowSummary, "→ 2026-08-18 20:") {
		t.Errorf("the fallback window for a still-open row is %q; it has to reach the now-line at "+
			"20:00, which is where timelineRowSegments draws that row's wait to — bounds are what "+
			"every control clamps against, so anything short of it is unreachable", openRowSummary)
	}

	// And with nothing readable, the view says so instead of fabricating a window.
	_, unreadableSummary := renderWith(unreadablePayload)
	if !strings.Contains(unreadableSummary, "nothing to place on a timeline") {
		t.Errorf("with no readable range and no readable row instants the summary reads %q; want the "+
			"existing decline rather than an invented window", unreadableSummary)
	}
}

func TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-901","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-902","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-903","createdTime":"2026-08-18T11:00:00Z","claimedTime":null,
	     "completedTime":null,"waitMinutes":60,"workMinutes":0,
	     "waitOpen":true,"workOpen":false,"hasWork":false,"anomaly":false}
	  ],
	  "projection": {
	    "confident": true,
	    "chainStart": "2026-08-18T12:00:00Z",
	    "queueEnd": "2026-08-18T14:30:00Z",
	    "windowSamples": 60, "windowSize": 60, "minimumSamples": 5,
	    "normalSamples": 55, "normalMinutes": 40,
	    "trivialSamples": 5, "trivialMinutes": 10,
	    "rows": [{"id":"REQ-903","startTime":"2026-08-18T12:00:00Z","endTime":"2026-08-18T12:40:00Z"}],
	    "excluded": [{"id":"REQ-904","reason":"waiting on an external condition"}],
	    "queueEndSource": "median"
	  }
	}`

	// One render per filter state, each from a fresh stub, so nothing carries over.
	probeDriver := `
function renderWithFilter(visibleIds) {
  [
    "timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
    "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
  ].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
  timelineStubVisibleIds = visibleIds;
  renderTimelineView();
  return {
    summary: timelineStubHosts["timeline-summary"].textContent || "",
    forecast: collectStubText(timelineStubHosts["timeline-forecast"]),
    excluded: collectStubText(timelineStubHosts["timeline-excluded"])
  };
}
function collectStubText(node) {
  var text = node.textContent || "";
  (node.children || []).forEach(function (child) { text += " " + collectStubText(child); });
  return text;
}
// A DRAINED queue, which is the only state that reaches the chainCount === 0
// branch. The whole board's projection is replaced rather than a second fixture
// added, because the contradiction is between two clauses of ONE paragraph and
// both have to be produced by one render.
function renderDrainedWithFilter(visibleIds) {
  boardData.timeline.projection.rows = [];
  return renderWithFilter(visibleIds);
}
process.stdout.write(JSON.stringify({
  unfiltered: renderWithFilter(null),
  filtered: renderWithFilter(["REQ-901"]),
  drainedUnfiltered: renderDrainedWithFilter(null),
  drainedFiltered: renderDrainedWithFilter(["REQ-901"])
}));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline forecast filter label", javascriptProbe)

	type renderedView struct {
		Summary  string `json:"summary"`
		Forecast string `json:"forecast"`
		Excluded string `json:"excluded"`
	}
	var rendered struct {
		Unfiltered        renderedView `json:"unfiltered"`
		Filtered          renderedView `json:"filtered"`
		DrainedUnfiltered renderedView `json:"drainedUnfiltered"`
		DrainedFiltered   renderedView `json:"drainedFiltered"`
	}
	if decodeError := json.Unmarshal(probeOutput, &rendered); decodeError != nil {
		t.Fatalf("decode rendered timeline views: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	// The fixture has to actually produce the disagreement, or the assertions
	// below pass against a view that was never filtered.
	if !strings.Contains(rendered.Unfiltered.Summary, "3 REQs") {
		t.Fatalf("the unfiltered render must show all three fixture rows, got summary %q", rendered.Unfiltered.Summary)
	}
	if !strings.Contains(rendered.Filtered.Summary, "1 REQ ") {
		t.Fatalf("the filtered render must show one row, got summary %q", rendered.Filtered.Summary)
	}

	if strings.Contains(rendered.Unfiltered.Forecast, "whole queue") {
		t.Errorf("an unfiltered view has nothing to disambiguate, so the forecast must carry no label; got %q",
			rendered.Unfiltered.Forecast)
	}
	if !strings.Contains(rendered.Filtered.Forecast, "whole queue") {
		t.Errorf("a filtered view forecasts the whole queue and must say so — this is the wiring, not the copy; got %q",
			rendered.Filtered.Forecast)
	}
	if !strings.Contains(rendered.Filtered.Excluded, "whole queue") {
		t.Errorf("the excluded list names REQ-904, which no visible row carries, and must say whose queue it lists; got %q",
			rendered.Filtered.Excluded)
	}
	// Both renders still carry the estimate: the label is added, not substituted.
	for viewName, view := range map[string]renderedView{
		"unfiltered": rendered.Unfiltered,
		"filtered":   rendered.Filtered,
	} {
		if !strings.Contains(view.Forecast, "Queue empties around") {
			t.Errorf("the %s forecast lost its estimate: %q", viewName, view.Forecast)
		}
	}

	// A DRAINED QUEUE, where the paragraph used to contradict itself inside one
	// sentence pair: "This covers the whole queue, not the rows shown." followed by
	// "Nothing left to schedule — every remaining REQ is listed below.", above a
	// single row, with the excluded paragraph under it naming a REQ that was not
	// listed anywhere.
	//
	// SETUP, ASSERTED: without this the two assertions below pass against a forecast
	// that never took the drained branch at all.
	if !strings.Contains(rendered.DrainedFiltered.Forecast, "Nothing left to schedule") {
		t.Fatalf("the drained fixture did not reach the nothing-left branch: %q",
			rendered.DrainedFiltered.Forecast)
	}
	if strings.Contains(rendered.DrainedFiltered.Forecast, "listed below") {
		t.Errorf("the forecast claims every remaining REQ is \"listed below\" while also saying it "+
			"covers the whole queue and not the rows shown, above one row: %q",
			rendered.DrainedFiltered.Forecast)
	}
	if !strings.Contains(rendered.DrainedFiltered.Forecast, "whole queue") {
		t.Errorf("the drained filtered forecast dropped the whole-queue note, which is the half "+
			"that is TRUE and is what the figures under it depend on: %q",
			rendered.DrainedFiltered.Forecast)
	}
	// Unfiltered, the sentence is accurate and must be left alone: the rows on
	// screen really are the whole queue there.
	if !strings.Contains(rendered.DrainedUnfiltered.Forecast, "listed below") {
		t.Errorf("with nothing filtered the rows ARE the whole queue, so the forecast should still "+
			"say so: %q", rendered.DrainedUnfiltered.Forecast)
	}
}

func TestJavaScriptBehaviorClipboardTitleSplicesPreserveMarkdownStructure(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceDeclarationAfter(t, indexHtml, "var inlineTicketTitleMaxLength ="),
		sliceDeclarationAfter(t, indexHtml, "var referencedTicketsGlossaryHeading ="),
		sliceBalancedBlockAfter(t, indexHtml, "function describeRequestStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function recordReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateTicketMentions("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildReferencedTicketsGlossary("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateClipboardPayload("),
	}
	cases := []struct{ name, title, wantShort string }{
		{"pipe", "Split the row | keep the pipe", "Split the row | keep the pipe"},
		{"backslash pipe", `Preserve \| and \\| and \\\| literally`, `Preserve \| and \\| and \\\| literally`},
		{"single cut", "Keep " + strings.Repeat("word ", 8) + "`command with many additional arguments and its close`", "Keep " + strings.Repeat("word ", 8) + "command with…"},
		{"double cut", "Keep " + strings.Repeat("word ", 8) + "``command with many additional arguments and its close``", "Keep " + strings.Repeat("word ", 8) + "command with…"},
		{"balanced code pipe", "Keep `left | right` readable", "Keep left | right readable"},
		{"code emphasis", "Keep `*important*` readable", "Keep *important* readable"},
		{"unmatched emphasis", "Keep `*important` readable", "Keep *important readable"},
		{"code link", "Keep `[label](target)` literal", "Keep [label](target) literal"},
		{"code entity", "Keep `&copy;` literal", "Keep &amp;copy; literal"},
	}
	bodies := []string{
		"| Reference | Unchanged | Last |\n| --- | --- | --- |\n| REQ-1108 | `author code` | final cell |\n",
		"Read REQ-1108, then `author code` and ``double code`` and final prose closing*.\n",
	}
	// Compare actual renderer structure, not a count of source delimiters: GFM
	// silently discards surplus cells, and an even backtick count can still be
	// an unmatched double-backtick delimiter.
	structurePattern := regexp.MustCompile(`</?[a-z][^>]*>`)
	for _, testCase := range cases {
		for bodyIndex, body := range bodies {
			t.Run(fmt.Sprintf("%s/%d", testCase.name, bodyIndex), func(t *testing.T) {
				mentions, encodeError := json.Marshal(collectDocumentTicketMentions(body, newCitationFixtureResolver()))
				if encodeError != nil {
					t.Fatal(encodeError)
				}
				probe := `var requestsById = {"REQ-1108": {title: ` + mustMarshalJSONString(t, testCase.title) + `, status: "pending"}}; var userRequestsById = {};` +
					strings.Join(functionBlocks, "\n") + `
process.stdout.write(JSON.stringify(annotateClipboardPayload([{text: ` + mustMarshalJSONString(t, body) + `, ticketMentions: ` + string(mentions) + `}], [])));`
				var payload string
				if decodeError := json.Unmarshal(runJavaScriptBehaviorProbe(t, "safe clipboard title", probe), &payload); decodeError != nil {
					t.Fatal(decodeError)
				}
				appendix := "\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n- REQ-1108 — " + testCase.title + " (pending)\n"
				if !strings.HasSuffix(payload, appendix) {
					t.Fatalf("full original appendix title changed: %q", payload)
				}
				annotatedBody := strings.TrimSuffix(payload, appendix)
				if !strings.Contains(annotatedBody, "REQ-1108 (-> ") {
					t.Fatalf("title expansion was suppressed: %q", annotatedBody)
				}
				originalHTML, renderError := renderMarkdownBodyToHtml(body)
				if renderError != nil {
					t.Fatal(renderError)
				}
				pastedHTML, renderError := renderMarkdownBodyToHtml(annotatedBody)
				if renderError != nil {
					t.Fatal(renderError)
				}
				if !strings.Contains(pastedHTML, "REQ-1108 (-&gt; "+testCase.wantShort+")") {
					t.Errorf("short title lost literal text: want %q in %s", testCase.wantShort, pastedHTML)
				}
				if !reflect.DeepEqual(structurePattern.FindAllString(originalHTML, -1), structurePattern.FindAllString(pastedHTML, -1)) {
					t.Errorf("splice changed rendered block/inline structure:\noriginal %s\npasted %s", originalHTML, pastedHTML)
				}
				if !strings.Contains(pastedHTML, "<code>author code</code>") {
					t.Errorf("title consumed author's code span: %s", pastedHTML)
				}
				if bodyIndex == 0 && !strings.Contains(pastedHTML, "<td>final cell</td>") {
					t.Errorf("title displaced a table cell: %s", pastedHTML)
				}
				if bodyIndex == 1 && !strings.Contains(pastedHTML, " and final prose closing*.") {
					t.Errorf("title consumed following prose: %s", pastedHTML)
				}
			})
		}
	}
}

// TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests drives the
// shipped renderActivity over a payload where one REQ owns two rows, which is
// the shape REQ-572 introduced. The summary has to report both numbers,
// because a transition count alone reads as a REQ count on a surface that used
// to be one row per REQ, and the empty states have to count the same unit.
func TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		// REQ-573 gave the REQ cell a detail button and the rows a selected
		// class, so renderActivity now reaches createElement and the selection
		// helpers; the counts this test is about are unchanged.
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function activityRowsWithin("),
		sliceBalancedBlockAfter(t, indexHtml, "function activityWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function selectedActivityRequestId("),
		sliceBalancedBlockAfter(t, indexHtml, "function applyActivitySelectionHighlight("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderActivity("),
	}
	javascriptProbe := `
// A stub node whose textContent setter drops its children, as the real one
// does: renderActivity clears the table body that way between renders, and a
// stub that kept them would let row counts accumulate across the cases below.
function makeStubNode() {
  var node = {
    childNodes: [],
    attributes: {},
    hidden: false,
    scope: "",
    stubText: "",
    classList: { toggle: function () {} },
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = attributeValue; },
    getAttribute: function (attributeName) {
      return Object.prototype.hasOwnProperty.call(this.attributes, attributeName) ? this.attributes[attributeName] : null;
    },
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
  Object.defineProperty(node, "textContent", {
    get: function () { return this.stubText; },
    set: function (nodeText) { this.stubText = nodeText; this.childNodes = []; }
  });
  return node;
}
var nodesById = {};
var document = {
  getElementById: function (nodeId) {
    if (!nodesById[nodeId]) { nodesById[nodeId] = makeStubNode(); }
    return nodesById[nodeId];
  },
  createElement: function (tagName) { var node = makeStubNode(); node.stubTag = tagName; return node; }
};
var viewState = { activityWindowHours: 24 };
// No drawer is open in any of these cases, so no row is ever selected here;
// TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest owns
// the selection behavior these two carry.
var currentDetailKind = "";
var currentDetailId = "";
var requestsById = {
  "REQ-800": { title: "Busy request", status: "completed" },
  "REQ-801": { title: "Quiet request", status: "claimed" }
};
var hiddenRequestIds = {};
function requestMatchesFilters(requestId) { return !hiddenRequestIds[requestId]; }
function makeInstantWithRelativeNode() { return null; }
function hoursAgo(hourCount) { return new Date(Date.now() - hourCount * 3600 * 1000).toISOString(); }
var threeTransitions = [
  { id: "REQ-800", stampField: "completed_at", stampAt: hoursAgo(1), transition: "completed" },
  { id: "REQ-801", stampField: "claimed_at", stampAt: hoursAgo(2), transition: "claimed" },
  { id: "REQ-800", stampField: "created_at", stampAt: hoursAgo(3), transition: "captured" }
];
var boardData = { activity: threeTransitions };
` + strings.Join(functionBlocks, "\n") + `
function renderCase() {
  nodesById = {};
  renderActivity();
  var tableBody = nodesById["activity-table-body"];
  return {
    summary: nodesById["activity-summary"].textContent,
    rowCount: tableBody.childNodes.length,
    requestAttributes: tableBody.childNodes.map(function (tableRow) { return tableRow.attributes["data-activity-request"]; }),
    emptyHidden: nodesById["activity-empty"].hidden,
    emptyText: nodesById["activity-empty"].textContent
  };
}
var results = {};
results.everyTransition = renderCase();
boardData.activity = [threeTransitions[0]];
results.singleTransition = renderCase();
boardData.activity = threeTransitions;
hiddenRequestIds = { "REQ-801": true };
results.oneRequestFiltered = renderCase();
hiddenRequestIds = { "REQ-800": true, "REQ-801": true };
results.everythingFiltered = renderCase();
hiddenRequestIds = {};
boardData.activity = [];
results.nothingInWindow = renderCase();
process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "activity summary counts", javascriptProbe)

	type activityRenderResult struct {
		Summary           string   `json:"summary"`
		RowCount          int      `json:"rowCount"`
		RequestAttributes []string `json:"requestAttributes"`
		EmptyHidden       bool     `json:"emptyHidden"`
		EmptyText         string   `json:"emptyText"`
	}
	var results struct {
		EveryTransition    activityRenderResult `json:"everyTransition"`
		SingleTransition   activityRenderResult `json:"singleTransition"`
		OneRequestFiltered activityRenderResult `json:"oneRequestFiltered"`
		EverythingFiltered activityRenderResult `json:"everythingFiltered"`
		NothingInWindow    activityRenderResult `json:"nothingInWindow"`
	}
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode activity render results: %v (output %q)", decodeError, probeOutput)
	}

	if results.EveryTransition.Summary != "3 transitions across 2 REQs in the last 24 hours" {
		t.Fatalf("summary for three transitions of two REQs = %q, want both counts", results.EveryTransition.Summary)
	}
	if results.EveryTransition.RowCount != 3 {
		t.Fatalf("rendered rows = %d, want 3 — one per transition, not one per REQ", results.EveryTransition.RowCount)
	}
	// REQ-573 highlights a REQ's sibling rows through this attribute, so the
	// repeated id is the contract rather than an oversight.
	wantAttributes := []string{"REQ-800", "REQ-801", "REQ-800"}
	if !reflect.DeepEqual(results.EveryTransition.RequestAttributes, wantAttributes) {
		t.Fatalf("data-activity-request values = %#v, want %#v (one REQ's rows repeat the id)",
			results.EveryTransition.RequestAttributes, wantAttributes)
	}
	if !results.EveryTransition.EmptyHidden {
		t.Fatalf("empty state showed while three transitions rendered")
	}
	if results.SingleTransition.Summary != "1 transition across 1 REQ in the last 24 hours" {
		t.Fatalf("singular summary = %q", results.SingleTransition.Summary)
	}
	if results.OneRequestFiltered.Summary != "2 transitions across 1 REQ in the last 24 hours (3 before filters)" {
		t.Fatalf("filtered summary = %q, want the filtered counts plus the transitions before filtering", results.OneRequestFiltered.Summary)
	}
	if results.EverythingFiltered.EmptyHidden ||
		results.EverythingFiltered.EmptyText != "3 transitions happened in this window, but the active filters hide all of them." {
		t.Fatalf("filters-hid-everything empty state = %q (hidden=%v)", results.EverythingFiltered.EmptyText, results.EverythingFiltered.EmptyHidden)
	}
	if results.NothingInWindow.EmptyHidden ||
		results.NothingInWindow.EmptyText != "No lifecycle transition falls inside the last 24 hours." {
		t.Fatalf("nothing-moved empty state = %q (hidden=%v)", results.NothingInWindow.EmptyText, results.NothingInWindow.EmptyHidden)
	}
}

// The Verify Findings strip sits outside the view panels, so every view switch
// used to leave it on screen — including the Activity view, where it only pushed
// the transitions table down (REQ-578). The rule belongs to the view switch, and
// it must not turn an empty strip into a visible one on the other views, so both
// halves are driven here through the real applyView.
func TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderVerifyFindingsStrip("),
		// The renderer's own helpers (REQ-579): the strip is a list of rows now,
		// and applyView still has to see what that renderer drew.
		sliceBalancedBlockAfter(t, indexHtml, "function formatFindingsSummary("),
		sliceBalancedBlockAfter(t, indexHtml, "function groupFindingsBySubject("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeFindingRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeSkippedProbeRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function applyView("),
	}
	javascriptProbe := `
function makeStubNode() {
  var node = {
    children: [],
    attributes: {},
    className: "",
    hidden: false,
    stubText: "",
    classList: { toggle: function () {} },
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = attributeValue; },
    appendChild: function (childNode) { this.children.push(childNode); return childNode; }
  };
  Object.defineProperty(node, "textContent", {
    get: function () { return this.stubText; },
    set: function (nodeText) { this.stubText = nodeText; this.children = []; }
  });
  return node;
}
var nodesById = {};
var document = {
  getElementById: function (nodeId) {
    if (!nodesById[nodeId]) { nodesById[nodeId] = makeStubNode(); }
    return nodesById[nodeId];
  },
  createElement: function (tagName) { var node = makeStubNode(); node.stubTag = tagName; return node; }
};
var viewState = { view: "board", lens: "flat" };
var renderedOnce = {
  userRequestLens: true, calendar: true, durations: true,
  timeline: true, activity: true, testing: true
};
var boardData = {};
// applyView's neighbours are stubbed: this probe is about the findings strip,
// and every renderer it can reach needs a live board to say anything.
function hasActiveVisibleFilters() { return false; }
function applyLens() {}
function updateUserRequestActivityVisibility() {}
function renderCalendar() {}
function renderDurationsView() {}
function renderTimelineView() {}
function renderActivity() {}
function renderTestingView() {}
` + strings.Join(functionBlocks, "\n") + `
function stripHiddenOnView(viewName) {
  viewState.view = viewName;
  applyView();
  return nodesById["board-findings"].hidden;
}
var results = {};

nodesById = {};
boardData = {
  verifyFindings: [
    { category: "WORKTREE-MERGE-STATE-UNDETERMINED", detail: "a REQ-506 worktree", remedy: "inspect it" },
    { category: "WORKTREE-PRESENT-RUN-IN-FLIGHT", detail: "the REQ-570 worktree", fixable: true }
  ],
  verifySkipped: ["one probe could not run"]
};
renderVerifyFindingsStrip();
results.withFindingsAfterRender = nodesById["board-findings"].hidden;
results.withFindingsOnBoard = stripHiddenOnView("board");
results.withFindingsOnActivity = stripHiddenOnView("activity");
results.withFindingsBackOnBoard = stripHiddenOnView("board");
results.withFindingsOnTimeline = stripHiddenOnView("timeline");

nodesById = {};
boardData = { verifyFindings: [], verifySkipped: [] };
renderVerifyFindingsStrip();
results.emptyAfterRender = nodesById["board-findings"].hidden;
results.emptyOnBoard = stripHiddenOnView("board");
results.emptyOnActivity = stripHiddenOnView("activity");

// Findings the renderer never drew: only the skipped probes have anything to
// say, and the strip must still survive a trip through the Activity view.
nodesById = {};
boardData = { verifyFindings: [], verifySkipped: ["one probe could not run"] };
renderVerifyFindingsStrip();
results.skippedOnlyOnActivity = stripHiddenOnView("activity");
results.skippedOnlyBackOnBoard = stripHiddenOnView("board");

process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "activity view hides the verify findings strip", javascriptProbe)

	var results struct {
		WithFindingsAfterRender bool `json:"withFindingsAfterRender"`
		WithFindingsOnBoard     bool `json:"withFindingsOnBoard"`
		WithFindingsOnActivity  bool `json:"withFindingsOnActivity"`
		WithFindingsBackOnBoard bool `json:"withFindingsBackOnBoard"`
		WithFindingsOnTimeline  bool `json:"withFindingsOnTimeline"`
		EmptyAfterRender        bool `json:"emptyAfterRender"`
		EmptyOnBoard            bool `json:"emptyOnBoard"`
		EmptyOnActivity         bool `json:"emptyOnActivity"`
		SkippedOnlyOnActivity   bool `json:"skippedOnlyOnActivity"`
		SkippedOnlyBackOnBoard  bool `json:"skippedOnlyBackOnBoard"`
	}
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode findings strip visibility results: %v (output %q)", decodeError, probeOutput)
	}

	if results.WithFindingsAfterRender || results.WithFindingsOnBoard {
		t.Fatalf("two findings rendered but the strip is hidden on the Board view (afterRender=%v, onBoard=%v)",
			results.WithFindingsAfterRender, results.WithFindingsOnBoard)
	}
	if !results.WithFindingsOnActivity {
		t.Fatalf("the Verify Findings strip is still visible on the Activity view")
	}
	if results.WithFindingsBackOnBoard || results.WithFindingsOnTimeline {
		t.Fatalf("the strip did not come back after the Activity view (board=%v, timeline=%v)",
			results.WithFindingsBackOnBoard, results.WithFindingsOnTimeline)
	}
	if !results.EmptyAfterRender || !results.EmptyOnBoard || !results.EmptyOnActivity {
		t.Fatalf("a strip with nothing to report became visible on a view switch (afterRender=%v, board=%v, activity=%v)",
			results.EmptyAfterRender, results.EmptyOnBoard, results.EmptyOnActivity)
	}
	if !results.SkippedOnlyOnActivity {
		t.Fatalf("a skipped-probes-only strip stayed visible on the Activity view")
	}
	if results.SkippedOnlyBackOnBoard {
		t.Fatalf("a skipped-probes-only strip did not come back on the Board view")
	}
}

// TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest drives
// the shipped renderActivity and the shipped selection helpers over the shape
// REQ-573 was raised for: two rows for one REQ and one row for another. Two
// halves are proved here. First, the REQ cell is a real button carrying the
// data-detail-* pair, which is the whole contract with the document-level
// delegation in board-controls.js — that is what opens the drawer, and adding a
// second opener is exactly what this REQ must not do. Second, the highlight is
// a SET: a REQ owns several rows, so clicking one marks all of them, a
// re-render restores them from the open drawer alone, and clearing the drawer
// clears them.
func TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function activityRowsWithin("),
		sliceBalancedBlockAfter(t, indexHtml, "function activityWindowPhrase("),
		sliceBalancedBlockAfter(t, indexHtml, "function selectedActivityRequestId("),
		sliceBalancedBlockAfter(t, indexHtml, "function applyActivitySelectionHighlight("),
		sliceBalancedBlockAfter(t, indexHtml, "function syncActivitySelectionToClick("),
		sliceBalancedBlockAfter(t, indexHtml, "function syncActivitySelectionToDrawer("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderActivity("),
	}
	javascriptProbe := `
// The textContent setter drops children, as the real one does — renderActivity
// clears the table body that way between renders. closest() supports only the
// [attribute] form the client actually passes: a stub that pretended to be a
// selector engine would keep passing after the selector changed.
function makeStubNode() {
  var node = {
    childNodes: [],
    attributes: {},
    presentClasses: {},
    hidden: false,
    scope: "",
    stubText: "",
    parentNode: null,
    setAttribute: function (attributeName, attributeValue) { this.attributes[attributeName] = attributeValue; },
    getAttribute: function (attributeName) {
      return Object.prototype.hasOwnProperty.call(this.attributes, attributeName) ? this.attributes[attributeName] : null;
    },
    appendChild: function (childNode) { childNode.parentNode = this; this.childNodes.push(childNode); return childNode; },
    closest: function (selector) {
      var attributeName = selector.slice(1, -1);
      if ("[" + attributeName + "]" !== selector) { throw new Error("unsupported selector " + selector); }
      var candidate = this;
      while (candidate) {
        if (Object.prototype.hasOwnProperty.call(candidate.attributes, attributeName)) { return candidate; }
        candidate = candidate.parentNode;
      }
      return null;
    }
  };
  node.classList = {
    toggle: function (className, shouldBePresent) {
      if (shouldBePresent) { node.presentClasses[className] = true; } else { delete node.presentClasses[className]; }
    },
    contains: function (className) { return !!node.presentClasses[className]; }
  };
  Object.defineProperty(node, "textContent", {
    get: function () { return this.stubText; },
    set: function (nodeText) { this.stubText = nodeText; this.childNodes = []; }
  });
  return node;
}
var nodesById = {};
var document = {
  getElementById: function (nodeId) {
    if (!nodesById[nodeId]) { nodesById[nodeId] = makeStubNode(); }
    return nodesById[nodeId];
  },
  createElement: function (tagName) { var node = makeStubNode(); node.stubTag = tagName; return node; }
};
var viewState = { activityWindowHours: 24 };
// The drawer's own state, which board-detail.js owns. The selection reads it
// rather than keeping a second copy, so the probe moves it the way the drawer
// does: openRequestDetail sets both, closeDrawer clears both.
var currentDetailKind = "";
var currentDetailId = "";
var requestsById = {
  "REQ-570": { title: "Delete the pending-heavy-testing status", status: "claimed" },
  "REQ-505": { title: "A quieter request", status: "completed" }
};
function requestMatchesFilters() { return true; }
function makeInstantWithRelativeNode() { return null; }
function hoursAgo(hourCount) { return new Date(Date.now() - hourCount * 3600 * 1000).toISOString(); }
var boardData = { activity: [
  { id: "REQ-570", stampField: "claimed_at", stampAt: hoursAgo(1), transition: "claimed" },
  { id: "REQ-505", stampField: "completed_at", stampAt: hoursAgo(2), transition: "completed" },
  { id: "REQ-570", stampField: "created_at", stampAt: hoursAgo(3), transition: "captured" }
] };
` + strings.Join(functionBlocks, "\n") + `
function currentRows() { return nodesById["activity-table-body"].childNodes; }
function reqCellButtonOf(tableRow) { return tableRow.childNodes[0].childNodes[0]; }
function selectionState() {
  return currentRows().map(function (tableRow) { return tableRow.classList.contains("is-activity-selected"); });
}
function readAcross(readOneRow) { return currentRows().map(readOneRow); }
var results = {};
renderActivity();
results.rowRequestIds = readAcross(function (tableRow) { return tableRow.getAttribute("data-activity-request"); });
results.reqCellButtonTags = readAcross(function (tableRow) { return reqCellButtonOf(tableRow).stubTag; });
results.reqCellButtonTypes = readAcross(function (tableRow) { return reqCellButtonOf(tableRow).type; });
results.reqCellButtonLabels = readAcross(function (tableRow) { return reqCellButtonOf(tableRow).textContent; });
results.detailKinds = readAcross(function (tableRow) { return reqCellButtonOf(tableRow).getAttribute("data-detail-kind"); });
results.detailIds = readAcross(function (tableRow) { return reqCellButtonOf(tableRow).getAttribute("data-detail-id"); });
results.beforeAnyClick = selectionState();

// A click on the first REQ-570 row. board-controls.js registers its delegation
// after this fragment, so the drawer state still names the PREVIOUS REQ while
// this listener runs — the probe moves it afterwards, as the delegation does.
syncActivitySelectionToClick({ target: reqCellButtonOf(currentRows()[0]) });
results.afterClickingTheFirstReq570Row = selectionState();
currentDetailKind = "req";
currentDetailId = "REQ-570";

// A window or filter change re-renders the table; the selection has to come
// back from the open drawer, with nothing remembered in between.
renderActivity();
results.afterReRender = selectionState();

syncActivitySelectionToClick({ target: reqCellButtonOf(currentRows()[1]) });
results.afterClickingTheReq505Row = selectionState();
currentDetailKind = "req";
currentDetailId = "REQ-505";

// A click that lands nowhere near the table re-reads the drawer, which is how
// the drawer's close button clears the highlight: closeDrawer runs first.
currentDetailKind = "";
currentDetailId = "";
syncActivitySelectionToClick({ target: makeStubNode() });
results.afterClosingTheDrawer = selectionState();

// Escape is the drawer's other close path and goes through the same re-read.
currentDetailKind = "req";
currentDetailId = "REQ-570";
syncActivitySelectionToDrawer();
results.afterReopeningWithoutAClick = selectionState();
currentDetailKind = "";
currentDetailId = "";
syncActivitySelectionToDrawer();
results.afterEscape = selectionState();

// A UR drawer is not a REQ selection: the id would collide with nothing here,
// but reading currentDetailId without its kind is the mistake worth pinning.
currentDetailKind = "ur";
currentDetailId = "REQ-570";
syncActivitySelectionToDrawer();
results.afterOpeningAUserRequestDrawer = selectionState();
process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "activity row selection", javascriptProbe)

	var results struct {
		RowRequestIds                  []string `json:"rowRequestIds"`
		ReqCellButtonTags              []string `json:"reqCellButtonTags"`
		ReqCellButtonTypes             []string `json:"reqCellButtonTypes"`
		ReqCellButtonLabels            []string `json:"reqCellButtonLabels"`
		DetailKinds                    []string `json:"detailKinds"`
		DetailIds                      []string `json:"detailIds"`
		BeforeAnyClick                 []bool   `json:"beforeAnyClick"`
		AfterClickingTheFirstReq570Row []bool   `json:"afterClickingTheFirstReq570Row"`
		AfterReRender                  []bool   `json:"afterReRender"`
		AfterClickingTheReq505Row      []bool   `json:"afterClickingTheReq505Row"`
		AfterClosingTheDrawer          []bool   `json:"afterClosingTheDrawer"`
		AfterReopeningWithoutAClick    []bool   `json:"afterReopeningWithoutAClick"`
		AfterEscape                    []bool   `json:"afterEscape"`
		AfterOpeningAUserRequestDrawer []bool   `json:"afterOpeningAUserRequestDrawer"`
	}
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode activity selection results: %v (output %q)", decodeError, probeOutput)
	}

	wantRowRequestIds := []string{"REQ-570", "REQ-505", "REQ-570"}
	if !reflect.DeepEqual(results.RowRequestIds, wantRowRequestIds) {
		t.Fatalf("data-activity-request values = %#v, want %#v", results.RowRequestIds, wantRowRequestIds)
	}
	// Keyboard reachability is the element, not a handler: a bare cell with a
	// click listener is unreachable by Tab, which is what the REQ rules out.
	if !reflect.DeepEqual(results.ReqCellButtonTags, []string{"button", "button", "button"}) ||
		!reflect.DeepEqual(results.ReqCellButtonTypes, []string{"button", "button", "button"}) {
		t.Fatalf("REQ cells are not keyboard-reachable buttons (tags %#v, types %#v)",
			results.ReqCellButtonTags, results.ReqCellButtonTypes)
	}
	if !reflect.DeepEqual(results.ReqCellButtonLabels, wantRowRequestIds) {
		t.Fatalf("REQ cell button labels = %#v, want the REQ ids %#v", results.ReqCellButtonLabels, wantRowRequestIds)
	}
	if !reflect.DeepEqual(results.DetailKinds, []string{"req", "req", "req"}) ||
		!reflect.DeepEqual(results.DetailIds, wantRowRequestIds) {
		t.Fatalf("the data-detail-* pair the board-controls.js delegation opens the drawer from is wrong (kinds %#v, ids %#v)",
			results.DetailKinds, results.DetailIds)
	}
	if !reflect.DeepEqual(results.BeforeAnyClick, []bool{false, false, false}) {
		t.Fatalf("rows rendered already selected with no drawer open: %#v", results.BeforeAnyClick)
	}
	wantReq570Selected := []bool{true, false, true}
	if !reflect.DeepEqual(results.AfterClickingTheFirstReq570Row, wantReq570Selected) {
		t.Fatalf("clicking one REQ-570 row selected %#v, want %#v — every row of that REQ, and no other REQ's row",
			results.AfterClickingTheFirstReq570Row, wantReq570Selected)
	}
	if !reflect.DeepEqual(results.AfterReRender, wantReq570Selected) {
		t.Fatalf("re-render dropped the selection for the open REQ: %#v, want %#v", results.AfterReRender, wantReq570Selected)
	}
	if !reflect.DeepEqual(results.AfterClickingTheReq505Row, []bool{false, true, false}) {
		t.Fatalf("clicking the other REQ left the first REQ's rows selected: %#v", results.AfterClickingTheReq505Row)
	}
	if !reflect.DeepEqual(results.AfterClosingTheDrawer, []bool{false, false, false}) {
		t.Fatalf("closing the drawer left rows highlighted: %#v", results.AfterClosingTheDrawer)
	}
	if !reflect.DeepEqual(results.AfterReopeningWithoutAClick, wantReq570Selected) {
		t.Fatalf("a drawer opened from somewhere other than the table did not mark its rows: %#v", results.AfterReopeningWithoutAClick)
	}
	if !reflect.DeepEqual(results.AfterEscape, []bool{false, false, false}) {
		t.Fatalf("Escape closed the drawer but left rows highlighted: %#v", results.AfterEscape)
	}
	if !reflect.DeepEqual(results.AfterOpeningAUserRequestDrawer, []bool{false, false, false}) {
		t.Fatalf("a UR drawer selected REQ rows by id alone: %#v", results.AfterOpeningAUserRequestDrawer)
	}
}

// REQ-579: the strip is a list of warnings, not a set of work items. A finding
// and a skipped probe are the same kind of thing to the reader — "verify has
// something to tell you" — so they share one row shape in one list, and the two
// visual languages that used to split them (a bordered card per finding, a
// bullet inside a collapsed disclosure per skipped probe) are gone.
//
// The weights come from the producer alone: `fixable` means `do-work cleanup`
// resolves it, and a skipped probe is a non-answer. Nothing here invents a
// severity the payload did not carry.
func TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function renderVerifyFindingsStrip("),
		sliceBalancedBlockAfter(t, indexHtml, "function formatFindingsSummary("),
		sliceBalancedBlockAfter(t, indexHtml, "function groupFindingsBySubject("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeFindingRow("),
		sliceBalancedBlockAfter(t, indexHtml, "function makeSkippedProbeRow("),
	}
	javascriptProbe := `
function makeStubNode(tagName) {
  var node = {
    stubTag: tagName,
    children: [],
    className: "",
    hidden: false,
    stubText: "",
    appendChild: function (childNode) { this.children.push(childNode); return childNode; }
  };
  Object.defineProperty(node, "textContent", {
    get: function () { return this.stubText; },
    set: function (nodeText) { this.stubText = nodeText; this.children = []; }
  });
  return node;
}
var createdTagNames = [];
var nodesById = {};
var document = {
  getElementById: function (nodeId) {
    if (!nodesById[nodeId]) { nodesById[nodeId] = makeStubNode("div"); }
    return nodesById[nodeId];
  },
  createElement: function (tagName) { createdTagNames.push(tagName); return makeStubNode(tagName); }
};
var boardData = {};
` + strings.Join(functionBlocks, "\n") + `
// Flatten a node's whole subtree so a row's class and its visible words can be
// read together — a row that carries the muted class but no text proves nothing.
function subtreeText(node) {
  if (node.children.length === 0) { return node.stubText; }
  return node.children.map(subtreeText).join(" ");
}
function describeChildren(node) {
  return node.children.map(function (childNode) {
    return {
      tag: childNode.stubTag,
      classes: childNode.className ? childNode.className.split(" ") : [],
      text: childNode.stubText,
      subtreeText: subtreeText(childNode)
    };
  });
}

// The subjectless finding is FIRST in producer order on purpose: grouping has to
// pull the two worktree rows above it, so the assertions below cannot pass on a
// renderer that simply echoes the payload order.
boardData = {
  verifyFindings: [
    {
      category: "CHECKPOINT-GHOST-REQUEST",
      detail: "do-work/CHECKPOINT.md names REQ-999, which exists nowhere",
      remedy: "edit that id out of the checkpoint"
    },
    {
      category: "WORKTREE-PRESENT-RUN-IN-FLIGHT",
      subject: "worktree-agent-REQ-506-focused-evidence",
      detail: "the worktree exists and its REQ is still in do-work/working/",
      remedy: "leave it in place"
    },
    {
      category: "MERGED-WORKTREE-LEFTOVER",
      subject: "worktree-agent-REQ-506-focused-evidence",
      detail: "the branch is already contained in HEAD",
      remedy: "cleanup Pass 5 removes it",
      fixable: true
    }
  ],
  verifySkipped: ["committed-queue-state probe for worktree-agent-REQ-506-focused-evidence: no such branch"]
};
renderVerifyFindingsStrip();

// The two hosts are display: contents, so their rows are one list on the page.
// Reading them back in that same order is what the reader sees. Read through
// getElementById, not the map: an id the renderer never touched must arrive as an
// empty node so the assertion below reports what is missing rather than dying on
// undefined.
var results = {
  stripHidden: document.getElementById("board-findings").hidden,
  headerCount: document.getElementById("board-findings-count").stubText,
  listChildren: describeChildren(document.getElementById("board-findings-cards"))
    .concat(describeChildren(document.getElementById("board-findings-skipped-list"))),
  createdTagNames: createdTagNames
};
process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "verify findings render as one row list", javascriptProbe)

	var results struct {
		StripHidden  bool   `json:"stripHidden"`
		HeaderCount  string `json:"headerCount"`
		ListChildren []struct {
			Tag         string   `json:"tag"`
			Classes     []string `json:"classes"`
			Text        string   `json:"text"`
			SubtreeText string   `json:"subtreeText"`
		} `json:"listChildren"`
		CreatedTagNames []string `json:"createdTagNames"`
	}
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode findings row results: %v (output %q)", decodeError, probeOutput)
	}

	if results.StripHidden {
		t.Fatal("the strip hid itself while it had two findings and a skipped probe to show")
	}
	// The skipped count joins the header only when there is one, so the reader
	// sees "checked, and here is what went unchecked" in one line.
	if results.HeaderCount != "3 findings · 1 probe not checked" {
		t.Errorf("header count = %q, want %q", results.HeaderCount, "3 findings · 1 probe not checked")
	}
	for _, createdTag := range results.CreatedTagNames {
		if createdTag == "details" || createdTag == "summary" {
			t.Errorf("the collapsed disclosure is still built (createElement(%q)) — skipped probes belong in the list", createdTag)
		}
	}

	// One list, in the order the reader sees it: the subject heading, its two
	// findings, the subjectless finding, then the skipped probe.
	if len(results.ListChildren) != 5 {
		t.Fatalf("the list holds %d children, want a subject heading plus four rows: %+v",
			len(results.ListChildren), results.ListChildren)
	}
	subjectHeading := results.ListChildren[0]
	if !hasClassName(subjectHeading.Classes, "board-findings-subject") {
		t.Fatalf("the list opens with %v, want the subject heading its two findings share", subjectHeading.Classes)
	}
	if subjectHeading.Text != "worktree-agent-REQ-506-focused-evidence" {
		t.Errorf("subject heading = %q, want the worktree both findings named", subjectHeading.Text)
	}

	rowMutedStates := []bool{}
	for _, listChild := range results.ListChildren[1:] {
		if hasClassName(listChild.Classes, "board-finding") {
			t.Errorf("a row still carries the old card class: %v", listChild.Classes)
		}
		if !hasClassName(listChild.Classes, "board-findings-row") {
			t.Fatalf("child %v is neither a row nor the one subject heading — the list must be flat", listChild.Classes)
		}
		rowMutedStates = append(rowMutedStates, hasClassName(listChild.Classes, "board-findings-row-muted"))
	}
	// Normal weight for the two findings a human must resolve; muted for the one a
	// command fixes and for the probe that never ran.
	wantMutedStates := []bool{false, true, false, true}
	if !reflect.DeepEqual(rowMutedStates, wantMutedStates) {
		t.Errorf("row muted states = %v, want %v (only fixable and skipped rows are muted)",
			rowMutedStates, wantMutedStates)
	}

	// Grouping reordered the payload: the subjectless finding came first and must
	// end up after the group, held away from it so it does not read as the last
	// heading's third row.
	subjectlessRow := results.ListChildren[3]
	if !strings.Contains(subjectlessRow.SubtreeText, "REQ-999") {
		t.Errorf("row 3 = %q, want the subjectless finding after the grouped ones", subjectlessRow.SubtreeText)
	}
	for _, detachedRow := range []int{3, 4} {
		if !hasClassName(results.ListChildren[detachedRow].Classes, "board-findings-row-detached") {
			t.Errorf("row %d joins the group above it: %v", detachedRow, results.ListChildren[detachedRow].Classes)
		}
	}

	fixableRow := results.ListChildren[2]
	if !strings.Contains(fixableRow.SubtreeText, "cleanup can fix") {
		t.Errorf("the fixable row lost its tag: %q", fixableRow.SubtreeText)
	}
	if !strings.Contains(fixableRow.SubtreeText, "cleanup Pass 5 removes it") {
		t.Errorf("the remedy is not on the row: %q", fixableRow.SubtreeText)
	}
	skippedRow := results.ListChildren[4]
	if !strings.Contains(skippedRow.SubtreeText, "not checked") {
		t.Errorf("the skipped row carries no not-checked chip: %q", skippedRow.SubtreeText)
	}
	if !strings.Contains(skippedRow.SubtreeText, "no such branch") {
		t.Errorf("the skipped row lost the probe text: %q", skippedRow.SubtreeText)
	}

	// The markup half of the same claim: the card grid and the disclosure are
	// gone from the shipped template, not merely unused by the renderer.
	findingsSection := sliceFindingsStripMarkup(t, indexHtml)
	for _, retiredMarkup := range []string{"<details", "<summary", "board-anomalies-cards", "board-findings-skipped-summary"} {
		if strings.Contains(findingsSection, retiredMarkup) {
			t.Errorf("the findings strip still ships %q:\n%s", retiredMarkup, findingsSection)
		}
	}
	// The two hosts sit inside one list element, and the CSS makes them
	// pass-through, so what ships is one list and not two stacked blocks.
	if !strings.Contains(findingsSection, `id="board-findings-rows"`) {
		t.Errorf("the findings strip has no row list wrapping its hosts:\n%s", findingsSection)
	}
	// The pass-through rule is the whole "one list" claim and no DOM probe can
	// see it: without it the two hosts are ordinary blocks and the rows inside
	// them stop being laid out by the list.
	groupRuleStart := strings.Index(indexHtml, ".board-findings-group {")
	if groupRuleStart < 0 {
		t.Fatal("the shipped page has no .board-findings-group rule")
	}
	groupRule := indexHtml[groupRuleStart:]
	if ruleEnd := strings.Index(groupRule, "}"); ruleEnd >= 0 {
		groupRule = groupRule[:ruleEnd]
	}
	if !strings.Contains(groupRule, "display: contents") {
		t.Errorf("the row hosts are not pass-through, so the rows render as two blocks: %q", groupRule)
	}
}

// hasClassName reports whether a stub node's split className list carries one
// exact class. Substring matching would answer yes for "board-finding" against
// "board-findings-row", which is the very distinction these assertions rest on.
func hasClassName(classNames []string, wantClassName string) bool {
	for _, className := range classNames {
		if className == wantClassName {
			return true
		}
	}
	return false
}

// sliceFindingsStripMarkup returns the `#board-findings` section's own markup, so
// a "no <details> here" assertion cannot be satisfied or broken by a disclosure
// somewhere else on the page.
func sliceFindingsStripMarkup(t *testing.T, indexHtml string) string {
	t.Helper()
	sectionStart := strings.Index(indexHtml, `id="board-findings"`)
	if sectionStart < 0 {
		t.Fatal("the generated page has no #board-findings strip")
	}
	sectionEnd := strings.Index(indexHtml[sectionStart:], "</section>")
	if sectionEnd < 0 {
		t.Fatal("the #board-findings strip is never closed")
	}
	return indexHtml[sectionStart : sectionStart+sectionEnd]
}
