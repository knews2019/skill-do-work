package publication

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
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
	fixture := `{"operation":"answer","answer":{"request_path":"do-work/queue/REQ-1.md","expected_status":"blocked","mode":"stakeholder","answers":[{"question_id":"Q-01","outcome":"answered","summary":"override"}],"override_capture":{"user_request_id":"UR-2","user_request":{"path":"do-work/user-requests/UR-2/input.md","payload":{"source_path":"payload/ur"}},"raw_input":{"source_path":"payload/raw"},"requests":[{"id":"REQ-2","user_request_id":"UR-2","file":{"path":"do-work/queue/REQ-2-change.md","payload":{"source_path":"payload/req"}},"reservation_path":"do-work/.req-reservations/REQ-2"}]}}}`
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
	writeFixture(t, root, "payload/ur", []byte("---\nid: UR-2\nrequests: [REQ-2]\n---\n"+string(containedOutsideBytes(raw, "\n"))+"\n"), 0o644)
	writeFixture(t, root, "payload/req", []byte("---\nid: REQ-2\nstatus: pending\nuser_request: UR-2\n---\n"), 0o644)
	writeFixture(t, root, "payload/blocked-history", []byte("## Blocked\n\n- Resolved after stakeholder override.\n"), 0o644)
	writeFixture(t, root, "payload/implementation", []byte("## Implementation\n\nNo code changes.\n"), 0o644)
	override := &CaptureManifest{UserRequestID: "UR-2", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-2/input.md", Payload: PayloadFile{SourcePath: "payload/ur"}}, RawInput: &PayloadFile{SourcePath: "payload/raw"}, Requests: []CaptureRequest{{ID: "REQ-2", UserRequestID: "UR-2", File: PublishedFile{Path: "do-work/queue/REQ-2-change.md", Payload: PayloadFile{SourcePath: "payload/req"}}, ReservationPath: "do-work/.req-reservations/REQ-2"}}}
	answer := &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "blocked", Mode: "stakeholder", ArchivePath: "do-work/archive/REQ-1-test.md", Answers: []QuestionAnswer{{QuestionID: "Q-01", Outcome: "answered", Summary: "override"}}, OverrideCapture: override, StakeholderTerminal: &StakeholderTerminalEvidence{BlockedHistory: PayloadFile{SourcePath: "payload/blocked-history"}, Implementation: PayloadFile{SourcePath: "payload/implementation"}}}
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
	for _, path := range []string{"do-work/archive/REQ-1-test.md", "do-work/queue/REQ-2-change.md", "do-work/user-requests/UR-2/input.md", "do-work/.req-reservations/REQ-2"} {
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
