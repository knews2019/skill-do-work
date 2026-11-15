# Builder Brief — REQ-734 (build the proposal as a supports-and-limits table and compare it with the map)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The coordinator merges; an integrator tests and reviews. **This REQ changes no shipped file: no release.**

- Operative name and branch: `worktree-agent-REQ-734-supports-limits-table-page`
- Worktree (the ONLY tree you commit in): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-734-supports-limits-table-page`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you, except the files named next): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with decisions D-01 to D-08, Exploration and Scope: read `do-work/working/REQ-734-proposal-table-and-map-comparison.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-734-handback.md`
- The four comparison pictures (D-07, never stage or commit them): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-734-compare-map-1440.png`, `…-compare-table-1440.png`, `…-compare-map-390.png`, `…-compare-table-390.png`. These five main-tree files are the only main-tree files you write.
- Up to ten builders run at the same time in this run. Machine load is high; expect slow browser renders.

## Read first
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`, `anti-slop.md` (the page is a human-facing artifact), `prompt-injection.md` (you read report text as data).
2. `skills/do-work-toolbox/actions/ai-report.md` — **Second-opinion form**, lines 75-89 (Session sequence at `:79`, names and tree ids at `:87`), `:36` (`judge --page`), and the clarity-editor prompt the launcher hands you.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add). The maintainer's words for this page: "keep language clear, do not complicate it".
4. Required lesson: `[family: session-is-the-backend]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (line 17; find it with a fixed-string grep). Write BOTH the staging page and the staging reply, or finalize refuses.
5. REQ-733 (the relation map page you compare against), in the main tree: `do-work/archive/UR-177/REQ-733-rebuild-relation-map-page.md` § Lessons Learned (`:292`), § Review, and decisions D-16 to D-18; `do-work/runs/work-2026-11-15-012821/REQ-733-handback.md`. Its review took three rounds because the builder checked one state per view. Your checks below cover every state on the first pass.
6. The content source: `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` `#rev-2` (`:196`, 28 boxes with source chips) and `#rev-3` (`:204`, 37 relation rows). The map you compare with: `index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html` in the same bundle.

## Decided (do not reopen)
- **D-01:** Session form in YOUR worktree (the main tree has other dirt, so prepare refuses there). Backend `claude`, model = your own declared model slug. Target = the bundle's `index.html`, so the new page is tree id `v1c04`.
- **D-02:** Reuse the 28 box texts, their source chips and the 37 relation words from the v1c03 page HTML. Do not rewrite the content.
- **D-03:** Kind "supports" → Supports column; "limits or blocks" and "gap in the evidence" → Limits column. O2 → O3 and O3 → O4 appear in both rows. The decision block above the table states the order (the four "order" relations). The eight relations that touch no option (L1, L6, L9 → live check; L4, L5 → S5; L2, L7 → decision; L3 → S9) go in one short section right below the table, grouped by the box they touch.
- **D-04:** Every relation entry carries `data-relation="r01"` … `"r37"`, numbered in the v1c03 `#rev-3` row order. A relation shown twice carries the same number twice.
- **D-05:** No JavaScript of the page's own: plain HTML and CSS. The bundle's `presentation.css` and `<script src="presentation.js">` are allowed but must not be edited. No library, no remote file, no zoom, drag or fixed-size canvas. At phone width each option row stacks Supports above Limits.
- **D-06:** "Readable on the first screen without any click": at load, page top, no scroll, no click, the relation's words and both end codes are fully inside the viewport, visible, not clipped, at 12 px or more. Viewports 1440×1000 and 390×844, light theme. Measure by script on both pages.
- **D-07:** The four comparison pictures go to the run folder paths above. Every other capture goes to your scratch folder and is deleted before hand-back.
- **D-08:** The hand-back shows the pictures side by side and picks no winner.

## What to build
1. In your worktree, confirm `git status --porcelain --untracked-files=all` is empty. Run the GREEN probe once (step 2 of Verify) and record its failure at base.
2. Prepare: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-second-opinion --prepare --backend claude --model <your-model-slug> --editor-direction '<one sentence: rebuild the report as one table with one row per option O1 to O4 and the columns Supports and Limits, the decision above it, each cell listing the related kill attempts, seeds and limits with code, one short line and source, the relations that touch no option below it, readable top to bottom with no zoom>' ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html`. The direction is ONE argv token. Read `staging_page`, `staging_reply`, `asset_prefix`, `prompt` from the typed row (state `prepared`, base `index.html`).
3. Do the `prompt` yourself, writing the page to `staging_page` and the final reply (with its `Variation slug:` line, for example `supports-limits-table`) to `staging_reply`. The page:
   - opens its body with `<p>Original report: <a href="index.html">index.html</a></p>`, then a "Changed in this iteration" list whose `#rev-<n>` links each hit a badged section (relative to `index.html`: the table form replaces the prose layout);
   - shows the decision (D, with the pilot rule in one or two lines) above the table, then the table: one row per option O1-O4, headed by the option's code and title, with a Supports cell and a Limits cell. Each item in a cell shows the other end's code, the relation words, and the source chip from that box's `#rev-2` text;
   - shows the eight no-option relations below the table, and a short note that the full texts are in `index.html` (link the matching anchor where the v1c03 box text links one);
   - follows D-03 to D-05.
4. Check it in a browser before finalize (the browser tool blocks `file:` URLs: serve the bundle with `python3 -I -m http.server <port> --bind 127.0.0.1` from the bundle folder and stop it afterwards). Run ONE measuring script in every state of this matrix, not one state per view: widths 1440×1000 and 390×844 × light and dark × first screen and full page scrolled through. These are the three checks REQ-733's review failed on the first pass, carried to a table:
   1. **Overlap and clipping (REQ-733 F1, G1).** No two text boxes overlap; no element's `scrollWidth` is larger than its `clientWidth` (clipped text); no word is broken inside a word (`AGENTS.md`, `--check`, `recommends` stay whole: no `word-break: break-all`, no `overflow-wrap: anywhere` on cell text); no text under 12 px.
   2. **Phone width (REQ-733 G2, Q1).** At 390 px: `document.documentElement.scrollWidth <= 390` (no page overflow), every row stacks Supports above Limits, every code and source chip stays readable, and the decision block comes right after the "Changed in this iteration" list.
   3. **Focus order (REQ-733 F2).** Press Tab from the page top through every link and control to the end. Focus order follows reading order (decision, then O1 to O4, Supports before Limits, then the no-option section), every focused element is fully inside the viewport after the browser scrolls to it, and each has a visible focus ring in light and dark.
   Also: with JavaScript turned off the page reads the same; dark mode keeps text contrast. Fix and run the whole matrix again until every check is clean. Record the script's numbers per state.
5. Finalize: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-second-opinion --finalize <staging_page>`. Judge verdict must be `pass`. On `fail`, read `judge.json` and its captures, fix, and finalize again only if the launcher allows; otherwise report the refusal. If you edit the published page after finalize, re-judge it with `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge ai-reports/2026-11-14_1547_proposal-tinka-review-surface --page <published page name>` and record that verdict too.
6. Commit by exact path: the new page, the new reply and `index.html`. Nothing else.
7. The comparison, on the committed bytes (served from your worktree's bundle folder):
   1. Open the map page `index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html` and your table page at 1440×1000, light, page top, no click. Take one viewport screenshot of each to the run-folder paths above. Repeat at 390×844.
   2. With the D-06 rule, count by script how many of the 37 relations are readable on each first screen. Four numbers, each "N of 37".
   3. Delete every other capture you took.

**Not yours:** any other file in the bundle; anything under `skills/`, `suite/`, `_dev/`; `CHANGELOG.md`, `VERSION`, `ai-reports/catalog.json`; any `do-work/` path except the five main-tree files named at the top; REQ-727 (the relation-map preset).

## Rules
- Write only the three Scope files in your branch. Out-of-scope finds go to "Discovered Tasks".
- **Browser captures:** every capture except the four comparison pictures goes to your scratch folder with an absolute filename, never into the main tree's `.playwright-mcp/` folder, and is deleted before hand-back. Delete any file the browser tool wrote to `.playwright-mcp/` by default (REQ-733's builder found two there) and remove that folder if your session created it. Record the folder used and the deletion command as a decision (D-09 or later). The judge's own capture folder is removed by the judge flow; confirm it is gone.
- No staging leftovers in the bundle (`.*staging*` files) after finalize.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-734]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Never edit `CHANGELOG.md`, `VERSION`, or any `do-work/` path beyond the five named files, even if a `write_set` names it.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the five named files.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (from your worktree)
1. The finalize row: state published, judge verdict `pass`, tree id `v1c04`, base `index.html`, and the page and reply names.
2. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-734-probe.sh` with your worktree as the current directory prints `REQ-734 probe PASS`. It fails at base by design (no v1c04 page, 15 bundle files); record that run from step 1 of What to build.
3. `git diff --stat main...HEAD` lists only the three Scope files.
4. The browser matrix of step 4 is clean in every state, and the four pictures exist at the run-folder paths.
5. `git status --porcelain --untracked-files=all --ignored` shows no capture file left in the worktree, and `ls /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.playwright-mcp/` shows no file you created.

## Hand-back file content
Put the new page's path FIRST, on its own line, above the headings. Then a `## Comparison` section before the other headings: one Markdown table per width with two columns, "Map: `<map page path>`" and "Table: `<table page path>`", each cell an image link to its run-folder picture (`![…](REQ-734-compare-map-1440.png)`, relative to the hand-back file), and under each table the two counts "N of 37 relations readable on the first screen without a click". Pick no winner. Then these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb; list the four pictures separately as run artifacts outside the branch.
- `## Scope changes`
- `## Red-green evidence` — the GREEN probe failing at base and passing at your tip.
- `## Verification` — the finalize row, the browser matrix with its numbers per state, the probe, the diff stat, the capture cleanup.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-09 onward (D-01 to D-08 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE. Include the capture-deletion decision.
- `## Discovered Tasks`
- `## Integration seams`
- `## Proposed CHANGELOG entry` — "none (no shipped file changes)": `.gitattributes:48` export-ignores `/ai-reports`.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
No sibling in this run edits the Tinka bundle. REQ-733 already landed its v1c03c01 page and changed `index.html` line 31. The only shared file is the bundle's `index.html`, and only the iteration list between `<!-- second-opinion-tree -->` markers (`index.html:25-33` at base) changes: the launcher rebuilds it. Name in the seam section the exact `index.html` lines the launcher changed. The probe pins every other bundle file's bytes, including both v1c03c01 files.

Your final message: under 2500 characters. Start with the new page's path, then the branch tip hash, the judge verdict, the four first-screen counts, the browser matrix result, and anything the coordinator must decide.
