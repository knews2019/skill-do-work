// Package archivefetch obtains the upstream suite tarball by whichever route works, and
// publishes it only whole.
//
// Route 1 is the anonymous tarball over HTTP, delegated to the shipped atomic-download
// primitive so it inherits that helper's retry and optional credentials. Route 2 is a
// shallow clone repacked with `git archive`.
//
// The git transport sits behind a different rate limiter than codeload, which is what makes
// route 2 load-bearing rather than decorative: a sustained 429 defeats retry alone.
// `git archive` — never a worktree copy — is mandatory, because only it honours
// `export-ignore`; cp -R, rsync, and tarring a clone all ship the maintainer-only tree into
// consumer installs, and nothing downstream catches that.
package archivefetch

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DefaultUpstreamURL is the branch tarball the suite ships against. DO_WORK_UPSTREAM_URL
// overrides it to route around a blocked host.
const DefaultUpstreamURL = "https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz"

// gitArchivePrefix stays independent of the selected branch. Both callers strip exactly
// one leading component, so the Git route must always emit exactly one directory here.
const gitArchivePrefix = "upstream/"

const branchArchiveMarker = "/archive/refs/heads/"

// Request names everything a fetch needs. AtomicDownloadScript is resolved by the caller,
// because only it knows which of the two shipped layouts it is running from.
type Request struct {
	ArchiveTargetPath     string
	UpstreamTarballURL    string
	UpstreamRepositoryURL string
	AtomicDownloadScript  string
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

// LocateAtomicDownloadScript probes the two mirror-relative depths the shell fetcher probed,
// so the primitive is found from either shipped layout.
func LocateAtomicDownloadScript(scriptDirectory string) string {
	for _, candidate := range []string{
		filepath.Join(scriptDirectory, "..", "scripts", "atomic-download.sh"),
		filepath.Join(scriptDirectory, "..", "skills", "do-work", "scripts", "atomic-download.sh"),
	} {
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() {
			return candidate
		}
	}
	return ""
}

// FetchArchive tries the HTTP route, then the git route, and leaves any pre-existing target
// untouched when both fail. On failure the error names both route outcomes and the escape
// hatch, matching the two stderr lines the shell fetcher printed.
func FetchArchive(ctx context.Context, request Request) (Result, error) {
	upstreamBranch, repositoryURL := deriveGitRoute(request.UpstreamTarballURL, request.UpstreamRepositoryURL)

	httpRouteOutcome := "not attempted"
	if request.AtomicDownloadScript == "" {
		httpRouteOutcome = "unavailable (atomic-download.sh not found beside this script)"
	} else if downloadThroughAtomicPrimitive(ctx, request) {
		return Result{RouteDescription: "upstream archive fetched over HTTP"}, nil
	} else {
		httpRouteOutcome = "failed (host unreachable, rate limited, or archive unreadable)"
	}

	gitRouteOutcome := fetchThroughGitRoute(ctx, request, repositoryURL, upstreamBranch)
	if gitRouteOutcome == "" {
		return Result{RouteDescription: fmt.Sprintf("upstream archive fetched with git (HTTP route %s)", httpRouteOutcome)}, nil
	}
	return Result{}, fmt.Errorf(
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

func downloadThroughAtomicPrimitive(ctx context.Context, request Request) bool {
	download := exec.CommandContext(ctx, "bash", request.AtomicDownloadScript,
		request.UpstreamTarballURL, request.ArchiveTargetPath)
	download.Stderr = nil
	if err := download.Run(); err != nil {
		return false
	}
	return archiveIsReadable(ctx, request.ArchiveTargetPath)
}

// fetchThroughGitRoute returns "" on success, or the outcome phrase describing why it did
// not run or did not complete. It stages into a temporary beside the target and only renames
// after the staged archive reads back, so a pre-existing target survives a total failure.
func fetchThroughGitRoute(ctx context.Context, request Request, repositoryURL, upstreamBranch string) string {
	if repositoryURL == "" {
		return "unavailable (no repository URL supplied and none derivable from the tarball URL)"
	}
	if _, err := exec.LookPath("git"); err != nil {
		return "unavailable (git is not installed)"
	}
	cloneDirectory, err := os.MkdirTemp("", "do-work-upstream-clone.*")
	if err != nil {
		return "failed (could not allocate private working paths)"
	}
	defer func() { _ = os.RemoveAll(cloneDirectory) }()
	stageFile, err := os.CreateTemp(filepath.Dir(request.ArchiveTargetPath),
		filepath.Base(request.ArchiveTargetPath)+".fetching.*")
	if err != nil {
		return "failed (could not allocate private working paths)"
	}
	stagePath := stageFile.Name()
	_ = stageFile.Close()
	stagePublished := false
	defer func() {
		if !stagePublished {
			_ = os.Remove(stagePath)
		}
	}()

	// git clone insists on creating the directory itself.
	if err := os.RemoveAll(cloneDirectory); err != nil {
		return "failed (could not allocate private working paths)"
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
		return "failed (clone, repack, or publication did not complete)"
	}

	stageHandle, err := os.OpenFile(stagePath, os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return "failed (clone, repack, or publication did not complete)"
	}
	archiveCommand := exec.CommandContext(ctx, "git", "-C", cloneDirectory,
		"archive", "--format=tar.gz", "--prefix="+gitArchivePrefix, "HEAD")
	archiveCommand.Stdout = stageHandle
	archiveErr := archiveCommand.Run()
	closeErr := stageHandle.Close()
	if archiveErr != nil || closeErr != nil {
		return "failed (clone, repack, or publication did not complete)"
	}
	if info, err := os.Stat(stagePath); err != nil || info.Size() == 0 {
		return "failed (clone, repack, or publication did not complete)"
	}
	if !archiveIsReadable(ctx, stagePath) {
		return "failed (clone, repack, or publication did not complete)"
	}
	if err := os.Rename(stagePath, request.ArchiveTargetPath); err != nil {
		return "failed (clone, repack, or publication did not complete)"
	}
	stagePublished = true
	return ""
}

// archiveIsReadable confirms a candidate really is a readable gzipped tar before it is
// published, so a rate-limit HTML body never lands where a suite archive should be.
func archiveIsReadable(ctx context.Context, archivePath string) bool {
	return exec.CommandContext(ctx, "tar", "tzf", archivePath).Run() == nil
}
