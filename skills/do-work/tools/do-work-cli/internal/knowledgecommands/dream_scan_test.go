package knowledgecommands

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestDreamScanProducesExactlySevenDeterministicFindingClasses(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "memory/MEMORY.md", "# Index\n- [[dangling]]\n- [[alpha]]\n")
	writeFixture(t, root, "memory/wiki/alpha.md", "---\nname: Alpha Decision\nupdated: 2025-01-01\n---\n# Alpha Decision\nToday see [[missing]].\n")
	writeFixture(t, root, "memory/wiki/alpha-decisions.md", "---\nname: Alpha Decisions\nupdated: 2026-09-01\n---\n# Alpha Decisions\nNo links.\n")
	writeFixture(t, root, "memory/wiki/unindexed.md", "---\nname: Unindexed\nupdated: 2026-09-01\n---\n# Unindexed\nNo links.\n")
	before := treeDigest(t, filepath.Join(root, "memory"))
	previousNow := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previousNow })
	result := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "memory"})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("outcome = %s: %+v", result.Outcome, result.Findings)
	}
	codes := map[string]bool{}
	for _, finding := range result.Findings {
		codes[finding.Code] = true
	}
	for _, code := range dreamFindingCodes {
		if !codes[code] {
			t.Errorf("missing finding class %s", code)
		}
	}
	if len(codes) != 7 {
		t.Fatalf("finding classes = %v, want exactly seven", codes)
	}
	if before != treeDigest(t, filepath.Join(root, "memory")) {
		t.Fatal("dream-scan changed bytes")
	}
}

func TestDreamScanOrderingIsStable(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "wiki/index.md", "# Index\n")
	writeFixture(t, root, "wiki/z.md", "# Z\nToday.\n")
	writeFixture(t, root, "wiki/a.md", "# A\nToday.\n")
	first := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "wiki"})
	second := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "wiki"})
	if findingsSignature(first.Findings) != findingsSignature(second.Findings) {
		t.Fatal("scan order changed")
	}
}

func TestDreamScanParsesBOMAndCRLFAndRefusesMalformedFrontmatter(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "memory/MEMORY.md", "# Index\n- [[alpha]]\n")
	writeFixture(t, root, "memory/alpha.md", "\ufeff---\r\nname: Alpha\r\nupdated: 2025-01-01\r\n---\r\n# Alpha\r\nToday.\r\n")
	result := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "memory"})
	if result.Outcome != resultmodel.OutcomeFindings {
		t.Fatalf("BOM/CRLF outcome = %s", result.Outcome)
	}
	codes := map[string]bool{}
	for _, finding := range result.Findings {
		codes[finding.Code] = true
	}
	if !codes["DREAM-STALE-PAGE"] || !codes["DREAM-RELATIVE-DATE"] {
		t.Fatalf("BOM/CRLF findings = %v", codes)
	}
	writeFixture(t, root, "memory/alpha.md", "---\nname: Alpha\n# no closing fence\n")
	malformed := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "memory"})
	if malformed.Outcome != resultmodel.OutcomeFailure || len(malformed.Findings) == 0 {
		t.Fatalf("malformed result = %+v", malformed)
	}
}
