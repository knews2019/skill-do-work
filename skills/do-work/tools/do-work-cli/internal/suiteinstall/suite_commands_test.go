package suiteinstall

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/managedsection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// runCommand drives a command through the real runtime, which is the seam the launchers
// hit. It returns the exit status, stdout (the rendered result and nothing else) and the
// narration stream.
func runCommand(t *testing.T, repositoryRoot, outputFormat, confirmation string, argv ...string) (int, string, string) {
	t.Helper()
	var standardOutput bytes.Buffer
	var narration bytes.Buffer
	handlers := Handlers(&narration, strings.NewReader(confirmation))
	arguments := append([]string{"--repo-root", repositoryRoot, "--format", outputFormat}, argv...)
	exitCode := commandruntime.NewRuntime(&standardOutput, handlers).Run(arguments)
	return exitCode, standardOutput.String(), narration.String()
}

func TestEveryCommandIsRegisteredUnderItsArgvToken(t *testing.T) {
	handlers := Handlers(nil, nil)
	for _, commandName := range []string{
		CommandInstallSuite, CommandUpdateSuite, CommandReplaceSection,
		CommandValidateManifest, CommandFetchArchive,
	} {
		if _, registered := handlers[commandName]; !registered {
			t.Errorf("command %q is not registered", commandName)
		}
	}
	if len(handlers) != 5 {
		t.Errorf("handler count = %d, want the five commands this package owns", len(handlers))
	}
}

func TestReplaceSectionReportsItsChangeAndExitsZero(t *testing.T) {
	directory := t.TempDir()
	sectionPath := filepath.Join(directory, "section.just")
	writeTestFile(t, sectionPath, "# >>> do-work:recipes >>>\nrun-kanban:\n    echo x\n# <<< do-work:recipes <<<\n")
	targetPath := filepath.Join(directory, "justfile")
	writeTestFile(t, targetPath, "custom:\n    echo custom\n")

	exitCode, standardOutput, narration := runCommand(t, directory, "json", "",
		CommandReplaceSection, "--target", targetPath, "--section-file", sectionPath)
	if exitCode != 0 {
		t.Fatalf("exit = %d\n%s", exitCode, standardOutput)
	}
	if narration != "" {
		t.Errorf("replace-section narrated on a clean run: %q", narration)
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal([]byte(standardOutput), &decoded); err != nil {
		t.Fatalf("decode result: %v\n%s", err, standardOutput)
	}
	if decoded.Command != CommandReplaceSection {
		t.Errorf("command = %q, want %q", decoded.Command, CommandReplaceSection)
	}
	if len(decoded.Changes) != 1 || decoded.Changes[0].Kind != managedsection.ChangeKindModified {
		t.Fatalf("changes = %#v", decoded.Changes)
	}
	if decoded.Rollback.Status != resultmodel.RollbackNotNeeded {
		t.Errorf("rollback status = %q", decoded.Rollback.Status)
	}
}

// A target that already matches byte for byte is a success with no change, which is what
// makes a reinstall byte-idempotent.
func TestReplaceSectionReportsNoChangeWhenTheTargetAlreadyMatches(t *testing.T) {
	directory := t.TempDir()
	sectionBytes := "# >>> do-work:recipes >>>\nrun-kanban:\n    echo x\n# <<< do-work:recipes <<<\n"
	sectionPath := filepath.Join(directory, "section.just")
	writeTestFile(t, sectionPath, sectionBytes)
	targetPath := filepath.Join(directory, "justfile")
	writeTestFile(t, targetPath, sectionBytes)

	exitCode, standardOutput, _ := runCommand(t, directory, "json", "",
		CommandReplaceSection, "--target", targetPath, "--section-file", sectionPath)
	if exitCode != 0 {
		t.Fatalf("exit = %d\n%s", exitCode, standardOutput)
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal([]byte(standardOutput), &decoded); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(decoded.Changes) != 0 {
		t.Errorf("an unchanged target reported changes: %#v", decoded.Changes)
	}
}

// A reserved-recipe collision is a refusal the consumer can resolve, so it exits 1 and its
// evidence is the phrase three behavioural suites already assert on.
func TestReservedRecipeCollisionRefusesWithExitOneAndTheCurrentPhrase(t *testing.T) {
	directory := t.TempDir()
	sectionPath := filepath.Join(directory, "section.just")
	writeTestFile(t, sectionPath, "# >>> do-work:recipes >>>\nrun-kanban:\n    echo x\nalias rk := run-kanban\n# <<< do-work:recipes <<<\n")
	targetPath := filepath.Join(directory, "justfile")
	writeTestFile(t, targetPath, "alias rk := other\nrun-kanban:\n    echo mine\n")

	exitCode, standardOutput, _ := runCommand(t, directory, "text", "",
		CommandReplaceSection, "--target", targetPath, "--section-file", sectionPath, "--reject-recipe-collisions")
	if exitCode != 1 {
		t.Fatalf("exit = %d, want 1\n%s", exitCode, standardOutput)
	}
	expectedLine := "finding " + managedsection.FailureReservedRecipeCollision +
		" [error]: target defines reserved Just recipe or alias outside managed section: rk, run-kanban"
	if !containsExactRenderedLine(standardOutput, expectedLine) {
		t.Errorf("stdout does not carry the exact finding line %q:\n%s", expectedLine, standardOutput)
	}
}

// A malformed marker pair is the caller's mistake to fix, not a refusal, so it exits 2.
func TestMalformedMarkersFailWithExitTwo(t *testing.T) {
	directory := t.TempDir()
	sectionPath := filepath.Join(directory, "section.just")
	writeTestFile(t, sectionPath, "# >>> do-work:recipes >>>\nrun-kanban:\n    echo x\n# <<< do-work:recipes <<<\n")
	targetPath := filepath.Join(directory, "justfile")
	writeTestFile(t, targetPath, "# >>> do-work:recipes >>>\nno end marker\n")

	exitCode, standardOutput, _ := runCommand(t, directory, "text", "",
		CommandReplaceSection, "--target", targetPath, "--section-file", sectionPath)
	if exitCode != 2 {
		t.Fatalf("exit = %d, want 2\n%s", exitCode, standardOutput)
	}
	if !strings.Contains(standardOutput, "target must contain exactly one begin marker and one end marker") {
		t.Errorf("stdout does not carry the current diagnostic:\n%s", standardOutput)
	}
}

func TestValidateManifestReportsItsSummaryLineOnStdout(t *testing.T) {
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	exitCode, standardOutput, narration := runCommand(t, sourceRoot, "text", "",
		CommandValidateManifest, "--root", sourceRoot)
	if exitCode != 0 {
		t.Fatalf("exit = %d\n%s", exitCode, standardOutput)
	}
	if narration != "" {
		t.Errorf("validate-manifest narrated on a clean run: %q", narration)
	}
	if !strings.Contains(standardOutput, "suite manifest valid: v"+fixtureSuiteVersion+" (4 modules)") {
		t.Errorf("stdout does not carry the summary line:\n%s", standardOutput)
	}
}

func TestValidateManifestFailsWithExitTwoAndNamesTheViolation(t *testing.T) {
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	writeTestFile(t, filepath.Join(sourceRoot, "suite", "modules.tsv"), "source\tdestination\textra\n")

	exitCode, standardOutput, _ := runCommand(t, sourceRoot, "text", "",
		CommandValidateManifest, "--root", sourceRoot)
	if exitCode != 2 {
		t.Fatalf("exit = %d, want 2\n%s", exitCode, standardOutput)
	}
	if !strings.Contains(standardOutput, "manifest header must be exactly: source<TAB>destination") {
		t.Errorf("stdout does not carry the current diagnostic:\n%s", standardOutput)
	}
}

// The winning route rides on stdout as a recorded change rather than as a finding, so a
// success outcome never carries an info-severity finding just to name a route.
func TestFetchArchiveRecordsTheWinningRouteAsAChange(t *testing.T) {
	toolDirectory, upstreamURL := newUpstreamArchiveServer(t, fixtureSuiteVersion)
	targetPath := filepath.Join(t.TempDir(), "upstream.tar.gz")
	t.Setenv("DO_WORK_TEST_TOOL_DIRECTORY", toolDirectory)

	// The command resolves its primitive from the running binary's location, which a package
	// test cannot move, so the fetch is driven through the package the command calls.
	exitCode, standardOutput, _ := runCommand(t, filepath.Dir(targetPath), "json", "",
		CommandFetchArchive, "--target", targetPath, "--url", upstreamURL, "--repo-url", "")
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal([]byte(standardOutput), &decoded); err != nil {
		t.Fatalf("decode result: %v\n%s", err, standardOutput)
	}
	if exitCode == 0 {
		if len(decoded.Changes) != 1 || decoded.Changes[0].Path != targetPath {
			t.Fatalf("a successful fetch did not record the archive path: %#v", decoded.Changes)
		}
		if !strings.HasPrefix(decoded.Changes[0].Detail, "upstream archive fetched") {
			t.Errorf("change detail does not name the route: %q", decoded.Changes[0].Detail)
		}
		return
	}
	// Without a reachable primitive both routes are unavailable, which must still be a
	// complete failure finding naming the escape hatch rather than a silent empty result.
	if exitCode != 2 {
		t.Fatalf("exit = %d, want 0 or 2\n%s", exitCode, standardOutput)
	}
	if len(decoded.Findings) != 1 || !strings.Contains(decoded.Findings[0].Evidence[0], "DO_WORK_UPSTREAM_URL") {
		t.Errorf("failure finding does not name the escape hatch: %#v", decoded.Findings)
	}
	if _, err := os.Stat(targetPath); err == nil {
		t.Errorf("a failed fetch published the archive path")
	}
}

func TestUnknownCommandOptionsFailWithARunnableNextArgv(t *testing.T) {
	directory := t.TempDir()
	exitCode, standardOutput, _ := runCommand(t, directory, "json", "",
		CommandValidateManifest, "--nonsense", "value")
	if exitCode != 2 {
		t.Fatalf("exit = %d, want 2\n%s", exitCode, standardOutput)
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal([]byte(standardOutput), &decoded); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(decoded.Findings) != 1 {
		t.Fatalf("findings = %#v", decoded.Findings)
	}
	nextArgv := decoded.Findings[0].NextArgv
	if len(nextArgv) == 0 || nextArgv[len(nextArgv)-1] != CommandValidateManifest {
		t.Errorf("next argv %v does not end in the command name", nextArgv)
	}
	for _, argument := range nextArgv {
		if argument == "<command>" {
			t.Errorf("next argv still carries the placeholder: %v", nextArgv)
		}
	}
}

// stdout carries the rendered result and nothing else. Every progress line, review diff and
// prompt goes to narration, which is what keeps --format json parseable.
func TestInstallNarrationNeverReachesStandardOutput(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
	archivePath := buildArchiveFromSourceTree(t, sourceRoot)

	exitCode, standardOutput, narration := runCommand(t, projectRoot, "json", "y\n",
		CommandInstallSuite, "--archive", archivePath)
	if exitCode != 0 {
		t.Fatalf("exit = %d\nstdout:\n%s\nnarration:\n%s", exitCode, standardOutput, narration)
	}
	var decoded resultmodel.CommandResult
	if err := json.Unmarshal([]byte(standardOutput), &decoded); err != nil {
		t.Fatalf("stdout is not one parseable result: %v\n%s", err, standardOutput)
	}
	if decoded.Command != CommandInstallSuite {
		t.Errorf("command = %q", decoded.Command)
	}
	for _, narrationOnlyPhrase := range []string{
		"Ready to install do-work suite v",
		"Reviewing the complete managed install before overwrite:",
		"Install this complete four-skill suite? [y/N]",
		"Installed do-work suite v",
	} {
		if !strings.Contains(narration, narrationOnlyPhrase) {
			t.Errorf("narration is missing %q:\n%s", narrationOnlyPhrase, narration)
		}
		if strings.Contains(standardOutput, narrationOnlyPhrase) {
			t.Errorf("narration phrase %q leaked onto stdout:\n%s", narrationOnlyPhrase, standardOutput)
		}
	}
}

func TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("built-binary signal integration is heavy-only")
	}
	binaryPath := buildTestCLIBinary(t)
	for _, commandName := range []string{CommandInstallSuite, CommandUpdateSuite} {
		for _, signalCase := range interruptingSignalCases() {
			t.Run(commandName+"/"+signalCase.name, func(t *testing.T) {
				projectRoot := newProjectRepository(t)
				arguments := []string{"--repo-root", projectRoot, "--format", "json", commandName}
				var installedVersionPath string
				if commandName == CommandInstallSuite {
					sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
					arguments = append(arguments, "--archive", buildArchiveFromSourceTree(t, sourceRoot))
				} else {
					skillRoot := installSuiteForUpdate(t, projectRoot, "0.100.0")
					_, upstreamURL := newUpstreamArchiveServer(t, "0.200.0")
					installedVersionPath = filepath.Join(skillRoot, "VERSION")
					arguments = append(arguments, "--skill-root", skillRoot, "--upstream-url", upstreamURL)
				}

				runBuiltCLIAtBlockedConfirmation(t, binaryPath, signalCase.signal, arguments)
				if commandName == CommandInstallSuite {
					for _, absentPath := range []string{".claude", "justfile", "CLAUDE.md"} {
						if _, err := os.Lstat(filepath.Join(projectRoot, absentPath)); err == nil {
							t.Errorf("%s created %s before confirmation", signalCase.name, absentPath)
						}
					}
				} else if installedVersion := readTestFile(t, installedVersionPath); installedVersion != "0.100.0\n" {
					t.Errorf("%s changed installed VERSION before confirmation: %q", signalCase.name, installedVersion)
				}
			})
		}
	}
}

// An interruption that lands AFTER the first write must still exit on the signal status,
// not on the rollback outcome's number. This case had no coverage: it recovers correctly and
// returns OutcomeRolledBack (exit 3), so before the single-exit-owner fix it raced the
// handler's 130 exactly as the pre-confirmation case did, and reported 3 whenever the
// ordinary return path won. 130 is the right answer for both — the process was killed by a
// signal, the rendered result still carries the rollback record, and the shell installer's
// interrupted status is what callers already assert.
func TestBuiltInstallExits130AfterRecoveringASignalInterruptedMidWriteInstall(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("built-binary signal integration is heavy-only")
	}
	binaryPath := buildTestCLIBinary(t)
	for _, signalCase := range interruptingSignalCases() {
		t.Run(signalCase.name, func(t *testing.T) {
			projectRoot := newProjectRepository(t)
			sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
			archivePath := buildArchiveFromSourceTree(t, sourceRoot)
			blockedMarkerPath := installPostWriteBlockingJust(t)

			narration := runBuiltCLIInterruptedAtMarker(t, binaryPath, signalCase.signal, blockedMarkerPath,
				[]string{"--repo-root", projectRoot, "--format", "json", CommandInstallSuite, "--archive", archivePath})

			if !strings.Contains(narration, "restored every managed path and the Git index to their exact pre-install state") {
				t.Errorf("%s exited before recovery reported completion:\n%s", signalCase.name, narration)
			}
			for _, absentPath := range []string{
				"justfile",
				"CLAUDE.md",
				filepath.Join(".claude", "settings.json"),
				filepath.Join(".claude", "skills", "do-work"),
			} {
				if _, err := os.Lstat(filepath.Join(projectRoot, absentPath)); err == nil {
					t.Errorf("%s left %s behind instead of recovering it", signalCase.name, absentPath)
				}
			}
		})
	}
}

// A signal that lands after the install is verified must NOT take the process. The install's
// last cancellable step is the installed Justfile's post-write validation; past it every
// managed byte is written and proven, cleanup skips recovery, and the caller is owed exit 0
// plus the rendered result on stdout. Without the installVerified guard this case exits 130
// with an empty stdout while narration reports a complete install (REQ-525 F1).
//
// The window between that last subprocess and the interrupted-exit owner is sub-millisecond,
// so a signal timed from outside loses it — an earlier attempt at this test signalled from
// the stub itself and the child had already exited 0 every time. The child is parked inside
// the window instead: its narration pipe is filled to capacity, so the one write it makes
// after setting installVerified — the success line — blocks until this test drains it, and
// nothing can advance to the exit owner until then.
func TestBuiltInstallExitsZeroWhenASignalLandsAfterTheInstallIsVerified(t *testing.T) {
	if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" {
		t.Skip("built-binary signal integration is heavy-only")
	}
	binaryPath := buildTestCLIBinary(t)
	for _, signalCase := range interruptingSignalCases() {
		t.Run(signalCase.name, func(t *testing.T) {
			projectRoot := newProjectRepository(t)
			sourceRoot := newSuiteSourceTree(t, fixtureSuiteVersion)
			archivePath := buildArchiveFromSourceTree(t, sourceRoot)
			reapMarkerPath := installReapMarkingJust(t)

			exitCode, standardOutput, narration := runBuiltCLISignalledOnBlockedNarration(t, binaryPath,
				signalCase.signal, reapMarkerPath,
				[]string{"--repo-root", projectRoot, "--format", "json", CommandInstallSuite, "--archive", archivePath})

			if exitCode != 0 {
				t.Fatalf("%s exit = %d, want 0; stdout carried %d bytes\nstdout:\n%s\nnarration:\n%s",
					signalCase.name, exitCode, len(standardOutput), standardOutput, narration)
			}
			var decoded resultmodel.CommandResult
			if err := json.Unmarshal([]byte(standardOutput), &decoded); err != nil {
				t.Fatalf("%s wrote no parseable result to stdout (%d bytes): %v\nnarration:\n%s",
					signalCase.name, len(standardOutput), err, narration)
			}
			if decoded.Outcome != resultmodel.OutcomeSuccess {
				t.Errorf("%s outcome = %q, want %q", signalCase.name, decoded.Outcome, resultmodel.OutcomeSuccess)
			}
			if decoded.Rollback.Status != resultmodel.RollbackNotNeeded {
				t.Errorf("%s rollback status = %q, want %q",
					signalCase.name, decoded.Rollback.Status, resultmodel.RollbackNotNeeded)
			}
			if installedVersion := readTestFile(t,
				filepath.Join(projectRoot, ".claude", "skills", "do-work", "VERSION")); installedVersion != fixtureSuiteVersion+"\n" {
				t.Errorf("%s installed VERSION = %q, want %q",
					signalCase.name, installedVersion, fixtureSuiteVersion+"\n")
			}
			if !strings.Contains(narration, "Installed do-work suite v"+fixtureSuiteVersion) {
				t.Errorf("%s did not report a complete install:\n%s", signalCase.name, narration)
			}
			if strings.Contains(narration, "recovering managed paths and Git index") {
				t.Errorf("%s recovered a verified install:\n%s", signalCase.name, narration)
			}
		})
	}
}

// interruptingSignalCases are the three signals RunInstall arms for. Every interruption test
// covers all three, because they share one handler and one exit status.
func interruptingSignalCases() []struct {
	name   string
	signal os.Signal
} {
	return []struct {
		name   string
		signal os.Signal
	}{
		{name: "HUP", signal: syscall.SIGHUP},
		{name: "INT", signal: syscall.SIGINT},
		{name: "TERM", signal: syscall.SIGTERM},
	}
}

// singleProcessorEnvironment pins the built child to one scheduler processor.
//
// The interrupted exit status used to be decided by whichever goroutine reached os.Exit
// first — the signal handler's 130 or the ordinary return path's own outcome number — and at
// default GOMAXPROCS the handler usually won, so the defect surfaced as an occasional
// failure under parallel load (3 of 7 full-module runs). At one processor the handler is
// starved for a slot and the ordinary path wins every time, which turned that flake class
// into a guaranteed pre-fix failure and makes these tests a real lock-in rather than a
// scheduling coincidence (REQ-525).
func singleProcessorEnvironment() []string {
	return append(os.Environ(), "GOMAXPROCS=1")
}

var (
	builtTestCLIBinaryOnce sync.Once
	builtTestCLIBinaryPath string
	builtTestCLIBinaryDir  string
	builtTestCLIBinaryErr  error
)

func TestMain(m *testing.M) {
	code := m.Run()
	if builtTestCLIBinaryDir != "" {
		_ = os.RemoveAll(builtTestCLIBinaryDir)
	}
	os.Exit(code)
}

func buildTestCLIBinary(t *testing.T) string {
	t.Helper()
	builtTestCLIBinaryOnce.Do(func() {
		moduleRoot, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			builtTestCLIBinaryErr = err
			return
		}
		temporaryDirectory, err := os.MkdirTemp("", "suiteinstall-cli-test-*")
		if err != nil {
			builtTestCLIBinaryErr = err
			return
		}
		builtTestCLIBinaryDir = temporaryDirectory
		builtTestCLIBinaryPath = filepath.Join(temporaryDirectory, "do-work-cli")
		command := exec.Command("go", "build", "-o", builtTestCLIBinaryPath, "./cmd/do-work-cli")
		command.Dir = moduleRoot
		if output, err := command.CombinedOutput(); err != nil {
			builtTestCLIBinaryErr = fmt.Errorf("build do-work-cli: %w: %s", err, output)
			return
		}
	})
	if builtTestCLIBinaryErr != nil {
		t.Fatalf("resolve test CLI binary: %v", builtTestCLIBinaryErr)
	}
	return builtTestCLIBinaryPath
}

func runBuiltCLIAtBlockedConfirmation(t *testing.T, binaryPath string, interrupt os.Signal, arguments []string) {
	t.Helper()
	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer stdinReader.Close()
	defer stdinWriter.Close()
	stderrPath := filepath.Join(t.TempDir(), "stderr")
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stderrFile.Close()
	stdoutFile, err := os.Create(filepath.Join(t.TempDir(), "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	defer stdoutFile.Close()

	command := exec.Command(binaryPath, arguments...)
	command.Env = singleProcessorEnvironment()
	command.Stdin = stdinReader
	command.Stdout = stdoutFile
	command.Stderr = stderrFile
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	_ = stdinReader.Close()
	waitForPromptInFile(t, command, stderrPath)
	// Keep stdin open and send no byte: the signal itself must release confirmation.
	if err := command.Process.Signal(interrupt); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("send %v: %v", interrupt, err)
	}
	waitResult := make(chan error, 1)
	go func() { waitResult <- command.Wait() }()
	select {
	case waitErr := <-waitResult:
		var exitError *exec.ExitError
		if !errors.As(waitErr, &exitError) || exitError.ExitCode() != 130 {
			output, _ := os.ReadFile(stderrPath)
			t.Fatalf("signal %v exit = %v, want 130\nstderr:\n%s", interrupt, waitErr, output)
		}
	case <-time.After(5 * time.Second):
		_ = command.Process.Kill()
		<-waitResult
		output, _ := os.ReadFile(stderrPath)
		t.Fatalf("signal %v did not stop blocked confirmation without input\nstderr:\n%s", interrupt, output)
	}
}

func waitForPromptInFile(t *testing.T, command *exec.Cmd, path string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		output, err := os.ReadFile(path)
		if err == nil && bytes.Contains(output, []byte("Install this complete four-skill suite? [y/N]")) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	_ = command.Process.Kill()
	_ = command.Wait()
	output, _ := os.ReadFile(path)
	t.Fatalf("built CLI did not reach confirmation prompt:\n%s", output)
}

// installPostWriteBlockingJust puts a `just` on PATH that succeeds for the pre-confirmation
// candidate check and then parks forever on its second invocation, which is the installed
// Justfile's post-write validation. `just` is invoked exactly twice per install, so the
// second invocation pins the child mid-transaction: writeStarted is set, every managed
// module and configuration file is already replaced, and recovery is required. It returns
// the marker path the stub touches once it is parked.
func installPostWriteBlockingJust(t *testing.T) string {
	t.Helper()
	stubDirectory := t.TempDir()
	counterPath := filepath.Join(stubDirectory, "invocations")
	blockedMarkerPath := filepath.Join(stubDirectory, "blocked")
	writeTestFile(t, filepath.Join(stubDirectory, "just"),
		"#!/usr/bin/env bash\nset -u\ncount=0\n[ ! -f \""+counterPath+"\" ] || count=\"$(cat \""+counterPath+"\")\"\n"+
			"count=$((count + 1))\nprintf '%s\\n' \"$count\" > \""+counterPath+"\"\n"+
			"[ \"$count\" -ne 1 ] || exit 0\n"+
			// exec so cancelling the work context kills the sleep itself, not a shell wrapper.
			"printf 'blocked\\n' > \""+blockedMarkerPath+"\"\nexec sleep 600\n")
	chmodTestFile(t, filepath.Join(stubDirectory, "just"), 0o755)
	t.Setenv("PATH", stubDirectory+string(os.PathListSeparator)+os.Getenv("PATH"))
	return blockedMarkerPath
}

// runBuiltCLIInterruptedAtMarker confirms the install, waits for the child to park at the
// named marker, signals it, and requires the interrupted status. It returns the narration so
// the caller can assert on recovery.
func runBuiltCLIInterruptedAtMarker(t *testing.T, binaryPath string, interrupt os.Signal,
	markerPath string, arguments []string) string {
	t.Helper()
	stderrPath := filepath.Join(t.TempDir(), "stderr")
	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stderrFile.Close()
	stdoutFile, err := os.Create(filepath.Join(t.TempDir(), "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	defer stdoutFile.Close()

	command := exec.Command(binaryPath, arguments...)
	command.Env = singleProcessorEnvironment()
	command.Stdin = strings.NewReader("y\n")
	command.Stdout = stdoutFile
	command.Stderr = stderrFile
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	waitForMarkerFile(t, command, markerPath, stderrPath)
	if err := command.Process.Signal(interrupt); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("send %v: %v", interrupt, err)
	}
	waitResult := make(chan error, 1)
	go func() { waitResult <- command.Wait() }()
	select {
	case waitErr := <-waitResult:
		narration, _ := os.ReadFile(stderrPath)
		var exitError *exec.ExitError
		if !errors.As(waitErr, &exitError) || exitError.ExitCode() != 130 {
			t.Fatalf("signal %v exit = %v, want 130\nstderr:\n%s", interrupt, waitErr, narration)
		}
		return string(narration)
	case <-time.After(30 * time.Second):
		_ = command.Process.Kill()
		<-waitResult
		narration, _ := os.ReadFile(stderrPath)
		t.Fatalf("signal %v did not stop the mid-write install\nstderr:\n%s", interrupt, narration)
		return ""
	}
}

// installReapMarkingJust puts a `just` on PATH that succeeds for the pre-confirmation
// candidate check and, on its second invocation — the installed Justfile's post-write
// validation, the install's last cancellable subprocess — records that the installer has
// collected its exit status, then exits 0. It returns that marker path.
//
// A detached grandchild writes the marker, and `kill -0` keeps succeeding for a zombie, so it
// fails only once the installer has reaped the stub. The marker therefore proves no
// subprocess is in flight, which is what makes it safe to signal: a signal delivered any
// earlier cancels the work context with a subprocess still running, which is the
// interrupted-mid-write case the preceding test covers.
func installReapMarkingJust(t *testing.T) string {
	t.Helper()
	stubDirectory := t.TempDir()
	counterPath := filepath.Join(stubDirectory, "invocations")
	reapMarkerPath := filepath.Join(stubDirectory, "reaped")
	writeTestFile(t, filepath.Join(stubDirectory, "just"), fmt.Sprintf(`#!/usr/bin/env bash
set -u
count=0
[ ! -f %[1]q ] || count="$(cat %[1]q)"
count=$((count + 1))
printf '%%s\n' "$count" > %[1]q
[ "$count" -ne 1 ] || exit 0
bash -c 'while kill -0 "$1" 2>/dev/null && [ "$SECONDS" -lt 60 ]; do :; done; printf reaped > "$2"' \
	marker-after-reap "$$" %[2]q &
exit 0
`, counterPath, reapMarkerPath))
	chmodTestFile(t, filepath.Join(stubDirectory, "just"), 0o755)
	t.Setenv("PATH", stubDirectory+string(os.PathListSeparator)+os.Getenv("PATH"))
	return reapMarkerPath
}

// runBuiltCLISignalledOnBlockedNarration drives the built CLI to a verified install, signals
// it while it is parked on a narration write it cannot complete, then releases that write and
// returns the exit status, stdout and the narration the installer produced after its
// confirmation.
//
// Both the confirmation and the release are held by this test, which is what removes every
// timing assumption: the pipe is filled before the install is allowed to start writing, and
// it is drained only after the signal has been delivered to a child that cannot leave the
// window.
func runBuiltCLISignalledOnBlockedNarration(t *testing.T, binaryPath string, interrupt os.Signal,
	reapMarkerPath string, arguments []string) (int, string, string) {
	t.Helper()
	confirmationReader, confirmationWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer confirmationWriter.Close()
	narrationReader, narrationWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer narrationReader.Close()
	defer narrationWriter.Close()
	pipeCapacity := measurePipeCapacity(t)
	stdoutPath := filepath.Join(t.TempDir(), "stdout")
	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		t.Fatal(err)
	}
	defer stdoutFile.Close()

	command := exec.Command(binaryPath, arguments...)
	command.Env = singleProcessorEnvironment()
	command.Stdin = confirmationReader
	command.Stdout = stdoutFile
	command.Stderr = narrationWriter
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	// The child holds its own duplicates of both pipe ends; the copies kept here are what
	// apply the back-pressure and later release it.
	_ = confirmationReader.Close()

	promptNarration := readNarrationThroughPrompt(t, command, narrationReader)
	fillerBytes := parkTheNextNarrationWrite(t, command, narrationReader, narrationWriter, pipeCapacity)
	if _, err := confirmationWriter.WriteString("y\n"); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("confirm the install: %v", err)
	}
	_ = confirmationWriter.Close()
	waitForReapMarker(t, command, reapMarkerPath, promptNarration)

	if err := command.Process.Signal(interrupt); err != nil {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("send %v: %v", interrupt, err)
	}
	// Closing this copy is what lets the drain below end at EOF; the child keeps its own.
	_ = narrationWriter.Close()
	drained := make(chan []byte, 1)
	go func() {
		remaining, _ := io.ReadAll(narrationReader)
		drained <- remaining
	}()
	waitResult := make(chan error, 1)
	go func() { waitResult <- command.Wait() }()
	select {
	case waitErr := <-waitResult:
		var exitError *exec.ExitError
		if waitErr != nil && !errors.As(waitErr, &exitError) {
			t.Fatalf("wait for the released install: %v", waitErr)
		}
	case <-time.After(60 * time.Second):
		_ = command.Process.Kill()
		<-waitResult
		t.Fatalf("signal %v left the install parked after its narration was released\nnarration:\n%s",
			interrupt, promptNarration)
	}
	remaining := <-drained
	if len(remaining) < fillerBytes {
		t.Fatalf("narration pipe returned %d bytes, fewer than the %d filler bytes written",
			len(remaining), fillerBytes)
	}
	return command.ProcessState.ExitCode(), readTestFile(t, stdoutPath),
		promptNarration + string(remaining[fillerBytes:])
}

// readNarrationThroughPrompt drains the narration pipe up to and including the confirmation
// prompt — the installer's last write before it reads stdin — which leaves the pipe empty for
// the fill that follows.
func readNarrationThroughPrompt(t *testing.T, command *exec.Cmd, narrationReader *os.File) string {
	t.Helper()
	if err := narrationReader.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var accumulated bytes.Buffer
	buffer := make([]byte, 4096)
	for !bytes.Contains(accumulated.Bytes(), []byte("Install this complete four-skill suite? [y/N]")) {
		count, readErr := narrationReader.Read(buffer)
		if count > 0 {
			accumulated.Write(buffer[:count])
		}
		if readErr != nil {
			_ = command.Process.Kill()
			_ = command.Wait()
			t.Fatalf("narration ended before the confirmation prompt: %v\n%s", readErr, accumulated.String())
		}
	}
	if err := narrationReader.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	return accumulated.String()
}

// measurePipeCapacity fills a throwaway pipe until it refuses another byte, which is this
// platform's pipe buffer size.
//
// The refusal has to be measured here rather than on the narration pipe itself: os.Pipe hands
// back poller-backed descriptors whose writes report EAGAIN, but exec takes the raw
// descriptor of any *os.File it gives a child, which puts that pipe back into blocking mode.
// A fill attempted on it after Start would block instead of refusing.
func measurePipeCapacity(t *testing.T) int {
	t.Helper()
	measuredReader, measuredWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer measuredReader.Close()
	defer measuredWriter.Close()
	rawConnection, err := measuredWriter.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	filler := bytes.Repeat([]byte("."), 4096)
	capacity := 0
	controlErr := rawConnection.Write(func(descriptor uintptr) bool {
		for {
			count, writeErr := syscall.Write(int(descriptor), filler)
			if count > 0 {
				capacity += count
			}
			if writeErr == nil || writeErr == syscall.EINTR {
				continue
			}
			if writeErr != syscall.EAGAIN {
				t.Errorf("measure the pipe buffer: %v", writeErr)
			}
			return true
		}
	})
	if controlErr != nil {
		t.Fatal(controlErr)
	}
	if capacity == 0 {
		t.Fatal("a fresh pipe accepted no bytes")
	}
	return capacity
}

// parkTheNextNarrationWrite leaves the narration pipe holding exactly one full buffer of
// filler, so whatever the installer writes next blocks until this test drains the pipe. It
// returns the filler byte count, which is also the offset the installer's own bytes resume
// at: pipe order is first-in-first-out, and each narrate call is a single sub-PIPE_BUF write
// that cannot interleave with the filler.
//
// The pipe must be empty for a write of exactly one buffer to complete without blocking, so
// emptiness is proven rather than assumed. It is provable here because the installer is
// parked on its confirmation read and cannot narrate again until this test confirms.
func parkTheNextNarrationWrite(t *testing.T, command *exec.Cmd,
	narrationReader, narrationWriter *os.File, pipeCapacity int) int {
	t.Helper()
	if err := narrationReader.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	leftover := make([]byte, 4096)
	if count, err := narrationReader.Read(leftover); !errors.Is(err, os.ErrDeadlineExceeded) {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("narration pipe was not drained to the prompt: read %d bytes (%v): %q",
			count, err, leftover[:count])
	}
	if err := narrationReader.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	count, err := narrationWriter.Write(bytes.Repeat([]byte("."), pipeCapacity))
	if err != nil {
		t.Fatalf("fill the narration pipe: %v", err)
	}
	return count
}

// waitForReapMarker blocks until the just stub records that the installer reaped it, which is
// the point past every cancellable step of the install.
func waitForReapMarker(t *testing.T, command *exec.Cmd, markerPath, narration string) {
	t.Helper()
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Lstat(markerPath); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	_ = command.Process.Kill()
	_ = command.Wait()
	t.Fatalf("built CLI did not finish the post-write Justfile validation\nnarration:\n%s", narration)
}

func waitForMarkerFile(t *testing.T, command *exec.Cmd, markerPath, narrationPath string) {
	t.Helper()
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Lstat(markerPath); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	_ = command.Process.Kill()
	_ = command.Wait()
	narration, _ := os.ReadFile(narrationPath)
	t.Fatalf("built CLI did not reach %s:\n%s", markerPath, narration)
}

func buildArchiveFromSourceTree(t *testing.T, sourceRoot string) string {
	t.Helper()
	stagingParent := t.TempDir()
	stagingRoot := filepath.Join(stagingParent, "skill-do-work-main")
	runTestCommand(t, "cp", "-R", sourceRoot, stagingRoot)
	archivePath := filepath.Join(stagingParent, "suite.tar.gz")
	runTestCommand(t, "tar", "czf", archivePath, "-C", stagingParent, "skill-do-work-main")
	return archivePath
}

func containsExactRenderedLine(rendered, wanted string) bool {
	for _, line := range strings.Split(rendered, "\n") {
		if line == wanted {
			return true
		}
	}
	return false
}

func runTestCommand(t *testing.T, name string, arguments ...string) {
	t.Helper()
	command := exec.Command(name, arguments...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("%s %s: %v: %s", name, strings.Join(arguments, " "), err, output)
	}
}
