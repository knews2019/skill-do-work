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

// judgeBundleDigest lists every path under the bundle with its bytes (a symlink
// by its target), so a test can prove the judge wrote nothing inside the bundle.
func judgeBundleDigest(t *testing.T, bundle string) map[string]string {
	t.Helper()
	digest := map[string]string{}
	err := filepath.WalkDir(bundle, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		switch {
		case entry.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(path)
			digest[path] = "link:" + target
			return err
		case !entry.IsDir():
			content, err := os.ReadFile(path)
			digest[path] = string(content)
			return err
		}
		digest[path] = "dir"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

func sameJudgeDigest(before, after map[string]string) bool {
	if len(before) != len(after) {
		return false
	}
	for path, content := range before {
		if after[path] != content {
			return false
		}
	}
	return true
}

// Pins the copy-and-swap workaround: to check a second-opinion page an agent had
// to copy the bundle and write the page over the copy's index.html. --page must
// render the named page in place: served_url names it (as the checked, cleaned path,
// so "./x.html" or a Windows "sub\x.html" never reaches the server raw), the phone
// overflow that only that page has is found (index.html alone passes), and the
// bundle is unchanged.
func TestAIReportJudgeRendersNamedPageInPlace(t *testing.T) {
	requireJudgeBrowser(t)
	repository := t.TempDir()
	bundle := filepath.Join(repository, "bundle")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	pages := map[string]string{
		"index.html": `<!doctype html><html><head><meta charset="utf-8"><style>body{margin:0}</style></head><body><p>clean</p></body></html>`,
		"index-v1c01-x.html": `<!doctype html><html><head><meta charset="utf-8"><style>body{margin:0}</style></head>` +
			`<body><div style="width:600px">wide second opinion</div></body></html>`,
	}
	for name, content := range pages {
		if err := os.WriteFile(filepath.Join(bundle, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	before := judgeBundleDigest(t, bundle)
	output := filepath.Join(t.TempDir(), "judge-output")
	result := handleAIReportJudge(commandruntime.ExecutionContext{RepositoryRoot: repository},
		[]string{"bundle", "--page", "./index-v1c01-x.html", "--out", output})
	judgeBytes, err := os.ReadFile(filepath.Join(output, "judge.json"))
	if err != nil {
		t.Fatalf("judge.json was not written: %v (result %+v)", err, result)
	}
	var report judgeReport
	if err := json.Unmarshal(judgeBytes, &report); err != nil {
		t.Fatalf("judge.json is not valid JSON: %v\n%s", err, judgeBytes)
	}
	if served, err := url.Parse(report.ServedURL); err != nil || served.Path != "/index-v1c01-x.html" {
		t.Fatalf("served_url %q does not name the cleaned page path", report.ServedURL)
	}
	if resultmodel.ExitCode(result.Outcome) == 0 || report.Verdict != "fail" {
		t.Fatalf("overflowing page: outcome %s verdict %q findings %+v", result.Outcome, report.Verdict, report.Findings)
	}
	phoneOverflow := false
	for _, finding := range report.Findings {
		if finding.Kind == judgeCodeOverflow && finding.Viewport == "phone" && finding.ScrollWidth >= 600 {
			phoneOverflow = true
		}
	}
	if !phoneOverflow {
		t.Fatalf("no phone overflow finding, so the captures did not show the named page: %+v", report.Findings)
	}
	for _, capture := range []string{"wide-light.png", "wide-dark.png", "phone-light.png", "phone-dark.png"} {
		if info, err := os.Stat(filepath.Join(output, capture)); err != nil || info.Size() == 0 {
			t.Fatalf("capture %s missing or empty: %v", capture, err)
		}
	}
	if !sameJudgeDigest(before, judgeBundleDigest(t, bundle)) {
		t.Fatal("the bundle changed during the --page run")
	}
}

// Pins the containment guard: --page must never render a file that is missing or
// that sits outside the bundle (by "..", an absolute path or a symlink). Each is
// refused with AI-REPORT-JUDGE-PAGE before any browser or judge.json.
func TestAIReportJudgeRefusesMissingOrOutsidePage(t *testing.T) {
	repository := t.TempDir()
	bundle := filepath.Join(repository, "bundle")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(repository, "outside.html")
	for path, content := range map[string]string{filepath.Join(bundle, "index.html"): "<!doctype html><title>report</title>",
		outside: "<!doctype html><title>outside</title>"} {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(outside, filepath.Join(bundle, "linked.html")); err != nil {
		t.Fatal(err)
	}
	before := judgeBundleDigest(t, bundle)
	for _, testCase := range []struct{ name, page string }{
		{"missing", "index-v1c01-missing.html"},
		{"dot-dot escape", "../outside.html"},
		{"absolute path", outside},
		{"symlink outside", "linked.html"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "judge-output")
			result := handleAIReportJudge(commandruntime.ExecutionContext{RepositoryRoot: repository},
				[]string{"bundle", "--page", testCase.page, "--out", output})
			if result.Outcome != resultmodel.OutcomeFailure {
				t.Fatalf("page %q: outcome %s, want failure (%+v)", testCase.page, result.Outcome, result.Findings)
			}
			typed := false
			for _, finding := range result.Findings {
				typed = typed || finding.Code == "AI-REPORT-JUDGE-PAGE"
			}
			if !typed {
				t.Fatalf("page %q: no AI-REPORT-JUDGE-PAGE finding in %+v", testCase.page, result.Findings)
			}
			if _, err := os.Stat(filepath.Join(output, "judge.json")); err == nil {
				t.Fatalf("page %q: a refusal wrote judge.json", testCase.page)
			}
			if !sameJudgeDigest(before, judgeBundleDigest(t, bundle)) {
				t.Fatalf("page %q: the bundle changed", testCase.page)
			}
		})
	}
}
