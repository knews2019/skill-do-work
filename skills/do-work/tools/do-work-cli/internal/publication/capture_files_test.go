package publication

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestRemediationF2CaptureRequiresCanonicalURAndREQDestinations(t *testing.T) {
	for _, test := range []struct {
		name    string
		urPath  string
		reqPath string
	}{
		{name: "noncanonical UR", urPath: "unrelated/UR-1/not-really-input.md", reqPath: "do-work/queue/REQ-1-work.md"},
		{name: "noncanonical REQ", urPath: "do-work/user-requests/UR-1/input.md", reqPath: "unrelated/REQ-1-work.md"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "payload/ur.md", []byte("---\nid: UR-1\nrequests: [REQ-1]\n---\n"), 0o644)
			writeFixture(t, root, "payload/req.md", []byte("---\nid: REQ-1\nstatus: pending\nuser_request: UR-1\n---\n"), 0o644)
			plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
				UserRequest: PublishedFile{Path: test.urPath, Payload: PayloadFile{SourcePath: "payload/ur.md"}},
				Requests:    []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: test.reqPath, Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-1"}},
			}})
			if plan.Refusal == nil {
				t.Fatalf("accepted noncanonical paths: %#v", plan)
			}
		})
	}
}

func TestRemediationMinorCaptureMutationOrderIsMarkerURAssetsREQFold(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "payload/ur.md", []byte("---\nid: UR-1\nrequests: [REQ-1]\n---\n"), 0o644)
	writeFixture(t, root, "payload/req.md", []byte("---\nid: REQ-1\nstatus: pending\nuser_request: UR-1\n---\n"), 0o644)
	writeFixture(t, root, "payload/asset", []byte("asset"), 0o644)
	writeFixture(t, root, "do-work/prose-backlog.md", []byte("old\n"), 0o644)
	writeFixture(t, root, "payload/fold-old", []byte("old\n"), 0o644)
	writeFixture(t, root, "payload/fold-new", []byte("new\n"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
		UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests:    []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-1-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-1"}},
		Assets:      []PublishedFile{{Path: "do-work/user-requests/UR-1/assets/a.bin", Payload: PayloadFile{SourcePath: "payload/asset"}}},
		Folds:       []ReplacementFile{{Path: "do-work/prose-backlog.md", ExpectedPayload: PayloadFile{SourcePath: "payload/fold-old"}, NewPayload: PayloadFile{SourcePath: "payload/fold-new"}}},
	}})
	if plan.Refusal != nil {
		t.Fatal(plan.Refusal)
	}
	got := make([]string, len(plan.Mutations))
	for index, mutation := range plan.Mutations {
		got[index] = mutation.Path
	}
	want := []string{"do-work/.req-reservations/REQ-1", "do-work/user-requests/UR-1/input.md", "do-work/user-requests/UR-1/assets/a.bin", "do-work/queue/REQ-1-work.md", "do-work/prose-backlog.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mutation order = %v, want %v", got, want)
	}
}

func TestRemediationF9CaptureRefusesCollisionAtEveryCreatePosition(t *testing.T) {
	for _, collisionPath := range []string{
		"do-work/.req-reservations/REQ-1",
		"do-work/user-requests/UR-1/input.md",
		"do-work/user-requests/UR-1/assets/a.bin",
		"do-work/queue/REQ-1-work.md",
	} {
		t.Run(collisionPath, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "payload/ur.md", []byte("---\nid: UR-1\nrequests: [REQ-1]\n---\n"), 0o644)
			writeFixture(t, root, "payload/req.md", []byte("---\nid: REQ-1\nstatus: pending\nuser_request: UR-1\n---\n"), 0o644)
			writeFixture(t, root, "payload/asset", []byte("asset"), 0o644)
			writeFixture(t, root, collisionPath, []byte("collision"), 0o644)
			plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
				UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
				Requests:    []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-1-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-1"}},
				Assets:      []PublishedFile{{Path: "do-work/user-requests/UR-1/assets/a.bin", Payload: PayloadFile{SourcePath: "payload/asset"}}},
			}})
			if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-COLLISION" {
				t.Fatalf("collision accepted at %s: %#v", collisionPath, plan.Refusal)
			}
		})
	}
}

func TestRemediationF9CaptureRefusesStaleFoldPreimage(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "payload/ur.md", []byte("---\nid: UR-1\nrequests: []\n---\n"), 0o644)
	writeFixture(t, root, "do-work/prose-backlog.md", []byte("current\n"), 0o644)
	writeFixture(t, root, "payload/fold-old", []byte("old\n"), 0o644)
	writeFixture(t, root, "payload/fold-new", []byte("new\n"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
		UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Folds:       []ReplacementFile{{Path: "do-work/prose-backlog.md", ExpectedPayload: PayloadFile{SourcePath: "payload/fold-old"}, NewPayload: PayloadFile{SourcePath: "payload/fold-new"}}},
	}})
	if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-FOLD-STALE" {
		t.Fatalf("stale fold accepted: %#v", plan.Refusal)
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
