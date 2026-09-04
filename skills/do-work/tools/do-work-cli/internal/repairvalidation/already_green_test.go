package repairvalidation

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/gateevidence"
)

const testGateFingerprint = "sha256:already-green-fixture"

func TestValidateAlreadyGreenRepairUsesIntakeAndRecordedEvidence(t *testing.T) {
	t.Parallel()
	repositoryRoot, requestPath := alreadyGreenRepository(t)
	validation, evidence, err := Validate(repositoryRoot, Options{
		RequestPath: requestPath,
		WriterLabel: "fixture:/repo",
		Now:         time.Date(2026, 9, 2, 5, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !validation.TDDAllowed || !validation.ReviewAllowed {
		t.Fatalf("canonical validation = %#v", validation)
	}
	if validation.IntakeFingerprint != testGateFingerprint || validation.ExpectedFingerprint != testGateFingerprint {
		t.Fatalf("fingerprint projection = %#v", validation)
	}
	if evidence == nil || !evidence.Matches || evidence.TargetRevision != validation.RecordedRevision {
		t.Fatalf("recorded evidence = %#v, validation = %#v", evidence, validation)
	}
	for _, required := range []string{requestPath, "do-work/archive/" + filepath.Base(requestPath)} {
		if !contains(validation.CanonicalCompletionPaths, required) {
			t.Fatalf("canonical paths %v missing %s", validation.CanonicalCompletionPaths, required)
		}
	}
}

func TestValidateAlreadyGreenRepairRejectsCoordinatedCurrentFingerprintMutation(t *testing.T) {
	t.Parallel()
	repositoryRoot, requestPath := alreadyGreenRepository(t)
	absolute := filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))
	contents, err := os.ReadFile(absolute)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.ReplaceAll(string(contents), testGateFingerprint, "sha256:coordinated-self-assertion")
	if err := os.WriteFile(absolute, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}

	validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if validation.TDDAllowed || validation.ReviewAllowed || !contains(validation.ReasonCodes, "REPAIR-INTAKE-NOT-DURABLE") {
		t.Fatalf("coordinated fingerprint mutation escaped durable intake proof: %#v", validation)
	}
}

func TestValidateAlreadyGreenRepairRejectsGreenRecordPredatingIntake(t *testing.T) {
	t.Parallel()
	repositoryRoot := t.TempDir()
	git(t, repositoryRoot, "init", "-q")
	git(t, repositoryRoot, "config", "user.name", "repair validator fixture")
	git(t, repositoryRoot, "config", "user.email", "repair-validator@example.invalid")
	writeFile(t, repositoryRoot, "seed.txt", "before intake\n")
	git(t, repositoryRoot, "add", "seed.txt")
	git(t, repositoryRoot, "commit", "-qm", "pre-intake baseline")
	evidence, err := gateevidence.RecordGreenGate(repositoryRoot, []string{"sh", "-c", "exit 0"})
	if err != nil {
		t.Fatalf("RecordGreenGate: %v", err)
	}
	requestPath := "do-work/working/REQ-701-already-green-repair.md"
	writeFile(t, repositoryRoot, requestPath, alreadyGreenRequestContents(evidence.RecordedRevision))
	git(t, repositoryRoot, "add", requestPath)
	git(t, repositoryRoot, "commit", "-qm", "intake and no-op after stale evidence")

	validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if validation.TDDAllowed || validation.ReviewAllowed || !contains(validation.ReasonCodes, "REPAIR-INTAKE-NOT-DURABLE") {
		t.Fatalf("pre-intake green record escaped durable intake proof: %#v", validation)
	}
}

func TestValidateAlreadyGreenRepairRejectsAmbiguousRequestIdentityForBothDecisions(t *testing.T) {
	t.Parallel()
	repositoryRoot, requestPath := alreadyGreenRepository(t)
	writeFile(t, repositoryRoot, "do-work/archive/REQ-701-duplicate.md", "---\nid: REQ-701\ntitle: Duplicate\nstatus: completed\n---\n")

	validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if validation.TDDAllowed || validation.ReviewAllowed || !contains(validation.ReasonCodes, "REPAIR-REQUEST-IDENTITY") {
		t.Fatalf("ambiguous request identity escaped core eligibility: %#v", validation)
	}
}

func TestValidateAlreadyGreenRepairRequiresExactNoOpBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		insertion string
	}{
		{name: "foreign line", insertion: "\nForeign no-op claim.\n"},
		{name: "foreign label", insertion: "\n- **Foreign authority:** self asserted\n"},
		{name: "prefix-like duplicate heading", insertion: "\n## Repository Gate Repair No-Op Extra\n\nSelf asserted.\n"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot, requestPath := alreadyGreenRepository(t)
			absolute := filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))
			contents, err := os.ReadFile(absolute)
			if err != nil {
				t.Fatal(err)
			}
			updated := strings.Replace(string(contents), "\n## Implementation Summary", test.insertion+"\n## Implementation Summary", 1)
			if err := os.WriteFile(absolute, []byte(updated), 0o644); err != nil {
				t.Fatal(err)
			}

			validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if validation.TDDAllowed || validation.ReviewAllowed || !contains(validation.ReasonCodes, "REPAIR-NO-OP-SHAPE") {
				t.Fatalf("additive no-op content escaped exact shape: %#v", validation)
			}
		})
	}
}

func TestValidateAlreadyGreenRepairUsesCanonicalCompletionPathsForStageAuthority(t *testing.T) {
	t.Parallel()
	repositoryRoot, requestPath := alreadyGreenRepository(t)
	requestAbsolute := filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))
	contents, err := os.ReadFile(requestAbsolute)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(requestAbsolute, append(contents, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, repositoryRoot, "add", requestPath)
	validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
	if err != nil {
		t.Fatalf("Validate exact stage: %v", err)
	}
	if !validation.ReviewAllowed {
		t.Fatalf("exact canonical request stage was refused: %#v", validation)
	}

	writeFile(t, repositoryRoot, "do-work/archive/REQ-999-unrelated.md", "unrelated\n")
	git(t, repositoryRoot, "add", "do-work/archive/REQ-999-unrelated.md")

	validation, _, err = Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !validation.TDDAllowed {
		t.Fatalf("stage-only blocker changed TDD decision: %#v", validation)
	}
	if validation.ReviewAllowed || !contains(validation.ReasonCodes, "REPAIR-STAGED-PATH-NOT-CANONICAL") {
		t.Fatalf("unrelated archive escaped exact result authorization: %#v", validation)
	}
}

func TestValidateAlreadyGreenRepairRejectsNeighboringShapes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		old  string
		new  string
	}{
		{"ordinary", "repository_gate_repair: true", "repository_gate_repair: false"},
		{"release metadata", "claimed_at: 2026-09-02T04:00:00Z", "claimed_at: 2026-09-02T04:00:00Z\nrelease_at: 2026-09-02T04:30:00Z"},
		{"fingerprint mismatch", "Expected diagnostic fingerprint:** " + testGateFingerprint, "Expected diagnostic fingerprint:** sha256:self-asserted"},
		{"conflicting folded intake", "## Repository Gate Repair No-Op", "## Repository Gate Repair Intake\n\n- **Parent:** REQ-702\n- **Gate command (argv JSON):** [\"sh\",\"-c\",\"exit 0\"]\n- **Direct exit status:** 1\n- **Diagnostic fingerprint:** sha256:conflict\n- **Repair dependency:** REQ-701\n\n## Repository Gate Repair No-Op"},
		{"gate argv mismatch", `- **Gate command:** ["sh","-c","exit 0"]`, `- **Gate command:** ["sh","-c","true"]`},
		{"missing recorded revision", "- **Recorded green revision:**", "- **Recorded green revision missing:**"},
		{"malformed no-op", "## Repository Gate Repair No-Op", "## Repair Note"},
		{"malformed summary", "None — verified repository-gate repair no-op.", "None."},
		{"malformed qualification", "durable gate evidence verified and project diff empty.", "gate looked green."},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			repositoryRoot, requestPath := alreadyGreenRepository(t)
			absolute := filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))
			contents, err := os.ReadFile(absolute)
			if err != nil {
				t.Fatal(err)
			}
			updated := strings.Replace(string(contents), test.old, test.new, 1)
			if updated == string(contents) {
				t.Fatalf("mutation %q changed nothing", test.name)
			}
			if err := os.WriteFile(absolute, []byte(updated), 0o644); err != nil {
				t.Fatal(err)
			}
			validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if validation.TDDAllowed || validation.ReviewAllowed || len(validation.ReasonCodes) == 0 {
				t.Fatalf("neighbor was accepted: %#v", validation)
			}
		})
	}
}

func TestValidateAlreadyGreenRepairRejectsProjectAndReleaseMutation(t *testing.T) {
	t.Parallel()
	for _, relativePath := range []string{"app.txt", "VERSION", "package-lock.json"} {
		relativePath := relativePath
		t.Run(relativePath, func(t *testing.T) {
			t.Parallel()
			repositoryRoot, requestPath := alreadyGreenRepository(t)
			writeFile(t, repositoryRoot, relativePath, "changed\n")
			validation, _, err := Validate(repositoryRoot, Options{RequestPath: requestPath, WriterLabel: "fixture:/repo"})
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if validation.TDDAllowed || validation.ReviewAllowed || !contains(validation.ReasonCodes, "REPAIR-PROJECT-DIFF-NONEMPTY") || !contains(validation.OffendingPaths, relativePath) {
				t.Fatalf("project/release mutation escaped: %#v", validation)
			}
		})
	}
}

func alreadyGreenRepository(t *testing.T) (string, string) {
	t.Helper()
	repositoryRoot := t.TempDir()
	git(t, repositoryRoot, "init", "-q")
	git(t, repositoryRoot, "config", "user.name", "repair validator fixture")
	git(t, repositoryRoot, "config", "user.email", "repair-validator@example.invalid")
	requestPath := "do-work/working/REQ-701-already-green-repair.md"
	intake := alreadyGreenIntakeContents()
	writeFile(t, repositoryRoot, requestPath, intake)
	git(t, repositoryRoot, "add", requestPath)
	git(t, repositoryRoot, "commit", "-qm", "fixture baseline")
	evidence, err := gateevidence.RecordGreenGate(repositoryRoot, []string{"sh", "-c", "exit 0"})
	if err != nil {
		t.Fatalf("RecordGreenGate: %v", err)
	}
	writeFile(t, repositoryRoot, requestPath, alreadyGreenRequestContents(evidence.RecordedRevision))
	return repositoryRoot, requestPath
}

func alreadyGreenIntakeContents() string {
	return `---
id: REQ-701
title: Already-green repair
status: claimed
route: C
repository_gate_repair: true
tdd: true
claimed_at: 2026-09-02T04:00:00Z
---

# Already-Green Repair

## Repository Gate Repair Intake

- **Parent:** REQ-700
- **Gate command (argv JSON):** ["sh","-c","exit 0"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** ` + testGateFingerprint + `
- **Repair dependency:** REQ-701
`
}

func alreadyGreenRequestContents(recordedRevision string) string {
	return alreadyGreenIntakeContents() + `
## Repository Gate Repair No-Op

- **Expected diagnostic fingerprint:** ` + testGateFingerprint + `
- **Gate command:** ["sh","-c","exit 0"]
- **Direct exit status:** 0
- **Recorded green revision:** ` + recordedRevision + `
- **Observed result:** green before implementation; repair already satisfied
- **Verified at:** 2026-09-02T04:40:00Z

## Implementation Summary

**Files changed:** None — verified repository-gate repair no-op.

**What was done:** Re-ran the repair's recorded canonical repository gate before source edits and confirmed it is already green; no implementation changes were necessary.

## Qualification

Passed — repository-gate repair no-op; durable gate evidence verified and project diff empty.
`
}

func git(t *testing.T, repositoryRoot string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return strings.TrimSpace(string(output))
}

func writeFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolute := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
