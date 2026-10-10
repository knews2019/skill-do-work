package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
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
// terminal appears at all — recently-done is history, not open work — which is
// why the recent-window is inert here. (The calendar is a different surface: it
// carries every REQ, open ones included, and is not a model for this digest.)
//
// This digest is read-only.
//
// `--format json` prints the same open tickets as machine facts for
// do-work-cli run-status, which cannot import this module: the column, the
// sentence bucketColumns wrote for it, and, for claimed tickets, the last
// activity attachRequestActivity computes for the HTML board. One home for the
// partition and the correlation; the reader only quotes them.

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
	PendingEarmarked    int
	Claimed             int
	NeedsInputOrBlocked int
}

// countOpenWork derives the headline breakdown from the bucketed columns.
func countOpenWork(board *Board) openWorkCounts {
	openCounts := openWorkCounts{
		Pending:             len(board.Columns.Pending),
		PendingReady:        len(board.Columns.PendingReady),
		PendingWaiting:      len(board.Columns.PendingWaiting),
		PendingEarmarked:    len(board.Columns.PendingEarmarked),
		Claimed:             len(board.Columns.Claimed),
		NeedsInputOrBlocked: len(board.Columns.NeedsInputOrBlocked),
	}
	openCounts.OpenTotal = openCounts.Pending + openCounts.Claimed + openCounts.NeedsInputOrBlocked
	return openCounts
}

// writeOpenWorkDigest renders the whole digest. Split from the command wrapper
// (which owns flag parsing and os.Exit) so the output shape is directly
// assertable — the same split writeBoardSummary uses.
func writeOpenWorkDigest(outputWriter io.Writer, board *Board) {
	openCounts := countOpenWork(board)

	fmt.Fprintf(outputWriter, "queue-kanban open work: %d open %s\n", openCounts.OpenTotal, pluralizeRequestNoun(openCounts.OpenTotal))
	fmt.Fprintf(outputWriter, "  pending %d (%d ready, %d waiting, %d earmarked) | claimed %d | needs-input/blocked %d\n",
		openCounts.Pending, openCounts.PendingReady, openCounts.PendingWaiting, openCounts.PendingEarmarked, openCounts.Claimed, openCounts.NeedsInputOrBlocked)

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

// openWorkFacts is the `open-work --format json` document. Keys are snake_case
// because the reader is do-work-cli, whose results use that case.
type openWorkFacts struct {
	GeneratedAt                string                 `json:"generated_at"`
	StaleClaimThresholdMinutes int                    `json:"stale_claim_threshold_minutes"`
	Requests                   []openWorkRequestFacts `json:"requests"`
}

type openWorkRequestFacts struct {
	Id                string   `json:"id"`
	Title             string   `json:"title"`
	Status            string   `json:"status"`
	Column            string   `json:"column"`
	PlacementReason   string   `json:"placement_reason"`
	UnmetDependencies []string `json:"unmet_dependencies"`
	AssignedTo        string   `json:"assigned_to"`
	LastActivityAt    string   `json:"last_activity_at,omitempty"`
	LastActivityKind  string   `json:"last_activity_kind,omitempty"`
	LastActivityPhase string   `json:"last_activity_phase,omitempty"`
}

// writeOpenWorkJSON renders every open ticket with the column the board put it
// in. Activity comes from the board's own correlation over one worktree-agent
// git read, the call attachVerifyFindingsAndRequestActivity makes.
func writeOpenWorkJSON(outputWriter io.Writer, board *Board, runner gitCommandRunner) error {
	var activityData generatedBoardData
	attachRequestActivity(&activityData, board, board.GeneratedAt, runner, readWorktreeAgentGitState(board.RepoRoot, runner))

	facts := openWorkFacts{
		GeneratedAt:                formatTimestamp(board.GeneratedAt),
		StaleClaimThresholdMinutes: int(staleClaimThreshold / time.Minute),
		Requests:                   []openWorkRequestFacts{},
	}
	for _, columnGroup := range []struct {
		columnName string
		tickets    []*RequestTicket
	}{
		{"pending-ready", board.Columns.PendingReady},
		{"pending-waiting", board.Columns.PendingWaiting},
		{"pending-earmarked", board.Columns.PendingEarmarked},
		{"claimed", board.Columns.Claimed},
		{"needs-input-or-blocked", board.Columns.NeedsInputOrBlocked},
	} {
		for _, ticket := range columnGroup.tickets {
			activity := activityData.RequestActivity[ticket.RequestId]
			facts.Requests = append(facts.Requests, openWorkRequestFacts{
				Id:                ticket.RequestId,
				Title:             ticket.Title,
				Status:            openWorkStatusLabel(ticket),
				Column:            columnGroup.columnName,
				PlacementReason:   ticket.PlacementReason,
				UnmetDependencies: append([]string{}, ticket.UnmetDependencies...),
				AssignedTo:        ticket.AssignedTo,
				LastActivityAt:    activity.LastActivityAt,
				LastActivityKind:  activity.LastActivityKind,
				LastActivityPhase: activity.LastActivityPhase,
			})
		}
	}
	encoder := json.NewEncoder(outputWriter)
	encoder.SetIndent("", "  ")
	return encoder.Encode(facts)
}
