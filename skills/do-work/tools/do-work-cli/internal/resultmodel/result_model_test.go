package resultmodel

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRenderersUseOneNormalizedResult(t *testing.T) {
	result := CommandResult{
		Command:        "doctor",
		Outcome:        OutcomeFindings,
		RepositoryRoot: "/tmp/example",
		Findings: []CommandFinding{{
			Code:                 "DIRTY-TARGET",
			Severity:             SeverityWarning,
			AffectedIDs:          []string{"REQ-406"},
			AffectedPaths:        []string{"tracked.txt"},
			Evidence:             []string{"tracked.txt has worktree changes"},
			Fixability:           FixabilityRefused,
			AutomationStopReason: "the target is already dirty",
			NextArgv:             []string{"git", "diff", "--", "tracked.txt"},
			NextJustRecipe:       "do-work-doctor",
			VerificationArgv:     []string{"git", "status", "--short", "--", "tracked.txt"},
		}},
		Rollback: RollbackResult{
			Status:  RollbackSucceeded,
			Actions: []string{"restored tracked.txt from HEAD"},
		},
	}

	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatalf("RenderResult JSON: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if decoded["schema_version"] != float64(SchemaVersion) {
		t.Fatalf("schema_version = %#v, want %d", decoded["schema_version"], SchemaVersion)
	}
	for _, field := range []string{"findings", "changes", "skipped_work"} {
		if decoded[field] == nil {
			t.Fatalf("%s must be a non-null collection", field)
		}
	}
	rollback, ok := decoded["rollback"].(map[string]any)
	if !ok || rollback["actions"] == nil || rollback["errors"] == nil {
		t.Fatalf("rollback must be a typed object with non-null collections: %#v", decoded["rollback"])
	}
	// RollbackStatus is a named type with one definition here; the wire form stays a plain
	// string so a second copy of these constants cannot appear elsewhere and drift.
	if rollback["status"] != string(RollbackSucceeded) {
		t.Fatalf("rollback status = %#v, want %q", rollback["status"], RollbackSucceeded)
	}

	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatalf("RenderResult text: %v", err)
	}
	for _, required := range []string{
		"doctor: findings", "DIRTY-TARGET", "tracked.txt has worktree changes",
		"git diff -- tracked.txt", "just do-work-doctor",
		"git status --short -- tracked.txt",
		"rollback: succeeded", "rollback action: restored tracked.txt from HEAD",
	} {
		if !strings.Contains(string(textOutput), required) {
			t.Errorf("text output missing %q:\n%s", required, textOutput)
		}
	}
}

func TestOutcomeExitCodes(t *testing.T) {
	tests := []struct {
		outcome CommandOutcome
		want    int
	}{
		{OutcomeSuccess, 0},
		{OutcomeFindings, 1},
		{OutcomeRefused, 1},
		{OutcomeFailure, 2},
		{OutcomeRolledBack, 3},
		{OutcomeRisk, 4},
	}
	for _, test := range tests {
		if got := ExitCode(test.outcome); got != test.want {
			t.Errorf("ExitCode(%q) = %d, want %d", test.outcome, got, test.want)
		}
	}
}

func TestRenderRejectsUnknownFormat(t *testing.T) {
	if _, err := RenderResult(CommandResult{Outcome: OutcomeSuccess}, OutputFormat("yaml")); err == nil {
		t.Fatal("RenderResult accepted unsupported output format")
	}
}

// The parity test covers findings; changes, skipped work and rollback errors render
// unasserted. A reader who only ever sees the text form would not notice one of these three
// sections disappearing, so each is pinned to the exact line it produces.
func TestTextRenderingNamesChangesSkippedWorkAndRollbackErrors(t *testing.T) {
	rendered, err := RenderResult(CommandResult{
		Command:        "install-suite",
		Outcome:        OutcomeRisk,
		RepositoryRoot: "/tmp/example",
		Changes: []RecordedChange{
			{Path: ".claude/skills/do-work", Kind: "created", Detail: "installed do-work suite v1.2.3"},
		},
		SkippedWork: []SkippedWork{
			{Code: "INSTALL-CANCELLED", Reason: "the single install confirmation was declined"},
		},
		Rollback: RollbackResult{
			Status:  RollbackIncomplete,
			Actions: []string{"restored justfile from the pre-install snapshot"},
			Errors:  []string{"could not restore the Git index"},
		},
	}, FormatText)
	if err != nil {
		t.Fatalf("RenderResult: %v", err)
	}
	for _, expectedLine := range []string{
		"install-suite: committed_state_risk",
		"repository: /tmp/example",
		"change .claude/skills/do-work [created]: installed do-work suite v1.2.3",
		"skipped INSTALL-CANCELLED: the single install confirmation was declined",
		"rollback: incomplete",
		"  rollback action: restored justfile from the pre-install snapshot",
		"  rollback error: could not restore the Git index",
	} {
		if !containsExactLine(string(rendered), expectedLine) {
			t.Errorf("text output is missing the exact line %q:\n%s", expectedLine, rendered)
		}
	}
}

// RollbackStatus has a fourth wire value — the empty string — for a result that never ran a
// Git transaction. Normalising it to not_needed would make every read-only command print a
// rollback line implying a mutation was possible, so the empty value is deliberate and both
// renderings must carry it through.
func TestAResultThatRanNoTransactionRendersNoRollbackLine(t *testing.T) {
	readOnlyResult := CommandResult{
		Command:        "validate-manifest",
		Outcome:        OutcomeSuccess,
		RepositoryRoot: "/tmp/example",
	}
	renderedText, err := RenderResult(readOnlyResult, FormatText)
	if err != nil {
		t.Fatalf("RenderResult text: %v", err)
	}
	if strings.Contains(string(renderedText), "rollback:") {
		t.Errorf("a result with no transaction printed a rollback line:\n%s", renderedText)
	}

	renderedJSON, err := RenderResult(readOnlyResult, FormatJSON)
	if err != nil {
		t.Fatalf("RenderResult JSON: %v", err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(renderedJSON, &decoded); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	if decoded.Rollback.Status != "" {
		t.Errorf("rollback status = %q, want the empty wire value", decoded.Rollback.Status)
	}
	if decoded.Rollback.Actions == nil || decoded.Rollback.Errors == nil {
		t.Errorf("rollback arrays must normalise to empty rather than null: %#v", decoded.Rollback)
	}
	if ExitCode(decoded.Outcome) != 0 {
		t.Errorf("a successful read-only command exits %d, want 0", ExitCode(decoded.Outcome))
	}
}

func containsExactLine(rendered, wanted string) bool {
	for _, line := range strings.Split(rendered, "\n") {
		if line == wanted {
			return true
		}
	}
	return false
}
