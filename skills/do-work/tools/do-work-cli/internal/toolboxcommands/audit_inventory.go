package toolboxcommands

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type auditInventory struct {
	files              []auditFile
	binary, unreadable int
}
type auditAggregateRow struct {
	extension           string
	files, lines, words int
}
type auditFolderRow struct {
	folder string
	files  int
}

func computeAuditInventory(root string, excludes []string) (auditInventory, error) {
	paths, err := auditTrackedFiles(root)
	if err != nil {
		return auditInventory{}, err
	}
	report := auditInventory{}
	for _, candidate := range paths {
		if auditExcluded(candidate, excludes) {
			continue
		}
		m, measureErr := measureAuditFile(root, candidate)
		if measureErr != nil {
			report.unreadable++
			continue
		}
		if m.binary {
			report.binary++
		}
		report.files = append(report.files, m)
	}
	return report, nil
}

func auditAggregates(files []auditFile) []auditAggregateRow {
	byExtension := map[string]*auditAggregateRow{}
	for _, file := range files {
		ext := path.Ext(file.path)
		if ext == "" {
			ext = "(none)"
		}
		row := byExtension[ext]
		if row == nil {
			row = &auditAggregateRow{extension: ext}
			byExtension[ext] = row
		}
		row.files++
		row.lines += file.lines
		row.words += file.words
	}
	rows := make([]auditAggregateRow, 0, len(byExtension))
	for _, row := range byExtension {
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].lines != rows[j].lines {
			return rows[i].lines > rows[j].lines
		}
		return rows[i].extension < rows[j].extension
	})
	return rows
}

func auditFolders(files []auditFile) []auditFolderRow {
	counts := map[string]int{}
	for _, file := range files {
		counts[path.Dir(file.path)]++
	}
	rows := make([]auditFolderRow, 0, len(counts))
	for folder, count := range counts {
		rows = append(rows, auditFolderRow{folder, count})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].files != rows[j].files {
			return rows[i].files > rows[j].files
		}
		return rows[i].folder < rows[j].folder
	})
	return rows
}

func renderAuditInventory(o auditOptions, report auditInventory) (*resultmodel.AuditMetricsResult, string) {
	payload := &resultmodel.AuditMetricsResult{Kind: o.kind, BinaryCount: report.binary, UnreadableCount: report.unreadable}
	var output strings.Builder
	if o.kind == "folders" {
		return renderAuditFolders(o, report, payload, &output)
	}
	fmt.Fprint(&output, "## Inventory — tracked files by extension\n\n| extension | files | lines | words |\n|---|---:|---:|---:|\n")
	lineTotal, wordTotal := 0, 0
	for _, row := range auditAggregates(report.files) {
		fmt.Fprintf(&output, "| %s | %d | %d | %d |\n", row.extension, row.files, row.lines, row.words)
		lineTotal += row.lines
		wordTotal += row.words
		payload.Aggregates = append(payload.Aggregates, resultmodel.AuditAggregate{Extension: row.extension, Files: row.files, Lines: row.lines, Words: row.words})
	}
	fmt.Fprintf(&output, "| **total** | %d | %d | %d |\n", len(report.files), lineTotal, wordTotal)
	if report.binary > 0 {
		fmt.Fprintf(&output, "\n%d binary file(s) excluded from line/word counts.\n", report.binary)
	}
	if report.unreadable > 0 {
		fmt.Fprintf(&output, "\nWARNING: %d tracked file(s) unreadable and skipped.\n", report.unreadable)
	}
	textFiles := []auditFile{}
	lines, words := []int{}, []int{}
	for _, file := range report.files {
		payload.Files = append(payload.Files, resultmodel.AuditFileMeasurement{Path: file.path, Lines: file.lines, Words: file.words, Binary: file.binary})
		if !file.binary {
			textFiles = append(textFiles, file)
			lines = append(lines, file.lines)
			words = append(words, file.words)
		}
	}
	lineSummary, wordSummary := auditSummary(lines), auditSummary(words)
	payload.Distributions = []resultmodel.AuditDistribution{{Metric: "file lines", Median: lineSummary.median, P90: lineSummary.p90, P95: lineSummary.p95, Max: lineSummary.max}, {Metric: "file words", Median: wordSummary.median, P90: wordSummary.p90, P95: wordSummary.p95, Max: wordSummary.max}}
	fmt.Fprintf(&output, "\n## File-Size Distributions (non-binary tracked files)\n\n| metric | median | p90 | p95 | max |\n|---|---:|---:|---:|---:|\n| file lines | %d | %d | %d | %d |\n| file words | %d | %d | %d | %d |\n", lineSummary.median, lineSummary.p90, lineSummary.p95, lineSummary.max, wordSummary.median, wordSummary.p90, wordSummary.p95, wordSummary.max)
	sort.Slice(textFiles, func(i, j int) bool {
		if textFiles[i].lines != textFiles[j].lines {
			return textFiles[i].lines > textFiles[j].lines
		}
		return textFiles[i].path < textFiles[j].path
	})
	fmt.Fprintf(&output, "\n## Largest Files — top %d by lines\n\n| path | lines | words |\n|---|---:|---:|\n", o.top)
	for i, file := range textFiles {
		if i >= o.top {
			break
		}
		fmt.Fprintf(&output, "| %s | %d | %d |\n", file.path, file.lines, file.words)
	}
	if o.watchLines != auditThresholdUnset || o.flagLines != auditThresholdUnset || o.watchWords != auditThresholdUnset || o.flagWords != auditThresholdUnset {
		fmt.Fprint(&output, "\n## Band Flags — files\n\n| path | metric | value | band |\n|---|---|---:|---|\n")
		count := 0
		for _, file := range textFiles {
			if band := auditBand(file.lines, o.watchLines, o.flagLines); band != "" {
				fmt.Fprintf(&output, "| %s | lines | %d | %s |\n", file.path, file.lines, band)
				payload.Bands = append(payload.Bands, resultmodel.AuditBand{Path: file.path, Metric: "lines", Value: file.lines, Band: band})
				count++
			}
			if band := auditBand(file.words, o.watchWords, o.flagWords); band != "" {
				fmt.Fprintf(&output, "| %s | words | %d | %s |\n", file.path, file.words, band)
				payload.Bands = append(payload.Bands, resultmodel.AuditBand{Path: file.path, Metric: "words", Value: file.words, Band: band})
				count++
			}
		}
		if count == 0 {
			fmt.Fprint(&output, "\n(no file exceeds a threshold)\n")
		}
	}
	return payload, output.String()
}

func renderAuditFolders(o auditOptions, report auditInventory, payload *resultmodel.AuditMetricsResult, output *strings.Builder) (*resultmodel.AuditMetricsResult, string) {
	rows := auditFolders(report.files)
	values := []int{}
	for _, row := range rows {
		values = append(values, row.files)
		payload.Folders = append(payload.Folders, resultmodel.AuditFolderMeasurement{Folder: row.folder, Files: row.files})
	}
	summary := auditSummary(values)
	payload.Distributions = []resultmodel.AuditDistribution{{Metric: "files per folder", Median: summary.median, P90: summary.p90, P95: summary.p95, Max: summary.max}}
	fmt.Fprintf(output, "## Folders — tracked files per folder (direct children)\n\n| metric | median | p90 | p95 | max |\n|---|---:|---:|---:|---:|\n| files per folder | %d | %d | %d | %d |\n\n## Most Crowded Folders — top %d by direct file count\n\n| folder | files |\n|---|---:|\n", summary.median, summary.p90, summary.p95, summary.max, o.top)
	for i, row := range rows {
		if i >= o.top {
			break
		}
		fmt.Fprintf(output, "| %s | %d |\n", row.folder, row.files)
	}
	if o.watchFiles != auditThresholdUnset || o.flagFiles != auditThresholdUnset {
		fmt.Fprint(output, "\n## Band Flags — folders\n\n| folder | metric | value | band |\n|---|---|---:|---|\n")
		count := 0
		for _, row := range rows {
			if band := auditBand(row.files, o.watchFiles, o.flagFiles); band != "" {
				fmt.Fprintf(output, "| %s | files | %d | %s |\n", row.folder, row.files, band)
				payload.Bands = append(payload.Bands, resultmodel.AuditBand{Path: row.folder, Metric: "files", Value: row.files, Band: band})
				count++
			}
		}
		if count == 0 {
			fmt.Fprint(output, "\n(no folder exceeds a threshold)\n")
		}
	}
	return payload, output.String()
}
