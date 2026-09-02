package knowledgecommands

import (
	"encoding/json"
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
	// Pin the clock: the redaction stamp below is asserted as a literal date.
	previous := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previous })
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
	// The retained script runs in a separate process on the real clock, so the
	// in-process nowUTC stub cannot reach it. Both sides therefore share one
	// captured wall-clock day. Retry once if UTC midnight crosses while the
	// external comparison is running rather than comparing two different days.
	for attempt := 0; attempt < 2; attempt++ {
		runNow := nowUTC()
		runDay := runNow.Format("2006-01-02")
		root := newMemoryRepository(t)
		for _, fixture := range []struct {
			daysAgo int
			text    string
		}{
			{7, "alpha beta seven"},
			{8, "alpha beta eight"},
			{30, "alpha beta thirty"},
			{31, "alpha beta thirty-one"},
		} {
			date := runNow.AddDate(0, 0, -fixture.daysAgo).Format("2006-01-02")
			writeInterviewFixture(t, root, "memory/logs/"+date+".md", "## 10:00 UTC note\n"+fixture.text+"\n")
		}
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
		if nowUTC().Format("2006-01-02") != runDay {
			continue
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
		return
	}
	t.Fatal("UTC day changed during both retained-script comparison attempts")
}

func TestForgetAndBootstrapNeverFollowPrivateSymlinksOrDiscloseOutsideBytes(t *testing.T) {
	root := newMemoryRepository(t)
	outside := filepath.Join(t.TempDir(), "outside-secret.md")
	const canary = "outside-secret-canary-417"
	if err := os.WriteFile(outside, []byte(canary+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "memory", "logs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "memory", "logs", "2026-09-01.md")); err != nil {
		t.Fatal(err)
	}
	forget := handleMemoryForget(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"outside secret"})
	forgetJSON, _ := json.Marshal(forget)
	if forget.Outcome != resultmodel.OutcomeFailure || strings.Contains(string(forgetJSON), canary) {
		t.Fatalf("forget followed/disclosed outside target: %s", forgetJSON)
	}
	if err := os.Remove(filepath.Join(root, "memory", "logs", "2026-09-01.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "memory", ".bootstrap-imported")); err != nil {
		t.Fatal(err)
	}
	bootstrap := handleMemoryBootstrap(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--confirm", "--manifest", filepath.Join(root, "missing.json")})
	bootstrapJSON, _ := json.Marshal(bootstrap)
	if bootstrap.Outcome != resultmodel.OutcomeFailure || strings.Contains(string(bootstrapJSON), canary) {
		t.Fatalf("bootstrap followed/disclosed outside sentinel: %s", bootstrapJSON)
	}
}

func TestRememberForgetAndStatusNeverPersistOrRediscloseProtectedTextInLedger(t *testing.T) {
	root := newMemoryRepository(t)
	const canary = "protected-fact-canary-417"
	remember := handleMemoryRemember(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--section", "notes", "--commit", canary})
	if remember.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("remember = %#v", remember)
	}
	ledgerPath := filepath.Join(root, "memory", "usage-ledger.jsonl")
	ledger, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(ledger), canary) || strings.Contains(string(ledger), "protected fact") {
		t.Fatalf("remember ledger retained fact: %s", ledger)
	}
	if err := os.WriteFile(ledgerPath, append(ledger, []byte(`{"ts":"2026-09-01T10:00:00Z","engine":"memory","event":"recall","query":"protected-fact-canary-417","hits":1,"source":"fixture","note":"raw-secret"}`+"\n")...), 0o600); err != nil {
		t.Fatal(err)
	}
	status := handleMemoryStatus(commandruntime.ExecutionContext{RepositoryRoot: root}, nil)
	statusJSON, _ := json.Marshal(status)
	if strings.Contains(string(statusJSON), canary) || strings.Contains(string(statusJSON), "raw-secret") {
		t.Fatalf("status redisclosed protected ledger fields: %s", statusJSON)
	}
	discovery := handleMemoryForget(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{canary})
	if len(discovery.Findings) == 0 {
		t.Fatalf("forget discovery = %#v", discovery)
	}
	id := discovery.Findings[0].AffectedIDs[0]
	forgotten := handleMemoryForget(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--confirm", "--match", id, canary})
	if forgotten.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("forget = %#v", forgotten)
	}
	ledger, _ = os.ReadFile(ledgerPath)
	lines := strings.Split(strings.TrimSpace(string(ledger)), "\n")
	if strings.Contains(lines[len(lines)-1], canary) || strings.Contains(lines[len(lines)-1], "protected fact") {
		t.Fatalf("forget ledger retained query: %s", lines[len(lines)-1])
	}
}

func TestMemoryAuditReportsRequiredProbesAndOldCitationDoesNotMakeActive(t *testing.T) {
	root := newMemoryRepository(t)
	writeInterviewFixture(t, root, ".claude/settings.json", `{"hooks":{"SessionStart":"memory-session-start.sh","Stop":"memory-stop-capture.sh"}}`)
	writeInterviewFixture(t, root, "memory/logs/2026-08-31.md", "## 09:00 UTC session capture abcdef01\n> captured\n\n## 10:00 UTC note\ncurated\n")
	writeInterviewFixture(t, root, "memory/usage-ledger.jsonl", strings.Join([]string{
		`{"ts":"2026-08-31T09:00:00Z","engine":"memory","event":"recall","hits":1}`,
		`{"ts":"2026-08-31T10:00:00Z","engine":"memory","event":"write","hits":1}`,
		`{"ts":"2026-08-31T11:00:00Z","engine":"memory","event":"other-new-event","hits":0}`,
		`{"ts":"2026-07-01T11:00:00Z","engine":"memory","event":"hit_cited","hits":1}`,
	}, "\n")+"\n")
	previous := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previous })
	result := handleMemoryAudit(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--engine", "memory"})
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Findings) != 1 {
		t.Fatalf("audit = %#v", result)
	}
	evidence := strings.Join(result.Findings[0].Evidence, " ")
	for _, required := range []string{"classification=Idle", "working_present=true", "section_fill=", "hook_start=true", "hook_stop=true", "log_days=1", "captures=1", "notes=1", "retrievals_28d=1", "hit_cited_28d=0", "weeks="} {
		if !strings.Contains(evidence, required) {
			t.Fatalf("audit missing %q: %s", required, evidence)
		}
	}
}

func TestMemoryAuditUsesExactFourteenDayBoundary(t *testing.T) {
	root := t.TempDir()
	ledger := filepath.Join(root, "usage-ledger.jsonl")
	writeInterviewFixture(t, root, "usage-ledger.jsonl", strings.Join([]string{
		`{"ts":"2026-08-18T12:00:00Z","event":"write"}`,
		`{"ts":"2026-08-18T12:00:00Z","event":"recall"}`,
		`{"ts":"2026-08-18T12:00:00Z","event":"hit_cited"}`,
		`{"ts":"2026-08-18T11:59:59Z","event":"hit_cited"}`,
	}, "\n")+"\n")
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	stats := collectLedgerAudit(ledger, "recall", now)
	if stats.events14 != 3 || stats.hitCited14 != 1 {
		t.Fatalf("14-day stats = events:%d cited:%d", stats.events14, stats.hitCited14)
	}
	if classification := classifyAudit(stats.newest, stats.events14, stats.hitCited14, now, false); classification != "Active" {
		t.Fatalf("classification = %s", classification)
	}
	if classification := classifyAudit(now.AddDate(0, 0, -30), 0, 0, now, false); classification != "Idle" {
		t.Fatalf("exact 30-day classification = %s", classification)
	}
	if classification := classifyAudit(now.AddDate(0, 0, -30).Add(-time.Second), 0, 0, now, false); classification != "Stale" {
		t.Fatalf("older than 30-day classification = %s", classification)
	}
}

func TestBKBAuditEngineCoversDiscoveryHistoryInboundLedgerAndReadOnlyParity(t *testing.T) {
	root := t.TempDir()
	runGitFixture(t, root, "init")
	runGitFixture(t, root, "config", "user.email", "fixture@example.com")
	runGitFixture(t, root, "config", "user.name", "BKB Fixture")
	writeInterviewFixture(t, root, "kb/wiki/_master_index.md", "Total articles: 1 | Topic clusters: 0\n")
	writeInterviewFixture(t, root, "kb/wiki/log.md", "2026-08-31 | garden | maintained\n")
	writeInterviewFixture(t, root, "kb/wiki/concepts/alpha.md", "# Alpha\n")
	writeInterviewFixture(t, root, "kb/raw/inbox/source.md", "raw\n")
	writeInterviewFixture(t, root, "kb/usage-ledger.jsonl", strings.Join([]string{
		`{"ts":"2026-08-31T09:00:00Z","event":"query"}`,
		`{"ts":"2026-08-31T10:00:00Z","event":"query"}`,
		`{"ts":"2026-08-31T11:00:00Z","event":"hit_cited"}`,
		`not-json`,
	}, "\n")+"\n")
	writeInterviewFixture(t, root, "docs/reference.md", "See kb/wiki/concepts/alpha.md and [[alpha]].\n")
	if err := os.MkdirAll(filepath.Join(root, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGitFixture(t, root, "add", ".")
	runGitFixture(t, root, "commit", "-m", "BKB history")
	previous := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previous })
	before := treeDigest(t, filepath.Join(root, "kb"))
	if inbound := countBKBInboundReferences(root, filepath.Join(root, "kb")); inbound != 1 {
		t.Fatalf("inbound=%d want 1", inbound)
	}
	classification, evidence := auditBKBEngine(root, filepath.Join(root, "kb"), "kb")
	if classification != "Active" {
		t.Fatalf("classification=%s evidence=%s", classification, evidence)
	}
	for _, required := range []string{"wiki_pages=3", "raw_inbox=1", "git_commits=1", "git_authors=1", "authors=BKB Fixture", "inbound_refs=1", "retrievals_28d=2", "hit_cited_28d=1", "malformed_ledger=1", "fairness=git_and_log_cover_pre_ledger"} {
		if !strings.Contains(evidence, required) {
			t.Fatalf("audit evidence missing %q: %s", required, evidence)
		}
	}
	contexts := []struct {
		root string
		args []string
	}{
		{root, []string{"--engine", "bkb"}},
		{root, []string{"--engine", "bkb", "--kb", "./kb"}},
		{root, []string{"--engine", "bkb", "--kb", filepath.Join(root, "kb")}},
		{filepath.Join(root, "subdir"), []string{"--engine", "bkb", "--kb", "../kb"}},
	}
	for _, fixture := range contexts {
		result := handleMemoryAudit(commandruntime.ExecutionContext{RepositoryRoot: fixture.root}, fixture.args)
		if result.Outcome != resultmodel.OutcomeSuccess || !resultHasFindingCode(result, "MEMORY-AUDIT-BKB") {
			t.Fatalf("discovery root=%q args=%v result=%#v", fixture.root, fixture.args, result)
		}
		textOutput, _ := resultmodel.RenderResult(result, resultmodel.FormatText)
		jsonOutput, _ := resultmodel.RenderResult(result, resultmodel.FormatJSON)
		if !strings.Contains(string(textOutput), "classification=Active") || !strings.Contains(string(jsonOutput), "classification=Active") {
			t.Fatalf("text/json classification parity failed: text=%s json=%s", textOutput, jsonOutput)
		}
	}
	if after := treeDigest(t, filepath.Join(root, "kb")); after != before {
		t.Fatal("BKB audit changed bytes")
	}
	if stale := classifyAudit(time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), 0, 0, nowUTC(), false); stale != "Stale" {
		t.Fatalf("stale boundary=%s", stale)
	}
	if idle := classifyAudit(nowUTC().AddDate(0, 0, -1), 1, 0, nowUTC(), false); idle != "Idle" {
		t.Fatalf("idle boundary=%s", idle)
	}
}

func TestMemorySingletonOptionsRejectRepetition(t *testing.T) {
	for _, arguments := range [][]string{
		{"--memory-root", "memory", "--path", "memory"},
		{"--section", "notes", "--section", "notes", "fact"},
		{"--replace", "one", "--replace-id", "two", "fact"},
		{"--engine", "memory", "--engine", "both"},
		{"--kb", "one", "--kb", "two"},
		{"--manifest", "one", "--manifest", "two"},
		{"--text", "one", "--query", "two"},
		{"--dry-run", "--dry-run", "fact"},
		{"--commit", "--commit", "fact"},
		{"--confirm", "--confirm", "fact"},
	} {
		if _, err := parseMemoryOptions(arguments, true); err == nil {
			t.Fatalf("repeated singleton accepted: %#v", arguments)
		}
	}
}

func TestMemoryRememberDryRunPlansWithoutChangingAnyBytes(t *testing.T) {
	root := newMemoryRepository(t)
	before := treeDigest(t, root)
	result := handleMemoryRemember(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--section", "notes", "--dry-run", "planned fact"})
	after := treeDigest(t, root)
	if result.Outcome != resultmodel.OutcomeSuccess || len(result.Changes) != 2 || before != after {
		t.Fatalf("dry run=%#v tree changed=%v", result, before != after)
	}
}
