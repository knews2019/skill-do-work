package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestNormalizeTestingStatusAliases asserts the alias map lands every natural
// spelling on the canonical vocabulary and leaves unknowns untouched (the
// caller decides they are unrecognized).
func TestNormalizeTestingStatusAliases(t *testing.T) {
	aliasCases := []struct {
		rawValue  string
		wantValue string
	}{
		{"in-testing", "in-testing"},
		{"In_Testing", "in-testing"},
		{"testing", "in-testing"},
		{"selected-for-testing", "in-testing"},
		{"tested", "tested"},
		{"returned", "returned"},
		{"returned-with-feedback", "returned"},
		{"Returned With Feedback", "returned"},
		{"", ""},
		{"bogus", "bogus"},
	}
	for _, aliasCase := range aliasCases {
		if got := normalizeTestingStatus(aliasCase.rawValue); got != aliasCase.wantValue {
			t.Errorf("normalizeTestingStatus(%q) = %q, want %q", aliasCase.rawValue, got, aliasCase.wantValue)
		}
	}
	if isKnownTestingStatus("") {
		t.Errorf("empty testing status must not count as a known enum member")
	}
	if !isKnownTestingStatus("in-testing") || !isKnownTestingStatus("tested") || !isKnownTestingStatus("returned") {
		t.Errorf("canonical testing statuses must all be known")
	}
}

// TestParseRequestTicketReadsTestingPlaceholders asserts the testing_* fields
// round-trip through the parser, including a double-quoted multiline feedback.
func TestParseRequestTicketReadsTestingPlaceholders(t *testing.T) {
	tmpDir := t.TempDir()
	reqFilePath := filepath.Join(tmpDir, "REQ-0101-testing-fields.md")
	reqFileContent := "---\n" +
		"id: REQ-0101\n" +
		"title: Fixture\n" +
		"status: completed\n" +
		"testing_status: returned\n" +
		"tested_by: \"Alice\"\n" +
		"testing_updated_at: 2026-07-17T10:00:00Z\n" +
		"testing_feedback: \"line one\\nline two\"\n" +
		"---\n\nbody\n"
	if writeErr := os.WriteFile(reqFilePath, []byte(reqFileContent), 0o644); writeErr != nil {
		t.Fatalf("write fixture: %v", writeErr)
	}

	ticket, parseErr := parseRequestTicket(reqFilePath, "queue")
	if parseErr != nil {
		t.Fatalf("parseRequestTicket: %v", parseErr)
	}
	if ticket.TestingStatus != "returned" {
		t.Errorf("TestingStatus = %q, want returned", ticket.TestingStatus)
	}
	if ticket.TestedBy != "Alice" {
		t.Errorf("TestedBy = %q, want Alice", ticket.TestedBy)
	}
	if ticket.TestingUpdatedAt == "" {
		t.Errorf("TestingUpdatedAt is empty, want a timestamp text")
	}
	if ticket.TestingFeedback != "line one\nline two" {
		t.Errorf("TestingFeedback = %q, want the two-line text restored", ticket.TestingFeedback)
	}
	if ticket.TestingStatusUnrecognized {
		t.Errorf("TestingStatusUnrecognized = true for a canonical value")
	}
}

// TestUnrecognizedTestingStatusFlagsAndWarns asserts the never-silently-drop
// leg: an off-vocabulary testing_status renders as not-tested with the invalid
// flag, and buildBoard raises a data warning naming the REQ.
func TestUnrecognizedTestingStatusFlagsAndWarns(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirErr := os.MkdirAll(queueDir, 0o755); mkdirErr != nil {
		t.Fatalf("mkdir: %v", mkdirErr)
	}
	reqFileContent := "---\nid: REQ-0102\ntitle: Fixture\nstatus: completed\ntesting_status: half-tested\n---\nbody\n"
	if writeErr := os.WriteFile(filepath.Join(queueDir, "REQ-0102-bad.md"), []byte(reqFileContent), 0o644); writeErr != nil {
		t.Fatalf("write fixture: %v", writeErr)
	}

	board, buildErr := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildErr != nil {
		t.Fatalf("buildBoard: %v", buildErr)
	}
	ticket := board.RequestsById["REQ-0102"]
	if ticket == nil {
		t.Fatalf("REQ-0102 not parsed")
	}
	if !ticket.TestingStatusUnrecognized {
		t.Errorf("TestingStatusUnrecognized = false, want true")
	}
	if ticket.TestingStatus != "" {
		t.Errorf("TestingStatus = %q, want empty (rendered as not tested)", ticket.TestingStatus)
	}
	warningFound := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "REQ-0102") && strings.Contains(warningText, "testing_status") {
			warningFound = true
		}
	}
	if !warningFound {
		t.Errorf("no testing_status warning for REQ-0102; warnings=%v", board.Warnings)
	}
}

// TestTestingProfilesStoreRoundTrip asserts append creates the store with its
// header, dedupes case-insensitively, validates names, and that
// loadTestingProfiles reads back file order.
func TestTestingProfilesStoreRoundTrip(t *testing.T) {
	testersFilePath := filepath.Join(t.TempDir(), "testers.md")

	profiles, appendErr := appendTestingProfile(testersFilePath, "  Alice  ")
	if appendErr != nil {
		t.Fatalf("append Alice: %v", appendErr)
	}
	if len(profiles) != 1 || profiles[0] != "Alice" {
		t.Fatalf("profiles after first append = %v, want [Alice]", profiles)
	}

	storedBytes, readErr := os.ReadFile(testersFilePath)
	if readErr != nil {
		t.Fatalf("read testers file: %v", readErr)
	}
	if !strings.HasPrefix(string(storedBytes), "# Testers") {
		t.Errorf("testers file missing explanatory header; content=%q", string(storedBytes))
	}

	if _, appendErr = appendTestingProfile(testersFilePath, "Bob"); appendErr != nil {
		t.Fatalf("append Bob: %v", appendErr)
	}
	profiles, appendErr = appendTestingProfile(testersFilePath, "alice") // case-insensitive dup → no-op
	if appendErr != nil {
		t.Fatalf("append duplicate alice: %v", appendErr)
	}
	if len(profiles) != 2 || profiles[0] != "Alice" || profiles[1] != "Bob" {
		t.Fatalf("profiles after dup append = %v, want [Alice Bob]", profiles)
	}
	if got := loadTestingProfiles(testersFilePath); len(got) != 2 || got[0] != "Alice" || got[1] != "Bob" {
		t.Fatalf("loadTestingProfiles = %v, want [Alice Bob]", got)
	}

	if _, appendErr = appendTestingProfile(testersFilePath, "  "); appendErr == nil {
		t.Errorf("empty tester name must be rejected")
	}
	if _, appendErr = appendTestingProfile(testersFilePath, "bad\nname"); appendErr == nil {
		t.Errorf("control characters in a tester name must be rejected")
	}
	if _, appendErr = appendTestingProfile(testersFilePath, strings.Repeat("x", 81)); appendErr == nil {
		t.Errorf("an over-long tester name must be rejected")
	}
}

// TestUpsertFrontmatterFieldsSurgicalEdit asserts the placeholder writer
// inserts missing keys, replaces existing ones in place, removes keys with
// their continuation lines, and leaves every other line — including the body —
// byte-identical.
func TestUpsertFrontmatterFieldsSurgicalEdit(t *testing.T) {
	reqFilePath := filepath.Join(t.TempDir(), "REQ-0103-upsert.md")
	originalContent := "---\n" +
		"id: REQ-0103\n" +
		"title: \"Keep: this exact title\"\n" +
		"status: completed\n" +
		"depends_on:\n" +
		"  - REQ-0001\n" +
		"  - REQ-0002\n" +
		"# a comment that must survive\n" +
		"---\n" +
		"\n## Body heading\n\nBody text stays untouched.\n"
	if writeErr := os.WriteFile(reqFilePath, []byte(originalContent), 0o644); writeErr != nil {
		t.Fatalf("write fixture: %v", writeErr)
	}

	// Insert all four placeholders.
	insertUpdates := buildTestingFieldUpdates(testingStatusReturned, "Alice", "fix the button\nsecond line", time.Date(2026, 7, 17, 10, 0, 0, 0, time.UTC))
	if upsertErr := upsertFrontmatterFields(reqFilePath, insertUpdates); upsertErr != nil {
		t.Fatalf("insert upsert: %v", upsertErr)
	}
	afterInsert, _ := os.ReadFile(reqFilePath)
	afterInsertText := string(afterInsert)
	for _, wantLine := range []string{
		"testing_status: returned",
		"tested_by: \"Alice\"",
		"testing_updated_at: 2026-07-17T10:00:00Z",
		"testing_feedback: \"fix the button\\nsecond line\"",
		"title: \"Keep: this exact title\"",
		"# a comment that must survive",
		"  - REQ-0002",
		"Body text stays untouched.",
	} {
		if !strings.Contains(afterInsertText, wantLine) {
			t.Errorf("after insert, file missing %q; content=\n%s", wantLine, afterInsertText)
		}
	}

	// The parser must read the placeholders back (proves the YAML stays valid).
	parsedTicket, parseErr := parseRequestTicket(reqFilePath, "queue")
	if parseErr != nil {
		t.Fatalf("parse after insert: %v", parseErr)
	}
	if parsedTicket.TestingFeedback != "fix the button\nsecond line" {
		t.Errorf("feedback round-trip = %q", parsedTicket.TestingFeedback)
	}
	if len(parsedTicket.DependsOn) != 2 {
		t.Errorf("depends_on lost in upsert: %v", parsedTicket.DependsOn)
	}

	// Transition to tested: status/tester replaced in place, feedback removed.
	testedUpdates := buildTestingFieldUpdates(testingStatusTested, "Bob", "", time.Date(2026, 7, 18, 9, 0, 0, 0, time.UTC))
	if upsertErr := upsertFrontmatterFields(reqFilePath, testedUpdates); upsertErr != nil {
		t.Fatalf("tested upsert: %v", upsertErr)
	}
	afterTested, _ := os.ReadFile(reqFilePath)
	afterTestedText := string(afterTested)
	if !strings.Contains(afterTestedText, "testing_status: tested") || !strings.Contains(afterTestedText, "tested_by: \"Bob\"") {
		t.Errorf("tested transition not written; content=\n%s", afterTestedText)
	}
	if strings.Contains(afterTestedText, "testing_feedback") {
		t.Errorf("stale testing_feedback survived the tested transition; content=\n%s", afterTestedText)
	}
	if strings.Count(afterTestedText, "testing_status:") != 1 {
		t.Errorf("testing_status duplicated by the upsert; content=\n%s", afterTestedText)
	}

	// Clear removes the whole track and restores the original file exactly.
	clearUpdates := buildTestingFieldUpdates(testingClearState, "", "", time.Now())
	if upsertErr := upsertFrontmatterFields(reqFilePath, clearUpdates); upsertErr != nil {
		t.Fatalf("clear upsert: %v", upsertErr)
	}
	afterClear, _ := os.ReadFile(reqFilePath)
	if string(afterClear) != originalContent {
		t.Errorf("clear did not restore the original file.\nwant:\n%s\ngot:\n%s", originalContent, string(afterClear))
	}
}

// TestUpsertFrontmatterFieldsConsumesDuplicateKeys asserts the duplicate-key
// contract: the YAML reader's recovery keeps the LAST occurrence of a repeated
// key, so an upsert must consume every occurrence — otherwise a transition
// looks successful but reads back as the untouched last value.
func TestUpsertFrontmatterFieldsConsumesDuplicateKeys(t *testing.T) {
	reqFilePath := filepath.Join(t.TempDir(), "REQ-0105-dup.md")
	duplicatedContent := "---\n" +
		"id: REQ-0105\n" +
		"status: completed\n" +
		"testing_status: in-testing\n" +
		"tested_by: \"Old\"\n" +
		"testing_status: tested\n" + // duplicate — the reader would keep this one
		"---\nbody\n"
	if writeErr := os.WriteFile(reqFilePath, []byte(duplicatedContent), 0o644); writeErr != nil {
		t.Fatalf("write fixture: %v", writeErr)
	}

	updates := buildTestingFieldUpdates(testingStatusReturned, "Alice", "needs work", time.Date(2026, 7, 17, 10, 0, 0, 0, time.UTC))
	if upsertErr := upsertFrontmatterFields(reqFilePath, updates); upsertErr != nil {
		t.Fatalf("upsert: %v", upsertErr)
	}
	updatedBytes, _ := os.ReadFile(reqFilePath)
	updatedText := string(updatedBytes)
	if got := strings.Count(updatedText, "testing_status:"); got != 1 {
		t.Fatalf("testing_status occurrences after upsert = %d, want 1; content=\n%s", got, updatedText)
	}
	parsedTicket, parseErr := parseRequestTicket(reqFilePath, "queue")
	if parseErr != nil {
		t.Fatalf("parse after upsert: %v", parseErr)
	}
	if parsedTicket.TestingStatus != "returned" || parsedTicket.TestedBy != "Alice" {
		t.Fatalf("reader sees %q by %q, want returned by Alice", parsedTicket.TestingStatus, parsedTicket.TestedBy)
	}

	// Clear must remove every occurrence too.
	if upsertErr := upsertFrontmatterFields(reqFilePath, buildTestingFieldUpdates(testingClearState, "", "", time.Now())); upsertErr != nil {
		t.Fatalf("clear upsert: %v", upsertErr)
	}
	clearedBytes, _ := os.ReadFile(reqFilePath)
	if strings.Contains(string(clearedBytes), "testing_") {
		t.Fatalf("clear left testing keys behind; content=\n%s", string(clearedBytes))
	}
}

// TestTestingApiRejectsUnfinishedReq asserts the pipeline-status gate: a
// non-clear transition on a pending REQ (no testing record) is a 409, while a
// requeued REQ that already carries a testing record may restart testing, and
// clear stays allowed everywhere.
func TestTestingApiRejectsUnfinishedReq(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	queueDir := filepath.Join(repoRoot, "do-work", "queue")

	// REQ-0001 is pending with no testing record: in-testing must be rejected.
	statusCode, apiResponse := postTestingApiJson(t, testServerFor(t, repoRoot), "/api/testing/status",
		map[string]string{"requestId": "REQ-0001", "testingStatus": "in-testing", "testedBy": "Alice"})
	if statusCode != http.StatusConflict || apiResponse.Ok {
		t.Fatalf("in-testing on a pending REQ: status=%d ok=%v, want 409", statusCode, apiResponse.Ok)
	}

	// A requeued REQ carrying a returned record may restart testing.
	requeuedContent := "---\nid: REQ-0004\ntitle: Requeued fix\nstatus: pending\n" +
		"testing_status: returned\ntested_by: \"Alice\"\ntesting_feedback: \"broken\"\n---\nbody\n"
	if writeErr := os.WriteFile(filepath.Join(queueDir, "REQ-0004-requeued.md"), []byte(requeuedContent), 0o644); writeErr != nil {
		t.Fatalf("write requeued fixture: %v", writeErr)
	}
	serverUrl := testServerFor(t, repoRoot)
	statusCode, apiResponse = postTestingApiJson(t, serverUrl, "/api/testing/status",
		map[string]string{"requestId": "REQ-0004", "testingStatus": "in-testing", "testedBy": "Alice"})
	if statusCode != http.StatusOK || !apiResponse.Ok {
		t.Fatalf("restart on a requeued REQ with a record: status=%d response=%+v, want 200", statusCode, apiResponse)
	}

	// Clear is always allowed — it only removes.
	statusCode, apiResponse = postTestingApiJson(t, serverUrl, "/api/testing/status",
		map[string]string{"requestId": "REQ-0001", "testingStatus": "clear"})
	if statusCode != http.StatusOK || !apiResponse.Ok {
		t.Fatalf("clear on a pending REQ: status=%d response=%+v, want 200", statusCode, apiResponse)
	}
}

// TestTestingApiRejectsSymlinkedReqFile asserts the write path refuses to
// follow a REQ-*.md symlink out of the do-work tree: the API errors and the
// symlink's target stays byte-identical.
func TestTestingApiRejectsSymlinkedReqFile(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	outsideTargetPath := filepath.Join(repoRoot, "outside-target.md")
	outsideContent := "---\nid: REQ-0042\ntitle: Outside file\nstatus: completed\n---\nprecious content\n"
	if writeErr := os.WriteFile(outsideTargetPath, []byte(outsideContent), 0o644); writeErr != nil {
		t.Fatalf("write outside target: %v", writeErr)
	}
	symlinkPath := filepath.Join(repoRoot, "do-work", "queue", "REQ-0042-link.md")
	if symlinkErr := os.Symlink(outsideTargetPath, symlinkPath); symlinkErr != nil {
		t.Skipf("cannot create symlinks on this platform: %v", symlinkErr)
	}

	statusCode, apiResponse := postTestingApiJson(t, testServerFor(t, repoRoot), "/api/testing/status",
		map[string]string{"requestId": "REQ-0042", "testingStatus": "in-testing", "testedBy": "Alice"})
	if statusCode != http.StatusBadRequest || apiResponse.Ok {
		t.Fatalf("write through a symlinked REQ: status=%d ok=%v, want 400", statusCode, apiResponse.Ok)
	}
	targetBytes, _ := os.ReadFile(outsideTargetPath)
	if string(targetBytes) != outsideContent {
		t.Fatalf("symlink target was modified:\n%s", string(targetBytes))
	}
}

// testServerFor starts a live board server over repoRoot and returns its base
// URL, closing it with the test.
func testServerFor(t *testing.T, repoRoot string) string {
	t.Helper()
	testServer := httptest.NewServer(newLiveBoardServer(repoRoot, 7*24*time.Hour))
	t.Cleanup(testServer.Close)
	return testServer.URL
}

// TestUpsertFrontmatterFieldsRejectsFencelessFile asserts a file without
// frontmatter is an error, never a guessed edit.
func TestUpsertFrontmatterFieldsRejectsFencelessFile(t *testing.T) {
	reqFilePath := filepath.Join(t.TempDir(), "REQ-0104-nofm.md")
	if writeErr := os.WriteFile(reqFilePath, []byte("# Just a body\n"), 0o644); writeErr != nil {
		t.Fatalf("write fixture: %v", writeErr)
	}
	updates := buildTestingFieldUpdates(testingStatusInTesting, "Alice", "", time.Now())
	if upsertErr := upsertFrontmatterFields(reqFilePath, updates); upsertErr == nil {
		t.Errorf("upsert on a frontmatter-less file must fail")
	}
}

// postTestingApiJson posts a JSON payload to a testing API path and decodes the
// response envelope.
func postTestingApiJson(t *testing.T, baseUrl string, apiPath string, payload any) (int, testingApiResponse) {
	t.Helper()
	payloadBytes, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		t.Fatalf("marshal payload: %v", marshalErr)
	}
	httpResponse, postErr := http.Post(baseUrl+apiPath, "application/json", bytes.NewReader(payloadBytes))
	if postErr != nil {
		t.Fatalf("POST %s: %v", apiPath, postErr)
	}
	defer httpResponse.Body.Close()
	var apiResponse testingApiResponse
	if decodeErr := json.NewDecoder(httpResponse.Body).Decode(&apiResponse); decodeErr != nil {
		t.Fatalf("decode %s response: %v", apiPath, decodeErr)
	}
	return httpResponse.StatusCode, apiResponse
}

// TestTestingApiProfileAndStatusHappyPath drives the live server end to end:
// add a tester profile, mark a completed REQ in-testing, then returned with
// feedback — asserting the markdown files (the database) carry the record.
func TestTestingApiProfileAndStatusHappyPath(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	doneReqPath := filepath.Join(queueDir, "REQ-0003-done.md")
	if writeErr := os.WriteFile(doneReqPath, []byte(fixtureReqFileContent("REQ-0003", "completed")), 0o644); writeErr != nil {
		t.Fatalf("write fixture REQ-0003: %v", writeErr)
	}

	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	statusCode, profileResponse := postTestingApiJson(t, testServer.URL, "/api/testing/profile",
		map[string]string{"name": "Alice"})
	if statusCode != http.StatusOK || !profileResponse.Ok {
		t.Fatalf("profile add: status=%d response=%+v", statusCode, profileResponse)
	}
	if len(profileResponse.Profiles) != 1 || profileResponse.Profiles[0] != "Alice" {
		t.Fatalf("profiles = %v, want [Alice]", profileResponse.Profiles)
	}

	statusCode, statusResponse := postTestingApiJson(t, testServer.URL, "/api/testing/status",
		map[string]string{"requestId": "REQ-0003", "testingStatus": "in-testing", "testedBy": "Alice"})
	if statusCode != http.StatusOK || !statusResponse.Ok {
		t.Fatalf("in-testing update: status=%d response=%+v", statusCode, statusResponse)
	}

	statusCode, statusResponse = postTestingApiJson(t, testServer.URL, "/api/testing/status",
		map[string]string{"requestId": "REQ-0003", "testingStatus": "returned", "testedBy": "Alice", "feedback": "button misaligned"})
	if statusCode != http.StatusOK || !statusResponse.Ok {
		t.Fatalf("returned update: status=%d response=%+v", statusCode, statusResponse)
	}

	updatedBytes, readErr := os.ReadFile(doneReqPath)
	if readErr != nil {
		t.Fatalf("read updated REQ: %v", readErr)
	}
	updatedText := string(updatedBytes)
	for _, wantLine := range []string{"testing_status: returned", "tested_by: \"Alice\"", "testing_feedback: \"button misaligned\""} {
		if !strings.Contains(updatedText, wantLine) {
			t.Errorf("updated REQ missing %q; content=\n%s", wantLine, updatedText)
		}
	}

	// The next board build must reflect the record and carry the profile list;
	// served data must also announce the live API.
	boardData := fetchServedBoardData(t, testServer.URL)
	if !boardData.LiveTestingApi {
		t.Errorf("served board data must set liveTestingApi")
	}
	if len(boardData.TestingProfiles) != 1 || boardData.TestingProfiles[0] != "Alice" {
		t.Errorf("served testingProfiles = %v, want [Alice]", boardData.TestingProfiles)
	}
	servedRequest := boardData.Requests["REQ-0003"]
	if servedRequest.TestingStatus != "returned" || servedRequest.TestedBy != "Alice" || servedRequest.TestingFeedback != "button misaligned" {
		t.Errorf("served REQ-0003 testing fields = %+v", servedRequest)
	}
}

// TestTestingApiValidation asserts the write guards: unknown states, missing
// tester, feedback-less returns, unknown REQ ids, non-POST methods, non-JSON
// content types, and cross-origin requests are all rejected.
func TestTestingApiValidation(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	rejectionCases := []struct {
		caseName   string
		payload    map[string]string
		wantStatus int
	}{
		{"unknown state", map[string]string{"requestId": "REQ-0001", "testingStatus": "half-done", "testedBy": "A"}, http.StatusBadRequest},
		{"missing tester", map[string]string{"requestId": "REQ-0001", "testingStatus": "in-testing"}, http.StatusBadRequest},
		{"returned without feedback", map[string]string{"requestId": "REQ-0001", "testingStatus": "returned", "testedBy": "A"}, http.StatusBadRequest},
		{"unknown REQ", map[string]string{"requestId": "REQ-9999", "testingStatus": "in-testing", "testedBy": "A"}, http.StatusNotFound},
	}
	for _, rejectionCase := range rejectionCases {
		statusCode, apiResponse := postTestingApiJson(t, testServer.URL, "/api/testing/status", rejectionCase.payload)
		if statusCode != rejectionCase.wantStatus || apiResponse.Ok {
			t.Errorf("%s: status=%d ok=%v, want status=%d ok=false",
				rejectionCase.caseName, statusCode, apiResponse.Ok, rejectionCase.wantStatus)
		}
	}

	// GET must be rejected.
	getResponse, getErr := http.Get(testServer.URL + "/api/testing/status")
	if getErr != nil {
		t.Fatalf("GET status endpoint: %v", getErr)
	}
	getResponse.Body.Close()
	if getResponse.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET status = %d, want 405", getResponse.StatusCode)
	}

	// A non-JSON content type must be rejected (the CSRF guard: a cross-origin
	// page cannot send application/json without a CORS preflight).
	formResponse, formErr := http.Post(testServer.URL+"/api/testing/profile",
		"application/x-www-form-urlencoded", strings.NewReader("name=Alice"))
	if formErr != nil {
		t.Fatalf("form POST: %v", formErr)
	}
	formResponse.Body.Close()
	if formResponse.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("form POST status = %d, want 415", formResponse.StatusCode)
	}

	// A mismatched Origin header must be rejected.
	crossOriginRequest, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/testing/profile",
		strings.NewReader(`{"name":"Alice"}`))
	crossOriginRequest.Header.Set("Content-Type", "application/json")
	crossOriginRequest.Header.Set("Origin", "https://evil.example")
	crossOriginResponse, crossOriginErr := http.DefaultClient.Do(crossOriginRequest)
	if crossOriginErr != nil {
		t.Fatalf("cross-origin POST: %v", crossOriginErr)
	}
	crossOriginResponse.Body.Close()
	if crossOriginResponse.StatusCode != http.StatusForbidden {
		t.Errorf("cross-origin POST status = %d, want 403", crossOriginResponse.StatusCode)
	}

	// Nothing may have been written by any rejected request.
	if _, statErr := os.Stat(filepath.Join(repoRoot, "do-work", "testers.md")); !os.IsNotExist(statErr) {
		t.Errorf("a rejected request created do-work/testers.md")
	}
}

// TestGenerateStaticBoardDataOmitsLiveTestingApi asserts a static snapshot
// never claims the write API (the frontend renders the testing view read-only
// from that signal).
func TestGenerateStaticBoardDataOmitsLiveTestingApi(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	board, buildErr := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildErr != nil {
		t.Fatalf("buildBoard: %v", buildErr)
	}
	boardData, projectErr := buildGeneratedBoardData(board)
	if projectErr != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectErr)
	}
	if boardData.LiveTestingApi {
		t.Errorf("static board data must not set liveTestingApi")
	}
}
