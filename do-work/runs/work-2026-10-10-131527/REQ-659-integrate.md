# Integrator brief: REQ-659 (frontmatter set and req append-section commands)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-659.

- REQ id: REQ-659. Working REQ path: `do-work/working/REQ-659-frontmatter-set-and-req-append-section.md`. UR: UR-145 (`do-work/user-requests/UR-145/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-659-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-659-frontmatter-set-append-section`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-659-frontmatter-set-append-section`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:24:46Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-659-frontmatter-set-append-section`, converted to UTC Z).
- Route B, tdd: true, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-659-probe.sh`.
- Integration order: second integrator. REQ-660 released 0.305.102 before you (archive do-work/archive/REQ-660-worktree-lifecycle-command.md). Re-read VERSION before the payloads (expect 0.305.103). REQ-658 integrates after you and closes UR-145.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your diff redefines how stamps and sections are written (work.md Living Logs, work-reference Stamps are append-only, the shared section order list in advance): tell your reviewer to run the Restatement Sweep over every other place that restates the stamp-writing or section-order rules.
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-659-frontmatter-set-and-req-append-section.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Builder note: one file outside the write boundary, `internal/corehelpers/commands_test.go` (handler count 21 to 22) was required by the new command (D-11); accept it in Qualification as a necessary scope addition.
- Seam: REQ-660 already landed `resultmodel/result_model.go` and `fan-out-reference.md` edits; resolve any conflict by keeping both.
- Lessons: REQ-660's integrator used the GitHub-link style already used in `lessons-do-work-cli.md`; follow the same style. REQ-658 re-points the UR-145 links when it closes the UR.
