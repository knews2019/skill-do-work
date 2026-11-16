# Review: REQ-739 (the second-opinion header keeps a side margin on a phone)

**Approve** — the header text now starts 16 px from the edge on the kit page and the REQ-737 final page, and the tstack v1c03 page keeps its own look.
Route B | merge range `540925d4..ef5dc795` (RED `1e478a2a`, GREEN `6584e473`, merge `ef5dc795`)

## What's built

- `secondOpinionDecorate` writes one `<style>` inside the header markers, right before the `aside`: `:where(body>aside[aria-label="Second-opinion direction"]){padding-inline:16px}`. The rule has zero specificity, so any page rule for the header wins. The `aside` open tag is byte-identical to 0.324.2.
- One new test, `TestSecondOpinionHeaderKeepsAPhoneSideMargin`, pins the rule, the unchanged tag, idempotence, the strip of a literal 0.324.1 block, and a plain decorate that leaves no `<style>` behind.
- Nothing is missing against the REQ. The release files (CHANGELOG pair, VERSION) are the integrator's, as Scope says.

## Decisions / risks for you

None. The one trade-off (a `<style>` element in the body) was named in the REQ's Exploration and D-01, and it is recorded below as a Nit.

## Findings

**Important:** None

**Minor:** None

**Nit:**
- F1: the launcher's `<style>` sits in the body (`report_second_opinion.go:1645`), a place the HTML standard does not allow for `<style>`; every browser applies it and the judge does not check placement, but a strict HTML validator would flag every directed page. The REQ chose this on purpose over a head `<style>` with its own strip matcher (D-01). — impact-negligible → report only

Known report-only items from earlier runs (REQ-737 review F1, empty space above the column title after a phone click; REQ-738 review F1, fold arrow hidden on pages that hide `summary` markers) were not re-reported.

## Requirements Checklist

- [x] Header text about 16 px from the left edge at 390 px on a page with no own rule, with a matching right margin — delivered. Kit page and v1c04: text left 16 px, right gap 16 px (measured). A page with the default 8 px body margin gets a 24 px box edge (builder P3, D-02, accepted as "about 16 px").
- [x] A page that styles `body>aside[aria-label="Second-opinion direction"]` keeps its own margin and padding — delivered. v1c03: margin `12px 16px`, padding `8px 14px`, text left 31 px, identical before and after.
- [x] Decorating twice gives the same bytes, and a second decorate strips a 0.324.0 or 0.324.1 block — delivered. New test checks 3 and 4, REQ-738 test check 5 (0.324.0 literal); my harness also found both scratch pages idempotent.
- [x] No sideways scroll at 390 px, closed and open — delivered (measured). `ai-report-judge` passes — builder evidence (both scratch pages `pass`), not rerun here.
- [x] Only the header block changes; metas, markers, fold, presets and iteration list stay — delivered. The diff touches the `block :=` line and the function comment only. Removing the new `<style>` from my rebuilt v1c04 gives the published page's exact bytes.
- [x] Pages under `ai-reports/` not rewritten — delivered (the range has 2 files, both Go).
- [x] Captures shown to the maintainer before done — delivered per the REQ's Qualification (approved in a question dialog, 2026-11-16).
- [x] UR-181 (one small request for a side margin on the launcher's header) — delivered.

## Code Review

- **Code quality.** One string and a 3-line comment change. The comment says why the rule is zero-specificity and why an inline padding would be wrong. The `<style>` body is masked by `secondOpinionRawHTML`, so its text can never forge a marker, and the existing strip removes it with no new code. `gofmt -l` empty and `go vet` clean at `ef5dc795`.
- **Edge cases checked by reading.** With no `<body>` tag the block goes to offset 0; the parser then puts the `<style>` in the implied head and the `aside` as the first body child, so `body>aside` still matches. A page with a generic `aside{padding:…}` rule beats the launcher rule; that is the stated design ("the page's own rule wins"). `:where()` is in every current engine; an engine without it drops the rule and the header falls back to the 0.324.2 look.
- **Tests.** The new test failed on check 1 with the rule removed and passes with it (I reproduced both on a scratch copy). `go test -count=1 -run '^TestSecondOpinion'` on a scratch copy of the module at `ef5dc795`: `ok` (15.1 s). The builder's whole-package run and the integrator's gate (128 s, 988 fast tests) are green.
- **Scope.** 2 of 5 declared files touched; the other 3 are integrator-only release files. No undeclared file. Decisions D-01 to D-14 cover every visible choice.
- **Risk.** Low. Every directed page gets the new `<style>`; the strip, idempotence and old-block cases are pinned.
- **P-A-U.** All three boxes are ticked.

### Domain Review (frontend)

- Renders without errors: no console error from the pages (one favicon 404 from the local server only).
- Responsive: measured at 390 x 844 and 1440 x 1000; the rule has no media query (D-03), and the numbers hold at both widths.
- Accessible: the `aside` label and the native `details`/`summary` fold are unchanged; the fold opens and closes by click.
- No animation, no dependency added.

## Restatement Sweep

This REQ is the run's last successful integration, so the trigger set is this diff's element plus REQ-738's recorded element.

- **Redefined here:** what sits between the `<!-- second-opinion-header -->` markers (a `<style>` element now comes before the `aside`; the `aside` tag is unchanged).
- **Inherited from REQ-738:** the header block's visible form (`Preset: <name> (full direction)` summary of a closed `details`).
- **Grep:** `Second-opinion direction`, `second-opinion-header`, `header block`, `direction block`, `launcher … header`, `padding-inline`, `:where(`, and `<style>`-placement rules over the tracked tree outside `ai-reports/` and `do-work/` (covers `skills/`, `_dev/`, `suite/`, `docs/`).
- **Results:** `skills/do-work-toolbox/actions/ai-report.md:89` describes only the visible header (preset on top, full direction in a closed `details`), still true. `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html:8` and `skills/do-work-toolbox/actions/ai-report-reference.md:74` ("page-only CSS goes in one `<style>` block") address the editor's page CSS, not the launcher's block, so they agree. Test literals at `report_second_opinion_test.go:1057` and `:1102` are intentional old blocks; `:1074` (REQ-738 check 7) still matches the unchanged tag. `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:18` now reads "folded it to a header of about 87 px", which matches v1c03 (63.4 px box plus 12 px margin above and below = 87.4 px, unchanged). The preset files under `skills/do-work-toolbox/docs/second-opinion-presets/` say nothing about the header. CHANGELOG lines 351 and 372 are history.
- **Stale restatements found:** none.

## Acceptance Testing

**Result: Pass** (implementation and integration stages)

Pages were decorated by the merged `secondOpinionDecorate` through a throwaway Go test in a scratch copy of the module at `ef5dc795` (head metas reused; the rebuilt v1c04 minus the new `<style>` equals the published bytes). "Before" = the same page with only the new `<style>` removed. Chromium through Playwright, light theme; the window shows a 15 px scroll bar, so `clientWidth` is 375 at 390.

| Page | Width | Text left (before) | Summary left | Right gap | margin / padding | Box height | Scroll closed / open |
|---|---|---|---|---|---|---|---|
| Kit skeleton | 390 | 16 (0) | 16 | 16 | 0 / 0 16px | 75.6 (75.6) | 375/375 / 375/375 |
| Kit skeleton | 1440 | 16 (0) | 16 | 16 | 0 / 0 16px | 75.6 (75.6) | 1425/1425 / 1425/1425 |
| v1c04 (REQ-737 final) | 390 | 16 (0) | 16 | 16 | 0 / 0 16px | 69.6 (69.6) | 375/375 / 375/375 |
| v1c04 | 1440 | 16 (0) | 16 | 16 | 0 / 0 16px | 69.6 (69.6) | 1425/1425 / 1425/1425 |
| v1c03 (own rule) | 390 | 31 (31) | 31 | 30 (30) | 12px 16px / 8px 14px (same) | 63.4 (63.4) | 375/375 / 375/375 |
| v1c03 | 1440 | 43.5 (43.5) | 43.5 | 42.5 (42.5) | 12px 28.5px / 8px 14px (same) | 63.4 (63.4) | 1425/1425 / 1425/1425 |

- The fold opened on the first click and closed on the second on every page and width.
- Kit page, both widths: "Open all details" opened every kit `details` in its section and left the header fold closed; no sideways scroll after it.
- The summary line lines up with the iteration line everywhere (16/16, 31/31, 43.5/43.5).
- Header heights did not change, so the map top does not move.
- 390 px captures viewed: kit and v1c04 headers sit 16 px in; v1c03 shows its own bordered card unchanged. Captures deleted after viewing.
- The v1c03 box height is 63.4 px here and 87.4 px in the builder's matrix; the builder's figure includes the 12 px top and bottom margin. Both agree that nothing moved.

## Suggested Additional Testing

1. Deployment and live acceptance: unassessed. After the release, run one real directed pass and open the new page at 390 px to confirm the 16 px edge on a live page.
2. Open one decorated page in Safari on an iPhone to confirm the `:where()` rule applies there too.
3. Run `ai-report-judge` on the first real directed page after release; this review relied on the builder's scratch-bundle `pass`.

## Anti-bloat

0 items the REQ did not name: no new helper, option or file; one new test function, named in D-04. Two of its five checks overlap existing asserts (check 2 repeats REQ-738's check 7 literal; check 3 repeats the idempotence assert of `TestSecondOpinionDirectedRunNamesAndMarksPage` on a page with tree marks), and D-04 asked for both.

## Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every requirement and constraint delivered |
| Code Quality | 100% | One string and its comment; F1 is a Nit with no score weight |
| Test Adequacy | 95% | RED-then-GREEN reproduced; the browser proof lives outside Go tests by design |
| Scope | 100% | 2 declared files, nothing else |
| Risk | Low | Every directed page changes, but the strip and idempotence cases are pinned |
| Acceptance | Pass | Measured at 390 and 1440 on three pages |

## Follow-ups created

None (1 finding report only)

Adjacent observation, not measured: the iteration list the launcher writes on `index.html` (`<nav aria-label="Second-opinion iterations">`) has no side margin of its own either.

## Review

**Overall: 98%** | <integrator fills>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 (Nit): the launcher's `<style>` sits in the body (`report_second_opinion.go:1645`), where the HTML standard does not allow `<style>`; browsers apply it and the judge does not check, but a strict validator would flag every directed page; chosen on purpose in D-01 — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: focused `^TestSecondOpinion` tests green at `ef5dc795` and RED reproduced without the rule; merged decorate on scratch kit, v1c04 and v1c03 pages in Chromium (390 x 844 and 1440 x 1000: text left 16 px on kit and v1c04, 31 px / 43.5 px on v1c03 with its own margin `12px 16px` and padding `8px 14px` unchanged, header heights unchanged, no sideways scroll closed or open, fold opens and closes, "Open all details" leaves the header fold closed); judge from builder evidence; deployment and live acceptance unassessed
**Restatement sweep:** redefined the header block's content between the markers (a zero-specificity `<style>` before an unchanged `aside`); inherited REQ-738's header visible form; grep over the tracked tree outside `ai-reports/` and `do-work/` finds `ai-report.md:89`, `skeleton.html:8`, `ai-report-reference.md:74`, `lessons-do-work-cli.md:18` and the test literals all agree; no stale restatement
**Suggested testing:** 3 items
**Follow-ups created:** None (1 finding report only)

*Reviewed by review-work action*
