---
id: REQ-530
title: '[impact-critical] Anchor every lifecycle gate that reads caller-authored text'
status: pending
created_at: 2026-09-03T09:45:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-528]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-528]
sweep: true
sweep_key: answer-line-marker-position-spoofing
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go
---

# Anchor Every Lifecycle Gate That Reads Caller-Authored Text

## What

REQ-528 fixed one instance of a shape: a lifecycle decision made by searching caller-authored text for a token, anywhere in that text. The shape-grep it ran while fixing that instance found three more, one of which drives a second terminal status in the same file. This sweep closes them all against the same invariant: **evidence a caller writes freely cannot be the anchor for a lifecycle write.**

## Instances

- **`internal/publication/answer.go:268` — a second terminal status, same file.** The stakeholder branch gates `status: completed`, `completed_at`, deletion of `blocked_by`/`blocked_at`, and the archive move on `bytes.Contains(ToLower(blockedHistory), "resolved")` **and** `bytes.Contains(ToLower(implementation), "no code")`. Both payloads are caller-authored prose and both tokens match anywhere, so `still not resolved` satisfies the first and `no code review yet` satisfies the second. The refusal message ("terminal evidence must carry resolved Blocked history and an Implementation no-code marker") names markers, but the code tests for substrings.
- **`internal/publication/answer.go:291` — non-terminal but still a lifecycle write.** The `blocked_by` linkage is gated on `bytes.Contains(reportsHistory, reportPath)` over caller-authored history text, unanchored to any line or position.
- **`internal/cleanup/cleanup_plan.go:236` — adjacent, and destructive.** A `do-work/CHECKPOINT.md` line is deleted when it merely *contains* `- <REQ-ID>:` plus the writer token, with no position anchor. Narrower than the others because both tokens must share one line, but it is the same unanchored-line-selection shape driving a destructive edit.
- **Bounded, listed for completeness rather than repair:** `internal/finalization/finalization_discovery.go:593` admits an unjournaled changelog tail into replay on `bytes.Contains(inserted, requestID)`. `singleInsertion` bounds the searched bytes to one verified diff hunk, so the containment is already anchored by its caller. Confirm that rather than assume it, then leave it.
- **Checked and deliberately not included:** `answer.go:406` (`findQuestionLine`) and `internal/corehelpers/checks.go:232` are unanchored but fail closed — the first requires exactly one matching line and refuses otherwise, the second only suppresses an `OK:` line in a report it composed itself.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- Findings **F-01 through F-04** from REQ-528's implementation shape-grep, which was asked for exactly because the fixed bug had a shape rather than being a one-off. REQ-528 fixed the reported instance only; this sweep owns the rest.
- The `impact-critical` token is carried from the first instance: it decides a terminal status and an archive move on evidence its own caller supplies.

## Detailed Requirements

- Every lifecycle or destructive decision that reads caller-authored text must anchor the token to a position the reader can attribute — a line start, a field boundary, or a separator the writer itself contributes. Containment anywhere in free text is not evidence.
- Where a token is genuinely a marker, test it as a marker: anchored, and at the position whatever wrote it places it. Where the refusal message already says "marker", the code must match that claim.
- An unattributable read must fail toward the non-destructive, non-terminal outcome, as REQ-528 established.
- Preserve every existing refusal code and typed result; these are matching fixes, not policy changes.
- Cover each instance with a test that forges the token in caller text and asserts the lifecycle write does **not** happen, plus a control proving the genuine path still does.

## Constraints

- Do not introduce a new schema field or change any stored format unless a matching fix genuinely cannot anchor the token; say so explicitly if you reach that conclusion.
- `internal/cleanup` and `internal/publication` are separate packages; honor the prime's **Package direction** rule rather than reaching for a shared helper that would violate it.
- Do not weaken `singleInsertion`'s bounding in `finalization_discovery.go`.

## Dependencies

Depends on REQ-528, which establishes the anchoring invariant and the fail-toward-non-terminal rule this applies.

## Red-Green Proof

**RED prompt/case:** Submit a terminal stakeholder disposition whose Blocked-history payload reads `still not resolved` and whose Implementation payload reads `no code review yet`, then inspect the resulting `status`.
**Why RED now:** both gates use `bytes.Contains` over caller-authored prose, so both negations satisfy them and the REQ reaches `completed` with `completed_at` set and the archive move planned.
**GREEN when:** that submission is refused with the existing `ANSWER-STAKEHOLDER-EVIDENCE-INVALID` code; a genuine terminal evidence pair still reaches `completed`; the `blocked_by` linkage and the CHECKPOINT line deletion are likewise anchored with forge-and-control tests; and `singleInsertion`'s bounding is confirmed intact.

---
*Source: REQ-528 implementation shape-grep, findings F-01 through F-04.*
