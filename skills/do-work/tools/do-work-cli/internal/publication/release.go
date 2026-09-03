package publication

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func BuildReleasePlan(repositoryRoot string, manifest Manifest) PublicationPlan {
	plan := PublicationPlan{Operation: OperationRelease, RepositoryRoot: repositoryRoot, CommitMessage: manifest.CommitMessage}
	release := manifest.Release
	if release == nil || len(release.Changelogs) == 0 {
		return refusedPlan(plan, "RELEASE-MANIFEST-MISSING", "release requires version and changelog targets", nil)
	}
	oldVersion, newVersion := release.OldVersion, release.NewVersion
	if len(release.Targets) > 0 {
		if oldVersion == "" {
			oldVersion = release.Targets[0].OldVersion
		}
		if newVersion == "" {
			newVersion = release.Targets[0].NewVersion
		}
	}
	if compareSemver(oldVersion, newVersion) != 1 {
		return refusedPlan(plan, "RELEASE-VERSION-NOT-INCREASING", "new version must be valid semver and strictly greater", nil)
	}
	for _, target := range release.Targets {
		path, pathError := containedPath(target.Path)
		if pathError != nil {
			return refusedPlan(plan, "RELEASE-PATH-UNSAFE", pathError.Error(), nil, target.Path)
		}
		if !release.MaintainerRelease {
			if gap := releaseTargetOwnershipGap(repositoryRoot, path, false); gap != "" {
				return refusedPlan(plan, "RELEASE-TARGET-OWNERSHIP-UNVERIFIED", ownershipRefusalReason(path, gap), nil, path)
			}
		}
		if target.OldVersion != oldVersion || target.NewVersion != newVersion {
			return refusedPlan(plan, "RELEASE-MIRROR-DISAGREEMENT", "all version mirrors must declare the same old and new versions", nil, path)
		}
		expectedBytes, _, expectedError := readPayload(repositoryRoot, target.ExpectedPayload)
		newBytes, _, newError := readPayload(repositoryRoot, target.NewPayload)
		if expectedError != nil || newError != nil {
			return refusedPlan(plan, "RELEASE-PAYLOAD-INVALID", firstError(expectedError, newError).Error(), nil, path)
		}
		currentBytes, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if readError != nil || !bytes.Equal(currentBytes, expectedBytes) {
			return refusedPlan(plan, "RELEASE-PREIMAGE-STALE", "version target does not match expected bytes", nil, path)
		}
		if !bytes.Contains(expectedBytes, []byte(oldVersion)) || !bytes.Contains(newBytes, []byte(newVersion)) {
			return refusedPlan(plan, "RELEASE-VERSION-EVIDENCE-MISSING", "version payloads do not carry the declared old/new values", nil, path)
		}
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: MutationReplace, Path: path, ExpectedBytes: expectedBytes, Contents: newBytes})
	}
	targetedPaths := map[string]bool{}
	for _, mutation := range plan.Mutations {
		targetedPaths[mutation.Path] = true
	}
	for _, requiredMirror := range release.RequiredMirrors {
		mirrorPath, mirrorError := containedPath(requiredMirror)
		if mirrorError != nil || !targetedPaths[mirrorPath] {
			return refusedPlan(plan, "RELEASE-MIRROR-UNDECLARED", "every discovered version mirror must be declared as a target", nil, requiredMirror)
		}
	}
	var changelogMirrorBytes []byte
	var changelogMirrorExpected []byte
	var changelogMirrorCreate *bool
	for _, changelog := range release.Changelogs {
		path, pathError := containedPath(changelog.Path)
		if pathError != nil {
			return refusedPlan(plan, "RELEASE-PATH-UNSAFE", pathError.Error(), nil, changelog.Path)
		}
		if !release.MaintainerRelease {
			if gap := releaseTargetOwnershipGap(repositoryRoot, path, changelog.Create); gap != "" {
				return refusedPlan(plan, "RELEASE-TARGET-OWNERSHIP-UNVERIFIED", ownershipRefusalReason(path, gap), nil, path)
			}
		}
		newBytes, _, newError := readPayload(repositoryRoot, changelog.NewPayload)
		if newError != nil {
			return refusedPlan(plan, "RELEASE-PAYLOAD-INVALID", newError.Error(), nil, path)
		}
		if changelogMirrorBytes == nil {
			changelogMirrorBytes = append([]byte(nil), newBytes...)
		} else if !bytes.Equal(changelogMirrorBytes, newBytes) {
			return refusedPlan(plan, "RELEASE-CHANGELOG-MIRROR-DIVERGED", "every declared changelog mirror must publish byte-identical caller-authored bytes", nil, path)
		}
		var expectedBytes []byte
		if changelog.Create {
			if pathExists(repositoryRoot, path) {
				return refusedPlan(plan, "RELEASE-CHANGELOG-COLLISION", "new changelog target already exists", nil, path)
			}
			if changelog.InsertionAnchor == "" || bytes.Count(newBytes, []byte(changelog.InsertionAnchor)) != 1 {
				return refusedPlan(plan, "RELEASE-CHANGELOG-ANCHOR", "bootstrap changelog anchor must occur exactly once in new bytes", nil, path)
			}
		} else {
			var expectedError error
			expectedBytes, _, expectedError = readPayload(repositoryRoot, changelog.ExpectedPayload)
			if expectedError != nil {
				return refusedPlan(plan, "RELEASE-PAYLOAD-INVALID", expectedError.Error(), nil, path)
			}
			currentBytes, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
			if readError != nil || !bytes.Equal(currentBytes, expectedBytes) {
				return refusedPlan(plan, "RELEASE-PREIMAGE-STALE", "changelog target does not match expected bytes", nil, path)
			}
			if changelog.InsertionAnchor == "" || bytes.Count(expectedBytes, []byte(changelog.InsertionAnchor)) != 1 {
				return refusedPlan(plan, "RELEASE-CHANGELOG-ANCHOR", "changelog insertion anchor must occur exactly once in the expected preimage", nil, path)
			}
		}
		if changelogMirrorCreate == nil {
			create := changelog.Create
			changelogMirrorCreate = &create
			changelogMirrorExpected = append([]byte(nil), expectedBytes...)
		} else if *changelogMirrorCreate != changelog.Create || !bytes.Equal(changelogMirrorExpected, expectedBytes) {
			return refusedPlan(plan, "RELEASE-CHANGELOG-MIRROR-DIVERGED", "every declared changelog mirror must have byte-identical expected and new bytes", nil, path)
		}
		if strings.TrimSpace(changelog.EntryKey) == "" || strings.TrimSpace(changelog.EntryTitle) == "" || bytes.Contains(expectedBytes, []byte(changelog.EntryKey)) || bytes.Contains(expectedBytes, []byte(changelog.EntryTitle)) || bytes.Count(newBytes, []byte(changelog.EntryKey)) != 1 || bytes.Count(newBytes, []byte(changelog.EntryTitle)) != 1 {
			return refusedPlan(plan, "RELEASE-CHANGELOG-ENTRY-DUPLICATE", "new changelog key and title must be absent from old bytes and unique in new bytes", nil, path)
		}
		kind := MutationReplace
		if changelog.Create {
			kind = MutationCreate
		}
		plan.Mutations = append(plan.Mutations, PlannedMutation{Kind: kind, Path: path, ExpectedBytes: expectedBytes, Contents: newBytes, Mode: 0o644})
	}
	plan = finalizePlan(plan)
	if plan.Refusal != nil {
		return plan
	}
	directories, topologyError := planCreatedDirectories(repositoryRoot, plan.TargetPaths)
	if topologyError != nil {
		return refusedPlan(plan, "RELEASE-TOPOLOGY-UNSAFE", topologyError.Error(), nil)
	}
	plan.CreatedDirectoryPaths = directories
	return plan
}

func compareSemver(oldVersion, newVersion string) int {
	oldParts, oldOK := parseSemver(oldVersion)
	newParts, newOK := parseSemver(newVersion)
	if !oldOK || !newOK {
		return 0
	}
	for index := range oldParts {
		if newParts[index] > oldParts[index] {
			return 1
		}
		if newParts[index] < oldParts[index] {
			return -1
		}
	}
	return 0
}

func parseSemver(value string) ([3]int, bool) {
	var result [3]int
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return result, false
	}
	for index, part := range parts {
		if part == "" || len(part) > 1 && part[0] == '0' {
			return result, false
		}
		parsed, err := strconv.Atoi(part)
		if err != nil || parsed < 0 {
			return result, false
		}
		result[index] = parsed
	}
	return result, true
}

func ownershipRefusalReason(path, gap string) string {
	return fmt.Sprintf("consumer release target %s is not proven project-owned: %s; only maintainer_release may mutate a suite package's own metadata", path, gap)
}

// releaseTargetOwnershipGap proves a consumer release target is a project-owned
// source instead of inferring that from the absence of known dependency
// directory names. It returns the ownership evidence that is missing or
// contradicted, or "" when ownership is established.
//
// Two declarations the repository itself owns supply the proof, and both are
// required:
//
//   - Git's index is the repository's own statement of what its sources are, so
//     a target that must already exist has to be tracked. A target the release
//     creates cannot be in the index yet, so its attestation comes from the
//     nearest ancestor directory that already exists — that directory must hold
//     tracked sources — plus `.gitignore`, the repository's own statement of
//     what is not source, not excluding the path.
//   - No directory between the repository root and the target's parent may carry
//     a `SKILL.md` package marker. That marker is how an installed suite package
//     declares it owns its own subtree, and is the same discriminator
//     `repositorymodel.FindRepositoryRoot` and queue-kanban's repo-root walk use
//     to tell an install apart from a directory merely named `do-work`.
//
// Neither half is sufficient alone. A committed vendored package satisfies the
// index claim and is exposed only by its marker; a distribution output or cache
// tree carries no marker and is exposed only by the index refusing to claim it.
// When Git cannot answer at all the target stays unproven and the release
// refuses, which is the fail-closed direction.
func releaseTargetOwnershipGap(repositoryRoot, path string, willCreate bool) string {
	if marker := enclosingPackageMarker(repositoryRoot, path); marker != "" {
		return fmt.Sprintf("%s declares an installed package that owns this subtree", marker)
	}
	if !willCreate {
		if !gitPathTracked(repositoryRoot, path) {
			return "the repository does not track it, so Git does not attest it as a project source"
		}
		return ""
	}
	if gitPathIgnored(repositoryRoot, path) {
		return "the repository's own ignore rules exclude it from its sources"
	}
	ancestor := nearestExistingAncestorDirectory(repositoryRoot, path)
	if !gitDirectoryHoldsTrackedSources(repositoryRoot, ancestor) {
		return fmt.Sprintf("the repository tracks no source in %s, so the new target's location is unattested", ancestor)
	}
	return ""
}

// enclosingPackageMarker returns the repository-relative `SKILL.md` of the
// outermost package that encloses path, or "" when no directory between the
// repository root and the target's parent carries one. The repository root
// itself is never examined: the marker means a nested unit owns its own
// subtree, and a repository that is itself a package still owns its own files.
func enclosingPackageMarker(repositoryRoot, path string) string {
	segments := strings.Split(filepath.ToSlash(path), "/")
	directory := ""
	for _, segment := range segments[:len(segments)-1] {
		if directory == "" {
			directory = segment
		} else {
			directory += "/" + segment
		}
		marker := directory + "/SKILL.md"
		if info, statError := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(marker))); statError == nil && !info.IsDir() {
			return marker
		}
	}
	return ""
}

// nearestExistingAncestorDirectory returns the repository-relative directory
// closest to path that already exists, or "." for the repository root.
func nearestExistingAncestorDirectory(repositoryRoot, path string) string {
	directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(path)))
	for directory != "." && directory != "/" {
		if info, statError := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(directory))); statError == nil && info.IsDir() {
			return directory
		}
		directory = filepath.ToSlash(filepath.Dir(filepath.FromSlash(directory)))
	}
	return "."
}

func gitDirectoryHoldsTrackedSources(repositoryRoot, directory string) bool {
	command := exec.Command("git", "-C", repositoryRoot, "ls-files", "--", directory)
	output, commandError := command.Output()
	return commandError == nil && len(bytes.TrimSpace(output)) > 0
}

func gitPathIgnored(repositoryRoot, path string) bool {
	command := exec.Command("git", "-C", repositoryRoot, "check-ignore", "-q", "--", path)
	return command.Run() == nil
}
