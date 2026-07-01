package main

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// splitFrontmatter separates a leading "---\n … \n---" YAML block from the
// Markdown body of a REQ/UR file. It returns the YAML text (the lines between
// the fences, without the fence lines), the body text (everything after the
// closing fence, kept verbatim for later Markdown rendering), and whether a
// frontmatter block was present at all.
//
// Files with no leading frontmatter return ("", originalContent, false) so the
// caller can skip gracefully. Parsing is purely mechanical (prefix/line splits,
// no regexp): the opening fence must be the very first line "---", and the
// closing fence is the next line that is exactly "---".
func splitFrontmatter(fileContent string) (yamlText string, bodyText string, hasFrontmatter bool) {
	// Normalize CRLF and strip a leading UTF-8 BOM so the fence checks below are
	// simple equality tests against "---".
	normalized := strings.ReplaceAll(fileContent, "\r\n", "\n")
	normalized = strings.TrimPrefix(normalized, "\ufeff")

	const openingFence = "---\n"
	if !strings.HasPrefix(normalized, openingFence) {
		return "", fileContent, false
	}

	afterOpening := normalized[len(openingFence):]
	lines := strings.Split(afterOpening, "\n")

	closingLineIndex := -1
	for lineIndex, line := range lines {
		if line == "---" {
			closingLineIndex = lineIndex
			break
		}
	}
	if closingLineIndex < 0 {
		// No closing fence — treat the file as having no frontmatter rather than
		// swallowing the whole document as YAML.
		return "", fileContent, false
	}

	yamlText = strings.Join(lines[:closingLineIndex], "\n")
	bodyText = strings.Join(lines[closingLineIndex+1:], "\n")
	return yamlText, bodyText, true
}

// parseFrontmatterFields unmarshals a YAML frontmatter block into a permissive
// map. A map (rather than a rigid struct) is used deliberately: REQ frontmatter
// in the wild mixes scalar, list, and occasionally numeric values for the same
// logical field, and a map lets the field coercion helpers (coerceScalarToString
// / coerceToStringList) accept whatever shape a given file happens to use without
// failing the whole parse on a type mismatch.
//
// One real file in the tree repeats a top-level key (two completed_at lines),
// which yaml.v3 rejects as a duplicate-key error. On any unmarshal error the
// function retries once against a de-duplicated copy that keeps the LAST value of
// each repeated top-level key, so a single malformed file is recovered rather
// than dropped.
func parseFrontmatterFields(yamlText string) (map[string]any, error) {
	fields := map[string]any{}
	if strings.TrimSpace(yamlText) == "" {
		return fields, nil
	}
	unmarshalError := yaml.Unmarshal([]byte(yamlText), &fields)
	if unmarshalError == nil {
		return fields, nil
	}

	deduplicated := dropDuplicateTopLevelKeys(yamlText)
	retryFields := map[string]any{}
	if retryError := yaml.Unmarshal([]byte(deduplicated), &retryFields); retryError == nil {
		return retryFields, nil
	}

	// Last resort: a handful of real files carry a malformed title line (a bare
	// colon, e.g. "title: Review fix: resolve …", or a quoted prefix followed by
	// more text, e.g. `title: "Clean up broken" button — …`). Strict YAML rejects
	// the whole block, which would silently drop the REQ's status, UR pointer, and
	// dependencies. A line-based extraction recovers the remaining top-level fields
	// so one bad line doesn't lose the record. unmarshalError is intentionally not
	// returned — recovery is the contract here.
	_ = unmarshalError
	return lenientFrontmatterFields(yamlText), nil
}

// lenientFrontmatterFields recovers top-level scalar and flow-list fields from a
// frontmatter block that strict YAML could not parse. It splits each unindented
// "key: value" line on its FIRST colon (so a colon inside the value survives in
// the value), unquotes a fully-quoted scalar, and expands a "[a, b]" flow list.
// Block (multi-line) lists and nested maps are not recovered — only the flat
// top-level fields the board model reads.
func lenientFrontmatterFields(yamlText string) map[string]any {
	fields := map[string]any{}
	for _, line := range strings.Split(yamlText, "\n") {
		key, isKey := topLevelKeyName(line)
		if !isKey {
			continue
		}
		rawValue := strings.TrimSpace(line[len(key)+1:])
		if rawValue == "" {
			continue
		}
		fields[key] = parseLenientScalarOrList(rawValue)
	}
	return fields
}

// parseLenientScalarOrList turns a raw frontmatter value into a flow-list slice
// when it is bracketed, otherwise into an unquoted scalar string.
func parseLenientScalarOrList(rawValue string) any {
	if strings.HasPrefix(rawValue, "[") && strings.HasSuffix(rawValue, "]") {
		inner := strings.TrimSpace(rawValue[1 : len(rawValue)-1])
		if inner == "" {
			return []any{}
		}
		parts := strings.Split(inner, ",")
		list := make([]any, 0, len(parts))
		for _, part := range parts {
			list = append(list, unquoteScalar(strings.TrimSpace(part)))
		}
		return list
	}
	return unquoteScalar(rawValue)
}

// unquoteScalar strips a single matching pair of surrounding quotes; a value
// that is only partially quoted (a malformed title) is returned verbatim.
func unquoteScalar(value string) string {
	if len(value) >= 2 {
		first := value[0]
		last := value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}

// dropDuplicateTopLevelKeys returns a copy of the YAML frontmatter with earlier
// occurrences of any repeated top-level key removed — together with their
// indented continuation lines — keeping the last value of each. Dropping only
// the "key:" line while leaving its block content (e.g. the "  - item" lines of
// a repeated depends_on list) would fold those orphaned lines into the previous
// field's value as a multiline scalar, silently corrupting it. Everything
// between a dropped key line and the next top-level key belongs to the dropped
// key, so it is dropped too. This is a narrow recovery for duplicate-key files,
// not a general YAML rewriter.
func dropDuplicateTopLevelKeys(yamlText string) string {
	lines := strings.Split(yamlText, "\n")

	lastIndexByKey := map[string]int{}
	for lineIndex, line := range lines {
		if key, isKey := topLevelKeyName(line); isKey {
			lastIndexByKey[key] = lineIndex
		}
	}

	keptLines := make([]string, 0, len(lines))
	droppingKeyBlock := false
	for lineIndex, line := range lines {
		if key, isKey := topLevelKeyName(line); isKey {
			droppingKeyBlock = lastIndexByKey[key] != lineIndex
			if droppingKeyBlock {
				continue
			}
		} else if droppingKeyBlock {
			continue // continuation line of a dropped earlier occurrence
		}
		keptLines = append(keptLines, line)
	}
	return strings.Join(keptLines, "\n")
}

// topLevelKeyName reports the key name of an unindented "key: value" line.
// Indented lines, list items, comments, and blanks are not top-level keys.
func topLevelKeyName(line string) (string, bool) {
	if line == "" {
		return "", false
	}
	firstByte := line[0]
	if firstByte == ' ' || firstByte == '\t' || firstByte == '-' || firstByte == '#' {
		return "", false
	}
	colonIndex := strings.Index(line, ":")
	if colonIndex <= 0 {
		return "", false
	}
	return line[:colonIndex], true
}
