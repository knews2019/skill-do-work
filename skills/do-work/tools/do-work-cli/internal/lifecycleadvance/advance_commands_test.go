package lifecycleadvance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
)

var (
	advanceCLIBinaryOnce sync.Once
	advanceCLIBinaryPath string
	advanceCLIBinaryErr  error
)

type advanceCommandResult struct {
	Outcome  string                  `json:"outcome"`
	Advance  *advanceCommandEvidence `json:"advance"`
	Findings []advanceCommandFinding `json:"findings"`
}

type advanceCommandEvidence struct {
	RequestID        string                   `json:"request_id"`
	RequestPath      string                   `json:"request_path"`
	TreeSection      string                   `json:"tree_section"`
	Status           string                   `json:"status"`
	Route            string                   `json:"route"`
	Phase            string                   `json:"phase"`
	PhaseKind        string                   `json:"phase_kind"`
	MissingEvidence  []advanceMissingEvidence `json:"missing_evidence"`
	NextArgv         []string                 `json:"next_argv"`
	VerificationArgv []string                 `json:"verification_argv"`
}

type advanceMissingEvidence struct {
	Kind     string `json:"kind"`
	Path     string `json:"path"`
	Field    string `json:"field"`
	Section  string `json:"section"`
	Expected string `json:"expected"`
}

type advanceCommandFinding struct {
	Code                 string   `json:"code"`
	AffectedIDs          []string `json:"affected_ids"`
	AffectedPaths        []string `json:"affected_paths"`
	Evidence             []string `json:"observed_evidence"`
	AutomationStopReason string   `json:"automation_stop_reason"`
	NextArgv             []string `json:"next_argv"`
	VerificationArgv     []string `json:"verification_argv"`
}

func TestAdvanceCommandPhaseMatrix(t *testing.T) {
	tests := []struct {
		name            string
		treeSection     string
		status          string
		frontmatter     string
		body            string
		wantPhase       string
		wantPhaseKind   string
		wantNextCommand string
		wantMissing     advanceMissingEvidence
	}{
		{
			name: "pending queue claim", treeSection: "queue", status: "pending",
			wantPhase: "claim", wantPhaseKind: "mechanical", wantNextCommand: "claim",
			wantMissing: advanceMissingEvidence{Kind: "field", Field: "status", Expected: "claimed"},
		},
		{
			name: "blocked queue re-probe", treeSection: "queue", status: "blocked", frontmatter: "blocked_by: service unavailable\nblocked_check: exit 0\n",
			wantPhase: "blocked-check", wantPhaseKind: "mechanical", wantNextCommand: "next",
			wantMissing: advanceMissingEvidence{Kind: "field", Field: "status", Expected: "pending after successful probe and unblock"},
		},
		{
			name: "claimed without triage", treeSection: "working", status: "claimed",
			wantPhase: "agent judgment: triage and open questions", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "field", Field: "route", Expected: "A, B, or C"},
		},
		{
			name: "route A needs estimate", treeSection: "working", status: "claimed", frontmatter: "route: A\n",
			body:      "## Triage\n\nRoute A.\n",
			wantPhase: "estimate-p50", wantPhaseKind: "mechanical", wantNextCommand: "estimate-p50",
			wantMissing: advanceMissingEvidence{Kind: "field", Field: "estimate.p50_active_minutes", Expected: "positive integer"},
		},
		{
			name: "route A records planning skip", treeSection: "working", status: "claimed", frontmatter: "route: A\nestimate:\n  p50_active_minutes: 5\n",
			body:      "## Triage\n\nRoute A.\n",
			wantPhase: "agent judgment: record planning not required", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Plan", Expected: "planning not required"},
		},
		{
			name: "route A implementation", treeSection: "working", status: "claimed", frontmatter: "route: A\nestimate:\n  p50_active_minutes: 5\n",
			body:      "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n",
			wantPhase: "agent judgment: implementation and summary", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Implementation Summary", Expected: "implemented file manifest"},
		},
		{
			name: "route A testing skips scope drift", treeSection: "working", status: "claimed", frontmatter: "route: A\nestimate:\n  p50_active_minutes: 5\n",
			body:      "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n\n## Implementation Summary\n\n- owned.go\n\n## Qualification\n\nPASS.\n",
			wantPhase: "agent judgment: testing", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Testing", Expected: "tests and gate evidence"},
		},
		{
			name: "route B records planning skip", treeSection: "working", status: "claimed", frontmatter: "route: B\nestimate:\n  p50_active_minutes: 20\n",
			body:      "## Triage\n\nRoute B.\n",
			wantPhase: "agent judgment: record planning not required", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Plan", Expected: "planning not required"},
		},
		{
			name: "route B exploration", treeSection: "working", status: "claimed", frontmatter: "route: B\nestimate:\n  p50_active_minutes: 20\n",
			body:      "## Triage\n\nRoute B.\n\n## Plan\n\nPlanning not required.\n",
			wantPhase: "agent judgment: exploration", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Exploration", Expected: "exploration findings"},
		},
		{
			name: "route C planning", treeSection: "working", status: "claimed", frontmatter: "route: C\nestimate:\n  p50_active_minutes: 30\n",
			body:      "## Triage\n\nRoute C.\n",
			wantPhase: "agent judgment: planning", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Plan", Expected: "validated Route C plan"},
		},
		{
			name: "route C exploration", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nestimate:\n  p50_active_minutes: 30\n",
			body:      "## Triage\n\nRoute C.\n\n## Plan\n\nSpecific plan.\n",
			wantPhase: "agent judgment: exploration", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Exploration", Expected: "exploration findings"},
		},
		{
			name: "route C scope", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nestimate:\n  p50_active_minutes: 30\n",
			body:      "## Triage\n\nRoute C.\n\n## Plan\n\nSpecific plan.\n\n## Exploration\n\nFound patterns.\n",
			wantPhase: "agent judgment: scope declaration", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Scope", Expected: "declared files and acceptance criteria"},
		},
		{
			name: "planned route C preflight", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			body:      routeCBodyThrough("Scope"),
			wantPhase: "preflight", wantPhaseKind: "mechanical", wantNextCommand: "preflight",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Implementation Summary", Expected: "preflight complete and implementation summarized"},
		},
		{
			name: "implemented without qualification", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			body:      routeCBodyThrough("Implementation Summary"),
			wantPhase: "qualify", wantPhaseKind: "mechanical", wantNextCommand: "qualify",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Qualification", Expected: "typed qualification result"},
		},
		{
			name: "qualified route C scope drift", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			body:      routeCBodyThrough("Qualification"),
			wantPhase: "scope-drift", wantPhaseKind: "mechanical", wantNextCommand: "scope-drift",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Testing", Expected: "tests and gate evidence"},
		},
		{
			name: "tested route C review", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			body:      routeCBodyThrough("Testing"),
			wantPhase: "agent judgment: review", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Review", Expected: "independent review verdict"},
		},
		{
			name: "reviewed route C lessons", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			body:      routeCBodyThrough("Review"),
			wantPhase: "agent judgment: lessons and orientation", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "section", Section: "Lessons Learned", Expected: "captured lessons"},
		},
		{
			name: "ready route C finalization", treeSection: "working", status: "claimed", frontmatter: "route: C\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 30\n",
			body:      routeCBodyThrough("Orientation"),
			wantPhase: "agent judgment: prepare finalization manifest", wantPhaseKind: "agent_judgment",
			wantMissing: advanceMissingEvidence{Kind: "file", Expected: "action-authored finalization manifest"},
		},
		{
			name: "archived without provenance", treeSection: "archive", status: "completed", frontmatter: "route: C\ncompleted_at: 2026-09-04T13:00:00Z\n",
			body:      routeCBodyThrough("Orientation"),
			wantPhase: "recover-finalization", wantPhaseKind: "mechanical", wantNextCommand: "recover-finalization",
			wantMissing: advanceMissingEvidence{Kind: "field", Field: "commit", Expected: "recorded implementation or primary commit"},
		},
		{
			name: "archived complete", treeSection: "archive", status: "completed", frontmatter: "route: C\ncompleted_at: 2026-09-04T13:00:00Z\ncommit: abc1234\n",
			body:      routeCBodyThrough("Orientation"),
			wantPhase: "complete", wantPhaseKind: "complete",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			requestPath := writeAdvanceRequest(t, repositoryRoot, test.treeSection, "REQ-703", test.status, test.frontmatter, test.body)
			result, _ := runAdvanceJSON(t, repositoryRoot, "REQ-703")
			if result.Outcome != "success" || result.Advance == nil {
				t.Fatalf("outcome=%q advance=%#v findings=%#v", result.Outcome, result.Advance, result.Findings)
			}
			advance := result.Advance
			if advance.RequestID != "REQ-703" || advance.RequestPath != requestPath || advance.TreeSection != test.treeSection || advance.Status != test.status {
				t.Fatalf("identity/state = %#v", advance)
			}
			if advance.Phase != test.wantPhase || advance.PhaseKind != test.wantPhaseKind {
				t.Fatalf("phase=%q kind=%q, want %q %q", advance.Phase, advance.PhaseKind, test.wantPhase, test.wantPhaseKind)
			}
			wantNextArgv := expectedAdvanceNextArgv(test.wantNextCommand, requestPath, advance.Route)
			if !reflect.DeepEqual(advance.NextArgv, wantNextArgv) {
				t.Fatalf("next argv=%#v, want %#v", advance.NextArgv, wantNextArgv)
			}
			if len(advance.MissingEvidence) == 0 {
				if test.wantMissing != (advanceMissingEvidence{}) {
					t.Fatalf("missing evidence empty, want %#v", test.wantMissing)
				}
			} else {
				got := advance.MissingEvidence[0]
				got.Path = ""
				if got != test.wantMissing {
					t.Fatalf("missing evidence=%#v, want %#v", got, test.wantMissing)
				}
			}
			if !reflect.DeepEqual(advance.VerificationArgv, []string{"do-work-cli", "--format", "json", "advance", "REQ-703"}) {
				t.Fatalf("verification argv=%#v", advance.VerificationArgv)
			}
		})
	}
}

func TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates(t *testing.T) {
	tests := []struct {
		name     string
		seed     func(*testing.T, string)
		wantCode string
	}{
		{
			name: "malformed request",
			seed: func(t *testing.T, root string) {
				writeAdvanceFile(t, root, "do-work/working/REQ-704-malformed.md", "not frontmatter\n")
			},
			wantCode: "ADVANCE-EVIDENCE-MISSING",
		},
		{
			name: "ambiguous request",
			seed: func(t *testing.T, root string) {
				writeAdvanceRequest(t, root, "queue", "REQ-704", "pending", "", "")
				writeAdvanceRequest(t, root, "archive", "REQ-704", "completed", "commit: abc1234\n", "")
			},
			wantCode: "ADVANCE-EVIDENCE-MISSING",
		},
		{
			name: "unknown route",
			seed: func(t *testing.T, root string) {
				writeAdvanceRequest(t, root, "working", "REQ-704", "claimed", "route: D\n", "## Triage\n")
			},
			wantCode: "ADVANCE-PHASE-UNKNOWN",
		},
		{
			name: "contradictory section order",
			seed: func(t *testing.T, root string) {
				writeAdvanceRequest(t, root, "working", "REQ-704", "claimed", "route: C\nestimate:\n  p50_active_minutes: 30\n", "## Triage\n\n## Implementation Summary\n\n- `owned.go` (modified)\n")
			},
			wantCode: "ADVANCE-EVIDENCE-MISSING",
		},
		{
			name: "impossible archive status",
			seed: func(t *testing.T, root string) {
				writeAdvanceRequest(t, root, "archive", "REQ-704", "claimed", "route: C\n", "## Triage\n")
			},
			wantCode: "ADVANCE-PHASE-UNKNOWN",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			test.seed(t, repositoryRoot)
			result, status := runAdvanceJSON(t, repositoryRoot, "REQ-704")
			if status != 1 || result.Outcome != "refused" || len(result.Findings) != 1 || result.Findings[0].Code != test.wantCode {
				t.Fatalf("status=%d result=%#v", status, result)
			}
			finding := result.Findings[0]
			if !reflect.DeepEqual(finding.AffectedIDs, []string{"REQ-704"}) || len(finding.VerificationArgv) == 0 || commandVerb(finding.VerificationArgv) != "advance" {
				t.Fatalf("finding lost typed consumer fields: %#v", finding)
			}
		})
	}
}

func TestAdvanceCommandIsByteForByteReadOnlyInTextAndJSON(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-705", "pending", "", "")
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\nKeep these bytes.\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Advance Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "advance@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")

	before := advanceTreeDigest(t, repositoryRoot)
	for _, format := range []string{"text", "json"} {
		command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", format, "advance", "REQ-705")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("advance %s failed: %v\n%s", format, err, output)
		}
	}
	after := advanceTreeDigest(t, repositoryRoot)
	if before != after {
		t.Fatalf("repository bytes changed: before=%s after=%s", before, after)
	}
	if status := strings.TrimSpace(string(runAdvanceGit(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=all"))); status != "" {
		t.Fatalf("advance dirtied Git state: %s", status)
	}
}

func TestAdvanceCommandRejectsInvalidArguments(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-706", "pending", "", "")
	for _, arguments := range [][]string{{}, {"REQ-706", "REQ-707"}, {"REQ-x"}} {
		commandArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "advance"}
		commandArguments = append(commandArguments, arguments...)
		command := exec.Command(advanceCLIBinary(t), commandArguments...)
		output, err := command.CombinedOutput()
		if err == nil {
			t.Fatalf("advance accepted %#v: %s", arguments, output)
		}
		var result advanceCommandResult
		if decodeError := json.Unmarshal(output, &result); decodeError != nil || result.Outcome != "failure" || len(result.Findings) != 1 || result.Findings[0].Code != "ADVANCE-USAGE" {
			t.Fatalf("arguments=%#v result=%#v decode=%v output=%s", arguments, result, decodeError, output)
		}
	}
}

func routeCBodyThrough(lastSection string) string {
	sections := []struct {
		name string
		body string
	}{
		{"Triage", "Route C."},
		{"Plan", "Specific plan."},
		{"Exploration", "Found patterns."},
		{"Scope", "**Files I will touch:**\n- `owned.go`"},
		{"Implementation Summary", "- `owned.go` (modified)"},
		{"Qualification", "Typed qualification passed."},
		{"Testing", "Gate passed."},
		{"Review", "Acceptance: Pass"},
		{"Lessons Learned", "Worth knowing."},
		{"Orientation", "Lifecycle command added."},
	}
	var body strings.Builder
	for _, section := range sections {
		fmt.Fprintf(&body, "## %s\n\n%s\n\n", section.name, section.body)
		if section.name == lastSection {
			break
		}
	}
	return body.String()
}

func writeAdvanceRequest(t *testing.T, repositoryRoot, treeSection, requestID, status, frontmatter, body string) string {
	t.Helper()
	requestPath := filepath.ToSlash(filepath.Join("do-work", treeSection, requestID+"-fixture.md"))
	contents := "---\nid: " + requestID + "\ntitle: Fixture " + requestID + "\nstatus: " + status + "\n" + frontmatter + "---\n\n" + body
	writeAdvanceFile(t, repositoryRoot, requestPath, contents)
	return requestPath
}

func writeAdvanceFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runAdvanceJSON(t *testing.T, repositoryRoot, requestID string) (advanceCommandResult, int) {
	t.Helper()
	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "advance", requestID)
	output, err := command.CombinedOutput()
	status := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
		} else {
			t.Fatalf("advance launch: %v", err)
		}
	}
	var result advanceCommandResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode advance result: %v\n%s", err, output)
	}
	return result, status
}

func advanceCLIBinary(t *testing.T) string {
	t.Helper()
	advanceCLIBinaryOnce.Do(func() {
		temporaryDirectory, err := os.MkdirTemp("", "advance-cli-test-*")
		if err != nil {
			advanceCLIBinaryErr = err
			return
		}
		advanceCLIBinaryPath = filepath.Join(temporaryDirectory, "do-work-cli")
		command := exec.Command("go", "build", "-o", advanceCLIBinaryPath, "../../cmd/do-work-cli")
		if output, buildError := command.CombinedOutput(); buildError != nil {
			advanceCLIBinaryErr = fmt.Errorf("build CLI: %w\n%s", buildError, output)
		}
	})
	if advanceCLIBinaryErr != nil {
		t.Fatal(advanceCLIBinaryErr)
	}
	return advanceCLIBinaryPath
}

func commandVerb(arguments []string) string {
	for index := 1; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--format" || argument == "--repo-root" {
			index++
			continue
		}
		if strings.HasPrefix(argument, "-") {
			continue
		}
		return argument
	}
	return ""
}

func expectedAdvanceNextArgv(command, requestPath, route string) []string {
	switch command {
	case "":
		return []string{}
	case "claim":
		return []string{"do-work-cli", "claim", "REQ-703", "--request-path", requestPath, "--provenance", "explicit-req", "--commit"}
	case "next":
		return []string{"do-work-cli", "--format", "json", "next", "REQ-703"}
	case "estimate-p50":
		return []string{"do-work-cli", "estimate-p50", "--route", route}
	case "preflight":
		return []string{"do-work-cli", "--format", "json", "preflight"}
	case "qualify", "scope-drift":
		return []string{"do-work-cli", "--format", "json", command, "--request-path", requestPath}
	case "recover-finalization":
		return []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"}
	default:
		return []string{"unexpected command " + command}
	}
}

func runAdvanceGit(t *testing.T, repositoryRoot string, arguments ...string) []byte {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = repositoryRoot
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return output
}

func advanceTreeDigest(t *testing.T, repositoryRoot string) string {
	t.Helper()
	paths := []string{}
	if err := filepath.WalkDir(repositoryRoot, func(path string, entry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if !entry.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	sort.Strings(paths)
	digest := sha256.New()
	for _, path := range paths {
		relativePath, err := filepath.Rel(repositoryRoot, path)
		if err != nil {
			t.Fatal(err)
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(digest, "%s\x00", filepath.ToSlash(relativePath))
		digest.Write(contents)
		digest.Write([]byte{0})
	}
	return hex.EncodeToString(digest.Sum(nil))
}
