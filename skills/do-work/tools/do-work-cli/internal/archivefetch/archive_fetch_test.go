package archivefetch

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The branch derivation decides which ref the git route clones. Substituting HEAD for a
// named branch that no longer exists would be a silent false success, so a URL shape the
// derivation does not recognise must produce no repository URL at all.
func TestGitRouteDerivationOnlyReadsBranchTarballUrls(t *testing.T) {
	tests := []struct {
		name               string
		tarballURL         string
		suppliedRepository string
		expectedBranch     string
		expectedRepository string
	}{
		{
			name:               "a GitHub branch tarball",
			tarballURL:         "https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz",
			expectedBranch:     "main",
			expectedRepository: "https://github.com/knews2019/skill-do-work.git",
		},
		{
			name:               "a branch name carrying a slash",
			tarballURL:         "https://github.com/owner/repo/archive/refs/heads/release/2.x.tar.gz",
			expectedBranch:     "release/2.x",
			expectedRepository: "https://github.com/owner/repo.git",
		},
		{
			name:               "an explicit repository URL is not overridden",
			tarballURL:         "https://github.com/owner/repo/archive/refs/heads/main.tar.gz",
			suppliedRepository: "git@example.com:owner/repo.git",
			expectedBranch:     "main",
			expectedRepository: "git@example.com:owner/repo.git",
		},
		{
			name:               "an unrecognised URL shape derives nothing",
			tarballURL:         "https://example.com/downloads/suite.tar.gz",
			expectedBranch:     "",
			expectedRepository: "",
		},
		{
			name:               "a branch path without the tarball suffix derives nothing",
			tarballURL:         "https://github.com/owner/repo/archive/refs/heads/main.zip",
			expectedBranch:     "",
			expectedRepository: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			branch, repository := deriveGitRoute(test.tarballURL, test.suppliedRepository)
			if branch != test.expectedBranch {
				t.Errorf("branch = %q, want %q", branch, test.expectedBranch)
			}
			if repository != test.expectedRepository {
				t.Errorf("repository URL = %q, want %q", repository, test.expectedRepository)
			}
		})
	}
}

func TestHttpRouteWinsWhenTheAtomicPrimitiveProducesAReadableArchive(t *testing.T) {
	directory := t.TempDir()
	sourceArchive := newFixtureArchive(t, directory, "payload.txt")
	downloadScript := writeExecutableScript(t, filepath.Join(directory, "atomic-download.sh"),
		"#!/usr/bin/env bash\nset -eu\ncp \""+sourceArchive+"\" \"$2\"\n")
	targetPath := filepath.Join(directory, "upstream.tar.gz")

	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:    targetPath,
		UpstreamTarballURL:   "https://example.com/suite.tar.gz",
		AtomicDownloadScript: downloadScript,
	})
	if err != nil {
		t.Fatalf("FetchArchive: %v", err)
	}
	if result.RouteDescription != "upstream archive fetched over HTTP" {
		t.Errorf("route = %q, want the HTTP route", result.RouteDescription)
	}
	if _, statErr := os.Stat(targetPath); statErr != nil {
		t.Errorf("the archive was not published: %v", statErr)
	}
}

// An unreadable body — a rate-limit HTML page, for one — must not be published as an archive.
func TestAnUnreadableHttpDownloadFallsThroughToTheGitRoute(t *testing.T) {
	directory := t.TempDir()
	repositoryPath := newFixtureRepository(t, directory)
	downloadScript := writeExecutableScript(t, filepath.Join(directory, "atomic-download.sh"),
		"#!/usr/bin/env bash\nset -eu\nprintf 'rate limited\\n' > \"$2\"\n")
	targetPath := filepath.Join(directory, "upstream.tar.gz")

	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "https://example.com/suite.tar.gz",
		UpstreamRepositoryURL: repositoryPath,
		AtomicDownloadScript:  downloadScript,
	})
	if err != nil {
		t.Fatalf("FetchArchive: %v", err)
	}
	if !strings.HasPrefix(result.RouteDescription, "upstream archive fetched with git (HTTP route failed") {
		t.Errorf("route = %q, want the git route naming the failed HTTP route", result.RouteDescription)
	}
	if !archiveIsReadable(context.Background(), targetPath) {
		t.Errorf("the published archive is not a readable tarball")
	}
	if names := archiveEntryNames(t, targetPath); !strings.Contains(names, gitArchivePrefix+"tracked.txt") {
		t.Errorf("archive entries = %q, want the tracked file under the Git prefix", names)
	}
}

// git archive honours export-ignore; a worktree copy does not. That is why the git route
// repacks rather than tars the clone.
func TestGitRouteHonoursExportIgnore(t *testing.T) {
	directory := t.TempDir()
	repositoryPath := newFixtureRepository(t, directory)
	writeRepositoryFile(t, repositoryPath, ".gitattributes", "maintainer-only.txt export-ignore\n")
	writeRepositoryFile(t, repositoryPath, "maintainer-only.txt", "not for consumers\n")
	runFixtureGit(t, repositoryPath, "add", "-A")
	runFixtureGit(t, repositoryPath, "commit", "-qm", "export-ignore fixture")
	targetPath := filepath.Join(directory, "upstream.tar.gz")

	if _, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "https://example.com/suite.tar.gz",
		UpstreamRepositoryURL: repositoryPath,
	}); err != nil {
		t.Fatalf("FetchArchive: %v", err)
	}
	names := archiveEntryNames(t, targetPath)
	if strings.Contains(names, "maintainer-only.txt") {
		t.Errorf("an export-ignored file reached the archive: %s", names)
	}
	if !strings.Contains(names, "tracked.txt") {
		t.Errorf("archive is missing its tracked content: %s", names)
	}
}

// Total failure must leave a pre-existing target byte-identical and drop no scratch file
// beside it, so a failed update cannot destroy an archive an operator already had.
func TestTotalFailurePreservesTheTargetAndLeavesNoScratch(t *testing.T) {
	directory := t.TempDir()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("existing archive bytes\n"), 0o644); err != nil {
		t.Fatalf("seed target: %v", err)
	}
	downloadScript := writeExecutableScript(t, filepath.Join(directory, "atomic-download.sh"),
		"#!/usr/bin/env bash\nexit 22\n")

	_, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "https://example.com/suite.tar.gz",
		UpstreamRepositoryURL: filepath.Join(directory, "no-such-repository"),
		AtomicDownloadScript:  downloadScript,
	})
	if err == nil {
		t.Fatalf("a total failure reported success")
	}
	for _, expectedFragment := range []string{
		"upstream archive could not be fetched",
		"HTTP route: failed",
		"Git route: failed",
		"DO_WORK_UPSTREAM_URL",
	} {
		if !strings.Contains(err.Error(), expectedFragment) {
			t.Errorf("failure report %q is missing %q", err.Error(), expectedFragment)
		}
	}
	preserved, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		t.Fatalf("read preserved target: %v", readErr)
	}
	if string(preserved) != "existing archive bytes\n" {
		t.Errorf("the pre-existing target was overwritten: %q", preserved)
	}
	entries, dirErr := os.ReadDir(directory)
	if dirErr != nil {
		t.Fatalf("read directory: %v", dirErr)
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".fetching.") {
			t.Errorf("a scratch file survived a total failure: %s", entry.Name())
		}
	}
}

// A repository URL that cannot be derived must report the route as unavailable rather than
// attempting a clone of nothing.
func TestUndeivableRepositoryUrlReportsAnUnavailableGitRoute(t *testing.T) {
	directory := t.TempDir()
	_, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:  filepath.Join(directory, "upstream.tar.gz"),
		UpstreamTarballURL: "https://example.com/downloads/suite.tar.gz",
	})
	if err == nil {
		t.Fatalf("a fetch with no route reported success")
	}
	if !strings.Contains(err.Error(), "Git route: unavailable (no repository URL supplied and none derivable") {
		t.Errorf("failure report %q does not name the underivable git route", err.Error())
	}
	if !strings.Contains(err.Error(), "HTTP route: unavailable (atomic-download.sh not found") {
		t.Errorf("failure report %q does not name the missing HTTP primitive", err.Error())
	}
}

func TestUpstreamUrlFromEnvironmentHonoursTheEscapeHatch(t *testing.T) {
	t.Setenv("DO_WORK_UPSTREAM_URL", "")
	if url := UpstreamURLFromEnvironment(); url != DefaultUpstreamURL {
		t.Errorf("default URL = %q, want %q", url, DefaultUpstreamURL)
	}
	t.Setenv("DO_WORK_UPSTREAM_URL", "https://mirror.example.com/suite.tar.gz")
	if url := UpstreamURLFromEnvironment(); url != "https://mirror.example.com/suite.tar.gz" {
		t.Errorf("overridden URL = %q", url)
	}
}

func newFixtureArchive(t *testing.T, directory, entryName string) string {
	t.Helper()
	contentDirectory := filepath.Join(directory, "archive-content")
	if err := os.MkdirAll(contentDirectory, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(contentDirectory, entryName), []byte("payload\n"), 0o644); err != nil {
		t.Fatalf("write archive entry: %v", err)
	}
	archivePath := filepath.Join(directory, "fixture.tar.gz")
	command := exec.Command("tar", "czf", archivePath, "-C", contentDirectory, entryName)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build fixture archive: %v: %s", err, output)
	}
	return archivePath
}

func newFixtureRepository(t *testing.T, directory string) string {
	t.Helper()
	repositoryPath := filepath.Join(directory, "upstream-repository")
	if err := os.MkdirAll(repositoryPath, 0o755); err != nil {
		t.Fatalf("mkdir repository: %v", err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(directory, "gitconfig"))
	t.Setenv("GIT_CONFIG_SYSTEM", filepath.Join(directory, "gitconfig-system"))
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
	runFixtureGit(t, directory, "init", "-q", repositoryPath)
	runFixtureGit(t, repositoryPath, "config", "user.email", "fixture@example.com")
	runFixtureGit(t, repositoryPath, "config", "user.name", "Archive Fetch Fixture")
	writeRepositoryFile(t, repositoryPath, "tracked.txt", "tracked content\n")
	runFixtureGit(t, repositoryPath, "add", "-A")
	runFixtureGit(t, repositoryPath, "commit", "-qm", "fixture")
	return repositoryPath
}

func writeRepositoryFile(t *testing.T, repositoryPath, relativePath, contents string) {
	t.Helper()
	fullPath := filepath.Join(repositoryPath, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func runFixtureGit(t *testing.T, workingDirectory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = workingDirectory
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(arguments, " "), err, output)
	}
}

func writeExecutableScript(t *testing.T, path, contents string) string {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o755); err != nil {
		t.Fatalf("write script %s: %v", path, err)
	}
	return path
}

func archiveEntryNames(t *testing.T, archivePath string) string {
	t.Helper()
	output, err := exec.Command("tar", "tzf", archivePath).CombinedOutput()
	if err != nil {
		t.Fatalf("list archive: %v: %s", err, output)
	}
	return string(output)
}
