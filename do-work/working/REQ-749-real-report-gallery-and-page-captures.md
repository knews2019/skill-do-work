---
id: REQ-749
title: 'Addendum: Build the gallery from ten real report pages'
status: claimed
created_at: 2026-11-18T17:43:03Z
user_request: UR-185
addendum_to: REQ-743
domain: frontend
prime_files: []
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-gallery-climb
related: [REQ-750, REQ-751, REQ-752]
claimed_at: 2026-11-18T17:55:31Z
---
# Addendum: Build the gallery from ten real report pages

## What
Give a new reader ten distinct reports to choose from at `docs/report-gallery/index.html`, with a screenshot, a copyable public name and a direct link to each report. Establish the true before-state first.

## Prior Implementation
REQ-743 built the gallery generator, sample discovery, presets and Add a sample. REQ-744 added evidence-based sample selection. Commits e143415e and 137d723d made the gallery wide and responsive and added ten named cards, real screenshot assets, clipboard controls and report links. Reuse those pieces where they meet the new checks.

At 137d723d, the curated cards open shipped `skills/do-work-toolbox/docs/report-gallery/*/example.html` through a separate report link. Six adapt repository evidence and four use labeled fictional data. The preview opens an image viewer, the names copy in title case and the card fragments begin with `#sample-`. The earlier Opus reviews and browser checks are useful design history. They did not use the new frozen corpus, answer key or blind rubric, so they provide no starting score for this run.

## Required result
Show exactly ten worked reports. Each card displays a screenshot of the actual
report page. Each report has a short, memorable name. Clicking the name copies
it. Clicking the preview opens the report. The gallery URL with #name scrolls
directly to that card. Names are unique, lowercase, two or three hyphenated words, at most 24 characters including the hyphens. The visible/copied name is also the bare URL fragment.

## Round 0 setup
1. Verify that playwright-cli can capture a screenshot. If it cannot and no
   equivalent browser capture tool works, stop and report the failed command.
2. Save screenshots of the original gallery before changing it.
3. Select ten real reports from ai-reports/. Include at least two examples of
   each kind: completed work, proposal, root cause, and architecture. If examples
   are missing, use do-work-toolbox ai-report with real archived REQs or URs.
4. Write a capture script that screenshots the top of each report's index.html
   at a fixed viewport and saves each image as a file next to the gallery.
   Record which report produced each image.

Use the current gallery immediately before this run as the original. Preserve its screenshots and the successful capture command (or fallback result) before altering the gallery. The four matched gallery views are 1440 px and 390 px in light and dark mode. Capture the source report at a consistent fixed viewport; the image displayed by the card must match the top of the `index.html` it opens. Preserve the source-to-image mapping for G2. Check legacy untyped reports against their evidence when selecting the ten. At capture time the catalog had no typed root-cause index page. Inspect actual bundles before using the requested writer fallback.

## Constraints
Never edit a report's content to improve its preview. Every preview must be an
actual page screenshot stored as an image file. Preserve the presets list and
the “Add a sample” section. Put open decisions in the final report. Support
each result with command output or evaluator output. A ten-second selection time is a goal. Only a timed test with new readers supports a measured ten-second result.

## Red-Green Proof
**RED prompt/case:** Follow a current card’s preview, copy its name and load the gallery with that copied fragment.
**Why RED now:** The committed markup opens a screenshot viewer, copies a spaced display title, and targets a prefixed fragment. Its report links point to shipped examples.
**GREEN when:** Each of the ten source reports satisfies the requested kind coverage. Its page image opens that exact report, the clicked name copies its lowercase slug, and loading that slug as the fragment scrolls to the card.
**Validation:** The user specified the acceptance cases, thresholds and evidence requirements in the verbatim request.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (25210 tokens, partially slugged). Its exact-basename and existence-versus-identity lessons apply to the gallery generator’s names and target paths.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [ ] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [ ] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.
