package doctor

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDoctorCommandIsRegisteredAndDefaultRunIsReadOnly(t *testing.T) {
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
