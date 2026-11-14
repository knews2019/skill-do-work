# Builder Brief — REQ-728 (ai-report-judge renders one named page of a bundle)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-728-judge-renders-one-named-page`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-728-judge-renders-one-named-page`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-05, Exploration and Scope: read `do-work/working/REQ-728-judge-renders-one-named-page.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-728-handback.md`
- Up to ten builders run at the same time in this run, each in its own worktree from the same main HEAD. Integration is serial. REQ-729 (relative `--out` rule of the same command) edits the same Go file in parallel; see Integration seams.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `_dev/primes/prime-action-files.md` (you edit an action file) and `_dev/primes/prime-releases.md`.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. Required lessons: `_dev/primes/lessons-releases.md`. Dropped only for budget, read the matching bullets with a fixed-string grep in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: `[family: guard-on-one-entry-path]`, `[family: existence-is-not-identity]`, `[family: silent-skip-reads-as-red]`.

## Decided (do not reopen)
- **D-01:** the option is `--page <name>` (also `--page=<name>`), parsed by `optionValue` beside `--out` and `--browser`.
- **D-02:** `--page` takes a slash-separated path relative to the bundle naming a regular `.html`/`.htm` file inside it. Refuse an absolute path, a `..` escape, a symlink whose target resolves outside the bundle, a directory, a missing file, or another extension with one new finding code `AI-REPORT-JUDGE-PAGE` (OutcomeFailure, so exit 2), before any browser lookup and before `--out` is created; no judge.json is written for a refusal, like today's usage refusals.
- **D-03:** the bundle check stays as it is: the bundle must hold `index.html` in both forms.
- **D-04:** `secondOpinionRender` and the finalize path stay as they are. Its one call to `judgeRender` (`report_second_opinion.go:1752`) only gains the `"index.html"` argument.
- **D-05:** no new judge.json field and no schema bump. With `--page`, `served_url` ends with the named page (URL-escaped path); `bundle` stays the bundle folder; the four capture names stay `wide-light.png`, `wide-dark.png`, `phone-light.png`, `phone-dark.png`.

## What to build
Two commits.

**Commit 1 (RED, tests only).** Append two test functions at the END of `report_judge_test.go`:
- `TestAIReportJudgeRendersNamedPageInPlace` (heavy: call `requireJudgeBrowser`). Fixture bundle with two pages: a clean `index.html` and an `index-v1c01-x.html` that overflows at phone width (for example a 600px-wide block, as the existing overflow test does). Run the handler with `--page index-v1c01-x.html --out <outside dir>`. Assert: `served_url` ends with `/index-v1c01-x.html`; the verdict is `fail` with a phone overflow finding (index.html alone would pass, so this proves the captures show the named page); the four captures exist; a digest (path plus bytes of every file) of the bundle is identical before and after. Name in the test comment the failure it pins.
- `TestAIReportJudgeRefusesMissingOrOutsidePage` (light, no browser; table of sub-tests): a missing name, `../outside.html` (the file exists beside the bundle), an absolute path, and a symlink inside the bundle pointing at a file outside it. Each must return OutcomeFailure with a finding code `AI-REPORT-JUDGE-PAGE`, write no judge.json, and leave the bundle digest unchanged. `runJudgeOnBundle` always reads judge.json, so call `handleAIReportJudge` directly here.
Run both and record the RED output (an unknown-argument usage refusal).

**Commit 2 (the option).**
1. `report_judge.go`: add the `--page` case to the parse loop (`:93-102`) and `[--page <file.html>]` to `judgeUsage` (`:39`). After the bundle is resolved (`:110-117`), check the page per D-02 (`filepath.IsLocal` on `filepath.FromSlash(page)`, then `os.Stat` regular file, then `filepath.EvalSymlinks` and `judgePathInside(resolvedPage, resolvedBundle)`). Add a constant `judgeCodePage = "AI-REPORT-JUDGE-PAGE"` beside the others (`:30-33`) and refuse through `judgeFailureResult` (`:279`) or a small sibling with a clearer stop reason. Default page is `index.html`.
2. `judgeRender` (`:159`) gains a page parameter; `report.ServedURL` (`:168`) becomes the origin plus the escaped page path. Pass the page from `handleAIReportJudge` (`:147`) and `"index.html"` from `report_second_opinion.go:1752`.
3. Update the header comment (`:24-27`) so it no longer says the judge renders only index.html.
4. `skills/do-work-toolbox/actions/ai-report.md`: the `judge <bundle-dir>` form line (`:36`) and Step 7 (`:181-184`) say that `--page <file.html>` renders that page of the bundle in place (a second-opinion page such as `index-v1c01-<name>.html`), with no bundle copy, and name the refusal code.
5. `skills/do-work-toolbox/docs/ai-report-guide.md`: one sentence in the judge paragraph (`:59`) and the command list line (`:80`).

**Not yours:** the `--out` resolution at `report_judge.go:127-138` (REQ-729); `report_index.go` (REQ-730); `VERSION`, both `CHANGELOG.md` files and every `do-work/` path (the integrator writes the release); `actions/help.md`, `stakeholder-report.md`, `architecture-report.md`; any `ai-reports/` bundle.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-728]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/toolboxcommands` prints nothing; `go vet ./internal/toolboxcommands/`; `DO_WORK_HEAVY_TESTS=1 go test -count=1 -v -run '^TestAIReportJudge' ./internal/toolboxcommands/` shows six PASS lines and no SKIP; `go test -count=1 ./...` passes (the REQ's GREEN condition).
2. `bash _dev/tests/action-shell-blocks.sh` (you may change a fenced command in an action), and `bash _dev/tests/shipped-package-reference-contract.sh`.
3. A manual run on a real second-opinion bundle, read-only: pick one `ai-reports/` bundle holding an `index-v1c*.html` page, run `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge <bundle> --page <that page> --out <a mktemp -d folder>`, confirm `served_url` names the page, and confirm `git status --short -- ai-reports/` prints nothing.
4. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-728-probe.sh` run with your worktree as the current directory prints `REQ-728 probe PASS`. It fails at base by design (the docs do not name `--page` and the two new tests do not exist).

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the two new tests failing after commit 1 and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-06 onward (D-01 to D-05 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — see below.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Render Check Takes One Named Page of a Bundle") and the line "What changes for an installed copy": `ai-report-judge --page <file.html>` renders one page of a bundle in place, so checking a second-opinion page no longer needs a bundle copy; a missing or outside page is refused with `AI-REPORT-JUDGE-PAGE`; runs without the option are unchanged. Expect a minor version bump (new option).
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
| Shared file | REQ-728 edits | Sibling that edits it in parallel |
| --- | --- | --- |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` | `judgeUsage` (`:39`), the `handleAIReportJudge` parse loop and the new page check after the bundle check, `judgeRender` signature and `ServedURL`, the header comment, one new constant | REQ-729: `judgeUsage` and the `--out` resolution in `handleAIReportJudge` (`:127-138`). Both edit `judgeUsage`, so the second integrator resolves a one-line conflict there. |
| `report_judge_test.go` | two new functions at the end | REQ-729 also appends tests at the end |
| `skills/do-work-toolbox/actions/ai-report.md` | the judge form line (`:36`) and Step 7 (`:181-184`) | REQ-729 (Step 7 `--out` rule), REQ-730 and REQ-731 (catalog prose) |
| `skills/do-work-toolbox/docs/ai-report-guide.md` | judge paragraph (`:59`) and command list (`:80`) | REQ-729, REQ-730, REQ-731 |

In the seam section, name which functions of `report_judge.go` you changed, with line ranges after your edit, so the second integrator can resolve the overlap with REQ-729. Add new text as its own sentence or bullet; never rewrite a neighbour's paragraph.

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
