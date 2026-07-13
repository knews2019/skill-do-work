# Tidy-Repo Action

> **Part of the do-work skill.** Tidy a repository's file layout so an outside coder sees a clean, intentional structure — declutter the root, fold legitimate strays into canonical homes, keep folders conceptually small — and move files **without breaking anything**: every reference (code imports, configs, scripts, doc links) is mapped before moving and rewritten after. Plan-first with an explicit consent gate before any file moves.

Map references → design target → present plan → *(user approves)* → `git mv` → rewrite references → verify. **Never move first and grep later.** The planning phases (Steps 1–5) are strictly read-only; moves and rewrites happen only after the user approves the plan.

## When to Use

**Use when:**
- User says "tidy-repo", "tidy the repo", "reorganize the repo", "restructure the layout", "declutter the root", or "the repo is a hot mess"
- The root is cluttered with docs, reports, scripts, or one-off files that belong in folders
- Stray or single-file folders should be folded into canonical homes (`docs/`, a reports folder, a scripts folder)

**Do NOT use when:**
- The problem is junk — temp/backup files, committed build artifacts, should-be-gitignored files → use `actions/stray-check.md` (deletes pollution; this action relocates legitimate files)
- The clutter is do-work's own bookkeeping — loose REQs, misplaced `do-work/` directories → use `actions/cleanup.md`
- The user wants code-structure refactoring (splitting modules, moving source between packages, renaming symbols) → capture as a REQ via `do-work capture-request:` — import-rewriting churn belongs in the work pipeline with tests, not a layout pass

## Input

`$ARGUMENTS` — all optional, combinable (e.g., `do-work tidy-repo docs/ plan`). The legacy `do-work file-reorg` alias routes here with the same arguments:

1. **Path scope** (e.g., `docs/`, `packages/api/`) — limit the tidy to a subtree. Default: whole repo.
2. **Mode token:**
   - *(default, no token)* — full guarded flow: plan (Steps 1–5), stop for approval, then execute + verify (Steps 6–8).
   - `plan` / `--plan-only` — produce the tidy-repo plan (Steps 1–5) and stop with no intentional project-file writes.

## Steps

### Step 1: Baseline

1. Run the test suite BEFORE touching anything and record the result (pass/fail counts). Post-tidy failures must be attributable to the approved moves, not pre-existing breakage. No test suite → note that and rely on the other Step 8 checks.
2. Run `git status --short --untracked-files=all` and record dirty paths. Planning may continue, but before execution stop for a decision if an approved move or reference rewrite overlaps a dirty path. Preserve unrelated user changes, keep them out of the plan, and report them explicitly.
3. Read the repo's CLAUDE.md / AGENTS.md / contributing docs for layout rules that already exist — the plan must respect them, not fight them.

### Step 2: Inventory, Classification & Target Design

1. Inventory **both** tracked and untracked files — `git ls-files` misses untracked clutter entirely:
   ```
   git ls-files | awk -F/ '{print $1}' | sort | uniq -c | sort -rn
   git ls-files --others --exclude-standard
   ```
   (The second command lists every untracked file individually and already drops correctly-ignored paths — do not substitute plain `git status --porcelain`, which collapses wholly-untracked directories into one row.)
2. Classify every root file and stray folder into exactly one bucket:
   - **Tool-mandated config** (stays in root): package.json / lockfiles, justfile / Makefile, test configs, `.gitignore`, env example files, README, CLAUDE.md / AGENTS.md, licenses.
   - **Executables / entry points** → the repo's existing executables folder (`cmd/`, `bin/`, `scripts/`). Respect what the repo already uses — don't invent a second convention.
   - **Durable documentation** (architecture docs, primes, lessons, specs, runbooks, reusable prompts) → `docs/`, with subfolders per kind (`docs/lessons-learned/`, `docs/specs/`, `docs/handoffs/`).
   - **Generated one-off reports** → the repo's report-output folder (e.g. `ai-reports/`). If a skill or tool auto-writes to a folder, that folder is a **fold-target** — merge strays into it, never move it. The `do-work/` tree itself is always a fixed point: never relocate or restructure it here.
   - **Historical records** (task-queue archives, old REQ/UR files, dated reports): relocate the folder **as a whole** if needed, but NEVER edit the contents — they are point-in-time records; stale paths inside them are correct history.
3. Treat folder size as a signal, not a quota: aim for ≤ ~10 conceptual entries per folder and review folders around ~25. Fold a single-file folder only when that improves ownership or discovery. Do NOT churn big code folders (`lib/`, `src/`, `tests/`) just to hit a number — moving code rewrites imports repo-wide for little presentational gain; flag it in the plan as an optional follow-up instead.
4. Draft the target layout and a numbered move list. Give every destination a reason grounded in ownership, discoverability, or an existing convention; exclude moves whose benefit is merely cosmetic relative to their reference churn.

### Step 3: Reference Mapping (before any move)

For EVERY file/folder on the move list, find who references it. Be exhaustive — this is the phase that makes the tidy safe. For a large map, parallelize the search categories only when the environment and user instructions allow it, following `crew-members/background-agents.md`; otherwise run them sequentially.

1. **Invocations & configs**: package.json (`main`, `scripts`), Makefile/justfile, shell scripts, CI configs, Dockerfiles, editor/test configs.
2. **Code imports**: `import`/`require` of moved files, including tests that import or spawn entry points by relative path.
3. **Internal path resolution of moved code** — the classic breakers:
   - relative imports (`./lib/...` becomes `../lib/...` one level deeper);
   - `__dirname` / `import.meta.url` used to locate SIBLING resources (static roots, asset dirs) — after the move these must resolve to the repo root, e.g. `path.resolve(script_dir, "..")`;
   - cwd-relative fs paths (usually fine if the run recipe stays "run from repo root" — verify);
   - spawns of sibling scripts.
4. **Doc links**: relative markdown links in moved docs (outbound) and links pointing AT moved docs (inbound). Depth changes differ per destination — a file moving to `docs/` needs one `../` added to root-relative links; a file moving to `docs/specs/` needs two. Links between files that move TOGETHER as siblings need NO change — a blanket "add `../` everywhere" pass breaks them.
5. **Agent/skill configs** (`.claude/`, `.codex/`, hooks, settings): check for hardcoded paths. Version-locked or vendored skills must NOT be edited — instead verify they discover files by glob (e.g. `**/prime-*.md`) and stay compatible with the target layout; if they can't, drop that move from the plan.
6. **Generated files**: update the generator, template, or source configuration rather than hand-editing output; regenerate only when the repository's normal workflow requires it.
7. Classify every hit as **live** (must update), **historical** (leave untouched), or **generated** (update its source). Flag dynamic construction, case-only renames, symlinks, and paths that escape the repo.

### Step 4: Adversarial Gap Pass

Run a "what did I miss" pass over the map before presenting it — fresh eyes or, when the environment and user instructions allow it, a second agent following `crew-members/background-agents.md`. Check the places the first sweep habitually skips: READMEs in subfolders, error-message strings, comments in code that cite paths, CI badge/status URLs, doc links written as bare code spans, generators, case-only renames, and symlinks.

### Step 5: Present the Plan (consent gate)

Present the tidy-repo plan in the Output Format below: move table, reference-rewrite map, dirty-path overlaps, risk notes, verification commands, and any flagged-but-excluded items (big code folders, incompatible vendored skills).

- **Mode `plan` / `--plan-only`:** stop here. Zero writes.
- **Default mode:** ask via your environment's ask-user prompt: `Execute this plan? [yes / trim it (list numbers to drop) / no]`. Execute exactly what was approved — nothing more. New move candidates noticed later go in a follow-up plan, not into this execution.

### Step 6: Execute

1. Use `git mv` for tracked files (history-preserving renames), batched by destination. Untracked files on the move list use plain `mv` (nothing to preserve). For a case-only rename on a case-insensitive filesystem, move through a temporary intermediate name.
2. Rewrite references from the Step 3 map with EXPLICIT per-file replacement pairs — scripted string replacement that prints a warning for any pair that matched zero times. Stop and investigate a zero match; never hide it with a blind repo-wide regex. Update a generator or source template before regenerated output.
3. Apply exactly the approved moves and mapped rewrites. Do not fix unrelated links, restyle content, or add newly noticed moves during execution.

### Step 7: Make the Structure Self-Enforcing

1. Update README, CLAUDE.md, AGENTS.md, or contributing docs only when they describe the old layout or the approved reorganization establishes a durable placement rule.
2. Express placement rules as conditions (where code, docs, reports, and generated output belong), with examples marked illustrative. Avoid a brittle closed allowlist unless an external tool truly mandates one.
3. Load `crew-members/anti-slop.md` before writing human-facing documentation and keep the edit limited to layout facts.
4. Treat hosting-provider boilerplate or broader documentation problems as follow-up work, not a mandatory side quest.

### Step 8: Verify (all of it, not a sample)

1. Test suite passes at the baseline count or better.
2. If an approved move can affect runtime path resolution, boot the app / run the primary entry point and smoke-test the touched routes or commands (static roots especially); otherwise mark this check not applicable.
3. Syntax-check moved scripts (`node --check`, `python -m py_compile`, …).
4. Run the repository's link checker if one exists. Otherwise verify EVERY relative link in moved Markdown and every changed document that points to a moved file; report the missing regression net instead of adding permanent tooling without approval.
5. Re-list the root and every changed folder; confirm the result reads clean against the Step 2 buckets.
6. Grep once more for old paths (e.g. `rg -n "\./old-name"`, bare `node old-entry.js`) to catch stragglers — excluding the historical records identified in Step 3.7.
7. Run `git diff --check` and inspect `git diff --summary --find-renames`. Confirm tracked moves appear as renames, historical contents and pre-existing dirty paths are unchanged, and no file outside the approved plan changed.

### Step 9: Report

Emit the post-execution report (Output Format below). Do **not** commit. Suggest `do-work commit` and state the atomicity requirement: renames + reference rewrites + required layout-doc updates belong in ONE commit, separate from unrelated work, so the move is revertable as a unit. Never push.

## Output Format

**The plan (Step 5):**

```markdown
# Tidy-Repo Plan

**Scope:** <path or "repo root">
**Baseline:** <N passed / M failed, or "no test suite">
**Dirty paths preserved:** <none or list>

## Moves

| # | From | To | Bucket | Live references |
|---|------|----|--------|-----------------|
| 1 | `WORKLOG.md` | `docs/worklog.md` | durable doc | 2 (README.md, justfile) |

## Reference Rewrites

- `README.md`: `WORKLOG.md` → `docs/worklog.md`
- `justfile`: `cat WORKLOG.md` → `cat docs/worklog.md`

## Flagged, Not Moved

- `lib/` (31 files) — over the folder ceiling, but splitting it rewrites imports repo-wide. Optional follow-up REQ.
- `.claude/skills/<vendored>/` — hardcodes `reports/`; version-locked, so `reports/` stays put.

## Risks

- <anything with __dirname/sibling resolution, cwd-relative paths, or dynamic loading>

## Verification

- <commands/checks that will prove the approved moves safe>
```

**The post-execution report (Step 9):** moves executed, rewrite pairs applied (and any zero-match warnings), verification results item-by-item against Step 8, preserved dirty paths, exclusions/follow-ups, and the commit suggestion.

## Rules

- **Steps 1–5 make no intentional project-file writes.** Execution happens only after the user approves the plan, and only the approved plan.
- **Preserve existing user changes.** Never overwrite, relocate, or silently absorb dirty paths; overlapping paths require a decision.
- **`git mv` only for tracked files** — a delete+add pair in the diff means history was dropped.
- **Historical records are never edited** — stale paths inside archives are correct history; rewriting them destroys the point-in-time record.
- **No blind repo-wide regex.** Every rewrite is an explicit per-file pair; zero-match pairs print a warning instead of silently passing.
- **Respect existing conventions** — fold into the folders the repo already uses; never invent a parallel convention.
- **Don't churn big code folders** to satisfy the size heuristic — flag as an optional follow-up.
- **Never edit version-locked/vendored skills** — verify glob compatibility or drop the move.
- **Never auto-commit, never push.** Suggest `do-work commit` with the one-atomic-commit note.

## Red Flags

- Files were moved before a reference map existed for them — the "never move first and grep later" invariant broke.
- `git status` after execution shows delete+add pairs for tracked files instead of renames — plain `mv` was used where `git mv` was required.
- Contents of an archive/historical folder were edited.
- A replacement pair matched zero times and no warning surfaced — the map and the tree disagree and it was swallowed.
- Repository guidance still describes an old layout affected by the approved moves, or unrelated documentation was rewritten as a side quest.
- Verification checked a sample of links "to be safe on time" instead of every affected relative link.
- The tidy was auto-committed, pushed, or bundled with unrelated changes.
- Plan-only mode intentionally changed a project file.

## Verification Checklist

- [ ] Test suite passes at the baseline count or better (baseline recorded in Step 1, compared in Step 8).
- [ ] Every touched runtime path-resolution site was smoke-tested, or the check was explicitly not applicable.
- [ ] Moved scripts pass a syntax check.
- [ ] Every relative link affected by the approved moves resolves; any missing permanent link checker is reported as follow-up work.
- [ ] Straggler grep for old paths is clean, excluding classified historical records.
- [ ] Root and every changed folder re-listed; result matches the approved plan's buckets.
- [ ] Steps 1–5 made no intentional project-file writes; execution ran only after explicit approval and covered exactly the approved list.
- [ ] Historical record contents and pre-existing dirty paths are byte-identical to before the tidy.
- [ ] No commit, no push; `do-work commit` suggested with the one-atomic-commit note.
