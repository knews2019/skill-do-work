```
hand REQ-704 to an integrator (it enters with advance REQ-704)
do-work run --coordinate --fan-out 1
These commands are sufficient; everything below them is context.
```

---

## Reference

**Run:** `do-work run UR-160 --coordinate`, run dir `do-work/runs/work-2026-10-20-200259/` (manifest, `PREDISPATCH-GUIDE.md`, `INTEGRATOR-GUIDE.md`, per-REQ briefs, probes, hand-backs, progress logs). UR-160 = directed second opinions (presets, free-text direction, linked iteration tree), decisions D1-D10 in `do-work/user-requests/UR-160/input.md`; D2 (shared project preset folder) is superseded by REQ-705's answer.

**Done this session:**
- REQ-703 (second opinions accept a named preset and a free-text editor direction): released 0.306.6, `a6c00290`, archived flat at `do-work/archive/REQ-703-second-opinion-preset-and-editor-direction.md`.
- REQ-705 (user presets live only inside the report bundle; the user's answer to the builder-decided confirm question): released 0.306.7, `ae016e17`, archived flat at `do-work/archive/REQ-705-confirm-uncommitted-preset-refusal.md`.

**In flight: REQ-704 (second opinions build a linked iteration tree on any base page)**, claimed in `do-work/working/REQ-704-second-opinion-linked-iteration-tree.md`.
- Merged: yes. Builder tip `50555523`, merge `92443689`, merge range `87614289..92443689` (keep `87614289` as `<pre>` on any re-merge).
- Written: Triage, Plan (D-01..D-13), Exploration, Scope, Pre-Flight, Implementation Summary, Decisions, Discovered Tasks, Qualification; stamps `dispatch_at`, `builder_handback_at`, `integration_at`.
- `recover`: `RECOVERY-TAKEOVER-AVAILABLE` (claim preserved). `advance REQ-704` reports phase `test-gate`.
- Remains, in order: repository gate at `92443689` (the first run was stopped mid-way by the coordinator at the user's handoff request; result unknown), the advance test gate with `REQ-704-probe.sh`, `## Testing`, review with the wave-end Restatement Sweep over REQ-703 and REQ-705, Lessons/Orientation, heavy plan and drain, link sweep, UR-160 closure finalization (next patch after 0.306.7), worktree cleanup. Then the coordinator runs `advance --checkpoint`.
- Uncommitted files: none (survey: both checkouts clean). The integrator brief is `do-work/runs/work-2026-10-20-200259/REQ-704-integrate.md`, with the shared guide beside it.

**Worktrees:**
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` (main, `9c297ea4` + this handoff commit): integration checkout.
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-704-second-opinion-linked-iteration-tree` at `50555523`: **ACTIVE**. Merged and clean, but REQ-704 is still in `do-work/working/`; review fixes go on this branch. Remove it only through the integrator's cleanup after finalization.

**Parallelism:** one member left, so `--fan-out 1`. The queue is otherwise empty; no gate edits were needed.

**Heads-up for the first ten minutes:**
- REQ-704 integrator: no `full-gate.lock` is left (the coordinator removed the stopped integrator's lock). Rerun the full gate from the start.
- REQ-704 integrator: closing UR-160 moves REQ-703 and REQ-705 into `do-work/archive/UR-160/`; re-point the REQ-703 lesson link in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and refresh `do-work/lessons-index.md` (`git grep -n "archive/REQ-70[345]" -- ':!do-work/'`). The integrate brief lists the commit paths.
- REQ-704 integrator: `report_second_opinion_test.go` ran 27.2 s against the 30 s per-file budget in REQ-705's gate; REQ-704's new tests are in `report_second_opinion_tree_test.go` (2.5 s). A budget red under load is rerun once.
- Coordinator: the old 2026-10-10 handoff in this file (UR-144) is replaced; its REQs were finished long ago.
