package lifecycletiming

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// The stream must be flat, append-safe, Git-private, and shared by a repository's
// worktrees: two boundary events land as two independent JSON lines under the Git
// common directory, and the second event's default start is the first event's end
// so a serial lifecycle chains without any caller supplying a timestamp.
func TestAppendTimingEventWritesFlatGitPrivateStreamAndChainsDefaultStart(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	clock := useSyntheticClock(t, "2026-09-05T00:34:20Z")

	clock.advance(90 * time.Second)
	first, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "recovery-selection", Operation: "queue selection and claim",
		StartedAt: mustParseTiming(t, "2026-09-05T00:34:20Z"), HasExplicitStart: true,
		ResponsibleAgent: "orchestrator",
	})
	if err != nil {
		t.Fatalf("first append: %v", err)
	}
	if first.StreamState != "appended" || first.EventCount != 1 {
		t.Fatalf("first = %#v", first)
	}

	clock.advance(300 * time.Second)
	second, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "planning", Operation: "route C plan",
	})
	if err != nil {
		t.Fatalf("second append: %v", err)
	}
	if second.EventCount != 2 {
		t.Fatalf("second event count = %d", second.EventCount)
	}

	commonDirectory := strings.TrimSpace(runTimingTestGit(t, repositoryRoot, "rev-parse", "--git-common-dir"))
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(repositoryRoot, commonDirectory)
	}
	wantPath := filepath.Join(commonDirectory, "do-work-timing", "work-2026-09-05-003420", "REQ-562.jsonl")
	if first.StreamPath != wantPath {
		t.Fatalf("stream path = %q, want %q", first.StreamPath, wantPath)
	}

	events := readTimingStream(t, wantPath)
	if len(events) != 2 {
		t.Fatalf("stream lines = %d, want 2", len(events))
	}
	if events[0].EventID == events[1].EventID {
		t.Fatalf("event ids are not separable: %q", events[0].EventID)
	}
	if events[0].StartedAt != "2026-09-05T00:34:20Z" || events[0].EndedAt != "2026-09-05T00:35:50Z" || events[0].ElapsedSeconds != 90 {
		t.Fatalf("first event window = %#v", events[0])
	}
	if events[0].ElapsedSource != "wall_clock_difference" || events[0].Outcome != "success" || events[0].ResponsibleAgent != "orchestrator" {
		t.Fatalf("first event evidence = %#v", events[0])
	}
	if events[1].StartedAt != events[0].EndedAt {
		t.Fatalf("second start = %q, want chained %q", events[1].StartedAt, events[0].EndedAt)
	}
	if events[1].ElapsedSeconds != 300 {
		t.Fatalf("second elapsed = %d, want 300", events[1].ElapsedSeconds)
	}
	if events[0].SchemaVersion != 1 || events[1].RunID != "work-2026-09-05-003420" || events[1].RequestID != "REQ-562" {
		t.Fatalf("identity fields = %#v %#v", events[0], events[1])
	}
}

// Timing evidence must stay bounded metadata. A command event keeps its exit
// status and an identity a reader recognizes, and no argument value ever reaches
// the durable stream.
func TestAppendTimingEventRedactsCommandArgumentsAndKeepsExitStatus(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	clock := useSyntheticClock(t, "2026-09-05T01:00:00Z")
	clock.advance(45 * time.Second)

	exitStatus := 2
	recorded, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "verification-gate", Operation: "repository gate",
		ExitStatus:  &exitStatus,
		CommandArgv: []string{"bash", "_dev/tests/maintainer-verify.sh", "--token", "s3cret-value"},
		Outcome:     "failure",
	})
	if err != nil {
		t.Fatalf("append: %v", err)
	}
	if recorded.RecordedEvent == nil || recorded.RecordedEvent.ExitStatus == nil || *recorded.RecordedEvent.ExitStatus != 2 {
		t.Fatalf("recorded = %#v", recorded.RecordedEvent)
	}
	if recorded.RecordedEvent.CommandIdentity != "bash (4 argv tokens)" {
		t.Fatalf("command identity = %q", recorded.RecordedEvent.CommandIdentity)
	}
	streamBytes, err := os.ReadFile(recorded.StreamPath)
	if err != nil {
		t.Fatalf("read stream: %v", err)
	}
	if bytes.Contains(streamBytes, []byte("s3cret-value")) || bytes.Contains(streamBytes, []byte("maintainer-verify.sh")) {
		t.Fatalf("stream leaked command arguments: %s", streamBytes)
	}
}

// The only place one process observes both ends of an event is a command this
// package runs itself, so that path must measure in-process and record the exit
// status without the caller reporting either.
func TestRunTimedCommandMeasuresInProcessAndReportsChildExitStatus(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	useSyntheticClock(t, "2026-09-05T01:10:00Z")

	var childOutput bytes.Buffer
	result, exitStatus, err := RunTimedCommand(repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "verification-gate", Operation: "repository gate",
	}, []string{"sh", "-c", "printf gate-output; exit 3"}, &childOutput)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if exitStatus != 3 {
		t.Fatalf("exit status = %d, want 3", exitStatus)
	}
	if childOutput.String() != "gate-output" {
		t.Fatalf("child output = %q", childOutput.String())
	}
	if result.RecordedEvent == nil {
		t.Fatal("no event recorded")
	}
	event := *result.RecordedEvent
	if event.ElapsedSource != "monotonic_in_process" {
		t.Fatalf("elapsed source = %q", event.ElapsedSource)
	}
	if event.ExitStatus == nil || *event.ExitStatus != 3 || event.Outcome != "failure" {
		t.Fatalf("command event = %#v", event)
	}
	if event.CommandIdentity != "sh (3 argv tokens)" {
		t.Fatalf("command identity = %q", event.CommandIdentity)
	}
}

// The fold is the whole point: one compact section carrying total observed time,
// time by category, the slowest stage and command events, and the wall time no
// event claimed. It must then remove the Git-private raw stream.
func TestFoldTimingSummaryWritesCompactSectionAndRemovesRawStream(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	requestPath := "do-work/working/REQ-562-record-lifecycle-timings.md"
	writeTimingTestFile(t, repositoryRoot, requestPath, "---\nid: REQ-562\n---\n\n# Timings\n\n## What\n\nBody.\n")
	seedDeterministicStream(t, repositoryRoot)

	summary, skipped, err := FoldTimingSummary(repositoryRoot, FoldRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562", RequestPath: requestPath,
	})
	if err != nil {
		t.Fatalf("fold: %v", err)
	}
	if len(skipped) != 0 {
		t.Fatalf("skipped = %#v", skipped)
	}
	if summary.StreamState != "folded" || !summary.SectionWritten {
		t.Fatalf("summary = %#v", summary)
	}
	if summary.TotalObservedSeconds != 5400 || summary.AttributedSeconds != 4500 || summary.UnattributedSeconds != 900 {
		t.Fatalf("totals: total=%d attributed=%d unattributed=%d",
			summary.TotalObservedSeconds, summary.AttributedSeconds, summary.UnattributedSeconds)
	}
	if len(summary.CategoryTotals) != 3 {
		t.Fatalf("category totals = %#v", summary.CategoryTotals)
	}
	if summary.CategoryTotals[0].Category != "builder-work" || summary.CategoryTotals[0].ElapsedSeconds != 2400 {
		t.Fatalf("slowest category = %#v", summary.CategoryTotals[0])
	}
	if summary.SlowestStage == nil || summary.SlowestStage.Operation != "builder dispatch to hand-back" {
		t.Fatalf("slowest stage = %#v", summary.SlowestStage)
	}
	if summary.SlowestCommand == nil || summary.SlowestCommand.Operation != "repository gate" {
		t.Fatalf("slowest command = %#v", summary.SlowestCommand)
	}

	folded := readTimingTestFile(t, repositoryRoot, requestPath)
	wantSection := strings.Join([]string{
		"## Timing",
		"",
		"Observed 2026-09-05T00:34:20Z to 2026-09-05T02:04:20Z: 1h 30m 00s total, 1h 15m 00s attributed across 4 events, 15m 00s unattributed.",
		"",
		"| Category | Elapsed | Events |",
		"| --- | --- | --- |",
		"| builder-work | 40m 00s | 1 |",
		"| verification-gate | 20m 00s | 1 |",
		"| planning | 15m 00s | 2 |",
		"",
		"Slowest stage: builder-work / builder dispatch to hand-back, 40m 00s, outcome success.",
		"Slowest command: verification-gate / repository gate, 20m 00s, exit 0, bash (2 argv tokens).",
		"",
	}, "\n")
	if !strings.HasSuffix(folded, wantSection) {
		t.Fatalf("folded request tail mismatch.\n--- got ---\n%s\n--- want suffix ---\n%s", folded, wantSection)
	}
	if !strings.Contains(folded, "## What\n\nBody.\n") {
		t.Fatalf("fold disturbed existing body:\n%s", folded)
	}
	if _, statError := os.Stat(summary.StreamPath); !os.IsNotExist(statError) {
		t.Fatalf("raw stream survived the fold: %v", statError)
	}
}

// A run that recorded nothing must finalize exactly as it did before this
// feature existed: no section, no failure, no byte changed.
func TestFoldTimingSummaryOnAbsentStreamLeavesTheRequestUnchanged(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	requestPath := "do-work/archive/REQ-448-historical.md"
	historical := "---\nid: REQ-448\nstatus: completed\n---\n\n# Historical\n\nNo timing stream was ever recorded.\n"
	writeTimingTestFile(t, repositoryRoot, requestPath, historical)

	summary, skipped, err := FoldTimingSummary(repositoryRoot, FoldRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-448", RequestPath: requestPath,
	})
	if err != nil {
		t.Fatalf("fold: %v", err)
	}
	if len(skipped) != 0 {
		t.Fatalf("skipped = %#v", skipped)
	}
	if summary.StreamState != "absent" || summary.SectionWritten || summary.EventCount != 0 {
		t.Fatalf("summary = %#v", summary)
	}
	if got := readTimingTestFile(t, repositoryRoot, requestPath); got != historical {
		t.Fatalf("historical request changed:\n%s", got)
	}
}

// A crash between a fold and the finalization commit leaves a section behind, so
// the next fold must replace it rather than stack a second one.
func TestFoldTimingSummaryReplacesAnExistingTimingSection(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	requestPath := "do-work/working/REQ-562-record-lifecycle-timings.md"
	writeTimingTestFile(t, repositoryRoot, requestPath,
		"---\nid: REQ-562\n---\n\n## Timing\n\nStale summary from an interrupted run.\n\n## Lessons Learned\n\nKeep this.\n")
	seedDeterministicStream(t, repositoryRoot)

	if _, _, err := FoldTimingSummary(repositoryRoot, FoldRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562", RequestPath: requestPath,
	}); err != nil {
		t.Fatalf("fold: %v", err)
	}
	folded := readTimingTestFile(t, repositoryRoot, requestPath)
	if strings.Count(folded, "## Timing") != 1 {
		t.Fatalf("expected exactly one Timing section:\n%s", folded)
	}
	if strings.Contains(folded, "Stale summary") {
		t.Fatalf("stale summary survived:\n%s", folded)
	}
	if !strings.Contains(folded, "## Lessons Learned\n\nKeep this.\n") {
		t.Fatalf("fold disturbed a later section:\n%s", folded)
	}
}

// Category and run identity are the two axes a reader separates concurrent work
// on, so both are validated before anything is written.
func TestAppendTimingEventRejectsUnknownCategoryAndUnsafeIdentity(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	useSyntheticClock(t, "2026-09-05T01:00:00Z")

	if _, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "critical-path", Operation: "not a stable category",
	}); err == nil {
		t.Fatal("expected an unknown category to be rejected")
	}
	if _, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "../../escape", RequestID: "REQ-562",
		Category: "planning", Operation: "traversal",
	}); err == nil {
		t.Fatal("expected a traversing run identity to be rejected")
	}
	if _, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "not-a-request",
		Category: "planning", Operation: "bad request id",
	}); err == nil {
		t.Fatal("expected a malformed request id to be rejected")
	}
}

// Concurrent runs and requests must never share a stream.
func TestConcurrentRunsAndRequestsKeepSeparateStreams(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	useSyntheticClock(t, "2026-09-05T01:00:00Z")

	first, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-A", RequestID: "REQ-562", Category: "planning", Operation: "plan",
	})
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	sameRunOtherRequest, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-A", RequestID: "REQ-563", Category: "planning", Operation: "plan",
	})
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	otherRun, err := AppendTimingEvent(repositoryRoot, EventRequest{
		RunIdentifier: "work-B", RequestID: "REQ-562", Category: "planning", Operation: "plan",
	})
	if err != nil {
		t.Fatalf("third: %v", err)
	}
	if first.StreamPath == sameRunOtherRequest.StreamPath || first.StreamPath == otherRun.StreamPath {
		t.Fatalf("streams collided: %q %q %q", first.StreamPath, sameRunOtherRequest.StreamPath, otherRun.StreamPath)
	}
}

// seedDeterministicStream writes the fixture the summary assertions are computed
// from: two planning boundaries, one delegated-agent wait, one command event, and
// a fifteen-minute hole nothing claimed.
func seedDeterministicStream(t *testing.T, repositoryRoot string) {
	t.Helper()
	clock := useSyntheticClock(t, "2026-09-05T00:34:20Z")

	clock.advance(600 * time.Second)
	appendOrFail(t, repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "planning", Operation: "route C plan",
		StartedAt: mustParseTiming(t, "2026-09-05T00:34:20Z"), HasExplicitStart: true,
	})
	clock.advance(2400 * time.Second)
	appendOrFail(t, repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "builder-work", Operation: "builder dispatch to hand-back",
		ResponsibleAgent: "implementation-builder",
	})
	// A fifteen-minute hole: the gate starts later than the previous event ended.
	clock.advance(2100 * time.Second)
	zeroExit := 0
	appendOrFail(t, repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "verification-gate", Operation: "repository gate",
		StartedAt: mustParseTiming(t, "2026-09-05T01:39:20Z"), HasExplicitStart: true,
		ExitStatus: &zeroExit, CommandArgv: []string{"bash", "verify.sh"},
	})
	clock.advance(300 * time.Second)
	appendOrFail(t, repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562",
		Category: "planning", Operation: "finalization manifest",
	})
}

func appendOrFail(t *testing.T, repositoryRoot string, eventRequest EventRequest) {
	t.Helper()
	if _, err := AppendTimingEvent(repositoryRoot, eventRequest); err != nil {
		t.Fatalf("append %s: %v", eventRequest.Operation, err)
	}
}

type syntheticClock struct{ current time.Time }

func (clock *syntheticClock) advance(duration time.Duration) {
	clock.current = clock.current.Add(duration)
}

func useSyntheticClock(t *testing.T, startInstant string) *syntheticClock {
	t.Helper()
	clock := &syntheticClock{current: mustParseTiming(t, startInstant)}
	previous := timingClockNow
	timingClockNow = func() time.Time { return clock.current }
	t.Cleanup(func() { timingClockNow = previous })
	return clock
}

func mustParseTiming(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parse %q: %v", value, err)
	}
	return parsed.UTC()
}

func readTimingStream(t *testing.T, streamPath string) []resultmodel.TimingEventRecord {
	t.Helper()
	contents, err := os.ReadFile(streamPath)
	if err != nil {
		t.Fatalf("read stream: %v", err)
	}
	records := []resultmodel.TimingEventRecord{}
	for _, line := range strings.Split(strings.TrimRight(string(contents), "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record resultmodel.TimingEventRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode %q: %v", line, err)
		}
		records = append(records, record)
	}
	return records
}

func newTimingRepository(t *testing.T) string {
	t.Helper()
	repositoryRoot := t.TempDir()
	runTimingTestGit(t, repositoryRoot, "init", "-q", "-b", "master")
	runTimingTestGit(t, repositoryRoot, "config", "user.email", "lifecycle-timing@example.invalid")
	runTimingTestGit(t, repositoryRoot, "config", "user.name", "Lifecycle Timing")
	writeTimingTestFile(t, repositoryRoot, "project.txt", "initial\n")
	runTimingTestGit(t, repositoryRoot, "add", "project.txt")
	runTimingTestGit(t, repositoryRoot, "commit", "-qm", "initial")
	return repositoryRoot
}

func runTimingTestGit(t *testing.T, repositoryRoot string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return string(output)
}

func writeTimingTestFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readTimingTestFile(t *testing.T, repositoryRoot, relativePath string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

// The synthetic clock can show that a wrapped command is LABELLED as measured in
// process; only the real clock can show that it actually is. This runs a command
// with a known floor and asserts the recorded window covers it, which is the one
// property a chained boundary event could otherwise fake.
func TestRunTimedCommandRealClockCoversTheChildsOwnDuration(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	result, exitStatus, err := RunTimedCommand(repositoryRoot, EventRequest{
		RunIdentifier: "work-real-clock", RequestID: "REQ-562",
		Category: "verification-gate", Operation: "slow probe",
	}, []string{"sh", "-c", "sleep 1"}, nil)
	if err != nil || exitStatus != 0 {
		t.Fatalf("run: exit=%d err=%v", exitStatus, err)
	}
	event := *result.RecordedEvent
	if event.ElapsedSeconds < 1 {
		t.Fatalf("elapsed = %ds, want at least the child's own second", event.ElapsedSeconds)
	}
	startInstant := mustParseTiming(t, event.StartedAt)
	endInstant := mustParseTiming(t, event.EndedAt)
	if endInstant.Sub(startInstant) < time.Second {
		t.Fatalf("recorded window %s to %s is shorter than the child ran", event.StartedAt, event.EndedAt)
	}
}

// A request body may show a fenced example of the section this writer produces.
// Matching that example and deleting everything up to the next heading destroys
// the closing fence and whatever follows it, which is the one failure this fold
// can never have.
func TestFoldTimingSummaryLeavesAFencedTimingExampleIntact(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	requestPath := "do-work/working/REQ-562-record-lifecycle-timings.md"
	body := "---\nid: REQ-562\n---\n\n# Example\n\n## What\n\nThe fold writes a section shaped like this:\n\n```markdown\n## Timing\n\nObserved an example window.\n```\n\nThat paragraph must survive the fold.\n\n## Lessons Learned\n\nKeep this too.\n"
	writeTimingTestFile(t, repositoryRoot, requestPath, body)
	seedDeterministicStream(t, repositoryRoot)

	if _, _, err := FoldTimingSummary(repositoryRoot, FoldRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562", RequestPath: requestPath,
	}); err != nil {
		t.Fatalf("fold: %v", err)
	}
	folded := readTimingTestFile(t, repositoryRoot, requestPath)
	if !strings.Contains(folded, "```markdown\n## Timing\n\nObserved an example window.\n```\n") {
		t.Fatalf("the fenced example was damaged:\n%s", folded)
	}
	if !strings.Contains(folded, "That paragraph must survive the fold.") {
		t.Fatalf("content after the fenced example was destroyed:\n%s", folded)
	}
	if !strings.Contains(folded, "## Lessons Learned\n\nKeep this too.\n") {
		t.Fatalf("a later section was destroyed:\n%s", folded)
	}
	if !strings.Contains(folded, "1h 30m 00s total") {
		t.Fatalf("the real summary was never written:\n%s", folded)
	}
}

// The fence rule is a condition over CommonMark's punctuation class, not a list
// of spellings, so a dialect fence this package has never heard of still hides
// what it wraps. The unpaired `---` above a request's source line is the
// counter-case: it encloses nothing, so it must not hide the real section.
func TestFoldTimingSummaryReadsFencesAsAConditionNotASpelling(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	unknownFence := "do-work/working/REQ-562-unknown-fence.md"
	writeTimingTestFile(t, repositoryRoot, unknownFence,
		"---\nid: REQ-562\n---\n\n:::note\n## Timing\n\nA container-fence example.\n:::\n\nSurvivor paragraph.\n")
	seedDeterministicStream(t, repositoryRoot)
	if _, _, err := FoldTimingSummary(repositoryRoot, FoldRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562", RequestPath: unknownFence,
	}); err != nil {
		t.Fatalf("fold: %v", err)
	}
	folded := readTimingTestFile(t, repositoryRoot, unknownFence)
	if !strings.Contains(folded, "A container-fence example.\n:::\n\nSurvivor paragraph.\n") {
		t.Fatalf("an unknown fence spelling did not hide its example:\n%s", folded)
	}

	sourceRule := "do-work/working/REQ-563-source-rule.md"
	writeTimingTestFile(t, repositoryRoot, sourceRule,
		"---\nid: REQ-563\n---\n\n# Body\n\n---\n*Source: \"do it\"*\n")
	seedDeterministicStream(t, repositoryRoot)
	for attempt := 0; attempt < 2; attempt++ {
		if _, _, err := FoldTimingSummary(repositoryRoot, FoldRequest{
			RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-563", RequestPath: sourceRule,
		}); err != nil {
			t.Fatalf("fold attempt %d: %v", attempt, err)
		}
		seedStreamForRequest(t, repositoryRoot, "REQ-563")
	}
	refolded := readTimingTestFile(t, repositoryRoot, sourceRule)
	if strings.Count(refolded, "## Timing") != 1 {
		t.Fatalf("a thematic break hid the real section, so the fold stacked a second one:\n%s", refolded)
	}
}

// The pipeline appends lesson content after the fold, so a folded file whose
// Timing section is last must still end with a newline or the next writer
// concatenates onto its final line.
func TestFoldTimingSummaryKeepsATrailingNewlineWhenTheSectionIsLast(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	requestPath := "do-work/working/REQ-562-trailing-newline.md"
	writeTimingTestFile(t, repositoryRoot, requestPath, "---\nid: REQ-562\n---\n\n## Timing\n\nStale summary from an interrupted run.\n")
	seedDeterministicStream(t, repositoryRoot)

	if _, _, err := FoldTimingSummary(repositoryRoot, FoldRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562", RequestPath: requestPath,
	}); err != nil {
		t.Fatalf("fold: %v", err)
	}
	folded := readTimingTestFile(t, repositoryRoot, requestPath)
	if !strings.HasSuffix(folded, "\n") || strings.HasSuffix(folded, "\n\n") {
		t.Fatalf("folded file must end with exactly one newline, got %q", folded[max(0, len(folded)-40):])
	}
}

// The fold writes to a path a caller supplies, so it must refuse a path that
// leaves the repository rather than rewriting a file the repository does not own.
func TestFoldTimingSummaryRefusesARequestPathOutsideTheRepository(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	seedDeterministicStream(t, repositoryRoot)
	outsideTarget := filepath.Join(filepath.Dir(repositoryRoot), "outside-target.md")
	if err := os.WriteFile(outsideTarget, []byte("---\nid: REQ-562\n---\n\n# Outside\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(outsideTarget) })

	for _, escapingPath := range []string{"../outside-target.md", outsideTarget, "do-work/../../outside-target.md"} {
		if _, _, err := FoldTimingSummary(repositoryRoot, FoldRequest{
			RunIdentifier: "work-2026-09-05-003420", RequestID: "REQ-562", RequestPath: escapingPath,
		}); err == nil {
			t.Fatalf("fold accepted a request path outside the repository: %q", escapingPath)
		}
	}
	survivor, err := os.ReadFile(outsideTarget)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(survivor), "## Timing") {
		t.Fatalf("the fold wrote outside the repository:\n%s", survivor)
	}
}

// The shipped prose claims the wrapper exits with the child's own status, so a
// signalled child and a command that never launched must report what a shell
// reports rather than collapsing to one number.
func TestRunTimedCommandReportsShellStatusesForSignalsAndMissingExecutables(t *testing.T) {
	repositoryRoot := newTimingRepository(t)
	_, signalledStatus, err := RunTimedCommand(repositoryRoot, EventRequest{
		RunIdentifier: "work-signals", RequestID: "REQ-562",
		Category: "verification-gate", Operation: "self-terminating probe",
	}, []string{"sh", "-c", "kill -TERM $$"}, nil)
	if err != nil {
		t.Fatalf("signalled run: %v", err)
	}
	if signalledStatus != 143 {
		t.Fatalf("signalled status = %d, want 143 (128 + SIGTERM)", signalledStatus)
	}

	_, missingStatus, err := RunTimedCommand(repositoryRoot, EventRequest{
		RunIdentifier: "work-signals", RequestID: "REQ-562",
		Category: "verification-gate", Operation: "missing executable",
	}, []string{"do-work-executable-that-does-not-exist"}, nil)
	if err == nil {
		t.Fatal("expected a launch failure for a missing executable")
	}
	if missingStatus != 127 {
		t.Fatalf("missing executable status = %d, want 127", missingStatus)
	}
}

// seedStreamForRequest writes the same deterministic fixture under another
// request id, so a re-fold has a stream to summarize.
func seedStreamForRequest(t *testing.T, repositoryRoot, requestID string) {
	t.Helper()
	clock := useSyntheticClock(t, "2026-09-05T00:34:20Z")
	clock.advance(600 * time.Second)
	appendOrFail(t, repositoryRoot, EventRequest{
		RunIdentifier: "work-2026-09-05-003420", RequestID: requestID,
		Category: "planning", Operation: "route C plan",
		StartedAt: mustParseTiming(t, "2026-09-05T00:34:20Z"), HasExplicitStart: true,
	})
}
