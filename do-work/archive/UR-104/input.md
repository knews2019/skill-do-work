---
id: UR-104
title: 'Two-tier maintainer gate: cut to the 80/20 core, add --heavy, never stop the factory for it'
created_at: 2026-09-03T14:49:02Z
requests: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
word_count: 720
---

# Two-Tier Maintainer Gate: Cut to the 80/20 Core, Add --heavy, Never Stop the Factory for It

## Summary

The gate is 301 s after today's deletions and still runs on every step. The maintainer wants the corpus cut to the tests that pin real failures at a proportionate cost, a `--heavy` tier that only they run by hand, and a pipeline that never stops for the heavy tier: a REQ touching a heavy surface finishes on the fast tier and appends one line to a standing Needs-Input REQ. REQ-519 and REQ-522 are superseded and cancelled.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-537 | A1: tier `maintainer-verify.sh` into a default fast gate and `--heavy`, with `--heavy-surfaces` and a two-tier self-test |
| REQ-538 | A2: cut queue-kanban tests: delete the two re-exec meta-tests, JavaScript probes skip in the fast tier through one knob, audit and parallelize the rest |
| REQ-539 | A3: cut the contract file to under 1,500 lines, split survivors per owner, classify probes fast/heavy by measured time, build do-work-cli once for the heavy aggregate |
| REQ-540 | A4: cut do-work-cli tests: keep transaction and recovery suites, delete duplicate matrices, move binary-building and signal tests behind `-short` |
| REQ-541 | A5: the heavy ask in the pipeline: after a green fast gate, a REQ whose diff touches a heavy surface appends one line to the standing "Heavy gate runs requested" sweep and the loop continues |
| REQ-542 | A6: the SessionStart hook launches the gate runner in the background when it is not already running |

## Batch Constraints

- Order: REQ-537, then 538, 539, 540, 541, 542. Each is one integrating commit with its version bump and changelog entry; changelog titles say what shipped.
- Land in place, never through `do-work run`: this batch edits the gate and the pipeline would run the gate three times to ship it. Prove each commit with exactly one `bash _dev/tests/gate-runner.sh --once`, plus one `--heavy` run at the end of the batch.
- Delete before you add. Every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. Nothing deleted returns as a sentence pin.
- Never touch `do-work/working/REQ-533-*.md` or any other session's claimed file; stage explicit paths.
- Cancel REQ-519 and REQ-522 through the canonical `cancel` transaction naming the landing hash. Do not mint REQs for findings met on the way; they go in the commit body or `do-work/prose-backlog.md`.
- Done for the batch: default gate under 120 s wall and exit 0; `--heavy` exit 0; `--self-test` exit 0; `check-green-gate` matches at HEAD after the runner's last run; the board shows the heavy-ask REQ in Needs Input after REQ-541's scenario; before and after wall times per tier recorded in the changelog.

## Full Verbatim Input

> ```
> [maintainer, message 1]
> basically I think we should disable such heavy weight tests, first trim it down aggresively per the 80/20 principle
> 
> [maintainer, message 2]
> then we could have a heavy weight test, but that should be executed only with user permission
> 
> [maintainer, message 3]
> go, cut aggressively and add the --heavy tier. PS: the factory should not stop on the heavy tier, those are being asked by placing the REQ into th eneed input column: [board capture: NEEDS INPUT - BLOCKED column, empty]
> 
> [maintainer, message 4]
> give me a comprehensive prompt (capture-request style) to implement the request in a new session
> 
> [assistant draft, approved by the maintainer with "ok, capture the request"]
> do-work capture-request: Two-tier maintainer gate: cut the test corpus to the 80/20 core, add `--heavy` for the rest, and never let the heavy tier stop the factory
> 
> Context (already landed today, do not redo)
> - 0.266.8 (fde204ef): an already-green gate repair runs the gate once; `check-green-gate --at-revision` lets qualification and review verify the recorded green revision instead of relaunching.
> - 0.266.9 (5e0e166c): six dead `_dev/tests` scripts and their fixture map deleted; the harness self-tests no longer run inside the aggregate; one board test run carries `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1` instead of a second lane; `_dev/tests/probe-batch.sh` runs the twelve behavioral sub-suites concurrently under job control; the Kanban write-surface prose pin is gone; `_dev/tests/gate-runner.sh` runs the gate once per new HEAD in the background and records green via `record-green-gate` (`--once` for a single run).
> - REQ-520, REQ-521, REQ-532 were cancelled as landed in place; REQ-518 is archived completed. REQ-519 and REQ-522 are superseded by this request: fold their intent here and cancel them.
> - Full gate today: 301 s. lint+vet ~10 s; aggregate contract suite 112 s, bounded by one lane where update-script (~60 s) and install-suite probes build do-work-cli into the same path; queue-kanban 149 s; do-work-cli 35 s wall (114 s CPU across 25 packages).
> - queue-kanban has 457 tests, 141 s of test time. TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes alone is 41.8 s and TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes 16.1 s (both re-execute the whole test binary); the 55 TestJavaScriptBehavior* probes total 48 s; everything else ~35 s. do-work-cli's slowest: corehelpers.TestInventoryMatchesRetainedPorcelainXYMatrix 12.6 s, publication.TestDeferGateRollsBackUntrackedCreateAndFoldTopologies 8.1 s, finalization.TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce 6.7 s. contract-regressions.sh is 8,451 lines with 294 sentence-assertion sites, 27 python heredocs, 6 git-init fixtures; top targets are extracted prose blocks ($verify_revalidation_block 19, $crash_recovery_block 10, $review_archived_input_block 9) and work-reference.md/work.md.
> 
> Maintainer decisions (settled)
> - D1 The fast tier IS the canonical gate. `bash _dev/tests/maintainer-verify.sh` with no flags is what the pipeline, the gate runner, and every REQ run. Target under 2 minutes wall, uncached.
> - D2 `--heavy` is opt-in and run only by the maintainer by hand. Nothing under skills/ or the gate runner ever passes --heavy.
> - D3 Cut aggressively by 80/20. A test survives only if it pins a real failure at a cost proportionate to that failure. A sentence pin on prose is not a real failure. A test that re-runs something another test covers is deleted, not moved to heavy.
> - D4 The factory never stops for the heavy tier. When a REQ's diff touches a heavy-lane surface, the loop finishes the REQ on the fast tier and appends one line to a single standing pending-answers REQ "Heavy gate runs requested" (sweep, sweep_key heavy-gate-requested) in the Needs Input column. The maintainer runs --heavy when they choose and ticks the line. One file, appended to; never a new REQ per occurrence, never a stopped run.
> - D5 Mechanics in scripts and Go, judgment in prose. No new prose that walks a shell sequence. No new sentence pins anywhere.
> 
> Requests A1 tier the gate; A2 cut queue-kanban; A3 cut the contract file; A4 cut do-work-cli tests; A5 heavy ask in the pipeline; A6 runner start. Constraints: one integrating commit per request in order, land in place (never through do-work run), one gate-runner --once per commit plus one --heavy at the end, delete before you add and list every deleted test with the failure it pinned, never touch another session's claimed file, cancel REQ-519 and REQ-522 with the landing hash, no new REQs for findings. Verification: default gate under 120 s exit 0, --heavy exit 0, --self-test exit 0, check-green-gate matches at HEAD, the board shows the heavy-ask REQ in Needs Input after the A5 scenario, before/after wall times in the changelog.
> 
> [maintainer, message 5]
> ok, capture the request
> ```

---
*Captured: 2026-09-03T14:49:02Z*
