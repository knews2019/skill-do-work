---
id: REQ-695
title: 'Addendum: run-status names a missing --run folder and stops recommending do-work run while an integrator works'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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
