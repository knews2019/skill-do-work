# Memory Value Action

> **Part of the do-work-knowledge skill.** Engine-agnostic value auditor for the parallel-memory experiment: scans the current repo for evidence that each memory engine — `actions/bkb.md`'s `kb/` and `actions/memory.md`'s `memory/` — is actually used and providing value, then renders a head-to-head verdict. Loaded lazily by `do-work-knowledge memory audit` (alias `memory value`); it has no routing row of its own.

**Strictly read-only.** This action never modifies, moves, or deletes anything — not even ledger appends. It reads files and git history, and it reports. Retiring an engine is a human decision made on this report, executed later as a normal maintenance pass.

`do-work-cli memory-audit --engine <bkb|memory|both>` is the sole executable authority for the probes and rubric classifications below. Consume its JSON findings and stop if the command is absent or fails; never reproduce the scans in prose. This action still owns the head-to-head interpretation, verdict, and recommendation because those are judgments rather than bookkeeping.

## When to Use

**Use when:**
- The user asks whether bkb (or the memory engine) is actually being used / worth keeping ("determine current value").
- The ADR-017 experiment window (~4 weeks of ledger data) is up and a verdict is due.

**Do NOT use when:**
- The user wants engine *content* health (broken links, orphan pages) → `do-work-knowledge bkb lint` / `do-work-knowledge memory status`.
- The user wants to consolidate or prune → `actions/dream.md` (manual, destructive).

## Input

`$ARGUMENTS` (optional): `bkb` or `memory` to audit one engine only; a `--kb <path>` flag passes through to KB location; empty → audit both + head-to-head. All commands run from `PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"`.

## Checks

Each check is an independent probe; run all that apply to the engines in scope. An engine whose store cannot be located short-circuits to **Absent** (run no further probes for it) — but for bkb, "cannot be located" means the full locating procedure below came up empty, never just "no `kb/` at the project root".

### bkb engine

**Locate the KB first** using `actions/bkb.md` → "Locating the Knowledge Base" — honor `--kb <path>`, then `kb/`, then `knowledge-base/`, then the parent-directory search, exactly as that section prescribes (it is the canonical procedure; do not re-derive it here). Call the result `<kb-root>`; only if that procedure finds nothing is the engine **Absent**. All probes below run against `<kb-root>` (shown as `kb/`-style paths for readability):

- **Existence & shape:** `find <kb-root>/wiki -name '*.md' 2>/dev/null | wc -l` wiki pages; `<kb-root>/raw/` inbox size.
- **Git activity:** `git log --oneline -- <kb-root>/ | wc -l` total commits; `git log -1 --format=%ci -- <kb-root>/` last-touched; `git log --format=%an -- <kb-root>/ | sort -u` distinct authors (human-touch signal). (A `<kb-root>` outside the project's git repo reports "no git history" — probe the ledger and log activity instead.)
- **Log activity:** entries in `<kb-root>/wiki/log.md` dated within the last 30 and 90 days (grep date headings).
- **Inbound references:** wiki pages cited from outside the KB — e.g. `grep -rl 'wiki/' --include='*.md' "$PROJECT_ROOT" | grep -v "^<kb-root>/"` and `[[wikilink]]` mentions outside `<kb-root>`. A wiki nobody links to is write-only.
- **Ledger stats:** `<kb-root>/usage-ledger.jsonl` per the shared procedure below.

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

The two ledgers also differ in durability, and the report must say so: `kb/usage-ledger.jsonl` is committed and travels with the repo, while `memory/usage-ledger.jsonl` is machine-local (git-excluded by `setup-memory`, so verbatim captures stay out of version control). A memory-engine count is therefore **this machine's** usage, not the team's — on a fresh clone or a second workstation it reads as zero from history that exists elsewhere. State the ledger's own age (`memory/usage-ledger.jsonl` mtime and first-line date) next to any memory-side count, and never read a thin memory ledger as disuse without checking whether the ledger itself is simply new here.

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
- Apply the fairness rule in every both-engine audit, and state it in the report — including the ledger-durability asymmetry (bkb's ledger is committed, memory's is machine-local).

## Red Flags

- The audit left `git status --porcelain` non-clean.
- A verdict that crowns a winner on `inject`/`capture`/`write` volume.
- bkb classified Stale/Absent purely from a missing ledger while `git log -- kb/` shows recent commits.

## Verification Checklist

- [ ] Report rendered with both engine sections (or the one requested), the head-to-head table, and a verdict paragraph.
- [ ] `git status --porcelain` unchanged by the audit.
- [ ] Every classification cites the probe evidence it rests on.
