package toolboxcommands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var toolboxNow = func() time.Time { return time.Now() }

func handleNote(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	rest, dryRun, commit, err := parseMutationFlags(arguments)
	if err != nil {
		return usageResult(CommandNote, err.Error())
	}
	text := normalizeNote(strings.Join(rest, " "))
	if text == "" {
		return usageResult(CommandNote, `Usage: do-work-note [--dry-run|--commit] "<text>"`)
	}
	relative := "do-work/notes.md"
	if confinementErr := validateNoLinkedAncestors(executionContext.RepositoryRoot, relative, false); confinementErr != nil {
		return usageResult(CommandNote, confinementErr.Error())
	}
	createdDirectories := absentTransactionDirectories(executionContext.RepositoryRoot, "do-work")
	line := fmt.Sprintf("- [%s] %s\n", toolboxNow().Format("2006-01-02"), text)
	result := runTransaction(CommandNote, executionContext.RepositoryRoot, []string{relative}, createdDirectories, dryRun, commit, "[do-work] Add note", func(recorder *gittransaction.MutationRecorder) error {
		if err := rootedMkdirAll(executionContext.RepositoryRoot, "do-work", 0o755); err != nil {
			return err
		}
		for _, directory := range createdDirectories {
			if err := recorder.RecordCreatedDirectory(directory); err != nil {
				return err
			}
		}
		current, readErr := rootedReadFile(executionContext.RepositoryRoot, relative)
		if os.IsNotExist(readErr) {
			if err := rootedPublishFile(executionContext.RepositoryRoot, relative, []byte(line), 0o644, false); err != nil {
				return err
			}
			return recorder.RecordCreated(relative)
		}
		if readErr != nil {
			return readErr
		}
		if err := rootedPublishFile(executionContext.RepositoryRoot, relative, append(current, []byte(line)...), 0o644, true); err != nil {
			return err
		}
		return recorder.RecordTouched(relative)
	})
	if result.Outcome == resultmodel.OutcomeSuccess {
		output := "Noted → do-work/notes.md\n  " + strings.TrimSuffix(line, "\n") + "\n"
		if dryRun {
			output = "Would note → do-work/notes.md\n  " + strings.TrimSuffix(line, "\n") + "\n"
		}
		result.ExactTextOutput = &output
	}
	return result
}

func normalizeNote(text string) string {
	text = strings.TrimSpace(text)
	if text == "add" {
		return ""
	}
	if strings.HasPrefix(text, "add ") {
		text = strings.TrimSpace(strings.TrimPrefix(text, "add "))
	}
	if len(text) >= 2 && ((text[0] == '\'' && text[len(text)-1] == '\'') || (text[0] == '"' && text[len(text)-1] == '"')) {
		text = strings.TrimSpace(text[1 : len(text)-1])
	}
	return text
}
