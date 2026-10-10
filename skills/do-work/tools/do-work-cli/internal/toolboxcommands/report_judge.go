package toolboxcommands

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// ai-report-judge serves one report bundle on a loopback port, renders its
// index.html in a headless Chromium-family engine at a wide and a phone width in
// light and dark, and fails on horizontal overflow or a same-origin href/src that
// answers 400 or more. Captures and judge.json go outside the bundle.

const (
	judgeCodeOverflow   = "AI-REPORT-JUDGE-OVERFLOW"
	judgeCodeBrokenLink = "AI-REPORT-JUDGE-BROKEN-LINK"
	judgeCodeSkipped    = "AI-REPORT-JUDGE-SKIPPED"
	judgeCodeError      = "AI-REPORT-JUDGE-ERROR"

	// judgeProtocolDeadline bounds every wait on the engine: one protocol reply,
	// one page settle, the engine's exit. A hung engine ends the run as an error.
	judgeProtocolDeadline = 30 * time.Second
	judgePollInterval     = 25 * time.Millisecond
	judgeUsage            = "Usage: ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]"
)

// aiReportJudgeBrowserNames and aiReportJudgeBrowserAppPaths are checked in that
// order when --browser is absent. They are a convenience, never a closed set:
// --browser makes any other Chromium-family engine usable.
var aiReportJudgeBrowserNames = []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "chrome"}
var aiReportJudgeBrowserAppPaths = []string{
	"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	"/Applications/Chromium.app/Contents/MacOS/Chromium",
}

var judgeViewports = []struct {
	name          string
	width, height int
}{{"wide", 1440, 1000}, {"phone", 390, 844}}

type judgeCapture struct {
	Viewport    string `json:"viewport"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Scheme      string `json:"scheme"`
	File        string `json:"file"`
	ScrollWidth int    `json:"scroll_width"`
	InnerWidth  int    `json:"inner_width"`
}

type judgeFinding struct {
	Kind        string `json:"kind"`
	Viewport    string `json:"viewport,omitempty"`
	Scheme      string `json:"scheme,omitempty"`
	ScrollWidth int    `json:"scroll_width,omitempty"`
	InnerWidth  int    `json:"inner_width,omitempty"`
	Attribute   string `json:"attribute,omitempty"`
	Target      string `json:"target,omitempty"`
	Status      int    `json:"status,omitempty"`
	Detail      string `json:"detail,omitempty"`
}

type judgeReport struct {
	SchemaVersion int            `json:"schema_version"`
	Bundle        string         `json:"bundle"`
	ServedURL     string         `json:"served_url"`
	Verdict       string         `json:"verdict"`
	Browser       string         `json:"browser"`
	Captures      []judgeCapture `json:"captures"`
	Findings      []judgeFinding `json:"findings"`
}

func handleAIReportJudge(ctx commandruntime.ExecutionContext, args []string) resultmodel.CommandResult {
	var bundle, output, browser string
	for index := 0; index < len(args); index++ {
		argument := args[index]
		var err error
		switch {
		case argument == "--out" || strings.HasPrefix(argument, "--out="):
			output, err = optionValue(args, &index, "--out")
		case argument == "--browser" || strings.HasPrefix(argument, "--browser="):
			browser, err = optionValue(args, &index, "--browser")
		case strings.HasPrefix(argument, "-") || bundle != "":
			err = fmt.Errorf("unexpected argument %q", argument)
		default:
			bundle = argument
		}
		if err != nil {
			return usageResult(CommandAIReportJudge, err.Error()+"; "+judgeUsage)
		}
	}
	if bundle == "" {
		return usageResult(CommandAIReportJudge, judgeUsage)
	}
	bundle = judgeAbsolutePath(ctx.RepositoryRoot, bundle)
	if info, err := os.Stat(filepath.Join(bundle, "index.html")); err != nil || !info.Mode().IsRegular() {
		return usageResult(CommandAIReportJudge, "the bundle must be a directory holding index.html: "+bundle)
	}
	resolvedBundle, err := filepath.EvalSymlinks(bundle)
	if err != nil {
		return usageResult(CommandAIReportJudge, err.Error())
	}
	if browser != "" {
		resolvedBrowser := judgeLookupBrowser(browser)
		if resolvedBrowser == "" {
			return usageResult(CommandAIReportJudge, "--browser names no runnable browser: "+browser)
		}
		browser = resolvedBrowser
	} else {
		browser = discoverJudgeBrowser()
	}
	if output == "" {
		if output, err = os.MkdirTemp("", "ai-report-judge-"); err != nil {
			return judgeFailureResult(judgeCodeError, "create the output directory: "+err.Error())
		}
	}
	output = judgeAbsolutePath(ctx.RepositoryRoot, output)
	if judgePathInside(judgeResolveExistingPrefix(output), resolvedBundle) {
		return usageResult(CommandAIReportJudge, "--out must be outside the bundle: "+output)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return usageResult(CommandAIReportJudge, "create the output directory: "+err.Error())
	}
	report := judgeReport{SchemaVersion: 1, Bundle: bundle, Captures: []judgeCapture{}, Findings: []judgeFinding{}}
	if browser == "" {
		report.Verdict = "skipped"
		report.Findings = append(report.Findings, judgeFinding{Kind: judgeCodeSkipped,
			Detail: "no browser engine was found; pass --browser <path>"})
		return judgeResult(report, output)
	}
	report.Browser = browser
	if err := judgeRender(&report, bundle, output, browser); err != nil {
		return judgeErrorResult(report, output, err.Error())
	}
	report.Verdict = "pass"
	if len(report.Findings) > 0 {
		report.Verdict = "fail"
	}
	return judgeResult(report, output)
}

// judgeRender owns the server and the engine: both stop before it returns, on
// every path, because each is closed by a defer set right after it starts.
func judgeRender(report *judgeReport, bundle, output, browser string) error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("open the bundle server: %w", err)
	}
	server := &http.Server{Handler: http.FileServer(http.Dir(bundle))}
	defer server.Close()
	go func() { _ = server.Serve(listener) }()
	origin := "http://" + listener.Addr().String()
	report.ServedURL = origin + "/index.html"

	session, err := startJudgeBrowser(browser)
	if err != nil {
		return err
	}
	defer session.stop()
	if err := session.attachToPage(); err != nil {
		return err
	}

	type pageLink struct{ Attribute, Value, Resolved string }
	var links []pageLink
	for _, viewport := range judgeViewports {
		for _, scheme := range []string{"light", "dark"} {
			capture, err := session.renderPass(report.ServedURL, viewport.name, viewport.width, viewport.height, scheme, output)
			if err != nil {
				return err
			}
			report.Captures = append(report.Captures, capture)
			if capture.ScrollWidth > capture.InnerWidth {
				report.Findings = append(report.Findings, judgeFinding{Kind: judgeCodeOverflow, Viewport: viewport.name,
					Scheme: scheme, ScrollWidth: capture.ScrollWidth, InnerWidth: capture.InnerWidth})
			}
			if links == nil {
				// Collected once, in the first (wide light) pass.
				linksJSON, err := session.evaluate(`Array.from(document.querySelectorAll("[href],[src]")).flatMap(function (element) {
  return ["href", "src"].filter(function (name) { return element.hasAttribute(name); }).map(function (name) {
    var value = element.getAttribute(name), resolved = "";
    try { resolved = new URL(value, document.baseURI).href; } catch (error) {}
    return {Attribute: name, Value: value, Resolved: resolved};
  });
})`)
				if err != nil {
					return err
				}
				links = []pageLink{}
				if err := json.Unmarshal(linksJSON, &links); err != nil {
					return fmt.Errorf("decode the page links: %w", err)
				}
			}
		}
	}
	client := &http.Client{Timeout: judgeProtocolDeadline}
	for _, link := range links {
		target, err := url.Parse(link.Resolved)
		if err != nil || target.Scheme+"://"+target.Host != origin {
			continue // another origin or a non-HTTP scheme: never fetched
		}
		target.Fragment = ""
		response, err := client.Get(target.String())
		if err != nil {
			return fmt.Errorf("fetch %s: %w", target, err)
		}
		_, _ = io.Copy(io.Discard, response.Body)
		response.Body.Close()
		if response.StatusCode >= 400 {
			report.Findings = append(report.Findings, judgeFinding{Kind: judgeCodeBrokenLink,
				Attribute: link.Attribute, Target: link.Value, Status: response.StatusCode})
		}
	}
	return nil
}

func judgeResult(report judgeReport, output string) resultmodel.CommandResult {
	judgePath := filepath.Join(output, "judge.json")
	encoded, err := json.MarshalIndent(report, "", "  ")
	if err == nil {
		err = os.WriteFile(judgePath, append(encoded, '\n'), 0o644)
	}
	if err != nil {
		return judgeFailureResult(judgeCodeError, "write judge.json: "+err.Error())
	}
	text := fmt.Sprintf("verdict=%s\njudge_json=%s\noutput_directory=%s\n", report.Verdict, judgePath, output)
	result := exactOutputResult(text, []resultmodel.RecordedChange{{Path: judgePath, Kind: "created", Detail: "ai-report-judge verdict " + report.Verdict}})
	for _, finding := range report.Findings {
		evidence := ""
		severity := resultmodel.SeverityError
		switch finding.Kind {
		case judgeCodeOverflow:
			evidence = fmt.Sprintf("%s %s: scroll width %d > inner width %d", finding.Viewport, finding.Scheme, finding.ScrollWidth, finding.InnerWidth)
		case judgeCodeBrokenLink:
			evidence = fmt.Sprintf("%s=%q answered HTTP %d", finding.Attribute, finding.Target, finding.Status)
		default:
			evidence = finding.Detail
			if finding.Kind == judgeCodeSkipped {
				severity = resultmodel.SeverityWarning
			}
		}
		result.Findings = append(result.Findings, toolboxFinding(CommandAIReportJudge, finding.Kind, severity,
			[]string{judgePath}, evidence, resultmodel.FixabilityManual, "the report is not render-verified clean"))
	}
	switch report.Verdict {
	case "pass":
		result.Outcome = resultmodel.OutcomeSuccess
	case "error":
		result.Outcome = resultmodel.OutcomeFailure
	default: // fail and skipped: a skip is never a pass
		result.Outcome = resultmodel.OutcomeFindings
	}
	return result
}

// judgeErrorResult records an engine or protocol failure in judge.json too, so
// the caller can read the cause there.
func judgeErrorResult(report judgeReport, output, detail string) resultmodel.CommandResult {
	report.Verdict = "error"
	report.Findings = append(report.Findings, judgeFinding{Kind: judgeCodeError, Detail: detail})
	return judgeResult(report, output)
}

func judgeFailureResult(code, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		toolboxFinding(CommandAIReportJudge, code, resultmodel.SeverityError, nil, evidence, resultmodel.FixabilityManual, "the render check could not run"),
	}}
}

func judgeAbsolutePath(root, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	return filepath.Clean(path)
}

// judgeResolveExistingPrefix resolves symlinks on the longest existing ancestor
// of path, so an --out that does not exist yet is compared by where it would land.
func judgeResolveExistingPrefix(path string) string {
	missing := ""
	for current := path; ; current = filepath.Dir(current) {
		if resolved, err := filepath.EvalSymlinks(current); err == nil {
			return filepath.Join(resolved, missing)
		}
		if filepath.Dir(current) == current {
			return path
		}
		missing = filepath.Join(filepath.Base(current), missing)
	}
}

func judgePathInside(path, directory string) bool {
	relative, err := filepath.Rel(directory, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func judgeLookupBrowser(name string) string {
	if info, err := os.Stat(name); err == nil && !info.IsDir() {
		return name
	}
	if resolved, err := exec.LookPath(name); err == nil {
		return resolved
	}
	return ""
}

func discoverJudgeBrowser() string {
	for _, name := range aiReportJudgeBrowserNames {
		if resolved, err := exec.LookPath(name); err == nil {
			return resolved
		}
	}
	for _, path := range aiReportJudgeBrowserAppPaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}
	return ""
}

// judgeBrowserSession is one headless engine driven over --remote-debugging-pipe:
// the engine reads NUL-terminated JSON commands on fd 3 and writes replies and
// events on fd 4. A small production copy of the board's browser test transport.
type judgeBrowserSession struct {
	browserCommand   *exec.Cmd
	commandWriter    *os.File
	replyReader      *os.File
	replyBuffer      *bufio.Reader
	profileDirectory string
	pageSessionID    string
	nextCommandID    int
}

func startJudgeBrowser(browser string) (*judgeBrowserSession, error) {
	profileDirectory, err := os.MkdirTemp("", "ai-report-judge-profile-")
	if err != nil {
		return nil, fmt.Errorf("create the browser profile: %w", err)
	}
	commandReader, commandWriter, err := os.Pipe()
	if err != nil {
		os.RemoveAll(profileDirectory)
		return nil, fmt.Errorf("open the command pipe: %w", err)
	}
	replyReader, replyWriter, err := os.Pipe()
	if err != nil {
		commandReader.Close()
		commandWriter.Close()
		os.RemoveAll(profileDirectory)
		return nil, fmt.Errorf("open the reply pipe: %w", err)
	}
	browserCommand := exec.Command(browser, "--headless", "--disable-gpu", "--no-sandbox", "--disable-dev-shm-usage",
		"--hide-scrollbars", "--no-first-run", "--user-data-dir="+profileDirectory, "--remote-debugging-pipe", "about:blank")
	browserCommand.ExtraFiles = []*os.File{commandReader, replyWriter}
	startError := browserCommand.Start()
	// The child holds its own copies; ours close either way so EOF reaches it.
	commandReader.Close()
	replyWriter.Close()
	if startError != nil {
		commandWriter.Close()
		replyReader.Close()
		os.RemoveAll(profileDirectory)
		return nil, fmt.Errorf("start the browser %s: %w", browser, startError)
	}
	return &judgeBrowserSession{browserCommand: browserCommand, commandWriter: commandWriter, replyReader: replyReader,
		replyBuffer: bufio.NewReader(replyReader), profileDirectory: profileDirectory}, nil
}

// stop closes the command pipe (the engine exits on EOF and reaps its own
// helpers), kills it only if it does not exit in time, waits for it, and removes
// the profile once nothing writes there.
func (session *judgeBrowserSession) stop() {
	session.commandWriter.Close()
	exited := make(chan error, 1)
	go func() { exited <- session.browserCommand.Wait() }()
	select {
	case <-exited:
	case <-time.After(judgeProtocolDeadline):
		_ = session.browserCommand.Process.Kill()
		<-exited
	}
	session.replyReader.Close()
	os.RemoveAll(session.profileDirectory)
}

func (session *judgeBrowserSession) call(method string, params map[string]any, pageScoped bool) (json.RawMessage, error) {
	session.nextCommandID++
	command := map[string]any{"id": session.nextCommandID, "method": method}
	if params != nil {
		command["params"] = params
	}
	if pageScoped {
		command["sessionId"] = session.pageSessionID
	}
	encoded, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	if _, err := session.commandWriter.Write(append(encoded, 0)); err != nil {
		return nil, fmt.Errorf("send %s: %w", method, err)
	}
	if err := session.replyReader.SetReadDeadline(time.Now().Add(judgeProtocolDeadline)); err != nil {
		return nil, err
	}
	for {
		message, err := session.replyBuffer.ReadBytes(0)
		if err != nil {
			return nil, fmt.Errorf("no reply to %s within %s: %w", method, judgeProtocolDeadline, err)
		}
		var reply struct {
			ID     int             `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  json.RawMessage `json:"error"`
		}
		if err := json.Unmarshal(bytes.TrimRight(message, "\x00"), &reply); err != nil {
			return nil, fmt.Errorf("undecodable protocol message after %s: %w", method, err)
		}
		if reply.ID != session.nextCommandID {
			continue // an event or another reply
		}
		if reply.Error != nil {
			return nil, fmt.Errorf("%s failed: %s", method, reply.Error)
		}
		return reply.Result, nil
	}
}

func (session *judgeBrowserSession) attachToPage() error {
	deadline := time.Now().Add(judgeProtocolDeadline)
	for {
		targetsJSON, err := session.call("Target.getTargets", nil, false)
		if err != nil {
			return err
		}
		var targets struct {
			TargetInfos []struct {
				TargetID string `json:"targetId"`
				Type     string `json:"type"`
			} `json:"targetInfos"`
		}
		if err := json.Unmarshal(targetsJSON, &targets); err != nil {
			return fmt.Errorf("decode the browser targets: %w", err)
		}
		for _, target := range targets.TargetInfos {
			if target.Type != "page" {
				continue
			}
			attachJSON, err := session.call("Target.attachToTarget", map[string]any{"targetId": target.TargetID, "flatten": true}, false)
			if err != nil {
				return err
			}
			var attachment struct {
				SessionID string `json:"sessionId"`
			}
			if err := json.Unmarshal(attachJSON, &attachment); err != nil || attachment.SessionID == "" {
				return fmt.Errorf("attaching to the page target returned no session id")
			}
			session.pageSessionID = attachment.SessionID
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("the browser opened no page target within %s", judgeProtocolDeadline)
		}
		time.Sleep(judgePollInterval)
	}
}

// evaluate runs one expression in the page and returns its value as JSON.
func (session *judgeBrowserSession) evaluate(expression string) (json.RawMessage, error) {
	evaluationJSON, err := session.call("Runtime.evaluate",
		map[string]any{"expression": expression, "returnByValue": true, "awaitPromise": true}, true)
	if err != nil {
		return nil, err
	}
	var evaluation struct {
		Result struct {
			Value json.RawMessage `json:"value"`
		} `json:"result"`
		ExceptionDetails json.RawMessage `json:"exceptionDetails"`
	}
	if err := json.Unmarshal(evaluationJSON, &evaluation); err != nil {
		return nil, err
	}
	if evaluation.ExceptionDetails != nil {
		return nil, fmt.Errorf("the page threw evaluating the judge script: %s", evaluation.ExceptionDetails)
	}
	if evaluation.Result.Value == nil {
		return nil, errors.New("the judge script produced no value")
	}
	return evaluation.Result.Value, nil
}

// renderPass loads the page at one viewport and colour scheme, measures its
// widths and writes a full-page capture named <viewport>-<scheme>.png.
func (session *judgeBrowserSession) renderPass(pageURL, viewport string, width, height int, scheme, output string) (judgeCapture, error) {
	capture := judgeCapture{Viewport: viewport, Width: width, Height: height, Scheme: scheme, File: viewport + "-" + scheme + ".png"}
	if _, err := session.call("Emulation.setDeviceMetricsOverride",
		map[string]any{"width": width, "height": height, "deviceScaleFactor": 1, "mobile": false}, true); err != nil {
		return capture, err
	}
	if _, err := session.call("Emulation.setEmulatedMedia",
		map[string]any{"features": []map[string]string{{"name": "prefers-color-scheme", "value": scheme}}}, true); err != nil {
		return capture, err
	}
	// The marker lives on the old document's window, so its absence proves the
	// new document loaded; the same URL is navigated four times.
	if _, err := session.evaluate(`window.aiReportJudgeStale = true`); err != nil {
		return capture, err
	}
	navigationJSON, err := session.call("Page.navigate", map[string]any{"url": pageURL}, true)
	if err != nil {
		return capture, err
	}
	var navigation struct {
		ErrorText string `json:"errorText"`
	}
	if json.Unmarshal(navigationJSON, &navigation) == nil && navigation.ErrorText != "" {
		return capture, fmt.Errorf("navigate to %s: %s", pageURL, navigation.ErrorText)
	}
	deadline := time.Now().Add(judgeProtocolDeadline)
	for {
		readyJSON, err := session.evaluate(`!window.aiReportJudgeStale && document.readyState === "complete"`)
		if err == nil && string(readyJSON) == "true" {
			break
		}
		if time.Now().After(deadline) {
			return capture, fmt.Errorf("%s %s pass: the page did not finish loading within %s (last error: %v)", viewport, scheme, judgeProtocolDeadline, err)
		}
		time.Sleep(judgePollInterval)
	}
	measureJSON, err := session.evaluate(`({ScrollWidth: document.documentElement.scrollWidth, InnerWidth: window.innerWidth, ScrollHeight: document.documentElement.scrollHeight})`)
	if err != nil {
		return capture, err
	}
	var measured struct{ ScrollWidth, InnerWidth, ScrollHeight int }
	if err := json.Unmarshal(measureJSON, &measured); err != nil {
		return capture, fmt.Errorf("decode the page measurement: %w", err)
	}
	capture.ScrollWidth, capture.InnerWidth = measured.ScrollWidth, measured.InnerWidth
	screenshotJSON, err := session.call("Page.captureScreenshot", map[string]any{"format": "png", "captureBeyondViewport": true,
		"clip": map[string]any{"x": 0, "y": 0, "scale": 1, "width": max(measured.ScrollWidth, width), "height": max(measured.ScrollHeight, height)}}, true)
	if err != nil {
		return capture, err
	}
	var screenshot struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(screenshotJSON, &screenshot); err != nil {
		return capture, err
	}
	imageBytes, err := base64.StdEncoding.DecodeString(screenshot.Data)
	if err != nil || len(imageBytes) == 0 {
		return capture, fmt.Errorf("%s %s pass: the screenshot was empty or undecodable", viewport, scheme)
	}
	return capture, os.WriteFile(filepath.Join(output, capture.File), imageBytes, 0o644)
}
