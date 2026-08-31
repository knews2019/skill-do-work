package publication

import (
	"encoding/json"
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

func TestRemediationF1RootedParentsRefuseCreateAndMoveSwapsOutsideRepository(t *testing.T) {
	for _, test := range []struct {
		name       string
		mutation   PlannedMutation
		swapParent string
		outside    string
	}{
		{name: "create destination", mutation: PlannedMutation{Kind: MutationCreate, Path: "do-work/queue/REQ-1.md", Contents: []byte("new\n"), Mode: 0o644}, swapParent: "do-work/queue", outside: "REQ-1.md"},
		{name: "move destination", mutation: PlannedMutation{Kind: MutationMove, Path: "do-work/working/REQ-1.md", DestinationPath: "do-work/archive/REQ-1.md", ExpectedBytes: []byte("old\n")}, swapParent: "do-work/archive", outside: "REQ-1.md"},
	} {
		t.Run(test.name, func(t *testing.T) {
			repositoryRoot := t.TempDir()
			outsideRoot := t.TempDir()
			outsidePath := filepath.Join(outsideRoot, test.outside)
			if err := os.WriteFile(outsidePath, []byte("protected outside\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			runGitFixture(t, repositoryRoot, "init", "-q")
			writeFixture(t, repositoryRoot, "do-work/working/REQ-1.md", []byte("old\n"), 0o644)
			if err := os.MkdirAll(filepath.Join(repositoryRoot, filepath.FromSlash(test.swapParent)), 0o755); err != nil {
				t.Fatal(err)
			}
			writeFixture(t, repositoryRoot, "seed", []byte("seed\n"), 0o644)
			runGitFixture(t, repositoryRoot, "add", ".")
			runGitFixture(t, repositoryRoot, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "seed")
			plan := finalizePlan(PublicationPlan{Operation: OperationAnswer, RepositoryRoot: repositoryRoot, Mutations: []PlannedMutation{test.mutation}})
			previous := beforePublicationMutation
			beforePublicationMutation = func(_ int, _ PlannedMutation) error {
				original := filepath.Join(repositoryRoot, filepath.FromSlash(test.swapParent))
				if err := os.Rename(original, original+"-original"); err != nil {
					return err
				}
				return os.Symlink(outsideRoot, original)
			}
			t.Cleanup(func() { beforePublicationMutation = previous })
			result := ApplyPlan(t.Context(), plan, false, false)
			if result.Outcome == resultmodel.OutcomeSuccess {
				t.Fatalf("parent swap reported success: %#v", result)
			}
			outsideBytes, err := os.ReadFile(outsidePath)
			if err != nil || string(outsideBytes) != "protected outside\n" {
				t.Fatalf("outside destination was touched: %v %q", err, outsideBytes)
			}
		})
	}
}

func TestRemediationF8FindingsCarryExactManifestProtocol(t *testing.T) {
	repositoryRoot := t.TempDir()
	manifestPath := "payload/exact answer manifest.json"
	writeFixture(t, repositoryRoot, manifestPath, []byte(`{"operation":"answer","answer":{"request_path":"missing","expected_status":"blocked","mode":"stakeholder","answers":[{"expected_line":"- [ ] Q?","outcome":"answered","summary":"yes"}]}}`), 0o644)
	result := handlePublicationCommand(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatJSON}, OperationAnswer, []string{"--manifest", manifestPath, "--dry-run"})
	if len(result.Findings) == 0 {
		t.Fatal("missing finding")
	}
	finding := result.Findings[0]
	for _, argv := range [][]string{finding.NextArgv, finding.VerificationArgv} {
		joined := strings.Join(argv, "\x00")
		if !strings.Contains(joined, manifestPath) || strings.Contains(joined, "<manifest.json>") {
			t.Fatalf("argv is not exact: %#v", argv)
		}
	}
	if finding.NextJustRecipe == "" {
		t.Fatal("publication finding has no actionable recipe equivalent")
	}
}

func TestRemediationF9TextAndJSONRenderTheSamePublicationFinding(t *testing.T) {
	result := commandFailure(t.TempDir(), OperationRelease, "RELEASE-TEST", "evidence")
	if result.Findings[0].NextJustRecipe == "" {
		t.Fatal("matrix requires an actionable recipe as well as renderer parity")
	}
	jsonBytes, err := resultmodel.RenderResult(result, resultmodel.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	textBytes, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{decoded.Findings[0].Code, decoded.Findings[0].NextArgv[len(decoded.Findings[0].NextArgv)-1], decoded.Findings[0].VerificationArgv[len(decoded.Findings[0].VerificationArgv)-2]} {
		if !strings.Contains(string(textBytes), value) {
			t.Fatalf("text renderer omitted %q: %s", value, textBytes)
		}
	}
}

func TestRemediationF9ReplaceAndMoveRollbackRestoreBytesAndModes(t *testing.T) {
	for _, test := range []struct {
		name      string
		mutation  PlannedMutation
		source    string
		wantBytes []byte
		wantMode  os.FileMode
	}{
		{name: "replace", mutation: PlannedMutation{Kind: MutationReplace, Path: "tracked.txt", ExpectedBytes: []byte("old\n"), Contents: []byte("new\n")}, source: "tracked.txt", wantBytes: []byte("old\n"), wantMode: 0o755},
		{name: "move", mutation: PlannedMutation{Kind: MutationMove, Path: "source.txt", DestinationPath: "archive/source.txt", ExpectedBytes: []byte("move-old\n")}, source: "source.txt", wantBytes: []byte("move-old\n"), wantMode: 0o755},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			runGitFixture(t, root, "init", "-q")
			writeFixture(t, root, test.source, test.wantBytes, test.wantMode)
			writeFixture(t, root, "seed", []byte("seed\n"), 0o644)
			if err := os.MkdirAll(filepath.Join(root, "archive"), 0o755); err != nil {
				t.Fatal(err)
			}
			runGitFixture(t, root, "add", ".")
			runGitFixture(t, root, "-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "seed")
			plan := finalizePlan(PublicationPlan{Operation: OperationRelease, RepositoryRoot: root, Mutations: []PlannedMutation{test.mutation, PlannedMutation{Kind: MutationCreate, Path: "z-last", Contents: []byte("last"), Mode: 0o644}}})
			previous := beforePublicationMutation
			beforePublicationMutation = func(index int, _ PlannedMutation) error {
				if index == 1 {
					return errors.New("rollback seam")
				}
				return nil
			}
			t.Cleanup(func() { beforePublicationMutation = previous })
			result := ApplyPlan(t.Context(), plan, false, false)
			if result.Outcome != resultmodel.OutcomeRolledBack || result.Rollback.Status != resultmodel.RollbackSucceeded {
				t.Fatalf("result = %#v", result)
			}
			got, err := os.ReadFile(filepath.Join(root, test.source))
			if err != nil || !reflect.DeepEqual(got, test.wantBytes) {
				t.Fatalf("restored bytes = %q, err=%v", got, err)
			}
			info, err := os.Stat(filepath.Join(root, test.source))
			if err != nil || info.Mode().Perm() != test.wantMode {
				t.Fatalf("restored mode = %v, err=%v", info, err)
			}
		})
	}
}

func TestRemediationF9CommitGuardsAndPostCommitRiskExposeExactRevert(t *testing.T) {
	t.Run("dirty index refused", func(t *testing.T) {
		root := initializedGitRepository(t)
		writeFixture(t, root, "staged.txt", []byte("staged\n"), 0o644)
		runGitFixture(t, root, "add", "staged.txt")
		plan := finalizePlan(PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: root, CommitMessage: "publication", Mutations: []PlannedMutation{{Kind: MutationCreate, Path: "target.txt", Contents: []byte("target\n"), Mode: 0o644}}})
		result := ApplyPlan(t.Context(), plan, false, true)
		if result.Outcome != resultmodel.OutcomeRefused {
			t.Fatalf("dirty index result = %#v", result)
		}
	})

	t.Run("post commit verification risk", func(t *testing.T) {
		root := initializedGitRepository(t)
		plan := finalizePlan(PublicationPlan{Operation: OperationCaptureFiles, RepositoryRoot: root, ManifestPath: "payload/manifest.json", CommitMessage: "publication", Mutations: []PlannedMutation{{Kind: MutationCreate, Path: "target.txt", Contents: []byte("target\n"), Mode: 0o644}}})
		previous := afterPublicationCommit
		afterPublicationCommit = func(PublicationPlan) error { return errors.New("post-commit risk seam") }
		t.Cleanup(func() { afterPublicationCommit = previous })
		result := ApplyPlan(t.Context(), plan, false, true)
		if result.Outcome != resultmodel.OutcomeRisk || len(result.Findings) != 1 || len(result.Findings[0].NextArgv) != 3 || result.Findings[0].NextArgv[0] != "git" || result.Findings[0].NextArgv[1] != "revert" || result.Findings[0].NextArgv[2] == "" {
			t.Fatalf("post-commit result = %#v", result)
		}
	})
}

func initializedGitRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGitFixture(t, root, "init", "-q")
	runGitFixture(t, root, "config", "user.name", "Test")
	runGitFixture(t, root, "config", "user.email", "test@example.com")
	writeFixture(t, root, "seed", []byte("seed\n"), 0o644)
	runGitFixture(t, root, "add", "seed")
	runGitFixture(t, root, "commit", "-qm", "seed")
	return root
}

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

func TestRemediationF9ApplyPlanDryRunAndCommitExposeSameExactTargets(t *testing.T) {
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
