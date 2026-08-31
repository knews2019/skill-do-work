---
id: REQ-435
title: 'Review fix: Complete the doctor-forensics delegation contract'
status: completed
domain: general
created_at: 2026-08-30T22:01:41Z
claimed_at: 2026-08-31T16:23:46Z
completed_at: 2026-08-31T17:41:35Z
commit: c1536cbf
route: B
user_request: UR-081
addendum_to: REQ-410
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-31T16:36:00Z
  basis:
    - Route B
    - 8-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
    - full-suite verification
tdd: true
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Complete the Doctor-Forensics Delegation Contract

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Delete unearned state totals, map every required report field to typed doctor evidence, route takeover judgment to Crash Recovery, and replace all live deleted-check pointers inside the frozen eight-file scope.
- [x] **[APPLY]:** Added the mixed-state authority/reference ratchet first, then repaired the action and six consumer groups without changing doctor/result production schemas.
- [x] **[UNIFY]:** Reviewed all eight changed files and ran focused/full CLI, queue-board, vet, reference, contract, exact-Go-1.25, ShellCheck, and scope checks; no debug artifacts remain.

## What
Make the natural-language forensics action executable end-to-end from its declared authorities, including report counts, actionable remedies, and valid cross-references.

## Context
Found during terminal re-review of REQ-410 after its single remediation. The action forbids deterministic rescanning and says to report only doctor evidence, while its output contract still requires queue/archive/working counts and remediation details absent from the typed doctor result. It also removed the former stuck-work reset procedure while `work-reference.md` and other consumers still point to numbered checks that no longer exist. An agent must currently rescan, omit required output, or follow stale pointers. Fold-first scan found no pending REQ or sweep in any UR that shares this delegation-contract root cause.

## Requirements
- Choose and document one complete ownership contract for queue/archive/working counts and every required deterministic report field without ad hoc rescanning.
- Ensure emitted doctor findings or the action's remaining judgment steps provide actionable remedies for the required report, including stuck-work handling.
- Replace stale numbered-check references in `work-reference.md` and every other in-scope consumer with stable canonical anchors or commands.
- Keep recurring-correction judgment and board-owned release verification outside doctor without duplicating their mechanics.
- Add contract coverage proving the action can produce its documented report using only its declared authorities.

## Constraints
- Delete unused queue/archive/working count requirements by default. Add typed counts only when a concrete consumer and regression test justify the lasting schema surface.

## Red-Green Proof
**RED prompt/case:** Execute the forensics action contract against a fixture with stuck work and mixed queue/archive/working states; assert every required report field and remedy has an authoritative source and every referenced anchor resolves.
**Why RED now:** Doctor's current result omits the required state counts and manual-reset procedure, while remaining documentation still references deleted numbered checks.
**GREEN when:** The action produces the documented report without independent mechanical scans, all remedies advance the user, and a reference contract finds no stale Check-number consumers.
**Validation:** Terminal review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** The default constraint narrows the ownership choice, but the action, reference consumers, doctor result, and contract tests must be explored before freezing the exact four-file cross-subsystem scope.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

No concrete consumer justifies adding queue/archive/working totals to the typed doctor result, so the default constraint applies: remove those unused report requirements. Doctor findings already carry every deterministic evidence and remediation field the action needs.

The canonical authority split is doctor for mechanical findings, `actions/work-reference.md` → **Crash Recovery (Step 1)** for stuck-work takeover judgment, the retained Recurring Corrections phase for nondeterministic grouping, and queue-kanban verify for release/queue invariants. A complete stale-reference repair requires eight current shipped files rather than the earlier four-file estimate.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/forensics.md` (modify) — delete unused state totals, map typed finding fields to report fields, and route stuck-work judgment to the stable Crash Recovery anchor.
- `skills/do-work/actions/work-reference.md` (modify) — replace deleted Check-number pointers with stable local anchors, timestamp authority, or doctor finding codes.
- `skills/do-work/actions/abandon.md` (modify) — replace the deleted failed-classification Check pointer with the stable doctor finding code.
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — replace deleted timestamp Check comments with the stable doctor authority.
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — repair user-visible recovery remedies and deleted-check comments.
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — replace the deleted timestamp-check comment with stable doctor authority.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modify) — repair the active maintenance instruction naming the deleted check.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go` (modify) — add the mixed-state authority and current-consumer reference ratchet.

**Files I will NOT touch:** doctor/result production Go schemas, the forensics guide, historical changelogs, archived REQs, run handbacks, or any queue/release metadata.

**Acceptance criteria (restated from REQ):**
- [ ] The forensics report is fully sourced from typed doctor findings plus its three declared judgment/board authorities, with no independent mechanical rescan.
- [ ] Unused queue/archive/working totals are absent from both report templates and no unearned result-schema fields are added.
- [ ] Stuck-work findings retain actionable doctor evidence and route takeover judgment to Crash Recovery without duplicating mechanics.
- [ ] Every live pointer to a deleted mechanical Check number in the frozen consumers is replaced by a stable anchor, command, or finding code.
- [ ] Recurring-correction judgment and board-owned verification remain outside doctor without duplicated mechanics.
- [ ] A mixed-state contract regression proves all required report fields, remedies, and reference destinations are authoritative.

## Implementation Summary

**What was done:** Forensics now assembles its report from typed doctor findings plus the stable Crash Recovery, Recurring Corrections, and queue-kanban verification authorities. Every board verify finding maps deterministically to Warnings and contributes one warning, while skipped/not-applicable probes remain coverage-only. Unsupported repository-state totals were deleted, and every live pointer to removed mechanical checks in the frozen consumers now resolves to a stable anchor, command, heading, or finding code.

**Files changed:**
- `skills/do-work/actions/forensics.md` (modified) — defines typed finding-to-report mapping, removes unused totals, and routes stuck-work judgment to Crash Recovery.
- `skills/do-work/actions/work-reference.md` (modified) — replaces deleted Check-number pointers with stable local anchors and doctor finding codes.
- `skills/do-work/actions/abandon.md` (modified) — uses the stable failed-classification finding code.
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified) — replaces deleted timestamp-check comments with stable doctor authority.
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — repairs recovery remedies and stranded-terminal references.
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — names doctor as the independent timestamp-evidence authority.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modified) — removes the stale third-copy maintenance instruction.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go` (modified) — adds mixed-state report-authority and six-consumer reference ratchets.

**Integration range:** `0acc9342..c1536cbf`

*Generated by work action from the builder hand-back*

## Initial Review

**Overall: 50%** | 2026-08-31T17:14:11Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 90% |
| Test Adequacy | 75% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Board-owned findings expose no authoritative severity or canonical category-to-severity mapping, so the promised severity sections and totals still require ad hoc judgment — impact-user-visible → returned for the one allowed remediation attempt.

**Minor findings:** 0 (report only)
**Acceptance:** Fail — doctor authority, total deletion, and reference repair pass, but the board-finding report path lacks a severity source and the contract test misses that surface.
**Suggested testing:** 2 items
**Follow-ups created:** None pending remediation; **sweeps appended to:** None

*Reviewed by review-work action*

## Remediation

**Initial review:** Acceptance failed at 50% because board-owned findings had no authoritative severity source or canonical mapping into the required report sections and totals.

**One allowed attempt:** Added an action-owned output-class rule: every queue-kanban verify finding is a Warning and contributes +1 warning; skipped and not-applicable probes appear only under skipped/unverified coverage, and fixability remains remedy metadata. A contract regression using the real `version-changelog-mismatch` category failed on missing placement and contribution before the action change, then passed after it.

**Remediation merge:** `c1536cbf` (branch commit `4f1568f4`), changing only `skills/do-work/actions/forensics.md` and `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go` within the frozen eight-file scope.

## Qualification

Passed — all eight declared shipped files are substantive and present in `0acc9342..c1536cbf`. The final action contract assigns deterministic report placement and counting to every doctor and board finding while leaving skipped/not-applicable probes as coverage-only evidence. P-A-U, stale-reference, debug-artifact, and exact-scope checks passed.

## Testing

**Red-green validation:** The mixed-state authority regression first failed because unsupported totals and deleted Check-number references had no authoritative source. The remediation regression then failed against the pre-remediation action because a real `version-changelog-mismatch` board finding had no severity placement or contribution. Both pass against the final tree.

**Merged-state checks:**
- Focused and full doctor CLI tests — PASS.
- Full queue-kanban board tests — PASS.
- `go vet ./...` in both Go modules — PASS.
- Shipped reference and contract regressions — PASS.
- `bash _dev/tests/do-work-cli-go125-compatibility.sh` — PASS (exact Go 1.25.0).
- Warning-level ShellCheck for `repair-req-timestamps.sh` — PASS.
- `DO_WORK_DIFF_RANGE="0acc9342..c1536cbf" skills/do-work/tools/checks/qualify.sh ...` — PASS.

## Review

**Overall: 98%** | 2026-08-31T17:38:09Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — the remediation closes the original finding with an authoritative, condition-based mapping for every board finding and coverage-only handling for skipped/not-applicable probes.
**Suggested testing:** None beyond the canonical repository gate.
**Follow-ups created:** None; **sweeps appended to:** None.

*Re-reviewed by review-work action after the one allowed remediation attempt*

## Lessons Learned

**What worked:** Keeping doctor, recurring judgment, crash recovery, and board verification as distinct authorities made the forensics report complete without growing an unused production schema.
**What didn't:** Delegating board verification without defining how its findings enter severity sections left the first contract incomplete even though the board output itself was authoritative.
**Worth knowing:** A delegated report source needs both evidence ownership and an explicit projection into the consumer's classification and totals; skipped coverage must be defined separately from findings.

## Orientation

Forensics now builds its deterministic report from typed doctor findings and queue-kanban verification, with explicit Warning mapping for board findings. Crash Recovery owns stuck-work takeover judgment, Recurring Corrections owns nondeterministic grouping, and unsupported queue/archive/working totals are no longer requested.
