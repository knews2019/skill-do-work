---
id: REQ-181
title: README understates install and implementation write boundaries
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T09:41:50Z
completed_at: 2026-08-15T09:48:44Z
commit: 86fbcca
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
route: A
kb_status: skipped
kb_entry:
---

# README Understates Install and Implementation Write Boundaries

## What

Correct the public README's absolute write-scope claim so adopters see the actual managed-install, durable queue-state, and explicitly scoped implementation boundaries before they trust the suite.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Replace only the stale existing-project Q&A answer with a compact statement of the reviewed/confirmed managed install paths, durable `do-work/` queue state, and REQ-scoped project-source write boundary. Leave the already-accurate Installation section unchanged; verify the old absolute claim is absent and the replacement agrees with the Builder Scope Fence and ADR-019.
- [x] **[APPLY]:** Replaced the stale Q&A answer in `README.md` only; preserved the accurate Installation section and made no runtime or boundary changes.
- [x] **[UNIFY]:** Reviewed the complete `README.md` diff (`1 insertion, 1 deletion`), ran `git diff --check`, confirmed the exact stale sentence is absent, and verified all three boundary statements against the Installation section, Builder Scope Fence, and ADR-019. No README-specific native linter or debug artifact applies to this prose-only change.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names `README.md`, supplies the contradicted sentence, and defines the three exact boundaries the replacement must communicate. The installation section already reflects the four-skill suite, so the remaining change is one focused public-prose correction.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `README.md` (modified)

**What was done:** Replaced the false `do-work/`-only claim with one compact answer that separates confirmed managed install/update paths, durable queue state, and explicitly invoked REQ-scoped project-source writes. The already-correct four-skill Installation section was preserved.

## Qualification

Passed — 1 implementation file verified, all 4 Detailed Requirements are satisfied in the current README, the one-line diff is substantive and scoped, P-A-U is complete, and data-flow checks are not applicable to prose.

## Testing

**Tests run:** exact stale-sentence absence check; three-boundary presence check; comparison with the Builder Scope Fence and ADR-019; `git diff --check`
**Result:** ✓ All passing

No runnable test is justified for this public-prose correction; the captured RED/GREEN evidence is the contract comparison itself.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-15T09:48:29Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the public README now states all three real write boundaries and no longer makes the contradicted absolute claim.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed inline by review-work action after the delegated reviewer timed out*

## Orientation

Existing-project adopters now see the suite's actual safety boundary before installation: confirmed managed suite/configuration paths, durable `do-work/` queue state, and explicitly invoked REQ-scoped source writes.

## Knowledge Handoff

- `kb_status: skipped`
- Route A prose correction had no non-trivial implementation lesson to promote.
