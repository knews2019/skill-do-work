package suiteinstall

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const fixtureSuiteVersion = "0.200.0"

// justfileTemplateBytes mirrors the shipped board template's shape: the WHOLE file is one
// managed section, which is what lets a fresh install's justfile equal the template byte for
// byte and lets the same file serve as both section and template.
const justfileTemplateBytes = `# >>> do-work:recipes >>>
run-kanban:
    echo kanban

run-kanban-cli:
    echo kanban-cli

kanban-static:
    echo static

kanban-summary:
    echo summary

run-do-work-update:
    echo update
# <<< do-work:recipes <<<
`

const agentInstructionsTemplateBytes = `<!-- >>> do-work:communication-style >>> -->
See .claude/skills/do-work/crew-members/communication-style.md
<!-- <<< do-work:communication-style <<< -->
`

const coreHooksFragmentBytes = `{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash \"${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/session-start.sh\""
          }
        ]
      }
    ]
  }
}
`

var fixtureModuleNames = []string{"do-work", "do-work-board", "do-work-knowledge", "do-work-toolbox"}

// newSuiteSourceTree builds a suite an installer accepts. It is synthesised rather than
// copied so a test can vary one thing about it without dragging the whole repository in.
func newSuiteSourceTree(t *testing.T, suiteVersion string) string {
	t.Helper()
	sourceRoot := t.TempDir()
	writeTestFile(t, filepath.Join(sourceRoot, "VERSION"), suiteVersion+"\n")
	manifest := "source\tdestination\n"
	for _, module := range fixtureModuleNames {
		manifest += "skills/" + module + "\t.claude/skills/" + module + "\n"
		writeTestFile(t, filepath.Join(sourceRoot, "skills", module, "SKILL.md"), "# "+module+"\n")
		writeTestFile(t, filepath.Join(sourceRoot, "skills", module, "payload.txt"), "payload for "+module+"\n")
	}
	writeTestFile(t, filepath.Join(sourceRoot, "suite", "modules.tsv"), manifest)
	writeTestFile(t, filepath.Join(sourceRoot, "skills", "do-work", "VERSION"), suiteVersion+"\n")
	writeTestFile(t, filepath.Join(sourceRoot, "skills", "do-work", "actions", "version.md"),
		"# Version Action\n\n**Current version**: "+suiteVersion+"\n")
	writeTestFile(t, filepath.Join(sourceRoot, "skills", "do-work-board", "justfile.template"), justfileTemplateBytes)
	writeTestFile(t, filepath.Join(sourceRoot, "skills", "do-work", "hooks", "hooks.json"), coreHooksFragmentBytes)
	writeTestFile(t, filepath.Join(sourceRoot, "skills", "do-work", "agent-instructions.template.md"), agentInstructionsTemplateBytes)
	return sourceRoot
}

func newProjectRepository(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "gitconfig"))
	t.Setenv("GIT_CONFIG_SYSTEM", filepath.Join(t.TempDir(), "gitconfig-system"))
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
	runTestGit(t, projectRoot, "init", "-q", ".")
	runTestGit(t, projectRoot, "config", "user.email", "fixture@example.com")
	runTestGit(t, projectRoot, "config", "user.name", "Suite Install Fixture")
	// t.TempDir can sit under a symlinked /tmp; the transaction resolves physically, so the
	// fixture has to hand back the same physical path it will report.
	physical, err := filepath.EvalSymlinks(projectRoot)
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}
	return physical
}

func runInstallFixture(t *testing.T, projectRoot, sourceRoot, confirmation string) (InstallResult, string) {
	t.Helper()
	var narration bytes.Buffer
	result := RunInstall(context.Background(), InstallOptions{
		ProjectRoot:         projectRoot,
		ExtractedSourceRoot: sourceRoot,
		Narration:           &narration,
		ConfirmationInput:   strings.NewReader(confirmation),
	})
	return result, narration.String()
}

func TestFreshInstallWritesFourModulesAndEveryManagedConfiguration(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)

	result, narration := runInstallFixture(t, projectRoot, sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, reason = %q\n%s", result.Outcome, result.FailureReason, narration)
	}
	for _, module := range fixtureModuleNames {
		skillFile := filepath.Join(projectRoot, ".claude", "skills", module, "SKILL.md")
		if info, err := os.Stat(skillFile); err != nil || info.Size() == 0 {
			t.Errorf("module %s was not installed: %v", module, err)
		}
	}
	// A fresh justfile is the board template byte for byte; synthesising those bytes instead
	// of copying them would break the shipped cmp.
	installedJustfile := readTestFile(t, filepath.Join(projectRoot, "justfile"))
	if installedJustfile != justfileTemplateBytes {
		t.Errorf("fresh justfile is not the template byte for byte:\n%q", installedJustfile)
	}
	if instructions := readTestFile(t, filepath.Join(projectRoot, "CLAUDE.md")); instructions != agentInstructionsTemplateBytes {
		t.Errorf("fresh CLAUDE.md is not the template byte for byte:\n%q", instructions)
	}
	settings := readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json"))
	if !strings.Contains(settings, "do-work/hooks/session-start.sh") {
		t.Errorf("settings did not receive the core SessionStart hook:\n%s", settings)
	}
	if len(result.Changes) != 7 {
		t.Errorf("changes = %#v, want four modules plus three configuration files", result.Changes)
	}
	if !strings.Contains(narration, "Installed do-work suite v"+fixtureSuiteVersion) {
		t.Errorf("narration does not report the installed version:\n%s", narration)
	}
}

func TestReinstallPreservesCustomBytesAndModesAndIsByteIdempotent(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)

	writeTestFile(t, filepath.Join(projectRoot, "justfile"),
		"custom-before:\n    echo before\n\n# >>> do-work:recipes >>>\nold-managed:\n    echo old\n# <<< do-work:recipes <<<\n\ncustom-after:\n    echo after\n")
	chmodTestFile(t, filepath.Join(projectRoot, "justfile"), 0o640)
	writeTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json"),
		`{"custom":{"keep":[1,2,3]},"hooks":{"Stop":[{"hooks":[{"type":"command","command":"bash \".claude/skills/do-work/hooks/pipeline-guard.sh\""},{"type":"command","command":"echo custom-stop"}]}]}}`+"\n")
	chmodTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json"), 0o600)
	writeTestFile(t, filepath.Join(projectRoot, "CLAUDE.md"),
		"# Consumer\n\nBefore.\n\n<!-- >>> do-work:communication-style >>> -->\nstale\n<!-- <<< do-work:communication-style <<< -->\n\nAfter.\n")

	result, narration := runInstallFixture(t, projectRoot, sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, reason = %q\n%s", result.Outcome, result.FailureReason, narration)
	}
	justfile := readTestFile(t, filepath.Join(projectRoot, "justfile"))
	for _, preserved := range []string{"custom-before:", "custom-after:"} {
		if !strings.Contains(justfile, preserved) {
			t.Errorf("reinstall dropped %q from the justfile:\n%s", preserved, justfile)
		}
	}
	if strings.Contains(justfile, "old-managed:") {
		t.Errorf("reinstall kept stale content inside the managed section:\n%s", justfile)
	}
	if mode := modeOf(t, filepath.Join(projectRoot, "justfile")); mode != 0o640 {
		t.Errorf("justfile mode = %o, want 640", mode)
	}
	settings := readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json"))
	if strings.Contains(settings, "pipeline-guard.sh") {
		t.Errorf("reinstall kept the retired pipeline guard:\n%s", settings)
	}
	if !strings.Contains(settings, "echo custom-stop") || !strings.Contains(settings, `"keep"`) {
		t.Errorf("reinstall dropped consumer settings state:\n%s", settings)
	}
	if mode := modeOf(t, filepath.Join(projectRoot, ".claude", "settings.json")); mode != 0o600 {
		t.Errorf("settings mode = %o, want 600", mode)
	}
	instructions := readTestFile(t, filepath.Join(projectRoot, "CLAUDE.md"))
	if !strings.Contains(instructions, "Before.") || !strings.Contains(instructions, "After.") {
		t.Errorf("reinstall changed CLAUDE.md outside the managed section:\n%s", instructions)
	}
	if strings.Contains(instructions, "stale") {
		t.Errorf("reinstall kept stale managed instructions:\n%s", instructions)
	}

	justfileSnapshot := justfile
	settingsSnapshot := settings
	if secondResult, secondNarration := runInstallFixture(t, projectRoot, sourceRoot, "y\n"); secondResult.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("idempotent reinstall failed: %q\n%s", secondResult.FailureReason, secondNarration)
	}
	if readTestFile(t, filepath.Join(projectRoot, "justfile")) != justfileSnapshot {
		t.Errorf("reinstall is not byte-idempotent for the justfile")
	}
	if readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json")) != settingsSnapshot {
		t.Errorf("reinstall is not byte-idempotent for settings.json")
	}
}

// Declining the single confirmation is a success with skipped work, not a refusal: that
// keeps exit 0 through the public shell path and lets update-suite tell cancelled from
// installed without a second outcome-to-number table.
func TestDecliningTheConfirmationChangesNothingAndReportsSkippedWork(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)

	result, narration := runInstallFixture(t, projectRoot, sourceRoot, "n\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, want success\n%s", result.Outcome, narration)
	}
	if !result.Cancelled {
		t.Errorf("a declined confirmation was not reported as cancelled")
	}
	if len(result.SkippedWork) != 1 || result.SkippedWork[0].Code != SkipCodeInstallCancelled {
		t.Fatalf("skipped work = %#v, want one %s", result.SkippedWork, SkipCodeInstallCancelled)
	}
	if resultmodel.ExitCode(result.Outcome) != 0 {
		t.Errorf("cancellation exits %d, want 0", resultmodel.ExitCode(result.Outcome))
	}
	for _, absentPath := range []string{".claude", "justfile", "CLAUDE.md"} {
		if _, err := os.Lstat(filepath.Join(projectRoot, absentPath)); err == nil {
			t.Errorf("cancellation created %s", absentPath)
		}
	}
	if !strings.Contains(narration, "Installation cancelled; no files were changed.") {
		t.Errorf("narration does not report the cancellation:\n%s", narration)
	}
}

// A post-write validation failure must restore every managed path to its exact prior bytes.
// The injection point is a `just` that succeeds for the candidate check and then fails for
// the installed file, which is the same shape the shell suite's flaky-just fixture uses.
func TestPostWriteValidationFailureRestoresExactManagedOriginals(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	if _, narration := mustInstall(t, projectRoot, sourceRoot); narration == "" {
		t.Fatalf("seed install produced no narration")
	}
	writeTestFile(t, filepath.Join(projectRoot, "extra-custom.txt"), "untouched\n")
	justfileBefore := readTestFile(t, filepath.Join(projectRoot, "justfile"))
	settingsBefore := readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json"))
	instructionsBefore := readTestFile(t, filepath.Join(projectRoot, "CLAUDE.md"))
	installFlakyJust(t)

	result, narration := runInstallFixture(t, projectRoot, sourceRoot, "y\n")
	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("the install reported success after post-write validation failed\n%s", narration)
	}
	if result.Outcome != resultmodel.OutcomeRolledBack {
		t.Fatalf("outcome = %q, want %q (reason %q)", result.Outcome, resultmodel.OutcomeRolledBack, result.FailureReason)
	}
	if result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Errorf("rollback status = %q, errors %v", result.Rollback.Status, result.Rollback.Errors)
	}
	if resultmodel.ExitCode(result.Outcome) != 3 {
		t.Errorf("a rolled-back install exits %d, want 3", resultmodel.ExitCode(result.Outcome))
	}
	if readTestFile(t, filepath.Join(projectRoot, "justfile")) != justfileBefore {
		t.Errorf("recovery did not restore the exact justfile bytes")
	}
	if readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json")) != settingsBefore {
		t.Errorf("recovery did not restore the exact settings bytes")
	}
	if readTestFile(t, filepath.Join(projectRoot, "CLAUDE.md")) != instructionsBefore {
		t.Errorf("recovery did not restore the exact CLAUDE.md bytes")
	}
	for _, module := range fixtureModuleNames {
		if _, err := os.Stat(filepath.Join(projectRoot, ".claude", "skills", module, "SKILL.md")); err != nil {
			t.Errorf("recovery did not restore module %s: %v", module, err)
		}
	}
	if !strings.Contains(narration, "restored every managed path and the Git index to their exact pre-install state") {
		t.Errorf("narration does not report a completed recovery:\n%s", narration)
	}
}

// A Justfile that already owns a reserved recipe outside the managed section is refused
// before the confirmation, so nothing is written and the diagnostic names the collision.
func TestReservedRecipeCollisionIsRefusedBeforeAnyConfirmation(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	writeTestFile(t, filepath.Join(projectRoot, "justfile"), "run-kanban:\n    echo mine\n")
	justfileBefore := readTestFile(t, filepath.Join(projectRoot, "justfile"))

	result, narration := runInstallFixture(t, projectRoot, sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("outcome = %q, want failure", result.Outcome)
	}
	if !strings.Contains(result.FailureReason, "target defines reserved Just recipe or alias outside managed section: run-kanban") {
		t.Errorf("failure reason does not name the collision: %q", result.FailureReason)
	}
	if strings.Contains(narration, "Install this complete four-skill suite?") {
		t.Errorf("the installer asked for confirmation before rejecting the collision:\n%s", narration)
	}
	if readTestFile(t, filepath.Join(projectRoot, "justfile")) != justfileBefore {
		t.Errorf("a refused collision changed the justfile")
	}
	if _, err := os.Lstat(filepath.Join(projectRoot, ".claude")); err == nil {
		t.Errorf("a refused collision installed modules")
	}
}

func TestANonGitProjectRootIsRefused(t *testing.T) {
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	result, _ := runInstallFixture(t, t.TempDir(), sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("outcome = %q, want failure", result.Outcome)
	}
	if !strings.Contains(result.FailureReason, "must be a Git repository") {
		t.Errorf("failure reason = %q", result.FailureReason)
	}
}

func TestASubdirectoryOfTheWorktreeIsRefused(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	subdirectory := filepath.Join(projectRoot, "subdir")
	if err := os.MkdirAll(subdirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	result, _ := runInstallFixture(t, subdirectory, sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("outcome = %q, want failure", result.Outcome)
	}
	if !strings.Contains(result.FailureReason, "must name the Git worktree root") {
		t.Errorf("failure reason = %q", result.FailureReason)
	}
}

// A managed destination reached through a symlink could redirect a write out of the project.
func TestASymlinkedManagedDestinationIsRefused(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	if err := os.MkdirAll(filepath.Join(projectRoot, ".claude", "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	elsewhere := t.TempDir()
	if err := os.Symlink(elsewhere, filepath.Join(projectRoot, ".claude", "skills", "do-work-board")); err != nil {
		t.Fatal(err)
	}
	result, _ := runInstallFixture(t, projectRoot, sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("outcome = %q, want failure", result.Outcome)
	}
	if !strings.Contains(result.FailureReason, "managed destination must not be a symlink") {
		t.Errorf("failure reason = %q", result.FailureReason)
	}
}

func mustInstall(t *testing.T, projectRoot, sourceRoot string) (InstallResult, string) {
	t.Helper()
	result, narration := runInstallFixture(t, projectRoot, sourceRoot, "y\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("seed install failed: %q\n%s", result.FailureReason, narration)
	}
	return result, narration
}

// installFlakyJust puts a `just` on PATH that succeeds once and fails afterwards, so the
// candidate passes validation and the installed file does not.
func installFlakyJust(t *testing.T) {
	t.Helper()
	stubDirectory := t.TempDir()
	counterPath := filepath.Join(stubDirectory, "invocations")
	writeTestFile(t, filepath.Join(stubDirectory, "just"),
		"#!/usr/bin/env bash\nset -u\ncount=0\n[ ! -f \""+counterPath+"\" ] || count=\"$(cat \""+counterPath+"\")\"\n"+
			"count=$((count + 1))\nprintf '%s\\n' \"$count\" > \""+counterPath+"\"\n[ \"$count\" -eq 1 ]\n")
	chmodTestFile(t, filepath.Join(stubDirectory, "just"), 0o755)
	t.Setenv("PATH", stubDirectory+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func chmodTestFile(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	if err := os.Chmod(path, mode); err != nil {
		t.Fatalf("chmod %s: %v", path, err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func modeOf(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return info.Mode().Perm()
}

func runTestGit(t *testing.T, workingDirectory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = workingDirectory
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(arguments, " "), err, output)
	}
}
