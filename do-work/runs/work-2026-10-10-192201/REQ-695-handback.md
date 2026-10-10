# REQ-695 hand-back (run-status names a missing --run folder and stops recommending do-work run while an integrator works)

- Branch: `worktree-agent-REQ-695-run-status-missing-run-integrator`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-695-run-status-missing-run-integrator`
- Base commit: `fd9a2378`
- Commits: `9fd11e18` `[REQ-695] name a missing run folder in run-status and stop recommending do-work run while an integrator works` (one commit, tests and fix together)

## File manifest

- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go` (modified): `integratorPhases` set; stat of the resolved run folder in `runStatusAt` (missing folder → empty folder to the three readers, display path + ` (not found; run fields absent)`); `report.RunDirectory` set from the display value; new C3 arm before the old one in `classifyRow` (empty `next_argv`, integrator remedy).
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go` (modified): full text form prints `run: <RunDirectory>` just before `next:` when set. `--watch` untouched.
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go` (modified): `lastActivityPhase` fixture field (empty → `implementation`), two new tests.
- `skills/do-work/actions/status.md` (modified): C3 remedy cell names the integrator exception (brief text verbatim).
- `skills/do-work/docs/status-guide.md` (modified): C3 "What to do" cell names the integrator exception (brief text verbatim).

## P-A-U

- **[PLAN]:** Follow the brief exactly. TDD: add the fixture field and the two tests first, run RED, then F3 (stat once after `resolveRunDirectory`, empty folder to readers, suffix on display path, `run:` text line in full form only) and F4 (one C3 arm keyed on the board's integrator phases and the existing quiet boundary, empty `next_argv`), then the two prose cells.
- **[APPLY]:** Done as planned, five files, no other path touched. The F4 arm condition spans two lines; no helper added (brief allows one only if the condition is unreadable; it is readable).
- **[UNIFY]:** `git diff fd9a2378 --stat`:
  ```
   skills/do-work/actions/status.md                   |  2 +-
   skills/do-work/docs/status-guide.md                |  2 +-
   .../do-work-cli/internal/runstatus/run_status.go   | 25 ++++++++++---
   .../internal/runstatus/run_status_render.go        |  3 ++
   .../internal/runstatus/run_status_test.go          | 42 ++++++++++++++++++++++
   5 files changed, 68 insertions(+), 6 deletions(-)
  ```
  Checks (from worktree root, after the change):
  - Probe `bash .../REQ-695-probe.sh`: exit 0, wall 1 s, printed `REQ-695 GREEN: missing run folder named, integrating C3 recommends no run, docs updated`.
  - `gofmt -l skills/do-work/tools/do-work-cli/internal/runstatus`: empty, exit 0.
  - `go vet -C skills/do-work/tools/do-work-cli ./internal/runstatus/`: exit 0.
  - `git diff --check`: exit 0.
  - Focused GREEN run: exit 0, wall under 1 s. Full `^TestRunStatus`: exit 0, wall 1 s.
  - Files reviewed: all five above, diff read line by line; no debug artifacts. Grep for other Go/shell callers of run-status text or `run_directory` outside `internal/runstatus/`: only `main.go` (dispatch), `resultmodel` (type), `nextselection` and the board (unrelated `run-status` mentions); no other test asserts on this text.

## Proof record

RED (tests added, production unchanged), `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestRunStatusMissingRunFolderIsNamedInTextAndJSON|TestRunStatusIntegratingHandbackDoesNotRecommendRun)$' ./internal/runstatus/`, exit 1:
```
run_status_test.go:309: run_directory = "do-work/runs/typo"
--- FAIL: TestRunStatusMissingRunFolderIsNamedInTextAndJSON (0.02s)
run_status_test.go:330: finding = {Code:C3 ... AutomationStopReason:the build is done and waits for integration; run `do-work run` NextArgv:[do-work run] ...}
--- FAIL: TestRunStatusIntegratingHandbackDoesNotRecommendRun (0.02s)
```
(The first test stops at its first assertion: no "not found" in `run_directory`. At base the same row also had `HandbackPresent` false, which the later assertion checks.)

GREEN, same command, exit 0:
```
--- PASS: TestRunStatusMissingRunFolderIsNamedInTextAndJSON (0.02s)
--- PASS: TestRunStatusIntegratingHandbackDoesNotRecommendRun (0.02s)
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/runstatus	0.351s
```
Full `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^TestRunStatus' ./internal/runstatus/`, exit 0: all 12 PASS (EtaCountsDown, LandedHandbackIsClassC3, UnmetDependencyIsClassC5, EarmarkedRequestIsClassC6, StaleClaimNamesRecoveryOnlyInTheStopReason, ListsRunLocalLockFileAndLeavesIt, JSONRowsUseTheCommandFindingShape, WatchFitsTwentyLines, ReportsQueueRowsWithoutARunDirectory, FailedRequestNamesForensicsNotClarify, and the two new ones). `ok ... 0.464s`.

## Decisions

- **D-06 DECIDE & STATE: the stat applies to any resolved folder, not only to an explicit `--run`.** The auto-picked newest `work-*` folder always exists (it comes from `ReadDir`), so the stat costs nothing there and needs no branch on where the path came from. Value: one code path. Risk: none.
- **D-07 DECIDE & STATE: a path that exists but is a file counts as "not found".** The brief says "stat fails or it is not a directory"; the note text stays the reviewer's. Value: run fields stay absent instead of reading a file as a folder. Risk: the word "not found" is slightly loose for a file; reversible wording.
- **D-08 DECIDE & STATE: test REQ numbers are REQ-820 and REQ-821.** REQ-801 to REQ-817 are taken in the file.
- **D-09 DECIDE & STATE: the F4 test reads the last text line as the `next:` line** rather than grepping the whole text, because the C3 remedy line itself must not be confused with the recommendation. It asserts the line starts with `next:` so a reordered render fails loudly.

## Discovered Tasks

- `run_status_render.go` `--watch` form never shows the run folder, so a mistyped `--run` under `/loop` is still silent there (kept out by the 20-line budget, D-03). A one-word marker on the header line would fit the budget. → report only
- The `integratorPhases` set duplicates the board's phase labels with no lock-step test; if `phaseMilestonesOf` renames a label, the arm silently stops firing and C3 falls back to `do-work run` (the pre-REQ behavior). A board-side fact (for example an `integrating` flag in `open-work` JSON) would remove the copy. → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file; families `canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`; none applied, no links or manifests touched).
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` traps `closed-enumeration-for-a-condition` (phase set is the board's fixed pipeline, cited at source) and `silent-skip-reads-as-red` (F3 itself).
- `lessons-do-work-cli.md` rules as restated in the brief: `destructive-next-argv` (the new arm's `next_argv` is empty), `silent-skip-reads-as-red`. Satellite not opened.

## Anti-bloat check

Diff stat above. Added beyond what the brief named: none. Named additions only: `integratorPhases` (the phase set, brief item 2), `lastActivityPhase` fixture field (brief item 3), the two named tests. Local `runDirectoryDisplay` in `runStatusAt` is a short-lived local carrying the display path the brief requires; it replaced the old `if runDirectory != "" { report.RunDirectory = ... }` block.

## Proposed CHANGELOG entry

**run-status names a missing run folder and stops recommending `do-work run` while an integrator works**

A mistyped `--run` used to read as "no run", quietly turning hand-back-landed rows into stale-claim rows without saying which folder was read. And a REQ under active integration was still told to run `do-work run`, inviting a second orchestrator.

- `run-status` text now prints `run: <folder>`; when the folder does not exist, JSON `run_directory` and that line add `(not found; run fields absent)`, and hand-back and manifest fields are absent rather than false. Still no refusal and no exit-code change.
- A C3 (hand-back landed) row whose last activity is in phase integration, review, remediation or re-review within the last 20 minutes has no `next_argv` and says an integrator is working on it now.
- `actions/status.md` and `docs/status-guide.md` C3 rows name the integrator exception.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

- [family: silent-skip-reads-as-red] [REQ-695: a report that resolves a path the user typed must say which path it read and whether it exists; a missing folder should read as "absent" plus a visible note, never as a false verdict on what the folder would hold](../../../../do-work/archive/UR-154/REQ-695-run-status-missing-run-and-integrating-c3.md)

## Integration seams and wall times

- No overlap with REQ-693, REQ-694, REQ-696 write sets. Shared seam only: the CHANGELOG entry and lesson bullet above, which the integrator writes.
- Wall times (four builders running): focused tests under 1 s, full `^TestRunStatus` 1 s, probe 1 s. No budget failures, no reruns.
