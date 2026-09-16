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
		hadComment := inComment
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
			begin := strings.Index(remainder, "<!--")
			if begin < 0 {
				break
			}
			inComment, hadComment, hiddenStart = true, true, start
			remainder = remainder[begin+4:]
		}
		if hadComment {
			continue
		}
		if length >= 3 && isEnclosingFenceCharacter(character) {
			openCharacter, openLength, hiddenStart = character, length, start
			continue
		}
		if name, found := strings.CutPrefix(line, "## "); found {
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
// how many times it repeats. Leading whitespace is skipped rather than measured,
// so an indented fence still counts as one.
func leadingPunctuationRun(line string) (byte, int) {
	content := strings.TrimLeft(strings.TrimRight(line, " \t\r"), " \t")
	if content == "" || !isMarkdownBlockPunctuation(content[0]) {
		return 0, 0
	}
	runLength := 0
	for runLength < len(content) && content[runLength] == content[0] {
		runLength++
	}
	return content[0], runLength
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
