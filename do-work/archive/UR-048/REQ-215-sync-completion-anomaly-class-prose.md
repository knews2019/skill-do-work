---
id: REQ-215
title: Sync completion-anomaly prose with the reversed-span class
status: completed
created_at: 2026-08-17T08:25:24Z
claimed_at: 2026-08-17T08:30:14Z
route: A
completed_at: 2026-08-17T08:32:42Z
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
write_set: [skills/do-work-board/docs/board-guide.md, skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go]
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

- [x] `skills/do-work-board/docs/board-guide.md:23` — "finished REQs whose completion instant can't be resolved" → broaden to cover broken completion bookkeeping including reversed spans
- [x] `skills/do-work-board/docs/board-guide.md:35` — chip legend "unresolvable completion instant, or a timestamp later than now" → same broadening
- [x] `skills/do-work-board/tools/queue-kanban/model.go` (never-silent warning suffix) — generic "fix: stamp completed_at: with a UTC ISO instant and/or a commit: field…" self-contradicts for a reversed span whose completed_at IS valid; make the suffix class-appropriate or defer to the reason's own fix text
- [x] `skills/do-work-board/tools/queue-kanban/web/board-cards.js:401-406` — comment "carry no honest completion instant… never sorted into Recently done" → correct for classes with a resolved instant

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** prime-kanban-board + crew loaded. Four anchored edits (guide strip definition + chip legend broadened to "broken completion bookkeeping incl. reversed span"; never-silent suffix defers to the per-class reason; web comment corrected incl. the dual-listing reality). One test consequence anticipated: any pin of the old suffix updates in the same commit.
- [x] **[APPLY]:** As planned; the anticipated test pin surfaced (completion_anomaly_test.go asserted "completed_at" in every warning — the exact self-contradiction being removed) and was updated with a REQ-215 note per the cross-REQ test-break rule; write_set extended accordingly.
- [x] **[UNIFY]:** Diff reviewed — board-guide.md (two lines), model.go (one string), board-cards.js (one comment), completion_anomaly_test.go (one assertion + note). go vet + module tests green; full verify FAIL set identical to baseline.

## Requirements

- Each instance reads true for all four anomaly classes (unparseable completed_at, undatable commit, neither field, reversed span).
- Comment/doc edits only, except the model.go warning-suffix string (behavioral text, covered by existing tests' fragment matching — keep fragments the tests pin, or update tests in the same commit).
- Go module tests stay green.

## Red-Green Proof
**RED prompt/case:** `board-guide.md:23` tells a user with a reversed-span anomaly that its "completion instant can't be resolved" — a false diagnosis; the never-silent line simultaneously says the completed_at is valid and tells them to re-stamp it.
**Why RED now:** Prose predates the class.
**GREEN when:** Every instance above describes the anomaly family accurately for all four classes; module tests green.
**Validation:** Review-generated (REQ-213 Important 2).

## Implementation Summary

**Files changed:**
- `skills/do-work-board/docs/board-guide.md` (modified) — strip definition and chip legend now cover "broken completion bookkeeping" including the reversed span
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — never-silent warning suffix defers to the per-class reason ("the reason above names the broken frontmatter field(s)…") instead of restating a completed_at fix that self-contradicted for reversed spans and commit-hash anomalies
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified) — anomalies-strip comment corrected: some classes carry a resolvable instant, and an in-window anomalous ticket may also appear under Recently done
- `skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go` (modified) — warning assertion updated from the removed "completed_at" fragment to the new routing suffix, with a REQ-215 note (intentional behavior change, cross-REQ test-break rule)

**What was done:** Every restatement of the anomaly family now reads true for all four classes.

## Decisions
<!-- D-XX counter: none used in Open Questions. -->
- **D-01** (DECIDE & STATE): The suffix defers to the reason rather than enumerating per-class fixes — same rationale as REQ-214's remedy (one authoritative fix text per class, stated where the class is detected).

## Qualification

Passed — qualify.sh exit 0; Route A (no Scope section, drift skip per contract). All four instances checked off with the edits traced.

## Testing

**Tests run:** go vet, module tests (-count=1), full maintainer-verify
**Result:** ✓ Green after the intentional pin update; FAIL set identical to the environment baseline.

**Red-green validation:** RED: TestCompletionAnomaliesFlaggedInBoardModel failed against the new suffix precisely because it pinned the removed self-contradiction — demonstrating the old text was load-bearing in the wrong direction. GREEN: assertion updated to the routing suffix; suite green; the live warning for REQ-091 no longer instructs re-stamping a valid completed_at.

*Verified by work action*

## Review

**Overall: 93%** | quick scan (Route A sweep, run inline by the orchestrator)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor (report only): the test-pin update landed in the same commit as the string change it validates — correct per the cross-REQ rule, but a reviewer diffing test-only changes would want the REQ-215 note, which is present. Acceptance: Pass — all four instances read true for all four anomaly classes.

*Reviewed by review-work action (quick scan)*

## Lessons Learned

**What worked:** Anticipating the test-pin consequence in [PLAN] before touching the string — the failure arrived exactly where predicted.
**What didn't:** Ran maintainer-verify from the Go module cwd once (exit 127, empty log) — relative-path invocations of the gate must run from the repo root.
**Worth knowing:** The warning's fix text now lives once per class in detectCompletionAnomaly's reasons; any new anomaly class must carry its own fix text in its reason.

## Orientation

**Now you can** read any completion-anomaly surface — guide, chip legend, warning line, web strip — and get a description that is true for all four classes, with the fix stated by the class that detected it. Leaf prose/comment sync; no map delta.

## Full Context
See `do-work/user-requests/UR-048/input.md`.
