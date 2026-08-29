package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// fixtureReqFileContent creates minimal REQ frontmatter + body for the serve
// fixture tree. Using a numeric id with leading zeros avoids collisions with the
// real REQ ids in the live tree.
func fixtureReqFileContent(requestId string, status string) string {
	return "---\nid: " + requestId + "\ntitle: Fixture " + requestId + "\nstatus: " + status + "\n---\n\n# " + requestId + "\n"
}

// createFixtureDoWorkTree creates a minimal do-work/queue/ tree under a temp dir
// and returns the repo root (the dir containing do-work/). Two REQs are seeded:
//
//	REQ-0001-alpha.md  status: pending
//	REQ-0002-beta.md   status: claimed
func createFixtureDoWorkTree(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	queueDir := filepath.Join(tmpDir, "do-work", "queue")
	if mkdirErr := os.MkdirAll(queueDir, 0o755); mkdirErr != nil {
		t.Fatalf("mkdir fixture tree: %v", mkdirErr)
	}
	alpha := []byte(fixtureReqFileContent("REQ-0001", "pending"))
	if writeErr := os.WriteFile(filepath.Join(queueDir, "REQ-0001-alpha.md"), alpha, 0o644); writeErr != nil {
		t.Fatalf("write fixture REQ-0001: %v", writeErr)
	}
	beta := []byte(fixtureReqFileContent("REQ-0002", "claimed"))
	if writeErr := os.WriteFile(filepath.Join(queueDir, "REQ-0002-beta.md"), beta, 0o644); writeErr != nil {
		t.Fatalf("write fixture REQ-0002: %v", writeErr)
	}
	return tmpDir
}

// TestServeHandlerRootReturnsBoardHtml asserts that GET "/" against the live
// board server returns 200 and the board HTML shell. The presence of the
// <title> text "do-work queue board" (from web/template.html) confirms the
// correct template was assembled and served.
func TestServeHandlerRootReturnsBoardHtml(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	resp, httpErr := http.Get(testServer.URL + "/")
	if httpErr != nil {
		t.Fatalf("GET /: %v", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET / status = %d, want 200", resp.StatusCode)
	}

	bodyText := readTestResponseBody(t, resp)
	if !strings.Contains(bodyText, "do-work queue board") {
		t.Fatalf("GET / body missing board HTML marker; body[:300]=%q", truncateText(bodyText, 300))
	}
	// The board shell must reference board-data.js for the browser to fetch.
	if !strings.Contains(bodyText, "board-data.js") {
		t.Fatalf("GET / body missing board-data.js reference; body[:300]=%q", truncateText(bodyText, 300))
	}
	// The page must name the project (the repo-root folder) so a viewer can tell
	// which project's board they are looking at.
	expectedProjectName := deriveProjectName(repoRoot)
	if !strings.Contains(bodyText, expectedProjectName) {
		t.Fatalf("GET / body missing project name %q; body[:300]=%q", expectedProjectName, truncateText(bodyText, 300))
	}
}

func TestStaticAndLivePagesUseIdenticalAssembledClient(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	assembledClient, assembleError := assembleBoardJavaScript(embeddedWebAssets)
	if assembleError != nil {
		t.Fatalf("assembleBoardJavaScript: %v", assembleError)
	}

	board, buildError := buildBoard(repoRoot, time.Now(), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatalf("buildBoard for static page: %v", buildError)
	}
	staticDirectory := t.TempDir()
	if generateError := generateStaticSite(staticDirectory, board); generateError != nil {
		t.Fatalf("generateStaticSite: %v", generateError)
	}
	staticPageBytes, readError := os.ReadFile(filepath.Join(staticDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read static index.html: %v", readError)
	}

	liveServer := httptest.NewServer(newLiveBoardServer(repoRoot, defaultRecentWindow))
	defer liveServer.Close()
	liveResponse, requestError := http.Get(liveServer.URL + "/")
	if requestError != nil {
		t.Fatalf("GET live /: %v", requestError)
	}
	defer liveResponse.Body.Close()
	livePage := readTestResponseBody(t, liveResponse)

	assembledClientText := string(assembledClient)
	extractAssembledClient := func(surfaceName string, pageText string) string {
		t.Helper()
		if inlineScriptCount := strings.Count(pageText, "<script>"); inlineScriptCount != 1 {
			t.Fatalf("%s page has %d classic inline scripts, want exactly 1", surfaceName, inlineScriptCount)
		}
		if strings.Contains(pageText, `type="module"`) {
			t.Fatalf("%s page changed the private classic client into a module", surfaceName)
		}
		for _, fragmentPath := range boardJavaScriptFragmentPaths {
			fragmentRequest := `src="` + filepath.Base(fragmentPath) + `"`
			if strings.Contains(pageText, fragmentRequest) {
				t.Fatalf("%s page requests fragment %s instead of inlining the assembled client", surfaceName, fragmentPath)
			}
		}
		if occurrenceCount := strings.Count(pageText, assembledClientText); occurrenceCount != 1 {
			t.Fatalf("%s page contains assembled client %d times, want exactly 1", surfaceName, occurrenceCount)
		}
		clientStart := strings.Index(pageText, assembledClientText)
		return pageText[clientStart : clientStart+len(assembledClientText)]
	}
	staticClient := extractAssembledClient("static", string(staticPageBytes))
	liveClient := extractAssembledClient("live", livePage)
	if staticClient != liveClient {
		t.Fatal("static and live pages embedded different assembled client bytes")
	}
}

// TestServeMtimeCacheInvalidatesOnStatusChange verifies the mtime cache:
//
//  1. A first /board-data.js request reflects the initial fixture state
//     (REQ-0001 in the pending column).
//  2. After the fixture file is rewritten with a new status and its mtime is
//     explicitly bumped, a second /board-data.js request rebuilds the board and
//     returns REQ-0001 in the claimed column.
func TestServeMtimeCacheInvalidatesOnStatusChange(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	reqFixturePath := filepath.Join(repoRoot, "do-work", "queue", "REQ-0001-alpha.md")

	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	// First request: REQ-0001 must appear in the pending column.
	boardData1 := fetchServedBoardData(t, testServer.URL)
	if !stringSliceContains(boardData1.Columns.Pending, "REQ-0001") {
		t.Fatalf("before status change: REQ-0001 not in pending column; pending=%v", boardData1.Columns.Pending)
	}
	if stringSliceContains(boardData1.Columns.Claimed, "REQ-0001") {
		t.Fatalf("before status change: REQ-0001 unexpectedly in claimed column")
	}

	// Rewrite the fixture file: change status from pending → claimed.
	rewrittenContent := []byte(fixtureReqFileContent("REQ-0001", "claimed"))
	if writeErr := os.WriteFile(reqFixturePath, rewrittenContent, 0o644); writeErr != nil {
		t.Fatalf("rewrite fixture REQ-0001: %v", writeErr)
	}
	// Explicitly advance the mtime by 2 seconds so the cache comparison sees a
	// difference even if the OS mtime resolution would otherwise produce the same
	// second-level timestamp as the initial write.
	futureModTime := time.Now().Add(2 * time.Second)
	if chtimesErr := os.Chtimes(reqFixturePath, futureModTime, futureModTime); chtimesErr != nil {
		t.Fatalf("chtimes fixture REQ-0001: %v", chtimesErr)
	}

	// Second request: mtime changed → cache miss → re-walk → REQ-0001 in claimed.
	boardData2 := fetchServedBoardData(t, testServer.URL)
	if stringSliceContains(boardData2.Columns.Pending, "REQ-0001") {
		t.Fatalf("after status change: REQ-0001 still in pending column (cache not invalidated)")
	}
	if !stringSliceContains(boardData2.Columns.Claimed, "REQ-0001") {
		t.Fatalf("after status change: REQ-0001 not in claimed column; claimed=%v", boardData2.Columns.Claimed)
	}
}

func TestServeMtimeCacheRefreshesTimelineAgainstCurrentTime(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	initialInstant := time.Date(2026, 8, 29, 10, 0, 0, 0, time.UTC)
	refreshedInstant := initialInstant.Add(30 * time.Minute)

	pendingTicket := &RequestTicket{
		RequestId: "REQ-9001",
		Title:     "Pending fixture",
		Status:    "pending",
		CreatedAt: initialInstant.Add(-2 * time.Hour).Format(time.RFC3339),
	}
	claimedTicket := &RequestTicket{
		RequestId: "REQ-9002",
		Title:     "Claimed fixture",
		Status:    "claimed",
		CreatedAt: initialInstant.Add(-3 * time.Hour).Format(time.RFC3339),
		ClaimedAt: initialInstant.Add(-2 * time.Hour).Format(time.RFC3339),
	}
	tickets := []*RequestTicket{pendingTicket, claimedTicket}
	for sampleIndex := 0; sampleIndex < timelineProjectionMinimumSamples; sampleIndex++ {
		claimedInstant := initialInstant.Add(-48*time.Hour + time.Duration(sampleIndex)*time.Hour)
		completedInstant := claimedInstant.Add(10 * time.Minute)
		tickets = append(tickets, &RequestTicket{
			RequestId:      fmt.Sprintf("REQ-%04d", 9100+sampleIndex),
			Title:          "Completed fixture",
			Status:         "completed",
			CreatedAt:      claimedInstant.Add(-time.Hour).Format(time.RFC3339),
			ClaimedAt:      claimedInstant.Format(time.RFC3339),
			CompletedAt:    completedInstant.Format(time.RFC3339),
			CompletionTime: completedInstant,
		})
	}

	cachedBoard := &Board{GeneratedAt: initialInstant, AllRequests: tickets}
	cachedBoardData, projectionError := buildGeneratedBoardData(cachedBoard)
	if projectionError != nil {
		t.Fatalf("project initial cached board data: %v", projectionError)
	}
	discovered, enumerateError := enumerateDoWorkTree(repoRoot)
	if enumerateError != nil {
		t.Fatalf("enumerate unchanged fixture tree: %v", enumerateError)
	}
	liveServer := newLiveBoardServer(repoRoot, defaultRecentWindow)
	liveServer.currentTime = func() time.Time { return refreshedInstant }
	liveServer.cachedFileMtimes = buildTreeMtimeFingerprint(discovered)
	liveServer.cachedBoard = cachedBoard
	liveServer.cachedBoardData = &cachedBoardData

	refreshedBoardData, refreshError := liveServer.refreshBoardData()
	if refreshError != nil {
		t.Fatalf("refresh unchanged cached board data: %v", refreshError)
	}
	if refreshedBoardData.GeneratedAt != formatTimestamp(refreshedInstant) {
		t.Fatalf("generatedAt = %q, want refreshed instant %q",
			refreshedBoardData.GeneratedAt, formatTimestamp(refreshedInstant))
	}
	if refreshedBoardData.Timeline.Now != formatTimestamp(refreshedInstant) {
		t.Fatalf("timeline now = %q, want refreshed instant %q",
			refreshedBoardData.Timeline.Now, formatTimestamp(refreshedInstant))
	}

	rowsById := map[string]generatedTimelineRow{}
	for _, row := range refreshedBoardData.Timeline.Rows {
		rowsById[row.RequestId] = row
	}
	if waitMinutes := rowsById[pendingTicket.RequestId].WaitMinutes; waitMinutes != 150 {
		t.Errorf("open wait = %.1f minutes, want 150", waitMinutes)
	}
	if workMinutes := rowsById[claimedTicket.RequestId].WorkMinutes; workMinutes != 150 {
		t.Errorf("open work = %.1f minutes, want 150", workMinutes)
	}
	if refreshedBoardData.Timeline.RangeEnd != formatTimestamp(refreshedInstant) {
		t.Errorf("timeline range end = %q, want refreshed instant %q",
			refreshedBoardData.Timeline.RangeEnd, formatTimestamp(refreshedInstant))
	}

	projection := refreshedBoardData.Timeline.Projection
	if !projection.Confident || len(projection.Rows) != 1 {
		t.Fatalf("projection = confident %v with %d rows, want one confident pending row",
			projection.Confident, len(projection.Rows))
	}
	if projection.ChainStart != formatTimestamp(refreshedInstant) {
		t.Errorf("projection chain start = %q, want %q",
			projection.ChainStart, formatTimestamp(refreshedInstant))
	}
	wantProjectionEnd := refreshedInstant.Add(10 * time.Minute)
	if projection.Rows[0].StartTime != formatTimestamp(refreshedInstant) ||
		projection.Rows[0].EndTime != formatTimestamp(wantProjectionEnd) {
		t.Errorf("projected row = %q to %q, want %q to %q",
			projection.Rows[0].StartTime, projection.Rows[0].EndTime,
			formatTimestamp(refreshedInstant), formatTimestamp(wantProjectionEnd))
	}
	if projection.QueueEnd != formatTimestamp(wantProjectionEnd) {
		t.Errorf("projection queue end = %q, want %q",
			projection.QueueEnd, formatTimestamp(wantProjectionEnd))
	}
}

func TestServeLazyMarkdownEndpointReturnsExactSources(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	// "Exact sources" means the whole file, fence included: REQ-089 made the Copy
	// payload round-trip back into a valid REQ file, and the live server shares the
	// static generator's projection, so both surfaces carry the same bytes.
	markdownData := fetchServedBoardMarkdownData(t, testServer.URL)
	wantFileText := "---\nid: REQ-0001\ntitle: Fixture REQ-0001\nstatus: pending\n---\n\n# REQ-0001\n"
	if got := markdownData.Requests["REQ-0001"]; got != wantFileText {
		t.Fatalf("served REQ Markdown = %q, want %q", got, wantFileText)
	}
}

// fetchServedBoardData requests /board-data.js from baseURL and decodes the
// window.queueKanbanBoardData = {...}; JS assignment into a generatedBoardData.
func fetchServedBoardData(t *testing.T, baseURL string) generatedBoardData {
	t.Helper()
	resp, httpErr := http.Get(baseURL + "/board-data.js")
	if httpErr != nil {
		t.Fatalf("GET /board-data.js: %v", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /board-data.js status = %d, want 200", resp.StatusCode)
	}

	rawText := readTestResponseBody(t, resp)

	// Strip the JS global assignment envelope: "window.queueKanbanBoardData = " + JSON + ";\n"
	const jsPrefix = "window.queueKanbanBoardData = "
	const jsSuffix = ";\n"
	jsonText := rawText
	if strings.HasPrefix(jsonText, jsPrefix) {
		jsonText = jsonText[len(jsPrefix):]
	}
	if strings.HasSuffix(jsonText, jsSuffix) {
		jsonText = jsonText[:len(jsonText)-len(jsSuffix)]
	}

	var boardData generatedBoardData
	if jsonErr := json.Unmarshal([]byte(jsonText), &boardData); jsonErr != nil {
		t.Fatalf("decode board data JSON: %v; raw[:200]=%q", jsonErr, truncateText(rawText, 200))
	}
	return boardData
}

// fetchServedBoardMarkdownData requests the lazily loaded raw-source payload
// and decodes its JavaScript assignment envelope.
func fetchServedBoardMarkdownData(t *testing.T, baseURL string) generatedBoardMarkdownData {
	t.Helper()
	resp, httpErr := http.Get(baseURL + "/board-markdown.js")
	if httpErr != nil {
		t.Fatalf("GET /board-markdown.js: %v", httpErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /board-markdown.js status = %d, want 200", resp.StatusCode)
	}

	rawText := readTestResponseBody(t, resp)
	const jsPrefix = "window.queueKanbanBoardMarkdownData = "
	const jsSuffix = ";\n"
	jsonText := strings.TrimSuffix(strings.TrimPrefix(rawText, jsPrefix), jsSuffix)

	var markdownData generatedBoardMarkdownData
	if jsonErr := json.Unmarshal([]byte(jsonText), &markdownData); jsonErr != nil {
		t.Fatalf("decode Markdown data JSON: %v; raw[:200]=%q", jsonErr, truncateText(rawText, 200))
	}
	return markdownData
}

// readTestResponseBody drains and returns the response body as a string.
func readTestResponseBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		t.Fatalf("read response body: %v", readErr)
	}
	return string(bodyBytes)
}

// truncateText returns at most maxLen bytes of s for diagnostic messages.
func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// TestResolveServeListenAddressBindsLoopbackByDefault asserts the security
// contract from the external review: the default address and bare port numbers
// (flag or env var) bind loopback only; LAN exposure requires an explicit
// host:port (or the deliberate host-less ":port" all-interfaces syntax).
func TestResolveServeListenAddressBindsLoopbackByDefault(t *testing.T) {
	t.Setenv(kanbanServePortEnvVar, "")

	addressCases := []struct {
		flagValue   string
		wantAddress string
	}{
		{"", "127.0.0.1:8090"},           // default
		{"9000", "127.0.0.1:9000"},       // bare flag port → loopback
		{":9000", ":9000"},               // explicit all-interfaces syntax passes through
		{"0.0.0.0:9000", "0.0.0.0:9000"}, // explicit host passes through
		{"192.168.1.5:9000", "192.168.1.5:9000"},
	}
	for _, addressCase := range addressCases {
		resolvedAddress := resolveServeListenAddress(addressCase.flagValue)
		if resolvedAddress != addressCase.wantAddress {
			t.Errorf("resolveServeListenAddress(%q) = %q, want %q",
				addressCase.flagValue, resolvedAddress, addressCase.wantAddress)
		}
	}
}

// TestResolveServeListenAddressEnvVarBarePortBindsLoopback asserts the env-var
// path applies the same loopback prefix as the flag path.
func TestResolveServeListenAddressEnvVarBarePortBindsLoopback(t *testing.T) {
	t.Setenv(kanbanServePortEnvVar, "9100")
	resolvedAddress := resolveServeListenAddress("")
	if resolvedAddress != "127.0.0.1:9100" {
		t.Errorf("resolveServeListenAddress with bare env port = %q, want %q", resolvedAddress, "127.0.0.1:9100")
	}
}

// TestDescribeListenExposureWarnsOnEveryNonLoopbackBind asserts the exposure
// warning fires for every network-reachable spelling — not only the host-less
// ":port" syntax. The board serves rendered REQ bodies, so 0.0.0.0:port,
// [::]:port, a LAN IP, and a non-localhost hostname must all warn; loopback
// spellings must stay silent.
func TestDescribeListenExposureWarnsOnEveryNonLoopbackBind(t *testing.T) {
	exposureCases := []struct {
		listenAddress      string
		wantDisplayAddress string
		wantWarning        bool
	}{
		{":9000", "localhost:9000", true},                // host-less all-interfaces syntax
		{"0.0.0.0:9000", "0.0.0.0:9000", true},           // IPv4 wildcard host
		{"[::]:9000", "[::]:9000", true},                 // IPv6 wildcard host
		{"192.168.1.5:9000", "192.168.1.5:9000", true},   // explicit LAN IP
		{"myhost.local:9000", "myhost.local:9000", true}, // non-localhost hostname — may resolve anywhere
		{"127.0.0.1:8090", "127.0.0.1:8090", false},
		{"localhost:8090", "localhost:8090", false},
		{"[::1]:8090", "[::1]:8090", false},
		{"not-an-address", "not-an-address", false}, // unparseable — net.Listen rejects it before any announce
	}
	for _, exposureCase := range exposureCases {
		displayAddress, exposureWarning := describeListenExposure(exposureCase.listenAddress)
		if displayAddress != exposureCase.wantDisplayAddress {
			t.Errorf("describeListenExposure(%q) display = %q, want %q",
				exposureCase.listenAddress, displayAddress, exposureCase.wantDisplayAddress)
		}
		if (exposureWarning != "") != exposureCase.wantWarning {
			t.Errorf("describeListenExposure(%q) warning = %q, wantWarning=%v",
				exposureCase.listenAddress, exposureWarning, exposureCase.wantWarning)
		}
		if exposureCase.wantWarning && !strings.Contains(exposureWarning, "REQ bodies") {
			t.Errorf("describeListenExposure(%q) warning must name the REQ-body exposure, got %q",
				exposureCase.listenAddress, exposureWarning)
		}
	}
}

// findFreeTcpPort asks the OS for an ephemeral port by binding to 127.0.0.1:0,
// reading the assigned port back, then releasing it immediately. There is a
// theoretical reuse race between release and the caller's own bind, but this
// is the standard Go testing idiom for "give me a free port" and is not flaky
// in practice within a single test process.
func findFreeTcpPort(t *testing.T) int {
	t.Helper()
	probeListener, listenErr := net.Listen("tcp", "127.0.0.1:0")
	if listenErr != nil {
		t.Fatalf("probe for a free port: %v", listenErr)
	}
	freePort := probeListener.Addr().(*net.TCPAddr).Port
	probeListener.Close()
	return freePort
}

// TestBindServeListenerAndAnnounceSkipsOpenerOnBindFailure is the regression
// test for the false-positive banner: when the port is already held, the bind
// must fail with an error and the browser opener must never be invoked —
// nothing "opens" a URL for a server that never came up.
func TestBindServeListenerAndAnnounceSkipsOpenerOnBindFailure(t *testing.T) {
	listenAddress := fmt.Sprintf("127.0.0.1:%d", findFreeTcpPort(t))

	blockingListener, listenErr := net.Listen("tcp", listenAddress)
	if listenErr != nil {
		t.Fatalf("pre-occupy %s: %v", listenAddress, listenErr)
	}
	defer blockingListener.Close()

	openerCallCount := 0
	stubOpener := func(string) { openerCallCount++ }

	_, bindErr := bindServeListenerAndAnnounce(listenAddress, "/tmp/fixture-repo", true, stubOpener)
	if bindErr == nil {
		t.Fatalf("bindServeListenerAndAnnounce on an occupied port: got nil error, want a bind failure")
	}
	if openerCallCount != 0 {
		t.Fatalf("browser opener called %d times on bind failure, want 0", openerCallCount)
	}
}

// TestBindServeListenerAndAnnounceOpensAfterSuccessWhenRequested asserts the
// happy path: a successful bind with openAfterBind=true invokes the opener
// exactly once with the printed board URL.
func TestBindServeListenerAndAnnounceOpensAfterSuccessWhenRequested(t *testing.T) {
	listenAddress := fmt.Sprintf("127.0.0.1:%d", findFreeTcpPort(t))
	wantUrl := fmt.Sprintf("http://%s", listenAddress)

	var openedUrls []string
	stubOpener := func(url string) { openedUrls = append(openedUrls, url) }

	listener, bindErr := bindServeListenerAndAnnounce(listenAddress, "/tmp/fixture-repo", true, stubOpener)
	if bindErr != nil {
		t.Fatalf("bindServeListenerAndAnnounce: %v", bindErr)
	}
	defer listener.Close()

	if len(openedUrls) != 1 || openedUrls[0] != wantUrl {
		t.Fatalf("browser opener calls = %v, want exactly one call with %q", openedUrls, wantUrl)
	}
}

// TestBindServeListenerAndAnnounceDefaultOffLeavesOpenerUntouched asserts
// that a successful bind with openAfterBind=false — the --open flag's default
// — never touches the opener seam at all.
func TestBindServeListenerAndAnnounceDefaultOffLeavesOpenerUntouched(t *testing.T) {
	listenAddress := fmt.Sprintf("127.0.0.1:%d", findFreeTcpPort(t))

	openerCallCount := 0
	stubOpener := func(string) { openerCallCount++ }

	listener, bindErr := bindServeListenerAndAnnounce(listenAddress, "/tmp/fixture-repo", false, stubOpener)
	if bindErr != nil {
		t.Fatalf("bindServeListenerAndAnnounce: %v", bindErr)
	}
	defer listener.Close()

	if openerCallCount != 0 {
		t.Fatalf("browser opener called %d times with openAfterBind=false, want 0", openerCallCount)
	}
}

// TestServeFingerprintTracksStrayRequestFiles asserts creating or removing a
// misplaced REQ file changes the mtime fingerprint. Without strays in the
// fingerprint, the file that is already invisible as a card would also fail to
// invalidate the cache — so the warning meant to reveal it would never appear
// (and a stale warning would never clear) until an unrelated file changed.
func TestServeFingerprintTracksStrayRequestFiles(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)

	enumerateAndFingerprint := func() map[string]time.Time {
		discovered, enumerateErr := enumerateDoWorkTree(repoRoot)
		if enumerateErr != nil {
			t.Fatalf("enumerateDoWorkTree: %v", enumerateErr)
		}
		return buildTreeMtimeFingerprint(discovered)
	}
	baselineFingerprint := enumerateAndFingerprint()

	strayDir := filepath.Join(repoRoot, "do-work", "user-requests", "UR-0001")
	if mkdirErr := os.MkdirAll(strayDir, 0o755); mkdirErr != nil {
		t.Fatalf("mkdir stray dir: %v", mkdirErr)
	}
	strayPath := filepath.Join(strayDir, "REQ-0031-example.md")
	if writeErr := os.WriteFile(strayPath, []byte(fixtureReqFileContent("REQ-0031", "pending")), 0o644); writeErr != nil {
		t.Fatalf("write stray REQ: %v", writeErr)
	}

	withStrayFingerprint := enumerateAndFingerprint()
	if _, tracked := withStrayFingerprint[strayPath]; !tracked {
		t.Fatalf("fingerprint does not track stray %s", strayPath)
	}
	if treeMtimeFingerprintsEqual(baselineFingerprint, withStrayFingerprint) {
		t.Fatal("fingerprint unchanged after adding a stray REQ — its warning would stay hidden")
	}

	if removeErr := os.Remove(strayPath); removeErr != nil {
		t.Fatalf("remove stray REQ: %v", removeErr)
	}
	if treeMtimeFingerprintsEqual(withStrayFingerprint, enumerateAndFingerprint()) {
		t.Fatal("fingerprint unchanged after removing the stray REQ — its warning would go stale")
	}
}
