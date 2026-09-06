package finalization

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/releaseownership"
)

// A release is a change to shipped files, and the converse held four times in one run:
// 0.305.17, .21, .22 and .23 bumped the version for changes under _dev/tests alone.
// Nothing refused them because the finalizer checked the release payload and never the
// implementation. This is that check.
//
// What ships is keyed on the declaration, not on a list: the module sources
// suite/modules.tsv names (a module row whose source carries its own VERSION), plus the
// declaration's own directory and the installer tools that consume it, which reach a
// consumer by the same route. A repository that declares no modules is a consumer
// project, whose releases are its own package's business, and is not guarded.
func releaseShippedChangeError(repositoryRoot string, manifest Manifest) error {
	tracked, err := enumerateTrackedReleasePaths(repositoryRoot)
	if err != nil {
		return fmt.Errorf("RELEASE-SHIPPED-CHANGE-UNVERIFIABLE: enumerate tracked paths: %w", err)
	}
	trackedSet := map[string]bool{}
	for _, path := range tracked {
		trackedSet[filepath.ToSlash(filepath.Clean(path))] = true
	}
	roots, err := releaseownership.DeclaredMaintainerReleaseRoots(trackedSet, headReleaseImage(repositoryRoot))
	if err != nil {
		return fmt.Errorf("RELEASE-SHIPPED-CHANGE-UNVERIFIABLE: %w", err)
	}
	if len(roots) == 0 {
		return nil
	}
	roots = append(roots, "suite", "tools")

	implementationPaths, err := implementationPathsForRelease(repositoryRoot, manifest)
	if err != nil {
		return fmt.Errorf("RELEASE-SHIPPED-CHANGE-UNVERIFIABLE: %w", err)
	}
	for _, path := range implementationPaths {
		path = filepath.ToSlash(filepath.Clean(path))
		if releaseownership.IsReleaseMetadataPath(path) {
			continue
		}
		for _, root := range roots {
			if path == root || strings.HasPrefix(path, root+"/") {
				return nil
			}
		}
	}
	return fmt.Errorf("RELEASE-WITHOUT-SHIPPED-CHANGE: the implementation changes no path under a shipped root (%s) other than release metadata; a maintainer-only change is not a release, so finalize without release_manifest_path", strings.Join(roots, ", "))
}

// implementationPathsForRelease lists what the implementation changed. Under
// supplied_commit provenance that is the named commit's first-parent diff: for a
// merge, what the merge brought in; for a plain commit, its own change. Under
// primary_commit provenance the finalization commit is the implementation, and its
// allowlist is the manifest's commit_paths.
func implementationPathsForRelease(repositoryRoot string, manifest Manifest) ([]string, error) {
	if manifest.ProvenanceMode != ProvenanceSuppliedCommit {
		return manifest.CommitPaths, nil
	}
	hash := manifest.ImplementationHash
	arguments := []string{"-C", repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", "--root"}
	if err := exec.Command("git", "-C", repositoryRoot, "rev-parse", "--verify", "--quiet", hash+"^1").Run(); err == nil {
		arguments = append(arguments, hash+"^1", hash)
	} else {
		arguments = append(arguments, hash)
	}
	output, err := exec.Command("git", arguments...).Output()
	if err != nil {
		return nil, fmt.Errorf("list implementation paths of %s: %w", hash, err)
	}
	paths := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths, nil
}
