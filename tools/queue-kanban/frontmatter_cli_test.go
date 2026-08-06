package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFrontmatterFixture seeds a REQ file with the given frontmatter body and
// returns its path.
func writeFrontmatterFixture(t *testing.T, frontmatterBody string) string {
	t.Helper()
	fixturePath := filepath.Join(t.TempDir(), "REQ-777-fixture.md")
	fixtureContent := "---\n" + frontmatterBody + "\n---\n\n# Body\n\nProse.\n"
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	return fixturePath
}

func TestReadFrontmatterFieldReturnsRawValue(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\nstatus: completed-with-issues\ndomain: back-end")

	value, found, readError := readFrontmatterField(fixturePath, "status")
	if readError != nil {
		t.Fatalf("readFrontmatterField: %v", readError)
	}
	if !found {
		t.Fatal("found = false, want true — status is present in the fixture")
	}
	if value != "completed-with-issues" {
		t.Fatalf("value = %q, want %q", value, "completed-with-issues")
	}
}

func TestReadFrontmatterFieldReportsAbsentField(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\nstatus: pending")

	value, found, readError := readFrontmatterField(fixturePath, "domain")
	if readError != nil {
		t.Fatalf("readFrontmatterField: %v", readError)
	}
	if found {
		t.Fatalf("found = true for an absent field, want false (value was %q)", value)
	}
}

// A file with no frontmatter at all is an error, not an absent field — the
// caller has to be able to tell "this REQ does not set domain" from "this is not
// a REQ", because the first is routine and the second is a broken input.
func TestReadFrontmatterFieldErrorsOnFileWithoutFrontmatter(t *testing.T) {
	fixturePath := filepath.Join(t.TempDir(), "not-a-req.md")
	if writeError := os.WriteFile(fixturePath, []byte("# Just a heading\n\nNo fences here.\n"), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	if _, _, readError := readFrontmatterField(fixturePath, "status"); readError == nil {
		t.Fatal("readFrontmatterField error = nil, want an error for a file with no frontmatter block")
	}
}

// The lenient recovery path is the whole reason prose cannot replicate this: a
// duplicate top-level key makes strict YAML fail, and the shipped parser keeps
// the last value rather than dropping the file.
func TestReadFrontmatterFieldRecoversDuplicateKey(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\ncompleted_at: 2026-01-01T00:00:00Z\ncompleted_at: 2026-02-02T00:00:00Z")

	value, found, readError := readFrontmatterField(fixturePath, "completed_at")
	if readError != nil {
		t.Fatalf("readFrontmatterField: %v", readError)
	}
	if !found || value != "2026-02-02T00:00:00Z" {
		t.Fatalf("value = %q (found=%v), want the LAST duplicate %q", value, found, "2026-02-02T00:00:00Z")
	}
}

func TestRunFrontmatterCommandGetPrintsValue(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\nstatus: completed\ndomain: back-end")

	var standardOut, standardErr strings.Builder
	exitCode := runFrontmatterCommand([]string{"get", fixturePath, "status"}, &standardOut, &standardErr)

	if exitCode != 0 {
		t.Fatalf("exit = %d, want 0 (stderr: %q)", exitCode, standardErr.String())
	}
	if standardOut.String() != "completed\n" {
		t.Fatalf("stdout = %q, want %q", standardOut.String(), "completed\n")
	}
}

// --normalize must send the canonical value to stdout and the contract's warning
// to stderr, so a caller can capture stdout as a clean single value.
func TestRunFrontmatterCommandNormalizeSeparatesValueFromWarning(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\ndomain: back-end")

	var standardOut, standardErr strings.Builder
	exitCode := runFrontmatterCommand([]string{"get", fixturePath, "domain", "--normalize"}, &standardOut, &standardErr)

	if exitCode != 0 {
		t.Fatalf("exit = %d, want 0", exitCode)
	}
	if standardOut.String() != "backend\n" {
		t.Fatalf("stdout = %q, want %q", standardOut.String(), "backend\n")
	}
	if standardErr.String() != "" {
		t.Fatalf("stderr = %q, want empty — back-end is a recognized alias, not a violation", standardErr.String())
	}
}

func TestRunFrontmatterCommandNormalizeWarnsOnUnrecognizedValue(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\ndomain: quantum")

	var standardOut, standardErr strings.Builder
	exitCode := runFrontmatterCommand([]string{"get", fixturePath, "domain", "--normalize"}, &standardOut, &standardErr)

	if exitCode != 0 {
		t.Fatalf("exit = %d, want 0 — an unrecognized value still resolves to the default", exitCode)
	}
	if standardOut.String() != "general\n" {
		t.Fatalf("stdout = %q, want the documented default %q", standardOut.String(), "general\n")
	}
	warningText := standardErr.String()
	for _, expectedFragment := range []string{"domain", "quantum", "not recognized", "general"} {
		if !strings.Contains(warningText, expectedFragment) {
			t.Fatalf("stderr %q missing %q — must match the contract's warning shape", warningText, expectedFragment)
		}
	}
}

// TestRunFrontmatterCommandNormalizeIsSilentOnFieldOutsideTheContract is REQ-118's
// stated RED case. isKnownSchemaFieldValue answers false both for a genuinely bad
// value and for a field the contract does not govern at all, so the warn branch
// fired on every --normalize of a timestamp, a title, or a path list — telling the
// reader the value was "not recognized" when actions/work-reference.md's Schema
// Read Contract explicitly places such fields OUTSIDE it, to be read verbatim
// ("no alias map, no case folding, no path canonicalization, no warning").
func TestRunFrontmatterCommandNormalizeIsSilentOnFieldOutsideTheContract(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t,
		"id: REQ-777\nstatus: completed\ncreated_at: 2026-08-05T15:53:39Z\ntitle: A plain title\nwrite_set: [a.go, b.go]")

	for _, fieldName := range []string{"created_at", "title", "id"} {
		t.Run(fieldName, func(t *testing.T) {
			var standardOut, standardErr strings.Builder
			exitCode := runFrontmatterCommand([]string{"get", fixturePath, fieldName, "--normalize"}, &standardOut, &standardErr)

			if exitCode != 0 {
				t.Fatalf("exit = %d, want 0", exitCode)
			}
			if standardOut.String() == "" {
				t.Fatalf("stdout is empty — the value must still print")
			}
			if standardErr.String() != "" {
				t.Fatalf("stderr = %q, want empty — %s has no canonical vocabulary, so the contract does not govern it and --normalize is a no-op",
					standardErr.String(), fieldName)
			}
		})
	}
}

// TestRunFrontmatterCommandNormalizeMatchesPlainGetOutsideTheContract states the
// no-op property directly: for a field the contract does not govern, --normalize
// must change nothing a caller can observe.
func TestRunFrontmatterCommandNormalizeMatchesPlainGetOutsideTheContract(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\ncreated_at: 2026-08-05T15:53:39Z")

	var plainOut, plainErr strings.Builder
	plainExit := runFrontmatterCommand([]string{"get", fixturePath, "created_at"}, &plainOut, &plainErr)
	var normalizedOut, normalizedErr strings.Builder
	normalizedExit := runFrontmatterCommand([]string{"get", fixturePath, "created_at", "--normalize"}, &normalizedOut, &normalizedErr)

	if plainExit != normalizedExit || plainOut.String() != normalizedOut.String() || plainErr.String() != normalizedErr.String() {
		t.Fatalf("--normalize changed observable output for a contract-less field:\n  plain: exit=%d out=%q err=%q\n  normalized: exit=%d out=%q err=%q",
			plainExit, plainOut.String(), plainErr.String(),
			normalizedExit, normalizedOut.String(), normalizedErr.String())
	}
}

// TestRunFrontmatterCommandInSetRejectsFieldOutsideTheContract covers the other
// flag. Both set names are `status` sets, so a membership test on a field with no
// canonical vocabulary is a question the command cannot answer — a usage error,
// not a silent "no".
func TestRunFrontmatterCommandInSetRejectsFieldOutsideTheContract(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\ncreated_at: 2026-08-05T15:53:39Z")

	var standardOut, standardErr strings.Builder
	exitCode := runFrontmatterCommand(
		[]string{"get", fixturePath, "created_at", "--in-set", "terminal-success"}, &standardOut, &standardErr)

	if exitCode != 2 {
		t.Fatalf("exit = %d, want 2 — a membership test on a field with no canonical vocabulary is a usage error, not a false answer", exitCode)
	}
	if standardOut.String() != "" {
		t.Fatalf("stdout = %q, want empty on a usage error", standardOut.String())
	}
	if !strings.Contains(standardErr.String(), "created_at") {
		t.Fatalf("stderr = %q, want it to name the offending field", standardErr.String())
	}
}

// --in-set is the membership check ~35 prose sites perform by hand. It prints
// nothing and answers via the exit code.
func TestRunFrontmatterCommandInSetExitCodes(t *testing.T) {
	testCases := []struct {
		statusValue  string
		setName      string
		wantExitCode int
	}{
		{"completed", "terminal-success", 0},
		{"completed-with-issues", "terminal-success", 0},
		{"failed", "terminal-success", 1},
		{"cancelled", "terminal-success", 1},
		{"pending", "terminal-success", 1},
		{"cancelled", "terminal-resolved", 0},
		{"completed-with-issues", "terminal-resolved", 0},
		{"failed", "terminal-resolved", 1},
		// An alias resolves before the membership test.
		{"done", "terminal-success", 0},
	}
	for _, testCase := range testCases {
		fixturePath := writeFrontmatterFixture(t, "id: REQ-777\nstatus: "+testCase.statusValue)
		var standardOut, standardErr strings.Builder
		exitCode := runFrontmatterCommand(
			[]string{"get", fixturePath, "status", "--in-set", testCase.setName}, &standardOut, &standardErr)
		if exitCode != testCase.wantExitCode {
			t.Fatalf("--in-set %s on status %q: exit = %d, want %d (stderr %q)",
				testCase.setName, testCase.statusValue, exitCode, testCase.wantExitCode, standardErr.String())
		}
		if standardOut.String() != "" {
			t.Fatalf("--in-set printed %q to stdout, want nothing", standardOut.String())
		}
	}
}

func TestRunFrontmatterCommandAbsentFieldExitsNonZeroWithEmptyStdout(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\nstatus: pending")

	var standardOut, standardErr strings.Builder
	exitCode := runFrontmatterCommand([]string{"get", fixturePath, "domain"}, &standardOut, &standardErr)

	if exitCode == 0 {
		t.Fatal("exit = 0 for an absent field, want non-zero")
	}
	if standardOut.String() != "" {
		t.Fatalf("stdout = %q for an absent field, want nothing", standardOut.String())
	}
}

func TestRunFrontmatterCommandRejectsUsageErrors(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-777\nstatus: pending")
	testCases := []struct {
		caseName  string
		arguments []string
	}{
		{"no verb", []string{}},
		{"unknown verb", []string{"set", fixturePath, "status"}},
		{"missing field", []string{"get", fixturePath}},
		{"leftover token", []string{"get", fixturePath, "status", "extra"}},
		{"unknown set name", []string{"get", fixturePath, "status", "--in-set", "nonsense"}},
		{"missing file", []string{"get", filepath.Join(t.TempDir(), "absent.md"), "status"}},
	}
	for _, testCase := range testCases {
		var standardOut, standardErr strings.Builder
		exitCode := runFrontmatterCommand(testCase.arguments, &standardOut, &standardErr)
		if exitCode != 2 {
			t.Fatalf("%s: exit = %d, want 2 (usage error)", testCase.caseName, exitCode)
		}
		if standardErr.String() == "" {
			t.Fatalf("%s: stderr empty, want a diagnostic", testCase.caseName)
		}
	}
}

// TestRunFrontmatterCommandWarnsOnUnrecognizedStatus closes the hole REQ-112
// shipped: status and testing_status are not rows in the contract table (their
// alias maps live in their own normalizers), and the first implementation forced
// recognized=true for them — so a typo printed to stdout, exited 0, and emitted
// no warning at all. That is the exact no-feedback path this command exists to
// replace. Found by Codex's review of PR #130.
func TestRunFrontmatterCommandWarnsOnUnrecognizedStatus(t *testing.T) {
	testCases := []struct {
		fieldName   string
		rawValue    string
		wantStdout  string
		wantWarning bool
	}{
		{"status", "completed", "completed\n", false},
		{"status", "done", "completed\n", false},         // alias still resolves silently
		{"status", "completedd", "completedd\n", true},   // typo must warn
		{"status", "in-progress", "in-progress\n", true}, // hand-edited value must warn
		{"status", "blocked-dependency-cycle", "blocked-dependency-cycle\n", false},
		{"testing_status", "in_testing", "in-testing\n", false},
		{"testing_status", "half-tested", "half-tested\n", true},
	}
	for _, testCase := range testCases {
		fixturePath := writeFrontmatterFixture(t, "id: REQ-998\n"+testCase.fieldName+": "+testCase.rawValue)
		var standardOut, standardErr strings.Builder
		exitCode := runFrontmatterCommand(
			[]string{"get", fixturePath, testCase.fieldName, "--normalize"}, &standardOut, &standardErr)

		if exitCode != 0 {
			t.Fatalf("%s=%q: exit = %d, want 0", testCase.fieldName, testCase.rawValue, exitCode)
		}
		if standardOut.String() != testCase.wantStdout {
			t.Fatalf("%s=%q: stdout = %q, want %q — an unrecognized value prints what was found, never a fabricated default",
				testCase.fieldName, testCase.rawValue, standardOut.String(), testCase.wantStdout)
		}
		gotWarning := standardErr.String() != ""
		if gotWarning != testCase.wantWarning {
			t.Fatalf("%s=%q: warning present = %v, want %v (stderr %q)",
				testCase.fieldName, testCase.rawValue, gotWarning, testCase.wantWarning, standardErr.String())
		}
		if testCase.wantWarning && !strings.Contains(standardErr.String(), testCase.rawValue) {
			t.Fatalf("%s=%q: warning %q must name the offending value",
				testCase.fieldName, testCase.rawValue, standardErr.String())
		}
	}
}

// An unrecognized status must also fail the membership test rather than slipping
// through as "not completed" without a word.
func TestRunFrontmatterCommandInSetWarnsOnUnrecognizedStatus(t *testing.T) {
	fixturePath := writeFrontmatterFixture(t, "id: REQ-998\nstatus: completedd")
	var standardOut, standardErr strings.Builder
	exitCode := runFrontmatterCommand(
		[]string{"get", fixturePath, "status", "--in-set", "terminal-success"}, &standardOut, &standardErr)
	if exitCode != 1 {
		t.Fatalf("exit = %d, want 1 — a typo is not a member", exitCode)
	}
	if standardErr.String() == "" {
		t.Fatal("stderr empty — an unrecognized status must warn even on the membership path")
	}
}
