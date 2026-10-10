# Integrator brief: REQ-688 (capture-files --example and the capture-reference fence fix)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-688.

- REQ id: REQ-688. Working REQ path: `do-work/working/REQ-688-capture-files-example-and-fence-fix.md`. UR: UR-153 (`do-work/user-requests/UR-153/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-688-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-688-capture-files-example-fence-fix`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-688-capture-files-example-fence-fix`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:23:59Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-688-capture-files-example-fence-fix`, converted to UTC Z).
- Route B, tdd: true, impact-user-visible, effort-mechanical. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-688-probe.sh`.
- Integration order: seventh integrator, first of UR-153. Re-read VERSION before the payloads (expect 0.305.108).
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (the capture fence rule: check clarify.md:106, which the builder reported states the rule without the three-backtick minimum).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-688-capture-files-example-and-fence-fix.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Released before the UR-153 members: REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed, archive/UR-145/), REQ-655 0.305.105, REQ-657 0.305.106, REQ-656 0.305.107 (UR-144 closed, archive/UR-144/). UR-153 order: REQ-688, REQ-692, REQ-687, REQ-690, REQ-691, then REQ-689 last (closes UR-153, carries the wave-end sweep). UR-153 stays open until REQ-689: archive flat and link any lesson bullet to the flat path; REQ-689's integrator re-points them.
- New tools on main you may use: `worktree merge` (REQ-660), `frontmatter set` and `req append-section` (REQ-659), `finalize --auto-manifest --emit` (REQ-658). The helper scripts also still work.
- Disk: the coordinator cleared the Go build cache earlier; check `df -h /System/Volumes/Data` before the gate and stop with a report if under 5 GB free.
- Stale path in your working REQ: `do-work/working/REQ-688-capture-files-example-and-fence-fix.md` (about line 65) cites REQ-661's flat archive path; REQ-661 now lives at `do-work/archive/UR-145/REQ-661-capture-files-init-manifest-skeleton.md`. Fix that path in your REQ file before qualify.
- Seam: REQ-692 later adds one clause near capture.md:119; you edit Step 5 only.
