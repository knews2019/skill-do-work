package hookcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestMemoryStopCaptureSelectsRealTextRedactsAndDeduplicates(t *testing.T) {
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
	} else if info.Mode().Perm() != 0o600 {
		t.Fatalf("capture mode=%v, want 0600", info.Mode().Perm())
	}
	_ = handleMemoryStopCapture(context, nil, strings.NewReader(input))
	second, _ := os.ReadFile(logPath)
	if string(second) != string(first) {
		t.Fatal("duplicate stop event appended a second capture")
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
