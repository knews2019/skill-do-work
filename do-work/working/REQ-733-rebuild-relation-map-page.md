---
id: REQ-733
title: 'Rebuild the relation map page: left to right, nine boxes first, readable labels'
status: claimed
created_at: 2026-11-14T23:21:00Z
user_request: UR-177
domain: frontend
prime_files: []
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-727]
batch: tinka-follow-ups
write_set: ["ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-*.html", "ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c03c01-*-reply.md", "ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html"]
claimed_at: 2026-11-14T23:28:20Z
---
# Rebuild the Relation Map Page: Left to Right, Nine Boxes First, Readable Labels
## What
Build one improved map page beside the existing one, as a new second-opinion page of the same report. The improved page fixes four problems the maintainer saw in a screenshot of the first map. Nothing under `skills/` changes.
## Why
The maintainer asked "is this the best way to present it?" about the first map page. It is not: at the full view the map used about one third of its frame, box text and line labels were too small to read, and selecting a box faded the rest too much. The improved page is the example the maintainer looks at before the `relation-map` preset (REQ-727) is built.
## Context
- The first map page is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`. It was generated from a data list of 28 boxes and 37 relations; the same content is reused.
- The four problems seen on a wide screen at 37% zoom: F1 the map is a tall column on a wide screen; F2 box text is too small; F3 the relation label is a few pixels high; F4 unrelated boxes fade until they cannot be read.
- A new second-opinion page is written with `ai-report second-opinion --editor-direction '<direction>' <page.html>`, with the first map page as the base, so the new page becomes its child in the iteration tree. The action is `skills/do-work-toolbox/actions/ai-report.md`, Second-opinion form, Session sequence.
## Detailed Requirements
- Lay the map out left to right in the order of the argument: evidence, then options, then what limits them, so a wide screen is used.
- The first view shows nine boxes: the decision, the four options and the four kill attempts. Seeds and limits appear when the reader clicks the box they belong to, and can be hidden again. A control shows everything at once.
- Line labels keep a readable size on screen at every zoom level.
- Selecting a box dims unrelated boxes only a little. Their titles stay readable.
- Keep what worked: drag to move, pinch or Control-scroll to zoom, buttons and keys, the detail panel with sources and relations, every box as text, every relation as a table row, and a page that works without JavaScript.
- Open the finished page in a browser, click one option and one kill attempt, and look at the result before calling it done. Record what was seen.
## Constraints
- No file under `skills/` or `suite/` changes, so this is no release.
- The original report page, the first map page and every other existing file in the bundle keep their bytes, except the iteration list that the second-opinion launcher rebuilds in `index.html`.
- The page passes `ai-report-judge` through the second-opinion finalize step.
- No library. Inline SVG and plain JavaScript only.
- Do not build the table form mentioned as an alternative. It was a doubt, not a request.
## Dependencies
None. REQ-727 (ship a relation-map preset) waits for this request.
## Red-Green Proof
**RED prompt/case:** Open the first map page on a screen 1440 pixels wide and press "Show everything".
**Why RED now:** The zoom is about 30 to 37%. Only the box codes can be read, a relation label is a few pixels high, and after a click the unrelated boxes cannot be read.
**GREEN when:** On the new page at the same width, the first view shows nine boxes whose titles can be read without zooming; a relation label can be read at every zoom level; a click on an option shows its seeds and limits; unrelated boxes stay readable; the second-opinion finalize step reports judge verdict `pass`.
**Validation:** Inferred during capture. The four fixes are the assistant's proposal, which the maintainer asked to capture; the proof was not confirmed as a separate question.
## Builder Guidance
- Certainty: the maintainer accepted the direction by asking to capture it. Sizes and the exact layout are the builder's.
- Scope cue from the maintainer: keep language clear, do not complicate it.
- Trap: second-opinion refuses to start when Git shows changes outside the bundle, and a claimed request file is such a change in the main tree. Run the prepare and finalize steps in this request's own worktree, where the tree is clean.
- The hand-back shows the maintainer the page path first. REQ-727 stays blocked until the maintainer has looked at this page.
## Full Context
See `do-work/user-requests/UR-177/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: convert it to a captured request, so we can run everything using do-work run*
