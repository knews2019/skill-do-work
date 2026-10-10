# REQ-660 re-review: review-fix delta 08f005e2..bd9ecb2d

Scope: the fix for the first review's F1 (a committed link deletes the main tree's ignored directory on merge) and F2 (wrong `pre` on a remediation re-merge, with no doc caveat). This is a quick scan of 3 files: `internal/cleanup/worktree_lifecycle.go`, `worktree_lifecycle_test.go`, `actions/fan-out-reference.md`. Reviewed at bd9ecb2d.

## Verdict

The delta closes F1 and F2. The guard runs before `git merge`. It refuses with no side effect, and the main tree's directory survives. The normal flow with a correct ignore line produces no warning and still works end to end. No new finding is above negligible.

## Evidence

- Focused tests: `go test -C skills/do-work/tools/do-work-cli -count=1 -run TestWorktree ./internal/cleanup/`. All 7 pass, including the new `TestWorktreeMergeRefusesCommittedLinkBeforeItReplacesTheMainTreePath`.
- F1 fixture, run with a binary built from bd9ecb2d in a `mktemp -d` repo, now removed. Setup: `worktree-links` = `node_modules`, `.gitignore` = `node_modules/`, a real `node_modules/pkg.txt` in the main tree.
  - `worktree new REQ-7`: outcome success, one warning `WORKTREE-LINK-NOT-IGNORED` for `node_modules`.
  - In the builder: `git add -A` and commit. The commit contains `feature.txt` and the `node_modules` symlink. This is the original hazard.
  - `worktree merge REQ-7`: exit 1, outcome refused, `WORKTREE-LINK-COMMITTED`, affected path `node_modules`. HEAD did not move. There is no `MERGE_HEAD`. `node_modules/pkg.txt` still reads `dep` and is still `!!` ignored.
  - Recovery: in the builder, `git rm --cached node_modules` and commit. The next `merge` succeeds, `node_modules/pkg.txt` survives, and `cleanup REQ-7` succeeds. The refusal message ("drop it from the branch first") leads to a path that works.
- Normal flow with ignore line `node_modules` (no slash). `worktree-links` also lists a nested `packages/app/node_modules`. `new` creates both links with no findings. The builder's `git add -A` commits only `feature.txt`. `merge` succeeds and gives a two-parent `[REQ-7] merge builder branch …` commit. `cleanup` succeeds and leaves only the `main` branch. Both main-tree dependency directories are intact.
- `next_argv` of the guard is `git -C <main> diff --name-only <pre>...<name> -- <link paths>`. This command is read-only.
- Test fails without the guard: in a scratch copy of the module with the guard disabled (`if false && …`), the new test fails at line 176 with `merge: exit 0 outcome success findings []`. This matches what the integrator observed.
- Docs: `fan-out-reference.md:53` is accurate. A `dir/` pattern matches only a directory, so it misses the symlink. The sequence "new warns, `git add -A` commits, merge refuses" matches what the fixture showed. `fan-out-reference.md:83` now lists the link refusal and says to keep the first `<pre>`. This agrees with step 1 (line 78) and with *Remediation re-merges* (line 93).

## New findings

- F9. The hand path's step 2 (`fan-out-reference.md:79`) has no matching link check. An integrator who merges by hand with a `dir/` ignore line can still lose the directory. The ignore-line sentence at line 53 covers the root cause, and the command is the documented way to merge. — impact-negligible → report only
- F10. The guard's pathspec is a configured path taken literally from `worktree-links`. This has two side effects. (a) If a listed path is tracked, or is missing from the main tree when `new` runs, `new` skips the link, and a builder's legitimate commit to that path is then refused at merge. (b) A glob character in a line acts as pathspec magic. Neither case fits how `worktree-links` is meant to be used, which is for ignored, untracked paths. — impact-negligible → report only
- F11. The shared test fixture `lifecycleRepository` uses `shared-deps/`, the bad ignore form. So every lifecycle test's `new` now emits `WORKTREE-LINK-NOT-IGNORED`, and the green-path test models the bad config. The tests still pass and the new test depends on this setup, so this is cosmetic. — impact-negligible → report only

### Re-review (delta 08f005e2..bd9ecb2d)

- Overall score: 92%
- Acceptance: Pass
- F1: closed. The F1 fixture refuses with `WORKTREE-LINK-COMMITTED` before `git merge`: HEAD does not move, there is no MERGE_HEAD, and the main `node_modules/` survives. With a guard-disabled copy, the new test fails.
- F2: closed. `fan-out-reference.md:83` says the reported `pre` on a re-merge is this merge's parent and that the first `<pre>` must be kept.
- F9 the hand-path step 2 has no link check, and only the ignore-line sentence protects a hand merge — impact-negligible → report only
- F10 the literal link pathspec refuses a legitimate commit to a tracked or skipped link path, and it treats glob characters as magic — impact-negligible → report only
- F11 the shared lifecycle fixture uses the `shared-deps/` form, so every `new` in the tests warns — impact-negligible → report only
