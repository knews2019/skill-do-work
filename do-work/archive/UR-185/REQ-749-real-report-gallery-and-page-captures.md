---
id: REQ-749
title: 'Addendum: Build the gallery from ten real report pages'
status: completed
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
route: B
estimate:
  p50_active_minutes: 30
write_set: [docs/report-gallery/, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go]
completed_at: 2026-11-18T18:03:50Z
commit: 171a0c386a3498608607c5f0f04f8a058973e1a9
release_at: 2026-11-18T18:03:50Z
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
- [x] **[PLAN]:** Inspect prior work and identify the evidence needed for this request.
- [x] **[APPLY]:** Carry out the scoped work under the batch constraints.
- [x] **[UNIFY]:** Verify results against the user’s criteria and retain the evidence.

## Full Context
See `do-work/user-requests/UR-185/input.md` for the complete protocol and prior-work credits.

## Triage

**Route: B** — The existing generator supplies the entry point. Source report selection and browser evidence require exploration.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

The original browser capture decoded every card image at both viewport sizes and in both themes. The original renderer used a dialog target and prefixed fragments. The archived duration lane audit explains the constant overflow Y coordinate and consequent reliance on labels. The queue velocity analysis traces recurring backlog to intake and review follow-ups using archive counts. Both are legacy untyped reports with causal evidence. Neither needed a newly authored replacement.

## Scope

**Files I will touch:**
- `docs/report-gallery/capture-previews.py`
- `docs/report-gallery/index.html`
- `docs/report-gallery/original/README.md`
- `docs/report-gallery/original/original-desktop-dark.png`
- `docs/report-gallery/original/original-desktop-light.png`
- `docs/report-gallery/original/original-phone-dark.png`
- `docs/report-gallery/original/original-phone-light.png`
- `docs/report-gallery/preview-manifest.json`
- `docs/report-gallery/previews/checkout-control-repo.png`
- `docs/report-gallery/previews/cli-command-map.png`
- `docs/report-gallery/previews/memory-note-routing.png`
- `docs/report-gallery/previews/overflow-labels.png`
- `docs/report-gallery/previews/phase-by-phase.png`
- `docs/report-gallery/previews/queue-refills.png`
- `docs/report-gallery/previews/release-dossier.png`
- `docs/report-gallery/previews/request-lifecycle.png`
- `docs/report-gallery/previews/second-opinion.png`
- `docs/report-gallery/previews/tinka-review.png`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go`

Change the repository gallery, capture inputs and evidence, generator, and focused generator regression test. Keep the shipped design lookup, project sample workflow, and preset list. Use only the ten selected source pages.

## Pre-Flight

The original gallery was captured at e26f4962. The focused regression failed against the original renderer with zero worked report cards.

## Implementation Summary

**Files changed:**
- `docs/report-gallery/capture-previews.py`
- `docs/report-gallery/index.html`
- `docs/report-gallery/original/README.md`
- `docs/report-gallery/original/original-desktop-dark.png`
- `docs/report-gallery/original/original-desktop-light.png`
- `docs/report-gallery/original/original-phone-dark.png`
- `docs/report-gallery/original/original-phone-light.png`
- `docs/report-gallery/preview-manifest.json`
- `docs/report-gallery/previews/checkout-control-repo.png`
- `docs/report-gallery/previews/cli-command-map.png`
- `docs/report-gallery/previews/memory-note-routing.png`
- `docs/report-gallery/previews/overflow-labels.png`
- `docs/report-gallery/previews/phase-by-phase.png`
- `docs/report-gallery/previews/queue-refills.png`
- `docs/report-gallery/previews/release-dossier.png`
- `docs/report-gallery/previews/request-lifecycle.png`
- `docs/report-gallery/previews/second-opinion.png`
- `docs/report-gallery/previews/tinka-review.png`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_page.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_test.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_gallery_worked.go`

The manifest pins ten report-to-image pairs and the browser capture viewport. The gallery generator uses the manifest when present. It renders the ten worked reports with direct preview links, copyable lowercase handles, and bare fragment targets. Project additions retain their existing tools and appear under Added samples. Isolated integration range is 29f6952f..8e1ac984. Step 9 records the originating merge 171a0c38 for release provenance because the finalizer evaluates the supplied merge’s first-parent diff and the last merge only amended repository captures.

## Red-green validation

TestAIReportGalleryShowsWorkedReportsAndKeepsProjectSamples failed against the original renderer (zero worked cards) and passed with the manifest-aware renderer. Browser verification exercised every report link, clipboard handle, and direct fragment. All ten preview images decoded and matched the paths in the manifest.

## Qualification

Pass. Advance satisfied qualification and scope drift on 29f6952f..8e1ac984. The renderer uses the same manifest as the capture command. Git reports no changes to source report bundles.

## Testing

- The full maintainer gate passed on the final integrated revision 8e1ac984.
- Focused Go gallery suite passed. The new regression separately fails with the unchanged renderer substituted.
- Playwright verified all ten preview paths and report navigations, clipboard contents, bare fragment targets, both themes at 1440 and 390 pixels, the preset and report-kind links, and copying the Add a sample command.
- A fresh capture produced the same ten PNGs byte for byte after waiting for report paint. The normal generator reproduced the committed HTML.

## Review

Acceptance passed in an inline review. The ten card sources are three completed-work pages, three proposals, two untyped causal investigations, and two architecture reports. Existing source HTML and assets were untouched. The new manifest is optional for existing gallery consumers and leaves the legacy renderer available to projects without one. The formal scored baseline and reader timing are outside this request.

## Lessons Learned

Report navigation can finish before a sticky header reaches its final composited position. Waiting two animation frames after layout and returning to the top removed a duplicated header observed in the first capture.

## Orientation

Regenerate previews with `python3 docs/report-gallery/capture-previews.py`. Regenerate the HTML with `skills/do-work/tools/do-work-cli.sh --format text ai-report-gallery`. The source-to-image manifest is `docs/report-gallery/preview-manifest.json`. Original desktop and phone captures and reproduction instructions are under `docs/report-gallery/original/`.
