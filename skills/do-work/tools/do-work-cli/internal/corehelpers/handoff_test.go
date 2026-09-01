package corehelpers

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestHandoffSurveyNamesIntegrationAndDirtyCheckout(t *testing.T) {
	repository := newGitFixture(t)
	branchOutput, branchError := gitOutput(repository, "branch", "--show-current")
	if branchError != nil {
		t.Fatal(branchError)
	}
	result := handleHandoffSurvey(testContext(repository), []string{"--integration-branch", string(bytes.TrimSpace(branchOutput))})
	if result.Outcome != "success" {
		t.Fatalf("result=%+v", result)
	}
	seen := map[string]bool{}
	for _, finding := range result.Findings {
		seen[finding.Code] = true
		if finding.Code == "HANDOFF-RECENT-COMMIT" {
			if len(finding.AffectedIDs) != 1 || len(finding.Evidence) != 1 || strings.Contains(finding.Evidence[0], "\n") {
				t.Fatalf("commit finding is not one structured record: %#v", finding)
			}
		}
		if finding.Code == "HANDOFF-WORKTREE" {
			if len(finding.AffectedPaths) != 1 || len(finding.Evidence) != 3 {
				t.Fatalf("worktree finding is not one structured record: %#v", finding)
			}
			for _, evidence := range finding.Evidence {
				if strings.Contains(evidence, "\n") {
					t.Fatalf("opaque multiline evidence survived: %#v", finding)
				}
			}
		}
	}
	if !seen["HANDOFF-INTEGRATION-BRANCH"] || !seen["HANDOFF-RECENT-COMMIT"] || !seen["HANDOFF-WORKTREE"] || !seen["HANDOFF-WORKTREE-CLEAN"] {
		t.Fatalf("codes=%v", seen)
	}
}

func TestHandoffStatusRowsAreTypedPerPathAndPreserveHostileNames(t *testing.T) {
	rows, err := parseHandoffStatus([]byte("R  renamed path\x00old path\x00?? line\nbreak\x00 M spaced name\x00"))
	if err != nil {
		t.Fatal(err)
	}
	want := []handoffStatusRow{
		{status: "R ", path: "renamed path", origin: "old path"},
		{status: "??", path: "line\nbreak"},
		{status: " M", path: "spaced name"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("rows=%#v want=%#v", rows, want)
	}
	for _, row := range rows {
		if row.path == "" || strings.Contains(row.status, "\n") {
			t.Fatalf("unstructured row: %#v", row)
		}
	}
	mutation := bytes.Replace([]byte("R  renamed path\x00old path\x00"), []byte("\x00old path\x00"), []byte("\x00"), 1)
	if _, err := parseHandoffStatus(mutation); err == nil {
		t.Fatal("rename-origin collapse mutation escaped")
	}
}
