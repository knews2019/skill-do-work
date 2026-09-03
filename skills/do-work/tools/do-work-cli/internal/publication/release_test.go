package publication

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildReleasePlanAcceptsParameterizedNonHouseChangelog(t *testing.T) {
	root := newReleaseRepository(t)
	writeFixture(t, root, "VERSION", []byte("1.2.3\n"), 0o644)
	writeFixture(t, root, "CHANGELOG.txt", []byte("History\n* old\n"), 0o644)
	commitReleaseFixtures(t, root, "VERSION", "CHANGELOG.txt")
	writeFixture(t, root, "payload/version-old", []byte("1.2.3\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.2.4\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.2.4 Fresh title\n* old\n"), 0o644)
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		Targets:    []ReleaseTarget{{Path: "VERSION", ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.2.3", NewVersion: "1.2.4"}},
		Changelogs: []ChangelogTarget{{Path: "CHANGELOG.txt", ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.2.4", EntryTitle: "Fresh title"}},
	}}
	plan := BuildReleasePlan(root, manifest)
	if plan.Refusal != nil || len(plan.Mutations) != 2 {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestBuildReleasePlanBootstrapsChangelogWithoutVersionFile(t *testing.T) {
	root := newReleaseRepository(t)
	writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
	commitReleaseFixtures(t, root, "README.md")
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-0.1.0 First release\n"), 0o644)
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{OldVersion: "0.0.0", NewVersion: "0.1.0", Changelogs: []ChangelogTarget{{Path: "CHANGELOG.txt", Create: true, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-0.1.0", EntryTitle: "First release"}}}}
	plan := BuildReleasePlan(root, manifest)
	if plan.Refusal != nil || plan.Mutations[0].Kind != MutationCreate {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestBuildReleasePlanRefusesNonMonotonicAndInstalledMetadata(t *testing.T) {
	root := t.TempDir()
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{Targets: []ReleaseTarget{{OldVersion: "2.0.0", NewVersion: "1.0.0"}}, Changelogs: []ChangelogTarget{{}}}}
	if plan := BuildReleasePlan(root, manifest); plan.Refusal == nil || plan.Refusal.Code != "RELEASE-VERSION-NOT-INCREASING" {
		t.Fatalf("refusal = %#v", plan.Refusal)
	}
}

// TestBuildReleasePlanRequiresAffirmativeProjectOwnedVersionTargets pins the
// REQ-461 invariant: a consumer release target is admitted only when the
// repository's own declarations attest it, never because its directory name is
// absent from a list of known dependency spellings. Every row is named by the
// class of location it stands for, because the paths themselves are only
// examples of that class — the previous denylist agreed with the fixtures of its
// day and still admitted `third_party/`, `dist/` and any cache tree spelled
// something other than `generated`.
func TestBuildReleasePlanRequiresAffirmativeProjectOwnedVersionTargets(t *testing.T) {
	for _, testCase := range []struct {
		name              string
		targetPath        string
		packageMarker     string
		trackTarget       bool
		maintainerRelease bool
		evidenceFragment  string
	}{
		{
			name:        "project-owned source the repository tracks",
			targetPath:  "VERSION",
			trackTarget: true,
		},
		{
			name:              "declared maintainer mirror inside the suite package",
			targetPath:        "skills/do-work/VERSION",
			packageMarker:     "skills/do-work/SKILL.md",
			trackTarget:       true,
			maintainerRelease: true,
		},
		{
			name:             "installed package under an unlisted dependency directory",
			targetPath:       "third_party/do-work/VERSION",
			packageMarker:    "third_party/do-work/SKILL.md",
			trackTarget:      true,
			evidenceFragment: "SKILL.md",
		},
		{
			name:             "distribution output tree",
			targetPath:       "dist/skills/do-work/VERSION",
			evidenceFragment: "does not track",
		},
		{
			name:             "arbitrarily named cache tree",
			targetPath:       ".stash-pantry/packages/do-work/VERSION",
			evidenceFragment: "does not track",
		},
		{
			name:             "committed vendored package",
			targetPath:       "vendor/do-work/VERSION",
			packageMarker:    "vendor/do-work/SKILL.md",
			trackTarget:      true,
			evidenceFragment: "SKILL.md",
		},
		{
			name:             "dependency install directory",
			targetPath:       "node_modules/do-work/VERSION",
			evidenceFragment: "does not track",
		},
		{
			name:             "agent skill install directory",
			targetPath:       ".claude/skills/do-work/VERSION",
			packageMarker:    ".claude/skills/do-work/SKILL.md",
			trackTarget:      true,
			evidenceFragment: "SKILL.md",
		},
		{
			name:             "in-tree suite package on a consumer release",
			targetPath:       "skills/do-work/VERSION",
			packageMarker:    "skills/do-work/SKILL.md",
			trackTarget:      true,
			evidenceFragment: "SKILL.md",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := newReleaseRepository(t)
			writeFixture(t, root, "CHANGELOG.md", []byte("History\n"), 0o644)
			writeFixture(t, root, testCase.targetPath, []byte("1.0.0\n"), 0o644)
			trackedPaths := []string{"CHANGELOG.md"}
			if testCase.packageMarker != "" {
				writeFixture(t, root, testCase.packageMarker, []byte("# package\n"), 0o644)
				trackedPaths = append(trackedPaths, testCase.packageMarker)
			}
			if testCase.trackTarget {
				trackedPaths = append(trackedPaths, testCase.targetPath)
			}
			commitReleaseFixtures(t, root, trackedPaths...)
			writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
			writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
			writeFixture(t, root, "payload/log-old", []byte("History\n"), 0o644)
			writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n"), 0o644)
			plan := BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
				MaintainerRelease: testCase.maintainerRelease, OldVersion: "1.0.0", NewVersion: "1.0.1",
				Targets:    []ReleaseTarget{{Path: testCase.targetPath, ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
				Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
			}})
			if testCase.evidenceFragment == "" {
				if plan.Refusal != nil || len(plan.Mutations) != 2 {
					t.Fatalf("project-owned target refused: %#v", plan)
				}
				return
			}
			if plan.Refusal == nil {
				t.Fatalf("target without ownership evidence accepted: %#v", plan)
			}
			if plan.Refusal.Code != "RELEASE-TARGET-OWNERSHIP-UNVERIFIED" {
				t.Fatalf("refusal code = %q (reason %q)", plan.Refusal.Code, plan.Refusal.Reason)
			}
			assertOwnershipRefusalNamesTargetAndEvidence(t, plan.Refusal, testCase.targetPath, testCase.evidenceFragment)
		})
	}
}

// TestBuildReleasePlanRefusesBootstrapChangelogInAnUnattestedLocation covers the
// created-target half of the ownership rule: a new changelog cannot be tracked
// yet, so the attestation has to come from the nearest location that already
// exists. An untracked distribution directory attests nothing.
func TestBuildReleasePlanRefusesBootstrapChangelogInAnUnattestedLocation(t *testing.T) {
	root := newReleaseRepository(t)
	writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
	commitReleaseFixtures(t, root, "README.md")
	writeFixture(t, root, "dist/skills/do-work/placeholder", []byte("built\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-0.1.0 First release\n"), 0o644)
	plan := BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{OldVersion: "0.0.0", NewVersion: "0.1.0", Changelogs: []ChangelogTarget{
		{Path: "dist/skills/do-work/CHANGELOG.md", Create: true, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-0.1.0", EntryTitle: "First release"},
	}}})
	if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-TARGET-OWNERSHIP-UNVERIFIED" {
		t.Fatalf("bootstrap changelog in an unattested location accepted: %#v", plan)
	}
	assertOwnershipRefusalNamesTargetAndEvidence(t, plan.Refusal, "dist/skills/do-work/CHANGELOG.md", "dist/skills/do-work")
}

// TestBuildReleasePlanPlansThisRepositorysOwnMaintainerReleaseShape reproduces
// the suite's own release: three of its five targets sit under the very
// `skills/do-work/` package marker that refuses a consumer release, and
// `maintainer_release` is the only door through which they may be mutated.
func TestBuildReleasePlanPlansThisRepositorysOwnMaintainerReleaseShape(t *testing.T) {
	root := newReleaseRepository(t)
	versionPaths := []string{"VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md"}
	writeFixture(t, root, "skills/do-work/SKILL.md", []byte("# do-work\n"), 0o644)
	writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "skills/do-work/VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "skills/do-work/actions/version.md", []byte("**Current version**: 1.0.0\n"), 0o644)
	writeFixture(t, root, "CHANGELOG.md", []byte("# Changelog\n\n## 1.0.0 — Seed (2026-01-01)\n"), 0o644)
	writeFixture(t, root, "skills/do-work/CHANGELOG.md", []byte("# Changelog\n\n## 1.0.0 — Seed (2026-01-01)\n"), 0o644)
	commitReleaseFixtures(t, root, "skills/do-work/SKILL.md", "VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md")
	writeFixture(t, root, "payload/plain-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/plain-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/action-old", []byte("**Current version**: 1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/action-new", []byte("**Current version**: 1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("# Changelog\n\n## 1.0.0 — Seed (2026-01-01)\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("# Changelog\n\n## 1.0.1 — Ownership Evidence (2026-09-03)\n\nDelivered.\n\n## 1.0.0 — Seed (2026-01-01)\n"), 0o644)
	targets := []ReleaseTarget{}
	for _, path := range versionPaths {
		expected, updated := "payload/plain-old", "payload/plain-new"
		if strings.HasSuffix(path, ".md") {
			expected, updated = "payload/action-old", "payload/action-new"
		}
		targets = append(targets, ReleaseTarget{Path: path, ExpectedPayload: PayloadFile{SourcePath: expected}, NewPayload: PayloadFile{SourcePath: updated}, OldVersion: "1.0.0", NewVersion: "1.0.1"})
	}
	changelogs := []ChangelogTarget{}
	for _, path := range []string{"CHANGELOG.md", "skills/do-work/CHANGELOG.md"} {
		changelogs = append(changelogs, ChangelogTarget{Path: path, ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "## 1.0.0", EntryKey: "## 1.0.1", EntryTitle: "Ownership Evidence"})
	}
	plan := BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		MaintainerRelease: true, OldVersion: "1.0.0", NewVersion: "1.0.1",
		RequiredMirrors: versionPaths, Targets: targets, Changelogs: changelogs,
	}})
	if plan.Refusal != nil || len(plan.Mutations) != 5 {
		t.Fatalf("maintainer release shape refused: %#v", plan)
	}
}

func TestRemediationF7ReleaseRefusesDivergentChangelogMirrors(t *testing.T) {
	// REQ-461 moved this test's installed-path rows into
	// TestBuildReleasePlanRequiresAffirmativeProjectOwnedVersionTargets: those
	// paths must now refuse for missing ownership evidence rather than for
	// matching a directory name, which is a different fixture shape.
	root := t.TempDir()
	writeFixture(t, root, "payload/root-log", []byte("History\nrelease-1.0.1 Root title\n"), 0o644)
	writeFixture(t, root, "payload/mirror-log", []byte("History\nrelease-1.0.1 Mirror title\n"), 0o644)
	plan := BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{MaintainerRelease: true, OldVersion: "1.0.0", NewVersion: "1.0.1", Changelogs: []ChangelogTarget{
		{Path: "CHANGELOG.md", Create: true, NewPayload: PayloadFile{SourcePath: "payload/root-log"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Root title"},
		{Path: ".codex/skills/do-work/CHANGELOG.md", Create: true, NewPayload: PayloadFile{SourcePath: "payload/mirror-log"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Mirror title"},
	}}})
	if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-CHANGELOG-MIRROR-DIVERGED" {
		t.Fatalf("divergent changelog mirrors accepted: %#v", plan.Refusal)
	}
}

func TestRemediationF10ReleaseReferenceHasNoActiveHandEditDirective(t *testing.T) {
	referencePath := filepath.Join("..", "..", "..", "..", "actions", "work-reference.md")
	contents, err := os.ReadFile(referencePath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(contents), "constraints on the hand edit") {
		t.Fatal("active lockfile hand-edit directive remains")
	}
}

func assertOwnershipRefusalNamesTargetAndEvidence(t *testing.T, refusal *Refusal, targetPath, evidenceFragment string) {
	t.Helper()
	if len(refusal.Paths) != 1 || refusal.Paths[0] != targetPath {
		t.Fatalf("refusal paths = %#v, want the single target %q", refusal.Paths, targetPath)
	}
	if !strings.Contains(refusal.Reason, targetPath) {
		t.Fatalf("refusal reason %q does not name the target %q", refusal.Reason, targetPath)
	}
	if !strings.Contains(refusal.Reason, evidenceFragment) {
		t.Fatalf("refusal reason %q does not name the missing ownership evidence %q", refusal.Reason, evidenceFragment)
	}
}

func newReleaseRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runReleaseGit(t, root, "init", "-q")
	runReleaseGit(t, root, "config", "user.name", "Release Fixture")
	runReleaseGit(t, root, "config", "user.email", "release@example.invalid")
	return root
}

func commitReleaseFixtures(t *testing.T, root string, paths ...string) {
	t.Helper()
	runReleaseGit(t, root, append([]string{"add", "--"}, paths...)...)
	runReleaseGit(t, root, "commit", "-qm", "seed release fixture")
}

func runReleaseGit(t *testing.T, root string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
}

// TestBuildReleasePlanRequiresAffirmativeAttestationForCreatedReleaseTargets
// pins the created-target half of the REQ-461 invariant. A target the release
// creates cannot be in the index yet, so its attestation has to come from the
// directory it is created in: that directory must already exist and the index
// must claim a source inside it. Walking further up the tree cannot substitute
// for that, because the repository root always holds tracked sources — so a
// root-anchored walk attests every not-yet-built tree, which is exactly the
// state a real repository is in before `dist/` or `.claude/skills/` exists.
// Every row is named by the class of location it stands for; the paths are only
// examples of that class.
func TestBuildReleasePlanRequiresAffirmativeAttestationForCreatedReleaseTargets(t *testing.T) {
	for _, testCase := range []struct {
		name             string
		setUpRepository  func(t *testing.T, root string)
		targetPath       string
		evidenceFragment string
	}{
		{
			name: "repository root holding tracked sources",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
			},
			targetPath: "CHANGELOG.md",
		},
		{
			name: "existing project directory holding tracked sources",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/guide.md", []byte("# Guide\n"), 0o644)
				commitReleaseFixtures(t, root, "docs/guide.md")
			},
			targetPath: "docs/CHANGELOG.md",
		},
		{
			name: "distribution output tree that does not exist yet",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
			},
			targetPath:       "dist/skills/do-work/CHANGELOG.md",
			evidenceFragment: "dist/skills/do-work",
		},
		{
			name: "agent skill install tree that does not exist yet",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
			},
			targetPath:       ".claude/skills/do-work/CHANGELOG.md",
			evidenceFragment: ".claude/skills/do-work",
		},
		{
			name: "arbitrarily named cache tree that does not exist yet",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
			},
			targetPath:       ".stash-pantry/packages/do-work/CHANGELOG.md",
			evidenceFragment: ".stash-pantry/packages/do-work",
		},
		{
			name: "existing directory holding only untracked build output",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
				writeFixture(t, root, "dist/skills/do-work/placeholder", []byte("built\n"), 0o644)
			},
			targetPath:       "dist/skills/do-work/CHANGELOG.md",
			evidenceFragment: "dist/skills/do-work",
		},
		{
			name: "directory the index tracks as a single object rather than as sources",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				writeFixture(t, root, "real-package/notes.md", []byte("notes\n"), 0o644)
				if err := os.MkdirAll(filepath.Join(root, "vendor"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(root, "real-package"), filepath.Join(root, "vendor", "lib")); err != nil {
					t.Skipf("filesystem does not support symlinks: %v", err)
				}
				commitReleaseFixtures(t, root, "README.md", "real-package/notes.md", "vendor/lib")
			},
			targetPath:       "vendor/lib/CHANGELOG.md",
			evidenceFragment: "vendor/lib",
		},
		{
			name: "installed package marker above an attested location",
			setUpRepository: func(t *testing.T, root string) {
				writeFixture(t, root, "packages/do-work/SKILL.md", []byte("# do-work\n"), 0o644)
				commitReleaseFixtures(t, root, "packages/do-work/SKILL.md")
			},
			targetPath:       "packages/do-work/CHANGELOG.md",
			evidenceFragment: "SKILL.md",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := newReleaseRepository(t)
			testCase.setUpRepository(t, root)
			plan := buildBootstrapChangelogReleasePlan(t, root, testCase.targetPath)
			if testCase.evidenceFragment == "" {
				if plan.Refusal != nil || len(plan.Mutations) != 1 || plan.Mutations[0].Kind != MutationCreate {
					t.Fatalf("attested created target refused: %#v", plan)
				}
				return
			}
			if plan.Refusal == nil {
				t.Fatalf("created target in an unattested location accepted: %#v", plan)
			}
			if plan.Refusal.Code != "RELEASE-TARGET-OWNERSHIP-UNVERIFIED" {
				t.Fatalf("refusal code = %q (reason %q)", plan.Refusal.Code, plan.Refusal.Reason)
			}
			assertOwnershipRefusalNamesTargetAndEvidence(t, plan.Refusal, testCase.targetPath, testCase.evidenceFragment)
		})
	}
}

// TestBuildReleasePlanOwnershipRefusalRemedyMatchesTheGap pins the actionability
// half of the REQ's fifth requirement. Naming the missing evidence is not enough
// when every refusal ends with the same remedy: `maintainer_release` is the
// answer only for a suite package's own marker, and pointing a caller whose
// target is merely uncommitted or ignored at that escape hatch steers them past
// the fix their gap actually needs.
func TestBuildReleasePlanOwnershipRefusalRemedyMatchesTheGap(t *testing.T) {
	for _, testCase := range []struct {
		name              string
		plan              func(t *testing.T, root string) PublicationPlan
		requiredRemedy    string
		forbiddenFragment string
	}{
		{
			name: "installed package marker owns the subtree",
			plan: func(t *testing.T, root string) PublicationPlan {
				writeFixture(t, root, "third_party/do-work/SKILL.md", []byte("# do-work\n"), 0o644)
				writeFixture(t, root, "third_party/do-work/VERSION", []byte("1.0.0\n"), 0o644)
				commitReleaseFixtures(t, root, "third_party/do-work/SKILL.md", "third_party/do-work/VERSION")
				return buildConsumerVersionReleasePlan(t, root, "third_party/do-work/VERSION")
			},
			requiredRemedy: "maintainer_release",
		},
		{
			name: "the existing target is not committed",
			plan: func(t *testing.T, root string) PublicationPlan {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
				writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
				return buildConsumerVersionReleasePlan(t, root, "VERSION")
			},
			requiredRemedy:    "commit",
			forbiddenFragment: "maintainer_release",
		},
		{
			name: "the repository's ignore rules exclude the created target",
			plan: func(t *testing.T, root string) PublicationPlan {
				writeFixture(t, root, ".gitignore", []byte("ignored-output/\n"), 0o644)
				commitReleaseFixtures(t, root, ".gitignore")
				writeFixture(t, root, "ignored-output/placeholder", []byte("built\n"), 0o644)
				return buildBootstrapChangelogReleasePlan(t, root, "ignored-output/CHANGELOG.md")
			},
			requiredRemedy:    "ignore rule",
			forbiddenFragment: "maintainer_release",
		},
		{
			name: "no tracked source attests the created target's location",
			plan: func(t *testing.T, root string) PublicationPlan {
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "README.md")
				return buildBootstrapChangelogReleasePlan(t, root, "dist/skills/do-work/CHANGELOG.md")
			},
			requiredRemedy:    "commit",
			forbiddenFragment: "maintainer_release",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := newReleaseRepository(t)
			plan := testCase.plan(t, root)
			if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-TARGET-OWNERSHIP-UNVERIFIED" {
				t.Fatalf("expected an ownership refusal, got %#v", plan)
			}
			if !strings.Contains(plan.Refusal.Reason, testCase.requiredRemedy) {
				t.Fatalf("refusal reason %q does not carry the remedy %q for this gap", plan.Refusal.Reason, testCase.requiredRemedy)
			}
			if testCase.forbiddenFragment != "" && strings.Contains(plan.Refusal.Reason, testCase.forbiddenFragment) {
				t.Fatalf("refusal reason %q offers %q, which does not resolve this gap", plan.Refusal.Reason, testCase.forbiddenFragment)
			}
		})
	}
}

// TestBuildReleasePlanRefusesEveryTargetWhenGitCannotAttestOwnership pins the
// fail-closed direction the whole ownership rule rests on. Ownership is proven
// from Git's index, so when Git cannot answer at all the target stays unproven
// and the release refuses rather than falling through to acceptance. Both halves
// need it: a target that must already exist, and one the release creates.
func TestBuildReleasePlanRefusesEveryTargetWhenGitCannotAttestOwnership(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		prepareRoot func(t *testing.T) string
	}{
		{
			name: "git is absent from PATH",
			prepareRoot: func(t *testing.T) string {
				root := newReleaseRepository(t)
				writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				commitReleaseFixtures(t, root, "VERSION", "README.md")
				t.Setenv("PATH", "")
				return root
			},
		},
		{
			name: "the release root is not a Git repository",
			prepareRoot: func(t *testing.T) string {
				root := t.TempDir()
				writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
				writeFixture(t, root, "README.md", []byte("# Project\n"), 0o644)
				return root
			},
		},
	} {
		t.Run(testCase.name+"/existing target", func(t *testing.T) {
			root := testCase.prepareRoot(t)
			plan := buildConsumerVersionReleasePlan(t, root, "VERSION")
			if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-TARGET-OWNERSHIP-UNVERIFIED" {
				t.Fatalf("unattested existing target accepted: %#v", plan)
			}
		})
		t.Run(testCase.name+"/created target", func(t *testing.T) {
			root := testCase.prepareRoot(t)
			plan := buildBootstrapChangelogReleasePlan(t, root, "CHANGELOG.md")
			if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-TARGET-OWNERSHIP-UNVERIFIED" {
				t.Fatalf("unattested created target accepted: %#v", plan)
			}
		})
	}
}

// buildBootstrapChangelogReleasePlan plans a consumer release whose only target
// is a changelog the release creates.
func buildBootstrapChangelogReleasePlan(t *testing.T, root, targetPath string) PublicationPlan {
	t.Helper()
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-0.1.0 First release\n"), 0o644)
	return BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{OldVersion: "0.0.0", NewVersion: "0.1.0", Changelogs: []ChangelogTarget{
		{Path: targetPath, Create: true, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-0.1.0", EntryTitle: "First release"},
	}}})
}

// buildConsumerVersionReleasePlan plans a consumer release whose single version
// target already exists on disk carrying 1.0.0.
func buildConsumerVersionReleasePlan(t *testing.T, root, targetPath string) PublicationPlan {
	t.Helper()
	writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("History\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n"), 0o644)
	return BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		OldVersion: "1.0.0", NewVersion: "1.0.1",
		Targets:    []ReleaseTarget{{Path: targetPath, ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
		Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
	}})
}
