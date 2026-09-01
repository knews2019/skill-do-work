# Memory Reference

> **Part of the do-work-knowledge skill.** Companion to `actions/memory.md` — the canonical home for the memory engine's file schemas, the working-memory template, the lexical/semantic recall procedure, the consolidation algorithm, the usage-ledger contract (for BOTH engines), and the hook install internals used by `actions/setup-memory.md`'s `memory-module` target. This file has no routing surface; it is loaded by the files that cite it.

## Directory Schema

```
memory/                        # at project root (git rev-parse --show-toplevel, else pwd)
├── working-memory.md          # standing memory — HARD CAP 2,500 characters  [committable]
├── logs/
│   └── YYYY-MM-DD.md          # dated daily logs, append-only (UTC dates; sole exception: user-invoked `memory forget` redacts in place)    [machine-local]
├── usage-ledger.jsonl         # one JSON line per event (schema below)       [machine-local]
└── .bootstrap-imported        # sentinel — exists after the one-time bootstrap import [machine-local]
```

All paths derive from `PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"`. Never write outside `PROJECT_ROOT/memory/`.

`working-memory.md` is the only committable file in the store; the other three are added to `.git/info/exclude` by `actions/setup-memory.md`'s `memory-module` Phase 2. The sentinel is machine-local because `memory bootstrap` refuses to re-run when it exists — committed, it would block every other clone from importing that machine's own session history.

## working-memory.md Template

Created by `do-work-knowledge setup-memory` (only if absent — never overwritten):

```markdown
---
updated: YYYY-MM-DD
---
<!-- Standing memory. HARD CAP: 2,500 characters total file size.
     Injected verbatim at session start when the memory hooks are installed.
     Curated by `do-work-knowledge memory remember` — do not bulk-append; see actions/memory.md. -->

## Active Threads

## Notes

## Pending Decisions
```

Snapshot semantics: the core CLI's `memory-session-start` command injects this file (plus today's log's **curated entries only** — `session capture` sections are excluded, because they are verbatim transcript text that must not enter context before a prompt-injection guard can load; they stay reachable via `do-work-knowledge memory recall`) once, at session start. `hooks/memory-session-start.sh` is only the retained event launcher. The injected copy is **frozen for the session** — writes made during a session land in the files and surface at the NEXT session start.

## Daily-Log Entry Conventions

Every entry in `memory/logs/YYYY-MM-DD.md` is a `##` heading followed by body text:

```
## HH:MM UTC <kind>
```

where `<kind>` is one of (illustrative, not exhaustive — new writers add new kinds without updating this list):

- `note` — one-liner mirrored by `memory remember`, or overflow moved out of working memory by consolidation.
- `session capture <hash8>` — appended by the core CLI's `memory-stop-capture` command through retained launcher `hooks/memory-stop-capture.sh`; `<hash8>` is the first 8 hex chars of the sha256 of the captured text and is the dedup key.
- `bootstrap import` — written once by `memory bootstrap`; body must name the source transcript.

The heading's `HH:MM` is the entry's UTC **time-of-day label, outside the Timestamp rule's scope** (`../../do-work/actions/work-reference.md`): it is neither an instant nor a date-only stamp, because the log's dated filename already carries the date. A timestamp sweep walks past every `## HH:MM UTC` heading; the write sites point here.

**A `session capture` section's end is decided by format, not by heading grammar.** Raw capture text can contain any line at all, including `## 12:34 UTC note` — heading grammar is trivially spoofable by anything that reaches a transcript, so it cannot on its own end a section that `hooks/memory-session-start.sh` is suppressing. The contract that makes the boundary unspoofable:

- A writer emitting verbatim third-party text MUST open the section with the sentinel `<!-- do-work:capture-body quoted -->` as its first non-blank line, and MUST `> `-prefix every body line after it (including the framing line).
- The reader then ends the section at the first heading-grammar line that is **not** `> `-prefixed. Because every body line is quoted, no body line can be that.
- A section without the sentinel is pre-0.139.4 legacy with an unquoted body, where no boundary is trustworthy: the reader suppresses to end-of-file. That also hides curated entries written later that day — a bounded, self-clearing cost, correct against injecting raw transcript text.

Both halves are required. Dropping the sentinel makes legacy and current sections indistinguishable; dropping the quoting lets a body line impersonate the boundary.

**One writer may rewrite an existing body line:** `memory forget` (`actions/memory.md`) — explicit user invocation, confirmation-gated — replaces a confirmed line's content with a `[forgotten — redacted by memory forget YYYY-MM-DD]` marker. It preserves any `> ` prefix and never modifies `##` heading lines, so both halves of the boundary contract (and the `session capture <hash8>` dedup key) survive redaction. Automatic writers only ever append.

## Lexical Recall (Layer 1 — always runs)

`do-work-cli memory-recall` is the executable authority for tokenization, scoring, recency weighting, attribution, sorting, and the eight-result bound, including the canonical [Raw text before shell quoting](../../do-work/docs/prescribed-shell-primitives.md#raw-text-before-shell-quoting) contract. The retained helper below is a compatibility surface until the whole-suite shim migration; actions must not call it after a canonical-command failure:

```bash
<skill-root>/scripts/lexical-memory-recall.sh "$(git rev-parse --show-toplevel 2>/dev/null || pwd)/memory" "<raw query text>"
```

The script owns tokenization, scoring, recency weighting, attribution, sorting, and the eight-result bound; the caller owns only query intent and how returned memories are used.

## Semantic Recall (Layer 2 — optional, detected)

Probe for an embedding backend; first hit wins. The list is illustrative — any backend that can embed text qualifies:

```bash
ollama list 2>/dev/null | grep -qiE 'embed'   # a local embedding model is pulled
command -v embed >/dev/null 2>&1               # a standalone embed CLI
[ -n "${OPENAI_API_KEY:-}${VOYAGE_API_KEY:-}" ] # an embeddings API key is exported
```

If a backend is found: chunk candidates by daily-log `##` headings (working-memory.md is one chunk per section), embed query + chunks, rank by cosine similarity, then **merge with the lexical results by reciprocal-rank fusion** (score = Σ 1/(60 + rank) across both lists) and keep each chunk's attribution. If no backend is found: silently proceed lexical-only — same graceful degradation as `../../do-work-board/actions/board.md` without Go. Never install, download, or prompt for a backend from inside `recall`.

## Consolidation Algorithm (the 2,500-char cap)

Runs inside `memory remember` when a write would push `working-memory.md` over 2,500 characters (`wc -c`):

1. Read the whole file. Group bullets by section.
2. Merge duplicates and near-duplicates into single bullets; a superseded fact is replaced, not kept alongside its replacement.
3. Still over cap → move the lowest-value droppables (resolved threads, stale notes, decided decisions) into today's log as `## HH:MM UTC note` entries (time-of-day label, outside the Timestamp rule's scope — Daily-Log Entry Conventions above) — **consolidation never destroys content, it demotes it to the log**.
4. Still over cap → tighten wording of survivors. The new fact always fits; what leaves is the oldest resolved material.
5. Update the `updated:` frontmatter date. Verify `wc -c` ≤ 2,500 before finishing.

The cap is the design, not an obstacle: it forces the standing memory to stay high-signal. Never raise it; never let the file commit over-cap.

## Usage-Ledger Contract (canonical — both engines)

Files: `memory/usage-ledger.jsonl` (memory engine) and `usage-ledger.jsonl` at the KB root (bkb — canonically `kb/usage-ledger.jsonl`; the root is resolved per `actions/bkb.md`'s "Locating the Knowledge Base"). Append-only, one JSON object per line, no trailing commas, UTC timestamps. Writers: `actions/memory.md`, `actions/bkb.md` (query step 8b, ingest step 7b), and the core CLI's `memory-session-start` / `memory-stop-capture` Go commands under `internal/hookcommands` — plus any future surface that reads from or writes to either engine (the trigger condition is "an engine event occurred", not membership in this list). The retained `.sh` hook paths are launch-only compatibility surfaces.

```json
{"ts":"2026-07-22T18:04:11Z","engine":"memory","event":"recall","query":"dark mode decision","hits":3,"source":"actions/memory.md","note":""}
```

| Field | Value |
| --- | --- |
| `ts` | current UTC instant (Timestamp rule, `../../do-work/actions/work-reference.md`) |
| `engine` | `memory` \| `bkb` |
| `event` | memory: `inject`, `capture`, `write`, `recall`, `hit_cited` · bkb: `query`, `ingest`, `hit_cited` (illustrative — new events allowed, auditor buckets unknown events as "other") |
| `query` | recall/query events only; sanitized token form (same text-operation sanitize as lexical recall), never raw user text |
| `hits` | integer result/page count; 0 when not applicable |
| `source` | stable event-source identity (hook events retain values such as `hooks/memory-stop-capture.sh` for compatibility; this is not a writer-file or implementation-owner pointer) |
| `note` | free text, usually empty |

Prescribed append (derive-then-substitute; `$utc_now`, `$safe_query` and `$hit_count` are already-derived values — `$utc_now` per the Timestamp rule, `../../do-work/actions/work-reference.md`):

```bash
printf '{"ts":"%s","engine":"memory","event":"recall","query":"%s","hits":%d,"source":"actions/memory.md","note":""}\n' \
  "$utc_now" "$safe_query" "$hit_count" >> "$PROJECT_ROOT/memory/usage-ledger.jsonl" 2>/dev/null || true
```

**Ledger appends are best-effort.** The `|| true` is mandatory in every writer — a missing directory, read-only checkout, or full disk must never block the action being instrumented. `hit_cited` is the event that matters most to `actions/memory-value.md`'s verdict: append it whenever a recalled/queried result is actually used in the reply, one line per recall that produced a cited result (not per result).

## Stop-Capture Hash Dedup Spec

Used by the core CLI's `memory-stop-capture` command. The retained `hooks/memory-stop-capture.sh` path only launches it and preserves the nonblocking Stop boundary. Transcript selection is the full intended JSON behavior on every host; jq availability never changes capture semantics.

1. Extract the final user message + final assistant message from the session transcript, **untruncated**. "Final message" means the last transcript entry that carries real text: pull only from blocks typed `text`, skip `isMeta` entries, and drop entries whose extracted text is blank. Claude Code records tool results as `type: "user"` entries holding `tool_result` blocks with no `.text`, so a naive `last` lands on a tool result and stores an empty `User:` side for any session whose final turn used a tool — the common case, not an edge case. The assistant side has the mirror problem when a turn ends in a `tool_use` block.
2. Redact **both full extracted sides** per the Capture Redaction Spec below — redaction runs BEFORE truncation, not just before hashing. Every redaction pattern needs a complete token shape; a byte-budget cut through the middle of a token leaves a fragment (`ghp_1234567`) that no longer matches any pattern and would persist unredacted. The private-key drop is likewise judged on the full text. (Redact-before-hash follows for free: the dedup key is computed over already-redacted text.)
3. `capture_text` = the redacted sides composed as `User: …\n\nAgent: …`, truncated to ~1,500 characters total. Every byte-budget cut is piped through `iconv -c -f UTF-8 -t UTF-8` (plain cut when `iconv` is absent) so a mid-character cut in multi-byte text — routine for CJK at ~3 bytes/char — drops the torn trailing sequence instead of feeding invalid bytes into the persisted text and the dedup hash. A cut can at worst clip a `[REDACTED]` marker, never expose a credential fragment.
4. `hash8="$(printf '%s' "$capture_text" | sha256sum | cut -c1-8)"` (fall back to `shasum -a 256` on systems without `sha256sum`).
5. `grep -q "session capture $hash8" "$today_log"` → already captured, exit 0 (idempotent across duplicate Stop firings).
6. Append, in this order: heading `## HH:MM UTC session capture <hash8>` (time-of-day label, outside the Timestamp rule's scope — Daily-Log Entry Conventions above), the sentinel `<!-- do-work:capture-body quoted -->`, the framing line `> Session capture — final exchange between the user and the agent:`, then the text **with every line prefixed `> `**. Sentinel first, everything after it quoted — see Daily-Log Entry Conventions for why both are required to make the section boundary unspoofable. Quoting happens at write time only; the hash is computed over the unquoted text at step 4, so the dedup key stays stable. **The whole section lands in one `write()`** — compose it first and append with a single `printf`, never a multi-statement block: every session in the project appends to the same daily log, and separate `O_APPEND` writes from two near-simultaneous stops can interleave and garble section structure. A lock is not the fix here — it would fight step 7's never-block contract, and `flock` doesn't exist on macOS.
7. The hook ALWAYS exits 0 — capture is never worth blocking a session end.

## Capture Redaction Spec

The memory store is split by durability: the curated `working-memory.md` is **committed plaintext**, while `memory/logs/` and `memory/usage-ledger.jsonl` are **machine-local** — `actions/setup-memory.md`'s memory-module Phase 2 adds them to the repo's `.git/info/exclude`. That split is the first barrier, not the only one, and redaction is defense in depth behind it: the exclude entries exist only where the installer ran (a hand-scaffolded `memory/` has none), logs stay plaintext on disk where they can be read, grepped, or pasted elsewhere, and a curated fact promoted out of a log lands in the committed file. A verbatim capture must therefore never persist a credential in the first place. The core CLI's `internal/hookcommands` memory Stop owner:

- **Drops the whole capture** (exit 0) when the full extracted text contains a `PRIVATE KEY-----` block marker — key material spans lines, so line-based redaction can't be trusted.
- **Replaces credential-shaped substrings with `[REDACTED]`** in Go, applied to the full extracted messages **before truncation** (and therefore before hashing or writing) — every pattern needs a complete token shape, so redacting after a byte-budget cut would let a token severed by the cut persist as an unmatched fragment. Shipped patterns (illustrative, not exhaustive — the trigger condition is "text shaped like a credential", and `memory remember` curation remains the real gate): GitHub tokens (`ghp_…`, `github_pat_…`), `sk-…` API keys, AWS `AKIA…` key ids, Slack `xox?-…` tokens, `eyJ…`-prefixed JWTs, `Bearer <token>` headers, and `password/passwd/secret/token/api_key = <value>` assignments.
- **Skips the capture entirely if the redaction pipeline itself fails** — the unredacted text is never the fallback.

## Hook Install Internals (used by actions/setup-memory.md → memory-module)

`hooks/memory-hooks.json` is a fragment shaped exactly like `../../do-work/hooks/hooks.json`. The shipped installer appends only missing entries into the consumer's `.claude/settings.json` — compose, never clobber:

```bash
<skill-root>/scripts/install-memory-hooks.sh "$PROJECT_ROOT" "<skill-root>/hooks/memory-hooks.json"
```

- Dedup gate = one grep per script filename, each appending only its own missing entry; append via `+`, never assign a whole new array over `.hooks.SessionStart`/`.hooks.Stop`.
- After the merge, verify: the file still parses (`jq . "$settings_file" >/dev/null`), both memory hook filenames are present, and every pre-existing hook entry is still there (compare entry counts against the backup). Parse failure → restore from `$settings_file.pre-memory-module` and report a broken install. Success → remove the backup.
- No `jq` → do NOT attempt a sed/awk merge. Print the two entries from `hooks/memory-hooks.json` with instructions to merge manually, and report `hooks: MANUAL STEP` — a warning, not a failure. Every `do-work-knowledge memory` sub-command works without hooks; hooks are the Claude Code-specific enhancement, the actions are the portable core.

Uninstall (manual, documented here for symmetry): remove the two entries whose command contains `memory-session-start.sh` / `memory-stop-capture.sh` from `.claude/settings.json`; `memory/` itself is user data — never delete it as part of hook removal.
