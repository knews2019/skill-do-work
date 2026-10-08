package main

import (
	"errors"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Per-request activity evidence: when was a REQ last worked on, and what was
// the longest stretch nobody touched it?
//
// The answer is the union of two observed streams — the REQ's own lifecycle
// stamps (lifecycleTimestampFields) and the git commits correlated to it — so a
// claimed card shows a growing "last activity" number when work stalls, and a
// done REQ's drawer states its "Largest gap between events" instead of the board
// guessing a pause from the span alone. The read is cheap (one windowed `git
// log` plus the owned worktree-agent branch tips from worktreeAgentGitState) and
// is run per response, never inside serve's mtime cache: commits land without
// any do-work file changing mtime, which is the REQ-284 shape.

// gitCommandRunner runs one read-only git command against repoRoot and returns
// its stdout. Injectable so tests feed canned output and never spawn git.
type gitCommandRunner func(repoRoot string, arguments ...string) ([]byte, error)

// runGitCommand is the live gitCommandRunner.
func runGitCommand(repoRoot string, arguments ...string) ([]byte, error) {
	if !gitBinaryAvailable() {
		return nil, errors.New("git binary not available")
	}
	return exec.Command("git", append([]string{"-C", repoRoot}, arguments...)...).Output()
}

// requestActivity is what the collector observed for one REQ.
type requestActivity struct {
	// The newest union event. Set only for open (claimed) requests: a done card
	// shows nothing new, and "last activity" on a finished REQ answers nothing.
	LastActivityAt    time.Time
	LastActivityKind  string // "stamp" or "commit"
	LastActivityPhase string // lowercased phase label ("dispatch"), "" before any phase stamp
	// The largest gap between consecutive union events; nil with fewer than two.
	LargestGap *activityGap
}

// activityGap is the longest stretch between two consecutive activity events.
// FromPhase is the phase the gap starts in, ToPhase the phase of the event
// that ended it.
type activityGap struct {
	Minutes   float64
	FromPhase string
	ToPhase   string
}

// requestPathPattern matches a REQ file in queue, working, or archive — archive
// nests under a UR folder — and every run artifact named for a REQ. Anything
// else under do-work/ (reservation markers, the lessons index) is not evidence
// for one REQ.
var requestPathPattern = regexp.MustCompile(`^do-work/(?:(?:queue|working|archive)/(?:[^/]+/)*(REQ-\d+)-[^/]*\.md|runs/[^/]+/(REQ-\d+)(?:[^/\d][^/]*)?)$`)

// requestSubjectPrefixPattern matches every bracketed [REQ-NNN] token in a
// commit subject. An unbracketed id is prose, not attribution.
var requestSubjectPrefixPattern = regexp.MustCompile(`\[(REQ-\d+)\]`)

// correlatedCommit is one record of the windowed git log.
type correlatedCommit struct {
	hash         string
	committedAt  time.Time
	parentHashes []string
	subject      string
	touchedPaths []string
}

// parseCorrelationLog reads `git log --format=%H%x00%cI%x00%P%x00%s --name-only`
// output. A header line is the only line that carries NUL bytes; every other
// non-blank line is a path the preceding commit touched.
func parseCorrelationLog(logOutput []byte) []*correlatedCommit {
	var commits []*correlatedCommit
	var current *correlatedCommit
	for _, line := range strings.Split(string(logOutput), "\n") {
		if strings.Contains(line, "\x00") {
			fields := strings.SplitN(line, "\x00", 4)
			if len(fields) != 4 {
				current = nil
				continue
			}
			committedAt, parsed := parseTimestamp(fields[1])
			if !parsed {
				current = nil
				continue
			}
			current = &correlatedCommit{
				hash:         fields[0],
				committedAt:  committedAt,
				parentHashes: strings.Fields(fields[2]),
				subject:      fields[3],
			}
			commits = append(commits, current)
			continue
		}
		if path := strings.TrimSpace(line); path != "" && current != nil {
			current.touchedPaths = append(current.touchedPaths, path)
		}
	}
	return commits
}

// correlateCommitsToRequests attributes each logged commit to REQ ids by any of:
// a touched REQ path, a [REQ-NNN] subject token, or membership in
// <merge>^1..<merge>^2 of a two-parent commit already matched by the first two
// rules. Ancestry is derived from the logged %P graph, so it reaches only as far
// back as the log window — which is all the board asks about.
func correlateCommitsToRequests(logOutput []byte) map[string][]time.Time {
	commits := parseCorrelationLog(logOutput)
	commitByHash := map[string]*correlatedCommit{}
	for _, commit := range commits {
		commitByHash[commit.hash] = commit
	}

	requestIdsByHash := map[string]map[string]bool{}
	attribute := func(commitHash string, requestId string) {
		if requestIdsByHash[commitHash] == nil {
			requestIdsByHash[commitHash] = map[string]bool{}
		}
		requestIdsByHash[commitHash][requestId] = true
	}
	for _, commit := range commits {
		for _, path := range commit.touchedPaths {
			if match := requestPathPattern.FindStringSubmatch(path); match != nil {
				attribute(commit.hash, match[1]+match[2])
			}
		}
		for _, match := range requestSubjectPrefixPattern.FindAllStringSubmatch(commit.subject, -1) {
			attribute(commit.hash, match[1])
		}
	}

	// Ancestry runs off the DIRECT matches only, snapshotted here before any
	// range is attributed: the loop below mutates requestIdsByHash, so a merge
	// reached through another merge's range (a builder merging main) must not
	// read its range-given ids back as direct ones and widen the attribution.
	directIdsByMergeHash := map[string][]string{}
	for _, commit := range commits {
		if len(commit.parentHashes) != 2 {
			continue
		}
		for requestId := range requestIdsByHash[commit.hash] {
			directIdsByMergeHash[commit.hash] = append(directIdsByMergeHash[commit.hash], requestId)
		}
	}
	for _, commit := range commits {
		directIds := directIdsByMergeHash[commit.hash]
		if len(directIds) == 0 {
			continue
		}
		firstParentAncestry := loggedAncestry(commitByHash, commit.parentHashes[0], nil)
		for rangeHash := range loggedAncestry(commitByHash, commit.parentHashes[1], firstParentAncestry) {
			for _, requestId := range directIds {
				attribute(rangeHash, requestId)
			}
		}
	}

	instantsById := map[string][]time.Time{}
	for _, commit := range commits {
		for requestId := range requestIdsByHash[commit.hash] {
			instantsById[requestId] = append(instantsById[requestId], commit.committedAt)
		}
	}
	return instantsById
}

// loggedAncestry walks parents from startHash through the logged graph,
// stopping at commits outside the log window and at any hash in stopAt.
func loggedAncestry(commitByHash map[string]*correlatedCommit, startHash string, stopAt map[string]bool) map[string]bool {
	reached := map[string]bool{}
	pending := []string{startHash}
	for len(pending) > 0 {
		hash := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		commit, logged := commitByHash[hash]
		if !logged || reached[hash] || stopAt[hash] {
			continue
		}
		reached[hash] = true
		pending = append(pending, commit.parentHashes...)
	}
	return reached
}

// activityEvent is one entry of a REQ's union stream.
type activityEvent struct {
	instant time.Time
	kind    string // "stamp" or "commit"
	phase   string // the phase this event opens (phase stamps) or sits in
	opens   bool   // a phase stamp: its own label becomes the current phase
}

// collectRequestActivity unions each ticket's stamps with its correlated
// commits and reports the newest event and the largest gap. `since` bounds the
// git log; `now` drops events past the clock-skew allowance, which only a
// broken clock can produce. ownedTipInstantsById adds each live builder
// branch's tip — a builder's commits before hand-back exist only there — but
// only for a branch that owns its tip (worktreeAgentGitState). A failed git
// read leaves the stamps alone, which is still true evidence.
func collectRequestActivity(repoRoot string, tickets []*RequestTicket, since time.Time, now time.Time, runner gitCommandRunner, ownedTipInstantsById map[string][]time.Time) map[string]requestActivity {
	commitInstantsById := map[string][]time.Time{}
	logOutput, logError := runner(repoRoot, "log", "--since="+since.UTC().Format(time.RFC3339),
		"--format=%H%x00%cI%x00%P%x00%s", "--name-only")
	if logError == nil {
		commitInstantsById = correlateCommitsToRequests(logOutput)
	}
	for requestId, tipInstants := range ownedTipInstantsById {
		commitInstantsById[requestId] = append(commitInstantsById[requestId], tipInstants...)
	}

	latestCredibleInstant := now.Add(futureTimestampSkewAllowance)
	activityById := map[string]requestActivity{}
	for _, ticket := range tickets {
		events := requestActivityEvents(ticket, commitInstantsById[ticket.RequestId], latestCredibleInstant)
		if len(events) == 0 {
			continue
		}
		activity := requestActivity{}
		if ticket.Status == "claimed" {
			newest := events[len(events)-1]
			activity.LastActivityAt = newest.instant
			activity.LastActivityKind = newest.kind
			activity.LastActivityPhase = newest.phase
		}
		for index := 1; index < len(events); index++ {
			gapMinutes := events[index].instant.Sub(events[index-1].instant).Minutes()
			if activity.LargestGap == nil || gapMinutes > activity.LargestGap.Minutes {
				activity.LargestGap = &activityGap{
					Minutes:   gapMinutes,
					FromPhase: events[index-1].phase,
					ToPhase:   events[index].phase,
				}
			}
		}
		activityById[ticket.RequestId] = activity
	}
	return activityById
}

// requestActivityEvents builds one REQ's sorted union stream, phase-labelled.
// Only the work window counts: from claimed_at (the queue wait before a claim
// is not idle work) to the last of completed_at and release_at for a finished
// REQ (a later sibling commit touching this archived file is not idle time
// inside it). An open REQ's window runs to now.
func requestActivityEvents(ticket *RequestTicket, commitInstants []time.Time, latestCredibleInstant time.Time) []activityEvent {
	phaseLabelByField := map[string]string{}
	for _, milestone := range phaseMilestonesOf(ticket) {
		phaseLabelByField[milestone.fieldName] = strings.ToLower(milestone.label)
	}

	windowStart, hasWindowStart := parseTimestamp(ticket.ClaimedAt)
	windowEnd := latestCredibleInstant
	if completedAt, parsed := parseTimestamp(ticket.CompletedAt); parsed && ticket.Status != "claimed" {
		windowEnd = completedAt
		if releaseAt, releaseParsed := parseTimestamp(ticket.ReleaseAt); releaseParsed && releaseAt.After(windowEnd) {
			windowEnd = releaseAt
		}
		if windowEnd.After(latestCredibleInstant) {
			windowEnd = latestCredibleInstant
		}
	}
	insideWindow := func(instant time.Time) bool {
		return !(hasWindowStart && instant.Before(windowStart)) && !instant.After(windowEnd)
	}

	var events []activityEvent
	for _, stamp := range lifecycleTimestampFields(ticket) {
		instant, parsed := parseTimestamp(stamp.RawValue)
		if !parsed || !insideWindow(instant) {
			continue
		}
		label, opensPhase := phaseLabelByField[stamp.FieldName]
		events = append(events, activityEvent{instant: instant, kind: "stamp", phase: label, opens: opensPhase})
	}
	for _, instant := range commitInstants {
		if insideWindow(instant) {
			events = append(events, activityEvent{instant: instant, kind: "commit"})
		}
	}
	// Stamps sort before commits at the same instant, so a commit made at the
	// moment a phase opened reads as inside that phase.
	sort.SliceStable(events, func(left, right int) bool {
		if !events[left].instant.Equal(events[right].instant) {
			return events[left].instant.Before(events[right].instant)
		}
		return events[left].kind == "stamp" && events[right].kind != "stamp"
	})
	currentPhase := ""
	for index := range events {
		if events[index].opens {
			currentPhase = events[index].phase
		}
		events[index].phase = currentPhase
	}
	return events
}

// generatedRequestActivity is one entry of the payload's requestActivity map.
type generatedRequestActivity struct {
	LastActivityAt     string                `json:"lastActivityAt,omitempty"`
	LastActivityKind   string                `json:"lastActivityKind,omitempty"`
	LastActivityPhase  string                `json:"lastActivityPhase,omitempty"`
	LargestActivityGap *generatedActivityGap `json:"largestActivityGap,omitempty"`
}

type generatedActivityGap struct {
	Minutes   float64 `json:"minutes"`
	FromPhase string  `json:"fromPhase"`
	ToPhase   string  `json:"toPhase"`
}

// attachRequestActivity collects activity for the board's claimed and
// recently-done REQs and folds it into the payload, over a worktree and branch
// read the caller already made (attachVerifyFindingsAndRequestActivity shares
// it with the verify probes). Older archived REQs are left out: reading their
// commits would mean walking history back to the oldest claim on the board.
func attachRequestActivity(data *generatedBoardData, board *Board, now time.Time, runner gitCommandRunner, gitState worktreeAgentGitState) {
	trackedTickets := append(append([]*RequestTicket{}, board.Columns.Claimed...), board.Columns.RecentlyDone...)
	if len(trackedTickets) == 0 {
		return
	}
	since := now.Add(-board.RecentWindow)
	for _, ticket := range trackedTickets {
		if claimedAt, parsed := parseTimestamp(ticket.ClaimedAt); parsed && claimedAt.Before(since) {
			since = claimedAt
		}
	}

	for requestId, activity := range collectRequestActivity(board.RepoRoot, trackedTickets, since, now, runner, gitState.ownedTipInstantsById) {
		entry := generatedRequestActivity{
			LastActivityAt:    formatTimestamp(activity.LastActivityAt),
			LastActivityKind:  activity.LastActivityKind,
			LastActivityPhase: activity.LastActivityPhase,
		}
		if activity.LargestGap != nil {
			entry.LargestActivityGap = &generatedActivityGap{
				Minutes:   activity.LargestGap.Minutes,
				FromPhase: activity.LargestGap.FromPhase,
				ToPhase:   activity.LargestGap.ToPhase,
			}
		}
		if entry.LastActivityAt == "" && entry.LargestActivityGap == nil {
			continue
		}
		if data.RequestActivity == nil {
			data.RequestActivity = map[string]generatedRequestActivity{}
		}
		data.RequestActivity[requestId] = entry
	}
}
