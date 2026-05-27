# Stray-Check

Repo-wide scan for orphan/junk files that pollute where they don't belong — leftover temp/backup files, committed build artifacts, files that should be gitignored, misplaced/duplicate/empty files, large blobs, AI scratch droppings, and best-effort dead code. Report-only by default; applies fixes only on explicit confirmation.

> **Not to be confused with cleanup or forensics.** `do-work cleanup` and `do-work forensics` both work on do-work's *own* bookkeeping (loose REQs, misplaced `do-work/` directories, stuck/hollow work). Stray-check works on the **whole repository** and deliberately skips the `do-work/` tree. If the clutter is do-work's own files, use cleanup; if you want pipeline health, use forensics.

## What it checks

| Check | What it detects | Severity |
|-------|-----------------|----------|
| **Stray temp/backup/OS files** | `*.tmp *.bak *.orig *.rej *.swp *~ .DS_Store Thumbs.db desktop.ini` | Warning (tracked) / Info (untracked) |
| **Tracked but should-be-gitignored** | Committed files that a `.gitignore` rule already matches | Warning |
| **Committed build/generated artifacts** | Tracked files under `dist/ build/ out/ target/ .next/ coverage/`, or `*.min.js`/`*.map` with a source sibling | Warning |
| **Committed secrets / sensitive files** | Tracked `.env *.pem *.key id_rsa credentials* *secret*` | **Critical** |
| **Misplaced files** | Nested project markers, wrong-language files, tests outside test dirs | Info |
| **Duplicate / old-copy files** | `* copy.*`, `*-old.*`, `*.v2.*`, `*-backup.*`, `* 2.*`, etc. | Warning |
| **Empty files / empty dirs** | Zero-byte files / childless dirs (excludes `.gitkeep`, `__init__.py`, `py.typed`) | Info |
| **Large binary blobs in git** | Tracked binaries > 5 MB (videos, archives, datasets) | Warning |
| **AI scratch artifacts** | `scratch.* notes.txt debug.log untitled* output.txt`, stray root `*.log`/draft `*.md` | Info / Warning |
| **Dead/unreferenced source files** | Source files with zero repo-wide references (best-effort) | Info only |

## Output

Markdown report organized by severity (`Critical Findings` / `Warnings` / `Info` / `Summary`), sections with no findings omitted. Each finding names a concrete path, a one-line reason, and a suggested fix; items that can be fixed automatically carry an `[auto-fixable]` tag. If nothing is found, the report says "All clear."

## Fix mode

By default stray-check is a dry run — it reports and stops. Add `fix` to act on findings:

1. It groups the auto-fixable findings by fix type — **delete (untracked)**, **`git rm --cached` (tracked)**, **append to `.gitignore`**.
2. It asks for explicit confirmation: `Apply these fixes? [all / numbers / none]`.
3. It applies **only** what you confirm.

Safety rails:
- Tracked files are removed with `git rm` (recoverable from history), never raw `rm`; `git rm --cached` is preferred when the goal is just to stop tracking.
- Never `git add -A`, never auto-commit — only touched paths are staged, and you commit when ready.
- Misplaced-file moves, large blobs, and dead-code candidates are **never** auto-fixed (reported only).
- Committed secrets get a loud warning: `git rm --cached` does **not** remove them from history — rotate the secret and scrub history.

## Key rules

- Read-only scan phase — zero writes until you confirm fixes.
- Skips the entire `do-work/` tree; defers misplaced `do-work/` directories to `do-work cleanup`.
- Respects `.gitignore` — a correctly-ignored untracked file is not pollution.

## Usage

```
do-work stray-check                 Report on the whole repo (dry run)
do-work stray-check src/            Limit the scan to a subtree
do-work stray-check fix             Scan, then apply fixes on confirmation
do-work stray-check report          Force pure read-only (suppress the fix prompt)
do-work find orphan files           Same thing
```

## When NOT to use

- Loose REQs or a misplaced `do-work/` directory → `do-work cleanup`
- do-work pipeline health (stuck/hollow REQs, orphaned URs) → `do-work forensics`
- Validating a single human-facing artifact (brief, report) → `do-work slop-check`
- Code-quality review of a change → `do-work code-review`
