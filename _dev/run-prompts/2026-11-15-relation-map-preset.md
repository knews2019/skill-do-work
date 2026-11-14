# Capture prompt: ship a relation-map preset for ai-report second-opinion

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.

## Gate

Capture this only after the maintainer has looked at the example page and wants it as a standard option. The example is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`. If the capture gives no sign of that, ask once and park the request.

## The request

Add one shipped second-opinion preset, `relation-map`, so `ai-report second-opinion --preset relation-map <bundle>` writes a page that shows a report as one zoomable map of its items and the relations between them.

### Context

- The maintainer wants report content that is easy to understand by zooming in and out and by seeing relations. The example page was built by hand with `--editor-direction` and passed the render judge.
- A preset is one Markdown file with a direction to the editor. The loader reads `<bundle>/second-opinion-presets/<name>.md`, then the shipped folder `skills/do-work-toolbox/docs/second-opinion-presets/`. No loader change is needed. `mindmap.md` there is the nearest existing preset: read it first and say in the new preset how the two differ (a mind map is a tree from one center; a relation map shows relations between any two items).

### What the preset must ask the editor for

- Boxes are the report's own items (for example options, risks, findings, limits, files, requests). Lines are relations the report states. The editor says on the page that the lines are its reading of the report.
- At most four kinds of line, each with one color and a plain name, and a key that explains them.
- Semantic zoom: zoomed out, a box shows its code and title; zoomed in, a summary; a click opens the detail with its sources and its relations, and fades unrelated boxes.
- Drag to move, pinch or Control-scroll to zoom, buttons and keys for the same, and a button that shows everything. Plain scrolling must still scroll the page.
- The page works without JavaScript: the map shows at full size and scrolls. Every box is also given as text, and every relation as one table row.
- No library. Inline SVG and plain JavaScript only, as the second-opinion prompt already requires.
- About 30 boxes is the upper limit for one map. With more, group items and say so.
- On a narrow screen the map starts on the main item, not on the whole picture.
- Every fact still comes from the report or the repository, as in `bold.md`.

### Files to change

- `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (new).
- `skills/do-work-toolbox/docs/ai-report-guide.md`: the preset list in "A Second Opinion".
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`: `TestSecondOpinionShippedPresetsLoadFromTheToolbox`, if it names the shipped presets.
- `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.

### What must not change

- The preset loader, the second-opinion prompt, the other presets, and every existing `ai-reports/` bundle.
- The report kit (`skeleton.html`, `presentation.css`, `presentation.js`).

### Verification

RED first: `--prepare --preset relation-map` on a fixture bundle is refused because the preset does not ship (the existing shipped-presets test is the place for it).

GREEN: the same call returns state `prepared` with `preset_file` in the shipped folder; the guide names the preset; one real page is built with the preset on a bundle other than the example, passes `ai-report-judge`, and is shown to the maintainer before the request is called done.
