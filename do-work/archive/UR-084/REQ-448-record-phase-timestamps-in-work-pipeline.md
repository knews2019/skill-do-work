---
id: REQ-448
title: 'Record per-phase timestamps through the work pipeline'
status: completed
route: C
created_at: 2026-08-31T20:38:40Z
user_request: UR-084
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-09-01T22:55:42Z
  basis:
    - Route C
    - 8-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - persistence changes
    - cross-route regression gates
    - full-suite verification
related: [REQ-449]
claimed_at: 2026-09-01T22:52:02Z
planning_at: 2026-09-01T22:56:07Z
dispatch_at: 2026-09-01T22:56:07Z
builder_handback_at: 2026-09-01T23:07:16Z
integration_at: 2026-09-01T23:07:16Z
review_at: 2026-09-01T23:10:36Z
release_at: 2026-09-01T23:18:12Z
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
completed_at: 2026-09-01T23:16:59Z
commit: 01a16769
---

# Record Per-Phase Timestamps Through the Work Pipeline

## What

Record optional per-phase timestamps on a REQ as it moves through the work pipeline — planning, dispatch, builder handback, integration, review, remediation, re-review, release — and display the phase breakdown at completion, so total wall time is not mislabeled as implementation time. `claimed_at` → `completed_at` stays the calibration span.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Use eight optional flat milestone fields (`planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at`, `re_review_at`, `release_at`). Stamp only successful observed events, strip all on recovery/blocked release, preserve claim-to-completion calibration, and render an ordered optional milestone breakdown without fabricating skipped phases.
- [x] **[APPLY]:** Added eight optional milestone observations across the canonical lifecycle schema, board parser, derived display payload, and shipped card/detail UI; retained claim-to-completion calibration unchanged.
- [x] **[UNIFY]:** Reviewed all ten declared files in scoped diff context. Verified lifecycle write/reset points, parser fields and future-stamp diagnostics, phase ordering and omission semantics, payload/UI wording, full-route and historical regressions, JavaScript syntax, and absence of debug artifacts.

## Detailed Requirements

- Record timestamps for planning, dispatch, builder handback, integration, review, remediation, re-review, and release. Phases a route skips (Route A/B never plan) simply get no stamp.
- "retain claimed_at → completed_at for calibration" — the calibration log (`skills/do-work/actions/work.md` Step 8 substep 7.5, `skills/do-work/actions/estimate-reference.md` § Calibration) keeps recording the single wall span; phase stamps are display/diagnosis, not a new calibration input.
- "display the phase breakdown so total wall time is not mislabeled as implementation time" — the completion report shows where the wall span went.
- "Keep historical REQs compatible when phase timestamps are absent" — fields are optional and additive; a REQ without them is fully valid and every existing reader behaves unchanged (precedent: the `estimate:` block's backwards-compatible contract, `skills/do-work/actions/work-reference.md` § Request File Schema).

## Constraints

- New fields must state their normalize/absent behavior in the Schema Read Contract (`skills/do-work/actions/work-reference.md` § Schema Read Contract) — absence semantics must not be left undefined.
- Crash recovery currently strips `claimed_at`/`route` (`skills/do-work/actions/work-reference.md` § crash recovery); the builder must make an explicit keep-or-strip decision for phase stamps so they never survive a crash as stale data by accident.
- If the board displays the breakdown, the Go parser changes in lock-step, same commit (`skills/do-work-board/tools/queue-kanban/model.go`, `durations.go`; parser lock-step rule in `skills/do-work/actions/work-reference.md`).
- Timestamps follow the existing Timestamp rule: current UTC instant, never local time with a Z suffix.
- Per the REQ-228 precedent (persisting derived forecasts disallowed), phase stamps are observations, not derivations — they sit on the allowed side of that line.

## Builder Guidance

Certainty: Firm on intent (phase visibility without disturbing calibration), Exploratory on representation — the builder chooses the field shape (e.g. a nested `phases:` block vs flat `*_at` keys) and which pipeline steps stamp which phase. Keep it simple; not every run needs every stamp.

## Red-Green Proof
**RED prompt/case:** A completed Route C REQ (e.g. REQ-411 in `do-work/runs/work-2026-08-31-165510/`) carries only `claimed_at` and `completed_at`; its wall span includes review, remediation, and re-review waits, and nothing in the REQ or the completion report distinguishes implementation time from pipeline time.
**Why RED now:** Only six `*_at` fields exist (`work-reference.md` § Request File Schema); no per-phase stamp is defined anywhere, and the calibration log's `wall_minutes` is the only duration ever computed.
**GREEN when:** A REQ completing through the work pipeline carries phase timestamps for the phases its route ran, the completion display shows the phase breakdown, and a historical REQ without the fields still parses, verifies, and logs calibration exactly as today.
**Validation:** Inferred during capture

## Full Context
See `do-work/user-requests/UR-084/input.md` for complete verbatim input.

---
*Source: "Record timestamps for planning, dispatch, builder handback, integration, review, remediation, re-review, and release; retain claimed_at → completed_at for calibration, but display the phase breakdown so total wall time is not mislabeled as implementation time. […] Keep historical REQs compatible when phase timestamps are absent."*

## Triage

**Route: C** - Complex

**Reasoning:** This adds persistent request-schema observations, changes crash-recovery semantics and every major work phase, and requires parser/display lock-step plus historical compatibility.

**Planning:** Required

## Plan

1. Define the eight flat optional milestone fields and their absent/invalid semantics in the canonical request schema; explicitly strip them during recovery and blocked release.
2. Add precise successful-event stamping points to the work pipeline and a completion breakdown that labels `claimed_at → completed_at` as calibration wall span, never implementation time.
3. Parse the optional fields in the board, derive ordered elapsed-since-previous milestone data while leaving duration calibration unchanged, and publish/render the breakdown only when observations exist.
4. Add model, duration, and production-JavaScript regressions for a full Route C flow, skipped phases, malformed/missing stamps, unchanged calibration, and historical records.

**Architectural decisions:** Flat top-level scalars preserve strict YAML and line-recovery compatibility. Milestones record completed events rather than inferred phase starts. Missing phases are omitted, not zero-filled. `release_at` follows successful release; `completed_at → release_at` is a separately labeled release tail.

**Validation warning:** The plan spans more than five underlying file tasks because schema, lifecycle, Go data, generated payload, and production JavaScript must change in lock-step; splitting would leave an internally inconsistent intermediate contract.

*Generated by Plan agent*

## Exploration

The canonical schema and recovery rules live in `work-reference.md`; stamping and completion reporting live in `work.md`. The board already centralizes claim-to-completion calculation in `durations.go`, ticket parsing in `model.go`, generated drawer data in `generate.go`, and drawer/card wording in the shipped web modules. Existing model/duration/JavaScript behavior tests provide the compatibility seams.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work-board/tools/queue-kanban/model.go`
- `skills/do-work-board/tools/queue-kanban/model_test.go`
- `skills/do-work-board/tools/queue-kanban/durations.go`
- `skills/do-work-board/tools/queue-kanban/durations_test.go`
- `skills/do-work-board/tools/queue-kanban/generate.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js`
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`

**Acceptance criteria:** All eight optional milestone stamps have precise write/reset/read semantics; observed phase breakdowns render without fabricated skipped phases; historical files remain valid; calibration continues to use only `claimed_at → completed_at`; user-visible copy calls that interval wall time.

## Root Cause

The work lifecycle persisted only claim and terminal timestamps, while the queue board presented their whole wall interval with implementation-oriented copy. There was no schema contract or parser path for the successful milestones inside that interval, so planning, hand-back, integration, review, remediation, and release time could neither be distinguished nor safely omitted on older requests.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)

**What was done:** Defined and wired eight optional successful-event timestamps, including recovery and blocked-release reset rules. The board now derives an ordered observed milestone breakdown, omits skipped or malformed observations, separates the release tail, and labels the unchanged claim-to-completion calibration interval as wall time.

## Qualification

Passed — 10 declared files verified against four acceptance criteria and the P-A-U loop. Historical records remain additive-compatible; invalid or absent observations create no fabricated phase; future stamps retain the existing diagnostic path; and calibration still calls the unchanged `measureImplementationSpan` over only `claimed_at` and `completed_at`. No unrelated paths were included.

## Testing

**Tests run:** focused phase parser/derivation/generated-payload/production-JavaScript tests; `node --check` on both changed JavaScript modules; `git diff --check`; `go test ./... -count=1` in `skills/do-work-board/tools/queue-kanban`; canonical maintainer verification
**Result:** ✓ All passing. The full queue-board module suite passed in 141.198s; canonical maintainer verification passed its ordinary and strict JavaScript lanes plus the full CLI suite. The optional external-browser lane was unavailable and skipped.

**Red-green validation:**
- Optional milestone model and phase-breakdown regressions describe data the baseline parser and payload did not expose; they pass with the new schema/parser/derivation path.
- Production JavaScript behavior now renders only observed milestones and calls the claim-to-completion interval wall time, including a separately measured release tail.

**New tests added:**
- Eight-field parsing plus historical absence compatibility
- Full Route C ordering and unchanged 75-minute calibration span
- Skipped/malformed phase omission without zero fabrication
- Generated payload and production detail-rendering behavior

**Existing tests updated (cross-REQ impact):**
- Done-card behavior expectations now assert `wall time` instead of `took`, matching the corrected meaning.

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-01T23:10:36Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 99% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — all successful milestones have explicit write/reset/read behavior, board output is additive for historical requests, and the calibration interval is no longer presented as active implementation time.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by work action*

## Orientation

The do-work lifecycle and queue board now expose where observed pipeline wall time went without changing duration calibration. The kanban-board prime remains current.
