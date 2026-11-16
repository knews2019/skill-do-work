---
id: REQ-738
title: 'Second-opinion direction block keeps the page content on the first phone screen'
status: claimed
created_at: 2026-11-15T12:45:58Z
user_request: UR-180
domain: frontend
prime_files: ["_dev/primes/prime-releases.md", "_dev/primes/prime-action-files.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-737]
batch: relation-map-phone-layout
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_tree_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-16T17:17:19Z
---
# Second-Opinion Direction Block Keeps the Page Content on the First Phone Screen
## What
The `ai-report second-opinion` launcher prints the whole direction text (preset text plus editor direction) as an always-open block above every directed page. With the 3,057-byte relation-map preset, that block is about 1,302 px tall on a 390 px phone, so the map starts on the third screen. Change how the launcher shows the direction so the page's own content, such as the relation map, starts inside the first phone screen.
## Context
- The block is written by `secondOpinionDecorate` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go`. Line 1634 builds the `Second-opinion preset:` paragraph and the `Direction:` paragraph (`white-space:pre-line`, the full `combinedText()`), inside an `<aside aria-label="Second-opinion direction">` between the `<!-- second-opinion-header -->` markers, right after the opening body tag.
- Source of the finding: REQ-727 (ship a relation-map preset for ai-report second-opinion), D-17 and review finding F2 in `do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md` (lines 198 and 298). The review (`do-work/runs/work-2026-11-15-012821/REQ-727-review.md` line 64) measured the block at 514 px (1440 px wide) and 1,302 px (390 px wide), with the map frame top at y = 853 and y = 2,069. Without the block, the frame would start near y = 339 and y = 767.
- REQ-727's builder suggested folding the direction after its first paragraph, for example in a `<details>` element (`do-work/runs/work-2026-11-15-012821/REQ-727-handback.md` line 96). That is one option, not a requirement.
## Detailed Requirements
- The top of a directed page shows a short header: the preset name (or `custom`) and, on a tree page, the iteration line with its parent link. The full direction text is no longer an always-open block above the page content.
- The full direction text stays readable on the page, for example in a closed `<details>` element or at the end of the page. The builder picks the form.
- The head metas stay as they are: `ai-report-second-opinion-preset`, `-direction`, `-id`, `-base`, `-model` and `-date`. The report catalog and the iteration list read them.
- Decorating stays idempotent: decorating a page twice gives the same bytes, and the markers still let a second decorate strip the old block, including a block written by 0.324.0.
- **This block is on the pages of every preset** (`bold`, `visual`, `eli5`, `mindmap`, `relation-map` and the rest) and on `--editor-direction` pages, not only on the map. Check all of them:
  - Update the existing tests that assert on the header: `report_second_opinion_test.go` (around lines 968-996, the expected header marks and escaping) and `report_second_opinion_tree_test.go` (around lines 71-84).
  - Check the report kit. `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` styles every `details` and `summary` element (lines 196-199), and `presentation.js` toggles `details` elements with its "Open all details" button (lines 44-50). A new header element must look right on a kit page and on a page without the kit, and must not be toggled or restyled by the kit by accident.
  - Keep the render judge passing: the block keeps `overflow-wrap:anywhere` or an equal guard, so a long direction token never scrolls the page sideways.
  - Update the contract text that describes the block: `skills/do-work-toolbox/actions/ai-report.md` line 89 ("a header block … showing the preset and the full direction") and `skills/do-work-toolbox/docs/ai-report-guide.md` line 124 ("the page shows the direction at its top").
## Constraints
- Only the launcher's header block changes. The preset loader, the second-opinion prompt, the presets, the iteration list on `index.html` and the head metas stay unchanged.
- The relation-map preset text is REQ-737 (the relation-map preset keeps box titles readable after a click on a phone), not this REQ.
- Decided at capture: do not rewrite existing published pages under `ai-reports/`. They keep their 0.324.0 header until a new pass rebuilds them. Measure the proof on a scratch copy.
- Release rules from `_dev/primes/prime-releases.md` apply: a CHANGELOG entry and a VERSION bump.
## Dependencies
None. REQ-737 is related but independent: it changes the preset text, not the launcher. A page that REQ-737 builds after this REQ ships gets the new header on its own.
## Red-Green Proof
**RED prompt/case:** A Go test in `report_second_opinion_test.go` decorates a page with the shipped relation-map preset text as the direction. It checks that the always-visible part of the header (outside any collapsed element) does not contain the full direction text.
**Why RED now:** `secondOpinionDecorate` writes the full direction as an always-open `<p style="white-space:pre-line">` inside the header, so the whole 3,057-byte preset is visible text above the page.
**GREEN when:** (1) The Go test passes: the full direction is on the page but not as always-open text above the content, the preset name is visible at the top, the metas are byte-identical to today's, and decorating twice gives the same bytes. (2) A scratch copy of `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`, decorated by the new code, opened at 390 x 844 px, has the map frame's top edge inside the first screen (y below 844) at load, and the header above the page content is about 80 px tall or less. At 1440 x 1000 px the header is about 80 px tall or less (514 px today). (3) The page never scrolls sideways at either width, and `ai-report-judge` passes it. (4) A page decorated with a short preset (for example `bold`) and one with `--editor-direction` only still show the preset or `custom` at the top and the full direction somewhere on the page.
**Validation:** Inferred during capture. "Map inside the first 844 px screen" is the user's own criterion ("must not push the map below the first screen"). The 80 px header size is the capture's estimate of what that needs (767 px map top without the block, plus a short header), so it is approximate.
## Builder Guidance
- Certainty: the outcome (the map on the first phone screen) is the user's. The form of the folded direction is the builder's choice. Show a phone capture of the top of the page to the maintainer before the REQ is called done.
- Measure the map top after decorating, not on the editor's staging page: the block is added at finalize.
- The change touches every directed page, so run the whole `toolboxcommands` package tests, not only the new test.
## Full Context
See `do-work/user-requests/UR-180/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (bare; the row is `slugged: partial`, so a targeted form is not legal) — cost 24,714 tokens against the 2,000-token budget — matches "a session that is its own second-opinion backend" and "writer-anchored-marker-position": the change edits how the second-opinion launcher writes the header between its markers.
- `_dev/primes/lessons-action-files.md` (bare; the row is `slugged: partial`) — cost 8,761 tokens against the 2,000-token budget — matches "downstream readers": the REQ edits the header contract stated in `skills/do-work-toolbox/actions/ai-report.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: relation-map preset needs a phone layout: readable titles after a click, and the launcher's preset text must not push the map below the first screen*
