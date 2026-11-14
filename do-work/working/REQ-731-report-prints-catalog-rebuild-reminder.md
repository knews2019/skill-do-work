---
id: REQ-731
title: 'ai-report prints the catalog rebuild command when a new report is missing from the catalog'
status: claimed
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
- [~] What does the check print, and for which path form? → **D-02**: Builder chose: the argument is the bundle path as the catalog stores it (`ai-reports/<folder>`, a trailing `/` trimmed); the printed line is one line naming the bundle and `do-work-toolbox ai-report index`, in the style of the organize hint (`report_organize.go:35`). A hand-made catalog (no generator marker) counts as not listing the bundle, so the line still prints and the rebuild then reports its own refusal. Reasoning: new bundles are always written directly under `ai-reports/`, so no month-folder form is needed. Value: one exact string the test pins. Risk: low; wording only.
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md*

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only, committed by the coordinator)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-731-preflight-probe.sh) and the repository gate green at df1bf423 (158 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*
