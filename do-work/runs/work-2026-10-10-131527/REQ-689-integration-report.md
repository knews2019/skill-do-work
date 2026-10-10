# REQ-689 integration report

- Release 0.305.113, "do-work run --coordinate Keeps the Main Session a Coordinator, With Run Rules, a Preflight and a Stall Check". Archive `do-work/archive/UR-153/REQ-689-run-coordinate-mode.md`. Finalization `7135c0fe` reached cleanup_complete on the first try. UR-153 is closed: REQ-687, 688, 690, 691, 692 and `input.md` are now in `archive/UR-153/`.
- Range: pre `cc049a01`. Merges `f8f32682`, `9cd669eb` (review fixes), `9f815f15` (re-check fixes; this is `commit:`). The merge was clean. Seams fixed in the merge: the stall loop reads `do-work status --watch` plus the progress-log tails, and uses only tails and commit times when status cannot report. `status-guide.md:34` example corrected.
- Review (wave end, sweep of all 11 siblings): 86%, Approve. I fixed F1-F14. F5 was adjusted so the handoff message still ends with its two lines. The fixes cover: the builder progress-log write rule, the stall loop sparing a running full gate, coordinated handoff claim lines, no rm line for full-gate.lock in status, run flags in the argument hint, Step 9 cleanup pointers, and the ai-report `--kind` input wording. Delta re-check: 95%, Approve. N1 and N2 fixed.
- Gate (reuse off), exit 0 at all three merges: 115 s at load 2.27, 117 s at load 4.97, 117 s at load 4.65. Advance test gate at `9f815f15`: green, probe exit 0.
- Heavy lane staged-skills: executed, exit 0, 32 s.
- Timing: 6 events, folded. No builder-work event (the hand-back landed long before).
- Lessons: one restated-mechanism-unchecked bullet. REQ-688, 690, 691 and 692 bullets re-pointed to archive/UR-153/. 3 index rows refreshed. contract-regressions.sh passes after finalization.
- Cleanup: worktree removed, branch deleted with -d, pruned. Optional batch packaging check: not run.
- Left dirty: only this report (untracked).
- Report only: N3 (a remote builder has no progress log); the `run_status.go:363` comment; still-stale items from earlier siblings (REQ-688 F2/F3, REQ-690 F8/F9, REQ-692 M2, REQ-660 F3); `run-simple-reqs.md` forwards only `--fan-out`.
