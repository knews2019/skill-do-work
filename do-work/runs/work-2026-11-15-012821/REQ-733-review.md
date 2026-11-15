# Independent review: REQ-733 (rebuild the relation map page: left to right, nine boxes first, readable labels)

## Review: REQ-733

**Approve with follow-ups.** The new page fixes the four problems on a wide screen at its first view and after one click. Two views still have relation labels lying over box titles: "Show everything" plus a click, and every click on a phone. Tab focus can also leave the visible frame.
Route B | merge `fd210a91` (range `c2f21c0d..fd210a91`)

### What's built
- A new second-opinion page, `index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html`, built on the v1c03 map. It has four columns, opens on the nine boxes, shows seeds and limits on click, and hides them on a second click. Labels sit in an unscaled overlay at 13 px, and unrelated boxes dim to 0.72.
- The page has a reply file. `index.html` gained one row inside its iteration list. No other file changed.

### Decisions / risks for you
- **Is this the standard option for REQ-727 (the relation-map preset)?** On a 1440-pixel screen the design works. At phone width a click covers the option titles with labels, which the builder's reply also asks about as Q1. Value of fixing F1 before REQ-727: the preset copies a label rule that works at every width. Risk of not fixing it: every map built by the preset carries the same phone defect.

### Findings

**Important:**
- F1. Labels cover box titles when the on-screen column gap is narrower than the label (`index-v1c03c01-…-left-to-right-map.html:389`). The label is at least 110 px wide, but the 210-unit gap is about 97 px at 46% and 88 px at 42%. Seen: at 1440×1000 after "Show everything", a click on O2 puts 11 labels over boxes. The first letters of the D, O1, O2, O3 and O4 titles are hidden ("lot O2 by hand", "and recipe, no ode"). K1 covers the K2, K3 and K4 titles. At 390×844 a click on O2 puts 14 labels over boxes, including O2's own title. The phone view is hard to read in light and dark. The builder checked "Show everything, then O3" only. Fix: seed `placed` in `placeLabels` with the on-screen rectangles of the shown boxes, so a label never lands on a box. A label with no free place is dropped, and the panel's relation list already holds every word. The simpler fix: skip the overlay when `(COLUMN_GAP - 14) * scale < 110`. — impact-user-visible → report only
- F2. Tab focus on an off-screen box leaves it off-screen (`:539` with `:322`). The browser scrolls the clipped viewport to the focused box before `focusin` measures it, so the box looks inside the frame and `frame()` is skipped. The scroll handler then resets the scroll to 0. Seen: at 149%, Tab to O3, O4 and K1 to K4 leaves each focused box 0% visible. Fix: in the `focusin` handler, set `viewport.scrollLeft = 0; viewport.scrollTop = 0;` before reading `getBoundingClientRect()`. — impact-user-visible → report only

**Minor:**
- F3. The panel keeps the last box after a deselect, and its reveal button shows a stale label (`:361-363`, `if (!detail) return;`). Seen: select O2, press Escape, then click "Hide its evidence and limits" three times. The map goes 9 → 17 → 9 boxes, but the button still reads "Hide…". After "Show everything" with nothing selected, the panel still offers "Hide its evidence and limits". Fix: when nothing is selected, restore the panel's first "Click a box" content, saved at script start. — impact-user-visible → report only
- F4. Without JavaScript, the zoom, "Nine boxes" and "Show everything" buttons show and do nothing. `.map-tools{display:flex}` (`:19`) beats the `hidden` attribute (`:107`). This is inherited from v1c03. Fix: add `.map-tools[hidden]{display:none}`. — impact-negligible → report only
- F5. The no-JavaScript map lost the words on its lines. v1c03 drew them as SVG text, but the overlay needs script. The relation table still lists all 37 relations, so no content is lost. — impact-negligible → report only
- F6. Labels pushed out of the frame are removed without notice (`:411`). At 30% zoom O2 shows 10 of its 13 labels, and at 24% it shows 9. The panel lists them all. — impact-negligible → report only

**Nit:**
- F7. Two `aria-label` texts are inherited from v1c03: "Live check: Live check: one lineage…" (`:159`) and "…Build nothing yet.. the decision." (the D box). — impact-negligible → report only

### Requirements Checklist
- [x] Left to right in argument order: four columns, at x=40 (P, S1-S9), x=550 (D, O1-O4), x=1060 (K1-K4) and x=1570 (L1-L9). Delivered.
- [x] The first view shows nine boxes at 1440: 9 of 28 at 88%, titles 15.8 px, all fully in the frame. On a phone the first view frames D and O1-O4 at 73% (titles 19.8 px), by decision D-11, and K1-K4 sit off to the right. Delivered, with a documented deviation on narrow screens.
- [x] Reveal on click and hide again: O2 goes 9 → 17 boxes, a second click goes back to 9, and P then shows L1, L6 and L9. Delivered.
- [x] One show-all control: "Show everything" (28 boxes at 46%), with "Nine boxes" to reset. Delivered.
- [~] Labels readable at every zoom: 13 px measured at 24% to 220%, so the size is delivered. Placement fails in some views: labels cover boxes in the full view and on a phone (F1). Partial.
- [x] Light dimming: opacity 0.72 against 0.3 in v1c03, and dimmed titles stay readable in both themes. Delivered.
- [x] Kept features: drag (a drag that starts on a box does not select it), Control-scroll zoom up to 240%, the + / − / 0 / arrow / Escape keys, the detail panel with sources and relations, 28 box texts and 37 relation rows identical to v1c03 by diff, and the no-JS map with all 28 boxes. Delivered, except F2 and F5.
- [x] No library: only the bundle's own `presentation.js`. Delivered.
- [x] No file outside the three: the diff stat shows 3 files, and the GREEN probe prints `REQ-733 probe PASS` (12 earlier files byte-identical, `index.html` changed only inside the tree markers). Delivered.
- [x] Browser check recorded by the builder, judge `pass` recorded at finalize. Delivered.
- [x] Table form not built (a constraint). Delivered.
- [x] P-A-U boxes all checked.

### Acceptance Testing

**Result: Partial** (stages covered: implementation and integration, in a real Chromium against the merged bytes served on 127.0.0.1:47734)
- **1440×1000, light:** the first view shows 9 boxes at 88%, titles 15.8 px, 0 box overlaps, no page overflow. O2 shows 17 boxes at 58%, 13 labels at 13 px, 0 labels on boxes, and O1, O4 and K3 at 0.72 and readable. K1 shows 3 labels, 13 boxes dimmed and readable. "Show everything" shows 28 boxes at 46% with 12.5 px titles and nothing clipped. "Show everything" plus O2 gives 11 labels on boxes (F1). "Nine boxes" returns to 9 at 88%.
- **1440×1000, dark:** the same numbers. "Show everything" plus K1 gives 4 labels over the K2, K3 and K4 titles (F1). Contrast is good.
- **390×844, light and dark:** the first view shows D and O1-O4 at 73% with 19.8 px titles, and no page overflow. O2 gives 17 boxes at 42% with 11.3 px titles and 14 labels over boxes, so the map is hard to read. K1 gives 4 labels over boxes, and "Show everything" does the same.
- **Keys:** + and − step from 58% to 220% and down to 24%, labels stay 13 px, 0 fits again, and the arrows move the map. Tab leaves the map after the last box, so there is no focus trap. Tab at high zoom loses the focused box (F2).
- **Links:** `#box-L1` on load opens P and then L1 (20 boxes) and selects L1. The panel's O2 link and the K1 link in "After the report" select their box and scroll the map into view.
- **No-JS:** all 28 boxes at 100%, scrollable inside the frame, 28 box texts and 37 relation rows. The tool bar shows dead buttons (F4).
- **Base v1c03 at 1440, "Show everything" plus O2:** 36% zoom, 10.1 px titles, 6.1 px labels, unrelated boxes at 0.3. The improvement claim holds: 46%, 12.5 px titles, 13 px labels, 0.72.

### Suggested Additional Testing
- Live acceptance: unassessed. The maintainer opens the page on their own screen and decides on REQ-727.
- Manual check on a real phone: pinch zoom and drag by touch. Only mouse and wheel were exercised here.
- Keyboard and screen-reader pass after F2: Tab through the boxes at 150% and check what is announced on reveal.
- Safari: the box titles use `-webkit-line-clamp`, and the label overlay uses `color-mix`.

### Scores (on the record, not the headline)

**Overall: 74%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 90% | All delivered on desktop. Label placement is partial (F1). The phone first view is D-11 by design. |
| Code Quality | 75% | Clear and small. Two real state bugs (F2, F3) and an undersized gap rule (F1). |
| Test Adequacy | 75% | The probe pins the byte constraints. The browser check missed the full view plus O2, and the phone overlap that the reply names. |
| Scope | 95% | Only the 3 files. The UI elements the REQ did not name are listed under Anti-bloat. |
| Risk | Low | A local, export-ignored bundle page. No shipped file. |
| Acceptance | Partial | Desktop works end to end. Labels cover boxes in the full view and on a phone. Tab focus can leave the frame. |

The average is (90 + 75 + 75 + 95) / 4 = 83.75. The Partial acceptance takes 10 points off, which gives 74%.

### Anti-bloat
- Files: 0 extra. Tests: 0 decorative (the probe is a run artifact). Libraries: 0.
- Controls the REQ did not name: 3. They are the "Showing N of 28 boxes" counter, the "+N" badge on boxes, and the panel button "Show/Hide its evidence and limits". The panel button repeats the click-again toggle and is the source of F3. The four v1c03 group toggles were removed (D-09, documented).
- Behaviours the REQ did not name: 4. They are the 42% zoom floor (`READABLE`), the phone start on the D and O column, the shortest-path auto-reveal for links, and focus-follow. All four are documented in decisions D-08, D-09 and D-11.
- Script functions: 15 (v1c03 had about 8). Each serves a named requirement or one of the behaviours above.

### Follow-ups created
None (7 findings report only). Proposed follow-up for the maintainer to capture if wanted: fix F1, F2 and F3 on this page before REQ-727 copies its label rule. Each fix is a few lines in the inline script.

---

## Review

**Overall: 74%** | 2026-11-15T00:13:37Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 75% |
| Test Adequacy | 75% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
- F1. Relation labels cover box titles when the on-screen column gap is narrower than the 110 px minimum label width (`index-v1c03c01-…-left-to-right-map.html:389`). At 1440 "Show everything" plus O2, 11 labels sit on boxes, over the D and O1-O4 titles. Every phone click puts 4 to 14 labels on boxes. Fix: seed `placed` with the shown box rectangles, or skip the overlay when `(COLUMN_GAP - 14) * scale < 110`. — impact-user-visible → report only
- F2. Tab focus on an off-screen box leaves it off-screen. `focusin` (`:539`) measures after the browser scrolled the clipped viewport, and then `:322` resets the scroll. At 149%, Tab to O3..K4 leaves each focused box 0% visible. Fix: reset the scroll to 0 before measuring in `focusin`. — impact-user-visible → report only

**Minor findings:** F3. After a deselect, the panel keeps the last box and its reveal button label goes stale (`:361-363`). Fix: restore the "Click a box" content when nothing is selected. — impact-user-visible → report only. F4. The no-JS tool bar shows dead buttons, because `.map-tools{display:flex}` (`:19`) beats `hidden` (`:107`). This is inherited. — impact-negligible → report only. F5. The no-JS map has no line words, but the relation table keeps them. — impact-negligible → report only. F6. Labels pushed off-frame are dropped without notice (`:411`), so O2 at 30% shows 10 of 13. — impact-negligible → report only. Nit F7. Inherited duplicated `aria-label` text (`:159`, D box). — impact-negligible → report only
**Acceptance:** Partial. Implementation and integration stages were checked in Chromium at 1440×1000 and 390×844, light and dark. The first view, O2, K1, "Show everything", "Nine boxes", the keys, links and no-JS all work. Labels cover titles in the full view plus a click and on a phone, and Tab focus can leave the frame.
**Restatement sweep:** nothing redefined (a bundle page; skipped by instruction, not the wave's last integration)
**Suggested testing:** 4 items
**Follow-ups created:** None (7 findings report only)
**Rendered page:** 1440 light and dark: 9 boxes at 88% (15.8 px titles). O2 gives 17 boxes at 58% with 13 labels, none on a box. "Show everything" gives 28 at 46% (12.5 px titles), but a click there puts labels over titles. 390 light and dark: D and O1-O4 at 73% (19.8 px). A click goes to 42% (11.3 px titles), and labels cover most option boxes. Base v1c03 at 1440 "Show everything": 36%, 10.1 px titles, 6.1 px labels, 0.3 dim.
**Anti-bloat:** 0 extra files, 0 decorative tests, 0 libraries. 3 controls the REQ did not name (box counter, "+N" badge, panel reveal button). 4 behaviours the REQ did not name, all in D-08, D-09 and D-11.

*Reviewed by review-work action*
