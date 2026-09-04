// Package archivefetch obtains the upstream suite tarball by whichever route works, and
// publishes it only whole.
//
// Route 1 is the anonymous tarball over HTTP. Route 2 is a shallow clone repacked with
// `git archive`. Both routes prepare and validate a private adjacent candidate before one
// shared publication step.
//
// The git transport sits behind a different rate limiter than codeload, which is what makes
// route 2 load-bearing rather than decorative: a sustained 429 defeats retry alone.
// `git archive` — never a worktree copy — is mandatory, because only it honours
// `export-ignore`; cp -R, rsync, and tarring a clone all ship the maintainer-only tree into
// consumer installs, and nothing downstream catches that.
package archivefetch

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DefaultUpstreamURL is the branch tarball the suite ships against. DO_WORK_UPSTREAM_URL
// overrides it to route around a blocked host.
const DefaultUpstreamURL = "https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz"

// gitArchivePrefix stays independent of the selected branch. Both callers strip exactly
// one leading component, so the Git route must always emit exactly one directory here.
const gitArchivePrefix = "upstream/"

const branchArchiveMarker = "/archive/refs/heads/"

// Request names everything a fetch needs. HTTP transfer is always in-process.
type Request struct {
	ArchiveTargetPath     string
	UpstreamTarballURL    string
	UpstreamRepositoryURL string
}

// Result reports which route produced the archive, in the wording the shell fetcher printed
// so the existing route assertions keep reading the same text.
type Result struct {
	RouteDescription string
}

// UpstreamURLFromEnvironment returns the configured upstream tarball URL, honouring the
// documented DO_WORK_UPSTREAM_URL escape hatch.
func UpstreamURLFromEnvironment() string {
	if configured := os.Getenv("DO_WORK_UPSTREAM_URL"); configured != "" {
		return configured
	}
	return DefaultUpstreamURL
}

// FetchArchive tries the HTTP route, then the git route, and leaves any pre-existing target
// untouched when both fail. The error names both route outcomes, including transport routes
// refused by target preflight, and the escape hatch the shell fetcher printed.
func FetchArchive(ctx context.Context, request Request) (Result, error) {
	if request.ArchiveTargetPath == "" || request.UpstreamTarballURL == "" {
		return Result{}, fmt.Errorf("archive target path and upstream tarball URL are required")
	}
	upstreamBranch, repositoryURL := deriveGitRoute(request.UpstreamTarballURL, request.UpstreamRepositoryURL)
	parentRoot, err := os.OpenRoot(filepath.Dir(request.ArchiveTargetPath))
	if err != nil {
		parentOutcome := fmt.Sprintf("not attempted (archive target parent is unavailable: %v)", err)
		return Result{}, newArchiveFetchFailure(parentOutcome, parentOutcome)
	}
	defer parentRoot.Close()
	targetSnapshot, err := inspectArchiveTarget(parentRoot, filepath.Base(request.ArchiveTargetPath))
	if err != nil {
		targetOutcome := fmt.Sprintf("not attempted (archive target is unsafe: %v)", err)
		return Result{}, newArchiveFetchFailure(targetOutcome, targetOutcome)
	}

	httpRouteOutcome := "not attempted"
	httpStageName, transfer := prepareDownloadCandidate(ctx, request.UpstreamTarballURL, parentRoot, targetSnapshot.targetName, archiveStreamReadable)
	if transfer.Err == nil {
		publishError := publishArchiveCandidate(parentRoot, targetSnapshot, httpStageName)
		_ = parentRoot.Remove(httpStageName)
		if publishError == nil {
			return Result{RouteDescription: "upstream archive fetched over HTTP"}, nil
		}
		httpRouteOutcome = "failed (validated archive could not be published safely)"
	} else {
		httpRouteOutcome = "failed (host unreachable, rate limited, or archive unreadable)"
	}

	gitStageName, gitRouteOutcome := prepareGitArchiveCandidate(ctx, parentRoot, targetSnapshot.targetName, repositoryURL, upstreamBranch)
	if gitRouteOutcome == "" {
		publishError := publishArchiveCandidate(parentRoot, targetSnapshot, gitStageName)
		_ = parentRoot.Remove(gitStageName)
		if publishError == nil {
			return Result{RouteDescription: fmt.Sprintf("upstream archive fetched with git (HTTP route %s)", httpRouteOutcome)}, nil
		}
		gitRouteOutcome = "failed (clone, repack, or publication did not complete)"
	}
	return Result{}, newArchiveFetchFailure(httpRouteOutcome, gitRouteOutcome)
}

func newArchiveFetchFailure(httpRouteOutcome, gitRouteOutcome string) error {
	return fmt.Errorf(
		"upstream archive could not be fetched. HTTP route: %s. Git route: %s. "+
			"Set DO_WORK_UPSTREAM_URL to a reachable archive URL to route around a blocked host.",
		httpRouteOutcome, gitRouteOutcome)
}

// deriveGitRoute reads the requested branch and repository out of a GitHub branch-tarball
// URL. Any other URL shape derives no repository URL at all rather than guessing one.
func deriveGitRoute(tarballURL, suppliedRepositoryURL string) (upstreamBranch, repositoryURL string) {
	repositoryURL = suppliedRepositoryURL
	markerIndex := strings.LastIndex(tarballURL, branchArchiveMarker)
	if markerIndex < 0 || !strings.HasSuffix(tarballURL, ".tar.gz") {
		return "", repositoryURL
	}
	upstreamBranch = strings.TrimSuffix(tarballURL[markerIndex+len(branchArchiveMarker):], ".tar.gz")
	repositoryBase := tarballURL[:markerIndex]
	if repositoryURL == "" {
		repositoryURL = repositoryBase + ".git"
	}
	return upstreamBranch, repositoryURL
}

type DownloadResult struct {
	StatusCode   int
	BytesWritten int64
	Attempts     int
	Err          error
}

var beforeAtomicDownloadPublish = func() {}
var beforeArchiveFetchPublish = func() {}

const (
	defaultAtomicRetryDelay  = 2 * time.Second
	defaultAtomicRetryBudget = 60 * time.Second
)

var (
	atomicHTTPClient  = &http.Client{}
	atomicRetryDelay  = defaultAtomicRetryDelay
	atomicRetryBudget = defaultAtomicRetryBudget
)

// DownloadAtomic performs curl-compatible initial-plus-three eligible retries and publishes
// a private adjacent file with no overwrite. Publication is the final commit point: after
// the rooted link succeeds, no failure path removes the destination pathname.
func DownloadAtomic(ctx context.Context, sourceURL, targetPath string) DownloadResult {
	if sourceURL == "" || targetPath == "" {
		return DownloadResult{Err: fmt.Errorf("source URL and target path are required")}
	}
	parentRoot, err := os.OpenRoot(filepath.Dir(targetPath))
	if err != nil {
		return DownloadResult{Err: err}
	}
	defer parentRoot.Close()
	targetName := filepath.Base(targetPath)
	if _, err := parentRoot.Lstat(targetName); err == nil {
		return DownloadResult{Err: fmt.Errorf("target already exists")}
	} else if !os.IsNotExist(err) {
		return DownloadResult{Err: err}
	}
	stageName, result := prepareDownloadCandidate(ctx, sourceURL, parentRoot, targetName, nil)
	if result.Err != nil {
		return result
	}
	defer func() { _ = parentRoot.Remove(stageName) }()
	beforeAtomicDownloadPublish()
	if err = parentRoot.Link(stageName, targetName); err != nil {
		result.Err = err
		return result
	}
	return result
}

func prepareDownloadCandidate(ctx context.Context, sourceURL string, parentRoot *os.Root, targetName string, validator func(io.Reader) error) (string, DownloadResult) {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	randomSuffix := make([]byte, 8)
	if _, err := rand.Read(randomSuffix); err != nil {
		return "", DownloadResult{Err: err}
	}
	stageName := fmt.Sprintf(".%s.fetching.%x", targetName, randomSuffix)
	keepCandidate := false
	defer func() {
		if !keepCandidate {
			_ = parentRoot.Remove(stageName)
		}
	}()
	statusCode := 0
	attempts := 0
	var lastError error
	var written int64
	startedAt := time.Now()
	for attempt := 0; attempt < 4; attempt++ {
		attempts = attempt + 1
		_ = parentRoot.Remove(stageName)
		stage, stageError := parentRoot.OpenFile(stageName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if stageError != nil {
			return "", DownloadResult{StatusCode: statusCode, Attempts: attempts, Err: stageError}
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
		if err != nil {
			_ = stage.Close()
			return "", DownloadResult{Attempts: attempts, Err: redactedTransferError(fmt.Errorf("build HTTP request: %w", err), token)}
		}
		if token != "" {
			request.Header.Set("Authorization", "Bearer "+token)
		}
		response, err := atomicHTTPClient.Do(request)
		if err == nil {
			statusCode = response.StatusCode
			if statusCode >= 200 && statusCode < 300 {
				written, err = io.Copy(stage, response.Body)
				if closeError := response.Body.Close(); err == nil {
					err = closeError
				}
				if err == nil {
					err = stage.Sync()
				}
			} else {
				_ = response.Body.Close()
				err = fmt.Errorf("HTTP status %d", statusCode)
			}
		}
		if closeError := stage.Close(); err == nil {
			err = closeError
		}
		if err == nil && statusCode >= 200 && statusCode < 300 && written > 0 {
			lastError = nil
			break
		}
		if err == nil && written == 0 {
			err = fmt.Errorf("downloaded body is empty")
		}
		lastError = err
		if attempt == 3 || !retryEligible(statusCode, err) {
			break
		}
		if time.Since(startedAt)+atomicRetryDelay > atomicRetryBudget {
			lastError = fmt.Errorf("retry budget exhausted after %d attempt(s): %w", attempts, lastError)
			break
		}
		select {
		case <-ctx.Done():
			return "", DownloadResult{StatusCode: statusCode, Attempts: attempts, Err: ctx.Err()}
		case <-time.After(atomicRetryDelay):
		}
	}
	if lastError != nil {
		return "", DownloadResult{StatusCode: statusCode, Attempts: attempts, Err: redactedTransferError(lastError, token)}
	}
	if validator != nil {
		stagedReader, openError := parentRoot.Open(stageName)
		if openError != nil {
			return "", DownloadResult{StatusCode: statusCode, Attempts: attempts, Err: openError}
		}
		validationError := validator(stagedReader)
		closeError := stagedReader.Close()
		if validationError == nil {
			validationError = closeError
		}
		if validationError != nil {
			return "", DownloadResult{StatusCode: statusCode, Attempts: attempts, Err: validationError}
		}
	}
	keepCandidate = true
	return stageName, DownloadResult{StatusCode: statusCode, BytesWritten: written, Attempts: attempts}
}

type archiveTargetSnapshot struct {
	targetName      string
	initiallyAbsent bool
	fileIdentity    os.FileInfo
	contentDigest   [sha256.Size]byte
}

func inspectArchiveTarget(parentRoot *os.Root, targetName string) (archiveTargetSnapshot, error) {
	if targetName == "" || targetName == "." {
		return archiveTargetSnapshot{}, fmt.Errorf("target name is required")
	}
	info, err := parentRoot.Lstat(targetName)
	if os.IsNotExist(err) {
		return archiveTargetSnapshot{targetName: targetName, initiallyAbsent: true}, nil
	}
	if err != nil {
		return archiveTargetSnapshot{}, err
	}
	if !info.Mode().IsRegular() {
		return archiveTargetSnapshot{}, fmt.Errorf("target must be absent or a regular non-symlink file")
	}
	fileIdentity, contentDigest, err := readRegularTarget(parentRoot, targetName)
	if err != nil {
		return archiveTargetSnapshot{}, err
	}
	return archiveTargetSnapshot{
		targetName:    targetName,
		fileIdentity:  fileIdentity,
		contentDigest: contentDigest,
	}, nil
}

func readRegularTarget(parentRoot *os.Root, targetName string) (os.FileInfo, [sha256.Size]byte, error) {
	var emptyDigest [sha256.Size]byte
	beforeOpen, err := parentRoot.Lstat(targetName)
	if err != nil {
		return nil, emptyDigest, err
	}
	if !beforeOpen.Mode().IsRegular() {
		return nil, emptyDigest, fmt.Errorf("target must remain a regular non-symlink file")
	}
	targetHandle, err := parentRoot.Open(targetName)
	if err != nil {
		return nil, emptyDigest, err
	}
	defer targetHandle.Close()
	openedInfo, err := targetHandle.Stat()
	if err != nil {
		return nil, emptyDigest, err
	}
	afterOpen, err := parentRoot.Lstat(targetName)
	if err != nil {
		return nil, emptyDigest, err
	}
	if !afterOpen.Mode().IsRegular() || !os.SameFile(beforeOpen, openedInfo) || !os.SameFile(openedInfo, afterOpen) {
		return nil, emptyDigest, fmt.Errorf("target changed while it was inspected")
	}
	digest, err := digestReaderContents(targetHandle)
	if err != nil {
		return nil, emptyDigest, err
	}
	afterRead, err := parentRoot.Lstat(targetName)
	if err != nil {
		return nil, emptyDigest, err
	}
	if !afterRead.Mode().IsRegular() || !os.SameFile(openedInfo, afterRead) {
		return nil, emptyDigest, fmt.Errorf("target changed while it was inspected")
	}
	return openedInfo, digest, nil
}

func digestReaderContents(contents io.Reader) ([sha256.Size]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, contents); err != nil {
		return [sha256.Size]byte{}, err
	}
	var digest [sha256.Size]byte
	copy(digest[:], hash.Sum(nil))
	return digest, nil
}

func publishArchiveCandidate(parentRoot *os.Root, targetSnapshot archiveTargetSnapshot, stageName string) error {
	beforeArchiveFetchPublish()
	if targetSnapshot.initiallyAbsent {
		if _, err := parentRoot.Lstat(targetSnapshot.targetName); err == nil {
			return fmt.Errorf("archive target was created concurrently")
		} else if !os.IsNotExist(err) {
			return err
		}
		return parentRoot.Link(stageName, targetSnapshot.targetName)
	}
	currentIdentity, currentDigest, err := readRegularTarget(parentRoot, targetSnapshot.targetName)
	if err != nil {
		return err
	}
	if !os.SameFile(targetSnapshot.fileIdentity, currentIdentity) || currentDigest != targetSnapshot.contentDigest {
		return fmt.Errorf("archive target changed before publication")
	}
	return parentRoot.Rename(stageName, targetSnapshot.targetName)
}

func archiveBytesReadable(contents []byte) error {
	return archiveStreamReadable(bytes.NewReader(contents))
}

func archiveStreamReadable(contents io.Reader) error {
	reader, err := gzip.NewReader(contents)
	if err != nil {
		return fmt.Errorf("archive is not gzip: %w", err)
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)
	entryCount := 0
	for {
		_, nextError := tarReader.Next()
		if nextError == io.EOF {
			break
		}
		if nextError != nil {
			return fmt.Errorf("archive tar stream is incomplete: %w", nextError)
		}
		entryCount++
		if _, copyError := io.Copy(io.Discard, tarReader); copyError != nil {
			return fmt.Errorf("archive entry is incomplete: %w", copyError)
		}
	}
	if entryCount == 0 {
		return fmt.Errorf("archive has no readable tar entry")
	}
	return nil
}

func retryEligible(status int, err error) bool {
	if err != nil && status == 0 {
		return true
	}
	switch status {
	case 408, 429, 500, 502, 503, 504:
		return true
	}
	return false
}
func redactedTransferError(err error, token string) error {
	message := err.Error()
	if token != "" {
		message = strings.ReplaceAll(message, token, "[REDACTED]")
	}
	return fmt.Errorf("%s", message)
}

// prepareGitArchiveCandidate returns a validated private candidate and an empty outcome on
// success, or the outcome phrase describing why the route did not run or did not complete.
func prepareGitArchiveCandidate(ctx context.Context, parentRoot *os.Root, targetName, repositoryURL, upstreamBranch string) (string, string) {
	if repositoryURL == "" {
		return "", "unavailable (no repository URL supplied and none derivable from the tarball URL)"
	}
	if _, err := exec.LookPath("git"); err != nil {
		return "", "unavailable (git is not installed)"
	}
	cloneDirectory, err := os.MkdirTemp("", "do-work-upstream-clone.*")
	if err != nil {
		return "", "failed (could not allocate private working paths)"
	}
	defer func() { _ = os.RemoveAll(cloneDirectory) }()
	randomSuffix := make([]byte, 8)
	if _, err := rand.Read(randomSuffix); err != nil {
		return "", "failed (could not allocate private working paths)"
	}
	stageName := fmt.Sprintf(".%s.fetching.%x", targetName, randomSuffix)
	stageHandle, err := parentRoot.OpenFile(stageName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", "failed (could not allocate private working paths)"
	}
	keepCandidate := false
	defer func() {
		_ = stageHandle.Close()
		if !keepCandidate {
			_ = parentRoot.Remove(stageName)
		}
	}()

	// git clone insists on creating the directory itself.
	if err := os.RemoveAll(cloneDirectory); err != nil {
		return "", "failed (could not allocate private working paths)"
	}
	// --single-branch --branch selects a named branch exactly, so a missing branch fails
	// rather than silently substituting the repository's HEAD.
	cloneArgs := []string{"clone", "--depth", "1", "--quiet"}
	if upstreamBranch != "" {
		cloneArgs = append(cloneArgs, "--single-branch", "--branch", upstreamBranch)
	}
	cloneArgs = append(cloneArgs, repositoryURL, cloneDirectory)
	clone := exec.CommandContext(ctx, "git", cloneArgs...)
	clone.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := clone.Run(); err != nil {
		return "", "failed (clone, repack, or publication did not complete)"
	}

	archiveCommand := exec.CommandContext(ctx, "git", "-C", cloneDirectory,
		"archive", "--format=tar.gz", "--prefix="+gitArchivePrefix, "HEAD")
	archiveCommand.Stdout = stageHandle
	archiveErr := archiveCommand.Run()
	closeErr := stageHandle.Close()
	if archiveErr != nil || closeErr != nil {
		return "", "failed (clone, repack, or publication did not complete)"
	}
	if info, err := parentRoot.Stat(stageName); err != nil || info.Size() == 0 {
		return "", "failed (clone, repack, or publication did not complete)"
	}
	stagedReader, err := parentRoot.Open(stageName)
	if err != nil {
		return "", "failed (clone, repack, or publication did not complete)"
	}
	validationError := archiveStreamReadable(stagedReader)
	closeValidationError := stagedReader.Close()
	if validationError == nil {
		validationError = closeValidationError
	}
	if validationError != nil {
		return "", "failed (clone, repack, or publication did not complete)"
	}
	keepCandidate = true
	return stageName, ""
}

// archiveIsReadable confirms a candidate really is a readable gzipped tar before it is
// published, so a rate-limit HTML body never lands where a suite archive should be.
func archiveIsReadable(ctx context.Context, archivePath string) bool {
	return exec.CommandContext(ctx, "tar", "tzf", archivePath).Run() == nil
}
