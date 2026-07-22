# Memory Value Action

> **Part of the do-work skill.** Engine-agnostic value auditor for the ADR-017 parallel-memory experiment: scans the current repo for evidence that each memory engine — `actions/bkb.md`'s `kb/` and `actions/memory.md`'s `memory/` — is actually used and providing value, then renders a head-to-head verdict. Loaded lazily by `do-work memory audit` (alias `memory value`); it has no routing row of its own. Justified with the rest of the trio in `decisions/records/adr-017-run-a-parallel-memory-engine-experiment-with-usage-ledgers.md`.

**Strictly read-only.** This action never modifies, moves, or deletes anything — not even ledger appends. It reads files and git history, and it reports. Retiring an engine is a human decision made on this report, executed later as a normal maintenance pass.

## When to Use

**Use when:**
- The user asks whether bkb (or the memory engine) is actually being used / worth keeping ("determine current value").
- The ADR-017 experiment window (~4 weeks of ledger data) is up and a verdict is due.

**Do NOT use when:**
- The user wants engine *content* health (broken links, orphan pages) → `do-work bkb lint` / `do-work memory status`.
- The user wants to consolidate or prune → `actions/dream.md` (manual, destructive).

## Input

`$ARGUMENTS` (optional): `bkb` or `memory` to audit one engine only; empty → audit both + head-to-head. All commands run from `PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"`.

## Checks

Each check is an independent probe; run all that apply to the engines in scope. An absent directory short-circuits that engine to **Absent** (run no further probes for it).

### bkb engine (`kb/`)

- **Existence & shape:** `kb/` present? `find kb/wiki -name '*.md' 2>/dev/null | wc -l` wiki pages; `kb/raw/` inbox size.
- **Git activity:** `git log --oneline -- kb/ | wc -l` total commits; `git log -1 --format=%ci -- kb/` last-touched; `git log --format=%an -- kb/ | sort -u` distinct authors (human-touch signal).
- **Log activity:** entries in `kb/wiki/log.md` dated within the last 30 and 90 days (grep date headings).
- **Inbound references:** wiki pages cited from outside `kb/` — e.g. `grep -rl 'kb/wiki/' --include='*.md' "$PROJECT_ROOT" | grep -v "^$PROJECT_ROOT/kb/"` and `[[wikilink]]` mentions outside `kb/`. A wiki nobody links to is write-only.
- **Ledger stats:** `kb/usage-ledger.jsonl` per the shared procedure below.

### memory engine (`memory/`)

- **Existence & shape:** `memory/working-memory.md` present? Character count vs the 2,500 cap; section fill (are the three `##` sections non-empty?); `updated:` frontmatter date.
- **Log cadence:** `ls memory/logs/*.md 2>/dev/null | wc -l` distinct days; newest log date; ratio of `session capture` headings to `note` headings (all-capture/no-note = nobody curates; all-note/no-capture = hooks not firing).
- **Hook wiring:** are `memory-session-start.sh` / `memory-stop-capture.sh` referenced in `.claude/settings.json`? (Unwired hooks explain an empty capture trail — an instrumentation gap, not absent value.)
- **Ledger stats:** `memory/usage-ledger.jsonl` per the shared procedure below.

### Ledger stats (shared procedure, both engines)

From the engine's `usage-ledger.jsonl` (tolerate absence — report "no ledger"):
- Events per week by `event` type over the trailing 4 weeks (bucket unknown event names as "other").
- Retrieval count: `recall` (memory) / `query` (bkb) events.
- **Hit-cited rate** = `hit_cited` ÷ retrieval count — the value signal.
- Age of the newest event.

### Fairness rule (mandatory)

bkb predates instrumentation: **absence of ledger evidence is not absence of value.** For bkb's pre-ledger window, weight `git log -- kb/` and `kb/wiki/log.md` history as the usage record, and say explicitly in the report which window each conclusion draws on. Conversely, don't credit the memory engine's `inject`/`capture` volume as value — those are automatic; only retrieval and citation count.

## Output Format

```
# Memory Engine Value Audit — <date>

## bkb (kb/)          → <Active | Idle | Stale | Absent>
<probe results, 1 line each>

## memory (memory/)   → <Active | Idle | Stale | Absent>
<probe results, 1 line each>

## Head-to-head
| Signal | bkb | memory |
| events/week (trailing 4w) | … | … |
| retrievals (recall/query) | … | … |
| hit-cited rate            | … | … |
| freshness (last activity) | … | … |
| human touch (git authors) | … | … |

## Verdict
<one paragraph: which engine is winning on CITED RETRIEVALS, what each is actually
being used for, any instrumentation gap found, and a recommendation —
keep both / retire one / fix instrumentation and re-audit on <date>.>
```

Classification rubric: **Active** = ≥3 non-automatic events (excluding `inject`/`capture`) in the last 14 days AND ≥1 `hit_cited` (or, for bkb pre-ledger, equivalent git/log.md evidence of use). **Idle** = structure exists, below the Active bar. **Stale** = no activity of any kind > 30 days. **Absent** = directory missing.

## Rules

- Read-only, no exceptions — this audit appends nothing, not even to the ledgers it reads.
- The verdict weighs cited retrievals, never raw write volume — a store nobody reads from is a landfill, not memory.
- Never recommend deletion as a done deal; the recommendation names the human decision and the maintenance path (`crew-members/maintenance.md`).
- Apply the fairness rule in every both-engine audit, and state it in the report.

## Red Flags

- The audit left `git status --porcelain` non-clean.
- A verdict that crowns a winner on `inject`/`capture`/`write` volume.
- bkb classified Stale/Absent purely from a missing ledger while `git log -- kb/` shows recent commits.

## Verification Checklist

- [ ] Report rendered with both engine sections (or the one requested), the head-to-head table, and a verdict paragraph.
- [ ] `git status --porcelain` unchanged by the audit.
- [ ] Every classification cites the probe evidence it rests on.
