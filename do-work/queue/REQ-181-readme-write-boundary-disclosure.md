---
id: REQ-181
title: README understates install and implementation write boundaries
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
effort_estimate: trivial
related: [REQ-182, REQ-183, REQ-184, REQ-185, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [README.md]
---

# README Understates Install and Implementation Write Boundaries

## What

Correct the public README's absolute write-scope claim so adopters see the actual managed-install, durable queue-state, and explicitly scoped implementation boundaries before they trust the suite.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The README tells existing-project adopters that the skill only creates files under `do-work/` and merely reads source. Install/update actually manage four `.claude/skills/` trees plus consented Just/settings sections, and `do-work run` can write project paths authorized by the active request's declared Scope.

## Context

- Audit priority: P1; impact 4; effort trivial.
- Root-cause key: `public-write-boundary-disclosure`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 1, audited at `58eb4f84f408dce1ec9828a07aef0b174930ce34`.
- Reproduce: `sed -n '24,35p;156,159p' README.md && sed -n '347,365p' skills/do-work/actions/work.md && sed -n '35,61p' decisions/records/adr-019-four-skill-suite-contract.md`.

## Detailed Requirements

- Replace the statement `The skill only creates files inside` at `README.md:158` at the audited SHA.
- Correct the install description at `README.md:24-31`, which names only `.claude/skills/do-work` rather than the managed four-skill suite.
- State the three actual boundaries:
  - reviewed and confirmed managed install paths;
  - durable core queue state under `do-work/`;
  - project-source writes only during explicitly invoked REQ implementation within that request's declared `## Scope`.
- Keep the wording consistent with the Builder Scope Fence in `skills/do-work/actions/work.md:357-365` and ADR-019.

## Constraints

- This is one direct public-prose correction.
- Do not add a generated safety-summary checker; its surface would cost more than the stable boundary statement.
- Do not change runtime behavior or expand the managed/write boundaries.

## Dependencies

None. The batch's other findings may proceed independently; shared file overlap is surfaced in frontmatter.

## Builder Guidance

Firm intent. Prefer the smallest accurate wording change and preserve the README's existing voice.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Read `README.md:24-31` and `:156-159`, then compare those claims with the Builder Scope Fence and ADR-019.
**Why RED now:** The README makes an absolute `do-work/`-only claim that omits managed suite installation and request-scoped project writes.
**GREEN when:** The README visibly distinguishes all three real boundaries and no longer makes the contradicted absolute claim; the exact old sentence is absent.
**Validation:** Confirmed by the user during verification on 2026-08-15. No runnable test is justified for this prose-only correction.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows all eight collapsed findings in the generated audit report. This request is row 01, labeled P1, impact 4, trivial effort, with the title matching this REQ.

## Full Context

See `do-work/user-requests/UR-041/input.md` for the complete verbatim request and batch constraints. See the canonical audit for the validated claim, evidence, and pre-emption record.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
