---
id: REQ-748
title: 'Run two rounds of five reports on one real evidence bundle'
status: claimed
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-gallery-feedback
related: [REQ-745, REQ-746, REQ-747]
depends_on: [REQ-747]
claimed_at: 2026-11-18T23:53:11Z
---
# Run two rounds of five reports on one real evidence bundle

## What
Carry out the first feedback trial with ten reports, five per round. The builder chooses a suitable existing real evidence bundle and records which one was used.

## Detailed Requirements
Keep the same bundle across both rounds. Deliver the first five for the user to mark, then pause. Generate the second five only after the user has marked the first five. Use the saved feedback and record lineage and changes. Check facts and units against the bundle, and disqualify an erroneous report even if the user liked it.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** Observe five first-round reports and no second-round generation before the user’s marks. After genuine user feedback, observe five second-round reports tied to the same evidence, with traceable changes and correctness decisions.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers gallery routing and report writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.
