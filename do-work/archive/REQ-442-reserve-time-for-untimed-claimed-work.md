---
id: REQ-442
title: 'Reserve forecast time for claimed work without a parseable stamp'
status: completed
route: A
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T22:03:12Z
  basis:
    - trivial short-circuit
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-443, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T22:02:33Z
completed_at: 2026-09-01T22:13:37Z
---

# Reserve Forecast Time for Claimed Work Without a Parseable Stamp

## What

When a claimed REQ lacks a parseable `claimed_at`, reserve one projected median span from `now` before scheduling pending work. Keep the existing timestamp defect diagnosis separate from this conservative forecast fallback.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** In `timelineChainStart`, reuse each claimed ticket's existing bucket median and fall back to `now` only when `claimed_at` is absent or malformed. Add test-first coverage for missing/malformed stamps, maximum finish across claims, dependent scheduling, parseable behavior, and independent verify findings.
- [x] **[APPLY]:** Replaced the untimed-claim skip with a conservative `now + existing bucket median` finish and added the combined forecast/diagnosis regression.
- [x] **[UNIFY]:** Reviewed `timeline.go` and `timeline_test.go` in full diff context; verified median reuse, maximum finish behavior, unchanged parseable claims, unchanged timestamp bytes/findings, and no debug artifacts. Focused projection tests pass.

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Reserve time for claimed work without a parseable stamp.`
- **Evidence:** `timeline.go:397-413` starts the chain at `now` and skips claimed tickets whose `claimed_at` cannot be parsed.
- **Origin / earned by:** `2daefd1c`/REQ-228 (Project the Remaining Queue onto the Timeline) defined `claimed_at + median` for timed work without defining the invalid-stamp fallback. A reproduced claimed prerequisite with no stamp placed its pending dependent exactly at `now` while separately emitting the timestamp defect.
- **Surface-cost:** Earned. One conservative fallback and missing/malformed regressions are cheaper than overlapping active work, declining the whole forecast, or adding timestamp-recovery machinery.

## Detailed Requirements

- Treat every claimed ticket with absent or malformed `claimed_at` as occupying `now` through `now +` its projected effort-bucket median.
- Use the maximum finish across multiple in-flight claims, preserving current behavior for parseable timestamps.
- Continue emitting timestamp-quality findings independently; do not invent or repair a stored timestamp.
- Keep pending and dependent work from beginning before the conservative untimed-claim finish.

## Constraints

- This is a forecast assumption, not metadata repair.
- Reuse existing projection medians and effort-bucket selection.

## Red-Green Proof

**RED prompt/case:** Generate a projection with completed samples, one claimed REQ missing or carrying malformed `claimed_at`, and one pending dependent.
**Why RED now:** The claimed ticket is skipped, so chain start and the dependent's start remain exactly `now`.
**GREEN when:** The dependent starts no earlier than `now +` the claimed ticket's projected median, timestamp diagnosis still appears, and parseable-claim projections remain unchanged.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. Use the existing median projection; do not add nested-frontmatter parsing or timestamp recovery.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 22 from the validated external feedback.*

## Triage

**Route: A** - Simple

**Reasoning:** The forecast already centralizes claimed-work reservation and bucket medians; the defect is one skip branch plus focused projection/diagnosis coverage.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

`timelineChainStart` treated an unparseable `claimed_at` as grounds to skip the claimed ticket entirely. Dependency placement already considered claimed prerequisites resolvable, so this skip made their dependents schedulable at `now` even though active work still occupied the serial forecast.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified)

**What was done:** Claimed tickets without a parseable claim stamp now reserve one full projected bucket-median span from `now`; parseable claims retain their original claim-time-plus-median finish. The regression covers absent and malformed stamps, multiple in-flight maxima, dependent placement, unchanged timestamp diagnosis, and no metadata repair.

## Qualification

Passed — 2 files verified, 4 requirements traced, P-A-U confirmed. The implementation reuses the existing effort-bucket projection, the timestamp verifier remains independent, and the unrelated knowledge-base batch was excluded from this REQ's evidence and staging.

## Testing

**Tests run:** focused timeline projection tests; `go vet ./... && go test ./... -count=1` in the queue-kanban module; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing, including strict JavaScript behavior and canonical maintainer verification. The optional browser lane was unavailable and skipped by the gate; this forecast-model change has no browser-only acceptance condition.

**Red-green validation:**
- `TestTimelineProjectionReservesTimeForUntimedClaimedWork`: ✗ before implementation (chain began at `now + 35m`, skipping the untimed 45-minute claim) → ✓ after (chain begins at the maximum conservative finish, `now + 45m`)

**New tests added:**
- Missing and malformed claim stamps across substantive/mechanical medians, multiple in-flight claims, dependent scheduling, parseable behavior, verify findings, and unchanged metadata

*Verified by work action*

## Review

**Overall: 98%** | 2026-09-01T22:13:06Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — untimed claims reserve their existing bucket median from `now`, the latest in-flight finish controls the chain, and diagnosis/metadata behavior remains separate.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

The board's serial queue forecast now reserves time for active claimed work even when its claim timestamp cannot be trusted. The Kanban board prime remains current.
