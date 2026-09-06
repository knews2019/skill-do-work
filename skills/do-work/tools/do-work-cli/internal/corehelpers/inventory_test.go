package corehelpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestInventoryStagedAdditionDeletedFromWorktreeIsDeletion(t *testing.T) {
	repository := newGitFixture(t)
	path := filepath.Join(repository, "ephemeral.txt")
	if err := os.WriteFile(path, []byte("ephemeral\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runFixtureGitCommand(t, repository, "add", "ephemeral.txt")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("AD path must be absent, stat err=%v", err)
	}
	porcelain := runGitOutput(t, repository, "status", "--porcelain=v1", "--untracked-files=all")
	if !strings.Contains(porcelain, "AD ephemeral.txt\n") {
		t.Fatalf("raw porcelain does not contain AD fixture: %q", porcelain)
	}

	retained := runRetainedInventory(t, repository)
	if !reflect.DeepEqual(retained, []inventoryRow{{Classification: "D", Path: "ephemeral.txt"}}) {
		t.Fatalf("retained AD rows=%+v, want D ephemeral.txt", retained)
	}
	result := handleInventory(testContext(repository), nil)
	if result.Outcome != resultmodel.OutcomeSuccess || !hasFinding(result, "INVENTORY-D") || hasFinding(result, "INVENTORY-A") {
		t.Fatalf("Go AD result=%#v, want only INVENTORY-D", result)
	}
}

// Pins the real failure: a Finder .DS_Store under do-work/ was listed as an
// untracked addition, finalization discovery read it as shared state with no
// owner, and `do-work run` stopped before queue selection.
func TestInventoryDropsUntrackedHiddenFilesUnderDoWork(t *testing.T) {
	repository := newGitFixture(t)
	for path, contents := range map[string]string{
		"do-work/.DS_Store":                   "finder\n",
		"do-work/queue/.REQ-1-fixture.md.swp": "editor\n",
		"do-work/.req-reservations/041":       "reserved\n",
		"do-work/queue/notes.txt":             "real\n",
		".DS_Store":                           "finder at root\n",
	} {
		absolute := filepath.Join(repository, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(absolute, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := readInventory(repository)
	if err != nil {
		t.Fatal(err)
	}
	want := []inventoryRow{
		{Classification: "A", Path: ".DS_Store"},
		{Classification: "A", Path: "do-work/.req-reservations/041"},
		{Classification: "A", Path: "do-work/queue/notes.txt"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("rows=%+v, want %+v", rows, want)
	}
}

func TestInventoryMatchesRetainedPorcelainXYMatrix(t *testing.T) {
	statuses := []string{"??"}
	for _, worktree := range []byte{'M', 'T', 'A', 'D', 'R', 'C'} {
		statuses = append(statuses, " "+string(worktree))
	}
	for _, index := range []byte{'M', 'T', 'A', 'D', 'R', 'C'} {
		statuses = append(statuses, string(index)+" ")
	}
	for _, index := range []byte{'M', 'T', 'A', 'R', 'C'} {
		for _, worktree := range []byte{'M', 'T', 'D', 'R', 'C'} {
			statuses = append(statuses, string([]byte{index, worktree}))
		}
	}
	statuses = append(statuses, "DD", "AU", "UD", "UA", "DU", "AA", "UU")
	if len(statuses) != 45 {
		t.Fatalf("porcelain matrix has %d statuses, want 45", len(statuses))
	}

	for _, status := range statuses {
		status := status
		t.Run(strings.ReplaceAll(status, " ", "underscore"), func(t *testing.T) {
			origin := ""
			if strings.ContainsAny(status, "RC") {
				origin = "origin.txt"
			}
			statusBytes := porcelainStatusBytes(status, "ordinary.txt", origin)
			goRows, err := parseInventoryBytes(statusBytes)
			if err != nil {
				t.Fatalf("status %q parse: %v", status, err)
			}
			want := expectedOrdinaryInventoryClass(status)
			if len(goRows) != 1 || goRows[0].Classification != want || goRows[0].Path != "ordinary.txt" || goRows[0].Origin != origin {
				t.Fatalf("status %q Go rows=%+v, want class=%s path=ordinary.txt origin=%q", status, goRows, want, origin)
			}
			assertInventoryFindingProjection(t, goRows)
		})
	}
}

func TestInventoryMatchesRetainedSecretOriginAndAmbiguityMatrix(t *testing.T) {
	tests := []struct {
		name        string
		statusBytes []byte
		want        []inventoryRow
	}{
		{"secret addition", porcelainStatusBytes("A ", ".env.local", ""), []inventoryRow{{"X", ".env.local", ""}}},
		{"secret deletion", porcelainStatusBytes("AD", ".env.local", ""), []inventoryRow{{"XD", ".env.local", ""}}},
		{"secret unmerged addition", porcelainStatusBytes("UA", ".env.local", ""), []inventoryRow{{"X", ".env.local", ""}}},
		{"secret unmerged deletion", porcelainStatusBytes("UD", ".env.local", ""), []inventoryRow{{"XD", ".env.local", ""}}},
		{"secret rename origin", porcelainStatusBytes("R ", "visible.txt", ".env"), []inventoryRow{{"XD", ".env", ""}, {"X", "visible.txt", ".env"}}},
		{"deleted secret rename destination", porcelainStatusBytes("RD", "visible.txt", ".env"), []inventoryRow{{"XD", ".env", ""}, {"XD", "visible.txt", ".env"}}},
		{"secret copy origin", porcelainStatusBytes("C ", "visible.txt", ".env"), []inventoryRow{{"X", "visible.txt", ".env"}}},
		{"deleted secret copy destination", porcelainStatusBytes("CD", "visible.txt", ".env"), []inventoryRow{{"XD", "visible.txt", ".env"}}},
		{"excluded promotes addition", append(porcelainStatusBytes("??", "ordinary.txt", ""), porcelainStatusBytes(" M", "secret.pem", "")...), []inventoryRow{{"X", "ordinary.txt", ""}, {"X", "secret.pem", ""}}},
		{"secret deletion promotes addition", append(porcelainStatusBytes("??", "ordinary.txt", ""), porcelainStatusBytes(" D", "secret.pem", "")...), []inventoryRow{{"X", "ordinary.txt", ""}, {"XD", "secret.pem", ""}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			goRows, err := parseInventoryBytes(test.statusBytes)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if !reflect.DeepEqual(goRows, test.want) {
				t.Fatalf("Go rows=%+v, want %+v", goRows, test.want)
			}
			assertInventoryFindingProjection(t, goRows)
		})
	}
}

func TestInventoryDifferentialComparatorRejectsClassPathAndOrderMutations(t *testing.T) {
	statusBytes := append(porcelainStatusBytes(" M", "first.txt", ""), porcelainStatusBytes(" D", "second.txt", "")...)
	goRows, err := parseInventoryBytes(statusBytes)
	if err != nil {
		t.Fatal(err)
	}
	expected := []inventoryRow{{"M", "first.txt", ""}, {"D", "second.txt", ""}}
	if err := compareInventoryProjection(goRows, expected); err != nil {
		t.Fatal(err)
	}
	mutations := []struct {
		name string
		rows []inventoryRow
	}{
		{"class", []inventoryRow{{"A", "first.txt", ""}, expected[1]}},
		{"path", []inventoryRow{{expected[0].Classification, "wrong.txt", ""}, expected[1]}},
		{"order", []inventoryRow{expected[1], expected[0]}},
	}
	for _, mutation := range mutations {
		if err := compareInventoryProjection(mutation.rows, expected); err == nil {
			t.Errorf("%s mutation escaped retained differential", mutation.name)
		}
	}
}

func TestParseInventoryBytesContract(t *testing.T) {
	t.Run("malformed short records rejected", func(t *testing.T) {
		for _, short := range [][]byte{
			[]byte("M\x00"),
			[]byte("AD\x00"),
			[]byte("M \x00"),
			[]byte("M  \x00"),
		} {
			_, err := parseInventoryBytes(short)
			if err == nil || !strings.Contains(err.Error(), "short porcelain record") {
				t.Errorf("expected short porcelain record error for %q, got %v", short, err)
			}
		}
	})

	t.Run("missing rename origin rejected", func(t *testing.T) {
		for _, invalid := range [][]byte{
			[]byte("R  destination.txt\x00"),
			[]byte("R  destination.txt\x00\x00"),
			[]byte("C  destination.txt\x00"),
			[]byte("C  destination.txt\x00\x00"),
		} {
			_, err := parseInventoryBytes(invalid)
			if err == nil || !strings.Contains(err.Error(), "rename origin missing") {
				t.Errorf("expected rename origin missing error for %q, got %v", invalid, err)
			}
		}
	})

	t.Run("record ordering preserved", func(t *testing.T) {
		statusBytes := append(
			append(porcelainStatusBytes(" M", "first.txt", ""), porcelainStatusBytes(" D", "second.txt", "")...),
			porcelainStatusBytes("??", "third.txt", "")...,
		)
		rows, err := parseInventoryBytes(statusBytes)
		if err != nil {
			t.Fatal(err)
		}
		want := []inventoryRow{
			{Classification: "M", Path: "first.txt", Origin: ""},
			{Classification: "D", Path: "second.txt", Origin: ""},
			{Classification: "A", Path: "third.txt", Origin: ""},
		}
		if !reflect.DeepEqual(rows, want) {
			t.Fatalf("rows=%+v, want %+v", rows, want)
		}
	})

	t.Run("metadata exclusions drop only untracked hidden files under do-work", func(t *testing.T) {
		statusBytes := append(
			append(porcelainStatusBytes("??", "do-work/.DS_Store", ""), porcelainStatusBytes("??", "do-work/queue/.REQ.swp", "")...),
			append(porcelainStatusBytes("??", "do-work/.req-reservations/042", ""), porcelainStatusBytes("??", ".DS_Store", "")...)...,
		)
		rows, err := parseInventoryBytes(statusBytes)
		if err != nil {
			t.Fatal(err)
		}
		want := []inventoryRow{
			{Classification: "A", Path: "do-work/.req-reservations/042", Origin: ""},
			{Classification: "A", Path: ".DS_Store", Origin: ""},
		}
		if !reflect.DeepEqual(rows, want) {
			t.Fatalf("rows=%+v, want %+v", rows, want)
		}
	})

	t.Run("cross-row secret promotion promotes additions only", func(t *testing.T) {
		statusBytes := append(
			append(porcelainStatusBytes("??", "first.txt", ""), porcelainStatusBytes("??", "second.txt", "")...),
			append(porcelainStatusBytes(" M", "modified.txt", ""), porcelainStatusBytes(" M", "secret.key", "")...)...,
		)
		rows, err := parseInventoryBytes(statusBytes)
		if err != nil {
			t.Fatal(err)
		}
		want := []inventoryRow{
			{Classification: "X", Path: "first.txt", Origin: ""},
			{Classification: "X", Path: "second.txt", Origin: ""},
			{Classification: "M", Path: "modified.txt", Origin: ""},
			{Classification: "X", Path: "secret.key", Origin: ""},
		}
		if !reflect.DeepEqual(rows, want) {
			t.Fatalf("rows=%+v, want %+v", rows, want)
		}
	})
}

func expectedOrdinaryInventoryClass(status string) string {
	if strings.Contains(status, "D") {
		return "D"
	}
	if status == "??" || status[0] == 'A' {
		return "A"
	}
	return "M"
}

func porcelainStatusBytes(status, path, origin string) []byte {
	output := append([]byte{}, []byte(status+" "+path)...)
	output = append(output, 0)
	if origin != "" {
		output = append(output, []byte(origin)...)
		output = append(output, 0)
	}
	return output
}

func runRetainedInventory(t *testing.T, repository string) []inventoryRow {
	t.Helper()
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	script := filepath.Clean(filepath.Join(workingDirectory, "..", "..", "..", "checks", "uncommitted-inventory.sh"))
	if _, err := os.Stat(script); err != nil {
		t.Fatalf("retained inventory script: %v", err)
	}
	command := exec.Command(script, repository)
	output, runError := command.CombinedOutput()
	if runError != nil {
		t.Fatalf("retained inventory: %v: %s", runError, output)
	}
	rows := []inventoryRow{}
	for _, line := range strings.Split(strings.TrimSuffix(string(output), "\n"), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			t.Fatalf("malformed retained row %q", line)
		}
		rows = append(rows, inventoryRow{Classification: parts[0], Path: parts[1]})
	}
	return rows
}

func compareInventoryProjection(goRows, retainedRows []inventoryRow) error {
	if len(goRows) != len(retainedRows) {
		return fmt.Errorf("row count Go=%d retained=%d: Go=%+v retained=%+v", len(goRows), len(retainedRows), goRows, retainedRows)
	}
	for index := range goRows {
		if goRows[index].Classification != retainedRows[index].Classification || goRows[index].Path != retainedRows[index].Path {
			return fmt.Errorf("row %d Go=%+v retained=%+v", index, goRows[index], retainedRows[index])
		}
	}
	return nil
}

func assertInventoryFindingProjection(t *testing.T, rows []inventoryRow) {
	t.Helper()
	findings := inventoryFindings(rows)
	if len(findings) != len(rows) {
		t.Fatalf("findings=%d rows=%d", len(findings), len(rows))
	}
	for index, row := range rows {
		wantEvidence := row.Classification
		if row.Origin != "" {
			wantEvidence += " from " + row.Origin
		}
		finding := findings[index]
		if finding.Code != "INVENTORY-"+row.Classification || !reflect.DeepEqual(finding.AffectedPaths, []string{row.Path}) || !reflect.DeepEqual(finding.Evidence, []string{wantEvidence}) {
			t.Errorf("row %+v finding=%+v", row, finding)
		}
	}
}

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

	retained := runRetainedInventory(t, repository)
	wantRetained := map[string]string{".envrc": "XD", "visible-config.txt": "X", "ordinary.txt": "X"}
	for _, row := range retained {
		if expected, ok := wantRetained[row.Path]; ok {
			if row.Classification != expected {
				t.Fatalf("retained %s=%s want %s (rows=%+v)", row.Path, row.Classification, expected, retained)
			}
			delete(wantRetained, row.Path)
		}
	}
	if len(wantRetained) != 0 {
		t.Fatalf("missing retained=%v rows=%+v", wantRetained, retained)
	}

	result := handleInventory(testContext(repository), nil)
	if result.Outcome != resultmodel.OutcomeSuccess || !hasFinding(result, "INVENTORY-XD") || !hasFinding(result, "INVENTORY-X") {
		t.Fatalf("Go result=%#v, want INVENTORY-XD and INVENTORY-X", result)
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

func TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest(t *testing.T) {
	// A checkout beneath a directory named "working" must not turn every archived
	// REQ into an in-flight one. In-flight-ness comes from the root being walked
	// (do-work/working versus do-work/archive), never from the checkout's path.
	// The fixture sits under a do-work/working/ directory so that no substring of
	// the absolute path, however specific, can stand in for the walked root.
	repository := filepath.Join(t.TempDir(), "do-work", "working", "project")
	if err := os.MkdirAll(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	writeMatrixFile(t, repository, "do-work/archive/REQ-905-blocked.md", "---\nid: REQ-905\nstatus: blocked\n---\n\n## Implementation Summary\n- `project.txt` (modified)\n")
	writeMatrixFile(t, repository, "do-work/working/REQ-906-blocked-in-flight.md", "---\nid: REQ-906\nstatus: blocked\n---\n\n## Implementation Summary\n- `in-flight.txt` (modified)\n")
	writeMatrixFile(t, repository, "project.txt", "uncommitted\n")
	writeMatrixFile(t, repository, "in-flight.txt", "uncommitted\n")
	associations, err := AssociateProjectPaths(repository, []string{"project.txt", "in-flight.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if owner := associations["project.txt"]; owner != "" {
		t.Fatalf("blocked archived request claimed project.txt through the checkout path: owner=%q", owner)
	}
	if owner := associations["in-flight.txt"]; owner != "REQ-906" {
		t.Fatalf("working request must claim in-flight.txt whatever its status says: owner=%q", owner)
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

func TestProtectedInventoryPreservesAssociationPartitionByClassification(t *testing.T) {
	t.Run("M A and D remain pathname candidates", func(t *testing.T) {
		repository := newGitFixture(t)
		writeInventoryOwner(t, repository, "first.txt", "second.txt", "added.txt")
		runFixtureGitCommand(t, repository, "add", "do-work")
		runFixtureGitCommand(t, repository, "commit", "-qm", "inventory owner")
		if err := os.WriteFile(filepath.Join(repository, "first.txt"), []byte("modified\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(filepath.Join(repository, "second.txt")); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repository, "added.txt"), []byte("added\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		start := handleProtectedInventory(testContext(repository), []string{"start"})
		for _, code := range []string{"INVENTORY-M", "INVENTORY-A", "INVENTORY-D"} {
			if !hasFinding(start, code) {
				t.Fatalf("start omitted %s: %#v", code, start)
			}
		}
		associate := handleProtectedInventory(testContext(repository), []string{"associate"})
		assertAssociatedPaths(t, associate, "first.txt", "second.txt", "added.txt")
	})

	t.Run("XD remains a pathname candidate while X is quarantined", func(t *testing.T) {
		repository := newGitFixture(t)
		for _, path := range []string{"secret.pem", "ordinary.txt"} {
			if err := os.WriteFile(filepath.Join(repository, path), []byte("fixture\n"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		writeInventoryOwner(t, repository, "secret.pem", "ordinary.txt", ".env.local")
		runFixtureGitCommand(t, repository, "add", ".")
		runFixtureGitCommand(t, repository, "commit", "-qm", "inventory owner and paths")
		if err := os.Remove(filepath.Join(repository, "secret.pem")); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repository, "ordinary.txt"), []byte("modified\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repository, ".env.local"), []byte("not-read\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		start := handleProtectedInventory(testContext(repository), []string{"start"})
		if !hasFinding(start, "INVENTORY-XD") || !hasFinding(start, "INVENTORY-X") {
			t.Fatalf("start classes=%#v", start)
		}
		associate := handleProtectedInventory(testContext(repository), []string{"associate"})
		assertAssociatedPaths(t, associate, "secret.pem", "ordinary.txt")
		for _, finding := range associate.Findings {
			if reflect.DeepEqual(finding.AffectedPaths, []string{".env.local"}) {
				t.Fatalf("X path escaped quarantine: %#v", associate)
			}
		}
	})
}

func writeInventoryOwner(t *testing.T, repository string, paths ...string) {
	t.Helper()
	requestPath := filepath.Join(repository, "do-work", "working", "REQ-999-inventory.md")
	if err := os.MkdirAll(filepath.Dir(requestPath), 0o755); err != nil {
		t.Fatal(err)
	}
	var contents strings.Builder
	contents.WriteString("---\nid: REQ-999\nstatus: claimed\n---\n\n## Implementation Summary\n")
	for _, path := range paths {
		fmt.Fprintf(&contents, "- `%s` (modified)\n", path)
	}
	if err := os.WriteFile(requestPath, []byte(contents.String()), 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertAssociatedPaths(t *testing.T, result resultmodel.CommandResult, paths ...string) {
	t.Helper()
	want := make(map[string]bool, len(paths))
	for _, path := range paths {
		want[path] = true
	}
	for _, finding := range result.Findings {
		if finding.Code == "ASSOCIATION-FOUND" && len(finding.AffectedPaths) == 1 {
			delete(want, finding.AffectedPaths[0])
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing associated paths %v: %#v", want, result)
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

func TestProtectedInventoryCompatibilityShimPreservesErrors(t *testing.T) {
	t.Setenv("DO_WORK_COMPATIBILITY_SHIM", "1")

	t.Run("preserves NO-DO-WORK-DIR", func(t *testing.T) {
		repository := newGitFixture(t)
		if err := os.WriteFile(filepath.Join(repository, "file.txt"), []byte("changed\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		start := handleProtectedInventory(testContext(repository), []string{"start"})
		if start.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("start failed: %#v", start)
		}
		associate := handleProtectedInventory(testContext(repository), []string{"associate"})
		rendered, err := resultmodel.RenderResult(associate, resultmodel.FormatText)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(rendered), "NO-DO-WORK-DIR") {
			t.Fatalf("expected NO-DO-WORK-DIR in rendered output, got: %q", string(rendered))
		}
	})

	t.Run("preserves PARSE-FAILED", func(t *testing.T) {
		repository := newGitFixture(t)
		if err := os.WriteFile(filepath.Join(repository, "file.txt"), []byte("changed\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		start := handleProtectedInventory(testContext(repository), []string{"start"})
		if start.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("start failed: %#v", start)
		}
		workingReq := filepath.Join(repository, "do-work", "working", "REQ-999.md")
		if err := os.MkdirAll(filepath.Dir(workingReq), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(workingReq, []byte("---\nid: REQ-999\nstatus: claimed\n---\n\n## Implementation Summary\n- `unmatched\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		associate := handleProtectedInventory(testContext(repository), []string{"associate"})
		rendered, err := resultmodel.RenderResult(associate, resultmodel.FormatText)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(rendered), "PARSE-FAILED") {
			t.Fatalf("expected PARSE-FAILED in rendered output, got: %q", string(rendered))
		}
	})

	t.Run("preserves walk error HELPER-USAGE finding", func(t *testing.T) {
		repository := newGitFixture(t)
		if err := os.WriteFile(filepath.Join(repository, "file.txt"), []byte("changed\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		start := handleProtectedInventory(testContext(repository), []string{"start"})
		if start.Outcome != resultmodel.OutcomeSuccess {
			t.Fatalf("start failed: %#v", start)
		}
		workingReq := filepath.Join(repository, "do-work", "working", "REQ-999.md")
		if err := os.MkdirAll(filepath.Dir(workingReq), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(workingReq, []byte("unreadable\n"), 0o000); err != nil {
			t.Fatal(err)
		}
		defer os.Chmod(workingReq, 0o644)
		associate := handleProtectedInventory(testContext(repository), []string{"associate"})
		rendered, err := resultmodel.RenderResult(associate, resultmodel.FormatText)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(rendered), "HELPER-USAGE") {
			t.Fatalf("expected HELPER-USAGE in rendered output, got: %q", string(rendered))
		}
	})
}
