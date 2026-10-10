package corehelpers

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const writerFixtureRequest = "---\nid: REQ-701\ntitle: fixture\nstatus: claimed\ndepends_on: [REQ-700]\nestimate:\n  p50_active_minutes: 10\n---\n# Fixture\n\n## What\nDo it.\n"

var canonicalInstantPattern = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z\n$`)

// runRegisteredWriter calls a command through the public registry, so a verb or
// command that was never registered fails the test instead of the build.
func runRegisteredWriter(t *testing.T, root, commandName string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	handler := Handlers()[commandName]
	if handler == nil {
		t.Fatalf("no handler registered for %q", commandName)
	}
	return handler(testContext(root), arguments)
}

func readFixture(t *testing.T, root, relativePath string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func writerOutput(result resultmodel.CommandResult) string {
	if result.ExactTextOutput == nil {
		return ""
	}
	return *result.ExactTextOutput
}

func assertRefusedUnchanged(t *testing.T, result resultmodel.CommandResult, root, relativePath, before, wantCode string) {
	t.Helper()
	if result.Outcome != resultmodel.OutcomeRefused || resultmodel.ExitCode(result.Outcome) != 1 {
		t.Fatalf("outcome=%s findings=%+v, want refused (exit 1)", result.Outcome, result.Findings)
	}
	if len(result.Findings) == 0 || result.Findings[0].Code != wantCode {
		t.Fatalf("findings=%+v, want code %s", result.Findings, wantCode)
	}
	if after := readFixture(t, root, relativePath); after != before {
		t.Fatalf("refusal changed the file:\n%s", after)
	}
}

// Pins the hand-heredoc failure: a second stamp silently overwrote review_at.
func TestFrontmatterSetStampIsAppendOnly(t *testing.T) {
	root := t.TempDir()
	requestPath := "do-work/working/REQ-701-fixture.md"
	writeMatrixFile(t, root, requestPath, writerFixtureRequest)

	first := runRegisteredWriter(t, root, CommandFrontmatter, "set", "REQ-701", "review_at", "--at", "now")
	if first.Outcome != resultmodel.OutcomeSuccess || !canonicalInstantPattern.MatchString(writerOutput(first)) {
		t.Fatalf("first stamp outcome=%s output=%q findings=%+v", first.Outcome, writerOutput(first), first.Findings)
	}
	stamped := readFixture(t, root, requestPath)
	if strings.Count(stamped, "\nreview_at: "+writerOutput(first)) != 1 {
		t.Fatalf("stamp not written once in canonical form:\n%s", stamped)
	}

	second := runRegisteredWriter(t, root, CommandFrontmatter, "set", "REQ-701", "review_at", "--at", "now")
	assertRefusedUnchanged(t, second, root, requestPath, stamped, "FRONTMATTER-STAMP-EXISTS")

	for _, refusal := range []struct {
		name      string
		arguments []string
		wantCode  string
	}{
		{"at on a non-stamp field", []string{"set", "REQ-701", "title", "--at", "now"}, "FRONTMATTER-AT-INVALID"},
		{"at with a word other than now", []string{"set", "REQ-701", "remediation_at", "--at", "later"}, "FRONTMATTER-AT-INVALID"},
		{"explicit stamp with a space", []string{"set", "REQ-701", "remediation_at", "2026-10-10 12:00:00"}, "FRONTMATTER-STAMP-NOT-CANONICAL"},
		{"explicit stamp with an offset", []string{"set", "REQ-701", "remediation_at", "2026-10-10T12:00:00+00:00"}, "FRONTMATTER-STAMP-NOT-CANONICAL"},
	} {
		t.Run(refusal.name, func(t *testing.T) {
			result := runRegisteredWriter(t, root, CommandFrontmatter, refusal.arguments...)
			assertRefusedUnchanged(t, result, root, requestPath, stamped, refusal.wantCode)
		})
	}

	// work.md stamps dispatch_at with an instant held earlier, so a canonical
	// explicit value is the accepted spelling.
	held := runRegisteredWriter(t, root, CommandFrontmatter, "set", "REQ-701", "dispatch_at", "2026-10-10T12:00:00Z")
	if held.Outcome != resultmodel.OutcomeSuccess || writerOutput(held) != "2026-10-10T12:00:00Z\n" {
		t.Fatalf("held canonical stamp outcome=%s output=%q", held.Outcome, writerOutput(held))
	}
}

// Pins the Frontmatter Quoting contract and the collapse of structured fields.
func TestFrontmatterSetWritesQuotedScalars(t *testing.T) {
	root := t.TempDir()
	requestPath := "do-work/working/REQ-701-fixture.md"
	writeMatrixFile(t, root, requestPath, writerFixtureRequest)

	awkward := "it's a: note # here"
	result := runRegisteredWriter(t, root, CommandFrontmatter, "set", requestPath, "review_note", awkward)
	if result.Outcome != resultmodel.OutcomeSuccess || writerOutput(result) != awkward+"\n" {
		t.Fatalf("outcome=%s output=%q findings=%+v", result.Outcome, writerOutput(result), result.Findings)
	}
	if contents := readFixture(t, root, requestPath); !strings.Contains(contents, "  p50_active_minutes: 10\nreview_note: 'it''s a: note # here'\n---\n") {
		t.Fatalf("new field not single-quoted at the end of the block:\n%s", contents)
	}
	readBack := runRegisteredWriter(t, root, CommandFrontmatter, "get", requestPath, "review_note")
	if writerOutput(readBack) != awkward+"\n" {
		t.Fatalf("get read back %q", writerOutput(readBack))
	}

	replaced := runRegisteredWriter(t, root, CommandFrontmatter, "set", "REQ-701", "title", "new title")
	if replaced.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("replace outcome=%s findings=%+v", replaced.Outcome, replaced.Findings)
	}
	if contents := readFixture(t, root, requestPath); !strings.Contains(contents, "\ntitle: 'new title'\nstatus: claimed\n") || strings.Count(contents, "\ntitle:") != 1 {
		t.Fatalf("existing scalar not replaced in place:\n%s", contents)
	}

	before := readFixture(t, root, requestPath)
	for _, field := range []string{"depends_on", "estimate"} {
		t.Run(field, func(t *testing.T) {
			result := runRegisteredWriter(t, root, CommandFrontmatter, "set", "REQ-701", field, "flat")
			assertRefusedUnchanged(t, result, root, requestPath, before, "FRONTMATTER-FIELD-STRUCTURED")
		})
	}
}

// Pins review F3: set wrote status: bogus-status, which looks validated but bypasses the lifecycle owners.
func TestFrontmatterSetRefusesLifecycleOwnedFields(t *testing.T) {
	root := t.TempDir()
	requestPath := "do-work/working/REQ-701-fixture.md"
	writeMatrixFile(t, root, requestPath, writerFixtureRequest)
	for _, owned := range []struct {
		field, value, wantOwner string
	}{
		{"status", "bogus-status", "advance"},
		{"id", "REQ-999", "never rewritten"},
	} {
		t.Run(owned.field, func(t *testing.T) {
			result := runRegisteredWriter(t, root, CommandFrontmatter, "set", "REQ-701", owned.field, owned.value)
			assertRefusedUnchanged(t, result, root, requestPath, writerFixtureRequest, "FRONTMATTER-FIELD-OWNED")
			if evidence := strings.Join(result.Findings[0].Evidence, " "); !strings.Contains(evidence, owned.wantOwner) {
				t.Fatalf("evidence %q does not name the owner %q", evidence, owned.wantOwner)
			}
		})
	}
}

// Pins the out-of-order section write that advance later refuses.
func TestRequestAppendSectionLandsInCanonicalOrder(t *testing.T) {
	cases := []struct {
		name, body, want string
	}{
		{
			name: "before a later canonical section",
			body: "## What\nDo it.\n\n## Review\nLooks good.\n",
			want: "## What\nDo it.\n\n## Testing\n\ntests ran\n\n## Review\nLooks good.\n",
		},
		{
			name: "at the end when no later canonical section exists",
			body: "## What\nDo it.\n\n## Why\nBecause.\n",
			want: "## What\nDo it.\n\n## Why\nBecause.\n\n## Testing\n\ntests ran\n",
		},
		{
			name: "a fenced heading is not an anchor",
			body: "## What\n```\n## Review\n```\n",
			want: "## What\n```\n## Review\n```\n\n## Testing\n\ntests ran\n",
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			requestPath := "do-work/working/REQ-701-fixture.md"
			frontmatter := "---\nid: REQ-701\nstatus: claimed\n---\n"
			writeMatrixFile(t, root, requestPath, frontmatter+test.body)
			writeMatrixFile(t, root, "testing.md", "\ntests ran\n\n")
			result := runRegisteredWriter(t, root, CommandRequest, "append-section", "REQ-701", "--section", "Testing", "--from", "testing.md")
			if result.Outcome != resultmodel.OutcomeSuccess {
				t.Fatalf("outcome=%s findings=%+v", result.Outcome, result.Findings)
			}
			if got := readFixture(t, root, requestPath); got != frontmatter+test.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, frontmatter+test.want)
			}
		})
	}

	t.Run("a name outside the canonical order refuses", func(t *testing.T) {
		root := t.TempDir()
		requestPath := "do-work/working/REQ-701-fixture.md"
		writeMatrixFile(t, root, requestPath, writerFixtureRequest)
		writeMatrixFile(t, root, "notes.md", "free text\n")
		result := runRegisteredWriter(t, root, CommandRequest, "append-section", "REQ-701", "--section", "Notes", "--from", "notes.md")
		assertRefusedUnchanged(t, result, root, requestPath, writerFixtureRequest, "SECTION-NOT-CANONICAL")
	})
}

// Pins the double append: a retried hand write left two Testing sections.
func TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses(t *testing.T) {
	root := t.TempDir()
	requestPath := "do-work/working/REQ-701-fixture.md"
	writeMatrixFile(t, root, requestPath, writerFixtureRequest)
	writeMatrixFile(t, root, "testing.md", "tests ran\n")
	writeMatrixFile(t, root, "other.md", "different evidence\n")
	writeMatrixFile(t, root, "headed.md", "## Review\n\nstray\n")
	arguments := []string{"append-section", "REQ-701", "--section", "Testing", "--from", "testing.md"}

	if first := runRegisteredWriter(t, root, CommandRequest, arguments...); first.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("first outcome=%s findings=%+v", first.Outcome, first.Findings)
	}
	written := readFixture(t, root, requestPath)
	second := runRegisteredWriter(t, root, CommandRequest, arguments...)
	if second.Outcome != resultmodel.OutcomeSuccess || resultmodel.ExitCode(second.Outcome) != 0 {
		t.Fatalf("repeat outcome=%s findings=%+v", second.Outcome, second.Findings)
	}
	if after := readFixture(t, root, requestPath); after != written || strings.Count(after, "## Testing") != 1 {
		t.Fatalf("repeat changed the file:\n%s", after)
	}

	conflict := runRegisteredWriter(t, root, CommandRequest, "append-section", "REQ-701", "--section", "Testing", "--from", "other.md")
	assertRefusedUnchanged(t, conflict, root, requestPath, written, "SECTION-CONFLICT")

	// Review F1: a body carrying its own ## heading wrote a stray section that
	// advance later refuses.
	headed := runRegisteredWriter(t, root, CommandRequest, "append-section", "REQ-701", "--section", "Qualification", "--from", "headed.md")
	assertRefusedUnchanged(t, headed, root, requestPath, written, "SECTION-BODY-HAS-HEADING")
}

// Pins review N1: an open fence or comment in the body hid the later ## Review from advance.
func TestRequestAppendSectionRefusesBodyThatHidesLaterSections(t *testing.T) {
	frontmatter := "---\nid: REQ-701\nstatus: claimed\n---\n"
	withReview := frontmatter + "## What\nDo it.\n\n## Review\nLooks good.\n"
	for _, test := range []struct {
		name, request, from string
	}{
		{"an open fence in the body", withReview, "```\n## x\n"},
		{"an open comment in the body", withReview, "text <!-- open\n"},
		// The end-of-file case the old single-copy re-check guarded: the
		// request's own open fence hides the appended heading.
		{"an open fence already at the end of the file", frontmatter + "## What\n```\nstill fenced\n", "tests ran\n"},
		// A net section count is fooled here: the open fence hides ## Review
		// and reveals the fenced example ## Hidden, so the count grows by one.
		{"an open fence that reveals a fenced example heading", withReview + "\n```\n## Hidden\n```\n", "```\nopen\n"},
		// Appended at the end, nothing follows yet, but the open fence hides
		// the ## Review written after it.
		{"an open fence in a body appended at the end", frontmatter + "## What\nDo it.\n", "```\nopen\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			requestPath := "do-work/working/REQ-701-fixture.md"
			writeMatrixFile(t, root, requestPath, test.request)
			writeMatrixFile(t, root, "testing.md", test.from)
			result := runRegisteredWriter(t, root, CommandRequest, "append-section", "REQ-701", "--section", "Testing", "--from", "testing.md")
			assertRefusedUnchanged(t, result, root, requestPath, test.request, "SECTION-WRITE-FAILED")
		})
	}
}

// Pins the resolver: an id must name exactly one working or queue file.
func TestRequestIDResolvesOneWorkingOrQueueFile(t *testing.T) {
	root := t.TempDir()
	writeMatrixFile(t, root, "do-work/working/REQ-701-fixture.md", writerFixtureRequest)
	writeMatrixFile(t, root, "do-work/working/REQ-702-twice.md", "---\nid: REQ-702\nstatus: claimed\n---\n")
	writeMatrixFile(t, root, "do-work/queue/REQ-702-twice.md", "---\nid: REQ-702\nstatus: pending\n---\n")
	writeMatrixFile(t, root, "do-work/archive/REQ-703-done.md", "---\nid: REQ-703\nstatus: completed\n---\n")
	writeMatrixFile(t, root, "testing.md", "tests ran\n")

	for _, refusal := range []struct {
		target, path, wantCode string
	}{
		{"REQ-704", "", "REQUEST-NOT-FOUND"},
		{"REQ-702", "do-work/queue/REQ-702-twice.md", "REQUEST-AMBIGUOUS"},
		{"REQ-703", "do-work/archive/REQ-703-done.md", "REQUEST-NOT-ACTIVE"},
		{"do-work/archive/REQ-703-done.md", "do-work/archive/REQ-703-done.md", "REQUEST-NOT-ACTIVE"},
	} {
		t.Run(refusal.target, func(t *testing.T) {
			before := ""
			if refusal.path != "" {
				before = readFixture(t, root, refusal.path)
			}
			check := func(result resultmodel.CommandResult) {
				t.Helper()
				if refusal.path != "" {
					assertRefusedUnchanged(t, result, root, refusal.path, before, refusal.wantCode)
				} else if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) == 0 || result.Findings[0].Code != refusal.wantCode {
					t.Fatalf("outcome=%s findings=%+v, want %s", result.Outcome, result.Findings, refusal.wantCode)
				}
			}
			check(runRegisteredWriter(t, root, CommandFrontmatter, "set", refusal.target, "review_at", "--at", "now"))
			if strings.HasPrefix(refusal.target, "REQ-") {
				check(runRegisteredWriter(t, root, CommandRequest, "append-section", refusal.target, "--section", "Testing", "--from", "testing.md"))
			}
		})
	}

	if result := runRegisteredWriter(t, root, CommandRequest, "append-section", "REQ-701", "--section", "Testing", "--from", "testing.md"); result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("working REQ did not resolve: outcome=%s findings=%+v", result.Outcome, result.Findings)
	}
}
