package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
)

// The browser behavior lane, beside the Node behavior lane in generate_test.go and
// deliberately built to the same shape rather than to a second convention.
//
// WHY IT EXISTS. Every font measurement in the Durations view got into the code by a
// person running a browser by hand and pasting the number into a comment with the
// build written beside it. That ritual is why those constants are stale — the repo's
// own comment admits the recorded box height "is NOT a supremum over the face space"
// — and it is why surveying operating systems by hand was ever proposed. A test that
// can ask a real engine for a text extent replaces the ritual. The lane's value is
// not the one probe below; it is that the next person needing a rendered measurement
// gets it from a test instead of from a browser session and a comment.
//
// WHICH DRIVER, AND WHY. A headless Chrome/Chromium invoked directly against a temp
// file, using `--dump-dom` and `--virtual-time-budget`. Chosen over Playwright and
// anything else needing npm because it requires only a binary most machines already
// have: a silent package-manager dependency is a build input nobody agreed to. The
// page writes its results into one DOM node and the Go side parses that node, so
// results cross the boundary as data this test can assert on rather than as an exit
// code alone. If a future probe genuinely needs more than `--dump-dom` can express,
// that is the moment to reach for a driver — and to say so here.
//
// NO WALL-CLOCK WAITS. Readiness is `--virtual-time-budget`, which advances the
// page's clock as fast as the work allows and then dumps, plus the sentinel the page
// writes into the result node. A `sleep` here is the flake this lane would be blamed
// for.
const (
	strictBrowserBehaviorDiagnostic = "queue-kanban: strict browser behavior lane executed zero probes"
	strictBrowserBehaviorMarker     = "QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR"
	strictBrowserBehaviorRunPattern = "^TestMaintainerStrictBrowserBehaviorLane$"

	// browserProbeBinaryOverride names a browser explicitly, for a machine whose
	// engine is not on PATH under any well-known name — a Playwright-managed build,
	// a distribution package in an unusual place, a pinned build under test.
	browserProbeBinaryOverride = "QUEUE_KANBAN_BROWSER"

	// browserProbeResultElementId is the single node the page writes into and the Go
	// side reads back. One node, one contract: everything the probe wants to report
	// is JSON inside it.
	browserProbeResultElementId = "queue-kanban-probe-result"
)

var browserBehaviorProbeCount atomic.Int64

// browserProbeWellKnownBinaries is checked after the override, in order. It is a
// convenience list, never a closed set — the override above is what makes an
// unlisted engine usable without editing this slice.
var browserProbeWellKnownBinaries = []string{
	"google-chrome",
	"google-chrome-stable",
	"chromium",
	"chromium-browser",
	"chrome",
}

// lookupBrowserForBehaviorProbe mirrors lookupNodeForJavaScriptProbe: consult the
// environment, then PATH, then SKIP. A machine with no browser still runs everything
// else in the suite — which is only safe because the strict lane below refuses to
// skip, so a missing browser can never quietly become a green run for the maintainer.
func lookupBrowserForBehaviorProbe(t *testing.T) string {
	t.Helper()
	if overriddenBrowser := strings.TrimSpace(os.Getenv(browserProbeBinaryOverride)); overriddenBrowser != "" {
		if _, statError := os.Stat(overriddenBrowser); statError == nil {
			return overriddenBrowser
		}
		resolvedOverride, lookupError := exec.LookPath(overriddenBrowser)
		if lookupError == nil {
			return resolvedOverride
		}
		// An override that names nothing is a mistake worth failing on, not a reason
		// to silently fall back: the caller asked for a specific engine.
		t.Fatalf("%s=%q names no runnable browser", browserProbeBinaryOverride, overriddenBrowser)
	}
	for _, candidateBinary := range browserProbeWellKnownBinaries {
		if resolvedBinary, lookupError := exec.LookPath(candidateBinary); lookupError == nil {
			return resolvedBinary
		}
	}
	t.Skipf("no browser is available (set %s to name one); skipping browser behavior probe",
		browserProbeBinaryOverride)
	return ""
}

// runBrowserBehaviorProbe renders pageHTML in a real engine and returns whatever the
// page wrote into its result node, as raw JSON text for the caller to unmarshal.
func runBrowserBehaviorProbe(t *testing.T, probeName string, pageHTML string) []byte {
	t.Helper()
	return runBrowserBehaviorProbeWithFlags(t, probeName, pageHTML)
}

// runBrowserBehaviorProbeWithFlags is the same probe with extra engine flags.
//
// It exists for one reason worth naming: without a colour-scheme flag Chromium
// resolves `prefers-color-scheme` to light, so every probe in this lane measures
// the light palette and NOTHING automated ever sees the dark one — which on this
// board is the `:root` base that the light block overrides. A view whose meaning
// is carried by colour needs both, and a one-time manual table is not a check
// that survives the next edit.
func runBrowserBehaviorProbeWithFlags(t *testing.T, probeName string, pageHTML string, extraFlags ...string) []byte {
	t.Helper()
	return runBrowserBehaviorProbeInDirectory(t, probeName, t.TempDir(), pageHTML, extraFlags...)
}

// runBrowserBehaviorProbeInDirectory writes the probe page into a directory the
// caller chose rather than a fresh temp one.
//
// It exists for the probes that must drive the REAL generated board: index.html
// loads board-data.js from beside itself, so a page copied into an empty temp
// directory renders an empty board and every assertion measures nothing. Pass
// the output of generateLiveSiteInDir and the probe runs against the page a user
// would actually open.
func runBrowserBehaviorProbeInDirectory(
	t *testing.T, probeName string, siteDirectory string, pageHTML string, extraFlags ...string,
) []byte {
	t.Helper()
	browserPath := lookupBrowserForBehaviorProbe(t)

	probeDirectory := siteDirectory
	pagePath := filepath.Join(probeDirectory, "probe.html")
	if writeError := os.WriteFile(pagePath, []byte(pageHTML), 0o644); writeError != nil {
		t.Fatalf("write %s probe page: %v", probeName, writeError)
	}
	// The profile never goes in the site directory: it is megabytes of engine
	// state, and a probe that leaves it beside index.html changes what the next
	// probe against the same directory sees.
	profileDirectory := t.TempDir()

	// --no-sandbox because CI and container users are routinely root, where the
	// sandbox refuses to start; --disable-gpu and --disable-dev-shm-usage because a
	// headless container has neither. --user-data-dir keeps concurrent probes from
	// fighting over one profile directory.
	probeArguments := []string{
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--user-data-dir=" + filepath.Join(profileDirectory, "profile"),
		"--virtual-time-budget=5000",
		"--dump-dom",
	}
	probeArguments = append(probeArguments, extraFlags...)
	probeArguments = append(probeArguments, "file://"+pagePath)
	probeCommand := exec.Command(browserPath, probeArguments...)
	browserBehaviorProbeCount.Add(1)
	probeOutput, probeError := probeCommand.CombinedOutput()
	if probeError != nil {
		t.Fatalf("execute %s browser behavior: %v\n%s", probeName, probeError, probeOutput)
	}

	resultText, extractError := extractBrowserProbeResult(string(probeOutput))
	if extractError != "" {
		t.Fatalf("read %s browser behavior result: %s\ndumped DOM:\n%s", probeName, extractError, probeOutput)
	}
	return []byte(resultText)
}

// browserProbeResultPattern finds the result node in the dumped DOM. The dump is
// serialized HTML, not a parse tree, so this reads the one node the contract defines
// rather than trying to parse the document.
var browserProbeResultPattern = regexp.MustCompile(
	`(?s)<[a-zA-Z]+[^>]*\bid="` + regexp.QuoteMeta(browserProbeResultElementId) + `"[^>]*>(.*?)</[a-zA-Z]+>`)

func extractBrowserProbeResult(dumpedDOM string) (string, string) {
	resultMatch := browserProbeResultPattern.FindStringSubmatch(dumpedDOM)
	if resultMatch == nil {
		return "", "the page never wrote its result node — the probe script threw, or the engine dumped before it ran"
	}
	resultText := strings.TrimSpace(resultMatch[1])
	if resultText == "" {
		return "", "the result node is empty"
	}
	// The sentinel: the page sets the node only after measuring, so an empty or
	// absent node means "did not finish", never "measured zero".
	if !strings.HasPrefix(resultText, "{") {
		return "", "the result node holds " + resultText + ", not the JSON the contract expects"
	}
	return resultText, ""
}

// markLabelTextExtent is what the page measures and the Go side asserts on.
type markLabelTextExtent struct {
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	FontFamily  string  `json:"fontFamily"`
	MeasuredPx  float64 `json:"measuredPx"`
	SampleLabel string  `json:"sampleLabel"`
}

// The one real probe. It measures a .durations-mark-label <text> at the board's own
// 11px through getBBox() and asserts the returned box is positive and finite.
//
// It names the failure it pins: a lane that renders nothing measurable is a lane that
// passes forever. getBBox() on an unrendered or detached element returns zeros, and a
// page that never ran returns no node at all — both fail here rather than reading as
// a successful measurement of nothing.
func TestBrowserBehaviorMarkLabelTextExtent(t *testing.T) {
	const sampleLabel = "REQ-042 · 3h 20m"
	probePage := `<!doctype html>
<html><head><meta charset="utf-8"><style>
  .durations-mark-label { font-size: 11px; font-family: ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, sans-serif; }
</style></head>
<body>
<svg id="probe-svg" width="600" height="80" xmlns="http://www.w3.org/2000/svg">
  <text id="probe-label" class="durations-mark-label" x="10" y="40">` + sampleLabel + `</text>
</svg>
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
  // Write the result node LAST and only once, so its presence is the sentinel that
  // the measurement completed. Any throw leaves the node empty and fails the Go side
  // with the dumped DOM attached, rather than reporting a zero measurement.
  (function () {
    var label = document.getElementById("probe-label");
    var box = label.getBBox();
    var computed = window.getComputedStyle(label);
    document.getElementById("` + browserProbeResultElementId + `").textContent =
      JSON.stringify({
        width: box.width,
        height: box.height,
        fontFamily: computed.fontFamily,
        measuredPx: parseFloat(computed.fontSize),
        sampleLabel: label.textContent
      });
  })();
</script>
</body></html>`

	resultJSON := runBrowserBehaviorProbe(t, "mark-label text extent", probePage)

	var measuredExtent markLabelTextExtent
	if unmarshalError := json.Unmarshal(resultJSON, &measuredExtent); unmarshalError != nil {
		t.Fatalf("parse probe result %s: %v", resultJSON, unmarshalError)
	}

	// Positive and finite. Zero is what an unrendered element reports, which is
	// exactly the silent-pass this assertion exists to prevent.
	if !(measuredExtent.Width > 0) || !(measuredExtent.Height > 0) {
		t.Fatalf("measured box is not positive: %+v — an unrendered element measures zero", measuredExtent)
	}
	if measuredExtent.Width > 10000 || measuredExtent.Height > 10000 {
		t.Fatalf("measured box is implausible: %+v", measuredExtent)
	}
	// The engine must actually have applied the stylesheet; a 11px rule that did not
	// take would make every measurement below meaningless.
	if measuredExtent.MeasuredPx != 11 {
		t.Fatalf("probe measured at %gpx, want 11px — the mark-label rule did not apply: %+v",
			measuredExtent.MeasuredPx, measuredExtent)
	}
	if measuredExtent.SampleLabel != sampleLabel {
		t.Fatalf("probe measured %q, want %q", measuredExtent.SampleLabel, sampleLabel)
	}
	// A 16-character label at 11px cannot be narrower than its own height.
	if measuredExtent.Width <= measuredExtent.Height {
		t.Fatalf("a %d-character label measured %gx%g — narrower than tall is not a real text extent",
			len(sampleLabel), measuredExtent.Width, measuredExtent.Height)
	}
	t.Logf("measured %q at %gpx in %s: %.2f x %.2f",
		measuredExtent.SampleLabel, measuredExtent.MeasuredPx, measuredExtent.FontFamily,
		measuredExtent.Width, measuredExtent.Height)
}

// The zero-probe guard. Without it a lane whose probes all skipped — no browser, a
// renamed test — reports green, which is the failure mode that makes a skippable lane
// dangerous. Mirrors TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes.
func TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes(t *testing.T) {
	strictCommand := exec.Command(os.Args[0], "-test.run=^TestBrowserBehavior", "-test.count=1")
	strictCommand.Env = testEnvironmentWithOverrides(
		os.Environ(),
		"PATH="+t.TempDir(),
		browserProbeBinaryOverride+"=",
		strictBrowserBehaviorMarker+"=1",
	)
	strictOutput, strictError := strictCommand.CombinedOutput()
	if strictError == nil {
		t.Fatalf("strict browser behavior lane exited zero without a browser; output:\n%s", strictOutput)
	}
	if !strings.Contains(string(strictOutput), strictBrowserBehaviorDiagnostic) {
		t.Fatalf("strict browser behavior lane output = %q, want %q", strictOutput, strictBrowserBehaviorDiagnostic)
	}
}

// The strict lane: when the maintainer selects it directly, a skip is a failure.
// This is what makes the ordinary skip above safe.
func TestMaintainerStrictBrowserBehaviorLane(t *testing.T) {
	testRunFlag := flag.Lookup("test.run")
	if testRunFlag == nil || testRunFlag.Value.String() != strictBrowserBehaviorRunPattern {
		t.Skip("maintainer strict browser behavior lane runs only when selected directly")
	}

	strictCommand := exec.Command(os.Args[0], "-test.run=^TestBrowserBehavior", "-test.count=1")
	strictCommand.Env = testEnvironmentWithOverrides(
		os.Environ(),
		strictBrowserBehaviorMarker+"=1",
	)
	strictOutput, strictError := strictCommand.CombinedOutput()
	if strictError != nil {
		t.Fatalf("strict browser behavior lane failed: %v\n%s", strictError, strictOutput)
	}
}
