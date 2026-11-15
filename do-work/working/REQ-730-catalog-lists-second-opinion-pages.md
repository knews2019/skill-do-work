---
id: REQ-730
title: 'Report catalog lists the second-opinion pages of each bundle'
status: claimed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "ai-reports/catalog.json", "ai-reports/index.html", "CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-14T23:28:19Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-14T23:30:35Z
  basis:
    - Route B
    - 8-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
dispatch_at: 2026-11-14T23:35:11Z
builder_handback_at: 2026-11-15T01:07:30Z
---
# Report Catalog Lists the Second-Opinion Pages of Each Bundle
## What
Under each bundle's row, the report catalog (`ai-reports/index.html` and `ai-reports/catalog.json`) lists that bundle's second-opinion pages with tree id, preset, model and date.
## Why
The catalog has one entry per bundle, so second-opinion pages cannot be found through it. The bundle `ai-reports/2026-11-14_1547_proposal-tinka-review-surface` holds `index.html` and three second-opinion pages (`v1c01`, `v1c02`, `v1c03`), and the catalog shows only the first.
## Context
- One entry per bundle: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`, around lines 166 to 172, `pickReportEntry`.
- Each second-opinion page's head already carries its tree id, preset, model and date as `ai-report-second-opinion-*` meta tags.
## Detailed Requirements
- `catalog.json` gains a list of second-opinion pages per bundle. The change is additive.
- `ai-reports/index.html` shows the pages under their bundle's row.
- A page without the tree meta tags (a plain second-opinion page) is still listed, by its file name.
- Files: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` and its test file; the catalog forms in `skills/do-work-toolbox/actions/ai-report.md`; "Report Catalog" in `skills/do-work-toolbox/docs/ai-report-guide.md`; `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
## Constraints
- Every existing field of `catalog.json` keeps its name and meaning, so `revise`, `find`, `organize` and `second-opinion latest` keep working.
- No bundle is edited. The catalog reads bundles and never writes into them.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** A Go test builds a fixture bundle with `index.html` and one `index-v1c01-*.html` page and asserts that the catalog entry lists that page.
**Why RED now:** The entry has no such list.
**GREEN when:** The test passes; `ai-reports/index.html` shows the pages under their bundle; `go test ./...` passes in the do-work-cli module; the existing catalog tests pass unchanged.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- Certainty: firm on the outcome. The field name and the HTML shape are the builder's choice.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes do-work-cli Go code for reports.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and the main file is named, but the page detection, the JSON field shape and the HTML rows must follow the existing catalog and second-opinion code, which needs exploration. One Go file, its test, two docs and the derived catalog pair change; no architecture change.

**Planning:** Not required

## Open Questions

The REQ has no open questions from capture. Pre-dispatch made these choices (unattended run, best judgment):

- [~] Which files count as second-opinion pages? → **D-01**: Builder chose: every regular file at the top level of the bundle whose name matches the published second-opinion page shape `^index-[A-Za-z0-9][A-Za-z0-9._-]*\.html$` (the same shape `secondOpinionBanner` accepts), minus the bundle's own entry file. That covers tree pages (`index-v1c01-...html`) and plain pages (`index-claude-claude-opus-5-5.html`); staging pages start with a dot and never match. Reasoning: one existing shape, no new grammar. Value: plain pages are listed by file name as the REQ requires. Risk: a hand-made `index-foo.html` in a bundle is also listed; low and visible, reversible by narrowing the pattern.
- [~] Field name and shape in `catalog.json`? → **D-02**: Builder chose: a new per-bundle field `second_opinions`, always an array (never null), each item `{ "page", "tree_id", "preset", "model", "date" }`, where `page` is the repository-relative path like `entry`, and the other four are read from the `ai-report-second-opinion-id/-preset/-model/-date` metas (each is the raw meta value, an empty string when the meta is absent; the catalog never invents a value). Order: byte sort of file names, which is tree pre-order. Reasoning: additive, matches the existing snake_case fields. Value: `revise`, `find`, `organize` and `second-opinion latest` read none of it. Risk: none for readers; a future rename would be a breaking change, so pick it once.
- [~] Is the real `ai-reports/catalog.json` / `ai-reports/index.html` pair regenerated? → **D-03**: Builder chose: yes. After the Go change passes, the builder runs `ai-report-index` in its worktree with the new code and commits the regenerated pair in its own commit, so `ai-reports/index.html` shows the pages (the REQ's GREEN line). Reasoning: the pair is tracked, derived and already stale (it lacks the 2026-11-14 bundles). Value: GREEN evidence on real data. Risk: a merge conflict with a sibling that also changes the pair; the integrator resolves it by rerunning `ai-report-index` on the merged tree, never by hand-merging.
- [~] How does the HTML page show them? → **D-04**: Builder chose: under the bundle's row, one nested list inside the title cell (`<ul class="second-opinions">`), one item per page: link to the page, then `tree id · preset · model · date`, skipping empty parts, and the file name as the link text when there is no tree id. Reasoning: smallest change to the existing table; no new column. Value: rows without pages render byte-identical to today. Risk: cosmetic only.

- [~] Where does the page-name shape live? → **D-05**: Builder chose: one named package variable in `report_index.go` (for example `reportSecondOpinionPagePattern`), and `secondOpinionBanner` in `report_second_opinion.go` uses it in place of its inline regex. Reasoning: one definition of "a published second-opinion page name" instead of two copies. Value: the catalog and the publisher cannot drift. Risk: a one-line edit in a file REQ-728/729 do not touch; no behaviour change.

<!-- D-XX counter: last used D-05. Next decision: D-06. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Required-lessons consult.** Read `_dev/primes/lessons-releases.md` in full (856 tokens, the only captured entry; nothing in it applies directly: no file moves and no history links). `do-work/lessons-index.md` matches two more rows. `_dev/primes/lessons-action-files.md` and `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` stay dropped for budget (both `slugged: partial`, so neither can be narrowed; see the drop section above, refreshed unchanged). Pre-dispatch still read three do-work-cli families by grep because they touch this change: `catalog-tie-order` (catalog order is minute-precise; do not change the sort), `closed-enumeration-for-a-condition` (D-01 keys on the published page-name shape, not a list of known models or presets), and `assertion-passes-when-feature-is-dead` (the RED test must fail on base for the missing list, not pass on an empty one).

**Where the change goes** (all paths under `skills/do-work/tools/do-work-cli/internal/toolboxcommands/` unless named):
- `report_index.go:66-76` `reportBundle`: add the new `second_opinions` field (D-02) as the last field. Keep it a non-nil slice, the way `LinkedIDs` is initialised at `report_index.go:158`, so the JSON shows `[]`, never `null`.
- `report_index.go:156-203` `readReportBundle`: after the entry is picked (`report_index.go:173`), list the bundle's second-opinion pages. `pickReportEntry` at `report_index.go:205-236` already does one `os.ReadDir` of the bundle; a separate small helper that reads the folder again is simpler than changing the picker's return shape. Exclude the entry file itself (D-01), because a bundle with no `index.html` may pick an `index-*.html` as its entry.
- Page shape: reuse the regex in `report_second_opinion.go:1533` (`^index-[A-Za-z0-9][A-Za-z0-9._-]*\.html$`, inlined in `secondOpinionBanner`). Hoist it into one named package variable used by both call sites (D-05). Tree pages also match `secondOpinionTreePagePattern`; plain pages (real example: `ai-reports/2026-10_3-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html`) match only the wider shape.
- Meta reading: read each page's metas the way `secondOpinionTreeList` does at `report_second_opinion.go:1670-1676`: `readReportMetas(secondOpinionVisibleHTML(data))`, so a meta inside a script, style or comment never counts. Keys: `ai-report-second-opinion-id`, `-preset`, `-model`, `-date` (written by `secondOpinionDecorate`, `report_second_opinion.go:1626-1630`).
- `report_index.go:363-421` `renderReportIndexPage`: the row is written at about `report_index.go:410-413`. Add the nested list inside the title cell only when the bundle has pages (D-04), so rows without pages stay byte-identical. Link each page with an href relative to `ai-reports/index.html`, as `reportPageLink` (`report_index.go:354-361`) does for the entry; escape every value with `html.EscapeString`.
- Readers that must keep working: `reportIndexFind` (`report_index.go:423`), second-opinion `latest` (`report_second_opinion.go:375`, calls `reportIndexWalk`), organize (regenerates the catalog). None of them reads the new field, and no reader uses `DisallowUnknownFields` on `catalog.json`.
- Tests: `report_index_test.go` has the fixture at `report_index_test.go:19-52`, `runReportIndex` at 54 and `readReportCatalog` at 63. `TestAIReportIndexCatalogsEveryBundleNamingStyleOnce` (100) compares a summary struct, not raw JSON, so the new field does not break it. Add new tests as new functions at the end of the file, with their own small fixture (do not extend the shared fixture: that changes the "only two new files" check's inputs for every test).
- Docs: `skills/do-work-toolbox/actions/ai-report.md:40-48` (catalog forms) and `skills/do-work-toolbox/docs/ai-report-guide.md:94-102` ("Report Catalog": the bullet that lists the metas read from each bundle is the natural place for one new bullet).
- Real data: `ai-reports/catalog.json` and `ai-reports/index.html` are tracked, generator-marked and stale (46 bundles; neither 2026-11-14 proposal is listed). `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/` holds `index.html` plus `index-v1c01-…`, `index-v1c02-…`, `index-v1c03-…`, so it is the real GREEN check (D-03).

**Concerns.** Reading every `index-*.html` makes the walk read more files on every `index`, `find` and `second-opinion latest`; page counts are small (at most a few per bundle), so no cache. A second-opinion `-reply.md` file is not a page and is not listed.

*Generated by pre-dispatch exploration (main session, no agent)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) — the new per-bundle list in the catalog JSON and the nested list in the HTML page
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modify) — new RED-then-GREEN tests at the end of the file
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify, D-05) — `secondOpinionBanner` uses the shared page-name variable instead of its inline regex; no behaviour change
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — catalog forms: the catalog also lists each bundle's second-opinion pages
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — "Report Catalog": the new field and what it reads
- `ai-reports/catalog.json` (modify, regenerated by the new code, D-03)
- `ai-reports/index.html` (modify, regenerated by the new code, D-03)
- `CHANGELOG.md` (modify, by the integrator only)
- `VERSION` (modify, by the integrator only)

**Files I will NOT touch:** every file inside any `ai-reports/` bundle; `reportIndexFind` behaviour (find still matches bundles only); `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` and Step 8 of `ai-report.md` (REQ-731's catalog reminder); `report_judge.go` (REQ-728, REQ-729).

**Acceptance criteria (restated from REQ):**
- [ ] `catalog.json` gains a list of second-opinion pages per bundle, with tree id, preset, model and date; the change is additive.
- [ ] `ai-reports/index.html` shows the pages under their bundle's row.
- [ ] A page without the tree meta tags (a plain second-opinion page) is still listed, by its file name.
- [ ] Every existing field of `catalog.json` keeps its name and meaning, so `revise`, `find`, `organize` and `second-opinion latest` keep working; the existing catalog tests pass unchanged.
- [ ] No bundle is edited; the catalog reads bundles and never writes into them.
- [ ] The RED Go test (fixture bundle with `index.html` and one `index-v1c01-*.html` page) fails at base and passes after; `go test ./...` passes in the do-work-cli module.

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only, committed by the coordinator)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-730-preflight-probe.sh) and the repository gate green at df1bf423 (158 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*
