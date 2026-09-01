---
id: REQ-470
title: 'Fold-first routing of unrelated gate failures into pending-answers REQs'
status: pending
created_at: 2026-09-01T04:29:16Z
user_request: UR-087
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-469]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-468, REQ-469, REQ-471, REQ-472]
batch: non-blocking-orchestration
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/capture-reference.md, _dev/tests/contract-regressions.sh]
---

# Fold-First Routing of Unrelated Gate Failures Into Pending-Answers REQs

## What

When the blocked set-aside (REQ-469) records unrelated canonical-gate failures, run the fold-first scan for each one: append it to a matching queued REQ when possible; otherwise mint a non-critical `pending-answers` REQ so the user can approve or reject fixing it. The failures must become queue work, never live only in `CHECKPOINT.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- "Run the fold-first scan for each unrelated gate failure. Append to a matching queued REQ when possible; otherwise create a non-critical `pending-answers` REQ so the user can approve or reject fixing it."
- Acceptance test: "The unrelated failures create or fold into `pending-answers` REQs instead of being stored only in `CHECKPOINT.md`."

## Constraints

- Use the canonical Fold-First Rule (`actions/capture-reference.md` → Fold-First Rule) as-is — this is a new call site, not a new mechanism. The rule's condition ("any flow about to mint a REQ from a finding… first scans the queue") already covers it; add the set-aside flow as a caller from `actions/work.md` rather than restating the ladder.
- A gate failure is a behavioral finding, so the prose-only destination will not normally apply; the minted no-match REQ is `pending-answers` with a consent question (approve/reject fixing the unrelated failure), consistent with the existing consent-gated pattern, and non-critical unless the failure is judged `impact-critical` (the rule's escalation clause then applies as written).
- Minted REQs follow the REQ Title Convention and carry their own judged `impact:`/`effort_estimate:`; provenance cites the set-aside REQ and UR per the fold citation format.
- Contract-regressions assertions pinning any edited prose change in the same commit.

## Dependencies

- REQ-469 (blocked set-aside) — this REQ wires the fold-first scan into the set-aside path REQ-469 creates.

## Builder Guidance

Certainty: Firm. Keep it small: one call-site paragraph in the set-aside flow plus whatever cross-reference the Fold-First Rule's caller list convention requires (the condition is the rule; avoid growing a closed enumeration).

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** After REQ-469's set-aside, nothing instructs routing the unrelated failing tests anywhere durable — today's hold records them at best in `CHECKPOINT.md` phase detail, which the next session's checkpoint rewrite or recovery can discard.
**Why RED now:** Unrelated gate failures have no path into the queue, so the gate stays red until a human notices and hand-captures the fixes.
**GREEN when:** The set-aside flow instructs the fold-first scan per failure with the pending-answers no-match destination; a contract assertion pins the call site; `bash _dev/tests/contract-regressions.sh` exits zero.
**Validation:** Inferred during capture (from the spec's acceptance tests)

## Full Context
See `do-work/user-requests/UR-087/input.md` for complete verbatim input.

---
*Source: UR-087 — "Run the fold-first scan for each unrelated gate failure…"*
