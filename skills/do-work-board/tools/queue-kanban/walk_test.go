package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// seedConsumerInstallLayout builds the consumer-install shape that broke root
// resolution: the skill vendored at skills/do-work/ (its top level ships
// SKILL.md) with the vendored tool nested inside it, and — when withQueue is
// true — the real queue at <consumer-root>/do-work/. It returns the consumer
// root and the vendored tool directory (the walk-up start point).
func seedConsumerInstallLayout(t *testing.T, withQueue bool) (string, string) {
	t.Helper()
	consumerRoot := t.TempDir()

	writeFixture := func(relativePath string, content string) {
		absolutePath := filepath.Join(consumerRoot, relativePath)
		if mkdirError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirError != nil {
			t.Fatalf("mkdir %s: %v", relativePath, mkdirError)
		}
		if writeError := os.WriteFile(absolutePath, []byte(content), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", relativePath, writeError)
		}
	}

	writeFixture(filepath.Join("skills", "do-work", "SKILL.md"), "# do-work skill entry point\n")
	if withQueue {
		writeFixture(filepath.Join("do-work", "queue", "REQ-9101-consumer.md"),
			"---\nid: REQ-9101\ntitle: Consumer fixture\nstatus: pending\n---\n\n## What\n\nBody for REQ-9101.\n")
	}

	toolDirectory := filepath.Join(consumerRoot, "skills", "do-work", "tools", "queue-kanban")
	if mkdirError := os.MkdirAll(toolDirectory, 0o755); mkdirError != nil {
		t.Fatalf("mkdir tool directory: %v", mkdirError)
	}
	return consumerRoot, toolDirectory
}

// TestResolveRepoRootSkipsSkillInstall reproduces the consumer-install bug:
// walking up from the vendored tool directory used to match skills/ (because
// skills/do-work exists) and silently render an empty board. A directory merely
// NAMED do-work that is a skill install (SKILL.md at its top level) must be
// skipped so the walk reaches the directory holding the actual queue.
func TestResolveRepoRootSkipsSkillInstall(t *testing.T) {
	consumerRoot, toolDirectory := seedConsumerInstallLayout(t, true)

	resolvedRoot, resolveError := resolveRepoRoot(toolDirectory)
	if resolveError != nil {
		t.Fatalf("resolveRepoRoot: %v", resolveError)
	}
	if resolvedRoot != consumerRoot {
		t.Fatalf("resolveRepoRoot = %s, want the consumer root %s", resolvedRoot, consumerRoot)
	}
}

// TestResolveRepoRootErrorsWhenOnlySkillInstallExists asserts the no-queue
// failure mode is a loud error, not a silently-empty board built from the
// skill install itself.
func TestResolveRepoRootErrorsWhenOnlySkillInstallExists(t *testing.T) {
	_, toolDirectory := seedConsumerInstallLayout(t, false)

	resolvedRoot, resolveError := resolveRepoRoot(toolDirectory)
	if resolveError == nil {
		t.Fatalf("resolveRepoRoot = %s, want an error when no queue-holding ancestor exists", resolvedRoot)
	}
}

// TestEnumerateSkipsAssetsRequestFiles guards the REQ-075 duplicate-id warning:
// a deliverable copy named REQ-NNN-*.md under a UR's assets/ folder is an
// attachment, not a queue ticket. Its filename-derived id would otherwise collide
// with the real REQ-NNN. The assets/ subtree must be pruned from REQ discovery —
// neither a RequestFile nor a flagged stray — while the real ticket is still found.
func TestEnumerateSkipsAssetsRequestFiles(t *testing.T) {
	repoRoot := t.TempDir()
	writeFixture := func(relativePath, content string) {
		absolutePath := filepath.Join(repoRoot, relativePath)
		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", relativePath, err)
		}
		if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", relativePath, err)
		}
	}
	realTicketRelativePath := filepath.Join("do-work", "archive", "UR-043", "REQ-075-fix.md")
	writeFixture(realTicketRelativePath,
		"---\nid: REQ-075\ntitle: Real ticket\nstatus: completed\nuser_request: UR-043\n---\n\n## What\n\nBody.\n")
	writeFixture(filepath.Join("do-work", "archive", "UR-043", "assets", "REQ-075-agent-prompt.md"),
		"# Agent prompt\n\nNot a queue ticket.\n")

	discovered, walkError := enumerateDoWorkTree(repoRoot)
	if walkError != nil {
		t.Fatalf("enumerateDoWorkTree: %v", walkError)
	}
	assetsMarker := string(filepath.Separator) + "assets" + string(filepath.Separator)
	for _, rf := range discovered.RequestFiles {
		if strings.Contains(rf.AbsolutePath, assetsMarker) {
			t.Fatalf("assets/ file discovered as a REQ: %s", rf.AbsolutePath)
		}
	}
	for _, stray := range discovered.StrayRequestFiles {
		if strings.Contains(stray.AbsolutePath, assetsMarker) {
			t.Fatalf("assets/ file flagged as a stray REQ: %s", stray.AbsolutePath)
		}
	}
	foundReal := false
	for _, rf := range discovered.RequestFiles {
		if strings.HasSuffix(rf.AbsolutePath, realTicketRelativePath) {
			foundReal = true
		}
	}
	if !foundReal {
		t.Fatalf("real REQ-075 ticket not discovered among %d RequestFiles", len(discovered.RequestFiles))
	}
}
