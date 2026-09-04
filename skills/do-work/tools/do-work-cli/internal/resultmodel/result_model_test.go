package resultmodel

import (
	"encoding/json"
	"reflect"
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
	for _, field := range []string{"findings", "changes", "skipped_work", "selected", "excluded"} {
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

func TestFinalizationTextAndJSONCarryTheSameOrderedEvidence(t *testing.T) {
	result := CommandResult{
		Command: "advance", Outcome: OutcomeRefused, RepositoryRoot: "/tmp/example",
		Finalizations: []FinalizationResult{{
			RequestID: "REQ-507", RequestPath: "do-work/working/REQ-507-tail.md", ArchivePath: "do-work/archive/REQ-507-tail.md",
			JournalPath: "/tmp/example/.git/do-work-finalization/REQ-507.json", Phase: "verified", TerminalStatus: "completed-with-issues",
			Resumed: true, Discovered: false, CommitPaths: []string{"do-work/working/REQ-507-tail.md", "do-work/archive/REQ-507-tail.md"},
			PrimaryCommit: "primary123", MetadataCommit: "metadata456", CreatedPrimaryCommit: "", CreatedMetadataCommit: "metadata456",
			BlockedPaths: []string{"foreign.txt"}, ReasonCodes: []string{"FINALIZATION-VERIFY"},
			NextArgv: []string{"do-work-cli", "recover-finalization"}, VerificationArgv: []string{"git", "status", "--short"}, CollectionArgv: []string{"git", "diff", "--", "foreign.txt"},
		}},
	}

	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"finalization REQ-507 [verified, completed-with-issues]", "request: do-work/working/REQ-507-tail.md", "archive: do-work/archive/REQ-507-tail.md",
		"journal: /tmp/example/.git/do-work-finalization/REQ-507.json", "resumed: true", "discovered: false",
		"commit paths: do-work/working/REQ-507-tail.md, do-work/archive/REQ-507-tail.md", "primary commit: primary123", "metadata commit: metadata456",
		"created primary commit:", "created metadata commit: metadata456", "blocked paths: foreign.txt", "reason codes: FINALIZATION-VERIFY",
		"next: do-work-cli recover-finalization", "verify: git status --short", "collect: git diff -- foreign.txt",
	} {
		if !strings.Contains(string(textOutput), expected) {
			t.Errorf("finalization text missing %q:\n%s", expected, textOutput)
		}
	}

	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Finalizations) != 1 || decoded.Finalization == nil || !reflect.DeepEqual(*decoded.Finalization, decoded.Finalizations[0]) || decoded.Finalizations[0].BlockedPaths == nil || decoded.Finalizations[0].ReasonCodes == nil {
		t.Fatalf("finalization JSON lost ordered/compatibility evidence: %#v", decoded)
	}

	ordered := result
	ordered.Finalizations = append(ordered.Finalizations, FinalizationResult{RequestID: "REQ-508", Phase: "cleanup_complete", TerminalStatus: "completed"})
	orderedText, err := RenderResult(ordered, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	firstOffset := strings.Index(string(orderedText), "finalization REQ-507")
	secondOffset := strings.Index(string(orderedText), "finalization REQ-508")
	if firstOffset < 0 || secondOffset <= firstOffset {
		t.Fatalf("finalization text changed record order:\n%s", orderedText)
	}
	orderedJSON, err := RenderResult(ordered, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	decoded = CommandResult{}
	if err := json.Unmarshal(orderedJSON, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Finalizations) != 2 || decoded.Finalization != nil || decoded.Finalizations[0].RequestID != "REQ-507" || decoded.Finalizations[1].RequestID != "REQ-508" {
		t.Fatalf("finalization JSON changed record order: %#v", decoded)
	}
}

func TestSelectionTextAndJSONCarryTheSameTypedCommands(t *testing.T) {
	result := CommandResult{
		Command: "next", Outcome: OutcomeSuccess, RepositoryRoot: "/tmp/example",
		Selected: []SelectionRecord{{
			RequestID: "REQ-007", RequestPath: "do-work/queue/REQ-007-ready.md", Title: "Ready work", Provenance: "explicit-req", OriginalStatus: "blocked",
			RequestPriority: "now",
			ProbeStatus:     ProbeSucceeded, ProbeAttempted: true, ProbeExitCode: 0, UnblockRequired: true, DependencyDepth: 0,
			EstimateMinutes: 10, EstimateKnown: true,
			NextArgv: []string{"do-work", "run", "REQ-007"}, NextJustRecipe: "do-work-run REQ-007",
			VerificationArgv: []string{"do-work-cli", "--format", "json", "next", "REQ-007"},
		}},
		Excluded: []SelectionExclusion{{
			RequestID: "REQ-008", RequestPath: "do-work/queue/REQ-008-waiting.md", Title: "Waiting work", Provenance: "ur-expanded", OriginalStatus: "pending",
			RequestPriority: "later",
			ProbeStatus:     ProbeNotApplicable, ProbeExitCode: -1,
			Code: "DEPENDENCIES-UNMET", Reason: "waits on REQ-007",
			ClaimEvidence: []SelectionClaimEvidence{{Source: "checkpoint", ClaimedAt: "2026-09-01T10:00:00Z", Writer: "host:/repo", Path: "do-work/CHECKPOINT.md", SourceLine: 12, HeaderText: "- REQ-008: Waiting work — claimed 2026-09-01T10:00:00Z — writer: host:/repo"}},
			NextArgv:      []string{"do-work", "run", "REQ-007"}, NextJustRecipe: "do-work-run REQ-007",
			VerificationArgv: []string{"do-work-cli", "--format", "json", "next", "REQ-008"},
		}},
		SelectionSummary: SelectionSummary{Pending: 2, TotalEstimatedMinutes: 10},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"selected REQ-007 [explicit-req, ordinary, depth 0, 10 min]",
		"request: do-work/queue/REQ-007-ready.md (original status: blocked)",
		"probe succeeded (attempted: true, exit: 0); unblock required",
		"excluded REQ-008 [ur-expanded, ordinary] DEPENDENCIES-UNMET: waits on REQ-007",
		"request: do-work/queue/REQ-008-waiting.md (original status: pending)",
		"probe not_applicable (attempted: false, exit: -1)",
		"claim checkpoint: claimed_at=2026-09-01T10:00:00Z writer=host:/repo path=do-work/CHECKPOINT.md line=12",
		"header: - REQ-008: Waiting work — claimed 2026-09-01T10:00:00Z — writer: host:/repo",
		"next: do-work run REQ-007", "just: just do-work-run REQ-007",
		"verify: do-work-cli --format json next REQ-008", "run_set: REQ-007",
	} {
		if !strings.Contains(string(textOutput), required) {
			t.Errorf("selection text missing %q:\n%s", required, textOutput)
		}
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Selected) != 1 || len(decoded.Excluded) != 1 || decoded.Selected[0].NextArgv[2] != "REQ-007" || decoded.Excluded[0].Code != "DEPENDENCIES-UNMET" || len(decoded.Excluded[0].ClaimEvidence) != 1 || decoded.Excluded[0].ClaimEvidence[0].SourceLine != 12 || decoded.Selected[0].RequestPath != "do-work/queue/REQ-007-ready.md" || !decoded.Selected[0].UnblockRequired || decoded.Selected[0].ProbeStatus != ProbeSucceeded || decoded.Selected[0].RequestPriority != "now" || decoded.Excluded[0].RequestPriority != "later" {
		t.Fatalf("selection JSON lost typed records: %#v", decoded)
	}
}

func TestNormalizeResultUsesEmptyClaimEvidenceArray(t *testing.T) {
	normalized := NormalizeResult(CommandResult{Selected: []SelectionRecord{{RequestID: "REQ-000"}}, Excluded: []SelectionExclusion{{RequestID: "REQ-001"}}})
	if normalized.Selected[0].RequestPriority != "next" || normalized.Excluded[0].RequestPriority != "next" {
		t.Fatalf("request priority defaults = selected %q excluded %q, want next", normalized.Selected[0].RequestPriority, normalized.Excluded[0].RequestPriority)
	}
	if normalized.Excluded[0].ClaimEvidence == nil {
		t.Fatal("claim evidence remained nil")
	}
	encoded, err := json.Marshal(normalized)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), `"claim_evidence":[]`) {
		t.Fatalf("normalized JSON = %s", encoded)
	}
}

func TestQueueAdvanceTextAndJSONPreserveTypedSessionState(t *testing.T) {
	result := CommandResult{
		Command: "advance", Outcome: OutcomeFindings,
		QueueAdvance: &QueueAdvanceResult{
			TargetTokens: []string{"UR-500"}, DispatchBound: 1, Partial: true,
			FrozenMembers:    []QueueAdvanceMember{{RequestID: "REQ-501", RequestPath: "do-work/working/REQ-501.md", Provenance: "ur-expanded", Consumed: true}},
			Claimed:          []QueueAdvanceMember{{RequestID: "REQ-501", RequestPath: "do-work/working/REQ-501.md", Provenance: "ur-expanded", Consumed: true}},
			Phases:           []QueueAdvancePhase{{RequestID: "REQ-501", Phase: "claim", Outcome: OutcomeSuccess}},
			ContinuationArgv: []string{"do-work-cli", "--format", "json", "advance", "UR-500", "--dispatch-bound", "1"},
			VerificationArgv: []string{"git", "status", "--short", "--", "do-work"},
		},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"queue advance: members=1 claimed=1 bound=1 partial=true",
		"REQ-501 claim: success",
		"continue: do-work-cli --format json advance UR-500 --dispatch-bound 1",
		"verify: git status --short -- do-work",
	} {
		if !strings.Contains(string(textOutput), expected) {
			t.Errorf("queue text missing %q:\n%s", expected, textOutput)
		}
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.QueueAdvance == nil || len(decoded.QueueAdvance.FrozenMembers) != 1 || len(decoded.QueueAdvance.Claimed) != 1 || len(decoded.QueueAdvance.Phases) != 1 || decoded.QueueAdvance.Phases[0].Changes == nil || decoded.QueueAdvance.Phases[0].Findings == nil || decoded.QueueAdvance.ContinuationArgv[4] != "UR-500" {
		t.Fatalf("queue JSON lost typed state: %#v", decoded.QueueAdvance)
	}
}

func TestGateEvidenceTextAndJSONCarryTheSameTypedState(t *testing.T) {
	result := CommandResult{
		Command: "check-green-gate", Outcome: OutcomeSuccess, RepositoryRoot: "/tmp/example",
		GateEvidence: &GateEvidenceResult{
			RepositoryIdentity: "/tmp/example/.git", GateCommand: []string{"bash", "verify.sh"},
			GateCommandSHA256: "abc123", RecordPath: "/tmp/example/.git/do-work-green-gates/abc123.json",
			RecordProvenance: "persisted_green_run", GateExitStatus: 0,
			RecordedRevision: "1111111", HeadRevision: "2222222", TargetRevision: "3333333",
			State:   GateEvidenceLogDescendantMatch,
			Matches: true, MatchBasis: "gate_log_only_descendant", BaselineRevision: "2222222",
		},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"gate evidence: state=gate_log_descendant_match matches=true basis=gate_log_only_descendant",
		"repository identity: /tmp/example/.git", "gate command: bash verify.sh (sha256: abc123, exit: 0)",
		"provenance: persisted_green_run", "recorded=1111111 head=2222222 baseline=2222222 target=3333333",
	} {
		if !strings.Contains(string(textOutput), expected) {
			t.Errorf("text output missing %q:\n%s", expected, textOutput)
		}
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.GateEvidence == nil || decoded.GateEvidence.State != GateEvidenceLogDescendantMatch || !decoded.GateEvidence.Matches || decoded.GateEvidence.BaselineRevision != "2222222" || decoded.GateEvidence.TargetRevision != "3333333" || !reflect.DeepEqual(decoded.GateEvidence.GateCommand, []string{"bash", "verify.sh"}) {
		t.Fatalf("JSON lost typed gate evidence: %#v", decoded.GateEvidence)
	}

	normalized := NormalizeResult(CommandResult{GateEvidence: &GateEvidenceResult{}})
	if normalized.GateEvidence.GateCommand == nil {
		t.Fatal("gate command remained nil")
	}
}

func TestHeavyVerificationRunTextAndJSONCarryTheSameTypedLanes(t *testing.T) {
	result := CommandResult{
		Command: "run-heavy-verification", Outcome: OutcomeSuccess, RepositoryRoot: "/tmp/example",
		HeavyVerificationRun: &HeavyVerificationRun{
			ManifestPath: "_dev/tests/heavy-lanes.json", ExecutionRevision: "4444444",
			Lanes: []HeavyLaneExecution{
				{LaneID: "queue-kanban-browser", CommandArgv: []string{"bash", "verify.sh", "--heavy-lane", "queue-kanban-browser"}, ExitStatus: 0, Skipped: true, WallSeconds: 2, Disposition: "executed", DispositionReason: "fingerprint_mismatch"},
				{LaneID: "update-script", CommandArgv: []string{"bash", "verify.sh", "--heavy-lane", "update-script"}, ExitStatus: 3, WallSeconds: 41, Disposition: "executed", DispositionReason: "no_prior_evidence"},
				{LaneID: "installer", CommandArgv: []string{"bash", "verify.sh", "--heavy-lane", "installer"}, ExitStatus: 0, WallSeconds: 0, Disposition: "reused", DispositionReason: "fingerprint_match", FingerprintSHA256: "abc123", EvidenceRevision: "3333333", EvidenceRecordedAt: "2026-09-04T18:00:00Z"},
			},
		},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	// The disposition renders beside every lane, and the inherited record is
	// named, so a reader of the text alone can tell which greens were measured.
	for _, expected := range []string{
		"heavy verification run: lanes=3",
		"manifest: _dev/tests/heavy-lanes.json",
		"execution revision: 4444444",
		"lane queue-kanban-browser: skipped in 2s [executed: fingerprint_mismatch] — bash verify.sh --heavy-lane queue-kanban-browser",
		"lane update-script: exit 3 in 41s [executed: no_prior_evidence] — bash verify.sh --heavy-lane update-script",
		"lane installer: exit 0 in 0s [reused: fingerprint_match] — bash verify.sh --heavy-lane installer",
		"inherited evidence: revision=3333333 recorded_at=2026-09-04T18:00:00Z",
	} {
		if !strings.Contains(string(textOutput), expected) {
			t.Errorf("text output missing %q:\n%s", expected, textOutput)
		}
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.HeavyVerificationRun == nil || decoded.HeavyVerificationRun.ExecutionRevision != "4444444" || len(decoded.HeavyVerificationRun.Lanes) != 3 {
		t.Fatalf("JSON lost the typed run: %#v", decoded.HeavyVerificationRun)
	}
	if decoded.HeavyVerificationRun.Lanes[0].WallSeconds != 2 || !decoded.HeavyVerificationRun.Lanes[0].Skipped || decoded.HeavyVerificationRun.Lanes[1].WallSeconds != 41 || decoded.HeavyVerificationRun.Lanes[1].ExitStatus != 3 {
		t.Fatalf("JSON lost per-lane state: %#v", decoded.HeavyVerificationRun.Lanes)
	}
	reusedLane := decoded.HeavyVerificationRun.Lanes[2]
	if reusedLane.Disposition != "reused" || reusedLane.DispositionReason != "fingerprint_match" || reusedLane.EvidenceRevision != "3333333" || reusedLane.EvidenceRecordedAt != "2026-09-04T18:00:00Z" || reusedLane.FingerprintSHA256 != "abc123" {
		t.Fatalf("JSON lost the reuse record: %#v", reusedLane)
	}

	normalized := NormalizeResult(CommandResult{HeavyVerificationRun: &HeavyVerificationRun{Lanes: []HeavyLaneExecution{{LaneID: "solo-lane"}}}})
	if normalized.HeavyVerificationRun.Lanes[0].CommandArgv == nil {
		t.Fatal("lane argv remained nil")
	}
	if normalized.HeavyVerificationRun.Lanes[0].Disposition != "executed" {
		t.Fatalf("an unstated disposition must normalize to executed, got %q", normalized.HeavyVerificationRun.Lanes[0].Disposition)
	}
	if NormalizeResult(CommandResult{HeavyVerificationRun: &HeavyVerificationRun{}}).HeavyVerificationRun.Lanes == nil {
		t.Fatal("lanes remained nil")
	}
}

func TestAlreadyGreenRepairTextAndJSONCarryTheSameTypedState(t *testing.T) {
	result := CommandResult{
		Command: "validate-already-green-repair", Outcome: OutcomeSuccess, RepositoryRoot: "/tmp/example",
		AlreadyGreenRepair: &AlreadyGreenRepairValidation{
			RequestID: "REQ-701", RequestPath: "do-work/working/REQ-701.md",
			TDDAllowed: true, ReviewAllowed: false,
			IntakeFingerprint: "sha256:intake", ExpectedFingerprint: "sha256:intake",
			GateCommand: []string{"bash", "verify.sh"}, RecordedRevision: "1111111",
			CanonicalCompletionPaths: []string{"do-work/archive/REQ-701.md", "do-work/working/REQ-701.md"},
			StagedPaths:              []string{"do-work/archive/REQ-999.md"}, ReasonCodes: []string{"REPAIR-STAGED-PATH-NOT-CANONICAL"},
			OffendingPaths: []string{"do-work/archive/REQ-999.md"}, Writer: "fixture:/repo", PlannedAt: "2026-09-02T05:00:00Z",
		},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"already-green repair: request=REQ-701 tdd_allowed=true review_allowed=false",
		"fingerprints: intake=sha256:intake expected=sha256:intake",
		"reason codes: REPAIR-STAGED-PATH-NOT-CANONICAL",
		"offending paths: do-work/archive/REQ-999.md",
	} {
		if !strings.Contains(string(textOutput), expected) {
			t.Errorf("text output missing %q:\n%s", expected, textOutput)
		}
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.AlreadyGreenRepair == nil || !decoded.AlreadyGreenRepair.TDDAllowed || decoded.AlreadyGreenRepair.ReviewAllowed || len(decoded.AlreadyGreenRepair.CanonicalCompletionPaths) != 2 {
		t.Fatalf("JSON lost repair validation: %#v", decoded.AlreadyGreenRepair)
	}
	if normalized := NormalizeResult(CommandResult{AlreadyGreenRepair: &AlreadyGreenRepairValidation{}}); normalized.AlreadyGreenRepair.GateCommand == nil || normalized.AlreadyGreenRepair.ReasonCodes == nil || normalized.AlreadyGreenRepair.OffendingPaths == nil {
		t.Fatalf("repair validation collections were not normalized: %#v", normalized.AlreadyGreenRepair)
	}
}

func TestAdvanceTextAndJSONCarryTheSameTypedLifecycleState(t *testing.T) {
	result := CommandResult{
		Command: "advance", Outcome: OutcomeSuccess, RepositoryRoot: "/tmp/example",
		Advance: &AdvanceLifecycleResult{
			RequestID: "REQ-503", RequestPath: "do-work/working/REQ-503-advance.md",
			TreeSection: "working", Status: "claimed", Route: "C",
			Phase: "qualify", PhaseKind: AdvancePhaseMechanical,
			MissingEvidence: []AdvanceMissingEvidence{{
				Kind: "section", Path: "do-work/working/REQ-503-advance.md",
				Section: "Qualification", Expected: "typed qualification result",
			}},
			NextArgv:         []string{"do-work-cli", "--format", "json", "qualify", "--request-path", "do-work/working/REQ-503-advance.md"},
			VerificationArgv: []string{"do-work-cli", "--format", "json", "advance", "REQ-503"},
			GateRecords: []AdvanceGateRecord{{
				RequestID: "REQ-503", RequestPath: "do-work/working/REQ-503-advance.md",
				GateID: "run-blocked-check", Provenance: AdvanceGateBaselineRecord,
				State: AdvanceGateFindings, Outcome: OutcomeFindings,
				FocusedTest: &FocusedTestResult{ExitStatus: 9, Launched: true, DiagnosticSHA256: "abc", BaselineState: FocusedBaselineNewRed},
			}},
		},
	}

	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"advance REQ-503 [working, claimed, route C]: qualify",
		"phase kind: mechanical",
		`missing section: do-work/working/REQ-503-advance.md section="Qualification" expected=typed qualification result`,
		"next: do-work-cli --format json qualify --request-path do-work/working/REQ-503-advance.md",
		"verify: do-work-cli --format json advance REQ-503",
		"gate run-blocked-check [baseline_record, findings]: findings",
		"focused test: status=9 baseline=new_red diagnostic=abc",
	} {
		if !strings.Contains(string(textOutput), expected) {
			t.Errorf("text output missing %q:\n%s", expected, textOutput)
		}
	}

	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Advance == nil || decoded.Advance.Phase != "qualify" || decoded.Advance.PhaseKind != AdvancePhaseMechanical ||
		!reflect.DeepEqual(decoded.Advance.NextArgv, result.Advance.NextArgv) || !reflect.DeepEqual(decoded.Advance.VerificationArgv, result.Advance.VerificationArgv) ||
		len(decoded.Advance.MissingEvidence) != 1 || decoded.Advance.MissingEvidence[0].Section != "Qualification" ||
		len(decoded.Advance.GateRecords) != 1 || decoded.Advance.GateRecords[0].FocusedTest == nil || decoded.Advance.GateRecords[0].FocusedTest.BaselineState != FocusedBaselineNewRed {
		t.Fatalf("JSON lost typed advance state: %#v", decoded.Advance)
	}

	normalized := NormalizeResult(CommandResult{Advance: &AdvanceLifecycleResult{}})
	if normalized.Advance.MissingEvidence == nil || normalized.Advance.NextArgv == nil || normalized.Advance.VerificationArgv == nil || normalized.Advance.GateRecords == nil {
		t.Fatalf("advance collections must normalize non-null: %#v", normalized.Advance)
	}
}

func TestRecoveryAndCheckpointTextAndJSONCarryTypedState(t *testing.T) {
	result := CommandResult{
		Command: "recover", Outcome: OutcomeSuccess,
		Recovery: &RecoveryResult{
			AuthorityMode: "take-over", TakeOverRequestID: "REQ-504", FinalizationPassed: true,
			Claims: []RecoveryClaimResult{{
				RequestID: "REQ-504", RequestPath: "do-work/working/REQ-504.md", Decision: "recovered", Recovered: true,
				CheckpointEvidence: []SelectionClaimEvidence{{Source: "checkpoint", Writer: "other:/repo"}},
			}},
			NextArgv: []string{"do-work-cli", "next"}, VerificationArgv: []string{"do-work-cli", "recover"},
		},
		Checkpoint: &CheckpointResult{CheckpointPath: "do-work/CHECKPOINT.md", PreservedClaims: 1, WrittenAt: "2026-09-04T12:00:00Z"},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"recovery [take-over]", "claim REQ-504", "checkpoint do-work/CHECKPOINT.md: preserved_claims=1"} {
		if !strings.Contains(string(textOutput), expected) {
			t.Fatalf("text output missing %q:\n%s", expected, textOutput)
		}
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Recovery == nil || len(decoded.Recovery.Claims) != 1 || decoded.Checkpoint == nil || decoded.Checkpoint.PreservedClaims != 1 {
		t.Fatalf("typed recovery/checkpoint state lost: %#v", decoded)
	}
	normalized := NormalizeResult(CommandResult{Recovery: &RecoveryResult{Claims: []RecoveryClaimResult{{}}}})
	if normalized.Recovery.NextArgv == nil || normalized.Recovery.VerificationArgv == nil || normalized.Recovery.Claims[0].CheckpointEvidence == nil {
		t.Fatalf("recovery collections must normalize non-null: %#v", normalized.Recovery)
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

func TestRefusalRemediesNeverNameTheInvokingCommand(t *testing.T) {
	tests := []struct {
		name        string
		result      CommandResult
		wantOutcome CommandOutcome
		wantNext    []string
	}{
		{
			name: "owned refusal becomes a set-aside",
			result: CommandResult{Command: "recover-finalization", Outcome: OutcomeRefused, Findings: []CommandFinding{{
				Code: "FINALIZATION-LIFECYCLE-APPLY", Fixability: FixabilityRefused, AffectedIDs: []string{"REQ-456"},
				AutomationStopReason: "lifecycle apply refused", NextArgv: []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"},
				VerificationArgv: []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"},
			}}},
			wantOutcome: OutcomeFindings,
			wantNext:    []string{},
		},
		{
			name: "different resolving verb is preserved",
			result: CommandResult{Command: "claim", Outcome: OutcomeRefused, Findings: []CommandFinding{{
				Code: "GIT-DIRTY-TARGET", AffectedIDs: []string{"REQ-513"},
				AutomationStopReason: "target is dirty", NextArgv: []string{"git", "diff", "--", "do-work/CHECKPOINT.md"},
				VerificationArgv: []string{"do-work-cli", "--format", "json", "claim", "REQ-513"},
			}}},
			wantOutcome: OutcomeRefused,
			wantNext:    []string{"git", "diff", "--", "do-work/CHECKPOINT.md"},
		},
		{
			name: "unowned self reference remains a global stop without a loop",
			result: CommandResult{Command: "cleanup", Outcome: OutcomeRefused, Findings: []CommandFinding{{
				Code: "CLEANUP-PREFLIGHT", Fixability: FixabilityRefused,
				AutomationStopReason: "shared cleanup target is dirty", NextArgv: []string{"do-work-cli", "cleanup", "--dry-run"},
				VerificationArgv: []string{"do-work-cli", "--format", "json", "cleanup", "--dry-run"},
			}}},
			wantOutcome: OutcomeRefused,
			wantNext:    []string{},
		},
		{
			name: "refusal finding inside findings outcome loses its self remedy",
			result: CommandResult{Command: "cleanup", Outcome: OutcomeFindings, Findings: []CommandFinding{{
				Code: "CLEANUP-GROUP-REFUSED", Fixability: FixabilityRefused, AffectedIDs: []string{"REQ-206"},
				AutomationStopReason: "scratch cannot participate in commit mode", NextArgv: []string{"do-work-cli", "cleanup"},
				VerificationArgv: []string{"git", "status", "--short"},
			}}},
			wantOutcome: OutcomeFindings,
			wantNext:    []string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			normalized := NormalizeResult(test.result)
			if normalized.Outcome != test.wantOutcome {
				t.Fatalf("outcome = %q, want %q", normalized.Outcome, test.wantOutcome)
			}
			if !reflect.DeepEqual(normalized.Findings[0].NextArgv, test.wantNext) {
				t.Fatalf("next argv = %#v, want %#v", normalized.Findings[0].NextArgv, test.wantNext)
			}
			if !reflect.DeepEqual(normalized.Findings[0].VerificationArgv, test.result.Findings[0].VerificationArgv) {
				t.Fatalf("verification argv changed: %#v", normalized.Findings[0].VerificationArgv)
			}
		})
	}
}

func TestRenderRejectsUnknownFormat(t *testing.T) {
	if _, err := RenderResult(CommandResult{Outcome: OutcomeSuccess}, OutputFormat("yaml")); err == nil {
		t.Fatal("RenderResult accepted unsupported output format")
	}
}

func TestProtocolOutputIsExactInTextAndTypedInJSON(t *testing.T) {
	exact := "hook bytes\nwithout generic prose\n"
	result := CommandResult{
		Command: "session-start", Outcome: OutcomeSuccess, ProtocolOutput: &exact,
		Findings: []CommandFinding{{Code: "HOOK-EVIDENCE", Severity: SeverityWarning}},
		Changes:  []RecordedChange{{Path: "do-work/.req-reservations/REQ-000001", Kind: "deleted", Detail: "matching committed request exists"}},
	}
	textOutput, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	if string(textOutput) != exact {
		t.Fatalf("text output = %q, want %q", textOutput, exact)
	}
	jsonOutput, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonOutput, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ProtocolOutput == nil || *decoded.ProtocolOutput != exact {
		t.Fatalf("JSON protocol output = %#v, want %q", decoded.ProtocolOutput, exact)
	}
	if len(decoded.Findings) != 1 || decoded.Findings[0].Code != "HOOK-EVIDENCE" || len(decoded.Changes) != 1 || decoded.Changes[0].Kind != "deleted" {
		t.Fatalf("JSON protocol projection lost typed operation evidence: %+v", decoded)
	}

	empty := ""
	textOutput, err = RenderResult(CommandResult{Command: "memory-stop-capture", Outcome: OutcomeSuccess, ProtocolOutput: &empty}, FormatText)
	if err != nil || len(textOutput) != 0 {
		t.Fatalf("explicit empty protocol output = %q, err=%v", textOutput, err)
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

func TestGateDeferralEvidenceHasTypedJSONAndTextParity(t *testing.T) {
	result := CommandResult{Command: "defer-gate", Outcome: OutcomeSuccess, GateDeferral: &GateDeferralResult{
		ParentID: "REQ-101", ParentPath: "do-work/queue/REQ-101.md", RepairID: "REQ-901", RepairPath: "do-work/queue/REQ-901.md",
		CheckpointPath: "do-work/CHECKPOINT.md", RepairOutcome: "folded", RepairDependency: "REQ-901",
		DiagnosticFingerprint: "sha256:red", SweepKey: "repository-gate-red", GateCommand: []string{"go", "test", "./..."}, GateExitStatus: 1,
	}}
	jsonBytes, err := RenderResult(result, FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var decoded CommandResult
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.GateDeferral == nil || decoded.GateDeferral.RepairOutcome != "folded" || !reflect.DeepEqual(decoded.GateDeferral.GateCommand, []string{"go", "test", "./..."}) {
		t.Fatalf("typed gate deferral = %#v", decoded.GateDeferral)
	}
	textBytes, err := RenderResult(result, FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"parent=REQ-101", "repair=REQ-901", "outcome=folded", "go test ./... (exit: 1)",
		"parent path: do-work/queue/REQ-101.md", "repair path: do-work/queue/REQ-901.md",
		"checkpoint path: do-work/CHECKPOINT.md", "repair dependency: REQ-901", "sweep key: repository-gate-red",
	} {
		if !strings.Contains(string(textBytes), required) {
			t.Fatalf("text omitted %q:\n%s", required, textBytes)
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

func TestExactTextAndAuditJSONShareOneResult(t *testing.T) {
	text := "## Inventory\n"
	result := CommandResult{Outcome: OutcomeSuccess, ExactTextOutput: &text, AuditMetrics: &AuditMetricsResult{Kind: "inventory"}}
	renderedText, err := RenderResult(result, FormatText)
	if err != nil || string(renderedText) != text {
		t.Fatalf("text=%q err=%v", renderedText, err)
	}
	renderedJSON, err := RenderResult(result, FormatJSON)
	var decoded CommandResult
	decodeErr := json.Unmarshal(renderedJSON, &decoded)
	if err != nil || decodeErr != nil || decoded.AuditMetrics == nil || decoded.AuditMetrics.Kind != "inventory" || strings.Contains(string(renderedJSON), "exact_text") {
		t.Fatalf("json=%s err=%v", renderedJSON, err)
	}
}
