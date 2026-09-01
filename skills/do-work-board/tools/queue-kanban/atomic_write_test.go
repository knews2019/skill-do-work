package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicallyPreservesCompleteMode(t *testing.T) {
	for _, test := range []struct {
		name string
		mode os.FileMode
	}{
		{name: "setuid", mode: 0o4640},
		{name: "setgid", mode: 0o2640},
		{name: "sticky", mode: 0o1640},
	} {
		t.Run(test.name, func(t *testing.T) {
			targetPath := filepath.Join(t.TempDir(), "request.md")
			if err := os.WriteFile(targetPath, []byte("old contents"), test.mode.Perm()); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(targetPath, queueKanbanGoModeFromUnix(test.mode)); err != nil {
				t.Fatal(err)
			}

			if err := writeFileAtomically(targetPath, []byte("complete replacement")); err != nil {
				t.Fatalf("writeFileAtomically: %v", err)
			}
			contents, err := os.ReadFile(targetPath)
			if err != nil {
				t.Fatal(err)
			}
			if string(contents) != "complete replacement" {
				t.Fatalf("contents = %q, want complete replacement", contents)
			}
			if mode := queueKanbanUnixModeOf(t, targetPath); mode != test.mode {
				t.Fatalf("mode = %04o, want %04o", mode, test.mode)
			}
		})
	}
}

func queueKanbanGoModeFromUnix(mode os.FileMode) os.FileMode {
	goMode := mode.Perm()
	if mode&0o4000 != 0 {
		goMode |= os.ModeSetuid
	}
	if mode&0o2000 != 0 {
		goMode |= os.ModeSetgid
	}
	if mode&0o1000 != 0 {
		goMode |= os.ModeSticky
	}
	return goMode
}

func queueKanbanUnixModeOf(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	mode := info.Mode().Perm()
	if info.Mode()&os.ModeSetuid != 0 {
		mode |= 0o4000
	}
	if info.Mode()&os.ModeSetgid != 0 {
		mode |= 0o2000
	}
	if info.Mode()&os.ModeSticky != 0 {
		mode |= 0o1000
	}
	return mode
}
