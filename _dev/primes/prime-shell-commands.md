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

- [REQ-247: a REQ that states two ceilings must say which wins — the ordering clamp and the introducing-commit ceiling are unsatisfiable together, and nobody caught it before build](../../do-work/archive/REQ-247-archive-timestamp-audit-tool-driven-by-git-commit-times.md#lessons-learned)
- [REQ-250: grep the same primitive across the file before calling a class closed — the genuinely silent instance lived in a different checker's copy of it; and pin every documented limitation with a fixture that can fail](../../do-work/archive/REQ-250-close-the-remaining-markdown-link-checker-gaps.md#lessons-learned)
- [REQ-246: derive a repair-side parser from the read-side detectors it claims parity with — a hand-rolled recognizer covers fewer shapes and can half-rewrite what it half-recognizes](../../do-work/archive/REQ-246-repair-wrong-queue-and-working-timestamps-from-the-session-hook.md#lessons-learned)
- [REQ-244: recognition broad, requirement narrow — a detector that only recognizes the spellings it already fixed locks in nothing](../../do-work/archive/REQ-244-cite-the-timestamp-rule-at-every-stamp-write-site.md#lessons-learned)
- [REQ-243: run the stated RED against pre-change code before building — a "RED" that was already red is the cheapest signal the work is half done](../../do-work/archive/REQ-243-check-that-shipped-markdown-pointers-resolve.md#lessons-learned)
- [REQ-238: every fix of this class trades a staleness risk that has a detector for a broken-link risk that has none](../../do-work/archive/REQ-238-point-present-work-at-the-canonical-independent-bytes-rationale.md#lessons-learned)
- [REQ-234: a derived count that almost reproduces the remembered figure is the trap, not the fix](../../do-work/archive/REQ-234-stop-the-shell-behavior-suite-counting-its-own-cases.md#lessons-learned)
- [REQ-230: a canonicalization suite proves a restatement is absent, never that the pointer replacing it resolves](../../do-work/archive/REQ-230-point-caller-docs-at-the-canonical-publication-rationale.md#lessons-learned)
- [REQ-180: use the tracked filename's exact case in shell test paths](../../do-work/archive/UR-040/REQ-180-contract-suite-justfile-case-mismatch.md#lessons-learned)
- [REQ-186: give identical baseline child invocations one required owner](../../do-work/archive/UR-041/REQ-186-baseline-suite-single-ownership.md#lessons-learned)
- [REQ-187: keep one maintainer command inventory and close aggregate self-test recursion with a fixture-only mode](../../do-work/archive/UR-041/REQ-187-canonical-local-maintainer-gate.md#lessons-learned)
- [REQ-193: lock complete shell-contract predicates so deletion or negation cannot survive a broad regex](../../do-work/archive/UR-043/REQ-193-keep-archived-urs-closed-during-review.md#lessons-learned)
- [REQ-196: distinguish live root-file inputs from intentional filename variants when enforcing tracked casing](../../do-work/archive/UR-041/REQ-196-lowercase-remaining-root-justfile-contract-paths.md#lessons-learned)
- [REQ-190: publish exclusive immutable snapshots before atomically refreshing mutable canonical output](../../do-work/archive/REQ-190-reduce-present-work-to-portfolio-only.md#lessons-learned)
- [REQ-191: keep optional development previews package-local, foreground-bound, and free of readiness guesses](../../do-work/archive/REQ-191-extract-explicit-present-video-action.md#lessons-learned)
- [REQ-192: separate executable command detection from explanatory prohibition prose and mutation-test the complete unsafe family](../../do-work/archive/REQ-192-migrate-presentation-routing-docs-and-contracts.md#lessons-learned)
- [REQ-198: treat directory publication as both a process-tree and filesystem-transaction boundary](../../do-work/archive/REQ-198-publish-generated-directory-only-after-image-success.md#lessons-learned)
- [REQ-199: verify exact destination semantics and independent immutable bytes, not only command success or inode identity](../../do-work/archive/REQ-199-publish-portfolio-snapshot-before-canonical-refresh.md#lessons-learned)
- [REQ-204: a batch owns the process tree it starts, and `mv` onto a directory nests instead of colliding](../../do-work/archive/REQ-204-harden-ai-report-generated-batch-lifecycle.md#lessons-learned)
- [REQ-205: `ln`/`mv` onto a directory nest and exit zero; a hard link is not a copy](../../do-work/archive/REQ-205-make-portfolio-publication-independent-and-exact.md#lessons-learned)
- [REQ-216: macOS bash 3.2 makes empty-array expansion fatal under `set -u`; use `set --` for optional arguments](../../do-work/archive/UR-049/REQ-216-harden-atomic-download-retry-and-credentials.md#lessons-learned)
- [REQ-217: only `git archive` honors `export-ignore`; a mirrored script must probe both sibling depths](../../do-work/archive/UR-049/REQ-217-add-upstream-archive-fetcher-with-git-fallback.md#lessons-learned)
- [REQ-208: pass data needles to grep with `--` — a pattern starting with `-` parses as options](../../do-work/archive/UR-047/REQ-208-deterministic-p50-estimator-script-reference-schema.md#lessons-learned)
- [REQ-225: `ln` refuses an occupied file but nests on an occupied directory; a bare rename's exit status is not proof of publication](../../do-work/archive/REQ-225-state-verified-exact-publication-as-a-shared-condition.md#lessons-learned)
- [REQ-229: a nested-publication probe must sit before whatever the success status authorizes, not merely after the write](../../do-work/archive/REQ-229-verify-the-published-path-in-download-and-screenshot-helpers.md#lessons-learned)
