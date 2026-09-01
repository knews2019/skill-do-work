package hookcommands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestMemorySessionStartIsSilentWithoutStore(t *testing.T) {
	result := handleMemorySessionStart(commandruntime.ExecutionContext{RepositoryRoot: t.TempDir()}, nil)
	if result.ProtocolOutput == nil || *result.ProtocolOutput != "" {
		t.Fatalf("output=%#v", result.ProtocolOutput)
	}
}

func TestMemorySessionStartFiltersQuotedAndLegacyCaptures(t *testing.T) {
	repository := t.TempDir()
	memoryRoot := filepath.Join(repository, "memory")
	if err := os.MkdirAll(filepath.Join(memoryRoot, "logs"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(memoryRoot, "working-memory.md"), []byte("standing\n"), 0o600)
	logText := "## 09:00 UTC note\nkeep one\n\n## 10:00 UTC session capture abcdef01\n\n" + captureBodySentinel + "\n> raw secret\n> ## 12:00 UTC note\n\n## 11:00 UTC note\nkeep two\n"
	_ = os.WriteFile(filepath.Join(memoryRoot, "logs", "2026-09-01.md"), []byte(logText), 0o600)
	originalClock := hookClock
	hookClock = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { hookClock = originalClock })
	result := handleMemorySessionStart(commandruntime.ExecutionContext{RepositoryRoot: repository}, nil)
	if result.ProtocolOutput == nil || strings.Contains(*result.ProtocolOutput, "raw secret") || !strings.Contains(*result.ProtocolOutput, "keep one") || !strings.Contains(*result.ProtocolOutput, "keep two") {
		t.Fatalf("filtered output:\n%s", valueOrEmpty(result.ProtocolOutput))
	}
	ledger, err := os.ReadFile(filepath.Join(memoryRoot, "usage-ledger.jsonl"))
	if err != nil || !strings.Contains(string(ledger), `"source":"hooks/memory-session-start.sh"`) {
		t.Fatalf("ledger=%q err=%v", ledger, err)
	}
	want := "<background-memory>\n" + memoryInstruction + "\n\nstanding\n\n" +
		"## Today's log (2026-09-01) — curated entries only; raw session captures load via `do-work-knowledge memory recall`\n" +
		"## 09:00 UTC note\nkeep one\n\n## 11:00 UTC note\nkeep two\n</background-memory>\n"
	if *result.ProtocolOutput != want {
		t.Fatalf("exact output=%q, want %q", *result.ProtocolOutput, want)
	}
}

func TestCuratedLogSuppressesLegacyCaptureToEOF(t *testing.T) {
	contents := "## 09:00 UTC note\nkeep\n## 10:00 UTC session capture old\nraw\n## 11:00 UTC note\nhide\n"
	filtered := curatedLog(contents)
	if filtered != "## 09:00 UTC note\nkeep" {
		t.Fatalf("filtered=%q", filtered)
	}
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
