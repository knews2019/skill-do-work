# Review: REQ-728 (ai-report-judge renders one named page of a bundle)

**Approve.** `--page` renders the named page in place, refuses bad pages before anything is written or started, and leaves the no-option run and the second-opinion finalize path unchanged.
Route B | merge `9618c836` (builder commits `6e40562f` RED tests, `661361a4` implementation), range `c9b1eb44..9618c836`

## What's built

- `ai-report-judge <bundle> --page <file.html>` (also `--page=<file.html>`) serves the bundle directory as before and opens the named page instead of `index.html`. `served_url` ends with the URL-escaped page.
- A page that is absolute, escapes with `..`, is not `.html`/`.htm`, is not a regular file, or resolves outside the bundle through a symlink is refused with `AI-REPORT-JUDGE-PAGE` (exit 2). The check runs before the browser lookup and before `--out` is created.
- `ai-report.md` (form line, one new Step 7 paragraph) and `ai-report-guide.md` (one sentence, command list) name the option. Release files are not in this range; the finalization writes them.

## Decisions / risks for you

None. D-04 (the second-opinion finalize path keeps its copy-and-swap) is a sound choice and is already recorded as a discovered task.

## Findings

**Important:** None.

**Minor:**
- F1. `report_judge.go:177` builds `ServedURL` from the raw `--page` text, not from the checked and cleaned path. On Windows, `--page sub\p.html` passes `judgeCheckPage` (it uses `filepath.FromSlash`, `IsLocal`, `Stat`), but the URL carries `%5C`; `http.Dir.Open` rejects the backslash through `filepath.Localize`, Chrome shows the server's error text, and the judge does not check the top-level response status, so the run can report `pass` on an error page. On macOS/Linux the only effect is cosmetic: `--page ./sub/../index.html` passed and recorded `served_url` `/./sub/../index.html`. The fix is one line: escape `filepath.ToSlash(localPage)` after `filepath.Clean`. — impact-user-visible → report only

**Nit:**
- F2. In `report_judge_test.go`, `sameJudgeDigest` does what the standard library's `maps.Equal` does, and `judgeBundleDigest` repeats the walk that `TestAIReportJudgePassesCleanBundleAndStopsServer` writes inline. — impact-negligible → report only
- F3. `judgeFailureResultWithStop` (D-07) exists only so that one caller gets a different stop text. One extra helper for one string. — impact-negligible → report only
- F4. `skills/do-work-toolbox/actions/help.md:18` still lists `ai-report judge <dir>` without `--page`. The Scope excluded `help.md` on purpose, and the line is a short summary, so this is cosmetic. — impact-negligible → report only

**Anti-bloat count (items the REQ did not name):** 1 production helper (`judgeFailureResultWithStop`, F3), 2 test helpers (`judgeBundleDigest`, `sameJudgeDigest`, F2). No new files, no options other than `--page`, no new judge.json fields. Both new tests pin a real failure. No decorative tests.

## Requirements Checklist

- [x] Option renders the named page without a bundle copy and without writing inside it: delivered (`judgeRender` serves `http.Dir(bundle)`; the test compares a full bundle digest before and after).
- [x] Served address ends with the named page and the captures show it: delivered (suffix assertion plus a 600 px phone overflow that only the named page has).
- [x] Missing or outside page refused with a typed finding: delivered (`AI-REPORT-JUDGE-PAGE`; missing, `..`, absolute, symlink, plus directory and wrong extension checked by hand).
- [x] No-option behavior unchanged (captures, judge.json, exit codes): delivered (`page` defaults to `index.html`, so `ServedURL` is `origin + "/index.html"`; `judgeFailureResult` output is unchanged; the four older judge tests pass).
- [x] Judge never writes inside the bundle: delivered (no write path added; refusals create no `--out`).
- [x] Second-opinion finalize path keeps its verdicts: delivered (`secondOpinionRender` passes `"index.html"`, the page its copy already writes; second-opinion tests pass).
- [x] Docs: judge form line, Step 7, guide: delivered. The new Step 7 paragraph sits after the render paragraph and before the "Judge width usage" paragraph, in the right section. It says `--page` is relative to the bundle; REQ-729's sentence names only `<bundle-dir>` and `--out <dir>` as relative to `--repo-root`. The two do not contradict. The merged `judgeUsage` keeps both texts (D-10).
- [ ] `CHANGELOG.md`, its mirror and `VERSION`: N/A for this range (Scope marks them integrator-only; the finalization release writes them).

## Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `DO_WORK_HEAVY_TESTS=1 go test -count=1 -v -run '^TestAIReportJudge' ./internal/toolboxcommands/` at `9618c836`: all 6 tests PASS, no SKIP (Google Chrome).
- `go test -count=1 -run SecondOpinion ./internal/toolboxcommands/`: ok (15.4 s).
- `go vet ./internal/toolboxcommands/` clean, `gofmt -l` empty, `GOOS=windows GOARCH=amd64 go vet ./internal/toolboxcommands` clean.
- Built CLI, run on a scratch bundle:
  - `--page 'my page #1.html'`: exit 1, `served_url` `/my%20page%20%231.html`, the page rendered and its missing image was reported. Escaping works.
  - `--page=sub/p.html`: exit 1, `served_url` `/sub/p.html`, its `../missing.png` was resolved against the sub-folder and reported.
  - `--page dir.html` (a directory): exit 2, `AI-REPORT-JUDGE-PAGE`, the `--out` directory was not created.
  - `--page notes.txt`: exit 2, `AI-REPORT-JUDGE-PAGE`.
  - `--page=` (empty): exit 2, `TOOLBOX-USAGE`.
  - no option: exit 0, `served_url` `/index.html`.
  - `--page ./sub/../index.html`: exit 0, cosmetic `served_url` (F1).
  - The bundle listing after all runs matched the files written before them.

## Suggested Additional Testing

1. Live acceptance: unassessed. Run `ai-report judge <real bundle> --page index-v1c01-<name>.html` on a real second-opinion bundle and confirm the four captures show that page (the builder reports one such manual run).
2. Deployment: unassessed. After release, confirm the installed consumer CLI accepts `--page`.
3. Windows: run `--page sub\p.html` on Windows to confirm F1, or fix F1 first.

## Scores (on the record, not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every requirement delivered; release files belong to finalization |
| Code Quality | 92% | Clear guard and threading; served URL uses the raw page text (F1) |
| Test Adequacy | 90% | Both tests pin real failures; no test for `--page=` form or an escaped name |
| Scope | 95% | Five declared files only; one small undeclared helper (F3, recorded as D-07) |
| Risk | Low | Windows-only false pass with backslash input (F1) |
| Acceptance | Pass | Implementation and integration stages |

## Follow-ups created

None (4 findings report only)

## Review

**Overall: 94%** | 2026-11-15T00:55:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `report_judge.go:177` builds `ServedURL` from the raw `--page` text, so on Windows `--page sub\p.html` passes the check but is served as `%5C`, `http.Dir` refuses it, and the judge can pass an error page. On Unix only `served_url` is cosmetic (`/./sub/../index.html`). Fix: escape `filepath.ToSlash` of the cleaned local path. — impact-user-visible → report only
- F2 (nit). `sameJudgeDigest` duplicates `maps.Equal`; `judgeBundleDigest` repeats an existing inline walk. — impact-negligible → report only
- F3 (nit). `judgeFailureResultWithStop` is one extra helper for one stop string (D-07). — impact-negligible → report only
- F4 (nit). `actions/help.md:18` lists `ai-report judge <dir>` without `--page` (excluded from Scope on purpose). — impact-negligible → report only
- Anti-bloat count: 1 production helper and 2 test helpers the REQ did not name (F2, F3); no new files, options, fields or decorative tests.

**Acceptance:** Pass. Implementation and integration stages: 6 judge tests PASS with Chrome, second-opinion tests ok, vet (darwin and windows) and gofmt clean, CLI edge cases checked (escaped name, sub-folder page, directory, wrong extension, empty value, no option).
**Restatement sweep:** nothing redefined (additive `--page` option; existing judge contract unchanged). Not the wave end, so no inherited sweep.
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Delta re-review (F1 fix)

**Approve. Overall: 95%** | 2026-11-15T01:04:06Z | delta `9618c836..b7e6a08e` (builder commit `1e80dacd`), 2 files, +10/-6

- F1 is fixed. `report_judge.go` runs `page = filepath.ToSlash(filepath.Clean(filepath.FromSlash(page)))` only after `judgeCheckPage` accepts the page, so the served path is the checked path. The default `index.html` stays `index.html`, so the no-option `served_url` stays `/index.html`. A Windows `sub\x.html` becomes `sub/x.html`, and `./x.html` becomes `x.html`. On Unix a literal backslash in a file name is left as it is, which is correct there. Nothing else in the handler changed.
- The test now passes `./index-v1c01-x.html` and requires the parsed `served_url` path to be exactly `/index-v1c01-x.html`. That pins the cosmetic half of F1. The Windows half cannot run on macOS and is covered by reading the code. Removing the `strings` import is correct, because the only use was replaced.
- Evidence at `b7e6a08e`: all 6 `TestAIReportJudge*` tests PASS with `DO_WORK_HEAVY_TESTS=1` and none skip. `GOOS=windows GOARCH=amd64 go vet ./internal/toolboxcommands` is clean.
- New findings: none. F2, F3 and F4 are unchanged and stay report only.
