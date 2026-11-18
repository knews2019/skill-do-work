# UR-185 Round 0

The frozen gallery is [baseline.html](../baseline.html). It shares the live gallery’s relative paths so its previews and report links work directly from the checkout. The capture manifest records the ten report identities, source hashes, and screenshot hashes. The gate verifies them against the rendered cards across four views.

Run from the repository root with `playwright-cli`, Google Chrome and Node Playwright installed. On this machine the Node package is in `/opt/homebrew/lib/node_modules/playwright`.

```sh
python3 docs/report-gallery/capture-previews.py
GALLERY_PLAYWRIGHT=/opt/homebrew/lib/node_modules/playwright node _dev/tests/report-gallery-gate.cjs \
  --gallery docs/report-gallery/baseline.html \
  --manifest docs/report-gallery/preview-manifest.json \
  --out "$(mktemp -d /tmp/req750-repro.XXXXXX)"
```

The committed views and capture hashes describe the frozen browser run. Verify their bytes before regenerating captures.

```sh
shasum -a 256 -c docs/report-gallery/baseline/SHA256SUMS
```

The original-gallery receipt preserves its four screenshots and red gate results. The separate hash fixture receipt shows a loadable image with the wrong capture hash failing G2. The green receipt covers each gate in each view. Reader needs and acceptable report identities were supplied by the coordinator before evaluation. Keep the answer key out of evaluator briefs.
