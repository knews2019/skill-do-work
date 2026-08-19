---
id: REQ-266
title: Name builds beside the JS renderer's measured face numbers
status: pending
created_at: 2026-08-18T20:07:08Z
status_changed_at: 2026-08-18T21:01:24Z
user_request: UR-051
addendum_to: REQ-252
domain: general
review_generated: true
sweep: true
sweep_key: durations-measured-face-constants-lack-provenance
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
---

# Name Builds Beside the JS Renderer's Measured Face Numbers

## What

`web/board-durations.js` presents measured face numbers (12.83 / 10.43 / 2.41 at `DURATIONS_LABEL_ROW_HEIGHT` and `DURATIONS_LABEL_TEXT_ASCENT`) as current fact with no build named — the same provenance gap REQ-252 closed in the Go files, on the JS surface its go/parser test cannot reach. Extend the rule: every browser-measured number in the JS comments names its build, and the mechanism that keeps it true (a JS-side check, or a stated review convention) is the builder's call, recorded either way.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-252's declared Discovered Task, verified by its review (F1b, gate: rule-change). Carries the sweep key because it is the same root cause as REQ-252 — that REQ is the claimed-and-archived sweep for the key, so this is a new file per the append rule, not an append. Created `pending-answers` per the generation-≥2 depth stop.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-252: the JS renderer's measured numbers carry no build provenance. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the Go-side comments now carry the builds and the JS numbers are near-duplicates of them.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.
