---
id: REQ-270
title: Carry a worktree builder's hand-back sections into Step 8
status: pending-answers
created_at: 2026-08-18T21:45:21Z
status_changed_at: 2026-08-18T21:45:21Z
user_request: UR-055
addendum_to: REQ-259
domain: general
review_generated: true
sweep: true
sweep_key: worktree-handback-sections-unread-by-step-8
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
---

# Carry a Worktree Builder's Hand-Back Sections into Step 8

## What

`skills/do-work/actions/work.md` Step 8 substep 4 reads `## Discovered Tasks` **from the REQ file**, describing it as "appended by the implementation agent as a separate section". In worktree dispatch mode the implementation agent **cannot write the REQ file** — the REQ lives in the main tree, which is exactly what `State stays home` forbids a builder to touch — so the section is never there, Step 8 finds nothing, and every out-of-scope find the builder recorded is silently dropped.

This is not hypothetical and it is not a builder error: it fired on REQ-259, whose builder correctly routed three Discovered Tasks to its hand-back per the dispatch brief. The section was absent from the REQ until the orchestrator noticed and transcribed it by hand. Nothing in the pipeline would have reported the loss.

**State the condition, not the instance.** `## Discovered Tasks` is the one that fired; the defect is the general shape — *any* Step 8 substep that expects to read a builder-authored section from the REQ file is silently disarmed under fan-out. Sweep the substeps for that shape rather than patching the one known case.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-259's independent review, Important finding 2 (gate: rule-change). Created `pending-answers` per the generation-≥2 cascade depth stop, since REQ-259 is itself `review_generated: true`; the reviewer assessed it as Important rather than critical-grade, so it does not pierce to `pending` on its own.

Worth knowing while it sits in the queue: **the loss is silent and unbounded in the past.** Every fan-out run before this fix depended on the orchestrator happening to transcribe by hand. This session's wave-1 REQs were transcribed manually once the review surfaced it, so nothing is lost here — but that is a person catching it, not the pipeline.

## Requirements

- Step 8's builder-authored-section reads work in worktree dispatch mode: the section is taken from the hand-back when the REQ file does not carry it.
- The rule is keyed on the **condition** (a Step 8 substep reading a section the builder authors) and marks any list of affected substeps illustrative, so a substep added later inherits the behavior instead of waiting to be remembered.
- The failure mode is loud rather than silent: a run that finds neither the REQ section nor a readable hand-back says so.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [ ] REQ-259's review found that a worktree builder's `## Discovered Tasks` never reach Step 8, because Step 8 reads that section from the REQ file and a worktree builder may not write it — so every out-of-scope find is silently dropped unless the orchestrator transcribes it by hand. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — rely on the orchestrator transcribing hand-back sections as part of integration, and say so in the dispatch contract instead of changing Step 8.
