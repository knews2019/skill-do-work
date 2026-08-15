# Prime: Prescribed Shell Commands

> Read this before writing or reviewing any shell block prescribed inside an action file,
> hook, or tool. Prime files are low noise, high value: every trap below was earned in a
> real debugging session, and new lessons in this domain get added here so future sessions
> find them. This file is maintainer-side (`_dev/` is export-ignored) — nothing shipped
> may cite it.

## Prescribed Shell Commands Must Surface What the Steps Consume

Action files are prose that prescribes shell behavior. When a step's logic iterates over the output of a command, the prescribed command must actually emit the items that logic consumes — a mismatch is invisible in the prose and only shows up when run against a real repo. Traps that have already bitten this skill:

- **`git status --porcelain` collapses wholly-untracked directories** into a single `?? dir/` row — it does not list the files inside. Any step that enumerates untracked files per-item (read each, check extension/size/name) must use `git status --porcelain --untracked-files=all` (`-uall`) or `git ls-files --others --exclude-standard`. The latter also drops correctly-ignored paths, so it doubles as the untracked ignore filter.
- **A blanket skip/exclude list applied _before_ a check silently neuters any check meant to fire inside the excluded set.** Scope skip-lists to the noise they actually target (untracked/ignored) and run tracked-file checks outside the exclusion — e.g. a committed `__pycache__/*.pyc` is correct-to-ignore when untracked but is exactly what a "committed build artifact" check should flag.
- **`git show --name-only` prints the commit header and message before the file list** — a message line can pass a filename grep and become a phantom path, and merge commits list no files at all. Use `git diff-tree --no-commit-id --name-only -r -m <commit>` (or `git show --name-only --format=`) when the output is consumed as file paths.
- **`git show` on a merge commit prints a combined diff that is usually empty** — so any consumer reading a REQ's `commit:` hash as a diff source silently sees nothing on worktree-merged work (the reviewer reads an empty diff as an empty REQ). Detect the second parent with `git rev-parse --verify -q '<sha>^2'` — quoted, since `^` is a glob operator in zsh and an escape character in cmd.exe — and use `git show --first-parent -m <sha>` when it succeeds.
- **Ignore patterns with an interior slash are root-anchored, while `git check-ignore` tests cwd-relative paths** — a guard that checks then appends can mismatch from a subdirectory (duplicate appends, path never ignored). Prefix with `**/` when the consumer may run below the repo root. Relatedly, never build `.git/`-internal paths from `--show-toplevel`; use `git rev-parse --git-path <name>` (worktree- and submodule-safe).
- **Never interpolate raw user text inside shell quoting.** A prescribed command like `$(echo '<user-slug>' | tr ...)` breaks on an apostrophe and is a command-injection vector. Derive a sanitized token as a text operation first, then substitute the already-safe value.
- **`diff -x PATTERN` matches basenames of files _and directories_.** Excluding a build artifact by bare name (`-x queue-kanban`) also excludes any same-named directory — silently blinding the diff to an entire source tree. Filter the diff's *output* for the specific artifact path instead (`| grep -v 'tools/queue-kanban/queue-kanban'`), or use a pattern that can only match the file.
- **`curl -o` writes the final path incrementally** — a mid-transfer failure leaves a non-empty partial file, so any presence- or size-gated consumer (`test -s` detect checks) reads the broken download as complete. Prescribe download-to-a-temp-name plus rename-on-success (`curl -o x.download … && mv x.download x || { rm -f x.download; false; }`); `--remove-on-error` needs curl ≥ 7.83, the rename works everywhere. The `; false` is not optional — `rm -f` on an absent path exits 0, so the plain `|| rm -f` form cleans up and then reports the failed download as a success.
- **Shell state does not survive between prescribed command blocks.** An action's steps run as separate shell invocations (often with a user-confirmation gate between them); a variable defined in one block — especially a `mktemp` random path — expands empty in the next, and an agent that "recovers" by re-running the earlier download can silently bypass a review the flow depends on. Blocks must re-derive what they need from deterministic paths and guard-check that inherited artifacts actually exist.

When a review finds a bug in prescribed-command logic, **grep the same primitive across all actions before calling it fixed** — these patterns are usually copy-pasted, so the fix is rarely local. (The first trap above had been copy-pasted into four action files; the audit only flagged one of them.)

The process-tree, complete-directory, current-invocation artifact, and opt-in authority boundaries are recorded together in [`../lessons/validated-runtime-boundaries.md`](../lessons/validated-runtime-boundaries.md). Read it when a shell helper launches descendants, publishes a directory or artifact, or gates a full-host backend.

## Closed Enumerations Go Stale

When a rule applies "whenever X happens" (load a guardrail, honor an enum, keep a guide in sync), state the trigger _condition_ in the rule's canonical home and mark any caller/value list as illustrative, not exhaustive. Hand-enumerated lists silently go stale the moment the set grows — one review traced four independent defects to this pattern (capture's stale domain enum, prompt-injection's five-caller list, the docs-exemption list, security.md's loader claims). When extending a set, grep for every other enumeration of it and update or generalize each one.

## Lessons

- [REQ-180: use the tracked filename's exact case in shell test paths](../../do-work/archive/UR-040/REQ-180-contract-suite-justfile-case-mismatch.md#lessons-learned)
- [REQ-186: give identical baseline child invocations one required owner](../../do-work/archive/REQ-186-baseline-suite-single-ownership.md#lessons-learned)
- [REQ-187: keep one maintainer command inventory and close aggregate self-test recursion with a fixture-only mode](../../do-work/archive/REQ-187-canonical-local-maintainer-gate.md#lessons-learned)
- [REQ-193: lock complete shell-contract predicates so deletion or negation cannot survive a broad regex](../../do-work/archive/UR-043/REQ-193-keep-archived-urs-closed-during-review.md#lessons-learned)
