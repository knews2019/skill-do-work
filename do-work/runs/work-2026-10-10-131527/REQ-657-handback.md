# Hand-back: REQ-657 (ai-report judge runs the render check as one bundled command)

- Branch: `worktree-agent-REQ-657-ai-report-judge-render-check`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-657-ai-report-judge-render-check`
- Base commit: `bd56c4b0`
- Commits:
  - `5b8d0143` [REQ-657] add ai-report-judge render check verb (Go verb, tests, registration)
  - `7492905f` [REQ-657] point render-check prose at ai-report-judge (seven Markdown files)
- Nothing under `do-work/` was staged or committed. No CHANGELOG, VERSION or mirror was touched.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (new): the `ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]` verb. Loopback `http.FileServer` on `127.0.0.1:0`, browser discovery, a small production DevTools pipe transport (fd 3/4, NUL-terminated JSON), four passes (wide 1440x1000 and phone 390x844, light and dark), the same-origin link check, `judge.json` and the exact text output.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (new): the four named tests.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modified): `CommandAIReportJudge = "ai-report-judge"` as the last const line and the last `Handlers()` entry.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modified): the registration list gains the constant; count 7 to 8.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): new Input sentence for `judge <bundle-dir>` after the `:30` line; Step 7 runs the command in a bash fence and keeps the judgement list and the two-pass SVG rule; Step 8 reads the `judge.json` verdict and removes the judge output directory; checklist line updated. Philosophy line unchanged.
- `skills/do-work-toolbox/actions/stakeholder-report.md` (modified): Step 5 calls the same command; one pass of the four captures.
- `skills/do-work-toolbox/actions/architecture-report.md` (modified): Step 5 runs the command on the draft's directory; the "network access disabled" sentence and the "never claim an unperformed check passed" clause stay.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): the render paragraph describes the command and its outputs; one `do-work-toolbox ai-report judge <dir>` line at the end of the Input block.
- `skills/do-work-toolbox/SKILL.md` (modified): new routing row directly below the `ai-report` row.
- `skills/do-work-toolbox/actions/help.md` (modified): new menu line directly below the `ai-report` line.
- `skills/do-work/docs/command-line-guide.md` (modified): the Toolbox line names `ai-report-judge` (see D-09).

## P-A-U

- **[PLAN]:** Follow pre-dispatch D-01 to D-08. Tests first (four named tests plus the registration count), RED as a compile failure. Then one file `report_judge.go`: argument parsing with `optionValue`; bundle and `--out` resolution against `RepositoryRoot`; `--out` refused inside the bundle after `EvalSymlinks`; `--browser` resolved or usage failure; discovery from two package-level lists; skip path writes `judge.json` and returns OutcomeFindings; `judgeRender` owns server and engine with a `defer` right after each starts; four passes via `Emulation.*`, `Page.navigate`, a readiness poll, measurement and `Page.captureScreenshot`; links collected in the wide light pass and fetched from our own server only. Then prose in the seven Markdown files. Lessons read: `lessons-releases.md` (whole), families `reaped-by-its-own-parent` and `silent-skip-reads-as-red` (targeted grep). The reaped-by-its-own-parent lesson shaped `stop()`: close the command pipe first so the engine exits on EOF and reaps its own helpers, and kill the leader only after the 30 s deadline. The silent-skip lesson shaped the skip path: `skipped` is OutcomeFindings (exit 1), announced as its own code `AI-REPORT-JUDGE-SKIPPED`, and never `pass`.
- **[APPLY]:** Done as planned in two commits. Deviations are listed under Decisions (D-09 to D-14).
- **[UNIFY]:** `git diff bd56c4b0 --stat`:

```
 skills/do-work-toolbox/SKILL.md                    |   1 +
 skills/do-work-toolbox/actions/ai-report.md        |  16 +-
 .../do-work-toolbox/actions/architecture-report.md |   8 +-
 skills/do-work-toolbox/actions/help.md             |   1 +
 .../do-work-toolbox/actions/stakeholder-report.md  |   8 +-
 skills/do-work-toolbox/docs/ai-report-guide.md     |   3 +-
 skills/do-work/docs/command-line-guide.md          |   2 +-
 .../internal/toolboxcommands/commands.go           |   2 +
 .../internal/toolboxcommands/commands_test.go      |   6 +-
 .../internal/toolboxcommands/report_judge.go       | 570 +++++++++++++++++++++
 .../internal/toolboxcommands/report_judge_test.go  | 167 ++++++
 11 files changed, 773 insertions(+), 11 deletions(-)
```

Exactly the eleven Scope files. Checks (run from the worktree root):

| Check | Exit | Wall time |
|---|---|---|
| `bash .../REQ-657-probe.sh` | 0, last line `REQ-657 probe: GREEN` | 3.4 s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal/toolboxcommands` | 0, no output | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands/` | 0 | <2 s |
| `DO_WORK_HEAVY_TESTS=1 go test ... -v -run '^TestAIReportJudge' ./internal/toolboxcommands/` | 0, four PASS, no SKIP | 2.8 s |
| `go test ... -short -run '^(TestAIReportJudge.*\|TestHandlersRegisterCanonicalToolboxCommands)$'` (fast-gate shape) | 0; three browser cases SKIP, skip case and registration PASS | 0.4 s |
| `go test ... -short ./internal/toolboxcommands/` (whole package, short) | 0 | 4.1 s |
| `shellcheck -e SC2154` on the three new bash fences (placeholders replaced by variables) | 0 | <1 s |
| `bash _dev/tests/shipped-package-reference-contract.sh` (citation check of the edited Markdown) | 0, `PASS` | 1 s |
| `git diff bd56c4b0 --check` | 0 | <1 s |

Files checked by reading the diff: all eleven above. No debug artifacts. No browser test needed a rerun under load.

## Proof record

1. **RED, command** (base `bd56c4b0`, before any edit): `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format json ai-report-judge <fixture>/overflow` and `.../broken` each exited **2** with finding `UNKNOWN-COMMAND`: `command "ai-report-judge" is not available: available commands are ...`.
2. **RED, tests**: `DO_WORK_HEAVY_TESTS=1 go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^TestAIReportJudge' ./internal/toolboxcommands/` failed to build: `undefined: discoverJudgeBrowser`, `undefined: judgeReport`, `undefined: handleAIReportJudge`, `undefined: judgeCodeOverflow`, `undefined: CommandAIReportJudge` ... `FAIL [build failed]`.
3. **GREEN, tests** (same command):

```
--- PASS: TestAIReportJudgeFailsOnPhoneWidthOverflow (0.75s)
--- PASS: TestAIReportJudgeFailsOnBrokenRelativeImage (0.69s)
--- PASS: TestAIReportJudgePassesCleanBundleAndStopsServer (0.62s)
--- PASS: TestAIReportJudgeReportsSkippedWithoutBrowser (0.00s)
ok  	.../internal/toolboxcommands	2.376s
```

No `--- SKIP` line. Discovery found `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome` with no flag.

**GREEN, command** (rerun of step 1's two fixtures, text format, `--out` in a scratch dir): each exited **1**.
- Overflow `judge.json` findings:
  `{"kind": "AI-REPORT-JUDGE-OVERFLOW", "viewport": "phone", "scheme": "light", "scroll_width": 600, "inner_width": 390}`
  `{"kind": "AI-REPORT-JUDGE-OVERFLOW", "viewport": "phone", "scheme": "dark", "scroll_width": 600, "inner_width": 390}`
  (wide captures: scroll width 1440, inner width 1440, no finding.)
- Broken image `judge.json` finding:
  `{"kind": "AI-REPORT-JUDGE-BROKEN-LINK", "attribute": "src", "target": "missing.png", "status": 404}`
- Clean fixture (`#x`, `https://example.invalid/` and `mailto:` links, dark-mode CSS): exit **0**, verdict `pass`, four PNGs; `phone-dark.png` inspected and shows the dark palette. After the run: `lsof -nP -iTCP -sTCP:LISTEN` had no listener on the served ports (51456, 51477, 51500); `pgrep -fl -- --remote-debugging-pipe` and `pgrep -fl ai-report-judge-profile` printed nothing; no `ai-report-judge-profile-*` directory remained.
- Extra manual cases: `--out` inside the bundle exits 2 (`TOOLBOX-USAGE ... --out must be outside the bundle`) and creates nothing in the bundle; `--browser /nope/chrome` exits 2 (`--browser names no runnable browser`); a directory without `index.html` exits 2; no arguments exits 2 with the usage line; `--browser /bin/echo` gives verdict `error`, exit 2, finding `AI-REPORT-JUDGE-ERROR: no reply to Target.getTargets within 30s: EOF` in 1.3 s, and `judge.json` is still written; `--format json` shows the `judge.json` path in `changes[0].path`.
4. **Mutation**: comparison changed to `capture.ScrollWidth > capture.InnerWidth+1000`; `TestAIReportJudgeFailsOnPhoneWidthOverflow` failed with `report_judge_test.go:71: overflowing bundle: outcome success verdict "pass" findings []`. Reverted (byte copy restored, `grep -c "InnerWidth+1000"` = 0) before the commit.

## Decisions

- **D-09** (DECIDE & STATE) `skills/do-work/docs/command-line-guide.md:26` lists Just recipes ("`just --list` is the live inventory installed from the managed template"), and the brief forbids a Just recipe. Appending `ai-report-judge` bare would claim a recipe that does not exist. I appended one sentence to the same line instead: "The toolbox CLI verb `ai-report-judge` has no recipe; actions call it through `do-work-cli.sh`." The probe's `grep -F ai-report-judge` passes. Reversible; if REQ-655 appends a recipe-less verb to the same line, the integrator can merge both into this sentence.
- **D-10** (DECIDE & STATE) The architecture-report draft is "one self-contained HTML document ... outside the reports directory" (`architecture-report.md:79`), a file, while the verb serves a directory holding `index.html`. The Step 5 text tells the agent to keep the draft as `index.html` in a directory of its own and pass that directory. `:79` itself is not edited. The publish command already takes `<draft-path>`, which is that `index.html`.
- **D-11** (DECIDE & STATE) Launch flags add `--hide-scrollbars` and `--no-first-run` to the board's set. `--hide-scrollbars` keeps `innerWidth` equal to the layout width, so a classic scrollbar cannot hide a small overflow under `scrollWidth > innerWidth`. The engine starts on `about:blank` and attaches to the first page target, and each pass navigates with `Page.navigate`.
- **D-12** (DECIDE & STATE) Readiness after `Page.navigate`: the same URL is loaded four times, so `readyState === "complete"` alone could read the previous document. Before each navigation the judge sets `window.aiReportJudgeStale = true` on the old window, then polls `!window.aiReportJudgeStale && document.readyState === "complete"` (evaluate errors during the swap count as not ready, bounded by 30 s). The marker lives only in the judge's own engine; the bundle is never modified.
- **D-13** (DECIDE & STATE) `--format json` hides the exact text output, so the `judge.json` path is reported as a `changes` entry (`kind: created`). The ai-report Step 7 prose tells the reader to find it there.
- **D-14** (DECIDE & STATE) `--browser` and the bundle check are validated before the default temporary output directory is created, so a usage failure leaves no empty temp directory (found in a manual run and fixed before the commit). A broken link's `target` is the attribute value as written (`missing.png`), not the resolved URL, so the finding names what the author must fix.

## Discovered Tasks

- `skills/do-work/docs/command-line-guide.md` mixes two inventories: the Just-recipe list and, now, a CLI verb without a recipe. If more recipe-less toolbox verbs arrive (REQ-655 is one candidate), the line may need a separate "CLI verbs without recipes" bullet. impact-low → report only
- `skills/do-work-toolbox/actions/architecture-report.md:79` still says "one self-contained HTML document as a draft" without naming the file; D-10 adds the `index.html`-in-its-own-directory rule only at Step 5. A later edit could move that rule to `:79`. impact-low → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, targeted: `[family: reaped-by-its-own-parent]` (REQ-543 bullet) and `[family: silent-skip-reads-as-red]` (REQ-634 and REQ-566 bullets).
- Dropped for budget per the REQ: `lessons-shell-commands.md`, `lessons-do-kanban.md`, the full `lessons-do-work-cli.md`, `lessons-action-files.md`. The primes `prime-shell-commands.md`, `prime-action-files.md` and `prime-releases.md` were read in full.

## Anti-bloat check

`git diff --stat`: see [UNIFY]. Names added that the brief did not name:

- `judgeCodeError` / `AI-REPORT-JUDGE-ERROR`: D-05 has an `error` verdict, and its finding needs a code so the cause shows in `judge.json` and the CLI result.
- `judgeProtocolDeadline`, `judgePollInterval`, `judgeUsage` constants: one deadline for every wait (brief: 30 s), one poll interval, one usage string reused by three refusals.
- `judgeViewports` variable: the two viewport sizes in one place, iterated by the pass loop.
- `judgeCapture`, `judgeFinding`, `judgeReport` types: the `judge.json` shape D-06 names.
- `judgeRender`: owns the server and engine so that their `defer`s fire on every return path.
- `judgeResult`: writes `judge.json` and maps verdict to outcome and findings; one place for all four verdicts.
- `judgeErrorResult`: the `error` verdict path, still writing `judge.json`.
- `judgeFailureResult`: OutcomeFailure with code `AI-REPORT-JUDGE-ERROR` for the two cases without a usable output directory (temp dir creation fails, `judge.json` write fails). `usageResult` hard-codes `TOOLBOX-USAGE`, which would mislabel these.
- `judgeAbsolutePath`: resolves the bundle and `--out` against `RepositoryRoot`.
- `judgeResolveExistingPrefix`, `judgePathInside`: the D-06 "refuse `--out` inside the bundle" check for an `--out` that does not exist yet (EvalSymlinks fails on a missing path).
- `judgeLookupBrowser`, `discoverJudgeBrowser`: D-03 explicit lookup and discovery.
- `judgeBrowserSession` with `startJudgeBrowser`, `stop`, `call`, `attachToPage`, `evaluate`, `renderPass`: the D-02 transport the brief asked to copy in small form.
- Test helpers `requireJudgeBrowser` (the heavy gate plus the no-engine skip, shared by three tests) and `runJudgeOnBundle` (fixture, handler call, `judge.json` decode, shared by all four).
- Launch flags `--hide-scrollbars`, `--no-first-run` (D-11).

No new flag, environment variable, Just recipe, Node file, `scripts/` directory or config.

## Proposed CHANGELOG entry

**ai-report Render Check Command**

The report render check is now one shipped command instead of steps each session retyped. It catches a report that scrolls sideways on a phone or links to a missing image before a stakeholder sees it.

- New do-work-cli verb `ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]`: serves the bundle on a free local port, captures 1440x1000 and 390x844 in light and dark, fails on horizontal overflow and on any same-origin `href`/`src` answering 400 or more, writes `judge.json` and four PNGs outside the bundle, and stops its server and browser on every exit path.
- No browser found gives verdict `skipped` with exit 1. A skip is never a pass.
- `ai-report` Step 7, `stakeholder-report` Step 5 and `architecture-report` Step 5 call the command; `ai-report judge <dir>` is routed, listed in help and documented in the guide.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

- [family: assertion-passes-when-feature-is-dead] [REQ-657: a readiness wait on `document.readyState === "complete"` passes immediately when the same URL is navigated again, because the old document is still complete; mark the old window before navigating and wait for the mark to disappear, or every pass after the first measures the previous render](../../../../do-work/archive/UR-144/REQ-657-ai-report-judge-render-check-command.md)

## Integration seams seen

- `commands.go` / `commands_test.go`: my constant and map entry are the last lines; count is 8. With REQ-655 the integrator keeps both entries and sets the count to 9.
- `SKILL.md` routing, `help.md` menu, `ai-report-guide.md` Input block and the `ai-report.md` Input section: my lines are new lines; existing lines are unchanged. The `ai-report.md` `:109` sentence (REQ-655's area) is untouched.
- `command-line-guide.md:26`: my sentence is appended to the line (D-09). If REQ-655 appends a verb to the same line, keep both.
- Test wall times: heavy run 2.4 to 3.1 s for the package (each browser case 0.6 to 2.1 s); fast-gate shape 0.2 s; whole package with `-short` 3.9 s. All far below the 30 s per-file budget.
