// Package requestmodel parses and edits REQ/UR Markdown documents without
// rewriting unrelated frontmatter or body bytes.
package requestmodel

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

var (
	ErrFrontmatterMissing  = errors.New("frontmatter is missing")
	ErrFrontmatterUnclosed = errors.New("frontmatter is not closed")
)

// FieldEvidence retains the selected top-level field and its decoded shapes.
type FieldEvidence struct {
	FieldName      string
	LineNumber     int
	RawValue       string
	ScalarValue    string
	ListValues     []string
	NestedValues   map[string]string
	DuplicateCount int
}

// RequestRecord is the typed REQ/UR evidence consumed by repository snapshots.
type RequestRecord struct {
	RequestID                           string
	RequestTitle                        string
	RequestStatus                       string
	OriginalStatus                      string
	StatusEvidence                      schemanormalization.FieldResult
	UserRequestID                       string
	DependsOn                           []string
	DependencySource                    string
	AddendumTo                          string
	AddendumSource                      string
	RelatedIDs                          []string
	RelatedSource                       string
	BatchName                           string
	BatchSource                         string
	SuggestedSpec                       string
	SuggestedSpecSource                 string
	WritePaths                          []string
	AssignedTo                          string
	CreatedAt                           string
	ClaimedAt                           string
	CompletedAt                         string
	ImplementationCommit                string
	HeavyVerifiedAt                     string
	HeavyVerifiedRevision               string
	CavemanValue                        string
	CavemanEvidence                     schemanormalization.FieldResult
	DomainValue                         string
	DomainEvidence                      schemanormalization.FieldResult
	RouteValue                          string
	RouteEvidence                       schemanormalization.FieldResult
	ImpactValue                         string
	ImpactEvidence                      schemanormalization.FieldResult
	RequestPriorityValue                string
	RequestPriorityEvidence             schemanormalization.FieldResult
	EffortEstimateValue                 string
	EffortEstimateEvidence              schemanormalization.FieldResult
	MaintenanceValue                    string
	MaintenanceEvidence                 schemanormalization.FieldResult
	TDDValue                            string
	TDDEvidence                         schemanormalization.FieldResult
	ErrorTypeValue                      string
	ErrorTypeEvidence                   schemanormalization.FieldResult
	KBStatusValue                       string
	KBStatusEvidence                    schemanormalization.FieldResult
	TestingStatusValue                  string
	TestingStatusEvidence               schemanormalization.FieldResult
	BuilderDecidedValue                 string
	BuilderDecidedEvidence              schemanormalization.FieldResult
	GateDeferredValue                   string
	GateDeferredEvidence                schemanormalization.FieldResult
	RepositoryGateRepairValue           string
	RepositoryGateRepairEvidence        schemanormalization.FieldResult
	DeferredImplementationBaseValue     string
	DeferredImplementationBaseEvidence  schemanormalization.FieldResult
	DeferredImplementationMergeValue    string
	DeferredImplementationMergeEvidence schemanormalization.FieldResult
	FieldEvidenceByName                 map[string]FieldEvidence
}

// RequestDocument is one lossless frontmatter-plus-Markdown document.
type RequestDocument struct {
	dataBytes          []byte
	bodyStartOffset    int
	closingFenceOffset int
	lineEnding         string
	warnings           []string
	fieldSpans         map[string][]fieldSpan
}

type fieldSpan struct {
	evidence   FieldEvidence
	lineStart  int
	lineEnd    int
	valueStart int
	valueEnd   int
	blockEnd   int
}

type documentLine struct {
	start      int
	contentEnd int
	end        int
}

// ParseDocument parses one complete REQ or UR document.
func ParseDocument(fileBytes []byte) (*RequestDocument, error) {
	document := &RequestDocument{dataBytes: append([]byte(nil), fileBytes...), fieldSpans: map[string][]fieldSpan{}}
	openingOffset := 0
	if bytes.HasPrefix(fileBytes, []byte{0xef, 0xbb, 0xbf}) {
		openingOffset = 3
	}
	openingEnd, openingLineEnding, foundOpening := fenceLineEnd(fileBytes, openingOffset)
	if !foundOpening {
		return nil, ErrFrontmatterMissing
	}
	document.lineEnding = openingLineEnding

	var frontmatterLines []documentLine
	lineStart := openingEnd
	for lineStart < len(fileBytes) {
		line := scanDocumentLine(fileBytes, lineStart)
		lineText := fileBytes[line.start:line.contentEnd]
		if bytes.Equal(lineText, []byte("---")) {
			document.closingFenceOffset = line.start
			document.bodyStartOffset = line.end
			break
		}
		frontmatterLines = append(frontmatterLines, line)
		lineStart = line.end
	}
	if document.bodyStartOffset == 0 {
		return nil, ErrFrontmatterUnclosed
	}
	if document.lineEnding == "" {
		document.lineEnding = "\n"
	}

	type keyLine struct {
		name       string
		colonIndex int
		lineIndex  int
	}
	var keyLines []keyLine
	for lineIndex, line := range frontmatterLines {
		lineBytes := fileBytes[line.start:line.contentEnd]
		if len(bytes.TrimSpace(lineBytes)) == 0 || lineBytes[0] == '#' || lineBytes[0] == ' ' || lineBytes[0] == '\t' || lineBytes[0] == '-' {
			continue
		}
		colonIndex := bytes.IndexByte(lineBytes, ':')
		if colonIndex <= 0 {
			document.warnings = append(document.warnings, fmt.Sprintf("malformed top-level frontmatter line %d retained", lineIndex+2))
			continue
		}
		fieldName := strings.TrimSpace(string(lineBytes[:colonIndex]))
		if fieldName == "" {
			continue
		}
		keyLines = append(keyLines, keyLine{name: fieldName, colonIndex: colonIndex, lineIndex: lineIndex})
	}
	for keyIndex, key := range keyLines {
		line := frontmatterLines[key.lineIndex]
		nextKeyStart := document.closingFenceOffset
		if keyIndex+1 < len(keyLines) {
			nextKeyStart = frontmatterLines[keyLines[keyIndex+1].lineIndex].start
		}
		valueStart := line.start + key.colonIndex + 1
		for valueStart < line.contentEnd && (fileBytes[valueStart] == ' ' || fileBytes[valueStart] == '\t') {
			valueStart++
		}
		valueEnd := scalarValueEnd(fileBytes, valueStart, line.contentEnd)
		rawValue := string(fileBytes[valueStart:valueEnd])
		blockEnd := line.end
		trimmedRawValue := strings.TrimSpace(rawValue)
		isLiteralBlock := strings.HasPrefix(trimmedRawValue, "|") || strings.HasPrefix(trimmedRawValue, ">")
		if rawValue == "" || isLiteralBlock {
			blockEnd = continuationBlockEnd(fileBytes, line.end, nextKeyStart, isLiteralBlock)
		}
		evidence := FieldEvidence{FieldName: key.name, LineNumber: key.lineIndex + 2, RawValue: rawValue}
		if rawValue != "" {
			if strings.HasPrefix(trimmedRawValue, "|") {
				evidence.ScalarValue = decodeLiteralBlock(trimmedRawValue, fileBytes[line.end:blockEnd])
			} else {
				evidence.ScalarValue = decodeScalar(rawValue)
			}
			if strings.HasPrefix(trimmedRawValue, "[") && strings.HasSuffix(trimmedRawValue, "]") {
				evidence.ListValues = decodeFlowList(rawValue)
			}
		} else {
			evidence.ListValues, evidence.NestedValues = decodeIndentedBlock(fileBytes[line.end:blockEnd])
		}
		span := fieldSpan{
			evidence: evidence, lineStart: line.start, lineEnd: line.end,
			valueStart: valueStart, valueEnd: valueEnd, blockEnd: blockEnd,
		}
		document.fieldSpans[key.name] = append(document.fieldSpans[key.name], span)
	}
	for fieldName, spans := range document.fieldSpans {
		if len(spans) <= 1 {
			continue
		}
		document.warnings = append(document.warnings, fmt.Sprintf("duplicate top-level key %q; last value wins", fieldName))
		for spanIndex := range spans {
			spans[spanIndex].evidence.DuplicateCount = len(spans)
		}
		document.fieldSpans[fieldName] = spans
	}
	return document, nil
}

// DocumentBytes returns the complete current document.
func (document *RequestDocument) DocumentBytes() []byte {
	return append([]byte(nil), document.dataBytes...)
}

// BodyBytes returns every byte after the closing frontmatter fence.
func (document *RequestDocument) BodyBytes() []byte {
	return append([]byte(nil), document.dataBytes[document.bodyStartOffset:]...)
}

// ReplaceBodySpan replaces one half-open span relative to BodyBytes. It reparses
// the document so later frontmatter edits retain correct absolute offsets.
func (document *RequestDocument) ReplaceBodySpan(startOffset, endOffset int, replacementBytes []byte) error {
	bodyLength := len(document.dataBytes) - document.bodyStartOffset
	if startOffset < 0 || endOffset < startOffset || endOffset > bodyLength {
		return fmt.Errorf("body span [%d,%d) is outside [0,%d)", startOffset, endOffset, bodyLength)
	}
	absoluteStart := document.bodyStartOffset + startOffset
	absoluteEnd := document.bodyStartOffset + endOffset
	return document.reparse(replaceByteSpan(document.dataBytes, absoluteStart, absoluteEnd, replacementBytes))
}

// ParseWarnings returns non-fatal recovery evidence such as duplicate keys.
func (document *RequestDocument) ParseWarnings() []string {
	return append([]string(nil), document.warnings...)
}

// FieldValue returns the last top-level occurrence of a field.
func (document *RequestDocument) FieldValue(fieldName string) (FieldEvidence, bool) {
	spans := document.fieldSpans[fieldName]
	if len(spans) == 0 {
		return FieldEvidence{}, false
	}
	evidence := spans[len(spans)-1].evidence
	if evidence.ListValues != nil {
		evidence.ListValues = append([]string{}, evidence.ListValues...)
	}
	if evidence.NestedValues != nil {
		evidence.NestedValues = cloneStringMap(evidence.NestedValues)
	}
	return evidence, true
}

// TypedRecord projects the parsed document into typed request evidence.
func (document *RequestDocument) TypedRecord() RequestRecord {
	status := document.scalarValue("status")
	statusEvidence := schemanormalization.NormalizeField("status", status)
	cavemanEvidence := schemanormalization.NormalizeField("caveman", document.scalarValue("caveman"))
	domainEvidence := schemanormalization.NormalizeField("domain", document.scalarValue("domain"))
	routeEvidence := schemanormalization.NormalizeField("route", document.scalarValue("route"))
	impactEvidence := schemanormalization.NormalizeField("impact", document.scalarValue("impact"))
	requestPriorityEvidence := schemanormalization.NormalizeField("priority", document.scalarValue("priority"))
	effortEstimateEvidence := schemanormalization.NormalizeField("effort_estimate", document.scalarValue("effort_estimate"))
	maintenanceEvidence := schemanormalization.NormalizeField("maintenance", document.scalarValue("maintenance"))
	tddEvidence := schemanormalization.NormalizeField("tdd", document.scalarValue("tdd"))
	errorTypeEvidence := schemanormalization.NormalizeField("error_type", document.scalarValue("error_type"))
	kbStatusEvidence := schemanormalization.NormalizeField("kb_status", document.scalarValue("kb_status"))
	testingStatusEvidence := schemanormalization.NormalizeField("testing_status", document.scalarValue("testing_status"))
	builderDecidedEvidence := schemanormalization.NormalizeField("builder_decided", document.scalarValue("builder_decided"))
	gateDeferredEvidence := schemanormalization.NormalizeField("gate_deferred", document.scalarValue("gate_deferred"))
	repositoryGateRepairEvidence := schemanormalization.NormalizeField("repository_gate_repair", document.scalarValue("repository_gate_repair"))
	deferredImplementationBaseEvidence := schemanormalization.NormalizeField("deferred_implementation_base", document.scalarValue("deferred_implementation_base"))
	deferredImplementationMergeEvidence := schemanormalization.NormalizeField("deferred_implementation_merge", document.scalarValue("deferred_implementation_merge"))
	dependencyValues, dependencySource := document.preferredList("depends_on")
	addendumValues, addendumSource := document.preferredList("addendum_to")
	relatedValues, relatedSource := document.preferredList("related")
	batchValue, batchSource := document.preferredScalar("batch")
	suggestedSpecValue, suggestedSpecSource := document.preferredScalar("suggested_spec")
	addendumValue := ""
	if len(addendumValues) > 0 {
		addendumValue = addendumValues[0]
	}
	fieldEvidenceByName := make(map[string]FieldEvidence, len(document.fieldSpans))
	for fieldName := range document.fieldSpans {
		if fieldEvidence, found := document.FieldValue(fieldName); found {
			fieldEvidenceByName[fieldName] = fieldEvidence
		}
	}
	return RequestRecord{
		RequestID: document.scalarValue("id"), RequestTitle: document.scalarValue("title"),
		RequestStatus: statusEvidence.ResolvedValue, OriginalStatus: status, StatusEvidence: statusEvidence,
		UserRequestID: document.scalarValue("user_request"),
		DependsOn:     dependencyValues, DependencySource: dependencySource,
		AddendumTo: addendumValue, AddendumSource: addendumSource,
		RelatedIDs: relatedValues, RelatedSource: relatedSource,
		BatchName: batchValue, BatchSource: batchSource,
		SuggestedSpec: suggestedSpecValue, SuggestedSpecSource: suggestedSpecSource,
		WritePaths: document.listValue("write_set"),
		AssignedTo: strings.TrimSpace(document.scalarValue("assigned_to")),
		CreatedAt:  document.scalarValue("created_at"), ClaimedAt: document.scalarValue("claimed_at"),
		CompletedAt: document.scalarValue("completed_at"), ImplementationCommit: strings.TrimSpace(document.scalarValue("commit")),
		HeavyVerifiedAt: document.scalarValue("heavy_verified_at"), HeavyVerifiedRevision: strings.TrimSpace(document.scalarValue("heavy_verified_revision")),
		CavemanValue: cavemanEvidence.ResolvedValue, CavemanEvidence: cavemanEvidence,
		DomainValue: domainEvidence.ResolvedValue, DomainEvidence: domainEvidence,
		RouteValue: routeEvidence.ResolvedValue, RouteEvidence: routeEvidence,
		ImpactValue: impactEvidence.ResolvedValue, ImpactEvidence: impactEvidence,
		RequestPriorityValue: requestPriorityEvidence.ResolvedValue, RequestPriorityEvidence: requestPriorityEvidence,
		EffortEstimateValue: effortEstimateEvidence.ResolvedValue, EffortEstimateEvidence: effortEstimateEvidence,
		MaintenanceValue: maintenanceEvidence.ResolvedValue, MaintenanceEvidence: maintenanceEvidence,
		TDDValue: tddEvidence.ResolvedValue, TDDEvidence: tddEvidence,
		ErrorTypeValue: errorTypeEvidence.ResolvedValue, ErrorTypeEvidence: errorTypeEvidence,
		KBStatusValue: kbStatusEvidence.ResolvedValue, KBStatusEvidence: kbStatusEvidence,
		TestingStatusValue: testingStatusEvidence.ResolvedValue, TestingStatusEvidence: testingStatusEvidence,
		BuilderDecidedValue: builderDecidedEvidence.ResolvedValue, BuilderDecidedEvidence: builderDecidedEvidence,
		GateDeferredValue: gateDeferredEvidence.ResolvedValue, GateDeferredEvidence: gateDeferredEvidence,
		RepositoryGateRepairValue: repositoryGateRepairEvidence.ResolvedValue, RepositoryGateRepairEvidence: repositoryGateRepairEvidence,
		DeferredImplementationBaseValue: deferredImplementationBaseEvidence.ResolvedValue, DeferredImplementationBaseEvidence: deferredImplementationBaseEvidence,
		DeferredImplementationMergeValue: deferredImplementationMergeEvidence.ResolvedValue, DeferredImplementationMergeEvidence: deferredImplementationMergeEvidence,
		FieldEvidenceByName: fieldEvidenceByName,
	}
}

// SetList replaces or appends one top-level field using the canonical inline-list form.
func (document *RequestDocument) SetList(fieldName string, values []string) error {
	if !validFieldName(fieldName) {
		return fmt.Errorf("invalid frontmatter field name %q", fieldName)
	}
	encodedItems := make([]string, len(values))
	for valueIndex, fieldValue := range values {
		if strings.ContainsAny(fieldValue, "\x00\r\n") {
			return fmt.Errorf("frontmatter list %q contains an unsupported control character", fieldName)
		}
		encodedItems[valueIndex] = encodeScalar(fieldValue, document.lineEnding)
	}
	encodedValue := "[" + strings.Join(encodedItems, ", ") + "]"
	updatedBytes := document.dataBytes
	spans := document.fieldSpans[fieldName]
	if len(spans) == 0 {
		fieldLine := []byte(fieldName + ": " + encodedValue + document.lineEnding)
		updatedBytes = replaceByteSpan(updatedBytes, document.closingFenceOffset, document.closingFenceOffset, fieldLine)
	} else {
		span := spans[len(spans)-1]
		if span.evidence.RawValue == "" {
			fieldText := fieldName + ": " + encodedValue + document.lineEnding
			updatedBytes = replaceByteSpan(updatedBytes, span.lineStart, span.blockEnd, []byte(fieldText))
		} else {
			updatedBytes = replaceByteSpan(updatedBytes, span.valueStart, span.valueEnd, []byte(encodedValue))
		}
	}
	return document.reparse(updatedBytes)
}

// SetScalar replaces or appends one top-level scalar field.
func (document *RequestDocument) SetScalar(fieldName string, fieldValue string) error {
	if !validFieldName(fieldName) {
		return fmt.Errorf("invalid frontmatter field name %q", fieldName)
	}
	if strings.ContainsAny(fieldValue, "\x00\r") {
		return fmt.Errorf("frontmatter scalar %q contains an unsupported control character", fieldName)
	}
	encodedValue := encodeScalar(fieldValue, document.lineEnding)
	updatedBytes := document.dataBytes
	spans := document.fieldSpans[fieldName]
	if len(spans) == 0 {
		fieldLine := []byte(fieldName + ": " + encodedValue + document.lineEnding)
		updatedBytes = replaceByteSpan(updatedBytes, document.closingFenceOffset, document.closingFenceOffset, fieldLine)
	} else {
		span := spans[len(spans)-1]
		if strings.Contains(fieldValue, "\n") || span.evidence.RawValue == "" {
			fieldText := fieldName + ": " + encodedValue + document.lineEnding
			updatedBytes = replaceByteSpan(updatedBytes, span.lineStart, span.blockEnd, []byte(fieldText))
		} else {
			updatedBytes = replaceByteSpan(updatedBytes, span.valueStart, span.valueEnd, []byte(encodedValue))
		}
	}
	return document.reparse(updatedBytes)
}

// DeleteField removes every byte belonging to a top-level field.
func (document *RequestDocument) DeleteField(fieldName string) error {
	spans := document.fieldSpans[fieldName]
	updatedBytes := document.dataBytes
	for spanIndex := len(spans) - 1; spanIndex >= 0; spanIndex-- {
		span := spans[spanIndex]
		updatedBytes = replaceByteSpan(updatedBytes, span.lineStart, span.blockEnd, nil)
	}
	return document.reparse(updatedBytes)
}

// ParseTimestamp accepts current and legacy timestamp shapes.
func ParseTimestamp(timestampText string) (time.Time, error) {
	trimmedTimestamp := strings.TrimSpace(timestampText)
	for _, timestampLayout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	} {
		if parsedTimestamp, parseError := time.Parse(timestampLayout, trimmedTimestamp); parseError == nil {
			return parsedTimestamp, nil
		}
	}
	return time.Time{}, fmt.Errorf("timestamp %q is not a supported REQ/UR timestamp", timestampText)
}

// CanonicalTimestamp renders the schema's UTC whole-second shape.
func CanonicalTimestamp(timestamp time.Time) string {
	return timestamp.UTC().Truncate(time.Second).Format("2006-01-02T15:04:05Z")
}

func fenceLineEnd(fileBytes []byte, openingOffset int) (int, string, bool) {
	if bytes.HasPrefix(fileBytes[openingOffset:], []byte("---\r\n")) {
		return openingOffset + 5, "\r\n", true
	}
	if bytes.HasPrefix(fileBytes[openingOffset:], []byte("---\n")) {
		return openingOffset + 4, "\n", true
	}
	return 0, "", false
}

func scanDocumentLine(fileBytes []byte, lineStart int) documentLine {
	newlineRelative := bytes.IndexByte(fileBytes[lineStart:], '\n')
	if newlineRelative < 0 {
		return documentLine{start: lineStart, contentEnd: len(fileBytes), end: len(fileBytes)}
	}
	newlineOffset := lineStart + newlineRelative
	contentEnd := newlineOffset
	if contentEnd > lineStart && fileBytes[contentEnd-1] == '\r' {
		contentEnd--
	}
	return documentLine{start: lineStart, contentEnd: contentEnd, end: newlineOffset + 1}
}

func scalarValueEnd(fileBytes []byte, valueStart int, lineEnd int) int {
	singleQuoted := false
	doubleQuoted := false
	escaped := false
	commentStart := lineEnd
	for byteIndex := valueStart; byteIndex < lineEnd; byteIndex++ {
		currentByte := fileBytes[byteIndex]
		if doubleQuoted && currentByte == '\\' && !escaped {
			escaped = true
			continue
		}
		if currentByte == '\'' && !doubleQuoted {
			if singleQuoted && byteIndex+1 < lineEnd && fileBytes[byteIndex+1] == '\'' {
				byteIndex++
				continue
			}
			singleQuoted = !singleQuoted
		} else if currentByte == '"' && !singleQuoted && !escaped {
			doubleQuoted = !doubleQuoted
		} else if currentByte == '#' && !singleQuoted && !doubleQuoted && (byteIndex == valueStart || unicode.IsSpace(rune(fileBytes[byteIndex-1]))) {
			commentStart = byteIndex
			break
		}
		escaped = false
	}
	for commentStart > valueStart && (fileBytes[commentStart-1] == ' ' || fileBytes[commentStart-1] == '\t') {
		commentStart--
	}
	return commentStart
}

func decodeScalar(rawValue string) string {
	trimmedValue := strings.TrimSpace(rawValue)
	if len(trimmedValue) >= 2 && trimmedValue[0] == '\'' && trimmedValue[len(trimmedValue)-1] == '\'' {
		return strings.ReplaceAll(trimmedValue[1:len(trimmedValue)-1], "''", "'")
	}
	if len(trimmedValue) >= 2 && trimmedValue[0] == '"' && trimmedValue[len(trimmedValue)-1] == '"' {
		if decodedValue, decodeError := strconv.Unquote(trimmedValue); decodeError == nil {
			return decodedValue
		}
	}
	return trimmedValue
}

func decodeFlowList(rawValue string) []string {
	trimmedValue := strings.TrimSpace(rawValue)
	trimmedValue = strings.TrimSpace(trimmedValue[1 : len(trimmedValue)-1])
	if trimmedValue == "" {
		return []string{}
	}
	var values []string
	itemStart := 0
	singleQuoted := false
	doubleQuoted := false
	for byteIndex := 0; byteIndex < len(trimmedValue); byteIndex++ {
		switch trimmedValue[byteIndex] {
		case '\'':
			if !doubleQuoted {
				singleQuoted = !singleQuoted
			}
		case '"':
			if !singleQuoted {
				doubleQuoted = !doubleQuoted
			}
		case ',':
			if !singleQuoted && !doubleQuoted {
				values = append(values, decodeScalar(trimmedValue[itemStart:byteIndex]))
				itemStart = byteIndex + 1
			}
		}
	}
	values = append(values, decodeScalar(trimmedValue[itemStart:]))
	return values
}

func decodeIndentedBlock(blockBytes []byte) ([]string, map[string]string) {
	var values []string
	nestedValues := map[string]string{}
	for lineStart := 0; lineStart < len(blockBytes); {
		line := scanDocumentLine(blockBytes, lineStart)
		lineText := strings.TrimSpace(string(blockBytes[line.start:line.contentEnd]))
		lineStart = line.end
		if lineText == "" || strings.HasPrefix(lineText, "#") {
			continue
		}
		if itemText, isItem := strings.CutPrefix(lineText, "- "); isItem {
			values = append(values, decodeLineScalar(itemText))
			continue
		}
		if colonIndex := strings.IndexByte(lineText, ':'); colonIndex > 0 {
			nestedValues[strings.TrimSpace(lineText[:colonIndex])] = decodeLineScalar(lineText[colonIndex+1:])
		}
	}
	if len(values) == 0 {
		values = nil
	}
	if len(nestedValues) == 0 {
		nestedValues = nil
	}
	return values, nestedValues
}

func continuationBlockEnd(fileBytes []byte, blockStart int, maximumEnd int, includeTrailingBlanks bool) int {
	blockEnd := blockStart
	for lineStart := blockStart; lineStart < maximumEnd; {
		line := scanDocumentLine(fileBytes, lineStart)
		if line.end > maximumEnd {
			line.end = maximumEnd
		}
		lineBytes := fileBytes[line.start:line.contentEnd]
		lineStart = line.end
		if len(bytes.TrimSpace(lineBytes)) == 0 {
			if includeTrailingBlanks {
				blockEnd = line.end
			}
			continue
		}
		if len(lineBytes) > 0 && (lineBytes[0] == ' ' || lineBytes[0] == '\t' || lineBytes[0] == '-') {
			blockEnd = line.end
			continue
		}
		break
	}
	return blockEnd
}

func decodeLiteralBlock(indicator string, blockBytes []byte) string {
	var rawLines []string
	minimumIndent := -1
	for lineStart := 0; lineStart < len(blockBytes); {
		line := scanDocumentLine(blockBytes, lineStart)
		lineText := string(blockBytes[line.start:line.contentEnd])
		lineStart = line.end
		if strings.TrimSpace(lineText) == "" {
			rawLines = append(rawLines, "")
			continue
		}
		rawLines = append(rawLines, lineText)
		indent := len(lineText) - len(strings.TrimLeft(lineText, " \t"))
		if minimumIndent < 0 || indent < minimumIndent {
			minimumIndent = indent
		}
	}
	if minimumIndent < 0 {
		minimumIndent = 0
	}
	for lineIndex, lineText := range rawLines {
		if len(lineText) >= minimumIndent {
			rawLines[lineIndex] = lineText[minimumIndent:]
		}
	}
	decodedValue := strings.Join(rawLines, "\n")
	if len(blockBytes) > 0 && blockBytes[len(blockBytes)-1] == '\n' {
		decodedValue += "\n"
	}
	switch {
	case strings.HasPrefix(indicator, "|-"):
		return strings.TrimRight(decodedValue, "\n")
	case strings.HasPrefix(indicator, "|+"):
		return decodedValue
	default:
		return strings.TrimRight(decodedValue, "\n") + "\n"
	}
}

func decodeLineScalar(rawValue string) string {
	valueBytes := []byte(strings.TrimSpace(rawValue))
	valueEnd := scalarValueEnd(valueBytes, 0, len(valueBytes))
	return decodeScalar(string(valueBytes[:valueEnd]))
}

func (document *RequestDocument) scalarValue(fieldName string) string {
	field, found := document.FieldValue(fieldName)
	if !found {
		return ""
	}
	return field.ScalarValue
}

func (document *RequestDocument) listValue(fieldName string) []string {
	field, found := document.FieldValue(fieldName)
	if !found {
		return nil
	}
	if field.ListValues != nil {
		return field.ListValues
	}
	if field.ScalarValue == "" {
		return nil
	}
	return []string{field.ScalarValue}
}

func (document *RequestDocument) preferredList(canonicalKey string) ([]string, string) {
	if _, found := document.fieldSpans[canonicalKey]; found {
		return document.listValue(canonicalKey), canonicalKey
	}
	for _, aliasKey := range schemanormalization.SchemaFieldAliases(canonicalKey) {
		if _, found := document.fieldSpans[aliasKey]; found {
			return document.listValue(aliasKey), aliasKey
		}
	}
	return nil, ""
}

func (document *RequestDocument) preferredScalar(canonicalKey string) (string, string) {
	if _, found := document.fieldSpans[canonicalKey]; found {
		return document.scalarValue(canonicalKey), canonicalKey
	}
	for _, aliasKey := range schemanormalization.SchemaFieldAliases(canonicalKey) {
		if _, found := document.fieldSpans[aliasKey]; found {
			return document.scalarValue(aliasKey), aliasKey
		}
	}
	return "", ""
}

func encodeScalar(fieldValue string, lineEnding string) string {
	if strings.Contains(fieldValue, "\n") {
		terminalNewlines := len(fieldValue) - len(strings.TrimRight(fieldValue, "\n"))
		indicator := "|-"
		if terminalNewlines == 1 {
			indicator = "|"
		} else if terminalNewlines > 1 {
			indicator = "|+"
		}
		trimmedValue := strings.TrimRight(fieldValue, "\n")
		encodedValue := indicator + lineEnding + "  " + strings.ReplaceAll(trimmedValue, "\n", lineEnding+"  ")
		if terminalNewlines > 1 {
			encodedValue += strings.Repeat(lineEnding, terminalNewlines-1)
		}
		return encodedValue
	}
	if fieldValue != "" && isPlainScalar(fieldValue) {
		return fieldValue
	}
	return "'" + strings.ReplaceAll(fieldValue, "'", "''") + "'"
}

func isPlainScalar(fieldValue string) bool {
	for _, character := range fieldValue {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || strings.ContainsRune("-_.:/", character) {
			continue
		}
		return false
	}
	switch strings.ToLower(fieldValue) {
	case "null", "true", "false", "yes", "no", "on", "off", "~":
		return false
	}
	return true
}

func validFieldName(fieldName string) bool {
	if fieldName == "" {
		return false
	}
	for _, character := range fieldName {
		if unicode.IsLetter(character) || unicode.IsDigit(character) || character == '_' || character == '-' {
			continue
		}
		return false
	}
	return true
}

func replaceByteSpan(sourceBytes []byte, startOffset int, endOffset int, replacementBytes []byte) []byte {
	updatedBytes := make([]byte, 0, len(sourceBytes)-(endOffset-startOffset)+len(replacementBytes))
	updatedBytes = append(updatedBytes, sourceBytes[:startOffset]...)
	updatedBytes = append(updatedBytes, replacementBytes...)
	updatedBytes = append(updatedBytes, sourceBytes[endOffset:]...)
	return updatedBytes
}

func (document *RequestDocument) reparse(updatedBytes []byte) error {
	updatedDocument, parseError := ParseDocument(updatedBytes)
	if parseError != nil {
		return parseError
	}
	*document = *updatedDocument
	return nil
}

func cloneStringMap(source map[string]string) map[string]string {
	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}
