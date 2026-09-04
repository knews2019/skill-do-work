# Restart Prompt — wave 2 of UR-106/UR-099 follow-on run (2026-09-04)

```
do-work run REQ-547 REQ-564 --fan-out 2

This command is sufficient; everything below it is context.

Pull request #180 is MERGED. All of wave 1 is on main: the gate retry,
the own-paths cleanliness change, and the per-REQ recovery change,
released as 0.277.0, 0.278.0 and 0.279.0. Do not reopen it and do not
push wave 2 onto that pull request - a merged PR is finished. Start the
branch fresh from main and open a new PR when wave 2 is ready.

Two REQs are claimed by this checkout (writer vm:/home/user/skill-do-work)
and both have COMPLETE builder branches on origin, still unmerged:
  worktree-agent-REQ-547-finalize-refuses-a-req-with-no-checkpoint-entry  (tip b3d25c8)
  worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours  (tip 19103f5)
Neither needs a builder re-dispatched. Both resume at the Step 6
hand-back merge. Merge them onto main, not onto the old branch.

If `recover` reports the claims as another writer's, take them over first:
  do-work-cli recover --take-over REQ-547
  do-work-cli recover --take-over REQ-564

Before wave 2, fix the two P1 findings on already-released code - they
are in the Reference section below and they are the first work.

Before running the repository gate, provision the toolchain this container did
not ship:
  export PATH="/root/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.1.linux-amd64/bin:/usr/local/bin:$PATH"
The gate needs Go 1.26.1, ShellCheck 0.11.0, and a `just` new enough for
[positional-arguments]. Without all three it cannot start, and a gate that
cannot start still reads as a red gate.

Each REQ's `## Handoff State` section carries the rest.
```

---

## Reference

### Wave 1 — done, do not revisit

| REQ | What it did | Released |
|---|---|---|
| REQ-559 | Rerun a red repository gate once before deferring or minting a repair REQ | 0.277.0 |
| REQ-560 | Hand-back and finalize check cleanliness only on the REQ's own paths | 0.278.0 |
| REQ-515 | Per-REQ recovery findings never stop the loop | 0.279.0 |
| REQ-561 | Not built — its code shipped in 0.273.0; cancelled against that release | — |

UR-106 and UR-099 both closed with all members archived.

### Wave 2 — in flight

**REQ-547** — stop finalize refusing a REQ that has no checkpoint entry.
- Branch `worktree-agent-REQ-547-finalize-refuses-a-req-with-no-checkpoint-entry` at `b3d25c8`, pushed, **unmerged**.
- Worktree `/home/user/skill-do-work-worktrees/worktree-agent-REQ-547-...` — **ACTIVE**, clean.
- Complete hand-back at `do-work/runs/work-2026-09-04-200249/REQ-547-handback.md`.
- Merge range: not yet captured. `<pre>` is whatever `git rev-parse --short HEAD` reads immediately before the first merge.
- Remaining: merge → qualify → gate → review → lessons → release → finalize.

**REQ-564** — reuse matching per-lane verification evidence for four hours.
- Branch `worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours` at `27b74c8`, pushed, **unmerged**.
- Worktree `/home/user/skill-do-work-worktrees/worktree-agent-REQ-564-...` — **ACTIVE**, clean.
- Complete hand-back at `do-work/runs/work-2026-09-04-200249/REQ-564-handback.md`. The builder finished after appearing interrupted; its `4526ab9` and `19103f5` sit on top of `27b74c8`, an orchestrator-authored WIP snapshot taken to survive the container. **Read the tip, not the snapshot.**
- Builder-reported verification: vet, gofmt, full module tests, Windows cross-compile and contract-regressions all clean, with twelve isolating reverts each proved red. No real heavy lane was run, so the merged-tree gate is still owed.

### Parallelism

`--fan-out 2` is safe. The two REQs are disjoint: REQ-547 is in `internal/requeststate/` and `internal/finalization/`; REQ-564 is in `internal/heavyverification/` and `internal/resultmodel/`. The only shared file is `skills/do-work/actions/work.md`, and they touch different sections. No dependency gate between them, so neither is on the other's critical path. Integration is serial regardless — merge, gate, review and release run one REQ at a time, which is where the wall clock goes.

Both are complete and unblocked; start with either. REQ-547 is the smaller merge.

### Heads-up list — things that will bite in the first ten minutes

- **The container ships the wrong toolchain.** Go 1.24.7, no ShellCheck, `just` 1.21.0. The gate needs Go 1.26.1, ShellCheck 0.11.0, and a `just` that understands `[positional-arguments]`. Provision all three before the gate, or it fails for reasons that have nothing to do with the code. — next session
- **One test file straddles its budget.** `_dev/tests/session-start-hook-behavior.sh` measured 27, 30, 32, 39 and 44 seconds across five runs today against a strict 30-second limit, with nobody touching it. It caused two full gate reruns. It is a flake, not a regression. Nothing in the queue covers it; worth capturing. — maintainer
- **Version numbers collided once already.** This branch and main both released from 0.275.3. Main took 0.275.4 and 0.276.0; this branch renumbered to 0.277.0/0.278.0/0.279.0. Before cutting the next release, re-read `CHANGELOG.md`'s first entry rather than assuming. — next session
- **REQ-508 is claimed by another checkout** (`t2s-Virtual-Machine.local:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`). Its checkpoint entry is FOREIGN — leave it byte-identical. — next session
- **Pull request #180 is merged**, so wave 1 is on main and that PR is closed history. Its review threads carry the two P1 findings below, already answered but not fixed. The same reviewer already found two genuine P1 bugs in REQ-515 that two full independent review passes had missed, so treat its findings as bug reports: verify each against the code, fix what is real, reply on the thread. Wave 2 needs a new pull request off main. — next session
- **Estimates run low on this queue.** Calibration: REQ-559 25→33, REQ-560 20→51, REQ-515 30→105, REQ-568 50→111. The gap is not building; it is merge, gate, review and remediation. Do not plan wave 2 against the P50 figures. — maintainer

### Open P1 findings on shipped code — do these first

Two P1 findings landed on knews2019/skill-do-work#180 after the handoff was written. Both are on REQ-515's code, **already released as 0.279.0**. I verified both against the source and replied on both threads; neither is fixed. They are the first work of the next session, ahead of wave 2.

**P1-A — recovery accepts a committed-risk transaction.** `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`, `requestScopedRefusal`. `gittransaction` defines `FailureCommittedRisk = "committed_state_risk"` and `exact_commit.go` returns it from seven sites *after HEAD has advanced*. `requestScopedRefusal` tests only `OutcomeFailure` and `RollbackIncomplete`; a rollback that restored files but not HEAD reports `succeeded`, so no guard fires, the record is accepted as request-scoped, and recovery drains later journals with an unknown commit in history. **This is a widening of what recovery accepts, which REQ-515's own Constraints forbid.** Fix shape: preserve committed risk through `advanceJournal`'s `FINALIZATION-PRIMARY-COMMIT` wrap instead of flattening it to a rolled-back outcome, then have `requestScopedRefusal` treat it as shared state. Pin it with a test that fails without the guard.

**P1-B — a set-aside journal's shared lifecycle images go stale.** Same file, `consumeRecoveryRecord`. When two prepared journals share `do-work/CHECKPOINT.md`, setting the first aside lets the second remove its own claim and commit; the first journal's preimage and postimage then match neither the original nor the current checkpoint, so its next recovery reports `FINALIZATION-LIFECYCLE-CONFLICT` with an incomplete rollback and the queue stops permanently. Fix shape: rebase the set-aside journal's shared lifecycle images as siblings settle, or make its later recovery tolerate sibling-owned checkpoint transitions.

**P1-B was already known and was misjudged.** This REQ's own re-review found it and scored it report-only, reasoning that it fails safe and still beats the old behaviour where the first refused record parked the whole queue. Both points are true and neither makes it acceptable — the promise this REQ shipped is "every other REQ still runs", and that promise holds only within one queue boundary. Treat it as P1.

### Queue-speed recommendation

`do-work run REQ-562` after wave 2. It records per-run lifecycle spans and reports the critical path, which is the only thing that turns the calibration gap above into something actionable. REQ-564 helps but is narrower than its title: this session's nine gate runs were fast-tier, not heavy lanes.

### Worktree verdicts

- `/home/user/skill-do-work` — integration checkout, branch `claude/req-559-560-515-547-564-g82uxl`, clean.
- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-547-finalize-refuses-a-req-with-no-checkpoint-entry` — **ACTIVE**, clean, unmerged.
- `/home/user/skill-do-work-worktrees/worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours` — **ACTIVE**, clean, unmerged.

No worktree is REMOVABLE. None was removed and no foreign claim was touched.
