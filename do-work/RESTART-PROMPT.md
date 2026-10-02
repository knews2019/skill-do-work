```
do-work run
This command is sufficient; everything below it is context.

Two REQs are already claimed by this checkout and part-way through the pipeline. Run canonical recover, then for each claimed REQ ask the classifier for its next phase (advance REQ-NNN --request-path <its do-work/working path>) and continue from there. Do not re-run phases whose sections already exist. Finish them one at a time: integration, release and archive are serial.

How to run this session without filling your own context. You are the orchestrator; every heavy read or long command happens outside your context:

1. Agents for everything that reads many files. Explore, build and review each run as a background agent (run_in_background). Give each one its full instructions in a brief file under this run's do-work/runs/<run>/ directory and pass only that absolute path in the prompt. Agents write their full output to a file in that run directory and reply with ONE line under 300 characters. Agent replies are cut at about 3500 characters, so never ask for a long reply.
2. Explore agents are read-only and cannot write files. Use a general-purpose agent when you need a findings file, or ask the Explore agent for under 3000 characters and save it yourself.
3. Builders work in their own git worktree outside the repo (../<repo>-worktrees/worktree-agent-REQ-NNN-<slug>). Their only main-tree write is their hand-back file. You merge.
4. Long commands go to the background with output in a log file in your scratchpad: the repository gate (bash _dev/tests/maintainer-verify.sh, about 2.5 minutes), heavy lanes, browser test runs. Read only the last lines of the log. Never pipe the gate to tail; redirect it to a file.
5. Never cat a large file. Action files are 500+ lines: grep their headings and read one section with sed. Parse do-work-cli JSON with a python one-liner that prints only outcome, phase, gate_id, state and finding codes.
6. Do not poll. Background jobs and agents notify you when they finish. To wait for a file or a log line, use one background until-loop, not repeated sleeps.
7. One gate at a time. Before starting the repository gate, check that no other maintainer-verify run is active (pgrep -f "^bash _dev/tests/maintainer-verify.sh").
8. Workflows are allowed for this queue if you prefer them to single agents: one workflow per REQ for explore, then build, then review, under 10 agents, each writing to the run directory. Integration (merge, qualify, test gate, finalize) stays in the main session.
9. Commit each coherent step (run artifacts before a merge, the session checkpoint at the end). Never push.
```

---

## Reference

Written 2026-10-02 at about 18:30Z by the session that ran `do-work run` on UR-132 (Board pages get URLs, and the Testing page shows free disk space). Run directory: `do-work/runs/work-2026-10-02-180407/`. Integration branch: `main`.

### In-flight REQs

**REQ-626 — Give every board page and lens its own URL** (`do-work/working/REQ-626-board-page-urls.md`)
- Merged: yes. Run-artifact commit 498dd997, builder commit 3100dea9, merge 505ec74a. Full merge range `498dd997..505ec74a`; the full merge hash for `supplied_commit` is `505ec74ac0d71f2554b52a8ef69259c12a7a6bcb`.
- Done: Triage, estimate, Plan skip, Exploration, Scope (synced, D-06), Pre-Flight, Implementation Summary, P-A-U boxes, Qualification (qualify + scope-drift satisfied), Testing (repository gate exit 0 at 505ec74a, 144s; focused probe satisfied), Decisions D-01 to D-06.
- Remaining, in order:
  1. Review. An independent reviewer was running when this was written. If `do-work/runs/work-2026-10-02-180407/REQ-626-review.md` exists, append its contents to the REQ as `## Review` and stamp `review_at`; otherwise run review-work in orchestrated mode on the range `498dd997..505ec74a`.
  2. Lessons Learned and Orientation (Route B: Lessons required). Lesson candidates: the probe-page URL check was copied into seven test files instead of one helper; Explore agents are read-only and cannot write their findings file.
  3. Heavy hold: record `commit:` with `record-commit-hash --request-path <path> --implementation-hash 505ec74ac0d71f2554b52a8ef69259c12a7a6bcb`, append `## Heavy Verification Plan` (lanes queue-kanban-javascript, queue-kanban-browser, staged-skills; range above).
  4. Continue REQ-627, then drain the heavy lanes at queue exhaustion from a detached checkout (`git worktree add --detach .git/work-run-<stamp>/drain-head HEAD`), with `QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"` exported, or the browser lane silently skips.
  5. Finalize with a release (VERSION is 0.305.60 now, so REQ-626 becomes 0.305.61 if it finalizes first). Helpers: `do-work/runs/work-2026-10-02-180407/helpers/prep-payloads.py` and `make-manifests.py` (set `REQ_ID`, `REQ_PATH`, `UR_ID=UR-132`, `UR_CLOSES=0` for the first REQ to finalize, `1` for the one that closes UR-132).
- Uncommitted files: none after the handoff commit (the working REQ, the hand-back and the review file are committed with this prompt).

**REQ-627 — Show free disk space on the Testing page** (`do-work/working/REQ-627-testing-page-free-disk-space.md`)
- Merged: no. No builder dispatched, no branch exists.
- Done: Triage (Route B), estimate (P50 25 min), Plan skip, Exploration (`do-work/runs/work-2026-10-02-180407/REQ-627-exploration.md`), Scope and write_set, required lesson `lessons-do-kanban.md#disk-space-blind-spot`.
- Remaining: Pre-Flight at the current tip (the 7917efbe pre-flight was for REQ-626; run the gate again after REQ-626's merge), then the builder brief and dispatch, from Step 6 of the work action on.
- Claimed early: the second queue-mode `advance` call of the run claimed it while REQ-626 was still building. The claim is valid; it only means `claimed_at` (18:04Z) overstates its queue wait.

### Parallelism

- REQ-626 and REQ-627 must not build at the same time: both write `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go`, and REQ-627's pre-flight must run on the tree that already contains REQ-626. REQ-626 is already merged, so only REQ-627 needs building and the question is moot; the resume command runs serially.
- No dependency gate is needed: REQ-626 is past building, and no other REQ is queued.
- Critical path: finish REQ-626 Review → hold → build REQ-627 → one heavy drain covering both → finalize both (two releases, 0.305.61 and 0.305.62, in finalize order).

### Worktree verdicts

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` — main checkout, ACTIVE.
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-626-board-page-urls` — ACTIVE, not removable yet. The branch is merged and the tree is clean, but REQ-626 is still in `do-work/working/`. After REQ-626 finalizes: `git worktree remove /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-626-board-page-urls && git branch -d worktree-agent-REQ-626-board-page-urls && git worktree prune`, run from `main`.

### Heads-up for the first ten minutes

- The repository gate's per-file 30s budget is close to its limit under load (slowest files 20.9s and 24.3s at load 7). Run it on a quiet machine; the first red run is retried once before any judgment.
- The qualifier reads backticked words in Scope and the Implementation Summary as file paths. Backtick only real paths there.
- `git merge` refuses while run artifacts are staged. Commit run artifacts first, then merge on an empty index.
- The focused-test probe must name the REQ's own tests (`-run ...`), not the whole package, or it exceeds the probe timeout.
- Finalization needs an empty index and `commit_paths` that are a superset of what the planner wants. Read the refusal for the missing paths and retry.
- `advance --checkpoint` leaves `do-work/CHECKPOINT.md` and `do-work/working/baseline.json` dirty at the end. Commit them as a final session checkpoint commit.
