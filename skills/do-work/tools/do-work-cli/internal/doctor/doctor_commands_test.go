package doctor

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDoctorCommandIsRegisteredAndDefaultRunIsReadOnly(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("go-run integration is heavy-only")
	}
	repositoryRoot := t.TempDir()
	requestPath := filepath.Join(repositoryRoot, "do-work", "queue", "REQ-001-clean.md")
	if err := os.MkdirAll(filepath.Dir(requestPath), 0o755); err != nil {
		t.Fatal(err)
	}
	original := []byte("---\nid: REQ-001\nstatus: pending\ncreated_at: 2026-08-30T20:00:00Z\n---\nBody\n")
	if err := os.WriteFile(requestPath, original, 0o644); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "doctor")
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("doctor returned %v:\n%s", runError, output)
	}
	after, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, original) {
		t.Fatal("default doctor changed a request")
	}
	if strings.Contains(string(output), "UNKNOWN-COMMAND") {
		t.Fatalf("doctor is not registered:\n%s", output)
	}
}

func TestCommittedRepairRediscoveryFailureReturnsRiskAndExactRevert(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-050-time.md", doctorRequest("REQ-050", "pending", "created_at: 2099-01-01T00:00:00Z\n", "Body"))
	commitDoctorFixture(t, repositoryRoot, "fixture")
	originalDiscover := discoverRepository
	discoveryCount := 0
	discoverRepository = func(root string) (*repositorymodel.RepositorySnapshot, error) {
		discoveryCount++
		if discoveryCount == 2 {
			return nil, errors.New("injected post-commit rediscovery failure")
		}
		return originalDiscover(root)
	}
	t.Cleanup(func() { discoverRepository = originalDiscover })

	result := handleDoctor(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--repair-timestamps", "--commit"})
	if result.Outcome != resultmodel.OutcomeRisk || resultmodel.ExitCode(result.Outcome) != 4 {
		t.Fatalf("post-commit rediscovery result = %#v", result)
	}
	headCommand := exec.Command("git", "-C", repositoryRoot, "rev-parse", "HEAD")
	headBytes, err := headCommand.Output()
	if err != nil {
		t.Fatal(err)
	}
	head := strings.TrimSpace(string(headBytes))
	found := false
	for _, finding := range result.Findings {
		if finding.Code == "DOCTOR-POST-REPAIR-DISCOVERY-FAILED" {
			found = true
			want := []string{"git", "revert", head}
			if strings.Join(finding.NextArgv, "\x00") != strings.Join(want, "\x00") {
				t.Fatalf("revert argv = %v, want %v", finding.NextArgv, want)
			}
		}
	}
	if !found {
		t.Fatalf("missing post-repair discovery finding: %#v", result.Findings)
	}
}

func TestForensicsActionDelegatesMechanicalChecksOnlyToDoctor(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeDoctorFixture(t, repositoryRoot, "do-work/working/REQ-060-stuck.md", doctorRequest("REQ-060", "in-progress", "created_at: 2026-08-20T00:00:00Z\nclaimed_at: 2026-08-20T01:00:00Z\n", "## Plan\nInterrupted"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-061-queued.md", doctorRequest("REQ-061", "pending", "created_at: 2026-08-29T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-062-terminal.md", doctorRequest("REQ-062", "completed", "created_at: 2026-08-19T00:00:00Z\ncompleted_at: 2026-08-21T00:00:00Z\n", "## Implementation Summary\n**Files changed:** none (no code changes needed)\n\n## Qualification\nVerified"))
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(t.Context(), snapshot, ScanOptions{Now: time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)})
	var stuckWork *resultmodel.CommandFinding
	for index := range result.Findings {
		if result.Findings[index].Code == "STUCK-WORK" {
			stuckWork = &result.Findings[index]
			break
		}
	}
	if stuckWork == nil {
		t.Fatalf("mixed-state doctor result omitted STUCK-WORK: %#v", result.Findings)
	}
	if strings.Join(stuckWork.AffectedIDs, "\n") != "REQ-060" || strings.Join(stuckWork.AffectedPaths, "\n") != "do-work/working/REQ-060-stuck.md" {
		t.Fatalf("STUCK-WORK lost exact identity: %#v", stuckWork)
	}
	if len(stuckWork.Evidence) == 0 || stuckWork.Fixability != resultmodel.FixabilityManual || len(stuckWork.NextArgv) == 0 || len(stuckWork.VerificationArgv) == 0 {
		t.Fatalf("STUCK-WORK is not report-ready: %#v", stuckWork)
	}

	actionPath := filepath.Join("..", "..", "..", "..", "actions", "forensics.md")
	actionBytes, err := os.ReadFile(actionPath)
	if err != nil {
		t.Fatal(err)
	}
	action := string(actionBytes)
	if !strings.Contains(action, "tools/do-work-cli.sh --repo-root <project-root> doctor") {
		t.Fatal("forensics action does not launch the canonical doctor")
	}
	for _, legacyMechanic := range []string{"repair-req-timestamps.sh", "audit-archive-timestamps.sh"} {
		if strings.Contains(action, legacyMechanic) {
			t.Fatalf("forensics action still delegates legacy mechanic %q", legacyMechanic)
		}
	}
	if strings.Count(action, "blanked-req-scan.sh") != 1 || !strings.Contains(action, "forensics must not execute it") {
		t.Fatal("blank scanner compatibility pointer must remain explicitly non-executable")
	}
	for _, fieldName := range []string{"code", "severity", "affected_ids", "affected_paths", "observed_evidence", "fixability", "automation_stop_reason", "next_argv", "verification_argv"} {
		if !strings.Contains(action, "`"+fieldName+"`") {
			t.Errorf("forensics report contract does not map typed field %q", fieldName)
		}
	}
	if !strings.Contains(action, "Crash Recovery (Step 1)") {
		t.Error("forensics action does not delegate stuck-work takeover judgment to Crash Recovery")
	}
	for _, retainedAuthority := range []string{"Recurring Corrections (judgment-owned)", "Release and Queue Invariants (board-owned)", "queue-kanban verify", "`skipped_work`", "## Skipped or Unverified Coverage", "0 critical, 0 warnings, 0 info items found."} {
		if !strings.Contains(action, retainedAuthority) {
			t.Errorf("forensics action omitted retained authority %q", retainedAuthority)
		}
	}
	for _, unsupportedTotal := range []string{"**Queue:**", "**Archive:**", "**Working:**"} {
		if strings.Contains(action, unsupportedTotal) {
			t.Errorf("forensics report still requires unsupported total %q", unsupportedTotal)
		}
	}

	deletedCheckReference := regexp.MustCompile(`(?i)\bchecks?\s+(?:[1-9]|1[123])\b`)
	consumerPaths := []string{
		filepath.Join("..", "..", "..", "..", "actions", "work-reference.md"),
		filepath.Join("..", "..", "..", "..", "actions", "abandon.md"),
		filepath.Join("..", "..", "..", "..", "scripts", "repair-req-timestamps.sh"),
		filepath.Join("..", "..", "..", "..", "..", "do-work-board", "tools", "queue-kanban", "verify.go"),
		filepath.Join("..", "..", "..", "..", "..", "do-work-board", "tools", "queue-kanban", "model.go"),
		filepath.Join("..", "..", "..", "..", "..", "do-work-board", "tools", "queue-kanban", "lessons-do-kanban.md"),
	}
	for _, consumerPath := range consumerPaths {
		contents, readError := os.ReadFile(consumerPath)
		if readError != nil {
			t.Fatal(readError)
		}
		if staleReference := deletedCheckReference.Find(contents); staleReference != nil {
			t.Errorf("%s retains deleted mechanical reference %q", consumerPath, staleReference)
		}
	}
	workReferenceBytes, err := os.ReadFile(consumerPaths[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(workReferenceBytes), "## Crash Recovery (Step 1)") {
		t.Fatal("forensics stuck-work authority does not resolve to the Crash Recovery heading")
	}
}

func TestForensicsActionPlacesAndCountsBoardFindingsDeterministically(t *testing.T) {
	actionPath := filepath.Join("..", "..", "..", "..", "actions", "forensics.md")
	actionBytes, err := os.ReadFile(actionPath)
	if err != nil {
		t.Fatal(err)
	}
	verifyPath := filepath.Join("..", "..", "..", "..", "..", "do-work-board", "tools", "queue-kanban", "verify.go")
	verifyBytes, err := os.ReadFile(verifyPath)
	if err != nil {
		t.Fatal(err)
	}
	boardFindingCategories := []string{"version-changelog-mismatch"}
	categoryDeclaration := regexp.MustCompile(`verifyCategoryVersionChangelogMismatch\s*=\s*"` + regexp.QuoteMeta(boardFindingCategories[0]) + `"`)
	if !categoryDeclaration.Match(verifyBytes) {
		t.Fatalf("board fixture category %q is not emitted by queue-kanban verify", boardFindingCategories[0])
	}

	action := string(actionBytes)
	const boardFindingMapping = "| Each finding emitted by `queue-kanban verify` | `## Warnings` | +1 warning |"
	mappingPresent := strings.Contains(action, boardFindingMapping)
	if !mappingPresent {
		t.Errorf("forensics action has no authoritative board-finding report path %q", boardFindingMapping)
	}
	boardFindingExample := "- **[" + boardFindingCategories[0] + "]** [board-emitted detail]"
	warningsStart := strings.Index(action, "\n## Warnings\n")
	infoStart := strings.Index(action, "\n## Info\n")
	exampleStart := strings.Index(action, boardFindingExample)
	exampleInWarnings := warningsStart >= 0 && exampleStart > warningsStart && infoStart > exampleStart
	if !exampleInWarnings {
		t.Errorf("board finding %q is not shown under ## Warnings", boardFindingCategories[0])
	}
	warningTotal := 0
	if mappingPresent && exampleInWarnings {
		warningTotal = len(boardFindingCategories)
	}
	if warningTotal != 1 {
		t.Fatalf("one board finding contributed %d to the warning total, want 1", warningTotal)
	}
}

func TestDoctorOptionGrammarKeepsMutationIntentExplicit(t *testing.T) {
	for _, arguments := range [][]string{{"--dry-run"}, {"--commit"}, {"--repair-timestamps", "--dry-run", "--commit"}, {"--unknown"}} {
		if _, err := parseCommandOptions(arguments); err == nil {
			t.Fatalf("parseCommandOptions(%v) succeeded", arguments)
		}
	}
	for _, arguments := range [][]string{nil, {"--repair-timestamps"}, {"--repair-timestamps", "--dry-run"}, {"--repair-timestamps", "--commit"}} {
		if _, err := parseCommandOptions(arguments); err != nil {
			t.Fatalf("parseCommandOptions(%v): %v", arguments, err)
		}
	}
}

func TestDoctorTextAndJSONRenderFromTheSameFindingResult(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-040-bad.md", doctorRequest("REQ-040", "pnding", "created_at: 2099-01-01T00:00:00Z\n", "Body"))
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(t.Context(), snapshot, ScanOptions{Now: time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)})
	textOutput, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	jsonOutput, err := resultmodel.RenderResult(result, resultmodel.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var rendered resultmodel.CommandResult
	if err := json.Unmarshal(jsonOutput, &rendered); err != nil {
		t.Fatal(err)
	}
	if len(rendered.Findings) != len(result.Findings) {
		t.Fatalf("JSON findings=%d result findings=%d", len(rendered.Findings), len(result.Findings))
	}
	for _, finding := range result.Findings {
		if !bytes.Contains(textOutput, []byte(finding.Code)) {
			t.Fatalf("text omitted %s:\n%s", finding.Code, textOutput)
		}
	}
}
