# Memory Action

> **Part of the do-work-knowledge skill.** Hermes-style session memory: a capped standing `memory/working-memory.md`, dated daily logs, and layered recall with source attribution. Invoked by `do-work-knowledge memory <sub-command>`, `remember <text>`, or `recall <query>`. It lives beside BKB in the knowledge package because the two engines share instrumentation and `actions/memory-value.md` compares them from one vantage point; optional hook setup is owned by `actions/setup-memory.md`.

**Philosophy.** This engine optimizes *capture*: never lose anything, at zero effort. The standing file is deliberately tiny (hard cap 2,500 characters) so it stays high-signal; everything else flows to append-only daily logs. It runs in parallel with `actions/bkb.md` (which optimizes *synthesis*) during the ADR-017 experiment — both engines write usage ledgers, and `memory audit` renders the head-to-head. Hooks (session-start injection, stop capture) are an optional Claude Code enhancement installed by `do-work-knowledge setup-memory`; **every sub-command below must work with no hooks installed.**

## When to Use

**Use when:**
- The user says "remember this/that …", "note for next time", or asks "what do you remember about …".
- The user wants session-persistent context without curating a wiki.
- The user asks for memory status, a one-time import of past session history (`bootstrap`), or the engine head-to-head (`audit`).

**Do NOT use when:**
- The content is wiki-grade knowledge to compile and interlink → `actions/bkb.md`.
- The user wants to consolidate/prune an arbitrary memory or wiki directory → `actions/dream.md`.
- The user is queueing work ("remember to fix X" = a task, not a fact) → `../../do-work/actions/capture.md`.

## Input

`$ARGUMENTS` = `<sub-command> [payload]`. Bare `memory` (or `memory help`) → print the Help Menu below.

**The sub-command is always present, including via the direct aliases.** `do-work-knowledge remember <text>`, `do-work-knowledge forget <text>`, and `do-work-knowledge recall <query>` route here with the alias preserved as the sub-command (`remember <text>` / `forget <text>` / `recall <query>`), and `what do you remember` arrives as `recall <query>` — the package router's memory row specifies that, because the router's default "strip the trigger, pass the rest" would otherwise hand this action a bare payload with nothing to dispatch on. If `$ARGUMENTS` ever arrives without a leading sub-command, **fall back to `recall`, never to `remember`** — and say which sub-command you assumed. Do not infer from sentence shape: real recall queries are usually noun phrases, not questions (`recall deployment decision` arrives as `deployment decision`), so a statement-shaped test would classify most reads as writes and silently mutate the store on a request to read it. The two errors are not symmetric — a wrong `recall` shows the user something they didn't ask for, a wrong `remember` writes a fact nobody asserted into standing memory. When the payload genuinely reads as something to store and an interactive prompt is available, ask which was meant (`crew-members/clear-questions.md`); with no prompt available, recall and say so. This is a fallback for a router bug, not the contract.

Locate the store first:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
MEMORY_DIR="$PROJECT_ROOT/memory"
```

If `$MEMORY_DIR/working-memory.md` is missing for any sub-command except `audit`: report that the memory engine isn't set up here and point at `do-work-knowledge setup-memory`, then stop.

## Sub-Commands

| Sub-command | Payload | What it does |
| --- | --- | --- |
| `remember` | the fact to store | Curate into `working-memory.md` (dedup, supersede, cap-enforce) + mirror to today's log |
| `forget` | the fact to remove | Confirmation-gated: remove the `working-memory.md` bullet AND redact matching daily-log lines in place |
| `recall` | the query | Layered search (lexical always; semantic when a backend is detected) with cited sources |
| `status` | — | Engine health: cap usage, freshness, log days, ledger tail |
| `bootstrap` | — | One-time, consent-gated import of prior session history into dated logs |
| `audit` (alias `value`) | optional focus | Lazy-loads the engine-vs-engine value auditor |

### remember <text>

1. Read the WHOLE `working-memory.md` first — never blind-append.
2. Place the fact in the right section (`## Active Threads` / `## Notes` / `## Pending Decisions`). If it duplicates an existing bullet, merge; if it supersedes one ("we now use X instead of Y"), replace the old bullet in place. An explicit ask to *forget* something is not a `remember` — route it to the `forget` sub-command below; supersede-in-place replaces a fact with its successor, `forget` erases a fact outright, and only the latter must also reach the logs.
3. Enforce the hard cap: if the file would exceed **2,500 characters** (`wc -c`), run the consolidation algorithm in `actions/memory-reference.md` — merge, then demote droppables to today's log, then tighten. The file must end ≤ 2,500 chars.
4. Mirror a one-liner to `memory/logs/<UTC date>.md` — the date-only shape in the Timestamp rule (`../../do-work/actions/work-reference.md`) — under a `## HH:MM UTC note` heading (create the file if needed).
5. Update the `updated:` frontmatter date. Append a `write` ledger event per `actions/memory-reference.md` (best-effort, never blocking).
6. Tell the user what was stored, and what (if anything) was merged, replaced, or demoted. Remind them once per conversation that writes surface at the NEXT session start (the injected snapshot is frozen).

### forget <text>

Removing the working-memory bullet alone is not forgetting: recall's Layer 1 searches the daily logs too, so a fact left there stays recallable forever. `forget` is the one named exception to the logs-are-append-only rule — explicit user invocation only, never an automatic writer, and nothing is touched before the user confirms.

1. Locate the fact everywhere: matching bullet(s) in `working-memory.md` AND matching lines in `memory/logs/*.md` (tokenize the payload with the same sanitization recall uses — `actions/memory-reference.md`, Lexical Recall). No match anywhere → report that and stop.
2. Load `crew-members/clear-questions.md`, then show the user exactly what would be removed or redacted — every matched line with its file and date — and ask for confirmation. Partial confirmation is fine (forget the bullet, keep a log line); redact only what was confirmed.
3. Remove the confirmed bullet(s) from `working-memory.md` and update its `updated:` frontmatter date.
4. Redact the confirmed log lines **in place**: replace each line's content with `[forgotten — redacted by memory forget YYYY-MM-DD]`. Never delete the line outright — the log must still show something stood there. Two structural constraints: a line inside a capture body keeps its `> ` prefix (the quoting contract in `actions/memory-reference.md` → Daily-Log Entry Conventions is what makes capture boundaries unspoofable), and `##` heading lines are never modified — a `session capture <hash8>` heading is the dedup key.
5. Append a `write` ledger event with `"note":"forget"` (best-effort, never blocking). Report what was removed and where, and state the known limit: `working-memory.md` is committed, so a forgotten bullet still exists in git history — `forget` scrubs the store, not the repo's past.

### recall <query>

**Empty query → broad recall, not an error and not a no-op.** `do-work-knowledge recall` with no payload, and the `what do you remember` phrasing that routes here with nothing after the verb (SKILL.md row 37), are asking *what's in there* rather than searching for a term. There is nothing to tokenize, so skip the search layers entirely: present all of `working-memory.md`, then the curated entries (`## HH:MM UTC` sections that are not `session capture` — same exclusion the session-start hook applies) from the most recent 3 log days, newest first, still citing every source per step 5. If that exceeds roughly 40 lines, summarize the older days and say how many entries were folded. Log the ledger event with an empty query string. Steps 2–4 below do not apply.

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
do-work-knowledge memory remember <text>   Curate a fact into working memory (2,500-char cap)
do-work-knowledge memory forget <text>     Confirmation-gated: remove a fact from working memory + redact it from the logs
do-work-knowledge memory recall <query>    Layered recall over working memory + daily logs, with cited sources
do-work-knowledge memory status            Engine health: cap usage, freshness, ledger tail
do-work-knowledge memory bootstrap         One-time import of past session history (consent-gated)
do-work-knowledge memory audit             Head-to-head value audit: this engine vs bkb (read-only)
do-work-knowledge setup-memory    Scaffold memory/ + optional SessionStart/Stop hooks
```

## Output Format

Each sub-command ends with a short plain-prose report: what was read, what changed (paths), and — for `recall` — the cited results. No tables of internals the user didn't ask for.

## Rules

- The 2,500-character cap on `working-memory.md` is HARD. Consolidate; never raise it, never leave the file over-cap.
- Consolidation runs only at `remember` time and only on `working-memory.md`. Never wire consolidation — or `actions/dream.md` — to a hook or timer; the only hook write is the append-only stop capture.
- Every sub-command works without hooks installed. Hooks are optional enhancement, actions are the portable core.
- Ledger appends are best-effort (`|| true`): never let instrumentation block or fail a sub-command.
- Never store secrets, tokens, or credentials in memory files. `working-memory.md` is committed plaintext; `memory/logs/` and `memory/usage-ledger.jsonl` are machine-local (git-excluded by `setup-memory`) but still plaintext on disk — local-only is not a licence to write a credential there.
- Daily logs are append-only for this action, with one named exception: `forget` — on explicit user invocation, after confirmation — redacts matching log lines in place (marker, not deletion). Automatic writers (hooks, consolidation demotions, `remember`'s mirror lines) only ever append. Rewriting log history beyond that is `actions/dream.md`'s job, on explicit invocation.
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
- A log line rewritten by anything other than a user-invoked, confirmed `forget` — or a `forget` that deleted lines outright instead of leaving the redaction marker.
- `memory/logs/` or `memory/usage-ledger.jsonl` showing up in `git status` as untracked — the local-ignore step of `setup-memory` was skipped or reverted; verbatim captures are one `git add -A` from the repo.

## Verification Checklist

- [ ] After `remember`: `wc -c memory/working-memory.md` ≤ 2,500 and the fact (or its merged form) is present in exactly one bullet.
- [ ] After `remember`: today's log gained one `## HH:MM UTC note` line; ledger gained one `write` line.
- [ ] After `forget`: the fact appears in no `working-memory.md` bullet and no un-redacted log line; every redacted line holds the marker (capture-body lines still `> `-prefixed); no `##` heading line changed; ledger gained a `write` line noting `forget` — and the user confirmed the exact lines shown before any write.
- [ ] After `recall`: every presented result shows path:line + date + heading; ledger gained a `recall` line (and `hit_cited` iff a result was used).
- [ ] After `bootstrap`: sentinel exists with a UTC date; nothing outside `memory/` changed (`git status --porcelain` confirms).
- [ ] `status` and `audit` changed no files.
