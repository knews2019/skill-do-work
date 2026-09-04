package finalization

import (
	"fmt"
	"os"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestREQ512TrackedFoldRequiresClosedBoundary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		append string
		want   bool
	}{
		{
			name:   "bounded review fold",
			append: "\n## Review Fold — REQ-812\n\nExact originating finding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n",
			want:   true,
		},
		{
			name:   "bounded recovery fold",
			append: "\n## Recovery Fold — REQ-812\n\nExact originating recovery.\n\n<!-- do-work:finalization-followup-fold-end kind=recovery request=REQ-812 -->\n",
			want:   true,
		},
		{
			name:   "unheaded foreign tail has no owned boundary",
			append: "\n## Review Fold — REQ-812\n\nExact originating finding.\n\nUnrelated unheaded paragraph.\n",
		},
		{name: "foreign bytes before", append: "\nforeign preface\n\n## Review Fold — REQ-812\n\nFinding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n"},
		{name: "unheaded bytes after", append: "\n## Review Fold — REQ-812\n\nFinding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\nforeign tail\n"},
		{name: "comment after", append: "\n## Review Fold — REQ-812\n\nFinding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n<!-- foreign -->\n"},
		{name: "malformed heading after", append: "\n## Review Fold — REQ-812\n\nFinding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n### Foreign\n"},
		{name: "delimiter after", append: "\n## Review Fold — REQ-812\n\nFinding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n---\n"},
		{name: "mismatched marker", append: "\n## Review Fold — REQ-812\n\nFinding.\n\n<!-- do-work:finalization-followup-fold-end kind=recovery request=REQ-812 -->\n"},
		{name: "duplicate marker", append: "\n## Review Fold — REQ-812\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n"},
		{name: "inline marker", append: "\n## Review Fold — REQ-812\n\nFinding.<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n"},
		{name: "second fold", append: "\n## Review Fold — REQ-812\n\nFinding.\n\n## Recovery Fold — REQ-812\n\nRecovery.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-812 -->\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			path := "do-work/queue/REQ-813-follow-up.md"
			before := "---\nid: REQ-813\ntitle: Follow-up\nstatus: pending\naddendum_to: REQ-812\n---\n\nOriginal body.\n"
			writeFinalizationFile(t, repositoryRoot, path, before)
			runFinalizationGit(t, repositoryRoot, "add", ".")
			runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
			writeFinalizationFile(t, repositoryRoot, path, before+test.append)
			if got := followupPathProves(repositoryRoot, path, "REQ-812"); got != test.want {
				t.Fatalf("followupPathProves() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestREQ512MalformedWorkspaceIdentityAndNPMRootLockFailClosed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		manifestPath string
		lockPath     string
		before       string
		after        string
		lock         string
	}{
		{
			name:         "cargo source missing package name",
			manifestPath: "skills/do-work/broken-rust/Cargo.toml",
			lockPath:     "skills/do-work/broken-rust/Cargo.lock",
			before:       "[package]\nversion = \"1.0.0\"\n",
			after:        "[package]\nversion = \"1.0.1\"\n",
			lock:         "[[package]]\nname = \"broken\"\nversion = \"1.0.0\"\n",
		},
		{
			name:         "uv source has malformed project name",
			manifestPath: "skills/do-work/broken-python/pyproject.toml",
			lockPath:     "skills/do-work/broken-python/uv.lock",
			before:       "[project]\nname = broken\nversion = \"1.0.0\"\n",
			after:        "[project]\nname = broken\nversion = \"1.0.1\"\n",
			lock:         "[[package]]\nname = \"broken\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n",
		},
		{
			name:         "npm root lock is malformed json",
			manifestPath: "skills/do-work/broken-npm/package.json",
			lockPath:     "skills/do-work/broken-npm/package-lock.json",
			before:       "{\"name\":\"broken\",\"version\":\"1.0.0\"}\n",
			after:        "{\"name\":\"broken\",\"version\":\"1.0.1\"}\n",
			lock:         "{\"name\":\"broken\",\"version\":\"1.0.0\",\"packages\":",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, test.manifestPath, test.before)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lock)
			seedSemanticLegacyTail(t, repositoryRoot)
			writeFinalizationFile(t, repositoryRoot, test.manifestPath, test.after)

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || len(result.Finalization.ReasonCodes) != 1 || result.Finalization.ReasonCodes[0] != "FINALIZATION-DISCOVERY-RELEASE-ENUMERATION" {
				t.Fatalf("malformed release authority did not fail closed: %#v", result)
			}
			if got := readFinalizationFile(t, repositoryRoot, test.manifestPath); got != test.after {
				t.Fatalf("refusal changed source bytes: %q", got)
			}
			if got := readFinalizationFile(t, repositoryRoot, test.lockPath); got != test.lock {
				t.Fatalf("refusal changed lock bytes: %q", got)
			}
		})
	}
}

func TestREQ512TrackedFoldForeignTailRefusesWithoutMutation(t *testing.T) {
	t.Parallel()
	repositoryRoot := newFinalizationRepository(t)
	archivePath := "do-work/archive/REQ-814-fixture.md"
	followupPath := "do-work/queue/REQ-815-follow-up.md"
	followupBefore := "---\nid: REQ-815\ntitle: Follow-up\nstatus: pending\naddendum_to: REQ-814\n---\n\nOriginal body.\n"
	writeFinalizationFile(t, repositoryRoot, archivePath, "---\nid: REQ-814\ntitle: Fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `owned.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, followupPath, followupBefore)
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, followupPath, followupBefore+"\n## Review Fold — REQ-814\n\nExact originating finding.\n\n<!-- do-work:finalization-followup-fold-end kind=review request=REQ-814 -->\nUnrelated unheaded paragraph.\n")
	bytesBefore := readFinalizationFile(t, repositoryRoot, followupPath)
	statusBefore := runFinalizationGit(t, repositoryRoot, "status", "--short")
	headBefore := runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", followupPath)
	if got := readFinalizationFile(t, repositoryRoot, followupPath); got != bytesBefore {
		t.Fatalf("refusal changed follow-up bytes: %q", got)
	}
	if got := runFinalizationGit(t, repositoryRoot, "status", "--short"); got != statusBefore {
		t.Fatalf("refusal changed status:\n%s\nwant:\n%s", got, statusBefore)
	}
	if got := runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"); got != headBefore {
		t.Fatalf("refusal changed HEAD: %s, want %s", got, headBefore)
	}
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-814")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(journalPath); !os.IsNotExist(err) {
		t.Fatalf("discovery refusal created journal %s: %v", journalPath, err)
	}
}

func TestREQ512ChangedWorkspaceRootsStillRequireRootLockMirrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		rootPath   string
		lockPath   string
		rootBefore string
		rootAfter  string
		lockBefore string
		lockAfter  string
	}{
		{
			name:       "npm",
			rootPath:   "skills/do-work/consumer/package.json",
			lockPath:   "skills/do-work/consumer/package-lock.json",
			rootBefore: "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"private\":true}\n",
			rootAfter:  "{\"name\":\"consumer\",\"version\":\"1.0.1\",\"private\":true}\n",
			lockBefore: "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"1.0.0\"}}}\n",
			lockAfter:  "{\"name\":\"consumer\",\"version\":\"1.0.1\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"1.0.1\"}}}\n",
		},
		{
			name:       "cargo",
			rootPath:   "skills/do-work/consumer-rust/Cargo.toml",
			lockPath:   "skills/do-work/consumer-rust/Cargo.lock",
			rootBefore: "[package]\nname = \"consumer\"\nversion = \"1.0.0\"\n",
			rootAfter:  "[package]\nname = \"consumer\"\nversion = \"1.0.1\"\n",
			lockBefore: "[[package]]\nname = \"consumer\"\nversion = \"1.0.0\"\n",
			lockAfter:  "[[package]]\nname = \"consumer\"\nversion = \"1.0.1\"\n",
		},
		{
			name:       "uv",
			rootPath:   "skills/do-work/consumer-python/pyproject.toml",
			lockPath:   "skills/do-work/consumer-python/uv.lock",
			rootBefore: "[project]\nname = \"consumer\"\nversion = \"1.0.0\"\n",
			rootAfter:  "[project]\nname = \"consumer\"\nversion = \"1.0.1\"\n",
			lockBefore: "[[package]]\nname = \"consumer\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n",
			lockAfter:  "[[package]]\nname = \"consumer\"\nversion = \"1.0.1\"\nsource = { editable = \".\" }\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, test.rootPath, test.rootBefore)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockBefore)
			seedSemanticLegacyTail(t, repositoryRoot)
			writeFinalizationFile(t, repositoryRoot, test.rootPath, test.rootAfter)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockAfter)

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
				t.Fatalf("changed root recovery = %#v", result)
			}
			for _, path := range []string{test.rootPath, test.lockPath} {
				if !containsFinalizationPath(result.Finalization.CommitPaths, path) {
					t.Fatalf("commit paths %v omit changed root path %s", result.Finalization.CommitPaths, path)
				}
			}
		})
	}
}

func TestREQ512WorkspaceMembersSelectChangedSourcesBeforeEqualVersionRoots(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		rootPath     string
		memberPath   string
		lockPath     string
		rootContents string
		memberBefore string
		memberAfter  string
		lockBefore   string
		lockAfter    string
	}{
		{
			name:         "npm",
			rootPath:     "skills/do-work/consumer/package.json",
			memberPath:   "skills/do-work/consumer/packages/widget/package.json",
			lockPath:     "skills/do-work/consumer/package-lock.json",
			rootContents: "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"private\":true,\"workspaces\":[\"packages/*\"]}\n",
			memberBefore: "{\"name\":\"widget\",\"version\":\"1.0.0\"}\n",
			memberAfter:  "{\"name\":\"widget\",\"version\":\"1.0.1\"}\n",
			lockBefore:   "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"1.0.0\"},\"packages/widget\":{\"name\":\"widget\",\"version\":\"1.0.0\"}}}\n",
			lockAfter:    "{\"name\":\"consumer\",\"version\":\"1.0.0\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"consumer\",\"version\":\"1.0.0\"},\"packages/widget\":{\"name\":\"widget\",\"version\":\"1.0.1\"}}}\n",
		},
		{
			name:         "cargo",
			rootPath:     "skills/do-work/consumer-rust/Cargo.toml",
			memberPath:   "skills/do-work/consumer-rust/crates/widget/Cargo.toml",
			lockPath:     "skills/do-work/consumer-rust/Cargo.lock",
			rootContents: "[package]\nname = \"consumer\"\nversion = \"1.0.0\"\n\n[workspace]\nmembers = [\"crates/*\"]\n",
			memberBefore: "[package]\nname = \"widget\"\nversion = \"1.0.0\"\n",
			memberAfter:  "[package]\nname = \"widget\"\nversion = \"1.0.1\"\n",
			lockBefore:   "[[package]]\nname = \"consumer\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"widget\"\nversion = \"1.0.0\"\n",
			lockAfter:    "[[package]]\nname = \"consumer\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"widget\"\nversion = \"1.0.1\"\n",
		},
		{
			name:         "uv",
			rootPath:     "skills/do-work/consumer-python/pyproject.toml",
			memberPath:   "skills/do-work/consumer-python/packages/widget/pyproject.toml",
			lockPath:     "skills/do-work/consumer-python/uv.lock",
			rootContents: "[project]\nname = \"consumer\"\nversion = \"1.0.0\"\n\n[tool.uv.workspace]\nmembers = [\"packages/*\"]\n",
			memberBefore: "[project]\nname = \"widget\"\nversion = \"1.0.0\"\n",
			memberAfter:  "[project]\nname = \"widget\"\nversion = \"1.0.1\"\n",
			lockBefore:   "[[package]]\nname = \"consumer\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"widget\"\nversion = \"1.0.0\"\nsource = { editable = \"packages/widget\" }\n",
			lockAfter:    "[[package]]\nname = \"consumer\"\nversion = \"1.0.0\"\nsource = { editable = \".\" }\n\n[[package]]\nname = \"widget\"\nversion = \"1.0.1\"\nsource = { editable = \"packages/widget\" }\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, test.rootPath, test.rootContents)
			writeFinalizationFile(t, repositoryRoot, test.memberPath, test.memberBefore)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockBefore)
			seedSemanticLegacyTail(t, repositoryRoot)
			rootBefore := readFinalizationFile(t, repositoryRoot, test.rootPath)
			writeFinalizationFile(t, repositoryRoot, test.memberPath, test.memberAfter)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockAfter)

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
				t.Fatalf("changed-source-first workspace recovery = %#v", result)
			}
			for _, path := range []string{test.memberPath, test.lockPath} {
				if !containsFinalizationPath(result.Finalization.CommitPaths, path) {
					t.Fatalf("commit paths %v omit changed workspace path %s", result.Finalization.CommitPaths, path)
				}
			}
			if containsFinalizationPath(result.Finalization.CommitPaths, test.rootPath) {
				t.Fatalf("commit paths %v select unchanged equal-version root %s", result.Finalization.CommitPaths, test.rootPath)
			}
			if got := readFinalizationFile(t, repositoryRoot, test.rootPath); got != rootBefore {
				t.Fatalf("unchanged root bytes = %q, want %q", got, rootBefore)
			}
		})
	}
}

func TestREQ512SharedWorkspaceLocksReplaceOnlyMultipleChangedMembers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		rootPath     string
		rootContents string
		lockPath     string
		lockBefore   string
		lockAfter    string
		memberPaths  []string
		memberBefore string
		memberAfter  string
	}{
		{
			name:         "npm",
			rootPath:     "skills/do-work/multi-npm/package.json",
			rootContents: "{\"name\":\"multi\",\"private\":true,\"workspaces\":[\"packages/*\"]}\n",
			lockPath:     "skills/do-work/multi-npm/package-lock.json",
			lockBefore:   "{\"name\":\"multi\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"multi\"},\"packages/alpha\":{\"name\":\"alpha\",\"version\":\"1.0.0\"},\"packages/beta\":{\"name\":\"beta\",\"version\":\"1.0.0\"},\"node_modules/dependency\":{\"name\":\"dependency\",\"version\":\"1.0.1\"}}}\n",
			lockAfter:    "{\"name\":\"multi\",\"lockfileVersion\":3,\"packages\":{\"\":{\"name\":\"multi\"},\"packages/alpha\":{\"name\":\"alpha\",\"version\":\"1.0.1\"},\"packages/beta\":{\"name\":\"beta\",\"version\":\"1.0.1\"},\"node_modules/dependency\":{\"name\":\"dependency\",\"version\":\"1.0.1\"}}}\n",
			memberPaths:  []string{"skills/do-work/multi-npm/packages/alpha/package.json", "skills/do-work/multi-npm/packages/beta/package.json"},
			memberBefore: "{\"name\":\"%s\",\"version\":\"1.0.0\"}\n",
			memberAfter:  "{\"name\":\"%s\",\"version\":\"1.0.1\"}\n",
		},
		{
			name:         "cargo",
			rootPath:     "skills/do-work/multi-rust/Cargo.toml",
			rootContents: "[workspace]\nmembers = [\"crates/*\"]\n",
			lockPath:     "skills/do-work/multi-rust/Cargo.lock",
			lockBefore:   "[[package]]\nname = \"alpha\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.1\"\nsource = \"registry+https://example.invalid\"\n",
			lockAfter:    "[[package]]\nname = \"alpha\"\nversion = \"1.0.1\"\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.1\"\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.1\"\nsource = \"registry+https://example.invalid\"\n",
			memberPaths:  []string{"skills/do-work/multi-rust/crates/alpha/Cargo.toml", "skills/do-work/multi-rust/crates/beta/Cargo.toml"},
			memberBefore: "[package]\nname = \"%s\"\nversion = \"1.0.0\"\n",
			memberAfter:  "[package]\nname = \"%s\"\nversion = \"1.0.1\"\n",
		},
		{
			name:         "uv",
			rootPath:     "skills/do-work/multi-python/pyproject.toml",
			rootContents: "[tool.uv.workspace]\nmembers = [\"packages/*\"]\n",
			lockPath:     "skills/do-work/multi-python/uv.lock",
			lockBefore:   "[[package]]\nname = \"alpha\"\nversion = \"1.0.0\"\nsource = { editable = \"packages/alpha\" }\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.0\"\nsource = { editable = \"packages/beta\" }\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.1\"\nsource = { registry = \"https://example.invalid\" }\n",
			lockAfter:    "[[package]]\nname = \"alpha\"\nversion = \"1.0.1\"\nsource = { editable = \"packages/alpha\" }\n\n[[package]]\nname = \"beta\"\nversion = \"1.0.1\"\nsource = { editable = \"packages/beta\" }\n\n[[package]]\nname = \"dependency\"\nversion = \"1.0.1\"\nsource = { registry = \"https://example.invalid\" }\n",
			memberPaths:  []string{"skills/do-work/multi-python/packages/alpha/pyproject.toml", "skills/do-work/multi-python/packages/beta/pyproject.toml"},
			memberBefore: "[project]\nname = \"%s\"\nversion = \"1.0.0\"\n",
			memberAfter:  "[project]\nname = \"%s\"\nversion = \"1.0.1\"\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, test.rootPath, test.rootContents)
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockBefore)
			for index, path := range test.memberPaths {
				name := []string{"alpha", "beta"}[index]
				writeFinalizationFile(t, repositoryRoot, path, fmt.Sprintf(test.memberBefore, name))
			}
			seedSemanticLegacyTail(t, repositoryRoot)
			for index, path := range test.memberPaths {
				name := []string{"alpha", "beta"}[index]
				writeFinalizationFile(t, repositoryRoot, path, fmt.Sprintf(test.memberAfter, name))
			}
			writeFinalizationFile(t, repositoryRoot, test.lockPath, test.lockAfter)

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
				t.Fatalf("shared lock recovery = %#v", result)
			}
			for _, path := range append(append([]string(nil), test.memberPaths...), test.lockPath) {
				if !containsFinalizationPath(result.Finalization.CommitPaths, path) {
					t.Fatalf("commit paths %v omit %s", result.Finalization.CommitPaths, path)
				}
			}
			if got := readFinalizationFile(t, repositoryRoot, test.lockPath); got != test.lockAfter {
				t.Fatalf("lock bytes changed unexpectedly:\n%s\nwant:\n%s", got, test.lockAfter)
			}
		})
	}
}
