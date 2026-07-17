package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// Testing-view write API — the ONE part of the board that writes to the tree,
// and it writes only the testing placeholders (testing.go): the testing_*
// frontmatter fields of a single REQ file, and the do-work/testers.md profile
// store. The work pipeline's fields (status, claimed_at, …) are never touched.
//
// There is deliberately no locking and no concurrency control: every write
// lands in the user's working tree, where git is the audit trail and the
// rollback mechanism. Serve mode is loopback-bound by default; the two
// browser-facing guards below (JSON content type + same-origin check) exist to
// stop a hostile web page from firing cross-site writes at localhost, not to
// implement authentication.

// testingApiMaxBodyBytes bounds a testing API request body. Feedback is capped
// separately (testingApiMaxFeedbackChars); this is the transport-level ceiling.
const testingApiMaxBodyBytes = 64 * 1024

// testingApiMaxFeedbackChars bounds the returned-with-feedback note. The
// feedback lives on one frontmatter line (escaped), so "a review document"
// belongs in the REQ body or a linked file, not here.
const testingApiMaxFeedbackChars = 10000

// testingProfileApiRequest is the POST /api/testing/profile payload.
type testingProfileApiRequest struct {
	ProfileName string `json:"name"`
}

// testingStatusApiRequest is the POST /api/testing/status payload.
// TestingStatus accepts the canonical vocabulary plus "clear" (remove the
// testing track entirely).
type testingStatusApiRequest struct {
	RequestId     string `json:"requestId"`
	TestingStatus string `json:"testingStatus"`
	TestedBy      string `json:"testedBy"`
	Feedback      string `json:"feedback"`
}

// testingApiResponse is the uniform JSON response envelope for both endpoints.
type testingApiResponse struct {
	Ok               bool     `json:"ok"`
	Error            string   `json:"error,omitempty"`
	Profiles         []string `json:"profiles,omitempty"`
	RequestId        string   `json:"requestId,omitempty"`
	TestingStatus    string   `json:"testingStatus,omitempty"`
	TestedBy         string   `json:"testedBy,omitempty"`
	TestingUpdatedAt string   `json:"testingUpdatedAt,omitempty"`
}

// writeTestingApiJson emits one testingApiResponse with the given HTTP status.
func writeTestingApiJson(responseWriter http.ResponseWriter, httpStatus int, apiResponse testingApiResponse) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	responseWriter.Header().Set("Cache-Control", "no-store")
	responseWriter.WriteHeader(httpStatus)
	_ = json.NewEncoder(responseWriter).Encode(apiResponse)
}

// writeTestingApiError is the error-shape shorthand.
func writeTestingApiError(responseWriter http.ResponseWriter, httpStatus int, errorText string) {
	writeTestingApiJson(responseWriter, httpStatus, testingApiResponse{Ok: false, Error: errorText})
}

// guardTestingApiWrite enforces the browser-facing write guards shared by both
// endpoints: POST only, a JSON content type (a cross-origin page cannot send
// application/json without a CORS preflight, which this server never grants),
// and — when the browser attaches an Origin header — an origin whose host
// matches the Host the request arrived on. Returns false after writing the
// error response when the request must be rejected.
func guardTestingApiWrite(responseWriter http.ResponseWriter, httpRequest *http.Request) bool {
	if httpRequest.Method != http.MethodPost {
		writeTestingApiError(responseWriter, http.StatusMethodNotAllowed, "POST required")
		return false
	}
	contentType := httpRequest.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(contentType)), "application/json") {
		writeTestingApiError(responseWriter, http.StatusUnsupportedMediaType, "Content-Type application/json required")
		return false
	}
	originHeader := httpRequest.Header.Get("Origin")
	if originHeader != "" && originHeader != "null" {
		parsedOrigin, parseError := url.Parse(originHeader)
		if parseError != nil || parsedOrigin.Host != httpRequest.Host {
			writeTestingApiError(responseWriter, http.StatusForbidden, "cross-origin writes are not allowed")
			return false
		}
	}
	return true
}

// decodeTestingApiBody decodes a size-capped JSON request body into target,
// writing the error response itself on failure.
func decodeTestingApiBody(responseWriter http.ResponseWriter, httpRequest *http.Request, target any) bool {
	bodyReader := http.MaxBytesReader(responseWriter, httpRequest.Body, testingApiMaxBodyBytes)
	if decodeError := json.NewDecoder(bodyReader).Decode(target); decodeError != nil {
		writeTestingApiError(responseWriter, http.StatusBadRequest, "invalid JSON body: "+decodeError.Error())
		return false
	}
	return true
}

// serveTestingProfileApi handles POST /api/testing/profile — add (or re-add,
// idempotently) a tester profile to do-work/testers.md and return the updated
// profile list.
func (liveServer *liveBoardServer) serveTestingProfileApi(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if !guardTestingApiWrite(responseWriter, httpRequest) {
		return
	}
	var profileRequest testingProfileApiRequest
	if !decodeTestingApiBody(responseWriter, httpRequest, &profileRequest) {
		return
	}

	testersFilePath := filepath.Join(liveServer.repoRoot, "do-work", testersFileRelativePath)
	updatedProfiles, appendError := appendTestingProfile(testersFilePath, profileRequest.ProfileName)
	if appendError != nil {
		writeTestingApiError(responseWriter, http.StatusBadRequest, appendError.Error())
		return
	}
	writeTestingApiJson(responseWriter, http.StatusOK, testingApiResponse{Ok: true, Profiles: updatedProfiles})
}

// serveTestingStatusApi handles POST /api/testing/status — write one REQ's
// testing placeholders. The REQ id is resolved against a fresh board build, so
// the file path written to always comes from the parsed tree, never from
// client input.
func (liveServer *liveBoardServer) serveTestingStatusApi(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	if !guardTestingApiWrite(responseWriter, httpRequest) {
		return
	}
	var statusRequest testingStatusApiRequest
	if !decodeTestingApiBody(responseWriter, httpRequest, &statusRequest) {
		return
	}

	requestedState := normalizeTestingStatus(statusRequest.TestingStatus)
	if requestedState != testingClearState && !isKnownTestingStatus(requestedState) {
		writeTestingApiError(responseWriter, http.StatusBadRequest, fmt.Sprintf(
			"testingStatus %q not recognized — expected in-testing, tested, returned, or clear", statusRequest.TestingStatus))
		return
	}

	testedBy := ""
	if requestedState != testingClearState {
		validatedTester, testerError := validateTesterProfileName(statusRequest.TestedBy)
		if testerError != nil {
			writeTestingApiError(responseWriter, http.StatusBadRequest, "testedBy: "+testerError.Error())
			return
		}
		testedBy = validatedTester
	}

	feedbackText := strings.TrimSpace(statusRequest.Feedback)
	if requestedState == testingStatusReturned && feedbackText == "" {
		writeTestingApiError(responseWriter, http.StatusBadRequest, "returned requires non-empty feedback")
		return
	}
	if len(feedbackText) > testingApiMaxFeedbackChars {
		writeTestingApiError(responseWriter, http.StatusBadRequest, fmt.Sprintf(
			"feedback is longer than %d characters — put the full write-up in the REQ body or a linked file", testingApiMaxFeedbackChars))
		return
	}

	// Resolve the REQ id → file path from a fresh parse of the tree (the same
	// build path every read uses), so a stale browser tab or a mistyped id can
	// never write outside the do-work tree.
	board, buildError := buildBoard(liveServer.repoRoot, time.Now(), liveServer.recentWindow, lookupGitCommitDate)
	if buildError != nil {
		log.Printf("queue-kanban serve: building board for testing update: %v", buildError)
		writeTestingApiError(responseWriter, http.StatusInternalServerError, "could not read the do-work tree")
		return
	}
	ticket := board.RequestsById[strings.TrimSpace(statusRequest.RequestId)]
	if ticket == nil {
		writeTestingApiError(responseWriter, http.StatusNotFound, fmt.Sprintf(
			"%q is not a REQ in the current do-work tree — reload the board", statusRequest.RequestId))
		return
	}

	// Only finished work enters testing. A non-clear transition needs the REQ
	// to be terminally successful — or to already carry a testing record (a
	// returned REQ that was requeued for the fix may legitimately restart
	// testing). Without this, a stale browser tab — or a direct API call —
	// could stamp testing state onto a pending/claimed REQ that hasn't been
	// built yet. `clear` stays allowed regardless: it only ever removes.
	hasExistingTestingRecord := ticket.TestingStatus != "" || ticket.TestingStatusUnrecognized
	if requestedState != testingClearState && !isCompletedStatus(ticket.Status) && !hasExistingTestingRecord {
		writeTestingApiError(responseWriter, http.StatusConflict, fmt.Sprintf(
			"%s has status %q — only finished REQs (completed / completed-with-issues) can enter testing; reload the board",
			ticket.RequestId, ticket.OriginalStatus))
		return
	}

	// The file path came from the parsed tree, but the tree itself is
	// untrusted checkout content — refuse symlinked targets (testing.go).
	if targetError := validateTestingWriteTarget(liveServer.repoRoot, ticket.FilePath); targetError != nil {
		log.Printf("queue-kanban serve: rejecting testing write for %s: %v", ticket.RequestId, targetError)
		writeTestingApiError(responseWriter, http.StatusBadRequest, targetError.Error())
		return
	}

	updateInstant := time.Now()
	fieldUpdates := buildTestingFieldUpdates(requestedState, testedBy, feedbackText, updateInstant)
	if upsertError := upsertFrontmatterFields(ticket.FilePath, fieldUpdates); upsertError != nil {
		log.Printf("queue-kanban serve: updating testing placeholders for %s: %v", ticket.RequestId, upsertError)
		writeTestingApiError(responseWriter, http.StatusInternalServerError, "could not update the REQ file: "+upsertError.Error())
		return
	}

	apiResponse := testingApiResponse{Ok: true, RequestId: ticket.RequestId}
	if requestedState != testingClearState {
		apiResponse.TestingStatus = requestedState
		apiResponse.TestedBy = testedBy
		apiResponse.TestingUpdatedAt = updateInstant.UTC().Format(time.RFC3339)
	}
	writeTestingApiJson(responseWriter, http.StatusOK, apiResponse)
}
