package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CompletionTimeSource records how a RequestTicket's completion instant was
// resolved, so a consumer can tell an authoritative frontmatter timestamp from a
// best-effort git or filesystem fallback.
type CompletionTimeSource string

const (
	// CompletionFromFrontmatter means the completed_at frontmatter value was used.
	CompletionFromFrontmatter CompletionTimeSource = "frontmatter"
	// CompletionFromGitLog means the commit hash's git committer date was used.
	CompletionFromGitLog CompletionTimeSource = "git"
	// CompletionUnresolved means no completion instant could be determined.
	// (File mtime is deliberately NOT a source: a clone, checkout, or tarball
	// extraction resets every mtime to "now", which stamped months-old REQs as
	// completed today and buried the real recent work.)
	CompletionUnresolved CompletionTimeSource = "unresolved"
)

// undatedCalendarDayKey buckets completed REQs whose completion instant could
// not be resolved. They stay visible on the calendar (never silently dropped)
// but carry no fabricated date.
const undatedCalendarDayKey = "undated"

// reservationStaleAfter is how old a reservation (status: reserved, written by
// do-work reserve — actions/reserve.md) may grow before the board flags it as
// stale and suggests recategorizing. Mirrors the 24h threshold in
// actions/work.md Step 1's stale-reservation check — keep the two in lock-step.
const reservationStaleAfter = 24 * time.Hour

// futureTimestampSkewAllowance is how far past the board's `now` a frontmatter
// timestamp may parse before it is flagged as future-dated. Two minutes absorbs
// ordinary clock skew between machines; anything beyond it is almost always a
// session that stamped local wall-clock time with a `Z` suffix (the Timestamp
// rule in actions/work-reference.md requires the current UTC instant). Mirrored
// by futureInstantSkewAllowanceMs in web/board.js — keep the two in lock-step.
const futureTimestampSkewAllowance = 2 * time.Minute

// RequestTicket is one parsed REQ-*.md file: its frontmatter fields (with status
// normalized and commit-hash variants collapsed), the raw Markdown body kept for
// later HTML rendering, where it was found in the tree, and how its completion
// instant was resolved.
type RequestTicket struct {
	RequestId      string // canonical "REQ-NNNN" (frontmatter id, else derived from the filename)
	Title          string
	Status         string // normalized status (complete/done/finished/closed → completed)
	OriginalStatus string // verbatim frontmatter status before normalization

	// Set by bucketColumns when the normalized status falls outside the Schema
	// Read Contract vocabulary (actions/work-reference.md). The ticket is parked
	// in Needs input / Blocked and the frontend highlights it as invalid with a
	// fix prompt — never silently dropped.
	StatusUnrecognized bool

	CreatedAt   string // raw frontmatter timestamp text, "" when absent
	ClaimedAt   string
	CompletedAt string // raw frontmatter completed_at text, "" when absent

	ReservedFor string // reserve action (do-work reserve): owning worktree/cloud-session label, "" when absent
	ReservedAt  string // raw frontmatter reserved_at text, "" when absent

	// Derived by bucketColumns for status "reserved": true when reserved_at is
	// missing, unparseable, or more than reservationStaleAfter before the
	// board's `now` — the owning session may be dead, so the frontend shows a
	// recategorize hint. Never read from frontmatter.
	ReservationStale bool

	CommitHash      string // resolved from commit / commit_hash / green_commit / commit_green / impl_commit
	CommitHashField string // the frontmatter key CommitHash came from, "" when absent
	UserRequestId   string // "UR-NNN" upward pointer (the reliable REQ→UR link), "" when absent
	Domain          string

	// Testing-track placeholders written by the board's testing view (see
	// testing.go). Orthogonal to Status: the work pipeline never reads them and
	// the board never writes Status.
	TestingStatus         string // normalized testing status ("" = not tested yet)
	OriginalTestingStatus string // verbatim frontmatter testing_status before normalization
	// Set when a non-empty testing_status normalizes to nothing in the canonical
	// vocabulary — the ticket renders as not-yet-tested with an invalid flag and
	// a data warning, never silently.
	TestingStatusUnrecognized bool
	TestedBy                  string // tester profile name, "" when absent
	TestingUpdatedAt          string // raw frontmatter timestamp text, "" when absent
	TestingFeedback           string // feedback text (present while testing_status is returned)

	DependsOn []string // canonical dependency REQ ids (depends_on wins over legacy dependencies)
	// blocked_by names the external condition a `status: blocked` REQ waits on
	// (e.g. "LM Studio running locally"). Free text is the modern shape and
	// parses as a one-element list; a legacy id-LIST value renders joined for
	// display only — it is NOT a dependency edge (dependency gating is DependsOn).
	BlockedBy    []string // external-condition text (or legacy ids), kept distinct from DependsOn
	BlockedAt    string   // raw frontmatter timestamp text for the blocked flip, "" when absent
	BlockedCheck string   // optional shell probe command (display only; the pipeline, not the board, runs it), "" when absent
	Related      []string // soft relations (not dependency edges)

	// Derived by annotateDependencyState after every ticket is parsed — never
	// read from frontmatter.
	UnmetDependencies []string // DependsOn entries that have not reached terminal success (a dangling id counts as unmet)
	Dependents        []string // REQ ids whose depends_on names this ticket, in id order — the reverse edge

	Route string
	Batch string

	BodyMarkdown string // raw Markdown body after the closing frontmatter fence

	FilePath    string // absolute path on disk
	TreeSection string // "queue" | "working" | "archive"

	CompletionTime       time.Time            // resolved completion instant (zero when unresolved)
	CompletionTimeSource CompletionTimeSource // how CompletionTime was resolved

	// Set by buildBoard (via detectCompletionAnomaly) for terminal-resolved
	// tickets whose completion bookkeeping is broken: no completion instant at
	// all, a completed_at that fails parseTimestamp, or a consulted commit-hash
	// field git cannot resolve. Anomalous tickets surface in the dedicated
	// CompletionAnomalies column — never sorted or counted as if completed
	// "now" (no fabricated instant), and never silently dropped.
	CompletionAnomaly       bool
	CompletionAnomalyReason string // names the broken field(s); "" when no anomaly

	// Derived by detectFutureTimestampFields against the board's `now` — never
	// read from frontmatter. Each entry is "<field> <raw value>" for a stamp
	// that parses to later than now + futureTimestampSkewAllowance — the
	// signature of local wall-clock time written with a `Z` suffix. The
	// frontend badges the card, and a board warning names the fix; nil when
	// every stamp is sane.
	FutureTimestampFields []string
}

// UserRequestTicket is one parsed UR input.md plus the REQ ids grouped under it.
// The grouping is built from each REQ's user_request upward pointer (the reliable
// link), so RequestIds is populated even for a UR whose input.md is missing — in
// which case InputFilePresent is false and the node is synthesized.
type UserRequestTicket struct {
	UserRequestId    string
	Title            string
	CreatedAt        string
	BodyMarkdown     string
	FilePath         string
	InputFilePresent bool
	RequestIds       []string // REQ ids that point here via user_request, in id order
}

// DependencyEdge is one resolved depends_on relationship: FromRequestId depends
// on ToRequestId. Resolved reports whether ToRequestId matched a parsed ticket
// (a dependency may point at an id not present in the current tree).
type DependencyEdge struct {
	FromRequestId string
	ToRequestId   string
	Resolved      bool
}

// DependencyGraph holds every depends_on edge plus forward and reverse adjacency
// lookups keyed by REQ id.
type DependencyGraph struct {
	Edges      []DependencyEdge
	DependsOn  map[string][]string // RequestId → ids it depends on
	Dependents map[string][]string // RequestId → ids that depend on it
}

// CalendarEntry places one terminally resolved REQ (completed*/cancelled) on
// the completion timeline, recording the resolved instant, how it was
// resolved, and a UTC day bucket key.
type CalendarEntry struct {
	RequestId      string
	CompletionTime time.Time
	TimeSource     CompletionTimeSource
	DayKey         string // "2006-01-02" UTC bucket
}

// QueueNote is one line of do-work/notes.md — a lightweight, dated next-step
// hint written by `do-work note`. A note is deliberately NOT a REQ: it has no
// frontmatter, no status, and no id, so it never enters a column, the calendar,
// or the dependency graph. The board surfaces notes the way `do-work roadmap`
// does — verbatim, in append order, above the queue.
type QueueNote struct {
	NoteDate string // "YYYY-MM-DD" from the standard `- [date] text` prefix, "" when the line carries none
	NoteText string // the note text with its bullet and date prefix stripped
}

// BoardColumns holds the active-work buckets. Completed REQs older than the
// recent window are NOT represented here — they live in Board.Calendar.
type BoardColumns struct {
	Pending             []*RequestTicket // status pending (the union of PendingReady and PendingWaiting)
	PendingReady        []*RequestTicket // pending with every depends_on target at terminal success — actionable now
	PendingWaiting      []*RequestTicket // pending with at least one unmet dependency — not yet actionable
	Claimed             []*RequestTicket // status claimed
	NeedsInputOrBlocked []*RequestTicket // pending-answers / blocked / blocked-* / failed
	RecentlyDone        []*RequestTicket // completed*/cancelled whose completion instant is within the window

	// Terminal-resolved tickets flagged CompletionAnomaly. Surfaced in EVERY
	// mode regardless of the recent window and no matter how old — the missing
	// completed_at is a bookkeeping bug the user wants to see and fix. A
	// ticket dated via git despite a broken field appears here AND (window
	// permitting) in RecentlyDone.
	CompletionAnomalies []*RequestTicket
}

// Board is the complete parsed model of the do-work queue: every ticket, the
// active board columns, the UR grouping, the dependency graph, and the
// completion calendar. It is the stable contract the generate/serve stages
// consume.
type Board struct {
	GeneratedAt  time.Time     // the `now` the board was built against
	RecentWindow time.Duration // window used to populate Columns.RecentlyDone
	ProjectName  string        // base name of the repo root (the parent project this board belongs to)

	AllRequests  []*RequestTicket          // every parsed REQ, in id order
	RequestsById map[string]*RequestTicket // RequestId → ticket (first occurrence wins)

	UserRequests     []*UserRequestTicket          // every UR (parsed or synthesized), in id order
	UserRequestsById map[string]*UserRequestTicket // UserRequestId → ticket

	Columns         BoardColumns
	DependencyGraph DependencyGraph
	Calendar        []CalendarEntry // terminally resolved REQs sorted most-recent-first (undated entries last)
	Notes           []QueueNote     // do-work/notes.md lines in append order (nil when the file is absent or empty)

	TestingProfiles []string // do-work/testers.md bullet lines in file order (nil when the file is absent or empty)

	Warnings []string // data-shape warnings (e.g. duplicate ids, unrecognized statuses, future-dated stamps) — surfaced, never silently dropped
}

// gitCommitDateLookup resolves a commit hash to its committer date. It is an
// injection point: the live path uses lookupGitCommitDate; tests pass a stub so
// completion resolution is deterministic and does not shell out to git.
type gitCommitDateLookup func(repoRoot string, commitHash string) (time.Time, bool)

// LoadBoard is the entry point the rest of the tool builds on. It resolves the
// repo root (walking up from the working directory when repoRootOverride is
// empty), walks the do-work tree, parses every REQ and UR, and assembles the
// board model. `now` anchors the recently-done window; recentWindow sizes it.
func LoadBoard(repoRootOverride string, now time.Time, recentWindow time.Duration) (*Board, error) {
	repoRoot, resolveError := resolveRepoRootOrDefault(repoRootOverride)
	if resolveError != nil {
		return nil, resolveError
	}
	return buildBoard(repoRoot, now, recentWindow, lookupGitCommitDate)
}

// buildBoard is the testable core of LoadBoard with the git lookup injected.
func buildBoard(repoRoot string, now time.Time, recentWindow time.Duration, gitLookup gitCommitDateLookup) (*Board, error) {
	discovered, enumerateError := enumerateDoWorkTree(repoRoot)
	if enumerateError != nil {
		return nil, enumerateError
	}

	board := &Board{
		GeneratedAt:      now,
		RecentWindow:     recentWindow,
		ProjectName:      deriveProjectName(repoRoot),
		RequestsById:     map[string]*RequestTicket{},
		UserRequestsById: map[string]*UserRequestTicket{},
	}

	var parsedTickets []*RequestTicket
	for _, reference := range discovered.RequestFiles {
		ticket, parseError := parseRequestTicket(reference.AbsolutePath, reference.TreeSection)
		if parseError != nil {
			continue // best-effort: skip an unreadable REQ file
		}
		parsedTickets = append(parsedTickets, ticket)
	}

	// Keep exactly one ticket per REQ id (with a warning per duplicate) BEFORE
	// building any view, so the columns, calendar, UR groups, and the id-keyed
	// JSON map all render the same copy instead of contradicting each other.
	board.AllRequests, board.Warnings = dedupeTicketsByRequestId(parsedTickets)

	for _, ticket := range board.AllRequests {
		if isTerminalResolvedStatus(ticket.Status) {
			completionTime, completionSource := resolveCompletionTime(ticket, repoRoot, gitLookup)
			ticket.CompletionTime = completionTime
			ticket.CompletionTimeSource = completionSource
			ticket.CompletionAnomaly, ticket.CompletionAnomalyReason = detectCompletionAnomaly(ticket)
		} else {
			ticket.CompletionTimeSource = CompletionUnresolved
		}
		ticket.FutureTimestampFields = detectFutureTimestampFields(ticket, now)
		if len(ticket.FutureTimestampFields) > 0 {
			board.Warnings = append(board.Warnings, fmt.Sprintf(
				"%s has future-dated timestamp(s): %s — later than the board's generation time (2min clock-skew allowance); likely local wall-clock time stamped with a Z suffix; fix: rewrite with the current UTC instant (date -u +%%Y-%%m-%%dT%%H:%%M:%%SZ)",
				ticket.RequestId, strings.Join(ticket.FutureTimestampFields, ", ")))
		}
		board.RequestsById[ticket.RequestId] = ticket
	}
	sortRequestTickets(board.AllRequests)

	for _, userRequestPath := range discovered.UserRequestFiles {
		userRequestTicket, parseError := parseUserRequestTicket(userRequestPath)
		if parseError != nil {
			continue
		}
		if _, exists := board.UserRequestsById[userRequestTicket.UserRequestId]; exists {
			continue // first input.md for a UR wins
		}
		board.UserRequestsById[userRequestTicket.UserRequestId] = userRequestTicket
		board.UserRequests = append(board.UserRequests, userRequestTicket)
	}

	linkRequestsToUserRequests(board)

	sortUserRequestTickets(board.UserRequests)
	for _, userRequestTicket := range board.UserRequests {
		sortRequestIdList(userRequestTicket.RequestIds)
	}

	board.DependencyGraph = buildDependencyGraph(board.AllRequests, board.RequestsById)

	// Dependency state must be annotated BEFORE bucketing: the Pending column is
	// split into ready vs waiting by each ticket's UnmetDependencies.
	board.Warnings = append(board.Warnings, annotateDependencyState(board)...)

	columns, columnWarnings := bucketColumns(board.AllRequests, now, recentWindow)
	board.Columns = columns
	board.Warnings = append(board.Warnings, columnWarnings...)
	board.Warnings = append(board.Warnings, collectTestingWarnings(board.AllRequests)...)
	board.Calendar = buildCalendar(board.AllRequests)
	board.Notes = loadQueueNotes(discovered.NotesFilePath)
	board.TestingProfiles = loadTestingProfiles(discovered.TestersFilePath)

	return board, nil
}

// loadQueueNotes reads do-work/notes.md into one QueueNote per note BULLET,
// preserving append order (the user curates the file by hand; `do-work note` only
// appends, so file order IS chronological order). It is best-effort: an absent or
// unreadable file yields no notes rather than failing the board build.
//
// Only bullet lines are notes. Real notes.md files in the wild carry a `#`
// heading, a wrapped prose preamble, and HTML-comment blocks recording pruned
// entries — an earlier "every non-blank line is a note" read rendered all of
// that as boilerplate notes. HTML comments are stripped BEFORE the bullet test,
// because a pruned-entries comment block is itself full of bullets that must
// not resurface on the board.
func loadQueueNotes(notesFilePath string) []QueueNote {
	if notesFilePath == "" {
		return nil
	}
	contentBytes, readError := os.ReadFile(notesFilePath)
	if readError != nil {
		return nil
	}

	var notes []QueueNote
	insideHtmlComment := false
	for _, rawLine := range strings.Split(string(contentBytes), "\n") {
		var visibleText string
		visibleText, insideHtmlComment = stripHtmlComments(rawLine, insideHtmlComment)

		trimmedLine := strings.TrimSpace(visibleText)
		if trimmedLine == "" {
			continue
		}
		note, isNote := parseQueueNoteLine(trimmedLine)
		if !isNote {
			continue
		}
		notes = append(notes, note)
	}
	return notes
}

// stripHtmlComments removes every `<!-- ... -->` span from one line, given
// whether the previous line left an unterminated comment open. It returns the
// visible remainder and the comment state for the next line. A comment that
// opens and closes several times on one line is handled by the loop; a comment
// that never closes swallows the rest of the file, which is what the Markdown
// renderer would do too.
func stripHtmlComments(rawLine string, insideHtmlComment bool) (string, bool) {
	const commentOpenMarker = "<!--"
	const commentCloseMarker = "-->"

	var visibleBuilder strings.Builder
	remainingText := rawLine
	for remainingText != "" {
		if insideHtmlComment {
			closeIndex := strings.Index(remainingText, commentCloseMarker)
			if closeIndex < 0 {
				return visibleBuilder.String(), true // comment continues past this line
			}
			remainingText = remainingText[closeIndex+len(commentCloseMarker):]
			insideHtmlComment = false
			continue
		}
		openIndex := strings.Index(remainingText, commentOpenMarker)
		if openIndex < 0 {
			visibleBuilder.WriteString(remainingText)
			break
		}
		visibleBuilder.WriteString(remainingText[:openIndex])
		remainingText = remainingText[openIndex+len(commentOpenMarker):]
		insideHtmlComment = true
	}
	return visibleBuilder.String(), insideHtmlComment
}

// parseQueueNoteLine splits one notes.md line into its optional date prefix and
// its text, reporting whether the line is a note at all. The canonical shape
// written by `do-work note` is `- [YYYY-MM-DD] text`; the bullet is what marks a
// line as a note, so a heading, a preamble sentence, a horizontal rule, or a
// frontmatter fence is skipped rather than rendered. The DATE is still optional —
// a hand-typed bullet whose date prefix drifted keeps its text.
func parseQueueNoteLine(trimmedLine string) (QueueNote, bool) {
	noteText := ""
	foundBullet := false
	for _, bulletPrefix := range []string{"- ", "* ", "+ "} {
		if strings.HasPrefix(trimmedLine, bulletPrefix) {
			noteText = strings.TrimSpace(strings.TrimPrefix(trimmedLine, bulletPrefix))
			foundBullet = true
			break
		}
	}
	if !foundBullet || noteText == "" {
		return QueueNote{}, false
	}

	noteDate := ""
	if strings.HasPrefix(noteText, "[") {
		closingBracketIndex := strings.Index(noteText, "]")
		if closingBracketIndex > 1 {
			candidateDate := noteText[1:closingBracketIndex]
			if isBareDateText(candidateDate) {
				noteDate = candidateDate
				noteText = strings.TrimSpace(noteText[closingBracketIndex+1:])
			}
		}
	}
	return QueueNote{NoteDate: noteDate, NoteText: noteText}, true
}

// isBareDateText reports whether text is exactly a `YYYY-MM-DD` calendar date.
// The length check is what keeps parseTimestamp's other accepted layouts (and a
// Markdown task-list marker like `[ ]` or `[x]`) from being mistaken for a note date.
func isBareDateText(text string) bool {
	if len(text) != len("2006-01-02") {
		return false
	}
	_, parsedOk := parseTimestamp(text)
	return parsedOk
}

// treeSectionPrecedence orders tree sections for duplicate-id resolution: the
// active copy (queue, then working) wins over the archive copy — it is the copy
// the work pipeline would act on, and it is the copy cleanup marks
// `blocked-archive-collision`, which lands the collision in the visible
// Needs-input/Blocked column.
func treeSectionPrecedence(treeSection string) int {
	switch treeSection {
	case "queue":
		return 0
	case "working":
		return 1
	case "archive":
		return 2
	default:
		return 3
	}
}

// dedupeTicketsByRequestId keeps one ticket per REQ id (queue > working >
// archive precedence) and reports every duplicate as a warning naming both
// file paths. Without this, an id present in two tree sections rendered in two
// views while the id-keyed JSON map could only carry one copy's content.
func dedupeTicketsByRequestId(parsedTickets []*RequestTicket) ([]*RequestTicket, []string) {
	winnersByRequestId := map[string]*RequestTicket{}
	var orderedWinners []*RequestTicket
	var duplicateWarnings []string

	for _, ticket := range parsedTickets {
		existing, exists := winnersByRequestId[ticket.RequestId]
		if !exists {
			winnersByRequestId[ticket.RequestId] = ticket
			orderedWinners = append(orderedWinners, ticket)
			continue
		}
		keptTicket, ignoredTicket := existing, ticket
		if treeSectionPrecedence(ticket.TreeSection) < treeSectionPrecedence(existing.TreeSection) {
			keptTicket, ignoredTicket = ticket, existing
		}
		if keptTicket != existing {
			winnersByRequestId[ticket.RequestId] = keptTicket
			for winnerIndex, winner := range orderedWinners {
				if winner == existing {
					orderedWinners[winnerIndex] = keptTicket
					break
				}
			}
		}
		duplicateWarnings = append(duplicateWarnings, fmt.Sprintf(
			"duplicate REQ id %s: showing the %s copy (%s); ignoring the %s copy (%s)",
			ticket.RequestId, keptTicket.TreeSection, keptTicket.FilePath,
			ignoredTicket.TreeSection, ignoredTicket.FilePath))
	}
	return orderedWinners, duplicateWarnings
}

// parseRequestTicket reads and parses a single REQ-*.md file into a ticket.
// Completion-time resolution is left to the caller (it needs the repo root and
// the git lookup).
func parseRequestTicket(filePath string, treeSection string) (*RequestTicket, error) {
	contentBytes, readError := os.ReadFile(filePath)
	if readError != nil {
		return nil, readError
	}

	yamlText, bodyText, hasFrontmatter := splitFrontmatter(string(contentBytes))
	fields := map[string]any{}
	if hasFrontmatter {
		parsedFields, parseError := parseFrontmatterFields(yamlText)
		if parseError == nil {
			fields = parsedFields
		}
	}

	requestId := coerceScalarToString(fields["id"])
	if requestId == "" {
		requestId = deriveRequestIdFromFilename(filePath)
	}
	originalStatus := coerceScalarToString(fields["status"])
	commitHashValue, commitHashField := resolveCommitHash(fields)
	originalTestingStatus := coerceScalarToString(fields["testing_status"])
	normalizedTestingStatus := normalizeTestingStatus(originalTestingStatus)
	testingStatusUnrecognized := false
	if originalTestingStatus != "" && !isKnownTestingStatus(normalizedTestingStatus) {
		// Outside the canonical vocabulary: render as not-yet-tested with an
		// invalid flag (collectTestingWarnings raises the matching data warning).
		normalizedTestingStatus = ""
		testingStatusUnrecognized = true
	}

	ticket := &RequestTicket{
		RequestId:                 requestId,
		Title:                     coerceScalarToString(fields["title"]),
		Status:                    normalizeStatus(originalStatus),
		OriginalStatus:            originalStatus,
		CreatedAt:                 coerceScalarToString(fields["created_at"]),
		ClaimedAt:                 coerceScalarToString(fields["claimed_at"]),
		CompletedAt:               coerceScalarToString(fields["completed_at"]),
		ReservedFor:               coerceScalarToString(fields["reserved_for"]),
		ReservedAt:                coerceScalarToString(fields["reserved_at"]),
		CommitHash:                commitHashValue,
		CommitHashField:           commitHashField,
		UserRequestId:             coerceScalarToString(fields["user_request"]),
		Domain:                    coerceScalarToString(fields["domain"]),
		TestingStatus:             normalizedTestingStatus,
		OriginalTestingStatus:     originalTestingStatus,
		TestingStatusUnrecognized: testingStatusUnrecognized,
		TestedBy:                  coerceScalarToString(fields["tested_by"]),
		TestingUpdatedAt:          coerceScalarToString(fields["testing_updated_at"]),
		TestingFeedback:           coerceScalarToString(fields["testing_feedback"]),
		DependsOn:                 resolveDependsOn(fields),
		BlockedBy:                 coerceToStringList(fields["blocked_by"]),
		BlockedAt:                 coerceScalarToString(fields["blocked_at"]),
		BlockedCheck:              coerceScalarToString(fields["blocked_check"]),
		Related:                   coerceToStringList(fields["related"]),
		Route:                     coerceScalarToString(fields["route"]),
		Batch:                     coerceScalarToString(fields["batch"]),
		BodyMarkdown:              bodyText,
		FilePath:                  filePath,
		TreeSection:               treeSection,
	}
	return ticket, nil
}

// parseUserRequestTicket reads and parses a single UR input.md into a ticket.
func parseUserRequestTicket(filePath string) (*UserRequestTicket, error) {
	contentBytes, readError := os.ReadFile(filePath)
	if readError != nil {
		return nil, readError
	}

	yamlText, bodyText, hasFrontmatter := splitFrontmatter(string(contentBytes))
	fields := map[string]any{}
	if hasFrontmatter {
		parsedFields, parseError := parseFrontmatterFields(yamlText)
		if parseError == nil {
			fields = parsedFields
		}
	}

	userRequestId := coerceScalarToString(fields["id"])
	if userRequestId == "" {
		userRequestId = deriveUserRequestIdFromPath(filePath)
	}

	return &UserRequestTicket{
		UserRequestId:    userRequestId,
		Title:            coerceScalarToString(fields["title"]),
		CreatedAt:        coerceScalarToString(fields["created_at"]),
		BodyMarkdown:     bodyText,
		FilePath:         filePath,
		InputFilePresent: true,
	}, nil
}

// linkRequestsToUserRequests groups every REQ under the UR it points at,
// synthesizing a minimal UR node (InputFilePresent=false) when a REQ references
// a UR whose input.md was not found.
func linkRequestsToUserRequests(board *Board) {
	for _, ticket := range board.AllRequests {
		if ticket.UserRequestId == "" {
			continue
		}
		userRequestTicket := board.UserRequestsById[ticket.UserRequestId]
		if userRequestTicket == nil {
			userRequestTicket = &UserRequestTicket{
				UserRequestId:    ticket.UserRequestId,
				InputFilePresent: false,
			}
			board.UserRequestsById[ticket.UserRequestId] = userRequestTicket
			board.UserRequests = append(board.UserRequests, userRequestTicket)
		}
		userRequestTicket.RequestIds = append(userRequestTicket.RequestIds, ticket.RequestId)
	}
}

// normalizeStatus collapses the legacy completion synonyms to "completed" and
// the won't-do synonyms to "cancelled" (both alias maps mirror the Schema Read
// Contract in actions/work-reference.md), lower-casing/trimming the rest.
// "completed-with-issues" is intentionally left as is — it is the other exact
// terminal-success state isCompletedStatus accepts.
func normalizeStatus(rawStatus string) string {
	normalized := strings.ToLower(strings.TrimSpace(rawStatus))
	switch normalized {
	case "complete", "done", "finished", "closed":
		return "completed"
	case "canceled", "abandoned", "wont-do", "wontfix":
		return "cancelled"
	default:
		return normalized
	}
}

// isCompletedStatus reports whether a normalized status is a terminal-success
// state — exactly "completed" or "completed-with-issues", the Terminal-success
// status set from actions/work-reference.md's Schema Read Contract. The exact
// match is deliberate: a prefix match would let a typo like
// "completed-wth-issues" bypass bucketColumns' unrecognized-status warning and
// silently enter the calendar / Recently done column.
func isCompletedStatus(normalizedStatus string) bool {
	return normalizedStatus == "completed" || normalizedStatus == "completed-with-issues"
}

// isCancelledStatus reports whether a normalized status is the terminal
// won't-do state written by the abandon action (do-work abandon). Cancelled is
// terminal but NOT successful — it shares the Recently-done column and the
// calendar with completed work, while success-only readers keep excluding it
// via isCompletedStatus.
func isCancelledStatus(normalizedStatus string) bool {
	return normalizedStatus == "cancelled"
}

// isTerminalResolvedStatus reports whether a normalized status is terminally
// resolved — the Terminal-resolved status set from actions/work-reference.md's
// Schema Read Contract: the terminal-success pair plus "cancelled". This set
// gates completion-time resolution, Recently-done bucketing, and the calendar.
func isTerminalResolvedStatus(normalizedStatus string) bool {
	return isCompletedStatus(normalizedStatus) || isCancelledStatus(normalizedStatus)
}

// isNeedsInputOrBlockedStatus reports whether a normalized status belongs in the
// Needs-input / Blocked column.
func isNeedsInputOrBlockedStatus(normalizedStatus string) bool {
	switch normalizedStatus {
	case "pending-answers",
		"blocked",
		"blocked-archive-collision",
		"blocked-dependency-cycle",
		"failed":
		return true
	default:
		return false
	}
}

// resolveCommitHash returns the first non-empty commit hash among the canonical
// field and its accepted variants, in priority order, plus the frontmatter key
// it came from (so an anomaly report can name the exact broken field).
func resolveCommitHash(fields map[string]any) (string, string) {
	for _, key := range []string{"commit", "commit_hash", "green_commit", "commit_green", "impl_commit"} {
		if value := coerceScalarToString(fields[key]); value != "" {
			return value, key
		}
	}
	return "", ""
}

// resolveDependsOn returns the canonical dependency list: depends_on when
// present, otherwise the legacy dependencies field.
func resolveDependsOn(fields map[string]any) []string {
	dependsOn := coerceToStringList(fields["depends_on"])
	if len(dependsOn) == 0 {
		dependsOn = coerceToStringList(fields["dependencies"])
	}
	return dependsOn
}

// resolveCompletionTime applies the completion fallback chain: frontmatter
// completed_at → the commit hash's git committer date → unresolved. The git
// step is best-effort (a nil or failing lookup leaves the completion undated).
// File mtime is deliberately NOT a fallback — a fresh clone, branch checkout,
// or tarball extraction resets every mtime to "now", which stamped months-old
// REQs as completed today; undated completions land in the calendar's
// "undated" bucket instead.
func resolveCompletionTime(ticket *RequestTicket, repoRoot string, gitLookup gitCommitDateLookup) (time.Time, CompletionTimeSource) {
	if ticket.CompletedAt != "" {
		if parsed, ok := parseTimestamp(ticket.CompletedAt); ok {
			return parsed, CompletionFromFrontmatter
		}
	}
	if ticket.CommitHash != "" && gitLookup != nil {
		if committedAt, ok := gitLookup(repoRoot, ticket.CommitHash); ok {
			return committedAt, CompletionFromGitLog
		}
	}
	return time.Time{}, CompletionUnresolved
}

// detectCompletionAnomaly inspects a terminal-resolved ticket AFTER its
// completion instant was resolved and reports whether its completion
// bookkeeping is broken, with a reason naming the broken field(s):
//
//   - a completed_at value that exists but fails parseTimestamp — flagged even
//     when a commit hash rescued the date, because the field is still wrong on
//     disk;
//   - a commit-hash field that was consulted (completed_at absent or
//     unparseable) but that git could not resolve;
//   - neither field present at all — a done-flip that forgot to stamp.
//
// A ticket whose completed_at parsed is never anomalous: the commit hash is not
// consulted on that path, so it is not re-validated here — doing so would cost
// one git subprocess per archived ticket for no board-behavior change.
func detectCompletionAnomaly(ticket *RequestTicket) (bool, string) {
	if !isTerminalResolvedStatus(ticket.Status) || ticket.CompletionTimeSource == CompletionFromFrontmatter {
		return false, ""
	}
	var brokenFieldReasons []string
	if ticket.CompletedAt != "" {
		// Non-empty but the source is not frontmatter ⇒ parseTimestamp rejected it.
		brokenFieldReasons = append(brokenFieldReasons, fmt.Sprintf(
			"completed_at %q does not parse as a timestamp", ticket.CompletedAt))
	}
	if ticket.CommitHash != "" && ticket.CompletionTimeSource != CompletionFromGitLog {
		commitFieldName := ticket.CommitHashField
		if commitFieldName == "" {
			commitFieldName = "commit"
		}
		// The lookup is best-effort and cannot distinguish its failure modes
		// here, so the reason must not blame the hash alone: "git could not
		// date it" covers an unknown hash AND a missing git binary / non-repo
		// tree (lookupGitCommitDate logs the missing-binary case once).
		brokenFieldReasons = append(brokenFieldReasons, fmt.Sprintf(
			"%s %q could not be dated — the hash is unknown to git, or git/the repository is unavailable", commitFieldName, ticket.CommitHash))
	}
	if len(brokenFieldReasons) > 0 {
		return true, strings.Join(brokenFieldReasons, "; ")
	}
	if ticket.CompletionTimeSource == CompletionUnresolved {
		return true, "terminal status but no completed_at and no resolvable commit hash"
	}
	return false, "" // dated via git with no broken sibling field
}

// lookupGitCommitDate resolves a commit hash to its committer date via
// `git -C <repoRoot> log -1 --format=%cI <hash>`. It is best-effort: a missing
// git binary, an unknown hash, or an unparseable date all return (zero, false).
// The hash comes from untrusted REQ frontmatter, so it is validated as plain
// hex before being placed in argv — an option-shaped value like "--all" or
// "--output=<path>" would otherwise be parsed by git as a flag (argument
// injection).
func lookupGitCommitDate(repoRoot string, commitHash string) (time.Time, bool) {
	trimmedHash := strings.TrimSpace(commitHash)
	if !isPlausibleCommitHash(trimmedHash) {
		return time.Time{}, false
	}
	if !gitBinaryAvailable() {
		return time.Time{}, false
	}
	command := exec.Command("git", "-C", repoRoot, "log", "-1", "--format=%cI", trimmedHash)
	output, runError := command.Output()
	if runError != nil {
		return time.Time{}, false
	}
	text := strings.TrimSpace(string(output))
	if text == "" {
		return time.Time{}, false
	}
	parsed, parseError := time.Parse(time.RFC3339, text)
	if parseError != nil {
		return time.Time{}, false
	}
	return parsed, true
}

// gitBinaryProbe caches a one-time PATH lookup for git, logging once when the
// binary is missing — so an operator running without git gets one clear line
// naming the real cause instead of a per-ticket anomaly that blames the hash,
// and the board skips one doomed subprocess per archived ticket.
var gitBinaryProbe struct {
	probeOnce sync.Once
	available bool
}

func gitBinaryAvailable() bool {
	gitBinaryProbe.probeOnce.Do(func() {
		_, lookupError := exec.LookPath("git")
		gitBinaryProbe.available = lookupError == nil
		if !gitBinaryProbe.available {
			log.Printf("queue-kanban: git binary not found on PATH — commit-hash completion dating is disabled")
		}
	})
	return gitBinaryProbe.available
}

// isPlausibleCommitHash reports whether text looks like an abbreviated or full
// git object hash: 4–64 hex digits (64 covers sha256-object repos). Anything
// else is rejected before it can reach a git argv.
func isPlausibleCommitHash(text string) bool {
	if len(text) < 4 || len(text) > 64 {
		return false
	}
	for _, character := range text {
		switch {
		case character >= '0' && character <= '9':
		case character >= 'a' && character <= 'f':
		case character >= 'A' && character <= 'F':
		default:
			return false
		}
	}
	return true
}

// detectFutureTimestampFields checks every timestamp the frontmatter can carry
// against the board's `now` and returns a "<field> <raw value>" entry for each
// one that parses to later than now + futureTimestampSkewAllowance. A future
// stamp is bookkeeping damage worth surfacing wherever the ticket renders: the
// usual cause is local wall-clock time written with a `Z` suffix, which makes
// elapsed-time math (queue wait, claim stopwatch) silently wrong until the wall
// clock catches up. Unparseable and absent values are not this check's concern
// — other paths (completion anomalies, reservation staleness) own those.
func detectFutureTimestampFields(ticket *RequestTicket, now time.Time) []string {
	timestampFields := []struct {
		fieldName string
		rawValue  string
	}{
		{"created_at", ticket.CreatedAt},
		{"claimed_at", ticket.ClaimedAt},
		{"completed_at", ticket.CompletedAt},
		{"blocked_at", ticket.BlockedAt},
		{"reserved_at", ticket.ReservedAt},
		{"testing_updated_at", ticket.TestingUpdatedAt},
	}
	skewHorizon := now.Add(futureTimestampSkewAllowance)
	var futureFieldEntries []string
	for _, timestampField := range timestampFields {
		parsedInstant, parsedOk := parseTimestamp(timestampField.rawValue)
		if parsedOk && parsedInstant.After(skewHorizon) {
			futureFieldEntries = append(futureFieldEntries,
				timestampField.fieldName+" "+strings.TrimSpace(timestampField.rawValue))
		}
	}
	return futureFieldEntries
}

// parseTimestamp parses the timestamp shapes seen across REQ frontmatter:
// RFC3339 with a Z or numeric offset, an offset-less datetime, and a bare date.
func parseTimestamp(text string) (time.Time, bool) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	} {
		if parsed, parseError := time.Parse(layout, trimmed); parseError == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

// bucketColumns sorts every ticket into the active-work columns by normalized
// status. Pending additionally splits on dependency readiness (annotated by
// annotateDependencyState, which must have run first): a pending ticket with no
// unmet dependency is what the work loop would actually claim next. The split is
// a view, not a status change — the gating is dynamic, so a waiting ticket stays
// `pending` on disk and becomes ready the moment its upstream completes.
// Terminally resolved tickets (completed*/cancelled) only enter
// RecentlyDone when their completion instant falls inside the window; older
// resolutions are left for the calendar. A status outside the known vocabulary
// is never silently dropped (Schema Read Contract, actions/work-reference.md):
// the ticket lands in Needs-input/Blocked so it stays visible, plus a warning
// naming the unrecognized status. Completion anomalies get the same
// never-silent guarantee: a terminal ticket flagged CompletionAnomaly enters
// the always-visible CompletionAnomalies column (window-independent — an
// unresolved completion would otherwise vanish from the board entirely) plus a
// warning carrying the reason and the concrete fix.
func bucketColumns(tickets []*RequestTicket, now time.Time, recentWindow time.Duration) (BoardColumns, []string) {
	var columns BoardColumns
	var statusWarnings []string
	for _, ticket := range tickets {
		switch {
		case ticket.Status == "pending":
			columns.Pending = append(columns.Pending, ticket)
			if len(ticket.UnmetDependencies) == 0 {
				columns.PendingReady = append(columns.PendingReady, ticket)
			} else {
				columns.PendingWaiting = append(columns.PendingWaiting, ticket)
			}
		case ticket.Status == "claimed":
			columns.Claimed = append(columns.Claimed, ticket)
		case ticket.Status == "reserved":
			// Allocated to a DIFFERENT worktree/cloud session (do-work reserve,
			// actions/reserve.md). Someone owns it, so it shares the Claimed
			// column — but the frontend grays it out because that someone is not
			// this board's session.
			columns.Claimed = append(columns.Claimed, ticket)
			if staleWarning := annotateReservationStaleness(ticket, now); staleWarning != "" {
				statusWarnings = append(statusWarnings, staleWarning)
			}
		case isNeedsInputOrBlockedStatus(ticket.Status):
			columns.NeedsInputOrBlocked = append(columns.NeedsInputOrBlocked, ticket)
		case isTerminalResolvedStatus(ticket.Status):
			if ticket.CompletionAnomaly {
				columns.CompletionAnomalies = append(columns.CompletionAnomalies, ticket)
				statusWarnings = append(statusWarnings, fmt.Sprintf(
					"%s (status %s) has a completion anomaly: %s — shown under Completion anomalies; fix: stamp completed_at: with a UTC ISO instant and/or a commit: field with a valid implementation commit hash in its frontmatter",
					ticket.RequestId, ticket.Status, ticket.CompletionAnomalyReason))
			}
			if isWithinRecentWindow(ticket.CompletionTime, now, recentWindow) {
				columns.RecentlyDone = append(columns.RecentlyDone, ticket)
			}
		default:
			ticket.StatusUnrecognized = true
			columns.NeedsInputOrBlocked = append(columns.NeedsInputOrBlocked, ticket)
			statusWarnings = append(statusWarnings, fmt.Sprintf(
				"%s has unrecognized status %q — shown under Needs input / Blocked; fix: edit its status: to a Schema Read Contract value (actions/work-reference.md) or run do-work forensics",
				ticket.RequestId, ticket.OriginalStatus))
		}
	}
	sort.SliceStable(columns.RecentlyDone, func(i, j int) bool {
		return columns.RecentlyDone[i].CompletionTime.After(columns.RecentlyDone[j].CompletionTime)
	})
	return columns, statusWarnings
}

// annotateReservationStaleness marks a reserved ticket stale when its
// reserved_at is missing, unparseable, or more than reservationStaleAfter
// before now, and returns a recategorize-suggestion warning ("" when the
// reservation is fresh). The suggestion mirrors actions/work.md Step 1's
// stale-reservation check: the owning session may be dead, but the decision
// stays with the user — the board never auto-releases.
func annotateReservationStaleness(ticket *RequestTicket, now time.Time) string {
	reservedInstant, parsedOk := parseTimestamp(ticket.ReservedAt)
	if !parsedOk {
		ticket.ReservationStale = true
		return fmt.Sprintf(
			"%s is reserved (for %q) but has no parseable reserved_at — treated as stale; recategorize: do-work release %s, do-work run %s, or leave it if that session is still active",
			ticket.RequestId, ticket.ReservedFor, ticket.RequestId, ticket.RequestId)
	}
	if now.Sub(reservedInstant) <= reservationStaleAfter {
		return ""
	}
	ticket.ReservationStale = true
	return fmt.Sprintf(
		"%s has been reserved for %q for more than 24h — recategorize: do-work release %s (back to queue), do-work run %s (claim here), or leave it if that session is still active",
		ticket.RequestId, ticket.ReservedFor, ticket.RequestId, ticket.RequestId)
}

// isWithinRecentWindow reports whether a completion instant is non-zero and falls
// within recentWindow before now (future instants are included to tolerate clock
// skew between the recording machine and the board host).
func isWithinRecentWindow(completionTime time.Time, now time.Time, recentWindow time.Duration) bool {
	if completionTime.IsZero() {
		return false
	}
	return completionTime.After(now.Add(-recentWindow))
}

// buildDependencyGraph builds the depends_on edge list and forward/reverse
// adjacency maps. An edge is Resolved when its target id is a parsed ticket.
func buildDependencyGraph(tickets []*RequestTicket, requestsById map[string]*RequestTicket) DependencyGraph {
	graph := DependencyGraph{
		DependsOn:  map[string][]string{},
		Dependents: map[string][]string{},
	}
	for _, ticket := range tickets {
		for _, dependencyId := range ticket.DependsOn {
			_, resolved := requestsById[dependencyId]
			graph.Edges = append(graph.Edges, DependencyEdge{
				FromRequestId: ticket.RequestId,
				ToRequestId:   dependencyId,
				Resolved:      resolved,
			})
			graph.DependsOn[ticket.RequestId] = append(graph.DependsOn[ticket.RequestId], dependencyId)
			graph.Dependents[dependencyId] = append(graph.Dependents[dependencyId], ticket.RequestId)
		}
	}
	return graph
}

// annotateDependencyState walks every depends_on edge once and fills in the two
// derived views the board renders from: each ticket's UnmetDependencies (the
// forward edge — what still has to land) and each ticket's Dependents (the
// reverse edge — what unblocks when this lands). It returns one warning per
// dangling dependency.
//
// A dependency is MET only when its target reached terminal SUCCESS — exactly
// the `completed` / `completed-with-issues` pair (actions/work-reference.md's
// Schema Read Contract, which the work loop's Step 1 selection scan gates on).
// Two consequences follow from that contract and are deliberate here:
//
//   - `cancelled` does NOT satisfy a dependency. The dependent presumably needed
//     the cancelled REQ's output, so it stays waiting until the user re-points
//     depends_on or abandons it too — it must never quietly read as ready.
//   - A dangling id (no such REQ anywhere in the tree) counts as UNMET, never as
//     satisfied. Failing open would silently promote a REQ into Ready on the
//     strength of a typo'd dependency. It is also surfaced as a warning, because
//     a pointer to nothing can never self-resolve — no amount of work clears it.
func annotateDependencyState(board *Board) []string {
	var danglingDependencyWarnings []string

	for _, ticket := range board.AllRequests {
		alreadySeenDependencyIds := map[string]bool{}
		for _, dependencyId := range ticket.DependsOn {
			if alreadySeenDependencyIds[dependencyId] {
				continue // a repeated depends_on entry must not double-count as a dependent
			}
			alreadySeenDependencyIds[dependencyId] = true

			dependencyTicket, dependencyExists := board.RequestsById[dependencyId]
			if !dependencyExists {
				ticket.UnmetDependencies = append(ticket.UnmetDependencies, dependencyId)
				danglingDependencyWarnings = append(danglingDependencyWarnings, fmt.Sprintf(
					"%s depends on %s, which is not in the do-work tree — treated as unmet",
					ticket.RequestId, dependencyId))
				continue
			}
			dependencyTicket.Dependents = append(dependencyTicket.Dependents, ticket.RequestId)
			if !isCompletedStatus(dependencyTicket.Status) {
				ticket.UnmetDependencies = append(ticket.UnmetDependencies, dependencyId)
			}
		}
	}

	for _, ticket := range board.AllRequests {
		sortRequestIdList(ticket.Dependents)
	}
	return danglingDependencyWarnings
}

// buildCalendar produces a completion-time-keyed index over every terminally
// resolved (completed*/cancelled) ticket, sorted most-recent-first. Tickets
// whose completion instant could not be resolved are kept — never silently
// dropped — under the trailing "undated" day bucket (the zero CompletionTime
// sorts them after every dated entry, and board.js falls back to rendering the
// raw day key as the group label).
func buildCalendar(tickets []*RequestTicket) []CalendarEntry {
	var entries []CalendarEntry
	for _, ticket := range tickets {
		if !isTerminalResolvedStatus(ticket.Status) {
			continue
		}
		if ticket.CompletionTime.IsZero() {
			entries = append(entries, CalendarEntry{
				RequestId:  ticket.RequestId,
				TimeSource: ticket.CompletionTimeSource,
				DayKey:     undatedCalendarDayKey,
			})
			continue
		}
		entries = append(entries, CalendarEntry{
			RequestId:      ticket.RequestId,
			CompletionTime: ticket.CompletionTime,
			TimeSource:     ticket.CompletionTimeSource,
			DayKey:         ticket.CompletionTime.UTC().Format("2006-01-02"),
		})
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].CompletionTime.After(entries[j].CompletionTime)
	})
	return entries
}

// deriveRequestIdFromFilename recovers the canonical REQ id from a filename like
// "REQ-1203-modal-shell.md" → "REQ-1203" when frontmatter has no id field.
func deriveRequestIdFromFilename(filePath string) string {
	baseName := strings.TrimSuffix(filepath.Base(filePath), ".md")
	parts := strings.Split(baseName, "-")
	if len(parts) >= 2 && strings.EqualFold(parts[0], "REQ") {
		return "REQ-" + parts[1]
	}
	return baseName
}

// deriveUserRequestIdFromPath recovers a UR id from the parent directory name of
// an input.md path (e.g. .../user-requests/UR-448/input.md → "UR-448").
func deriveUserRequestIdFromPath(filePath string) string {
	return filepath.Base(filepath.Dir(filePath))
}

// sortRequestTickets orders tickets by numeric REQ id (REQ-9 before REQ-100).
func sortRequestTickets(tickets []*RequestTicket) {
	sort.SliceStable(tickets, func(i, j int) bool {
		return requestIdLess(tickets[i].RequestId, tickets[j].RequestId)
	})
}

// sortUserRequestTickets orders URs by numeric UR id.
func sortUserRequestTickets(tickets []*UserRequestTicket) {
	sort.SliceStable(tickets, func(i, j int) bool {
		return identifierLess(tickets[i].UserRequestId, tickets[j].UserRequestId)
	})
}

// sortRequestIdList orders a slice of REQ ids in place by numeric id.
func sortRequestIdList(ids []string) {
	sort.SliceStable(ids, func(i, j int) bool {
		return requestIdLess(ids[i], ids[j])
	})
}

// requestIdLess compares two "REQ-NNNN" ids numerically, falling back to a plain
// string comparison when either lacks a numeric suffix.
func requestIdLess(left string, right string) bool {
	return identifierLess(left, right)
}

// identifierLess compares two "PREFIX-NNNN" identifiers by their numeric suffix
// when both expose one, otherwise lexically. It works for both REQ- and UR- ids.
func identifierLess(left string, right string) bool {
	leftNumber, leftOk := numericIdSuffix(left)
	rightNumber, rightOk := numericIdSuffix(right)
	if leftOk && rightOk && leftNumber != rightNumber {
		return leftNumber < rightNumber
	}
	if leftOk != rightOk {
		return leftOk
	}
	return left < right
}

// numericIdSuffix extracts the integer following the last hyphen of an id such
// as "REQ-1207" → 1207. It reports false when there is no parseable suffix.
func numericIdSuffix(identifier string) (int, bool) {
	hyphenIndex := strings.LastIndex(identifier, "-")
	if hyphenIndex < 0 || hyphenIndex == len(identifier)-1 {
		return 0, false
	}
	number, parseError := strconv.Atoi(identifier[hyphenIndex+1:])
	if parseError != nil {
		return 0, false
	}
	return number, true
}

// coerceScalarToString turns a YAML scalar of any concrete type into a trimmed
// string, so frontmatter fields that are sometimes quoted, sometimes bare, and
// occasionally numeric all read back as strings. yaml.v3 decodes ISO-8601
// timestamp values (created_at / completed_at / claimed_at) into time.Time when
// the target is `any`; those are normalized back to RFC3339 (UTC) so the raw
// field text stays a real timestamp that parseTimestamp can re-read — not Go's
// default "2006-01-02 15:04:05 +0000 UTC" rendering.
func coerceScalarToString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case time.Time:
		return typed.UTC().Format(time.RFC3339)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
}

// coerceToStringList normalizes a YAML value into a string slice: a sequence maps
// element-wise, a bare scalar becomes a one-element slice, and nil/empty yields
// nil. Empty elements are dropped.
func coerceToStringList(value any) []string {
	switch typed := value.(type) {
	case nil:
		return nil
	case []any:
		list := make([]string, 0, len(typed))
		for _, element := range typed {
			text := coerceScalarToString(element)
			if text != "" {
				list = append(list, text)
			}
		}
		if len(list) == 0 {
			return nil
		}
		return list
	default:
		text := coerceScalarToString(typed)
		if text == "" {
			return nil
		}
		return []string{text}
	}
}
