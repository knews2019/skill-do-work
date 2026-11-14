---
id: REQ-728
title: 'ai-report-judge renders one named page of a bundle'
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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

