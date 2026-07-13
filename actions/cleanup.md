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
- User wants to *delete* or discard work — cleanup only reorganizes work items (URs, REQs), never deletes them. (The lone exception is Pass 4, which sweeps *consumed* run scratch — a `Status: complete` directory under `do-work/runs/` — after its findings have been promoted. That is spent scratch, not work.)

## When This Runs

- **Automatically** at the end of every work loop (after all pending REQs are processed)
- **Manually** when the user invokes it (e.g., `do-work cleanup`, `do-work tidy`)

## Steps

Five passes, in order:

### Pass 0: Sweep Finished Queue Items

Scan `do-work/queue/` and the working directory for REQs with terminal statuses that should have been archived but weren't — typically from manual work, different agents, or legacy sessions that completed outside the standard work pipeline.

1. **Glob `do-work/queue/REQ-*.md`**
2. **Read each REQ's frontmatter** `status` field
3. **If status is any terminal value** — `completed`, `completed-with-issues`, `failed`, `cancelled`, or any non-standard terminal status (`done`, `finished`, `closed`, `canceled`, `abandoned`, `wont-do`):
   - **Normalize non-standard statuses** before moving: change `done` → `completed`, `finished` → `completed`, `closed` → `completed`, `canceled`/`abandoned`/`wont-do` → `cancelled` in frontmatter
   - Move the REQ to `do-work/archive/` root (Pass 1 and Pass 2 will then consolidate it into the correct UR folder)
   - Report: `Swept REQ-NNN from do-work/queue/ (was status: {original}) → archive`
4. **Leave `pending`, `pending-answers`, and `claimed` REQs untouched** — those are active queue items
5. **Also check `do-work/working/`** — if any REQ there has a terminal status (`completed`, `completed-with-issues`, `done`, `finished`, `closed`, `failed`, `cancelled`), it was finished but never moved out. Same treatment: normalize status, move to `do-work/archive/` root, report it.

### Pass 1: Close Completed User Requests

Check `do-work/user-requests/` for UR folders that are ready to archive.

For each UR folder in `do-work/user-requests/`:

1. Read `input.md` and parse the `requests` array from frontmatter (e.g., `[REQ-044, REQ-045, REQ-046]`)
2. For each REQ ID in the array, check if it exists with a **terminal-resolved status** (`completed`, `completed-with-issues`, or `cancelled` — see `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set) in ANY of these locations:
   - `do-work/archive/UR-NNN/` (already consolidated)
   - `do-work/archive/` root (loose in archive)
   If the same REQ-ID is found in **both** locations simultaneously, flag it and leave the UR in `user-requests/` untouched: `⚠ Duplicate: REQ-NNN found in both archive/ root and archive/UR-NNN/. Resolve manually, then re-run cleanup.`
3. If **ALL** REQs are terminally resolved — `completed`, `completed-with-issues`, or `cancelled` (and no duplicates flagged):
   - Gather any loose completed/cancelled REQ files from `do-work/archive/` root into the UR folder
   - Move the entire UR folder to `do-work/archive/UR-NNN/`
   - Report: `Archived UR-NNN (all N REQs resolved)` — when any were cancelled, say so: `(N-K complete, K cancelled)`
4. If **NOT all** REQs are terminally resolved:
   - Leave the UR folder in `user-requests/` — it's not ready yet
   - Report: `UR-NNN still open (X/Y REQs complete)`

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

Fan-out actions (code-review, deep-explore, multi-REQ work — see `crew-members/background-agents.md`) each delete their own `do-work/runs/<action>-<ts>/` directory once its findings are consumed. This pass is the **safety net** for runs abandoned after they finished but before their owner deleted them — e.g. a session that crashed between synthesis and cleanup.

1. **Glob `do-work/runs/*/`** (each is one run directory).
2. **Read each run's `manifest.md`** and check its `Status:` line.
3. **If `Status: complete`** — the run's findings were already synthesized and promoted; the directory is spent scratch. **Delete it.** Report: `Swept run dir do-work/runs/{name} (Status: complete)`.
4. **If the manifest is missing, or `Status:` is anything other than `complete`** (e.g. `in-progress`) — **leave it untouched** and report: `Left run dir do-work/runs/{name} (incomplete — may be resumable)`. A crashed run with unfinished dimensions is recoverable from its files (see `crew-members/background-agents.md` recovery procedure); never delete it here.

This is the one place cleanup deletes rather than reorganizes, and it is scoped strictly to **consumed run scratch** — a `Status: complete` directory under `do-work/runs/` only. URs, REQs, and every other `do-work/` artifact are still only ever moved, never deleted.

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
  - Repointed: 39 doc links in 5 files
  - Still open: UR-015 (2/4 REQs complete)
```

When files were moved but no referrers were found, still print `Repointed: none` — the line is evidence the repoint step ran.

If nothing needed fixing:
```
Archive is clean. No loose files or pending closures found.
```

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

**No loose REQ or CONTEXT files should exist directly in `do-work/archive/` after cleanup.**

## Commit (Git repos only)

After all passes complete, if any files were moved or consolidated, commit the structural changes.

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
# git add docs/prime-foo.md docs/prime-bar.md  (so the repoint lands in the same commit as the moves it repairs)

git commit -m "$(cat <<'EOF'
do-work: cleanup — consolidated {N} REQs, closed {M} URs

- Archived: {list of UR-NNN closed}
- Consolidated: {X} loose REQs into UR folders
- Legacy: {Y} items moved to archive/legacy/
- Fixed: {Z} misplaced folders
- Repointed: {W} doc links

EOF
)"
```

**Format:** `do-work: cleanup — consolidated {N} REQs, closed {M} URs` — adjust the counts and bullet list to reflect what actually changed. Omit bullet categories where the count is zero.

If nothing was moved (archive was already clean), skip the commit entirely.

Do not use `git add -A` or `git add .` — stage only paths within `do-work/archive/`, `do-work/user-requests/`, any `do-work/queue/` or working/ REQs swept by Pass 0, any misplaced `do-work/` directories relocated by Pass 3a, and the specific doc files rewritten by the repoint step. Don't bypass pre-commit hooks.

## What This Action Does NOT Do

- Delete any files — only moves them into the right location
- Modify file contents or frontmatter — files are relocated as-is. Exceptions: Pass 0 normalizes non-standard terminal statuses (`done` → `completed`, etc.) in frontmatter before moving, and the Repoint Documentation Links step rewrites link targets in docs that reference moved files.
- Touch **active** files in `do-work/queue/` (the queue) or `do-work/working/` — `pending`, `pending-answers`, and `claimed` REQs are actions/work.md's responsibility. Exceptions: Pass 0 sweeps REQs with terminal statuses (`completed`, `done`, `failed`, etc.) from `do-work/queue/` and working/ to archive — that's recovering stranded finished work, not queue processing. Pass 3a relocates queue and working items from **misplaced** `do-work/` trees (created in the wrong directory) back to the canonical root — that's error recovery.
- Archive UR folders that still have pending/in-progress REQs
- Process any REQ files (use actions/work.md for that)

## Common Rationalizations

Guard against these during cleanup:

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This REQ is probably done" | Check the actual status in frontmatter and verify against git history | Premature archival loses in-progress work |
| "Close enough to completed — archive it" | Only archive REQs with terminal status (completed, failed, cancelled) | Non-terminal REQs belong in the queue, not the archive |
| "This UR folder looks empty, delete it" | Check if REQs reference it via `user_request` field | Empty UR folders may have REQs still in the queue or working/ |
| "The archive structure is fine, skip reorganization" | Run all 4 passes even if the archive looks clean | Loose files accumulate gradually — what looks clean may have orphans |

## Red Flags

- REQ with terminal status (completed/failed/cancelled) still in `do-work/queue/` or `do-work/working/`
- UR archived but some of its REQs still pending in the queue
- Duplicate REQs found in multiple locations (queue + archive, or working + archive)
- UR folder in archive with no REQ files inside
- A UR whose REQs are all `completed-with-issues` never closes (stays in `user-requests/`) — Pass 1 is filtering on the literal `completed` instead of the terminal-resolved set (`completed`, `completed-with-issues`, or `cancelled`; see `actions/work-reference.md`)
- A UR held open forever by a `cancelled` REQ — same bug class: `cancelled` is terminally resolved and must count toward UR closure
- A moved file still referenced by its old path in tracked markdown after cleanup — the repoint step was skipped or missed a referrer

## Verification Checklist

- [ ] All 4 consolidation passes attempted
- [ ] No terminal-status REQs remain in `do-work/queue/` or `do-work/working/`
- [ ] Every archived REQ with `user_request` field is inside its UR folder
- [ ] No empty UR folders remain in archive (unless REQs are still pending elsewhere)
- [ ] Every moved file's old path greps to zero hits in tracked markdown outside `do-work/`
