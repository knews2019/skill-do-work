package rootedfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func openTestRoot(t *testing.T) (*os.Root, string) {
	t.Helper()
	directory := t.TempDir()
	root, err := os.OpenRoot(directory)
	if err != nil {
		t.Fatalf("open root: %v", err)
	}
	t.Cleanup(func() { _ = root.Close() })
	return root, directory
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// Rename must move within the root and must refuse a source whose directory is a
// symlink out of the root: the Go 1.25 method refuses it, and the rollback quarantine
// relies on that to never move a foreign file.
func TestRenameStaysInsideTheRoot(t *testing.T) {
	root, directory := openTestRoot(t)
	writeTestFile(t, filepath.Join(directory, "nested", "object"), "payload")
	if err := os.Mkdir(filepath.Join(directory, "quarantine"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := Rename(root, filepath.Join("nested", "object"), filepath.Join("quarantine", "object")); err != nil {
		t.Fatalf("rename inside root: %v", err)
	}
	moved, err := os.ReadFile(filepath.Join(directory, "quarantine", "object"))
	if err != nil || string(moved) != "payload" {
		t.Fatalf("renamed payload = %q, %v", moved, err)
	}

	outside := t.TempDir()
	writeTestFile(t, filepath.Join(outside, "victim"), "foreign")
	if err := os.Symlink(outside, filepath.Join(directory, "escape")); err != nil {
		t.Fatal(err)
	}
	err = Rename(root, filepath.Join("escape", "victim"), "stolen")
	if err == nil {
		t.Fatal("rename through an escaping symlink succeeded")
	}
	if _, statError := os.Stat(filepath.Join(outside, "victim")); statError != nil {
		t.Fatalf("foreign file was moved: %v", statError)
	}
	if err := Rename(root, filepath.Join("quarantine", "object"), filepath.Join("escape", "planted")); err == nil {
		t.Fatal("rename into an escaping directory succeeded")
	}
	if err := Rename(root, filepath.Join(directory, "quarantine", "object"), "absolute"); err == nil {
		t.Fatal("absolute source path was accepted")
	}
}

// A root whose name no longer identifies its directory must be refused, regardless
// of how the Go version names nested roots. Our nested opener remains usable.
func TestNestedRootsRefuseUnresolvableNames(t *testing.T) {
	root, directory := openTestRoot(t)
	writeTestFile(t, filepath.Join(directory, "nested", "object"), "payload")
	goNested, err := root.OpenRoot("nested")
	if err != nil {
		t.Fatal(err)
	}
	defer goNested.Close()
	if err := os.Rename(filepath.Join(directory, "nested"), filepath.Join(directory, "held")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(directory, "nested", "object"), "replacement")
	if err := Rename(goNested, "object", "moved"); !errors.Is(err, errRootUnresolvable) {
		t.Fatalf("rename inside an unresolvable root error = %v, want unresolvable root", err)
	}
	for _, name := range []string{"nested", "held"} {
		if _, err := root.Lstat(filepath.Join(name, "object")); err != nil {
			t.Fatalf("object moved despite refusal: %v", err)
		}
	}
	nested, err := OpenRoot(root, "nested")
	if err != nil {
		t.Fatalf("rootedfs.OpenRoot: %v", err)
	}
	defer nested.Close()
	if err := Rename(nested, "object", "moved"); err != nil {
		t.Fatalf("rename inside a rootedfs-nested root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(directory, "nested", "moved")); err != nil {
		t.Fatalf("object was not moved: %v", err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(directory, "escape")); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenRoot(root, "escape"); err == nil {
		t.Fatal("OpenRoot followed a symlink out of the root")
	}
}

// Link is the exclusive publish step of a download: it must create the target once and
// refuse when the target appeared in the meantime.
func TestLinkPublishesExclusively(t *testing.T) {
	root, directory := openTestRoot(t)
	writeTestFile(t, filepath.Join(directory, ".stage"), "archive")
	if err := Link(root, ".stage", "archive.tar.gz"); err != nil {
		t.Fatalf("first link: %v", err)
	}
	published, err := os.ReadFile(filepath.Join(directory, "archive.tar.gz"))
	if err != nil || string(published) != "archive" {
		t.Fatalf("published payload = %q, %v", published, err)
	}
	err = Link(root, ".stage", "archive.tar.gz")
	if !errors.Is(err, fs.ErrExist) {
		t.Fatalf("second link error = %v, want ErrExist", err)
	}
	if err := Link(root, ".stage", filepath.Join("missing-parent", "archive.tar.gz")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("link into a missing parent error = %v, want ErrNotExist", err)
	}
}

// MkdirAll must create nested parents through the root, accept what already exists, and
// refuse to treat a regular file as a directory.
func TestMkdirAllCreatesThroughTheRoot(t *testing.T) {
	root, directory := openTestRoot(t)
	nested := filepath.Join("one", "two", "three")
	if err := MkdirAll(root, nested, 0o755); err != nil {
		t.Fatalf("mkdir all: %v", err)
	}
	if info, err := os.Stat(filepath.Join(directory, nested)); err != nil || !info.IsDir() {
		t.Fatalf("nested directory missing: %v", err)
	}
	if err := MkdirAll(root, nested, 0o755); err != nil {
		t.Fatalf("mkdir all on an existing tree: %v", err)
	}
	if err := MkdirAll(root, ".", 0o755); err != nil {
		t.Fatalf("mkdir all on the root itself: %v", err)
	}
	writeTestFile(t, filepath.Join(directory, "file"), "")
	if err := MkdirAll(root, filepath.Join("file", "child"), 0o755); err == nil {
		t.Fatal("a regular file was accepted as a parent directory")
	}
	if err := MkdirAll(root, filepath.Join("..", "escaped"), 0o755); err == nil {
		t.Fatal("a path above the root was created")
	}
}

// Replacing an enumerated directory must not redirect deletion into a sibling.
func TestRemoveAllKeepsEnumeratedDirectoryBound(t *testing.T) {
	root, directory := openTestRoot(t)
	writeTestFile(t, filepath.Join(directory, "owned", "matching"), "owned")
	writeTestFile(t, filepath.Join(directory, "unrelated", "matching"), "foreign")
	previous := afterRemoveAllEnumeration
	t.Cleanup(func() { afterRemoveAllEnumeration = previous })
	afterRemoveAllEnumeration = func(path string) {
		if path != "owned" {
			return
		}
		if err := os.Rename(filepath.Join(directory, "owned"), filepath.Join(directory, "held")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("unrelated", filepath.Join(directory, "owned")); err != nil {
			t.Fatal(err)
		}
	}
	if err := RemoveAll(root, "owned"); err != nil {
		t.Fatal(err)
	}
	contents, err := ReadFile(root, "unrelated/matching")
	if err != nil || string(contents) != "foreign" {
		t.Fatalf("unrelated file changed: %q, %v", contents, err)
	}
	if _, err := root.Lstat("held/matching"); !os.IsNotExist(err) {
		t.Fatalf("enumerated file was not removed: %v", err)
	}
}
