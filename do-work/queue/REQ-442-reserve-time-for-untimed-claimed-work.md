---
id: REQ-442
title: 'Reserve forecast time for claimed work without a parseable stamp'
status: pending
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
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-443, REQ-444]
batch: accepted-feedback-regressions
---

# Reserve Forecast Time for Claimed Work Without a Parseable Stamp

## What

When a claimed REQ lacks a parseable `claimed_at`, reserve one projected median span from `now` before scheduling pending work. Keep the existing timestamp defect diagnosis separate from this conservative forecast fallback.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
