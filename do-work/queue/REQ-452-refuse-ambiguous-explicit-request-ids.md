---
id: REQ-452
title: 'Refuse ambiguous explicit request IDs'
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
effort_estimate: effort-mechanical
related: [REQ-450, REQ-451, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Refuse Ambiguous Explicit Request IDs

## What

Preserve duplicate queue-record collision evidence when resolving numeric request IDs, and return an ambiguity exclusion when an explicit `REQ-NNN` token cannot identify exactly one file. Explicit targeting may bypass documented dependency, assignment, and impact gates, but never repository identity ambiguity.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this explicit-target duplicate-ID ambiguity root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #4 — P2 — source:** `internal/nextselection/next_targets.go:30`

> ````text
> [P2] Reject ambiguous IDs even for explicit targets — [prj].claude/skills/do-work/
> tools/do-work-cli/internal/nextselection/next_targets.go:30-30
> When two queued records normalize to the same REQ number, this assignment silently overwrites the first entry; because explicit
> candidates later bypass graph ambiguity checks, next REQ-NNN selects whichever duplicate was encountered last and returns its
> exact path as runnable. Explicit targeting only overrides the documented dependency, assignment, and impact gates, not an
> unresolved record identity (.claude/skills/do-work/actions/work-reference.md:394-399), so detect the duplicate here and emit an
> ambiguity exclusion.
> ````

- **Finding #8 — P1 — source:** `internal/nextselection/next_selection.go:196-202`

> ````text
> [P1] Refuse ambiguous explicitly targeted request IDs — [prj].claude/skills/do-work/tools/do-work-
> cli/internal/nextselection/next_selection.go:196-202
> When duplicate queue records claim the same REQ id, requestByNumber has already selected one arbitrary file, and explicit
> provenance skips the node.IsAmbiguous guard here. Thus next REQ-NNN can authorize processing whichever duplicate happened to
> overwrite the map entry; explicit targeting should bypass dependency gates, not repository identity ambiguity.
> ````

- **Finding #18 — P2 — source:** `internal/nextselection/next_targets.go:29-30`

> ````text
> [P2] Refuse explicit IDs that resolve to duplicate queue files — [prj].claude/skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go:29-30
> When two queued files normalize to the same numeric REQ id, this assignment silently keeps the last one discovered. next REQ-NNN then selects that arbitrary file because explicit provenance bypasses the graph's IsAmbiguous guard, even though the token cannot
> distinguish the records. Preserve collision evidence and return an ambiguity exclusion instead, consistent with prime-do-work-cli.md:14-16 (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L14-L16).
> ````

- **Evidence:** `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go:24-31` builds a numeric map whose assignment overwrites an earlier duplicate. Explicit provenance then skips ambiguity at `next_selection.go:196-203`, despite collision evidence already being available in `internal/repositorymodel/repository_model.go:78-93,239-268`. The chosen file therefore depends on discovery order.
- **Surface-cost result:** Earned — the repository already computes collision evidence. Reusing it and adding one duplicate-explicit-target replay is cheaper than arbitrary request authorization.

## Detailed Requirements

- Preserve every path contributing to a duplicate normalized numeric request ID.
- Reject explicit targeting when the token resolves to more than one queue record.
- Emit a typed ambiguity exclusion rather than returning either duplicate as runnable.
- Keep explicit targeting's documented dependency, assignment, and impact overrides unchanged.
- Make the result independent of repository discovery order.

## Constraints

- Reuse normalized collision evidence rather than inventing a second identity rule.

## Dependencies

No request prerequisite. Shared selector files with other UR-085 requests do not establish dependency ordering.

## Builder Guidance

Certainty level: Firm. Identity ambiguity is a repository integrity condition, not an eligibility gate.

## Red-Green Proof

**RED prompt/case:** Place two queue files that normalize to the same numeric REQ ID, explicitly request that `REQ-NNN`, and repeat with reversed discovery order.
**Why RED now:** Numeric-map assignment silently keeps one path and the explicit path skips the graph's ambiguity check.
**GREEN when:** Both replays return the same typed ambiguity exclusion containing collision evidence, and neither duplicate is selected.
**Validation:** User confirmed after validate-feedback accepted Findings #4/#8/#18.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Findings #4, #8, and #18, captured by UR-085.*
