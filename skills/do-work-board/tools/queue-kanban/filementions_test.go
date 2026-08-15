package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// stubGitLookupNever is a gitCommitDateLookup that never resolves — file-mention
// tests do not care about completion dating.
func stubGitLookupNever(string, string) (time.Time, bool) {
	return time.Time{}, false
}

// writeFixtureRepoFile creates a file (and its parent dirs) under repoRoot.
func writeFixtureRepoFile(t *testing.T, repoRoot string, relativePath string, content string) {
	t.Helper()
	absolutePath := filepath.Join(repoRoot, relativePath)
	if mkdirErr := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirErr != nil {
		t.Fatalf("mkdir for %s: %v", relativePath, mkdirErr)
	}
	if writeErr := os.WriteFile(absolutePath, []byte(content), 0o644); writeErr != nil {
		t.Fatalf("write %s: %v", relativePath, writeErr)
	}
}

// TestCollectRepoFileMentionsClassifiesExistence verifies the build-time
// existence map behind the drawer's file links: a mentioned path that exists
// maps to true, a mentioned path that does not exist maps to false, and a
// placeholder mention with template tokens never enters the map at all.
func TestCollectRepoFileMentionsClassifiesExistence(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	writeFixtureRepoFile(t, repoRoot, "docs/example-guide.md", "# guide\n")

	reqBody := "---\nid: REQ-0031\ntitle: Mentions\nstatus: pending\n---\n\n" +
		"Read `docs/example-guide.md` first; `docs/never-written.md` is aspirational.\n" +
		"Placeholder `<store>/.injection-receipts/<YYYY-MM>_receipts.json` must stay plain.\n"
	writeFixtureRepoFile(t, repoRoot, "do-work/queue/REQ-0031-mentions.md", reqBody)

	board, buildErr := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, stubGitLookupNever)
	if buildErr != nil {
		t.Fatalf("buildBoard: %v", buildErr)
	}
	boardData, projectErr := buildGeneratedBoardData(board)
	if projectErr != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectErr)
	}

	if exists, mapped := boardData.RepoFileMentions["docs/example-guide.md"]; !mapped || !exists {
		t.Errorf("docs/example-guide.md: want mapped true, got mapped=%v exists=%v", mapped, exists)
	}
	if exists, mapped := boardData.RepoFileMentions["docs/never-written.md"]; !mapped || exists {
		t.Errorf("docs/never-written.md: want mapped false, got mapped=%v exists=%v", mapped, exists)
	}
	for mentionPath := range boardData.RepoFileMentions {
		if strings.ContainsAny(mentionPath, "<>") {
			t.Errorf("placeholder path leaked into the mention map: %q", mentionPath)
		}
	}
	// Static board data must never claim the live file endpoint exists.
	if boardData.LiveFileApi {
		t.Errorf("static board data must not set liveFileApi")
	}
}

// TestServeFileEndpointServesRepoFileReadOnly exercises GET /file end-to-end
// over a real loopback listener: the served board data advertises liveFileApi,
// and a repo-relative text file comes back verbatim as text/plain.
func TestServeFileEndpointServesRepoFileReadOnly(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	fileContent := "# guide\n\nline two\n"
	writeFixtureRepoFile(t, repoRoot, "docs/example-guide.md", fileContent)

	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	boardData := fetchServedBoardData(t, testServer.URL)
	if !boardData.LiveFileApi {
		t.Errorf("served board data must set liveFileApi")
	}

	resp, httpErr := http.Get(testServer.URL + "/file?path=docs%2Fexample-guide.md")
	if httpErr != nil {
		t.Fatalf("GET /file: %v", httpErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /file status = %d, want 200", resp.StatusCode)
	}
	if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
		t.Errorf("GET /file Content-Type = %q, want text/plain for non-PNG content", contentType)
	}
	if bodyText := readTestResponseBody(t, resp); bodyText != fileContent {
		t.Errorf("GET /file body = %q, want %q", bodyText, fileContent)
	}
}

// TestServeFileEndpointRedirectsHtmlToFolderPreview verifies that active HTML
// never executes on the live board origin. The target is an ephemeral preview
// origin rooted at the HTML file's containing folder, which is what lets local
// relative assets resolve without exposing the board's testing APIs.
func TestServeFileEndpointRedirectsHtmlToFolderPreview(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	writeFixtureRepoFile(t, repoRoot, "reports/example/index.html", "<!doctype html><title>Example</title>")

	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	defer shutdownHtmlPreviewManagerForTest(t, liveServer.htmlPreviews)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	noRedirectClient := testServer.Client()
	noRedirectClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, httpErr := noRedirectClient.Get(testServer.URL + "/file?path=reports%2Fexample%2Findex.html")
	if httpErr != nil {
		t.Fatalf("GET HTML /file: %v", httpErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("GET HTML /file status = %d, want 307", resp.StatusCode)
	}
	previewUrl, parseErr := url.Parse(resp.Header.Get("Location"))
	if parseErr != nil {
		t.Fatalf("parse HTML preview Location: %v", parseErr)
	}
	boardUrl, parseErr := url.Parse(testServer.URL)
	if parseErr != nil {
		t.Fatalf("parse board URL: %v", parseErr)
	}
	if previewUrl.Scheme != "http" || previewUrl.Host == "" {
		t.Errorf("HTML preview Location = %q, want an absolute HTTP URL", previewUrl.String())
	}
	if previewUrl.Host == boardUrl.Host {
		t.Errorf("HTML preview shares board origin %q", boardUrl.Host)
	}

	writeFixtureRepoFile(t, repoRoot, "reports/example/not-html.txt", "<!doctype html><title>Still text</title>")
	textResponse, textError := testServer.Client().Get(testServer.URL + "/file?path=reports%2Fexample%2Fnot-html.txt")
	if textError != nil {
		t.Fatalf("GET HTML-like text /file: %v", textError)
	}
	defer textResponse.Body.Close()
	if contentType := textResponse.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
		t.Errorf("HTML-like .txt Content-Type = %q, want inert text/plain", contentType)
	}
}

// TestServeFileEndpointRendersPng verifies that a captured screenshot opened
// through the live file endpoint carries an image content type instead of
// exposing its binary bytes as page text.
func TestServeFileEndpointRendersPng(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	imagePath := filepath.Join(repoRoot, "do-work/user-requests/UR-001/assets/REQ-001-screenshot.png")
	if mkdirErr := os.MkdirAll(filepath.Dir(imagePath), 0o755); mkdirErr != nil {
		t.Fatalf("mkdir for PNG: %v", mkdirErr)
	}
	var imageBytes bytes.Buffer
	if encodeErr := png.Encode(&imageBytes, image.NewNRGBA(image.Rect(0, 0, 1, 1))); encodeErr != nil {
		t.Fatalf("encode PNG fixture: %v", encodeErr)
	}
	if writeErr := os.WriteFile(imagePath, imageBytes.Bytes(), 0o644); writeErr != nil {
		t.Fatalf("write PNG fixture: %v", writeErr)
	}

	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	resp, httpErr := http.Get(testServer.URL +
		"/file?path=do-work%2Fuser-requests%2FUR-001%2Fassets%2FREQ-001-screenshot.png")
	if httpErr != nil {
		t.Fatalf("GET PNG /file: %v", httpErr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET PNG /file status = %d, want 200", resp.StatusCode)
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "image/png" {
		t.Errorf("GET PNG /file Content-Type = %q, want image/png", contentType)
	}
	if responseBytes := []byte(readTestResponseBody(t, resp)); !bytes.Equal(responseBytes, imageBytes.Bytes()) {
		t.Errorf("GET PNG /file body changed: got %d bytes, want %d", len(responseBytes), imageBytes.Len())
	}

	writeFixtureRepoFile(t, repoRoot, "docs/not-an-image.png", "<svg onload='alert(1)'></svg>")
	spoofedResponse, spoofedError := http.Get(testServer.URL + "/file?path=docs%2Fnot-an-image.png")
	if spoofedError != nil {
		t.Fatalf("GET mislabeled PNG /file: %v", spoofedError)
	}
	defer spoofedResponse.Body.Close()
	if contentType := spoofedResponse.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
		t.Errorf("GET mislabeled PNG /file Content-Type = %q, want inert text/plain", contentType)
	}
}

// TestProductionServeFileEndpointUsesListenerAuthority exercises /file through
// the post-bind production wrapper: the listener's assigned :0 port works,
// while a matching hostile Host/Origin pair with that same port is rejected
// before the repo-file route can run.
func TestProductionServeFileEndpointUsesListenerAuthority(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	writeFixtureRepoFile(t, repoRoot, "docs/example-guide.md", "# guide\n")
	testServer := productionLiveBoardTestServer(t, repoRoot, "127.0.0.1:0")

	acceptedResponse, acceptedError := http.Get(testServer.URL + "/file?path=docs%2Fexample-guide.md")
	if acceptedError != nil {
		t.Fatalf("accepted production GET /file: %v", acceptedError)
	}
	acceptedResponse.Body.Close()
	if acceptedResponse.StatusCode != http.StatusOK {
		t.Fatalf("accepted production GET /file status = %d, want 200", acceptedResponse.StatusCode)
	}

	listenerPort := testServer.Listener.Addr().(*net.TCPAddr).Port
	hostileAuthority := fmt.Sprintf("example.com:%d", listenerPort)
	hostileRequest, requestError := http.NewRequest(http.MethodGet,
		testServer.URL+"/file?path=docs%2Fexample-guide.md", nil)
	if requestError != nil {
		t.Fatalf("build hostile /file request: %v", requestError)
	}
	hostileRequest.Host = hostileAuthority
	hostileRequest.Header.Set("Origin", "http://"+hostileAuthority)
	hostileResponse, hostileError := testServer.Client().Do(hostileRequest)
	if hostileError != nil {
		t.Fatalf("hostile production GET /file: %v", hostileError)
	}
	hostileResponse.Body.Close()
	if hostileResponse.StatusCode != http.StatusForbidden {
		t.Errorf("hostile production GET /file status = %d, want 403", hostileResponse.StatusCode)
	}
}

// TestServeFileEndpointRejectsEscapesAndMissing locks the /file guards:
// absolute paths and ".." traversal are 400, a path not present is 404, and a
// symlink pointing outside the repo is refused even though its own path looks
// repo-relative.
func TestServeFileEndpointRejectsEscapesAndMissing(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	outsideDir := t.TempDir()
	outsideSecretPath := filepath.Join(outsideDir, "secret.txt")
	if writeErr := os.WriteFile(outsideSecretPath, []byte("secret"), 0o644); writeErr != nil {
		t.Fatalf("write outside secret: %v", writeErr)
	}
	if symlinkErr := os.Symlink(outsideSecretPath, filepath.Join(repoRoot, "escape-hop.txt")); symlinkErr != nil {
		t.Fatalf("create escape symlink: %v", symlinkErr)
	}

	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)
	testServer := httptest.NewServer(liveServer)
	defer testServer.Close()

	refusedCases := []struct {
		queryPath  string
		wantStatus int
	}{
		{"../outside.txt", http.StatusBadRequest},
		{"do-work/../../outside.txt", http.StatusBadRequest},
		{"/etc/hosts", http.StatusBadRequest},
		{"docs/no-such-file.md", http.StatusNotFound},
		{"escape-hop.txt", http.StatusBadRequest},
		{"", http.StatusBadRequest},
	}
	for _, refusedCase := range refusedCases {
		requestUrl := testServer.URL + "/file?path=" + strings.ReplaceAll(refusedCase.queryPath, "/", "%2F")
		resp, httpErr := http.Get(requestUrl)
		if httpErr != nil {
			t.Fatalf("GET /file?path=%s: %v", refusedCase.queryPath, httpErr)
		}
		bodyText := readTestResponseBody(t, resp)
		resp.Body.Close()
		if resp.StatusCode != refusedCase.wantStatus {
			t.Errorf("GET /file?path=%q status = %d, want %d (body %q)",
				refusedCase.queryPath, resp.StatusCode, refusedCase.wantStatus, truncateText(bodyText, 120))
		}
	}
}

// TestServeFileEndpointIsLoopbackOnly asserts a non-loopback peer gets 403 —
// LAN-exposing the board must not also expose the whole repo as plain text.
func TestServeFileEndpointIsLoopbackOnly(t *testing.T) {
	repoRoot := createFixtureDoWorkTree(t)
	writeFixtureRepoFile(t, repoRoot, "docs/example-guide.md", "# guide\n")
	liveServer := newLiveBoardServer(repoRoot, 7*24*time.Hour)

	lanRequest := httptest.NewRequest(http.MethodGet, "/file?path=docs%2Fexample-guide.md", nil)
	lanRequest.RemoteAddr = "203.0.113.9:44321"
	recorder := httptest.NewRecorder()
	liveServer.ServeHTTP(recorder, lanRequest)
	if recorder.Code != http.StatusForbidden {
		t.Errorf("non-loopback GET /file status = %d, want 403", recorder.Code)
	}

	loopbackRequest := httptest.NewRequest(http.MethodGet, "/file?path=docs%2Fexample-guide.md", nil)
	loopbackRequest.RemoteAddr = "127.0.0.1:55555"
	loopbackRecorder := httptest.NewRecorder()
	liveServer.ServeHTTP(loopbackRecorder, loopbackRequest)
	if loopbackRecorder.Code != http.StatusOK {
		t.Errorf("loopback GET /file status = %d, want 200", loopbackRecorder.Code)
	}
}
