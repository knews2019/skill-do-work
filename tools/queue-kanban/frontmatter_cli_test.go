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
