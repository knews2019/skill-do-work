```
do-work run --fan-out 2
This command is sufficient; everything below it is context.
```

---

## Reference

- **Integration state:** `main` at `156f7c6f` before this handoff commit. This session shipped REQ-408 (`ac2e3acd`, metadata `e6488553`), REQ-409 (`a57bf51e`, metadata `208c2178`), and REQ-410 (`210d1459`, metadata `156f7c6f`). Version 0.251.0 and the canonical maintainer gate are green.
- **In-flight REQs:** none. REQ-411 was stopped before planning or implementation, released from `do-work/working/` to `status: pending`, and has no branch, worktree, merge commit, or merge range. Its only retained work is the frozen 100-minute low-confidence estimate; fresh planner/explorer output is required.
- **Uncommitted files at survey time:** the REQ-411 claim/release bookkeeping, `do-work/CHECKPOINT.md`, the confirmed REQ-427 answer, and reservation markers REQ-428 through REQ-435. The handoff commit includes all of them plus this prompt; no implementation diff is being handed over.
- **Worktree verdict:** ACTIVE — `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` on `main` is the sole checkout and is the integration tree the new session should use. There are no builder worktrees and no merged or unmerged `worktree-agent-*` branches. Nothing is removable.
- **Parallelism:** fan-out 2 is intentional. REQ-426 and REQ-427 are disjoint installer follow-ups; REQ-428 and REQ-429 are disjoint dependency-graph/request-model fixes. REQ-430 through REQ-433 must remain serial because they all change cleanup planning/application; the queue encodes `REQ-431 depends_on REQ-430`, `REQ-432 depends_on REQ-431`, and `REQ-433 depends_on REQ-432`. REQ-434 can pair with cleanup work. REQ-435 must finish before REQ-411 because both own doctor/forensics-to-work-action contracts; REQ-411's dependency field encodes that gate. REQ-412 waits for both REQ-411 and REQ-433 so request-state/archive transactions build on corrected selection and cleanup semantics.
- **Critical path:** close REQ-428, REQ-429, and REQ-435, then REQ-411 → REQ-412 → REQ-413 → REQ-414 → REQ-415 → REQ-416 → REQ-417 → REQ-418 → REQ-419 → REQ-420. The cleanup review-fix chain REQ-430 → REQ-431 → REQ-432 → REQ-433 joins at REQ-412.
- **First-ten-minutes heads-up:** REQ-410 is intentionally `completed-with-issues`; REQ-434 owns the unsupported timestamp-ordering-anchor defect and REQ-435 owns the incomplete doctor/forensics delegation contract. REQ-427 is no longer awaiting clarification: the confirmed answer is to lower only the core installer/updater floor to Go 1.23.0, leaving the optional board module at Go 1.26.1. The user-owned edit is now committed queue state, not a dirty-tree surprise.
