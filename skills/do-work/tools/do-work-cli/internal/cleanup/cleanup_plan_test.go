package cleanup

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestBuildPlanCoversStrandedArchiveURClosureAndConsumedRuns(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-101-done.md", cleanupRequest("REQ-101", "done", "UR-081"))
	writeCleanupFile(t, repositoryRoot, "do-work/user-requests/UR-081/input.md", "---\nid: UR-081\n---\nInput\n")
	writeCleanupFile(t, repositoryRoot, "do-work/runs/spent/manifest.md", "Status: consumed\n")
	writeCleanupFile(t, repositoryRoot, "do-work/runs/live/manifest.md", "Status: synthesized\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(snapshot)
	if !planHasOperation(plan, OperationMove, "do-work/queue/REQ-101-done.md", "do-work/archive/UR-081/REQ-101-done.md") {
		t.Fatalf("stranded request move missing: %#v", plan.Groups)
	}
	if !planHasOperation(plan, OperationMove, "do-work/user-requests/UR-081/input.md", "do-work/archive/UR-081/input.md") {
		t.Fatalf("UR closure missing: %#v", plan.Groups)
	}
	if !planHasOperation(plan, OperationDelete, "do-work/runs/spent/manifest.md", "") {
		t.Fatalf("consumed run sweep missing: %#v", plan.Groups)
	}
	if planHasOperation(plan, OperationDelete, "do-work/runs/live/manifest.md", "") {
		t.Fatal("non-consumed run was planned for deletion")
	}
}

func TestBuildPlanConsolidatesLooseArchiveWithoutURIntoLegacy(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/archive/REQ-102-old.md", cleanupRequest("REQ-102", "completed", ""))
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(snapshot)
	if !planHasOperation(plan, OperationMove, "do-work/archive/REQ-102-old.md", "do-work/archive/legacy/REQ-102-old.md") {
		t.Fatalf("legacy move missing: %#v", plan.Groups)
	}
}

func TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if dotIndex := strings.IndexByte(hostname, '.'); dotIndex >= 0 {
		hostname = hostname[:dotIndex]
	}
	writeCleanupFile(t, repositoryRoot, "do-work/working/REQ-109-done.md", cleanupRequest("REQ-109", "completed", ""))
	checkpoint := "# Session Checkpoint\n\n## Session Notes\n\n" +
		"- The real `## In Progress (interrupted)` section follows.\n\n" +
		"## In Progress (interrupted)\n\n" +
		"- REQ-109: own — writer: " + hostname + ":" + repositoryRoot + "\n" +
		"  Last known state: implementing\n" +
		"  Key files being modified: cleanup_plan.go\n" +
		"- REQ-109: foreign — writer: other:/checkout\n" +
		"  Last known state: foreign work\n" +
		"  Key files being modified: foreign.go\n"
	writeCleanupFile(t, repositoryRoot, "do-work/CHECKPOINT.md", checkpoint)
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(snapshot)
	var checkpointContents string
	for _, group := range plan.Groups {
		for _, operation := range group.Operations {
			if operation.Kind == OperationReplace && operation.SourcePath == "do-work/CHECKPOINT.md" {
				checkpointContents = string(operation.Contents)
			}
		}
	}
	if strings.Contains(checkpointContents, "own —") || strings.Contains(checkpointContents, "implementing") || strings.Contains(checkpointContents, "cleanup_plan.go") {
		t.Fatalf("checkpoint retained own entry bytes = %q", checkpointContents)
	}
	if !strings.Contains(checkpointContents, "The real `## In Progress (interrupted)` section follows.") ||
		!strings.Contains(checkpointContents, "foreign —") ||
		!strings.Contains(checkpointContents, "Last known state: foreign work") ||
		!strings.Contains(checkpointContents, "Key files being modified: foreign.go") {
		t.Fatalf("checkpoint replacement = %q", checkpointContents)
	}
}

func TestURClosureUsesVacuousMembershipExcludesHoldAndReportsMissingCaptureIDs(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/user-requests/UR-201/input.md", "---\nid: UR-201\nrequests: [REQ-999]\n---\nInput\n")
	writeCleanupFile(t, repositoryRoot, "do-work/archive/hold/REQ-202-held.md", cleanupRequest("REQ-202", "pending", "UR-201"))
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(snapshot)
	if !planHasOperation(plan, OperationMove, "do-work/user-requests/UR-201/input.md", "do-work/archive/UR-201/input.md") {
		t.Fatalf("zero-active-member UR did not close: %#v", plan.Groups)
	}
	foundMissing := false
	for _, finding := range plan.Findings {
		if finding.Code == "CLEANUP-CAPTURE-ID-MISSING" {
			foundMissing = true
		}
	}
	if !foundMissing {
		t.Fatalf("missing capture-array ID was not reported: %#v", plan.Findings)
	}
}

func TestLegacyContextIsConsolidatedAndMisplacedItemsAreIndependent(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/archive/CONTEXT-001-old.md", "legacy context\n")
	writeCleanupFile(t, repositoryRoot, "nested/do-work/queue/REQ-203-one.md", cleanupRequest("REQ-203", "pending", ""))
	writeCleanupFile(t, repositoryRoot, "nested/do-work/queue/REQ-204-two.md", cleanupRequest("REQ-204", "pending", ""))
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-203-one.md", cleanupRequest("REQ-203", "pending", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(snapshot)
	if !planHasOperation(plan, OperationMove, "do-work/archive/CONTEXT-001-old.md", "do-work/archive/legacy/CONTEXT-001-old.md") {
		t.Fatalf("legacy context move missing: %#v", plan.Groups)
	}
	var itemGroups int
	for _, group := range plan.Groups {
		if strings.HasPrefix(group.Code, "RELOCATE-") {
			itemGroups++
		}
	}
	if itemGroups != 2 {
		t.Fatalf("misplaced items share a conflict domain: groups=%#v", plan.Groups)
	}
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("conflict-aware result = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "queue", "REQ-204-two.md")); err != nil {
		t.Fatalf("nonconflicting misplaced item did not relocate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "nested", "do-work", "queue", "REQ-203-one.md")); err != nil {
		t.Fatalf("conflicting item did not remain: %v", err)
	}
}

func TestMisplacedArchivedURMovesNonconflictingSiblingAndRefusesOnlyConflict(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/archive/user-requests/UR-301/input.md", "misplaced input\n")
	writeCleanupFile(t, repositoryRoot, "do-work/archive/user-requests/UR-301/REQ-301-done.md", cleanupRequest("REQ-301", "completed", "UR-301"))
	writeCleanupFile(t, repositoryRoot, "do-work/archive/UR-301/input.md", "canonical input\n")
	commitCleanupFixture(t, repositoryRoot)

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plan := BuildPlan(snapshot)
	var groupCodes []string
	for _, group := range plan.Groups {
		if strings.HasPrefix(group.Code, "FIX-ARCHIVE-UR-301-") {
			groupCodes = append(groupCodes, group.Code)
			if len(group.Operations) != 1 {
				t.Fatalf("group %q has %d operations, want one per conflict domain", group.Code, len(group.Operations))
			}
		}
	}
	if len(groupCodes) != 2 || groupCodes[0] != "FIX-ARCHIVE-UR-301-REQ-301-done.md" || groupCodes[1] != "FIX-ARCHIVE-UR-301-input.md" {
		t.Fatalf("misplaced archived UR group codes = %#v", groupCodes)
	}

	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("partial-merge outcome = %q, want %q", result.Outcome, resultmodel.OutcomeFindings)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "archive", "UR-301", "REQ-301-done.md")); err != nil {
		t.Fatalf("nonconflicting sibling did not move: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work", "archive", "user-requests", "UR-301", "REQ-301-done.md")); !os.IsNotExist(err) {
		t.Fatalf("nonconflicting sibling source still exists: %v", err)
	}
	if contents, err := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "archive", "UR-301", "input.md")); err != nil || string(contents) != "canonical input\n" {
		t.Fatalf("canonical input changed: contents=%q err=%v", contents, err)
	}
	if contents, err := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "archive", "user-requests", "UR-301", "input.md")); err != nil || string(contents) != "misplaced input\n" {
		t.Fatalf("conflicting input source changed: contents=%q err=%v", contents, err)
	}

	var conflictEvidence []string
	for _, finding := range result.Findings {
		if finding.Code == "CLEANUP-GROUP-REFUSED" && len(finding.AffectedPaths) == 1 && finding.AffectedPaths[0] == "do-work/archive/UR-301/input.md" {
			conflictEvidence = finding.Evidence
		}
	}
	if len(conflictEvidence) != 1 || !strings.Contains(conflictEvidence[0], "destination already exists; cleanup never overwrites") {
		t.Fatalf("conflict evidence = %#v, findings=%#v", conflictEvidence, result.Findings)
	}
}

func planHasOperation(plan CleanupPlan, kind OperationKind, source, destination string) bool {
	for _, group := range plan.Groups {
		for _, operation := range group.Operations {
			if operation.Kind == kind && operation.SourcePath == source && operation.DestinationPath == destination {
				return true
			}
		}
	}
	return false
}

func cleanupRequest(id, status, userRequest string) string {
	userLine := ""
	if userRequest != "" {
		userLine = "user_request: " + userRequest + "\n"
	}
	return "---\nid: " + id + "\nstatus: " + status + "\n" + userLine + "---\nBody\n"
}

func cleanupRepository(t *testing.T) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	repositoryRoot := t.TempDir()
	runCleanupGit(t, repositoryRoot, "init", "-q")
	runCleanupGit(t, repositoryRoot, "config", "user.name", "Cleanup Test")
	runCleanupGit(t, repositoryRoot, "config", "user.email", "cleanup@example.invalid")
	return repositoryRoot
}

func writeCleanupFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitCleanupFixture(t *testing.T, repositoryRoot string) {
	t.Helper()
	runCleanupGit(t, repositoryRoot, "add", "-A")
	runCleanupGit(t, repositoryRoot, "commit", "-q", "-m", "fixture")
}

func runCleanupGit(t *testing.T, repositoryRoot string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}
