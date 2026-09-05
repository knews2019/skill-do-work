package finalization

import (
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// TestREQ565AmbiguousWorkspaceSourceIdentityFailsClosed pins the first gap left
// open by REQ-512: a changed Cargo or uv source whose applicable identity is
// duplicated or competing used to release under whichever name was listed
// first, so the mirror belonging to the other declared name silently stayed on
// the old version. Every fixture below updates only the first-listed name's
// lock entry, which is exactly what the first-match parser accepted.
func TestREQ565AmbiguousWorkspaceSourceIdentityFailsClosed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		manifestPath   string
		lockPath       string
		manifestBefore string
		manifestAfter  string
		lockBefore     string
		lockAfter      string
	}{
		{
			name:           "cargo declares the package name twice",
			manifestPath:   "skills/do-work/ambiguous-rust/Cargo.toml",
			lockPath:       "skills/do-work/ambiguous-rust/Cargo.lock",
			manifestBefore: "[package]\nname = \"alpha\"\nname = \"beta\"\nversion = \"1.0.0\"\n",
			manifestAfter:  "[package]\nname = \"alpha\"\nname = \"beta\"\nversion = \"1.0.1\"\n",
			lockBefore:     "[[package]]\nname = \"alpha\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\n",
			lockAfter:      "[[package]]\nname = \"alpha\"\nversion = \"1.0.1\"\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\n",
		},
		{
			name:           "cargo repeats the package section",
			manifestPath:   "skills/do-work/repeated-rust/Cargo.toml",
			lockPath:       "skills/do-work/repeated-rust/Cargo.lock",
			manifestBefore: "[package]\nname = \"alpha\"\nversion = \"1.0.0\"\n\n[dependencies]\nserde = \"1\"\n\n[package]\nname = \"beta\"\n",
			manifestAfter:  "[package]\nname = \"alpha\"\nversion = \"1.0.1\"\n\n[dependencies]\nserde = \"1\"\n\n[package]\nname = \"beta\"\n",
			lockBefore:     "[[package]]\nname = \"alpha\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\n",
			lockAfter:      "[[package]]\nname = \"alpha\"\nversion = \"1.0.1\"\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\n",
		},
		{
			name:           "uv declares the project name twice",
			manifestPath:   "skills/do-work/ambiguous-python/pyproject.toml",
			lockPath:       "skills/do-work/ambiguous-python/uv.lock",
			manifestBefore: "[project]\nname = \"alpha\"\nname = \"beta\"\nversion = \"1.0.0\"\n",
			manifestAfter:  "[project]\nname = \"alpha\"\nname = \"beta\"\nversion = \"1.0.1\"\n",
			lockBefore:     "[[package]]\nname = \"alpha\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n",
			lockAfter:      "[[package]]\nname = \"alpha\"\nversion = \"1.0.1\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n",
		},
		{
			name:           "uv project and poetry sections compete",
			manifestPath:   "skills/do-work/competing-python/pyproject.toml",
			lockPath:       "skills/do-work/competing-python/uv.lock",
			manifestBefore: "[project]\nname = \"alpha\"\nversion = \"1.0.0\"\n\n[tool.poetry]\nname = \"beta\"\n",
			manifestAfter:  "[project]\nname = \"alpha\"\nversion = \"1.0.1\"\n\n[tool.poetry]\nname = \"beta\"\n",
			lockBefore:     "[[package]]\nname = \"alpha\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n",
			lockAfter:      "[[package]]\nname = \"alpha\"\nversion = \"1.0.1\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, test.manifestPath, test.manifestBefore)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockBefore)
			seedSemanticLegacyTail(t, repositoryRoot)
			writeFinalizationFile(t, repositoryRoot, test.manifestPath, test.manifestAfter)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockAfter)

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			assertReleaseEnumerationRefusal(t, result)
			if got := readFinalizationFile(t, repositoryRoot, test.manifestPath); got != test.manifestAfter {
				t.Fatalf("refusal changed source bytes: %q", got)
			}
			if got := readFinalizationFile(t, repositoryRoot, test.lockPath); got != test.lockAfter {
				t.Fatalf("refusal changed lock bytes: %q", got)
			}
		})
	}
}

// TestREQ565ChangedNPMRootRequiresEveryPresentRootVersionCopy pins the second
// gap left open by REQ-512: root copy counting keyed on "already equals the old
// version", so a structurally present copy holding any other value dropped out
// of the required mirror set instead of refusing. Presence is the obligation.
func TestREQ565ChangedNPMRootRequiresEveryPresentRootVersionCopy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		lockBefore string
		lockAfter  string
		wantRefuse bool
	}{
		{
			name:       "top-level root copy is stale",
			lockBefore: "{\"name\":\"consumer\",\"version\":\"0.9.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"1.0.0\"}}}\n",
			lockAfter:  "{\"name\":\"consumer\",\"version\":\"0.9.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"1.0.1\"}}}\n",
			wantRefuse: true,
		},
		{
			name:       "both root copies are stale",
			lockBefore: "{\"name\":\"consumer\",\"version\":\"0.9.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"0.9.0\"}}}\n",
			lockAfter:  "{\"name\":\"consumer\",\"version\":\"0.9.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"0.9.0\"}}}\n",
			wantRefuse: true,
		},
		{
			name:       "only the top-level root copy is present",
			lockBefore: "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"node_modules/dependency\":{\"name\":\"dependency\",\"version\":\"2.0.0\"}}}\n",
			lockAfter:  "{\"name\":\"consumer\",\"version\":\"1.0.1\",\"lockfileVersion\":3,\"packages\":{\"node_modules/dependency\":{\"name\":\"dependency\",\"version\":\"2.0.0\"}}}\n",
		},
	}
	rootPath := "skills/do-work/root-npm/package.json"
	lockPath := "skills/do-work/root-npm/package-lock.json"
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, rootPath, "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"private\":true}\n")
			writeFinalizationFile(t, repositoryRoot, lockPath, test.lockBefore)
			seedSemanticLegacyTail(t, repositoryRoot)
			writeFinalizationFile(t, repositoryRoot, rootPath, "{\"name\":\"consumer\",\"version\":\"1.0.1\",\"private\":true}\n")
			writeFinalizationFile(t, repositoryRoot, lockPath, test.lockAfter)

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			if test.wantRefuse {
				assertReleaseEnumerationRefusal(t, result)
				if got := readFinalizationFile(t, repositoryRoot, lockPath); got != test.lockAfter {
					t.Fatalf("refusal changed lock bytes: %q", got)
				}
				return
			}
			if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
				t.Fatalf("single present root copy recovery = %#v", result)
			}
			for _, path := range []string{rootPath, lockPath} {
				if !containsFinalizationPath(result.Finalization.CommitPaths, path) {
					t.Fatalf("commit paths %v omit changed root path %s", result.Finalization.CommitPaths, path)
				}
			}
		})
	}
}

func assertReleaseEnumerationRefusal(t *testing.T, result resultmodel.CommandResult) {
	t.Helper()
	if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || len(result.Finalization.ReasonCodes) != 1 || result.Finalization.ReasonCodes[0] != "FINALIZATION-DISCOVERY-RELEASE-ENUMERATION" {
		t.Fatalf("release identity gap did not fail closed: %#v", result)
	}
}
