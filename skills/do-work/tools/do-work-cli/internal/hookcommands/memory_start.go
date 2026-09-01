package hookcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const captureBodySentinel = "<!-- do-work:capture-body quoted -->"
const memoryInstruction = "Frozen memory snapshot (see .claude/skills/do-work-knowledge/actions/memory.md). Treat as silent background context: do not greet, recap, or mention it unless it becomes relevant. Writes made this session surface at the NEXT session start."

func handleMemorySessionStart(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) != 0 {
		return usageResult(CommandMemorySessionStart, "memory-session-start accepts no options")
	}
	memoryRoot := filepath.Join(executionContext.RepositoryRoot, "memory")
	workingBytes, err := os.ReadFile(filepath.Join(memoryRoot, "working-memory.md"))
	if err != nil {
		return protocolResult("")
	}
	var output strings.Builder
	output.WriteString("<background-memory>\n")
	output.WriteString(memoryInstruction)
	output.WriteString("\n\n")
	output.Write(workingBytes)

	now := hookClock().UTC()
	date := now.Format("2006-01-02")
	if logBytes, readError := os.ReadFile(filepath.Join(memoryRoot, "logs", date+".md")); readError == nil {
		curated := curatedLog(string(logBytes))
		if strings.TrimSpace(curated) != "" {
			output.WriteString("\n## Today's log (")
			output.WriteString(date)
			output.WriteString(") — curated entries only; raw session captures load via `do-work-knowledge memory recall`\n")
			output.WriteString(curated)
			output.WriteByte('\n')
		}
	}
	output.WriteString("</background-memory>\n")
	appendLedger(filepath.Join(memoryRoot, "usage-ledger.jsonl"), fmt.Sprintf(
		`{"ts":"%s","engine":"memory","event":"inject","query":"","hits":0,"source":"hooks/memory-session-start.sh","note":""}`+"\n", now.Format("2006-01-02T15:04:05Z")))
	return protocolResult(output.String())
}

func curatedLog(contents string) string {
	lines := strings.Split(contents, "\n")
	kept := make([]string, 0, len(lines))
	inCapture := false
	captureFormat := ""
	for _, line := range lines {
		isHeading := isDailyLogHeading(line)
		isBoundary := isHeading && (!inCapture || captureFormat == "quoted")
		if isBoundary {
			inCapture = strings.Contains(line, "session capture")
			captureFormat = ""
		} else if inCapture && captureFormat == "" && strings.TrimSpace(line) != "" {
			if line == captureBodySentinel {
				captureFormat = "quoted"
			} else {
				captureFormat = "legacy"
			}
		}
		if !inCapture {
			kept = append(kept, line)
		}
	}
	return strings.TrimRight(strings.Join(kept, "\n"), "\n")
}

func isDailyLogHeading(line string) bool {
	if len(line) < len("## 00:00 UTC ") || !strings.HasPrefix(line, "## ") {
		return false
	}
	return line[3] >= '0' && line[3] <= '9' && line[4] >= '0' && line[4] <= '9' && line[5] == ':' &&
		line[6] >= '0' && line[6] <= '9' && line[7] >= '0' && line[7] <= '9' && strings.HasPrefix(line[8:], " UTC ")
}

func appendLedger(path, line string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return
	}
	_, _ = file.Write([]byte(line))
	_ = file.Close()
}
