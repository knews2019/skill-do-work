# REQ-750 browser gate

Run from the repository root with Node.js, the Playwright Node module and
Google Chrome. Point `GALLERY_PLAYWRIGHT` at the module directory if it
is installed outside Node's normal module search path. Use `--channel chromium`
for a Playwright-managed Chromium installation.

```sh
node _dev/tests/report-gallery-gate.cjs \
  --gallery docs/report-gallery/index.html \
  --manifest /path/to/capture-manifest.json \
  --out /path/to/evidence
```

The runner loads the gallery from a file URL. It writes four full-page PNGs
(`1440-light.png`, `1440-dark.png`, `390-light.png`, `390-dark.png`) and
`gates.json` with failure arrays per gate and view. Screenshots use a
900 px viewport height and include the entire gallery. Exit 0 means all
seven gates passed in every view, 1 means a gate failed, and 2 means a
setup or runner error. Use a fresh output directory for each run.

## Capture manifest from REQ-749

Supply JSON with a fixed `viewport: { width, height }`, `scrollY: 0`, and ten
`captures` rows. Each row has four string fields:

| Field | Value |
| --- | --- |
| `report` | Repository-relative `ai-reports/.../index.html` path |
| `image` | Repository-relative path to the captured preview image file, or an absolute path for fixtures |
| `reportSha256` | Lowercase SHA-256 of the source HTML bytes at capture |
| `imageSha256` | Lowercase SHA-256 of the emitted screenshot bytes |

Generate the manifest in the capture script as it visits each source at
scroll position zero and saves the fixed viewport screenshot. The gate
trusts this capture record, then matches the preview destination, the
browser-selected image URL, decoded image dimensions, and both file
hashes to it. Keep the trusted record separate from candidate markup.
Identical image hashes across reports are allowed. If REQ-749 emits a
different schema, adapt its output to this schema using its trusted capture
records. A run without `--manifest` exercises the browser checks and fails G2.

## DOM hooks

- Put `data-gallery-card` on each of the ten worked cards. Give each card
  or its fragment target the bare name as an HTML `id`.
- Put `data-gallery-name` on the clickable name control. Its visible text
  is the copied name.
- Put `data-gallery-preview` on the preview anchor. Its `href` is the
  source report and it contains the rendered `img`. The click may open
  the current tab or a new tab.
- For another markup structure, pass `--selectors` with a JSON object
  overriding `card`, `name`, `preview`, `image` or `controls` with Playwright
  selector strings. Name and preview are scoped to each card, and image
  is scoped to the preview. Include custom controls in the controls
  selector. A `report` selector scoped to the card can supply the expected
  navigation URL from a separate anchor when testing legacy markup.

The gate checks the clipboard after each click, loads the bare fragment
from a fresh navigation, and allows document-bottom scroll clamping
when the target is still visible. G6 checks linked files, HTTP responses,
named/ID anchors and image decoding. G7 opens native disclosures to test
their visible controls and checks layout width, control rectangles, sampled hit targets and
Playwright actionability. The saved screenshots show the default state.

## Original gallery red run

```sh
GALLERY_PLAYWRIGHT=/opt/homebrew/lib/node_modules/playwright \
  node _dev/tests/report-gallery-gate.cjs \
  --selectors _dev/tests/fixtures/report-gallery-legacy-selectors.json \
  --out /tmp/req750-final-red
```
