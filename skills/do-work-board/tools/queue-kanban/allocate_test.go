package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
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

// Every successful allocation reserves its number. Sequential callers must
// therefore receive distinct ids even when neither has written a REQ file yet.
func TestNextRequestNumberReservesNumberForTheNextCall(t *testing.T) {
	repoRoot := writeAllocationFixture(t, []string{"do-work/queue/REQ-042-only.md"})

	firstNumber, firstError := nextRequestNumber(repoRoot)
	if firstError != nil {
		t.Fatalf("first nextRequestNumber: %v", firstError)
	}
	secondNumber, secondError := nextRequestNumber(repoRoot)
	if secondError != nil {
		t.Fatalf("second nextRequestNumber: %v", secondError)
	}
	if firstNumber != 43 || secondNumber != 44 {
		t.Fatalf("sequential reservations = (%d, %d), want (43, 44)", firstNumber, secondNumber)
	}

	for _, reservedNumber := range []int{43, 44} {
		reservationPath := filepath.Join(repoRoot, "do-work", ".req-reservations", requestReservationFileName(reservedNumber))
		if _, statError := os.Stat(reservationPath); statError != nil {
			t.Errorf("reservation %d missing at %s: %v", reservedNumber, reservationPath, statError)
		}
	}
}

// Exclusive marker creation is the concurrency boundary: separate allocator
// processes that all see the same highest REQ must race safely and each return a
// different number.
func TestNextRequestNumberConcurrentProcessesReserveDistinctNumbers(t *testing.T) {
	repoRoot := writeAllocationFixture(t, nil)
	const callerCount = 16

	allocatedNumbers := make(chan int, callerCount)
	allocationErrors := make(chan error, callerCount)
	var processes sync.WaitGroup
	for processIndex := 0; processIndex < callerCount; processIndex++ {
		processes.Add(1)
		go func() {
			defer processes.Done()
			allocatorProcess := exec.Command(os.Args[0], "-test.run=^TestNextRequestNumberProcessHelper$", "--", repoRoot)
			allocatorProcess.Env = testEnvironmentWithOverrides(os.Environ(), "QUEUE_KANBAN_ALLOCATOR_HELPER=1", strictJavaScriptBehaviorMarker+"=")
			processOutput, processError := allocatorProcess.CombinedOutput()
			if processError != nil {
				allocationErrors <- fmt.Errorf("allocator process: %w: %s", processError, strings.TrimSpace(string(processOutput)))
				return
			}
			allocatedNumber, parseError := strconv.Atoi(strings.TrimSpace(string(processOutput)))
			if parseError != nil {
				allocationErrors <- fmt.Errorf("allocator output %q is not one decimal number: %w", processOutput, parseError)
				return
			}
			allocatedNumbers <- allocatedNumber
		}()
	}
	processes.Wait()
	close(allocatedNumbers)
	close(allocationErrors)

	for allocateError := range allocationErrors {
		t.Errorf("concurrent nextRequestNumber: %v", allocateError)
	}
	var sortedNumbers []int
	for allocatedNumber := range allocatedNumbers {
		sortedNumbers = append(sortedNumbers, allocatedNumber)
	}
	if len(sortedNumbers) != callerCount {
		t.Fatalf("successful allocations = %d, want %d", len(sortedNumbers), callerCount)
	}
	sort.Ints(sortedNumbers)
	for numberIndex, allocatedNumber := range sortedNumbers {
		expectedNumber := numberIndex + 1
		if allocatedNumber != expectedNumber {
			t.Fatalf("sorted allocation %d = %d, want %d; all=%v", numberIndex, allocatedNumber, expectedNumber, sortedNumbers)
		}
	}
}

// TestNextRequestNumberProcessHelper is re-entered through the test binary so
// the concurrency test exercises OS processes rather than goroutines sharing
// one Go runtime. It exits immediately to keep stdout to exactly one decimal
// number and a newline, matching the next-req command contract.
func TestNextRequestNumberProcessHelper(t *testing.T) {
	if os.Getenv("QUEUE_KANBAN_ALLOCATOR_HELPER") != "1" {
		return
	}
	argumentSeparator := -1
	for argumentIndex, argumentValue := range os.Args {
		if argumentValue == "--" {
			argumentSeparator = argumentIndex
			break
		}
	}
	if argumentSeparator < 0 || argumentSeparator+1 >= len(os.Args) {
		fmt.Fprintln(os.Stderr, "missing allocator repository root")
		os.Exit(2)
	}
	allocatedNumber, allocateError := nextRequestNumber(os.Args[argumentSeparator+1])
	if allocateError != nil {
		fmt.Fprintln(os.Stderr, allocateError)
		os.Exit(2)
	}
	fmt.Println(allocatedNumber)
	os.Exit(0)
}

// Reservation markers remain authoritative even when no matching REQ exists —
// an interrupted capture consumes a number rather than making it reusable.
func TestNextRequestNumberAdvancesPastExistingReservation(t *testing.T) {
	repoRoot := writeAllocationFixture(t, []string{"do-work/queue/REQ-042-only.md"})
	reservationDirectory := filepath.Join(repoRoot, "do-work", requestReservationDirectoryName)
	if mkdirError := os.Mkdir(reservationDirectory, 0o755); mkdirError != nil {
		t.Fatalf("mkdir reservation directory: %v", mkdirError)
	}
	if writeError := os.WriteFile(filepath.Join(reservationDirectory, requestReservationFileName(88)), nil, 0o644); writeError != nil {
		t.Fatalf("write existing reservation: %v", writeError)
	}

	allocatedNumber, allocateError := nextRequestNumber(repoRoot)
	if allocateError != nil {
		t.Fatalf("nextRequestNumber: %v", allocateError)
	}
	if allocatedNumber != 89 {
		t.Fatalf("nextRequestNumber = %d, want 89 after reserved REQ-88", allocatedNumber)
	}
}

// A reservation directory symlink could redirect next-req's new write surface
// outside do-work. Fail closed instead of following it.
func TestNextRequestNumberRejectsSymlinkedReservationDirectory(t *testing.T) {
	repoRoot := writeAllocationFixture(t, nil)
	outsideDirectory := t.TempDir()
	reservationDirectory := filepath.Join(repoRoot, "do-work", requestReservationDirectoryName)
	if symlinkError := os.Symlink(outsideDirectory, reservationDirectory); symlinkError != nil {
		t.Skipf("symlinks unavailable: %v", symlinkError)
	}

	if _, allocateError := nextRequestNumber(repoRoot); allocateError == nil {
		t.Fatal("nextRequestNumber succeeded through a symlinked reservation directory")
	}
	outsideEntries, readError := os.ReadDir(outsideDirectory)
	if readError != nil {
		t.Fatalf("read outside directory: %v", readError)
	}
	if len(outsideEntries) != 0 {
		t.Fatalf("nextRequestNumber wrote outside do-work through a symlink: %v", outsideEntries)
	}
}

// The parent metadata directory is untrusted too: a repository-local do-work
// symlink must not redirect reservation markers anywhere outside the checkout.
func TestNextRequestNumberRejectsDoWorkRootOutsideRepository(t *testing.T) {
	repoRoot := t.TempDir()
	outsideDirectory := t.TempDir()
	if mkdirError := os.Mkdir(filepath.Join(outsideDirectory, "queue"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir outside queue: %v", mkdirError)
	}
	if symlinkError := os.Symlink(outsideDirectory, filepath.Join(repoRoot, "do-work")); symlinkError != nil {
		t.Skipf("symlinks unavailable: %v", symlinkError)
	}

	if _, allocateError := nextRequestNumber(repoRoot); allocateError == nil {
		t.Fatal("nextRequestNumber succeeded through a do-work symlink outside the repository")
	}
	if _, statError := os.Lstat(filepath.Join(outsideDirectory, requestReservationDirectoryName)); !os.IsNotExist(statError) {
		t.Fatalf("nextRequestNumber created an outside reservation path: %v", statError)
	}
}

// Reserving a number must not rewrite or add a REQ in queue/. The only new entry
// belongs in the dedicated reservation store.
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
		t.Errorf("queue entry count changed from %d to %d; reservations must not create REQ files", len(beforeEntries), len(afterEntries))
	}
}
