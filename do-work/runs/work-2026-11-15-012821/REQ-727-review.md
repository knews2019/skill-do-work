# Review: REQ-727 (ship a relation-map preset for ai-report second-opinion)

**Approve.** The preset ships, loads through the shipped-presets test, and produced a real page whose seven D-04 checks I re-measured in a browser with the same results as the builder.
Route B | range `9d909325..9d07ce2f` (builder merge `eec3bad0` plus catalog rebuild `9d07ce2f`)

### What's built

- `ai-report second-opinion --preset relation-map` now finds a shipped 3,057-byte preset: one opening paragraph (how a relation map differs from a mind map, lines are the editor's reading, facts come from the report or repository) and 20 one-sentence checks.
- The guide names the preset, the shipped-presets test names it, and tree page v1c03 of the tstack review-workflows bundle was built with `--preset relation-map` alone. The catalog lists that page.
- Still missing by design: the report kit's closing band never suggests the preset (constraint D-07), and the maintainer has not yet looked at the page (D-08).

### Decisions / risks for you

- **D-16, phone map after a click.** A click on a phone zooms the map out to scale 0.205 to 0.342, so box titles are 4.5 to 7.5 px. The panel under the map still lists every relation. No acceptance criterion becomes false: check 17 and the Scope criterion cover only the phone's *first* view, and I measured that view at 12.9 px titles with D, O1, O2 and O3 fully in the frame. Value of keeping it: check 15 keeps every related box in view. Risk: after a click, the phone map is an overview, not something you can read. A fix is a change to the preset text and needs your taste (F1).
- **D-17, direction block above the page.** The launcher prints the whole preset above the page: 514 px at 1440 and 1,302 px at 390. The map frame starts at y = 853 and y = 2,069. No acceptance criterion requires the map on the first screen, and D-02 asked only that the builder report these heights, which it did. The 514 and 1,302 px heights already include this page's own CSS (D-13), which shrinks the block. A page generated without that CSS would have a taller block. The fix belongs in the launcher (F2).

### Findings

**Important:** None.

**Minor:**
- F1: On a phone, a click zooms the map to scale 0.342 (D, O1) or 0.205 (O2, O3), so box titles are 7.5 or 4.5 px and 0 to 1 relation words show. This happens because preset check 15 (`skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md:17`) keeps every related box in the frame. The first-view criteria still hold — impact-user-visible → report only
- F2: The launcher prints the full direction above every directed page (`report_second_opinion.go:1634`). With this preset, the block is 514 px at 1440 and 1,302 px at 390, so a phone reader meets the map on the third screen. This is the builder's discovered task and needs a launcher change — impact-user-visible → report only
- F3 (restatement sweep): The kit's closing band offers shipped presets only as *teaching methods* (`skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:177-190`, `skills/do-work-toolbox/actions/ai-report-reference.md:91`). Its layout line builds only with `--editor-direction`. So the new preset, which is a form and not a teaching method, is never suggested where a report proposes another layout. The text is not false, since it calls the presets "a starting menu", and the REQ's constraints keep both files unchanged. This is the builder's D-07 discovered task — impact-negligible → report only

**Nit:**
- F4: The guide sentence now chains three "and" clauses: "`bold` and `visual` change how the page looks, and `eli5`, … and `schematics` change how it teaches, and `relation-map` draws …" (`skills/do-work-toolbox/docs/ai-report-guide.md:124`). It is correct, but a full stop before "`relation-map` draws" would read more easily — impact-negligible → report only
- F5: Check 19 of the preset contains Markdown backticks (`` `[hidden]` ``). The launcher prints the direction as plain text, so the page and its meta tag show the backticks literally (page line 145). No other shipped preset uses backticks — impact-negligible → report only
- F6: The test's loop variable is still named `method` (`report_second_opinion_test.go:1026`) after the list became "named presets". Leaving it keeps the diff minimal — impact-negligible → report only

### Requirements Checklist

- [x] Shipped-presets test names `relation-map`, RED first then GREEN — delivered (`report_second_opinion_test.go:1026`; RED and GREEN output are in the hand-back; I re-ran the test: `ok`)
- [x] Guide's "A Second Opinion" list names `relation-map` — delivered (`ai-report-guide.md:124`)
- [x] Preset says how a relation map differs from a mind map — delivered (opening paragraph)
- [x] Preset asks for every Detailed Requirement as changed by Addendum 1 — delivered. Coverage: own items and stated relations as the editor's reading (opening), four line kinds with a key (4), semantic zoom (7), click detail with sources and relations (8), light dimming (8), drag, pinch, Control-scroll, buttons, keys and plain scroll (11), Show everything (10), no-JS map plus text and table (18, 19), no library (1), about 30 boxes else group (3), narrow screen starts on the main item (17), left to right (2), about nine first-view boxes (5), readable labels (12), facts from the report or repository (opening)
- [x] The seven REQ-733/REQ-734 defects become checks — delivered (6, 12-13, 14, 15, 16, 17-18, 19)
- [x] A decision records what a preset cannot guarantee — delivered (D-05, measured further in D-16 and D-17)
- [x] One real page on a bundle other than the example, with judge `pass`. The hand-back opens with its path, and a browser table covers both widths — delivered. The "shown to the maintainer" part is deferred by D-08 (escalated), not skipped
- [x] Loader, prompt, other presets, kit and existing bundle files unchanged except the iteration list — delivered (`git diff --stat 9d909325..9d07ce2f` shows 8 files, all declared)
- [x] Preset tests pass, gofmt clean — delivered (`gofmt -l` prints nothing; `go vet` clean)

### Acceptance Testing

**Result: Pass** (implementation and integration stages: test run on the merged tree plus a browser re-measurement of the published page)

- `go test -count=1 -run '^TestSecondOpinionShippedPresetsLoadFromTheToolbox$' ./internal/toolboxcommands/`: `ok`. `gofmt -l` on the test file printed nothing, and `go vet ./internal/toolboxcommands/` was clean.
- The preset agrees with the loader (`report_second_opinion.go:195-229`): a regular file, 3,057 bytes, no `{{`, no leading or trailing whitespace beyond one final newline, which `TrimSpace` removes.
- Catalog row (`ai-reports/catalog.json`) matches the page's meta tags: page path, `tree_id` v1c03, preset `relation-map`, model `claude claude-opus-5-5`, date 2026-11-15. `ai-reports/index.html` and the bundle's iteration list agree.
- Fact sample against the bundle report: 7 boxes and 3 relations all match the report.
  - Boxes: S1 (three providers, no explicit round cap), S3 (one remediation, then completed-with-issues), S6 (seven other skills; pr-wait inspected, not executed), O1 ("No implementation cost"), O3 ("High recurring cost", routine work pays the same overhead), L1 (source analysis, not a benchmark), Q2 (environments; use available harness capabilities).
  - Relations: D recommends O2, S4 → O2 "reuse its verdicts", S6 → P4 "from takeover".
  - The R710 correction is true: `do-work/archive/UR-163/REQ-710-*.md` is `completed`, `CHANGELOG.md:209` is 0.310.0 on 2026-11-05, and `ai-evaluation-loop-reference.md:60,69` state five total rounds or two stalled rounds, ending as `completed-with-issues`.
- Browser re-measurement. I used a Node Playwright script, light theme, with the bundle served on 127.0.0.1. It clicked every shown box from the first view and from "Show everything" (33 clicks at 1440, 27 at 390). Measured with `getBoundingClientRect`. Numbers in parentheses are the builder's.

| Check | 1440x1000 | 390x844 |
| --- | --- | --- |
| (a) Relation words at load, no click | 12 of 12, scale 0.564, title 12.4 px (same) | 2 of 12, scale 0.588, title 12.9 px, D/O1/O2/O3 fully in frame (2 of 12, 12.9 px) |
| (b) Words on boxes / words on words: load, every click, Show everything, O2 | 0 / 0 in 35 views (0 / 0) | 0 / 0 in 29 views (0 / 0) |
| (c) Words or column titles cut at the frame edge; words broken inside a label or box title | 0 cut, 0 broken (0, 0) | 0 cut, 0 broken (0, 0) |
| (d) Clicked or related shown box outside the frame after a click | 0 of 33 clicks (0 of 33) | 0 of 27 clicks (0 of 27) |
| (e) Tab at high zoom: focused box fully in frame, frame scroll 0 | 23 of 23 at 2.50 (23 of 23 at 2.50) | 23 of 23, settles at 1.08 (23 of 23 at 1.03) |
| (f) No JavaScript | tool bar `display:none`; 28 of 28 words as SVG `<text>`, max 0.9 px from their own line; 0 on boxes, 0 overlapping; frame scrolls 1800x2480; 23 box texts, 28 rows; page overflow 0 (same) | same; page overflow 0; 0 of 28 table rows wider than the screen (same) |
| Direction block height | 514 px, frame top y = 853 (514, 853) | 1,302 px, frame top y = 2,069 (1,303, 2,069) |
| Phone after one click | — | D, O1 0.342 (7.5 px); O2, O3 0.205 (4.5 px); 0 to 1 words (0.326 and 0.195; 4 to 7 px) |

- Other measurements match the builder: O2 click at 1440 gives scale 0.453 with 7 of 11 words, and "Show everything" gives 5 of 28 at 1440. With O2 selected, O1 dims to opacity 0.7. Escape and a second click both return the panel to "Click a box". The page logged no errors.
- The phone scales differ by about 5 % because my headless Chromium has overlay scrollbars, so the map frame is 356 px wide. The builder's numbers fit a frame about 15 px narrower, a classic page scrollbar. The conclusion of D-16 does not change.
- Clean-up: the server is stopped, the browser is closed, no screenshot was taken, and `.playwright-mcp/` does not exist.

### Suggested Additional Testing

- Deployment: unassessed. The release (VERSION, both CHANGELOG files) is the integrator's finalization step. After it, run `ai-report second-opinion --prepare --preset relation-map` from an installed copy and confirm `preset_file` points at the installed toolbox folder.
- Live acceptance: unassessed. The maintainer opens the v1c03 page on a real phone and judges D-16: whether a click should keep every related box in view or keep the titles readable.
- Edge case: run the preset on a report with more than 30 items to see whether an editor groups items and says so (check 3). No page exercises that yet.

### Scores (on the record — not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All criteria met. "Shown to the maintainer" is deferred by D-08 |
| Code Quality | 95% | Preset clear and loader-safe. Guide sentence is long (F4) |
| Test Adequacy | 90% | One name in an existing test, RED then GREEN. Page behaviour is proved by browser measurement, not a test |
| Scope | 100% | 8 files, all declared. Guarded paths unchanged |
| Risk | None | Docs, one test list entry and a report page |
| Acceptance | Pass | Test run plus browser re-measurement at two widths and with no JavaScript |

### Follow-ups created

None (6 findings report only)

## Review

**Overall: 95%** | 2026-11-15T11:45:07Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: phone click zooms the map to scale 0.205 to 0.342, so titles are 4.5 to 7.5 px (D-16, preset check 15 at `relation-map.md:17`). The first-view criteria still hold — impact-user-visible → report only
- F2: launcher direction block is 514 px at 1440 and 1,302 px at 390 above the map (D-17, `report_second_opinion.go:1634`) — impact-user-visible → report only
- F3: kit closing band offers presets only as teaching methods, and its layout line uses only `--editor-direction`, so relation-map is never suggested (`skeleton.html:177-190`, `ai-report-reference.md:91`). Not false, and the files are frozen by the constraints — impact-negligible → report only
- F4 (nit): guide sentence chains three "and" clauses (`ai-report-guide.md:124`) — impact-negligible → report only
- F5 (nit): preset check 19's backticks print literally on the page and in the meta tag (`relation-map.md:21`) — impact-negligible → report only
- F6 (nit): test loop variable still named `method` (`report_second_opinion_test.go:1026`) — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages. The shipped-presets test passed, gofmt and vet were clean, and the browser re-measurement at 1440x1000 and 390x844 plus a no-JS pass matched the builder's numbers. Phone scales were within 5 % because of scrollbar width.
**Restatement sweep:** redefined: the shipped preset set (new `relation-map`) and the shipped-presets test's named list with its "learning-method" → "named preset" wording. Swept `skills/`, `docs/`, `_dev/` and `README.md` for preset names, "learning-method", "teaching method" and preset counts. Only `ai-report-guide.md:124` and the test list name or count presets, and both are updated. No "learning-method" text is left outside the CHANGELOG history. `ai-report.md:22,81-83`, `ai-report-reference.md:91`, `skeleton.html:180,190`, the prompt (`{{EDITOR_DIRECTION}}` at :20), `help.md` (no second-opinion line) and the Go loader (folder glob, no list) agree. F3 is the only gap. Inherited at wave end: REQ-703's name shape, placeholder, header-before-backlink and finalize rule agree with v1c03 (page lines 125 and 148). REQ-705's bundle-then-shipped location, `preset_file` meaning and catalog walk agree (`report_second_opinion.go:197`, test :1026, `report_index.go`). REQ-733 recorded nothing. REQ-734 recorded nothing of its own, and its F4 line `help.md:18` now carries `--page`. Unread for the sweep: none.
**Anti-bloat:** 0 unnamed helpers, 0 options, 0 extra files, 0 decorative tests. Preset check 20 (open the page in a browser) is not named in the REQ, but it serves the D-05 limit, and the page's D-13 direction-block CSS lives in the generated page only.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*
