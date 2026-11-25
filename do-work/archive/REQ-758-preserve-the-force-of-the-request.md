---
id: REQ-758
title: 'Preserve the force of the request during capture'
status: completed
dispatch_at: 2026-11-25T23:12:35Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-25T23:12:22Z
  basis:
  - trivial short-circuit
route: A
created_at: 2026-11-25T22:57:02Z
user_request: UR-190
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-mechanical
batch: skill-instruction-maintenance
related: [REQ-759, REQ-760]
required_lessons: ['_dev/primes/lessons-action-files.md#condition-preserving-prose-extraction']
write_set: [skills/do-work/actions/capture.md]
claimed_at: 2026-11-25T23:11:32Z
builder_handback_at: 2026-11-25T23:15:48Z
integration_at: 2026-11-25T23:15:48Z
review_at: 2026-11-25T23:32:41Z
commit: 5e742d16399dd3a6e232cb0317fd16a891a4d004
heavy_verified_at: 2026-11-25T23:34:07Z
heavy_verified_revision: 77c78cd721c24f0c49ed4f1901cd656a60575060
completed_at: 2026-11-25T23:34:08Z
release_at: 2026-11-25T23:34:08Z
---
# Preserve the force of the request during capture
## What
Delete “Mentioned once in passing: Still a requirement. Capture it.” Trim the closing capture checklists only where they repeat the pre-publication comparison.

## Why
An example can help the builder without becoming an acceptance criterion. A firm constraint can appear only once. Keep the verbatim request and the final comparison of its outcomes, constraints, supplied checks, and builder latitude with each prepared REQ or fold destination. Preserve distinct publication and linkage checks.

## Red-Green Proof
**RED prompt/case:** Replay the Durations request in `do-work/archive/UR-069/input.md:43`. The overflow count is required and the number of visible rows is the builder’s call. Today the blanket rule can turn the example row count into a fixed requirement.
**GREEN when:** The exact blanket rule in `skills/do-work/actions/capture.md:244` is gone. The walkthrough preserves both decisions with their original force, and the final reread at line 210 accounts for them once.
**Validation:** Inferred from the accepted finding. UR-069 illustrates the distinction. Its capture preserved it.

## Full Context
See `do-work/user-requests/UR-190/input.md` for the verbatim capture instruction.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A. Reviewed the action-file prime and the required condition-preserving prose lesson. Remove the blanket rule and narrow the overlapping checklist item.
- [x] **[APPLY]:** Merged the builder's change to the single declared file, `skills/do-work/actions/capture.md`.
- [x] **[UNIFY]:** Reviewed the full merge diff and its stat. Only `capture.md` changed, Step 5's reread and publication path are intact, and the range passes `git diff --check` with no debug artifacts.
*Source: Finding 1, pasted audit, rank 1, no severity supplied. Accepted at efcd3618. Surface-cost: N/A. Delete the passing-mention rule. Trim the closing capture checklists where they repeat the comparison already performed before writing the queue.*

## Triage

Route A. The capture action names the exact sentence and identifies which closing checks are duplicates.

## Plan

**Planning not required** — Route A, direct implementation.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified). No test files changed.

**What was done:** Removed the instruction that made passing mentions mandatory. The final reread still compares requested outcomes, constraints, supplied checks, and builder discretion against the prepared REQs and folds. The closing checklist now checks that each folded item names its destination.

**Implementation range:** `514d330b984d18f66297d1a27ff140cf14175101..5e742d16399dd3a6e232cb0317fd16a891a4d004`. The merge's second parent is the handed-back commit `dc6c016bb29a9160cae34637b29ddcfce017b96f`.

## Decisions

- **D-01.** Narrowed the mixed checklist item to the fold linkage check. This preserves the link from each folded request to its recorded destination, including a prose backlog, while Step 5 owns the comparison to the original input.

## Qualification

Passed. Request-bound `advance` qualification accepted the merged range with no findings. The range changes only the declared capture action, matches the builder's diff, and retains the final reread, verbatim input, publication, and linkage checks.

## Testing

**Tests run:** Read-only UR-069 replay, `bash _dev/tests/shipped-package-reference-contract.sh` via the request-bound focused gate, and `bash _dev/tests/maintainer-verify.sh` on the merged tree.

**Result:** The replay found the blanket rule at the baseline and confirmed its removal after the merge. Against the archived UR-069 clarification and REQ-346, the walkthrough requires an explicit count for every UR without a row and leaves the row count to the builder. The final reread and verbatim-input, publication, and linkage checks passed. The focused gate used a one-line Bash invocation because the probe runner evaluates probe text through `sh`. The standard maintainer gate passed.

**Repository gate runs:** First exit 1 in the queue-kanban test stage. The immediate retry exited 1 when the fast-stage evidence probe exceeded its 30-second file budget. The fast-stage probe passed alone. The isolated board stage reproduced a timeout in the nested browser lane of `TestLegacyStrictEntryPointsRejectZeroProbes`, which enables that lane internally even with `QUEUE_KANBAN_BROWSER_PROBES=off` in the parent. The single regression test passed without changes in 26.68 seconds under the same fast flags. The next full run used the unchanged canonical argv.

**New tests added:** None. The replay exercised the merged instructions through a walkthrough with assertions, without writing live capture records.

**Heavy verification plan:**
- Range `514d330b984d18f66297d1a27ff140cf14175101..5e742d16399dd3a6e232cb0317fd16a891a4d004`.
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. The capture action is under the shipped `skills` tree.

*Verified by work action*

## Review

**Overall: 97%** | 2026-11-25T23:32:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:**
- Verification's `Nit` description still uses “passing mention” in the action and user guide. Agents could read these descriptions as classifying a clause by how often it appears. Both files are within REQ-759's declared verification scope — impact-rule-change → report only

**Minor findings:** The handed-back commit subject is missing the bracketed REQ marker. Its SHA is retained as the merge's second parent, and the merge itself carries `[REQ-758]` — impact-rule-change → report only

**Acceptance:** Pass — implementation and integration. The merged change matches the UR-190 finding. The UR-069 walkthrough preserves the mandatory remainder and the builder's choice of row count. The retained checks cover verbatim input, REQ linkage, fold destinations, and canonical publication. The focused and standard maintainer gates passed.

**Restatement sweep:** redefined how capture treats passing mentions and what the closing fold check verifies. Searched capture, verification, their guides, and contract tests. The verifier wording is reported above. REQ-757 and REQ-759 were still in `working/`, so the wave-end sweep did not apply.

**Suggested testing:** 0 items

**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Orientation

Capture preserves firm constraints and builder discretion in its final comparison, with a separate checklist item for fold destinations.

## Heavy Verification Plan

- Base revision `514d330b984d18f66297d1a27ff140cf14175101`
- Target revision `5e742d16399dd3a6e232cb0317fd16a891a4d004`
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. The changed capture action matches the `skills` subtree.

## Heavy Verification Result

- Target revision `5e742d16399dd3a6e232cb0317fd16a891a4d004`
- Execution revision `77c78cd721c24f0c49ed4f1901cd656a60575060`
- `staged-skills` — executed, exit 0, 44 seconds, no skip. Command `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`.

## Timing

Observed 2026-11-25T23:19:53Z to 2026-11-25T23:32:15Z: 12m 22s total, 7m 05s attributed across 3 events, 5m 17s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 7m 05s | 3 |

Slowest command: verification-gate / Canonical maintainer verification after isolated checks, 3m 16s, exit 0, bash (2 argv tokens).
