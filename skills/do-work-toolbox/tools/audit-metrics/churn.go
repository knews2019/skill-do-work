package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// The `churn` and `hotspots` subcommands: how often each file changes, and
// churn × size as the MVP complexity proxy. History comes from one
// `git log -M -C --find-copies-harder --name-status --format=%H
// --since=<window>` pass, parsed newest-first with rename normalization: an R
// entry counts as a touch of the (resolved) new path AND records old→current,
// so every older touch of the old name is attributed to the file's current
// path instead of a dead one. Copy detection covers staged migrations
// (copy-first, delete-the-original-later — how the 2026-08-08 skills/
// restructure actually happened, which -M alone cannot see because the old
// path outlives the copy commit): a dead path that was copy-sourced to a
// surviving file gets its history reassigned to that survivor. Paths deleted
// outright with no surviving copy fall out of the final filter against
// `git ls-files` — only files that still exist can be hotspots. Shallow clones
// are DETECTED and reported in the output, never silently truncated: the
// numbers still print, with a warning that they undercount.

// churnReport is the aggregated result of one history pass.
type churnReport struct {
	ShallowClone  bool
	CommitCount   int
	SinceWindow   string
	TouchesByPath map[string]int
}

// computeChurnReport runs the history pass. Excludes apply to the
// rename-resolved current paths (the ceremony files a caller excludes are named
// by their current paths).
func computeChurnReport(repoRoot string, sinceWindow string, excludePrefixes []string) (churnReport, error) {
	report := churnReport{SinceWindow: sinceWindow, TouchesByPath: map[string]int{}}

	shallowClone, shallowError := isShallowRepository(repoRoot)
	if shallowError != nil {
		return churnReport{}, shallowError
	}
	report.ShallowClone = shallowClone

	// core.quotepath=off keeps non-ASCII paths literal so they match ls-files.
	// --find-copies-harder is what lets a staged migration's copy commit show
	// as C even though its source file was untouched in that commit (~seconds
	// on a repo this size, and the audit runs occasionally, not per keystroke).
	logOutput, logError := runGitCommand(repoRoot, "-c", "core.quotepath=off",
		"log", "-M", "-C", "--find-copies-harder", "--name-status", "--format=%H", "--since="+sinceWindow)
	if logError != nil {
		return churnReport{}, logError
	}

	// Walking newest→oldest, aliasToCurrent maps an old name to the current
	// path it eventually became, so older touches land on the live file.
	// copySourceToTarget remembers where a copy came from; it only matters
	// later, for sources that turn out to be dead today.
	aliasToCurrent := map[string]string{}
	copySourceToTarget := map[string]string{}
	for _, logLine := range strings.Split(logOutput, "\n") {
		if logLine == "" {
			continue
		}
		if isCommitHashLine(logLine) {
			report.CommitCount++
			continue
		}
		statusFields := strings.Split(logLine, "\t")
		statusCode := statusFields[0]
		switch {
		case len(statusCode) > 0 && (statusCode[0] == 'R' || statusCode[0] == 'C') && len(statusFields) >= 3:
			currentPath := resolveCurrentPath(aliasToCurrent, statusFields[2])
			report.TouchesByPath[currentPath]++
			if statusCode[0] == 'R' {
				aliasToCurrent[statusFields[1]] = currentPath
			} else if _, alreadyRecorded := copySourceToTarget[statusFields[1]]; !alreadyRecorded {
				// Nearest-to-present copy wins when a source was copied twice.
				copySourceToTarget[statusFields[1]] = currentPath
			}
		case statusCode == "D" && len(statusFields) >= 2:
			// A deletion is not a touch of living code; if the path never
			// comes back, the tracked-set filter below drops its history.
		case len(statusFields) >= 2:
			report.TouchesByPath[resolveCurrentPath(aliasToCurrent, statusFields[1])]++
		}
	}

	// Two distinct sets, deliberately: livePathSet (every tracked path, no
	// excludes) answers "is this path alive today?" for the reassignment
	// below; reportPathSet (excludes applied) decides what the report shows.
	// Building one excluded set for both jobs made an excluded-but-live copy
	// source look dead, silently transferring its whole history to the
	// surviving copy (PR #139 review finding).
	trackedPaths, listError := listTrackedFiles(repoRoot)
	if listError != nil {
		return churnReport{}, listError
	}
	livePathSet := map[string]bool{}
	for _, trackedPath := range trackedPaths {
		livePathSet[trackedPath] = true
	}
	reportPathSet := map[string]bool{}
	for _, trackedPath := range applyPathExcludes(trackedPaths, excludePrefixes) {
		reportPathSet[trackedPath] = true
	}
	// Staged-migration reassignment: a path that is dead today but was
	// copy-sourced to a surviving file hands its history to that survivor —
	// the copy was the real "rename", just split across commits. Aliveness is
	// judged against livePathSet: an excluded-but-live path keeps its own
	// history (and is then simply not reported). Paths deleted outright with
	// no surviving copy drop in the filter below.
	for touchedPath, touchCount := range report.TouchesByPath {
		if livePathSet[touchedPath] {
			continue
		}
		if survivingPath, found := resolveSurvivingCopy(copySourceToTarget, livePathSet, touchedPath); found {
			report.TouchesByPath[survivingPath] += touchCount
		}
	}
	for touchedPath := range report.TouchesByPath {
		if !reportPathSet[touchedPath] {
			delete(report.TouchesByPath, touchedPath)
		}
	}
	return report, nil
}

// resolveSurvivingCopy follows the copy chain from a dead path to a file that
// is tracked today. The hop bound guards against a pathological copy cycle.
func resolveSurvivingCopy(copySourceToTarget map[string]string, trackedSet map[string]bool, deadPath string) (string, bool) {
	candidatePath := deadPath
	for hop := 0; hop <= len(copySourceToTarget); hop++ {
		targetPath, hasCopy := copySourceToTarget[candidatePath]
		if !hasCopy {
			return "", false
		}
		if trackedSet[targetPath] {
			return targetPath, true
		}
		candidatePath = targetPath
	}
	return "", false
}

// resolveCurrentPath follows the alias chain to the path's current name. The
// hop bound guards against a pathological rename cycle in the parsed history.
func resolveCurrentPath(aliasToCurrent map[string]string, historicalPath string) string {
	currentPath := historicalPath
	for hop := 0; hop <= len(aliasToCurrent); hop++ {
		nextPath, hasAlias := aliasToCurrent[currentPath]
		if !hasAlias {
			return currentPath
		}
		currentPath = nextPath
	}
	return currentPath
}

// isCommitHashLine reports whether a log line is a %H commit hash (40 hex for
// SHA-1, 64 for SHA-256) rather than a name-status entry.
func isCommitHashLine(logLine string) bool {
	if len(logLine) != 40 && len(logLine) != 64 {
		return false
	}
	for _, character := range logLine {
		isHexDigit := (character >= '0' && character <= '9') || (character >= 'a' && character <= 'f')
		if !isHexDigit {
			return false
		}
	}
	return true
}

// sortedChurnEntries returns path/touch pairs sorted by touches descending,
// path ascending on ties — a stable order the audit can diff between runs.
type churnEntry struct {
	Path    string
	Touches int
}

func sortedChurnEntries(touchesByPath map[string]int) []churnEntry {
	var churnEntries []churnEntry
	for touchedPath, touchCount := range touchesByPath {
		churnEntries = append(churnEntries, churnEntry{Path: touchedPath, Touches: touchCount})
	}
	sort.Slice(churnEntries, func(left, right int) bool {
		if churnEntries[left].Touches != churnEntries[right].Touches {
			return churnEntries[left].Touches > churnEntries[right].Touches
		}
		return churnEntries[left].Path < churnEntries[right].Path
	})
	return churnEntries
}

// writeShallowWarning prints the shallow-clone warning line. Its own function
// so churn and hotspots stay word-for-word identical.
func writeShallowWarning(outputWriter io.Writer, shallowClone bool) {
	if shallowClone {
		fmt.Fprintf(outputWriter, "> WARNING: shallow clone detected — history is truncated at the shallow boundary; counts below UNDERCOUNT reality.\n\n")
	}
}

// writeChurnReport renders the top-N churn table.
func writeChurnReport(outputWriter io.Writer, report churnReport, topCount int) {
	fmt.Fprintf(outputWriter, "## Churn — commits touching each file (since %s)\n\n", report.SinceWindow)
	writeShallowWarning(outputWriter, report.ShallowClone)
	fmt.Fprintf(outputWriter, "Commits scanned: %d. Rename-normalized; deleted paths dropped.\n\n", report.CommitCount)
	fmt.Fprintf(outputWriter, "| path | commits |\n|---|---:|\n")
	for index, entry := range sortedChurnEntries(report.TouchesByPath) {
		if index >= topCount {
			break
		}
		fmt.Fprintf(outputWriter, "| %s | %d |\n", entry.Path, entry.Touches)
	}
}

// hotspotEntry is one churn × size row.
type hotspotEntry struct {
	Path         string
	Touches      int
	Lines        int
	HotspotScore int
}

// computeHotspotEntries joins churn touches with current line counts. Size is
// the MVP complexity proxy; binary files measure zero lines and so score zero.
func computeHotspotEntries(repoRoot string, report churnReport) ([]hotspotEntry, error) {
	var hotspotEntries []hotspotEntry
	for touchedPath, touchCount := range report.TouchesByPath {
		measurement, measureError := measureTrackedFile(repoRoot, touchedPath)
		if measureError != nil {
			continue // tracked but unreadable in this worktree — no size, no score
		}
		hotspotEntries = append(hotspotEntries, hotspotEntry{
			Path:         touchedPath,
			Touches:      touchCount,
			Lines:        measurement.Lines,
			HotspotScore: touchCount * measurement.Lines,
		})
	}
	sort.Slice(hotspotEntries, func(left, right int) bool {
		if hotspotEntries[left].HotspotScore != hotspotEntries[right].HotspotScore {
			return hotspotEntries[left].HotspotScore > hotspotEntries[right].HotspotScore
		}
		return hotspotEntries[left].Path < hotspotEntries[right].Path
	})
	return hotspotEntries, nil
}

// writeHotspotsReport renders the top-N churn × size table.
func writeHotspotsReport(outputWriter io.Writer, report churnReport, hotspotEntries []hotspotEntry, topCount int) {
	fmt.Fprintf(outputWriter, "## Hotspots — churn × size (since %s)\n\n", report.SinceWindow)
	writeShallowWarning(outputWriter, report.ShallowClone)
	fmt.Fprintf(outputWriter, "Commits scanned: %d. Score = commits × current lines.\n\n", report.CommitCount)
	fmt.Fprintf(outputWriter, "| path | commits | lines | score |\n|---|---:|---:|---:|\n")
	for index, entry := range hotspotEntries {
		if index >= topCount {
			break
		}
		fmt.Fprintf(outputWriter, "| %s | %d | %d | %d |\n", entry.Path, entry.Touches, entry.Lines, entry.HotspotScore)
	}
}
