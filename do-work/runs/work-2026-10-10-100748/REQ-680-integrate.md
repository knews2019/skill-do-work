# Integrator brief: REQ-680 (archive fetch test keeps an absent target absent)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-680.

- REQ id: REQ-680. Working REQ path: `do-work/working/REQ-680-archive-fetch-absent-target-test.md`. UR: UR-152 (`do-work/user-requests/UR-152/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-100748/` (run id `work-2026-10-10-100748`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-100748/REQ-680-handback.md`. Builder summary: one commit `9b8dee8b` on base `b629e5cd`; `archive_fetch_test.go` only: `TestTotalFailurePreservesTheTargetAndLeavesNoScratch` becomes a table with `pre-existing target` and `absent target` rows (63+/42-, mostly re-indentation). archivefetch tests 0 (3.9 s), `absent_target` PASS, probe 0, gofmt/vet clean. Production code unchanged.
- Coordinator ruling on D-03 (mutation proof): the brief asked for a mutation that fails only the new case; none exists, because removing private staging also overwrites a pre-existing target and breaks about 15 other tests sharing the function. Record in `## Testing`: "mutation (stage straight to the target, no cleanup) fails the new case on its own assertion line; the pre-existing row and other staging tests fail too, as expected for that mutation". The new row still pins a distinct failure: a failed fetch must never create a target that did not exist.
- Discovered task from the hand-back (`TestMissingArchiveTargetParentReportsUnattemptedRoutesInTextAndJSON` shows SKIP in the package run): `→ report only`; note why it skips if you see it in one read, no investigation.
- Operative name (branch = worktree basename): `worktree-agent-REQ-680-archive-fetch-absent-target-test`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-680-archive-fetch-absent-target-test`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T10:26:48Z`.
- Route A, tdd: false, impact-negligible, effort-mechanical, domain testing. Test-gate probe: `do-work/runs/work-2026-10-10-100748/REQ-680-probe.sh`.
- Integration order (coordinator ruling): hand-back order. Order: REQ-681 (released 0.305.94), then you, then REQ-683, REQ-684, REQ-679, REQ-685, REQ-682, and REQ-686 last. Re-read VERSION before the payloads (expect 0.305.95).
- Not the wave's last successful integration. Do not pass the wave-end fact. Restatement sweep: none expected (test-only).
- Manifest: the coordinator may have left dispatch stamps for other REQs pending; if `manifest.md` is clean, leave it.
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, cleanup); the coordinator does after you report.
- Lessons: none proposed by the builder.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-680-archive-fetch-absent-target-test.md`.
- Changelog title direction: plain words, what shipped, for example "Archive Fetch Test Proves a Failed Fetch Never Creates a Target That Did Not Exist"; the builder's proposed entry is in the hand-back.
