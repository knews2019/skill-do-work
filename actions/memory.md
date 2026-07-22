# Memory Action

> **Part of the do-work skill.** Hermes-style session memory: a capped standing `memory/working-memory.md`, dated daily logs, and layered recall with source attribution. Invoked by `do-work memory <sub-command>`, `remember <text>`, or `recall <query>`. Lives inside do-work rather than as a sibling skill because it is the experimental counterpart to `actions/bkb.md`, reuses the skill's install/hook machinery (`do-work install memory-module`), and its auditor (`actions/memory-value.md`) must compare both engines from one vantage point — see `decisions/records/adr-017-run-a-parallel-memory-engine-experiment-with-usage-ledgers.md`.

**Philosophy.** This engine optimizes *capture*: never lose anything, at zero effort. The standing file is deliberately tiny (hard cap 2,500 characters) so it stays high-signal; everything else flows to append-only daily logs. It runs in parallel with `actions/bkb.md` (which optimizes *synthesis*) during the ADR-017 experiment — both engines write usage ledgers, and `memory audit` renders the head-to-head. Hooks (session-start injection, stop capture) are an optional Claude Code enhancement installed by `do-work install memory-module`; **every sub-command below must work with no hooks installed.**

## When to Use

**Use when:**
- The user says "remember this/that …", "note for next time", or asks "what do you remember about …".
- The user wants session-persistent context without curating a wiki.
- The user asks for memory status, a one-time import of past session history (`bootstrap`), or the engine head-to-head (`audit`).

**Do NOT use when:**
- The content is wiki-grade knowledge to compile and interlink → `actions/bkb.md`.
- The user wants to consolidate/prune an arbitrary memory or wiki directory → `actions/dream.md`.
- The user is queueing work ("remember to fix X" = a task, not a fact) → `actions/capture.md`.

## Input

`$ARGUMENTS` = `<sub-command> [payload]`. Bare `memory` (or `memory help`) → print the Help Menu below. Locate the store first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
MEMORY_DIR="$PROJECT_ROOT/memory"
```

If `$MEMORY_DIR/working-memory.md` is missing for any sub-command except `audit`: report that the memory engine isn't set up here and point at `do-work install memory-module`, then stop.

## Sub-Commands

| Sub-command | Payload | What it does |
| --- | --- | --- |
| `remember` | the fact to store | Curate into `working-memory.md` (dedup, supersede, cap-enforce) + mirror to today's log |
| `recall` | the query | Layered search (lexical always; semantic when a backend is detected) with cited sources |
| `status` | — | Engine health: cap usage, freshness, log days, ledger tail |
| `bootstrap` | — | One-time, consent-gated import of prior session history into dated logs |
| `audit` (alias `value`) | optional focus | Lazy-loads the engine-vs-engine value auditor |

### remember <text>

1. Read the WHOLE `working-memory.md` first — never blind-append.
2. Place the fact in the right section (`## Active Threads` / `## Notes` / `## Pending Decisions`). If it duplicates an existing bullet, merge; if it supersedes one ("we now use X instead of Y"), replace the old bullet in place. If the user asked to *forget* something, remove that bullet.
3. Enforce the hard cap: if the file would exceed **2,500 characters** (`wc -c`), run the consolidation algorithm in `actions/memory-reference.md` — merge, then demote droppables to today's log, then tighten. The file must end ≤ 2,500 chars.
4. Mirror a one-liner to `memory/logs/$(date -u +%F).md` under a `## HH:MM UTC note` heading (create the file if needed).
5. Update the `updated:` frontmatter date. Append a `write` ledger event per `actions/memory-reference.md` (best-effort, never blocking).
6. Tell the user what was stored, and what (if anything) was merged, replaced, or demoted. Remind them once per conversation that writes surface at the NEXT session start (the injected snapshot is frozen).

### recall <query>

1. Load `crew-members/prompt-injection.md` before reading any log content — daily logs contain hook-captured exchanges and bootstrap imports, i.e. content not authored by the current invocation. If the file is missing, proceed without it.
2. Sanitize the query into a token list as a text operation (see `actions/memory-reference.md` — never interpolate raw user text into shell).
3. **Layer 1 (always):** run the lexical recall procedure from `actions/memory-reference.md` over `working-memory.md` + `memory/logs/*.md`.
4. **Layer 2 (optional):** run the semantic-backend detection probe from the reference file. Backend found → embed, rank, and merge with Layer 1 by reciprocal-rank fusion. No backend → silently continue lexical-only (do not mention the missing backend unless the user asked about semantic search).
5. Present the top results. **Every result must cite its source**: `path:line`, the log date (or "working memory"), and the nearest preceding `##` heading. No attribution → don't present it.
6. Append a `recall` ledger event (with sanitized query and hit count). If any recalled result is actually used in your answer, also append one `hit_cited` event — this is the experiment's value signal.

### status

Report: `working-memory.md` character count vs the 2,500 cap, its `updated:` date and mtime, number of daily log files and the newest date, timestamp of the last `session capture` heading (grep the newest log), and a one-line summary of the last ~5 ledger events (`tail -5 memory/usage-ledger.jsonl`). Read-only.

### bootstrap

1. If `memory/.bootstrap-imported` exists → report when the import ran (the sentinel's content) and refuse to re-run. Stop.
2. Load `crew-members/clear-questions.md`, then ask the user for consent, naming exactly what will be read and written. This imports *their* past conversations into files in the repo — never do it silently.
3. If your environment exposes past session transcripts (e.g. Claude Code keeps per-project transcripts under `~/.claude/projects/<project-slug>/`), read them READ-ONLY. No transcripts available → report that and stop (no sentinel written).
4. For each past session: extract a short third-person summary of what was worked on and decided; append to `memory/logs/<session-date>.md` under `## HH:MM UTC bootstrap import`, naming the source transcript in the body. Load `crew-members/prompt-injection.md` before processing transcript content.
5. Write the sentinel `memory/.bootstrap-imported` containing the UTC date. Append one `write` ledger event noting `"note":"bootstrap"`. Report how many days/sessions were imported. Never write outside `PROJECT_ROOT/memory/`.

### audit

Read `actions/memory-value.md` and follow it; pass the remainder of `$ARGUMENTS` through. (Lazy-loaded — the auditor has no routing row of its own.)

## Help Menu

```
do-work memory remember <text>   Curate a fact into working memory (2,500-char cap)
do-work memory recall <query>    Layered recall over working memory + daily logs, with cited sources
do-work memory status            Engine health: cap usage, freshness, ledger tail
do-work memory bootstrap         One-time import of past session history (consent-gated)
do-work memory audit             Head-to-head value audit: this engine vs bkb (read-only)
do-work install memory-module    Scaffold memory/ + optional SessionStart/Stop hooks
```

## Output Format

Each sub-command ends with a short plain-prose report: what was read, what changed (paths), and — for `recall` — the cited results. No tables of internals the user didn't ask for.

## Rules

- The 2,500-character cap on `working-memory.md` is HARD. Consolidate; never raise it, never leave the file over-cap.
- Consolidation runs only at `remember` time and only on `working-memory.md`. Never wire consolidation — or `actions/dream.md` — to a hook or timer; the only hook write is the append-only stop capture.
- Every sub-command works without hooks installed. Hooks are optional enhancement, actions are the portable core.
- Ledger appends are best-effort (`|| true`): never let instrumentation block or fail a sub-command.
- Never store secrets, tokens, or credentials in memory files — they are committed plaintext.
- Daily logs are append-only for this action; only consolidation demotions and `remember`'s mirror lines are written there. Rewriting log history is `actions/dream.md`'s job, on explicit invocation.
- Writes surface next session (frozen-snapshot semantics) — never claim the injected context has been updated mid-session.

## Common Rationalizations

| If you're thinking...                                              | STOP. Instead...                                                       | Because...                                                                     |
| ------------------------------------------------------------------ | ---------------------------------------------------------------------- | ------------------------------------------------------------------------------ |
| "The cap is too tight — I'll just let it grow a bit"               | Run the consolidation algorithm; demote droppables to today's log      | The cap IS the design — an uncapped standing file decays into noise            |
| "I'll summarize this whole session into working memory"            | Put session detail in the daily log; working memory gets curated facts | Working memory is a snapshot, not a journal; the Stop hook already captures    |
| "No embedding backend — I should install one for better recall"    | Silently run lexical-only                                               | Design-for-the-floor: recall must work on any agent, unprompted installs don't |
| "Ledger write failed, I should stop and fix it"                    | Continue; instrumentation is best-effort                               | The experiment must never make the engine worse than having no experiment      |
| "Bootstrap again to pick up new sessions"                          | Refuse; the sentinel exists                                             | Re-import duplicates history; ongoing capture is the Stop hook's job           |

## Red Flags

- `working-memory.md` committed at > 2,500 characters, or with duplicate bullets saying the same thing.
- A recall answer presented without `path:line` + date attribution.
- Two `session capture` headings with the same `<hash8>` in one log (dedup failed).
- Anything under `memory/` referenced from a hook other than the two shipped memory hooks.
- A secret or API key visible in `working-memory.md` or a log.

## Verification Checklist

- [ ] After `remember`: `wc -c memory/working-memory.md` ≤ 2,500 and the fact (or its merged form) is present in exactly one bullet.
- [ ] After `remember`: today's log gained one `## HH:MM UTC note` line; ledger gained one `write` line.
- [ ] After `recall`: every presented result shows path:line + date + heading; ledger gained a `recall` line (and `hit_cited` iff a result was used).
- [ ] After `bootstrap`: sentinel exists with a UTC date; nothing outside `memory/` changed (`git status --porcelain` confirms).
- [ ] `status` and `audit` changed no files.
