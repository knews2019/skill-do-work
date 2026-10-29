package toolboxcommands

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/rootedfs"
)

const CommandPeriodicCheck = "periodic-check-history"
const periodicHistoryDirectory = "do-work/deliverables/periodic-check"

var periodicIdentifierPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
var periodicRequestPattern = regexp.MustCompile(`^REQ-[0-9]{3,}$`)

type periodicChecklist struct {
	CompatibilityVersion string                    `json:"version"`
	CaptureEnabled       bool                      `json:"capture"`
	CheckRecords         []periodicCheckDefinition `json:"checks"`
}
type periodicCheckDefinition struct {
	StableID             string `json:"id"`
	CompatibilityVersion string `json:"version"`
	CadenceSeconds       int64  `json:"cadence_seconds"`
	QuickProfile         bool   `json:"quick"`
	CheckCommand         string `json:"command"`
	PrerequisiteText     string `json:"prerequisites"`
	ResultInterpretation string `json:"interpretation"`
}
type periodicFinding struct {
	StableID        string `json:"id"`
	FindingSummary  string `json:"summary"`
	EvidenceText    string `json:"evidence"`
	RepairRequest   string `json:"repair_request,omitempty"`
	FindingDelta    string `json:"delta,omitempty"`
	CaptureAction   string `json:"capture_action,omitempty"`
	PreviousRequest string `json:"previous_request,omitempty"`
	CaptureNote     string `json:"capture_note,omitempty"`
}
type periodicCheckResult struct {
	CheckID        string            `json:"check_id"`
	CheckCompleted bool              `json:"completed"`
	EvidenceText   string            `json:"evidence"`
	CheckFindings  []periodicFinding `json:"findings"`
	RecordedAt     time.Time         `json:"recorded_at"`
}
type periodicRun struct {
	SchemaVersion      int                   `json:"schema_version"`
	RecordPath         string                `json:"path"`
	StartedAt          time.Time             `json:"started_at"`
	FinishedAt         *time.Time            `json:"finished_at,omitempty"`
	SourceRevision     string                `json:"revision"`
	DirtyTree          string                `json:"dirty_tree"`
	RunProfile         string                `json:"profile"`
	ChecklistConfig    periodicChecklist     `json:"checklist"`
	CaptureEnabled     bool                  `json:"capture"`
	UnattendedMode     bool                  `json:"unattended"`
	ConfigurationError string                `json:"configuration_error,omitempty"`
	RunOutcome         string                `json:"outcome"`
	CheckRecords       []periodicCheckResult `json:"checks"`
	MissingCoverage    []string              `json:"missing_coverage"`
	ReportPath         string                `json:"report"`
}

func periodicDecode(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("expected one JSON value")
	}
	return nil
}

func periodicLoadChecklist(root string) (periodicChecklist, error) {
	var checklist periodicChecklist
	data, err := rootedReadFile(root, "prime-periodic-check.md")
	if err != nil {
		return checklist, fmt.Errorf("read prime-periodic-check.md: %w", err)
	}
	const marker = "```periodic-check\n"
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if strings.Count(text, marker) != 1 {
		return checklist, fmt.Errorf("prime needs exactly one periodic-check fenced JSON block")
	}
	body := strings.SplitN(text, marker, 2)[1]
	end := strings.Index(body, "\n```")
	if end < 0 {
		return checklist, fmt.Errorf("unclosed periodic-check block")
	}
	if err := periodicDecode([]byte(body[:end]), &checklist); err != nil {
		return checklist, err
	}
	if checklist.CompatibilityVersion == "" || len(checklist.CheckRecords) == 0 {
		return checklist, fmt.Errorf("checklist needs version and checks")
	}
	seen := map[string]bool{}
	for _, check := range checklist.CheckRecords {
		if !periodicIdentifierPattern.MatchString(check.StableID) || seen[check.StableID] || check.CompatibilityVersion == "" || check.CadenceSeconds <= 0 || check.CadenceSeconds > int64((3650*24*time.Hour)/time.Second) || strings.TrimSpace(check.CheckCommand) == "" || strings.TrimSpace(check.ResultInterpretation) == "" {
			return checklist, fmt.Errorf("invalid or duplicate check %q: require id, version, positive cadence_seconds (up to ten years), command and interpretation", check.StableID)
		}
		seen[check.StableID] = true
	}
	return checklist, nil
}

func periodicSelected(checklist periodicChecklist, profile string) []periodicCheckDefinition {
	var selected []periodicCheckDefinition
	for _, check := range checklist.CheckRecords {
		if profile == "full" || check.QuickProfile {
			selected = append(selected, check)
		}
	}
	return selected
}

func periodicReadRun(root, path string) (periodicRun, error) {
	var run periodicRun
	if filepath.ToSlash(filepath.Dir(path)) != periodicHistoryDirectory || !periodicIdentifierPattern.MatchString(filepath.Base(path)) {
		return run, fmt.Errorf("invalid run path")
	}
	data, err := rootedReadFile(root, path+"/record.json")
	if err != nil {
		return run, err
	}
	if err = periodicDecode(data, &run); err != nil {
		return run, err
	}
	if run.SchemaVersion != 1 || run.RecordPath != path || run.ReportPath != path+"/report.md" || run.StartedAt.IsZero() || (run.RunProfile != "full" && run.RunProfile != "quick") {
		return run, fmt.Errorf("invalid run record: %s", path)
	}
	return run, nil
}

func periodicReadHistory(root string) ([]periodicRun, error) {
	if err := validateNoLinkedAncestors(root, periodicHistoryDirectory, true); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(root, periodicHistoryDirectory))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var history []periodicRun
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".pending-") {
			continue
		}
		if !entry.IsDir() {
			return nil, fmt.Errorf("unexpected history entry: %s", entry.Name())
		}
		run, err := periodicReadRun(root, periodicHistoryDirectory+"/"+entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read history: %w", err)
		}
		history = append(history, run)
	}
	sort.Slice(history, func(i, j int) bool {
		if history[i].StartedAt.Equal(history[j].StartedAt) {
			return history[i].RecordPath > history[j].RecordPath
		}
		return history[i].StartedAt.After(history[j].StartedAt)
	})
	return history, nil
}

func periodicCompatible(left, right periodicRun, checkID string) bool {
	if left.RunProfile != right.RunProfile || left.ChecklistConfig.CompatibilityVersion != right.ChecklistConfig.CompatibilityVersion {
		return false
	}
	version := ""
	for _, check := range left.ChecklistConfig.CheckRecords {
		if check.StableID == checkID {
			version = check.CompatibilityVersion
		}
	}
	for _, check := range right.ChecklistConfig.CheckRecords {
		if check.StableID == checkID {
			return version != "" && version == check.CompatibilityVersion
		}
	}
	return false
}

func periodicPrevious(history []periodicRun, current periodicRun, checkID string) *periodicCheckResult {
	var latest *periodicCheckResult
	for _, prior := range history {
		if prior.RecordPath == current.RecordPath || !periodicCompatible(current, prior, checkID) {
			continue
		}
		for _, result := range prior.CheckRecords {
			if result.CheckID == checkID && result.CheckCompleted && (latest == nil || result.RecordedAt.After(latest.RecordedAt)) {
				copy := result
				latest = &copy
			}
		}
	}
	return latest
}

// Finding identity is supplied by the project interpreter; timestamps and volatile
// evidence are deliberately excluded from change detection.
func periodicCompare(history []periodicRun, current periodicRun, result *periodicCheckResult) {
	previous := periodicPrevious(history, current, result.CheckID)
	for index := range result.CheckFindings {
		finding := &result.CheckFindings[index]
		finding.FindingDelta, finding.CaptureAction = "new", "create"
		if !result.CheckCompleted {
			finding.FindingDelta, finding.CaptureAction = "unassessed", "defer"
			continue
		}
		if previous != nil {
			for _, old := range previous.CheckFindings {
				if old.StableID != finding.StableID || old.FindingDelta == "resolved" {
					continue
				}
				finding.FindingDelta = "unchanged"
				if old.FindingSummary != finding.FindingSummary {
					finding.FindingDelta = "changed"
				}
				finding.RepairRequest = old.RepairRequest
				if old.RepairRequest != "" {
					finding.CaptureAction = "inspect-request"
				}
			}
		}
		if finding.FindingDelta == "new" && previous != nil {
			// A successful equivalent check without the finding is resolution evidence.
			var linkedAt time.Time
			for _, prior := range history {
				if prior.RecordPath == current.RecordPath || !periodicCompatible(current, prior, result.CheckID) {
					continue
				}
				for _, check := range prior.CheckRecords {
					if check.CheckID != result.CheckID || !check.CheckCompleted || check.RecordedAt.After(previous.RecordedAt) {
						continue
					}
					for _, old := range check.CheckFindings {
						if old.StableID == finding.StableID && old.RepairRequest != "" && check.RecordedAt.After(linkedAt) {
							linkedAt = check.RecordedAt
							finding.CaptureAction, finding.PreviousRequest = "recurrence", old.RepairRequest
						}
					}
				}
			}
		}
	}
	if result.CheckCompleted && previous != nil {
		for _, old := range previous.CheckFindings {
			if old.FindingDelta == "resolved" {
				continue
			}
			found := false
			for _, finding := range result.CheckFindings {
				if finding.StableID == old.StableID {
					found = true
				}
			}
			if !found {
				old.FindingDelta, old.CaptureAction = "resolved", "none"
				result.CheckFindings = append(result.CheckFindings, old)
			}
		}
	}
}

func periodicRefresh(run *periodicRun) {
	run.RunOutcome = "incomplete"
	run.MissingCoverage = nil
	findings := false
	for _, definition := range periodicSelected(run.ChecklistConfig, run.RunProfile) {
		covered := false
		for _, check := range run.CheckRecords {
			if check.CheckID == definition.StableID {
				covered = check.CheckCompleted
				for _, finding := range check.CheckFindings {
					if finding.FindingDelta != "resolved" {
						findings = true
					}
				}
			}
		}
		if !covered {
			run.MissingCoverage = append(run.MissingCoverage, definition.StableID)
		}
	}
	if run.FinishedAt != nil && run.ConfigurationError == "" && len(run.MissingCoverage) == 0 {
		run.RunOutcome = "clean"
		if findings {
			run.RunOutcome = "completed-with-findings"
		}
	}
}

func periodicSaveRun(root string, run periodicRun, replace bool) error {
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	return rootedPublishFile(root, run.RecordPath+"/record.json", append(data, '\n'), 0o644, replace)
}

func periodicTimestamp(timestamp time.Time, now time.Time) string {
	age := now.Sub(timestamp).Round(time.Second)
	if age < 0 {
		return timestamp.UTC().Format(time.RFC3339) + " (clock is ahead of this machine)"
	}
	return fmt.Sprintf("%s (%s ago)", timestamp.UTC().Format(time.RFC3339), age)
}

func periodicStatus(root string, history []periodicRun, checklist periodicChecklist, configErr error) string {
	var output strings.Builder
	now := toolboxNow().UTC()
	for _, profile := range []string{"full", "quick"} {
		fmt.Fprintf(&output, "%s\n", profile)
		for _, kind := range []string{"attempt", "completed", "clean"} {
			label := "no recorded run"
			var latest time.Time
			for _, run := range history {
				if run.RunProfile != profile {
					continue
				}
				if kind != "attempt" && (run.FinishedAt == nil || run.RunOutcome == "incomplete") {
					continue
				}
				if kind == "clean" && run.RunOutcome != "clean" {
					continue
				}
				stamp := run.StartedAt
				if kind != "attempt" {
					stamp = *run.FinishedAt
				}
				if !stamp.After(latest) {
					continue
				}
				latest = stamp
				label = periodicTimestamp(stamp, now) + " — " + run.RunOutcome + " — " + run.RecordPath
				if run.ChecklistConfig.CompatibilityVersion != checklist.CompatibilityVersion {
					label += " (historical checklist)"
				}
			}
			fmt.Fprintf(&output, "  last %s: %s\n", kind, label)
		}
		current := periodicRun{RunProfile: profile, ChecklistConfig: checklist}
		for _, definition := range periodicSelected(checklist, profile) {
			previous := periodicPrevious(history, current, definition.StableID)
			due := "due: no equivalent completed check"
			if previous != nil {
				deadline := previous.RecordedAt.Add(time.Duration(definition.CadenceSeconds) * time.Second)
				due = "due " + deadline.UTC().Format(time.RFC3339)
				if !now.Before(deadline) {
					due += " (overdue)"
				}
				for _, finding := range previous.CheckFindings {
					if finding.FindingDelta != "resolved" {
						fmt.Fprintf(&output, "  unresolved %s/%s: %s\n", definition.StableID, finding.StableID, finding.FindingSummary)
					}
				}
			}
			fmt.Fprintf(&output, "  %s: %s\n", definition.StableID, due)
		}
		if profile == "quick" && len(periodicSelected(checklist, profile)) == 0 {
			output.WriteString("  quick subset is undefined\n")
		}
	}
	if configErr != nil {
		fmt.Fprintf(&output, "Configuration unavailable: %v\n", configErr)
	}
	output.WriteString("History covers records available in this checkout; unshared runs on other computers are unknown.\n")
	return output.String()
}

func periodicReport(run periodicRun) string {
	var output strings.Builder
	fmt.Fprintf(&output, "# Periodic check: %s\n\nProfile: %s\nStarted: %s\nRevision: %s\nReport: %s\n", run.RunOutcome, run.RunProfile, run.StartedAt.Format(time.RFC3339), run.SourceRevision, run.ReportPath)
	if run.ConfigurationError != "" {
		fmt.Fprintf(&output, "\nConfiguration: %s\n", run.ConfigurationError)
	}
	for _, check := range run.CheckRecords {
		fmt.Fprintf(&output, "\n## %s\n\nCompleted: %t\nEvidence: %s\n", check.CheckID, check.CheckCompleted, check.EvidenceText)
		for _, finding := range check.CheckFindings {
			fmt.Fprintf(&output, "\n- %s — %s: %s\n  Evidence: %s\n  Repair: %s; previous: %s; capture: %s %s\n", finding.StableID, finding.FindingDelta, finding.FindingSummary, finding.EvidenceText, finding.RepairRequest, finding.PreviousRequest, finding.CaptureAction, finding.CaptureNote)
		}
	}
	fmt.Fprintf(&output, "\nMissing coverage: %s\n", strings.Join(run.MissingCoverage, ", "))
	return output.String()
}

func handlePeriodicCheck(ctx commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	output, err := periodicCommand(ctx.RepositoryRoot, arguments)
	if err != nil {
		return usageResult(CommandPeriodicCheck, err.Error())
	}
	return exactOutputResult(output, nil)
}

func periodicCommand(root string, arguments []string) (string, error) {
	if len(arguments) == 0 {
		return "", fmt.Errorf("usage: periodic-check-history status|start|record|link|finish [--profile full|quick] [--run relative-path] [--input JSON-file] [--capture|--report-only] [--unattended]")
	}
	operation, profile, runPath, input := arguments[0], "full", "", ""
	capture, reportOnly, unattended := false, false, false
	for index := 1; index < len(arguments); index++ {
		switch arguments[index] {
		case "--profile", "--run", "--input":
			name := arguments[index]
			value, err := optionValue(arguments, &index, name)
			if err != nil {
				return "", err
			}
			switch name {
			case "--profile":
				profile = value
			case "--run":
				runPath = value
			case "--input":
				input = value
			}
		case "--capture":
			capture = true
		case "--report-only":
			reportOnly = true
		case "--unattended":
			unattended = true
		default:
			return "", fmt.Errorf("unknown option %s", arguments[index])
		}
	}
	if profile != "full" && profile != "quick" {
		return "", fmt.Errorf("profile must be full or quick")
	}
	history, err := periodicReadHistory(root)
	if err != nil {
		return "", err
	}
	if operation == "status" || operation == "start" {
		checklist, configErr := periodicLoadChecklist(root)
		output := periodicStatus(root, history, checklist, configErr)
		if operation == "status" {
			return output, nil
		}
		if profile == "quick" && configErr == nil && len(periodicSelected(checklist, profile)) == 0 {
			configErr = fmt.Errorf("quick subset is undefined")
		}
		revision, err := auditGit(root, "rev-parse", "HEAD")
		if err != nil {
			return "", err
		}
		dirty, err := auditGit(root, "status", "--porcelain=v1", "--untracked-files=all")
		if err != nil {
			return "", err
		}
		suffix := make([]byte, 12)
		if _, err := rand.Read(suffix); err != nil {
			return "", err
		}
		now := toolboxNow().UTC()
		path := periodicHistoryDirectory + "/" + now.Format("20060102-150405") + "-" + hex.EncodeToString(suffix)
		run := periodicRun{SchemaVersion: 1, RecordPath: path, StartedAt: now, SourceRevision: strings.TrimSpace(revision), DirtyTree: dirty, RunProfile: profile, ChecklistConfig: checklist, CaptureEnabled: (capture || checklist.CaptureEnabled) && !reportOnly, UnattendedMode: unattended, ReportPath: path + "/report.md"}
		if configErr != nil {
			run.ConfigurationError = configErr.Error()
			run.FinishedAt = &now
		}
		periodicRefresh(&run)
		if err := rootedMkdirAll(root, periodicHistoryDirectory, 0o755); err != nil {
			return "", err
		}
		if err := periodicPublishInitial(root, run); err != nil {
			return "", err
		}
		return output + "Run: " + path + "\n" + periodicReport(run), nil
	}
	if operation != "record" && operation != "link" && operation != "finish" {
		return "", fmt.Errorf("unknown operation %s", operation)
	}
	run, err := periodicReadRun(root, runPath)
	if err != nil {
		return "", err
	}
	// Exclusive writer claim prevents a late scan writer from overwriting a finish.
	lockPath := runPath + "/writer-lock"
	if _, err := rootedMkdirExclusive(root, lockPath, 0o700); err != nil {
		return "", fmt.Errorf("run writer busy or interrupted; preserve this attempt and start a new one: %w", err)
	}
	defer os.Remove(filepath.Join(root, lockPath))
	run, err = periodicReadRun(root, runPath)
	if err != nil {
		return "", err
	}
	if run.FinishedAt != nil {
		return "", fmt.Errorf("completed attempts are immutable")
	}
	switch operation {
	case "record":
		data, err := os.ReadFile(input)
		if err != nil {
			return "", err
		}
		var check periodicCheckResult
		if err := periodicDecode(data, &check); err != nil {
			return "", err
		}
		selected := false
		for _, definition := range periodicSelected(run.ChecklistConfig, run.RunProfile) {
			if definition.StableID == check.CheckID {
				selected = true
			}
		}
		if !selected || strings.TrimSpace(check.EvidenceText) == "" {
			return "", fmt.Errorf("check must belong to selected profile and include evidence or missing-coverage reason")
		}
		for _, previous := range run.CheckRecords {
			if previous.CheckID == check.CheckID {
				return "", fmt.Errorf("check already recorded; use a new attempt to repeat it")
			}
		}
		seen := map[string]bool{}
		for _, finding := range check.CheckFindings {
			if finding.StableID == "" || seen[finding.StableID] || finding.FindingSummary == "" || finding.EvidenceText == "" || finding.RepairRequest != "" || finding.FindingDelta != "" || finding.CaptureAction != "" || finding.PreviousRequest != "" || finding.CaptureNote != "" {
				return "", fmt.Errorf("findings require unique ids, summary and evidence; derived fields must be omitted")
			}
			seen[finding.StableID] = true
		}
		check.RecordedAt = toolboxNow().UTC()
		periodicCompare(history, run, &check)
		periodicResolveLinks(root, &check)
		run.CheckRecords = append(run.CheckRecords, check)
	case "link":
		data, err := os.ReadFile(input)
		if err != nil {
			return "", err
		}
		var link struct {
			CheckID            string `json:"check_id"`
			FindingID          string `json:"finding_id"`
			RequestID          string `json:"request_id"`
			CaptureDisposition string `json:"disposition"`
			CaptureNote        string `json:"note"`
		}
		if err := periodicDecode(data, &link); err != nil {
			return "", err
		}
		if link.CaptureNote == "" || (link.RequestID != "" && !periodicRequestPattern.MatchString(link.RequestID)) {
			return "", fmt.Errorf("link requires note and a valid request id when supplied")
		}
		switch link.CaptureDisposition {
		case "linked", "created", "verify-resolution", "deferred", "report-only":
		default:
			return "", fmt.Errorf("invalid capture disposition")
		}
		if (link.CaptureDisposition == "linked" || link.CaptureDisposition == "created" || link.CaptureDisposition == "verify-resolution") && link.RequestID == "" {
			return "", fmt.Errorf("capture disposition requires request id")
		}
		if link.CaptureDisposition == "created" && !run.CaptureEnabled {
			return "", fmt.Errorf("capture was not authorized for this attempt")
		}
		if link.RequestID != "" {
			if _, err := periodicRequestStatus(root, link.RequestID); err != nil {
				return "", err
			}
		}
		found := false
		for ci := range run.CheckRecords {
			for fi := range run.CheckRecords[ci].CheckFindings {
				finding := &run.CheckRecords[ci].CheckFindings[fi]
				if run.CheckRecords[ci].CheckID == link.CheckID && finding.StableID == link.FindingID && finding.FindingDelta != "resolved" {
					if link.CaptureDisposition == "created" && !run.CheckRecords[ci].CheckCompleted {
						return "", fmt.Errorf("incomplete check cannot authorize new capture")
					}
					if finding.CaptureNote != "" {
						return "", fmt.Errorf("capture disposition already recorded")
					}
					finding.RepairRequest, finding.CaptureAction, finding.CaptureNote = link.RequestID, link.CaptureDisposition, link.CaptureNote
					found = true
				}
			}
		}
		if !found {
			return "", fmt.Errorf("active finding not found")
		}
	case "finish":
		now := toolboxNow().UTC()
		for ci := range run.CheckRecords {
			for fi := range run.CheckRecords[ci].CheckFindings {
				finding := &run.CheckRecords[ci].CheckFindings[fi]
				if finding.FindingDelta != "resolved" && finding.CaptureNote == "" {
					if run.CaptureEnabled {
						finding.CaptureAction = "deferred"
						finding.CaptureNote = "Capture not completed; inspect existing requests before retrying."
					} else {
						finding.CaptureAction = "report-only"
						finding.CaptureNote = "Capture disabled."
					}
				}
			}
		}
		run.FinishedAt = &now
	}
	periodicRefresh(&run)
	if operation == "finish" {
		// A previous interrupted finish may already have published the identical report.
		report := []byte(periodicReport(run))
		old, readErr := rootedReadFile(root, run.ReportPath)
		if readErr == nil && !bytes.Equal(old, report) {
			return "", fmt.Errorf("report already exists with different content")
		}
		if os.IsNotExist(readErr) {
			if err := rootedPublishFile(root, run.ReportPath, report, 0o644, false); err != nil {
				return "", err
			}
		} else if readErr != nil {
			return "", readErr
		}
	}
	if err := periodicSaveRun(root, run, true); err != nil {
		return "", err
	}
	return periodicReport(run), nil
}

// Resolve known links across queue, working, and archive without changing any REQ.
// Root-cause matches without an existing link remain the agent's capture judgment.
func periodicResolveLinks(root string, result *periodicCheckResult) {
	for index := range result.CheckFindings {
		finding := &result.CheckFindings[index]
		if finding.RepairRequest == "" || finding.FindingDelta == "resolved" {
			continue
		}
		status, err := periodicRequestStatus(root, finding.RepairRequest)
		finding.CaptureAction = "defer"
		if err != nil {
			continue
		}
		switch status {
		case "pending", "pending-answers", "in-progress", "working", "blocked":
			finding.CaptureAction = "link"
		case "completed", "completed-with-issues":
			finding.CaptureAction = "verify-resolution"
		}
	}
}

func periodicRequestStatus(root, requestID string) (string, error) {
	var statuses []string
	for _, directory := range []string{"do-work/queue", "do-work/working", "do-work/archive"} {
		if err := validateNoLinkedAncestors(root, directory, true); err != nil {
			return "", err
		}
		err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, walkErr error) error {
			if os.IsNotExist(walkErr) {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("request lookup contains a symbolic link")
			}
			if entry.IsDir() || !strings.HasPrefix(entry.Name(), "REQ-") || !strings.HasSuffix(entry.Name(), ".md") {
				return nil
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			data, err := rootedReadFile(root, relative)
			if err != nil {
				return err
			}
			document, err := requestmodel.ParseDocument(data)
			if err != nil {
				return err
			}
			record := document.TypedRecord()
			if record.RequestID == requestID {
				statuses = append(statuses, record.RequestStatus)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}
	if len(statuses) != 1 {
		return "", fmt.Errorf("repair request is absent or ambiguous: %s", requestID)
	}
	return statuses[0], nil
}

// Publish a complete initial record before history readers can discover it.
// A killed pre-publication start leaves only private staging, never false coverage.
func periodicPublishInitial(repositoryRoot string, run periodicRun) error {
	staging := periodicHistoryDirectory + "/.pending-" + filepath.Base(run.RecordPath)
	if _, err := rootedMkdirExclusive(repositoryRoot, staging, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return err
	}
	if err := rootedPublishFile(repositoryRoot, staging+"/record.json", append(data, '\n'), 0o644, false); err != nil {
		return err
	}
	if run.FinishedAt != nil {
		if err := rootedPublishFile(repositoryRoot, staging+"/report.md", []byte(periodicReport(run)), 0o644, false); err != nil {
			return err
		}
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	if _, err := root.Lstat(run.RecordPath); !os.IsNotExist(err) {
		return fmt.Errorf("run destination already exists or cannot be inspected")
	}
	return rootedfs.Rename(root, staging, run.RecordPath)
}
