package toolboxcommands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMutationFlags(t *testing.T) {
	rest, dry, commit, err := parseMutationFlags([]string{"--dry-run", "x"})
	if err != nil || !dry || commit || len(rest) != 1 || rest[0] != "x" {
		t.Fatalf("unexpected parse: %v %v %v %v", rest, dry, commit, err)
	}
	if _, _, _, err := parseMutationFlags([]string{"--dry-run", "--commit"}); err == nil {
		t.Fatal("combined flags accepted")
	}
}

func TestRemediationRootedBoundaryRejectsLinkedAncestor(t *testing.T) {
	repository := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(repository, "reports")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := repositoryPath(repository, filepath.Join(repository, "reports", "run", "index.html")); err == nil {
		t.Fatal("linked ancestor passed repository confinement")
	}
	if _, err := os.Stat(filepath.Join(outside, "run")); !os.IsNotExist(err) {
		t.Fatal("confinement validation wrote outside repository")
	}
}
