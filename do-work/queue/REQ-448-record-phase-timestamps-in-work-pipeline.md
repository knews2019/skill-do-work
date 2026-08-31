---
id: REQ-448
title: 'Record per-phase timestamps through the work pipeline'
status: pending
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
related: [REQ-449]
---

# Record Per-Phase Timestamps Through the Work Pipeline

## What

Record optional per-phase timestamps on a REQ as it moves through the work pipeline — planning, dispatch, builder handback, integration, review, remediation, re-review, release — and display the phase breakdown at completion, so total wall time is not mislabeled as implementation time. `claimed_at` → `completed_at` stays the calibration span.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
