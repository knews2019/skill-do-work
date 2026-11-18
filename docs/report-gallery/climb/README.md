# Gallery climb receipts

The round directories hold candidate patches, original gate screenshots, capture and gate command output, Opus transcripts and run identities, and sealed choices. The decision files map each choice to an immutable source report path, record both versions' five scores in rubric order, and hash the neutral source and screenshots released to Opus.

Each candidate ran the frozen capture command inside a separate worktree. The temporary recapture outputs were restored to the frozen manifest and image bytes before the gate ran against the exact gallery and image set proposed for use on main. The screenshot inputs for the first two judges used the committed Round 0 gate views for the incumbent.

```sh
python3 docs/report-gallery/capture-previews.py
GALLERY_PLAYWRIGHT=/opt/homebrew/lib/node_modules/playwright node _dev/tests/report-gallery-gate.cjs \
  --gallery docs/report-gallery/index.html \
  --manifest docs/report-gallery/preview-manifest.json \
  --out "$round_views"
```

For each green round the operator ran the prepared `blind_round.py stage` with the two HTMLs and four views per version, in 1440-light, 1440-dark, 390-light and 390-dark order, and retained its randomized mapping privately. The reviewer called `blind_round.py release` after writing all eight first choices and before reading either HTML. The score decision uses only the paired results from that reviewer. No reader-selection time was measured.
