---
id: REQ-481
title: 'One-time required-lessons stamp of the pending queue'
status: cancelled
created_at: 2026-09-01T11:04:54Z
user_request: UR-089
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-478]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-477, REQ-478, REQ-479]
batch: lessons-transfer-routing
completed_at: 2026-09-01T12:10:11Z
---

# One-Time Required-Lessons Stamp of the Pending Queue

## What

When REQ-478's stamping mechanism lands, run the same stamping decision once over every pending, unassigned REQ already in `do-work/queue/` and stamp them. A single pass, not a standing behavior.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Run REQ-478's stamping decision (index read, relevance ranking, token budget, entry forms, full-only targeting) once over every `pending`, unassigned (`assigned_to` absent) REQ in `do-work/queue/`.
- Each retroactively stamped REQ notes the retroactive pass in its body.
- REQ-457 is the motivating case — it is queued for a family whose lesson already exists (final-boundary identity/rollback in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`).
- A single pass, not a standing behavior: this REQ adds no recurring mechanism, hook, or schedule.

## Constraints

- Edit only pending, unclaimed queue REQs — the same legal window the fold-first append uses; never a claimed, blocked-holding, or archived file.
- A no-match decision is a valid outcome; never invent a stamp.

## Dependencies

Depends on REQ-478 (the stamping mechanism and `required_lessons` field this pass applies).

## Builder Guidance

Certainty level: Firm. This is the backfill that lets the current queue — the REQs the 2026-09-01 analysis showed relearning recorded families — benefit from the mechanism instead of only future captures.

## Red-Green Proof

**RED prompt/case:** After REQ-478 lands, inspect `do-work/queue/REQ-457-record-cleanup-move-destinations-after-exclusive-creation.md`: no `required_lessons` despite an existing same-family lesson set in the satellite.
**Why RED now:** The stamping mechanism runs only at capture time; nothing revisits already-queued REQs.
**GREEN when:** Every pending, unassigned queue REQ has been through the stamping decision exactly once; stamped REQs (REQ-457 included, unless the decision honestly finds no match) carry the field plus a body note of the retroactive pass.
**Validation:** User confirmed (v4 revision via validate-feedback, 2026-09-01).

## Full Context

See `do-work/user-requests/UR-089/input.md` for complete verbatim input.

---
*Source: UR-089 (v4 triage: lessons-routing refinements and fold-gate shape condition)*

## Cancelled

- **When:** 2026-09-01T12:10:11Z
- **Why:** superseded by claim-time consult (REQ-479 addendum, 2026-09-01)
- **Decided by:** user, via `do-work abandon`
