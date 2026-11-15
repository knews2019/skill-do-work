---
id: REQ-734
title: 'Build the proposal as a supports-and-limits table and compare it with the map'
status: completed
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
builder_handback_at: 2026-11-15T03:14:46Z
integration_at: 2026-11-15T03:15:01Z
review_at: 2026-11-15T03:37:01Z
heavy_verified_at: 2026-11-15T03:37:01Z
heavy_verified_revision: 13322db24340462b2871754cb14be84165655c84
kb_status: pending
commit: 13322db24340462b2871754cb14be84165655c84
completed_at: 2026-11-15T03:37:22Z
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
- [x] **[PLAN]:** generate the page from the v1c03 page's own box texts and relation rows with a scratch Python script, so all 37 relations keep their words and their r01-r37 numbers; place them by D-03; measure every state in a browser before finalize. (from the builder hand-back)
- [x] **[APPLY]:** prepare in the worktree, write the staging page and reply, three matrix runs to clean, finalize (judge pass), one gutter fix, matrix and re-judge pass, commit by exact path, comparison pictures and counts. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat main...HEAD` lists only the v1c04 page, the v1c04 reply and `index.html` (+1 line); the probe pins every other bundle file's bytes and passes. (from the builder hand-back)
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

(D-09 to D-15 from the builder hand-back)

- **D-09** (DECIDE & STATE): each entry reads "From → To: relation words", then "Code: box title" and one source chip. The brief asks for the other end's code, the words and the source. Many relation words need their target to make sense ("K2 → O3: is not answered by"), and a bare code tells a reader nothing, so the entry also names the other box by its v1c03 title. The source chip is the first chip of that box's `#rev-2` text; the full chip lists stay in the v1c03 page and in `index.html`.
- **D-10** (DECIDE & STATE): captures. Every working capture went to the builder's scratch folder `build-REQ-734/shots/` and was deleted. The Playwright tool wrote one console log to the main tree's `.playwright-mcp/`, a folder that did not exist before the builder's session; the builder deleted the log and the folder. The four comparison pictures are the only captures kept.
- **D-11** (DECIDE & STATE): order relations. The decision block shows r01 to r04. Each option row's header repeats the order relations that touch it: O1 gets r02, O2 gets r01 and r03, O3 gets r03 and r04, O4 gets r04. Repeats carry the same number (D-04). O1 has no supports or limits in the report, so both its cells say "None" instead of an empty cell.
- **D-12** (DECIDE & STATE): one CSS edit after finalize. The launcher's direction block (`<aside>` after `<body>`) had no side gutter on this page. The builder added `body > aside { padding: 10px clamp(16px, 2.6vw, 60px) 0; }` and `body > aside p { margin: 0 0 .3em; }` to the page's own `<style>`, then re-judged with `ai-report-judge --page` (pass) and re-ran the browser matrix (clean).
- **D-13** (DECIDE & STATE): no `presentation.js`. The page links `presentation.css` only. It needs no behavior, so it has no script at all and reads the same with JavaScript off.
- **D-14** (DECIDE & STATE): a short closing section with the three corrections and the 0.316.0 / 0.317.0 note, each re-checked by the builder against the repository.
- **D-15** (DECIDE & STATE): relations on the map page are counted by the same generic rule as on the table page (words plus both end codes visible as text in the viewport). The map shows its line words only after a click, so its first-screen count is 0 at both widths by design.

(integrator)

- **D-16** (DECIDE & STATE): the catalog rebuild is the integrator's merge follow-up, as the integrate brief asks: `ai-report-index` on the merged tree gave 50 bundles and added the v1c04 row under the Tinka bundle, committed as `0ccdf98d` (`ai-reports/catalog.json`, `ai-reports/index.html`). Scope says the builder does not touch `ai-reports/catalog.json`; the builder did not, so Qualification judges the builder range `d31159be..c812ea4b` and the gate and heavy lanes run at `0ccdf98d`, the tree that holds both.

- **D-17** (DECIDE & STATE): Scope edited at integration, entries only. The two glob entries now name the files the launcher created, because the scope check matches exact paths and a launcher-named file cannot be known before finalize. The catalog pair is added as the integrator's follow-up (D-16); "Files I will NOT touch" still holds for the builder, whose range touched only the three bundle files.

- **D-18** (DECIDE & STATE): review F1 and F2 are fixed on the page only, on the builder branch (`e760b5b5`), re-merged with the same pre as `13322db2`. The reply file `v1c04-claude-claude-opus-5-5-custom-supports-limits-table-reply.md` keeps its bytes, although its line 15 repeats the old catalog claim: it is the launcher-published record of what the editor said at its base `18bd280a`, and the page is what a reader opens. The two `ai-report.md:87` chips on K4 stay, because they are the original report's own citations (D-02).
- **D-19** (DECIDE & STATE): no lessons-satellite bullet. The builder's proposed lesson (check the published page, not only the staging page) is a tip for operating the second-opinion session, not a lesson about changing the CLI code, so no satellite owns it. It is kept in `## Lessons Learned`; the durable fix is the discovered task for a gutter in the launcher's direction block.
- **D-20** (DECIDE & STATE): `skills/do-work-toolbox/actions/help.md:18` (review F4) is an incomplete one-line summary, not a false restatement: it omits `--page` and says nothing untrue. It goes to the coordinator as a report-only discovered task and does not widen this no-release REQ.

<!-- D-XX counter: last used D-20. Next decision: D-21. -->

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
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c04-claude-claude-opus-5-5-custom-supports-limits-table.html` (new; declared as the glob index-v1c04-*.html, concrete name set by the launcher at finalize, D-17) — the supports-and-limits table page, written by the editor to the staging page and named by the launcher at finalize
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c04-claude-claude-opus-5-5-custom-supports-limits-table-reply.md` (new; declared as the glob v1c04-*-reply.md, D-17) — the editor's final reply, named by the launcher
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html` (modify, launcher only) — the iteration list between the second-opinion-tree markers gains the v1c04 row

- `ai-reports/catalog.json` (modify; integrator merge follow-up, D-16, D-17)
- `ai-reports/index.html` (modify; integrator merge follow-up, D-16, D-17)

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

## Implementation Summary

**Files changed:**
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c04-claude-claude-opus-5-5-custom-supports-limits-table.html` (created; written to the staging page by the builder, named and published by the second-opinion launcher, one CSS edit after finalize per D-12)
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c04-claude-claude-opus-5-5-custom-supports-limits-table-reply.md` (created; the editor's final reply, named by the launcher)
- `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html` (modified by the launcher only; one line added inside the iteration list)
- `ai-reports/catalog.json` (modified; regenerated by the integrator, D-16)
- `ai-reports/index.html` (modified; regenerated by the integrator, D-16)

**What was done:** Built the Tinka proposal as a fourth tree page `v1c04`: the decision above a table with one row per option (O1 to O4), a Supports column and a Limits column, and a short section below the table for the eight relations that touch no option. Every one of the 37 relations from the v1c03 map carries its `data-relation` number r01 to r37. The page has no script and no zoom or drag. The builder took the first screen of the map page (`v1c03c01`) and the table page at 1440 × 1000 and 390 × 844 and wrote them to the run folder as `REQ-734-compare-{map,table}-{1440,390}.png`; the hand-back shows them side by side, counts the relations readable on each first screen (map 0 and 0, table 4 and 0), and picks no winner.

## Discovered Tasks

(from the builder hand-back)

- The standalone re-judge (`ai-report-judge --page`) left its capture folder under the system temp folder (`judge.json` and four PNGs). The command reports it as a created output, so this may be by design. → report only
- The launcher's `<aside>` direction block carries no side padding, so every tree page that does not style `body > aside` shows it flush at x = 0 (the v1c03c01 map page has the same look). A gutter in the launcher's inserted block would fix it for every page. → report only

(from the integration and review, run-wide Restatement Sweep)

- `skills/do-work-toolbox/actions/help.md:18` reads `ai-report judge <dir>   Render-check a report bundle (widths, themes, links)` and omits the `--page <file.html>` form that REQ-728 (render check takes one named page) added; `skills/do-work-toolbox/docs/ai-report-guide.md:80` and `skills/do-work-toolbox/actions/ai-report.md:36` show it. Incomplete, not false (D-20); a fix is a shipped change, so it needs its own REQ and release. → report only (follow-up REQ)
- `skills/do-work-toolbox/crew-members/coding-guardrails.md` lacks the "Earned defense" paragraph of `skills/do-work/crew-members/coding-guardrails.md:56`; the gap dates from REQ-170 (`84c1b81d`), not from this run. → report only
- The hand-back's first line is a bare page path before its heading, and it says the map page's 390 first screen reaches the first section, which starts below the screen (review F3). Run artifact only. → report only

## Qualification

**Gate records:** `advance --diff-range d31159be..0ccdf98d` returns `success` with no finding after D-17. The first call (builder range `d31159be..c812ea4b`) flagged only scope bookkeeping: the two glob Scope entries did not match the launcher-named files, and the catalog pair was outside the builder range. No debug artifact, no stray file.

**Requirement trace** (diff `d31159be..0ccdf98d`, 5 files, 193 insertions, 1 deletion; checked in the real files at `0ccdf98d`):
- One row per option, decision above: the page opens with the decision heading "Pilot O2 by hand. Build nothing yet." and the order relations r01-r04, then one table (`<th scope="col">` Option, Supports, Limits) with four `<th scope="row">` rows O1 to O4. ✓
- Two main columns with code, short line and source per entry: each entry reads "From → To: relation words", the other box's code and title, and one source chip (D-09). ✓
- Relations that touch no option, none of 37 lost: the section "Relations that touch no option" groups the eight by Live check, S5, Decision and S9 (D-03); `data-relation` values r01-r37 all present (37 distinct, 43 entries with the D-11 repeats); the GREEN probe passes. ✓
- No zooming or dragging, reads top to bottom at both widths: no `<script>` element, no transform or canvas; the only "zoom" strings are the direction text and a link to the v1c03 "zoomable relation map" page. The builder's browser matrix shows page width 1425 at 1440 and 375 at 390 (no overflow), rows stacked at 390. ✓
- Comparison: the hand-back shows the four first-screen pictures side by side per width with both page paths and the counts (map 0/37 and 0/37, table 4/37 and 0/37), and picks no winner. ✓
- Finalize judge `pass`, every other bundle file keeps its bytes, `index.html` changes only in the iteration list: the probe pins 14 blobs and the outside-the-list hash; the merged `index.html` diff is one added `<li>`. ✓
- Pictures not in the bundle: the probe counts 17 bundle files; the four PNGs are in `do-work/runs/work-2026-11-15-012821/`. ✓

**Scope comparison:** declared (after D-17) = touched: the v1c04 page, the v1c04 reply and the bundle `index.html` (builder), plus `ai-reports/catalog.json` and `ai-reports/index.html` (integrator, D-16). Nothing under `skills/` or `suite/` changed, so no release.

## Testing

**Tests run:**
- Repository gate `bash _dev/tests/maintainer-verify.sh` (stage reuse off) at `0ccdf98d` (merge `c812ea4b` plus the catalog rebuild) in the run's detached checkout: exit 0, 124 s, 1-minute load 2.73 before start. Recorded as the green gate with `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0`.
- GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-734-probe.sh`, run by `advance --probe-file`: exit 0 (`REQ-734 probe PASS`). It checks the 14 pinned bundle blobs, `index.html` outside the iteration list, exactly one v1c04 page and reply, 17 bundle files (no picture in the bundle), the v1c04 id and base metas, the Supports and Limits headings, r01 to r37, no script and no remote stylesheet, and the iteration-list link.
- Render judge on the new page alone, `ai-report-judge ai-reports/2026-11-14_1547_proposal-tinka-review-surface --page index-v1c04-claude-claude-opus-5-5-custom-supports-limits-table.html --out <integrator scratch>/judge` on the merged tree: verdict `pass`, no findings, nothing written into the bundle. The scratch output folder (judge.json and four captures) was deleted afterwards.
- From the hand-back: finalize judge `pass`; re-judge after the D-12 CSS edit `pass`; a browser matrix at 1440 and 390, light and dark, first screen and full page: 0 overlaps, 0 clipped, 0 broken words, no text under 12 px, minimum contrast 5.51, no page overflow, rows stacked at 390, 12 Tab stops all in view.

- After the review fixes (D-18), at `13322db2`: repository gate exit 0, 122 s, load 2.49, recorded with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh`; the probe passes; `ai-report-judge --page` on the fixed page `pass`, no findings, scratch output deleted; `ai-report-index` changed nothing.

**Result:** pass.

**Red-green validation:** tdd false. RED at base `18bd280a`: the probe failed with no v1c04 page, no reply, and 15 bundle files. GREEN at the builder tip `6d50d565` and again at `0ccdf98d`.

**New/updated tests:** none in the repository; the probe is a run artifact.

**Heavy verification plan:** range `d31159be33a1498602f516e1a0c1d36d37f86507..13322db24340462b2871754cb14be84165655c84` (re-planned after the review fixes; same six lanes as at `0ccdf98d`). No lane covers `ai-reports/`, so coverage is uncertain and six lanes are selected, each argv `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane>`, reason "coverage is uncertain" for the four uncovered `ai-reports/` paths: `queue-kanban-javascript`, `queue-kanban-browser`, `do-work-cli-integrations`, `staged-skills`, `updater`, `installer`.

## Review

**Overall: 96%** | 2026-11-15T03:37:01Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Verdict:** Approve (first pass 94% on `d31159be..0ccdf98d`; delta re-review 96% on `0ccdf98d..13322db2`). Full report: `do-work/runs/work-2026-11-15-012821/REQ-734-review.md`.

**Important findings:**
- F1. The page's catalog correction ("neither Tinka proposal is in `ai-reports/catalog.json`, grep count 0") was true at the builder's base `18bd280a` and false after the merge (count 11). Fixed on the page in `e760b5b5` (D-18); the reply keeps its bytes as the editor's record — impact-user-visible → fixed

**Minor findings:**
- F2. The page's direction-block correction cited `ai-report.md:87`; the sentence is now at line 89. Fixed in `e760b5b5` — impact-negligible → fixed
- F3 (nit). The hand-back's first line is a bare path, and it says the map page's 390 first screen shows the start of the first section, which starts at y = 1023 — impact-negligible → report only
- F4 (nit, restatement sweep). `skills/do-work-toolbox/actions/help.md:18` lists `ai-report judge <dir>` without the `--page <file.html>` form REQ-728 (render check takes one named page) added (D-20) — impact-negligible → report only

**Acceptance:** Pass. The reviewer's script matched all 37 relations against v1c03 `#rev-3` (words, both end codes, kind, D-03 column; 0 mismatches) and its own browser count agrees with the builder's: map 0/37 and 0/37, table 4/37 and 0/37.
**Restatement sweep:** nothing redefined (this diff touches only export-ignored `ai-reports/` files). Wave end: inherited elements of REQ-723, REQ-724 and REQ-732 and the named candidates for REQ-728, REQ-729, REQ-730 and REQ-731 swept in `skills/` and `suite/`; one incomplete restatement, F4 (`skills/do-work-toolbox/actions/help.md:18`). Unread for the sweep: REQ-729, REQ-726, REQ-725 (their lines record no element).
**Anti-bloat:** 0 helpers, 0 options, 0 extra files, 0 decorative tests that the REQ did not name.
**Map vs table:** At 1440 the table page shows 4 of 37 relations on the first screen (the order of the options under the decision) and the map page 0; the map starts just below the first screen, shows a line's words only after a click on a box, and keeps seeds, limits and box texts behind clicks or collapsed groups, so without a click its relations are readable only in the "Every relation" table 4 to 6 screens down. The table hides nothing behind a click: all 37 relations are readable by scrolling within about 4 screens, and it leaves out the full box texts and extra source chips, which it links to `index.html`. At 390 both show 0; the table's relations stack over about 9.5 screens, while the map is a 1910-px drawing in a 331-px window that needs panning, its relation table needs a sideways swipe on most rows, and the page is about 11.5 screens long.
**Follow-ups created:** None (F1, F2 fixed; F3, F4 report only)

## Heavy Verification Plan

Base `d31159be33a1498602f516e1a0c1d36d37f86507`, target `13322db24340462b2871754cb14be84165655c84`. No lane covers `ai-reports/`, so coverage is uncertain and six lanes are selected; each argv is `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane>`, reason "coverage is uncertain" for the four uncovered `ai-reports/` paths.
- `queue-kanban-javascript`
- `queue-kanban-browser`
- `do-work-cli-integrations`
- `staged-skills`
- `updater`
- `installer`

## Heavy Verification Result

Target `13322db24340462b2871754cb14be84165655c84`, execution revision `13322db24340462b2871754cb14be84165655c84` (run's detached checkout, `QUEUE_KANBAN_BROWSER` set, load 3.24 at start). All exit 0, no findings.
- `queue-kanban-javascript`: executed, 6 s
- `queue-kanban-browser`: executed, 72 s
- `do-work-cli-integrations`: executed, 56 s
- `staged-skills`: executed, 30 s
- `updater`: executed, 66 s
- `installer`: executed, 26 s

## Lessons Learned

**What worked:** generating the page from the v1c03 page's own box and relation texts kept all 37 relations word for word, and the `data-relation` numbers made "none is lost" a script check for the probe and the reviewer. One browser script for every state, run before finalize, found the contrast and broken-word defects early.

**What didn't:** two corrections on the page were true at the builder's base and false after the merge, because siblings of the same run landed in between (REQ-730 rebuilt the catalog, REQ-728 and REQ-729 moved lines in `ai-report.md`). A dated claim about the repository needs the revision it was checked at, or it goes stale at merge. The launcher's inserted direction block also changed the published page's first screen (no side gutter), which the staging-page check could not see.

**Worth knowing:** run the browser check on the published page too, not only on the staging page (builder's proposed lesson, D-19). For the map-or-table choice, the reviewer's plain comparison is in `## Review` above.

## Orientation

Report bundles under `ai-reports/` (export-ignored, never shipped): a second-opinion tree page is published by the launcher from a staging page and listed in the bundle's iteration list and in the regenerated `ai-reports/catalog.json` pair.

## Timing

Observed 2026-11-15T00:47:40Z to 2026-11-15T03:37:02Z: 2h 49m 22s total, 2h 40m 20s attributed across 8 events, 9m 02s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 2h 08m 03s | 1 |
| review | 17m 27s | 1 |
| verification-gate | 9m 44s | 3 |
| exploration-preflight | 4m 05s | 1 |
| planning | 52s | 1 |
| handback-merge | 9s | 1 |

Slowest stage: builder-work / builder worktree build, 2h 08m 03s, outcome success.
