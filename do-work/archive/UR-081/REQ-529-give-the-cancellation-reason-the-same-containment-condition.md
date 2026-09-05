---
id: REQ-529
title: 'Give the cancellation reason the same containment condition'
status: cancelled
created_at: 2026-09-03T03:10:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-460]
sweep: true
sweep_key: markdown-delimiter-containment-prefix-gaps
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
status_changed_at: 2026-09-03T10:39:11Z
completed_at: 2026-09-03T11:47:25Z
---

# Give the Cancellation Reason the Same Containment Condition

## What

REQ-460 made the answer-summary inline-or-contain decision condition-complete. The same decision has a second home that has no condition at all: `cancellationReasonBlock` inlines any cancellation reason after `- **Why:** ` on the sole test that it contains no newline.

## Instances

- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:895-903` — `cancellationReasonBlock` inlines on `!strings.Contains(reason, "\n")` alone. A one-line reason of `***` or `## Notes` is written straight into the document.
- That package already has the machinery: its own `containedOutsideText` (`state_apply.go:905`) and its own `validateOutsideText` (`state_plan.go:230`). Only the structural judgment is missing.
- `skills/do-work/actions/abandon.md:66` says "A **safe** one-line reason stays inline after `**Why:**`" — the prose asks for the delimiter judgment the code never makes.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding F2** — `impact-rule-change` — from REQ-460's independent review, found by its restatement sweep across every inline-or-contain decision keyed on a newline test. One home is now condition-complete and its sibling is not.
- This is the `[family: alternate-writer-contract-drift]` shape the CLI prime already names: changing a stored-format contract in one writer without sweeping the others leaves repository behavior split.

## Open Questions

- [x] Should the two seams share one predicate, or keep separate ones with a shared test corpus? Sharing means a new exported helper crossing the `publication` → `requeststate` boundary, or a third package both import; the CLI prime's **Package direction** section constrains which imports are allowed, so this is a structural call rather than a mechanical one. → Confirmed: one shared predicate
  - Recommended: a small shared package (or an exported predicate in whichever package the prime's direction rule permits importing), so the condition is stated exactly once — the whole point of REQ-460.
  - Also: duplicate the predicate in `requeststate` with a shared table-driven corpus, accepting two implementations pinned by one fixture set.
  - Also: leave them separate and unpinned — rejected, since that is the drift being fixed.

## Detailed Requirements

- A one-line cancellation reason that can be read as document structure or a delimiter must be carried through the existing contained-text path, not inlined.
- Reuse this package's existing `containedOutsideText` and `validateOutsideText` rather than adding a third containment writer.
- Preserve lossless bytes for safe plain reasons and for contained reasons alike.
- Table-driven cases spanning the same structural classes REQ-460 pinned, so the two seams cannot drift again without a test failing.

## Constraints

- Do not weaken the existing newline refusal or the C0/DEL validation in `state_plan.go`.
- Honor the CLI prime's **Package direction** rule when deciding where a shared predicate may live.

## Dependencies

No request prerequisite. REQ-460 established the condition this applies.

## Red-Green Proof

**RED prompt/case:** Cancel a REQ with the one-line reason `***`, and again with `## Notes`, then read the archived REQ body.
**Why RED now:** `cancellationReasonBlock` tests only for a newline, so both reasons are written inline after `- **Why:** ` as document structure.
**GREEN when:** Both reasons land through the contained path with bytes intact, an ordinary plain reason still inlines, and the structural corpus is shared with the answer-summary seam so a future divergence fails a test.

---
*Source: REQ-460 independent review finding F2.*


## Answer Notes

- 2026-09-03 - [ ] Should the two seams share one predicate, or keep separate ones with a shared test corpus? Sharing means a new exported helper crossing the `publication` → `requeststate` boundary, or a third package both import; the CLI prime's **Package direction** section constrains which imports are allowed, so this is a structural call rather than a mechanical one.: Confirmed: one shared predicate
> ```
> One shared predicate. A single implementation prevents answer summaries and cancellation reasons from drifting, and the small shared dependency is accepted. This clarification does not add another containment writer; the existing containment and validation paths remain in use.
> ```

## Cancelled

- **When:** 2026-09-03T11:47:25Z
- **Why:** folded into REQ-466 section Folded From REQ-529 as an acceptance criterion (maintainer decision, 2026-09-03 roadmap triage)
- **Decided by:** user, via `do-work abandon`
