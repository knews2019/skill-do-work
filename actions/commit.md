# Commit Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to commit accumulated uncommitted files. Analyzes changes, associates them with existing REQs for traceability, groups the rest semantically, and commits everything in small atomic batches. User-facing walkthrough: [`docs/commit-guide.md`](../docs/commit-guide.md).

Unlike the commit steps embedded in other actions (capture Step 7, work.md's Commit Phase, review-work standalone, cleanup), this action handles files that accumulated outside the normal pipeline — manual edits, ad-hoc fixes, or work done between do-work runs.

**Commit pathway deconfliction:** Three actions can commit archived REQs: (1) actions/work.md's Commit Phase commits the REQ + implementation after completion, (2) review-work standalone commits the REQ after appending a Review section, (3) this action commits leftover files traced to archived REQs. This action only discovers files via `git status` — if work or review-work already committed a file, it won't appear here. No double-commit risk exists as long as the prior actions committed cleanly. If a prior commit was interrupted, this action may pick up the leftovers — that's the intended behavior.

## When to Use

**Use when:**
- User wants to commit accumulated uncommitted files with REQ tracing
- User says "commit", "commit changes", "save changes", or "save work"
- Files have accumulated outside the normal pipeline (manual edits, ad-hoc fixes)

**Do NOT use when:**
- User just wants to *understand* uncommitted changes — route to actions/inspect.md instead
- Committing as part of actions/work.md (work.md has its own commit step)

## When This Runs

- **Manually** when the user invokes it (e.g., `do-work commit`, `do-work save work`)

## Steps

```
commit action
  │
  ├── Preflight ── not a git repo? → exit
  │                 clean tree? → "Nothing to commit" → exit
  │
  ├── Read Changes ── diffs for modified, contents for new, paths for deleted
  │
  ├── Associate with REQs ── match files to archived REQ Implementation Summaries
  │
  ├── Group Unassociated ── semantic clustering (1-5 files per group)
  │
  ├── Commit ── REQ-linked groups first, then unassociated groups
  │
  └── Report ── summary table of all commits
```

### Step 1: Preflight

Check for git with `git rev-parse --git-dir 2>/dev/null`. If not a git repo, report and exit.

Run `git status --porcelain --untracked-files=all` to get all uncommitted changes — staged, unstaged, and untracked. The `--untracked-files=all` (`-uall`) flag matters: plain `git status --porcelain` collapses a wholly-untracked directory into a single `?? dir/` row, so Step 2 would never see (or would try to "read") the files inside a new untracked folder. With `-uall`, every untracked file is listed individually.

If the working tree is clean, report "Nothing to commit" and exit.

Categorize each file by its status:
- **Modified** (M) — existing files with changes
- **Added** (??, A) — new or untracked files
- **Deleted** (D) — removed files

**Exclude dangerous files** from all subsequent steps:
- `.env`, `.env.*` — environment variables
- `credentials.*`, `*credentials*` — credential files
- `*.pem`, `*.key`, `*.p12`, `*.pfx` — certificates and keys
- `*.secret`, `*secret*` — secret files

If any files are excluded, collect them for the final report. Do not silently skip them — the user needs to know.

### Step 2: Read Changes

Build a semantic understanding of each uncommitted file:

- **Modified files**: Read the `git diff` for each file. Understand what changed and why.
- **New/untracked files**: Read the file contents. Skip binary files (detect by extension: images, compiled assets, archives). For large files (>500 lines), read the first 100 lines and last 50 lines to understand purpose.
- **Deleted files**: Note the path and what the file likely was (infer from path and name).

The goal is to understand each file well enough to group it with related changes and write a meaningful commit message.

### Step 3: Associate with REQs

Scan `do-work/archive/` for completed REQs that might own some of the uncommitted files:

1. Glob for `do-work/archive/**/REQ-*.md` — find all archived REQs
2. For each archived REQ:
   - Read the frontmatter — check for `commit:` field and a terminal-success `status` (`completed` or `completed-with-issues` — see `actions/work-reference.md`'s Terminal-success status set)
   - Read the `## Implementation Summary` section — extract the list of files created/modified
3. Also check `do-work/working/` for in-flight REQs with file lists

Match uncommitted files against these file lists by path. A file is associated with a REQ if it appears in that REQ's Implementation Summary (created, modified, or referenced).

**Conflict resolution:** If a file matches multiple REQs, associate it with the most recently completed one (latest `completed_at` timestamp).

**Partial matches count.** If 3 out of 5 files in a REQ's Implementation Summary are among the uncommitted files, group all 3 under that REQ.

Files that don't match any REQ remain unassociated and move to Step 4.

### Step 4: Group Unassociated Files

Cluster the remaining files into semantic groups of 1-5 files each:

1. **Read the diffs/contents** from Step 2 for each unassociated file
2. **Identify logical changes** — files that work together toward a single purpose:
   - A component and its test file
   - Multiple files in the same module touching the same feature
   - Config file changes that go together
   - Documentation updates related to the same topic
3. **Use directory proximity as a secondary signal** — files in the same directory are more likely related, but don't group unrelated changes just because they're neighbors
4. **Assign a short descriptive label** to each group (e.g., "API client error handling", "Test coverage for auth module", "Config and dependency updates")

**When uncertain, prefer smaller groups.** Two commits of 2 files each is better than one commit of 4 loosely-related files.

**Single-file groups are fine.** A standalone change that doesn't relate to anything else gets its own commit.

### Step 5: Commit

Commit each group in order — REQ-associated groups first, then unassociated groups.

**REQ-associated commits** (one per REQ):

```bash
git add src/specific-file.ts src/other-file.ts

git commit -m "$(cat <<'EOF'
[REQ-NNN] {REQ title} — additional changes

Traced-to: do-work/archive/UR-NNN/REQ-NNN-slug.md

- Modified src/specific-file.ts
- Added src/other-file.ts

EOF
)"
```

**Format:** `[{REQ id}] {REQ title} — additional changes` + `Traced-to:` line pointing to the archived REQ + file list bullets. Note: this format intentionally differs from actions/work.md's primary commit format (`[{id}] {title} (Route {route})` + `Implements:`). The `— additional changes` suffix and `Traced-to:` prefix signal these are supplementary commits for files that missed the original work commit, not the primary implementation commit.

**Unassociated commits** (one per semantic group):

```bash
git add src/specific-file.ts src/other-file.ts

git commit -m "$(cat <<'EOF'
[do-work] {descriptive label}

- Modified src/specific-file.ts
- Added src/other-file.ts

EOF
)"
```

**Format:** `[do-work] {descriptive label}` + file list bullets.

**Rules:**
- Stage specific files per group — never `git add -A` or `git add .`
- Do not bypass pre-commit hooks — fix issues and retry
- One commit per group — keep them atomic
- List each file in the commit body with its action (Modified, Added, Deleted)

### Step 6: Report

Print a summary of all commits:

```
Committed {N} groups ({M} files):
  abc1234  [REQ-003] Dark Mode — additional changes (3 files)
  def5678  [do-work] API client error handling (2 files)
  ghi9012  [do-work] Test coverage for auth module (4 files)
  jkl3456  [do-work] Config and dependency updates (5 files)
```

If files were excluded:

```
Excluded (potential secrets):
  .env.local — skipped
  credentials.json — skipped
```

If nothing was committed (all files were excluded):

```
No files committed. All uncommitted files matched exclusion patterns.

Excluded:
  .env.local — potential secrets
```

## Error Handling

| Situation | Action |
|-----------|--------|
| Not a git repo | Report "Not a git repository" and exit |
| Clean working tree | Report "Nothing to commit" and exit |
| Pre-commit hook failure | Fix the underlying issue, re-stage, and retry as a **new** commit. Do NOT use `--no-verify` to skip hooks — fix the root cause. |
| File matches multiple REQs | Associate with the most recently completed REQ (`completed_at` timestamp) |
| Ambiguous semantic grouping | Prefer smaller groups (1-2 files) over larger uncertain groups |
| Binary files in untracked | Skip reading contents, group by directory proximity and filename |
| Very large number of files (50+) | Process normally but warn the user: "Large changeset — {N} files across {M} commits. Review the commit log." |
| All files excluded | Report the exclusions clearly, commit nothing |

## What This Action Does NOT Do

- Create REQ files — it only traces back to existing archived REQs
- Modify archived REQ files — `Traced-to:` is in the commit message only, not written to the REQ
- Push to remote — only creates local commits
- Handle interactive staging (`git add -p`) — it commits complete files
- Replace the commit steps in other actions — those remain for their specific pipelines
- Stage `.env`, credentials, keys, or other secret files — these are always excluded

## Common Rationalizations

Guard against these when committing:

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "No REQ matches — just commit everything together" | Group unassociated files by semantic purpose (feature, fix, config, etc.) | Even outside the pipeline, commits should be atomic and meaningful |
| "The commit message doesn't need a REQ reference" | Include REQ reference when REQs exist — it's the traceability link | Without REQ references, the trail of intent is broken |

## Red Flags

- `.env`, credentials, or secret files staged for commit
- Single commit with >20 files (likely needs splitting)
- Commit message has no REQ reference when matching REQs exist in the system
- Files from multiple unrelated REQs grouped in a single commit
- Uncommitted files belonging to a `completed-with-issues` REQ aren't associated to it — Step 3 is filtering on the literal `completed` instead of the terminal-success set (`completed` or `completed-with-issues`; see `actions/work-reference.md`)

## Verification Checklist

- [ ] Every commit traces to a REQ or a clear semantic group
- [ ] No credential or secret files committed (.env, *.key, *.pem, credentials.*)
- [ ] Commit messages follow the established format
- [ ] Each commit is atomic — one logical change per commit
- [ ] All excluded files reported to the user with reason
