---
id: REQ-484
title: '[impact-rule-change] Add same-UR verification-append destination to Fold-First Rule'
status: pending
created_at: 2026-09-01T12:11:03Z
user_request: UR-091
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-420]
batch: lessons-transfer-routing
---

# Add Same-UR Verification-Append Destination to Fold-First Rule

## What

Add a new destination to the Fold-First Rule between destination 2 (convert) and
destination 3 (prose backlog): a goal-shaped, non-critical finding from inside UR-N
appends to UR-N's pending downstream verification/parity/acceptance member as a
`## Folded From` acceptance-criteria section — never a sweep conversion; the
destination keeps its `write_set`, `estimate`, frontmatter, and non-sweep shape.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in
any UR sharing this root cause; REQ-480, which loosened destination 2's conversion
gate for the same evidence, was cancelled 2026-09-01 as the wrong mechanism before
this capture, and this REQ replaces it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Amend `skills/do-work/actions/capture-reference.md` § Fold-First Rule with the new
  destination. Eligible destination: a pending, unclaimed (`assigned_to` absent)
  same-UR downstream member whose stated scope is verification, parity, or
  acceptance over the code the finding names. Conditions, all required:
  - (a) the finding is **goal-shaped**, not defect-shaped — it restates, refines, or
    extends the destination's acceptance goals rather than reporting behavior broken
    in what currently ships. "Shipped behavior" means what currently ships to
    consumers: a divergence in not-yet-authoritative migration code, whose legacy
    implementation still ships until the destination REQ converts it, is goal-shaped
    relative to the migration's acceptance. A shipped-behavior defect never folds
    behind a gate regardless of impact and is minted standalone as today.
  - (b) the destination's dependency chain is **alive** — every `depends_on` member
    is terminal-successful or present in the queue/working set; a `failed` or
    `cancelled` member keeps the destination ineligible.
  - (c) the finding's judged impact is not `impact-critical` — critical keeps
    current behavior (own REQ, `status: pending`, prominent pierce line) at any depth.
- The append form is a `## Folded From REQ-NNN / <source> (YYYY-MM-DD)` section
  restating the finding as acceptance criteria or named fixtures with the fold date
  and source recorded in the destination body, so a stalled chain carrying folds
  stays visible to `actions/review-work.md` Step 10 (and any board surface showing
  gated REQs). Never the destination-2 sweep conversion.
- Update every flow that restates fold-destination eligibility:
  `actions/review-work.md` Step 10 and `actions/work-reference.md`'s follow-up flows
  cite the Fold-First Rule by name — verify none restates the gate.
- Add a contract-regression predicate (`_dev/tests/contract-regressions.sh`) pinning
  conditions (a)–(c) if the builder judges the wording load-bearing.

## Constraints

- Destination 2's conversion path stays byte-for-byte in meaning for dependency-ready
  same-root-cause matches; destinations 1 and 3 are untouched.
- This widens fold eligibility only; escalation and the critical pierce never narrow.

## Builder Guidance

Certainty level: Firm — conditions decided with the maintainer (2026-09-01 session,
approved plan). Latitude: exact wording and placement within the Fold-First Rule.

## Red-Green Proof

**RED prompt/case:** Re-run the fold-first scan for REQ-414's re-review finding 4
(differential parity) with REQ-420 pending behind REQ-419: no destination accepts
it — destination 2 both declares the gated match ineligible and, were the gate
lifted, would convert the terminal parity REQ into a sweep and clear its `write_set`.
**Why RED now:** That exact path minted REQ-464/REQ-465, which the maintainer folded
into REQ-420 by hand as `## Folded From` acceptance criteria and cancelled (commit
593c5145); the 2026-09-01 disposition repeated the same hand-fold for REQ-467,
REQ-473, REQ-474, and REQ-476.
**GREEN when:** The same scan routes such a finding through the new destination
automatically — appended, recorded with date and source, no conversion — and the
amended rule states conditions (a)–(c) with the failed/cancelled-member exclusion.
**Validation:** User confirmed (approved plan, 2026-09-01 session).

## Full Context

See `do-work/user-requests/UR-091/input.md` for complete verbatim input.

---
*Source: UR-091 (Add a same-UR verification-append destination to the Fold-First Rule)*
