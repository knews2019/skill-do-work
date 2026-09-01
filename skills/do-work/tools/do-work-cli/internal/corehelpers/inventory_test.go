package corehelpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
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

func TestProtectedInventoryPersistsLaterXAndRequiresStartedState(t *testing.T) {
	repository := newGitFixture(t)
	if err := os.WriteFile(filepath.Join(repository, "first.txt"), []byte("changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	start := handleProtectedInventory(testContext(repository), []string{"start"})
	if start.Outcome != "success" || !hasFinding(start, "INVENTORY-M") {
		t.Fatalf("start=%#v", start)
	}
	quarantineOutput, err := gitOutput(repository, "rev-parse", "--git-path", "do-work-commit-secret-quarantine")
	if err != nil {
		t.Fatal(err)
	}
	quarantinePath := strings.TrimSpace(string(quarantineOutput))
	if !filepath.IsAbs(quarantinePath) {
		quarantinePath = filepath.Join(repository, quarantinePath)
	}
	if contents, err := os.ReadFile(quarantinePath); err != nil || len(contents) != 0 {
		t.Fatalf("initial quarantine=%q err=%v", contents, err)
	}
	secretPath := filepath.Join(repository, ".env.local")
	if err := os.WriteFile(secretPath, []byte("TOKEN=x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	associate := handleProtectedInventory(testContext(repository), []string{"associate"})
	if associate.Outcome != "success" {
		t.Fatalf("associate=%#v", associate)
	}
	contents, err := os.ReadFile(quarantinePath)
	if err != nil || string(contents) != ".env.local\n" {
		t.Fatalf("persisted quarantine=%q err=%v", contents, err)
	}
	if err := os.Remove(secretPath); err != nil {
		t.Fatal(err)
	}
	second := handleProtectedInventory(testContext(repository), []string{"associate"})
	if second.Outcome != "success" {
		t.Fatalf("second associate=%#v", second)
	}
	contents, _ = os.ReadFile(quarantinePath)
	if string(contents) != ".env.local\n" {
		t.Fatalf("durable quarantine was lost: %q", contents)
	}

	unstarted := newGitFixture(t)
	if err := os.WriteFile(filepath.Join(unstarted, "first.txt"), []byte("changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	missing := handleProtectedInventory(testContext(unstarted), []string{"associate"})
	if missing.Outcome != "failure" {
		t.Fatalf("missing start state accepted: %#v", missing)
	}
}

func TestProtectedInventoryCleanStatusAndXDIsNotQuarantined(t *testing.T) {
	repository := newGitFixture(t)
	clean := handleProtectedInventory(testContext(repository), []string{"start"})
	if clean.Outcome != "findings" || !hasFinding(clean, "INVENTORY-CLEAN") {
		t.Fatalf("clean=%#v", clean)
	}
	secret := filepath.Join(repository, "secret.pem")
	if err := os.WriteFile(secret, []byte("secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runFixtureGitCommand(t, repository, "add", "secret.pem")
	runFixtureGitCommand(t, repository, "commit", "-qm", "secret fixture")
	if err := os.Remove(secret); err != nil {
		t.Fatal(err)
	}
	result := handleProtectedInventory(testContext(repository), []string{"start"})
	if result.Outcome != "success" || !hasFinding(result, "INVENTORY-XD") {
		t.Fatalf("XD start=%#v", result)
	}
	quarantineOutput, _ := gitOutput(repository, "rev-parse", "--git-path", "do-work-commit-secret-quarantine")
	quarantinePath := strings.TrimSpace(string(quarantineOutput))
	if !filepath.IsAbs(quarantinePath) {
		quarantinePath = filepath.Join(repository, quarantinePath)
	}
	if contents, err := os.ReadFile(quarantinePath); err != nil || len(contents) != 0 {
		t.Fatalf("XD was quarantined: %q err=%v", contents, err)
	}
}

func TestAssociateEmptyAndProtectedDryRunPreserveStatusAndState(t *testing.T) {
	repository := newGitFixture(t)
	empty := filepath.Join(repository, "empty-paths")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	association := handleAssociate(testContext(repository), []string{"--paths-file", empty})
	if association.Outcome != "findings" || !hasFinding(association, "ASSOCIATION-NO-CANDIDATES") {
		t.Fatalf("empty association=%#v", association)
	}
	if err := os.WriteFile(filepath.Join(repository, ".env.local"), []byte("TOKEN=x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result := handleProtectedInventory(testContext(repository), []string{"start", "--dry-run"})
	if result.Outcome != "success" || !hasFinding(result, "INVENTORY-X") {
		t.Fatalf("protected dry-run=%#v", result)
	}
	quarantineOutput, _ := gitOutput(repository, "rev-parse", "--git-path", "do-work-commit-secret-quarantine")
	quarantinePath := strings.TrimSpace(string(quarantineOutput))
	if !filepath.IsAbs(quarantinePath) {
		quarantinePath = filepath.Join(repository, quarantinePath)
	}
	if _, err := os.Stat(quarantinePath); !os.IsNotExist(err) {
		t.Fatalf("protected dry-run wrote quarantine: %v", err)
	}
}

func hasFinding(result resultmodel.CommandResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
