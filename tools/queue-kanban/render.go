package main

import (
	"bytes"

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
// must collapse to a space (standard Markdown) rather than a forced <br>. Raw
// HTML embedded in a body stays ESCAPED (WithUnsafe is NOT set), so a malicious
// or accidental `<script>` in a Markdown body is rendered as inert text rather
// than executed in the shareable artifact.
var markdownToHtmlRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

// renderMarkdownBodyToHtml renders one Markdown body to a safe HTML fragment.
// An empty body yields an empty string; a render failure yields the error so the
// caller can decide whether to skip the body or fail the whole generate.
func renderMarkdownBodyToHtml(markdownBody string) (string, error) {
	if markdownBody == "" {
		return "", nil
	}
	var htmlBuffer bytes.Buffer
	if convertError := markdownToHtmlRenderer.Convert([]byte(markdownBody), &htmlBuffer); convertError != nil {
		return "", convertError
	}
	return htmlBuffer.String(), nil
}
