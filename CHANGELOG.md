# Changelog

What's new, what's better, what's different. Most recent stuff on top.

> Older releases (0.65.x through the entry below the 20th here) live in [`CHANGELOG-archive.md`](./CHANGELOG-archive.md) — kept in the git repo but excluded from the distribution tarball. Tarball-installed copies (no local `.git`, no archive file on disk) can browse it at <https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-archive.md>. Pre-0.65 release notes live one hop further back: `CHANGELOG-2026-spring.md` and `CHANGELOG-pre-0.50.md` at commit [`bf15fe2`](https://github.com/knews2019/skill-do-work/tree/bf15fe2).

---

## 0.162.4 — Reconcile Reservation, Maintainer-Contract, and Token-Canonicalization Docs (2026-08-01)

Codex re-review of the final tree caught three consistency gaps the UR/exclusive-session work left behind:

- **Reservation override.** `actions/reserve.md` and `actions/prime-req-reservation.md` still said "only the default full-queue scan honors reservations" — but a `do-work run UR-NNN` scoped run now honors them too. Corrected to "only explicit per-REQ naming overrides a reservation".
- **Maintainer contract.** `CLAUDE.md`'s board lock-step rule still told maintainers "which REQs may co-dispatch stays with the dispatch gate" — machinery REQ-069 deleted. Reworded to the exclusive-session, display-only framing.
- **Token canonicalization.** The Target ID Resolution contract advertised case-insensitive `req-42` but callers glob/compare against zero-padded upper-case stored ids. Added a canonicalize-before-lookup step (uppercase prefix, numeric-value digit match) so `req-42` actually resolves to `REQ-042`.

## 0.162.3 — Sweep the Last "Overlaps Everything" Mis-Framings (2026-08-01)

Verification found two more spots that inherited the deleted dispatch gate's "absent write_set ⇒ overlaps everything" wording — `actions/capture-reference.md` and a malformed-glob parenthetical in `actions/work-reference.md`. Both corrected so all five board-story files agree: an absent/empty `write_set` gets **no** overlaps badge (unknown, not conflict), and the docs make no false universal claim about malformed globs.

## 0.162.2 — Finish the Board Comment Sweep (2026-08-01)

Re-verification caught one `tools/queue-kanban/model.go` comment (`isWriteSetOverlapCandidateStatus`) still using the removed "a dispatcher could put these in flight together" framing while its four siblings had been refreshed. Reworded it to the exclusive-session, contention-heads-up wording so the file tells one story. Comment-only; the Go tests still pass.

## 0.162.1 — Correct the Board Overlap Semantics After the Exclusive-Session Cut (2026-08-01)

An adversarial verification pass over the batch caught a factual error introduced by the concurrency-machinery removal, plus stale comments. Fixes:

- `actions/work-reference.md` said an absent/empty `write_set` makes the board's overlaps badge "read as overlaps everything." The board does the opposite — it shows **no** badge for an empty set (absence = unknown, not conflict). Corrected the schema note and the crash-recovery rationale to match `actions/board.md`, `tools/queue-kanban/prime-do-kanban.md`, and the board code.
- `actions/work.md` `--wave` restatement now says "any targeting token (`REQ-`/`UR-`)" instead of "targeted REQ IDs", matching the authoritative Input guard.
- Refreshed `tools/queue-kanban/model.go` and `web/board.js` **comments/tooltip** that still described a "dispatch gate" / "co-dispatch" decision the exclusive-session model removed (comment-only; the board's overlap display is unchanged and its Go tests still pass).

## 0.162.0 — REQ IDs Accepted by `do-work roadmap` (2026-08-01)

Closes the inverse of the asymmetry this batch fixed: `roadmap` took `UR-NNN` as a scope token but not `REQ-NNN`, so `do-work roadmap REQ-067` silently returned a whole-queue survey. Now every id-taking action in the skill accepts both prefixes.

- `do-work roadmap REQ-NNN` scopes to a single REQ — its status, dependency position, feasibility read, and its UR siblings for context. Deliberately thin (per-REQ detail stays with `do-work inspect`).
- Multiple id tokens resolve to their union; the soft unrecognized-argument fallback stays, so a genuinely unknown token still yields the full survey with a note.
- Cites the shared Target ID Resolution contract for token shapes rather than restating them.

## 0.161.1 — UR Token Precedence and Reserve Mode-Table Fixes (2026-08-01)

Two corrections from PR review of the UR-id work.

- `do-work reserve`/`release` mode table now lets either `REQ-NNN` or `UR-NNN` lead — a UR-only `do-work reserve UR-011 for cloud-alpha` or `do-work release UR-011` no longer misses every row (or gets read as a free-text label).
- `do-work run` now defines precedence for a REQ reached both by explicit name and by UR expansion (`do-work run UR-011 REQ-068`): explicit naming wins, so the deduped target takes the named branch (bypasses `depends_on`, claims a reservation).

## 0.161.0 — Exclusive-Session Model Replaces the Concurrency Machinery (2026-08-01)

The work pipeline now states plainly what it always assumed: one `do-work` session, one active REQ, one coder context. The ~6,500 words of orchestrator-lock, parallel-dispatch, and co-dispatch-re-validation machinery that existed to detect and recover unsupported concurrent runs are gone, replaced by a short operating rule. Behavior for the normal single-session run is unchanged; the pipeline is just far smaller and easier to follow.

- Removed: the `Concurrent-Orchestrator Lock Guard` section, the Step 1 parallel-dispatch gate and serial-only rule, the Step 3 / Step 5.5 co-dispatch re-validations, the crash-recovery concurrency gate, cleanup's Pass 0 live-claim gate, and every orchestrator-lock heartbeat/claim touchpoint.
- Added a `## Execution Model — Exclusive Session` rule: unexpected repo state matters only when it blocks the active REQ; the coder stops after three consecutive fix attempts; and **read-only actions (roadmap, board, inspect, forensics, recap, reviews) may run in parallel** — the boundary governs writers only.
- Kept: crash recovery of an interrupted single session, single-builder worktree isolation, and the `write_set` field (now display-only, feeding the board's overlaps badge; the board tool is untouched).

## 0.160.0 — UR IDs Accepted by `do-work abandon` and `do-work reserve` (2026-08-01)

The last two REQ-only actions now take a UR. `do-work abandon UR-011` cancels the UR's cancellable members (including any `failed` one still holding it open) behind a single itemized confirmation; `do-work reserve UR-011 for cloud-alpha` reserves its pending members; `do-work release UR-011` returns them. Both cite REQ-067's Target ID Resolution contract instead of re-deriving the rule.

- Each action resolves tokens in front of its existing per-target loop, so every per-REQ gate (status refusals, the `claimed`/`reserved` extra confirmations, reserve's pending-only capture, release's reserved-only touch) applies unchanged to expanded members.
- Bulk cancel is protected by one prompt that itemizes every resolved target with a total count — no `--yes` or per-member bypass.
- `release <token>` precedence is stated explicitly: a `UR-` token resolves as an id, so a reservation *label* literally named `UR-011` is released by naming its REQ ids, not by `release UR-011`.

## 0.159.0 — UR IDs Accepted by `do-work run` (2026-08-01)

`do-work run UR-011` now works: a UR argument expands to its member REQs and runs them in dependency order, instead of being rejected as an unrecognized token. A new shared Target ID Resolution contract defines the token grammar once so the other id-taking actions can cite it rather than each restating it.

- New `Target ID Resolution` contract in `actions/work-reference.md` — `REQ-`/`UR-` token shapes (case-insensitive), UR→REQ expansion by scanning `user_request:` frontmatter (never the `requests:` array), and an empty resolution that stops the action instead of falling through to a full-queue run.
- `do-work run` accepts `UR-NNN` alongside `REQ-NNN`; an explicitly-named REQ bypasses `depends_on`, a UR-expanded member does not. `--wave` stays mutually exclusive with any targeting token.
- The unrecognized-argument guard is unchanged — a bad token like `REG-042` still errors, never a silent full-queue run.

## 0.158.1 — Pre-Flight Tells a Failed Test Run From a Test Command That Never Ran (2026-08-01)

The pre-flight check records a test baseline so the reviewer can tell a pre-existing failure from a regression the builder just introduced. It was recording one even when the test command never launched — a typo'd command exits 127, which looked exactly like a red suite, so a session could start with a baseline that described no test run at all and later excuse a real regression as "already failing before we started". The whole point of the check is attribution, and that was the one way it could get attribution wrong.

- A command that can't be launched (exit 126/127) now reports `WARN: could not run the test command — no baseline recorded`, writes no failures file, clears any stale one, and records `"launched": false`. The review step reads that field and refuses to compare rather than comparing against a fiction.
- `preflight.sh "npx jest --silent"` works now. Passing the test command as one quoted string used to die with an opaque `command not found` — which is how the bug above got triggered in the first place. The quoted form runs via `sh -c`, so `"cd app && npm test"` works too.
- The baseline path is written once instead of spelled two different ways, and the script's header documents that it runs from the project root like its siblings.

## 0.158.0 — An Answer You Give Mid-Run Is Written Into the REQ (2026-08-01)

The question rulebook governed how to ask you something and said nothing about where your answer had to land — so a long run could stop, ask two decisions that were genuinely yours, get detailed answers, and write none of it down. The next builder started fresh, found one question still unanswered, and re-decided it off the stored recommendation whose reasoning you had just rejected. Same outcome, wrong reason, no record.

- New `crew-members/clear-questions.md` Principle 8: an answer obtained interactively outlives the transcript or it may as well not exist — write it into the durable record before acting on it, and capture any new work it implies as its own REQ. Scoped to decisions the work is later read from; a plain proceed/abort gate needs nothing extra.
- `actions/work.md` Step 3.5 gained the missing third branch. It had two — the builder decides (`- [~]` plus a D-XX), or the question waits for `do-work clarify`. Now: a decision escalated mid-run and answered is written in by the orchestrator before dispatch as `- [x]`, never `- [~]` and never with a D-XX, because it is your decision and not the builder's guess.
- `actions/clarify.md` Step 4's `- [x] question → answer` form is now a named entry point, canonical for any caller that obtains an answer — cited by name so a renumbered step can't dangle. Step 5's status flip stays clarify-only; a REQ already in flight has no `pending-answers` to leave.
- Three assertions in `_dev/tests/contract-regressions.sh` pin the principle, the branch, and the named format so a later cleanup pass can't drop one half and leave the others pointing at nothing.

Mid-run questions are deliberately still allowed — asking was the right call. The bug was losing the answer.

## 0.157.1 — Board Pills and Drawer Fields Wrap Instead of Overflowing (2026-07-31)

Long badge values on the Kanban board — most visibly a card's "blocked by \<condition\>" pill — used to run off the card edge because badges were locked to a single line. Pills now wrap onto extra lines, and the drawer's metadata values break long unbroken tokens (write_set paths, blocked_check commands) instead of widening past the drawer.

- `.badge` drops `white-space: nowrap`: wraps within the card, keeps the small-caps label ("BLOCKED BY", "NEEDS") on one piece.
- Drawer `detail-meta` values get `min-width: 0` + `overflow-wrap: anywhere` so a single long path can't stretch the grid.
- Testing-view meta chips wrap the same way.

## 0.157.0 — Source Repo Now Tracks Its Own Queue and Knowledge Base (2026-07-31)

This repo had been keeping its own `do-work/` and `kb/` untracked via a local `.git/info/exclude` entry, on the theory that committing them would leak the maintainer's queue into consumer installs. That theory was already false — `.gitattributes` `export-ignore` plus the tar `--exclude` flags do that job, and the tarball has never contained either directory. Meanwhile the blanket ignore was costing real safety: several of the skill's own guards only work on tracked files, so the repo that dogfoods do-work was running with them silently disabled.

- `do-work/` and `kb/` are now tracked here, matching the Trail of Intent the skill tells consumers to commit. That re-arms `record-commit-hash.sh`'s HEAD size-floor and numstat data-loss guards, and makes `do-work cleanup` Pass 6 blanked-REQ recovery functional instead of a permanent no-op. It also ends a hybrid state where 25 archived REQs were tracked and 54 weren't.
- `/kb export-ignore` added to `.gitattributes`, and `--exclude='kb'` added to all six install/update tar invocations (README, `actions/version.md` ×3, `tools/do-work-update.sh` ×2). Consumers' own `kb/` was already safe — extraction never deleted it — this keeps *upstream's* KB from landing in their skill directory.
- New contract regressions assert `/do-work export-ignore`, `/kb export-ignore`, and the `kb` tar exclude. Those lines are now the only barrier to shipping the maintainer's archive, so they get a ratchet rather than a comment.
- Only genuinely-transient runtime state stays locally excluded: the orchestrator lock and its mutex files, plus preflight's `do-work/working/baseline.json` and `baseline-failures.txt`.
- Caught on the way in: archived `REQ-034` had 833 KB of raw Explore-agent session JSONL pasted into its `## Exploration` section where the agent's summary belonged. Verbatim session capture — UUIDs, prompt text, local paths — has no business in permanent history, so the block is replaced by a note explaining the removal; the REQ's Scope/Verification/Decisions sections are untouched (844 KB → 32 KB).
- Fixed a stale claim in `actions/note.md` that `do-work/pipeline.json` is kept out of git "via the shipped `.gitignore`" — `.gitignore` is itself export-ignored, so no `.gitignore` ever ships. `actions/pipeline.md` Step 4's `.git/info/exclude` entry is the real mechanism.

## 0.156.1 — Verify Pins Removed Lines to the Parent's Frontmatter, Updater Warns About Uncommitted Edits (2026-07-31)

Two follow-ups from an external review of 0.155–0.156. Both were narrow, and both had a way of losing content quietly.

- `record-commit-hash.sh --verify` used to accept *any* removed line starting with `commit:` as "the old field being replaced". Archived REQs that quote the schema have `commit:` lines in their **body**, so a hook that deleted one while the write-back inserted the frontmatter line netted +1/−1 and passed — with a message claiming the patch was a single line. The removal side is now measured against the parent commit's actual frontmatter: zero removals on an insert, exactly that line on a replace.
- The updater's recovery instructions read like a full undo. They aren't: git restores what was **committed**, so an uncommitted edit to a shipped file dies at the extraction now that no rollback copy is kept. It now says that before the confirmation prompt — with the files named and `git stash` suggested — and repeats in the recovery output that the printed `git checkout` won't bring those edits back. The `git clean` line also flags that a root install's shipped paths hold project-owned files.

## 0.156.0 — Update Script Keeps No Rollback Copy (2026-07-31)

`just run-do-work-update` no longer duplicates your whole install before extracting. Version control is the undo, and copying a tracked tree on every run buys nothing git does not already hold — it just left a `do-work.preupdate-<timestamp>.bak` sitting in the project as untracked noise.

- The `cp -R` snapshot and the automatic restore are gone. The prompt now reads "Files are overwritten in place and no rollback copy is kept."
- A failure inside the destructive region reports the partial install and prints recovery commands you can paste — `git checkout --` scoped to the shipped paths git actually tracks, then `git clean -nd --` to review what the extraction added before deleting it.
- If the install is **not** tracked in git (a project that gitignores `.claude/`, or no repo at all), you get told that before the confirmation prompt, because there nothing can be restored.
- New `_dev/tests/update-script-behavior.sh` runs the real updater against a synthetic install with a stubbed upstream fetch: the happy path leaves no `.bak`, a mid-extraction failure reports instead of restoring, and a declined update changes nothing. A contract check fails the build if a `cp -R` snapshot is ever reintroduced.

## 0.155.0 — Commit-Hash Verify Inspects the Committed Patch, Partial Restores Fail Loudly, Updater Flags Stale Files (2026-07-31)

Three fixes from an external review of the data-loss guards shipped in 0.151–0.153. The headline: `--verify` was making a promise it couldn't keep against the commonest kind of pre-commit hook.

- `record-commit-hash.sh --verify` now asserts the **committed patch** — HEAD introduced exactly one new line for the REQ, and that line is the `commit:` field. It used to compare the committed blob against the worktree, which proves nothing when a hook rewrites the file and re-stages it: both sides move together, so a body could be silently gutted while the sizes agreed. Where the patch can't be isolated (root commit, merge HEAD, file added by that same commit, or `--verify` run too late) it now says so and labels the weaker guarantee instead of reporting a clean pass.
- `blanked-req-scan.sh --restore` no longer reports a partial repair as a repair. Content restored but its recorded hash rejected is its own outcome now: exit 1, a `FAIL:` line, and the write-back's own diagnosis passed through rather than swallowed. Previously that file counted as restored and `do-work cleanup` Pass 6 was told everything was fixed — then committed it with provenance pointing at nothing.
- `do-work update` lists files sitting in your install that upstream no longer ships. The extraction only overwrites, so a deleted action or check used to survive downstream forever while the post-update audit reported clean — the old filter dropped "upstream removed this" and "that's your `.DS_Store`" on the same line. Reported, never auto-deleted: your own file in a shipped directory looks identical from there. (The `prompts/`/`interviews/` pre-clean already covered those two directories; this covers the rest.)
- Four new probes in the guard fixture, including a re-staging hook that rewrites the body without changing what a blob read-back can see.

## 0.154.0 — Project Justfile Left Alone, Failed Updates Roll Back, Lock Race Closed (2026-07-31)

Three follow-ups from a review of the update script and the orchestrator lock. The headline: a failed `do-work update` now restores itself instead of leaving you a partial install and a path to fix by hand.

- `do-work update` no longer overwrites a project-owned justfile when the skill lives at the project root. It records which of `justfile`/`Justfile`/`.justfile` the project actually uses — by real directory entry, since a case-insensitive filesystem makes `[ -f justfile ]` match a `Justfile` — and restores that exact name and content across the extraction. A nested install's justfile is the skill's own and still gets refreshed.
- Any failure inside the destructive region now rolls the install back automatically: shipped paths are restored from the rollback copy, files the update added are cleared, and the copy is kept as the audit trail. Previously it printed the backup path and left the partial install in place.
- The lock mutex's remaining lost-update window is closed, not just narrowed. The staged lock image now lives *inside* the mutex directory, so an evicting `rm -rf` destroys it and the evicted owner's publishing rename can only fail — and that rename now fails closed (exit 3, re-acquire) instead of swallowing the error.

## 0.153.1 — Guard Fixture Lints Clean (2026-07-31)

Housekeeping in the commit-hash test fixture, plus one assertion that was written but never wired up.

- The scan probes now assert the *recoverable* byte count reported from the pre-blanking commit, not just that a file was found — the number an operator actually decides on. It was computed and then dropped on the floor.
- `_dev/tests/record-commit-hash-guards.sh` is shellcheck-clean at default severity.

## 0.153.0 — Cleanup Restores Blanked Archived REQs (2026-07-31)

Forensics could tell you an archived REQ's content had been destroyed; now cleanup can put it back. The recovery used to be hand-rolled git archaeology, one file at a time, against a deadline — the lost content only survives until `git gc` collects it.

- New `### Pass 6` in `do-work cleanup`: shows a dry run first (which file, which commit, how many bytes, which hash goes back), asks before writing, and restores only what you approve. Unattended runs report and stop.
- `tools/checks/blanked-req-scan.sh --restore` does the work — temp file plus atomic rename, refuses to write recovered content that's itself empty, and re-applies `commit:` by calling `record-commit-hash.sh` rather than hand-editing frontmatter, so the guards come along.
- `--dry-run` writes nothing and keeps the finding exit code; a completed repair exits 0, because a fixed thing is not a finding.
- A file git has no non-empty version of is reported as a permanent loss, never silently skipped.
- Seven restore probes added to the git fixture suite, including the full six-file incident reproduction and a byte-identity assertion against the pre-blanking blob.

## 0.152.0 — Forensics Detects Blanked REQ Files (2026-07-30)

`do-work forensics` can now tell "this REQ's content was destroyed" apart from "this REQ has a typo in its status." Previously a 0-byte archived REQ showed up as an `unrecognized status ''` warning whose suggested fix — edit the status field — would have written over a file that needed recovering first.

- New check 13 flags any REQ/UR file that is 0 bytes or has lost its frontmatter, as **Critical**, with the recovery commit and the recorded implementation hash already resolved from git history.
- New `tools/checks/blanked-req-scan.sh` does the scanning and the history walk. Read-only, so forensics keeps its never-modifies-anything contract; `--porcelain` emits machine-readable records.
- Check 11 now skips files with no parseable frontmatter, so a destroyed file is reported once, with the remedy that fits.
- A file with no non-empty version anywhere in history is reported as unrecoverable rather than silently passed over.

## 0.151.0 — Guarded Commit-Hash Write-Back (2026-07-30)

The Step 9 "record commit hash" step used to be prose — write the hash into the archived REQ's `commit:` field, then commit. In a repo using do-work, that free-form edit truncated six archived REQ files to 0 bytes, destroying 9 KB to 26 KB of decision trail each, with commit messages that claimed success. It's a script now, and the guards run before anything is staged.

- New `tools/checks/record-commit-hash.sh`: edits only the `commit:` line inside the frontmatter block, and refuses to write unless the rewrite changed exactly that one line. A `commit:` quoted in body prose is structurally unreachable.
- Guards for the real failure shapes — an already-blanked or truncated REQ (checked against its size in `HEAD`), an unterminated frontmatter block, duplicate `commit:` keys, CRLF, a hash the repo can't resolve, and the literal `<hash>` placeholder pasted straight out of the docs.
- Running it twice is a no-op, but an edit that never got committed (a pre-commit hook rejected it) is detected and reported for committing rather than silently skipped.
- A `--verify` mode reads the committed blob back, which is the only way to catch a content-mutating pre-commit hook rewriting the file after every other guard passed.
- Works the same on the worktree-dispatch path, where the hash is the `--no-ff` merge commit; degrades cleanly outside git or where `do-work/` is git-excluded.
- First behavioral test fixture in the repo: `_dev/tests/record-commit-hash-guards.sh` builds a throwaway git repo and asserts each guard actually fires.

## 0.150.15 — Update Script: Guard Project Docs, Clean Non-Interactive Cancel, Rollback Pointer on Failure (2026-07-30)

Four fixes to `do-work update`'s helper script after a code review of the new just shortcut. The headline: running it where the skill *is* the project root (a dev repo or direct clone) no longer risks deleting your project's own `CLAUDE.md`/`AGENTS.md`, and a piped or CI invocation now cancels cleanly instead of dying mid-prompt.

- Non-interactive stdin (piped, CI, `</dev/null`) now defaults to No and exits 0 — previously the bare `read` returned non-zero at EOF and `set -e` aborted before the cancel branch could run.
- The stale-vendored-doc cleanup (`rm` of `CLAUDE.md`/`AGENTS.md`) runs only for a nested install; when `skill_root == project_root` those are the project's own instruction files and are left untouched.
- `justfile` joined `shipped_paths`, so its overwrite now shows in the pre-confirmation diff and the uncommitted-changes warning instead of happening silently off-list.
- A failure mid-update (e.g. ENOSPC during extraction — the `cp -R` backup just doubled the on-disk size) now always prints the rollback-copy path, closing the one path that left a half-updated install with no pointer to the backup.

## 0.150.14 — Update Check No Longer False-Fails on Stray Local Files (2026-07-30)

`do-work update` could report a clean, successful update as broken — and steer you toward a needless rollback — whenever a Finder `.DS_Store`, an editor swapfile, or any other local-only file sat under a shipped path. The post-update integrity check now ignores install-only extras while still catching genuine extraction failures.

- Post-update check drops `--new-file` and filters install-only "extras" with a fixed-string match, so a stray `.DS_Store` or `*.swp` no longer trips the false `exit 1`.
- The filter is metacharacter-proof (`grep -vF`), so a `[`, `+`, or `.bak` in the install path can't defeat it.
- Added an explicit missing-file check so a wholly-absent shipped file or directory is still caught as a real failure — a gap that dropping `--new-file` would otherwise have opened.

## 0.150.13 — Project-Local Just Update Shortcut (2026-07-29)

Projects that install the do-work just recipes can now update their local skill without spending an agent turn.

- `do-work install just-kanban` now adds `just run-do-work-update`, which checks the upstream version, shows the installed-versus-upstream diff, asks before overwriting, creates a rollback copy, and preserves the runtime `do-work/` directory.
- Existing justfile recipes detect the new command as drift and offer the same consent-gated upgrade as other shipped recipe changes.

## 0.150.12 — Board No Longer Flags assets/ Deliverable Copies as Duplicate REQs (2026-07-29)

The Kanban board (`do-work board`) walked every `REQ-*.md` file, including deliverable copies parked under a UR's `assets/` folder. Those attachments have no frontmatter `id`, so their id fell back to the filename and collided with the real ticket — producing a spurious "duplicate REQ id" data warning (or a phantom card, for a uniquely-named asset).

- The walk now prunes any `assets/` folder at any depth, alongside the existing `deliverables/` and `runs/` exclusions.

## 0.150.11 — Abandon Resolves a Failed REQ So Its UR Can Close (2026-07-29)

A `failed` REQ had no way out — nothing in the skill could move it off `failed`, so any User Request holding one stayed open forever (the gap 0.150.10 uncovered). Now `do-work abandon` cancels an already-archived failed REQ in place, flipping it to `cancelled` so its UR can close, while keeping the failure record intact.

- `do-work abandon` accepts a `failed` REQ (at `do-work/archive/` root or `legacy/`) and cancels it in place — no move, and the failure signal (`error`/`error_type` plus a `## Cancelled` "Previously: failed" note) is preserved.
- **The one thing to know:** a completed follow-up REQ recovers the work but never flips the original out of `failed` — cancelling is the only transition, needed whether or not a follow-up ran. This corrects the 0.150.10 note that framed a follow-up as resolving it.
- `cleanup` Pass 1 and `forensics` Check 6 now point you at `do-work abandon` when a failure is holding a UR open; the failure-resolution rule lives canonically in work-reference.md's Terminal-resolved statement, with the three closure readers deferring to it (no restated copies to drift).
- No board or schema change — `cancelled` was already terminally-resolved; the change is purely which inputs abandon accepts and how the resolution is documented.

## 0.150.10 — UR Closure Keys on the Terminal-Resolved Set Everywhere (2026-07-29)

work.md Step 8's archive table was the last reader still counting `failed` as closing a User Request — cleanup Pass 1 and forensics Check 4 already keyed on work-reference.md's terminal-resolved set, so the two halves of the pipeline disagreed on whether a failed REQ holds its UR open. Now all three readers cite the one canonical set: a `failed` REQ keeps its UR open until a follow-up resolves it.

- Step 8's row cites the set instead of restating it; the canonical paragraph's caller list is now marked illustrative and includes forensics Check 4.
- docs/cleanup-guide.md and docs/forensics-guide.md stop describing closure as "all REQs archived" (a failed REQ is archived but not resolved).
- Discovered en route: nothing in the skill can resolve a `failed` REQ at all — queued as a follow-up decision (REQ-060, pending-answers).

## 0.150.9 — Blocked-Flip Guard Judges Worktree Builders by Their Branch (2026-07-29)

Step 8's blocked-vs-failed call used `git diff` on the main tree to ask "did edits land this attempt?" — but a worktree builder commits on its own branch, so the main tree always reads clean and real work got wrongly parked as `blocked`. The guard now reads the builder's branch in worktree mode.

- Case order: a completed hand-back merge (`<merge_hash>` held) proves edits landed; otherwise a quoted `git rev-parse --verify -q '<operative_name>'` existence probe (a missing branch = genuine before-any-work, and `rev-list` on it would exit fatal instead of printing a count); only then `git rev-list --count HEAD..<operative_name>` decides.
- Judged from git state, never the builder's handed-back manifest; serial-mode behavior untouched.

## 0.150.8 — Merge-Aware Diff Reads for Worktree-Merged REQs (2026-07-29)

Every consumer that reads a REQ's `commit:` hash as a diff source now detects the worktree `--no-ff` merge case and diffs against the first parent — plain `git show` on a merge prints a near-empty combined diff, so standalone reviews and receipts of worktree-merged work silently saw nothing.

- Shared idiom at all sites: detect via `git rev-parse --verify -q '<sha>^2'` (quoted — `^` is special in zsh/cmd.exe), then `git show --first-parent -m <sha>`; ordinary serial commits unchanged.
- Covered: review-work.md Get-the-Diff + Two-Modes table, present-work.md (three sites incl. the interactive explainer's receipt), pipeline.md's Completion-Report bullet + pipeline-reference.md's rendering template, ai-report.md's Verify-It-Yourself spec.
- The maintainer shell-trap catalog gains the merge-commit/empty-combined-diff trap.

## 0.150.7 — Lessons-Capture Honors a Prime's Inline-Only Marker (2026-07-29)

The pipeline's Lessons-capture step now inlines a lesson into a prime file that declares itself inline-only, instead of appending an archive link that would be dead in every consumer install.

- Both write sites (`actions/work.md` Step 8, `actions/review-work.md` standalone twin) branch on the prime's `## Lessons` marker comment ("inlined, not linked").
- Keyed off the marker condition in the prime's header, never a hand-list of primes (Closed Enumerations rule).
- The normal (non-marked) link path is unchanged.

## 0.150.6 — Forensics Check 4 Keys UR Closure on user_request Scan (2026-07-29)

`forensics` Check 4 (Orphaned URs) was keying UR closure on the capture-time `requests:` array — the same stale-list bug REQ-048 already fixed in `cleanup` Pass 1. That let it false-positive on UR-007 today, warning to archive a UR that still has six pending follow-up REQs.

- Check 4 now derives UR membership by scanning `user_request` across `queue/`, `working/`, `archive/` root, and `archive/UR-NNN/`, gating on `work-reference.md`'s Terminal-resolved status set by pointer (mirroring the REQ-048 fix), instead of testing whether the `requests:` array's ids all live under `archive/`.
- The `requests:` array is no longer read as the closure predicate; the live UR-007 false positive is gone.

## 0.150.5 — Route A Keeps Its Capture-Seeded write_set (2026-07-29)

Small doc-accuracy fix found during REQ-045: `capture-reference.md` said the pipeline's Scope step "firms up and overwrites" a REQ's `write_set` — true only for Routes B and C, since a Route A REQ never runs that step.

- Notes that a Route A REQ keeps its capture-seeded `write_set` for the whole run, and that value is what `work.md` Step 3 re-validates for disjointness when co-dispatched (per REQ-045).

## 0.150.4 — Board User Guide (2026-07-29)

The board's features were only documented in the agent-facing action file — no linkable human tour existed. New `docs/board-guide.md` covers what a user actually sees: modes, columns, badges, and the Testing view.

- Covers serve/static/summary modes, the four board columns, the Notes and Completion-anomalies strips, the toolbar, the card drawer, and the Testing view's columns and per-card actions.
- The `overlaps` badge gets its own subsection on the four ways it can under-report (no `write_set` declared, `*`/`**` glob quirks, identical malformed patterns, directory entries never badging files inside them) so it isn't misread as a safety guarantee.
- Keeps the human-tour/agent-contract boundary: feature facts sourced from `actions/board.md`, but none of its build/dispatch internals are duplicated here.

## 0.150.3 — Display-Only Overlap-Annotation Invariant Ratchet (2026-07-29)

The board's write-set overlap badge is display-only by design — it must never affect column placement — but that invariant was only protected by one Go test plus prose. `contract-regressions.sh` now pins it on both sides.

- Ratchet asserts `annotateWriteSetOverlap` runs after `bucketColumns` in `model.go`'s `buildBoard`, and that `board.md`'s Rules block keeps its "display-only, never column logic" wording.
- Anchors are call-site and heading-scoped (not file-wide greps), so a hoisted call or a relocated doc claim fails loud with a message naming the file and the fix.
- Red-green verified across five mutation scenarios in a sandboxed clone; live tree stays green.

## 0.150.2 — Board Badge Render-Path Test Assertion (2026-07-29)

The overlap badge's frontend render path had zero test coverage — only the Go-side annotation logic was tested. `generate_test.go` now asserts the badge's render tokens actually make it into the generated board HTML, so a regression in that path fails loudly instead of shipping silent.

- New test anchors on the inlined `web/board.js`/`web/board.css` tokens (`badge-write-overlap`, `writeSetOverlaps`, the drawer row, the CSS rule) rather than rendered DOM, so it holds regardless of live queue contents.
- Red-green verified: mutating the anchor token fails the test; restoring it passes.

## 0.150.1 — Doc-Accuracy Fixes: Legacy-Suppression Comment and Board Glob Miss-Classes (2026-07-29)

Three documentation inaccuracies from the deep review, all comment/prose with no behavior change.

- The memory hook's legacy-suppression comment claimed it "self-clears as soon as the next capture is written" — it actually suppresses to end-of-file and clears at the next UTC day's fresh log; the comment now matches the awk.
- The board's "a malformed glob pattern matches nothing" overstated: `writeSetPatternsIntersect` short-circuits on literal equality first, so two REQs declaring the identical malformed pattern still badge. Aligned across `model.go` (source), `board.md`, `work-reference.md`, and the board prime.
- Added the directory-entry case (`actions/` never badges `actions/board.md`) to the board's illustrative miss-class list, kept explicitly illustrative.

## 0.150.0 — Review Restatement Sweep (2026-07-29)

Six concurrency-spec REQs passed adversarial review at 86 to 98 percent, yet a later pass found every top defect was the same class: a contract changed in its canonical home while a restatement elsewhere kept the old meaning. The review step now forces a sweep for exactly that.

- `review-work.md` Step 6 gains a required **Restatement Sweep**: when a diff redefines something other text restates (a token, a field's semantics, what a hash holds, a command's output shape), grep every other statement/consumer and flag stale ones as findings — including in files outside the REQ's Scope (routed to follow-ups, not scored as builder scope drift). `work.md` Step 7 cross-references it.
- Trigger is condition-based, not a token list (Closed Enumerations Go Stale); a proportionality guard skips diffs that redefine nothing. Inherited by both pipeline and standalone `do-work review`.

## 0.149.3 — Cleanup Keys UR Closure on user_request, Not the Stale requests Array (2026-07-29)

`cleanup` Pass 1 decided a UR was done by reading its capture-time `requests:` array — but review-spawned and addendum follow-ups carry `user_request:` without ever being added to that array, so a UR with pending follow-ups could be archived out from under them. Pass 1 now uses the same predicate `work` Step 8 does.

- `cleanup` Pass 1 derives UR membership by scanning `user_request:` across `queue/`, `working/`, `archive/` root, and `archive/UR-NNN/`, gating on the terminal-resolved set (with `failed` holding the UR open); the `requests:` array is now a report-only cross-check.
- `capture` documents `requests:` as the capture-time record only — never the closure predicate. Two more readers of the old predicate were found and queued as follow-ups (forensics Check 4; a `failed`-status contradiction in work Step 8).

## 0.149.2 — Lock Mutex Re-Verifies Ownership Before Publishing (2026-07-29)

The serialized-lock mutex could evict a slow-but-live owner on the one-minute age check, and that owner's already-staged lock write would still land — clobbering its successor and losing a claim, the exact failure the mutex exists to prevent.

- The prescribed block now re-checks the mutex owner token immediately before the publishing `mv`; on a mismatch it discards the staged temp file, writes nothing, and re-acquires (exit 3) — narrowing the lost-update window from model-round-trip scale to the instant before the rename.
- The mtime-reclaim comment now says the age check proves age, not death, and points at the re-check as what makes eviction safe. The one-minute bound and the fixed-mtime property are unchanged; serial and single-session runs behave identically.

## 0.149.1 — Worktree Cleanup Uses the Recorded Operative Name (2026-07-29)

REQ-038 taught a crash-recovered worktree REQ to re-dispatch under a fresh name variant, but every later step still re-derived the original slug-based name — so after a collision the merge and cleanup targeted the *leftover*, not the builder's actual worktree.

- The name `git worktree add` actually succeeded with is now the REQ's held **operative name**, used by the hand-back merge, Step 8 cleanup, and the own-session crash-sweep exemption — never re-derived from the slug. (The merge site was the sharper bug: it would have integrated the wrong branch.)
- One variant scheme (incrementing numeric token) replaces the free counter-or-timestamp choice; the `worktree-agent-REQ-NNN-` prefix invariant is preserved. No-collision and serial behavior are unchanged.

## 0.149.0 — Dispatch Re-Validation: Full Route Coverage and the Serialization Loser (2026-07-29)

REQ-036's write-set re-validation was written against the Route B/C pipeline only, leaving a co-dispatched Route A REQ building under an unvalidated hint. This states one covering invariant and fixes three coherence gaps around it.

- Every co-dispatched REQ now gets exactly one post-dispatch re-validation, and its route picks which: Routes B/C at Step 5.5, Route A at Step 3 (serialize-only — it has no `## Scope` to hold a partition). The three previously-contradicting sentences now agree.
- The serialization "loser" is defined (the REQ at the re-check is held, never a dispatched sibling mid-build), with a deadlock guard for the two-discoverer case; dispatch-time partitions are persisted into `write_set` so a sibling's re-check compares against the real subset.
- The absent-`write_set` gloss is reworded condition-first and reconciled with REQ-044's conditional recovery clear; a new contract-regression ratchet pins the whole contract.

## 0.148.0 — Lock Claim Coherence: Dispatch Record as the Recompute Source (2026-07-29)

Four coherence defects survived REQ-035's move to a canonical `claimed_reqs` list. The fix names a single source of truth for the recompute — the orchestrator's in-memory dispatch record — instead of a `working/` listing.

- The heartbeat recompute reads the session's dispatch record, never a directory listing and never the lock's own previous `claimed_reqs`. Step 2's claim-before-move append enters the id into that record before the file moves, so the refresh carrying the claim no longer erases it.
- A known-dead builder's id leaves the record (so its REQ is reclaimable that session), explicitly distinguished from an ordinary failed build that keeps its claim through remediation and Step 8.
- Crash Recovery clears `write_set` only for Scope-mirrored sets and preserves capture-seeded / Route-A sets; a stale proceed-anyway gate restatement is aligned and pinned by a new block-scoped contract-regression ratchet.

## 0.147.0 — Worktree Merge Range: Fail-Loud Validation and Seam-in-Range Merge (2026-07-29)

The worktree-dispatch merge range `<pre>..<merge_hash>` had six confirmed defects in how it was produced, validated, and restated — most dangerously, `qualify.sh` printed OK on a broken range instead of failing. All six are fixed as one coherence contract.

- `qualify.sh` now hard-FAILs (exit 1, naming the range) on an unresolvable `DO_WORK_DIFF_RANGE` instead of reading an empty diff and passing vacuously; serial mode is byte-unchanged.
- The integration seam is folded into the merge commit (`git merge --no-ff --no-commit` → apply seam → commit), so it provably lands inside the range and `commit:` still records the merge.
- Remediation re-merges get a defined cumulative range (`<pre1>..<merge_hash2>`); Step 6 gains an imperative orchestrator-side hand-back merge instruction; the hash-writeback block gains a worktree carve-out so it can't record the changelog commit's hash.

## 0.146.3 — Board Overlap Badge Uses OS-Independent Glob Matching (2026-07-29)

The Kanban board's write-set overlap badge matched globs with `filepath.Match`, whose separator is `\` on Windows — so `*` could wrongly cross `/` and the badge would misjudge contention off-platform. It now uses `path.Match` (correct for slash-separated repo-relative paths), and the glob dialect is finally written down where readers meet the field.

- `writeSetPatternsIntersect` uses `path.Match`; the doc comment, `actions/board.md`, the board prime, and the `write_set` schema line all state the dialect: `*` never crosses `/`, `**` is not recursive, malformed patterns match nothing on the board — the dispatch gate still treats an unexpandable/overlapping glob as overlapping, so a board false-negative never loosens it.
- New tests pin the slash boundary (with a same-segment positive control) and malformed-pattern behavior; the badge stays display-only (no schema or column-placement change).

## 0.146.2 — Note Worktree Isolation in the Harness-Tier Guide (2026-07-29)

Someone sizing a harness against `background-agents.md`'s three fan-out rungs had no way to learn that per-builder git-worktree isolation exists or where it's documented. A short cross-reference now closes that gap.

- Added a "worktree isolation is a separate axis" note after the harness rungs, pointing at `actions/work-reference.md` → Worktree Dispatch Mode; the rungs and the file's load contract are unchanged.

## 0.146.1 — Worktree Name-Collision Handling on Re-Dispatch (2026-07-29)

Worktree dispatch mode names each builder's worktree and branch deterministically from the REQ id, and the crash sweep reports (never deletes) an unmerged leftover — so a crash-recovered REQ would re-dispatch straight into the name its own leftover still holds and fail to start. Re-dispatch now sidesteps the occupied name instead of deadlocking on it.

- On a name collision at creation, dispatch under a fresh unique variant (an incrementing `-2`/`-3` or short timestamp token), keeping the `worktree-agent-REQ-NNN-` prefix so sweeps still correlate both names to the REQ.
- The crash sweep now states that a reported unmerged leftover doesn't block re-dispatch — the collision variant covers it, and the two coexist until cleanup Pass 5 resolves the leftover.

## 0.146.0 — Worktree Merge Placement and Evidence Re-Pointing (2026-07-29)

Worktree dispatch mode said who merges and how, but never *when* in the pipeline — and after a merge the main tree is clean, so the qualify check, the review step, and the commit step all read an empty diff and quietly passed nothing. The merge now has a fixed place in the sequence and a defined range, and every evidence step reads that range instead of the post-merge-clean tree.

- The orchestrator merges each builder branch at hand-back (end of Step 6, before the Implementation Summary), and captures the range `<pre>..<merge_hash>` around it — stated once in Worktree Dispatch Mode and consumed by qualify (Step 6.3), review (Step 7), post-merge verification (Step 8), and Step 9's validation.
- `tools/checks/qualify.sh` gained an optional `DO_WORK_DIFF_RANGE` env var; unset, it reads the working+staged diff exactly as before, so serial runs are byte-for-byte unchanged.
- Step 9 is reconciled for merged work: it stages only the changelog/version/metadata (the implementation is already in the merge commit) and records the merge commit's hash.
- Fixed a latent gap where `work.md` and `work-reference.md` disagreed on the post-merge verification default (now both say per-merge whenever more than one REQ is in flight).

## 0.145.0 — Re-Validate Write-Set Disjointness When Scope Firms It (2026-07-29)

The parallel-dispatch gate decided co-dispatch on capture's write-set guess, but Step 5.5 then rewrote that field from each REQ's real scope with no second look — so two REQs seeded as disjoint could both quietly claim the same file once their scopes firmed. Step 5.5 now re-checks disjointness before it commits the field, and a dispatch-time partition directive survives the mirror instead of being erased.

- New Step 5.5 re-validation runs only under co-dispatch: it re-checks the firmed scope against every other in-flight REQ's current `write_set` and serializes or partitions the loser before its builder starts — the same check Step 6 already runs before a mid-build write-set extension.
- The Step 1 gate and Step 5.5 now agree that both steps enforce, at different times (the gate on capture's hint, Step 5.5 after firming); Step 4's plan-validation flag is documented as a warning, not the enforcement point.
- The Step 6 write-boundary bullet clarifies that an absent `write_set` handed to a builder means "dispatched serially, full-scope freedom," never "write nothing."
- Serial (floor) runs are completely unaffected — every new clause is gated to the parallel-dispatch path.

## 0.144.0 — Concurrent Claim Tracking in the Orchestrator Lock (2026-07-29)

The orchestrator lock could only name one in-flight REQ per session, so the moment a single orchestrator dispatched more than one builder at once, Crash Recovery and cleanup mistook the siblings for abandoned crash artifacts and re-queued them mid-build. The lock now tracks every concurrent claim, and the recovery and cleanup gates honor the whole set — the session's own claims included — so parallel dispatch is finally safe inside the skill's own protocol.

- New canonical `claimed_reqs` list on the holder and each coexisting-session entry; the old `claimed_req` stays as a derived legacy mirror (`claimed_reqs[0]`), so older readers and the serial default are completely untouched.
- Crash Recovery now gates on freshness alone and skips any file in a live claim set — the session's own included — so a Step 10 → Step 1 loop no longer strips and re-queues a still-building sibling.
- Per-merge post-merge verification becomes the default whenever more than one REQ is in flight, and `cleanup.md`'s Pass 0 live-claim gate reads the whole claim list.
- Contract-regression ratchets pin both the field's presence and the same-story gate phrasing across `work.md`, `work-reference.md`, and `cleanup.md`.

## 0.143.0 — Capture Slicing Nudge and Board Write-Set Overlap Badges (2026-07-29)

Two upstream levers for parallel-friendly queues: capture's slicing convention now prefers boundaries that give each REQ its own files (declaring unavoidable overlap in `write_set`), and the Kanban board shows an `overlaps` badge on pending/claimed cards whose declared write-sets could touch the same files.

- The overlap annotation computes in Go after column bucketing — structurally display-only (badge + drawer rows, never column logic; co-dispatch decisions stay with the work pipeline's gate) — and is glob-aware in both directions.
- Drawer gains "Write set" and linked "Overlapping write sets" rows; no badge on a REQ without a declared `write_set` (unknown, not safe — the gate's serialize reading is documented alongside).
- The three stale "no overlap computation on the board" claims (board action, maintainer doc, board prime) updated in lock-step.

## 0.142.0 — Worktree Dispatch Mode with Defined Cleanup Ownership (2026-07-29)

The work pipeline now documents running builders in orchestrator-created git worktrees: each builder commits on its own `worktree-agent-REQ-NNN-*` branch, the orchestrator stays the sole writer of the main tree and merges in dependency order, and nothing archives until the merged state re-passes the REQ's checks. Every leftover now has an owner.

- Happy path: the archive step removes the worktree and deletes the branch with `git branch -d` from the integration branch — a refusal is the signal a merge was skipped or lost, so never `-D`/`--force`.
- Crash path: Crash Recovery sweeps `worktree-agent-*` leftovers — merged ones removed mechanically, unmerged ones only reported; discarding unmerged work belongs to cleanup's new consent-gated Pass 5 (its first interactive pass — six passes now, mirrored in the cleanup guide).
- `do-work/` state stays in the main tree only; builders get their brief in the dispatch prompt (and treat any committed `do-work/` snapshot in a worktree as absent).
- Ships honestly single-builder: co-dispatching several worktree builders waits on the lock's multi-claim work (queued follow-up), and four contract-regression ratchets pin the naming, the `-d` assertion, the post-merge gate, and the consent-gated pass.

## 0.141.0 — Write-Set Declarations and a Parallel-Dispatch Gate (2026-07-29)

The queue schema gains an optional `write_set` field (repo-relative paths/globs a REQ expects to write; absent means it overlaps everything), and the work pipeline gains an opt-in dispatch gate: advanced harnesses may co-dispatch dependency-ready REQs whose write-sets are pairwise disjoint. The serial default is untouched.

- Serial-only resource classes: REQs writing ordered/generated resources (migrations, lockfiles, generated bundles — illustrative list) never co-dispatch, disjoint or not.
- Builders get a write boundary: out-of-set needs are a stop-and-report to the orchestrator, never a silent write; the Scope declaration one-directionally mirrors into `write_set` so the two can't drift.
- Chosen over timed per-file locks — a TTL expires over a live slow agent and hands the file to a second writer (the 0.140.4 mutex-break defect class).
- `tools/queue-kanban` parses `write_set` into the board payload (display only), and three new contract-regression ratchets pin the gate text and the parser lock-step.

## 0.140.4 — Owner-Checked Lock Mutex, Atomic Capture Appends (2026-07-28)

Two accepted findings from an external concurrency review. The lock mutex could be forcibly broken after 15 seconds even when its owner was live mid-write — dangerous now that the critical section legitimately spans a model round-trip — and concurrent Stop-hook captures could interleave their writes in the shared daily log.

- `actions/work-reference.md`: removed the 15-second attempt-count mutex break — the one-minute mtime check (a verified stale-owner bound) is now the only reclaim path. The winner records an owner token in the mutex so release can't delete a successor's mutex, and an `mkdir` failure with no contender present reports and stops instead of spinning forever.
- `hooks/memory-stop-capture.sh` (+ spec in `actions/memory-reference.md`): each capture section is composed first and appended in a single `printf` — one atomic `O_APPEND` write — so near-simultaneous stops from sessions sharing the daily log can no longer garble section structure. No lock: `flock` doesn't exist on macOS and the hook must never block session end.

## 0.140.3 — Loopback-Only Board Writes, In-Mutex Stale Revalidation, Redact-Before-Truncate (2026-07-28)

Five accepted findings from an external review, all verified against the code before fixing. The two that mattered most: a LAN-exposed board's testing endpoints accepted writes from any machine (the Origin check only fires when a browser sends one), and a stale-lock takeover judged staleness on a pre-mutex read, so a holder that heartbeated in the gap could be overwritten and its in-flight REQ re-queued.

- Kanban testing writes (`/api/testing/profile`, `/api/testing/status`) now require a loopback peer, same as `/file` — a network-exposed board is read-only.
- The stale-lock takeover re-confirms the holder's identity and recomputes its heartbeat age from the fresh read inside the mutex before overwriting; the user-gated take-over keeps the identity check, and coexisting-session prune ages come from the same fresh read. Wording now ratcheted by `_dev/tests/contract-regressions.sh`.
- The memory Stop hook redacts credentials (and judges the private-key drop) on the full extracted messages *before* truncation — a byte-budget cut can no longer sever a token into an unmatched, persisted fragment. Spec in `actions/memory-reference.md` reordered to match; regression probe reproduces the straddling-token case.
- Stray (misplaced) REQ files now feed the live board's mtime fingerprint, so their warning appears and clears without waiting for an unrelated file change.
- Testing updates preserve the REQ file's existing permission bits across the atomic rewrite instead of forcing 0644.

The remaining four accepted findings from the same external triage 0.140.1 started on. The big one: a live Route C build could go 45+ minutes between lock touchpoints and get taken over — and re-queued — mid-build.

- Heartbeat refreshes are no longer phase-boundaries-only: refresh before dispatching and when each long-running agent returns (explore, plan, build, review), plus at any pause once 15 minutes have passed since the last lock write (`actions/work-reference.md`, `actions/work.md` Steps 4–6). The 45-minute stale threshold's rationale now states the schedule it depends on.
- Cleanup's Pass 0 gained a live-claim gate: any REQ freshly claimed in the orchestrator lock by another session is exempt from the queue and `working/` sweeps — a coexisting session flips its REQ terminal *before* moving it, and that window used to look like an abandoned sweep candidate (`actions/cleanup.md`).
- Next-step suggestions are sourced from SKILL.md's routing triggers, not Action Dispatch names — `do-work capture-requests` didn't route; the capture form is `do-work capture-request: <text>` (`next-steps.md`).
- Lock acquisition and pipeline-state setup now run the `git ls-files` tracked-path check behind their ignore appends, reporting the `git rm --cached` remedy when an earlier session committed the file — an ignore rule can't rescue an already-indexed path (`actions/work-reference.md`, `actions/pipeline.md`; the memory-module installer already had this check since 0.139.2).

## 0.140.1 — Claim-Before-Move Closes the Lock Races; Checkpoints and Temp Files Stop Lingering (2026-07-28)

Four accepted findings from external feedback triage, all in the concurrent-orchestrator machinery. The claim/recovery race and the acquisition race could each let two sessions fight over one REQ; the other two were slow leaks.

- Step 2 now claims the REQ in the orchestrator lock *before* moving it into `working/`, and the Crash Recovery gate lists `working/` *before* reading the lock — the paired ordering means a file can never be observed unclaimed by a live scan (`actions/work.md`, `actions/work-reference.md`). Archive-time gets the explicit mirror order: move out first, then clear `claimed_req`.
- Lock acquisition re-validates existence inside the serialized mutex — a session that loses the empty-queue race now falls into the existing-lock decision tree instead of silently overwriting the winner's holder slot (`actions/work-reference.md`).
- The local-ignore prescription for the lock now uses the glob `do-work/orchestrator-lock.json.*`, covering orphaned PID-suffixed `.tmp` files as well as the mutex directory.
- Checkpoint deletion is scoped to the files a session is allowed to recover, so a coexisting session's live claim in `working/` no longer keeps a stale `CHECKPOINT.md` alive forever (`actions/work.md`).

## 0.140.0 — UTF-8-Safe Capture Truncation; Forget Now Scrubs the Daily Logs Too (2026-07-28)

Two accepted findings from external feedback triage. The stop-capture hook could tear a multi-byte character in half, and "forget" only forgot half the store.

- The stop-capture hook's byte-budget cuts now pipe through `iconv -c -f UTF-8 -t UTF-8` (plain cut when iconv is absent). A raw `head -c` cut lands mid-character routinely on CJK text (~3 bytes/char), persisting invalid bytes into the log and the dedup hash — and on macOS the torn sequence made BSD sed fail in the redaction pipeline, silently dropping the whole capture.
- `memory forget` is now an explicit, confirmation-gated sub-command (with a `do-work forget` alias) instead of one clause inside `remember`. It removes the working-memory bullet AND redacts matching daily-log lines in place with a `[forgotten — …]` marker, since recall searches the logs too. It is the one named exception to the logs-are-append-only rule, scoped to explicit user invocation — automatic writers still only append, capture-body lines keep their `> ` quoting, and heading lines (the dedup key) are never touched.

## 0.139.4 — Capture Boundaries Can't Be Spoofed; Tracked Logs Caught Before the Install Fast Path (2026-07-28)

Four follow-up fixes, all found by reviewing the previous two rounds' fixes. Two of them were holes the earlier fixes left open rather than new bugs.

- Capture sections now open with a sentinel and quote every body line, and the session-start filter ends a section only at an unquoted heading. Heading grammar alone was still spoofable: raw capture text can contain `## 12:34 UTC note`, which ended the section and injected the rest. Legacy sections written before this release have unquoted bodies where no boundary is trustworthy, so they're suppressed to end-of-file — that hides curated entries written after them that day, and self-clears with the next capture.
- The memory-module install detects a tracked raw store in Phase 1 instead of Phase 4. A fully-wired install with a committed log plus an ignore rule took the "already installed" early return, so the check added in 0.139.2 never ran in the one scenario it was written for.
- A `memory` payload that arrives with no sub-command now falls back to `recall`, never `remember`. Real recall queries are usually noun phrases, so the previous sentence-shape test would have classified most reads as writes and mutated the store on a request to read it.
- `do-work recall` with an empty query — including the `what do you remember` phrasing — is now defined: it presents working memory plus recent curated log entries instead of searching for nothing.

## 0.139.3 — Orchestrator Lock Updates Are Serialized; Long Prompts No Longer Eat the Agent Reply (2026-07-28)

Three fixes from a second review pass. The lock guard turns out to have had a lost-update hole since coexisting sessions were added — measured at 19 of 20 concurrent writes discarded.

- Every write to `do-work/orchestrator-lock.json` now runs inside a `mkdir` mutex and lands via temp-file-plus-rename. Multiple sessions write that file (a coexisting session refreshing its entry, the holder's prune, a take-over), so "each writer touches only its own fields" was never enough: a lost `claimed_req` makes Crash Recovery re-queue a REQ that's actively being built. The mutex self-heals after a minute and never blocks the pipeline.
- The Stop hook budgets the user and assistant sides of a capture separately. Truncating the combined string meant a long prompt silently dropped the entire assistant reply — the half holding the decisions. Each side is guaranteed half the budget and a short side yields its slack; cuts are marked `[truncated]`. This was reachable only because 0.139.2 started populating the user side.
- `do-work run` holds the orchestrator lock through `cleanup` and its commit instead of releasing first. Cleanup sweeps `do-work/queue/` and `do-work/working/`, so releasing early let a departing session sweep and commit an arriving session's just-claimed REQ.

## 0.139.2 — Raw Session Captures Stay Out of Session Start; Tool-Using Sessions Capture the Human Prompt (2026-07-28)

Four review fixes to the memory engine, its installer, and the router. The important one: a Markdown heading inside a captured exchange could end the capture section early, so the rest of the raw transcript was injected as curated memory at the next session start — exactly what that filter exists to prevent.

- Capture bodies are now blockquoted at write time and the session-start filter only treats `## HH:MM UTC …` as a section boundary, so a captured `## Findings` can't end the section. The reader-side rule is what protects logs already on disk.
- The Stop hook now finds the real final exchange. Claude Code records tool results as `type: "user"` entries, so any session whose last turn used a tool was storing an empty `User:` side — 7 of 8 recent transcripts on the author's machine. Text is pulled only from `text` blocks, `isMeta` entries are skipped, and blank entries are dropped; the same fix repairs assistant turns that end in a tool call.
- `do-work remember X` and `do-work recall X` keep their verb as the sub-command instead of arriving as a bare payload, and `what do you remember` maps to `recall`.
- The memory-module install verifies with `git ls-files`, not just `git check-ignore`. Ignore rules don't apply to tracked files, so a repair install over a committed log reported "logs ignored: OK" while the Stop hook kept appending to it. The canonical local-ignore snippet in `crew-members/background-agents.md` now carries the same caveat so the next copy-paste doesn't reintroduce it.

## 0.139.1 — Bootstrap Sentinel Stays Machine-Local; "Remember To" Routes to Capture (2026-07-27)

Three review fixes to the memory engine and the installer. The big one: the bootstrap sentinel was committable, so one machine's import would have silently blocked every other clone from importing its own history.

- `do-work install memory-module` now adds `memory/.bootstrap-imported` to `.git/info/exclude` alongside `memory/logs/` and the usage ledger, and verifies it. `memory bootstrap` refuses to re-run when the sentinel exists — committed, that refusal would have followed the repo to every other machine.
- Routing: `remember to fix X` is queued work, not a fact. SKILL.md row 37 now sends task-shaped `remember` phrasings to capture; `actions/memory.md` already documented the boundary, but the router decides the route before the action file is ever read.
- The `bowser` skill download's cleanup returns failure (`|| { rm -f …; false; }`) instead of reporting a failed download as a success — matching the `ui-design` and `ideation-adhd` install rows.

## 0.139.0 — Parallel Memory Engine: memory module, install target, and value auditor (2026-07-25)

A second, capture-first memory engine now runs alongside bkb so real usage data — not theory — decides which one earns its keep (ADR-017). Both engines log usage, and a new auditor renders the head-to-head verdict.

- New `do-work memory` action (`remember` / `recall` / `status` / `bootstrap` / `audit`): a 2,500-char capped `memory/working-memory.md`, dated daily logs, and layered recall (lexical always; semantic only when an embedding backend is detected) with cited sources. Companion schemas in `actions/memory-reference.md`.
- New `do-work install memory-module` target scaffolds `memory/` and merges optional SessionStart/Stop hooks into `.claude/settings.json` — composing with existing hook entries (backup + parse-verify), never clobbering them. The Stop hook appends a hash-deduplicated capture of each session's final exchange and never blocks a session end.
- New `do-work memory audit` (`actions/memory-value.md`): read-only, engine-agnostic value audit of bkb and the memory engine — git/log history probes, usage-ledger stats, hit-cited rate as the verdict signal, with an explicit fairness rule for bkb's pre-instrumentation era.
- bkb `query`/`ingest` now append best-effort usage-ledger events so the comparison sees both engines.
- Review hardening (PR #122): the Stop hook redacts credential-shaped text before persisting captures; session start injects only curated memory (raw captures stay behind `memory recall`, which loads the prompt-injection guardrail); the hook merge checks and appends each settings entry independently; the installer repairs any missing scaffold component; and the auditor locates the KB via bkb's locating contract instead of assuming `kb/`.
- Raw captures and the usage ledger are machine-local: the installer adds `memory/logs/` and `memory/usage-ledger.jsonl` to the repo's `.git/info/exclude` (never your committable `.gitignore`), so only the curated `working-memory.md` is shareable. Redaction stays as a second line of defense.

## 0.138.3 — Skill Downloads Are Atomic: Temp File, Then Rename (2026-07-25)

The 0.138.1 fix caught zero-byte downloads, but `curl -o` writes the final path incrementally — a connection dropped mid-transfer left a non-empty partial `SKILL.md` that `test -s` read as a complete install, unrepairable by re-running. Downloads now land in a `SKILL.md.download` temp name and only rename into place on success.

- Applies to all three curl-based skill installs (`ui-design`, `ideation-adhd`, `bowser`); the temp file is removed on failure, so nothing is left behind either way.
- Chosen over `curl --remove-on-error`, which needs curl ≥ 7.83 — the rename works on any curl.

## 0.138.2 — Install Target Renamed to ideation-adhd (2026-07-25)

The 0.138.0 target ships as `install ideation-adhd` — the name now says what it does (the "adhd" is the upstream skill's metaphor for its branching style, not the substance). `install adhd` and the `adhd-mode` spellings still work as aliases.

- The install **folder** stays `.claude/skills/adhd/` — it must match the upstream frontmatter `name:` field so `/adhd` auto-discovers.

## 0.138.1 — Install Detect Treats a Zero-Byte Skill File as Absent (2026-07-25)

An interrupted download could leave a zero-byte `SKILL.md` that the `ls`-based detect read as "already installed", making the failed install unrepairable by re-running. Review caught it on the new `adhd` target; the same copy-pasted primitive was fixed in `ui-design` and `bowser` too.

- Detect commands for the single-file skill targets now use `test -s` (non-empty), so a re-run over a failed download repairs it instead of stopping at Phase 1.
- The never-overwrite rule is scoped to non-empty files: reinstalling over a zero-byte file is repair, not overwrite.

## 0.138.0 — New Install Target: adhd Divergent-Ideation Skill (2026-07-25)

`do-work install adhd` vendors the [adhd skill](https://github.com/UditAkhourii/adhd) (MIT) into the project — parallel divergent ideation that branches a named problem across distinct cognitive frames, then scores, clusters, and deepens the top candidates. Complements `scan-ideas` (repo-grounded) with deliberately unconventional exploration; feed the winners to `capture-request:`.

- Single self-contained `SKILL.md` installed project-scoped to `.claude/skills/adhd/` — folder name matches upstream so `/adhd` auto-discovers; no global npm install.
- Same manifest-driven detect → install → verify → report shape as `ui-design`; idempotent, never overwrites an existing copy.
- Routing accepts `install adhd` (also the `install adhd-mode` / `install adhd mode` / `setup adhd` spellings — the target normalizes after the install verb; bare `adhd` without the verb is not a route).

## 0.137.0 — Clarify Opens Each Question With a Plain-Language Story (2026-07-24)

Answering pending questions used to mean remembering what REQ-025 was about — often days after the work happened. Now every question arrives with a short story above it: what you asked for, what the builder ran into, and why the call is yours. The decision block underneath is unchanged.

- Questions are presented in three layers — a 1–4 sentence story, then the existing `Decision / Value / Risk / Also` block, then the builder's original wording and file paths *only if you ask*.
- Layer one is written to be read aloud: no file paths, no bare identifiers, no CamelCase, one idea per sentence. Any technical term used lower down gets paraphrased in the story first.
- Blocked REQs waiting on an external condition now get a one-line "what it was for" too — those are the ones you've had the longest to forget.
- New red flags catch the failure mode this invites: a story that just restates the question is padding, not context.

## 0.136.1 — Board Flags REQ Files Found Outside the Scanned Sections (2026-07-24)

A REQ that lands somewhere other than `queue/`, `working/`, or `archive/` — say a work agent that archived to `do-work/user-requests/UR-NNN/` instead of `do-work/archive/` — used to vanish from the board with no trace. Now the walk catches it and raises a data warning instead of silently dropping it.

- The board now emits a warning naming the misplaced REQ, its exact path, and how to fix it (move into `archive/` or `queue/`) — shown in the web warnings banner and the `board summary` output.
- A stray REQ is only flagged, never rendered as a card, so its off-vocab location can't masquerade as a real column placement.

## 0.136.0 — Maintainer Docs No Longer Ship to Consumer Installs (2026-07-23)

The repo's own `CLAUDE.md` and `AGENTS.md` were landing in every consumer install, where Claude Code auto-loads the nested `CLAUDE.md` on every skill-file read — a ~2.5k-word context tax whose commit protocol (bump the version, add a changelog entry) is actively wrong advice inside someone else's project. They're maintainer docs, not skill content, so they no longer ship.

- Both files are now `export-ignore`'d, and `do-work update` deletes the stale copies that installs ≤0.135.x left behind (tar extraction never removes files dropped upstream).
- Every shipped file's citation of the skill's own CLAUDE.md was reworded to be self-contained or point at a shipped home (e.g. `actions/kb-lessons-handoff.md` for the KB handoff contract) — 14 sites across actions, crew-members, and hooks.
- New contract-regression checks keep it that way: the export-ignore lines must exist, and shipped files must not cite the unshipped docs.

## 0.135.0 — Board Drawer Links Every REQ/UR Mention, URL, and File Path (2026-07-23)

The detail drawer's cross-references are now real, obviously-styled links instead of plain text or button-shaped chips. File paths get existence-checked at build time, so a stale reference is visible at a glance.

- Every REQ/UR id in the drawer is a link: the UR drawer's "REQ ids" row, the REQ drawer's "User request" / "Depends on" / "Unblocks" / "Blocked by" rows, and any `REQ-…`/`UR-…` mention inside a rendered body (only ids actually on the board — unknown mentions stay text). Short mentions resolve compound card ids (`REQ-031` → `UR-002-REQ-031`).
- All links are visibly links: accent color + underline (the old "User request" chip looked like a badge).
- File paths in code spans are checked against the repo at board-build time: existing files render as blue links that open read-only via the live server's new `GET /file` endpoint (loopback-only, repo-contained, always text/plain); missing files render red with a "Not found in this repository" tooltip — in static snapshots too, where the existence verdict is baked into the data.
- URLs in code spans become clickable, and every http(s) link in a body opens in a new tab instead of navigating the board away.

## 0.134.0 — Pending-Card Timer Tracks the Last Transition, Not Capture Time (2026-07-23)

A REQ answered via `do-work clarify` went back to `pending` but its board card kept counting from capture time ("queued … · 28m" seconds after the flip). Pending-tier cards now time from the last state change.

- New optional `status_changed_at` frontmatter field, stamped on any status flip that has no dedicated `*_at` stamp of its own — clarify's answered/unblock/discovered-task flips, work Step 1's probe unblock, reserve release, manual resets. Schema + board parser updated in lock-step.
- The board's state timer for pending-tier cards resolves: `status_changed_at` → the later of `created_at` / file mtime (only when mtime beats capture by >5min, so untouched cards still read "queued") — verb "updated" whenever the instant isn't capture time.
- Dedicated stamps stay authoritative: mtime never drives the claim/blocked/reserved stopwatches (the pipeline edits the file all through a claim), and mtime remains banned from completion dating.
- `status_changed_at` is covered by the 0.133.0 future-stamp guard.

## 0.133.0 — UTC Stamping Rule + Board Defense Against Future-Dated Timestamps (2026-07-22)

A session stamped `claimed_at` with local wall-clock time plus a `Z` suffix, and the board's stopwatch froze at "0s" until the wall clock caught up. Now the instructions say exactly how to stamp, and the board calls out bad stamps instead of rendering them silently.

- New Timestamp rule in `actions/work-reference.md` (Full Frontmatter): every `*_at` field is the current UTC instant from `date -u +%Y-%m-%dT%H:%M:%SZ` — never local time with a `Z` suffix. Every write site that said just `<timestamp>` / `<now>` (work claim + terminal flips + blocked flip, capture templates, clarify, reserve, abandon) now says so or points at the rule.
- queue-kanban flags any frontmatter timestamp parsing later than board generation time + 2min skew: a "⚠ future stamp" card badge naming the field(s) plus a data warning (banner in serve/static, listed in summary) with the fix.
- The stopwatch renders an honest "⚠ clock skew" marker (with a tooltip explaining the likely cause) instead of a dead-looking "0s" when its instant is beyond the skew allowance, and recovers to normal ticking once the clock catches up.
- New forensics check 12 sweeps all REQ frontmatter for future-dated timestamps.

## 0.132.0 — Every Non-Terminal Card Gets a Live State Timer (2026-07-22)

The claimed-card stopwatch from 0.131.0 now covers every state: any card that isn't done shows when it entered its current state plus a ticking elapsed timer — "queued … · 3d 04h", "blocked … · 2h 15m", "reserved … · 12m 30s" — so you can see at a glance where time is going before a task lands in Done.

- Pending / pending-answers / failed cards count from `created_at` (labeled "queued" — time since capture; those states write no transition instant of their own), blocked cards from `blocked_at`, reserved from `reserved_at`, claimed from `claimed_at` as before.
- Durations past a day render as "3d 04h" instead of a wall of hours.
- Drawer parity: "Blocked since" and a new "Reserved" row tick while that hold is live, degrading to the plain instant when the field is a stale leftover.

---

## 0.131.1 — Kill Boards by Executable, Bound Probes Everywhere, Guard a Symlinked do-work (2026-07-21)

Four review follow-ups, each pinned by a new contract-regression check.

- `just run-kanban` (and the installed recipe) now identifies a stale board by the listener's **executable path** via `lsof -d txt`, not by grepping argv — a process that merely mentions "queue-kanban" in an argument can no longer be killed by mistake. No `lsof` info still means refuse-to-kill.
- The work pipeline's `blocked_check` probe keeps its 30-second bound even where GNU `timeout`/`gtimeout` don't exist (stock macOS): a background-and-poll fallback kills a hung probe and reports exit 124 instead of running unbounded.
- Roadmap's **Ready** bucket now requires normalized `status: pending`, so `blocked`/`reserved`/`claimed` REQs can't classify as Ready.
- The board's testing writes (REQ placeholders *and* `testers.md`) verify that `do-work/` itself resolves inside the repository — a symlinked `do-work/` pointing elsewhere is refused. If you deliberately symlink `do-work/` outside your repo, the board's Testing view will decline to write there (the work pipeline is unaffected).

## 0.131.0 — Claimed Cards Show Claim Time and a Ticking Duration (2026-07-21)

A claimed card used to sit in its column with no time at all. It now shows when it was claimed plus a live stopwatch — "claimed Jul 21, 19:44 UTC · 4m 01s" — so you can see at a glance how long the current REQ has been in flight.

- The duration ticks every second (s → "Xm YYs" → "Xh YYm") with tabular digits so the line doesn't jitter.
- The detail drawer gets a matching "Claimed" row: ticking stopwatch while the claim is live, plain instant + relative label if a stale `claimed_at` lingers on any other status.

## 0.130.4 — Testing View: Safer Writes, Sturdier Feedback Form (2026-07-21)

Hardening for the board's one write surface and the client flows that drive it.

- Server: all testing writes are serialized behind a mutex (concurrent add-tester or status posts raced their read-modify-write cycles); REQ frontmatter updates land via atomic temp-file-and-rename so a crash can never leave a zero-byte REQ; the testers file opens with O_APPEND.
- Client: typed feedback survives testing-view re-renders and failed posts (the form now closes only on server confirmation); Clear no longer requires a tester profile (it only removes fields); compound ids like UR-002-REQ-031 sort by their REQ number; the Recently-done and testing date windows anchor to the wall clock instead of page-generation time, so a long-open tab keeps meaning "last 24 hours".

## 0.130.3 — Board Names the Real Cause When git Is Unavailable (2026-07-21)

Running the board without a `git` binary (or outside a repo) used to produce a per-ticket anomaly blaming each commit hash. Now a one-time probe logs a single clear line ("git binary not found on PATH"), skips the doomed per-ticket subprocesses, and the anomaly reason says the hash could not be dated rather than asserting it is invalid.

## 0.130.2 — Blocked-REQ Probe Works on Stock macOS (2026-07-21)

Two fixes to the blocked-status machinery in `actions/work.md`.

- The `blocked_check` probe no longer assumes GNU `timeout` exists: it resolves `timeout` → `gtimeout` → unwrapped, so stock macOS (which ships neither) probes the condition instead of failing on exit 127 and wrongly reporting "probe failed".
- The mid-run blocked flip's "no edits landed" test now excludes `do-work/` from its porcelain/diff check — the REQ's own bookkeeping is always dirty mid-run, so the unscoped check could never read clean and silently defeated the flip.

## 0.130.1 — Check Scripts Stop Misreading Diffs and Scope Lists (2026-07-21)

Two bugfix rounds for the shipped verification scripts from the work-pipeline hardening.

- `tools/checks/qualify.sh`: `grep -q` on a piped `git diff` could SIGPIPE the pipeline and mark genuinely-changed files as absent (false WARNs); the diff file list is now computed once. The debug-artifact grep now excludes `do-work/` at the pathspec level, so REQ prose merely *mentioning* console.log/TODO no longer FAILs clean implementations.
- `tools/checks/scope-drift.sh`: inline `**Files I will touch:** \`a\`, \`b\`` lists now parse (previously only bullet lists did — an inline list silently turned the whole check into a SKIP); a touch-list header with zero parseable paths is now a FAIL instead of a silent SKIP; drift path lists print one per line unsplit.

## 0.130.0 — Relative Times Next to Every Board Timestamp (2026-07-21)

Every timestamp on the Kanban board now carries a live relative label — "done Jul 21, 16:24 UTC · 6min ago" — so you can tell at a glance how fresh a card is without doing UTC math.

- Covers card done/cancelled lines, the testing view's tester chips, the detail drawer's Created / Blocked since / Completed / Testing updated rows, the reserved/blocked badge tooltips, and the "Generated" page header.
- Labels tick every second client-side (s → min → h → d buckets), so a tab left open stays accurate; tooltips get a render-time snapshot.

## 0.129.1 — Kill-Stale Guard Matches Cross-Repo Kanban Binaries (2026-07-21)

`just run-kanban` can now restart the board when the port is held by a queue-kanban instance started from *another* repo. Other projects' recipes build the same tool under different names (e.g. `build/go-bin-queue-kanban`), and the old guard only killed a process named exactly `queue-kanban` — so the recipe refused and failed instead of reclaiming the port.

- The kill-stale check now substring-matches `queue-kanban` against the listener's full command line (`ps -o args=`, which unlike `comm` isn't truncated on Linux) and echoes what it stopped.
- Unrelated listeners are still left alone and named in the error, exactly as before.
- Applies to both the shipped `just-kanban` install template (`actions/install.md`) and this repo's own justfile; already-installed projects get the fix as a consent-gated upgrade offer on re-running `do-work install just-kanban`.

## 0.129.0 — Blocked-on-External-Condition Status (2026-07-19)

REQs can now wait on an external condition — LM Studio being up, a designer answering, credentials getting provisioned — instead of being mislabeled as "needs clarification" or dying as a failed environment error. A new `blocked` status names the condition and gets its own badge on the board.

- New `status: blocked` with a free-text `blocked_by` condition (plus optional `blocked_at` and a `blocked_check` shell probe). Distinct from `pending-answers` (a question for you) and `depends_on` (a wait on another REQ).
- `do-work run` re-probes each blocked REQ's `blocked_check` at scan time and auto-unblocks on exit 0 — the same "resolves dynamically" feel as dependency gating. The probe runs the repo-authored command safely (scratch file + `timeout`, fail-closed) and never halts the run.
- Mid-run, when a builder hits a missing external precondition before any edits land, the pipeline flips the REQ to `blocked` and moves on — instead of forcing a `failed` + follow-up cycle.
- `do-work capture` emits `blocked` when a task states it waits on something external; `do-work clarify` now also lets you confirm a blocked condition is met; `do-work abandon` / `roadmap` / `forensics` / `cleanup` all recognize it.
- The Kanban board shows blocked REQs in the *Needs input · Blocked* column with a distinct "blocked by: …" badge and drawer rows (condition, since, probe).

## 0.128.1 — Built-In Preferences Reference Doc (2026-07-18)

The operating nudges people paste at the start of every run — "keep writing lessons learned," "commit often," "I'm AFK, don't block on questions" — are almost all already the skill's defaults. A new reference doc maps each common nudge to where that behavior already lives, so you can stop re-typing them.

- New `docs/standing-preferences.md`: a table of common nudges → the built-in behavior and its home (lessons learned, discovered tasks, YAGNI, per-REQ atomic commits, background agents, non-blocking `pending-answers` questions).
- Calls out the two nudges that are deliberately *not* defaults — an unbounded queue drain (declined in ADR-006/014) and a backgrounded commit — so expectations match reality.
- New README Q&A entry pointing at the reference.

## 0.128.0 — Board Surfaces Completion Anomalies (2026-07-18)

A done REQ with no `completed_at` and no resolvable commit hash used to vanish from the live board — terminal, but with no instant to place it in Recently done. Those are bookkeeping bugs, and now the board shouts about them instead of hiding them.

- New always-visible "Completion anomalies" strip on the board (every view, immune to the recent-window toggle and filters); each card carries an `anomaly` badge, the reason, and the concrete fix, echoed in the detail drawer and the data-warnings banner.
- Three anomaly shapes detected: neither field present, a `completed_at` that doesn't parse (flagged even when the commit hash rescues the date), and a commit-hash field git can't resolve — the reason names the exact broken field.
- Anomalous tickets are never dated "now": no fabricated instant, no Recently-done membership, no mtime fallback — dated tickets keep the existing window behavior unchanged.
- Headless too: `queue-kanban summary` prints `completion anomalies : N` and lists the offending REQ ids.
- Prevention at the source: `actions/work.md`'s done/fail flips and `actions/work-reference.md`'s frontmatter template now make the `completed_at` + `commit` stamp an explicit hard rule on every terminal flip.

## 0.127.0 — Testing View Sorts Newest First and Filters by Date (2026-07-18)

With hundreds of finished REQs, the Ready-to-test column buried the work you just shipped at the bottom. Testing columns now read newest-first, and a date filter narrows them to a window.

- All four testing columns sort most-recent-first — by last testing activity, falling back to the REQ's completion instant; unknown dates sink to the bottom, ties break toward the higher REQ id.
- New date filter in the shared filter bar (visible only on the Testing view): Any date / Last 24 hours / Last 7 days / Last 30 days / Older than 30 days. It joins the existing search/domain/status filters and the Clear button, and never touches the Board or Calendar views.
- This repo's own queue data: four duplicate REQ ids (UR-003's doc-diet stream had reused REQ-015..018 from the earlier kanban stream) renumbered to REQ-021..024, with every frontmatter and audit-trail reference repointed — the board's duplicate-id warnings are gone.

## 0.126.1 — Testing View Review Fixes: Status Gate, Duplicate Keys, Symlink Guard (2026-07-17)

Four PR-review catches on the new Testing view (thanks, Codex review on #119) — all hardening the write path before it ships.

- The status API now rejects non-`clear` transitions on unfinished REQs (409): only terminal-success REQs — or REQs already carrying a testing record, so a returned-then-requeued REQ can restart testing — accept testing writes. A stale browser tab can no longer stamp `in-testing` onto a pending REQ.
- The frontmatter upsert consumes **every** occurrence of a duplicated key, not just the first — the YAML reader keeps the last occurrence, so a first-only edit could look successful yet read back unchanged.
- Testing writes refuse symlinked targets: the REQ file must be a regular file whose parent resolves inside `do-work/`, and `testers.md` gets the same guard — a hostile checkout can't redirect a write outside the tree.
- The Testing view keeps REQs with an *invalid* `testing_status` visible even after their pipeline status leaves terminal-success (the record and its invalid flag no longer vanish on requeue).

## 0.126.0 — Board Testing View: Track Who Tested Which REQ (2026-07-17)

With thousands of REQs, "done" told you nothing about whether anyone actually tested it. The kanban board now has a Testing view (next to Board / Calendar) where a tester picks their profile, selects a finished REQ, and marks it in-testing, tested, or returned with feedback — and the record lives in the markdown itself, so git is the audit trail.

- New Testing view in `do-work board` serve mode: four columns (Ready to test → In testing → Returned with feedback → Tested) over every terminal-success REQ, with per-card actions and an inline feedback form.
- The markdown files are the database: actions write `testing_status` / `tested_by` / `testing_updated_at` / `testing_feedback` placeholder frontmatter into the REQ file via new loopback-only `/api/testing/*` endpoints (surgical line-level upsert — everything else in the file stays byte-identical). No locking by design — changes land in the working tree and commit like any other edit.
- Tester profiles are add-or-select in the view's toolbar, stored as plain bullets in `do-work/testers.md` (created on first use, hand-editable).
- The main Board view shows a `testing` badge on any card carrying a record, and the detail drawer lists the testing meta, so testing state is visible without switching views. Static snapshots render the view read-only.
- Schema Read Contract gains the `testing_status` vocabulary (normalize-and-warn like every other enum — an off-vocabulary value renders as not-tested with an invalid flag and a data warning).

## 0.125.2 — Crew Member Renamed: karpathy.md → coding-guardrails.md (2026-07-16)

The always-loaded implementation crew member is now named for what it does, not for a person — Andrej Karpathy is more than four coding rules. The source attribution inside the file stays.

- `crew-members/karpathy.md` → `crew-members/coding-guardrails.md`; H1 retitled to "Coding Guardrails Crew Member".
- All live references updated (SKILL.md, CLAUDE.md, README, actions, specs, sibling crew files); review-work's audit heading is now "Coding-Guardrails Principle Check".
- Historical records (ADRs, archives) left as written; ADR-003 gained a one-line rename pointer.

## 0.125.1 — Reservation Review Fixes: UR Closure, Release Routing, Roadmap Section (2026-07-16)

Three PR-review catches on the new reservation feature (thanks, Codex review on #118).

- Step 8's UR-finalization check now holds a UR open for **any non-terminal** sibling — a reserved REQ no longer lets its UR archive out from under it.
- `do-work release REQ-042` now actually releases: the router passes `release <rest>` for the `release`/`unreserve` triggers so the reserve action enters release mode instead of trying to reserve the bare ID.
- The roadmap report gained the promised `## Reserved (Other Sessions)` section (with the stale-reservation recategorize hint), a reserved total, and a matching next-step line.

## 0.125.0 — REQ Reservations for Other Worktrees and Cloud Sessions (2026-07-16)

You can now reserve pending REQs for a different worktree or cloud session (`do-work reserve REQ-042 for cloud-alpha`) so the local work loop walks past them. Unlike a claim, a reservation stays in `do-work/queue/` — crash recovery can't steal it — and it travels to sibling checkouts via a normal git sync.

- New `reserved` status in the Schema Read Contract, with `reserved_for` (owner label) and `reserved_at` frontmatter; new `actions/reserve.md` (reserve / release / list).
- The default queue scan skips reserved REQs; targeted `do-work run REQ-NNN` claims them — that's how the owning session picks up its slice.
- Reservations older than 24 hours are flagged as stale everywhere they render (work-loop queue summary, exit summary, forensics, roadmap, board) with a recategorize suggestion — release, claim here, or leave it. Never auto-released.
- The Kanban board shows reserved REQs grayed out in the Claimed column with a "reserved for" badge and a stale marker.
- Intent and contract recorded in `actions/prime-req-reservation.md`.

## 0.124.4 — Qualify and Scope-Drift Checks Tightened Against False Passes (2026-07-15)

Second Codex review round on #117 caught three ways the new checks could be fooled; all three are closed.

- `qualify.sh` no longer counts the previous commit's diff as current work — a no-op builder can't pass on the back of the last REQ's changes.
- `(deleted)` summary entries now need deletion evidence in the working/staged diff, not just disk absence — a typo'd path no longer qualifies.
- `scope-drift.sh` reads only the "Files I will touch" list, so documenting out-of-scope files in "Files I will NOT touch" no longer reports false drift.

## 0.124.3 — Portable Check Scripts and Stale-Baseline Cleanup (2026-07-15)

Two PR-review fixes to the new tools/checks/ scripts (thanks, Codex review on #117).

- Replaced GNU-only `grep -P` extraction (and `\s` ERE classes) with portable `sed`/`grep -E [[:space:]]` — the checks now run on BSD/macOS grep, matching the skill's any-environment contract.
- `preflight.sh` deletes a stale `baseline-failures.txt` when the baseline passes, so Step 6.5 can never misclassify a new regression as pre-existing.

## 0.124.2 — Regrowth Ratchets: Router Word Budget and Sibling-Skill Gate (2026-07-15)

Two guards so the bloat this cleanup removed can't quietly come back.

- Contract tests now fail any commit that pushes SKILL.md past 2,650 words (post-diet count + ~10% headroom); the prescribed fix is a merge or lazy-load, never a bigger budget.
- CLAUDE.md: every NEW action must state why it belongs inside do-work rather than a sibling skill — reviewers reject additions without the justification.

## 0.124.1 — Extraction Plans for the Three Relocatable Subsystems (2026-07-15)

Plan-only release: grep-verified extraction plans for the prompt library, the interview framework, and bkb+dream now live in `decisions/audits/2026-07-15-relocation-extraction-plans.md` (maintainer docs, not shipped). No files moved; nothing changes for consumers in this release.

- Each plan names the target sibling repo, the full manifest with word counts, every inbound-reference seam to cut, and a migration note for git-clone and tarball installs.
- Recommended sequence: prompts → interview → bkb+dream (~47k words would leave the shipped skill if all three run).

## 0.124.0 — Mechanical Work-Loop Checks Ship as Scripts (2026-07-15)

Four parts of the work loop that were pure shell-logic-in-prose are now shipped executables under `tools/checks/`, so they run the same way every time instead of being re-derived from paragraphs. Judgment stays in the prose; mechanics move to code.

- `archive-collision.sh` (Step 2.0, full), `preflight.sh` (Step 5.75, full — also records a machine-readable test baseline for Step 6.5), `scope-drift.sh` (Step 5.5's review-time comparison), `qualify.sh` (Step 6.3's items 1/4/5 + the only-do-work-paths rule).
- work.md steps shrink to pointers + the judgment that remains; every pointer has a script-missing fallback.
- Contract tests now assert the pointers and scripts stay in sync.

## 0.123.2 — Small Actions State Each Guard Once (2026-07-15)

Four action files said the same rules two to seven times over (commit.md stated the .env exclusion in seven places). The guard content survives — stated once, in the section that owns it.

- note.md and scan-ideas.md: Common Rationalizations / Red Flags / Verification Checklist removed — every row mapped 1:1 onto the files' own Rules (mappings recorded in REQ-023, renumbered from REQ-017 in 0.127.0).
- commit.md: step-recap Checklist and "Common mistakes" blocks removed; generic git-advice rationalization rows dropped; the REQ-traceability rows and the hard-won terminal-status Red Flag stay.
- quick-wins.md: two generic rationalization rows dropped; the scan-breadth and dynamic-reference rows stay.

## 0.123.1 — Changelog Trimmed to the Newest 20 Entries (2026-07-15)

The live changelog was 162 entries (~24k words of shipped payload) while the version action only ever reads the newest five. Older entries moved verbatim to `CHANGELOG-archive.md`, which stays in the git repo but is export-ignored from the distribution tarball.

- Live file keeps the newest 20 entries; everything older is in the archive.
- Tarball installs (no `.git`, no archive file) can browse the archive on GitHub — link in the header.
- `actions/version.md`'s "last 5 releases" read is unaffected (first ~80 lines).

## 0.123.0 — Router Diet: One Routing Table, Help Menu Loads Lazily (2026-07-15)

SKILL.md dropped from ~5,500 to ~2,400 words with zero routing changes. The router used to enumerate the action set five times; now the priority table (with the old Verb Reference's disambiguation folded into its Notes column) and the Action Dispatch table are the only two, and the help menu lives in its own action file that loads only when you actually ask for help.

- Actions bullet list deleted — each action file's own blockquote already carries its description.
- Verb Reference merged into the routing table; every trigger verb and precedence rule preserved.
- Help menu + per-command help moved to `actions/help.md` (new `help` dispatch row).
- Every invocation now loads ~3,100 fewer words of router text before your content is touched.

## 0.122.0 — AI-Report Render-Judge Pass and SVG Design Rules (2026-07-14)

The ai-report action now looks at its own output before shipping: when browser automation is available it serves the report over HTTP, takes full-page light+dark screenshots, and judges them against an explicit layout rubric — catching the dead-gutter columns, SVG label collisions, and buried-lede layouts that read fine in source and broke on screen.

- New mandatory Step 7 "Render and Judge": HTTP serve (never `file://` — it screenshots blank in headless Chrome), full-page light AND dark captures (dark via browser color-scheme emulation), fix-and-re-render loop with two passes minimum when any SVG has text labels; graceful footer disclosure when browser automation is absent
- Six-dimension judge rubric applied to the screenshot, not the source: width usage, table shape, diagram informativeness, emphasis hierarchy, theme robustness, SVG label collisions/clipping
- Data-viz rules for hand-authored SVGs: single-hue ordinal ramps for ordered data, ink-colored labels with identity swatches, above/below label lanes with edge-aware text anchors, stat-tile typography
- Reports commit to one coherent aesthetic direction per report via characterful system font stacks (CDN allowlist unchanged: Tailwind + Mermaid only)
- Matching Red Flags, Common Rationalizations, and Verification Checklist entries; user guide updated to match

## 0.121.1 — Recoverable Runs and Leaner Board Loading (2026-07-13)

Fan-out runs now distinguish “assembled” from “delivered,” so an interrupted review or exploration can resume without cleanup deleting its only result. The board also keeps exact-copy Markdown out of the initial payload until someone actually presses Copy.

- Added `in-progress` → `synthesized` → `consumed` run states, persisted code-review reports, a root deep-explore manifest, and consumed-only cleanup with explicit staging for deleted run paths
- Aligned cleanup's five-pass documentation, changelog-title examples, and prime's interactive questions with their canonical contracts
- Moved raw REQ/UR Markdown into lazy `board-markdown.js`; the current tree's initial `board-data.js` is 43% smaller while generated and live boards still copy exact source text

## 0.121.0 — Tidy-Repo Rename and Safer Layout Planning (2026-07-13)

`file-reorg` is now `tidy-repo`: a clearer name for the same reference-safe repository-layout job, with the old command retained as a compatibility alias. The workflow is tighter about what belongs in a layout pass and more careful around real-world repositories that already have local changes, generators, or platform-sensitive paths.

- Renamed `actions/file-reorg.md` to `actions/tidy-repo.md` and promoted `do-work tidy-repo [path] [plan]` across routing, help, dispatch, README, and next-step guidance
- Added an explicit target-design step, dirty-path overlap handling, generated-source mapping, case-only rename handling, and post-move diff verification
- Made README/CLAUDE edits conditional on actual layout drift; unrelated link fixes, boilerplate rewrites, and permanent link-checker creation are follow-up work instead of mandatory side effects
- Preserved `do-work file-reorg` as a legacy alias so existing prompts keep working

## 0.120.0 — Run Dirs Are Committed, Then Cleaned Up on Consumption (2026-07-13)

Fan-out run directories (`do-work/runs/`) are no longer gitignored transient scratch — they're now committable, so a review or exploration is visible and doesn't get silently lost mid-run. In exchange, the run dir gets deleted the moment its findings are consumed (synthesized and promoted to a report, REQs, or deliverables), which keeps `do-work/runs/` from growing without bound. That whole create → inspect → promote → delete lifecycle is now part of the job, not an afterthought.

- `.gitignore` no longer excludes `do-work/runs/` (`do-work/pipeline.json` stays excluded — it's live state, not work).
- `crew-members/background-agents.md` is the canonical lifecycle: run dirs are committable (step 1) and deleted once consumed (new step 5). The old `.git/info/exclude` append for run dirs is gone.
- `code-review` and `deep-explore` now delete their run/session directory as the final step, after promoting anything worth keeping into `do-work/deliverables/`.
- `cleanup` gains a safety-net pass that sweeps abandoned `Status: complete` run dirs (and leaves incomplete, possibly-resumable ones alone).
- The shared local-ignore snippet still used by `pipeline.json`, the vendored `last30days` engine, and build artifacts moved to a dedicated section in `background-agents.md`; its former callers point there.

## 0.119.0 — Board Drawer Copy Button (2026-07-11)

The Kanban board's ticket drawer gets a Copy button next to Close: one click puts the open REQ's (or UR's) raw Markdown on the clipboard, ready to paste into chat, email, or another ticket without losing headings, checkboxes, or links.

- The data island now ships `bodyMarkdown` beside the pre-rendered `bodyHtml`, so the copy is the ticket's source text, not scraped HTML.
- Transient feedback ("Copied ✓" / "Copy failed") resets on every drawer open; a hidden-textarea fallback covers contexts where the async Clipboard API is missing or denied (file://, plain http).

## 0.118.0 — Cleanup Repoints Doc Links to Moved Files (2026-07-11)

Cleanup's consolidation passes move REQ files around the archive, which used to silently break any doc that linked to them (one consumer repo hit 39 broken prime-doc links). Cleanup now records every move's old → new path and rewrites the referring links itself.

- New `Repoint Documentation Links` step in `actions/cleanup.md`: after all passes, filename-grep tracked markdown outside `do-work/` for each moved file and rewrite link targets from the per-move mapping — preserving `#anchors`, skipping bare prose mentions, tracked files only by design.
- Summary gains a `Repointed: N doc links in M files` line (`Repointed: none` when nothing referenced the moved files, so the step visibly ran).
- The cleanup commit stages the rewritten docs alongside the moves they repair; `docs/cleanup-guide.md` documents the behavior.

## 0.117.1 — Retroactive Descriptive Changelog Titles (2026-07-11)

The descriptive-title convention from 0.117.0 now applies to the whole file: all 152 pre-0.117.0 codename headings ("The Red Pen", "The Court Scribe", …) were rewritten to say what each release delivered. Bodies are untouched — only the heading titles changed.

- Every `## X.Y.Z — The [Codename] (date)` heading from 0.65.0 through 0.115.0 replaced with a short descriptive title derived from that entry's own body.
- Verified no duplicate titles across the file and no codename headings remain.
- CLAUDE.md's "leave pre-0.117.0 entries as-is" note removed — it no longer applies.

## 0.117.0 — Board View Filters (2026-07-11)

The board's By-UR lens rendered the entire archive — after months of history it was an archive dump, not a work view. Every view now filters: a shared search + domain/status bar in the topbar, and an Active/All toggle that hides fully resolved URs by default.

- By-UR lens defaults to Active (URs with at least one unresolved REQ); a footer note counts the hidden resolved URs, and All brings them back.
- Shared filter bar applies to whichever view is active: search matches REQ/UR ids and titles, domain and status selects populate from the data. Column and UR counts read "shown / total" while filtering; the calendar hides days with no matches.
- A search hit on a UR header keeps its whole group visible (domain/status still filter the cards inside).

## 0.116.1 — Clear Questions in Review-Work Follow-Ups (2026-07-11)

0.116.0 required cold-reader question authoring in work.md's follow-ups but missed the copy-paste sibling: review-work's ambiguous-requirements follow-ups emit the same `Recommended:`/`Also:` template. A grep for every `pending-answers` authoring site found this one remaining gap.

- `actions/review-work.md` ambiguous-requirements follow-ups now load `crew-members/clear-questions.md` and author Open Questions for a cold reader (gloss shorthand, state why the decision is the user's — Principle 7), matching work.md Step 8.

## 0.116.0 — Escalated Questions Explain Themselves (2026-07-11)

Escalated questions were reaching the user written in builder shorthand — technically asked, practically unanswerable. Now clarity is enforced at both ends: builders author Open Questions for a cold reader, and clarify rewrites what slips through.

- `actions/clarify.md` Step 3 now loads `crew-members/clear-questions.md` and rewrites stored question text to its contract instead of rendering it verbatim.
- New clear-questions Principle 7: an escalated question must say why the decision is the user's — the rule that forced the escalation and what silently deciding would have cost.
- `actions/work.md` Step 8 and the follow-up template in `actions/work-reference.md` require Open Questions destined for clarify to meet the contract at authoring time.

## 0.115.0 — Board Flags Invalid REQ Statuses (2026-07-10)

The Kanban board now marks a REQ whose `status:` is outside the schema vocabulary as *invalid* — red status, an INVALID pill on the card, and a drawer note telling you exactly how to fix it — instead of letting it blend in with normal blocked tickets. Came out of triaging review feedback: the live-tree bucketing test contradicted the board's own deliberate catch-all and would have failed on any off-vocabulary status.

- `bucketColumns` flags off-vocabulary tickets (`StatusUnrecognized`), and its warning now carries the fix prompt (edit `status:` per the Schema Read Contract, or run `do-work forensics`).
- New forensics check 11 sweeps queue/working/archive for unrecognized statuses — the mechanical fix path the board's warning points at.
- `TestLiveTreeColumnBucketingMatchesStatus` now asserts the real invariant (unrecognized statuses legitimately live in Needs-input *when flagged*), plus a seeded synthetic regression test so the live queue can't mask it.

## 0.114.0 — Retire the Weekly-Signal-Diff Prompt (2026-07-10)

Retired the `weekly-signal-diff` prompt from the library. It graduated into the consumer project's own `wsd-skill` (as `daily-signal-diff`, driven by the `wsd-full` / `wsd-go` / `wsd-refresh` family) months ago — the shipped copy was a stale duplicate that every `do-work update` kept reinstalling.

- Removed `prompts/weekly-signal-diff.md` and `prompts/weekly-signal-diff-personal.md`; dropped their rows from `prompts/README.md`.
- `decisions/imported-specs/2026-04-17_improve-weekly-diff-skill.md` gained a Status footer recording the removal; changelog history stays as-is.
- The `**Runnable:**` header key in `actions/prompts.md` is generic and remains — it just no longer has a shipped opt-out example.

## 0.113.2 — Drawer Formatting for Questions and Prose (2026-07-10)

The drawer was mashing a REQ's Open Questions into one run-on paragraph and stretching prose across the whole panel. Both readable now.

- `Recommended:` / `Also:` / `Value:` / `Risk:` / `→` continuation lines render on their own lines instead of lazily merging into the question sentence (fenced code blocks stay verbatim).
- Markdown body text caps at ~90 characters per line, so a wide drawer no longer means 200-character lines.

## 0.113.1 — Notes Strip Parses Only Bullet Lines (2026-07-10)

The Notes strip was reading a real `notes.md` as eighteen notes when it held two. Only bullet lines are notes now.

- The `#` heading, the prose preamble, and horizontal rules are skipped instead of rendered as notes.
- `<!-- ... -->` comment blocks are stripped **before** the bullet test — that's where pruned entries get parked, and their bullets were resurfacing on the board.
- `do-work roadmap` and `do-work note` carry the same rule, so every reader of `notes.md` agrees on what a note is.

## 0.113.0 — Board Dependency Graph: Ready vs Waiting (2026-07-10)

The board finally draws the dependency graph it was already parsing. Pending now separates what you can pick up right now from what's still waiting on an upstream REQ, and every card tells you how much is waiting on *it*.

- **Ready vs. Waiting.** The Pending column splits in two. When nothing is waiting, it stays a flat list — no new headers for a queue without dependencies.
- **Unblocks N.** A card carrying that badge is the one to work on: N unresolved REQs are waiting for it. The full list is in the detail drawer.
- **Dangling dependencies are now loud.** A `depends_on` pointing at a REQ that isn't in the tree fails closed (the dependent stays waiting, never quietly ready) and raises a data warning — it can never self-resolve.
- Dependency chips show met (struck through) vs. unmet (amber), and the drawer lists each dependency with the status that decides it. `cancelled` never satisfies gating, matching the work loop.
- `do-work board summary` now prints the ready / waiting breakdown.

## 0.112.0 — Notes Strip on the Kanban Board (2026-07-10)

Your `do-work note` hints now show up on the Kanban board, not just in `do-work roadmap`. They sit in a collapsible Notes strip above the columns, so the thing you told yourself to check next is visible while you're staring at the queue.

- `do-work board` reads `do-work/notes.md` and renders each line with its date, in append order.
- The strip stays visible in the calendar view too, and disappears entirely when there are no notes.
- Notes render as plain text, never Markdown — they're hints, not tickets, so they get no column, no calendar entry, and no detail drawer.
- Serve mode watches `notes.md`, so appending a note and reloading the page shows it.

## 0.111.0 — Versioned Changelog Entries in Target Repos (2026-07-09)

Changelog entries in unversioned repos came out keyed by date alone, so nothing told you whether an entry was a typo fix or a rewrite. Every entry now carries a version and a date, and the number is earned — bumped by what the change actually did to people using the code.

- Entry key is always `## X.Y.Z — The [Codename] (YYYY-MM-DD)`
- Version source resolves in order: a version file in the repo (bumped and staged with the REQ commit), release tags (read, never created — a tag is a human's release call), or the changelog's own counter seeded at `0.1.0` for repos with no version at all
- Bump size reads the delivered change: breaking a consumer is major, a new user-invocable capability is minor, everything else is patch. Ties break downward; below `1.0.0` a breaking change bumps the minor, so a seeded repo never silently promotes itself to a `1.0.0` release
- Fixes a duplicate-header bug on the versioned path, which reused the repo's current version for every entry instead of bumping it
- Guards added for disagreeing version files (leave them alone, fall back to the counter, report it) and for out-of-band releases (bump from whichever source is higher)
- The commit's "did we actually stage an implementation?" check now knows the version file is bookkeeping, not implementation — so a lone version bump can't masquerade as delivered work

## 0.110.0 — Work Pipeline Writes Target-Repo Changelogs (2026-07-07)

This changelog was the only one do-work ever kept — every target repo's history lived in commit messages nobody rereads. Now the work pipeline writes a changelog entry in every repo it works in, by default, in the house voice (picked from a six-voice side-by-side style lab over four real entries).

- New **Changelog Entry Procedure (Step 9)** in `actions/work-reference.md`: house-style contract (value-first lead + technical bullets), `## YYYY-MM-DD — The [Codename]` keys for unversioned repos, the repo's own version when it has one — never invented
- Bootstrap when `CHANGELOG.md` is missing; an existing changelog in a different format wins over the house voice
- Successful REQs only — failed and cancelled work gets no entry; `CHANGELOG.md` joins the explicit staging list and doesn't count as implementation in the commit validation check
- Wired into `actions/work.md`'s Commit Phase; entries load `crew-members/anti-slop.md` like any human-facing artifact
