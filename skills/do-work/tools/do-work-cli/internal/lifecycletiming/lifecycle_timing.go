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

// RunTimedCommand executes argv, measures it in-process, and records one command
// event. It returns the child's exit status alongside the typed projection. This
// is the only path where one process observes both ends of an event, so it is
// the only path that reports a monotonic elapsed source.
func RunTimedCommand(repositoryRoot string, eventRequest EventRequest, argv []string, childOutput io.Writer) (resultmodel.LifecycleTimingResult, int, error) {
	if len(argv) == 0 {
		return resultmodel.LifecycleTimingResult{}, 0, errors.New("a timed command requires argv after --")
	}
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
		return resultmodel.LifecycleTimingResult{}, 0, fmt.Errorf("launch %s: %w", filepath.Base(argv[0]), runError)
	}
	exitStatus := command.ProcessState.ExitCode()
	if exitStatus < 0 {
		exitStatus = 128
	}

	eventRequest.CommandArgv = argv
	eventRequest.ExitStatus = &exitStatus
	eventRequest.Outcome = "success"
	if exitStatus != 0 {
		eventRequest.Outcome = "failure"
	}
	if eventRequest.Operation == "" {
		eventRequest.Operation = redactCommandIdentity(argv)
	}
	// The start is derived from the measured elapsed rather than from the
	// previous event's end, so wall time spent before the launch stays
	// unattributed instead of being silently folded into the command.
	measuredElapsed := endInstant.Sub(startInstant)
	eventRequest.StartedAt = endInstant.UTC().Add(-measuredElapsed)
	eventRequest.HasExplicitStart = true

	elapsedSeconds := int(measuredElapsed.Round(time.Second) / time.Second)
	result, err := appendMeasuredEvent(repositoryRoot, eventRequest, endInstant.UTC(), elapsedSourceMonotonic, &elapsedSeconds)
	return result, exitStatus, err
}

// appendMeasuredEvent is the single write path. Every field is validated and
// bounded before one JSON line is appended to the stream.
func appendMeasuredEvent(repositoryRoot string, eventRequest EventRequest, endInstant time.Time, elapsedSource string, measuredSeconds *int) (resultmodel.LifecycleTimingResult, error) {
	runIdentifier := eventRequest.RunIdentifier
	if runIdentifier == "" {
		runIdentifier = "standalone"
	}
	if !timingRunIDPattern.MatchString(runIdentifier) {
		return resultmodel.LifecycleTimingResult{}, fmt.Errorf("run identity %q is not a single safe path segment", runIdentifier)
	}
	if !timingRequestIDPattern.MatchString(eventRequest.RequestID) {
		return resultmodel.LifecycleTimingResult{}, fmt.Errorf("request identity %q is not REQ-NNN", eventRequest.RequestID)
	}
	if !vocabularyContains(TimingCategoryVocabulary, eventRequest.Category) {
		return resultmodel.LifecycleTimingResult{}, fmt.Errorf("category %q is not one of %s", eventRequest.Category, strings.Join(TimingCategoryVocabulary, ", "))
	}
	operation := boundEvidenceText(eventRequest.Operation)
	if operation == "" {
		return resultmodel.LifecycleTimingResult{}, errors.New("an operation name is required")
	}
	outcome := eventRequest.Outcome
	if outcome == "" {
		outcome = "success"
	}
	if !vocabularyContains(TimingOutcomeVocabulary, outcome) {
		return resultmodel.LifecycleTimingResult{}, fmt.Errorf("outcome %q is not one of %s", outcome, strings.Join(TimingOutcomeVocabulary, ", "))
	}

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
	if strings.TrimSpace(foldRequest.RequestPath) == "" {
		return resultmodel.LifecycleTimingResult{}, nil, errors.New("a request path is required")
	}
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
// it replaces an existing one in place and otherwise appends at the end. Every
// byte outside that span is preserved.
func replaceTimingSection(documentBytes []byte, section string) []byte {
	lines := strings.Split(string(documentBytes), "\n")
	sectionStart := -1
	for index, line := range lines {
		if strings.TrimRight(line, " \t\r") == timingSectionHeading {
			sectionStart = index
			break
		}
	}
	if sectionStart < 0 {
		document := string(documentBytes)
		if !strings.HasSuffix(document, "\n") {
			document += "\n"
		}
		if !strings.HasSuffix(document, "\n\n") {
			document += "\n"
		}
		return []byte(document + section)
	}
	sectionEnd := len(lines)
	for index := sectionStart + 1; index < len(lines); index++ {
		if strings.HasPrefix(lines[index], "## ") {
			sectionEnd = index
			break
		}
	}
	replacement := strings.Split(strings.TrimSuffix(section, "\n"), "\n")
	if sectionEnd < len(lines) {
		replacement = append(replacement, "")
	}
	rebuilt := append([]string{}, lines[:sectionStart]...)
	rebuilt = append(rebuilt, replacement...)
	rebuilt = append(rebuilt, lines[sectionEnd:]...)
	return []byte(strings.Join(rebuilt, "\n"))
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
