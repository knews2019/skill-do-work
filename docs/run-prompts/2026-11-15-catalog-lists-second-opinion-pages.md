# Capture prompt: the report catalog shows second-opinion pages and new reports

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.

## The request

Two gaps in the report catalog (`ai-reports/index.html` and `ai-reports/catalog.json`), both seen on 2026-11-14.

1. **Second-opinion pages are not in the catalog.** It has one entry per bundle (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`, around lines 166 to 172, `pickReportEntry`). The bundle `ai-reports/2026-11-14_1547_proposal-tinka-review-surface` holds `index.html` and three second-opinion pages (`v1c01`, `v1c02`, `v1c03`), and the catalog cannot show them. Under each bundle's row, list its second-opinion pages with tree id, preset, model and date. Each page's head already carries these as `ai-report-second-opinion-*` meta tags.
2. **A new report is missing from the catalog until someone rebuilds it.** Ordinary creation writes only the bundle. On 2026-11-14 two new proposal bundles were absent from `catalog.json` (46 bundles, no match for either). First read why creation does not rebuild the catalog (`skills/do-work-toolbox/actions/completed-work-presentation-reference.md`, Collision-Safe Publication). If that reason still holds, do the smaller change: Step 8 prints the `ai-report index` line when the new bundle is not in the catalog. If it does not hold, creation rebuilds the catalog as the revise form does. Park the choice as a pending answer when the reason is unclear.

### Files to change

- `report_index.go` and its test file in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`.
- `skills/do-work-toolbox/actions/ai-report.md` (catalog forms, Step 8) and `skills/do-work-toolbox/docs/ai-report-guide.md` ("Report Catalog").
- `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.

### What must not change

- `catalog.json` changes are additive: every existing field keeps its name and meaning, so `revise`, `find`, `organize` and `second-opinion latest` keep working.
- No bundle is edited. The catalog reads bundles and never writes into them.
- A page without the tree meta tags (a plain second-opinion page) is still listed, by its file name.

### Verification

RED first: a Go test builds a fixture bundle with `index.html` and one `index-v1c01-*.html` page and asserts that the catalog entry lists that page; it fails because the entry has no such list.

GREEN: the test passes; `ai-reports/index.html` shows the pages under their bundle; a second test covers the chosen fix for gap 2; `go test ./...` passes in the do-work-cli module; existing catalog tests pass unchanged.
