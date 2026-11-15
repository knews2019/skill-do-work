---
id: REQ-728
title: 'ai-report-judge renders one named page of a bundle'
status: completed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-729, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-14T23:28:12Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-14T23:30:23Z
  basis:
    - Route B
    - 6-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
dispatch_at: 2026-11-14T23:35:10Z
builder_handback_at: 2026-11-15T00:48:21Z
integration_at: 2026-11-15T00:48:57Z
review_at: 2026-11-15T01:04:46Z
kb_status: pending
commit: b7e6a08e20d22b31aaa3024e1ee287fe24ced170
heavy_verified_at: 2026-11-15T01:04:52Z
heavy_verified_revision: b7e6a08e20d22b31aaa3024e1ee287fe24ced170
completed_at: 2026-11-15T01:05:05Z
release_at: 2026-11-15T01:05:05Z
---
# ai-report-judge Renders One Named Page of a Bundle
## What
Add an option to `ai-report-judge`, for example `--page <file.html>`, that renders one named page of a bundle in place. Today the judge renders only `index.html`.
## Why
To check or capture a second-opinion page such as `index-v1c01-<name>.html`, an agent must copy the whole bundle and write that page over the copy's `index.html`. That copy step was needed three times on 2026-11-14.
## Context
- The judge renders only `index.html`: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go`, lines 25, 111 and 168.
- `secondOpinionRender` in `report_second_opinion.go` (from line 1706) already does the copy-and-swap internally. Reuse or replace it as the builder judges.
## Detailed Requirements
- The new option renders the named page of the bundle, without copying the bundle and without writing inside it.
- A named page that is missing or outside the bundle is refused with a typed finding.
- Files: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` and its test file; the judge form in `skills/do-work-toolbox/actions/ai-report.md`; `skills/do-work-toolbox/docs/ai-report-guide.md` where the judge is described; `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
## Constraints
- Without the new option the judge behaves exactly as today: same captures, same `judge.json`, same exit codes.
- The judge never writes inside the bundle. Bundles stay byte-identical.
- The second-opinion finalize path keeps its verdicts for the same pages.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** A Go test runs the judge with `--page index-v1c01-x.html` on a fixture bundle with two pages.
**Why RED now:** The judge fails on the unknown option.
**GREEN when:** The test passes; the result's served address ends with the named page and the captures show that page; a missing or outside page is refused with a typed finding; the fixture bundle's digest is unchanged; `go test ./...` passes in the do-work-cli module.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- Certainty: firm on the outcome. The option's name is the builder's choice; `--page` is a suggestion.
- The relative `--out` rule of the same command is its own request in this batch.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes do-work-cli Go code for reports.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add `--page` beside `--out`/`--browser`; check the page after the bundle resolves (IsLocal, extension, regular file, EvalSymlinks + judgePathInside); thread the page into `judgeRender` for `ServedURL`; one-argument change at the second-opinion caller; docs name the option and the refusal code. (from the builder hand-back)
- [x] **[APPLY]:** Done as planned in two commits (`6e40562f` RED tests, `661361a4` implementation); no other files. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat 1de00505..HEAD` showed exactly 5 files (ai-report.md 4 lines, ai-report-guide.md 4, report_judge.go 47, report_judge_test.go +149, report_second_opinion.go 2); checked each diff for debug output, gofmt, and that the `--out` block (REQ-729 area) is untouched. (from the builder hand-back)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-judge-renders-one-named-page.md*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and the files are named, but the change adds a path-containment rule (missing or outside page, symlinks) to a browser-driving command, shares `report_judge.go` with REQ-729, and must keep the second-opinion finalize path and today's no-option behavior unchanged. Exploration of the render path and its one internal caller is needed before building.

**Planning:** Not required

## Open Questions

No `## Open Questions` section was captured. Choices the capture left to the builder are recorded here as decisions so the builder does not re-derive them.

- [~] Option name? → **D-01**: Builder chose: `--page <name>`, also `--page=<name>`, matching the existing `--out`/`--browser` parsing through `optionValue`. Reasoning: the REQ suggests `--page` and the parser already handles both forms. Value: one obvious name, no new parsing shape. Risk: low, a rename later is a one-line change plus docs.
- [~] What does `--page` accept? → **D-02**: Builder chose: a path relative to the bundle, slash-separated, naming a regular `.html`/`.htm` file inside it (a basename like `index-v1c01-x.html` is the common case; a sub-folder page is allowed). An absolute path, a path that escapes with `..`, a symlink whose target resolves outside the bundle, a directory, or a missing file is refused with one new typed finding code `AI-REPORT-JUDGE-PAGE` (OutcomeFailure, exit 2, no browser launched, no judge.json written, like today's usage refusals). Reasoning: the REQ demands "a typed finding"; the generic TOOLBOX-USAGE code does not say which input was wrong. Value: callers can branch on the code. Risk: low; a new code is additive.
- [~] Must `index.html` still exist when `--page` names another page? → **D-03**: Builder chose: yes, keep the bundle check (a bundle is a directory holding index.html) unchanged for both forms. Reasoning: a second-opinion page always sits in a bundle that has index.html, and keeping the check avoids a second code path. Value: no-option behavior provably unchanged. Risk: low.
- [~] Reuse `secondOpinionRender` or change it? → **D-04**: Builder chose: leave `secondOpinionRender` and the finalize path untouched; add a page parameter to `judgeRender` (callers pass "index.html" where they render index today). Reasoning: the finalize path judges page bytes through a copy whose index.html is the new page, and the constraint says its verdicts must not change; switching it to in-place rendering would change what relative links resolve against for pages in sub-folders and is not asked for. Value: zero risk to second-opinion verdicts. Risk: the copy-and-swap stays as duplicated machinery; a later REQ can switch finalize to `--page` rendering once this ships.
- [~] judge.json shape with `--page`? → **D-05**: Builder chose: no new field and no schema bump; `served_url` ends with the named page (URL-escaped path) and `bundle` stays the bundle directory. Captures keep their four names. Reasoning: the REQ names the served address as the proof; adding a field is not asked. Value: readers of judge.json keep working. Risk: low; a reader that wants the page name reads the end of `served_url`.

<!-- D-XX counter: last used D-05. Next decision: D-06. -->

## Decisions

(from the builder hand-back; D-01 to D-05 are in Open Questions above)
- **D-06** (DECIDE & STATE): the extension check (`.html`/`.htm`, case-insensitive) runs before the file stat, so `--page notes.txt` is refused as "not an HTML file" even when it exists. All refusals share the one code `AI-REPORT-JUDGE-PAGE`; the evidence text says which rule failed.
- **D-07** (DECIDE & STATE): the refusal uses a small sibling `judgeFailureResultWithStop` so the stop reason reads "--page names no HTML page inside the bundle" instead of "the render check could not run"; `judgeFailureResult` now delegates to it (same output for existing callers).
- **D-08** (DECIDE & STATE): the page check runs before the browser lookup and before `--out` is created, so a refusal writes nothing anywhere (the test asserts no judge.json and an unchanged bundle digest).
- **D-09** (DECIDE & STATE): the refusal table has the four cases the brief names; directory and wrong-extension refusals are covered by the code but not given their own sub-tests (YAGNI; same code path).

(integrator)
- **D-10** (DECIDE & STATE): the merge conflict in `judgeUsage` keeps both texts: the builder's ` [--page <file.html>]` after `<bundle-dir>` and REQ-729's closing sentence "a relative <bundle-dir> or --out starts at --repo-root". They do not contradict: that sentence names only `<bundle-dir>` and `--out`, and the new Step 7 paragraph says `--page` is relative to the bundle. Value: both members' text survives unchanged. Risk: low; a reader of the one-line usage could guess `--page` also starts at `--repo-root`, but the refusal message and Step 7 say otherwise.
- **D-11** (DECIDE & STATE, integrator, after review): fix review finding F1 in this REQ instead of reporting it. After `judgeCheckPage` passes, the handler cleans the page to slash form (`filepath.ToSlash(filepath.Clean(filepath.FromSlash(page)))`) and serves that, so `served_url` and the browser use the checked path, not the raw text. Fixed on the builder branch (`1e80dacd`), re-merged with the same base as `b7e6a08e`. Value: no false pass on Windows (the CLI is vetted for `GOOS=windows`) for a backslash page, and a clean `served_url` everywhere. Risk: low; the default `index.html` is unchanged by the clean step. Cost: one re-gate and one heavy re-drain.
<!-- D-XX counter: last used D-11. Next decision: D-12. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Done by the pre-dispatch agent (no subagent) at main HEAD 1389e2be.

**Required-lessons consult.** `required_lessons` holds `_dev/primes/lessons-releases.md` (856 tokens, read in full). Its two families (canonical links after a file move, manifest ownership) do not apply: this REQ moves no file and edits no manifest. The index rows for `lessons-do-work-cli.md` (24500 tokens) and `lessons-action-files.md` (8761 tokens) match, but both are `slugged: partial`, so they cannot be narrowed and stay dropped for budget (section below, unchanged). Builder: read the do-work-cli lessons bullets for `guard-on-one-entry-path`, `existence-is-not-identity` and `silent-skip-reads-as-red` by grep if time allows; they match this change (a containment guard, a page path checked by existence, a heavy test that skips).

**Render path (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go`).**
- `report_judge.go:39` `judgeUsage` is the usage string; add `[--page <file.html>]`. REQ-729 edits the same constant (relative `--out` rule) — a text seam.
- `report_judge.go:88-109` `handleAIReportJudge` parses args with `optionValue` (`commands.go:91`, handles `--x v` and `--x=v`). Add a `--page` case beside `--out`/`--browser`.
- `report_judge.go:110-117` resolves the bundle against `ctx.RepositoryRoot`, requires `index.html`, then `filepath.EvalSymlinks(bundle)` into `resolvedBundle`. The page check belongs right after line 117: reject absolute paths and `..` escapes (`filepath.IsLocal` on the OS-form path), join to the bundle, `os.Stat` for a regular file, `filepath.EvalSymlinks` the page and require `judgePathInside(resolvedPage, resolvedBundle)` (`report_judge.go:307`). Refusal = new code `AI-REPORT-JUDGE-PAGE` via `judgeFailureResult` (`report_judge.go:279`) or a sibling helper, before any browser lookup or `--out` creation.
- `report_judge.go:132-138` resolves and creates `--out` — REQ-729's area; do not touch.
- `report_judge.go:147` calls `judgeRender(&report, bundle, output, browser)`; `report_judge.go:159-168` `judgeRender` serves `http.Dir(bundle)` and hard-codes `report.ServedURL = origin + "/index.html"` (line 168). Add a page argument: `ServedURL = origin + "/" + (&url.URL{Path: slashPage}).EscapedPath()` so a name with spaces still loads. Links are resolved against `document.baseURI`, so relative links of a sub-folder page resolve correctly with no other change.
- `report_judge.go:509-510` `renderPass` names captures `<viewport>-<scheme>.png`; keep those names (D-05).
- `report_judge.go:24-27` header comment says "renders its index.html"; update it.

**Second caller.** `report_second_opinion.go:1706-1771` `secondOpinionRender` copies the bundle into a scratch folder, writes the new page as the copy's `index.html`, and calls `judgeRender(&report, renderBundle, captures, browser)` at line 1752. It is the only other caller (grep `judgeRender`). Under D-04 it changes by one argument (`"index.html"`) and nothing else. Its tests stub `secondOpinionJudge` (`report_second_opinion.go:46`; `report_second_opinion_test.go:208,697`; `report_second_opinion_tree_test.go:32`), so its verdict behavior is pinned by those tests plus the unchanged call.

**Tests (`report_judge_test.go`, 167 lines).** Browser cases call `requireJudgeBrowser` (line 22: skips unless `DO_WORK_HEAVY_TESTS=1` and a browser is found; Google Chrome is at `/Applications/Google Chrome.app` on this machine). `runJudgeOnBundle` (line 35) builds a bundle with index.html plus extra files and always reads judge.json, so a refusal test needs its own small call (a refusal writes no judge.json). `TestAIReportJudgePassesCleanBundleAndStopsServer` (line 109) shows the bundle byte-check pattern; `TestAIReportJudgeReportsSkippedWithoutBrowser` (line 152) is the light no-browser case.

**Docs.** `skills/do-work-toolbox/actions/ai-report.md:36` (the `judge <bundle-dir>` form line) and `:176-186` (Step 7, the command block at 181 and its paragraph at 184); `skills/do-work-toolbox/docs/ai-report-guide.md:59` (judge paragraph) and `:80` (command list). `actions/help.md:18`, `stakeholder-report.md:44` and `architecture-report` call the judge without the option and need no change.

**Release.** `VERSION` is 0.318.1. A new option is user-visible, so expect a minor bump (`work-reference.md` → Changelog Entry Procedure); the integrator writes `CHANGELOG.md`, its mirror `skills/do-work/CHANGELOG.md`, and `VERSION`.

*Generated by pre-dispatch agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (modify) — the page option, its containment check with a typed refusal, and the page argument of judgeRender
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (modify) — two new tests appended at the end
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — one call site passes index.html to judgeRender
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — judge form line and Step 7 name the option
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — judge paragraph and command list name the option
- `CHANGELOG.md` (modify, integrator only) — release entry
- `skills/do-work/CHANGELOG.md` (modify, integrator only) — byte-identical mirror
- `VERSION` (modify, integrator only) — minor bump

**Files I will NOT touch:** the relative out-directory resolution in report_judge.go (REQ-729), `report_index.go` (REQ-730), `actions/help.md`, `actions/stakeholder-report.md`, `actions/architecture-report.md`, any `ai-reports/` bundle, any `do-work/` path.

**Acceptance criteria (restated from REQ):**
- [ ] The judge with the page option renders the named page of the bundle, without copying the bundle and without writing inside it.
- [ ] The result's served address ends with the named page and the captures show that page.
- [ ] A named page that is missing or outside the bundle is refused with a typed finding.
- [ ] The fixture bundle's digest is unchanged after the run.
- [ ] Without the option the judge behaves exactly as today: same captures, same judge.json, same exit codes.
- [ ] The second-opinion finalize path keeps its verdicts for the same pages; `go test ./...` passes in the do-work-cli module.

## Pre-Flight

**Git:** ✓ No uncommitted project files at dispatch (run artifacts under do-work/ only, committed by the coordinator)
**Tests baseline:** ✓ Focused pre-build probe passing (do-work/runs/work-2026-11-15-012821/REQ-728-preflight-probe.sh) and the repository gate green at df1bf423 (158 s, detached checkout, stage reuse off)
**Dependencies:** ✓ Installed

*Checked by work action (coordinator, run work-2026-11-15-012821)*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (modified: header comment, judgeCodePage constant, judgeUsage, --page parse case, page check after the bundle check, judgeRender page parameter and ServedURL, new judgeFailureResultWithStop and judgeCheckPage)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (modified: two tests and two digest helpers appended; strings import added)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified: the one judgeRender call passes "index.html")
- `skills/do-work-toolbox/actions/ai-report.md` (modified: the judge form line; one new paragraph in Step 7)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified: one sentence in the judge paragraph; the command list line)

**What was done:** `ai-report-judge --page <file.html>` renders one named page of a bundle in place: the server still serves the bundle directory, and the browser opens the named page instead of `index.html`. A page that is absolute, escapes with `..`, has no `.html`/`.htm` extension, is not a regular file, or resolves outside the bundle through a symlink is refused with `AI-REPORT-JUDGE-PAGE` before any browser starts or `--out` is created. Runs without the option are unchanged, and the second-opinion finalize path keeps its copy-and-swap (D-04).

## Discovered Tasks

(from the builder hand-back)
- `secondOpinionRender` could switch from copy-and-swap to `--page`-style in-place rendering (already named in D-04). Not impact-critical → report only.

## Qualification

**Gate records:** `advance REQ-728 --diff-range c9b1eb44..9618c836` returned the `qualify` gate `satisfied`. The first call refused `QUALIFY-CLAIMED-PATH-MISSING` ten times because the Implementation Summary file list carried backticked identifiers (`judgeCodePage`, `--page` and others); the backticks inside the parentheses were removed, nothing else changed. The `scope-drift` gate still lists `SCOPE-DECLARED-NOT-TOUCHED` for `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` and `VERSION`: the Scope marks those three "integrator only", and the finalization release writes them. No `SCOPE-UNDECLARED-TOUCH` remains.

**Diff:** `git diff c9b1eb44..9618c836 --stat`: 5 files, 198 insertions, 13 deletions. Declared vs touched: all five project files of the Scope are touched, nothing outside it; the three release files come at finalization. The `--out` resolution block (REQ-729's area) is byte-unchanged; the only shared line is `judgeUsage`, resolved by D-10.

**Requirement trace:**
- Renders the named page without a bundle copy and without writing inside it: `handleAIReportJudge` keeps `page := "index.html"` as the default and passes `page` to `judgeRender`, which still serves `http.Dir(bundle)` and opens `origin + "/" + EscapedPath(page)` (report_judge.go:177). No copy step exists on this path. `TestAIReportJudgeRendersNamedPageInPlace` compares a full bundle digest before and after.
- Served address ends with the named page and the captures show it: the same test asserts `served_url` ends with `/index-v1c01-x.html`, and asserts a phone overflow of at least 600 px that only the named page has (index.html is clean), plus four non-empty captures.
- Missing or outside page refused with a typed finding: `judgeCheckPage` (report_judge.go:300-317) refuses non-local paths, non-HTML extensions, non-regular files, and symlinks that resolve outside the resolved bundle, all as `AI-REPORT-JUDGE-PAGE`. `TestAIReportJudgeRefusesMissingOrOutsidePage` covers missing, `..` escape, absolute path, and a symlink to outside, and asserts no judge.json and an unchanged bundle.
- No-option behavior unchanged: the default page is `index.html`, so `ServedURL` is the same `origin + "/index.html"` as before; the four existing judge tests (overflow, broken image, clean pass with server stop, skipped without browser) stay in the GREEN probe.
- Second-opinion finalize path keeps its verdicts: `secondOpinionRender` changes by one argument (`"index.html"`, the page its copy-and-swap already writes), and the probe runs the four second-opinion judge-verdict tests.
- Docs: `ai-report.md` names `--page` in the judge form line and in a new Step 7 paragraph (relative to the bundle, refusal code, exit 2, no judge.json); `ai-report-guide.md` adds one sentence and the command-list form. The REQ-729 sentence ("a relative `<bundle-dir>` or `--out <dir>` starts at `--repo-root`") sits in the paragraph above and does not contradict the new one, which says `--page` is relative to the bundle.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (repository gate, stage reuse off) at merge `9618c836e1ce24fa7851cb9e3b23ac251d52a405` in the run's detached checkout, started 2026-11-15T00:49:54Z with the 1-minute load at 3.27 and no other gate running; GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-728-probe.sh` through `advance`, then once directly from ROOT.
**Result:** ✓ Gate exit 0 in 133 s (single run, no retry; go-test budget: 982 tests, slowest file 17.21 s, limit 30 s). Probe exit 0 (`REQ-728 probe PASS`, 5.8 s); `advance` recorded the green-gate and run-blocked-check records as satisfied.

**Red-green validation (tdd: true, from the hand-back):**
- `TestAIReportJudgeRendersNamedPageInPlace` and `TestAIReportJudgeRefusesMissingOrOutsidePage`: ✗ after the RED commit `6e40562f` (`TOOLBOX-USAGE ... unexpected argument "--page"`; all four refusal sub-tests "no AI-REPORT-JUDGE-PAGE finding") → ✓ after `661361a4` and at merge `9618c836` (run with `DO_WORK_HEAVY_TESTS=1`, PASS lines, no SKIP).
- `REQ-728-probe.sh`: ✗ at base `1de00505` (`ai-report.md does not name --page`) → ✓ at builder tip and at merge.
- Builder checks (from the hand-back): `gofmt -l` and `go vet ./internal/toolboxcommands/` clean; `go test -count=1 ./...` in do-work-cli all ok twice; `bash _dev/tests/action-shell-blocks.sh` and `bash _dev/tests/shipped-package-reference-contract.sh` pass; a manual `--page` run on a real second-opinion page served the named page and left `ai-reports/` clean.

**New tests added:**
- `TestAIReportJudgeRendersNamedPageInPlace` (report_judge_test.go): pins that `--page` renders the named page in place (served address, a page-only phone overflow, four captures, unchanged bundle digest).
- `TestAIReportJudgeRefusesMissingOrOutsidePage` (report_judge_test.go): pins the containment guard (missing, `..`, absolute, symlink outside → `AI-REPORT-JUDGE-PAGE`, no judge.json, unchanged bundle).

**Heavy verification plan:**
- Range: c9b1eb445acd98769aa0fb42d06caefe95dfd5ea..9618c836e1ce24fa7851cb9e3b23ac251d52a405
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — report_judge.go, report_judge_test.go and report_second_opinion.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three Go files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three Go files matched subtree skills/do-work/tools/do-work-cli

**Re-gate after review fix (D-11):** merge `b7e6a08e20d22b31aaa3024e1ee287fe24ced170` (range `c9b1eb44..b7e6a08e`), started 2026-11-15T00:57:38Z at 1-minute load 2.02 with no other gate running: gate exit 0, wall 137 s (slowest file 17.77 s, limit 30 s); probe PASS; green-gate record re-written with `record-green-gate` because `advance` was already in the review phase. Red-green for the fix: `TestAIReportJudgeRendersNamedPageInPlace` with `--page ./index-v1c01-x.html` ✗ before (`served_url "http://127.0.0.1:54239/./index-v1c01-x.html" does not name the cleaned page path`) → ✓ after; all six `TestAIReportJudge*` PASS with Chrome; `go vet` clean for darwin and `GOOS=windows`. The heavy plan for `c9b1eb44..b7e6a08e` selects the same four lanes.

*Verified by work action*

## Review

**Overall: 95%** | 2026-11-15T01:04:46Z — Approve. Independent reviewer (opus): full report `do-work/runs/work-2026-11-15-012821/REQ-728-review.md`, first pass 94% on `c9b1eb44..9618c836`, delta re-review 95% on `9618c836..b7e6a08e`.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:**
None

**Minor findings:**
- F1. `ServedURL` was built from the raw `--page` text, so on Windows `--page sub\x.html` passed the check but reached the server as `%5C` and could pass an error page; on Unix `served_url` read `/./x.html`. — impact-user-visible → fixed in this REQ (D-11, builder-branch commit `1e80dacd`, re-merge `b7e6a08e`, delta re-review approved, re-gate and heavy re-drain green)
- F2 (nit). `sameJudgeDigest` duplicates `maps.Equal`; `judgeBundleDigest` repeats an existing inline walk. — impact-negligible → report only
- F3 (nit). `judgeFailureResultWithStop` is one extra helper for one stop string (D-07). — impact-negligible → report only
- F4 (nit). `skills/do-work-toolbox/actions/help.md:18` lists `ai-report judge <dir>` without `--page` (excluded from Scope on purpose). — impact-negligible → report only
- Anti-bloat count: 1 production helper and 2 test helpers the REQ did not name (F2, F3); no new files, options, judge.json fields or decorative tests.

**Acceptance:** Pass. Six `TestAIReportJudge*` tests PASS with Chrome (no SKIP), second-opinion verdict tests ok, `go vet` clean for darwin and windows, CLI edge cases checked by the reviewer (escaped name, sub-folder page, directory named `.html`, wrong extension, empty value, no option).
**Restatement sweep:** not the wave's last integration; nothing redefined (additive option).
**Follow-ups created:** None (F1 fixed; F2-F4 report only)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: c9b1eb445acd98769aa0fb42d06caefe95dfd5ea
- Target: b7e6a08e20d22b31aaa3024e1ee287fe24ced170
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — report_judge.go, report_judge_test.go and report_second_opinion.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three Go files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three Go files matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target: b7e6a08e20d22b31aaa3024e1ee287fe24ced170; execution revision b7e6a08e20d22b31aaa3024e1ee287fe24ced170 (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set, no `HEAVY-RUN-LANE-SKIPPED`)
- do-work-cli-integrations: executed, exit 0, 56 s
- staged-skills: executed, exit 0, 32 s
- updater: executed, exit 0, 65 s
- installer: executed, exit 0, 28 s
- Earlier drain at the first merge 9618c836 (before the F1 fix): same four lanes executed, exit 0 (66 s, 34 s, 67 s, 27 s).

## Lessons Learned

**What worked:** The page check sits on the one entry path and runs before the browser lookup and before `--out` exists, so every refusal writes nothing; the refusal test proves it with a bundle digest. The expected `judgeUsage` conflict with REQ-729 resolved by keeping both texts.

**What didn't:** The first merge built `served_url` from the raw `--page` text instead of the checked path. Review finding F1 caught it: `./x.html` served as `/./x.html`, and a Windows `sub\x.html` reached the file server as `%5C` and could pass an error page. Fixed on the builder branch (D-11) and pinned by an exact-path assertion.

**Worth knowing:** After a path passes a containment check, use the cleaned form of that same path for everything downstream (URL, file open), not the raw input. The builder proposed no lesson bullet, so no satellite changed.

## Orientation

`ai-report-judge` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go`) renders `index.html` or the `--page` it is given, served in place from the bundle; the second-opinion finalize path still renders through its own copy-and-swap in `report_second_opinion.go`.

## Timing

Observed 2026-11-14T23:35:10Z to 2026-11-15T01:03:57Z: 1h 28m 47s total, 1h 43m 39s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 1h 13m 12s | 1 |
| verification-gate | 19m 00s | 4 |
| review | 10m 57s | 1 |
| handback-merge | 30s | 1 |

Slowest stage: builder-work / builder worktree build, 1h 13m 12s, outcome success.
