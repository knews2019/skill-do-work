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
// CommonMark on six separate constructs. Go can answer it exactly: render.go
// already parses every one of these bodies with goldmark to produce the
// drawer's HTML, so this file walks the same AST once more and ships the
// answer as positions. filementions.go is the same move for the same reason:
// the client cannot stat the filesystem, so Go ships repoFileMentions.
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
// This MUST stay in lock-step with bodyMentionPattern in web/board-detail.js,
// which scans the drawer's rendered text with the same three alternatives in
// the same order. TestTicketMentionPatternAgreesBetweenGoAndTheClient drives
// both over one shared corpus so a drift fails rather than silently changing
// what a paste says.
var bodyTicketMentionPattern = regexp.MustCompile(
	`(https?://[^\s<>"')\]]+)` +
		`|((?:[A-Za-z0-9_@-]+(?:\.[A-Za-z0-9_-]+)*/)+[A-Za-z0-9_@-][A-Za-z0-9_@.-]*\.[A-Za-z][A-Za-z0-9]{0,7})` +
		`|(\b(?:UR-\d+-REQ-\d+[a-z]?|REQ-\d+[a-z]?|UR-\d+)\b)`)

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
//
// A mention covered by none of the three (a link reference definition, a raw
// HTML block, a fence's own backticks) is not annotatable at all and is
// dropped: goldmark keeps no prose text there, so there is nothing to expand.
type mentionSurface int

const (
	surfaceProse mentionSurface = iota
	surfaceCodeSpan
	surfaceCodeBlock
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

// textNodeSurface reports whether a prose text node is actually the contents of
// a code span. Walking ancestors rather than checking the immediate parent
// because an extension may wrap a span's contents.
func textNodeSurface(textNode ast.Node) mentionSurface {
	for ancestor := textNode.Parent(); ancestor != nil; ancestor = ancestor.Parent() {
		if ancestor.Kind() == ast.KindCodeSpan {
			return surfaceCodeSpan
		}
	}
	return surfaceProse
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
// TestTicketMentionResolutionAgreesBetweenGoAndTheClient drives both over one
// shared corpus, so whichever side drifts alone fails.
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
// JavaScript string, where an em dash is one unit and a byte offset would land
// mid-word in any body containing one. Every REQ body in this repo contains one.
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
		// Submatch 3 is the ticket-id alternative; a match on 1 (URL) or 2
		// (repo-relative path) claimed the run and names no ticket.
		mentionStart, mentionStop := matchIndexes[6], matchIndexes[7]
		if mentionStart < 0 {
			continue
		}
		mentionText := bodyMarkdown[mentionStart:mentionStop]
		surface, covered := surfaceAt(surfaces, mentionStart, mentionStop)
		if !covered {
			continue // goldmark keeps no text here: a link reference definition, a raw HTML block, a fence's own characters.
		}

		kind, resolvedId := resolver.resolve(mentionText)
		if kind == "" {
			if surface == surfaceCodeBlock || resolver.isAmbiguous(mentionText) {
				continue
			}
			mentions = append(mentions, generatedTicketMention{
				Offset: documentOffsets.at(mentionStart),
				Length: documentOffsets.at(mentionStop) - documentOffsets.at(mentionStart),
				Id:     mentionText,
			})
			continue
		}

		expand := surface == surfaceProse && !expandedIds[resolvedId]
		if expand {
			expandedIds[resolvedId] = true
		}
		mentions = append(mentions, generatedTicketMention{
			Offset: documentOffsets.at(mentionStart),
			Length: documentOffsets.at(mentionStop) - documentOffsets.at(mentionStart),
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
