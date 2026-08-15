---
id: REQ-186
title: Required baseline verification executes two child suites twice
status: pending
created_at: 2026-08-15T07:13:20Z
user_request: UR-041
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec: refactor
depends_on: []
maintenance: true
effort_estimate: trivial
related: [REQ-181, REQ-182, REQ-183, REQ-184, REQ-185, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [CLAUDE.md, _dev/tests/contract-regressions.sh, _dev/tests/staged-skills-contract.sh]
---

# Required Baseline Verification Executes Two Child Suites Twice

## What

Give each duplicated baseline child suite one owner: remove the aggregate's redundant direct prescribed-shell edge and stop requiring a second standalone shipped-reference invocation when it has no distinct mode or fixture.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`prescribed-shell-scripts-behavior.sh` runs directly and through `staged-skills-contract.sh` inside the aggregate. `shipped-package-reference-contract.sh` runs inside the aggregate and again under the documented hand-back baseline without a different mode or fixture.

## Context

- Audit priority: P3; impact 2; effort trivial.
- Root-cause key: `baseline-suite-single-ownership`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 6.
- Reproduce: `rg -n 'prescribed-shell-scripts-behavior|staged-skills-contract|shipped-package-reference-contract' CLAUDE.md _dev/tests/contract-regressions.sh _dev/tests/staged-skills-contract.sh`.

## Detailed Requirements

- Delete `_dev/tests/contract-regressions.sh`'s direct invocation of `prescribed-shell-scripts-behavior.sh`.
- Retain `_dev/tests/staged-skills-contract.sh` as the owner of the prescribed-shell behavior suite for standalone-semantics coverage.
- Keep `shipped-package-reference-contract.sh` inside the aggregate.
- Narrow the documented hand-back baseline so standalone child invocations are focused and not a second mandatory execution of an identical child.
- Preserve every distinct late aggregate check and final pass marker.

## Constraints

- Delete before adding; this is a duplicate-edge removal, not a harness redesign.
- Do not add a generic test-graph registry or parser.
- Do not absorb the separate REQ-180 filename-casing defect or the already-covered structural harness-bloat review.

## Dependencies

None. Coordinate the shared `CLAUDE.md` surface with REQ-187 if both are implemented in one batch.

## Builder Guidance

Firm intent. History identifies the correct ownership edges, so no new lock-in framework is earned; verify the observable execution markers directly.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Trace the three child-suite names across `CLAUDE.md`, the aggregate, and the staged suite.
**Why RED now:** Two child suites each have two required ownership edges with no distinct mode or fixture.
**GREEN when:** The aggregate executes each required behavior once through its intended owner, the hand-back instructions do not re-run an identical child, and all distinct late checks still execute and pass.
**Validation:** Confirmed by the user during verification on 2026-08-15. TDD is false because the finding explicitly rejects new generic graph machinery for two direct deletions.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 06, labeled P3, impact 2, trivial effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 6 in the canonical audit.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*
