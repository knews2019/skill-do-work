# REQ-730 Hand-back (the report catalog lists the second-opinion pages of each bundle)

Branch `worktree-agent-REQ-730-catalog-lists-second-opinion-pages`, tip `5a19ca8b`.

## Commits
- `f00c54d5` [REQ-730] RED: catalog tests for second-opinion pages under their bundle
- `359bb427` [REQ-730] catalog lists each bundle's second-opinion pages under its row
- `5a19ca8b` [REQ-730] regenerate the report catalog pair with second-opinion pages

## Files
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go`
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (one line, D-05)
- modified `skills/do-work-toolbox/actions/ai-report.md` (one sentence in "Catalog forms: index and find")
- modified `skills/do-work-toolbox/docs/ai-report-guide.md` (one new bullet in "Report Catalog")
- modified `ai-reports/catalog.json`, `ai-reports/index.html` (regenerated, D-03)

## Scope changes
None. Field names, href form and test names match the probe's assumptions.

## Red-green evidence
- After commit 1 (`go test -run TestAIReportIndexLists`):
  - `--- FAIL: TestAIReportIndexListsSecondOpinionPagesUnderTheirBundle` — `report_index_test.go:362: bundle without pages: second_opinions="", want []`
  - `--- FAIL: TestAIReportIndexListsPlainSecondOpinionPageByFileName` — `report_index_test.go:428: fallback entry listed as a page or field missing` (field absent, so nil)
- After commit 2: both `--- PASS`, with all five existing `TestAIReportIndex*`, both `TestAIReportOrganize*` and `TestSecondOpinionBannerOnlyRealBodyAndPerFileIdempotency` passing.
- GREEN probe at tip `5a19ca8b`: `REQ-730 probe PASS`. At base it fails by design (the two test names do not exist, so its step 1 fails; the tracked catalog has no `second_opinions`). I did not check out base to rerun it; the RED run above shows the same missing field.

## Verification
- `gofmt -l ./internal/toolboxcommands`: empty. `go vet ./internal/toolboxcommands/`: clean.
- `go test -count=1 ./...` (do-work-cli module): all packages ok, 1:52 wall time; a second run (inside the background command below) printed no non-ok line.
- `bash _dev/tests/contract-regressions.sh`: "Contract regression checks passed."
- `bash _dev/tests/audit-lockins.sh`: "Audit lock-in regressions passed."
- `ai-report-index` in the worktree: exit 0, "ai-reports/catalog.json: 49 bundles", no organize hint, no refusal. 5 bundles now carry a nested list.
- `git diff --name-only $(git merge-base HEAD main) -- ai-reports/`: only `ai-reports/catalog.json` and `ai-reports/index.html`.
- Tinka row in the regenerated `ai-reports/index.html`:
  `<tr><td class="report-date">2026-11-14T15:47</td><td><a href="2026-11-14_1547_proposal-tinka-review-surface/index.html">Tinka as a review canvas for ai-report bundles</a><ul class="second-opinions"><li><a href="2026-11-14_1547_proposal-tinka-review-surface/index-v1c01-claude-claude-fable-5-2-custom-worked-example-tally.html">v1c01</a> · custom · claude claude-fable-5-2 · 2026-11-14</li><li><a href="…/index-v1c02-agy-gemini-pro-orchestrator-custom-parallel-review-session-timelines.html">v1c02</a> · custom · agy gemini-pro-orchestrator · 2026-11-14</li><li><a href="…/index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html">v1c03</a> · custom · claude claude-fable-5-2 · 2026-11-14</li></ul></td><td class="report-ids"></td></tr>`
  (`…` = `2026-11-14_1547_proposal-tinka-review-surface`, shortened here only.)

## P-A-U
- PLAN: pattern variable + `reportSecondOpinionPage` type + a re-reading helper; field initialised to `[]`; nested list appended after the successor span in the title cell; D-05 one-line swap.
- APPLY: done as planned; one CSS line `ul.second-opinions{…}` added to the page style (this changes the page header bytes for every catalog, rows without pages are unchanged).
- UNIFY: `git diff --stat` base..tip: 7 files, 409+/54-; checked the Go diff (no debug output, gofmt clean), the two doc lines, and that the `ai-reports/` diff is the catalog pair only.

## Decisions
- D-06 DECIDE & STATE: a page read error leaves the page listed with empty metas (not dropped, not fatal), matching how an unreadable entry still yields a bundle.
- D-07 DECIDE & STATE: the nested list sits after the "superseded by" span inside the title cell, so a superseded row keeps its existing text first.
- D-08 DECIDE & STATE: the HTML list shows the raw preset value, empty skipped (D-04); it does not substitute "plain" the way `secondOpinionTreeList` does, so the catalog never invents a value (D-02).

## Discovered Tasks
- None.

## Integration seams
- `report_index.go`: new var `reportSecondOpinionPagePattern` (:51-53); `reportBundle` gained the last field `SecondOpinions` plus new type `reportSecondOpinionPage` (:70-92); `readReportBundle` (:171-, init at the `bundle :=` line and one call before `reportKind`); `renderReportIndexPage` (:379-, one CSS line and the row `Fprintf` gained a `%s` + argument); new functions at the end `readReportSecondOpinionPages` (:488) and `renderReportSecondOpinionList` (:516). `handleAIReportIndex` and the usage string are untouched (REQ-731).
- `report_index_test.go`: new functions at the end only: `writeReportIndexFiles`, `TestAIReportIndexListsSecondOpinionPagesUnderTheirBundle`, `TestAIReportIndexListsPlainSecondOpinionPageByFileName`.
- `report_second_opinion.go`: `secondOpinionBanner` first line (:1533). The `judgeRender` call (REQ-728) is untouched.
- `ai-report.md`: "Catalog forms: index and find", one sentence after the bash block. Step 8 and the closing summary untouched.
- `ai-report-guide.md`: "Report Catalog", one bullet after the meta-tags bullet.
- `ai-reports/catalog.json`, `ai-reports/index.html`: on any conflict rerun `ai-report-index` on the merged tree.

## Proposed CHANGELOG entry
**Report Catalog Lists Each Bundle's Second-Opinion Pages**

What changes for an installed copy: `ai-report index` now writes a `second_opinions` list per bundle into `ai-reports/catalog.json` (additive; every existing field is unchanged) and shows those pages under their bundle's row in `ai-reports/index.html`; the next `ai-report index` run in a project picks them up, and no bundle is changed.

## Proposed lesson
- none

## Lessons read
- `_dev/primes/lessons-releases.md` (full); `lessons-do-work-cli.md` bullets `[family: catalog-tie-order]`, `[family: closed-enumeration-for-a-condition]`, `[family: assertion-passes-when-feature-is-dead]` by fixed-string grep. Applied: the RED test checks raw JSON bytes for `[]` so an absent field cannot pass; pages are keyed on the name shape, not a list of models.
