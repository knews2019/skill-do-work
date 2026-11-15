## Review: REQ-734

**Approve.** The table page delivers every requirement, and all 37 relations match the v1c03 rows word for word and end for end. One correction on the page went false when the branch merged.
Route B | `d31159be..0ccdf98d` (builder commit `6d50d565`, merge `c812ea4b`, catalog rebuild `0ccdf98d`)

### What's built
- A new page `index-v1c04-claude-claude-opus-5-5-custom-supports-limits-table.html` has the decision first, then one table row per option (O1 to O4) with Supports and Limits columns, then the eight relations that touch no option. It has no script, no zoom and no drag.
- The hand-back shows the map page and the table page side by side at 1440 and 390, with both paths and the counts. It picks no winner.
- Still missing: nothing the REQ asked for. The "Corrections" box on the page says something that is false at HEAD (F1).

### Decisions / risks for you
- None for this REQ. The map-versus-table choice is yours. It decides what REQ-727 (ship a relation-map preset) becomes. See "Map vs table" below.

### Findings

**Important:**
- F1. The page's third correction is false at HEAD. `index-v1c04-…-supports-limits-table.html` (closing "Corrections" list) and `v1c04-…-reply.md:15` both say that `ai-reports/catalog.json` has no Tinka bundle ("grep -c -i tinka: 0", "neither Tinka proposal from 2026-11-14 is in it"). That was true at the builder's base `18bd280a`, where the count is 0. It was already false at the merge base `d31159be`, where the count is 10, and at HEAD `0ccdf98d` the count is 11. Both Tinka bundles and the v1c04 page itself are now listed. A reader of the page is told to distrust a catalog that is in fact complete. Fixing it means editing the page (and the reply) and judging it again with `--page` — impact-user-visible → report only (follow-up REQ)

**Minor:**
- F2. The page's second correction cites `ai-report.md:87` for the visible direction block. That sentence is now at `skills/do-work-toolbox/actions/ai-report.md:89`: it was at 87 at `18bd280a`, and REQ-728 and REQ-729 added lines above it. Line 87 now holds the second-opinion `latest` rule. K4's chip `ai-report.md:87` is copied from v1c03 content (the REQ's D-02 rule), so it is not the builder's error — impact-negligible → report only

**Nit:**
- F3. The hand-back has two small accuracy problems. Its first line is a bare page path before the `#` heading. It also says that at 390 "the first screen of both pages holds the launcher block, the change list and the start of the first section". On the map page the first section starts at y = 1023, below the 844-px screen, so that first screen ends at the sticky top bar. The builder's picture shows this correctly — impact-negligible → report only
- F4 (restatement sweep). `skills/do-work-toolbox/actions/help.md:18` reads `ai-report judge <dir>   Render-check a report bundle (widths, themes, links)`. It omits the `--page <file.html>` form that REQ-728 added. The guide's command list (`skills/do-work-toolbox/docs/ai-report-guide.md:80`) and the action (`skills/do-work-toolbox/actions/ai-report.md:36`) both show that form. The help line is incomplete, not wrong — impact-negligible → report only (follow-up REQ)

### Requirements Checklist

- [x] The decision sits above the table, with one row per option O1 to O4. The decision block is `#rev-1`, and the row headers are `<th scope="row">` for `row-O1` to `row-O4`. Delivered.
- [x] There are two main columns, Supports and Limits. Each entry has the relation (both end codes and the words), the other box's code and v1c03 title, and one source chip (D-09). I checked every title and chip against the first chip of the v1c03 `#rev-2` box text. Delivered.
- [x] The relations that touch no option are shown in "Relations that touch no option" (`#rev-3`). The 8 entries are grouped by Live check (L1, L6, L9), S5 (L4, L5), Decision (L2, L7) and S9 (L3), exactly the D-03 list. Delivered.
- [x] All 37 relations are present and correct. My script read the 37 v1c03 `#rev-3` rows and the 43 v1c04 entries: 37 distinct ids, plus 6 repeats of the order relations (D-11). For every entry, the words, the From code, the To code and the kind class match the row with the same number. 0 mismatches and 0 missing. Delivered.
- [x] Column placement follows D-03. All 8 "supports" relations sit in the Supports cell of the option they touch (r12 `O4 → K1` sits in O4 Supports). All 14 "limits or blocks" relations and the 3 option-touching gaps (r34, r35, r36) sit in the Limits cell. The 4 order relations sit in the decision block and are repeated in the row header of each option they touch. No relation is missing from a row it touches. Delivered.
- [x] There is no zoom or drag, and the page reads top to bottom. It has no `<script>`, no canvas and no transform. My browser run measured scroll width 1425 at 1440 and 375 at 390, so nothing overflows, and the rows stack at 390. Delivered.
- [x] The hand-back shows the pictures side by side for each width, with both paths and the counts, and picks no winner. The four PNGs are 1440×1000 and 390×844 and live in the run folder, not the bundle. Delivered.
- [x] The finalize judge passed, and so did the re-judge after the D-12 CSS edit. Every other bundle file keeps its bytes: I ran the probe and it printed `REQ-734 probe PASS`, exit 0. `index.html` gained one `<li>` inside the iteration list. Delivered.
- [x] Nothing under `skills/` or `suite/` changed, so this is no release. `.gitattributes:52` marks `/ai-reports` export-ignore. The REQ cites line 48, which is stale, but harmless. Delivered.
- [x] UR-178 asks to compare the two pages. The pictures, the counts and the paragraph below do that. Delivered.

### Acceptance Testing

**Result: Pass** (implementation and integration stages; no deployment stage, because nothing ships).
- I ran `REQ-734-probe.sh` at HEAD: PASS.
- My Python script compared the 37 relations against v1c03 `#rev-3`, checking the words, both end codes, the kind and the column. Result: 0 problems.
- Browser check: I served the bundle with `python3 -I -m http.server` on 127.0.0.1 and ran Playwright Chromium at 1440×1000 and 390×844 in light, with a dark glance at 1440. My first-screen count (words plus both end codes fully in the viewport, not clipped, at 12 px or more, opacity above 0.5) gives: map 0/37 and 0/37, table 4/37 (r01 to r04) and 0/37. **I agree with the builder's four numbers.** Dark mode on both pages looks clean, with readable contrast.
- I checked the page's three corrections and its 0.316.0 / 0.317.0 note against the repository. The screenshot count (4 files) is true. The 0.316.0 / 0.317.0 note matches `CHANGELOG.md:82-96`. The catalog claim is false (F1). The line citation drifted (F2).

### Suggested Additional Testing
- Live acceptance: unassessed. The maintainer reads both pages and picks the format. That choice is the REQ's real outcome, and it unblocks REQ-727.
- Manual check: read the table at phone width in a real phone browser. The layout measured clean, but the table is about 9.5 screens long at 390.

### Restatement sweep (wave end; this REQ is the run's last successful integration)
- This diff redefines nothing shipped. All five changed files are under the export-ignored `ai-reports/`.
- Inherited elements, taken from earlier members' `**Restatement sweep:**` lines:
  - REQ-723 (archive-collision check ignores evidence files in assets folders): the scan set of the archive-collision command. No shipped restatement of the scan set exists, so nothing is stale.
  - REQ-724 (crew files carry the rules for status explanations, user decisions and pasted commands): the DECIDE & STATE pointer, the clear-questions option shape, the per-agent model rule, the numeric-target and redesign lines, and the stall-check wording. The `skills/do-work-toolbox/crew-members/` and `skills/do-work-knowledge/crew-members/` copies of `background-agents.md`, `clear-questions.md` and `coding-guardrails.md` differ from `skills/do-work/crew-members/` only in relative link paths, apart from one older gap noted below. `communication-style.md` has one home. Consistent.
  - REQ-732 (allow docs/run-prompts/ and move the run prompts there): the allowed `docs/` folders and the home of the run prompts. No shipped file mentions `_dev/run-prompts/`. `skills/do-work-toolbox/actions/folder-structure.md:38` lists `run-prompts/` under `docs/`, and `skills/do-work/crew-members/communication-style.md:28` points there. Consistent.
  - REQ-728 (render check takes one named page, `--page`) and REQ-730 (catalog lists each bundle's second-opinion pages): both recorded "nothing redefined", but I checked the candidates the brief named. The `--page` form appears in the usage text (`report_judge.go:41`), `ai-report.md:36` and `:190`, and `ai-report-guide.md:59` and `:80`, all consistent. The one gap is `help.md:18` (F4). `second_opinions` appears in `report_index.go:81`, `ai-report.md:48` and `ai-report-guide.md:105`, all consistent.
  - REQ-731 (new reports remind you to rebuild the catalog): recorded "nothing redefined". `ai-report.md:208` and `ai-report-guide.md:102` agree with `report_index.go:546-551`. `stakeholder-report.md:51-54` and `architecture-report.md:119-122` run the organize check without the catalog check. That is REQ-731's own D-03 out-of-scope decision and was already listed as a discovered task, so it is not a stale restatement.
  - REQ-733 (rebuild the relation map page): nothing redefined.
- Unread for the sweep: REQ-729, REQ-726 and REQ-725. Their lines say "not run" or "skipped" and record no element. I still checked REQ-729's rule because the brief named it: the usage text and `ai-report.md:186` both say a relative `--out` starts at `--repo-root`. `stakeholder-report.md:44` and `architecture-report.md:100` pass no `--out`. Consistent.
- Stale shipped restatements found: one, F4 (`skills/do-work-toolbox/actions/help.md:18`).
- Outside this run: the toolbox `coding-guardrails.md` lacks the "Earned defense" paragraph (`skills/do-work/crew-members/coding-guardrails.md:56`). That gap dates from REQ-170 (`84c1b81d`), not from this wave.

### Anti-bloat
The REQ did not name any of these, and the count of each is 0: helpers, options, files and decorative tests. The repository has no new test. The probe is a run artifact the procedure calls for. The closing "Corrections" section was not named in the REQ, but the second-opinion editor prompt asks for corrections (`skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md`), so it is not bloat. F1 lives in that section. The `.nw` nowrap span is a CSS class that fixes broken words (recorded in the hand-back), not a helper.

### Map vs table
At 1440×1000, the table page shows more relations on the first screen: 4 of 37 (the order of the options under the decision), against 0 of 37 on the map page. The map page's first screen holds the launcher block, the change list and the page title with its intro. The map itself starts at y = 1143, just below the screen. It opens on nine boxes. A line's words appear only after a click on a box, and the seed and limit boxes appear only after a click or "Show everything". The box texts sit in five collapsed groups. All 37 relations can be read without a click only in the depth-3 "Every relation" table, about 4 to 6 screens down (the page is about 6 screens long). The table page hides nothing behind a click or in a collapsed group. All 37 relations are readable by scrolling alone, within about 4 screens. What it leaves out entirely is the full box texts and the extra source chips: it keeps one line and one chip per box and links to `index.html` for the rest. At 390×844, neither page shows a relation on the first screen (0 and 0). The table page's first screen reaches the decision heading, and its order relations start on the second screen. Its rows stack, and all 37 relations are readable by vertical scrolling over about 9.5 screens. The map page's first screen ends at its top bar, and its title starts on the second screen. The map is about three screens down, and it shows a 1910-px-wide drawing through a 331-px window, so reading it needs panning or zooming. Its "Every relation" table sits in a sideways-scrolling box (557 px of table in 335 px), so most rows need a sideways swipe to show the To code. Without that swipe, my script read only 3 of 37. The map page is about 11.5 screens long. I agree with the builder's counts: map 0/37 and 0/37, table 4/37 and 0/37.

### Scores (on the record — not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every REQ and UR requirement delivered; all 37 relations match |
| Code Quality | 85% | Clean, script-free page; one claim stale at merge (F1), one citation drifted (F2) |
| Test Adequacy | 90% | Probe pins bytes and ids; judge pass; browser matrix; red-green at base |
| Scope | 100% | Declared (after D-17) equals touched; D-12 edit recorded |
| Risk | None | Export-ignored report files only |
| Acceptance | Pass | Implementation and integration stages exercised |

### Follow-ups created
None (4 findings report only)

---

## Review

**Overall: 94%** | 2026-11-15T03:26:41Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Verdict:** Approve

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1. The page's and the reply's third correction ("ai-reports/catalog.json has no Tinka bundle, grep -c -i tinka: 0") is false at HEAD. The count was 0 at the builder's base `18bd280a`, 10 at merge base `d31159be` and 11 at `0ccdf98d`, where both Tinka bundles and v1c04 are listed. Files: `index-v1c04-…-supports-limits-table.html` "Corrections" list and `v1c04-…-reply.md:15` — impact-user-visible → report only (follow-up REQ)

**Minor findings:**
- F2. The page's direction-block correction cites `ai-report.md:87`. That sentence is now at `skills/do-work-toolbox/actions/ai-report.md:89`, because REQ-728 and REQ-729 added lines above it — impact-negligible → report only
- F3 (nit). The hand-back's first line is a bare path before its heading. It also says the 390 first screen of the map page shows "the start of the first section", but the map page's title starts at y = 1023, below the 844-px screen — impact-negligible → report only
- F4 (nit, restatement sweep). `skills/do-work-toolbox/actions/help.md:18` lists `ai-report judge <dir>` without the `[--page <file.html>]` form that `ai-report-guide.md:80` and `ai-report.md:36` show since REQ-728 — impact-negligible → report only (follow-up REQ)

**Acceptance:** Pass — implementation and integration stages. The probe passes. All 37 relations match v1c03 `#rev-3` in words, both end codes and D-03 column. In a real browser at 1440×1000 and 390×844 (light, with a dark glance), my first-screen counts agree with the builder's: map 0/37 and 0/37, table 4/37 and 0/37.
**Restatement sweep:** nothing redefined (this diff touches only export-ignored `ai-reports/` files). At the wave end, I swept the inherited elements of REQ-723, REQ-724 and REQ-732 and the brief's named candidates for REQ-728, REQ-729, REQ-730 and REQ-731 in `skills/` and `suite/`. One stale restatement: F4 (`skills/do-work-toolbox/actions/help.md:18`). Unread for the sweep: REQ-729, REQ-726, REQ-725 (their lines record no element).
**Anti-bloat:** 0 helpers, 0 options, 0 extra files, 0 decorative tests that the REQ did not name. The closing corrections section comes from the second-opinion editor prompt.
**Map vs table:** At 1440, the table page shows 4 of 37 relations on the first screen (the order of the options under the decision), and the map page shows 0. The map starts just below the first screen, opens on nine boxes, and shows a line's words only after a click on a box. Its seeds and limits appear only after a click or "Show everything", and its box texts sit in five collapsed groups. Its relations can be read without a click only in the depth-3 table, about 4 to 6 screens down. The table page hides nothing behind a click or in a collapsed group, and all 37 relations are readable by scrolling within about 4 screens. What it leaves out entirely is the full box texts and extra source chips, which are linked to `index.html`. At 390, both pages show 0 relations on the first screen. The table's relations start on the second screen and stack in about 9.5 screens. The map page's title starts on the second screen, and its drawing about three screens down, 1910 px wide in a 331-px window, so reading it needs panning or zooming. Its relation table needs a sideways swipe on most rows. The map page is about 11.5 screens long. I agree with the builder's counts.
**Suggested testing:** 2 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Delta re-review

**Overall: 96%** | 2026-11-15 | **Verdict: Approve** | range `d31159be..13322db2`, delta `0ccdf98d..13322db2`

The delta is one file with 2 lines changed, both in the page's "Corrections" list. The builder commit is `e760b5b5`, re-merged as `13322db2`.

- F1: fixed. The catalog correction now says that when the page was written, neither Tinka proposal was in the catalog, and that a later rebuild added both. Its chip names `18bd280a`. I checked that `git show 18bd280a:ai-reports/catalog.json | grep -c -i tinka` prints 0. The text is now true at HEAD.
- F2: fixed. The chip now reads `ai-report.md:89`, and line 89 at `13322db2` holds the "header block between `<!-- second-opinion-header -->` markers" sentence. The K4 chips `ai-report.md:87` (page lines 108 and 109) are the original report's own citations: `index.html` at `18bd280a` carries them twice. Leaving them unchanged follows the REQ's D-02 rule.
- Nothing else moved. My relation script still finds all 37 relations with 0 mismatches, and the probe passes at HEAD (`REQ-734 probe PASS`). The catalog pair is unchanged since `0ccdf98d`. The coordinator reports the re-judge as pass and the gate at exit 0; I did not re-run either.
- D-18, which leaves `v1c04-…-reply.md:15` unchanged: accepted. The reply is the launcher-published record of what the editor said at its base, and rewriting it would falsify that record. The reader page now carries the true statement.
- `help.md:18`: I agree with the coordinator's judgment. The line omits an option rather than stating something false, so F4 stays a nit, impact-negligible → report only (discovered task).
- Open: F3 and F4 (both nits, report only).

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Follow-ups created:** None (2 open findings report only; F1 and F2 fixed)
