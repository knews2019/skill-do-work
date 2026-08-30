package requestmodel

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"
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

func TestTypedRecordCarriesEveryNormalizedSchemaFieldAndGenericEvidence(t *testing.T) {
	fixture := []byte("---\n" +
		"id: REQ-42\n" +
		"domain: back_end\n" +
		"route: b\n" +
		"impact: surprising\n" +
		"effort_estimate: trivial\n" +
		"maintenance: yes\n" +
		"tdd: test_first\n" +
		"error_type: spec\n" +
		"kb_status: skip\n" +
		"testing_status: returned with feedback\n" +
		"builder_decided: true\n" +
		"unknown_field: keep-me\n" +
		"---\nBody\n")
	document, err := ParseDocument(fixture)
	if err != nil {
		t.Fatal(err)
	}
	record := document.TypedRecord()
	resolvedValues := []string{
		record.DomainValue, record.RouteValue, record.ImpactValue,
		record.EffortEstimateValue, record.MaintenanceValue, record.TDDValue,
		record.ErrorTypeValue, record.KBStatusValue, record.TestingStatusValue,
		record.BuilderDecidedValue,
	}
	wantValues := []string{"backend", "B", "impact-user-visible", "effort-mechanical", "true", "true", "spec", "skipped", "returned", "true"}
	if !reflect.DeepEqual(resolvedValues, wantValues) {
		t.Fatalf("normalized values = %v, want %v", resolvedValues, wantValues)
	}
	if record.ImpactEvidence.IsRecognized || record.ImpactEvidence.OriginalValue != "surprising" || record.ImpactEvidence.WarningMessage == "" {
		t.Fatalf("impact evidence = %#v", record.ImpactEvidence)
	}
	if !record.DomainEvidence.IsRecognized || !record.RouteEvidence.IsRecognized || !record.MaintenanceEvidence.IsRecognized || !record.TDDEvidence.IsRecognized || !record.ErrorTypeEvidence.IsRecognized || !record.KBStatusEvidence.IsRecognized || !record.TestingStatusEvidence.IsRecognized || !record.BuilderDecidedEvidence.IsRecognized {
		t.Fatalf("recognized field evidence missing: %#v", record)
	}
	unknownField, found := record.FieldEvidenceByName["unknown_field"]
	if !found || unknownField.ScalarValue != "keep-me" {
		t.Fatalf("generic unknown field evidence = %#v, found=%v", unknownField, found)
	}
}
