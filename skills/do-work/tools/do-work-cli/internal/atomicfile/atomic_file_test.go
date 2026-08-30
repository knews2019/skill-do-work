package atomicfile

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestReplaceExistingPublishesWholeContentsAndPreservesMode(t *testing.T) {
	targetPath := filepath.Join(t.TempDir(), "request.md")
	if err := os.WriteFile(targetPath, []byte("old"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := ReplaceExisting(targetPath, []byte("new document")); err != nil {
		t.Fatalf("ReplaceExisting: %v", err)
	}
	contents, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "new document" {
		t.Fatalf("contents = %q, want complete replacement", contents)
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("mode = %o, want 640", info.Mode().Perm())
	}
}

func TestReplaceExistingRefusesUnsafeAndChangedTargets(t *testing.T) {
	testDirectory := t.TempDir()
	realPath := filepath.Join(testDirectory, "real")
	if err := os.WriteFile(realPath, []byte("outside"), 0o600); err != nil {
		t.Fatal(err)
	}
	symlinkPath := filepath.Join(testDirectory, "link")
	if err := os.Symlink(realPath, symlinkPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := ReplaceExisting(symlinkPath, []byte("changed")); !errors.Is(err, ErrUnsafeTarget) {
		t.Fatalf("symlink error = %v, want ErrUnsafeTarget", err)
	}

	targetPath := filepath.Join(testDirectory, "race")
	if err := os.WriteFile(targetPath, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	previousHook := beforeAtomicReplace
	beforeAtomicReplace = func(path string) {
		beforeAtomicReplace = previousHook
		replacementPath := path + ".changed"
		if err := os.WriteFile(replacementPath, []byte("other writer"), 0o600); err != nil {
			t.Fatalf("write competing target: %v", err)
		}
		if err := os.Rename(replacementPath, path); err != nil {
			t.Fatalf("replace competing target: %v", err)
		}
	}
	t.Cleanup(func() { beforeAtomicReplace = previousHook })
	if err := ReplaceExisting(targetPath, []byte("ours")); !errors.Is(err, ErrTargetChanged) {
		t.Fatalf("changed-target error = %v, want ErrTargetChanged", err)
	}
	contents, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "other writer" {
		t.Fatalf("changed target overwritten: %q", contents)
	}
}

func TestReplaceExistingDetectsInPlaceContentMutation(t *testing.T) {
	targetPath := filepath.Join(t.TempDir(), "in-place")
	if err := os.WriteFile(targetPath, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	originalInfo, err := os.Stat(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	previousHook := beforeAtomicReplace
	beforeAtomicReplace = func(path string) {
		beforeAtomicReplace = previousHook
		if writeError := os.WriteFile(path, []byte("mutated!"), 0o600); writeError != nil {
			t.Fatalf("mutate target in place: %v", writeError)
		}
		if timeError := os.Chtimes(path, originalInfo.ModTime(), originalInfo.ModTime()); timeError != nil {
			t.Fatalf("restore target timestamp: %v", timeError)
		}
	}
	t.Cleanup(func() { beforeAtomicReplace = previousHook })
	if replaceError := ReplaceExisting(targetPath, []byte("our data")); !errors.Is(replaceError, ErrTargetChanged) {
		t.Fatalf("in-place mutation error = %v, want ErrTargetChanged", replaceError)
	}
	contents, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "mutated!" {
		t.Fatalf("in-place mutation overwritten: %q", contents)
	}
}

func TestCreateExclusiveAllowsExactlyOneConcurrentReservation(t *testing.T) {
	reservationPath := filepath.Join(t.TempDir(), "REQ-000123")
	start := make(chan struct{})
	errorsSeen := make(chan error, 8)
	var waitGroup sync.WaitGroup
	for workerIndex := 0; workerIndex < 8; workerIndex++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			<-start
			errorsSeen <- CreateExclusive(reservationPath, []byte("reserved\n"), 0o644)
		}()
	}
	close(start)
	waitGroup.Wait()
	close(errorsSeen)
	successCount := 0
	for createError := range errorsSeen {
		if createError == nil {
			successCount++
			continue
		}
		if !errors.Is(createError, os.ErrExist) {
			t.Errorf("loser error = %v, want os.ErrExist", createError)
		}
	}
	if successCount != 1 {
		t.Fatalf("successful reservations = %d, want 1", successCount)
	}
}
