package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
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
// The existing DevTools pipe speaks directly to the browser binary; no package
// manager or external driver is required. Measurement probes read their completed
// result node over that same channel, and interaction probes can dispatch trusted
// input. Chrome 151.0.7922.174 on this host printed --dump-dom output but did not exit, so
// process termination cannot be the measurement-completion signal (REQ-375).
// Readiness is a bounded wait for the result node, not a fixed sleep. Every probe
// still asserts its own result and contributes to the strict lane's zero-probe guard.
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
	// Virtual-time budgets belonged to the dump-DOM command. The shared protocol
	// transport waits for explicit completion on the real page clock instead.
	browserFlags := make([]string, 0, len(extraFlags))
	for _, browserFlag := range extraFlags {
		if !strings.HasPrefix(browserFlag, "--virtual-time-budget=") {
			browserFlags = append(browserFlags, browserFlag)
		}
	}
	session := startTrustedInputBrowserSession(t, probeName, siteDirectory, pageHTML, browserFlags...)
	defer session.closeBrowserSession()
	session.waitForPageCondition(t, "completed browser result node",
		`document.getElementById("`+browserProbeResultElementId+`") && document.getElementById("`+browserProbeResultElementId+`").textContent.trim() !== ""`)
	var resultText string
	session.decodeResult(t, "completed probe result", session.evaluateInPage(t,
		`document.getElementById("`+browserProbeResultElementId+`").textContent`), &resultText)
	resultText = strings.TrimSpace(resultText)
	if !strings.HasPrefix(resultText, "{") {
		t.Fatalf("%s result node holds %q, not the JSON object the contract expects", probeName, resultText)
	}
	return []byte(resultText)
}

// ---- the trusted-input transport -------------------------------------------
//
// This protocol session originally enabled trusted input (REQ-341); measurement
// probes now share it. A page script alone can dispatch only synthetic events,
// which cannot establish capture on a pointer id the engine did not issue.
// Keeping the browser alive lets the tests drive the real click/capture path.
//
// WHAT THIS IS. Chromium's `--remote-debugging-pipe` moves the DevTools Protocol off
// the WebSocket endpoint and onto a pair of inherited file descriptors: the engine
// reads commands on fd 3 and writes replies and events on fd 4, each message a JSON
// object terminated by a NUL byte. `os/exec`'s ExtraFiles supplies the descriptors
// and `encoding/json` frames the messages, so the lane gains real input with NO new
// module dependency — which a WebSocket client would have cost.
//
// WHAT IT COSTS. The engine stays alive for the length of the probe, so this
// transport cannot use `--virtual-time-budget`: virtual time races ahead of input
// that has not been sent yet. Readiness is therefore a POLLED CONDITION with a
// deadline (waitForPageCondition), never a fixed sleep.
//
// FAILING LOUDLY. Every protocol exchange is bounded by a read deadline and every
// failure is a t.Fatalf carrying the engine's stderr, because the failure mode this
// lane must never have is a probe that measured nothing and reported green. The
// engine-missing SKIP is unchanged: this transport resolves its browser through
// lookupBrowserForBehaviorProbe like every other probe, and counts itself in
// browserBehaviorProbeCount, so the strict lane's refusal to skip still covers it.
const (
	// browserProbeTrustedInputDeadline bounds one protocol exchange. A hung engine
	// must fail the probe, not the whole `go test` timeout.
	browserProbeTrustedInputDeadline = 30 * time.Second

	// browserProbePageConditionDeadline bounds waitForPageCondition. Generous
	// because the real board is a three-hundred-row render and this lane shares a
	// machine with a dozen other probes; bounded because a condition that never
	// comes true has to fail rather than hang.
	browserProbePageConditionDeadline = 45 * time.Second

	// browserProbePageConditionPollInterval is how often a pending condition is
	// re-asked. It is a poll interval, not a wait on a duration: nothing is ever
	// assumed to be ready because this much time passed.
	browserProbePageConditionPollInterval = 25 * time.Millisecond

	// browserProbeGestureSettleDeadline bounds a wait on something a gesture was
	// supposed to PRODUCE, as opposed to something the engine guarantees. The engine
	// synthesizes a click in the task after the release, so the honest budget is
	// milliseconds and this is three orders of magnitude of headroom on it. Short on
	// purpose: when the produced thing never comes, that IS the finding, and a probe
	// that sits on it for the full condition deadline reports a timeout to a reader
	// who needed the property.
	browserProbeGestureSettleDeadline = 5 * time.Second

	// browserProbePageFileName is the page all probes write into the site
	// directory, and the suffix every measurement's location.href must carry.
	browserProbePageFileName = "probe.html"
)

// trustedInputBrowserSession is one live engine with one attached page.
type trustedInputBrowserSession struct {
	probeName            string
	browserCommand       *exec.Cmd
	commandPipeWriter    *os.File
	eventPipeReader      *os.File
	eventPipeBuffer      *bufio.Reader
	standardErrorPath    string
	pageSessionId        string
	nextCommandId        int
	browserAlreadyClosed bool
}

// startTrustedInputBrowserSession renders pageHTML in a real engine and leaves it
// running with a protocol channel attached, so the caller can drive it.
//
// siteDirectory is the caller's, not a fresh temp one, for the same reason
// runBrowserBehaviorProbeInDirectory takes one: index.html loads board-data.js from
// beside itself, and a page copied somewhere empty renders an empty board.
func startTrustedInputBrowserSession(
	t *testing.T, probeName string, siteDirectory string, pageHTML string, extraFlags ...string,
) *trustedInputBrowserSession {
	t.Helper()
	browserPath := lookupBrowserForBehaviorProbe(t)

	pagePath := filepath.Join(siteDirectory, browserProbePageFileName)
	if writeError := os.WriteFile(pagePath, []byte(pageHTML), 0o644); writeError != nil {
		t.Fatalf("write %s probe page: %v", probeName, writeError)
	}
	// The profile and the engine's stderr both stay out of the site directory: one
	// is megabytes of engine state and the other is a file the board would then
	// try to render.
	scratchDirectory := t.TempDir()
	standardErrorPath := filepath.Join(scratchDirectory, "browser-stderr.log")
	standardErrorFile, standardErrorCreateError := os.Create(standardErrorPath)
	if standardErrorCreateError != nil {
		t.Fatalf("create %s browser stderr log: %v", probeName, standardErrorCreateError)
	}
	defer standardErrorFile.Close()

	// fd 3 is where the engine READS commands and fd 4 is where it WRITES replies;
	// ExtraFiles[0] lands on fd 3 in the child and ExtraFiles[1] on fd 4.
	commandPipeReader, commandPipeWriter, commandPipeError := os.Pipe()
	if commandPipeError != nil {
		t.Fatalf("open %s command pipe: %v", probeName, commandPipeError)
	}
	eventPipeReader, eventPipeWriter, eventPipeError := os.Pipe()
	if eventPipeError != nil {
		t.Fatalf("open %s event pipe: %v", probeName, eventPipeError)
	}

	// Container-safe browser flags. Keep the engine alive for protocol replies,
	// with the real clock used by both readiness and trusted input.
	probeArguments := []string{
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--user-data-dir=" + filepath.Join(scratchDirectory, "profile"),
		"--remote-debugging-pipe",
	}
	probeArguments = append(probeArguments, extraFlags...)
	probeArguments = append(probeArguments, "file://"+pagePath)
	browserCommand := exec.Command(browserPath, probeArguments...)
	browserCommand.ExtraFiles = []*os.File{commandPipeReader, eventPipeWriter}
	browserCommand.Stderr = standardErrorFile

	browserBehaviorProbeCount.Add(1)
	startError := browserCommand.Start()
	// Our copies of the child's ends are closed either way, so a failed start does
	// not leak descriptors and a live child gets EOF when we close our writer.
	commandPipeReader.Close()
	eventPipeWriter.Close()
	if startError != nil {
		commandPipeWriter.Close()
		eventPipeReader.Close()
		t.Fatalf("start %s browser with a protocol channel: %v", probeName, startError)
	}

	session := &trustedInputBrowserSession{
		probeName:         probeName,
		browserCommand:    browserCommand,
		commandPipeWriter: commandPipeWriter,
		eventPipeReader:   eventPipeReader,
		eventPipeBuffer:   bufio.NewReader(eventPipeReader),
		standardErrorPath: standardErrorPath,
	}
	t.Cleanup(session.closeBrowserSession)
	session.attachToProbePage(t)
	return session
}

// attachToProbePage waits for the engine's page target to BE the probe page, attaches
// a protocol session to it, and does not return until that page has finished loading.
//
// The target is matched on its URL rather than being taken as "the first page", because
// the engine opens its first tab on about:blank and navigates a moment later.
// about:blank reports readyState "complete" immediately, so a session attached to it
// sails through the load wait and then answers every question about the wrong document.
func (session *trustedInputBrowserSession) attachToProbePage(t *testing.T) {
	t.Helper()
	var pageTargetId string
	var lastSeenTargetURL string
	attachDeadline := time.Now().Add(browserProbePageConditionDeadline)
	for pageTargetId == "" {
		var discoveredTargets struct {
			TargetInfos []struct {
				TargetId string `json:"targetId"`
				Type     string `json:"type"`
				Url      string `json:"url"`
			} `json:"targetInfos"`
		}
		session.decodeResult(t, "Target.getTargets",
			session.callDevToolsMethod(t, "Target.getTargets", nil, false), &discoveredTargets)
		for _, targetInfo := range discoveredTargets.TargetInfos {
			if targetInfo.Type != "page" {
				continue
			}
			lastSeenTargetURL = targetInfo.Url
			if strings.HasSuffix(targetInfo.Url, "/"+browserProbePageFileName) {
				pageTargetId = targetInfo.TargetId
			}
		}
		if pageTargetId != "" {
			break
		}
		if time.Now().After(attachDeadline) {
			t.Fatalf("%s probe: the engine opened no page target on the probe page within %s "+
				"(the last page target it had was on %q)\n%s",
				session.probeName, browserProbePageConditionDeadline, lastSeenTargetURL,
				session.browserStandardError())
		}
		time.Sleep(browserProbePageConditionPollInterval)
	}

	var attachment struct {
		SessionId string `json:"sessionId"`
	}
	session.decodeResult(t, "Target.attachToTarget",
		session.callDevToolsMethod(t, "Target.attachToTarget",
			map[string]any{"targetId": pageTargetId, "flatten": true}, false), &attachment)
	if attachment.SessionId == "" {
		t.Fatalf("%s probe: attaching to the page target returned no session id\n%s",
			session.probeName, session.browserStandardError())
	}
	session.pageSessionId = attachment.SessionId

	// The target reports the probe page's URL as soon as the navigation commits in the
	// browser process, which is BEFORE the renderer has swapped documents — so the
	// first few evaluates still answer for about:blank. Every measurement after this
	// point goes through evaluateInPage, which refuses a wrong href outright, so the
	// settle is done once here on the raw channel rather than by softening that rule.
	settleDeadline := time.Now().Add(browserProbePageConditionDeadline)
	for {
		documentStateText := session.evaluateStringInPage(t,
			`JSON.stringify({href: location.href, readyState: document.readyState})`)
		var documentState struct {
			Href       string `json:"href"`
			ReadyState string `json:"readyState"`
		}
		session.decodeResult(t, "document settle", json.RawMessage(documentStateText), &documentState)
		if strings.HasSuffix(documentState.Href, "/"+browserProbePageFileName) &&
			documentState.ReadyState == "complete" {
			return
		}
		if time.Now().After(settleDeadline) {
			t.Fatalf("%s probe: the attached page was still %q (readyState %q) after %s\n%s",
				session.probeName, documentState.Href, documentState.ReadyState,
				browserProbePageConditionDeadline, session.browserStandardError())
		}
		time.Sleep(browserProbePageConditionPollInterval)
	}
}

// callDevToolsMethod sends one protocol command and returns its result.
//
// sessionScoped picks the addressee: browser-level commands (Target.*) carry no
// session id, page-level ones (Runtime.*, Input.*) carry the attached page's. Replies
// and events interleave on one pipe, so this reads until it sees its own id.
func (session *trustedInputBrowserSession) callDevToolsMethod(
	t *testing.T, method string, params map[string]any, sessionScoped bool,
) json.RawMessage {
	t.Helper()
	session.nextCommandId++
	commandId := session.nextCommandId
	command := map[string]any{"id": commandId, "method": method}
	if params != nil {
		command["params"] = params
	}
	if sessionScoped {
		command["sessionId"] = session.pageSessionId
	}
	commandBytes, marshalError := json.Marshal(command)
	if marshalError != nil {
		t.Fatalf("%s probe: encode %s: %v", session.probeName, method, marshalError)
	}
	if _, writeError := session.commandPipeWriter.Write(append(commandBytes, 0)); writeError != nil {
		t.Fatalf("%s probe: send %s over the protocol pipe: %v\n%s",
			session.probeName, method, writeError, session.browserStandardError())
	}

	replyDeadline := time.Now().Add(browserProbeTrustedInputDeadline)
	for {
		if deadlineError := session.eventPipeReader.SetReadDeadline(replyDeadline); deadlineError != nil {
			t.Fatalf("%s probe: set a read deadline for %s: %v",
				session.probeName, method, deadlineError)
		}
		messageBytes, readError := session.eventPipeBuffer.ReadBytes(0)
		if readError != nil {
			t.Fatalf("%s probe: no reply to %s within %s (%v) — the protocol channel is the "+
				"transport, so this is a transport failure and not a measurement\n%s",
				session.probeName, method, browserProbeTrustedInputDeadline, readError,
				session.browserStandardError())
		}
		var reply struct {
			Id     int             `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  json.RawMessage `json:"error"`
		}
		if decodeError := json.Unmarshal(bytes.TrimRight(messageBytes, "\x00"), &reply); decodeError != nil {
			t.Fatalf("%s probe: the engine sent something that is not a protocol message after %s: %v\n%s",
				session.probeName, method, decodeError, messageBytes)
		}
		// Everything else on this pipe is a protocol EVENT or another command's
		// reply; only our own id answers this call.
		if reply.Id != commandId {
			continue
		}
		if reply.Error != nil {
			t.Fatalf("%s probe: %s failed: %s\n%s",
				session.probeName, method, reply.Error, session.browserStandardError())
		}
		return reply.Result
	}
}

func (session *trustedInputBrowserSession) decodeResult(
	t *testing.T, method string, resultJSON json.RawMessage, destination any,
) {
	t.Helper()
	if decodeError := json.Unmarshal(resultJSON, destination); decodeError != nil {
		t.Fatalf("%s probe: decode the %s result %s: %v",
			session.probeName, method, resultJSON, decodeError)
	}
}

// evaluateInPage runs one JavaScript expression in the probe page and returns its
// value as JSON for the caller to unmarshal.
//
// The expression is wrapped so every result carries the page it was measured on.
// That is the prime's render-evidence rule made mechanical: a probe on this
// transport cannot forget to report location.href, because the transport reports it
// and checks it on every single call rather than once at the start.
//
// The wrapper resolves the expression as a PROMISE before reading location.href, which
// buys two things. An expression may await the page's own scheduling — settling a
// pending render before a gesture is dispatched needs exactly that — and the href is
// then read at the instant the value was produced rather than at the instant the
// expression started, which is the stronger version of the same evidence rule.
func (session *trustedInputBrowserSession) evaluateInPage(
	t *testing.T, expression string,
) json.RawMessage {
	t.Helper()
	envelopeText := session.evaluateStringInPage(t,
		"Promise.resolve(("+expression+")).then(function (probeValue) {"+
			" return JSON.stringify({href: location.href, value: probeValue}); })")
	var envelope struct {
		Href  string          `json:"href"`
		Value json.RawMessage `json:"value"`
	}
	if decodeError := json.Unmarshal([]byte(envelopeText), &envelope); decodeError != nil {
		t.Fatalf("%s probe: decode the envelope %q from %s: %v",
			session.probeName, envelopeText, expression, decodeError)
	}
	// Render evidence, on every measurement rather than once: a page that navigated
	// out from under the probe answers confidently about somebody else's document.
	if !strings.HasSuffix(envelope.Href, "/"+browserProbePageFileName) {
		t.Fatalf("%s probe: measured on %q, not the probe page — every number from this call "+
			"describes a document this test did not render", session.probeName, envelope.Href)
	}
	if envelope.Value == nil {
		t.Fatalf("%s probe: evaluating %s produced undefined; a probe that measures nothing must "+
			"not read as a probe that measured", session.probeName, expression)
	}
	return envelope.Value
}

// evaluateStringInPage is the raw half of evaluateInPage: it runs an expression that
// must produce a string and returns it, with no envelope and no render-evidence check.
//
// Only two callers may use it — evaluateInPage itself, and the attach settle, which
// cannot check the href because deciding whether the href is right yet is its whole
// job. Everything a probe measures goes through evaluateInPage.
func (session *trustedInputBrowserSession) evaluateStringInPage(
	t *testing.T, expression string,
) string {
	t.Helper()
	var evaluation struct {
		Result struct {
			Value json.RawMessage `json:"value"`
		} `json:"result"`
		ExceptionDetails json.RawMessage `json:"exceptionDetails"`
	}
	session.decodeResult(t, "Runtime.evaluate",
		session.callDevToolsMethod(t, "Runtime.evaluate", map[string]any{
			"expression":    expression,
			"returnByValue": true,
			"awaitPromise":  true,
		}, true), &evaluation)
	if evaluation.ExceptionDetails != nil {
		t.Fatalf("%s probe: the page threw evaluating %s: %s",
			session.probeName, expression, evaluation.ExceptionDetails)
	}
	var stringValue string
	if decodeError := json.Unmarshal(evaluation.Result.Value, &stringValue); decodeError != nil {
		t.Fatalf("%s probe: evaluating %s returned %s, not the string the transport expects",
			session.probeName, expression, evaluation.Result.Value)
	}
	return stringValue
}

// waitForPageCondition re-asks a boolean expression until it holds. The wait is on
// the CONDITION and the deadline only bounds it, so a probe never proceeds because
// enough time passed — it proceeds because the page says it is ready.
//
// Use this for a condition the engine GUARANTEES will arrive: a load completing, a
// dispatched event being delivered. For a condition that is itself the thing under
// test, use pageConditionHoldsWithin and say what its absence means — a timeout
// message names the wait, and the reader needs the property.
func (session *trustedInputBrowserSession) waitForPageCondition(
	t *testing.T, conditionDescription string, booleanExpression string,
) {
	t.Helper()
	if !session.pageConditionHoldsWithin(
		t, conditionDescription, booleanExpression, browserProbePageConditionDeadline) {
		t.Fatalf("%s probe: waited %s for %s and it never became true (expression %s)",
			session.probeName, browserProbePageConditionDeadline, conditionDescription,
			booleanExpression)
	}
}

// pageConditionHoldsWithin is waitForPageCondition without the verdict: it reports
// whether the condition came true inside its own budget and leaves what that means to
// the caller.
func (session *trustedInputBrowserSession) pageConditionHoldsWithin(
	t *testing.T, conditionDescription string, booleanExpression string, conditionBudget time.Duration,
) bool {
	t.Helper()
	conditionDeadline := time.Now().Add(conditionBudget)
	for {
		var conditionHolds bool
		session.decodeResult(t, "condition "+conditionDescription,
			session.evaluateInPage(t, "!!("+booleanExpression+")"), &conditionHolds)
		if conditionHolds {
			return true
		}
		if time.Now().After(conditionDeadline) {
			return false
		}
		time.Sleep(browserProbePageConditionPollInterval)
	}
}

// The three gestures. What arrives in the page is a TRUSTED event: isTrusted is
// true, the ENGINE issues the pointerId, and setPointerCapture on that id
// establishes a real capture that really retargets the synthesized click. That last
// sentence is the entire reason this transport exists.
//
// Coordinates are viewport CSS pixels, the same frame getBoundingClientRect reports
// in — which is why a caller must scroll its target on screen before aiming at it.
func (session *trustedInputBrowserSession) pressTrustedMouseAt(t *testing.T, viewportX, viewportY float64) {
	t.Helper()
	session.dispatchTrustedMouseEvent(t, "mousePressed", viewportX, viewportY, "left", 1)
}

func (session *trustedInputBrowserSession) dragTrustedMouseTo(t *testing.T, viewportX, viewportY float64) {
	t.Helper()
	session.dispatchTrustedMouseEvent(t, "mouseMoved", viewportX, viewportY, "left", 1)
}

func (session *trustedInputBrowserSession) releaseTrustedMouseAt(t *testing.T, viewportX, viewportY float64) {
	t.Helper()
	session.dispatchTrustedMouseEvent(t, "mouseReleased", viewportX, viewportY, "left", 0)
}

func (session *trustedInputBrowserSession) dispatchTrustedMouseEvent(
	t *testing.T, eventType string, viewportX, viewportY float64, mouseButton string, heldButtons int,
) {
	t.Helper()
	session.callDevToolsMethod(t, "Input.dispatchMouseEvent", map[string]any{
		"type":       eventType,
		"x":          viewportX,
		"y":          viewportY,
		"button":     mouseButton,
		"buttons":    heldButtons,
		"clickCount": 1,
	}, true)
}

// browserStandardError returns what the engine wrote to stderr, tail-first for the
// failure messages above. It is read from a file rather than an in-process buffer so
// that reading it while the child is still writing is not a data race.
func (session *trustedInputBrowserSession) browserStandardError() string {
	standardErrorBytes, readError := os.ReadFile(session.standardErrorPath)
	if readError != nil {
		return "browser stderr unavailable: " + readError.Error()
	}
	const standardErrorTailBytes = 4000
	if len(standardErrorBytes) > standardErrorTailBytes {
		standardErrorBytes = standardErrorBytes[len(standardErrorBytes)-standardErrorTailBytes:]
	}
	return "browser stderr:\n" + string(standardErrorBytes)
}

// closeBrowserSession shuts the engine down. Closing the command pipe is the
// protocol's own goodbye — the engine exits when its reader hits EOF — and the kill
// is there so a wedged engine cannot hold the test binary open.
func (session *trustedInputBrowserSession) closeBrowserSession() {
	if session.browserAlreadyClosed {
		return
	}
	session.browserAlreadyClosed = true
	session.commandPipeWriter.Close()
	browserExited := make(chan error, 1)
	go func() { browserExited <- session.browserCommand.Wait() }()
	select {
	case <-browserExited:
	case <-time.After(browserProbeTrustedInputDeadline):
		_ = session.browserCommand.Process.Kill()
		<-browserExited
	}
	session.eventPipeReader.Close()
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

type ticketMentionLinkProbe struct {
	DetailId      string `json:"detailId"`
	DetailKind    string `json:"detailKind"`
	TooltipTitle  string `json:"tooltipTitle"`
	ExpandedTitle string `json:"expandedTitle"`
	InsideCode    bool   `json:"insideCode"`
	Text          string `json:"text"`
}

type ticketMissingMentionProbe struct {
	TagName      string `json:"tagName"`
	Text         string `json:"text"`
	TooltipTitle string `json:"tooltipTitle"`
	IsAnchor     bool   `json:"isAnchor"`
}

type ticketGlossaryBrowserRow struct {
	Identifier string `json:"identifier"`
	DetailKind string `json:"detailKind"`
	Title      string `json:"title"`
	Status     string `json:"status"`
}

type ticketDrawerProbeSnapshot struct {
	DetailKind           string                      `json:"detailKind"`
	DetailId             string                      `json:"detailId"`
	BodyText             string                      `json:"bodyText"`
	BodyLinks            []ticketMentionLinkProbe    `json:"bodyLinks"`
	MissingMentions      []ticketMissingMentionProbe `json:"missingMentions"`
	MetaLinks            []ticketMentionLinkProbe    `json:"metaLinks"`
	GlossaryHidden       bool                        `json:"glossaryHidden"`
	GlossaryRows         []ticketGlossaryBrowserRow  `json:"glossaryRows"`
	HorizontalOverflow   float64                     `json:"horizontalOverflow"`
	TitleSpanFontFamily  string                      `json:"titleSpanFontFamily"`
	IdSpanFontFamily     string                      `json:"idSpanFontFamily"`
	TitleSpanColour      string                      `json:"titleSpanColour"`
	IdSpanColour         string                      `json:"idSpanColour"`
	MissingMentionColour string                      `json:"missingMentionColour"`
}

type ticketDrawerProbeResult struct {
	LocationHref string `json:"locationHref"`
	// The scheme the ENGINE resolved, not the one the flag asked for: a flag this
	// build silently ignores would let one palette be measured twice.
	ResolvedScheme    string                    `json:"resolvedScheme"`
	SurfaceColour     string                    `json:"surfaceColour"`
	RequestDrawer     ticketDrawerProbeSnapshot `json:"requestDrawer"`
	UserRequestDrawer ticketDrawerProbeSnapshot `json:"userRequestDrawer"`
	NoReferenceDrawer ticketDrawerProbeSnapshot `json:"noReferenceDrawer"`
	ConsoleErrors     []string                  `json:"consoleErrors"`
}

func ticketMentionFixtureTicket(requestID string, title string, status string, bodyMarkdown string) *RequestTicket {
	return &RequestTicket{
		RequestId:           requestID,
		Title:               title,
		Status:              status,
		OriginalStatus:      status,
		Domain:              "frontend",
		OriginalDomain:      "frontend",
		TreeSection:         "queue",
		CreatedAt:           "2026-08-20T12:00:00Z",
		FrontmatterMarkdown: "---\nid: " + requestID + "\ntitle: '" + title + "'\nstatus: " + status + "\n---\n",
		BodyMarkdown:        bodyMarkdown,
	}
}

// REQ-374 renders a body that cites one id twice, one inside a code span, a user
// request, an id no record answers to, and an AMBIGUOUS segment two compound ids
// share. Driving the real generated board is what makes the claim about pixels
// checkable: the Node lane can prove the fragment shape, only an engine can say
// the expansion did not push the drawer into horizontal overflow, and only a
// real drawer reopen can prove the glossary does not survive into the next
// ticket.
func TestBrowserBehaviorDrawerTicketTitlesAndGlossary(t *testing.T) {
	citedRequestTitle := "Keep the timeline forecast honest about ordering and timings"
	citingRequest := ticketMentionFixtureTicket("REQ-500", "Cites other tickets", "claimed",
		"## Notes\n\nRead REQ-1679 lessons, `REQ-1108` in code, UR-074 for context, "+
			"plus REQ-9999 and REQ-042 stay honest.\n\nLater the REQ-1679 note matters again.\n")
	citingRequest.UserRequestId = "UR-074"
	citingRequest.DependsOn = []string{"REQ-1679"}
	// BlockedBy and WriteSetOverlaps reach makeTicketLinkList, which serves three
	// of the five meta rows the acceptance criterion names. Without them that
	// function is never called in this fixture, and flipping its expandTitle
	// argument to false passed the whole suite.
	citingRequest.BlockedBy = []string{"REQ-1108"}
	citingRequest.WriteSetOverlaps = []string{"REQ-600"}
	citedRequest := ticketMentionFixtureTicket("REQ-1679", citedRequestTitle, "completed", "Cited body.\n")
	citedRequest.TreeSection = "archive"
	backtickedRequest := ticketMentionFixtureTicket("REQ-1108", "Short one", "pending", "Backticked body.\n")
	firstAmbiguous := ticketMentionFixtureTicket("UR-001-REQ-042", "First half of an ambiguous pair", "pending", "A.\n")
	secondAmbiguous := ticketMentionFixtureTicket("UR-002-REQ-042", "Second half of an ambiguous pair", "pending", "B.\n")
	noReferenceRequest := ticketMentionFixtureTicket("REQ-600", "Cites nothing", "pending",
		"Nothing to cross-reference here.\n")

	citingUserRequest := &UserRequestTicket{
		UserRequestId:       "UR-074",
		Title:               "Ticket ids should carry their titles",
		FrontmatterMarkdown: "---\nid: UR-074\ntitle: 'Ticket ids should carry their titles'\nrequests: [REQ-500]\n---\n",
		BodyMarkdown:        "The ask, in full. See REQ-1108 for the shape.\n",
		InputFilePresent:    true,
	}

	fixtureBoard := &Board{
		GeneratedAt: time.Now().UTC(),
		ProjectName: "REQ-374 ticket mention probe",
		AllRequests: []*RequestTicket{
			citingRequest, citedRequest, backtickedRequest, firstAmbiguous, secondAmbiguous, noReferenceRequest,
		},
		UserRequests:     []*UserRequestTicket{citingUserRequest},
		UserRequestsById: map[string]*UserRequestTicket{citingUserRequest.UserRequestId: citingUserRequest},
		Columns: BoardColumns{
			Pending:      []*RequestTicket{backtickedRequest, firstAmbiguous, secondAmbiguous, noReferenceRequest},
			PendingReady: []*RequestTicket{backtickedRequest, firstAmbiguous, secondAmbiguous, noReferenceRequest},
			Claimed:      []*RequestTicket{citingRequest},
		},
	}
	linkRequestsToUserRequests(fixtureBoard)

	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate ticket mention fixture: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read ticket mention fixture: %v", readError)
	}

	errorCaptureStub := `
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

        function describeTicketLink(linkNode) {
          var titleSpan = linkNode.querySelector(".ticket-link-title");
          return {
            detailId: linkNode.dataset.detailId,
            detailKind: linkNode.dataset.detailKind,
            tooltipTitle: linkNode.getAttribute("title") || "",
            expandedTitle: titleSpan ? titleSpan.textContent : "",
            insideCode: Boolean(linkNode.closest("code")),
            text: linkNode.textContent
          };
        }

        function describeDrawer() {
          var drawerNode = document.getElementById("detail-drawer");
          var bodyNode = document.getElementById("detail-body");
          var glossaryNode = document.getElementById("detail-glossary");
          var firstTitleSpan = bodyNode.querySelector(".ticket-link-title");
          var firstIdSpan = bodyNode.querySelector(".ticket-link-id");
          var firstMissingMention = bodyNode.querySelector(".ticket-missing");
          return {
            detailKind: document.getElementById("detail-kind").textContent,
            detailId: document.getElementById("detail-id").textContent,
            bodyText: bodyNode.textContent,
            bodyLinks: Array.from(bodyNode.querySelectorAll("a.ticket-link")).map(describeTicketLink),
            missingMentions: Array.from(bodyNode.querySelectorAll(".ticket-missing")).map(function (node) {
              return {
                tagName: node.tagName,
                text: node.textContent,
                tooltipTitle: node.getAttribute("title") || "",
                isAnchor: node.tagName === "A"
              };
            }),
            metaLinks: Array.from(document.getElementById("detail-meta").querySelectorAll("a.ticket-link"))
              .map(describeTicketLink),
            glossaryHidden: glossaryNode.hidden,
            glossaryRows: Array.from(glossaryNode.querySelectorAll(".detail-glossary-term")).map(function (termNode) {
              var definitionNode = termNode.nextElementSibling;
              return {
                identifier: termNode.textContent.trim(),
                detailKind: termNode.querySelector("a.ticket-link").dataset.detailKind,
                title: definitionNode.querySelector(".detail-glossary-name").textContent,
                status: definitionNode.querySelector(".detail-glossary-status").textContent
              };
            }),
            horizontalOverflow: drawerNode.scrollWidth - drawerNode.clientWidth,
            titleSpanFontFamily: firstTitleSpan ? getComputedStyle(firstTitleSpan).fontFamily : "",
            idSpanFontFamily: firstIdSpan ? getComputedStyle(firstIdSpan).fontFamily : "",
            titleSpanColour: firstTitleSpan ? getComputedStyle(firstTitleSpan).color : "",
            idSpanColour: firstIdSpan ? getComputedStyle(firstIdSpan).color : "",
            missingMentionColour: firstMissingMention ? getComputedStyle(firstMissingMention).color : ""
          };
        }

        function openCard(requestId) {
          var card = document.querySelector('.req-card[data-detail-id="' + requestId + '"]');
          if (!card) {
            throw new Error("no card for " + requestId);
          }
          card.click();
          return waitFor(function () {
            return document.getElementById("detail-id").textContent === requestId;
          }, requestId + " drawer");
        }

        async function runProbe() {
          var result = {
            locationHref: location.href,
            resolvedScheme: matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
            surfaceColour: getComputedStyle(document.body).backgroundColor
          };

          await openCard("REQ-500");
          result.requestDrawer = describeDrawer();

          document.querySelector('#detail-body a[data-detail-kind="ur"][data-detail-id="UR-074"]').click();
          await waitFor(function () {
            return document.getElementById("detail-id").textContent === "UR-074";
          }, "UR-074 drawer");
          result.userRequestDrawer = describeDrawer();

          await openCard("REQ-600");
          result.noReferenceDrawer = describeDrawer();

          result.consoleErrors = window.__queueKanbanProbeErrors.slice();
          resultNode.textContent = JSON.stringify(result);
        }

        runProbe().catch(function (probeError) {
          resultNode.textContent = JSON.stringify({
            locationHref: location.href,
            resolvedScheme: matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
            consoleErrors: window.__queueKanbanProbeErrors.concat([String((probeError && probeError.message) || probeError)])
          });
        });
      })();
`

	probePage := string(indexBytes)
	clientScriptOpen := "    <script>\n"
	clientScriptOpenIndex := strings.LastIndex(probePage, clientScriptOpen)
	if clientScriptOpenIndex < 0 {
		t.Fatal("generated page has no client opening for the ticket mention stub")
	}
	errorStubIndex := clientScriptOpenIndex + len(clientScriptOpen)
	probePage = probePage[:errorStubIndex] + errorCaptureStub + probePage[errorStubIndex:]
	clientCloseIndex := strings.LastIndex(probePage, "})();")
	if clientCloseIndex < 0 {
		t.Fatal("generated page has no client close for the ticket mention probe")
	}
	clientCloseIndex += len("})();")
	probePage = probePage[:clientCloseIndex] + "\n" + probeScript + probePage[clientCloseIndex:]

	// BOTH palettes. This board is dark-first — :root is the dark palette and
	// @media (prefers-color-scheme: light) overrides it — and Chromium resolves
	// light with no flag, so a single run leaves the other palette checked by
	// nothing. The expansion's ink hierarchy is a claim about both.
	for _, colourScheme := range []struct {
		name string
		flag string
	}{
		{name: "light", flag: "--blink-settings=preferredColorScheme=1"},
		{name: "dark", flag: "--blink-settings=preferredColorScheme=0"},
	} {
		t.Run(colourScheme.name, func(t *testing.T) {
			assertDrawerTicketTitlesAndGlossary(
				t, siteDirectory, probePage, colourScheme.name, colourScheme.flag, citedRequestTitle)
		})
	}
}

func assertDrawerTicketTitlesAndGlossary(
	t *testing.T, siteDirectory string, probePage string,
	schemeName string, schemeFlag string, citedRequestTitle string,
) {
	t.Helper()
	resultJSON := runBrowserBehaviorProbeInDirectory(
		t, "drawer ticket titles ("+schemeName+")", siteDirectory, probePage,
		"--virtual-time-budget=30000", schemeFlag,
	)
	var result ticketDrawerProbeResult
	if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
		t.Fatalf("decode ticket mention probe: %v\n%s", decodeError, resultJSON)
	}
	if !strings.HasSuffix(result.LocationHref, "/"+browserProbePageFileName) {
		t.Fatalf("ticket mention probe measured %q, not its probe page", result.LocationHref)
	}
	if result.ResolvedScheme != schemeName {
		t.Fatalf("asked the engine for the %s palette and it resolved %s; this build ignores the "+
			"colour-scheme flag, so one palette would be measured twice and reported as two",
			schemeName, result.ResolvedScheme)
	}
	if len(result.ConsoleErrors) != 0 {
		t.Fatalf("ticket mention browser errors: %q", result.ConsoleErrors)
	}

	requestDrawer := result.RequestDrawer
	wantBodyLinks := []ticketMentionLinkProbe{
		{
			DetailId: "REQ-1679", DetailKind: "req", TooltipTitle: citedRequestTitle,
			ExpandedTitle: citedRequestTitle, InsideCode: false,
			Text: "REQ-1679 " + citedRequestTitle,
		},
		{DetailId: "REQ-1108", DetailKind: "req", InsideCode: true, Text: "REQ-1108"},
		{
			DetailId: "UR-074", DetailKind: "ur", TooltipTitle: "Ticket ids should carry their titles",
			ExpandedTitle: "Ticket ids should carry their titles",
			Text:          "UR-074 Ticket ids should carry their titles",
		},
		{DetailId: "REQ-1679", DetailKind: "req", Text: "REQ-1679"},
	}
	if !reflect.DeepEqual(requestDrawer.BodyLinks, wantBodyLinks) {
		t.Errorf("REQ-500 body links = %#v,\nwant %#v", requestDrawer.BodyLinks, wantBodyLinks)
	}
	wantMissingMentions := []ticketMissingMentionProbe{
		{TagName: "SPAN", Text: "REQ-9999", TooltipTitle: "Not found in this queue", IsAnchor: false},
	}
	if !reflect.DeepEqual(requestDrawer.MissingMentions, wantMissingMentions) {
		t.Errorf("REQ-500 unresolved mentions = %#v, want %#v", requestDrawer.MissingMentions, wantMissingMentions)
	}
	if !strings.Contains(requestDrawer.BodyText, "REQ-042 stay honest") {
		t.Errorf("REQ-500 body = %q, want the ambiguous segment left as plain prose", requestDrawer.BodyText)
	}
	for _, metaLink := range requestDrawer.MetaLinks {
		if metaLink.ExpandedTitle == "" || metaLink.TooltipTitle == "" {
			t.Errorf("meta row link %s carries no title: %#v", metaLink.DetailId, metaLink)
		}
	}
	// Every meta row the criterion names, by the function that renders it:
	// makeDependencyDetailList (Depends on), makeTicketLinkList (Blocked by,
	// Overlapping write sets) and the direct User request call. Naming the ids
	// rather than counting is what makes this bite — a count passes while a whole
	// row silently loses its titles.
	metaLinkTitles := map[string]string{}
	for _, metaLink := range requestDrawer.MetaLinks {
		metaLinkTitles[metaLink.DetailId] = metaLink.ExpandedTitle
	}
	for _, wantMetaRow := range []struct{ detailId, renderedBy string }{
		{"REQ-1679", "makeDependencyDetailList (Depends on)"},
		{"REQ-1108", "makeTicketLinkList (Blocked by)"},
		{"REQ-600", "makeTicketLinkList (Overlapping write sets)"},
		{"UR-074", "the direct User request call"},
	} {
		if metaLinkTitles[wantMetaRow.detailId] == "" {
			t.Errorf("meta row link %s carries no expanded title — %s stopped expanding",
				wantMetaRow.detailId, wantMetaRow.renderedBy)
		}
	}
	wantRequestGlossary := []ticketGlossaryBrowserRow{
		{Identifier: "REQ-1679", DetailKind: "req", Title: citedRequestTitle, Status: "completed"},
		{Identifier: "REQ-1108", DetailKind: "req", Title: "Short one", Status: "pending"},
		{Identifier: "UR-074", DetailKind: "ur", Title: "Ticket ids should carry their titles", Status: "user request"},
	}
	if requestDrawer.GlossaryHidden {
		t.Error("REQ-500 cited five ids and its glossary stayed hidden")
	}
	if !reflect.DeepEqual(requestDrawer.GlossaryRows, wantRequestGlossary) {
		t.Errorf("REQ-500 glossary = %#v,\nwant %#v", requestDrawer.GlossaryRows, wantRequestGlossary)
	}
	// The expansion's whole reason for truncating is that an over-long title must
	// not widen the drawer at its default width.
	if requestDrawer.HorizontalOverflow > 1 {
		t.Errorf("expanded ticket titles pushed the drawer %vpx into horizontal overflow", requestDrawer.HorizontalOverflow)
	}
	if requestDrawer.TitleSpanFontFamily == requestDrawer.IdSpanFontFamily {
		t.Errorf("title span kept the id's font (%q) — .ticket-link's mono was not reset",
			requestDrawer.TitleSpanFontFamily)
	}
	if requestDrawer.TitleSpanColour == requestDrawer.IdSpanColour {
		t.Errorf("title span and id span share colour %q — the title is not one step dimmer",
			requestDrawer.TitleSpanColour)
	}
	// The dimmer step is only worth having if the title is still readable, and a
	// broken reference only reads as broken if it is not the body's own ink.
	surfaceLuminance, surfaceKnown := relativeLuminanceOfCSSColour(result.SurfaceColour)
	titleLuminance, titleKnown := relativeLuminanceOfCSSColour(requestDrawer.TitleSpanColour)
	if !surfaceKnown || !titleKnown {
		t.Fatalf("could not read the %s palette colours: surface %q, title %q",
			schemeName, result.SurfaceColour, requestDrawer.TitleSpanColour)
	}
	if titleContrast := contrastRatio(titleLuminance, surfaceLuminance); titleContrast < 4.5 {
		t.Errorf("expanded title contrast against the page surface = %.2f:1 in the %s palette, want 4.5:1 for body text",
			titleContrast, schemeName)
	}
	if requestDrawer.MissingMentionColour == "" {
		t.Error("the unresolved mention rendered no colour to measure")
	} else if requestDrawer.MissingMentionColour == requestDrawer.TitleSpanColour {
		t.Errorf("a broken reference is drawn in the same ink as an ordinary title (%q) in the %s palette",
			requestDrawer.MissingMentionColour, schemeName)
	}

	// The UR drawer linkifies its own body, and its glossary replaces — never
	// inherits — the one the REQ drawer left behind.
	userRequestDrawer := result.UserRequestDrawer
	if userRequestDrawer.DetailKind != "UR" {
		t.Fatalf("second drawer kind = %q, want UR", userRequestDrawer.DetailKind)
	}
	wantUserRequestGlossary := []ticketGlossaryBrowserRow{
		{Identifier: "REQ-1108", DetailKind: "req", Title: "Short one", Status: "pending"},
	}
	if !reflect.DeepEqual(userRequestDrawer.GlossaryRows, wantUserRequestGlossary) {
		t.Errorf("UR-074 glossary = %#v, want only its own citation %#v",
			userRequestDrawer.GlossaryRows, wantUserRequestGlossary)
	}
	if len(userRequestDrawer.BodyLinks) != 1 || userRequestDrawer.BodyLinks[0].ExpandedTitle != "Short one" {
		t.Errorf("UR-074 body links = %#v, want one expanded REQ-1108 link", userRequestDrawer.BodyLinks)
	}

	// A body that cited nothing gets no glossary at all — including right after a
	// body that did.
	noReferenceDrawer := result.NoReferenceDrawer
	if !noReferenceDrawer.GlossaryHidden || len(noReferenceDrawer.GlossaryRows) != 0 {
		t.Errorf("REQ-600 cited nothing but kept a glossary: hidden=%v rows=%#v",
			noReferenceDrawer.GlossaryHidden, noReferenceDrawer.GlossaryRows)
	}
}

// A result is JSON text, not serialized HTML. Clipboard arrows and entity-like
// title text must cross the browser boundary without an extra encode/decode pass.
func TestBrowserBehaviorProbeResultKeepsLiteralText(t *testing.T) {
	const literalText = "<em>&amp;</em> -> & \"quoted\""
	page := `<!doctype html><pre id="` + browserProbeResultElementId + `"></pre><script>
 document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({text: ` + mustMarshalJSONString(t, literalText) + `});
</script>`
	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(runBrowserBehaviorProbe(t, "literal result text", page), &result); err != nil {
		t.Fatal(err)
	}
	if result.Text != literalText {
		t.Fatalf("result text changed across browser transport: got %q, want %q", result.Text, literalText)
	}
}
