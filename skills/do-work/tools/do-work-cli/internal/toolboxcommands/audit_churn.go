package toolboxcommands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type auditChurnReport struct {
	shallow     bool
	commits     int
	window      string
	touches     map[string]int
	unavailable []string
}
type auditChurnRow struct {
	path    string
	commits int
}
type auditHotspotRow struct {
	path                  string
	commits, lines, score int
}

func computeAuditChurn(root, window string, excludes []string) (auditChurnReport, error) {
	report := auditChurnReport{window: window, touches: map[string]int{}}
	shallow, err := auditGit(root, "rev-parse", "--is-shallow-repository")
	if err != nil {
		return report, err
	}
	report.shallow = strings.TrimSpace(shallow) == "true"
	logOutput, err := auditGit(root, "-c", "core.quotepath=off", "log", "-M", "-C", "--find-copies-harder", "--name-status", "--format=%H", "--since="+window)
	if err != nil {
		return report, err
	}
	aliases, copies := map[string]string{}, map[string]string{}
	for _, line := range strings.Split(logOutput, "\n") {
		if line == "" {
			continue
		}
		if auditHashLine(line) {
			report.commits++
			continue
		}
		fields := strings.Split(line, "\t")
		status := fields[0]
		switch {
		case len(status) > 0 && (status[0] == 'R' || status[0] == 'C') && len(fields) >= 3:
			current := auditResolveAlias(aliases, fields[2])
			report.touches[current]++
			if status[0] == 'R' {
				aliases[fields[1]] = current
			} else {
				if _, exists := copies[fields[1]]; !exists {
					copies[fields[1]] = current
				}
			}
		case status == "D" && len(fields) >= 2:
		case len(fields) >= 2:
			report.touches[auditResolveAlias(aliases, fields[1])]++
		}
	}
	tracked, err := auditTrackedFiles(root)
	if err != nil {
		return report, err
	}
	live, visible := map[string]bool{}, map[string]bool{}
	for _, p := range tracked {
		live[p] = true
		if !auditExcluded(p, excludes) {
			visible[p] = true
		}
	}
	for p, count := range report.touches {
		if live[p] {
			continue
		}
		if survivor, ok := auditSurvivingCopy(copies, live, p); ok {
			report.touches[survivor] += count
		}
	}
	for p := range report.touches {
		if !visible[p] {
			delete(report.touches, p)
		}
	}
	return report, nil
}

func auditHashLine(line string) bool {
	if len(line) != 40 && len(line) != 64 {
		return false
	}
	for _, r := range line {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}
func auditResolveAlias(aliases map[string]string, p string) string {
	for i := 0; i <= len(aliases); i++ {
		next, ok := aliases[p]
		if !ok {
			return p
		}
		p = next
	}
	return p
}
func auditSurvivingCopy(copies map[string]string, live map[string]bool, p string) (string, bool) {
	for i := 0; i <= len(copies); i++ {
		next, ok := copies[p]
		if !ok {
			return "", false
		}
		if live[next] {
			return next, true
		}
		p = next
	}
	return "", false
}
func sortedAuditChurn(touches map[string]int) []auditChurnRow {
	rows := make([]auditChurnRow, 0, len(touches))
	for p, c := range touches {
		rows = append(rows, auditChurnRow{p, c})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].commits != rows[j].commits {
			return rows[i].commits > rows[j].commits
		}
		return rows[i].path < rows[j].path
	})
	return rows
}

func renderAuditChurn(o auditOptions, report auditChurnReport) (*resultmodel.AuditMetricsResult, string) {
	payload := &resultmodel.AuditMetricsResult{Kind: o.kind, SinceWindow: report.window, ShallowClone: report.shallow, CommitCount: report.commits}
	var output strings.Builder
	if o.kind == "churn" {
		fmt.Fprintf(&output, "## Churn — commits touching each file (since %s)\n\n", report.window)
		writeAuditShallow(&output, report.shallow)
		fmt.Fprintf(&output, "Commits scanned: %d. Rename-normalized; deleted paths dropped.\n\n| path | commits |\n|---|---:|\n", report.commits)
		for i, row := range sortedAuditChurn(report.touches) {
			payload.Churn = append(payload.Churn, resultmodel.AuditChurnMeasurement{Path: row.path, Commits: row.commits})
			if i < o.top {
				fmt.Fprintf(&output, "| %s | %d |\n", row.path, row.commits)
			}
		}
		return payload, output.String()
	}
	rows := []auditHotspotRow{}
	unavailable := []string{}
	for p, c := range report.touches {
		m, err := measureAuditFile(o.repoRoot, p)
		if err != nil {
			unavailable = append(unavailable, p)
			continue
		}
		rows = append(rows, auditHotspotRow{p, c, m.lines, c * m.lines})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].score != rows[j].score {
			return rows[i].score > rows[j].score
		}
		return rows[i].path < rows[j].path
	})
	sort.Strings(unavailable)
	payload.UnavailablePaths = append(payload.UnavailablePaths, unavailable...)
	fmt.Fprintf(&output, "## Hotspots — churn × size (since %s)\n\n", report.window)
	writeAuditShallow(&output, report.shallow)
	if len(unavailable) > 0 {
		fmt.Fprintf(&output, "> WARNING: numeric hotspot ranking is incomplete — %d churn-bearing tracked path(s) could not be measured.\n\n", len(unavailable))
	}
	fmt.Fprintf(&output, "Commits scanned: %d. Score = commits × current lines.\n\n| path | commits | lines | score |\n|---|---:|---:|---:|\n", report.commits)
	for i, row := range rows {
		payload.Hotspots = append(payload.Hotspots, resultmodel.AuditHotspot{Path: row.path, Commits: row.commits, Lines: row.lines, Score: row.score})
		if i < o.top {
			fmt.Fprintf(&output, "| %s | %d | %d | %d |\n", row.path, row.commits, row.lines, row.score)
		}
	}
	if len(unavailable) > 0 {
		fmt.Fprint(&output, "\n## NOT-MEASURED — unavailable tracked paths\n\n| path | commits | current lines | score |\n|---|---:|---:|---:|\n")
		for _, p := range unavailable {
			fmt.Fprintf(&output, "| %s | %d | NOT-MEASURED | NOT-MEASURED |\n", p, report.touches[p])
		}
	}
	return payload, output.String()
}

func writeAuditShallow(output *strings.Builder, shallow bool) {
	if shallow {
		fmt.Fprint(output, "> WARNING: shallow clone detected — history is truncated at the shallow boundary; counts below UNDERCOUNT reality.\n\n")
	}
}
