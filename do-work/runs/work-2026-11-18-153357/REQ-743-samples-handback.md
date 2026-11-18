# REQ-743 sample handback

Branch `t3/build-report-gallery-samples`  
Commit `813c0517daf17670072dd9a199e026db26c420a2`

## Files

```text
skills/do-work-toolbox/docs/report-gallery/parts-over-time/example.html
skills/do-work-toolbox/docs/report-gallery/parts-over-time/sample.md
skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/example.html
skills/do-work-toolbox/docs/report-gallery/rhythm-over-cycle/sample.md
skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/example.html
skills/do-work-toolbox/docs/report-gallery/two-measures-side-by-side/sample.md
skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/example.html
skills/do-work-toolbox/docs/report-gallery/verdict-with-checks/sample.md
```

The worked examples use original, visibly labelled sample data. The membership charts share a month selector and label both endpoints. The dispatch verdict leads with the result, then six named checks and one evidence disclosure. The repair charts show stacked totals, same-scale parts, and hatching through 14 May. The library charts switch between stacked bands and lines, with six same-scale rows, labelled peaks, and one time selector. All four pages load the existing kit through relative links and keep custom styles inline.

## Verification

- Loaded the four pages as local files in a temporary copy of the installed toolbox layout. Headless Chrome and axe-core passed all 48 page, width, and theme combinations at 320, 390, 661, 768, 1280, and 1440 pixels.
- Exercised the selectors, chart alignment and markers, view switch, deck mode, and evidence disclosures. All four disclosures also passed with JavaScript disabled. Reviewed rendered phone and desktop captures.
- Checked the agreed YAML record schema, strict HTML parsing, local kit links, the shipped-package reference contract, and `git diff --check`. The commit contains the eight listed files.

## Open issues

No open sample issues. Generated gallery and install/update verification remain part of REQ-743 integration.
