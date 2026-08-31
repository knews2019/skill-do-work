package publication

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildReleasePlanAcceptsParameterizedNonHouseChangelog(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "VERSION", []byte("1.2.3\n"), 0o644)
	writeFixture(t, root, "CHANGELOG.txt", []byte("History\n* old\n"), 0o644)
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
	root := t.TempDir()
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

func TestRemediationF7ReleaseRefusesInstalledTreesAndDivergentChangelogMirrors(t *testing.T) {
	for _, path := range []string{".codex/skills/do-work/VERSION", ".claude/skills/do-work/VERSION", "generated/skills/do-work/VERSION", "vendor/do-work/VERSION", "vendored/do-work/VERSION"} {
		t.Run(path, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, path, []byte("1.0.0\n"), 0o644)
			writeFixture(t, root, "payload/old", []byte("1.0.0\n"), 0o644)
			writeFixture(t, root, "payload/new", []byte("1.0.1\n"), 0o644)
			writeFixture(t, root, "payload/log", []byte("History\nrelease-1.0.1 Title\n"), 0o644)
			plan := BuildReleasePlan(root, Manifest{Operation: OperationRelease, Release: &ReleaseManifest{OldVersion: "1.0.0", NewVersion: "1.0.1",
				Targets:    []ReleaseTarget{{Path: path, ExpectedPayload: PayloadFile{SourcePath: "payload/old"}, NewPayload: PayloadFile{SourcePath: "payload/new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
				Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", Create: true, NewPayload: PayloadFile{SourcePath: "payload/log"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
			}})
			if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-INSTALLED-METADATA-REFUSED" {
				t.Fatalf("installed path accepted: %#v", plan.Refusal)
			}
		})
	}

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
