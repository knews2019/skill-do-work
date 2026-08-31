package nextselection

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestMixedSuccessfulProbesCarryExactUnblockTargets(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-501-ordinary.md", "REQ-501", "pending", "")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-502-blocked.md", "REQ-502", "blocked", "blocked_by: service alpha\nblocked_check: printf alpha\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-503-ordinary.md", "REQ-503", "pending", "")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-504-blocked.md", "REQ-504", "blocked", "blocked_by: service beta\nblocked_check: printf beta\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	limit := 4
	result := Select(snapshot, dependencygraph.BuildGraph(snapshot), SelectionOptions{FanOutLimit: &limit}, func([]byte, int) (int, error) {
		return 0, nil
	})
	result.Command = "next"
	jsonBytes, err := resultmodel.RenderResult(result, resultmodel.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	type unblockRecord struct {
		RequestID       string `json:"request_id"`
		RequestPath     string `json:"request_path"`
		OriginalStatus  string `json:"original_status"`
		ProbeStatus     string `json:"probe_status"`
		UnblockRequired bool   `json:"unblock_required"`
	}
	var decoded struct {
		Selected []unblockRecord `json:"selected"`
	}
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Selected) != 4 {
		t.Fatalf("selected records = %#v", decoded.Selected)
	}
	for _, record := range decoded.Selected {
		wantPath := "do-work/queue/" + record.RequestID + map[string]string{
			"REQ-501": "-ordinary.md", "REQ-502": "-blocked.md", "REQ-503": "-ordinary.md", "REQ-504": "-blocked.md",
		}[record.RequestID]
		blocked := record.RequestID == "REQ-502" || record.RequestID == "REQ-504"
		if record.RequestPath != wantPath {
			t.Errorf("%s request_path = %q, want %q", record.RequestID, record.RequestPath, wantPath)
		}
		if blocked && (record.OriginalStatus != "blocked" || record.ProbeStatus != "succeeded" || !record.UnblockRequired) {
			t.Errorf("blocked record lacks exact unblock evidence: %#v", record)
		}
		if !blocked && (record.OriginalStatus != "pending" || record.ProbeStatus != "not_applicable" || record.UnblockRequired) {
			t.Errorf("ordinary record is falsely marked for unblock: %#v", record)
		}
	}
	textBytes, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"do-work/queue/REQ-502-blocked.md", "probe succeeded", "unblock required", "do-work/queue/REQ-504-blocked.md"} {
		if !strings.Contains(string(textBytes), required) {
			t.Errorf("text handoff missing %q:\n%s", required, textBytes)
		}
	}
	boundedLimit := 3
	bounded := Select(snapshot, dependencygraph.BuildGraph(snapshot), SelectionOptions{FanOutLimit: &boundedLimit}, func([]byte, int) (int, error) {
		return 0, nil
	})
	for _, exclusion := range bounded.Excluded {
		if exclusion.RequestID == "REQ-504" && exclusion.Code == "FAN-OUT-LIMIT" {
			if exclusion.RequestPath != "do-work/queue/REQ-504-blocked.md" || exclusion.ProbeStatus != resultmodel.ProbeSucceeded || !exclusion.UnblockRequired {
				t.Fatalf("fan-out exclusion lost successful-probe evidence: %#v", exclusion)
			}
			return
		}
	}
	t.Fatalf("bounded mixed selection omitted REQ-504 fan-out evidence: %#v", bounded.Excluded)
}

func TestExplicitREQOverridesDependencyAssignmentAndNegligibleFilters(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-201-target.md", "REQ-201", "pending", "user_request: UR-201\ndepends_on: [REQ-999]\nassigned_to: cloud-alpha\nimpact: impact-negligible\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-202-dependency.md", "REQ-202", "pending", "user_request: UR-202\ndepends_on: [REQ-999]\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-203-negligible.md", "REQ-203", "pending", "user_request: UR-203\nimpact: impact-negligible\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	graph := dependencygraph.BuildGraph(snapshot)

	defaultResult := Select(snapshot, graph, SelectionOptions{SkipImpactNegligible: true}, nil)
	if len(defaultResult.Selected) != 0 {
		t.Fatalf("default selection did not retain assignment gate: %#v", defaultResult)
	}
	assertExclusionCode(t, defaultResult, "REQ-201", "ASSIGNED-ELSEWHERE")
	explicitResult := Select(snapshot, graph, SelectionOptions{TargetTokens: []string{"REQ-201"}, SkipImpactNegligible: true}, nil)
	if got := selectedRequestIDsFromModel(explicitResult.Selected); !equalStrings(got, []string{"REQ-201"}) {
		t.Fatalf("explicit override selected %v, want REQ-201; exclusions=%#v", got, explicitResult.Excluded)
	}
	if explicitResult.Selected[0].Provenance != ProvenanceExplicit {
		t.Fatalf("explicit provenance lost: %#v", explicitResult.Selected[0])
	}
	userRequestResult := Select(snapshot, graph, SelectionOptions{TargetTokens: []string{"UR-201"}, SkipImpactNegligible: true}, nil)
	if len(userRequestResult.Selected) != 0 {
		t.Fatalf("UR expansion incorrectly inherited explicit overrides: %#v", userRequestResult.Selected)
	}
	assertExclusionCode(t, userRequestResult, "REQ-201", "ASSIGNED-ELSEWHERE")
	for _, test := range []struct {
		identifier  string
		userRequest string
		wantURCode  string
	}{
		{"REQ-202", "UR-202", "DEPENDENCY-MISSING"},
		{"REQ-203", "UR-203", "IMPACT-NEGLIGIBLE"},
	} {
		direct := Select(snapshot, graph, SelectionOptions{TargetTokens: []string{test.identifier}, SkipImpactNegligible: true}, nil)
		if len(direct.Selected) != 1 || direct.Selected[0].RequestID != test.identifier {
			t.Fatalf("direct %s did not preserve override: %#v", test.identifier, direct)
		}
		expanded := Select(snapshot, graph, SelectionOptions{TargetTokens: []string{test.userRequest}, SkipImpactNegligible: true}, nil)
		assertExclusionCode(t, expanded, test.identifier, test.wantURCode)
	}
}

func TestBlockedProbeReceivesExactBytesAndDoesNotMutateRequest(t *testing.T) {
	repositoryRoot := t.TempDir()
	probe := "printf 'service ready\\n' >/dev/null\n"
	relativePath := "do-work/queue/REQ-210-blocked.md"
	writeCommandRequest(t, repositoryRoot, relativePath, "REQ-210", "blocked", "blocked_by: local service\nblocked_check: |\n  printf 'service ready\\n' >/dev/null\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	graph := dependencygraph.BuildGraph(snapshot)
	original := append([]byte(nil), snapshot.RequestFiles[0].ContentBytes...)
	var received []byte
	result := Select(snapshot, graph, SelectionOptions{}, func(probeBytes []byte, timeoutSeconds int) (int, error) {
		received = append([]byte(nil), probeBytes...)
		if timeoutSeconds != 30 {
			t.Fatalf("timeout = %d, want 30", timeoutSeconds)
		}
		return 0, nil
	})
	if !bytes.Equal(received, []byte(probe)) {
		t.Fatalf("probe bytes = %q, want exact %q", received, probe)
	}
	if len(result.Selected) != 1 || result.SelectionSummary.Probed != 1 || result.SelectionSummary.ProbeSucceeded != 1 {
		t.Fatalf("successful probe did not select request: %#v", result)
	}
	after := snapshot.RequestFiles[0].ContentBytes
	if !bytes.Equal(original, after) || snapshot.RequestFiles[0].TypedRecord.RequestStatus != "blocked" {
		t.Fatal("read-only selection mutated blocked request evidence")
	}
	diskBytes, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	if err != nil || !bytes.Equal(original, diskBytes) {
		t.Fatalf("read-only selection changed the request on disk: %v", err)
	}
	timedOut := Select(snapshot, graph, SelectionOptions{}, func([]byte, int) (int, error) { return 124, nil })
	assertExclusionCode(t, timedOut, "REQ-210", "BLOCKED-PROBE-FAILED")
}

func TestBlockedProbeOutcomesRemainDistinctPerRecord(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-601-missing.md", "REQ-601", "blocked", "blocked_by: missing probe\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-602-failed.md", "REQ-602", "blocked", "blocked_check: fail-probe\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-603-timeout.md", "REQ-603", "blocked", "blocked_check: timeout-probe\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-604-success.md", "REQ-604", "blocked", "blocked_check: success-probe\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	limit := 4
	result := Select(snapshot, dependencygraph.BuildGraph(snapshot), SelectionOptions{FanOutLimit: &limit}, func(probeBytes []byte, _ int) (int, error) {
		switch strings.TrimSpace(string(probeBytes)) {
		case "fail-probe":
			return 23, nil
		case "timeout-probe":
			return 124, nil
		case "success-probe":
			return 0, nil
		default:
			t.Fatalf("unexpected probe bytes %q", probeBytes)
			return 125, nil
		}
	})
	result.Command = "next"
	jsonBytes, err := resultmodel.RenderResult(result, resultmodel.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	var rendered resultmodel.CommandResult
	if err := json.Unmarshal(jsonBytes, &rendered); err != nil {
		t.Fatal(err)
	}
	assertProbeExclusion(t, rendered, "REQ-601", resultmodel.ProbeMissing, false, -1)
	assertProbeExclusion(t, rendered, "REQ-602", resultmodel.ProbeFailed, true, 23)
	assertProbeExclusion(t, rendered, "REQ-603", resultmodel.ProbeTimedOut, true, 124)
	if len(rendered.Selected) != 1 {
		t.Fatalf("selected = %#v, want only successful probe", rendered.Selected)
	}
	success := rendered.Selected[0]
	if success.RequestID != "REQ-604" || success.RequestPath != "do-work/queue/REQ-604-success.md" || success.OriginalStatus != "blocked" || success.ProbeStatus != resultmodel.ProbeSucceeded || !success.ProbeAttempted || success.ProbeExitCode != 0 || !success.UnblockRequired {
		t.Fatalf("successful probe evidence = %#v", success)
	}
	textBytes, err := resultmodel.RenderResult(result, resultmodel.FormatText)
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"probe missing", "probe failed", "probe timed_out", "probe succeeded"} {
		if !strings.Contains(string(textBytes), status) {
			t.Errorf("text output omitted distinct %q outcome:\n%s", status, textBytes)
		}
	}
}

func assertProbeExclusion(t *testing.T, result resultmodel.CommandResult, identifier string, probeStatus resultmodel.SelectionProbeStatus, attempted bool, exitCode int) {
	t.Helper()
	for _, exclusion := range result.Excluded {
		if exclusion.RequestID != identifier {
			continue
		}
		if exclusion.RequestPath == "" || exclusion.OriginalStatus != "blocked" || exclusion.ProbeStatus != probeStatus || exclusion.ProbeAttempted != attempted || exclusion.ProbeExitCode != exitCode || exclusion.UnblockRequired {
			t.Fatalf("%s probe evidence = %#v", identifier, exclusion)
		}
		return
	}
	t.Fatalf("missing exclusion %s in %#v", identifier, result.Excluded)
}

func TestWaveDepthAndFanOutAreSeparateSelectionAxes(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-301-root.md", "REQ-301", "pending", "")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-302-root.md", "REQ-302", "pending", "")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-303-child.md", "REQ-303", "pending", "depends_on: [REQ-301]\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	graph := dependencygraph.BuildGraph(snapshot)
	wave := 0
	limit := 1
	result := Select(snapshot, graph, SelectionOptions{WaveDepth: &wave, FanOutLimit: &limit}, nil)
	if got := selectedRequestIDsFromModel(result.Selected); !equalStrings(got, []string{"REQ-301"}) {
		t.Fatalf("wave/fan-out selection = %v, want first root", got)
	}
	assertExclusionCode(t, result, "REQ-302", "FAN-OUT-LIMIT")
	assertExclusionCode(t, result, "REQ-303", "WAVE-MISMATCH")
}

func TestSimpleSelectionRetainsSpecializedVetoesAndFrozenEstimate(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-401-good.md", "REQ-401", "pending", "effort_estimate: trivial\nestimate:\n  p50_active_minutes: 15\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-402-maintenance.md", "REQ-402", "pending", "effort_estimate: effort-mechanical\nmaintenance: true\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-403-security.md", "REQ-403", "pending", "effort_estimate: effort-mechanical\ndomain: security\n")
	writeCommandRequest(t, repositoryRoot, "do-work/queue/REQ-404-critical.md", "REQ-404", "pending", "effort_estimate: effort-mechanical\nimpact: impact-critical\n")
	snapshot, _ := repositorymodel.DiscoverRepository(repositoryRoot)
	result := Select(snapshot, dependencygraph.BuildGraph(snapshot), SelectionOptions{SimpleOnly: true}, nil)
	if got := selectedRequestIDsFromModel(result.Selected); !equalStrings(got, []string{"REQ-401"}) {
		t.Fatalf("simple selected = %v; result=%#v", got, result)
	}
	if result.Selected[0].EstimateMinutes != 15 || !result.Selected[0].EstimateKnown {
		t.Fatalf("frozen estimate lost: %#v", result.Selected[0])
	}
	assertExclusionCode(t, result, "REQ-402", "MAINTENANCE-JUDGMENT")
	assertExclusionCode(t, result, "REQ-403", "SECURITY-RISK")
	assertExclusionCode(t, result, "REQ-404", "IMPACT-CRITICAL")
}

func selectedRequestIDsFromModel(records []resultmodel.SelectionRecord) []string {
	requestIDs := make([]string, len(records))
	for index := range records {
		requestIDs[index] = records[index].RequestID
	}
	return requestIDs
}

func assertExclusionCode(t *testing.T, result resultmodel.CommandResult, identifier, code string) {
	t.Helper()
	for _, exclusion := range result.Excluded {
		if exclusion.RequestID == identifier && exclusion.Code == code {
			return
		}
	}
	t.Fatalf("missing %s/%s exclusion in %#v", identifier, code, result.Excluded)
}
