---
id: REQ-633
title: 'Panel B and the calibration log exclude by largest stamp gap, not raw span'
status: completed
route: B
estimate:
  p50_active_minutes: 40
  confidence: medium
  basis:
  - Route B
  - 18-file write set
  - 3 subsystems involved
  - 6 acceptance criteria
  calculated_at: 2026-10-05T20:56:04Z
created_at: 2026-10-05T20:11:59Z
user_request: UR-135
review_at: 2026-10-05T21:50:42Z
integration_at: 2026-10-05T21:45:38Z
builder_handback_at: 2026-10-05T21:45:31Z
dispatch_at: 2026-10-05T21:30:07Z
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: ["REQ-632"]
batch: board-activity-evidence
depends_on: [REQ-632]
write_set: ["skills/do-work-board/tools/queue-kanban/durations.go", "skills/do-work-board/tools/queue-kanban/durations_test.go", "skills/do-work-board/tools/queue-kanban/generate.go", "skills/do-work-board/tools/queue-kanban/generate_test.go", "skills/do-work-board/tools/queue-kanban/timeline.go", "skills/do-work-board/tools/queue-kanban/timeline_test.go", "skills/do-work-board/tools/queue-kanban/web/board-durations.js", "skills/do-work-board/tools/queue-kanban/web/board-timeline.js", "skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js", "skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go", "skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go", "skills/do-work-board/docs/board-guide.md", "skills/do-work/actions/estimate-reference.md", "skills/do-work/actions/work-reference.md", "skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go", "skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go", "skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go", "_dev/backfill-calibration-gap-column.py", "skills/do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go", "_dev/tests/contracts/queue-kanban.sh"]
claimed_at: 2026-10-05T20:55:46Z
commit: 0951df366d67a642b0ee1555d5c65cd871748343
heavy_verified_at: 2026-10-05T21:55:50Z
heavy_verified_revision: 0951df366d67a642b0ee1555d5c65cd871748343
completed_at: 2026-10-05T21:56:01Z
release_at: 2026-10-05T21:56:01Z
---
# Panel B and the Calibration Log Exclude by Largest Stamp Gap, Not Raw Span

## What
Replace the 4-hour raw-span exclusion with a rule on the largest gap between consecutive lifecycle stamps, 2 hours, applied identically by the two readers that claim one definition (`skills/do-work/actions/estimate-reference.md` → Calibration and `skills/do-work-board/tools/queue-kanban/durations.go`), rename the verdict string from `paused` to `idle-gap` in every reader, and add a `max_stamp_gap_minutes` column to `do-work/calibration-log.tsv` so the next re-fit can apply the same rule.

## Why
The 4h rule throws away long continuous sessions. Two of 190 were excluded at fit time; in the maintainer's repo three of the last five would be. As Route C runs grow, the raw-span rule discards the most informative samples and the timeline forecast inherits the bias. The calibration log carries only `wall_minutes`, so without a gap column the next re-fit cannot apply a gap rule at all. The maintainer chose one constant, 2h, with the phase named in the label, and asked that the archive's per-phase distribution be checked before the number is fixed.

## Detailed Requirements
1. Archive check first. Before changing code, compute the per-phase p95 of stamp gaps over every REQ in `do-work/archive/` using the intervals `buildPhaseBreakdown` emits, and record the table in this REQ's plan. If any phase other than dispatch → handback shows real multi-hour holes, stop and raise an Open Question proposing a second tier (builder phase 2h, other phases 45 min) instead of coding one constant. Otherwise proceed with one constant.
2. Rule (`durations.go`): `activityGapCeiling = 2 * time.Hour` replaces `analysisOutlierCeiling`. `dayMedianExclusionReason` takes the largest gap between consecutive parsed lifecycle stamps (the whole claimed → completed span when a REQ has no phase stamps) and returns `idle-gap`, `reversed`, or empty. Stamps only, on purpose: the calibration-log appender runs inside the lifecycle transaction and must not shell out to git, and both readers must agree byte for byte.
3. Rename the verdict string everywhere it is read: `generate.go` (the `implementationSpanReason` comment and the two assignment sites), `web/board-durations.js` (summary sentence, quantile skip, axis title, mark tooltip, day tooltip, table cell), `web/board-timeline.js` (forecast exclusion and assumptions text), `web/board-user-request-summary.js` (excluded count), `timeline.go` comments, `skills/do-work-board/docs/board-guide.md`. Panel B prose becomes "largest idle gap over 2h · excluded from day medians". The "2h" reaches the client as text shipped from Go, never as a number the client re-applies (lesson REQ-219).
4. Doc: rewrite the read-time rule in `estimate-reference.md` → Calibration to "exclude a span when the largest gap between consecutive lifecycle stamps exceeds 2h, or the span is negative", add one sentence naming `durations.go` as its second reader, and note that the 188-of-190 figure predates this rule. Update the `durations.go` header comment that calls itself the second reader so the two say the same thing.
5. Calibration log: add `max_stamp_gap_minutes` as the last column. Both appenders (`requeststate/state_plan.go` and `state_apply.go`) compute it from the REQ's stamps at write time. Ship a one-off idempotent backfill that fills old rows from archive frontmatter (location is the builder's call: `skills/do-work/scripts/` if it ships, else `_dev/`). The header row is the version marker; no separate marker.
6. Tests: a boundary pair straddling 2h by one minute (lesson REQ-374: only a straddling pair catches a second ceiling); a 4h21m span with a 67-minute largest gap is kept; a Route A REQ with no phase stamps falls back to the full span; a new TSV row carries the column; the backfill is idempotent on a second run; every JavaScript behaviour test that asserted `paused` now asserts `idle-gap` and the new prose.
7. Release: changelog entry and version bump per `_dev/primes/prime-releases.md`.

## Constraints
- Panel B reads archive frontmatter, not the TSV, so its medians shift as soon as the rule lands. That is accepted.
- The card-side gap from REQ-632 (Board cards show last correlated activity) uses stamps ∪ commits; this rule uses stamps only. A commit can only shrink a gap, so Panel B may exclude a REQ whose drawer shows a smaller gap. Nothing on the card shows the exclusion, so nothing visibly disagrees; say this once in the doc.
- `reversed` keeps its name and meaning.
- Shipped files change in three packages (board, core action, core CLI), so this is a release; the backfill script, if shipped, follows `_dev/primes/prime-shell-commands.md`.

## Dependencies
Depends on REQ-632 (Board cards show last correlated activity and drop the assumed-pause badge): the badge text field must already be gone, and the drawer row this REQ's wording refers to must exist.

## Builder Guidance
Certainty is high on the shape (one constant, stamps-only, both readers, new column) and medium on the number: 2h is the maintainer's pick pending the archive check in requirement 1. Builder latitude: the backfill script's language and location, the exact Panel B sentence wording, and whether `ImplementationSpan.ExclusionReason` is renamed in Go or only its values.

## Red-Green Proof
**RED prompt/case:** A completed REQ with a 4h21m span and a 67-minute largest stamp gap, and a second completed REQ with a 2h01m gap between `dispatch_at` and `builder_handback_at` inside a 3h span.
**Why RED now:** Panel B excludes the first (`excludedReason: "paused"`) and keeps the second; `do-work/calibration-log.tsv` has no column from which a re-fit could tell them apart.
**GREEN when:** The first is kept and the second is excluded with reason `idle-gap`; the Panel B summary names "largest idle gap over 2h"; the TSV header ends with `max_stamp_gap_minutes` and every row has a value; `grep -rn '"paused"' skills/do-work-board/tools/queue-kanban/*.go skills/do-work-board/tools/queue-kanban/web/*.js` returns nothing; `estimate-reference.md` and `durations.go` state the same rule.
**Validation:** User confirmed (plan approved 2026-10-05)

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: REQ-374 (only a straddling pair catches a second ceiling) and REQ-219 (ship the verdict, never the inputs) govern this change and are restated in requirements 3 and 6.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7095 tokens, over the 2000 budget; `slugged: partial`). Matched: changing queue-kanban model and UI.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18027 tokens, over the 2000 budget; `slugged: partial`). Matched: this REQ changes the calibration-log appenders in `requeststate`.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over the 2000 budget; `slugged: partial`). Matched: a shipped backfill script is shipped shell.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One gap definition in two places: the board's largest gap over its phase list without release, and the CLI's over its own stamp list plus completed_at; sorted by time, largest consecutive difference. One shared CLI row formatter decides the row shape from the log's first line. The rule's words ship from Go as text. (From the builder hand-back.)
- [x] **[APPLY]:** As planned; the write set grew by the shared formatter file, one finalization test file (D-04) and one maintainer contract check (D-02), recorded in Scope.
- [x] **[UNIFY]:** git diff --stat 22 files; gofmt and go vet clean in both modules; board package with JavaScript probes ok, whole CLI module ok; contract check and shipped-reference contract pass; no debug artifacts. Orchestrator cross-checked APPLY against git diff --stat 27e8248b..0951df36 (22 files, all in the updated Scope).

## Full Context
See `do-work/user-requests/UR-135/input.md` for complete verbatim input.

*Source: C. The Panel B calibration exclusion can keep its own rule, but it should be decided from the same gap evidence, not from the raw span … Both readers (estimate-reference.md and durations.go) must change together — the doc calls itself the single definition.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The REQ fixes the rule, the constant, the renamed verdict, every reader and the new column decision by decision (UR-135 session), so no Plan agent is needed. It spans three packages (board, core estimate doc, core CLI appenders) plus a one-off backfill, and its first requirement is an archive measurement, so exploration runs that measurement and locates the appenders and every reader of the verdict before dispatch.

**Planning:** Not required

## Open Questions

- [x] The archive check found multi-hour idle holes outside the builder phase (integration to review, review to completed). Should Panel B use one 2h largest-gap rule, the two-tier fallback (2h builder phase, 45 min elsewhere), or should this REQ be cancelled and the 4h span rule kept? → One 2h rule for every phase (recommended option).

Answered 2026-10-05 by the maintainer in the work session, after seeing the archive numbers: 539 measurable REQs; today's 4h span rule excludes 36; one 2h gap rule excludes 51 (the same 36 plus 15, 12 of them REQs with no phase stamps judged on their whole span) and keeps none the 4h rule drops in this repo; the two-tier rule would exclude 126. Reasoning: the post-merge holes are real idle time, so excluding them is correct, and long continuous runs in consumer repos (the brief's 4h21m REQ-2248 with a largest gap near 65 minutes) are kept. Out of scope: the second tier, and any per-phase threshold.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors and the archive table: `do-work/runs/work-2026-10-05-203159/REQ-633-exploration.md`. The findings that change what the builder writes:

- **Requirement 1 archive check: the stop condition fired** and the maintainer chose one 2h rule (Open Questions above). Per interval over 539 measurable REQs: dispatch to builder handback 1 of 72 over 2h (p95 43 min, and that one is a rewritten claim stamp); integration to review 9 of 67 over 2h (p95 712 min); review to completed 4 of 56. Today's 4h rule excludes 36; the 2h gap rule excludes 51 and keeps none of the 36.
- **A third calibration row builder exists:** finalization/finalization_discovery.go:470-491 rebuilds the row and requires the new file to be exactly the old file plus that row. If it is not updated, finalization cannot account for the log append, and any header rewrite breaks its prefix proof. One shared formatter for the appender (state_plan.go:400-435), the verifier (state_apply.go:479-506) and this proof.
- **work-reference.md:263 must change:** it says phase stamps never change calibration.
- **The core CLI has no list of phase stamp names;** the appender needs one, read through the typed record's field evidence. completed_at is the transaction's now at write time.
- **Definition to pin in both readers:** exclude release_at (phaseMilestonesOf lists it after Completed); 18 REQs have stamps out of declared order, and the 2h verdict is the same either way while the gap number differs.
- **The timeline renders the raw verdict string** (board-timeline.js:918 and :955), so the rename reaches the UI as "idle-gap excluded" automatically.
- **Calibration log:** 308 rows, all 5 columns, every row maps to one archived REQ, so a backfill fills all of them (43 get a gap over 120 min). Leave wall_minutes alone (24 rows already disagree with frontmatter). Board verify.go:1577 reads only columns 0 and 3.
- **Backfill:** no precedent for one-off migrations; shipped scripts are thin launchers into Go, and shell date parsing differs between macOS and Linux. Python in _dev/, skipping when the header already ends with the new column. Consumer repos keep the old header, so the appender's behaviour for an old header must be decided.
- **Tests:** board targeted run with JavaScript probes 22s; requeststate 5s; finalization package about 32s across its two biggest files (each under 30s). Tests asserting the verdict or the prose are listed in the findings file (javascript_behavior_a_test.go:619, :1394; javascript_behavior_d_test.go:594; generate_test.go:3702, :3709, :2456; durations_test.go:475-517, :744-788).

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — activityGapCeiling replaces analysisOutlierCeiling; the exclusion reads the largest gap between chronologically sorted lifecycle stamps (release excluded), returns idle-gap, reversed or empty; header comment
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — straddling pair at 2h, 4h21m kept, no-phase fallback, ceiling-derived tests rewritten
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — verdict comments and values
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — REQ-909 kept, REQ-902 and paused fixtures renamed
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modify) — comments
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modify) — only if a fixture's verdict changes
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — summary sentence, axis title, mark tooltip, table cell
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — forecast assumptions text
- `skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js` (modify) — comment and excluded wording
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modify) — Panel B and forecast prose assertions
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (modify) — fixture verdict
- `skills/do-work-board/docs/board-guide.md` (modify) — excluded wording; one sentence on stamps-only vs stamps plus commits
- `skills/do-work/actions/estimate-reference.md` (modify) — Calibration rule, column list, second-reader sentence, the 188-of-190 note
- `skills/do-work/actions/work-reference.md` (modify) — the line that says phase stamps never change calibration
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modify) — appender writes max_stamp_gap_minutes when the header carries it
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modify) — verifier rebuilds the same row
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modify) — verifier and old-header cases
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modify) — third row builder uses the shared formatter
- `_dev/backfill-calibration-gap-column.py` (new) — one-off idempotent backfill
- `skills/do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go` (new) — the shared calibration row formatter used by the appender, the verifier and the finalization proof
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modify) — the finalization proof under both header shapes (D-04)
- `_dev/tests/contracts/queue-kanban.sh` (modify) — pins the CLI's stamp list to the board's phase list (D-02)

**Files I will NOT touch:** do-work/calibration-log.tsv (the orchestrator runs the backfill in the main tree after this REQ's finalization, because REQ-632's finalization appends a row to it first), web/board-cards.js and activity_correlation.go (REQ-632's card and drawer evidence, stamps plus commits, stays as is), board verify.go (reads only columns 0 and 3, safe with an appended column), every release path (written by finalization), do-work/lessons-index.md.

**Acceptance criteria (restated from REQ):**
- [ ] A completed REQ with a 4h21m span and a 67-minute largest stamp gap is kept; a REQ with a 2h01m gap inside a 3h span is excluded with reason idle-gap; a pair straddling 2h by one minute pins the constant
- [ ] A REQ with no phase stamps falls back to the full claimed-to-completed span
- [ ] No board Go or JS file in the package or web folder ships the verdict string "paused"; Panel B says largest idle gap over 2h, shipped from Go as text
- [ ] estimate-reference.md and durations.go state the same rule; work-reference.md no longer says phase stamps never affect calibration
- [ ] A new calibration file's header ends with max_stamp_gap_minutes and each new row carries the value; all three row builders (appender, verifier, finalization proof) agree through one formatter; an old-header file gets old-shape rows
- [ ] The backfill fills every row of this repo's log and is a no-op on a second run

## Pre-Flight

**Git:** ⚠ 1 pre-existing uncommitted file outside do-work/ — `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, REQ-632's pending lesson entry (its finalization commits it); not in this REQ's write set, excluded from this REQ's evidence. Integration tip 9dacf899 (REQ-633 claim), which contains REQ-632's merge b004ecd6.
**Tests baseline:** ✓ focused board and requeststate tests green (`do-work/runs/work-2026-10-05-203159/REQ-633-probe.sh`, launched by advance)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 9dacf899, both Go stages EXECUTING, gate wall 118s; green-gate record satisfied
**Dependencies:** ✓ Go toolchain and python3 present; no new module dependency planned

*Checked by work action*

## Implementation Summary

**Files changed:**
- `_dev/backfill-calibration-gap-column.py` (new)
- `_dev/tests/contracts/queue-kanban.sh` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js` (modified)
- `skills/do-work/actions/estimate-reference.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go` (new)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modified)

**What was done:** Panel B's day-median exclusion now reads the largest gap between chronologically sorted lifecycle stamps (release excluded; the whole claimed-to-completed span when a REQ has no phase stamps) against one 2h constant, with verdict idle-gap or reversed; the rule's words ship from Go as text and every reader of the old paused verdict is renamed. The calibration log gains max_stamp_gap_minutes, written by one shared formatter that the appender, the verifier and the finalization proof all use, with the log's header deciding the row shape so existing five-column logs keep five-column rows. estimate-reference.md and durations.go state the same rule, work-reference.md no longer says phase stamps never affect calibration, and a maintainer-only Python backfill fills the column idempotently. Merge range 27e8248b..0951df36 (builder commit 43024f7e, merge 0951df36).

## Qualification

**Diff range:** 27e8248b..0951df36 (builder commit 43024f7e, merge 0951df36)
**Gate records:** qualify satisfied; scope-drift satisfied after Scope and write_set were brought in line with the hand-back (two declared test files were not needed; the formatter file, one finalization test file and one maintainer contract check were added, D-02 and D-04).
**Warnings judged:** QUALIFY-NEW-FILE-UNWIRED on `_dev/backfill-calibration-gap-column.py` — expected: a standalone maintainer script run by hand, not imported. Seven QUALIFY-REPORTER-OUTPUT lines in the same script — its intended report to the person running it (row counts, dry-run notice, refusal on an unexpected header), not leftover debugging.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. Requirement 1's archive table and the maintainer's answer are in Open Questions and Exploration. One 2h constant; the largest gap over chronologically sorted stamps with release excluded and a whole-span fallback; verdicts idle-gap and reversed; the rule text ships from Go; every listed reader renamed; the estimate-reference rule, column list and second-reader sentence; the work-reference line corrected; one formatter for all three row builders with the header deciding the shape; an idempotent backfill with a dry run. The live log was not backfilled on the branch, as instructed; the orchestrator runs it after finalization.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back; APPLY cross-checked against git diff --stat 27e8248b..0951df36 (22 files, all in the updated Scope).
**Live data flow:** dayMedianExclusionReason feeds both implementationSpanReason and excludedReason in generate.go; the web readers take the shipped durations.exclusionRule text; requeststate and finalization call the one formatter.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 0951df36
**Result:** ✓ All passing — exit 0, gate wall 120s; stage queue-kanban-fast-tests (420 tests, wall 36s, slowest file 18.58s < 30s); stage do-work-cli-fast-tests (870 tests, wall 59s, slowest file finalization_recovery_test.go 20.25s < 30s, up from about 19.7s with the new proof test). The contract checks include the new stamp-list pin. Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-05-203159/REQ-633-probe.sh` (board Durations, Timeline, ImplementationSpan, UserRequestProgress, DoneCard and Calibration tests with JavaScript probes on, plus the requeststate package) → exit 0, advance probe record satisfied.

**Red-green validation:** traced to `## Red-Green Proof`; RED taken by the builder against a compiling stub, GREEN at builder commit 43024f7e:
- TestImplementationSpanVerdictReadsTheLargestStampGap (the RED pair, the straddling one-minute pair at 2h, the no-phase fallback, a rewritten claim): ✗ the 4h21m span with a 67-minute gap read "paused", want kept; the 2h01m phase gap inside a 3h span read kept, want "idle-gap" → ✓
- TestGeneratedRequestCarriesTheDoneCardImplementationSpan: ✗ REQ-909 "paused", want kept; durations.exclusionRule empty → ✓
- TestJavaScriptBehaviorDurationsHeadlineRollingMedianAndCadenceTicks and TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions: ✗ old assumed-pause prose → ✓ new prose
- TestCompleteCalibrationRowFollowsTheLogHeader: ✗ a new log got the five-column header and a six-column log got a five-column row → ✓ (the five-column-log case passed before and after: it is the old-header guard)
- TestArchivedCalibrationVerifierReadsTheHeaderShape: ✗ six-column row under a six-column header rejected → ✓
- TestCalibrationAppendProofFollowsTheLogHeader: ✗ a six-column row did not prove, a five-column row under a six-column header did → ✓
- Backfill fixture test: ✗ log unchanged against an empty stub → ✓ first run fills the column and leaves wall_minutes alone, second run changes nothing, dry run writes nothing

**New tests added:**
- the four Go tests above plus the backfill fixture test (scratchpad), and the stamp-list contract check in `_dev/tests/contracts/queue-kanban.sh`

**Existing tests updated (cross-REQ impact):**
- durations_test.go and generate_test.go (from REQ-632 and earlier Panel B work): paused fixtures now read idle-gap or kept under the new rule — intentional
- javascript_behavior_a_test.go and javascript_behavior_d_test.go: verdict string and Panel B prose — intentional

**Render evidence (builder):** headless Chrome on the Durations page, 30-day window: "Panel B excludes 32 spans from its medians (largest idle gap over 2h, or a negative span from a broken stamp)"; the timeline group reads "1 idle-gap excluded". All-time on this repo: 36 excluded before, 51 after, none moved from excluded to kept.
**Backfill dry run on the main tree (builder):** 309 rows, 309 filled, 43 with a gap over 120 minutes, 0 unresolved; file unchanged.

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 27e8248b..0951df36
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban and an unmapped _dev script
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — files changed under skills/do-work-board/tools/queue-kanban and an unmapped _dev script
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage uncertain for the unmapped _dev/backfill-calibration-gap-column.py
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage uncertain for the unmapped _dev/backfill-calibration-gap-column.py
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage uncertain for the unmapped _dev/backfill-calibration-gap-column.py
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage uncertain for the unmapped _dev/backfill-calibration-gap-column.py

*Verified by work action*

## Review

**Overall: 94%** | 2026-10-05T21:49:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 93% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- Rounding at the boundary differs between the two readers: the board excludes when the exact gap is over 2h (`durations.go:284`, `largestStampGap > activityGapCeiling`), while the CLI writes the gap floored to whole minutes (`calibration_row.go:51`) and `estimate-reference.md:94-96` gives no comparison for the column. A re-fit using `max_stamp_gap_minutes > 120` keeps a gap in [120m01s, 121m) that Panel B excludes. On this repo all 309 backfilled rows agree with the board's verdict, so it is theoretical today. Fix with one doc clause (`> 120` on the floored column means "121 minutes or more"), or compare in seconds. — impact-negligible → report only
- `estimate-reference.md` tells the next re-fit to apply the gap rule, but by design an existing five-column log never gets the column: consumer repos keep the old header and the backfill is in `_dev/`, so it does not ship. The doc does not say what a re-fit does with five-column rows (for example, fall back to the span rule or skip them). The header-decides-shape choice is documented in the doc and the formatter comment, but it is not among the hand-back's D-01..D-05. — impact-rule-change → report only
- The REQ's Constraint asked for the stamps-only vs stamps-plus-commits sentence "once in the doc". It appears in two shipped docs, `skills/do-work-board/docs/board-guide.md:57` and `skills/do-work/actions/estimate-reference.md:96`, plus the `durations.go` header comment. Keep one doc copy and point to it from the other. — impact-negligible → report only
- Operational trap that nothing documents: `calibrationAppendProves` (`finalization_discovery.go:470-479`) requires the working-tree log to start with the HEAD bytes. The backfill rewrites every line, so if it runs and is not committed on its own before the next REQ completes, that REQ's finalization cannot prove its calibration append and refuses. The script's docstring should say: run it, then commit it alone before the next completion. — impact-user-visible → report only
- Stale restatement: the comment at `skills/do-work-board/tools/queue-kanban/generate_test.go:2194` still says the sweep covers "everything under the four-hour read-time ceiling that admits a sample at all". Under the gap rule, spans over 4h are admitted, for example REQ-909 at 261 minutes. Behaviour is unaffected because Panel B clamps medians at 45 minutes. — impact-negligible → report only
- Nit: the comment at `durations_browser_probe_test.go:797` still calls its fixtures "paused spans". — impact-negligible → report only

**Acceptance:** Pass — I built the board at 0951df36 and ran generate against the repo. Result: 540 samples, 51 `idle-gap`, 489 kept, `durations.exclusionRule` = "largest idle gap over 2h", and no `"paused"` in board-data.js. The backfill dry run wrote nothing (md5 unchanged). On a scratch copy the backfill is idempotent and leaves columns 1-5 byte-identical. Its gaps agree with the board's verdict on 309 of 309 rows.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

### Review detail (orchestrated report)

**Approve** — The gap rule, the rename and the new log column all behave as specified. The two gap definitions agree on real data. The six findings are minor documentation and comment gaps.
Route B | 43024f7e merged as 0951df36

#### Check-by-check
1. Gap definitions agree:
   - Stamp set: board `phaseMilestonesOf` minus release_at; CLI `CalibrationGapStampFields` plus completed_at. The contract check pins the two lists and passes.
   - release_at is excluded on both sides.
   - Both sides sort by instant and take the largest consecutive difference.
   - Fallback: a phaseless REQ falls back to the whole claimed→completed span on both sides.
   - Both use the same parser layouts.
   - The only difference is the boundary rounding (Minor 1).
   - Cross-check: the backfilled column agrees with board verdicts on all 309 log rows.
2. One formatter: `requestmodel.FormatCalibrationRow` is called by the appender (`state_plan.go:408,425`), the verifier (`state_apply.go:497`) and the finalization proof (`finalization_discovery.go:483`).
   - Header-decides-shape is tested for a new log, a 6-column log and a 5-column log through the command seam (`TestCompleteCalibrationRowFollowsTheLogHeader`, which includes the verifier), plus `TestArchivedCalibrationVerifierReadsTheHeaderShape`.
   - The prefix proof is tested for an existing 6-column and 5-column log (`TestCalibrationAppendProofFollowsTheLogHeader`).
   - The new-log proof gap predates this REQ and is reported in the hand-back.
3. Restatement Sweep: I grepped skills/, _dev/primes, _dev/tests and the docs for paused / assumed pause / four-hour / 4h / analysisOutlierCeiling / implementationSpanReason / excludedReason / calibration-log columns.
   - Updated: Go, all three web readers, board-guide, estimate-reference (rule, columns, second reader, 188-of-190 note) and work-reference:263.
   - Correctly left alone: the remaining "paused" hits are unrelated (estimate scope, hold folder, the four-hour heavy-lane reuse) or quote history (the REQ-374 lesson, the timeline_test Red-Green quote).
   - Stale: two test comments (Minor 5 and the Nit).
   - Board `verify.go:1616` reads at least 4 columns, so it is safe with 6.
4. The Constraint sentence is present, but twice (Minor 3).
5. Backfill: idempotent through the header check, wall_minutes untouched (verified), the dry run writes nothing, and it uses an atomic replace.
6. Tests: each named test states its pin. The straddling pair is derived from `activityGapCeiling` (−1m / at / +1m). Also covered: 4h21m with a 67-minute gap is kept, the phaseless fallback, release_at ignored, and the rewritten-claim sort. The vacuity guard covers idle-gap. The builder recorded RED for all of them.
7. Release (changelog and version): left to finalization, as the hand-back states.

Focused tests rerun: requestmodel ok, requeststate ok (6.4s), the finalization proof test ok, board Duration/ImplementationSpan/Timeline/JS subsets with JavaScript probes ok (7.3s), and `_dev/tests/contracts/queue-kanban.sh` passed.

#### Suggested Additional Testing
- After finalization, run the backfill without `--dry-run` and commit it alone. Then complete one REQ and confirm the new row has 6 columns and finalization proves the append.
- Open the Durations page and check the axis title width and the mark tooltip text at phone width.
- Optional: one board test that pins the agreement between the board verdict and a floored CLI gap at 120m30s, if Minor 1 is resolved in favour of exact comparison.

#### Scores
| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Requirements 1-6 delivered; requirement 7 (release) is owned by finalization |
| Code Quality | 93% | One formatter, the rule ships as text, clear comments; the rounding contract is not stated |
| Test Adequacy | 90% | Straddling pair, seam tests at all three writers; no cross-reader test at the boundary |
| Scope | 95% | Three added files recorded in Scope (D-02, D-04) |
| Risk | Low | Fail-closed finalization if the backfill is left uncommitted |
| Acceptance | Pass | Generated board shows 51 excluded and the shipped rule text |

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-05-203159/REQ-633-handback.md`).

- D-01 (DECIDE & STATE): the verifier and the finalization proof parse the estimate as an integer through the shared formatter instead of comparing the raw string, so a padded estimate cannot make the appender and the verifier disagree.
- D-02 (DECIDE & STATE): the CLI's stamp list is pinned to the board's phase list by a grep in `_dev/tests/contracts/queue-kanban.sh`, not a Go test, so no shipped test reads a sibling module's source. It runs only in the maintainer gate.
- D-03 (DECIDE & STATE): Panel B wording — summary "largest idle gap over 2h, or a negative span from a broken stamp"; table and marks "excluded from day median — largest idle gap over 2h"; the axis title stays short ("idle-gap and broken spans excluded"); the timeline reads "Idle-gap and reversed spans are excluded from both medians".
- D-04 (ESCALATE, resolved by the orchestrator): a finalization proof test was added to finalization_recovery_test.go, outside the original write set. Value: the third row builder is pinned under both header shapes. Risk: one file outside scope; easy to revert. Orchestrator judgment: accepted and added to Scope; the proof change needs its own test and that file is where the package's calibration-proof tests live.
- D-05 (DECIDE & STATE): backfill rows with no single archive match get an empty value and are reported, never guessed (0 on this repo).

Orchestrator decisions:
- D-06 (DECIDE & STATE): the log's header decides the row shape, and an existing header is never rewritten, because finalization's calibration proof requires the new log to be exactly the old bytes plus one row. Consequence (review finding on old logs): a consumer log with the five-column header never gains the column, and the next re-fit there falls back to the raw span for those rows. Reversible: a shipped migration could rewrite the header in a commit of its own.
- D-07 (DECIDE & STATE): the backfill runs in the main tree after this REQ's finalization and is committed on its own before any other REQ completes, because REQ-632's finalization appended a row first and the next finalization's proof needs the committed bytes as its prefix.
- D-08 (DECIDE & STATE): the six review findings stay report only per their impact tokens; none is critical.

## Discovered Tasks

From the builder's hand-back and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- The board excludes on the exact gap while the log stores whole minutes rounded down; the doc does not say to compare the column with "at least 121". — impact-negligible → report only
- A five-column consumer log never gains max_stamp_gap_minutes, and estimate-reference.md gives the next re-fit no fallback for those rows. — impact-rule-change → report only
- The stamps-only versus drawer sentence appears in board-guide.md and estimate-reference.md, not once. — impact-negligible → report only
- After a backfill the log must be committed alone before the next finalization, or that finalization's calibration proof refuses; nothing documents this. — impact-user-visible → report only
- Stale comments: generate_test.go:2194 (four-hour ceiling) and durations_browser_probe_test.go:797 (paused spans). — impact-negligible → report only
- Finalization's calibration proof cannot prove the creation of a brand-new log (the header is not part of the expected row); pre-existing. — impact-negligible → report only

## Lessons Learned

**What worked:** Measuring the archive before coding: the REQ's own stop condition fired, and the maintainer chose from real numbers (36 excluded today, 51 under one 2h rule, 126 under two tiers) in one question. Routing all three calibration row builders through one formatter closed the alternate-writer trap exploration found in finalization, which the REQ's write set had missed.
**What didn't:** The REQ's premise (the 4h rule throws away long continuous sessions) has no case in this repo; the gap rule excludes more here, and its benefit shows only in consumer repos with long continuous runs. The write set also missed work-reference.md and the finalization proof, both found only by exploration.
**Worth knowing:** A calibration-log header is never rewritten in place, because finalization proves the append as old bytes plus one row; the header is the version marker and decides the row shape for every writer. The backfill is a maintainer-only one-off and must be committed on its own.

## Orientation

Now Panel B's day medians and the calibration log judge a long REQ by its largest gap between lifecycle stamps (over 2h means idle-gap) instead of its raw span, and the log records that gap for the next estimator re-fit; lives in the board's durations rule and the core CLI's calibration writers (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`). [MAP CHANGED]: the calibration log has a sixth column decided by its header, written by one shared formatter (requestmodel/calibration_row.go) that the appender, verifier and finalization proof all call. Prime spot-check: the listed primes' referenced paths exist; none names the old 4h rule.

## Heavy Verification Plan

- Base revision: 27e8248b0b9dafd9397fb42b1b54ac2d55c8cf58
- Target revision: 0951df366d67a642b0ee1555d5c65cd871748343 (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — board files changed
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — board files changed
- do-work-cli-integrations — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage uncertain for the unmapped backfill script
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed
- updater — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage uncertain for the unmapped backfill script
- installer — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage uncertain for the unmapped backfill script

## Heavy Verification Result

- Target revision: 0951df366d67a642b0ee1555d5c65cd871748343
- Execution revision: 0951df366d67a642b0ee1555d5c65cd871748343 (detached drain checkout `.git/work-run-2026-10-05-203159/drain-head`, run 21:50:49Z to 21:55:36Z, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 7s
- queue-kanban-browser: exit 0, executed (not skipped), 78s
- do-work-cli-integrations: exit 0, executed, 67s
- staged-skills: exit 0, executed, 38s
- updater: exit 0, executed, 64s
- installer: exit 0, executed, 28s

Green: every selected lane present, exit 0, none skipped, none reused.

## Timing

Observed 2026-10-05T20:56:00Z to 2026-10-05T21:55:50Z: 59m 50s total, 1h 00m 59s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| exploration-preflight | 33m 14s | 1 |
| builder-work | 15m 24s | 1 |
| verification-gate | 7m 49s | 2 |
| review | 4m 24s | 1 |
| handback-merge | 8s | 1 |

Slowest stage: exploration-preflight / triage, estimate, exploration agent, archive check, maintainer question, scope, pre-flight gate, 33m 14s, outcome success.
