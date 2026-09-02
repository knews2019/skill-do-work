package requeststate

import (
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestSixLifecycleTransitionsProduceRunnablePlans(t *testing.T) {
	now := time.Date(2026, 8, 31, 20, 30, 0, 0, time.UTC)
	tests := []struct {
		name, path, status, extra string
		options                   StateOptions
	}{
		{"claim", "do-work/queue/REQ-101.md", "pending", "", StateOptions{Transition: TransitionClaim, RequestID: "REQ-101", Provenance: ProvenanceDefault, Now: now}},
		{"recover", "do-work/working/REQ-106.md", "claimed", "claimed_at: 2026-08-31T20:00:00Z\n", StateOptions{Transition: TransitionRecover, RequestID: "REQ-106", CheckpointAbsent: true, AssumeSoleWriter: true, DryRun: true, Now: now}},
		{"unblock", "do-work/queue/REQ-102.md", "blocked", "blocked_by: service ready\nblocked_check: test -e ready\n", StateOptions{Transition: TransitionUnblock, RequestID: "REQ-102", OriginalStatus: "blocked", ProbeStatus: resultmodel.ProbeSucceeded, UnblockRequired: true, UnblockSource: UnblockProbe, Now: now}},
		{"complete", "do-work/working/REQ-103.md", "claimed", "claimed_at: 2026-08-31T20:00:00Z\n", StateOptions{Transition: TransitionComplete, RequestID: "REQ-103", TerminalStatus: "completed", ImplementationHash: "abcdef0", WriterLabel: "host:/repo", Now: now}},
		{"fail", "do-work/working/REQ-104.md", "claimed", "claimed_at: 2026-08-31T20:00:00Z\n", StateOptions{Transition: TransitionFail, RequestID: "REQ-104", FailureError: "tests failed", FailureType: "code", WriterLabel: "host:/repo", Now: now}},
		{"cancel", "do-work/queue/REQ-105.md", "pending", "", StateOptions{Transition: TransitionCancel, RequestID: "REQ-105", CancellationConfirmed: true, DependentDisposition: "leave", CancellationReason: "superseded", Now: now}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			writeStateRequest(t, repositoryRoot, test.path, test.options.RequestID, test.status, test.extra)
			snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
			if err != nil {
				t.Fatal(err)
			}
			plan := BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), test.options)
			if !plan.Runnable() {
				t.Fatalf("plan refusal = %#v", plan.Refusal)
			}
			if len(plan.TargetPaths) == 0 || len(plan.Changes) == 0 {
				t.Fatalf("incomplete plan: %#v", plan)
			}
		})
	}
}

func TestClaimPreservesExplicitDependencyBypassButGatesURExpansion(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeStateRequest(t, repositoryRoot, "do-work/queue/REQ-201.md", "REQ-201", "pending", "depends_on: [REQ-999]\nassigned_to: other\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	graph := dependencygraph.BuildGraph(snapshot)
	explicit := BuildPlan(snapshot, graph, StateOptions{Transition: TransitionClaim, RequestID: "REQ-201", Provenance: ProvenanceExplicit})
	if !explicit.Runnable() {
		t.Fatalf("explicit claim lost override: %#v", explicit.Refusal)
	}
	expanded := BuildPlan(snapshot, graph, StateOptions{Transition: TransitionClaim, RequestID: "REQ-201", Provenance: ProvenanceURExpanded})
	if expanded.Refusal == nil || expanded.Refusal.Code != "CLAIM-DEPENDENCY-FAILED" {
		t.Fatalf("UR-expanded dependency gate = %#v", expanded.Refusal)
	}
}

func TestFailAcceptsOnlyCanonicalErrorTypes(t *testing.T) {
	for _, errorType := range []string{"intent", "spec", "code", "environment"} {
		t.Run("accepts "+errorType, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			writeStateRequest(t, repositoryRoot, "do-work/working/REQ-220.md", "REQ-220", "claimed", "")
			snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
			plan := BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), StateOptions{Transition: TransitionFail, RequestID: "REQ-220", FailureError: "failed", FailureType: errorType})
			if !plan.Runnable() {
				t.Fatalf("canonical error_type refused: %#v", plan.Refusal)
			}
		})
	}
	for _, errorType := range []string{"", "CODE", "security", "code "} {
		t.Run("refuses "+errorType, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			writeStateRequest(t, repositoryRoot, "do-work/working/REQ-221.md", "REQ-221", "claimed", "")
			snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
			plan := BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), StateOptions{Transition: TransitionFail, RequestID: "REQ-221", FailureError: "failed", FailureType: errorType})
			if plan.Refusal == nil || plan.Refusal.Code != "FAIL-CLASSIFICATION-INVALID" {
				t.Fatalf("noncanonical error_type plan = %#v", plan)
			}
		})
	}
}

func TestCancelReasonPreflightRejectsUnsafeControlBytes(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeStateRequest(t, repositoryRoot, "do-work/queue/REQ-222.md", "REQ-222", "pending", "")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plan := BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), StateOptions{Transition: TransitionCancel, RequestID: "REQ-222", CancellationConfirmed: true, DependentDisposition: "leave", CancellationReason: "unsafe\x01text"})
	if plan.Refusal == nil || plan.Refusal.Code != "CANCEL-REASON-UNSAFE" {
		t.Fatalf("unsafe reason plan = %#v", plan)
	}
}

func TestCancelMultilineReasonRequiresSafeSummary(t *testing.T) {
	for _, summary := range []string{"", "unsafe\nsummary", "unsafe\x7fsummary"} {
		repositoryRoot := t.TempDir()
		writeStateRequest(t, repositoryRoot, "do-work/queue/REQ-223.md", "REQ-223", "pending", "")
		snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
		plan := BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), StateOptions{Transition: TransitionCancel, RequestID: "REQ-223", CancellationConfirmed: true, DependentDisposition: "leave", CancellationReason: "line one\nline two", CancellationSummary: summary})
		if plan.Refusal == nil || (plan.Refusal.Code != "CANCEL-REASON-SUMMARY-MISSING" && plan.Refusal.Code != "CANCEL-REASON-UNSAFE") {
			t.Fatalf("unsafe/missing summary %q plan = %#v", summary, plan)
		}
	}
}
