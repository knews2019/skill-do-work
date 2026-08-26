package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// citationFixtureRequestIds and citationFixtureUserRequestIds are the board every
// test below resolves against. UR-001-REQ-042 and UR-002-REQ-042 share a REQ
// segment on purpose: "REQ-042" is the ambiguous case the resolver must refuse
// to guess at, and it is only ambiguous while BOTH are present.
var citationFixtureRequestIds = []string{
	"REQ-1679", "REQ-1108", "REQ-1685", "REQ-500", "REQ-501",
	"UR-001-REQ-042", "UR-002-REQ-042", "UR-003-REQ-077",
}

var citationFixtureUserRequestIds = []string{"UR-074", "UR-075"}

func newCitationFixtureResolver() *ticketMentionResolver {
	return newTicketMentionResolver(citationFixtureRequestIds, citationFixtureUserRequestIds)
}

// describeTicketMentions renders an emitted index one line per mention, with the
// text each offset actually points at sliced back out of the document. A bare
// list of numbers proves nothing about whether the numbers are right; slicing
// the document with them does, and it is what a failure message needs to be
// readable.
func describeTicketMentions(documentText string, mentions []generatedTicketMention) []string {
	documentRunes := utf16RunesOf(documentText)
	described := make([]string, 0, len(mentions))
	for _, mention := range mentions {
		expandMarker := "quoted"
		if mention.Expand {
			expandMarker = "EXPAND"
		}
		kindLabel := mention.Kind
		if kindLabel == "" {
			kindLabel = "missing"
		}
		described = append(described, fmt.Sprintf("%s %s %s %q",
			kindLabel, mention.Id, expandMarker,
			string(documentRunes[mention.Offset:mention.Offset+mention.Length])))
	}
	return described
}

// utf16RunesOf re-expresses a Go string as the code-unit sequence a JavaScript
// string is, so a test can index it the way the client will. Surrogate pairs are
// split, which is exactly the point: a byte offset and a UTF-16 offset diverge
// at the first non-ASCII character, and this is how a test can see it.
func utf16RunesOf(sourceText string) []rune {
	var codeUnits []rune
	for _, currentRune := range sourceText {
		if currentRune > 0xFFFF {
			// Two units. Their individual values never matter here — only that
			// the count matches what the client indexes.
			codeUnits = append(codeUnits, 0xD800, 0xDC00)
			continue
		}
		codeUnits = append(codeUnits, currentRune)
	}
	return codeUnits
}

// Every construct the deleted client scanner got wrong, in one document, with
// the emitted index pinned line by line. The six rows of REQ-383's table are
// here: a blockquoted fence, a backtick in a fence's info string, a fence opened
// as a list item, a code span crossing a newline, a four-space indented block,
// and a link reference definition. Five of them the scanner reached and read
// wrongly; the indented block it could not see at all, because it matches fence
// LINES and an indented block has none.
func TestCollectDocumentTicketMentionsClassifiesEveryQuotedConstruct(t *testing.T) {
	documentText := "---\n" +
		"id: REQ-500\n" +
		"depends_on: [REQ-1679]\n" +
		"---\n" +
		"\n" +
		"Prose cites REQ-1679 and `REQ-1108` in a span, and REQ-1679 again.\n" +
		"\n" +
		"> ````text\n" +
		"> Quoted REQ-1685 verbatim, plus dead REQ-8888.\n" +
		"> ````\n" +
		"\n" +
		"A paragraph closing the blockquote.\n" +
		"\n" +
		"- ```yaml\n" +
		"  depends_on: [REQ-1685]\n" +
		"  ```\n" +
		"\n" +
		"A paragraph closing the list.\n" +
		"\n" +
		"    indented REQ-1108 block, plus dead REQ-8887\n" +
		"\n" +
		"```yaml REQ-1685-template\n" +
		"depends_on: [REQ-500]\n" +
		"```\n" +
		"\n" +
		"```lang`invalid\n" +
		"That line is prose, so REQ-501 under it is a real reference.\n" +
		"\n" +
		"the example reads\n" +
		"`- REQ-500 a worked example — the second\n" +
		"finding for REQ-1108` matters.\n" +
		"\n" +
		"[REQ-1685]: do-work/queue/REQ-1685-x.md\n" +
		"\n" +
		"Dead REQ-9999 in prose, ambiguous REQ-042, and short REQ-077.\n"

	wantMentions := []string{
		`req REQ-1679 EXPAND "REQ-1679"`,      // prose, first mention
		`req REQ-1108 quoted "REQ-1108"`,      // code span: glossed, never expanded
		`req REQ-1679 quoted "REQ-1679"`,      // prose repeat: one expansion per id
		`req REQ-1685 quoted "REQ-1685"`,      // blockquoted fence — the containment contract's preserved words
		`req REQ-1685 quoted "REQ-1685"`,      // list-item fence
		`req REQ-1108 quoted "REQ-1108"`,      // four-space indented block
		`req REQ-1685 quoted "REQ-1685"`,      // a fence's INFO STRING is not prose either
		`req REQ-500 quoted "REQ-500"`,        // that fence's contents
		`req REQ-501 EXPAND "REQ-501"`,        // a backtick in the info string makes that line prose
		`req REQ-500 quoted "REQ-500"`,        // code span, opening line
		`req REQ-1108 quoted "REQ-1108"`,      // same span, continuation line
		`missing REQ-9999 quoted "REQ-9999"`,  // dead id in prose earns an appendix line
		`req UR-003-REQ-077 EXPAND "REQ-077"`, // short segment resolves to the one card carrying it
	}

	mentions := collectDocumentTicketMentions(documentText, newCitationFixtureResolver())
	gotMentions := describeTicketMentions(documentText, mentions)
	if !reflect.DeepEqual(gotMentions, wantMentions) {
		t.Errorf("emitted mentions:\n got  %s\n want %s",
			strings.Join(gotMentions, "\n       "), strings.Join(wantMentions, "\n       "))
	}
}

// The client indexes a JavaScript string, where an offset counts UTF-16 code
// units. A Go byte offset is a different number the moment a body carries an em
// dash — which every REQ body in this repo does — and splicing at one lands
// mid-word. Neither half of this is theoretical: the fixture puts an astral
// emoji and an em dash in front of the mention, and an emoji in the fence, so a
// byte offset misses by six.
func TestCollectDocumentTicketMentionsMeasuresOffsetsInUtf16CodeUnits(t *testing.T) {
	// "---\n" (4) + "title: 😀 dash\n" (7 + 2 + 6 = 15) + "---\n" (4) = 23 units.
	frontmatterText := "---\ntitle: 😀 dash\n---\n"
	// "Emoji " (6) + "😀" (2) + " and an em dash " (16) + "—" (1) + " before " (8) = 33 units.
	bodyText := "Emoji 😀 and an em dash — before REQ-1679.\n"
	documentText := frontmatterText + bodyText

	const wantOffset = 23 + 33
	mentions := collectDocumentTicketMentions(documentText, newCitationFixtureResolver())
	if len(mentions) != 1 {
		t.Fatalf("mentions = %#v, want exactly one", mentions)
	}
	if mentions[0].Offset != wantOffset {
		t.Errorf("offset = %d, want %d UTF-16 code units", mentions[0].Offset, wantOffset)
	}
	if byteOffset := strings.Index(documentText, "REQ-1679"); byteOffset == wantOffset {
		t.Fatalf("the fixture carries no character that separates bytes from code units at %d — the assertion above cannot fail", byteOffset)
	}
	// The offset is only right if it slices the mention back out of the string
	// the CLIENT holds, which is the code-unit sequence, not the bytes.
	documentRunes := utf16RunesOf(documentText)
	if sliced := string(documentRunes[mentions[0].Offset : mentions[0].Offset+mentions[0].Length]); sliced != "REQ-1679" {
		t.Errorf("slicing the client's string at the shipped offset gave %q, want %q", sliced, "REQ-1679")
	}
}

// Offsets are measured from the start of the WHOLE document, so they carry the
// frontmatter fence in front of them. That arithmetic is the one place this
// change can silently corrupt a paste — an offset short by the fence length
// splices a title into the YAML — so it gets its own test rather than riding
// along in another. All three shapes splitFrontmatter recognises are here,
// because each yields a different body start.
func TestCollectDocumentTicketMentionsOffsetsCarryTheFrontmatterFence(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		documentText string
	}{
		{"lf fence", "---\nid: REQ-500\n---\nBody cites REQ-1679 here.\n"},
		{"crlf fence", "---\r\nid: REQ-500\r\n---\r\nBody cites REQ-1679 here.\r\n"},
		{"no fence at all", "# Notes\n\nBody cites REQ-1679 here.\n"},
		{"unclosed fence is all body", "---\nstatus: pending\nBody cites REQ-1679 here.\n"},
		{"non-ascii inside the fence", "---\ntitle: a — dash\n---\nBody cites REQ-1679 here.\n"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			mentions := collectDocumentTicketMentions(testCase.documentText, newCitationFixtureResolver())
			if len(mentions) != 1 {
				t.Fatalf("mentions = %#v, want exactly one", mentions)
			}
			documentRunes := utf16RunesOf(testCase.documentText)
			sliced := string(documentRunes[mentions[0].Offset : mentions[0].Offset+mentions[0].Length])
			if sliced != "REQ-1679" {
				t.Errorf("offset %d slices %q, want %q — the fence length is wrong",
					mentions[0].Offset, sliced, "REQ-1679")
			}
		})
	}
}

// A mention inside the frontmatter fence is NEVER annotated: depends_on, related
// and user_request live there, and a title spliced into one stops the paste
// parsing as YAML. That is the invariant the whole Copy feature rests on.
func TestCollectDocumentTicketMentionsNeverTouchesTheFrontmatterFence(t *testing.T) {
	documentText := "---\nid: REQ-500\ndepends_on: [REQ-1679, REQ-1108]\nuser_request: UR-074\n---\n" +
		"Body cites REQ-1685 only.\n"
	mentions := collectDocumentTicketMentions(documentText, newCitationFixtureResolver())
	fenceLength := strings.Index(documentText, "Body cites")
	for _, mention := range mentions {
		if mention.Offset < fenceLength {
			t.Errorf("mention %s at offset %d sits inside the %d-character frontmatter fence",
				mention.Id, mention.Offset, fenceLength)
		}
	}
	if len(mentions) != 1 || mentions[0].Id != "REQ-1685" {
		t.Errorf("mentions = %#v, want only the body's REQ-1685", mentions)
	}
}

// Compound card ids, the short-segment index, and the refusal to guess. An
// ambiguous segment is NOT a dead reference — the board holds records that match
// — so it earns neither an expansion nor an appendix line, and a caller that
// flags unresolved ids must leave it alone.
func TestTicketMentionResolverNeverGuessesAnAmbiguousSegment(t *testing.T) {
	resolver := newCitationFixtureResolver()
	for _, testCase := range []struct {
		mentionText   string
		wantKind      string
		wantId        string
		wantAmbiguous bool
	}{
		{"REQ-1679", "req", "REQ-1679", false},
		{"UR-074", "ur", "UR-074", false},
		{"UR-001-REQ-042", "req", "UR-001-REQ-042", false},
		{"REQ-077", "req", "UR-003-REQ-077", false},
		{"REQ-042", "", "", true},
		{"REQ-9999", "", "", false},
	} {
		gotKind, gotId := resolver.resolve(testCase.mentionText)
		if gotKind != testCase.wantKind || gotId != testCase.wantId {
			t.Errorf("resolve(%q) = (%q, %q), want (%q, %q)",
				testCase.mentionText, gotKind, gotId, testCase.wantKind, testCase.wantId)
		}
		if gotAmbiguous := resolver.isAmbiguous(testCase.mentionText); gotAmbiguous != testCase.wantAmbiguous {
			t.Errorf("isAmbiguous(%q) = %v, want %v", testCase.mentionText, gotAmbiguous, testCase.wantAmbiguous)
		}
	}
	// An ambiguous mention reaches neither the index nor the appendix.
	mentions := collectDocumentTicketMentions("Compare REQ-042 with REQ-042 again.\n", resolver)
	if len(mentions) != 0 {
		t.Errorf("mentions = %#v, want none — the board refuses to pick between two cards", mentions)
	}
}

func mustMarshalProbeJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, encodeError := json.Marshal(value)
	if encodeError != nil {
		t.Fatalf("encode %v: %v", value, encodeError)
	}
	return string(encoded)
}

// ---- the two resolvers, pinned -------------------------------------------
// Resolution exists twice on purpose. The drawer resolves in the browser
// against a rendered DOM; the clipboard resolves in Go against source bytes,
// because only Go can parse the Markdown that decides which mentions may be
// annotated at all. Collapsing them would mean shipping a Markdown parser to
// the client or resolving the drawer's mentions at build time, and neither is
// this REQ's job.
//
// So they are PINNED rather than merged, the way REQ-248 pins shared geometry:
// one corpus, both implementations, compared both ways. A change to either side
// alone fails here — including a change that makes one side answer MORE, which
// an assertion written only in the drifting side's direction would miss.

// ticketMentionAgreementCorpus is the shared input both sides must answer
// identically about. Each line names something that has been got wrong: a
// boundary, an alternative claiming a run, a compound id, an ambiguous segment.
var ticketMentionAgreementCorpus = []string{
	"Read REQ-1679 lessons and REQ-1679 again.",
	"Compound UR-001-REQ-042 and its short form REQ-042.",
	"A letter suffix REQ-1679a and a longer number REQ-16790.",
	"No boundary in xREQ-1679 or REQ-1679x.",
	"A URL https://example.com/REQ-1679 claims the whole run.",
	"A path do-work/archive/UR-075/REQ-378-title.md claims the whole run.",
	"Punctuation (REQ-1679), REQ-1679. and [REQ-1679].",
	"Not paths: and/or, TLS1.2/1.3, 2.0/5.75.",
	"User requests UR-074 and UR-999.",
	"A short segment REQ-077 and a dead id REQ-9999.",
	"A definition line [REQ-1685]: do-work/queue/REQ-1685-x.md",
}

type ticketMentionAgreementProbeResult struct {
	PatternMatches [][]string        `json:"patternMatches"`
	Resolutions    map[string]string `json:"resolutions"`
}

// describeGoPatternMatches renders one corpus line's matches as
// "<alternative>:<text>", so a drift shows as a readable diff rather than a pair
// of index arrays.
func describeGoPatternMatches(corpusLine string) []string {
	described := []string{}
	for _, matchIndexes := range bodyTicketMentionPattern.FindAllStringSubmatchIndex(corpusLine, -1) {
		for alternative := 1; alternative <= 3; alternative++ {
			start, stop := matchIndexes[alternative*2], matchIndexes[alternative*2+1]
			if start < 0 {
				continue
			}
			described = append(described, fmt.Sprintf("%d:%s", alternative, corpusLine[start:stop]))
		}
	}
	return described
}

// describeGoResolution renders one mention's verdict. "ambiguous" is a verdict
// of its own, not a flavour of "missing": the board holds records that match and
// refuses to pick one, so a caller flagging dead references must skip it.
func describeGoResolution(resolver *ticketMentionResolver, mentionText string) string {
	kind, resolvedId := resolver.resolve(mentionText)
	if kind != "" {
		return kind + ":" + resolvedId
	}
	if resolver.isAmbiguous(mentionText) {
		return "ambiguous"
	}
	return "missing"
}

func TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo(t *testing.T) {
	indexHtml := generateLiveSite(t)

	// Every mention any corpus line contains, so neither side can pass by
	// answering about fewer of them than the other.
	resolutionInputs := []string{}
	seenInputs := map[string]bool{}
	for _, corpusLine := range ticketMentionAgreementCorpus {
		for _, described := range describeGoPatternMatches(corpusLine) {
			if !strings.HasPrefix(described, "3:") || seenInputs[described[2:]] {
				continue
			}
			seenInputs[described[2:]] = true
			resolutionInputs = append(resolutionInputs, described[2:])
		}
	}
	if len(resolutionInputs) < 8 {
		t.Fatalf("the corpus yielded only %d ticket mentions — too few to pin two resolvers", len(resolutionInputs))
	}

	javascriptProbe := `
var requestsById = {};
` + mustMarshalProbeJSON(t, citationFixtureRequestIds) + `.forEach(function (id) { requestsById[id] = { title: id }; });
var userRequestsById = {};
` + mustMarshalProbeJSON(t, citationFixtureUserRequestIds) + `.forEach(function (id) { userRequestsById[id] = { title: id }; });
` + strings.Join([]string{
		sliceDeclarationAfter(t, indexHtml, "var bodyMentionPattern ="),
		sliceBalancedBlockAfter(t, indexHtml, "function buildRequestIdByReqSegment("),
		sliceDeclarationAfter(t, indexHtml, "var requestIdByReqSegment ="),
		sliceBalancedBlockAfter(t, indexHtml, "function resolveTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function isAmbiguousTicketMention("),
	}, "\n") + `

var patternMatches = ` + mustMarshalProbeJSON(t, ticketMentionAgreementCorpus) + `.map(function (corpusLine) {
  var described = [];
  var matchResult;
  bodyMentionPattern.lastIndex = 0;
  while ((matchResult = bodyMentionPattern.exec(corpusLine)) !== null) {
    for (var alternative = 1; alternative <= 3; alternative += 1) {
      if (matchResult[alternative] !== undefined) {
        described.push(alternative + ":" + matchResult[alternative]);
      }
    }
  }
  return described;
});

var resolutions = {};
` + mustMarshalProbeJSON(t, resolutionInputs) + `.forEach(function (mentionText) {
  var ticketTarget = resolveTicketMention(mentionText);
  resolutions[mentionText] = ticketTarget
    ? ticketTarget.kind + ":" + ticketTarget.id
    : (isAmbiguousTicketMention(mentionText) ? "ambiguous" : "missing");
});

process.stdout.write(JSON.stringify({ patternMatches: patternMatches, resolutions: resolutions }));`

	probeOutput := runJavaScriptBehaviorProbe(t, "ticket mention agreement", javascriptProbe)
	var probeResult ticketMentionAgreementProbeResult
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode ticket mention agreement: %v (output %q)", decodeError, probeOutput)
	}

	if len(probeResult.PatternMatches) != len(ticketMentionAgreementCorpus) {
		t.Fatalf("the client answered about %d corpus lines, want %d",
			len(probeResult.PatternMatches), len(ticketMentionAgreementCorpus))
	}
	for corpusIndex, corpusLine := range ticketMentionAgreementCorpus {
		wantMatches := describeGoPatternMatches(corpusLine)
		if !reflect.DeepEqual(probeResult.PatternMatches[corpusIndex], wantMatches) {
			t.Errorf("pattern drift on %q:\n client %v\n Go     %v",
				corpusLine, probeResult.PatternMatches[corpusIndex], wantMatches)
		}
	}

	resolver := newCitationFixtureResolver()
	if len(probeResult.Resolutions) != len(resolutionInputs) {
		t.Errorf("the client resolved %d mentions, Go asked about %d", len(probeResult.Resolutions), len(resolutionInputs))
	}
	for _, mentionText := range resolutionInputs {
		clientVerdict, answered := probeResult.Resolutions[mentionText]
		if !answered {
			t.Errorf("the client left %q unanswered", mentionText)
			continue
		}
		if goVerdict := describeGoResolution(resolver, mentionText); clientVerdict != goVerdict {
			t.Errorf("resolver drift on %q: client %q, Go %q", mentionText, clientVerdict, goVerdict)
		}
	}
}

// ---- the whole tree, not a fixture ---------------------------------------

// Every mention shipped for the real board must slice the id it claims back out
// of the document it was measured on. A fixture can only prove the arithmetic on
// the shapes someone thought to write down; this proves it on all 450-odd real
// documents, including the ones carrying emoji, em dashes and CRLF endings.
//
// The check is deliberately the CLIENT's operation — index the code-unit
// sequence, slice, compare — because an offset that is self-consistent in Go and
// wrong in JavaScript is exactly the failure this class of change makes.
func TestBuildGeneratedBoardMarkdownDataLocatesEveryTicketMention(t *testing.T) {
	markdownData := buildGeneratedBoardMarkdownData(liveBoard(t))

	checkedMentionCount := map[string]int{}
	checkDocument := func(documentKind string, documentId string, documentText string, mentions []generatedTicketMention) {
		documentRunes := utf16RunesOf(documentText)
		previousOffset := -1
		for _, mention := range mentions {
			checkedMentionCount[documentKind]++
			if mention.Offset < 0 || mention.Offset+mention.Length > len(documentRunes) {
				t.Errorf("%s: mention %s at %d+%d falls outside a %d-unit document",
					documentId, mention.Id, mention.Offset, mention.Length, len(documentRunes))
				continue
			}
			// Ascending order is what makes the client's descending splice safe.
			if mention.Offset <= previousOffset {
				t.Errorf("%s: mention %s at %d is not after the previous mention at %d",
					documentId, mention.Id, mention.Offset, previousOffset)
			}
			previousOffset = mention.Offset

			sliced := string(documentRunes[mention.Offset : mention.Offset+mention.Length])
			if !bodyTicketMentionPattern.MatchString(sliced) || len(bodyTicketMentionPattern.FindString(sliced)) != len(sliced) {
				t.Errorf("%s: offset %d slices %q, which is not a whole ticket id", documentId, mention.Offset, sliced)
			}
			// A dead reference reports the words that were WRITTEN — there is no
			// record to report a canonical id from.
			if mention.Kind == "" && sliced != mention.Id {
				t.Errorf("%s: dead reference reports %q but the document says %q", documentId, mention.Id, sliced)
			}
		}
	}

	for requestId, documentText := range markdownData.Requests {
		checkDocument("req", requestId, documentText, markdownData.RequestMentions[requestId])
	}
	for userRequestId, documentText := range markdownData.UserRequests {
		checkDocument("ur", userRequestId, documentText, markdownData.UserRequestMentions[userRequestId])
	}

	// A map that never reached the payload would pass every assertion above
	// without executing one of them — and the two maps are filled by separate
	// loops, so one floor covering both would let either go missing.
	for documentKind, floorCount := range map[string]int{"req": 500, "ur": 100} {
		if checkedMentionCount[documentKind] < floorCount {
			t.Fatalf("only %d %s mentions reached the payload — the live tree carries far more, so that map is not being shipped",
				checkedMentionCount[documentKind], documentKind)
		}
	}
}

// The frontmatter fence of every real document survives the walk untouched:
// nothing before a body's first byte is ever offered as annotatable. Whole-tree
// rather than fixture, because the fences in the tree are not uniform — some
// carry a BOM, some CRLF, one repeats a key, several hold malformed titles.
func TestBuildGeneratedBoardMarkdownDataNeverOffersToAnnotateAFence(t *testing.T) {
	board := liveBoard(t)
	markdownData := buildGeneratedBoardMarkdownData(board)

	checkedFenceCount := 0
	for _, ticket := range board.AllRequests {
		if ticket.FrontmatterMarkdown == "" {
			continue
		}
		checkedFenceCount++
		fenceUnitLength := utf16LengthOf(ticket.FrontmatterMarkdown)
		for _, mention := range markdownData.RequestMentions[ticket.RequestId] {
			if mention.Offset < fenceUnitLength {
				t.Errorf("%s: mention %s at %d sits inside its %d-unit frontmatter fence",
					ticket.RequestId, mention.Id, mention.Offset, fenceUnitLength)
			}
		}
	}
	if checkedFenceCount < 100 {
		t.Fatalf("only %d fenced REQ files were checked — the assertion above never ran on a real fence", checkedFenceCount)
	}
}
