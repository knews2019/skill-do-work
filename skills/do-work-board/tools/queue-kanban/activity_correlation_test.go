package main

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// gitRunnerThatFindsNothing stands in for git in fixtures that only care about
// stamps: every command fails, the way it does outside a repository, so the
// activity collector falls back to the lifecycle stamps alone and no git runs.
func gitRunnerThatFindsNothing(string, ...string) ([]byte, error) {
	return nil, errors.New("no git in this fixture")
}

// cannedLogRecord renders one commit exactly as the collector's
// `git log --format=%H%x00%cI%x00%P%x00%s --name-only` prints it: a NUL-joined
// header line, a blank line, then one touched path per line. A merge prints no
// paths under plain --name-only, which an empty touchedPaths reproduces.
func cannedLogRecord(commitHash string, committedAt time.Time, parentHashes string, subject string, touchedPaths ...string) string {
	record := commitHash + "\x00" + committedAt.Format(time.RFC3339) + "\x00" + parentHashes + "\x00" + subject + "\n"
	if len(touchedPaths) > 0 {
		record += "\n" + strings.Join(touchedPaths, "\n") + "\n"
	}
	return record
}

// cannedGitRunner answers the git reads behind a board response — the windowed
// log, the worktree list, the integration branch, and the worktree-agent branch
// listings with their tip dates — from fixed text, so attribution tests never
// spawn git. A branch named in mergedBranches has its tip reachable from
// the integration branch, which is what `--no-merged` filters on.
type cannedGitRunner struct {
	mutex                 sync.Mutex
	logOutput             func() string
	branchTips            map[string]time.Time
	mergedBranches        map[string]bool
	logCallCount          int
	callCount             int
	worktreeListCallCount int
}

func (runner *cannedGitRunner) run(_ string, arguments ...string) ([]byte, error) {
	runner.mutex.Lock()
	defer runner.mutex.Unlock()
	runner.callCount++
	if len(arguments) == 0 {
		return nil, errors.New("no git subcommand")
	}
	branchNames := make([]string, 0, len(runner.branchTips))
	for branchName := range runner.branchTips {
		branchNames = append(branchNames, branchName)
	}
	sort.Strings(branchNames)
	switch arguments[0] {
	case "worktree":
		runner.worktreeListCallCount++
		return []byte("worktree /repo\nHEAD c0\nbranch refs/heads/main\n\n"), nil
	case "rev-parse":
		return []byte("main\n"), nil
	case "branch":
		return []byte(strings.Join(branchNames, "\n") + "\n"), nil
	case "for-each-ref":
		var listing strings.Builder
		for _, branchName := range branchNames {
			if runner.mergedBranches[branchName] && slices.Contains(arguments, "--no-merged=main") {
				continue
			}
			listing.WriteString(branchName + " " + runner.branchTips[branchName].Format(time.RFC3339) + "\n")
		}
		return []byte(listing.String()), nil
	case "log":
		runner.logCallCount++
		if runner.logOutput == nil {
			return nil, nil
		}
		return []byte(runner.logOutput()), nil
	}
	return nil, errors.New("unexpected git subcommand " + arguments[0])
}

var activityFixtureNow = time.Date(2026, 10, 5, 18, 0, 0, 0, time.UTC)

// The brief's case: a commit whose subject carries no [REQ-NNN] prefix but which
// touches the REQ file must still count, or a claimed card reads idle while the
// orchestrator is actively writing its trail.
func TestRequestActivityAttributesACommitThatTouchesTheRequestFileWithoutAPrefix(t *testing.T) {
	claimedTicket := &RequestTicket{
		RequestId:  "REQ-701",
		Status:     "claimed",
		ClaimedAt:  activityFixtureNow.Add(-3 * time.Hour).Format(time.RFC3339),
		DispatchAt: activityFixtureNow.Add(-2 * time.Hour).Format(time.RFC3339),
	}
	commitInstant := activityFixtureNow.Add(-10 * time.Minute)
	runner := &cannedGitRunner{logOutput: func() string {
		return cannedLogRecord("c1", commitInstant, "c0", "docs(do-work): tidy the queue notes",
			"do-work/working/REQ-701-board-cards.md")
	}}

	activityById := collectRequestActivity("/repo", []*RequestTicket{claimedTicket},
		activityFixtureNow.Add(-4*time.Hour), activityFixtureNow, runner.run, nil)
	activity, present := activityById["REQ-701"]
	if !present {
		t.Fatalf("REQ-701 carries no activity; a path-only commit must attribute")
	}
	if !activity.LastActivityAt.Equal(commitInstant) || activity.LastActivityKind != "commit" {
		t.Errorf("last activity = %v (%s), want the path-only commit at %v (commit)",
			activity.LastActivityAt, activity.LastActivityKind, commitInstant)
	}
	if activity.LastActivityPhase != "dispatch" {
		t.Errorf("last activity phase = %q, want %q — a commit sits in the interval its newest preceding stamp opened",
			activity.LastActivityPhase, "dispatch")
	}
}

// Every [REQ-NNN] token in a subject counts, not only the first: a shared fix
// committed as "[REQ-701] [REQ-702] …" is evidence for both cards. A bare id
// without brackets is prose, not a prefix.
func TestRequestActivityAttributesEveryPrefixTokenInASubject(t *testing.T) {
	commitInstant := activityFixtureNow.Add(-20 * time.Minute)
	instantsById := correlateCommitsToRequests([]byte(
		cannedLogRecord("c1", commitInstant, "c0", "[REQ-701] [REQ-702] fold the shared parser fix; see REQ-703") +
			cannedLogRecord("c0", commitInstant.Add(-time.Hour), "", "seed", "README.md")))

	for _, requestId := range []string{"REQ-701", "REQ-702"} {
		if instants := instantsById[requestId]; len(instants) != 1 || !instants[0].Equal(commitInstant) {
			t.Errorf("%s attributed instants = %v, want exactly the prefixed commit", requestId, instants)
		}
	}
	if instants := instantsById["REQ-703"]; len(instants) != 0 {
		t.Errorf("REQ-703 attributed %v from an unbracketed mention", instants)
	}
}

// Builder commits on a worktree branch carry no REQ path and need not carry a
// prefix. They reach the board only as the second-parent side of the merge that
// brought them in; without the ancestry rule a builder working for an hour shows
// nothing until its merge lands.
func TestRequestActivityAttributesBuilderCommitsReachableOnlyThroughAMatchedMergesSecondParent(t *testing.T) {
	baseInstant := activityFixtureNow.Add(-5 * time.Hour)
	logOutput := cannedLogRecord("merge1", baseInstant.Add(4*time.Hour), "main1 build2",
		"[REQ-701] merge builder branch worktree-agent-REQ-701-board-cards") +
		cannedLogRecord("build2", baseInstant.Add(3*time.Hour), "build1", "write the parser", "src/parser.go") +
		cannedLogRecord("main1", baseInstant.Add(150*time.Minute), "base0", "unrelated main commit", "README.md") +
		cannedLogRecord("build1", baseInstant.Add(2*time.Hour), "base0", "red tests first", "src/parser_test.go") +
		// An UNMATCHED merge: its second parent must not attribute to anyone.
		cannedLogRecord("merge9", baseInstant.Add(90*time.Minute), "base0 other1", "merge a branch nobody claimed") +
		cannedLogRecord("other1", baseInstant.Add(80*time.Minute), "base0", "someone else's work", "src/other.go") +
		cannedLogRecord("base0", baseInstant, "older9", "base", "go.mod")

	instantsById := correlateCommitsToRequests([]byte(logOutput))
	got := map[time.Time]bool{}
	for _, instant := range instantsById["REQ-701"] {
		got[instant] = true
	}
	for name, wantInstant := range map[string]time.Time{
		"the matched merge":        baseInstant.Add(4 * time.Hour),
		"the builder tip":          baseInstant.Add(3 * time.Hour),
		"the first builder commit": baseInstant.Add(2 * time.Hour),
	} {
		if !got[wantInstant] {
			t.Errorf("REQ-701 lacks %s at %v; attributed %v", name, wantInstant, instantsById["REQ-701"])
		}
	}
	for name, unwantedInstant := range map[string]time.Time{
		"the first parent's own commit": baseInstant.Add(150 * time.Minute),
		"the shared merge base":         baseInstant,
		"an unmatched merge's work":     baseInstant.Add(80 * time.Minute),
	} {
		if got[unwantedInstant] {
			t.Errorf("REQ-701 wrongly attributed %s at %v", name, unwantedInstant)
		}
	}
	if len(instantsById["REQ-701"]) != 3 {
		t.Errorf("REQ-701 attributed %d instants, want 3: %v", len(instantsById["REQ-701"]), instantsById["REQ-701"])
	}
}

// A builder that has committed but not yet handed back is visible only on its
// live worktree-agent branch, which the windowed main-line log cannot reach.
func TestRequestActivityCountsALiveWorktreeAgentBranchTip(t *testing.T) {
	claimedTicket := &RequestTicket{
		RequestId:  "REQ-701",
		Status:     "claimed",
		ClaimedAt:  activityFixtureNow.Add(-3 * time.Hour).Format(time.RFC3339),
		DispatchAt: activityFixtureNow.Add(-2 * time.Hour).Format(time.RFC3339),
	}
	tipInstant := activityFixtureNow.Add(-5 * time.Minute)
	runner := &cannedGitRunner{
		branchTips: map[string]time.Time{
			"worktree-agent-REQ-701-board-cards":  tipInstant,
			"worktree-agent-REQ-999-someone-else": activityFixtureNow.Add(-time.Minute),
		},
	}

	activity := collectRequestActivity("/repo", []*RequestTicket{claimedTicket},
		activityFixtureNow.Add(-4*time.Hour), activityFixtureNow, runner.run,
		readWorktreeAgentGitState("/repo", runner.run).ownedTipInstantsById)["REQ-701"]
	if !activity.LastActivityAt.Equal(tipInstant) || activity.LastActivityKind != "commit" {
		t.Errorf("last activity = %v (%s), want the live branch tip at %v (commit)",
			activity.LastActivityAt, activity.LastActivityKind, tipInstant)
	}
}

// Commits and stamps are ONE event stream: a commit inside the longest stamp gap
// splits it, and a sibling's commit that touches this REQ's archived file after
// completion (REQ-631's completion commit touched REQ-630's file) is not idle
// time inside this REQ's work.
func TestRequestActivityLargestGapUnitesCommitsWithStamps(t *testing.T) {
	claimInstant := activityFixtureNow.Add(-30 * time.Hour)
	doneTicket := &RequestTicket{
		RequestId:         "REQ-702",
		Status:            "completed",
		CreatedAt:         claimInstant.Add(-20 * time.Hour).Format(time.RFC3339),
		ClaimedAt:         claimInstant.Format(time.RFC3339),
		DispatchAt:        claimInstant.Add(10 * time.Minute).Format(time.RFC3339),
		BuilderHandbackAt: claimInstant.Add(130 * time.Minute).Format(time.RFC3339),
		IntegrationAt:     claimInstant.Add(180 * time.Minute).Format(time.RFC3339),
		CompletedAt:       claimInstant.Add(200 * time.Minute).Format(time.RFC3339),
	}
	runner := &cannedGitRunner{logOutput: func() string {
		return cannedLogRecord("late1", claimInstant.Add(26*time.Hour), "c9", "[REQ-799] complete: a sibling",
			"do-work/archive/UR-099/REQ-702-done-earlier.md", "do-work/archive/UR-099/REQ-799-sibling.md") +
			cannedLogRecord("mid1", claimInstant.Add(70*time.Minute), "c0", "[REQ-702] builder progress", "src/a.go")
	}}

	activity := collectRequestActivity("/repo", []*RequestTicket{doneTicket},
		claimInstant.Add(-time.Hour), activityFixtureNow, runner.run, nil)["REQ-702"]
	if activity.LargestGap == nil {
		t.Fatalf("REQ-702 has no largest gap")
	}
	// Stamps alone, dispatch→handback is the 120-minute gap. The commit at +70
	// splits it into two 60-minute halves; the first half wins the tie and sits
	// wholly inside the dispatch interval.
	wantGap := activityGap{Minutes: 60, FromPhase: "dispatch", ToPhase: "dispatch"}
	if *activity.LargestGap != wantGap {
		t.Errorf("largest gap = %+v, want %+v", *activity.LargestGap, wantGap)
	}
	if !activity.LastActivityAt.IsZero() {
		t.Errorf("a done REQ carries last activity %v; only open requests state it", activity.LastActivityAt)
	}
}

// The REQ-632 Red-Green case (b) through the production payload: a 4h21m span
// whose largest gap is 67 minutes between dispatch and builder handback.
func TestGeneratedPayloadCarriesTheLargestActivityGap(t *testing.T) {
	boardData := buildImplementationSpanFixturePayload(t)
	activity, present := boardData.RequestActivity["REQ-909"]
	if !present || activity.LargestActivityGap == nil {
		t.Fatalf("REQ-909 payload carries no largestActivityGap: %+v", boardData.RequestActivity)
	}
	wantGap := generatedActivityGap{Minutes: 67, FromPhase: "dispatch", ToPhase: "builder handback"}
	if *activity.LargestActivityGap != wantGap {
		t.Errorf("REQ-909 largestActivityGap = %+v, want %+v", *activity.LargestActivityGap, wantGap)
	}
	if activity.LastActivityAt != "" {
		t.Errorf("done REQ-909 ships lastActivityAt %q; done cards show nothing new", activity.LastActivityAt)
	}
}

// REQ-284's rule applied to the new read: git history changes while every file
// mtime stays identical, so the activity read must run on EVERY /board-data.js
// response, outside refreshBoardData's mtime cache.
func TestServeReadsRequestActivityOnEveryResponse(t *testing.T) {
	claimInstant := time.Now().UTC().Add(-3 * time.Hour).Truncate(time.Second)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/working/REQ-711-live-card.md", "---\nid: REQ-711\ntitle: live card\nstatus: claimed\nclaimed_at: " +
			claimInstant.Format(time.RFC3339) + "\n---\n\n# live card\n"},
	})
	commitInstant := claimInstant.Add(time.Hour)
	runner := &cannedGitRunner{logOutput: func() string {
		return cannedLogRecord("c1", commitInstant, "c0", "touch", "do-work/working/REQ-711-live-card.md")
	}}
	liveServer := newLiveBoardServer(repoRoot, defaultRecentWindow)
	liveServer.liveGitRunner = runner.run
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	first := fetchServedBoardData(t, testServer.URL).RequestActivity["REQ-711"]
	if first.LastActivityAt != commitInstant.Format(time.RFC3339) {
		t.Fatalf("first response lastActivityAt = %q, want %q", first.LastActivityAt, commitInstant.Format(time.RFC3339))
	}
	commitInstant = commitInstant.Add(30 * time.Minute)
	second := fetchServedBoardData(t, testServer.URL).RequestActivity["REQ-711"]
	if second.LastActivityAt != commitInstant.Format(time.RFC3339) {
		t.Errorf("second response lastActivityAt = %q, want the new commit %q — the read was cached by file mtime",
			second.LastActivityAt, commitInstant.Format(time.RFC3339))
	}
	if runner.logCallCount != 2 {
		t.Errorf("git log ran %d times for two responses, want 2", runner.logCallCount)
	}
}

// REQ-636's Red-Green case: a builder branch cut from the integration branch with
// no commits of its own points at the integration commit. That commit is not
// evidence for the REQ, so the card's last activity stays its claim stamp.
func TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn(t *testing.T) {
	claimInstant := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/working/REQ-701-unowned-tip.md", "---\nid: REQ-701\ntitle: unowned tip\nstatus: claimed\nclaimed_at: " +
			claimInstant.Format(time.RFC3339) + "\n---\n\n# unowned tip\n"},
	})
	runner := &cannedGitRunner{
		branchTips:     map[string]time.Time{"worktree-agent-REQ-701-x": claimInstant.Add(50 * time.Minute)},
		mergedBranches: map[string]bool{"worktree-agent-REQ-701-x": true},
	}
	liveServer := newLiveBoardServer(repoRoot, defaultRecentWindow)
	liveServer.liveGitRunner = runner.run
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	activity := fetchServedBoardData(t, testServer.URL).RequestActivity["REQ-701"]
	if activity.LastActivityAt != claimInstant.Format(time.RFC3339) || activity.LastActivityKind != "stamp" {
		t.Errorf("last activity = %q (%s), want the claim stamp %q (stamp) — the branch tip is the integration commit",
			activity.LastActivityAt, activity.LastActivityKind, claimInstant.Format(time.RFC3339))
	}
	if activity.LargestActivityGap != nil {
		t.Errorf("largest gap = %+v, want none — an unowned tip must not split the idle time", *activity.LargestActivityGap)
	}
}

// One served response lists worktrees once and reads the worktree-agent branches
// with a fixed number of git commands, however many branches exist. Verify and
// activity used to list both twice and then run one `git log -1` per branch.
func TestServedResponseListsWorktreesAndBranchesOnce(t *testing.T) {
	claimInstant := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/working/REQ-701-counted.md", "---\nid: REQ-701\ntitle: counted\nstatus: claimed\nclaimed_at: " +
			claimInstant.Format(time.RFC3339) + "\n---\n\n# counted\n"},
	})
	gitCallsForBranchCount := func(branchCount int) (totalCalls int, worktreeListCalls int) {
		runner := &cannedGitRunner{branchTips: map[string]time.Time{}}
		for branchIndex := 0; branchIndex < branchCount; branchIndex++ {
			runner.branchTips[fmt.Sprintf("worktree-agent-REQ-%d-x", 801+branchIndex)] = claimInstant
		}
		liveServer := newLiveBoardServer(repoRoot, defaultRecentWindow)
		liveServer.liveGitRunner = runner.run
		testServer := httptest.NewServer(liveServer)
		defer testServer.Close()
		fetchServedBoardData(t, testServer.URL)
		return runner.callCount, runner.worktreeListCallCount
	}

	oneBranchCalls, oneBranchWorktreeLists := gitCallsForBranchCount(1)
	fourBranchCalls, fourBranchWorktreeLists := gitCallsForBranchCount(4)
	if oneBranchWorktreeLists != 1 || fourBranchWorktreeLists != 1 {
		t.Errorf("worktree listings per response = %d and %d, want 1 each", oneBranchWorktreeLists, fourBranchWorktreeLists)
	}
	if oneBranchCalls != fourBranchCalls {
		t.Errorf("git commands per response = %d with 1 branch and %d with 4, want a constant", oneBranchCalls, fourBranchCalls)
	}
}

// Archive files nest under a UR folder and run artifacts carry suffixes; both
// must attribute, while queue metadata that merely names a REQ must not.
func TestRequestPathPatternMatchesOnlyRequestFilesAndRunArtifacts(t *testing.T) {
	for path, wantId := range map[string]string{
		"do-work/archive/UR-134/REQ-630-takeover-message.md":   "REQ-630",
		"do-work/queue/REQ-633-panel-b.md":                     "REQ-633",
		"do-work/runs/work-2026-10-05-203159/REQ-632-brief.md": "REQ-632",
		"do-work/runs/work-2026-10-05-203159/REQ-632-probe.sh": "REQ-632",
		"do-work/.req-reservations/REQ-632":                    "",
		"do-work/archive/UR-134/input.md":                      "",
		"src/do-work/working/REQ-632-not-the-queue.md":         "",
	} {
		match := requestPathPattern.FindStringSubmatch(path)
		gotId := ""
		if match != nil {
			gotId = match[1] + match[2]
		}
		if gotId != wantId {
			t.Errorf("%s attributed to %q, want %q", path, gotId, wantId)
		}
	}
}
