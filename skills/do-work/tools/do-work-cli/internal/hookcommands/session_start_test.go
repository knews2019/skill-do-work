package hookcommands

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
)

func TestSessionStartRendersExactBanner(t *testing.T) {
	repository := t.TempDir()
	skillRoot := filepath.Join(t.TempDir(), "skill")
	if err := os.MkdirAll(filepath.Join(skillRoot, "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillRoot, "actions", "version.md"), []byte("**Current version**: 9.8.7\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repository, "do-work", "queue"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"REQ-001-one.md", "REQ-002-two.md"} {
		if err := os.WriteFile(filepath.Join(repository, "do-work", "queue", name), []byte("---\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	result := handleSessionStart(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--skill-root", skillRoot})
	want := "do-work v9.8.7 loaded. 2 pending REQ(s). Say 'do-work help' for commands.\n"
	if result.ProtocolOutput == nil || *result.ProtocolOutput != want {
		t.Fatalf("output=%#v, want %q", result.ProtocolOutput, want)
	}
}

func TestSessionStartCombinesReservationAndTimestampHousekeeping(t *testing.T) {
	repository := t.TempDir()
	skillRoot := filepath.Join(t.TempDir(), "skill")
	if err := os.MkdirAll(filepath.Join(skillRoot, "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(skillRoot, "actions", "version.md"), []byte("**Current version**: 9.8.7\n"), 0o600)
	queueRoot := filepath.Join(repository, "do-work", "queue")
	reservationRoot := filepath.Join(repository, "do-work", ".req-reservations")
	if err := os.MkdirAll(queueRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(reservationRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(queueRoot, "REQ-001-fixture.md")
	requestBytes := []byte("---\nid: REQ-001\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z\n---\nbody\n")
	_ = os.WriteFile(requestPath, requestBytes, 0o600)
	runHookGit(t, repository, "init", "-q")
	runHookGit(t, repository, "config", "user.email", "fixture@example.invalid")
	runHookGit(t, repository, "config", "user.name", "Fixture")
	runHookGit(t, repository, "add", "do-work/queue/REQ-001-fixture.md")
	runHookGit(t, repository, "commit", "-qm", "land request")
	_ = os.WriteFile(requestPath, append(requestBytes, []byte("dirty\n")...), 0o600)
	_ = os.WriteFile(filepath.Join(reservationRoot, "REQ-000001"), nil, 0o600)
	mtime := time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
	_ = os.Chtimes(requestPath, mtime, mtime)
	originalClock := hookClock
	hookClock = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { hookClock = originalClock })

	result := handleSessionStart(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--skill-root", skillRoot})
	want := "do-work v9.8.7 loaded. 1 pending REQ(s). Say 'do-work help' for commands.\n" +
		"do-work: removed 1 stale REQ reservation marker(s) from do-work/.req-reservations/ — stage and commit the deletion(s).\n" +
		"do-work: repaired do-work/queue/REQ-001-fixture.md created_at: 2093-01-01T00:00:00Z -> 2026-08-10T12:00:00Z (file mtime)\n" +
		"do-work: repaired 1 detectably wrong timestamp(s) — review and commit the correction(s) with the next housekeeping commit.\n"
	if result.ProtocolOutput == nil || *result.ProtocolOutput != want {
		t.Fatalf("output=%q, want %q", valueOrEmpty(result.ProtocolOutput), want)
	}
	if len(result.Changes) != 2 {
		t.Fatalf("typed JSON seam lost housekeeping changes: %+v", result.Changes)
	}
	if _, err := os.Stat(filepath.Join(reservationRoot, "REQ-000001")); !os.IsNotExist(err) {
		t.Fatalf("reservation survived: %v", err)
	}
	updated, _ := os.ReadFile(requestPath)
	if !strings.Contains(string(updated), "created_at: 2026-08-10T12:00:00Z") {
		t.Fatalf("timestamp was not repaired:\n%s", updated)
	}
}

func runHookGit(t *testing.T, repository string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repository}, arguments...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
}

func TestSessionStartPreservesMultipleVersionLineQuirk(t *testing.T) {
	repository := t.TempDir()
	if err := os.Mkdir(filepath.Join(repository, "do-work"), 0o755); err != nil {
		t.Fatal(err)
	}
	skillRoot := t.TempDir()
	if err := os.Mkdir(filepath.Join(skillRoot, "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(skillRoot, "actions", "version.md"), []byte("**Current version**: 1\n**Current version**: 2\n"), 0o600)
	result := handleSessionStart(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--skill-root=" + skillRoot})
	want := "do-work v1\n2 loaded. 0 pending REQ(s). Say 'do-work help' for commands.\n"
	if result.ProtocolOutput == nil || *result.ProtocolOutput != want {
		t.Fatalf("output=%#v, want %q", result.ProtocolOutput, want)
	}
}

func TestRetainedSessionStartDifferentialTimestampRefusalShape(t *testing.T) {
	repository := t.TempDir()
	skillRoot := filepath.Join(t.TempDir(), "skill")
	if err := os.MkdirAll(filepath.Join(skillRoot, "actions"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(skillRoot, "actions", "version.md"), []byte("**Current version**: 9.8.7\n"), 0o600)
	queueRoot := filepath.Join(repository, "do-work", "queue")
	if err := os.MkdirAll(queueRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(queueRoot, "REQ-009-refused.md"), []byte("---\nid: REQ-009\nstatus: pending\ncreated_at: 2093-01-01T00:00:00.123Z\n---\n"), 0o600)
	result := handleSessionStart(commandruntime.ExecutionContext{RepositoryRoot: repository}, []string{"--skill-root", skillRoot})
	wantOutput := "do-work v9.8.7 loaded. 1 pending REQ(s). Say 'do-work help' for commands.\n"
	if result.ProtocolOutput == nil || *result.ProtocolOutput != wantOutput || len(result.Changes) != 0 {
		t.Fatalf("refusal protocol/effects drifted: %+v", result)
	}
	codes := map[string]bool{}
	for _, finding := range result.Findings {
		codes[finding.Code] = true
	}
	if !codes["TIMESTAMP-REPAIR-REFUSED"] || !codes["TIMESTAMP-FUTURE"] {
		t.Fatalf("typed refusal findings missing: %+v", result.Findings)
	}
}
