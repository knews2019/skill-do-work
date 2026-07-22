---
title: "ADR-017: Run a Parallel Memory Engine Experiment With Usage Ledgers"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-07-22
sources:
  - actions/memory.md (new action — dispatcher)
  - actions/memory-reference.md (new companion — schemas, ledger contract)
  - actions/memory-value.md (new action — engine-agnostic value auditor)
  - hooks/memory-session-start.sh, hooks/memory-stop-capture.sh, hooks/memory-hooks.json (new optional hooks)
  - actions/install.md (5th target — memory-module)
  - actions/bkb.md (query/ingest ledger instrumentation)
related:
  - page: adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb
    rel: complements
  - page: adr-014-considered-declined-autonomous-loop-until-done
    rel: complements
created: 2026-07-22
updated: 2026-07-22
confidence: medium
---

# ADR-017: Run a Parallel Memory Engine Experiment With Usage Ledgers

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] (complements)

## Context

The skill already ships a memory system: `actions/bkb.md` compiles raw sources into a curated `kb/wiki/` — retrieval memory that optimizes *synthesis* (ADR-010). An external "Claude Code Memory Plan" proposes the opposite trade: optimize *capture* — a small standing working-memory file injected at session start, dated daily logs auto-appended by a Stop hook, and layered recall over the raw history. Persist-everything capture is both its selling point (zero-effort, nothing lost) and its known weakness (signal decays as the store grows, contradictions are never reconciled).

Rather than argue the trade on theory, the user wants to **battle-test both engines in parallel and let usage data decide**. That requires (a) the new engine to exist alongside bkb without either replacing the other, (b) both engines to leave comparable usage evidence, and (c) a neutral auditor that reads the evidence and renders a verdict — including the user's specific ask for a "bkb-determine-current-value" scan of whether the existing KB is actually used.

Constraints from the skill's own rules: every new action must justify not being a sibling skill (CLAUDE.md); SKILL.md's routing table has an enforced word budget with ~178 words of headroom; no new compiled tools (`tools/queue-kanban` is the sole toolchain exception); `actions/dream.md` is destructive and must never run from a hook; action prose must stay platform-agnostic while hooks are a Claude Code-specific enhancement.

## Decision

**Ship the experimental memory engine inside do-work (not as a sibling skill), instrument both engines with append-only JSONL usage ledgers, and judge the experiment with an engine-agnostic auditor after ~4 weeks of ledger data.** Concretely:

- **In-skill, one routing row.** `actions/memory.md` is a sub-command dispatcher (`remember | recall | status | bootstrap | audit`) with a companion `actions/memory-reference.md`; the auditor `actions/memory-value.md` is lazy-loaded by `memory audit` and takes no routing surface. In-skill because the engine reuses do-work's install/hook machinery (`do-work install memory-module` becomes the 5th install target), because the auditor must compare both engines from one vantage point, and because the experiment's whole premise is a fair fight inside one skill — a sibling skill would fragment the instrumentation and skew the comparison. This paragraph is the CLAUDE.md-required justification for all three files.
- **Usage ledgers on both engines.** One JSON line per event (`inject`, `capture`, `write`, `recall`, `hit_cited` for the memory engine; `query`, `ingest`, `hit_cited` for bkb) appended to `memory/usage-ledger.jsonl` / `kb/usage-ledger.jsonl`. Canonical schema lives in `actions/memory-reference.md`; appends are best-effort and never block the action they instrument. The verdict signal is **hit-cited rate** (recalled content actually used in an answer), not raw write volume — a store nobody reads from is not memory, it's a landfill.
- **Stop-hook capture is verbatim-tail + hash, not an LLM summary.** A shell Stop hook cannot summarize, and the alternative — `{"decision":"block"}` to force the agent to write a summary — was rejected: it delays every session end and risks stop-loops. The hook appends the final exchange (truncated, third-person-framed, sha256-prefix deduplicated) to today's log and always exits 0. Curation is `memory remember`'s job plus cap-driven consolidation of the 2,500-char working-memory file. Tradeoff accepted: captures are raw, not distilled.
- **Recall is layered, floor-first.** Lexical grep scoring over working memory + logs always works (design-for-the-floor); a semantic layer activates only when an embedding backend is detected (ollama, an `embed` CLI, or an API key), degrading silently like `actions/board.md` does without Go. No vector database ships with the skill.
- **dream stays manual.** Nothing hook-side consolidates or prunes; the only hook write is the append-only capture. The engine's own consolidation runs at `remember` time, driven by the cap, on `memory/working-memory.md` only.
- **Exit criterion.** After roughly 4 weeks of ledger data, `do-work memory audit` renders the head-to-head. The auditor classifies each engine Active/Idle/Stale/Absent, weights bkb's pre-instrumentation era via `git log -- kb/` and `kb/wiki/log.md` history (absence of ledger ≠ absence of value), and recommends keep-both / retire-one / fix-instrumentation. Retirement is a human decision — the auditor never deletes anything.

## Alternatives

1. **Judge on theory, don't build.** Rejected — the capture-vs-synthesis trade is genuinely empirical; which engine wins depends on how this user actually works.
2. **Sibling skill vendored via install (like last30days).** Rejected for the experiment phase — split instrumentation, two update channels, and the auditor would live in one skill judging another. Revisit if the new engine wins and grows heavy.
3. **Full-fidelity vector memory (require an embedding backend).** Rejected — breaks design-for-the-floor and would make the new engine unavailable in exactly the environments bkb runs in, rigging the usage comparison.
4. **Replace bkb outright.** Rejected — bkb's compile step (contradiction resolution, dedup into a wiki) is the thing persist-everything capture cannot do; the engines optimize different ends of the pipeline and may both earn their keep.

## Consequences

do-work temporarily carries two memory engines — accepted bloat, bounded by the exit criterion and by the anti-bloat gate this ADR satisfies. SKILL.md spends one routing row (~80 words of its remaining budget). bkb's `query`/`ingest` steps gain best-effort ledger appends. `hooks/` grows two scripts and a fragment; the install target must compose hook entries into a consumer's `.claude/settings.json` without clobbering existing ones (append-only jq merge with backup — a new class of install risk, mitigated by parse-verify and entry-count checks). When the experiment concludes, the loser's removal is a normal maintenance pass (`crew-members/maintenance.md` applies) and gets its own changelog entry.

## References

- [actions/memory.md](../../actions/memory.md) — the dispatcher
- [actions/memory-reference.md](../../actions/memory-reference.md) — schemas, ledger contract, recall shell
- [actions/memory-value.md](../../actions/memory-value.md) — the auditor
- [actions/bkb.md](../../actions/bkb.md) — the incumbent engine + its instrumentation
- [actions/install.md](../../actions/install.md) — the memory-module install target
- [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] — why bkb is synthesis-first memory
