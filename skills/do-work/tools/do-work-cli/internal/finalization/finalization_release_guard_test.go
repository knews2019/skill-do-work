package finalization

import (
	"strings"
	"testing"
)

// seedMaintainerSuite makes the fixture repository a maintainer checkout: a module
// declaration and the module's own VERSION, which is what makes a root shipped.
func seedMaintainerSuite(t *testing.T, repositoryRoot string) {
	t.Helper()
	writeFinalizationFile(t, repositoryRoot, "suite/modules.tsv", "source\tdestination\nskills/do-work\t.claude/skills/do-work\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/VERSION", "1.0.0\n")
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/CHANGELOG.md", "# Changelog\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed maintainer suite")
}

func commitTouching(t *testing.T, repositoryRoot string, paths ...string) string {
	t.Helper()
	for _, path := range paths {
		content := "changed by " + path + "\n"
		if path == "suite/modules.tsv" {
			// The declaration must stay readable: an unreadable one fails closed as
			// RELEASE-SHIPPED-CHANGE-UNVERIFIABLE, which is not what this case tests.
			content = "source\tdestination\nskills/do-work\t.claude/skills/do-work\nskills/do-work-extra\t.claude/skills/do-work-extra\n"
		}
		writeFinalizationFile(t, repositoryRoot, path, content)
	}
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "implementation")
	return strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))
}

// The failure this pins: four versions (0.305.17, .21, .22, .23) were released for
// commits that changed only _dev/tests, and nothing refused them.
func TestReleaseRefusedWhenTheImplementationShipsNothing(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedMaintainerSuite(t, repositoryRoot)

	cases := []struct {
		name    string
		paths   []string
		refused bool
	}{
		{"maintainer-only test change", []string{"_dev/tests/some-probe.sh"}, true},
		{"release metadata alone is not a shipped change", []string{"skills/do-work/VERSION", "skills/do-work/CHANGELOG.md"}, true},
		{"a shipped action file", []string{"_dev/tests/some-probe.sh", "skills/do-work/actions/work.md"}, false},
		{"the suite declaration", []string{"suite/modules.tsv"}, false},
		{"the installer tools", []string{"tools/install-do-work-suite.sh"}, false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			hash := commitTouching(t, repositoryRoot, testCase.paths...)
			err := releaseShippedChangeError(repositoryRoot, Manifest{ProvenanceMode: ProvenanceSuppliedCommit, ImplementationHash: hash})
			if testCase.refused && (err == nil || !strings.Contains(err.Error(), "RELEASE-WITHOUT-SHIPPED-CHANGE")) {
				t.Fatalf("expected RELEASE-WITHOUT-SHIPPED-CHANGE, got %v", err)
			}
			if !testCase.refused && err != nil {
				t.Fatalf("shipped change refused: %v", err)
			}
		})
	}
}

// A merge is judged by what it brought in (its first-parent diff), not by the union
// across both parents, which would credit the merge with main's own shipped changes.
func TestReleaseGuardReadsAMergeByItsFirstParent(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedMaintainerSuite(t, repositoryRoot)
	base := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))
	commitTouching(t, repositoryRoot, "skills/do-work/actions/shipped-on-main.md")
	runFinalizationGit(t, repositoryRoot, "checkout", "-q", "-b", "maintainer-only", base)
	commitTouching(t, repositoryRoot, "_dev/tests/branch-probe.sh")
	runFinalizationGit(t, repositoryRoot, "checkout", "-q", "-")
	runFinalizationGit(t, repositoryRoot, "merge", "-q", "--no-ff", "-m", "merge maintainer-only", "maintainer-only")
	merge := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))

	err := releaseShippedChangeError(repositoryRoot, Manifest{ProvenanceMode: ProvenanceSuppliedCommit, ImplementationHash: merge})
	if err == nil || !strings.Contains(err.Error(), "RELEASE-WITHOUT-SHIPPED-CHANGE") {
		t.Fatalf("a merge bringing in only _dev/tests was accepted as a release: %v", err)
	}
}

func TestReleaseGuardUsesCommitPathsUnderPrimaryCommitProvenance(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	seedMaintainerSuite(t, repositoryRoot)
	refused := releaseShippedChangeError(repositoryRoot, Manifest{ProvenanceMode: ProvenancePrimaryCommit, CommitPaths: []string{"do-work/archive/REQ-1.md", "skills/do-work/VERSION"}})
	if refused == nil || !strings.Contains(refused.Error(), "RELEASE-WITHOUT-SHIPPED-CHANGE") {
		t.Fatalf("record plus metadata accepted as a release: %v", refused)
	}
	writeFinalizationFile(t, repositoryRoot, "skills/do-work/actions/work.md", "shipped implementation\n")
	if err := releaseShippedChangeError(repositoryRoot, Manifest{ProvenanceMode: ProvenancePrimaryCommit, CommitPaths: []string{"do-work/archive/REQ-1.md", "skills/do-work/actions/work.md"}}); err != nil {
		t.Fatalf("shipped change refused: %v", err)
	}
}

// A consumer project declares no modules; its releases are its own package's and the
// guard must stay out of the way.
func TestReleaseGuardIgnoresARepositoryWithoutAModuleDeclaration(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	hash := commitTouching(t, repositoryRoot, "src/app.js")
	if err := releaseShippedChangeError(repositoryRoot, Manifest{ProvenanceMode: ProvenanceSuppliedCommit, ImplementationHash: hash}); err != nil {
		t.Fatalf("consumer release refused: %v", err)
	}
}

// Manifests carry runtime configuration as well as a release version. The guard
// must judge the actual edit for both committed and pending implementations.
func TestReleaseGuardDistinguishesFunctionalManifestEdits(t *testing.T) {
	for _, mode := range []string{ProvenanceSuppliedCommit, ProvenancePrimaryCommit} {
		for _, fixture := range []struct{ name, before, functional, metadata string }{
			{"package.json", `{"name":"demo","version":"1.0.0","dependencies":{"dep":"1.0.0"},"main":"old.js"}`, `{"name":"demo","version":"1.0.1","dependencies":{"dep":"2.0.0"},"main":"old.js"}`, `{"name":"demo","version":"1.0.1","dependencies":{"dep":"1.0.0"},"main":"old.js"}`},
			{"Cargo.toml", "[package]\nname = \"demo\"\nversion = \"1.0.0\"\n[dependencies]\ndep = \"1.0.0\"\n", "[package]\nname = \"demo\"\nversion = \"1.0.1\"\n[dependencies]\ndep = \"2.0.0\"\n", "[package]\nname = \"demo\"\nversion = \"1.0.1\"\n[dependencies]\ndep = \"1.0.0\"\n"},
			{"pyproject.toml", "[project]\nname = \"demo\"\nversion = \"1.0.0\"\ndependencies = [\"dep==1.0.0\"]\n", "[project]\nname = \"demo\"\nversion = \"1.0.1\"\ndependencies = [\"dep==2.0.0\"]\n", "[project]\nname = \"demo\"\nversion = \"1.0.1\"\ndependencies = [\"dep==1.0.0\"]\n"},
		} {
			configuration := strings.Replace(fixture.before, `"main":"old.js"`, `"main":"new.js"`, 1)
			if fixture.name == "Cargo.toml" {
				configuration = strings.Replace(fixture.before, "[package]", "[package]\nbuild = \"build.rs\"", 1)
			} else if fixture.name == "pyproject.toml" {
				configuration = fixture.before + "[build-system]\nrequires = [\"setuptools\"]\nbuild-backend = \"setuptools.build_meta\"\n"
			}
			for _, edit := range []struct {
				name, after string
				functional  bool
			}{
				{"dependencies", strings.Replace(fixture.functional, "1.0.1", "1.0.0", 1), true},
				{"configuration", configuration, true},
				{"version-only", fixture.metadata, false},
			} {
				t.Run(mode+"/"+fixture.name+"/"+edit.name, func(t *testing.T) {
					repositoryRoot := newFinalizationRepository(t)
					seedMaintainerSuite(t, repositoryRoot)
					path := "skills/do-work/" + fixture.name
					writeFinalizationFile(t, repositoryRoot, path, fixture.before)
					runFinalizationGit(t, repositoryRoot, "add", path)
					runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed manifest")
					writeFinalizationFile(t, repositoryRoot, path, edit.after)
					manifest := Manifest{ProvenanceMode: mode, CommitPaths: []string{path}}
					if mode == ProvenanceSuppliedCommit {
						runFinalizationGit(t, repositoryRoot, "add", path)
						runFinalizationGit(t, repositoryRoot, "commit", "-qm", "edit manifest")
						manifest.ImplementationHash = strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))
						// Later worktree content must not change the named commit's verdict.
						writeFinalizationFile(t, repositoryRoot, path, fixture.before)
					}
					err := releaseShippedChangeError(repositoryRoot, manifest)
					if edit.functional && err != nil {
						t.Fatalf("functional manifest edit refused: %v", err)
					}
					if !edit.functional && (err == nil || !strings.Contains(err.Error(), "RELEASE-WITHOUT-SHIPPED-CHANGE")) {
						t.Fatalf("version-only edit: expected release refusal, got %v", err)
					}
				})
			}
		}
	}
}
