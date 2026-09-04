package heavyverification

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const heavyTestManifest = `{
  "schema_version": 1,
  "lanes": [
    {
      "id": "browser-behavior",
      "argv": ["bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "browser-behavior"],
      "coverage": [{"kind": "suffix-under", "root": "web", "suffix": ".js"}]
    },
    {
      "id": "update-script",
      "argv": ["bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "update-script"],
      "coverage": [{"kind": "exact", "path": "scripts/update.sh"}]
    }
  ],
  "non_heavy_coverage": [{"kind": "subtree", "path": "docs"}]
}`

func TestPlanSelectsOnlyAffectedLanesAndExplainsThem(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "web/app.js", "changed\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "known change")

	plan, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior"}) {
		t.Fatalf("selected lanes = %v, want browser-behavior only", got)
	}
	if plan.Uncertain || len(plan.UncoveredPaths) != 0 {
		t.Fatalf("known change reported uncertainty: %#v", plan)
	}
	if len(plan.SelectedLanes[0].Reasons) != 1 || !strings.Contains(plan.SelectedLanes[0].Reasons[0], "web/app.js") {
		t.Fatalf("selection reasons = %v, want changed path", plan.SelectedLanes[0].Reasons)
	}
}

func TestPlanSelectsMultipleAffectedLanes(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "web/app.js", "changed\n")
	writeHeavyTestFile(t, repositoryRoot, "scripts/update.sh", "changed\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "two lanes")

	plan, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior", "update-script"}) {
		t.Fatalf("selected lanes = %v", got)
	}
}

func TestPlanMatchesBothRenameEndpoints(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "scripts/update.sh", "same bytes\n")
	commitHeavyTestChanges(t, repositoryRoot, "add source")
	baseRevision = runHeavyTestGit(t, repositoryRoot, "rev-parse", "HEAD")
	if err := os.MkdirAll(filepath.Join(repositoryRoot, "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	runHeavyTestGit(t, repositoryRoot, "mv", "scripts/update.sh", "web/update.js")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "rename source")

	plan, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior", "update-script"}) {
		t.Fatalf("rename selected lanes = %v, want both endpoint owners", got)
	}
	if !reflect.DeepEqual(plan.ChangedPaths, []string{"scripts/update.sh", "web/update.js"}) {
		t.Fatalf("rename paths = %v", plan.ChangedPaths)
	}
}

func TestPlanPreservesNewlineInChangedPath(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	changedPath := "web/line\nbreak.js"
	writeHeavyTestFile(t, repositoryRoot, changedPath, "changed\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "newline path")

	plan, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plan.ChangedPaths, []string{changedPath}) {
		t.Fatalf("changed paths = %q", plan.ChangedPaths)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior"}) {
		t.Fatalf("newline path selected lanes = %v", got)
	}
}

func TestPlanUnknownPathSelectsAllWithExplicitUncertainty(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "mystery.bin", "unknown\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "unknown path")

	plan, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior", "update-script"}) {
		t.Fatalf("fallback selected lanes = %v", got)
	}
	if !plan.Uncertain || !reflect.DeepEqual(plan.UncoveredPaths, []string{"mystery.bin"}) {
		t.Fatalf("unknown path did not remain explicit: %#v", plan)
	}
	for _, lane := range plan.SelectedLanes {
		if len(lane.Reasons) != 1 || !strings.Contains(lane.Reasons[0], "coverage is uncertain") {
			t.Fatalf("fallback reasons for %s = %v", lane.LaneID, lane.Reasons)
		}
	}
}

func TestPlanForceAllAndNoHeavySelection(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "docs/guide.md", "documentation\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "non-heavy path")
	manifestPath := heavyTestManifestPath(repositoryRoot)

	ordinaryPlan, err := Plan(repositoryRoot, manifestPath, baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(ordinaryPlan.SelectedLanes) != 0 || ordinaryPlan.Uncertain {
		t.Fatalf("non-heavy change selected work: %#v", ordinaryPlan)
	}

	forcedPlan, err := Plan(repositoryRoot, manifestPath, baseRevision, targetRevision, true)
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedHeavyLaneIDs(forcedPlan); !reflect.DeepEqual(got, []string{"browser-behavior", "update-script"}) {
		t.Fatalf("force-all selected lanes = %v", got)
	}
	if !forcedPlan.ForcedAll {
		t.Fatalf("force-all plan did not record override: %#v", forcedPlan)
	}
}

func TestPlanRejectsInvalidManifest(t *testing.T) {
	invalidManifest := strings.Replace(heavyTestManifest, `"update-script"`, `"browser-behavior"`, 1)
	repositoryRoot, baseRevision := newHeavyTestRepository(t, invalidManifest)
	_, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, "HEAD", false)
	if err == nil {
		t.Fatal("invalid manifest was accepted")
	}
}

func TestPlanUsesManifestBytesFromTargetRevision(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "web/app.js", "changed\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "known change")
	narrowedManifest := strings.Replace(heavyTestManifest, `"root": "web"`, `"root": "elsewhere"`, 1)
	writeHeavyTestFile(t, repositoryRoot, "heavy-lanes.json", narrowedManifest)

	plan, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior"}) {
		t.Fatalf("dirty manifest narrowed committed selection: %v", got)
	}
}

func TestPlanRevalidationUnionsHistoricalRangesAgainstExecutionManifest(t *testing.T) {
	repositoryRoot := t.TempDir()
	runHeavyTestGit(t, repositoryRoot, "init", "-q")
	runHeavyTestGit(t, repositoryRoot, "config", "user.name", "Heavy Test")
	runHeavyTestGit(t, repositoryRoot, "config", "user.email", "heavy@example.test")
	writeHeavyTestFile(t, repositoryRoot, "seed.txt", "seed\n")
	firstBase := commitHeavyTestChanges(t, repositoryRoot, "seed without manifest")
	writeHeavyTestFile(t, repositoryRoot, "scripts/update.sh", "same bytes\n")
	firstTarget := commitHeavyTestChanges(t, repositoryRoot, "historical updater")
	secondBase := firstTarget
	if err := os.MkdirAll(filepath.Join(repositoryRoot, "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	runHeavyTestGit(t, repositoryRoot, "mv", "scripts/update.sh", "web/update.js")
	secondTarget := commitHeavyTestChanges(t, repositoryRoot, "historical rename")
	writeHeavyTestFile(t, repositoryRoot, "heavy-lanes.json", heavyTestManifest)
	executionRevision := commitHeavyTestChanges(t, repositoryRoot, "add current lane manifest")

	plan, err := PlanRevalidation(repositoryRoot, heavyTestManifestPath(repositoryRoot), []resultmodel.HeavySourceRange{
		{BaseRevision: firstBase, TargetRevision: firstTarget},
		{BaseRevision: secondBase, TargetRevision: secondTarget},
	}, executionRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Mode != "historical-revalidation" || plan.ManifestRevision != executionRevision || plan.ExecutionRevision != executionRevision {
		t.Fatalf("revalidation identity = %#v", plan)
	}
	if !reflect.DeepEqual(plan.ChangedPaths, []string{"scripts/update.sh", "web/update.js"}) {
		t.Fatalf("unioned paths = %v", plan.ChangedPaths)
	}
	if got := selectedHeavyLaneIDs(plan); !reflect.DeepEqual(got, []string{"browser-behavior", "update-script"}) {
		t.Fatalf("selected lanes = %v", got)
	}
	if !reflect.DeepEqual(plan.SourceRanges, []resultmodel.HeavySourceRange{
		{BaseRevision: firstBase, TargetRevision: firstTarget},
		{BaseRevision: secondBase, TargetRevision: secondTarget},
	}) {
		t.Fatalf("source ranges = %#v", plan.SourceRanges)
	}
}

func TestPlanRevalidationRejectsMissingExecutionManifestAndNonAncestorTarget(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "web/app.js", "branch change\n")
	branchTarget := commitHeavyTestChanges(t, repositoryRoot, "side branch target")
	runHeavyTestGit(t, repositoryRoot, "checkout", "-q", "-b", "execution", baseRevision)
	writeHeavyTestFile(t, repositoryRoot, "docs/guide.md", "execution branch\n")
	executionRevision := commitHeavyTestChanges(t, repositoryRoot, "execution branch")

	_, err := PlanRevalidation(repositoryRoot, heavyTestManifestPath(repositoryRoot), []resultmodel.HeavySourceRange{{BaseRevision: baseRevision, TargetRevision: branchTarget}}, executionRevision, false)
	if err == nil || !strings.Contains(err.Error(), "ancestor") {
		t.Fatalf("non-ancestor target error = %v", err)
	}

	if err := os.Remove(heavyTestManifestPath(repositoryRoot)); err != nil {
		t.Fatal(err)
	}
	missingManifestRevision := commitHeavyTestChanges(t, repositoryRoot, "remove manifest")
	_, err = PlanRevalidation(repositoryRoot, heavyTestManifestPath(repositoryRoot), []resultmodel.HeavySourceRange{{BaseRevision: baseRevision, TargetRevision: executionRevision}}, missingManifestRevision, false)
	if err == nil || !strings.Contains(err.Error(), "manifest") {
		t.Fatalf("missing execution manifest error = %v", err)
	}
}

func TestPlanRevalidationUnknownPathSelectsAllAndStrictPlanRemainsExact(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	writeHeavyTestFile(t, repositoryRoot, "mystery.bin", "unknown\n")
	targetRevision := commitHeavyTestChanges(t, repositoryRoot, "unknown historical path")

	revalidation, err := PlanRevalidation(repositoryRoot, heavyTestManifestPath(repositoryRoot), []resultmodel.HeavySourceRange{{BaseRevision: baseRevision, TargetRevision: targetRevision}}, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if !revalidation.Uncertain || !reflect.DeepEqual(selectedHeavyLaneIDs(revalidation), []string{"browser-behavior", "update-script"}) {
		t.Fatalf("revalidation fail-closed plan = %#v", revalidation)
	}
	strict, err := Plan(repositoryRoot, heavyTestManifestPath(repositoryRoot), baseRevision, targetRevision, false)
	if err != nil {
		t.Fatal(err)
	}
	if strict.Mode != "" || len(strict.SourceRanges) != 0 || strict.ManifestRevision != "" || strict.ExecutionRevision != "" {
		t.Fatalf("strict planner contract changed: %#v", strict)
	}
}

func TestPlanUsesTargetRevisionManifestModeAndRejectsTargetSymlink(t *testing.T) {
	repositoryRoot, baseRevision := newHeavyTestRepository(t, heavyTestManifest)
	manifestPath := heavyTestManifestPath(repositoryRoot)
	if err := os.Chmod(manifestPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(repositoryRoot, manifestPath, baseRevision, "HEAD", false); err != nil {
		t.Fatalf("working-copy mode drift overrode regular target metadata: %v", err)
	}
	if err := os.Remove(manifestPath); err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(repositoryRoot, manifestPath, baseRevision, "HEAD", false); err != nil {
		t.Fatalf("missing working copy overrode regular target metadata: %v", err)
	}

	if err := os.Symlink(filepath.Join(repositoryRoot, "seed.txt"), manifestPath); err != nil {
		t.Fatal(err)
	}
	symlinkTarget := commitHeavyTestChanges(t, repositoryRoot, "symlink manifest")
	if err := os.Remove(manifestPath); err != nil {
		t.Fatal(err)
	}
	writeHeavyTestFile(t, repositoryRoot, "heavy-lanes.json", heavyTestManifest)
	if _, err := Plan(repositoryRoot, manifestPath, baseRevision, symlinkTarget, false); err == nil {
		t.Fatal("target-revision symlink manifest was accepted through a regular working copy")
	}
}

func TestPlanRejectsMalformedTrailingJSON(t *testing.T) {
	malformedManifest := heavyTestManifest + " trailing"
	secondRoot, secondBase := newHeavyTestRepository(t, malformedManifest)
	if _, err := Plan(secondRoot, heavyTestManifestPath(secondRoot), secondBase, "HEAD", false); err == nil {
		t.Fatal("malformed trailing JSON was accepted")
	}
}

// TestDecodeManifestRefusesUnusableFingerprintDeclarations keeps a fingerprint
// block that cannot decide reuse out of the manifest, rather than letting it
// decode into a digest that quietly covers less than it claims.
func TestDecodeManifestRefusesUnusableFingerprintDeclarations(t *testing.T) {
	for _, testCase := range []struct{ name, fingerprint string }{
		{"no toolchain probe", `{"toolchain_probes": []}`},
		{"empty probe argv", `{"toolchain_probes": [[]]}`},
		{"empty probe token", `{"toolchain_probes": [["go", ""]]}`},
		{"duplicate environment variable", `{"toolchain_probes": [["go", "version"]], "environment_variables": ["CI", "CI"]}`},
		{"blank environment variable", `{"toolchain_probes": [["go", "version"]], "environment_variables": [" "]}`},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			manifest := `{"schema_version": 1, "lanes": [{"id": "lane", "argv": ["true"], ` +
				`"coverage": [{"kind": "exact", "path": "seed.txt"}], "fingerprint": ` + testCase.fingerprint + `}]}`
			if _, err := decodeManifest([]byte(manifest)); err == nil {
				t.Fatalf("manifest with %s was accepted", testCase.name)
			}
		})
	}
	// The declaration stays optional: a lane without one still decodes, it just
	// never reuses evidence.
	withoutFingerprint := `{"schema_version": 1, "lanes": [{"id": "lane", "argv": ["true"], "coverage": [{"kind": "exact", "path": "seed.txt"}]}]}`
	manifest, err := decodeManifest([]byte(withoutFingerprint))
	if err != nil {
		t.Fatalf("a lane without fingerprint inputs no longer decodes: %v", err)
	}
	if manifest.Lanes[0].Fingerprint != nil {
		t.Fatalf("absent fingerprint decoded into %#v", manifest.Lanes[0].Fingerprint)
	}
}

func TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint(t *testing.T) {
	repositoryRoot, err := filepath.Abs("../../../../../..")
	if err != nil {
		t.Fatal(err)
	}
	manifestContents, err := os.ReadFile(filepath.Join(repositoryRoot, "_dev/tests/heavy-lanes.json"))
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := decodeManifest(manifestContents)
	if err != nil {
		t.Fatal(err)
	}
	wantLaneIDs := []string{
		"queue-kanban-javascript", "queue-kanban-browser", "do-work-cli-integrations",
		"staged-skills", "updater", "installer",
	}
	maintainerScript, err := os.ReadFile(filepath.Join(repositoryRoot, "_dev/tests/maintainer-verify.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Lanes) != len(wantLaneIDs) {
		t.Fatalf("manifest lanes = %d, want %d", len(manifest.Lanes), len(wantLaneIDs))
	}
	for laneIndex, wantLaneID := range wantLaneIDs {
		lane := manifest.Lanes[laneIndex]
		if lane.ID != wantLaneID {
			t.Fatalf("lane %d id = %q, want %q", laneIndex, lane.ID, wantLaneID)
		}
		wantArgv := []string{"env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", wantLaneID}
		if !reflect.DeepEqual(lane.Argv, wantArgv) {
			t.Fatalf("lane %s argv = %v", lane.ID, lane.Argv)
		}
		if !strings.Contains(string(maintainerScript), "    "+wantLaneID+")") {
			t.Fatalf("maintainer dispatcher has no case for %s", wantLaneID)
		}
	}
	if !strings.Contains(string(maintainerScript), "--heavy)\n") {
		t.Fatal("maintainer dispatcher no longer preserves --heavy force-all")
	}
	// Complete lanes retain toolchain evidence. Browser runtime discovery is
	// deliberately uncertain, so that lane must execute.
	for _, lane := range manifest.Lanes {
		if lane.ID == "queue-kanban-browser" {
			continue
		}
		if lane.Fingerprint == nil || len(lane.Fingerprint.ToolchainProbes) == 0 {
			t.Fatalf("lane %s declares no fingerprint toolchain probes, so its evidence can never be reused", lane.ID)
		}
	}

	coverageCases := []struct {
		path        string
		wantLaneIDs []string
	}{
		{"skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", []string{"do-work-cli-integrations", "staged-skills", "updater", "installer"}},
		{"VERSION", []string{"installer"}},
		{"README.md", []string{"installer"}},
		{"suite/modules.tsv", []string{"staged-skills", "updater", "installer"}},
		{"tools/install-do-work-suite.sh", []string{"updater", "installer"}},
		{"tools/validate-suite-manifest.sh", []string{"updater", "installer"}},
		{"tools/replace-text-section.sh", []string{"updater", "installer"}},
		{"tools/fetch-upstream-archive.sh", []string{"updater"}},
		{"skills/do-work/tools/do-work-update.sh", []string{"staged-skills", "updater"}},
		{"skills/do-work/tools/do-work-cli.sh", []string{"staged-skills", "updater"}},
		{"skills/do-work/hooks/hooks.json", []string{"staged-skills", "updater"}},
		{"skills/do-work/agent-instructions.template.md", []string{"staged-skills", "updater", "installer"}},
		{"skills/do-work-board/justfile.template", []string{"staged-skills", "updater", "installer"}},
		{"skills/do-work/actions/version.md", []string{"staged-skills", "updater"}},
	}
	for _, coverageCase := range coverageCases {
		if got := manifestSelectedLaneIDs(manifest, coverageCase.path); !reflect.DeepEqual(got, coverageCase.wantLaneIDs) {
			t.Errorf("%s selected lanes = %v, want %v", coverageCase.path, got, coverageCase.wantLaneIDs)
		}
	}
}

func manifestSelectedLaneIDs(manifest laneManifest, changedPath string) []string {
	selectedLaneIDs := []string{}
	for _, lane := range manifest.Lanes {
		for _, rule := range lane.Coverage {
			if rule.matches(changedPath) {
				selectedLaneIDs = append(selectedLaneIDs, lane.ID)
				break
			}
		}
	}
	return selectedLaneIDs
}

func selectedHeavyLaneIDs(plan resultmodel.HeavyVerificationPlan) []string {
	identities := make([]string, 0, len(plan.SelectedLanes))
	for _, lane := range plan.SelectedLanes {
		identities = append(identities, lane.LaneID)
	}
	return identities
}

func newHeavyTestRepository(t *testing.T, manifestContents string) (string, string) {
	t.Helper()
	repositoryRoot := t.TempDir()
	runHeavyTestGit(t, repositoryRoot, "init", "-q")
	runHeavyTestGit(t, repositoryRoot, "config", "user.name", "Heavy Test")
	runHeavyTestGit(t, repositoryRoot, "config", "user.email", "heavy@example.test")
	writeHeavyTestFile(t, repositoryRoot, "seed.txt", "seed\n")
	writeHeavyTestFile(t, repositoryRoot, "heavy-lanes.json", manifestContents)
	return repositoryRoot, commitHeavyTestChanges(t, repositoryRoot, "seed")
}

func heavyTestManifestPath(repositoryRoot string) string {
	return filepath.Join(repositoryRoot, "heavy-lanes.json")
}

func writeHeavyTestFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	path := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func commitHeavyTestChanges(t *testing.T, repositoryRoot, message string) string {
	t.Helper()
	runHeavyTestGit(t, repositoryRoot, "add", "-A")
	runHeavyTestGit(t, repositoryRoot, "commit", "-qm", message)
	return runHeavyTestGit(t, repositoryRoot, "rev-parse", "HEAD")
}

func runHeavyTestGit(t *testing.T, repositoryRoot string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(arguments, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}
