package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// requestFileReference is one discovered REQ-*.md path plus which section of the
// do-work tree it was found under ("queue", "working", or "archive"). The
// section is recorded for provenance; status — not location — drives column
// bucketing.
type requestFileReference struct {
	AbsolutePath string
	TreeSection  string
}

// strayRequestFile is a REQ-*.md discovered UNDER do-work/ but OUTSIDE the
// scanned sections (queue/, working/, archive/) — e.g. a work agent that
// archived to do-work/user-requests/UR-NNN/ instead of do-work/archive/. Such a
// file has no card and is invisible on the board; buildBoard turns each one into
// a data warning so a misplaced REQ is flagged, never silently dropped.
type strayRequestFile struct {
	AbsolutePath string
	RelativePath string // path relative to do-work/, for a location-pinpointing warning
}

// discoveredTreeFiles is the raw output of walking the do-work tree: every REQ
// file reference, every UR input.md path, and the do-work/notes.md path when it
// exists — all before any parsing.
type discoveredTreeFiles struct {
	RequestFiles      []requestFileReference
	StrayRequestFiles []strayRequestFile // REQ-*.md found outside queue/working/archive — flagged, not shown as cards
	UserRequestFiles  []string
	NotesFilePath     string // absolute path to do-work/notes.md, "" when absent
	TestersFilePath   string // absolute path to do-work/testers.md, "" when absent
}

// resolveRepoRoot walks upward from startDirectory until it finds a directory
// that contains a `do-work/` subdirectory, and returns that directory as the
// repo root. A `do-work/` that is a skill install rather than a queue tree
// (see isSkillInstallDirectory) is skipped and the walk continues upward. It
// errors only if the filesystem root is reached without finding one.
func resolveRepoRoot(startDirectory string) (string, error) {
	currentDirectory := startDirectory
	for {
		candidate := filepath.Join(currentDirectory, "do-work")
		if info, statError := os.Stat(candidate); statError == nil && info.IsDir() && !isSkillInstallDirectory(candidate) {
			return currentDirectory, nil
		}
		parentDirectory := filepath.Dir(currentDirectory)
		if parentDirectory == currentDirectory {
			return "", fmt.Errorf("queue-kanban: no do-work/ directory found walking up from %s", startDirectory)
		}
		currentDirectory = parentDirectory
	}
}

// isSkillInstallDirectory reports whether a directory named `do-work` is the
// do-work skill's install tree rather than a queue tree. Consumer repos vendor
// the skill at a path itself named do-work (e.g. `.claude/skills/do-work/`), so
// an upward walk from the vendored tool would otherwise match the install's
// parent and silently build an empty board while the real queue sits further
// up. SKILL.md at the directory's top level is the discriminator: every skill
// install ships it; a queue tree never contains one.
func isSkillInstallDirectory(candidateDirectory string) bool {
	_, statError := os.Stat(filepath.Join(candidateDirectory, "SKILL.md"))
	return statError == nil
}

// deriveProjectName returns a human-facing project name for a repo root: the base
// name of its absolute path (e.g. "/Users/t2/2code/g1w-game-find-the-difference"
// → "g1w-game-find-the-difference"). It falls back to the un-absolutized base when
// filepath.Abs fails, and to "do-work" when even that collapses to "." or "/", so
// the board title is never blank.
func deriveProjectName(repoRoot string) string {
	resolvedRoot := repoRoot
	if absoluteRoot, absError := filepath.Abs(repoRoot); absError == nil {
		resolvedRoot = absoluteRoot
	}
	projectName := filepath.Base(resolvedRoot)
	if projectName == "." || projectName == string(filepath.Separator) {
		return "do-work"
	}
	return projectName
}

// resolveRepoRootOrDefault returns the override when it is non-empty, otherwise
// it resolves the repo root by walking up from the current working directory.
func resolveRepoRootOrDefault(repoRootOverride string) (string, error) {
	if strings.TrimSpace(repoRootOverride) != "" {
		return repoRootOverride, nil
	}
	workingDirectory, getwdError := os.Getwd()
	if getwdError != nil {
		return "", getwdError
	}
	return resolveRepoRoot(workingDirectory)
}

// enumerateDoWorkTree walks repoRoot/do-work and collects every REQ-*.md file
// (from queue/, working/, and the entire archive/** subtree — handling both the
// flat archive/UR-NNN/ shape and the banded archive/UR-NNN-MMM/ shape with its
// nested UR-NNN/ subfolders) plus every UR input.md (from user-requests/** and
// archive/**) and the top-level notes.md written by `do-work note`.
//
// The do-work/deliverables/ and do-work/runs/ subtrees are skipped entirely.
// The kb/wiki/sources/ mirror lives OUTSIDE do-work and so is never reached —
// walking only under do-work is what keeps the REQ count from roughly doubling.
func enumerateDoWorkTree(repoRoot string) (discoveredTreeFiles, error) {
	var discovered discoveredTreeFiles

	doWorkDirectory := filepath.Join(repoRoot, "do-work")
	info, statError := os.Stat(doWorkDirectory)
	if statError != nil || !info.IsDir() {
		return discovered, fmt.Errorf("queue-kanban: do-work directory not found at %s", doWorkDirectory)
	}

	walkError := filepath.WalkDir(doWorkDirectory, func(path string, dirEntry fs.DirEntry, entryError error) error {
		if entryError != nil {
			// Best-effort: skip an unreadable entry rather than aborting the whole walk.
			return nil
		}

		relativePath, relativeError := filepath.Rel(doWorkDirectory, path)
		if relativeError != nil {
			return nil
		}
		topSection := strings.Split(relativePath, string(filepath.Separator))[0]

		if dirEntry.IsDir() {
			if path != doWorkDirectory && isSkippedSection(topSection, dirEntry.Name()) {
				return fs.SkipDir
			}
			return nil
		}

		baseName := dirEntry.Name()
		switch {
		case strings.HasPrefix(baseName, "REQ-") && strings.HasSuffix(baseName, ".md"):
			switch topSection {
			case "queue", "working", "archive":
				discovered.RequestFiles = append(discovered.RequestFiles, requestFileReference{
					AbsolutePath: path,
					TreeSection:  topSection,
				})
			default:
				// A REQ file that the walk reached but that lives outside the
				// scanned sections (the pruned deliverables/ and runs/ subtrees
				// never get here — they are SkipDir'd as directories). It would
				// otherwise vanish; record it so buildBoard can flag it.
				discovered.StrayRequestFiles = append(discovered.StrayRequestFiles, strayRequestFile{
					AbsolutePath: path,
					RelativePath: relativePath,
				})
			}
		case baseName == "input.md":
			switch topSection {
			case "user-requests", "archive":
				discovered.UserRequestFiles = append(discovered.UserRequestFiles, path)
			}
		case relativePath == "notes.md":
			// Only the notes file at the top level of do-work/ — a notes.md
			// nested under a UR or archive folder is somebody's scratch file,
			// not the queue's note list.
			discovered.NotesFilePath = path
		case relativePath == testersFileRelativePath:
			// Same top-level-only rule as notes.md: do-work/testers.md is the
			// testing view's profile store (see testing.go).
			discovered.TestersFilePath = path
		}
		return nil
	})

	return discovered, walkError
}

// isSkippedSection reports whether a directory should be pruned from the walk.
// The deliverables (reports, not REQs) and runs (run logs) sections are excluded
// per the data model, and any hidden directory (a leading dot, e.g. .git) is
// skipped defensively.
func isSkippedSection(topSection string, directoryName string) bool {
	if topSection == "deliverables" || topSection == "runs" {
		return true
	}
	return strings.HasPrefix(directoryName, ".")
}
