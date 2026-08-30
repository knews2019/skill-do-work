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

// newUpstreamArchiveServer stages a suite tree as a tarball and returns a download script
// that publishes it, which is how an update is driven without a network.
func newUpstreamArchiveServer(t *testing.T, suiteVersion string) (toolDirectory string, upstreamURL string) {
	t.Helper()
	sourceRoot := newSuiteSourceTree(t, suiteVersion)
	stagingParent := t.TempDir()
	stagingRoot := filepath.Join(stagingParent, "skill-do-work-main")
	if output, err := exec.Command("cp", "-R", sourceRoot, stagingRoot).CombinedOutput(); err != nil {
		t.Fatalf("stage upstream tree: %v: %s", err, output)
	}
	archivePath := filepath.Join(stagingParent, "upstream.tar.gz")
	if output, err := exec.Command("tar", "czf", archivePath, "-C", stagingParent, "skill-do-work-main").CombinedOutput(); err != nil {
		t.Fatalf("build upstream archive: %v: %s", err, output)
	}

	// LocateAtomicDownloadScript probes <tools>/../scripts and <tools>/../skills/do-work/scripts.
	toolDirectory = filepath.Join(t.TempDir(), "tools")
	scriptPath := filepath.Join(filepath.Dir(toolDirectory), "scripts", "atomic-download.sh")
	writeTestFile(t, scriptPath, "#!/usr/bin/env bash\nset -eu\ncp \""+archivePath+"\" \"$2\"\n")
	chmodTestFile(t, scriptPath, 0o755)
	if err := os.MkdirAll(toolDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	return toolDirectory, "https://example.invalid/suite.tar.gz"
}

// installSuiteForUpdate seeds a project with an installed suite so the updater has something
// to compare against, then returns the installed skill root.
func installSuiteForUpdate(t *testing.T, projectRoot, installedVersion string) string {
	t.Helper()
	sourceRoot := newSuiteSourceTree(t, installedVersion)
	mustInstall(t, projectRoot, sourceRoot)
	return filepath.Join(projectRoot, ".claude", "skills", "do-work")
}

func runUpdateFixture(t *testing.T, projectRoot, skillRoot, toolDirectory, upstreamURL, confirmation string) (UpdateResult, string) {
	t.Helper()
	var narration bytes.Buffer
	result := RunUpdate(context.Background(), UpdateOptions{
		ProjectRoot:        projectRoot,
		InstalledSkillRoot: skillRoot,
		UpstreamURL:        upstreamURL,
		ToolDirectory:      toolDirectory,
		Narration:          &narration,
		ConfirmationInput:  strings.NewReader(confirmation),
	})
	return result, narration.String()
}

func TestAnUpToDateSuiteReportsSkippedWorkWithoutInstalling(t *testing.T) {
	projectRoot := newProjectRepository(t)
	skillRoot := installSuiteForUpdate(t, projectRoot, fixtureSuiteVersion)
	toolDirectory, upstreamURL := newUpstreamArchiveServer(t, fixtureSuiteVersion)
	settingsBefore := readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json"))

	result, narration := runUpdateFixture(t, projectRoot, skillRoot, toolDirectory, upstreamURL, "y\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, reason = %q\n%s", result.Outcome, result.FailureReason, narration)
	}
	if len(result.SkippedWork) != 1 || result.SkippedWork[0].Code != SkipCodeUpdateAlreadyCurrent {
		t.Fatalf("skipped work = %#v, want one %s", result.SkippedWork, SkipCodeUpdateAlreadyCurrent)
	}
	if !strings.Contains(narration, "You're up to date (v"+fixtureSuiteVersion+")") {
		t.Errorf("narration does not report the up-to-date state:\n%s", narration)
	}
	if strings.Contains(narration, "Install this complete four-skill suite?") {
		t.Errorf("an up-to-date check asked for an install confirmation:\n%s", narration)
	}
	if readTestFile(t, filepath.Join(projectRoot, ".claude", "settings.json")) != settingsBefore {
		t.Errorf("an up-to-date check rewrote settings.json")
	}
}

func TestAnOlderUpstreamIsRefused(t *testing.T) {
	projectRoot := newProjectRepository(t)
	skillRoot := installSuiteForUpdate(t, projectRoot, "0.300.0")
	toolDirectory, upstreamURL := newUpstreamArchiveServer(t, "0.200.0")

	result, _ := runUpdateFixture(t, projectRoot, skillRoot, toolDirectory, upstreamURL, "y\n")
	if result.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("outcome = %q, want failure", result.Outcome)
	}
	if !strings.Contains(result.FailureReason, "upstream version v0.200.0 is older than installed v0.300.0") {
		t.Errorf("failure reason = %q", result.FailureReason)
	}
}

// The updater no longer runs the installer as a subprocess, so a declined confirmation comes
// back as a value rather than through DO_WORK_INSTALL_CANCEL_EXIT_STATUS. That env var is
// gone, and this is the behaviour that replaced it.
func TestACancelledUpdateIsSuccessWithSkippedWorkAndNoEnvironmentVariable(t *testing.T) {
	projectRoot := newProjectRepository(t)
	skillRoot := installSuiteForUpdate(t, projectRoot, "0.100.0")
	toolDirectory, upstreamURL := newUpstreamArchiveServer(t, "0.200.0")
	installedVersionBefore := readTestFile(t, filepath.Join(skillRoot, "VERSION"))

	result, narration := runUpdateFixture(t, projectRoot, skillRoot, toolDirectory, upstreamURL, "n\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, reason = %q\n%s", result.Outcome, result.FailureReason, narration)
	}
	if resultmodel.ExitCode(result.Outcome) != 0 {
		t.Errorf("a cancelled update exits %d, want 0", resultmodel.ExitCode(result.Outcome))
	}
	if len(result.SkippedWork) != 1 || result.SkippedWork[0].Code != SkipCodeUpdateCancelled {
		t.Fatalf("skipped work = %#v, want one %s", result.SkippedWork, SkipCodeUpdateCancelled)
	}
	if !strings.Contains(narration, "Update cancelled; no files were changed.") {
		t.Errorf("narration does not report the cancellation:\n%s", narration)
	}
	if readTestFile(t, filepath.Join(skillRoot, "VERSION")) != installedVersionBefore {
		t.Errorf("a cancelled update changed the installed version")
	}
}

func TestAnAvailableUpdateInstallsAndVerifiesTheNewVersion(t *testing.T) {
	projectRoot := newProjectRepository(t)
	skillRoot := installSuiteForUpdate(t, projectRoot, "0.100.0")
	toolDirectory, upstreamURL := newUpstreamArchiveServer(t, "0.200.0")

	result, narration := runUpdateFixture(t, projectRoot, skillRoot, toolDirectory, upstreamURL, "y\n")
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, reason = %q\n%s", result.Outcome, result.FailureReason, narration)
	}
	if result.RemoteVersion != "0.200.0" || result.LocalVersion != "0.100.0" {
		t.Errorf("versions = local %q remote %q", result.LocalVersion, result.RemoteVersion)
	}
	if installed := readTestFile(t, filepath.Join(skillRoot, "VERSION")); installed != "0.200.0\n" {
		t.Errorf("installed VERSION = %q, want the upstream version", installed)
	}
	if !strings.Contains(narration, "Updated to v0.200.0 at "+projectRoot) {
		t.Errorf("narration does not report the completed update:\n%s", narration)
	}
}

// The updater refuses to touch a shared install: the skill it is updating has to live inside
// the project it was pointed at.
func TestASkillOutsideTheProjectIsRefused(t *testing.T) {
	projectRoot := newProjectRepository(t)
	sharedSkillRoot := filepath.Join(t.TempDir(), "shared", "do-work")
	writeTestFile(t, filepath.Join(sharedSkillRoot, "SKILL.md"), "# do-work\n")
	writeTestFile(t, filepath.Join(sharedSkillRoot, "actions", "version.md"), "**Current version**: 0.100.0\n")
	toolDirectory, upstreamURL := newUpstreamArchiveServer(t, "0.200.0")

	result, _ := runUpdateFixture(t, projectRoot, sharedSkillRoot, toolDirectory, upstreamURL, "y\n")
	if result.Outcome != resultmodel.OutcomeFailure {
		t.Fatalf("outcome = %q, want failure", result.Outcome)
	}
	if !strings.Contains(result.FailureReason, "refusing to update a shared install") {
		t.Errorf("failure reason = %q", result.FailureReason)
	}
}

func TestSemanticVersionComparisonOrdersEveryComponent(t *testing.T) {
	tests := []struct {
		localVersion  string
		remoteVersion string
		expected      int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "1.2.4", 1},
		{"1.2.4", "1.2.3", -1},
		{"1.9.0", "1.10.0", 1},
		{"2.0.0", "1.99.99", -1},
		{"0.100.0", "0.200.0", 1},
	}
	for _, test := range tests {
		if actual := compareSemanticVersions(test.localVersion, test.remoteVersion); actual != test.expected {
			t.Errorf("compareSemanticVersions(%q, %q) = %d, want %d",
				test.localVersion, test.remoteVersion, actual, test.expected)
		}
	}
}
