package requestmodel

import "strings"

// VisibleSection identifies an unquoted level-two section in original body bytes.
// End excludes an unclosed fenced/comment region so destructive writers retain it.
type VisibleSection struct {
	Name  string
	Start int
	End   int
}

// VisibleSections discovers sections without normalizing line endings or offsets.
// It shares the timing writer's conservative dialect-fence rule: enclosing runs
// of punctuation hide their contents, while headings and thematic breaks do not.
// A heading or fence counts only after CommonMark's block indent (see
// blockIndentContent). On a line that starts outside a comment, a fence opener is
// decided before any comment scan, because the rest of that line is its info
// string, and a "<!--" inside an inline code span on the same line is literal.
func VisibleSections(body []byte) []VisibleSection {
	lines := strings.SplitAfter(string(body), "\n")
	sections := []VisibleSection{}
	openCharacter, openLength, hiddenStart := byte(0), 0, -1
	inComment := false
	offset := 0
	for _, rawLine := range lines {
		start := offset
		offset += len(rawLine)
		line := strings.TrimSuffix(rawLine, "\n")
		character, length := leadingPunctuationRun(line)
		if openLength > 0 {
			if character == openCharacter && length >= openLength {
				openCharacter, openLength, hiddenStart = 0, 0, -1
			}
			continue
		}
		startedInComment := inComment
		if !startedInComment && length >= 3 && isEnclosingFenceCharacter(character) {
			openCharacter, openLength, hiddenStart = character, length, start
			continue
		}
		remainder := line
		for {
			if inComment {
				end := strings.Index(remainder, "-->")
				if end < 0 {
					break
				}
				remainder = remainder[end+3:]
				inComment, hiddenStart = false, -1
			}
			begin := commentOpenerIndex(remainder)
			if begin < 0 {
				break
			}
			inComment, hiddenStart = true, start
			remainder = remainder[begin+4:]
		}
		if startedInComment {
			continue
		}
		// Inline comments do not hide a heading's visible prefix. Keep its raw
		// name so boundary recognition does not broaden generated-section ownership.
		content, isBlockStart := blockIndentContent(line)
		if !isBlockStart {
			continue
		}
		if name, found := strings.CutPrefix(content, "## "); found {
			name = strings.TrimRight(name, " \t\r")
			if name == "" {
				continue
			}
			if len(sections) > 0 {
				sections[len(sections)-1].End = start
			}
			sections = append(sections, VisibleSection{Name: name, Start: start, End: len(body)})
		}
	}
	if hiddenStart >= 0 && len(sections) > 0 {
		sections[len(sections)-1].End = hiddenStart
	}
	return sections
}

// leadingPunctuationRun returns the ASCII punctuation mark a line opens with and
// how many times it repeats. The mark counts only inside CommonMark's block
// indent (see blockIndentContent); deeper indentation is code, so it opens or
// closes no fence.
func leadingPunctuationRun(line string) (byte, int) {
	content, isBlockStart := blockIndentContent(strings.TrimRight(line, " \t\r"))
	if !isBlockStart || content == "" || !isMarkdownBlockPunctuation(content[0]) {
		return 0, 0
	}
	runLength := 0
	for runLength < len(content) && content[runLength] == content[0] {
		runLength++
	}
	return content[0], runLength
}

// blockIndentContent strips CommonMark's block indent: zero to three spaces.
// A tab anywhere in the indentation, or a fourth space, makes the line indented
// code, so it reports false and nothing on it can open a heading or a fence. The
// indent is measured before anything is trimmed, because the indent itself is
// the structure being tested.
func blockIndentContent(line string) (string, bool) {
	spaces := 0
	for spaces < len(line) && line[spaces] == ' ' {
		spaces++
	}
	if spaces > 3 || spaces < len(line) && line[spaces] == '\t' {
		return "", false
	}
	return line[spaces:], true
}

// commentOpenerIndex finds the first "<!--" outside an inline code span. A run of
// n backticks opens a span that closes at the next run of exactly n backticks on
// the same text, and everything between is literal. A run with no closer is
// literal text itself, so a "<!--" after it still opens a comment. It returns -1
// when no opener is visible.
func commentOpenerIndex(text string) int {
	for index := 0; index < len(text); {
		if strings.HasPrefix(text[index:], "<!--") {
			return index
		}
		if text[index] != '`' {
			index++
			continue
		}
		runLength := backtickRunLength(text, index)
		index += runLength
		for search := index; search < len(text); {
			if text[search] != '`' {
				search++
				continue
			}
			closeLength := backtickRunLength(text, search)
			if closeLength == runLength {
				index = search + closeLength
				break
			}
			search += closeLength
		}
	}
	return -1
}

func backtickRunLength(text string, start int) int {
	end := start
	for end < len(text) && text[end] == '`' {
		end++
	}
	return end - start
}

// isMarkdownBlockPunctuation is CommonMark's own ASCII punctuation class, taken
// wholesale rather than narrowed to the marks today's fence syntax happens to
// use. The narrower set would be an enumeration to revisit whenever a dialect
// adds a construct, which is exactly how a fence-blind classifier reopens.
func isMarkdownBlockPunctuation(value byte) bool {
	return value >= '!' && value <= '/' ||
		value >= ':' && value <= '@' ||
		value >= '[' && value <= '`' ||
		value >= '{' && value <= '~'
}

// isEnclosingFenceCharacter excludes the marks CommonMark defines as line
// constructs that never enclose anything: thematic breaks use "-", "_" and "*",
// setext underlines use "-" and "=", and "#" opens an ATX heading. Treating
// those as fences would let the unpaired "---" a request carries above its
// source line swallow the rest of the document, or an ordinary "### " heading
// swallow the request's own Timing section, so that section could never be
// found again. The exclusion is by construct, not by run length: no run of "#"
// encloses anything, since seven or more open nothing at all.
func isEnclosingFenceCharacter(value byte) bool {
	if value == 0 || !isMarkdownBlockPunctuation(value) {
		return false
	}
	return value != '-' && value != '_' && value != '*' && value != '=' && value != '#'
}
