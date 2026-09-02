package nextselection

import (
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

func TestTargetResolutionPreservesMixedTokenOrderAndExplicitProvenance(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-042-first.md", "REQ-042", "pending", "user_request: UR-011\ndepends_on: [REQ-099]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-043-second.md", "REQ-043", "pending", "user_request: UR-011\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	graph := dependencygraph.BuildGraph(snapshot)
	candidates, exclusions := resolveTargets(snapshot, graph, SelectionOptions{TargetTokens: []string{"Ur-11", "req-42"}})
	if len(exclusions) != 0 || len(candidates) != 2 {
		t.Fatalf("resolution = %#v, exclusions=%#v", candidates, exclusions)
	}
	if candidates[0].RequestID != "REQ-043" || candidates[0].Provenance != ProvenanceUserRequest {
		t.Fatalf("UR expansion did not retain its token position: %#v", candidates[0])
	}
	if candidates[1].RequestID != "REQ-042" || candidates[1].Provenance != ProvenanceExplicit {
		t.Fatalf("explicit duplicate was not reserved for its caller position: %#v", candidates[1])
	}

	_, missing := resolveTargets(snapshot, graph, SelectionOptions{TargetTokens: []string{"REQ-999", "UR-999"}})
	if len(missing) != 2 || missing[0].Code != "TARGET-NOT-FOUND" || missing[1].Code != "TARGET-NOT-FOUND" {
		t.Fatalf("missing targets were not actionable exclusions: %#v", missing)
	}
	for _, exclusion := range missing {
		if len(exclusion.NextArgv) == 0 || len(exclusion.VerificationArgv) == 0 || exclusion.NextJustRecipe == "" {
			t.Fatalf("missing target has incomplete commands: %#v", exclusion)
		}
	}
}
