```
do-work run

This command is sufficient; everything below it is context.

Three requests are claimed in do-work/working/ and all three are merged, verified green and
reviewed. They are held for the heavy-lane drain and nothing about them needs rebuilding. Seven
more sit pending in do-work/queue/. Run the loop; it will drain the held three at queue
exhaustion and finalize each one.

Four things to know before you start, because each will otherwise cost you time:

1. The fast gate now caches. `bash _dev/tests/maintainer-verify.sh` prints one line per Go test
   stage saying EXECUTING or REUSED with a reason, and a gate wall line. On a matching tree it
   takes about 21 seconds instead of about 96. That reuse is NOT yet trustworthy for changes under
   do-work/: REQ-592 in the queue is the impact-critical fix, and until it lands you can disable
   reuse entirely for any run with DO_WORK_FAST_STAGE_REUSE=off, which is verified to work and
   needs no code change. Prefer that for any gate run whose verdict you are going to rely on.

2. Canonical recover currently REFUSES with FINALIZATION-DISCOVERY-AMBIGUOUS. Its blocked_paths are
   three untracked hand-back files under do-work/runs/ that belong to OTHER live sessions in this
   same checkout (REQ-588, REQ-589, REQ-590). They are not finalization tails and they are not
   yours. Judge it and continue, exactly as the previous session did — do not delete them, do not
   commit them, and do not let the refusal park the queue. If they are still there and still
   foreign when you finish, say so in your own hand-off.

3. Never run the canonical gate while another gate process is running on this machine. Check with
   `ps -Ao args= | grep -c '[m]aintainer-verify'` and with `uptime` first. Five gate runs in the
   previous session failed on per-test-file wall-clock budgets, on four different files, none of
   them touched by any request in the run — purely from sibling sessions' load. The same files
   pass in a quiet window. A load average above about 10 predicts a budget failure.

4. Capture the gate's own exit status directly. Never pipe it to tail: the shell reports the
   pipe's status, and a red gate reads as green. That mistake cost the previous session a wasted
   cycle.

The browser heavy lane needs QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
at drain time. Without it the lane reports skipped, and a skip is not a pass.
```

---

## Reference

### What is done, and what remains

Three requests are claimed, merged, verified green, independently reviewed, and held at Step 7.7
for the heavy-lane drain. None needs rebuilding. Each has its full pipeline record in its own file
in `do-work/working/`, and every explorer, planner, builder and reviewer report for the run is
committed under `do-work/runs/work-2026-09-05-170806/`.

| REQ | What it delivered | Merge range | `commit:` | Review |
|---|---|---|---|---|
| REQ-583 | Tests that pin three evidence-gate behaviours which could previously be deleted with the package staying green | `a22ddfcf..722f5ada` | `722f5ada` | 96%, Pass, 6 findings all report-only |
| REQ-587 | The Timeline view scrolls once, with the date axis pinned, like the Activity view | `93ec7792..8fad73b2` | `8fad73b2` | 92%, Pass, 6 findings all report-only |
| REQ-591 | The fast gate skips a stage whose inputs have not moved; the SessionStart probe stopped relinking a byte-identical binary nine times | `c2a74d2f..fcf07ea4` (cumulative, one remediation) | `fcf07ea4` | 60%, **Partial**, Risk **Critical**, one impact-critical finding now queued as REQ-592 |

Remaining, all `pending`: REQ-592 (the critical fix above), REQ-486, REQ-552, REQ-554, REQ-555,
REQ-556, REQ-557, REQ-558.

### The one thing that is actually wrong right now

REQ-591 shipped a real false green and it is live in this repository's gate. `_dev/tests/fast-stages.json`
declares `do-work/` as a tree no fast stage reads. Both stages read it. The reviewer reproduced it
end to end: warm evidence store, one newline appended to `do-work/archive/UR-003/input.md`, the
do-work-cli stage reports `REUSED`, the gate prints `Maintainer verification passed.` and exits 0 —
while `TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass` fails on that same tree.

REQ-592 carries the fix and the RED/GREEN case. Note that two existing assertions currently pin the
wrong behaviour and must move with it: `fast_stage_evidence_test.go`'s `queue state changed` case
and `_dev/tests/fast-stage-reuse-behavior.sh`'s `queue state alone still reuses` case.

Mitigation until then, verified working: `DO_WORK_FAST_STAGE_REUSE=off`.

### Parallelism

The paste block says `do-work run` with no `--fan-out`, deliberately. Six of the seven pending
requests — REQ-552, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558 — each append one assertion to the
same file, `_dev/tests/audit-lockins.sh`. Their `write_set` fields overlap on it, so any two of them
running concurrently produce a merge conflict on every hand-back. That constraint is not encoded as
dependency gates because it is a merge-collision constraint rather than a real dependency, and
inventing `depends_on` edges to express it would corrupt the dependency graph's meaning. Serial is
the honest encoding, and the gate is now fast enough that serial costs little.

The dependency gates that ARE real are already in the queue and need nothing from you: REQ-555
depends on REQ-554, REQ-557 depends on REQ-550 and REQ-552, REQ-558 depends on REQ-557.

Critical path: REQ-592 first — it fixes the gate you are about to rely on. Then the audit batch in
dependency order. REQ-486 (collapsible UR groups with progress summaries) is the only board-side
request and the only one that could safely run beside an audit request, but it is also the largest
single piece of work left and is `priority: later`.

### Pre-exploration already done for four pending requests

Read these before triaging; they were produced this run and each re-verified its request's claims
against HEAD rather than against the audited commit:

- `do-work/runs/work-2026-09-05-170806/REQ-556-exploration.md` — **the request's baseline is stale.**
  It claims 9 prose sites; there are 7 at HEAD. Two were already removed by the REQ-504→REQ-506
  chain. It also establishes that "keep one sentence naming the three finding codes" is an
  *addition*, not a retention: those codes appear in no shipped prose today.
- `do-work/runs/work-2026-09-05-170806/REQ-552-exploration.md` — both exec sites still present at
  HEAD, the claim holds; includes the exact replacement primitives and a paste-ready lock-in
  assertion. The lock-in must keep `--glob '!*_test.go'` or it is red on day one.
- `do-work/runs/work-2026-09-05-170806/REQ-554-exploration.md` — the shared-line count, the
  destination guide's structure, and the ratchet constants that need re-baselining.
- `do-work/runs/work-2026-09-05-170806/REQ-486-exploration.md` — the prior UR-lens implementation,
  where a shared summary function belongs, and whether the board can already read the nested P50.

### Worktree verdicts

| Path | Verdict |
|---|---|
| `../skill-do-work2-worktrees/worktree-agent-REQ-583-…` | **ACTIVE** — merged and clean, but its claim is still in `do-work/working/`. Remove only after finalization: `git worktree remove ../skill-do-work2-worktrees/worktree-agent-REQ-583-pin-the-evidence-gate-remedy-redirection-guard-and-interrupted-path` then `git branch -d worktree-agent-REQ-583-pin-the-evidence-gate-remedy-redirection-guard-and-interrupted-path` |
| `../skill-do-work2-worktrees/worktree-agent-REQ-587-…` | **ACTIVE** — same; remove after finalization |
| `../skill-do-work2-worktrees/worktree-agent-REQ-591-…` | **ACTIVE** — same; remove after finalization |
| `.git/work-run-20260905-1201/worktree-agent-REQ-573-activity-drawer` | **FOREIGN** — unmerged, clean, from an earlier run that is not this session's. Leave it. `actions/cleanup.md` Pass 5 owns it. |
| the previous session's scratchpad measurement worktrees | **already removed** — two detached, branchless worktrees under the session scratchpad, taken out with `git worktree remove` plus `git worktree prune` before this handoff was finished. `git worktree list` should show only the four in the rows above plus the main checkout. |

All three builder branches are in `git branch --merged main`. Do not `-D` any of them; `-d`'s
refusal is the only assertion that the integration actually happened.

### Heads-up list — what bites in the first ten minutes

- **The gate caches now and its reuse is wrong for `do-work/` changes.** You. See above.
- **`recover` refuses on three foreign hand-backs.** You, by judging and continuing. Not a stop.
- **Another session shares this checkout and commits every few minutes.** Expect `HEAD` to move
  under you between a gate run and its green record; check `git diff --name-only <gate-rev> HEAD`
  and re-run only if source moved. Expect uncommitted foreign edits in `skills/` — leave them, and
  work in a worktree so they cannot reach your build.
- **A pre-existing intermittent:** `TestLaneMutationCannotPublishOrReuseSuccess/commit=true` failed
  once in four full-gate runs and passes 6/6 in isolation at both revisions. Independently judged
  pre-existing by the reviewer, not caused by REQ-591. If you see it, it is known.
- **One line of the REQ-591 exploration is wrong and is corrected in place.** Its section 8 says
  `/usr/bin/time -p` under-reports process-tree CPU and that bash's `time` keyword is needed. That
  did not reproduce — both agree, including on the toolchain case that matters. An orchestrator
  correction with the numbers is appended to `do-work/runs/work-2026-09-05-170806/REQ-591-exploration.md`.
  Do not build a measurement protocol on the original claim.
- **REQ-556's own baseline is stale** (9 claimed, 7 actual). Re-verify before building, and expect
  to record the new baseline in the lock-in rather than the request's number.
