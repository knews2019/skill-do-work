# REQ-731 Hand-back (ai-report prints the catalog rebuild command when a new report is missing from the catalog)

Branch: `worktree-agent-REQ-731-report-catalog-rebuild-reminder`, tip `0e798128`, base `1de00505`.

## Commits
- `a6ac1be6` [REQ-731] RED: ai-report-index --check prints the rebuild line only when a bundle is missing
- `0e798128` [REQ-731] ai-report prints the catalog rebuild command when a new report is missing from the catalog

## Files
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (new test at end of file)
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (one branch in `handleAIReportIndex`, usage string, new `reportIndexCheck` at end of file)
- modified `skills/do-work-toolbox/actions/ai-report.md` (Step 8: one paragraph and one fenced command after the organize check)
- modified `skills/do-work-toolbox/docs/ai-report-guide.md` (Report Catalog: one new bullet)

## Scope changes
None. CHANGELOG.md and VERSION left to the integrator.

## Red-green evidence
- After commit 1: `go test -count=1 -run TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing ./internal/toolboxcommands/` → FAIL: `--check [ai-reports/2026-11-15_1100_b] outcome=failure findings=[{Code:TOOLBOX-USAGE ... Evidence:[Usage: ai-report-index [--find <topic...>]]`.
- GREEN probe at commit 1 (base plus RED test): exit 1, `FAIL github.com/.../toolboxcommands`.
- After commit 2: `go test -count=1 ./internal/toolboxcommands/` → `ok ... 30.364s`; GREEN probe → `REQ-731 GREEN probe: PASS` (exit 0).

## Verification
1. `gofmt -l ./internal/toolboxcommands` → empty. `go vet ./internal/toolboxcommands/` → clean. `go test -count=1 ./internal/toolboxcommands/` → ok.
2. `bash _dev/tests/contract-regressions.sh` → exit 1 twice, both ONLY per-file time budgets under load (18 live sessions, load average 39):
   - run 1: `FAIL: _dev/tests/contracts/core-checks.sh took 36s` (every probe passed).
   - run 2: core-checks passed in 8s; `FAIL: _dev/tests/contracts/replace-text-section.sh took 47s`. No other failure lines in either run. A different file each run, so load, not this change (neither file touches report code). Integrator: rerun on a quiet machine.
   - `bash _dev/tests/audit-lockins.sh` → `Audit lock-in regressions passed.` exit 0.
3. Manual scratch Git repo: bundle `a` + `ai-report-index` → `ai-reports/catalog.json: 1 bundles`; added `b`; `--check ai-reports/2026-11-15_1100_b` → `ai-reports/2026-11-15_1100_b is not in ai-reports/catalog.json yet: run do-work-toolbox ai-report index to rebuild the catalog.` exit 0; `--check ai-reports/2026-11-15_1000_a` → no output, exit 0.
4. GREEN probe → `REQ-731 GREEN probe: PASS`.

## P-A-U
- PLAN: new read-only `--check <bundle>` branch decoding only `catalog.json` (model: organize `--check`); one Step 8 paragraph + command; one guide bullet.
- APPLY: as planned; test first (RED commit), then feature.
- UNIFY: `git diff --stat HEAD~2` = 4 files, +106/-1; checked report_index.go (one branch, usage, function at end, no debug), test at end of file, ai-report.md (only lines after the organize block), guide (one bullet added, neighbours untouched).

## Decisions
- D-04 DECIDE & STATE: hand-made catalog listing the bundle prints nothing. The REQ file's D-02 text says a hand-made catalog counts as not listing; the brief (later) says decode the same and print nothing when listed. Followed the brief: the check only answers "is it listed".
- D-05 DECIDE & STATE: Step 8 says "Skip it for a revise", because the revise form reaches Step 8 for its commit rule but has already regenerated the catalog.
- D-06 DECIDE & STATE: no Verification Checklist clause added; the Step 8 paragraph is the single home, and a checklist line would restate it.
- D-07 DECIDE & STATE: `--check` with an empty argument or more than one argument returns the usage result (tested with `--check`, `--check ""`, `--check a b`).

## Discovered Tasks
- `skills/do-work-toolbox/actions/architecture-report.md:119-122` and `skills/do-work-toolbox/actions/stakeholder-report.md:51-54` also write bundles and run the organize check, but do not run the new catalog check (D-03). Note architecture-report caps output at "the one allowed line beyond the eight", so adding the catalog line there needs a wording decision.

## Integration seams
- `report_index.go`: `handleAIReportIndex` lines 78-103 (new `--check` branch at 96-101, usage string at 102); new `reportIndexCheck` at 475-495 (end of file). No change to `reportCatalog`/`reportBundle`; decoding uses default `encoding/json`, so REQ-730's new `second_opinions` field is ignored.
- `report_index_test.go`: new `TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing` at 317-385 (end of file); uses existing helpers only.
- `ai-report.md`: Step 8, new paragraph + fenced block at lines 201-206, after the organize `--check` block, before `## Output Format`.
- `ai-report-guide.md`: "Report Catalog", one new bullet at line 102 (between "hand-made file is refused" and "three optional meta tags").

## Proposed CHANGELOG entry
**New Reports Remind You to Rebuild the Catalog**
- What changes for an installed copy: after an ordinary `ai-report`, the summary prints one `ai-report index` line while the new report is missing from `ai-reports/catalog.json` (or the catalog does not exist yet). `ai-report-index --check <bundle>` is new and writes nothing. Creating a report still never rebuilds or commits the catalog.

## Proposed lesson
none

## Lessons read
- `_dev/primes/lessons-releases.md` (families `canonical-link-outlives-its-target`, `manifest-ownership-vs-edit-content`): no match, no file moves or manifest edits here.
- Dropped-for-budget bullets read: `[family: catalog-tie-order]` and `[family: commit-preflight-before-side-effect]` in `lessons-do-work-cli.md` (no match: the check never writes and does not depend on order), `[family: alternate-writer-contract-drift]` in `lessons-action-files.md` (applied: grepped for organize-check restatements; found the two writers listed under Discovered Tasks).
