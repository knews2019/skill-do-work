package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// embeddedWebAssets holds the hand-authored static frontend (HTML shell + CSS +
// JS fragments) that `generate` inlines into a single self-contained index.html. REQ-1209
// (`serve`) embeds and re-serves the SAME web/ directory unchanged, so the shape
// here is the shared contract.
//
//go:embed web/template.html web/board.css web/*.js
var embeddedWebAssets embed.FS

// Inline placeholder tokens in web/template.html. They are deliberately
// comment-shaped so the template stays valid HTML/CSS/JS on its own (it can be
// opened during development without the inlining step). generate replaces each
// exactly once.
const (
	inlineStylePlaceholder         = "/* INLINE_BOARD_STYLES */"
	inlineScriptPlaceholder        = "/* INLINE_BOARD_SCRIPT */"
	boardJavaScriptPlaceholder     = "/* INLINE_BOARD_FRAGMENTS */"
	boardJavaScriptPlaceholderLine = boardJavaScriptPlaceholder + "\n"
	boardJavaScriptShellPath       = "web/board.js"
	generatedAtDisplayPlaceholder  = "GENERATED_AT_DISPLAY"
	projectNamePlaceholder         = "PROJECT_NAME"
)

// boardJavaScriptFragmentPaths is the sole execution-order manifest for the
// private classic-script client. The wildcard embed makes authored inventory
// testable; it never decides runtime order.
var boardJavaScriptFragmentPaths = [...]string{
	"web/board-core.js",
	"web/board-filters.js",
	"web/board-cards.js",
	"web/board-calendar.js",
	"web/board-durations.js",
	"web/board-timeline.js",
	"web/board-testing.js",
	"web/board-detail.js",
	"web/board-controls.js",
	"web/board-clipboard.js",
}

// boardDataJsFilename is the sibling file written next to index.html that assigns
// the board JSON to window.queueKanbanBoardData. A <script src="board-data.js">
// in index.html loads it before the assembled client so the board renders offline from file://.
const boardDataJsFilename = "board-data.js"

// boardMarkdownJsFilename carries raw REQ/UR Markdown for the drawer Copy button.
// It is deliberately NOT referenced by index.html: board-clipboard.js injects it
// on the first copy, keeping duplicate raw bodies out of the initial board payload while
// preserving file:// support (a dynamically-added script works without fetch/CORS).
const boardMarkdownJsFilename = "board-markdown.js"

// generatedBoardData is the JSON data island embedded in the static page. It is
// the single source of truth the client-side script renders every view from, so
// the board works with zero network once the file is open.
type generatedBoardData struct {
	GeneratedAt      string                          `json:"generatedAt"`
	Columns          generatedColumns                `json:"columns"`
	RequestOrder     []string                        `json:"requestOrder"`
	Requests         map[string]generatedRequest     `json:"requests"`
	UserRequestOrder []string                        `json:"userRequestOrder"`
	UserRequests     map[string]generatedUserRequest `json:"userRequests"`
	Calendar         []generatedCalendarEntry        `json:"calendar"`
	Durations        generatedDurations              `json:"durations"`
	Timeline         generatedTimeline               `json:"timeline"`
	Notes            []generatedNote                 `json:"notes,omitempty"` // do-work/notes.md lines — rendered as a strip above the queue

	// Verify findings carried into the page so a human looking at the board sees
	// what `queue-kanban verify` sees (REQ-284). Three categories are suppressed
	// before they get here — see attachVerifyFindings. VerifySkipped is never
	// dropped: a skipped probe rendering as nothing reads as "checked and clean".
	VerifyFindings []generatedVerifyFinding `json:"verifyFindings,omitempty"`
	VerifySkipped  []string                 `json:"verifySkipped,omitempty"`

	Warnings []string `json:"warnings,omitempty"` // data-shape warnings (e.g. duplicate ids, unrecognized statuses, future-dated stamps) — rendered as a banner

	TestingProfiles []string `json:"testingProfiles,omitempty"` // do-work/testers.md profiles for the testing view's tester picker
	// True only when served by the live server (serve.go sets it): the testing
	// view's write actions need the /api/testing/* endpoints, so a static
	// snapshot renders the view read-only.
	LiveTestingApi bool `json:"liveTestingApi,omitempty"`
	// True only when served by the live server (serve.go sets it): the drawer
	// turns existing file-path mentions into GET /file?path=… links, which need
	// the server behind them — a static snapshot has no server, so there only
	// the missing-file marker applies.
	LiveFileApi bool `json:"liveFileApi,omitempty"`
	// File-path mentions found in REQ/UR bodies → whether the file exists in
	// the repo (checked at build time by collectRepoFileMentions). The drawer
	// links only paths mapped true and flags paths mapped false as missing.
	RepoFileMentions map[string]bool `json:"repoFileMentions,omitempty"`
}

// generatedColumns lists the active-board buckets as REQ id slices. RecentlyDone
// is the generate-time default-window snapshot; the client recomputes it from the
// calendar for the 24h/48h/7d toggle, so this slice is just the initial paint.
// PendingReady and PendingWaiting partition Pending — the full list is kept so a
// consumer that ignores dependency readiness still sees every pending ticket.
type generatedColumns struct {
	Pending             []string `json:"pending"`
	PendingReady        []string `json:"pendingReady"`
	PendingWaiting      []string `json:"pendingWaiting"`
	Claimed             []string `json:"claimed"`
	NeedsInputOrBlocked []string `json:"needsInputOrBlocked"`
	RecentlyDone        []string `json:"recentlyDone"`
	// Terminal REQs whose completion instant is broken (see RequestTicket
	// .CompletionAnomaly). Rendered as an always-visible anomaly strip,
	// independent of the client-side recently-done window.
	CompletionAnomalies []string `json:"completionAnomalies"`
}

// generatedRequest is one REQ card's full payload, including its pre-rendered
// Markdown body so the detail drawer opens with zero network.
type generatedRequest struct {
	RequestId             string   `json:"id"`
	Title                 string   `json:"title"`
	Status                string   `json:"status"`
	OriginalStatus        string   `json:"originalStatus"`
	StatusUnrecognized    bool     `json:"statusUnrecognized,omitempty"`
	Error                 string   `json:"error,omitempty"`
	ErrorType             string   `json:"errorType,omitempty"`
	OriginalErrorType     string   `json:"originalErrorType,omitempty"`
	ErrorTypeUnrecognized bool     `json:"errorTypeUnrecognized,omitempty"`
	Domain                string   `json:"domain"`
	OriginalDomain        string   `json:"originalDomain,omitempty"`
	DomainUnrecognized    bool     `json:"domainUnrecognized,omitempty"`
	UserRequestId         string   `json:"userRequestId"`
	DependsOn             []string `json:"dependsOn"`
	UnmetDependencies     []string `json:"unmetDependencies"`
	Dependents            []string `json:"dependents"`
	BlockedBy             []string `json:"blockedBy"`
	BlockedAt             string   `json:"blockedAt,omitempty"`
	BlockedCheck          string   `json:"blockedCheck,omitempty"`
	Related               []string `json:"related"`
	WriteSet              []string `json:"writeSet"`
	// The session a pending REQ is earmarked for (see RequestTicket.AssignedTo).
	// Verbatim and display only — a card badge and a drawer row, never column or
	// dispatch meaning.
	AssignedTo string `json:"assignedTo,omitempty"`
	// Other pending/claimed REQ ids whose write_set could touch the same files
	// (see RequestTicket.WriteSetOverlaps). Display only — the card badge and a
	// drawer row; no column or dispatch meaning.
	WriteSetOverlaps  []string `json:"writeSetOverlaps,omitempty"`
	Route             string   `json:"route"`
	OriginalRoute     string   `json:"originalRoute,omitempty"`
	RouteUnrecognized bool     `json:"routeUnrecognized,omitempty"`
	// Whether anyone would ever notice the work (see RequestTicket.Impact).
	// Display only — a card chip rendered for every value except the
	// impact-user-visible default, plus a drawer row; never column or scheduling
	// meaning. Same raw provenance as domain/route so the chip can say what was
	// declared, not just what it normalized to.
	Impact             string `json:"impact,omitempty"`
	OriginalImpact     string `json:"originalImpact,omitempty"`
	ImpactUnrecognized bool   `json:"impactUnrecognized,omitempty"`
	// Triage bit separating small mechanical fixes from real work (see
	// RequestTicket.EffortEstimate). Display only — a card chip rendered only
	// when effort-mechanical, plus a drawer row; never column or scheduling
	// meaning. Carries the same raw provenance as domain/route so the chip can
	// say what was declared, not just what it normalized to.
	EffortEstimate             string `json:"effortEstimate,omitempty"`
	OriginalEffortEstimate     string `json:"originalEffortEstimate,omitempty"`
	EffortEstimateUnrecognized bool   `json:"effortEstimateUnrecognized,omitempty"`
	// Sweep marker + instance counts (see RequestTicket.Sweep).
	// Display only — a card chip and a drawer row; no column or scheduling
	// meaning.
	Sweep                bool   `json:"sweep,omitempty"`
	SweepInstancesOpen   int    `json:"sweepInstancesOpen,omitempty"`
	SweepInstancesDone   int    `json:"sweepInstancesDone,omitempty"`
	Batch                string `json:"batch"`
	TreeSection          string `json:"treeSection"`
	CreatedAt            string `json:"createdAt"`
	ClaimedAt            string `json:"claimedAt"`
	CompletedAt          string `json:"completedAt"`
	StatusChangedAt      string `json:"statusChangedAt,omitempty"` // last no-dedicated-stamp status flip (see RequestTicket.StatusChangedAt)
	FileModifiedAt       string `json:"fileModifiedAt,omitempty"`  // file mtime at generation, RFC3339 — state-timer fallback only, never completion dating
	CompletionTime       string `json:"completionTime"`
	CompletionTimeSource string `json:"completionTimeSource"`

	// The implementation span the Recently-Done card states: completed_at −
	// claimed_at in minutes, with the read-time rule's verdict already applied
	// Go-side (durations.go's measureImplementationSpan). Present only for a REQ
	// that reached terminal SUCCESS and carries both parseable frontmatter
	// stamps, so `hasImplementationSpan` false is a real "unmeasured" rather than
	// a span of zero. `implementationSpanReason` is "paused" or "reversed", empty
	// when the span reads plainly — the four-hour ceiling itself never crosses to
	// the client, which would make the renderer a second definition of it.
	HasImplementationSpan bool `json:"hasImplementationSpan,omitempty"`
	// Deliberately NOT omitempty: a genuine zero-minute span is possible (identical
	// stamps, or date-only stamps on both fields, which parseTimestamp accepts), and
	// omitempty would drop that 0 while hasImplementationSpan still shipped true —
	// leaving the client to multiply undefined and render "took NaNs". The flag above
	// is what carries presence; this field carries the value, including zero.
	ImplementationSpanMinutes float64 `json:"implementationSpanMinutes"`
	ImplementationSpanReason  string  `json:"implementationSpanReason,omitempty"`

	CompletionAnomaly       bool   `json:"completionAnomaly,omitempty"`
	CompletionAnomalyReason string `json:"completionAnomalyReason,omitempty"`

	// "<field> <raw value>" per frontmatter stamp parsing past the board's
	// generation time + skew allowance (see RequestTicket.FutureTimestampFields).
	FutureTimestampFields []string `json:"futureTimestampFields,omitempty"`

	BodyHtml string `json:"bodyHtml"`

	TestingStatus             string `json:"testingStatus,omitempty"`
	OriginalTestingStatus     string `json:"originalTestingStatus,omitempty"`
	TestingStatusUnrecognized bool   `json:"testingStatusUnrecognized,omitempty"`
	TestedBy                  string `json:"testedBy,omitempty"`
	TestingUpdatedAt          string `json:"testingUpdatedAt,omitempty"`
	TestingFeedback           string `json:"testingFeedback,omitempty"`
}

// generatedUserRequest is one UR node for the by-UR lens, with its grouped REQ
// ids and pre-rendered input.md body.
type generatedUserRequest struct {
	UserRequestId    string   `json:"id"`
	Title            string   `json:"title"`
	InputFilePresent bool     `json:"inputFilePresent"`
	RequestIds       []string `json:"requestIds"`
	BodyHtml         string   `json:"bodyHtml"`
}

// generatedBoardMarkdownData is the lazy raw-source payload used only by the
// drawer Copy button. Keeping it separate avoids loading every Markdown body
// twice (source + rendered HTML) during the initial board paint.
type generatedBoardMarkdownData struct {
	Requests     map[string]string `json:"requests"`
	UserRequests map[string]string `json:"userRequests"`
}

// generatedNote is one do-work/notes.md line. The text stays plain — notes are
// rendered with textContent, not as Markdown, so a hand-typed `<script>` in a
// note is inert and no renderer is needed for a one-line hint.
type generatedNote struct {
	NoteDate string `json:"date"`
	NoteText string `json:"text"`
}

// generatedCalendarEntry plots one REQ on the calendar timeline. `status` rides
// along so a chip colours itself and recentlyDoneIds filters the array without a
// per-entry requestsById lookup; `entryTime` is a completion instant, a claim
// instant, or empty, per the entry's day bucket (see buildCalendar).
type generatedCalendarEntry struct {
	RequestId  string `json:"id"`
	Status     string `json:"status"`
	EntryTime  string `json:"entryTime"`
	DayKey     string `json:"dayKey"`
	TimeSource string `json:"timeSource"`
}

// generatedDurations is the Durations view's data: one measured sample per
// archived REQ that carries both stamps, and one entry per active day. Panels A
// and B disagree about which samples count on purpose — see durations.go.
type generatedDurations struct {
	Samples []generatedDurationSample `json:"samples"`
	Days    []generatedDurationDay    `json:"days"`
}

// generatedDurationSample is one REQ's raw, signed wall span. `excludedReason`
// is "paused" or "reversed" when the calibration's read-time rule holds it out
// of the day medians, and empty when it counts — panel A plots it either way.
// Direct-label placement is NOT here: the renderer decides it, because sizing a
// label needs the width the engine actually draws (REQ-292). Nothing in this
// payload says which marks get a label or how many did not.
type generatedDurationSample struct {
	RequestId      string  `json:"id"`
	Route          string  `json:"route"`
	CompletionTime string  `json:"completionTime"`
	DayKey         string  `json:"dayKey"`
	WallMinutes    float64 `json:"wallMinutes"`
	ExcludedReason string  `json:"excludedReason,omitempty"`
}

// generatedDurationDay carries both figures a day needs: the ruled median (with
// its sample size) and the unruled completion count. `hasMedian` false is a real
// state — every sample that day was rule-excluded — and is not a median of zero.
type generatedDurationDay struct {
	DayKey         string  `json:"dayKey"`
	DayTime        string  `json:"dayTime"`
	MedianMinutes  float64 `json:"medianMinutes"`
	HasMedian      bool    `json:"hasMedian"`
	KeptCount      int     `json:"keptCount"`
	CompletedCount int     `json:"completedCount"`
}

// generatedTimeline is the Timeline view's data: one row per REQ that carries a
// parseable created_at, plus the range those rows span and the single instant
// every open span was measured against.
type generatedTimeline struct {
	Rows       []generatedTimelineRow      `json:"rows"`
	RangeStart string                      `json:"rangeStart"`
	RangeEnd   string                      `json:"rangeEnd"`
	Now        string                      `json:"now"`
	Projection generatedTimelineProjection `json:"projection"`
}

// generatedTimelineProjection is the forward half of the view: where each
// unstarted REQ lands in a serial chain, which REQs cannot be scheduled at all,
// and every figure the forecast rests on so the view can state its assumptions
// rather than imply them. `confident` false means the history was too thin —
// `rows` and `queueEnd` are then empty and `declinedReason` says why.
type generatedTimelineProjection struct {
	Rows           []generatedTimelineProjectedRow `json:"rows"`
	Excluded       []generatedTimelineExclusion    `json:"excluded"`
	QueueEnd       string                          `json:"queueEnd,omitempty"`
	ChainStart     string                          `json:"chainStart"`
	WindowSize     int                             `json:"windowSize"`
	WindowSamples  int                             `json:"windowSamples"`
	TrivialSamples int                             `json:"trivialSamples"`
	NormalSamples  int                             `json:"normalSamples"`
	TrivialMinutes float64                         `json:"trivialMinutes"`
	NormalMinutes  float64                         `json:"normalMinutes"`
	MinimumSamples int                             `json:"minimumSamples"`
	Confident      bool                            `json:"confident"`
	DeclinedReason string                          `json:"declinedReason,omitempty"`
}

// generatedTimelineProjectedRow is one unstarted REQ's forecast slot, keyed by id
// onto the measured row it belongs to.
type generatedTimelineProjectedRow struct {
	RequestId string `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Bucket    string `json:"bucket"`
	Position  int    `json:"position"`
}

// generatedTimelineExclusion names a REQ the projection refuses to schedule and
// why, in plain words rather than by echoing its status.
type generatedTimelineExclusion struct {
	RequestId string `json:"id"`
	Reason    string `json:"reason"`
}

// generatedTimelineRow carries timing only. Title, route, status and domain stay
// in `requests` where every other view already reads them — a second copy here
// would be a second thing to keep in step, and the client holds both.
//
// `waitMinutes` and `workMinutes` are SIGNED: a negative span is the board's
// reversed-stamp anomaly, reported through `anomaly`/`anomalyReason` copied from
// the ticket, and clamping it here would erase the very thing the reader needs
// to see.
type generatedTimelineRow struct {
	RequestId     string  `json:"id"`
	CreatedTime   string  `json:"createdTime"`
	ClaimedTime   string  `json:"claimedTime,omitempty"`
	CompletedTime string  `json:"completedTime,omitempty"`
	WaitMinutes   float64 `json:"waitMinutes"`
	WaitOpen      bool    `json:"waitOpen,omitempty"`
	HasWork       bool    `json:"hasWork"`
	WorkMinutes   float64 `json:"workMinutes"`
	WorkOpen      bool    `json:"workOpen,omitempty"`
	Anomaly       bool    `json:"anomaly,omitempty"`
	AnomalyReason string  `json:"anomalyReason,omitempty"`
}

type staticSiteOutput struct {
	TargetName   string
	PayloadBytes []byte
}

type staticSiteBackup struct {
	TargetPath string
	BackupPath string
}

// generateStaticSite writes a three-file static board into outputDirectory:
//   - index.html — the page shell with CSS + the assembled client inlined; references board-data.js
//   - board-data.js — the initial board payload (including pre-rendered body HTML)
//   - board-markdown.js — raw REQ/UR bodies, loaded lazily on the first Copy click
//
// All three files together are self-contained and open directly from disk (file://) or
// any static server with zero build steps.
func generateStaticSite(outputDirectory string, board *Board) error {
	return generateStaticSiteWithPublisher(outputDirectory, board, os.Rename)
}

// generateStaticSiteWithPublisher keeps the real generate path intact while
// allowing tests to seed a publication rename failure deterministically.
func generateStaticSiteWithPublisher(outputDirectory string, board *Board, publishFile func(string, string) error) error {
	if strings.TrimSpace(outputDirectory) == "" {
		return fmt.Errorf("queue-kanban: generate requires a non-empty --out directory")
	}

	boardData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		return buildError
	}
	// The snapshot carries what verify sees at generate time. buildGeneratedBoardData
	// keeps its board-only signature — a dozen tests call it — so the findings are
	// folded in here, at the two real callers, rather than threaded through it.
	attachVerifyFindings(&boardData, board, time.Now())
	boardMarkdownData := buildGeneratedBoardMarkdownData(board)

	boardDataJs, encodeError := encodeBoardDataForJsAssignment(boardData)
	if encodeError != nil {
		return encodeError
	}

	boardMarkdownJs, markdownEncodeError := encodeBoardMarkdownForJsAssignment(boardMarkdownData)
	if markdownEncodeError != nil {
		return markdownEncodeError
	}

	pageHtml, assembleError := assembleStaticPage(board.GeneratedAt, board.ProjectName)
	if assembleError != nil {
		return assembleError
	}

	if mkdirError := os.MkdirAll(outputDirectory, 0o755); mkdirError != nil {
		return fmt.Errorf("queue-kanban: cannot create --out directory %s: %w", outputDirectory, mkdirError)
	}

	staticOutputs := [3]staticSiteOutput{
		{TargetName: boardDataJsFilename, PayloadBytes: []byte(boardDataJs)},
		{TargetName: boardMarkdownJsFilename, PayloadBytes: []byte(boardMarkdownJs)},
		{TargetName: "index.html", PayloadBytes: []byte(pageHtml)},
	}
	return publishStaticSiteOutputs(outputDirectory, staticOutputs, publishFile)
}

// publishStaticSiteOutputs publishes exactly the static board's three files.
// The renames are not cross-file atomic, but a handled failure restores every
// target to its pre-invocation bytes before returning.
func publishStaticSiteOutputs(outputDirectory string, staticOutputs [3]staticSiteOutput, publishFile func(string, string) error) (returnError error) {
	privateDirectory, temporaryDirectoryError := os.MkdirTemp(outputDirectory, ".queue-kanban-static-")
	if temporaryDirectoryError != nil {
		return fmt.Errorf("queue-kanban: cannot create private static-site directory: %w", temporaryDirectoryError)
	}
	removePrivateDirectory := true
	defer func() {
		if !removePrivateDirectory {
			return
		}
		if cleanupError := os.RemoveAll(privateDirectory); cleanupError != nil {
			returnError = errors.Join(returnError, fmt.Errorf("queue-kanban: cannot remove private static-site directory %s: %w", privateDirectory, cleanupError))
		}
	}()

	stagedPaths := make([]string, 0, len(staticOutputs))
	for outputIndex, staticOutput := range staticOutputs {
		stagedPath := filepath.Join(privateDirectory, fmt.Sprintf("staged-%d-%s", outputIndex, staticOutput.TargetName))
		if writeError := os.WriteFile(stagedPath, staticOutput.PayloadBytes, 0o644); writeError != nil {
			return fmt.Errorf("queue-kanban: cannot stage %s: %w", staticOutput.TargetName, writeError)
		}
		stagedPaths = append(stagedPaths, stagedPath)
	}

	backups := make([]staticSiteBackup, 0, len(staticOutputs))
	for outputIndex, staticOutput := range staticOutputs {
		targetPath := filepath.Join(outputDirectory, staticOutput.TargetName)
		if _, statError := os.Lstat(targetPath); statError != nil {
			if os.IsNotExist(statError) {
				continue
			}
			restoreError := restoreStaticSiteTargets(nil, backups)
			if restoreError != nil {
				removePrivateDirectory = false
			}
			return errors.Join(
				fmt.Errorf("queue-kanban: cannot inspect %s before publication: %w", targetPath, statError),
				restoreError,
			)
		}

		backupPath := filepath.Join(privateDirectory, fmt.Sprintf("backup-%d-%s", outputIndex, staticOutput.TargetName))
		if backupError := os.Rename(targetPath, backupPath); backupError != nil {
			restoreError := restoreStaticSiteTargets(nil, backups)
			if restoreError != nil {
				removePrivateDirectory = false
			}
			return errors.Join(
				fmt.Errorf("queue-kanban: cannot back up %s: %w", targetPath, backupError),
				restoreError,
			)
		}
		backups = append(backups, staticSiteBackup{TargetPath: targetPath, BackupPath: backupPath})
	}

	publishedTargets := make([]string, 0, len(staticOutputs))
	for outputIndex, staticOutput := range staticOutputs {
		targetPath := filepath.Join(outputDirectory, staticOutput.TargetName)
		if publicationError := publishFile(stagedPaths[outputIndex], targetPath); publicationError != nil {
			targetsToRemove := append(publishedTargets, targetPath)
			restoreError := restoreStaticSiteTargets(targetsToRemove, backups)
			if restoreError != nil {
				removePrivateDirectory = false
			}
			return errors.Join(
				fmt.Errorf("queue-kanban: cannot publish %s: %w", targetPath, publicationError),
				restoreError,
			)
		}
		publishedTargets = append(publishedTargets, targetPath)
	}

	return nil
}

func restoreStaticSiteTargets(publishedTargets []string, backups []staticSiteBackup) error {
	var restoreErrors []error
	for _, targetPath := range publishedTargets {
		if removeError := os.Remove(targetPath); removeError != nil && !os.IsNotExist(removeError) {
			restoreErrors = append(restoreErrors, fmt.Errorf("remove replacement %s: %w", targetPath, removeError))
		}
	}
	for _, backup := range backups {
		if restoreError := os.Rename(backup.BackupPath, backup.TargetPath); restoreError != nil {
			restoreErrors = append(restoreErrors, fmt.Errorf("restore %s: %w", backup.TargetPath, restoreError))
		}
	}
	if len(restoreErrors) == 0 {
		return nil
	}
	return fmt.Errorf("queue-kanban: static-site rollback failed: %w", errors.Join(restoreErrors...))
}

// generatedVerifyFinding is one verify finding as the page consumes it. `fixable`
// keeps verify's exact meaning — `do-work cleanup` can mechanically resolve it —
// because an inflated count sends the reader to a command that will not help.
type generatedVerifyFinding struct {
	Category string `json:"category"`
	Detail   string `json:"detail"`
	Remedy   string `json:"remedy,omitempty"`
	Fixable  bool   `json:"fixable,omitempty"`
}

// boardRenderedVerifyCategories are the findings the board already shows by other
// means, so forwarding them would print the same prose a second or third time:
// duplicate ids and stray files arrive through board.Warnings, completion
// anomalies additionally through their own column and a per-card badge. An
// unrecognized status already reaches the warning banner, its parked Needs input
// / Blocked column, and that card's invalid-status badge.
//
// Suppression happens here, in the producer, so the client can render the list
// blindly and no second copy of this judgment lives in JavaScript.
var boardRenderedVerifyCategories = map[string]bool{
	verifyCategoryCompletionAnomaly:         true,
	verifyCategoryDuplicateRequestId:        true,
	verifyCategoryStrayRequestFile:          true,
	verifyCategoryUnrecognizedRequestStatus: true,
}

// attachVerifyFindings runs the probe set against an already-built board and folds
// the result into the payload. Both callers use it — generate for the static
// snapshot and serve per request — so the suppression list and the path reduction
// below have exactly one home.
func attachVerifyFindings(data *generatedBoardData, board *Board, now time.Time) {
	report := collectVerifyFindings(board.RepoRoot, board, now)
	for _, finding := range report.Findings {
		if boardRenderedVerifyCategories[finding.Category] {
			continue
		}
		data.VerifyFindings = append(data.VerifyFindings, generatedVerifyFinding{
			Category: finding.Category,
			Detail:   reduceAbsolutePaths(finding.Detail, board.RepoRoot),
			Remedy:   reduceAbsolutePaths(finding.Remedy, board.RepoRoot),
			Fixable:  finding.Fixable,
		})
	}
	for _, skipped := range report.SkippedProbes {
		data.VerifySkipped = append(data.VerifySkipped, reduceAbsolutePaths(skipped, board.RepoRoot))
	}
}

// reduceAbsolutePaths strips this machine's filesystem layout out of text bound for
// the page. A static snapshot is shareable, and a path that means something here
// means something else — or nothing — on the machine that opens it. The worktree
// probe is the one that makes this necessary: its detail carries `git worktree list
// --porcelain` output, which is absolute at the source.
//
// The repo root becomes a repo-relative path; any other absolute path is replaced
// wholesale, because a path outside the repo says even more about this machine than
// one inside it. The CLI report keeps its absolute paths — they are useful next to a
// shell on the machine that produced them.
func reduceAbsolutePaths(text string, repoRoot string) string {
	if text == "" {
		return text
	}
	reduced := text
	if repoRoot != "" {
		reduced = strings.ReplaceAll(reduced, repoRoot+string(os.PathSeparator), "")
		reduced = strings.ReplaceAll(reduced, repoRoot+"/", "")
		reduced = strings.ReplaceAll(reduced, repoRoot, ".")
	}
	return remainingAbsolutePath.ReplaceAllString(reduced, "${1}<path outside this repository>")
}

// Any surviving absolute path: a POSIX root-anchored run, or a Windows drive letter
// followed by EITHER slash. Git on Windows emits both — `C:\Users\...` from some
// commands and `C:/Users/...` from others (rev-parse --show-toplevel among them) — and
// a drive pattern that accepts only the backslash form lets the forward-slash one
// through untouched: `C:` is not a boundary character, so the POSIX branch cannot
// match the `/` after it either, and the whole path ships inside a static board.
//
// The leading group is load-bearing. RE2 has no lookbehind, so the boundary is
// captured and put back rather than asserted — without it the pattern matches the
// `/` INSIDE an already-relative path and turns `do-work/calibration-log.tsv` into
// `do-work<path outside this repository>`, mangling exactly the relative paths the
// repo-root reduction just produced.
var remainingAbsolutePath = regexp.MustCompile(`(\A|[\s"'(\[])((?:[A-Za-z]:[\\/]|/)[^\s"'` + "`" + `]*)`)

// buildGeneratedBoardData projects the parsed Board into the JSON data island,
// pre-rendering every REQ and UR body to HTML along the way.
func buildGeneratedBoardData(board *Board) (generatedBoardData, error) {
	data := generatedBoardData{
		GeneratedAt:      formatTimestamp(board.GeneratedAt),
		Warnings:         board.Warnings,
		TestingProfiles:  board.TestingProfiles,
		RepoFileMentions: collectRepoFileMentions(board),
		Requests:         map[string]generatedRequest{},
		UserRequests:     map[string]generatedUserRequest{},
		Columns: generatedColumns{
			Pending:             requestIdsOf(board.Columns.Pending),
			PendingReady:        requestIdsOf(board.Columns.PendingReady),
			PendingWaiting:      requestIdsOf(board.Columns.PendingWaiting),
			Claimed:             requestIdsOf(board.Columns.Claimed),
			NeedsInputOrBlocked: requestIdsOf(board.Columns.NeedsInputOrBlocked),
			RecentlyDone:        requestIdsOf(board.Columns.RecentlyDone),
			CompletionAnomalies: requestIdsOf(board.Columns.CompletionAnomalies),
		},
	}

	for _, ticket := range board.AllRequests {
		bodyHtml, renderError := renderMarkdownBodyToHtml(ticket.BodyMarkdown)
		if renderError != nil {
			return generatedBoardData{}, fmt.Errorf("queue-kanban: rendering %s body: %w", ticket.RequestId, renderError)
		}
		data.RequestOrder = append(data.RequestOrder, ticket.RequestId)
		// Terminal SUCCESS only: cancelled work shares the Recently-Done column
		// but never took an implementation span worth stating.
		implementationSpan := ImplementationSpan{}
		if isCompletedStatus(ticket.Status) {
			implementationSpan = measureImplementationSpan(ticket)
		}
		data.Requests[ticket.RequestId] = generatedRequest{
			RequestId:                  ticket.RequestId,
			Title:                      ticket.Title,
			Status:                     ticket.Status,
			OriginalStatus:             ticket.OriginalStatus,
			StatusUnrecognized:         ticket.StatusUnrecognized,
			Error:                      ticket.Error,
			ErrorType:                  ticket.ErrorType,
			OriginalErrorType:          ticket.OriginalErrorType,
			ErrorTypeUnrecognized:      ticket.ErrorTypeUnrecognized,
			Domain:                     ticket.Domain,
			OriginalDomain:             ticket.OriginalDomain,
			DomainUnrecognized:         ticket.DomainUnrecognized,
			UserRequestId:              ticket.UserRequestId,
			DependsOn:                  ticket.DependsOn,
			UnmetDependencies:          ticket.UnmetDependencies,
			Dependents:                 ticket.Dependents,
			BlockedBy:                  ticket.BlockedBy,
			BlockedAt:                  ticket.BlockedAt,
			BlockedCheck:               ticket.BlockedCheck,
			Related:                    ticket.Related,
			WriteSet:                   ticket.WriteSet,
			WriteSetOverlaps:           ticket.WriteSetOverlaps,
			AssignedTo:                 ticket.AssignedTo,
			Route:                      ticket.Route,
			OriginalRoute:              ticket.OriginalRoute,
			RouteUnrecognized:          ticket.RouteUnrecognized,
			Impact:                     ticket.Impact,
			OriginalImpact:             ticket.OriginalImpact,
			ImpactUnrecognized:         ticket.ImpactUnrecognized,
			EffortEstimate:             ticket.EffortEstimate,
			OriginalEffortEstimate:     ticket.OriginalEffortEstimate,
			EffortEstimateUnrecognized: ticket.EffortEstimateUnrecognized,
			Sweep:                      ticket.Sweep,
			SweepInstancesOpen:         ticket.SweepInstancesOpen,
			SweepInstancesDone:         ticket.SweepInstancesDone,
			Batch:                      ticket.Batch,
			TreeSection:                ticket.TreeSection,
			CreatedAt:                  ticket.CreatedAt,
			ClaimedAt:                  ticket.ClaimedAt,
			CompletedAt:                ticket.CompletedAt,
			StatusChangedAt:            ticket.StatusChangedAt,
			FileModifiedAt:             formatTimestamp(ticket.FileModifiedAt),
			CompletionTime:             formatTimestamp(ticket.CompletionTime),
			CompletionTimeSource:       string(ticket.CompletionTimeSource),

			HasImplementationSpan:     implementationSpan.StampsParsed,
			ImplementationSpanMinutes: implementationSpan.WallMinutes,
			ImplementationSpanReason:  implementationSpan.ExclusionReason,

			CompletionAnomaly:       ticket.CompletionAnomaly,
			CompletionAnomalyReason: ticket.CompletionAnomalyReason,

			FutureTimestampFields: ticket.FutureTimestampFields,

			BodyHtml: bodyHtml,

			TestingStatus:             ticket.TestingStatus,
			OriginalTestingStatus:     ticket.OriginalTestingStatus,
			TestingStatusUnrecognized: ticket.TestingStatusUnrecognized,
			TestedBy:                  ticket.TestedBy,
			TestingUpdatedAt:          ticket.TestingUpdatedAt,
			TestingFeedback:           ticket.TestingFeedback,
		}
	}

	for _, userRequest := range board.UserRequests {
		bodyHtml, renderError := renderMarkdownBodyToHtml(userRequest.BodyMarkdown)
		if renderError != nil {
			return generatedBoardData{}, fmt.Errorf("queue-kanban: rendering %s body: %w", userRequest.UserRequestId, renderError)
		}
		data.UserRequestOrder = append(data.UserRequestOrder, userRequest.UserRequestId)
		data.UserRequests[userRequest.UserRequestId] = generatedUserRequest{
			UserRequestId:    userRequest.UserRequestId,
			Title:            userRequest.Title,
			InputFilePresent: userRequest.InputFilePresent,
			RequestIds:       userRequest.RequestIds,
			BodyHtml:         bodyHtml,
		}
	}

	for _, note := range board.Notes {
		data.Notes = append(data.Notes, generatedNote{
			NoteDate: note.NoteDate,
			NoteText: note.NoteText,
		})
	}

	for _, entry := range board.Calendar {
		data.Calendar = append(data.Calendar, generatedCalendarEntry{
			RequestId:  entry.RequestId,
			Status:     entry.Status,
			EntryTime:  formatTimestamp(entry.EntryTime),
			DayKey:     entry.DayKey,
			TimeSource: string(entry.TimeSource),
		})
	}

	durationAggregate := buildDurationAggregate(board.AllRequests)
	for _, sample := range durationAggregate.Samples {
		data.Durations.Samples = append(data.Durations.Samples, generatedDurationSample{
			RequestId:      sample.RequestId,
			Route:          sample.Route,
			CompletionTime: formatTimestamp(sample.CompletionTime),
			DayKey:         sample.DayKey,
			WallMinutes:    sample.WallMinutes,
			ExcludedReason: sample.DayMedianExclusion,
		})
	}

	// The board's generation instant is the view's one now: it measures every
	// open span here and positions the client's now-line, so the two cannot
	// disagree about where "still running" ends.
	timelineAggregate := buildTimelineAggregate(board.AllRequests, board.GeneratedAt)
	for _, row := range timelineAggregate.Rows {
		data.Timeline.Rows = append(data.Timeline.Rows, generatedTimelineRow{
			RequestId:     row.RequestId,
			CreatedTime:   formatTimestamp(row.CreatedTime),
			ClaimedTime:   formatTimestamp(row.ClaimedTime),
			CompletedTime: formatTimestamp(row.CompletedTime),
			WaitMinutes:   row.WaitMinutes,
			WaitOpen:      row.WaitOpen,
			HasWork:       row.HasWork,
			WorkMinutes:   row.WorkMinutes,
			WorkOpen:      row.WorkOpen,
			Anomaly:       row.Anomaly,
			AnomalyReason: row.AnomalyReason,
		})
	}
	data.Timeline.RangeStart = formatTimestamp(timelineAggregate.RangeStart)
	data.Timeline.RangeEnd = formatTimestamp(timelineAggregate.RangeEnd)
	data.Timeline.Now = formatTimestamp(timelineAggregate.Now)

	timelineProjection := buildTimelineProjection(board.AllRequests, durationAggregate, board.GeneratedAt)
	for _, row := range timelineProjection.Rows {
		data.Timeline.Projection.Rows = append(data.Timeline.Projection.Rows, generatedTimelineProjectedRow{
			RequestId: row.RequestId,
			StartTime: formatTimestamp(row.StartTime),
			EndTime:   formatTimestamp(row.EndTime),
			Bucket:    row.Bucket,
			Position:  row.Position,
		})
	}
	for _, exclusion := range timelineProjection.Excluded {
		data.Timeline.Projection.Excluded = append(data.Timeline.Projection.Excluded, generatedTimelineExclusion{
			RequestId: exclusion.RequestId,
			Reason:    exclusion.Reason,
		})
	}
	data.Timeline.Projection.QueueEnd = formatTimestamp(timelineProjection.QueueEnd)
	data.Timeline.Projection.ChainStart = formatTimestamp(timelineProjection.ChainStart)
	data.Timeline.Projection.WindowSize = timelineProjection.WindowSize
	data.Timeline.Projection.WindowSamples = timelineProjection.WindowSampleCount
	data.Timeline.Projection.TrivialSamples = timelineProjection.TrivialSampleCount
	data.Timeline.Projection.NormalSamples = timelineProjection.NormalSampleCount
	data.Timeline.Projection.TrivialMinutes = timelineProjection.TrivialMedianMinutes
	data.Timeline.Projection.NormalMinutes = timelineProjection.NormalMedianMinutes
	data.Timeline.Projection.MinimumSamples = timelineProjection.MinimumSamples
	data.Timeline.Projection.Confident = timelineProjection.Confident
	data.Timeline.Projection.DeclinedReason = timelineProjection.DeclinedReason
	for _, day := range durationAggregate.Days {
		data.Durations.Days = append(data.Durations.Days, generatedDurationDay{
			DayKey:         day.DayKey,
			DayTime:        formatTimestamp(day.DayTime),
			MedianMinutes:  day.MedianMinutes,
			HasMedian:      day.HasMedian,
			KeptCount:      day.KeptCount,
			CompletedCount: day.CompletedCount,
		})
	}

	return data, nil
}

// buildGeneratedBoardMarkdownData projects each ticket's FILE TEXT — the original
// frontmatter fence followed by the body — into a compact id-keyed payload. It is
// the whole file, not just the body, so a Copy from the drawer round-trips: the
// paste can be saved straight back as a valid REQ or UR file.
//
// It is separated from buildGeneratedBoardData so the initial page does not
// download source text that is only needed after a Copy click. The wire shape is
// unchanged (map[string]string) — only the value grew by the fence.
func buildGeneratedBoardMarkdownData(board *Board) generatedBoardMarkdownData {
	markdownData := generatedBoardMarkdownData{
		Requests:     map[string]string{},
		UserRequests: map[string]string{},
	}
	for _, ticket := range board.AllRequests {
		markdownData.Requests[ticket.RequestId] = ticket.FrontmatterMarkdown + ticket.BodyMarkdown
	}
	for _, userRequest := range board.UserRequests {
		// A synthesized UR (no input.md on disk) has no file text to offer. The
		// frontend treats key PRESENCE as "the real file is available" and copies
		// the value verbatim, so an empty entry here would put an empty string on
		// the clipboard instead of triggering the rendered-text fallback.
		if !userRequest.InputFilePresent {
			continue
		}
		markdownData.UserRequests[userRequest.UserRequestId] = userRequest.FrontmatterMarkdown + userRequest.BodyMarkdown
	}
	return markdownData
}

// assembleBoardJavaScript replaces the private shell's single placeholder with
// the eight manifest-ordered closure fragments. Each fragment owns exactly one
// trailing LF; this seam owns the blank line between fragments and before boot.
func assembleBoardJavaScript(webAssets fs.FS) ([]byte, error) {
	shellText, shellError := fs.ReadFile(webAssets, boardJavaScriptShellPath)
	if shellError != nil {
		return nil, fmt.Errorf("queue-kanban: reading embedded JavaScript shell: %w", shellError)
	}
	placeholderBytes := []byte(boardJavaScriptPlaceholder)
	if placeholderCount := bytes.Count(shellText, placeholderBytes); placeholderCount != 1 {
		return nil, fmt.Errorf("queue-kanban: JavaScript shell has %d fragment placeholders, want exactly 1", placeholderCount)
	}
	placeholderLineBytes := []byte(boardJavaScriptPlaceholderLine)
	if placeholderLineCount := bytes.Count(shellText, placeholderLineBytes); placeholderLineCount != 1 {
		return nil, fmt.Errorf("queue-kanban: JavaScript shell has %d canonical fragment placeholder lines, want exactly 1", placeholderLineCount)
	}

	fragmentTexts := make([][]byte, 0, len(boardJavaScriptFragmentPaths))
	seenFragmentPaths := make(map[string]bool, len(boardJavaScriptFragmentPaths))
	for _, fragmentPath := range boardJavaScriptFragmentPaths {
		if seenFragmentPaths[fragmentPath] {
			return nil, fmt.Errorf("queue-kanban: duplicate JavaScript fragment path %s", fragmentPath)
		}
		seenFragmentPaths[fragmentPath] = true

		fragmentText, fragmentError := fs.ReadFile(webAssets, fragmentPath)
		if fragmentError != nil {
			return nil, fmt.Errorf("queue-kanban: reading embedded JavaScript fragment %s: %w", fragmentPath, fragmentError)
		}
		if !bytes.HasSuffix(fragmentText, []byte("\n")) ||
			bytes.HasSuffix(fragmentText, []byte("\n\n")) ||
			bytes.HasSuffix(fragmentText, []byte("\r\n")) {
			return nil, fmt.Errorf("queue-kanban: JavaScript fragment %s must end with exactly one LF", fragmentPath)
		}
		fragmentTexts = append(fragmentTexts, fragmentText)
	}

	assembledFragments := bytes.Join(fragmentTexts, []byte("\n"))
	assembledFragments = append(assembledFragments, '\n')
	assembledClient := bytes.Replace(shellText, placeholderLineBytes, assembledFragments, 1)
	if bytes.Contains(assembledClient, placeholderBytes) {
		return nil, fmt.Errorf("queue-kanban: assembled JavaScript retained a fragment placeholder")
	}
	return assembledClient, nil
}

// assembleStaticPage inlines the CSS and assembled classic client into the HTML
// template, producing the index.html string. The JSON data island is NOT inlined
// here — it lives in the sibling board-data.js file (written by generateStaticSite)
// and is loaded via <script src="board-data.js"> already present in the template.
// projectName labels which repo this board belongs to (the parent folder name);
// it is HTML-escaped before substitution so an exotic folder name can never break
// out of the <title>/identity markup. Every PROJECT_NAME token is replaced.
func assembleStaticPage(generatedAt time.Time, projectName string) (string, error) {
	templateText, templateError := embeddedWebAssets.ReadFile("web/template.html")
	if templateError != nil {
		return "", fmt.Errorf("queue-kanban: reading embedded template: %w", templateError)
	}
	styleText, styleError := embeddedWebAssets.ReadFile("web/board.css")
	if styleError != nil {
		return "", fmt.Errorf("queue-kanban: reading embedded stylesheet: %w", styleError)
	}
	scriptText, scriptError := assembleBoardJavaScript(embeddedWebAssets)
	if scriptError != nil {
		return "", scriptError
	}

	page := string(templateText)
	page = strings.ReplaceAll(page, projectNamePlaceholder, html.EscapeString(projectName))
	page = strings.Replace(page, generatedAtDisplayPlaceholder, displayGeneratedAt(generatedAt), 1)
	page = strings.Replace(page, inlineStylePlaceholder, string(styleText), 1)
	page = strings.Replace(page, inlineScriptPlaceholder, string(scriptText), 1)
	return page, nil
}

// encodeBoardDataForJsAssignment marshals boardData as a JavaScript global
// assignment: window.queueKanbanBoardData = <JSON>;
// HTML escaping is disabled so body HTML (e.g. <h2 id=…>) survives verbatim
// inside the .js file. The </script> neutralization used for inline <script>
// data islands is not needed here because board-data.js is not HTML-parsed.
func encodeBoardDataForJsAssignment(boardData generatedBoardData) (string, error) {
	var jsonBuffer bytes.Buffer
	encoder := json.NewEncoder(&jsonBuffer)
	encoder.SetEscapeHTML(false)
	if encodeError := encoder.Encode(boardData); encodeError != nil {
		return "", fmt.Errorf("queue-kanban: encoding board data for js file: %w", encodeError)
	}
	jsonText := strings.TrimRight(jsonBuffer.String(), "\n")
	return "window.queueKanbanBoardData = " + jsonText + ";\n", nil
}

// encodeBoardMarkdownForJsAssignment emits the lazy raw-source payload as a
// script assignment so both HTTP and file:// boards can load it on demand.
func encodeBoardMarkdownForJsAssignment(markdownData generatedBoardMarkdownData) (string, error) {
	var jsonBuffer bytes.Buffer
	encoder := json.NewEncoder(&jsonBuffer)
	encoder.SetEscapeHTML(false)
	if encodeError := encoder.Encode(markdownData); encodeError != nil {
		return "", fmt.Errorf("queue-kanban: encoding lazy Markdown data for js file: %w", encodeError)
	}
	jsonText := strings.TrimRight(jsonBuffer.String(), "\n")
	return "window.queueKanbanBoardMarkdownData = " + jsonText + ";\n", nil
}

// requestIdsOf projects a column's tickets to their REQ ids, preserving order.
func requestIdsOf(tickets []*RequestTicket) []string {
	ids := make([]string, 0, len(tickets))
	for _, ticket := range tickets {
		ids = append(ids, ticket.RequestId)
	}
	return ids
}

// formatTimestamp renders an instant as RFC3339 UTC, or "" for the zero time so
// the JSON carries an empty string the client can test rather than a bogus year.
func formatTimestamp(instant time.Time) string {
	if instant.IsZero() {
		return ""
	}
	return instant.UTC().Format(time.RFC3339)
}

// displayGeneratedAt formats the board's generation instant for the human-facing
// "Generated …" line in the top bar.
func displayGeneratedAt(instant time.Time) string {
	return instant.UTC().Format("2006-01-02 15:04 MST")
}
