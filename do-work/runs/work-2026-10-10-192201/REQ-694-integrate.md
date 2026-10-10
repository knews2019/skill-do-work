# Integrator brief: REQ-694 (append-section hidden-section guard and frontmatter set owned fields)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-694.

- REQ id: REQ-694. Working REQ path: `do-work/working/REQ-694-append-section-hidden-sections-and-set-status-guard.md`. UR: UR-154 (`do-work/user-requests/UR-154/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-192201/` (run id `work-2026-10-10-192201`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-192201/REQ-694-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-694-append-section-and-set-guards`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-694-append-section-and-set-guards`. Base `fd9a2378`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T19:26:34Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-694-append-section-and-set-guards`, converted to UTC Z).
- Route A, tdd: true, impact-user-visible, effort-mechanical. Test-gate probe: `do-work/runs/work-2026-10-10-192201/REQ-694-probe.sh`.
- Integration order: second integrator. REQ-696 released 0.305.114 before you. Re-read VERSION before the payloads (expect 0.305.115). REQ-695 and then REQ-693 (last, closes UR-154) follow.
- Not the wave's last successful integration. Do not pass the wave-end fact.
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-694-append-section-hidden-sections-and-set-status-guard.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- UR-154 stays open: archive flat; link any lesson to the flat path.
- COORDINATOR RULING (do this before the first merge): the builder showed the decided count-only check can be fooled (a body that opens a fence, with a fenced example heading lower in the file: ## Review is hidden, the example heading becomes visible, the count still grows by one, the write succeeds). Replace the count check with: every section heading visible before the write is still visible after, and the only new visible heading is the appended section's name. Commit it on the builder branch `[REQ-694] require every prior section to stay visible after an append`, add the builder's repro as one more row in the existing N1 table test, and record it as a decision (it overrides the brief's decided design because the decided design does not fix the named failure). Then merge once.
- Seam: REQ-693 later edits internal/corehelpers/handoff.go in the same package; you add no package-level names.
