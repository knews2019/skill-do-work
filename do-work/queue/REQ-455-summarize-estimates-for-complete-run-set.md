---
id: REQ-455
title: 'Summarize estimates for the complete run set'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-450, REQ-451, REQ-452, REQ-453, REQ-454, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Summarize Estimates For The Complete Run Set

## What

Compute estimate totals, unknown-estimate evidence, and critical-path information over the complete authoritative run set rather than only the immediately selected subset. Retain estimate and dependency evidence for members deferred behind prerequisites so callers never need a prohibited queue rescan.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this incomplete-run estimate-summary root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #14 — P2 — source:** `internal/nextselection/next_selection.go:54-59`

> ````text
> [P2] Summarize estimates for the complete run set — [prj].claude/skills/do-work/tools/do-work-cli/
> internal/nextselection/next_selection.go:54-59
> The summary adds estimates only from Selected, but default serial selection limits that collection to one ready request and
> dependency-blocked members carry no estimate evidence in Excluded. Consequently multi-REQ runs underreport total effort and
> unknown estimates, and the authoritative result cannot produce the required full-run critical-path summary without the prohibited
> queue rescan.
> ````

- **Evidence:** `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go:54-60` totals only `Selected`. Dependency-blocked candidates exit at lines `210-216` before contributing estimate evidence. `internal/resultmodel/result_model.go:105-120` gives `SelectionExclusion` neither estimate nor dependency information, and the summary at lines `124-136` has no critical path. The required complete-run summary is specified at `skills/do-work/actions/work.md:172-183`.
- **Surface-cost result:** N/A — direct correction of the authoritative result model and execution contract.

## Detailed Requirements

- Define the complete run set for default, targeted REQ, targeted UR, and multi-token selection modes.
- Aggregate known estimates across that complete set, not only immediately runnable records.
- Count and identify unknown estimates across the complete set.
- Preserve enough dependency evidence to compute and report the full-run critical path.
- Keep selected-versus-deferred state distinct from membership in the run set.
- Maintain text and JSON result parity without queue rescans.

## Constraints

- Do not inflate parallel critical-path duration by simply summing every member.
- Preserve the existing immediately selected records and exclusion reasons for execution decisions.

## Dependencies

No request prerequisite. This root cause is independently accepted even though its model may overlap REQ-453 and REQ-454.

## Builder Guidance

Certainty level: Firm on the required complete-run outcome; choose the smallest result-model extension that preserves membership, estimates, and dependency edges.

## Red-Green Proof

**RED prompt/case:** Select a multi-REQ run with one immediately ready request, one pending dependent with a known estimate, one dependent with an unknown estimate, and a parallel branch.
**Why RED now:** The summary sees only `Selected`; deferred members lose estimate and dependency evidence.
**GREEN when:** The authoritative text and JSON summaries report complete-set totals, the unknown member, and the correct dependency critical path without rescanning the queue.
**Validation:** User confirmed after validate-feedback accepted Finding #14.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #14, captured by UR-085.*
