Confirmed staging output: [index-v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines.html](index-v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines.html)

### Clarity Changes
- **Redrew options O1–O4 as four parallel six-phase timelines (`#rev-1`):** Mapped one full review session from the end of `ai-report second-opinion` (`T1`) through orientation (`T2`), full-depth reading (`T3`), feedback capture (`T4`), drift check and extraction (`T5`), and the start of the next `ai-report revise` (`T6`), marking every step as `chat`, `real page`, or `canvas` with a `MOCKUP — proposal` SVG diagram and an interactive surface filter.
- **Tallied session surfaces, extra steps, and option pricing (`#rev-2`):** Placed O1–O4 side by side to show that O2 and O3 use all three surfaces (`chat`, `canvas`, and `real page`) because captures hide 38–46% of text in closed `<details>` (`K1`) and Tinka rejects non-loopback hosts (`K2`), while pricing O2's ~20 manual steps (~1 hr) against O3 and deferred O4.
- **Anchored pilot stop rules (`R1–R3`), kill attempts (`K1–K4`), and live captures to session phases (`#rev-3`):** Linked each pilot gate (`R1` in `T4`, `R2` in `T2/T4`, `R3` in `T1/T5`), each kill attempt, and all four `screenshots/` PNGs directly to the timeline phases they govern.
- **Consolidated verified corrections, open questions (`Q1–Q7`), folded evidence, and capture lines (`#rev-4`):** Placed verified repository corrections and `Q1–Q7` directly above four expandable `<details>` panels holding `S1–S9`, the live-check command log and limits `L1–L9`, the printed O1–O4 capture lines, and the glossary and source map.

### Corrections to the Original Report
- **`screenshots/` image count (`index.html:531`):** The original footer states *"The five images in `screenshots/` are real captures"*, whereas `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/screenshots/` holds **4** PNG files (`canvas-lineage-overview.png`, `tile-reply-markdown-full-view.png`, `tile-capture-full-view.png`, `s1-probe-html-tile-full-view.png`), matching the 4 figures in `#proof`.
- **Catalog snapshot vs. bundles on disk (`index.html:67, 424`):** `ai-reports/catalog.json` (generated `2026-11-04`) lists **46** bundles, while **49** bundle directories sit under `ai-reports/` on disk because three newer bundles (`2026-11-11_1411_...`, `2026-11-14_1452_...`, `2026-11-14_1547_...`) post-date the last catalog run; ordinary creation and `second-opinion` do not update `catalog.json` (`skills/do-work-toolbox/actions/ai-report.md:42-48, 85`).
- **Repository `HEAD` progression (`index.html:38, 401, 531`):** `VERSION` remains `0.315.5`, while `HEAD` has moved past `9de5839` through `7a7d5cea`, `f438ac98`, merge `d2100647`, and sibling `v1c01` commit `6bd7405d` (none modifying `skills/`).
- **Local explanation for `Q6` (`index.html:72, 462`):** `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:873-960` publishes sibling pages before running the render judge and inserts `<!-- second-opinion-tree -->` into `index.html` only on judge `pass`; commit `dce5dd1f` (`2026-11-04`) recorded only the two `v1c01` files in `2026-11-04_1524_proposal-tstack-review-workflows`.
- **Verification boundary:** Local repository claims (`report_index.go`, `report_judge.go`, `report_organize.go`, `report_second_opinion.go`, `folder_structure.go`, `html_preview.go`, `ai-report.md`, `ai-report-guide.md`) were verified directly; external Tinka `77ba29c` claims and container measurements rely on the original report because external services and network fetches are outside this pass's boundary.

### Open Questions
- **Q1:** Keep Tinka files outside the repository (`~/.local/share/do-work-canvas/<repo>/`, advised) or under `.tinka/` with `.git/info/exclude`?
- **Q2:** Should O3 (`ai-report-canvas`) demand a pinned Tinka checkout (`S5`) or accept any installed Tinka?
- **Q3:** Confirm or adjust R1's 50% canvas-before-chat threshold before starting the 2-lineage pilot.
- **Q4:** Keep only wide-light captures on the canvas (1 tile/page) or include dark/phone captures (4 tiles/page)?
- **Q5:** Should `ai-report-judge` add single-page targeting for second-opinion pages so O2 no longer needs 3 scratch copies?
- **Q6:** Back-fill or leave untouched the unlinked `v1c01` page in `2026-11-04_1524_proposal-tstack-review-workflows`?
- **Q7:** Would adding second-opinion pages under each bundle in `ai-reports/index.html` solve enough of the lineage map without Tinka?

Variation slug: parallel review session timelines
