---
id: REQ-750
title: 'Establish the seven gallery checks and freeze the scored baseline'
status: completed
route: B
estimate:
  p50_active_minutes: 55
  confidence: medium
write_set:
  - _dev/tests/fixtures/report-gallery-legacy-selectors.json
  - _dev/tests/report-gallery-gate.cjs
  - _dev/tests/report-gallery-gate.md
  - docs/report-gallery/baseline.html
  - docs/report-gallery/baseline/README.md
  - docs/report-gallery/baseline/SHA256SUMS
  - docs/report-gallery/baseline/answer-key.json
  - docs/report-gallery/baseline/capture.log
  - docs/report-gallery/baseline/choice-verification.json
  - docs/report-gallery/baseline/judge-identity.json
  - docs/report-gallery/baseline/judge-output.md
  - docs/report-gallery/baseline/reader-needs.json
  - docs/report-gallery/baseline/red-mismatched-image/exit-status
  - docs/report-gallery/baseline/red-mismatched-image/gates.json
  - docs/report-gallery/baseline/red-mismatched-image/run.log
  - docs/report-gallery/baseline/red-original/1440-dark.png
  - docs/report-gallery/baseline/red-original/1440-light.png
  - docs/report-gallery/baseline/red-original/390-dark.png
  - docs/report-gallery/baseline/red-original/390-light.png
  - docs/report-gallery/baseline/red-original/exit-status
  - docs/report-gallery/baseline/red-original/gates.json
  - docs/report-gallery/baseline/red-original/run.log
  - docs/report-gallery/baseline/rubric.txt
  - docs/report-gallery/baseline/sealed-choices.json
  - docs/report-gallery/baseline/views/1440-dark.png
  - docs/report-gallery/baseline/views/1440-light.png
  - docs/report-gallery/baseline/views/390-dark.png
  - docs/report-gallery/baseline/views/390-light.png
  - docs/report-gallery/baseline/views/exit-status
  - docs/report-gallery/baseline/views/gates.json
  - docs/report-gallery/baseline/views/run.log
  - docs/report-gallery/capture-previews.py
  - docs/report-gallery/index.html
  - docs/report-gallery/preview-manifest.json
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go
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
claimed_at: 2026-11-18T18:04:11Z
completed_at: 2026-11-18T18:05:51Z
commit: c3fd91f4b198da6c1769ea97bdf43b6fdab55ee7
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
- [x] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [x] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [x] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.

## Triage

Route B. The prepared gate and real-gallery implementation define the scope. The scored baseline requires an independent blind evaluation.

## Plan

Planning not required for Route B. Wire the prepared gate to the capture manifest and normal generator, run the browser matrix and freeze the green build before dispatch.

## Exploration

The prepared legacy red run exposed fragment and preview issues. The real-report gallery already had direct links and bare fragments. Its unused dialog image had no source, so the generator now emits the dialog only when gallery samples use it. The mismatch receipt exercises G2 independently of image loading.

## Scope

**Files I will touch:**
- `_dev/tests/fixtures/report-gallery-legacy-selectors.json`
- `_dev/tests/report-gallery-gate.cjs`
- `_dev/tests/report-gallery-gate.md`
- `docs/report-gallery/baseline.html`
- `docs/report-gallery/baseline/README.md`
- `docs/report-gallery/baseline/SHA256SUMS`
- `docs/report-gallery/baseline/answer-key.json`
- `docs/report-gallery/baseline/capture.log`
- `docs/report-gallery/baseline/choice-verification.json`
- `docs/report-gallery/baseline/judge-identity.json`
- `docs/report-gallery/baseline/judge-output.md`
- `docs/report-gallery/baseline/reader-needs.json`
- `docs/report-gallery/baseline/red-mismatched-image/exit-status`
- `docs/report-gallery/baseline/red-mismatched-image/gates.json`
- `docs/report-gallery/baseline/red-mismatched-image/run.log`
- `docs/report-gallery/baseline/red-original/1440-dark.png`
- `docs/report-gallery/baseline/red-original/1440-light.png`
- `docs/report-gallery/baseline/red-original/390-dark.png`
- `docs/report-gallery/baseline/red-original/390-light.png`
- `docs/report-gallery/baseline/red-original/exit-status`
- `docs/report-gallery/baseline/red-original/gates.json`
- `docs/report-gallery/baseline/red-original/run.log`
- `docs/report-gallery/baseline/rubric.txt`
- `docs/report-gallery/baseline/sealed-choices.json`
- `docs/report-gallery/baseline/views/1440-dark.png`
- `docs/report-gallery/baseline/views/1440-light.png`
- `docs/report-gallery/baseline/views/390-dark.png`
- `docs/report-gallery/baseline/views/390-light.png`
- `docs/report-gallery/baseline/views/exit-status`
- `docs/report-gallery/baseline/views/gates.json`
- `docs/report-gallery/baseline/views/run.log`
- `docs/report-gallery/capture-previews.py`
- `docs/report-gallery/index.html`
- `docs/report-gallery/preview-manifest.json`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go`

Preserve all source reports, the presets and the sample workflow. Scope is UR-185 / REQ-750.

## Pre-Flight

The original-gallery and mismatched-image runs returned exit 1. See the red receipts under `docs/report-gallery/baseline/`.

## Implementation Summary

**Files changed:**
- `_dev/tests/fixtures/report-gallery-legacy-selectors.json`
- `_dev/tests/report-gallery-gate.cjs`
- `_dev/tests/report-gallery-gate.md`
- `docs/report-gallery/baseline.html`
- `docs/report-gallery/baseline/README.md`
- `docs/report-gallery/baseline/SHA256SUMS`
- `docs/report-gallery/baseline/answer-key.json`
- `docs/report-gallery/baseline/capture.log`
- `docs/report-gallery/baseline/choice-verification.json`
- `docs/report-gallery/baseline/judge-identity.json`
- `docs/report-gallery/baseline/judge-output.md`
- `docs/report-gallery/baseline/reader-needs.json`
- `docs/report-gallery/baseline/red-mismatched-image/exit-status`
- `docs/report-gallery/baseline/red-mismatched-image/gates.json`
- `docs/report-gallery/baseline/red-mismatched-image/run.log`
- `docs/report-gallery/baseline/red-original/1440-dark.png`
- `docs/report-gallery/baseline/red-original/1440-light.png`
- `docs/report-gallery/baseline/red-original/390-dark.png`
- `docs/report-gallery/baseline/red-original/390-light.png`
- `docs/report-gallery/baseline/red-original/exit-status`
- `docs/report-gallery/baseline/red-original/gates.json`
- `docs/report-gallery/baseline/red-original/run.log`
- `docs/report-gallery/baseline/rubric.txt`
- `docs/report-gallery/baseline/sealed-choices.json`
- `docs/report-gallery/baseline/views/1440-dark.png`
- `docs/report-gallery/baseline/views/1440-light.png`
- `docs/report-gallery/baseline/views/390-dark.png`
- `docs/report-gallery/baseline/views/390-light.png`
- `docs/report-gallery/baseline/views/exit-status`
- `docs/report-gallery/baseline/views/gates.json`
- `docs/report-gallery/baseline/views/run.log`
- `docs/report-gallery/capture-previews.py`
- `docs/report-gallery/index.html`
- `docs/report-gallery/preview-manifest.json`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go`

The normal generator reproduces the frozen HTML byte for byte. The manifest records hashes for the ten fixed sources and captures. The G1–G7 runner reads the actual browser clipboard, follows bare fragments and preview links, and checks the full page in each view. The baseline directory contains the original and green browser evidence, frozen inputs, and sealed evaluation.

## Qualification

The exact main integration range passed canonical qualification and scope drift after changing terminal log writes to explicit CLI stdout. The first qualification pass had flagged two output statements as library instrumentation. The output text and gate checks are unchanged.

## Testing

- The legacy gallery and a loadable but mismatched capture both returned exit 1. The committed receipts preserve the failures by gate and view. The mismatch independently failed G2.
- The real gallery returned exit 0 and passed G1–G7 at 1440 and 390 pixels in both themes. The main checkout rerun reproduced every gate result.
- The canonical focused probe passed the freeze checksums, byte-for-byte live and frozen HTML comparison, and focused Go gallery tests. The normal generator reproduced the scored HTML. Source report bundles are unchanged.

## Review

Fresh T3 Claude Opus 5.5 read the four anonymous screenshots, sealed four first choices, then read the released full HTML and rubric. R1–R5 scored 17, 15, 16, 17, and 14 for 79/100. All four first choices matched the coordinator's key. The judge located the weakest point at the top-row third card, `#release-dossier`, where the name, preview headline, and subtitle suggest different subjects. The raw judge response, T3 identity, sealed choices, and verification are in `docs/report-gallery/baseline/`. Reader selection time was not measured.

## Lessons Learned

A hidden image element with no source still fails an actual image-decode audit. Emit the sample lightbox only when the gallery has samples. Capture provenance needs source and PNG hashes together to catch a replacement that still decodes.

## Orientation

Use `docs/report-gallery/baseline/README.md` for exact reproduction commands. `docs/report-gallery/baseline.html` is the complete scored snapshot. `docs/report-gallery/preview-manifest.json` freezes the ten source identities and capture hashes. Retain the sealed evaluation and fixed needs for the next round. Keep the answer key private from evaluators.
