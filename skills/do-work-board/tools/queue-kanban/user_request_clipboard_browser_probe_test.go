package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type userRequestCopyAllProbeResult struct {
	LocationHref               string            `json:"locationHref"`
	DrawerButtonLabels         []string          `json:"drawerButtonLabels"`
	GroupedRequestIds          []string          `json:"groupedRequestIds"`
	PlainCopyPayload           string            `json:"plainCopyPayload"`
	CopyAllPayload             string            `json:"copyAllPayload"`
	MissingMemberFeedback      string            `json:"missingMemberFeedback"`
	RequestDetailCopyAllHidden bool              `json:"requestDetailCopyAllHidden"`
	NoMemberRequestIds         []string          `json:"noMemberRequestIds"`
	NoMemberPayload            string            `json:"noMemberPayload"`
	SyntheticUrFeedback        string            `json:"syntheticUrFeedback"`
	SuccessfulFeedback         map[string]string `json:"successfulFeedback"`
	ClipboardWrites            []string          `json:"clipboardWrites"`
	ConsoleErrors              []string          `json:"consoleErrors"`
}

func userRequestCopyFixture(
	userRequestID string, title string, requestArrayText string, bodyMarker string,
) *UserRequestTicket {
	return &UserRequestTicket{
		UserRequestId:       userRequestID,
		Title:               title,
		FrontmatterMarkdown: "---\nid: " + userRequestID + "\ntitle: '" + title + "'\nrequests: " + requestArrayText + "\n---\n",
		BodyMarkdown:        "# " + title + "\n\n" + bodyMarker + "\n",
		InputFilePresent:    true,
	}
}

// REQ-368 drives the generated Board and opens UR details through the real By-UR
// UI. The fixture's capture-time requests array names only REQ-710, while three
// REQs point upward from queue, working, and archive. Only the board's grouped
// userRequest.requestIds set can therefore produce the displayed and copied set.
func TestBrowserBehaviorUserRequestCopyAllIncludesGroupedRequests(t *testing.T) {
	groupedUserRequest := userRequestCopyFixture("UR-700", "Grouped work", "[REQ-710]", "GROUPED-UR-BODY")
	noMemberUserRequest := userRequestCopyFixture("UR-800", "No members", "[]", "NO-MEMBER-UR-BODY")
	// REQ-379: the same body copied through two buttons with different payloads.
	// REQ-710 is a grouped member, so Copy all excludes it while the plain Copy —
	// whose payload is the UR file alone — still owes the reader its title.
	groupedUserRequest.BodyMarkdown += "\nCovers REQ-710 and UR-800, plus REQ-9999.\n"

	queueRequest := boardColumnCopyFixtureTicket("REQ-702", "Queue member", "pending", "frontend", "QUEUE-MEMBER-BODY")
	queueRequest.UserRequestId = groupedUserRequest.UserRequestId
	queueRequest.TreeSection = "queue"
	workingRequest := boardColumnCopyFixtureTicket("REQ-710", "Working member", "claimed", "backend", "WORKING-MEMBER-BODY")
	workingRequest.UserRequestId = groupedUserRequest.UserRequestId
	workingRequest.TreeSection = "working"
	archiveRequest := boardColumnCopyFixtureTicket("REQ-725", "Archive member", "completed", "docs", "ARCHIVE-MEMBER-BODY")
	archiveRequest.UserRequestId = groupedUserRequest.UserRequestId
	archiveRequest.TreeSection = "archive"
	syntheticMember := boardColumnCopyFixtureTicket("REQ-901", "Synthetic UR member", "pending", "general", "SYNTHETIC-MEMBER-BODY")
	syntheticMember.UserRequestId = "UR-900"
	syntheticMember.TreeSection = "queue"

	fixtureBoard := &Board{
		ProjectName: "REQ-368 UR copy-all probe",
		AllRequests: []*RequestTicket{
			workingRequest, archiveRequest, queueRequest, syntheticMember,
		},
		UserRequests: []*UserRequestTicket{groupedUserRequest, noMemberUserRequest},
		UserRequestsById: map[string]*UserRequestTicket{
			groupedUserRequest.UserRequestId:  groupedUserRequest,
			noMemberUserRequest.UserRequestId: noMemberUserRequest,
		},
		Columns: BoardColumns{
			Pending:      []*RequestTicket{queueRequest, syntheticMember},
			PendingReady: []*RequestTicket{queueRequest, syntheticMember},
			Claimed:      []*RequestTicket{workingRequest},
		},
	}
	linkRequestsToUserRequests(fixtureBoard)
	for _, userRequest := range fixtureBoard.UserRequests {
		sortRequestIdList(userRequest.RequestIds)
	}

	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate UR copy-all fixture: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read UR copy-all fixture: %v", readError)
	}

	clipboardStub := `
      window.__queueKanbanClipboardWrites = [];
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

        function openUserRequest(userRequestId) {
          // The head folds the group in both readings since REQ-486, so the
          // drawer trigger is the sibling Details button, not the head.
          var trigger = document.querySelector('.ur-group-detail[data-detail-id="' + userRequestId + '"]');
          if (!trigger) {
            throw new Error("no By-UR trigger for " + userRequestId);
          }
          trigger.click();
          return waitFor(function () {
            return document.getElementById("detail-kind").textContent === "UR" &&
              document.getElementById("detail-id").textContent === userRequestId;
          }, userRequestId + " drawer");
        }

        function groupedRequestIds() {
          return Array.from(
            document.querySelectorAll('#detail-meta [data-detail-kind="req"][data-detail-id]')
          ).map(function (link) {
            return link.dataset.detailId;
          });
        }

        function copyWithSuccess(button, failureLabel) {
          var previousWriteCount = window.__queueKanbanClipboardWrites.length;
          button.click();
          return waitFor(function () {
            return window.__queueKanbanClipboardWrites.length === previousWriteCount + 1 &&
              button.textContent === "Copied ✓";
          }, failureLabel).then(function () {
            return {
              payload: window.__queueKanbanClipboardWrites[previousWriteCount],
              feedback: button.textContent
            };
          });
        }

        async function runProbe() {
          document.querySelector('#lens-group [data-lens-target="user-request"]:not([data-ur-cards])').click();
          await waitFor(function () {
            return !document.getElementById("user-request-lens").hidden;
          }, "By-UR lens");

          await openUserRequest("UR-700");
          var plainButton = document.getElementById("detail-copy");
          var copyAllButton = document.getElementById("detail-copy-all");
          var result = {
            locationHref: location.href,
            drawerButtonLabels: Array.from(document.querySelectorAll(".detail-head-actions button")).map(function (button) {
              return button.textContent.trim();
            }),
            groupedRequestIds: groupedRequestIds(),
            successfulFeedback: {},
            consoleErrors: window.__queueKanbanProbeErrors
          };

          var plainCopy = await copyWithSuccess(plainButton, "plain UR copy");
          result.plainCopyPayload = plainCopy.payload;
          result.successfulFeedback.plain = plainCopy.feedback;

          var copyAll = await copyWithSuccess(copyAllButton, "UR copy all");
          result.copyAllPayload = copyAll.payload;
          result.successfulFeedback.copyAll = copyAll.feedback;

          var missingMemberRaw = window.queueKanbanBoardMarkdownData.requests["REQ-725"];
          delete window.queueKanbanBoardMarkdownData.requests["REQ-725"];
          var writeCountBeforeFailure = window.__queueKanbanClipboardWrites.length;
          copyAllButton.click();
          await waitFor(function () {
            return copyAllButton.textContent === "Copy failed";
          }, "missing grouped member failure");
          result.missingMemberFeedback = copyAllButton.textContent;
          if (window.__queueKanbanClipboardWrites.length !== writeCountBeforeFailure) {
            throw new Error("missing grouped member wrote a partial clipboard payload");
          }
          window.queueKanbanBoardMarkdownData.requests["REQ-725"] = missingMemberRaw;

          document.querySelector('#detail-meta [data-detail-kind="req"][data-detail-id="REQ-702"]').click();
          await waitFor(function () {
            return document.getElementById("detail-kind").textContent === "REQ";
          }, "REQ drawer");
          result.requestDetailCopyAllHidden = copyAllButton.hidden;

          document.querySelector('[data-ur-activity="all"]').click();
          await openUserRequest("UR-800");
          result.noMemberRequestIds = groupedRequestIds();
          var noMemberCopy = await copyWithSuccess(copyAllButton, "no-member UR copy all");
          result.noMemberPayload = noMemberCopy.payload;
          result.successfulFeedback.noMember = noMemberCopy.feedback;

          await openUserRequest("UR-900");
          var writeCountBeforeSyntheticFailure = window.__queueKanbanClipboardWrites.length;
          copyAllButton.click();
          await waitFor(function () {
            return copyAllButton.textContent === "Copy failed";
          }, "synthetic UR failure");
          result.syntheticUrFeedback = copyAllButton.textContent;
          if (window.__queueKanbanClipboardWrites.length !== writeCountBeforeSyntheticFailure) {
            throw new Error("synthetic UR wrote a REQ-only clipboard payload");
          }

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
		t.Fatal("generated page has no client opening for UR clipboard stub")
	}
	clipboardStubIndex := clientScriptOpenIndex + len(clientScriptOpen)
	probePage = probePage[:clipboardStubIndex] + clipboardStub + probePage[clipboardStubIndex:]
	clientCloseIndex := strings.LastIndex(probePage, "})();")
	if clientCloseIndex < 0 {
		t.Fatal("generated page has no client close for UR clipboard probe")
	}
	clientCloseIndex += len("})();")
	probePage = probePage[:clientCloseIndex] + "\n" + probeScript + probePage[clientCloseIndex:]

	resultJSON := runBrowserBehaviorProbeInDirectory(
		t, "UR copy all", siteDirectory, probePage, "--virtual-time-budget=30000",
	)
	var result userRequestCopyAllProbeResult
	if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
		t.Fatalf("decode UR copy-all probe: %v\n%s", decodeError, resultJSON)
	}
	if !strings.HasSuffix(result.LocationHref, "/"+browserProbePageFileName) {
		t.Fatalf("UR copy all measured %q, not its probe page", result.LocationHref)
	}
	if len(result.ConsoleErrors) != 0 {
		t.Fatalf("UR copy-all browser errors: %q", result.ConsoleErrors)
	}
	if !reflect.DeepEqual(result.DrawerButtonLabels, []string{"Copy", "Copy all", "Close"}) {
		t.Fatalf("UR drawer controls = %q, want separate Copy / Copy all / Close", result.DrawerButtonLabels)
	}
	wantGroupedRequestIds := []string{"REQ-702", "REQ-710", "REQ-725"}
	if !reflect.DeepEqual(result.GroupedRequestIds, wantGroupedRequestIds) {
		t.Fatalf("displayed grouped REQs = %v, want all-tree order %v", result.GroupedRequestIds, wantGroupedRequestIds)
	}
	// FrontmatterMarkdown is spliced in verbatim on purpose: it carries
	// "requests: [REQ-710]", so a fence the annotator touched fails right here.
	annotatedGroupedBody := "# Grouped work\n\nGROUPED-UR-BODY\n\n" +
		"Covers REQ-710 (-> Working member) and UR-800 (-> No members), plus REQ-9999.\n"
	wantPlainCopyPayload := groupedUserRequest.FrontmatterMarkdown + annotatedGroupedBody +
		"\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- REQ-710 — Working member (claimed)\n" +
		"- UR-800 — No members (user request)\n" +
		"- REQ-9999 — not found in this queue\n"
	if result.PlainCopyPayload != wantPlainCopyPayload {
		t.Fatalf("plain UR Copy changed payload:\n got %q\nwant %q", result.PlainCopyPayload, wantPlainCopyPayload)
	}
	// Copy all carries REQ-710's own file, so its title is already in the paste
	// and the appendix drops it — the same body, a different reference list.
	wantCopyAllPayload := groupedUserRequest.FrontmatterMarkdown + annotatedGroupedBody +
		queueRequest.FrontmatterMarkdown + queueRequest.BodyMarkdown +
		workingRequest.FrontmatterMarkdown + workingRequest.BodyMarkdown +
		archiveRequest.FrontmatterMarkdown + archiveRequest.BodyMarkdown +
		"\n---\n\n" + referencedRequestsGlossaryHeading + "\n\n" +
		"- UR-800 — No members (user request)\n" +
		"- REQ-9999 — not found in this queue\n"
	if result.CopyAllPayload != wantCopyAllPayload {
		t.Fatalf("UR Copy all changed membership/order/content:\n got %q\nwant %q", result.CopyAllPayload, wantCopyAllPayload)
	}
	if result.MissingMemberFeedback != "Copy failed" {
		t.Fatalf("missing grouped payload feedback = %q, want Copy failed", result.MissingMemberFeedback)
	}
	if !result.RequestDetailCopyAllHidden {
		t.Fatal("UR Copy all remained visible in a REQ detail drawer")
	}
	if len(result.NoMemberRequestIds) != 0 {
		t.Fatalf("no-member UR displayed grouped REQs: %v", result.NoMemberRequestIds)
	}
	noMemberRaw := noMemberUserRequest.FrontmatterMarkdown + noMemberUserRequest.BodyMarkdown
	if result.NoMemberPayload != noMemberRaw {
		t.Fatalf("no-member Copy all = %q, want its UR payload only %q", result.NoMemberPayload, noMemberRaw)
	}
	if result.SyntheticUrFeedback != "Copy failed" {
		t.Fatalf("synthetic UR feedback = %q, want Copy failed", result.SyntheticUrFeedback)
	}
	if len(result.ClipboardWrites) != 3 {
		t.Fatalf("successful clipboard writes = %d, want plain + grouped + no-member", len(result.ClipboardWrites))
	}
	if !reflect.DeepEqual(result.SuccessfulFeedback, map[string]string{
		"plain": "Copied ✓", "copyAll": "Copied ✓", "noMember": "Copied ✓",
	}) {
		t.Fatalf("successful feedback = %#v", result.SuccessfulFeedback)
	}
}
