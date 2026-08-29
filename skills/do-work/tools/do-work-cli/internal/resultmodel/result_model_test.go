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
	if !ok || rollback["actions"] == nil {
		t.Fatalf("rollback must be a typed object with non-null actions: %#v", decoded["rollback"])
	}

	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatalf("RenderResult text: %v", err)
	}
	for _, required := range []string{
		"doctor: findings", "DIRTY-TARGET", "tracked.txt has worktree changes",
		"git diff -- tracked.txt", "just do-work-doctor",
		"git status --short -- tracked.txt",
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
