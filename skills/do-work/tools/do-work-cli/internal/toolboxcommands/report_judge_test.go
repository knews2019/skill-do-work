package toolboxcommands

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io/fs"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// requireJudgeBrowser gates the browser-launching cases like the other heavy
// do-work-cli tests and skips when discovery finds no engine on this machine.
func requireJudgeBrowser(t *testing.T) {
	t.Helper()
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("ai-report-judge browser cases are heavy-only")
	}
	if discoverJudgeBrowser() == "" {
		t.Skip("no browser engine found for ai-report-judge")
	}
}

// runJudgeOnBundle writes index.html (plus extra files) into a fresh bundle, runs
// the verb through its handler with an output directory outside the bundle, and
// returns the result, the decoded judge.json, the bundle and the output directory.
func runJudgeOnBundle(t *testing.T, indexHTML string, extraFiles map[string][]byte) (resultmodel.CommandResult, judgeReport, string, string) {
	t.Helper()
	repository := t.TempDir()
	bundle := filepath.Join(repository, "bundle")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{"index.html": []byte(indexHTML)}
	for name, content := range extraFiles {
		files[name] = content
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(bundle, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	output := filepath.Join(t.TempDir(), "judge-output")
	result := handleAIReportJudge(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"bundle", "--out", output})
	var report judgeReport
	judgeBytes, err := os.ReadFile(filepath.Join(output, "judge.json"))
	if err != nil {
		t.Fatalf("judge.json was not written: %v (result %+v)", err, result)
	}
	if err := json.Unmarshal(judgeBytes, &report); err != nil {
		t.Fatalf("judge.json is not valid JSON: %v\n%s", err, judgeBytes)
	}
	return result, report, bundle, output
}

// Pins the main failure the verb exists for: a report that fits a desktop but
// scrolls sideways on a phone must fail, and judge.json must name the width.
func TestAIReportJudgeFailsOnPhoneWidthOverflow(t *testing.T) {
	requireJudgeBrowser(t)
	result, report, _, _ := runJudgeOnBundle(t, `<!doctype html><html><head><meta charset="utf-8"><style>body{margin:0}</style></head>`+
		`<body><div style="width:600px">wide content</div></body></html>`, nil)
	if resultmodel.ExitCode(result.Outcome) == 0 || report.Verdict != "fail" {
		t.Fatalf("overflowing bundle: outcome %s verdict %q findings %+v", result.Outcome, report.Verdict, report.Findings)
	}
	phoneOverflow := false
	for _, finding := range report.Findings {
		if finding.Kind != judgeCodeOverflow {
			continue
		}
		if finding.Viewport == "wide" {
			t.Fatalf("600px content was reported as overflowing at 1440: %+v", finding)
		}
		if finding.Viewport == "phone" && finding.ScrollWidth >= 600 && finding.InnerWidth == 390 {
			phoneOverflow = true
		}
	}
	if !phoneOverflow {
		t.Fatalf("no phone overflow finding with scroll width >= 600 and inner width 390: %+v", report.Findings)
	}
}

// Pins the broken relative asset: a missing image must fail with the target and
// the HTTP status, not pass because the page still rendered.
func TestAIReportJudgeFailsOnBrokenRelativeImage(t *testing.T) {
	requireJudgeBrowser(t)
	result, report, _, _ := runJudgeOnBundle(t, `<!doctype html><html><head><meta charset="utf-8"></head>`+
		`<body><p>report</p><img src="missing.png" alt=""></body></html>`, nil)
	if resultmodel.ExitCode(result.Outcome) == 0 || report.Verdict != "fail" {
		t.Fatalf("broken image bundle: outcome %s verdict %q findings %+v", result.Outcome, report.Verdict, report.Findings)
	}
	for _, finding := range report.Findings {
		if finding.Kind == judgeCodeBrokenLink && finding.Target == "missing.png" && finding.Status == 404 {
			return
		}
	}
	t.Fatalf("no broken-link finding for missing.png with status 404: %+v", report.Findings)
}

// Pins the clean path: exit 0 with four real captures, the bundle left byte for
// byte as it was, the external link never fetched, and the server gone after return.
func TestAIReportJudgePassesCleanBundleAndStopsServer(t *testing.T) {
	requireJudgeBrowser(t)
	var pngBytes bytes.Buffer
	if err := png.Encode(&pngBytes, image.NewRGBA(image.Rect(0, 0, 8, 8))); err != nil {
		t.Fatal(err)
	}
	indexHTML := `<!doctype html><html><head><meta charset="utf-8"><style>body{margin:0}</style></head><body>` +
		`<p><a href="#details">details</a> <a href="https://example.invalid/">external</a></p>` +
		`<img src="pixel.png" alt="pixel"><h2 id="details">Details</h2></body></html>`
	result, report, bundle, output := runJudgeOnBundle(t, indexHTML, map[string][]byte{"pixel.png": pngBytes.Bytes()})
	if resultmodel.ExitCode(result.Outcome) != 0 || report.Verdict != "pass" || len(report.Findings) != 0 {
		t.Fatalf("clean bundle: outcome %s verdict %q findings %+v", result.Outcome, report.Verdict, report.Findings)
	}
	for _, capture := range []string{"wide-light.png", "wide-dark.png", "phone-light.png", "phone-dark.png"} {
		info, err := os.Stat(filepath.Join(output, capture))
		if err != nil || info.Size() == 0 {
			t.Fatalf("capture %s missing or empty: %v", capture, err)
		}
	}
	bundleFiles := map[string]string{}
	_ = filepath.WalkDir(bundle, func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			content, _ := os.ReadFile(path)
			bundleFiles[path] = string(content)
		}
		return nil
	})
	if len(bundleFiles) != 2 || bundleFiles[filepath.Join(bundle, "index.html")] != indexHTML ||
		bundleFiles[filepath.Join(bundle, "pixel.png")] != pngBytes.String() {
		t.Fatalf("the bundle changed during the judge run: %v", len(bundleFiles))
	}
	served, err := url.Parse(report.ServedURL)
	if err != nil || served.Host == "" {
		t.Fatalf("judge.json served_url %q is not a URL: %v", report.ServedURL, err)
	}
	if connection, dialError := net.DialTimeout("tcp", served.Host, 2*time.Second); dialError == nil {
		connection.Close()
		t.Fatalf("the bundle server at %s still accepts connections after the run", served.Host)
	}
}

// Pins "a skip is not a pass": with no engine found the run must say skipped,
// still write judge.json, and exit non-zero. It never launches a browser.
func TestAIReportJudgeReportsSkippedWithoutBrowser(t *testing.T) {
	savedNames, savedAppPaths := aiReportJudgeBrowserNames, aiReportJudgeBrowserAppPaths
	aiReportJudgeBrowserNames, aiReportJudgeBrowserAppPaths = nil, nil
	t.Cleanup(func() { aiReportJudgeBrowserNames, aiReportJudgeBrowserAppPaths = savedNames, savedAppPaths })
	t.Setenv("PATH", t.TempDir())
	result, report, _, _ := runJudgeOnBundle(t, `<!doctype html><title>report</title>`, nil)
	if resultmodel.ExitCode(result.Outcome) == 0 || report.Verdict != "skipped" {
		t.Fatalf("no-browser run: outcome %s verdict %q", result.Outcome, report.Verdict)
	}
	for _, finding := range result.Findings {
		if finding.Code == judgeCodeSkipped {
			return
		}
	}
	t.Fatalf("no %s finding in %+v", judgeCodeSkipped, result.Findings)
}
