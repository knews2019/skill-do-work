package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// generateLiveSiteInDir builds the board against the REAL do-work tree and writes
// the static site into a temp dir, returning the output directory path. The git
// lookup is stubbed (mirrors board_live_test.go) so the whole-tree build is fast
// and deterministic. Tests that need board-data.js or other sibling files should
// call this helper directly.
func generateLiveSiteInDir(t *testing.T) string {
	t.Helper()
	workingDirectory, getwdError := os.Getwd()
	if getwdError != nil {
		t.Fatalf("getwd: %v", getwdError)
	}
	repoRoot, resolveError := resolveRepoRoot(workingDirectory)
	if resolveError != nil {
		t.Fatalf("resolveRepoRoot: %v", resolveError)
	}
	stubGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, stubGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	outputDirectory := t.TempDir()
	if generateError := generateStaticSite(outputDirectory, board); generateError != nil {
		t.Fatalf("generateStaticSite: %v", generateError)
	}
	return outputDirectory
}

// generateLiveSite builds the board and returns the index.html contents. It is a
// convenience wrapper over generateLiveSiteInDir for tests that only need the
// main HTML page.
func generateLiveSite(t *testing.T) string {
	t.Helper()
	outputDirectory := generateLiveSiteInDir(t)
	indexPath := filepath.Join(outputDirectory, "index.html")
	indexBytes, readError := os.ReadFile(indexPath)
	if readError != nil {
		t.Fatalf("reading generated index.html: %v", readError)
	}
	return string(indexBytes)
}

func TestGenerateWritesSelfContainedIndex(t *testing.T) {
	indexHtml := generateLiveSite(t)

	// The page must be self-contained: CSS + JS inlined, no CDN / external asset.
	if !strings.Contains(indexHtml, "<style>") {
		t.Fatalf("generated page has no inlined <style> block")
	}
	for _, externalMarker := range []string{
		`src="http`,
		`src='http`,
		`href="http`,
		`<link rel="stylesheet"`,
		"cdn.",
	} {
		if strings.Contains(indexHtml, externalMarker) {
			t.Fatalf("generated page is not self-contained: found external reference %q", externalMarker)
		}
	}
	// The inlined behaviour script must be present (a known function name).
	if !strings.Contains(indexHtml, "renderColumns") {
		t.Fatalf("inlined board.js behaviour is missing from the page")
	}
	// The display placeholder must have been resolved.
	if strings.Contains(indexHtml, "GENERATED_AT_DISPLAY") {
		t.Fatalf("GENERATED_AT_DISPLAY placeholder was not substituted")
	}
}

func TestGenerateRendersColumnHeaders(t *testing.T) {
	indexHtml := generateLiveSite(t)
	for _, columnHeader := range []string{
		"Pending",
		"Claimed",
		"Needs input",
		"Recently done",
	} {
		if !strings.Contains(indexHtml, columnHeader) {
			t.Fatalf("column header %q not found in the generated page", columnHeader)
		}
	}
}

func TestGenerateEmbedsLivePendingCards(t *testing.T) {
	// After REQ-1213 the card data (including REQ IDs) lives in board-data.js. The
	// expected ids are derived from the live board, not hard-coded — the old test
	// pinned REQ-1207..1210 from the source monorepo, which don't exist in this
	// extraction. Exact seeded-card coverage lives in the synthetic board tests.
	board := liveBoard(t)
	if len(board.Columns.Pending) == 0 {
		t.Skip("no pending REQs in the live tree; nothing to assert")
	}

	outputDirectory := generateLiveSiteInDir(t)
	boardDataPath := filepath.Join(outputDirectory, "board-data.js")
	boardDataBytes, readError := os.ReadFile(boardDataPath)
	if readError != nil {
		t.Fatalf("reading board-data.js: %v", readError)
	}
	boardDataJs := string(boardDataBytes)

	checks := 0
	for _, ticket := range board.Columns.Pending {
		if !strings.Contains(boardDataJs, ticket.RequestId) {
			t.Fatalf("live pending id %q not found in board-data.js", ticket.RequestId)
		}
		checks++
		if checks >= 25 {
			break // a representative sample is enough
		}
	}
}

func TestGenerateEmitsBlockedFields(t *testing.T) {
	// The synthetic tree seeds REQ-9006 as status: blocked with a free-text
	// blocked_by, a blocked_at, and a blocked_check. Those must survive into the
	// generated payload so the frontend can render the "blocked by" badge/drawer.
	board := syntheticBoard(t)
	generatedData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}
	blockedRequest, present := generatedData.Requests["REQ-9006"]
	if !present {
		t.Fatalf("REQ-9006 (blocked) missing from generated requests")
	}
	if blockedRequest.Status != "blocked" {
		t.Fatalf("REQ-9006 status = %q, want blocked", blockedRequest.Status)
	}
	if len(blockedRequest.BlockedBy) != 1 || blockedRequest.BlockedBy[0] != "LM Studio running locally" {
		t.Fatalf("REQ-9006 blockedBy = %+v, want [\"LM Studio running locally\"]", blockedRequest.BlockedBy)
	}
	if blockedRequest.BlockedCheck == "" || blockedRequest.BlockedAt == "" {
		t.Fatalf("REQ-9006 blockedCheck/blockedAt not populated: check=%q at=%q", blockedRequest.BlockedCheck, blockedRequest.BlockedAt)
	}
	// The fields must also survive JSON marshaling under their camelCase keys.
	marshaledBytes, marshalError := json.Marshal(blockedRequest)
	if marshalError != nil {
		t.Fatalf("marshal generated request: %v", marshalError)
	}
	marshaledJson := string(marshaledBytes)
	for _, expectedKey := range []string{`"blockedBy"`, `"blockedAt"`, `"blockedCheck"`} {
		if !strings.Contains(marshaledJson, expectedKey) {
			t.Fatalf("generated JSON missing %s: %s", expectedKey, marshaledJson)
		}
	}
}

func TestGenerateIndexHtmlUnderSizeBudget(t *testing.T) {
	// The JSON data island (all pre-rendered REQ bodies) must be externalized to
	// board-data.js so index.html stays well under 1 MB. Before REQ-1213 the
	// monolithic file weighed ~14 MB.
	const maxIndexHtmlBytes = 1 << 20 // 1 MiB
	indexHtml := generateLiveSite(t)
	actualBytes := len(indexHtml)
	if actualBytes >= maxIndexHtmlBytes {
		t.Fatalf("index.html is %d bytes (%.1f MB) — exceeds the 1 MB budget; externalize the JSON data island to board-data.js",
			actualBytes, float64(actualBytes)/(1<<20))
	}
}

func TestGenerateHasCalendarAndUserRequestLensHooks(t *testing.T) {
	outputDirectory := generateLiveSiteInDir(t)

	indexPath := filepath.Join(outputDirectory, "index.html")
	indexBytes, readError := os.ReadFile(indexPath)
	if readError != nil {
		t.Fatalf("reading generated index.html: %v", readError)
	}
	indexHtml := string(indexBytes)

	if !strings.Contains(indexHtml, `data-view-target="calendar"`) {
		t.Fatalf("calendar view hook not found")
	}
	if !strings.Contains(indexHtml, `data-lens-target="user-request"`) {
		t.Fatalf("by-UR lens toggle hook not found")
	}

	// Calendar day-keyed completion entries live in the externalized board-data.js.
	boardDataPath := filepath.Join(outputDirectory, "board-data.js")
	boardDataBytes, bdReadError := os.ReadFile(boardDataPath)
	if bdReadError != nil {
		t.Fatalf("reading board-data.js: %v", bdReadError)
	}
	if !strings.Contains(string(boardDataBytes), `"dayKey"`) {
		t.Fatalf("calendar entries (dayKey) not found in board-data.js")
	}
}

func TestGenerateEmbedsGoldmarkRenderedBody(t *testing.T) {
	// After REQ-1213 the JSON data island (including pre-rendered bodies) lives in
	// board-data.js, not in index.html. Read the sibling file for assertions.
	outputDirectory := generateLiveSiteInDir(t)
	boardDataPath := filepath.Join(outputDirectory, "board-data.js")
	boardDataBytes, readError := os.ReadFile(boardDataPath)
	if readError != nil {
		t.Fatalf("reading board-data.js: %v", readError)
	}
	boardDataJs := string(boardDataBytes)

	// Every REQ body in this repo has `## ` headings; goldmark (with auto heading
	// IDs) renders them to `<h2 id="...">`. Asserting the id form proves the
	// marker came from a rendered REQ body — not from the page chrome.
	if !strings.Contains(boardDataJs, `<h2 id=`) {
		t.Fatalf("no goldmark-rendered `<h2 id=` body heading found in board-data.js")
	}
	// The data island must carry pre-rendered bodies under the bodyHtml key.
	if !strings.Contains(boardDataJs, `"bodyHtml"`) {
		t.Fatalf("board-data.js has no bodyHtml field")
	}
}

func TestGenerateSeparatesRawMarkdownForLazyCopy(t *testing.T) {
	// Copy still needs exact source, but shipping it beside bodyHtml nearly
	// doubles the initial payload. Raw bodies belong in a lazy sibling script.
	outputDirectory := generateLiveSiteInDir(t)

	boardDataBytes, readError := os.ReadFile(filepath.Join(outputDirectory, "board-data.js"))
	if readError != nil {
		t.Fatalf("reading board-data.js: %v", readError)
	}
	if strings.Contains(string(boardDataBytes), `"bodyMarkdown"`) {
		t.Fatalf("board-data.js still carries bodyMarkdown — raw source must stay out of the initial payload")
	}

	boardMarkdownBytes, markdownReadError := os.ReadFile(filepath.Join(outputDirectory, boardMarkdownJsFilename))
	if markdownReadError != nil {
		t.Fatalf("reading %s: %v", boardMarkdownJsFilename, markdownReadError)
	}
	if !strings.HasPrefix(string(boardMarkdownBytes), "window.queueKanbanBoardMarkdownData = ") {
		t.Fatalf("%s does not assign the lazy Markdown global", boardMarkdownJsFilename)
	}

	indexBytes, indexReadError := os.ReadFile(filepath.Join(outputDirectory, "index.html"))
	if indexReadError != nil {
		t.Fatalf("reading generated index.html: %v", indexReadError)
	}
	if !strings.Contains(string(indexBytes), `id="detail-copy"`) {
		t.Fatalf("detail drawer Copy button (id=\"detail-copy\") not found in index.html")
	}
	if strings.Contains(string(indexBytes), `<script src="board-markdown.js"></script>`) {
		t.Fatalf("index.html eagerly loads board-markdown.js; raw source must load only after Copy")
	}
	if !strings.Contains(string(indexBytes), `markdownScript.src = "board-markdown.js"`) {
		t.Fatalf("inlined board.js has no lazy board-markdown.js loader")
	}
}

func TestBuildGeneratedBoardMarkdownDataKeepsExactSources(t *testing.T) {
	board := &Board{
		AllRequests: []*RequestTicket{
			{RequestId: "REQ-1", BodyMarkdown: "## What\n\n- [ ] keep formatting\n"},
		},
		UserRequests: []*UserRequestTicket{
			{UserRequestId: "UR-1", BodyMarkdown: "# Original request\n\nExact text.\n"},
		},
	}

	markdownData := buildGeneratedBoardMarkdownData(board)
	if got := markdownData.Requests["REQ-1"]; got != board.AllRequests[0].BodyMarkdown {
		t.Fatalf("REQ raw Markdown changed: got %q, want %q", got, board.AllRequests[0].BodyMarkdown)
	}
	if got := markdownData.UserRequests["UR-1"]; got != board.UserRequests[0].BodyMarkdown {
		t.Fatalf("UR raw Markdown changed: got %q, want %q", got, board.UserRequests[0].BodyMarkdown)
	}
}

func TestRenderMarkdownBodyToHtmlHeadingsAndTaskLists(t *testing.T) {
	body := "## What\n\nA paragraph.\n\n- [ ] unchecked item\n- [x] checked item\n"
	rendered, renderError := renderMarkdownBodyToHtml(body)
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if !strings.Contains(rendered, "<h2") {
		t.Fatalf("expected an <h2> from a ## heading, got: %s", rendered)
	}
	if !strings.Contains(rendered, `type="checkbox"`) {
		t.Fatalf("expected GFM task-list checkboxes, got: %s", rendered)
	}
}

func TestRenderMarkdownQuestionOptionsKeepTheirOwnLines(t *testing.T) {
	// The Open Questions format (actions/capture.md) indents Recommended:/Also:
	// continuation lines under the checkbox item; plain Markdown would lazily
	// merge them into the question paragraph. The renderer must emit a <br>
	// before each so they stay separate visual lines in the drawer.
	body := "## Open Questions\n\n" +
		"- [ ] Should I process this as a new task?\n" +
		"  Recommended: Yes, add to queue.\n" +
		"  Also: No, discard it.\n"
	rendered, renderError := renderMarkdownBodyToHtml(body)
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if strings.Count(rendered, "<br") != 2 {
		t.Fatalf("expected 2 hard breaks (before Recommended: and Also:), got: %s", rendered)
	}
	if !strings.Contains(rendered, `type="checkbox"`) {
		t.Fatalf("checkbox item must survive the option-line preprocessing, got: %s", rendered)
	}
}

func TestRenderMarkdownLeavesCodeFencesVerbatim(t *testing.T) {
	// A fenced block whose content happens to start with an option keyword must
	// not have hard-break backslashes injected into its verbatim content.
	body := "```\nsome output\nRecommended: not a question option\n```\n"
	rendered, renderError := renderMarkdownBodyToHtml(body)
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if strings.Contains(rendered, "\\") || strings.Contains(rendered, "<br") {
		t.Fatalf("code fence content must stay verbatim, got: %s", rendered)
	}
}

func TestRenderMarkdownEscapesRawHtml(t *testing.T) {
	rendered, renderError := renderMarkdownBodyToHtml("a <script>alert(1)</script> b")
	if renderError != nil {
		t.Fatalf("renderMarkdownBodyToHtml: %v", renderError)
	}
	if strings.Contains(rendered, "<script>") {
		t.Fatalf("raw <script> should be escaped, got: %s", rendered)
	}
}

// TestEncodeBoardDataJsAssignmentPreservesRawHtml covers the one encoder both
// generate and serve actually ship (board-data.js is a plain .js file, never
// HTML-parsed, so no </script> neutralization is involved): the assignment
// wrapper must be exact and pre-rendered body HTML must survive unescaped
// (SetEscapeHTML off — the goldmark proof the GREEN test greps for).
func TestEncodeBoardDataJsAssignmentPreservesRawHtml(t *testing.T) {
	data := generatedBoardData{
		Requests: map[string]generatedRequest{
			"REQ-1": {RequestId: "REQ-1", BodyHtml: "<h2>Lessons & Notes</h2>"},
		},
	}
	encoded, encodeError := encodeBoardDataForJsAssignment(data)
	if encodeError != nil {
		t.Fatalf("encodeBoardDataForJsAssignment: %v", encodeError)
	}
	if !strings.HasPrefix(encoded, "window.queueKanbanBoardData = ") {
		t.Fatalf("expected the window.queueKanbanBoardData assignment prefix: %s", encoded)
	}
	if !strings.HasSuffix(encoded, ";\n") {
		t.Fatalf("expected the assignment to end with a semicolon + newline: %s", encoded)
	}
	if !strings.Contains(encoded, "<h2>Lessons & Notes</h2>") {
		t.Fatalf("expected pre-rendered HTML to survive verbatim (HTML escaping off): %s", encoded)
	}
	escapedLessThan := "\\u003c"
	escapedAmpersand := "\\u0026"
	if strings.Contains(encoded, escapedLessThan) || strings.Contains(encoded, escapedAmpersand) {
		t.Fatalf("body HTML was unicode-escaped by the JSON encoder: %s", encoded)
	}
}

func TestEncodeBoardMarkdownJsAssignmentRoundTripsRawSource(t *testing.T) {
	want := generatedBoardMarkdownData{
		Requests:     map[string]string{"REQ-1": "## What\n\nA <literal> & text.\n"},
		UserRequests: map[string]string{"UR-1": "# Ask\n\nCopy me.\n"},
	}
	encoded, encodeError := encodeBoardMarkdownForJsAssignment(want)
	if encodeError != nil {
		t.Fatalf("encodeBoardMarkdownForJsAssignment: %v", encodeError)
	}

	const prefix = "window.queueKanbanBoardMarkdownData = "
	if !strings.HasPrefix(encoded, prefix) || !strings.HasSuffix(encoded, ";\n") {
		t.Fatalf("unexpected lazy Markdown assignment envelope: %q", encoded)
	}
	jsonText := strings.TrimSuffix(strings.TrimPrefix(encoded, prefix), ";\n")
	var got generatedBoardMarkdownData
	if decodeError := json.Unmarshal([]byte(jsonText), &got); decodeError != nil {
		t.Fatalf("decode lazy Markdown assignment: %v", decodeError)
	}
	if got.Requests["REQ-1"] != want.Requests["REQ-1"] || got.UserRequests["UR-1"] != want.UserRequests["UR-1"] {
		t.Fatalf("raw Markdown did not round-trip: got %#v, want %#v", got, want)
	}
}

// TestRecentlyDoneWindowDefaultsTo24h asserts that a fresh board load defaults
// the RECENTLY DONE column to the 24h window: the 24h toggle button must carry
// aria-pressed="true" and the 7d (168h) button must NOT be the default-active one.
// The assertion also verifies that the inlined board.js initialises windowHours to
// 24, not 168, so the JS runtime agrees with the HTML button state on load.
func TestRecentlyDoneWindowDefaultsTo24h(t *testing.T) {
	indexHtml := generateLiveSite(t)

	// The 24h button must be the active one on load.
	activeMarker24h := `data-window-hours="24" aria-pressed="true"`
	if !strings.Contains(indexHtml, activeMarker24h) {
		t.Fatalf("24h window button is not the default-active toggle: expected %q in the generated page", activeMarker24h)
	}

	// The 7d button must NOT carry aria-pressed="true" (it is the old default).
	staleActive7d := `data-window-hours="168" aria-pressed="true"`
	if strings.Contains(indexHtml, staleActive7d) {
		t.Fatalf("7d window button is still marked as the default-active toggle: %q must not appear in the generated page", staleActive7d)
	}

	// The inlined board.js JS default must match the HTML button state.
	jsDefaultWindow24h := "windowHours: 24"
	if !strings.Contains(indexHtml, jsDefaultWindow24h) {
		t.Fatalf("board.js windowHours default is not 24: expected %q in the inlined script", jsDefaultWindow24h)
	}
	jsDefaultWindow168 := "windowHours: 168"
	if strings.Contains(indexHtml, jsDefaultWindow168) {
		t.Fatalf("board.js still initialises windowHours to 168: %q must not appear in the inlined script", jsDefaultWindow168)
	}
}
