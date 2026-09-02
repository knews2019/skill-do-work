package publication

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecodeManifestIsStrictAndOperationSpecific(t *testing.T) {
	valid := `{"operation":"capture-files","capture":{"user_request_id":"UR-1","user_request":{"path":"do-work/user-requests/UR-1/input.md","payload":{"source_path":"payload"}},"requests":[]}}`
	if _, err := DecodeManifest(strings.NewReader(valid), OperationCaptureFiles); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []string{
		`{"operation":"capture-files","capture":{"user_request_id":"UR-1","user_request":{"path":"x","payload":{"source_path":"p"}},"requests":[]},"unknown":true}`,
		`{"operation":"capture-files","capture":{"user_request_id":"UR-1","user_request":{"path":"x","payload":{"source_path":"p"}},"requests":[]},"answer":{"request_path":"x"}}`,
		`{"operation":"answer","capture":{"user_request_id":"UR-1","user_request":{"path":"x","payload":{"source_path":"p"}},"requests":[]}}`,
	} {
		if _, err := DecodeManifest(strings.NewReader(fixture), OperationCaptureFiles); err == nil {
			t.Fatalf("accepted %s", fixture)
		}
	}
}

func TestDecodeManifestAcceptsOnlyTypedDeferGateBody(t *testing.T) {
	valid := `{"operation":"defer-gate","defer_gate":{"parent_id":"REQ-1","parent_path":"do-work/working/REQ-1.md","expected_parent":{"source_path":"parent"},"expected_status":"claimed","checkpoint_path":"do-work/CHECKPOINT.md","expected_checkpoint":{"source_path":"checkpoint"},"writer_label":"host:/repo","gate_command":["go","test"],"gate_exit_status":1,"diagnostic_fingerprint":"fingerprint","diagnostic_evidence":["failed"],"sweep_key":"gate-fingerprint","repair_id":"REQ-2","repair_path":"do-work/queue/REQ-2.md","repair_title":"Repair gate","repair_created_at":"2026-09-02T01:02:03Z","reservation_path":"do-work/.req-reservations/REQ-2"}}`
	if _, err := DecodeManifest(strings.NewReader(valid), OperationDeferGate); err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeManifest(strings.NewReader(strings.Replace(valid, `}}`, `,"invented":true}}`, 1)), OperationDeferGate); err == nil {
		t.Fatal("unknown defer-gate field accepted")
	}
}

func TestReadPayloadRejectsSymlinkSource(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "payload/real", []byte("data"), 0o644)
	if err := os.Symlink("real", filepath.Join(root, "payload/link")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readPayload(root, PayloadFile{SourcePath: "payload/link"}); err == nil {
		t.Fatal("symlink payload accepted")
	}
}

func TestContainedPathRejectsAbsoluteAndEscapingPaths(t *testing.T) {
	for _, path := range []string{"", "/tmp/x", "../x", "do-work/../../x"} {
		if _, err := containedPath(path); err == nil {
			t.Errorf("accepted %q", path)
		}
	}
}

func TestOutsideTextContainmentUsesLongerFence(t *testing.T) {
	raw := []byte("line\n````\n## heading")
	if err := validateOutsideBytes(raw); err != nil {
		t.Fatal(err)
	}
	contained := string(containedOutsideBytes(raw, "\n"))
	if !strings.HasPrefix(contained, "> `````\n") || !strings.Contains(contained, "> ## heading\n") {
		t.Fatalf("contained = %q", contained)
	}
	if err := validateOutsideBytes([]byte{'x', 0}); err == nil {
		t.Fatal("control byte accepted")
	}
	if err := validateOutsideBytes([]byte{'x', '\r'}); err == nil {
		t.Fatal("carriage-return control byte accepted")
	}
}
