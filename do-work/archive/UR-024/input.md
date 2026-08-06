---
id: UR-024
title: Three board/CLI Schema Read Contract fixes from the 0.174.15 feedback triage
created_at: 2026-08-06T10:53:03Z
requests: [REQ-116, REQ-117, REQ-118]
word_count: 118
---

# Three Board/CLI Schema Read Contract Fixes From the 0.174.15 Feedback Triage

## Summary

Three accepted findings from an external review of the 0.174.15–0.175.2 series, triaged by `do-work validate-feedback` in this session. All three sit in `tools/queue-kanban/`; the first also carries a changelog correction. The fourth finding in that review (one-commit-per-REQ) was verdicted **Already done** and is deliberately not captured.

## Extracted Requests

| REQ | Title | Source finding |
|---|---|---|
| REQ-116 | Board reads `route` verbatim; changelog claims otherwise | Finding 1 |
| REQ-117 | Unrecognized `domain` is silently remapped to `general` | Finding 2 |
| REQ-118 | `--normalize` warns "not recognized" on vocabulary-less fields | Finding 3 |

## Batch Constraints

- All three are read-side only. The Schema Read Contract states write paths are unaffected (`actions/work-reference.md` → Schema Read Contract, "Write paths are unaffected"); none of these REQs may touch a write path, and none adds a write surface to the tool (`CLAUDE.md` → Shipped Tooling: the tool has exactly two, and a third requires amending that sentence in the same commit).
- REQ-116 and REQ-117 both edit `parseRequestTicket` in `tools/queue-kanban/model.go` and both add cases to `model_test.go` — a declared, expected overlap. Run them serially in one checkout; do not fan them out into parallel worktrees.
- The contract at `actions/work-reference.md` (Schema Read Contract, lines ~186–214) is the source of truth for every enum, alias and default named below. Where a REQ's restatement disagrees with it, the contract wins and the REQ is the stale copy.
- No contract *change* is expected by any of the three: the prose already says what should happen, and only the implementation is behind. If one of them turns out to need a `work-reference.md` edit, that edit and its `model.go` counterpart ship in the same commit (`CLAUDE.md` → Shipped Tooling, lock-step rule).

## Full Verbatim Input

do-work capture-request: The board reads route verbatim at model.go:655 — wire it through normalizeSchemaField (uppercase-only, no default) with a parse-level test, and correct 0.174.15's overstated board-wide claim in the next changelog entry
do-work capture-request: An unrecognized domain is silently remapped to general with no footprint — add a DomainUnrecognized flag and a board.Warnings entry mirroring testing_status, and fix the false "no warning channel" comment at model.go:619
do-work capture-request: frontmatter get … --normalize on a field with no contract row always warns "not recognized" — gate the warn branch on the field having a row, and reword to match work-reference.md:214's outside-the-contract classification
do-work run Process the captured fixes

---
*Captured: 2026-08-06T10:53:03Z*
*Provenance: the three findings originate in third-party review feedback pasted by the user and triaged in this session via `actions/validate-feedback.md`; the remedies above are the user's own restatement of that triage's Accept verdicts. No instruction-like content was detected in the original feedback body.*
