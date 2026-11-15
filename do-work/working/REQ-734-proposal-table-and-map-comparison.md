---
id: REQ-734
title: 'Build the proposal as a supports-and-limits table and compare it with the map'
status: claimed
created_at: 2026-11-14T23:21:51Z
user_request: UR-178
domain: frontend
prime_files: []
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-733]
related: [REQ-727, REQ-733]
batch: tinka-follow-ups
write_set: ["ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c04-*.html", "ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c04-*-reply.md", "ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html"]
claimed_at: 2026-11-15T00:47:35Z
route: B
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#session-is-the-backend]
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-15T00:49:48Z
  basis:
    - Route B
    - 3-file write set
    - 7 acceptance criteria
dispatch_at: 2026-11-15T01:06:43Z
---
# Build the Proposal as a Supports-and-Limits Table and Compare It With the Map
## What
Build one more page of the same report: a table with one row per option and the columns "supports" and "limits". Then put the table page and the improved map page (REQ-733, rebuild the relation map page) side by side, so the maintainer can see which one tells the story better.
## Why
The maintainer's words: "basically compare the two, so I can see which one is better to tell the story". The assistant had doubted the map: the report is one decision with reasons for and against, and a table would show the same relations with no zooming.
## Context
- The report is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html`. Its first map page, `index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`, holds the content as 28 boxes and 37 relations, each also as text and as a table row. The table page reuses that content.
- A new page is written with `ai-report second-opinion --editor-direction '<direction>' ai-reports/2026-11-14_1547_proposal-tinka-review-surface`, with `index.html` as the base. The action is `skills/do-work-toolbox/actions/ai-report.md`, Second-opinion form, Session sequence.
## Detailed Requirements
- One row per option, O1 to O4, with the decision above the table.
- Two main columns: what supports the option and what limits it. Each cell lists the kill attempts, seeds and limits that relate to the option, each with its code, one short line and its source.
- The relations that do not touch an option (for example a limit that weakens a seed) are shown too, in a place the builder chooses, and none of the 37 relations is lost.
- No zooming and no dragging. The page reads top to bottom on a wide screen and on a phone.
- The comparison: take one picture of each page's first screen at the same width, 1440 pixels, and one at phone width. The hand-back shows the two pictures side by side with both page paths, and says for each page how many of the relations can be read on the first screen without any click.
- The hand-back does not pick a winner. The maintainer picks.
## Constraints
- No file under `skills/` or `suite/` changes, so this is no release.
- Every existing file in the bundle keeps its bytes, except the iteration list that the second-opinion launcher rebuilds in `index.html`.
- The page passes `ai-report-judge` through the second-opinion finalize step.
- The pictures for the comparison are not committed into the bundle. They live in the run's own folder.
## Dependencies
REQ-733 (rebuild the relation map page) completes first, because the comparison needs the improved map.
## Red-Green Proof
**RED prompt/case:** Look in the report folder for a page that shows the options with what supports and what limits each one, in one table.
**Why RED now:** No such page exists. The relations can be read only on the map, or as 37 separate rows.
**GREEN when:** The table page exists and passes the render judge; every one of the 37 relations appears on it; the hand-back shows the map and the table side by side at 1440 pixels and at phone width, with both paths.
**Validation:** Inferred during capture, from the maintainer's own words quoted above.
## Builder Guidance
- Certainty: firm. The table form is quoted by the maintainer; the cell layout is the builder's.
- Scope cue from the maintainer: keep language clear, do not complicate it.
- Trap: second-opinion refuses to start when Git shows changes outside the bundle, and a claimed request file is such a change in the main tree. Run the prepare and finalize steps in this request's own worktree, where the tree is clean.
- The maintainer's choice between map and table decides what REQ-727 (ship a relation-map preset) becomes.
## Full Context
See `do-work/user-requests/UR-178/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: build this one too, a table with one row per option and columns supports and limits; compare the two*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (one table page, one comparison) and the bundle is named, but the builder must read the v1c03 page's 28 boxes and 37 relations, run the second-opinion prepare and finalize sequence, and measure both pages in a browser.

**Planning:** Not required

## Decisions

- **D-01** (DECIDE & STATE): The builder runs the Session form of `ai-report second-opinion` in its own worktree with `--backend claude --model <the builder's own declared model slug>`, recorded as `session-declared`. The target is the bundle's `index.html` (the REQ names it as the base), so the new page is the fourth child of `v1`: tree id `v1c04`, name `index-v1c04-claude-<model>-custom-<slug>.html`, reply `v1c04-claude-<model>-custom-<slug>-reply.md`.
- **D-02** (DECIDE & STATE): Content comes from the v1c03 page's HTML: the 28 box texts with their source chips (`#rev-2`) and the 37 relation rows (`#rev-3`). Text is reused, not rewritten. The v1c03c01 page holds the same content.
- **D-03** (DECIDE & STATE): Column rule. A relation of kind "supports" goes in the Supports column; "limits or blocks" and "gap in the evidence" go in the Limits column. A relation between two options (the order rows O2 → O3, O3 → O4) appears in both options' rows. The decision block above the table states the order and holds the four "order" relations. The eight relations that touch no option (L1, L6, L9 on the live check; L4, L5 on S5; L2, L7 on the decision; L3 on S9) go in one short section right below the table, grouped by the box they touch.
- **D-04** (DECIDE & STATE): Every relation entry on the page carries `data-relation="r01"` to `data-relation="r37"`, numbered in the order of the v1c03 `#rev-3` table rows. An entry shown twice carries the same number twice. This makes "none of the 37 is lost" checkable by a script, not by reading.
- **D-05** (DECIDE & STATE): The page needs no JavaScript: plain HTML and CSS, the bundle's `presentation.css` allowed, no library, no remote file. At phone width each option row stacks Supports above Limits. No zoom, drag or fixed-size canvas.
- **D-06** (DECIDE & STATE): "Readable on the first screen without any click" means: at load, page top, no scroll, no click, the relation's words and both of its end codes are fully inside the viewport, not hidden or clipped, at 12 px or more. Viewports are 1440×1000 and 390×844 (the judge's sizes), light theme. The builder measures this with a script on both pages and reports four numbers (two pages × two widths) out of 37.
- **D-07** (DECIDE & STATE): The four comparison pictures (map and table, 1440×1000 and 390×844, first screen) are the REQ's deliverable, and the REQ says they live in the run's own folder. The builder writes them to `do-work/runs/work-2026-11-15-012821/REQ-734-compare-{map,table}-{1440,390}.png` in the main tree and never stages them; the coordinator commits them with the run artifacts. Every other browser capture goes to the builder's scratch folder and is deleted before hand-back. This narrows the coordinator's "captures go to scratch only" instruction for these four files only; the coordinator may override.
- **D-08** (DECIDE & STATE): The hand-back shows the two pictures side by side as a two-column Markdown table per width, with both page paths, and picks no winner.

<!-- D-XX counter: last used D-08. Next decision: D-09. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

## Exploration

Done by the pre-dispatch agent at main `b23c1771`, by direct reads (no subagent).

**Required-lessons consult.** `do-work/lessons-index.md:12` names "a session that is its own second-opinion backend" for `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`. Narrowed to `#session-is-the-backend` (line 17): the agent writes BOTH `staging_page` and `staging_reply` before `--finalize`, or finalize refuses. No other row matches: nothing under `skills/`, no kanban, release or shell change. Nothing dropped for budget. REQ-733's `## Lessons Learned` (`do-work/archive/UR-177/REQ-733-rebuild-relation-map-page.md:292`) applies to the method: measure the page with a script in a real browser, at both widths and in every state, before finalize; its review took three rounds (labels over titles, phone width, Tab focus off-screen).

**The content source.** `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`:
- `:196` `#rev-2` every box as text: 28 `<article class="detail g-…" id="detail-<code>">` items inside `<details>` groups, each with a title, text, and a `.detail-sources` chip that names its source in `index.html`.
- `:204` `#rev-3` every relation as a table row: 37 rows with From, Relation words, To and Kind (4 order, 8 supports, 14 limits or blocks, 11 gap in the evidence). The eight rows that touch no option are L1, L6, L9 → live check; L4, L5 → S5; L2, L7 → decision; L3 → S9. The row order fixes the `r01`-`r37` numbers (D-04).
- The v1c03c01 page (REQ-733) has the same content in a left-to-right map. Its top, `:79-81` (`Original report` link, then the `Changed in this iteration` nav), shows the shape the launcher's backlink and change-list checks want.

**The launcher.** `skills/do-work-toolbox/actions/ai-report.md:79` (Session sequence) and `:87` (names, tree ids, metas, iteration list). A bundle or `index.html` target means base `index.html` (tree id `v1`); its next free child is `v1c04`. Code: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` — `secondOpinionOutside` `:526` refuses any Git change outside the bundle (so prepare runs in the builder's clean worktree); `secondOpinionPrepare` `:581` writes the staging page and its record inside the bundle. Finalize runs the render judge (`report_judge.go:56`, Chrome at 1440×1000 and 390×844, light and dark) and after `pass` rebuilds the list between `<!-- second-opinion-tree -->` markers (`index.html:25-33` today). REQ-728 is merged (`9618c836`): `ai-report-judge <bundle> --page <file.html>` now re-judges one page in place (`ai-report.md:36`).

**Shipped check.** `.gitattributes:48` marks `/ai-reports` export-ignore and `suite/modules.tsv` ships only the four `skills/` packages, so this request changes no shipped file: no release.

*Generated by the pre-dispatch agent (direct reads)*

## Scope

**Files I will touch:**
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c04-*.html` (new) — the supports-and-limits table page, written by the editor to the staging page and named by the launcher at finalize
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c04-*-reply.md` (new) — the editor's final reply, named by the launcher
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html` (modify, launcher only) — the iteration list between the second-opinion-tree markers gains the v1c04 row

The four comparison pictures (D-07) are run artifacts in the run folder, written by the builder and committed by the coordinator; they are not part of the branch.

**Files I will NOT touch:** every other file in the bundle (the original report body, the v1c01, v1c02, v1c03 and v1c03c01 pages and replies, presentation.css, presentation.js, screenshots), anything under skills/ or suite/, CHANGELOG.md, VERSION, ai-reports/catalog.json, and any other do-work/ path.

**Acceptance criteria (restated from REQ):**
- [ ] One row per option, O1 to O4, with the decision above the table.
- [ ] Two main columns, Supports and Limits; each cell lists the related kill attempts, seeds and limits, each with its code, one short line and its source.
- [ ] The relations that touch no option are shown too, and all 37 relations appear on the page (r01 to r37).
- [ ] No zooming and no dragging; the page reads top to bottom at 1440 px and at 390 px, with no page overflow.
- [ ] The hand-back shows the map and the table first screens side by side at 1440 and at 390, with both page paths, and the count of relations readable on each first screen without a click; it picks no winner.
- [ ] Finalize reports judge verdict `pass`; every other bundle file keeps its bytes; `index.html` changes only inside the iteration list.
- [ ] The comparison pictures are not in the bundle; they are in the run folder.

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-734-preflight-probe.sh) and the repository gate green at the REQ-728 merge b7e6a08e (137 s)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*
