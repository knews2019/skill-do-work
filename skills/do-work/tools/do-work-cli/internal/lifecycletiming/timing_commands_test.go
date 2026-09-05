package lifecycletiming

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// The orchestrator reaches this package only through argv, so the whole
// record-then-fold path is exercised at that seam.
func TestTimingCommandsRecordThenFoldThroughArgv(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	requestPath := "do-work/working/REQ-562-record-lifecycle-timings.md"
	writeTimingTestFile(t, repositoryRoot, requestPath, "---\nid: REQ-562\n---\n\n# Body\n")
	clock := useSyntheticClock(t, "2026-09-05T00:34:20Z")
	handlers := Handlers(io.Discard)
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatText}

	clock.advance(600 * time.Second)
	recorded := handlers[CommandRecordTimingEvent](executionContext, []string{
		"--request", "REQ-562", "--run", "work-2026-09-05-003420",
		"--category", "builder-work", "--operation", "builder dispatch to hand-back",
		"--started-at", "2026-09-05T00:34:20Z", "--agent", "implementation-builder",
	})
	if recorded.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("record outcome = %s: %#v", recorded.Outcome, recorded.Findings)
	}
	if recorded.LifecycleTiming == nil || recorded.LifecycleTiming.RecordedEvent == nil {
		t.Fatalf("record result = %#v", recorded.LifecycleTiming)
	}
	if recorded.LifecycleTiming.RecordedEvent.ElapsedSeconds != 600 {
		t.Fatalf("elapsed = %d", recorded.LifecycleTiming.RecordedEvent.ElapsedSeconds)
	}

	clock.advance(300 * time.Second)
	folded := handlers[CommandFoldTimingSummary](executionContext, []string{
		"--request", "REQ-562", "--run", "work-2026-09-05-003420", "--request-path", requestPath,
	})
	if folded.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("fold outcome = %s: %#v", folded.Outcome, folded.Findings)
	}
	if folded.LifecycleTiming == nil || !folded.LifecycleTiming.SectionWritten {
		t.Fatalf("fold result = %#v", folded.LifecycleTiming)
	}
	if len(folded.Changes) != 1 || folded.Changes[0].Path != requestPath {
		t.Fatalf("fold changes = %#v", folded.Changes)
	}
	body := readTimingTestFile(t, repositoryRoot, requestPath)
	if !strings.Contains(body, "## Timing") || !strings.Contains(body, "10m 00s total") {
		t.Fatalf("folded body:\n%s", body)
	}
}

// A wrapped command must keep its own output off the result envelope and return
// the child's status as the process status.
func TestRunTimedCommandArgvKeepsChildOutputOffTheResultEnvelope(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	useSyntheticClock(t, "2026-09-05T00:34:20Z")
	var childOutput bytes.Buffer
	handlers := Handlers(&childOutput)
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatText}

	result := handlers[CommandRunTimedCommand](executionContext, []string{
		"--request", "REQ-562", "--run", "work-2026-09-05-003420",
		"--category", "verification-gate", "--operation", "repository gate",
		"--", "sh", "-c", "printf gate-said-this; exit 7",
	})
	if result.Outcome != resultmodel.OutcomeFindings || result.ExitCodeOverride != 7 {
		t.Fatalf("result = %s exit=%d", result.Outcome, result.ExitCodeOverride)
	}
	if childOutput.String() != "gate-said-this" {
		t.Fatalf("child output = %q", childOutput.String())
	}
	rendered, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if bytes.Contains(rendered, []byte("gate-said-this")) {
		t.Fatalf("child output leaked into the result envelope:\n%s", rendered)
	}
	if !bytes.Contains(rendered, []byte("sh (3 argv tokens)")) {
		t.Fatalf("rendered result lost the command identity:\n%s", rendered)
	}
}

// Every usage mistake must refuse before anything is written, and name the
// runnable argv that fixes it.
func TestTimingCommandsRefuseIncompleteArgv(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	handlers := Handlers(io.Discard)
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatText}

	missingCategory := handlers[CommandRecordTimingEvent](executionContext, []string{"--request", "REQ-562", "--operation", "x"})
	if missingCategory.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("missing category outcome = %s", missingCategory.Outcome)
	}
	missingArgv := handlers[CommandRunTimedCommand](executionContext, []string{
		"--request", "REQ-562", "--category", "verification-gate", "--operation", "gate",
	})
	if missingArgv.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("missing argv outcome = %s", missingArgv.Outcome)
	}
	missingRequestPath := handlers[CommandFoldTimingSummary](executionContext, []string{"--request", "REQ-562"})
	if missingRequestPath.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("missing request path outcome = %s", missingRequestPath.Outcome)
	}
}

// The shipped prose promises the wrapper exits with the child's own status, so
// the process status a caller reads must carry a signal and a failed launch too.
func TestRunTimedCommandArgvCarriesSignalAndLaunchStatusesToTheProcessStatus(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	handlers := Handlers(io.Discard)
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatText}

	signalled := handlers[CommandRunTimedCommand](executionContext, []string{
		"--request", "REQ-562", "--run", "work-argv-signals", "--category", "verification-gate",
		"--operation", "self-terminating probe", "--", "sh", "-c", "kill -TERM $$",
	})
	if signalled.ExitCodeOverride != 143 {
		t.Fatalf("signalled exit override = %d, want 143", signalled.ExitCodeOverride)
	}
	missing := handlers[CommandRunTimedCommand](executionContext, []string{
		"--request", "REQ-562", "--run", "work-argv-signals", "--category", "verification-gate",
		"--operation", "missing executable", "--", "do-work-executable-that-does-not-exist",
	})
	if missing.Outcome != resultmodel.OutcomeFailure || missing.ExitCodeOverride != 127 {
		t.Fatalf("missing executable = %s exit=%d, want failure with 127", missing.Outcome, missing.ExitCodeOverride)
	}
}

// The fold is reachable from argv, so its repository-root confinement has to
// hold at that seam and not only inside the writer.
func TestFoldTimingSummaryArgvRefusesAnEscapingRequestPath(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	handlers := Handlers(io.Discard)
	executionContext := commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatText}

	refused := handlers[CommandFoldTimingSummary](executionContext, []string{
		"--request", "REQ-562", "--run", "work-smoke", "--request-path", "../outside-target.md",
	})
	if refused.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("escaping request path outcome = %s", refused.Outcome)
	}
}

func TestInvalidTimingOptionsNeverLaunchChild(t *testing.T) {
	for _, testCase := range []struct{ name, requestID, runID, category, operation string }{
		{"category", "REQ-562", "work-test", "unknown-category", "gate"},
		{"request", "../REQ-562", "work-test", "verification-gate", "gate"},
		{"missing request", "", "work-test", "verification-gate", "gate"},
		{"run", "REQ-562", "../outside", "verification-gate", "gate"},
		{"operation", "REQ-562", "work-test", "verification-gate", "\t"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			repositoryRoot := newTimingRepository(t)
			arguments := []string{"--run", testCase.runID, "--category", testCase.category, "--operation", testCase.operation}
			if testCase.requestID != "" {
				arguments = append(arguments, "--request", testCase.requestID)
			}
			arguments = append(arguments, "--", "sh", "-c", "echo ran > child-output")
			result := handleRunTimedCommand(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, arguments, io.Discard)
			if result.Outcome != resultmodel.OutcomeFailure {
				t.Fatalf("invalid invocation accepted: %#v", result)
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, "child-output")); !os.IsNotExist(err) {
				t.Fatalf("child ran before rejection: %v", err)
			}
			if len(result.Findings) != 1 || result.Findings[0].Code != "TIMING-USAGE" {
				t.Fatalf("wrong validation error: %#v", result.Findings)
			}
		})
	}
}

func TestTimingRecordingFailurePreservesChildExitAndDoesNotReportLaunchFailure(t *testing.T) {
	for _, childStatus := range []int{0, 7, 127} {
		t.Run(fmt.Sprintf("exit=%d", childStatus), func(t *testing.T) {
			repositoryRoot := newTimingRepository(t)
			streamPath, pathError := streamPathFor(repositoryRoot, "work-test", "REQ-562")
			if pathError != nil {
				t.Fatal(pathError)
			}
			if err := os.MkdirAll(streamPath, 0700); err != nil {
				t.Fatal(err)
			}
			var outputBuffer bytes.Buffer
			commandRuntime := commandruntime.NewRuntime(&outputBuffer, Handlers(io.Discard))
			exitStatus := commandRuntime.Run([]string{
				"--repo-root", repositoryRoot, "--format", "json", CommandRunTimedCommand,
				"--request", "REQ-562", "--run", "work-test", "--category", "verification-gate",
				"--", "sh", "-c", fmt.Sprintf("echo ran > child-output; exit %d", childStatus),
			})
			if exitStatus != childStatus {
				t.Fatalf("wrapper exit=%d, child exit=%d: %s", exitStatus, childStatus, outputBuffer.String())
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, "child-output")); err != nil {
				t.Fatalf("child did not run: %v", err)
			}
			outputText := outputBuffer.String()
			if !strings.Contains(outputText, "TIMED-COMMAND-RECORDING-FAILED") || strings.Contains(outputText, "TIMED-COMMAND-LAUNCH-FAILED") || !strings.Contains(outputText, fmt.Sprintf("exited %d", childStatus)) {
				t.Fatalf("recording failure hid observed execution: %s", outputText)
			}
		})
	}
}
