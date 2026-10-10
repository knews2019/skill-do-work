```
do-work run, coordinator shape: this session coordinates only; background agents run the queue.
Builders run in parallel git worktrees, one per REQ, up to capacity. Integrators run one at a time
(merge, gate, review, release). Follow skills/do-work/actions/fan-out-reference.md, section
"Delegated integration — the coordinator shape".

1. Five REQs are already claimed in do-work/working/ by the previous session. They have no sections,
   branches or worktrees yet. Continue each claim first; each call is read-only and names the next phase:
   advance REQ-655
   advance REQ-657
   advance REQ-658
   advance REQ-659
   advance REQ-660
   Build all five in parallel worktrees. Integrate them one at a time, in this order:
   REQ-660, REQ-659, REQ-658, REQ-655, REQ-657.
2. Before step 3, check do-work/queue/. REQ-654, 661, 667, 668, 669, 670 and 674 must no longer be
   `pending`: session skill-do-work2-b1 is cancelling and recapturing them as fewer, smaller REQs.
   If any of them is still `pending`, do not run step 3. Stop and tell the user.
3. do-work run --fan-out 4

These commands are sufficient; everything below them is context.
```

---

## Reference

Written 2026-10-10T12:57Z by the session that claimed the UR-144 to UR-149 batch, then handed off before any build started.

**What happened**

- At 12:52Z this session ran `recover` (clean) and `advance --fan-out 12`. That claimed twelve REQs, in commits 9d844aaf through 6a28ed59.
- The maintainer then chose the aggressive queue simplification in session skill-do-work2-b1. That session asked for seven claims back. All seven went back to the queue with `recover --take-over`, one commit each: REQ-654 2b109946, REQ-661 ba3975ef, REQ-667 87565a26, REQ-668 bcdb86c0, REQ-669 49a93977, REQ-670 49b76c10, REQ-674 61e02428.
- No workflow, builder, worktree, branch, run directory or REQ section was created. The empty run directory `do-work/runs/work-2026-10-10-125225/` was deleted before the requeue.

**Claims left in `do-work/working/`** (claimed by this checkout, no sections, nothing merged or uncommitted)

- REQ-655: `ai-report index` and `find` build a catalog of report bundles. Go under `internal/toolboxcommands/`, plus toolbox ai-report docs.
- REQ-657: `ai-report judge` runs the render check as a bundled command. Toolbox `scripts/judge.mjs`, plus ai-report docs.
- REQ-658: `finalize --auto-manifest` builds the mechanical manifest fields. Go under `internal/finalization/`, plus work.md and work-reference.md.
- REQ-659: `frontmatter set` and `req append-section`. Go `corehelpers/commands.go`, the board's `frontmatter_cli.go`, plus work.md and work-reference.md.
- REQ-660: `worktree new/status/merge/cleanup`. A new Go command, plus fan-out-reference.md.

**Parallelism**

- Safe to build all five at once: no `depends_on` edges between them.
- Shared prose: 655 and 657 both edit `skills/do-work-toolbox/actions/ai-report.md`, `help.md`, `SKILL.md` and `docs/ai-report-guide.md`. 658 and 659 both edit `work.md` and `work-reference.md`. 658, 659 and 660 all add bullets to `lessons-do-work-cli.md`. These are merge seams for the serial integrators, not reasons to serialize the builds, so no queue gate was added.
- Integration order (step 1): the three tool REQs 660, 659 and 658 first. They feed the later coordinator REQs being recaptured. 655 and 657 follow, and 657 merges last of the ai-report pair.
- REQ-656 (ai-report revise) currently `depends_on: [REQ-655]`. Session b1 plans to drop that edge. If the edge is still there, 656 waits for 655 to archive.

**Worktree verdicts**

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` (main): ACTIVE, the shared checkout, clean at 61e02428.
- `/private/tmp/do-work-portable-workflows` (branch `codex/portable-workflow-actions`, 7676a183, clean): FOREIGN, another tool's work. Leave it as it is.

**Heads-up for the first ten minutes**

- Session skill-do-work2-b1 is live in this same checkout. It will commit cancels and new captures under `do-work/`. Commit by exact path only, and re-read the highest REQ/UR id before writing one. (Next session.)
- The recaptured REQs will have new ids. Old REQ-662 to 666 and 671 to 673 may be cancelled, or folded into them. Trust the queue on disk, not this list. (Next session.)
- Integrators must wait for a quiet machine (1-minute load under 5, no running `maintainer-verify.sh`) inside one Bash sleep loop, never by ending their turn. (Next session's integrator briefs.)
