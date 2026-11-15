ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html

# REQ-733 hand-back: rebuild the relation map page (left to right, nine boxes first, readable labels)

Branch `worktree-agent-REQ-733-relation-map-left-to-right`, tip `ed6a3a8a`. Judge verdict `pass`. No shipped file changed, so there is no release.

## Commits
- `ed6a3a8a` [REQ-733] relation map page v1c03c01: left to right, nine boxes first, readable labels

## Files
All three files are in `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/`:
- `index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html`: created. The editor wrote it to the staging page, and the launcher named it, added its metas and its header block.
- `v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map-reply.md`: created. The editor's final reply, named by the launcher.
- `index.html`: modified by the launcher only. It gained one row in the iteration list.

## Scope changes
None. Only the three Scope files changed. I wrote no other file in the worktree.

## Red-green evidence
- **Probe RED at base** (worktree at `1de00505`, before prepare): `REQ-733-probe.sh` printed `FAIL: expected one index-v1c03c01-*.html, found 0` and `FAIL: expected one v1c03c01-*-reply.md, found 0`, exit 1.
- **Probe GREEN at tip** (`ed6a3a8a`, run before the commit on the same bytes): `REQ-733 probe PASS`, exit 0.
- **RED view.** I opened the v1c03 page at 1440×1000 (light), pressed "Show everything" and clicked O2. The zoom was 36%. A far-zoom title measured about 10 px on screen, a relation label about 6 px (17 px SVG text × 0.36), and unrelated boxes were at opacity 0.3.
- **GREEN view.** On the new page at the same size, the first view shows 9 of 28 boxes at 88%. Titles are 15.8 px on screen. After "Show everything" (46%), titles are 12.4 px. Relation labels are 13 px at every zoom. Unrelated boxes are at opacity 0.72.

## Verification
- **Finalize row:** outcome `success`, state `ran`, judge verdict `pass`, `judge_findings: []`, `outside_paths: []`, tree id `v1c03c01`, banner `inserted`, `model_evidence: session-declared`. The base page is the v1c03 page. The page and reply names are the ones listed under Files.
- **Browser check.** I served the bundle on 127.0.0.1 because the browser tool blocks `file:` URLs. The page was the staging page, which has the same content as the published page. At each step a script measured box overlaps, label overlaps, labels outside the viewport, clipped titles and horizontal page overflow. I also looked at the captures myself.
  - At 1440×1000, light and dark: the first view shows the nine boxes at 88% with readable titles and summaries.
  - Click on O2: the view goes to 58%, 17 boxes are shown (the live check and S1-S7 appear on the left), and 13 labels sit beside the related boxes. No label covers O2, and K3, O1 and O4 dim but stay readable.
  - Second click on O2: back to 9 boxes.
  - "Nine boxes", then K1: 3 labels, "is the only option that removes K1", "a picture hides 38 to 46% of the text" and "is not answered by". Nothing new is revealed, because K1 has no evidence or limit relations.
  - "Show everything", then O3: 28 boxes at 46%, 9 labels, all readable.
  - In every step: 0 box overlaps, 0 label overlaps, 0 labels outside the frame, 0 clipped titles, no page overflow. The titles of all 28 boxes were also checked at 46%, 60%, 78%, 101%, 132%, 171% and 223%: none is clipped.
  - At 390×844 (phone), light and dark: the map opens on the decision and options column at 73% (titles 19.7 px). After a click the view never zooms below 42% and centers on the selected box. It has 0 overlaps of boxes or labels, but the O2 view is crowded with labels (see Q1 in the reply).
  - Also checked: the +, − and 0 keys, Tab across 12 boxes (the frame stays steady), the panel's "Show/Hide its evidence and limits" button, Escape, the detail link from L1 to P, and a cold load of `#box-L1` (it opens O2 and P to reach L1).
  - No-JavaScript view: all 28 boxes and all 37 lines are shown, the tool bar is hidden, the map scrolls, the table has 37 rows and the text section has 28 details.
- **Probe:** `REQ-733 probe PASS`.
- **Diff stat:** `git diff --stat main...HEAD` lists 3 files with 587 insertions: the new page (566 lines), `index.html` (+1) and the reply (+20).
- **Capture cleanup:**
  - All browser captures went to the scratch folder `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/build-REQ-733/` and were deleted with `rm -f "$S"/*.png`. 0 PNG files remain.
  - The browser tool wrote 2 files to the main tree's `.playwright-mcp/` on the first navigation, a console log and a snapshot. That folder did not exist before (mtime 01:43, the time of my first navigation). I deleted both files by name and removed the empty folder. `ls .playwright-mcp/` now reports "No such file or directory".
  - The worktree `git status --porcelain --untracked-files=all --ignored` is empty.
  - The judge's folders (`ai-report-second-opinion-judge-*`, `ai-report-judge-profile-*`) are gone from `$TMPDIR`. Its code removes them with `defer os.RemoveAll`.
  - No `.*staging*` file is left in the bundle.
  - The local server is stopped, and the browser is closed.

## P-A-U
- PLAN: generate the page from the v1c03 page's own HTML (28 boxes, 37 relations, the text, table and after-report sections reused byte for byte) with a scratch Python generator. Put the boxes in four columns, draw labels in an unscaled overlay, and show a "core plus expanded" set of boxes.
- APPLY: Session prepare in the worktree, staging page and reply written, browser check repeated until clean, then finalize (judge pass) and one commit by exact path.
- UNIFY: `git diff --stat main...HEAD` shows only the three Scope files. I read the published page's head metas (id `v1c03c01`, base = the v1c03 page), the reply's first line, and the `index.html` diff, which is one `<li>` inside the tree markers. No staging leftover, no capture, no remote script or stylesheet (probe check 5).

## Decisions
- **D-06** (DECIDE & STATE): Browser captures went only to the scratch folder above and were deleted (`rm -f "$S"/*.png`). The browser tool's automatic `.playwright-mcp/` files in the main tree (2 files, created by my session) were deleted by name, and the empty folder was removed. Nothing was saved in the worktree.
- **D-07** (DECIDE & STATE): The browser tool blocks `file:` URLs, so I served the bundle with `python3 -I -m http.server 47733 --bind 127.0.0.1` from the bundle folder during the check, then stopped it. It wrote nothing.
- **D-08** (DECIDE & STATE): Which boxes a click reveals comes from the relations, not from a hand list. A click shows the boxes that relate to the clicked box and are not core, for example O2 shows P and S1-S7. A revealed box can be clicked in turn (P shows L1, L6, L9). A link to a hidden box opens the shortest path to it from the nine.
- **D-09** (DECIDE & STATE): The four group toggles of v1c03 (Options, Kill attempts, Seeds, Limits) are replaced by "Nine boxes" and "Show everything", plus a "Showing N of 28 boxes" count. The click-to-reveal model does the toggles' job, and two ways to hide boxes would conflict. The zoom buttons and all keys are kept. The 0 key now fits the shown boxes.
- **D-10** (DECIDE & STATE): Labels are an HTML overlay, not counter-scaled SVG text. A label shows only for the selected box's lines and sits beside the other box of the line (a line between two stacked boxes puts its label in the gap left of the column). Overlapping labels are pushed down. With all 37 labels shown at once they could not be readable without overlap.
- **D-11** (DECIDE & STATE): On a screen narrower than 600 px, the first view frames the decision and options column. The view never zooms itself out below 42%. This answers v1c03's Q3: the nine boxes at the 42% fit would be about 11 px titles in 120 px boxes.
- **D-12** (DECIDE & STATE): The model slug is `claude-opus-5-5`, this session's declared model, recorded as `session-declared`. The launcher's date meta reads `2026-11-14`, which is the UTC date at finalize.
- **D-13** (DECIDE & STATE): Box summaries and the status text (for example "verified in code") show from 75% zoom and up. Below that, a box shows its code and a title of up to 3 lines. All boxes are 300×160 world px, so no title or summary is clipped at any zoom.

## Discovered Tasks
- For REQ-727 (the relation-map preset): a preset needs a rule for any report, saying which kinds go in which column and which boxes form the first view. This page hard-codes both (reply Q4). Columns on a phone stay crowded when a box has many relations (reply Q1).
- The browser tool writes `.playwright-mcp/` files into the main tree on every navigation unless each call names a file. Parallel builders share that folder.

## Integration seams
- The only shared file is `index.html`. The launcher inserted one line, line 31: `<li style="margin-left:1.5em"><a href="index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html">v1c03c01</a> · custom · claude claude-opus-5-5 · 2026-11-14</li>`. The iteration list now runs from line 25 (`<body><!-- second-opinion-tree -->`) to line 33 (`<!-- /second-opinion-tree -->`).
- REQ-734 builds on this bundle after REQ-733 lands. The launcher rebuilds that block again, and nothing outside it changes (probe check 2).

## Proposed CHANGELOG entry
None. No shipped file changed: `.gitattributes:48` export-ignores `/ai-reports`.

## Proposed lesson
- When a map moves by CSS transform inside an `overflow:hidden` viewport, pin the viewport's `scrollLeft` and `scrollTop` to 0. Focus on an off-screen box (Tab, or a test tool's click) otherwise scrolls the clipped viewport, and the map goes blank. The v1c03 page has the same latent bug.

## Lessons read
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:17` `[family: session-is-the-backend]`: I wrote both `staging_page` and `staging_reply` before finalize, and finalize did not refuse.
