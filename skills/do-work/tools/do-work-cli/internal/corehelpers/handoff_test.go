package corehelpers

import (
	"bytes"
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
	}
	if !seen["HANDOFF-INTEGRATION-BRANCH"] || !seen["HANDOFF-WORKTREE-CLEAN"] {
		t.Fatalf("codes=%v", seen)
	}
}
