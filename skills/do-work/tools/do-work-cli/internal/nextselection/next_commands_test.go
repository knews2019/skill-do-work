package nextselection

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type commandSelectionRecord struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code"`
}

type commandSelectionResult struct {
	Selected []commandSelectionRecord `json:"selected"`
	Excluded []commandSelectionRecord `json:"excluded"`
}

func TestNextCommandMixedFixture(t *testing.T) {
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
		"REQ-109": "STATUS-NOT-PENDING",
	}
	for _, exclusion := range result.Excluded {
		deleteIfEqual(wantExclusions, exclusion.RequestID, exclusion.Code)
	}
	if len(wantExclusions) != 0 {
		t.Fatalf("missing exclusions: %v; result=%#v", wantExclusions, result.Excluded)
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
