package toolboxcommands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	reportGalleryPreviewWidth  = 1040
	reportGalleryPreviewHeight = 740
)

// Retain the entire catalog entry, including its extension and optional month
// folder, so different bundles or entry files never share a preview.
func reportGalleryPreviewPath(entry string) string {
	return reportGalleryFolder + "/previews/" + entry + ".png"
}

func captureReportGalleryPreviews(repositoryRoot string, repository *os.Root, args []string) resultmodel.CommandResult {
	bundles, _, err := reportIndexWalk(repositoryRoot)
	if err != nil {
		return reportGalleryFailure(reportsFolder, err)
	}
	worked, err := readWorkedReports(repository)
	if err != nil {
		return reportGalleryFailure(reportGalleryFolder, err)
	}
	curated := make(map[string]bool, len(worked))
	for _, card := range worked {
		curated[card.ReportPath] = true
	}
	var entries []string
	for _, source := range reportGalleryGeneratedSources(repository, bundles) {
		entry := source.bundle.Entry
		if !curated[entry] && (len(args) == 0 || entry == args[0]) {
			entries = append(entries, entry)
		}
	}
	if len(entries) == 0 {
		if len(args) == 1 {
			return reportGalleryFailure(args[0], fmt.Errorf("no non-curated report has primary entry %q", args[0]))
		}
		return exactOutputResult("", nil)
	}
	browser := discoverJudgeBrowser()
	if browser == "" {
		return reportGalleryFailure(reportGalleryFolder, errors.New("no Chromium browser was found for gallery preview capture"))
	}

	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return reportGalleryFailure(reportGalleryFolder, err)
	}
	server := &http.Server{Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Render a Markdown-only legacy entry as the browser's plain-text page.
		if strings.EqualFold(path.Ext(request.URL.Path), ".md") {
			writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		serveReportGalleryAsset(writer, request, repositoryRoot, repository, "")
	})}
	defer server.Close()
	go func() { _ = server.Serve(listener) }()
	origin := "http://" + listener.Addr().String()

	session, err := startJudgeBrowser(browser)
	if err != nil {
		return reportGalleryFailure(reportGalleryFolder, err)
	}
	defer session.stop()
	if err := session.attachToPage(); err != nil {
		return reportGalleryFailure(reportGalleryFolder, err)
	}
	var output strings.Builder
	var changes []resultmodel.RecordedChange
	var findings []resultmodel.CommandFinding
	for _, entry := range entries {
		pageURL := origin + (&url.URL{Path: "/" + entry}).EscapedPath()
		change, err := session.captureReportGalleryPreview(repositoryRoot, repository, entry, pageURL)
		if err != nil {
			findings = append(findings, reportGalleryFailure(entry, fmt.Errorf("capture preview of %s: %w", entry, err)).Findings...)
			continue
		}
		fmt.Fprintln(&output, change.Path)
		changes = append(changes, change)
	}
	result := exactOutputResult(output.String(), changes)
	if len(findings) > 0 {
		result.Outcome = resultmodel.OutcomeFailure
		result.Findings = findings
	}
	return result
}

func (session *judgeBrowserSession) captureReportGalleryPreview(repositoryRoot string, repository *os.Root, entry, pageURL string) (resultmodel.RecordedChange, error) {
	var change resultmodel.RecordedChange
	if err := session.loadCapturePage(pageURL, reportGalleryPreviewWidth, reportGalleryPreviewHeight, "light", true); err != nil {
		return change, err
	}
	pageJSON, err := session.evaluate(`(async () => {
		await document.fonts.ready;
		window.scrollTo({top: 0, left: 0, behavior: 'instant'});
		await new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve)));
		return {
			url: location.href,
			scrollY: window.scrollY,
			status: performance.getEntriesByType('navigation')[0]?.responseStatus
		};
	})()`)
	if err != nil {
		return change, err
	}
	var rendered struct {
		URL     string  `json:"url"`
		ScrollY float64 `json:"scrollY"`
		Status  int     `json:"status"`
	}
	if err := json.Unmarshal(pageJSON, &rendered); err != nil {
		return change, fmt.Errorf("decode the preview page: %w", err)
	}
	if rendered.URL != pageURL || rendered.ScrollY != 0 {
		return change, fmt.Errorf("expected the top of %s, got %s at scroll offset %g", pageURL, rendered.URL, rendered.ScrollY)
	}
	if rendered.Status >= http.StatusBadRequest {
		return change, fmt.Errorf("load %s: HTTP %d", pageURL, rendered.Status)
	}
	image, err := session.capturePagePNG(reportGalleryPreviewWidth, reportGalleryPreviewHeight)
	if err != nil {
		return change, err
	}
	preview := reportGalleryPreviewPath(entry)
	if err := rootedMkdirAll(repositoryRoot, path.Dir(preview), 0o755); err != nil {
		return change, err
	}
	_, err = repository.Lstat(preview)
	exists := err == nil
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return change, err
	}
	if err := rootedPublishFile(repositoryRoot, preview, image, 0o644, exists); err != nil {
		return change, err
	}
	kind := "created"
	if exists {
		kind = "modified"
	}
	return resultmodel.RecordedChange{Path: preview, Kind: kind, Detail: "browser screenshot of " + entry}, nil
}
