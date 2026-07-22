# Memory Reference

> **Part of the do-work skill.** Companion to `actions/memory.md` — the canonical home for the memory engine's file schemas, the working-memory template, the lexical/semantic recall procedure, the consolidation algorithm, the usage-ledger contract (for BOTH engines), and the hook install internals used by `actions/install.md`'s `memory-module` target. This file has no routing surface; it is loaded by the files that cite it.

## Directory Schema

```
memory/                        # at project root (git rev-parse --show-toplevel, else pwd)
├── working-memory.md          # standing memory — HARD CAP 2,500 characters
├── logs/
│   └── YYYY-MM-DD.md          # dated daily logs, append-only (UTC dates)
├── usage-ledger.jsonl         # one JSON line per event (schema below)
└── .bootstrap-imported        # sentinel — exists after the one-time bootstrap import
```

All paths derive from `PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"`. Never write outside `PROJECT_ROOT/memory/`.

## working-memory.md Template

Created by `do-work install memory-module` (only if absent — never overwritten):

```markdown
---
updated: YYYY-MM-DD
---
<!-- Standing memory. HARD CAP: 2,500 characters total file size.
     Injected verbatim at session start when the memory hooks are installed.
     Curated by `do-work memory remember` — do not bulk-append; see actions/memory.md. -->

## Active Threads

## Notes

## Pending Decisions
```

Snapshot semantics: hooks inject this file (plus today's log) once, at session start. The injected copy is **frozen for the session** — writes made during a session land in the files and surface at the NEXT session start.

## Daily-Log Entry Conventions

Every entry in `memory/logs/YYYY-MM-DD.md` is a `##` heading followed by body text:

```
## HH:MM UTC <kind>
```

where `<kind>` is one of (illustrative, not exhaustive — new writers add new kinds without updating this list):

- `note` — one-liner mirrored by `memory remember`, or overflow moved out of working memory by consolidation.
- `session capture <hash8>` — appended by `hooks/memory-stop-capture.sh`; `<hash8>` is the first 8 hex chars of the sha256 of the captured text and is the dedup key.
- `bootstrap import` — written once by `memory bootstrap`; body must name the source transcript.

## Lexical Recall (Layer 1 — always runs)

Design-for-the-floor: grep + arithmetic only. Sanitize the query FIRST as a text operation (CLAUDE.md rule: never interpolate raw user text inside shell quoting), then substitute the already-safe value:

```bash
PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
MEMORY_DIR="$PROJECT_ROOT/memory"
# 1. Derive a sanitized token list from the user's query: lowercase, strip everything
#    but [a-z0-9 ], split on whitespace, drop tokens shorter than 3 chars.
#    Do this as a text operation on the query BEFORE building any command.
# 2. For each safe token:
#      grep -inH -- "$token" "$MEMORY_DIR/working-memory.md" "$MEMORY_DIR"/logs/*.md 2>/dev/null
# 3. Score each matching line: (number of distinct query tokens hitting that line) × recency weight.
#    Recency weight by source file:
#      working-memory.md                → 4 (curated beats raw)
#      logs/<date>.md, date ≤ 7 days old  → 3
#      logs/<date>.md, date ≤ 30 days old → 2
#      logs/<date>.md, older              → 1
#    (Log age comes from the filename's YYYY-MM-DD, not mtime.)
# 4. Emit the top 8 lines by score. EVERY result must carry attribution:
#    path:line, the log date (or "working memory"), and the nearest preceding ## heading.
```

Steps 3–4 are scoring/formatting the agent performs on the grep output — they need no additional shell state, so nothing carries between command blocks (CLAUDE.md: shell state does not survive between prescribed blocks).

## Semantic Recall (Layer 2 — optional, detected)

Probe for an embedding backend; first hit wins. The list is illustrative — any backend that can embed text qualifies:

```bash
ollama list 2>/dev/null | grep -qiE 'embed'   # a local embedding model is pulled
command -v embed >/dev/null 2>&1               # a standalone embed CLI
[ -n "${OPENAI_API_KEY:-}${VOYAGE_API_KEY:-}" ] # an embeddings API key is exported
```

If a backend is found: chunk candidates by daily-log `##` headings (working-memory.md is one chunk per section), embed query + chunks, rank by cosine similarity, then **merge with the lexical results by reciprocal-rank fusion** (score = Σ 1/(60 + rank) across both lists) and keep each chunk's attribution. If no backend is found: silently proceed lexical-only — same graceful degradation as `actions/board.md` without Go. Never install, download, or prompt for a backend from inside `recall`.

## Consolidation Algorithm (the 2,500-char cap)

Runs inside `memory remember` when a write would push `working-memory.md` over 2,500 characters (`wc -c`):

1. Read the whole file. Group bullets by section.
2. Merge duplicates and near-duplicates into single bullets; a superseded fact is replaced, not kept alongside its replacement.
3. Still over cap → move the lowest-value droppables (resolved threads, stale notes, decided decisions) into today's log as `## HH:MM UTC note` entries — **consolidation never destroys content, it demotes it to the log**.
4. Still over cap → tighten wording of survivors. The new fact always fits; what leaves is the oldest resolved material.
5. Update the `updated:` frontmatter date. Verify `wc -c` ≤ 2,500 before finishing.

The cap is the design, not an obstacle: it forces the standing memory to stay high-signal. Never raise it; never let the file commit over-cap.

## Usage-Ledger Contract (canonical — both engines)

Files: `memory/usage-ledger.jsonl` (memory engine) and `usage-ledger.jsonl` at the KB root (bkb — canonically `kb/usage-ledger.jsonl`; the root is resolved per `actions/bkb.md`'s "Locating the Knowledge Base"). Append-only, one JSON object per line, no trailing commas, UTC timestamps. Writers: `actions/memory.md`, `actions/bkb.md` (query step 8b, ingest step 7b), `hooks/memory-session-start.sh`, `hooks/memory-stop-capture.sh` — plus any future surface that reads from or writes to either engine (the trigger condition is "an engine event occurred", not membership in this list).

```json
{"ts":"2026-07-22T18:04:11Z","engine":"memory","event":"recall","query":"dark mode decision","hits":3,"source":"actions/memory.md","note":""}
```

| Field | Value |
| --- | --- |
| `ts` | `date -u +%Y-%m-%dT%H:%M:%SZ` |
| `engine` | `memory` \| `bkb` |
| `event` | memory: `inject`, `capture`, `write`, `recall`, `hit_cited` · bkb: `query`, `ingest`, `hit_cited` (illustrative — new events allowed, auditor buckets unknown events as "other") |
| `query` | recall/query events only; sanitized token form (same text-operation sanitize as lexical recall), never raw user text |
| `hits` | integer result/page count; 0 when not applicable |
| `source` | file path of the writer (e.g. `hooks/memory-stop-capture.sh`) |
| `note` | free text, usually empty |

Prescribed append (derive-then-substitute; `$safe_query` and `$hit_count` are already-sanitized values):

```bash
utc_now="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
printf '{"ts":"%s","engine":"memory","event":"recall","query":"%s","hits":%d,"source":"actions/memory.md","note":""}\n' \
  "$utc_now" "$safe_query" "$hit_count" >> "$PROJECT_ROOT/memory/usage-ledger.jsonl" 2>/dev/null || true
```

**Ledger appends are best-effort.** The `|| true` is mandatory in every writer — a missing directory, read-only checkout, or full disk must never block the action being instrumented. `hit_cited` is the event that matters most to `actions/memory-value.md`'s verdict: append it whenever a recalled/queried result is actually used in the reply, one line per recall that produced a cited result (not per result).

## Stop-Capture Hash Dedup Spec

Used by `hooks/memory-stop-capture.sh`:

1. `capture_text` = final user message + final assistant message from the session transcript, truncated to ~1,500 characters total.
2. `hash8="$(printf '%s' "$capture_text" | sha256sum | cut -c1-8)"` (fall back to `shasum -a 256` on systems without `sha256sum`).
3. `grep -q "session capture $hash8" "$today_log"` → already captured, exit 0 (idempotent across duplicate Stop firings).
4. Append heading `## HH:MM UTC session capture <hash8>`, framing line `Session capture — final exchange between the user and the agent:`, then the text.
5. The hook ALWAYS exits 0 — capture is never worth blocking a session end.

## Hook Install Internals (used by actions/install.md → memory-module)

`hooks/memory-hooks.json` is a fragment shaped exactly like `hooks/hooks.json`. The install target **appends** its entries into the consumer's `.claude/settings.json` — compose, never clobber:

```bash
settings_file="$PROJECT_ROOT/.claude/settings.json"
[ -f "$settings_file" ] || printf '{}\n' > "$settings_file"
if grep -q 'memory-session-start.sh' "$settings_file"; then
  echo "memory hooks already present — skipping merge"
else
  cp "$settings_file" "$settings_file.pre-memory-module"   # backup BEFORE touching anything
  merged_settings="$settings_file.merge-tmp"
  jq --slurpfile frag "<skill-root>/hooks/memory-hooks.json" \
     '.hooks.SessionStart = ((.hooks.SessionStart // []) + $frag[0].hooks.SessionStart)
    | .hooks.Stop        = ((.hooks.Stop        // []) + $frag[0].hooks.Stop)' \
     "$settings_file" > "$merged_settings" && mv "$merged_settings" "$settings_file"
fi
```

- Dedup gate = grep for the script filename; append via `+`, never assign a whole new array over `.hooks.SessionStart`/`.hooks.Stop`.
- After the merge, verify: the file still parses (`jq . "$settings_file" >/dev/null`), both memory hook filenames are present, and every pre-existing hook entry is still there (compare entry counts against the backup). Parse failure → restore from `$settings_file.pre-memory-module` and report a broken install. Success → remove the backup.
- No `jq` → do NOT attempt a sed/awk merge. Print the two entries from `hooks/memory-hooks.json` with instructions to merge manually, and report `hooks: MANUAL STEP` — a warning, not a failure. Every `do-work memory` sub-command works without hooks; hooks are the Claude Code-specific enhancement, the actions are the portable core.

Uninstall (manual, documented here for symmetry): remove the two entries whose command contains `memory-session-start.sh` / `memory-stop-capture.sh` from `.claude/settings.json`; `memory/` itself is user data — never delete it as part of hook removal.
