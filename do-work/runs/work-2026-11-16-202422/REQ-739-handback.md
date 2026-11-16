# REQ-739 hand-back (the second-opinion header keeps a side margin on a phone)

Branch `worktree-agent-REQ-739-header-phone-side-margin`, tip `6584e473`, base `e71b4aa2` (VERSION 0.324.2). Written 2026-11-16T20:36Z.

## Commits
- `1e478a2a` [REQ-739] RED: pin a zero-specificity side margin on the second-opinion header
- `6584e473` [REQ-739] GREEN: the second-opinion header keeps a 16 px side margin

## Files
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`: new `TestSecondOpinionHeaderKeepsAPhoneSideMargin` (5 checks), placed right after `TestSecondOpinionHeaderFoldsTheFullDirection`.
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go`: the `block :=` string of `secondOpinionDecorate` gains `<style>:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}</style>` right before the `aside`; the function comment names the rule and why it is zero-specificity.
- Main tree, not committed: this hand-back and `REQ-739-captures/` (5 PNG files).

## Scope changes
none. No existing assert changed.

## Red-green evidence
- RED after commit 1: `report_second_opinion_test.go:1091: check 1: the header block has no zero-specificity side-margin rule, so its text touches a phone's edge: <aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere"><p><strong>Iteration v1c03</strong>…` then `--- FAIL: TestSecondOpinionHeaderKeepsAPhoneSideMargin (0.00s)` (assertion failure, not a compile error).
- GREEN after commit 2: `--- PASS: TestSecondOpinionHeaderKeepsAPhoneSideMargin (0.00s)` and `--- PASS: TestSecondOpinionHeaderFoldsTheFullDirection (0.00s)`.
- GREEN probe at base (`e71b4aa2`, test file stashed): `REQ-739 probe FAIL: report_second_opinion_test.go has no TestSecondOpinionHeaderKeepsAPhoneSideMargin`. At the RED commit: FAIL on the new test's check 1. At tip `6584e473`: `REQ-739 probe PASS (4 header tests, aside tag unchanged)`.

## Verification
1. `gofmt -l` on both Go files: empty. `go vet ./internal/toolboxcommands/`: clean. `go test -count=1 ./internal/toolboxcommands/`: `ok … 23.901s` (whole package, after the final comment edit).
2. `bash _dev/tests/shipped-package-reference-contract.sh`: PASS (1.5 s). `bash _dev/tests/contract-regressions.sh`: "Contract regression checks passed." (26 s). `bash _dev/tests/quiet-grep-pipeline-audit.sh`: passed (245 tracked shell files).
3. `git diff --stat main...HEAD`: only `report_second_opinion.go` (8 lines, +5 −3) and `report_second_opinion_test.go` (+33). `git status --short`: clean (the scratch test file `zz_req739_scratch_test.go` is deleted; it was never committed; no folder under `ai-reports/`).
4. `git diff main...HEAD -- …/report_second_opinion.go`: only the 4-line function comment of `secondOpinionDecorate` and the `block :=` line.
5. GREEN probe: PASS at tip (see above).

**Scratch pages.** Folder `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/build-REQ-739/site/`. Decorated by the new code through a throwaway Go test. P2, P4, P5 were rebuilt from their own head metas; all six metas equal the published bytes on each. The "before (0.324.2)" copies are the new pages with only the new `<style>` element removed; the one for P5 is byte-identical to the published v1c04 page, which shows the metas rebuild the page exactly. P4 also has its published copy (0.324.0 header, unfolded) as `b0-v1c03.html`.

**Measurement script:** `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/build-REQ-739/measure.js` (run through the Playwright `browser_run_code` tool, light theme). Raw output: `…/build-REQ-739/measure-output.txt` in the same folder, copied below. The Chromium window shows a 15 px scroll bar, so `clientWidth` is 375 at 390.

```
390 | P1 kit new | {"textLeft":16,"rightGap":16,"headerHeight":75.6,"margin":"0px","padding":"0px 16px","mapTop":null,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":16,"summaryLeft":16,"iterationLines":1,"summaryLines":1,"kitCard":{"detailsBorder":"0px","detailsRadius":"0px","summaryWeight":"400","summarySize":"17px","bodySize":"17px","openAllOpensKit":true,"headerFoldStaysClosed":true}}
390 | P1 kit before(0.324.2) | {"textLeft":0,"rightGap":0,"headerHeight":75.6,"margin":"0px","padding":"0px","mapTop":null,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":0,"summaryLeft":0,"iterationLines":1,"summaryLines":1,"kitCard":{"detailsBorder":"0px","detailsRadius":"0px","summaryWeight":"400","summarySize":"17px","bodySize":"17px","openAllOpensKit":true,"headerFoldStaysClosed":true}}
390 | P2 v1c01 new | {"textLeft":16,"rightGap":16,"headerHeight":75.6,"margin":"0px","padding":"0px 16px","mapTop":null,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":16,"summaryLeft":16,"iterationLines":1,"summaryLines":1}
390 | P2 v1c01 before(0.324.2) | {"textLeft":0,"rightGap":0,"headerHeight":75.6,"margin":"0px","padding":"0px","mapTop":null,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":0,"summaryLeft":0,"iterationLines":1,"summaryLines":1}
390 | P3 plain new | {"textLeft":37.8,"rightGap":24,"headerHeight":22.5,"margin":"0px","padding":"0px 16px","mapTop":null,"scrollClosed":"390/390","scrollOpen":"390/390","iterationLeft":24,"summaryLeft":24,"iterationLines":1,"summaryLines":1}
390 | P3 plain before(0.324.2) | {"textLeft":21.8,"rightGap":8,"headerHeight":22.5,"margin":"0px","padding":"0px","mapTop":null,"scrollClosed":"390/390","scrollOpen":"390/390","iterationLeft":8,"summaryLeft":8,"iterationLines":1,"summaryLines":1}
390 | P4 v1c03 new | {"textLeft":31,"rightGap":31,"headerHeight":87.4,"margin":"12px 16px","padding":"8px 14px","mapTop":830,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":31,"summaryLeft":31,"iterationLines":1,"summaryLines":1}
390 | P4 v1c03 published(0.324.0) | {"textLeft":31,"rightGap":31,"headerHeight":1326.7,"margin":"12px 16px","padding":"8px 14px","mapTop":2069.3,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":31,"summaryLeft":null,"iterationLines":1,"summaryLines":null}
390 | P4 v1c03 before(0.324.2) | {"textLeft":31,"rightGap":31,"headerHeight":87.4,"margin":"12px 16px","padding":"8px 14px","mapTop":830,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":31,"summaryLeft":31,"iterationLines":1,"summaryLines":1}
390 | P5 v1c04 new | {"textLeft":16,"rightGap":16,"headerHeight":69.6,"margin":"0px","padding":"0px 16px","mapTop":398.3,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":16,"summaryLeft":16,"iterationLines":1,"summaryLines":1}
390 | P5 v1c04 published(=0.324.2) | {"textLeft":0,"rightGap":0,"headerHeight":69.6,"margin":"0px","padding":"0px","mapTop":398.3,"scrollClosed":"375/375","scrollOpen":"375/375","iterationLeft":0,"summaryLeft":0,"iterationLines":1,"summaryLines":1}
1440 | P1 kit new | {"textLeft":16,"rightGap":16,"headerHeight":75.6,"margin":"0px","padding":"0px 16px","mapTop":null,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":16,"summaryLeft":16,"iterationLines":1,"summaryLines":1,"kitCard":{"detailsBorder":"0px","detailsRadius":"0px","summaryWeight":"400","summarySize":"17px","bodySize":"17px","openAllOpensKit":true,"headerFoldStaysClosed":true}}
1440 | P1 kit before(0.324.2) | {"textLeft":0,"rightGap":0,"headerHeight":75.6,"margin":"0px","padding":"0px","mapTop":null,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":0,"summaryLeft":0,"iterationLines":1,"summaryLines":1,"kitCard":{"detailsBorder":"0px","detailsRadius":"0px","summaryWeight":"400","summarySize":"17px","bodySize":"17px","openAllOpensKit":true,"headerFoldStaysClosed":true}}
1440 | P2 v1c01 new | {"textLeft":16,"rightGap":16,"headerHeight":75.6,"margin":"0px","padding":"0px 16px","mapTop":null,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":16,"summaryLeft":16,"iterationLines":1,"summaryLines":1}
1440 | P2 v1c01 before(0.324.2) | {"textLeft":0,"rightGap":0,"headerHeight":75.6,"margin":"0px","padding":"0px","mapTop":null,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":0,"summaryLeft":0,"iterationLines":1,"summaryLines":1}
1440 | P3 plain new | {"textLeft":37.8,"rightGap":24,"headerHeight":22.5,"margin":"0px","padding":"0px 16px","mapTop":null,"scrollClosed":"1440/1440","scrollOpen":"1440/1440","iterationLeft":24,"summaryLeft":24,"iterationLines":1,"summaryLines":1}
1440 | P3 plain before(0.324.2) | {"textLeft":21.8,"rightGap":8,"headerHeight":22.5,"margin":"0px","padding":"0px","mapTop":null,"scrollClosed":"1440/1440","scrollOpen":"1440/1440","iterationLeft":8,"summaryLeft":8,"iterationLines":1,"summaryLines":1}
1440 | P4 v1c03 new | {"textLeft":43.5,"rightGap":43.5,"headerHeight":87.4,"margin":"12px 28.5px","padding":"8px 14px","mapTop":401.9,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":43.5,"summaryLeft":43.5,"iterationLines":1,"summaryLines":1}
1440 | P4 v1c03 published(0.324.0) | {"textLeft":43.5,"rightGap":43.5,"headerHeight":538.5,"margin":"12px 28.5px","padding":"8px 14px","mapTop":853,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":43.5,"summaryLeft":null,"iterationLines":1,"summaryLines":null}
1440 | P4 v1c03 before(0.324.2) | {"textLeft":43.5,"rightGap":43.5,"headerHeight":87.4,"margin":"12px 28.5px","padding":"8px 14px","mapTop":401.9,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":43.5,"summaryLeft":43.5,"iterationLines":1,"summaryLines":1}
1440 | P5 v1c04 new | {"textLeft":16,"rightGap":16,"headerHeight":69.6,"margin":"0px","padding":"0px 16px","mapTop":384.9,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":16,"summaryLeft":16,"iterationLines":1,"summaryLines":1}
1440 | P5 v1c04 published(=0.324.2) | {"textLeft":0,"rightGap":0,"headerHeight":69.6,"margin":"0px","padding":"0px","mapTop":384.9,"scrollClosed":"1425/1425","scrollOpen":"1425/1425","iterationLeft":0,"summaryLeft":0,"iterationLines":1,"summaryLines":1}
```

**Matrix (new code; before in brackets).**

| Page | Width | Left text edge x | Right gap | Header height | margin / padding | Map top y | Sideways scroll closed / open |
| --- | --- | --- | --- | --- | --- | --- | --- |
| P1 kit | 390 | 16 [0] | 16 [0] | 75.6 [75.6] | 0 / 0 16px | none | none / none |
| P1 kit | 1440 | 16 [0] | 16 [0] | 75.6 | 0 / 0 16px | none | none / none |
| P2 v1c01 | 390 | 16 [0] | 16 [0] | 75.6 [75.6] | 0 / 0 16px | none | none / none |
| P2 v1c01 | 1440 | 16 [0] | 16 [0] | 75.6 | 0 / 0 16px | none | none / none |
| P3 plain | 390 | 24 box, 37.8 first word [8, 21.8] | 24 [8] | 22.5 [22.5] | 0 / 0 16px | none | none / none |
| P3 plain | 1440 | 24 box, 37.8 first word [8, 21.8] | 24 [8] | 22.5 | 0 / 0 16px | none | none / none |
| P4 v1c03 | 390 | 31 [31] | 31 [31] | 87.4 [87.4] | 12px 16px / 8px 14px [same] | 830 [830] | none / none |
| P4 v1c03 | 1440 | 43.5 [43.5] | 43.5 [43.5] | 87.4 [87.4] | 12px 28.5px / 8px 14px [same] | 401.9 [401.9] | none / none |
| P5 v1c04 | 390 | 16 [0] | 16 [0] | 69.6 [69.6] | 0 / 0 16px | 398.3 [398.3] | none / none |
| P5 v1c04 | 1440 | 16 [0] | 16 [0] | 69.6 [69.6] | 0 / 0 16px | 384.9 [384.9] | none / none |

- P3 has no tree marks, so its first text box is the word "Preset:" after the "▸" fold marker; the header box edge (and the marker) is at 24 px, which is D-02's figure (8 px body margin + 16 px).
- Fold alignment: the "▸ Preset" summary starts at the same x as the iteration line on every page and width (16/16, 31/31, 43.5/43.5).
- No header line wrapped: the iteration line and the summary are 1 line each on every page at 390, and every header height is unchanged.
- P1 kit checks: header `details` border 0 px, radius 0 px, summary weight 400 and 17 px (body 17 px), so it is not the kit card; "Open all details" opened all kit `details` in its section and left the header fold closed.
- Judge: `ai-report-judge` on the scratch bundle (outside the repository, accepted) for P4 `index-v1c03-…-review-loop-map.html` and P5 `index-v1c04-…-phone-first-map.html`: both `outcome: success`, "ai-report-judge verdict pass", no findings. Verdicts in `…/build-REQ-739/judge-v1c03/judge.json` and `judge-v1c04/judge.json`.
- Test-file budget: `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/` took 16.98 s before the change and 15.47 s after.

## Captures for the maintainer
Folder `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-202422/REQ-739-captures/`, all 390 x 844, light theme, fold closed:
- `P1-kit-page-390.png`: the report kit skeleton; "Iteration v1c01" and "▸ Preset: bold" now start 16 px from the edge, near the kit's own 20 px content edge.
- `P2-custom-v1c01-390.png`: the tstack v1c01 custom page (`body{margin:0}`, no header rule); header text at 16 px.
- `P3-plain-page-390.png`: a page with no CSS; the header box starts at 24 px (default 8 px body margin + 16 px), the body text at 8 px.
- `P4-v1c03-own-header-rule-390.png`: the page with its own header rule; the bordered header card looks exactly as before (12 px top, 16 px sides, padding 8 px 14 px).
- `P5-v1c04-req737-final-390.png`: REQ-737's final relation-map page; header text at 16 px, aligned with the page's own 16 px content edge; map top unchanged at y = 398.

Form chosen: D-01 as recommended, a `<style>` inside the markers with `:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}`. It met GREEN (2) and (3), so no other form was tried. Forms rejected (from D-01): an inline padding or margin on the `aside`; a `<style>` in the head; a class with a rule; inline padding on the `p` and `summary`.

## P-A-U
- PLAN: add the RED test with 5 checks, then one `<style>` with a `:where()` rule inside the markers before the `aside`; prove on scratch pages P1-P5 at 390 and 1440.
- APPLY: two commits, the test and the block string plus its comment; nothing else in the launcher moved.
- UNIFY: `git diff --stat main...HEAD` lists only `report_second_opinion.go` (+5 −3, comment and `block :=` line) and `report_second_opinion_test.go` (+33, one new function); no scratch test, no debug output, no `ai-reports/` change.

## Decisions
- **D-08** (DECIDE & STATE): The "before" copies for P1, P2, P4 and P5 are the new pages with only the new `<style>` removed, which is exactly what 0.324.2 writes; for P4 the published 0.324.0 copy is measured too. The published v1c03 page still carries the 0.324.0 unfolded header (1326.7 px tall at 390), so only its margin and padding are comparable; height and map top are compared against the 0.324.2 copy (830 at 390, as REQ-738 measured).
- **D-09** (DECIDE & STATE): The judge accepted the scratch bundle outside the repository, so no copy went into the worktree's `ai-reports/`.
- **D-10** (DECIDE & STATE): Captures were written straight into the captures folder by `page.screenshot`; the only file this session created in the main tree's `.playwright-mcp/` was `console-2026-11-16T20-34-09-044Z.log` (a favicon 404), deleted by name. The folder is now empty.
- **D-11** (DECIDE & STATE): P3 was decorated with an editor direction and no tree marks (a plain `custom` run), so it has no iteration line; its capture is kept because it looks different from P1 (24 px edge, no iteration line).

## Discovered Tasks
- The published tstack v1c03 page still shows the 0.324.0 unfolded header (1326.7 px tall at 390, map top y = 2069). Pages are not rewritten by design; a new pass would rebuild it. No action taken.

## Integration seams
- Page rule wins: P4 margin `12px 16px` / padding `8px 14px` at 390 and `12px 28.5px` / `8px 14px` at 1440, equal to both before copies. Map top y: P4 830 (390) and 401.9 (1440); P5 398.3 (390) and 384.9 (1440); all unchanged. Header heights unchanged (P4 87.4, P5 69.6); no header line wrapped on the phone.
- `git diff --stat main...HEAD -- ai-reports/` prints nothing.
- `^TestSecondOpinion` time: 16.98 s before, 15.47 s after.
- Exact block for a tree page with preset `relation-map` (direction shortened with "…"):

```
<!-- second-opinion-header --><style>:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}</style><aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere"><p><strong>Iteration v1c03</strong>, built on <a href="index.html">index.html</a></p><details style="background:none;border:0;border-radius:0;padding:0;margin:0"><summary style="cursor:pointer;font:inherit;padding:0;margin:.25em 0;border:0"><strong>Preset:</strong> relation-map (full direction)</summary><p style="white-space:pre-line">Make a relation map: one zoomable inline SVG map whose boxes…</p></details></aside><!-- /second-opinion-header -->
```

## Proposed CHANGELOG entry
Title: **Second-Opinion Header Keeps a Side Margin on Phones** (no earlier entry uses "side margin").

The short header that `ai-report second-opinion` writes above a directed page no longer touches the screen edge on a phone.

What changes for an installed copy: a page that `ai-report second-opinion` builds now gives its short header about 16 px of side room on pages that do not style the header themselves; a page with its own rule for the header keeps its own look; the head metas and the fold are unchanged; pages already published keep their old header until a new pass rebuilds them.

## Proposed lesson
- [family: writer-anchored-marker-position] REQ-739: a launcher that adds style to a block a page may also style must put the rule in a zero-specificity `:where()` selector inside its own markers, never inline on the element: an inline style beats the page's own rule, while a `<style>` between the markers is masked and stripped with the block, so idempotence and the old-block strip need no new code.

## Lessons read
- `_dev/primes/lessons-releases.md` (required), `_dev/primes/prime-releases.md`.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` lines 105-107 `[family: writer-anchored-marker-position]` and lines 17-20 `[family: session-is-the-backend]`, by fixed-string grep.
- `skills/do-work/crew-members/testing.md`, `shared-principles.md`; the REQ's Exploration and Decisions.
