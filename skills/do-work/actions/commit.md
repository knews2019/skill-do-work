# Commit Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to commit accumulated uncommitted files. Analyzes changes, associates them with existing REQs for traceability, groups the rest semantically, and commits everything in small atomic batches. User-facing walkthrough: [`docs/commit-guide.md`](../docs/commit-guide.md).

Unlike the commit steps embedded in other actions (capture Step 7, work.md's Commit Phase, review-work standalone, cleanup), this action handles files that accumulated outside the normal pipeline — manual edits, ad-hoc fixes, or work done between do-work runs.

**Commit pathway deconfliction:** `actions/work.md` delegates its complete lifecycle, archive, release, exact commit, and provenance tail to canonical finalization. Standalone review may commit its own appended report; this action handles only leftover files discovered through `git status`. It never replaces, retries, or reconstructs a lifecycle transaction.

## When to Use

**Use when:**
- User wants to commit accumulated uncommitted files with REQ tracing
- User says "commit", "commit changes", "save changes", or "save work"
- Files have accumulated outside the normal pipeline (manual edits, ad-hoc fixes)

**Do NOT use when:**
- User just wants to *understand* uncommitted changes — route to `../../do-work-toolbox/actions/inspect.md` instead
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

Before protected association or reading any uncommitted path, invoke the canonical recovery composition:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json recover
```

Continue only on typed `success`, and read the ordered `finalizations` one record at a time: a record carrying `FINALIZATION-SET-ASIDE` in its `reason_codes` is one REQ recovery excluded, and the remaining records still count as settled when their `blocked_paths` and `reason_codes` are empty. Recovery commits all safe finalization groups in its returned order, preserves unfinished working claims without authority, and refuses before this action selects ordinary groups when staged, protected, shared, or multiply-owned evidence remains ambiguous. Group only the changes left after recovery; never re-associate a recovered path.

Start the protected inventory wrapper; it owns the worktree-safe run quarantine and delegates low-level classification to the existing checks:

```bash
<skill-root>/scripts/protected-inventory.sh start
```

It gates on `git rev-parse --git-dir`, enumerates every uncommitted path, and prints one `<tag>\t<path>` row per file:

- **M** — modified (a renamed path is tagged M too: read its diff, don't re-read it as new content)
- **A** — added, covering both staged-new and untracked
- **D** — deleted
- **X** — non-deleted excluded path: secret-shaped, secret-derived, or an ambiguous addition beside `X`/`XD`; fully excluded
- **XD** — deleted secret-shaped path; eligible for a deletion-only commit

Secret-shaped matching is case-insensitive and applies to the basename only: `.env*` or `*.env`, `*credentials*`, `*.pem`/`*.key`/`*.p12`/`*.pfx`, and `*secret*`. Thus `.ENV`, `AuthCredentials.json`, `private.PEM`, and `UPPER-SECRET.txt` are all secret-shaped, while an ordinary file beneath a directory named `secrets/` is not classified from the directory name alone.

Exit 1 means the working tree is clean — report "Nothing to commit" and exit. Exit 2 means this is not a git repo — report and exit.

**`X` rows are reported, never skipped silently.** Carry them into the final report so the user knows a secret-shaped, secret-derived, or ambiguity-quarantined file is sitting uncommitted; never read, diff, stage, or commit them.

Start a run-level quarantine with the first inventory. The deterministic file comes from `git rev-parse --git-path`, so every later command block re-derives the same worktree-safe location per the canonical [State across command blocks](../docs/prescribed-shell-primitives.md#state-across-command-blocks) rule. Before consuming any inventory, add every current `X` path to that set and overlay the full set on the current rows, changing any matching M/A path back to `X`. Apply this before every later read, REQ association, or staging decision in Steps 2, 3, and 5. The invariant is **once `X`, always `X` for this action run**: resetting a pre-staged rename destination must not make its contents readable when the next inventory can see only an addition.

If an X path was already staged before this action began, stop before making any commit. Remove it from the index without reading it (for example, `git reset -- <path>`), then re-run the inventory using the retained quarantine — re-derive its Git-private path, do not truncate it, append the new X rows, and overlay the full set before consuming the result. The action must not let a pre-staged X path ride along with an otherwise safe group.

**`XD` rows are also reported, but only the deletion may proceed.** Never read, diff, reconstruct, or retrieve the former contents from the index, a commit, or any other Git history. The path and its current deletion state are the complete inspection surface.

If the script is missing or will not run, do it by hand from the complete NUL-delimited output of `git -c status.renames=copies status --porcelain=v1 --untracked-files=all -z`; the command-line copy setting is required even when repository configuration disables rename detection. First classify every path as M/A/D, including each rename or copy record's second NUL-delimited origin path, then lowercase its basename with `tr '[:upper:]' '[:lower:]'` and apply the patterns above. Change a non-deleted secret-shaped tag to X and a deleted one to XD; a secret-shaped rename origin is XD and its destination is X, while a destination copied from a secret-shaped origin is X without an XD source row. Buffer the complete classification before using it. If the finished inventory contains any X or XD and any remaining A, provenance is ambiguous — Git cannot identify a copy when both source and destination are untracked — so change every A to X. Finally add all X paths to the run-level quarantine and overlay that quarantine before any row is consumed. The `-uall` flag is not optional; preserve the canonical [Per-file untracked inventory](../docs/prescribed-shell-primitives.md#per-file-untracked-inventory) contract or files beneath a brand-new directory can escape the exclusion scan. That is a secret-leak path, and `../../do-work-toolbox/actions/stray-check.md`'s Red Flags record that it has been hit.

### Step 2: Read Changes

Build a semantic understanding of each uncommitted file:

- **Modified files**: Read the `git diff` for each file. Understand what changed and why.
- **New/untracked files**: Read the file contents. Skip binary files (detect by extension: images, compiled assets, archives). For large files (>500 lines), read the first 100 lines and last 50 lines to understand purpose.
- **Deleted files (`D`)**: Note the path and what the file likely was (infer from path and name).
- **Deleted secret-shaped files (`XD`)**: Note only the path and that it is deleted. Do not run `git diff`, `git show`, `git log -p`, or any equivalent that reads, reconstructs, or displays former contents.

The goal is to understand each file well enough to group it with related changes and write a meaningful commit message.

### Step 3: Associate with REQs

Feed the current protected M/A/D/XD paths into association, excluding every path quarantined as `X` during the run:

```bash
<skill-root>/scripts/protected-inventory.sh associate
```

The wrapper re-derives the repository root and moves paths through files rather than interpolating them into shell source. It appends the new X rows to Step 1's quarantine before filtering, so both current X rows and paths excluded by an earlier inventory stay out. M/A/D/XD participate in association only when the path has never been X. The delegated check scans `do-work/archive/**/REQ-*.md` and `do-work/working/REQ-*.md`, reads each REQ's `## Implementation Summary` file list, and prints one `<owner>\t<path>` row per candidate — a `REQ-NNN` id, or `-` for unassociated. Exit 1 means there were no candidates other than X; continue with the reported X rows only. Exit 2 means a usage error or no `do-work/` directory; skip REQ tracing and send the remaining safe M/A/D/XD files to Step 4.

What the script settles, so this prose no longer has to:

- **Terminal-success matching honors the Schema Read Contract's aliases**, so `completed`, `completed-with-issues`, and `complete`/`done`/`finished`/`closed` all qualify. Testing only for the literal `completed` is the bug in the Red Flags below — it drops every remediated-with-issues REQ, and its files then never get associated.
- **In-flight `working/` REQs are included** regardless of status, since a claimed REQ is never terminal.
- **Conflict resolution:** a path claimed by two REQs goes to the one with the latest `completed_at`. An archived REQ outranks an in-flight one.
- **`do-work/` metadata paths are excluded** from association, matching `tools/checks/scope-drift.sh`.

**Partial matches count.** If 3 out of 5 files in a REQ's Implementation Summary are among the uncommitted files, group all 3 under that REQ.

Files that come back `-` remain unassociated and move to Step 4.

If the script is missing or will not run, do it by hand: glob both directories, read each REQ's `status` (accepting every alias above) and `## Implementation Summary` list, path-match, and tie-break on the latest `completed_at`.

### Step 4: Group Unassociated Files

Cluster the remaining files into semantic groups of 1-5 files each:

1. **Use the safe evidence from Step 2** for each unassociated file — diffs/contents for M/A, path metadata for D/XD
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

For each XD path, re-run the protected association immediately before staging and proceed only while the exact row is still XD and the path has never been X. Then invoke the exact-deletion guard, which checks cached name/status without reading content:

```bash
<skill-root>/scripts/stage-exact-deletion.sh "$path"
```

This accepts a path already staged as one exact deletion without `git add -u`; otherwise it stages the tracked deletion and verifies the same exact cached metadata afterward. Never inspect a staged diff. If the path is no longer XD, stop and reclassify it. X paths are never staged by any command, and no commit proceeds while an X path remains pre-staged.

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

**Rules:** see `## Rules` below for the staging/hook guard. One commit per group — keep them atomic. List each file in the commit body with its action (Modified, Added, Deleted).

Committing an XD deletion removes the path from the new tree only. It does **not** erase the secret from Git history and does **not** rotate or revoke the credential; report those as separate remediation needs when applicable.

### Step 6: Report

Remove the run's Git-private quarantine file before reporting: `rm -f "$(git rev-parse --git-path do-work-commit-secret-quarantine)"`. Do this on every exit after Step 1; the next run truncates the same deterministic path defensively, but a completed action leaves no scratch state behind.

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

Report committed or still-pending XD paths separately as **secret-shaped deletions (contents not inspected)** and repeat that a deletion does not clean Git history or rotate credentials.

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
| Pre-commit hook failure | Fix the underlying issue, re-stage, and retry as a **new** commit (see `## Rules` — never `--no-verify`). |
| File matches multiple REQs | Associate with the most recently completed REQ (`completed_at` timestamp) |
| Ambiguous semantic grouping | Prefer smaller groups (1-2 files) over larger uncertain groups |
| Binary files in untracked | Skip reading contents, group by directory proximity and filename |
| Very large number of files (50+) | Process normally but warn the user: "Large changeset — {N} files across {M} commits. Review the commit log." |
| All paths are X | Report the exclusions clearly, commit nothing |

## What This Action Does NOT Do

- Create REQ files — it only traces back to existing archived REQs
- Modify archived REQ files — `Traced-to:` is in the commit message only, not written to the REQ
- Push to remote — only creates local commits
- Handle interactive staging (`git add -p`) — it commits complete files
- Replace the commit steps in other actions — those remain for their specific pipelines
- Stage or commit secret file contents — X is always excluded; XD permits only the deletion entry

## Rules

**Canonical statement of the commit-staging/hook guard** — every other action in this skill that stages or commits files points here rather than restating it:

- **Never `git add -A` or `git add .`.** Stage only the specific files that belong to the commit you're making. A blanket add risks sweeping in secrets, `.env` files, or unrelated in-progress changes from other work — the whole point of grouping files by REQ/semantic purpose (Steps 3-4) is defeated if staging ignores those groups.
- **X and XD have different boundaries.** Never read, diff, stage, or commit X. For XD, stage only the verified deletion with the Step 5 procedure; never read, diff, or reconstruct its former contents.
- **Never bypass a failing pre-commit hook** with `--no-verify` (or signing with `--no-gpg-sign`). Fix the underlying issue, re-stage, and retry as a **new** commit — never amend past a hook failure.

## Red Flags

- An X path staged for commit, or an XD path whose staged state is anything other than deletion
- Commit message has no REQ reference when matching REQs exist in the system
- Files from multiple unrelated REQs grouped in a single commit
- Uncommitted files belonging to a terminal-success REQ aren't associated to it — Step 3 is filtering on the literal `completed` instead of the full success set (`completed`, `completed-with-issues`, plus `complete`/`done`/`finished`/`closed` aliases; see `actions/work-reference.md`)

## Verification Checklist

- [ ] Every commit traces to a REQ or a clear semantic group
- [ ] No secret contents committed; every X path stayed excluded and every committed XD path contributed only a deletion
- [ ] Commit messages follow the established format
- [ ] Each commit is atomic — one logical change per commit
- [ ] All excluded files reported to the user with reason
