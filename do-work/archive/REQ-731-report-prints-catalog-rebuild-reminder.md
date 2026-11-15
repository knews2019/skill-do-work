---
id: REQ-731
title: 'ai-report prints the catalog rebuild command when a new report is missing from the catalog'
status: completed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-730, REQ-732]
batch: tinka-follow-ups
write_set: [skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go, skills/do-work-toolbox/actions/ai-report.md, skills/do-work-toolbox/docs/ai-report-guide.md, CHANGELOG.md, VERSION]
claimed_at: 2026-11-14T23:28:19Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
    - Route B
    - 6-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
  calculated_at: 2026-11-14T23:30:40Z
dispatch_at: 2026-11-14T23:35:12Z
builder_handback_at: 2026-11-15T01:23:03Z
integration_at: 2026-11-15T01:24:20Z
review_at: 2026-11-15T01:46:56Z
commit: 15ad4b3e3dae6992bf369b2d8181f9b94711f94f
heavy_verified_at: 2026-11-15T01:51:10Z
heavy_verified_revision: 15ad4b3e3dae6992bf369b2d8181f9b94711f94f
kb_status: pending
completed_at: 2026-11-15T01:51:37Z
release_at: 2026-11-15T01:51:37Z
---
# ai-report Prints the Catalog Rebuild Command When a New Report Is Missing From the Catalog
## What
After `ai-report` creates a report, its closing summary prints the one command that rebuilds the catalog, when the new bundle is not in `ai-reports/catalog.json`.
## Why
Ordinary creation writes only the bundle, so a new report is missing from the catalog until someone rebuilds it. On 2026-11-14 two new proposal bundles were absent from `catalog.json` (46 bundles, no match for either), and they still are.
## Context
- The revise form rebuilds the catalog and commits it with the bundle. Ordinary creation does not.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md`, Collision-Safe Publication, was read at capture. It forbids overwriting existing artifacts and gives no reason against a reminder.
## Detailed Requirements
- Step 8 of `skills/do-work-toolbox/actions/ai-report.md` checks whether the new bundle is in the catalog, and prints the `ai-report index` line when it is not.
- The check is a program step, not an agent's guess. Use an existing command when one can answer it, in the way Step 8 already runs `ai-report-organize --check`.
- The line is printed just before the final bundle-path line, like the organize line.
- The guide's "Report Catalog" section says that a new report enters the catalog at the next rebuild and that the summary reminds you.
## Constraints
- Creating a report does not rebuild or commit the catalog. That was the maintainer's choice at capture.
- No bundle is edited.
- When the catalog files do not exist yet, the same line is printed.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** Create a report in a project whose catalog exists, and read the closing summary.
**Why RED now:** The summary does not mention the catalog, and the new bundle is absent from `catalog.json`.
**GREEN when:** The summary prints the rebuild command when the new bundle is absent from `catalog.json`, and prints nothing about the catalog when the bundle is in it.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Open Questions
- [x] How should a new report get into the catalog: by a printed reminder, or by a rebuild on every new report? → Print a reminder (Recommended).
  Note (2026-11-14): The maintainer chose the reminder: no extra file changes, so two sessions that create reports at the same time cannot collide. The accepted risk: the catalog stays stale until someone runs the command. A rebuild on every new report is out of scope.
- [~] Which program answers "is the new bundle in catalog.json"? → **D-01**: Builder chose: a new read-only flag `ai-report-index --check <bundle>` in Go, with Go tests. It reads only `ai-reports/catalog.json`, prints one rebuild line when the bundle path is absent, the catalog file is missing, or it cannot be parsed, and prints nothing when the path is listed. It never walks or writes. Reasoning: no existing command reads `catalog.json` (`--find` walks the folders in memory and never reads the catalog, `report_index.go:422`), and the REQ's Builder Guidance names a new flag on an existing command as the allowed route; it mirrors `ai-report-organize --check` (`report_organize.go:97`). Value: one deterministic line, testable, no jq dependency. Risk: one more flag on the CLI surface and a Go change that shares `report_index.go` with REQ-730; easy to revert.
- [~] What does the check print, and for which path form? → **D-02**: Builder chose: the argument is the bundle path as the catalog stores it (`ai-reports/<folder>`, a trailing `/` trimmed); the printed line is one line naming the bundle and `do-work-toolbox ai-report index`, in the style of the organize hint (`report_organize.go:35`). A hand-made catalog (no generator marker) counts as not listing the bundle, so the line still prints and the rebuild then reports its own refusal. *(Superseded by D-04: a hand-made catalog is decoded the same way, so a bundle it lists prints nothing; one it does not list still prints the line.)* Reasoning: new bundles are always written directly under `ai-reports/`, so no month-folder form is needed. Value: one exact string the test pins. Risk: low; wording only.
- [~] Do `architecture-report` and `stakeholder-report`, which also write bundles and run `ai-report-organize --check`, get the same line? → **D-03**: Builder chose: no, out of scope; the REQ names Step 8 of `ai-report.md` only. The builder lists them as a discovered task. Value: keeps the change small. Risk: those reports stay missing from the catalog until a rebuild, as today.
<!-- D-XX counter: last used D-03. Next decision: D-04. -->
## Builder Guidance
- Certainty: the choice is the maintainer's. How the check is done is the builder's.
- If the check needs a new flag on an existing command, that is Go code with a Go test; say so in the request and add the files to the scope.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason (claim time): the request adds a flag to the do-work-cli `ai-report-index` command (families `catalog-tie-order`, `commit-preflight-before-side-effect`).
---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (one reminder line in Step 8), but no existing command can answer whether a bundle is in `catalog.json`, so the check needs a new read-only `--check` flag on `ai-report-index` (Go + Go test) plus prose in the action and the guide. Two subsystems, about six files, known patterns (`ai-report-organize --check`).

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:78` `handleAIReportIndex` routes no args (rebuild + organize hint) and `--find <topic>`; usage string at line 96. A new `--check <bundle>` branch goes here, with its function added at the end of the file. `reportCatalogFile` (line 27) and the `reportCatalog` / `reportBundle` types (lines 59-76, JSON field `path`) are what the check decodes. `reportIndexFind` (line 422) walks folders and never reads `catalog.json`, so no existing command can answer "is this bundle in the catalog" (D-01).
- `report_organize.go:97-104` is the pattern to copy: `--check` returns `exactOutputResult(line+"\n", nil)` or an empty string, writes nothing. Hint wording style: `report_organize.go:35`.
- `report_index_test.go` helpers: `reportIndexFixture` (line 19), `runReportIndex(t, root, args...)` (line 54), `snapshotReportTree` (line 76) for the "writes nothing" assertion. `report_organize_test.go:297` (`TestReportOrganizeHintPrintsOnlyWhenWorkRemains`) is the model test for a `--check` line.
- `skills/do-work-toolbox/actions/ai-report.md:188-200` Step 8: the summary paragraph, then the organize check whose line prints just before the final bundle-path line. The new check goes beside it with the same placement rule. Catalog forms section: lines 40-48. Verification Checklist: lines 214-223.
- `skills/do-work-toolbox/docs/ai-report-guide.md:94-102` "Report Catalog": add one sentence that a new report enters the catalog at the next rebuild and the summary reminds you.
- Sibling seam: REQ-730 changes `report_index.go` (`pickReportEntry`, `reportBundle`, catalog rendering) and appends tests to `report_index_test.go` in the same wave. Decoding must stay tolerant of extra JSON fields (default `encoding/json` behavior), and edits stay in `handleAIReportIndex` plus new functions at file end.
- Out of scope (D-03): `skills/do-work-toolbox/actions/architecture-report.md:122` and `stakeholder-report.md:54` also run the organize check after writing a bundle.
- Lessons consult: `_dev/primes/lessons-releases.md` read (no matching family: no file moves, no archive links). `lessons-do-work-cli.md` dropped for budget; the builder reads its `catalog-tie-order` and `commit-preflight-before-side-effect` bullets.

*Generated by pre-dispatch agent (explored directly, no subagent)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) — new read-only check flag on ai-report-index, usage string updated
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modify) — new test function TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing at the end of the file
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — Step 8 runs the catalog check and prints its line before the final path line; catalog forms and checklist mention it
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — Report Catalog says a new report enters at the next rebuild and the summary reminds you
- `CHANGELOG.md` (modify, integrator only) — release entry
- `VERSION` (modify, integrator only) — version bump

**Files I will NOT touch:** `report_organize.go`, `architecture-report.md`, `stakeholder-report.md`, any `ai-reports/` bundle, `ai-reports/catalog.json`, `ai-reports/index.html`, `completed-work-presentation-reference.md`.

**Acceptance criteria (restated from REQ):**
- [ ] Step 8 of `ai-report.md` runs a program check whether the new bundle is in `ai-reports/catalog.json` and prints the `ai-report index` line when it is not.
- [ ] The check is an existing command with a new read-only flag (D-01), covered by a Go test.
- [ ] The line prints just before the final bundle-path line, like the organize line.
- [ ] When the catalog files do not exist yet, the same line prints.
- [ ] When the bundle is in the catalog, nothing about the catalog prints.
- [ ] Creating a report still does not rebuild or commit the catalog, and no bundle is edited; the guide's "Report Catalog" section says a new report enters at the next rebuild and the summary reminds you.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) New read-only `--check <bundle>` branch decoding only `catalog.json` (model: organize `--check`); one Step 8 paragraph + command; one guide bullet.
- [x] **[APPLY]:** (from the builder hand-back) As planned; test first (RED commit `a6ac1be6`), then feature (`0e798128`).
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat HEAD~2` = 4 files, +106/-1; checked report_index.go (one branch, usage, function at end, no debug), test at end of file, ai-report.md (only lines after the organize block), guide (one bullet added, neighbours untouched). Integrator re-check after the merge: `gofmt -l` empty, `go vet` clean, `go test -count=1 ./internal/toolboxcommands/` ok (26.9 s).
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md*

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only, committed by the coordinator)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-731-preflight-probe.sh) and the repository gate green at df1bf423 (158 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)

**What was done:** A new check flag on ai-report-index (handleAIReportIndex branch, new function reportIndexCheck, a new Go test at the end of the test file) reads only ai-reports/catalog.json for the given bundle path and prints one ai-report index rebuild line when the bundle is not listed or the catalog is missing or unreadable, else nothing; it never writes. Step 8 of the ai-report action runs it for an ordinary report and prints the line just before the final bundle-path line. The hand-back merge (range 4feebfee..15ad4b3e) conflicted with REQ-730 (catalog lists second-opinion pages) at the ends of both Go files; both members' functions and tests were kept.

## Decisions

(from the builder hand-back)
- D-04 DECIDE & STATE: hand-made catalog listing the bundle prints nothing. The REQ file's D-02 text says a hand-made catalog counts as not listing; the brief (later) says decode the same and print nothing when listed. Followed the brief: the check only answers "is it listed".
- D-05 DECIDE & STATE: Step 8 says "Skip it for a revise", because the revise form reaches Step 8 for its commit rule but has already regenerated the catalog.
- D-06 DECIDE & STATE: no Verification Checklist clause added; the Step 8 paragraph is the single home, and a checklist line would restate it.
- D-07 DECIDE & STATE: `--check` with an empty argument or more than one argument returns the usage result (tested with `--check`, `--check ""`, `--check a b`).

(integrator)
- D-08 DECIDE & STATE: merge conflicts with REQ-730 at the ends of `report_index.go` and `report_index_test.go` resolved by keeping both sides in order (REQ-730's functions/tests first, REQ-731's after). The usage string auto-merged to `[--find <topic...> | --check <bundle>]`.
- D-09 DECIDE & STATE: release bump is minor (a new read-only CLI flag is additive), matching REQ-728 and REQ-730 in this run.
<!-- D-XX counter: last used D-09. Next decision: D-10. -->

## Discovered Tasks

(from the builder hand-back)
- `skills/do-work-toolbox/actions/architecture-report.md:119-122` and `skills/do-work-toolbox/actions/stakeholder-report.md:51-54` also write bundles and run the organize check, but do not run the new catalog check (D-03). Note architecture-report caps output at "the one allowed line beyond the eight", so adding the catalog line there needs a wording decision. → report only


## Qualification

**Gate records:** `advance --diff-range 4feebfee..15ad4b3e` (qualify, provenance merged_range) left two warnings only: SCOPE-DECLARED-NOT-TOUCHED for `CHANGELOG.md` and `VERSION`, which the Scope marks "integrator only" and which land in the finalization transaction. A first run also refused QUALIFY-CLAIMED-PATH-MISSING / SCOPE-UNDECLARED-TOUCH because the Files changed list carried backticked function names after each path; the list now holds paths only. No debug artifacts: the diff has no TODO, FIXME or print lines.

**Requirement trace (against `git diff 4feebfee..15ad4b3e`, 4 files, +106/-1):**
- Step 8 checks whether the new bundle is in the catalog and prints the `ai-report index` line when it is not: `skills/do-work-toolbox/actions/ai-report.md` Step 8, new paragraph and fenced `ai-report-index --check ai-reports/<folder>` command after the organize check. ✓
- The check is a program step on an existing command: `report_index.go` `handleAIReportIndex` gains `--check <bundle>` (exactly one non-empty argument, trailing `/` trimmed), answered by `reportIndexCheck`, which decodes only `ai-reports/catalog.json`. ✓
- The line prints just before the final bundle-path line, like the organize line: the paragraph says "Print its line, if any, with the organize line, just before the final bundle-path line." ✓
- The guide's "Report Catalog" says a new report enters at the next rebuild and the summary reminds you: new bullet in `skills/do-work-toolbox/docs/ai-report-guide.md`. ✓
- Constraint, creating a report never rebuilds or commits the catalog, and no bundle is edited: `reportIndexCheck` only reads; the test snapshots the whole `ai-reports/` tree before and after every `--check` call and fails on any change. ✓
- Constraint, catalog missing prints the same line: `os.ReadFile` error returns the line; the test's first case (no catalog) pins it. ✓
- GREEN when: prints the line when absent, nothing when listed: the test pins both, and a manual run on the merged tree confirms both (see Testing).

**Scope comparison (Route B):** declared 6 files; touched the 4 project files declared. `CHANGELOG.md` and `VERSION` follow at finalization. No undeclared file. The "Files I will NOT touch" list holds: `report_organize.go`, `architecture-report.md`, `stakeholder-report.md`, bundles and the catalog pair are untouched. The Scope line for `ai-report.md` also said "catalog forms and checklist mention it"; the builder chose a single Step 8 home instead (D-06), which I accept: a checklist line would restate the paragraph.

## Testing

**Tests run:**
- `go test -count=1 ./internal/toolboxcommands/` (in `skills/do-work/tools/do-work-cli`) after the conflict resolution → `ok` (26.9 s); `gofmt -l` empty, `go vet` clean.
- Repository gate `bash _dev/tests/maintainer-verify.sh` via `helpers/gate.sh` at merge `15ad4b3e3dae6992bf369b2d8181f9b94711f94f` (detached drain-head checkout, stage reuse off), started at 1-minute load 4.53 → `GATE-EXIT=0 WALL=141s`.
- GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-731-probe.sh` through `advance --probe-file` → raw probe status 0; green gate recorded for the merge revision.
- Contract regressions on a quiet machine (inside the gate's aggregate contract suite): `Contract regression checks passed.`, `core-checks.sh 7s`, `replace-text-section.sh 9s` (limit 30 s). The builder's two red runs at load 39 (core-checks 36 s, replace-text-section 47 s) were the per-file budget under load only.
- Manual check on the merged tree (main checkout, `--format text`):
  - `ai-report-index --check ai-reports/2026-11-14_1547_proposal-tinka-review-surface` (listed in the 49-bundle catalog) → no output, exit 0.
  - `ai-report-index --check ai-reports/2026-11-15_0100_new-report-not-indexed` (not in the catalog) → `ai-reports/2026-11-15_0100_new-report-not-indexed is not in ai-reports/catalog.json yet: run do-work-toolbox ai-report index to rebuild the catalog.`, exit 0. Every real top-level bundle in this checkout is already listed, so the unlisted case uses a path with no folder; the check reads only the catalog, and the Go test covers a real new bundle folder.

**Result:** ✓ All passing

**Red-green validation:**
- `TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing` (`report_index_test.go`): ✗ at RED commit `a6ac1be6` (`TOOLBOX-USAGE`, `Usage: ai-report-index [--find <topic...>]`) → ✓ at `0e798128` and at the merge (from the hand-back, re-run here).

**New tests added:**
- `TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing`: no catalog → line; bundle added after the catalog → line (with and without trailing `/`); listed bundle → empty; after rebuild → empty; `--check`, `--check ""`, `--check a b` → usage; no `--check` call changes the `ai-reports/` tree.

**Heavy verification plan:**
- Range: 4feebfee6fe5954ec9286d2840f2e0eeeb41d0ff..15ad4b3e3dae6992bf369b2d8181f9b94711f94f
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — report_index.go and its test matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — report_index.go and its test matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — report_index.go and its test matched subtree skills/do-work/tools/do-work-cli

*Verified by work action (integrator, run work-2026-11-15-012821)*

## Review

**Overall: 95%** | 2026-11-15T01:46:56Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `--check` does not normalize its argument (`report_index.go:112-113`), so `./ai-reports/<folder>`, an absolute path, `//` or `ai-reports/<folder>/index.html` print a false rebuild line for a listed bundle. The forms Step 8 prescribes work. — impact-user-visible → report only
- F2. The REQ's D-02 entry (line 60) still says a hand-made catalog counts as not listing the bundle. D-04 and the code do the opposite, and D-02 was not amended. — impact-negligible → report only
- F3 (Nit). The test reaches "catalog without the new bundle" by deleting and rewriting bundle `b` (`report_index_test.go:474-488`). Writing `b` after the first rebuild is shorter. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: focused Go tests ok at `15ad4b3e`, plus a scratch-repo run of missing, listed, new, hand-made and unparseable catalogs. Merge hunks match the builder's diff byte for byte.
**Restatement sweep:** nothing redefined (a new read-only flag only; the usage string has no restatement outside `report_index.go`, checked with grep)
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

Full reviewer report: `do-work/runs/work-2026-11-15-012821/REQ-731-review.md` (reviewer agent, model opus, orchestrated mode, Route B standard review).

## Heavy Verification Plan

- Base: 4feebfee6fe5954ec9286d2840f2e0eeeb41d0ff
- Target: 15ad4b3e3dae6992bf369b2d8181f9b94711f94f
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — report_index.go and report_index_test.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — report_index.go and report_index_test.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — report_index.go and report_index_test.go matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target: 15ad4b3e3dae6992bf369b2d8181f9b94711f94f; execution revision 15ad4b3e3dae6992bf369b2d8181f9b94711f94f (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set, started at 1-minute load 4.98)
- do-work-cli-integrations: executed, exit 0, 66 s
- staged-skills: executed, exit 0, 37 s
- updater: executed, exit 0, 68 s
- installer: executed, exit 0, 28 s
- No `HEAVY-RUN-LANE-SKIPPED`; run findings empty.

## Lessons Learned

- **What worked:** the builder kept every edit inside `handleAIReportIndex` and at the ends of both Go files, as the Exploration asked, so the conflict with REQ-730 was two end-of-file appends that resolved by keeping both sides in order.
- **What didn't:** the first qualify run refused because the Implementation Summary's file list carried backticked function names after each path; the qualifier reads every backticked word there as a claimed path. Keep that list to paths only.
- **Worth knowing:** the builder's contract-regressions reds at load 39 (two different files over the 30 s budget) passed at load 4.5 inside the gate (7 s and 9 s). The builder proposed no lesson bullet, so no satellite changed.

## Orientation

`ai-report` Step 8 now runs two read-only CLI checks before the final bundle-path line: `ai-report-organize --check` and `ai-report-index --check <bundle>`; the catalog check lives in `report_index.go` beside the catalog rebuild and reads only `ai-reports/catalog.json`.

## Timing

Observed 2026-11-14T23:35:12Z to 2026-11-15T01:51:10Z: 2h 15m 58s total, 1h 58m 49s attributed across 5 events, 17m 09s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 1h 47m 52s | 1 |
| verification-gate | 6m 33s | 2 |
| review | 3m 12s | 1 |
| handback-merge | 1m 12s | 1 |

Slowest stage: builder-work / builder worktree build, 1h 47m 52s, outcome success.
Slowest command: verification-gate / heavy drain, 3m 41s, exit 0, .
