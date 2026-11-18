---
id: REQ-745
title: 'Save likes, dislikes and notes on reports and figures'
status: pending
created_at: 2026-11-18T15:50:47Z
user_request: UR-184
domain: frontend
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-feedback
related: [REQ-746, REQ-747, REQ-748]
depends_on: [REQ-744]
---
# Save likes, dislikes and notes on reports and figures

## What
Open the gallery with `just run-report-gallery`. From that page mark a whole report or an individual figure as liked or disliked. A short note is optional.

## Detailed Requirements
Keep report and figure feedback distinguishable, including when a user likes a report and dislikes one figure in it. Persist marks so both the reopened page and the round-generation agent can read the same feedback.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** Open with the Just recipe, like a report without a note, dislike one figure with a note, close and reopen. Both distinct targets, their marks and the note are restored and machine-readable.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers gallery routing and report writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.
