---
id: REQ-754
title: 'Addendum: second-opinion runs backends in a throwaway Git worktree, so uncommitted files outside the report never block it'
status: claimed
route: C
dispatch_at: 2026-11-19T23:43:26Z
planning_at: 2026-11-19T23:41:11Z
required_lessons: ["_dev/primes/lessons-releases.md"]
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-11-19T23:39:27Z
  basis:
    - Route C
    - 5-file write set
    - 8 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
created_at: 2026-11-19T23:37:12Z
user_request: UR-187
addendum_to: REQ-698
domain: backend
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
claimed_at: 2026-11-19T23:38:44Z
---
# Addendum: Second-Opinion Runs Backends in a Throwaway Git Worktree

## What
Change `ai-report second-opinion` (REQ-698, the command that asks another model to rebuild a report bundle) so an uncommitted file outside the report bundle never stops it. A subprocess backend (codex, gemini, agy, claude) runs in a new, temporary Git worktree; the launcher checks that worktree, copies only the allowed new files back into the real bundle, and always removes the worktree. Session mode drops the outside-bundle Git status check.

## Prior Implementation
REQ-698 (ai-report second-opinion) shipped the subcommand in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (merge `31e7fcdf`). The backend process runs with its working directory set to the real bundle (`secondOpinionCapture`, called from `secondOpinionRun`). `secondOpinionOutside` / `secondOpinionOutsideCheck` run `git status --porcelain=v1 -z --untracked-files=all` on the real checkout and refuse any path outside the bundle, both before and after the backend run, and also in `secondOpinionPrepare` and `secondOpinionFinalize` (session mode). REQ-698 decision D-02 deliberately kept this strict check and refused a status-delta comparison, because a dirty tracked file can change again without its status code changing. `secondOpinionPreserved` separately checks that `index.html` and existing bundle assets are unchanged except the staging page and prefixed new assets. Backend locks today: codex runs with `-s workspace-write`, claude with `--permission-mode acceptEdits`, gemini with `--approval-mode auto_edit`, agy with `--dangerously-skip-permissions`.

## Why
The current check cannot tell a change the user already had from a change the model made, so it requires a clean tree. Any unrelated uncommitted edit (for example a skill file or a gallery page) blocks a report review, and the user has to commit unrelated work first. A fresh worktree starts clean, so any change in it outside the bundle can only come from the run.

## Detailed Requirements
- Subprocess mode: create a new Git worktree for each backend run, on a detached HEAD so the run cannot move a branch.
- Copy the current bundle (including uncommitted bundle files) into the worktree before the run; the backend's working directory is the bundle inside the worktree.
- After the run, check the worktree: any changed or new path outside the bundle fails the run with a finding naming the paths, as today. Inside the bundle, the existing preservation rule applies (only the staging page, its reply and prefixed new assets may be new; `index.html` and existing assets unchanged).
- Copy back into the real bundle only the files the launcher would publish today (page, reply, prefixed new assets). The existing refusal behaviour that keeps a finished page and names it stays, with the kept page landing in the real bundle.
- Remove the worktree after every run on every exit path: success, refusal, backend failure, timeout, and interrupt (Ctrl-C / signal). Nothing the run created stays registered in `git worktree list` or on disk.
- Remove the outside-bundle Git status check on the real checkout from subprocess mode. Uncommitted files outside the bundle in the real checkout are never read as a failure.
- Session mode (`--prepare` / `--finalize`): drop the outside-bundle Git status check entirely. Keep the bundle preservation check.
- Update the shipped instructions that say a run refuses uncommitted changes outside the bundle (for example in `skills/do-work-toolbox/actions/ai-report.md`, the ai-report guide and the second-opinion row fields such as `outside_paths`), so no text still tells the user to commit unrelated work before a run.

## Constraints
- The worktree is not a security boundary: a backend with no lock can still write by absolute path anywhere on disk. Do not claim otherwise in docs or findings.
- Keep the existing backend locks (`-s workspace-write` and the other flags); this REQ does not change them.
- Delete the code the worktree replaces instead of keeping both paths.

## Red-Green Proof
**RED prompt/case:** In a test repository, the user has an uncommitted edit to `README.md` outside the report bundle. A fake codex backend on a temporary PATH runs `ai-report-second-opinion` on the bundle.
**Why RED now:** The run refuses before the backend starts with `outside-bundle changes: README.md`.
**GREEN when:** The run publishes the second-opinion page, `README.md` is byte-identical to before, and `git worktree list` shows no leftover worktree. Second case: a fake backend that writes `../../README.md` from its working directory fails with a finding naming `README.md`, the real `README.md` is untouched, and the worktree is still removed.
**Validation:** User confirmed

## Builder Guidance
The user picked a Git worktree over a temporary copy of only the bundle after hearing the trade-offs: a worktree checks out the whole repository and shares refs and stash with the main checkout, and an uncommitted bundle must be copied in. The detached HEAD requirement addresses the shared-refs risk; other choices (worktree location, how cleanup survives a signal) are the builder's.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 25210 tokens, `slugged: partial` so no targeted form is legal; it matches on the `session-is-the-backend` and `native-backend-identity` families from earlier second-opinion work.

- `_dev/primes/lessons-action-files.md` — 8851 tokens, `slugged: partial`; matches documentation of alternate session/subprocess paths.

## Full Context
See `do-work/user-requests/UR-187/input.md` for complete verbatim input.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Isolate probes and generation in a detached worktree, validate copied bundle, import allowed files exclusively, publish using existing real-bundle tail, then remove the worktree.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: If we are talking about this level of untrustworthiness, wouldn't then it make more sense to create a new work tree, work on that and then check it before merging it in? Also, make sure to clean up any work tree after the work is done.*

## Triage

**Route: C** - Complex

**Reasoning:** The change spans subprocess isolation, bundle import and publication, cancellation cleanup, and the session entry points.

**Planning:** Required

## Plan

Create one owned temporary detached Git worktree per subprocess backend. Materialize the current report snapshot exactly, including uncommitted additions, edits and deletions. Use a fresh isolated snapshot for inode-aware preservation. Run probes and generation there with worktree paths in the prompt and outside-path checks against that worktree. Revalidate the real bundle before importing only allowed new regular files with exclusive publication. Keep the shared publication tail in the real bundle and preserve finished refusal output there. Cleanup follows process reaping and uses an independent bounded context. Remove all real-checkout outside-path checks, including the shared publication tail. Update the action and guide.

Test first with a dirty README and fake backend, then cover outside writes, session mode, copied bundle fidelity, failure/retention, cancellation and signal cleanup, collisions, and existing adapters. Run focused second-opinion tests, native static checks and the repository gate.

*Generated by Plan agent and checked against every detailed requirement.*

## Exploration

The shared publication tail currently checks outside Git status again after the judge passes. Session mode needs that check removed as well as prepare/finalize checks. Preservation compares inode identity, so real and isolated bundles require separate snapshots. Existing rooted exclusive publication and owned-process capture can be reused. Copying only an overlay would resurrect tracked bundle deletions, so materialization must replace the isolated bundle contents. Each new root artifact is allowed only by the existing staging/prefix rule. Publication remains recoverable per file.

The lessons index was consulted at claim. The CLI satellite remains over the shared selection budget and partially slugged. Its relevant sections still load under the prime touch rules. Release lessons fit and were read in full. The action satellite is also partially slugged and over budget (8851 tokens), matching alternate-writer documentation consistency.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify)
- `skills/do-work-toolbox/actions/ai-report.md` (modify)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify)

**Acceptance criteria (restated from REQ):**
- [ ] Each subprocess backend uses a detached temporary worktree.
- [ ] Current bundle bytes include uncommitted additions, edits and deletions.
- [ ] Outside worktree changes fail with named paths and preserve real outside files.
- [ ] Only allowed new files return, with existing refusal retention behavior.
- [ ] Success, refusal, backend failure, deadline and signal exits clean up owned worktrees.
- [ ] Real checkout dirt never blocks subprocess or session second opinions.
- [ ] Session bundle preservation and backend permission flags remain effective.
- [ ] Shipped instructions describe the new behavior.

## Pre-Flight

**Git:** Clean before claim. Current dirty paths are this run’s request and run artifacts.

**Tests baseline:** Focused TestSecondOpinion checks passed. Repository baseline `bash _dev/tests/maintainer-verify.sh` passed, exit 0, 173 seconds. Request-bound advance preflight and green-gate records are satisfied.

**Dependencies:** Existing Go and ShellCheck toolchains passed the native gate.
