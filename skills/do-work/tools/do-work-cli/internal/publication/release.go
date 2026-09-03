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
			if gap := releaseTargetOwnershipGap(repositoryRoot, path, false); !gap.ownershipProven() {
				return refusedPlan(plan, "RELEASE-TARGET-OWNERSHIP-UNVERIFIED", gap.refusalReason(path), nil, path)
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
			if gap := releaseTargetOwnershipGap(repositoryRoot, path, changelog.Create); !gap.ownershipProven() {
				return refusedPlan(plan, "RELEASE-TARGET-OWNERSHIP-UNVERIFIED", gap.refusalReason(path), nil, path)
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

// releaseOwnershipGap names the ownership declaration a release target is
// missing and the action that supplies it. The zero value means ownership is
// proven. The remedy travels with the gap rather than being appended to every
// refusal, because `maintainer_release` resolves exactly one of these gaps: a
// caller whose target is merely uncommitted or ignored needs to commit it or
// drop the ignore rule, and offering the maintainer escape hatch there would
// steer them past the fix their gap actually needs.
type releaseOwnershipGap struct {
	MissingEvidence string
	Remedy          string
}

func (gap releaseOwnershipGap) ownershipProven() bool {
	return gap.MissingEvidence == ""
}

func (gap releaseOwnershipGap) refusalReason(path string) string {
	return fmt.Sprintf("consumer release target %s is not proven project-owned: %s; %s", path, gap.MissingEvidence, gap.Remedy)
}

// releaseTargetOwnershipGap proves a consumer release target is a project-owned
// source instead of inferring that from the absence of known dependency
// directory names. It returns the ownership evidence that is missing or
// contradicted, or the zero gap when ownership is established.
//
// Two declarations the repository itself owns supply the proof, and both are
// required:
//
//   - Git's index is the repository's own statement of what its sources are, so
//     a target that must already exist has to be tracked. A target the release
//     creates cannot be in the index yet, so its attestation comes from the
//     directory it is created in: that directory must already exist and the
//     index must claim a source inside it, and `.gitignore` — the repository's
//     own statement of what is not source — must not exclude the path. The
//     attestation deliberately does not walk further up the tree, because the
//     repository root always holds tracked sources: a root-anchored walk
//     attests every not-yet-built `dist/`, cache or install tree, which is the
//     state a real repository is in before those trees exist.
//   - No directory between the repository root and the target's parent may carry
//     a `SKILL.md` package marker. That marker is how an installed suite package
//     declares it owns its own subtree, and is the same discriminator
//     `repositorymodel.FindRepositoryRoot` and queue-kanban's repo-root walk use
//     to tell an install apart from a directory merely named `do-work`. It is
//     checked for created and existing targets alike.
//
// Neither half is sufficient alone. A committed vendored package satisfies the
// index claim and is exposed only by its marker; a distribution output or cache
// tree carries no marker and is exposed only by the index refusing to claim it.
// When Git cannot answer at all the target stays unproven and the release
// refuses, which is the fail-closed direction.
func releaseTargetOwnershipGap(repositoryRoot, path string, willCreate bool) releaseOwnershipGap {
	if marker := enclosingPackageMarker(repositoryRoot, path); marker != "" {
		return releaseOwnershipGap{
			MissingEvidence: fmt.Sprintf("%s declares an installed package that owns this subtree", marker),
			Remedy:          "only maintainer_release may mutate a suite package's own metadata",
		}
	}
	if !willCreate {
		if !gitPathTracked(repositoryRoot, path) {
			return releaseOwnershipGap{
				MissingEvidence: "the repository does not track it, so Git does not attest it as a project source",
				Remedy:          "commit the target so Git attests it as a project source, then release again",
			}
		}
		return releaseOwnershipGap{}
	}
	if gitPathIgnored(repositoryRoot, path) {
		return releaseOwnershipGap{
			MissingEvidence: "the repository's own ignore rules exclude it from its sources",
			Remedy:          "drop the ignore rule covering the target, or release a target the repository claims as a source",
		}
	}
	parentDirectory := releaseTargetParentDirectory(path)
	parentDescription := releaseDirectoryDescription(parentDirectory)
	locationRemedy := fmt.Sprintf("commit a project source in %s first, or create the target where the repository already tracks sources", parentDescription)
	if !directoryExists(repositoryRoot, parentDirectory) {
		return releaseOwnershipGap{
			MissingEvidence: fmt.Sprintf("its parent directory %s does not exist, so no project directory attests the new target's location", parentDescription),
			Remedy:          locationRemedy,
		}
	}
	if !gitDirectoryHoldsTrackedSources(repositoryRoot, parentDirectory) {
		return releaseOwnershipGap{
			MissingEvidence: fmt.Sprintf("the repository tracks no source in %s, so the new target's location is unattested", parentDescription),
			Remedy:          locationRemedy,
		}
	}
	return releaseOwnershipGap{}
}

// enclosingPackageMarker returns the repository-relative `SKILL.md` of the
// outermost package that encloses path, or "" when no directory between the
// repository root and the target's parent carries one. The repository root
// itself is never examined: the marker means a nested unit owns its own
// subtree, and a repository that is itself a package still owns its own files.
//
// Any filesystem entry at that name counts. A legitimate install cannot carry a
// `SKILL.md` directory — both `suitemanifest.ValidateSuite` and the install
// transaction require a non-empty regular file — and on a case-insensitive
// filesystem an ancestor spelled `skill.md` is seen as a marker too. Both cases
// refuse a release that could have been allowed rather than admitting one that
// should not be, which is the direction to err in here.
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
		if _, statError := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(marker))); statError == nil {
			return marker
		}
	}
	return ""
}

// releaseTargetParentDirectory returns the repository-relative directory the
// target sits directly in, or "." when the target sits at the repository root.
func releaseTargetParentDirectory(path string) string {
	parent := filepath.ToSlash(filepath.Dir(filepath.FromSlash(path)))
	if parent == "/" {
		return "."
	}
	return parent
}

func releaseDirectoryDescription(directory string) string {
	if directory == "." {
		return "the repository root"
	}
	return directory
}

func directoryExists(repositoryRoot, directory string) bool {
	info, statError := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(directory)))
	return statError == nil && info.IsDir()
}

// gitDirectoryHoldsTrackedSources reports whether the index claims at least one
// path inside directory. An index entry equal to the directory itself — a
// submodule gitlink, or a tracked symlink standing where a directory appears to
// be — is not a source the repository keeps in that directory and attests
// nothing about a target created beneath it.
func gitDirectoryHoldsTrackedSources(repositoryRoot, directory string) bool {
	command := exec.Command("git", "-C", repositoryRoot, "ls-files", "-z", "--", directory)
	output, commandError := command.Output()
	if commandError != nil {
		return false
	}
	insidePrefix := ""
	if directory != "." {
		insidePrefix = directory + "/"
	}
	for _, entry := range strings.Split(string(output), "\x00") {
		if entry != "" && entry != directory && strings.HasPrefix(entry, insidePrefix) {
			return true
		}
	}
	return false
}

func gitPathIgnored(repositoryRoot, path string) bool {
	command := exec.Command("git", "-C", repositoryRoot, "check-ignore", "-q", "--", path)
	return command.Run() == nil
}
