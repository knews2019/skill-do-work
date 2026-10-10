---
id: REQ-657
title: 'ai-report judge runs the render check as a bundled command'
status: completed
created_at: 2026-10-09T21:10:00Z
user_request: UR-144
domain: testing
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
priority: later
related: [REQ-654, REQ-655, REQ-656]
batch: ai-report-modes
claimed_at: 2026-10-10T12:52:22Z
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go", "skills/do-work/docs/command-line-guide.md", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/architecture-report.md", "skills/do-work-toolbox/actions/stakeholder-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-10-10T13:20:05Z
  basis:
    - Route B
    - 10-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
builder_handback_at: 2026-10-10T14:03:37Z
integration_at: 2026-10-10T16:03:24Z
review_at: 2026-10-10T16:14:57Z
kb_status: pending
heavy_verified_at: 2026-10-10T16:24:12Z
heavy_verified_revision: fe0b19e2f383b5370c742af65f893aec8f388248
commit: fe0b19e2f383b5370c742af65f893aec8f388248
completed_at: 2026-10-10T16:24:46Z
release_at: 2026-10-10T16:24:46Z
---
# ai-report Judge Runs the Render Check as a Bundled Command
## What
`ai-report judge <dir>` turns the existing Step 7 render check into one shipped command, so it is not retyped per session. It serves the bundle over HTTP on a free port, captures wide and phone widths in light and dark, fails on horizontal overflow, checks relative links and images for 404s, writes a machine-readable result outside the bundle, and stops the server. `architecture-report` and `stakeholder-report` call the same command.
## Why
Report A4 (UR-144 input): a background `python3 -m http.server` for report QA was started 58 times across 22 sessions, and plans retype "capture the report in light and dark at wide and narrow widths and check for horizontal overflow". The check lives in a standalone `make-ai-report-with-screenshot` skill copied into consumer repos; one machine has eleven copies in five versions (501 to 599 lines) that agree only on the HTTP server step.
## Priority
`priority: later` comes from the report's own words: "A4 (optional, lower priority)".
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/ai-report.md:113`: "serve the report folder over HTTP—never `file://`—and take full-page screenshots in both light and dark rendering contexts. Save judge captures outside the report bundle..." No phone width, no overflow check, no link check, no machine-readable result.
- `skills/do-work-toolbox/scripts/` does not exist. The repo's existing browser checks are Go tests that drive a browser through `os/exec` (`skills/do-work-board/tools/queue-kanban/*_browser_probe_test.go`).
## Detailed Requirements
1. One command that, given a bundle directory: picks a free port and serves the bundle over HTTP; captures 1440x1000 and 390x844 in light and dark; fails when `scrollWidth > innerWidth`; checks relative `href` and `src` targets for 404s; writes `judge.json` outside the bundle; stops the server on every exit path.
2. Exit non-zero on any overflow or broken relative link, and name each finding in `judge.json`.
3. `ai-report.md` Step 7 points at the command instead of describing the steps in prose. `architecture-report.md` and `stakeholder-report.md` call the same command where they render-check.
4. Wire `judge` into `skills/do-work-toolbox/docs/ai-report-guide.md`, `skills/do-work-toolbox/SKILL.md:23` routing and `skills/do-work-toolbox/actions/help.md:13`.
5. A focused test: a fixture page with a horizontal overflow at 390px and a fixture page with a broken relative image link each make the command exit non-zero with the finding in `judge.json`.
6. Read `_dev/primes/prime-shell-commands.md` before writing any shipped shell. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Only the render check is in scope; do not merge the rest of the standalone screenshot-report skill.
- The bundle stays untouched: captures and `judge.json` go outside it.
- No new queue fields and no new statuses.
- This REQ is its own release; do not fold REQ-654, REQ-655 or REQ-656 into it.
## Assumptions (recorded at capture, no questions asked)
- Tool home and language are the builder's call. The report suggests `skills/do-work-toolbox/scripts/judge.mjs` plus a serve helper, but the toolbox has no `scripts/` folder and the repo's browser checks are Go driving a browser by `os/exec`. Prefer a toolbox CLI verb in that style over adding a Node runtime dependency; the acceptance behaviour is the contract, not the file name. `write_set` is omitted because the home is not decided.
- When no browser is available, the command reports a distinct "skipped" result and never reports a pass; a skip is not a pass.
- `judge.json` and the captures go to an output directory given as an argument, defaulting to a fresh temporary directory, and the command prints that path.
- Crawling covers relative links only; external URLs are not fetched.
- Deleting the drifted standalone copies in consumer repos is operator work for the user after release, not part of this REQ.
## Dependencies
None. Independent of REQ-654, REQ-655 and REQ-656.
## Builder Guidance
Certainty is high on behaviour, low on tool home. Latitude: implementation language, port selection, browser discovery, and the exact `judge.json` shape.
## Red-Green Proof
**RED prompt/case:** A fixture bundle whose page overflows horizontally at 390px wide, and a fixture bundle with `<img src="missing.png">`; run `ai-report judge <dir>` on each.
**Why RED now:** No `judge` form exists; the check is prose at `ai-report.md:113` with no phone width, overflow check, link check or result file.
**GREEN when:** Each run exits non-zero and writes `judge.json` naming the overflow (with width) or the broken relative link; a clean fixture exits zero; the server is stopped after each run; with no browser available the run reports skipped, not passed.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing shipped shell and prescribed command blocks (the action fences that call the verb); family `assertion-passes-when-feature-is-dead` fits a render check that must never pass when the browser is missing.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8815 tokens, over budget; `slugged: partial`). Matching reason: its row covers browser behaviour testing, the same browser-driving mechanism; family `disk-space-blind-spot` covers screenshot output growth.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens, over budget; `slugged: partial`). Matching reason: the verb is new do-work-cli internals; families `reaped-by-its-own-parent` (stopping the engine) and `silent-skip-reads-as-red` (a skip is not a pass) fit.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over budget; `slugged: partial`). Matching reason: its row covers changing action routing (`SKILL.md:23`) and the owning prime `prime-action-files.md` is in `prime_files`.
## Full Context
See `do-work/user-requests/UR-144/input.md` for complete verbatim input (sections Request A4, What happened, Where the behaviour lives today, Proposed direction A4, Acceptance check). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Follow pre-dispatch D-01 to D-08. Tests first (four named tests plus the registration count), RED as a compile failure. Then one file `report_judge.go`: argument parsing with `optionValue`; bundle and `--out` resolution against `RepositoryRoot`; `--out` refused inside the bundle after `EvalSymlinks`; `--browser` resolved or usage failure; discovery from two package-level lists; skip path writes `judge.json` and returns OutcomeFindings; `judgeRender` owns server and engine with a `defer` right after each starts; four passes via `Emulation.*`, `Page.navigate`, a readiness poll, measurement and `Page.captureScreenshot`; links collected in the wide light pass and fetched from our own server only. Then prose in the seven Markdown files. Lessons read: `lessons-releases.md` (whole), families `reaped-by-its-own-parent` and `silent-skip-reads-as-red` (targeted grep). The reaped-by-its-own-parent lesson shaped `stop()`: close the command pipe first so the engine exits on EOF and reaps its own helpers, and kill the leader only after the 30 s deadline. The silent-skip lesson shaped the skip path: `skipped` is OutcomeFindings (exit 1), announced as its own code `AI-REPORT-JUDGE-SKIPPED`, and never `pass`.
- [x] **[APPLY]:** (from the builder hand-back) Done as planned in two commits. Deviations are listed under Decisions (D-09 to D-14).
- [x] **[UNIFY]:** (from the builder hand-back) `git diff bd56c4b0 --stat`:

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
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A4: "A4 (optional, lower priority). `ai-report judge <dir>`, the existing Step 7 render check as a bundled script so it is not retyped per session."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The behaviour is fully specified by the REQ and the UR's acceptance check; only the tool home and the browser mechanism were open, and exploration settles both by following the existing toolbox CLI verb pattern (`internal/toolboxcommands/`) and the board's browser-over-`os/exec` pattern. No architectural change, one new verb plus prose wiring in five action/doc files.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Done by the pre-dispatch agent at main `bd56c4b0` (no Explore subagent; the files are few and named).

**Where the verb goes.** The toolbox's deterministic phases are Go verbs in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`. `commands.go:11-19` declares the verb-name constants and `commands.go:21-31` registers them in `Handlers()`; `cmd/do-work-cli/main.go` loops over `toolboxcommands.Handlers()`, so registering the handler is the whole wiring. Handler signature: `func(commandruntime.ExecutionContext, []string) resultmodel.CommandResult` (`internal/commandruntime/command_runtime.go:14-19`); tests call the handler directly with `commandruntime.ExecutionContext{RepositoryRoot: dir}` (`architecture_test.go:37,61`). Helpers to reuse: `usageResult`, `toolboxFinding`, `optionValue`, `exactOutputResult` (`commands.go:33-100`). Exit codes come only from `resultmodel.ExitCode` (`internal/resultmodel/result_model.go:647-662`): success 0, findings/refused 1, failure 2. The module is standard-library only (`go.mod` has no `require`), so there is no HTML parser and no WebSocket library: link discovery has to come from the browser, and the browser transport has to be the pipe form.
**The browser pattern to copy.** `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` drives Chromium through `os/exec` with `--remote-debugging-pipe` (DevTools protocol on fd 3/4, NUL-terminated JSON): launch flags `:276-284`, launch and pipes `:255-320`, attach to the page target `:323-400`, call a method `:403-458`, decode `:460-481`, evaluate `:483-551`, wait for a page condition `:553-592`, close and reap `:641-672`. Browser discovery is `:69-95` (override env, then a well-known-name list documented as illustrative, never closed). The light/dark switch there is `--blink-settings=preferredColorScheme=0|1` (`:1019-1024`); inside one protocol session `Emulation.setEmulatedMedia` with `prefers-color-scheme` does the same per pass. This is test code in another module, so the judge needs its own small production copy of the transport (launch, attach, call, evaluate, close), not an import.
**The heavy-test convention.** Browser-launching Go tests in do-work-cli gate on `testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1"` (`internal/nextselection/next_commands_test.go:49`). The fast gate runs `-short`; the heavy gate and the `do-work-cli-integrations` heavy lane set `DO_WORK_HEAVY_TESTS=1` (`_dev/tests/maintainer-verify.sh:783,866`).
**Prose that changes.**
- `skills/do-work-toolbox/actions/ai-report.md:111-115` (Step 7: the serve-and-screenshot sentence at `:113` becomes the command; the judgement list at `:115` stays, applied to the captures), `:119` (Step 8: "every relative asset resolves" and "stop the temporary HTTP server" are now the command's job), `:30` Input (add the `judge <dir>` form as its own sentence, do not rewrite the `:30` line), `:139` checklist line.
- `skills/do-work-toolbox/actions/stakeholder-report.md:39-41` (Step 5 Render-Check) calls the command.
- `skills/do-work-toolbox/actions/architecture-report.md:97` ("Open the HTML locally in a browser when available ...") calls the command; the "network access disabled" sentence stays.
- `skills/do-work-toolbox/SKILL.md:23` routing row, `skills/do-work-toolbox/actions/help.md:13` menu line, `skills/do-work-toolbox/docs/ai-report-guide.md:55` (render paragraph) and `:63-70` (Input block).
- `skills/do-work/docs/command-line-guide.md:26` lists every toolbox verb; it gains `ai-report-judge`.
**Lock-in test that must move.** `internal/toolboxcommands/commands_test.go:5-15` lists every toolbox verb and asserts `len(handlers) != 7`; it gains the new constant and the count goes to 8.
**Not needed.** No Just recipe: `skills/do-work-board/justfile.template:147-169` has recipes for some toolbox verbs, but no test couples verbs to recipes and the REQ names none; the actions call `do-work-cli.sh` directly the way `architecture-report.md:41` does. No prescribed-shell case: `_dev/tests/prescribed-shell-cases/` covers shim scripts, not a `do-work-cli.sh` call in an action fence.
**Lessons.** The consult added `_dev/primes/lessons-releases.md` (666 tokens, full). The over-budget satellites are listed below under Dropped for Budget. Families worth a targeted grep by the builder even though the satellite is dropped: `reaped-by-its-own-parent` and `silent-skip-reads-as-red` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (killing the browser, and a skip that must never read as a pass).

### Pre-dispatch decisions (DECIDE & STATE; the REQ gave the builder latitude on each)
- **D-01** Tool home: a Go toolbox verb `ai-report-judge` in `internal/toolboxcommands/report_judge.go`, invoked as `ai-report judge <dir>` from the action. No Node runtime, no `skills/do-work-toolbox/scripts/`. Reason: the REQ's assumption prefers it, and the CLI is already the toolbox's deterministic engine. Value: one shipped binary path, no new runtime. Risk: a production copy of a small DevTools pipe transport (about 150 lines); reversible.
- **D-02** Browser mechanism: one headless Chromium-family process per run over `--remote-debugging-pipe`. Per pass: `Emulation.setDeviceMetricsOverride` (1440x1000 or 390x844, scale 1, mobile false), `Emulation.setEmulatedMedia` (`prefers-color-scheme` light or dark), `Page.navigate` to the served `index.html`, wait for `document.readyState === "complete"`, `Runtime.evaluate` for `document.documentElement.scrollWidth` and `window.innerWidth`, `Page.captureScreenshot` full page (`captureBeyondViewport` with a clip to the content height). Not the `--screenshot` / `--dump-dom` flags: they cannot measure without injecting script into the served page. Every protocol wait has a timeout (30 s is fine) so a hung engine ends the run instead of hanging it.
- **D-03** Browser discovery: `--browser <path>` first (a value that names nothing is a usage failure, never a silent fallback), then well-known names on `PATH` (`google-chrome`, `google-chrome-stable`, `chromium`, `chromium-browser`, `chrome`), then the standard macOS app-bundle executables for Google Chrome and Chromium. The lists are a convenience, documented as illustrative; `--browser` makes any other engine usable. No new environment variable. The lists are package-level variables so the skip test can empty them.
- **D-04** What is checked: the bundle's `index.html` only, no crawl into linked pages. Links: every element carrying `href` or `src` on the rendered page (collected once, in the wide light pass), each resolved against the page URL; a target on the judge's own origin is fetched from the judge's server with its fragment removed, and a status of 400 or more is a broken-link finding; other origins and non-HTTP schemes (`mailto:`, `data:`, `javascript:`) are skipped, never fetched.
- **D-05** Verdicts and exit codes: `pass` (outcome success, exit 0); `fail` (outcome findings, exit 1; finding codes `AI-REPORT-JUDGE-OVERFLOW` naming viewport, scheme, scroll width and inner width, and `AI-REPORT-JUDGE-BROKEN-LINK` naming attribute, target and status); `skipped` when no browser is found (outcome findings, exit 1, code `AI-REPORT-JUDGE-SKIPPED`, `judge.json` still written, link check not run); `error` when the engine starts but the protocol, navigation or a capture fails or times out (outcome failure, exit 2). A skip never exits 0.
- **D-06** Output: `--out <dir>` (created when missing), default a fresh `os.MkdirTemp` directory; refused with a usage failure when it resolves inside the bundle. It holds `judge.json` and `wide-light.png`, `wide-dark.png`, `phone-light.png`, `phone-dark.png`. The engine profile goes in its own temporary directory, removed on exit. Exact text output: `verdict=<v>`, `judge_json=<absolute path>`, `output_directory=<absolute path>` lines. `judge.json` carries at least: `schema_version`, `bundle`, `served_url`, `verdict`, `browser`, `captures` (viewport, width, height, scheme, file, scroll_width, inner_width), `findings` (kind plus the fields in D-05).
- **D-07** Server: `net.Listen("tcp", "127.0.0.1:0")` plus `http.FileServer` over the bundle; closed by a `defer` set right after the listener opens, and the engine is killed and reaped by a `defer` set right after it starts, so every return path stops both.
- **D-08** Tests: one RED case per named failure (overflow at 390, broken relative image), one clean bundle that passes and proves the server stopped and the bundle is unchanged, and one no-browser run that reports skipped. Browser cases use the heavy-test gate above and skip when no engine is found; the skip case never launches a browser and always runs.
<!-- D-XX counter: last used D-08. Next decision: D-09. -->

*Generated by the pre-dispatch agent (exploration done directly)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (new): the ai-report-judge verb, its HTTP server, browser discovery, the DevTools pipe transport, judge.json
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (new): overflow, broken relative image, clean bundle with server stopped, no-browser skipped
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modify): verb constant and Handlers registration
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modify): registration lock-in lists the new verb, count 7 to 8
- `skills/do-work/docs/command-line-guide.md` (modify): toolbox verb list gains ai-report-judge
- `skills/do-work-toolbox/actions/ai-report.md` (modify): judge form in Input, Step 7 and Step 8 call the command, checklist line
- `skills/do-work-toolbox/actions/architecture-report.md` (modify): the browser inspection sentence calls the command
- `skills/do-work-toolbox/actions/stakeholder-report.md` (modify): Step 5 Render-Check calls the command
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify): render paragraph and Input block
- `skills/do-work-toolbox/SKILL.md` (modify): routing phrases for the judge form
- `skills/do-work-toolbox/actions/help.md` (modify): menu line for the judge form

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (the pattern is copied, not imported or refactored), `skills/do-work-board/justfile.template` and the root justfile (no Just recipe), `_dev/tests/` (no prescribed-shell case, no gate change), `skills/do-work-toolbox/actions/ai-report-reference.md` and `completed-work-presentation-reference.md` (REQ-687's area), `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` and version mirrors, anything under `do-work/` (the release and the REQ record are the integrator's).

**Acceptance criteria (restated from REQ):**
- [ ] One command, given a bundle directory, picks a free port and serves the bundle over HTTP, captures 1440x1000 and 390x844 in light and dark, checks `scrollWidth > innerWidth`, checks relative `href` and `src` targets for 404s, writes `judge.json` outside the bundle, and stops the server on every exit path.
- [ ] It exits non-zero on any overflow or broken relative link, and `judge.json` names each finding (overflow with its width; broken link with its target).
- [ ] A clean bundle exits zero; the bundle's files are unchanged after the run; the server no longer accepts connections after the run.
- [ ] With no browser available the run reports `skipped`, distinct from a pass, and never exits zero.
- [ ] `judge.json` and the captures go to `--out <dir>` or a fresh temporary directory, and the command prints that path.
- [ ] `ai-report.md` Step 7 points at the command instead of describing the steps in prose; `architecture-report.md` and `stakeholder-report.md` call the same command where they render-check.
- [ ] `judge` is wired into `docs/ai-report-guide.md`, the `SKILL.md` routing row and `actions/help.md`.
- [ ] Focused tests: a fixture page that overflows at 390px and a fixture page with `<img src="missing.png">` each make the command exit non-zero with the finding in `judge.json`.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/stakeholder-report.md` (modified)
- `skills/do-work-toolbox/actions/architecture-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (modified)

**What was done:** Added the toolbox CLI verb `ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]`. It serves the bundle with a loopback `http.FileServer` on `127.0.0.1:0`, finds a Chromium-family browser (`--browser` first, then illustrative name lists), drives it over a small DevTools pipe transport, and takes four full-page captures (1440x1000 and 390x844, light and dark). It fails on `scrollWidth > innerWidth` and on any same-origin `href`/`src` answering 400 or more, writes `judge.json` plus the four PNGs to `--out` or a fresh temporary directory outside the bundle, and stops its server and browser through `defer`s on every return path. Verdicts: `pass` (exit 0), `fail` and `skipped` (exit 1, a skip is never a pass), `error` (exit 2). The verb is registered in `commands.go`. `ai-report.md` Step 7 and Step 8, `stakeholder-report.md` Step 5 and `architecture-report.md` Step 5 call the command; the guide, the toolbox routing table, help and `command-line-guide.md` name it.

Hand-back merge `aea1c4d3..3f8fd77c` (integrator): five files conflicted with REQ-655 (ai-report index and find, already on main) and were resolved by keeping both texts: the toolbox `SKILL.md` keeps REQ-655's extended `ai-report` row and adds the judge row below it; `help.md` and the guide's Input block keep both new lines; `ai-report.md` keeps the judge Input sentence followed by REQ-655's `### Catalog forms` subsection, and the checklist keeps REQ-657's render line and REQ-655's catalog exception; `commands_test.go` lists both verbs and the handler count is 9. The `ai-report.md` `$ARGUMENTS` line now says the `judge`, `index` and `find` forms create no report (REQ-655 review F2, the stale line next to this merge hunk); the other two stale lines that F2 named (`ai-report.md` Output Format, `ai-report-guide.md:21`) are outside this merge and stay as they are.

After the first heavy drain and the review (builder-branch commits `4021963c` and `1ae7a03b` by the integrator, re-merged with the same `<pre>` as `fe0b19e2`): the toolbox `SKILL.md` routes `ai-report` in one row again, holding REQ-655's and REQ-657's phrases, because the staged-skills heavy contract requires each public action to be routed exactly once (`_dev/tests/staged-skills-contract.sh:787`) and the two-row form failed it. Review fixes: `architecture-report.md` Step 5 treats the `../` prior-bundle link finding as expected and checks that link after publication (F1); `ai-report.md` Step 7, Step 8 and the checklist disclose an exit-2 check failure in the footer like `skipped` (F2); `command-line-guide.md` names both recipe-less verbs, `ai-report-index` and `ai-report-judge` (F3).

## Decisions

(from the builder hand-back, verbatim; D-01 to D-08 are in `## Exploration`)

- **D-09** (DECIDE & STATE) `skills/do-work/docs/command-line-guide.md:26` lists Just recipes ("`just --list` is the live inventory installed from the managed template"), and the brief forbids a Just recipe. Appending `ai-report-judge` bare would claim a recipe that does not exist. I appended one sentence to the same line instead: "The toolbox CLI verb `ai-report-judge` has no recipe; actions call it through `do-work-cli.sh`." The probe's `grep -F ai-report-judge` passes. Reversible; if REQ-655 appends a recipe-less verb to the same line, the integrator can merge both into this sentence.
- **D-10** (DECIDE & STATE) The architecture-report draft is "one self-contained HTML document ... outside the reports directory" (`architecture-report.md:79`), a file, while the verb serves a directory holding `index.html`. The Step 5 text tells the agent to keep the draft as `index.html` in a directory of its own and pass that directory. `:79` itself is not edited. The publish command already takes `<draft-path>`, which is that `index.html`.
- **D-11** (DECIDE & STATE) Launch flags add `--hide-scrollbars` and `--no-first-run` to the board's set. `--hide-scrollbars` keeps `innerWidth` equal to the layout width, so a classic scrollbar cannot hide a small overflow under `scrollWidth > innerWidth`. The engine starts on `about:blank` and attaches to the first page target, and each pass navigates with `Page.navigate`.
- **D-12** (DECIDE & STATE) Readiness after `Page.navigate`: the same URL is loaded four times, so `readyState === "complete"` alone could read the previous document. Before each navigation the judge sets `window.aiReportJudgeStale = true` on the old window, then polls `!window.aiReportJudgeStale && document.readyState === "complete"` (evaluate errors during the swap count as not ready, bounded by 30 s). The marker lives only in the judge's own engine; the bundle is never modified.
- **D-13** (DECIDE & STATE) `--format json` hides the exact text output, so the `judge.json` path is reported as a `changes` entry (`kind: created`). The ai-report Step 7 prose tells the reader to find it there.
- **D-14** (DECIDE & STATE) `--browser` and the bundle check are validated before the default temporary output directory is created, so a usage failure leaves no empty temp directory (found in a manual run and fixed before the commit). A broken link's `target` is the attribute value as written (`missing.png`), not the resolved URL, so the finding names what the author must fix.

## Discovered Tasks

(from the builder hand-back; impact tokens are the builder's)

- `skills/do-work/docs/command-line-guide.md` mixes two inventories: the Just-recipe list and, now, a CLI verb without a recipe. If more recipe-less toolbox verbs arrive (REQ-655 is one candidate), the line may need a separate "CLI verbs without recipes" bullet. impact-low → report only
- `skills/do-work-toolbox/actions/architecture-report.md:79` still says "one self-contained HTML document as a draft" without naming the file; D-10 adds the `index.html`-in-its-own-directory rule only at Step 5. A later edit could move that rule to `:79`. impact-low → report only

## Qualification

**Gate records (`advance --diff-range aea1c4d3..3f8fd77c`):** `qualify` satisfied (success), `scope-drift` satisfied (success). Two `QUALIFY-NEW-FILE-UNWIRED` warnings, on `report_judge.go` and `report_judge_test.go`, judged false positives: `commands.go:33` registers `handleAIReportJudge`, and a `_test.go` file is found by the Go test runner by convention. No debug-artifact, P-A-U or output-primitive findings.

**Scope:** the declared `write_set` (11 paths) equals the touched set in `git diff aea1c4d3..3f8fd77c --stat` (11 files, 774 insertions, 11 deletions). The queue guard printed nothing before the merge. Five files conflicted with REQ-655 and were resolved by keeping both texts (see `## Implementation Summary`); after the resolution the merged tree builds, vets, is gofmt-clean, and `go test -short ./internal/toolboxcommands/` passes.

**Requirement trace (read against the merged files):**
1. One command: `handleAIReportJudge` (`report_judge.go:88`). Free port: `net.Listen("tcp", "127.0.0.1:0")` plus `http.FileServer` over the bundle (`:160-166`). Captures: `judgeViewports` 1440x1000 and 390x844 (`:51-54`) times light and dark, through `Emulation.setDeviceMetricsOverride`, `Emulation.setEmulatedMedia` and `Page.captureScreenshot` (`renderPass`, `:509`). Overflow: `capture.ScrollWidth > capture.InnerWidth` (`:188`). Links: every `[href],[src]` collected in the wide light pass, resolved against `document.baseURI`, fetched only when the origin is the judge's own server, fragment removed, status 400 or more is a finding (`:192-228`). Stop: `defer server.Close()` right after the server starts and `defer session.stop()` right after the engine starts (`:165`, `:174`); `stop` closes the command pipe, waits up to 30 s, then kills and reaps, and removes the profile.
2. Non-zero and named: `fail` maps to OutcomeFindings (exit 1); overflow findings carry viewport, scheme, scroll width and inner width; broken-link findings carry attribute, target as written, and status (`judgeResult`, `:232`).
3. Clean bundle: `pass` maps to OutcomeSuccess (exit 0). The bundle is only read through `http.Dir`; the readiness marker lives on the engine's window, never in a file (D-12). The clean-bundle test checks the server no longer accepts connections.
4. No browser: `skipped` with code `AI-REPORT-JUDGE-SKIPPED`, OutcomeFindings (exit 1), `judge.json` still written (`:140-145`, `:265`).
5. Output: `--out` or `os.MkdirTemp`, refused inside the bundle after symlink resolution of the longest existing prefix (`:127-135`); the text output prints `verdict=`, `judge_json=` and `output_directory=`; JSON output names `judge.json` in `changes` (D-13).
6. Prose: `ai-report.md` Step 7 runs the command in a bash fence and keeps the judgement list and the two-pass SVG rule; Step 8 reads the verdict and removes the output directory; `stakeholder-report.md` Step 5 and `architecture-report.md` Step 5 call the same command (architecture keeps "network access disabled" and "never claim an unperformed check passed").
7. Wiring: the guide's render paragraph and Input block, the toolbox `SKILL.md` routing row (`ai-report judge`, `render check`, `judge report`), the `help.md` menu line, and `command-line-guide.md` (D-09: a sentence, not a recipe name).
8. Focused tests: `TestAIReportJudgeFailsOnPhoneWidthOverflow`, `TestAIReportJudgeFailsOnBrokenRelativeImage`, `TestAIReportJudgePassesCleanBundleAndStopsServer`, `TestAIReportJudgeReportsSkippedWithoutBrowser`; the three browser cases run only under `DO_WORK_HEAVY_TESTS=1`, and the probe sets it.

**Builder deviations judged:** D-09 holds (`command-line-guide.md:26` is a Just-recipe inventory, so a bare verb name would claim a recipe). D-10 is the smallest change that fits a directory-serving verb to a single-file draft. D-11 (`--hide-scrollbars`) keeps `innerWidth` equal to the layout width so a classic scrollbar cannot hide overflow. D-12 is needed because the same URL is navigated four times. D-13 and D-14 are local and reversible.

**Constraints:** render check only; the bundle stays untouched; no new queue field or status; own release.

**Re-merge, cumulative range `aea1c4d3..fe0b19e2`:** `advance` takes no qualify input at this phase, so the integrator ran the `qualify --request-path <P> --diff-range aea1c4d3..fe0b19e2` handler directly: success, only the judged `QUALIFY-NEW-FILE-UNWIRED` warnings. The touched set is still the 11 declared files (774 insertions, 12 deletions); the delta touches four Markdown files only. The queue guard before the re-merge printed nothing.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `3f8fd77c` (machine quiet before launch: 1-minute load 4.82, no other gate running; load 6.04 at the end), then `advance REQ-657 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-657-probe.sh`.
**Result:** ✓ Repository gate passed on the first run (exit 0, gate wall 150 s; do-work-cli 902 tests in 73 s, slowest file `internal/finalization/finalization_recovery_test.go` 22.45 s under the 30 s limit; queue-kanban 420 tests in 49 s, slowest 21.08 s). Probe exit 0: it runs the focused tests with `DO_WORK_HEAVY_TESTS=1`, fails on any `--- SKIP:`, and checks the prose wiring. Gate records `green-gate`, `scope-drift` and `run-blocked-check` satisfied. A direct probe rerun at `3f8fd77c` gave the same five `--- PASS:` lines with a real Chrome (overflow 0.88 s, broken image 0.85 s, clean bundle 0.90 s, skipped 0.00 s, registration 0.00 s) and no skip.

**Red-green validation:** (from the builder hand-back, traced to `## Red-Green Proof`)
- `TestAIReportJudgeFailsOnPhoneWidthOverflow`, `TestAIReportJudgeFailsOnBrokenRelativeImage`, `TestAIReportJudgePassesCleanBundleAndStopsServer`, `TestAIReportJudgeReportsSkippedWithoutBrowser` (`internal/toolboxcommands/report_judge_test.go`): ✗ before the production code (build failure: `undefined: discoverJudgeBrowser`, `undefined: judgeReport`, `undefined: handleAIReportJudge`, `undefined: CommandAIReportJudge`; at base the launcher answered `ai-report-judge` with `UNKNOWN-COMMAND`, exit 2) → ✓ after. The overflow fixture overflows at 390 px (scroll width 600) in both phone passes; the broken fixture's `missing.png` answers 404. Mutation check: with the comparison changed to `ScrollWidth > InnerWidth+1000` the overflow test failed; the mutation was reverted before the commit.
- `TestHandlersRegisterCanonicalToolboxCommands`: updated for the new verb; after the merge with REQ-655 it lists both verbs and expects 9 handlers, and passes.

**New tests added:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go` (four tests; three need a browser and run only under `DO_WORK_HEAVY_TESTS=1`, the skip case always runs)

**Existing tests updated:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (registration list and count, 9 after the merge)

**Heavy verification plan:**
- Range: aea1c4d35f4769fc0020ea2db73ed3bf62115681..3f8fd77cd4cc6269ebe79bc045196541467f249d
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

**After the re-merge (`fe0b19e2`):** the first gate run (load 2.13 before launch) failed after 29 s because the disk was full (`No space left on device` in several stages; 131 MB free, machine-wide, not caused by this change). The integrator trimmed Go build-cache entries unused for more than a day (about 2.3 GB, the same age rule Go's own cache trim uses) and reran: exit 0, gate wall 170 s (load 4.11 before, 8.81 at the end), do-work-cli 902 tests in 84 s, slowest file `internal/finalization/finalization_recovery_test.go` 26.74 s under the 30 s limit; queue-kanban 420 tests in 53 s. `advance` refuses gate input past this phase, so the green record was written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` at `fe0b19e2`, and `REQ-657-probe.sh` was run directly: GREEN, five `--- PASS:` lines with real Chrome, no skip. The heavy plan for `aea1c4d3..fe0b19e2` selects the same four lanes.

*Verified by work action*

## Review

**Overall: 80%** | 2026-10-10T16:14:57Z

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

**Delta re-review** | 2026-10-10T16:24:12Z | delta `3f8fd77c..fe0b19e2` (builder-branch commits `4021963c`, `1ae7a03b`)

**Overall: 91%** (Requirements 92% → 96%, Acceptance Partial → Pass; others unchanged). F1, F2 and F3 closed; the routing part of F8 closed (one `ai-report` row holding all nine phrases). F4, F5, F6, F7 and the code part of F8 stay open, report only.
- F9 (nit) `ai-report.md:137`: Step 8 checks the last `judge.json` verdict, but some exit-2 paths (temp-directory failure, usage refusals) write no `judge.json`; Step 7's "report its finding" from the command result covers them. impact-negligible → report only

**Follow-ups created:** None (7 findings report only)

## Lessons Learned

**What worked:** Copying the board's DevTools pipe transport in a small production form kept the module standard-library only. Closing the command pipe first and killing only after a deadline left no engine process or profile behind in normal and timeout runs. A mutation check (`ScrollWidth > InnerWidth+1000`) proved the overflow test can fail.
**What didn't:** A second routing row for the `ai-report judge` sub-form failed the staged-skills heavy lane, which requires each public action to be routed exactly once; the fast gate and the builder's checks did not run that contract, so it surfaced only in the heavy drain. The architecture-report call site could never reach `pass` because the judge serves only the draft directory and the prior-bundle `../` link always answers 404 (review F1).
**Worth knowing:** `document.readyState === "complete"` passes at once when the same URL is navigated again, because the old document is still complete; mark the old window and wait for the mark to disappear (D-12). A sub-form of an existing action extends that action's routing row; it never gets its own row.

## Orientation

Now `ai-report-judge` (`toolboxcommands/report_judge.go`) is the one render check behind `ai-report` Step 7, `stakeholder-report` Step 5 and `architecture-report` Step 5; it needs a Chromium-family browser and reports `skipped` without one. [MAP CHANGED]: a new toolbox CLI verb that drives a headless browser from shipped code. Touched primes (`prime-shell-commands.md`, `prime-action-files.md`, `prime-releases.md`) name no path this change removed.

## Heavy Verification Plan

- Base: aea1c4d35f4769fc0020ea2db73ed3bf62115681
- Target: fe0b19e2f383b5370c742af65f893aec8f388248
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: fe0b19e2f383b5370c742af65f893aec8f388248
- Execution revision: fe0b19e2f383b5370c742af65f893aec8f388248 (detached checkout under `.git/work-run-work-2026-10-10-131527/`, `QUEUE_KANBAN_BROWSER` set)
- do-work-cli-integrations: exit 0, reused (executed at `3f8fd77c`: exit 0, 79 s; the delta touches no file under its subtree)
- staged-skills: exit 0, executed, 40 s
- updater: exit 0, reused (executed at `3f8fd77c`: exit 0, 72 s)
- installer: exit 0, reused (executed at `3f8fd77c`: exit 0, 35 s)
- The first drain at `3f8fd77c` was red on staged-skills only (41 s): "do-work-toolbox must route ai-report exactly once (found 2)". Fixed by the one-row routing in `4021963c`.

## Timing

Observed 2026-10-10T16:02:19Z to 2026-10-10T16:23:22Z: 21m 03s total, 23m 22s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 15m 00s | 4 |
| review | 6m 50s | 1 |
| handback-merge | 1m 32s | 2 |

Slowest stage: review / review, 6m 50s, outcome success.
Slowest command: verification-gate / heavy drain (staged-skills red), 4m 19s, exit 1, .

Notes: the builder-work event was skipped because the hand-back had already landed when integration began (`actions/fan-out-reference.md` → Landed hand-back); `builder_handback_at` is the builder commit's committer date. The 37 s delta re-review has no event (the recorder takes no end time). The first re-merge gate run that failed on a full disk is inside the "repository gate at re-merge" event.
