## Review: REQ-738

**Approve**: a directed page now opens with a short header, and the full direction sits in a closed fold. On a 390 x 844 phone the relation map starts inside the first screen.
Route B | 943c4385 (merge of the builder branch; reviewed range f995337a..943c4385)

### What's built
- Every directed second-opinion page (every preset, `custom`, tree and plain) shows the iteration line with its parent link, then one line "▸ Preset: <name> (full direction)". That line is the `summary` of a closed `details` that holds the full direction with its line breaks.
- Markers, head metas, strip calls and the insert position did not change. A page with a 0.324.0 header gets the new header on its next decorate, and decorating twice gives the same bytes.
- The two contract sentences describe the fold: `skills/do-work-toolbox/actions/ai-report.md:89` and `skills/do-work-toolbox/docs/ai-report-guide.md:124`.

### Decisions / risks for you
- The map frame top is y = 830 on an 844 px phone screen, so only 14 px of the frame shows at load. That meets the user's words ("must not push the map below the first screen"). The page's own intro and change list cause most of the remaining distance (767 px without any header), not the launcher. Value: the launcher header now costs 87 px, down from 1,302 px. Risk: none from this diff. A taller first screen would need a page-layout change in the preset (REQ-737's area), not here.
- The header is 87.4 px (63.4 px box plus 24 px margins from the tstack page's own rule), above the capture's estimate of "about 80". D-03 (accept up to about 100 px when the map top is under 844) covers this, and numeric targets are approximate.

### Findings

**Important:**
- None.

**Minor:**
- F1. `report_second_opinion.go:1636`: the header's inline styles cancel the kit's card rules, but not a page rule that hides the summary marker (`summary{list-style:none}`, `::-webkit-details-marker`). 17 published files under `ai-reports/` hide the marker for their own folds. An editor can copy that rule into a directed page, and the header fold then shows no ▸. The words "(full direction)" and the pointer cursor still hint that it opens, and Tab plus Enter still work. A partial fix is `display:list-item;list-style:revert` inline on the `summary`. — impact-user-visible → report only
- F2. Restatement: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:18` still says in the present tense that "the launcher's direction block, printed above the page at finalize, moves the map down (514 px at 1440, 1,302 px at 390)". After this change the block is 87 px. The lesson's advice (measure again after finalize) still holds. Proposal for the integrator: add "until 0.324.x (REQ-738 folds it to about 87 px)" when it records this REQ's lesson. — impact-negligible → report only
- F3. Confirmed pre-existing gap, already in the REQ's Discovered Tasks: on a page with no rule for the header (report kit, plain page) the header has no side gutter. Measured `asideLeft` 0 px on the kit page at 390, and the text touches the screen edge (my capture `kit-390.png`). 0.324.0 has the same gap, so this diff does not cause it. — impact-user-visible → report only

**Nit:**
- F4. `TestSecondOpinionHeaderFoldsTheFullDirection` does not pin the inline styles that cancel the kit's `details`/`summary` card. A later edit that drops them passes the Go test and brings back a card of about 60 px on kit pages. Only a browser check catches it today. — impact-negligible → report only

### Requirements Checklist
- [x] Short header at the top: preset name or `custom`, plus the iteration line with the parent link on a tree page — delivered (`report_second_opinion.go:1632-1637`; test checks 2 and 6; browser: summary on one line at both widths).
- [x] The full direction is no longer an always-open block, and it stays readable on the page — delivered (closed `details`; fold text 3,056 characters = the trimmed 3,057-byte preset; `white-space:pre-line` kept).
- [x] Head metas unchanged — delivered (test check 3, a literal six-meta string; my scratch decorate of the real v1c03 page gives six metas byte-equal to the published ones).
- [x] Idempotent, and a 0.324.0 block is stripped — delivered (test checks 4 and 5; on the real v1c03 page I checked decorate-twice equality and that a plain strip gives the same bytes as a strip of the published page).
- [x] Every preset and `--editor-direction` — delivered (one shared slot branch; check 6 covers `bold` and `custom`; I decorated a kit page with `bold` and a plain page with editor direction only).
- [x] Existing header tests updated where needed — delivered as "not needed": both existing tests pass unchanged (D-08), and none asserts the old visible words.
- [x] Report kit: no card look, not toggled by "Open all details" — delivered (measured; see Acceptance).
- [x] Long-token guard kept — delivered (`overflow-wrap:anywhere` on the `aside`, check 7; a 400-character token gives 0 px sideways scroll with the fold closed and open).
- [x] Contract text updated — delivered (both sentences).
- [x] Map frame top under 844 at 390 x 844; header about 80 px — delivered within D-03 (87.4 px, y = 830).
- [ ] CHANGELOG and VERSION — N/A in this range (written at finalization, D-05).
- [x] Phone capture shown to the maintainer — recorded in the Qualification trace (approved 2026-11-16).

### Acceptance Testing

**Result: Pass** (stages: implementation and integration)
- `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands/` on the main tree at 943c4385: ok in 21.7 s.
- I decorated scratch pages with the merged `secondOpinionDecorate`. A throwaway test file in a detached worktree at 943c4385 did the decorate, and the worktree was removed afterwards. Pages: (S1) the published tstack v1c03 relation-map page, rebuilt from its metas; (KIT) `skeleton.html` with `presentation.css`/`presentation.js` beside it, preset `bold`, tree v1c01; (PLAIN) a bare page, editor direction only, with a 400-character token, `&`, `<b>` and a parent tree page. Pages were served from 127.0.0.1 and measured in Playwright Chromium, light theme.
- S1 at 390 x 844: header 87.4 px (box 63.4), map frame top y = 830, 0 px sideways scroll, summary on one line. At 1440 x 1000: header 87.4 px, map top y = 402, 0 px sideways scroll. Fold opened: 1,326.7 px / y = 2,069 at 390 and 538.5 px / y = 853 at 1440, still 0 px sideways scroll.
- Keyboard at both widths: Tab stop 1 is the parent link and stop 2 is the summary. Enter opens the fold and Space closes it.
- JavaScript off (390): the header is 87.4 px closed (map top 748). Enter on the focused summary opens the fold, and a click closes it. The same holds on the kit page.
- Kit page at both widths: the header `details` has a transparent background, 0 px border, 0 px radius and 0 px padding. The summary has weight 400, inherits 17 px, has 0 px padding and no bottom border when open. Header 75.6 px. "Open all details" opened the kit's three sections (`true,true,true`) and left the header fold closed.
- PLAIN: header 75.5 px at 390 and 57 px at 1440. 0 px sideways scroll with the fold closed and open.
- No page errors. The render judge was not re-run by me; the builder's judge run on the scratch S1 bundle passed.

### Suggested Additional Testing
- Deployment: unassessed. Install the next release into a consumer copy and confirm `ai-report second-opinion` ships the new header.
- Live acceptance: unassessed. Run a real `ai-report second-opinion --preset relation-map` pass after REQ-737 lands. Then measure the finalized page at 390 x 844 (map frame top below 844) and run `ai-report-judge` on it.
- Manual: view one directed page in Safari on iOS, where the `details` marker and the tap target of the `summary` render differently from Chromium.
- Edge case: a page whose own CSS hides summary markers (F1). Check that readers still find the fold.

### Scores (on the record — not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every requirement delivered; header size within the approximate target per D-03 |
| Code Quality | 95% | One-branch change, comment kept in sync; marker hiding by page CSS not cancelled (F1) |
| Test Adequacy | 90% | RED-before-GREEN shown (975066d9 → 152d6bbb); kit-cancel styles unpinned (F4) |
| Scope | 100% | 4 of 8 declared paths touched, all declared; the rest are finalization-only or not needed (D-08) |
| Risk | Low | Shared decorate path for every directed page; markers, metas and strips unchanged |
| Acceptance | Pass | Implementation + integration, browser-measured at both widths |

### Follow-ups created
None (4 findings report only)

### Self-validation
- I checked error and edge paths, not only the happy path: JavaScript off, keyboard, a long token, a 0.324.0 block on the real page, a plain strip, and the kit button.
- Regression consumers: the catalog and the iteration list read head metas only, and the metas are byte-equal. The strip-only call at `report_second_opinion.go:840` depends on the markers, which are unchanged. The published pages' own scripts that toggle `details` are scoped to a section or depth (`#rev-4 details`, `section.querySelectorAll`), so none opens the header fold.
- Proposed lesson for the integrator: the builder's `[family: session-is-the-backend]` line in the hand-back is accurate. Pair it with the F2 amendment to line 18 of the same file.

### Anti-bloat
Helpers, options, files or decorative tests the REQ did not name: none. The new test is the one D-04 named. Its two closures (`shipped`, `header`) serve only that test.

## Review

**Overall: 96%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: page CSS that hides summary markers (17 published files under `ai-reports/` do this for their own folds) would also hide the header fold's ▸; the inline styles do not cancel it (`report_second_opinion.go:1636`) — impact-user-visible → report only
- F2: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:18` still states the 0.324.0 block size (514 px / 1,302 px) in the present tense; now about 87 px — impact-negligible → report only
- F3: confirmed pre-existing gap from the REQ's Discovered Tasks: no side gutter for the header on kit and plain pages (`asideLeft` 0 px on the kit page at 390) — impact-user-visible → report only
- F4 (Nit): the new test does not pin the kit-cancelling inline styles on `details`/`summary` — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: focused Go tests green at 943c4385; merged decorate on scratch S1/kit/plain pages measured in Chromium (390 x 844: header 87.4 px, map top y = 830; 1440 x 1000: header 87.4 px; 0 px sideways scroll closed and open; Tab+Enter opens; works with JS off; kit card cancelled and "Open all details" leaves the header closed); deployment and live acceptance unassessed
**Restatement sweep:** redefined the header block's visible form ("Second-opinion preset:" / "Direction:" paragraphs → "Preset: <name> (full direction)" summary of a closed details); grep over `skills/`, `_dev/`, `suite/` finds the two contract sentences updated in this diff and the test's intentional 0.324.0 literal; one stale restatement at `lessons-do-work-cli.md:18` (F2)
**Suggested testing:** 4 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*
