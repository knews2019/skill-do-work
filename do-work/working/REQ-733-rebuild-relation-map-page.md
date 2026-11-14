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
<!-- D-XX counter: last used D-05. Next decision: D-06. -->

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only, committed by the coordinator)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-733-preflight-probe.sh) and the repository gate green at df1bf423 (158 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*
