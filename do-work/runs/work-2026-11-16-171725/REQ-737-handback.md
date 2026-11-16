# REQ-737 Hand-back — relation-map preset keeps box titles readable after a click on a phone

`ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map.html`

Branch `worktree-agent-REQ-737-relation-map-phone-titles`, tip `8ef2084e`, base `6f00ab48` (VERSION 0.324.0). Done and verified: at 390 x 844 the smallest clicked or related title after a click is 12.32 px in 33 clicks (v1c03: 4.52 px), and every related box outside the frame is one Show button away (70 of 70 steps work).

## Commits

- `e9d30a3f` [REQ-737] relation-map preset: readable titles after a click on a phone
- `8ef2084e` [REQ-737] proof page v1c04: relation map with readable titles after a phone click

## Files

- modified `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` — check 15 gains one phone sentence; the desktop sentence is unchanged word for word.
- created `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map.html` (resolves the Scope glob `index-v1c04-*.html`) — the proof page, written by the launcher at finalize.
- created `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map-reply.md` (resolves the Scope glob `v1c04-*-reply.md`) — the editor's reply.
- modified `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` — the launcher added one `v1c04` line to the iteration list.
- Main tree, not committed: this hand-back, `REQ-737-progress.log`, and `REQ-737-captures/` (script, three raw outputs, eight PNGs).

## Scope changes

None.

## Red-green evidence

Same script, same server, 390 x 844, light theme, 33 clicks each (10 from the first view, 23 after "Show everything"):

| Page | Map scale after a click | Smallest clicked title | Smallest related title in the frame | Related titles out of view | Of those, one named step works |
| --- | --- | --- | --- | --- | --- |
| v1c03 (RED, 0.324.0 preset) | 0.205 to 0.342 | 4.52 px | 4.52 px | 0 of 87 | n/a |
| v1c04 published (GREEN) | 0.56 every click | 12.32 px | 12.32 px | 70 of 87 | 70 of 70, each at 12.32 px in view after the tap |

GREEN probe (`REQ-737-probe.sh`):
- At base: its step 2 fails by design. Checked on the base file: `git show 7f2c4b94:<preset> | grep numbered | grep click | grep -c 'phone|600 px'` prints 0.
- After commit 1: `REQ-737 probe FAIL: no relation-map page other than v1c03 in ai-reports/...` (step 4, as expected). Steps 1 to 3 passed, including the Go test.
- At the tip: `REQ-737 probe PASS (3289 bytes, 20 checks, ai-reports/.../index-v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map.html)`.

## Verification

Commands, from the worktree:
1. `go test -count=1 -v -run '^TestSecondOpinionShippedPresetsLoadFromTheToolbox$' ./internal/toolboxcommands/`: `--- PASS: TestSecondOpinionShippedPresetsLoadFromTheToolbox (0.35s)`, `ok`.
2. `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0. `bash _dev/tests/contract-regressions.sh`: exit 0, `Contract regression checks passed.` One run each, no rerun.
3. `git diff --stat main...HEAD -- skills/`: only `relation-map.md | 2 +-`. `git diff --stat main...HEAD`: that file, the v1c04 page (675+), its reply (20+) and `index.html` (1+).
4. Standalone re-judge with `--page index-v1c04-…-phone-readable-map.html --out <scratch>/judge`: exit 0, `success`, findings `[]`. Output folder deleted.
5. GREEN probe: PASS (above).

Finalize row: `state ran`, `judge_verdict pass`, `judge_findings []`, `banner_state inserted`, `tree_id v1c04`, `model_evidence session-declared`, `preset_file skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (worktree). The staging page, reply and sidecar were consumed by finalize; `git status` was clean after it.

Preset size: 3,057 bytes before, 3,289 after (target about 3,500, probe bound 3,600). 20 numbered checks, no `{{`.

Browser matrix (`measure-relation-map.js`, raw output `measure-output-staging.json`, `measure-output-v1c04.json`, `measure-output-v1c03.json`):

| Row | Staging 390 | Published 390 | Staging 1440 | Published 1440 |
| --- | --- | --- | --- | --- |
| Clicks (mouse + synthetic for boxes outside the screen) | 33 (26 + 7) | 33 (26 + 7) | 33 (33 + 0) | 33 (33 + 0) |
| Smallest clicked title, all in view | 12.32 px, 33 of 33 | 12.32 px, 33 of 33 | 9.97 px (O2), 33 of 33 | 9.97 px, 33 of 33 |
| Smallest related title in the frame | 12.32 px (17 in view) | 12.32 px (17) | 9.97 px (87 of 87) | 9.97 px (87 of 87) |
| Related out of view → named step works | 70 → 70 | 70 → 70 | 0 | 0 |
| Clicked or related box out of frame | 70 (by design, phone rule) | 70 | 0 of 33 clicks | 0 of 33 clicks |
| Sideways page scroll | 0 | 0 | 0 | 0 |
| Words on boxes / words broken / narrow labels | 0 / 0 / 0 | 0 / 0 / 0 | 0 / 0 / 0 | 0 / 0 / 0 |
| Relation words in view per click | 0 to 5 | 0 to 5 | 1 to 7 | 1 to 7 |
| First view: titles, words | 12.93 px, D O1 O2 O3 in view, 2 of 12 words | same | 12.41 px, 10 boxes, 12 of 12 words | same |
| Tab: focused box in frame, frame scroll 0 | 23 of 23 at 1.08, 0 | same | 23 of 23 at 1.56, 0 | same |
| List: rows = relations, unnamed, rows wider than screen | 28 = 28, 0, 0 | same | 28 = 28, 0, 0 | same |
| No JavaScript: tools hidden, words on their line, on boxes, boxes as text, rows | yes, 28 of 28, 0, 23, 28 | same | same | same |
| Launcher header block height | 0 (none yet) | 1,892 px | 0 | 726 px |
| Map frame top y | 801 | 2,713 | 353 | 1,099 |

Desktop 9.97 px is the 0.324.0 desktop rule at work (O2 with 11 relations fits at 0.453); check 15's desktop sentence is unchanged, and v1c03 measures the same 9.97 px.

## Captures for the maintainer

Layout form chosen: **form 1, focus and panel jumps.** On a phone a click keeps the map at 0.56 (titles 12.3 px). The frame holds the clicked box with the related boxes that fit, and it slides them toward the frame edge away from the others, so the lines and relation words toward them stay in view. A line above the map says "N of M related boxes in view; the Show buttons in the detail panel bring the others into view". Each relation in the panel has a Show button that moves the map to that box and scrolls the map onto the screen; the selection and panel stay.

Rejected:
- Form 2 (stacked related boxes): it needs a second layout engine in every generated page, and it moves boxes out of their argument columns (check 2) while selected.
- Form 3 (no map move): the related boxes would not be shown at all after a click, which drops check 8 on phones.

Files in `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-captures/` (published v1c04 unless named; each shot scrolls the map tool row to the top of the screen first, because the launcher's header fills the first 1,892 px on a phone):
- `v1c04-390-first-view.png` — phone first view: D, O1, O2, O3 at 12.9 px, 2 relation words.
- `v1c04-390-click-O2.png` — after a click on O2 (11 relations): O2 and D at 12.3 px, 5 relation words to the right, "1 of 11 related boxes in view" line, panel starts below.
- `v1c04-390-click-O1.png` — after a click on O1 (2 relations): O1 and D framed, the word "already covers routine work" toward S3.
- `v1c04-390-click-O2-then-reach-S1.png` — O2 selected, then Show S1 tapped in the panel: S1 at 12.3 px with its lines and words toward O2 at the right edge.
- `v1c04-390-show-everything-click-S6.png` — "Show everything", then a click on S6: S6 at 12.3 px with all 5 relation words toward P1 to P5.
- `v1c04-1440-first-view.png` — desktop first view, unchanged from v1c03 apart from the new count line.
- `v1c04-1440-click-O2.png` — desktop click on O2: all 11 related boxes in the frame, panel with Show buttons.
- `v1c03-390-click-O2.png` — the before page after a click on O2: scale 0.205, titles 4.5 px.

## P-A-U

- PLAN: commit 1 adds one phone sentence to check 15; commit 2 builds v1c04 through the Session form, measured RED on v1c03 and GREEN on staging and published pages with one kept script.
- APPLY: two commits as planned; one page revision round (the first phone framing showed only the clicked box, so I added the slide toward left-out boxes, then the same slide for the Show step).
- UNIFY: `git diff --stat main...HEAD` shows 4 files: `relation-map.md` (1 line), the v1c04 page, its reply, `index.html` (1 line). I read the preset hunk and the `index.html` hunk; the page diff against v1c03 is the builder script's edits only. No staging file, debug output or capture is in the diff.

## Decisions

- **D-11** (ESCALATE): Layout form 1 with a slide toward the left-out boxes, as described under Captures. Value: titles stay at 12.3 px after every phone click, relation words toward off-frame boxes stay readable, and the preset needed one sentence (232 bytes). Risk: with this map's column spacing, a 356 px frame at a readable scale holds one column, so most phone clicks show the clicked box plus same-column neighbours, and 70 of 87 related boxes need a Show tap. If the maintainer prefers form 2, it is a new build pass on this branch, and the preset sentence changes.
- **D-12** (DECIDE & STATE): Check 15 keeps its desktop sentence byte for byte and adds the phone sentence in the same check, so no check 21. Checks 8, 17 and 18 do not contradict it: related boxes are still shown (only some sit outside the frame), and the first view and table rows are untouched.
- **D-13** (DECIDE & STATE): The phone sentence says "every related title in the frame" and "a button that brings that box into view on the map". "Into view" covers the page scroll, because on a phone the panel sits below the map.
- **D-14** (DECIDE & STATE): The page reuses v1c03's content and map technique, rebuilt by a scratch script that strips the launcher's decoration and the page's own CSS rule for the direction block (D-07), then edits the click code. The facts in its corrections were re-checked against the current repository (line numbers unchanged).
- **D-15** (DECIDE & STATE): The measurement script is a Node Playwright script (global `playwright`, run with `NODE_PATH="$(npm root -g)"`), not the browser MCP tool, so the integrator can rerun it with one command. It therefore created nothing in the main tree's `.playwright-mcp/`; that folder does not exist, and nothing was deleted there.
- **D-16** (DECIDE & STATE): The script clicks a box with the real mouse when its title shows on screen and with dispatched pointer events otherwise (7 boxes in the phone first view lie outside the frame). The row records which.
- **D-17** (DECIDE & STATE): The panel button changed from "Open X" (select X) to "Show X" (move the map, keep the selection), because the preset now asks for a button that brings the box into view. Tapping X on the map still opens it.
- **D-18** (DECIDE & STATE): The model slug is `claude-opus-5-5`, recorded as `session-declared`. Commit trailers use Opus 5.5 per the brief's current text (an earlier copy said Fable 5.2).
- **D-19** (DECIDE & STATE): My first finalize call repeated `--backend/--model/--preset` and was refused as `TOOLBOX-USAGE`; `--finalize <staging-page>` alone is the accepted form. Nothing changed on disk from the refusal.

## Discovered Tasks

- The brief's finalize line ("the same launcher line with `--finalize` in place of the prepare flags") is refused by the launcher: finalize takes only the staging path — impact-negligible → report only
- With the 0.324.0 header block unstyled (D-07), the direction block is 1,892 px tall at 390 (REQ-727 measured 1,302 px with its own styling), so the phone map starts on the fourth screen until REQ-738 lands — impact-user-visible → report only

## Integration seams

1. Prepare: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-second-opinion --prepare --backend claude --model claude-opus-5-5 --preset relation-map ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1`. Finalize: `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-second-opinion --finalize ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/.claude-claude-opus-5-5-staging-1831528706.html`. Bundle `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1`, page id `v1c04`, page `index-v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map.html`, reply `v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map-reply.md`, model slug `claude-opus-5-5`.
2. Server, from the bundle folder of the tree to measure: `python3 -I -m http.server 8771 --bind 127.0.0.1`. Script: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-captures/measure-relation-map.js`. Command (zsh): `NODE_PATH="$(npm root -g)" node /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-737-captures/measure-relation-map.js http://127.0.0.1:8771/index-v1c04-claude-claude-opus-5-5-relation-map-phone-readable-map.html > measure-output-v1c04.json`; the before page is `http://127.0.0.1:8771/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`. Add `CAPTURE_DIR=<folder> CAPTURE_PREFIX=<name>` to also save the PNGs. `summary.green` per size is the GREEN verdict.
3. The GREEN numbers do not depend on the header block: staging (no block) and published (1,892 px block) give the same title sizes, steps and counts, because the script scrolls the map tool row to the top before each click and the page scrolls the frame onto the screen after a phone click. Numbers that change when the block gets shorter: `firstView.frameTopY` (2,713 at 390, 1,099 at 1440) and `firstView.headerBlockHeight` (1,892, 726).

## Proposed CHANGELOG entry

Title checked: no earlier entry uses "Phone Click" or "Readable After a Phone".

```
## 0.324.1 — Relation-Map Preset Keeps Titles Readable After a Phone Click (2026-11-16)

- On a phone, a click on a box keeps its title and the related titles in view at about 12 px, instead of zooming out to 4 to 7 px.
- The detail panel gives each related box a Show button that brings it into view on the map.
- What changes for an installed copy: a page made with `ai-report second-opinion --preset relation-map` keeps box titles readable after a click on a phone. Existing pages are not rebuilt, and nothing else changes. This removes the first limit named in the 0.324.0 entry's "Known limits" line.
```

## Proposed lesson

- [family: session-is-the-backend] A readable-scale floor on a phone can leave no room for a second column: with 300-unit boxes and 300-unit gaps, a 356 px frame at 0.56 holds one column, so "frame the related boxes that fit" framed only the clicked box until the view slid toward the left-out boxes. Look at a phone capture of the click view, not only the size numbers, before finalize.

## Lessons read

- `_dev/primes/prime-releases.md` and `_dev/primes/lessons-releases.md` in full (canonical-link sweep is the integrator's at UR-180 close).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, the two `session-is-the-backend` bullets by fixed-string grep (0.306.4 split at the generation call; 0.324.0 measure every click, then again after finalize). Both applied.

## Addendum

Second build pass, after the maintainer kept form 1 and asked for the map in the top half of the first phone screen and no relation label cut at an edge (REQ `## Addendum (mid-run)`). The new proof page is `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c05-claude-claude-opus-5-5-relation-map-phone-first-map.html`. Branch tip `4b3d6fad`. All earlier sections above describe the first pass and stay as written.

### New commits

- `b8542cb5` [REQ-737] relation-map preset: phone map starts in the top half, words stay clear of frame edge and arrowheads
- `4b3d6fad` [REQ-737] proof page v1c05: phone map in the top half of the first screen, words clear of arrowheads

### Preset

3,289 → 3,426 bytes (probe bound 3,600), 20 checks in order, no `{{`. Two checks changed, each with the rule where a generated page needs it:
- Check 12: "A relation word stays readable at every zoom, sits wholly inside the frame and never covers a box or an arrowhead, because shown boxes and line ends count as taken space."
- Check 17: "Below 600 px wide, the map frame starts in the top half of the first screen, the first view frames the main item at the frame's top, its titles read without zoom or sideways panning, and the page never scrolls sideways."

### What the page changed (v1c05 against v1c04)

- Phone only: the intro paragraph moves below the map (CSS `order` in a flex column, so the DOM order and desktop are unchanged), the change list runs inline, the title is 1.2 rem, and the zoom buttons show "+" and "−" (their aria-labels still say "Zoom in" and "Zoom out").
- All widths: the note "the lines are my reading of the report" moves from above the map to beside the key, below the map.
- Phone first view: the main boxes sit at the frame's top (the column title 12 px under the frame edge) instead of centered.
- Word placement: each drawn line's arrowhead is taken space, 7 line widths long outside the box edge plus 5 px. Words may also shift a quarter of their width along the line, which keeps desktop at 12 of 12 first-view words.

### Changed matrix rows (published pages, same script)

| Row | v1c04 390 | v1c05 390 | v1c04 1440 | v1c05 1440 |
| --- | --- | --- | --- | --- |
| Words on an arrowhead (first view + 33 clicks) | 10 (D, O2, Q1, L1, L3, Q2 clicks) | 0 | 11 (O2 twice, O3) | 0 |
| Words cut at the frame edge | 0 | 0 | 0 | 0 |
| First-view relation words | 2 of 12 | 3 of 12 | 12 of 12 | 12 of 12 |
| Relation words in view per click | 0 to 5 | 0 to 6 | 1 to 7 | 1 to 6 |
| Frame top y at load, no launcher header | 801 | 313 | 353 | 299 |
| Frame top y at load, 87.4 px short header (derived) | 888 | 400 | 440 | 375 |
| Frame top y at load, 0.324.0 header (published) | 2,713 | 2,311 | 1,099 | 1,070 |
| `summary.green` | false (arrowheads, frame top) | true | true | true |

Unchanged on v1c05: 33 clicks per width, phone titles 12.32 px after every click (clicked and in-frame related), 70 of 70 Show steps work, 0 sideways scroll, 0 words on boxes or broken, Tab 23 of 23, list 28 = 28 rows with 0 too wide, no-JavaScript view as before, desktop 0 of 33 clicks with a box out of frame. Staging and published v1c05 give the same numbers apart from the header rows.

The `v1c04 390` arrowhead count is the measured form of the cut label in `v1c04-390-click-O2.png`. "open: hold or hand back?" was not cut by the frame edge (that count was already 0). Its first letters sat under the arrowhead stack on O2's right edge, so it read as cut.

How the script measures the short-header frame top: `firstView.frameTopYWithShortHeader = frameTopY − pageTopY + SHORT_HEADER_PX`. Here `pageTopY` is where the page's own `.page` content starts at load, below whatever header the launcher printed. `SHORT_HEADER_PX` defaults to 87.4 and can be set from the environment. On a page with no header (staging) this is frameTopY + 87.4. On the published page it replaces the 0.324.0 header (1,966 px at 390) with 87.4 px, and both give 400. After REQ-738 lands, the integrator can compare it with the real `frameTopY`. The phone GREEN verdict now also requires `frameTopYWithShortHeader ≤ 422`, 0 words on arrowheads and 0 words cut at the frame edge.

Raw outputs (captures folder): `measure-output-staging-v1c05.json`, `measure-output-v1c05.json`, `measure-output-v1c04-rerun.json` (v1c04 measured with the updated script).

### Verification (second pass, from the worktree)

- `TestSecondOpinionShippedPresetsLoadFromTheToolbox`: PASS. `shipped-package-reference-contract.sh`: exit 0. `contract-regressions.sh`: exit 0.
- `git diff --stat main...HEAD -- skills/`: only `relation-map.md | 6 +++---`. Whole branch: 6 files (the preset, the v1c04 and v1c05 pages and replies, `index.html`).
- Finalize: `judge_verdict pass`, `tree_id v1c05`, `model_evidence session-declared`, `banner_state inserted`. Standalone judge `--page index-v1c05-…-phone-first-map.html`: `success`, findings `[]`, output deleted.
- GREEN probe: `REQ-737 probe PASS (3426 bytes, 20 checks, …v1c04…)`. The probe names the first relation-map page other than v1c03.
- HTTP server stopped. No `.playwright-mcp/` files were created.

### New captures (in `REQ-737-captures/`)

- `v1c05-390-load-without-launcher-header.png` — the first phone screen at load, as it looks under no header. The frame starts at y = 313, so with REQ-738's 87 px header it starts at about 400. This page body was captured from the staging copy, which matches the published page except for the header.
- `v1c05-390-first-view.png` — published page, map tool row scrolled to the top: D, O1, O2, O3 at 12.9 px at the frame's top, 3 relation words.
- `v1c05-390-click-O2.png` — published page after a click on O2: 6 relation words in view, "open: hold or hand back?" clear of the arrowheads.

### Decisions (second pass)

- **D-20** (DECIDE & STATE): The rebuilt page is a new tree page, v1c05, not a rebuild of v1c04 in place. The launcher never overwrites a page, and an in-place rebuild means deleting a published file and editing the launcher's iteration list by hand. v1c04 stays on the branch as the first-pass page. If the integrator wants only one proof page, deleting v1c04 and its `index.html` line is a two-path change; the probe passes either way.
- **D-21** (DECIDE & STATE): The two new rules go into checks 12 and 17, no check 21. The frame-top rule is about the phone first view, which check 17 owns. The arrowhead rule is a case of "never covers a box", which check 12 owns.
- **D-22** (DECIDE & STATE): The phone intro moves below the map with CSS order instead of being cut or folded, so no words leave the page and the desktop layout keeps its order.
- **D-23** (DECIDE & STATE): "Cut at the frame edge" in the addendum turned out to be the label under the arrowhead stack. Both are now measured. The edge count was already 0, and the arrowhead count went from 10 to 0 on the phone.

### Proposed CHANGELOG entry (replaces the one above)

```
## 0.324.1 — Relation-Map Preset Keeps Titles Readable After a Phone Click (2026-11-16)

- On a phone, a click on a box keeps its title and the related titles in view at about 12 px, instead of zooming out to 4 to 7 px; the detail panel's Show buttons bring the other related boxes into view.
- On a phone, the map now starts in the top half of the first screen, and relation words keep clear of arrowheads and of the map's edge.
- What changes for an installed copy: a page made with `ai-report second-opinion --preset relation-map` keeps box titles readable after a click on a phone. Existing pages are not rebuilt, and nothing else changes. This removes both limits named in the 0.324.0 entry's "Known limits" line, together with the launcher's short header.
```

## Addendum (pass 3)

The one proof page is `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html`, built with the merged 0.324.1 launcher and REQ-738's short header. Branch tip `9375501d`. On the real page at 390 x 844, the map frame's top is at y = 398 at load. That number is measured, not derived. The earlier sections and addendum describe passes 1 and 2 and stay as written. Their v1c04 and v1c05 pages are no longer on the branch.

### Commits

- `437232b4` [REQ-737] merge main (0.324.1, REQ-738 short launcher header) before the final rebuild. `git merge --no-ff main` from main tip `ab5e07c1`, with no conflicts.
- `9f4a0200` [REQ-737] drop the first-pass proof pages before the final rebuild. Removes both pages, both replies and their two `index.html` lines.
- `9375501d` [REQ-737] proof page v1c04: relation map built with the 0.324.1 launcher.

`git diff --stat main...HEAD` now shows 4 files: `relation-map.md` (3 lines changed), the new page (692+), its reply (20+) and `index.html` (1+).

### Page id

- **D-24** (DECIDE & STATE): The final page took id **v1c04** again, because the launcher hands out the next id not present in the bundle once v1c04 and v1c05 were deleted. Its file name differs from the first-pass v1c04 by its slug (`phone-first-map`, not `phone-readable-map`), so no link from an earlier commit's page names this file. Prepare (staging `.claude-claude-opus-5-5-staging-2753489393.html`) and `--finalize <staging-page>` alone were run, from the second-pass staging content rebuilt by the same scratch builder. Finalize: `judge_verdict pass`, `tree_id v1c04`, `model_evidence session-declared`, `banner_state inserted`.

### Measured numbers (published page, `measure-output-final.json`)

| Row | 390 x 844 | 1440 x 1000 |
| --- | --- | --- |
| Launcher header block height (the `aside`, without its margins) | 70 px | 70 px |
| Page content starts at y (header plus its margins) | 86 | 98 |
| **Map frame top y at load, measured** | **398** (top half: ≤ 422) | 385 |
| `frameTopYWithShortHeader` (pass-2 derivation, for comparison) | 400 | 375 |
| Clicks, smallest clicked title, all in view | 33, 12.32 px, 33 of 33 | 33, 9.97 px, 33 of 33 |
| Smallest related title in the frame | 12.32 px | 9.97 px |
| Related out of view → Show step works | 70 → 70 | 0 |
| Clicked or related box out of frame | 70 (phone rule) | 0 of 33 clicks |
| Words on arrowheads / cut at frame / on boxes / broken | 0 / 0 / 0 / 0 | 0 / 0 / 0 / 0 |
| First view: titles, relation words | 12.93 px (D, O1, O2, O3), 3 of 12 | 12.41 px (10 boxes), 12 of 12 |
| Sideways scroll | 0 | 0 |
| Tab: in frame, frame scroll 0 | 23 of 23 at 1.08 | 23 of 23 at 1.56 |
| List rows = relations, too wide | 28 = 28, 0 | 28 = 28, 0 |
| No JavaScript: tools hidden, words on lines, boxes as text, rows | yes, 28 of 28, 23, 28 | same |
| `summary.green` | true | true |

### Verification (from the worktree)

- `TestSecondOpinionShippedPresetsLoadFromTheToolbox`: PASS. `shipped-package-reference-contract.sh`: exit 0. `contract-regressions.sh`: exit 0.
- `git diff --stat main...HEAD -- skills/`: only `relation-map.md`.
- Standalone judge `--page index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html`: `success`, findings `[]`; output folder deleted.
- GREEN probe: `REQ-737 probe PASS (3426 bytes, 20 checks, ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html)`.
- HTTP server stopped. No `.playwright-mcp/` folder exists.

### Captures (in `REQ-737-captures/`)

- `final-390-load.png` — the first phone screen at load, with no scrolling. It shows the short header, the backlink, the change list, the title and the tools, and the map frame from y = 398 with D, O1 and O2 visible.
- `final-390-click-O2.png` — after a click on O2, with the map tool row at the top: titles at 12.3 px, 6 relation words, "open: hold or hand back?" clear of the arrowheads.
- `final-1440-load.png` — desktop at load: the header and the full first view, with 12 of 12 relation words.

### Discovered Tasks (pass 3)

- REQ-738's short header has no side gutter at 390 px: "Iteration v1c04" and the summary triangle start at x = 0 (`final-390-load.png`). The page may not style the block (D-07), so this is a launcher change. impact-user-visible → report only.
