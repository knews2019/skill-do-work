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
			// The body must come back VERBATIM — CRLF endings preserved — while
			// the YAML text is normalized for the parser. A normalized body once
			// made the fence-by-subtraction arithmetic steal body bytes into the
			// Copy payload's fence (the CRLF Copy-corruption bug).
			name:          "crlf body stays verbatim, yaml is normalized",
			input:         "---\r\nid: REQ-4\r\n---\r\nline one\r\nline two\r\n",
			wantHas:       true,
			wantYamlHas:   "id: REQ-4",
			wantBodyExact: "line one\r\nline two\r\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			yamlText, bodyText, bodyStartOffset, hasFrontmatter := splitFrontmatter(testCase.input)
			if hasFrontmatter != testCase.wantHas {
				t.Fatalf("hasFrontmatter = %v, want %v", hasFrontmatter, testCase.wantHas)
			}
			if !testCase.wantHas {
				if bodyText != testCase.input {
					t.Fatalf("no-frontmatter body should equal original input")
				}
				return
			}
			if got := testCase.input[:bodyStartOffset] + bodyText; got != testCase.input {
				t.Fatalf("fence prefix + body must reassemble the original file byte-for-byte:\ngot  %q\nwant %q", got, testCase.input)
			}
			if strings.Contains(yamlText, "\r") {
				t.Fatalf("yaml %q must be CRLF-normalized for the parser", yamlText)
			}
			if testCase.wantYamlHas != "" && !strings.Contains(yamlText, testCase.wantYamlHas) {
				t.Fatalf("yaml %q missing %q", yamlText, testCase.wantYamlHas)
			}
			if testCase.wantBodyHas != "" && !strings.Contains(bodyText, testCase.wantBodyHas) {
				t.Fatalf("body %q missing %q", bodyText, testCase.wantBodyHas)
			}
			if testCase.wantBodyExact != "" || testCase.wantBodyHas == "" {
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
		{
			// External-review regression: a malformed title used to make the
			// lenient fallback silently drop a block-style depends_on — the
			// board lost dependency edges with no warning.
			name: "block-list depends_on survives lenient recovery",
			yamlText: strings.Join([]string{
				`id: REQ-1201`,
				`title: "Wire up" the board — with block-list deps`,
				`status: pending`,
				`user_request: UR-450`,
				`depends_on:`,
				`  - REQ-1199`,
				`  - REQ-1200`,
			}, "\n"),
			wantStatus:     "pending",
			wantUserReq:    "UR-450",
			wantDependsOn:  []string{"REQ-1199", "REQ-1200"},
			wantTitleHasIt: "board",
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

// TestLenientFrontmatterFieldsBlockListShapes exercises the block-list
// recovery directly: column-zero items (valid YAML), a blank line between
// items, a quoted item, and a following "key: value" line terminating the
// list without being swallowed by it.
func TestLenientFrontmatterFieldsBlockListShapes(t *testing.T) {
	yamlText := strings.Join([]string{
		`title: "Broken" title forcing lenient recovery`,
		`depends_on:`,
		`- REQ-1`,
		``,
		`- "REQ-2"`,
		`domain: [frontend]`,
		`status: pending`,
	}, "\n")
	fields := lenientFrontmatterFields(yamlText)
	if got := coerceToStringList(fields["depends_on"]); !reflect.DeepEqual(got, []string{"REQ-1", "REQ-2"}) {
		t.Fatalf("depends_on = %v, want [REQ-1 REQ-2]", got)
	}
	if got := coerceToStringList(fields["domain"]); !reflect.DeepEqual(got, []string{"frontend"}) {
		t.Fatalf("domain = %v, want [frontend] (list terminator line must still parse as its own key)", got)
	}
	if got := coerceScalarToString(fields["status"]); got != "pending" {
		t.Fatalf("status = %q, want pending", got)
	}
}

// TestLenientFrontmatterFieldsBareKeyWithoutItemsStaysAbsent pins the
// non-list bare-key behavior: a "key:" line followed by anything other than
// list items (a nested map, another key) contributes no field, exactly as
// before block-list recovery existed.
func TestLenientFrontmatterFieldsBareKeyWithoutItemsStaysAbsent(t *testing.T) {
	yamlText := strings.Join([]string{
		`title: "Broken" title forcing lenient recovery`,
		`metadata:`,
		`  type: user`,
		`tags:`,
		`status: pending`,
	}, "\n")
	fields := lenientFrontmatterFields(yamlText)
	if _, present := fields["metadata"]; present {
		t.Fatalf("metadata = %v, want absent (nested maps are not recovered)", fields["metadata"])
	}
	if _, present := fields["tags"]; present {
		t.Fatalf("tags = %v, want absent (bare key with no items)", fields["tags"])
	}
	if got := coerceScalarToString(fields["status"]); got != "pending" {
		t.Fatalf("status = %q, want pending", got)
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

// userTextFrontmatterFixture wraps one value line in a frontmatter block that
// also carries a NESTED map (estimate:). The nesting is the point: the lenient
// recovery in frontmatter.go is flat and top-level only, so a block that falls
// to it comes back without estimate at all. Asserting estimate survived is how
// these tests tell "the strict parser read this" from "the salvage path did".
func userTextFrontmatterFixture(valueLine string) string {
	return strings.Join([]string{
		`id: REQ-999`,
		`status: pending`,
		valueLine,
		`user_request: UR-068`,
		`estimate:`,
		`  p50_active_minutes: 30`,
		`  confidence: medium`,
	}, "\n")
}

// TestFrontmatterQuotingOfUserText pins the property the schema's Frontmatter
// Quoting contract exists for (skills/do-work/actions/work-reference.md ->
// Request File Schema): a frontmatter value carrying user-typed text comes back
// BYTE-IDENTICAL, through the strict parser, without the rest of the block
// being handed to the last-resort recovery.
//
// The failure it locks out was silent in both directions. A title written the
// way the schema used to prescribe — a bare double-quoted scalar — loses a
// typed double quote: `title: "Fix: A " # B"` is VALID YAML that reads back as
// `Fix: A`, the remainder taken as a comment, with no error to notice. Where it
// fails the strict parse instead, the WHOLE block drops to recovery and the
// nested estimate: map goes with it.
func TestFrontmatterQuotingOfUserText(t *testing.T) {
	// Each case is text a user typed. wantLine is that text written the way the
	// contract requires; the parse must return userText unchanged.
	requiredForm := []struct {
		name     string
		userText string
		wantLine string
	}{
		{
			// The REQ-344 case: a double quote, a colon and a hash in one title.
			name:     "double quote, colon and hash",
			userText: `Fix: A " # B`,
			wantLine: `title: 'Fix: A " # B'`,
		},
		{
			// The common shape — every REQ minted here carries a bracketed
			// impact tag, and the comma is what recovery eats.
			name:     "bracketed impact tag with a comma",
			userText: `[impact-negligible] Retitle export, again [v2]`,
			wantLine: `title: '[impact-negligible] Retitle export, again [v2]'`,
		},
		{
			// The one escape the single-quoted form has: an apostrophe doubles.
			name:     "apostrophe doubled inside the scalar",
			userText: `Don't ship "it": #1 blocker`,
			wantLine: `title: 'Don''t ship "it": #1 blocker'`,
		},
		{
			// blocked_check is a shell command the user supplied verbatim.
			name:     "a supplied shell probe with a hash and quotes",
			userText: `sh -c 'echo "#1" | grep -q 1'`,
			wantLine: `blocked_check: 'sh -c ''echo "#1" | grep -q 1'''`,
		},
	}

	for _, testCase := range requiredForm {
		t.Run("required form/"+testCase.name, func(t *testing.T) {
			fieldKey, _, _ := strings.Cut(testCase.wantLine, ":")
			fields, parseError := parseFrontmatterFields(userTextFrontmatterFixture(testCase.wantLine))
			if parseError != nil {
				t.Fatalf("parse: %v", parseError)
			}
			if got := coerceScalarToString(fields[fieldKey]); got != testCase.userText {
				t.Fatalf("%s round-trip = %q, want %q", fieldKey, got, testCase.userText)
			}
			if fields["estimate"] == nil {
				t.Fatalf("%s: nested estimate: map is gone — the block fell to the last-resort recovery instead of parsing strictly", fieldKey)
			}
		})
	}

	// The forms the contract forbids, kept as the reason it exists. The
	// predicate is that the RECORD does not come back whole — the text is
	// altered, or the block falls to recovery and its nested estimate: map is
	// dropped, or both. Which half breaks varies by shape and is not asserted,
	// and no particular corrupted value is asserted either, so widening the
	// recovery parser (which REQ-344 explicitly does not do) could not turn any
	// of these into a false failure.
	forbiddenForm := []struct {
		name     string
		userText string
		line     string
		why      string
	}{
		{
			name:     "double-quoted scalar, typed quote then hash",
			userText: `Fix: A " # B`,
			line:     `title: "Fix: A " # B"`,
			why:      "valid YAML that silently truncates the title at the typed quote",
		},
		{
			name:     "double-quoted scalar, typed quote mid-value",
			userText: `[impact-negligible] Fix "quoted" key: strip # now`,
			line:     `title: "[impact-negligible] Fix "quoted" key: strip # now"`,
			why:      "strict parse rejects the whole block, so recovery answers and the nested estimate: map is dropped",
		},
		{
			name:     "unquoted scalar with a bracketed tag",
			userText: `[impact-negligible] Retitle export, again [v2]`,
			line:     `title: [impact-negligible] Retitle export, again [v2]`,
			why:      "recovery re-reads it as a flow list and eats the comma",
		},
	}

	for _, testCase := range forbiddenForm {
		t.Run("forbidden form/"+testCase.name, func(t *testing.T) {
			fields, parseError := parseFrontmatterFields(userTextFrontmatterFixture(testCase.line))
			if parseError != nil {
				t.Fatalf("parse: %v", parseError)
			}
			titleSurvived := coerceScalarToString(fields["title"]) == testCase.userText
			if titleSurvived && fields["estimate"] != nil {
				t.Fatalf("the whole record survived %s (%s) — the Frontmatter Quoting contract's premise no longer holds; re-derive the rule before relaxing it", testCase.line, testCase.why)
			}
		})
	}
}

// TestFrontmatterQuotingBlockScalarCarriesANewline pins the contract's second
// form. No quoted scalar on one line can carry a newline, and the recovery
// parser does not recover block scalars at all — so this form exists only while
// the strict parse succeeds, which is what the estimate assertion checks.
func TestFrontmatterQuotingBlockScalarCarriesANewline(t *testing.T) {
	userText := "line one: with a colon\nline two \" with a quote"
	yamlText := strings.Join([]string{
		`id: REQ-999`,
		`status: pending`,
		`blocked_by: |-`,
		`  line one: with a colon`,
		`  line two " with a quote`,
		`estimate:`,
		`  p50_active_minutes: 30`,
	}, "\n")

	fields, parseError := parseFrontmatterFields(yamlText)
	if parseError != nil {
		t.Fatalf("parse: %v", parseError)
	}
	if got := coerceScalarToString(fields["blocked_by"]); got != userText {
		t.Fatalf("blocked_by round-trip = %q, want %q", got, userText)
	}
	if fields["estimate"] == nil {
		t.Fatalf("nested estimate: map is gone — the block fell to the recovery parser, which cannot read block scalars")
	}
}
