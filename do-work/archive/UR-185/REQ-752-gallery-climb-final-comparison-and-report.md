---
id: REQ-752
title: 'Compare the frozen endpoints and publish the gallery climb evidence'
status: completed
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
    - Route B
    - 17-file write set
    - 7 acceptance criteria
  calculated_at: 2026-11-18T18:26:00Z
write_set:
  - docs/report-gallery/final/baseline-full-opus.md
  - docs/report-gallery/final/choice-verification.json
  - docs/report-gallery/final/final-gates.log
  - docs/report-gallery/final/generate.py
  - docs/report-gallery/final/index.html
  - docs/report-gallery/final/judge-identity.json
  - docs/report-gallery/final/judge-output.md
  - docs/report-gallery/final/parity.log
  - docs/report-gallery/final/report-check.log
  - docs/report-gallery/final/report.json
  - docs/report-gallery/final/reports.json
  - docs/report-gallery/final/round-7-recheck/exit-status
  - docs/report-gallery/final/round-7-recheck/frozen-baseline-check.log
  - docs/report-gallery/final/round-7-recheck/gate.stdout
  - docs/report-gallery/final/round-7-recheck/gates.json
  - docs/report-gallery/final/round-7-recheck/verification.log
  - docs/report-gallery/final/sealed-choices.json
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
claimed_at: 2026-11-18T18:24:29Z
completed_at: 2026-11-18T18:28:20Z
commit: f726f11cd7bbba4fe985e861ea1e6c0590cc180a
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
- [x] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [x] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [x] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.

## Triage

Route B. The frozen protocol specifies the endpoints and acceptance criteria. UX decisions belong in the report.

## Plan

Compare the randomized frozen endpoints, verify first choices and publish the evidence with a local report.

## Exploration

Round 4 was the retained endpoint and its recorded captures match the neutral packet.

## Scope

**Files I will touch:**
- `docs/report-gallery/final/baseline-full-opus.md`
- `docs/report-gallery/final/choice-verification.json`
- `docs/report-gallery/final/final-gates.log`
- `docs/report-gallery/final/generate.py`
- `docs/report-gallery/final/index.html`
- `docs/report-gallery/final/judge-identity.json`
- `docs/report-gallery/final/judge-output.md`
- `docs/report-gallery/final/parity.log`
- `docs/report-gallery/final/report-check.log`
- `docs/report-gallery/final/report.json`
- `docs/report-gallery/final/reports.json`
- `docs/report-gallery/final/round-7-recheck/exit-status`
- `docs/report-gallery/final/round-7-recheck/frozen-baseline-check.log`
- `docs/report-gallery/final/round-7-recheck/gate.stdout`
- `docs/report-gallery/final/round-7-recheck/gates.json`
- `docs/report-gallery/final/round-7-recheck/verification.log`
- `docs/report-gallery/final/sealed-choices.json`

## Implementation Summary

The fresh Opus 5.5 pair scored 73 and 80. Both sealed choice sets matched 4/4. Added the complete Round 0 transcript as separate evidence and included the Round 7 recheck. The report lists the three open UX and timing decisions.

**Files changed:**
- `docs/report-gallery/final/baseline-full-opus.md`
- `docs/report-gallery/final/choice-verification.json`
- `docs/report-gallery/final/final-gates.log`
- `docs/report-gallery/final/generate.py`
- `docs/report-gallery/final/index.html`
- `docs/report-gallery/final/judge-identity.json`
- `docs/report-gallery/final/judge-output.md`
- `docs/report-gallery/final/parity.log`
- `docs/report-gallery/final/report-check.log`
- `docs/report-gallery/final/report.json`
- `docs/report-gallery/final/reports.json`
- `docs/report-gallery/final/round-7-recheck/exit-status`
- `docs/report-gallery/final/round-7-recheck/frozen-baseline-check.log`
- `docs/report-gallery/final/round-7-recheck/gate.stdout`
- `docs/report-gallery/final/round-7-recheck/gates.json`
- `docs/report-gallery/final/round-7-recheck/verification.log`
- `docs/report-gallery/final/sealed-choices.json`

## Qualification

Canonical qualification and scope drift passed against the exact integrated range. The report is reproducible and its historical and paired scores have separate sources.

## Testing

Final G1–G7 passed 28/28, and the canonical focused probe passed. Chrome opened the report from disk at 390 and 1440 px in both themes, resolved 43 links and decoded all 12 screenshots. Both generators reproduced exact HTML, and the baseline checksum verification passed.

## Review

Both A/B choice arrays equal the first sealed record. Each handle resolves in the HTML for its label to the matching frozen report ID. The report distinguishes the new 73/80 pair from the historical 79 and 75 results and records the actual stop after three consecutive reverts.

## Lessons Learned

A blind pair supports a same-evaluator comparison even when historical judges varied. A correct untimed choice establishes selection accuracy for that evaluator.

## Orientation

Open docs/report-gallery/final/index.html from disk. It links the evidence and records D1–D3. The endpoint remains Round 4. This documentation-only change carries no package release.
