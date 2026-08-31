package publication

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
	manifest := Manifest{Operation: OperationAnswer, Answer: &AnswerManifest{RequestPath: "do-work/queue/REQ-1-test.md", ExpectedStatus: "blocked", Mode: "stakeholder", Answers: []QuestionAnswer{{QuestionID: "Q-01", Outcome: "answered", Summary: "blue"}}, ArchivePath: "do-work/archive/UR-1/REQ-1-test.md", CloseUserRequest: true, UserRequestPath: "do-work/user-requests/UR-1", ArchiveDirectory: "do-work/archive/UR-1"}}
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
