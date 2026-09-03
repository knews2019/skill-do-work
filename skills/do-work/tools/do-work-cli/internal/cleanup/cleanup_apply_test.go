package cleanup

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
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

func TestCleanupPreflightRemediationMatchesEveryReachableFailureKind(t *testing.T) {
	tests := []struct {
		name              string
		groupCode         string
		affectedID        string
		setup             func(*testing.T) (CleanupPlan, ApplyOptions)
		wantFailureCode   string
		wantReason        string
		wantAffectedPaths []string
		wantNextArgv      []string
		wantVerifyArgv    []string
	}{
		{
			name:       "not git repository",
			groupCode:  "NOT-GIT-GROUP",
			affectedID: "REQ-501",
			setup: func(t *testing.T) (CleanupPlan, ApplyOptions) {
				repositoryRoot := t.TempDir()
				return CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{{Code: "NOT-GIT-GROUP", AffectedID: "REQ-501", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "target.txt", Contents: []byte("replacement\n")}}}}}, ApplyOptions{}
			},
			wantFailureCode: "GIT-NOT-GIT-REPOSITORY",
			wantReason:      "mutating commands require a Git repository",
			wantNextArgv:    []string{"do-work-cli", "--repo-root", "<git-repository>", "cleanup"},
			wantVerifyArgv:  []string{"git", "rev-parse", "--show-toplevel"},
		},
		{
			name:       "invalid empty targets",
			groupCode:  "INVALID-GROUP",
			affectedID: "REQ-502",
			setup: func(t *testing.T) (CleanupPlan, ApplyOptions) {
				repositoryRoot := cleanupRepository(t)
				return CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{{Code: "INVALID-GROUP", AffectedID: "REQ-502"}}}, ApplyOptions{}
			},
			wantFailureCode: "GIT-INVALID-OPTIONS",
			wantReason:      "at least one exact target path is required",
			wantNextArgv:    []string{"do-work-cli", "--format", "text", "cleanup"},
			wantVerifyArgv:  []string{"do-work-cli", "--format", "json", "cleanup"},
		},
		{
			name:       "dirty target",
			groupCode:  "DIRTY-TARGET-GROUP",
			affectedID: "REQ-503",
			setup: func(t *testing.T) (CleanupPlan, ApplyOptions) {
				repositoryRoot := cleanupRepository(t)
				writeCleanupFile(t, repositoryRoot, "target.txt", "original\n")
				commitCleanupFixture(t, repositoryRoot)
				writeCleanupFile(t, repositoryRoot, "target.txt", "user edit\n")
				return CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{{Code: "DIRTY-TARGET-GROUP", AffectedID: "REQ-503", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "target.txt", Contents: []byte("replacement\n")}}}}}, ApplyOptions{}
			},
			wantFailureCode:   "GIT-DIRTY-TARGET",
			wantReason:        "already dirty",
			wantAffectedPaths: []string{"target.txt"},
			wantNextArgv:      []string{"git", "status", "--short", "--", "target.txt"},
			wantVerifyArgv:    []string{"git", "diff", "--quiet", "--exit-code", "--", "target.txt"},
		},
		{
			name:       "dirty index",
			groupCode:  "DIRTY-INDEX-GROUP",
			affectedID: "REQ-504",
			setup: func(t *testing.T) (CleanupPlan, ApplyOptions) {
				repositoryRoot := cleanupRepository(t)
				writeCleanupFile(t, repositoryRoot, "target.txt", "original\n")
				writeCleanupFile(t, repositoryRoot, "unrelated.txt", "original\n")
				commitCleanupFixture(t, repositoryRoot)
				writeCleanupFile(t, repositoryRoot, "unrelated.txt", "staged\n")
				runCleanupGit(t, repositoryRoot, "add", "unrelated.txt")
				return CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{{Code: "DIRTY-INDEX-GROUP", AffectedID: "REQ-504", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "target.txt", Contents: []byte("replacement\n")}}}}}, ApplyOptions{Commit: true, CommitMessage: "cleanup"}
			},
			wantFailureCode: "GIT-DIRTY-INDEX",
			wantReason:      "--commit requires an empty existing index",
			wantNextArgv:    []string{"git", "diff", "--cached", "--name-only"},
			wantVerifyArgv:  []string{"git", "diff", "--cached", "--quiet", "--exit-code"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			plan, options := test.setup(t)
			result := ApplyPlan(context.Background(), plan, options)
			if result.Outcome != resultmodel.OutcomeFindings || len(result.Findings) != 1 || len(result.Changes) != 0 {
				t.Fatalf("preflight result = %#v", result)
			}
			finding := result.Findings[0]
			evidence := strings.Join(finding.Evidence, " ")
			if finding.Code != "CLEANUP-GROUP-REFUSED" ||
				!reflect.DeepEqual(finding.AffectedIDs, []string{test.affectedID}) ||
				!reflect.DeepEqual(finding.AffectedPaths, test.wantAffectedPaths) ||
				!strings.Contains(evidence, test.groupCode) ||
				!strings.Contains(evidence, test.wantFailureCode) ||
				!strings.Contains(evidence, test.wantReason) ||
				!reflect.DeepEqual(finding.NextArgv, test.wantNextArgv) ||
				!reflect.DeepEqual(finding.VerificationArgv, test.wantVerifyArgv) {
				t.Fatalf("preflight finding = %#v", finding)
			}
		})
	}
}

func TestApplyPlanRefusesNestedRepositoryRootBeforeMutation(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	sourcePath := "do-work/queue/REQ-438-done.md"
	destinationPath := "do-work/archive/REQ-438-done.md"
	writeCleanupFile(t, repositoryRoot, sourcePath, cleanupRequest("REQ-438", "done", ""))
	commitCleanupFixture(t, repositoryRoot)
	nestedRoot := filepath.Join(repositoryRoot, "nested")
	if err := os.MkdirAll(nestedRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	headBefore := runCleanupGit(t, repositoryRoot, "rev-parse", "HEAD")

	plan := CleanupPlan{RepositoryRoot: nestedRoot, Groups: []OperationGroup{{
		Code:       "ARCHIVE-REQ-438",
		AffectedID: "REQ-438",
		Operations: []CleanupOperation{{
			Kind:            OperationMove,
			SourcePath:      sourcePath,
			DestinationPath: destinationPath,
		}},
	}}}
	result := ApplyPlan(context.Background(), plan, ApplyOptions{Commit: true, CommitMessage: "nested cleanup"})
	if result.Outcome != resultmodel.OutcomeFindings || len(result.Findings) != 1 || len(result.Changes) != 0 {
		t.Fatalf("nested-root cleanup = %#v", result)
	}
	evidence := strings.Join(result.Findings[0].Evidence, " ")
	if !strings.Contains(evidence, "GIT-INVALID-OPTIONS") || !strings.Contains(evidence, nestedRoot) || !strings.Contains(evidence, repositoryRoot) {
		t.Fatalf("nested-root cleanup evidence = %#v", result.Findings[0])
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(sourcePath))); err != nil {
		t.Fatalf("nested-root refusal changed source: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(destinationPath))); !os.IsNotExist(err) {
		t.Fatalf("nested-root refusal created destination: %v", err)
	}
	if headAfter := runCleanupGit(t, repositoryRoot, "rev-parse", "HEAD"); headAfter != headBefore {
		t.Fatalf("nested-root refusal changed HEAD: before %s after %s", headBefore, headAfter)
	}
	if status := runCleanupGit(t, repositoryRoot, "status", "--short"); status != "" {
		t.Fatalf("nested-root refusal changed repository state: %q", status)
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
	for _, path := range []string{"middle.txt", "transitive.txt", "missing.txt", "empty-required.txt", "duplicate-one.txt", "duplicate-two.txt", "duplicate-dependent.txt", "cycle-a.txt", "cycle-b.txt", "repeated.txt", "collision-source.txt", "archive/collision.txt"} {
		writeCleanupFile(t, repositoryRoot, path, "clean\n")
	}
	commitCleanupFixture(t, repositoryRoot)
	writeCleanupFile(t, repositoryRoot, "direct.txt", "user edit\n")

	plan := CleanupPlan{RepositoryRoot: repositoryRoot, Groups: []OperationGroup{
		{Code: "DIRECT", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "direct.txt", Contents: []byte("replacement\n")}}},
		{Code: "MIDDLE", RequiredGroupCodes: []string{"DIRECT"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "middle.txt", Contents: []byte("clean\n")}}},
		{Code: "TRANSITIVE", RequiredGroupCodes: []string{"MIDDLE"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "transitive.txt", Contents: []byte("clean\n")}}},
		{Code: "MISSING", RequiredGroupCodes: []string{"ABSENT"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "missing.txt", Contents: []byte("clean\n")}}},
		{Code: "EMPTY-REQUIRED", RequiredGroupCodes: []string{""}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "empty-required.txt", Contents: []byte("clean\n")}}},
		{Code: "DUPLICATE", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "duplicate-one.txt", Contents: []byte("clean\n")}}},
		{Code: "DUPLICATE", Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "duplicate-two.txt", Contents: []byte("clean\n")}}},
		{Code: "DUPLICATE-DEPENDENT", RequiredGroupCodes: []string{"DUPLICATE"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "duplicate-dependent.txt", Contents: []byte("clean\n")}}},
		{Code: "CYCLE-A", RequiredGroupCodes: []string{"CYCLE-B"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "cycle-a.txt", Contents: []byte("clean\n")}}},
		{Code: "CYCLE-B", RequiredGroupCodes: []string{"CYCLE-A"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "cycle-b.txt", Contents: []byte("clean\n")}}},
		{Code: "REPEATED", RequiredGroupCodes: []string{"SAFE", "SAFE"}, Operations: []CleanupOperation{{Kind: OperationReplace, SourcePath: "repeated.txt", Contents: []byte("clean\n")}}},
		{Code: "COLLISION", Operations: []CleanupOperation{{Kind: OperationMove, SourcePath: "collision-source.txt", DestinationPath: "archive/collision.txt"}}},
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
		"EMPTY-REQUIRED":      "<empty>",
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

	structuralRefusalCount := 0
	duplicateRefusalCount := 0
	for _, finding := range result.Findings {
		if finding.Code != "CLEANUP-GROUP-REFUSED" {
			continue
		}
		evidence := strings.Join(finding.Evidence, " ")
		isDuplicateIdentity := strings.Contains(evidence, "DUPLICATE: group code is duplicated")
		isPrerequisiteRefusal := strings.Contains(evidence, ": prerequisite ")
		if !isDuplicateIdentity && !isPrerequisiteRefusal {
			continue
		}
		structuralRefusalCount++
		for commandName, commandArgv := range map[string][]string{"next": finding.NextArgv, "verification": finding.VerificationArgv} {
			if len(commandArgv) == 0 {
				t.Errorf("%s structural refusal has empty %s argv: %#v", evidence, commandName, finding)
				continue
			}
			for argumentIndex, argument := range commandArgv {
				if argument == "" {
					t.Errorf("%s structural refusal %s argv[%d] is empty: %#v", evidence, commandName, argumentIndex, commandArgv)
				}
			}
		}
		if isDuplicateIdentity {
			duplicateRefusalCount++
			if !reflect.DeepEqual(finding.NextArgv, []string{"git", "status", "--short"}) {
				t.Errorf("duplicate identity next argv = %#v", finding.NextArgv)
				continue
			}
			if !reflect.DeepEqual(finding.VerificationArgv, []string{"do-work-cli", "cleanup", "--dry-run"}) {
				t.Errorf("duplicate identity verification argv = %#v", finding.VerificationArgv)
			}
			command := exec.Command(finding.NextArgv[0], finding.NextArgv[1:]...)
			command.Dir = repositoryRoot
			if output, err := command.CombinedOutput(); err != nil {
				t.Errorf("duplicate identity diagnostic failed: %v\n%s", err, output)
			}
		} else {
			wantArgv := []string{"do-work-cli", "cleanup", "--dry-run"}
			if !reflect.DeepEqual(finding.NextArgv, wantArgv) || !reflect.DeepEqual(finding.VerificationArgv, wantArgv) {
				t.Errorf("prerequisite refusal argv = next %#v verify %#v, want %#v", finding.NextArgv, finding.VerificationArgv, wantArgv)
			}
		}
	}
	if structuralRefusalCount != 10 || duplicateRefusalCount != 2 {
		t.Errorf("structural refusal counts = %d total, %d duplicate; want 10 total, 2 duplicate: %#v", structuralRefusalCount, duplicateRefusalCount, result.Findings)
	}

	for groupCode, wantCommands := range map[string]struct {
		nextArgv   []string
		verifyArgv []string
	}{
		"DIRECT": {
			nextArgv:   []string{"git", "status", "--short", "--", "direct.txt"},
			verifyArgv: []string{"git", "diff", "--quiet", "--exit-code", "--", "direct.txt"},
		},
		"COLLISION": {
			nextArgv:   []string{"git", "status", "--short", "--", "archive/collision.txt"},
			verifyArgv: []string{"do-work-cli", "cleanup", "--dry-run"},
		},
	} {
		foundExactPath := false
		for _, finding := range result.Findings {
			evidence := strings.Join(finding.Evidence, " ")
			if finding.Code == "CLEANUP-GROUP-REFUSED" && strings.HasPrefix(evidence, groupCode+": ") {
				foundExactPath = true
				if !reflect.DeepEqual(finding.AffectedPaths, wantCommands.nextArgv[4:]) || !reflect.DeepEqual(finding.NextArgv, wantCommands.nextArgv) || !reflect.DeepEqual(finding.VerificationArgv, wantCommands.verifyArgv) {
					t.Errorf("%s path-bearing refusal = paths %#v next %#v verify %#v", groupCode, finding.AffectedPaths, finding.NextArgv, finding.VerificationArgv)
				}
			}
		}
		if !foundExactPath {
			t.Errorf("%s path-bearing refusal missing: %#v", groupCode, result.Findings)
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
	if err := moveWithoutOverwrite(repositoryRoot, "source.txt", "archive/destination.txt", nil); err == nil {
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

			if err := moveWithoutOverwrite(repositoryRoot, "source.txt", "archive/destination.txt", nil); err != nil {
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

func TestConsumedUntrackedRunCommitRefusesScratchOnlyAndMixedDeletion(t *testing.T) {
	for _, test := range []struct {
		name        string
		withTracked bool
	}{
		{name: "scratch only"},
		{name: "beside eligible tracked cleanup", withTracked: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := cleanupRepository(t)
			writeCleanupFile(t, repositoryRoot, "README.md", "tracked\n")
			if test.withTracked {
				writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-206-done.md", cleanupRequest("REQ-206", "completed", ""))
			}
			commitCleanupFixture(t, repositoryRoot)
			originalHead := runCleanupGit(t, repositoryRoot, "rev-parse", "HEAD")

			scratchContents := map[string][]byte{
				"do-work/runs/spent/manifest.md": []byte("Status: consumed\n"),
				"do-work/runs/spent/output.txt":  []byte("spent\n"),
			}
			for relativePath, contents := range scratchContents {
				writeCleanupFile(t, repositoryRoot, relativePath, string(contents))
			}
			snapshot, discoveryErr := repositorymodel.DiscoverRepository(repositoryRoot)
			if discoveryErr != nil {
				t.Fatal(discoveryErr)
			}
			result := ApplyPlan(context.Background(), BuildPlan(snapshot), ApplyOptions{Commit: true, CommitMessage: "cleanup with spent scratch"})
			if result.Outcome != resultmodel.OutcomeFindings {
				t.Fatalf("commit result = %#v", result)
			}

			foundRefusal := false
			for _, finding := range result.Findings {
				if finding.Code != "CLEANUP-GROUP-REFUSED" || !strings.Contains(strings.Join(finding.Evidence, " "), "--commit cannot delete entirely untracked consumed scratch") {
					continue
				}
				foundRefusal = true
				wantPaths := []string{"do-work/runs/spent/manifest.md", "do-work/runs/spent/output.txt"}
				if !reflect.DeepEqual(finding.AffectedPaths, wantPaths) {
					t.Errorf("refusal paths = %#v, want %#v", finding.AffectedPaths, wantPaths)
				}
				if !reflect.DeepEqual(finding.NextArgv, []string{"do-work-cli", "cleanup"}) {
					t.Errorf("refusal remediation = %#v", finding.NextArgv)
				}
			}
			if !foundRefusal {
				t.Fatalf("commit-mode scratch refusal missing: %#v", result.Findings)
			}
			for relativePath, wantContents := range scratchContents {
				contents, readErr := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
				if readErr != nil || !reflect.DeepEqual(contents, wantContents) {
					t.Errorf("scratch %s = %q, %v; want byte-for-byte %q", relativePath, contents, readErr, wantContents)
				}
			}

			currentHead := runCleanupGit(t, repositoryRoot, "rev-parse", "HEAD")
			if test.withTracked {
				if currentHead == originalHead {
					t.Fatal("eligible tracked cleanup was not committed")
				}
				if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/archive/REQ-206-done.md")); err != nil {
					t.Fatalf("eligible tracked cleanup did not apply: %v", err)
				}
			} else if currentHead != originalHead {
				t.Fatalf("scratch-only refusal changed HEAD from %s to %s", originalHead, currentHead)
			}
		})
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

// Two writers can race the same absent destination. The move used to register the
// destination before creating it, so the loser's EEXIST rolled back by deleting the
// winner's published file.
func TestLosingMoveWriterRollbackPreservesTheWinnersDestination(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-601-done.md", "ours\n")
	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	commitCleanupFixture(t, repositoryRoot)
	destinationPath := filepath.Join(repositoryRoot, "do-work", "archive", "REQ-601-done.md")

	result := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"do-work/queue/REQ-601-done.md", "do-work/archive/REQ-601-done.md"},
	}, func(recorder *gittransaction.MutationRecorder) error {
		// The winner publishes the destination after this plan was built and preflighted.
		if err := os.WriteFile(destinationPath, []byte("winner\n"), 0o644); err != nil {
			return err
		}
		return applyOperation(repositoryRoot, recorder, CleanupOperation{
			Kind:            OperationMove,
			SourcePath:      "do-work/queue/REQ-601-done.md",
			DestinationPath: "do-work/archive/REQ-601-done.md",
		})
	})

	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("the losing writer reported success: %#v", result)
	}
	if contents, err := os.ReadFile(destinationPath); err != nil || string(contents) != "winner\n" {
		t.Fatalf("losing rollback destroyed the winner's file: %q %v", contents, err)
	}
	if contents, err := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "queue", "REQ-601-done.md")); err != nil || string(contents) != "ours\n" {
		t.Fatalf("losing writer lost its own source: %q %v", contents, err)
	}
}

// Registration now happens between the exclusive create and the source deletion, which adds
// one new failure point: a registration that fails must undo exactly the object this process
// created and leave the source untouched.
func TestFailedDestinationRegistrationRemovesOnlyTheCreatedDestination(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-602-done.md", "ours\n")
	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	sourcePath := filepath.Join(repositoryRoot, "do-work", "queue", "REQ-602-done.md")
	destinationPath := filepath.Join(repositoryRoot, "do-work", "archive", "REQ-602-done.md")

	registrationError := errors.New("registration refused")
	registered := false
	moveError := moveWithoutOverwrite(repositoryRoot, "do-work/queue/REQ-602-done.md", "do-work/archive/REQ-602-done.md", func() error {
		registered = true
		// The exclusive create is the ownership event, so the object must already be
		// published when registration runs.
		if contents, err := os.ReadFile(destinationPath); err != nil || string(contents) != "ours\n" {
			t.Fatalf("destination was registered before it was created: %q %v", contents, err)
		}
		if _, err := os.Lstat(sourcePath); err != nil {
			t.Fatalf("source deleted before registration: %v", err)
		}
		return registrationError
	})

	if !registered || !errors.Is(moveError, registrationError) {
		t.Fatalf("registered=%v move error=%v", registered, moveError)
	}
	if _, err := os.Lstat(destinationPath); !os.IsNotExist(err) {
		t.Fatalf("unregistered destination remains: %v", err)
	}
	if contents, err := os.ReadFile(sourcePath); err != nil || string(contents) != "ours\n" {
		t.Fatalf("source did not survive a failed registration: %q %v", contents, err)
	}

	// The same failure through the real recorder seam: an undeclared destination cannot be
	// registered, so the transaction fails with the source intact and nothing published.
	result := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"do-work/queue/REQ-602-done.md"},
	}, func(recorder *gittransaction.MutationRecorder) error {
		return applyOperation(repositoryRoot, recorder, CleanupOperation{
			Kind:            OperationMove,
			SourcePath:      "do-work/queue/REQ-602-done.md",
			DestinationPath: "do-work/archive/REQ-602-done.md",
		})
	})
	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("an unregistrable destination reported success: %#v", result)
	}
	if _, err := os.Lstat(destinationPath); !os.IsNotExist(err) {
		t.Fatalf("unregistrable destination remains: %v", err)
	}
	if contents, err := os.ReadFile(sourcePath); err != nil || string(contents) != "ours\n" {
		t.Fatalf("source lost to an unregistrable destination: %q %v", contents, err)
	}
}
