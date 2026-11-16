# REQ-738 Hand-back — the second-opinion launcher folds the full direction under a short header

Branch `worktree-agent-REQ-738-fold-second-opinion-direction`, tip `6432bb22`, base `6f00ab48` (VERSION 0.324.0).

## Commits

- `975066d9` [REQ-738] RED: pin the folded second-opinion direction under a short header
- `152d6bbb` [REQ-738] fold the second-opinion direction under a short header
- `6432bb22` [REQ-738] say the second-opinion header folds the full direction

## Files

- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` — new `TestSecondOpinionHeaderFoldsTheFullDirection` (7 checks) and the `html` import.
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` — only the slot branch inside `secondOpinionDecorate` (one line plus one comment line) and the function comment.
- modified `skills/do-work-toolbox/actions/ai-report.md` — line 89, one clause.
- modified `skills/do-work-toolbox/docs/ai-report-guide.md` — line 124, one clause.

## Scope changes

- `report_second_opinion_tree_test.go`: declared, not needed (D-08). Its asserts (parent link, id and base metas, plain run has no marks) hold for the new form.
- The existing asserts in `TestSecondOpinionDirectedRunNamesAndMarksPage` also hold unchanged.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`: not touched, integrator only (D-05).

## Red-green evidence

- RED after commit 1 (`go test -count=1 -run '^TestSecondOpinionHeaderFoldsTheFullDirection$' ./internal/toolboxcommands/`): `--- FAIL ... report_second_opinion_test.go:1030: check 1: the header has no closed details fold, so the whole direction is always-open text: <aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere"><p><strong>Iteration v1c03</strong>, built on <a href="index.html">index.html</a></p><p><strong>Second-opinion preset:</strong> relation-map</p><p style="white-space:pre-line"><strong>Directio…` — fails on check 1, not on compile or a missing file.
- GREEN after commit 2: `--- PASS: TestSecondOpinionHeaderFoldsTheFullDirection (0.00s)`, with `TestSecondOpinionDirectedRunNamesAndMarksPage` (0.71 s) and `TestSecondOpinionTreeTargetsNamesAndParentHeader` (0.95 s) also PASS.
- GREEN probe at base `main` (6f00ab48, detached scratch worktree, removed after): `REQ-738 probe FAIL: report_second_opinion_test.go has no TestSecondOpinionHeaderFoldsTheFullDirection`, exit 1.
- GREEN probe at tip `6432bb22`: `REQ-738 probe PASS (3 header tests, 2 contract sentences)`.

## Verification

1. `gofmt -l` on the three Go files: prints nothing. `go vet ./internal/toolboxcommands/`: clean. `go test -count=1 ./internal/toolboxcommands/`: `ok … 31.078s` (real 31.80 s).
2. `shipped-package-reference-contract.sh`: PASS. `contract-regressions.sh`: passed. `quiet-grep-pipeline-audit.sh`: passed (240 files). `DO_WORK_MAINTAINER_TIER=heavy staged-skills-contract.sh`: PASS in 53 s.
3. `git diff --stat main...HEAD`: the 4 files above only. `git status --short`: clean (no scratch test file, no folder under `ai-reports/`).
4. The launcher diff touches only `secondOpinionDecorate`'s comment and its slot branch.
5. GREEN probe: PASS (above).

Measurement script: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f2615d23-6af3-4e32-b77d-010b0c2bd2d4/scratchpad/build-REQ-738/measure-req738.js` (sha1 `06d9862c5447f0b6f47419c35a83dbbb8556e345`). The browser tool refuses files outside the main tree, so the same text ran inline. The click sweep in it used synthetic click events, which the map ignores (it selects on pointerdown/pointerup), so the sweep reran once with real mouse clicks; both outputs are below. Chromium headless, light theme. `clientWidth` is 375 at 390 and 1425 at 1440 because of the scrollbar.

Raw output, main run (abridged only where marked):

```
390:  s1Load {headerPx 87.4, asideBox 63.4, mapFrameTop 830, sideways 0, foldOpen false, ownRule "1px solid rgb(195, 201, 208) / 12.8px / 10px"}
      s1Open {headerPx 1326.7, asideBox 1302.7, mapFrameTop 2069, sideways 0, foldOpen true, lineBreaks {whiteSpace pre-line, renderedLines 22, sourceLines 22}}
      tab [{A[index.html] inView true, focusRing "auto 1px", focusVisible true}, {SUMMARY inView true, focusRing "auto 1px", focusVisible true}, {enterOpens true}, {spaceCloses true}]
      noJs {headerPx 87.4, asideBox 63.4, mapFrameTop 748, sideways 0, foldOpen false, foldOpensByClick true, listUnderMap "118 items in #rev-3, visible true", jsClass false}
      s2 {summary "Preset: custom (full direction)", foldHoldsDirection true, ownDetails 6, opensByClick 6, ownDetailsStyleSameAsPublished true, headerPx 75.6, publishedHeaderPx 170, sideways 0}
      s3 {headerPx 75.6, detailsBorder 0px, detailsRadius 0px, detailsPadding 0px, detailsBackground rgba(0,0,0,0), summaryWeight 400, summarySize 17px, summaryText "Preset: bold (full direction)",
          afterOpenAll {kitDetailsOpen "true,true,true", headerFoldOpen false}, openedFold {summaryBorderBottom "0px none", summaryMarginBottom 4.25px}}
      s4 {sidewaysClosed 0, sidewaysOpen 0}
1440: s1Load {headerPx 87.4, asideBox 63.4, mapFrameTop 402, sideways 0, foldOpen false, ownRule same}
      s1Open {headerPx 538.5, asideBox 514.5, mapFrameTop 853, sideways 0, lineBreaks {pre-line, 22, 22}}
      tab, noJs (mapFrameTop 360), s2 (headerPx 75.6, publishedHeaderPx 115.6), s3, s4: same values as 390
scriptErrors []
(synthetic-click sweep: every box "selected" nothing; discarded, rerun below)
```

Click sweep rerun (real mouse clicks; title px = font size × screen scale; out of frame = selected or related boxes outside the frame):

```
390 first view: S1 S2 S3 L2 Q1 R710 not in the frame at start (the phone first view shows D, O1, O2, O3).
  D 7.2 | O1 O2 O3 L1 7.2 | 0 | sideways 0;  O1 7.2 | S3 D 7.2 | 0 | 0;  O2 4.3 | 11 related 4.3 | 0 | 0;  O3 4.3 | S1 D L3 4.3 | 0 | 0
390 after Show everything: all 23 clicked; titles 7.2 (S1-S6, D, O1, J, P1, P4, P5, L1-L3, Q1-Q3, R710) or 4.3 (O2, O3, P2, P3); out of frame 0; sideways 0 every click.
1440 first view: all 10 clicked; titles 15.2 (S1, D), 16.5 (S2, S3, O1, L2, Q1, R710), 10 (O2), 12.1 (O3); out of frame 0; sideways 0.
1440 after Show everything: 20 of 23 clicked; titles 10-16.5; out of frame 0; sideways 0. S1, D and L1 were not selected because the click point at the visible part's centre hit another element (a line or label); harness limit, the map is unchanged.
errors []
```

The eight measurements:

| # | Check | 390 x 844 | 1440 x 1000 |
| --- | --- | --- | --- |
| 1 | S1 at load, fold closed | header 87.4 px (box 63.4 + 12 + 12 margins), map frame top y = 830, sideways 0 | header 87.4 px, map frame top y = 402, sideways 0 |
| | 0.324.0 for comparison | 1,302 px, y = 2,069 | 514 px, y = 853 |
| 2 | S1 fold opened | header 1,326.7 px, sideways 0, pre-line, 22 lines as in the source | header 538.5 px, sideways 0, 22 lines |
| 3 | Click sweep | no new sideways scroll, no script error; numbers above | same |
| 4 | Tab order | 1st stop parent link, 2nd stop summary, both in view, focus ring `auto`; Enter opens, Space closes | same |
| 5 | JavaScript off | short header shows (87.4 px), click opens the fold, the box list `#rev-3` is there (118 items) | same, map top 360 |
| 6 | S2 custom | summary "Preset: custom (full direction)", fold text equals the direction meta, own 6 `details` styled as on the published page and all open by click; header 75.6 px (published 170) | header 75.6 px (published 115.6) |
| 7 | S3 kit page | no kit card (border 0, radius 0, padding 0, weight 400, size 17 px = inherited); header 75.6 px; "Open all details" opens the kit's 3, header fold stays closed; opened fold has no dashed border | same |
| 8 | S4 400-char token | sideways 0 closed and open | same |

D-03 reading: the header is 87.4 px with margins (63.4 px box), so above "about 80" but under 100, and the map frame top is y = 830, below 844. Met under D-03. Only 14 px of the map frame shows on the phone's first screen.

S1 metas: byte-identical to the published v1c03 page (3,473 bytes); S2 metas byte-identical too (469 bytes). The scratch decorator was a throwaway `zz_req738_scratch_test.go`, deleted before the last check.

Render judge on the scratch S1 bundle (outside the repository, accepted with no copy into `ai-reports/`): `outcome success`, change detail `ai-report-judge verdict pass`.

Timing for the 30-second per-file budget, `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/`: before the GREEN change 19.229 s (with the RED test failing in 0.00 s), after 23.629 s. The new test runs in 0.00 s, so the difference is machine load from the second builder.

## Captures for the maintainer

All in `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-16-171725/REQ-738-captures/`:

- `S1-390x844-load.png` — the tstack relation-map page on a phone at load: two-line header ("Iteration v1c03, built on index.html" / "▸ Preset: relation-map (full direction)"), then the page; the map frame starts at the bottom edge (y = 830).
- `S1-390x844-fold-open.png` — same page after a tap on the summary: the full relation-map direction with its line breaks.
- `S1-1440x1000-load.png` — desktop at load: the two-line header, and the map with its first view fully on the first screen.
- `S3-390x844-kit-load.png` — a report-kit page (bold preset): the header has no kit card look.
- `S2-390x844-load.png` — the custom page: "▸ Preset: custom (full direction)" above the page's own styled content.

Form chosen (D-09): the iteration line stays a `p`, and the preset name is the `summary` of a closed `details` that holds the direction: `▸ Preset: relation-map (full direction)`.
Forms rejected, measured on S1 at 390 x 844 with the page's own header rule:
- D-01 as written (iteration `p`, preset `p`, "Full direction" summary): three lines, map top y = 848.
- Iteration and preset in one `p` ("…, built on index.html, preset relation-map") plus "Full direction": that `p` wraps to two lines, y = 848.
- A shorter one-`p` wording ("Iteration v1c03 of index.html, preset relation-map"): y = 830, but it fits by only a few pixels, so a longer preset name or tree id wraps it.
- Summary "Second-opinion preset: relation-map, full direction": wraps to two lines, y = 848.

## P-A-U

- PLAN: keep the markers, metas, strip calls and insert position; replace only the slot paragraphs with a closed, kit-neutral `details`; pin with a direct-call test; measure on scratch copies.
- APPLY: three commits: RED test, the slot branch of `secondOpinionDecorate` and its comment, two contract clauses.
- UNIFY: `git diff --stat main...HEAD` shows 4 files (launcher +3/−1 lines, test +75, `ai-report.md` 1 line, `ai-report-guide.md` 1 line); no debug output, no scratch file, `ai-reports/` unchanged; gofmt, vet and the package tests are clean.

## Decisions

- **D-09** (DECIDE & STATE): The preset name moved into the `summary` ("Preset: <slot> (full direction)") instead of its own `p`. D-01's form measured y = 848 on the phone, above the 844 limit. The iteration line keeps its `p`, the `aside`, its label and the long-token guard stay. The test reads the summary as visible and the `details` body after `</summary>` as the folded part.
- **D-10** (DECIDE & STATE): The visible label is now "Preset:" instead of "Second-opinion preset:", and the folded text no longer starts with "Direction:". The `aside` label "Second-opinion direction" and the "(full direction)" words carry that context, and the shorter label keeps the summary on one phone line for every shipped preset name (checked with `visual-studyguide`).
- **D-11** (DECIDE & STATE): The click sweep reran with real mouse clicks because the map ignores synthetic click events. Boxes outside the frame at the start (phone first view) and three desktop boxes whose centre click hit another element are reported, not forced. This REQ does not change the map (D-07).
- **D-12** (DECIDE & STATE): The judge accepted the scratch bundle outside the repository, so nothing was copied into `ai-reports/`.
- **D-13** (DECIDE & STATE): I deleted one file my session created in the main tree's `.playwright-mcp/`: `console-2026-11-16T17-33-03-035Z.log` (one favicon 404 entry). The captures were written straight to the captures folder.
- **D-14** (DECIDE & STATE): Commit trailers use `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`, following the brief as updated on disk and the dispatch prompt.

## Discovered Tasks

- On a page without its own header rule (the report kit, S3; the custom page, S2), the header text has no side gutter: it starts at x = 0 on a phone. 0.324.0 has the same gap. A fix would add a padding to the `aside`, but an inline padding would override the tstack page's own rule, so it needs its own decision.

## Integration seams

- **Exact new block, tree page, preset `relation-map`** (direction shortened as `…`):
  `<!-- second-opinion-header --><aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere"><p><strong>Iteration v1c03</strong>, built on <a href="index.html">index.html</a></p><details style="background:none;border:0;border-radius:0;padding:0;margin:0"><summary style="cursor:pointer;font:inherit;padding:0;margin:.25em 0;border:0"><strong>Preset:</strong> relation-map (full direction)</summary><p style="white-space:pre-line">Make a relation map: one zoomable inline SVG map …</p></details></aside><!-- /second-opinion-header -->`
  With no `</head>` on the page, the metas sit inside the block before the `aside`, as in 0.324.0.
- **Header height:** 87.4 px with margins (63.4 px box) at both 390 and 1440. Map frame top on S1: y = 830 at 390 x 844, y = 402 at 1440 x 1000. Measure the header by `body>aside[aria-label="Second-opinion direction"]` (box height plus computed top and bottom margins) and the map by `#map-frame`.
- **The page's own rule:** the element, label and `p` children stay. The tstack v1c03 rule applies as before (1 px border, 10 px radius, 12.8 px font, `p` margin .25em). The summary takes the same font through `font:inherit` and the same .25em margin. See `S1-390x844-load.png`.
- **Published pages:** `git diff --stat main...HEAD -- ai-reports/` prints nothing.
- **Test-file budget:** `-run '^TestSecondOpinion'` took 19.229 s before the change and 23.629 s after it, on a machine shared with REQ-737's builder. The new test runs in 0.00 s.

## Proposed CHANGELOG entry

Title: `Second-Opinion Pages Fold the Full Direction Under a Short Header` (no earlier entry uses it).

A directed second-opinion page now opens with a two-line header instead of the whole direction text. On a phone, the relation-map page's map now starts on the first screen (y = 830, was 2,069).

- The header shows the iteration line with its parent link and "▸ Preset: <name> (full direction)". A click or Enter on that line opens the full direction with its line breaks. It works without JavaScript, and the report kit neither restyles it nor opens it with "Open all details".
- What changes for an installed copy: a page that `ai-report second-opinion` builds with a preset or an editor direction now opens with a short header, the full direction is behind a fold on the same page, the head metas are unchanged, and pages already published keep their old header until a new pass rebuilds them.

## Proposed lesson

- [family: session-is-the-backend] REQ-738: a fold meant to keep content on the first phone screen has to be measured with the page's own CSS, at the real width. The recommended three-line header missed 844 px by 4 px, and only putting the preset name in the summary line made it fit. Count lines, not elements.

## Lessons read

- `_dev/primes/lessons-releases.md` (required; no archive move or manifest edit here).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` lines 103-105 (`writer-anchored-marker-position`: no whitespace outside the markers, strip on the masked surface, idempotent decorate) and line 18 (`session-is-the-backend`: measure after decorating). Both followed: markers and strips unchanged, check 4 and check 5 pin idempotence and the old-block strip, and the measurement used decorated scratch pages.
