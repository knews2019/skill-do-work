---
id: REQ-751
title: 'Run blind Opus comparisons and keep qualifying gallery changes'
status: completed
route: B
estimate:
  p50_active_minutes: 120
  confidence: medium
write_set:
  - docs/report-gallery/climb/README.md
  - docs/report-gallery/climb/claim.json
  - docs/report-gallery/climb/final-endpoint-check.log
  - docs/report-gallery/climb/frozen-baseline-check.log
  - docs/report-gallery/climb/round-1/candidate.patch
  - docs/report-gallery/climb/round-1/capture.log
  - docs/report-gallery/climb/round-1/decision.json
  - docs/report-gallery/climb/round-1/gate.log
  - docs/report-gallery/climb/round-1/judge-identity.json
  - docs/report-gallery/climb/round-1/judge-output.md
  - docs/report-gallery/climb/round-1/mapping.json
  - docs/report-gallery/climb/round-1/sealed-choices.json
  - docs/report-gallery/climb/round-1/views/1440-dark.png
  - docs/report-gallery/climb/round-1/views/1440-light.png
  - docs/report-gallery/climb/round-1/views/390-dark.png
  - docs/report-gallery/climb/round-1/views/390-light.png
  - docs/report-gallery/climb/round-1/views/gates.json
  - docs/report-gallery/climb/round-2/candidate.patch
  - docs/report-gallery/climb/round-2/capture.log
  - docs/report-gallery/climb/round-2/decision.json
  - docs/report-gallery/climb/round-2/gate.log
  - docs/report-gallery/climb/round-2/judge-identity.json
  - docs/report-gallery/climb/round-2/judge-output.md
  - docs/report-gallery/climb/round-2/mapping.json
  - docs/report-gallery/climb/round-2/recapture-diff.log
  - docs/report-gallery/climb/round-2/sealed-choices.json
  - docs/report-gallery/climb/round-2/views/1440-dark.png
  - docs/report-gallery/climb/round-2/views/1440-light.png
  - docs/report-gallery/climb/round-2/views/390-dark.png
  - docs/report-gallery/climb/round-2/views/390-light.png
  - docs/report-gallery/climb/round-2/views/gates.json
  - docs/report-gallery/climb/round-3/candidate.patch
  - docs/report-gallery/climb/round-3/capture.log
  - docs/report-gallery/climb/round-3/decision.json
  - docs/report-gallery/climb/round-3/gate.log
  - docs/report-gallery/climb/round-3/judge-identity.json
  - docs/report-gallery/climb/round-3/judge-output.md
  - docs/report-gallery/climb/round-3/mapping.json
  - docs/report-gallery/climb/round-3/recapture-diff.log
  - docs/report-gallery/climb/round-3/sealed-choices.json
  - docs/report-gallery/climb/round-3/views/1440-dark.png
  - docs/report-gallery/climb/round-3/views/1440-light.png
  - docs/report-gallery/climb/round-3/views/390-dark.png
  - docs/report-gallery/climb/round-3/views/390-light.png
  - docs/report-gallery/climb/round-3/views/gates.json
  - docs/report-gallery/climb/round-4/candidate.patch
  - docs/report-gallery/climb/round-4/capture.log
  - docs/report-gallery/climb/round-4/decision.json
  - docs/report-gallery/climb/round-4/gate.log
  - docs/report-gallery/climb/round-4/judge-identity.json
  - docs/report-gallery/climb/round-4/judge-output.md
  - docs/report-gallery/climb/round-4/mapping.json
  - docs/report-gallery/climb/round-4/recapture-diff.log
  - docs/report-gallery/climb/round-4/sealed-choices.json
  - docs/report-gallery/climb/round-4/views/1440-dark.png
  - docs/report-gallery/climb/round-4/views/1440-light.png
  - docs/report-gallery/climb/round-4/views/390-dark.png
  - docs/report-gallery/climb/round-4/views/390-light.png
  - docs/report-gallery/climb/round-4/views/gates.json
  - docs/report-gallery/climb/round-5/candidate.patch
  - docs/report-gallery/climb/round-5/capture.log
  - docs/report-gallery/climb/round-5/gate.log
  - docs/report-gallery/climb/round-5/recapture-diff.log
  - docs/report-gallery/climb/round-5/views/1440-dark.png
  - docs/report-gallery/climb/round-5/views/1440-light.png
  - docs/report-gallery/climb/round-5/views/390-dark.png
  - docs/report-gallery/climb/round-5/views/390-light.png
  - docs/report-gallery/climb/round-5/views/gates.json
  - docs/report-gallery/climb/round-6/candidate.patch
  - docs/report-gallery/climb/round-6/capture.log
  - docs/report-gallery/climb/round-6/decision.json
  - docs/report-gallery/climb/round-6/gate.log
  - docs/report-gallery/climb/round-6/judge-identity.json
  - docs/report-gallery/climb/round-6/judge-output.md
  - docs/report-gallery/climb/round-6/mapping.json
  - docs/report-gallery/climb/round-6/recapture-diff.log
  - docs/report-gallery/climb/round-6/sealed-choices.json
  - docs/report-gallery/climb/round-6/views/1440-dark.png
  - docs/report-gallery/climb/round-6/views/1440-light.png
  - docs/report-gallery/climb/round-6/views/390-dark.png
  - docs/report-gallery/climb/round-6/views/390-light.png
  - docs/report-gallery/climb/round-6/views/gates.json
  - docs/report-gallery/climb/round-7/candidate.patch
  - docs/report-gallery/climb/round-7/capture.log
  - docs/report-gallery/climb/round-7/decision.json
  - docs/report-gallery/climb/round-7/gate.log
  - docs/report-gallery/climb/round-7/judge-identity.json
  - docs/report-gallery/climb/round-7/judge-output.md
  - docs/report-gallery/climb/round-7/mapping.json
  - docs/report-gallery/climb/round-7/recapture-diff.log
  - docs/report-gallery/climb/round-7/sealed-choices.json
  - docs/report-gallery/climb/round-7/views/1440-dark.png
  - docs/report-gallery/climb/round-7/views/1440-light.png
  - docs/report-gallery/climb/round-7/views/390-dark.png
  - docs/report-gallery/climb/round-7/views/390-light.png
  - docs/report-gallery/climb/timing-fold.json
  - docs/report-gallery/gallery-climb-log.md
  - docs/report-gallery/index.html
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go
created_at: 2026-11-18T17:43:03Z
user_request: UR-185
domain: ui-design
prime_files: []
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: false
maintenance: false
batch: report-gallery-climb
related: [REQ-749, REQ-750, REQ-752]
depends_on: [REQ-750]
claimed_at: 2026-11-18T18:07:04Z
completed_at: 2026-11-18T18:22:53Z
commit: 702711ea
release_at: 2026-11-18T18:22:53Z
---
# Run blind Opus comparisons and keep qualifying gallery changes

## What
Run the bounded hill climb on the frozen gallery, using one candidate change per commit and a fresh Opus 5.5 evaluator via T3 Code for each scored round.

## Detailed Requirements
Hill-climbing loop — Repeat up to 12 rounds
1. Start from the current incumbent. Take its evaluator's single weakest point.
2. Make one candidate change aimed at that weakness. One idea, one commit.
3. Run the capture script, save the candidate screenshots, and run the gates.
   If a gate fails, record which one, revert only the candidate change, and
   count the round as reverted.
4. If the gates pass, run the blind A/B evaluation above.
5. Keep the candidate only if this round's judge scores it at least 3 points
   above the incumbent AND no rubric line drops by more than 2 points. A kept
   candidate becomes the incumbent. Otherwise revert only the candidate change.
6. Append one line to gallery-climb-log.md with the round, change, gate result,
   both totals, and keep/revert decision. Link the evidence. For a gate failure,
   record “unscored” rather than inventing totals.
7. After two consecutive reverts, use the next round for one larger idea, such
   as a different card layout. Apply the same gates and keep rule.

Stop when the incumbent reaches 90/100 with all gates passing and all four
reader needs matched, when three consecutive rounds are reverted, or when
12 rounds are complete. Record which condition ended the run.

## Blind comparison procedure
Use a fresh Opus 5.5 evaluator through T3 Code for each round. Randomly assign
the incumbent and candidate to A and B. Show matching screenshots at the same
viewports and themes. Include all ten cards. Give the evaluator the four reader
needs and have it record its first choice for each version before showing the
HTML. Then supply both versions' HTML and the rubric. Keep the version order,
change description, prior scores, and builder's notes hidden from the evaluator.

Score each line from 0 to 20, for a total of 100:
R1  Previews — Can a reader distinguish the reports from the images alone?
R2  Names — Does each name convey the subject and seem easy to recall?
R3  Choice — How easily can a reader select a report for each of the four needs?
R4  Range — Do the examples demonstrate distinctly different report kinds and shapes?
R5  Attention — Does each visible element earn its place? Deduct for filler and
    repeated text.

Return the five scores, total, and single weakest point with its page location
for each version. Record the four first choices against the frozen answer key.
Claim a measured ten-second result only if a timed test with new readers was run.

Keep the answer key, A/B mapping, prior scores, change description and builder notes out of judge materials, paths and HTML metadata. Preserve the same judge’s recorded screenshot-only first choices for both versions before releasing the HTML and rubric. The coordinator unmasks A/B and checks choices privately against the frozen key. Preserve first choices, evaluator output, individual line deltas, commands and candidate identity even for reverted rounds. Compare the two totals from that round’s judge for each keep decision. If the Round 0 incumbent already meets the success condition, record the stop with zero candidate rounds.

The frozen inputs from REQ-750 are the ten report identities, capture script, gate script, four reader needs, acceptable-answer key and rubric. Keep them fixed through the final comparison. Preserve the frozen baseline for REQ-752 even after a keep.

## Constraints
Never edit a report's content to improve its preview. Every preview must be an
actual page screenshot stored as an image file. Preserve the presets list and
the “Add a sample” section. Put open decisions in the final report. Support
each result with command output or evaluator output. A ten-second selection time is a goal. Only a timed test with new readers supports a measured ten-second result.

## Red-Green Proof
**RED prompt/case:** Inspect a proposed keep that is two points better overall or loses three points on one rubric line. Inspect a candidate with a failed gate.
**Why RED now:** Earlier Opus approval and design iterations have no blinded paired scores, frozen answer key or qualifying round log under these rules.
**GREEN when:** Those scored changes revert, a failed gate counts as an unscored revert, and every executed round has an evidence-linked decision. The log identifies the exact stopping condition.
**Validation:** The user specified the acceptance cases, thresholds and evidence requirements in the verbatim request.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [x] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [x] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.

## Triage

Route B. The gallery endpoint and blind procedure were established by REQ-750. This request owns the bounded comparison loop and exact generator parity.

## Plan

Apply the frozen capture and gate to each isolated change, obtain sealed screenshot-first choices and paired scores, and retain the qualifying endpoint and exact generator output.

## Exploration

The frozen manifest already maps each screenshot to a report. The gallery generator can keep those identities and bytes while presenting the cards in groups.

## Scope

**Files I will touch:**
- `docs/report-gallery/climb/README.md`
- `docs/report-gallery/climb/claim.json`
- `docs/report-gallery/climb/final-endpoint-check.log`
- `docs/report-gallery/climb/frozen-baseline-check.log`
- `docs/report-gallery/climb/round-1/candidate.patch`
- `docs/report-gallery/climb/round-1/capture.log`
- `docs/report-gallery/climb/round-1/decision.json`
- `docs/report-gallery/climb/round-1/gate.log`
- `docs/report-gallery/climb/round-1/judge-identity.json`
- `docs/report-gallery/climb/round-1/judge-output.md`
- `docs/report-gallery/climb/round-1/mapping.json`
- `docs/report-gallery/climb/round-1/sealed-choices.json`
- `docs/report-gallery/climb/round-1/views/1440-dark.png`
- `docs/report-gallery/climb/round-1/views/1440-light.png`
- `docs/report-gallery/climb/round-1/views/390-dark.png`
- `docs/report-gallery/climb/round-1/views/390-light.png`
- `docs/report-gallery/climb/round-1/views/gates.json`
- `docs/report-gallery/climb/round-2/candidate.patch`
- `docs/report-gallery/climb/round-2/capture.log`
- `docs/report-gallery/climb/round-2/decision.json`
- `docs/report-gallery/climb/round-2/gate.log`
- `docs/report-gallery/climb/round-2/judge-identity.json`
- `docs/report-gallery/climb/round-2/judge-output.md`
- `docs/report-gallery/climb/round-2/mapping.json`
- `docs/report-gallery/climb/round-2/recapture-diff.log`
- `docs/report-gallery/climb/round-2/sealed-choices.json`
- `docs/report-gallery/climb/round-2/views/1440-dark.png`
- `docs/report-gallery/climb/round-2/views/1440-light.png`
- `docs/report-gallery/climb/round-2/views/390-dark.png`
- `docs/report-gallery/climb/round-2/views/390-light.png`
- `docs/report-gallery/climb/round-2/views/gates.json`
- `docs/report-gallery/climb/round-3/candidate.patch`
- `docs/report-gallery/climb/round-3/capture.log`
- `docs/report-gallery/climb/round-3/decision.json`
- `docs/report-gallery/climb/round-3/gate.log`
- `docs/report-gallery/climb/round-3/judge-identity.json`
- `docs/report-gallery/climb/round-3/judge-output.md`
- `docs/report-gallery/climb/round-3/mapping.json`
- `docs/report-gallery/climb/round-3/recapture-diff.log`
- `docs/report-gallery/climb/round-3/sealed-choices.json`
- `docs/report-gallery/climb/round-3/views/1440-dark.png`
- `docs/report-gallery/climb/round-3/views/1440-light.png`
- `docs/report-gallery/climb/round-3/views/390-dark.png`
- `docs/report-gallery/climb/round-3/views/390-light.png`
- `docs/report-gallery/climb/round-3/views/gates.json`
- `docs/report-gallery/climb/round-4/candidate.patch`
- `docs/report-gallery/climb/round-4/capture.log`
- `docs/report-gallery/climb/round-4/decision.json`
- `docs/report-gallery/climb/round-4/gate.log`
- `docs/report-gallery/climb/round-4/judge-identity.json`
- `docs/report-gallery/climb/round-4/judge-output.md`
- `docs/report-gallery/climb/round-4/mapping.json`
- `docs/report-gallery/climb/round-4/recapture-diff.log`
- `docs/report-gallery/climb/round-4/sealed-choices.json`
- `docs/report-gallery/climb/round-4/views/1440-dark.png`
- `docs/report-gallery/climb/round-4/views/1440-light.png`
- `docs/report-gallery/climb/round-4/views/390-dark.png`
- `docs/report-gallery/climb/round-4/views/390-light.png`
- `docs/report-gallery/climb/round-4/views/gates.json`
- `docs/report-gallery/climb/round-5/candidate.patch`
- `docs/report-gallery/climb/round-5/capture.log`
- `docs/report-gallery/climb/round-5/gate.log`
- `docs/report-gallery/climb/round-5/recapture-diff.log`
- `docs/report-gallery/climb/round-5/views/1440-dark.png`
- `docs/report-gallery/climb/round-5/views/1440-light.png`
- `docs/report-gallery/climb/round-5/views/390-dark.png`
- `docs/report-gallery/climb/round-5/views/390-light.png`
- `docs/report-gallery/climb/round-5/views/gates.json`
- `docs/report-gallery/climb/round-6/candidate.patch`
- `docs/report-gallery/climb/round-6/capture.log`
- `docs/report-gallery/climb/round-6/decision.json`
- `docs/report-gallery/climb/round-6/gate.log`
- `docs/report-gallery/climb/round-6/judge-identity.json`
- `docs/report-gallery/climb/round-6/judge-output.md`
- `docs/report-gallery/climb/round-6/mapping.json`
- `docs/report-gallery/climb/round-6/recapture-diff.log`
- `docs/report-gallery/climb/round-6/sealed-choices.json`
- `docs/report-gallery/climb/round-6/views/1440-dark.png`
- `docs/report-gallery/climb/round-6/views/1440-light.png`
- `docs/report-gallery/climb/round-6/views/390-dark.png`
- `docs/report-gallery/climb/round-6/views/390-light.png`
- `docs/report-gallery/climb/round-6/views/gates.json`
- `docs/report-gallery/climb/round-7/candidate.patch`
- `docs/report-gallery/climb/round-7/capture.log`
- `docs/report-gallery/climb/round-7/decision.json`
- `docs/report-gallery/climb/round-7/gate.log`
- `docs/report-gallery/climb/round-7/judge-identity.json`
- `docs/report-gallery/climb/round-7/judge-output.md`
- `docs/report-gallery/climb/round-7/mapping.json`
- `docs/report-gallery/climb/round-7/recapture-diff.log`
- `docs/report-gallery/climb/round-7/sealed-choices.json`
- `docs/report-gallery/climb/round-7/views/1440-dark.png`
- `docs/report-gallery/climb/round-7/views/1440-light.png`
- `docs/report-gallery/climb/round-7/views/390-dark.png`
- `docs/report-gallery/climb/round-7/views/390-light.png`
- `docs/report-gallery/climb/timing-fold.json`
- `docs/report-gallery/gallery-climb-log.md`
- `docs/report-gallery/index.html`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`

## Pre-Flight

The frozen baseline checksums verified before round evaluation. The established Round 0 gate evidence and fresh evaluator score define the starting incumbent.

## Implementation Summary

- Seven isolated candidate rounds. Two met the comparative keep rule, four reverted on score and one reverted unscored after a gate failure. Three consecutive reverts ended the climb.
- Round 4 (`58f81933`) is the final incumbent. It scored 75 versus 72, matched all four frozen needs and passed all 28 gate checks. The 90-point goal was unmet. Round 3 qualified on score but was ineligible for final integration until Round 4 restored two identifiable diagnosis examples.
- Normal `ai-report-gallery` generation reproduced the exact adjudicated Round 4 HTML. Both diagnostic examples appear together. The frozen baseline checksum check passed, and all four archived endpoint screenshots match the judge inputs pixel for pixel.
- Main is the local integration target. The coordinator owns the separate endpoint comparison.

Review the [climb log](../../docs/report-gallery/gallery-climb-log.md) and [endpoint checks](../../docs/report-gallery/climb/final-endpoint-check.log).
- `docs/report-gallery/climb/README.md`
- `docs/report-gallery/climb/claim.json`
- `docs/report-gallery/climb/final-endpoint-check.log`
- `docs/report-gallery/climb/frozen-baseline-check.log`
- `docs/report-gallery/climb/round-1/candidate.patch`
- `docs/report-gallery/climb/round-1/capture.log`
- `docs/report-gallery/climb/round-1/decision.json`
- `docs/report-gallery/climb/round-1/gate.log`
- `docs/report-gallery/climb/round-1/judge-identity.json`
- `docs/report-gallery/climb/round-1/judge-output.md`
- `docs/report-gallery/climb/round-1/mapping.json`
- `docs/report-gallery/climb/round-1/sealed-choices.json`
- `docs/report-gallery/climb/round-1/views/1440-dark.png`
- `docs/report-gallery/climb/round-1/views/1440-light.png`
- `docs/report-gallery/climb/round-1/views/390-dark.png`
- `docs/report-gallery/climb/round-1/views/390-light.png`
- `docs/report-gallery/climb/round-1/views/gates.json`
- `docs/report-gallery/climb/round-2/candidate.patch`
- `docs/report-gallery/climb/round-2/capture.log`
- `docs/report-gallery/climb/round-2/decision.json`
- `docs/report-gallery/climb/round-2/gate.log`
- `docs/report-gallery/climb/round-2/judge-identity.json`
- `docs/report-gallery/climb/round-2/judge-output.md`
- `docs/report-gallery/climb/round-2/mapping.json`
- `docs/report-gallery/climb/round-2/recapture-diff.log`
- `docs/report-gallery/climb/round-2/sealed-choices.json`
- `docs/report-gallery/climb/round-2/views/1440-dark.png`
- `docs/report-gallery/climb/round-2/views/1440-light.png`
- `docs/report-gallery/climb/round-2/views/390-dark.png`
- `docs/report-gallery/climb/round-2/views/390-light.png`
- `docs/report-gallery/climb/round-2/views/gates.json`
- `docs/report-gallery/climb/round-3/candidate.patch`
- `docs/report-gallery/climb/round-3/capture.log`
- `docs/report-gallery/climb/round-3/decision.json`
- `docs/report-gallery/climb/round-3/gate.log`
- `docs/report-gallery/climb/round-3/judge-identity.json`
- `docs/report-gallery/climb/round-3/judge-output.md`
- `docs/report-gallery/climb/round-3/mapping.json`
- `docs/report-gallery/climb/round-3/recapture-diff.log`
- `docs/report-gallery/climb/round-3/sealed-choices.json`
- `docs/report-gallery/climb/round-3/views/1440-dark.png`
- `docs/report-gallery/climb/round-3/views/1440-light.png`
- `docs/report-gallery/climb/round-3/views/390-dark.png`
- `docs/report-gallery/climb/round-3/views/390-light.png`
- `docs/report-gallery/climb/round-3/views/gates.json`
- `docs/report-gallery/climb/round-4/candidate.patch`
- `docs/report-gallery/climb/round-4/capture.log`
- `docs/report-gallery/climb/round-4/decision.json`
- `docs/report-gallery/climb/round-4/gate.log`
- `docs/report-gallery/climb/round-4/judge-identity.json`
- `docs/report-gallery/climb/round-4/judge-output.md`
- `docs/report-gallery/climb/round-4/mapping.json`
- `docs/report-gallery/climb/round-4/recapture-diff.log`
- `docs/report-gallery/climb/round-4/sealed-choices.json`
- `docs/report-gallery/climb/round-4/views/1440-dark.png`
- `docs/report-gallery/climb/round-4/views/1440-light.png`
- `docs/report-gallery/climb/round-4/views/390-dark.png`
- `docs/report-gallery/climb/round-4/views/390-light.png`
- `docs/report-gallery/climb/round-4/views/gates.json`
- `docs/report-gallery/climb/round-5/candidate.patch`
- `docs/report-gallery/climb/round-5/capture.log`
- `docs/report-gallery/climb/round-5/gate.log`
- `docs/report-gallery/climb/round-5/recapture-diff.log`
- `docs/report-gallery/climb/round-5/views/1440-dark.png`
- `docs/report-gallery/climb/round-5/views/1440-light.png`
- `docs/report-gallery/climb/round-5/views/390-dark.png`
- `docs/report-gallery/climb/round-5/views/390-light.png`
- `docs/report-gallery/climb/round-5/views/gates.json`
- `docs/report-gallery/climb/round-6/candidate.patch`
- `docs/report-gallery/climb/round-6/capture.log`
- `docs/report-gallery/climb/round-6/decision.json`
- `docs/report-gallery/climb/round-6/gate.log`
- `docs/report-gallery/climb/round-6/judge-identity.json`
- `docs/report-gallery/climb/round-6/judge-output.md`
- `docs/report-gallery/climb/round-6/mapping.json`
- `docs/report-gallery/climb/round-6/recapture-diff.log`
- `docs/report-gallery/climb/round-6/sealed-choices.json`
- `docs/report-gallery/climb/round-6/views/1440-dark.png`
- `docs/report-gallery/climb/round-6/views/1440-light.png`
- `docs/report-gallery/climb/round-6/views/390-dark.png`
- `docs/report-gallery/climb/round-6/views/390-light.png`
- `docs/report-gallery/climb/round-6/views/gates.json`
- `docs/report-gallery/climb/round-7/candidate.patch`
- `docs/report-gallery/climb/round-7/capture.log`
- `docs/report-gallery/climb/round-7/decision.json`
- `docs/report-gallery/climb/round-7/gate.log`
- `docs/report-gallery/climb/round-7/judge-identity.json`
- `docs/report-gallery/climb/round-7/judge-output.md`
- `docs/report-gallery/climb/round-7/mapping.json`
- `docs/report-gallery/climb/round-7/recapture-diff.log`
- `docs/report-gallery/climb/round-7/sealed-choices.json`
- `docs/report-gallery/climb/round-7/views/1440-dark.png`
- `docs/report-gallery/climb/round-7/views/1440-light.png`
- `docs/report-gallery/climb/round-7/views/390-dark.png`
- `docs/report-gallery/climb/round-7/views/390-light.png`
- `docs/report-gallery/climb/timing-fold.json`
- `docs/report-gallery/gallery-climb-log.md`
- `docs/report-gallery/index.html`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`

## Qualification

Canonical qualification passed against the exact main integration range. The declared write set and retained implementation files agree.

## Testing

The selected Round 4 gate run passed all 28 checks, and the canonical focused probe passed checksums, exact regeneration, screenshot pixel matching, and Go gallery tests. Round 5 failed G2, G5 and G6 in all four views and was reverted unscored. The Round 7 runner printed all 28 PASS results and produced all four screenshots, then stalled after the checks during browser shutdown. Its stdout is the durable gate evidence. The candidate reverted on score.

## Review

Six fresh T3 Claude Opus 5.5 evaluators sealed four choices per variant before opening the HTML and rubric. Each sealed set matched 4/4. The winning Round 4 comparison improved Choice by two and Attention by one, holding the other three lines steady. The final HTML preserves all ten reports and at least two visible examples per kind. Exact HTML parity preserves the redundant closing div also present in its adjudicated source. The four view files are exact copies of the matching gate screenshots.

## Lessons Learned

Treat visible kind coverage as an integration requirement independent of the numeric keep rule. A proposal can also explain a cause, and classifying it as proposal-only can remove a real example from the diagnosis group. Keep card metadata edits separate from identity paths.

## Orientation

The coordinator can compare the Round 0 frozen files against the exact endpoint files identified in the climb log. The final valid gallery lives at `docs/report-gallery/index.html`; its matching images and gate results are in `docs/report-gallery/climb/round-4/views/`. This handback stops after REQ-751.
