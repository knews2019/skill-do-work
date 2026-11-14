Confirmed staging output: [`index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html`](index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html)

### Clarity Changes
- **The map reads left to right (`#rev-1`):** four columns in the order of the argument: evidence (the live check and the seeds S1 to S9), the decision and the options O1 to O4, the kill attempts K1 to K4, and the limits L1 to L9. On a 1440-pixel-wide screen the map now uses the width of its frame instead of a tall column.
- **Nine boxes first (`#rev-2`):** the map opens on the decision, the four options and the four kill attempts at about 88% zoom, so every title reads at about 16 pixels. A number such as +8 on a box counts the evidence and limits behind it. A click shows them, a second click hides them, "Show everything" shows all 28 boxes and "Nine boxes" goes back to the start.
- **Readable line words (`#rev-3`):** the words of the selected box's lines sit in a layer that is not zoomed, so they stay 13 pixels at every zoom. Each label sits beside the other box of its line, and labels never overlap each other.
- **Light dimming (`#rev-3`):** a selected box dims the other boxes to 72% opacity instead of 30%, so their titles stay readable.
- **Kept:** drag, pinch and Control-scroll zoom, the zoom buttons and keys, the detail panel with sources and relations, every box as text, every relation as a table row, and the no-JavaScript view, which shows all 28 boxes and 37 lines at full size. The box texts, the 37 relations and the after-the-report section are the same as on the v1c03 page.

### Corrections to the Original Report
- No new corrections. The three corrections of the v1c03 page stay in "What happened after the report, and corrections": the footer's five images (the folder holds four), the direction text that is now also shown in a visible block, and the catalog page that lists only the bundles that existed at its last rebuild.
- **Not verified, as before:** every line reference into Tinka, every output of the live check, and the folded-text shares of 38 to 46%.

### Open Questions
- **Q1:** On a phone the map opens on the decision and the options at about 73% zoom. A click on O2 then shows 13 line labels on a 350-pixel-wide map, and they cover most of the boxes behind them. Should a narrow screen show fewer labels, for example only those in view?
- **Q2:** "Show everything" at 1440 pixels is at about 46% zoom, where titles read at about 12 pixels. Is that enough, or should the full view keep only codes and leave titles to the zoomed-in view?
- **Q3:** The limits L1, L6, L9 (to the live check) and L4, L5 (to S5) cross the whole map, because the limits column is on the far right. Should limits that concern the evidence sit beside it instead?
- **Q4:** The layout is still hand-made for this report. A relation-map preset needs a rule for any report: which kinds go in which column, and which boxes are shown first.

Variation slug: left to right map
