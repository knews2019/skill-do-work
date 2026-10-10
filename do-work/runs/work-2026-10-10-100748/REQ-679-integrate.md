# Integrator brief: REQ-679 (timeline browser probes that depend on live queue dates use the fixed fixture, and the lost Previous and Next assertions return)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-679.

- REQ id: REQ-679. Working REQ path: `do-work/working/REQ-679-timeline-probes-fixed-data.md`. UR: UR-152 (`do-work/user-requests/UR-152/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-100748/` (run id `work-2026-10-10-100748`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-100748/REQ-679-handback.md`. Builder summary: one commit `bab949ea` on base `b629e5cd`; `timeline_browser_probe_test.go` only (+86/-45). The shared range-end fixture gained two rows (a short completed REQ in the typed week, an open REQ since 2026-07-28); the prose-window probe builds from it; Previous/Next presses and the outright forward refusal were added. Timeline run with browser env: 38 pass, 0 fail, 22 `TestBrowserBehavior*` pass, 49 s; the 23 SKIP lines are `TestJavaScriptBehaviorTimeline*` heavy-only. Prose red/green: with the 54 REQs dated in that week deleted in a scratch worktree, the base probe failed and the new one passed. Step assertion red: a throwaway one-day Previous move failed the new assertion at `:2279`, reverted.
- Coordinator rulings: (1) D-02 accepted: none of the nine other probes moved, because none depends on dates; `:3564` deliberately uses the real board and moving it would need a large fixture (YAGNI). (2) Probe check: `REQ-679-probe.sh` runs the two named tests with `QUEUE_KANBAN_BROWSER_PROBES=on` and requires a `--- PASS:` line for each, so a skip fails it. Before the test gate, confirm the Previous/Next assertions (around `:2279`) sit inside one of those two named tests; if they sit in another test (for example the trailing-windows probe), add that test name to `probe_tests` in the probe (a run artifact) and say so in `## Testing`.
- Operative name (branch = worktree basename): `worktree-agent-REQ-679-timeline-probes-fixed-data`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-679-timeline-probes-fixed-data`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T10:26:48Z`.
- Route B, tdd: false, impact-negligible, effort-substantive, domain testing. Pre-flight green gate already recorded by the coordinator. Test-gate probe: `do-work/runs/work-2026-10-10-100748/REQ-679-probe.sh`.
- Integration order (coordinator ruling): hand-back order. Earlier integrators of this run bump VERSION first; re-read VERSION right before the payloads. REQ-686 stays last.
- Not the wave's last successful integration. Do not pass the wave-end fact. Restatement sweep: none expected (test-only).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, cleanup); the coordinator does after you report.
- Lessons: the builder proposed a bullet in the hand-back; judge it for truth and place it in `_dev/primes/lessons-kanban-board.md` beside its family (or the satellite it names if better), flat archive link.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-679-timeline-probes-fixed-data.md`.
- Changelog title direction: plain words, what shipped, for example "Timeline Browser Probes No Longer Depend on the Live Queue's Dates"; the builder's proposed entry is in the hand-back.
