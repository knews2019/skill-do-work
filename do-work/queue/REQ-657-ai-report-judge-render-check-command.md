---
id: REQ-657
title: 'ai-report judge runs the render check as a bundled command'
status: pending
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
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing shipped shell and publication scripts; family `assertion-passes-when-feature-is-dead` fits a render check that must never pass when the browser is missing.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8334 tokens, over budget; `slugged: partial`). Matching reason: family `disk-space-blind-spot` and the board's browser-probe lessons cover the same browser-driving mechanism.
## Full Context
See `do-work/user-requests/UR-144/input.md` for complete verbatim input (sections Request A4, What happened, Where the behaviour lives today, Proposed direction A4, Acceptance check). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A4: "A4 (optional, lower priority). `ai-report judge <dir>`, the existing Step 7 render check as a bundled script so it is not retyped per session."*
