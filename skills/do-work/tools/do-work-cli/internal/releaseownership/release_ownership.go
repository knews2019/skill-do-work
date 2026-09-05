// Package releaseownership answers one question for every release reader:
// which tracked files does this repository's release actually own?
//
// Two readers ask it. The publication planner asks before refusing a manifest
// that would leave a mirror behind on the old version, and finalization's
// recovery discovery asks before associating release metadata with a journal.
// They must never disagree about what a mirror is, so the rule lives here once.
//
// Ownership is affirmative evidence: the repository root, a suite module
// declared in suite/modules.tsv, or a workspace member a tracked project
// manifest claims. A shared version number is NOT evidence — an independently
// versioned component may legitimately sit on the same value as the
// application, and reading that coincidence as membership refuses a valid
// release or drags an unrelated component into it.
package releaseownership

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ReadImage returns the bytes of a repository-relative path in the image the
// caller judges, and whether that image holds the file. Finalization reads
// HEAD, because it classifies work already committed; the release planner
// reads the working tree, because that is the state it is about to rewrite.
type ReadImage func(path string) ([]byte, bool)

// Ownership partitions a repository's tracked files by affirmative release
// ownership evidence.
type Ownership struct {
	// MetadataPaths holds every release-metadata path the release owns.
	MetadataPaths map[string]bool
	// ProjectManifests holds every package.json / Cargo.toml / pyproject.toml
	// the release owns, which is what selects a workspace lock entry.
	ProjectManifests map[string]bool
}

// IsReleaseMetadataPath reports whether a path carries release metadata at all.
// It is a shape question, asked before ownership.
func IsReleaseMetadataPath(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, "CHANGELOG") || base == "VERSION" || base == "package.json" || base == "package-lock.json" || base == "Cargo.toml" || base == "Cargo.lock" || base == "pyproject.toml" || base == "uv.lock" || path == "skills/do-work/actions/version.md"
}

// NormalizedDirectory returns a path's containing directory, with the
// repository root spelled as the empty string.
func NormalizedDirectory(path string) string {
	directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(path)))
	if directory == "." {
		return ""
	}
	return directory
}

// AffirmativeOwnership classifies every tracked path against the repository's
// declared release topology. An unreadable declaration is an error, never an
// empty answer: a classification that cannot run must not read as "unowned".
func AffirmativeOwnership(tracked []string, trackedSet map[string]bool, readImage ReadImage) (Ownership, error) {
	ownedRoots := map[string]bool{"": true}
	ownedManifests := map[string]bool{}
	maintainerRoots, err := declaredMaintainerReleaseRoots(trackedSet, readImage)
	if err != nil {
		return Ownership{}, err
	}
	for _, manifestPath := range tracked {
		manifestPath = filepath.ToSlash(filepath.Clean(manifestPath))
		base := filepath.Base(manifestPath)
		if base != "package.json" && base != "Cargo.toml" && base != "pyproject.toml" {
			continue
		}
		if _, exists := readImage(manifestPath); !exists {
			return Ownership{}, fmt.Errorf("read tracked project manifest %s", manifestPath)
		}
		directory := NormalizedDirectory(manifestPath)
		if directory == "" || pathWithinReleaseRoots(manifestPath, maintainerRoots) {
			ownedRoots[directory] = true
			ownedManifests[manifestPath] = true
		}
	}
	for promoted := true; promoted; {
		promoted = false
		for _, manifestPath := range tracked {
			manifestPath = filepath.ToSlash(filepath.Clean(manifestPath))
			base := filepath.Base(manifestPath)
			if base != "package.json" && base != "Cargo.toml" && base != "pyproject.toml" || ownedManifests[manifestPath] {
				continue
			}
			if _, _, ok := FindOwnedWorkspaceOwner(trackedSet, ownedManifests, manifestPath, base, readImage); ok {
				directory := NormalizedDirectory(manifestPath)
				ownedRoots[directory] = true
				ownedManifests[manifestPath] = true
				promoted = true
			}
		}
	}

	ownedPaths := map[string]bool{}
	for _, path := range tracked {
		path = filepath.ToSlash(filepath.Clean(path))
		maintainerOwned := pathWithinReleaseRoots(path, maintainerRoots)
		if IsReleaseMetadataPath(path) && (ownedRoots[NormalizedDirectory(path)] || maintainerOwned) {
			ownedPaths[path] = true
		}
	}
	return Ownership{MetadataPaths: ownedPaths, ProjectManifests: ownedManifests}, nil
}

func pathWithinReleaseRoots(path string, roots []string) bool {
	for _, root := range roots {
		if path == root || strings.HasPrefix(path, root+"/") {
			return true
		}
	}
	return false
}

// declaredMaintainerReleaseRoots reads suite/modules.tsv, the suite's only
// source declaration for its sibling packages. A module is a release root only
// when it carries its own tracked VERSION file.
func declaredMaintainerReleaseRoots(trackedSet map[string]bool, readImage ReadImage) ([]string, error) {
	if !trackedSet["suite/modules.tsv"] {
		return nil, nil
	}
	contents, exists := readImage("suite/modules.tsv")
	if !exists {
		return nil, fmt.Errorf("read tracked suite/modules.tsv")
	}
	roots := []string{}
	lines := strings.Split(strings.TrimSpace(string(contents)), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "source\tdestination" {
		return nil, fmt.Errorf("suite/modules.tsv has no source/destination header")
	}
	for _, line := range lines[1:] {
		fields := strings.Split(line, "\t")
		if len(fields) != 2 {
			return nil, fmt.Errorf("suite/modules.tsv has an invalid module row")
		}
		source := filepath.ToSlash(filepath.Clean(strings.TrimSpace(fields[0])))
		if source == "." || source == ".." || strings.HasPrefix(source, "../") || !trackedSet[source+"/VERSION"] {
			continue
		}
		roots = append(roots, source)
	}
	return sortedUniqueRoots(roots), nil
}

// FindOwnedWorkspaceOwner walks up from a project manifest to the nearest
// already-owned manifest of the same ecosystem that claims it as a workspace
// member, returning that owner's path and the member's path relative to it.
func FindOwnedWorkspaceOwner(tracked, ownedManifests map[string]bool, manifestPath, manifestName string, readImage ReadImage) (string, string, bool) {
	memberDirectory := filepath.ToSlash(filepath.Dir(manifestPath))
	for ownerDirectory := filepath.ToSlash(filepath.Dir(memberDirectory)); ; ownerDirectory = filepath.ToSlash(filepath.Dir(ownerDirectory)) {
		if ownerDirectory == "." {
			ownerDirectory = ""
		}
		ownerPath := filepath.ToSlash(filepath.Join(ownerDirectory, manifestName))
		if tracked[ownerPath] && ownedManifests[ownerPath] {
			if ownerBytes, exists := readImage(ownerPath); exists {
				relative, relativeError := filepath.Rel(filepath.Dir(filepath.FromSlash(ownerPath)), filepath.FromSlash(memberDirectory))
				if relativeError == nil {
					relative = filepath.ToSlash(relative)
					for _, pattern := range workspaceMemberPatterns(manifestName, ownerBytes) {
						if workspacePatternMatches(pattern, relative) {
							return ownerPath, relative, true
						}
					}
				}
			}
		}
		if ownerDirectory == "" {
			break
		}
	}
	return "", "", false
}

func workspaceMemberPatterns(manifestName string, contents []byte) []string {
	if manifestName == "package.json" {
		var document struct {
			Workspaces json.RawMessage `json:"workspaces"`
		}
		if json.Unmarshal(contents, &document) != nil || len(document.Workspaces) == 0 {
			return nil
		}
		var direct []string
		if json.Unmarshal(document.Workspaces, &direct) == nil {
			return direct
		}
		var nested struct {
			Packages []string `json:"packages"`
		}
		if json.Unmarshal(document.Workspaces, &nested) == nil {
			return nested.Packages
		}
		return nil
	}
	section := "workspace"
	if manifestName == "pyproject.toml" {
		section = "tool.uv.workspace"
	}
	return tomlSectionStringArray(contents, section, "members")
}

func workspacePatternMatches(pattern, memberPath string) bool {
	pattern = strings.Trim(strings.TrimSpace(filepath.ToSlash(pattern)), "/")
	memberPath = strings.Trim(strings.TrimSpace(filepath.ToSlash(memberPath)), "/")
	if pattern == memberPath {
		return true
	}
	if strings.HasSuffix(pattern, "/**") {
		return strings.HasPrefix(memberPath+"/", strings.TrimSuffix(pattern, "**"))
	}
	matched, err := filepath.Match(filepath.FromSlash(pattern), filepath.FromSlash(memberPath))
	return err == nil && matched
}

func tomlSectionStringArray(contents []byte, section, key string) []string {
	sectionPattern := regexp.MustCompile(`(?m)^\[` + regexp.QuoteMeta(section) + `\][ \t]*\r?$`)
	location := sectionPattern.FindIndex(contents)
	if location == nil {
		return nil
	}
	sectionBytes := contents[location[1]:]
	if next := regexp.MustCompile(`(?m)^\[[^\r\n]+\][ \t]*\r?$`).FindIndex(sectionBytes); next != nil {
		sectionBytes = sectionBytes[:next[0]]
	}
	arrayPattern := regexp.MustCompile(`(?s)(?:^|\n)[ \t]*` + regexp.QuoteMeta(key) + `[ \t]*=[ \t]*\[(.*?)\]`)
	match := arrayPattern.FindSubmatch(sectionBytes)
	if len(match) != 2 {
		return nil
	}
	quoted := regexp.MustCompile(`["']([^"']+)["']`).FindAllSubmatch(match[1], -1)
	values := make([]string, 0, len(quoted))
	for _, value := range quoted {
		values = append(values, string(value[1]))
	}
	return values
}

func sortedUniqueRoots(roots []string) []string {
	seen := map[string]bool{}
	unique := []string{}
	for _, root := range roots {
		if !seen[root] {
			seen[root] = true
			unique = append(unique, root)
		}
	}
	sort.Strings(unique)
	return unique
}
