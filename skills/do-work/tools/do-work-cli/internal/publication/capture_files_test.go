package publication

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCaptureCollidesWithQueueKanbanReservationForSameNumber(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "do-work/queue/REQ-481-existing.md", []byte("---\nid: REQ-481\nstatus: pending\n---\n"), 0o644)
	boardModule := filepath.Clean("../../../../../do-work-board/tools/queue-kanban")
	command := exec.Command("go", "run", ".", "next-req", "--repo-root", root)
	command.Dir = boardModule
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("queue-kanban next-req: %v\n%s", err, output)
	}
	if number := strings.TrimSpace(string(output)); number != "482" {
		t.Fatalf("next-req output = %q, want 482", number)
	}

	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-482"}), 0o644)
	writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-482", "UR-1"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{
		UserRequestID: "UR-1", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests: []CaptureRequest{{ID: "REQ-482", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-482-new.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-482"}},
	}})
	if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-COLLISION" {
		t.Fatalf("capture did not collide with board reservation: %#v", plan.Refusal)
	}
}

func TestCaptureLegacyReservationAliasBlocksCanonicalCreate(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "do-work/.req-reservations/REQ-000482", nil, 0o644)
	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-482"}), 0o644)
	writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-482", "UR-1"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{
		UserRequestID: "UR-1", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests: []CaptureRequest{{ID: "REQ-482", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-482-new.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-482"}},
	}})
	if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-COLLISION" || !reflect.DeepEqual(plan.Refusal.Paths, []string{"do-work/.req-reservations/REQ-000482"}) {
		t.Fatalf("legacy alias did not block capture: %#v", plan.Refusal)
	}
}

func TestCaptureRefusesFixedSixReservationManifestPath(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-482"}), 0o644)
	writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-482", "UR-1"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{
		UserRequestID: "UR-1", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests: []CaptureRequest{{ID: "REQ-482", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-482-new.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-000482"}},
	}})
	if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-RESERVATION-MISMATCH" {
		t.Fatalf("fixed-six manifest path accepted: %#v", plan.Refusal)
	}
}

func TestCaptureIgnoresWhitespaceWrappedReservationLikeNames(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "do-work/.req-reservations/ REQ-482 ", nil, 0o644)
	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-482"}), 0o644)
	writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-482", "UR-1"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{
		UserRequestID: "UR-1", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests: []CaptureRequest{{ID: "REQ-482", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-482-new.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-482"}},
	}})
	if plan.Refusal != nil {
		t.Fatalf("whitespace-wrapped filename gained reservation authority: %#v", plan.Refusal)
	}
}

func TestBuildCapturePlanBootstrapsAbsentDoWorkAndPinsRawInput(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeFixture(t, repositoryRoot, "payload/raw.txt", []byte("hello\n````\n## outside"), 0o600)
	raw, _ := os.ReadFile(filepath.Join(repositoryRoot, "payload/raw.txt"))
	ur := append(canonicalURFixture("UR-1", []string{"REQ-1"}), containedOutsideBytes(raw, "\n")...)
	ur = append(ur, '\n')
	writeFixture(t, repositoryRoot, "payload/ur.md", ur, 0o644)
	writeFixture(t, repositoryRoot, "payload/req.md", canonicalREQFixture("REQ-1", "UR-1"), 0o644)
	manifest := Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
		UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}}, RawInput: &PayloadFile{SourcePath: "payload/raw.txt"},
		Requests: []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-1-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-001"}},
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
	writeFixture(t, repositoryRoot, "payload/ur.md", canonicalURFixture("UR-1", nil), 0o644)
	plan := BuildCapturePlan(repositoryRoot, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1", UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}}}})
	if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-TOPOLOGY-UNSAFE" {
		t.Fatalf("refusal = %#v", plan.Refusal)
	}
}

func TestBuildCapturePlanRefusesNoncanonicalCaptureSchema(t *testing.T) {
	validUR := canonicalURFixture("UR-1", []string{"REQ-1"})
	validREQ := canonicalREQFixture("REQ-1", "UR-1")
	testCases := []struct {
		name     string
		urBytes  []byte
		reqBytes []byte
		wantCode string
	}{
		{name: "UR missing required title", urBytes: bytesReplace(validUR, "title: 'Captured request'\n", ""), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR requests scalar", urBytes: bytesReplace(validUR, "requests: [REQ-1]", "requests: REQ-1"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR noncanonical created timestamp", urBytes: bytesReplace(validUR, "2026-09-04T12:00:00Z", "2026-09-04T15:00:00+03:00"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR unsafe title scalar", urBytes: bytesReplace(validUR, "title: 'Captured request'", "title: Captured: request"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR duplicate membership", urBytes: bytesReplace(validUR, "requests: [REQ-1]", "requests: [REQ-1, REQ-1]"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR extra membership", urBytes: bytesReplace(validUR, "requests: [REQ-1]", "requests: [REQ-1, REQ-2]"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR missing membership keeps specific finding", urBytes: bytesReplace(validUR, "requests: [REQ-1]", "requests: []"), reqBytes: validREQ, wantCode: "CAPTURE-UR-MEMBERSHIP-MISSING"},
		{name: "UR negative count", urBytes: bytesReplace(validUR, "word_count: 2", "word_count: -1"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR fractional count", urBytes: bytesReplace(validUR, "word_count: 2", "word_count: 1.5"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "UR duplicate field", urBytes: bytesReplace(validUR, "word_count: 2", "word_count: 2\nword_count: 3"), reqBytes: validREQ, wantCode: "CAPTURE-UR-SCHEMA-INVALID"},
		{name: "REQ missing required domain", urBytes: validUR, reqBytes: bytesReplace(validREQ, "domain: general\n", ""), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ read alias is not writable", urBytes: validUR, reqBytes: bytesReplace(validREQ, "depends_on: []", "dependencies: []"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ list field is scalar", urBytes: validUR, reqBytes: bytesReplace(validREQ, "depends_on: []", "depends_on: REQ-2"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ enum alias is not writable", urBytes: validUR, reqBytes: bytesReplace(validREQ, "domain: general", "domain: back_end"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ boolean alias is not writable", urBytes: validUR, reqBytes: bytesReplace(validREQ, "tdd: true", "tdd: yes"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ maintenance alias", urBytes: validUR, reqBytes: bytesReplace(validREQ, "maintenance: false", "maintenance: on"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ effort alias", urBytes: validUR, reqBytes: bytesReplace(validREQ, "maintenance: false", "maintenance: false\neffort_estimate: trivial"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ route case", urBytes: validUR, reqBytes: bytesReplace(validREQ, "maintenance: false", "maintenance: false\nroute: b"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ noncanonical created timestamp", urBytes: validUR, reqBytes: bytesReplace(validREQ, "2026-09-04T12:00:00Z", "2026-09-04"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ unsafe title scalar", urBytes: validUR, reqBytes: bytesReplace(validREQ, "title: 'Work title'", "title: \"Work title\""), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ blocked fields missing", urBytes: validUR, reqBytes: bytesReplace(validREQ, "status: pending", "status: blocked"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ blocked fields on pending", urBytes: validUR, reqBytes: bytesReplace(validREQ, "status: pending", "status: pending\nblocked_by: 'service unavailable'\nblocked_at: 2026-09-04T12:00:00Z"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
		{name: "REQ impact title mismatch", urBytes: validUR, reqBytes: bytesReplace(validREQ, "maintenance: false", "maintenance: false\nimpact: impact-rule-change"), wantCode: "CAPTURE-REQ-SCHEMA-INVALID"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "payload/ur.md", testCase.urBytes, 0o644)
			writeFixture(t, root, "payload/req.md", testCase.reqBytes, 0o644)
			plan := BuildCapturePlan(root, canonicalCaptureManifest("UR-1", "REQ-1"))
			if plan.Refusal == nil || plan.Refusal.Code != testCase.wantCode {
				t.Fatalf("refusal = %#v, want %s", plan.Refusal, testCase.wantCode)
			}
			if len(plan.Mutations) != 0 || len(plan.TargetPaths) != 0 {
				t.Fatalf("schema refusal leaked a partial mutation plan: %#v", plan)
			}
		})
	}
}

func TestBuildCapturePlanAcceptsReorderedExactMembership(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-2", "REQ-1"}), 0o644)
	writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-1", "UR-1"), 0o644)
	writeFixture(t, root, "payload/second.md", canonicalREQFixture("REQ-2", "UR-1"), 0o644)
	manifest := canonicalCaptureManifest("UR-1", "REQ-1")
	manifest.Capture.Requests = append(manifest.Capture.Requests, CaptureRequest{ID: "REQ-2", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-2-second.md", Payload: PayloadFile{SourcePath: "payload/second.md"}}, ReservationPath: "do-work/.req-reservations/REQ-002"})
	if plan := BuildCapturePlan(root, manifest); plan.Refusal != nil {
		t.Fatalf("reordered exact membership refused: %#v", plan.Refusal)
	}
}

func TestBuildCapturePlanEnforcesAuthoredFieldsWithoutChangingJudgments(t *testing.T) {
	validREQ := canonicalREQFixture("REQ-1", "UR-1")
	testCases := []struct {
		name     string
		reqBytes []byte
		accepted bool
	}{
		{"quoted apostrophe", bytesReplace(validREQ, "'Work title'", "'User''s title'"), true},
		{"literal title", bytesReplace(validREQ, "title: 'Work title'", "title: |-\n  Work title\n  Continued"), true},
		{"bare apostrophe", bytesReplace(validREQ, "'Work title'", "'User's title'"), false},
		{"plain title", bytesReplace(validREQ, "'Work title'", "Work title"), false},
		{"empty title", bytesReplace(validREQ, "'Work title'", "''"), false},
		{"duplicate title", bytesReplace(validREQ, "title: 'Work title'", "title: 'First'\ntitle: 'Work title'"), false},
		{"blocked without optional probe", bytesReplace(validREQ, "status: pending", "status: blocked\nblocked_by: 'API unavailable'\nblocked_at: 2026-09-04T12:00:00Z"), true},
		{"pending answers", bytesReplace(validREQ, "status: pending", "status: pending-answers"), true},
		{"blocked unsafe condition", bytesReplace(validREQ, "status: pending", "status: blocked\nblocked_by: API unavailable\nblocked_at: 2026-09-04T12:00:00Z"), false},
		{"blocked missing time", bytesReplace(validREQ, "status: pending", "status: blocked\nblocked_by: 'API unavailable'"), false},
		{"safe assignment", bytesReplace(validREQ, "maintenance: false", "maintenance: false\nassigned_to: 'Ada'"), true},
		{"unsafe assignment", bytesReplace(validREQ, "maintenance: false", "maintenance: false\nassigned_to: Ada"), false},
		{"impact mirror", bytesReplace(bytesReplace(validREQ, "maintenance: false", "maintenance: false\nimpact: impact-rule-change"), "'Work title'", "'[impact-rule-change] Work title'"), true},
		{"default impact", bytesReplace(validREQ, "maintenance: false", "maintenance: false\nimpact: impact-user-visible"), true},
		{"default impact tag", bytesReplace(validREQ, "'Work title'", "'[impact-user-visible] Work title'"), false},
		{"scalar addendum", bytesReplace(validREQ, "maintenance: false", "maintenance: false\naddendum_to: REQ-2"), true},
		{"list addendum", bytesReplace(validREQ, "maintenance: false", "maintenance: false\naddendum_to: [REQ-2]"), false},
	}
	for _, aliasKey := range []string{"amends", "parent", "amendment_to", "dependencies", "batch_name", "related_reqs", "spec_hint", "suggested-spec"} {
		testCases = append(testCases, struct {
			name     string
			reqBytes []byte
			accepted bool
		}{"alias " + aliasKey, bytesReplace(validREQ, "maintenance: false", "maintenance: false\n"+aliasKey+": REQ-2"), false})
	}
	for _, fieldName := range []string{"prime_files", "required_lessons", "depends_on", "related", "write_set"} {
		baseBytes := bytesReplace(validREQ, fieldName+": []\n", "")
		for _, fieldValue := range []string{"REQ-2", "[]", "[REQ-2]", "\n  - REQ-2"} {
			testCases = append(testCases, struct {
				name     string
				reqBytes []byte
				accepted bool
			}{fieldName + " " + fieldValue, bytesReplace(baseBytes, "maintenance: false", "maintenance: false\n"+fieldName+": "+fieldValue), fieldValue != "REQ-2"})
		}
	}
	for _, timestamp := range []string{"2026-09-04", "2026-09-04 12:00:00Z", "2026-09-04T15:00:00+03:00", "2026-09-04T12:00:00.1Z", "2026-02-30T12:00:00Z", "bad"} {
		testCases = append(testCases, struct {
			name     string
			reqBytes []byte
			accepted bool
		}{"timestamp " + timestamp, bytesReplace(validREQ, "2026-09-04T12:00:00Z", timestamp), false})
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-1"}), 0o644)
			writeFixture(t, root, "payload/req.md", testCase.reqBytes, 0o644)
			plan := BuildCapturePlan(root, canonicalCaptureManifest("UR-1", "REQ-1"))
			if testCase.accepted {
				if plan.Refusal != nil {
					t.Fatalf("valid authored record refused: %#v", plan.Refusal)
				}
			} else if plan.Refusal == nil || plan.Refusal.Code != "CAPTURE-REQ-SCHEMA-INVALID" {
				t.Fatalf("invalid authored record refusal = %#v", plan.Refusal)
			}
			if _, err := os.Stat(filepath.Join(root, "do-work")); !os.IsNotExist(err) {
				t.Fatalf("planning wrote queue state: %v", err)
			}
		})
	}
}

func TestBuildCapturePlanAcceptsPublishedCaptureExamples(t *testing.T) {
	referenceBytes, err := os.ReadFile("../../../../actions/capture-reference.md")
	if err != nil {
		t.Fatal(err)
	}
	exampleFor := func(heading string) []byte {
		t.Helper()
		_, section, found := strings.Cut(string(referenceBytes), heading+"\n")
		if !found {
			t.Fatalf("missing example %q", heading)
		}
		_, fencedBody, found := strings.Cut(section, "```markdown\n")
		if !found {
			t.Fatalf("missing example fence %q", heading)
		}
		example, _, found := strings.Cut(fencedBody, "\n```\n")
		if !found {
			t.Fatalf("unclosed example %q", heading)
		}
		return []byte(example + "\n")
	}
	simpleREQ := exampleFor("### Simple REQ")
	complexREQ := append(append([]byte(nil), simpleREQ...), exampleFor("### Complex REQ (additional sections)")...)
	urExample := exampleFor("### UR input.md")
	for _, testCase := range []struct {
		name                     string
		reqBytes, urBytes        []byte
		requestID, userRequestID string
	}{
		{"simple", simpleREQ, canonicalURFixture("UR-001", []string{"REQ-001"}), "REQ-001", "UR-001"},
		{"complex", complexREQ, canonicalURFixture("UR-001", []string{"REQ-001"}), "REQ-001", "UR-001"},
		{"UR input", canonicalREQFixture("REQ-020", "UR-005"), urExample, "REQ-020", "UR-005"},
		{"addendum", exampleFor("## Addendum REQ Template"), canonicalURFixture("UR-006", []string{"REQ-021"}), "REQ-021", "UR-006"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "payload/ur.md", testCase.urBytes, 0o644)
			writeFixture(t, root, "payload/req.md", testCase.reqBytes, 0o644)
			if plan := BuildCapturePlan(root, canonicalCaptureManifest(testCase.userRequestID, testCase.requestID)); plan.Refusal != nil {
				t.Fatalf("copyable example refused: %#v", plan.Refusal)
			}
		})
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
			writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-1"}), 0o644)
			writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-1", "UR-1"), 0o644)
			plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
				UserRequest: PublishedFile{Path: test.urPath, Payload: PayloadFile{SourcePath: "payload/ur.md"}},
				Requests:    []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: test.reqPath, Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-001"}},
			}})
			if plan.Refusal == nil {
				t.Fatalf("accepted noncanonical paths: %#v", plan)
			}
		})
	}
}

func TestRemediationMinorCaptureMutationOrderIsMarkerURAssetsREQFold(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-1"}), 0o644)
	writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-1", "UR-1"), 0o644)
	writeFixture(t, root, "payload/asset", []byte("asset"), 0o644)
	writeFixture(t, root, "do-work/prose-backlog.md", []byte("old\n"), 0o644)
	writeFixture(t, root, "payload/fold-old", []byte("old\n"), 0o644)
	writeFixture(t, root, "payload/fold-new", []byte("new\n"), 0o644)
	plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
		UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests:    []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-1-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-001"}},
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
	want := []string{"do-work/.req-reservations/REQ-001", "do-work/user-requests/UR-1/input.md", "do-work/user-requests/UR-1/assets/a.bin", "do-work/queue/REQ-1-work.md", "do-work/prose-backlog.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mutation order = %v, want %v", got, want)
	}
}

func TestRemediationF9CaptureRefusesCollisionAtEveryCreatePosition(t *testing.T) {
	for _, collisionPath := range []string{
		"do-work/.req-reservations/REQ-001",
		"do-work/user-requests/UR-1/input.md",
		"do-work/user-requests/UR-1/assets/a.bin",
		"do-work/queue/REQ-1-work.md",
	} {
		t.Run(collisionPath, func(t *testing.T) {
			root := t.TempDir()
			writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", []string{"REQ-1"}), 0o644)
			writeFixture(t, root, "payload/req.md", canonicalREQFixture("REQ-1", "UR-1"), 0o644)
			writeFixture(t, root, "payload/asset", []byte("asset"), 0o644)
			writeFixture(t, root, collisionPath, []byte("collision"), 0o644)
			plan := BuildCapturePlan(root, Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{UserRequestID: "UR-1",
				UserRequest: PublishedFile{Path: "do-work/user-requests/UR-1/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
				Requests:    []CaptureRequest{{ID: "REQ-1", UserRequestID: "UR-1", File: PublishedFile{Path: "do-work/queue/REQ-1-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: "do-work/.req-reservations/REQ-001"}},
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
	writeFixture(t, root, "payload/ur.md", canonicalURFixture("UR-1", nil), 0o644)
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

func canonicalCaptureManifest(userRequestID, requestID string) Manifest {
	reservationPath, _ := canonicalReservationPath(requestID)
	return Manifest{Operation: OperationCaptureFiles, Capture: &CaptureManifest{
		UserRequestID: userRequestID,
		UserRequest:   PublishedFile{Path: "do-work/user-requests/" + userRequestID + "/input.md", Payload: PayloadFile{SourcePath: "payload/ur.md"}},
		Requests:      []CaptureRequest{{ID: requestID, UserRequestID: userRequestID, File: PublishedFile{Path: "do-work/queue/" + requestID + "-work.md", Payload: PayloadFile{SourcePath: "payload/req.md"}}, ReservationPath: reservationPath}},
	}}
}

func canonicalURFixture(userRequestID string, requestIDs []string) []byte {
	return []byte("---\n" +
		"id: " + userRequestID + "\n" +
		"title: 'Captured request'\n" +
		"created_at: 2026-09-04T12:00:00Z\n" +
		"requests: [" + strings.Join(requestIDs, ", ") + "]\n" +
		"word_count: 2\n" +
		"---\n\n# Captured request\n")
}

func canonicalREQFixture(requestID, userRequestID string) []byte {
	return []byte("---\n" +
		"id: " + requestID + "\n" +
		"title: 'Work title'\n" +
		"status: pending\n" +
		"created_at: 2026-09-04T12:00:00Z\n" +
		"user_request: " + userRequestID + "\n" +
		"domain: general\n" +
		"prime_files: []\n" +
		"tdd: true\n" +
		"depends_on: []\n" +
		"maintenance: false\n" +
		"---\n\n# Work title\n")
}

func bytesReplace(source []byte, oldText, newText string) []byte {
	return []byte(strings.Replace(string(source), oldText, newText, 1))
}
