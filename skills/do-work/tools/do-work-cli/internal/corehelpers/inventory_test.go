package corehelpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInventoryClassificationProtectsSecretProvenance(t *testing.T) {
	tests := []struct{ status, path, origin, want string }{{"??", "source.go", "", "A"}, {"R ", "public.txt", ".env", "X"}, {" D", "secret.pem", "", "XD"}, {" M", "source.go", "", "M"}}
	for _, test := range tests {
		if got := classifyInventory(test.status, test.path, test.origin); got != test.want {
			t.Errorf("%q %q from %q=%s want %s", test.status, test.path, test.origin, got, test.want)
		}
	}
}

func TestInventoryPreservesSecretRenameOriginAndQuarantinesAmbiguousAdds(t *testing.T) {
	repository := newGitFixture(t)
	secret := filepath.Join(repository, ".envrc")
	if err := os.WriteFile(secret, []byte("TOKEN=x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runFixtureGitCommand(t, repository, "add", ".envrc")
	runFixtureGitCommand(t, repository, "commit", "-qm", "secret fixture")
	runFixtureGitCommand(t, repository, "mv", ".envrc", "visible-config.txt")
	if err := os.WriteFile(filepath.Join(repository, "ordinary.txt"), []byte("copy\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	rows, err := readInventory(repository)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{".envrc": "XD", "visible-config.txt": "X", "ordinary.txt": "X"}
	for _, row := range rows {
		if expected, ok := want[row.Path]; ok {
			if row.Classification != expected {
				t.Fatalf("%s=%s want %s (rows=%+v)", row.Path, row.Classification, expected, rows)
			}
			delete(want, row.Path)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing=%v rows=%+v", want, rows)
	}
}

func runFixtureGitCommand(t *testing.T, repository string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repository}, args...)...)
	command.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com", "GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

func TestTerminalSuccessAliases(t *testing.T) {
	for _, status := range []string{"completed", "completed-with-issues", "complete", "done", "finished", "closed", " DONE "} {
		if !terminalSuccessStatus(status) {
			t.Errorf("status %q was not terminal success", status)
		}
	}
	if terminalSuccessStatus("cancelled") {
		t.Fatal("cancelled request was treated as successful ownership evidence")
	}
}
