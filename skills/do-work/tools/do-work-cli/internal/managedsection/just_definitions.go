package managedsection

import (
	"bytes"
	"regexp"
)

// The Just definition scanner exists for one job: deciding whether a consumer's justfile
// already defines a name the managed section reserves, without requiring `just` on PATH. It
// is a state machine over Just's five literal forms because a colon inside a string default
// is not a recipe header, and a naive line scan reads it as one.
var (
	justIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*`)
	justAliasPattern      = regexp.MustCompile(`^alias[ \t]+([A-Za-z_][A-Za-z0-9_-]*)[ \t]*:=`)
)

var byteOrderMark = []byte{0xef, 0xbb, 0xbf}

// JustDefinitionNames reports every recipe and alias name the source defines. A BOM is
// stripped from line zero's classification view only; the raw bytes are never touched, which
// is what keeps marker matching and definition scanning deliberately asymmetric.
func JustDefinitionNames(data []byte) map[string]struct{} {
	definitionNames := map[string]struct{}{}
	var activeDelimiter []byte
	var pendingDefinitionLines [][]byte
	lines, _ := splitLinesKeepingEnds(data)
	for lineIndex, line := range lines {
		classificationLine := line
		if lineIndex == 0 && bytes.HasPrefix(line, byteOrderMark) {
			classificationLine = line[len(byteOrderMark):]
		}
		lineStartsInsideLiteral := activeDelimiter != nil
		if lineStartsInsideLiteral {
			pendingDefinitionLines = append(pendingDefinitionLines, classificationLine)
		} else {
			pendingDefinitionLines = [][]byte{classificationLine}
		}
		activeDelimiter = justMultilineLiteralState(classificationLine, activeDelimiter)
		if lineStartsInsideLiteral && activeDelimiter != nil {
			continue
		}
		definitionSource := classificationLine
		if lineStartsInsideLiteral {
			definitionSource = bytes.Join(pendingDefinitionLines, nil)
		}
		if definitionName := justDefinitionName(definitionSource); definitionName != nil {
			definitionNames[string(definitionName)] = struct{}{}
		}
		if activeDelimiter == nil {
			pendingDefinitionLines = nil
		}
	}
	return definitionNames
}

// justDefinitionName returns the recipe or alias name a definition source declares, or nil
// when the source is a comment, a body line, an assignment, or an unterminated header.
func justDefinitionName(line []byte) []byte {
	body := lineBody(line)
	if len(body) == 0 || body[0] == ' ' || body[0] == '\t' || body[0] == '#' {
		return nil
	}
	if aliasMatch := justAliasPattern.FindSubmatch(body); aliasMatch != nil {
		return aliasMatch[1]
	}

	nameOffset := 0
	if body[0] == '@' {
		nameOffset = 1
	}
	if nameOffset >= len(body) {
		return nil
	}
	nameMatch := justIdentifierPattern.Find(body[nameOffset:])
	if nameMatch == nil {
		return nil
	}

	remainder := body[nameOffset+len(nameMatch):]
	var activeDelimiter []byte
	index := 0
	for index < len(remainder) {
		if activeDelimiter != nil {
			if justDelimiterMatches(remainder, index, activeDelimiter) {
				index += len(activeDelimiter)
				activeDelimiter = nil
			} else {
				index++
			}
			continue
		}
		if openingDelimiter := justOpeningDelimiter(remainder, index); openingDelimiter != nil {
			activeDelimiter = openingDelimiter
			index += len(openingDelimiter)
			continue
		}
		// The first bare colon that is not `:=` makes this a recipe header.
		if remainder[index] == ':' {
			if index+1 < len(remainder) && remainder[index+1] == '=' {
				return nil
			}
			return nameMatch
		}
		index++
	}
	return nil
}

// justMultilineLiteralState carries literal state across lines. A line that starts OUTSIDE a
// literal and begins with space, tab or `#`, or is empty, resets the state — which is how a
// recipe body or a comment can never open one.
func justMultilineLiteralState(line []byte, activeDelimiter []byte) []byte {
	body := lineBody(line)
	if activeDelimiter == nil && (len(body) == 0 || body[0] == ' ' || body[0] == '\t' || body[0] == '#') {
		return nil
	}
	index := 0
	for index < len(body) {
		if activeDelimiter != nil {
			if justDelimiterMatches(body, index, activeDelimiter) {
				index += len(activeDelimiter)
				activeDelimiter = nil
			} else {
				index++
			}
			continue
		}
		if body[index] == '#' {
			break
		}
		if openingDelimiter := justOpeningDelimiter(body, index); openingDelimiter != nil {
			activeDelimiter = openingDelimiter
			index += len(openingDelimiter)
			continue
		}
		index++
	}
	return activeDelimiter
}

// justDelimiterMatches decides whether the active literal closes here. Two rules earn their
// place: a triple backtick adjacent to another backtick is content rather than a delimiter,
// and only the cooked forms (`"` and `"""`) honour a backslash escape — Just's raw strings
// and backtick commands do not.
func justDelimiterMatches(body []byte, index int, activeDelimiter []byte) bool {
	if index >= len(body) || !bytes.HasPrefix(body[index:], activeDelimiter) {
		return false
	}
	if string(activeDelimiter) == "```" {
		if (index > 0 && body[index-1] == '`') ||
			(index+len(activeDelimiter) < len(body) && body[index+len(activeDelimiter)] == '`') {
			return false
		}
	}
	if string(activeDelimiter) == `"` || string(activeDelimiter) == `"""` {
		backslashCount := 0
		for backslashIndex := index - 1; backslashIndex >= 0 && body[backslashIndex] == '\\'; backslashIndex-- {
			backslashCount++
		}
		if backslashCount%2 == 1 {
			return false
		}
	}
	return true
}

// justOpeningDelimiter recognises Just's literal openers in priority order, so a triple form
// is never mistaken for its single-character prefix.
func justOpeningDelimiter(body []byte, index int) []byte {
	if index >= len(body) {
		return nil
	}
	remainder := body[index:]
	if bytes.HasPrefix(remainder, []byte("'''")) {
		return []byte("'''")
	}
	if bytes.HasPrefix(remainder, []byte(`"""`)) {
		return []byte(`"""`)
	}
	if bytes.HasPrefix(remainder, []byte("```")) &&
		(index == 0 || body[index-1] != '`') &&
		(index+3 == len(body) || body[index+3] != '`') {
		return []byte("```")
	}
	switch body[index] {
	case '"', '\'', '`':
		return []byte{body[index]}
	}
	return nil
}
