# Cleanup

Consolidates the archive — moves loose files into the right places, closes completed URs, organizes legacy items. Runs automatically at the end of every work loop, or manually on demand.

The canonical implementation is `<skill-root>/tools/do-work-cli.sh --repo-root <project-root> cleanup`. It plans and applies the same ordered policy without an LLM. Use `--dry-run` to inspect exact changes, `--format json` before the command for structured evidence, and `--commit` to create one guarded exact-path commit. Use an exact consent token for a blanked record or an explicit do-work builder discard. Run `--discard-worktree` separately from `--commit`.

## What it does

Seven passes, in order (matching `actions/cleanup.md`):

### Pass 0: Sweep finished queue items
Moves terminal-status REQs (canonical statuses and documented aliases, including `complete`/`done`/`finished`/`closed` → `completed`) from `do-work/queue/` and `working/` into `archive/`.

### Pass 1: Close completed User Requests
When every REQ for a UR is resolved (`completed`, `completed-with-issues`, or `cancelled`), moves the entire UR folder from `user-requests/` into `archive/UR-NNN/`. Being archived isn't enough — a `failed` REQ holds its UR open until it is cancelled with `do-work abandon REQ-NNN` (or `do-work abandon UR-NNN`, which reaches the UR's still-open members — including that `failed` one — in a single itemized confirmation). A completed follow-up REQ does the recovery work but never flips the original out of `failed`, so cancelling the failed REQ is the only thing that closes the UR — do it once the follow-up has done the work, or when no follow-up is wanted.

### Pass 2: Consolidate loose REQ files
Moves REQ files sitting in `archive/` root into their UR folders (`archive/UR-NNN/`). REQs without a UR reference go to `archive/legacy/`.

### Pass 3: Fix misplaced directories
Detects `do-work/` directories accidentally created in subdirectories and relocates them. Catches UR folders nested under `archive/user-requests/` and moves them up.

### Pass 4: Sweep consumed run scratch
Deletes only `do-work/runs/*/` directories whose root `manifest.md` says `Status: consumed`. In-progress runs remain resumable, while `synthesized` and legacy `complete` runs are preserved because their assembled output may not have reached the user yet. Tracked deletions are staged by their exact run path.

An entirely untracked consumed run is spent scratch but has no Git rollback source. Cleanup therefore revalidates its exact regular-file inventory and consumed manifest immediately before deletion and labels the result as a non-rollback spent-scratch deletion. Mixed or durable dirty targets remain refused.

### Pass 5: Remove orphaned worktrees
Checks every linked checkout, including T3 and detached worktrees, plus leftover local `t3/*` and `worktree-agent-*` branches. Cleanup removes a clean candidate when merging its tip would leave the files on local `main` unchanged. This includes squash merges. Repositories without `main` use the primary checkout’s branch. The primary checkout, invoking checkout, locked worktrees, local edits and conflicting changes are protected. Cleanup also protects named do-work lanes with an active or ambiguous request. Use `cleanup --dry-run` to see each planned removal and retained candidate.

### Pass 6: Restore blanked archived REQs
Finds archived REQ and UR files that have lost their content — 0 bytes, or no frontmatter left — and offers to write it back from git history. That damage is the signature of an unguarded `commit:` write-back: the file was complete at its implementation commit and the metadata commit right after it replaced the whole thing with nothing. Cleanup shows you a dry run first (which file, which commit the content comes from, how many bytes, which hash goes back into `commit:`), then asks before writing anything; running unattended, it only reports. Recovery has a deadline — the lost content survives in git only until a `git gc` collects it — so a file reported here is worth acting on. Restoring overwrites the file, so anything written to it since it was blanked is discarded.

### After all passes: repoint doc links
Every move above changes a file's path. Cleanup tracks each old → new path and rewrites links in the repo's tracked markdown outside `do-work/` that pointed at the moved file (e.g. a prime doc's Lessons link to an archived REQ), so consolidation doesn't leave broken links behind. The summary reports `Repointed: N doc links in M files` (or `Repointed: none`).

## Result

```
do-work/
├── archive/
│   ├── UR-001/              # Self-contained: input.md + completed REQs
│   │   ├── input.md
│   │   ├── REQ-001-done.md
│   │   └── REQ-002-done.md
│   ├── UR-002/
│   └── legacy/              # Standalone REQs without UR references
│       └── REQ-010-done.md
├── user-requests/            # Only open URs remain here
├── queue/                    # Only active queue items remain here
│   └── (pending REQs)
```

## Key rules

- Preserves work items. Pass 4 sweeps consumed run scratch, and Pass 5 removes stale worktrees and builder branches.
- No content modification except normalizing non-standard statuses (`done` → `completed`), repointing doc links to moved files, and — only with your say-so — restoring a blanked archived REQ from git history
- Skips active queue items (`pending`, `claimed`)

## Usage

```
do-work cleanup
do-work tidy
do-work consolidate
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> cleanup --dry-run
```
