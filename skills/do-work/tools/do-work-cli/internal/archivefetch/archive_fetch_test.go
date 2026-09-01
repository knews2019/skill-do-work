package archivefetch

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
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

func TestDownloadAtomicStillRefusesAnOccupiedTarget(t *testing.T) {
	directory := t.TempDir()
	targetPath := filepath.Join(directory, "download.bin")
	if err := os.WriteFile(targetPath, []byte("owned bytes\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		requestCount.Add(1)
		_, _ = response.Write([]byte("replacement"))
	}))
	defer server.Close()

	result := DownloadAtomic(context.Background(), server.URL, targetPath)
	if result.Err == nil {
		t.Fatal("occupied generic download target was replaced")
	}
	if requestCount.Load() != 0 {
		t.Fatalf("download made %d request(s) after the occupied-target refusal", requestCount.Load())
	}
	contents, readError := os.ReadFile(targetPath)
	if readError != nil || string(contents) != "owned bytes\n" {
		t.Fatalf("target=%q err=%v", contents, readError)
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

func TestHttpRouteReplacesAnUnchangedRegularArchiveAfterValidation(t *testing.T) {
	directory := t.TempDir()
	sourceArchive := newFixtureArchive(t, directory, "replacement.txt")
	archiveBytes, readError := os.ReadFile(sourceArchive)
	if readError != nil {
		t.Fatal(readError)
	}
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = response.Write(archiveBytes)
	}))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("old archive bytes\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:  targetPath,
		UpstreamTarballURL: server.URL,
	})
	if err != nil {
		t.Fatalf("FetchArchive: %v", err)
	}
	if result.RouteDescription != "upstream archive fetched over HTTP" {
		t.Fatalf("route = %q, want the HTTP route", result.RouteDescription)
	}
	contents, readError := os.ReadFile(targetPath)
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(contents, archiveBytes) {
		t.Fatal("HTTP route did not publish the validated replacement")
	}
	info, statError := os.Stat(targetPath)
	if statError != nil {
		t.Fatal(statError)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%o, want private candidate mode 600", info.Mode().Perm())
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

func TestGitFallbackSlashBranchExtractsAtRootAfterOneComponentStrip(t *testing.T) {
	directory := t.TempDir()
	repositoryPath := newFixtureRepository(t, directory)
	runFixtureGit(t, repositoryPath, "switch", "-q", "-c", "release/2.x")
	writeRepositoryFile(t, repositoryPath, "VERSION", "release-2.x\n")
	runFixtureGit(t, repositoryPath, "add", "VERSION")
	runFixtureGit(t, repositoryPath, "commit", "-qm", "slash branch fixture")

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = response.Write([]byte("rate limited\n"))
	}))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    server.URL + "/owner/repo/archive/refs/heads/release/2.x.tar.gz",
		UpstreamRepositoryURL: repositoryPath,
	})
	if err != nil {
		t.Fatalf("FetchArchive: %v", err)
	}
	if !strings.HasPrefix(result.RouteDescription, "upstream archive fetched with git (HTTP route failed") {
		t.Fatalf("route = %q, want Git after failed HTTP", result.RouteDescription)
	}

	extractionRoot := filepath.Join(directory, "extracted")
	if err := os.Mkdir(extractionRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	extract := exec.Command("tar", "xzf", targetPath, "--strip-components=1", "-C", extractionRoot)
	if output, err := extract.CombinedOutput(); err != nil {
		t.Fatalf("extract Git fallback archive: %v: %s", err, output)
	}
	version, err := os.ReadFile(filepath.Join(extractionRoot, "VERSION"))
	if err != nil {
		t.Fatalf("read root VERSION after one-component extraction: %v", err)
	}
	if string(version) != "release-2.x\n" {
		t.Fatalf("root VERSION = %q, want exact slash branch content", version)
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
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		_, _ = response.Write([]byte("not an archive\n"))
	}))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("existing archive bytes\n"), 0o644); err != nil {
		t.Fatalf("seed target: %v", err)
	}
	_, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    server.URL,
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
	preservedInfo, statError := os.Stat(targetPath)
	if statError != nil {
		t.Fatal(statError)
	}
	if preservedInfo.Mode().Perm() != 0o644 {
		t.Errorf("the pre-existing target mode changed: %o", preservedInfo.Mode().Perm())
	}
	assertNoArchiveScratch(t, directory)
}

func TestFetchArchiveRefusesUnsafeTargetsUnchanged(t *testing.T) {
	t.Run("symlink", func(t *testing.T) {
		directory := t.TempDir()
		protectedPath := filepath.Join(directory, "protected.tar.gz")
		if err := os.WriteFile(protectedPath, []byte("protected\n"), 0o640); err != nil {
			t.Fatal(err)
		}
		targetPath := filepath.Join(directory, "upstream.tar.gz")
		if err := os.Symlink(protectedPath, targetPath); err != nil {
			t.Fatal(err)
		}
		_, fetchError := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: "://invalid"})
		if fetchError == nil {
			t.Fatal("symlink target was accepted")
		}
		for _, expected := range []string{"HTTP route: not attempted", "Git route: not attempted", "archive target is unsafe"} {
			if !strings.Contains(fetchError.Error(), expected) {
				t.Errorf("failure %q is missing %q", fetchError, expected)
			}
		}
		contents, readError := os.ReadFile(protectedPath)
		if readError != nil || string(contents) != "protected\n" {
			t.Fatalf("protected target=%q err=%v", contents, readError)
		}
		if info, err := os.Lstat(targetPath); err != nil || info.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("symlink target changed: info=%v err=%v", info, err)
		}
		assertNoArchiveScratch(t, directory)
	})

	t.Run("directory", func(t *testing.T) {
		directory := t.TempDir()
		targetPath := filepath.Join(directory, "upstream.tar.gz")
		if err := os.Mkdir(targetPath, 0o750); err != nil {
			t.Fatal(err)
		}
		if _, err := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: "://invalid"}); err == nil {
			t.Fatal("directory target was accepted")
		}
		if info, err := os.Lstat(targetPath); err != nil || !info.IsDir() {
			t.Fatalf("directory target changed: info=%v err=%v", info, err)
		}
		assertNoArchiveScratch(t, directory)
	})

	t.Run("special file", func(t *testing.T) {
		directory := t.TempDir()
		targetPath := filepath.Join(directory, "upstream.tar.gz")
		listener, listenError := net.Listen("unix", targetPath)
		if listenError != nil {
			t.Skipf("Unix sockets unavailable: %v", listenError)
		}
		defer listener.Close()
		if _, err := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: "://invalid"}); err == nil {
			t.Fatal("special-file target was accepted")
		}
		if info, err := os.Lstat(targetPath); err != nil || info.Mode()&os.ModeSocket == 0 {
			t.Fatalf("special-file target changed: info=%v err=%v", info, err)
		}
		assertNoArchiveScratch(t, directory)
	})
}

func TestMissingArchiveTargetParentReportsUnattemptedRoutesInTextAndJSON(t *testing.T) {
	directory := t.TempDir()
	missingParent := filepath.Join(directory, "missing-parent")
	targetPath := filepath.Join(missingParent, "upstream.tar.gz")
	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		requestCount.Add(1)
		_, _ = response.Write([]byte("network route should not run"))
	}))
	defer server.Close()

	for _, outputFormat := range []string{"text", "json"} {
		t.Run(outputFormat, func(t *testing.T) {
			command := exec.Command("go", "run", "./cmd/do-work-cli", "--format", outputFormat,
				"fetch-archive", "--target", targetPath, "--url", server.URL)
			command.Dir = filepath.Join("..", "..")
			var standardOutput bytes.Buffer
			var standardError bytes.Buffer
			command.Stdout = &standardOutput
			command.Stderr = &standardError
			if err := command.Run(); err == nil {
				t.Fatal("fetch-archive unexpectedly accepted a missing target parent")
			}

			projectedEvidence := standardOutput.String()
			if outputFormat == "json" {
				var result struct {
					Findings []struct {
						ObservedEvidence []string `json:"observed_evidence"`
					} `json:"findings"`
				}
				if err := json.Unmarshal(standardOutput.Bytes(), &result); err != nil {
					t.Fatalf("decode JSON projection: %v\nstdout=%s\nstderr=%s", err, standardOutput.String(), standardError.String())
				}
				if len(result.Findings) != 1 {
					t.Fatalf("findings=%d, want 1", len(result.Findings))
				}
				projectedEvidence = strings.Join(result.Findings[0].ObservedEvidence, "\n")
			}
			for _, expected := range []string{
				"HTTP route: not attempted",
				"Git route: not attempted",
				"Set DO_WORK_UPSTREAM_URL",
			} {
				if !strings.Contains(projectedEvidence, expected) {
					t.Errorf("%s projection %q is missing %q", outputFormat, projectedEvidence, expected)
				}
			}
		})
	}
	if requestCount.Load() != 0 {
		t.Fatalf("parent preflight made %d network request(s)", requestCount.Load())
	}
	if _, err := os.Lstat(missingParent); !os.IsNotExist(err) {
		t.Fatalf("missing target parent was mutated: %v", err)
	}
}

func TestAbsentArchiveTargetRacePreservesTheCompetingCreation(t *testing.T) {
	directory := t.TempDir()
	sourceArchive := newFixtureArchive(t, directory, "payload.txt")
	archiveBytes := readFixtureBytes(t, sourceArchive)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write(archiveBytes) }))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	setBeforeArchiveFetchPublish(t, func() {
		if err := os.WriteFile(targetPath, []byte("competing creation\n"), 0o640); err != nil {
			t.Fatal(err)
		}
	})

	if _, err := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: server.URL}); err == nil {
		t.Fatal("competing creation was overwritten")
	}
	contents, readError := os.ReadFile(targetPath)
	if readError != nil || string(contents) != "competing creation\n" {
		t.Fatalf("competing target=%q err=%v", contents, readError)
	}
	assertNoArchiveScratch(t, directory)
}

func TestRegularArchiveTargetRacesPreserveTheCurrentObject(t *testing.T) {
	tests := []struct {
		name   string
		change func(*testing.T, string)
		want   string
	}{
		{
			name: "in-place mutation",
			change: func(t *testing.T, targetPath string) {
				if err := os.WriteFile(targetPath, []byte("mutated in place\n"), 0o640); err != nil {
					t.Fatal(err)
				}
			},
			want: "mutated in place\n",
		},
		{
			name: "replacement",
			change: func(t *testing.T, targetPath string) {
				if err := os.Remove(targetPath); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(targetPath, []byte("competing replacement\n"), 0o640); err != nil {
					t.Fatal(err)
				}
			},
			want: "competing replacement\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			sourceArchive := newFixtureArchive(t, directory, "payload.txt")
			archiveBytes := readFixtureBytes(t, sourceArchive)
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write(archiveBytes) }))
			defer server.Close()
			targetPath := filepath.Join(directory, "upstream.tar.gz")
			if err := os.WriteFile(targetPath, []byte("original bytes\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			setBeforeArchiveFetchPublish(t, func() { test.change(t, targetPath) })

			if _, err := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: server.URL}); err == nil {
				t.Fatal("concurrently changed target was overwritten")
			}
			contents, readError := os.ReadFile(targetPath)
			if readError != nil || string(contents) != test.want {
				t.Fatalf("current target=%q err=%v", contents, readError)
			}
			assertNoArchiveScratch(t, directory)
		})
	}
}

func TestRegularArchiveTargetRemovalBeforePublicationIsNotRecreated(t *testing.T) {
	directory := t.TempDir()
	sourceArchive := newFixtureArchive(t, directory, "payload.txt")
	archiveBytes := readFixtureBytes(t, sourceArchive)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write(archiveBytes) }))
	defer server.Close()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("original bytes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	setBeforeArchiveFetchPublish(t, func() {
		if err := os.Remove(targetPath); err != nil {
			t.Fatal(err)
		}
	})

	if _, err := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: server.URL}); err == nil {
		t.Fatal("removed target was recreated")
	}
	if _, err := os.Lstat(targetPath); !os.IsNotExist(err) {
		t.Fatalf("removed target was changed: %v", err)
	}
	assertNoArchiveScratch(t, directory)
}

func TestArchiveFetchParentSwapCannotRedirectRegularReplacement(t *testing.T) {
	directory := t.TempDir()
	parent := filepath.Join(directory, "parent")
	held := filepath.Join(directory, "held")
	outside := filepath.Join(directory, "outside")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	targetPath := filepath.Join(parent, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("old archive\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	protectedPath := filepath.Join(outside, "upstream.tar.gz")
	if err := os.WriteFile(protectedPath, []byte("protected\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	sourceArchive := newFixtureArchive(t, directory, "replacement.txt")
	archiveBytes := readFixtureBytes(t, sourceArchive)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) { _, _ = response.Write(archiveBytes) }))
	defer server.Close()
	setBeforeArchiveFetchPublish(t, func() {
		if err := os.Rename(parent, held); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, parent); err != nil {
			t.Fatal(err)
		}
	})

	result, err := FetchArchive(context.Background(), Request{ArchiveTargetPath: targetPath, UpstreamTarballURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	if result.RouteDescription != "upstream archive fetched over HTTP" {
		t.Fatalf("route=%q", result.RouteDescription)
	}
	protectedContents := readFixtureBytes(t, protectedPath)
	if string(protectedContents) != "protected\n" {
		t.Fatalf("outside target changed: %q", protectedContents)
	}
	publishedContents := readFixtureBytes(t, filepath.Join(held, "upstream.tar.gz"))
	if !bytes.Equal(publishedContents, archiveBytes) {
		t.Fatal("replacement was not confined to the opened parent")
	}
	assertNoArchiveScratch(t, held)
}

func TestGitFallbackReplacesAnUnchangedRegularArchiveAfterValidation(t *testing.T) {
	directory := t.TempDir()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("old archive bytes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "://invalid",
		UpstreamRepositoryURL: newFixtureRepository(t, directory),
	})
	if err != nil {
		t.Fatalf("FetchArchive: %v", err)
	}
	if !strings.HasPrefix(result.RouteDescription, "upstream archive fetched with git (HTTP route failed") {
		t.Fatalf("route = %q, want Git after failed HTTP", result.RouteDescription)
	}
	if !archiveIsReadable(context.Background(), targetPath) {
		t.Fatal("Git route did not publish a readable replacement")
	}
	if names := archiveEntryNames(t, targetPath); !strings.Contains(names, gitArchivePrefix+"tracked.txt") {
		t.Fatalf("archive entries = %q, want the Git route payload", names)
	}
	info, statError := os.Stat(targetPath)
	if statError != nil {
		t.Fatal(statError)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%o, want private candidate mode 600", info.Mode().Perm())
	}
}

func TestGitFallbackPreservesAConcurrentRegularTargetReplacement(t *testing.T) {
	directory := t.TempDir()
	targetPath := filepath.Join(directory, "upstream.tar.gz")
	if err := os.WriteFile(targetPath, []byte("old archive bytes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	setBeforeArchiveFetchPublish(t, func() {
		if err := os.Remove(targetPath); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(targetPath, []byte("competing Git-route replacement\n"), 0o640); err != nil {
			t.Fatal(err)
		}
	})

	_, err := FetchArchive(context.Background(), Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    "://invalid",
		UpstreamRepositoryURL: newFixtureRepository(t, directory),
	})
	if err == nil {
		t.Fatal("Git fallback overwrote a competing replacement")
	}
	contents, readError := os.ReadFile(targetPath)
	if readError != nil || string(contents) != "competing Git-route replacement\n" {
		t.Fatalf("current target=%q err=%v", contents, readError)
	}
	assertNoArchiveScratch(t, directory)
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

func setBeforeArchiveFetchPublish(t *testing.T, action func()) {
	t.Helper()
	previousHook := beforeArchiveFetchPublish
	callCount := 0
	beforeArchiveFetchPublish = func() {
		if callCount == 0 {
			action()
		}
		callCount++
	}
	t.Cleanup(func() { beforeArchiveFetchPublish = previousHook })
}

func assertNoArchiveScratch(t *testing.T, directory string) {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read directory: %v", err)
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".fetching.") {
			t.Errorf("archive scratch survived: %s", entry.Name())
		}
	}
}

func readFixtureBytes(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return contents
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
