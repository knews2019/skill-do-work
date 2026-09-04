package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestAdvanceQueueModeClaimsDefaultAndExplicitTargets(t *testing.T) {
	for _, test := range []struct {
		name      string
		arguments []string
		wantID    string
	}{
		{name: "default", wantID: "REQ-801"},
		{name: "explicit", arguments: []string{"REQ-802"}, wantID: "REQ-802"},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := newAdvanceQueueRepository(t)
			writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-801", "pending", "priority: now\n", "")
			writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-802", "pending", "priority: later\n", "")
			commitAdvanceQueueFixture(t, repositoryRoot)

			result, status := runAdvanceQueueJSON(t, repositoryRoot, test.arguments...)
			if status != 0 || result.Outcome != "success" || result.QueueAdvance == nil {
				t.Fatalf("status=%d result=%#v", status, result)
			}
			if len(result.QueueAdvance.Claimed) != 1 || result.QueueAdvance.Claimed[0].RequestID != test.wantID {
				t.Fatalf("claimed=%#v, want %s", result.QueueAdvance.Claimed, test.wantID)
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", test.wantID+"-fixture.md")); err != nil {
				t.Fatalf("working request missing: %v", err)
			}
			if dirty := string(runAdvanceGit(t, repositoryRoot, "status", "--porcelain=v1")); dirty != "" {
				t.Fatalf("claim left dirty tree: %q", dirty)
			}
			assertAdvanceCommitPaths(t, repositoryRoot, []string{
				"do-work/CHECKPOINT.md",
				"do-work/queue/" + test.wantID + "-fixture.md",
				"do-work/working/" + test.wantID + "-fixture.md",
			})
			if test.name == "default" && (len(result.QueueAdvance.FrozenMembers) != 1 || len(result.QueueAdvance.ContinuationArgv) != 0) {
				t.Fatalf("fresh default must retain canonical one-record bound: %#v", result.QueueAdvance)
			}
		})
	}
}

type advanceQueueCommandResult struct {
	Outcome      string                `json:"outcome"`
	QueueAdvance *advanceQueueEvidence `json:"queue_advance"`
}

type advanceQueueEvidence struct {
	TargetTokens     []string             `json:"target_tokens"`
	DispatchBound    int                  `json:"dispatch_bound"`
	FrozenMembers    []advanceQueueMember `json:"frozen_members"`
	Claimed          []advanceQueueMember `json:"claimed"`
	Phases           []advanceQueuePhase  `json:"phases"`
	Partial          bool                 `json:"partial"`
	ContinuationArgv []string             `json:"continuation_argv"`
}

type advanceQueueMember struct {
	RequestID   string `json:"request_id"`
	RequestPath string `json:"request_path"`
	Provenance  string `json:"provenance"`
	Consumed    bool   `json:"consumed"`
}

type advanceQueuePhase struct {
	RequestID string                  `json:"request_id"`
	Phase     string                  `json:"phase"`
	Outcome   string                  `json:"outcome"`
	Findings  []advanceCommandFinding `json:"findings"`
}

func TestAdvanceQueueTargetLedgerProjectsBeforeBounding(t *testing.T) {
	repositoryRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-821", "pending", "user_request: UR-820\npriority: now\n", "")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-829", "pending", "user_request: UR-820\npriority: later\n", "")
	commitAdvanceQueueFixture(t, repositoryRoot)

	first, status := runAdvanceQueueJSON(t, repositoryRoot, "UR-820", "--fan-out", "1")
	if status != 0 || first.QueueAdvance == nil || len(first.QueueAdvance.Claimed) != 1 || first.QueueAdvance.Claimed[0].RequestID != "REQ-821" {
		t.Fatalf("first queue advance = status %d %#v", status, first)
	}
	if first.QueueAdvance.DispatchBound != 1 || !reflect.DeepEqual(queueMemberIDs(first.QueueAdvance.FrozenMembers), []string{"REQ-821", "REQ-829"}) {
		t.Fatalf("frozen ledger = %#v", first.QueueAdvance)
	}
	if containsArgument(first.QueueAdvance.ContinuationArgv, "--fan-out") {
		t.Fatalf("continuation retained observation bound: %#v", first.QueueAdvance.ContinuationArgv)
	}
	completeClaimedFixture(t, repositoryRoot, "REQ-821")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-801", "pending", "user_request: UR-820\npriority: now\n", "")
	runAdvanceGit(t, repositoryRoot, "add", "do-work")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "integrate first and add later UR member")

	second, secondStatus := runAdvanceContinuationJSON(t, repositoryRoot, first.QueueAdvance.ContinuationArgv)
	if secondStatus != 0 || second.QueueAdvance == nil || len(second.QueueAdvance.Claimed) != 1 || second.QueueAdvance.Claimed[0].RequestID != "REQ-829" {
		t.Fatalf("continuation did not project frozen members before bound: status=%d %#v", secondStatus, second)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "queue", "REQ-801-fixture.md")); err != nil {
		t.Fatalf("later out-of-ledger UR member was consumed: %v", err)
	}
}

func TestAdvanceQueueClaimsURChainAndForkAcrossStatelessContinuation(t *testing.T) {
	repositoryRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-831", "pending", "user_request: UR-830\n", "")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-832", "pending", "user_request: UR-830\ndepends_on: [REQ-831]\n", "")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-833", "pending", "user_request: UR-830\ndepends_on: [REQ-831]\n", "")
	commitAdvanceQueueFixture(t, repositoryRoot)

	first, status := runAdvanceQueueJSON(t, repositoryRoot, "UR-830")
	if status != 0 || first.QueueAdvance == nil || !reflect.DeepEqual(queueMemberIDs(first.QueueAdvance.FrozenMembers), []string{"REQ-831", "REQ-832", "REQ-833"}) || first.QueueAdvance.FrozenMembers[1].Provenance != "ur-expanded" {
		t.Fatalf("UR chain was not frozen: status=%d %#v", status, first)
	}
	completeClaimedFixture(t, repositoryRoot, "REQ-831")
	runAdvanceGit(t, repositoryRoot, "add", "do-work")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "integrate prerequisite")
	second, secondStatus := runAdvanceContinuationJSON(t, repositoryRoot, first.QueueAdvance.ContinuationArgv)
	if secondStatus != 0 || len(second.QueueAdvance.Claimed) != 1 || second.QueueAdvance.Claimed[0].RequestID != "REQ-832" || len(second.QueueAdvance.ContinuationArgv) == 0 {
		t.Fatalf("UR continuation = status %d %#v", secondStatus, second)
	}
	completeClaimedFixture(t, repositoryRoot, "REQ-832")
	runAdvanceGit(t, repositoryRoot, "add", "do-work")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "integrate first fork")
	third, thirdStatus := runAdvanceContinuationJSON(t, repositoryRoot, second.QueueAdvance.ContinuationArgv)
	if thirdStatus != 0 || len(third.QueueAdvance.Claimed) != 1 || third.QueueAdvance.Claimed[0].RequestID != "REQ-833" || len(third.QueueAdvance.ContinuationArgv) != 0 {
		t.Fatalf("UR fork continuation = status %d %#v", thirdStatus, third)
	}
}

func TestAdvanceQueueRunsBlockedProbeThenUnblocksAndClaims(t *testing.T) {
	repositoryRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-841", "blocked", "blocked_by: waiting\nblocked_check: exit 0\n", "")
	commitAdvanceQueueFixture(t, repositoryRoot)

	result, status := runAdvanceQueueJSON(t, repositoryRoot, "REQ-841")
	if status != 0 || result.QueueAdvance == nil || !reflect.DeepEqual(queuePhaseNames(result.QueueAdvance.Phases), []string{"unblock", "claim"}) {
		t.Fatalf("successful probe advance = status %d %#v", status, result)
	}
	if subjects := strings.Fields(string(runAdvanceGit(t, repositoryRoot, "log", "-2", "--format=%s"))); len(subjects) == 0 {
		t.Fatal("unblock and claim commits were not created")
	}

	failingRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, failingRoot, "queue", "REQ-842", "blocked", "blocked_by: waiting\nblocked_check: exit 9\n", "")
	commitAdvanceQueueFixture(t, failingRoot)
	failure, failureStatus := runAdvanceQueueJSON(t, failingRoot, "REQ-842")
	if failureStatus != 1 || failure.QueueAdvance == nil || len(failure.QueueAdvance.Claimed) != 0 || len(failure.QueueAdvance.Phases) != 1 || failure.QueueAdvance.Phases[0].Phase != "selection" {
		t.Fatalf("failed probe advance = status %d %#v", failureStatus, failure)
	}
	if contents := readAdvanceQueueFile(t, failingRoot, "do-work/queue/REQ-842-fixture.md"); !strings.Contains(contents, "status: blocked") {
		t.Fatalf("failed probe changed request:\n%s", contents)
	}
}

func TestAdvanceQueueCommitsNestedArchiveCollisionAndCycleHolds(t *testing.T) {
	collisionRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, collisionRoot, "queue", "REQ-851", "pending", "", "")
	writeAdvanceRequest(t, collisionRoot, "archive/UR-850", "REQ-851", "completed", "commit: abc1234\n", "")
	commitAdvanceQueueFixture(t, collisionRoot)
	collision, collisionStatus := runAdvanceQueueJSON(t, collisionRoot, "REQ-851")
	if collisionStatus != 0 || collision.QueueAdvance == nil || !reflect.DeepEqual(queuePhaseNames(collision.QueueAdvance.Phases), []string{"archive-collision-hold"}) {
		t.Fatalf("collision hold = status %d %#v", collisionStatus, collision)
	}
	if contents := readAdvanceQueueFile(t, collisionRoot, "do-work/queue/REQ-851-fixture.md"); !strings.Contains(contents, "status: blocked-archive-collision") {
		t.Fatalf("collision status not held:\n%s", contents)
	}

	cycleRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, cycleRoot, "queue", "REQ-861", "pending", "user_request: UR-860\ndepends_on: [REQ-862]\n", "")
	writeAdvanceRequest(t, cycleRoot, "queue", "REQ-862", "pending", "user_request: UR-860\ndepends_on: [REQ-861]\n", "")
	commitAdvanceQueueFixture(t, cycleRoot)
	cycle, cycleStatus := runAdvanceQueueJSON(t, cycleRoot)
	if cycleStatus != 0 || cycle.QueueAdvance == nil || !reflect.DeepEqual(queuePhaseNames(cycle.QueueAdvance.Phases), []string{"dependency-cycle-hold", "dependency-cycle-hold"}) {
		t.Fatalf("cycle holds = status %d %#v", cycleStatus, cycle)
	}
	for _, requestID := range []string{"REQ-861", "REQ-862"} {
		if contents := readAdvanceQueueFile(t, cycleRoot, "do-work/queue/"+requestID+"-fixture.md"); !strings.Contains(contents, "status: blocked-dependency-cycle") {
			t.Fatalf("%s cycle status not held:\n%s", requestID, contents)
		}
	}
}

func TestAdvanceQueueReportsCommittedPartialClaimBeforeDirtyRefusal(t *testing.T) {
	repositoryRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-871", "pending", "user_request: UR-870\n", "")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-872", "pending", "user_request: UR-870\n", "")
	commitAdvanceQueueFixture(t, repositoryRoot)
	secondPath := filepath.Join(repositoryRoot, "do-work", "queue", "REQ-872-fixture.md")
	if err := os.WriteFile(secondPath, []byte(readAdvanceQueueFile(t, repositoryRoot, "do-work/queue/REQ-872-fixture.md")+"\nlocal dirt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, status := runAdvanceQueueJSON(t, repositoryRoot, "UR-870", "--fan-out", "2")
	if status != 1 || result.Outcome != "findings" || result.QueueAdvance == nil || !result.QueueAdvance.Partial || len(result.QueueAdvance.Claimed) != 1 || result.QueueAdvance.Claimed[0].RequestID != "REQ-871" {
		t.Fatalf("partial refusal = status %d %#v", status, result)
	}
	if !reflect.DeepEqual(queuePhaseNames(result.QueueAdvance.Phases), []string{"claim", "claim"}) || result.QueueAdvance.Phases[1].Outcome == "success" {
		t.Fatalf("partial phase truth lost: %#v", result.QueueAdvance.Phases)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "working", "REQ-871-fixture.md")); err != nil {
		t.Fatalf("prior committed claim rolled back: %v", err)
	}
}

func TestAdvanceQueueWaveFanOutAndHostileTokens(t *testing.T) {
	repositoryRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-881", "pending", "", "")
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-882", "pending", "", "")
	commitAdvanceQueueFixture(t, repositoryRoot)
	result, status := runAdvanceQueueJSON(t, repositoryRoot, "--wave", "0", "--fan-out", "2")
	if status != 0 || result.QueueAdvance == nil || len(result.QueueAdvance.Claimed) != 2 {
		t.Fatalf("wave fan-out = status %d %#v", status, result)
	}

	hostileRoot := newAdvanceQueueRepository(t)
	writeAdvanceRequest(t, hostileRoot, "queue", "REQ-891", "pending", "", "")
	commitAdvanceQueueFixture(t, hostileRoot)
	hostile, hostileStatus := runAdvanceQueueJSON(t, hostileRoot, "REQ-891;touch-owned")
	if hostileStatus != 2 || hostile.Outcome != "failure" {
		t.Fatalf("hostile token accepted: status %d %#v", hostileStatus, hostile)
	}
	if _, err := os.Stat(filepath.Join(hostileRoot, "touch-owned")); !os.IsNotExist(err) {
		t.Fatalf("hostile token reached shell: %v", err)
	}
}

func newAdvanceQueueRepository(t *testing.T) string {
	t.Helper()
	repositoryRoot := t.TempDir()
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Queue Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance-queue@example.invalid")
	return repositoryRoot
}

func commitAdvanceQueueFixture(t *testing.T, repositoryRoot string) {
	t.Helper()
	runAdvanceGit(t, repositoryRoot, "add", "do-work")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
}

func assertAdvanceCommitPaths(t *testing.T, repositoryRoot string, want []string) {
	t.Helper()
	got := strings.Fields(string(runAdvanceGit(t, repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD")))
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("committed paths = %#v, want %#v", got, want)
	}
}

func runAdvanceQueueJSON(t *testing.T, repositoryRoot string, arguments ...string) (advanceQueueCommandResult, int) {
	t.Helper()
	commandArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "advance"}
	commandArguments = append(commandArguments, arguments...)
	command := exec.Command(advanceCLIBinary(t), commandArguments...)
	output, err := command.CombinedOutput()
	status := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
		} else {
			t.Fatalf("advance launch: %v", err)
		}
	}
	var result advanceQueueCommandResult
	if decodeError := json.Unmarshal(output, &result); decodeError != nil {
		t.Fatalf("decode advance result: %v\n%s", decodeError, output)
	}
	return result, status
}

func runAdvanceContinuationJSON(t *testing.T, repositoryRoot string, continuationArgv []string) (advanceQueueCommandResult, int) {
	t.Helper()
	if len(continuationArgv) < 4 || continuationArgv[0] != "do-work-cli" || continuationArgv[3] != "advance" {
		t.Fatalf("invalid continuation argv: %#v", continuationArgv)
	}
	return runAdvanceQueueJSON(t, repositoryRoot, continuationArgv[4:]...)
}

func completeClaimedFixture(t *testing.T, repositoryRoot, requestID string) {
	t.Helper()
	workingPath := filepath.Join(repositoryRoot, "do-work", "working", requestID+"-fixture.md")
	contents, readError := os.ReadFile(workingPath)
	if readError != nil {
		t.Fatal(readError)
	}
	updated := strings.Replace(string(contents), "status: claimed", "status: completed", 1)
	archivePath := filepath.Join(repositoryRoot, "do-work", "archive", requestID+"-fixture.md")
	if makeError := os.MkdirAll(filepath.Dir(archivePath), 0o755); makeError != nil {
		t.Fatal(makeError)
	}
	if writeError := os.WriteFile(workingPath, []byte(updated), 0o644); writeError != nil {
		t.Fatal(writeError)
	}
	if renameError := os.Rename(workingPath, archivePath); renameError != nil {
		t.Fatal(renameError)
	}
}

func queueMemberIDs(members []advanceQueueMember) []string {
	requestIDs := make([]string, len(members))
	for memberIndex, member := range members {
		requestIDs[memberIndex] = member.RequestID
	}
	return requestIDs
}

func queuePhaseNames(phases []advanceQueuePhase) []string {
	names := make([]string, len(phases))
	for phaseIndex, phase := range phases {
		names[phaseIndex] = phase.Phase
	}
	return names
}

func containsArgument(arguments []string, wanted string) bool {
	for _, argument := range arguments {
		if argument == wanted {
			return true
		}
	}
	return false
}

func readAdvanceQueueFile(t *testing.T, repositoryRoot, relativePath string) string {
	t.Helper()
	contents, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if readError != nil {
		t.Fatalf("read %s: %v", relativePath, readError)
	}
	return string(contents)
}
