```
do-work run UR-115 UR-117
This command is sufficient; everything below it is context.

You are resuming a do-work run in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 handed off at 2026-09-05T00:11:20Z. Run one REQ at a time (serial default). Order is already encoded in the queue: REQ-574 (repository gate repair, claimed by this checkout's writer label, resumes at the planning phase; a plan draft sits at do-work/runs/work-2026-09-04-232225/REQ-574-plan.md, validate it before adopting) → REQ-572 (gate_deferred, depends_on REQ-574, saved implementation range 7ad53bff..fbdcd35e already merged on main; run the saved-range resume proof, then only qualification, focused tests, the gate and review) → REQ-573 (depends_on REQ-572) → REQ-578 (no dependency; UR-117).

The repository gate (bash _dev/tests/maintainer-verify.sh) is red on time only: four do-work-cli test files exceed the 30-second per-file budget (internal/corehelpers/inventory_test.go, internal/publication/defer_gate_test.go, internal/finalization/finalization_recovery_test.go, internal/finalization/finalization_req499_test.go); every test passes. Three tests own most of it: TestInventoryMatchesRetainedPorcelainXYMatrix 53s, TestDeferGateRollsBackUntrackedCreateAndFoldTopologies 27s, TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce 22s (per-test log: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/83001520-580f-484b-a90e-0f2e11fcabac/scratchpad/slow-tests-574.log, may be gone). REQ-574 must give those files real headroom (target under 15s per file with one gate running) by speeding the fixtures, splitting by genuine test family, or moving true integration scenarios to the heavy lane; never by raising the budget. Run every gate and heavy lane from a detached worktree (git worktree add --detach), because other sessions dirty the main tree. Expect recover to refuse on foreign live paths (other sessions' do-work/runs/* and do-work/working/REQ-577/REQ-506 files): judge them foreign, leave them byte-identical, and continue.
```

---

## Reference

### In-flight REQ (ACTIVE, this checkout's writer label)

- **REQ-574** — Repository gate repair: bring do-work-cli test files under the 30s per-file budget. Claimed 2026-09-04T23:59:43Z at 8269a0bb, Route C, estimate P50 40 min, `## Triage` written, no `## Plan` yet. A plan draft exists at do-work/runs/work-2026-09-04-232225/REQ-574-plan.md; it was NOT consumed into the REQ (read it as input, validate it, then write ## Plan and planning_at yourself). Not merged, no builder branch exists, nothing uncommitted of its own (owner evidence committed at 05ab21d5). Evidence: `do-work/working/REQ-574-bring-do-work-cli-test-files-under-the-30s-budget.md`, run dir `do-work/runs/work-2026-09-04-232225/`.

### Deferred and queued (this UR set)

- **REQ-572** — merged on main at fbdcd35e (range `7ad53bff..fbdcd35e`, builder head 2d8beb40), qualified, deferred by defer-gate at 35af5c23 behind REQ-574. Frontmatter carries `gate_deferred`, `deferred_implementation_base/merge`. Its `## Qualification` and `## Implementation Summary` sections are present; on resume the classifier will treat them as done, but the reference requires rerunning qualification, focused tests, the gate and review before completion. Hand-back with builder decisions D-01 to D-07: `do-work/runs/work-2026-09-04-232225/REQ-572-handback.md` (untracked scratch, keep).
- **REQ-573** — pending, depends_on REQ-572. Open the detail drawer from an Activity row and highlight every row of the same REQ. A read-only prep another session wrote (`do-work/runs/work-2026-09-05-005615/activity-prep.md`, committed) proposes its scope including `web/board-detail.js`.
- **REQ-578** — pending, no dependency, UR-117: hide the verify-findings strip while the Activity view is active (view switch in `board-controls.js`, not the renderer).

### Worktrees (survey at handoff)

- `skill-do-work2-worktrees/worktree-agent-REQ-570-…` — FOREIGN (other session; REQ-570 released as 0.287.0, its cleanup is theirs).
- `skill-do-work2-worktrees/gate-REQ-570` — FOREIGN (other session's detached gate checkout).
- `.git/work-run-20260905/worktree-agent-REQ-506-focused-evidence` (branch codex/…) — FOREIGN, clean.
- `.git/work-run-20260905/worktree-agent-REQ-577-launcher-fixture` (branch codex/…) — FOREIGN, clean.
- REQ-572's builder worktree and branch were removed after the merge (no leftover). No REMOVABLE worktree belongs to this session.

### Parallelism

- Critical path: REQ-574 → REQ-572 → REQ-573. REQ-578 is independent but shares `javascript_behavior_c_test.go` with REQ-573; run it serially after REQ-573 to avoid an append conflict in the merge (no gate added, per the write_set rule).
- Serial only. The gate's per-file budget is load-sensitive: with two or three concurrent gate processes the four files measured 31 to 60s; with one they sat at 17 to 28s (history in `do-work/test-durations.tsv`, last column is concurrent gate count).

### Heads-up (first ten minutes)

- Two other sessions share this checkout: REQ-577 (ShellCheck launcher-fixture repair, claimed 23:58Z) and REQ-506 (deferred behind REQ-577) are theirs; their run dirs `do-work/runs/work-2026-09-05-005615/` and `work-2026-09-05-020017/` and `do-work/CHECKPOINT.md` header edits are foreign. Owner: those sessions.
- `do-work/CHECKPOINT.md` lists REQ-574 and REQ-577 under In Progress with the same writer label (host:path); recover cannot tell sessions apart. Only REQ-574 is yours. Owner: next session.
- `stash@{0}: On main: do-work recovery: preserve interrupted REQ-539 orchestration metadata` is parked and easy to pop by accident from any worktree. Owner: the maintainer (apply or drop).
- The pre-build green-gate record for REQ-572 was never recorded; the test-gate phase needs a direct gate run with `--gate-exit-status`. Owner: next session.
- REQ-574's plan draft was written by a Plan agent whose validation never ran; run the Route C plan validation before adopting it. Owner: next session.
- `do-work/runs/work-2026-09-04-232225/REQ-572-probe.sh` is the focused-test probe for REQ-572's test gate (untracked scratch). Owner: next session.
