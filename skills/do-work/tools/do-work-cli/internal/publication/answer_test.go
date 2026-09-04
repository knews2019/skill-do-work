package publication

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestRemediationF3TerminalFollowUpUsesExistingArchivedUR(t *testing.T) {
	root := t.TempDir()
	request := []byte("---\nid: REQ-2\nstatus: pending-answers\nuser_request: UR-1\nbuilder_decided: true\n---\n## Open Questions\n\n- [ ] Keep default?\n")
	writeFixture(t, root, "do-work/queue/REQ-2-follow-up.md", request, 0o644)
	writeFixture(t, root, "do-work/archive/UR-1/input.md", []byte("already archived\n"), 0o644)
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{
		RequestPath: "do-work/queue/REQ-2-follow-up.md", ExpectedStatus: "pending-answers", Mode: "clarify", CloseUserRequest: true,
		ArchivePath: "do-work/archive/UR-1/REQ-2-follow-up.md", UserRequestPath: "do-work/user-requests/UR-1", ArchiveDirectory: "do-work/archive/UR-1",
		Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Keep default?", Outcome: "confirmed", Summary: "keep default"}},
	}}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
	if plan.Refusal != nil {
		t.Fatalf("archived-UR fallback refused: %#v", plan.Refusal)
	}
	if len(plan.Mutations) != 1 || plan.Mutations[0].DestinationPath != "do-work/archive/UR-1/REQ-2-follow-up.md" {
		t.Fatalf("mutations = %#v", plan.Mutations)
	}
}

func TestRemediationF4StakeholderManifestExpressesCoupledEvidence(t *testing.T) {
	fixture := `{"operation":"answer","answer":{"request_path":"do-work/queue/REQ-1.md","expected_status":"blocked","mode":"stakeholder","answers":[{"question_id":"Q-01","outcome":"answered","summary":"blue"}],"report":{"path":"ai-reports/REQ-1/report.md","payload":{"source_path":"payload/report"}},"stakeholder_report":{"blocked_by":"ai-reports/REQ-1/report.md","reports_history":{"source_path":"payload/reports-history"}},"stakeholder_terminal":{"blocked_history":{"source_path":"payload/blocked-history"},"implementation":{"source_path":"payload/implementation"}}}}`
	if _, err := DecodeManifest(strings.NewReader(fixture), OperationAnswer); err != nil {
		t.Fatalf("typed stakeholder evidence rejected: %v", err)
	}
}

func TestRemediationF4StakeholderPartialAndTerminalPublishRequiredHistory(t *testing.T) {
	t.Run("partial report linkage", func(t *testing.T) {
		root := t.TempDir()
		request := []byte("---\nid: REQ-1\nstatus: blocked\nblocked_by: old-report.md\nblocked_at: 2026-08-01T00:00:00Z\n---\n## Open Questions\n\n- [ ] **Q-01** — Color?\n- [ ] **Q-02** — Size?\n")
		writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
		writeFixture(t, root, "payload/report", []byte("fresh report\n"), 0o644)
		writeFixture(t, root, "payload/reports-history", []byte("## Reports\n\n- ai-reports/REQ-1/fresh.md — partial stakeholder answers.\n"), 0o644)
		plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{
			RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "blocked", Mode: "stakeholder",
			Answers:           []QuestionAnswer{{QuestionID: "Q-01", Outcome: "answered", Summary: "blue"}},
			Report:            &PublishedFile{Path: "ai-reports/REQ-1/fresh.md", Payload: PayloadFile{SourcePath: "payload/report"}},
			StakeholderReport: &StakeholderReportEvidence{BlockedBy: "ai-reports/REQ-1/fresh.md", ReportsHistory: PayloadFile{SourcePath: "payload/reports-history"}},
		}}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
		if plan.Refusal != nil {
			t.Fatal(plan.Refusal)
		}
		var requestBytes []byte
		for _, mutation := range plan.Mutations {
			if mutation.Path == "do-work/queue/REQ-1-test.md" {
				requestBytes = mutation.Contents
			}
		}
		if !bytes.Contains(requestBytes, []byte("blocked_by: ai-reports/REQ-1/fresh.md")) || !bytes.Contains(requestBytes, []byte("## Reports")) {
			t.Fatalf("partial request bytes = %s", requestBytes)
		}
	})

	t.Run("terminal history", func(t *testing.T) {
		root := t.TempDir()
		request := []byte("---\nid: REQ-1\nstatus: blocked\nblocked_by: old-report.md\nblocked_at: 2026-08-01T00:00:00Z\n---\n## Open Questions\n\n- [ ] **Q-01** — Color?\n")
		writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
		writeFixture(t, root, "payload/blocked-history", []byte("## Blocked\n\n- Resolved 2026-09-01 after stakeholder answer.\n"), 0o644)
		writeFixture(t, root, "payload/implementation", []byte("## Implementation\n\nNo code changes.\n"), 0o644)
		plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{
			RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "blocked", Mode: "stakeholder", ArchivePath: "do-work/archive/REQ-1-test.md",
			Answers:             []QuestionAnswer{{QuestionID: "Q-01", Outcome: "answered", Summary: "blue"}},
			StakeholderTerminal: &StakeholderTerminalEvidence{BlockedHistory: PayloadFile{SourcePath: "payload/blocked-history"}, Implementation: PayloadFile{SourcePath: "payload/implementation"}},
		}}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
		if plan.Refusal != nil {
			t.Fatal(plan.Refusal)
		}
		contents := plan.Mutations[0].Contents
		if !bytes.Contains(contents, []byte("## Blocked")) || !bytes.Contains(contents, []byte("## Implementation")) || bytes.Contains(contents, []byte("blocked_by:")) {
			t.Fatalf("terminal request bytes = %s", contents)
		}
	})
}

func TestRemediationF5DelimiterSummaryRequiresMatchingRawPayload(t *testing.T) {
	root := t.TempDir()
	request := []byte("---\nid: REQ-1\nstatus: pending-answers\n---\n## Open Questions\n\n- [ ] Choice?\n")
	writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Choice?", Outcome: "answered", Summary: "## injected heading"}}}}, time.Now())
	if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-RAW-PAYLOAD-REQUIRED" {
		t.Fatalf("delimiter-shaped inline summary accepted: %#v", plan)
	}
	writeFixture(t, root, "payload/delimiter-answer", []byte("## injected heading"), 0o644)
	plan = BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Choice?", Outcome: "answered", Summary: "## injected heading", RawAnswer: &PayloadFile{SourcePath: "payload/delimiter-answer"}}}}}, time.Now())
	if plan.Refusal != nil {
		t.Fatalf("matching raw delimiter answer refused: %#v", plan.Refusal)
	}
	if bytes.Contains(plan.Mutations[0].Contents, []byte("→ ## injected heading")) || !bytes.Contains(plan.Mutations[0].Contents, []byte("> ## injected heading")) {
		t.Fatalf("delimiter answer not safely contained: %s", plan.Mutations[0].Contents)
	}
}

func TestRemediationF6OverrideUsesStructuredCaptureManifest(t *testing.T) {
	fixture := `{"operation":"answer","answer":{"request_path":"do-work/queue/REQ-1.md","expected_status":"blocked","mode":"stakeholder","answers":[{"question_id":"Q-01","outcome":"answered","summary":"override"}],"override_capture":{"user_request_id":"UR-2","user_request":{"path":"do-work/user-requests/UR-2/input.md","payload":{"source_path":"payload/ur"}},"raw_input":{"source_path":"payload/raw"},"requests":[{"id":"REQ-2","user_request_id":"UR-2","file":{"path":"do-work/queue/REQ-2-change.md","payload":{"source_path":"payload/req"}},"reservation_path":"do-work/.req-reservations/REQ-002"}]}}}`
	if _, err := DecodeManifest(strings.NewReader(fixture), OperationAnswer); err != nil {
		t.Fatalf("structured override capture rejected: %v", err)
	}
}

func TestRemediationF6StructuredOverridePlansAndRollsBackWithAnswer(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init", "-q")
	request := []byte("---\nid: REQ-1\nstatus: blocked\n---\n## Open Questions\n\n- [ ] **Q-01** — Override?\n")
	writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
	writeFixture(t, root, "seed", []byte("seed\n"), 0o644)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "seed")
	raw := []byte("stakeholder changed scope")
	writeFixture(t, root, "payload/raw", raw, 0o644)
	urPayload := append(canonicalURFixture("UR-2", []string{"REQ-2"}), containedOutsideBytes(raw, "\n")...)
	urPayload = append(urPayload, '\n')
	writeFixture(t, root, "payload/ur", urPayload, 0o644)
	writeFixture(t, root, "payload/req", canonicalREQFixture("REQ-2", "UR-2"), 0o644)
	writeFixture(t, root, "payload/blocked-history", []byte("## Blocked\n\n- Resolved after stakeholder override.\n"), 0o644)
	writeFixture(t, root, "payload/implementation", []byte("## Implementation\n\nNo code changes.\n"), 0o644)
	override := &CaptureManifest{UserRequestID: "UR-2", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-2/input.md", Payload: PayloadFile{SourcePath: "payload/ur"}}, RawInput: &PayloadFile{SourcePath: "payload/raw"}, Requests: []CaptureRequest{{ID: "REQ-2", UserRequestID: "UR-2", File: PublishedFile{Path: "do-work/queue/REQ-2-change.md", Payload: PayloadFile{SourcePath: "payload/req"}}, ReservationPath: "do-work/.req-reservations/REQ-002"}}}
	answer := &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "blocked", Mode: "stakeholder", ArchivePath: "do-work/archive/REQ-1-test.md", Answers: []QuestionAnswer{{QuestionID: "Q-01", Outcome: "answered", Summary: "override"}}, OverrideCapture: override, StakeholderTerminal: &StakeholderTerminalEvidence{BlockedHistory: PayloadFile{SourcePath: "payload/blocked-history"}, Implementation: PayloadFile{SourcePath: "payload/implementation"}}}
	noncanonicalOverride := *override
	noncanonicalOverride.Requests = append([]CaptureRequest(nil), override.Requests...)
	noncanonicalOverride.Requests[0].ReservationPath = "do-work/.req-reservations/REQ-000002"
	noncanonicalAnswer := *answer
	noncanonicalAnswer.OverrideCapture = &noncanonicalOverride
	noncanonicalPlan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &noncanonicalAnswer}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
	if noncanonicalPlan.Refusal == nil || noncanonicalPlan.Refusal.Code != "ANSWER-OVERRIDE-CAPTURE-CAPTURE-RESERVATION-MISMATCH" {
		t.Fatalf("noncanonical override reservation accepted: %#v", noncanonicalPlan.Refusal)
	}
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: answer}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
	if plan.Refusal != nil {
		t.Fatal(plan.Refusal)
	}
	previous := beforePublicationMutation
	beforePublicationMutation = func(index int, _ PlannedMutation) error {
		if index == len(plan.Mutations)-1 {
			return errors.New("override rollback seam")
		}
		return nil
	}
	t.Cleanup(func() { beforePublicationMutation = previous })
	result := ApplyPlan(t.Context(), plan, false, false)
	if result.Outcome != resultmodel.OutcomeRolledBack || result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Fatalf("rollback result = %#v", result)
	}
	got, err := os.ReadFile(filepath.Join(root, "do-work/queue/REQ-1-test.md"))
	if err != nil || !bytes.Equal(got, request) {
		t.Fatalf("answer target not restored: %v %q", err, got)
	}
	for _, path := range []string{"do-work/archive/REQ-1-test.md", "do-work/queue/REQ-2-change.md", "do-work/user-requests/UR-2/input.md", "do-work/.req-reservations/REQ-002"} {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Fatalf("override path survived rollback %s: %v", path, err)
		}
	}
}

func TestRemediationF6GenericOverrideCreatesAreRefused(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "do-work/queue/REQ-1-test.md", []byte("---\nid: REQ-1\nstatus: pending-answers\n---\n- [ ] Choice?\n"), 0o644)
	writeFixture(t, root, "payload/arbitrary", []byte("arbitrary"), 0o644)
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Choice?", Outcome: "answered", Summary: "yes"}}, OverrideCreates: []PublishedFile{{Path: "arbitrary/file", Payload: PayloadFile{SourcePath: "payload/arbitrary"}}}}}, time.Now())
	if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-OVERRIDE-UNSTRUCTURED" {
		t.Fatalf("generic override create accepted: %#v", plan.Refusal)
	}
}

func TestBuildAnswerPlanUsesWholeRecordDispositionAndContainsRawBytes(t *testing.T) {
	root := t.TempDir()
	request := []byte("---\nid: REQ-1\nstatus: pending-answers\nbuilder_decided: true\n---\n## Open Questions\n\n- [ ] First choice?\n- [ ] Second choice?\n")
	writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
	writeFixture(t, root, "payload/answer", []byte("line one\n## not a heading\n````"), 0o600)
	manifest := Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] First choice?", Outcome: "confirmed", Summary: "keep it", RawAnswer: &PayloadFile{SourcePath: "payload/answer"}}}}}
	plan := BuildAnswerPlan(root, manifest, time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC))
	if plan.Refusal != nil {
		t.Fatalf("refusal = %#v", plan.Refusal)
	}
	if len(plan.Mutations) != 1 || !bytes.Contains(plan.Mutations[0].Contents, []byte("status: pending-answers")) || !bytes.Contains(plan.Mutations[0].Contents, []byte("> ## not a heading")) {
		t.Fatalf("mutation = %#v", plan.Mutations)
	}
}

func TestBuildAnswerPlanTerminalClosureInventoriesAssets(t *testing.T) {
	root := t.TempDir()
	request := []byte("---\nid: REQ-1\nstatus: blocked\nuser_request: UR-1\n---\n- [ ] **Q-01** — Color?\n")
	writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
	writeFixture(t, root, "do-work/user-requests/UR-1/input.md", []byte("input"), 0o644)
	writeFixture(t, root, "do-work/user-requests/UR-1/assets/pic.bin", []byte{1, 2, 3}, 0o600)
	writeFixture(t, root, "payload/blocked-history", []byte("## Blocked\n\n- Resolved after stakeholder answer.\n"), 0o644)
	writeFixture(t, root, "payload/implementation", []byte("## Implementation\n\nNo code changes.\n"), 0o644)
	manifest := Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "blocked", Mode: "stakeholder", Answers: []QuestionAnswer{{QuestionID: "Q-01", Outcome: "answered", Summary: "blue"}}, ArchivePath: "do-work/archive/UR-1/REQ-1-test.md", CloseUserRequest: true, UserRequestPath: "do-work/user-requests/UR-1", ArchiveDirectory: "do-work/archive/UR-1"}}
	manifest.Answer.StakeholderTerminal = &StakeholderTerminalEvidence{BlockedHistory: PayloadFile{SourcePath: "payload/blocked-history"}, Implementation: PayloadFile{SourcePath: "payload/implementation"}}
	plan := BuildAnswerPlan(root, manifest, time.Now())
	if plan.Refusal != nil {
		t.Fatalf("refusal = %#v", plan.Refusal)
	}
	foundAsset := false
	for _, mutation := range plan.Mutations {
		if mutation.Path == "do-work/user-requests/UR-1/assets/pic.bin" {
			foundAsset = true
		}
	}
	if !foundAsset {
		t.Fatalf("asset absent from mutations: %#v", plan.Mutations)
	}
	if _, err := os.Stat(filepath.Join(root, "do-work/user-requests/UR-1/assets/pic.bin")); err != nil {
		t.Fatal(err)
	}
}

func TestBuildAnswerPlanDerivesClarifyDispositionFromWholeRecord(t *testing.T) {
	tests := []struct {
		name       string
		builder    bool
		questions  string
		answers    []QuestionAnswer
		wantStatus string
		archive    bool
	}{
		{"open wins", true, "- [ ] First?\n- [ ] Second?\n", []QuestionAnswer{{ExpectedLine: "- [ ] First?", Outcome: "discarded", Summary: "not needed"}}, "pending-answers", false},
		{"all discarded", false, "- [ ] First?\n", []QuestionAnswer{{ExpectedLine: "- [ ] First?", Outcome: "discarded", Summary: "not needed"}}, "cancelled", true},
		{"builder confirmed", true, "- [ ] First?\n", []QuestionAnswer{{ExpectedLine: "- [ ] First?", Outcome: "confirmed", Summary: "use default"}}, "completed", true},
		{"ordinary answered", false, "- [ ] First?\n", []QuestionAnswer{{ExpectedLine: "- [ ] First?", Outcome: "answered", Summary: "choose blue"}}, "pending", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			builder := ""
			if test.builder {
				builder = "builder_decided: true\n"
			}
			request := []byte("---\nid: REQ-1\nstatus: pending-answers\n" + builder + "---\n## Open Questions\n\n" + test.questions)
			writeFixture(t, root, "do-work/queue/REQ-1-test.md", request, 0o644)
			answer := &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: test.answers}
			if test.archive {
				answer.ArchivePath = "do-work/archive/REQ-1-test.md"
			}
			plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: answer}, time.Date(2026, 8, 31, 1, 2, 3, 0, time.UTC))
			if plan.Refusal != nil {
				t.Fatalf("refusal=%#v", plan.Refusal)
			}
			contents := plan.Mutations[0].Contents
			if !bytes.Contains(contents, []byte("status: "+test.wantStatus)) {
				t.Fatalf("contents=%s", contents)
			}
			if !bytes.Contains(contents, []byte("## Answer Notes")) {
				t.Fatal("dated note missing")
			}
		})
	}
}

func TestHeavyTestingAnswerCompletesOnGreenAndRequeuesOnFailure(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		outcome      string
		summary      string
		wantVerified bool
	}{
		{name: "green", outcome: "confirmed", summary: "yes, exit 0", wantVerified: true},
		{name: "failed", outcome: "answered", summary: "no, browser lane failed"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			runGitFixture(t, root, "init", "-q")
			writeFixture(t, root, "implementation.txt", []byte("implementation\n"), 0o644)
			runGitFixture(t, root, "add", ".")
			runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "implementation")
			targetRevision := strings.TrimSpace(runGitFixtureOutput(t, root, "rev-parse", "HEAD"))
			writeFixture(t, root, "execution.txt", []byte("execution\n"), 0o644)
			runGitFixture(t, root, "add", ".")
			runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "execution")
			executionRevision := strings.TrimSpace(runGitFixtureOutput(t, root, "rev-parse", "HEAD"))
			requestPath := "do-work/queue/REQ-3-heavy.md"
			question := "- [ ] Run heavy tests at " + targetRevision + "; did they exit 0?"
			request := []byte("---\nid: REQ-3\nstatus: pending-heavy-testing\nclaimed_at: 2026-09-03T18:00:00Z\ncommit: " + targetRevision + "\n---\n## Open Questions\n\n" + question + "\n")
			writeFixture(t, root, requestPath, request, 0o644)
			laneExit := 1
			if testCase.wantVerified {
				laneExit = 0
			}
			answer := &AnswerManifest{
				RequestPath: requestPath, ExpectedStatus: "pending-heavy-testing", Mode: "heavy-testing",
				Answers:      []QuestionAnswer{{ExpectedLine: question, Outcome: testCase.outcome, Summary: testCase.summary}},
				HeavyTesting: &HeavyTestingEvidence{TargetRevision: targetRevision, ExecutionRevision: executionRevision, Lanes: []HeavyLaneResult{{LaneID: "browser", CommandArgv: []string{"bash", "verify.sh", "browser"}, ExitStatus: laneExit, WallSeconds: 26}}},
			}
			plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: answer}, time.Date(2026, 9, 3, 20, 0, 0, 0, time.UTC))
			if plan.Refusal != nil || len(plan.Mutations) != 1 {
				t.Fatalf("heavy-testing plan = %#v", plan)
			}
			if !bytes.Contains(plan.Mutations[0].Contents, []byte("status: pending")) {
				t.Fatalf("heavy-testing contents = %s", plan.Mutations[0].Contents)
			}
			if plan.Mutations[0].DestinationPath != "" || bytes.Contains(plan.Mutations[0].Contents, []byte("claimed_at:")) {
				t.Fatalf("heavy-testing should requeue without a stale claim: %#v\n%s", plan.Mutations[0], plan.Mutations[0].Contents)
			}
			if testCase.wantVerified {
				for _, expected := range []string{"heavy_verified_at: 2026-09-03T20:00:00Z", "heavy_verified_revision: " + executionRevision, "## Heavy Verification Result", "- browser: exit 0, 26s — "} {
					if !bytes.Contains(plan.Mutations[0].Contents, []byte(expected)) {
						t.Errorf("green result omitted %q:\n%s", expected, plan.Mutations[0].Contents)
					}
				}
			} else {
				if bytes.Contains(plan.Mutations[0].Contents, []byte("heavy_verified_")) {
					t.Fatalf("failed result retained green evidence:\n%s", plan.Mutations[0].Contents)
				}
				if !bytes.Contains(plan.Mutations[0].Contents, []byte("- browser: exit 1, 26s — ")) {
					t.Errorf("failed result omitted the red lane's wall time:\n%s", plan.Mutations[0].Contents)
				}
			}
		})
	}
}

// newHeavyAnswerFixture builds a two-commit repository holding one
// pending-heavy-testing request whose commit is the first revision, and
// returns the pieces a heavy-testing answer manifest needs.
func newHeavyAnswerFixture(t *testing.T, requestID string) (root, requestPath, question, targetRevision, executionRevision string) {
	t.Helper()
	root = t.TempDir()
	runGitFixture(t, root, "init", "-q")
	writeFixture(t, root, "implementation.txt", []byte("implementation\n"), 0o644)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "implementation")
	targetRevision = strings.TrimSpace(runGitFixtureOutput(t, root, "rev-parse", "HEAD"))
	writeFixture(t, root, "execution.txt", []byte("execution\n"), 0o644)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "execution")
	executionRevision = strings.TrimSpace(runGitFixtureOutput(t, root, "rev-parse", "HEAD"))
	requestPath = "do-work/queue/" + requestID + "-heavy.md"
	question = "- [ ] Run heavy tests at " + targetRevision + "; did they exit 0?"
	request := []byte("---\nid: " + requestID + "\nstatus: pending-heavy-testing\nclaimed_at: 2026-09-03T18:00:00Z\ncommit: " + targetRevision + "\n---\n## Open Questions\n\n" + question + "\n")
	writeFixture(t, root, requestPath, request, 0o644)
	return root, requestPath, question, targetRevision, executionRevision
}

func TestHeavyTestingAnswerRejectsNegativeWallSeconds(t *testing.T) {
	root, requestPath, question, targetRevision, executionRevision := newHeavyAnswerFixture(t, "REQ-5")
	answer := &AnswerManifest{
		RequestPath: requestPath, ExpectedStatus: "pending-heavy-testing", Mode: "heavy-testing",
		Answers:      []QuestionAnswer{{ExpectedLine: question, Outcome: "confirmed", Summary: "green"}},
		HeavyTesting: &HeavyTestingEvidence{TargetRevision: targetRevision, ExecutionRevision: executionRevision, Lanes: []HeavyLaneResult{{LaneID: "browser", CommandArgv: []string{"bash", "verify.sh", "browser"}, ExitStatus: 0, WallSeconds: -1}}},
	}
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: answer}, time.Now())
	if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-HEAVY-EVIDENCE-INVALID" {
		t.Fatalf("negative wall seconds refusal = %#v", plan.Refusal)
	}
}

// TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane pins the answer half of
// R1: a lane that did not execute can never stamp heavy_verified_*, and the
// recorded evidence still carries the skip's wall time.
func TestHeavyTestingAnswerRefusesConfirmedWithSkippedLane(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		outcome    string
		wantCode   string
		wantResult string
	}{
		{name: "confirmed", outcome: "confirmed", wantCode: "ANSWER-HEAVY-EVIDENCE-NOT-GREEN"},
		{name: "answered", outcome: "answered", wantResult: "- browser: skipped, 1s — "},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root, requestPath, question, targetRevision, executionRevision := newHeavyAnswerFixture(t, "REQ-6")
			answer := &AnswerManifest{
				RequestPath: requestPath, ExpectedStatus: "pending-heavy-testing", Mode: "heavy-testing",
				Answers:      []QuestionAnswer{{ExpectedLine: question, Outcome: testCase.outcome, Summary: "browser lane skipped"}},
				HeavyTesting: &HeavyTestingEvidence{TargetRevision: targetRevision, ExecutionRevision: executionRevision, Lanes: []HeavyLaneResult{{LaneID: "browser", CommandArgv: []string{"bash", "verify.sh", "browser"}, ExitStatus: 0, Skipped: true, WallSeconds: 1}}},
			}
			plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: answer}, time.Now())
			if testCase.wantCode != "" {
				if plan.Refusal == nil || plan.Refusal.Code != testCase.wantCode {
					t.Fatalf("skipped-lane refusal = %#v", plan.Refusal)
				}
				return
			}
			if plan.Refusal != nil || len(plan.Mutations) != 1 {
				t.Fatalf("skipped-lane plan = %#v", plan)
			}
			if !bytes.Contains(plan.Mutations[0].Contents, []byte(testCase.wantResult)) {
				t.Fatalf("skipped-lane result omitted %q:\n%s", testCase.wantResult, plan.Mutations[0].Contents)
			}
			if bytes.Contains(plan.Mutations[0].Contents, []byte("heavy_verified_")) {
				t.Fatalf("skipped-lane result stamped green evidence:\n%s", plan.Mutations[0].Contents)
			}
		})
	}
}

func TestHeavyTestingAnswerRejectsMismatchedEvidence(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init", "-q")
	writeFixture(t, root, "seed", []byte("seed\n"), 0o644)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "seed")
	targetRevision := strings.TrimSpace(runGitFixtureOutput(t, root, "rev-parse", "HEAD"))
	writeFixture(t, root, "later", []byte("later\n"), 0o644)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "later")
	laterRevision := strings.TrimSpace(runGitFixtureOutput(t, root, "rev-parse", "HEAD"))
	requestPath := "do-work/queue/REQ-4-heavy.md"
	question := "- [ ] Run heavy tests; did they exit 0?"
	writeFixture(t, root, requestPath, []byte("---\nid: REQ-4\nstatus: pending-heavy-testing\ncommit: "+targetRevision+"\n---\n## Open Questions\n\n"+question+"\n"), 0o644)
	answer := &AnswerManifest{RequestPath: requestPath, ExpectedStatus: "pending-heavy-testing", Mode: "heavy-testing", Answers: []QuestionAnswer{{ExpectedLine: question, Outcome: "confirmed", Summary: "green"}}, HeavyTesting: &HeavyTestingEvidence{TargetRevision: laterRevision, ExecutionRevision: laterRevision, Lanes: []HeavyLaneResult{{LaneID: "lane", CommandArgv: []string{"true"}, ExitStatus: 0}}}}
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: answer}, time.Now())
	if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-HEAVY-EVIDENCE-MISMATCH" {
		t.Fatalf("mismatched evidence refusal = %#v", plan.Refusal)
	}
}

func runGitFixtureOutput(t *testing.T, repositoryRoot string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(arguments, " "), err, output)
	}
	return string(output)
}

func TestBuildAnswerPlanRefusesArchivedReply(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "do-work/archive/REQ-1-test.md", []byte("---\nid: REQ-1\nstatus: completed\n---\n- [ ] Old?\n"), 0o644)
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/archive/REQ-1-test.md", ExpectedStatus: "completed", Mode: "stakeholder", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Old?", Outcome: "answered", Summary: "late"}}}}, time.Now())
	if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-ARCHIVED-READ-ONLY" {
		t.Fatalf("refusal=%#v", plan.Refusal)
	}
}

// TestSummaryContainmentDecidesByMarkdownStructureCondition pins the condition, not a prefix
// list: a one-line summary needs canonical containment exactly when a Markdown reader could
// take it for the document's own structure. Each row names the structural class it stands
// for — the strings are fixtures for that class, never the definition of the accepted set.
func TestSummaryContainmentDecidesByMarkdownStructureCondition(t *testing.T) {
	structuralClasses := []struct {
		class   string
		summary string
	}{
		{"ATX heading", "## injected heading"},
		{"setext heading underline, equals run", "==="},
		{"setext heading underline, single equals", "="},
		{"thematic break, hyphens", "---"},
		{"thematic break, asterisks", "***"},
		{"thematic break, underscores", "___"},
		{"thematic break, spaced asterisks", "* * *"},
		{"thematic break, spaced hyphens", "- - -"},
		{"thematic break, long hyphen run", "-----"},
		{"frontmatter fence, TOML form", "+++"},
		{"fenced code, backticks", "```go"},
		{"fenced code, tildes", "~~~"},
		{"block quote", "> quoted answer"},
		{"bullet list marker, hyphen", "- first option"},
		{"bare bullet list marker, hyphen", "-"},
		{"bare bullet list marker, asterisk", "*"},
		{"bare bullet list marker, plus", "+"},
		{"task list checkbox", "- [ ] a question the board would adopt"},
		{"ordered list marker, period", "1. first option"},
		{"bare ordered list marker, period", "1."},
		{"bare ordered list marker, paren", "1)"},
		{"ordered list marker, multi-digit", "1234."},
		{"HTML block", "<div>"},
		{"link reference definition", "[label]: https://example.test"},
		{"footnote definition", "[^1]: a footnote"},
		{"table row", "| option | cost |"},
		{"table delimiter row", "|---|---|"},
		{"indented code block, spaces", "    four spaces of code"},
		{"indented block, tab", "\ttab-indented"},
		{"directive or admonition fence", ":::note"},
		{"math block fence", "$$"},
	}
	proseClasses := []struct {
		class   string
		summary string
	}{
		{"plain prose", "keep the default"},
		{"prose with an internal hyphen", "use the x-height metric"},
		{"prose opening with a digit", "42 requests remain"},
		{"prose opening with a digit and a later period", "3 files changed. ship it"},
		{"prose carrying a delimiter mid-line", "use --- as the separator"},
		{"prose carrying backticks mid-line", "run `go test` first"},
		{"prose whose first word carries a special character", "don't fence it in"},
		{"non-ASCII prose", "Überschrift bleibt"},
	}
	for _, structural := range structuralClasses {
		t.Run("structural/"+structural.class, func(t *testing.T) {
			if !summaryRequiresContainment(structural.summary) {
				t.Fatalf("%s accepted inline: %q", structural.class, structural.summary)
			}
		})
	}
	for _, prose := range proseClasses {
		t.Run("prose/"+prose.class, func(t *testing.T) {
			if summaryRequiresContainment(prose.summary) {
				t.Fatalf("%s forced into containment: %q", prose.class, prose.summary)
			}
		})
	}
}

// TestBuildAnswerPlanCarriesStructuralSummaryContainedAndProseInline proves the condition at
// the caller seam for a class the earlier prefix list missed, and that neither branch loses a
// byte: the structural summary lands contained and verbatim, the prose summary lands inline.
func TestBuildAnswerPlanCarriesStructuralSummaryContainedAndProseInline(t *testing.T) {
	t.Run("thematic break refuses an inline summary", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "do-work/queue/REQ-1-test.md", []byte("---\nid: REQ-1\nstatus: pending-answers\n---\n## Open Questions\n\n- [ ] Choice?\n"), 0o644)
		plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Choice?", Outcome: "answered", Summary: "***"}}}}, time.Now())
		if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-RAW-PAYLOAD-REQUIRED" {
			t.Fatalf("thematic-break summary accepted inline: %#v", plan.Refusal)
		}
	})

	t.Run("thematic break lands contained and lossless", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "do-work/queue/REQ-1-test.md", []byte("---\nid: REQ-1\nstatus: pending-answers\n---\n## Open Questions\n\n- [ ] Choice?\n"), 0o644)
		writeFixture(t, root, "payload/thematic-break", []byte("***"), 0o644)
		plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Choice?", Outcome: "answered", Summary: "***", RawAnswer: &PayloadFile{SourcePath: "payload/thematic-break"}}}}}, time.Now())
		if plan.Refusal != nil {
			t.Fatalf("matching raw payload refused: %#v", plan.Refusal)
		}
		contents := plan.Mutations[0].Contents
		if !bytes.Contains(contents, []byte("→ See contained answer note")) {
			t.Fatalf("structural summary not replaced by the contained label: %s", contents)
		}
		if bytes.Contains(contents, []byte("→ ***")) || !bytes.Contains(contents, []byte("> ***")) {
			t.Fatalf("thematic break not safely contained: %s", contents)
		}
	})

	t.Run("prose stays inline and verbatim", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "do-work/queue/REQ-1-test.md", []byte("---\nid: REQ-1\nstatus: pending-answers\n---\n## Open Questions\n\n- [ ] Choice?\n"), 0o644)
		summary := "use the x-height metric, not --- as a separator"
		plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "pending-answers", Mode: "clarify", Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Choice?", Outcome: "answered", Summary: summary}}}}, time.Now())
		if plan.Refusal != nil {
			t.Fatalf("plain prose summary refused: %#v", plan.Refusal)
		}
		if !bytes.Contains(plan.Mutations[0].Contents, []byte("→ "+summary)) {
			t.Fatalf("prose summary not written inline verbatim: %s", plan.Mutations[0].Contents)
		}
	})
}

// clarifyRequestFixturePath is the queue path every resolved-disposition case below answers
// against; the cases differ only in what an earlier round left on the question lines.
const clarifyRequestFixturePath = "do-work/queue/REQ-1-test.md"

func writeClarifyFixture(t *testing.T, root string, builderDecided bool, questionSection string) {
	t.Helper()
	builderMarker := ""
	if builderDecided {
		builderMarker = "builder_decided: true\n"
	}
	writeFixture(t, root, clarifyRequestFixturePath, []byte("---\nid: REQ-1\nstatus: pending-answers\n"+builderMarker+"---\n## Open Questions\n\n"+questionSection), 0o644)
}

// buildClarifyRound answers one clarify round against the current fixture bytes. Chaining two
// rounds through the writer keeps the earlier round's resolved line exactly as production
// composes it instead of a hand-authored imitation of it.
func buildClarifyRound(t *testing.T, root string, answers ...QuestionAnswer) PublicationPlan {
	t.Helper()
	plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{
		RequestPath: clarifyRequestFixturePath, ExpectedStatus: "pending-answers", Mode: "clarify", Answers: answers,
	}}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
	if plan.Refusal != nil {
		t.Fatalf("clarify round refused: %#v", plan.Refusal)
	}
	return plan
}

// rewriteFixtureWithCarriageReturns republishes the current fixture with CRLF line endings, so
// a case can exercise the format whose trailing carriage return the readers trim.
func rewriteFixtureWithCarriageReturns(t *testing.T, root string) {
	t.Helper()
	current, readError := os.ReadFile(filepath.Join(root, filepath.FromSlash(clarifyRequestFixturePath)))
	if readError != nil {
		t.Fatalf("fixture unreadable: %v", readError)
	}
	writeFixture(t, root, clarifyRequestFixturePath, bytes.ReplaceAll(current, []byte("\n"), []byte("\r\n")), 0o644)
}

// requireClarifyStatus checks the status a round would publish together with the two effects
// that status decides: a terminal status stamps completed_at and moves the REQ into the
// archive, a non-terminal one leaves it in the queue with no completion timestamp.
func requireClarifyStatus(t *testing.T, plan PublicationPlan, wantStatus string) {
	t.Helper()
	mutation := plan.Mutations[0]
	// A CRLF document publishes CRLF frontmatter, so the status line is matched against
	// line-ending-normalized bytes; every other assertion here is ending-independent.
	normalized := bytes.ReplaceAll(mutation.Contents, []byte("\r\n"), []byte("\n"))
	if !bytes.Contains(normalized, []byte("status: "+wantStatus+"\n")) {
		t.Fatalf("want status %q, published document = %s", wantStatus, mutation.Contents)
	}
	terminal := wantStatus == "cancelled" || wantStatus == "completed"
	if terminal != bytes.Contains(mutation.Contents, []byte("completed_at:")) {
		t.Fatalf("completed_at presence does not match status %q: %s", wantStatus, mutation.Contents)
	}
	if terminal != (mutation.Kind == MutationMove) || terminal != strings.HasPrefix(mutation.DestinationPath, "do-work/archive/") {
		t.Fatalf("archive disposition does not match status %q: %#v", wantStatus, mutation)
	}
}

// TestBuildAnswerPlanRefusesDispositionForgedByAnswerText is REQ-528's lock-in. An answer
// summary is user text and lands immediately after the " → " separator the writer appended, so
// a summary reading "keep it → Discarded: not really" puts the marker bytes on a line whose
// question was answered, not discarded. Matching the marker anywhere on the line turns that
// text into a terminal status and archives the REQ; only the writer's own position counts.
func TestBuildAnswerPlanRefusesDispositionForgedByAnswerText(t *testing.T) {
	t.Run("discarded marker in an answer summary does not cancel the REQ", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, false, "- [ ] Keep the flag?\n- [ ] Drop the table?\n")
		firstRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Keep the flag?", Outcome: "answered", Summary: "keep it → Discarded: not really"})
		writeFixture(t, root, clarifyRequestFixturePath, firstRound.Mutations[0].Contents, 0o644)
		secondRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Drop the table?", Outcome: "discarded", Summary: "not needed"})
		requireClarifyStatus(t, secondRound, "pending")
	})

	t.Run("confirmed marker in an answer summary does not complete the REQ", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, true, "- [ ] Keep the flag?\n- [ ] Drop the table?\n")
		firstRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Keep the flag?", Outcome: "answered", Summary: "keep it → Confirmed: not really"})
		writeFixture(t, root, clarifyRequestFixturePath, firstRound.Mutations[0].Contents, 0o644)
		secondRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Drop the table?", Outcome: "confirmed", Summary: "use the default"})
		requireClarifyStatus(t, secondRound, "pending")
	})
}

// TestBuildAnswerPlanKeepsGenuineMultiRoundDispositionsTerminal is the other half of the same
// contract: a genuinely uniform disposition spread over two rounds must still reach its
// terminal status, so the position fix cannot be a blanket refusal of prior-round lines.
func TestBuildAnswerPlanKeepsGenuineMultiRoundDispositionsTerminal(t *testing.T) {
	t.Run("every question discarded across two rounds still cancels", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, false, "- [ ] Keep the flag?\n- [ ] Drop the table?\n")
		firstRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Keep the flag?", Outcome: "discarded", Summary: "the flag is gone"})
		writeFixture(t, root, clarifyRequestFixturePath, firstRound.Mutations[0].Contents, 0o644)
		secondRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Drop the table?", Outcome: "discarded", Summary: "not needed"})
		requireClarifyStatus(t, secondRound, "cancelled")
	})

	t.Run("every question confirmed across two rounds still completes", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, true, "- [ ] Keep the flag?\n- [ ] Drop the table?\n")
		firstRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Keep the flag?", Outcome: "confirmed", Summary: "keep it"})
		writeFixture(t, root, clarifyRequestFixturePath, firstRound.Mutations[0].Contents, 0o644)
		secondRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Drop the table?", Outcome: "confirmed", Summary: "use the default"})
		requireClarifyStatus(t, secondRound, "completed")
	})
}

// TestBuildAnswerPlanJudgesPriorRoundResolvedLinesAtTheWriterPosition covers the lines this
// invocation did not write. A question resolved in an earlier round is already "- [x]" in the
// body, and clarify's own contract has humans editing these files, so the line can carry
// anything. Each row supplies one prior line and discards the one remaining question, which
// makes the published status the whole verdict for that line.
func TestBuildAnswerPlanJudgesPriorRoundResolvedLinesAtTheWriterPosition(t *testing.T) {
	tests := []struct {
		name       string
		priorLine  string
		wantStatus string
		carriage   bool
	}{
		{"disposition at the writer position is the real one", "- [x] Keep the flag? → Discarded: not needed", "cancelled", false},
		{"summary merely mentioning the marker is not a disposition", "- [x] Keep the flag? → keep it → Discarded: not really", "pending", false},
		{"answered line carries no disposition at all", "- [x] Keep the flag? → keep it", "pending", false},
		{"resolved line carrying no separator leaves no position to read", "- [x] Keep the flag?", "pending", false},
		{"question text carrying the separator leaves no identifiable position", "- [x] Cancel → Discarded: yes? → Discarded: not needed", "pending", false},
		{"CRLF document still reads the disposition at the writer position", "- [x] Keep the flag? → Discarded: not needed", "cancelled", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeClarifyFixture(t, root, false, test.priorLine+"\n- [ ] Drop the table?\n")
			if test.carriage {
				rewriteFixtureWithCarriageReturns(t, root)
			}
			round := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Drop the table?", Outcome: "discarded", Summary: "not needed"})
			requireClarifyStatus(t, round, test.wantStatus)
		})
	}
}

// TestBuildAnswerPlanRefusesAnAnsweredSummaryOpeningWithADispositionLabel closes the half of
// REQ-528 the position anchor cannot reach. An answered summary lands at exactly the position
// the writer would place "Confirmed: " or "Discarded: ", and the published line records nothing
// about which of the two wrote those bytes — so a summary that opens with a disposition label
// is a disposition as far as every later reader is concerned, and a following round that
// discards the last open question cancels and archives the REQ. Position anchoring cannot
// separate them; only the write side can, by refusing the collision it would create.
func TestBuildAnswerPlanRefusesAnAnsweredSummaryOpeningWithADispositionLabel(t *testing.T) {
	for _, labelPrefix := range []string{discardedLabelPrefix, confirmedLabelPrefix} {
		t.Run("answered summary opening with "+strconv.Quote(labelPrefix), func(t *testing.T) {
			root := t.TempDir()
			writeClarifyFixture(t, root, true, "- [ ] Keep the flag?\n- [ ] Drop the table?\n")
			plan := BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{
				RequestPath: clarifyRequestFixturePath, ExpectedStatus: "pending-answers", Mode: "clarify",
				Answers: []QuestionAnswer{{ExpectedLine: "- [ ] Keep the flag?", Outcome: "answered", Summary: labelPrefix + "not really"}},
			}}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
			if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-SUMMARY-INVALID" {
				t.Fatalf("forgeable answered summary accepted: refusal=%#v", plan.Refusal)
			}
			if !strings.Contains(plan.Refusal.Reason, strconv.Quote(labelPrefix)) {
				t.Fatalf("refusal reason does not name the label it rejected: %q", plan.Refusal.Reason)
			}
			if len(plan.Mutations) != 0 {
				t.Fatalf("refused plan still carries mutations: %#v", plan.Mutations)
			}
		})
	}

	// The guard is keyed on the one combination that can forge a verdict, not on the words: with
	// the outcome carrying the disposition, the writer's own label already occupies the position
	// and the summary behind it cannot change what any reader attributes to the line.
	t.Run("the same label inside a discarded summary stays legal and still cancels", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, false, "- [ ] Keep the flag?\n")
		round := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Keep the flag?", Outcome: "discarded", Summary: confirmedLabelPrefix + "keep the flag"})
		if !bytes.Contains(round.Mutations[0].Contents, []byte(resolvedQuestionSeparator+discardedLabelPrefix+confirmedLabelPrefix+"keep the flag")) {
			t.Fatalf("discarded summary not written verbatim behind the writer's label: %s", round.Mutations[0].Contents)
		}
		requireClarifyStatus(t, round, "cancelled")
	})
}

// TestBuildAnswerPlanRefusesADeclaredArchivePathAgainstANonTerminalVerdict makes the blocked
// terminal verdict visible. When a prior round's discard summary carries the project's own
// "A → B" style, the resolved line holds two separators and no reader can attribute its
// disposition, so a round that discards the last open question derives pending rather than
// cancelled. That direction is deliberate, but silence about it is not: the caller declared the
// archive path it had already computed, and publishing pending while ignoring that declaration
// sends a REQ whose every question the user discarded back to the queue to be built.
func TestBuildAnswerPlanRefusesADeclaredArchivePathAgainstANonTerminalVerdict(t *testing.T) {
	const priorDiscardSummary = "superseded by A → B"
	priorLine := "- [x] Keep the flag?" + resolvedQuestionSeparator + discardedLabelPrefix + priorDiscardSummary

	buildDiscardingRound := func(t *testing.T, root string) PublicationPlan {
		t.Helper()
		return BuildAnswerPlan(root, Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{
			RequestPath: clarifyRequestFixturePath, ExpectedStatus: "pending-answers", Mode: "clarify",
			ArchivePath: "do-work/archive/REQ-1-test.md",
			Answers:     []QuestionAnswer{{ExpectedLine: "- [ ] Drop the table?", Outcome: "discarded", Summary: "not needed"}},
		}}, time.Date(2026, 9, 1, 1, 2, 3, 0, time.UTC))
	}

	t.Run("the refusal names the line that blocked the terminal verdict", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, false, "- [ ] Keep the flag?\n- [ ] Drop the table?\n")
		// Compose the prior line through the writer so the blocking line is exactly what
		// production writes for a genuine discard whose summary carries the separator.
		firstRound := buildClarifyRound(t, root, QuestionAnswer{ExpectedLine: "- [ ] Keep the flag?", Outcome: "discarded", Summary: priorDiscardSummary})
		writeFixture(t, root, clarifyRequestFixturePath, firstRound.Mutations[0].Contents, 0o644)
		if !bytes.Contains(firstRound.Mutations[0].Contents, []byte(priorLine)) {
			t.Fatalf("prior round did not compose the expected blocking line: %s", firstRound.Mutations[0].Contents)
		}
		plan := buildDiscardingRound(t, root)
		if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-ARCHIVE-PATH-MISMATCH" {
			t.Fatalf("declared terminal archive path silently ignored: refusal=%#v", plan.Refusal)
		}
		if len(plan.Mutations) != 0 {
			t.Fatalf("refused plan still carries mutations: %#v", plan.Mutations)
		}
		for _, want := range []string{"pending", strconv.Quote(priorLine)} {
			if !strings.Contains(plan.Refusal.Reason, want) {
				t.Fatalf("refusal reason %q does not carry %s", plan.Refusal.Reason, want)
			}
		}
	})

	// The evidence quotes the line, so a CRLF document is where a dropped carriage-return trim
	// becomes observable: the same blocking line must read identically in both formats.
	t.Run("CRLF evidence quotes the blocking line without its carriage return", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, false, priorLine+"\n- [ ] Drop the table?\n")
		rewriteFixtureWithCarriageReturns(t, root)
		plan := buildDiscardingRound(t, root)
		if plan.Refusal == nil || plan.Refusal.Code != "ANSWER-ARCHIVE-PATH-MISMATCH" {
			t.Fatalf("declared terminal archive path silently ignored: refusal=%#v", plan.Refusal)
		}
		if !strings.Contains(plan.Refusal.Reason, strconv.Quote(priorLine)) {
			t.Fatalf("refusal reason %q does not quote the blocking line as written", plan.Refusal.Reason)
		}
	})

	// The declaration is only refused where it disagrees: the same terminal path stays accepted
	// when every resolved line does carry an attributable discard.
	t.Run("a declared archive path matching a genuine cancellation is accepted", func(t *testing.T) {
		root := t.TempDir()
		writeClarifyFixture(t, root, false, "- [x] Keep the flag?"+resolvedQuestionSeparator+discardedLabelPrefix+"not needed\n- [ ] Drop the table?\n")
		plan := buildDiscardingRound(t, root)
		if plan.Refusal != nil {
			t.Fatalf("genuine cancellation refused: %#v", plan.Refusal)
		}
		requireClarifyStatus(t, plan, "cancelled")
	})
}
