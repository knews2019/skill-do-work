package publication

import (
	"reflect"
	"testing"
)

func TestPublicationPlanSortsTargetsAndRejectsDuplicates(t *testing.T) {
	plan := finalizePlan(PublicationPlan{RepositoryRoot: t.TempDir(), Operation: OperationCaptureFiles, Mutations: []PlannedMutation{
		{Kind: MutationCreate, Path: "do-work/queue/REQ-2.md"},
		{Kind: MutationCreate, Path: "do-work/queue/REQ-1.md"},
	}})
	if plan.Refusal != nil {
		t.Fatal(plan.Refusal)
	}
	if !reflect.DeepEqual(plan.TargetPaths, []string{"do-work/queue/REQ-1.md", "do-work/queue/REQ-2.md"}) {
		t.Fatalf("targets = %v", plan.TargetPaths)
	}
	duplicate := finalizePlan(PublicationPlan{RepositoryRoot: t.TempDir(), Operation: OperationCaptureFiles, Mutations: []PlannedMutation{
		{Kind: MutationCreate, Path: "do-work/queue/REQ-1.md"},
		{Kind: MutationReplace, Path: "do-work/queue/REQ-1.md"},
	}})
	if duplicate.Refusal == nil || duplicate.Refusal.Code != "PUBLICATION-DUPLICATE-TARGET" {
		t.Fatalf("duplicate refusal = %#v", duplicate.Refusal)
	}
}
