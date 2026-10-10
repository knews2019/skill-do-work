---
id: REQ-657
title: 'ai-report judge runs the render check as a bundled command'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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
