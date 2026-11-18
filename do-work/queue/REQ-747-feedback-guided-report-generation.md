---
id: REQ-747
title: 'Generate traceable report variations with GPT-6.2-astra in T3 Code'
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
related: [REQ-745, REQ-746, REQ-748]
depends_on: [REQ-745, REQ-746]
---
# Generate traceable report variations with GPT-6.2-astra in T3 Code

## What
Provide an agent command that reads saved feedback and generates the next round using GPT-6.2-astra through T3 Code.

## Detailed Requirements
Each generated report records its parent report, the specific feedback marks used, and the changes it made. Keep report and figure marks attributable to their original targets. For the initial generation record the real source and any absence of prior marks honestly. Every generated design gets the same copyable named-design starter prompt. Factual or unit errors disqualify a report regardless of design approval.

## Constraints
This batch builds on REQ-744 (the gallery-aware report writer). The builder chooses and records the existing real evidence bundle for the first trial. The second round waits for the user's actual first-round feedback.

## Red-Green Proof
**RED prompt/case:** Try the workflow below with the original gallery and writer.
**Why RED now:** The captured gallery and ranking features provide sample selection and second opinions. They do not implement this requested behavior.
**GREEN when:** With known marks, generate a variation through the named model and provider. Inspect its parent, marks and change record and try its starter prompt. A liked example with an incorrect unit is disqualified.
**Validation:** The user specified the observable behavior and confirmed who chooses the evidence bundle. The concrete proof is inferred during capture.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers gallery routing and report writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the linked request and evidence, then record the approach.
- [ ] **[APPLY]:** Implement or execute the scoped work.
- [ ] **[UNIFY]:** Verify the outcome and the evidence for the acceptance case.

## Full Context
See `do-work/user-requests/UR-184/input.md` for the verbatim request and clarification.
