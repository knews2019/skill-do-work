---
id: REQ-738
title: 'Second-opinion direction block keeps the page content on the first phone screen'
status: completed
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
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-11-16T17:21:08Z
  basis:
    - Route B
    - 8-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
claimed_at: 2026-11-16T17:17:19Z
dispatch_at: 2026-11-16T17:31:30Z
builder_handback_at: 2026-11-16T19:38:30Z
integration_at: 2026-11-16T19:38:42Z
kb_status: pending
heavy_verified_at: 2026-11-16T19:50:56Z
heavy_verified_revision: 943c43858982436dcdae54ec9ccc80caed0e5468
review_at: 2026-11-16T19:47:13Z
commit: 943c43858982436dcdae54ec9ccc80caed0e5468
completed_at: 2026-11-16T19:51:50Z
release_at: 2026-11-16T19:51:50Z
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
- [x] **[PLAN]:** Keep the markers, metas, strip calls and insert position; replace only the slot paragraphs with a closed, kit-neutral details; pin with a direct-call test; measure on scratch copies. (from the builder hand-back)
- [x] **[APPLY]:** Three commits: RED test (975066d9), the slot branch of secondOpinionDecorate and its comment (152d6bbb), two contract clauses (6432bb22). (from the builder hand-back)
- [x] **[UNIFY]:** git diff --stat main...HEAD shows 4 files (launcher +3/-1 lines, test +75, ai-report.md 1 line, ai-report-guide.md 1 line); no debug output, no scratch file, ai-reports/ unchanged; gofmt, vet and the package tests are clean. (from the builder hand-back)
*Source: relation-map preset needs a phone layout: readable titles after a click, and the launcher's preset text must not push the map below the first screen*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and the code location is named (one function, `secondOpinionDecorate`). The form of the folded direction, the report-kit rules that touch it and the tests that pin the header still need discovery, so exploration comes before the build.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Where the change goes**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:1602-1648`, function `secondOpinionDecorate`. It strips old blocks (1607-1608), strips copied metas (1610-1618), builds the metas (1625-1631), builds the visible paragraphs (1631 the iteration line with the parent link, 1634 the preset paragraph and the always-open `Direction:` paragraph), and inserts one block right after the opening body tag (1641-1646). The function comment at 1597-1601 says "the slot and the full direction" and must change with the code.
- Only lines 1631, 1634 and 1641 decide what a reader sees. The metas, the strip calls and the insert position stay as they are.
- Callers: line 841 (strip only, before the backlink check) and line 864 (the real decorate at publish). Neither needs a change.
- `secondOpinionStripBlocks` (1563-1589) removes everything between the two marker comments, whatever is inside. So a block written by 0.324.0 is stripped by the new code with no extra work, as long as the new block keeps the same two markers.

**Tests that pin the header today**
- `report_second_opinion_test.go:942` `TestSecondOpinionDirectedRunNamesAndMarksPage`. Line 968 asserts the two metas, both markers and the free text "add a timeline" on a `bold` plus editor-direction page. Lines 990-1001 assert idempotence, the metas before the real head end with escaping, and that a plain decorate strips back to the input bytes. These asserts hold for any form that keeps the markers and the metas, so they may need no edit. Baseline at 7f2c4b94: PASS in 1.65 s.
- `report_second_opinion_tree_test.go:43` `TestSecondOpinionTreeTargetsNamesAndParentHeader`. Line 79 asserts the parent link and the id and base metas on a child page; line 84 asserts a plain run has no marks. Baseline: PASS in 2.41 s.
- No test asserts the words "Second-opinion preset:" or "Direction:", so the visible wording is free.
- The whole second-opinion test family runs close to the 30-second per-file budget (REQ-705's lesson, quoted in the REQ-727 brief). A new test that calls `secondOpinionDecorate` directly costs milliseconds; a new test that runs the fake backend costs about 2 s.

**Other readers of the block (grep over `skills`, `_dev`, `suite`)**
- No Go code reads the visible block. The catalog and the iteration list read the head metas only (`readReportMetas`, used at line 1675).
- Contract text: `skills/do-work-toolbox/actions/ai-report.md:89` ("a header block between the markers right after its opening body tag, showing the preset and the full direction") and `skills/do-work-toolbox/docs/ai-report-guide.md:124` ("the page shows the direction at its top"). The same line 89 also says "A tree page's header names its id and links its parent page", which stays true.
- **A published page styles the block by selector.** REQ-727's page (`ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`) carries its own rule `body>aside[aria-label="Second-opinion direction"]{...font-size:.8rem;border:1px solid...}` and `... p{margin:.25em 0}` (REQ-727 D-13). A future editor can copy that. So the outer element should stay an `aside` with that exact label, a direct child of body, and the short header text should stay in `p` elements.

**Report kit**
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css:196-199` styles every `details` (surface background, 2 px border, 14 px radius, padding 4px 18px), `details + details` (top margin), every `summary` (pointer, bold 1.1rem display font, padding 12px 0) and `details[open] summary` (dashed bottom border, 12 px bottom margin). An unstyled `details` in the header would get a card look and about 60 px of height on a kit page. Inline styles on the header's own `details` and `summary` win over these rules.
- `presentation.css:361` (`.depth-three details > *:not(summary)`) cannot match: the header is not inside `.depth-three`.
- `presentation.js:44-50`: the "Open all details" button toggles the `details` inside its own closest `section`. The header sits right after the body tag, outside every section, so the button cannot toggle it.
- `presentation.css:211-214` print rules do not touch `details` or `aside`.
- The three second-opinion pages in the tstack bundle link neither kit file (they are self-contained), and two of them (v1c01, v1c02, both `custom`) hold five or six `details` of their own. The kit check therefore needs a scratch page made from `skeleton.html` with the kit files beside it.

**The proof page**
- The v1c03 page above is 110,479 bytes, has no `details`, and holds a 0.324.0 block right after its body tag. Decorating it with the new code is the old-block strip case. Its metas give the inputs to rebuild the same direction: preset `relation-map`, id `v1c03`, base `index.html`, model `claude claude-opus-5-5`, date `2026-11-15`; the preset text is `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (3,057 bytes).
- No command re-decorates a published page. The finalize path publishes a new page. A scratch copy needs a direct call of `secondOpinionDecorate` (D-02).

**Lessons read**
- `_dev/primes/lessons-releases.md` (required): no archive move and no manifest edit here; the release is the integrator's.
- By fixed-string grep, not whole files: `[family: writer-anchored-marker-position]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:103-105` (a launcher block carries no whitespace outside its markers; strip on the masked surface; decorating twice gives the same bytes) and `[family: session-is-the-backend]` line 18 (measure again after finalize, because the block moves the map). `_dev/primes/lessons-action-files.md:15` (a changed contract must reach every reader: the grep above is that sweep).
- `do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md` Lessons Learned, D-13, D-17 and F2.

*Generated by the pre-dispatch agent (no Explore subagent; the guide forbids subagents)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — the header block built in secondOpinionDecorate and its function comment
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — the new RED test for the folded direction; existing header asserts only if the new form breaks them
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_tree_test.go` (modify) — the tree header asserts, only if the new form breaks them (D-08)
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — the sentence in line 89 that describes the header block
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — the clause in line 124 that says where the page shows the direction
- `CHANGELOG.md` (modify, integrator only) — release entry
- `skills/do-work/CHANGELOG.md` (modify, integrator only) — byte-identical mirror
- `VERSION` (modify, integrator only) — version bump (D-05)

**Files I will NOT touch:** the preset loader, the prompt builder, the iteration list and every other function of the launcher file; the second-opinion prompt; every preset; the report kit (skeleton, presentation CSS and JS); the judge; every published file under the reports folder (the proof is a scratch copy, D-02); anything under the queue folder.

**Acceptance criteria (restated from REQ):**
- [ ] A Go test, RED first, decorates a page with the shipped relation-map preset text and shows that the always-visible part of the header does not contain the full direction text.
- [ ] The full direction is still on the page, the preset name (or custom) is visible at the top, and a tree page shows its iteration line with the parent link.
- [ ] The head metas are byte-identical to today's: preset, direction, id, base, model and date.
- [ ] Decorating twice gives the same bytes, and a page that holds a 0.324.0 block decorates to the same bytes as a clean page.
- [ ] A scratch copy of the tstack v1c03 page, decorated by the new code, at 390 x 844 has the map frame top inside the first screen (y below 844) and a header of about 80 px or less; at 1440 x 1000 the header is about 80 px or less.
- [ ] The page never scrolls sideways at either width, the block keeps its long-token guard, and the render judge passes the scratch copy.
- [ ] A short-preset page (bold) and an editor-direction-only page still show the preset or custom at the top and the full direction on the page; on a report-kit page the header is not restyled by the kit and not toggled by the "Open all details" button.
- [ ] The two contract sentences describe the new header.
- [ ] The maintainer is shown a phone capture of the top of the page before the REQ is called done (D-06).

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-16-171725/REQ-738-preflight-probe.sh) and the repository gate green at 6f00ab48 (134 s, detached checkout, stage reuse off, load 1.9). Three earlier base runs (at 7f2c4b94 and 6f00ab48) were red only on the 30 s per-file budget, each on a different file (finalization_recovery_test.go 31.15 s at load 28; report_second_opinion_test.go 30.89 s at load 11; finalization_req499_test.go 34.79 s at load 6), from other projects and the builders loading the machine.
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-16-171725)*

## Decisions

(pre-dispatch)

- **D-01** (DECIDE & STATE): Recommended form, which the builder may change if the measurement fails (it then records its own decision and names the forms it rejected). The block stays in the same place, between the same two markers, as `<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere">` and a direct child of body, because a published page already styles that selector (REQ-727 D-13). Inside it: one short always-visible part in `p` elements (the iteration line with the parent link on a tree page, and the preset name or `custom`), then one closed `details` element whose `summary` is a short label such as "Full direction" and whose body is the `white-space:pre-line` direction text. The header's `details` and `summary` carry inline styles that cancel the report kit's card look (background, border, radius, padding, font, the open-state border and margin). No script, so it works with JavaScript off. Rejected: a block at the end of the page (a second insert position and a second place to strip, and a reader at the top has no hint it exists); metas only (the direction is then not readable on the page, which the REQ requires); a CSS line clamp (no way to open it without a script).
- **D-02** (DECIDE & STATE): The browser proof needs a page decorated by the new code, and no command re-decorates a published page. The builder writes one throwaway Go test file in its worktree that calls `secondOpinionDecorate` on a copy of a page and writes the result to its scratch folder, with the direction and tree values rebuilt from the page's own metas and the shipped preset file. The file is never committed and is deleted before hand-back. It also checks that the scratch copy's head metas equal the published page's, which proves the rebuild used the same direction. Published pages under the reports folder are not rewritten (the REQ's capture decision).
- **D-03** (DECIDE & STATE): "About 80 px or less" is the capture's estimate and approximate. The measured value is the height of the header `aside` box plus its margins at load, with the fold closed. Up to about 100 px counts as met when the map frame top is below y = 844 on the phone; the hand-back gives the exact numbers. The open state is measured too and reported, with no limit: a reader chose to open it.
- **D-04** (DECIDE & STATE): The RED test is one new function, `TestSecondOpinionHeaderFoldsTheFullDirection`, in `report_second_opinion_test.go`. It calls `secondOpinionDecorate` directly (milliseconds, so the file stays inside its 30-second budget). It reads the shipped relation-map preset text from the toolbox folder and pins: the direction text is absent from the header once the folded element is removed, and present inside it; the preset name is in the visible part; the six metas equal a literal expected string; decorating twice gives the same bytes; a page holding a literal 0.324.0 block decorates to the same bytes as the clean page; a `bold` page and an editor-direction-only page show `bold` or `custom` in the visible part; the long-token guard is on the block. The GREEN probe names this test.
- **D-05** (DECIDE & STATE): The integrator picks the bump. Suggested: a patch (0.324.0 to 0.324.1), because an existing header changes its form and no option is added. The builder only proposes the changelog entry.
- **D-06** (ESCALATE): The REQ says to show a phone capture of the top of the page to the maintainer before the REQ is called done. The run is unattended, so the builder saves its final captures under the run folder's `REQ-738-captures/` and the coordinator shows them before integration starts. Value: the maintainer sees the chosen form on a real page before it ships to every directed page. Risk: if the maintainer wants another form, the change is one function and one test, fully reversible in a follow-up; pages published in between keep the form they got until a new pass rebuilds them.
- **D-07** (DECIDE & STATE): The run's browser matrix (click every shown box from each starting view) is about map pages. This REQ does not change a map, so the builder's one script measures the header and the map top at load on the scratch copy, then runs the click sweep on that copy once per width only to show the header change moved nothing else (sideways page scroll, boxes out of frame). It adds what the new element can break: Tab order (the `summary` is a new first focus stop), the no-JavaScript view, the opened fold with a long token, and a report-kit page with the "Open all details" button.
- **D-08** (DECIDE & STATE): The REQ lists both test files as files to update. Exploration shows their header asserts (markers, metas, the parent link, the free text) hold for the D-01 form. The builder edits them only when its form breaks an assert; an untouched file is then reported under Scope changes as declared but not needed, not as drift to fix.

(from the builder hand-back)

- **D-09** (DECIDE & STATE): The preset name moved into the `summary` ("Preset: <slot> (full direction)") instead of its own `p`. D-01's form measured y = 848 on the phone, above the 844 limit. The iteration line keeps its `p`, the `aside`, its label and the long-token guard stay. The test reads the summary as visible and the `details` body after `</summary>` as the folded part.
- **D-10** (DECIDE & STATE): The visible label is now "Preset:" instead of "Second-opinion preset:", and the folded text no longer starts with "Direction:". The `aside` label "Second-opinion direction" and the "(full direction)" words carry that context, and the shorter label keeps the summary on one phone line for every shipped preset name (checked with `visual-studyguide`).
- **D-11** (DECIDE & STATE): The click sweep reran with real mouse clicks because the map ignores synthetic click events. Boxes outside the frame at the start (phone first view) and three desktop boxes whose centre click hit another element are reported, not forced. This REQ does not change the map (D-07).
- **D-12** (DECIDE & STATE): The judge accepted the scratch bundle outside the repository, so nothing was copied into `ai-reports/`.
- **D-13** (DECIDE & STATE): I deleted one file my session created in the main tree's `.playwright-mcp/`: `console-2026-11-16T17-33-03-035Z.log` (one favicon 404 entry). The captures were written straight to the captures folder.
- **D-14** (DECIDE & STATE): Commit trailers use `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`, following the brief as updated on disk and the dispatch prompt.

(integrator)

- **D-15** (DECIDE & STATE): Re-passing the recorded pre-flight argv to confirm the phase rewrote `do-work/working/baseline.json` to name this REQ's pre-flight probe. That file was committed by REQ-737's pre-flight and the brief says to leave it alone, so I restored it with `git checkout -- do-work/working/baseline.json` right away; the green-gate record the re-pass wrote lives under `.git/` and stays.
- **D-16** (DECIDE & STATE): Review finding F2 (a lesson line still stated the 0.324.0 header size in the present tense) is fixed in the same release commit as this REQ's lesson: line 18 of `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` now says the block moved the map down until 0.324.1 folded it. It is one clause in a file this finalization already writes, so no follow-up is needed.
- **D-17** (DECIDE & STATE): The brief asks for the integration report in this finalization's commit, but the report must name the finalization commit, which does not exist until the commit is made. The report is written after finalization and committed by exact path as one plain `[REQ-738] integration report` commit, the same way a lesson-link fix is committed after finalization.

<!-- D-XX counter: last used D-17. Next decision: D-18. -->

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified) — the slot branch of secondOpinionDecorate and its function comment
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified) — new TestSecondOpinionHeaderFoldsTheFullDirection (7 checks) and the html import
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — one clause of the header-block sentence
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — one clause about where the page shows the direction

**What was done:** The launcher's header block now shows the iteration line with its parent link and one summary line "Preset: <name> (full direction)" of a closed, inline-styled details element that holds the full direction, so the 3,057-byte relation-map direction no longer pushes the map off a phone's first screen (map top y = 830 at 390 x 844, was 2,069). Markers, metas, strip calls and insert position are unchanged; the two contract sentences say the direction is folded.

## Discovered Tasks

(from the builder hand-back)

- **impact-user-visible** On a page without its own header rule (the report kit, S3; the custom page, S2), the header text has no side gutter: it starts at x = 0 on a phone. 0.324.0 has the same gap. A fix would add a padding to the `aside`, but an inline padding would override the tstack page's own rule, so it needs its own decision. → report only

## Qualification

**Gate records:** `advance --diff-range f995337a..943c4385` → qualify gate satisfied (after the Implementation Summary stopped backticking identifiers, which the qualifier reads as claimed paths). Scope-drift reports four declared-not-touched paths and no undeclared touch: `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` and `VERSION` are written by finalization (D-05); `report_second_opinion_tree_test.go` was declared conditionally and its asserts hold for the new form (D-08). Accepted.

**Requirement trace (against `git diff f995337a..943c4385` and the files at 943c4385):**
- Short header at the top: the block now holds the iteration `p` with the parent link (tree pages) and a `summary` "Preset: <slot> (full direction)"; the full direction is no longer always-open text (report_second_opinion.go lines 1634-1636).
- Full direction still on the page: inside the closed `details` as `p style="white-space:pre-line"`, the same escaped `text` as before.
- Head metas unchanged: the metas code path is untouched by the diff; test check 3 pins a literal six-meta string.
- Idempotent and strips a 0.324.0 block: markers, strip calls and insert position untouched; checks 4 and 5 pin both, check 5 with the literal 0.324.0 block.
- Every preset and editor direction: the change is in the shared slot branch, so it applies to every slot; check 6 covers `bold` and an editor-direction-only page (`custom`).
- Report kit: inline styles on `details` and `summary` cancel the kit card rules (S3 measured: border 0, radius 0, padding 0, weight 400); the header sits outside every `section`, so "Open all details" leaves it alone (S3 measured `headerFoldOpen false`).
- Long-token guard kept: `overflow-wrap:anywhere` on the `aside` (check 7); S4 400-char token measured sideways 0 open and closed.
- Contract text: `ai-report.md` line 89 and `ai-report-guide.md` line 124 describe the fold; the GREEN probe pins both. `git grep "Second-opinion preset:" -- skills/` finds only the test's literal 0.324.0 block, which is intentional.
- Phone proof: builder measured header 87.4 px and map frame top y = 830 at 390 x 844 (under D-03's "about 100 px" reading), header 87.4 px at 1440 x 1000; the maintainer saw the captures and approved the form on 2026-11-16.

**Scope comparison:** declared 8 paths, touched 4 project paths, all declared. No debug output, no scratch file, no file under `ai-reports/` or the queue folder in the range.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at 943c4385 in the run's detached checkout, stage reuse off (`helpers/gate.sh`), started 2026-11-16T19:39:41Z at 1-minute load 1.95 with no other gate running; then `advance REQ-738 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-11-16-171725/REQ-738-probe.sh`.
**Result:** ✓ Gate passed on the first run (wall 129 s, 987 Go tests in the fast stage, slowest file `finalization_recovery_test.go` 17.65 s, under the 30 s budget). GREEN probe passed under advance (`run-blocked-check` and `green-gate` satisfied). Builder: `go test -count=1 ./internal/toolboxcommands/` ok in 31.08 s, `-run '^TestSecondOpinion'` 23.6 s after the change (19.2 s before, machine shared with a second builder; the new test runs in 0.00 s).

**Red-green validation:**
- `report_second_opinion_test.go` TestSecondOpinionHeaderFoldsTheFullDirection: ✗ after RED commit 975066d9 (check 1: "the header has no closed details fold, so the whole direction is always-open text") → ✓ after 152d6bbb (PASS 0.00 s).
- GREEN probe `REQ-738-probe.sh`: ✗ at base 6f00ab48 ("has no TestSecondOpinionHeaderFoldsTheFullDirection", exit 1) → ✓ at 6432bb22 and again under advance at 943c4385.

**New tests added:**
- TestSecondOpinionHeaderFoldsTheFullDirection in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (7 checks: fold closed and holds the whole direction, visible slot and parent link, metas byte-identical, idempotent, 0.324.0 block stripped, `bold` and `custom` pages, label and long-token guard).

**Existing tests updated (cross-REQ impact):** none. TestSecondOpinionDirectedRunNamesAndMarksPage and TestSecondOpinionTreeTargetsNamesAndParentHeader pass unchanged (D-08).

**Heavy verification plan:**
- Range: f995337a9c4d503f27b349d0b39a30451813abba..943c43858982436dcdae54ec9ccc80caed0e5468
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the launcher and its test matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four changed paths matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the launcher and its test matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the launcher and its test matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action (integrator, run work-2026-11-16-171725)*

## Review

**Overall: 96%** | 2026-11-16T19:47:13Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: page CSS that hides summary markers (17 published files under `ai-reports/` do this for their own folds) would also hide the header fold's ▸; the inline styles do not cancel it (`report_second_opinion.go:1636`) — impact-user-visible → report only
- F2: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:18` still states the 0.324.0 block size (514 px / 1,302 px) in the present tense; now about 87 px — impact-negligible → report only
- F3: confirmed pre-existing gap from the REQ's Discovered Tasks: no side gutter for the header on kit and plain pages (`asideLeft` 0 px on the kit page at 390) — impact-user-visible → report only
- F4 (Nit): the new test does not pin the kit-cancelling inline styles on `details`/`summary` — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: focused Go tests green at 943c4385; merged decorate on scratch S1/kit/plain pages measured in Chromium (390 x 844: header 87.4 px, map top y = 830; 1440 x 1000: header 87.4 px; 0 px sideways scroll closed and open; Tab+Enter opens; works with JS off; kit card cancelled and "Open all details" leaves the header closed); deployment and live acceptance unassessed
**Restatement sweep:** redefined the header block's visible form ("Second-opinion preset:" / "Direction:" paragraphs → "Preset: <name> (full direction)" summary of a closed details); grep over `skills/`, `_dev/`, `suite/` finds the two contract sentences updated in this diff and the test's intentional 0.324.0 literal; one stale restatement at `lessons-do-work-cli.md:18` (F2)
**Suggested testing:** 4 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: f995337a9c4d503f27b349d0b39a30451813abba
- Target: 943c43858982436dcdae54ec9ccc80caed0e5468
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the launcher and its test matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four changed paths matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the launcher and its test matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the launcher and its test matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target: 943c43858982436dcdae54ec9ccc80caed0e5468; execution revision 943c43858982436dcdae54ec9ccc80caed0e5468 (the run's detached checkout, browser variable set, started 2026-11-16T19:47:31Z at 1-minute load 1.92)
- do-work-cli-integrations: executed, exit 0, 61 s
- staged-skills: executed, exit 0, 36 s
- updater: executed, exit 0, 66 s
- installer: executed, exit 0, 26 s

## Lessons Learned

**What worked:** Measuring the candidate header forms on the real tstack page with its own CSS at 390 x 844 before committing to one. The pre-dispatch recommendation (D-01: iteration `p`, preset `p`, "Full direction" summary) measured y = 848, 4 px past the first screen; putting the preset name into the summary line (D-09) made it fit at y = 830. A direct call of the decorate function kept the new test at 0.00 s inside a file that sits near its 30 s budget.
**What didn't:** Synthetic click events do not select boxes on the relation map (it listens for pointer events), so the builder's first click sweep proved nothing and had to be rerun with real mouse clicks (D-11).
**Worth knowing:** Inline styles on the header's own `details` and `summary` are what keep the report kit's card look off it; no test pins them yet (review F4). A page rule that hides summary markers would also hide the header's ▸ (review F1).

## Orientation

The `ai-report second-opinion` launcher's header block (`secondOpinionDecorate` in the toolbox commands of the do-work CLI) is the only launcher output a reader sees above a directed page; it now folds the direction, and the head metas remain the machine-readable record.

## Timing

Observed 2026-11-16T19:38:07Z to 2026-11-16T19:50:56Z: 12m 49s total, 11m 01s attributed across 4 events, 1m 48s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 5m 42s | 2 |
| review | 4m 43s | 1 |
| handback-merge | 36s | 1 |

Slowest stage: review / independent review, 4m 43s, outcome success.
