package knowledgecommands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const workingMemoryFixture = `---
updated: 2026-08-01
---

# Working Memory

## Active Threads

- Ship the command platform

## Notes

## Pending Decisions
`

func newMemoryRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.com")
	runGitFixture(t, root, "config", "user.name", "Fixture")
	writeInterviewFixture(t, root, "memory/working-memory.md", workingMemoryFixture)
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-m", "fixture")
	if err := os.WriteFile(filepath.Join(root, ".git", "info", "exclude"), []byte("/memory/logs/\n/memory/usage-ledger.jsonl\n/memory/.bootstrap-imported\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestMemoryRecallScoringStatusAndBroadRecall(t *testing.T) {
	root := newMemoryRepository(t)
	writeInterviewFixture(t, root, "memory/logs/2026-08-31.md", "## 10:00 UTC note\ncommand platform uses Go\n\n## 11:00 UTC session capture deadbeef\n<!-- do-work:capture-body quoted -->\n> raw command platform transcript\n")
	previous := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previous })

	recall := handleMemoryRecall(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"command platform"})
	hasScoredHit := false
	for _, finding := range recall.Findings {
		if strings.Contains(strings.Join(finding.Evidence, " "), "score=") {
			hasScoredHit = true
		}
	}
	if recall.Outcome != resultmodel.OutcomeSuccess || !hasScoredHit {
		t.Fatalf("recall = %#v", recall)
	}
	broad := handleMemoryRecall(commandruntime.ExecutionContext{RepositoryRoot: root}, nil)
	for _, finding := range broad.Findings {
		if strings.Contains(strings.Join(finding.Evidence, " "), "raw command platform transcript") {
			t.Fatalf("broad recall exposed capture: %#v", broad)
		}
	}
	status := handleMemoryStatus(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--memory-root", "./memory"})
	if status.Outcome != resultmodel.OutcomeSuccess || !strings.Contains(status.Findings[0].Evidence[0], "cap=2500") {
		t.Fatalf("status = %#v", status)
	}
	if !strings.Contains(strings.Join(status.Findings[0].NextArgv, "\x00"), "--memory-root\x00./memory") {
		t.Fatalf("status lost relative memory-root spelling: %#v", status.Findings[0].NextArgv)
	}
}

func TestMemoryRememberUsesExplicitSectionAndCommitsOnlyWorkingMemory(t *testing.T) {
	root := newMemoryRepository(t)
	result := handleMemoryRemember(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--section", "notes", "--commit", "Prefer exact typed results"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("remember = %#v", result)
	}
	working, _ := os.ReadFile(filepath.Join(root, "memory", "working-memory.md"))
	if !strings.Contains(string(working), "- Prefer exact typed results") {
		t.Fatalf("working memory = %q", working)
	}
	logs, _ := filepath.Glob(filepath.Join(root, "memory", "logs", "*.md"))
	if len(logs) != 1 {
		t.Fatalf("logs = %#v", logs)
	}
	commit := strings.TrimSpace(runGitFixture(t, root, "rev-parse", "HEAD"))
	if paths := strings.TrimSpace(runGitFixture(t, root, "show", "--pretty=", "--name-only", commit)); paths != "memory/working-memory.md" {
		t.Fatalf("commit paths = %q", paths)
	}
}

func TestMemoryForgetRequiresContentBoundConfirmationAndPreservesCaptureQuote(t *testing.T) {
	root := newMemoryRepository(t)
	writeInterviewFixture(t, root, "memory/logs/2026-08-31.md", "## 10:00 UTC session capture deadbeef\n<!-- do-work:capture-body quoted -->\n> Ship the command platform\n")
	discovery := handleMemoryForget(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"command platform"})
	if discovery.Outcome != resultmodel.OutcomeFindings || len(discovery.Findings) < 2 {
		t.Fatalf("discovery = %#v", discovery)
	}
	ids := []string{}
	for _, finding := range discovery.Findings {
		ids = append(ids, finding.AffectedIDs...)
	}
	arguments := []string{"--confirm"}
	for _, id := range ids {
		arguments = append(arguments, "--match", id)
	}
	arguments = append(arguments, "command platform")
	result := handleMemoryForget(commandruntime.ExecutionContext{RepositoryRoot: root}, arguments)
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("forget = %#v", result)
	}
	logBytes, _ := os.ReadFile(filepath.Join(root, "memory/logs/2026-08-31.md"))
	if !strings.Contains(string(logBytes), "> [forgotten — redacted by memory forget 2026-09-01]") || strings.Contains(string(logBytes), "> Ship the command platform") {
		t.Fatalf("log = %q", logBytes)
	}
}

func TestMemoryBootstrapConsumesApprovedManifestOnceAndRefusesCommit(t *testing.T) {
	root := newMemoryRepository(t)
	manifest := filepath.Join(root, "imports.json")
	writeInterviewFixture(t, root, "imports.json", `[{"date":"2026-08-30","time":"09:15","source":"session-a.jsonl","summary":"The user chose Go for deterministic commands."}]`)
	refused := handleMemoryBootstrap(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--manifest", manifest, "--confirm", "--commit"})
	if refused.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("bootstrap commit = %#v", refused)
	}
	result := handleMemoryBootstrap(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--manifest", manifest, "--confirm"})
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("bootstrap = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(root, "memory", ".bootstrap-imported")); err != nil {
		t.Fatal(err)
	}
	again := handleMemoryBootstrap(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--manifest", manifest, "--confirm"})
	if again.Outcome != resultmodel.OutcomeRefused {
		t.Fatalf("second bootstrap = %#v", again)
	}
}

func TestPrivatePublicationDoesNotFollowParentSwapOutsideRepository(t *testing.T) {
	root := newMemoryRepository(t)
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "memory", "logs"), 0o755); err != nil {
		t.Fatal(err)
	}
	rootedCreateTestHook = func(repositoryRoot, relative string) {
		if !strings.Contains(relative, "logs/") {
			return
		}
		if err := os.Rename(filepath.Join(repositoryRoot, "memory", "logs"), filepath.Join(repositoryRoot, "memory", "owned-logs")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(repositoryRoot, "memory", "logs")); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { rootedCreateTestHook = nil })

	result := handleMemoryRemember(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--section", "notes", "Never escape the repository"})
	if result.Outcome != resultmodel.OutcomeRolledBack && result.Outcome != resultmodel.OutcomeRisk {
		t.Fatalf("parent swap result = %#v", result)
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("rooted publication wrote outside repository: %#v", entries)
	}
}

func TestLexicalRecallMatchesRetainedScriptAtRecencyBoundaries(t *testing.T) {
	root := newMemoryRepository(t)
	for _, fixture := range []struct {
		date, text string
	}{
		{"2026-08-25", "alpha beta seven"},
		{"2026-08-24", "alpha beta eight"},
		{"2026-08-02", "alpha beta thirty"},
		{"2026-08-01", "alpha beta thirty-one"},
	} {
		writeInterviewFixture(t, root, "memory/logs/"+fixture.date+".md", "## 10:00 UTC note\n"+fixture.text+"\n")
	}
	previous := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previous })
	hits, err := lexicalMemoryRecall(filepath.Join(root, "memory"), "memory", "ALPHA!!! beta alpha")
	if err != nil {
		t.Fatal(err)
	}
	scriptPath, err := filepath.Abs("../../../../../do-work-knowledge/scripts/lexical-memory-recall.sh")
	if err != nil {
		t.Fatal(err)
	}
	output, err := exec.Command("bash", scriptPath, filepath.Join(root, "memory"), "ALPHA!!! beta alpha").CombinedOutput()
	if err != nil {
		t.Fatalf("retained recall script: %v\n%s", err, output)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) != len(hits) {
		t.Fatalf("script lines=%d Go hits=%d\n%s\n%#v", len(lines), len(hits), output, hits)
	}
	for index, hit := range hits {
		fields := strings.Split(lines[index], "\t")
		if len(fields) != 5 || fields[0] != fmt.Sprint(hit.Score) || fields[4] != hit.Content {
			t.Fatalf("differential row %d script=%q Go=%#v", index, lines[index], hit)
		}
	}
}
