package lifecycletiming

import (
	"bytes"
	"io"
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
