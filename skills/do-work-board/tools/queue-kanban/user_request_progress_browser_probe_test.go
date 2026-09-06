package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The pixel half of REQ-486. Everything here is a claim no JavaScript assertion
// can hold: where a box actually lands, whether two boxes overlap, whether text
// was clipped, and what a colour measures against the ground it is painted on.
//
// The lane drives the REAL generated page — index.html beside its board-data.js
// — because a summary strip measured in an empty temp directory is a strip in a
// board with no user requests in it.

type userRequestProgressRect struct {
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
}

type userRequestProgressTextBox struct {
	ClassName   string  `json:"className"`
	Text        string  `json:"text"`
	ScrollWidth float64 `json:"scrollWidth"`
	ClientWidth float64 `json:"clientWidth"`
	Contrast    float64 `json:"contrast"`
	Color       string  `json:"color"`
}

type userRequestProgressProbeResult struct {
	LocationHref      string                       `json:"locationHref"`
	UserAgent         string                       `json:"userAgent"`
	ResolvedScheme    string                       `json:"resolvedScheme"`
	BodyBackground    string                       `json:"bodyBackground"`
	ViewportWidth     float64                      `json:"viewportWidth"`
	HeadRowRect       userRequestProgressRect      `json:"headRowRect"`
	StripRect         userRequestProgressRect      `json:"stripRect"`
	GroupRect         userRequestProgressRect      `json:"groupRect"`
	MetricRects       []userRequestProgressRect    `json:"metricRects"`
	TitleRect         userRequestProgressRect      `json:"titleRect"`
	DetailRect        userRequestProgressRect      `json:"detailRect"`
	OverlappingPairs  []string                     `json:"overlappingPairs"`
	EscapedMetrics    []string                     `json:"escapedMetrics"`
	TextBoxes         []userRequestProgressTextBox `json:"textBoxes"`
	HeaderMetricTexts []string                     `json:"headerMetricTexts"`
	DrawerMetricTexts []string                     `json:"drawerMetricTexts"`
	HeadTagName       string                       `json:"headTagName"`
	HeadTabIndex      int                          `json:"headTabIndex"`
	DetailTagName     string                       `json:"detailTagName"`
	DetailTabIndex    int                          `json:"detailTabIndex"`
	HeadPrecedesDetil bool                         `json:"headPrecedesDetail"`
	HeadFocusable     bool                         `json:"headFocusable"`
	DetailFocusable   bool                         `json:"detailFocusable"`
	LargeMemberCount  int                          `json:"largeMemberCount"`
	IdListExpanded    bool                         `json:"idListExpanded"`
	InputRowRect      userRequestProgressRect      `json:"inputRowRect"`
	DrawerRect        userRequestProgressRect      `json:"drawerRect"`
	ConsoleErrors     []string                     `json:"consoleErrors"`
}

// userRequestProgressFixtureBoard composes the two user requests this lane needs
// across queue, working and archive, the way userRequestCopyFixture's board
// does: one small UR whose members exercise every branch the strip can render,
// and one 43-member UR for the drawer-containment claim.
func userRequestProgressFixtureBoard(t *testing.T) *Board {
	t.Helper()

	summaryUserRequest := userRequestCopyFixture("UR-950", "Progress summary probe", "[REQ-950]", "SUMMARY-UR-BODY")
	largeUserRequest := userRequestCopyFixture("UR-951", "Forty-three grouped requests", "[]", "LARGE-UR-BODY")

	// A completion the Go side can actually measure a span for: terminal success
	// with a parseable origin and completion. The browser never re-measures it —
	// it reads hasImplementationSpan / implementationSpanMinutes.
	shippedRequest := boardColumnCopyFixtureTicket("REQ-950", "Shipped member", "completed", "frontend", "SHIPPED-BODY")
	shippedRequest.UserRequestId = summaryUserRequest.UserRequestId
	shippedRequest.TreeSection = "archive"
	shippedRequest.CreatedAt = "2026-08-20T09:00:00Z"
	shippedRequest.ClaimedAt = "2026-08-20T09:30:00Z"
	shippedRequest.CompletedAt = "2026-08-20T11:15:00Z"

	// A live claim, so the header and the drawer both carry the ticking attribute
	// and the strip renders a growing figure rather than a static one.
	runningRequest := boardColumnCopyFixtureTicket("REQ-951", "Running member", "claimed", "backend", "RUNNING-BODY")
	runningRequest.UserRequestId = summaryUserRequest.UserRequestId
	runningRequest.TreeSection = "working"
	runningRequest.ClaimedAt = time.Now().UTC().Add(-95 * time.Minute).Format(time.RFC3339)

	// The nested forecast REQ-486 taught the board to read, so the Remaining
	// figure has a real saved number behind it rather than only the fallback.
	forecastRequest := boardColumnCopyFixtureTicket("REQ-952", "Forecast member", "pending", "docs", "FORECAST-BODY")
	forecastRequest.UserRequestId = summaryUserRequest.UserRequestId
	forecastRequest.HasEstimateP50ActiveMinutes = true
	forecastRequest.EstimateP50ActiveMinutes = 145

	// A cancellation, so the strip renders its unmeasured qualifier: the widest
	// text the layout has to survive is the one carrying disclosure.
	cancelledRequest := boardColumnCopyFixtureTicket("REQ-953", "Cancelled member", "cancelled", "general", "CANCELLED-BODY")
	cancelledRequest.UserRequestId = summaryUserRequest.UserRequestId
	cancelledRequest.TreeSection = "archive"
	cancelledRequest.CompletedAt = "2026-08-21T10:00:00Z"

	allRequests := []*RequestTicket{shippedRequest, runningRequest, forecastRequest, cancelledRequest}
	var largeMembers []*RequestTicket
	for memberIndex := 0; memberIndex < 43; memberIndex++ {
		member := boardColumnCopyFixtureTicket(
			fmt.Sprintf("REQ-%d", 1000+memberIndex),
			fmt.Sprintf("Grouped member %02d", memberIndex),
			"pending", "general", "LARGE-MEMBER-BODY")
		member.UserRequestId = largeUserRequest.UserRequestId
		largeMembers = append(largeMembers, member)
	}
	allRequests = append(allRequests, largeMembers...)

	fixtureBoard := &Board{
		ProjectName:  "REQ-486 UR progress probe",
		AllRequests:  allRequests,
		UserRequests: []*UserRequestTicket{summaryUserRequest, largeUserRequest},
		UserRequestsById: map[string]*UserRequestTicket{
			summaryUserRequest.UserRequestId: summaryUserRequest,
			largeUserRequest.UserRequestId:   largeUserRequest,
		},
		Columns: BoardColumns{
			Pending:      append([]*RequestTicket{forecastRequest}, largeMembers...),
			PendingReady: append([]*RequestTicket{forecastRequest}, largeMembers...),
			Claimed:      []*RequestTicket{runningRequest},
			RecentlyDone: []*RequestTicket{shippedRequest},
		},
	}
	linkRequestsToUserRequests(fixtureBoard)
	for _, userRequest := range fixtureBoard.UserRequests {
		sortRequestIdList(userRequest.RequestIds)
	}
	return fixtureBoard
}

// measureUserRequestProgressAtWidth sets the layout viewport on an already-open
// engine and asks the page to measure once.
//
// Two things it exists for. --window-size cannot express 320: Chromium clamps a
// window to 500 CSS px, lays the page out at 500, and a probe asking for 320
// then reports a pass for a width it never saw — a silent skip wearing a green.
// Emulation.setDeviceMetricsOverride is the only way to reach the narrow case,
// and the caller checks the width the page reports back, so an engine that
// ignored the override fails loudly.
//
// It also keeps this probe cheap on a shared machine. One engine per theme
// measures all three widths and the page hands its result back as a resolved
// promise, so a case costs two protocol round trips rather than an engine launch
// plus a poll loop. Under the full heavy gate that poll loop was enough traffic
// to time a Runtime.evaluate out at 30s — which fails as a transport error and
// says nothing at all about the layout.
func measureUserRequestProgressAtWidth(
	t *testing.T, session *trustedInputBrowserSession, viewportWidth int,
) []byte {
	t.Helper()
	session.callDevToolsMethod(t, "Emulation.setDeviceMetricsOverride", map[string]any{
		"width":             viewportWidth,
		"height":            1100,
		"deviceScaleFactor": 1,
		"mobile":            false,
	}, true)
	var resultText string
	session.decodeResult(t, "UR progress measurement",
		session.evaluateInPage(t, "window.__queueKanbanMeasure()"), &resultText)
	resultText = strings.TrimSpace(resultText)
	if !strings.HasPrefix(resultText, "{") {
		t.Fatalf("the page returned %q, not the JSON object the contract expects", resultText)
	}
	return []byte(resultText)
}

// REQ-486: the progress strip has to survive three widths in both themes without
// colliding with the UR title or the Details control, without clipping any of
// its own text, and without dropping below the contrast floor. None of that is
// visible to a Node probe: a jsdom-free DOM stub has no layout and no colour.
func TestBrowserBehaviorUserRequestProgressStripSurvivesEveryWidth(t *testing.T) {
	fixtureBoard := userRequestProgressFixtureBoard(t)
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate UR progress fixture: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read UR progress fixture: %v", readError)
	}

	errorStub := `
      window.__queueKanbanProbeErrors = [];
      window.addEventListener("error", function (event) {
        window.__queueKanbanProbeErrors.push(event.message || "window error");
      });
      var queueKanbanOriginalConsoleError = console.error;
      console.error = function () {
        window.__queueKanbanProbeErrors.push(Array.prototype.join.call(arguments, " "));
        queueKanbanOriginalConsoleError.apply(console, arguments);
      };
`
	probeScript := `
      (function () {
        var resultNode = document.createElement("pre");
        resultNode.id = "` + browserProbeResultElementId + `";
        document.body.appendChild(resultNode);

        function waitFor(predicate, failureLabel) {
          return new Promise(function (resolve, reject) {
            var attempts = 0;
            function poll() {
              if (predicate()) { resolve(); return; }
              attempts += 1;
              if (attempts > 400) { reject(new Error("timed out waiting for " + failureLabel)); return; }
              setTimeout(poll, 10);
            }
            poll();
          });
        }

        function rectOf(element) {
          var rect = element.getBoundingClientRect();
          return { top: rect.top, right: rect.right, bottom: rect.bottom, left: rect.left };
        }
        function rectsIntersect(first, second) {
          return first.left < second.right - 0.5 && second.left < first.right - 0.5 &&
            first.top < second.bottom - 0.5 && second.top < first.bottom - 0.5;
        }
        function rectContains(outer, inner) {
          return inner.left >= outer.left - 0.5 && inner.right <= outer.right + 0.5 &&
            inner.top >= outer.top - 0.5 && inner.bottom <= outer.bottom + 0.5;
        }
        function parseColorChannels(colorText) {
          var numbers = String(colorText).match(/[\d.]+/g) || [];
          return numbers.slice(0, 3).map(Number);
        }
        function relativeLuminance(channels) {
          var linear = channels.map(function (channel) {
            var normalized = channel / 255;
            return normalized <= 0.03928 ? normalized / 12.92 : Math.pow((normalized + 0.055) / 1.055, 2.4);
          });
          return 0.2126 * linear[0] + 0.7152 * linear[1] + 0.0722 * linear[2];
        }
        // Measured against the BODY ground, which is the surface behind this
        // board's content, never a --surface token read out of the stylesheet.
        function contrastAgainstBody(colorText) {
          var foreground = relativeLuminance(parseColorChannels(colorText));
          var background = relativeLuminance(parseColorChannels(getComputedStyle(document.body).backgroundColor));
          var lighter = Math.max(foreground, background);
          var darker = Math.min(foreground, background);
          return (lighter + 0.05) / (darker + 0.05);
        }

        function groupFor(userRequestId) {
          var detailButton = document.querySelector('.ur-group-detail[data-detail-id="' + userRequestId + '"]');
          if (!detailButton) { throw new Error("no By-UR group for " + userRequestId); }
          return detailButton.closest(".ur-group");
        }

        function nextFrame() {
          return new Promise(function (resolve) { requestAnimationFrame(function () { resolve(); }); });
        }

        // Re-callable on purpose. The Go side changes the layout viewport through
        // the protocol and calls this once per width, so ONE engine measures all
        // three instead of three engines measuring one each. Every step is
        // idempotent: the lens click lands on the lens already showing, and
        // reopening a drawer rebuilds the same rows.
        async function measureProbePage() {
          await nextFrame();
          await nextFrame();
          document.querySelector('#lens-group [data-lens-target="user-request"]:not([data-ur-cards])').click();
          await waitFor(function () {
            return !document.getElementById("user-request-lens").hidden;
          }, "By-UR lens");
          document.querySelector('[data-ur-activity="all"]').click();
          await waitFor(function () {
            return document.querySelector('.ur-group-detail[data-detail-id="UR-951"]') !== null;
          }, "both UR groups");

          var group = groupFor("UR-950");
          var headRow = group.querySelector(".ur-group-row");
          var strip = group.querySelector(".ur-summary");
          if (!strip) { throw new Error("the By UR group rendered no progress strip"); }
          var titleNode = group.querySelector(".ur-title");
          var detailButton = group.querySelector(".ur-group-detail");
          var headButton = group.querySelector(".ur-group-head");
          var metricNodes = Array.from(strip.querySelectorAll(".ur-summary-metric"));

          var titleRect = rectOf(titleNode);
          var detailRect = rectOf(detailButton);
          var groupRect = rectOf(group);
          var metricRects = metricNodes.map(rectOf);

          var overlappingPairs = [];
          var namedBoxes = [{ name: "ur-title", rect: titleRect }, { name: "ur-group-detail", rect: detailRect }];
          metricRects.forEach(function (rect, metricIndex) {
            namedBoxes.push({ name: "metric-" + metricIndex, rect: rect });
          });
          for (var firstIndex = 0; firstIndex < namedBoxes.length; firstIndex++) {
            for (var secondIndex = firstIndex + 1; secondIndex < namedBoxes.length; secondIndex++) {
              if (rectsIntersect(namedBoxes[firstIndex].rect, namedBoxes[secondIndex].rect)) {
                overlappingPairs.push(namedBoxes[firstIndex].name + " x " + namedBoxes[secondIndex].name);
              }
            }
          }
          var escapedMetrics = [];
          metricRects.forEach(function (rect, metricIndex) {
            if (!rectContains(groupRect, rect)) { escapedMetrics.push("metric-" + metricIndex); }
          });

          var textBoxes = Array.from(strip.querySelectorAll(".ur-summary-label, .ur-summary-value")).map(function (node) {
            return {
              className: node.className,
              text: node.textContent,
              scrollWidth: node.scrollWidth,
              clientWidth: node.clientWidth,
              contrast: contrastAgainstBody(getComputedStyle(node).color),
              color: getComputedStyle(node).color
            };
          });
          var headerMetricTexts = metricNodes.map(function (node) { return node.textContent; });

          // The same UR in the drawer: the two surfaces have to state the same
          // five figures in the page a person actually opens.
          detailButton.click();
          await waitFor(function () {
            return document.getElementById("detail-id").textContent === "UR-950";
          }, "UR-950 drawer");
          var drawerMetricTexts = Array.from(
            document.querySelectorAll('#detail-meta [data-ur-summary-metric]')
          ).map(function (node) {
            return node.previousElementSibling === null ? node.textContent : node.textContent;
          });

          // The 43-member UR: the id list starts expanded, and the input.md row
          // below it still has to be inside the panel.
          var largeDetailButton = document.querySelector('.ur-group-detail[data-detail-id="UR-951"]');
          largeDetailButton.click();
          await waitFor(function () {
            return document.getElementById("detail-id").textContent === "UR-951";
          }, "UR-951 drawer");
          var idListButton = Array.from(document.querySelectorAll("#detail-meta .detail-fold")).find(function (button) {
            return button.textContent.indexOf("REQ ids") === 0;
          });
          var largeMemberCount = document.querySelectorAll('#detail-meta [data-detail-kind="req"][data-detail-id]').length;
          var inputRow = Array.from(document.querySelectorAll("#detail-meta dt")).find(function (term) {
            return term.textContent.trim() === "input.md";
          });

          return JSON.stringify({
            locationHref: location.href,
            userAgent: navigator.userAgent,
            resolvedScheme: matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
            bodyBackground: getComputedStyle(document.body).backgroundColor,
            viewportWidth: document.documentElement.clientWidth,
            headRowRect: rectOf(headRow),
            stripRect: rectOf(strip),
            groupRect: groupRect,
            metricRects: metricRects,
            titleRect: titleRect,
            detailRect: detailRect,
            overlappingPairs: overlappingPairs,
            escapedMetrics: escapedMetrics,
            textBoxes: textBoxes,
            headerMetricTexts: headerMetricTexts,
            drawerMetricTexts: drawerMetricTexts,
            headTagName: headButton.tagName,
            headTabIndex: headButton.tabIndex,
            detailTagName: detailButton.tagName,
            detailTabIndex: detailButton.tabIndex,
            headPrecedesDetail:
              (headButton.compareDocumentPosition(detailButton) & Node.DOCUMENT_POSITION_FOLLOWING) !== 0,
            headFocusable: (headButton.focus(), document.activeElement === headButton),
            detailFocusable: (largeDetailButton.focus(), document.activeElement === largeDetailButton),
            largeMemberCount: largeMemberCount,
            idListExpanded: idListButton ? idListButton.getAttribute("aria-expanded") === "true" : false,
            inputRowRect: inputRow ? rectOf(inputRow) : { top: 0, right: 0, bottom: 0, left: 0 },
            drawerRect: rectOf(document.getElementById("detail-drawer")),
            consoleErrors: window.__queueKanbanProbeErrors.slice()
          });
        }

        // Resolves rather than rejects on failure: the caller reads one JSON
        // object either way, and a probe reporting its own error is easier to act
        // on than a protocol exception. The result node is still written so the
        // page carries its own evidence.
        window.__queueKanbanMeasure = function () {
          return measureProbePage().then(function (resultText) {
            resultNode.textContent = resultText;
            return resultText;
          }, function (probeError) {
            var failureText = JSON.stringify({
              locationHref: location.href,
              consoleErrors: window.__queueKanbanProbeErrors.concat([String((probeError && probeError.message) || probeError)])
            });
            resultNode.textContent = failureText;
            return failureText;
          });
        };
      })();
`

	probePage := string(indexBytes)
	clientScriptOpen := "    <script>\n"
	clientScriptOpenIndex := strings.LastIndex(probePage, clientScriptOpen)
	if clientScriptOpenIndex < 0 {
		t.Fatal("generated page has no client opening for the UR progress stub")
	}
	stubIndex := clientScriptOpenIndex + len(clientScriptOpen)
	probePage = probePage[:stubIndex] + errorStub + probePage[stubIndex:]
	clientCloseIndex := strings.LastIndex(probePage, "})();")
	if clientCloseIndex < 0 {
		t.Fatal("generated page has no client close for the UR progress probe")
	}
	clientCloseIndex += len("})();")
	probePage = probePage[:clientCloseIndex] + "\n" + probeScript + probePage[clientCloseIndex:]

	for _, theme := range []struct {
		name       string
		colourFlag string
	}{
		{name: "light", colourFlag: "--force-light-mode"},
		{name: "dark", colourFlag: "--force-dark-mode"},
	} {
		theme := theme
		t.Run(theme.name, func(t *testing.T) {
			// One engine per THEME measures all three widths. Without a
			// colour-scheme flag Chromium resolves prefers-color-scheme to light,
			// so every case would measure the light palette and nothing automated
			// would ever see the dark one — which on this board is the :root base
			// the light block overrides.
			session := startTrustedInputBrowserSession(
				t, "UR progress strip "+theme.name, siteDirectory, probePage,
				"--window-size=1280,1100", theme.colourFlag)
			defer session.closeBrowserSession()
			for _, viewportWidth := range []int{320, 768, 1280} {
				viewportWidth := viewportWidth
				t.Run(fmt.Sprintf("%dpx", viewportWidth), func(t *testing.T) {
					checkUserRequestProgressStrip(t, session, viewportWidth, theme.name)
				})
			}
		})
	}
}

// The assertions for one theme at one width. Split out of the loop above so the
// engine can be reused across widths without nesting the whole body two levels
// deeper.
func checkUserRequestProgressStrip(
	t *testing.T, session *trustedInputBrowserSession, viewportWidth int, wantScheme string,
) {
	t.Helper()
	resultJSON := measureUserRequestProgressAtWidth(t, session, viewportWidth)
	var result userRequestProgressProbeResult
	if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
		t.Fatalf("decode UR progress probe: %v\n%s", decodeError, resultJSON)
	}
	if len(result.ConsoleErrors) != 0 {
		t.Fatalf("UR progress probe browser errors: %q", result.ConsoleErrors)
	}
	if !strings.HasSuffix(result.LocationHref, "/"+browserProbePageFileName) {
		t.Fatalf("probe measured %q, not its own page", result.LocationHref)
	}
	// The build these numbers came from, recorded beside them: a green here
	// is evidence about this engine, not a compatibility claim.
	t.Logf("%s: %s, resolved scheme %s, body ground %s, viewport %.0f CSS px",
		wantScheme, result.UserAgent, result.ResolvedScheme, result.BodyBackground, result.ViewportWidth)

	// Without the colour-scheme flag Chromium resolves prefers-color-scheme
	// to light and the dark palette is never measured by anything.
	// The width the numbers below actually describe.
	if int(result.ViewportWidth) != viewportWidth {
		t.Fatalf("the engine laid the page out at %.0f CSS px, not the %d this case names; "+
			"every measurement below would describe a width nobody asked about",
			result.ViewportWidth, viewportWidth)
	}
	if result.ResolvedScheme != wantScheme {
		t.Fatalf("engine resolved the %s case as %q; the dark palette would never be measured",
			wantScheme, result.ResolvedScheme)
	}

	// The strip is its own row under the head, at every width.
	if result.StripRect.Top < result.HeadRowRect.Bottom-0.5 {
		t.Fatalf("summary strip top %.1f sits above the head row bottom %.1f — the metrics are sharing the title's line",
			result.StripRect.Top, result.HeadRowRect.Bottom)
	}
	if len(result.MetricRects) != 5 {
		t.Fatalf("the strip laid out %d metric boxes, want the five figures the request names", len(result.MetricRects))
	}
	if len(result.OverlappingPairs) != 0 {
		t.Fatalf("boxes collide at %d CSS px: %v", viewportWidth, result.OverlappingPairs)
	}
	if len(result.EscapedMetrics) != 0 {
		t.Fatalf("metric boxes escaped their .ur-group at %d CSS px: %v", viewportWidth, result.EscapedMetrics)
	}

	// Nothing is clipped, and everything clears the contrast floor against
	// the ground it is actually painted on.
	if len(result.TextBoxes) != 10 {
		t.Fatalf("measured %d label/value boxes, want a label and a value for each of the five metrics",
			len(result.TextBoxes))
	}
	for _, textBox := range result.TextBoxes {
		if textBox.ScrollWidth > textBox.ClientWidth+1 {
			t.Errorf("%s %q is clipped at %d CSS px (scroll %.1f > client %.1f)",
				textBox.ClassName, textBox.Text, viewportWidth, textBox.ScrollWidth, textBox.ClientWidth)
		}
		if textBox.Contrast < 4.5 {
			t.Errorf("%s %q measures %.2f:1 in %s against %s, below the 4.5:1 floor",
				textBox.ClassName, textBox.Text, textBox.Contrast, result.ResolvedScheme, result.BodyBackground)
		}
	}

	// Both surfaces, in the page a person opens.
	if len(result.DrawerMetricTexts) != 5 {
		t.Fatalf("the UR drawer rendered %d summary values, want five", len(result.DrawerMetricTexts))
	}
	for metricIndex, headerText := range result.HeaderMetricTexts {
		if !strings.HasSuffix(headerText, result.DrawerMetricTexts[metricIndex]) {
			t.Fatalf("header metric %d reads %q while the drawer reads %q; the two surfaces render one rollup",
				metricIndex, headerText, result.DrawerMetricTexts[metricIndex])
		}
	}

	// Keyboard reachability. Trusted Tab is a browser default action this
	// lane does not synthesize; being real focusable <button> elements in
	// fold-then-Details document order is what delivers the order.
	if result.HeadTagName != "BUTTON" || result.DetailTagName != "BUTTON" {
		t.Fatalf("fold control is <%s> and Details is <%s>; both must be real buttons",
			result.HeadTagName, result.DetailTagName)
	}
	if result.HeadTabIndex < 0 || result.DetailTabIndex < 0 {
		t.Fatalf("tabIndex fold=%d details=%d; a negative index takes the control out of the tab ring",
			result.HeadTabIndex, result.DetailTabIndex)
	}
	if !result.HeadPrecedesDetil || !result.HeadFocusable || !result.DetailFocusable {
		t.Fatalf("fold precedes Details=%v, fold focusable=%v, Details focusable=%v",
			result.HeadPrecedesDetil, result.HeadFocusable, result.DetailFocusable)
	}

	// The drawer with a 43-member UR: the reason the id list was capped.
	if result.LargeMemberCount != 43 || !result.IdListExpanded {
		t.Fatalf("large UR drawer listed %d ids with the fold expanded=%v, want 43 and an open list",
			result.LargeMemberCount, result.IdListExpanded)
	}
	if result.InputRowRect.Bottom > result.DrawerRect.Bottom+0.5 || result.InputRowRect.Top < result.DrawerRect.Top-0.5 {
		t.Errorf("the input.md row sits at %.1f-%.1f outside the drawer panel %.1f-%.1f at %d CSS px; "+
			"the summary rows plus an expanded 43-id list pushed it out again",
			result.InputRowRect.Top, result.InputRowRect.Bottom,
			result.DrawerRect.Top, result.DrawerRect.Bottom, viewportWidth)
	}
}
