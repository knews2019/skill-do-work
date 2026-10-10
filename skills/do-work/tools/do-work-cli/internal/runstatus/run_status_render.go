package runstatus

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// renderReport is the text (and --watch) form of the same rows the JSON result
// carries: a table, one remedy line per row that is not C1, the run-local
// files, and last the single most useful next command.
func renderReport(rows []statusRow, report *resultmodel.RunStatusResult, watch bool, now time.Time) string {
	var output strings.Builder
	nextLine := "next: nothing to run now; check again with `do-work status`"
	for _, row := range rows {
		if len(row.finding.NextArgv) > 0 {
			nextLine = "next: " + strings.Join(row.finding.NextArgv, " ")
			break
		}
	}
	if watch {
		classCounts := []string{}
		for _, class := range classPrecedence {
			count := 0
			for _, row := range rows {
				if row.record.Class == class {
					count++
				}
			}
			if count > 0 {
				classCounts = append(classCounts, fmt.Sprintf("%s %d", class, count))
			}
		}
		fmt.Fprintf(&output, "status %s: %d open rows (%s), %d run-local files\n", now.Format("15:04Z"), len(rows), strings.Join(classCounts, ", "), len(report.RunLocalFiles))
		for _, row := range rows {
			fmt.Fprintf(&output, "%s %s %s, activity %s, ETA %s\n", row.record.RequestID, row.record.Class, phaseOf(row.record), minutesText(row.record.MinutesSinceActivity), row.record.EtaText)
		}
		output.WriteString(nextLine + "\n")
		return output.String()
	}
	if len(rows) == 0 {
		output.WriteString("No claimed, blocked, waiting or earmarked REQs.\n")
	} else {
		table := tabwriter.NewWriter(&output, 0, 0, 2, ' ', 0)
		fmt.Fprintln(table, "REQ\tTITLE\tPHASE\tSINCE ACTIVITY\tETA\tCLASS")
		for _, row := range rows {
			fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\t%s %s\n", row.record.RequestID, shortTitle(row.record.Title), phaseOf(row.record),
				minutesText(row.record.MinutesSinceActivity), row.record.EtaText, row.record.Class, classNames[row.record.Class])
		}
		_ = table.Flush()
		output.WriteString("\n")
		for _, row := range rows {
			if row.record.Class != "C1" {
				fmt.Fprintf(&output, "%s %s %s: %s.\n", row.record.RequestID, row.record.Class, classNames[row.record.Class], row.remedy)
			}
		}
	}
	if len(report.RunLocalFiles) > 0 {
		output.WriteString(runLocalLabel + ":\n")
		for _, localFile := range report.RunLocalFiles {
			fmt.Fprintf(&output, "  %s  %d min old  %s\n", localFile.Path, localFile.AgeMinutes, fallbackDash(localFile.FirstLine))
		}
	}
	output.WriteString(nextLine + "\n")
	return output.String()
}

func phaseOf(record resultmodel.RunStatusRow) string {
	if record.LastActivityPhase != "" {
		return record.LastActivityPhase
	}
	if record.Column != "claimed" {
		return record.Column
	}
	return "-"
}

func minutesText(minutes *int) string {
	if minutes == nil {
		return "-"
	}
	return fmt.Sprintf("%d min", *minutes)
}

// shortTitle cuts a title to six words for the table.
func shortTitle(title string) string {
	words := strings.Fields(title)
	if len(words) <= 6 {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:6], " ") + "…"
}

func fallbackDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
