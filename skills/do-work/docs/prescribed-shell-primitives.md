# Prescribed Shell Primitives

This is the canonical shipped rationale and executable-home contract for deterministic mechanics used across do-work actions. New command-platform mechanics live behind `tools/do-work-cli.sh`; retained scripts remain compatibility/parity surfaces only where stated. Callers keep one line of local intent plus an invocation. A caller-specific gate always wins; this guide owns the shared primitive, not the action's policy.

## Shipped executable homes

| Canonical executable route | Mechanics owned |
|---|---|
| `tools/do-work-cli.sh … show-commit-diff` | Ordinary versus first-parent merge display |
| `tools/do-work-cli.sh … add-local-git-exclude` | Worktree-safe local exclude resolution and idempotent append |
| `tools/do-work-cli.sh … atomic-download` | Private adjacent download and rename-on-success publication |
| `tools/do-work-cli.sh … capture-screenshot` | Unique verified copy, no-clobber link, dispatch-owned cleanup |
| `tools/do-work-cli.sh … run-blocked-check` | GNU timeout selection and isolated stock-Bash process-group timeout/cleanup |
| `tools/do-work-cli.sh … protected-inventory` | Run-level secret quarantine around the existing inventory/association checks |
| `tools/do-work-cli.sh … stage-exact-deletion` | Cached-metadata-only exact deletion staging |
| `tools/do-work-cli.sh … lexical-memory-recall` | Query sanitization, lexical ranking, and attribution |
| `tools/do-work-cli.sh … install-memory-hooks` | Independent hook merge, verification, and rollback |
| `tools/do-work-cli.sh … record-timing-event` | Lifecycle timing: UTC stamping, elapsed derivation, command redaction, and the folded per-request summary |
| `tools/do-work-cli.sh … generate-report-image` | Backend selection, launched-process-tree ownership, verified exact invocation-private publication, and exact opt-in agentic scratch |
| `tools/do-work-cli.sh … generate-report-image-batch` | Parallel batch launch, retained per-image statuses, launched-process-tree ownership, and verified all-or-nothing directory publication |
| `tools/do-work-cli.sh … publish-portfolio-summary` | Verified single-source canonical refresh and snapshot-first exclusive publication |
| `tools/do-work-cli.sh … install-last30days` | Complete-payload validation and verified exact transactional project-local publication/repair |

The table names where each mechanic is **owned**. Where a route also ships a retained `scripts/*.sh` launcher of the same name, that launcher is what this guide's own prose and the shipped actions **invoke**, and the two are not interchangeable: most launchers translate a legacy positional call into the flags the subcommand requires, so the same arguments passed straight to the command are rejected. `scripts/protected-inventory.sh` goes further and sets `DO_WORK_COMPATIBILITY_SHIM=1`, which selects the `<tag>\t<path>` output that `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` parse one row per file from. Change the command's flags or its output shape and you have changed the launcher's contract; fix both together.

`tools/install-do-work-suite.sh` is a compatibility launcher over the `do-work-cli` `install-suite` command, which owns the install transaction. It stays self-contained in one respect only: `--print-bootstrap-command` prints a literal heredoc and needs no Go toolchain, because that snippet has to run before anything is installed. Everything else the installer does requires Go 1.25.0 or newer. Atomic REQ reservation remains owned only by the board package's Go tool; it has no shell twin.

## Lifecycle timing

A step that needs a lifecycle stage or a material external command timed calls `tools/do-work-cli.sh … record-timing-event` for a completed interval, or `tools/do-work-cli.sh … run-timed-command` for a command the suite launches and measures itself. Both own the clock, so no caller derives a timestamp with `date` and no caller invents a second duration format. The wrapper attaches the child's own stdout and stderr to the console rather than a pipe, and exits with the child's status: the child's own code, 128 plus the signal number when a signal killed it, and 127 when the command never launched. That fidelity is what lets it wrap a gate whose exit status is the caller's evidence.

Timing evidence is metadata only. A command reaches the stream as its executable's base name plus an argv token count, never as arguments, so a token or a user-controlled path cannot land in durable evidence. `tools/do-work-cli.sh … fold-timing-summary` turns one run's stream into a single `## Timing` section and deletes the stream; per-test durations stay with the project's own test-duration log rather than being re-derived here.

## Per-file untracked inventory

Any step that inspects untracked paths individually must receive individual files. Plain `git status --porcelain` may collapse a wholly untracked directory into one directory row, hiding every file beneath it.

- Prefer `git ls-files --others --exclude-standard` when only untracked, non-ignored files are needed.
- Use `git status --porcelain --untracked-files=all` (`-uall`) when tracked and untracked states must share one inventory.
- When filenames may contain spaces, quotes, newlines, or rename/copy provenance, consume `git ... -z`; never store NUL-delimited output in command substitution, and consume the second path carried by rename/copy records.
- A check for tracked files that should be ignored is separate: feed tracked paths to `git check-ignore --no-index`. Do not apply an untracked skip list before a tracked-artifact check.

The complete secret-aware inventory and REQ association ship behind `scripts/protected-inventory.sh`. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` invoke the wrapper; [Protected inventory fallbacks](#protected-inventory-fallbacks) below is the single home for its tag legend, its per-tag reading rules, its association semantics, and the by-hand procedure each mode falls back to.

## Protected inventory fallbacks

`scripts/protected-inventory.sh start` prints one `<tag>\t<path>` row per uncommitted file; its `associate` mode prints one `<owner>\t<path>` row per candidate. Both modes resolve the run-level quarantine file with `git rev-parse --git-path <quarantine name>` and exit 2 when that resolution fails, which is what happens outside a Git repository. This section owns what a tag means, how each tag class is read, what association settles, and the by-hand procedure for each mode. It does not own what a caller then does with an `X` or `XD` row: a writing action may stage an `XD` deletion where a read-only action may not, so each caller states its own handling and its own exit-status responses.

Tag legend:

- **M** — modified (a renamed path is tagged M too: read its diff, don't re-read it as new content)
- **A** — added, covering both staged-new and untracked
- **D** — deleted
- **X** — non-deleted excluded path: secret-shaped, secret-derived, or an ambiguous addition beside `X`/`XD`
- **XD** — deleted secret-shaped path

Secret-shaped matching is case-insensitive and applies to the basename only: `.env*` or `*.env`, `*credentials*`, `*.pem`/`*.key`/`*.p12`/`*.pfx`, and `*secret*`. Thus `.ENV`, `AuthCredentials.json`, `private.PEM`, and `UPPER-SECRET.txt` are all secret-shaped, while an ordinary file beneath a directory named `secrets/` is not classified from the directory name alone.

Read each tag class this way, and no other way:

- **Modified files**: Read the `git diff` for each file. Understand what changed and why.
- **New/untracked files**: Read the file contents. Skip binary files (detect by extension: images, compiled assets, archives). For large files (>500 lines), read the first 100 lines and last 50 lines to understand purpose.
- **Deleted files (`D`)**: Note the path and what the file likely was (infer from path and name).
- **Deleted secret-shaped files (`XD`)**: Note only the path and that it is deleted. Do not run `git diff`, `git show`, `git log -p`, or any equivalent that reads, reconstructs, or displays former contents.

If the script is missing or will not run, build the inventory by hand from the complete NUL-delimited output of `git -c status.renames=copies status --porcelain=v1 --untracked-files=all -z`; the command-line copy setting is required even when repository configuration disables rename detection. First classify every path as M/A/D, including each rename or copy record's second NUL-delimited origin path, then lowercase its basename with `tr '[:upper:]' '[:lower:]'` and apply the patterns above. Change a non-deleted secret-shaped tag to X and a deleted one to XD; a secret-shaped rename origin is XD and its destination is X, while a destination copied from a secret-shaped origin is X without an XD source row. Buffer the complete classification before using it. If the finished inventory contains any X or XD and any remaining A, provenance is ambiguous — Git cannot identify a copy when both source and destination are untracked — so change every A to X. Finally add all X paths to the run-level quarantine and overlay that quarantine before any row is consumed. The `-uall` flag is not optional; preserve the [Per-file untracked inventory](#per-file-untracked-inventory) contract above or files beneath a brand-new directory can escape the exclusion scan. That is a secret-leak path, and `../../do-work-toolbox/actions/stray-check.md`'s Red Flags record that it has been hit.

What `associate` settles, so no caller's prose has to:

- **Terminal-success matching honors the Schema Read Contract's aliases**, so `completed`, `completed-with-issues`, and `complete`/`done`/`finished`/`closed` all qualify. Testing only for the literal `completed` drops every remediated-with-issues REQ, and its files then never get associated.
- **In-flight `working/` REQs are included** regardless of status, since a claimed REQ is never terminal.
- **Conflict resolution:** a path claimed by two REQs goes to the one with the latest `completed_at`. An archived REQ outranks an in-flight one.
- **`do-work/` metadata paths are excluded** from association, matching `tools/checks/scope-drift.sh`.
- **Partial matches count.** If 3 out of 5 files in a REQ's Implementation Summary are among the uncommitted files, group all 3 under that REQ.

If the script is missing or will not run, associate by hand: glob `do-work/archive/**/REQ-*.md` and `do-work/working/REQ-*.md`, read each REQ's `status` (accepting every terminal-success alias listed above) and `## Implementation Summary` list, path-match, and tie-break on the latest `completed_at`.

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

Invoke `tools/do-work-cli.sh --repo-root <project-root> publish-portfolio-summary` with one retained source and the action-selected mode. The canonical command copies and verifies that source into a private file adjacent to the canonical target, once per output. `--canonical-only` atomically replaces only the canonical file. `--with-snapshot` first publishes an exclusive snapshot from its own verified copy, advances occupied candidates with numeric suffixes, and only then atomically replaces the canonical file from the same bytes. Missing, failed, or malformed canonical tooling stops the caller; the retained toolbox script is not a fallback.

The two outputs carry identical bytes but never share storage: a snapshot linked to the canonical file would follow every later in-place edit of it. Each publication makes the [Verified exact publication](#verified-exact-publication) check, and the canonical command's answers are that a snapshot candidate occupied by a directory advances to the next suffix, a canonical path occupied by a directory fails closed, and neither leaves a private file nested inside it.

An exclusive snapshot failure leaves the prior canonical unchanged. A later canonical replacement failure leaves the new snapshot published and reports that partial outcome. Existing snapshots are never truncated, replaced, or automatically removed.

## Report image batch publication

Generate a report's images through `tools/do-work-cli.sh --repo-root <project-root> generate-report-image-batch <report-directory> <style-brief> <target-name>:<prompt> …` rather than orchestrating the retained scripts yourself. Each pair splits on its first colon, and a target name must be a bare filename because the canonical command joins it to its own invocation-private staging directory adjacent to `generated/`. Missing, failed, or malformed canonical tooling stops the caller; a retained script is not a fallback.

The batch launches one helper per image, retains every PID and status, and waits each one even after an earlier failure — a bare `wait` would discard the mixed statuses that decide which images are current. An image is current only when its own helper status is zero and its staged target is non-empty; failed targets are removed. Publication happens once, as a single same-filesystem rename of the complete verified batch, and only when at least one image is current. `generated/` must be absent both before staging and immediately before the rename, and the rename makes the [Verified exact publication](#verified-exact-publication) check: a nested stage means someone else owns `generated/`, so the batch discards only its own stage, leaves that directory untouched, and exits nonzero.

An all-failed batch is not an error. It removes its exact private directory and returns a typed successful fallback outcome so the caller uses hand-authored diagrams. A publication success returns the verified directory in the canonical result; callers never infer freshness from stdout emptiness or target presence.

The batch owns the process tree it starts. Each helper is launched under job control so it leads its own process group, and that group is signalled only when it verifies as the helper's own and not the batch's — an unverified group degrades to bare-PID signalling, because the only group it could otherwise hit is its own. An interrupted batch terminates, escalates, and reaps everything it launched *before* staging is removed; nothing it started may keep writing into a directory it is about to delete.

## Raw text before shell quoting

Never place raw user or imported text inside a shell command string. Derive a sanitized token as a text operation first, then substitute only that safe token into the command. An apostrophe is enough to break naïve single-quote interpolation; command syntax inside the text can then execute rather than remain data.

## Diff output filtering

Do not use a bare `diff -x NAME` to hide one generated file: `-x` matches file and directory basenames, so a same-named source directory disappears too. Run the diff without that exclusion and filter the emitted artifact path specifically.

## State across command blocks

Prescribed command blocks are independent invocations. Shell variables, traps, and random `mktemp` paths from an earlier block are unavailable later. Every block re-derives deterministic paths and re-validates inherited artifacts. When a value cannot be derived (for example a captured merge endpoint), the action carries the literal in its durable/context record and re-types it; it never expands a variable assumed to survive.
