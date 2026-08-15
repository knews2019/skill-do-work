package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func shutdownHtmlPreviewManagerForTest(t *testing.T, previewManager *htmlFolderPreviewManager) {
	t.Helper()
	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), time.Second)
	defer cancelShutdown()
	if shutdownError := previewManager.shutdown(shutdownContext); shutdownError != nil {
		t.Errorf("shutdown HTML previews: %v", shutdownError)
	}
}

func readPreviewResponse(t *testing.T, response *http.Response) string {
	t.Helper()
	defer response.Body.Close()
	responseBytes, readError := io.ReadAll(response.Body)
	if readError != nil {
		t.Fatalf("read preview response: %v", readError)
	}
	return string(responseBytes)
}

// TestHtmlFolderPreviewServesAuthoredFolder verifies the core browser-facing
// contract: one folder is one origin, its root-relative and relative resources
// keep their MIME types, and a different folder gets a different origin.
func TestHtmlFolderPreviewServesAuthoredFolder(t *testing.T) {
	repoRoot := t.TempDir()
	firstFolder := filepath.Join(repoRoot, "reports", "first")
	secondFolder := filepath.Join(repoRoot, "reports", "second")
	writeFixtureRepoFile(t, repoRoot, "reports/first/index.html",
		"<!doctype html><link rel='stylesheet' href='/assets/site.css'><script src='scripts/app.js'></script>")
	writeFixtureRepoFile(t, repoRoot, "reports/first/other.htm", "<!doctype html><title>Other</title>")
	writeFixtureRepoFile(t, repoRoot, "reports/first/assets/site.css", "body { color: rgb(1, 2, 3); }")
	writeFixtureRepoFile(t, repoRoot, "reports/first/scripts/app.js", "window.previewLoaded = true;")
	writeFixtureRepoFile(t, repoRoot, "reports/first/data/state.json", `{"ready":true}`)
	writeFixtureRepoFile(t, repoRoot, "reports/second/index.html", "<!doctype html><title>Second</title>")

	previewManager := newHtmlFolderPreviewManager()
	defer shutdownHtmlPreviewManagerForTest(t, previewManager)
	firstUrl, firstError := previewManager.previewUrlForHtmlFile(filepath.Join(firstFolder, "index.html"))
	if firstError != nil {
		t.Fatalf("first preview URL: %v", firstError)
	}
	otherUrl, otherError := previewManager.previewUrlForHtmlFile(filepath.Join(firstFolder, "other.htm"))
	if otherError != nil {
		t.Fatalf("same-folder preview URL: %v", otherError)
	}
	secondUrl, secondError := previewManager.previewUrlForHtmlFile(filepath.Join(secondFolder, "index.html"))
	if secondError != nil {
		t.Fatalf("second preview URL: %v", secondError)
	}

	parsePreviewUrl := func(label string, previewUrl string) *url.URL {
		t.Helper()
		parsedUrl, parseError := url.Parse(previewUrl)
		if parseError != nil {
			t.Fatalf("parse %s preview URL %q: %v", label, previewUrl, parseError)
		}
		return parsedUrl
	}
	firstParsed := parsePreviewUrl("first", firstUrl)
	otherParsed := parsePreviewUrl("same-folder", otherUrl)
	secondParsed := parsePreviewUrl("second", secondUrl)
	if firstParsed.Host != otherParsed.Host {
		t.Errorf("same folder got different preview origins: %q and %q", firstParsed.Host, otherParsed.Host)
	}
	if firstParsed.Host == secondParsed.Host {
		t.Errorf("different folders share preview origin %q", firstParsed.Host)
	}

	pageResponse, pageError := http.Get(firstUrl)
	if pageError != nil {
		t.Fatalf("GET preview HTML: %v", pageError)
	}
	pageBody := readPreviewResponse(t, pageResponse)
	if pageResponse.StatusCode != http.StatusOK || !strings.HasPrefix(pageResponse.Header.Get("Content-Type"), "text/html") {
		t.Errorf("preview HTML status/type = %d %q, want 200 text/html",
			pageResponse.StatusCode, pageResponse.Header.Get("Content-Type"))
	}
	if !strings.Contains(pageBody, "assets/site.css") {
		t.Errorf("preview HTML body changed: %q", pageBody)
	}
	if policy := pageResponse.Header.Get("Content-Security-Policy"); policy != "frame-ancestors 'none'" {
		t.Errorf("preview CSP = %q, want frame isolation without resource restrictions", policy)
	}
	if openerPolicy := pageResponse.Header.Get("Cross-Origin-Opener-Policy"); openerPolicy != "same-origin" {
		t.Errorf("preview COOP = %q, want same-origin", openerPolicy)
	}
	if corsHeader := pageResponse.Header.Get("Access-Control-Allow-Origin"); corsHeader != "" {
		t.Errorf("preview unexpectedly grants CORS with %q", corsHeader)
	}

	assetCases := []struct {
		resourceReference string
		wantTypeFragment  string
		wantBody          string
	}{
		{"/assets/site.css", "text/css", "body { color: rgb(1, 2, 3); }"},
		{"scripts/app.js", "javascript", "window.previewLoaded = true;"},
		{"data/state.json", "application/json", `{"ready":true}`},
	}
	for _, assetCase := range assetCases {
		assetUrl, resolveError := firstParsed.Parse(assetCase.resourceReference)
		if resolveError != nil {
			t.Fatalf("resolve %s: %v", assetCase.resourceReference, resolveError)
		}
		assetResponse, getError := http.Get(assetUrl.String())
		if getError != nil {
			t.Fatalf("GET %s: %v", assetCase.resourceReference, getError)
		}
		assetBody := readPreviewResponse(t, assetResponse)
		if assetResponse.StatusCode != http.StatusOK {
			t.Errorf("GET %s status = %d, want 200", assetCase.resourceReference, assetResponse.StatusCode)
		}
		if contentType := assetResponse.Header.Get("Content-Type"); !strings.Contains(contentType, assetCase.wantTypeFragment) {
			t.Errorf("GET %s Content-Type = %q, want %q", assetCase.resourceReference, contentType, assetCase.wantTypeFragment)
		}
		if assetBody != assetCase.wantBody {
			t.Errorf("GET %s body = %q, want %q", assetCase.resourceReference, assetBody, assetCase.wantBody)
		}
	}
}

// TestHtmlFolderPreviewGuards locks the read-only folder boundary. Directories
// never list, non-loopback and non-read requests fail, and neither lexical nor
// symlink traversal can escape the selected HTML folder.
func TestHtmlFolderPreviewGuards(t *testing.T) {
	previewRoot := t.TempDir()
	writeFixtureRepoFile(t, previewRoot, "index.html", "<!doctype html><title>Root</title>")
	if mkdirError := os.Mkdir(filepath.Join(previewRoot, "without-index"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir without-index: %v", mkdirError)
	}
	outsideRoot := t.TempDir()
	outsideFile := filepath.Join(outsideRoot, "secret.txt")
	if writeError := os.WriteFile(outsideFile, []byte("secret"), 0o644); writeError != nil {
		t.Fatalf("write outside file: %v", writeError)
	}
	if symlinkError := os.Symlink(outsideFile, filepath.Join(previewRoot, "escape.txt")); symlinkError != nil {
		t.Fatalf("create escape symlink: %v", symlinkError)
	}

	previewHandler := &htmlFolderPreviewHandler{directoryRoot: previewRoot}
	testCases := []struct {
		caseName   string
		method     string
		requestUrl string
		remoteAddr string
		wantStatus int
	}{
		{"folder root index", http.MethodGet, "http://preview.test/", "127.0.0.1:41001", http.StatusOK},
		{"head is allowed", http.MethodHead, "http://preview.test/index.html", "127.0.0.1:41002", http.StatusOK},
		{"directory never lists", http.MethodGet, "http://preview.test/without-index/", "127.0.0.1:41003", http.StatusNotFound},
		{"post is refused", http.MethodPost, "http://preview.test/index.html", "127.0.0.1:41004", http.StatusMethodNotAllowed},
		{"LAN peer is refused", http.MethodGet, "http://preview.test/index.html", "203.0.113.4:41005", http.StatusForbidden},
		{"lexical traversal is refused", http.MethodGet, "http://preview.test/../secret.txt", "127.0.0.1:41006", http.StatusBadRequest},
		{"symlink traversal is refused", http.MethodGet, "http://preview.test/escape.txt", "127.0.0.1:41007", http.StatusBadRequest},
	}
	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.requestUrl, nil)
			request.RemoteAddr = testCase.remoteAddr
			responseRecorder := httptest.NewRecorder()
			previewHandler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d; body=%q", responseRecorder.Code, testCase.wantStatus, responseRecorder.Body.String())
			}
		})
	}
}

func TestHtmlFolderPreviewShutdownClosesOrigins(t *testing.T) {
	previewRoot := t.TempDir()
	htmlPath := filepath.Join(previewRoot, "index.html")
	if writeError := os.WriteFile(htmlPath, []byte("<!doctype html><title>Close me</title>"), 0o644); writeError != nil {
		t.Fatalf("write HTML: %v", writeError)
	}
	previewManager := newHtmlFolderPreviewManager()
	previewUrl, previewError := previewManager.previewUrlForHtmlFile(htmlPath)
	if previewError != nil {
		t.Fatalf("preview URL: %v", previewError)
	}
	response, getError := http.Get(previewUrl)
	if getError != nil {
		t.Fatalf("GET before shutdown: %v", getError)
	}
	response.Body.Close()

	shutdownHtmlPreviewManagerForTest(t, previewManager)
	if _, restartError := previewManager.previewUrlForHtmlFile(htmlPath); restartError == nil {
		t.Error("closed preview manager accepted a new preview")
	}
	freshClient := &http.Client{
		Timeout:   250 * time.Millisecond,
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	if responseAfterShutdown, getAfterShutdownError := freshClient.Get(previewUrl); getAfterShutdownError == nil {
		responseAfterShutdown.Body.Close()
		t.Error("preview origin still accepted a new connection after shutdown")
	}
}
