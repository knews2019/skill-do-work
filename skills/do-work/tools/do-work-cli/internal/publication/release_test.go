package publication

import (
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
