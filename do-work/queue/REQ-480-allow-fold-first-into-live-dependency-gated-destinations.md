---
id: REQ-480
title: '[impact-rule-change] Allow fold-first conversion into live dependency-gated destinations'
status: pending
created_at: 2026-09-01T10:47:44Z
user_request: UR-088
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
impact: impact-rule-change
effort_estimate: effort-mechanical
write_set: [skills/do-work/actions/capture-reference.md, skills/do-work/actions/review-work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh]
---

# Allow Fold-First Conversion into Live Dependency-Gated Destinations

## What

Amend the Fold-First Rule's destination 2 so a matching non-sweep REQ that is dependency-gated can still receive a fold when the destination's chain is alive and the finding is not critical, instead of being treated as "no match" and minting a duplicate.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Amend `skills/do-work/actions/capture-reference.md` § Fold-First Rule, destination 2: a matching non-sweep REQ that is dependency-gated is fold-eligible when all of:
  - (a) the root cause matches;
  - (b) the destination's dependency chain is alive — every `depends_on` member is terminal-successful or present in the queue/working set (a `failed` or `cancelled` member keeps the destination ineligible);
  - (c) the finding's judged impact is not `impact-critical` (critical keeps current behavior — never deferred behind a gate; the existing escalation rule still applies to folds).
- Keep the unassigned (`assigned_to` absent) requirement unchanged.
- Update any restatement of the eligibility condition: `skills/do-work/actions/review-work.md` Step 10 and `skills/do-work/actions/work-reference.md` follow-up flows cite the rule by name — verify none restates the gate.
- Add a contract-regression predicate (`_dev/tests/contract-regressions.sh`) if the builder judges the new wording load-bearing.

## Constraints

- This widens fold eligibility only; it never narrows escalation, never folds a critical finding behind a gate, and never changes destination 1 (sweep append) or destination 3 (prose backlog).

## Builder Guidance

Certainty level: Firm — the exact conditions were decided with the maintainer (2026-09-01 session).

## Red-Green Proof

**RED prompt/case:** Re-run the fold-first scan for REQ-414's re-review finding 4 (differential parity) with REQ-420 pending behind REQ-419: destination 2 declares the match ineligible ("otherwise treat it as no match") and the flow mints a new REQ.
**Why RED now:** That exact path minted REQ-464 and REQ-465, which the maintainer then folded into REQ-420 by hand and cancelled via the canonical cancel transaction (commit 593c5145) — each of their bodies records "cannot receive this finding under the fold-first conversion rule".
**GREEN when:** The same scan folds the finding into REQ-420 (alive chain, non-critical impact, unassigned), and the amended rule states the three conditions with the failed/cancelled-member exclusion.
**Validation:** User confirmed (approved plan, 2026-09-01 session).

## Full Context

See `do-work/user-requests/UR-088/input.md` for complete verbatim input.

---
*Source: UR-088 (Lessons routing with token-budgeted mandatory reads and a fold-gate fix)*

## Addendum (2026-09-01)

User added (v4 revision, validate-feedback Finding 8 — Accept):

> ```
> (b) the finding is goal-shaped, not defect-shaped — it restates, refines, or extends the destination's acceptance goals rather than reporting behavior that is broken in what currently ships; a shipped-behavior defect keeps today's behavior (minted standalone) regardless of impact, because a fix must never wait behind a gate; [...] A fold accepted into a gated destination is recorded in the destination body with the fold date and source, so a stalled chain carrying folds is visible to review-work Step 10 (and to the board, if it already surfaces gated REQs). [...] add a contract-regression predicate pinning conditions (b)–(d) if the builder judges the wording load-bearing.
> ```

- New condition inserted as (b): the finding must be **goal-shaped** — it restates, refines, or extends the destination's acceptance goals. A **shipped-behavior defect** (behavior broken in what currently ships) keeps today's behavior — minted standalone — regardless of impact, because a fix must never wait behind a gate. The original conditions (b) and (c) re-letter to (c) and (d); the amended rule states all four.
- A fold accepted into a gated destination is recorded in the destination body with the fold date and source, so a stalled chain carrying folds stays visible to review-work Step 10 (and to the board, if it already surfaces gated REQs).
- The optional contract-regression predicate pins conditions (b)–(d).
- The Red-Green Proof is unchanged and still valid: REQ-464/465 were goal-shaped (duplicates of REQ-420's acceptance goals) — exactly the case (b) admits.
- Provenance: validate-feedback 2026-09-01, Finding 8. Surface-cost: Earned for the fold-record line — the REQ-464/465 cleanup required reconstructing fold provenance by hand (commit 593c5145); the record is one line per accepted fold. The shape condition itself narrows scope back to the original rule's rationale (N/A).
