# Cleanup Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to tidy the archive, or automatically at the end of the work loop. Consolidates loose files and ensures the archive is well-organized. User-facing walkthrough: [`docs/cleanup-guide.md`](../docs/cleanup-guide.md).

The archive should be a collection of self-contained UR folders, each containing their original input and all related REQ files. Over time, REQ files can end up loose in the archive root — either from intermediate archival (when not all REQs were done yet) or from legacy requests predating the UR system. This action fixes that.

## When to Use

**Use when:**
- User wants to tidy the archive — organize loose REQs into UR folders
- User says "cleanup", "clean up", "tidy", or "consolidate"
- Automatically at the end of the work loop

**Do NOT use when:**
- User wants *diagnostics* on pipeline health — route to actions/forensics.md instead
- User wants to *delete* or discard work — cleanup only reorganizes work items (URs, REQs), never deletes them. Three narrow exceptions: Pass 4 sweeps *consumed* run scratch (a `Status: consumed` directory under `do-work/runs/`) after its findings have been promoted — spent scratch, not work; Pass 5 removes an orphaned `worktree-agent-*` git worktree and branch mechanically only when it is clean, merged, and its exact REQ is positively settled outside `do-work/working/`, and otherwise only with the user's explicit consent; and Pass 6 *restores* content to archived REQs that were blanked by an unguarded write-back — it writes work back in rather than reorganizing it, which is why it asks first.

## When This Runs

- **Automatically** at the end of every work loop (after all pending REQs are processed)
- **Manually** when the user invokes it (e.g., `do-work cleanup`, `do-work tidy`)

## Steps

Resolve the project root once, then invoke the canonical deterministic command:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> cleanup
```

Pass `--dry-run` for an evidence-only plan, `--commit` for one exact-path cleanup commit, and repeat exact consent tokens as `--restore-blanked <repo-relative-path>` or `--discard-worktree <worktree-agent-* branch>`. `--commit` and `--discard-worktree` are intentionally separate runs: a Git commit cannot roll back a forceful worktree discard. Global `--format json` goes before `cleanup`. A missing launcher, build failure, nonzero exit, or malformed result stops the action with the command's actionable finding; do not fall back to free-form file or Git mutation. The command owns discovery, Passes 0–6, link repointing, target guards, mutation, rollback, and optional commit. The pass descriptions below remain the human policy contract the command implements.

Seven passes, in order:

### Pass 0: Sweep Finished Queue Items

Scan `do-work/queue/` and the working directory for REQs with terminal statuses that should have been archived but weren't — typically from manual work, different agents, or legacy sessions that completed outside the standard work pipeline.

**Safe because it sweeps only terminal statuses.** A REQ carrying `completed`, `failed`, `cancelled` or any of their variants is finished by definition — whoever finished it, and however many builders were running at the time — so moving it out of `working/` takes nothing away from anyone. That is the durable reason, and it holds at any builder count. **Do not argue this pass's safety from the Execution Model's ownership rule**: that rule is about who *releases* (`actions/work-reference.md` → **Execution Model — Claim Anywhere, One Releaser**), and it does not say a claimed `working/` file is this session's — any checkout may have claimed it — `actions/work-reference.md` → **Crash Recovery (Step 1)** exists precisely because it may not be. Pass 0 never touches a `claimed` REQ (step 4 below), which is what keeps the two rules from colliding. When cleanup runs at the end of a work loop (`actions/work.md` Step 10), Step 8 has already moved out every REQ that run finished — plural under fan-out (`actions/work-reference.md` → **Worktree Dispatch Mode** → Fan-Out Dispatch) — so no in-flight REQ sits terminal-in-`working/` for Pass 0 to sweep out from under a builder.

1. **Glob `do-work/queue/REQ-*.md`**
2. **Read each REQ's frontmatter** `status` field
3. **If status is a canonical terminal value or a documented status alias** — normalize the alias under the Schema Read Contract before moving, then treat its canonical terminal value as the archival decision:
   - **Normalize non-standard statuses** before moving according to the Schema Read Contract's `status` row; do not maintain a second alias list here
   - Move the REQ to `do-work/archive/` root (Pass 1 and Pass 2 will then consolidate it into the correct UR folder)
   - Report: `Swept REQ-NNN from do-work/queue/ (was status: {original}) → archive`
4. **Leave `pending`, `pending-answers`, `blocked`, and `claimed` REQs untouched** — those are active queue items (`blocked` waits on an external condition)
5. **Also check `do-work/working/`** — if any REQ there has a canonical terminal status or a documented status alias, it was finished but never moved out (a crashed prior run). Apply the same Schema Read Contract normalization, move it to `do-work/archive/` root, and report it. Moving it out of `working/` also drops its `## In Progress (interrupted)` entry from `do-work/CHECKPOINT.md` — **this checkout's own-label entry only; one carrying another checkout's `writer:` label records a claim held elsewhere and is left alone.** `actions/work-reference.md` → **In-Progress Record (Step 2)** states when that removal is owed, how the label is derived, and why the removal is part of the move rather than a later sweep.

### Pass 1: Close Completed User Requests

Check `do-work/user-requests/` for UR folders that are ready to archive.

**The closure condition, stated — not a stored list.** A UR is ready to archive only when **every REQ carrying `user_request: UR-NNN` in its frontmatter, wherever it currently sits, is terminally resolved.** Membership is derived by scanning that field; it is never read off the UR's `requests:` array. That array is the capture-time record of the REQs capture itself created (`actions/capture.md` Step 5) and nothing maintains it afterward — review-spawned follow-ups, addendum REQs, and clarify-derived REQs all carry `user_request:` without ever being appended to it. Keying on the array archives a UR whose follow-ups are still queued and strands their UR reference. This is the same predicate `actions/work.md` Step 8 evaluates on its archive path; the two readers must not drift apart.

For each UR folder in `do-work/user-requests/`:

1. **Collect the UR's REQs** by reading the `user_request` field of every `REQ-*.md` in all four locations and keeping those whose value is this UR's id:
   - `do-work/queue/REQ-*.md` (pending, pending-answers, blocked, claimed)
   - `do-work/working/REQ-*.md` (in flight)
   - `do-work/archive/REQ-*.md` (loose in archive root — non-recursive; `archive/legacy/` REQs have no `user_request` by definition)
   - `do-work/archive/UR-NNN/REQ-*.md` (already consolidated)
2. **Normalize each collected REQ's status, then check the terminal-resolved set** (see `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set; don't restate or fork either list here). Any status outside it holds the UR open, **`failed` included** (how a `failed` REQ is resolved so it leaves this held-open state is defined at that canonical statement — do not re-derive it here).
   If the same REQ-ID is found in **both** `do-work/archive/` root and `do-work/archive/UR-NNN/`, flag it and leave the UR in `user-requests/` untouched: `⚠ Duplicate: REQ-NNN found in both archive/ root and archive/UR-NNN/. Resolve manually, then re-run cleanup.`
3. If **ALL** collected REQs are terminally resolved (and no duplicates flagged):
   - Gather any loose completed/cancelled REQ files from `do-work/archive/` root into the UR folder
   - Move the entire UR folder to `do-work/archive/UR-NNN/`
   - Report: `Archived UR-NNN (all N REQs resolved)` — when any were cancelled, say so: `(N-K complete, K cancelled)`
4. If **NOT all** collected REQs are terminally resolved:
   - Leave the UR folder in `user-requests/` — it's not ready yet
   - Report: `UR-NNN still open (X/Y REQs complete)`. When the only unresolved members are `failed` REQs, also name them and the exit: cancel each with `do-work abandon REQ-NNN` (`actions/abandon.md`), which flips it to `cancelled` in place so the UR closes on the next cleanup. This is the one transition out of `failed` — a completed follow-up REQ does the recovery work but never resolves the original, so the failed REQ must be cancelled either way. Without this line a UR held open by a failure looks stuck with no stated way out.
5. **Report-only cross-check against the array.** Any REQ id listed in the UR's `requests:` array that the scan in step 1 found in none of the four locations is a missing file, not an open REQ: report `⚠ REQ-NNN listed in UR-NNN's requests: array but found nowhere`. It must never hold the UR open — a stale array entry that wedges closure is the failure this pass was rewritten to avoid.

### Pass 2: Consolidate Loose REQ Files in Archive

Check `do-work/archive/` root for any `REQ-*.md` files that should be inside a UR folder.

For each loose `REQ-*.md` file directly in `do-work/archive/` (not inside a subfolder):

1. Read its frontmatter and check for a `user_request` field
2. **If it has `user_request: UR-NNN`:**
   - Check if `do-work/archive/UR-NNN/` exists
   - If yes: move the REQ file into that UR folder
   - If no: check if `do-work/user-requests/UR-NNN/` exists (UR still open — leave the REQ in archive root for now; it will be consolidated when the UR is fully complete and archived by Pass 1)
   - If the UR folder doesn't exist anywhere: report a warning — `REQ-XXX references UR-NNN but no UR folder found`
3. **If it has NO `user_request` field (legacy/standalone):**
   - Move it to `do-work/archive/legacy/` (create the folder if needed)
   - Report: `Moved REQ-XXX to archive/legacy/ (no UR reference)`

### Pass 3a: Misplaced do-work Directories Elsewhere in the Repo

Scan for `do-work/` directories created inside utility subdirectories instead of the project root. This happens when an agent's working directory drifts into a subdirectory (e.g., during a refactor) and the next capture creates `do-work/` relative to that location. Once the misplaced directory exists, subsequent sessions keep writing there — silently diverging from the canonical queue.

1. **Detect directories, not file patterns.** Search for any directory named `do-work/` anywhere in the repo EXCEPT the project root. Look for the directory itself — don't rely on specific file patterns inside it, since a misplaced tree may contain only `user-requests/`, only `working/`, only assets, or any partial subset of the normal structure.
2. For each misplaced `do-work/` found, inspect its known subtrees (`archive/`, `user-requests/`, `working/`, and `queue/` REQ files). Relocate preserving internal structure:
   - **Queue REQ files** (`do-work/queue/REQ-*.md`): move to canonical `do-work/queue/REQ-*.md`. **Before moving**, check if a REQ with the same number already exists at the canonical location (Pass 0 sweeps terminal-status REQs, but a misplaced `do-work/` may have a REQ with a status Pass 0 doesn't touch, such as `pending`). Conflict = same REQ number exists at both locations — report and leave the misplaced copy in place for manual resolution.
   - **`user-requests/UR-NNN/`**: move entire folder to canonical `do-work/user-requests/UR-NNN/`. Conflict = same UR number exists at both locations.
   - **`archive/UR-NNN/`**: move entire folder to canonical `do-work/archive/UR-NNN/`. Conflict = same UR number exists at both locations.
   - **`working/REQ-*.md`**: move to canonical `do-work/working/REQ-*.md`. Conflict = same REQ number exists at both locations.
   - **Other files/dirs**: move to matching path under canonical `do-work/`. Conflict = same path already exists.
   - **Conflict handling**: when the same item exists at both locations, do NOT overwrite — report the conflict with both paths and leave the misplaced copy in place for manual resolution.
   - Report: `Found misplaced do-work/ at {path} — relocated {N} items to project root` (and list any conflicts separately)
3. After relocating all non-conflicting contents, remove the misplaced `do-work/` directory if empty. If conflicts remain, leave it in place.

### Pass 3b: Misplaced Folders Within the Archive

Check for UR folders that ended up in wrong locations within the archive.

1. Check if `do-work/archive/user-requests/` exists (this is a common mistake — the entire `user-requests/` dir got moved instead of individual UR folders)
2. If it exists, for each `UR-NNN/` folder inside it:
   - If `do-work/archive/UR-NNN/` does NOT already exist: move it up to `do-work/archive/UR-NNN/`
   - If `do-work/archive/UR-NNN/` DOES already exist: merge contents (move files from the misplaced folder into the correct one)
   - Report: `Fixed misplaced UR-NNN (was in archive/user-requests/)`
3. If `do-work/archive/user-requests/` is now empty, remove it

Also check for and consolidate any loose CONTEXT-*.md files:
- Move to `do-work/archive/legacy/` alongside legacy REQs

### Pass 4: Sweep Consumed Run Directories

Fan-out actions (code-review, deep-explore, multi-REQ work — see `crew-members/background-agents.md`) each delete their own `do-work/runs/<action>-<ts>/` directory once its findings are consumed. This pass is the **safety net** for the narrow case where an owner recorded `Status: consumed` after delivery/promotion but the session ended before the following deletion. Synthesis alone never qualifies.

1. **Glob `do-work/runs/*/`** (each is one run directory).
2. **Read each run's `manifest.md`** and check its `Status:` line.
3. **If `Status: consumed`** — the run's findings were delivered/promoted and the directory is spent scratch. **Delete it.** Record the exact deleted path for the Commit section and report: `Swept run dir do-work/runs/{name} (Status: consumed)`.
4. **If the manifest is missing, or `Status:` is anything other than `consumed`** (`in-progress`, `synthesized`, or the legacy `complete`) — **leave it untouched** and report its actual status. An in-progress run may need missing work; a synthesized/legacy-complete run may contain an assembled output the user never received. See `crew-members/background-agents.md` for the recovery branches; never infer consumption from synthesis.

This is the only pass that deletes anything **under `do-work/`** rather than reorganizing it, and it is scoped strictly to **consumed run scratch** — a `Status: consumed` directory under `do-work/runs/` only. URs, REQs, and every other `do-work/` artifact are still only ever moved, never deleted. (Pass 5 also deletes, but its objects are git worktrees and branches outside `do-work/`, and it only discards unmerged work with the user's explicit consent.)

Tracked run files use the normal Git transaction and rollback. An entirely untracked consumed run cannot be restored from Git, so the command revalidates the exact file inventory and `Status: consumed` immediately before deleting it and labels those changes `non-rollback spent-scratch deletion`. This narrow exception follows the spent-scratch invariant above and does not weaken the dirty-target rule for any durable or mixed tracked/untracked group.

### Pass 5: Orphaned Worktrees (consent-gated)

Removes the git worktrees and branches left behind when `actions/work.md`'s worktree dispatch mode is interrupted (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**). Only `worktree-agent-REQ-NNN-*` names are in scope — that naming convention is what makes a leftover attributable to a REQ; every other worktree is the user's own and is never touched.

**Non-interactive runs never delete a consent-required leftover.** <a id="interactivity-test"></a>A run is *interactive* only when both hold: your environment can put a prompt in front of a human (an ask-user mechanism exists), **and** the run was not launched unattended (subagent, cron, CI, `--yes`-style automation). An unconfirmed or ambiguous result counts as non-interactive. The automatic end-of-work-loop invocation (`actions/work.md` Step 10) is exactly that case — unattended, often dispatched as a subagent — so it may remove a leftover proved clean, merged, and settled, but it only reports every other leftover.

1. **`git worktree prune` first** — it drops administrative entries whose directories are already gone, so the enumeration below doesn't report ghosts. Not a git repo, or no `worktree` subcommand → skip the whole pass.
2. **Enumerate and attribute** with `git worktree list --porcelain`, `git branch --list 'worktree-agent-*'`, and one fresh repository discovery. A branch with no worktree and a worktree with no branch each count as a leftover. Read an exact REQ id only from the `worktree-agent-REQ-NNN-*` prefix; a missing, malformed, unreadable, or multiply claimed identity is not proof that the run finished.
3. **Remove only when all three facts are positive:** the worktree is clean, its branch or detached head is merged into the current integration `HEAD`, and the exact readable, unambiguous REQ is outside `do-work/working/`. Then use `git worktree remove <path>` (no `--force`) and `git branch -d <branch>` **from that integration branch**, and report `Removed merged worktree <name>`. A stale initial discovery is not enough: Pass 5 discovers again after Passes 0–4 have applied, so a terminal REQ those passes just moved out of `working/` can qualify while a still-active one cannot.
4. **Any missing fact requires consent.** Dirty, unmerged, still-working, absent, ambiguous, malformed, and unreadable cases all use the same consent-required finding. Never `-D`, never `--force` without the exact `--discard-worktree <name>` consent token. List the leftover with its branch, worktree path, REQ id when available, and failed facts; then load `crew-members/clear-questions.md` and ask with your environment's ask-user prompt whether to discard it, naming what is lost if the answer is yes. Delete only what the user explicitly names. No answer — or a non-interactive run — leaves it in place, reported.

Unlike Pass 4, nothing here is spent scratch by construction: a `worktree-agent-*` lane without all three positive facts can still be active or hold the only copy of a builder's work. That is why this is the one pass that asks before deleting.

### Pass 6: Restore Blanked Archived REQs (consent-gated)

Recovers archived REQ and UR files whose content was destroyed by an unguarded `commit:` write-back — the file survives as 0 bytes or with its frontmatter gone, so the decision trail it held is only in git history. (Origin incident: six archived REQs, 8 KB to 26 KB each, were truncated to nothing in a consumer repo; each blanking commit reported `1 file changed, N deletions(-)` and nothing read that number.) `tools/checks/record-commit-hash.sh` is the guard that prevents new damage; this pass repairs damage that already happened.

**Runs last, after the passes above have moved files.** A blanked REQ that Pass 1 or Pass 2 just consolidated into a UR folder is scanned at its final path, so the restore writes where the file now lives and the Commit section stages one path per file rather than a move plus an edit.

**The window is finite.** The lost content lives only in unreferenced git objects until a `git gc` collects them, so a detected-but-deferred restore can become unrecoverable. Report that alongside any finding the user declines.

**Non-interactive runs report only, never write.** Same interactivity test as Pass 5 above (*Non-interactive runs report only*); an ambiguous result counts as non-interactive. The automatic end-of-work-loop invocation (`actions/work.md` Step 10) is that case.

1. **Detect.** Run the shipped scanner — read-only, and a no-op exit when `do-work/` holds nothing damaged:

   ```bash
   <skill-root>/tools/checks/blanked-req-scan.sh
   ```

   Nothing found → the pass is done, report nothing. If the script is missing, skip the pass and report that it is missing; do not hand-roll the git archaeology.

2. **Show the plan.** Run the same scanner with `--restore --dry-run`. It names each file, the commit its content would be recovered from, the byte count, and the `commit:` hash it would re-apply — and writes nothing.

3. **Ask.** Load `crew-members/clear-questions.md` and ask with your environment's ask-user prompt whether to restore the listed files. State what a *yes* does: each file's content is overwritten with the recovered blob, so any edit made to the file since it was blanked is discarded. Restore only what the user approves.

4. **Restore.** On approval, run `<skill-root>/tools/checks/blanked-req-scan.sh --restore`. It writes each file via a temp file in the file's own directory plus an atomic rename, refuses to write recovered content that is itself empty, and re-applies `commit:` by calling `tools/checks/record-commit-hash.sh` — never by editing frontmatter itself. Exit 0 means every damaged file was fully repaired — content back *and* the recorded hash re-applied. Exit 1 means at least one was not, with a `SKIP:`/`FAIL:` line naming which and why, and it covers two different states: a file nothing could be written for, and a **partial** repair where the content is back but its recorded `commit:` hash would not apply. Read which. Record each restored path for the Commit section and report: `Restored <path> (<N> bytes from <sha>, commit: <hash>)`.

**Never commit a partial repair as-is.** A partially repaired file looks healthy — full content, parseable frontmatter — while carrying `commit:` provenance that points at the wrong commit or nothing at all, which is precisely the kind of quiet falsehood in the Trail of Intent this pass exists to undo. Resolve it first by running `tools/checks/record-commit-hash.sh <path> <hash>` and reading its output (the failing guard is named there), or hand the file back to the user with that output. Commit the cleanly restored paths either way; hold the partial one out of the commit until its hash is right.

5. **A file with no recovery source is reported, not restored.** The scanner emits `-` for a file git has no prior content for (never committed non-empty, or history rewritten past it). Say so plainly — it is a permanent loss and the user should check backups.

Unlike every pass above, this one writes *into* work items rather than moving them, which is why it asks: the recovered blob is the pre-blanking commit's version, so it silently discards anything written to the file afterward.

### Repoint Documentation Links

Durable docs outside `do-work/` may link to files the passes above just moved (e.g. a prime doc's `## Lessons` section linking `[REQ-987](../do-work/archive/REQ-987-slug.md)`). The move is the only moment both the old and new path are known, so repointing is part of cleanup — not a separate "find broken links" sweep afterward.

1. **As any pass moves a file, record its old → new repo-relative path.** This applies to every move cleanup makes, whichever pass makes it — the passes above are the current set, not a closed list.
2. **After all passes**, for each moved file, search the repo's tracked markdown outside `do-work/` for references to it. Match on the **filename** — REQ filenames are unique, and referrers use relative paths, so matching the full old path would miss them:

   ```bash
   git grep -l -F 'REQ-987-slug.md' -- '*.md' ':(exclude)do-work'
   ```

   `-F` because filenames contain dots; `git grep` searches tracked files only, so untracked noise and `do-work/` internals are excluded by construction.
3. **For each hit, rewrite the link target** to the correct relative path from the linking file's directory to the file's new location, using the old → new mapping from step 1. Three guards:
   - **Preserve anchors.** A link like `REQ-987-slug.md#lessons-learned` keeps its `#fragment` suffix — rewrite only the path portion of the target.
   - **Rewrite path occurrences, not bare mentions.** The filename grep also hits prose that mentions `REQ-987-slug.md` with no path component. Rewrite occurrences of the old path (any relative spelling of it); leave a bare filename with no path component alone — never graft a path onto a prose mention.
   - **Tracked files only, by design.** `git grep` won't see an untracked, not-yet-committed doc that links to a moved file. That scope is deliberate (link-checking tests validate tracked files); the repoint does not guarantee zero broken links in untracked drafts.

Risk note: a bad rewrite could mangle a doc, but it's git-reversible, the change is reviewable in the cleanup commit diff, and any link-checking test the repo runs doubles as the regression detector.

## Reporting

Print a summary at the end:

```
Archive cleanup complete:
  - Swept: 3 finished REQs from do-work/queue/, 1 from working/
  - Archived: UR-011 (3 REQs), UR-004 (8 REQs)
  - Consolidated: 5 loose REQs into their UR folders
  - Legacy: 24 REQs moved to archive/legacy/
  - Misplaced do-work/: relocated 7 REQs, 6 URs from exp/g3-segment-anything/do-work/
  - Fixed: 1 misplaced UR folder in archive
  - Swept runs: 2 consumed run directories
  - Worktrees: removed 2 clean, merged, settled worktree-agent-* leftovers; 1 unfinished or unattributed leftover reported (awaiting consent)
  - Restored: 2 blanked archived REQs (REQ-1282, REQ-1287); 1 unrecoverable reported
  - Repointed: 39 doc links in 5 files
  - Still open: UR-015 (2/4 REQs complete)
```

When files were moved but no referrers were found, still print `Repointed: none` — the line is evidence the repoint step ran.

If nothing needed fixing:
```
Archive is clean. No loose files or pending closures found.
```

Print that line only when Passes 5 and 6 also came back empty. A clean archive says nothing about worktrees or blanked content — a well-organized archive full of 0-byte REQs satisfies every other pass. When Pass 5 found leftovers, or Pass 6 found damage, each reports it even if every other pass was a no-op.

## Archive Structure After Cleanup

```
do-work/archive/
├── UR-001/                    # Self-contained: input + all REQs
│   ├── input.md
│   ├── assets/
│   ├── REQ-018-feature.md
│   └── REQ-019-feature.md
├── UR-002/
│   ├── input.md
│   └── REQ-024-feature.md
├── legacy/                    # REQs and CONTEXT docs without UR references
│   ├── REQ-001-old-task.md
│   ├── REQ-002-old-task.md
│   └── CONTEXT-001-batch.md
└── hold/                      # Items on hold (paused by user — cleanup skips these)
```

**No loose REQ or CONTEXT files should exist directly in `do-work/archive/` after cleanup.** Structure is not the only property worth checking: a correctly-placed REQ can still be a 0-byte file, which is what Pass 6 looks for.

## Commit (Git repos only)

After all passes complete, if any files were moved, consolidated, repointed, or swept from `do-work/runs/`, commit the structural changes.

Check for git with `git rev-parse --git-dir 2>/dev/null`. If not a git repo, skip.

```bash
# Stage all paths affected by cleanup (moves show as delete + add)
# Include do-work/queue/ and working/ if Pass 0 swept any finished REQs
git add do-work/archive/ do-work/user-requests/
# If Pass 0 swept REQs from do-work/queue/ or working/, also stage those paths:
# git add do-work/queue/REQ-NNN-*.md do-work/working/REQ-NNN-*.md  (the deletion side of the moves)
# If Pass 3a found misplaced directories, also stage those paths:
# git add exp/g3-segment-anything/do-work/  (the deletion side of the move)
# If the repoint step rewrote doc links, also stage each rewritten doc file:
# git add <project-root>/docs/prime-foo.md <project-root>/docs/prime-bar.md  (so the repoint lands in the same commit as the moves it repairs)
# If Pass 4 deleted a tracked consumed run, stage only that exact deletion prefix:
# git add -u -- do-work/runs/code-review-2026-07-13-143012/
# Repeat for each swept run. `-u` stages tracked modifications/deletions only,
# so neighboring live or untracked runs cannot be pulled into the commit.
# If Pass 6 restored any blanked REQs, stage each restored file by its exact path:
# git add do-work/archive/UR-010/REQ-1282-slug.md  (one path per restored file)

git commit -m "$(cat <<'EOF'
do-work: cleanup — consolidated {N} REQs, closed {M} URs

- Archived: {list of UR-NNN closed}
- Consolidated: {X} loose REQs into UR folders
- Legacy: {Y} items moved to archive/legacy/
- Fixed: {Z} misplaced folders
- Repointed: {W} doc links
- Swept runs: {R} consumed run directories
- Restored: {S} blanked archived REQs

EOF
)"
```

**Format:** `do-work: cleanup — consolidated {N} REQs, closed {M} URs` — adjust the counts and bullet list to reflect what actually changed. Omit bullet categories where the count is zero.

If nothing was moved, rewritten, restored, or swept (archive and run scratch were already clean), skip the commit entirely.

Stage only paths within `do-work/archive/`, `do-work/user-requests/`, any `do-work/queue/` or working/ REQs swept by Pass 0, any misplaced `do-work/` directories relocated by Pass 3a, the specific doc files rewritten by the repoint step, exact consumed-run deletion prefixes from Pass 4 via `git add -u -- <deleted-run-path>`, and the exact path of each file Pass 6 restored — never a blanket `git add -A`/`.` (see `actions/commit.md` § Rules for the full staging/hook guard).

**Pass 5 stages nothing.** Its objects are git worktrees and branches — refs and directories outside the index — so removing them produces no staged change and never, on its own, makes a commit necessary.

## What This Action Does NOT Do

- Delete work items — only consumed run scratch (`Status: consumed`) is deleted outright; URs, REQs, and other durable artifacts are moved. Pass 5 removes orphaned `worktree-agent-*` worktrees and branches, but mechanically only when they are clean, merged, and their exact REQ is settled outside `do-work/working/`; every other case needs the user's explicit consent
- Modify file contents or frontmatter — files are relocated as-is. Exceptions: Pass 0 normalizes non-standard terminal statuses (`done` → `completed`, etc.) in frontmatter before moving; the Repoint Documentation Links step rewrites link targets in docs that reference moved files; and Pass 6, with the user's consent, rewrites a blanked file's whole content from git history and re-applies its `commit:` field through `tools/checks/record-commit-hash.sh`.
- Touch **active** files in `do-work/queue/` (the queue) or `do-work/working/` — `pending`, `pending-answers`, `blocked`, and `claimed` REQs are actions/work.md's responsibility. Exceptions: Pass 0 sweeps REQs with terminal statuses (`completed`, `done`, `failed`, etc.) from `do-work/queue/` and working/ to archive — that's recovering stranded finished work, not queue processing. Pass 3a relocates queue and working items from **misplaced** `do-work/` trees (created in the wrong directory) back to the canonical root — that's error recovery. Pass 5 reads exact REQ identity and location as deletion evidence but never changes a REQ file; finding the REQ in `working/` forbids automatic removal.
- Archive UR folders that still have pending/in-progress REQs
- Process any REQ files (use actions/work.md for that)

## Common Rationalizations

Guard against these during cleanup:

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This REQ is probably done" | Check the actual status in frontmatter and verify against git history | Premature archival loses in-progress work |
| "Close enough to completed — archive it" | Only archive REQs with terminal status (completed, failed, cancelled) | Non-terminal REQs belong in the queue, not the archive |
| "This UR folder looks empty, delete it" | Check if REQs reference it via `user_request` field | Empty UR folders may have REQs still in the queue or working/ |
| "The archive structure is fine, skip reorganization" | Run all 7 passes even if the archive looks clean | Loose files, consumed run scratch, orphaned worktrees, and blanked REQ content accumulate independently — any of them can need cleanup, and a tidy archive of 0-byte files looks clean to every pass but Pass 6 |
| "That archived REQ is empty, so it must have been an empty REQ" | Run Pass 6's scanner before believing it | A 0-byte archived REQ is the signature of an unguarded `commit:` write-back, not of an empty request — six real REQs of 8–26 KB were lost that way, and their content is recoverable only until `git gc` runs |
| "This `worktree-agent-*` branch is ancient, clean, and merged, so it is obviously abandoned" | Require its exact REQ to be readable, unambiguous, and outside `do-work/working/`; otherwise ask (Pass 5) | Age, cleanliness, and merged-ness do not prove the run finished. A live REQ can have a clean merged builder checkout before its next source edit, and forced deletion discards its execution lane. |

## Red Flags

- REQ with terminal status (completed/failed/cancelled) still in `do-work/queue/` or `do-work/working/`
- UR archived but some of its REQs still pending in the queue — Pass 1 keyed on the UR's `requests:` array (a capture-time record) instead of scanning `user_request:` frontmatter, so review-spawned and addendum follow-ups were invisible to it
- Duplicate REQs found in multiple locations (queue + archive, or working + archive)
- UR folder in archive with no REQ files inside
- A UR whose REQs are all `completed-with-issues` never closes (stays in `user-requests/`) — Pass 1 is filtering on the literal `completed` instead of the terminal-resolved set (`completed`, `completed-with-issues`, or `cancelled`; see `actions/work-reference.md`)
- A UR held open forever by a `cancelled` REQ — same bug class: `cancelled` is terminally resolved and must count toward UR closure
- A moved file still referenced by its old path in tracked markdown after cleanup — the repoint step was skipped or missed a referrer
- `worktree-agent-*` branches or worktrees still present after cleanup ran with no accompanying report line — Pass 5 was skipped entirely (leftovers accumulate silently; the observed case was orphan branches sitting for over half an hour after their run died)
- A 0-byte or frontmatter-less REQ in `do-work/archive/` after cleanup reported a clean archive — Pass 6 was skipped, or it ran non-interactively and its report was not read
- A restored REQ whose `commit:` field was hand-edited rather than written by `tools/checks/record-commit-hash.sh` — the repair bypassed the very guard that exists because free-form edits at this step caused the damage

## Verification Checklist

- [ ] All 7 cleanup passes attempted (Passes 0–6; Pass 3 includes 3a and 3b)
- [ ] No terminal-status REQs remain in `do-work/queue/` or `do-work/working/`
- [ ] Every archived REQ with `user_request` field is inside its UR folder
- [ ] No empty UR folders remain in archive (unless REQs are still pending elsewhere)
- [ ] Every UR Pass 1 closed had its membership derived from the `user_request:` frontmatter scan across `do-work/queue/`, `do-work/working/`, `do-work/archive/` root, and `do-work/archive/UR-NNN/` — never from `input.md`'s `requests:` array
- [ ] Every moved file's old path greps to zero hits in tracked markdown outside `do-work/`
- [ ] Every `do-work/runs/` directory deleted by Pass 4 had `Status: consumed`; `in-progress`, `synthesized`, legacy `complete`, and missing-manifest runs remain
- [ ] Every tracked consumed-run deletion is staged by its exact path
- [ ] Every `worktree-agent-*` leftover Pass 5 removed automatically was clean, merged, and tied to one exact readable REQ outside `do-work/working/`, then came off a successful `git worktree remove` + `git branch -d` (never `-D`/`--force`); every other case was left in place and reported or discarded only after the user named it
- [ ] Pass 6's scanner was run, and every file it restored was approved by the user first, written by `--restore` (not by hand), and staged by its exact path; every unrecoverable file was reported as a permanent loss
