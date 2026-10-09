---
id: REQ-660
title: 'worktree new, status, merge and cleanup wrap the builder worktree lifecycle the actions already define'
status: pending
created_at: 2026-10-09T21:14:04Z
user_request: UR-145
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-658, REQ-659, REQ-661]
batch: cli-ergonomics
---
# worktree new, status, merge and cleanup Wrap the Builder Worktree Lifecycle the Actions Already Define
## What
`do-work-cli worktree new|status|merge|cleanup REQ-N` runs the builder worktree lifecycle that `actions/fan-out-reference.md` already specifies step by step: naming with the collision suffix, creation outside the repo, optional per-repo links such as `node_modules`, the hand-back merge with its queue guard, and cleanup. Never `--force`, never `-D`; a refusal is reported and stops.
## Why
Report item 3 (UR-145 input, "What happened"): 85 hand-run `git worktree add` commands for a builder in about 31 sessions over 14 days. One consumer (2026-10-07) chained `git worktree add -b $name $WT/$name main && ln -s $PWD/node_modules $WT/$name/node_modules && ln -s $PWD/env.vars $WT/$name/env.vars`. Nothing in `actions/` or `docs/` covers dependency directories, so every session adds its own `ln -s` lines.
## Verified Facts (checked at 0.305.87)
The report cites `actions/work-reference.md:416-445` at 0.305.84; 0.305.87 moved that text unchanged into `skills/do-work/actions/fan-out-reference.md`:
- `:37` naming: directory basename and branch are the same string `worktree-agent-REQ-NNN-<suffix>`, suffix derived from the filename slug as a text operation, never by piping REQ text through `tr`/`sed`.
- `:39` collision: append `-2`, `-3`; never delete or force.
- `:43` the name actually created is the operative name for every later step.
- `:47` worktrees live outside the repo (`../<repo>-worktrees/...`), never nested.
- `:59` integrate with `git merge --no-ff`, never rebase.
- `:73` hand-back step 2: queue guard `git diff --name-only <pre>...<operative_name> -- do-work/` before the merge, then `git merge --no-ff --no-commit <operative_name>`; `Already up to date.` means an empty hand-back.
- `:93` cleanup: `git worktree remove <path>` (no `--force`), `git branch -d <operative_name>` from the integration branch, `git worktree prune`. "Never `-D`, never `--force`. Report the refusal and stop."
- `:16` the branch rung (no second directory) uses the same name as a plain branch.
## Detailed Requirements
1. `worktree new REQ-N` derives the name and collision suffix exactly as `fan-out-reference.md:37-39` says, creates the branch and worktree from the integration branch outside the repo (`:47`), links the paths listed in an optional per-repo config, writes the brief skeleton, and prints the operative name.
2. The optional per-repo link config is a new file with one repo-relative path per line (the report's example: `do-work/worktree-links`, with lines such as `node_modules` and `env.vars`). Absent file means no links.
3. `worktree status` lists each `worktree-agent-REQ-*` worktree with ahead/behind against the integration branch, dirty or clean, and last-commit age.
4. `worktree merge REQ-N` runs the `fan-out-reference.md:73` sequence: clean index, `<pre>` capture, queue guard, `git merge --no-ff --no-commit`, commit with a `[REQ-N]` message, then prints `<pre>` and `<merge_hash>`.
5. `worktree cleanup REQ-N` removes the links, runs `git worktree remove` and `git branch -d` from the integration branch, then `git worktree prune`.
6. Never `--force`, never `-D`. Any refusal is reported and the command stops.
7. Change the action lines in `fan-out-reference.md` that describe these hand sequences to name the commands, keeping the hand path valid.
8. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Naming, merge and cleanup rules do not change. This REQ wraps them.
- No new frontmatter field and no new status.
- `merge` refuses when the builder branch committed anything under `do-work/` (the queue guard).
- This REQ is its own release; do not fold REQ-658, REQ-659 or REQ-661 into it.
## Assumptions (recorded at capture, no questions asked)
- The integration branch is the branch checked out in the main working tree when `new` runs, overridable by a flag. `merge` and `cleanup` refuse unless that integration branch is checked out, because `branch -d` from any other branch gives a false merged-ness answer (`fan-out-reference.md:93`).
- The operative name is not persisted (`fan-out-reference.md:43` says nothing persists it). `merge` and `cleanup` take REQ-N and find the operative name by enumerating `worktree-agent-REQ-N-*` worktrees and branches. Zero or several matches refuse; an explicit `--name <operative_name>` flag resolves several.
- The worktree parent directory is the sibling `../<repo>-worktrees/` that `fan-out-reference.md:47` names first.
- The link config location in the report (`do-work/worktree-links`) is an example. The builder may place it elsewhere if `do-work/` holding consumer config conflicts with an existing rule, and records the choice in Decisions. Linked paths are expected to be git-ignored in the worktree; `cleanup` removes only links it created.
- "Writes the brief skeleton" means the builder dispatch brief that `fan-out-reference.md` describes. If no fixed brief template exists there, `new` skips this step and the builder records that in Decisions; no new template is invented.
- On a merge conflict, `merge` stops with the merge in progress and reports the conflicted paths. Resolving a conflict stays with the integrator, as today.
- An empty hand-back (`Already up to date.`) makes `merge` report empty and exit non-zero without committing.
- Only the worktree rung is in scope. The branch rung (`fan-out-reference.md:16`) stays prose.
## Dependencies
None. Independent of REQ-658, REQ-659 and REQ-661.
## Builder Guidance
Certainty is high on the git sequence; it is already specified line by line. Lower on the link config home and the brief skeleton. Latitude: flag names, the status table layout, and how the operative name is found.
## Red-Green Proof
**RED prompt/case:** In a fixture repo with a queued REQ-N, run `worktree new REQ-N`, commit a change on the builder branch, run `worktree merge REQ-N`, then `worktree cleanup REQ-N`.
**Why RED now:** do-work-cli has no `worktree` command (`do-work-cli help` lists none).
**GREEN when:** `new` then `cleanup` leaves no worktree, branch or link behind. `merge` refuses when the builder branch committed anything under `do-work/`. A dirty worktree makes `cleanup` refuse. No command ever runs `--force` or `-D`.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers do-work-cli internals and destructive next-step argv; family `destructive-next-argv` fits a command that must never force a removal.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: its index row covers argv and quoting; the name derivation must stay a text operation, never REQ text in a shell line.
## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 3, What happened, Where the behaviour lives today Item 3, Proposed direction 3, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 3: "`worktree new|status|merge|cleanup REQ-N`: the builder worktree lifecycle that `actions/work-reference.md` already specifies step by step."*
