# Setup Memory Action

> Explicitly scaffolds the project memory store and, on Claude Code, composes the optional knowledge hooks into existing settings. Suite installation never invokes this action automatically.

## Preconditions

Run only when the user explicitly asks for `do-work-knowledge setup-memory` or its documented `install memory-module` alias. Load `crew-members/clear-questions.md` before any confirmation.

Resolve the project and knowledge skill roots:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
KNOWLEDGE_ROOT="<directory containing this package's SKILL.md>"
```

The two deterministic legacy-to-modular hook migrations are:

```text
.claude/skills/do-work/hooks/memory-session-start.sh -> .claude/skills/do-work-knowledge/hooks/memory-session-start.sh
.claude/skills/do-work/hooks/memory-stop-capture.sh  -> .claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh
```

Replace only those exact command substrings. Never rewrite unrelated hook commands or whole hook arrays.

## Phase 1: Inspect before writing

Check all four store components: non-empty `memory/working-memory.md`, `memory/logs/`, `memory/usage-ledger.jsonl`, and the machine-local ignore coverage for logs, ledger, and `.bootstrap-imported`.

If Git tracks any of these raw-store paths, stop and name them:

```bash
git -C "$PROJECT_ROOT" ls-files -- memory/logs memory/usage-ledger.jsonl memory/.bootstrap-imported
```

Tell the user to review and run `git rm --cached <paths>` themselves. Never untrack or delete memory on their behalf; an ignore rule cannot protect an already indexed transcript.

Inspect `$PROJECT_ROOT/.claude/settings.json` when present. Invalid JSON is a hard stop before scaffolding or hook edits. Record its bytes and existing hook entry counts for later verification.

## Phase 2: Scaffold without overwriting

Create missing directories and the empty ledger. Create `memory/working-memory.md` from `actions/memory-reference.md` → **working-memory.md Template**, with today's UTC date, only when the file is absent or empty. Never overwrite standing memory.

Add these patterns independently to the repository's local Git exclude file only when not already ignored:

```text
**/memory/logs/
**/memory/usage-ledger.jsonl
**/memory/.bootstrap-imported
```

Use `git rev-parse --git-path info/exclude`; do not edit the project's `.gitignore`. Outside Git, report that machine-local ignore protection is unavailable and continue with the portable hookless actions.

## Phase 3: Offer optional hooks

If the project has no `.claude/` directory, report `hooks: n/a` and skip. Otherwise show that the optional hooks inject the frozen curated snapshot at session start and append a redacted final exchange at stop. Ask whether to enable them. A missing or ambiguous answer means no hook change.

When enabled, use `hooks/memory-hooks.json` from this package. Its commands must target:

```text
.claude/skills/do-work-knowledge/hooks/memory-session-start.sh
.claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh
```

If `jq` is unavailable, leave settings unchanged, print the fragment, and report `hooks: MANUAL STEP`. Do not text-patch JSON.

With `jq` available:

1. Create `{}` only when settings is absent; otherwise require it to parse.
2. Back up the exact original to `settings.json.pre-memory-module` before the first mutation.
3. Migrate each known legacy command substring above to its exact modular target. No other string changes.
4. Gate each hook independently by its script filename, then append only the missing fragment entry with array `+`; never assign a replacement hooks object.
5. Write through a sibling temporary file and rename only after `jq` validates the candidate.
6. Verify both filenames, every pre-existing unrelated hook, and entry counts. On any failure, restore the exact backup and report failure.
7. Remove the backup only after verification succeeds.

The complete append expression and fragment shape live in `actions/memory-reference.md` → **Hook Install Internals**. Follow them rather than inventing a second merge algorithm.

## Phase 4: Verify and report

Verify:

- `memory/working-memory.md` exists, is non-empty, and was not overwritten;
- `memory/logs/` and `memory/usage-ledger.jsonl` exist;
- Git tracks none of the three raw-store paths;
- local ignore covers logs, ledger, and sentinel independently;
- settings parses and retains every prior unrelated hook;
- when enabled, both commands point to `do-work-knowledge`, with no known legacy memory-hook command left;
- when declined or unsupported, settings is byte-identical to its pre-action state.

Report scaffolding and hooks separately. Never call hookless memory unusable: every memory subcommand remains functional without automatic capture.

## Rules

- Fresh suite installation does not call this action and does not enable memory capture.
- Never read transcript history during setup; `memory bootstrap` owns that separate consent gate.
- Never commit raw logs, the ledger, or bootstrap sentinel.
- Never delete user memory, rewrite `.gitignore`, clobber settings, or weaken exact rollback.
- Hook writers always redact before truncating and always exit zero; setup does not alter those scripts.
