package runstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var fixtureNow = time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)

const waitingReason = "Waiting on dependencies: it stays under Pending until every depends_on target is source-ready (terminally successful, or claimed with a commit)."

// fixtureRequest is one open REQ: the frontmatter written to disk and the row
// the board's `open-work --format json` would print for it.
type fixtureRequest struct {
	id, status, column, placementReason, extraFrontmatter string
	unmetDependencies                                     []string
	assignedTo                                            string
	lastActivityAt                                        time.Time
}

// writeRunStatusFixture is the one fixture builder: it writes each REQ under
// do-work/working (claimed) or do-work/queue, and the hand-written board facts
// file in the shape queue-kanban prints. It returns the repository root and the
// board facts path.
func writeRunStatusFixture(t *testing.T, requests ...fixtureRequest) (string, string) {
	t.Helper()
	repositoryRoot := t.TempDir()
	boardRequests := []map[string]any{}
	for _, request := range requests {
		folder := "queue"
		if request.status == "claimed" {
			folder = "working"
		}
		writeFixtureFile(t, repositoryRoot, fmt.Sprintf("do-work/%s/%s-fixture.md", folder, request.id),
			"---\nid: "+request.id+"\ntitle: Fixture request "+request.id+" with a long descriptive title\nstatus: "+request.status+"\n"+request.extraFrontmatter+"---\n# Fixture\n")
		boardRequest := map[string]any{
			"id": request.id, "title": "Fixture request " + request.id + " with a long descriptive title", "status": request.status,
			"column": request.column, "placement_reason": request.placementReason,
			"unmet_dependencies": append([]string{}, request.unmetDependencies...), "assigned_to": request.assignedTo,
		}
		if !request.lastActivityAt.IsZero() {
			boardRequest["last_activity_at"] = request.lastActivityAt.Format(time.RFC3339)
			boardRequest["last_activity_kind"] = "commit"
			boardRequest["last_activity_phase"] = "implementation"
		}
		boardRequests = append(boardRequests, boardRequest)
	}
	boardFacts, err := json.Marshal(map[string]any{"generated_at": fixtureNow.Format(time.RFC3339), "stale_claim_threshold_minutes": 180, "requests": boardRequests})
	if err != nil {
		t.Fatal(err)
	}
	boardFactsPath := filepath.Join(t.TempDir(), "board-facts.json")
	if err := os.WriteFile(boardFactsPath, boardFacts, 0o600); err != nil {
		t.Fatal(err)
	}
	return repositoryRoot, boardFactsPath
}

func writeFixtureFile(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func claimedFixture(id string, claimedAt time.Time, extra string) fixtureRequest {
	return fixtureRequest{id: id, status: "claimed", column: "claimed", placementReason: "Claimed: a run picked it up. The claim is a record, not a lock.",
		extraFrontmatter: "claimed_at: " + claimedAt.Format(time.RFC3339) + "\n" + extra, lastActivityAt: fixtureNow.Add(-5 * time.Minute)}
}

func runFixture(t *testing.T, repositoryRoot string, now time.Time, arguments ...string) resultmodel.CommandResult {
	t.Helper()
	result := runStatusAt(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot, Format: resultmodel.FormatText}, arguments, now)
	if result.Outcome != resultmodel.OutcomeSuccess || result.RunStatus == nil {
		t.Fatalf("outcome = %s, findings = %+v", result.Outcome, result.Findings)
	}
	return result
}

func findingFor(t *testing.T, result resultmodel.CommandResult, requestID string) resultmodel.CommandFinding {
	t.Helper()
	for _, finding := range result.Findings {
		if len(finding.AffectedIDs) == 1 && finding.AffectedIDs[0] == requestID {
			return finding
		}
	}
	t.Fatalf("no finding for %s in %+v", requestID, result.Findings)
	return resultmodel.CommandFinding{}
}

func textOf(result resultmodel.CommandResult) string {
	if result.ExactTextOutput == nil {
		return ""
	}
	return *result.ExactTextOutput
}

// The RED case of the REQ: a claimed REQ with a 60-minute estimate, claimed 45
// minutes ago, has about 15 minutes left; at 75 minutes it is over by 15 and
// never prints a negative ETA.
func TestRunStatusEtaCountsDownFromTheFrozenEstimate(t *testing.T) {
	claimedAt := fixtureNow.Add(-45 * time.Minute)
	repositoryRoot, boardFacts := writeRunStatusFixture(t, claimedFixture("REQ-801", claimedAt, "estimate:\n  p50_active_minutes: 60\n"))

	onTime := runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts)
	row := onTime.RunStatus.Rows[0]
	if row.EtaMinutes == nil || *row.EtaMinutes != 15 || row.EtaText != "15 min" || !strings.Contains(textOf(onTime), "15 min") {
		t.Fatalf("at 45 of 60 min: row = %+v\n%s", row, textOf(onTime))
	}

	overdue := runFixture(t, repositoryRoot, claimedAt.Add(75*time.Minute), "--board-facts", boardFacts)
	row = overdue.RunStatus.Rows[0]
	if row.EtaMinutes != nil || row.EtaText != "over estimate by 15 min" || !strings.Contains(textOf(overdue), "over estimate by 15 min") || strings.Contains(textOf(overdue), "-15") {
		t.Fatalf("at 75 of 60 min: row = %+v\n%s", row, textOf(overdue))
	}
}

// A hand-back that landed in the run directory means the build is done and
// waits for integration, whatever the claim's age says.
func TestRunStatusLandedHandbackIsClassC3(t *testing.T) {
	repositoryRoot, boardFacts := writeRunStatusFixture(t, claimedFixture("REQ-802", fixtureNow.Add(-30*time.Minute), ""))
	runDirectory := "do-work/runs/work-2026-10-10-113000"
	writeFixtureFile(t, repositoryRoot, runDirectory+"/manifest.md", "| REQ | Builder | Dispatch instant |\n| --- | --- | --- |\n| REQ-802 | builder-a | 2026-10-10T11:31:00Z |\n")
	writeFixtureFile(t, repositoryRoot, runDirectory+"/REQ-802-handback.md", "# Hand-back\n")

	result := runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts)
	finding := findingFor(t, result, "REQ-802")
	row := result.RunStatus.Rows[0]
	if finding.Code != "C3" || !reflect.DeepEqual(finding.NextArgv, []string{"do-work", "run"}) {
		t.Fatalf("finding = %+v", finding)
	}
	if row.HandbackPresent == nil || !*row.HandbackPresent || row.DispatchedAt != "2026-10-10T11:31:00Z" || row.ManifestRow != "| REQ-802 | builder-a | 2026-10-10T11:31:00Z |" {
		t.Fatalf("row = %+v", row)
	}
}

// A blocked REQ waiting on an unmet dependency is the board's Waiting group:
// the row quotes the board's own reason and names the REQ it waits on, and
// offers no command, because only the upstream REQ finishing moves it.
func TestRunStatusUnmetDependencyIsClassC5QuotingTheWaitingReason(t *testing.T) {
	repositoryRoot, boardFacts := writeRunStatusFixture(t, fixtureRequest{id: "REQ-803", status: "blocked", column: "pending-waiting",
		placementReason: waitingReason, unmetDependencies: []string{"REQ-999"}, extraFrontmatter: "depends_on: [REQ-999]\nblocked_by: [vendor key]\n"})

	finding := findingFor(t, runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts), "REQ-803")
	evidence := strings.Join(finding.Evidence, "\n")
	if finding.Code != "C5" || len(finding.NextArgv) != 0 || !strings.Contains(evidence, waitingReason) || !strings.Contains(evidence, "REQ-999") {
		t.Fatalf("finding = %+v", finding)
	}
}

// An earmarked REQ is skipped by another session's default run; the row
// quotes why and offers the by-name run that overrides it.
func TestRunStatusEarmarkedRequestIsClassC6QuotingTheEarmarkedReason(t *testing.T) {
	earmarkedReason := "Earmarked for cloud-alpha: an advisory claim marker, not a lock. Another session's default run skips it; naming it explicitly overrides that."
	repositoryRoot, boardFacts := writeRunStatusFixture(t, fixtureRequest{id: "REQ-804", status: "pending", column: "pending-earmarked",
		placementReason: earmarkedReason, assignedTo: "cloud-alpha", extraFrontmatter: "assigned_to: cloud-alpha\n"})

	finding := findingFor(t, runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts), "REQ-804")
	if finding.Code != "C6" || !reflect.DeepEqual(finding.NextArgv, []string{"do-work", "run", "REQ-804"}) || !strings.Contains(strings.Join(finding.Evidence, "\n"), earmarkedReason) {
		t.Fatalf("finding = %+v", finding)
	}
}

// destructive-next-argv: next_argv is followed literally, so a claim past the
// threshold offers the read-only history inspection, and the takeover command
// appears only in the stop reason, with what it destroys.
func TestRunStatusStaleClaimNamesRecoveryOnlyInTheStopReason(t *testing.T) {
	repositoryRoot, boardFacts := writeRunStatusFixture(t, claimedFixture("REQ-805", fixtureNow.Add(-200*time.Minute), ""))

	finding := findingFor(t, runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts), "REQ-805")
	wantArgv := []string{"git", "log", "--full-history", "--", "do-work/working/REQ-805-fixture.md"}
	if finding.Code != "C7" || !reflect.DeepEqual(finding.NextArgv, wantArgv) {
		t.Fatalf("finding = %+v", finding)
	}
	if !strings.Contains(finding.AutomationStopReason, "do-work run-with-recovery REQ-805") || !strings.Contains(finding.AutomationStopReason, "only if you know the run that claimed it is gone") {
		t.Fatalf("stop reason = %q", finding.AutomationStopReason)
	}
	if strings.Contains(strings.Join(finding.NextArgv, " ")+strings.Join(finding.Evidence, " "), "run-with-recovery") {
		t.Fatalf("run-with-recovery leaked outside the stop reason: %+v", finding)
	}
}

// A lock the suite does not define is listed with its age and first line, and
// run-status leaves it where it is: it may belong to another checkout.
func TestRunStatusListsRunLocalLockFileAndLeavesIt(t *testing.T) {
	repositoryRoot, boardFacts := writeRunStatusFixture(t, claimedFixture("REQ-806", fixtureNow.Add(-30*time.Minute), ""))
	runDirectory := "do-work/runs/work-2026-10-10-113000"
	writeFixtureFile(t, repositoryRoot, runDirectory+"/manifest.md", "| REQ |\n")
	writeFixtureFile(t, repositoryRoot, runDirectory+"/REQ-806-brief.md", "# Brief\n")
	writeFixtureFile(t, repositoryRoot, runDirectory+"/full-gate.lock", "pid 4242 bash _dev/tests/maintainer-verify.sh\nstarted 2026-10-10T11:48:00Z\n")
	lockPath := filepath.Join(repositoryRoot, filepath.FromSlash(runDirectory), "full-gate.lock")
	if err := os.Chtimes(lockPath, fixtureNow.Add(-12*time.Minute), fixtureNow.Add(-12*time.Minute)); err != nil {
		t.Fatal(err)
	}

	result := runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts)
	wantFiles := []resultmodel.RunLocalFile{{Path: runDirectory + "/full-gate.lock", AgeMinutes: 12, FirstLine: "pid 4242 bash _dev/tests/maintainer-verify.sh", Label: "run-local, not interpreted"}}
	if !reflect.DeepEqual(result.RunStatus.RunLocalFiles, wantFiles) {
		t.Fatalf("run-local files = %+v", result.RunStatus.RunLocalFiles)
	}
	if !strings.Contains(textOf(result), "full-gate.lock") {
		t.Fatalf("text does not list the lock:\n%s", textOf(result))
	}
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatalf("the lock is gone after run-status: %v", err)
	}
}

// forensics and status share one result model: through the real command
// runtime, each JSON finding carries every key doctor's findings carry, and
// its code is the row's class.
func TestRunStatusJSONRowsUseTheCommandFindingShape(t *testing.T) {
	repositoryRoot, boardFacts := writeRunStatusFixture(t, claimedFixture("REQ-807", time.Now().UTC().Add(-10*time.Minute), ""))
	var output bytes.Buffer
	exitCode := commandruntime.NewRuntime(&output, Handlers()).Run([]string{"--repo-root", repositoryRoot, "--format", "json", "run-status", "--board-facts", boardFacts})
	var decoded struct {
		Outcome   string           `json:"outcome"`
		Findings  []map[string]any `json:"findings"`
		RunStatus struct {
			Rows []struct {
				Class string `json:"class"`
			} `json:"rows"`
		} `json:"run_status"`
	}
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil || exitCode != 0 || decoded.Outcome != "success" {
		t.Fatalf("exit=%d err=%v\n%s", exitCode, err, output.String())
	}
	if len(decoded.Findings) != 1 || len(decoded.RunStatus.Rows) != 1 {
		t.Fatalf("want one finding and one row:\n%s", output.String())
	}
	for _, key := range []string{"code", "severity", "affected_ids", "affected_paths", "observed_evidence", "fixability", "automation_stop_reason", "next_argv", "next_just_recipe", "verification_argv"} {
		if _, present := decoded.Findings[0][key]; !present {
			t.Fatalf("finding lacks %q:\n%s", key, output.String())
		}
	}
	if decoded.Findings[0]["code"] != decoded.RunStatus.Rows[0].Class || decoded.Findings[0]["severity"] != "info" {
		t.Fatalf("code %v is not the row class %q", decoded.Findings[0]["code"], decoded.RunStatus.Rows[0].Class)
	}
}

// --watch is a /loop body: eight open REQs fit in 20 lines, every REQ shows.
func TestRunStatusWatchFitsTwentyLinesForEightOpenRequests(t *testing.T) {
	requests := []fixtureRequest{}
	for index := 0; index < 8; index++ {
		requests = append(requests, claimedFixture(fmt.Sprintf("REQ-81%d", index), fixtureNow.Add(-time.Duration(20+index*30)*time.Minute), "estimate:\n  p50_active_minutes: 90\n"))
	}
	repositoryRoot, boardFacts := writeRunStatusFixture(t, requests...)
	writeFixtureFile(t, repositoryRoot, "do-work/runs/work-2026-10-10-113000/full-gate.lock", "pid 4242\n")

	watchText := textOf(runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts, "--watch"))
	if lineCount := strings.Count(strings.TrimRight(watchText, "\n"), "\n") + 1; lineCount > 20 {
		t.Fatalf("--watch printed %d lines:\n%s", lineCount, watchText)
	}
	for _, request := range requests {
		if !strings.Contains(watchText, request.id) {
			t.Fatalf("--watch omits %s:\n%s", request.id, watchText)
		}
	}
}

// No run directory is not an error: the queue-side rows still print and the
// run fields are absent, not false.
func TestRunStatusReportsQueueRowsWithoutARunDirectory(t *testing.T) {
	repositoryRoot, boardFacts := writeRunStatusFixture(t, claimedFixture("REQ-808", fixtureNow.Add(-10*time.Minute), ""))

	result := runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts)
	if result.RunStatus.RunDirectory != "" || len(result.RunStatus.Rows) != 1 || result.RunStatus.Rows[0].HandbackPresent != nil || findingFor(t, result, "REQ-808").Code != "C1" {
		t.Fatalf("run status = %+v", result.RunStatus)
	}
}
