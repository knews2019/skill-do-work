package suiteinstall

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func buildTestCLIBinary(t *testing.T) string {
	t.Helper()
	moduleRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := filepath.Join(t.TempDir(), "do-work-cli")
	command := exec.Command("go", "build", "-o", binaryPath, "./cmd/do-work-cli")
	command.Dir = moduleRoot
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build do-work-cli: %v: %s", err, output)
	}
	return binaryPath
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
