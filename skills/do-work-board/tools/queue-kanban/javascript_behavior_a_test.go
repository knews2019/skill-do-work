package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestJavaScriptBehaviorAssembledClientSyntax(t *testing.T) {
	behaviorProbeCountBefore := javaScriptBehaviorProbeCount.Load()
	defer func() {
		if behaviorProbeCountAfter := javaScriptBehaviorProbeCount.Load(); behaviorProbeCountAfter != behaviorProbeCountBefore {
			t.Errorf("assembled syntax changed behavior-probe count from %d to %d",
				behaviorProbeCountBefore, behaviorProbeCountAfter)
		}
	}()
	assembledClient, assembleError := assembleBoardJavaScript(embeddedWebAssets)
	if assembleError != nil {
		t.Fatalf("assembleBoardJavaScript: %v", assembleError)
	}
	nodePath := lookupNodeForJavaScriptProbe(t)
	syntaxCommand := exec.Command(nodePath, "--check", "-")
	syntaxCommand.Stdin = bytes.NewReader(assembledClient)
	syntaxOutput, syntaxError := syntaxCommand.CombinedOutput()
	if syntaxError != nil {
		t.Fatalf("node --check assembled client: %v\n%s", syntaxError, syntaxOutput)
	}
}

func TestJavaScriptBehaviorClipboardAnnotatesBodiesAndAppendsOneGlossary(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function describeRequestStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildRequestIdByReqSegment("),
		sliceBalancedBlockAfter(t, indexHtml, "function resolveTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function isAmbiguousTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function recordReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateTicketMentions("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildReferencedTicketsGlossary("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateClipboardPayload("),
	}
	declarationBlocks := []string{
		sliceDeclarationAfter(t, indexHtml, "var inlineTicketTitleMaxLength ="),
		sliceDeclarationAfter(t, indexHtml, "var requestIdByReqSegment ="),
		sliceDeclarationAfter(t, indexHtml, "var referencedTicketsGlossaryHeading ="),
	}

	longTitle := "Make every referenced request identifier in a drawer body carry its own title"
	shortenedLongTitle := "Make every referenced request identifier in a drawer body…"
	exactlySixtyTitle := "Keep the timeline forecast honest about ordering and timings"

	// One document carrying every exclusion at once. REQ-1679 sits in the fence
	// AND in the body, `REQ-1108` sits in a code span before its prose mention,
	// REQ-1685 sits in a fenced block before its prose mention, REQ-8888/REQ-8887
	// are dead ids inside fenced blocks, REQ-9999 is a dead id in prose, and
	// REQ-378/UR-075 ride inside a repo-relative path.
	hostDocument := "---\n" +
		"id: REQ-500\n" +
		"depends_on: [REQ-1679]\n" +
		"user_request: UR-074\n" +
		"---\n" +
		"# Host document\n" +
		"\n" +
		"Read REQ-1679 lessons and REQ-1679 again, plus `REQ-1108` and REQ-1108, and UR-074.\n" +
		"\n" +
		"```yaml\n" +
		"depends_on: [REQ-1685, REQ-8888]\n" +
		"```\n" +
		"\n" +
		"~~~text\n" +
		"REQ-8887 illustration\n" +
		"~~~\n" +
		"\n" +
		"Trailing REQ-1685 mention and REQ-9999 too.\n" +
		"See do-work/archive/UR-075/REQ-378-title.md for the path case.\n"

	wantAnnotatedHostDocument := "---\n" +
		"id: REQ-500\n" +
		"depends_on: [REQ-1679]\n" +
		"user_request: UR-074\n" +
		"---\n" +
		"# Host document\n" +
		"\n" +
		"Read REQ-1679 (-> " + exactlySixtyTitle + ") lessons and REQ-1679 again, plus `REQ-1108` and " +
		"REQ-1108 (-> Short one), and UR-074 (-> Ticket ids should carry their titles).\n" +
		"\n" +
		"```yaml\n" +
		"depends_on: [REQ-1685, REQ-8888]\n" +
		"```\n" +
		"\n" +
		"~~~text\n" +
		"REQ-8887 illustration\n" +
		"~~~\n" +
		"\n" +
		"Trailing REQ-1685 (-> " + shortenedLongTitle + ") mention and REQ-9999 too.\n" +
		"See do-work/archive/UR-075/REQ-378-title.md for the path case.\n"

	// The second half of the concatenation trap: its fence must survive the join
	// untouched, and its body's REQ-1679 must expand even though the first
	// document already spent that id — first-mention memory is per document.
	secondDocument := "---\n" +
		"id: REQ-501\n" +
		"depends_on: [REQ-1679, REQ-9998]\n" +
		"---\n" +
		"# Second host document\n" +
		"\n" +
		"Nothing but REQ-1679 here.\n"
	wantAnnotatedSecondDocument := "---\n" +
		"id: REQ-501\n" +
		"depends_on: [REQ-1679, REQ-9998]\n" +
		"---\n" +
		"# Second host document\n" +
		"\n" +
		"Nothing but REQ-1679 (-> " + exactlySixtyTitle + ") here.\n"

	unclosedFenceDocument := "---\nid: REQ-1685\nRead REQ-1108 here\n"
	carriageReturnDocument := "---\r\nid: REQ-500\r\n---\r\nBody REQ-1108 here\r\n"
	fencelessDocument := "# Notes\n\nSee REQ-1108 twice, REQ-1108.\n"
	loneFenceDocument := "---\n"
	ambiguousDocument := "Compare REQ-042 with REQ-042 again.\n"
	noReferenceDocument := "# Plain\n\nNothing here.\n"

	// The outside-text containment contract (actions/clarify.md Step 4) writes
	// every UR's Full Verbatim Input as a BLOCKQUOTED fence, and the contract
	// promises the text stays byte-identical apart from the containment bytes.
	// Two URs in this repo hold ticket ids inside one — UR-075's carries 21 —
	// so annotating inside it rewrites the user's own preserved words.
	blockquotedFenceDocument := "---\nid: REQ-500\n---\n\n" +
		"Prose cites REQ-1679 once.\n\n" +
		"> ````text\n" +
		"> The user pasted REQ-1108 and REQ-1685 here verbatim.\n" +
		"> ````\n\n" +
		"Trailing prose cites REQ-1108.\n"

	// CommonMark forbids a backtick anywhere in a BACKTICK fence's info string,
	// so this line is prose and the ids under it are real references. The Go
	// renderer already agrees (TestRenderMarkdownInvalidBacktickInfoRemainsQuestionProse).
	// Treating it as a fence opened a block that never opens and swallowed every
	// reference until EOF.
	invalidInfoStringDocument := "---\nid: REQ-501\n---\n\n" +
		"```lang`invalid\n" +
		"This line is prose and REQ-1679 in it is a real reference.\n"

	// A fence can open directly as a list item. The prefix stripper's own comment
	// promised list markers before the code handled them, so the promise and the
	// behaviour disagreed — the comment was right and the code was not.
	listItemFenceDocument := "---\nid: REQ-500\n---\n\n" +
		"- ```yaml\n" +
		"  depends_on: [REQ-1679]\n" +
		"  ```\n\n" +
		"Prose after the list cites REQ-1108.\n"

	// A code span may cross a line break — CommonMark closes it on the matching
	// run anywhere in the paragraph. Line-by-line scanning read the opener as a
	// stray backtick and expanded the id inside. Live in REQ-380's body.
	// The id sits on the CONTINUATION line on purpose. With it only on the
	// opening line, dropping the cross-line carry changes nothing observable and
	// the mutation passes — a vacuous assertion, which is the failure this suite
	// has now shipped twice.
	multiLineCodeSpanDocument := "---\nid: REQ-501\n---\n\n" +
		"the example reads\n" +
		"`- REQ-1679 a quoted worked example — the second\n" +
		"finding for REQ-1108` matters.\n\n" +
		"Trailing prose cites REQ-1108 again.\n"

	// The stub board and the Go resolver are built from ONE list of ids, so the
	// positions spliced in below were computed against exactly the records the
	// client looks titles up in. Two lists would let the halves disagree
	// silently, which is the failure this whole probe exists to catch.
	clipboardProbeResolver := newTicketMentionResolver(
		[]string{"REQ-1679", "REQ-1108", "REQ-1685", "REQ-500", "REQ-501", "UR-001-REQ-042", "UR-002-REQ-042"},
		[]string{"UR-074"},
	)
	probeDocument := func(documentText string) string {
		return clipboardProbeDocument(t, documentText, clipboardProbeResolver)
	}

	javascriptProbe := `
var requestsById = {
  "REQ-1679": { title: ` + mustMarshalJSONString(t, exactlySixtyTitle) + `, status: "completed" },
  "REQ-1108": { title: "Short one", status: "pending" },
  "REQ-1685": { title: ` + mustMarshalJSONString(t, longTitle) + `, status: "claimed" },
  "REQ-500": { title: "Host document", status: "claimed" },
  "REQ-501": { title: "Second host document", status: "pending" },
  "UR-001-REQ-042": { title: "First half of an ambiguous pair", status: "pending" },
  "UR-002-REQ-042": { title: "Second half of an ambiguous pair", status: "pending" }
};
var userRequestsById = {
  "UR-074": { title: "Ticket ids should carry their titles" }
};
` + strings.Join(functionBlocks, "\n") + "\n" + strings.Join(declarationBlocks, "\n") + `

var hostDocument = ` + probeDocument(hostDocument) + `;
var secondDocument = ` + probeDocument(secondDocument) + `;
var annotatedHost = annotateTicketMentions(hostDocument.text, hostDocument.ticketMentions);
var joinedPayload = annotateClipboardPayload([hostDocument, secondDocument], ["REQ-500", "REQ-501"]);
var glossaryHeadingCount = joinedPayload.split(referencedTicketsGlossaryHeading).length - 1;

process.stdout.write(JSON.stringify({
  annotatedHostDocument: annotatedHost.text,
  hostReferencedIds: annotatedHost.referencedTickets.map(function (entry) { return entry.id; }),
  joinedPayload: joinedPayload,
  glossaryHeadingCount: glossaryHeadingCount,
  excludedPayload: annotateClipboardPayload(
    [hostDocument, secondDocument], ["REQ-500", "REQ-501", "REQ-1679", "REQ-1108"]
  ),
  unclosedFencePayload: annotateClipboardPayload([` + probeDocument(unclosedFenceDocument) + `], []),
  carriageReturnPayload: annotateClipboardPayload([` + probeDocument(carriageReturnDocument) + `], ["REQ-500"]),
  fencelessPayload: annotateClipboardPayload([` + probeDocument(fencelessDocument) + `], ["REQ-1108"]),
  loneFencePayload: annotateClipboardPayload([` + probeDocument(loneFenceDocument) + `], []),
  ambiguousPayload: annotateClipboardPayload([` + probeDocument(ambiguousDocument) + `], []),
  noReferencePayload: annotateClipboardPayload([` + probeDocument(noReferenceDocument) + `], []),
  blockquotedFencePayload: annotateClipboardPayload([` + probeDocument(blockquotedFenceDocument) + `], ["REQ-500"]),
  invalidInfoStringPayload: annotateClipboardPayload([` + probeDocument(invalidInfoStringDocument) + `], ["REQ-501"]),
  listItemFencePayload: annotateClipboardPayload([` + probeDocument(listItemFenceDocument) + `], ["REQ-500"]),
  multiLineCodeSpanPayload: annotateClipboardPayload([` + probeDocument(multiLineCodeSpanDocument) + `], ["REQ-501"])
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "clipboard ticket annotation", javascriptProbe)
	var probeResult clipboardAnnotationProbeResult
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode clipboard annotation behavior: %v (output %q)", decodeError, probeOutput)
	}

	if probeResult.AnnotatedHostDocument != wantAnnotatedHostDocument {
		t.Errorf("annotated host document:\n got %q\nwant %q", probeResult.AnnotatedHostDocument, wantAnnotatedHostDocument)
	}
	wantHostReferencedIds := []string{"REQ-1679", "REQ-1108", "UR-074", "REQ-1685", "REQ-9999"}
	if !reflect.DeepEqual(probeResult.HostReferencedIds, wantHostReferencedIds) {
		t.Errorf("host references = %v, want first-mention order %v (a fenced or path-borne id must not appear)",
			probeResult.HostReferencedIds, wantHostReferencedIds)
	}

	wantGlossary := "\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- REQ-1679 — " + exactlySixtyTitle + " (completed)\n" +
		"- REQ-1108 — Short one (pending)\n" +
		"- UR-074 — Ticket ids should carry their titles (user request)\n" +
		"- REQ-1685 — " + longTitle + " (claimed)\n" +
		"- REQ-9999 — not found in this queue\n"
	wantJoinedPayload := wantAnnotatedHostDocument + wantAnnotatedSecondDocument + wantGlossary
	if probeResult.JoinedPayload != wantJoinedPayload {
		t.Errorf("joined clipboard payload:\n got %q\nwant %q", probeResult.JoinedPayload, wantJoinedPayload)
	}
	if probeResult.GlossaryHeadingCount != 1 {
		t.Errorf("glossary heading appeared %d times, want exactly one appendix at the end", probeResult.GlossaryHeadingCount)
	}

	wantExcludedGlossary := "\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- UR-074 — Ticket ids should carry their titles (user request)\n" +
		"- REQ-1685 — " + longTitle + " (claimed)\n" +
		"- REQ-9999 — not found in this queue\n"
	wantExcludedPayload := wantAnnotatedHostDocument + wantAnnotatedSecondDocument + wantExcludedGlossary
	if probeResult.ExcludedPayload != wantExcludedPayload {
		t.Errorf("payload with excluded ids:\n got %q\nwant %q", probeResult.ExcludedPayload, wantExcludedPayload)
	}

	// No closing fence means everything is body, exactly as splitFrontmatter
	// decides on the Go side. Reading it as an unterminated fence would skip the
	// whole document and annotate nothing.
	wantUnclosedFencePayload := "---\nid: REQ-1685 (-> " + shortenedLongTitle + ")\nRead REQ-1108 (-> Short one) here\n" +
		"\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- REQ-1685 — " + longTitle + " (claimed)\n" +
		"- REQ-1108 — Short one (pending)\n"
	if probeResult.UnclosedFencePayload != wantUnclosedFencePayload {
		t.Errorf("unclosed-fence payload:\n got %q\nwant %q", probeResult.UnclosedFencePayload, wantUnclosedFencePayload)
	}

	// CRLF endings survive byte-for-byte: the body is never normalized, only
	// extended.
	wantCarriageReturnPayload := "---\r\nid: REQ-500\r\n---\r\nBody REQ-1108 (-> Short one) here\r\n" +
		"\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- REQ-1108 — Short one (pending)\n"
	if probeResult.CarriageReturnPayload != wantCarriageReturnPayload {
		t.Errorf("CRLF payload:\n got %q\nwant %q", probeResult.CarriageReturnPayload, wantCarriageReturnPayload)
	}

	// The drawer's rendered-text fallback has no fence at all, and a repeat
	// mention there stays bare.
	wantFencelessPayload := "# Notes\n\nSee REQ-1108 (-> Short one) twice, REQ-1108.\n"
	if probeResult.FencelessPayload != wantFencelessPayload {
		t.Errorf("fence-less payload:\n got %q\nwant %q", probeResult.FencelessPayload, wantFencelessPayload)
	}
	if probeResult.LoneFencePayload != loneFenceDocument {
		t.Errorf("lone-fence payload = %q, want the document unchanged", probeResult.LoneFencePayload)
	}
	// Ambiguous is not missing: the board holds records that match and refuses to
	// pick one, so it earns neither an expansion nor a not-found line.
	if probeResult.AmbiguousPayload != ambiguousDocument {
		t.Errorf("ambiguous payload = %q, want the document unchanged", probeResult.AmbiguousPayload)
	}
	if probeResult.NoReferencePayload != noReferenceDocument {
		t.Errorf("payload citing nothing = %q, want no appendix at all", probeResult.NoReferencePayload)
	}

	// A fence inside a blockquote is a fence. The outside-text containment
	// contract writes every UR's Full Verbatim Input this way and promises the
	// text stays byte-identical apart from the containment bytes, so an id
	// inside one is preserved words, not a reference. UR-075 holds 21 of them.
	if strings.Contains(probeResult.BlockquotedFence, "REQ-1108 (-> Short one)\n> ") ||
		strings.Contains(probeResult.BlockquotedFence, "> The user pasted REQ-1108 (") {
		t.Errorf("a blockquoted fence was annotated — the containment contract's preserved text was rewritten:\n%s",
			probeResult.BlockquotedFence)
	}
	if !strings.Contains(probeResult.BlockquotedFence, "> The user pasted REQ-1108 and REQ-1685 here verbatim.\n") {
		t.Errorf("the blockquoted verbatim line is not byte-identical:\n%s", probeResult.BlockquotedFence)
	}
	// Prose on either side of that block still expands, or the fix would have
	// been to stop annotating rather than to recognise the container.
	if !strings.Contains(probeResult.BlockquotedFence, "Prose cites REQ-1679 (-> ") {
		t.Errorf("prose before a blockquoted fence lost its expansion:\n%s", probeResult.BlockquotedFence)
	}

	// A fence opened as a list item is a fence. The prefix stripper's comment
	// promised list markers before the code stripped them, so ids inside such a
	// block were expanded as prose and its closer could be misread as a new
	// opener, suppressing annotation of everything after it.
	if strings.Contains(probeResult.ListItemFence, "depends_on: [REQ-1679 (") {
		t.Errorf("an id inside a list-item fence was expanded:\n%s", probeResult.ListItemFence)
	}
	if !strings.Contains(probeResult.ListItemFence, "Prose after the list cites REQ-1108 (-> ") {
		t.Errorf("prose after a list-item fence lost its expansion — the closer was misread as an opener:\n%s",
			probeResult.ListItemFence)
	}

	// A code span may cross a line break. Reading the opener as a stray backtick
	// expanded the id inside a quoted worked example — live in REQ-380's body.
	// Neither id inside the span may expand — the one on the opening line or the
	// one on the continuation line. The continuation id is the discriminator:
	// without the cross-line carry it is treated as prose and expands, and then
	// the trailing prose mention becomes a repeat and stays bare, so both halves
	// of this pair flip together.
	if strings.Contains(probeResult.MultiLineCodeSpan, "REQ-1679 (") ||
		strings.Contains(probeResult.MultiLineCodeSpan, "finding for REQ-1108 (") {
		t.Errorf("an id inside a code span crossing a newline was expanded:\n%s", probeResult.MultiLineCodeSpan)
	}
	if !strings.Contains(probeResult.MultiLineCodeSpan, "finding for REQ-1108` matters.") {
		t.Errorf("the code span's continuation line is not byte-identical:\n%s", probeResult.MultiLineCodeSpan)
	}
	if !strings.Contains(probeResult.MultiLineCodeSpan, "Trailing prose cites REQ-1108 (-> ") {
		t.Errorf("prose after a multi-line code span lost its expansion:\n%s", probeResult.MultiLineCodeSpan)
	}

	// CommonMark forbids a backtick in a backtick fence's info string, so the
	// line is prose and what follows it is not fenced. Accepting it opened a
	// block that never opens and left every later reference bare.
	if !strings.Contains(probeResult.InvalidInfoString, "REQ-1679 (-> ") {
		t.Errorf("an invalid backtick info string opened a fence that CommonMark calls prose, "+
			"so the reference under it was left bare:\n%s", probeResult.InvalidInfoString)
	}
}

func TestJavaScriptBehaviorDurationsWindowSelectionRefreshesOnlyDurations(t *testing.T) {
	indexHTML := generateLiveSite(t)
	transitionFunction := sliceBalancedBlockAfter(t, indexHTML, "function applyDurationsWindowSelection(")
	javascriptProbe := `
var viewState = { windowHours: 24, view: "durations" };
var renderedOnce = { durations: true };
var chosenWindows = [];
var activeWindows = [];
var renderCount = 0;
function setDurationsWindow(windowName) { chosenWindows.push(windowName); }
function setActiveButton(selector, attributeName, attributeValue) { activeWindows.push(attributeValue); }
function renderDurationsView() { renderCount += 1; }
` + transitionFunction + `
["30", "90", "all"].forEach(function (windowName) { applyDurationsWindowSelection(windowName); });
var visibleState = { renderCount: renderCount, rendered: renderedOnce.durations };
viewState.view = "board";
applyDurationsWindowSelection("30");
process.stdout.write(JSON.stringify({
  chosenWindows: chosenWindows,
  activeWindows: activeWindows,
  windowHours: viewState.windowHours,
  visibleState: visibleState,
  hiddenRenderCount: renderCount,
  hiddenRendered: renderedOnce.durations
}));`
	probeOutput := runJavaScriptBehaviorProbe(t, "Durations window transitions", javascriptProbe)
	var result struct {
		ChosenWindows []string `json:"chosenWindows"`
		ActiveWindows []string `json:"activeWindows"`
		WindowHours   int      `json:"windowHours"`
		VisibleState  struct {
			RenderCount int  `json:"renderCount"`
			Rendered    bool `json:"rendered"`
		} `json:"visibleState"`
		HiddenRenderCount int  `json:"hiddenRenderCount"`
		HiddenRendered    bool `json:"hiddenRendered"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode Durations-window transition: %v (output %q)", decodeError, probeOutput)
	}
	wantWindows := []string{"30", "90", "all", "30"}
	if strings.Join(result.ChosenWindows, ",") != strings.Join(wantWindows, ",") ||
		strings.Join(result.ActiveWindows, ",") != strings.Join(wantWindows, ",") {
		t.Fatalf("Durations selections = chosen %#v active %#v, want %#v", result.ChosenWindows, result.ActiveWindows, wantWindows)
	}
	if result.WindowHours != 24 {
		t.Fatalf("Durations selection changed viewState.windowHours to %d, want unchanged 24", result.WindowHours)
	}
	if result.VisibleState.RenderCount != 3 || !result.VisibleState.Rendered {
		t.Fatalf("three visible selections produced %#v, want one render per selection and fresh state", result.VisibleState)
	}
	if result.HiddenRenderCount != 3 || result.HiddenRendered {
		t.Fatalf("hidden selection produced renderCount=%d rendered=%v, want no render and stale cache", result.HiddenRenderCount, result.HiddenRendered)
	}
}

func TestJavaScriptBehaviorTestingDoneWindowIsViewSpecific(t *testing.T) {
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function hasActiveVisibleFilters("),
		sliceBalancedBlockAfter(t, indexHtml, "function formatFilteredCount("),
		sliceBalancedBlockAfter(t, indexHtml, "function columnEmptyText("),
		sliceBalancedBlockAfter(t, indexHtml, "function fillColumn("),
		sliceBalancedBlockAfter(t, indexHtml, "function fillTestingColumn("),
	}
	javascriptProbe := `
var filterState = { searchText: "", domain: "", status: "", doneWindow: "168" };
var viewState = { view: "board" };
var nodesBySelector = {};
function makeNode() {
  return {
    childNodes: [],
    textContent: "",
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
}
var document = {
  createElement: function () { return makeNode(); },
  querySelector: function (selector) {
    if (!nodesBySelector[selector]) {
      nodesBySelector[selector] = makeNode();
    }
    return nodesBySelector[selector];
  }
};
` + strings.Join(functionBlocks, "\n") + `
fillColumn("board", [], null, 1);
var boardCopy = nodesBySelector['[data-cards="board"]'].childNodes[0].textContent;
viewState.view = "testing";
fillColumn("hidden-board", [], null, 1);
var hiddenBoardCopy = nodesBySelector['[data-cards="hidden-board"]'].childNodes[0].textContent;
fillTestingColumn("testing-ready", [], 1);
var testingCopy = nodesBySelector['[data-cards="testing-ready"]'].childNodes[0].textContent;
process.stdout.write(JSON.stringify([boardCopy, hiddenBoardCopy, testingCopy]));`
	probeOutput := runJavaScriptBehaviorProbe(t, "testing empty-copy decision", javascriptProbe)
	var results []string
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode assembled-client empty-copy results: %v (output %q)", decodeError, probeOutput)
	}
	if len(results) != 3 {
		t.Fatalf("empty-copy result count = %d, want 3: %#v", len(results), results)
	}
	if results[0] != "Nothing here" {
		t.Fatalf("Board empty copy with only doneWindow = %q, want Nothing here", results[0])
	}
	if results[1] != "Nothing here" {
		t.Fatalf("hidden Board empty copy during Testing view = %q, want Nothing here", results[1])
	}
	if results[2] != "No matches" {
		t.Fatalf("Testing empty copy with doneWindow = %q, want No matches", results[2])
	}
}

func TestJavaScriptBehaviorDurationsHeadlineRollingMedianAndCadenceTicks(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}
	headlineJSON, encodeError := json.Marshal(durationHeadlineFixtureData(t))
	if encodeError != nil {
		t.Fatalf("encode headline fixture: %v", encodeError)
	}
	rollingPayloads := map[string]generatedDurations{
		"six":   durationRollingFixtureData(t, 6),
		"seven": durationRollingFixtureData(t, 7),
		"eight": durationRollingFixtureData(t, 8),
	}
	rollingJSON, encodeError := json.Marshal(rollingPayloads)
	if encodeError != nil {
		t.Fatalf("encode rolling fixtures: %v", encodeError)
	}

	probeDriver := `
function resetDurationsHosts() {
  ["durations-chart", "durations-summary", "durations-stat-median", "durations-stat-p90",
   "durations-stat-active-days", "durations-stat-reqs-per-day", "durations-readout",
   "durations-table-body"].forEach(function (nodeId) {
    var nodeName = nodeId.indexOf("stat-") >= 0 ? "dd" : "div";
    durationsStubHosts[nodeId] = makeStubNode(nodeName);
  });
}
function nodeText(node) {
  return (node.children || []).map(function (child) {
    return child.textContent !== undefined ? child.textContent : nodeText(child);
  }).join("");
}
function captureRender(payload, windowName) {
  resetDurationsHosts();
  boardData = { durations: payload };
  setDurationsWindow(windowName);
  renderDurationsView();
  var svg = durationsStubHosts["durations-chart"].children[0];
  var rollingPaths = [], rollingMarkers = [], panelBBars = [], countTicks = [], countGridlines = [];
  (svg.children || []).forEach(function (child, childIndex) {
    var attributes = child.attributes || {};
    var className = String(attributes["class"] || "");
    if (className === "durations-bar") { panelBBars.push({ childIndex: childIndex }); }
    if (className === "durations-rolling-line") {
      rollingPaths.push({ d: attributes.d || "", childIndex: childIndex });
    }
    if (className === "durations-rolling-marker") {
      rollingMarkers.push({ cx: Number(attributes.cx), cy: Number(attributes.cy), childIndex: childIndex });
    }
    if (attributes["data-durations-count-tick"] === "true") {
      countTicks.push({ text: nodeText(child), y: Number(attributes.y) });
    }
    if (attributes["data-durations-count-grid"] === "midpoint") {
      countGridlines.push({ y: Number(attributes.y1), childIndex: childIndex });
    }
  });
  return {
    stats: [
      durationsStubHosts["durations-stat-median"].textContent,
      durationsStubHosts["durations-stat-p90"].textContent,
      durationsStubHosts["durations-stat-active-days"].textContent,
      durationsStubHosts["durations-stat-reqs-per-day"].textContent
    ],
    summary: durationsStubHosts["durations-summary"].textContent,
    ariaLabel: svg.attributes["aria-label"],
    visibleTitles: (svg.children || []).filter(function (child) {
      return child.stubName === "text" && String(child.attributes["class"] || "").indexOf("durations-axis-title") >= 0;
    }).map(nodeText),
    rollingPaths: rollingPaths,
    rollingMarkers: rollingMarkers,
    panelBBars: panelBBars,
    countTicks: countTicks,
    countGridlines: countGridlines
  };
}
process.stdout.write(JSON.stringify({
  headline30: captureRender(` + string(headlineJSON) + `, "30"),
  headline90: captureRender(` + string(headlineJSON) + `, "90"),
  headlineAll: captureRender(` + string(headlineJSON) + `, "all"),
  rollingSix: captureRender(` + string(rollingJSON) + `.six, "all"),
  rollingSeven: captureRender(` + string(rollingJSON) + `.seven, "all"),
  rollingEight: captureRender(` + string(rollingJSON) + `.eight, "all")
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "Durations headline, rolling median, and cadence ticks",
		durationsRenderDomStubPreamble+string(rendererFragment)+probeDriver)

	type capturedRender struct {
		Stats         []string `json:"stats"`
		Summary       string   `json:"summary"`
		AriaLabel     string   `json:"ariaLabel"`
		VisibleTitles []string `json:"visibleTitles"`
		RollingPaths  []struct {
			D          string `json:"d"`
			ChildIndex int    `json:"childIndex"`
		} `json:"rollingPaths"`
		RollingMarkers []struct {
			CX         float64 `json:"cx"`
			CY         float64 `json:"cy"`
			ChildIndex int     `json:"childIndex"`
		} `json:"rollingMarkers"`
		PanelBBars []struct {
			ChildIndex int `json:"childIndex"`
		} `json:"panelBBars"`
		CountTicks []struct {
			Text string  `json:"text"`
			Y    float64 `json:"y"`
		} `json:"countTicks"`
		CountGridlines []struct {
			Y          float64 `json:"y"`
			ChildIndex int     `json:"childIndex"`
		} `json:"countGridlines"`
	}
	var result struct {
		Headline30   capturedRender `json:"headline30"`
		Headline90   capturedRender `json:"headline90"`
		HeadlineAll  capturedRender `json:"headlineAll"`
		RollingSix   capturedRender `json:"rollingSix"`
		RollingSeven capturedRender `json:"rollingSeven"`
		RollingEight capturedRender `json:"rollingEight"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode Durations headline/rolling result: %v (output starts %q)",
			decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}

	for _, windowCase := range []struct {
		name      string
		got       capturedRender
		wantStats []string
	}{
		{name: "30", got: result.Headline30, wantStats: []string{"25.0 min", "5h 50m", "3 / 30", "2.0"}},
		{name: "90", got: result.Headline90, wantStats: []string{"35.0 min", "5h 30m", "5 / 90", "1.6"}},
		{name: "all", got: result.HeadlineAll, wantStats: []string{"45.0 min", "5h 10m", "7 / 121", "1.4"}},
	} {
		if !reflect.DeepEqual(windowCase.got.Stats, windowCase.wantStats) {
			t.Errorf("%s-day headline stats = %#v, want %#v", windowCase.name, windowCase.got.Stats, windowCase.wantStats)
		}
	}
	wantExclusionSentence := "Panel B excludes 3 spans from its medians (over four hours is an assumed pause, negative is a broken stamp); panel A still plots them."
	if !strings.Contains(result.Headline30.Summary, wantExclusionSentence) {
		t.Errorf("summary exclusion rule changed: %q", result.Headline30.Summary)
	}

	if len(result.RollingSix.RollingMarkers) != 0 || len(result.RollingSix.RollingPaths) != 0 {
		t.Errorf("six eligible days drew %d markers and %d paths, want neither",
			len(result.RollingSix.RollingMarkers), len(result.RollingSix.RollingPaths))
	}
	if len(result.RollingSeven.RollingMarkers) != 1 || len(result.RollingSeven.RollingPaths) != 0 {
		t.Errorf("seven eligible days drew %d markers and %d paths, want one marker and no path",
			len(result.RollingSeven.RollingMarkers), len(result.RollingSeven.RollingPaths))
	}
	if len(result.RollingEight.RollingMarkers) != 2 || len(result.RollingEight.RollingPaths) != 1 {
		t.Fatalf("eight eligible days drew %d markers and %d paths, want two markers and one path",
			len(result.RollingEight.RollingMarkers), len(result.RollingEight.RollingPaths))
	}
	if !strings.Contains(result.RollingEight.VisibleTitles[2], "trailing 7-active-day median") ||
		!strings.Contains(result.RollingEight.AriaLabel, "trailing 7-active-day median") {
		t.Errorf("Panel B title/accessibility copy does not name trailing 7-active-day median: titles=%q aria=%q",
			result.RollingEight.VisibleTitles, result.RollingEight.AriaLabel)
	}
	medianTop := durationsRendererConstant(t, "DURATIONS_MEDIAN_TOP")
	medianBottom := durationsRendererConstant(t, "DURATIONS_MEDIAN_BOTTOM")
	wantRolling := []struct {
		day    time.Time
		median float64
	}{
		{day: time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC), median: 40},
		{day: time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC), median: 50},
	}
	timeStart := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	timeEnd := time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)
	marginLeft := durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")
	plotWidth := durationsRendererConstant(t, "DURATIONS_VIEW_WIDTH") - marginLeft - durationsRendererConstant(t, "DURATIONS_MARGIN_RIGHT")
	for markerIndex, marker := range result.RollingEight.RollingMarkers {
		wantX := marginLeft + (wantRolling[markerIndex].day.Add(12*time.Hour).Sub(timeStart).Seconds()/timeEnd.Sub(timeStart).Seconds())*plotWidth
		wantY := medianBottom - (math.Min(wantRolling[markerIndex].median, 45)/45)*(medianBottom-medianTop)
		if math.Abs(marker.CX-wantX) > 0.11 || math.Abs(marker.CY-wantY) > 0.11 {
			t.Errorf("rolling marker %d = (%.2f, %.2f), want active-day trailing point (%.2f, %.2f)",
				markerIndex, marker.CX, marker.CY, wantX, wantY)
		}
	}
	lastBarIndex := result.RollingEight.PanelBBars[len(result.RollingEight.PanelBBars)-1].ChildIndex
	if result.RollingEight.RollingPaths[0].ChildIndex <= lastBarIndex ||
		result.RollingEight.RollingMarkers[0].ChildIndex <= result.RollingEight.RollingPaths[0].ChildIndex {
		t.Errorf("rolling draw order paths=%+v markers=%+v last bar=%d; want bars, path, markers",
			result.RollingEight.RollingPaths, result.RollingEight.RollingMarkers, lastBarIndex)
	}

	wantCountTicks := map[string]float64{
		"0":   durationsRendererConstant(t, "DURATIONS_COUNT_BOTTOM") + durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP"),
		"2.5": (durationsRendererConstant(t, "DURATIONS_COUNT_TOP")+durationsRendererConstant(t, "DURATIONS_COUNT_BOTTOM"))/2 + durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP"),
		"5":   durationsRendererConstant(t, "DURATIONS_COUNT_TOP") + durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP"),
	}
	if len(result.RollingEight.CountTicks) != len(wantCountTicks) {
		t.Fatalf("Panel C ticks = %+v, want zero, exact midpoint, and peak", result.RollingEight.CountTicks)
	}
	for _, tick := range result.RollingEight.CountTicks {
		wantY, exists := wantCountTicks[tick.Text]
		if !exists || math.Abs(tick.Y-wantY) > 0.01 {
			t.Errorf("Panel C tick %q at %.2f, want exact tick map %v", tick.Text, tick.Y, wantCountTicks)
		}
	}
	if len(result.RollingEight.CountGridlines) != 1 ||
		math.Abs(result.RollingEight.CountGridlines[0].Y-(durationsRendererConstant(t, "DURATIONS_COUNT_TOP")+durationsRendererConstant(t, "DURATIONS_COUNT_BOTTOM"))/2) > 0.01 {
		t.Errorf("Panel C midpoint gridlines = %+v, want one at exact half height", result.RollingEight.CountGridlines)
	}
}

func TestJavaScriptBehaviorDurationsPanelASpreadAndDailyDistribution(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	dayStart := time.Date(2026, 8, 10, 9, 0, 0, 0, time.UTC)
	fixtureMinutes := []float64{0, 5, 15, 30, 45, 60, 300, -20}
	fixtureTickets := make([]*RequestTicket, 0, len(fixtureMinutes)+2)
	for sampleIndex, minutes := range fixtureMinutes {
		completedAt := dayStart.Add(time.Duration(sampleIndex) * time.Hour)
		claimedAt := completedAt.Add(-time.Duration(minutes * float64(time.Minute)))
		fixtureTickets = append(fixtureTickets, durationTicket(
			fmt.Sprintf("REQ-%03d", sampleIndex+1), "B",
			claimedAt.Format(time.RFC3339), completedAt.Format(time.RFC3339),
		))
	}
	// A missing-route sample proves that the lower ordinary opacity does not
	// erase the outlined unknown category.
	unknownCompletedAt := dayStart.Add(24*time.Hour + 2*time.Hour)
	fixtureTickets = append(fixtureTickets, durationTicket(
		"REQ-901", "", unknownCompletedAt.Add(-10*time.Minute).Format(time.RFC3339), unknownCompletedAt.Format(time.RFC3339),
	))
	thirdDayCompletedAt := dayStart.Add(48*time.Hour + 3*time.Hour)
	fixtureTickets = append(fixtureTickets, durationTicket(
		"REQ-902", "A", thirdDayCompletedAt.Add(-20*time.Minute).Format(time.RFC3339), thirdDayCompletedAt.Format(time.RFC3339),
	))

	aggregate := buildDurationAggregate(fixtureTickets)
	generatedData, buildError := buildGeneratedBoardData(&Board{AllRequests: fixtureTickets})
	if buildError != nil {
		t.Fatalf("build generated board data: %v", buildError)
	}
	durationsJSON, encodeError := json.Marshal(generatedData.Durations)
	if encodeError != nil {
		t.Fatalf("encode durations payload: %v", encodeError)
	}

	probeDriver := `
renderDurationsView();
var svg = durationsStubHosts["durations-chart"].children[0];
var ticks = [], circles = [], paths = [];
(svg.children || []).forEach(function (childNode, childIndex) {
  var attributes = childNode.attributes || {};
  if (childNode.stubName === "text" && String(attributes["class"] || "").indexOf("durations-tick") !== -1) {
    ticks.push({ text: ((childNode.children || [])[0] || {}).textContent || "", x: Number(attributes.x), y: Number(attributes.y) });
  }
  if (childNode.stubName === "circle" && String(attributes["class"] || "").indexOf("durations-mark") !== -1) {
    circles.push({
      cx: Number(attributes.cx), cy: Number(attributes.cy), opacity: attributes.opacity || "",
      fill: attributes.fill || "", class: attributes["class"] || "", childIndex: childIndex
    });
  }
  if (childNode.stubName === "path") {
    paths.push({ class: attributes["class"] || "", d: attributes.d || "", childIndex: childIndex });
  }
});
process.stdout.write(JSON.stringify({ ticks: ticks, circles: circles, paths: paths }));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "durations panel A spread and distribution",
		durationsRenderDomStubPreamble+
			"var boardData = { durations: "+string(durationsJSON)+" };\n"+
			string(rendererFragment)+probeDriver)

	var drawn struct {
		Ticks []struct {
			Text string  `json:"text"`
			X    float64 `json:"x"`
			Y    float64 `json:"y"`
		} `json:"ticks"`
		Circles []struct {
			CX         float64 `json:"cx"`
			CY         float64 `json:"cy"`
			Opacity    string  `json:"opacity"`
			Fill       string  `json:"fill"`
			Class      string  `json:"class"`
			ChildIndex int     `json:"childIndex"`
		} `json:"circles"`
		Paths []struct {
			Class      string `json:"class"`
			D          string `json:"d"`
			ChildIndex int    `json:"childIndex"`
		} `json:"paths"`
	}
	if decodeError := json.Unmarshal(probeOutput, &drawn); decodeError != nil {
		t.Fatalf("decode panel A geometry: %v (output starts %q)", decodeError, string(probeOutput[:min(len(probeOutput), 400)]))
	}
	if len(drawn.Circles) != len(aggregate.Samples) {
		t.Fatalf("drew %d marks for %d samples", len(drawn.Circles), len(aggregate.Samples))
	}

	// The tick set is exact, including the new 5-minute foothold. Its y positions
	// follow sqrt(minutes / 60), not the old linear scale.
	panelATicks := map[string]float64{}
	for _, tick := range drawn.Ticks {
		if math.Abs(tick.X-(durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT")-8)) < 0.01 &&
			tick.Y <= durationsRendererConstant(t, "DURATIONS_MAIN_BOTTOM")+durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP")+0.01 {
			panelATicks[tick.Text] = tick.Y - durationsRendererConstant(t, "DURATIONS_TICK_BASELINE_DROP")
		}
	}
	wantTickMinutes := []float64{0, 5, 15, 30, 45, 60}
	if len(panelATicks) != len(wantTickMinutes)+1 { // plus the separate 60+ overflow-lane tick
		t.Fatalf("Panel A ticks = %v, want 60+ and exactly 0, 5, 15, 30, 45, 60", panelATicks)
	}
	mainTop := durationsRendererConstant(t, "DURATIONS_MAIN_TOP")
	mainBottom := durationsRendererConstant(t, "DURATIONS_MAIN_BOTTOM")
	for _, minutes := range wantTickMinutes {
		label := strconv.FormatFloat(minutes, 'f', -1, 64)
		gotY, exists := panelATicks[label]
		if !exists {
			t.Errorf("Panel A has no %s-minute tick", label)
			continue
		}
		wantY := mainBottom - math.Sqrt(minutes/60)*(mainBottom-mainTop)
		if math.Abs(gotY-wantY) > 0.01 {
			t.Errorf("%s-minute tick y=%.3f, want sqrt-scale y=%.3f", label, gotY, wantY)
		}
	}

	ordinaryMark := drawn.Circles[1]
	if ordinaryMark.Opacity == "" || ordinaryMark.Opacity == "1" {
		t.Errorf("ordinary mark opacity = %q, want a stated lower opacity", ordinaryMark.Opacity)
	}
	reversedMark := drawn.Circles[7]
	if reversedMark.Fill != "var(--durations-critical)" || reversedMark.Opacity != "" {
		t.Errorf("reversed mark fill/opacity = %q / %q, want undimmed critical red", reversedMark.Fill, reversedMark.Opacity)
	}
	unknownMark := drawn.Circles[8]
	if !strings.Contains(unknownMark.Class, "durations-mark-unknown") || unknownMark.Opacity != "" {
		t.Errorf("unknown mark class/opacity = %q / %q, want an undimmed outlined unknown", unknownMark.Class, unknownMark.Opacity)
	}

	if len(drawn.Paths) != 2 {
		t.Fatalf("drew %d Panel A distribution paths, want one p25-p75 ribbon and one median line: %+v", len(drawn.Paths), drawn.Paths)
	}
	pathByClass := map[string]struct {
		D          string
		ChildIndex int
	}{}
	for _, path := range drawn.Paths {
		pathByClass[path.Class] = struct {
			D          string
			ChildIndex int
		}{path.D, path.ChildIndex}
		if path.D == "" || strings.Contains(path.D, "NaN") || strings.Contains(path.D, "Infinity") {
			t.Errorf("%q path has invalid geometry %q", path.Class, path.D)
		}
	}
	ribbon, ribbonExists := pathByClass["durations-quantile-ribbon"]
	median, medianExists := pathByClass["durations-quantile-median"]
	if !ribbonExists || !medianExists {
		t.Fatalf("distribution path classes = %v, want ribbon and median", pathByClass)
	}
	if ribbon.ChildIndex >= drawn.Circles[0].ChildIndex || median.ChildIndex >= drawn.Circles[0].ChildIndex {
		t.Errorf("distribution paths at child indexes %d/%d are not behind first circle at %d",
			ribbon.ChildIndex, median.ChildIndex, drawn.Circles[0].ChildIndex)
	}
}

func TestJavaScriptBehaviorTimelineZoomHoldsTheAnchorInstant(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineZoomedWindow(") + `
var boundStart = 0;
var boundEnd = 30 * 24 * 3600 * 1000;   // a 30-day board
var startWindow = { windowStartMs: boundStart, windowEndMs: boundEnd };

function anchorInstant(window, fraction) {
  return window.windowStartMs + (window.windowEndMs - window.windowStartMs) * fraction;
}

// Zoom in three times at the same off-centre anchor; the instant under it must
// not move.
var anchorFraction = 0.25;
var wantAnchor = anchorInstant(startWindow, anchorFraction);
var zoomed = startWindow;
for (var step = 0; step < 3; step++) {
  zoomed = timelineZoomedWindow(zoomed.windowStartMs, zoomed.windowEndMs, 1.6, anchorFraction, boundStart, boundEnd);
}
var anchorDriftMs = Math.abs(anchorInstant(zoomed, anchorFraction) - wantAnchor);

// Zooming all the way back out clamps to the bounds rather than overshooting.
var wideOpen = zoomed;
for (var back = 0; back < 12; back++) {
  wideOpen = timelineZoomedWindow(wideOpen.windowStartMs, wideOpen.windowEndMs, 1 / 1.6, 0.5, boundStart, boundEnd);
}

// Zooming all the way in stops at the floor rather than collapsing to zero.
var deep = startWindow;
for (var deeper = 0; deeper < 40; deeper++) {
  deep = timelineZoomedWindow(deep.windowStartMs, deep.windowEndMs, 1.6, 0.5, boundStart, boundEnd);
}

process.stdout.write(JSON.stringify({
  anchorDriftMs: anchorDriftMs,
  widestSpanMs: wideOpen.windowEndMs - wideOpen.windowStartMs,
  boundSpanMs: boundEnd - boundStart,
  deepestSpanMs: deep.windowEndMs - deep.windowStartMs,
  minSpanMs: TIMELINE_MIN_SPAN_MS,
  withinBounds: wideOpen.windowStartMs >= boundStart && wideOpen.windowEndMs <= boundEnd
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline zoom", javascriptProbe)
	var zoomResult struct {
		AnchorDriftMs float64 `json:"anchorDriftMs"`
		WidestSpanMs  float64 `json:"widestSpanMs"`
		BoundSpanMs   float64 `json:"boundSpanMs"`
		DeepestSpanMs float64 `json:"deepestSpanMs"`
		MinSpanMs     float64 `json:"minSpanMs"`
		WithinBounds  bool    `json:"withinBounds"`
	}
	if decodeError := json.Unmarshal(probeOutput, &zoomResult); decodeError != nil {
		t.Fatalf("decode timeline zoom behavior: %v (output %q)", decodeError, probeOutput)
	}
	if zoomResult.AnchorDriftMs > 1 {
		t.Fatalf("the anchored instant drifted %.0f ms over three zoom steps; it must stay put",
			zoomResult.AnchorDriftMs)
	}
	if zoomResult.WidestSpanMs != zoomResult.BoundSpanMs {
		t.Fatalf("zooming out settled at %.0f ms, want the full bound span %.0f ms",
			zoomResult.WidestSpanMs, zoomResult.BoundSpanMs)
	}
	if !zoomResult.WithinBounds {
		t.Fatal("the zoomed-out window escaped its bounds")
	}
	if zoomResult.DeepestSpanMs != zoomResult.MinSpanMs {
		t.Fatalf("zooming in settled at %.0f ms, want the %.0f ms floor",
			zoomResult.DeepestSpanMs, zoomResult.MinSpanMs)
	}
}

func TestJavaScriptBehaviorTimelineUserRequestGroupsUseOnlyListedMembers(t *testing.T) {
	indexHTML := generateLiveSite(t)
	javascriptProbe := rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_ROW_HEIGHT") + "\n" +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_GROUP_HEADER_HEIGHT") + "\n" +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_UNKNOWN_USER_REQUEST_NAME") + "\n" +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineFormatSpanMinutes(") + "\n" +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineGroupWindowRows(") + "\n" +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineGroupDetailText(") + "\n" +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineGroupMetricText(") + "\n" +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineFlattenGroups(") + `
var nowMs = Date.UTC(2026, 7, 24, 14, 0);
var windowStartMs = Date.UTC(2026, 7, 24, 7, 0);
var windowEndMs = Date.UTC(2026, 7, 24, 15, 0);
var rows = [
  { id: "REQ-505", claimedTime: "2026-08-24T10:00:00Z", completedTime: "2026-08-24T12:00:00Z", hasWork: true },
  { id: "REQ-504", completedTime: "2026-08-24T12:30:00Z", hasWork: false },
  { id: "REQ-503", claimedTime: "2026-08-24T08:00:00Z", hasWork: true, workOpen: true },
  { id: "REQ-502", claimedTime: "2026-08-24T09:00:00Z", completedTime: "2026-08-24T13:00:00Z", hasWork: true },
  { id: "REQ-501", waitOpen: true, hasWork: false }
];
var requestsById = {
  "REQ-505": { userRequestId: "UR-202" },
  "REQ-504": { userRequestId: "" },
  "REQ-503": { userRequestId: "UR-101" },
  "REQ-502": { userRequestId: "UR-202" },
  "REQ-501": { userRequestId: "" }
};
var samples = [
  { id: "REQ-505", wallMinutes: 120 },
  { id: "REQ-504", wallMinutes: 20 },
  { id: "REQ-503", wallMinutes: 360, excludedReason: "paused" },
  { id: "REQ-502", wallMinutes: -10, excludedReason: "reversed" }
];
function summarize(groups) {
  return groups.map(function (group) {
    return {
      label: group.label,
      ids: group.members.map(function (member) { return member.row.id; }),
      elapsedMinutes: group.elapsedMinutes,
      earliestClaimMs: group.earliestClaimMs,
      latestCompletionMs: group.latestCompletionMs,
      running: group.running,
      acceptedWorkMinutes: group.acceptedWorkMinutes,
      acceptedWorkCount: group.acceptedWorkCount,
      excludedReasons: group.excludedReasons,
      unavailableWorkCount: group.unavailableWorkCount,
      unresolvedCompletionCount: group.unresolvedCompletionCount,
      elapsedUnavailableReason: group.elapsedUnavailableReason,
      metricText: timelineGroupMetricText(group)
    };
  });
}
var allGroups = timelineGroupWindowRows(
  rows, requestsById, samples, nowMs, windowStartMs, windowEndMs);
var narrowRows = [
  { id: "REQ-601", claimedTime: "2026-08-24T09:00:00Z", completedTime: "2026-08-24T13:00:00Z", hasWork: true }
];
var narrowGroups = timelineGroupWindowRows(
  narrowRows,
  { "REQ-601": { userRequestId: "UR-NARROW" } },
  [{ id: "REQ-601", wallMinutes: 240 }],
  nowMs,
  Date.UTC(2026, 7, 24, 10, 0),
  Date.UTC(2026, 7, 24, 12, 0)
);
var unresolvedRows = [
  { id: "REQ-701", claimedTime: "2026-08-24T10:00:00Z", hasWork: true }
];
var mixedRows = [
  { id: "REQ-702", claimedTime: "2026-08-24T09:00:00Z", completedTime: "2026-08-24T11:00:00Z", hasWork: true },
  unresolvedRows[0]
];
var unresolvedRequests = {
  "REQ-701": { userRequestId: "UR-ENDPOINT" },
  "REQ-702": { userRequestId: "UR-ENDPOINT" }
};
process.stdout.write(JSON.stringify({
  all: summarize(allGroups),
  flattenedRowIndexes: timelineFlattenGroups(allGroups).items
    .filter(function (item) { return item.kind === "request"; })
    .map(function (item) { return item.rowIndex; }),
  windowSubset: summarize(timelineGroupWindowRows(
    [rows[3], rows[4]], requestsById, samples, nowMs, windowStartMs, windowEndMs)),
  narrow: summarize(narrowGroups),
  unresolved: summarize(timelineGroupWindowRows(
    unresolvedRows, unresolvedRequests, [], nowMs, windowStartMs, windowEndMs)),
  mixedUnresolved: summarize(timelineGroupWindowRows(
    mixedRows, unresolvedRequests, [{ id: "REQ-702", wallMinutes: 120 }],
    nowMs, windowStartMs, windowEndMs))
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline user-request grouping", javascriptProbe)
	type renderedGroup struct {
		Label                     string         `json:"label"`
		Ids                       []string       `json:"ids"`
		ElapsedMinutes            *float64       `json:"elapsedMinutes"`
		EarliestClaimMS           float64        `json:"earliestClaimMs"`
		LatestCompletionMS        float64        `json:"latestCompletionMs"`
		Running                   bool           `json:"running"`
		AcceptedWorkMinutes       float64        `json:"acceptedWorkMinutes"`
		AcceptedWorkCount         int            `json:"acceptedWorkCount"`
		ExcludedReasons           map[string]int `json:"excludedReasons"`
		UnavailableWorkCount      int            `json:"unavailableWorkCount"`
		UnresolvedCompletionCount int            `json:"unresolvedCompletionCount"`
		ElapsedUnavailableReason  string         `json:"elapsedUnavailableReason"`
		MetricText                string         `json:"metricText"`
	}
	var groupingResult struct {
		All                 []renderedGroup `json:"all"`
		FlattenedRowIndexes []int           `json:"flattenedRowIndexes"`
		WindowSubset        []renderedGroup `json:"windowSubset"`
		Narrow              []renderedGroup `json:"narrow"`
		Unresolved          []renderedGroup `json:"unresolved"`
		MixedUnresolved     []renderedGroup `json:"mixedUnresolved"`
	}
	if decodeError := json.Unmarshal(probeOutput, &groupingResult); decodeError != nil {
		t.Fatalf("decode timeline user-request grouping: %v (output %q)", decodeError, probeOutput)
	}
	if got := []string{groupingResult.All[0].Label, groupingResult.All[1].Label, groupingResult.All[2].Label}; !reflect.DeepEqual(got, []string{"UR-202", "UR-101", "No UR recorded"}) {
		t.Errorf("group order = %v, want first-seen URs with the no-UR group last", got)
	}
	if got := groupingResult.All[0].Ids; !reflect.DeepEqual(got, []string{"REQ-505", "REQ-502"}) {
		t.Errorf("UR-202 members = %v, want newest-first input order", got)
	}
	if !reflect.DeepEqual(groupingResult.FlattenedRowIndexes, []int{0, 3, 2, 1, 4}) {
		t.Errorf("flattened REQ indexes = %v, want grouped display order with headers omitted",
			groupingResult.FlattenedRowIndexes)
	}
	if groupingResult.All[0].ElapsedMinutes == nil || *groupingResult.All[0].ElapsedMinutes != 240 {
		t.Errorf("closed UR elapsed = %v, want 240 minutes from earliest claim to latest completion",
			groupingResult.All[0].ElapsedMinutes)
	}
	if groupingResult.All[0].AcceptedWorkMinutes != 120 || groupingResult.All[0].AcceptedWorkCount != 1 ||
		groupingResult.All[0].ExcludedReasons["reversed"] != 1 {
		t.Errorf("closed UR work verdict = %#v, want one accepted 120-minute sample and one reversed exclusion",
			groupingResult.All[0])
	}
	if !groupingResult.All[1].Running || groupingResult.All[1].ElapsedMinutes == nil ||
		*groupingResult.All[1].ElapsedMinutes != 360 {
		t.Errorf("open UR elapsed = %#v, want six hours to the frozen now and running", groupingResult.All[1])
	}
	if groupingResult.All[2].ElapsedMinutes != nil {
		t.Errorf("no-claim group elapsed = %v, want unavailable rather than created_at fallback",
			*groupingResult.All[2].ElapsedMinutes)
	}
	if !groupingResult.All[2].Running {
		t.Error("the no-claim group contains an open wait but is not visibly classified as running")
	}
	if groupingResult.All[2].AcceptedWorkMinutes != 0 || groupingResult.All[2].UnavailableWorkCount != 2 {
		t.Errorf("no-UR work detail = %#v, want no accepted claim-to-completion work and two unavailable samples",
			groupingResult.All[2])
	}
	if len(groupingResult.WindowSubset) != 2 || len(groupingResult.WindowSubset[0].Ids) != 1 ||
		len(groupingResult.WindowSubset[1].Ids) != 1 {
		t.Errorf("window subset grouped %#v; headers must count only the rows passed after windowing",
			groupingResult.WindowSubset)
	}
	if len(groupingResult.Narrow) != 1 || groupingResult.Narrow[0].ElapsedMinutes == nil ||
		*groupingResult.Narrow[0].ElapsedMinutes != 120 ||
		groupingResult.Narrow[0].AcceptedWorkMinutes != 120 ||
		groupingResult.Narrow[0].EarliestClaimMS != float64(time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC).UnixMilli()) ||
		groupingResult.Narrow[0].LatestCompletionMS != float64(time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC).UnixMilli()) {
		t.Errorf("narrow-window group = %#v, want both the 240-minute claim span and work sample clipped to the window's 120 minutes",
			groupingResult.Narrow)
	}
	for caseName, groups := range map[string][]renderedGroup{
		"isolated": groupingResult.Unresolved,
		"mixed":    groupingResult.MixedUnresolved,
	} {
		if len(groups) != 1 || groups[0].ElapsedMinutes != nil ||
			groups[0].UnresolvedCompletionCount != 1 ||
			groups[0].ElapsedUnavailableReason != "completion endpoint unavailable" ||
			!strings.Contains(groups[0].MetricText, "completion endpoint unavailable") {
			t.Errorf("%s unresolved-completion group = %#v, want no partial elapsed and an explicit endpoint-unavailable reason",
				caseName, groups)
		}
	}
}

func TestJavaScriptBehaviorTimelineNarrowRowsDrawOneMarker(t *testing.T) {
	indexHtml := generateLiveSite(t)
	// The threshold is READ FROM THE RENDERER. Restating 7 here would let the
	// shipped constant drift to 1 with this test still green — REQ-265's lesson,
	// and REQ-322 shipped exactly that mistake in this file.
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPLIT_WIDTH", "TIMELINE_MIN_SEGMENT_WIDTH") +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineWaitEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineWorkEndMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineRowSegments(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineCollapsedRowMark(") + `
var nowMs = Date.UTC(2026, 7, 23, 12, 0);
// Fit all over three months at the shipped plot width: 90 days across 1200px.
var windowStartMs = Date.UTC(2026, 5, 1);
var windowEndMs = Date.UTC(2026, 7, 30);
var plotWidthPx = 1200;
var windowSpanMs = windowEndMs - windowStartMs;
var msPerPixel = windowSpanMs / plotWidthPx;

function markFor(row, projectedRow, spanStartMs, spanEndMs, widthPx) {
  return timelineCollapsedRowMark(
    timelineRowSegments(row, nowMs, projectedRow),
    spanStartMs === undefined ? windowStartMs : spanStartMs,
    spanEndMs === undefined ? windowEndMs : spanEndMs,
    widthPx === undefined ? plotWidthPx : widthPx);
}

// A REQ whose whole life is four pixels wide at this zoom: the wait and the work
// each floor to TIMELINE_MIN_SEGMENT_WIDTH and together claim more room than the
// REQ occupies.
var fourPixelStartMs = Date.UTC(2026, 6, 1);
var narrowRow = {
  id: "REQ-001",
  createdTime: new Date(fourPixelStartMs).toISOString(),
  claimedTime: new Date(fourPixelStartMs + msPerPixel * 2).toISOString(),
  completedTime: new Date(fourPixelStartMs + msPerPixel * 4).toISOString(),
  hasWork: true, waitMinutes: 60, workMinutes: 60
};
// The same REQ at a zoom where four pixels became four hundred.
var wideWindowStartMs = fourPixelStartMs - msPerPixel * 10;
var wideWindowEndMs = fourPixelStartMs + msPerPixel * 14;

// A row with reversed stamps. Its break markers are the point of drawing it, so
// the caller excludes it — but the mark function must not be the thing that
// silently rescues a caller that forgets.
var reversedRow = {
  id: "REQ-002",
  createdTime: new Date(fourPixelStartMs).toISOString(),
  claimedTime: new Date(fourPixelStartMs - msPerPixel * 2).toISOString(),
  completedTime: new Date(fourPixelStartMs + msPerPixel).toISOString(),
  hasWork: true, waitMinutes: -60, workMinutes: 30
};

// A REQ still waiting: ONE drawn segment, so there is no split to withdraw and
// no marker to collapse to, however narrow it is.
var singleSegmentRow = {
  id: "REQ-003",
  createdTime: new Date(nowMs - msPerPixel).toISOString(),
  waitOpen: true, waitMinutes: 5
};

// The unparseable row. timelineRowSegments hands it the -Infinity → Infinity
// sentinel so it is listed in every window; collapsing that to one marker would
// draw a bar across the entire chart.
var unparseableRow = { id: "REQ-004", createdTime: "not-a-date", waitMinutes: 0 };

// A pending REQ whose forecast bar sits right beside its open wait, both inside
// the same handful of pixels. Collapsing it would draw one SOLID marker over a
// span that is mostly projection — work claimed as measured.
var projectedNarrowRow = {
  id: "REQ-005",
  createdTime: new Date(nowMs - msPerPixel).toISOString(),
  waitOpen: true, waitMinutes: 5
};
var narrowProjection = {
  id: "REQ-005",
  startTime: new Date(nowMs).toISOString(),
  endTime: new Date(nowMs + msPerPixel * 2).toISOString()
};

var narrowMark = markFor(narrowRow);
var wideMark = markFor(narrowRow, undefined, wideWindowStartMs, wideWindowEndMs);

process.stdout.write(JSON.stringify({
  splitWidth: TIMELINE_MIN_SPLIT_WIDTH,
  segmentWidth: TIMELINE_MIN_SEGMENT_WIDTH,
  narrowCollapsed: narrowMark !== null,
  narrowMarkStartIso: narrowMark ? new Date(narrowMark.startMs).toISOString() : "",
  narrowMarkEndIso: narrowMark ? new Date(narrowMark.endMs).toISOString() : "",
  narrowRowCreatedIso: narrowRow.createdTime,
  narrowRowCompletedIso: narrowRow.completedTime,
  wideCollapsed: wideMark !== null,
  reversedCollapsed: markFor(reversedRow) !== null,
  singleSegmentCollapsed: markFor(singleSegmentRow) !== null,
  unparseableCollapsed: markFor(unparseableRow) !== null,
  // A row sitting entirely outside the window has no visible width at all, so it
  // is at the collapsing end of the scale rather than the splitting end.
  offscreenCollapsed: markFor(narrowRow, undefined, Date.UTC(2027, 0, 1), Date.UTC(2027, 1, 1)) !== null,
  // Two segments, four pixels: collapsible by width alone. The caller is what
  // has to spare it.
  projectedNarrowCollapsibleByWidth: markFor(projectedNarrowRow, narrowProjection) !== null
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline narrow rows", javascriptProbe)
	var markResult struct {
		SplitWidth             float64 `json:"splitWidth"`
		SegmentWidth           float64 `json:"segmentWidth"`
		NarrowCollapsed        bool    `json:"narrowCollapsed"`
		NarrowMarkStartIso     string  `json:"narrowMarkStartIso"`
		NarrowMarkEndIso       string  `json:"narrowMarkEndIso"`
		NarrowRowCreatedIso    string  `json:"narrowRowCreatedIso"`
		NarrowRowCompletedIso  string  `json:"narrowRowCompletedIso"`
		WideCollapsed          bool    `json:"wideCollapsed"`
		ReversedCollapsed      bool    `json:"reversedCollapsed"`
		SingleSegmentCollapsed bool    `json:"singleSegmentCollapsed"`
		UnparseableCollapsed   bool    `json:"unparseableCollapsed"`
		OffscreenCollapsed     bool    `json:"offscreenCollapsed"`

		ProjectedNarrowCollapsibleByWidth bool `json:"projectedNarrowCollapsibleByWidth"`
	}
	if decodeError := json.Unmarshal(probeOutput, &markResult); decodeError != nil {
		t.Fatalf("decode timeline narrow row behavior: %v (output %q)", decodeError, probeOutput)
	}

	// THE THRESHOLD, pinned to what a two-segment bar physically needs: two
	// floored segments plus a pixel of boundary between them. Dropping
	// TIMELINE_MIN_SPLIT_WIDTH below that would make the collapse fire only where
	// the split was already invisible, which is the defect this REQ removed.
	if wantFloor := 2*markResult.SegmentWidth + 1; markResult.SplitWidth < wantFloor {
		t.Errorf("the split threshold is %g and a readable two-segment bar needs %g "+
			"(2 x %g floored segments plus a boundary); below it the collapse never fires "+
			"where the split is unreadable", markResult.SplitWidth, wantFloor, markResult.SegmentWidth)
	}

	// The pair that makes the rule about WIDTH and not about the row. Same REQ,
	// same stamps, two zooms.
	if !markResult.NarrowCollapsed {
		t.Error("a REQ whose whole span is four pixels wide kept its wait/work split; " +
			"two floored segments in four pixels is a split the pixels cannot carry")
	}
	if markResult.WideCollapsed {
		t.Error("the same REQ collapsed to one marker at a zoom where its span is hundreds " +
			"of pixels wide; the collapse must cost nothing once there is room to split")
	}
	// The marker has to cover the row's real extent, not a floored stub anchored
	// at one end: the reader still reads its POSITION against the gridlines.
	if markResult.NarrowMarkStartIso != markResult.NarrowRowCreatedIso ||
		markResult.NarrowMarkEndIso != markResult.NarrowRowCompletedIso {
		t.Errorf("the collapsed marker covers %s → %s, want the row's own extent %s → %s",
			markResult.NarrowMarkStartIso, markResult.NarrowMarkEndIso,
			markResult.NarrowRowCreatedIso, markResult.NarrowRowCompletedIso)
	}

	if markResult.SingleSegmentCollapsed {
		t.Error("a REQ drawing one segment was reported as collapsible; there is no split " +
			"to withdraw, and collapsing it would replace its open wait with a closed marker")
	}
	// Same guard as the case above, reached by a different route: the sentinel
	// segment timelineRowSegments emits for an unreadable created_at arrives alone,
	// so "one segment" is what spares it. Collapsing it would draw one marker
	// across the whole chart.
	if markResult.UnparseableCollapsed {
		t.Error("a row with an unreadable created_at was reported as collapsible; its segment " +
			"is the -Infinity sentinel, and collapsing it draws one bar across the whole chart")
	}
	if !markResult.OffscreenCollapsed {
		t.Error("a row with no visible width in the window was reported as splittable; " +
			"zero pixels cannot carry two segments")
	}
	// Not an assertion about the mark function's own judgement — it has no way to
	// know — but a record of what the caller must keep doing. renderVisibleRows
	// excludes broken rows before asking, and this pins the reason: asked
	// directly, the function WOULD collapse one.
	if !markResult.ReversedCollapsed {
		t.Error("a reversed-stamp row stopped being collapsible by width alone; if that is " +
			"now handled here, renderVisibleRows's own broken-stamp guard is dead code")
	}
	// The forecast is the second distinction the collapse would erase, and the
	// same pair applies: collapsible by width alone, spared by the caller.
	if !markResult.ProjectedNarrowCollapsibleByWidth {
		t.Error("a narrow pending REQ with a forecast stopped being collapsible by width " +
			"alone; if that is now handled in the mark function, the caller's projection " +
			"guard is dead code")
	}
	// And the other half of that pair: the guard has to still be at the call site.
	// This is a source check rather than a behavioral one because the collapse
	// decision is the pure function above and the exclusion is the caller's — the
	// failure it names is a broken row quietly drawn as one clean marker, with its
	// break markers, the only reason it is on the chart, gone.
	renderVisibleRowsBody := sliceBalancedBlockAfter(t, indexHtml, "function renderVisibleRows(")
	if !strings.Contains(renderVisibleRowsBody, "rowHasBrokenStamps") {
		t.Error("renderVisibleRows no longer excludes broken-stamp rows before asking whether " +
			"to collapse; a narrow reversed row would draw one clean marker instead of the " +
			"break markers that are the reason it is drawn at all")
	}
	if !strings.Contains(renderVisibleRowsBody, "rowHasBrokenStamps || projectedRow") {
		t.Error("renderVisibleRows no longer excludes forecast-carrying rows before asking " +
			"whether to collapse; a narrow pending REQ would draw one solid marker over a span " +
			"that is mostly projection, claiming work that has not happened")
	}
}

func TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineForecastDomStub +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatSpanMinutes(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatStamp(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function clearTimelineForecast(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function renderTimelineForecast(") + `
var confidentProjection = {
  confident: true,
  chainStart: "2026-06-20T12:00:00Z",
  queueEnd: "2026-06-20T14:30:00Z",
  windowSamples: 60, windowSize: 60, minimumSamples: 5,
  normalSamples: 55, normalMinutes: 40,
  trivialSamples: 5, trivialMinutes: 10,
  rows: [{ id: "REQ-401" }, { id: "REQ-402" }],
  excluded: [{ id: "REQ-404", reason: "waiting on an external condition" }]
};
renderTimelineForecast(confidentProjection, false);
var confidentText = collectText(stubNodes["timeline-forecast"]);
var confidentExcludedText = collectText(stubNodes["timeline-excluded"]);

stubNodes["timeline-forecast"] = makeStubNode();
stubNodes["timeline-excluded"] = makeStubNode();
renderTimelineForecast({
  confident: false,
  declinedReason: "only 2 completed REQs inside the read-time rule; 5 are needed before a median means anything",
  rows: [], excluded: [], windowSamples: 2, minimumSamples: 5
}, false);
var declinedText = collectText(stubNodes["timeline-forecast"]);

// The same projection with filters ON. The rows are a subset; the forecast is
// not, and the copy has to say so.
stubNodes["timeline-forecast"] = makeStubNode();
stubNodes["timeline-excluded"] = makeStubNode();
renderTimelineForecast(confidentProjection, true);
var filteredText = collectText(stubNodes["timeline-forecast"]);
var filteredExcludedText = collectText(stubNodes["timeline-excluded"]);

// Declining with filters on carries the same label: the history it declined on
// is the whole queue's, not the subset's.
stubNodes["timeline-forecast"] = makeStubNode();
stubNodes["timeline-excluded"] = makeStubNode();
renderTimelineForecast({
  confident: false,
  declinedReason: "only 2 completed REQs inside the read-time rule; 5 are needed before a median means anything",
  rows: [], excluded: [], windowSamples: 2, minimumSamples: 5
}, true);
var filteredDeclinedText = collectText(stubNodes["timeline-forecast"]);

// The no-rows path clears both nodes without rendering anything: a forecast left
// standing beside "no REQ matches" describes rows that are not on screen.
clearTimelineForecast();
var clearedText = collectText(stubNodes["timeline-forecast"]);
var clearedExcludedText = collectText(stubNodes["timeline-excluded"]);

process.stdout.write(JSON.stringify({
  confidentText: confidentText,
  confidentExcludedText: confidentExcludedText,
  declinedText: declinedText,
  filteredText: filteredText,
  filteredExcludedText: filteredExcludedText,
  filteredDeclinedText: filteredDeclinedText,
  clearedText: clearedText,
  clearedExcludedText: clearedExcludedText
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline forecast", javascriptProbe)
	var forecastResult struct {
		ConfidentText         string `json:"confidentText"`
		ConfidentExcludedText string `json:"confidentExcludedText"`
		DeclinedText          string `json:"declinedText"`
		FilteredText          string `json:"filteredText"`
		FilteredExcludedText  string `json:"filteredExcludedText"`
		FilteredDeclinedText  string `json:"filteredDeclinedText"`
		ClearedText           string `json:"clearedText"`
		ClearedExcludedText   string `json:"clearedExcludedText"`
	}
	if decodeError := json.Unmarshal(probeOutput, &forecastResult); decodeError != nil {
		t.Fatalf("decode timeline forecast behavior: %v (output %q)", decodeError, probeOutput)
	}

	// Each of these is a separate requirement, so each is asserted separately —
	// a single "contains everything" check would report one failure for any of
	// four different regressions.
	for _, wantFragment := range []struct {
		requirement string
		fragment    string
	}{
		{"the end instant itself", "2026-06-20 14:30 UTC"},
		{"the window's sample size", "last 60 completed REQs"},
		{"each bucket's sample count and median", "55 substantive at 40 min"},
		{"the serial assumption", "one REQ at a time"},
		{"the no-parallelism assumption", "no parallel builders"},
		{"the static-queue assumption", "queue that stops growing"},
		{"the read-time rule's exclusions", "Paused and reversed spans are excluded"},
	} {
		if !strings.Contains(forecastResult.ConfidentText, wantFragment.fragment) {
			t.Fatalf("the forecast sentence does not state %s (wanted %q in %q)",
				wantFragment.requirement, wantFragment.fragment, forecastResult.ConfidentText)
		}
	}
	if !strings.Contains(forecastResult.ConfidentExcludedText, "REQ-404") ||
		!strings.Contains(forecastResult.ConfidentExcludedText, "waiting on an external condition") {
		t.Fatalf("the excluded list must name every unschedulable REQ and its reason; got %q",
			forecastResult.ConfidentExcludedText)
	}

	if strings.TrimSpace(forecastResult.ClearedText) != "" ||
		strings.TrimSpace(forecastResult.ClearedExcludedText) != "" {
		t.Fatalf("clearing left forecast %q and excluded %q; a filter matching no rows must leave neither standing beside \"no REQ matches\"",
			forecastResult.ClearedText, forecastResult.ClearedExcludedText)
	}

	if strings.Contains(forecastResult.DeclinedText, "Queue empties") {
		t.Fatalf("thin history produced an end date: %q", forecastResult.DeclinedText)
	}
	if !strings.Contains(forecastResult.DeclinedText, "No end estimate") ||
		!strings.Contains(forecastResult.DeclinedText, "5 are needed") {
		t.Fatalf("declining must say so and carry the reason; got %q", forecastResult.DeclinedText)
	}

	// REQ-305: rows are filtered, the projection never is. With a subset on
	// screen the forecast schedules the whole queue and the excluded list names
	// IDs no visible row carries, so the copy has to name its own population.
	// The label has to read correctly alone, because this paragraph is the one
	// people screenshot.
	if !strings.Contains(forecastResult.FilteredText, "whole queue") {
		t.Errorf("with filters on, the forecast must say it covers the whole queue rather than the rows shown; got %q",
			forecastResult.FilteredText)
	}
	if !strings.Contains(forecastResult.FilteredExcludedText, "whole queue") {
		t.Errorf("with filters on, the excluded list must say it lists the whole queue's exclusions — it names IDs no visible row carries; got %q",
			forecastResult.FilteredExcludedText)
	}
	if !strings.Contains(forecastResult.FilteredDeclinedText, "whole queue") {
		t.Errorf("a declined forecast declined on the whole queue's history, and must say so under a filter too; got %q",
			forecastResult.FilteredDeclinedText)
	}
	// The label is added, never substituted: everything the unfiltered sentence
	// promised is still in the filtered one.
	if !strings.Contains(forecastResult.FilteredText, "Queue empties around") ||
		!strings.Contains(forecastResult.FilteredText, "no parallel builders") {
		t.Errorf("the filtered forecast must still carry the estimate and its assumptions; got %q",
			forecastResult.FilteredText)
	}
	// And the unfiltered copy is untouched — the settled case stays settled.
	if strings.Contains(forecastResult.ConfidentText, "whole queue") ||
		strings.Contains(forecastResult.ConfidentExcludedText, "whole queue") {
		t.Errorf("with no filter active there is nothing to disambiguate, so the label must not appear; got %q / %q",
			forecastResult.ConfidentText, forecastResult.ConfidentExcludedText)
	}
}

func TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened(t *testing.T) {
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
    "REQ-603": { status: "claimed", title: "beta running", domain: "general" },
    "REQ-604": { status: "completed", title: "gamma archived", domain: "general" }
  },
  userRequests: {
    "UR-401": { requestIds: ["REQ-601", "REQ-602"], title: "alpha request", inputFilePresent: true },
    "UR-402": { requestIds: ["REQ-603"], title: "beta request", inputFilePresent: true },
    "UR-403": { requestIds: ["REQ-604"], title: "gamma request", inputFilePresent: false }
  },
  userRequestOrder: ["UR-401", "UR-402", "UR-403"],
  calendar: [
    { id: "REQ-602", completionTime: "2026-08-15T06:00:00Z" },
    { id: "REQ-604", completionTime: "2026-08-01T06:00:00Z" }
  ]
};
var requestsById = boardData.requests;
var userRequestsById = boardData.userRequests;
var viewState = { view: "board", lens: "user-request", windowHours: 24 };
var filterState = { searchText: "", domain: "", status: "", userRequestActivity: "all" };
var userRequestCardsFolded = true;

function makeNode() {
  return {
    childNodes: [],
    dataset: {},
    attributes: {},
    listeners: {},
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
}
var userRequestLensNode = makeNode();
var document = {
  getElementById: function (nodeId) { return nodeId === "user-request-lens" ? userRequestLensNode : null; },
  createElement: function () { return makeNode(); }
};
function makeRequestCard(requestId) { return { className: "req-card", requestId: requestId }; }
` + strings.Join(functionBlocks, "\n") + `
function collectByClassName(node, wantedClassName, found) {
  found = found || [];
  if (node.className === wantedClassName) { found.push(node); }
  (node.childNodes || []).forEach(function (child) { collectByClassName(child, wantedClassName, found); });
  return found;
}
function collectDrawerTriggers(node, found) {
  found = found || [];
  if (node.dataset && node.dataset.detailKind === "ur") { found.push(node.dataset.detailId); }
  (node.childNodes || []).forEach(function (child) { collectDrawerTriggers(child, found); });
  return found;
}
function renderGroups() {
  userRequestLensNode = makeNode();
  renderUserRequestLens();
  return userRequestLensNode.childNodes.filter(function (node) { return node.className === "ur-group"; });
}
function headOf(group) { return collectByClassName(group, "ur-group-head")[0]; }
function describeGroups(groups) {
  return groups.map(function (group) {
    var cardIds = [];
    collectByClassName(group, "ur-group-cards").forEach(function (cardsNode) {
      cardsNode.childNodes.forEach(function (card) { cardIds.push(card.requestId); });
    });
    return {
      userRequestId: collectByClassName(group, "ur-id")[0].textContent,
      expanded: headOf(group).getAttribute("aria-expanded") || "",
      cardIds: cardIds,
      drawerTriggers: collectDrawerTriggers(group)
    };
  });
}

var foldedGroups = renderGroups();
var foldedInitial = describeGroups(foldedGroups);
headOf(foldedGroups[0]).dispatch("click");
var afterOpen = describeGroups(foldedGroups);
headOf(foldedGroups[0]).dispatch("click");
var afterClose = describeGroups(foldedGroups);

filterState.userRequestActivity = "active";
filterState.status = "pending";
var scopedFoldedGroups = renderGroups();
headOf(scopedFoldedGroups[0]).dispatch("click");
var scopedFolded = describeGroups(scopedFoldedGroups);
userRequestCardsFolded = false;
var scopedByUserRequest = describeGroups(renderGroups());

process.stdout.write(JSON.stringify({
  foldedInitial: foldedInitial,
  afterOpen: afterOpen,
  afterClose: afterClose,
  scopedFolded: scopedFolded,
  scopedByUserRequest: scopedByUserRequest
}));
`
	probeOutput := runJavaScriptBehaviorProbe(t, "URs-only fold", javascriptProbe)

	var result struct {
		FoldedInitial       []renderedUserRequestRow `json:"foldedInitial"`
		AfterOpen           []renderedUserRequestRow `json:"afterOpen"`
		AfterClose          []renderedUserRequestRow `json:"afterClose"`
		ScopedFolded        []renderedUserRequestRow `json:"scopedFolded"`
		ScopedByUserRequest []renderedUserRequestRow `json:"scopedByUserRequest"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode URs-only fold output: %v (output %q)", decodeError, probeOutput)
	}

	// Folded: one row per UR, no cards, and the drawer still reachable from each row.
	wantUserRequestIds := []string{"UR-401", "UR-402", "UR-403"}
	if len(result.FoldedInitial) != len(wantUserRequestIds) {
		t.Fatalf("URs only rendered %d rows, want %d: %#v", len(result.FoldedInitial), len(wantUserRequestIds), result.FoldedInitial)
	}
	for rowIndex, row := range result.FoldedInitial {
		if row.UserRequestId != wantUserRequestIds[rowIndex] {
			t.Fatalf("URs only row %d = %q, want %q", rowIndex, row.UserRequestId, wantUserRequestIds[rowIndex])
		}
		if len(row.CardIds) != 0 {
			t.Fatalf("URs only row %s rendered cards %#v before it was opened, want none", row.UserRequestId, row.CardIds)
		}
		if row.Expanded != "false" {
			t.Fatalf("URs only row %s aria-expanded = %q, want \"false\"", row.UserRequestId, row.Expanded)
		}
		if len(row.DrawerTriggers) != 1 || row.DrawerTriggers[0] != row.UserRequestId {
			t.Fatalf("URs only row %s drawer triggers = %#v, want exactly one for itself", row.UserRequestId, row.DrawerTriggers)
		}
	}

	// Opening one row reveals exactly that UR's cards and leaves the others folded.
	openedRow := result.AfterOpen[0]
	if openedRow.Expanded != "true" {
		t.Fatalf("opened row %s aria-expanded = %q, want \"true\"", openedRow.UserRequestId, openedRow.Expanded)
	}
	if strings.Join(openedRow.CardIds, ",") != "REQ-601,REQ-602" {
		t.Fatalf("opened row %s cards = %#v, want REQ-601 and REQ-602", openedRow.UserRequestId, openedRow.CardIds)
	}
	for _, stillFolded := range result.AfterOpen[1:] {
		if len(stillFolded.CardIds) != 0 {
			t.Fatalf("row %s unfolded with a sibling: cards %#v", stillFolded.UserRequestId, stillFolded.CardIds)
		}
	}

	// Activating it again folds it back.
	closedRow := result.AfterClose[0]
	if closedRow.Expanded != "false" || len(closedRow.CardIds) != 0 {
		t.Fatalf("re-activated row = %#v, want aria-expanded false and no cards", closedRow)
	}

	// Active scope plus a status filter must decide identically in both lenses.
	foldedScopedIds := userRequestIdsOf(result.ScopedFolded)
	byUserRequestScopedIds := userRequestIdsOf(result.ScopedByUserRequest)
	if strings.Join(foldedScopedIds, ",") != strings.Join(byUserRequestScopedIds, ",") {
		t.Fatalf("URs only showed %#v under Active+status:pending, by-UR showed %#v; the two lenses must hide the same URs",
			foldedScopedIds, byUserRequestScopedIds)
	}
	if strings.Join(foldedScopedIds, ",") != "UR-401" {
		t.Fatalf("Active+status:pending showed %#v, want only UR-401", foldedScopedIds)
	}
	if strings.Join(result.ScopedFolded[0].CardIds, ",") != "REQ-601" {
		t.Fatalf("opened row under a status filter showed %#v, want only the matching REQ-601", result.ScopedFolded[0].CardIds)
	}
	// The fold is the folded lens's alone: a by-UR head announces no expanded state.
	if result.ScopedByUserRequest[0].Expanded != "" {
		t.Fatalf("by-UR head carries aria-expanded=%q; the fold must not leak into the always-open lens",
			result.ScopedByUserRequest[0].Expanded)
	}
}

func TestJavaScriptBehaviorTimelineTrailingWindowAnchorsBeforeDisplayPadding(t *testing.T) {
	indexHTML := generateLiveSite(t)
	boundsStartToken := "var boundStartMs = Date.parse(timeline.rangeStart);"
	boundsEndToken := "if (!timelineViewState.fitted ||"
	boundsStart := strings.Index(indexHTML, boundsStartToken)
	if boundsStart < 0 {
		t.Fatalf("generated page has no production timeline bounds block beginning %q", boundsStartToken)
	}
	boundsEndOffset := strings.Index(indexHTML[boundsStart:], boundsEndToken)
	if boundsEndOffset < 0 {
		t.Fatalf("generated page has no production timeline bounds block ending %q", boundsEndToken)
	}
	boundsSource := indexHTML[boundsStart : boundsStart+boundsEndOffset]
	applySource := sliceBalancedBlockAfter(t, indexHTML, "function applyTrailingWindow(")

	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS") +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineZoomedWindow(") + "\n" +
		sliceBalancedBlockAfter(t, indexHTML, "function timelineTrailingWindow(") + `
var recordedStartMs = Date.UTC(2026, 0, 1);
var meaningfulEndMs = recordedStartMs + 95 * TIMELINE_DAY_MS;
var nowMs = meaningfulEndMs + 10 * TIMELINE_DAY_MS;
var timeline = {
  rangeStart: new Date(recordedStartMs).toISOString(),
  rangeEnd: new Date(meaningfulEndMs).toISOString()
};
var projection = { queueEnd: "" };
var filterMatchedRows = [{}];
(function () {
` + boundsSource + "\n" + applySource + `
var timelineViewState = { windowStartMs: boundStartMs, windowEndMs: boundEndMs };
applyTrailingWindow("1");
process.stdout.write(JSON.stringify({
  recordedStartMs: recordedStartMs,
  meaningfulEndMs: meaningfulEndMs,
  displayStartMs: boundStartMs,
  displayEndMs: boundEndMs,
  windowStartMs: timelineViewState.windowStartMs,
  windowEndMs: timelineViewState.windowEndMs,
  dayMs: TIMELINE_DAY_MS
}));
}());`

	probeOutput := runJavaScriptBehaviorProbe(t, "production trailing window with display padding", javascriptProbe)
	var result struct {
		RecordedStartMs float64 `json:"recordedStartMs"`
		MeaningfulEndMs float64 `json:"meaningfulEndMs"`
		DisplayStartMs  float64 `json:"displayStartMs"`
		DisplayEndMs    float64 `json:"displayEndMs"`
		WindowStartMs   float64 `json:"windowStartMs"`
		WindowEndMs     float64 `json:"windowEndMs"`
		DayMs           float64 `json:"dayMs"`
	}
	if decodeError := json.Unmarshal(probeOutput, &result); decodeError != nil {
		t.Fatalf("decode production trailing-window padding behavior: %v (output %q)", decodeError, probeOutput)
	}

	if result.DisplayStartMs >= result.RecordedStartMs || result.DisplayEndMs <= result.MeaningfulEndMs+result.DayMs {
		t.Fatalf("production bounds did not retain the 95-day range's cosmetic padding: recorded %.0f→%.0f, display %.0f→%.0f",
			result.RecordedStartMs, result.MeaningfulEndMs, result.DisplayStartMs, result.DisplayEndMs)
	}
	if result.WindowEndMs != result.MeaningfulEndMs || result.WindowStartMs != result.MeaningfulEndMs-result.DayMs {
		t.Fatalf("production Last day landed on %.0f→%.0f, want the day ending at meaningful endpoint %.0f; the display end %.0f is cosmetic padding",
			result.WindowStartMs, result.WindowEndMs, result.MeaningfulEndMs, result.DisplayEndMs)
	}
}

func TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant(t *testing.T) {
	indexHtml := generateLiveSite(t)
	javascriptProbe := timelineProbePreamble(t, "TIMELINE_MIN_SPAN_MS", "TIMELINE_DAY_MS",
		"TIMELINE_AXIS_TICK_COUNT", "TIMELINE_AXIS_TICK_LIMIT") +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_WEEK_ALIGNMENT_MS") + "\n" +
		rendererDeclarationLine(t, "web/board-timeline.js", "TIMELINE_MONTHS") + "\n" +
		rendererBracketDeclaration(t, "web/board-timeline.js", "TIMELINE_AXIS_TICK_STEPS") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickStepSpanMs(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTickStep(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineTickAtOrBefore(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineSteppedTick(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineAxisTicks(") + "\n" +
		sliceBalancedBlockAfter(t, indexHtml, "function timelineFormatAxisTick(") + `
// The REAL tick source and the REAL formatter, called the way renderAxis calls
// them — including passing the gap that positioned the ticks rather than deriving
// one here.
function axisTicks(name, startMs, endMs) {
  var chosen = timelineAxisTicks(startMs, endMs);
  var ticks = chosen.instants.map(function (tickMs) {
    var instant = new Date(tickMs);
    return {
      epochMs: tickMs,
      label: timelineFormatAxisTick(tickMs, chosen.gapMs, startMs, endMs),
      dayOfMonth: instant.getUTCDate(),
      hour: instant.getUTCHours(),
      minute: instant.getUTCMinutes(),
      year: instant.getUTCFullYear()
    };
  });
  return {
    name: name, ticks: ticks, startMs: startMs, endMs: endMs, gapMs: chosen.gapMs,
    // Whether every tick sits on the week boundary the SHIPPED axis aligns to.
    // Both constants are read off the renderer rather than restated here, so
    // moving Monday moves this with it or fails here (REQ-322).
    everyTickOnAWeekBoundary: chosen.instants.every(function (tickMs) {
      return (tickMs - TIMELINE_WEEK_ALIGNMENT_MS) % (7 * TIMELINE_DAY_MS) === 0;
    }),
    tickWeekdays: chosen.instants.map(function (tickMs) { return new Date(tickMs).getUTCDay(); })
  };
}

var mondayMs = Date.UTC(2026, 7, 17);      // 17 Aug 2026 is a Monday
process.stdout.write(JSON.stringify({
  windows: [
    // Where the Now button lands: a window covering the now-line and the
    // forecast's queue-empty instant, which on a healthy queue is well under an
    // hour, so the span settles near the view's floor and the start is wherever
    // "now" fell — 11:26, not the top of an hour.
    axisTicks("Now", Date.UTC(2026, 7, 18, 11, 26), Date.UTC(2026, 7, 18, 11, 26) + TIMELINE_MIN_SPAN_MS),
    axisTicks("Day", Date.UTC(2026, 7, 18), Date.UTC(2026, 7, 19)),
    // A free zoom between the period levels: not a whole number of anything.
    axisTicks("free zoom, four days", Date.UTC(2026, 7, 15), Date.UTC(2026, 7, 19)),
    axisTicks("Week", mondayMs, mondayMs + 7 * TIMELINE_DAY_MS),
    axisTicks("Month", Date.UTC(2026, 7, 1), Date.UTC(2026, 8, 1)),
    axisTicks("Fit all", Date.UTC(2026, 3, 7), Date.UTC(2026, 7, 18)),
    // Three months, which is what Fit all measures on this repo's own board. It
    // picks the FORTNIGHT rung, and is the only fixture here that does — without
    // it the week-boundary alignment of that rung is never checked.
    axisTicks("Fit all, three months", Date.UTC(2026, 4, 27), Date.UTC(2026, 7, 25)),
    // Fit all is the whole capture history, and it only grows. Once it crosses a
    // calendar year one day-and-month comes round twice.
    axisTicks("Fit all across two years", Date.UTC(2025, 7, 18), Date.UTC(2027, 7, 18)),
    // Nine days, crossing New Year. Shorter than a year and still ambiguous
    // without one — the case the old spanMs >= TIMELINE_YEAR_MS threshold missed.
    axisTicks("across New Year", Date.UTC(2026, 11, 28), Date.UTC(2027, 0, 6))
  ]
}));`

	probeOutput := runJavaScriptBehaviorProbe(t, "timeline axis labels", javascriptProbe)
	var axisResult struct {
		Windows []struct {
			Name  string `json:"name"`
			Ticks []struct {
				EpochMs    float64 `json:"epochMs"`
				Label      string  `json:"label"`
				DayOfMonth int     `json:"dayOfMonth"`
				Hour       int     `json:"hour"`
				Minute     int     `json:"minute"`
				Year       int     `json:"year"`
			} `json:"ticks"`
			StartMs                  float64 `json:"startMs"`
			EndMs                    float64 `json:"endMs"`
			GapMs                    float64 `json:"gapMs"`
			EveryTickOnAWeekBoundary bool    `json:"everyTickOnAWeekBoundary"`
			TickWeekdays             []int   `json:"tickWeekdays"`
		} `json:"windows"`
	}
	if decodeError := json.Unmarshal(probeOutput, &axisResult); decodeError != nil {
		t.Fatalf("decode timeline axis behavior: %v (output %q)", decodeError, probeOutput)
	}

	// What each window's labels have to look like. The three period windows and
	// Fit all are here to hold their EXISTING labels: this is a formatting fix,
	// and the formats that were already right may not move.
	const (
		axisLabelWithTime = "date and time"
		axisLabelDateOnly = "date alone"
		axisLabelWithYear = "date and year"
	)
	wantAxisLabelShape := map[string]string{
		"Now":                   axisLabelWithTime,
		"Day":                   axisLabelWithTime,
		"free zoom, four days":  axisLabelDateOnly,
		"Week":                  axisLabelDateOnly,
		"Month":                 axisLabelDateOnly,
		"Fit all":               axisLabelDateOnly,
		"Fit all, three months": axisLabelDateOnly,
		// Both of these cross a calendar year, which is what earns the year — not
		// being longer than 365 days. "across New Year" is nine days long.
		"Fit all across two years": axisLabelWithYear,
		"across New Year":          axisLabelWithYear,
	}
	if len(axisResult.Windows) != len(wantAxisLabelShape) {
		t.Fatalf("the probe drove %d windows, want the %d named", len(axisResult.Windows), len(wantAxisLabelShape))
	}

	weekGapsSeen := 0
	for _, window := range axisResult.Windows {
		labelShape, isNamed := wantAxisLabelShape[window.Name]
		if !isNamed {
			t.Fatalf("the probe drove an unnamed window %q", window.Name)
		}
		renderedLabels := make([]string, 0, len(window.Ticks))
		distinctLabels := map[string]bool{}
		for _, tick := range window.Ticks {
			renderedLabels = append(renderedLabels, tick.Label)
			distinctLabels[tick.Label] = true
		}
		// Two ticks at different instants reading the same label is what makes
		// the axis unreadable rather than merely imprecise.
		if len(distinctLabels) != len(window.Ticks) {
			t.Fatalf("the %s window draws %d ticks with only %d distinct labels: %q",
				window.Name, len(window.Ticks), len(distinctLabels), renderedLabels)
		}
		// THE TICKS THEMSELVES. Ascending, inside the window, and there at all —
		// without this the label assertions below would pass over an empty axis.
		if len(window.Ticks) < 3 {
			t.Fatalf("the %s window drew %d ticks; an axis with fewer than three is not one",
				window.Name, len(window.Ticks))
		}
		for tickIndex, tick := range window.Ticks {
			if tick.EpochMs < window.StartMs || tick.EpochMs > window.EndMs {
				t.Fatalf("the %s window drew a tick at %s, outside the window it describes",
					window.Name, tick.Label)
			}
			if tickIndex > 0 && tick.EpochMs <= window.Ticks[tickIndex-1].EpochMs {
				t.Fatalf("the %s window's ticks are not strictly ascending at %q", window.Name, tick.Label)
			}
		}
		// A WEEK-LONG GAP LANDS ON THE WEEK BOUNDARY TIMELINE_WEEK_ALIGNMENT_MS NAMES.
		// Aligning it to the epoch instead gives Thursdays — still midnights, still
		// distinct, still inside the window, so every other assertion here passes
		// and the axis silently draws a week that starts on the wrong day.
		const oneWeekMs = 7 * 24 * 60 * 60 * 1000
		if window.GapMs == oneWeekMs || window.GapMs == 2*oneWeekMs {
			weekGapsSeen++
			if !window.EveryTickOnAWeekBoundary {
				t.Fatalf("the %s window uses a %.0f-day gap but its ticks fall on weekdays %v; a "+
					"week-long gap has to land on TIMELINE_WEEK_ALIGNMENT_MS",
					window.Name, window.GapMs/float64(24*60*60*1000), window.TickWeekdays)
			}
		}
		// A LABEL WITH NO TIME IS A CLAIM OF MIDNIGHT. This is the assertion the old
		// version of this test did not make, and the whole reason the week axis could
		// print "9 Jul" for a tick at 9 Jul 12:00 with the suite green.
		if labelShape != axisLabelWithTime {
			for _, tick := range window.Ticks {
				if tick.Hour != 0 || tick.Minute != 0 {
					t.Fatalf("the %s window labels the tick at %02d:%02d on the %dth as %q, with no "+
						"time in it — a date-only label claims midnight, so a tick that is not at "+
						"midnight may not have one",
						window.Name, tick.Hour, tick.Minute, tick.DayOfMonth, tick.Label)
				}
			}
		}
		// Every number in the label has to be one the instant carries. Matching
		// the whole label also pins the shape, so a window cannot quietly gain
		// or lose a component the reader relies on.
		for _, tick := range window.Ticks {
			var wantLabelPattern string
			switch labelShape {
			case axisLabelWithTime:
				wantLabelPattern = fmt.Sprintf(`^%d [A-Z][a-z]{2} %02d:%02d$`, tick.DayOfMonth, tick.Hour, tick.Minute)
			case axisLabelDateOnly:
				wantLabelPattern = fmt.Sprintf(`^%d [A-Z][a-z]{2}$`, tick.DayOfMonth)
			case axisLabelWithYear:
				wantLabelPattern = fmt.Sprintf(`^%d [A-Z][a-z]{2} %d$`, tick.DayOfMonth, tick.Year)
			}
			if !regexp.MustCompile(wantLabelPattern).MatchString(tick.Label) {
				t.Fatalf("the %s window renders the tick at %d/%02d:%02d/%d as %q, want %s matching %s",
					window.Name, tick.DayOfMonth, tick.Hour, tick.Minute, tick.Year,
					tick.Label, labelShape, wantLabelPattern)
			}
		}
	}
	// The week-boundary assertion above is inside a conditional, so it is worth
	// nothing if no fixture window ever picks a week-long gap. The Month window is
	// the one that does; if the ladder is re-tuned so none of them do, this says so
	// instead of quietly passing.
	// Both rungs, counted separately: one fixture hitting the 7-day rung leaves the
	// 14-day rung's alignment unchecked, which is exactly how a mutation of it
	// passed the first time this ran.
	if weekGapsSeen < 2 {
		t.Errorf("only %d fixture window(s) chose a week-length gap; both the 7-day and the "+
			"14-day rung need one, or the alignment of the unvisited rung is never checked",
			weekGapsSeen)
	}
}

func TestJavaScriptBehaviorTimelineRefusesToRenderAgainstAnUnmeasurableHost(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	// Twelve rows, so an eight-row viewport is visibly a truncation rather than a
	// coincidence, spread over four hours so a 120px plot is visibly a crush.
	timelinePayload := `{
	  "now": "2026-08-18T13:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-901","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:10:00Z","completedTime":"2026-08-18T09:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-902","createdTime":"2026-08-18T09:20:00Z","claimedTime":"2026-08-18T09:30:00Z","completedTime":"2026-08-18T10:00:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-903","createdTime":"2026-08-18T09:40:00Z","claimedTime":"2026-08-18T09:50:00Z","completedTime":"2026-08-18T10:20:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-904","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:10:00Z","completedTime":"2026-08-18T10:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-905","createdTime":"2026-08-18T10:20:00Z","claimedTime":"2026-08-18T10:30:00Z","completedTime":"2026-08-18T11:00:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-906","createdTime":"2026-08-18T10:40:00Z","claimedTime":"2026-08-18T10:50:00Z","completedTime":"2026-08-18T11:20:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-907","createdTime":"2026-08-18T11:00:00Z","claimedTime":"2026-08-18T11:10:00Z","completedTime":"2026-08-18T11:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-908","createdTime":"2026-08-18T11:20:00Z","claimedTime":"2026-08-18T11:30:00Z","completedTime":"2026-08-18T12:00:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-909","createdTime":"2026-08-18T11:40:00Z","claimedTime":"2026-08-18T11:50:00Z","completedTime":"2026-08-18T12:20:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-910","createdTime":"2026-08-18T12:00:00Z","claimedTime":"2026-08-18T12:10:00Z","completedTime":"2026-08-18T12:40:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-911","createdTime":"2026-08-18T12:10:00Z","claimedTime":"2026-08-18T12:20:00Z","completedTime":"2026-08-18T12:50:00Z","waitMinutes":10,"workMinutes":30,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-912","createdTime":"2026-08-18T12:20:00Z","claimedTime":"2026-08-18T12:30:00Z","completedTime":"2026-08-18T12:55:00Z","waitMinutes":10,"workMinutes":25,"waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false}
	  ]
	}`

	// Render measurable, then unmeasurable, then measurable again — the three
	// states a reader passes through when they resize the browser on another view
	// and come back. Both renders after the first go through the same entry point
	// the resize listener uses.
	probeDriver := `
function drawnRowIds() {
  var ids = [];
  (function walk(node) {
    (node.children || []).forEach(function (childNode) {
      var attributes = childNode.attributes || {};
      if (childNode.stubName === "g" && attributes["data-detail-id"]) {
        ids.push(attributes["data-detail-id"]);
        return;
      }
      walk(childNode);
    });
  })(timelineStubHosts["timeline-scroll"]);
  return ids;
}
function countDescendants(node, stubName) {
  var found = 0;
  (node.children || []).forEach(function (childNode) {
    if (childNode.stubName === stubName) { found++; }
    found += countDescendants(childNode, stubName);
  });
  return found;
}
function hostSize(widthPx, heightPx) {
  var host = timelineStubHosts["timeline-scroll"];
  host.clientWidth = widthPx;
  host.clientHeight = heightPx;
  host.getBoundingClientRect = function () {
    return { width: widthPx, height: heightPx, left: 0, top: 0 };
  };
}
// The stub's textContent is a plain property, so the renderer's own
// "scrollHost.textContent = \"\"" does not drop its children — every other probe
// in this lane renders once and never notices. Three renders do, so the fixture
// clears what a real DOM would have cleared. This is stub bookkeeping, not a
// production behaviour: getting it wrong made the second render's row list read
// as the first render's.
function clearRenderedHosts() {
  ["timeline-scroll", "timeline-axis", "timeline-table-body"].forEach(function (hostId) {
    timelineStubHosts[hostId].children = [];
  });
  timelineStubHosts["timeline-summary"].textContent = "";
}
function snapshot() {
  return {
    rowIds: drawnRowIds(),
    summary: timelineStubHosts["timeline-summary"].textContent,
    // The axis SVG SHELL is created by renderTimelineView before renderAll runs,
    // so counting the host's children counts a container that is always there.
    // What matters is whether any tick was drawn INTO it.
    axisTicks: countDescendants(timelineStubHosts["timeline-axis"], "text")
  };
}
hostSize(900, 400);
renderTimelineView();
var measured = snapshot();

hostSize(0, 0);
clearRenderedHosts();
renderTimelineView();
var unmeasurable = snapshot();

hostSize(900, 400);
clearRenderedHosts();
renderTimelineView();
var remeasured = snapshot();

process.stdout.write(JSON.stringify({
  measured: measured, unmeasurable: unmeasurable, remeasured: remeasured
}));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline unmeasurable host", javascriptProbe)

	type renderSnapshot struct {
		RowIds    []string `json:"rowIds"`
		Summary   string   `json:"summary"`
		AxisTicks int      `json:"axisTicks"`
	}
	var hostResult struct {
		Measured     renderSnapshot `json:"measured"`
		Unmeasurable renderSnapshot `json:"unmeasurable"`
		Remeasured   renderSnapshot `json:"remeasured"`
	}
	if decodeError := json.Unmarshal(probeOutput, &hostResult); decodeError != nil {
		t.Fatalf("decode timeline unmeasurable-host behavior: %v (output %q)", decodeError, probeOutput)
	}

	// SETUP, ASSERTED: without a full first render there is nothing for the
	// unmeasurable render to be compared against.
	if hostResult.Measured.AxisTicks == 0 {
		t.Fatal("the measurable render drew no axis tick labels, so the zero-tick assertion below " +
			"would pass against any render at all")
	}
	if len(hostResult.Measured.RowIds) != 12 {
		t.Fatalf("the measurable render drew %d rows, want all 12 fixture rows; the probe is not "+
			"measuring a full chart", len(hostResult.Measured.RowIds))
	}
	if !strings.Contains(hostResult.Measured.Summary, "12 REQs in the window") {
		t.Fatalf("the measurable render's summary is %q, want it to name all 12 rows", hostResult.Measured.Summary)
	}

	// THE DEFECT. Eight rows and a 120px plot are not a smaller truth, they are a
	// measurement of a box that does not exist.
	if len(hostResult.Unmeasurable.RowIds) != 0 {
		t.Fatalf("rendering against a zero-width host drew %d rows (%v); an unmeasurable host is "+
			"\"not yet\", not a 120px plot with an eight-row viewport",
			len(hostResult.Unmeasurable.RowIds), hostResult.Unmeasurable.RowIds)
	}
	if hostResult.Unmeasurable.Summary != "" {
		t.Fatalf("rendering against a zero-width host wrote the summary %q; it must not describe a "+
			"window it could not lay out", hostResult.Unmeasurable.Summary)
	}
	if hostResult.Unmeasurable.AxisTicks != 0 {
		t.Fatalf("rendering against a zero-width host drew %d axis tick labels", hostResult.Unmeasurable.AxisTicks)
	}

	// And the numbers come back whole once the host has a box again, which is what
	// the ResizeObserver triggers in a real engine.
	if len(hostResult.Remeasured.RowIds) != len(hostResult.Measured.RowIds) {
		t.Fatalf("after the host regained its box the render drew %d rows, want the %d it drew before",
			len(hostResult.Remeasured.RowIds), len(hostResult.Measured.RowIds))
	}
	if hostResult.Remeasured.Summary != hostResult.Measured.Summary {
		t.Fatalf("after the host regained its box the summary reads %q, want the %q it read before",
			hostResult.Remeasured.Summary, hostResult.Measured.Summary)
	}
}

func TestJavaScriptBehaviorTimelineNoMatchStateRetiresTheToolbar(t *testing.T) {
	rendererFragment, readError := embeddedWebAssets.ReadFile("web/board-timeline.js")
	if readError != nil {
		t.Fatalf("read web/board-timeline.js: %v", readError)
	}

	timelinePayload := `{
	  "now": "2026-08-18T12:00:00Z",
	  "rangeStart": "2026-08-18T09:00:00Z",
	  "rangeEnd": "2026-08-18T13:00:00Z",
	  "rows": [
	    {"id":"REQ-921","createdTime":"2026-08-18T09:00:00Z","claimedTime":"2026-08-18T09:30:00Z",
	     "completedTime":"2026-08-18T10:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false},
	    {"id":"REQ-922","createdTime":"2026-08-18T10:00:00Z","claimedTime":"2026-08-18T10:30:00Z",
	     "completedTime":"2026-08-18T11:00:00Z","waitMinutes":30,"workMinutes":30,
	     "waitOpen":false,"workOpen":false,"hasWork":true,"anomaly":false}
	  ]
	}`

	// The stub's querySelectorAll has to answer the control selectors, so the
	// controls are real stub nodes the driver can inspect and press.
	probeDriver := `
var stubControls = ["period-prev", "period-day", "zoom-fit", "range-start"].map(function (name) {
  var control = makeStubNode(name === "range-start" ? "input" : "button");
  control.controlName = name;
  control.disabled = false;
  control.onclick = null;
  return control;
});
document.querySelectorAll = function (selector) {
  if (String(selector).indexOf(".timeline-periods button") !== -1) { return stubControls; }
  if (String(selector).indexOf("[data-timeline-period]") !== -1) { return []; }
  return [];
};
document.getElementById = function (nodeId) {
  if (nodeId === "timeline-zoom-fit") { return stubControls[2]; }
  return timelineStubHosts[nodeId] || null;
};

function controlState() {
  return stubControls.map(function (control) {
    return { name: control.controlName, wired: typeof control.onclick === "function", disabled: !!control.disabled };
  });
}

timelineStubVisibleIds = null;
renderTimelineView();
var matched = { controls: controlState(), summary: timelineStubHosts["timeline-summary"].textContent };

// Nothing matches. The early return fires.
["timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
 "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state",
 "board-main"
].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
timelineStubVisibleIds = [];
renderTimelineView();
var noMatch = { controls: controlState(), summary: timelineStubHosts["timeline-summary"].textContent };

// And back: a filter that matches again must restore every control.
["timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
 "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state",
 "board-main"
].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
timelineStubVisibleIds = null;
renderTimelineView();
var matchedAgain = { controls: controlState(), summary: timelineStubHosts["timeline-summary"].textContent };

process.stdout.write(JSON.stringify({ matched: matched, noMatch: noMatch, matchedAgain: matchedAgain }));
`

	javascriptProbe := timelineRenderDomStubPreamble +
		"var boardData = { timeline: " + timelinePayload + " };\n" +
		string(rendererFragment) +
		probeDriver
	probeOutput := runJavaScriptBehaviorProbe(t, "timeline no-match toolbar", javascriptProbe)

	type controlState struct {
		Name     string `json:"name"`
		Wired    bool   `json:"wired"`
		Disabled bool   `json:"disabled"`
	}
	type renderState struct {
		Controls []controlState `json:"controls"`
		Summary  string         `json:"summary"`
	}
	var toolbarResult struct {
		Matched      renderState `json:"matched"`
		NoMatch      renderState `json:"noMatch"`
		MatchedAgain renderState `json:"matchedAgain"`
	}
	if decodeError := json.Unmarshal(probeOutput, &toolbarResult); decodeError != nil {
		t.Fatalf("decode timeline no-match toolbar behavior: %v (output %q)", decodeError, probeOutput)
	}

	// SETUP, ASSERTED: without these the checks below are measuring nothing.
	if len(toolbarResult.Matched.Controls) == 0 {
		t.Fatal("the probe found no controls, so nothing below is measured")
	}
	if !strings.Contains(toolbarResult.NoMatch.Summary, "No REQ matches the current filters") {
		t.Fatalf("the second render did not take the no-match path (summary %q)",
			toolbarResult.NoMatch.Summary)
	}
	if !strings.Contains(toolbarResult.Matched.Summary, "REQs in the window") {
		t.Fatalf("the first render did not draw a chart (summary %q)", toolbarResult.Matched.Summary)
	}
	// The "still carries a handler" check below is vacuous for any control this stub
	// never wires, so at least one has to arrive wired or that half proves nothing.
	wiredAfterAMatchingRender := 0
	for _, control := range toolbarResult.Matched.Controls {
		if control.Wired {
			wiredAfterAMatchingRender++
		}
	}
	if wiredAfterAMatchingRender == 0 {
		t.Fatal("no control was wired by the matching render, so the handler half of this test " +
			"cannot fail; give the stub's getElementById a control the renderer wires")
	}

	for _, control := range toolbarResult.NoMatch.Controls {
		if control.Wired {
			t.Errorf("after the no-match render the %s control still carries a handler; it belongs "+
				"to the previous render, whose rows the filter excluded", control.Name)
		}
		if !control.Disabled {
			t.Errorf("after the no-match render the %s control is still pressable; a control that "+
				"cannot act must say so rather than doing nothing", control.Name)
		}
	}
	for _, control := range toolbarResult.MatchedAgain.Controls {
		if control.Disabled {
			t.Errorf("the %s control is still disabled after the filter matched again", control.Name)
		}
	}
}

func TestJavaScriptBehaviorDetailRendersOnlyObservedPhaseBreakdown(t *testing.T) {
	boardData := buildImplementationSpanFixturePayload(t)
	phasePayload, encodeError := json.Marshal(boardData.Requests["REQ-901"].PhaseBreakdown)
	if encodeError != nil {
		t.Fatalf("encode phase payload: %v", encodeError)
	}
	indexHtml := generateLiveSite(t)
	functionBlocks := []string{
		sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
		sliceBalancedBlockAfter(t, indexHtml, "function appendPhaseBreakdownRows("),
		sliceBalancedBlockAfter(t, indexHtml, "function formatElapsedDuration("),
	}
	javascriptProbe := `
function makeNode(tagName) {
  return {
    tagName: tagName,
    className: "",
    textContent: "",
    childNodes: [],
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; }
  };
}
var document = {
  createElement: function (tagName) { return makeNode(tagName); },
  createTextNode: function (text) { return { nodeType: "text", text: text, childNodes: [] }; }
};
var futureInstantSkewAllowanceMs = 120000;
var clockSkewMarkerText = "clock skew";
function makeInstantWithRelativeNode(isoText) { return document.createTextNode(isoText); }
function formatShortInstant(isoText) { return isoText; }
function nodeText(node) {
  if (node.nodeType === "text") { return node.text; }
  return (node.textContent || "") + node.childNodes.map(nodeText).join("");
}
var rows = [];
function appendMetaRow(label, value) { rows.push({ label: label, value: typeof value === "string" ? value : nodeText(value) }); }
` + strings.Join(functionBlocks, "\n") + `
appendPhaseBreakdownRows(` + string(phasePayload) + `);
process.stdout.write(JSON.stringify(rows));`

	probeOutput := runJavaScriptBehaviorProbe(t, "detail phase breakdown", javascriptProbe)
	var rows []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}
	if decodeError := json.Unmarshal(probeOutput, &rows); decodeError != nil {
		t.Fatalf("decode phase rows: %v (output %q)", decodeError, probeOutput)
	}
	wantLabels := []string{
		"Phase · Planning", "Phase · Dispatch", "Phase · Builder handback", "Phase · Integration",
		"Phase · Review", "Phase · Completed", "Phase · Release",
	}
	if len(rows) != len(wantLabels) {
		t.Fatalf("rendered %d phase rows, want %d: %#v", len(rows), len(wantLabels), rows)
	}
	for index, wantLabel := range wantLabels {
		if rows[index].Label != wantLabel {
			t.Errorf("phase row %d label = %q, want %q", index, rows[index].Label, wantLabel)
		}
		if !strings.Contains(rows[index].Value, " wall since ") {
			t.Errorf("phase row %d omits honest wall-span wording: %q", index, rows[index].Value)
		}
	}
	if !strings.Contains(rows[len(rows)-1].Value, "wall since Completed") {
		t.Errorf("release tail was not separated from completion: %q", rows[len(rows)-1].Value)
	}
}
