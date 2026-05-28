# Stray-Check Action

> **Part of the do-work skill.** Repo-wide scan for orphan/junk files that pollute where they don't belong — temp/backup files, committed build artifacts, should-be-gitignored files, misplaced/duplicate/empty files, large blobs, AI scratch droppings, and best-effort dead code. Report-only by default; applies fixes only on explicit confirmation. User-facing walkthrough: [`docs/stray-check-guide.md`](../docs/stray-check-guide.md).

A hygiene scanner for the **whole repository**, not do-work's own bookkeeping. The scan phase is strictly read-only: it reports findings grouped by severity, each with a suggested fix. Fixes (deleting untracked junk, `git rm --cached` for tracked artifacts, appending to `.gitignore`) happen **only after the user explicitly confirms** — never silently. It **skips the entire `do-work/` tree**; that's the cleanup and forensics actions' territory.

## When to Use

**Use when:**
- User says "stray files", "orphan files", "junk", "what doesn't belong", "file hygiene", or "stray-check"
- User wants to find files polluting the repo — leftover temp/backup files, committed build output, files that should be gitignored
- A pre-commit or pre-release cleanliness sweep before sharing the repo

**Do NOT use when:**
- The clutter is do-work's own files — loose REQs, misplaced `do-work/` directories → use the **cleanup** action
- The user wants do-work pipeline health — stuck/hollow REQs, orphaned URs → use the **forensics** action
- The user wants to validate a single human-facing artifact (brief, report) → use the **slop-check** action
- The user wants code-quality review of a change → use the **code-review** action

## Input

`$ARGUMENTS` — all optional:

1. **Path scope** (e.g., `src/`, `packages/api/`) — limit the scan to a subtree. Default: repo root.
2. **Mode token:**
   - *(default, no token)* — **report only** (dry run). Print findings and stop; do not offer fixes.
   - `fix` / `--fix` — run the scan, then enter the guarded fix flow (Step 5).
   - `report` / `--report-only` — force pure read-only; suppress the fix prompt even if fixable items exist.

A path and a mode token can be combined (e.g., `do-work stray-check src/ fix`).

## Steps

### Step 1: Scope & Setup

1. Check for git: `git rev-parse --git-dir 2>/dev/null`. If not a git repo, run filesystem-only checks (categories 1, 5, 6, 7, 8, 9, 10) and skip the git-dependent ones (2, 3, 4 in their tracked form).
2. Resolve the path argument to a scan root (default: repo root).
3. Build the **noise skip-list** — `node_modules/`, `vendor/`, `.git/`, `.venv/`, `venv/`, `__pycache__/`, `.cache/`, `.pytest_cache/`, `.mypy_cache/`, plus any untracked path already matched by `.gitignore` (those are *correctly* ignored — not pollution). **The skip-list applies only to untracked/ignored content.** A **tracked (committed)** file inside one of these dirs is *not* skipped — it stays subject to the tracked-artifact checks (categories 2, 3, 4). A committed `__pycache__/foo.pyc` or `dist/bundle.js` is exactly the pollution category 3 hunts for, so never let the skip-list filter a tracked path.
4. **Always skip the entire `do-work/` tree.** If you find a misplaced `do-work/` directory elsewhere in the repo, do not handle it here — note it once as "see `do-work cleanup` (Pass 3a)" and move on.

### Step 2: Inventory

- **In a git repo:** the source of truth is `git ls-files` for **tracked** files and `git ls-files --others --exclude-standard` for **untracked** files. Tag each file accordingly. **Do not use plain `git status --porcelain` for the untracked inventory** — it collapses a wholly-untracked directory into a single `?? dir/` row, so junk like `tmp/debug.log` inside a brand-new directory never reaches the filename/extension/size/content checks. `git ls-files --others --exclude-standard` lists every untracked file individually *and* already drops paths matched by `.gitignore` (those are correctly ignored — not pollution), so it doubles as the untracked ignore filter; no separate `git check-ignore` pass is needed for untracked paths. (If you prefer `git status`, you must pass `--untracked-files=all` / `-uall` to get the same per-file expansion.) For the **tracked but should-be-gitignored** check (category 2) you still feed `git ls-files` into `git check-ignore --no-index --stdin`: by default `git check-ignore` consults the index and never reports an already-tracked file, so without `--no-index` that category would silently find nothing.
- **Outside git:** walk the filesystem under the scan root, honoring the skip-list.
- **Empty directories** (category 7) — neither `git ls-files` nor `git ls-files --others` emits directories, so empty dirs are invisible to the git-based inventory. Run a separate filesystem pass to find them: `find <scan-root> -type d -empty` and prune the skip-list paths (`-not -path '*/node_modules/*' -not -path '*/.git/*'` etc.). Outside-git mode already walks the filesystem, so the same `find` works there too.

Skip binary files by extension for content-based checks (`.png .jpg .jpeg .gif .webp .ico .pdf .zip .tar .gz .tgz .7z .exe .dll .so .a .o .pyc .class .jar .whl .mp4 .mov .woff .woff2`). For large text files (>500 lines), sample the first 100 and last 50 lines — never full-read a large blob.

### Step 3: Run the Checks

Run each category below. For every finding, record: **path** (tracked/untracked), **category**, **severity**, **one-line reason**, **suggested fix**, and whether it's **auto-fixable**.

| # | Category | Detection | Severity | Auto-fixable? |
|---|----------|-----------|----------|---------------|
| 1 | **Stray temp/backup/OS files** | `*.tmp *.bak *.orig *.rej *.swp *.swo *~ .DS_Store Thumbs.db desktop.ini` | Warning (tracked) / Info (untracked) | Yes — delete (untracked) or `git rm` + gitignore (tracked) |
| 2 | **Tracked but should-be-gitignored** | feed `git ls-files` into `git check-ignore --no-index --stdin` (committed yet covered by an ignore rule). `--no-index` is required — plain `git check-ignore` never reports tracked paths | Warning | Yes — `git rm --cached` + ensure rule in `.gitignore` |
| 3 | **Committed build/generated artifacts** | tracked files under `dist/ build/ out/ target/ .next/ __pycache__/ coverage/ .nuxt/`, or `*.min.js *.min.css *.map` that have a source sibling. Run this against the **tracked** inventory (`git ls-files`) — the Step 1 skip-list must not filter these out, since several of these dirs are *on* the skip-list yet a committed file inside them is exactly the pollution to flag | Warning | Yes — `git rm --cached` + gitignore the dir/pattern |
| 4 | **Committed secrets / sensitive files** | tracked `.env .env.* *.pem *.key *.p12 *.pfx id_rsa id_dsa credentials* *secret*` | **Critical** | Partial — offer `git rm --cached` + gitignore, but **flag loudly: the secret is already in git history; rotate it and scrub history.** Never silently delete |
| 5 | **Misplaced files (folder cohesion)** | nested project markers (`package.json`/`go.mod`/`pyproject.toml`/`Cargo.toml` in a non-root subdir of a single-project repo), a wrong-language file in an otherwise single-language tree, a test file outside the project's test dirs | Info | No — suggest a move (manual; moving breaks imports) |
| 6 | **Duplicate / old-copy files** | name patterns: `* copy.*`, `*copy.*`, `*-old.* *.old *-backup.* *-bak.* *.v2.* *-final.* *-deprecated.* *(1).* * 2.*` | Warning | Untracked → delete on confirm; tracked → review manually first |
| 7 | **Empty files / empty dirs** | size-0 files (from the file inventory) or childless directories (from the `find -type d -empty` pass in Step 2), **excluding** the intentional-empty allowlist (`__init__.py .gitkeep .keep py.typed .npmignore`, plus any directory containing only allowlisted files) | Info | Yes — delete on confirm (respect allowlist) |
| 8 | **Large binary blobs in git** | tracked binary (by extension) larger than **5 MB** — videos, archives, datasets, model weights | Warning | No — suggest Git LFS or removal (may be an intentional asset) |
| 9 | **AI scratch artifacts** | `scratch.* notes.txt tmp.* temp.* debug.log untitled* Untitled* output.txt =* nul .aider*`, plus stray `*.log` or draft `*.md` at the repo root that don't match a known doc name (`README CHANGELOG CONTRIBUTING LICENSE AGENTS CLAUDE`) | Info / Warning | Yes — delete (untracked) / `git rm` (tracked) on confirm |
| 10 | **Dead/unreferenced source files** *(best-effort)* | for each source file, grep the repo for its basename / module path; zero references → candidate. **Exclude** entrypoints (`main index app __init__ conftest setup` and config-declared bins/scripts), test files, and anything dynamically loadable | **Info only** — never Critical | No — verify before removing |

**Category 10 caveat (always state it):** dynamic imports, reflection, framework auto-discovery (routes, plugins, migrations), and CLI entrypoints produce false positives. Report dead-code candidates as Info with "verify before removing" — never auto-fix them.

### Step 4: Report

Emit the severity-grouped report (see Output Format). Group findings under `## Critical Findings`, `## Warnings`, `## Info`, `## Summary`. **Omit any section with no findings.** Every finding names a concrete path, a one-line reason, and a suggested fix; auto-fixable items carry an `[auto-fixable]` tag. End with the summary line. If nothing was found, print the all-clear.

### Step 5: Apply Fixes (Guarded)

**Skip this step entirely** unless the mode is `fix` / `--fix` (the default and `report` / `--report-only` are both report-only), or if there are no auto-fixable findings.

Otherwise:

1. Collect the auto-fixable findings and group them by fix type: **delete (untracked)**, **`git rm --cached` (tracked)**, **append to `.gitignore`**.
2. Present the grouped list and ask, via your environment's ask-user prompt: `Apply these fixes? [all / list the numbers / none]`.
3. Apply **only** the items the user confirms. Support selective application ("all", a number list like `1,3,5`, or "none").
4. Constraints:
   - **Never delete a tracked file with raw `rm`** — use `git rm` (recoverable from history). Prefer `git rm --cached` when the goal is "stop tracking", not "destroy".
   - **Never `git add -A` / `git add .`** — stage only the paths you touched.
   - **Never auto-commit.** After applying, print exactly what changed and suggest `do-work commit`.
   - **Never touch** misplaced-file moves (cat 5), large blobs (cat 8), or dead-code candidates (cat 10) — those are report-only by design.
   - **Secrets (cat 4):** if confirmed, do `git rm --cached` + gitignore, and **repeat the warning** that the secret remains in history and must be rotated.

## Output Format

```markdown
# Stray-Check Report

**Scan date:** <ISO 8601 timestamp>
**Scope:** <path or "repo root">
**Scanned:** <N> tracked + <N> untracked files (<M> skipped: vendored / ignored / do-work)

## Critical Findings

- **[Committed Secret]** `config/.env` (tracked) — environment file with credentials is committed. **Suggested fix:** `git rm --cached config/.env`, add `config/.env` to `.gitignore`, then **ROTATE the secret** — it's already in git history.

## Warnings

- **[Build Artifact]** `dist/bundle.js` (tracked) — generated build output is committed. [auto-fixable] **Suggested fix:** `git rm --cached dist/bundle.js`, gitignore `dist/`.
- **[Should Be Ignored]** `app.log` (tracked) — matches a `.gitignore` rule but is committed. [auto-fixable] **Suggested fix:** `git rm --cached app.log`.

## Info

- **[Empty File]** `src/utils/helpers.ts` (tracked) — zero bytes. [auto-fixable] **Suggested fix:** delete, or add intended content.
- **[Dead Code?]** `src/legacy/oldHelper.ts` (tracked) — no imports found repo-wide. Verify (dynamic import / entrypoint?) before removing.

## Summary

<N> critical, <N> warnings, <N> info items found.
<1-2 sentence recommendation>
```

If `fix` mode and auto-fixable items exist, follow the Summary with the consent prompt:

```
<M> items are auto-fixable. Apply? [all / numbers / none]
```

If nothing was found:

```
# Stray-Check Report

**Scan date:** <ISO 8601 timestamp>
**Scope:** <path or "repo root">

All clear — no stray, misplaced, or orphan files detected.
```

## Rules

- **The scan phase makes zero writes.** Findings are reported; fixes happen only after explicit confirmation (Step 5), and never in `report` mode.
- **Skip the entire `do-work/` tree** and defer misplaced `do-work/` directories to actions/cleanup.md. This action owns repo-wide hygiene, not do-work's bookkeeping.
- **Tracked files are removed with `git rm`, never raw `rm`** — prefer `git rm --cached` when the intent is to stop tracking rather than destroy.
- **Never `git add -A`**, never auto-commit, never touch paths outside the scan root.
- **Secrets are flagged loudly** with a history-retention + rotation warning; `git rm --cached` does not remove them from history — say so.
- **Dead-code findings are Info-only and never auto-fixed.** Always attach the false-positive caveat.
- **Respect `.gitignore`**: an untracked file already covered by an ignore rule is correct, not pollution — don't report it.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This `.bak` is obviously safe to delete — just nuke it" | Report it, and only delete on confirmation; for tracked files use `git rm` | Untracked deletes are unrecoverable; the user may have kept the backup deliberately |
| "It's unused, so it's dead code — delete it" | Report as Info-only with the dynamic-import caveat; never auto-remove | Reflection, route auto-discovery, and CLI entrypoints have zero static references yet are live |
| "I'll just `git rm --cached` the secret and we're done" | Remove it, but warn the secret is still in history and must be rotated | History retains the blob; untracking ≠ scrubbing |
| "A big binary in git is fine, skip it" | Flag it as a Warning and suggest Git LFS | Large blobs bloat every clone forever; the user should decide consciously |
| "There's a `do-work/` dir in a subfolder — I'll relocate it" | Note it and defer to `do-work cleanup` (Pass 3a) | cleanup owns do-work's own files; double-handling risks conflicts |
| "Let me just stage everything and commit the cleanup" | Stage only touched paths; never auto-commit | `git add -A` sweeps unrelated changes; the user commits when ready |

## Red Flags

- A file was deleted or `.gitignore` was modified without an explicit user confirmation — the scan phase must be read-only.
- The report flagged files inside `node_modules/`, `dist/` (untracked + gitignored), or `do-work/` — skip-list or do-work exclusion failed.
- A committed artifact inside a skip-listed dir (e.g. `__pycache__/x.pyc`, `dist/bundle.js` that is *tracked*) was **not** flagged by category 3 — the skip-list wrongly filtered a tracked path. The skip-list is for untracked/ignored noise only.
- Untracked junk inside a brand-new directory (e.g. `tmp/debug.log`) was missed — the untracked inventory used plain `git status --porcelain` (which collapses the dir to `?? tmp/`) instead of `git ls-files --others --exclude-standard` / `-uall`.
- A dead-code candidate was reported as Critical or auto-removed — category 10 is Info-only, never auto-fixed.
- A tracked file was removed with raw `rm` instead of `git rm` — unrecoverable.
- `git add -A` or an auto-commit appeared — staging/commit must be scoped and user-driven.
- A committed secret was reported without a rotation warning — the most important part of the finding is missing.
- The report lists "some files" or generic descriptions instead of concrete paths — findings must name paths.

## Verification Checklist

- [ ] Scan phase made **zero** writes; any fixes were applied only after explicit confirmation (or skipped in `report` mode).
- [ ] The entire `do-work/` tree was skipped; misplaced `do-work/` dirs were deferred to cleanup, not handled here.
- [ ] Findings grouped under `## Critical Findings`, `## Warnings`, `## Info`, `## Summary`; empty sections omitted.
- [ ] Every finding names a concrete path and includes a suggested fix; auto-fixable items tagged `[auto-fixable]`.
- [ ] Untracked files already matched by `.gitignore` were **not** reported as pollution.
- [ ] The untracked inventory listed files **individually** (via `git ls-files --others --exclude-standard` or `git status --porcelain -uall`), not collapsed untracked directories — junk inside brand-new dirs was seen.
- [ ] Tracked files inside skip-listed dirs (e.g. a committed `__pycache__/*.pyc`) still reached the tracked-artifact checks (categories 2/3/4); the skip-list did not filter tracked paths.
- [ ] Dead-code candidates are Info-only with the false-positive caveat; never auto-fixed.
- [ ] Any committed secret carries a history-retention + rotation warning.
- [ ] No `git add -A`, no auto-commit; tracked removals used `git rm`.
