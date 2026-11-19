package toolboxcommands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAIReportGalleryPreviewRejectsHTTPErrorPage(t *testing.T) {
	requireJudgeBrowser(t)
	const entry = "ai-reports/moved/index.html"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Security-Policy", "sandbox allow-scripts")
		http.NotFound(writer, request)
	}))
	t.Cleanup(server.Close)

	root := t.TempDir()
	repository, err := os.OpenRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { repository.Close() })
	session, err := startJudgeBrowser(discoverJudgeBrowser())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(session.stop)
	if err := session.attachToPage(); err != nil {
		t.Fatal(err)
	}

	_, err = session.captureReportGalleryPreview(root, repository, entry, server.URL+"/"+entry)
	if err == nil || !strings.Contains(err.Error(), "HTTP 404") {
		t.Fatalf("expected a 404 diagnostic, got %v", err)
	}
	if _, err := repository.Stat(reportGalleryPreviewPath(entry)); !os.IsNotExist(err) {
		t.Fatalf("a failed capture must not write a preview, got %v", err)
	}
}
