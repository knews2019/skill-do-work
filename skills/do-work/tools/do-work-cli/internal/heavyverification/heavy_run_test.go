//go:build unix

// Every test in this file executes a real lane, and the runner fails closed on
// a platform that cannot prove descendant process ownership, so the whole file
// is scoped to the platforms where a lane can run at all.

package heavyverification

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const heavyRunTestManifest = `{
  "schema_version": 1,
  "lanes": [
    {
      "id": "green-lane",
      "argv": ["sh", "lanes/green.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/green.sh"}]
    },
    {
      "id": "red-lane",
      "argv": ["sh", "lanes/red.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/red.sh"}]
    },
    {
      "id": "skip-lane",
      "argv": ["sh", "lanes/skip.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/skip.sh"}]
    },
    {
      "id": "skip-then-fail-lane",
      "argv": ["sh", "lanes/skip-then-fail.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/skip-then-fail.sh"}]
    },
    {
      "id": "slow-lane",
      "argv": ["sh", "lanes/slow.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/slow.sh"}]
    }
  ],
  "non_heavy_coverage": [{"kind": "subtree", "path": "docs"}]
}`

// greenLaneMarkerFile is written by the green lane so a test can prove whether
// that lane executed at all.
const greenLaneMarkerFile = "green-lane.ran"

// slowLaneChildFile records the pid of the descendant the slow lane leaves
// behind, so the timeout test can prove the whole group was terminated.
const slowLaneChildFile = "slow-lane-child.pid"

// newHeavyRunRepository commits a manifest whose lanes are real scripts: one
// green lane that takes a measurable second, one red lane, one lane that
// announces a skip, one lane that announces a skip and then fails anyway, and
// one lane that backgrounds a TERM-deaf descendant.
func newHeavyRunRepository(t *testing.T) string {
	t.Helper()
	repositoryRoot, _ := newHeavyTestRepository(t, heavyRunTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "lanes/green.sh", "echo ran > "+greenLaneMarkerFile+"\nsleep 1\nexit 0\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/red.sh", "exit 3\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/skip.sh", "printf 'SKIP: no browser is available\\n'\nexit 0\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/skip-then-fail.sh", "printf 'SKIP: no browser is available\\n'\nexit 4\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/slow.sh", "(trap '' TERM; sleep 30) &\necho $! > "+slowLaneChildFile+"\nwait\n")
	writeHeavyTestFile(t, repositoryRoot, "do-work/queue/REQ-900-orchestrator-bookkeeping.md", "queue state\n")
	commitHeavyTestChanges(t, repositoryRoot, "lane scripts")
	return repositoryRoot
}

// runHeavyLanes discards the lane output tee rather than letting it reach this
// process's stderr. The fixture lanes below print `SKIP:` announcements on
// purpose, and the heavy lane that runs this package's tests is watched for
// exactly that prefix on its own output.
func runHeavyLanes(t *testing.T, repositoryRoot string, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	return runHeavyVerificationLanes(
		commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot},
		append([]string{"--manifest", "heavy-lanes.json"}, arguments...),
		io.Discard,
	)
}

func TestRunLanesRecordsExitStatusAndWallSeconds(t *testing.T) {
	repositoryRoot := newHeavyRunRepository(t)
	headRevision := runHeavyTestGit(t, repositoryRoot, "rev-parse", "HEAD")

	result := runHeavyLanes(t, repositoryRoot, "--lane", "red-lane", "--lane", "green-lane")

	if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerificationRun == nil {
		t.Fatalf("run result = %#v", result)
	}
	run := *result.HeavyVerificationRun
	if run.ExecutionRevision != headRevision {
		t.Fatalf("execution revision = %q, want HEAD %q", run.ExecutionRevision, headRevision)
	}
	if run.ManifestPath != "heavy-lanes.json" {
		t.Fatalf("manifest path = %q", run.ManifestPath)
	}
	if len(run.Lanes) != 2 || run.Lanes[0].LaneID != "green-lane" || run.Lanes[1].LaneID != "red-lane" {
		t.Fatalf("lanes did not run in manifest order: %#v", run.Lanes)
	}
	greenLane := run.Lanes[0]
	if greenLane.ExitStatus != 0 || greenLane.Skipped || greenLane.WallSeconds < 1 {
		t.Fatalf("green lane = %#v", greenLane)
	}
	if !reflect.DeepEqual(greenLane.CommandArgv, []string{"sh", "lanes/green.sh"}) {
		t.Fatalf("green lane argv = %v", greenLane.CommandArgv)
	}
	redLane := run.Lanes[1]
	if redLane.ExitStatus != 3 || redLane.Skipped || redLane.WallSeconds < 0 {
		t.Fatalf("red lane = %#v", redLane)
	}
	if len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-RUN-LANE-RED" || result.Findings[0].Severity != resultmodel.SeverityWarning {
		t.Fatalf("red lane finding = %#v", result.Findings)
	}
	if !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), "red-lane") {
		t.Fatalf("red lane finding does not name the lane: %#v", result.Findings[0])
	}
}

func TestRunLanesMarksSkipFromExplicitSkipLine(t *testing.T) {
	repositoryRoot := newHeavyRunRepository(t)

	result := runHeavyLanes(t, repositoryRoot, "--lane", "skip-lane")

	if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerificationRun == nil {
		t.Fatalf("run result = %#v", result)
	}
	lanes := result.HeavyVerificationRun.Lanes
	if len(lanes) != 1 || !lanes[0].Skipped || lanes[0].ExitStatus != 0 {
		t.Fatalf("skip lane = %#v", lanes)
	}
	if len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-RUN-LANE-SKIPPED" || result.Findings[0].Severity != resultmodel.SeverityWarning {
		t.Fatalf("skip finding = %#v", result.Findings)
	}
	if !strings.Contains(strings.Join(result.Findings[0].Evidence, " "), "SKIP: no browser is available") {
		t.Fatalf("skip finding lost the announcing line: %#v", result.Findings[0])
	}
}

// TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine pins the misreport a
// heavy run hit at revision 6646ba51: the do-work-cli lane exited non-zero with
// a failing test, printed a fixture's `SKIP:` line along the way, and was
// recorded as a lane that did not run. A skip line is an announcement; a
// non-zero exit status is a result, and the result decides.
func TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine(t *testing.T) {
	repositoryRoot := newHeavyRunRepository(t)

	result := runHeavyLanes(t, repositoryRoot, "--lane", "skip-then-fail-lane")

	if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerificationRun == nil {
		t.Fatalf("run result = %#v", result)
	}
	lanes := result.HeavyVerificationRun.Lanes
	if len(lanes) != 1 || lanes[0].Skipped || lanes[0].ExitStatus != 4 {
		t.Fatalf("a lane that ran and failed was recorded as skipped: %#v", lanes)
	}
	if len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-RUN-LANE-RED" {
		t.Fatalf("failed lane finding = %#v", result.Findings)
	}
}

func TestRunLanesRefusesUnknownLaneBeforeExecuting(t *testing.T) {
	repositoryRoot := newHeavyRunRepository(t)

	result := runHeavyLanes(t, repositoryRoot, "--lane", "green-lane", "--lane", "missing-lane")

	if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-RUN-LANE-UNKNOWN" {
		t.Fatalf("unknown lane result = %#v", result)
	}
	if result.HeavyVerificationRun != nil {
		t.Fatalf("unknown lane produced a run: %#v", result.HeavyVerificationRun)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, greenLaneMarkerFile)); !os.IsNotExist(err) {
		t.Fatalf("the known lane executed before the unknown lane was refused (stat err = %v)", err)
	}
}

// TestRunLanesRefusesDirtyTrackedTree pins both halves of the dirty-tree rule:
// queue bookkeeping under do-work/ is never lane input, so it cannot hold the
// drain back, while any other tracked change makes the run unattributable.
func TestRunLanesRefusesDirtyTrackedTree(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		dirtyPath   string
		wantRefusal bool
	}{
		{name: "queue bookkeeping runs", dirtyPath: "do-work/queue/REQ-900-orchestrator-bookkeeping.md"},
		{name: "source change refuses", dirtyPath: "seed.txt", wantRefusal: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			repositoryRoot := newHeavyRunRepository(t)
			writeHeavyTestFile(t, repositoryRoot, testCase.dirtyPath, "locally modified\n")

			result := runHeavyLanes(t, repositoryRoot, "--lane", "red-lane")

			if !testCase.wantRefusal {
				if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerificationRun == nil {
					t.Fatalf("dirt under do-work/ held the run back: %#v", result)
				}
				return
			}
			if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "HEAVY-RUN-DIRTY-TREE" {
				t.Fatalf("dirty tree result = %#v", result)
			}
			if result.HeavyVerificationRun != nil {
				t.Fatalf("dirty tree produced a run: %#v", result.HeavyVerificationRun)
			}
		})
	}
}

func TestRunLanesTerminatesTimedOutLaneGroup(t *testing.T) {
	repositoryRoot := newHeavyRunRepository(t)

	result := runHeavyLanes(t, repositoryRoot, "--lane", "slow-lane", "--lane-timeout-seconds", "1")

	if result.Outcome != resultmodel.OutcomeSuccess || result.HeavyVerificationRun == nil {
		t.Fatalf("timed out run result = %#v", result)
	}
	lanes := result.HeavyVerificationRun.Lanes
	if len(lanes) != 1 || lanes[0].ExitStatus != HeavyLaneTimeoutStatus {
		t.Fatalf("timed out lane = %#v", lanes)
	}
	descendantPID := readSlowLaneChildPID(t, repositoryRoot)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(descendantPID, 0); err != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("descendant %d survived the terminated lane group", descendantPID)
}

func readSlowLaneChildPID(t *testing.T, repositoryRoot string) int {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, slowLaneChildFile))
	if err != nil {
		t.Fatalf("read slow lane child pid: %v", err)
	}
	descendantPID, err := strconv.Atoi(strings.TrimSpace(string(contents)))
	if err != nil {
		t.Fatalf("parse slow lane child pid %q: %v", contents, err)
	}
	return descendantPID
}

// TestRunLanesWithoutLaneTimeoutUsesTheDefaultBound pins the zero-value rule in
// RunLanes. The 1800-second default is chosen by the CLI argument parser, so an
// in-process caller that leaves LaneTimeout unset once armed time.NewTimer(0)
// and had its lane terminated while it was still starting — a green lane
// reported red, and only under load. green-lane sleeps a full second, far past
// an instant bound and far below the default.
func TestRunLanesWithoutLaneTimeoutUsesTheDefaultBound(t *testing.T) {
	repositoryRoot := newHeavyRunRepository(t)

	run, _, err := RunLanes(LaneRunRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: "heavy-lanes.json",
		LaneIDs: []string{"green-lane"}, LaneOutputWriter: io.Discard,
	})
	if err != nil {
		t.Fatalf("run with no LaneTimeout: %v", err)
	}
	if len(run.Lanes) != 1 {
		t.Fatalf("lanes = %#v", run.Lanes)
	}
	lane := run.Lanes[0]
	if lane.ExitStatus != 0 {
		t.Fatalf("lane with no LaneTimeout = %#v, want exit 0; exit %d is the lane-timeout status, meaning the unset field armed an instant bound",
			lane, HeavyLaneTimeoutStatus)
	}
}
