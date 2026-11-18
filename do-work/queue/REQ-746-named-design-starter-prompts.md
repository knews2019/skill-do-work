---
id: REQ-746
title: 'Copy a starter prompt that selects a specific report design'
status: pending
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-feedback
related: [REQ-745, REQ-747, REQ-748]
depends_on: [REQ-744]
---
# Copy a starter prompt that selects a specific report design

## What
Give every design a copyable do-work-toolbox ai-report starter prompt. Pasting it makes the writer use the named design directly.

## Detailed Requirements
Cover each shipped sample, each user-added sample and each generated report design. Show the full usable prompt and a copy control on the gallery page. An explicit selection wins even when automatic ranking would choose another sample. Preserve automatic ranking for ordinary requests.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** For each of the three design sources, copy the starter prompt and use it for a target whose automatic match would differ. The writer uses the named design with the new target’s own evidence.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers gallery routing and report writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.
