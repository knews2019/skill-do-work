---
id: REQ-750
title: 'Establish the seven gallery checks and freeze the scored baseline'
status: pending
created_at: 2026-11-18T17:43:03Z
user_request: UR-185
domain: testing
prime_files: []
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-climb
related: [REQ-749, REQ-751, REQ-752]
depends_on: [REQ-749]
---
# Establish the seven gallery checks and freeze the scored baseline

## What
Make Round 0 reproducible. Run the checks against the source-linked gallery, record the reader tasks and obtain the first Opus score.

## Detailed Requirements
5. Write a gate script that exits 0 only when all seven checks pass:
   G1  Exactly ten example cards appear.
   G2  Each card displays the nonempty screenshot captured from the report it opens.
   G3  Each name is unique, lowercase, two or three hyphenated words, at most
       24 characters in total.
   G4  Clicking each name copies that exact name. Verify the clipboard in a browser.
   G5  Opening the gallery at #name scrolls to that card. Its preview opens the
       correct report.
   G6  Every link resolves and every image loads.
   G7  At 1440 px and 390 px, in both light and dark mode, names and controls
       remain visible and usable, with no overlap or horizontal overflow.
6. Write four concrete reader needs, one per report kind. Record acceptable
   report choices for each before evaluation. Keep that answer key from the judge.
7. Bring the gallery through all seven gates. Save its screenshots and get a
   baseline score from a fresh Opus 5.5 evaluator through T3 Code. This is the
   first incumbent. Freeze the ten report identities, capture script, gate
   script, reader needs, answer key, and rubric for the rest of the run.

Record command output per check, viewport and theme. G4 must read the browser clipboard after the click. For G5 start from a URL containing the exact copied name and follow the preview to the intended report. G6 covers every gallery link and image. Use source identity and capture output for G2 so a valid image attached to the wrong report fails. Preserve a reproducible frozen baseline gallery, its matching views with all ten cards and its original-before-change screenshots for later comparison.

Give each reader need one or more acceptable report identities before the evaluator sees it. Keep the answer key private from every evaluator. The baseline is evaluated on its own with the screenshot-first protocol below; the randomized A/B comparison applies once a pair exists. Retain the five line scores, total, first choices, privately checked matches and located weakness. The first incumbent exists only after G1–G7 pass and this new evaluation is recorded.

## Frozen evaluator protocol and rubric
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

## Constraints
Never edit a report's content to improve its preview. Every preview must be an
actual page screenshot stored as an image file. Preserve the presets list and
the “Add a sample” section. Put open decisions in the final report. Support
each result with command output or evaluator output. A ten-second selection time is a goal. Only a timed test with new readers supports a measured ten-second result.

## Red-Green Proof
**RED prompt/case:** Run the proposed checks on the original gallery, including following a preview or loading its bare name as a fragment. A mismatched report/screenshot pair must also fail G2 even when the file loads.
**Why RED now:** Earlier browser checks covered a lightbox, prefixed anchors and shipped examples. There is no committed G1–G7 gate or frozen score for this protocol.
**GREEN when:** The gate exits zero only when every named check passes. Linked command output and a fresh T3 Code Opus 5.5 result establish the frozen first incumbent, four needs and private key.
**Validation:** The user specified the acceptance cases, thresholds and evidence requirements in the verbatim request.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` (9341 tokens, partially slugged). The gate must read the files and exercise the behaviors its output claims to verify.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [ ] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [ ] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.
