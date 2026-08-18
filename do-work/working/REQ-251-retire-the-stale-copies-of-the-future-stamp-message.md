---
id: REQ-251
title: Retire the stale copies of the future-stamp message
status: claimed
created_at: 2026-08-18T13:55:32Z
claimed_at: 2026-08-18T18:25:40Z
route: A
status_changed_at: 2026-08-18T13:55:32Z
user_request: UR-055
addendum_to: REQ-245
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/verify_test.go
- skills/do-work-board/tools/queue-kanban/timestamp.go
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T18:26:21Z
  basis:
    - trivial short-circuit
---

# Retire the Stale Copies of the Future-Stamp Message

## What

REQ-245 rewrote the future-stamp diagnosis in five renderers. Two copies of the old wording survive outside its write set, both harmless to behaviour and both misleading to the next person who greps for the message.

## Instances

- [ ] **`verify_test.go:1186` and `:1230` hold hand-typed copies of the retired reversed-span reason as fixture literals.** They are input, never asserted against `detectCompletionAnomaly` — REQ-245's review read the surrounding asserts and confirmed they check `"REQ-9330"` and `"is earlier than claimed_at"` only. So nothing is broken. But the message they copy has now moved twice, and the next person greping for it will find two copies of a sentence that no longer exists anywhere.
- [ ] **`timestamp.go:35` — "the exact corruption the Timestamp rule warns about" — is now strictly stronger than its own test's claim.** It is largely defensible, since `formatCanonicalTimestamp` can only prevent the timezone cause; REQ-245's twin comment at `timestamp_test.go:43-46` already says so explicitly ("one of the two corruptions… This is the corruption a correct writer rules out"). Bring the source comment in line with the test that describes it.

## Requirements

- No behaviour change and no new checks. This is text.
- After the change, a grep for the retired wording across shipped source returns nothing but `do-work/` history, `CHANGELOG.md` release notes and `ai-reports/` narrative — all correctly frozen.

## Context

Both were found by REQ-245's builder and its independent reviewer, and both were deliberately left out of REQ-245 rather than widening that REQ a third time.

---

## Triage

**Route: A** - Simple

**Reasoning:** Two named stale text sites with line numbers and the replacement direction stated; no behaviour change, no new checks.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
