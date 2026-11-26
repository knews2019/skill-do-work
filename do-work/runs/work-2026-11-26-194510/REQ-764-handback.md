# REQ-764 · Product Launch Brief

The shipped **Product Launch Brief** sample is ready for integration on branch `t3/product-launch-brief`, commit `b61df115c5eb714de0d35aeba33451f9ba5971b7`. The implementation checkout is clean.

## Changed files

| Repository path | Change |
| --- | --- |
| `docs/report-gallery/index.html` | Regenerated with one new Sample designs entry. |
| `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/example.html` | Offline City Postcards worked example with a two-sided collectible, a collection example, an ordered newsroom workflow, and three linked launch decisions. |
| `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/sample.md` | Shipped record with the gallery question and the three concrete Borrow moves. |
| `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/first-ferry.webp` | Original local illustration used for the card, collection and workflow. |
| `skills/do-work-toolbox/docs/report-gallery/product-launch-brief/preview.png` | Authentic 1040 × 740 Chrome capture of the same report with JavaScript disabled, showing both faces in full. |

## Checks and observations

- Before implementation the named entry was absent from the static gallery and shipped catalog. After regeneration it appeared once under **Sample designs**. The CLI listed it and resolved both `Product Launch Brief` and `product-launch-brief` to the same sample record.
- Offline Chrome opened the actual report at 1440 and 390 pixels in light and dark. The desktop premise sits beside the card and the phone stacks the premise, trial, collection, workflow and decisions. The same ferry image and title appear in the three sections. The question about who covers the audience producer’s phone check jumps to **Check and schedule**. All three decision links target their corresponding steps.
- Desktop Tab showed a visible 3-pixel focus outline. Enter revealed the excerpt and Space returned to the front without moving focus. Emulated touch flipped the card both ways on the phone. With reduced motion the face changed immediately with zero active animations. Chrome’s accessibility tree exposed only the active face in each state.
- The preview image showed the illustrated front and excerpt together in the gallery’s thumbnail and modal. At both widths the generated gallery opened the actual report and copied the design name and the exact `do-work-toolbox ai-report --design "Product Launch Brief"` starter. The disclosure showed the three Borrow moves.
- `go test ./internal/toolboxcommands -run 'Test.*(ReportGallery|ReportJudge)'` passed. The repository’s `report-gallery-gate.cjs` passed all seven gates in all four sizes and themes. The staged diff passed `git diff --cached --check`.
- The [offline browser script](/tmp/req764-checks.BciyrNMe/check-launch-brief.cjs), [desktop front](/tmp/req764-checks.BciyrNMe/report-1440-light-front.png), [phone back](/tmp/req764-checks.BciyrNMe/report-390-dark-back.png), [phone gallery preview](/tmp/req764-checks.BciyrNMe/gallery-preview-390.png) and [gallery gate results](/tmp/req764-checks.BciyrNMe/gallery-gate/gates.json) are temporary verification evidence.

## Decisions

- The City Postcards premise, ferry excerpt and collection are fictional. The card flip is interactive and the collection illustrates a saved item. The proposed weekday cadence is labeled as a pilot assumption.
- Applied the user’s design screenshot and the Claude Opus 5.5 review from the brief. The sample uses whitespace and rules between sections, concise action labels, one accent for the control and focus, and ordered steps with named roles. The Borrow notes tell future writers where to put the product trial, how to carry the same example into the workflow, and how to connect each decision to a step.
- The private design reference was inaccessible to the browsing tool. The public sample follows the approved brief and uses original copy and artwork. The local WebP was generated with the built-in image tool and resized for the report.

Image prompt used

> Original illustration for a collectible card in a fictional local newspaper feature, landscape 3:2. A small cream-and-rust commuter ferry with deep teal hull approaches a wooden city pier just after sunrise. An amber sun sits low above a misty waterfront, understated warehouse roofs and a small red harbor beacon are visible, with long teal reflections in the water. Vintage editorial travel art, screen-printed gouache, bold simplified silhouettes and subtle paper grain. Full-bleed artwork with the ferry prominent, legible at thumbnail size. Palette of rust orange, aged ivory, petrol teal and soft sage. No text, letters, symbols, logos, borders or watermarks.

## Outstanding concerns

Touch verification used Chrome mobile emulation. Reuse verification covered the shipped gallery starter and the CLI’s design resolution. Release and finalization belong to integration.

## Discovered Tasks

None.
