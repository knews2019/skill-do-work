package corehelpers

import "testing"

func TestScopeParserUsesOnlyFirstBacktickedPath(t *testing.T) {
	contents := "## Scope\n- `source.go` (modified) — keeps `flex-wrap` behavior\n\n## Implementation Summary\n- `source.go` (modified)\n"
	declared, found, err := firstBacktickedPaths(contents, "Scope", true)
	if err != nil || !found {
		t.Fatalf("parse: %v found=%v", err, found)
	}
	if len(declared) != 1 || declared[0] != "source.go" {
		t.Fatalf("declared=%q", declared)
	}
}
func TestAssociationParserRetainsEveryClosedToken(t *testing.T) {
	contents := "## Implementation Summary\n- `one.go` and `two.go` (modified)\n"
	paths, _, err := allBacktickedPaths(contents, "Implementation Summary")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "one.go" || paths[1] != "two.go" {
		t.Fatalf("paths=%q", paths)
	}
}

func TestScopeDriftReadsAnnotatedScopeAndEverySummaryPath(t *testing.T) {
	contents := "## Scope (Files I will touch)\n- `one.go` — keeps `prose-token`\n- `two.go`\n\n## Implementation Summary\n- `one.go` and `two.go` (modified)\n"
	declared, found, err := firstBacktickedPaths(contents, "Scope", true)
	if err != nil || !found || len(declared) != 2 {
		t.Fatalf("declared=%q found=%v err=%v", declared, found, err)
	}
	implemented, found, err := allBacktickedPaths(contents, "Implementation Summary")
	if err != nil || !found || len(implemented) != 2 {
		t.Fatalf("implemented=%q found=%v err=%v", implemented, found, err)
	}
}
