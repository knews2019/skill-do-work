package suiteinstall

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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
