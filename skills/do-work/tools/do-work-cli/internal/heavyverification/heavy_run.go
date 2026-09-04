package heavyverification

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/ownedprocess"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	// HeavyLaneTimeoutStatus is the exit status recorded for a lane whose
	// process group was terminated at the timeout bound.
	HeavyLaneTimeoutStatus = 124
	// laneSkipPrefix is how a lane announces that it did not actually run. The
	// condition is the prefix on a line of the lane's own output, never a list
	// of lane names, so a lane added later inherits it.
	laneSkipPrefix = "SKIP:"
	// queueStatePrefix is the one tree the dirty-tree refusal ignores: queue
	// bookkeeping is never lane input, and the orchestrator routinely carries
	// uncommitted edits there at the moment a drain runs.
	queueStatePrefix           = "do-work/"
	laneTerminationGracePeriod = time.Second
	interruptedLaneStatusFloor = 128
)

// laneRunRefusal carries the typed code a caller reports for a condition that
// stops the run before or instead of executing lanes.
type laneRunRefusal struct {
	Code  string
	Cause error
}

func (refusal laneRunRefusal) Error() string { return refusal.Cause.Error() }

func refuseLaneRun(code, format string, arguments ...any) error {
	return laneRunRefusal{Code: code, Cause: fmt.Errorf(format, arguments...)}
}

// LaneRunRefusalCode reports the typed code for an error RunLanes returned.
// Anything without its own code is unverifiable rather than a known refusal.
func LaneRunRefusalCode(err error) string {
	var refusal laneRunRefusal
	if errors.As(err, &refusal) {
		return refusal.Code
	}
	return "HEAVY-RUN-UNVERIFIABLE"
}

// RunLanes executes the named manifest lanes at HEAD, one at a time in manifest
// order, and records each lane's exit status, skip state, and wall time. It
// does not plan (that is Plan), does not judge the stored plan prose, and does
// not build an answer manifest; a red or skipped lane is evidence the caller
// records rather than a failure of this command. Lane output is teed to
// laneOutputWriter so a caller watching a long drain still sees progress.
func RunLanes(repositoryRoot, manifestPath string, requestedLaneIDs []string, laneTimeout time.Duration, laneOutputWriter io.Writer) (resultmodel.HeavyVerificationRun, []resultmodel.CommandFinding, error) {
	_, manifestRelativePath, err := resolveManifestPath(repositoryRoot, manifestPath)
	if err != nil {
		return resultmodel.HeavyVerificationRun{}, nil, err
	}
	executionRevision, err := resolveRevision(repositoryRoot, "HEAD")
	if err != nil {
		return resultmodel.HeavyVerificationRun{}, nil, fmt.Errorf("resolve repository HEAD: %w", err)
	}
	if err := refuseDirtyTrackedTree(repositoryRoot); err != nil {
		return resultmodel.HeavyVerificationRun{}, nil, err
	}
	manifest, err := readManifestAtRevision(repositoryRoot, manifestRelativePath, executionRevision)
	if err != nil {
		return resultmodel.HeavyVerificationRun{}, nil, err
	}
	orderedLanes, err := selectManifestLanes(manifest, requestedLaneIDs)
	if err != nil {
		return resultmodel.HeavyVerificationRun{}, nil, err
	}
	run := resultmodel.HeavyVerificationRun{
		ManifestPath:      manifestRelativePath,
		ExecutionRevision: executionRevision,
		Lanes:             []resultmodel.HeavyLaneExecution{},
	}
	findings := []resultmodel.CommandFinding{}
	interruptions := make(chan os.Signal, 1)
	signal.Notify(interruptions, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interruptions)
	for _, lane := range orderedLanes {
		execution, skipLine, interrupted, err := runOneLane(repositoryRoot, lane, laneTimeout, laneOutputWriter, interruptions)
		if err != nil {
			return resultmodel.HeavyVerificationRun{}, nil, err
		}
		run.Lanes = append(run.Lanes, execution)
		switch {
		case execution.Skipped:
			findings = append(findings, laneSkippedFinding(execution, skipLine))
		case execution.ExitStatus != 0:
			findings = append(findings, laneRedFinding(execution))
		}
		if interrupted {
			// Remaining lanes are left unrun rather than reported: a lane
			// absent from the run can never be read as a green lane.
			break
		}
	}
	return run, findings, nil
}

// refuseDirtyTrackedTree refuses when a tracked file outside the queue tree
// differs from HEAD, because a lane result cannot be attributed to a revision
// the working tree no longer matches.
func refuseDirtyTrackedTree(repositoryRoot string) error {
	status, err := runGit(repositoryRoot, "status", "--porcelain", "-z", "--untracked-files=no")
	if err != nil {
		return fmt.Errorf("inspect the working tree: %w", err)
	}
	offendingPaths := []string{}
	fields := strings.Split(string(status), "\x00")
	for fieldIndex := 0; fieldIndex < len(fields); fieldIndex++ {
		entry := fields[fieldIndex]
		if len(entry) < 4 {
			continue
		}
		// A rename or copy is followed by its source path in its own field;
		// both ends decide whether this change belongs to the queue tree.
		if entry[0] == 'R' || entry[0] == 'C' || entry[1] == 'R' || entry[1] == 'C' {
			fieldIndex++
			if fieldIndex < len(fields) && !strings.HasPrefix(fields[fieldIndex], queueStatePrefix) {
				offendingPaths = append(offendingPaths, fields[fieldIndex])
			}
		}
		if changedPath := entry[3:]; !strings.HasPrefix(changedPath, queueStatePrefix) {
			offendingPaths = append(offendingPaths, changedPath)
		}
	}
	if len(offendingPaths) > 0 {
		sort.Strings(offendingPaths)
		return refuseLaneRun("HEAVY-RUN-DIRTY-TREE", "commit or stash %s before running heavy lanes; a lane result must be attributable to HEAD", strings.Join(offendingPaths, ", "))
	}
	return nil
}

// selectManifestLanes returns the requested lanes in manifest order, refusing
// before anything executes when the manifest declares no lane by that name.
func selectManifestLanes(manifest laneManifest, requestedLaneIDs []string) ([]manifestLane, error) {
	unmatchedLaneIDs := map[string]bool{}
	for _, laneID := range requestedLaneIDs {
		unmatchedLaneIDs[laneID] = true
	}
	orderedLanes := []manifestLane{}
	for _, lane := range manifest.Lanes {
		if unmatchedLaneIDs[lane.ID] {
			delete(unmatchedLaneIDs, lane.ID)
			orderedLanes = append(orderedLanes, lane)
		}
	}
	if len(unmatchedLaneIDs) > 0 {
		unknownLaneIDs := make([]string, 0, len(unmatchedLaneIDs))
		for laneID := range unmatchedLaneIDs {
			unknownLaneIDs = append(unknownLaneIDs, laneID)
		}
		sort.Strings(unknownLaneIDs)
		return nil, refuseLaneRun("HEAVY-RUN-LANE-UNKNOWN", "the heavy-lane manifest declares no lane named %s", strings.Join(unknownLaneIDs, ", "))
	}
	return orderedLanes, nil
}

// runOneLane executes a single lane inside its own process group. It reports
// the typed execution, the line a skipping lane announced itself with, and
// whether a terminating signal ended the run.
func runOneLane(repositoryRoot string, lane manifestLane, laneTimeout time.Duration, laneOutputWriter io.Writer, interruptions <-chan os.Signal) (resultmodel.HeavyLaneExecution, string, bool, error) {
	command := exec.Command(lane.Argv[0], lane.Argv[1:]...)
	command.Dir = repositoryRoot
	// Fail closed: a lane can spawn a test server or a browser, and terminating
	// only the leader would leave those running past the timeout this bound
	// exists to enforce.
	if !ownedprocess.ConfigureGroup(command) {
		return resultmodel.HeavyLaneExecution{}, "", false, refuseLaneRun("HEAVY-RUN-UNOWNED-PROCESS", "heavy lanes cannot run here: descendant process ownership cannot be proved on this platform")
	}
	skipWatcher := &laneSkipWatcher{downstream: laneOutputWriter}
	// One writer for both streams keeps os/exec on a single pipe, so the
	// watcher sees the lane's output in the order the lane produced it.
	command.Stdout = skipWatcher
	command.Stderr = skipWatcher
	startedAt := time.Now()
	if err := command.Start(); err != nil {
		return resultmodel.HeavyLaneExecution{}, "", false, refuseLaneRun("HEAVY-RUN-LANE-UNLAUNCHABLE", "start heavy lane %s: %v", lane.ID, err)
	}
	waited := make(chan error, 1)
	go func() { waited <- command.Wait() }()
	laneTimer := time.NewTimer(laneTimeout)
	defer laneTimer.Stop()
	interrupted := false
	exitStatus := 0
	select {
	case waitError := <-waited:
		exitStatus = laneExitStatus(waitError)
	case <-laneTimer.C:
		// TerminateGroup returns only once the group is gone, so the next lane
		// never starts beside a descendant of the lane that just timed out.
		_ = ownedprocess.TerminateGroup(command.Process.Pid, laneTerminationGracePeriod)
		<-waited
		exitStatus = HeavyLaneTimeoutStatus
	case received := <-interruptions:
		_ = ownedprocess.TerminateGroup(command.Process.Pid, laneTerminationGracePeriod)
		<-waited
		exitStatus = interruptedLaneStatus(received)
		interrupted = true
	}
	skipWatcher.FlushPendingLine()
	execution := resultmodel.HeavyLaneExecution{
		LaneID:      lane.ID,
		CommandArgv: append([]string(nil), lane.Argv...),
		ExitStatus:  exitStatus,
		Skipped:     skipWatcher.SkipLine() != "",
		WallSeconds: int(time.Since(startedAt).Seconds()),
	}
	return execution, skipWatcher.SkipLine(), interrupted, nil
}

func laneExitStatus(waitError error) int {
	if waitError == nil {
		return 0
	}
	var exitError *exec.ExitError
	if errors.As(waitError, &exitError) {
		if status := exitError.ExitCode(); status >= 0 {
			return status
		}
		return interruptedLaneStatusFloor
	}
	return 1
}

func interruptedLaneStatus(received os.Signal) int {
	if signalNumber, isSignal := received.(syscall.Signal); isSignal {
		return interruptedLaneStatusFloor + int(signalNumber)
	}
	return interruptedLaneStatusFloor + int(syscall.SIGINT)
}

func laneSkippedFinding(execution resultmodel.HeavyLaneExecution, skipLine string) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{
		Code: "HEAVY-RUN-LANE-SKIPPED", Severity: resultmodel.SeverityWarning,
		Evidence:             []string{fmt.Sprintf("heavy lane %s did not run after %ds: %s", execution.LaneID, execution.WallSeconds, skipLine)},
		Fixability:           resultmodel.FixabilityManual,
		AutomationStopReason: "a lane that did not run is not evidence that the lane is green",
		NextArgv:             append([]string(nil), execution.CommandArgv...),
		VerificationArgv:     []string{"do-work-cli", "--format", "json", CommandRunHeavyVerification, "--lane", execution.LaneID},
	}
}

func laneRedFinding(execution resultmodel.HeavyLaneExecution) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{
		Code: "HEAVY-RUN-LANE-RED", Severity: resultmodel.SeverityWarning,
		Evidence:             []string{fmt.Sprintf("heavy lane %s exited %d after %ds", execution.LaneID, execution.ExitStatus, execution.WallSeconds)},
		Fixability:           resultmodel.FixabilityManual,
		AutomationStopReason: "a red heavy lane is a result to record, not a failure of this command",
		NextArgv:             append([]string(nil), execution.CommandArgv...),
		VerificationArgv:     []string{"do-work-cli", "--format", "json", CommandRunHeavyVerification, "--lane", execution.LaneID},
	}
}

// laneSkipWatcher tees a lane's combined output to the caller's writer while
// watching for the line a lane prints when it did not run.
type laneSkipWatcher struct {
	downstream  io.Writer
	pendingLine []byte
	skipLine    string
}

func (watcher *laneSkipWatcher) Write(payload []byte) (int, error) {
	watcher.pendingLine = append(watcher.pendingLine, payload...)
	for {
		newlineIndex := bytes.IndexByte(watcher.pendingLine, '\n')
		if newlineIndex < 0 {
			break
		}
		watcher.inspectLine(string(watcher.pendingLine[:newlineIndex]))
		watcher.pendingLine = append([]byte(nil), watcher.pendingLine[newlineIndex+1:]...)
	}
	if watcher.downstream != nil {
		// The tee is diagnostic only; the skip verdict comes from the scan
		// above, so a writer that rejects a chunk must not stop the lane.
		_, _ = watcher.downstream.Write(payload)
	}
	return len(payload), nil
}

// FlushPendingLine inspects output the lane left without a trailing newline.
func (watcher *laneSkipWatcher) FlushPendingLine() {
	if len(watcher.pendingLine) > 0 {
		watcher.inspectLine(string(watcher.pendingLine))
		watcher.pendingLine = nil
	}
}

func (watcher *laneSkipWatcher) inspectLine(line string) {
	if watcher.skipLine == "" && strings.HasPrefix(strings.TrimSuffix(line, "\r"), laneSkipPrefix) {
		watcher.skipLine = strings.TrimSuffix(line, "\r")
	}
}

func (watcher *laneSkipWatcher) SkipLine() string { return watcher.skipLine }

// defaultLaneTimeoutSeconds bounds one lane at half an hour, well above the
// slowest recorded lane and well below an unattended run hanging overnight.
const defaultLaneTimeoutSeconds = 1800
