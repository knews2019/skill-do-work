---
id: REQ-186
title: Required baseline verification executes two child suites twice
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T11:25:40Z
completed_at: 2026-08-15T11:35:07Z
commit: 0ab2b79
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
route: A
kb_status: promoted
kb_entry: REQ-186-required-baseline-verification-executes-.md
---

# Required Baseline Verification Executes Two Child Suites Twice

## What

Give each duplicated baseline child suite one owner: remove the aggregate's redundant direct prescribed-shell edge and stop requiring a second standalone shipped-reference invocation when it has no distinct mode or fixture.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Delete the aggregate's redundant prescribed-shell edge and the maintainer baseline's redundant standalone shipped-reference edge; preserve the staged owner, aggregate-owned shipped-reference check, late checks, and final pass marker without adding graph machinery.
- [x] **[APPLY]:** Removed the aggregate's direct prescribed-shell invocation and narrowed the documented hand-back baseline to the aggregate. Left the staged suite unchanged as the prescribed-shell owner.
- [x] **[UNIFY]:** Reviewed the two-file deletion diff and the three-name ownership trace. Verified the direct staged suite, full aggregate through its late installer and final marker, shell syntax, ShellCheck-enabled aggregate lint, and `git diff --check`; no debug artifacts or additions remain.

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

## Triage

**Route A** — this is a three-site ownership correction with an exact deletion-first repair: remove the aggregate's direct prescribed-shell edge, keep the staged-suite owner unchanged, and remove the redundant standalone shipped-reference requirement from the maintainer baseline.

## Plan

1. Delete the direct `prescribed-shell-scripts-behavior.sh` invocation and failure branch at the aggregate's start; retain the staged suite's existing invocation unchanged.
2. Narrow `CLAUDE.md`'s integrating baseline to the aggregate, which already owns `shipped-package-reference-contract.sh`; preserve the instruction to run focused suites when their distinct coverage is touched.
3. Trace the three ownership names after the edit, run the aggregate and staged suites directly, verify the late aggregate pass marker, and review the three-file boundary for accidental additions.

## Implementation Summary

**Files changed:**
- `CLAUDE.md` (modified) — makes the aggregate the single required hand-back baseline instead of mandating a second identical shipped-reference child run
- `_dev/tests/contract-regressions.sh` (modified) — removes the direct prescribed-shell child edge so the staged suite is its sole owner

**Verified unchanged:** `_dev/tests/staged-skills-contract.sh` still invokes `prescribed-shell-scripts-behavior.sh` for standalone-semantics coverage.

**Behavior:** The required aggregate executes prescribed-shell behavior once through the staged contract and shipped-reference once through its own aggregate edge. Distinct late checks and the final aggregate pass marker remain intact.

## Testing

- Ownership trace across `CLAUDE.md`, the aggregate, and staged suite — PASS; prescribed-shell has one execution owner and shipped-reference remains aggregate-owned
- `bash _dev/tests/staged-skills-contract.sh` — PASS, including 22 prescribed-shell cases
- `bash _dev/tests/contract-regressions.sh` — PASS, including shipped-reference, staged-skills, suite-installer, and `Contract regression checks passed.`
- `bash -n` on the aggregate, staged, prescribed-shell, and shipped-reference scripts — PASS
- Aggregate shell-block lint with ShellCheck enabled — PASS
- `git diff --check` — PASS

## Qualification

- **Boundary:** PASS — Route A uses the captured write set; only `CLAUDE.md` and `_dev/tests/contract-regressions.sh` changed, while the declared staged owner remained intentionally byte-unchanged.
- **Mechanical checks:** PASS — `qualify.sh` found both implementation files, complete P-A-U evidence, and no debug artifacts; `scope-drift.sh` correctly skips Route A.
- **Substance and traceability:** PASS — the two duplicate edges named by the finding are deleted, with no new registry, parser, or ownership prose added.
- **Wiring/data flow:** PASS — the aggregate still reaches shipped-reference directly and prescribed-shell through staged-skills; the installer and final pass marker execute afterward.

## Review

**Result:** Approve — Acceptance: Pass
**Overall score:** 100%

- **Requirements (100%):** Both duplicate ownership edges are removed and every retained owner/late check satisfies the stated contract.
- **Code quality (100%):** The repair is pure subtraction with no replacement registry, parser, or prose machinery.
- **Test adequacy (100%):** Direct staged and aggregate runs prove the child, failure-propagation path, late installer, and final pass marker remain reachable.
- **Scope (100%):** Only the two files requiring deletion changed; the surviving staged owner remained byte-unchanged.

**Important findings:** None.
**Minor findings:** None.
**Explicit remediation:** None.

## Lessons Learned

- A required aggregate should give each identical child invocation one owner. If a nested suite already preserves the needed standalone semantics and failure propagation, a second direct edge adds runtime without adding evidence.
- Maintainer hand-back instructions should name the aggregate baseline once and reserve standalone child runs for genuinely distinct modes, fixtures, or touched-area focus.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.
