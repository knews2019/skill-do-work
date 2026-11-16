# Review: REQ-737 (relation-map preset keeps box titles readable after a click on a phone)

**Approve.** The preset change does what the REQ and the maintainer's mid-run addendum ask. My own browser run on the final page gives the same numbers as the builder's pass-3 hand-back, value for value.
Route B | range `38b278f2..3f6ebf38` (builder merge `8e6d94e9` plus catalog commit `3f6ebf38`) | orchestrated mode, standard review

## What's built

- The shipped preset `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` changes three checks. Check 15 adds a phone sentence: titles stay at about 12 px and a Show button reaches each related box. Check 12 keeps relation words inside the frame and off arrowheads. Check 17 puts the phone map frame in the top half of the first screen. The file grows from 3,057 to 3,426 bytes and still has checks 1 to 20 in order.
- One proof page, v1c04 `phone-first-map`, was built with that preset on the tstack bundle. Its reply, the bundle's iteration list line and the catalog pair (`ai-reports/catalog.json`, `ai-reports/index.html`) list it.
- Still missing: nothing the REQ names. The release (CHANGELOG, VERSION) is finalization's job and has not run yet. VERSION reads 0.324.1.

## Decisions / risks for you

- None block release. F1 below is the one visible rough edge left on the proof page. The maintainer approved captures that already show it (D-25).

## Browser re-measure (my run, not the hand-back's)

Server: `python3 -I -m http.server 8773 --bind 127.0.0.1` from the bundle folder. I ran the builder's script `do-work/runs/work-2026-11-16-171725/REQ-737-captures/measure-relation-map.js` on the published v1c04 page (Chromium 146, measured 2026-11-16T20:07Z). I also ran three small scratch scripts of my own: load and header fold, drag and zoom, and the empty space above column titles. Afterwards I stopped the server and deleted my PNG captures.

| Row | Mine 390 x 844 | Hand-back 390 | Mine 1440 x 1000 | Hand-back 1440 |
| --- | --- | --- | --- | --- |
| Map frame top y at load | 398 | 398 | 385 | 385 |
| Launcher header height / page content starts at y | 70 / 86 | 70 / 86 | 70 / 98 | 70 / 98 |
| Clicks (mouse + synthetic) | 33 (26 + 7) | 33 (26 + 7) | 33 (33 + 0) | 33 (33 + 0) |
| Map scale after a click | 0.56 every click | 0.56 | 0.453 to 0.75 | same |
| Smallest clicked title / smallest related title in frame | 12.32 / 12.32 px | 12.32 / 12.32 | 9.97 / 9.97 px | 9.97 / 9.97 |
| Related out of view, then Show step works | 70, then 70 | 70, then 70 | 0 | 0 |
| Clicked or related box out of frame | 70 (phone rule) | 70 | 0 of 33 clicks | 0 |
| Words on arrowheads / cut at frame / on boxes / broken | 0 / 0 / 0 / 0 | same | 0 / 0 / 0 / 0 | same |
| First view: titles, relation words | 12.93 px (D, O1, O2, O3), 3 of 12 | same | 12.41 px (10 boxes), 12 of 12 | same |
| Sideways scroll / Tab 23 of 23 in frame / list 28 = 28 rows / no-JS view | 0 / yes / yes / passes | same | same | same |
| `summary.green` | true | true | true | true |

Own checks:
- Header fold: closed at load at both widths. A click on the summary opens it, and the aside grows from 70 px to 1,966 px (phone) and 751 px (desktop). A second click closes it.
- Phone load (`own-390-load.png`, viewed): the short header, backlink, change list, title, map heading and tools sit above the frame. The frame starts at y = 398 with the column title "What is proposed" and boxes D, O1 and O2 visible. Maintainer target: y about 422 or less. Met.
- Phone click on O2 (`review-390-click-O2.png` and `own-390-click-O2.png`, viewed): D and O2 titles read at 12.3 px. The line "1 of 11 related boxes in view; the Show buttons in the detail panel bring the others into view" sits above the frame. "open: hold or hand back?" sits clear of the arrowhead stack, and no word touches the frame edge.
- Check 12 beyond clicks: I made 5 drags and 3 zoom-in presses per width, and 0 relation words were cut at the frame edge after any of them.
- The page's head meta and the header's fold both carry the committed preset text byte for byte (trailing newline aside). The page was built with the final preset.
- `TestSecondOpinionShippedPresetsLoadFromTheToolbox`: `ok` (0.713 s), run in the main tree at `3f6ebf38`.

Report-only facts from the brief:
- (a) This still holds on the final page, and it is wider than one click. Gap from the frame top to the column title after a phone click: O2 88 px, O1 155 px, D and O3 21 px, and after O2 then "Show S1" 223 px of a 620 px frame. The first view has 12 px. Recorded as F1.
- (b) REQ-738's header has no side gutter: "Iteration v1c04" and the summary triangle start at x = 0. This happens at 390 px and also at 1440 px (`asideLeft` 0 at both widths). It belongs to REQ-738 (its review F3), and the coordinator captures it as its own request. Recorded as F3, report only.

## Findings

**Important:** None.

**Minor:**
- F1: After a phone click or a Show tap, the framing centers the kept boxes vertically. Empty space then stays above the column title inside the frame: 88 px for O2, 155 px for O1, 223 px after O2 then Show S1. Only the first view pins the column to the frame's top (`index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html:583`, `view.ty = 12 - bb.y * view.s`). `frameSelection` (line 519) and `bringIntoView` (line 564) use `centerOn` only. No preset check names this. The maintainer's "less empty space above the columns" applies to the first screen, which meets it (12 px). — impact-user-visible → report only
- F2: The hand-back's proposed CHANGELOG entry is stale. It is titled `0.324.1`, which REQ-738 already released. The pass-2 text also says this release "removes both limits named in the 0.324.0 entry's Known limits line", but the second limit (map on the third phone screen) was removed by 0.324.1. The finalization entry should claim only the first limit (a phone click shrinks titles) and use the next free version. — impact-negligible → report only

**Report-only, belongs to REQ-738 (not scored here):**
- F3: REQ-738's short header has no side gutter on a phone or on desktop: x = 0 at 390 and 1440. This is REQ-738 review F3, and the coordinator captures it as a new request. — impact-user-visible → report only

**Nit:**
- F4: Check 15's phone sentence says "readable titles win" without naming the rule it overrides. The rule ("never leaves ... outside the frame") is the sentence just before it, so a reader finds it. The proof page applied it correctly. — impact-negligible → report only

**Anti-bloat count:** 0 extra helpers, options, files or decorative tests. The range touches exactly 6 files: the preset, the page, its reply, the bundle's `index.html` line and the catalog pair. The page's own script (framing, `slideFrom`, `bringFrameIntoView`, `bringIntoView`, Show buttons) is the expected proof-page code. The measurement script and probe are run artifacts under `do-work/runs/` and were not shipped.

## Requirements Checklist

- [x] Below 600 px, a click leaves the clicked title and the shown related titles readable without manual zoom: delivered, 12.32 px in 33 of 33 clicks (v1c03 had 4.52 px)
- [x] Check 15 extended for phones only, and related boxes that do not fit are reachable one way the preset names: delivered, Show buttons, 70 of 70 steps work
- [x] Desktop behavior of check 15 unchanged: delivered, first sentence byte-identical, 0 of 33 desktop clicks leave a box out of frame
- [x] The other 19 checks still hold, especially 17, 18 and 19: delivered (measured rows above)
- [x] Change stays in the preset text and near 3,057 bytes (at most about 3,500), no check removed: delivered, 3,426 bytes, 20 checks
- [x] One real page built with the changed preset and measured: delivered (v1c04, launcher finalize judge pass, standalone judge success per the hand-back)
- [x] Maintainer saw phone captures before done: delivered (D-25)
- [x] Addendum: map starts in the top half of the first phone screen with the short header: delivered, y = 398
- [x] Addendum: shorter intro on phones, less empty space above the columns at first view: delivered, intro moved below the map, 12 px gap at first view (after-click gap is F1)
- [x] Addendum: no relation label runs off the frame edge: delivered, 0 cut after every click and after drags and zooms
- [x] Addendum: form 1 and the Show taps kept: delivered
- [x] Constraints: only the preset changes among shipped files, v1c03 and other bundles untouched: delivered (`git diff --stat` shows 6 files)
- [ ] Release (CHANGELOG entry, VERSION bump): N/A at review, owned by finalization

## Restatement sweep

Redefined by this diff: preset checks 12, 15 and 17 (the phone click rule, words off arrowheads and inside the frame, the phone map frame in the top half). I grepped `skills/`, `_dev/primes/`, every `lessons-*.md` and both CHANGELOGs for `relation-map`, `relation map`, the phone limit wording and the check text:
- `skills/do-work-toolbox/docs/ai-report-guide.md:124`: names the preset in one clause about what it draws, with no check or phone wording. Agrees.
- `skills/do-work-toolbox/actions/ai-report.md`: no relation-map or click restatement.
- CLI help and Go code: no relation-map text outside `report_second_opinion_test.go` (the shipped-presets list, which only loads the file).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:18-19`: about measuring method and header height, not the checks. Agrees. REQ-738's F2 there (stale 514 / 1,302 px) is already resolved in line 18's current text.
- `CHANGELOG.md` / `skills/do-work/CHANGELOG.md` 0.324.0 "Known limits" still names the phone-click limit. That entry is released, so it is noted only, never rewritten. The next release entry should say it removes that limit (F2). The 0.324.1 entry's "y = 830" was measured on v1c03 and is a released historical figure. Noted only.

Inherited from REQ-738 (wave end): the header block's visible form (closed `details` with "Preset: <name> (full direction)"). The final v1c04 page carries that exact header, and its fold works at both widths. No stale restatement found.

## Acceptance Testing

**Result: Pass.** Stages covered: implementation and integration (the merged tree at `3f6ebf38`, served locally, measured at both widths, plus the shipped-presets Go test). I did not rerun the repository gate. The integrator's gate at `3f6ebf38` passed (REQ `## Testing`).

## Suggested Additional Testing

1. Deployment: unassessed. After finalization, check that the new CHANGELOG entry and its byte-identical mirror name only the phone-click limit as removed, under the next free version (F2).
2. Live acceptance: unassessed. Open the v1c04 page on a real phone (iOS Safari and Android Chrome) and click O2, then tap "Show S1". Confirm the 12 px titles and the empty space of F1 look acceptable on a real device.
3. Manual: open the page in dark theme at 390 px. The script measures the light theme only.

## Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every REQ and addendum requirement delivered; release owned by finalization |
| Code Quality | 90% | Clear proof-page code; after-click vertical centering leaves empty space (F1) |
| Test Adequacy | 90% | Kept script, RED on v1c03 and GREEN on v1c04; reproduced exactly; light theme only |
| Scope | 100% | 6 files, all declared; nothing extra |
| Risk | None | Prose preset plus one report page |
| Acceptance | Pass | Implementation + integration, measured at 390 and 1440 |

## Follow-ups created

None (4 findings report only)

---

## Review

**Overall: 95%** | 2026-11-16T20:10:57Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: After a phone click or Show tap, the framing centers vertically and leaves empty space above the column title in the 620 px frame: 88 px (O2), 155 px (O1), 223 px (O2 then Show S1). Only the first view pins to the top (`index-v1c04-…-phone-first-map.html:583`; `frameSelection` :519, `bringIntoView` :564 use `centerOn`). The first screen meets the addendum (12 px). — impact-user-visible → report only
- F2: The hand-back's proposed CHANGELOG entry is titled 0.324.1 (already REQ-738's) and claims both 0.324.0 limits. The release should claim only the phone-click limit. — impact-negligible → report only
- F3 (REQ-738's, not scored): the short header has no side gutter, x = 0 at 390 and 1440 (REQ-738 review F3). The coordinator captures it as a new request. — impact-user-visible → report only
- F4 (Nit): check 15's "readable titles win" leaves the overridden rule implicit. It is the sentence just before it. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages. My rerun of `measure-relation-map.js` at `3f6ebf38` matches the hand-back exactly: frame top y = 398 (phone) and 385 (desktop); 12.32 px titles in 33 of 33 phone clicks; 70 of 70 Show steps; 0 words on arrowheads or cut at the frame (also after drags and zooms); header fold closed at load and opens; green at both widths; shipped-presets Go test ok. Anti-bloat: 0 extras.
**Restatement sweep:** redefined relation-map checks 12, 15 and 17 (phone click rule, words off arrowheads and inside the frame, phone frame in the top half). No stale shipped restatement in `ai-report-guide.md:124`, `ai-report.md`, CLI help or the lessons files. The released 0.324.0 "Known limits" phone-click line and the 0.324.1 "y = 830" figure are noted only. Inherited from REQ-738: the header form; the v1c04 page agrees.
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*
