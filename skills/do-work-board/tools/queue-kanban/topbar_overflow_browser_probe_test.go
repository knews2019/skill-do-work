package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The top bar's fit inside the viewport, measured in a real engine.
//
// WHY A BROWSER. "Can the reader reach the controls" is not a fact about any
// string the renderer writes. REQ-586 folded the identity onto one `nowrap`
// line and pinned the result with markup assertions, which passed while the
// bar's content wanted 905px inside an 800px viewport and `body { overflow:
// hidden }` clipped the last 105px of the controls with no scrollbar to reach
// them. Nothing short of a laid-out page can see that: it is the product of the
// identity's text width, the controls' wrapped width, the bar's padding and gap,
// and which media query is in force.
//
// Two widths, because the failure and the feature live on opposite sides of one
// breakpoint. At 900px the bar must stack, and both halves must sit inside the
// viewport. At 1400px the bar must still be the single line REQ-586 delivered —
// a fix that reached the narrow band by giving up the wide one would pass a
// one-width probe.
//
// Following REQ-291's lesson, every box this probe compares is asserted to have
// been laid out before any comparison is trusted: an unrendered top bar reports
// a right edge of 0 and reads as "inside the viewport".

// topbarOverflowProbeStackedWidth sits inside the 761-999px band that REQ-586
// left laying the identity and the controls out side by side. It is the width
// the regression was measured at, one step in from the 999px edge.
const topbarOverflowProbeStackedWidth = 900

// topbarOverflowProbeOneLineWidth is the width REQ-586 measured its 126px -> 68px
// improvement at, so it is the width that says whether the improvement survives.
const topbarOverflowProbeOneLineWidth = 1400

// topbarOverflowProbeLongProjectName is longer than any real repository directory
// name, and is written into the identity at measurement time. It is what proves
// the ellipsis cap on .board-project is doing the work rather than this repo's
// own short name happening to fit.
const topbarOverflowProbeLongProjectName = "a-project-directory-name-far-longer-than-any-real-checkout-uses"

// topbarFitMeasurement is everything one width reports back in a single node.
// Edges rather than a pair of booleans, so a failure says how far past the
// viewport the bar reached instead of only that it did.
type topbarFitMeasurement struct {
	ViewportWidth float64 `json:"viewportWidth"`

	// DocumentScrollWidth is the page's own overflow verdict. Anything wider
	// than the viewport is content the reader cannot scroll to, because the
	// body hides horizontal overflow.
	DocumentScrollWidth float64 `json:"documentScrollWidth"`

	TopBarHeight   float64 `json:"topBarHeight"`
	IdentityLeft   float64 `json:"identityLeft"`
	IdentityRight  float64 `json:"identityRight"`
	IdentityBottom float64 `json:"identityBottom"`
	ControlsLeft   float64 `json:"controlsLeft"`
	ControlsRight  float64 `json:"controlsRight"`
	ControlsTop    float64 `json:"controlsTop"`
	ControlsWidth  float64 `json:"controlsWidth"`

	// ProjectNameTruncated is whether the project name is being clipped by its
	// own max-width, read as scrollWidth exceeding clientWidth.
	ProjectNameTruncated bool `json:"projectNameTruncated"`
}

func TestBrowserBehaviorTopBarStaysInsideTheViewport(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)

	siteDirectory, indexHTML := topbarOverflowProbeGeneratedBoard(t)

	t.Run("stacked below the breakpoint", func(t *testing.T) {
		measured := topbarOverflowProbeMeasureAtWidth(
			t, siteDirectory, indexHTML, topbarOverflowProbeStackedWidth, false)

		if measured.ControlsRight > measured.ViewportWidth {
			t.Errorf("at %.0fpx the controls end at %.1fpx, so their last %.1fpx are clipped by the body's hidden overflow and no scrollbar reaches them: %+v",
				measured.ViewportWidth, measured.ControlsRight,
				measured.ControlsRight-measured.ViewportWidth, measured)
		}
		if measured.DocumentScrollWidth > measured.ViewportWidth {
			t.Errorf("at %.0fpx the page's content is %.1fpx wide: %+v",
				measured.ViewportWidth, measured.DocumentScrollWidth, measured)
		}
		// The stacking rule is what makes the fit above possible at this width,
		// so it is asserted rather than assumed: a bar that happened to fit side
		// by side on a wider control set would pass the edge checks and break
		// again the next time a control group is added.
		if measured.ControlsTop < measured.IdentityBottom {
			t.Errorf("at %.0fpx the controls start %.1fpx above the identity's baseline box bottom (%.1fpx), so the bar is still one row and the 999px stacking rule did not apply: %+v",
				measured.ViewportWidth, measured.ControlsTop, measured.IdentityBottom, measured)
		}
	})

	t.Run("one line above the breakpoint", func(t *testing.T) {
		measured := topbarOverflowProbeMeasureAtWidth(
			t, siteDirectory, indexHTML, topbarOverflowProbeOneLineWidth, false)

		if measured.ControlsRight > measured.ViewportWidth {
			t.Errorf("at %.0fpx the controls end at %.1fpx, past the viewport: %+v",
				measured.ViewportWidth, measured.ControlsRight, measured)
		}
		// REQ-586's delivered shape: the identity and the controls share one row.
		// Read as a horizontal relationship rather than a height threshold, which
		// would break every time a control group changes its own wrapping.
		if measured.ControlsLeft < measured.IdentityRight {
			t.Errorf("at %.0fpx the controls start at %.1fpx, left of the identity's right edge %.1fpx, so the bar is no longer the one line REQ-586 delivered: %+v",
				measured.ViewportWidth, measured.ControlsLeft, measured.IdentityRight, measured)
		}
		if measured.ProjectNameTruncated {
			t.Errorf("at %.0fpx this fixture's ordinary project name is being truncated, so the cap on .board-project is too tight for real names: %+v",
				measured.ViewportWidth, measured)
		}
	})

	t.Run("a very long project name is capped, not carried", func(t *testing.T) {
		measured := topbarOverflowProbeMeasureAtWidth(
			t, siteDirectory, indexHTML, topbarOverflowProbeOneLineWidth, true)

		if !measured.ProjectNameTruncated {
			t.Errorf("a %d-character project name is not being truncated, so nothing bounds how wide the identity's nowrap line can grow: %+v",
				len(topbarOverflowProbeLongProjectName), measured)
		}
		if measured.ControlsRight > measured.ViewportWidth {
			t.Errorf("a %d-character project name pushed the controls to %.1fpx, past the %.0fpx viewport: %+v",
				len(topbarOverflowProbeLongProjectName), measured.ControlsRight,
				measured.ViewportWidth, measured)
		}
	})
}

// topbarOverflowProbeGeneratedBoard builds one real static board for every width
// to reuse. The bar's width depends on the control groups, not on the queue, so
// a small fixture measures the same layout a full queue does.
func topbarOverflowProbeGeneratedBoard(t *testing.T) (string, string) {
	t.Helper()

	stampBase := time.Now().UTC().Add(-2 * time.Hour)
	fixtureFiles := make([]verifyFixtureFile, 0, 3)
	for requestIndex := 0; requestIndex < 3; requestIndex++ {
		requestID := fmt.Sprintf("REQ-%03d", 700+requestIndex)
		fixtureFiles = append(fixtureFiles, verifyFixtureFile{
			RelativePath: "do-work/archive/" + requestID + "-topbar-fit.md",
			Content: "---\n" +
				"id: " + requestID + "\n" +
				"title: Top bar fit fixture " + requestID + "\n" +
				"status: completed\n" +
				"created_at: " + stampBase.Format("2006-01-02T15:04:05Z") + "\n" +
				"completed_at: " + stampBase.Format("2006-01-02T15:04:05Z") + "\n" +
				"---\n",
		})
	}
	repoRoot := writeVerifyFixture(t, fixtureFiles)

	board, buildError := buildBoard(repoRoot, time.Now().UTC(), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatal(buildError)
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, board); generateError != nil {
		t.Fatal(generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatal(readError)
	}
	return siteDirectory, string(indexBytes)
}

// topbarOverflowProbeMeasureAtWidth opens the board in its own engine at one
// window width and reports the bar's geometry. A separate engine per width, not
// a resize: the window size is the flag the media queries resolve against when
// the page first lays out, and a resized window is a different measurement to
// argue about than the one a reader with that screen actually gets.
func topbarOverflowProbeMeasureAtWidth(
	t *testing.T, siteDirectory string, indexHTML string, viewportWidth int, useLongProjectName bool,
) topbarFitMeasurement {
	t.Helper()

	session := startTrustedInputBrowserSession(t,
		fmt.Sprintf("top bar fit at %dpx", viewportWidth), siteDirectory, indexHTML,
		fmt.Sprintf("--window-size=%d,900", viewportWidth))
	defer session.closeBrowserSession()

	session.waitForPageCondition(t, "the top bar and its controls",
		`document.querySelector('.board-topbar') && document.querySelector('.board-controls')`)

	if useLongProjectName {
		session.evaluateInPage(t, fmt.Sprintf(
			`(document.getElementById('board-project').textContent = %q, "renamed")`,
			topbarOverflowProbeLongProjectName))
	}

	var measured topbarFitMeasurement
	session.decodeResult(t, "top bar fit",
		session.evaluateInPage(t, topbarFitProbeExpression()), &measured)

	// REQ-291's lesson: refuse to reason about a box that was never laid out.
	// Without this an unrendered bar reports every edge at 0, which passes every
	// "inside the viewport" check in this file.
	if measured.TopBarHeight <= 0 || measured.ControlsWidth <= 0 {
		t.Fatalf("the top bar was not laid out (height %.1f, controls width %.1f), so nothing measured here means anything: %+v",
			measured.TopBarHeight, measured.ControlsWidth, measured)
	}
	if measured.ViewportWidth <= 0 {
		t.Fatalf("the page reports a viewport of %.1fpx: %+v", measured.ViewportWidth, measured)
	}

	t.Logf("top bar fit viewport=%.0f documentScrollWidth=%.1f height=%.1f identity=[%.1f,%.1f] controls=[%.1f,%.1f] controlsTop=%.1f projectTruncated=%t",
		measured.ViewportWidth, measured.DocumentScrollWidth, measured.TopBarHeight,
		measured.IdentityLeft, measured.IdentityRight,
		measured.ControlsLeft, measured.ControlsRight, measured.ControlsTop,
		measured.ProjectNameTruncated)

	return measured
}

// topbarFitProbeExpression reads the bar's geometry after one animation frame, so
// a project name written in the same task has been laid out before it is measured.
func topbarFitProbeExpression() string {
	return `new Promise(function (resolveMeasurement) {
  requestAnimationFrame(function () {
    requestAnimationFrame(function () {
      var topBar = document.querySelector('.board-topbar');
      var identity = document.querySelector('.board-identity');
      var controls = document.querySelector('.board-controls');
      var projectName = document.getElementById('board-project');
      var topBarBox = topBar.getBoundingClientRect();
      var identityBox = identity.getBoundingClientRect();
      var controlsBox = controls.getBoundingClientRect();
      resolveMeasurement({
        viewportWidth: window.innerWidth,
        documentScrollWidth: document.documentElement.scrollWidth,
        topBarHeight: topBarBox.height,
        identityLeft: identityBox.left,
        identityRight: identityBox.right,
        identityBottom: identityBox.bottom,
        controlsLeft: controlsBox.left,
        controlsRight: controlsBox.right,
        controlsTop: controlsBox.top,
        controlsWidth: controlsBox.width,
        projectNameTruncated: projectName.scrollWidth > projectName.clientWidth + 1
      });
    });
  });
})`
}
