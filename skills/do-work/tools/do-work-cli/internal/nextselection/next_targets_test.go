package nextselection

import (
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestExplicitREQRefusesAmbiguousQueueIdentityIndependentOfDiscoveryOrder(t *testing.T) {
	repositoryRoot := t.TempDir()
	firstPath := "do-work/queue/REQ-452-first.md"
	secondPath := "do-work/queue/REQ-0452-second.md"
	writeCommandRequest(t, repositoryRoot, firstPath, "REQ-452", "pending", "depends_on: [REQ-999]\nassigned_to: cloud-alpha\nimpact: impact-negligible\n")
	writeCommandRequest(t, repositoryRoot, secondPath, "REQ-0452", "pending", "")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}

	var firstExclusion *resultmodel.SelectionExclusion
	for _, replay := range []struct {
		name    string
		reverse bool
	}{
		{name: "discovery order"},
		{name: "reversed discovery order", reverse: true},
	} {
		t.Run(replay.name, func(t *testing.T) {
			replaySnapshot := *snapshot
			replaySnapshot.RequestFiles = append([]*repositorymodel.RequestFile(nil), snapshot.RequestFiles...)
			if replay.reverse {
				for left, right := 0, len(replaySnapshot.RequestFiles)-1; left < right; left, right = left+1, right-1 {
					replaySnapshot.RequestFiles[left], replaySnapshot.RequestFiles[right] = replaySnapshot.RequestFiles[right], replaySnapshot.RequestFiles[left]
				}
			}
			result := Select(&replaySnapshot, dependencygraph.BuildGraph(&replaySnapshot), SelectionOptions{
				TargetTokens:         []string{"REQ-452"},
				SkipImpactNegligible: true,
			}, nil)
			if len(result.Selected) != 0 {
				t.Fatalf("ambiguous explicit target selected an arbitrary record: %#v", result.Selected)
			}
			if len(result.Excluded) != 1 {
				t.Fatalf("exclusions = %#v, want one typed ambiguity exclusion", result.Excluded)
			}
			exclusion := result.Excluded[0]
			if exclusion.Code != "DEPENDENCY-AMBIGUOUS" || exclusion.Provenance != ProvenanceExplicit || exclusion.RequestPath != "" {
				t.Fatalf("ambiguity exclusion = %#v", exclusion)
			}
			for _, collisionPath := range []string{firstPath, secondPath} {
				if !strings.Contains(exclusion.Reason, collisionPath) {
					t.Errorf("ambiguity reason omitted collision path %q: %q", collisionPath, exclusion.Reason)
				}
			}
			if firstExclusion == nil {
				firstExclusion = &exclusion
			} else if !reflect.DeepEqual(*firstExclusion, exclusion) {
				t.Fatalf("discovery order changed exclusion:\nfirst: %#v\nnext:  %#v", *firstExclusion, exclusion)
			}
		})
	}
}

func TestExplicitREQRefusesNormalizedFrontmatterIdentityCollisionIndependentOfDiscoveryOrder(t *testing.T) {
	repositoryRoot := t.TempDir()
	firstPath := "do-work/queue/REQ-900-first.md"
	secondPath := "do-work/queue/REQ-901-second.md"
	writeCommandRequest(t, repositoryRoot, firstPath, "REQ-452", "pending", "depends_on: [REQ-999]\nassigned_to: cloud-alpha\nimpact: impact-negligible\n")
	writeCommandRequest(t, repositoryRoot, secondPath, "REQ-0452", "pending", "")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}

	var firstExclusion *resultmodel.SelectionExclusion
	for _, replay := range []struct {
		name    string
		reverse bool
	}{
		{name: "discovery order"},
		{name: "reversed discovery order", reverse: true},
	} {
		t.Run(replay.name, func(t *testing.T) {
			replaySnapshot := *snapshot
			replaySnapshot.RequestFiles = append([]*repositorymodel.RequestFile(nil), snapshot.RequestFiles...)
			if replay.reverse {
				for left, right := 0, len(replaySnapshot.RequestFiles)-1; left < right; left, right = left+1, right-1 {
					replaySnapshot.RequestFiles[left], replaySnapshot.RequestFiles[right] = replaySnapshot.RequestFiles[right], replaySnapshot.RequestFiles[left]
				}
			}
			result := Select(&replaySnapshot, dependencygraph.BuildGraph(&replaySnapshot), SelectionOptions{
				TargetTokens:         []string{"REQ-452"},
				SkipImpactNegligible: true,
			}, nil)
			if len(result.Selected) != 0 {
				t.Fatalf("ambiguous explicit target selected an arbitrary record: %#v", result.Selected)
			}
			if len(result.Excluded) != 1 {
				t.Fatalf("exclusions = %#v, want one typed ambiguity exclusion", result.Excluded)
			}
			exclusion := result.Excluded[0]
			if exclusion.Code != "DEPENDENCY-AMBIGUOUS" || exclusion.Provenance != ProvenanceExplicit || exclusion.RequestPath != "" {
				t.Fatalf("ambiguity exclusion = %#v", exclusion)
			}
			for _, collisionPath := range []string{firstPath, secondPath} {
				if !strings.Contains(exclusion.Reason, collisionPath) {
					t.Errorf("ambiguity reason omitted collision path %q: %q", collisionPath, exclusion.Reason)
				}
			}
			if firstExclusion == nil {
				firstExclusion = &exclusion
			} else if !reflect.DeepEqual(*firstExclusion, exclusion) {
				t.Fatalf("discovery order changed exclusion:\nfirst: %#v\nnext:  %#v", *firstExclusion, exclusion)
			}
		})
	}
}

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

func TestTargetResolutionOrdersUREntriesByRequestPriorityButKeepsExplicitOrder(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-851-later.md", "REQ-851", "pending", "user_request: UR-850\npriority: later\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-852-now.md", "REQ-852", "pending", "user_request: UR-850\npriority: now\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	graph := dependencygraph.BuildGraph(snapshot)
	for _, testCase := range []struct {
		name   string
		tokens []string
		want   []string
	}{
		{"UR expansion", []string{"UR-850"}, []string{"REQ-852", "REQ-851"}},
		{"explicit caller order", []string{"REQ-851", "REQ-852"}, []string{"REQ-851", "REQ-852"}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			candidates, exclusions := resolveTargets(snapshot, graph, SelectionOptions{TargetTokens: testCase.tokens})
			if len(exclusions) != 0 {
				t.Fatalf("exclusions = %#v", exclusions)
			}
			got := make([]string, 0, len(candidates))
			for _, candidate := range candidates {
				got = append(got, candidate.RequestID)
			}
			if !equalStrings(got, testCase.want) {
				t.Fatalf("resolved order = %v, want %v", got, testCase.want)
			}
		})
	}
}
