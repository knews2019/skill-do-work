Confirmed staging output: [`index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`](index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html)

### Clarity Changes
- **One zoomable map of the whole report (`#rev-1`):** 28 boxes (the decision, the live check, options O1 to O4, kill attempts K1 to K4, seeds S1 to S9, limits L1 to L9) and 37 relation lines. Zoomed out, a box shows its code and title. Zoomed in, it shows a summary. A click opens the detail, the sources and the relations in a panel, and fades every box that is not related.
- **Relations by kind:** each line has one of four colors: order between options, supports, limits or blocks, gap in the evidence. The words of a line show when its box is selected or when the map is zoomed in far.
- **Every box as text (`#rev-2`):** the same content in folds, for reading in order, for search and for a screen reader. The panel shows these texts.
- **Every relation in one table (`#rev-3`):** one row per line of the map.
- **After the report (`#rev-4`):** the changelog shows that the canvas offer shipped in 0.316.0 and was removed in 0.317.0, because a real canvas showed each page as a thin unreadable strip. The page links this to K1.

### Corrections to the Original Report
- **Footer, "The five images in `screenshots/`":** the folder holds four images and the page shows four (`ls screenshots/`).
- **O1, "only the direction text survives, as a meta tag":** the new page also shows the direction in a visible block at its top (`skills/do-work-toolbox/actions/ai-report.md:87`).
- **O1, "The catalog page lists every bundle":** it lists the bundles that existed at its last rebuild; neither proposal about Tinka from 2026-11-14 is in `ai-reports/catalog.json`.
- **Not verified:** every line reference into Tinka, every output of the live check, and the folded-text shares of 38 to 46%. Their sources are outside the files this pass reads.
- **This page's own reading, not the original's:** the 37 relation lines. The original states them in sentences; this page turned each into one line and one label. The box texts are shortened from the original.

### Open Questions
- **Q1:** Is the map readable enough when everything is shown? On a 1440-pixel-wide screen the full view is at about 30 to 36% zoom, where titles are small and only the codes read easily.
- **Q2:** The layout is hand-made for this report (columns for seeds, options, kill attempts). A report with other content needs its own layout, or an automatic one. Which is wanted if this becomes a preset?
- **Q3:** On a phone the map starts zoomed on O2, not on the whole picture, because the whole picture would be unreadable. Is that the right start?
- **Q4:** Some limits (L6, L9) relate to the whole check, not to one box. The map attaches them to the live check. Is that reading right?

Variation slug: zoomable relation map
