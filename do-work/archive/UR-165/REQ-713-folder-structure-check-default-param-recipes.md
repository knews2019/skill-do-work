---
id: REQ-713
title: 'folder-structure-check sees justfile recipes with default parameter values'
status: completed
created_at: 2026-11-09T21:47:36Z
user_request: UR-165
domain: backend
prime_files: []
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-09T21:52:19Z
  basis:
    - trivial short-circuit
claimed_at: 2026-11-09T21:51:56Z
route: A
dispatch_at: 2026-11-09T21:52:25Z
builder_handback_at: 2026-11-09T21:53:42Z
integration_at: 2026-11-09T21:53:48Z
review_at: 2026-11-09T22:02:40Z
kb_status: pending
commit: 013ea54efafa3140594defdc2acdef58537ba0ea
heavy_verified_at: 2026-11-09T22:05:25Z
heavy_verified_revision: 013ea54efafa3140594defdc2acdef58537ba0ea
completed_at: 2026-11-09T22:05:33Z
release_at: 2026-11-09T22:05:33Z
---
# folder-structure-check sees justfile recipes with default parameter values
## What
Fix `folder-structure-check` so a justfile recipe whose parameters have default values (for example `run-kanban $port="8090":`) is recognized as a recipe and its body lines are counted. Today the header regex `justRecipeHeader` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go` stops at the `=`, so the recipe is skipped.
## Detailed Requirements
- Recognize recipe headers with quoted or unquoted default parameter values, `*args`/`+args` variadics and `$`-exported parameters.
- Keep the existing exclusions: `name := value` variables and the set/alias/export/import/mod directives are never recipes.
- Add a case to `TestFolderStructureCheckReportsEachGapClass` with a multi-line recipe that has a quoted default parameter, and show it fails without the fix.
- Layout rule this enforces: a justfile recipe holds no logic and calls `cmd/<name>.sh` in one line.
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository; the root justfile's `run-kanban $port="8090":` recipe has 3 body lines.
**Why RED now:** The output reports only `justfile: kanban-static`; `run-kanban` is missing because the header regex fails on the `=`.
**GREEN when:** The output reports `FOLDER-STRUCTURE-JUST-LOGIC justfile: run-kanban`, the new test passes, and `go test ./internal/toolboxcommands/` passes.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24067 tokens, `slugged: partial`, so it cannot be narrowed; matched because the fix is in do-work-cli internals.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Replace the header regex with a quote-aware scan: leading name, then the first colon outside quotes; a colon followed by an equals sign is an assignment. Extend the two folder-structure tests first.
- [x] **[APPLY]:** Changed only folder_structure.go and folder_structure_test.go, RED commit 06761c55 then GREEN b9250783.
- [x] **[UNIFY]:** git diff --stat shows the two declared files (+42 -4). gofmt -l and go vet clean on internal/toolboxcommands; go test for the package ok. No debug output added.
*Source: do-work capture-request: Fix folder-structure-check so it sees justfile recipes whose parameters have default values. Today justRecipeHeader in skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go stops at the "=" inside a parameter like `run-kanban $port="8090":`, so that recipe's 3 body lines are never counted (layout rule: a justfile recipe holds no logic and calls cmd/<name>.sh in one line). Keep the existing exclusions for `name := value` variables and set/alias/export/import/mod directives. Add a case to TestFolderStructureCheckReportsEachGapClass with a multi-line recipe that has a quoted default parameter, and show the test fails without the fix. Verify: go test ./internal/toolboxcommands/ passes, and running the check on this repo reports `justfile: run-kanban`.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the file, the regex and the test to extend. It is a one-function bug fix with a clear repro.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- D-01 (DECIDE & STATE): A quote-aware scan replaces the regex, because a regex cannot skip a colon inside a quoted default (`url='http://localhost'`) without lookarounds, which Go's regexp package does not have. Superseded in part by D-02: the scan is the existing one, not a new one.
- D-02 (DECIDE & STATE): Reuse `justDefinitionName` from `internal/managedsection` through an exported one-line wrapper, although the REQ names only `folder_structure.go`. One scanner means one place to fix Just syntax (review F1).

## Discovered Tasks

- The comment at the top of `internal/managedsection/just_definitions.go` and line 12 of `skills/do-work/tools/prime-do-work-update.md` still say the scanner has one job; it now has two callers (review F6). impact-negligible → report only
- The body-line count still misreads an unindented line inside a multi-line string, a header or dependency list continued with a backslash, and a first-line recipe after a byte-order mark (review F3, F4). All are older than this REQ. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go` (modified)

**What was done:** The justfile check no longer uses a regex that stopped at the first equals sign. It now asks the scanner the module already ships for the managed justfile section (`justDefinitionName`, reached through the new exported wrapper `JustLineDefinitionName`), which knows Just's quoted, triple-quoted and backtick forms. Recipes whose parameters have default values are counted, assignments and `set`/`export`/`import`/`mod` lines still define nothing, and the check's own header regex and directive list are deleted. The first version of the fix added a second scanner; review finding F1 replaced it with this reuse before release.

## Qualification

**Range:** `14a13316..013ea54e` (two merges of the REQ branch: RED 06761c55, GREEN b9250783, review fix 08852425; first merge 218404ff)
**Mechanical gates:** qualify ✓ satisfied on `14a13316..218404ff`; the review fix was re-read by the orchestrator and re-checked by the reviewer on `218404ff..013ea54e`. Route A has no Scope comparison.

**Requirement trace (orchestrator, read against the diff):**
- Default parameter values: the new fixture recipe `run-kanban $port="8090" url='http://localhost' +rest: lint` (quoted default, a default holding a colon, an exported parameter, a variadic and a dependency) is reported with its three body lines.
- Exclusions kept: `port := "8080"`, `url:="http://localhost"` and `set shell := [...]` in the passing-layout fixture report nothing. After the review fix the shared scanner returns nothing for assignments and for `set`, `export`, `import` and `mod` lines; an `alias` line returns a name that has no body lines, so it cannot become a gap.
- No false positive for a one-line recipe with defaults: `serve $port="8090" *args:` under a `[positional-arguments]` attribute passes.
- The new case lives in `TestFolderStructureCheckReportsEachGapClass` and failed at the RED commit.

**Live flow:** the test calls the registered `folder-structure-check` handler, and the built CLI run on this repository now reports `justfile: run-kanban` beside `justfile: kanban-static`.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/2026-11-09-layout-queue/REQ-713-probe.sh` (the `FolderStructure` and handler-registration tests in `internal/toolboxcommands`), then `bash _dev/tests/maintainer-verify.sh` from a detached checkout at each merge revision
**Result:** ✓ Focused probe green. Repository gate exit 0 at the first merge 218404ff (gate wall 125s) and exit 0 at the final merge 013ea54e (gate wall 122s; do-work-cli module 54s, 974 tests; slowest file `finalization_recovery_test.go` 16.0s, under the 30s budget). Both Go stages reported EXECUTING, not REUSED.

**Red-green validation:**
- `TestFolderStructureCheckReportsEachGapClass` (`folder_structure_test.go`): ✗ at commit 06761c55, the wanted `FOLDER-STRUCTURE-JUST-LOGIC justfile: run-kanban` was missing → ✓ at commit b9250783 and at the final merge. This is the captured RED case: a recipe header `run-kanban $port="8090" ...:` with three body lines that the check skipped.
- The captured GREEN on this repository: the built CLI at 013ea54e reports `justfile: run-kanban` (3 body lines) beside `justfile: kanban-static`.
- No false positive: `TestFolderStructureCheckPassesTheTargetLayout` holds `url:="http://localhost"` and a one-line recipe `serve $port="8090" *args:` and still reports zero gaps.

**Existing tests updated:**
- `TestFolderStructureCheckPassesTheTargetLayout` and `TestFolderStructureCheckReportsEachGapClass` (new fixture lines; no assertion was loosened)

**Heavy verification plan:**
- Range: 14a13316eab7ad46f115c7e68a286276af4e244b..013ea54efafa3140594defdc2acdef58537ba0ea
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the changed Go files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — every changed path matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the changed Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the changed Go files matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 96%** | 2026-11-09T22:02:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Reviewed range: `14a13316..013ea54e` (first review covered `14a13316..218404ff`; the re-check covered `218404ff..013ea54e`).

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 Duplicate justfile header scanner: the first fix added `justRecipeName` in `folder_structure.go`, a weaker copy of `justDefinitionName` in `internal/managedsection/just_definitions.go`. FIXED in the re-check (commit 08852425): `folder_structure.go:226` now calls the exported wrapper `managedsection.JustLineDefinitionName`, and `justRecipeName`, `justLeadingName` and `justDirective` are deleted. The reuse is correct for line-by-line use: the loop still skips blank, comment and indented lines first; `set`, `export`, `import`, `mod` and assignment lines return ""; an `alias` line returns the alias name, which has no body lines and cannot become a gap. `managedsection` imports no internal package, so there is no import cycle — impact-rule-change → report only

**Minor findings:**
- F2 Missed recipe for a default with an odd number of `\"` or a triple-quoted default holding its own quote character. FIXED by the F1 change: `multi sep="\"":` and `multi a="""x"y""":` are now reported — impact-negligible → report only
- F3 False gap, older than this REQ and still present: an unindented `Name:` line inside a multi-line string assignment followed by two indented lines is reported as a recipe (`justfile: Usage`), and a dependency list continued with `\` onto an indented line is counted as a body line. The wrapper classifies one line at a time, so the multi-line literal state that `JustDefinitionNames` keeps is not used here — impact-negligible → report only
- F4 Missed recipe, partly fixed: recipes named `import`, `export`, `set`, `mod` or `alias` are now counted (just 1.46.0 accepts them; real directives are still skipped). Still missed, as before this REQ: a header continued with `\` before its colon, a header whose default string spans lines, and a first-line recipe in a justfile that starts with a byte-order mark — impact-negligible → report only
- F5 (nit) Untested backtick and escape branches. CLOSED by the F1 change: `TestJustDefinitionScannerClassifiesEveryLiteralForm` in `internal/managedsection/just_definitions_test.go` already pins those forms — impact-negligible → report only
- F6 (nit) Stale purpose statement after the reuse: the comment at `internal/managedsection/just_definitions.go:8-11` says the scanner "exists for one job" (reserved-name collisions), and `skills/do-work/tools/prime-do-work-update.md:12` calls it the "reserved-recipe collision scanner". It now has a second caller, `folder-structure-check`; only the wrapper comment at `just_definitions.go:56-58` says so — impact-negligible → report only

**Acceptance:** Pass — Implementation and Integration stages on the merged tree at 013ea54e: `go test -count=1 -run FolderStructure ./internal/toolboxcommands/` ok, `go test -count=1 ./internal/managedsection/` ok, `gofmt -l` and `go vet` clean on both packages; RED reproduced on an extracted copy of 06761c55 (`TestFolderStructureCheckReportsEachGapClass` fails, `justfile: run-kanban` missing); the built CLI on this repository reports `justfile: run-kanban` (3 body lines) beside `justfile: kanban-static`; 32 justfile edge cases compared against `just 1.46.0 --dump`, the pre-REQ binary, the first-fix binary and the final binary: the final binary differs from the first fix on three cases only, all toward the real parser, and no case that was right before is wrong now.
**Restatement sweep:** redefined what `folder-structure-check` counts as a justfile recipe header, and widened the purpose of the managed-section Just scanner to a second caller; `skills/do-work-toolbox/actions/folder-structure.md:46` and `:71`, `CHANGELOG.md:25`, `skills/do-work/CHANGELOG.md:25` and `README.md:126` still agree and none of them lists the deleted directive names; the two purpose statements in F6 are stale.
**Suggested testing:** 2 items — (1) Integration: this review ran only the `FolderStructure` tests and the `managedsection` package, so take the full `go test ./internal/toolboxcommands/` and `./internal/suiteinstall/` results from the orchestrator's gate run; (2) Deployment: unassessed, run `folder-structure-check` from an installed suite copy after the release and confirm `justfile: run-kanban` is reported.
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: 14a13316eab7ad46f115c7e68a286276af4e244b
- Target revision: 013ea54efafa3140594defdc2acdef58537ba0ea
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the changed Go files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — every changed path matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the changed Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the changed Go files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: 013ea54efafa3140594defdc2acdef58537ba0ea
- Execution revision: 013ea54efafa3140594defdc2acdef58537ba0ea (detached checkout)
- do-work-cli-integrations: exit 0, 55s, executed
- staged-skills: exit 0, 34s, executed
- updater: exit 0, 65s, executed
- installer: exit 0, 26s, executed

## Lessons Learned

**What worked:** Launching the reviewer right after qualification, while the gate ran. It compared 29 justfile cases against the real `just` binary and found the existing scanner in four minutes.
**What didn't:** The first fix wrote a second justfile header scan inside `folder_structure.go`. The module already had `justDefinitionName` in `internal/managedsection`, which handles escaped quotes and triple-quoted strings that the new scan got wrong. Searching the module for "just" before writing a parser would have found it. The reuse cost one more merge and one more gate run.
**Worth knowing:** `folder-structure-check` reads a justfile line by line, so it still misreads an unindented line inside a multi-line string and a header continued with a backslash. `JustDefinitionNames` in the same package carries literal state across lines; a caller that needs that precision should build on it, not on the per-line wrapper.

## Orientation

Now `folder-structure-check` counts the body lines of every justfile recipe, including recipes whose parameters have default values such as `run-kanban $port="8090":`. It lives in the toolbox commands of the do-work CLI and reuses the justfile scanner of the managed-section package. No prime file covers this area, so there is no staleness check to run.
