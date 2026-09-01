package corehelpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrivateCopyPublishesExactBytesWithoutOverwrite(t *testing.T) {
	directory := t.TempDir()
	source := filepath.Join(directory, "source.png")
	target := filepath.Join(directory, "out", "capture.png")
	if err := os.WriteFile(source, []byte("image bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := publishPrivateCopy(source, target); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(target)
	if err != nil || string(contents) != "image bytes" {
		t.Fatalf("contents=%q err=%v", contents, err)
	}
	info, _ := os.Stat(target)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%o", info.Mode().Perm())
	}
	if err := publishPrivateCopy(source, target); err == nil {
		t.Fatal("occupied target was overwritten")
	}
}
func TestPrivateCopyRefusesSymlinkSource(t *testing.T) {
	directory := t.TempDir()
	real := filepath.Join(directory, "real")
	_ = os.WriteFile(real, []byte("x"), 0o600)
	link := filepath.Join(directory, "link")
	if err := os.Symlink(real, link); err != nil {
		t.Skip(err)
	}
	if err := publishPrivateCopy(link, filepath.Join(directory, "target")); err == nil {
		t.Fatal("symlink source accepted")
	}
}

func TestPrivateCopyParentSwapCannotRedirectPublicationOrCleanup(t *testing.T) {
	directory := t.TempDir()
	source := filepath.Join(directory, "source")
	parent := filepath.Join(directory, "parent")
	held := filepath.Join(directory, "held")
	outside := filepath.Join(directory, "outside")
	_ = os.WriteFile(source, []byte("new"), 0o600)
	_ = os.Mkdir(parent, 0o755)
	_ = os.Mkdir(outside, 0o755)
	protected := filepath.Join(outside, "capture")
	_ = os.WriteFile(protected, []byte("protected"), 0o600)
	originalHook := beforePrivateCopyPublish
	defer func() { beforePrivateCopyPublish = originalHook }()
	beforePrivateCopyPublish = func() {
		if err := os.Rename(parent, held); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, parent); err != nil {
			t.Fatal(err)
		}
	}
	if err := publishPrivateCopy(source, filepath.Join(parent, "capture")); err != nil {
		t.Fatal(err)
	}
	if contents, _ := os.ReadFile(protected); string(contents) != "protected" {
		t.Fatalf("outside file changed: %q", contents)
	}
	if contents, _ := os.ReadFile(filepath.Join(held, "capture")); string(contents) != "new" {
		t.Fatalf("rooted publication missing: %q", contents)
	}
}
