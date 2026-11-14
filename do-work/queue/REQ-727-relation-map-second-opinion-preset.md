---
id: REQ-727
title: 'Ship a relation-map preset for ai-report second-opinion'
status: blocked
blocked_by: 'Not looked yet, hold it: the maintainer looks at the example map page and confirms that it should become a standard option'
blocked_at: 2026-11-14T22:47:58Z
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
depends_on: [REQ-733]
related: [REQ-733, REQ-726, REQ-728, REQ-729, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "CHANGELOG.md", "VERSION"]
---
# Ship a Relation-Map Preset for ai-report second-opinion
## What
Add one shipped second-opinion preset, `relation-map`, so `ai-report second-opinion --preset relation-map <bundle>` writes a page that shows a report as one zoomable map of its items and the relations between them.
## Why
The maintainer wants report content that is easy to understand by zooming in and out and by seeing relations. One example page was built by hand and passed the render judge.
## Context
- The example is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`, built with `--editor-direction`.
- A preset is one Markdown file with a direction to the editor. The loader reads `<bundle>/second-opinion-presets/<name>.md`, then the shipped folder `skills/do-work-toolbox/docs/second-opinion-presets/`. No loader change is needed.
- `mindmap.md` in that folder is the nearest existing preset. A mind map is a tree from one center; a relation map shows relations between any two items.
## Detailed Requirements
- Read `mindmap.md` first, and say in the new preset how the two differ.
- Boxes are the report's own items (for example options, risks, findings, limits, files, requests). Lines are relations the report states. The editor says on the page that the lines are its reading of the report.
- At most four kinds of line, each with one color and a plain name, and a key that explains them.
- Semantic zoom: zoomed out, a box shows its code and title; zoomed in, a summary; a click opens the detail with its sources and its relations, and fades unrelated boxes.
- Drag to move, pinch or Control-scroll to zoom, buttons and keys for the same, and a button that shows everything. Plain scrolling must still scroll the page.
- The page works without JavaScript: the map shows at full size and scrolls. Every box is also given as text, and every relation as one table row.
- No library. Inline SVG and plain JavaScript only, as the second-opinion prompt already requires.
- About 30 boxes is the upper limit for one map. With more, group items and say so.
- On a narrow screen the map starts on the main item, not on the whole picture.
- Every fact still comes from the report or the repository, as in `bold.md`.
- Files: `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (new); the preset list in "A Second Opinion" in `skills/do-work-toolbox/docs/ai-report-guide.md`; `TestSecondOpinionShippedPresetsLoadFromTheToolbox` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`, if it names the shipped presets; `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
## Constraints
- The preset loader, the second-opinion prompt, the other presets and every existing `ai-reports/` bundle stay unchanged.
- The report kit (`skeleton.html`, `presentation.css`, `presentation.js`) stays unchanged.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** `ai-report second-opinion --prepare --preset relation-map` on a fixture bundle, through the existing shipped-presets test.
**Why RED now:** The call is refused because no preset named `relation-map` ships.
**GREEN when:** The same call returns state `prepared` with `preset_file` in the shipped folder; the guide names the preset; one real page is built with the preset on a bundle other than the example, passes `ai-report-judge`, and is shown to the maintainer before the request is called done.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Open Questions
- [x] Has the maintainer looked at the example map page, and should it become a standard option? → Not looked yet, hold it.
  Note (2026-11-14): The request is captured as blocked and waits for the maintainer. The maintainer was told the risk of a yes: the map was readable for 28 boxes, and larger reports are untested. `do-work clarify` releases it when the maintainer confirms.
## Builder Guidance
- Certainty: the direction is the maintainer's; the page itself has not been approved. Do not build before the block is released.
- Scope cue from the maintainer on 2026-11-14: keep language clear, do not complicate it.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Addendum (2026-11-14)

Captured with UR-177 (rebuild the relation map page). The maintainer saw the first map page and asked whether it is the best way to present the content. It is not, so REQ-733 builds an improved page first.

- This request now depends on REQ-733.
- The example the maintainer looks at before releasing the block is the page REQ-733 builds, not the first map page.
- The preset asks the editor for the improved design: a left-to-right layout, a first view of about nine main boxes with the rest opening on click, line labels that stay readable at every zoom level, and only light dimming of unrelated boxes. Where this differs from the Detailed Requirements above, this addendum wins.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-relation-map-preset.md*
