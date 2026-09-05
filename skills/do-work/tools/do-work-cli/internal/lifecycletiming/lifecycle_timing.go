// Package lifecycletiming is the one canonical writer of do-work's per-request
// lifecycle timing stream. It owns every timestamp and duration mechanic so no
// action file, hook or script has to reach for `date` and invent a second
// format.
//
// The stream is flat on purpose. Each event names its own category, its own
// window and its own outcome, and nothing about any other event, so appending
// one line can never invalidate another and no parent/child model is needed to
// read it. Aggregation is therefore a sum over independent rows, and the gap
// between rows is honest unattributed wall time rather than something a nesting
// model has to explain away.
//
// Per-test durations are deliberately out of scope here: a wrapped test command
// is one command event, and the project's own test-duration log stays the single
// owner of what happened inside it.
package lifecycletiming

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// timingStreamSchemaVersion is stamped on every appended event so a later reader
// can tell an old row from a new one without guessing from field presence.
const timingStreamSchemaVersion = 1

// timingDirectoryName sits under the Git common directory, so every worktree of
// one repository shares one stream root and nothing is ever committed.
const timingDirectoryName = "do-work-timing"

const timingSectionHeading = "## Timing"

// fenceRunMinimumLength is CommonMark's own floor for a fence or thematic break.
const fenceRunMinimumLength = 3

// Shell-compatible statuses for a wrapped command that did not exit normally.
const (
	timedCommandNotLaunchedStatus = 127
	signalExitStatusBase          = 128
)

// Elapsed sources distinguish a duration this process measured end to end from
// one derived by subtracting two wall-clock instants.
const (
	elapsedSourceMonotonic = "monotonic_in_process"
	elapsedSourceWallClock = "wall_clock_difference"
)

// Stream states are the typed answer to "what did this call do to the stream".
const (
	streamStateAppended = "appended"
	streamStateAbsent   = "absent"
	streamStateFolded   = "folded"
)

// timingClockNow is this package's single clock. Production reads the real wall
// clock, whose readings also carry Go's monotonic reading, so a start and end
// observed inside one process subtract monotonically. Tests replace it with a
// synthetic clock to make durations, aggregation and the folded summary exact.
var timingClockNow = time.Now

// TimingCategoryVocabulary is the stable closed set of lifecycle stages an event
// may be filed under. It mirrors the work pipeline's own major boundaries; an
// unrecognized value is refused rather than normalized, because a category
// nobody sums is worse than no event at all.
var TimingCategoryVocabulary = []string{
	"recovery-selection",
	"planning",
	"exploration-preflight",
	"builder-work",
	"handback-merge",
	"verification-gate",
	"review",
	"remediation",
	"finalization",
	"cleanup",
}

// TimingOutcomeVocabulary is the closed set of event outcomes.
var TimingOutcomeVocabulary = []string{"success", "failure", "refused", "skipped"}

var (
	timingRequestIDPattern = regexp.MustCompile(`^REQ-[0-9]+$`)
	// A run identity becomes a directory name, so it must be a single safe path
	// segment: leading alphanumeric rules out "." and ".." and any leading dash.
	timingRunIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)
)

// evidenceTextLimit bounds every free-text field that reaches durable evidence.
const evidenceTextLimit = 120

// EventRequest is one caller's description of a completed interval. Everything
// clock-shaped is derived here, never supplied, except StartedAt — which a
// caller may pin to a milestone it already recorded so a gap before it stays
// visible as unattributed time.
type EventRequest struct {
	RunIdentifier    string
	RequestID        string
	Category         string
	Operation        string
	StartedAt        time.Time
	HasExplicitStart bool
	Outcome          string
	Revision         string
	ResponsibleAgent string
	ExitStatus       *int
	CommandArgv      []string
}

// FoldRequest names the stream to summarize and the request file that receives
// the compact `## Timing` section.
type FoldRequest struct {
	RunIdentifier string
	RequestID     string
	RequestPath   string
}

// AppendTimingEvent records one completed interval and returns the typed
// projection of the append.
func AppendTimingEvent(repositoryRoot string, eventRequest EventRequest) (resultmodel.LifecycleTimingResult, error) {
	endInstant := timingClockNow().UTC()
	return appendMeasuredEvent(repositoryRoot, eventRequest, endInstant, elapsedSourceWallClock, nil)
}

var (
	errTimedCommandInput = errors.New("invalid timed command input")
	errTimingRecording   = errors.New("timing recording failed")
)

// RunTimedCommand executes argv, measures it in-process, and records one command
// event. It returns the child's exit status alongside the typed projection. This
// is the only path where one process observes both ends of an event, so it is
// the only path that reports a monotonic elapsed source.
func RunTimedCommand(repositoryRoot string, eventRequest EventRequest, argv []string, childOutput io.Writer) (resultmodel.LifecycleTimingResult, int, error) {
	if len(argv) == 0 {
		return resultmodel.LifecycleTimingResult{}, 0, fmt.Errorf("%w: a timed command requires argv after --", errTimedCommandInput)
	}
	if eventRequest.Operation == "" {
		eventRequest.Operation = redactCommandIdentity(argv)
	}
	// The child supplies its outcome; only caller-owned inputs can be checked now.
	eventRequest.Outcome = "success"
	validatedRequest, validationError := validateEventRequest(eventRequest)
	if validationError != nil {
		return resultmodel.LifecycleTimingResult{}, 0, fmt.Errorf("%w: %v", errTimedCommandInput, validationError)
	}
	eventRequest = validatedRequest
	if childOutput == nil {
		childOutput = io.Discard
	}
	command := exec.Command(argv[0], argv[1:]...)
	command.Dir = repositoryRoot
	command.Stdout = childOutput
	command.Stderr = childOutput

	startInstant := timingClockNow()
	runError := command.Run()
	endInstant := timingClockNow()

	var exitError *exec.ExitError
	switch {
	case runError == nil:
	case errors.As(runError, &exitError):
	default:
		return resultmodel.LifecycleTimingResult{}, timedCommandNotLaunchedStatus,
			fmt.Errorf("launch %s: %w", filepath.Base(argv[0]), runError)
	}
	exitStatus := childExitStatus(command.ProcessState)

	eventRequest.CommandArgv = argv
	eventRequest.ExitStatus = &exitStatus
	eventRequest.Outcome = "success"
	if exitStatus != 0 {
		eventRequest.Outcome = "failure"
	}
	// The start is derived from the measured elapsed rather than from the
	// previous event's end, so wall time spent before the launch stays
	// unattributed instead of being silently folded into the command.
	measuredElapsed := endInstant.Sub(startInstant)
	eventRequest.StartedAt = endInstant.UTC().Add(-measuredElapsed)
	eventRequest.HasExplicitStart = true

	elapsedSeconds := int(measuredElapsed.Round(time.Second) / time.Second)
	result, err := appendMeasuredEvent(repositoryRoot, eventRequest, endInstant.UTC(), elapsedSourceMonotonic, &elapsedSeconds)
	if err != nil {
		return result, exitStatus, fmt.Errorf("%w: command exited %d; %v", errTimingRecording, exitStatus, err)
	}
	return result, exitStatus, nil
}

// childExitStatus reports what a shell reports for this child, which is what the
// shipped prose promises: the child's own code when it exited normally, and 128
// plus the signal number when a signal killed it. os.ProcessState offers no
// portable signal accessor, so the reading is taken through an anonymous
// interface that only a platform with signals satisfies; a platform without them
// keeps the bare base, still non-zero and still not a code a child chose. A
// command that never launched is reported by the caller as 127, the shell's own
// "command not found".
func childExitStatus(state *os.ProcessState) int {
	if code := state.ExitCode(); code >= 0 {
		return code
	}
	if waitStatus, readable := state.Sys().(interface {
		Signaled() bool
		Signal() syscall.Signal
	}); readable && waitStatus.Signaled() {
		return signalExitStatusBase + int(waitStatus.Signal())
	}
	return signalExitStatusBase
}

// validateEventRequest shares caller-input validation between measured commands
// and externally recorded events. It does no I/O and cannot run a child.
func validateEventRequest(eventRequest EventRequest) (EventRequest, error) {
	runIdentifier := eventRequest.RunIdentifier
	if runIdentifier == "" {
		runIdentifier = "standalone"
	}
	if !timingRunIDPattern.MatchString(runIdentifier) {
		return EventRequest{}, fmt.Errorf("run identity %q is not a single safe path segment", runIdentifier)
	}
	if !timingRequestIDPattern.MatchString(eventRequest.RequestID) {
		return EventRequest{}, fmt.Errorf("request identity %q is not REQ-NNN", eventRequest.RequestID)
	}
	if !vocabularyContains(TimingCategoryVocabulary, eventRequest.Category) {
		return EventRequest{}, fmt.Errorf("category %q is not one of %s", eventRequest.Category, strings.Join(TimingCategoryVocabulary, ", "))
	}
	operation := boundEvidenceText(eventRequest.Operation)
	if operation == "" {
		return EventRequest{}, errors.New("an operation name is required")
	}
	outcome := eventRequest.Outcome
	if outcome == "" {
		outcome = "success"
	}
	if !vocabularyContains(TimingOutcomeVocabulary, outcome) {
		return EventRequest{}, fmt.Errorf("outcome %q is not one of %s", outcome, strings.Join(TimingOutcomeVocabulary, ", "))
	}

	eventRequest.RunIdentifier = runIdentifier
	eventRequest.Operation = operation
	eventRequest.Outcome = outcome
	return eventRequest, nil
}

// appendMeasuredEvent is the single write path. Every field is validated and
// bounded before one JSON line is appended to the stream.
func appendMeasuredEvent(repositoryRoot string, eventRequest EventRequest, endInstant time.Time, elapsedSource string, measuredSeconds *int) (resultmodel.LifecycleTimingResult, error) {
	validatedRequest, validationError := validateEventRequest(eventRequest)
	if validationError != nil {
		return resultmodel.LifecycleTimingResult{}, validationError
	}
	eventRequest = validatedRequest
	runIdentifier, operation, outcome := eventRequest.RunIdentifier, eventRequest.Operation, eventRequest.Outcome

	streamPath, err := streamPathFor(repositoryRoot, runIdentifier, eventRequest.RequestID)
	if err != nil {
		return resultmodel.LifecycleTimingResult{}, err
	}
	existing, _, err := readStream(streamPath)
	if err != nil {
		return resultmodel.LifecycleTimingResult{}, err
	}

	startInstant := eventRequest.StartedAt.UTC()
	if !eventRequest.HasExplicitStart {
		startInstant = defaultStartInstant(existing, endInstant)
	}
	if startInstant.After(endInstant) {
		return resultmodel.LifecycleTimingResult{}, fmt.Errorf("event start %s is after its end %s",
			formatTimingInstant(startInstant), formatTimingInstant(endInstant))
	}
	elapsedSeconds := int(endInstant.Sub(startInstant) / time.Second)
	if measuredSeconds != nil {
		elapsedSeconds = *measuredSeconds
	}

	record := resultmodel.TimingEventRecord{
		SchemaVersion:    timingStreamSchemaVersion,
		EventID:          fmt.Sprintf("%04d", len(existing)+1),
		RunID:            runIdentifier,
		RequestID:        eventRequest.RequestID,
		Category:         eventRequest.Category,
		Operation:        operation,
		StartedAt:        formatTimingInstant(startInstant),
		EndedAt:          formatTimingInstant(endInstant),
		ElapsedSeconds:   elapsedSeconds,
		ElapsedSource:    elapsedSource,
		Outcome:          outcome,
		Revision:         boundEvidenceText(eventRequest.Revision),
		ResponsibleAgent: boundEvidenceText(eventRequest.ResponsibleAgent),
		ExitStatus:       eventRequest.ExitStatus,
		CommandIdentity:  redactCommandIdentity(eventRequest.CommandArgv),
	}
	if err := appendStreamLine(streamPath, record); err != nil {
		return resultmodel.LifecycleTimingResult{}, err
	}
	return resultmodel.LifecycleTimingResult{
		RequestID: record.RequestID, RunID: record.RunID, StreamPath: streamPath,
		StreamState: streamStateAppended, RecordedEvent: &record, EventCount: len(existing) + 1,
	}, nil
}

// defaultStartInstant chains a boundary event to the previous event's end, which
// is what makes a serial lifecycle recordable without any caller holding a
// timestamp. An empty stream starts its first event at the same instant it ends,
// so a caller that forgot to pin a start records a zero rather than a fiction.
func defaultStartInstant(existing []resultmodel.TimingEventRecord, endInstant time.Time) time.Time {
	for index := len(existing) - 1; index >= 0; index-- {
		if parsed, err := time.Parse(time.RFC3339, existing[index].EndedAt); err == nil {
			return parsed.UTC()
		}
	}
	return endInstant
}

// FoldTimingSummary summarizes the stream, writes one `## Timing` section into
// the request file, and removes the Git-private raw stream. A stream that was
// never written is not an error: the request keeps exactly the bytes it had,
// which is how every historical request and every uninstrumented run behaves.
func FoldTimingSummary(repositoryRoot string, foldRequest FoldRequest) (resultmodel.LifecycleTimingResult, []resultmodel.SkippedWork, error) {
	runIdentifier := foldRequest.RunIdentifier
	if runIdentifier == "" {
		runIdentifier = "standalone"
	}
	if !timingRunIDPattern.MatchString(runIdentifier) {
		return resultmodel.LifecycleTimingResult{}, nil, fmt.Errorf("run identity %q is not a single safe path segment", runIdentifier)
	}
	if !timingRequestIDPattern.MatchString(foldRequest.RequestID) {
		return resultmodel.LifecycleTimingResult{}, nil, fmt.Errorf("request identity %q is not REQ-NNN", foldRequest.RequestID)
	}
	confinedRequestPath, pathError := repositoryRelativeRequestPath(foldRequest.RequestPath)
	if pathError != nil {
		return resultmodel.LifecycleTimingResult{}, nil, pathError
	}
	foldRequest.RequestPath = confinedRequestPath
	streamPath, err := streamPathFor(repositoryRoot, runIdentifier, foldRequest.RequestID)
	if err != nil {
		return resultmodel.LifecycleTimingResult{}, nil, err
	}
	records, unreadableLines, err := readStream(streamPath)
	if err != nil {
		return resultmodel.LifecycleTimingResult{}, nil, err
	}
	summary := summarizeStream(records)
	summary.RequestID = foldRequest.RequestID
	summary.RunID = runIdentifier
	summary.StreamPath = streamPath
	summary.RequestPath = foldRequest.RequestPath

	skipped := []resultmodel.SkippedWork{}
	if unreadableLines > 0 {
		skipped = append(skipped, resultmodel.SkippedWork{
			Code:   "TIMING-STREAM-LINE-UNREADABLE",
			Reason: fmt.Sprintf("%d stream line(s) did not decode and were left out of the summary", unreadableLines),
		})
	}
	if len(records) == 0 {
		summary.StreamState = streamStateAbsent
		return summary, skipped, nil
	}

	absoluteRequestPath := filepath.Join(repositoryRoot, filepath.FromSlash(foldRequest.RequestPath))
	documentBytes, err := os.ReadFile(absoluteRequestPath)
	if err != nil {
		return summary, skipped, fmt.Errorf("read request file: %w", err)
	}
	folded := replaceTimingSection(documentBytes, RenderTimingSection(summary))
	if err := atomicfile.ReplaceExisting(absoluteRequestPath, folded); err != nil {
		return summary, skipped, fmt.Errorf("write timing section: %w", err)
	}
	if err := os.Remove(streamPath); err != nil && !os.IsNotExist(err) {
		return summary, skipped, fmt.Errorf("remove raw timing stream: %w", err)
	}
	// A run directory with no streams left is residue; its removal fails
	// harmlessly while any sibling request is still recording.
	_ = os.Remove(filepath.Dir(streamPath))

	summary.StreamState = streamStateFolded
	summary.SectionWritten = true
	return summary, skipped, nil
}

// repositoryRelativeRequestPath confines the one file this command rewrites to
// the repository the caller named. The fold takes its target from argv, so an
// absolute path or one that climbs out of the root would let it rewrite a file
// the repository does not own; every sibling writer in this module states the
// same confinement.
func repositoryRelativeRequestPath(requestPath string) (string, error) {
	trimmed := strings.TrimSpace(requestPath)
	if trimmed == "" {
		return "", errors.New("a request path is required")
	}
	if filepath.IsAbs(trimmed) {
		return "", fmt.Errorf("request path %q must be repository-relative", requestPath)
	}
	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("request path %q escapes or names the repository root", requestPath)
	}
	return filepath.ToSlash(cleaned), nil
}

// summarizeStream reduces the flat stream to the four answers the maintainer
// asked for: how long the run was observed for, where the time went by category,
// which single stage and which single command cost the most, and how much wall
// time no event claimed.
func summarizeStream(records []resultmodel.TimingEventRecord) resultmodel.LifecycleTimingResult {
	summary := resultmodel.LifecycleTimingResult{EventCount: len(records), CategoryTotals: []resultmodel.TimingCategoryTotal{}}
	if len(records) == 0 {
		return summary
	}
	var observedStart, observedEnd time.Time
	totalsByCategory := map[string]resultmodel.TimingCategoryTotal{}
	for index := range records {
		record := records[index]
		if startInstant, err := time.Parse(time.RFC3339, record.StartedAt); err == nil {
			if observedStart.IsZero() || startInstant.Before(observedStart) {
				observedStart = startInstant.UTC()
			}
		}
		if endInstant, err := time.Parse(time.RFC3339, record.EndedAt); err == nil {
			if observedEnd.IsZero() || endInstant.After(observedEnd) {
				observedEnd = endInstant.UTC()
			}
		}
		summary.AttributedSeconds += record.ElapsedSeconds
		total := totalsByCategory[record.Category]
		total.Category = record.Category
		total.ElapsedSeconds += record.ElapsedSeconds
		total.EventCount++
		totalsByCategory[record.Category] = total

		if record.ExitStatus == nil {
			if summary.SlowestStage == nil || record.ElapsedSeconds > summary.SlowestStage.ElapsedSeconds {
				candidate := record
				summary.SlowestStage = &candidate
			}
			continue
		}
		if summary.SlowestCommand == nil || record.ElapsedSeconds > summary.SlowestCommand.ElapsedSeconds {
			candidate := record
			summary.SlowestCommand = &candidate
		}
	}
	if !observedStart.IsZero() && !observedEnd.IsZero() {
		summary.ObservedStart = formatTimingInstant(observedStart)
		summary.ObservedEnd = formatTimingInstant(observedEnd)
		summary.TotalObservedSeconds = int(observedEnd.Sub(observedStart) / time.Second)
	}
	summary.UnattributedSeconds = summary.TotalObservedSeconds - summary.AttributedSeconds
	if summary.UnattributedSeconds < 0 {
		summary.UnattributedSeconds = 0
	}
	for _, total := range totalsByCategory {
		summary.CategoryTotals = append(summary.CategoryTotals, total)
	}
	sort.Slice(summary.CategoryTotals, func(left, right int) bool {
		if summary.CategoryTotals[left].ElapsedSeconds != summary.CategoryTotals[right].ElapsedSeconds {
			return summary.CategoryTotals[left].ElapsedSeconds > summary.CategoryTotals[right].ElapsedSeconds
		}
		return summary.CategoryTotals[left].Category < summary.CategoryTotals[right].Category
	})
	return summary
}

// RenderTimingSection renders the compact Markdown section a fold writes. It
// ends with a newline and carries no heading above `## Timing`.
func RenderTimingSection(summary resultmodel.LifecycleTimingResult) string {
	var section strings.Builder
	section.WriteString(timingSectionHeading + "\n\n")
	fmt.Fprintf(&section, "Observed %s to %s: %s total, %s attributed across %d events, %s unattributed.\n\n",
		summary.ObservedStart, summary.ObservedEnd,
		formatElapsedSeconds(summary.TotalObservedSeconds),
		formatElapsedSeconds(summary.AttributedSeconds),
		summary.EventCount,
		formatElapsedSeconds(summary.UnattributedSeconds))
	section.WriteString("| Category | Elapsed | Events |\n| --- | --- | --- |\n")
	for _, total := range summary.CategoryTotals {
		fmt.Fprintf(&section, "| %s | %s | %d |\n", total.Category, formatElapsedSeconds(total.ElapsedSeconds), total.EventCount)
	}
	section.WriteString("\n")
	if summary.SlowestStage != nil {
		fmt.Fprintf(&section, "Slowest stage: %s / %s, %s, outcome %s.\n",
			summary.SlowestStage.Category, summary.SlowestStage.Operation,
			formatElapsedSeconds(summary.SlowestStage.ElapsedSeconds), summary.SlowestStage.Outcome)
	}
	if summary.SlowestCommand != nil {
		fmt.Fprintf(&section, "Slowest command: %s / %s, %s, exit %d, %s.\n",
			summary.SlowestCommand.Category, summary.SlowestCommand.Operation,
			formatElapsedSeconds(summary.SlowestCommand.ElapsedSeconds),
			*summary.SlowestCommand.ExitStatus, summary.SlowestCommand.CommandIdentity)
	}
	return section.String()
}

// replaceTimingSection puts exactly one `## Timing` section into the document:
// it replaces the request's own section in place and otherwise appends at the
// end. Every byte outside that span is preserved, and the result always ends
// with exactly one newline because the pipeline appends lesson content after
// this writer runs.
func replaceTimingSection(documentBytes []byte, section string) []byte {
	lines := strings.Split(string(documentBytes), "\n")
	fencedLine, unclosedFenceStart := markFencedLines(lines)

	sectionStart := -1
	for index, line := range lines {
		if fencedLine[index] {
			continue
		}
		if strings.TrimRight(line, " \t\r") == timingSectionHeading {
			sectionStart = index
			break
		}
	}
	if sectionStart < 0 {
		document := string(documentBytes)
		document = strings.TrimRight(document, "\n") + "\n\n"
		return []byte(document + section)
	}

	sectionEnd := len(lines)
	for index := sectionStart + 1; index < len(lines); index++ {
		if fencedLine[index] {
			continue
		}
		if strings.HasPrefix(strings.TrimLeft(lines[index], " \t"), "## ") {
			sectionEnd = index
			break
		}
	}
	// An unclosed fence below the heading has no closing line to stop at, so the
	// span would otherwise run to the end of the file and take that fence and
	// everything under it with it. Stop at the opener instead: leaving a stray
	// paragraph behind is recoverable, deleting a caller's content is not.
	if unclosedFenceStart > sectionStart && unclosedFenceStart < sectionEnd {
		sectionEnd = unclosedFenceStart
	}

	rebuilt := append([]string{}, lines[:sectionStart]...)
	rebuilt = append(rebuilt, strings.Split(strings.TrimSuffix(section, "\n"), "\n")...)
	if sectionEnd < len(lines) {
		rebuilt = append(rebuilt, "")
		rebuilt = append(rebuilt, lines[sectionEnd:]...)
	}
	return []byte(strings.TrimRight(strings.Join(rebuilt, "\n"), "\n") + "\n")
}

// markFencedLines reports, for each line, whether it belongs to a fenced region
// a caller wrote, and where a fence that never closes begins (-1 when none does).
// A heading inside such a region is an example, never this writer's own section.
//
// The rule is a condition, not a list of fence spellings. A fence is a run of
// three or more of one ASCII punctuation mark, closed by a run of the same mark
// at least as long, so a dialect fence this package has never heard of still
// hides what it wraps. Only the marks CommonMark itself defines as line
// constructs that enclose nothing are excluded, and an unclosed opener fences
// everything after it, which fails toward appending a duplicate section rather
// than toward deleting a caller's bytes.
func markFencedLines(lines []string) ([]bool, int) {
	fencedLine := make([]bool, len(lines))
	openCharacter, openLength, openIndex := byte(0), 0, -1
	for index, line := range lines {
		runCharacter, runLength := leadingPunctuationRun(line)
		if openLength > 0 {
			fencedLine[index] = true
			if runCharacter == openCharacter && runLength >= openLength {
				openCharacter, openLength, openIndex = 0, 0, -1
			}
			continue
		}
		if runLength >= fenceRunMinimumLength && isEnclosingFenceCharacter(runCharacter) {
			openCharacter, openLength, openIndex = runCharacter, runLength, index
			fencedLine[index] = true
		}
	}
	return fencedLine, openIndex
}

// leadingPunctuationRun returns the ASCII punctuation mark a line opens with and
// how many times it repeats. Leading whitespace is skipped rather than measured,
// so an indented fence still counts as one.
func leadingPunctuationRun(line string) (byte, int) {
	content := strings.TrimLeft(strings.TrimRight(line, " \t\r"), " \t")
	if content == "" || !isMarkdownBlockPunctuation(content[0]) {
		return 0, 0
	}
	runLength := 0
	for runLength < len(content) && content[runLength] == content[0] {
		runLength++
	}
	return content[0], runLength
}

// isMarkdownBlockPunctuation is CommonMark's own ASCII punctuation class, taken
// wholesale rather than narrowed to the marks today's fence syntax happens to
// use. The narrower set would be an enumeration to revisit whenever a dialect
// adds a construct, which is exactly how a fence-blind classifier reopens.
func isMarkdownBlockPunctuation(value byte) bool {
	return value >= '!' && value <= '/' ||
		value >= ':' && value <= '@' ||
		value >= '[' && value <= '`' ||
		value >= '{' && value <= '~'
}

// isEnclosingFenceCharacter excludes the four marks CommonMark defines as line
// constructs that never enclose anything: thematic breaks use "-", "_" and "*",
// and setext underlines use "-" and "=". Treating those as fences would let the
// unpaired "---" a request carries above its source line swallow the rest of the
// document, so the request's own section could never be found again.
func isEnclosingFenceCharacter(value byte) bool {
	if value == 0 || !isMarkdownBlockPunctuation(value) {
		return false
	}
	return value != '-' && value != '_' && value != '*' && value != '='
}

// streamPathFor resolves the Git common directory so every linked worktree of one
// repository appends to the same private stream and nothing lands in the index.
func streamPathFor(repositoryRoot, runIdentifier, requestID string) (string, error) {
	commonDirectoryBytes, err := exec.Command("git", "-C", repositoryRoot, "rev-parse", "--git-common-dir").Output()
	if err != nil {
		return "", fmt.Errorf("resolve Git common directory: %w", err)
	}
	commonDirectory := strings.TrimSpace(string(commonDirectoryBytes))
	if commonDirectory == "" {
		return "", errors.New("resolve Git common directory: git reported no path")
	}
	if !filepath.IsAbs(commonDirectory) {
		commonDirectory = filepath.Join(repositoryRoot, commonDirectory)
	}
	return filepath.Join(commonDirectory, timingDirectoryName, runIdentifier, requestID+".jsonl"), nil
}

// readStream decodes the append-only stream. An undecodable line is counted and
// skipped rather than failing the read: instrumentation must never be the reason
// a finished request cannot be finalized.
func readStream(streamPath string) ([]resultmodel.TimingEventRecord, int, error) {
	contents, err := os.ReadFile(streamPath)
	if os.IsNotExist(err) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, fmt.Errorf("read timing stream: %w", err)
	}
	records := []resultmodel.TimingEventRecord{}
	unreadableLines := 0
	for _, line := range strings.Split(string(contents), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var record resultmodel.TimingEventRecord
		if decodeError := json.Unmarshal([]byte(line), &record); decodeError != nil {
			unreadableLines++
			continue
		}
		records = append(records, record)
	}
	return records, unreadableLines, nil
}

func appendStreamLine(streamPath string, record resultmodel.TimingEventRecord) error {
	if err := os.MkdirAll(filepath.Dir(streamPath), 0o700); err != nil {
		return fmt.Errorf("create timing stream directory: %w", err)
	}
	encoded, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode timing event: %w", err)
	}
	streamFile, err := os.OpenFile(streamPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open timing stream: %w", err)
	}
	defer streamFile.Close()
	if _, err := streamFile.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("append timing event: %w", err)
	}
	return nil
}

// redactCommandIdentity keeps only what a reader needs to recognize a command:
// the executable's base name and how many argv tokens it carried. Arguments can
// carry secrets, paths and user-controlled text, so no argument value ever
// reaches durable evidence.
func redactCommandIdentity(argv []string) string {
	if len(argv) == 0 {
		return ""
	}
	executable := boundEvidenceText(filepath.Base(argv[0]))
	if executable == "" {
		executable = "unnamed-executable"
	}
	return fmt.Sprintf("%s (%d argv tokens)", executable, len(argv))
}

// boundEvidenceText keeps free text printable, single-line and short, so an
// operation name or agent label cannot smuggle control characters or an
// unbounded payload into the stream or the folded section.
func boundEvidenceText(value string) string {
	cleaned := strings.Map(func(character rune) rune {
		if character == '\n' || character == '\r' || character == '\t' {
			return ' '
		}
		if unicode.IsControl(character) || character == '|' {
			return -1
		}
		return character
	}, value)
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	runes := []rune(cleaned)
	if len(runes) > evidenceTextLimit {
		return string(runes[:evidenceTextLimit])
	}
	return cleaned
}

func formatTimingInstant(instant time.Time) string {
	return instant.UTC().Format(time.RFC3339)
}

// formatElapsedSeconds renders a duration the same way every time, so a folded
// section is byte-stable for one stream.
func formatElapsedSeconds(totalSeconds int) string {
	if totalSeconds < 0 {
		totalSeconds = 0
	}
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	switch {
	case hours > 0:
		return fmt.Sprintf("%dh %02dm %02ds", hours, minutes, seconds)
	case minutes > 0:
		return fmt.Sprintf("%dm %02ds", minutes, seconds)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}

func vocabularyContains(vocabulary []string, value string) bool {
	for _, candidate := range vocabulary {
		if candidate == value {
			return true
		}
	}
	return false
}
