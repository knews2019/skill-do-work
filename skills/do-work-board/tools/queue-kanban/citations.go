package main

import (
	"regexp"
	"sort"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ---- where a ticket id may be annotated -----------------------------------
//
// The board's Copy button pastes a REQ/UR file with the first mention of each
// referenced id carrying that ticket's title. Deciding WHICH mentions may carry
// one is a Markdown question: an id inside a fenced block, an indented block, a
// code span or a link reference definition is quoted text and must keep its
// exact bytes, and an id in prose is a real reference.
//
// The browser client cannot answer that question — it has raw bytes and no
// parser — and the hand-rolled scanner it used to carry disagreed with
// CommonMark on six separate constructs. Go can answer it exactly, because the
// parser is already here: render.go builds a goldmark renderer for the drawer's
// HTML, and this file reuses THAT PARSER on the raw stored body.
//
// The same parser, deliberately not the same parse. renderMarkdownBodyToHtml
// preprocesses its input first (insertQuestionOptionHardBreaks, which rewrites
// 92 of this repo's 376 request bodies), and every byte it inserts shifts the
// offsets after it — so sharing one AST would ship positions measured against
// text the client never receives. See collectDocumentTicketMentions.
//
// filementions.go is the same move for the same reason: the client cannot stat
// the filesystem, so Go ships repoFileMentions.
//
// The division that follows from that, and the one to keep: MARKDOWN KNOWLEDGE
// LIVES IN GO, THE CLIENT SPLICES. Nothing downstream of this file re-derives
// block structure.

// bodyTicketMentionPattern finds every mention a body can carry, with the
// alternation order acting as priority: an absolute URL, then a repo-relative
// file path, then a REQ/UR id (compound form first, so "UR-002-REQ-031" is one
// mention rather than two). The first two alternatives exist only to CLAIM
// those runs: an id inside a URL or inside a path names neither, and expanding
// one would rewrite a link or a path into something that no longer resolves.
//
// The file-path alternative is COMPOSED from repoFileMentionPattern rather than
// restated. Two Go regexps that must describe the same syntax are two things to
// keep in step, and that obligation has failed twice in this file's history —
// so there is one definition of a repo-relative path in Go, and the only reason
// this one is wrapped is the capture group the alternation needs.
//
// What is left to keep in step is the WIRE: bodyMentionPattern in
// web/board-detail.js scans the drawer's rendered text with the same three
// alternatives in the same order, and a browser cannot import a Go variable.
// TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo drives both
// over one shared corpus in both directions, so a drift fails rather than
// silently changing what a paste says.
var bodyTicketMentionPattern = regexp.MustCompile(
	`(https?://[^\s<>"')\]]+)` +
		`|(` + repoFileMentionPattern.String() + `)` +
		`|(?P<ticket>UR-\d+-REQ-\d+[a-z]?|REQ-\d+[a-z]?|UR-\d+)`)

// Boundaries are checked AFTER consuming a candidate. Anchoring the regexp
// lets a rejected compound backtrack to its UR prefix. Only ASCII letters and
// digits continue an id; underscore is punctuation, including Markdown emphasis.
func isTicketMentionBoundary(sourceText string, offset int) bool {
	if offset < 0 || offset >= len(sourceText) {
		return true
	}
	currentByte := sourceText[offset]
	return !(currentByte >= 'A' && currentByte <= 'Z' ||
		currentByte >= 'a' && currentByte <= 'z' ||
		currentByte >= '0' && currentByte <= '9')
}

// ticketMentionGroupIndex locates the ticket-id alternative by NAME. Reading it
// by fixed position was safe only while this pattern owned all three groups —
// and it stopped owning them the moment it began composing repoFileMentionPattern
// from filementions.go. One capturing group added there would silently shift
// every later index, and the build would ship offsets spanning a directory
// segment of a file path. A name cannot shift.
var ticketMentionGroupIndex = bodyTicketMentionPattern.SubexpIndex("ticket")

// requestIdSegmentPattern pulls the REQ segment out of a compound card id
// ("UR-002-REQ-031" → "REQ-031"). Mirrors the literal in
// buildRequestIdByReqSegment (web/board-core.js).
var requestIdSegmentPattern = regexp.MustCompile(`(?i)REQ-\d+[a-z]?`)

// mentionSurface classifies the text a mention sits in. The three values drive
// two INDEPENDENT suppressions, which is why one boolean will not do:
//
//	surfaceProse     — a real reference: title expands, a dead id is reported.
//	surfaceCodeSpan  — a code run must not be contaminated with prose, so no
//	                   title; an inline `REQ-005` is still a reference, so a
//	                   dead id is still reported.
//	surfaceCodeBlock — fenced or indented. REQ bodies print templates and
//	                   worked examples here, so an id that answers to nothing is
//	                   an illustration: no title AND no dead-id report.
//	surfaceLinkLabel — a link's own text or an image's alt text. Nothing at all:
//	                   no title, no dead-id report, no appendix line.
//
// surfaceLinkLabel exists for two reasons that happen to have one answer. A
// title spliced into a SHORTCUT reference use — `[REQ-001]` answered by a
// `[REQ-001]: …` definition elsewhere — orphans it from its definition, so the
// paste silently loses a link it had; protecting the definition (which
// goldmark keeps no text for) and not the use fixes half the bug. And the
// drawer already skips anchor text outright (`parentElement.closest("a")` in
// linkifyDetailBody), so annotating it here would make the paste and the
// drawer say different things about the same body. Ticket ids written as
// Markdown links are REQ-382's subject; until it lands, both surfaces leave
// them alone.
//
// A mention covered by none of these (a link reference DEFINITION, a raw HTML
// block, a fence's own backticks) is not annotatable either and is dropped:
// goldmark keeps no prose text there, so there is nothing to expand.
type mentionSurface int

const (
	surfaceProse mentionSurface = iota
	surfaceCodeSpan
	surfaceCodeBlock
	surfaceLinkLabel
)

// surfaceRange is one half-open byte range of a body, classified. Ranges never
// overlap: they come from disjoint AST segments.
type surfaceRange struct {
	start   int
	stop    int
	surface mentionSurface
}

// collectMentionSurfaces walks a parsed body once and returns every byte range
// whose surface is known, sorted by start offset.
//
// Two node families, two extraction paths, because goldmark models them
// differently — and the difference is not cosmetic: Lines() PANICS on an inline
// node (goldmark/ast/inline.go:38), so a CodeSpan's extent has to come from its
// child *ast.Text segments instead.
//
//   - Block nodes carry line segments. A fenced or indented code block's lines
//     are the quoted text; a fenced block's Info is the only part of its opening
//     fence line that can hold anything but fence characters.
//   - Inline *ast.Text nodes are every byte goldmark keeps as prose, wherever it
//     sits — a paragraph, a heading, a table cell, a list item, inside emphasis
//     or a link. A Text node under a CodeSpan is the code-span surface.
//
// Container nesting needs no handling at all: by the time the AST exists, the
// parser has already resolved blockquote and list prefixes, so a fence inside a
// blockquote is simply a fence.
func collectMentionSurfaces(bodySource []byte, bodyRoot ast.Node) []surfaceRange {
	var surfaces []surfaceRange
	appendSegment := func(segment text.Segment, surface mentionSurface) {
		if segment.Stop > segment.Start {
			surfaces = append(surfaces, surfaceRange{start: segment.Start, stop: segment.Stop, surface: surface})
		}
	}
	appendLines := func(node ast.Node, surface mentionSurface) {
		lines := node.Lines()
		for lineIndex := 0; lineIndex < lines.Len(); lineIndex++ {
			appendSegment(lines.At(lineIndex), surface)
		}
	}

	_ = ast.Walk(bodyRoot, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch typedNode := node.(type) {
		case *ast.FencedCodeBlock:
			appendLines(typedNode, surfaceCodeBlock)
			if typedNode.Info != nil {
				appendSegment(typedNode.Info.Segment, surfaceCodeBlock)
			}
		case *ast.CodeBlock:
			appendLines(typedNode, surfaceCodeBlock)
		case *ast.Text:
			appendSegment(typedNode.Segment, textNodeSurface(typedNode))
		}
		return ast.WalkContinue, nil
	})

	sort.Slice(surfaces, func(leftIndex, rightIndex int) bool {
		return surfaces[leftIndex].start < surfaces[rightIndex].start
	})
	return surfaces
}

// textNodeSurface reports what a prose text node is really inside — a link
// label, a code span, or open prose. Walking the whole ancestor chain rather
// than checking the immediate parent, because an extension may wrap the
// contents of either.
//
// A LINK ANYWHERE UP THE CHAIN WINS, even over a code span nested inside it.
// The drawer's rule is `parentElement.closest("a")`, which returns early for
// every text node under an anchor without looking at what else encloses it; a
// code span inside a link label earns no drawer glossary line, so it must not
// earn a clipboard one either. Taking the innermost construct instead read as
// the narrower, more careful claim and was simply a different answer from the
// drawer's — which is the one thing this pair may never be.
func textNodeSurface(textNode ast.Node) mentionSurface {
	enclosingSurface := surfaceProse
	for ancestor := textNode.Parent(); ancestor != nil; ancestor = ancestor.Parent() {
		switch ancestor.Kind() {
		// Link and Image only. An ast.AutoLink holds its text in a private
		// `value *Text` field rather than as a child, so no walked *ast.Text can
		// ever have one as an ancestor and an AutoLink arm here would be dead
		// code. Autolinked ids are already excluded anyway: the URL alternative
		// of the mention pattern claims the whole run.
		case ast.KindLink, ast.KindImage:
			return surfaceLinkLabel
		case ast.KindCodeSpan:
			enclosingSurface = surfaceCodeSpan
		}
	}
	return enclosingSurface
}

// surfaceAt reports the surface of the byte range [start, stop), and whether
// any single range covers it. A mention straddling two ranges is not covered:
// it is half prose and half something else, which no annotation can be correct
// about.
func surfaceAt(surfaces []surfaceRange, start int, stop int) (mentionSurface, bool) {
	candidateIndex := sort.Search(len(surfaces), func(index int) bool {
		return surfaces[index].start > start
	}) - 1
	if candidateIndex < 0 {
		return surfaceProse, false
	}
	candidate := surfaces[candidateIndex]
	if start < candidate.start || stop > candidate.stop {
		return surfaceProse, false
	}
	return candidate.surface, true
}

// ---- resolving a mention to a board record --------------------------------

// ticketMentionResolver answers "which board record does this written id name?"
// Ported from board-core.js (buildRequestIdByReqSegment / resolveTicketMention /
// isAmbiguousTicketMention) rather than shared, because the drawer resolves the
// same mentions in the browser against the same board.
//
// The two copies are pinned, not merged:
// TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo drives both
// over one shared corpus, so whichever side drifts alone fails.
type ticketMentionResolver struct {
	requestIds     map[string]bool
	userRequestIds map[string]bool
	// REQ segment → the one card id carrying it. An empty value marks an
	// AMBIGUOUS segment: two cards share it, and the board never guesses.
	requestIdByReqSegment map[string]string
}

func newTicketMentionResolver(requestIds []string, userRequestIds []string) *ticketMentionResolver {
	resolver := &ticketMentionResolver{
		requestIds:            make(map[string]bool, len(requestIds)),
		userRequestIds:        make(map[string]bool, len(userRequestIds)),
		requestIdByReqSegment: map[string]string{},
	}
	for _, requestId := range requestIds {
		resolver.requestIds[requestId] = true
	}
	for _, userRequestId := range userRequestIds {
		resolver.userRequestIds[userRequestId] = true
	}
	for _, requestId := range requestIds {
		segmentMatch := requestIdSegmentPattern.FindString(requestId)
		if segmentMatch == "" || segmentMatch == requestId {
			continue // not a compound id — the short form IS the card id
		}
		segmentKey := upperAsciiId(segmentMatch)
		if _, alreadyIndexed := resolver.requestIdByReqSegment[segmentKey]; alreadyIndexed {
			resolver.requestIdByReqSegment[segmentKey] = "" // ambiguous — never guess
			continue
		}
		resolver.requestIdByReqSegment[segmentKey] = requestId
	}
	return resolver
}

// resolve returns the record kind ("req" or "ur") and the full board id a
// mention names, or two empty strings when nothing on the board answers to it.
func (resolver *ticketMentionResolver) resolve(mentionText string) (string, string) {
	if resolver.requestIds[mentionText] {
		return "req", mentionText
	}
	if resolver.userRequestIds[mentionText] {
		return "ur", mentionText
	}
	if segmentTargetId := resolver.requestIdByReqSegment[upperAsciiId(mentionText)]; segmentTargetId != "" {
		return "req", segmentTargetId
	}
	return "", ""
}

// isAmbiguous reports a segment two cards share. Such a mention is NOT a dead
// reference — the board holds records that match and refuses to pick one — so a
// caller flagging unresolved ids must leave it alone or the never-guess rule
// turns into a false alarm.
func (resolver *ticketMentionResolver) isAmbiguous(mentionText string) bool {
	segmentTargetId, indexed := resolver.requestIdByReqSegment[upperAsciiId(mentionText)]
	return indexed && segmentTargetId == ""
}

// upperAsciiId upper-cases a ticket id for segment lookup. Ticket ids are ASCII
// by construction, and a locale-free byte fold is what the client's
// String.prototype.toUpperCase does to them.
func upperAsciiId(ticketId string) string {
	upperBytes := []byte(ticketId)
	for byteIndex, currentByte := range upperBytes {
		if currentByte >= 'a' && currentByte <= 'z' {
			upperBytes[byteIndex] = currentByte - ('a' - 'A')
		}
	}
	return string(upperBytes)
}

// ---- the shipped index ----------------------------------------------------

// generatedTicketMention is one REQ/UR id written in a ticket body, located in
// the CLIPBOARD DOCUMENT (frontmatter fence + body) the same build ships, and
// already resolved. The client splices at Offset+Length and needs no pattern,
// no resolver and no block reasoning of its own.
//
// Offset and Length are UTF-16 CODE UNITS, not bytes: the client indexes a
// JavaScript string, and the two measurements diverge at the first byte that is
// not ASCII — after which a byte offset lands mid-word. An em dash costs three
// bytes and one unit; an emoji costs four bytes and two. Nearly every body here
// carries one (372 of this board's 376 requests hold an em dash), and the four
// that do not are exactly the ones where a byte offset would have looked right.
//
// Kind is "req", "ur", or empty for a mention no board record answers to — the
// dead reference a paste's reader learns about from the appendix, since plain
// text has no red. Id is the resolved board id, or the written text when Kind
// is empty. Expand marks the one occurrence per id that may carry a title: the
// first in prose. It is never set on a quoted mention or a dead one.
type generatedTicketMention struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Kind   string `json:"kind,omitempty"`
	Id     string `json:"id"`
	Expand bool   `json:"expand,omitempty"`
}

// collectDocumentTicketMentions locates every annotatable ticket mention in one
// clipboard document — the WHOLE document the client holds, fence and all, so
// the offsets it returns index the exact string the client will splice into.
//
// The fence itself is never scanned: it carries depends_on, related and
// user_request, and annotating any of them would stop the paste parsing as
// YAML. splitFrontmatter draws that line, and it is the same function the file
// was parsed with, so the fence a paste round-trips through and the fence
// skipped here cannot come apart.
//
// The body is parsed as it is STORED. renderMarkdownBodyToHtml preprocesses its
// input first (insertQuestionOptionHardBreaks appends hard breaks), which shifts
// every offset after the first break it adds — so the two must not share a
// parse, and this one takes the raw text.
func collectDocumentTicketMentions(documentText string, resolver *ticketMentionResolver) []generatedTicketMention {
	_, bodyMarkdown, bodyStartOffset, _ := splitFrontmatter(documentText)
	if bodyMarkdown == "" {
		return nil
	}
	bodySource := []byte(bodyMarkdown)
	bodyRoot := markdownToHtmlRenderer.Parser().Parse(text.NewReader(bodySource))
	surfaces := collectMentionSurfaces(bodySource, bodyRoot)

	documentOffsets := newUtf16OffsetCursor(bodyMarkdown, utf16LengthOf(documentText[:bodyStartOffset]))
	expandedIds := map[string]bool{}
	var mentions []generatedTicketMention

	for _, matchIndexes := range bodyTicketMentionPattern.FindAllStringSubmatchIndex(bodyMarkdown, -1) {
		// A match on the URL or repo-path alternative claimed the run and names
		// no ticket, so the named group is absent and its indexes are -1.
		mentionStart := matchIndexes[2*ticketMentionGroupIndex]
		mentionStop := matchIndexes[2*ticketMentionGroupIndex+1]
		if mentionStart < 0 || !isTicketMentionBoundary(bodyMarkdown, mentionStart-1) ||
			!isTicketMentionBoundary(bodyMarkdown, mentionStop) {
			continue
		}
		mentionText := bodyMarkdown[mentionStart:mentionStop]
		surface, covered := surfaceAt(surfaces, mentionStart, mentionStop)
		if !covered || surface == surfaceLinkLabel {
			// Either goldmark keeps no text here — a link reference definition, a
			// raw HTML block, a fence's own characters — or it does and touching
			// it would break something (see surfaceLinkLabel).
			continue
		}

		kind, resolvedId := resolver.resolve(mentionText)
		if kind == "" && (surface == surfaceCodeBlock || resolver.isAmbiguous(mentionText)) {
			continue
		}

		// Read the cursor exactly twice, ascending. Building these inline in the
		// composite literal below called at(start) again AFTER at(stop), which
		// takes the rewind branch and rescans the whole prefix — the quadratic
		// the cursor exists to avoid, paid twice per mention.
		mentionOffset := documentOffsets.at(mentionStart)
		mentionEndOffset := documentOffsets.at(mentionStop)

		if kind == "" {
			mentions = append(mentions, generatedTicketMention{
				Offset: mentionOffset,
				Length: mentionEndOffset - mentionOffset,
				Id:     mentionText,
			})
			continue
		}

		expand := surface == surfaceProse && !expandedIds[resolvedId]
		if expand {
			expandedIds[resolvedId] = true
		}
		mentions = append(mentions, generatedTicketMention{
			Offset: mentionOffset,
			Length: mentionEndOffset - mentionOffset,
			Kind:   kind,
			Id:     resolvedId,
			Expand: expand,
		})
	}
	return mentions
}

// ---- byte offsets → JavaScript string offsets -----------------------------

// utf16OffsetCursor converts ascending byte offsets in a Go string to the
// UTF-16 code-unit offsets a JavaScript string is indexed by, carrying a fixed
// base (the frontmatter fence the body sits behind).
//
// It is a cursor rather than a function because a body is scanned once in
// order: re-measuring the prefix per mention would be quadratic on the longest
// bodies in the tree.
type utf16OffsetCursor struct {
	sourceText string
	byteOffset int
	unitOffset int
}

func newUtf16OffsetCursor(sourceText string, baseUnitOffset int) *utf16OffsetCursor {
	return &utf16OffsetCursor{sourceText: sourceText, unitOffset: baseUnitOffset}
}

// at returns the UTF-16 offset of a byte offset. Offsets must arrive in
// ascending order; an out-of-order offset would silently under-count, so it
// restarts from the beginning rather than answering wrongly.
func (cursor *utf16OffsetCursor) at(byteOffset int) int {
	if byteOffset < cursor.byteOffset {
		base := cursor.unitOffset - utf16LengthOf(cursor.sourceText[:cursor.byteOffset])
		cursor.byteOffset = 0
		cursor.unitOffset = base
	}
	cursor.unitOffset += utf16LengthOf(cursor.sourceText[cursor.byteOffset:byteOffset])
	cursor.byteOffset = byteOffset
	return cursor.unitOffset
}

// utf16LengthOf counts the UTF-16 code units a string occupies once it is a
// JavaScript string. Every rune outside the Basic Multilingual Plane costs two;
// an invalid byte decodes to one U+FFFD, which is exactly what encoding/json
// writes for it, so the count still matches what the client receives.
func utf16LengthOf(sourceText string) int {
	unitCount := 0
	for _, currentRune := range sourceText {
		if currentRune > 0xFFFF {
			unitCount += 2
			continue
		}
		unitCount++
	}
	return unitCount
}
