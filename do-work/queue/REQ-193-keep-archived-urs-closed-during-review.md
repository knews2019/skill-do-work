---
id: REQ-193
title: Keep archived URs closed during standalone review
status: pending
created_at: 2026-08-15T09:12:04Z
user_request: UR-043
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-194]
batch: closed-ur-documentation-hardening
write_set: [skills/do-work/actions/review-work.md, skills/do-work/actions/work.md, _dev/tests/contract-regressions.sh]
---

# Keep Archived URs Closed During Standalone Review

## What

Co-locate the archived-UR invariant with standalone review's archived-input read so a reviewer can create a follow-up without moving or reopening the already-closed User Request folder. Make the downstream archive instruction explicit enough that the completed follow-up returns to that folder in place.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The invariant already exists in abandon and capture, but standalone review loads neither action. The board correctly caught the resulting misplaced REQs; prevention belongs beside the review instruction that sends an agent into the archive.

## Context

- Original priority: Primary; no P-level severity was supplied.
- Verbatim claim: “Missing: any statement that the archived UR folder stays put.”
- Evidence: `skills/do-work/actions/review-work.md:60` reads archived input and `:412` queues follow-ups, while `skills/do-work/actions/abandon.md:129` and `skills/do-work/actions/capture.md:156` state the closed-folder invariant.
- Surface-cost: Earned. Commit `1323982` records a real misplaced-REQ replay case, and the long-lived addition is a co-located instruction plus one focused contract assertion.

## Detailed Requirements

- Add the invariant immediately beside Step 3's fallback read from `do-work/archive/UR-NNN/input.md`.
- State that reading archived input is context-only and grants no authority to move, reopen, or re-consolidate the archived UR folder.
- Preserve the established review-follow-up model: the new REQ keeps the same `user_request`, carries `review_generated: true`, and goes into `do-work/queue/`.
- Distinguish standalone review of a closed UR from orchestrated review while the UR is still open.
- Make the work/archive path explicit for a review-generated REQ whose UR already lives in `archive/`: move the completed REQ into that existing UR folder in place and never move the UR folder itself.
- Preserve review's existing narrow exception that permits appending post-work review metadata to an archived REQ.
- Add a co-located contract-regression assertion that fails if the archived-input block loses the stays-closed instruction.

## Constraints

- This is documentation/lifecycle hardening, not a claim that queue-kanban behaved incorrectly.
- State the rule at the review consumer that needs it; do not restate it across unrelated actions.
- Do not create a new UR for standalone-review follow-ups. The user confirmed the same-UR/stays-closed model.
- Do not edit downstream consumer repositories; suite update propagation is outside this REQ.

## Dependencies

None. REQ-194 depends on this REQ's canonical lifecycle statement.

## Builder Guidance

Firm intent. Prefer the smallest wording and archive-branch clarification that make the same-UR/stays-closed path unambiguous.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Add a `review_archived_input_block` probe to `_dev/tests/contract-regressions.sh` that extracts `review-work.md` Step 3 and requires the archived UR to stay closed while follow-ups go to the queue.
**Why RED now:** The extracted block names the archive fallback but contains no prohibition on moving or reopening that folder, so the focused assertion fails.
**GREEN when:** The same probe passes because the invariant is co-located in Step 3, the work archive text covers the already-archived review-generated case in place, and the baseline contract suites remain green.
**Validation:** User confirmed the contract assertion and same-UR/stays-closed lifecycle on 2026-08-15. `tdd: false` because the primary change is operating documentation; the named check still supplies finding-closure proof.

## Assets

None.

## Full Context

See `do-work/user-requests/UR-043/input.md` for the complete verbatim request, validated evidence, and batch constraints.

---
*Source: documentation-hardening finding validated through `do-work-toolbox validate-feedback`; original priority Primary.*
