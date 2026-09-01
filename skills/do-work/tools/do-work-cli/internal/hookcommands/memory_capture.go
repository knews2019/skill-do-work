package hookcommands

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const captureBudgetBytes = 1500
const captureTextBudget = captureBudgetBytes - 15
const captureSideFloor = captureTextBudget / 2

var stopLoopGuard = regexp.MustCompile(`"stop_hook_active"[[:space:]]*:[[:space:]]*true`)
var credentialPatterns = []struct {
	pattern     *regexp.Regexp
	replacement string
}{
	{regexp.MustCompile(`(?:gh[pousr]|github_pat)_[A-Za-z0-9_]{16,}`), "[REDACTED]"},
	{regexp.MustCompile(`sk-[A-Za-z0-9_-]{16,}`), "[REDACTED]"},
	{regexp.MustCompile(`AKIA[0-9A-Z]{16}`), "[REDACTED]"},
	{regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]{10,}`), "[REDACTED]"},
	{regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{5,}\.[A-Za-z0-9_-]{5,}`), "[REDACTED]"},
	{regexp.MustCompile(`(?i)Bearer[[:space:]]+[A-Za-z0-9._~+/=-]{16,}`), "Bearer [REDACTED]"},
	{regexp.MustCompile(`(?i)((?:password|passwd|secret|token|api[_-]?key)["']?[[:space:]]*[:=][[:space:]]*)["']?[^[:space:]"']{6,}`), "${1}[REDACTED]"},
}

type hookInput struct {
	TranscriptPath string `json:"transcript_path"`
}

type transcriptEntry struct {
	Type    string `json:"type"`
	IsMeta  bool   `json:"isMeta"`
	Message struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
}

type transcriptBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func handleMemoryStopCapture(executionContext commandruntime.ExecutionContext, arguments []string, input io.Reader) resultmodel.CommandResult {
	if len(arguments) != 0 {
		return usageResult(CommandMemoryStopCapture, "memory-stop-capture accepts no options")
	}
	emptyResult := protocolResult("")
	inputBytes, err := io.ReadAll(input)
	if err != nil || stopLoopGuard.Match(inputBytes) {
		return emptyResult
	}
	var event hookInput
	if json.Unmarshal(inputBytes, &event) != nil || event.TranscriptPath == "" {
		return emptyResult
	}
	logsRoot := filepath.Join(executionContext.RepositoryRoot, "memory", "logs")
	if info, statError := os.Stat(logsRoot); statError != nil || !info.IsDir() {
		return emptyResult
	}
	userText, assistantText, parseError := finalTranscriptTexts(event.TranscriptPath)
	if parseError != nil {
		return emptyResult
	}
	if strings.Contains(userText+" "+assistantText, "PRIVATE KEY-----") {
		return emptyResult
	}
	userText = redactCredentials(userText)
	assistantText = redactCredentials(assistantText)
	userText, assistantText = budgetCaptureSides(userText, assistantText)
	captureText := "User: " + userText + "\n\nAgent: " + assistantText
	captureText = truncateUTF8(captureText, captureBudgetBytes, false)
	if strings.TrimSpace(captureText) == "" || captureText == "User: \n\nAgent: " {
		return emptyResult
	}
	hashBytes := sha256.Sum256([]byte(captureText))
	hash := hex.EncodeToString(hashBytes[:])[:8]
	now := hookClock().UTC()
	todayLog := filepath.Join(logsRoot, now.Format("2006-01-02")+".md")
	if existing, readError := os.ReadFile(todayLog); readError == nil && bytes.Contains(existing, []byte("session capture "+hash)) {
		return emptyResult
	}
	quoted := "> " + strings.ReplaceAll(captureText, "\n", "\n> ")
	section := fmt.Sprintf("\n## %s UTC session capture %s\n\n%s\n> Session capture — final exchange between the user and the agent:\n>\n%s\n",
		now.Format("15:04"), hash, captureBodySentinel, quoted)
	logFile, openError := os.OpenFile(todayLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if openError != nil {
		return emptyResult
	}
	_, writeError := logFile.Write([]byte(section))
	closeError := logFile.Close()
	if writeError != nil || closeError != nil {
		return emptyResult
	}
	appendLedger(filepath.Join(executionContext.RepositoryRoot, "memory", "usage-ledger.jsonl"), fmt.Sprintf(
		`{"ts":"%s","engine":"memory","event":"capture","query":"","hits":0,"source":"hooks/memory-stop-capture.sh","note":"%s"}`+"\n",
		now.Format("2006-01-02T15:04:05Z"), hash))
	return emptyResult
}

func finalTranscriptTexts(path string) (string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	userText, assistantText := "", ""
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var entry transcriptEntry
		if json.Unmarshal(scanner.Bytes(), &entry) != nil {
			return "", "", fmt.Errorf("malformed transcript")
		}
		if entry.IsMeta || (entry.Type != "user" && entry.Type != "assistant") {
			continue
		}
		text, contentError := messageText(entry.Message.Content)
		if contentError != nil {
			return "", "", contentError
		}
		if strings.TrimSpace(text) == "" {
			continue
		}
		if entry.Type == "user" {
			userText = text
		} else {
			assistantText = text
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	return userText, assistantText, nil
}

func messageText(content json.RawMessage) (string, error) {
	if len(content) == 0 || bytes.Equal(content, []byte("null")) {
		return "", nil
	}
	var text string
	if json.Unmarshal(content, &text) == nil {
		return text, nil
	}
	var blocks []transcriptBlock
	if err := json.Unmarshal(content, &blocks); err != nil {
		return "", err
	}
	texts := []string{}
	for _, block := range blocks {
		if block.Type == "text" {
			texts = append(texts, block.Text)
		}
	}
	return strings.Join(texts, " "), nil
}

func redactCredentials(text string) string {
	for _, redaction := range credentialPatterns {
		text = redaction.pattern.ReplaceAllString(text, redaction.replacement)
	}
	return text
}

func budgetCaptureSides(userText, assistantText string) (string, string) {
	if len(userText)+len(assistantText) <= captureTextBudget {
		return userText, assistantText
	}
	if len(userText) <= captureSideFloor {
		assistantText = truncateUTF8(assistantText, captureTextBudget-len(userText), true)
	} else if len(assistantText) <= captureSideFloor {
		userText = truncateUTF8(userText, captureTextBudget-len(assistantText), true)
	} else {
		userText = truncateUTF8(userText, captureSideFloor, true)
		assistantText = truncateUTF8(assistantText, captureSideFloor, true)
	}
	return userText, assistantText
}

func truncateUTF8(text string, budget int, mark bool) string {
	if len(text) <= budget {
		return text
	}
	suffix := ""
	if mark {
		suffix = " [truncated]"
		budget -= len(suffix)
	}
	if budget < 0 {
		budget = 0
	}
	cut := []byte(text)
	if len(cut) > budget {
		cut = cut[:budget]
	}
	for len(cut) > 0 && !utf8.Valid(cut) {
		cut = cut[:len(cut)-1]
	}
	return string(cut) + suffix
}
