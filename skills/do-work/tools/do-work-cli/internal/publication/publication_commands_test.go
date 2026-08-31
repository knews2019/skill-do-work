package publication

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestHandlersRegisterEveryPublicationCommand(t *testing.T) {
	handlers := Handlers()
	for _, commandName := range []string{"capture-files", "answer", "release"} {
		handler, found := handlers[commandName]
		if !found {
			t.Fatalf("missing handler %q", commandName)
		}
		result := handler(commandruntime.ExecutionContext{RepositoryRoot: t.TempDir(), Format: resultmodel.FormatJSON}, nil)
		if result.Findings[0].Code == "UNKNOWN-COMMAND" {
			t.Fatalf("%s still routes as unknown", commandName)
		}
	}
}

func TestApplyPlanDryRunAndCommitExposeSameExactTargets(t *testing.T) {
	repositoryRoot := t.TempDir()
	runGitFixture(t, repositoryRoot, "init", "-q")
	runGitFixture(t, repositoryRoot, "config", "user.name", "Test")
	runGitFixture(t, repositoryRoot, "config", "user.email", "test@example.com")
	writeFixture(t, repositoryRoot, "seed", []byte("seed\n"), 0o644)
	runGitFixture(t, repositoryRoot, "add", "seed")
	runGitFixture(t, repositoryRoot, "commit", "-qm", "seed")
	plan := PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: repositoryRoot, CommitMessage: "capture exact targets", Mutations: []PlannedMutation{
		{Kind: MutationCreate, Path: "do-work/queue/REQ-1-test.md", Contents: []byte("request\n"), Mode: 0o750},
		{Kind: MutationCreate, Path: "do-work/.req-reservations/REQ-1", Contents: []byte("REQ-1\n"), Mode: 0o644},
	}}
	plan = finalizePlan(plan)
	plan.CreatedDirectoryPaths, _ = planCreatedDirectories(repositoryRoot, plan.TargetPaths)
	dryRun := ApplyPlan(t.Context(), plan, true, false)
	if _, err := os.Lstat(filepath.Join(repositoryRoot, "do-work")); !os.IsNotExist(err) {
		t.Fatalf("dry-run changed tree: %v", err)
	}
	applied := ApplyPlan(t.Context(), plan, false, true)
	if dryRun.Outcome != resultmodel.OutcomeSuccess || applied.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("dry=%#v applied=%#v", dryRun, applied)
	}
	dryPaths, appliedPaths := changePaths(dryRun.Changes), changePaths(applied.Changes)
	if !reflect.DeepEqual(dryPaths, appliedPaths) || !reflect.DeepEqual(appliedPaths, plan.TargetPaths) {
		t.Fatalf("dry=%v applied=%v targets=%v", dryPaths, appliedPaths, plan.TargetPaths)
	}
	command := exec.Command("git", "-C", repositoryRoot, "diff-tree", "--root", "--no-commit-id", "--name-only", "-r", "HEAD")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	committed := strings.Fields(string(output))
	sort.Strings(committed)
	if !reflect.DeepEqual(committed, plan.TargetPaths) {
		t.Fatalf("committed=%v targets=%v", committed, plan.TargetPaths)
	}
	info, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-1-test.md"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o750 {
		t.Fatalf("mode=%v", info.Mode())
	}
}

func changePaths(changes []resultmodel.RecordedChange) []string {
	paths := make([]string, len(changes))
	for index := range changes {
		paths[index] = changes[index].Path
	}
	sort.Strings(paths)
	return paths
}

func TestApplyCapturePlanRollsBackEveryCreatedFileAndDirectory(t *testing.T) {
	repositoryRoot := t.TempDir()
	runGitFixture(t, repositoryRoot, "init", "-q")
	writeFixture(t, repositoryRoot, "seed", []byte("seed\n"), 0o644)
	runGitFixture(t, repositoryRoot, "add", "seed")
	runGitFixture(t, repositoryRoot, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "seed")
	plan := PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: repositoryRoot, Mutations: []PlannedMutation{
		{Kind: MutationCreate, Path: "do-work/.req-reservations/REQ-1", Contents: []byte("REQ-1\n"), Mode: 0o644},
		{Kind: MutationCreate, Path: "do-work/user-requests/UR-1/input.md", Contents: []byte("input"), Mode: 0o644},
		{Kind: MutationCreate, Path: "do-work/queue/REQ-1-test.md", Contents: []byte("request"), Mode: 0o644},
	}}
	plan = finalizePlan(plan)
	plan.CreatedDirectoryPaths, _ = planCreatedDirectories(repositoryRoot, plan.TargetPaths)
	previous := beforePublicationMutation
	beforePublicationMutation = func(index int, _ PlannedMutation) error {
		if index == 2 {
			return errors.New("injected publication failure")
		}
		return nil
	}
	t.Cleanup(func() { beforePublicationMutation = previous })
	result := ApplyPlan(t.Context(), plan, false, false)
	if result.Outcome != resultmodel.OutcomeRolledBack || result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Lstat(filepath.Join(repositoryRoot, "do-work")); !os.IsNotExist(err) {
		t.Fatalf("do-work survived rollback: %v", err)
	}
}

func runGitFixture(t *testing.T, repositoryRoot string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
}
