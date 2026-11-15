# Independent delta re-review: REQ-733 (rebuild the relation map page: left to right, nine boxes first, readable labels)

## Review: REQ-733 (delta re-review)

**Request changes: one more small fix.** The fix commit `4be874a1` closes F1, F2, F3, F4 and F6 from the first review. In every default view I measured, no label covers a box and no two labels overlap. The fix itself adds two new defects. Below about 53% zoom, labels get narrower than their longest word, so words break in the middle. The 55% zoom floor after a click also moves a view that already showed everything, and it pushes the related boxes out of the frame.
Route B | merge `567903a5` (delta `fd210a91..567903a5`, full range `c2f21c0d..567903a5`)

### What's built
- The delta changes only the inline script of `index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html` and adds one CSS rule (29 lines added, 11 removed).
- Shown boxes now count as taken space for labels. Label width follows the on-screen column gap. A click zooms in to at least 55% (`LABEL_SCALE`). A label pushed more than one box height from its line end is dropped. The count line says "N of M line words in view" when labels are dropped.
- `focusin` resets the viewport scroll before it measures (F2). The panel restores its "Click a box" content when nothing is selected (F3). `.map-tools[hidden]{display:none}` hides the dead tool bar without JavaScript (F4).

### Decisions / risks for you
- **Is the phone behavior acceptable?** At 390 px wide, every click shows 0 labels ("0 of 13 line words in view" for O2). Dragging 130 px shows 6 of 13 labels, all readable and none on a box. I judge this acceptable as decision D-16's answer to the reply's Q1, because the panel lists every relation and the count line says how many words are hidden. Value: no label ever covers a title. Risk: a phone reader sees lines leave the frame with no words until they drag. G2's fix makes this a little better (the "0" key at 42% already shows 8 of 13).

### Findings

**Important:**
- G1. Labels get narrower than their longest word below about 53% zoom, so words break in the middle (`:401`, `const labelWidth = Math.min(190, Math.floor(COLUMN_GAP * scale) - 12);`). The first review's rule kept 110 px as the smallest width. The new rule gives 85 px at 46%, 76 px at 42%, 50 px at 30% and 30 px at 20%. Seen: "Show everything", then O2, then − (46%): "recomme|nds", "AGENTS.|md", "comment|s". On a phone, O2 then "0" (42%): "recomm|ends", "compute|r", "compare|d". At 20% each surviving label is one letter per line, about 20 lines tall. These zooms come from −, Control-scroll, pinch, and the "0" key on a phone. The REQ says "Relation labels stay readable on screen at every zoom level". Fix: `Math.min(190, Math.max(110, Math.floor(COLUMN_GAP * scale) - 12))`. The taken-space rule then drops a label that does not fit, and the count line already reports it. — impact-user-visible → report only
- G2. The 55% floor moves a view that already showed the box and every related box, and it pushes those related boxes out of the frame (`:465-466`). `frame()` used to move only when something was off screen. Now it also refits at 55% and centers on the clicked box when the related boxes do not fit at 55%. Seen at 1440×1000, after "Show everything" (46%, all 28 boxes in the frame): a click on P leaves 0 of its 4 related boxes in the frame (O2, L1, L6, L9) and shows 1 of 4 labels. L1 and L6 each lose P, with 0 of 1 labels. A cold load of `#box-L1` also lands at 55% with P off screen and 0 of 1 labels. K1 jumps to 75%, not 55%, because the largest zoom stays `Math.max(scale, 0.75)`, and 15 boxes end up cut at the frame edge. On a phone the floor gives 0 labels after every click, while 42% shows 8 of 13. Fix: raise the zoom to `LABEL_SCALE` only when the box and its related boxes fit at that zoom, and then zoom only to `LABEL_SCALE`: `if (!inside || (scale < floor && fitsAt(ids, floor))) fitIds(ids, inside ? floor : Math.max(scale, 0.75), floor, id);`. Here `fitsAt` is the `fitted` test from `fitIds`. — impact-user-visible → report only

**Minor:**
- G3. The count line is an `aria-live="polite"` region, and the script now rewrites it on every `apply()`, including every drag and zoom frame (`:394`, `:429`). Before the fix it changed only in `render()`. Seen: 20 text changes during a 20-step drag, and its text was the same before and after. A screen reader may announce it again and again. Fix: build the string first, and assign it only when it differs from `mapCount.textContent`. — impact-user-visible → report only
- F5 (carried). The map without JavaScript has no line words. The relation table keeps all 37. — impact-negligible → report only

**Nit:**
- F7 (carried). Two duplicated `aria-label` texts inherited from v1c03 ("Live check: Live check: …" at `:160`, and "…Build nothing yet.. the decision." on the D box). — impact-negligible → report only
- G4. "line words" on the count line is a new two-word term. The map's help text says "the words of its lines". "8 of 13 relation words in view" would match the panel's "Relations" heading. — impact-negligible → report only

### Fixed findings, re-tested in the browser
- F1 fixed in every default view: 0 labels on boxes and 0 label overlaps in all 32 measured views (2 sizes × 2 themes × 8 states). 1440 "Show everything" plus O2 now shows 13 of 13 labels at 58%. Before the fix, 11 of them sat on boxes.
- F2 fixed: at 172%, Tab moves through D, O1-O4 and K1-K4, and each focused box is 100% in the frame. Tab after K4 leaves the map, so there is no focus trap. With K1 selected, Shift+Tab through O4-O1 keeps each focused box 100% in the frame.
- F3 fixed: O2, then Escape, shows "Click a box" with no reveal button. "Show everything" with nothing selected does the same, and so does a click on an empty part of the map.
- F4 fixed: without JavaScript the tool bar computes `display: none`. All 28 boxes and 38 table rows are there, and there is no page overflow.
- F6 fixed: dropped labels now show as "N of M line words in view".

### Requirements Checklist
- [x] Left to right in argument order: four columns (P and S, then D and O, then K, then L). Delivered.
- [x] The first view shows nine boxes: at 1440, 9 of 28 at 88% with 15.8 px titles, all in the frame. On a phone, D and O1-O4 at 73% with 19.8 px titles (decision D-11). Delivered.
- [x] Reveal on click and hide again, and one show-all control: O2 goes 9 → 17 boxes. The panel button goes 17 → 9 → 17 with a correct label ("Show its evidence and limits (8)" / "Hide its evidence and limits"). "Show everything" gives 28 at 46%, "Nine boxes" gives 9 at 88%. Delivered.
- [~] Labels readable at every zoom: 13 px at every zoom. At 55% and above every shown label is whole and off the boxes. Below about 53% words break in the middle (G1). Partial.
- [x] Light dimming: 0.72 opacity in both themes, and the dimmed titles stay readable. Delivered.
- [x] Kept features: drag (a drag that starts on a box keeps the selection), Control-scroll (20% to 240%), + / − / 0 / arrows / Escape, the panel, the hash link, the no-JS map. Delivered, except the behavior change in G2.
- [x] No library, no file outside the three: the delta touches 1 file, and `git status --short -- ai-reports` is empty. Delivered.
- [x] Browser check recorded, judge `pass` for the builder bytes (`ed6a3a8a`). The fixed bytes cannot be judged one page at a time yet (decision D-17). Delivered with the D-17 gap.

### Acceptance Testing

**Result: Partial** (stages covered: implementation and integration, in a real Chromium against the merged bytes served on 127.0.0.1:47736)
- **1440×1000, light and dark (same numbers):** the first view shows 9 boxes at 88%, 15.8 px titles, no title cut off, no page overflow. O2 shows 17 boxes at 58%, 13 of 13 labels 109 px wide (up to 6 lines), 0 on boxes. K1 shows 3 of 3 labels at 88%. O3 shows 12 boxes at 55%, 8 of 9 labels, and the count line reads "Showing 12 of 28 boxes · 8 of 9 line words in view". "Show everything" shows 28 boxes at 46% with 12.5 px titles. "Show everything" plus O2 gives 58% and 13 of 13 labels, with S8 cut at the frame edge and 10 boxes outside the frame. "Show everything" plus K1 gives 75% with 15 boxes cut at the frame edge (G2). "Show everything" plus P leaves all 4 related boxes outside the frame (G2). "Nine boxes" returns to 9 at 88%.
- **390×844, light and dark:** the first view shows D and O1-O4 at 73% with 19.8 px titles, and no page overflow. O2, K1 and O3 each go to 55% with 14.9 px titles and 0 labels ("0 of 13", "0 of 3", "0 of 9 line words in view"). "Show everything" gives 42% (11.3 px titles), and a click then gives 55% with 0 labels. A 130 px drag shows 6 of 13 labels for O2, 103 px wide and readable. "0" at 42% shows 8 of 13 labels, but 4 of them have words broken in the middle (G1).
- **Keys and zoom:** + steps to 69% and 86% (labels 132 and 168 px wide), − goes down to 20% (labels 30 px wide, one letter per line, G1), and Control-scroll covers 20% to 240%. The labels never sit on a box.
- **Links:** a cold load of `#box-L1` selects L1, shows 20 boxes, and puts L1 in the frame at 55%. P, its only related box, is outside the frame (G2).
- **No-JS:** the tool bar is hidden, with 28 boxes and 38 table rows.

### Suggested Additional Testing
- Live acceptance: unassessed. The maintainer opens the page on their own screen and decides on REQ-727 (the relation-map preset).
- After G1 and G2 are fixed: rerun "Show everything" plus P, L1 and K1 at 1440, and − down to 20% with O2 selected.
- On a real phone: drag and pinch by touch, and check that the count line helps a reader find the hidden words.
- A screen reader during a drag, to confirm G3 and then its fix.

### Scores (on the record, not the headline)

**Overall: 72%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 88% | Every criterion delivered at the default views. Readable labels are partial below about 53% zoom (G1). |
| Code Quality | 75% | A small, readable fix that closes five findings. It brings two behavior regressions (G1, G2) and the live-region chatter (G3). |
| Test Adequacy | 70% | The integrator's numbers match what I measured for the views it named. It did not try zoom below 55%, or a click in "Show everything" on a box whose related boxes are far apart. |
| Scope | 95% | One file, and the delta stays inside the findings. |
| Risk | Low | A local, export-ignored bundle page, with no shipped file. |
| Acceptance | Partial | Default views are clean at both sizes and in both themes. Labels break words below about 53% zoom, and the floor hides related boxes. |

The average is (88 + 75 + 70 + 95) / 4 = 82. The Partial acceptance takes 10 points off, which gives 72%.

### Anti-bloat
- Files: 0 extra. Tests: 0 decorative. Libraries: 0. New controls: 0.
- `LABEL_SCALE`: one named constant with a comment that explains the number. It earns its place, but its use in `frame()` causes G2.
- The two new `fitIds` parameters: `smallest` carries the floor and is needed for the rule. `centerId` makes a Tab-focused box center on itself instead of on O2 when nothing is selected, so it earns its place too.
- The count suffix: it replaces the silent drop of F6, so it earns its place. It causes G3 and uses a new term (G4).
- Two new state variables (`panelStart`, `countText`) and the "strayed" drop rule: each one serves a named finding.

### Follow-ups created
None (6 findings report only). Proposed fix for the integrator, about 4 lines in the inline script: G1 (restore the 110 px smallest label width at `:401`) and G2 (apply the 55% floor only when the box and its related boxes fit at 55%, `:465-466`). Optional: G3.

---

## Review

**Overall: 72%** | 2026-11-15T00:30:34Z

| Dimension | Score |
|-----------|-------|
| Requirements | 88% |
| Code Quality | 75% |
| Test Adequacy | 70% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
- G1. Below about 53% zoom the label width (`index-v1c03c01-…-left-to-right-map.html:401`) is narrower than the longest word, so words break in the middle ("recomme|nds" at 46%, "compute|r" at 42% on a phone), and at 20% each label is one letter per line. Fix: `Math.min(190, Math.max(110, Math.floor(COLUMN_GAP * scale) - 12))`, and the taken-space rule drops what does not fit. — impact-user-visible → report only
- G2. The 55% floor in `frame()` (`:465-466`) moves a view that already showed everything and pushes the related boxes out: "Show everything" plus P leaves 0 of 4 related boxes in the frame, L1, L6 and a cold `#box-L1` lose P, and K1 jumps to 75%. On a phone it gives 0 labels after every click. Fix: raise the zoom to `LABEL_SCALE` only when the box and its related boxes fit at that zoom, and then only to `LABEL_SCALE`. — impact-user-visible → report only

**Minor findings:** G3. The `aria-live` count line is rewritten on every drag and zoom frame (20 changes in a 20-step drag) (`:394`, `:429`). Fix: assign only when the text changes. — impact-user-visible → report only. F5 (carried). The no-JS map has no line words, but the relation table keeps all 37. — impact-negligible → report only. Nit F7 (carried). Duplicated `aria-label` texts (`:160`, D box). — impact-negligible → report only. Nit G4. "line words" is a new term, and "relation words" would match the panel. — impact-negligible → report only
**Acceptance:** Partial. Implementation and integration stages were checked in Chromium at 1440×1000 and 390×844, light and dark. F1, F2, F3, F4 and F6 are fixed. Every default view has 0 labels on boxes, but labels break words below about 53% zoom (G1), and the 55% floor hides related boxes after a click in "Show everything" (G2).
**Restatement sweep:** nothing redefined (a bundle page, skipped by instruction because this is not the wave's last integration)
**Suggested testing:** 4 items
**Follow-ups created:** None (6 findings report only)
**Rendered page:** 1440 light and dark: 9 boxes at 88% (15.8 px titles). O2 gives 17 boxes at 58%, 13 of 13 labels, 0 on boxes. K1 gives 3 of 3 labels. O3 gives 8 of 9 labels at 55%. "Show everything" gives 28 at 46%, then O2 gives 58% and 13 of 13 labels. "Show everything" plus K1 gives 75% with 15 boxes cut, and plus P gives 0 of 4 related boxes in the frame. 390 light and dark: D and O1-O4 at 73% (19.8 px). Every click gives 55% and 0 labels, a 130 px drag shows 6 of 13, and "0" at 42% shows 8 of 13 with 4 words broken. No page overflow at either size.
**Anti-bloat:** 0 extra files, tests, libraries or controls. `LABEL_SCALE`, `smallest`, `centerId` and the count suffix each serve a named finding. `LABEL_SCALE`'s use in `frame()` causes G2, and the suffix causes G3.

*Reviewed by review-work action*
