---
id: REQ-215
title: Sync completion-anomaly prose with the reversed-span class
status: pending
created_at: 2026-08-17T08:25:24Z
user_request: UR-048
addendum_to: REQ-213
review_generated: true
sweep: true
sweep_key: anomaly-class-prose-predates-reversed-span
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
depends_on: []
maintenance: false
effort_estimate: normal
write_set: [skills/do-work-board/docs/board-guide.md, skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/web/board-cards.js]
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-08-17T08:25:24Z
  basis:
    - Route A
    - 4-file write set
---

# Sync Completion-Anomaly Prose with the Reversed-Span Class

## What

Shipped prose still defines a completion anomaly as "completion instant can't be resolved" — stale since REQ-213 added a class whose instant resolves fine (the reversed span). One root cause, several instances; sweep them.

**Finding provenance (REQ-213 review, Important 2, gate: user-visible, sweep-consolidated):** root cause — "all this prose predates a class whose completion instant resolves fine."

## Instances

- [ ] `skills/do-work-board/docs/board-guide.md:23` — "finished REQs whose completion instant can't be resolved" → broaden to cover broken completion bookkeeping including reversed spans
- [ ] `skills/do-work-board/docs/board-guide.md:35` — chip legend "unresolvable completion instant, or a timestamp later than now" → same broadening
- [ ] `skills/do-work-board/tools/queue-kanban/model.go` (never-silent warning suffix) — generic "fix: stamp completed_at: with a UTC ISO instant and/or a commit: field…" self-contradicts for a reversed span whose completed_at IS valid; make the suffix class-appropriate or defer to the reason's own fix text
- [ ] `skills/do-work-board/tools/queue-kanban/web/board-cards.js:401-406` — comment "carry no honest completion instant… never sorted into Recently done" → correct for classes with a resolved instant

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- Each instance reads true for all four anomaly classes (unparseable completed_at, undatable commit, neither field, reversed span).
- Comment/doc edits only, except the model.go warning-suffix string (behavioral text, covered by existing tests' fragment matching — keep fragments the tests pin, or update tests in the same commit).
- Go module tests stay green.

## Red-Green Proof
**RED prompt/case:** `board-guide.md:23` tells a user with a reversed-span anomaly that its "completion instant can't be resolved" — a false diagnosis; the never-silent line simultaneously says the completed_at is valid and tells them to re-stamp it.
**Why RED now:** Prose predates the class.
**GREEN when:** Every instance above describes the anomaly family accurately for all four classes; module tests green.
**Validation:** Review-generated (REQ-213 Important 2).

## Full Context
See `do-work/user-requests/UR-048/input.md`.
