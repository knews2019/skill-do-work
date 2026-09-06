# Prescribed Shell Primitives

This is the canonical shipped rationale and executable-home contract for deterministic mechanics used across do-work actions. New command-platform mechanics live behind `tools/do-work-cli.sh`; retained scripts remain compatibility/parity surfaces only where stated. Callers keep one line of local intent plus an invocation. A caller-specific gate always wins; this guide owns the shared primitive, not the action's policy.

## Shipped executable homes

| Canonical executable route | Mechanics owned |
|---|---|
| `tools/do-work-cli.sh … show-commit-diff` | Ordinary versus first-parent merge display |
| `tools/do-work-cli.sh … add-local-git-exclude` | Worktree-safe local exclude resolution and idempotent append |
| `tools/do-work-cli.sh … atomic-download` | Private adjacent download and rename-on-success publication |
| `tools/do-work-cli.sh … capture-screenshot` | Unique verified copy, no-clobber link, dispatch-owned cleanup |
| `tools/do-work-cli.sh … run-blocked-check` | Owned probe process-group timeout and kill escalation, first-hand launch and timeout evidence, bounded diagnostic identity, and focused-test baseline comparison |
| `tools/do-work-cli.sh … protected-inventory` | Run-level secret quarantine around the existing inventory/association checks |
| `tools/do-work-cli.sh … stage-exact-deletion` | Cached-metadata-only exact deletion staging |
| `tools/do-work-cli.sh … lexical-memory-recall` | Query sanitization, lexical ranking, and attribution |
| `tools/do-work-cli.sh … install-memory-hooks` | Per-event hook gating, order-preserving settings merge, pre-mutation backup, and rename publication |
| `tools/do-work-cli.sh … record-timing-event` | Lifecycle timing: UTC stamping, elapsed derivation or chaining, command redaction, and the append-only per-request stream |
| `tools/do-work-cli.sh … generate-report-image` | Backend selection, launched-process-tree ownership, verified exact invocation-private publication, and exact opt-in agentic scratch |
| `tools/do-work-cli.sh … generate-report-image-batch` | Parallel batch launch, retained per-image statuses, launched-process-tree ownership, and verified all-or-nothing directory publication |
| `tools/do-work-cli.sh … publish-portfolio-summary` | Verified single-source canonical refresh and snapshot-first exclusive publication |
| `tools/do-work-cli.sh … install-last30days` | Complete-payload validation and verified exact transactional project-local publication/repair |

The table names where each mechanic is **owned**. Where a route also ships a retained `scripts/*.sh` launcher of the same name, that launcher is what this guide's own prose and the shipped actions **invoke**, and the two are not interchangeable: most launchers translate a legacy positional call into the flags the subcommand requires, so the same arguments passed straight to the command are rejected. `scripts/protected-inventory.sh` forwards global flags (`--repo-root` and `--format`) ahead of the command token, accepting `--repo-root <root>` (or `--repo-root=<root>`) either before or after the mode token (`start` or `associate`) so the wrapper works from any directory; remaining arguments reach the subcommand as written, so a positional after the mode is rejected with `unknown option`, and a caller spells subcommand flags itself (`--quarantine-name <name>`). What that launcher adds is `DO_WORK_COMPATIBILITY_SHIM=1`, which selects the `<tag>\t<path>` output that `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` parse one row per file from. Change the command's flags or its output shape and you have changed the launcher's contract; fix both together.

`tools/install-do-work-suite.sh` is a compatibility launcher over the `do-work-cli` `install-suite` command, which owns the install transaction. It stays self-contained in one respect only: `--print-bootstrap-command` prints a literal heredoc and needs no Go toolchain, because that snippet has to run before anything is installed. Everything else the installer does requires Go 1.25.0 or newer. Atomic REQ reservation remains owned only by the board package's Go tool; it has no shell twin.

## Lifecycle timing

A step that needs a lifecycle stage or a material external command timed calls `tools/do-work-cli.sh … record-timing-event` for a completed interval, or `tools/do-work-cli.sh … run-timed-command` for a command the suite launches and measures itself. Neither accepts a duration — elapsed seconds are always derived by the command — so no caller invents a second duration format. `run-timed-command` owns both ends of its interval. `record-timing-event` stamps the end itself and, given no start, chains the start to the previous event's end (or to its own end on an empty stream); `--started-at <RFC3339>` pins the start to an instant the caller already held, which is how `actions/work.md` records a delegated builder wait, and the instant it holds comes from the Timestamp rule in `actions/work-reference.md`, whose POSIX floor is `date -u`. The wrapper hands the child the CLI's own stderr handle for both its stdout and its stderr — a file descriptor, not a pipe — which keeps the CLI's stdout free for the rendered result: redirect the CLI's stdout and you capture the result and none of the child's output, which follows the CLI's stderr. It exits with the child's status: the child's own code, 128 plus the signal number when a signal killed it, and 127 when the command never launched. That fidelity is what lets it wrap a gate whose exit status is the caller's evidence.

A command reaches the stream as its executable's base name plus an argv token count, never as arguments, so no argument value lands in durable evidence. That redaction covers the command and nothing else. `--operation` is free caller text: it reaches the stream as written, and the slowest stage's and the slowest command's operation is printed verbatim into the folded `## Timing` section, which stays in the request file after the stream is gone — so keep tokens, secrets and user-controlled paths out of an operation name. `--agent` and `--revision` are free text too and reach the stream, but not the folded section. All three are only stripped of control characters and `|`, collapsed to one line, and cut to 120 characters; that is bounding, not redaction. `tools/do-work-cli.sh … fold-timing-summary` turns one request's stream — a run holds one per REQ — into a single `## Timing` section and deletes that stream, leaving a sibling request's stream in place; per-test durations stay with the project's own test-duration log rather than being re-derived here.

## Per-file untracked inventory

Any step that inspects untracked paths individually must receive individual files. Plain `git status --porcelain` may collapse a wholly untracked directory into one directory row, hiding every file beneath it.

- Prefer `git ls-files --others --exclude-standard` when only untracked, non-ignored files are needed.
- Use `git status --porcelain --untracked-files=all` (`-uall`) when tracked and untracked states must share one inventory.
- When filenames may contain spaces, quotes, newlines, or rename/copy provenance, consume `git ... -z`; never store NUL-delimited output in command substitution, and consume the second path carried by rename/copy records.
- A check for tracked files that should be ignored is separate: feed tracked paths to `git check-ignore --no-index`. Do not apply an untracked skip list before a tracked-artifact check.

The complete secret-aware inventory and REQ association ship behind `scripts/protected-inventory.sh`. `actions/commit.md` and `../../do-work-toolbox/actions/inspect.md` invoke the wrapper; [Protected inventory fallbacks](#protected-inventory-fallbacks) below is the single home for its tag legend, its per-tag reading rules, its association semantics, and the by-hand procedure each mode falls back to.

## Protected inventory fallbacks

`scripts/protected-inventory.sh start` prints one `<tag>\t<path>` row per uncommitted file, except an untracked hidden file under `do-work/`, which it drops as editor or operating-system metadata; its `associate` mode prints one `<owner>\t<path>` row for each candidate a REQ claims that the quarantine does not hold, and no row at all for the rest — so a missing row means either that no REQ claims the path or that an earlier inventory quarantined it, and those are different states. Both modes resolve the run-level quarantine file with `git rev-parse --git-path <quarantine name>` and exit 2 when that resolution fails, which is what happens outside a Git repository. This section owns what a tag means, how each tag class is read, what association settles, and the by-hand procedure for each mode. It does not own what a caller then does with an `X` or `XD` row: a writing action may stage an `XD` deletion where a read-only action may not, so each caller states its own handling and its own exit-status responses.

Tag legend:

- **M** — modified (a renamed path is tagged M too: read its diff, don't re-read it as new content)
- **A** — added, covering both staged-new and untracked
- **D** — deleted
- **X** — non-deleted excluded path: secret-shaped, secret-derived, or an ambiguous addition beside `X`/`XD`
- **XD** — deleted excluded path: secret-shaped, or derived from a secret-shaped origin

Secret-shaped matching is case-insensitive and applies to the basename only: `.env*` or `*.env`, `*credential*`, `*.pem`/`*.key`/`*.p12`/`*.pfx`, and `*secret*`. Thus `.ENV`, `credential.txt`, `AuthCredentials.json`, `private.PEM`, and `UPPER-SECRET.txt` are all secret-shaped, while an ordinary file beneath a directory named `secrets/` is not classified from the directory name alone.

Read each tag class this way, and no other way:

- **Modified files**: Read the `git diff` for each file. Understand what changed and why.
- **New/untracked files**: Read the file contents. Skip binary files (detect by extension: images, compiled assets, archives). For large files (>500 lines), read the first 100 lines and last 50 lines to understand purpose.
- **Deleted files (`D`)**: Note the path and what the file likely was (infer from path and name).
- **Excluded files (`X`)**: Note the path and the tag only. Do not read, diff, stage, or commit the file, and do not relax that because the name looks ordinary — an addition promoted by the ambiguity rule below carries the same prohibition as a secret-shaped one.
- **Deleted excluded files (`XD`)**: Note only the path and that it is deleted. Do not run `git diff`, `git show`, `git log -p`, or any equivalent that reads, reconstructs, or displays former contents.

If the script is missing or will not run, build the inventory by hand from the NUL-delimited output of `git -c status.renames=copies status --porcelain=v1 --untracked-files=all -z`, dropping any untracked path whose basename begins with `.` beneath `do-work/` — that is editor or operating-system metadata, not pipeline state, and the tool drops it before classifying; the command-line copy setting is required even when repository configuration disables rename detection. First classify each record's own path, in the code's own order: any status containing `D` is D, otherwise a `??` or an index-column `A` is A, otherwise M — so a rename or copy destination is M only when its record reports no deletion. Then lowercase its basename with `tr '[:upper:]' '[:lower:]'` and apply the patterns above. Run each rename or copy record's second NUL-delimited origin path through the same patterns; it gets a row of its own only when the record is a rename and the origin is secret-shaped, in which case that row is XD. Change a non-deleted secret-shaped tag to X and a deleted one to XD; a destination copied from a secret-shaped origin is X without an XD source row. Buffer the complete classification before using it. If the finished inventory contains any X or XD and any remaining A, provenance is ambiguous — Git cannot identify a copy when both source and destination are untracked — so change every A to X. Finally write all X paths to the run-level quarantine, replacing whatever an earlier run left there, and overlay that set on the rows before any is consumed — `start` only writes the file, and the union with a later run's X paths belongs to `associate`. The `-uall` flag is not optional; preserve the [Per-file untracked inventory](#per-file-untracked-inventory) contract above or files beneath a brand-new directory can escape the exclusion scan. That is a secret-leak path, and `../../do-work-toolbox/actions/stray-check.md`'s Red Flags record that it has been hit.

What `associate` settles, so no caller's prose has to:

- **Terminal-success matching honors the Schema Read Contract's aliases**, so `completed`, `completed-with-issues`, and `complete`/`done`/`finished`/`closed` all qualify. Testing only for the literal `completed` drops every remediated-with-issues REQ, and its files then never get associated.
- **In-flight `working/` REQs are included** regardless of status, since a claimed REQ is never terminal.
- **Conflict resolution:** a path claimed by two REQs goes to the one with the latest `completed_at`. A REQ with a parseable `completed_at` beats one without, whichever root is read first; when both parse equal, or neither parses, the first claim found stands (`working/` before `archive/`, name order within a root).
- **`do-work/` metadata paths are excluded** from association, matching `tools/checks/scope-drift.sh`.
- **Partial matches count.** If 3 out of 5 files in a REQ's Implementation Summary are among the uncommitted files, group all 3 under that REQ.

If the script is missing or will not run, associate by hand: collect every `REQ-*.md` under `do-work/working/` and `do-work/archive/` at any depth with `find do-work/working do-work/archive -type f -name 'REQ-*.md' 2>/dev/null | LC_ALL=C sort` — `-type f` skips the symlinks the walk skips, the redirect survives a project that has not archived anything yet, and the sort makes the tie-break below reproducible, because `find` returns raw directory order where the walk reads each directory name-sorted; without `shopt -s globstar` a `**` glob degrades to a single `*`, so it requires exactly one intervening directory and misses the files sitting directly in `do-work/archive/` — read each REQ's `status`, counting a `working/` REQ whatever its status says and an `archive/` REQ only on a terminal-success alias listed above, read its `## Implementation Summary` list, path-match, and tie-break on the latest `completed_at`.

## Merge-aware commit diff

Before treating a REQ's `commit:` hash as a diff source, invoke `scripts/show-commit-diff.sh <commit>`. It resolves the argument with `git rev-parse --verify <commit>^{commit}`, tests for a second parent with `git rev-parse --verify -q <commit>^2` — both handed to `git` as argv, with no shell in between to quote for — and emits `git show --first-parent -m` for a merge or ordinary `git show` otherwise. A normal show of a merge produces a combined view that is commonly empty, so it cannot be used as evidence that the REQ changed nothing.

An orchestrated worktree build is different: the owner passes `<pre>..<merge_hash>` and the reader uses that range. Do not rediscover an endpoint from `HEAD`.

## Commit file listing

When command output will be consumed strictly as paths, start from:

```bash
git diff-tree --no-commit-id --name-only -r -m <commit>
```

This prints no commit header and no message, so message text cannot become a phantom filename, and `-m` is what makes it list a merge's paths at all — without it a merge prints nothing. Two inputs need more than the form above: the repository's first commit prints nothing without `--root`, and a filename carrying a double quote, a backslash, a control character or (under the default `core.quotePath`) a non-ASCII byte comes back C-quoted without `-z`; a plain space does not trigger quoting. On a merge, `-m` prints a path once per parent that changed it, so a path changed on both sides appears twice — de-duplicate before counting. `git show --name-only --format=` is acceptable for a commit known to have one parent when a caller genuinely needs `show`, and it lists a root commit's paths without extra flags; never point it at a merge, where it prints the combined diff — nothing for a conflict-free merge, only the conflict-resolved paths otherwise — and a path-only consumer reads either as the whole truth. No shipped file runs the form above as printed; the three Go readers each carry what their own input needs: gate evidence walks every commit after a green gate with `-m -z`, the Git transaction reads back the commit it just made with `--root -z` and no `-m`, and finalization matches candidate commits against a recorded diff digest with neither.

## Local Git ignore

For genuinely transient paths that must stay local to the consuming checkout, invoke the shipped helper rather than editing the project's committed `.gitignore`:

```bash
scripts/add-local-git-exclude.sh <path> '**/<path>'
```

Use `git rev-parse --git-path`; constructing `<repo>/.git/info/exclude` breaks linked worktrees and submodules. The `**/` prefix keeps an interior-slash pattern and a cwd-relative `git check-ignore` probe aligned from subdirectories. The command makes an untracked path ignorable, not a tracked path safe: a caller whose requirement is “must never be committed” also checks `git ls-files -- <path>` and asks the user to untrack it when necessary. A mere build artifact can skip that tracked-file check.

## Verified exact publication

Whenever a publication's destination could already be occupied, the publishing step must be able to say what actually landed there; a rename's or link's exit status alone is not that answer, and shell and syscall disagree about what it means. In hand-written shell, `ln` and `mv` treat a directory standing in the destination's place as a container rather than a collision, so the payload lands inside it under the private staging name, the command exits zero, and the destination is still the directory it was; a block that reads only the status records a publication that never happened, over a private file abandoned in someone else's directory. None of the shipped publications below runs `ln` or `mv`. Each publishes through Go's `Rename` or `Link`, which cannot nest anything: `link(2)` refuses every existing destination, and `rename(2)` refuses a directory — but silently replaces an occupying regular file, so a rename is only as safe as the check made on the occupant before it. One shipped publication then proves the result: `capture-screenshot` links its stage into place, stats both names, and compares them with `os.SameFile`. The three others below make a misplaced write impossible in advance instead — an exclusive `mkdir` or `link` onto a name already checked free, or a rename guarded by a check of what stands in the destination — and each one's section says which.

The trigger is that condition and not the identity of any one helper: it holds for a single file and for a whole staged directory, for every publication described below, and for any publication added later. What a helper does *about* an occupied destination is its own policy and stays in its own section — advancing to the next candidate, failing closed, and discarding only its own stage are each correct answers to the same check.

## Atomic download publication

Never download incrementally into the final path when presence or size is later treated as success. The shipped helper downloads to a private adjacent temporary file, publishes by rename only after curl succeeds, and preserves failures:

```bash
scripts/atomic-download.sh "$source_url" "$target_path"
```

The helper retries transient failures itself (`--retry 3 --retry-delay 2 --retry-max-time 60`), so a rate-limited host — a sustained codeload 429, for instance — does not fail a caller that would have succeeded a moment later. Plain `--retry` has treated 429 as transient since curl 7.51.0; `--retry-all-errors` is deliberately not used because it would raise the required curl version to 7.71 without adding anything here.

Credentials are opt-in. When `GH_TOKEN` or `GITHUB_TOKEN` is non-empty the helper sends `Authorization: Bearer <token>`; absent or empty, the request goes out exactly as it would without them. Callers get both behaviors by using the helper rather than writing their own `curl`.

Cleanup never converts a failed download into success. When review occurs between download and publication, later command blocks must re-derive the deterministic reviewed path and verify it exists; they must not silently download again.

The two publications answer the [Verified exact publication](#verified-exact-publication) check differently, and only one of them verifies. The screenshot install byte-compares its private copy against the source, publishes it with a link that will not clobber, and then confirms the published name is the same file it staged; that ordering is what protects the staged source, since a staged capture is removed only after a publication that verifiably happened. The download verifies nothing afterwards — it stats the published target only to report a byte count, and discards that error — it refuses a target that is already a directory before it fetches, and then trusts the rename, which cannot nest inside an occupying directory because a rename over one fails. Both remove their own private file on any failure, and both leave an occupying directory untouched and exit nonzero.

## Portfolio summary publication

Invoke `tools/do-work-cli.sh --repo-root <project-root> publish-portfolio-summary` with one retained source and the action-selected mode. The canonical command reads that source once, requires it to be a regular file, and writes every output from those same bytes: each output is staged in a private dot-file beside its own destination — the snapshot in the snapshot directory, the canonical file in its own — written, synced and closed, then linked or renamed into place. Nothing reads a stage back; what guards a replacement is that the canonical file is still, by inode and by bytes, the regular file the command read just before staging. `--canonical-only` atomically replaces only the canonical file. `--with-snapshot` first publishes the snapshot with an exclusive link onto a candidate it checked free, advancing an occupied candidate to the next numeric suffix (`-2`, `-3`, …), and only then atomically replaces the canonical file from the same bytes. Missing, failed, or malformed canonical tooling stops the caller; there is no script to fall back to — neither skill ships one for this publication.

The two outputs carry identical bytes but never share storage: a snapshot linked to the canonical file would follow every later in-place edit of it. Each publication answers the [Verified exact publication](#verified-exact-publication) check by making a misplaced write impossible rather than by re-reading what it wrote: an occupied snapshot candidate, directory or not, advances to the next suffix; a canonical path occupied by anything but a regular file fails closed; and neither leaves a private file nested inside the occupant.

An exclusive snapshot failure leaves the prior canonical unchanged. A later canonical replacement failure leaves the new snapshot published and reports that partial outcome. Existing snapshots are never truncated, replaced, or automatically removed.

## Report image batch publication

Generate a report's images through `tools/do-work-cli.sh --repo-root <project-root> generate-report-image-batch <report-directory> <style-brief> <target-name>:<prompt> …` — there is no shipped script for this publication to orchestrate by hand. Each pair splits on its first colon, and a target name must be a bare filename because the canonical command joins it to its own invocation-private staging directory, which it creates under the system temporary directory rather than beside `generated/`. Missing, failed, or malformed canonical tooling stops the caller; there is no script to fall back to.

The batch launches one helper per image concurrently, keeps each helper's own outcome in its own slot, and joins every one of them even after an earlier failure — a single combined status would discard the per-image outcomes that decide which images are current. An image is current only when its own helper status is zero and its staged target is non-empty; failed targets are removed. Publication is all-or-nothing but it is not a rename. The batch claims `generated/` with an exclusive `mkdir`, which is what makes the [Verified exact publication](#verified-exact-publication) check — the directory must be absent both before staging and at that moment, and a directory that already exists means someone else owns it, so the batch discards only its own stage, leaves `generated/` untouched, and exits nonzero. Having claimed it, the batch writes each verified image into it one at a time, and removes the whole claimed directory if any one of those writes fails, so a partly published `generated/` never survives. It publishes only when at least one image is current.

An all-failed batch is not an error. It removes its exact private directory and returns a typed successful fallback outcome so the caller uses hand-authored diagrams. A publication success returns the verified directory in the canonical result; callers never infer freshness from stdout emptiness or target presence.

The batch owns the process tree it starts. Each helper is launched as the leader of its own process group, and that group is signalled only when it verifies as the helper's own and not the batch's — an unverified group degrades to bare-PID signalling, because the only group it could otherwise hit is its own. An interrupted batch terminates, escalates, and reaps everything it launched *before* staging is removed; nothing it started may keep writing into a directory it is about to delete.

## Raw text before shell quoting

Never place raw user or imported text inside a shell command string. Derive a sanitized token as a text operation first, then substitute only that safe token into the command. An apostrophe is enough to break naïve single-quote interpolation; command syntax inside the text can then execute rather than remain data.

## Diff output filtering

Do not use a bare `diff -x NAME` to hide one generated file: `-x` matches file and directory basenames, so a same-named source directory disappears too. Run the diff without that exclusion and filter the emitted artifact path specifically.

## State across command blocks

Prescribed command blocks are independent invocations. Shell variables, traps, and random `mktemp` paths from an earlier block are unavailable later. Every block re-derives deterministic paths and re-validates inherited artifacts. When a value cannot be derived (for example a captured merge endpoint), the action carries the literal in its durable/context record and re-types it; it never expands a variable assumed to survive.
