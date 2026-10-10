---
id: REQ-660
title: 'worktree new, status, merge and cleanup wrap the builder worktree lifecycle the actions already define'
status: claimed
created_at: 2026-10-09T21:14:04Z
user_request: UR-145
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
  - Route B
  - 6-file write set
  - 2 subsystems involved
  - 8 acceptance criteria
  calculated_at: 2026-10-10T13:19:16Z
related: [REQ-658, REQ-659, REQ-661]
batch: cli-ergonomics
required_lessons: [_dev/primes/lessons-releases.md]
claimed_at: 2026-10-10T12:52:20Z
route: B
write_set: ["skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go", "skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go", "skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/actions/fan-out-reference.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md"]
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
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers do-work-cli internals and the next step a finding suggests; families `destructive-next-argv`, `commit-preflight-before-side-effect`, `opaque-evidence-projection` and `fixture-cost-is-subprocess-spawning` fit a command that must never force a removal and must refuse before any side effect. The pre-dispatch brief quotes the four family rules the builder needs.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: its index row covers argv and quoting; the name derivation must stay a text operation, never REQ text in a shell line.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over budget; `slugged: partial`). Matching reason: its index row covers changing action routing and restated mechanisms; this REQ edits `actions/fan-out-reference.md` to name the new commands.

## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 3, What happened, Where the behaviour lives today Item 3, Proposed direction 3, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 3: "`worktree new|status|merge|cleanup REQ-N`: the builder worktree lifecycle that `actions/work-reference.md` already specifies step by step."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is specified line by line (`fan-out-reference.md` already spells out every git step, and the Assumptions settle the design questions), so no planning pass is needed. What needs discovery is where the command lives in do-work-cli and which existing worktree helpers it can reuse.

**Planning:** Not required


## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Required-lessons consult: `do-work/lessons-index.md` read. Kept `_dev/primes/lessons-releases.md` (666 tokens, matches the release requirement; the release is the integrator's). Dropped three partial satellites for budget (section above). No listed file was missing.

**Where the command goes.** `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` registers command families by package `Handlers()` maps. `internal/cleanup/` already owns worktree evidence (`lessons-do-work-cli.md` § Package routing: "plans safe Passes 0–4, consent-gated repairs, link repointing, and worktree evidence") and already has the helpers this command needs in `internal/cleanup/cleanup_git.go`:
- `parseWorktrees` (`:302`), NUL-safe parse of `git worktree list --porcelain -z` into `worktreeRecord{Name, Path, Head}`, naming a record by its branch or by a `worktree-agent-` basename.
- `worktreeClean` (`:345`), `git -C <path> status --porcelain=v1 -z --untracked-files=all` is empty.
- `requestIDFromWorktree` (`:351`), `worktree-agent-REQ-660-x` to `REQ-660`, rejecting `REQ-66` vs `REQ-660` prefix confusion.
- `cleanupGit` / `cleanupGitBytes` run git with a context in the repository root.
`internal/cleanup/cleanup_commands.go:21` is `Handlers()` returning `{"cleanup": handleCleanup}`; adding `"worktree"` there needs no `main.go` edit. `ApplyWorktreeRepairs` (`cleanup_git.go:209`) is Pass 5 and passes `--force` / `-D` under consent; the new command must not call it.

**REQ lookup.** `repositorymodel.DiscoverRepository(root)` returns `RequestsByID map[string][]*RequestFile` with `RelativePath`; `new` resolves REQ-N to exactly one file there and derives the suffix from the basename after `REQ-NNN-` and before `.md`. `merge` and `cleanup` must not need the REQ file: cleanup runs after finalization, when the REQ is already in `do-work/archive/`. They enumerate `worktree-agent-REQ-N-*` branches and worktrees instead.

**Result shape.** `internal/resultmodel/result_model.go:609-644` `CommandResult` has typed optional blocks (`LifecycleTiming`, `GateEvidence`, ...) and `renderText` prints each (`:1231` for timing). There is no generic slot for an operative name, a `<pre>`/`<merge_hash>` pair, or status rows; `RecordedChange{Path, Kind, Detail}` would force consumers to parse `Detail` strings (`opaque-evidence-projection` family in `prime-do-work-cli.md`).

**Git behaviour checked on this machine (git in PATH, 2026-10-10).** A symlink named `node_modules` is NOT matched by a `node_modules/` ignore line, so it shows as `?? node_modules` in the worktree and `git worktree remove` (no force) refuses with "contains modified or untracked files". `git merge --no-ff --no-commit <branch>` with nothing new prints "Already up to date.", exits 0 and leaves no `MERGE_HEAD`.

**Prose to change.** `skills/do-work/actions/fan-out-reference.md`: `### Naming` (`:35-41`, derivation and collision), `### The operative name` (`:43-45`), `### Where worktrees live` (`:47-49`), `### When to merge` steps 1 to 4 (`:63-77`), `### Cleanup — happy path` (`:93-95`). No contract test or `_dev/tests/` script pins any of this text (grep for `Never \`-D\``, `worktree-agent-REQ-NNN`, `Already up to date` finds nothing outside the file). The dispatch brief has no fixed template in that file: `:141` lists its contents only ("REQ body, worktree path, branch name, never-touch list, the commit subject rule, hand-back format").

**Siblings in this run.** REQ-689 (coordinate mode prose) edits `fan-out-reference.md` too, in `### Delegated integration` and `### Run directory, briefs and hand-backs` (`:124-151`), not the sections above. REQ-658 (finalize auto-manifest) and REQ-690 (run-status) cite `result_model.go`; either may add a field beside the new one. REQ-690 also reads `worktree-agent-REQ-NNN-*` branch tip age for its own status view; no shared helper is planned.

**Decisions made at pre-dispatch (best judgment, no question asked).**
- D-01: The command lives in `internal/cleanup/` as a new file plus a `"worktree"` entry in `cleanup.Handlers()`, reusing the four helpers above. Reasoning: reuse beats a second worktree parser, and the package already owns worktree evidence. Value: no `main.go` edit, no duplicated porcelain parser. Risk: the package name reads narrower than its content; reversible by moving the file.
- D-02: A typed `resultmodel` block (one `CommandResult` field plus its type and a short `renderText` block) carries operative name, worktree path, integration branch, links, `pre`, `merge_hash`, empty/conflict state, and status rows. Reasoning: the coordinator reads these from JSON; string-packed `RecordedChange.Detail` is the opaque projection the prime warns against. Value: exact fields for consumers. Risk: an adjacent-line merge seam with REQ-658/REQ-690 in `result_model.go`; the integrator resolves it by keeping both.
- D-03: `new` writes no brief skeleton. Reasoning: the REQ's Assumption says skip when no fixed template exists, and `fan-out-reference.md:141` is a content list, not a template; `new` also has no run directory to write into. Value: no invented template. Risk: the coordinator still writes briefs by hand, as today.
- D-04: The link config is `do-work/worktree-links`, read from the main tree only, one repo-relative path per line, blank lines skipped. Reasoning: the report's example; `do-work/` already holds root-level consumer files (`lessons-index.md`, `calibration-log.tsv`) and `new` runs as the orchestrator in the main tree, so *State stays home* is kept. A missing source path is skipped with a warning finding; an existing destination in the new worktree is never overwritten (warning, skipped). `cleanup` removes only links that are symlinks pointing at the main-tree path of a configured line.
- D-05: `merge` commits immediately after a clean `--no-commit` merge, with the subject `[REQ-N] merge builder branch <operative_name>`. A hand-back that carries integration seam lines uses the hand steps 2 to 4 (the prose stays valid). Reasoning: YAGNI; a seam flag can come later if the hand path proves costly. Value: smallest command. Risk: seam-bearing hand-backs keep the manual path.
- D-06: `merge` and `cleanup` treat the branch checked out in the main tree as the integration branch and refuse on detached HEAD or when that branch itself starts with `worktree-agent-`. Only `new` takes an override flag for its base. Reasoning: nothing persists the integration branch (`:43`), and `branch -d` is only meaningful from the branch merged into (`:93`).
<!-- D-XX counter: last used D-06. Next decision: D-07. -->

*Generated by pre-dispatch exploration*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go` (new) — the worktree new, status, merge and cleanup subcommands
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (new) — fixture-repository tests for the GREEN conditions
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go` (modify) — register the worktree command in Handlers
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — typed worktree result block and its text rendering
- `skills/do-work/actions/fan-out-reference.md` (modify) — name the commands beside the hand sequences they wrap
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify) — the internal/cleanup package-routing bullet names the new command

**Files I will NOT touch:** the main.go registration file, `skills/do-work/actions/cleanup.md` (Pass 5 policy unchanged), `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (REQ-659's file), any CHANGELOG, VERSION or version mirror (the integrator's release).

**Acceptance criteria (restated from REQ):**
- [ ] `worktree new REQ-N` derives `worktree-agent-REQ-N-<suffix>` from the filename slug as a text operation, appends `-2`, `-3` on a collision without deleting or forcing anything, creates branch and worktree from the integration branch under the sibling `<repo>-worktrees/` directory, links the paths in the optional `do-work/worktree-links`, and prints the operative name.
- [ ] An absent link config means no links.
- [ ] `worktree status` lists each `worktree-agent-REQ-*` worktree with ahead/behind against the integration branch, dirty or clean, and last-commit age.
- [ ] `worktree merge REQ-N` refuses a non-empty index, captures `pre`, refuses when the builder branch committed anything under `do-work/` (queue guard, before the merge), runs `git merge --no-ff --no-commit`, commits `[REQ-N] merge builder branch <operative_name>`, and prints `pre` and `merge_hash`. An empty hand-back reports empty and exits non-zero without a commit. A conflict stops with the merge in progress and lists the conflicted paths.
- [ ] `worktree cleanup REQ-N` refuses a dirty worktree (ignoring only its own links) before any side effect, removes its links, runs `git worktree remove` and `git branch -d` from the integration branch, then `git worktree prune`; `new` then `cleanup` leaves no worktree, branch or link behind.
- [ ] No subcommand ever passes `--force` or `-D`; every refusal is reported and the command stops.
- [ ] `fan-out-reference.md` names each command beside the hand sequence it wraps, and the hand sequence stays valid.
- [ ] Zero or several `worktree-agent-REQ-N-*` matches refuse in `merge` and `cleanup`; `--name <operative_name>` resolves several.
