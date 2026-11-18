---
id: REQ-752
title: 'Compare the frozen endpoints and publish the gallery climb evidence'
status: pending
created_at: 2026-11-18T17:43:03Z
user_request: UR-185
domain: general
prime_files: []
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-gallery-climb
related: [REQ-749, REQ-750, REQ-751]
depends_on: [REQ-751]
---
# Compare the frozen endpoints and publish the gallery climb evidence

## What
Report what the gallery climb achieved using comparable start and end scores from one fresh blind evaluator.

## Detailed Requirements
Give one verdict line with start score, end score, rounds run, changes kept,
and changes reverted. For comparable start and end scores, have one fresh
Opus 5.5 judge score the frozen baseline and final incumbent together, blind.
Show the original, baseline, and final screenshots. Include a table of all ten
names, kinds, and report paths. Link the climb log and gate evidence. List
open decisions as D1, D2, and so on.

Use a new Opus 5.5 evaluator through T3 Code. Present the frozen baseline and final incumbent together under the randomized screenshot-first procedure and exact R1–R5 rubric from REQ-750. Include all ten cards in matched 1440/390 px light/dark views. Record first choices for each version before releasing either HTML; keep the key, mapping, chronology, earlier scores, change descriptions and builder notes private. Retain both five-line scores, totals, located weakest points and the privately checked four needs per version. The verdict line’s start and end scores come from this judge; the historical judges’ totals remain in the round log. Report the actual stopping condition and count candidate rounds separately from Round 0 and the final comparison. Preserve links to original, baseline and final captures and show the screenshots in the report. List any unresolved decisions as D1, D2 and so on, or state that none remain.

## Constraints
Never edit a report's content to improve its preview. Every preview must be an
actual page screenshot stored as an image file. Preserve the presets list and
the “Add a sample” section. Put open decisions in the final report. Support
each result with command output or evaluator output. A ten-second selection time is a goal. Only a timed test with new readers supports a measured ten-second result.

## Red-Green Proof
**RED prompt/case:** Try to support an improvement claim from the previous unscored Opus reviews or from different judges’ historical totals.
**Why RED now:** Neither supplies a common blind endpoint comparison or the requested evidence report.
**GREEN when:** A reader can trace both verdict scores to the same new judge, reconcile the round counts against the climb log, open each of the ten real reports and inspect the three sets of screenshots and gate output.
**Validation:** The user specified the acceptance cases, thresholds and evidence requirements in the verbatim request.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [ ] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [ ] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.
