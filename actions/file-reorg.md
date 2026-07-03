# File-Reorg Action

> **Part of the do-work skill.** Reorganize a repository's file layout so an outside coder sees a clean, intentional structure — declutter the root, fold stray and single-file folders into canonical homes, keep folders conceptually small — and move files **without breaking anything**: every reference (code imports, configs, scripts, doc links) is mapped before moving and rewritten after. Plan-first with an explicit consent gate before any file moves.

Map references → design target → present plan → *(user approves)* → `git mv` → rewrite references → verify. **Never move first and grep later.** The planning phases (Steps 1–5) are strictly read-only; moves and rewrites happen only after the user approves the plan.

## When to Use

**Use when:**
- User says "file-reorg", "reorganize the repo", "restructure the layout", "declutter the root", "tidy the layout", "the repo is a hot mess"
- The root is cluttered with docs, reports, scripts, or one-off files that belong in folders
- Stray or single-file folders should be folded into canonical homes (`docs/`, a reports folder, a scripts folder)

**Do NOT use when:**
- The problem is junk — temp/backup files, committed build artifacts, should-be-gitignored files → use `actions/stray-check.md` (deletes pollution; this action relocates legitimate files)
- The clutter is do-work's own bookkeeping — loose REQs, misplaced `do-work/` directories → use `actions/cleanup.md`
- The user wants code-structure refactoring (splitting modules, moving source between packages, renaming symbols) → capture as a REQ via `do-work capture-request:` — import-rewriting churn belongs in the work pipeline with tests, not a layout pass

## Input

`$ARGUMENTS` — all optional, combinable (e.g., `do-work file-reorg docs/ plan`):

1. **Path scope** (e.g., `docs/`, `packages/api/`) — limit the reorg to a subtree. Default: whole repo.
2. **Mode token:**
   - *(default, no token)* — full guarded flow: plan (Steps 1–5), stop for approval, then execute + verify (Steps 6–8).
   - `plan` / `--plan-only` — produce the reorg plan (Steps 1–5) and stop. Zero writes.

## Steps

### Step 1: Baseline

1. Run the test suite BEFORE touching anything and record the result (pass/fail counts). Post-reorg failures must be attributable to the reorg, not pre-existing breakage. No test suite → note that and rely on the other Step 8 checks.
2. Check `git status` — if there are staged or uncommitted changes, **stop and ask**: the user should commit them first (suggest `do-work commit`) or explicitly accept mixing. The reorg must land as its own atomic commit.
3. Read the repo's CLAUDE.md / AGENTS.md / contributing docs for layout rules that already exist — the plan must respect them, not fight them.

### Step 2: Inventory & Classification

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
3. Apply the folder-size heuristic: aim for ≤ ~10 conceptual entries per folder; treat ~25 as the hard ceiling. Fold single-file folders into a sibling. Do NOT churn big code folders (`lib/`, `src/`, `tests/`) just to hit a number — moving code rewrites imports repo-wide for little presentational gain; flag it in the plan as an optional follow-up instead.

### Step 3: Reference Mapping (before any move)

For EVERY file/folder on the move list, find who references it. Be exhaustive — this is the phase that makes the reorg safe. If your environment supports subagents, fan the search categories out in parallel following the durability pattern in `crew-members/background-agents.md`; otherwise run them sequentially.

1. **Invocations & configs**: package.json (`main`, `scripts`), Makefile/justfile, shell scripts, CI configs, Dockerfiles, editor/test configs.
2. **Code imports**: `import`/`require` of moved files, including tests that import or spawn entry points by relative path.
3. **Internal path resolution of moved code** — the classic breakers:
   - relative imports (`./lib/...` becomes `../lib/...` one level deeper);
   - `__dirname` / `import.meta.url` used to locate SIBLING resources (static roots, asset dirs) — after the move these must resolve to the repo root, e.g. `path.resolve(script_dir, "..")`;
   - cwd-relative fs paths (usually fine if the run recipe stays "run from repo root" — verify);
   - spawns of sibling scripts.
4. **Doc links**: relative markdown links in moved docs (outbound) and links pointing AT moved docs (inbound). Depth changes differ per destination — a file moving to `docs/` needs one `../` added to root-relative links; a file moving to `docs/specs/` needs two. Links between files that move TOGETHER as siblings need NO change — a blanket "add `../` everywhere" pass breaks them.
5. **Agent/skill configs** (`.claude/`, hooks, settings): check for hardcoded paths. Version-locked or vendored skills must NOT be edited — instead verify they discover files by glob (e.g. `**/prime-*.md`) and stay compatible with the target layout; if they can't, drop that move from the plan.
6. Classify every hit **live** (must update) vs **historical** (archives, old reports — leave untouched).

### Step 4: Adversarial Gap Pass

Run a "what did I miss" pass over the map before presenting it — fresh eyes or a second agent. Check the places the first sweep habitually skips: READMEs in subfolders, error-message strings, comments in code that cite paths, CI badge/status URLs, doc links written as bare code spans instead of markdown links.

### Step 5: Present the Plan (consent gate)

Present the reorg plan in the Output Format below: move table, reference-rewrite map, risk notes, and any flagged-but-excluded items (big code folders, incompatible vendored skills).

- **Mode `plan` / `--plan-only`:** stop here. Zero writes.
- **Default mode:** ask via your environment's ask-user prompt: `Execute this plan? [yes / trim it (list numbers to drop) / no]`. Execute exactly what was approved — nothing more. New move candidates noticed later go in a follow-up plan, not into this execution.

### Step 6: Execute

1. `git mv` only for tracked files (history-preserving renames), batched by destination. Untracked files on the move list use plain `mv` (nothing to preserve).
2. Rewrite references from the Step 3 map with EXPLICIT per-file replacement pairs — scripted string replacement that prints a warning for any pair that matched zero times. No blind repo-wide regex.
3. While touching a file, fix pre-existing broken links you find (note them for the commit message), but do not restyle unrelated content.

### Step 7: Make the Structure Self-Enforcing

1. Update README and CLAUDE.md to describe the new layout.
2. Add a **root-hygiene rule** to CLAUDE.md (root allowlist + where new code/docs/reports go) so the structure holds for future sessions.
3. If the README is still hosting-provider boilerplate, rewrite it: what the project is, quick start, layout table, doc entry points.
4. README and CLAUDE.md prose is a human-facing artifact — load `crew-members/anti-slop.md` before writing it.

### Step 8: Verify (all of it, not a sample)

1. Test suite passes at the baseline count or better.
2. Boot the app / run the primary entry point; smoke-test the routes or commands whose path resolution you touched (static roots especially).
3. Syntax-check moved scripts (`node --check`, `python -m py_compile`, …).
4. Link-check EVERY relative markdown link in docs, primes, and READMEs. If the repo has no such test, ADD one as a permanent regression net (walk the tree, resolve `[label](relative/path)` targets, skip http/anchors and links escaping the repo).
5. Re-list the root and every changed folder; confirm the result reads clean against the Step 2 buckets.
6. Grep once more for old paths (e.g. `rg -n "\./old-name"`, bare `node old-entry.js`) to catch stragglers — excluding the historical records identified in Step 3.6.

### Step 9: Report

Emit the post-execution report (Output Format below). Do **not** commit. Suggest `do-work commit` and state the atomicity requirement: renames + reference rewrites + doc updates belong in ONE commit, separate from unrelated work, so the move is revertable as a unit; note any pre-existing broken links fixed along the way so the commit message can mention them. Never push.

## Output Format

**The plan (Step 5):**

```markdown
# Reorg Plan

**Scope:** <path or "repo root">
**Baseline:** <N passed / M failed, or "no test suite">

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
```

**The post-execution report (Step 9):** moves executed, rewrite pairs applied (and any zero-match warnings), verification results item-by-item against Step 8, pre-existing broken links fixed, and the commit suggestion.

## Rules

- **Steps 1–5 make zero writes.** Execution happens only after the user approves the plan, and only the approved plan.
- **`git mv` only for tracked files** — a delete+add pair in the diff means history was dropped.
- **Historical records are never edited** — stale paths inside archives are correct history; rewriting them destroys the point-in-time record.
- **No blind repo-wide regex.** Every rewrite is an explicit per-file pair; zero-match pairs print a warning instead of silently passing.
- **Respect existing conventions** — fold into the folders the repo already uses; never invent a parallel convention.
- **Don't churn big code folders** to satisfy the size heuristic — flag as an optional follow-up.
- **Never edit version-locked/vendored skills** — verify glob compatibility or drop the move.
- **Never auto-commit, never push.** Suggest `do-work commit` with the one-atomic-commit note.

## Common Rationalizations

| If you're thinking...                                                     | STOP. Instead...                                                                     | Because...                                                                                             |
| ------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| "The reference map looks complete — skip the adversarial pass"            | Run Step 4; check subfolder READMEs, error strings, path-citing comments             | The first sweep reliably misses non-import references; those are exactly the ones that break silently  |
| "One repo-wide sed will catch every old path"                             | Explicit per-file replacement pairs with zero-match warnings                          | A global regex rewrites historical records and coincidental matches; a zero-match pair is a mapping bug it would hide |
| "Every moved doc needs `../` added to its links"                          | Compute depth per destination; sibling files moving together need no change          | Depth deltas differ per subfolder — a blanket pass breaks exactly the links that were fine             |
| "These archived files reference old paths — I'll fix them while I'm here" | Leave historical records untouched; classify them in Step 3.6 and exclude them       | They're point-in-time records; stale paths inside them are correct history                             |
| "`src/` has 40 files, way over the ceiling — split it up"                 | Flag it as an optional follow-up in the plan                                          | Moving code rewrites imports repo-wide for presentational gain; that churn needs the work pipeline     |
| "Tests were already red, so I'll skip the baseline"                       | Record the baseline anyway and compare counts after                                   | Without the baseline you can't attribute post-reorg failures — every red test becomes your suspect    |
| "The vendored skill hardcodes a path — quick patch and done"              | Don't touch it; verify glob discovery or drop that move from the plan                 | Version-locked content gets overwritten on update; your patch dies and the break comes back            |
| "While executing I found two more files that should move — I'll add them" | Execute exactly the approved plan; new candidates go in a follow-up plan              | The user approved a specific move list; silent scope growth breaks the consent gate                    |

## Red Flags

- Files were moved before a reference map existed for them — the "never move first and grep later" invariant broke.
- `git status` after execution shows delete+add pairs for tracked files instead of renames — plain `mv` was used where `git mv` was required.
- Contents of an archive/historical folder were edited.
- A replacement pair matched zero times and no warning surfaced — the map and the tree disagree and it was swallowed.
- README or CLAUDE.md still describe the old layout after execution.
- Verification checked a sample of links "to be safe on time" instead of every relative link.
- The reorg was auto-committed, pushed, or bundled with unrelated changes.
- Plan-only mode produced any file change at all.

## Verification Checklist

- [ ] Test suite passes at the baseline count or better (baseline recorded in Step 1, compared in Step 8).
- [ ] Primary entry point boots; every touched path-resolution site was smoke-tested.
- [ ] Moved scripts pass a syntax check.
- [ ] Every relative markdown link in docs/primes/READMEs resolves; a link-check regression net exists (added if the repo had none).
- [ ] Straggler grep for old paths is clean, excluding classified historical records.
- [ ] Root and every changed folder re-listed; result matches the approved plan's buckets.
- [ ] Steps 1–5 made zero writes; execution ran only after explicit approval and covered exactly the approved list.
- [ ] Historical record contents are byte-identical to before the reorg.
- [ ] No commit, no push; `do-work commit` suggested with the one-atomic-commit note.
