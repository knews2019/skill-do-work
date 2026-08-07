package main

import (
	"strings"
	"testing"
)

// renderOpenWorkDigest is the shared render-to-string helper for these tests.
func renderOpenWorkDigest(t *testing.T, board *Board) string {
	t.Helper()
	var digestBuffer strings.Builder
	writeOpenWorkDigest(&digestBuffer, board)
	return digestBuffer.String()
}

// The open total must be exactly the non-terminal tickets — no more (a completed
// REQ is not open work) and no fewer (a status outside the vocabulary is parked in
// needs-input/blocked and must still be counted, never silently dropped). Derived
// from the parsed statuses rather than from a literal so the assertion follows the
// Schema Read Contract instead of freezing today's fixture arithmetic.
func TestOpenWorkCountsEveryNonTerminalTicketExactlyOnce(t *testing.T) {
	board := syntheticBoard(t)
	openCounts := countOpenWork(board)

	expectedOpenTotal := 0
	for _, ticket := range board.AllRequests {
		if !isTerminalResolvedStatus(ticket.Status) {
			expectedOpenTotal++
		}
	}

	if openCounts.OpenTotal != expectedOpenTotal {
		t.Fatalf("open total = %d, want %d (the non-terminal tickets in the tree)", openCounts.OpenTotal, expectedOpenTotal)
	}
	if openCounts.Pending+openCounts.Claimed+openCounts.NeedsInputOrBlocked != openCounts.OpenTotal {
		t.Fatalf("bucket split %d+%d+%d does not sum to the open total %d",
			openCounts.Pending, openCounts.Claimed, openCounts.NeedsInputOrBlocked, openCounts.OpenTotal)
	}
	if openCounts.PendingReady+openCounts.PendingWaiting != openCounts.Pending {
		t.Fatalf("ready %d + waiting %d does not sum to pending %d", openCounts.PendingReady, openCounts.PendingWaiting, openCounts.Pending)
	}
}

// The headline carries the open total plus the per-bucket breakdown — the whole
// point of the digest over a browser tab is that the first line answers the
// question.
func TestOpenWorkDigestHeadlineCarriesTotalAndBreakdown(t *testing.T) {
	board := syntheticBoard(t)
	digestOutput := renderOpenWorkDigest(t, board)

	firstLine := strings.SplitN(digestOutput, "\n", 2)[0]
	if !strings.Contains(firstLine, "3 open REQs") {
		t.Fatalf("headline does not report the open total (1 pending + 1 claimed + 1 blocked): %q", firstLine)
	}
	for _, breakdownFragment := range []string{
		"pending 1 (1 ready, 0 waiting)",
		"claimed 1",
		"needs-input/blocked 1",
	} {
		if !strings.Contains(digestOutput, breakdownFragment) {
			t.Fatalf("digest is missing the breakdown fragment %q:\n%s", breakdownFragment, digestOutput)
		}
	}
}

// Claimed REQs are listed with their titles: an id alone does not tell the reader
// what is in flight, which is the reason this subcommand exists next to summary.
func TestOpenWorkDigestListsClaimedRequestsWithTitles(t *testing.T) {
	board := syntheticBoard(t)
	digestOutput := renderOpenWorkDigest(t, board)

	claimedSection := sectionAfterHeading(t, digestOutput, "claimed (1)")
	if !strings.Contains(claimedSection, "REQ-9002") || !strings.Contains(claimedSection, "Fixture REQ-9002") {
		t.Fatalf("claimed section does not carry the claimed REQ id AND title:\n%s", claimedSection)
	}
	// The pending REQ must NOT be listed by title — the digest lists claimed and
	// needs-input tickets per-ticket and pending as a count, so it stays scannable.
	if strings.Contains(digestOutput, "REQ-9001") {
		t.Fatalf("digest lists pending REQ-9001 per-ticket; pending is a count only:\n%s", digestOutput)
	}
	// Nothing terminal belongs in an open-work digest.
	for _, terminalRequestId := range []string{"REQ-9003", "REQ-9004", "REQ-9005"} {
		if strings.Contains(digestOutput, terminalRequestId) {
			t.Fatalf("digest mentions terminal %s — recently-done is history, not open work:\n%s", terminalRequestId, digestOutput)
		}
	}
}

// A needs-input/blocked line must name the status that parked it there, so a
// question waiting on the user (pending-answers) is distinguishable from an
// external condition (blocked) without opening the file. The external condition
// itself rides along when blocked_by names one.
func TestOpenWorkDigestShowsNeedsInputStatusAndBlockedCondition(t *testing.T) {
	board := syntheticBoard(t)
	digestOutput := renderOpenWorkDigest(t, board)

	needsInputSection := sectionAfterHeading(t, digestOutput, "needs input / blocked (1)")
	for _, expectedFragment := range []string{"REQ-9006", "blocked", "waiting on: LM Studio running locally"} {
		if !strings.Contains(needsInputSection, expectedFragment) {
			t.Fatalf("needs-input section is missing %q:\n%s", expectedFragment, needsInputSection)
		}
	}
	// blocked_check is display-only for the board and not shown here at all — the
	// work pipeline runs the probe, never this digest.
	if strings.Contains(digestOutput, "curl -sf") {
		t.Fatalf("digest printed the blocked_check probe; it must not suggest the digest runs it:\n%s", digestOutput)
	}
}

// An off-vocabulary status keeps its verbatim frontmatter text, tagged invalid.
// bucketColumns parks unrecognized statuses in this column so they stay visible;
// rendering them as a tidy label would undo that.
func TestOpenWorkDigestFlagsUnrecognizedStatusVerbatim(t *testing.T) {
	boardWithInvalidStatus := &Board{}
	boardWithInvalidStatus.Columns.NeedsInputOrBlocked = []*RequestTicket{{
		RequestId:          "REQ-4242",
		Title:              "Half-migrated ticket",
		Status:             "in-progress",
		OriginalStatus:     "in-progress",
		StatusUnrecognized: true,
	}}

	digestOutput := renderOpenWorkDigest(t, boardWithInvalidStatus)
	if !strings.Contains(digestOutput, "invalid:in-progress") {
		t.Fatalf("unrecognized status is not flagged with its verbatim value:\n%s", digestOutput)
	}
}

// Empty sections print "(none)" rather than vanishing: a missing section reads as
// "I did not check", which is the opposite of what a fast confidence check is for.
func TestOpenWorkDigestPrintsEmptySectionsExplicitly(t *testing.T) {
	digestOutput := renderOpenWorkDigest(t, &Board{})

	if !strings.Contains(digestOutput, "0 open REQs") {
		t.Fatalf("empty board does not report zero open REQs:\n%s", digestOutput)
	}
	if strings.Count(digestOutput, "(none)") != 2 {
		t.Fatalf("both empty sections must print (none):\n%s", digestOutput)
	}
}

// Parse warnings are surfaced as a count with a pointer at summary — never
// dropped, never re-implemented here.
func TestOpenWorkDigestPointsAtSummaryForWarnings(t *testing.T) {
	boardWithoutWarnings := renderOpenWorkDigest(t, &Board{})
	if strings.Contains(boardWithoutWarnings, "warnings:") {
		t.Fatalf("clean board printed a warnings line:\n%s", boardWithoutWarnings)
	}

	boardWithWarnings := &Board{Warnings: []string{"REQ-1 has unrecognized status \"wip\""}}
	digestOutput := renderOpenWorkDigest(t, boardWithWarnings)
	if !strings.Contains(digestOutput, "warnings: 1") || !strings.Contains(digestOutput, "queue-kanban summary") {
		t.Fatalf("warnings are not surfaced with a pointer at summary:\n%s", digestOutput)
	}
}

// sectionAfterHeading returns the digest text from a section heading up to the
// next blank line, so a section's assertions cannot accidentally be satisfied by
// text belonging to the other section.
func sectionAfterHeading(t *testing.T, digestOutput string, sectionHeading string) string {
	t.Helper()
	headingIndex := strings.Index(digestOutput, sectionHeading)
	if headingIndex < 0 {
		t.Fatalf("digest has no %q section:\n%s", sectionHeading, digestOutput)
	}
	sectionBody := digestOutput[headingIndex+len(sectionHeading):]
	if sectionEnd := strings.Index(sectionBody, "\n\n"); sectionEnd >= 0 {
		sectionBody = sectionBody[:sectionEnd]
	}
	return sectionBody
}
