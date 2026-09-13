package finalization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
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
			metadataOnly, err := implementationIsReleaseMetadataOnly(repositoryRoot, manifest, path)
			if err != nil {
				return fmt.Errorf("RELEASE-SHIPPED-CHANGE-UNVERIFIABLE: %w", err)
			}
			if metadataOnly {
				continue
			}
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
// pending changes must also belong to the manifest's commit_paths allowlist.
func implementationPathsForRelease(repositoryRoot string, manifest Manifest) ([]string, error) {
	if manifest.ProvenanceMode != ProvenanceSuppliedCommit {
		allowed := map[string]bool{}
		for _, path := range manifest.CommitPaths {
			allowed[filepath.ToSlash(filepath.Clean(path))] = true
		}
		paths := []string{}
		for _, arguments := range [][]string{
			{"diff", "--name-only", "--no-renames", "-z", "HEAD", "--"},
			{"ls-files", "--others", "--exclude-standard", "-z", "--"},
		} {
			output, err := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...).Output()
			if err != nil {
				return nil, fmt.Errorf("list pending implementation changes: %w", err)
			}
			for _, path := range strings.Split(string(output), "\x00") {
				if allowed[path] {
					paths = append(paths, path)
				}
			}
		}
		return paths, nil
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

// Ownership asks whether a file carries metadata; this guard instead asks
// whether the implementation changed anything besides its release version.
func implementationIsReleaseMetadataOnly(repositoryRoot string, manifest Manifest, path string) (bool, error) {
	switch filepath.Base(path) {
	case "package.json", "Cargo.toml", "pyproject.toml":
	default:
		return true, nil
	}
	beforeRevision := "HEAD"
	var after FileImage
	var err error
	if manifest.ProvenanceMode == ProvenanceSuppliedCommit {
		beforeRevision = manifest.ImplementationHash + "^1"
		after, err = committedManifestImage(repositoryRoot, manifest.ImplementationHash, path)
		if err == nil && exec.Command("git", "-C", repositoryRoot, "rev-parse", "--verify", "--quiet", beforeRevision).Run() != nil {
			// An initial commit adds its manifests without a preimage.
			return false, nil
		}
	} else {
		after, err = currentImage(repositoryRoot, path)
	}
	if err != nil {
		return false, err
	}
	before, err := committedManifestImage(repositoryRoot, beforeRevision, path)
	if err != nil || !before.Exists || !after.Exists {
		return false, err
	}
	return bytes.Equal(manifestWithoutReleaseVersion(path, before.Bytes), manifestWithoutReleaseVersion(path, after.Bytes)), nil
}

func committedManifestImage(repositoryRoot, revision, path string) (FileImage, error) {
	listing, err := exec.Command("git", "-C", repositoryRoot, "--literal-pathspecs", "ls-tree", "-z", revision, "--", path).Output()
	if err != nil {
		return FileImage{}, fmt.Errorf("inspect manifest %s at %s: %w", path, revision, err)
	}
	if len(listing) == 0 {
		return FileImage{Path: path}, nil
	}
	contents, err := exec.Command("git", "-C", repositoryRoot, "show", revision+":"+path).Output()
	if err != nil {
		return FileImage{}, fmt.Errorf("read manifest %s at %s: %w", path, revision, err)
	}
	return FileImage{Path: path, Exists: true, Bytes: contents}, nil
}

func manifestWithoutReleaseVersion(path string, contents []byte) []byte {
	if filepath.Base(path) == "package.json" {
		var document map[string]json.RawMessage
		if json.Unmarshal(contents, &document) != nil || document == nil {
			return contents
		}
		delete(document, "version")
		encoded, _ := json.Marshal(document)
		return encoded
	}
	// Keep TOML bytes outside the project's version assignment intact, including
	// dependency versions that happen to equal the package's own version.
	lines := strings.Split(string(contents), "\n")
	section := ""
	versionAssignment := regexp.MustCompile(`^(version[ \t]*=[ \t]*)["'][^"']*["']([ \t]*(?:#.*)?)$`)
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "["), "]"))
			continue
		}
		if filepath.Base(path) == "Cargo.toml" && section == "package" || filepath.Base(path) == "pyproject.toml" && (section == "project" || section == "tool.poetry") {
			if versionAssignment.MatchString(trimmed) {
				lines[index] = versionAssignment.ReplaceAllString(trimmed, `${1}""${2}`)
			}
		}
	}
	return []byte(strings.Join(lines, "\n"))
}
