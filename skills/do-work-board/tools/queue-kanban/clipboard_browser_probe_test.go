package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type boardColumnCopyProbeResult struct {
	LocationHref           string            `json:"locationHref"`
	ButtonCount            int               `json:"buttonCount"`
	InitialDisabled        map[string]bool   `json:"initialDisabled"`
	InitialPendingIds      []string          `json:"initialPendingIds"`
	InitialPendingPayload  string            `json:"initialPendingPayload"`
	ClaimedIds             []string          `json:"claimedIds"`
	ClaimedPayload         string            `json:"claimedPayload"`
	FilteredPendingIds     []string          `json:"filteredPendingIds"`
	FilteredPendingPayload string            `json:"filteredPendingPayload"`
	FilteredDisabled       map[string]bool   `json:"filteredDisabled"`
	RecentlyDone24Ids      []string          `json:"recentlyDone24Ids"`
	RecentlyDone24Payload  string            `json:"recentlyDone24Payload"`
	RecentlyDone48Ids      []string          `json:"recentlyDone48Ids"`
	RecentlyDone48Payload  string            `json:"recentlyDone48Payload"`
	SuccessfulFeedback     map[string]string `json:"successfulFeedback"`
	FailedFeedback         string            `json:"failedFeedback"`
	ClipboardWrites        []string          `json:"clipboardWrites"`
	ConsoleErrors          []string          `json:"consoleErrors"`
}

func boardColumnCopyFixtureTicket(
	requestID string, title string, status string, domain string, bodyMarker string,
) *RequestTicket {
	return &RequestTicket{
		RequestId:           requestID,
		Title:               title,
		Status:              status,
		OriginalStatus:      status,
		Domain:              domain,
		OriginalDomain:      domain,
		TreeSection:         "queue",
		CreatedAt:           "2026-08-20T12:00:00Z",
		FrontmatterMarkdown: "---\nid: " + requestID + "\ntitle: '" + title + "'\nstatus: " + status + "\n---\n",
		BodyMarkdown:        "# " + title + "\n\n" + bodyMarker + "\n",
	}
}

// REQ-367 drives the real generated Board rather than calling a composition
// helper directly. The clicked column is therefore the authority: filters,
// Pending's Ready/Waiting grouping, and the Recently Done window all change the
// rendered cards before the clipboard code reads their ids back in DOM order.
func TestBrowserBehaviorBoardColumnCopyAll(t *testing.T) {
	probeNow := time.Now().UTC()
	pendingBeta := boardColumnCopyFixtureTicket("REQ-102", "Pending beta", "pending", "docs", "BETA-BODY")
	pendingAlpha := boardColumnCopyFixtureTicket("REQ-101", "Pending alpha", "pending", "frontend", "ALPHA-BODY")
	pendingWaiting := boardColumnCopyFixtureTicket("REQ-103", "Pending waiting", "pending", "docs", "WAITING-BODY")
	claimed := boardColumnCopyFixtureTicket("REQ-201", "Claimed only", "claimed", "frontend", "CLAIMED-BODY")
	recent24 := boardColumnCopyFixtureTicket("REQ-401", "Done recently", "completed", "docs", "DONE-24-BODY")
	recent48 := boardColumnCopyFixtureTicket("REQ-402", "Done yesterday", "completed", "frontend", "DONE-48-BODY")
	tooOld := boardColumnCopyFixtureTicket("REQ-403", "Done too old", "completed", "docs", "DONE-OLD-BODY")
	// REQ-379: bodies carry ticket mentions so the copied payload exercises the
	// annotation end to end. REQ-101 sits INSIDE the Pending payload (expanded
	// inline, dropped from the glossary), REQ-201 sits outside it, and REQ-9999
	// answers to no record at all. Claimed and Recently Done keep mention-free
	// bodies, which is what pins "no references, no appendix".
	pendingBeta.BodyMarkdown += "\nSee REQ-101 and REQ-9999, then REQ-101 again.\n"
	pendingAlpha.BodyMarkdown += "\nBlocked by REQ-201.\n"
	recent24.CompletionTime = probeNow.Add(-2 * time.Hour)
	recent48.CompletionTime = probeNow.Add(-36 * time.Hour)
	tooOld.CompletionTime = probeNow.Add(-8 * 24 * time.Hour)
	for _, terminalTicket := range []*RequestTicket{recent24, recent48, tooOld} {
		terminalTicket.TreeSection = "archive"
		terminalTicket.CompletedAt = terminalTicket.CompletionTime.Format(time.RFC3339)
	}

	fixtureBoard := &Board{
		GeneratedAt: probeNow,
		ProjectName: "REQ-367 Board column copy probe",
		AllRequests: []*RequestTicket{
			pendingAlpha, pendingBeta, pendingWaiting, claimed, recent24, recent48, tooOld,
		},
		Columns: BoardColumns{
			Pending:        []*RequestTicket{pendingBeta, pendingAlpha, pendingWaiting},
			PendingReady:   []*RequestTicket{pendingBeta, pendingAlpha},
			PendingWaiting: []*RequestTicket{pendingWaiting},
			Claimed:        []*RequestTicket{claimed},
			RecentlyDone:   []*RequestTicket{recent24},
		},
		Calendar: []CalendarEntry{
			{RequestId: recent24.RequestId, Status: recent24.Status, EntryTime: recent24.CompletionTime},
			{RequestId: recent48.RequestId, Status: recent48.Status, EntryTime: recent48.CompletionTime},
			{RequestId: tooOld.RequestId, Status: tooOld.Status, EntryTime: tooOld.CompletionTime},
		},
	}

	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate Board column copy fixture: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read Board column copy fixture: %v", readError)
	}

	clipboardStub := `
      window.__queueKanbanClipboardWrites = [];
      window.__queueKanbanClipboardReject = false;
      window.__queueKanbanProbeErrors = [];
      window.addEventListener("error", function (event) {
        window.__queueKanbanProbeErrors.push(event.message || "window error");
      });
      var queueKanbanOriginalConsoleError = console.error;
      console.error = function () {
        window.__queueKanbanProbeErrors.push(Array.prototype.join.call(arguments, " "));
        queueKanbanOriginalConsoleError.apply(console, arguments);
      };
      Object.defineProperty(navigator, "clipboard", {
        configurable: true,
        value: {
          writeText: function (clipboardText) {
            if (window.__queueKanbanClipboardReject) {
              return Promise.reject(new Error("probe clipboard rejection"));
            }
            window.__queueKanbanClipboardWrites.push(clipboardText);
            return Promise.resolve();
          }
        }
      });
`
	probeScript := `
      (function () {
        var resultNode = document.createElement("pre");
        resultNode.id = "` + browserProbeResultElementId + `";
        document.body.appendChild(resultNode);

        function copyButton(columnKey) {
          return document.querySelector('[data-copy-column="' + columnKey + '"]');
        }
        function visibleIds(columnKey) {
          var columnNode = copyButton(columnKey).closest(".kanban-column");
          return Array.from(columnNode.querySelectorAll('.req-card[data-detail-id]')).map(function (card) {
            return card.dataset.detailId;
          });
        }
        function disabledState() {
          var state = {};
          document.querySelectorAll("[data-copy-column]").forEach(function (button) {
            state[button.dataset.copyColumn] = button.disabled;
          });
          return state;
        }
        function waitFor(predicate, failureLabel) {
          return new Promise(function (resolve, reject) {
            var attempts = 0;
            function poll() {
              if (predicate()) {
                resolve();
                return;
              }
              attempts += 1;
              if (attempts > 200) {
                reject(new Error("timed out waiting for " + failureLabel));
                return;
              }
              setTimeout(poll, 10);
            }
            poll();
          });
        }
        function copyColumn(columnKey) {
          var button = copyButton(columnKey);
          var previousWriteCount = window.__queueKanbanClipboardWrites.length;
          button.click();
          return waitFor(function () {
            return window.__queueKanbanClipboardWrites.length === previousWriteCount + 1 &&
              button.textContent === "Copied ✓";
          }, columnKey + " clipboard write").then(function () {
            return {
              payload: window.__queueKanbanClipboardWrites[previousWriteCount],
              feedback: button.textContent
            };
          });
        }
        function dispatchValue(node, nextValue, eventName) {
          node.value = nextValue;
          node.dispatchEvent(new Event(eventName, {bubbles: true}));
        }

        async function runProbe() {
          var result = {
            locationHref: location.href,
            buttonCount: document.querySelectorAll("[data-copy-column]").length,
            initialDisabled: disabledState(),
            initialPendingIds: visibleIds("pending"),
            successfulFeedback: {},
            consoleErrors: window.__queueKanbanProbeErrors
          };

          var pendingCopy = await copyColumn("pending");
          result.initialPendingPayload = pendingCopy.payload;
          result.successfulFeedback.pending = pendingCopy.feedback;

          result.claimedIds = visibleIds("claimed");
          var claimedCopy = await copyColumn("claimed");
          result.claimedPayload = claimedCopy.payload;
          result.successfulFeedback.claimed = claimedCopy.feedback;

          dispatchValue(document.getElementById("filter-search"), "beta", "input");
          dispatchValue(document.getElementById("filter-domain"), "docs", "change");
          dispatchValue(document.getElementById("filter-status"), "pending", "change");
          result.filteredPendingIds = visibleIds("pending");
          result.filteredDisabled = disabledState();
          var filteredCopy = await copyColumn("pending");
          result.filteredPendingPayload = filteredCopy.payload;
          result.successfulFeedback.filtered = filteredCopy.feedback;

          document.getElementById("filter-clear").click();
          result.recentlyDone24Ids = visibleIds("recentlyDone");
          var recent24Copy = await copyColumn("recentlyDone");
          result.recentlyDone24Payload = recent24Copy.payload;
          result.successfulFeedback.recent24 = recent24Copy.feedback;

          document.querySelector('[data-window-hours="48"]').click();
          result.recentlyDone48Ids = visibleIds("recentlyDone");
          var recent48Copy = await copyColumn("recentlyDone");
          result.recentlyDone48Payload = recent48Copy.payload;
          result.successfulFeedback.recent48 = recent48Copy.feedback;

          window.__queueKanbanClipboardReject = true;
          document.execCommand = function () { return false; };
          copyButton("claimed").click();
          await waitFor(function () {
            return copyButton("claimed").textContent === "Copy failed";
          }, "failed clipboard feedback");
          result.failedFeedback = copyButton("claimed").textContent;
          result.clipboardWrites = window.__queueKanbanClipboardWrites.slice();
          result.consoleErrors = window.__queueKanbanProbeErrors.slice();
          resultNode.textContent = JSON.stringify(result);
        }

        runProbe().catch(function (probeError) {
          resultNode.textContent = JSON.stringify({
            locationHref: location.href,
            consoleErrors: window.__queueKanbanProbeErrors.concat([String(probeError && probeError.message || probeError)])
          });
        });
      })();
`

	probePage := string(indexBytes)
	clientScriptOpen := "    <script>\n"
	clientScriptOpenIndex := strings.LastIndex(probePage, clientScriptOpen)
	if clientScriptOpenIndex < 0 {
		t.Fatal("generated page has no client opening for clipboard stub")
	}
	clipboardStubIndex := clientScriptOpenIndex + len(clientScriptOpen)
	probePage = probePage[:clipboardStubIndex] + clipboardStub + probePage[clipboardStubIndex:]
	clientCloseIndex := strings.LastIndex(probePage, "})();")
	if clientCloseIndex < 0 {
		t.Fatal("generated page has no client close for clipboard probe")
	}
	clientCloseIndex += len("})();")
	probePage = probePage[:clientCloseIndex] + "\n" + probeScript + probePage[clientCloseIndex:]

	resultJSON := runBrowserBehaviorProbeInDirectory(
		t, "Board column copy all", siteDirectory, probePage, "--virtual-time-budget=30000",
	)
	var result boardColumnCopyProbeResult
	if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
		t.Fatalf("decode Board column copy probe: %v\n%s", decodeError, resultJSON)
	}
	if !strings.HasSuffix(result.LocationHref, "/"+browserProbePageFileName) {
		t.Fatalf("Board column copy measured %q, not its probe page", result.LocationHref)
	}
	if len(result.ConsoleErrors) != 0 {
		t.Fatalf("Board column copy browser errors: %q", result.ConsoleErrors)
	}
	if result.ButtonCount != 4 {
		t.Fatalf("column copy button count = %d, want exactly four flat Board columns", result.ButtonCount)
	}
	if !result.InitialDisabled["needsInputOrBlocked"] || result.InitialDisabled["pending"] ||
		result.InitialDisabled["claimed"] || result.InitialDisabled["recentlyDone"] {
		t.Fatalf("initial disabled states = %#v, want only empty Needs input/Blocked disabled", result.InitialDisabled)
	}

	wantPendingIds := []string{"REQ-102", "REQ-101", "REQ-103"}
	if !reflect.DeepEqual(result.InitialPendingIds, wantPendingIds) {
		t.Fatalf("initial Pending display order = %v, want Ready then Waiting order %v", result.InitialPendingIds, wantPendingIds)
	}
	// Written out rather than recomputed: every FrontmatterMarkdown below is the
	// fixture's own bytes, so a fence that gained an expansion fails here, and the
	// annotated bodies are literals so the expansion shape cannot drift silently.
	annotatedPendingBetaBody := "# Pending beta\n\nBETA-BODY\n\nSee REQ-101 (-> Pending alpha) and REQ-9999, then REQ-101 again.\n"
	annotatedPendingAlphaBody := "# Pending alpha\n\nALPHA-BODY\n\nBlocked by REQ-201 (-> Claimed only).\n"
	wantPendingPayload :=
		pendingBeta.FrontmatterMarkdown + annotatedPendingBetaBody +
			pendingAlpha.FrontmatterMarkdown + annotatedPendingAlphaBody +
			pendingWaiting.FrontmatterMarkdown + pendingWaiting.BodyMarkdown +
			"\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
			"- REQ-9999 — not found in this queue\n" +
			"- REQ-201 — Claimed only (claimed)\n"
	if result.InitialPendingPayload != wantPendingPayload {
		t.Fatalf("initial Pending clipboard payload changed order/content:\n got %q\nwant %q", result.InitialPendingPayload, wantPendingPayload)
	}
	// A mention-free column copies its files and gains no appendix at all.
	if !reflect.DeepEqual(result.ClaimedIds, []string{"REQ-201"}) ||
		result.ClaimedPayload != claimed.FrontmatterMarkdown+claimed.BodyMarkdown {
		t.Fatalf("single-card Claimed copy = ids %v payload %q", result.ClaimedIds, result.ClaimedPayload)
	}
	// The same body, copied from a filtered column that no longer carries
	// REQ-101: the exclusion tracks the payload, so REQ-101 now earns a line.
	wantFilteredPendingPayload := pendingBeta.FrontmatterMarkdown + annotatedPendingBetaBody +
		"\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- REQ-101 — Pending alpha (pending)\n" +
		"- REQ-9999 — not found in this queue\n"
	if !reflect.DeepEqual(result.FilteredPendingIds, []string{"REQ-102"}) ||
		result.FilteredPendingPayload != wantFilteredPendingPayload {
		t.Fatalf("filtered Pending copy = ids %v payload %q, want %q",
			result.FilteredPendingIds, result.FilteredPendingPayload, wantFilteredPendingPayload)
	}
	for _, hiddenColumnKey := range []string{"claimed", "needsInputOrBlocked", "recentlyDone"} {
		if !result.FilteredDisabled[hiddenColumnKey] {
			t.Errorf("filtered-empty %s copy button remained enabled: %#v", hiddenColumnKey, result.FilteredDisabled)
		}
	}
	if !reflect.DeepEqual(result.RecentlyDone24Ids, []string{"REQ-401"}) ||
		result.RecentlyDone24Payload != recent24.FrontmatterMarkdown+recent24.BodyMarkdown {
		t.Fatalf("24h Recently Done copy = ids %v payload %q", result.RecentlyDone24Ids, result.RecentlyDone24Payload)
	}
	wantRecent48Ids := []string{"REQ-401", "REQ-402"}
	wantRecent48Payload := recent24.FrontmatterMarkdown + recent24.BodyMarkdown +
		recent48.FrontmatterMarkdown + recent48.BodyMarkdown
	if !reflect.DeepEqual(result.RecentlyDone48Ids, wantRecent48Ids) || result.RecentlyDone48Payload != wantRecent48Payload {
		t.Fatalf("48h Recently Done copy = ids %v payload %q, want %v / %q",
			result.RecentlyDone48Ids, result.RecentlyDone48Payload, wantRecent48Ids, wantRecent48Payload)
	}
	if len(result.ClipboardWrites) != 5 {
		t.Fatalf("successful clipboard writes = %d, want five non-vacuous copy actions", len(result.ClipboardWrites))
	}
	for copySurface, feedbackText := range result.SuccessfulFeedback {
		if feedbackText != "Copied ✓" {
			t.Errorf("%s success feedback = %q", copySurface, feedbackText)
		}
	}
	if len(result.SuccessfulFeedback) != 5 || result.FailedFeedback != "Copy failed" {
		t.Fatalf("feedback coverage = success %#v failure %q", result.SuccessfulFeedback, result.FailedFeedback)
	}
}
