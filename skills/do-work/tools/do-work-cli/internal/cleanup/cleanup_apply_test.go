package cleanup

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDryRunReportsExactArchiveMoveAndApplyPerformsIt(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-103-finished.md", cleanupRequest("REQ-103", "finished", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plan := BuildPlan(snapshot)
	dryResult := ApplyPlan(context.Background(), plan, ApplyOptions{DryRun: true})
	if dryResult.Outcome != resultmodel.OutcomeSuccess || len(dryResult.Changes) != 1 || !strings.Contains(dryResult.Changes[0].Detail, "do-work/queue/REQ-103-finished.md") {
		t.Fatalf("dry-run result = %#v", dryResult)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-103-finished.md")); err != nil {
		t.Fatalf("dry-run mutated source: %v", err)
	}
	applyResult := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if applyResult.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("apply result = %#v", applyResult)
	}
	archivedPath := filepath.Join(repositoryRoot, "do-work/archive/REQ-103-finished.md")
	contents, err := os.ReadFile(archivedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "status: completed") {
		t.Fatalf("terminal alias not normalized: %s", contents)
	}
}

func TestDirtyGroupIsRefusedWithoutBlockingIndependentSafeGroup(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-104-done.md", cleanupRequest("REQ-104", "done", ""))
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-105-done.md", cleanupRequest("REQ-105", "done", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	plan := BuildPlan(snapshot)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-104-done.md", cleanupRequest("REQ-104", "done", "")+"user edit\n")
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/archive/REQ-105-done.md")); err != nil {
		t.Fatalf("safe group did not apply: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-104-done.md")); err != nil {
		t.Fatalf("dirty group moved: %v", err)
	}
}

func TestURClosureWaitsForRequiredMemberArchival(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-430-zulu.md", cleanupRequest("REQ-430", "completed", "UR-430"))
	writeCleanupFile(t, repositoryRoot, "do-work/working/REQ-429-alpha.md", cleanupRequest("REQ-429", "completed", "UR-430"))
	writeCleanupFile(t, repositoryRoot, "do-work/user-requests/UR-430/input.md", "---\nid: UR-430\n---\nInput\n")
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-431-unrelated.md", cleanupRequest("REQ-431", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)

	snapshot, discoveryErr := repositorymodel.DiscoverRepository(repositoryRoot)
	if discoveryErr != nil {
		t.Fatal(discoveryErr)
	}
	plan := BuildPlan(snapshot)
	var closureGroup OperationGroup
	for _, group := range plan.Groups {
		if group.Code == "CLOSE-UR-430" {
			closureGroup = group
			break
		}
	}
	wantPrerequisites := []string{"ARCHIVE-REQ-429", "ARCHIVE-REQ-430"}
	if !reflect.DeepEqual(closureGroup.RequiredGroupCodes, wantPrerequisites) {
		t.Fatalf("CLOSE-UR-430 prerequisites = %#v, want %#v", closureGroup.RequiredGroupCodes, wantPrerequisites)
	}

	writeCleanupFile(t, repositoryRoot, "do-work/working/REQ-429-alpha.md", cleanupRequest("REQ-429", "completed", "UR-430")+"user edit\n")
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	for _, retainedPath := range []string{"do-work/working/REQ-429-alpha.md", "do-work/user-requests/UR-430/input.md"} {
		if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(retainedPath))); err != nil {
			t.Fatalf("required input %s did not remain: %v", retainedPath, err)
		}
	}
	for _, refusedDestination := range []string{"do-work/archive/UR-430/REQ-429-alpha.md", "do-work/archive/UR-430/input.md"} {
		if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(refusedDestination))); !os.IsNotExist(err) {
			t.Fatalf("refused destination %s exists: %v", refusedDestination, err)
		}
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/archive/REQ-431-unrelated.md")); err != nil {
		t.Fatalf("unrelated safe group did not apply: %v", err)
	}
	foundClosureBlocker := false
	for _, finding := range result.Findings {
		if finding.Code == "CLEANUP-GROUP-REFUSED" && strings.Contains(strings.Join(finding.Evidence, " "), "CLOSE-UR-430") && strings.Contains(strings.Join(finding.Evidence, " "), "ARCHIVE-REQ-429") {
			foundClosureBlocker = true
		}
	}
	if !foundClosureBlocker {
		t.Fatalf("closure refusal did not name its blocking member: %#v", result.Findings)
	}
}

func TestURClosureAppliesAfterAllRequiredMemberArchival(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-432-one.md", cleanupRequest("REQ-432", "completed", "UR-432"))
	writeCleanupFile(t, repositoryRoot, "do-work/working/REQ-433-two.md", cleanupRequest("REQ-433", "completed", "UR-432"))
	writeCleanupFile(t, repositoryRoot, "do-work/user-requests/UR-432/input.md", "---\nid: UR-432\n---\nInput\n")
	commitCleanupFixture(t, repositoryRoot)

	snapshot, discoveryErr := repositorymodel.DiscoverRepository(repositoryRoot)
	if discoveryErr != nil {
		t.Fatal(discoveryErr)
	}
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("apply result = %#v", result)
	}
	for _, archivedPath := range []string{
		"do-work/archive/UR-432/REQ-432-one.md",
		"do-work/archive/UR-432/REQ-433-two.md",
		"do-work/archive/UR-432/input.md",
	} {
		if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(archivedPath))); err != nil {
			t.Fatalf("required archive output %s missing: %v", archivedPath, err)
		}
	}
}

func TestOperationGroupPrerequisitesFailClosed(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "direct.txt", "original\n")
	writeCleanupFile(t, repositoryRoot, "safe.txt", "safe\n")
	for _, path := range []string{"middle.txt", "transitive.txt", "missing.txt", "duplicate-one.txt", "duplicate-two.txt", "duplicate-dependent.txt", "cycle-a.txt", "cycle-b.txt", "repeated.txt"} {
		writeCleanupFile(t, repositoryRoot, path, "clean\n")
	}
	commitCleanupFixture(t, repositoryRoot)
	writeCleanupFile(t, repositoryRoot, "direct.txt", "user edit\n")

	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{
		{Code: "DIRECT", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "direct.txt", Contents: []byte("replacement\n")}}},
		{Code: "MIDDLE", RequiredGroupCodes: []string{"DIRECT"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "middle.txt", Contents: []byte("clean\n")}}},
		{Code: "TRANSITIVE", RequiredGroupCodes: []string{"MIDDLE"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "transitive.txt", Contents: []byte("clean\n")}}},
		{Code: "MISSING", RequiredGroupCodes: []string{"ABSENT"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "missing.txt", Contents: []byte("clean\n")}}},
		{Code: "DUPLICATE", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "duplicate-one.txt", Contents: []byte("clean\n")}}},
		{Code: "DUPLICATE", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "duplicate-two.txt", Contents: []byte("clean\n")}}},
		{Code: "DUPLICATE-DEPENDENT", RequiredGroupCodes: []string{"DUPLICATE"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "duplicate-dependent.txt", Contents: []byte("clean\n")}}},
		{Code: "CYCLE-A", RequiredGroupCodes: []string{"CYCLE-B"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "cycle-a.txt", Contents: []byte("clean\n")}}},
		{Code: "CYCLE-B", RequiredGroupCodes: []string{"CYCLE-A"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "cycle-b.txt", Contents: []byte("clean\n")}}},
		{Code: "REPEATED", RequiredGroupCodes: []string{"SAFE", "SAFE"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "repeated.txt", Contents: []byte("clean\n")}}},
		{Code: "SAFE", Operations: []CleanupOperation{{Kind: OperationMove, SourcePath: "safe.txt", DestinationPath: "archive/safe.txt"}}},
	}}
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "archive/safe.txt")); err != nil {
		t.Fatalf("independent safe group did not apply: %v", err)
	}
	for dependentCode, blockerCode := range map[string]string{
		"MIDDLE":              "DIRECT",
		"TRANSITIVE":          "DIRECT",
		"MISSING":             "ABSENT",
		"DUPLICATE-DEPENDENT": "DUPLICATE",
		"CYCLE-A":             "CYCLE-A",
		"CYCLE-B":             "CYCLE-A",
		"REPEATED":            "SAFE",
	} {
		foundRefusal := false
		for _, finding := range result.Findings {
			evidence := strings.Join(finding.Evidence, " ")
			if finding.Code == "CLEANUP-GROUP-REFUSED" && strings.Contains(evidence, dependentCode+": prerequisite ") && strings.Contains(evidence, blockerCode) {
				foundRefusal = true
				break
			}
		}
		if !foundRefusal {
			t.Errorf("%s refusal did not name blocker %s: %#v", dependentCode, blockerCode, result.Findings)
		}
	}
}

func TestCleanupCommitContainsOnlyExactTouchedPaths(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-110-done.md", cleanupRequest("REQ-110", "completed", ""))
	writeCleanupFile(t, repositoryRoot, "unrelated.txt", "initial\n")
	commitCleanupFixture(t, repositoryRoot)
	writeCleanupFile(t, repositoryRoot, "unrelated.txt", "user dirt\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{Commit: true, CommitMessage: "cleanup exact move"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("commit result = %#v", result)
	}
	paths := strings.Fields(runCleanupGit(t, repositoryRoot, "show", "--no-renames", "--pretty=format:", "--name-only", "HEAD"))
	sort.Strings(paths)
	want := []string{"do-work/archive/REQ-110-done.md", "do-work/queue/REQ-110-done.md"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("commit paths = %#v, want %#v", paths, want)
	}
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, "unrelated.txt"))
	if err != nil || string(contents) != "user dirt\n" {
		t.Fatalf("unrelated dirt changed: %q, %v", contents, err)
	}
}

func TestMovePublicationNeverOverwritesAnExistingDestination(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCleanupFile(t, repositoryRoot, "source.txt", "source\n")
	writeCleanupFile(t, repositoryRoot, "archive/destination.txt", "user content\n")
	if err := moveWithoutOverwrite(repositoryRoot, "source.txt", "archive/destination.txt"); err == nil {
		t.Fatal("move overwrote an existing destination")
	}
	destinationContents, err := os.ReadFile(filepath.Join(repositoryRoot, "archive", "destination.txt"))
	if err != nil || string(destinationContents) != "user content\n" {
		t.Fatalf("destination = %q, %v", destinationContents, err)
	}
	sourceContents, err := os.ReadFile(filepath.Join(repositoryRoot, "source.txt"))
	if err != nil || string(sourceContents) != "source\n" {
		t.Fatalf("source = %q, %v", sourceContents, err)
	}
}

func TestMovePublicationPreservesCompleteSourceMode(t *testing.T) {
	for _, test := range []struct {
		name string
		mode os.FileMode
	}{
		{name: "setuid", mode: 0o4640},
		{name: "setgid", mode: 0o2640},
		{name: "sticky", mode: 0o1640},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			writeCleanupFile(t, repositoryRoot, "source.txt", "source\n")
			if err := os.Chmod(filepath.Join(repositoryRoot, "source.txt"), cleanupGoModeFromUnix(test.mode)); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(filepath.Join(repositoryRoot, "archive"), 0o755); err != nil {
				t.Fatal(err)
			}

			if err := moveWithoutOverwrite(repositoryRoot, "source.txt", "archive/destination.txt"); err != nil {
				t.Fatalf("moveWithoutOverwrite: %v", err)
			}
			if _, err := os.Stat(filepath.Join(repositoryRoot, "source.txt")); !os.IsNotExist(err) {
				t.Fatalf("source remains after move: %v", err)
			}
			destinationPath := filepath.Join(repositoryRoot, "archive", "destination.txt")
			contents, err := os.ReadFile(destinationPath)
			if err != nil || string(contents) != "source\n" {
				t.Fatalf("destination = %q, %v", contents, err)
			}
			if mode := cleanupUnixModeOf(t, destinationPath); mode != test.mode {
				t.Fatalf("mode = %04o, want %04o", mode, test.mode)
			}
		})
	}
}

func cleanupGoModeFromUnix(mode os.FileMode) os.FileMode {
	goMode := mode.Perm()
	if mode&0o4000 != 0 {
		goMode |= os.ModeSetuid
	}
	if mode&0o2000 != 0 {
		goMode |= os.ModeSetgid
	}
	if mode&0o1000 != 0 {
		goMode |= os.ModeSticky
	}
	return goMode
}

func cleanupUnixModeOf(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	mode := info.Mode().Perm()
	if info.Mode()&os.ModeSetuid != 0 {
		mode |= 0o4000
	}
	if info.Mode()&os.ModeSetgid != 0 {
		mode |= 0o2000
	}
	if info.Mode()&os.ModeSticky != 0 {
		mode |= 0o1000
	}
	return mode
}

func TestRollbackDoesNotClaimAppliedChangesAndPreservesCanonicalRoots(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-205-done.md", cleanupRequest("REQ-205", "completed", ""))
	for _, rootPath := range []string{"do-work/archive", "do-work/working", "do-work/user-requests", "do-work/runs"} {
		if err := os.MkdirAll(filepath.Join(repositoryRoot, filepath.FromSlash(rootPath)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	commitCleanupFixture(t, repositoryRoot)
	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{
		{Code: "forced-rollback", Operations: []CleanupOperation{
			{Kind: OperationMove, SourcePath: "do-work/queue/REQ-205-done.md", DestinationPath: "do-work/archive/REQ-205-done.md"},
			{Kind: OperationDelete, SourcePath: "do-work/queue/missing-after-preflight.md"},
		}},
	}}
	result := ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeRolledBack {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	for _, change := range result.Changes {
		if strings.HasPrefix(change.Detail, "applied ") {
			t.Fatalf("rollback claimed applied change: %#v", result.Changes)
		}
	}
	for _, rootPath := range []string{"do-work", "do-work/queue", "do-work/working", "do-work/user-requests", "do-work/archive", "do-work/runs"} {
		if info, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(rootPath))); err != nil || !info.IsDir() {
			t.Fatalf("canonical root %s removed: %v", rootPath, err)
		}
	}
}

func TestConsumedUntrackedRunIsDeletedWithTruthfulNonRollbackEvidence(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/runs/spent/manifest.md", "Status: consumed\n")
	writeCleanupFile(t, repositoryRoot, "do-work/runs/spent/output.txt", "spent\n")
	writeCleanupFile(t, repositoryRoot, "README.md", "tracked\n")
	runCleanupGit(t, repositoryRoot, "add", "README.md")
	runCleanupGit(t, repositoryRoot, "commit", "-q", "-m", "fixture")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("untracked consumed cleanup = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/runs/spent")); !os.IsNotExist(err) {
		t.Fatalf("spent run remains: %v", err)
	}
	foundEvidence := false
	for _, change := range result.Changes {
		if strings.Contains(change.Detail, "non-rollback spent-scratch") {
			foundEvidence = true
		}
	}
	if !foundEvidence {
		t.Fatalf("non-rollback boundary not reported: %#v", result.Changes)
	}
}

func TestConsumedUntrackedRunCommitRefusesDirtyIndexBeforeDeletion(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "unrelated.txt", "initial\n")
	commitCleanupFixture(t, repositoryRoot)

	scratchContents := map[string][]byte{
		"do-work/runs/spent/manifest.md": []byte("Status: consumed\n"),
		"do-work/runs/spent/output.txt":  []byte("spent\n"),
	}
	for relativePath, contents := range scratchContents {
		writeCleanupFile(t, repositoryRoot, relativePath, string(contents))
	}
	writeCleanupFile(t, repositoryRoot, "unrelated.txt", "staged user work\n")
	runCleanupGit(t, repositoryRoot, "add", "unrelated.txt")

	snapshot, discoveryErr := repositorymodel.DiscoverRepository(repositoryRoot)
	if discoveryErr != nil {
		t.Fatal(discoveryErr)
	}
	plan := BuildPlan(snapshot)
	result := ApplyPlan(context.Background(), plan, ApplyOptions{Commit: true, CommitMessage: "cleanup spent scratch"})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Errorf("commit result = %#v", result)
	}
	foundCommitGuard := false
	for _, finding := range result.Findings {
		if finding.Code == "CLEANUP-GROUP-REFUSED" && strings.Contains(strings.Join(finding.Evidence, " "), "--commit requires an empty existing index") {
			foundCommitGuard = true
			if !reflect.DeepEqual(finding.NextArgv, []string{"git", "diff", "--cached", "--name-only"}) ||
				!reflect.DeepEqual(finding.VerificationArgv, []string{"git", "diff", "--cached", "--quiet", "--exit-code"}) {
				t.Errorf("commit guard remediation = next %#v verify %#v", finding.NextArgv, finding.VerificationArgv)
			}
		}
	}
	if !foundCommitGuard {
		t.Errorf("commit guard refusal missing: %#v", result.Findings)
	}
	scratchRetained := true
	for relativePath, wantContents := range scratchContents {
		contents, readErr := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
		if readErr != nil || !reflect.DeepEqual(contents, wantContents) {
			t.Errorf("scratch %s = %q, %v; want byte-for-byte %q", relativePath, contents, readErr, wantContents)
			scratchRetained = false
		}
	}
	if !scratchRetained {
		return
	}

	runCleanupGit(t, repositoryRoot, "restore", "--staged", "--", "unrelated.txt")
	result = ApplyPlan(context.Background(), plan, ApplyOptions{})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("non-commit cleanup result = %#v", result)
	}
	if _, statErr := os.Stat(filepath.Join(repositoryRoot, "do-work/runs/spent")); !os.IsNotExist(statErr) {
		t.Fatalf("spent scratch remained eligible without --commit: %v", statErr)
	}
}

func TestCommittedRiskPreservesExactRevertArgvAndDoesNotClaimApplied(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-211-done.md", cleanupRequest("REQ-211", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{Commit: true, CommitMessage: "cleanup", PostCommitVerify: func(context.Context, string) error {
		return errors.New("forced post-commit risk")
	}})
	if result.Outcome != resultmodel.OutcomeRisk {
		t.Fatalf("outcome = %s", result.Outcome)
	}
	if len(result.Findings) == 0 || len(result.Findings[len(result.Findings)-1].NextArgv) != 3 || result.Findings[len(result.Findings)-1].NextArgv[0] != "git" || result.Findings[len(result.Findings)-1].NextArgv[1] != "revert" {
		t.Fatalf("revert evidence = %#v", result.Findings)
	}
	for _, change := range result.Changes {
		if strings.HasPrefix(change.Detail, "applied ") {
			t.Fatalf("risk claimed applied: %#v", result.Changes)
		}
	}
}
