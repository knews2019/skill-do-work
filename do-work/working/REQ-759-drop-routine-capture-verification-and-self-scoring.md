---
id: REQ-759
title: 'Drop routine capture verification and self-scoring'
status: claimed
created_at: 2026-11-25T22:57:02Z
user_request: UR-190
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
batch: skill-instruction-maintenance
related: [REQ-758, REQ-760]
required_lessons: ['_dev/primes/lessons-action-files.md#cross-action-exception-closure']
write_set: [skills/do-work/next-steps.md, skills/do-work/actions/capture.md, skills/do-work/actions/verify-requests.md, skills/do-work/docs/verify-requests-guide.md]
claimed_at: 2026-11-25T23:11:32Z
---
# Drop routine capture verification and self-scoring
## What
Remove the default verification suggestion for a new capture and the extra nudge for long or complex captures. An explicit verification reports the original clause, which REQ changed or missed it, and the correction or remaining ambiguity. Delete the confidence percentages, score bands, averaging, and post-repair re-score. Align the user guide with that output.

## Why
What matters to the maintainer is how often they have to restore an already-stated requirement. A discrepancy the verifier fixes itself, a user correction, a question already answered in the request, and a later change of preference have different meanings. This change concerns verification and its report. Building a metrics system is outside its scope.

## Constraints
Keep explicit verification and repair, comparison with folded destinations, and the gap severities. Preserve the verification that the user opts into through `validate-feedback --capture` and its gap-based run decision at `skills/do-work-toolbox/actions/validate-feedback.md:121`. Report corrections as applied only when a repair actually occurred.

## Red-Green Proof
**RED prompt/case:** Follow a new capture through `skills/do-work/next-steps.md:27` and `skills/do-work/actions/capture.md:223`. Then replay `do-work/archive/UR-011/input.md:75` against the verification report at `skills/do-work/actions/verify-requests.md:118` and the re-score at line 171.
**GREEN when:** Ordinary capture can go straight to build on the user's instruction. Explicit verification identifies that the capture described a hard stop while the original said the agent warned and recovered. The report shows the clause and disposition without a numerical score. UR-065, by contrast, records a confirmed ambiguity about blocked work in the Timeline forecast (`do-work/archive/UR-065/input.md:108`).
**Validation:** Inferred from the accepted finding and its archived cases.

## Full Context
See `do-work/user-requests/UR-190/input.md` for the verbatim capture instruction.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Finding 2, pasted audit, rank 2, no severity supplied. Accepted at efcd3618. Surface-cost: N/A. Remove the suggestion to run `verify-requests` after every new capture. When verification is requested, report the original clause, what the capture changed or missed, and the correction. Delete the percentage table and re-scoring step.*
