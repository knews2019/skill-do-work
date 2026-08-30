package doctor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
)

func TestTimestampAuditUsesEffectiveTopLevelFieldAndRefusesUnsupportedRepairShapes(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	contents := "---\r\nid: REQ-020\r\ncreated_at: 2026-08-30T19:00:00Z\r\n  calculated_at: 2099-01-01T00:00:00Z\r\ncompleted_at: 2099-01-01T00:00:00.500Z # keep\r\ncreated_at: 2026-08-30T19:30:00Z # effective\r\n---\r\nBody\r\n"
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-020-time.md", contents)
	commitDoctorFixture(t, repositoryRoot, "fixture")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plans, findings := BuildTimestampPlan(context.Background(), snapshot, time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC))
	if len(plans) != 0 {
		t.Fatalf("unsupported fraction was planned: %#v", plans)
	}
	foundRefusal := false
	for _, finding := range findings {
		if finding.Code == "TIMESTAMP-REPAIR-REFUSED" && strings.Contains(strings.Join(finding.Evidence, " "), "line 5") {
			foundRefusal = true
		}
	}
	if !foundRefusal {
		t.Fatalf("missing source-line refusal: %#v", findings)
	}
	if got, err := os.ReadFile(filepath.Join(repositoryRoot, "do-work/archive/REQ-020-time.md")); err != nil || string(got) != contents {
		t.Fatalf("audit changed bytes: err=%v bytes=%q", err, got)
	}
}

func TestTimestampRepairAcceptsQuotedASCIIInnerPaddingOnly(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-021-single.md", doctorRequest("REQ-021", "completed", "created_at: ' 2099-01-01T00:00:00Z '\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-022-double.md", doctorRequest("REQ-022", "completed", "created_at: \"\t2099-01-01 00:00:00 \"\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-023-offset.md", doctorRequest("REQ-023", "completed", "created_at: ' 2099-01-01T00:00:00+01:00 '\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-024-fraction.md", doctorRequest("REQ-024", "completed", "created_at: \" 2099-01-01T00:00:00.5Z \"\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-025-unicode.md", doctorRequest("REQ-025", "completed", "created_at: \"\u00a02099-01-01T00:00:00Z\u00a0\"\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-026-nested.md", doctorRequest("REQ-026", "completed", "created_at: \"'2099-01-01T00:00:00Z'\"\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-027-unclosed.md", doctorRequest("REQ-027", "completed", "created_at: '2099-01-01T00:00:00Z\n", "Body"))
	commitDoctorFixture(t, repositoryRoot, "fixture")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	plans, findings := BuildTimestampPlan(context.Background(), snapshot, time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC))
	planned := map[string]bool{}
	for _, plan := range plans {
		planned[plan.RelativePath] = true
	}
	for _, path := range []string{"do-work/archive/REQ-021-single.md", "do-work/archive/REQ-022-double.md"} {
		if !planned[path] {
			t.Fatalf("quoted ASCII-padded timestamp was not planned: %s, plans=%#v", path, plans)
		}
	}
	refused := map[string]bool{}
	for _, finding := range findings {
		if finding.Code == "TIMESTAMP-REPAIR-REFUSED" && len(finding.AffectedPaths) == 1 {
			refused[finding.AffectedPaths[0]] = true
		}
	}
	for _, path := range []string{"do-work/archive/REQ-023-offset.md", "do-work/archive/REQ-024-fraction.md", "do-work/archive/REQ-025-unicode.md", "do-work/archive/REQ-026-nested.md", "do-work/archive/REQ-027-unclosed.md"} {
		if planned[path] {
			t.Fatalf("unsupported timestamp shape was planned: %s", path)
		}
		if !refused[path] {
			t.Fatalf("unsupported timestamp shape was not safely refused: %s, findings=%#v", path, findings)
		}
	}
}
