package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplitFrontmatter(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		wantHas       bool
		wantYamlHas   string // substring expected in the YAML block
		wantBodyHas   string // substring expected in the body
		wantBodyExact string // when set, the body must equal this exactly
	}{
		{
			name:        "standard frontmatter then body",
			input:       "---\nid: REQ-1\nstatus: completed\n---\n\n# Title\n\nBody text.\n",
			wantHas:     true,
			wantYamlHas: "id: REQ-1",
			wantBodyHas: "# Title",
		},
		{
			name:    "no frontmatter at all",
			input:   "# Just a heading\n\nNo fences here.\n",
			wantHas: false,
		},
		{
			name:    "opening fence but no closing fence",
			input:   "---\nid: REQ-2\nstill going\nand going\n",
			wantHas: false,
		},
		{
			name:          "empty body after closing fence",
			input:         "---\nid: REQ-3\n---",
			wantHas:       true,
			wantYamlHas:   "id: REQ-3",
			wantBodyExact: "",
		},
		{
			name:        "crlf line endings are normalized",
			input:       "---\r\nid: REQ-4\r\n---\r\nbody\r\n",
			wantHas:     true,
			wantYamlHas: "id: REQ-4",
			wantBodyHas: "body",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			yamlText, bodyText, hasFrontmatter := splitFrontmatter(testCase.input)
			if hasFrontmatter != testCase.wantHas {
				t.Fatalf("hasFrontmatter = %v, want %v", hasFrontmatter, testCase.wantHas)
			}
			if !testCase.wantHas {
				if bodyText != testCase.input {
					t.Fatalf("no-frontmatter body should equal original input")
				}
				return
			}
			if testCase.wantYamlHas != "" && !strings.Contains(yamlText, testCase.wantYamlHas) {
				t.Fatalf("yaml %q missing %q", yamlText, testCase.wantYamlHas)
			}
			if testCase.wantBodyHas != "" && !strings.Contains(bodyText, testCase.wantBodyHas) {
				t.Fatalf("body %q missing %q", bodyText, testCase.wantBodyHas)
			}
			if testCase.wantBodyExact != "" || testCase.name == "empty body after closing fence" {
				if bodyText != testCase.wantBodyExact {
					t.Fatalf("body = %q, want exactly %q", bodyText, testCase.wantBodyExact)
				}
			}
		})
	}
}

func TestParseFrontmatterFieldsRecoversDuplicateKeys(t *testing.T) {
	// Mirrors the lone real file with two completed_at lines: the parser must
	// recover and keep the LAST value rather than dropping the whole block.
	yamlText := strings.Join([]string{
		"id: REQ-1034",
		"status: completed",
		"completed_at: 2026-06-10T12:40:16Z",
		"claimed_at: 2026-06-10T12:04:56Z",
		"completed_at: 2026-06-10T14:00:00Z",
	}, "\n")

	fields, parseError := parseFrontmatterFields(yamlText)
	if parseError != nil {
		t.Fatalf("expected duplicate-key recovery, got error: %v", parseError)
	}
	if got := coerceScalarToString(fields["completed_at"]); got != "2026-06-10T14:00:00Z" {
		t.Fatalf("completed_at = %q, want the last value 2026-06-10T14:00:00Z", got)
	}
	if got := coerceScalarToString(fields["status"]); got != "completed" {
		t.Fatalf("status = %q, want completed", got)
	}
}

func TestParseFrontmatterFieldsRecoversDuplicateBlockListKeys(t *testing.T) {
	// A repeated BLOCK-LIST key is the dangerous duplicate shape: dropping only
	// the earlier "depends_on:" line would orphan its "  - item" lines, which
	// YAML then folds into the preceding field (here `id`) as a multiline
	// scalar — silently corrupting the REQ id. The whole earlier block must go.
	yamlText := strings.Join([]string{
		"id: REQ-1",
		"depends_on:",
		"  - REQ-2",
		"  - REQ-3",
		"status: pending",
		"depends_on:",
		"  - REQ-4",
	}, "\n")

	fields, parseError := parseFrontmatterFields(yamlText)
	if parseError != nil {
		t.Fatalf("expected duplicate-key recovery, got error: %v", parseError)
	}
	if got := coerceScalarToString(fields["id"]); got != "REQ-1" {
		t.Fatalf("id = %q, want REQ-1 (earlier block's list items leaked into the preceding field)", got)
	}
	if got := coerceScalarToString(fields["status"]); got != "pending" {
		t.Fatalf("status = %q, want pending", got)
	}
	if got := coerceToStringList(fields["depends_on"]); !reflect.DeepEqual(got, []string{"REQ-4"}) {
		t.Fatalf("depends_on = %v, want the last block's [REQ-4]", got)
	}
}

func TestParseFrontmatterFieldsRecoversMalformedTitleLine(t *testing.T) {
	// Two real malformed-title shapes that strict YAML rejects: a quoted prefix
	// with trailing text, and a bare colon inside the value. Both must still
	// surface status, user_request, and depends_on.
	testCases := []struct {
		name           string
		yamlText       string
		wantStatus     string
		wantUserReq    string
		wantDependsOn  []string
		wantTitleHasIt string
	}{
		{
			name: "quoted prefix then trailing text",
			yamlText: strings.Join([]string{
				`id: REQ-1150`,
				`title: "Clean up broken" button — scan then delete`,
				`status: complete`,
				`commit: 096dacba`,
				`user_request: UR-419`,
				`depends_on: [REQ-1147, REQ-1148]`,
			}, "\n"),
			wantStatus:     "complete",
			wantUserReq:    "UR-419",
			wantDependsOn:  []string{"REQ-1147", "REQ-1148"},
			wantTitleHasIt: "button",
		},
		{
			name: "bare colon inside the title value",
			yamlText: strings.Join([]string{
				`id: REQ-082`,
				`title: Review fix: resolve node version mismatch`,
				`status: completed`,
				`user_request: UR-026`,
			}, "\n"),
			wantStatus:     "completed",
			wantUserReq:    "UR-026",
			wantDependsOn:  nil,
			wantTitleHasIt: "Review fix: resolve",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fields, parseError := parseFrontmatterFields(testCase.yamlText)
			if parseError != nil {
				t.Fatalf("expected lenient recovery, got error: %v", parseError)
			}
			if got := coerceScalarToString(fields["status"]); got != testCase.wantStatus {
				t.Fatalf("status = %q, want %q", got, testCase.wantStatus)
			}
			if got := coerceScalarToString(fields["user_request"]); got != testCase.wantUserReq {
				t.Fatalf("user_request = %q, want %q", got, testCase.wantUserReq)
			}
			if got := coerceToStringList(fields["depends_on"]); !reflect.DeepEqual(got, testCase.wantDependsOn) {
				t.Fatalf("depends_on = %v, want %v", got, testCase.wantDependsOn)
			}
			if got := coerceScalarToString(fields["title"]); !strings.Contains(got, testCase.wantTitleHasIt) {
				t.Fatalf("title %q missing %q", got, testCase.wantTitleHasIt)
			}
		})
	}
}

func TestParseFrontmatterFieldsEmptyIsNotAnError(t *testing.T) {
	fields, parseError := parseFrontmatterFields("")
	if parseError != nil {
		t.Fatalf("empty yaml should not error, got %v", parseError)
	}
	if len(fields) != 0 {
		t.Fatalf("expected no fields, got %d", len(fields))
	}
}

func TestCoerceScalarToString(t *testing.T) {
	testCases := []struct {
		name  string
		value any
		want  string
	}{
		{"nil", nil, ""},
		{"string trimmed", "  hello  ", "hello"},
		{"int", 42, "42"},
		{"float whole", 8.0, "8"},
		{"bool", true, "true"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := coerceScalarToString(testCase.value); got != testCase.want {
				t.Fatalf("coerceScalarToString(%v) = %q, want %q", testCase.value, got, testCase.want)
			}
		})
	}
}

func TestCoerceToStringList(t *testing.T) {
	testCases := []struct {
		name  string
		value any
		want  []string
	}{
		{"nil", nil, nil},
		{"sequence", []any{"REQ-1", "REQ-2"}, []string{"REQ-1", "REQ-2"}},
		{"sequence drops empties", []any{"REQ-1", "", "REQ-3"}, []string{"REQ-1", "REQ-3"}},
		{"bare scalar wraps", "REQ-9", []string{"REQ-9"}},
		{"empty scalar", "   ", nil},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := coerceToStringList(testCase.value)
			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("coerceToStringList(%v) = %v, want %v", testCase.value, got, testCase.want)
			}
		})
	}
}
