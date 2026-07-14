# Cleanup

Consolidates the archive — moves loose files into the right places, closes completed URs, organizes legacy items. Runs automatically at the end of every work loop, or manually on demand.

## What it does

Five passes, in order (matching `actions/cleanup.md`):

### Pass 0: Sweep finished queue items
Moves terminal-status REQs (`completed`, `completed-with-issues`, `failed`, `cancelled`, plus normalized aliases like `done`/`finished`/`closed`/`abandoned`/`wont-do`) from `do-work/queue/` and `working/` into `archive/`.

### Pass 1: Close completed User Requests
When all REQs for a UR are archived, moves the entire UR folder from `user-requests/` into `archive/UR-NNN/`.

### Pass 2: Consolidate loose REQ files
Moves REQ files sitting in `archive/` root into their UR folders (`archive/UR-NNN/`). REQs without a UR reference go to `archive/legacy/`.

### Pass 3: Fix misplaced directories
Detects `do-work/` directories accidentally created in subdirectories and relocates them. Catches UR folders nested under `archive/user-requests/` and moves them up.

### Pass 4: Sweep consumed run scratch
Deletes only `do-work/runs/*/` directories whose root `manifest.md` says `Status: consumed`. In-progress runs remain resumable, while `synthesized` and legacy `complete` runs are preserved because their assembled output may not have reached the user yet. Tracked deletions are staged by their exact run path.

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

- Deletes no work items — only run scratch explicitly marked `Status: consumed`
- No content modification except normalizing non-standard statuses (`done` → `completed`) and repointing doc links to moved files
- Skips active queue items (`pending`, `claimed`)

## Usage

```
do-work cleanup
do-work tidy
do-work consolidate
```
