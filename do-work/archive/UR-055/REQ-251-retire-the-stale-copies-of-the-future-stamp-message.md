---
id: REQ-251
title: Retire the stale copies of the future-stamp message
status: completed
created_at: 2026-08-18T13:55:32Z
claimed_at: 2026-08-18T18:25:40Z
completed_at: 2026-08-18T18:41:26Z
commit: 96bb593
kb_status: promoted
kb_entry: REQ-251-retire-the-stale-copies-of-the-future-st.md
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

---

## Triage note

Route A confirmed in practice — two files, text only, no behaviour change.

## Implementation Summary

**What was done:** Retired the two surviving copies of the pre-REQ-245 future-stamp wording. The `verify_test.go` fixture literals (input-only; assertions unchanged) now derive their reason from production via the existing `reversedSpanAnomalyReason(t)` helper, so a future message change cannot strand another copy; `timestamp.go`'s `formatCanonicalTimestamp` comment now matches its test twin's claim ("one of the two corruptions…") instead of overstating ("the exact corruption…").

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified) — two fixture literals replaced with `reversedSpanAnomalyReason(t)` calls
- `skills/do-work-board/tools/queue-kanban/timestamp.go` (modified) — doc comment aligned with its test twin

*Integrated by orchestrator from builder hand-back; merge range `37ac5ce..96bb593`.*

## Decisions

- **D-01 (DECIDE & STATE, builder):** fixtures derive the reason from `reversedSpanAnomalyReason(t)` instead of a freshly pasted literal — the REQ exists because hand-typed copies strand when the message moves (twice already); a helper call ties the fixtures to production permanently. Orchestrator concurs: this closes the drift *class*, where the literal the brief imagined would have closed the instance.
- **D-02 (DECIDE & STATE, builder):** the source comment mirrors the test twin's wording rather than inventing new phrasing, so the two comments state one claim in one voice.

## Qualification

Passed — 2 files verified in merge range `37ac5ce..96bb593`; both requirements traced; P-A-U audited (text-only diff, no debug artifacts). Orchestrator independently re-ran both retired-wording greps over `skills/` (zero hits) and the affected package tests (`ok`, 3.4s).

## Review

**Overall: 95%** | 2026-08-18T18:40:30Z (Route A quick scan, orchestrated inline)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

- [x] No behaviour change, no new checks — text only: verified by diff (comment + two fixture-literal replacements; assertions byte-identical).
- [x] Post-change grep for the retired wording returns nothing in shipped source — reproduced independently by the orchestrator over `skills/` (both retired phrasings, zero hits; survivors are only frozen `do-work/` history and this REQ's own quotation, as required).
- [x] `reversedSpanAnomalyReason(t)` exists in the same package, is already exercised by `timestamp_test.go`'s pairing tests, and produces the identical timestamps the fixtures used — verified by reading and by re-running `TestVerify*` / `TestFormatCanonicalTimestamp*` (ok, exit 0).
- Full gate: exit 0 observed by the builder on its final commit; the orchestrator's post-merge gate run covers the merged tree below.

**Findings:** none Important. One Nit: the REQ's own queue file quotes the phrase it retires (by design — the requirement text needs it), so future greps hit the archived REQ; acceptable and documented here.

**Acceptance: Pass.** No follow-ups.

*Reviewed by review-work action (orchestrated inline, Route A depth; merge range `37ac5ce..96bb593`)*

## Lessons Learned

**What worked:** Deriving fixture text from production (helper call) instead of re-pasting the current wording — the third copy would have stranded at the next message move exactly as the first two did.

**Worth knowing:** `reversedSpanAnomalyReason(t)` in `timestamp_test.go` is the canonical way for any test in the package to obtain the production reversed-span reason; new fixtures should call it, never paste.

## Orientation

Now a grep for the future-stamp diagnosis finds only the production wording — the two stale copies are gone and the fixtures track production mechanically. Lives in the board tool's timestamp-diagnostics area. Leaf change; map unchanged.
