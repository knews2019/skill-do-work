package main

import (
	"fmt"
	"io"
	"path"
	"sort"
)

// The `inventory` and `folders` subcommands: tracked-file counts and sizes.
// Inventory aggregates by extension, summarizes file-lines/file-words
// distributions, lists the top-N largest files, and — only when band flags are
// passed — the files over each threshold. Folders does the same for
// files-per-folder (direct children of each directory, not recursive: the
// audit asks "which folder is crowded", and recursive counts would charge every
// parent for its whole subtree). Renderers take io.Writer so the markdown
// shape is assertable; the command wrappers in main.go own flags and os.Exit.

// inventoryReport is everything computed from one tracked-file walk.
type inventoryReport struct {
	Files           []fileMeasurement
	BinaryCount     int
	UnreadableCount int
}

// computeInventoryReport lists tracked files (excludes applied) and measures
// each one. Tracked-but-unreadable files (e.g. deleted from the worktree) are
// skipped and counted rather than failing the whole run — the audit wants the
// numbers that exist, with the gap visible.
func computeInventoryReport(repoRoot string, excludePrefixes []string) (inventoryReport, error) {
	trackedPaths, listError := listTrackedFiles(repoRoot)
	if listError != nil {
		return inventoryReport{}, listError
	}
	report := inventoryReport{}
	for _, trackedPath := range applyPathExcludes(trackedPaths, excludePrefixes) {
		measurement, measureError := measureTrackedFile(repoRoot, trackedPath)
		if measureError != nil {
			report.UnreadableCount++
			continue
		}
		if measurement.Binary {
			report.BinaryCount++
		}
		report.Files = append(report.Files, measurement)
	}
	return report, nil
}

// textFiles returns the non-binary measurements — the only ones whose
// line/word numbers mean anything.
func (report inventoryReport) textFiles() []fileMeasurement {
	var textOnly []fileMeasurement
	for _, measurement := range report.Files {
		if !measurement.Binary {
			textOnly = append(textOnly, measurement)
		}
	}
	return textOnly
}

// extensionAggregate is one row of the by-extension table.
type extensionAggregate struct {
	Extension string
	FileCount int
	LineTotal int
	WordTotal int
}

// aggregateByExtension groups measurements by filename extension. Binary files
// count as files but contribute no lines/words.
func aggregateByExtension(measurements []fileMeasurement) []extensionAggregate {
	aggregatesByExtension := map[string]*extensionAggregate{}
	for _, measurement := range measurements {
		extension := path.Ext(measurement.Path)
		if extension == "" {
			extension = "(none)"
		}
		aggregate, exists := aggregatesByExtension[extension]
		if !exists {
			aggregate = &extensionAggregate{Extension: extension}
			aggregatesByExtension[extension] = aggregate
		}
		aggregate.FileCount++
		aggregate.LineTotal += measurement.Lines
		aggregate.WordTotal += measurement.Words
	}
	var aggregates []extensionAggregate
	for _, aggregate := range aggregatesByExtension {
		aggregates = append(aggregates, *aggregate)
	}
	sort.Slice(aggregates, func(left, right int) bool {
		if aggregates[left].LineTotal != aggregates[right].LineTotal {
			return aggregates[left].LineTotal > aggregates[right].LineTotal
		}
		return aggregates[left].Extension < aggregates[right].Extension
	})
	return aggregates
}

// fileBandThresholds carries the inventory band flags. Any field left at
// bandThresholdUnset simply never fires; hasAny gates the whole section.
type fileBandThresholds struct {
	WatchLines int
	FlagLines  int
	WatchWords int
	FlagWords  int
}

func (thresholds fileBandThresholds) hasAny() bool {
	return thresholds.WatchLines != bandThresholdUnset || thresholds.FlagLines != bandThresholdUnset ||
		thresholds.WatchWords != bandThresholdUnset || thresholds.FlagWords != bandThresholdUnset
}

// writeInventoryReport renders the full inventory as pasteable markdown. The
// band section appears ONLY when a band flag was passed — no flag, no band.
func writeInventoryReport(outputWriter io.Writer, report inventoryReport, thresholds fileBandThresholds, topCount int) {
	fmt.Fprintf(outputWriter, "## Inventory — tracked files by extension\n\n")
	fmt.Fprintf(outputWriter, "| extension | files | lines | words |\n|---|---:|---:|---:|\n")
	lineTotal, wordTotal := 0, 0
	for _, aggregate := range aggregateByExtension(report.Files) {
		fmt.Fprintf(outputWriter, "| %s | %d | %d | %d |\n", aggregate.Extension, aggregate.FileCount, aggregate.LineTotal, aggregate.WordTotal)
		lineTotal += aggregate.LineTotal
		wordTotal += aggregate.WordTotal
	}
	fmt.Fprintf(outputWriter, "| **total** | %d | %d | %d |\n", len(report.Files), lineTotal, wordTotal)
	if report.BinaryCount > 0 {
		fmt.Fprintf(outputWriter, "\n%d binary file(s) excluded from line/word counts.\n", report.BinaryCount)
	}
	if report.UnreadableCount > 0 {
		fmt.Fprintf(outputWriter, "\nWARNING: %d tracked file(s) unreadable and skipped.\n", report.UnreadableCount)
	}

	textMeasurements := report.textFiles()
	var lineValues, wordValues []int
	for _, measurement := range textMeasurements {
		lineValues = append(lineValues, measurement.Lines)
		wordValues = append(wordValues, measurement.Words)
	}
	fmt.Fprintf(outputWriter, "\n## File-Size Distributions (non-binary tracked files)\n\n")
	fmt.Fprintf(outputWriter, "| metric | median | p90 | p95 | max |\n|---|---:|---:|---:|---:|\n")
	writeDistributionRow(outputWriter, "file lines", summarizeDistribution(lineValues))
	writeDistributionRow(outputWriter, "file words", summarizeDistribution(wordValues))

	sort.Slice(textMeasurements, func(left, right int) bool {
		if textMeasurements[left].Lines != textMeasurements[right].Lines {
			return textMeasurements[left].Lines > textMeasurements[right].Lines
		}
		return textMeasurements[left].Path < textMeasurements[right].Path
	})
	fmt.Fprintf(outputWriter, "\n## Largest Files — top %d by lines\n\n", topCount)
	fmt.Fprintf(outputWriter, "| path | lines | words |\n|---|---:|---:|\n")
	for index, measurement := range textMeasurements {
		if index >= topCount {
			break
		}
		fmt.Fprintf(outputWriter, "| %s | %d | %d |\n", measurement.Path, measurement.Lines, measurement.Words)
	}

	if thresholds.hasAny() {
		writeFileBandSection(outputWriter, textMeasurements, thresholds)
	}
}

// writeFileBandSection lists every file over a threshold as path/metric/value/
// band rows — the mechanical output the audit pastes.
func writeFileBandSection(outputWriter io.Writer, measurements []fileMeasurement, thresholds fileBandThresholds) {
	fmt.Fprintf(outputWriter, "\n## Band Flags — files\n\n")
	fmt.Fprintf(outputWriter, "| path | metric | value | band |\n|---|---|---:|---|\n")
	flaggedCount := 0
	for _, measurement := range measurements {
		if band := bandLabelForValue(measurement.Lines, thresholds.WatchLines, thresholds.FlagLines); band != "" {
			fmt.Fprintf(outputWriter, "| %s | lines | %d | %s |\n", measurement.Path, measurement.Lines, band)
			flaggedCount++
		}
		if band := bandLabelForValue(measurement.Words, thresholds.WatchWords, thresholds.FlagWords); band != "" {
			fmt.Fprintf(outputWriter, "| %s | words | %d | %s |\n", measurement.Path, measurement.Words, band)
			flaggedCount++
		}
	}
	if flaggedCount == 0 {
		fmt.Fprintf(outputWriter, "\n(no file exceeds a threshold)\n")
	}
}

// writeDistributionRow prints one metric's distribution table row.
func writeDistributionRow(outputWriter io.Writer, metricName string, summary distributionSummary) {
	fmt.Fprintf(outputWriter, "| %s | %d | %d | %d | %d |\n", metricName, summary.Median, summary.P90, summary.P95, summary.Max)
}

// folderCount is one directory and its direct tracked-file count.
type folderCount struct {
	Folder    string
	FileCount int
}

// countFilesPerFolder buckets measurements by their direct parent directory
// (repo root renders as "."). Direct children only — see the file-top comment.
func countFilesPerFolder(measurements []fileMeasurement) []folderCount {
	countsByFolder := map[string]int{}
	for _, measurement := range measurements {
		countsByFolder[path.Dir(measurement.Path)]++
	}
	var folderCounts []folderCount
	for folder, fileCount := range countsByFolder {
		folderCounts = append(folderCounts, folderCount{Folder: folder, FileCount: fileCount})
	}
	sort.Slice(folderCounts, func(left, right int) bool {
		if folderCounts[left].FileCount != folderCounts[right].FileCount {
			return folderCounts[left].FileCount > folderCounts[right].FileCount
		}
		return folderCounts[left].Folder < folderCounts[right].Folder
	})
	return folderCounts
}

// writeFoldersReport renders the files-per-folder distribution, top-N crowded
// folders, and — only when band flags were passed — the folders over each
// threshold.
func writeFoldersReport(outputWriter io.Writer, report inventoryReport, watchFiles int, flagFiles int, topCount int) {
	folderCounts := countFilesPerFolder(report.Files)
	var countValues []int
	for _, folder := range folderCounts {
		countValues = append(countValues, folder.FileCount)
	}
	fmt.Fprintf(outputWriter, "## Folders — tracked files per folder (direct children)\n\n")
	fmt.Fprintf(outputWriter, "| metric | median | p90 | p95 | max |\n|---|---:|---:|---:|---:|\n")
	writeDistributionRow(outputWriter, "files per folder", summarizeDistribution(countValues))

	fmt.Fprintf(outputWriter, "\n## Most Crowded Folders — top %d by direct file count\n\n", topCount)
	fmt.Fprintf(outputWriter, "| folder | files |\n|---|---:|\n")
	for index, folder := range folderCounts {
		if index >= topCount {
			break
		}
		fmt.Fprintf(outputWriter, "| %s | %d |\n", folder.Folder, folder.FileCount)
	}

	if watchFiles != bandThresholdUnset || flagFiles != bandThresholdUnset {
		fmt.Fprintf(outputWriter, "\n## Band Flags — folders\n\n")
		fmt.Fprintf(outputWriter, "| folder | metric | value | band |\n|---|---|---:|---|\n")
		flaggedCount := 0
		for _, folder := range folderCounts {
			if band := bandLabelForValue(folder.FileCount, watchFiles, flagFiles); band != "" {
				fmt.Fprintf(outputWriter, "| %s | files | %d | %s |\n", folder.Folder, folder.FileCount, band)
				flaggedCount++
			}
		}
		if flaggedCount == 0 {
			fmt.Fprintf(outputWriter, "\n(no folder exceeds a threshold)\n")
		}
	}
}
