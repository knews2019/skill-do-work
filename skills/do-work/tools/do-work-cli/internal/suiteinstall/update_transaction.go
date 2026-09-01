package suiteinstall

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/archivefetch"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/suitemanifest"
)

// SkipCodeUpdateAlreadyCurrent marks an update that had nothing to do.
const SkipCodeUpdateAlreadyCurrent = "UPDATE-ALREADY-CURRENT"

// SkipCodeUpdateCancelled marks an update whose install confirmation was declined.
const SkipCodeUpdateCancelled = "UPDATE-CANCELLED"

// UpdateOptions names what an update needs. The installed suite is discovered from the
// project root rather than supplied, because the updater refuses to touch a shared install.
type UpdateOptions struct {
	ProjectRoot        string
	InstalledSkillRoot string
	UpstreamURL        string
	ToolDirectory      string
	Narration          io.Writer
	ConfirmationInput  io.Reader
}

// UpdateResult reports the update in resultmodel's vocabulary. It reuses the install
// transaction's result rather than translating it twice.
type UpdateResult struct {
	Outcome       resultmodel.CommandOutcome
	LocalVersion  string
	RemoteVersion string
	Changes       []resultmodel.RecordedChange
	SkippedWork   []resultmodel.SkippedWork
	Rollback      resultmodel.RollbackResult
	FailurePaths  []string
	FailureReason string
}

// RunUpdate fetches, extracts and validates the upstream suite once, compares versions, and
// hands the already-extracted tree to the install transaction in-process. Because the
// installer is no longer a subprocess, cancellation comes back as a value rather than as an
// exit status, which is what retires DO_WORK_INSTALL_CANCEL_EXIT_STATUS.
func RunUpdate(ctx context.Context, options UpdateOptions) UpdateResult {
	narrate := func(format string, arguments ...any) {
		if options.Narration != nil {
			fmt.Fprintf(options.Narration, format+"\n", arguments...)
		}
	}

	projectRoot, skillRoot, err := resolveUpdateRoots(options)
	if err != nil {
		return updateFailure("%s", err.Error())
	}
	localVersion, err := readInstalledVersion(filepath.Join(skillRoot, "actions", "version.md"))
	if err != nil {
		return updateFailure("could not read a semantic local version")
	}

	updateTmp, err := os.MkdirTemp("", "do-work-update.*")
	if err != nil {
		return updateFailure("could not allocate a private update workspace")
	}
	defer func() { _ = os.RemoveAll(updateTmp) }()
	upstreamTarball := filepath.Join(updateTmp, "upstream.tar.gz")
	freshUpstream := filepath.Join(updateTmp, "fresh")
	if err := os.MkdirAll(freshUpstream, 0o755); err != nil {
		return updateFailure("could not allocate a private update workspace")
	}

	narrate("Checking do-work updates…")
	upstreamURL := options.UpstreamURL
	if upstreamURL == "" {
		upstreamURL = archivefetch.UpstreamURLFromEnvironment()
	}
	if _, err := archivefetch.FetchArchive(ctx, archivefetch.Request{
		ArchiveTargetPath:  upstreamTarball,
		UpstreamTarballURL: upstreamURL,
	}); err != nil {
		return updateFailure("upstream archive could not be fetched by any route; no files were changed")
	}
	if err := runQuietly(ctx, "", "tar", "xzf", upstreamTarball, "-C", freshUpstream, "--strip-components=1"); err != nil {
		return updateFailure("upstream archive could not be extracted; no files were changed")
	}
	validation, err := suitemanifest.ValidateSuite(freshUpstream)
	if err != nil {
		return updateFailure("suite manifest validation failed; no files were changed: %v", err)
	}
	remoteVersion := validation.SuiteVersion

	switch compareSemanticVersions(localVersion, remoteVersion) {
	case 0:
		narrate("You're up to date (v%s)", localVersion)
		return UpdateResult{
			Outcome:       resultmodel.OutcomeSuccess,
			LocalVersion:  localVersion,
			RemoteVersion: remoteVersion,
			Rollback:      resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
			SkippedWork: []resultmodel.SkippedWork{{
				Code:   SkipCodeUpdateAlreadyCurrent,
				Reason: "the installed suite is already v" + localVersion,
			}},
		}
	case -1:
		return updateFailure("upstream version v%s is older than installed v%s", remoteVersion, localVersion)
	}

	narrate("Update available: v%s (you have v%s), archive layout: four-module suite.", remoteVersion, localVersion)
	installResult := RunInstall(ctx, InstallOptions{
		ProjectRoot:         projectRoot,
		ExtractedSourceRoot: freshUpstream,
		ToolDirectory:       options.ToolDirectory,
		Narration:           options.Narration,
		ConfirmationInput:   options.ConfirmationInput,
	})
	if installResult.Cancelled {
		narrate("Update cancelled; no files were changed.")
		return UpdateResult{
			Outcome:       resultmodel.OutcomeSuccess,
			LocalVersion:  localVersion,
			RemoteVersion: remoteVersion,
			Rollback:      resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
			SkippedWork: []resultmodel.SkippedWork{{
				Code:   SkipCodeUpdateCancelled,
				Reason: "the single install confirmation was declined; no files were changed",
			}},
		}
	}
	if installResult.Outcome != resultmodel.OutcomeSuccess {
		return UpdateResult{
			Outcome:       installResult.Outcome,
			LocalVersion:  localVersion,
			RemoteVersion: remoteVersion,
			Changes:       installResult.Changes,
			Rollback:      installResult.Rollback,
			FailurePaths:  installResult.FailurePaths,
			FailureReason: "full-suite installation failed; managed paths were recovered: " + installResult.FailureReason,
		}
	}

	installedVersion, err := readInstalledVersion(
		filepath.Join(projectRoot, ".claude", "skills", "do-work", "actions", "version.md"))
	if err != nil || installedVersion != remoteVersion {
		return updateFailure("post-update version verification failed (expected v%s, found v%s)",
			remoteVersion, orUnknown(installedVersion))
	}
	narrate("Updated to v%s at %s using the four-module suite.", remoteVersion, projectRoot)
	return UpdateResult{
		Outcome:       resultmodel.OutcomeSuccess,
		LocalVersion:  localVersion,
		RemoteVersion: remoteVersion,
		Changes:       installResult.Changes,
		Rollback:      resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
	}
}

// resolveUpdateRoots enforces the two guards that keep an update project-local: the project
// must be its own Git worktree root, and the skill being updated must live inside it, so a
// shared install can never be clobbered from a consumer project.
func resolveUpdateRoots(options UpdateOptions) (projectRoot string, skillRoot string, err error) {
	if options.ProjectRoot == "" {
		return "", "", fmt.Errorf("usage: do-work-update.sh --project-root <project-root>")
	}
	info, statErr := os.Stat(options.ProjectRoot)
	if statErr != nil || !info.IsDir() {
		return "", "", fmt.Errorf("project root does not exist: %s", options.ProjectRoot)
	}
	projectRoot, err = physicalPath(options.ProjectRoot)
	if err != nil {
		return "", "", fmt.Errorf("project root does not exist: %s", options.ProjectRoot)
	}
	skillRoot, err = physicalPath(options.InstalledSkillRoot)
	if err != nil {
		return "", "", fmt.Errorf("the installed do-work skill could not be located")
	}
	if !strings.HasPrefix(skillRoot, projectRoot+string(filepath.Separator)) {
		return "", "", fmt.Errorf(
			"skill is outside this project (%s is not within %s); refusing to update a shared install",
			skillRoot, projectRoot)
	}
	if err := requireNonEmptyRegularFile(filepath.Join(skillRoot, "SKILL.md")); err != nil {
		return "", "", fmt.Errorf("SKILL.md is missing at %s", skillRoot)
	}
	if info, statErr := os.Stat(filepath.Join(skillRoot, "actions", "version.md")); statErr != nil || !info.Mode().IsRegular() {
		return "", "", fmt.Errorf("actions/version.md is missing at %s", skillRoot)
	}

	gitRoot, gitErr := runCapturing(context.Background(), projectRoot, "git", "rev-parse", "--show-toplevel")
	if gitErr != nil || strings.TrimSpace(gitRoot) == "" {
		return "", "", fmt.Errorf("the project must be a Git repository so a failed suite update can be recovered")
	}
	physicalGitRoot, gitPathErr := physicalPath(strings.TrimSpace(gitRoot))
	if gitPathErr != nil {
		return "", "", fmt.Errorf("the project must be a Git repository so a failed suite update can be recovered")
	}
	if physicalGitRoot != projectRoot {
		return "", "", fmt.Errorf("--project-root must name the Git worktree root (%s)", physicalGitRoot)
	}
	return projectRoot, skillRoot, nil
}

func physicalPath(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return filepath.Abs(resolved)
}

// readInstalledVersion requires a plain semantic version, which is the check both the shell
// updater and the manifest validator applied before trusting a version to compare.
func readInstalledVersion(actionVersionPath string) (string, error) {
	version, markerCount, err := suitemanifest.ReadActionVersion(actionVersionPath)
	if err != nil || markerCount == 0 {
		return "", fmt.Errorf("could not read a Current version line from %s", actionVersionPath)
	}
	if !suitemanifest.IsSemanticVersion(version) {
		return "", fmt.Errorf("%s does not carry a plain semantic version", actionVersionPath)
	}
	return version, nil
}

// compareSemanticVersions returns 1 when remote is newer, -1 when it is older, 0 when the
// three numeric components match.
func compareSemanticVersions(localVersion, remoteVersion string) int {
	localParts := strings.SplitN(localVersion, ".", 3)
	remoteParts := strings.SplitN(remoteVersion, ".", 3)
	for component := 0; component < 3; component++ {
		localNumber := numericComponent(localParts, component)
		remoteNumber := numericComponent(remoteParts, component)
		if remoteNumber > localNumber {
			return 1
		}
		if remoteNumber < localNumber {
			return -1
		}
	}
	return 0
}

func numericComponent(parts []string, index int) int {
	if index >= len(parts) {
		return 0
	}
	number, err := strconv.Atoi(parts[index])
	if err != nil {
		return 0
	}
	return number
}

func updateFailure(format string, arguments ...any) UpdateResult {
	return UpdateResult{
		Outcome:       resultmodel.OutcomeFailure,
		Rollback:      resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
		FailureReason: fmt.Sprintf(format, arguments...),
	}
}
