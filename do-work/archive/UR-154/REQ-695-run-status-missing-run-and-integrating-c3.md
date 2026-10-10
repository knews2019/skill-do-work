---
id: REQ-695
title: 'Addendum: run-status names a missing --run folder and stops recommending do-work run while an integrator works'
status: completed
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
addendum_to: REQ-690
domain: backend
prime_files: ["_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-693, REQ-694, REQ-696]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go", "skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go", "skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go", "skills/do-work/actions/status.md", "skills/do-work/docs/status-guide.md"]
claimed_at: 2026-10-10T19:22:01Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-10-10T19:25:00Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
dispatch_at: 2026-10-10T19:27:42Z
builder_handback_at: 2026-10-10T19:28:50Z
integration_at: 2026-10-10T20:07:17Z
review_at: 2026-10-10T20:56:03Z
kb_status: pending
commit: bf9ed3d6b25c480346aeceac0511c5420fb07024
heavy_verified_at: 2026-10-10T20:56:14Z
heavy_verified_revision: bf9ed3d6b25c480346aeceac0511c5420fb07024
completed_at: 2026-10-10T20:56:44Z
release_at: 2026-10-10T20:56:44Z
---
# Addendum: run-status Names a Missing --run Folder and Stops Recommending do-work run While an Integrator Works

## What

Two gaps the REQ-690 (do-work status action and run-status report) review left as report only. A mistyped `--run` silently reads as "no run", so hand-backs count as absent and C3 (hand-back landed) rows become C7 (claim past threshold) rows whose text names `run-with-recovery`, and the text never says which run folder it read (review F3). C3 also recommends `do-work run` while an integrator is working on that REQ right now (review F4), which invites a second orchestrator.

## Prior Implementation

REQ-690 shipped `do-work-cli run-status` in `skills/do-work/tools/do-work-cli/internal/runstatus/`, archived at commit `65820f45`. Decision D-19 (coordinator ruling) removed the `--run` stat and refusal, and D-20 declined the reviewer's optional display line. `classifyRow` sends every claimed REQ with a landed hand-back to C3 with `next_argv ["do-work","run"]`; the row already carries `last_activity_phase` and minutes since last activity.

## Detailed Requirements

- F3, display only (D-19 stays: no refusal, no exit-code change): when the `--run` folder does not exist, the report says so. JSON `run_directory` gains a "not found; run fields absent" note (or an equivalent field), and the text output prints a `run: <dir>` line. The reviewer's exact suggestion is in `REQ-690-review.md` under F3.
- F4: a C3 row whose last activity is in phase `integration` and recent (inside the existing quiet boundary) does not recommend `do-work run`. It says an integrator is working on it and carries no runnable `next_argv`. Other C3 rows are unchanged.

## Red-Green Proof

**RED prompt/case:** (1) `run-status --run do-work/runs/typo` on a fixture with a claimed REQ, a landed hand-back in the real run folder, and a claim older than the threshold. (2) A fixture C3 row with last activity in phase `integration` 5 minutes ago.
**Why RED now:** (1) the row is C7 and the output never names the folder it read. (2) the row says "run `do-work run`".
**GREEN when:** (1) the text and JSON say the run folder was not found. (2) the row says an integrator is working on it and its `next_argv` is not `do-work run`.
**Validation:** Inferred during capture

## Constraints

- Keep the D-19 ruling: a missing `--run` folder is reported, never refused.
- One test per failure in `run_status_test.go`.
- `next_argv` stays read-only or empty; never a destructive command.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high (F3 reproduced on the live repo, F4 seen live on REQ-690 itself; both re-checked at HEAD `0e8e0ef9`). Latitude: the builder picks whether F4 is a new class or a C3 variant; prefer the smaller change.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 20802 tokens, bare only (`slugged: partial`); matches the runstatus path (families `destructive-next-argv`, `silent-skip-reads-as-red`). Its `destructive-next-argv`, `silent-skip-reads-as-red`, `closed-enumeration-for-a-condition` and `cross-module-fact-handoff` bullets were read at pre-dispatch and their rules are restated in the builder brief.
- `_dev/primes/lessons-action-files.md` — 8128 tokens, bare only (`slugged: partial`); matches the status-contract cell edits in `skills/do-work/actions/status.md` (D-04).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/archive/UR-153/REQ-690-run-status-action.md` → `## Review` and D-19, D-20; `do-work/runs/work-2026-10-10-131527/REQ-690-review.md` (F3, F4); `do-work/runs/work-2026-10-10-131527/REQ-690-integration-report.md`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Follow the brief exactly. TDD: add the fixture field and the two tests first, run RED, then F3 (stat once after `resolveRunDirectory`, empty folder to readers, suffix on display path, `run:` text line in full form only) and F4 (one C3 arm keyed on the board's integrator phases and the existing quiet boundary, empty `next_argv`), then the two prose cells.
- [x] **[APPLY]:** Done as planned, five files, no other path touched. The F4 arm condition spans two lines; no helper added (brief allows one only if the condition is unreadable; it is readable).
- [x] **[UNIFY]:** `git diff fd9a2378 --stat`: 5 files, 68 insertions, 6 deletions. Probe `REQ-695-probe.sh` exit 0; `gofmt -l` on `internal/runstatus` empty; `go vet ./internal/runstatus/` exit 0; `git diff --check` exit 0; focused GREEN run and full `^TestRunStatus` exit 0. All five files reviewed line by line, no debug artifacts. No other Go or shell caller asserts on run-status text or `run_directory` outside `internal/runstatus/`. (from the builder hand-back)
*Source: UR-154 R4 — "a mistyped `--run` silently turns C3 (hand-back landed) rows into C7 (claim past threshold) rows (review F3); C3 recommends `do-work run` even while an integrator runs (F4)."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and the files are named, but the change also touches the shipped C3 remedy wording in `skills/do-work/actions/status.md` and `skills/do-work/docs/status-guide.md`, and the fixture builder hard-codes the activity phase, so the exact seams needed exploration. Route B also gives the coordinator a pre-flight baseline for the focused runstatus tests.

**Planning:** Not required

**Pre-dispatch decisions (unattended, no `## Open Questions` section exists, so none came from capture):**
- **D-01 DECIDE & STATE: F4 is a C3 variant, not a new class.** The REQ gives the builder this latitude and asks for the smaller change. The row keeps code `C3`, gets an empty `next_argv`, and its remedy says an integrator is working on it. Value: no change to `classPrecedence`, `classNames`, or the class tables. Risk: the class name "hand-back landed, not integrated" is slightly early for a row that is mid-integration; the remedy text carries the truth. Reversible.
- **D-02 DECIDE & STATE: F4's guard covers the phase the integrator opens and every later integrator-owned phase (integration, review, remediation, re-review), not only `integration`.** `integration_at` is stamped after the hand-back merges, and `review_at`, `remediation_at`, `re_review_at` follow it (`actions/work.md:310`, `:376`). A claimed REQ with a landed hand-back whose recent activity sits in any of these is being worked by the integrator, so `do-work run` invites the same second orchestrator the review's F4 names. The phase labels are the board's fixed pipeline (`skills/do-work-board/tools/queue-kanban/durations.go:264` `phaseMilestonesOf`, lowercased). Value: closes the whole failure, not one phase of it. Risk: one string set to keep in step with the board's milestone list; the test pins `integration`, the REQ's RED case. Reversible.
- **D-03 DECIDE & STATE: F3 uses the reviewer's text (suffix on `run_directory`, a `run: <dir>` text line), not a new JSON field.** The suffix stays inside the declared Go files; a new field would need `internal/resultmodel/result_model.go`. The `run:` line prints only in the full text form, not in `--watch` (the 20-line `--watch` budget is pinned by a test). D-19 stays: no refusal, no exit-code change, the row class is not changed. Value: smallest change that names the folder read and says it is missing. Risk: a JSON reader that treats `run_directory` as a bare path gets a path plus a note when the folder is missing; no reader exists in the repository (grep for `run_directory` finds only the result model). Reversible.
- **D-04 DECIDE & STATE: the two shipped C3 table cells are updated.** `skills/do-work/actions/status.md:48` and `skills/do-work/docs/status-guide.md:26` say C3's remedy is `do-work run`; after D-01 that is false while an integrator works. One cell each. No sibling REQ in this batch touches either file (REQ-696's write set is forensics.md, README.md, clarify.md, capture.md and the contract test).

<!-- D-XX counter: last used D-05 (D-05 is recorded under ## Exploration). Next decision: D-06. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

All paths below are under `skills/do-work/tools/do-work-cli/internal/runstatus/` unless named in full. Line numbers are at main HEAD `fd9a2378`.

**F3 (missing `--run` folder):**
- `run_status.go:115` calls `resolveRunDirectory` (`:210-228`), which joins a relative `--run` onto the repository root and returns it with no stat (D-19). An empty string means no run folder at all.
- The resolved folder feeds three readers: `readManifestLines` (`:132`, `:230-239`), `buildRow` (`:140`; the hand-back stat and manifest row at `:290-302`), and `listRunLocalFiles` (`:149`, `:365-387`). Each one guards only on `runDirectory != ""`, so a missing folder gives `handback_present: false` (not absent), no manifest row, and no run-local files.
- `run_status.go:150-152` sets `report.RunDirectory` (JSON `run_directory`, `internal/resultmodel/result_model.go:645`, `omitempty`). Nothing else in the repository reads `run_directory`.
- `run_status_render.go:15-69` `renderReport`: the full text form never prints the run folder. The final `output.WriteString(nextLine + "\n")` is at `:67`. The `--watch` form returns early at `:41-42`.
- Existing pin: `TestRunStatusReportsQueueRowsWithoutARunDirectory` (`run_status_test.go:275-284`) says "the run fields are absent, not false" when there is no run folder. D-05 below keeps that contract for a missing `--run` folder too.

**F4 (C3 while an integrator works):**
- `classifyRow` (`run_status.go:317-361`), C3 arm at `:325-327`: any claimed row with `HandbackPresent` true gets `next_argv ["do-work","run"]` and the remedy "the build is done and waits for integration; run `do-work run`". It reads neither `LastActivityPhase` nor `MinutesSinceActivity`, both already on the record (`buildRow` `:245`, `:281-283`).
- `quietActivityBoundary` (`run_status.go:39`, 20 min) is the existing "recent" boundary used by C2 (`:348`).
- Phase labels come from the board: `skills/do-work-board/tools/queue-kanban/activity_correlation.go:181` copies the newest event's phase, and the labels are `phaseMilestonesOf` (`skills/do-work-board/tools/queue-kanban/durations.go:264-277`) lowercased: claimed, planning, dispatch, builder handback, integration, review, remediation, re-review, completed, release. `integration_at` is stamped after the hand-back merges (`skills/do-work/actions/work.md:310`); `review_at`, `remediation_at`, `re_review_at` follow (`skills/do-work/actions/work.md:376`).
- `renderReport` picks the `next:` line from the first row with a non-empty `next_argv` (`run_status_render.go:17-23`), so a C3 row with an empty `next_argv` drops out of the `next:` line with no render change.

**Tests and fixtures:**
- `writeRunStatusFixture` (`run_status_test.go:35-67`) hard-codes `last_activity_phase` to `implementation` (`:54`). An F4 test needs a phase field on `fixtureRequest` (`:24-29`) that defaults to the current value when empty, so no other test changes.
- `claimedFixture` (`:80-83`) puts last activity 5 minutes before `fixtureNow`; that fits the F4 RED case as is.
- `TestRunStatusLandedHandbackIsClassC3` (`:134-149`) is the "other C3 rows are unchanged" pin and must stay green untouched.
- `TestRunStatusStaleClaimNamesRecoveryOnlyInTheStopReason` (`:181-195`) shows the C7 shape a mistyped `--run` produces today (claim 200 min old, threshold 180).
- `TestRunStatusWatchFitsTwentyLinesForEightOpenRequests` (`:256-273`) pins the 20-line `--watch` budget, so the `run:` line stays out of `--watch` (D-03).

**Shipped prose that states the C3 remedy:** `skills/do-work/actions/status.md:48` and `skills/do-work/docs/status-guide.md:26`. No contract test greps either cell.

**Lessons that apply (read at pre-dispatch from the dropped `lessons-do-work-cli.md`):** `destructive-next-argv` (`next_argv` is followed literally, so the F4 row carries an empty `next_argv`, never a command), `silent-skip-reads-as-red` (a silent fallback reads as a verdict, which is exactly F3), `closed-enumeration-for-a-condition` (D-02's phase set is the board's fixed pipeline order, cited at its source, not a sample of examples).

- **D-05 DECIDE & STATE: a missing `--run` folder also leaves the run fields absent, not false.** The note the REQ asks for says "run fields absent", and the existing pin says run fields are absent, not false, when no run folder exists. Today a missing folder still stats the hand-back and reports `handback_present: false`. The builder passes an empty folder to the three readers when the stat fails, and keeps the typed path for display. Value: the JSON note is true, and a mistyped folder looks the same as no folder plus a visible note. Risk: none beyond D-19; no class changes, because a false and an absent hand-back both skip C3.

*Explored directly by the pre-dispatch agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go` (modify) — F3 stat and "not found" note, run fields absent for a missing folder (D-05); F4 C3 variant with no next argv while an integrator works (D-01, D-02)
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go` (modify) — print the run folder line in the full text form (D-03)
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go` (modify) — one phase field on the fixture, one test per failure (F3, F4)
- `skills/do-work/actions/status.md` (modify) — C3 table cell names the integrator exception (D-04)
- `skills/do-work/docs/status-guide.md` (modify) — C3 table cell names the integrator exception (D-04)

**Files I will NOT touch:** `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (no new JSON field, D-03), the queue-kanban board module, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, any version mirror, anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] When the run-status run folder option names a folder that does not exist, JSON `run_directory` carries a "not found; run fields absent" note and the text output prints a `run: <dir>` line; no refusal and no exit-code change (D-19 stays).
- [ ] The text output names the run folder it read (`run: <dir>`) whenever there is one.
- [ ] A C3 row whose last activity is in phase `integration` and inside the 20-minute quiet boundary does not recommend `do-work run`: it says an integrator is working on it and carries no runnable `next_argv`.
- [ ] Other C3 rows are unchanged (`TestRunStatusLandedHandbackIsClassC3` stays green untouched).
- [ ] One test per failure in `run_status_test.go`; `next_argv` stays read-only or empty.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go` (modified)
- `skills/do-work/actions/status.md` (modified)
- `skills/do-work/docs/status-guide.md` (modified)
- `skills/do-work/actions/fan-out-reference.md` (modified, review F3)

**What was done:** F3: `runStatusAt` stats the resolved run folder once. When it is missing or not a directory, the three readers get an empty folder (run fields absent, not false) and `run_directory` becomes the display path plus ` (not found; run fields absent)`. No refusal, no exit-code change (D-19 stays). The full text form prints `run: <folder>` just before `next:`; `--watch` is unchanged. F4: a new C3 arm before the old one in `classifyRow` matches a claimed REQ with a landed hand-back whose last activity is in an integrator phase (`integratorPhases`: integration, review, remediation, re-review) and inside the 20-minute quiet boundary. The row stays C3, carries an empty `next_argv`, and says an integrator is working on it now. The C3 cells in `actions/status.md` and `docs/status-guide.md` name the exception. Builder commit `9fd11e18`; integrator commit `3ea1e4ae` on the builder branch after review (F1: the integrator test also checks that a quiet integration row, 25 min old, offers `do-work run` again; F3: the coordinator stall-loop sentence in `actions/fan-out-reference.md` excludes a C3 row an integrator is working on; F4: `actions/status.md` Step 4 relays the `run:` line; F7: the `status.md` C3 cell names the four phases). Tests: `TestRunStatusMissingRunFolderIsNamedInTextAndJSON`, `TestRunStatusIntegratingHandbackDoesNotRecommendRun`, plus a `lastActivityPhase` fixture field (empty means `implementation`).

## Decisions

(from the builder hand-back, verbatim)

- **D-06 DECIDE & STATE: the stat applies to any resolved folder, not only to an explicit `--run`.** The auto-picked newest `work-*` folder always exists (it comes from `ReadDir`), so the stat costs nothing there and needs no branch on where the path came from. Value: one code path. Risk: none.
- **D-07 DECIDE & STATE: a path that exists but is a file counts as "not found".** The brief says "stat fails or it is not a directory"; the note text stays the reviewer's. Value: run fields stay absent instead of reading a file as a folder. Risk: the word "not found" is slightly loose for a file; reversible wording.
- **D-08 DECIDE & STATE: test REQ numbers are REQ-820 and REQ-821.** REQ-801 to REQ-817 are taken in the file.
- **D-09 DECIDE & STATE: the F4 test reads the last text line as the `next:` line** rather than grepping the whole text, because the C3 remedy line itself must not be confused with the recommendation. It asserts the line starts with `next:` so a reordered render fails loudly.

- **D-10 DECIDE & STATE (integrator, after review):** fixed review F1, F3, F4 and F7 on the builder branch (`3ea1e4ae`, merged as `bf9ed3d6`). F3 edits `skills/do-work/actions/fan-out-reference.md`, which is outside the declared write set: its stall-loop sentence told a coordinator that every C3 row means "brief the next integrator", which would start the second integrator this REQ exists to prevent. One clause, no other change. Value: shipped prose agrees with the new C3 arm. Risk: none beyond a wider diff.
- **D-11 DECIDE & STATE (integrator, after review):** review F2 (a contract check that `integratorPhases` matches the board's `phaseMilestonesOf` labels) stays report only. A drifted label makes the arm stop firing, which falls back to the pre-REQ `do-work run` advice, not to a wrong command; a new cross-module shell check is outside the write set. The reviewer agreed in the delta re-review. Value: no new test scaffolding. Risk: the arm can go quiet after a board rename.

<!-- D-XX counter: last used D-11. Next decision: D-12. -->

## Discovered Tasks

(from the builder hand-back, verbatim)

- `run_status_render.go` `--watch` form never shows the run folder, so a mistyped `--run` under `/loop` is still silent there (kept out by the 20-line budget, D-03). A one-word marker on the header line would fit the budget. → report only
- The `integratorPhases` set duplicates the board's phase labels with no lock-step test; if `phaseMilestonesOf` renames a label, the arm silently stops firing and C3 falls back to `do-work run` (the pre-REQ behavior). A board-side fact (for example an `integrating` flag in `open-work` JSON) would remove the copy. → report only

## Qualification

**Gate records:** `advance --diff-range d2991804..5d4e136c` returned `qualify` and `scope-drift` both `satisfied` (provenance `merged_range`), no findings.

**Requirement trace** (against `git diff d2991804..5d4e136c`, 5 files, +68 -6, and the merged files):
- F3, missing folder: `runStatusAt` stats the folder after `resolveRunDirectory`; a stat error or a non-directory sets the reader folder to `""` (hand-back, manifest and run-local readers all skip, so `handback_present` is absent, not false) and `run_directory` to `<path> (not found; run fields absent)`. No new error return, so no refusal and no exit-code change (D-19).
- F3, text names the folder: `renderReport` writes `run: <RunDirectory>` before `next:` whenever `RunDirectory` is set; the `--watch` branch returns earlier and is untouched (20-line budget, D-03).
- F4: the new arm precedes the old C3 arm and needs a claimed row, a landed hand-back, a phase in `integratorPhases`, and `MinutesSinceActivity` under `quietActivityBoundary` (20 min). It sets class C3, leaves `next_argv` empty, and the remedy says an integrator is working on it. Because `renderReport` takes `next:` from the first row with a non-empty `next_argv`, the row drops out of the `next:` line.
- Other C3 rows: the old arm is unchanged; `TestRunStatusLandedHandbackIsClassC3` is untouched and passes.
- One test per failure: `TestRunStatusMissingRunFolderIsNamedInTextAndJSON` (F3) and `TestRunStatusIntegratingHandbackDoesNotRecommendRun` (F4). `next_argv` is empty on the new arm, never a command.
- Prose: the C3 cells in `skills/do-work/actions/status.md` and `skills/do-work/docs/status-guide.md` name the integrator exception (D-04).

**After review:** cumulative range `d2991804..bf9ed3d6`, 6 files. `advance` was then in the review phase and refused gate input, so the delta was checked by hand: no debug lines, `gofmt -l` empty, `git diff --check` clean.

**Scope:** Route B. `## Scope` declares five files; the first merge touched exactly those five. The review fix adds `skills/do-work/actions/fan-out-reference.md` (one clause, review F3, D-10). `result_model.go`, the board module, changelogs and `do-work/` are untouched. Focused `go test ./internal/runstatus/` on main at `5d4e136c`: exit 0.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at the merge `5d4e136c` (1-minute load 3.13 at start, wall 140 s, exit 0), recorded through `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-192201/REQ-695-probe.sh`: `green-gate`, `run-blocked-check` and `scope-drift` satisfied, GREEN probe exit 0. Focused `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/runstatus/` exit 0 on main.
Gate again at the review-fix merge `bf9ed3d6` (load 4.99 at start, wall 153 s, exit 0), recorded with `record-green-gate` (recorded revision `bf9ed3d6b25c480346aeceac0511c5420fb07024`) because `advance` was in the review phase; GREEN probe exit 0 at `bf9ed3d6`.
**Result:** ✓ All passing (do-work-cli 917 Go tests, queue-kanban 423, no test file over the 30 s budget).

**Red-green validation:** (traces to `## Red-Green Proof`, from the builder hand-back)
- `TestRunStatusMissingRunFolderIsNamedInTextAndJSON` (RED case 1): ✗ before (`run_directory = "do-work/runs/typo"`, no "not found") → ✓ after.
- `TestRunStatusIntegratingHandbackDoesNotRecommendRun` (RED case 2): ✗ before (C3 with `NextArgv:[do-work run]`) → ✓ after.
- Quiet-integration sub-check in `TestRunStatusIntegratingHandbackDoesNotRecommendRun` (review F1): ✗ with the quiet-boundary condition removed from the new arm (mutation, restored) → ✓ with it.
- Full `^TestRunStatus` (12 tests) exit 0 after, including the untouched `TestRunStatusLandedHandbackIsClassC3` and `TestRunStatusWatchFitsTwentyLinesForEightOpenRequests`.

**New tests added:**
- `TestRunStatusMissingRunFolderIsNamedInTextAndJSON` and `TestRunStatusIntegratingHandbackDoesNotRecommendRun` in `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go`; fixture field `lastActivityPhase` (empty keeps `implementation`).

**Heavy verification plan:**
- Range: d2991804a6032d0ce10dd0de313395118f95b29d..bf9ed3d6b25c480346aeceac0511c5420fb07024 (same four lanes as the first plan at 5d4e136c)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the three runstatus Go files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all six changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three runstatus Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three runstatus Go files matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 95%** | 2026-10-10T20:56:03Z (first review 93% at 5d4e136c; re-review of the fix delta 5d4e136c..bf9ed3d6 at 95%)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. Both REQ-690 review gaps are closed: a missing `--run` folder is named and its run fields are absent, and a C3 row an integrator is working on no longer recommends `do-work run`. Full report: `do-work/runs/work-2026-10-10-192201/REQ-695-review.md`.

**Important findings:** None

**Minor findings:**
- F1: the 20-minute boundary on the integrator C3 arm was unpinned (boundary-drop mutation survived). Fixed in `3ea1e4ae` (quiet-integration sub-check), reviewer confirmed the mutation now fails. Resolved.
- F2: `integratorPhases` has no lock-step check against the board's labels; a Go import is impossible (separate module, `package main`). Kept report only (D-11): a drift falls back to the pre-REQ advice. — impact-negligible → report only
- F3: `actions/fan-out-reference.md` stall loop said every C3 row means "brief the next integrator". Fixed in `3ea1e4ae` (D-10). Resolved.
- F4: `actions/status.md` Step 4 relay list omitted the `run:` line. Fixed in `3ea1e4ae`. Resolved.
- F5: `--watch` still never shows a missing run folder (the builder's discovered task). — impact-user-visible → report only
- F6 (nit): prose in the `run_directory` path field (D-03); a permission error reads "not found" (D-07); the `IsDir` check and the three later phases are unpinned. — impact-negligible → report only
- F7 (nit): `status.md` said "phase integration or later". Fixed in `3ea1e4ae` (names the four phases). Resolved.

**Acceptance:** Pass. Focused `^TestRunStatus` 12/12 exit 0 at `bf9ed3d6`; mutations: new arm disabled and boundary dropped are both killed.
**Restatement sweep:** `status.md`, `status-guide.md` and `fan-out-reference.md` now agree on C3 and the `run:` line; no reader of `run_directory` outside `internal/runstatus` and `resultmodel`. Not the wave's last integration, so no wave-end sweep.
**Anti-bloat:** 0 unnamed additions; no new options, files or decorative tests.
**Suggested testing:** 3 items (live `do-work status --run do-work/runs/typo`; live status during an integrator's review phase; `--run` at a regular file and an unreadable folder)
**Follow-ups created:** None (remaining findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Passing an empty folder to every reader when the stat fails reused the existing "no run folder" branch, so the run fields became absent with no change in the three readers. Keying the integrator arm on the board's own phase label plus the existing quiet boundary needed no new class.
**What didn't:** The first test pinned only that an integrating row stops recommending `do-work run`; dropping the 20-minute boundary still passed, which would leave a dead integrator's REQ stuck on "an integrator is working" forever (review F1).
**Worth knowing:** A guard that suppresses a remedy needs a test on both sides of its time boundary. A report that resolves a path the user typed must print which path it read and whether it exists; a missing folder should read as absent plus a visible note, never as a false verdict on what the folder would hold.

## Orientation

Now `do-work-cli run-status` (`internal/runstatus/`, see `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`) names the run folder it read and flags a missing one, and its C3 class has an integrator-at-work variant with no `next_argv`; `actions/status.md`, `docs/status-guide.md` and `actions/fan-out-reference.md` describe it. Primes checked: both `prime_files` paths still exist.

## Heavy Verification Plan

- Base: d2991804a6032d0ce10dd0de313395118f95b29d
- Target: bf9ed3d6b25c480346aeceac0511c5420fb07024
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the three runstatus Go files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all six changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three runstatus Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three runstatus Go files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: bf9ed3d6b25c480346aeceac0511c5420fb07024
- Execution revision: bf9ed3d6b25c480346aeceac0511c5420fb07024 (detached drain checkout, `QUEUE_KANBAN_BROWSER` set, 1-minute load 8.28 at start). An earlier drain at the first merge `5d4e136c` was also green (72 s, 40 s, 69 s, 32 s).
- do-work-cli-integrations: exit 0, executed, 76 s
- staged-skills: exit 0, executed, 54 s
- updater: exit 0, executed, 69 s
- installer: exit 0, executed, 36 s

## Timing

Observed 2026-10-10T20:07:04Z to 2026-10-10T20:55:26Z: 48m 22s total, 57m 51s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 44m 40s | 1 |
| verification-gate | 12m 58s | 4 |
| handback-merge | 13s | 1 |

Slowest stage: review / review, fix on builder branch, delta re-review (spans the re-gate and re-drain), 44m 40s, outcome success.

Notes: no builder-work event, because the hand-back had already landed when the integrator started (dispatch 2026-10-10T19:27:42Z, builder commit 2026-10-10T19:28:50Z). The review event has no end flag, so it spans the review, the fix on the builder branch, the second gate and drain, and the delta re-review; it overlaps two verification-gate events, which is why attributed time exceeds the total.
