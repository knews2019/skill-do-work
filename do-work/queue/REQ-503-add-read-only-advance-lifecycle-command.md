---
id: REQ-503
title: 'Add the read-only advance lifecycle command'
status: pending
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-489, REQ-498, REQ-499, REQ-500, REQ-501, REQ-502]
batch: orchestrator-simplification
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: [skills/do-work/tools/do-work-cli/internal/lifecycleadvance/, skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Add the Read-Only advance Lifecycle Command

## What
Add `do-work-cli advance REQ-NNN`: for the REQ's route, report the current lifecycle phase, the evidence still missing, and the exact next command, and refuse an impossible transition with a typed finding. Read-only in this REQ: it composes the existing commands (`next`, `claim`, `estimate-p50`, `preflight`, `scope-drift`, `qualify`, `run-blocked-check`, `complete`, `finalize`, `recover-finalization`) into one state machine but mutates nothing itself.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Every later move in this chain deletes prose by pointing at `advance`. Landing it read-only first means it changes no behavior, so it needs no prose deletion and every other REQ can depend on it.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- Phases per route (A, B, C) derived from the current `work.md` step order; the report's step table is the source, with mechanical phases CLI-driven and judgment phases reported as "agent judgment: <what>".
- Output is typed: phase, missing evidence with the file or field that would satisfy it, `next_argv`, `verification_argv`; text and JSON formats like every other command.
- Refusals are typed findings (for example `ADVANCE-EVIDENCE-MISSING`, `ADVANCE-PHASE-UNKNOWN`); no paragraphs.
- No mutation in this REQ; `advance` never writes a REQ, the checkpoint, or Git.
- Prime file `prime-do-work-cli.md` documents the command and its phase table.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Root of the chain: waits on REQ-489, REQ-498, REQ-499, REQ-500, REQ-501 and REQ-502 so recovery and finalization exist before the state machine composes them.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Run `do-work-cli --format json advance REQ-NNN` against a fixture REQ in each of: queue, claimed without triage, planned Route C, implemented without qualification, archived without provenance.
**Why RED now:** The command does not exist.
**GREEN when:** Each fixture returns the expected phase, the missing evidence, and a `next_argv` that names the existing command for that phase; an unknown state returns a typed refusal; `go test ./internal/lifecycleadvance` passes and no file outside the write set changes.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2299 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on families `cross-action-exception-closure` and `opaque-evidence-projection`.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*
