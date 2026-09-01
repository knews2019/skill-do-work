package hookcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestMemoryStopCaptureSelectsRealTextRedactsAndDeduplicates(t *testing.T) {
	originalUmask := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(originalUmask) })
	repository := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repository, "memory", "logs"), 0o755); err != nil {
		t.Fatal(err)
	}
	transcript := filepath.Join(repository, "transcript.jsonl")
	contents := `{"type":"user","message":{"content":"keep ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ"}}` + "\n" +
		`{"type":"assistant","message":{"content":[{"type":"text","text":"answer"},{"type":"tool_use"}]}}` + "\n" +
		`{"type":"user","message":{"content":[{"type":"tool_result"}]}}` + "\n"
	if err := os.WriteFile(transcript, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	originalClock := hookClock
	hookClock = func() time.Time { return time.Date(2026, 9, 1, 12, 34, 56, 0, time.UTC) }
	t.Cleanup(func() { hookClock = originalClock })
	input := fmt.Sprintf(`{"transcript_path":%q}`, transcript)
	context := commandruntime.ExecutionContext{RepositoryRoot: repository}
	result := handleMemoryStopCapture(context, nil, strings.NewReader(input))
	if result.ProtocolOutput == nil || *result.ProtocolOutput != "" {
		t.Fatalf("result=%+v", result)
	}
	logPath := filepath.Join(repository, "memory", "logs", "2026-09-01.md")
	first, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "\n## 12:34 UTC session capture 1e06ae68\n\n" + captureBodySentinel +
		"\n> Session capture — final exchange between the user and the agent:\n>\n" +
		"> User: keep [REDACTED]\n> \n> Agent: answer\n"
	if string(first) != want {
		t.Fatalf("capture=%q, want %q", first, want)
	}
	if info, statError := os.Stat(logPath); statError != nil {
		t.Fatalf("stat capture: %v", statError)
	} else if info.Mode().Perm() != 0o644 {
		t.Fatalf("capture mode=%v, want legacy umask-derived 0644", info.Mode().Perm())
	}
	if len(result.Changes) != 2 || result.Changes[0].Path != "memory/logs/2026-09-01.md" || result.Changes[1].Path != "memory/usage-ledger.jsonl" {
		t.Fatalf("typed result lost capture effects: %+v", result.Changes)
	}
	_ = handleMemoryStopCapture(context, nil, strings.NewReader(input))
	second, _ := os.ReadFile(logPath)
	if string(second) != string(first) {
		t.Fatal("duplicate stop event appended a second capture")
	}
}

// These cases are preserved from the retained jq/slurp behavior. The table is a
// differential characterization oracle: changing the Go parser must keep every
// status/effect decision below, including jq-accepted blank separators.
func TestRetainedStopDifferentialTranscriptMatrix(t *testing.T) {
	tests := []struct {
		name       string
		transcript []byte
		wantUser   string
		wantAgent  string
		wantError  bool
	}{
		{"string content", []byte("{\"type\":\"user\",\"message\":{\"content\":\"ask\"}}\n{\"type\":\"assistant\",\"message\":{\"content\":\"answer\"}}\n"), "ask", "answer", false},
		{"blank separator", []byte("{\"type\":\"user\",\"message\":{\"content\":\"ask\"}}\n\n{\"type\":\"assistant\",\"message\":{\"content\":\"answer\"}}\n"), "ask", "answer", false},
		{"tool blocks and meta skipped", []byte("{\"type\":\"user\",\"isMeta\":true,\"message\":{\"content\":\"meta\"}}\n{\"type\":\"user\",\"message\":{\"content\":[{\"type\":\"text\",\"text\":\"real\"}]} }\n{\"type\":\"assistant\",\"message\":{\"content\":[{\"type\":\"tool_use\"}]}}\n"), "real", "", false},
		{"malformed JSON", []byte("{\"type\":\"user\"}\nnot-json\n"), "", "", true},
		{"invalid UTF-8 JSON", append([]byte("{\"type\":\"user\",\"message\":{\"content\":\""), []byte{0xff, '"', '}', '}', '\n'}...), "", "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "transcript.jsonl")
			if err := os.WriteFile(path, test.transcript, 0o600); err != nil {
				t.Fatal(err)
			}
			user, agent, err := finalTranscriptTexts(path)
			if (err != nil) != test.wantError || user != test.wantUser || agent != test.wantAgent {
				t.Fatalf("user=%q agent=%q err=%v", user, agent, err)
			}
		})
	}
}

func TestRetainedStopDifferentialUmaskMatrix(t *testing.T) {
	for _, test := range []struct {
		name  string
		umask int
		want  os.FileMode
	}{{"ordinary", 0o022, 0o644}, {"restrictive", 0o077, 0o600}} {
		t.Run(test.name, func(t *testing.T) {
			originalUmask := syscall.Umask(test.umask)
			defer syscall.Umask(originalUmask)
			repository := t.TempDir()
			_ = os.MkdirAll(filepath.Join(repository, "memory", "logs"), 0o755)
			transcript := filepath.Join(repository, "transcript.jsonl")
			_ = os.WriteFile(transcript, []byte("{\"type\":\"user\",\"message\":{\"content\":\"ask\"}}\n{\"type\":\"assistant\",\"message\":{\"content\":\"answer\"}}\n"), 0o600)
			input := fmt.Sprintf(`{"transcript_path":%q}`, transcript)
			_ = handleMemoryStopCapture(commandruntime.ExecutionContext{RepositoryRoot: repository}, nil, strings.NewReader(input))
			for _, path := range []string{filepath.Join(repository, "memory", "logs", hookClock().UTC().Format("2006-01-02")+".md"), filepath.Join(repository, "memory", "usage-ledger.jsonl")} {
				info, err := os.Stat(path)
				if err != nil {
					t.Fatalf("stat %s: %v", path, err)
				}
				if info.Mode().Perm() != test.want {
					t.Fatalf("%s mode=%v, want %v", path, info.Mode().Perm(), test.want)
				}
			}
		})
	}
}

func TestMemoryStopCaptureReportsAppendFailureOnlyInTypedResult(t *testing.T) {
	repository := t.TempDir()
	logs := filepath.Join(repository, "memory", "logs")
	_ = os.MkdirAll(logs, 0o755)
	transcript := filepath.Join(repository, "transcript.jsonl")
	_ = os.WriteFile(transcript, []byte("{\"type\":\"user\",\"message\":{\"content\":\"ask\"}}\n{\"type\":\"assistant\",\"message\":{\"content\":\"answer\"}}\n"), 0o600)
	today := hookClock().UTC().Format("2006-01-02") + ".md"
	_ = os.Mkdir(filepath.Join(logs, today), 0o755)
	input := fmt.Sprintf(`{"transcript_path":%q}`, transcript)
	result := handleMemoryStopCapture(commandruntime.ExecutionContext{RepositoryRoot: repository}, nil, strings.NewReader(input))
	if result.ProtocolOutput == nil || *result.ProtocolOutput != "" || len(result.Findings) != 1 || result.Findings[0].Code != "MEMORY-CAPTURE-APPEND-SKIPPED" {
		t.Fatalf("result=%+v", result)
	}
}

func TestMemoryStopCaptureConcurrentDistinctWritesRemainWhole(t *testing.T) {
	repository := t.TempDir()
	_ = os.MkdirAll(filepath.Join(repository, "memory", "logs"), 0o755)
	const captures = 8
	var wait sync.WaitGroup
	for index := 0; index < captures; index++ {
		transcript := filepath.Join(repository, fmt.Sprintf("transcript-%d.jsonl", index))
		_ = os.WriteFile(transcript, []byte(fmt.Sprintf("{\"type\":\"user\",\"message\":{\"content\":\"ask-%d\"}}\n{\"type\":\"assistant\",\"message\":{\"content\":\"answer-%d\"}}\n", index, index)), 0o600)
		wait.Add(1)
		go func(path string) {
			defer wait.Done()
			input := fmt.Sprintf(`{"transcript_path":%q}`, path)
			_ = handleMemoryStopCapture(commandruntime.ExecutionContext{RepositoryRoot: repository}, nil, strings.NewReader(input))
		}(transcript)
	}
	wait.Wait()
	logPath := filepath.Join(repository, "memory", "logs", hookClock().UTC().Format("2006-01-02")+".md")
	contents, err := os.ReadFile(logPath)
	if err != nil || strings.Count(string(contents), " UTC session capture ") != captures || strings.Count(string(contents), captureBodySentinel) != captures {
		t.Fatalf("concurrent capture structure drifted: err=%v headings=%d sentinels=%d", err, strings.Count(string(contents), " UTC session capture "), strings.Count(string(contents), captureBodySentinel))
	}
}

func TestCaptureRedactsEveryCredentialFamily(t *testing.T) {
	tests := []string{
		"ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"sk-ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"AKIAABCDEFGHIJKLMNOP",
		"xoxb-1234567890ABCDE",
		"eyJABCDEFGHIJK.ABCDE.FGHIJ",
		"Bearer ABCDEFGHIJKLMNOP",
		"api_key = abcdefghijklmnop",
	}
	for _, credential := range tests {
		redacted := redactCredentials("before " + credential + " after")
		if strings.Contains(redacted, credential) || !strings.Contains(redacted, "[REDACTED]") {
			t.Errorf("credential %q was not redacted: %q", credential, redacted)
		}
	}
}

func TestMemoryStopCaptureRawLoopGuardWinsOverMalformedJSON(t *testing.T) {
	repository := t.TempDir()
	_ = os.MkdirAll(filepath.Join(repository, "memory", "logs"), 0o755)
	_ = handleMemoryStopCapture(commandruntime.ExecutionContext{RepositoryRoot: repository}, nil, strings.NewReader(`not-json "stop_hook_active" : true`))
	entries, _ := os.ReadDir(filepath.Join(repository, "memory", "logs"))
	if len(entries) != 0 {
		t.Fatalf("loop guard wrote %d files", len(entries))
	}
}

func TestCaptureBudgetPreservesUTF8(t *testing.T) {
	user, assistant := budgetCaptureSides(strings.Repeat("界", 800), strings.Repeat("🙂", 800))
	composed := "User: " + user + "\n\nAgent: " + assistant
	if len(composed) > captureBudgetBytes || !utf8.ValidString(composed) || !strings.Contains(composed, "[truncated]") {
		t.Fatalf("bytes=%d valid=%t", len(composed), utf8.ValidString(composed))
	}
}

func TestRetainedStopDifferentialRedactionBeforeBudgetBoundary(t *testing.T) {
	credential := "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fullUser := strings.Repeat("x", captureTextBudget-17) + credential
	redactedUser := redactCredentials(fullUser)
	user, agent := budgetCaptureSides(redactedUser, "answer")
	composed := "User: " + user + "\n\nAgent: " + agent
	if strings.Contains(composed, "ghp_") || !strings.Contains(composed, "[REDACTED]") || len(composed) > captureBudgetBytes || !utf8.ValidString(composed) {
		t.Fatalf("redaction/budget boundary drifted: bytes=%d text=%q", len(composed), composed)
	}
}

func TestMemoryStopCaptureDropsPrivateKeyAndMalformedTranscript(t *testing.T) {
	for _, transcriptText := range []string{
		`{"type":"user","message":{"content":"-----BEGIN PRIVATE KEY-----"}}` + "\n",
		`{"type":"user","message":{"content":"valid"}}` + "\nnot json\n",
	} {
		repository := t.TempDir()
		_ = os.MkdirAll(filepath.Join(repository, "memory", "logs"), 0o755)
		path := filepath.Join(repository, "transcript.jsonl")
		_ = os.WriteFile(path, []byte(transcriptText), 0o600)
		input := fmt.Sprintf(`{"transcript_path":%q}`, path)
		_ = handleMemoryStopCapture(commandruntime.ExecutionContext{RepositoryRoot: repository}, nil, strings.NewReader(input))
		entries, _ := os.ReadDir(filepath.Join(repository, "memory", "logs"))
		if len(entries) != 0 {
			t.Fatalf("unsafe transcript wrote %d files", len(entries))
		}
	}
}
