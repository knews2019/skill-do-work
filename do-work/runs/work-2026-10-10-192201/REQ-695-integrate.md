# Integrator brief: REQ-695 (run-status missing --run folder and integrator-at-work C3)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-695.

- REQ id: REQ-695. Working REQ path: `do-work/working/REQ-695-run-status-missing-run-and-integrating-c3.md`. UR: UR-154 (`do-work/user-requests/UR-154/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-192201/` (run id `work-2026-10-10-192201`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-192201/REQ-695-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-695-run-status-missing-run-integrator`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-695-run-status-missing-run-integrator`. Base `fd9a2378`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T19:27:42Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-695-run-status-missing-run-integrator`, converted to UTC Z).
- Route B, tdd: true, impact-user-visible, effort-mechanical. Test-gate probe: `do-work/runs/work-2026-10-10-192201/REQ-695-probe.sh`.
- Integration order: third integrator. REQ-696 0.305.114 and REQ-694 0.305.115 released before you. Re-read VERSION before the payloads (expect 0.305.116). REQ-693 (last, closes UR-154) follows.
- Not the wave's last successful integration. Do not pass the wave-end fact.
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-695-run-status-missing-run-and-integrating-c3.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- UR-154 stays open: archive flat; link any lesson to the flat path.
- Seam: REQ-693 later edits WorktreeStatusRow in resultmodel/result_model.go right above RunStatusResult; you should not need that file. status.md and status-guide.md C3 cells are yours alone in this run.
- Builder's report-only risk (the integrator-phase list copies the board's phase names with no test keeping them in step): have the reviewer judge whether one assertion against the board's phase constant is cheap; if not, report only.
