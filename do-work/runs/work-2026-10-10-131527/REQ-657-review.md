## Review: REQ-657 (ai-report judge runs the render check as one bundled command)

**Approve with follow-ups**: the `ai-report-judge` verb works as specified, but the architecture-report call site can never reach `pass` on a report that links to its prior bundle.
Route B | merge `3f8fd77c` (range `aea1c4d3..3f8fd77c`)

### What's built
- A do-work-cli verb `ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]`. It serves the bundle on `127.0.0.1:0`, renders `index.html` at 1440x1000 and 390x844 in light and dark through a headless Chromium-family engine, fails on `scrollWidth > innerWidth` and on same-origin `href`/`src` answering 400 or more, writes `judge.json` plus four PNGs outside the bundle, and stops the server and the engine on every return path. Verdicts: `pass` exit 0, `fail` and `skipped` exit 1, `error` exit 2.
- `ai-report.md` Step 7 and Step 8, `stakeholder-report.md` Step 5 and `architecture-report.md` Step 5 call the verb. The toolbox routing table, help menu, ai-report guide and command-line guide name it.
- Still missing: architecture-report's relative link to the prior bundle always fails the check (F1).

### Decisions / risks for you
- F1 needs a choice. Smallest change (recommended): a prose edit in `architecture-report.md` Step 5 that says the judge serves only the draft directory, so a `BROKEN-LINK` finding whose target starts with `../` is expected there, and that link is checked after publication in Step 6 (which already says "Verify the final HTML from that location, including its relative links"). Value: one sentence, no code. Risk: the agent must tell expected `../` findings from real ones. Alternative: a code change that reports targets escaping the bundle root as their own finding kind that does not fail the verdict. Value: the verdict stays mechanical. Risk: a new finding kind and test, for one consumer.

### Findings

**Important:**
- F1 `skills/do-work-toolbox/actions/architecture-report.md:93,103`. Step 4 requires "Relative links to the prior bundle must work from the final bundle location", but the judge serves only the draft directory, so `../<prior>/index.html` resolves to `/<prior>/index.html` on the judge's server and answers 404. Reproduced: a bundle `reports/new` linking to an existing sibling `reports/prior/index.html` gave `verdict=fail`, `AI-REPORT-JUDGE-BROKEN-LINK href="../prior/index.html" answered HTTP 404`, exit 1. Step 5 says "fix and rerun until the verdict is `pass`", which is unreachable on every non-first architecture report and invites deleting a link Step 4 requires. — impact-user-visible → report only

**Minor:**
- F2 `skills/do-work-toolbox/actions/ai-report.md:131,137`. Exit 2 handling differs across the three callers. ai-report says "report its finding and never treat the layout as verified", while Step 8 accepts only `pass` or `skipped`; `stakeholder-report.md:47` and `architecture-report.md:103` ship with a "not render-verified" note. A persistent `error` leaves ai-report with no way to finish Step 8. Reproduced cause: a page with an unreachable external image (`http://10.255.255.1/x.png`) never reaches `readyState complete`, and the run ends `verdict=error`, exit 2, after 31 s. — impact-low → report only
- F3 `skills/do-work/docs/command-line-guide.md:26`. "The toolbox CLI verb `ai-report-judge` has no recipe" reads as if it is the only verb without a recipe. REQ-655's `ai-report-index` (the ai-report catalog verb already on main) has none either and is not named. Builder decision D-09 expected the integrator to merge both into this sentence, and the merge did not. — impact-low → report only
- F4 `skills/do-work-toolbox/actions/architecture-report.md:103`. "Test with network access disabled so rendering does not depend on GitHub or a CDN" now sits beside a command that has no offline mode and never fetches other origins in its link check, so a CDN-dependent draft passes the judge. The sentence predates this REQ but now reads as if the command covers it. — impact-low → report only
- F5 `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go:128,350`. Leftover temporary files. When the CLI is killed with SIGTERM, the engine exits (pipe EOF, checked with `pgrep`), but `ai-report-judge-profile-*` stays in `$TMPDIR` (reproduced). The default output directory (four full-page PNGs, which grow with page height: 1440x6036 for a 6000 px page) is removed only by ai-report Step 8, and the stakeholder-report and architecture-report steps never remove it. — impact-low → report only
- F6 `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go:22-30`. Test gaps. No automated test covers the `--out`-inside-bundle refusal (a stated constraint) or the error path stopping the server and engine. Both were checked by hand only. The three browser cases skip when discovery finds no engine, and no environment switch makes a missing engine fail (the board has `QUEUE_KANBAN_BROWSER` for this). On a heavy lane without Chrome the run is green with skips, the `silent-skip-reads-as-red` family. — impact-low → report only

**Nit:**
- F7 `skills/do-work-toolbox/actions/ai-report.md:12`. The Philosophy line still says "When browser automation is available, inspect full-page light and dark renders", with no phone width and the old condition wording. — impact-negligible → report only
- F8 `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go:136-137,273`. Small inconsistencies and bloat. A failure to create `--out` is labelled `TOOLBOX-USAGE`, while the same failure for the default temp directory (`:128-129`) uses `AI-REPORT-JUDGE-ERROR`, which is the reason D-14 added `judgeFailureResult`. `judgeErrorResult` is a 4-line wrapper with one caller. `SKILL.md:24-25` has two rows routing to the same file. — impact-negligible → report only

### Requirements Checklist

- [x] One command: free port, HTTP serve, 1440x1000 and 390x844 in light and dark, `scrollWidth > innerWidth`, relative `href`/`src` 404 check, `judge.json` outside the bundle, server stopped on every exit path. Delivered (`report_judge.go:88-230`; `defer server.Close()` at `:165`, `defer session.stop()` at `:174`).
- [x] Non-zero exit on overflow or broken link, each finding named in `judge.json`. Delivered (overflow with viewport, scheme, scroll and inner width; broken link with attribute, target as written, status).
- [x] Clean bundle exits 0, bundle unchanged, server closed after the run. Delivered (test plus manual run).
- [x] No browser gives `skipped`, exit 1, never a pass. Delivered (`:140-145`, `:265`).
- [x] `--out` or a fresh temp directory, path printed. Delivered. Text output prints `judge_json=` and `output_directory=`. JSON output names `judge.json` in `changes[0].path` because `ExactTextOutput` is `json:"-"` (D-13), and ai-report Step 7 tells the reader to look there.
- [x] ai-report Step 7 points at the command; architecture-report and stakeholder-report call it. Delivered, with F1 on the architecture-report call site.
- [x] Wired into the guide, `SKILL.md` routing and `help.md`. Delivered; both REQ-655 and REQ-657 lines survive the merge and read correctly.
- [x] Focused tests for overflow at 390 px and `<img src="missing.png">`. Delivered.
- [x] UR-144 A4 (the user's request for a bundled render check): matches the proposed direction. D-04 narrows "crawls" to the one `index.html` page, which is recorded.

### Acceptance Testing

**Result: Partial** (stages covered: Implementation, Integration)
- `DO_WORK_HEAVY_TESTS=1 go test -count=1 -run '^(TestAIReportJudge|TestHandlersRegisterCanonicalToolboxCommands)' ./internal/toolboxcommands/`: five PASS, no SKIP, 2.8 s, real Chrome. `go test -short ./internal/toolboxcommands/` passes. `gofmt -l` is empty. `go vet` and `git diff --check` are clean.
- CLI on scratch fixtures, text and JSON formats. Clean bundle: exit 0, `pass`. Overflow (`<pre>` 705 px wide): exit 1, two phone findings, no wide finding. Broken `missing.css`, `other.html`, `gone.png`: exit 1, three findings with status 404. Mixed bundle: a space in a file name, a directory link without a trailing slash, a query link, protocol-relative, `javascript:`, `data:` and `mailto:` links all handled correctly. A 6000 px page gives full-height captures.
- Refusals all exit 2 with `TOOLBOX-USAGE` and create nothing in the bundle: `--out` nested inside the bundle, `--out` through a symlink to the bundle, the bundle given through a symlink, `--out` equal to the bundle, `--browser /nope/chrome`, a missing `--browser` value, `--out=`, no arguments, a second positional argument, a directory without `index.html`.
- Cleanup: after the normal and timeout runs, no `ai-report-judge-profile` process or directory remained. SIGTERM case: see F5.
- The verb itself works end-to-end. The architecture-report call site fails on non-first reports (F1).

### Suggested Additional Testing

- Deployment: unassessed. Install the release into a consumer repo and confirm `<skill-root>/../do-work/tools/do-work-cli.sh ... ai-report-judge` resolves from the toolbox skill root.
- Live acceptance: unassessed. Run `do-work-toolbox ai-report` end to end on a real completed REQ with screenshots, and confirm the Step 7 rerun loop and the Step 8 cleanup.
- A second architecture report on a repository that already has one (F1).
- Linux engine discovery (`chromium` or `google-chrome` on `PATH`). Only the macOS app-bundle path was exercised.
- A machine without any Chromium-family engine: confirm the `skipped` footer flow in all three actions.

### Restatement and anti-bloat

- **Restatement sweep:** the diff redefines the render-check step (old: "serve over HTTP and take light and dark screenshots when browser automation is available"; new: the `ai-report-judge` verb with four verdicts and `judge.json`). Grep of `skills/` for `over HTTP`, `http.server`, `browser automation`, `render-verif`, `render check`, `judge capture` and `temporary HTTP`: no remaining copy of the old serve-and-screenshot prose. The stale or conflicting restatements are F1 (architecture-report Step 4 relative links), F2 (exit-2 handling), F3 (command-line-guide), F4 (offline sentence) and F7 (Philosophy line). `ai-report.md:77` (capturing a running app into `screenshots/`) is a different step and is correct as written.
- **Anti-bloat count:** one new 570-line file and one 167-line test file. New names: 4 finding codes, 3 constants, 3 package variables, 4 types, 1 handler, 9 functions, 6 session methods, 2 test helpers, 4 tests, 2 launch flags. All trace to D-01 to D-14 except the one-caller `judgeErrorResult` (F8). The production copy of the DevTools pipe transport is about 230 lines, against D-01's estimate of about 150; it is accepted by D-01. No new environment variable, recipe, runtime, config or decorative test. Each test names the failure it pins.

### Scores (on the record, not the headline)

**Overall: 80%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 92% | All eight criteria delivered; the architecture-report call site misbehaves (F1) |
| Code Quality | 88% | Bounded waits, defers right after each start, clean pipe-close shutdown; small label inconsistency (F8) |
| Test Adequacy | 85% | RED-GREEN plus mutation proof; no test for the `--out` refusal or error-path cleanup; silent skip without an engine (F6) |
| Scope | 97% | Touched set equals the 11-path write_set; D-09 to D-14 recorded |
| Risk | Low | Loopback-only server, `http.Dir` traversal-safe; temp leftovers (F5) |
| Acceptance | Partial | Verb end-to-end Pass; architecture-report non-first runs cannot reach `pass` |

Average (92+88+85+97)/4 = 90.5, minus 10 for Acceptance Partial = 80.5, recorded as 80%.

### Follow-ups created
None (8 findings report only)

## Review

**Overall: 80%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 97% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `architecture-report.md:93,103`: Step 4 requires relative links to the prior bundle, but the judge serves only the draft directory, so `../<prior>/index.html` always answers 404 (reproduced, exit 1) and "rerun until the verdict is `pass`" is unreachable on every non-first architecture report — impact-user-visible → report only

**Minor findings:** F2 `ai-report.md:131,137` exit-2 handling differs from stakeholder-report and architecture-report, and Step 8 accepts only `pass`/`skipped`, so a persistent `error` (e.g. a hanging external image, reproduced: 31 s, exit 2) blocks ai-report — impact-low → report only; F3 `command-line-guide.md:26` names `ai-report-judge` as the recipe-less verb but omits REQ-655's `ai-report-index` — impact-low → report only; F4 `architecture-report.md:103` "test with network access disabled" sits beside a command with no offline mode that skips other origins — impact-low → report only; F5 `report_judge.go:128,350` SIGTERM leaves `ai-report-judge-profile-*` in `$TMPDIR`, and default output dirs are removed only by ai-report Step 8 — impact-low → report only; F6 `report_judge_test.go:22-30` no test for the `--out` refusal or error-path cleanup, and browser cases skip silently without an engine — impact-low → report only; F7 `ai-report.md:12` Philosophy line keeps the old light/dark-only, "browser automation" wording — impact-negligible → report only; F8 `report_judge.go:136-137,273` `--out` creation failure labelled `TOOLBOX-USAGE`, one-caller `judgeErrorResult`, two routing rows to one file — impact-negligible → report only
**Acceptance:** Partial — Implementation and Integration stages: focused heavy tests pass with real Chrome, and CLI runs on scratch fixtures confirm overflow, broken links, refusals, both output formats and cleanup; the architecture-report call site fails on non-first reports (F1)
**Restatement sweep:** redefined the render-check step (serve-and-screenshot prose replaced by the `ai-report-judge` verb, its four verdicts and `judge.json`); stale restatements recorded as F1, F2, F3, F4, F7
**Suggested testing:** 5 items
**Follow-ups created:** None (8 findings report only)

*Reviewed by review-work action*
