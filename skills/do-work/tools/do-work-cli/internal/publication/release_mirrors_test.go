package publication

import (
	"reflect"
	"testing"
)

// The 0.282.0 release wrote root VERSION and CHANGELOG.md but left
// skills/do-work/actions/version.md and skills/do-work/CHANGELOG.md on the old
// version, because the manifest author listed the mirrors by hand. The planner
// must find every tracked mirror still carrying the old version and refuse a
// manifest that does not declare it.
func TestBuildReleasePlanRefusesUndeclaredMirrorsLeftOnOldVersion(t *testing.T) {
	root := initializedGitRepository(t)
	for _, path := range []string{"VERSION", "skills/do-work/VERSION"} {
		writeFixture(t, root, path, []byte("1.0.0\n"), 0o644)
	}
	writeFixture(t, root, "skills/do-work/actions/version.md", []byte("# Version Action\n\n**Current version**: 1.0.0\n"), 0o644)
	writeFixture(t, root, "suite/modules.tsv", []byte("source\tdestination\nskills/do-work\t.claude/skills/do-work\n"), 0o644)
	for _, path := range []string{"CHANGELOG.md", "skills/do-work/CHANGELOG.md"} {
		writeFixture(t, root, path, []byte("History\n* old\n"), 0o644)
	}
	writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/action-old", []byte("# Version Action\n\n**Current version**: 1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/action-new", []byte("# Version Action\n\n**Current version**: 1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n* old\n"), 0o644)
	runGitFixture(t, root, "add", "-A")

	versionTarget := func(path, oldPayload, newPayload string) ReleaseTarget {
		return ReleaseTarget{Path: path, ExpectedPayload: PayloadFile{SourcePath: oldPayload}, NewPayload: PayloadFile{SourcePath: newPayload}, OldVersion: "1.0.0", NewVersion: "1.0.1"}
	}
	changelogTarget := func(path string) ChangelogTarget {
		return ChangelogTarget{Path: path, ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}
	}
	partial := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		MaintainerRelease: true, OldVersion: "1.0.0", NewVersion: "1.0.1",
		ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		Targets:             []ReleaseTarget{versionTarget("VERSION", "payload/version-old", "payload/version-new")},
		Changelogs:          []ChangelogTarget{changelogTarget("CHANGELOG.md")},
	}}
	plan := BuildReleasePlan(root, partial)
	if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-MIRROR-UNDECLARED" {
		t.Fatalf("partial mirror set accepted: %#v", plan.Refusal)
	}
	wantPaths := []string{"skills/do-work/CHANGELOG.md", "skills/do-work/VERSION", "skills/do-work/actions/version.md"}
	if !reflect.DeepEqual(plan.Refusal.Paths, wantPaths) {
		t.Fatalf("undeclared paths = %v, want %v", plan.Refusal.Paths, wantPaths)
	}

	complete := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		MaintainerRelease: true, OldVersion: "1.0.0", NewVersion: "1.0.1",
		ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		RequiredMirrors:     []string{"skills/do-work/VERSION", "skills/do-work/actions/version.md", "skills/do-work/CHANGELOG.md"},
		Targets: []ReleaseTarget{
			versionTarget("VERSION", "payload/version-old", "payload/version-new"),
			versionTarget("skills/do-work/VERSION", "payload/version-old", "payload/version-new"),
			versionTarget("skills/do-work/actions/version.md", "payload/action-old", "payload/action-new"),
		},
		Changelogs: []ChangelogTarget{changelogTarget("CHANGELOG.md"), changelogTarget("skills/do-work/CHANGELOG.md")},
	}}
	if plan := BuildReleasePlan(root, complete); plan.Refusal != nil || len(plan.Mutations) != 5 {
		t.Fatalf("complete mirror set refused: %#v", plan)
	}
}

// An independently versioned component may legitimately sit on the release's
// own version — a fresh monorepo starts every package at the same number. The
// coincidence is not ownership evidence, so a root-only release must still
// plan, rather than refusing until the component joins someone else's bump.
func TestBuildReleasePlanIgnoresIndependentComponentSharingTheOldVersion(t *testing.T) {
	root := initializedGitRepository(t)
	writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "CHANGELOG.md", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "tools/independent/VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n* old\n"), 0o644)
	runGitFixture(t, root, "add", "-A")
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		OldVersion: "1.0.0", NewVersion: "1.0.1", ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		Targets:    []ReleaseTarget{{Path: "VERSION", ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
		Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
	}}
	if plan := BuildReleasePlan(root, manifest); plan.Refusal != nil || len(plan.Mutations) != 2 {
		t.Fatalf("independent component on the same version blocked the release: %#v", plan)
	}
}

// A workspace member the root manifest claims IS owned, so leaving its VERSION
// on the old value is the forgotten-mirror case the check exists for.
func TestBuildReleasePlanRefusesUndeclaredWorkspaceMemberMirror(t *testing.T) {
	root := initializedGitRepository(t)
	writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "CHANGELOG.md", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "package.json", []byte("{\"name\":\"root\",\"workspaces\":[\"packages/*\"]}\n"), 0o644)
	writeFixture(t, root, "packages/widget/package.json", []byte("{\"name\":\"widget\",\"version\":\"1.0.0\"}\n"), 0o644)
	writeFixture(t, root, "packages/widget/VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n* old\n"), 0o644)
	runGitFixture(t, root, "add", "-A")
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		OldVersion: "1.0.0", NewVersion: "1.0.1", ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		Targets:    []ReleaseTarget{{Path: "VERSION", ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
		Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
	}}
	plan := BuildReleasePlan(root, manifest)
	if plan.Refusal == nil || plan.Refusal.Code != "RELEASE-MIRROR-UNDECLARED" {
		t.Fatalf("owned workspace member mirror accepted: %#v", plan.Refusal)
	}
	if want := []string{"packages/widget/VERSION"}; !reflect.DeepEqual(plan.Refusal.Paths, want) {
		t.Fatalf("undeclared paths = %v, want %v", plan.Refusal.Paths, want)
	}
}

// A tracked VERSION on some other version is an independently versioned tool,
// not a mirror; an untracked file is scratch. Neither is required.
func TestBuildReleasePlanIgnoresVersionFilesThatAreNotMirrors(t *testing.T) {
	root := initializedGitRepository(t)
	writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "CHANGELOG.md", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "tools/board/VERSION", []byte("9.9.9\n"), 0o644)
	writeFixture(t, root, "docs/version.md", []byte("no current version line\n"), 0o644)
	writeFixture(t, root, "CHANGELOG-archive.md", []byte("Older history\n"), 0o644)
	writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-old", []byte("History\n* old\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n* old\n"), 0o644)
	runGitFixture(t, root, "add", "-A")
	writeFixture(t, root, "scratch/VERSION", []byte("1.0.0\n"), 0o644)
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		OldVersion: "1.0.0", NewVersion: "1.0.1", ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		Targets:    []ReleaseTarget{{Path: "VERSION", ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
		Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", ExpectedPayload: PayloadFile{SourcePath: "payload/log-old"}, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
	}}
	if plan := BuildReleasePlan(root, manifest); plan.Refusal != nil || len(plan.Mutations) != 2 {
		t.Fatalf("non-mirror version files were required: %#v", plan)
	}
}

// Without a tracked-file listing the mirror check cannot run, and a check that
// cannot run must refuse rather than pass silently.
func TestBuildReleasePlanRefusesWhenMirrorEnumerationFails(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "VERSION", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-old", []byte("1.0.0\n"), 0o644)
	writeFixture(t, root, "payload/version-new", []byte("1.0.1\n"), 0o644)
	writeFixture(t, root, "payload/log-new", []byte("History\nrelease-1.0.1 Title\n"), 0o644)
	manifest := Manifest{Operation: OperationRelease, Release: &ReleaseManifest{
		OldVersion: "1.0.0", NewVersion: "1.0.1", ProjectOwnedTargets: []string{"VERSION", "CHANGELOG.md"},
		Targets:    []ReleaseTarget{{Path: "VERSION", ExpectedPayload: PayloadFile{SourcePath: "payload/version-old"}, NewPayload: PayloadFile{SourcePath: "payload/version-new"}, OldVersion: "1.0.0", NewVersion: "1.0.1"}},
		Changelogs: []ChangelogTarget{{Path: "CHANGELOG.md", Create: true, NewPayload: PayloadFile{SourcePath: "payload/log-new"}, InsertionAnchor: "History\n", EntryKey: "release-1.0.1", EntryTitle: "Title"}},
	}}
	if plan := BuildReleasePlan(root, manifest); plan.Refusal == nil || plan.Refusal.Code != "RELEASE-MIRROR-ENUMERATION" {
		t.Fatalf("enumeration failure passed silently: %#v", plan.Refusal)
	}
}
