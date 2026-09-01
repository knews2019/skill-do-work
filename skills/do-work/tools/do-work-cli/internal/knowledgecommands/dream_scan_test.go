package knowledgecommands

import (
	"os"
	"path/filepath"
	"strings"
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

func TestDreamScanFlatDepthDuplicateStemsAndDanglingAffectedPath(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "memory/MEMORY.md", "# Index\n- [[dangling]]\n- [[alpha]]\n")
	writeFixture(t, root, "memory/wiki/alpha.md", "# Alpha\n")
	writeFixture(t, root, "memory/wiki/nested/ignored.md", "Today [[missing]].\n")
	result := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "memory"})
	for _, finding := range result.Findings {
		if strings.Contains(strings.Join(finding.AffectedPaths, " "), "nested") {
			t.Fatalf("recursive page was scanned: %+v", finding)
		}
		if finding.Code == "DREAM-DANGLING-INDEX" && (len(finding.AffectedPaths) != 1 || finding.AffectedPaths[0] != "memory/MEMORY.md") {
			t.Fatalf("dangling affected path=%v", finding.AffectedPaths)
		}
	}
	writeFixture(t, root, "memory/wiki/alpha.md.md", "# Duplicate normalized stem\n")
	duplicate := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: root}, []string{"--path", "memory"})
	if duplicate.Outcome != resultmodel.OutcomeFailure || !strings.Contains(duplicate.Findings[0].Evidence[0], "duplicate page stem") {
		t.Fatalf("duplicate=%+v", duplicate)
	}
}

func TestDreamScanAbsoluteTargetBoundariesVocabularyAndPhysicalRoot(t *testing.T) {
	root := t.TempDir()
	physical := filepath.Join(root, "physical")
	if err := os.MkdirAll(physical, 0o755); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(root, "alias")
	if err := os.Symlink(physical, alias); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, physical, "memory/MEMORY.md", "# Index\n- [[ninety]]\n- [[ninety-one]]\n")
	writeFixture(t, physical, "memory/wiki/ninety.md", "---\nupdated: 2026-06-03\n---\n# Ninety\nYesterday, today, tomorrow, tonight, last week, next month, this evening, a few days ago, recently, just now, earlier today, the other day. [[ninety-one]]\n")
	writeFixture(t, physical, "memory/wiki/ninety-one.md", "---\nupdated: 2026-06-02\n---\n# Ninety One\n[[ninety]]\n")
	previousNow := nowUTC
	nowUTC = func() time.Time { return time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { nowUTC = previousNow })
	target := filepath.Join(alias, "memory")
	result := handleDreamScan(commandruntime.ExecutionContext{RepositoryRoot: alias}, []string{"--path", target})
	stale := []string{}
	relativeEvidence := ""
	for _, finding := range result.Findings {
		if strings.Join(finding.NextArgv, "\x00") != strings.Join([]string{"do-work-cli", CommandDreamScan, "--path", target}, "\x00") {
			t.Fatalf("target lost: %v", finding.NextArgv)
		}
		if finding.Code == "DREAM-STALE-PAGE" {
			stale = append(stale, strings.Join(finding.AffectedPaths, ""))
		}
		if finding.Code == "DREAM-RELATIVE-DATE" {
			relativeEvidence = strings.Join(finding.Evidence, " ")
		}
	}
	if len(stale) != 1 || !strings.Contains(stale[0], "ninety-one.md") {
		t.Fatalf("90/91 boundary stale=%v", stale)
	}
	for _, vocabulary := range []string{"yesterday", "today", "tomorrow", "tonight", "last week", "next month", "this evening", "a few days ago", "recently", "just now", "earlier today", "the other day"} {
		if !strings.Contains(relativeEvidence, vocabulary) {
			t.Errorf("missing vocabulary %q in %q", vocabulary, relativeEvidence)
		}
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
