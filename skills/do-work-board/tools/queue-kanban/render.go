package main

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// markdownToHtmlRenderer converts a REQ/UR Markdown body to HTML at generate
// time so the static page needs no runtime renderer. GFM is enabled for tables,
// strikethrough, autolinks, and — importantly for the REQ format — task-list
// checkboxes (the `- [ ]` / `- [x]` AI Execution State items render as disabled
// checkbox inputs). Hard-wraps are deliberately NOT enabled: REQ bodies wrap
// their prose at a column width, so a single source newline is a soft wrap that
// must collapse to a space (standard Markdown) rather than a forced <br>. The
// one exception is the Open Questions option format, handled by
// insertQuestionOptionHardBreaks below — those continuation lines are
// semantically separate lines, not wrapped prose. Raw
// HTML embedded in a body stays ESCAPED (WithUnsafe is NOT set), so a malicious
// or accidental `<script>` in a Markdown body is rendered as inert text rather
// than executed in the shareable artifact.
var markdownToHtmlRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

// questionOptionLinePrefixes are the continuation-line keywords of the REQ
// Open Questions format (see actions/capture.md "Open Questions" and
// actions/clarify.md's Recommended:/Also: fallback). In the source they sit as
// indented lines under a `- [ ]` checkbox item; plain Markdown treats them as
// lazy paragraph continuations and collapses them into the question sentence.
// The preprocessor below re-establishes them as their own visual lines.
var questionOptionLinePrefixes = []string{
	"Recommended:",
	"Also:",
	"Value:",
	"Risk:",
	"→",
}

// insertQuestionOptionHardBreaks appends a Markdown hard-break (`\`) to any
// line directly followed by a question-option continuation line, so each
// `Recommended:` / `Also:` / `Value:` / `Risk:` / `→` line renders on its own
// line instead of merging into the question paragraph. Lines inside fenced
// code blocks are left untouched — a fence's content is verbatim by contract.
func insertQuestionOptionHardBreaks(markdownBody string) string {
	bodyLines := strings.Split(markdownBody, "\n")
	insideCodeFence := false
	for lineIndex, currentLine := range bodyLines {
		trimmedLine := strings.TrimSpace(currentLine)
		if strings.HasPrefix(trimmedLine, "```") || strings.HasPrefix(trimmedLine, "~~~") {
			insideCodeFence = !insideCodeFence
			continue
		}
		if insideCodeFence || lineIndex == 0 {
			continue
		}
		startsWithOptionPrefix := false
		for _, optionPrefix := range questionOptionLinePrefixes {
			if strings.HasPrefix(trimmedLine, optionPrefix) {
				startsWithOptionPrefix = true
				break
			}
		}
		if !startsWithOptionPrefix {
			continue
		}
		previousLine := bodyLines[lineIndex-1]
		// Only a non-blank previous line forms a paragraph this line would lazily
		// continue; after a blank line the option line already starts fresh.
		if strings.TrimSpace(previousLine) == "" {
			continue
		}
		if strings.HasSuffix(previousLine, "\\") || strings.HasSuffix(previousLine, "  ") {
			continue
		}
		bodyLines[lineIndex-1] = previousLine + "\\"
	}
	return strings.Join(bodyLines, "\n")
}

// renderMarkdownBodyToHtml renders one Markdown body to a safe HTML fragment.
// An empty body yields an empty string; a render failure yields the error so the
// caller can decide whether to skip the body or fail the whole generate.
func renderMarkdownBodyToHtml(markdownBody string) (string, error) {
	if markdownBody == "" {
		return "", nil
	}
	var htmlBuffer bytes.Buffer
	preprocessedBody := insertQuestionOptionHardBreaks(markdownBody)
	if convertError := markdownToHtmlRenderer.Convert([]byte(preprocessedBody), &htmlBuffer); convertError != nil {
		return "", convertError
	}
	return htmlBuffer.String(), nil
}
