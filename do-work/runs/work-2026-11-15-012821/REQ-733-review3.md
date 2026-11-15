# Independent second delta re-review: REQ-733 (rebuild the relation map page: left to right, nine boxes first, readable labels)

## Review: REQ-733 (second delta re-review)

**Approve.** The fix commit `3eb2ab4b` closes G1, G2, G3 and G4 from the second review. F1, F2, F3, F4 and F6 from the first review stay fixed. In a real Chromium I found 0 words broken inside a label, 0 labels on a box and 0 label overlaps in every view, at every zoom from 20% to 240%. A click never moved a related box out of a view that already showed it. One small new finding stays: H1, an uneven zoom after a manual zoom-out.
Route B | merge `b6a1e1a7` (delta `567903a5..b6a1e1a7`, full range `c2f21c0d..b6a1e1a7`)

### What's built
- The delta changes only the inline script of `index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html` (17 lines added, 14 removed).
- A label is now 110 to 190 px wide (`:400`). The taken-space rule drops a label that has no free place, and the count line reports it.
- A click zooms in to `LABEL_SCALE` (0.575) only when the box and its related boxes all fit at that zoom (`:466-468`). Otherwise the builder's framing rule applies: move only when something is off screen.
- `setCount` (`:431`) writes the live count line only when its text changes. "line words" became "relation words". `fitIds` lost its unused floor parameter.

### Decisions / risks for you
- **Phone behavior (decision D-16, then D-18).** At 390 px a click lands at 42% and shows 0 to 3 relation words, with the count note "N of M relation words in view". The panel lists every relation. This meets the REQ: every label shown is whole and readable, and no relation is lost. Value: no label covers a title. Risk: a phone reader drags to find the words.

### Findings

**Minor:**
- H1. After a manual zoom-out, a click zooms in for some boxes and not for others (`:466-469`). The zoom-in runs only when the related boxes fit at 57.5%. When they do not, the view stays at the reader's zoom, even when a closer zoom would still fit them all. Seen at 1440×1000, "Show everything", − twice (27%): a click on S2, K1 or O2 jumps to 57%. A click on P stays at 27% with 0 of 4 relation words, and O3 stays at 27% with 2 of 9. From the 46% "Show everything" view, P shows 1 of 4 and O3 shows 4 of 9. Fix: compute `fitted` once (as in `fitIds`), and zoom in to `Math.min(LABEL_SCALE, fitted)` whenever `selected && scale < LABEL_SCALE && fitted > scale && fitted >= READABLE`. — impact-user-visible → report only
- F5 (carried). The map without JavaScript has no relation words on its lines. The relation table keeps all 37 relations. — impact-negligible → report only

**Nit:**
- H2. `fitsForWords` (`:466-467`) repeats the fit test from `fitIds` (`:447`) with a literal `48` in place of `2 * pad`. If one changes, the two drift. A shared `fittedScale(ids)` helper would serve both, and H1's fix. — impact-negligible → report only
- F7 (carried, inherited from v1c03). Duplicated text: "Live check: Live check: …" in the `aria-label` at `:160` and in 4 panel relation links at `:214`, and "…Build nothing yet.. the decision." at `:170`. v1c03 has the same counts. — impact-negligible → report only

### Fixed findings, re-tested in the browser
- G1 fixed: label width 110 px at every zoom below 60%. 0 broken words with O2, O3 or P selected while stepping − from 57% to 20% at 1440 and from 42% to 20% at 390. 0 broken words during Control-scroll from 20% to 240% and with + up to 240% (labels up to 176 px wide). At 20% to 37%, 1 or 2 labels stay, whole and readable. The count line names the rest.
- G2 fixed: I clicked every shown box whose related boxes were all in the frame before the click. 1440: 9 in the first view, 28 in "Show everything" (46%), 17 with O2 open, and 28 at 27%. 390: 2, 1, 2 and 13. Every click kept the box and all its related boxes in the frame (0 failures in 100 clicks). "Show everything" plus P, L1 and L6 stays at 46% with their related boxes in the frame. K1 and O2 zoom in to 57% with all related boxes in the frame.
- G3 fixed: a 20-step drag wrote the count line once (when one label left the frame). 20 arrow-key pans that end where they started wrote it 0 times.
- G4 fixed: the count line says "relation words". The help text at `:205` says "the words of its lines". That is a sentence, not a term, so it does not conflict.
- F1 fixed: 0 labels on boxes and 0 label overlaps in every measured view.
- F2 fixed: at 193%, Tab moves D, O1-O4 and K1-K4, and each focused box is 100% in the frame. The next Tab leaves the map.
- F3 fixed: Escape after O2 returns the panel to "Click a box" with no reveal button.
- F4 fixed: without JavaScript the tool bar computes `display: none`. All 28 boxes and 38 table rows show, and the page has no horizontal overflow at either width.
- F6 fixed: dropped labels are reported as "N of M relation words in view".

### Requirements Checklist
- [x] Left to right in argument order: four columns (P and S, D and O, K, L). Delivered.
- [x] First view shows nine boxes: 1440 shows 9 of 28 at 88% with 15.8 px titles. 390 shows D and O1-O4 at 73% with 19.8 px titles (D-11). Delivered.
- [x] Reveal on click, hide again, one show-all control: O2 goes 9 → 17 boxes. The panel button goes "Hide its evidence and limits" → "Show its evidence and limits (8)" → "Hide…" with 9 → 17 → 9 boxes. "Show everything" gives 28 at 46%, and "Nine boxes" returns to 9 at 88%. Delivered.
- [x] Labels readable at every zoom: 13 px text, 110 to 190 px wide, no word broken, no label on a box, from 20% to 240%. Labels that do not fit are dropped and counted. Delivered.
- [x] Light dimming: unrelated boxes stay readable in light and dark (captures checked). Delivered.
- [x] Kept features: drag (a drag that starts on a box keeps the selection), Control-scroll, + / − / 0 / arrows / Escape, the panel, the hash link, the no-JS map. Delivered.
- [x] No library, nothing outside the bundle: the delta touches 1 file, and `git status --short -- ai-reports` is empty. Delivered.
- [x] Browser check recorded. The judge `pass` covers the builder bytes (`ed6a3a8a`), and the fixed bytes cannot be judged one page at a time yet (D-17). Delivered with the D-17 gap.

### Acceptance Testing

**Result: Pass** (stages covered: implementation and integration, in Chromium against the merged bytes served on 127.0.0.1:47737)
- **1440×1000, light and dark (same numbers):** first view 9 boxes at 88%. O2: 17 boxes at 58%, 13 of 13 labels. K1: 3 of 3 at 88%. O3: 12 boxes at 49%, 6 of 9 labels, all 9 related boxes in the frame. P (after O2): 20 boxes at 49%, 1 of 4 labels, all related boxes in the frame. "Show everything": 28 boxes at 46%, 12.5 px titles. Then O2 and K1 go to 57% with all labels, and P and L1 stay at 46% with all related boxes in the frame. "Nine boxes" returns to 9 at 88%.
- **390×844, light and dark:** first view D and O1-O4 at 73%. O2 and O3 land at 42% with 3 of 13 and 3 of 9 labels. K1 and P land at 42% with 0 labels and the count note. "Show everything" stays at 42%, centered on O2. No page overflow.
- **Keys and zoom:** − down to 20%, Control-scroll from 20% to 240%, + up to 240%, and "0" refits. No broken words, no labels on boxes.
- **Links:** a cold load of `#box-L1` in a new page selects L1 and shows 20 boxes. At 1440 it lands at 49% with P in the frame. At 390 it lands at 42% centered on L1, with P outside the frame and "0 of 1 relation words in view" (D-11).
- **No-JS:** tool bar hidden, 28 boxes, 38 table rows.
- The integrator's numbers match mine, except its "20 arrow-key pans write the count line 3 times". My 20 pans that return to the start wrote it 0 times. Both results show the G3 fix works.

### Suggested Additional Testing
- Live acceptance: unassessed. The maintainer opens the page on their own screen, then decides on REQ-727 (the relation-map preset).
- On a real phone: drag and pinch by touch, and check that the count note helps a reader find the hidden words.
- A screen reader during a drag, to confirm the count line is announced once.

### Scores (on the record, not the headline)

**Overall: 89%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 92% | Every criterion delivered. A phone shows few relation words per click, as recorded in D-16 and D-18. |
| Code Quality | 88% | A small fix with accurate comments that closes four findings. H1 is an uneven edge, and H2 is a repeated formula. |
| Test Adequacy | 82% | The integrator's view numbers match mine. It did not try a click after a manual zoom-out (H1). No automated test, which fits a bundle page. |
| Scope | 95% | One file, and the delta stays inside G1-G4. |
| Risk | Low | A local, export-ignored bundle page with no shipped file. |
| Acceptance | Pass | Implementation and integration stages pass at both sizes and in both themes. |

The average is (92 + 88 + 82 + 95) / 4 = 89.25, which rounds to 89%. Acceptance Pass takes no points off.

### Anti-bloat
- Files: 0 extra. Tests: 0 decorative. Libraries: 0. New controls: 0.
- `setCount`: one line, and it serves G3.
- `fitsForWords`: needed for G2. It repeats the fit formula (H2).
- The delta removes one `fitIds` parameter, so it shrinks the code it touches.

### Follow-ups created
None (4 findings report only). Optional for the integrator, about 4 lines: H1 and H2 together through one `fittedScale(ids)` helper.

---

## Review

**Overall: 89%** | 2026-11-15T00:43:48Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 88% |
| Test Adequacy | 82% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
- None

**Minor findings:** H1. After a manual zoom-out, a click zooms in to 57% only when the related boxes fit there, so at 27% S2 jumps to 57% but P stays at 27% with 0 of 4 relation words (`index-v1c03c01-…-left-to-right-map.html:466-469`). Fix: zoom in to `Math.min(LABEL_SCALE, fitted)` whenever `scale < LABEL_SCALE && fitted > scale && fitted >= READABLE`. — impact-user-visible → report only. F5 (carried). The no-JS map has no relation words on its lines, but the relation table keeps all 37. — impact-negligible → report only. Nit H2. `fitsForWords` (`:466-467`) repeats the `fitIds` fit test with a literal 48. Fix: one `fittedScale(ids)` helper. — impact-negligible → report only. Nit F7 (carried, inherited). "Live check: Live check: …" at `:160` and in 4 panel links at `:214`, and "yet.. the decision" at `:170`. — impact-negligible → report only
**Acceptance:** Pass. Implementation and integration stages checked in Chromium at 1440×1000 and 390×844, light and dark. G1-G4 are fixed, F1-F4 and F6 stay fixed. There are 0 broken words, 0 labels on boxes and 0 label overlaps from 20% to 240%, and 0 of 100 clicks moved a related box out of a view that showed it.
**Restatement sweep:** nothing redefined (a bundle page, skipped by instruction because this is not the wave's last integration)
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)
**Rendered page:** 1440 light and dark: 9 boxes at 88% (15.8 px titles). O2 gives 17 boxes at 58% with 13 of 13 labels. K1 gives 3 of 3. O3 gives 6 of 9 at 49%. "Show everything" gives 28 at 46%. Then O2 and K1 go to 57% with every label, and P and L1 stay at 46% with every related box in the frame. 390 light and dark: D and O1-O4 at 73% (19.8 px). Each click lands at 42% with 0 to 3 whole labels and the count note. No page overflow at either size.
**Anti-bloat:** 0 extra files, tests, libraries or controls. `setCount` serves G3 and `fitsForWords` serves G2. One `fitIds` parameter was removed.

*Reviewed by review-work action*
