package heavyverification

import (
	"reflect"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestPlanHeavyVerificationHandlerProjectsTypedPlan(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "web/app.js", "changed\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "known change")
	result := handlePlanHeavyVerification(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{
		"--manifest", "heavy-lanes.json", "--base-revision", baseRevision, "--target-revision", targetRevision,
	})
	if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerification == nil {
		t.Fatalf("handler result = %#v", result)
	}
	if got := selectedHeavyLaneIDs(*result.HeavyVerification); !reflect.DeepEqual(got, []string{"browser-behavior"}) {
		t.Fatalf("handler selected lanes = %v", got)
	}
}

func TestPlanHeavyVerificationHandlerRejectsIncompleteArguments(t *testing.T) {
	result := handlePlanHeavyVerification(commandruntime.ExecutionContext{RepositoryRoot: t.TempDir()}, []string{"--base-revision", "HEAD"})
	if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-PLAN-USAGE" {
		t.Fatalf("usage result = %#v", result)
	}
}

func TestPlanHeavyRevalidationHandlerProjectsMultipleRanges(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "web/app.js", "changed\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "known change")
	result := handlePlanHeavyRevalidation(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{
		"--manifest", "heavy-lanes.json",
		"--source-range", baseRevision + ".." + targetRevision,
		"--source-range", baseRevision + ".." + targetRevision,
		"--execution-revision", targetRevision,
	})
	if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerification == nil {
		t.Fatalf("handler result = %#v", result)
	}
	if len(result.HeavyVerification.SourceRanges) != 2 || result.HeavyVerification.Mode != "historical-revalidation" {
		t.Fatalf("handler plan = %#v", result.HeavyVerification)
	}
}

func TestRunHeavyVerificationHandlerRejectsMissingLane(t *testing.T) {
	result := handleRunHeavyVerification(commandruntime.ExecutionContext{RepositoryRoot: t.TempDir()}, []string{"--manifest", "heavy-lanes.json"})
	if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-RUN-USAGE" {
		t.Fatalf("usage result = %#v", result)
	}
}

func TestHandlersRegisterRunHeavyVerification(t *testing.T) {
	if _, registered := Handlers()[CommandRunHeavyVerification]; !registered {
		t.Fatalf("Handlers() = %v, missing %s", Handlers(), CommandRunHeavyVerification)
	}
}

func TestPlanHeavyRevalidationHandlerRejectsIncompleteArguments(t *testing.T) {
	result := handlePlanHeavyRevalidation(commandruntime.ExecutionContext{RepositoryRoot: t.TempDir()}, []string{"--execution-revision", "HEAD"})
	if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-REVALIDATION-USAGE" {
		t.Fatalf("usage result = %#v", result)
	}
}
