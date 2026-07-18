package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// embeddedWebAssets holds the hand-authored static frontend (HTML shell + CSS +
// JS) that `generate` inlines into a single self-contained index.html. REQ-1209
// (`serve`) embeds and re-serves the SAME web/ directory unchanged, so the shape
// here is the shared contract.
//
//go:embed web/template.html web/board.css web/board.js
var embeddedWebAssets embed.FS

// Inline placeholder tokens in web/template.html. They are deliberately
// comment-shaped so the template stays valid HTML/CSS/JS on its own (it can be
// opened during development without the inlining step). generate replaces each
// exactly once.
const (
	inlineStylePlaceholder        = "/* INLINE_BOARD_STYLES */"
	inlineScriptPlaceholder       = "/* INLINE_BOARD_SCRIPT */"
	generatedAtDisplayPlaceholder = "GENERATED_AT_DISPLAY"
	projectNamePlaceholder        = "PROJECT_NAME"
)

// boardDataJsFilename is the sibling file written next to index.html that assigns
// the board JSON to window.queueKanbanBoardData. A <script src="board-data.js">
// in index.html loads it before board.js so the board renders offline from file://.
const boardDataJsFilename = "board-data.js"

// boardMarkdownJsFilename carries raw REQ/UR Markdown for the drawer Copy button.
// It is deliberately NOT referenced by index.html: board.js injects it on the
// first copy, keeping duplicate raw bodies out of the initial board payload while
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
	Notes            []generatedNote                 `json:"notes,omitempty"`    // do-work/notes.md lines — rendered as a strip above the queue
	Warnings         []string                        `json:"warnings,omitempty"` // duplicate ids / unrecognized statuses — rendered as a banner

	TestingProfiles []string `json:"testingProfiles,omitempty"` // do-work/testers.md profiles for the testing view's tester picker
	// True only when served by the live server (serve.go sets it): the testing
	// view's write actions need the /api/testing/* endpoints, so a static
	// snapshot renders the view read-only.
	LiveTestingApi bool `json:"liveTestingApi,omitempty"`
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
	RequestId            string   `json:"id"`
	Title                string   `json:"title"`
	Status               string   `json:"status"`
	OriginalStatus       string   `json:"originalStatus"`
	StatusUnrecognized   bool     `json:"statusUnrecognized,omitempty"`
	Domain               string   `json:"domain"`
	UserRequestId        string   `json:"userRequestId"`
	DependsOn            []string `json:"dependsOn"`
	UnmetDependencies    []string `json:"unmetDependencies"`
	Dependents           []string `json:"dependents"`
	BlockedBy            []string `json:"blockedBy"`
	Related              []string `json:"related"`
	Route                string   `json:"route"`
	Batch                string   `json:"batch"`
	TreeSection          string   `json:"treeSection"`
	CreatedAt            string   `json:"createdAt"`
	ClaimedAt            string   `json:"claimedAt"`
	CompletedAt          string   `json:"completedAt"`
	ReservedFor          string   `json:"reservedFor,omitempty"`
	ReservedAt           string   `json:"reservedAt,omitempty"`
	ReservationStale     bool     `json:"reservationStale,omitempty"`
	CompletionTime       string   `json:"completionTime"`
	CompletionTimeSource string   `json:"completionTimeSource"`

	CompletionAnomaly       bool   `json:"completionAnomaly,omitempty"`
	CompletionAnomalyReason string `json:"completionAnomalyReason,omitempty"`

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

// generatedCalendarEntry plots one completed REQ on the completion timeline.
type generatedCalendarEntry struct {
	RequestId      string `json:"id"`
	CompletionTime string `json:"completionTime"`
	DayKey         string `json:"dayKey"`
	TimeSource     string `json:"timeSource"`
}

// generateStaticSite writes a three-file static board into outputDirectory:
//   - index.html — the page shell with CSS + board.js inlined; references board-data.js
//   - board-data.js — the initial board payload (including pre-rendered body HTML)
//   - board-markdown.js — raw REQ/UR bodies, loaded lazily on the first Copy click
//
// All three files together are self-contained and open directly from disk (file://) or
// any static server with zero build steps.
func generateStaticSite(outputDirectory string, board *Board) error {
	if strings.TrimSpace(outputDirectory) == "" {
		return fmt.Errorf("queue-kanban: generate requires a non-empty --out directory")
	}

	boardData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		return buildError
	}
	boardMarkdownData := buildGeneratedBoardMarkdownData(board)

	if mkdirError := os.MkdirAll(outputDirectory, 0o755); mkdirError != nil {
		return fmt.Errorf("queue-kanban: cannot create --out directory %s: %w", outputDirectory, mkdirError)
	}

	boardDataJs, encodeError := encodeBoardDataForJsAssignment(boardData)
	if encodeError != nil {
		return encodeError
	}
	boardDataPath := filepath.Join(outputDirectory, boardDataJsFilename)
	if writeError := os.WriteFile(boardDataPath, []byte(boardDataJs), 0o644); writeError != nil {
		return fmt.Errorf("queue-kanban: cannot write %s: %w", boardDataPath, writeError)
	}

	boardMarkdownJs, markdownEncodeError := encodeBoardMarkdownForJsAssignment(boardMarkdownData)
	if markdownEncodeError != nil {
		return markdownEncodeError
	}
	boardMarkdownPath := filepath.Join(outputDirectory, boardMarkdownJsFilename)
	if writeError := os.WriteFile(boardMarkdownPath, []byte(boardMarkdownJs), 0o644); writeError != nil {
		return fmt.Errorf("queue-kanban: cannot write %s: %w", boardMarkdownPath, writeError)
	}

	pageHtml, assembleError := assembleStaticPage(board.GeneratedAt, board.ProjectName)
	if assembleError != nil {
		return assembleError
	}
	indexPath := filepath.Join(outputDirectory, "index.html")
	if writeError := os.WriteFile(indexPath, []byte(pageHtml), 0o644); writeError != nil {
		return fmt.Errorf("queue-kanban: cannot write %s: %w", indexPath, writeError)
	}
	return nil
}

// buildGeneratedBoardData projects the parsed Board into the JSON data island,
// pre-rendering every REQ and UR body to HTML along the way.
func buildGeneratedBoardData(board *Board) (generatedBoardData, error) {
	data := generatedBoardData{
		GeneratedAt:     formatTimestamp(board.GeneratedAt),
		Warnings:        board.Warnings,
		TestingProfiles: board.TestingProfiles,
		Requests:        map[string]generatedRequest{},
		UserRequests:    map[string]generatedUserRequest{},
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
		data.Requests[ticket.RequestId] = generatedRequest{
			RequestId:            ticket.RequestId,
			Title:                ticket.Title,
			Status:               ticket.Status,
			OriginalStatus:       ticket.OriginalStatus,
			StatusUnrecognized:   ticket.StatusUnrecognized,
			Domain:               ticket.Domain,
			UserRequestId:        ticket.UserRequestId,
			DependsOn:            ticket.DependsOn,
			UnmetDependencies:    ticket.UnmetDependencies,
			Dependents:           ticket.Dependents,
			BlockedBy:            ticket.BlockedBy,
			Related:              ticket.Related,
			Route:                ticket.Route,
			Batch:                ticket.Batch,
			TreeSection:          ticket.TreeSection,
			CreatedAt:            ticket.CreatedAt,
			ClaimedAt:            ticket.ClaimedAt,
			CompletedAt:          ticket.CompletedAt,
			ReservedFor:          ticket.ReservedFor,
			ReservedAt:           ticket.ReservedAt,
			ReservationStale:     ticket.ReservationStale,
			CompletionTime:       formatTimestamp(ticket.CompletionTime),
			CompletionTimeSource: string(ticket.CompletionTimeSource),

			CompletionAnomaly:       ticket.CompletionAnomaly,
			CompletionAnomalyReason: ticket.CompletionAnomalyReason,

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
			RequestId:      entry.RequestId,
			CompletionTime: formatTimestamp(entry.CompletionTime),
			DayKey:         entry.DayKey,
			TimeSource:     string(entry.TimeSource),
		})
	}

	return data, nil
}

// buildGeneratedBoardMarkdownData projects raw Markdown into a compact id-keyed
// payload. It is separated from buildGeneratedBoardData so the initial page does
// not download source text that is only needed after a Copy click.
func buildGeneratedBoardMarkdownData(board *Board) generatedBoardMarkdownData {
	markdownData := generatedBoardMarkdownData{
		Requests:     map[string]string{},
		UserRequests: map[string]string{},
	}
	for _, ticket := range board.AllRequests {
		markdownData.Requests[ticket.RequestId] = ticket.BodyMarkdown
	}
	for _, userRequest := range board.UserRequests {
		markdownData.UserRequests[userRequest.UserRequestId] = userRequest.BodyMarkdown
	}
	return markdownData
}

// assembleStaticPage inlines the CSS and board.js into the HTML template,
// producing the index.html string. The JSON data island is NOT inlined here —
// it lives in the sibling board-data.js file (written by generateStaticSite)
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
	scriptText, scriptError := embeddedWebAssets.ReadFile("web/board.js")
	if scriptError != nil {
		return "", fmt.Errorf("queue-kanban: reading embedded script: %w", scriptError)
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
