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
| `../../do-work-knowledge/scripts/lexical-memory-recall.sh` | Query sanitization, lexical ranking, and attribution |
| `../../do-work-knowledge/scripts/install-memory-hooks.sh` | Independent hook merge, verification, and rollback |
| `../../do-work-toolbox/scripts/generate-report-image.sh` | Backend selection, launched-process-tree ownership, verified exact invocation-private publication, and exact opt-in agentic scratch |
| `../../do-work-toolbox/scripts/generate-report-image-batch.sh` | Parallel batch launch, retained per-image statuses, launched-process-tree ownership, and verified all-or-nothing directory publication |
| `../../do-work-toolbox/scripts/publish-portfolio-summary.sh` | Verified single-source canonical refresh and snapshot-first exclusive publication |
| `../../do-work-toolbox/scripts/install-last30days.sh` | Complete-payload validation and verified exact transactional project-local publication/repair |

`tools/install-do-work-suite.sh` is a compatibility launcher over the `do-work-cli` `install-suite` command, which owns the install transaction. It stays self-contained in one respect only: `--print-bootstrap-command` prints a literal heredoc and needs no Go toolchain, because that snippet has to run before anything is installed. Everything else the installer does requires Go 1.25.0 or newer. Atomic REQ reservation remains owned only by the board package's Go tool; it has no shell twin.

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

## Verified exact publication

Whenever a publication's destination could already be occupied, the publishing step verifies the path it actually wrote; the rename's or link's exit status is not that proof. `ln` and `mv` treat a directory standing in the destination's place as a container rather than a collision, so the payload lands inside it under the private staging name, the command exits zero, and the destination is still the directory it was. A caller that reads only the status records a publication that never happened, over a private file abandoned in someone else's directory.

The trigger is that condition and not the identity of any one helper: it holds for a single file and for a whole staged directory, for every publication described below, and for any publication added later. What a helper does *about* a nested write is its own policy and stays in its own section — advancing to the next candidate, failing closed, and discarding only its own stage are each correct answers to the same check.

## Atomic download publication

Never download incrementally into the final path when presence or size is later treated as success. The shipped helper downloads to a private adjacent temporary file, publishes by rename only after curl succeeds, and preserves failures:

```bash
scripts/atomic-download.sh "$source_url" "$target_path"
```

The helper retries transient failures itself (`--retry 3 --retry-delay 2 --retry-max-time 60`), so a rate-limited host — a sustained codeload 429, for instance — does not fail a caller that would have succeeded a moment later. Plain `--retry` has treated 429 as transient since curl 7.51.0; `--retry-all-errors` is deliberately not used because it would raise the required curl version to 7.71 without adding anything here.

Credentials are opt-in. When `GH_TOKEN` or `GITHUB_TOKEN` is non-empty the helper sends `Authorization: Bearer <token>`; absent or empty, the request goes out exactly as it would without them. Callers get both behaviors by using the helper rather than writing their own `curl`.

Cleanup never converts a failed download into success. When review occurs between download and publication, later command blocks must re-derive the deterministic reviewed path and verify it exists; they must not silently download again.

Publication here makes the [Verified exact publication](#verified-exact-publication) check, and so does the screenshot install that shares these mechanics: each verifies the path it actually wrote, removes only its own nested artifact, and exits nonzero with the occupying directory untouched. For the screenshot install that ordering is what protects the staged source, since a staged capture is removed only after a publication that verifiably happened.

## Portfolio summary publication

Invoke `../../do-work-toolbox/scripts/publish-portfolio-summary.sh` with one retained source and the action-selected mode. The helper copies and verifies that source into a private file adjacent to the canonical target, once per output. `--canonical-only` atomically replaces only the canonical file. `--with-snapshot` first publishes an exclusive snapshot from its own verified copy, advances occupied candidates with numeric suffixes, and only then atomically replaces the canonical file from the same bytes.

The two outputs carry identical bytes but never share storage: a snapshot linked to the canonical file would follow every later in-place edit of it. Each publication makes the [Verified exact publication](#verified-exact-publication) check, and this helper's answers to it are that a snapshot candidate occupied by a directory advances to the next suffix, a canonical path occupied by a directory fails closed, and neither leaves a private file nested inside it.

An exclusive snapshot failure leaves the prior canonical unchanged. A later canonical replacement failure leaves the new snapshot published and reports that partial outcome. Existing snapshots are never truncated, replaced, or automatically removed.

## Report image batch publication

Generate a report's images through `../../do-work-toolbox/scripts/generate-report-image-batch.sh <report-directory> <style-brief> <target-name>:<prompt> …` rather than orchestrating the per-image helper yourself. Each pair splits on its first colon, and a target name must be a bare filename because the batch joins it to its own invocation-private staging directory adjacent to `generated/`.

The batch launches one helper per image, retains every PID and status, and waits each one even after an earlier failure — a bare `wait` would discard the mixed statuses that decide which images are current. An image is current only when its own helper status is zero and its staged target is non-empty; failed targets are removed. Publication happens once, as a single same-filesystem rename of the complete verified batch, and only when at least one image is current. `generated/` must be absent both before staging and immediately before the rename, and the rename makes the [Verified exact publication](#verified-exact-publication) check: a nested stage means someone else owns `generated/`, so the batch discards only its own stage, leaves that directory untouched, and exits nonzero.

An all-failed batch is not an error. It removes its exact private directory, prints nothing, and exits zero so the caller falls back to hand-authored diagrams. Success prints the published directory on stdout; every diagnostic goes to stderr so the caller can read the path directly.

The batch owns the process tree it starts. Each helper is launched under job control so it leads its own process group, and that group is signalled only when it verifies as the helper's own and not the batch's — an unverified group degrades to bare-PID signalling, because the only group it could otherwise hit is its own. An interrupted batch terminates, escalates, and reaps everything it launched *before* staging is removed; nothing it started may keep writing into a directory it is about to delete.

## Raw text before shell quoting

Never place raw user or imported text inside a shell command string. Derive a sanitized token as a text operation first, then substitute only that safe token into the command. An apostrophe is enough to break naïve single-quote interpolation; command syntax inside the text can then execute rather than remain data.

## Diff output filtering

Do not use a bare `diff -x NAME` to hide one generated file: `-x` matches file and directory basenames, so a same-named source directory disappears too. Run the diff without that exclusion and filter the emitted artifact path specifically.

## State across command blocks

Prescribed command blocks are independent invocations. Shell variables, traps, and random `mktemp` paths from an earlier block are unavailable later. Every block re-derives deterministic paths and re-validates inherited artifacts. When a value cannot be derived (for example a captured merge endpoint), the action carries the literal in its durable/context record and re-types it; it never expands a variable assumed to survive.
