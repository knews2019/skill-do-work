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

<!-- D-XX counter: last used D-08. Next decision: D-09. -->
