package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unicode/utf16"
)

// citationFixtureRequestIds and citationFixtureUserRequestIds are the board every
// test below resolves against. UR-001-REQ-042 and UR-002-REQ-042 share a REQ
// segment on purpose: "REQ-042" is the ambiguous case the resolver must refuse
// to guess at, and it is only ambiguous while BOTH are present.
var citationFixtureRequestIds = []string{
	"REQ-1679", "REQ-1108", "REQ-1685", "REQ-500", "REQ-501",
	"UR-001-REQ-042", "UR-002-REQ-042", "UR-003-REQ-077",
}

var citationFixtureUserRequestIds = []string{"UR-003", "UR-074", "UR-075"}

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
// string is, so a test can index it the way the client will.
//
// It uses unicode/utf16 from the standard library ON PURPOSE. Hand-rolling the
// same "> 0xFFFF costs two" loop that utf16LengthOf uses would check the
// production algorithm against a copy of itself, and every offset assertion in
// this file would pass for a shared misconception. The stdlib encoder is the
// independent second opinion; utf8.RuneCountInString is not, because it counts
// the wrong unit.
func utf16RunesOf(sourceText string) []rune {
	encodedUnits := utf16.Encode([]rune(sourceText))
	codeUnits := make([]rune, 0, len(encodedUnits))
	for _, encodedUnit := range encodedUnits {
		codeUnits = append(codeUnits, rune(encodedUnit))
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
		"Dead REQ-9999 in prose, ambiguous REQ-042, and short REQ-077.\n" +
		"\n" +
		"A dead id in a code span: `do-work run REQ-8886` still earns its line.\n" +
		"\n" +
		"Blocked on UR-003-REQ-077 written out, then REQ-077 again.\n"

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
		// A dead id in a CODE SPAN is still reported, and that is the only
		// behaviour separating surfaceCodeSpan from surfaceCodeBlock — where a
		// dead id is an illustration and reported by neither. Without this row,
		// collapsing the two surfaces passes the whole suite.
		`missing REQ-8886 quoted "REQ-8886"`,
		// The compound form written out, then its short form. Both name the same
		// record, so the SECOND must not expand — which pins that first-mention
		// memory is keyed by the resolved id rather than by the written text.
		// Every other repeat in this fixture repeats the same characters, so
		// keying on mentionText passes without these two.
		`req UR-003-REQ-077 quoted "UR-003-REQ-077"`,
		`req UR-003-REQ-077 quoted "REQ-077"`,
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
	"Fixed in _REQ-1679_ last week.",
	"_tracked under UR-003-REQ-077_",
	"Adjacent REQ-1679_UR-003 and _UR-003_.",
	"Suffixes _REQ-1679a_ and _UR-003-REQ-077a_.",
	"Invalid UR-003-REQ-077ab, UR-003-REQ-077A, and xUR-003-REQ-077.",
	"Compound UR-001-REQ-042 and its short form REQ-042.",
	"A letter suffix REQ-1679a and a longer number REQ-16790.",
	"No boundary in xREQ-1679 or REQ-1679x.",
	"A URL https://example.com/REQ-1679 claims the whole run.",
	"A path do-work/archive/UR-075/REQ-378-title.md claims the whole run.",
	"Scoped node_modules/@scope/pkg/index.js and retina assets/@2x/sprite.png.",
	"Digits in extensions: docs/media/recording.mp4 and data/report.h5.",
	"Every punctuation at once: a_b-c.d/e@f/g_h-i.j.txt",
	"Near misses: trailing/dot. and digit/start.9ext stay prose.",
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

// REQ-385: prove both the emitted offsets and the drawer's real consumption
// loop. A pattern-only probe misses the drawer retrying inside a rejected
// compound and linking its REQ segment. The clipboard uses the real Go index.
func TestJavaScriptBehaviorTicketMentionUnderscoreBoundaries(t *testing.T) {
	indexHtml := generateLiveSite(t)
	testCases := []struct {
		Source            string                   `json:"source"`
		DocumentText      string                   `json:"documentText"`
		InsideFencedBlock bool                     `json:"insideFencedBlock"`
		Want              []string                 `json:"-"`
		WantBody          string                   `json:"-"`
		Mentions          []generatedTicketMention `json:"mentions"`
	}{
		{Source: "Fixed in _REQ-1679_ last week.", Want: []string{`req REQ-1679 EXPAND "REQ-1679"`}, WantBody: "Fixed in _REQ-1679 (Fix mentions)_ last week."},
		{Source: "_tracked under UR-003-REQ-077_", Want: []string{`req UR-003-REQ-077 EXPAND "UR-003-REQ-077"`}, WantBody: "_tracked under UR-003-REQ-077 (Compound work)_"},
		{Source: "😀 — REQ-1679_UR-003", Want: []string{`req REQ-1679 EXPAND "REQ-1679"`, `ur UR-003 EXPAND "UR-003"`}, WantBody: "😀 — REQ-1679 (Fix mentions)_UR-003 (Ship the widget)"},
		{Source: "_UR-003_", Want: []string{`ur UR-003 EXPAND "UR-003"`}, WantBody: "_UR-003 (Ship the widget)_"},
		{Source: "_REQ-1679a_ _UR-003-REQ-077a_", Want: []string{`missing REQ-1679a quoted "REQ-1679a"`, `missing UR-003-REQ-077a quoted "UR-003-REQ-077a"`}},
		{Source: "UR-003-REQ-077ab UR-003-REQ-077A xUR-003-REQ-077 1UR-003-REQ-077", Want: []string{}},
		{Source: "xREQ-1679 REQ-1679ab UR-003x", Want: []string{}},
		// Suppressing a missing compound in a code block must not retry its
		// resolvable inner REQ segment (REQ-077 names UR-003-REQ-077 here).
		{Source: "_UR-999-REQ-077_", InsideFencedBlock: true, Want: []string{}},
		{Source: "UR-003-REQ-077ab then REQ-1679.", Want: []string{`req REQ-1679 EXPAND "REQ-1679"`}, WantBody: "UR-003-REQ-077ab then REQ-1679 (Fix mentions)."},
	}
	for index := range testCases {
		testCase := &testCases[index]
		testCase.DocumentText = testCase.Source
		if testCase.InsideFencedBlock {
			testCase.DocumentText = "```\n" + testCase.Source + "\n```\n"
		}
		testCase.Mentions = collectDocumentTicketMentions(testCase.DocumentText, newCitationFixtureResolver())
		if got := describeTicketMentions(testCase.DocumentText, testCase.Mentions); !reflect.DeepEqual(got, testCase.Want) {
			t.Errorf("Go mentions for %q:\n got %v\nwant %v", testCase.Source, got, testCase.Want)
		}
		if testCase.WantBody == "" {
			testCase.WantBody = testCase.DocumentText
		}
	}

	javascriptProbe := `
var requestsById = {
  "REQ-1679": { title: "Fix mentions", status: "completed" },
  "UR-003-REQ-077": { title: "Compound work", status: "pending" }
};
var userRequestsById = { "UR-003": { title: "Ship the widget" } };
var document = {
  createDocumentFragment: function () { return { appendChild: function () {} }; },
  createTextNode: function (text) { return text; }
};
var linkedMentions;
function makeTicketLink(kind, id, text, expand) {
  linkedMentions.push(kind + " " + id + (expand ? " EXPAND " : " quoted ") + JSON.stringify(text));
  return {};
}
function makeMissingTicketMention(text) {
  linkedMentions.push("missing " + text + " quoted " + JSON.stringify(text));
  return {};
}
` + strings.Join([]string{
		sliceDeclarationAfter(t, indexHtml, "var bodyMentionPattern ="),
		sliceDeclarationAfter(t, indexHtml, "var inlineTicketTitleMaxLength ="),
		sliceBalancedBlockAfter(t, indexHtml, "function buildRequestIdByReqSegment("),
		sliceDeclarationAfter(t, indexHtml, "var requestIdByReqSegment ="),
		sliceBalancedBlockAfter(t, indexHtml, "function resolveTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function isAmbiguousTicketMention("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeRequestStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildLinkifiedFragment("),
		sliceBalancedBlockAfter(t, indexHtml, "function recordReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateTicketMentions("),
	}, "\n") + `
var results = ` + mustMarshalProbeJSON(t, testCases) + `.map(function (testCase) {
  linkedMentions = [];
  var renderState = { expandedTicketKeys: {}, glossaryKeys: {}, glossaryEntries: [] };
  buildLinkifiedFragment(testCase.source, testCase.insideFencedBlock, testCase.insideFencedBlock, renderState);
  var annotation = annotateTicketMentions(testCase.documentText, testCase.mentions);
  return { mentions: linkedMentions, body: annotation.text,
    glossary: renderState.glossaryEntries.map(function (entry) { return entry.id; }),
    appendix: annotation.referencedTickets.map(function (entry) { return entry.id; }) };
});
process.stdout.write(JSON.stringify(results));`
	probeOutput := runJavaScriptBehaviorProbe(t, "underscore ticket boundaries", javascriptProbe)
	var results []struct {
		Mentions []string `json:"mentions"`
		Body     string   `json:"body"`
		Glossary []string `json:"glossary"`
		Appendix []string `json:"appendix"`
	}
	if decodeError := json.Unmarshal(probeOutput, &results); decodeError != nil {
		t.Fatalf("decode underscore ticket boundaries: %v (output %q)", decodeError, probeOutput)
	}
	if len(results) != len(testCases) {
		t.Fatalf("client answered %d cases, want %d", len(results), len(testCases))
	}
	for index, result := range results {
		testCase := testCases[index]
		if !reflect.DeepEqual(result.Mentions, testCase.Want) {
			t.Errorf("drawer mentions for %q:\n got %v\nwant %v", testCase.Source, result.Mentions, testCase.Want)
		}
		if result.Body != testCase.WantBody {
			t.Errorf("clipboard body for %q:\n got %q\nwant %q", testCase.Source, result.Body, testCase.WantBody)
		}
		wantGlossary, wantAppendix := []string{}, []string{}
		for _, described := range testCase.Want {
			fields := strings.Fields(described)
			wantAppendix = append(wantAppendix, fields[1])
			if fields[0] != "missing" {
				wantGlossary = append(wantGlossary, fields[1])
			}
		}
		if !reflect.DeepEqual(result.Glossary, wantGlossary) || !reflect.DeepEqual(result.Appendix, wantAppendix) {
			t.Errorf("references for %q: drawer %v want %v; clipboard %v want %v", testCase.Source,
				result.Glossary, wantGlossary, result.Appendix, wantAppendix)
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

// The file-path syntax has exactly ONE definition in Go. citations.go composes
// repoFileMentionPattern instead of restating it, so a change to the file
// scanner reaches the mention scanner automatically — the lock-step obligation
// that failed twice in this file's history is gone rather than documented.
//
// The property under test is AGREEMENT ON VALUE, not the presence of a
// composition expression (REQ-289: grep the value, not the constant name).
// Re-inlining a literal that happens to be byte-identical is harmless and must
// pass; re-inlining one that differs by a single character is the drift this
// exists to catch, so the corpus below has to discriminate at every character
// class the pattern uses. It did not on the first attempt: three mutations
// survived because no fixture carried an "@" in a directory segment or a digit
// in an extension.
var filePathAgreementCorpus = []string{
	"do-work/archive/UR-075/REQ-378-title.md",
	"skills/do-work-board/tools/queue-kanban/citations.go",
	"node_modules/@scope/pkg/index.js",   // "@" in a directory segment
	"assets/@2x/sprite-sheet.png",        // "@" beside a digit
	"docs/media/recording.mp4",           // a digit inside the extension
	"data/exports/report.h5",             // a digit ENDING the extension
	"do-work/queue/REQ-380-cross-ref.md", // a digit inside a path segment
	"a_b-c.d/e@f/g_h-i.j.txt",            // every allowed punctuation at once
	"and/or",                             // no extension — not a path
	"TLS1.2/1.3",                         // a version ratio — not a path
	"2.0/5.75",                           // a numeric ratio — not a path
	"paths/<placeholder>/file.md",        // an angle-bracket placeholder
	"trailing/dot.",                      // a dot with no extension after it
	"digit/start.9ext",                   // an extension starting with a digit
}

func TestBodyTicketMentionPatternComposesTheOneFilePathDefinition(t *testing.T) {
	// Submatch 2 is the composed pattern's file-path alternative. Both scanners
	// answer about the same subject, and every answer must match.
	filePathAlternativeMatch := func(subject string) string {
		matchIndexes := bodyTicketMentionPattern.FindStringSubmatchIndex(subject)
		if matchIndexes == nil || matchIndexes[4] < 0 {
			return ""
		}
		return subject[matchIndexes[4]:matchIndexes[5]]
	}

	discriminatingCount := 0
	for _, subject := range filePathAgreementCorpus {
		fileScannerMatch := repoFileMentionPattern.FindString(subject)
		if fileScannerMatch != "" {
			discriminatingCount++
		}
		if mentionScannerMatch := filePathAlternativeMatch(subject); mentionScannerMatch != fileScannerMatch {
			t.Errorf("on %q the mention scanner reads the path as %q and the file scanner as %q — the two Go definitions have drifted",
				subject, mentionScannerMatch, fileScannerMatch)
		}
	}

	// A corpus of nothing but non-paths would agree trivially, and every
	// assertion above would pass without exercising the pattern at all.
	if discriminatingCount < 8 {
		t.Fatalf("only %d corpus lines matched as paths — the agreement above is nearly vacuous", discriminatingCount)
	}

	// Composition, not coincidence: the whole of repoFileMentionPattern appears
	// in the alternation exactly once, as its own group.
	composedPattern := bodyTicketMentionPattern.String()
	if !strings.Contains(composedPattern, "|("+repoFileMentionPattern.String()+")|") {
		t.Errorf("bodyTicketMentionPattern does not carry repoFileMentionPattern as its second alternative:\n got %s", composedPattern)
	}
}

// A ticket id written as a link — or as a link's label answered by a definition
// elsewhere — is left exactly as the author wrote it, on BOTH surfaces.
//
// Two failures, one cause. Splicing a title into `[REQ-1679]` when a
// `[REQ-1679]: …` definition exists elsewhere orphans the use from the
// definition, and the paste silently loses a link it had — REQ-379's F6 from
// the other side, since the definition itself is already protected. And the
// drawer skips anchor text outright (`parentElement.closest("a")`), so
// annotating it here would make the paste and the drawer disagree about the
// same body. Ticket ids in link syntax are REQ-382's subject; until then,
// neither surface touches them.
func TestCollectDocumentTicketMentionsLeavesLinkSyntaxAlone(t *testing.T) {
	documentText := "---\nid: REQ-500\n---\n\n" +
		// The em dash is load-bearing: describeTicketMentions slices every
		// reported offset back out of the document, so a body that is pure ASCII
		// cannot tell a byte offset from a UTF-16 one. This fixture puts a
		// multi-byte character BEFORE the mentions and before the definition
		// block, which is the arrangement that makes the two diverge.
		"A shortcut use [REQ-1679], a collapsed use [REQ-1108][], an inline link\n" +
		"[REQ-1685](https://example.test/x) — and an image ![REQ-501 diagram](y.png).\n\n" +
		"Prose afterwards cites REQ-1685 and REQ-1679 for real.\n\n" +
		"[REQ-1679]: https://example.test/a\n" +
		"[REQ-1108]: https://example.test/b\n"

	mentions := collectDocumentTicketMentions(documentText, newCitationFixtureResolver())
	gotMentions := describeTicketMentions(documentText, mentions)
	wantMentions := []string{
		// Only the two prose mentions survive, and each is the FIRST prose sighting
		// of its id — a link label must not spend the one expansion an id gets.
		`req REQ-1685 EXPAND "REQ-1685"`,
		`req REQ-1679 EXPAND "REQ-1679"`,
	}
	if !reflect.DeepEqual(gotMentions, wantMentions) {
		t.Errorf("emitted mentions:\n got  %s\n want %s",
			strings.Join(gotMentions, "\n       "), strings.Join(wantMentions, "\n       "))
	}
}

// A link anywhere up the chain wins, even over a code span nested inside it.
// The drawer's rule is `parentElement.closest("a")`, which returns early for
// every text node under an anchor without looking at what else encloses it — so
// a backticked id inside a link label earns no drawer glossary line and must
// earn no clipboard appendix line. Taking the innermost construct instead reads
// as the more careful claim and is simply a DIFFERENT answer from the drawer's,
// which is the one thing these two may never give.
func TestTextNodeSurfaceLetsALinkWinOverAnEnclosedCodeSpan(t *testing.T) {
	documentText := "Prose cites REQ-1679 and a link [`REQ-1108` in code](https://example.test/x).\n"
	mentions := collectDocumentTicketMentions(documentText, newCitationFixtureResolver())
	gotMentions := describeTicketMentions(documentText, mentions)
	wantMentions := []string{`req REQ-1679 EXPAND "REQ-1679"`}
	if !reflect.DeepEqual(gotMentions, wantMentions) {
		t.Errorf("emitted mentions:\n got  %s\n want %s",
			strings.Join(gotMentions, "\n       "), strings.Join(wantMentions, "\n       "))
	}
}

// utf16OffsetCursor answers a scripted sequence of offsets, forwards and — the
// defensive path — backwards. A cursor that mis-restores its base on a rewind
// returns confident wrong numbers rather than failing, so both directions are
// driven against the independent stdlib measurement.
func TestUtf16OffsetCursorAnswersBothDirections(t *testing.T) {
	sourceText := "Emoji 😀 dash — REQ-1679 then 😀 again REQ-1108 end."
	const baseUnitOffset = 7
	expectedUnitOffsetAt := func(byteOffset int) int {
		return baseUnitOffset + len(utf16.Encode([]rune(sourceText[:byteOffset])))
	}

	firstMentionStart := strings.Index(sourceText, "REQ-1679")
	secondMentionStart := strings.Index(sourceText, "REQ-1108")
	cursor := newUtf16OffsetCursor(sourceText, baseUnitOffset)
	for _, byteOffset := range []int{
		0, firstMentionStart, firstMentionStart + len("REQ-1679"),
		secondMentionStart, secondMentionStart + len("REQ-1108"), len(sourceText),
		// Backwards, twice, and then forwards again: the rewind must leave the
		// cursor usable, not just answer once.
		firstMentionStart, 0, secondMentionStart, len(sourceText),
	} {
		if gotUnitOffset := cursor.at(byteOffset); gotUnitOffset != expectedUnitOffsetAt(byteOffset) {
			t.Errorf("cursor.at(%d) = %d, want %d", byteOffset, gotUnitOffset, expectedUnitOffsetAt(byteOffset))
		}
	}
}

// The mention walk reads the cursor exactly twice per mention, ascending.
// Building the offsets inline in the composite literal called at(start) again
// after at(stop), which takes the rewind branch and rescans the whole prefix —
// the quadratic the cursor exists to avoid, paid twice per mention on the
// longest bodies in the tree.
//
// Asserted on the source text because the cost is invisible in the OUTPUT: the
// offsets are identical either way, so no behavioural assertion can see it and
// a timing assertion would be flaky. This is the same greppable-value technique
// TestClipboardCarriesNoMarkdownScanner uses on the client.
func TestCollectDocumentTicketMentionsReadsTheOffsetCursorTwicePerMention(t *testing.T) {
	citationsSource, readError := os.ReadFile("citations.go")
	if readError != nil {
		t.Fatalf("reading citations.go: %v", readError)
	}
	functionStart := strings.Index(string(citationsSource), "func collectDocumentTicketMentions(")
	if functionStart < 0 {
		t.Fatal("collectDocumentTicketMentions not found — the assertion below would be vacuous")
	}
	functionEnd := strings.Index(string(citationsSource)[functionStart:], "\n}\n")
	if functionEnd < 0 {
		t.Fatal("could not find the end of collectDocumentTicketMentions")
	}
	functionBody := string(citationsSource)[functionStart : functionStart+functionEnd]

	if cursorReadCount := strings.Count(functionBody, "documentOffsets.at("); cursorReadCount != 2 {
		t.Errorf("collectDocumentTicketMentions reads the offset cursor %d times, want 2 (one for the mention start, one for its end)",
			cursorReadCount)
	}
}

// The one skew the split leaves, and the client's guard for it.
//
// D3 co-located the offsets with the document text so those two can never come
// from different builds. TITLES were not co-located and cannot be: they live in
// board-data.js, which the page loaded once, while /board-markdown.js re-walks
// the tree on the Copy click. So in serve mode a REQ created since page load
// resolves against the fresh tree and has no title in the stale snapshot.
//
// The client used to splice the empty string, writing a bare " ()" into the
// paste and an appendix line reading "- REQ-384 —  (not in tree)". Now it
// leaves the mention as the author wrote it and the appendix says which side is
// stale. Saying "not found in this queue" would be false — the queue has it.
func TestJavaScriptBehaviorClipboardSurvivesATitleTheSnapshotDoesNotHaveYet(t *testing.T) {
	indexHtml := generateLiveSite(t)

	// Resolved by the build, absent from the page: exactly what a REQ created
	// between page load and the Copy click looks like to the client.
	documentText := "---\nid: REQ-500\n---\n\nBlocked on REQ-1679 until it lands.\n"
	// The index is computed by the REAL walk against a resolver that knows
	// REQ-1679; the page stub below does not. That IS the skew — a fresh
	// board-markdown.js meeting a stale board-data.js — rather than a
	// hand-written stand-in for it.
	staleSnapshotMentions, encodeError := json.Marshal(
		collectDocumentTicketMentions(documentText, newCitationFixtureResolver()))
	if encodeError != nil {
		t.Fatalf("encode probe ticket mentions: %v", encodeError)
	}
	if !strings.Contains(string(staleSnapshotMentions), `"id":"REQ-1679","expand":true`) {
		t.Fatalf("the walk did not offer REQ-1679 for expansion, so the guard below is never reached: %s", staleSnapshotMentions)
	}

	javascriptProbe := `
var requestsById = { "REQ-500": { title: "Host document", status: "claimed" } };
var userRequestsById = {};
` + strings.Join([]string{
		sliceDeclarationAfter(t, indexHtml, "var inlineTicketTitleMaxLength ="),
		sliceBalancedBlockAfter(t, indexHtml, "function describeRequestStatus("),
		sliceBalancedBlockAfter(t, indexHtml, "function ticketTitleFor("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeTicketTitle("),
		sliceBalancedBlockAfter(t, indexHtml, "function shortTicketTitle("),
		sliceDeclarationAfter(t, indexHtml, "var referencedTicketsGlossaryHeading ="),
		sliceBalancedBlockAfter(t, indexHtml, "function recordReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateTicketMentions("),
		sliceBalancedBlockAfter(t, indexHtml, "function describeReferencedTicket("),
		sliceBalancedBlockAfter(t, indexHtml, "function buildReferencedTicketsGlossary("),
		sliceBalancedBlockAfter(t, indexHtml, "function annotateClipboardPayload("),
	}, "\n") + `
process.stdout.write(JSON.stringify({ payload: annotateClipboardPayload(
  [{ text: ` + mustMarshalJSONString(t, documentText) + `, ticketMentions: ` + string(staleSnapshotMentions) + ` }], ["REQ-500"]
) }));`

	probeOutput := runJavaScriptBehaviorProbe(t, "stale-snapshot title", javascriptProbe)
	var probeResult struct {
		Payload string `json:"payload"`
	}
	if decodeError := json.Unmarshal(probeOutput, &probeResult); decodeError != nil {
		t.Fatalf("decode stale-snapshot behavior: %v (output %q)", decodeError, probeOutput)
	}

	if strings.Contains(probeResult.Payload, "REQ-1679 ()") {
		t.Errorf("an empty title was spliced as \"()\":\n%s", probeResult.Payload)
	}
	if !strings.Contains(probeResult.Payload, "Blocked on REQ-1679 until it lands.\n") {
		t.Errorf("the body line was not left as the author wrote it:\n%s", probeResult.Payload)
	}
	// The appendix is where a reader learns why the title is missing, and it must
	// not claim the id is unknown to the queue — the build just resolved it.
	if !strings.Contains(probeResult.Payload, "- REQ-1679 — added since this board was loaded — reload to see its title\n") {
		t.Errorf("the appendix does not say which side is stale:\n%s", probeResult.Payload)
	}
	if strings.Contains(probeResult.Payload, "not found in this queue") {
		t.Errorf("the appendix claims the queue has no such record, but the build resolved it:\n%s", probeResult.Payload)
	}
}

// utf16LengthOf is measured against what actually reaches the browser: the
// string after encoding/json has written it and a JSON reader has read it back.
// That round trip is not the identity — encoding/json replaces every invalid
// UTF-8 byte and every lone surrogate with U+FFFD — so a length computed on the
// Go-side bytes can differ from the length of the JavaScript string the client
// indexes, and every offset after it would be wrong.
//
// The oracle is unicode/utf16 rather than a second copy of utf16LengthOf's own
// loop: a helper that mirrors the implementation agrees with it about a shared
// misconception, which is how a whole suite of offset assertions can pass while
// every offset is wrong.
func TestUtf16LengthMatchesWhatTheClientReceives(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		sourceText string
	}{
		{"ascii", "REQ-1679 plain"},
		{"em dash", "REQ-1679 — the character every REQ body carries"},
		{"astral emoji", "REQ-1679 😀 beyond the BMP"},
		{"combining marks", "REQ-1679 été with combining acutes"},
		{"line and paragraph separators", "REQ-1679 \u2028 and \u2029"},
		{"byte order mark", "\ufeffREQ-1679 after a BOM"},
		{"invalid utf-8", "REQ-1679 \xff\xfe raw bytes"},
		// A GENUINE U+FFFD, three bytes wide, which invalid input decodes to as
		// well. Without it a length that counted a replacement character by its
		// byte width passed, because every invalid byte above is one byte wide.
		{"literal replacement character", "REQ-1679 \ufffd pasted from somewhere"},
		{"lone surrogate in wtf-8", "REQ-1679 \xed\xa0\x80 lone surrogate"},
		{"nul and control bytes", "REQ-1679 \x00\x01\x1f controls"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			encoded, encodeError := json.Marshal(testCase.sourceText)
			if encodeError != nil {
				t.Fatalf("encode: %v", encodeError)
			}
			var receivedText string
			if decodeError := json.Unmarshal(encoded, &receivedText); decodeError != nil {
				t.Fatalf("decode: %v", decodeError)
			}
			wantUnitCount := len(utf16.Encode([]rune(receivedText)))
			if gotUnitCount := utf16LengthOf(testCase.sourceText); gotUnitCount != wantUnitCount {
				t.Errorf("utf16LengthOf(%q) = %d, but the client receives a %d-unit string",
					testCase.sourceText, gotUnitCount, wantUnitCount)
			}
		})
	}
}

// The index has to survive the trip to disk, not just the builder.
//
// The end-to-end pin already exists — TestBrowserBehaviorBoardColumnCopyAll
// clicks a real Copy button on a real generated board and asserts a literal
// annotated payload, spliced titles and glossary included. But that lane SKIPS
// when no browser is available (browser_probe_test.go), so on a machine without
// one the whole seam between buildGeneratedBoardMarkdownData and the client
// goes unchecked: dropping the two maps after the build passes everything else.
//
// This is the cheap guard for that case. It reads the file the page actually
// loads, so it fails if the index is built and then not shipped.
func TestGeneratedBoardMarkdownFileShipsTheTicketMentionIndex(t *testing.T) {
	outputDirectory := generateLiveSiteInDir(t)
	markdownJs, readError := os.ReadFile(filepath.Join(outputDirectory, "board-markdown.js"))
	if readError != nil {
		t.Fatalf("reading generated board-markdown.js: %v", readError)
	}

	const assignmentPrefix = "window.queueKanbanBoardMarkdownData = "
	jsonText, foundPrefix := strings.CutPrefix(string(markdownJs), assignmentPrefix)
	if !foundPrefix {
		t.Fatalf("board-markdown.js does not open with %q", assignmentPrefix)
	}
	var shippedPayload struct {
		Requests        map[string]string                   `json:"requests"`
		RequestMentions map[string][]generatedTicketMention `json:"requestMentions"`
	}
	if decodeError := json.Unmarshal([]byte(strings.TrimSuffix(strings.TrimSpace(jsonText), ";")), &shippedPayload); decodeError != nil {
		t.Fatalf("decoding the shipped payload: %v", decodeError)
	}

	if len(shippedPayload.RequestMentions) == 0 {
		t.Fatal("board-markdown.js ships no requestMentions at all — the client would paste every document unannotated")
	}
	// Not just present: pointing at the right characters in the document shipped
	// beside it. A map that survived the build but lost its pairing is the
	// failure that corrupts a paste rather than merely flattening it.
	checkedDocumentCount := 0
	for requestId, mentions := range shippedPayload.RequestMentions {
		documentText, hasDocument := shippedPayload.Requests[requestId]
		if !hasDocument {
			t.Errorf("%s has mentions but no document in the same payload", requestId)
			continue
		}
		documentRunes := utf16RunesOf(documentText)
		previousOffset := -1
		for _, mention := range mentions {
			// Ascending order is what makes the client's descending splice safe,
			// and it is a property of the SHIPPED list — a payload that lost the
			// ordering somewhere between the builder and the file splices titles
			// at offsets that have already moved.
			if mention.Offset <= previousOffset {
				t.Errorf("%s: shipped mention %s at %d is not after the previous one at %d",
					requestId, mention.Id, mention.Offset, previousOffset)
			}
			previousOffset = mention.Offset
			if mention.Offset+mention.Length > len(documentRunes) {
				t.Errorf("%s: mention %s at %d+%d overruns its %d-unit document",
					requestId, mention.Id, mention.Offset, mention.Length, len(documentRunes))
				continue
			}
			sliced := string(documentRunes[mention.Offset : mention.Offset+mention.Length])
			if bodyTicketMentionPattern.FindString(sliced) != sliced {
				t.Errorf("%s: offset %d slices %q, which is not a whole ticket id", requestId, mention.Offset, sliced)
			}
			checkedDocumentCount++
		}
	}
	if checkedDocumentCount < 500 {
		t.Fatalf("only %d shipped mentions were checked — the live tree carries thousands", checkedDocumentCount)
	}
}
