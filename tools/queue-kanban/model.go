package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

// RequestTicket is one parsed REQ-*.md file: its frontmatter fields (with status
// normalized and commit-hash variants collapsed), the raw Markdown body kept for
// later HTML rendering, where it was found in the tree, and how its completion
// instant was resolved.
type RequestTicket struct {
	RequestId      string // canonical "REQ-NNNN" (frontmatter id, else derived from the filename)
	Title          string
	Status         string // normalized status (complete/done/finished/closed → completed)
	OriginalStatus string // verbatim frontmatter status before normalization

	CreatedAt   string // raw frontmatter timestamp text, "" when absent
	ClaimedAt   string
	CompletedAt string // raw frontmatter completed_at text, "" when absent

	CommitHash    string // resolved from commit / commit_hash / green_commit / commit_green / impl_commit
	UserRequestId string // "UR-NNN" upward pointer (the reliable REQ→UR link), "" when absent
	Domain        string

	DependsOn []string // canonical dependency REQ ids (depends_on wins over legacy dependencies)
	BlockedBy []string // legacy blocked_by ids, kept distinct from DependsOn
	Related   []string // soft relations (not dependency edges)

	Route    string
	Batch    string
	Severity string

	BodyMarkdown string // raw Markdown body after the closing frontmatter fence

	FilePath    string // absolute path on disk
	TreeSection string // "queue" | "working" | "archive"

	CompletionTime       time.Time            // resolved completion instant (zero when unresolved)
	CompletionTimeSource CompletionTimeSource // how CompletionTime was resolved
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

// CalendarEntry places one completed REQ on the completion timeline, recording
// the resolved instant, how it was resolved, and a UTC day bucket key.
type CalendarEntry struct {
	RequestId      string
	CompletionTime time.Time
	TimeSource     CompletionTimeSource
	DayKey         string // "2006-01-02" UTC bucket
}

// BoardColumns holds the active-work buckets. Completed REQs older than the
// recent window are NOT represented here — they live in Board.Calendar.
type BoardColumns struct {
	Pending             []*RequestTicket // status pending
	Claimed             []*RequestTicket // status claimed
	NeedsInputOrBlocked []*RequestTicket // pending-answers / blocked-* / failed
	RecentlyDone        []*RequestTicket // completed* whose completion instant is within the window
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
	Calendar        []CalendarEntry // completed REQs sorted most-recent-first (undated entries last)

	Warnings []string // data-shape warnings (duplicate ids, unrecognized statuses) — surfaced, never silently dropped
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
		if isCompletedStatus(ticket.Status) {
			completionTime, completionSource := resolveCompletionTime(ticket, repoRoot, gitLookup)
			ticket.CompletionTime = completionTime
			ticket.CompletionTimeSource = completionSource
		} else {
			ticket.CompletionTimeSource = CompletionUnresolved
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

	columns, columnWarnings := bucketColumns(board.AllRequests, now, recentWindow)
	board.Columns = columns
	board.Warnings = append(board.Warnings, columnWarnings...)
	board.DependencyGraph = buildDependencyGraph(board.AllRequests, board.RequestsById)
	board.Calendar = buildCalendar(board.AllRequests)

	return board, nil
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

	ticket := &RequestTicket{
		RequestId:      requestId,
		Title:          coerceScalarToString(fields["title"]),
		Status:         normalizeStatus(originalStatus),
		OriginalStatus: originalStatus,
		CreatedAt:      coerceScalarToString(fields["created_at"]),
		ClaimedAt:      coerceScalarToString(fields["claimed_at"]),
		CompletedAt:    coerceScalarToString(fields["completed_at"]),
		CommitHash:     resolveCommitHash(fields),
		UserRequestId:  coerceScalarToString(fields["user_request"]),
		Domain:         coerceScalarToString(fields["domain"]),
		DependsOn:      resolveDependsOn(fields),
		BlockedBy:      coerceToStringList(fields["blocked_by"]),
		Related:        coerceToStringList(fields["related"]),
		Route:          coerceScalarToString(fields["route"]),
		Batch:          coerceScalarToString(fields["batch"]),
		Severity:       coerceScalarToString(fields["severity"]),
		BodyMarkdown:   bodyText,
		FilePath:       filePath,
		TreeSection:    treeSection,
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
// lower-cases/trims the rest. "completed-with-issues" is intentionally left as
// is — it is already a completed* state recognized by isCompletedStatus.
func normalizeStatus(rawStatus string) string {
	normalized := strings.ToLower(strings.TrimSpace(rawStatus))
	switch normalized {
	case "complete", "done", "finished", "closed":
		return "completed"
	default:
		return normalized
	}
}

// isCompletedStatus reports whether a normalized status is any completed* state
// (covers "completed" and "completed-with-issues").
func isCompletedStatus(normalizedStatus string) bool {
	return strings.HasPrefix(normalizedStatus, "completed")
}

// isNeedsInputOrBlockedStatus reports whether a normalized status belongs in the
// Needs-input / Blocked column.
func isNeedsInputOrBlockedStatus(normalizedStatus string) bool {
	switch normalizedStatus {
	case "pending-answers",
		"blocked-archive-collision",
		"blocked-dependency-cycle",
		"failed":
		return true
	default:
		return false
	}
}

// resolveCommitHash returns the first non-empty commit hash among the canonical
// field and its accepted variants, in priority order.
func resolveCommitHash(fields map[string]any) string {
	for _, key := range []string{"commit", "commit_hash", "green_commit", "commit_green", "impl_commit"} {
		if value := coerceScalarToString(fields[key]); value != "" {
			return value
		}
	}
	return ""
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
// status. Completed* tickets only enter RecentlyDone when their completion
// instant falls inside the window; older completions are left for the calendar.
// A status outside the known vocabulary is never silently dropped (Schema Read
// Contract, actions/work-reference.md): the ticket lands in Needs-input/Blocked
// so it stays visible, plus a warning naming the unrecognized status.
func bucketColumns(tickets []*RequestTicket, now time.Time, recentWindow time.Duration) (BoardColumns, []string) {
	var columns BoardColumns
	var statusWarnings []string
	for _, ticket := range tickets {
		switch {
		case ticket.Status == "pending":
			columns.Pending = append(columns.Pending, ticket)
		case ticket.Status == "claimed":
			columns.Claimed = append(columns.Claimed, ticket)
		case isNeedsInputOrBlockedStatus(ticket.Status):
			columns.NeedsInputOrBlocked = append(columns.NeedsInputOrBlocked, ticket)
		case isCompletedStatus(ticket.Status):
			if isWithinRecentWindow(ticket.CompletionTime, now, recentWindow) {
				columns.RecentlyDone = append(columns.RecentlyDone, ticket)
			}
		default:
			columns.NeedsInputOrBlocked = append(columns.NeedsInputOrBlocked, ticket)
			statusWarnings = append(statusWarnings, fmt.Sprintf(
				"%s has unrecognized status %q — shown under Needs input / Blocked",
				ticket.RequestId, ticket.OriginalStatus))
		}
	}
	sort.SliceStable(columns.RecentlyDone, func(i, j int) bool {
		return columns.RecentlyDone[i].CompletionTime.After(columns.RecentlyDone[j].CompletionTime)
	})
	return columns, statusWarnings
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

// buildCalendar produces a completion-time-keyed index over every completed*
// ticket, sorted most-recent-first. Tickets whose completion instant could not
// be resolved are kept — never silently dropped — under the trailing "undated"
// day bucket (the zero CompletionTime sorts them after every dated entry, and
// board.js falls back to rendering the raw day key as the group label).
func buildCalendar(tickets []*RequestTicket) []CalendarEntry {
	var entries []CalendarEntry
	for _, ticket := range tickets {
		if !isCompletedStatus(ticket.Status) {
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
