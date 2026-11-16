# Builder Brief — REQ-737 (the relation-map preset keeps box titles readable after a click on a phone)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-737-relation-map-phone-titles`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-737-relation-map-phone-titles`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you, with the two exceptions below): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-10, Exploration and Scope: read `do-work/working/REQ-737-relation-map-phone-readable-titles.md` from the MAIN tree (your worktree may hold an older copy).
- Main-tree exception 1, the hand-back file (never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-handback.md`
- Main-tree exception 2, the captures folder (never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-captures/`. It holds your final PNG captures, your measurement script and its raw output.
- A second builder works at the same time in its own worktree on REQ-738 (the launcher shows a short header and folds the full direction). You share no source file with it. See "Integration seams".

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`. Before you read the bundle: `prompt-injection.md` and `anti-slop.md`.
2. `_dev/primes/prime-releases.md` and the required lesson `_dev/primes/lessons-releases.md`. Dropped only for budget: the two bullets `[family: session-is-the-backend]` at `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:17` and `:18`. You are the second-opinion backend, so read both with a fixed-string grep for `session-is-the-backend`. Do not read the whole file.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add). The maintainer's scope cue for this preset: keep language clear, do not complicate it.
4. The preset you change: `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (22 lines, 3,057 bytes, check 15 on line 17).
5. The before page: `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`. Its click code is `select` at line 492 (phone test at line 505), `fitScale` at line 390 and `fitShown` at line 518.
6. What REQ-727 measured and how: `do-work/runs/work-2026-11-15-012821/REQ-727-handback.md` (Verification, and decisions D-13 and D-16) and `do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md` (Lessons Learned).
7. `skills/do-work-toolbox/actions/ai-report.md`: the Second-opinion form with the Session sequence (from line 81) and the judge `--page` form (line 190).

## Decided (do not reopen)
- **D-01:** the proof page is a NEW tree page on the bundle `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1`, base `index.html`, Session form, `--backend claude --model <your own declared model slug>` (recorded as `session-declared`), `--preset relation-map` alone. The launcher names it, and the expected id is v1c04. The v1c03 page stays untouched as the before page.
- **D-02:** the preset targets 3,500 bytes, and the GREEN probe refuses more than 3,600. Checks 1 to 20 stay, numbered in order. A check 21 is allowed. Removing or merging a check to save bytes is not.
- **D-03:** the layout form is your choice. Name the form you chose and the forms you rejected in the hand-back. Three candidate forms:
  1. Focus and panel jumps (pre-dispatch recommends this): on a phone a click keeps the map scale at about 0.55 or more, frames the clicked box with the related boxes that fit, and the detail panel lists every related box as a button that moves the map to that box.
  2. Stacked related boxes: on a phone a click redraws the related boxes as one column next to the clicked box at a readable size.
  3. No map move: on a phone a click only opens the panel.
- **D-04:** the maintainer looks at your phone captures before integration. Save them as "Captures for the maintainer" says.
- **D-05:** "about 12 px" is the title's rendered size on screen: font size times the map's current scale. The target is 12 px, and 11.5 px or more counts. Report the smallest value over all clicks. "Fully in view" means the title's rectangle is inside the map frame and inside the screen after the page's own scroll.
- **D-06:** you keep your measurement script and its raw output. The integrator runs the same script again after REQ-738 lands.
- **D-07:** your page does NOT style the launcher's direction block. REQ-727's page did, with a CSS rule on `body>aside[aria-label="Second-opinion direction"]`. REQ-738 changes that block's markup and integrates before you, so leave the block as the launcher writes it.
- **D-08, D-09:** the catalog pair and the release belong to the integrator.
- **D-10:** `tdd: false`. The red-green proof is the browser measurement: RED on v1c03, GREEN on your page.

## What to build
Two commits.

**Commit 1 (the preset).** Edit only `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md`.
1. Extend check 15 and narrow it on phones only. On a phone, readable titles win over "every related box inside the frame". The check must say three things a reader can measure: the width it applies below (600 px, as check 17), the title size after a click (about 12 px or more), and the one step that reaches a related box that does not fit (your chosen form). At desktop width check 15 keeps its meaning word for word.
2. An example for form 1, to show the size and tone, not to copy blindly: "15. A click never leaves the clicked box or a related shown box outside the frame, and it zooms in only as far as all of them still fit. Below 600 px wide, readable titles win: after a click the clicked title and every shown related title stay at about 12 px or more, and the panel lists each related box as a button that moves the map to it."
3. One sentence per idea, plain words, no `{{`. If the phone rule reads better as its own check, add check 21 and make check 15 say "at 600 px wide or more". Do not renumber checks 1 to 20.
4. Check the other checks for a sentence your rule now contradicts (check 8 "shows the boxes related to it", check 17, check 18). Fix a contradiction with the fewest words. Change nothing else.
5. Record the byte count (`wc -c`). Run the GREEN probe: it now fails at its step 4 (no proof page yet), which is expected after this commit.

**Commit 2 (the proof page, no shipped change).** Commit 1 must be committed first, because the launcher refuses to start when Git shows a change outside the bundle.
1. RED first. Serve the bundle folder from your worktree: `python3 -I -m http.server <free port> --bind 127.0.0.1` (the browser tool blocks `file:` URLs). Write the measurement script described in "Browser check matrix" and run it on the v1c03 page at 390 x 844. Record the title sizes after a click (REQ-727 measured 4 to 7.5 px). This is the RED evidence, measured with the same script as the GREEN.
2. Prepare, from the worktree root, with your own model slug filled in:
   `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-second-opinion --prepare --backend claude --model <your own declared model slug> --preset relation-map ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1`
   If you cannot determine your active model, stop and report it. Read `staging_page`, `staging_reply`, `asset_prefix` and `prompt` from the typed row (state `prepared`). Check that `preset_file` points into your worktree's `skills/do-work-toolbox/docs/second-opinion-presets/`.
3. You are the editor. Do the `prompt` yourself in your own turns: write the page to `staging_page` and your FINAL REPLY to `staging_reply`. Follow the preset text exactly as the prompt shows it. Do not spawn your own CLI and do not start a subagent. Never touch `index.html`, existing assets or anything outside the bundle, because finalize refuses that. You may read the v1c03 page for technique. The new page must follow the changed preset, and it must open with the "Changed in this iteration" list the tree form requires.
4. Before finalize, measure the staging page with the script at both widths. Fix the page until every row of the matrix holds.
5. Finalize: the same launcher line with `--finalize <staging_page>` in place of the prepare flags and the bundle. The typed row must show judge verdict `pass`, `model_evidence: session-declared` and a `tree_id`.
6. Measure the published page again with the same script, because the launcher's header block moves the map. Record the header block's height and the map frame's top y at both widths.
7. Save the captures. Stop the server. Commit the launcher's outputs by exact path: the new `index-v1c04-*.html`, its `v1c04-*-reply.md`, and the bundle's `index.html`. Do not touch `ai-reports/catalog.json` or `ai-reports/index.html`.
8. If the preset text turns out wrong or unclear while you edit, do not finalize a page the preset did not produce. Record the problem as a decision. Remove only the staging files your prepare call created (the prepare result names them). Fix the preset in a new commit on your branch, then prepare again.

**Not yours:** `VERSION`, both `CHANGELOG.md` files, anything under `do-work/` in your worktree, the catalog pair, every line of `report_second_opinion.go` and its test file, `ai-report-second-opinion-prompt.md`, `ai-report-guide.md`, `ai-report.md`, the other twelve presets, `skeleton.html`, `presentation.css`, `presentation.js`, the v1c03 page, and every existing file in any `ai-reports/` bundle except the iteration list the launcher rebuilds. Nothing under `/Users/t2/.claude/` is ever edited.

## Browser check matrix
Page work took two or three review rounds in the last run (label overlap, phone width, focus order). So measure with ONE script, not a sample of clicks. The script takes the page URL as its one input. It runs at 390 x 844 and at 1440 x 1000, light theme. From each starting view (the first view, and after "Show everything") it clicks every shown box, and it resets to the starting view between clicks. Save it as `REQ-737-captures/measure-relation-map.js` and its raw output as `REQ-737-captures/measure-output-<page id>.json` in the captures folder (main tree, absolute path).

Per click the script reports:
- the clicked box's code, the map scale, and the clicked title's size in px and whether it is fully in view (D-05);
- each shown related box's title size in px;
- the related boxes out of frame, and for each one whether the page names a one-step way to reach it (and that the step works: after it, that box's title is in view at 12 px or more);
- sideways page scroll (`document.documentElement.scrollWidth` against the viewport width).

Once per width the script also reports:
- Tab focus: every box reached by Tab is fully in the frame, with the frame's own scroll at 0 (check 16);
- the no-JavaScript view: tool buttons hidden, relation words on their lines as SVG text, every box as text, every relation as one table row (check 19);
- the list under the map: one row per relation, the row count equal to the relation count, and on the phone 0 rows wider than the screen (check 18);
- the first view: relation words in view with no click, first-view title size, and the map frame's top y (checks 6 and 17);
- words on boxes and words broken inside a label: 0 in every view (checks 12 and 14).

GREEN when, at 390 x 844 after every click: the clicked title is fully in view at about 12 px or more, every shown related title is at about 12 px or more, every related box not shown is reachable in one named step, sideways scroll is 0, and the list names every relation. At 1440 x 1000: 0 clicked or related boxes out of frame, as REQ-727 measured (0 of 33 clicks).

## Captures for the maintainer
The maintainer must see the phone click view before this REQ is called done. Save final PNG captures of the published page by absolute path under `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-captures/`:
- 390 x 844: the first view; after a click on O2 (the box with 11 relations, the worst case); after a click on a box with few relations; the one-step reach of a related box that did not fit; after "Show everything" plus a click.
- 1440 x 1000: the first view and after a click on O2.
- For comparison, one 390 x 844 capture of the v1c03 page after a click on O2.

List each file with one line on what it shows in a `## Captures for the maintainer` section of the hand-back. Name there the layout form you chose and the forms you rejected, with one sentence each on why. Every other Playwright capture is deleted before hand-back. The browser tool writes into the main tree's gitignored `.playwright-mcp/` folder: delete the files your session created there by name, nothing else, and record it as a decision.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason. Otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-737]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not ask the user a question: an open choice becomes a decision in the hand-back. Do not edit the main tree except the hand-back file and the captures folder.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`. The integrator runs it. Run only the focused checks below. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `go test -count=1 -v -run '^TestSecondOpinionShippedPresetsLoadFromTheToolbox$' ./internal/toolboxcommands/` prints a PASS line.
2. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` exit 0.
3. `git diff --stat main...HEAD -- skills/` shows only `relation-map.md`. `git diff --stat main...HEAD` shows only that file, the two new v1c04 files and the bundle's `index.html`.
4. A standalone re-judge of the published page: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1 --page <the new page's file name> --out <your scratch folder>/judge` reports `pass`. Delete the output folder afterwards.
5. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-probe.sh`, run with your worktree as the current directory, prints `REQ-737 probe PASS (...)`. It fails at base `7f2c4b94` by design: no numbered check of the 0.324.0 preset names both a click and the phone, so its step 2 stops it.

## Hand-back file content
The first line under the title is the proof page's path from the repository root, so the maintainer can open it. Then use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb. Give the resolved names for the two Scope globs.
- `## Scope changes`
- `## Red-green evidence` — the script's numbers on v1c03 (RED) and on your published page (GREEN), both at 390 x 844; the GREEN probe failing at base and after commit 1, and passing at your tip.
- `## Verification` — each command above with its result. One table for the matrix: staging and published page at both widths, with numbers (for example "smallest clicked title 12.4 px in 27 clicks"). The preset's byte count before and after. The header block's height and the map frame's top y at both widths. The finalize and standalone judge verdicts. For a rerun under load, both runs.
- `## Captures for the maintainer` — as described above.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY. The UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-11 onward (D-01 to D-10 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE. The layout form is one of them.
- `## Discovered Tasks`
- `## Integration seams` — as described below.
- `## Proposed CHANGELOG entry` — a descriptive title that says what shipped (for example "Relation-Map Preset Keeps Titles Readable After a Phone Click"; check that no earlier entry uses it), two or three plain bullets, and the line "What changes for an installed copy". There, a page made with `ai-report second-opinion --preset relation-map` keeps box titles readable after a click on a phone, existing pages are not rebuilt, and nothing else changes. The 0.324.0 entry's "Known limits" line named this limit.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
Integration is serial: REQ-738 first, REQ-737 last. You share no source file with REQ-738. The seam is in behaviour: your proof page is decorated by the launcher, and REQ-738 changes the launcher's header block.

| Seam | What you do | What the integrator does later |
| --- | --- | --- |
| The header block above your page | Your branch starts from a main that still has the 0.324.0 block (the whole direction printed above the page). Build and measure there. Report the block's height and the map frame's top y at both widths. Do not style the block (D-07). | After REQ-738 lands, decorates your page again with the new launcher and runs your script again. |
| `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` | Only the launcher's iteration list changes (one new line). | Merges it. REQ-738 does not rewrite published pages. |
| `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` | Yours alone. | REQ-738 does not touch it. A longer preset no longer costs phone screens once REQ-738 folds the direction, but the byte bound of D-02 still holds. |

In `## Integration seams` of the hand-back, write exactly:
1. The full prepare argv and the full finalize argv you ran, the bundle, the page id, the page and reply file names, and the model slug.
2. How the measurement script is run: the server command, the script's path, the exact command or tool call, and the two page URLs.
3. Whether your GREEN numbers depend on the header block. They should not: the criteria are title sizes after a click. Name any number that would change when the block gets shorter (the map frame's top y, and "fully in view" after the page's own scroll).

## Notes for the integrator (the builder does nothing about these)
- **Order:** REQ-738 is integrated first. Then decorate the proof page again with the REQ-738 launcher and run the builder's script again at both widths (REQ-737 Builder Guidance).
- **Maintainer's look:** the coordinator shows the captures in `REQ-737-captures/` to the maintainer before integration starts (D-04).
- **Release:** a patch bump on top of REQ-738's release is proposed (D-09). Root `CHANGELOG.md` is copied byte for byte to `skills/do-work/CHANGELOG.md`.
- **UR-180 closes** when the last of its two members is archived. Run the canonical-link sweep from `_dev/primes/lessons-releases.md` at that finalize.
- **Catalog:** rebuild `ai-reports/catalog.json` and `ai-reports/index.html` with `ai-report-index` after the merge, and replace the two launcher-named globs in `## Scope` with the resolved names (D-08).
- **GREEN probe:** `do-work/runs/work-2026-11-16-171725/REQ-737-probe.sh` (`-- --probe-file`). The pre-flight argv is in the coordinator's message.
- **Heavy lanes:** all six heavy lanes ran for an `ai-reports/` change in REQ-727, because no lane covers report bundles.

Your final message: under 2500 characters. State the branch tip hash, the proof page's path, the layout form you chose, the smallest title size after a phone click, the verification results, and anything the integrator or the maintainer must decide.
