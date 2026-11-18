package toolboxcommands

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const galleryTestRoot = "DO_WORK_REPORT_GALLERY_TEST_ROOT"

func TestReportGalleryFigureIDsFollowFiguresAcrossReordering(t *testing.T) {
	const diagram = `<figure><svg role="img" aria-labelledby="split-title"><title id="split-title">Procedure split</title></svg><figcaption>Measured steps</figcaption></figure>`
	const chart = `<svg role="img" aria-label="Gate time"><path d="M0 2"></path></svg>`
	const image = `<IMG src="screenshots/replay.png" alt="Replay" />`
	const hidden = `<!-- <figure id="comment"></figure> --><script>const demo = '<figure id="script"></figure>';</script><svg aria-hidden="true"><title>Decoration</title></svg>`
	const href = "../../ai-reports/phase-by-phase/index.html"
	original := reportGalleryFigures([]byte(diagram+hidden+chart+image), href)
	reordered := reportGalleryFigures([]byte(image+chart+diagram), href)
	if len(original) != 3 || len(reordered) != 3 || original[0].ID != "split-title" ||
		original[0].Label != "Procedure split" || original[0] != reordered[2] ||
		original[1] != reordered[1] || original[2] != reordered[0] {
		t.Fatalf("unstable figure identity: %v reordered to %v", original, reordered)
	}
}

func TestAIReportGalleryServeRoundTrip(t *testing.T) {
	if root := os.Getenv(galleryTestRoot); root != "" {
		result := runGalleryFixture(t, root, "--serve", "--no-open")
		if result.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("live gallery exited: %+v", result.Findings)
		}
		return
	}

	root := t.TempDir()
	const report = "ai-reports/2026-11_1-report/2026-11-18_0900_new-proposal/index.html"
	const otherReport = "ai-reports/2026-11-17_0900_another-report/index.html"
	writeGalleryFixture(t, root, report, `<title>New proposal</title><figure id="procedure-split"><figcaption>Procedure split</figcaption></figure>`)
	writeGalleryFixture(t, root, otherReport, `<title>Another report</title><figure id="procedure-split"><figcaption>Other visual</figcaption></figure>`)
	writeGalleryFixture(t, root, ReportGalleryPageFile, "Do not overwrite this hand-made page.")
	writeGalleryFixture(t, root, "private.txt", "Outside the report and gallery directories.")
	getPage, stop := startGalleryFixture(t, root)
	webURL, err := url.Parse(getPage)
	if err != nil {
		t.Fatal(err)
	}
	origin := webURL.Scheme + "://" + webURL.Host
	endpoint := origin + reportGalleryFeedbackAPI
	response, err := http.Get(getPage)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	page, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{report, "Procedure split", `data-target-kind="report"`, `data-figure-id="procedure-split"`} {
		if response.StatusCode != http.StatusOK || !strings.Contains(string(page), expected) {
			t.Fatalf("live gallery status %d, missing %q", response.StatusCode, expected)
		}
	}
	whole := reportGalleryEntry{reportGalleryTarget{Kind: "report", ReportPath: report}, "like", ""}
	figure := reportGalleryEntry{reportGalleryTarget{Kind: "figure", ReportPath: report, FigureID: "procedure-split"}, "dislike", "Label the arrows.\nKeep the title café."}
	sibling := reportGalleryEntry{reportGalleryTarget{Kind: "figure", ReportPath: otherReport, FigureID: "procedure-split"}, "like", "Different report"}
	store := filepath.Join(root, reportGalleryFeedbackFile)
	runGalleryFeedbackRequest(t, endpoint, "", http.MethodGet, nil, http.StatusOK)
	unknown := figure
	unknown.FigureID = "not-on-this-page"
	runGalleryFeedbackRequest(t, endpoint, origin, http.MethodPut, unknown, http.StatusBadRequest)
	runGalleryFeedbackRequest(t, endpoint, "https://other.example", http.MethodPut, whole, http.StatusForbidden)
	if _, err := os.Lstat(store); !os.IsNotExist(err) {
		t.Fatalf("load or rejected save touched the record: %v", err)
	}
	for _, entry := range []reportGalleryEntry{whole, figure, sibling} {
		acknowledgement := runGalleryFeedbackRequest(t, endpoint, origin, http.MethodPut, entry, http.StatusOK)
		var actual reportGalleryEntry
		if err := json.Unmarshal(acknowledgement, &actual); err != nil || actual != entry {
			t.Fatalf("save acknowledged %s, want %+v: %v", acknowledgement, entry, err)
		}
	}
	figure.Note = "Retain the caption, revise the diagram."
	runGalleryFeedbackRequest(t, endpoint, origin, http.MethodPut, figure, http.StatusOK)
	want := reportGalleryFeedback{1, []reportGalleryEntry{whole, figure, sibling}}
	saved, err := os.ReadFile(store)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(saved, []byte(`"note": ""`)) {
		t.Fatalf("optional note should be absent: %s", saved)
	}
	assertGalleryFeedback(t, saved, want)
	runGalleryFeedbackRequest(t, origin+"/private.txt", "", http.MethodGet, nil, http.StatusNotFound)
	runGalleryFeedbackRequest(t, origin+"/"+reportGalleryFeedbackFile, "", http.MethodGet, nil, http.StatusNotFound)
	runGalleryFeedbackRequest(t, origin+"/"+report, "", http.MethodGet, nil, http.StatusOK)

	stop()
	reopened, stopReopened := startGalleryFixture(t, root)
	defer stopReopened()
	reopenedOrigin := strings.TrimSuffix(reopened, "/"+ReportGalleryPageFile)
	assertGalleryFeedback(t, runGalleryFeedbackRequest(t, reopenedOrigin+reportGalleryFeedbackAPI, "", http.MethodGet, nil, http.StatusOK), want)
	handMade, err := os.ReadFile(filepath.Join(root, ReportGalleryPageFile))
	if err != nil || string(handMade) != "Do not overwrite this hand-made page." {
		t.Fatalf("live gallery touched static output: %q, %v", handMade, err)
	}

	// The server must reject a bad file and a linked write directory without
	// claiming a successful save or replacing the user's data.
	writeGalleryFixture(t, root, reportGalleryFeedbackFile, "{bad JSON")
	runGalleryFeedbackRequest(t, reopenedOrigin+reportGalleryFeedbackAPI, reopenedOrigin, http.MethodPut, whole, http.StatusConflict)
	broken, err := os.ReadFile(store)
	if err != nil || string(broken) != "{bad JSON" {
		t.Fatalf("failed save replaced malformed data: %q, %v", broken, err)
	}
	outside := t.TempDir()
	writeGalleryFixture(t, outside, filepath.Base(store), string(saved))
	if err := os.Rename(filepath.Dir(store), filepath.Join(root, "old-feedback")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Dir(store)); err != nil {
		t.Fatal(err)
	}
	runGalleryFeedbackRequest(t, reopenedOrigin+reportGalleryFeedbackAPI, reopenedOrigin, http.MethodPut, whole, http.StatusConflict)
	outsideData, err := os.ReadFile(filepath.Join(outside, filepath.Base(store)))
	if err != nil || !bytes.Equal(outsideData, saved) {
		t.Fatalf("feedback escaped the project: %s, %v", outsideData, err)
	}
}

func runGalleryFeedbackRequest(t *testing.T, address, origin, method string, payload any, status int) []byte {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(method, address, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	if origin != "" {
		request.Header.Set("Origin", origin)
	}
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil || response.StatusCode != status {
		t.Fatalf("%s %s returned %d, want %d: %s (%v)", method, address, response.StatusCode, status, content, err)
	}
	return content
}

func assertGalleryFeedback(t *testing.T, content []byte, want reportGalleryFeedback) {
	t.Helper()
	var actual reportGalleryFeedback
	if err := json.Unmarshal(content, &actual); err != nil || !reflect.DeepEqual(actual, want) {
		t.Fatalf("feedback = %s, want %+v (%v)", content, want, err)
	}
}

func startGalleryFixture(t *testing.T, root string) (string, func()) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestAIReportGalleryServeRoundTrip$")
	command.Env = append(os.Environ(), galleryTestRoot+"="+root)
	stdout, err := command.StdoutPipe()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	command.Stderr = command.Stdout
	if err := command.Start(); err != nil {
		cancel()
		t.Fatal(err)
	}
	stop := sync.OnceFunc(func() {
		cancel()
		_ = command.Wait()
	})
	t.Cleanup(stop)
	scanner := bufio.NewScanner(stdout)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if address, found := strings.CutPrefix(line, "Report gallery: "); found {
			return address, stop
		}
		lines = append(lines, line)
	}
	t.Fatalf("live gallery did not announce a URL: %s", strings.Join(lines, "\n"))
	return "", stop
}
