# REQ-505 Builder Brief

Implement the validated Route C plan in `do-work/working/REQ-505-move-selection-and-claim-behind-advance.md` from an isolated builder worktree. Do not edit any `do-work/` path. Read `CLAUDE.md`, the request, all listed prime files, `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `maintenance.md`, and the action/shell/CLI lesson satellites before writing.

## Boundary

- Add a queue-mode `advance` composition for fresh default, explicit, UR, wave, fan-out, and stateless targeted replay.
- Preserve the existing read-only working/archive classifier and checkpoint-only mutation mode.
- Carry frozen target membership, provenance, consumed state, non-fan-out flags, and dispatch bound as typed data and tokenized continuation argv. Observe canonically, project frozen membership, then bound.
- Use the recursive repository snapshot for nested archive collisions. Do not reuse or widen the shallow legacy archive-collision helper.
- Add exact guarded request-state holds for archive collisions and dependency cycles, canonical unblock from selector-supplied successful probe evidence, and a fresh per-request claim plan/commit after each mutation.
- Report partial success truthfully if a later claim refuses; never imply prior committed claims rolled back.
- Collapse procedural selection/claim prose and align every declared reader, while leaving evidence-gate execution to REQ-506 and finalization to REQ-507.

## TDD

Write public RED fixtures first. Cover default, explicit override, UR chain/fork, wave/fan-out, REQ-453 projection-before-bounding with a later member, nested archive collision, dependency cycle, blocked-probe success and failure, stale/dirty partial refusal, hostile tokens, exact committed footprints, and unchanged working/archive/checkpoint contracts. The captured sentence-predicate RED is stale after REQ-504; document the adaptation rather than manufacturing prose failures.

## Exact Scope

The request frontmatter and `## Scope` list the complete 22-path ceiling. Touch only paths actually needed inside that ceiling; report each created, modified, or deleted path exactly. If a required file lies outside it, stop and report before editing.

## Handback

Commit the coherent implementation on your operative branch, leave the worktree clean, and write `do-work/runs/work-2026-09-04-164003/REQ-505-handback.md` in the owner checkout. Include RED/GREEN evidence with timings, exact manifest, P-A-U, decisions, lessons read, discovered tasks, integration seams, and the commit hash. Never commit the handback or any `do-work/` path on the builder branch.
