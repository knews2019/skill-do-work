package publication

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildCapturePlanBootstrapsAbsentDoWorkAndPinsRawInput(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeFixture(t, repositoryRoot, "payload/raw.txt", []byte("hello\n````\n## outside"), 0o600)
	raw, _ := os.ReadFile(filepath.Join(repositoryRoot, "payload/raw.txt"))
	ur := []byte("---\nid: UR-1\nrequests: [REQ-1]\n---\n# Input\n\n" + string(containedOutsideBytes(raw, "\n")) + "\n")
	writeFixture(t, repositoryRoot, "payload/ur.md", ur, 0o644)
	writeFixture(t, repositoryRoot, "payload/req.md", []byte("---\nid: REQ-1\nstatus: pending\nuser_request: UR-1\n---\n# Work\n"), 0o644)
	manifest := Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
		UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}}, RawInput: &PayloadFile{SourcePath: "payload/raw.txt"},
		Requests: []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-1-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-1"}},
	}}
	plan := BuildCapturePlan(repositoryRoot, manifest)
	if plan.Refusal != nil {
		t.Fatalf("refusal = %#v", plan.Refusal)
	}
	if len(plan.Mutations) != 3 {
		t.Fatalf("mutations = %#v", plan.Mutations)
	}
	if len(plan.CreatedDirectoryPaths) == 0 || plan.CreatedDirectoryPaths[0] != "do-work" {
		t.Fatalf("created dirs = %v", plan.CreatedDirectoryPaths)
	}
}

func TestBuildCapturePlanRefusesUnsafeBootstrapTopology(t *testing.T) {
	repositoryRoot := t.TempDir()
	if err := os.Symlink(t.TempDir(), filepath.Join(repositoryRoot, "do-work")); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, repositoryRoot, "payload/ur.md", []byte("---\nid: UR-1\nrequests: []\n---\n"), 0o644)
	plan := BuildCapturePlan(repositoryRoot, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}}}})
	if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-TOPOLOGY-UNSAFE" {
		t.Fatalf("refusal = %#v", plan.Refusal)
	}
}

func writeFixture(t *testing.T, root, path string, contents []byte, mode os.FileMode) {
	t.Helper()
	absolute := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, contents, mode); err != nil {
		t.Fatal(err)
	}
}
