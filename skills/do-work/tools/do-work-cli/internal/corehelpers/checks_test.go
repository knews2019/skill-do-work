package corehelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestQualificationArtifactDetectorFindsEveryMarkerWithoutMatchingItsSource(t *testing.T) {
	markers := []string{
		strings.Join([]string{"de", "bugger"}, ""),
		strings.Join([]string{"TO", "DO"}, ""),
		strings.Join([]string{"FIX", "ME"}, ""),
	}
	for _, marker := range markers {
		if !qualificationDebugArtifactPattern.MatchString("prefix " + marker + " suffix") {
			t.Errorf("marker %q was not detected", marker)
		}
	}
	source, err := os.ReadFile("checks.go")
	if err != nil {
		t.Fatal(err)
	}
	if qualificationDebugArtifactPattern.Match(source) {
		t.Fatal("artifact detector matches its own implementation source")
	}
}

func TestPreflightReportsMissingNodeAndPythonDependencies(t *testing.T) {
	repository := t.TempDir()
	if err := os.WriteFile(filepath.Join(repository, "package.json"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "requirements.txt"), []byte("fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("VIRTUAL_ENV", "")
	result := handlePreflight(testContext(repository), nil)
	if result.Outcome != "success" || !hasFinding(result, "PREFLIGHT-NODE-MODULES-MISSING") || !hasFinding(result, "PREFLIGHT-VIRTUALENV-MISSING") {
		t.Fatalf("result=%#v", result)
	}
}

func TestQualificationTreatsMovedMarkerAsWarningAndOrdersPaths(t *testing.T) {
	repository := newGitFixture(t)
	marker := strings.Join([]string{"TO", "DO"}, "")
	sourcePath := filepath.Join(repository, "artifact.go")
	if err := os.WriteFile(sourcePath, []byte(marker+"\nclean\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runFixtureGitCommand(t, repository, "add", "artifact.go")
	runFixtureGitCommand(t, repository, "commit", "-qm", "artifact baseline")
	if err := os.WriteFile(sourcePath, []byte("clean\n"+marker+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	requestPath := writeQualificationRequest(t, repository, []string{"artifact.go"})
	result := handleQualify(testContext(repository), []string{"--request-path", requestPath})
	if result.Outcome != "success" || !hasFinding(result, "QUALIFY-DEBUG-ARTIFACT-RELOCATED") || hasFinding(result, "QUALIFY-DEBUG-ARTIFACT") {
		t.Fatalf("relocation result=%#v", result)
	}

	for _, path := range []string{"zeta.go", "alpha.go"} {
		if err := os.WriteFile(filepath.Join(repository, path), []byte(marker+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	requestPath = writeQualificationRequest(t, repository, []string{"zeta.go", "alpha.go"})
	result = handleQualify(testContext(repository), []string{"--request-path", requestPath})
	artifactPaths := []string{}
	for _, finding := range result.Findings {
		if finding.Code == "QUALIFY-DEBUG-ARTIFACT" {
			artifactPaths = append(artifactPaths, finding.AffectedPaths[0])
		}
	}
	if strings.Join(artifactPaths, ",") != "alpha.go,zeta.go" {
		t.Fatalf("artifact order=%q result=%#v", artifactPaths, result)
	}
}

func TestQualificationReporterLibraryRenameAndBinaryBoundaries(t *testing.T) {
	repository := newGitFixture(t)
	outputCall := strings.Join([]string{"console.", "log('x')"}, "")
	if err := os.WriteFile(filepath.Join(repository, "library.js"), []byte("export const value = 1;\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runFixtureGitCommand(t, repository, "add", "library.js")
	runFixtureGitCommand(t, repository, "commit", "-qm", "library baseline")
	runFixtureGitCommand(t, repository, "mv", "library.js", "renamed.js")
	if err := os.WriteFile(filepath.Join(repository, "renamed.js"), []byte("process.exit(0);\n"+outputCall+";\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	reporter := "reporter.py"
	if err := os.WriteFile(filepath.Join(repository, reporter), []byte("sys.exit(0)\n"+strings.Join([]string{"pri", "nt('ok')\n"}, "")), 0o600); err != nil {
		t.Fatal(err)
	}
	binary := "fixture.bin"
	marker := strings.Join([]string{"FIX", "ME"}, "")
	if err := os.WriteFile(filepath.Join(repository, binary), append([]byte{0}, []byte(marker)...), 0o600); err != nil {
		t.Fatal(err)
	}
	requestPath := writeQualificationRequest(t, repository, []string{"renamed.js", reporter, binary})
	result := handleQualify(testContext(repository), []string{"--request-path", requestPath})
	if !hasFinding(result, "QUALIFY-LIBRARY-OUTPUT") || !hasFinding(result, "QUALIFY-REPORTER-OUTPUT") {
		t.Fatalf("output boundaries=%#v", result)
	}
	for _, finding := range result.Findings {
		if finding.Code == "QUALIFY-DEBUG-ARTIFACT" && len(finding.AffectedPaths) > 0 && finding.AffectedPaths[0] == binary {
			t.Fatalf("binary content was scanned: %#v", finding)
		}
	}
}

func writeQualificationRequest(t *testing.T, repository string, paths []string) string {
	t.Helper()
	relative := "do-work/working/REQ-999-qualification.md"
	absolute := filepath.Join(repository, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	var summary strings.Builder
	for _, path := range paths {
		fmt.Fprintf(&summary, "- `%s` (modified)\n", path)
	}
	contents := "---\nid: REQ-999\nstatus: claimed\n---\n\n## AI Execution State (P-A-U Loop)\n- [x] **[PLAN]:** done\n- [x] **[APPLY]:** done\n- [x] **[UNIFY]:** done\n\n## Implementation Summary\n" + summary.String()
	if err := os.WriteFile(absolute, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return relative
}
