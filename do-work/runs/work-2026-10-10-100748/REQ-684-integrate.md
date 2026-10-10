# Integrator brief: REQ-684 (qualify tells a missing Implementation Summary section apart from a section that names no files)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-684.

- REQ id: REQ-684. Working REQ path: `do-work/working/REQ-684-qualify-summary-messages.md`. UR: UR-152 (`do-work/user-requests/UR-152/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-100748/` (run id `work-2026-10-10-100748`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-100748/REQ-684-handback.md`. Builder summary: one commit `2a1a1789` on base `b629e5cd`; `checks.go` `handleQualify` now says "Implementation Summary section not found" or "Implementation Summary lists no backticked file paths"; finding code `QUALIFY-SUMMARY-MISSING`, severity and fixability unchanged; two tests in `checks_test.go`, both red at base and green after. corehelpers tests 0 (11.6 s), probe 0, gofmt/vet clean. 35+/1-.
- Operative name (branch = worktree basename): `worktree-agent-REQ-684-qualify-summary-messages`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-684-qualify-summary-messages`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T10:28:41Z`.
- Route A, tdd: true, impact-user-visible, effort-mechanical, domain backend. Test-gate probe: `do-work/runs/work-2026-10-10-100748/REQ-684-probe.sh`.
- Guide section 12 "The merged build runs your own tools" applies: your own qualify call runs the new messages.
- Integration order (coordinator ruling): hand-back order. Earlier integrators of this run bump VERSION first; re-read VERSION right before the payloads. REQ-686 stays last.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer records its own **Restatement sweep:** line: grep actions, docs, crew files and lessons for the old "missing or empty" wording or any description of when `QUALIFY-SUMMARY-MISSING` fires, and report what still quotes the old text.
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, cleanup); the coordinator does after you report.
- Lessons: none proposed by the builder.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-684-qualify-summary-messages.md`.
- Changelog title direction: plain words, what shipped, for example "Qualify Says Whether the Implementation Summary Is Missing or Names No Files"; the builder's proposed entry is in the hand-back.
