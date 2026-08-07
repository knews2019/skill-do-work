package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeAllocationFixture seeds a synthetic do-work tree and returns its repo
// root. Each entry is a path relative to the repo root; content is irrelevant to
// allocation (the number lives in the filename), so a stub body is written.
func writeAllocationFixture(t *testing.T, relativePaths []string) string {
	t.Helper()
	repoRoot := t.TempDir()
	// An empty tree still needs do-work/ to exist — enumerateDoWorkTree errors
	// without it, and "no REQs yet" must return 1, not an error.
	if mkdirError := os.MkdirAll(filepath.Join(repoRoot, "do-work", "queue"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir do-work/queue: %v", mkdirError)
	}
	for _, relativePath := range relativePaths {
		absolutePath := filepath.Join(repoRoot, relativePath)
		if mkdirError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirError != nil {
			t.Fatalf("mkdir for %s: %v", relativePath, mkdirError)
		}
		body := "---\nid: " + reqIdFromFixturePath(relativePath) + "\nstatus: pending\ntitle: fixture\n---\n\n# Fixture\n"
		if writeError := os.WriteFile(absolutePath, []byte(body), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", relativePath, writeError)
		}
	}
	return repoRoot
}

// reqIdFromFixturePath pulls "REQ-NNN" off a fixture filename so the stub
// frontmatter agrees with the filename — a disagreement would be testing the
// wrong thing.
func reqIdFromFixturePath(relativePath string) string {
	baseName := filepath.Base(relativePath)
	if len(baseName) < 7 {
		return baseName
	}
	return baseName[:7]
}

func TestNextRequestNumber(t *testing.T) {
	testCases := []struct {
		caseName       string
		fixturePaths   []string
		expectedNumber int
	}{
		{
			caseName:       "highest id is REQ-070 in the queue",
			fixturePaths:   []string{"do-work/queue/REQ-069-earlier.md", "do-work/queue/REQ-070-latest.md"},
			expectedNumber: 71,
		},
		{
			caseName:       "empty tree starts at 1",
			fixturePaths:   nil,
			expectedNumber: 1,
		},
		{
			caseName:       "a gap in the sequence is tolerated, never filled",
			fixturePaths:   []string{"do-work/queue/REQ-001-first.md", "do-work/queue/REQ-070-latest.md"},
			expectedNumber: 71,
		},
		{
			caseName: "the max spans queue, working, and nested archive UR folders",
			fixturePaths: []string{
				"do-work/queue/REQ-005-queued.md",
				"do-work/working/REQ-011-claimed.md",
				"do-work/archive/REQ-020-legacy.md",
				"do-work/archive/UR-012/REQ-069-consolidated.md",
			},
			expectedNumber: 70,
		},
		{
			caseName:       "zero-padding width does not cap the number",
			fixturePaths:   []string{"do-work/queue/REQ-0142-wide.md"},
			expectedNumber: 143,
		},
		{
			caseName: "pruned subtrees do not contribute",
			fixturePaths: []string{
				"do-work/queue/REQ-007-real.md",
				"do-work/deliverables/REQ-900-report-copy.md",
				"do-work/runs/work-2026-08-03/REQ-901-brief.md",
				"do-work/archive/UR-003/assets/REQ-902-deliverable-copy.md",
			},
			expectedNumber: 8,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			repoRoot := writeAllocationFixture(t, testCase.fixturePaths)
			allocatedNumber, allocateError := nextRequestNumber(repoRoot)
			if allocateError != nil {
				t.Fatalf("nextRequestNumber: unexpected error: %v", allocateError)
			}
			if allocatedNumber != testCase.expectedNumber {
				t.Errorf("nextRequestNumber = %d, want %d", allocatedNumber, testCase.expectedNumber)
			}
		})
	}
}

// A REQ file whose frontmatter id disagrees with its filename must not let the
// allocator hand out a number that is already taken. Both are consulted; the
// higher wins.
func TestNextRequestNumberHonorsFrontmatterIdAboveFilename(t *testing.T) {
	repoRoot := t.TempDir()
	queueDirectory := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDirectory, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	mismatchedBody := "---\nid: REQ-088\nstatus: pending\ntitle: renamed file, original id\n---\n\n# Body\n"
	if writeError := os.WriteFile(filepath.Join(queueDirectory, "REQ-012-renamed.md"), []byte(mismatchedBody), 0o644); writeError != nil {
		t.Fatalf("write: %v", writeError)
	}

	allocatedNumber, allocateError := nextRequestNumber(repoRoot)
	if allocateError != nil {
		t.Fatalf("nextRequestNumber: unexpected error: %v", allocateError)
	}
	if allocatedNumber != 89 {
		t.Errorf("nextRequestNumber = %d, want 89 (frontmatter id REQ-088 outranks filename REQ-012)", allocatedNumber)
	}
}

// The allocator is read-only with respect to the queue (requirement 1): it must
// not create, move, or rewrite anything under do-work/.
func TestNextRequestNumberLeavesTheQueueUntouched(t *testing.T) {
	repoRoot := writeAllocationFixture(t, []string{"do-work/queue/REQ-042-only.md"})
	queueFile := filepath.Join(repoRoot, "do-work", "queue", "REQ-042-only.md")

	beforeContent, readError := os.ReadFile(queueFile)
	if readError != nil {
		t.Fatalf("read before: %v", readError)
	}
	beforeEntries, listError := os.ReadDir(filepath.Join(repoRoot, "do-work", "queue"))
	if listError != nil {
		t.Fatalf("list before: %v", listError)
	}

	if _, allocateError := nextRequestNumber(repoRoot); allocateError != nil {
		t.Fatalf("nextRequestNumber: %v", allocateError)
	}

	afterContent, readError := os.ReadFile(queueFile)
	if readError != nil {
		t.Fatalf("read after: %v", readError)
	}
	if string(afterContent) != string(beforeContent) {
		t.Error("nextRequestNumber rewrote a queue file; it must be read-only toward the queue")
	}
	afterEntries, listError := os.ReadDir(filepath.Join(repoRoot, "do-work", "queue"))
	if listError != nil {
		t.Fatalf("list after: %v", listError)
	}
	if len(afterEntries) != len(beforeEntries) {
		t.Errorf("queue entry count changed from %d to %d; the allocator must not create files", len(beforeEntries), len(afterEntries))
	}
}
