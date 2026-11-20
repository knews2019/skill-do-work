---
id: REQ-754
title: 'Addendum: second-opinion runs backends in a throwaway Git worktree, so uncommitted files outside the report never block it'
status: completed
route: C
heavy_verified_at: 2026-11-20T00:05:27Z
heavy_verified_revision: d1156bc6d8493ab87bd673cf662bff9ed1006819
review_at: 2026-11-20T00:00:47Z
commit: d1156bc6d8493ab87bd673cf662bff9ed1006819
kb_status: pending
builder_handback_at: 2026-11-19T23:59:01Z
integration_at: 2026-11-19T23:59:01Z
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
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree_test.go"]
claimed_at: 2026-11-19T23:38:44Z
completed_at: 2026-11-20T00:05:27Z
release_at: 2026-11-20T00:05:27Z
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
- [x] **[APPLY]:** Six scoped source, documentation and test files implemented and merged in d1156bc6.
- [x] **[UNIFY]:** Reviewed all six files and the complete diff. Verified cleanup ordering, current-bundle materialization, exclusive import, session checks, test assertions and documentation. Builder formatting, vet, race, compatibility builds and repository checks passed. Integrated canonical gate and all selected heavy lanes passed.
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

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree_test.go` (new) — authorized regression split for per-file test budget

**Acceptance criteria (restated from REQ):**
- [x] Each subprocess backend uses a detached temporary worktree.
- [x] Current bundle bytes include uncommitted additions, edits and deletions.
- [x] Outside worktree changes fail with named paths and preserve real outside files.
- [x] Only allowed new files return, with existing refusal retention behavior.
- [x] Success, refusal, backend failure, deadline and signal exits clean up owned worktrees.
- [x] Real checkout dirt never blocks subprocess or session second opinions.
- [x] Session bundle preservation and backend permission flags remain effective.
- [x] Shipped instructions describe the new behavior.

## Pre-Flight

**Git:** Clean before claim. Current dirty paths are this run’s request and run artifacts.

**Tests baseline:** Focused TestSecondOpinion checks passed. Repository baseline `bash _dev/tests/maintainer-verify.sh` passed, exit 0, 173 seconds. Request-bound advance preflight and green-gate records are satisfied.

**Dependencies:** Existing Go and ShellCheck toolchains passed the native gate.

## Decisions

- D-01 (DECIDE & STATE). Replace the isolated bundle completely to preserve tracked deletions and take separate inode-aware snapshots.
- D-02 (DECIDE & STATE). Use existing rooted exclusive publication and process reaping. Preserve recoverable per-file publication.
- D-03 (DECIDE & STATE). Canonicalize temporary paths and use a separate 30-second cleanup context. Report cleanup errors with a failed run.
- D-04 (DECIDE & STATE). Preserve a finished staging page exclusively even when concurrent real-bundle edits prevent publication.
- D-05 (DECIDE & STATE). Split new tests into a dedicated worktree test file to satisfy the repository per-file runtime budget.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree_test.go` (new)

**What was done:** Subprocess second opinions now probe and generate in temporary detached Git worktrees containing the current report. Validated new files return through exclusive copy-back and the existing publication tail, with finished refusal output retained. Independent bounded cleanup follows process reaping on all handled exits. Session mode permits unrelated outside edits while enforcing report preservation. Instructions and regression tests cover the new behavior.

## Qualification

Request-bound qualification and scope-drift gates passed for `02f3a5e1..d1156bc6`. All six source, test and documentation files were inspected. The two new-file wiring warnings are false positives for Go package compilation and test discovery. `secondOpinionRun` directly calls the worktree helper, and Go discovers the regression functions. P-A-U and scope evidence are complete.

## Testing

The captured RED/GREEN proof is implemented as `TestSecondOpinionWorktreeDirtyReadme`. Before source changes it failed with `outside-bundle changes: README.md`. After implementation both cases pass, proving dirty tracked README publication, detached HEAD, an outside-write finding, unchanged real README bytes and removal of worktree registration and scratch files.

Run from the CLI module, `go test ./internal/toolboxcommands -run '^TestSecondOpinionWorktreeDirtyReadme$' -count=1` changed from exit 1 to exit 0. The full focused command `go test ./internal/toolboxcommands -run '^TestSecondOpinion' -count=1` passed on the merged main tree in 33.026 seconds. Request-bound test-gate records report scope, focused tests and the repository green gate satisfied. Each test file stays under the repository 30-second file budget.

`bash _dev/tests/maintainer-verify.sh` ran directly and unmodified on merged main, exit 0 in 139 seconds. It checked ShellCheck, formatting, contracts, both Go modules' vet and fast tests. Builder race tests, Go 1.24 and Windows cross-builds also passed. The builder ran a state-output-adjusted repository script, and the parent canonical run supplies the integrated repository evidence.

The final retention fix was also test-first. Concurrent real index replacement and asset collision failed the finished-page assertion before the fix and pass after it. New coverage includes exact bundle additions/edits/deletions, empty directories, isolated probe writes, original replacement, unprefixed and symlink assets, exclusive copy-back, session final publication, deadline, cancellation, SIGINT, setup failure and cleanup failure reporting.

The prior REQ-698 pre-existing-dirt refusal expectation was intentionally replaced by the new positive regression. Existing adapters, permission flags, model verification and judge behavior remain covered. Native protocol fixtures and the existing render-judge seam supply behavior evidence. Live paid model backends were not invoked.

The heavy planner selected CLI integrations, staged skills, updater and installer. All four lanes executed successfully at the merged revision, as recorded below.

## Review

Independent review approved `02f3a5e1..d1156bc6` with no findings. Requirements, code, tests and scope scored 100%; acceptance passed for local implementation and integration. The reviewer read all six files and existing publication/process/rooted helpers and checked each requirement against the original UR.

Independent merged-tree acceptance ran dirty README publication, outside-write refusal, session publication, SIGINT and setup cancellation in 2.763 seconds, exit 0. Whitespace validation passed. The restatement sweep found no stale active instructions requiring unrelated work to be committed. Dated historical records retain their original wording. Windows evidence is cross-build/vet, and live model behavior uses protocol fixtures.

## Lessons Learned

**What worked**
- Reusing owned-process reaping and exclusive rooted publication kept cleanup and output recovery consistent with existing behavior.
- Replacing the isolated bundle preserved tracked deletions as well as uncommitted additions and edits.

**What did not work**
- Comparing a macOS temporary path directly with Git registration missed its canonical `/private` spelling. Canonicalization corrected cleanup.
- Refusing changed real bundles before retaining their finished new page lost recoverable output. Exclusive staging retention fixed that case.

**Worth knowing**
- Real and isolated snapshots need separate inode identities. Teardown needs a fresh bounded context after cancellation.
- A cleanup refusal is reported as a failed run with its path. Absolute-path writes and shared Git metadata remain outside worktree isolation.

## Orientation

[MAP CHANGED] The second-opinion subsystem now permits unrelated uncommitted project work. Subprocess probes and generation run in temporary detached worktrees, and only permitted report output returns to the real report. Session reviews enforce report preservation without requiring outside files to be clean. Touched prime links remain valid.

## Heavy Verification Plan

Base `02f3a5e11ba8cd58ae410fedfca50f39ac081f84`, target `d1156bc6d8493ab87bd673cf662bff9ed1006819`.

- `do-work-cli-integrations` — CLI subtree changed.
- `staged-skills` — shipped skill paths changed.
- `updater` — CLI subtree changed.
- `installer` — CLI subtree changed.

Each lane uses `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane-id>`. The typed plan has no uncovered paths and no uncertainty. The runner executes these lanes at the retained target before finalization.

## Heavy Verification Result

Target and execution revision `d1156bc6d8493ab87bd673cf662bff9ed1006819`. All lanes executed freshly, exited 0 and were not skipped.

- CLI integrations passed in 69 seconds.
- Staged skills passed in 34 seconds.
- Updater passed in 71 seconds.
- Installer passed in 32 seconds.

## Review

**Later review and retrospective** | 2026-11-20T01:30:48Z | reviewed revision 4e48c2ef (0.331.2)

**Verdict: Request changes.** Normal worktree isolation works, but the cleanup and outside-change promises have reproduced gaps. This evidence qualifies the earlier approval; that earlier record remains historical.

- F1 [P2] A readable 0555 asset directory is copied with its mode intact, then both Git removal and os.RemoveAll fail. The page can already be linked while the run reports cleanup failure and leaves scratch files. This is a regression introduced by REQ-754. — impact-user-visible → explicitly captured as REQ-755 under UR-188.
- F2 [P2] The reused git status check misses a committed outside README edit and an ignored outside file. Both runs were accepted with outside_paths=[]. These gaps predate REQ-754; the writes remained in the isolated checkout. — impact-user-visible → explicitly captured as REQ-756 under UR-188.

**Evidence:** Existing TestSecondOpinion tests passed in 38.532 seconds. Three focused Go overlay probes failed on the reviewed revision. The cleanup probe passed against pre-change source 02f3a5e1; the other two failed there too. Original probes and results are preserved as UR-188 assets. Paid models and real-browser acceptance were not exercised. Fixture leftovers were cleaned up.

### Retrospective

**What worked:** Detached worktrees solved the unrelated-uncommitted-files blocker. Existing failure, cancellation and signal tests exercised real lifecycle paths. Comparing the same new probes against the old source distinguished a new regression from inherited defects.

**What missed the failures:** Cleanup tests varied exit paths but omitted directory permissions. Outside-change tests trusted Git status semantics without checking ignored files or a backend commit. Broad green gates proved their covered scenarios, so the earlier unqualified approval overstated the evidence.

**Review handling:** The later review initially kept findings in chat and probes in temporary storage. This append and UR-188 preserve them. The two repairs remain separate so their root causes and closure evidence are independently reviewable.

**Change for this repair:** Start with the three demonstrated failures, add directly related controls, and require RED/GREEN evidence in the normal Go harness. Recheck every original promise after integration and verify both filesystem cleanup and Git registration cleanup. Keep isolation claims bounded: unrestricted absolute-path access and shared Git metadata are outside its security guarantees.
