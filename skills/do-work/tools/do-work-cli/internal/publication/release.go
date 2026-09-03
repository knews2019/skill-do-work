package publication

import (
	"bytes"
	"os"
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
	if refusal := validateReleaseOwnership(release); refusal != nil {
		plan.Refusal = refusal
		return plan
	}
	for _, target := range release.Targets {
		path, pathError := containedPath(target.Path)
		if pathError != nil {
			return refusedPlan(plan, "RELEASE-PATH-UNSAFE", pathError.Error(), nil, target.Path)
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
	var changelogMirrorBytes []byte
	var changelogMirrorExpected []byte
	var changelogMirrorCreate *bool
	for _, changelog := range release.Changelogs {
		path, pathError := containedPath(changelog.Path)
		if pathError != nil {
			return refusedPlan(plan, "RELEASE-PATH-UNSAFE", pathError.Error(), nil, changelog.Path)
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

func validateReleaseOwnership(release *ReleaseManifest) *Refusal {
	mutationPaths := map[string]bool{}
	mutationOrder := []string{}
	for _, target := range release.Targets {
		path, err := containedPath(target.Path)
		if err != nil {
			return &Refusal{Code: "RELEASE-PATH-UNSAFE", Reason: err.Error(), Paths: []string{target.Path}}
		}
		if !mutationPaths[path] {
			mutationOrder = append(mutationOrder, path)
		}
		mutationPaths[path] = true
	}
	for _, changelog := range release.Changelogs {
		path, err := containedPath(changelog.Path)
		if err != nil {
			return &Refusal{Code: "RELEASE-PATH-UNSAFE", Reason: err.Error(), Paths: []string{changelog.Path}}
		}
		if !mutationPaths[path] {
			mutationOrder = append(mutationOrder, path)
		}
		mutationPaths[path] = true
	}

	projectOwned := map[string]bool{}
	projectOwnedOrder := []string{}
	for _, declaredPath := range release.ProjectOwnedTargets {
		path, err := containedPath(declaredPath)
		if err != nil {
			return &Refusal{Code: "RELEASE-PATH-UNSAFE", Reason: err.Error(), Paths: []string{declaredPath}}
		}
		if projectOwned[path] {
			return &Refusal{Code: "RELEASE-OWNERSHIP-DUPLICATE", Reason: "project-owned release target is declared more than once; declare each normalized path exactly once", Paths: []string{path}}
		}
		projectOwned[path] = true
		projectOwnedOrder = append(projectOwnedOrder, path)
	}

	mirrors := map[string]bool{}
	mirrorOrder := []string{}
	for _, declaredPath := range release.RequiredMirrors {
		path, err := containedPath(declaredPath)
		if err != nil {
			return &Refusal{Code: "RELEASE-PATH-UNSAFE", Reason: err.Error(), Paths: []string{declaredPath}}
		}
		if !release.MaintainerRelease {
			return &Refusal{Code: "RELEASE-MIRROR-ILLEGAL", Reason: "required_mirrors is maintainer-only; classify consumer release targets as project_owned_targets", Paths: []string{path}}
		}
		if mirrors[path] {
			return &Refusal{Code: "RELEASE-OWNERSHIP-DUPLICATE", Reason: "maintainer mirror is declared more than once; declare each normalized path exactly once", Paths: []string{path}}
		}
		if projectOwned[path] {
			return &Refusal{Code: "RELEASE-OWNERSHIP-OVERLAP", Reason: "release target cannot be both project-owned and a maintainer mirror", Paths: []string{path}}
		}
		mirrors[path] = true
		mirrorOrder = append(mirrorOrder, path)
	}

	for _, path := range mutationOrder {
		if !projectOwned[path] && !mirrors[path] {
			return &Refusal{Code: "RELEASE-OWNERSHIP-MISSING", Reason: "release target lacks exact ownership evidence; declare it once in project_owned_targets or, for maintainer releases only, required_mirrors", Paths: []string{path}}
		}
	}
	for _, path := range projectOwnedOrder {
		if !mutationPaths[path] {
			return &Refusal{Code: "RELEASE-OWNERSHIP-EXTRA", Reason: "project-owned declaration does not name a version or changelog target", Paths: []string{path}}
		}
	}
	for _, path := range mirrorOrder {
		if !mutationPaths[path] {
			return &Refusal{Code: "RELEASE-OWNERSHIP-EXTRA", Reason: "maintainer-mirror declaration does not name a version or changelog target", Paths: []string{path}}
		}
	}
	return nil
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
