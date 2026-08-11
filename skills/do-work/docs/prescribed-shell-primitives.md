# Prescribed Shell Primitives

This is the canonical shipped rationale and fallback contract for shell fragments used across do-work actions. Callers keep the command they execute and one line of local intent, then point here instead of copying the failure-mode explanation. A caller-specific gate always wins; this guide owns the shared primitive, not the action's policy.

## Per-file untracked inventory

Any step that inspects untracked paths individually must receive individual files. Plain `git status --porcelain` may collapse a wholly untracked directory into one directory row, hiding every file beneath it.

- Prefer `git ls-files --others --exclude-standard` when only untracked, non-ignored files are needed.
- Use `git status --porcelain --untracked-files=all` (`-uall`) when tracked and untracked states must share one inventory.
- When filenames may contain spaces, quotes, newlines, or rename/copy provenance, consume `git ... -z`; never store NUL-delimited output in command substitution, and consume the second path carried by rename/copy records.
- A check for tracked files that should be ignored is separate: feed tracked paths to `git check-ignore --no-index`. Do not apply an untracked skip list before a tracked-artifact check.

The complete secret-aware inventory and REQ association already ship as `tools/checks/uncommitted-inventory.sh` and `tools/checks/associate-files.sh`. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` call those scripts; their manual fallback must preserve `-uall`, NUL parsing, rename/copy consumption, secret quarantine, and the scripts' documented exit meanings.

## Merge-aware commit diff

Before treating a REQ's `commit:` hash as a diff source, detect a merge commit with `git rev-parse --verify -q '<commit>^2'` (keep the revision quoted). For a merge, use `git show --first-parent -m <commit>`; for an ordinary commit, use `git show <commit>`. A normal show of a merge produces a combined view that is commonly empty, so it cannot be used as evidence that the REQ changed nothing.

An orchestrated worktree build is different: the owner passes `<pre>..<merge_hash>` and the reader uses that range. Do not rediscover an endpoint from `HEAD`.

## Commit file listing

When command output will be consumed strictly as paths, use:

```bash
git diff-tree --no-commit-id --name-only -r -m <commit>
```

This emits paths for ordinary and merge commits without letting commit-header or message text become phantom filenames. `git show --name-only --format=` is acceptable when a caller genuinely needs show, but the diff-tree form above is the suite default for path-only consumers.

## Local Git ignore

For genuinely transient paths that must stay local to the consuming checkout, append to Git's own exclude file rather than the project's committed `.gitignore`:

```bash
exclude_file="$(git rev-parse --git-path info/exclude 2>/dev/null)" || exclude_file=""
if [ -n "$exclude_file" ]; then
  git check-ignore -q <path> 2>/dev/null || echo '**/<path>' >> "$exclude_file"
fi
```

Use `git rev-parse --git-path`; constructing `<repo>/.git/info/exclude` breaks linked worktrees and submodules. The `**/` prefix keeps an interior-slash pattern and a cwd-relative `git check-ignore` probe aligned from subdirectories. The command makes an untracked path ignorable, not a tracked path safe: a caller whose requirement is “must never be committed” also checks `git ls-files -- <path>` and asks the user to untrack it when necessary. A mere build artifact can skip that tracked-file check.

## Atomic download publication

Never download incrementally into the final path when presence or size is later treated as success. Download to a deterministic sibling temporary name, publish by rename only after curl succeeds, and make cleanup preserve the failure status:

```bash
curl -fsSL -o "$target_path.download" "$source_url" \
  && mv "$target_path.download" "$target_path" \
  || { rm -f "$target_path.download"; false; }
```

The final `false` is required because cleanup can succeed after the download failed. When review occurs between download and publication, later command blocks must re-derive the deterministic temporary path and verify it exists; they must not silently download again.

## Raw text before shell quoting

Never place raw user or imported text inside a shell command string. Derive a sanitized token as a text operation first, then substitute only that safe token into the command. An apostrophe is enough to break naïve single-quote interpolation; command syntax inside the text can then execute rather than remain data.

## Diff output filtering

Do not use a bare `diff -x NAME` to hide one generated file: `-x` matches file and directory basenames, so a same-named source directory disappears too. Run the diff without that exclusion and filter the emitted artifact path specifically.

## State across command blocks

Prescribed command blocks are independent invocations. Shell variables, traps, and random `mktemp` paths from an earlier block are unavailable later. Every block re-derives deterministic paths and re-validates inherited artifacts. When a value cannot be derived (for example a captured merge endpoint), the action carries the literal in its durable/context record and re-types it; it never expands a variable assumed to survive.

