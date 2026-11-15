---
id: REQ-733
title: 'Rebuild the relation map page: left to right, nine boxes first, readable labels'
status: completed
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
route: B
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#session-is-the-backend]
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-14T23:32:26Z
  basis:
    - Route B
    - 3-file write set
    - 7 acceptance criteria
dispatch_at: 2026-11-14T23:35:12Z
builder_handback_at: 2026-11-15T00:02:44Z
integration_at: 2026-11-15T00:02:58Z
review_at: 2026-11-15T00:45:36Z
kb_status: pending
commit: b6a1e1a72f5338729c2e8acaead6f711a14212e6
heavy_verified_at: 2026-11-15T00:45:47Z
heavy_verified_revision: b6a1e1a72f5338729c2e8acaead6f711a14212e6
completed_at: 2026-11-15T00:46:26Z
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
- [x] **[PLAN]:** Generate the page from the v1c03 page's own HTML (28 boxes, 37 relations, the text, table and after-report sections reused byte for byte) with a scratch Python generator. Put the boxes in four columns, draw labels in an unscaled overlay, and show a 'core plus expanded' set of boxes. (from the builder hand-back)
- [x] **[APPLY]:** Session prepare in the worktree, staging page and reply written, browser check repeated until clean, then finalize (judge pass) and one commit by exact path. (from the builder hand-back)
- [x] **[UNIFY]:** git diff --stat main...HEAD shows only the three Scope files. The builder read the published page's head metas (id v1c03c01, base = the v1c03 page), the reply's first line, and the index.html diff (one <li> inside the tree markers). No staging leftover, no capture, no remote script or stylesheet. (from the builder hand-back)
*Source: convert it to a captured request, so we can run everything using do-work run*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and the bundle is named, but the builder must discover the existing map page's structure, the second-opinion prepare/finalize sequence, and the judge's checks before rebuilding the page.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

## Exploration

Done by the pre-dispatch agent at main `1389e2be`, by direct reads (no subagent).

**Required-lessons consult.** `do-work/lessons-index.md` has one row whose condition matches: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` names "a session that is its own second-opinion backend". Narrowed to `#session-is-the-backend` (line 17): the session form splits at the generation call; `--prepare` records the bundle digests in a sidecar, the agent writes BOTH `staging_page` and `staging_reply`, and `--finalize` re-checks and publishes. No other row matches: nothing under `skills/`, no kanban, release or shell change. Nothing dropped for budget.

**The base page.** `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` (424 lines, 86 KB):
- `:70` the launcher's metas (`-preset custom`, `-direction`, `-id v1c03`, `-base index.html`, `-model`, `-date`); `:71` the launcher's `<!-- second-opinion-header -->` block. The editor never writes these; the launcher adds them for the new page.
- `:72` `<p>Original report: <a href="index.html">index.html</a></p>` must open the body (the launcher's backlink check, `report_second_opinion.go:1486`). `:73-81` the "Changed in this iteration" list: every `#rev-<n>` link must hit a badged section, or the judge finding `AI-REPORT-SECOND-OPINION-CHANGE-LIST` withholds the `index.html` link.
- `:94-195` `#rev-1` the map: tool bar `:97-107` (zoom −/+, "Show everything", four group toggles), viewport `:109`, world with one SVG `viewBox="0 0 1940 2040"` at `:112` (four arrow markers by relation kind: order, supports, limits, gap) and 28 absolutely positioned `<button class="box g-…">` boxes at `:151-177`. The layout is a tall column: limits at x=20 and x=1620, seeds at x=400, options at x=1000, kills at x=1620, y up to 1850. This is the cause of F1.
- `:196` `#rev-2` every box as text; `:204` `#rev-3` every relation as a table row (37 rows); `:252` `#rev-4` after the report. Content to reuse; the original data list that generated the page is not in the repository, so the builder takes boxes and relations from this page's HTML.
- `:285` uses the bundle's shared `presentation.js`; `:286-424` the inline map script: `clampScale` `:297` (0.12 to 2.6), `fit` `:305`, `select` `:325-342` (fades unrelated boxes: F4), drag `:347-368`, wheel zoom `:375`, keys `:380`, group toggles `:402`. Relation labels are SVG text inside the scaled world, so they shrink with zoom (F3).
- The reply `v1c03-…-reply.md` lists Q1-Q4. Q1 (unreadable full view) is this request; Q3 (phone start) and Q4 (L6/L9 attached to the live check) stay the builder's call.

**The launcher.** `skills/do-work-toolbox/actions/ai-report.md:79` (Session sequence) and `:87` (names, tree ids, metas, iteration list). A page built on the v1c03 page gets tree id `v1c03c01`, name `index-v1c03c01-claude-<model>-custom-<slug>.html` and reply `v1c03c01-claude-<model>-custom-<slug>-reply.md`; the slug comes from the reply's `Variation slug:` line. Code: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` — `secondOpinionOutside` `:526` refuses any `git status --porcelain --untracked-files=all` entry outside the bundle (ignored files such as `.playwright-mcp/` do not count); `secondOpinionPrepare` `:581` writes `.<prefix>staging-*.html` and its `.json` record inside the bundle; finalize publishes, removes the staging leaves, runs the render judge (`report_judge.go`, headless Chrome at 1440x1000 and 390x844, light and dark) and, after `pass`, rebuilds the list between `<!-- second-opinion-tree -->` markers in `index.html` (`index.html:25-32` today).

**Shipped check.** `.gitattributes:48` marks `/ai-reports` export-ignore and `suite/modules.tsv` ships only the four `skills/` packages, so this request changes no shipped file: no release.

*Generated by the pre-dispatch agent (direct reads)*

## Scope

**Files I will touch:**
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-*.html` (new) — the rebuilt map page, written by the editor to the staging page and named by the launcher at finalize
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c03c01-*-reply.md` (new) — the editor's final reply, named by the launcher
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html` (modify, launcher only) — the iteration list between the second-opinion-tree markers gains the v1c03c01 row

**Files I will NOT touch:** every other file in the bundle (the original report body, v1c01 to v1c03 pages and replies, presentation.css, presentation.js, screenshots), anything under skills/ or suite/, CHANGELOG.md, VERSION, ai-reports/catalog.json, and any do-work/ path.

**Acceptance criteria (restated from REQ):**
- [ ] The map runs left to right in argument order: evidence, then options, then what limits them, using a wide screen.
- [ ] The first view shows nine boxes: the decision, the four options and the four kill attempts.
- [ ] Clicking a box reveals the seeds and limits that belong to it, and they can be hidden again; one control shows everything.
- [ ] Relation labels stay readable on screen at every zoom level.
- [ ] Selecting a box dims unrelated boxes only a little; their titles stay readable.
- [ ] Drag, pinch or Control-scroll zoom, buttons and keys, the detail panel with sources and relations, every box as text, every relation as a table row, and a no-JavaScript fallback all still work.
- [ ] The page was opened in a browser at 1440 px wide, one option and one kill attempt were clicked, and what was seen is recorded; finalize reports judge verdict `pass`; no file outside the three above changed; no library is loaded.

## Decisions

- **D-01** (DECIDE & STATE): The builder runs the Session form of `ai-report second-opinion` in its own worktree with `--backend claude --model <the builder's own declared model slug>` and `--editor-direction` describing the four fixes, with the v1c03 page as the target. The model in the new file name is the builder's, recorded as `session-declared`.
- **D-02** (DECIDE & STATE): Columns from left to right: evidence (seeds S1-S9 and the live check P), then the decision and options O1-O4, then kill attempts K1-K4, then limits L1-L9. The first view and its "fit" frame only the nine boxes. Exact sizes and positions are the builder's.
- **D-03** (DECIDE & STATE): Relation labels keep a constant on-screen size by counter-scaling with the zoom (or an HTML overlay); the builder picks the method. A target of about 12 to 14 CSS pixels on screen.
- **D-04** (DECIDE & STATE): Box content comes from the v1c03 page's HTML (28 boxes, 37 relations, the `#rev-2` texts and `#rev-3` rows), because the generating data list is not in the repository. Text is reused, not rewritten.
- **D-05** (DECIDE & STATE): Browser captures for the visual check are written outside the main tree (the builder's scratch folder or the worktree's ignored `.playwright-mcp/`), and every capture is deleted before hand-back; nothing is left in the main tree's `.playwright-mcp/`.
*D-06 to D-13 are from the builder hand-back (`do-work/runs/work-2026-11-15-012821/REQ-733-handback.md`).*

- **D-06** (DECIDE & STATE): Browser captures went only to the scratch folder above and were deleted (`rm -f "$S"/*.png`). The browser tool's automatic `.playwright-mcp/` files in the main tree (2 files, created by my session) were deleted by name, and the empty folder was removed. Nothing was saved in the worktree.
- **D-07** (DECIDE & STATE): The browser tool blocks `file:` URLs, so I served the bundle with `python3 -I -m http.server 47733 --bind 127.0.0.1` from the bundle folder during the check, then stopped it. It wrote nothing.
- **D-08** (DECIDE & STATE): Which boxes a click reveals comes from the relations, not from a hand list. A click shows the boxes that relate to the clicked box and are not core, for example O2 shows P and S1-S7. A revealed box can be clicked in turn (P shows L1, L6, L9). A link to a hidden box opens the shortest path to it from the nine.
- **D-09** (DECIDE & STATE): The four group toggles of v1c03 (Options, Kill attempts, Seeds, Limits) are replaced by "Nine boxes" and "Show everything", plus a "Showing N of 28 boxes" count. The click-to-reveal model does the toggles' job, and two ways to hide boxes would conflict. The zoom buttons and all keys are kept. The 0 key now fits the shown boxes.
- **D-10** (DECIDE & STATE): Labels are an HTML overlay, not counter-scaled SVG text. A label shows only for the selected box's lines and sits beside the other box of the line (a line between two stacked boxes puts its label in the gap left of the column). Overlapping labels are pushed down. With all 37 labels shown at once they could not be readable without overlap.
- **D-11** (DECIDE & STATE): On a screen narrower than 600 px, the first view frames the decision and options column. The view never zooms itself out below 42%. This answers v1c03's Q3: the nine boxes at the 42% fit would be about 11 px titles in 120 px boxes.
- **D-12** (DECIDE & STATE): The model slug is `claude-opus-5-5`, this session's declared model, recorded as `session-declared`. The launcher's date meta reads `2026-11-14`, which is the UTC date at finalize.
- **D-13** (DECIDE & STATE): Box summaries and the status text (for example "verified in code") show from 75% zoom and up. Below that, a box shows its code and a title of up to 3 lines. All boxes are 300×160 world px, so no title or summary is clipped at any zoom.

- **D-14** (DECIDE & STATE, integrator): The qualifier's four scope findings (`SCOPE-DECLARED-NOT-TOUCHED` on the two globs, `SCOPE-UNDECLARED-TOUCH` on the two resolved names) are one fact seen twice: the Scope list declares the launcher-named files as globs, and the qualifier compares paths literally. Both resolved names match their globs, so there is no scope drift. The Scope list stays as written.
- **D-15** (DECIDE & STATE, integrator): `QUALIFY-NEW-FILE-UNWIRED` on the reply file is an accepted exception. The second-opinion launcher writes the reply as a companion record beside each page, and no page links it: the v1c03 reply is unlinked too.

- **D-16** (DECIDE & STATE, integrator): The first review scored 74% (below the 75% pass line) with F1 (labels cover box titles after "Show everything" and on a phone) and F2 (Tab focus leaves a box off-screen). The integrator fixed F1, F2, F3 (stale panel after deselect), F4 (no-JS tool bar shows dead buttons) and F6 (dropped labels without notice) on the builder branch (`4be874a1`) and re-merged with the same pre (`567903a5`). F1 rule: shown boxes are taken space; label width follows the column gap (gap x zoom - 12 px, at most 190); a selection zooms in to at least 55%; a label pushed more than one box height from its line end is dropped; the count line says "N of M line words in view". On a phone this answers the reply's Q1 with "fewer labels, only where they fit": the first click shows the boxes, and dragging brings labels into the gaps. F5 and F7 stay report only. The editor's reply file is not edited: it is the editor's record, and its Q1 is answered here.
- **D-17** (DECIDE & STATE, integrator): The finalize judge `pass` covers the builder's bytes (`ed6a3a8a`). `ai-report-judge <bundle>` renders only a bundle's `index.html` today (one-page judging is REQ-728, not merged), so the fixed page cannot be re-judged now. The fix changes only the inline script and adds one CSS rule (`.map-tools[hidden]`); the integrator's browser measurement and the delta re-review stand in for the judge.

- **D-18** (DECIDE & STATE, integrator): The delta re-review of `567903a5` scored 72% with two defects the D-16 fix caused: G1 (labels narrower than their longest word broke words) and G2 (the 55% floor moved a view that already showed everything, losing related boxes; on a phone it gave 0 labels after every click). Fixed on the builder branch (`3eb2ab4b`) and re-merged with the same pre (`b6a1e1a7`): a label is 110 to 190 px wide; a click zooms in to 0.575 (where a 110 px label fits a column gap) only when the box and its related boxes all fit at that zoom, otherwise the builder's framing rule applies; the live count line is written only when its text changes (G3); "line words" became "relation words" (G4); `fitIds` lost its now-unused floor parameter. D-16's 55% floor is superseded. On a phone a click again lands at 42% as the builder designed (D-11), with the relation words that fit (0 to 3) and the count note.

- **D-19** (DECIDE & STATE, integrator): H1 (uneven zoom-in after a manual zoom-out) and H2 (the fit formula repeated in `frame()`) from the third review stay report only: the page passed at 89%, both are edge polish, and REQ-727 (the relation-map preset) will rewrite this framing rule for any report.
- **D-20** (DECIDE & STATE, integrator): The builder's proposed lesson (pin a transformed viewport's scroll to 0) gets no satellite bullet. No lessons satellite covers hand-built report pages, and the third round showed the pin alone was not enough (F2: reset the scroll before measuring in `focusin` too). The corrected lesson is kept in this REQ's `## Lessons Learned` and carried to REQ-727 through Discovered Tasks.

<!-- D-XX counter: last used D-20. Next decision: D-21. -->

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only, committed by the coordinator)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-733-preflight-probe.sh) and the repository gate green at df1bf423 (158 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*

## Implementation Summary

**Files changed:**
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html` (new)
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map-reply.md` (new)
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html` (modified, launcher only: one row in the iteration list)

**What was done:** The builder ran the Session form of `ai-report second-opinion` on the first map page (v1c03) in its own worktree and wrote a new child page v1c03c01. The new map runs in four columns left to right, opens on the nine core boxes, reveals seeds and limits on click, draws relation labels in an unscaled HTML overlay, and dims unrelated boxes to opacity 0.72. Finalize reported judge verdict `pass`.

**Integrator review fixes (D-16, D-18):** two commits on the builder branch, `4be874a1` and `3eb2ab4b`, change only the new page's inline script and add one CSS rule: labels never cover shown boxes and never break a word, a click zooms in only when the related boxes still fit, Tab keeps the focused box in view, the panel resets after a deselect, the tool bar stays hidden without JavaScript, and the count line says how many relation words are in view. Final merge `b6a1e1a7`; the page is 584 lines.

## Discovered Tasks

*(from the builder hand-back)*

- For REQ-727 (the relation-map preset): a preset needs a rule for any report, saying which kinds go in which column and which boxes form the first view. This page hard-codes both (reply Q4). Columns on a phone stay crowded when a box has many relations (reply Q1).
- The browser tool writes `.playwright-mcp/` files into the main tree on every navigation unless each call names a file. Parallel builders share that folder.

*(from the integrator and reviewers)*

- For REQ-727 (the relation-map preset): carry this page's label rules (shown boxes are taken space, 110 to 190 px labels, a count of the labels that did not fit, zoom in only when the related boxes still fit) and the `focusin` scroll reset (Lessons Learned). H1 and H2 from the third review (uneven zoom-in after a manual zoom-out, the fit formula written twice) belong in that rewrite. → report only
- F5 (no relation words on the no-JavaScript map; the relation table keeps all 37) and F7 (duplicated "Live check: Live check:" label text, inherited from v1c03). → report only

## Qualification

**Gate records (`advance --diff-range c2f21c0d..fd210a91`):** `qualify` satisfied. Findings judged: `QUALIFY-NEW-FILE-UNWIRED` on the reply file (accepted, D-15); `SCOPE-DECLARED-NOT-TOUCHED` x2 and `SCOPE-UNDECLARED-TOUCH` x2 (glob versus resolved name, no drift, D-14). No debug artifacts: the diff holds only the three bundle files.

**Diff:** `git diff c2f21c0d..fd210a91 --stat` lists 3 files, 587 insertions: the new page (566 lines), the reply (20 lines), `index.html` (+1 line inside the `second-opinion-tree` markers).

**Requirement trace (checked in the merged files):**
- Left to right in argument order: the 28 boxes sit in four columns, x=40 (live check P and seeds S1-S9), x=550 (decision D and options O1-O4), x=1060 (kills K1-K4), x=1570 (limits L1-L9).
- Nine boxes first: exactly D, O1-O4 and K1-K4 carry `data-core="true"`; the tool bar has "Nine boxes" and "Show everything".
- Seeds and limits on click, hidden again: the script shows a box when it is core, when "Show everything" is on, or when an expanded shown box relates to it (comment at the script head); the builder's browser check recorded click and second click on O2.
- Readable labels at every zoom: labels are an HTML overlay outside the scaled world (D-10), 13 px on screen per the builder's measurement.
- Light dimming: the dim rule is `opacity:.72` (v1c03 used 0.3).
- Kept features: drag, wheel and keys code present; the detail panel, the `#rev-2` text section, the 37-row relation table and the no-JavaScript view are reused from v1c03 (D-04). `viewport.scrollLeft/scrollTop` is pinned to 0.
- Browser check recorded by the builder at 1440x1000 and 390x844, light and dark, with O2 and K1 clicks; finalize judge verdict `pass`. The independent reviewer re-opens the page (## Review).
- Constraints: no file under `skills/` or `suite/`; the probe's byte check of the 12 earlier bundle files and the hash of `index.html` outside the tree block pass (`REQ-733 probe PASS` at fd210a91); only `presentation.js` (the bundle's own) is loaded, no library.

**Scope comparison:** declared 3 entries (2 globs + `index.html`), touched 3 files, each inside a declared entry. Nothing outside the bundle.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at `fd210a910f39fa3a3715fdf03531e87a53b4492d` (helpers/gate.sh, detached checkout, stage reuse off, 1-minute load 4.28 at start); GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-733-probe.sh` through `advance --probe-file`.
**Result:** ✓ Gate exit 0, wall 133 s, single run; probe `REQ-733 probe PASS`; green-gate record written for the merge revision.

**Re-gate after review fixes (D-16):** merge `567903a5b009356167b06cb74f4acc515ebfa893` (range `c2f21c0d..567903a5`), gate exit 0, wall 137 s, load 4.82 at start; probe PASS; green-gate record re-written with `record-green-gate`.

**Re-gate after second review fixes (D-18):** merge `b6a1e1a72f5338729c2e8acaead6f711a14212e6` (range `c2f21c0d..b6a1e1a7`), gate exit 0, wall 131 s, load 4.44 at start; probe PASS; green-gate record re-written.

**Integrator browser check of the second fix:** at 1440x1000, O2 13 of 13 labels at 58% and K1 3 of 3, with every related box in the frame; after "Show everything", O2 and K1 zoom to 57% with all labels, while P and L1 stay at 46% with every related box in the frame. At 390x844 every click lands at 42% with 0 to 3 labels and 0 on boxes. No label is narrower than its content (0 overflowing labels in every view). A cold load of `#box-L1` selects L1 with P in the frame. 20 arrow-key pans write the count line 3 times.

**Integrator browser check of the first fix (Playwright, 127.0.0.1, captures deleted):** at 1440x1000, O2 shows 13 of 13 labels at 58% with 0 labels on boxes and 0 label overlaps, in the nine-box view and after "Show everything"; K1 3 of 3; O3 8 of 9 at 55%. At 390x844 every click leaves 0 labels on boxes; the auto view shows 0 to 2 labels and dragging left shows 6 of 13 for O2. Tab at 3 zoom steps in: O1 to K4 each 100% in view (F2 fixed). Escape returns the panel to "Click a box" (F3). Without JavaScript the tool bar computes `display: none` (F4).

**Red-green validation:** *(from the builder hand-back; tdd: false, page check rather than unit test)*
- `REQ-733-probe.sh`: ✗ at base `1de00505` (no v1c03c01 page or reply) → ✓ at `ed6a3a8a` and again at merge `fd210a91`.
- Browser view at 1440x1000: ✗ v1c03 "Show everything" at 36% zoom, labels about 6 px, unrelated boxes at opacity 0.3 → ✓ v1c03c01 first view 9 of 28 boxes at 88%, titles 15.8 px, labels 13 px at every zoom, unrelated boxes at opacity 0.72.

**New tests added:** none (an ai-reports bundle page; the probe is a run artifact, not a repository test).

**Heavy verification plan:**
- Range: c2f21c0dabd2f1678e66807a6bff369c9fcc3867..b6a1e1a72f5338729c2e8acaead6f711a14212e6 (the same six lanes were selected at fd210a91 and 567903a5 and were green at both)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for the new page and `index.html`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for the new page and `index.html`
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for the new page and `index.html`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for the new page and `index.html`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for the new page and `index.html`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for the new page and `index.html`

*Verified by work action (integrator, run work-2026-11-15-012821)*

## Review

*Reviewed 2026-11-15T00:45:36Z by three independent opus reviewer runs (orchestrated mode, Route B). Round 1 at merge `fd210a91`: 74%, request changes (F1 labels over box titles, F2 Tab focus off-screen, F3-F7 minor). Round 2 at `567903a5` after integrator fix `4be874a1`: 72%, request changes (G1 broken words, G2 zoom floor lost related boxes, G3-G4 minor). Round 3 at `b6a1e1a7` after fix `3eb2ab4b`: 89%, approve. Full reports: `do-work/runs/work-2026-11-15-012821/REQ-733-review.md`, `REQ-733-review2.md`, `REQ-733-review3.md`. The block below is round 3.*

**Overall: 89%** | 2026-11-15T00:43:48Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 88% |
| Test Adequacy | 82% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
- None

**Minor findings:** H1. After a manual zoom-out, a click zooms in to 57% only when the related boxes fit there, so at 27% S2 jumps to 57% but P stays at 27% with 0 of 4 relation words (`index-v1c03c01-…-left-to-right-map.html:466-469`). Fix: zoom in to `Math.min(LABEL_SCALE, fitted)` whenever `scale < LABEL_SCALE && fitted > scale && fitted >= READABLE`. — impact-user-visible → report only. F5 (carried). The no-JS map has no relation words on its lines, but the relation table keeps all 37. — impact-negligible → report only. Nit H2. `fitsForWords` (`:466-467`) repeats the `fitIds` fit test with a literal 48. Fix: one `fittedScale(ids)` helper. — impact-negligible → report only. Nit F7 (carried, inherited). "Live check: Live check: …" at `:160` and in 4 panel links at `:214`, and "yet.. the decision" at `:170`. — impact-negligible → report only
**Acceptance:** Pass. Implementation and integration stages checked in Chromium at 1440×1000 and 390×844, light and dark. G1-G4 are fixed, F1-F4 and F6 stay fixed. There are 0 broken words, 0 labels on boxes and 0 label overlaps from 20% to 240%, and 0 of 100 clicks moved a related box out of a view that showed it.
**Restatement sweep:** nothing redefined (a bundle page, skipped by instruction because this is not the wave's last integration)
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)
**Rendered page:** 1440 light and dark: 9 boxes at 88% (15.8 px titles). O2 gives 17 boxes at 58% with 13 of 13 labels. K1 gives 3 of 3. O3 gives 6 of 9 at 49%. "Show everything" gives 28 at 46%. Then O2 and K1 go to 57% with every label, and P and L1 stay at 46% with every related box in the frame. 390 light and dark: D and O1-O4 at 73% (19.8 px). Each click lands at 42% with 0 to 3 whole labels and the count note. No page overflow at either size.
**Anti-bloat:** 0 extra files, tests, libraries or controls. `setCount` serves G3 and `fitsForWords` serves G2. One `fitIds` parameter was removed.

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: `c2f21c0dabd2f1678e66807a6bff369c9fcc3867`
- Target: `b6a1e1a72f5338729c2e8acaead6f711a14212e6`
- Coverage is uncertain for the new page and `index.html`, so the planner selected six lanes (the same set at each merge):
  - queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
  - queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
  - do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
  - staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
  - updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
  - installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Heavy Verification Result

- Target: `b6a1e1a72f5338729c2e8acaead6f711a14212e6`; execution revision `b6a1e1a72f5338729c2e8acaead6f711a14212e6` (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set).
- queue-kanban-javascript: executed, exit 0, 7 s
- queue-kanban-browser: executed, exit 0, 86 s
- do-work-cli-integrations: executed, exit 0, 72 s
- staged-skills: executed, exit 0, 43 s
- updater: executed, exit 0, 77 s
- installer: executed, exit 0, 39 s
- Earlier drains, superseded: all six lanes exit 0 at `fd210a91` and at `567903a5`.

## Lessons Learned

**What worked:** measuring the page in a real browser with a script (labels on boxes, label overlaps, words broken inside labels, related boxes in the frame) at both widths caught every defect the reviewers named, and the same script confirmed each fix in minutes.

**What didn't:** the builder's browser check tried one click per view, so it missed labels over titles after "Show everything" plus a click and on a phone. The first integrator fix (a fixed zoom floor and narrower labels) traded one defect for two (G1, G2); checking "related boxes still in the frame" and "no word broken in a label" before committing would have caught both.

**Worth knowing:**
- On a map moved by CSS transform inside an `overflow:hidden` viewport, focus on an off-screen box scrolls the clipped viewport before `focusin` runs. Pinning `scrollLeft`/`scrollTop` to 0 in a scroll handler is not enough: reset them in `focusin` before measuring the box too.
- Unscaled HTML labels over a scaled map need the shown boxes as taken space, a minimum width of the longest word, and a count of the labels that did not fit.
- The finalize judge judges the bytes it finalizes; a later edit to a second-opinion page cannot be re-judged one page at a time until REQ-728 lands (D-17).

## Orientation

`ai-reports/` bundles are export-ignored and ship nothing. A second-opinion page is written by the `ai-report second-opinion` launcher (`skills/do-work-toolbox/actions/ai-report.md`, Session sequence) and joins the bundle's iteration tree in `index.html`.

## Timing

Observed 2026-11-14T23:29:27Z to 2026-11-15T00:45:53Z: 1h 16m 26s total, 1h 28m 30s attributed across 9 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 38m 12s | 1 |
| builder-work | 27m 33s | 1 |
| verification-gate | 16m 14s | 4 |
| planning | 3m 47s | 1 |
| exploration-preflight | 2m 25s | 1 |
| handback-merge | 19s | 1 |

Slowest stage: review / independent review (three rounds with two integrator fixes), 38m 12s, outcome success.
Slowest command: verification-gate / heavy drain (final, b6a1e1a7), 9m 07s, exit 0, .
