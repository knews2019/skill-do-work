package toolboxcommands

import "testing"

func TestAuditPercentilesAndBands(t *testing.T) {
	s := auditSummary([]int{1, 2, 3, 4, 5})
	if s.median != 3 || s.p90 != 5 || s.p95 != 5 || s.max != 5 {
		t.Fatalf("summary=%+v", s)
	}
	if auditBand(10, 10, 20) != "" || auditBand(11, 10, 20) != "WATCH" || auditBand(21, 10, 20) != "FLAG" {
		t.Fatal("strict band edge regressed")
	}
}

func TestParseAuditOptions(t *testing.T) {
	o, err := parseAuditOptions(".", []string{"inventory", "--exclude-path", "CHANGELOG.md", "--top-count=3"})
	if err != nil || o.top != 3 || len(o.excludes) != 1 {
		t.Fatalf("options=%+v err=%v", o, err)
	}
	if _, err := parseAuditOptions(".", []string{"inventory", "leftover"}); err == nil {
		t.Fatal("leftover accepted")
	}
}
