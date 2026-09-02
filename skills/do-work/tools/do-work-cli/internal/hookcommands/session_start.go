package hookcommands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/doctor"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var hookClock = func() time.Time { return time.Now().UTC() }

func handleSessionStart(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	skillRoot := ""
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--skill-root" && index+1 < len(arguments) {
			index++
			skillRoot = arguments[index]
			continue
		}
		if strings.HasPrefix(argument, "--skill-root=") && strings.TrimPrefix(argument, "--skill-root=") != "" {
			skillRoot = strings.TrimPrefix(argument, "--skill-root=")
			continue
		}
		return usageResult(CommandSessionStart, "session-start requires --skill-root <installed-core-root>")
	}
	if skillRoot == "" {
		return usageResult(CommandSessionStart, "session-start requires --skill-root <installed-core-root>")
	}

	var output strings.Builder
	fmt.Fprintf(&output, "do-work v%s loaded. %d pending REQ(s). Say 'do-work help' for commands.\n",
		installedVersion(skillRoot), directQueueCount(executionContext.RepositoryRoot))
	cleanup := corehelpers.CleanupReservations(executionContext.RepositoryRoot)
	operations := []resultmodel.CommandResult{cleanup}
	if len(cleanup.Changes) > 0 {
		fmt.Fprintf(&output, "do-work: removed %d stale REQ reservation marker(s) from do-work/.req-reservations/ — stage and commit the deletion(s).\n", len(cleanup.Changes))
	}

	snapshot, discoveryError := repositorymodel.DiscoverRepository(executionContext.RepositoryRoot)
	if discoveryError == nil {
		finalizationTailFindings := doctor.FinalizationTailFindings(context.Background(), snapshot)
		if len(finalizationTailFindings) > 0 {
			operations = append(operations, resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: finalizationTailFindings})
			if len(finalizationTailFindings[0].AffectedIDs) > 0 {
				fmt.Fprintf(&output, "do-work: unfinished finalization for %s — 'do-work run' resumes it; 'do-work run-with-recovery' if this checkout is the only writer.\n", finalizationTailFindings[0].AffectedIDs[0])
			}
		}
		plans, planFindings := doctor.BuildTimestampPlanForScope(context.Background(), snapshot, hookClock(), doctor.TimestampScopeActive)
		operations = append(operations, resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: planFindings})
		repair := doctor.ApplyUncommittedTimestampPlans(snapshot, plans)
		operations = append(operations, repair)
		applied := map[string]int{}
		for _, change := range repair.Changes {
			applied[change.Path]++
		}
		repairCount := 0
		for _, plan := range plans {
			remaining := applied[plan.RelativePath]
			for _, change := range plan.Changes {
				if remaining == 0 {
					break
				}
				fmt.Fprintf(&output, "do-work: repaired %s %s: %s -> %s (%s)\n", plan.RelativePath, change.FieldName, change.OldValue, change.NewValue, change.Source)
				remaining--
				repairCount++
			}
		}
		for _, finding := range repair.Findings {
			if len(finding.AffectedPaths) > 0 && len(finding.Evidence) > 0 {
				fmt.Fprintf(&output, "do-work: FAILED to repair %s — %s\n", finding.AffectedPaths[0], finding.Evidence[0])
			}
		}
		if repairCount > 0 {
			fmt.Fprintf(&output, "do-work: repaired %d detectably wrong timestamp(s) — review and commit the correction(s) with the next housekeeping commit.\n", repairCount)
		}
	}
	return protocolResult(output.String(), operations...)
}

func installedVersion(skillRoot string) string {
	contents, err := os.ReadFile(filepath.Join(skillRoot, "actions", "version.md"))
	if err != nil {
		return "unknown"
	}
	versions := []string{}
	for _, line := range strings.Split(string(contents), "\n") {
		const prefix = "**Current version**:"
		if strings.HasPrefix(line, prefix) {
			versions = append(versions, strings.TrimLeftFunc(line[len(prefix):], unicode.IsSpace))
		}
	}
	if len(versions) == 0 {
		return "unknown"
	}
	return strings.Join(versions, "\n")
}

func directQueueCount(repositoryRoot string) int {
	queueRoot := filepath.Join(repositoryRoot, "do-work", "queue")
	info, err := os.Lstat(queueRoot)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return 0
	}
	entries, err := os.ReadDir(queueRoot)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if matched, _ := filepath.Match("REQ-*.md", entry.Name()); matched {
			count++
		}
	}
	return count
}
