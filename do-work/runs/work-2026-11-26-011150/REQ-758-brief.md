# REQ-758 builder brief

Implement the request below in your T3-assigned checkout and branch. Record their absolute path and name in the handback. Read CLAUDE.md, the applicable primes and the implementation crew rules in skills/do-work/actions/work-implementation-reference.md, including required lesson entries. Consult do-work/lessons-index.md and the linked user request. Preserve sibling changes. Only this request is yours. Use the REQ's own validation, capture the baseline before edits, follow red/green where specified, run focused checks, review the complete diff, and commit the implementation with REQ-758 in the subject. Do not publish, merge or modify queue records, shared version files or CHANGELOG.

Write the completed handback to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-758-handback.md in the main checkout. Provide operative branch, SHA, base, file manifest, changed behavior, focused commands and results, baseline/red evidence, decisions and discovered tasks (explicitly say when none), integration seams and any limitations. Append UTC start/end lines to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-758-progress.log for long operations. These two main-tree paths are the only shared writes. No other agent owns your implementation. Full gate and release are serial integration work. No project run-policy file.

---
id: REQ-758
title: 'Preserve the force of the request during capture'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Finding 1, pasted audit, rank 1, no severity supplied. Accepted at efcd3618. Surface-cost: N/A. Delete the passing-mention rule. Trim the closing capture checklists where they repeat the comparison already performed before writing the queue.*

## Triage

Route A. The capture action names the exact sentence and identifies which closing checks are duplicates.

## Plan

**Planning not required** — Route A, direct implementation.
