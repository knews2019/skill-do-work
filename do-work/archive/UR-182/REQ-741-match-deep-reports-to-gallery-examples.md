---
id: REQ-741
title: 'Choose a gallery example before drafting each deep report'
status: cancelled
created_at: 2026-11-18T15:04:36Z
user_request: UR-182
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-sample-gallery
related: [REQ-740, REQ-742]
depends_on: [REQ-740]
completed_at: 2026-11-18T15:28:06Z
---
# Choose a gallery example before drafting each deep report

## What
Before composing a detailed ai-report, consider the available gallery samples, including additions, and choose the closest fit for the reader’s question and the available evidence. Adapt the useful parts using editorial judgment.

## Why
The reader should be able to understand the main finding in most reports.

## Detailed Requirements
Apply the chosen approach to the report itself. The author may combine or simplify ideas when that helps the reader. Keep the report grounded in its own evidence. For the agreed example, show the two easily confused measures side by side and label what each measures.

## Constraints
The screenshots illustrate reading patterns. They supply design references, and each generated report supplies its own facts.

## Red-Green Proof
**RED prompt/case:** Save the inflation sample and request a detailed report where two measures could be confused.
**Why RED now:** The current report action starts from the kit and prescribes a narrative without consulting a gallery of saved examples.
**GREEN when:** The writer consults the gallery before drafting and uses the comparison to clarify the two measures with the new subject’s evidence.
**Validation:** User confirmed the end-to-end inflation example during capture.

## Builder Guidance
Use the existing report gallery demonstration, report action and second-opinion workflow as starting points. The choice of storage, matching technique and page design belongs to the builder.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers report routing and alternate writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read relevant rules and record the approach.
- [ ] **[APPLY]:** Implement the agreed behavior.
- [ ] **[UNIFY]:** Verify the diff and the acceptance example.

## Full Context
See `do-work/user-requests/UR-182/input.md` for the original request and the confirmed example.

## Cancelled

- **When:** 2026-11-18T15:28:06Z
- **Why:** Replaced on 2026-11-18 by a two-request rewrite the user chose: one request for the report gallery and one for the report writer that uses it.
- **Decided by:** user, via `do-work abandon`
