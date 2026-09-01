package archivefetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
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

func TestDownloadAtomicRetriesThreeTimesAfterInitial429AndUsesTokenPrecedence(t *testing.T) {
	setAtomicRetryTiming(t, time.Millisecond, time.Second)
	directory := t.TempDir()
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer preferred" {
			t.Errorf("authorization=%q", request.Header.Get("Authorization"))
		}
		if attempts.Add(1) < 4 {
			response.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = response.Write([]byte("complete body"))
	}))
	defer server.Close()
	t.Setenv("GH_TOKEN", "preferred")
	t.Setenv("GITHUB_TOKEN", "fallback")
	target := filepath.Join(directory, "download.bin")
	result := DownloadAtomic(context.Background(), server.URL, target)
	if result.Err != nil {
		t.Fatal(result.Err)
	}
	if result.Attempts != 4 {
		t.Fatalf("attempts=%d", result.Attempts)
	}
	contents, _ := os.ReadFile(target)
	if string(contents) != "complete body" {
		t.Fatalf("contents=%q", contents)
	}
	info, _ := os.Stat(target)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%o", info.Mode().Perm())
	}
}

func TestDownloadAtomicUsesFixedRetryDelayAndNoWholeRequestTimeout(t *testing.T) {
	if atomicHTTPClient.Timeout != 0 {
		t.Fatalf("HTTP client imposes a shorter whole-request timeout: %s", atomicHTTPClient.Timeout)
	}
	setAtomicRetryTiming(t, 20*time.Millisecond, time.Second)
	requestTimes := []time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		requestTimes = append(requestTimes, time.Now())
		if len(requestTimes) < 3 {
			response.WriteHeader(http.StatusTooManyRequests)
			return
		}
		for index := 0; index < 4; index++ {
			_, _ = response.Write(make([]byte, 256*1024))
			time.Sleep(5 * time.Millisecond)
		}
	}))
	defer server.Close()
	target := filepath.Join(t.TempDir(), "streamed.bin")
	result := DownloadAtomic(context.Background(), server.URL, target)
	if result.Err != nil || result.Attempts != 3 || result.BytesWritten != 1024*1024 {
		t.Fatalf("result=%#v", result)
	}
	for index := 1; index < len(requestTimes); index++ {
		if delay := requestTimes[index].Sub(requestTimes[index-1]); delay < 15*time.Millisecond {
			t.Fatalf("retry delay %s is shorter than configured fixed delay", delay)
		}
	}
}

func TestDownloadAtomicNeverLeaksTokenInFailure(t *testing.T) {
	t.Setenv("GH_TOKEN", "top-secret-token")
	result := DownloadAtomic(context.Background(), "://top-secret-token", filepath.Join(t.TempDir(), "target"))
	if result.Err == nil {
		t.Fatal("invalid URL succeeded")
	}
	if strings.Contains(result.Err.Error(), "top-secret-token") {
		t.Fatalf("token leaked: %v", result.Err)
	}
}

func TestDownloadAtomicParentSwapCannotOverwriteOutsideTarget(t *testing.T) {
	directory := t.TempDir()
	parent := filepath.Join(directory, "parent")
	held := filepath.Join(directory, "held")
	outside := filepath.Join(directory, "outside")
	_ = os.Mkdir(parent, 0o755)
	_ = os.Mkdir(outside, 0o755)
	protected := filepath.Join(outside, "target")
	_ = os.WriteFile(protected, []byte("protected"), 0o600)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write([]byte("download")) }))
	defer server.Close()
	originalHook := beforeAtomicDownloadPublish
	defer func() { beforeAtomicDownloadPublish = originalHook }()
	beforeAtomicDownloadPublish = func() {
		if err := os.Rename(parent, held); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, parent); err != nil {
			t.Fatal(err)
		}
	}
	result := DownloadAtomic(context.Background(), server.URL, filepath.Join(parent, "target"))
	if result.Err != nil {
		t.Fatal(result.Err)
	}
	if contents, _ := os.ReadFile(protected); string(contents) != "protected" {
		t.Fatalf("outside target changed: %q", contents)
	}
	if contents, _ := os.ReadFile(filepath.Join(held, "target")); string(contents) != "download" {
		t.Fatalf("rooted download missing: %q", contents)
	}
}

func TestHttpRouteWinsWhenTheAtomicPrimitiveProducesAReadableArchive(t *testing.T) {
	directory := t.TempDir()
	sourceArchive := newFixtureArchive(t, directory, "payload.txt")
	archiveBytes, readError := os.ReadFile(sourceArchive)
	if readError != nil {
		t.Fatal(readError)
	}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write(archiveBytes) }))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")

	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:  targetPath,
		UpstreamTarballURL: server.URL,
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
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write([]byte("rate limited\n")) }))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")

	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    server.URL,
		UpstreamRepositoryURL: repositoryPath,
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
	_, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "https://example.com/suite.tar.gz",
		UpstreamRepositoryURL: filepath.Join(directory, "no-such-repository"),
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

func TestGitFallbackNeverOverwritesAnOccupiedTarget(t *testing.T) {
	directory := t.TempDir()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("owned bytes\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "://invalid",
		UpstreamRepositoryURL: newFixtureRepository(t, directory),
	})
	if err == nil {
		t.Fatal("occupied publication unexpectedly succeeded")
	}
	contents, readErr := os.ReadFile(targetPath)
	if readErr != nil || string(contents) != "owned bytes\n" {
		t.Fatalf("target=%q err=%v", contents, readErr)
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
	if !strings.Contains(err.Error(), "HTTP route: failed") {
		t.Errorf("failure report %q does not name the failed in-process HTTP route", err.Error())
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

func setAtomicRetryTiming(t *testing.T, delay, budget time.Duration) {
	t.Helper()
	previousDelay, previousBudget := atomicRetryDelay, atomicRetryBudget
	atomicRetryDelay, atomicRetryBudget = delay, budget
	t.Cleanup(func() {
		atomicRetryDelay, atomicRetryBudget = previousDelay, previousBudget
	})
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
