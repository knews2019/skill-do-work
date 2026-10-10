package main

import (
	"bytes"
	"slices"
	"testing"
	"time"
)

// Pins two failures. (1) request-commits must credit exactly what the board's
// activity correlation credits: a second copy of the rule would drift, and an
// unbracketed id in prose would make trace call a REQ built when nobody
// committed for it. (2) The log must run with no --since window and no
// pathspec: a window hides older shipped work, and a do-work/ pathspec drops
// every merge (lesson 0.305.67). The 0 count marks bookkeeping commits.
func TestRequestCommitsListsOnlyCommitsTheBoardCredits(t *testing.T) {
	commitInstant := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	var loggedArguments []string
	cannedRunner := func(_ string, arguments ...string) ([]byte, error) {
		loggedArguments = arguments
		return []byte(
			cannedLogRecord("c4", commitInstant, "[REQ-007] add the export button",
				"src/export.go", "src/export_test.go", "do-work/working/REQ-007-export-button.md") +
				cannedLogRecord("c3", commitInstant.Add(-time.Hour), "docs(do-work): note the export follow-up",
					"do-work/archive/UR-001/REQ-007-export-button.md") +
				cannedLogRecord("c2", commitInstant.Add(-2*time.Hour), "fix parser; see REQ-007", "src/parser.go") +
				cannedLogRecord("c1", commitInstant.Add(-3*time.Hour), "[REQ-007] claim request lifecycle",
					"do-work/working/REQ-007-export-button.md")), nil
	}

	var standardOut, standardErr bytes.Buffer
	exitCode := runRequestCommitsCommand([]string{"--repo-root", "/repo", "REQ-007"}, &standardOut, &standardErr, cannedRunner)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr %q)", exitCode, standardErr.String())
	}
	wantOutput := "request_id\tcommit\tcommitted_at\tpaths_outside_do_work\tsubject\n" +
		"REQ-007\tc4\t2026-10-09T12:00:00Z\t2\t[REQ-007] add the export button\n" +
		"REQ-007\tc3\t2026-10-09T11:00:00Z\t0\tdocs(do-work): note the export follow-up\n" +
		"REQ-007\tc1\t2026-10-09T09:00:00Z\t0\t[REQ-007] claim request lifecycle\n"
	if standardOut.String() != wantOutput {
		t.Errorf("stdout =\n%s\nwant\n%s", standardOut.String(), wantOutput)
	}
	wantArguments := []string{"log", "--format=%H%x00%cI%x00%s", "--name-only"}
	if !slices.Equal(loggedArguments, wantArguments) {
		t.Errorf("git arguments = %q, want %q (no --since window, no pathspec)", loggedArguments, wantArguments)
	}
}

// Pins a misspelled or misplaced argument being read as "no commits": trace
// would then report a built ask as `not started`. Anything that is not REQ-
// plus digits (a UR id here) is a usage error with nothing on stdout.
func TestRequestCommitsRefusesAnArgumentThatIsNotARequestId(t *testing.T) {
	runnerCalled := false
	cannedRunner := func(string, ...string) ([]byte, error) {
		runnerCalled = true
		return nil, nil
	}

	var standardOut, standardErr bytes.Buffer
	exitCode := runRequestCommitsCommand([]string{"--repo-root", "/repo", "UR-012"}, &standardOut, &standardErr, cannedRunner)

	if exitCode != 2 {
		t.Errorf("exit code = %d, want 2", exitCode)
	}
	if standardOut.Len() != 0 {
		t.Errorf("stdout = %q, want nothing", standardOut.String())
	}
	if runnerCalled {
		t.Errorf("git ran for a refused argument")
	}
}
