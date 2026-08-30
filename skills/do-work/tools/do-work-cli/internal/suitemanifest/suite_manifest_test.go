package suitemanifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureVersion = "0.184.0"

// newValidSuiteFixture builds the smallest tree ValidateSuite accepts, so each rejection
// case below can break exactly one thing about it.
func newValidSuiteFixture(t *testing.T) string {
	t.Helper()
	archiveRoot := t.TempDir()
	writeFixtureFile(t, filepath.Join(archiveRoot, "VERSION"), fixtureVersion+"\n")
	writeFixtureFile(t, filepath.Join(archiveRoot, "suite", "modules.tsv"),
		"source\tdestination\n"+
			"skills/do-work\t.claude/skills/do-work\n"+
			"skills/do-work-board\t.claude/skills/do-work-board\n"+
			"skills/do-work-knowledge\t.claude/skills/do-work-knowledge\n"+
			"skills/do-work-toolbox\t.claude/skills/do-work-toolbox\n")
	for _, module := range []string{"do-work", "do-work-board", "do-work-knowledge", "do-work-toolbox"} {
		writeFixtureFile(t, filepath.Join(archiveRoot, "skills", module, "SKILL.md"), "# "+module+"\n")
	}
	writeFixtureFile(t, filepath.Join(archiveRoot, "skills", "do-work", "VERSION"), fixtureVersion+"\n")
	writeFixtureFile(t, filepath.Join(archiveRoot, "skills", "do-work", "actions", "version.md"),
		"# Version Action\n\n**Current version**: "+fixtureVersion+"\n")
	return archiveRoot
}

func writeFixtureFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestCanonicalFixtureValidatesAndReportsItsFourModules(t *testing.T) {
	result, err := ValidateSuite(newValidSuiteFixture(t))
	if err != nil {
		t.Fatalf("ValidateSuite: %v", err)
	}
	if result.SuiteVersion != fixtureVersion {
		t.Errorf("suite version = %q, want %q", result.SuiteVersion, fixtureVersion)
	}
	if len(result.Modules) != 4 {
		t.Errorf("module count = %d, want 4", len(result.Modules))
	}
	expectedSummary := "suite manifest valid: v" + fixtureVersion + " (4 modules)"
	if result.SummaryLine() != expectedSummary {
		t.Errorf("summary = %q, want %q", result.SummaryLine(), expectedSummary)
	}
}

// Every rejection below is one suite-manifest-contract.sh drives end to end today. Naming the
// expected phrase keeps the Go port's diagnostics identical to the shell validator's.
func TestMalformedSuitesAreRejectedWithTheirCurrentDiagnostic(t *testing.T) {
	tests := []struct {
		name             string
		corrupt          func(t *testing.T, archiveRoot string)
		expectedFragment string
	}{
		{
			name: "an extra manifest column",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "suite", "modules.tsv"), "source\tdestination\textra\n")
			},
			expectedFragment: "manifest header must be exactly: source<TAB>destination",
		},
		{
			name: "an absolute source",
			corrupt: func(t *testing.T, archiveRoot string) {
				replaceManifestRow(t, archiveRoot, "skills/do-work\t", "/skills/do-work\t")
			},
			expectedFragment: "source must be relative",
		},
		{
			name: "a traversing source",
			corrupt: func(t *testing.T, archiveRoot string) {
				replaceManifestRow(t, archiveRoot, "skills/do-work\t", "skills/../skills/do-work\t")
			},
			expectedFragment: "source traverses directories",
		},
		{
			name: "a traversing destination",
			corrupt: func(t *testing.T, archiveRoot string) {
				replaceManifestRow(t, archiveRoot, "\t.claude/skills/do-work\n", "\t.claude/../skills/do-work\n")
			},
			expectedFragment: "destination traverses directories",
		},
		{
			name: "a carriage return in a row",
			corrupt: func(t *testing.T, archiveRoot string) {
				replaceManifestRow(t, archiveRoot, ".claude/skills/do-work\n", ".claude/skills/do-work\r\n")
			},
			expectedFragment: "contains a carriage return",
		},
		{
			name: "a module mapped outside its required destination",
			corrupt: func(t *testing.T, archiveRoot string) {
				replaceManifestRow(t, archiveRoot, "\t.claude/skills/do-work\n", "\t.claude/skills/elsewhere\n")
			},
			expectedFragment: "outside its required destination",
		},
		{
			name: "a duplicated source row",
			corrupt: func(t *testing.T, archiveRoot string) {
				appendManifestRow(t, archiveRoot, "skills/do-work\t.claude/skills/duplicate\n")
			},
			expectedFragment: "duplicates source skills/do-work",
		},
		{
			name: "a missing module directory",
			corrupt: func(t *testing.T, archiveRoot string) {
				if err := os.RemoveAll(filepath.Join(archiveRoot, "skills", "do-work-board")); err != nil {
					t.Fatalf("remove module: %v", err)
				}
			},
			expectedFragment: "must be a real directory, not a missing path or symlink",
		},
		{
			name: "an empty SKILL.md",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "skills", "do-work-toolbox", "SKILL.md"), "")
			},
			expectedFragment: "SKILL.md must be a non-empty regular file",
		},
		{
			name: "a symlinked module root",
			corrupt: func(t *testing.T, archiveRoot string) {
				modulePath := filepath.Join(archiveRoot, "skills", "do-work-knowledge")
				if err := os.RemoveAll(modulePath); err != nil {
					t.Fatalf("remove module: %v", err)
				}
				if err := os.Symlink(filepath.Join(archiveRoot, "skills", "do-work"), modulePath); err != nil {
					t.Fatalf("symlink module: %v", err)
				}
			},
			expectedFragment: "must be a real directory, not a missing path or symlink",
		},
		{
			name: "a non-semantic VERSION",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "VERSION"), "not-a-version\n")
			},
			expectedFragment: "VERSION must be a plain semantic version (X.Y.Z)",
		},
		{
			name: "a VERSION carrying a second line",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "VERSION"), fixtureVersion+"\nextra\n")
			},
			expectedFragment: "VERSION must contain exactly one newline-terminated line",
		},
		{
			name: "a core VERSION that disagrees with the suite VERSION",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "skills", "do-work", "VERSION"), "9.9.9\n")
			},
			expectedFragment: "skills/do-work/VERSION mismatch",
		},
		{
			name: "two Current version lines",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "skills", "do-work", "actions", "version.md"),
					"**Current version**: "+fixtureVersion+"\n**Current version**: "+fixtureVersion+"\n")
			},
			expectedFragment: "must contain exactly one Current version line",
		},
		{
			name: "an action version that disagrees with the suite VERSION",
			corrupt: func(t *testing.T, archiveRoot string) {
				writeFixtureFile(t, filepath.Join(archiveRoot, "skills", "do-work", "actions", "version.md"),
					"**Current version**: 9.9.9\n")
			},
			expectedFragment: "version.md mismatch",
		},
		{
			name: "a missing manifest",
			corrupt: func(t *testing.T, archiveRoot string) {
				if err := os.Remove(filepath.Join(archiveRoot, "suite", "modules.tsv")); err != nil {
					t.Fatalf("remove manifest: %v", err)
				}
			},
			expectedFragment: "suite/modules.tsv must be a regular file in the suite root",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			archiveRoot := newValidSuiteFixture(t)
			test.corrupt(t, archiveRoot)
			_, err := ValidateSuite(archiveRoot)
			if err == nil {
				t.Fatalf("a malformed suite was accepted")
			}
			if !strings.Contains(err.Error(), test.expectedFragment) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), test.expectedFragment)
			}
		})
	}
}

// The validator answers a question; it must never edit the tree it is asked about.
func TestValidationDoesNotModifyTheArchive(t *testing.T) {
	archiveRoot := newValidSuiteFixture(t)
	before := treeFingerprint(t, archiveRoot)
	if _, err := ValidateSuite(archiveRoot); err != nil {
		t.Fatalf("ValidateSuite: %v", err)
	}
	if after := treeFingerprint(t, archiveRoot); after != before {
		t.Errorf("validation modified the archive:\nbefore %s\nafter  %s", before, after)
	}
}

// The two Current-version parsers in the shell suite disagreed deliberately; the Go port
// keeps the validator's stricter one, which takes the whole remainder of the line.
func TestReadActionVersionTakesTheWholeRemainderAndCountsEveryMarker(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "version.md")
	writeFixtureFile(t, path, "**Current version**: 1.2.3-rc1 trailing\n**Current version**: 9.9.9\n")
	version, markerCount, err := ReadActionVersion(path)
	if err != nil {
		t.Fatalf("ReadActionVersion: %v", err)
	}
	if version != "1.2.3-rc1 trailing" {
		t.Errorf("version = %q, want the whole remainder of the first marker line", version)
	}
	if markerCount != 2 {
		t.Errorf("marker count = %d, want 2", markerCount)
	}
}

func replaceManifestRow(t *testing.T, archiveRoot, oldText, newText string) {
	t.Helper()
	manifestPath := filepath.Join(archiveRoot, "suite", "modules.tsv")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	updated := strings.Replace(string(data), oldText, newText, 1)
	if updated == string(data) {
		t.Fatalf("manifest rewrite %q -> %q matched nothing", oldText, newText)
	}
	writeFixtureFile(t, manifestPath, updated)
}

func appendManifestRow(t *testing.T, archiveRoot, row string) {
	t.Helper()
	manifestPath := filepath.Join(archiveRoot, "suite", "modules.tsv")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	writeFixtureFile(t, manifestPath, string(data)+row)
}

func treeFingerprint(t *testing.T, root string) string {
	t.Helper()
	var fingerprint strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fingerprint.WriteString(path)
		fingerprint.WriteString(info.Mode().String())
		if info.Mode().IsRegular() {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			fingerprint.Write(data)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return fingerprint.String()
}
