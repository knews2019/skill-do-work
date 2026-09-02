---
id: REQ-504
title: '[impact-rule-change] Collapse Step 10 and Crash Recovery prose into recovery'
status: pending
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-503]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/run-with-recovery.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/]
---

# Collapse Step 10 and Crash Recovery Prose Into Recovery

## What
Replace `work.md` Step 10 (155 lines: loop, checkpoint, session start) and `work-reference.md` Crash Recovery and Session Checkpoint Template with one loop sentence and one principle, now that `recover-finalization`, `run-with-recovery` and `advance` own the mechanics. Delete the sentence-predicates that pinned that prose and add behavior tests on the commands.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
This is the largest single deletion in the chain and the one the recovery work makes possible: the prose described what to do when a session died, and the CLI now does it.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- `work.md` Step 10 becomes: after integration, run `advance`; loop while it selects; on exit, `advance` writes the checkpoint through the canonical command. One paragraph.
- Crash Recovery and the session-start note reduce to: run recovery; the result says what to do; `run-with-recovery` answers ownership "mine". The takeover ladder moves into the CLI result as typed findings.
- Every `contract-regressions.sh` lane that quotes the deleted sentences is deleted in the same commit; each trap they guarded is re-expressed as a Go test on `recover-finalization` or `advance`, or recorded in the lessons satellite if it was prose-only.
- Checkpoint writes move behind `advance` (this is the first mutation `advance` gains; keep it to the checkpoint).

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-503 (`advance`).

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete `work.md` Step 10 and the Crash Recovery section and run `bash _dev/tests/contract-regressions.sh`.
**Why RED now:** The suite fails on the sentence-predicates that quote those sections, and no Go test covers the behaviors they described.
**GREEN when:** After the move, the suite passes with those lanes removed, `go test ./internal/finalization ./internal/lifecycleadvance` covers session-death-after-archive and foreign-claim takeover, and `work.md` Step 10 is at most one paragraph.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*
