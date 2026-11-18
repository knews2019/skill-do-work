---
id: REQ-742
title: 'Offer tailored second-opinion prompts for the finished report'
status: pending
created_at: 2026-11-18T15:04:36Z
user_request: UR-182
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-sample-gallery
related: [REQ-740, REQ-741]
depends_on: [REQ-741]
addendum_to: REQ-712
---
# Offer tailored second-opinion prompts for the finished report

## What
Once the detailed report is built, offer ready-to-use prompts for distinct second opinions on that specific report. Choose the best available preset for what each variation helps the reader see.

## Why
The reader should be able to understand the main finding in most reports.

## Detailed Requirements
Use the completed report, the gallery match, and the available presets to choose meaningful alternative views. Point each prompt at this report and explain briefly what the alternative clarifies. The deliverable is the suggested prompts; running the variations is a later user action.

## Constraints
The screenshots illustrate reading patterns. They supply design references, and each generated report supplies its own facts.

## Red-Green Proof
**RED prompt/case:** Inspect the finished comparison report and try its suggested lines as requests for alternate views.
**Why RED now:** The existing close offers a different layout and a teaching method, but does not choose a fitting preset for each variation in light of the gallery-informed report.
**GREEN when:** The prompts target the current report and name usable, fitting presets for distinct alternate views of its evidence.
**Validation:** User confirmed the end-to-end inflation example during capture.

## Prior Implementation
REQ-712 (choose the best teaching method for the report, commit 0c06eb10) added the teaching-method suggestion beside the layout suggestion in the report skeleton and design rules. Extend those recommendations using the available second-opinion presets and the current report.

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
