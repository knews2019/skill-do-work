---
id: UR-098
title: 'Simplify the orchestrator: principles in prose, mechanics in the CLI'
created_at: 2026-09-02T14:37:54Z
requests: [REQ-503, REQ-504, REQ-505, REQ-506, REQ-507, REQ-508, REQ-509, REQ-510]
word_count: 95
---

# Simplify the Orchestrator: Principles in Prose, Mechanics in the CLI

## Summary

Move every mechanical step of the run action into the Go CLI behind one lifecycle command, `advance`, and rewrite `work.md` as principles keyed on conditions. Each move is one REQ that carries a four-part write set: the owning CLI command, the deleted prose, the deleted sentence-predicates in the contract suite, and a new behavior test. The analysis and the per-step classification are in `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70).

The chain runs after the recovery work: its root waits on REQ-489, REQ-498, REQ-499, REQ-500, REQ-501 and REQ-502.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-503 | Add the read-only `advance` lifecycle command |
| REQ-504 | Collapse Step 10, Crash Recovery and session-start prose into recovery |
| REQ-505 | Move selection and claim (Steps 1, 2.0, 2) behind `advance` |
| REQ-506 | Run the evidence gates (Steps 3.6, 5.75, 6.3, 6.5 gate) from `advance` |
| REQ-507 | Hand the archive and commit tails (Steps 8, 9) to `finalize` |
| REQ-508 | Reduce capture templates to minimal examples backed by the schema layer |
| REQ-509 | Merge the Common Rationalizations tables into one crew member |
| REQ-510 | Sweep `work-reference.md` sections whose contract is now a CLI behavior test |

## Batch Constraints

- One step per REQ, never a rewrite of `work.md`. Each REQ names all four parts of a move in its write set; a review refuses a move missing any part.
- Judgment stays prose: route criteria, planning, the Finding-Closure Ratchet, review depth, Fold-First minting, lessons, the one-broken-pipe principle.
- `advance` emits typed findings, never paragraphs; it must not become a second prose surface.
- The floor agent (reads files, runs shell) must still complete a run using only `advance` output plus the remaining prose.
- Serial chain: every REQ depends on the one before it, because all of them write `_dev/tests/contract-regressions.sh` and most write `work.md`.
- Run in one session; no second orchestrator during the migration.

## Full Verbatim Input

> ```
> and also since now we'll have rwr how can we simplify the skill?
> we should have more principles for the LLMs not exact steps, for that we have the go script.
> 
> I see the plan but I don't know what you want to do with it. I asked how it can be simplified. but now I don't know what you want todo , so how about making a do-work ai-report to show me how to simplify the current do-work ur/req orchestrator
> 
> can you capture-request with dependency setup so it runs after the others are done?
> ```

---
*Captured: 2026-09-02T14:37:54Z*
