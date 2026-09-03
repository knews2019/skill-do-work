package nextselection

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type commandSelectionRecord struct {
	RequestID         string `json:"request_id"`
	Code              string `json:"code"`
	SelectionPriority string `json:"selection_priority"`
	RequestPriority   string `json:"priority"`
}

func TestNextCommandProjectsGateSelectionPriority(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("go-run integration is heavy-only")
	}
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-811-ordinary.md", "REQ-811", "pending", "priority: now\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-812-deferred.md", "REQ-812", "pending", "gate_deferred: true\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-813-repair.md", "REQ-813", "pending", "repository_gate_repair: true\n")
	command := exec.Command("go", "run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", "json", "next", "--fan-out", "3")
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("next command returned %v:\n%s", runError, output)
	}
	var result commandSelectionResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatal(err)
	}
	if got := selectedRequestIDs(result.Selected); !equalStrings(got, []string{"REQ-813", "REQ-812", "REQ-811"}) {
		t.Fatalf("priority ids = %v", got)
	}
	if result.Selected[0].SelectionPriority != PriorityRepositoryGateRepair || result.Selected[1].SelectionPriority != PriorityDeferredParent || result.Selected[2].SelectionPriority != PriorityOrdinary {
		t.Fatalf("priority evidence = %#v", result.Selected)
	}
	if result.Selected[0].RequestPriority != RequestPriorityNext || result.Selected[1].RequestPriority != RequestPriorityNext || result.Selected[2].RequestPriority != RequestPriorityNow {
		t.Fatalf("request priority evidence = %#v", result.Selected)
	}
}

func TestNextCommandProjectsPriorityOnSelectedAndFanOutExcluded(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("go-run integration is heavy-only")
	}
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-821-later.md", "REQ-821", "pending", "priority: later\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-822-now.md", "REQ-822", "pending", "priority: now\n")
	command := exec.Command("go", "run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", "json", "next", "--fan-out", "1")
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("next command returned %v:\n%s", runError, output)
	}
	var result commandSelectionResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Selected) != 1 || result.Selected[0].RequestID != "REQ-822" || result.Selected[0].RequestPriority != RequestPriorityNow {
		t.Fatalf("selected priority projection = %#v", result.Selected)
	}
	if len(result.Excluded) != 1 || result.Excluded[0].RequestID != "REQ-821" || result.Excluded[0].Code != "FAN-OUT-LIMIT" || result.Excluded[0].RequestPriority != RequestPriorityLater {
		t.Fatalf("excluded priority projection = %#v", result.Excluded)
	}
}

type commandSelectionResult struct {
	Selected []commandSelectionRecord `json:"selected"`
	Excluded []commandSelectionRecord `json:"excluded"`
}

func TestNextCommandMixedFixture(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("go-run integration is heavy-only")
	}
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/archive/REQ-900-done.md", "REQ-900", "completed", "")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-101-root.md", "REQ-101", "pending", "estimate:\n  p50_active_minutes: 10\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-102-chain.md", "REQ-102", "pending", "depends_on: [REQ-101]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-103-cycle.md", "REQ-103", "pending", "depends_on: [REQ-104]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-104-cycle.md", "REQ-104", "pending", "depends_on: [REQ-103]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-105-assigned.md", "REQ-105", "pending", "assigned_to: cloud-alpha\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-106-negligible.md", "REQ-106", "pending", "impact: impact-negligible\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-107-ready.md", "REQ-107", "pending", "depends_on: [REQ-900]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-108-missing.md", "REQ-108", "pending", "depends_on: [REQ-999]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-109-claimed.md", "REQ-109", "claimed", "claimed_at: 2026-08-31T10:00:00Z\n")

	command := exec.Command("go", "run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", "json", "next", "--fan-out", "2", "--skip-impact-negligible")
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("next command returned %v:\n%s", runError, output)
	}
	var result commandSelectionResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode next JSON: %v\n%s", err, output)
	}
	if got := selectedRequestIDs(result.Selected); !equalStrings(got, []string{"REQ-101", "REQ-107"}) {
		t.Fatalf("selected ids = %v, want [REQ-101 REQ-107]", got)
	}
	wantExclusions := map[string]string{
		"REQ-102": "DEPENDENCIES-UNMET",
		"REQ-103": "DEPENDENCY-CYCLE",
		"REQ-104": "DEPENDENCY-CYCLE",
		"REQ-105": "ASSIGNED-ELSEWHERE",
		"REQ-106": "IMPACT-NEGLIGIBLE",
		"REQ-108": "DEPENDENCY-MISSING",
		"REQ-109": "ALREADY-CLAIMED",
	}
	for _, exclusion := range result.Excluded {
		deleteIfEqual(wantExclusions, exclusion.RequestID, exclusion.Code)
	}
	if len(wantExclusions) != 0 {
		t.Fatalf("missing exclusions: %v; result=%#v", wantExclusions, result.Excluded)
	}
}

func TestNextCommandTreatsEmptyInlineDependencyListAsReady(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("go-run integration is heavy-only")
	}
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-110-empty-dependency-list.md", "REQ-110", "pending", "depends_on: []\n")

	command := exec.Command("go", "run", "../../cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", "json", "next")
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("next command returned %v:\n%s", runError, output)
	}
	var result commandSelectionResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode next JSON: %v\n%s", err, output)
	}
	if got := selectedRequestIDs(result.Selected); !equalStrings(got, []string{"REQ-110"}) {
		t.Fatalf("selected ids = %v, want [REQ-110]; exclusions=%#v", got, result.Excluded)
	}
	for _, exclusion := range result.Excluded {
		if exclusion.RequestID == "REQ-110" && exclusion.Code == "DEPENDENCY-MISSING" {
			t.Fatalf("empty list was treated as a missing dependency: %#v", exclusion)
		}
	}
}

func TestNextCommandUsesInProcessBlockedProbeAuthority(t *testing.T) {
	status, err := RunBlockedProbe([]byte("exit 19"), 2)
	if err != nil || status != 19 {
		t.Fatalf("raw status=%d err=%v", status, err)
	}
}

func writeCommandRequest(t *testing.T, repositoryRoot, relativePath, requestID, status, extra string) {
	t.Helper()
	requestPath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(requestPath), 0o755); err != nil {
		t.Fatal(err)
	}
	contents := "---\nid: " + requestID + "\ntitle: Fixture " + requestID + "\nstatus: " + status + "\n" + extra + "---\nBody\n"
	if err := os.WriteFile(requestPath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func selectedRequestIDs(records []commandSelectionRecord) []string {
	requestIDs := make([]string, len(records))
	for index, record := range records {
		requestIDs[index] = record.RequestID
	}
	return requestIDs
}

func deleteIfEqual(values map[string]string, key, actual string) {
	if values[key] == actual {
		delete(values, key)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
