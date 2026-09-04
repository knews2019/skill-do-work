package requestmodel

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

func TestParseDocumentRetainsBytesAndRecoversLiveShapes(t *testing.T) {
	fixture := []byte("\xef\xbb\xbf---\r\n" +
		"id: REQ-007\r\n" +
		"title: Review fix: keep colon # malformed YAML but recoverable\r\n" +
		"status: pending\r\n" +
		"dependencies:\r\n  - REQ-001\r\n  - 'REQ-002'\r\n" +
		"depends_on: [REQ-003, \"REQ-004\"]\r\n" +
		"estimate:\r\n  p50_active_minutes: 25\r\n  confidence: low\r\n" +
		"unknown_field: ' keep me '\r\n" +
		"id: REQ-008\r\n" +
		"---\r\n\r\n# Body\r\n\r\nKeep \xff bytes.\r\n")
	document, err := ParseDocument(fixture)
	if err != nil {
		t.Fatalf("ParseDocument: %v", err)
	}
	if !bytes.Equal(document.DocumentBytes(), fixture) {
		t.Fatal("parse changed original bytes")
	}
	if !bytes.Equal(document.BodyBytes(), []byte("\r\n# Body\r\n\r\nKeep \xff bytes.\r\n")) {
		t.Fatalf("body = %q", document.BodyBytes())
	}
	record := document.TypedRecord()
	if record.RequestID != "REQ-008" || record.RequestTitle != "Review fix: keep colon" || record.RequestStatus != "pending" {
		t.Fatalf("record identity = %#v", record)
	}
	if !reflect.DeepEqual(record.DependsOn, []string{"REQ-003", "REQ-004"}) || record.DependencySource != "depends_on" {
		t.Fatalf("dependency evidence = %v from %q", record.DependsOn, record.DependencySource)
	}
	estimate, found := document.FieldValue("estimate")
	if !found || !reflect.DeepEqual(estimate.NestedValues, map[string]string{"p50_active_minutes": "25", "confidence": "low"}) {
		t.Fatalf("estimate = %#v", estimate)
	}
	idField, found := document.FieldValue("id")
	if !found || idField.DuplicateCount != 2 || len(document.ParseWarnings()) == 0 {
		t.Fatalf("duplicate evidence = %#v warnings=%v", idField, document.ParseWarnings())
	}
}

func TestSetScalarChangesOnlyAuthorizedBytes(t *testing.T) {
	original := []byte("---\r\nid: REQ-9\r\nstatus: pending   # pipeline state\r\nunknown: [a, b]\r\n---\r\nBODY\r\n")
	document, err := ParseDocument(original)
	if err != nil {
		t.Fatal(err)
	}
	if err := document.SetScalar("status", "completed-with-issues"); err != nil {
		t.Fatalf("SetScalar: %v", err)
	}
	want := bytes.Replace(original, []byte("status: pending"), []byte("status: completed-with-issues"), 1)
	if !bytes.Equal(document.DocumentBytes(), want) {
		t.Fatalf("edited bytes:\n%q\nwant:\n%q", document.DocumentBytes(), want)
	}
	if err := document.SetScalar("assigned_to", "Cloud: O'Brien"); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(document.DocumentBytes(), []byte("assigned_to: 'Cloud: O''Brien'\r\n---")) {
		t.Fatalf("appended quoted scalar missing: %q", document.DocumentBytes())
	}
	if err := document.DeleteField("assigned_to"); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(document.DocumentBytes(), want) {
		t.Fatal("DeleteField did not restore all unrelated bytes")
	}
}

func TestSetListCanonicalizesOwnedFieldAndPreservesOutsideBytes(t *testing.T) {
	original := []byte("---\r\nid: REQ-9\r\ndepends_on:\r\n  - REQ-1\r\nunknown: [keep, me]\r\n---\r\nBODY\r\n")
	document, err := ParseDocument(original)
	if err != nil {
		t.Fatal(err)
	}
	if err := document.SetList("depends_on", []string{"REQ-1", "REQ-2"}); err != nil {
		t.Fatal(err)
	}
	want := []byte("---\r\nid: REQ-9\r\ndepends_on: [REQ-1, REQ-2]\r\nunknown: [keep, me]\r\n---\r\nBODY\r\n")
	if !bytes.Equal(document.DocumentBytes(), want) {
		t.Fatalf("SetList bytes = %q, want %q", document.DocumentBytes(), want)
	}
	if err := document.SetList("related", []string{}); err != nil {
		t.Fatal(err)
	}
	if field, found := document.FieldValue("related"); !found || field.ListValues == nil || len(field.ListValues) != 0 {
		t.Fatalf("empty list evidence = %#v, found=%t", field, found)
	}
}

func TestRequiredLessonsSurviveNormalizedWriterRoundTrip(t *testing.T) {
	original := []byte("---\n" +
		"id: REQ-478\n" +
		"status: pending\n" +
		"required_lessons: [skills/do-work/actions/lessons-actions.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#final-boundary-identity]\n" +
		"---\nBody\n")
	document, err := ParseDocument(original)
	if err != nil {
		t.Fatal(err)
	}

	// TypedRecord exercises schemanormalization before the requestmodel writer
	// performs an ordinary state transition.
	_ = document.TypedRecord()
	if err := document.SetScalar("status", "claimed"); err != nil {
		t.Fatal(err)
	}

	wantLessons := []byte("required_lessons: [skills/do-work/actions/lessons-actions.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#final-boundary-identity]\n")
	if bytes.Count(document.DocumentBytes(), wantLessons) != 1 {
		t.Fatalf("required_lessons changed during normalized write: %q", document.DocumentBytes())
	}
	want := bytes.Replace(original, []byte("status: pending"), []byte("status: claimed"), 1)
	if !bytes.Equal(document.DocumentBytes(), want) {
		t.Fatalf("normalized write changed unrelated bytes:\n%q\nwant:\n%q", document.DocumentBytes(), want)
	}
}

func TestLiteralBlockScalarAndUnrelatedCommentsSurviveFieldEdits(t *testing.T) {
	original := []byte("---\nblocked_by: |-\n  first line\n  second: line\n# keep this comment\nstatus: pending\n---\nBody\n")
	document, err := ParseDocument(original)
	if err != nil {
		t.Fatal(err)
	}
	blockedBy, found := document.FieldValue("blocked_by")
	if !found || blockedBy.ScalarValue != "first line\nsecond: line" {
		t.Fatalf("blocked_by = %#v", blockedBy)
	}
	if err := document.DeleteField("blocked_by"); err != nil {
		t.Fatal(err)
	}
	want := []byte("---\n# keep this comment\nstatus: pending\n---\nBody\n")
	if !bytes.Equal(document.DocumentBytes(), want) {
		t.Fatalf("DeleteField bytes = %q, want %q", document.DocumentBytes(), want)
	}
}

func TestSetScalarRoundTripsMultilineChomping(t *testing.T) {
	document, err := ParseDocument([]byte("---\r\nid: REQ-1\r\n---\r\nBody\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fieldValue := range []string{"first\nsecond", "one newline\n", "two newlines\n\n", "\n"} {
		if err := document.SetScalar("blocked_by", fieldValue); err != nil {
			t.Fatal(err)
		}
		field, found := document.FieldValue("blocked_by")
		if !found || field.ScalarValue != fieldValue {
			t.Fatalf("multiline scalar round trip = %q, want %q; bytes=%q", field.ScalarValue, fieldValue, document.DocumentBytes())
		}
	}
}

func TestParseDocumentRejectsMissingAndUnclosedFences(t *testing.T) {
	if _, err := ParseDocument([]byte("# body\n")); !errors.Is(err, ErrFrontmatterMissing) {
		t.Fatalf("missing fence error = %v", err)
	}
	if _, err := ParseDocument([]byte("---\nid: REQ-1\n")); !errors.Is(err, ErrFrontmatterUnclosed) {
		t.Fatalf("unclosed fence error = %v", err)
	}
}

func TestTimestampCompatibilityAndCanonicalWriting(t *testing.T) {
	for _, timestampText := range []string{"2026-08-30T20:21:22Z", "2026-08-30T23:21:22+03:00", "2026-08-30T20:21:22", "2026-08-30 20:21:22", "2026-08-30"} {
		if _, err := ParseTimestamp(timestampText); err != nil {
			t.Errorf("ParseTimestamp(%q): %v", timestampText, err)
		}
	}
	localTime := time.Date(2026, 8, 30, 23, 21, 22, 999_000_000, time.FixedZone("EEST", 3*60*60))
	if got := CanonicalTimestamp(localTime); got != "2026-08-30T20:21:22Z" {
		t.Fatalf("CanonicalTimestamp = %q", got)
	}
}

func TestReplaceBodySpanPreservesEveryOutsideByte(t *testing.T) {
	original := []byte("\xef\xbb\xbf---\r\nid: REQ-9\r\n---\r\nhead \xff\r\nquestion\r\ntail\r\n")
	document, err := ParseDocument(original)
	if err != nil {
		t.Fatal(err)
	}
	body := document.BodyBytes()
	start := bytes.Index(body, []byte("question"))
	if err := document.ReplaceBodySpan(start, start+len("question"), []byte("answer")); err != nil {
		t.Fatal(err)
	}
	want := bytes.Replace(original, []byte("question"), []byte("answer"), 1)
	if !bytes.Equal(document.DocumentBytes(), want) {
		t.Fatalf("edited bytes = %q, want %q", document.DocumentBytes(), want)
	}
	if err := document.ReplaceBodySpan(len(document.BodyBytes()), len(document.BodyBytes()), []byte("appended")); err != nil {
		t.Fatal(err)
	}
	if !bytes.HasSuffix(document.DocumentBytes(), []byte("appended")) {
		t.Fatal("append-at-end was not applied")
	}
	before := document.DocumentBytes()
	if err := document.ReplaceBodySpan(-1, 0, nil); err == nil {
		t.Fatal("negative span accepted")
	}
	if err := document.ReplaceBodySpan(2, 1, nil); err == nil {
		t.Fatal("reversed span accepted")
	}
	if err := document.ReplaceBodySpan(0, len(document.BodyBytes())+1, nil); err == nil {
		t.Fatal("out-of-range span accepted")
	}
	if !bytes.Equal(document.DocumentBytes(), before) {
		t.Fatal("refused span changed bytes")
	}
}

func TestTypedRecordCarriesEveryNormalizedSchemaFieldAndGenericEvidence(t *testing.T) {
	fixture := []byte("---\n" +
		"id: REQ-42\n" +
		"status: done\n" +
		"commit:  abc123  \n" +
		"builder_decided: true\n" +
		"gate_deferred: yes\n" +
		"repository_gate_repair: true\n" +
		"deferred_implementation_base: abc123\n" +
		"deferred_implementation_merge: def456\n" +
		"caveman: light\n" +
		"domain: back_end\n" +
		"effort_estimate: trivial\n" +
		"error_type: spec\n" +
		"impact: surprising\n" +
		"priority: urgent\n" +
		"kb_status: skip\n" +
		"maintenance: yes\n" +
		"route: b\n" +
		"tdd: test_first\n" +
		"testing_status: returned with feedback\n" +
		"unknown_field: keep-me\n" +
		"---\nBody\n")
	document, err := ParseDocument(fixture)
	if err != nil {
		t.Fatal(err)
	}
	record := document.TypedRecord()
	if record.ImplementationCommit != "abc123" {
		t.Fatalf("implementation commit = %q", record.ImplementationCommit)
	}
	projections := []struct {
		schemaField   string
		valueField    string
		evidenceField string
		rawValue      string
		resolvedValue string
		recognized    bool
	}{
		{"builder_decided", "BuilderDecidedValue", "BuilderDecidedEvidence", "true", "true", true},
		{"caveman", "CavemanValue", "CavemanEvidence", "light", "lite", true},
		{"deferred_implementation_base", "DeferredImplementationBaseValue", "DeferredImplementationBaseEvidence", "abc123", "abc123", true},
		{"deferred_implementation_merge", "DeferredImplementationMergeValue", "DeferredImplementationMergeEvidence", "def456", "def456", true},
		{"domain", "DomainValue", "DomainEvidence", "back_end", "backend", true},
		{"effort_estimate", "EffortEstimateValue", "EffortEstimateEvidence", "trivial", "effort-mechanical", true},
		{"error_type", "ErrorTypeValue", "ErrorTypeEvidence", "spec", "spec", true},
		{"gate_deferred", "GateDeferredValue", "GateDeferredEvidence", "yes", "true", true},
		{"impact", "ImpactValue", "ImpactEvidence", "surprising", "impact-user-visible", false},
		{"kb_status", "KBStatusValue", "KBStatusEvidence", "skip", "skipped", true},
		{"maintenance", "MaintenanceValue", "MaintenanceEvidence", "yes", "true", true},
		{"priority", "RequestPriorityValue", "RequestPriorityEvidence", "urgent", "next", false},
		{"repository_gate_repair", "RepositoryGateRepairValue", "RepositoryGateRepairEvidence", "true", "true", true},
		{"route", "RouteValue", "RouteEvidence", "b", "B", true},
		{"status", "RequestStatus", "StatusEvidence", "done", "completed", true},
		{"tdd", "TDDValue", "TDDEvidence", "test_first", "true", true},
		{"testing_status", "TestingStatusValue", "TestingStatusEvidence", "returned with feedback", "returned", true},
	}
	projectedFieldNames := make([]string, 0, len(projections))
	recordValue := reflect.ValueOf(record)
	for _, projection := range projections {
		projectedFieldNames = append(projectedFieldNames, projection.schemaField)
		valueField := recordValue.FieldByName(projection.valueField)
		if !valueField.IsValid() {
			t.Errorf("%s value field %s is missing", projection.schemaField, projection.valueField)
			continue
		}
		if valueField.String() != projection.resolvedValue {
			t.Errorf("%s normalized value = %q, want %q", projection.schemaField, valueField.String(), projection.resolvedValue)
		}
		evidenceField := recordValue.FieldByName(projection.evidenceField)
		if !evidenceField.IsValid() {
			t.Errorf("%s evidence field %s is missing", projection.schemaField, projection.evidenceField)
			continue
		}
		evidence, ok := evidenceField.Interface().(schemanormalization.FieldResult)
		if !ok {
			t.Errorf("%s evidence field %s has type %s, want schemanormalization.FieldResult", projection.schemaField, projection.evidenceField, evidenceField.Type())
			continue
		}
		if evidence.FieldName != projection.schemaField || evidence.OriginalValue != projection.rawValue || evidence.ResolvedValue != projection.resolvedValue || evidence.IsRecognized != projection.recognized {
			t.Errorf("%s evidence = %#v", projection.schemaField, evidence)
		}
		genericEvidence, found := record.FieldEvidenceByName[projection.schemaField]
		if !found || genericEvidence.ScalarValue != projection.rawValue {
			t.Errorf("%s generic evidence = %#v, found=%v", projection.schemaField, genericEvidence, found)
		}
	}
	if !reflect.DeepEqual(projectedFieldNames, schemanormalization.SchemaFieldNames()) {
		t.Fatalf("projected schema fields = %v, want every contracted field %v", projectedFieldNames, schemanormalization.SchemaFieldNames())
	}
	if record.ImpactEvidence.WarningMessage == "" {
		t.Fatalf("impact evidence = %#v", record.ImpactEvidence)
	}
	unknownField, found := record.FieldEvidenceByName["unknown_field"]
	if !found || unknownField.ScalarValue != "keep-me" {
		t.Fatalf("generic unknown field evidence = %#v, found=%v", unknownField, found)
	}
}

func TestEffectiveFieldEvidenceCarriesExactSourceLine(t *testing.T) {
	document, err := ParseDocument([]byte("---\ncreated_at: 2026-08-01 # shadowed\n  calculated_at: 2099-01-01\ncreated_at: '2026-08-02' # effective\n---\nBody\n"))
	if err != nil {
		t.Fatal(err)
	}
	field, found := document.FieldValue("created_at")
	if !found || field.LineNumber != 4 || field.RawValue != "'2026-08-02'" || field.DuplicateCount != 2 {
		t.Fatalf("effective source evidence = %#v, found=%v", field, found)
	}
}

func TestTypedRecordProjectsAliasBackedCaptureFieldsWithCanonicalPrecedence(t *testing.T) {
	document, err := ParseDocument([]byte("---\n" +
		"id: REQ-200\n" +
		"dependencies: [REQ-001]\n" +
		"parent: REQ-002\n" +
		"related_reqs: [REQ-003, REQ-004]\n" +
		"batch_name: legacy-batch\n" +
		"spec_hint: legacy-spec\n" +
		"related: [REQ-005]\n" +
		"batch: canonical-batch\n" +
		"suggested_spec: canonical-spec\n" +
		"---\n"))
	if err != nil {
		t.Fatal(err)
	}
	record := document.TypedRecord()
	if !reflect.DeepEqual(record.DependsOn, []string{"REQ-001"}) || record.DependencySource != "dependencies" {
		t.Fatalf("dependencies = %v from %q", record.DependsOn, record.DependencySource)
	}
	if record.AddendumTo != "REQ-002" || record.AddendumSource != "parent" {
		t.Fatalf("addendum = %q from %q", record.AddendumTo, record.AddendumSource)
	}
	if !reflect.DeepEqual(record.RelatedIDs, []string{"REQ-005"}) || record.RelatedSource != "related" {
		t.Fatalf("related = %v from %q", record.RelatedIDs, record.RelatedSource)
	}
	if record.BatchName != "canonical-batch" || record.BatchSource != "batch" {
		t.Fatalf("batch = %q from %q", record.BatchName, record.BatchSource)
	}
	if record.SuggestedSpec != "canonical-spec" || record.SuggestedSpecSource != "suggested_spec" {
		t.Fatalf("suggested spec = %q from %q", record.SuggestedSpec, record.SuggestedSpecSource)
	}
}

func TestTypedRecordSelectsEachReadAliasAndEmptyCanonicalKey(t *testing.T) {
	for _, canonicalKey := range []string{"depends_on", "addendum_to", "related", "batch", "suggested_spec"} {
		for _, aliasKey := range schemanormalization.SchemaFieldAliases(canonicalKey) {
			for _, canonicalPresent := range []bool{false, true} {
				frontmatter := "---\nid: REQ-200\n" + aliasKey + ": REQ-001\n"
				if canonicalPresent {
					frontmatter += canonicalKey + ":\n"
				}
				document, err := ParseDocument([]byte(frontmatter + "---\n"))
				if err != nil {
					t.Fatal(err)
				}
				record := document.TypedRecord()
				var value, source string
				switch canonicalKey {
				case "depends_on":
					value, source = strings.Join(record.DependsOn, ","), record.DependencySource
				case "addendum_to":
					value, source = record.AddendumTo, record.AddendumSource
				case "related":
					value, source = strings.Join(record.RelatedIDs, ","), record.RelatedSource
				case "batch":
					value, source = record.BatchName, record.BatchSource
				case "suggested_spec":
					value, source = record.SuggestedSpec, record.SuggestedSpecSource
				}
				wantValue, wantSource := "REQ-001", aliasKey
				if canonicalPresent {
					wantValue, wantSource = "", canonicalKey
				}
				if value != wantValue || source != wantSource {
					t.Errorf("%s canonical-present=%v: %q from %q, want %q from %q", aliasKey, canonicalPresent, value, source, wantValue, wantSource)
				}
			}
		}
	}
}

func TestTypedRecordPreservesEmptyInlineListEvidence(t *testing.T) {

	tests := []struct {
		name             string
		frontmatter      string
		wantDependencies []string
		wantListPresent  bool
	}{
		{"absent", "", nil, false},
		{"scalar", "depends_on: REQ-101\n", []string{"REQ-101"}, false},
		{"populated flow list", "depends_on: [REQ-101, REQ-102]\n", []string{"REQ-101", "REQ-102"}, true},
		{"empty flow list", "depends_on: []\n", []string{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document, err := ParseDocument([]byte("---\nid: REQ-100\nstatus: pending\n" + test.frontmatter + "---\nBody\n"))
			if err != nil {
				t.Fatal(err)
			}
			record := document.TypedRecord()
			if !reflect.DeepEqual(record.DependsOn, test.wantDependencies) {
				t.Fatalf("DependsOn = %#v, want %#v", record.DependsOn, test.wantDependencies)
			}
			field, found := document.FieldValue("depends_on")
			if found != (test.frontmatter != "") {
				t.Fatalf("FieldValue found = %v", found)
			}
			if found && (field.ListValues != nil) != test.wantListPresent {
				t.Fatalf("ListValues presence = %v, want %v", field.ListValues != nil, test.wantListPresent)
			}
		})
	}
}
