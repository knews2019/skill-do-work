package main

import (
	"fmt"
	"io"
	"strings"
)

// The `open-work` subcommand: the terminal answer to "what is in flight right
// now?", meant to be read in two seconds without a browser.
//
// It is deliberately NOT a second `summary`. Summary is the parser's smoke test —
// every column's count, completion anomalies, calendar entries, dependency edges,
// warnings — and its shape is asserted by tests and relayed by
// actions/forensics.md. This digest answers a different question and so carries
// per-ticket lines instead of counts alone: the open total (everything the board
// buckets outside the terminal-resolved set), every claimed REQ with its title,
// and every needs-input/blocked REQ with the status that put it there. Nothing
// terminal appears at all — recently-done and the calendar are history, not open
// work — which is why the recent-window is inert here.
//
// Read-only, like every subcommand but the two named write surfaces (the board's
// testing view and `next-version`).

// openWorkCounts is the headline breakdown: the open total plus the per-bucket
// split behind it. Open means "not terminally resolved" — bucketColumns puts
// every parsed ticket in exactly one of pending / claimed / needs-input-or-blocked
// / terminal, so summing the first three is exhaustive by construction rather
// than by a hand-maintained status list that would go stale the next time the
// Schema Read Contract grows a value (actions/work-reference.md).
type openWorkCounts struct {
	OpenTotal           int
	Pending             int
	PendingReady        int
	PendingWaiting      int
	Claimed             int
	NeedsInputOrBlocked int
}

// countOpenWork derives the headline breakdown from the bucketed columns.
func countOpenWork(board *Board) openWorkCounts {
	counts := openWorkCounts{
		Pending:             len(board.Columns.Pending),
		PendingReady:        len(board.Columns.PendingReady),
		PendingWaiting:      len(board.Columns.PendingWaiting),
		Claimed:             len(board.Columns.Claimed),
		NeedsInputOrBlocked: len(board.Columns.NeedsInputOrBlocked),
	}
	counts.OpenTotal = counts.Pending + counts.Claimed + counts.NeedsInputOrBlocked
	return counts
}

// writeOpenWorkDigest renders the whole digest. Split from the command wrapper
// (which owns flag parsing and os.Exit) so the output shape is directly
// assertable — the same split writeBoardSummary uses.
func writeOpenWorkDigest(outputWriter io.Writer, board *Board) {
	counts := countOpenWork(board)

	fmt.Fprintf(outputWriter, "queue-kanban open work: %d open %s\n", counts.OpenTotal, pluralizeRequestNoun(counts.OpenTotal))
	fmt.Fprintf(outputWriter, "  pending %d (%d ready, %d waiting) | claimed %d | needs-input/blocked %d\n",
		counts.Pending, counts.PendingReady, counts.PendingWaiting, counts.Claimed, counts.NeedsInputOrBlocked)

	writeClaimedSection(outputWriter, board.Columns.Claimed)
	writeNeedsInputSection(outputWriter, board.Columns.NeedsInputOrBlocked)

	// A count, not the warning text: the digest's job is to say "the parse found
	// problems, go read them", not to become summary. Silently dropping them
	// would be the one unacceptable option — a REQ with a malformed status is
	// exactly the kind of thing a fast check exists to notice.
	if len(board.Warnings) > 0 {
		fmt.Fprintf(outputWriter, "\nwarnings: %d (run `queue-kanban summary` for the details)\n", len(board.Warnings))
	}
}

// writeClaimedSection lists each claimed REQ as id + title — the work someone
// (or some session) already picked up.
func writeClaimedSection(outputWriter io.Writer, claimedTickets []*RequestTicket) {
	fmt.Fprintf(outputWriter, "\nclaimed (%d)\n", len(claimedTickets))
	if len(claimedTickets) == 0 {
		fmt.Fprintf(outputWriter, "  (none)\n")
		return
	}
	requestIdWidth := widestRequestId(claimedTickets)
	for _, claimedTicket := range claimedTickets {
		fmt.Fprintf(outputWriter, "  %-*s  %s\n", requestIdWidth, claimedTicket.RequestId, ticketTitleOrPlaceholder(claimedTicket))
	}
}

// writeNeedsInputSection lists each needs-input/blocked REQ with the status that
// parked it there, so the reader can tell a question waiting on them
// (pending-answers) from an external condition (blocked) without opening a file.
func writeNeedsInputSection(outputWriter io.Writer, needsInputTickets []*RequestTicket) {
	fmt.Fprintf(outputWriter, "\nneeds input / blocked (%d)\n", len(needsInputTickets))
	if len(needsInputTickets) == 0 {
		fmt.Fprintf(outputWriter, "  (none)\n")
		return
	}
	requestIdWidth := widestRequestId(needsInputTickets)
	statusLabelWidth := 0
	for _, needsInputTicket := range needsInputTickets {
		if labelWidth := len(openWorkStatusLabel(needsInputTicket)); labelWidth > statusLabelWidth {
			statusLabelWidth = labelWidth
		}
	}
	for _, needsInputTicket := range needsInputTickets {
		fmt.Fprintf(outputWriter, "  %-*s  %-*s  %s%s\n",
			requestIdWidth, needsInputTicket.RequestId,
			statusLabelWidth, openWorkStatusLabel(needsInputTicket),
			ticketTitleOrPlaceholder(needsInputTicket),
			blockedConditionSuffix(needsInputTicket))
	}
}

// openWorkStatusLabel is the status text shown for a needs-input/blocked ticket.
// An off-vocabulary status renders as the verbatim frontmatter value tagged
// invalid — bucketColumns parks unrecognized statuses in this column precisely so
// they stay visible, and collapsing them to a tidy label here would undo that.
func openWorkStatusLabel(ticket *RequestTicket) string {
	if ticket.StatusUnrecognized {
		return fmt.Sprintf("invalid:%s", ticket.OriginalStatus)
	}
	return ticket.Status
}

// blockedConditionSuffix appends the external condition a blocked REQ waits on
// when `blocked_by` names one — the single most useful word after the status
// itself. Display only: the digest never runs `blocked_check` (the work pipeline
// does), and blocked_by is never a dependency edge.
func blockedConditionSuffix(ticket *RequestTicket) string {
	if len(ticket.BlockedBy) == 0 {
		return ""
	}
	return fmt.Sprintf("  (waiting on: %s)", strings.Join(ticket.BlockedBy, ", "))
}

// ticketTitleOrPlaceholder keeps a titleless REQ from printing as a bare id with
// trailing whitespace.
func ticketTitleOrPlaceholder(ticket *RequestTicket) string {
	if strings.TrimSpace(ticket.Title) == "" {
		return "(untitled)"
	}
	return ticket.Title
}

// pluralizeRequestNoun keeps the headline reading like English at a count of one.
func pluralizeRequestNoun(requestCount int) string {
	if requestCount == 1 {
		return "REQ"
	}
	return "REQs"
}

// widestRequestId is the column width for the id gutter, so ids of different
// lengths (REQ-72 next to REQ-1318) still line up.
func widestRequestId(tickets []*RequestTicket) int {
	widestWidth := 0
	for _, ticket := range tickets {
		if idWidth := len(ticket.RequestId); idWidth > widestWidth {
			widestWidth = idWidth
		}
	}
	return widestWidth
}
