# REQ-633 Exploration — Panel B and the calibration log exclude by largest stamp gap

Scripts (scratchpad, not committed): `archive_gaps.py` (main table), `archive_gaps_variants.py` (declared vs sorted order, two-tier rule), `tsv_coverage.py` (backfill coverage). They mirror `durations.go`: `phaseMilestonesOf` order (durations.go:244-257), `buildPhaseBreakdown` (durations.go:264-310: nil when no optional stamp parses, declared order, no sorting, `release_at` counts as an optional stamp), `measureImplementationSpan` (durations.go:202-218: earliest origin-eligible stamp → `completed_at`), `isCompletedStatus` (model.go:1032), `parseTimestamp` layouts (model.go:1680-1696).

## F1. Requirement 1 archive check — the stop condition fires

Corpus: 608 REQ files under `do-work/archive/` (recursive), 558 terminal-success, 539 with a measurable span. 415 of 539 have no optional phase stamp (full-span fallback). Minutes, linear-interpolated percentiles.

| interval (buildPhaseBreakdown label pair) | n | p50 | p95 | max | >45m | >2h | negative |
|---|---|---|---|---|---|---|---|
| Claimed -> Planning | 40 | 3.7 | 45.2 | 62.1 | 2 | 0 | 11 |
| Claimed -> Dispatch | 42 | 5.8 | 24.9 | 40.2 | 0 | 0 | 2 |
| Planning -> Dispatch | 33 | 3.3 | 16.5 | 107.2 | 1 | 0 | 0 |
| Claimed -> Builder handback | 3 | 6.5 | 20.5 | 22.0 | 0 | 0 | 0 |
| **Dispatch -> Builder handback** | 72 | 11.2 | 42.7 | 403.6 | 4 | **1** | 1 |
| Builder handback -> Integration | 69 | 0.0 | 2.8 | 8.8 | 0 | 0 | 1 |
| Claimed -> Review | 6 | 2.5 | 98.2 | 126.5 | 1 | 1 | 0 |
| Planning -> Review | 1 | 78.4 | 78.4 | 78.4 | 1 | 0 | 0 |
| Builder handback -> Review | 1 | 6.6 | 6.6 | 6.6 | 0 | 0 | 0 |
| **Integration -> Review** | 67 | 11.7 | **712.5** | 999.9 | 14 | **9** | 1 |
| Review -> Remediation | 19 | 6.5 | 28.4 | 34.3 | 0 | 0 | 4 |
| Remediation -> Re-review | 15 | 7.6 | 28.9 | 67.7 | 1 | 0 | 0 |
| Claimed -> Completed | 33 | 48.3 | 449.3 | 450.1 | 17 | 9 | 0 |
| Planning -> Completed | 6 | 120.8 | 379.6 | 400.2 | 4 | 3 | 0 |
| Dispatch -> Completed | 3 | 41.2 | 42.1 | 42.2 | 0 | 0 | 0 |
| Builder handback -> Completed | 5 | 116.3 | 210.2 | 228.7 | 5 | 2 | 0 |
| Integration -> Completed | 2 | 731.0 | 961.8 | 987.4 | 2 | 2 | 0 |
| **Review -> Completed** | 56 | 4.7 | **476.4** | 869.6 | 8 | **4** | 0 |
| Remediation -> Completed | 4 | 37.0 | 274.4 | 316.1 | 1 | 1 | 0 |
| Re-review -> Completed | 15 | 2.0 | 31.2 | 46.6 | 1 | 0 | 0 |
| Completed -> Release | 105 | 0.0 | 1.4 | 5.9 | 0 | 0 | 0 |

Non-builder intervals (release excluded): n=420, 31 over 2h, 58 over 45 min.

**Verdict:** the multi-hour holes are NOT in the builder phase. Dispatch → Builder handback has one >2h interval in 72 (REQ-504, and that one is a reorder artifact: its `claimed_at` was rewritten at 22:27 after an integration at 16:31; sorted chronologically the hole is Integration → Claimed). The real holes sit after the work merges: Integration → Review (9 over 2h, p95 712 min, e.g. REQ-624 integration 2026-09-08T22:57 → review 2026-09-09T09:00, an overnight wait for review) and Review → Completed (4 over 2h). Requirement 1's stop condition ("any phase other than dispatch → handback shows real multi-hour holes") is met. The REQ's premise (builder phase long, others short) is reversed by the data.

### What each candidate rule does (539 samples, 0 reversed)

| rule | excluded | notes |
|---|---|---|
| current: span > 4h | 36 | |
| one constant: largest gap > 2h | 51 | all 36 current exclusions stay excluded; **0 long spans are rescued**; 15 newly excluded with span ≤ 4h (12 of them phase-less fallback REQs at 122–222 min) |
| two-tier: Dispatch→Builder handback > 2h, any other interval > 45 min | 126 | the Open Question's proposed alternative; excludes 2.5× the 2h rule |
| largest gap > 45 min (one constant) | 127 | |

Of the 51 the 2h rule excludes: 20 via the full-span fallback (no phase stamps), 9 Claimed→Completed (REQs with only `release_at`, e.g. the 2026-09-05 batch claimed 00:33–00:41 and completed 08:02–08:07), 8 Integration→Review, 4 Review→Completed, 3 Planning→Completed, 2 Integration→Completed, 2 Builder handback→Completed, 1 Remediation→Completed, 1 Claimed→Review, 1 Dispatch→Builder handback.

So the Why ("the 4h rule throws away long continuous sessions") has no instance in this archive: every span over 4h also has a gap over 2h. The new rule's effect here is purely stricter (36 → 51). The 4h21m / 67-min case exists only as the test fixture REQ-909 (generate_test.go:3659-3667).

Full per-REQ table (every REQ with largest gap > 45 min or span > 4h) is in the script output `archive_gaps.out` in the scratchpad; the top rows:

| REQ | route | span | largest gap | interval | current | 2h rule |
|---|---|---|---|---|---|---|
| REQ-483 | A | 1014 | 1000 | Integration -> Review | paused | idle-gap |
| REQ-485 | C | 1023 | 996 | Integration -> Review | paused | idle-gap |
| REQ-502 | A | 994 | 987 | Integration -> Completed | paused | idle-gap |
| REQ-507 | C | 965 | 921 | Integration -> Review | paused | idle-gap |
| REQ-475 | C | 900 | 870 | Review -> Completed | paused | idle-gap |
| REQ-624 | A | 615 | 603 | Integration -> Review | paused | idle-gap |
| REQ-504 | C | 422 | 404 | Dispatch -> Builder handback (reorder) | paused | idle-gap |
| REQ-592 | B | 277 | 229 | Builder handback -> Completed | paused | idle-gap |
| REQ-594 | C | 199 | 141 | Planning -> Completed | kept | idle-gap |
| REQ-552 | B | 153 | 136 | Builder handback -> Completed | kept | idle-gap |
| REQ-543 | B | 179 | 126 | Claimed -> Review | kept | idle-gap |
| REQ-457 | B | 123 | 123 | Claimed -> Completed | kept | idle-gap |
| REQ-554 | B | 156 | 116 | Builder handback -> Completed | kept | kept |
| REQ-486 | C | 180 | 107 | Planning -> Dispatch | kept | kept |

### Definition details the builder must pin

- **Declared order vs chronological sort:** 18 REQs have stamps out of declared order (rewritten claims, e.g. REQ-483 claim 2026-09-04T14:31 after a dispatch on 09-03). The 2h verdict is identical either way on all 539 (0 differ); the gap NUMBER differs for those 18. Pick one and state it in both readers and the CLI. Sorting makes a negative interval impossible and measures the real hole; declared order matches `buildPhaseBreakdown`. 18 REQs carry a negative phase interval with a non-negative span.
- **`release_at` must be excluded** from the gap rule. `phaseMilestonesOf` includes it after Completed (durations.go:255), so reusing `buildPhaseBreakdown` directly would add the Completed→Release interval (max 5.9 min here, harmless today, wrong in principle).
- **Origin before `claimed_at`:** 11 REQs (REQ-475, 483, 485, 502–507, 567, 572) have an origin stamp more than 2h before `claimed_at`. The board span starts at the earliest origin (durations.go:186-191); the TSV `wall_minutes` is `completed − claimed_at` (state_plan.go:415-416). Over the stamp set both stamps are included, so the gap rule is consistent; the wall figures already differ (pre-existing).
- **Fallback:** a REQ with no phase stamp has origin = `claimed_at`, so "full span" is the same number in the board and the TSV. A REQ with only `release_at` gets a breakdown of Claimed → Completed, which equals the fallback.

## F2. Readers of the verdict and the ceiling (skills/)

Go:
- durations.go:21-31 header + `analysisOutlierCeiling`; :41-43 `DayMedianExclusion` comment; :128 `ExclusionReason` comment; :216 call site passes `wallSpan` only — must pass the ticket (or its stamps) instead; :312-322 `dayMedianExclusionReason(wallSpan)` → new signature over stamps, returns `idle-gap`/`reversed`/"".
- generate.go:250-257 `implementationSpanReason` comment ("paused"); :265 field; :368-370 `excludedReason` comment; :864 and :924 assignment sites (values only change).
- timeline.go:250-251 and :343-345 comments ("paused or reversed span").
- activity_correlation.go:19 comment mentions "pause from the span alone" (REQ-632 text, probably fine).
- board-cards.js:70 reads only `"reversed"` — unaffected.

JS:
- web/board-durations.js:12-15 comment; :647-670 summary sentence ("over four hours is an assumed pause"); :772 quantile skip (truthy test, no string); :1054 axis title ("paused and broken spans excluded"); :1265-1268 mark tooltip (`=== "paused"`); :1434-1437 table cell (`=== "paused"`).
- web/board-timeline.js:874-876 counts by reason string; :918-919 renders `count + " " + reason + " excluded"` and :955-956 `reason + " excluded×"` — the raw verdict string reaches the UI, so renaming produces "idle-gap excluded" automatically; :1526 forecast assumptions text "Paused and reversed spans are excluded from both medians."
- web/board-user-request-summary.js:13 comment, :105-107 truthy check + comment "(an assumed pause, or reversed stamps)".

Docs:
- skills/do-work-board/docs/board-guide.md:41 "`N excluded` … (an assumed pause, or reversed stamps)". board-guide.md:57 (REQ-632's idle-gap drawer paragraph) is where the Constraints sentence about stamps-only vs stamps ∪ commits fits.
- skills/do-work/actions/work-reference.md:263 — **not in write_set, must change**: "None of these fields changes calibration: its only span remains the just-archived REQ's `claimed_at` → `completed_at`." The new column makes phase stamps feed calibration.

## F3. Calibration-log writers — there are three row builders, not two

- Writer: state_plan.go:400-435 `planCalibration`. Runs only on `TransitionComplete` with an estimate block; header written only when the file is empty (state_plan.go:426-428, header `req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at`); row at :433 `fmt.Sprintf("%s\t%s\t%d\t%d\t%s\n", id, route, estimate, wallMinutes, CanonicalTimestamp(now))`. `completed_at` is `plan.Options.Now`, not yet in the record. Phase stamps are not typed fields on the record; read them via `plan.Target.TypedRecord.FieldEvidenceByName["dispatch_at"].ScalarValue` (request_model.go:23-27, :91). The core CLI has **no list of phase stamp names** anywhere (grep for `builder_handback_at` in non-test core Go: none), so a new list is needed — a second copy of `phaseMilestonesOf`; pin the two together with a test.
- Verifier: state_apply.go:479-506 `verifyArchivedCalibrationEvidence` rebuilds the row from the archived file and requires `HasSuffix(plan.CalibrationBytes, wantRow)`.
- **Third, not in write_set:** finalization/finalization_discovery.go:470-491 `calibrationAppendProves` rebuilds the same row and requires the post-image to be exactly `before + want` (`HasPrefix(after, before)` at :476). If the row gains a column and this is not updated, finalization can no longer attribute the TSV append. Any header rewrite of an old file also breaks the prefix proof. The do-work-cli prime names this trap (`alternate-writer-contract-drift`). Recommend one shared row formatter (e.g. in requestmodel) used by all three.
- Header handling for an existing old-header file is undecided by the REQ ("the header row is the version marker"): today the appender never rewrites a header and a test writes a headerless file (state_apply_test.go:855-858). Consumer repos will keep the 5-column header unless something upgrades it.
- Tests pinning row/header: state_apply_test.go:78 (`"REQ-303\tC\t60\t60\t2026-08-31T21:00:00Z"`, Contains — still passes if a column is appended after it only if the assertion includes the trailing tab/newline; it does not, so it keeps passing); state_apply_test.go:855-933 (mode test); finalization_recovery_test.go:219, :533, :563; finalization_req499_test.go:627, :635; board verify_test.go:1851-1899, :2472.
- Other readers: board verify.go:1577-1675 `appendCalibrationLogFindings` skips lines starting `req_id\t`, needs ≥4 columns, reads only columns 0 and 3 — safe with an appended column. estimate-reference.md:94 lists the columns (update it). No reader in `_dev/` scripts (only lessons satellites mention the file). Run helpers under do-work/runs/*/helpers/make-manifests.py only name the path.
- Current file: do-work/calibration-log.tsv, 309 lines = header + 308 rows, all 5 columns, no duplicate ids. Every row maps to exactly one archive file with parseable `claimed_at`/`completed_at`, so a backfill fills all 308. 43 rows would get a gap > 120. 24 rows' `wall_minutes` already disagree with frontmatter (the known verify finding); the backfill must not touch `wall_minutes`.

## F4. Doc text to rewrite

- skills/do-work/actions/estimate-reference.md:90 `## Calibration`; :92 the 188-of-190 sentence ("excluding spans over 4 hours or negative"); :94 TSV column list + "outlier rule is applied when reading"; :98 "exclude spans > 4h or negative as assumed pauses". estimate-reference.md does not itself say it is the single definition; durations.go:25-27 says it ("The rule is stated once, in … estimate-reference.md → Calibration; this is its second reader, not a second definition").

## F5. Tests asserting "paused" or Panel B prose

- javascript_behavior_a_test.go:619 `TestJavaScriptBehaviorDurationsHeadlineRollingMedianAndCadenceTicks` — exact summary sentence.
- javascript_behavior_a_test.go:950 `TestJavaScriptBehaviorTimelineUserRequestGroupsUseOnlyListedMembers` — fixture `excludedReason: "paused"`; assertion at :1052-1053 checks only `reversed`.
- javascript_behavior_a_test.go:1394 `TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions` — "Paused and reversed spans are excluded".
- javascript_behavior_d_test.go:594 `TestJavaScriptBehaviorUserRequestProgressSumsAcceptedSpansAndLiveClaims` — fixture `implementationSpanReason: "paused"`; :649 comment.
- javascript_behavior_b_test.go:1727 (REQ-632 badge-absence row, wording only "assumed-pause badge"); no "paused" value assert.
- timeline_test.go:499-541 `TestTimelineProjectionChainsTheQueueSerially` — comments; fixture spans decide.
- durations_test.go:84, :116-117, :160, :222-233, :475-517 (ceiling straddle derived from `analysisOutlierCeiling`), :576, :744-788 (card/chart agreement, requires witnessing "paused").
- generate_test.go:3575 `implementationSpanOverCeilingFixtureSpan = 18h`, :3605-3607 guard against `analysisOutlierCeiling`, :3702 REQ-902 "paused", :3709 REQ-909 (4h21m, stamps 10:05/10:35/11:05/12:12/13:05/13:55/14:26 → largest gap 67 min) "paused" → must become "" — this fixture is already the REQ's RED case. REQ-901 (160-min span, largest gap 105) stays kept. generate_test.go:2456-2461 "Five paused spans" fixture.
- Latest durations (do-work/test-durations.tsv): javascript_behavior_a 0.90 s, _b 0.86 s, _d 0.49 s (2026-09-08, probably probes off), timeline_test 0.00, durations_test 0.04, generate_test 7.09, verify_test 4.40, state_apply_test 6.15, state_plan_test 0.59, finalization_recovery_test 18.49, finalization_req499_test 13.65.

## F6. Test commands (measured 2026-10-05)

- Board, targeted with JS probes: `cd skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run 'Duration|Timeline|ImplementationSpan|UserRequestProgress|DoneCard|Calibration' .` — ok, 22.3 s wall. Whole package ≈ 58 s (REQ-632 exploration).
- Core requeststate: `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/requeststate/` — ok, 5.3 s.
- Finalization (third row builder): `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/finalization/` — recovery 18.5 s + req499 13.7 s per the durations file; run it, it is not in the REQ's plan.

## F7. Where the backfill fits

- `skills/do-work/scripts/` holds shipped shell; the timestamp tools there (repair-req-timestamps.sh) are thin launchers into a do-work-cli Go subcommand. The do-work-cli prime says Go is canonical for deterministic operations. So a shipped backfill means a do-work-cli subcommand plus a launcher — heavy for a one-off.
- `_dev/` has no scripts directory; the only tracked Python there is `_dev/tests/heavy-runtime-fingerprint.py`; run helpers in do-work/runs/*/helpers/ are Python. No Go programs live in `_dev/`.
- Shell is a poor fit: RFC3339-with-offset parsing differs between BSD and GNU `date`.
- Recommendation: Python in `_dev/` (maintainer-only, idempotent: skip when the header already ends with `max_stamp_gap_minutes`), reusing the stamp list and order from F1. Consequence: consumer repos keep the old header, so the appender needs a defined behaviour for an old header (O1 append the column anyway; O2 rewrite the header line in place and update `calibrationAppendProves`; O3 skip the column when the header is old). That choice belongs in Decisions.

## Release

VERSION and skills/do-work/VERSION are 0.305.66; REQ-632 is merged but not yet finalized, so the preimage depends on its finalization landing first.
