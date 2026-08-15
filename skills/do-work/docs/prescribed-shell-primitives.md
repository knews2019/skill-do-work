# Prescribed Shell Primitives

This is the canonical shipped rationale and executable-home contract for shell used across do-work actions. Reusable mechanics live in shipped scripts; callers keep one line of local intent plus an invocation. A caller-specific gate always wins; this guide owns the shared primitive, not the action's policy.

## Shipped executable homes

| Script | Mechanics owned |
|---|---|
| `scripts/show-commit-diff.sh` | Ordinary versus first-parent merge display |
| `scripts/add-local-git-exclude.sh` | Worktree-safe local exclude resolution and idempotent append |
| `scripts/atomic-download.sh` | Private adjacent download and rename-on-success publication |
| `scripts/capture-screenshot.sh` | Unique verified copy, no-clobber link, dispatch-owned cleanup |
| `scripts/run-blocked-check.sh` | GNU timeout selection and isolated stock-Bash process-group timeout/cleanup |
| `scripts/protected-inventory.sh` | Run-level secret quarantine around the existing inventory/association checks |
| `scripts/stage-exact-deletion.sh` | Cached-metadata-only exact deletion staging |
| `../do-work-knowledge/scripts/lexical-memory-recall.sh` | Query sanitization, lexical ranking, and attribution |
| `../do-work-knowledge/scripts/install-memory-hooks.sh` | Independent hook merge, verification, and rollback |
| `../do-work-toolbox/scripts/generate-report-image.sh` | Backend selection, invocation-private adjacent publication, and exact opt-in agentic scratch |
| `../do-work-toolbox/scripts/publish-portfolio-summary.sh` | Verified single-source canonical refresh and snapshot-first exclusive publication |
| `../do-work-toolbox/scripts/install-last30days.sh` | Complete-payload validation and transactional project-local publication/repair |

`tools/install-do-work-suite.sh` remains self-contained because it is the bootstrap that installs these packages. Atomic REQ reservation remains owned only by the board package's Go tool; it has no shell twin.

## Per-file untracked inventory

Any step that inspects untracked paths individually must receive individual files. Plain `git status --porcelain` may collapse a wholly untracked directory into one directory row, hiding every file beneath it.

- Prefer `git ls-files --others --exclude-standard` when only untracked, non-ignored files are needed.
- Use `git status --porcelain --untracked-files=all` (`-uall`) when tracked and untracked states must share one inventory.
- When filenames may contain spaces, quotes, newlines, or rename/copy provenance, consume `git ... -z`; never store NUL-delimited output in command substitution, and consume the second path carried by rename/copy records.
- A check for tracked files that should be ignored is separate: feed tracked paths to `git check-ignore --no-index`. Do not apply an untracked skip list before a tracked-artifact check.

The complete secret-aware inventory and REQ association ship behind `scripts/protected-inventory.sh`, which orchestrates `tools/checks/uncommitted-inventory.sh` and `tools/checks/associate-files.sh` without duplicating their low-level logic. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` invoke the wrapper; their manual fallback must preserve `-uall`, NUL parsing, rename/copy consumption, secret quarantine, and the scripts' documented exit meanings.

## Merge-aware commit diff

Before treating a REQ's `commit:` hash as a diff source, invoke `scripts/show-commit-diff.sh <commit>`. It detects the quoted second parent and emits `git show --first-parent -m` for a merge or ordinary `git show` otherwise. A normal show of a merge produces a combined view that is commonly empty, so it cannot be used as evidence that the REQ changed nothing.

An orchestrated worktree build is different: the owner passes `<pre>..<merge_hash>` and the reader uses that range. Do not rediscover an endpoint from `HEAD`.

## Commit file listing

When command output will be consumed strictly as paths, use:

```bash
git diff-tree --no-commit-id --name-only -r -m <commit>
```

This emits paths for ordinary and merge commits without letting commit-header or message text become phantom filenames. `git show --name-only --format=` is acceptable when a caller genuinely needs show, but the diff-tree form above is the suite default for path-only consumers.

## Local Git ignore

For genuinely transient paths that must stay local to the consuming checkout, invoke the shipped helper rather than editing the project's committed `.gitignore`:

```bash
scripts/add-local-git-exclude.sh <path> '**/<path>'
```

Use `git rev-parse --git-path`; constructing `<repo>/.git/info/exclude` breaks linked worktrees and submodules. The `**/` prefix keeps an interior-slash pattern and a cwd-relative `git check-ignore` probe aligned from subdirectories. The command makes an untracked path ignorable, not a tracked path safe: a caller whose requirement is “must never be committed” also checks `git ls-files -- <path>` and asks the user to untrack it when necessary. A mere build artifact can skip that tracked-file check.

## Atomic download publication

Never download incrementally into the final path when presence or size is later treated as success. The shipped helper downloads to a private adjacent temporary file, publishes by rename only after curl succeeds, and preserves failures:

```bash
scripts/atomic-download.sh "$source_url" "$target_path"
```

Cleanup never converts a failed download into success. When review occurs between download and publication, later command blocks must re-derive the deterministic reviewed path and verify it exists; they must not silently download again.

## Portfolio summary publication

Invoke `../do-work-toolbox/scripts/publish-portfolio-summary.sh` with one retained source and the action-selected mode. The helper copies and verifies that source into a private file adjacent to the canonical target. `--canonical-only` atomically replaces only the canonical file. `--with-snapshot` first hard-links that private verified file to an exclusive snapshot candidate, advances occupied candidates with numeric suffixes, and only then atomically replaces the canonical file from the same bytes.

An exclusive snapshot failure leaves the prior canonical unchanged. A later canonical replacement failure leaves the new snapshot published and reports that partial outcome. Existing snapshots are never truncated, replaced, or automatically removed.

## Raw text before shell quoting

Never place raw user or imported text inside a shell command string. Derive a sanitized token as a text operation first, then substitute only that safe token into the command. An apostrophe is enough to break naïve single-quote interpolation; command syntax inside the text can then execute rather than remain data.

## Diff output filtering

Do not use a bare `diff -x NAME` to hide one generated file: `-x` matches file and directory basenames, so a same-named source directory disappears too. Run the diff without that exclusion and filter the emitted artifact path specifically.

## State across command blocks

Prescribed command blocks are independent invocations. Shell variables, traps, and random `mktemp` paths from an earlier block are unavailable later. Every block re-derives deterministic paths and re-validates inherited artifacts. When a value cannot be derived (for example a captured merge endpoint), the action carries the literal in its durable/context record and re-types it; it never expands a variable assumed to survive.
