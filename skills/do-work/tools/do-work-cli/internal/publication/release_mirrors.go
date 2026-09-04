package publication

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Undeclared-mirror admission. A release manifest names its version and
// changelog targets by hand, and a hand-written list forgets: the 0.282.0
// release wrote root VERSION and CHANGELOG.md and left the two shipped mirrors
// on the old version until a repair commit. So the planner discovers the mirror
// set itself and refuses a manifest that would leave a mirror behind.
//
// The set is a condition, not a list: every tracked file this planner reads as
// a plain version file that still carries the old version, plus every tracked
// changelog whose bytes equal the declared changelog preimage. It matches the
// rule finalization's recovery-time discovery applies, so plan-time admission
// and recovery never disagree about what a mirror is.

const currentVersionLinePrefix = "**Current version**: "

// undeclaredReleaseMirrors lists the tracked mirrors the manifest did not
// declare, sorted. A failed enumeration or read returns a refusal instead of
// an empty list, because a check that cannot run must not pass.
func undeclaredReleaseMirrors(repositoryRoot, oldVersion string, declared map[string]bool, changelogPreimage []byte) ([]string, *Refusal) {
	tracked, err := trackedRepositoryPaths(repositoryRoot)
	if err != nil {
		return nil, &Refusal{Code: "RELEASE-MIRROR-ENUMERATION", Reason: "tracked-file enumeration failed, so undeclared mirrors cannot be ruled out: " + err.Error()}
	}
	undeclared := []string{}
	for _, path := range tracked {
		if declared[path] || !releaseMirrorCandidate(path) {
			continue
		}
		contents, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if readError != nil {
			return nil, &Refusal{Code: "RELEASE-MIRROR-ENUMERATION", Reason: "tracked release mirror could not be read: " + readError.Error(), Paths: []string{path}}
		}
		if isReleaseMirror(path, contents, oldVersion, changelogPreimage) {
			undeclared = append(undeclared, path)
		}
	}
	sort.Strings(undeclared)
	return undeclared, nil
}

// releaseMirrorCandidate keeps the read to files the mirror rule can apply to.
func releaseMirrorCandidate(path string) bool {
	base := filepath.Base(path)
	return base == "VERSION" || base == "version.md" || strings.HasPrefix(base, "CHANGELOG")
}

func isReleaseMirror(path string, contents []byte, oldVersion string, changelogPreimage []byte) bool {
	if strings.HasPrefix(filepath.Base(path), "CHANGELOG") {
		return changelogPreimage != nil && bytes.Equal(contents, changelogPreimage)
	}
	value, ok := plainVersionFileValue(path, contents)
	return ok && value == oldVersion
}

// plainVersionFileValue reads a VERSION file's whole content or a version
// action file's `**Current version**:` line as a semantic version.
func plainVersionFileValue(path string, contents []byte) (string, bool) {
	var value string
	switch filepath.Base(path) {
	case "VERSION":
		value = strings.TrimSpace(string(contents))
	case "version.md":
		for _, line := range strings.Split(string(contents), "\n") {
			if strings.HasPrefix(line, currentVersionLinePrefix) {
				value = strings.TrimSpace(strings.TrimPrefix(line, currentVersionLinePrefix))
				break
			}
		}
	}
	if _, ok := parseSemver(value); !ok {
		return "", false
	}
	return value, true
}

// trackedRepositoryPaths lists Git-tracked files as slash-separated paths
// relative to the repository root.
func trackedRepositoryPaths(repositoryRoot string) ([]string, error) {
	output, err := exec.Command("git", "-C", repositoryRoot, "ls-files", "-z", "--full-name").Output()
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, entry := range bytes.Split(output, []byte{0}) {
		if len(entry) > 0 {
			paths = append(paths, filepath.ToSlash(string(entry)))
		}
	}
	return paths, nil
}
