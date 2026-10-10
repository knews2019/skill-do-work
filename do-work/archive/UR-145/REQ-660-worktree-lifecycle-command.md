---
id: REQ-660
title: 'worktree new, status, merge and cleanup wrap the builder worktree lifecycle the actions already define'
status: completed
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
builder_handback_at: 2026-10-10T13:30:43Z
integration_at: 2026-10-10T13:33:38Z
review_at: 2026-10-10T14:04:09Z
remediation_at: 2026-10-10T14:05:28Z
re_review_at: 2026-10-10T14:38:04Z
route: B
kb_status: pending
write_set: ["skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go", "skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go", "skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/actions/fan-out-reference.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md"]
commit: bd9ecb2d5b118c34380b7e7fc8a7c0c3690b73ad
heavy_verified_at: 2026-10-10T14:38:27Z
heavy_verified_revision: bd9ecb2d5b118c34380b7e7fc8a7c0c3690b73ad
completed_at: 2026-10-10T14:39:20Z
release_at: 2026-10-10T14:39:20Z
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
- [x] **[PLAN]:** (from the builder hand-back) Add the typed result block first so the tests compile, write the five tests against `Handlers()` through `commandruntime.NewRuntime`, run RED (missing handler), then implement one file in `internal/cleanup/`. All git work runs against the main tree, found as the first `git worktree list --porcelain -z` record. Each subcommand checks every precondition before its first side effect: `new` checks request uniqueness, base branch, link config and the name before `git worktree add`; `merge` checks MERGE_HEAD, index, ahead count and queue guard before `git merge`; `cleanup` checks merged-ness and dirt (owned links filtered out) before removing any link. Refusal `next_argv` is always a read-only git command.
- [x] **[APPLY]:** (from the builder hand-back) Implemented as planned. One correction during GREEN: `RequestFile.FilenameID` is the normalized id (not the raw filename prefix), so the suffix is derived by stripping `^REQ-[0-9]+-?` from the basename with a regex instead.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff bd56c4b0 --stat`: 6 files, 789 insertions, 3 deletions. `gofmt -l` on `internal/cleanup` and `internal/resultmodel`: no output. `go vet` on both packages: exit 0. `go test -count=1 ./internal/resultmodel/` exit 0 (0.81 s); `./internal/cleanup/` exit 0 (16.3 s); the five named tests exit 0 (2.65 s); `REQ-660-probe.sh` exit 0 (4.61 s); `git diff --check` exit 0. All six changed files read in the final diff; no debug artifacts; `"--force"` and `"-D"` do not appear in `worktree_lifecycle.go`.
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

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/actions/fan-out-reference.md` (modified)
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified)

**What was done:** Added a `do-work-cli worktree` command with `new`, `status`, `merge` and `cleanup` subcommands that run the builder worktree lifecycle `fan-out-reference.md` already defines, checking every precondition before the first side effect and never passing `--force` or `-D`. Results travel in a new typed `worktree` block of the command result (operative name, paths, links, `pre`/`merge_hash`, conflicted paths, status rows); five fixture-repository tests cover the GREEN conditions; `fan-out-reference.md` names each command beside the hand sequence it wraps, and the `internal/cleanup/` routing bullet in the lessons satellite names the new command. After review (remediation re-merge `bd9ecb2d`, builder-branch commit `09745c7a` by the integrator): `merge` refuses a builder commit that touches a configured link path (`WORKTREE-LINK-COMMITTED`) before `git merge`, `new` warns when a created link is not git-ignored (`WORKTREE-LINK-NOT-IGNORED`), a sixth test pins the refusal, and `fan-out-reference.md` names the ignore-line form and keeps the first `<pre>` on a remediation re-merge.

## Decisions

(from the builder hand-back, verbatim)

- D-07: `new` returns outcome `success` (exit 0) when the worktree was created, even when a link was skipped or a collision suffix was used. Those are `warning` findings (`WORKTREE-LINK-SKIPPED`, `WORKTREE-NAME-COLLISION`). Reason: exit 1 would tell a script the worktree was not created. DECIDE & STATE.
- D-08: Refusal findings use a read-only `git ...` argv as `next_argv` (for example `git -C <main> status --short`, the queue-guard diff itself), never `do-work-cli worktree status`. Reason: `NormalizeResult` clears a refusal's `next_argv` whose verb equals the command (`worktree`), so that suggestion would silently disappear. `worktree status` stays in `verification_argv`. DECIDE & STATE.
- D-09: A merge that stops on conflicts returns outcome `findings` (exit 1, code `WORKTREE-MERGE-CONFLICT`, `conflicted_paths` set), not `refused`, because a merge is now in progress (a side effect happened). A `git merge` error with no conflicted paths returns `failure` with git's stderr. DECIDE & STATE.
- D-10: `status` uses each worktree's `HEAD` (from the porcelain record) for ahead/behind and last-commit time, not the branch name. For an attached worktree these are the same commit; a detached builder worktree still gets a row. DECIDE & STATE.
- D-11: `status` filters owned links out of the dirty check, using the same helper as `cleanup`. Without this, every worktree with a link reads as dirty, and test 1 requires "clean" after `new` created a link. So `worktreeClean` is not reused; a filtered porcelain read (`dirtyPaths`) replaces it in this file only. DECIDE & STATE.
- D-12: `status` refuses on a detached main tree, like `merge` and `cleanup`, because there is no integration branch to count against. Only `merge` and `cleanup` refuse a `worktree-agent-*` branch checked out in the main tree; `new` and `status` do not. DECIDE & STATE.
- D-13: `merge` also refuses when `MERGE_HEAD` exists (brief item: "a merge is already in progress"), checked first, before the index check. DECIDE & STATE.
- D-14: The suffix keeps only `[a-z0-9-]` after lowercasing: other characters are dropped, not replaced by `-`, and leading/trailing `-` are trimmed. Spaces therefore join words (`odd name` becomes `oddname`). If the slug is empty, the name is `worktree-agent-REQ-N`. Reversible. DECIDE & STATE.
- D-15: Link config lines are trimmed and cleaned. An absolute path, a `..` segment, or a line that cleans to `.` refuses the whole command before anything is created. There is no comment syntax (YAGNI). DECIDE & STATE.
- D-16: All git work runs in the main tree (the first porcelain record), including `merge`, even when the command is started from a linked tree. REQ discovery and the link config also read the main tree. This follows "State stays home". DECIDE & STATE.
- D-17 (integrator, after review): Review F1 (impact-critical: a committed link replaced the main tree's ignored directory on merge) and F2 (remediation `pre`) were fixed on the builder branch and re-merged with the same `<pre>` instead of minting a follow-up REQ. Reasoning: the fix is about 20 lines plus one test, it closes a reproduced data-loss path in the command this REQ ships, and releasing the known defect would advertise it. F3 to F8 stay report only. DECIDE & STATE.

## Discovered Tasks

(from the builder hand-back; the builder's `impact-minor` token is not in the impact vocabulary, so the integrator re-stamped each line by the two questions in `actions/review-work.md` Step 10)

- **impact-user-visible** `cleanup` Pass 5 (`ApplyWorktreeRepairs` in `internal/cleanup/cleanup_git.go`) uses `worktreeClean`, which counts the symlinks `worktree new` creates as dirt. A merged leftover that still has its links would be reported as dirty and need consent there, even though `worktree cleanup` would remove it. Confirmed by code reading: `worktreeClean` is an unfiltered `status --porcelain`. → report only
- **impact-user-visible** The Crash Recovery sweep (`actions/work-reference.md` → Crash Recovery (Step 1)) has the same dirt reading for linked leftovers. It is not checked in code here. → report only
- **impact-negligible** `worktree cleanup` leaves the empty `<repo>-worktrees/` parent directory after the last worktree is removed. The hand path does the same. → report only

## Qualification

**Gate records (`advance --diff-range bcf12cf3..08f005e2`):** `qualify` satisfied (success), `scope-drift` satisfied (success). Two `QUALIFY-NEW-FILE-UNWIRED` warnings, judged false positives: `worktree_lifecycle.go` is a same-package Go file whose `handleWorktree` is registered in `cleanup_commands.go` `Handlers()`, and `worktree_lifecycle_test.go` is a `_test.go` file the Go test runner discovers by convention. No debug-artifact, P-A-U or output-primitive findings.

**Scope:** declared `write_set` (6 paths) equals the touched set in `git diff bcf12cf3..08f005e2 --stat` (6 files, 789 insertions, 3 deletions). No `do-work/` path on the builder branch (queue guard printed nothing before the merge).

**Requirement trace (read against the merged files):**
1. `worktree new REQ-N`: `createWorktree` resolves exactly one request file through `repositorymodel.DiscoverRepository`, derives the suffix from the basename as a text operation (`strings.Map` to `[a-z0-9-]`), loops `-2`, `-3` while `lifecycleNameTaken` sees a branch, worktree or path with that name (nothing deleted), runs `git worktree add -b <name> <repo>-worktrees/<name> <base>` from the main tree's branch (or `--from`), links configured paths, returns `worktree.operative_name`. The brief skeleton is skipped per D-03 (no fixed template in `fan-out-reference.md`), which the REQ's Assumptions allow.
2. Link config: `do-work/worktree-links`, one path per line; absent file returns no links (`readLinkConfig`, `os.ErrNotExist`); absolute or `..` lines refuse before any side effect.
3. `worktree status`: one row per `worktree-agent-*` record with ahead/behind (`rev-list --left-right --count`), dirty (owned links filtered), last-commit unix time; text render shows age.
4. `worktree merge`: refuses MERGE_HEAD, staged paths, zero ahead (empty hand-back, no commit, exit non-zero), and any `do-work/` path in `<pre>...<name>` before `git merge --no-ff --no-commit`; commits `[REQ-N] merge builder branch <name>`; reports `pre` and `merge_hash`. Conflict leaves the merge in progress with `conflicted_paths`.
5. `worktree cleanup`: refuses unmerged (`merge-base --is-ancestor`) and dirty (owned links excluded) before any side effect, then removes owned links, `git worktree remove` (no force), `git branch -d`, `git worktree prune`; refuses a detached or builder-branch main tree, so `-d` runs from the integration branch.
6. `grep` of `worktree_lifecycle.go`: no `--force`, no `-D`; every refusal returns through `stop`.
7. `fan-out-reference.md`: one added paragraph each in Naming, Where worktrees live, When to merge (after step 4) and one sentence in Cleanup — happy path; every hand sequence is unchanged.
8. Release: handled at finalization (patch bump, own changelog entry).

**Red-Green Proof trace:** the five tests in `worktree_lifecycle_test.go` cover the captured GREEN conditions: `new` then `cleanup` leaves nothing behind, `merge` refuses a `do-work/` commit, `merge` reports an empty hand-back without a commit, `cleanup` refuses a dirty worktree before any side effect, and `new` appends `-2` on a collision.

**Remediation re-merge (review F1, F2), cumulative range `bcf12cf3..bd9ecb2d`:** `advance` no longer accepts qualify input at this phase (`ADVANCE-GATE-INPUT-IRRELEVANT`), so the integrator ran the same `qualify --request-path <P> --diff-range bcf12cf3..bd9ecb2d` handler directly: success, only the two judged `QUALIFY-NEW-FILE-UNWIRED` warnings above. Touched set is still the 6 declared files (`git diff bcf12cf3..bd9ecb2d --stat`: 840 insertions, 3 deletions); the delta touches `worktree_lifecycle.go`, `worktree_lifecycle_test.go` and `fan-out-reference.md` only. Queue guard before the re-merge printed nothing.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `08f005e2` (machine quiet before launch: 1-minute load 2.98, no other gate running; load rose to 7.43 during the run), then `advance REQ-660 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-660-probe.sh`.
**Result:** ✓ Repository gate passed on the first run (exit 0, gate wall 131 s, do-work-cli 887 tests, slowest file `internal/finalization/finalization_recovery_test.go` 19.49 s under the 30 s limit). Probe exit 0: the five named tests reported `--- PASS` (none skipped), `worktree_lifecycle.go` has no `"--force"` or `"-D"`, and `do-work-cli help` lists `worktree`. Gate records `green-gate`, `scope-drift` and `run-blocked-check` satisfied.

**After the review-fix re-merge (`bd9ecb2d`):** the same gate argv ran again on a quiet machine (load 2.04 before launch, 5.29 at the end): exit 0, gate wall 135 s, 888 do-work-cli tests, slowest file 20.80 s. `advance` refuses gate input past this phase, so the green record was written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` at `bd9ecb2d`, and `REQ-660-probe.sh` was run directly: exit 0. `go test -count=1 ./internal/cleanup/ ./internal/resultmodel/` in the builder worktree: exit 0.

**Red-green validation:** (from the builder hand-back, traced to `## Red-Green Proof`)
- `TestWorktreeLifecycleNewMergeCleanupLeavesNothingBehind`, `TestWorktreeMergeRefusesBuilderCommitUnderDoWork`, `TestWorktreeMergeReportsEmptyHandBackWithoutCommit`, `TestWorktreeCleanupRefusesDirtyWorktreeBeforeAnySideEffect`, `TestWorktreeNewAppendsNumericSuffixOnCollision` (`internal/cleanup/worktree_lifecycle_test.go`): ✗ before the handler existed (exit 1, each `worktree [new REQ-41] carried no typed worktree result: exit 2` with the runtime's `UNKNOWN-COMMAND` finding) → ✓ after (exit 0, no skips).
- `TestWorktreeMergeRefusesCommittedLinkBeforeItReplacesTheMainTreePath` (review F1): ✗ before the fix (exit 1, `new` carried no `WORKTREE-LINK-NOT-IGNORED`; with only the merge guard disabled, `merge: exit 0 outcome success findings []`) → ✓ after (exit 0).
- Queue-guard mutation check: with the guard disabled, `TestWorktreeMergeRefusesBuilderCommitUnderDoWork` failed (`merge: exit 0 outcome success findings []`); restored, it passed.

**New tests added:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (five tests through the real `commandruntime` seam with `--format json`)
- `TestWorktreeMergeRefusesCommittedLinkBeforeItReplacesTheMainTreePath` (same file, integrator's review fix for F1)

**Heavy verification plan:**
- Range: bcf12cf300e2694b62c5ee36846b186beeff0062..bd9ecb2d5b118c34380b7e7fc8a7c0c3690b73ad (re-planned after the review-fix re-merge; same four lanes as the first plan at `08f005e2`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 60%** | 2026-10-10T14:04:09Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 80% |
| Test Adequacy | 80% |
| Scope | 95% |
| Risk | Critical |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `worktree merge` merges a builder commit containing a configured link: with a `node_modules/` ignore line and `git add -A`, the main tree's ignored directory is deleted and replaced by a self-looping committed symlink (reproduced; fix = link-path guard beside the queue guard at `worktree_lifecycle.go:291-298`, a `check-ignore` warning in `new`, one doc sentence) — impact-critical → follow-up REQ for the integrator to create (reviewer created none, per brief)
- F2 the remediation re-merge reports this merge's parent as `pre`, and the new `fan-out-reference.md:83` paragraph omits the keep-the-first-`<pre>` rule — impact-user-visible → report only
- F3 linked worktrees read as dirty in queue-kanban `verify.go:1109-1116` (false "uncommitted changes" message), `corehelpers/handoff.go:87`, cleanup Pass 5 `cleanup_git.go:345`, and the Crash Recovery prose; all fail toward consent — impact-user-visible → report only

**Minor findings:** F4 `status`/`cleanup` exit 2 on a prunable record and hide every row — impact-user-visible → report only; F5 `cleanup` removes links before a locked-worktree refusal from `git worktree remove` — impact-negligible → report only; F6 link config accepts `do-work` paths and links the live queue into a builder tree — impact-negligible → report only; F7 (nit) `dirtyPaths` rename skip is wrong after `splitZero` sorts, so evidence paths only — impact-negligible → report only; F8 (nit) anti-bloat: `links_skipped`/`empty_hand_back` repeat findings, four failure codes could be two — impact-negligible → report only
**Acceptance:** Pass — Implementation and Integration stages: focused tests exit 0, and the fixture new → commit → status → merge → cleanup, empty hand-back, queue guard, dirty and unmerged cleanup, conflict, and bad-argument refusals all behave as specified, with read-only `next_argv` only
**Restatement sweep:** redefined the builder worktree add / hand-back merge / cleanup hand sequences (now wrapped by `do-work-cli worktree`), the new `do-work/worktree-links` file, and what "dirty" means for a linked builder worktree; stale: `fan-out-reference.md:83` (remediation `<pre>`, F2), Where worktrees live (ignore-line form, F1), dirt readers in `verify.go`, `handoff.go`, `cleanup_git.go` and Crash Recovery (F3); consistent: `crew-members/background-agents.md:186`, `actions/cleanup.md` Pass 5, `actions/work.md:574-575`, `actions/restart-with-parallel-handoff.md:55`, `docs/cleanup-guide.md`
**Suggested testing:** 4 items
**Follow-ups created:** None by the reviewer. F1 is impact-critical and goes to the integrator to create (7 other findings report only)

Full report: `do-work/runs/work-2026-10-10-131527/REQ-660-review.md`.

*Reviewed by review-work action*

### Re-review (delta 08f005e2..bd9ecb2d) | 2026-10-10T14:38:04Z

F1 and F2 were fixed on the builder branch and re-merged with the same `<pre>` (D-17), so no follow-up REQ was created for F1.

- Overall: 92% | Acceptance: Pass (F1 fixture against the merged code: `merge` refuses `WORKTREE-LINK-COMMITTED` before `git merge`, HEAD unmoved, no MERGE_HEAD, main `node_modules/` intact; `node_modules` ignore line gives no warning and new → merge → cleanup works; the new test fails with the guard disabled)
- F1: closed. F2: closed.
- F9 the hand-path step 2 has no link check, and only the ignore-line sentence protects a hand merge — impact-negligible → report only
- F10 the literal link pathspec refuses a legitimate commit to a tracked or skipped link path, and it treats glob characters as magic — impact-negligible → report only
- F11 the shared lifecycle fixture uses the `shared-deps/` form, so every `new` in the tests warns — impact-negligible → report only

Full re-review report: `do-work/runs/work-2026-10-10-131527/REQ-660-rereview.md`.

## Lessons Learned

**What worked:** Reusing the cleanup package's porcelain parser and git runner kept the command to one new file, and checking every precondition before the first side effect made each refusal leave the repository exactly as it was (the dirty-cleanup test pins this). Counting `HEAD..<branch>` instead of reading `git merge` output caught the empty hand-back, because `git merge --no-ff --no-commit` prints "Already up to date." and exits 0.
**What didn't:** The first build trusted the REQ's assumption that linked paths are git-ignored. A `node_modules/` ignore line matches only a directory, never the symlink `worktree new` creates, so a builder's `git add -A` committed the link and merging it made git delete the main tree's ignored `node_modules/` directory (review F1, reproduced). The qualify and scope gates read the range only after the merge, so no gate could have seen it. Fixed by a link-path guard before `git merge` plus a `check-ignore` warning in `new`.
**Worth knowing:** A refusal finding whose `next_argv` verb equals the command's own name is cleared by `NormalizeResult`, so `worktree` refusals point at a read-only `git` argv and keep `worktree status` in `verification_argv` (D-08). Git treats ignored files as expendable on merge and checkout: any tool that creates untracked paths inside a worktree must guard the merge, not rely on the ignore file. Cleanup Pass 5, the board's verify probe and the handoff survey still read owned links as dirt (F3, report only).

## Orientation

Now a coordinator can create, list, merge and remove builder worktrees with `do-work-cli worktree new|status|merge|cleanup REQ-N`, with optional dependency links from `do-work/worktree-links`; it lives in the do-work-cli `internal/cleanup` package beside the worktree evidence helpers and is documented in `actions/fan-out-reference.md` → Worktree Dispatch Mode, where the hand sequences stay valid. [MAP CHANGED]: a new CLI command and a new optional per-repository config file. Prime spot-check: `prime-shell-commands.md`, `prime-action-files.md` and `prime-releases.md` name no path this change moved or removed; none is stale.

## Heavy Verification Plan

- Base: bcf12cf300e2694b62c5ee36846b186beeff0062
- Target: bd9ecb2d5b118c34380b7e7fc8a7c0c3690b73ad
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: bd9ecb2d5b118c34380b7e7fc8a7c0c3690b73ad
- Execution revision: bd9ecb2d5b118c34380b7e7fc8a7c0c3690b73ad (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-660`, `QUEUE_KANBAN_BROWSER` set, removed afterwards)
- do-work-cli-integrations: exit 0, executed, 86 s
- staged-skills: exit 0, executed, 47 s
- updater: exit 0, executed, 74 s
- installer: exit 0, executed, 30 s
- Earlier drain at the first merge `08f005e2` (before the review fix): all four lanes exit 0, executed (73, 39, 66, 29 s).

## Timing

Observed 2026-10-10T13:33:24Z to 2026-10-10T14:37:57Z: 1h 04m 33s total, 1h 09m 39s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 55m 40s | 2 |
| verification-gate | 13m 44s | 4 |
| handback-merge | 15s | 1 |

Slowest stage: review / re-review of the review-fix delta, 29m 07s, outcome success.

Note: no builder-work event was recorded. The hand-back had landed before this integrator started (builder commit 13:30:43Z, dispatch 13:23:49Z, about 7 minutes of build), so per fan-out-reference.md → Landed hand-back the event is skipped rather than charging the builder with the wait. The review row includes the first review (about 26 minutes) and the delta re-review (about 29 minutes).
