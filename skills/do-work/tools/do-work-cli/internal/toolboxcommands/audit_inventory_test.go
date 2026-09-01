package toolboxcommands

import "testing"

func TestAuditAggregatesAndFolders(t *testing.T) {
	files := []auditFile{{path: "a/x.go", lines: 3, words: 4}, {path: "a/y.go", lines: 2, words: 3}, {path: "README", lines: 10, words: 12}}
	a := auditAggregates(files)
	if len(a) != 2 || a[0].extension != "(none)" || a[1].extension != ".go" {
		t.Fatalf("aggregates=%+v", a)
	}
	f := auditFolders(files)
	if len(f) != 2 || f[0].folder != "a" || f[0].files != 2 {
		t.Fatalf("folders=%+v", f)
	}
}
