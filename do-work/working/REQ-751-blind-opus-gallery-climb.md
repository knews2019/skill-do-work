---
id: REQ-751
title: 'Run blind Opus comparisons and keep qualifying gallery changes'
status: claimed
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
- [ ] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [ ] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [ ] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.
