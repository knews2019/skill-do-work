# REQ-728 Hand-back (ai-report-judge renders one named page of a bundle)

Branch `worktree-agent-REQ-728-judge-renders-one-named-page`, tip `661361a4`, base `1de00505`.

## Commits
- `6e40562f` [REQ-728] RED: judge tests for --page rendering and page refusals
- `661361a4` [REQ-728] ai-report-judge --page renders one named page of a bundle in place

## Files
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (header comment, `judgeCodePage` constant, `judgeUsage`, `--page` parse case, page check after the bundle check, `judgeRender` page parameter and `ServedURL`, new `judgeFailureResultWithStop` and `judgeCheckPage`)
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (two tests plus two small digest helpers appended; `strings` import added)
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (line 1752 passes `"index.html"`)
- modified `skills/do-work-toolbox/actions/ai-report.md` (judge form line; one new paragraph in Step 7 after the command paragraph)
- modified `skills/do-work-toolbox/docs/ai-report-guide.md` (one sentence at the end of the judge paragraph; command list line)

## Scope changes
None. CHANGELOG.md, skills/do-work/CHANGELOG.md and VERSION are left to the integrator as instructed.

## Red-green evidence
- After commit 1 (`6e40562f`): `DO_WORK_HEAVY_TESTS=1 go test -count=1 -v -run 'TestAIReportJudgeRendersNamedPageInPlace|TestAIReportJudgeRefusesMissingOrOutsidePage'` → both FAIL with `TOOLBOX-USAGE ... unexpected argument "--page"; Usage: ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]` (all four refusal sub-tests: "no AI-REPORT-JUDGE-PAGE finding").
- After commit 2 (`661361a4`): both PASS, including the four sub-tests (missing, dot-dot escape, absolute path, symlink outside).
- GREEN probe at base `1de00505`: fails at its first check, because `git show 1de00505:<file> | grep -c -- --page` is 0 for both `ai-report.md` and `ai-report-guide.md` (and the two tests do not exist). At tip: `REQ-728 probe PASS`.

## Verification
1. `gofmt -l ./internal/toolboxcommands` → empty. `go vet ./internal/toolboxcommands/` → clean. `DO_WORK_HEAVY_TESTS=1 go test -count=1 -v -run '^TestAIReportJudge'` → six PASS lines (FailsOnPhoneWidthOverflow, FailsOnBrokenRelativeImage, PassesCleanBundleAndStopsServer, ReportsSkippedWithoutBrowser, RendersNamedPageInPlace, RefusesMissingOrOutsidePage), no SKIP. `go test -count=1 ./...` in do-work-cli → every package `ok` (run twice, both all ok; no per-file budget failure).
2. `bash _dev/tests/action-shell-blocks.sh` → "Shell-block lint passed: 88 fenced blocks and 35 shipped shell files; ShellCheck enabled." `bash _dev/tests/shipped-package-reference-contract.sh` → PASS.
3. Manual run: `ai-report-judge ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1 --page index-v1c01-codex-gpt-6.1-terra-custom-matrix-to-implementation.html --out <mktemp -d in scratch>` → exit 1, `served_url` `http://127.0.0.1:54254/index-v1c01-codex-gpt-6.1-terra-custom-matrix-to-implementation.html`, four captures plus judge.json; the one finding is the expected revise back-link `../2026-11-04_1524_proposal-tstack-review-workflows/index.html` (404, because only this folder is served). `git status --short -- ai-reports/` printed nothing.
4. GREEN probe from the worktree → `REQ-728 probe PASS` (rerun after the header-comment reflow amend: PASS, six judge PASS lines).

## P-A-U
- PLAN: add `--page` beside `--out`/`--browser`; check the page after the bundle resolves (IsLocal, extension, regular file, EvalSymlinks + judgePathInside); thread the page into `judgeRender` for `ServedURL`; one-argument change at the second-opinion caller; docs name the option and the refusal code.
- APPLY: done as planned in two commits; no other files.
- UNIFY: `git diff --stat 1de00505..HEAD` shows exactly 5 files (ai-report.md 4 lines, ai-report-guide.md 4, report_judge.go 47, report_judge_test.go +149, report_second_opinion.go 2); checked each diff for debug output, gofmt, and that the `--out` block (REQ-729's area) is untouched.

## Decisions
- D-06 DECIDE & STATE: the extension check (`.html`/`.htm`, case-insensitive) runs before the file stat, so `--page notes.txt` is refused as "not an HTML file" even when it exists. All refusals share the one code `AI-REPORT-JUDGE-PAGE`; the evidence text says which rule failed.
- D-07 DECIDE & STATE: the refusal uses a small sibling `judgeFailureResultWithStop` so the stop reason reads "--page names no HTML page inside the bundle" instead of "the render check could not run"; `judgeFailureResult` now delegates to it (same output for existing callers).
- D-08 DECIDE & STATE: the page check runs before the browser lookup and before `--out` is created, so a refusal writes nothing anywhere (the test asserts no judge.json and an unchanged bundle digest).
- D-09 DECIDE & STATE: the refusal table has the four cases the brief names; directory and wrong-extension refusals are covered by the code but not given their own sub-tests (YAGNI; same code path).

## Discovered Tasks
- None blocking. Possible follow-up already named in D-04 of the request: `secondOpinionRender` could switch from copy-and-swap to `--page` style in-place rendering; not done here by decision.

## Integration seams
`report_judge.go` at tip `661361a4`, line ranges after my edit:
- header comment `:24-28` (reworded and reflowed, now five lines).
- constant `judgeCodePage` `:35`; `judgeUsage` `:41` — I inserted ` [--page <file.html>]` after `<bundle-dir>`. REQ-729 appended a sentence at the END of the same string, so the merge is a one-line conflict: keep my insertion and REQ-729's appended sentence.
- `handleAIReportJudge` `:90-163`: new `page := "index.html"` (`:92`), `--page` case in the parse loop (`:99-100`), page check `:123-125` right after `resolvedBundle` and before the browser lookup, and the `judgeRender(&report, bundle, page, output, browser)` call. The `--out` resolution block (REQ-729) is byte-unchanged; it only moved down by five lines.
- `judgeRender` `:166-239`: signature gains `page`, `ServedURL` built from it (`:177`).
- new functions `judgeFailureResultWithStop` `:292-296` and `judgeCheckPage` `:300-317`; `judgeFailureResult` `:288-290` now delegates.
- `report_judge_test.go`: my additions are appended after `TestAIReportJudgeReportsSkippedWithoutBrowser`, plus a `"strings"` import; REQ-729 also appends at the end, so keep both blocks (and dedupe the import if REQ-729 added `strings` too).
- `ai-report.md`: the judge form line (`:36`) and a new standalone paragraph after Step 7's command paragraph; no neighbour text rewritten. `ai-report-guide.md`: one sentence appended to the judge paragraph and the `judge <dir>` command list line.

## Proposed CHANGELOG entry
Title: **Render Check Takes One Named Page of a Bundle** (minor bump: new option)

What changes for an installed copy: `ai-report-judge --page <file.html>` renders one page of a bundle in place, so checking a second-opinion page no longer needs a bundle copy; a missing or outside page is refused with `AI-REPORT-JUDGE-PAGE`; runs without the option are unchanged.

## Proposed lesson
none

## Lessons read
- `_dev/primes/lessons-releases.md` (full; its families do not apply: no file move, no manifest edit).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` by fixed-string grep: `[family: guard-on-one-entry-path]` (the page guard sits in the one entry point all `--page` values pass through), `[family: existence-is-not-identity]` (existence alone is not enough: the check also requires a regular file, an HTML extension and a resolved path inside the bundle), `[family: silent-skip-reads-as-red]` (the heavy test was run with `DO_WORK_HEAVY_TESTS=1` and the probe requires PASS, not SKIP).
