# REQ-727 hand-back: ship a relation-map preset for ai-report second-opinion

ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html

Branch `worktree-agent-REQ-727-relation-map-preset`, tip `f2edff06a5bcacbe90ee61b17d6765661fbd9d92`, built on `445950c4`. Nothing pushed or merged.

## Commits

- `773d18a7` [REQ-727] RED: the shipped-presets test names relation-map
- `688d9e6f` [REQ-727] GREEN: ship the relation-map second-opinion preset
- `f2edff06` [REQ-727] build tree page v1c03 with the relation-map preset

## Files

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`: modified. The function comment, the list comment, the named list and the error text changed (4 lines).
- `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md`: created (3,057 bytes, 22 lines).
- `skills/do-work-toolbox/docs/ai-report-guide.md`: modified (one clause at line 124).
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`: created by the launcher at finalize.
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/v1c03-claude-claude-opus-5-5-relation-map-review-loop-map-reply.md`: created by the launcher at finalize.
- `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html`: modified by the launcher (one `v1c03` line in the iteration list, line 6).

The resolved names for the two Scope globs (D-11) are the two `v1c03` names above.

## Scope changes

None. No file outside the Scope list changed. `VERSION`, both `CHANGELOG.md` files, the catalog pair and every `do-work/` path are untouched on the branch.

## Red-green evidence

- RED after `773d18a7`: `go test -count=1 -run '^TestSecondOpinionShippedPresetsLoadFromTheToolbox$' ./internal/toolboxcommands/` printed `--- FAIL: TestSecondOpinionShippedPresetsLoadFromTheToolbox (0.17s)` and `report_second_opinion_test.go:1028: named preset "relation-map" does not ship in …/skills/do-work-toolbox/docs/second-opinion-presets`.
- GREEN after `688d9e6f`: the same command printed `--- PASS: TestSecondOpinionShippedPresetsLoadFromTheToolbox (0.18s)`. The loop over shipped files prepared `relation-map` with state `prepared` and `preset_file` in the shipped folder.
- GREEN probe at base: its first check needs the preset file, and `git cat-file -e 445950c4:skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` fails (`exists on disk, but not in '445950c4'`). I did not create a base checkout to run the whole probe (D-21).
- GREEN probe after `688d9e6f`: `REQ-727 probe FAIL: no second-opinion page built with --preset relation-map outside the example bundle` (check 5, as expected).
- GREEN probe at tip `f2edff06`: `REQ-727 probe PASS (ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html)`, exit 0.

## Verification

Commands, from the worktree:
- `gofmt -l internal/toolboxcommands/report_second_opinion_test.go`: prints nothing.
- `go vet ./internal/toolboxcommands/`: clean.
- `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/`: `ok … 18.957s` (19.3 s wall). One run, no rerun needed.
- `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0, `shipped package reference contract: PASS`.
- `bash _dev/tests/contract-regressions.sh`: exit 0, `Contract regression checks passed.`
- `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`: exit 0, `staged skills contract: PASS` (41.0 s wall).
- `git diff --stat main...HEAD -- …/report_second_opinion.go …/ai-report-second-opinion-prompt.md …/ai-report-kit/`: prints nothing.
- `git diff --stat main...HEAD -- skills/do-work-toolbox/docs/second-opinion-presets/`: only `relation-map.md | 22 +`.
- Finalize judge: `judge_verdict = pass`, `judge_findings = []`, `banner_state = inserted`, `tree_id = v1c03`, `model_evidence = session-declared`.
- Standalone re-judge: `ai-report-judge … --page index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html --out <scratch>/judge`: exit 0, `findings: []`, `ai-report-judge verdict pass`. Output folder deleted.
- Preset size: 3,057 bytes (target about 3,500 or less, D-02). It contains no `{{`.

Browser measurements. The method: a Playwright script served the bundle with `python3 -I -m http.server 8737 --bind 127.0.0.1`, light theme, at 1440x1000 and 390x844. Dark theme was also checked at load. The script clicked every shown box from the first view and from "Show everything", swept zoom from 0.15 to 2.50 with O2 selected (22 views), and tabbed through every box at high zoom. Scale and title size are given as world scale and the 22 px title's screen size.

| Check (D-04) | Staging 1440 | Staging 390 | Published 1440 | Published 390 |
| --- | --- | --- | --- | --- |
| (1) Relation words at load, no click | 12 of 12 (scale 0.564) | 2 of 12 (D, O1-O3 framed; the other words' lines leave the frame) | 12 of 12 | 2 of 12 |
| (2) Words on boxes: load, every click, Show everything, zoom sweep | 0 in 1 + 34 + 22 views; 0 word overlaps; 0 words cut at the frame edge; 0 cut column titles | 0 in 1 + 28 + 22 views; 0 overlaps; 0 cut | 0 in 1 + 34 + 22 views; 0 overlaps; 0 cut | 0 in 1 + 28 + 22 views; 0 overlaps; 0 cut |
| (3) Words broken inside a label | 0 (labels break only at spaces; checked in every view above) | 0 | 0 | 0 |
| (4) Clicked or related shown box outside the frame | 0 of 33 clicks | 0 of 27 clicks | 0 of 33 | 0 of 27 |
| (5) Tab at high zoom, focused box fully in frame, frame scroll 0 | 23 of 23 at 2.50 | 23 of 23 at 1.03 | 23 of 23 at 2.50 | 23 of 23 at 1.03 |
| (6) Phone: first-view titles, page overflow, table rows | n/a; page overflow 0 | title 12.9 px with no pan (D, O1, O2, O3 fully in frame); page overflow 0; 0 of 28 table rows wider than the screen | page overflow 0 | 12.9 px; overflow 0; 0 of 28 rows wider |
| (7) No JavaScript | tool bar `display:none`; 28 of 28 words on lines, 0 on boxes, 0 overlapping; map 1800x2480 scrolls in its frame; 23 boxes as text; 28 table rows; page overflow 0 | same | same | same |

Other measurements, the same on both pages:
- Desktop first-view title: 12.4 px. A desktop click lands at 0.45 to 0.75 (O2 with 11 relations: 7 of 11 words in view at 0.453).
- "Show everything": 23 of 23 boxes in the frame. Words in view: 5 of 28 at 1440 (scale 0.261) and 2 of 28 at 390 (scale 0.195). The count line names the rest.
- Light dimming: unrelated O1 at opacity 0.7 with O2 selected, and its title is still shown.
- Escape and a second click on the selected box both return the panel to "Click a box".
- A plain wheel over the map scrolls the page, and Control plus the wheel zooms the map.
- No box body overflows: summaries were checked at every zoom of the sweep.
- Direction block height on the published page: 514 px at 1440 and 1,303 px at 390. The map frame starts at y = 853 at 1440 (its first 147 px are on the first screen) and at y = 2,069 at 390 (the first phone screen is the direction text). No text leaves the first screen sideways at either width.
- Phone after one click: D and O1 land at 0.326, O2 and O3 at 0.195, with 0 to 1 relation words in view (D-16).

Browser captures were written to the main tree's `.playwright-mcp/` by the browser tool and deleted by name. That folder did not exist before this session, so it was removed once empty. The server was stopped.

## P-A-U

- PLAN: Commit 1 adds a RED test list entry. Commit 2 adds a 20-check preset under 3,500 bytes and one guide clause. Commit 3 builds a Session second opinion on the D-01 bundle, measured in a browser before and after finalize.
- APPLY: three commits as planned. The page came from a generator in my scratch folder, written by following the preset text the prompt showed. The preset needed no change during editing (D-18).
- UNIFY: `git diff --stat main...HEAD` shows six files: the page (617+), `index.html` (1+), the reply (19+), `ai-report-guide.md` (1 line), `relation-map.md` (22+) and the test (4 lines). I read every hunk. No debug output, no staging file and no `.playwright-mcp` file is in the diff. The guarded paths (loader, prompt, kit, other presets) show no diff.

## Decisions

- **D-13** (DECIDE & STATE): The page styles the launcher's direction block with smaller text (0.8 rem), a border and the page's side gutter, through a `body>aside[aria-label="Second-opinion direction"]` rule in the page's own CSS. Its words and order are unchanged. This answers the REQ-734 lesson that the inserted block had no side gutter. The heights in Verification are with this style.
- **D-14** (DECIDE & STATE): REQ-710 ("shipped in 0.310.0") is a box in the "what tests or limits it" column, not the "known" column. Its relations go to O2, J and Q1. From the known column, its line to Q1 would cross the proposed column's boxes, and it answers the report rather than supporting it. The page says these page codes are its own (D, J, P1 to P5, L1 to L3, REQ-710).
- **D-15** (DECIDE & STATE): The map draws only lines between neighboring columns or inside one column. So the link from S3 to L2 (both about `work.md`'s reviewer fallback) is not a line; L2's sources name `work.md:146` instead.
- **D-16** (ESCALATE): On a phone, check 15 (a click keeps the clicked box and every related shown box in the frame) makes a click zoom out to 0.195 to 0.326. At that scale 0 to 1 relation words fit and titles are 4 to 7 px; the panel below the map lists every relation. This is the phone limit D-05 names, now measured. Value: the rule keeps every related box in view, which REQ-733's G2 showed matters. Risk: the phone map is a navigation overview after a click, not a reading view. A fix is a phone-specific layout, such as stacking the related boxes. That is a preset change, reversible in one text file, and it needs the maintainer's taste.
- **D-17** (DECIDE & STATE): The preset stays at 3,057 bytes with 20 checks. On a phone the launcher prints it as 1,303 px of text above the page, so the map starts on the third screen. Removing checks would drop the D-04 defects the preset exists to prevent. The launcher's block is out of scope (Discovered Tasks).
- **D-18** (DECIDE & STATE): The preset text was clear enough to build from, so the finalized page is the one the preset produced. No re-prepare was needed. Where the editor still made choices (D-14, D-15, the R710 to Q1 words "answers: no hold" that fit beside their arc), they are the editor's reading, which the preset allows.
- **D-19** (DECIDE & STATE): The model slug is `claude-opus-5-5`, this session's declared model, recorded as `session-declared`. The launcher's date meta reads 2026-11-15.
- **D-20** (DECIDE & STATE): Repository files are cited as text paths, not links, because the prompt allows canonical links only when remote reachability is verified. The page's corrections cite current line numbers next to the report's pinned `ccea405e` lines.
- **D-21** (DECIDE & STATE): The GREEN probe was not run in a separate base checkout. Its first check needs the preset file, and `git cat-file` shows that file absent at `445950c4`, which is the base failure the brief predicts.

## Discovered Tasks

- The report kit's closing band (`skeleton.html`) does not suggest the `relation-map` preset (D-07) — impact-negligible → report only
- The launcher prints the full direction above every directed page (`report_second_opinion.go:1634`). With the 3,057-byte relation-map preset, that block is 1,303 px tall at 390 px wide, so a phone reader meets the map on the third screen. Folding the direction after its first paragraph (for example in a `<details>`) would be a launcher change — impact-user-visible → report only
- The render judge does not measure label overlap, broken words, zoom framing, focus or no-JS line words (D-05). The browser script used here (relation words in view, words on boxes, words cut at the frame, related boxes out of frame after every click, Tab framing) could become a judge check for map pages — impact-negligible → report only

## Integration seams

| Shared file | What I changed | Lines after my edit |
| --- | --- | --- |
| `skills/do-work-toolbox/docs/ai-report-guide.md` | In the "A Second Opinion" paragraph, "`schematics` change how it teaches." became "`schematics` change how it teaches, and `relation-map` draws the report's items and the relations it states between them as one zoomable map." Nothing else changed; the REQ-728 to REQ-731 lines keep their bytes. | 124 |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` | Function comment ("a named preset the user curated (UR-164, UR-176)"), list comment, `"relation-map"` added to the named list, error text "named preset %q does not ship". | 1004, 1025-1026, 1028 |
| `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index.html` | The launcher added one `v1c03` line to the iteration list. | 6 |

## Proposed CHANGELOG entry

Title checked: no earlier entry uses "relation-map" or "relation map".

```
## 0.324.0 — Relation-Map Preset for ai-report Second Opinion (2026-11-15)

A report can now be read as one zoomable map of its items and the relations it states between them.

- `ai-report second-opinion --preset relation-map` asks the editor for a left-to-right map of the report's items (evidence, then the decision and options, then limits and open questions), a first view of about nine boxes with the words on their lines, and a click that opens a box's sources and relations.
- The preset turns the defects found while building the first map pages into checks: words never cover a box or break inside a word, a click keeps related boxes in view, focus keeps the box in view, the phone view needs no sideways pan, and the page works without JavaScript.
- The guide lists the new preset.

What changes for an installed copy: `ai-report second-opinion --preset relation-map` writes a zoomable map of the report's items and their relations, and the guide lists it. Nothing else changes.
```

## Proposed lesson

- [family: session-is-the-backend] (`skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`) **Measure a generated map page with one script that clicks every shown box from each starting view, not a sample of clicks.** REQ-733's builder missed defects that only one click in one view showed. Here one Playwright script clicked all boxes from the first view and from "Show everything" at two widths (122 views with the zoom sweep), on the staging page and again on the published page. It found the one missing first-view word and the cut column titles before finalize. Re-measuring after finalize is needed because the launcher's direction block moves the map (514 px at 1440, 1,303 px at 390).

## Lessons read

- `_dev/primes/prime-releases.md` (full) and `_dev/primes/lessons-releases.md` (required, full; the `canonical-link-outlives-its-target` family applies to the integrator's UR-176 close, D-12).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:17`, the `[family: session-is-the-backend]` bullet, read with a fixed-string grep. Its point that `staging_reply` must be written in own-turns mode was followed.
- `do-work/archive/UR-177/REQ-733-rebuild-relation-map-page.md` (Decisions D-08 to D-20, Discovered Tasks, Lessons Learned), the three REQ-733 reviews (findings F1 to F7, G1 to G4, H1, H2), and `do-work/archive/UR-178/REQ-734-proposal-table-and-map-comparison.md` (Review, "Map vs table", Lessons Learned).
- Crew members: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`, `prompt-injection.md` and `anti-slop.md` (both loaded before reading the bundle). No prompt injection was seen in the bundle.
