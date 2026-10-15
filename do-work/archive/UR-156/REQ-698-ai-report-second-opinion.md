---
id: REQ-698
title: 'ai-report second-opinion: a second model rebuilds a report bundle as index-<cli>-<model-slug>.html (session, targeted and batch modes)'
status: completed
created_at: 2026-10-15T11:02:38Z
user_request: UR-156
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "skills/do-work-toolbox/docs/prime-agy-usage.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-*.html", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/*-reply.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go", "skills/do-work/tools/do-work-cli.sh", "_dev/tests/do-work-cli-launcher-behavior.sh", "skills/do-work-toolbox/actions/ai-report-reference.md", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/gemini-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/agy-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-claude-opus-5-5-reply.md", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-gpt-6-astra-reply.md", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md", "do-work/lessons-index.md"]
estimate:
  p50_active_minutes: 75
  confidence: low
  calculated_at: 2026-10-15T11:02:38Z
  basis:
    - Route C
    - 10-file write set
    - 4 new files
    - 4 subsystems involved
    - 16 acceptance criteria
    - browser evidence
    - cross-route regression gates
claimed_at: 2026-10-15T11:42:48Z
route: C
planning_at: 2026-10-15T11:46:56Z
dispatch_at: 2026-10-15T11:53:14Z
builder_handback_at: 2026-10-15T12:44:02Z
integration_at: 2026-10-15T12:45:19Z
review_at: 2026-10-15T13:01:16Z
remediation_at: 2026-10-15T13:08:52Z
re_review_at: 2026-10-15T13:16:34Z
commit: 31e7fcdfd8fc7a934a5ca3f13202e66a9f8c3200
kb_status: pending
kb_entry: ''
heavy_verified_at: 2026-10-15T13:27:56Z
heavy_verified_revision: 31e7fcdfd8fc7a934a5ca3f13202e66a9f8c3200
completed_at: 2026-10-15T13:28:14Z
release_at: 2026-10-15T13:28:14Z
---
# ai-report second-opinion: A Second Model Rebuilds a Report Bundle
## What
Add the `ai-report second-opinion` form to the toolbox `ai-report` action. It takes an existing report bundle under `ai-reports/` (a folder holding `index.html`), hands it with a clarity-editor prompt to another model, and gets back a clarity-edited sibling page `index-<cli>-<model-slug>.html` in the same bundle, plus `<cli>-<model-slug>-reply.md` holding that model's corrections and open questions. The mechanical part is a do-work-cli subcommand `ai-report-second-opinion` that runs the backend, derives the file name from the model that actually ran, checks that nothing outside the bundle changed, render-judges the new page, and inserts a marker-wrapped banner into `index.html` linking to it.
## Why
The user wants a second opinion on AI-written reports from a different model, so a reader gets an easier, verified version next to the original. Migrated from the consumer repository iwr2-imprint-wright (UR-061 / REQ-287 there, cancelled in favour of this capture) because the feature belongs to the suite and the consumer will receive it through `just run-do-work-update`.
## Detailed Requirements
- Trigger and forms. The official trigger phrase is `ai-report second-opinion`, which parallels `do-work ai-report` (ai-report creates a report; ai-report second-opinion gives a second opinion on an existing one). Add the trigger words to the `skills/do-work-toolbox/SKILL.md` routing row for `ai-report.md` and to the action's description so any AI coder (Claude Code, Codex, Antigravity/agy, Gemini) picks them up. Forms, bundle defaulting to `latest` (newest bundle in the regenerated catalog, as the revise form resolves it): `ai-report second-opinion [<bundle-dir>|latest]` (session mode), `ai-report second-opinion --backend <cli> [--model <slug>] [<bundle-dir>|latest]` (targeted mode), `ai-report second-opinion --all [<bundle-dir>|latest]` (batch mode). Natural phrasings map to the same forms: "generate a second opinion" is session mode, "Codex, generate a second opinion" is targeted mode, "generate second opinions from all models" is batch mode.
- Action form in `skills/do-work-toolbox/actions/ai-report.md` (a "Second-opinion form" section beside the revise form) and a guide entry in `skills/do-work-toolbox/docs/ai-report-guide.md`. The form tells the calling agent: in session mode pass its own CLI and model slug (`--backend <own-cli> --model <own-slug>`), and stop with a message when its own model cannot be determined; in batch mode pass its own CLI as `--exclude <own-cli>`. The final reply lists the output file path(s), 3 to 5 bullets on clarity changes, the corrections list, and open questions, quoting each `<cli>-<model-slug>-reply.md`.
- Prompt file `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md`: the user's clarity-editor prompt in substance (role, input, using the repository, output, goal, rules 1 to 8, process a to e, final reply, as written in the UR input) with placeholders for the absolute bundle path, the absolute repository root, the output file name, the new-asset prefix `<cli>-<model-slug>-`, and the backlink requirement (the new page opens with a short "Original report: index.html" block linking to `index.html` by relative path; the model never touches `index.html`). Resolve the draft's contradiction between "runs offline" and "library from a CDN" by adopting the Report Design Rules in `ai-report-reference.md`: CDN allowed only for Tailwind CSS and Mermaid, everything else inline.
- Subcommand `ai-report-second-opinion` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go`, registered in `commands.go` beside `ai-report-judge`, invoked as `<skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format json ai-report-second-opinion [--backend <cli> [--model <slug>] | --all --exclude <own-cli>] [<bundle-dir>|latest]`. It resolves the bundle (a directory directly under `ai-reports/` holding `index.html`; `latest` as above), fills the prompt placeholders per run, and runs the modes: session mode is `--backend <own-cli> --model <own-slug>` with no fallback; with no `--backend` and no `--all` (a plain terminal with no active model) it uses the configured order codex, gemini, agy and runs the first one installed; targeted mode runs only the named backend and stops with a finding when it is not installed or the model slug cannot be resolved, never substituting; batch mode detects the installed backends among codex, gemini, agy, drops the excluded one, and runs them one at a time, each with its own post-run git status check, judge check, banner insert and reply file, reporting and skipping a failed or missing backend, continuing, and ending with a per-backend summary (ran, file name, judge verdict, or the reason it was skipped).
- Per backend: codex runs `codex exec -s workspace-write` with the working directory set to the bundle folder so the sandbox blocks writes elsewhere, never `--dangerously-bypass-approvals-and-sandbox`; gemini runs headless (`gemini -p ... -m <slug> --approval-mode auto_edit`); agy runs headless per `prime-agy-usage.md` (`-p`, explicit `--model`, `--print-timeout`, resolved model verified from the stream's `init` event); claude runs headless (`claude -p --model <slug>`) with the bundle as working directory when it is the session's own CLI. gemini, agy and claude have no write fence and rely on the post-run check. The model slug is read from the backend that actually ran, and the subcommand writes `index-<cli>-<slug>.html` itself; the model never chooses the file name. An existing file of that name is never overwritten; a numeric suffix is added. New assets the model creates must carry the `<cli>-<slug>-` prefix. The backend's final reply is saved as `<cli>-<slug>-reply.md` in the bundle.
- After every backend run, `git status --porcelain` must list no changed or new path outside the bundle folder; otherwise that run fails with a finding naming the paths. Then the new page is render-checked by reusing the judge in-process (`judgeRender`) on a scratch copy where the page is `index.html` together with the assets it references; the verdict is reported and the scratch copy removed. `skipped` and `error` are reported as such, never as pass.
- Banner in the original, inserted by the subcommand only, after the git status check passed and the judge verdict for that page is `pass`: one small banner block saying a second-opinion version is available, with a relative link to `index-<cli>-<slug>.html`, wrapped in `<!-- second-opinion-banner:<file name> -->` ... `<!-- /second-opinion-banner -->` markers and inserted right after the opening `<body>` tag. Markers are per file name, so `index.html` may carry one banner per second-opinion page; when a marker for that file name already exists the subcommand skips it and says so. When either check fails, `index.html` is left untouched and the result says why. The banner is the only change the subcommand may make to `index.html`.
- Tests in `report_second_opinion_test.go` pin: mode selection and the no-substitution refusal, file-name derivation from a fake backend's reported model, the numeric suffix on collision, the outside-bundle git status refusal, banner insertion and idempotency per file name, and the batch summary with one skipped backend. Fake backend executables on a temporary PATH stand in for codex, gemini, agy and claude.
- Shipped-file hygiene: justify the form in the action's description blockquote, keep every cross-package citation at the literal relative depth, run `_dev/tests/contract-regressions.sh`, `_dev/tests/shipped-package-reference-contract.sh` and `just maintainer-verify`, and release per `_dev/primes/prime-releases.md`.
- Prove it on `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/`. Session mode from Claude Code (`--backend claude --model <claude-slug>`) produces `index-claude-<slug>.html` there; targeted mode (`--backend codex`) produces `index-codex-<slug>.html`; batch mode from Claude Code (`--all --exclude claude`) produces one page and one reply file per installed backend among codex, gemini and agy and prints the per-backend summary. In all three: `git status` shows only those new files inside the bundle plus the banner change in `index.html`; `git diff -- <bundle>/index.html` shows only marker-wrapped banner blocks, one per produced page; a second run reports every banner as already present and leaves `index.html` unchanged; every judge verdict is recorded in the Testing section.
## Constraints
- The model never modifies `index.html`, existing assets, or anything outside the bundle; `index.html` stays byte-identical except the subcommand's marker-wrapped banner blocks.
- Session and targeted modes never substitute another backend or model; the codex, gemini, agy order applies only to the terminal fallback with no backend given.
- Read `skills/do-work-toolbox/docs/prime-agy-usage.md` before wiring the gemini or agy invocation.
- Follow `_dev/primes/prime-action-files.md` for the action form and router row, `_dev/primes/prime-shell-commands.md` for any prescribed shell block, and `_dev/primes/prime-releases.md` for the release commit.
## Dependencies
None. `ai-report-judge` already exists in the do-work CLI and is reused in-process.
## Builder Guidance
High certainty on the trigger, the three modes, the file naming, the banner rule and the checks; they were decided with the user over four rounds in the consumer repository. Challenge and resolution, recorded here because the user is not watching: the consumer REQ asked for a bash launcher in a project-owned skill folder; upstream, the suite's rule that programs beat prose for anything mechanical and the fact that the judge already lives in Go make a do-work-cli subcommand the right shape, with the action form carrying only judgment. Latitude on how each backend's resolved model is read (codex: `codex exec --json` session header; gemini: its JSON output; agy: the `init` event; claude: its JSON output), on how installed backends are detected, and on the banner markup. Session mode from Claude Code implies `claude` is a runnable backend even though the user's batch set is codex, gemini, agy; keep the batch set as written. If `claude -p` cannot run headless with a chosen model here, document the finding and prove session mode from another installed AI coder.
## Red-Green Proof
**RED prompt/case:** `skills/do-work/tools/do-work-cli.sh --repo-root <root> --format json ai-report-second-opinion --backend codex ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism` (targeted), the same with `--backend claude --model <claude-slug>` from a Claude Code session (session mode), and with `--all --exclude claude` (batch); plus `go test ./internal/toolboxcommands/ -run SecondOpinion` in `skills/do-work/tools/do-work-cli`.
**Why RED now:** the subcommand, the action form, the prompt file and the tests do not exist; the CLI answers unknown command and no second-opinion page can be produced.
**GREEN when:** the tests pass; targeted mode exits 0 and `index-codex-<slug>.html` exists, opens, and starts with a backlink to `index.html`; session mode from Claude Code produces `index-claude-<slug>.html` using Claude Code's own model; batch mode from Claude Code produces one `index-<cli>-<slug>.html` and one `<cli>-<slug>-reply.md` per installed non-calling backend and ends with the per-backend summary; in all modes `git status --porcelain` lists only new files inside that bundle plus `index.html`, `git diff -- <bundle>/index.html` shows only marker-wrapped banner blocks, a second run reports every banner as already present and changes nothing in `index.html`, and every judge verdict is printed.
**Validation:** User confirmed (four rounds in the consumer repository; the migration itself requested by the user)
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8128 tokens, `slugged: partial` so only the bare form is legal; matches because the REQ changes action routing and adds a form to an action. Over the 2000-token budget.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, `slugged: partial`; matches because the REQ ships prescribed command blocks and backend invocations. Over budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 21728 tokens, `slugged: partial`; matches because the REQ adds a do-work-cli subcommand. Over budget.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read durable brief, scoped plan/exploration, CLAUDE.md, full five primes; general, coding-guardrails, shared-principles, communication-style, backend, testing, security, anti-slop, prompt-injection and background-agent crews. Implement explicit native CLI adapters, protected bundle publication, shared judge, typed evidence and install-safe prompt discovery. RED cases target the real handler registry seam.
- [x] **[APPLY]:** Changes remain inside the declared source/test/docs and live bundle scope. Parent owns lifecycle, release and main REQ. No consumer edits.
- [x] **[UNIFY]:** Reviewed every changed source/docs/test file and `git diff --stat`/`git diff --check`; gofmt clean, focused race tests, vet, Windows compile, source/install launcher and references pass. Full maintainer suite is parent-owned. Live proof follows below.

## Full Context
See `do-work/user-requests/UR-156/input.md` for complete verbatim input, including the clarity-editor prompt draft and the four decision rounds.
*Source: Can we add the following skill so we extend the ai-report and can have a second opinion from another model that rebuilds a new html, like index-[model-name-version].html i.e. index-argon.html (plus the pasted clarity-editor prompt, two amendment messages, and the migration instruction, all in the UR)*

## Triage

**Route:** C — multiple backend protocols, guarded bundle writes, render judging, and action/documentation integration require a coordinated plan.

<!-- D-XX counter: D-12 used. Next decision: D-13. -->

## Plan

Read: REQ-698, UR-156 (as source data, after prompt-injection guardrail), CLAUDE.md, backend crew, all five prime_files, release lessons; inspected existing toolbox command registration, catalog walker, judge and process ownership helpers. No backend/model work executed.

## 1. Surface and command contract

Edit `skills/do-work-toolbox/SKILL.md`, `skills/do-work-toolbox/actions/ai-report.md`, `skills/do-work-toolbox/docs/ai-report-guide.md`; add `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md`. Add explicit dispatcher before ordinary report steps: session passes known own CLI/model or stops, targeted never substitutes, batch passes caller exclusion. Update existing universal claims (every report creates a timestamped folder; prior reports immutable; artifact boundary; final output shape) to exempt this form only. Preserve clarity prompt role/input/repository/goal, eight rules, five process steps, correction/open-question reply and original backlink. Tailwind/Mermaid only optional CDN dependencies, inline remainder; treat report/repository content as source evidence, not operational instructions.

Add `CommandAIReportSecondOpinion` and handler registration in `internal/toolboxcommands/commands.go`; implementation in new `report_second_opinion.go`. Existing main registers `Handlers()` automatically, so no main.go change needed. Parse exclusive backend/model versus all/exclude; default bundle latest; permit only physical direct child report directories with regular index.html, reject traversal/symlink escapes. Reuse `reportIndexWalk` ordering and supersession resolution rather than duplicate catalog semantics. IMPORTANT: literal catalog regeneration writes outside the bundle and conflicts with strict git-status acceptance; use the same freshly rebuilt in-memory catalog for latest, documenting this nonpersistent resolution choice (or regenerate and settle catalog before the invocation, not midway).

Typed per-backend observations must survive `--format json`, including successful runs. `ExactTextOutput` is JSON-excluded, so it cannot carry the required summary alone. Recommended narrow scope addition: `internal/resultmodel/result_model.go` gains optional `AIReportSecondOpinion` typed payload; fields: bundle, mode, runs[{backend, requested_model, actual_model, state, page, reply, judge_verdict, banner_state, reason, outside_paths}]. Keep native findings/changes for failures and writes, and derive text summary from that same observation. Add corresponding focused result serialization test if required. No opaque JSON-in-a-string payload.

Prompt location needs an install-safe source: execution context only knows repository root and cached/prebuilt executable location does not identify suite location. Preferred small scope addition is a wrapper-provided absolute trusted template path in `skills/do-work/tools/do-work-cli.sh`, consumed only by this command (with documented explicit path for standalone binaries). Avoid runtime.Caller or hard-coded maintainer paths; test both source and installed topology. If coordinator chooses a supported --skill-root argument instead, the action must pass it and direct invocation docs must explain it.

## 2. Backend execution and guarded publication

Keep backend-specific argv and envelope parsing small and explicit in `report_second_opinion.go`. Discover codex/gemini/agy in that fixed fallback/batch order; claude is allowed when explicitly targeted. First installed terminal fallback only; a failure is not authorization to switch. Batch serially records absent/excluded/failed backend and continues. Use argument arrays, bundle cwd, stdin closed, separate stdout/stderr, bounded cancellation with proven owned-process cleanup. Reuse existing ownedprocess machinery rather than shell command strings. Codex uses exec -s workspace-write, never bypass flags; Gemini uses -p, -m and auto_edit; AGY uses explicit catalog-resolved latest Flash Low/Medium default and print-timeout, stream init plus SUCCESS terminal event; Claude uses -p --model and structured output. Validate installed protocol through read-only help/model probes, then a bounded live protocol observation before defining fake adapters. Coordinator checked all four CLIs available; Codex supports --json, --output-schema, --output-last-message and --ephemeral. Official documented Codex thread.started carries only thread_id, so do not build a fake model header that the real CLI never emits. For Codex inspect actual stderr session header or a session turn_context tied to the returned thread_id; require verifiable run-specific provenance and avoid unrelated-session/default-config inference.

Resolve a backend model before constructing model-dependent asset prefix, then verify actual run identity from native structured backend metadata. Never use model-authored prose or merely the requested flag as proof of actual identity. Fail with specific finding on absent identity, ambiguous identity, failed/truncated result, or explicit model substitution. Alias normalization must be backed by backend metadata, not a guessed mapping. The REQ's suggested Codex JSON header is an assumption to verify, not permission to synthesize one; if installed Codex cannot report actual model, report that backend limitation and retain refusal behavior.

Launcher chooses a unique staging output name and publishes it exclusively under `index-<backend>-<actual-slug>[-N].html` after validating metadata. Model never chooses destination. Use paired collision suffix for reply too (`<backend>-<slug>[-N]-reply.md`), so repeat generations cannot overwrite prior replies. Reject unsafe model slug path characters; don't silently collapse different identities onto the same normalized name. Require newly written assets to use verified prefix and preserve all pre-existing bundle files byte-for-byte and mode/type. Snapshot existing bundle files before each run: git status alone cannot detect forbidden changes to an already-dirty original or ignored assets. Refuse symlink outputs and guarded publication races; never auto-revert unrelated writes.

Strictly inspect `git status --porcelain -z --untracked-files=all` before and after each backend; any outside-bundle path (including both rename sides) is a refusal naming paths. Early preflight avoids wasting a model run on an already-dirty tree. Do not weaken to status deltas: a dirty tracked file can change again without a status-code change. Failure leaves original unchanged by launcher and clearly records any backend violation.

Render on an owned scratch copy with generated page as index.html and required local assets. Use discoverJudgeBrowser and judgeRender in process; retain structured verdict/findings before deleting scratch and captures. Judge pass alone does not validate original backlink semantics because scratch index.html points to the candidate: explicitly check the candidate backlink targets the real original and exists. Validate original bytes still match pre-run observation before injecting the sole permitted marker-wrapped banner immediately after body. A missing body or changed original is a finding; no rewrite/formatting. Skipped/error/fail never inserts banner.

Resolved conflict: full generation reruns must allocate a new numeric name, hence legitimately append a new banner. “Second run changes nothing” applies to reinserting a banner for an already-produced filename, not invoking a fresh generation. Test per-filename idempotency separately and record this smallest faithful interpretation; adding a reuse cache or silently skipping generation would conflict with collision semantics.

## 3. Evidence, integration and release

Write RED tests first in new `report_second_opinion_test.go` with fake executables on isolated temporary PATH. Pin modes, no backend/model substitution, actual-model metadata determines naming, unknown/mismatching model refuses, collision preserves page and reply, strict outside dirt preflight/postflight (including rename/untracked), unchanged pre-existing assets/original, unsafe paths/symlinks, required backlink, judge fail/skipped/error leaving original untouched, exact banner insertion and per-file idempotency, and serial batch summary with absent/failed backend plus successful remaining run. Use deterministic fake judge seam for fast state-machine tests; existing real judge/browser tests cover rendering, with one integration proof using a real browser.

Run focused SecondOpinion RED then GREEN, gofmt/go vet, package tests, `_dev/tests/contract-regressions.sh`, `_dev/tests/shipped-package-reference-contract.sh`, and prescribed maintainer gate. Use ordinary heavy-lane plan based on exact implementation diff; no duplicate full suites without changes.

For live proof on `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/`, first commit implementation and perform runs from a clean isolated checkout; claimed REQ metadata in the coordinator otherwise correctly fails strict outside-bundle acceptance. Capture explicit session backend/model proof only for a model genuinely known to the calling session; do not pretend this Codex session is Claude Code. Exercise targeted Codex and batch with actual available engines, preserving every unavailable/refused/judge verdict. Import only owned bundle outputs after checking their diff; no broad add. Verify original diff consists solely of banners and artifact links open. Required live backend limitations remain documented findings rather than fabricated passes. Finalize with release prime/mirrors and canonical lifecycle; parent owns release/archive/cleanup.

**Plan validation:** All captured requirements map to the three steps. The typed JSON result and install-safe prompt handoff are necessary wiring for the requested JSON command and installed suite, not additional product behavior. Native Codex text-mode stderr was observed reporting its actual model; JSON mode did not, so use the verified header instead of an invented JSON fixture.

*Generated by Plan agent; validated by orchestrator*

## Decisions

- **D-01 — DECIDE & STATE:** Banner idempotency applies to a given generated filename. A fresh generation uses the required numeric suffix and receives its own banner; this reconciles the collision rule with the contradictory whole-rerun no-change sentence.
- **D-02 — DECIDE & STATE:** Preserve the strict outside-bundle Git-status check and perform live acceptance in the clean isolated builder checkout. Resolve latest through the existing catalog walker in memory, avoiding a catalog write outside the permitted bundle.
- **D-03 — DECIDE & STATE:** Extend the declared implementation scope to the typed result model, its focused serialization test if needed, and the launcher for install-safe prompt discovery. These files are necessary for the requested structured CLI output and shipped prompt.

- **D-04 — DECIDE & STATE:** Extend scope for the exact handler-count test, launcher source/install test, shared report reference restatements if needed, and model-prefixed live assets required by the original proof. These are demonstrated integration seams, not extra features.

- D-05 — DECIDE & STATE: Use one short native read-only model-resolution probe (90-second bound) before generation, then require identical native model metadata on generation. AGY uses its zero-turn `/model` command. Codex native text startup header, Claude modelUsage and Gemini stats.models provide actual identity; absent/ambiguous/mismatched identity refuses. A default CLI model is never represented as this chat's active model.
- D-06 — DECIDE & STATE: Models fill an unpredictable launcher-reserved staging leaf. The launcher publishes the verified model-named final page/reply exclusively and maps that exact staging leaf to the final page name inside the saved reply. New assets remain directly in the bundle under the verified backend/model prefix. Captures use disposable scratch outside the repository and typed findings survive cleanup.
- D-07 — DECIDE & STATE: Reject symlinks/special files in report inputs and new outputs; snapshot every pre-existing bundle file's inode, complete mode and bytes. Do not auto-revert a backend violation. Preserve strict pre/post Git status, including both rename paths, and report named violations.

- D-08 — DECIDE & STATE: Save relative published page links in replies by replacing the exact absolute reserved stage path before its leaf name. Prompt source-map references use repository-relative text unless a canonical remote link is independently known to be published; a local commit and configured remote are insufficient. Applied only to new generated artifacts, with disclosure and a separate post-edit render judge.
- D-09 — DECIDE & STATE: Backlink validation accepts reordered anchor attributes, unquoted hrefs and harmless inline formatting, but rejects duplicate/wrong hrefs, misleading href text inside other attributes, wrong labels and backlinks following main content. A focused regression proved the old parser too strict; the original AGY candidate was not retained, so its failure cannot be confidently attributed to model noncompliance.

- **D-10 — DECIDE & STATE:** Add the four concrete generated artifact paths to Scope and write_set as expansions of the predeclared globs. The scope comparator is literal; this records existing authorization without broadening it.

- **D-12 — DECIDE & STATE:** Extend the parent-owned completion scope to the CLI lesson satellite, its generalized prime trap and lessons index. The confirmed recurrence of reaped-by-its-own-parent requires a prime amendment under general.md; native model protocol evidence is a new targeted lesson. Defer actual lesson writes until the finalization plan establishes the archive target.

## Exploration


No repository edits and no model turns. Read the plan and full claimed REQ. The planned shape fits existing ownership; keep this one command plus a typed result and small wrapper path contract, not a backend framework.

## Exact reuse seams

- `internal/toolboxcommands/commands.go`: add one constant/handler; main already registers `Handlers()`. **Also update `commands_test.go`**: explicit constant list and exact count 9 currently make registration of the new command fail.
- `report_index.go:81 reportIndexWalk(repo)`: rebuilt in-memory catalog, excludes hidden folders and symlink directory entries, sorts date descending then path ascending; fills `Supersedes`/`SupersededBy` with newest successor. Reuse for latest and forward-chain resolution, with cycle detection. It admits folders without index.html, so validate chosen physical direct child and regular index independently. Do not call `reportIndexWrite`: it dirties catalog/index outside the allowed bundle.
- `report_judge.go:159 judgeRender(&report,bundle,out,browser)`: owns server/browser lifetime, fills captures/findings, returns operational error. It **does not assign verdict**. The new command must set skipped if browser discovery empty; error on returned error; fail if findings exist; pass only otherwise. `discoverJudgeBrowser()` already finds PATH Chromium and macOS Chrome. Retain findings in the new typed payload before deleting scratch; merely retaining verdict loses actionable evidence. Browser checks only same-origin links, so a separate required original backlink check is needed. Existing real-render tests are gated by DO_WORK_HEAVY_TESTS and testing.Short.
- `mutation.go:132 rootedPublishFile(repo,relative,bytes,mode,replace)`: existing no-linked-ancestor, pinned confined-parent, exclusive hard-link creation and inode/content-checked replacement. Reuse for final page/reply publication; for banner, first compare original against the pre-backend bytes/mode/type, then reuse replacement. Its internal replacement snapshot starts at function invocation, so it does not replace that pre-backend preservation check.
- `report_image.go:336 mutationSignalContext()`: existing HUP/INT/TERM cancellation with conventional exit status. Reuse and set ExitCodeOverride when interrupted.
- `report_image_process.go`: `configureOwnedProcess(*exec.Cmd)` establishes group or fails closed; `completedProcessResult(err)` decodes exit. Existing `runOwnedProcess` discards stdout/stderr, so cannot establish model identity. A narrow captured-output runner should use those helpers plus `ownedprocess.TerminateGroup(pid,grace)`, wait before scratch cleanup, bound total runtime, close stdin. Do not use CombinedOutput because metadata and authored output have different authority.

## Template and serialization

- Existing source/install pattern: `knowledgecommands/interview_commands.go:680 locateKnowledgeRoot` accepts explicit root, otherwise checks `.claude/skills/do-work-knowledge` then `skills/do-work-knowledge` under repo. This is useful precedent for explicit override, but misses user-installed skills unrelated to repository. The compatibility wrapper knows its physical script_dir before both go-tool and prebuilt exec branches; export one absolute template path derived from `script_dir/../../do-work-toolbox/docs/...` there. Tests must exercise installed layout as well as source. Avoid runtime.Caller and executable-relative lookup: Go cache and prebuilt cache relocate binary.
- `resultmodel.CommandResult` already supports optional typed command payloads (AuditMetrics, HeavyVerification, etc.). Add SecondOpinion pointer and typed per-run fields. `ExactTextOutput` is JSON-excluded; text must derive from same payload. Include model requested/actual, state/reason, page/reply, judge verdict and findings, banner state, outside paths. Normalize optional slices to [] when payload present; JSON serialization test should assert complete batch including missing/failed/successful rows. Existing `result_model_test.go:742` tests typed AuditMetrics plus exact text without opaque duplication.

## Tests and integration

- Reuse `toolboxTestRepository`, `toolboxTestGit` in report_image_test.go; fake shell executable files mode0755 and package seam restored with t.Cleanup are established convention. The requested temporary PATH executables should be outside fixture repo to avoid strict outside-dirt preflight refusing the test backend itself. Keep real git available in temporary PATH or resolve git before narrowing it.
- `report_image_process_test.go` covers cancellation and leader-first death hazards. Add a focused captured-runner cancellation test if introducing a new runner rather than blindly assuming reused group setup proves teardown.
- Cover both sides of porcelain rename rows using -z --untracked-files=all; never compare only pre/post status strings. Snapshot all existing bundle files (including ignored files); git cannot prove preserved bytes. Batch must take a fresh snapshot per backend after earlier legitimate banners/assets.
- Use one fake judge seam for deterministic fast run states, plus one real browser integration. Judge skip/error/fail must all leave original bytes unchanged.
- Adjacent markup can impersonate markers/body tags; restrict banner identification/insertion to real HTML context enough to avoid a comment/string marker appearing as an existing banner. Idempotency is per already-produced filename, while a new generation uses collision suffix and legitimately adds a new banner.
- Broaden write set only for demonstrated seams: `commands_test.go`, resultmodel files, wrapper and its source/install contract test if using wrapper path. No need to touch main.go, catalog implementation, or judge implementation.

## Restatement sweep sites

`actions/ai-report.md`: description line3, principle11, argument dispatcher32, boundary141, output163, checklist179. `docs/ai-report-guide.md`: every-writing-invocation claim9, artifact boundary21, catalog immutable claim90, command forms70-73. These need explicit second-opinion exception. `actions/ai-report-reference.md` is the actual shared reference location (not docs/); check its collision/publication rules for operative restatements before final diff. A narrow second-opinion dispatcher must stop before ordinary report Steps1-8 and preserve normal revise behavior.

## Native identity constraint

Parent-provided observed evidence: Codex JSON mode does not emit model identity; text stderr startup header contains model/provider/session id. Keep metadata parser limited to native startup header before model-authored/tool output, do not search arbitrary whole stderr for model lines. A separate default-resolution probe is not sufficient proof of the execution model; validate the actual generation run too. Explicit slug mismatch is refusal, never substitution. Prefix construction before actual identity must use an observed resolved selection or a staged protocol; do not infer identity from a requested flag. No backend probes or model turns were launched by this exploration.

## Read evidence

Read fully: all five prime_files (action-files, shell-commands, releases, do-work-cli, agy usage), required `_dev/primes/lessons-releases.md`, `_dev/lessons/validated-runtime-boundaries.md`, REQ and supplied plan. Touch reads: action lessons Template/Agent Compatibility and alternate-writer/restated-mechanism families; shell lessons Prescribed Commands/Unchecked Exit/SIGPIPE/Every Flag; CLI lessons Package routing/Package direction/Scoped command contracts/Verify plus reaped-by-its-own-parent, final-boundary-identity, opaque-evidence-projection, smoke-vs-characterization, shipped-module-test-self-containment and time-is-not-authorization entries. Large initial satellite output was truncated, so do not record those whole satellites as fully read. Re-read relevant named sections in bounded reads. No release changes made.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/ai-report.md` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work-toolbox/SKILL.md` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work-toolbox/docs/ai-report-guide.md` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-*.html` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/*-reply.md` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work/tools/do-work-cli.sh` — command, documentation, deterministic validation, or required live proof.
- `_dev/tests/do-work-cli-launcher-behavior.sh` — command, documentation, deterministic validation, or required live proof.
- `skills/do-work-toolbox/actions/ai-report-reference.md` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-*` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/gemini-*` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/agy-*` — command, documentation, deterministic validation, or required live proof.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-*` — command, documentation, deterministic validation, or required live proof.

- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-claude-opus-5-5-reply.md` — exact expansion of the predeclared live-artifact glob.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-gpt-6-astra-reply.md` — exact expansion of the predeclared live-artifact glob.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html` — exact expansion of the predeclared live-artifact glob.
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html` — exact expansion of the predeclared live-artifact glob.

- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — parent-owned post-review lesson recurrence update.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — parent-owned post-review lesson recurrence update.
- `do-work/lessons-index.md` — parent-owned post-review lesson recurrence update.

**Files I will NOT touch:** Other report bundles, catalog output, unrelated CLI commands, consumer repositories, release/lifecycle metadata (orchestrator-owned).

**Acceptance criteria (restated from REQ):**
- [ ] Router/action/guide expose second-opinion session, targeted, terminal fallback and sequential batch forms; preserve original report creation/revision.
- [ ] Full clarity-editor prompt preserves source rules, steps and reply, and substitutes trusted absolute repository/bundle, launcher filename, prefix and backlink.
- [ ] Resolve latest with existing catalog order; accept only a physical direct-child bundle with regular index.html.
- [ ] Session/targeted never substitute; unknown session model stops; terminal selects first installed codex/gemini/agy once; batch excludes caller and continues through missing/failing engines.
- [ ] Invoke backend headlessly with required flags and bundle cwd, derive actual identity from native run metadata, verify explicit model selection.
- [ ] Launcher publishes unique model-derived page and paired reply; new assets prefixed; preserve every existing bundle file.
- [ ] Strict Git status before/after names all outside-bundle dirt; no banner on failure.
- [ ] Reuse judgeRender on disposable scratch, report pass/fail/skipped/error and findings; only pass allows launcher-only marker banner immediately after body.
- [ ] Per-filename banner insertion is idempotent; fresh generation allocates numeric suffix.
- [ ] Fake executable tests prove selection, identity, collision, outside dirt, preservation, judge states, banner and batch behavior with observed RED/GREEN.
- [ ] Live proof uses the specified REQ-632 bundle with truthful backend/model/session and judge evidence; original diff is banner-only.
- [ ] Contract/reference/maintainer gates pass and release follows prime.

## Pre-Flight

**Git:** Only orchestrator-owned REQ metadata changed; builder checkout clean at 61ca1e9c.
**Tests baseline:** `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run 'TestAIReport|TestHandlersRegister' -count=1` launched and passed.
**Canonical gate:** `bash _dev/tests/maintainer-verify.sh` exit 0, timed baseline 126 seconds. `advance` returned satisfied preflight and green-gate records for this exact REQ/path at 61ca1e9c7ca9090fae5b227aec3aecacfdd341fe.
**Dependencies:** All four backend executables found. Gemini live authentication currently refuses UNSUPPORTED_CLIENT; preserve this as an honest batch failure rather than substitute another backend.

*Checked by work action*


## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified; parent-owned recurring process-lifetime trap)
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified; parent-owned native identity and process-lifetime lessons)
- `_dev/tests/do-work-cli-launcher-behavior.sh` (modified)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-claude-opus-5-5-reply.md` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-gpt-6-astra-reply.md` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html` (new)
- `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html` (modified)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (new)
- `skills/do-work/tools/do-work-cli.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (new)

**What was done:** Added the shipped second-opinion action and Go command with explicit native model verification, serial batch results, strict bundle preservation, supervised process-group completion, unique page/reply publication and pass-only original banners. Published Codex and pinned Claude proofs; the Codex source-map links received a disclosed verifier correction and separate post-edit judge. The original differs only by two banners; existing evidence files are unchanged.

**Integration:** Captured first base `3b593d22`; initial merge `d0106c31`, final remediation merge `31e7fcdf`. Builder commits `406ff3c4`, `9fa76500`, `2f5cd87a`, `69a2a998`, `2abd7939`; no integration seams or do-work edits on builder branch.

## Discovered Tasks

None outside this request.

## Builder Evidence

Read `_dev/primes/lessons-releases.md` whole-satellite. No missing required entry.
Touch-conditional reads: `_dev/primes/lessons-action-files.md` Template, Agent Compatibility, alternate-writer-contract-drift/restated-mechanism families; `_dev/primes/lessons-shell-commands.md` prescribed commands, unchecked exit, SIGPIPE, every-flag caller sections; `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` Package routing/direction, Scoped command contracts, Verify, final-boundary-identity, opaque-evidence-projection, reaped-by-its-own-parent, smoke-vs-characterization, shipped-module-test-self-containment families. Large initial combined outputs were truncated; the relevant sections/families were subsequently read in bounded reads. Also read `_dev/lessons/validated-runtime-boundaries.md` fully. Do not claim entire large satellites were read.

- RED: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run SecondOpinion -count=1`: all initial cases failed at `ai-report-second-opinion command is absent` (assertion failure, not compile failure), `/private/tmp/req698-evidence/builder-red.log`.
- GREEN implementation focused PASS7.084s (`builder-publish-test.log`), reply portability focused PASS7.595s (`builder-reply-portability-test.log`).
- Backlink regression RED: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run TestSecondOpinionBacklinkAllowsSemantic -count=1` failed three valid-equivalent examples; `builder-backlink-red.log`. GREEN: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands ./internal/resultmodel -run SecondOpinion -count=1` PASS7.375s/0.363s (`builder-backlink-green.log`).
- Final race package checks: `go -C skills/do-work/tools/do-work-cli test -race ./internal/toolboxcommands ./internal/resultmodel`, PASS toolbox11.485s/resultmodel2.010s (`builder-final-recheck.log`), after backlink and naming edits. Prior race checks also passed.
- `go -C skills/do-work/tools/do-work-cli vet ./...`: final exit0 (`builder-final-vet.log`).
- `GOOS=windows GOARCH=amd64 go -C skills/do-work/tools/do-work-cli test -c ./internal/toolboxcommands -o /private/tmp/req698-evidence/toolbox-windows.test.exe`: exit0.
- `bash _dev/tests/do-work-cli-launcher-behavior.sh`: exit0, source and installed prompt topology fixtures pass.
- `bash _dev/tests/shipped-package-reference-contract.sh`: exit0, PASS.
- `bash _dev/tests/contract-regressions.sh`: first restricted-sandbox attempt exit1 with unavailable cached executables and `ps: Operation not permitted`; exact retry with escalation exit0, `Contract regression checks passed`. No cache implementation changes.
- Fake executable coverage: all4native adapters; terminal chooses once; no backend/model substitution; generation model switch; filename/reply collisions; latest/direct-child/symlink refusal; pre/post outside dirt and rename source; asset/original preservation; required backlink; pass/fail/skipped/error banner gating; serial batch missing/excluded/failed states; typed JSON findings; cancellation; actual-body banner and per-filename idempotency; required headless flags.

### Live backend evidence

All commands use `skills/do-work/tools/do-work-cli.sh --repo-root /private/tmp/worktree-agent-REQ-698-ai-report-second-opinion --format json ai-report-second-opinion` and bundle `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism`.

| Mode / arguments | Native identity | Outcome | Judge | Banner | Outside paths |
| --- | --- | --- | --- | --- | --- |
| Targeted `--backend codex` | `gpt-6-astra` | Exit0; page/reply published | pass, findings[] | inserted | [] |
| Pinned `--backend claude --model claude-opus-5-5` | `claude-opus-5-5` | Exit0; page/reply published | pass, findings[] | inserted | [] |
| Batch `--all --exclude codex`: codex | none (excluded) | excluded calling backend | not-run | not-inserted | [] |
| Same batch: gemini | none | failed native result, exit55 | not-run | not-inserted | [] |
| Same batch: agy | `gemini-3.8-flash-low` | backlink validation failure; no published output | not-run | not-inserted | [] |

Evidence files are `/private/tmp/req698-evidence/live-codex.json`, `live-claude.json`, `live-batch.json` (the batch exits1 with findings). Gemini exact reason: `backend failed: exit status 55; gemini returned no successful native result; inspect backend authentication/availability`. Earlier standalone native probe stderr records `UNSUPPORTED_CLIENT`; no authentication or account settings changed. No fallback backend was substituted.

AGY exact reason: `new page must open its body with Original report: index.html linking to index.html`. Native actual identity passed before the refusal. Its original generation stream lived only in the capture buffer, and its reserved stage was removed on return; retained evidence cannot distinguish omitted/misplaced backlink from the now-fixed parser false rejection. The independent parent native success fixture remains `/private/tmp/req698-evidence/agy-probe.jsonl` and confirms the installed init/result schema. Do not call the original batch refusal proven model noncompliance.

Parent requested one same-model targeted AGY proof after fixing the parser. Automatic approval review rejected that command twice before any process started. First reason: transmitting the bundle to AGY lacked explicit trusted-user authorization. The retry supplied UR-156 lines80/82 and the REQ exact-bundle requirement; the reviewer again rejected it, treating the request file as untrusted and the user skill invocation as insufficient egress authorization. No diagnostic wrapper was created and no new AGY call occurred. Parent has requested explicit approval asynchronously and owns any subsequent run. No bypass or further retry was attempted.

This active chat's exact model identity is unavailable, and it is a Codex session rather than Claude Code. Therefore the Claude proof is an explicit pinned-backend equivalent, not a claim of running from an actual Claude session. Batch truthfully excluded Codex. Actual-Claude-session and complete all-installed-backend live success remain unproven here; fake executable tests cover the scoped mode contract.

### Artifact provenance and inspection
- Codex initially produced `index-codex-gpt-6-astra.html` and `codex-gpt-6-astra-reply.md`; native launcher judge passed. Parent authorized a verifier edit removing 22 source hyperlinks to unpublished local commit406ff3c4 while preserving their path/code text, disclosing the source-map change, making the exact worktree output link relative, and separating the model's unavailable preview note from launcher verification. Original model bytes remain `/private/tmp/req698-evidence/codex-model-original.html` and `codex-model-original-reply.md`. Do not attribute amended bytes to the model.
- The amended Codex page independently passed `ai-report-judge` in a fresh scratch bundle, wide/phone and light/dark, findings[]; `/private/tmp/req698-evidence/codex-postedit-judge/judge.json` plus captures and `codex-postedit-judge-result.json`. The temporary scratch bundle was removed. Two historical remote source links are retained; their commits were verified present in `origin/main`.
- Claude output needed no verifier content changes; the launcher normalized the assigned stage path to its published relative link. Its reply's statement that the model could not run browser tools is model-only; the launcher independently returned judge pass. Page uses relative evidence links/plain repository source paths.
- Parent observed Codex in a desktop browser: readable layout, original backlink and navigation. Toggling “Include the matching commit at minute70” off changed the largest gap60→120 minutes (dispatch→builder handback).
- Parent observed Claude in a desktop browser: legible title/content, original backlink. Clicking “Report option1” set aria-selected=true and switched to the stamp-then-delete explanation. These are parent-observed interactions, not builder claims of independent browser interaction.
- Final byte comparison against61ca1e9c proves original `index.html` differs only by exactly two marker-wrapped banners; all three pre-existing evidence files are byte-identical. Reply file links resolve within the bundle and contain no worktree/staging path. Evidence: `/private/tmp/req698-evidence/artifact-invariants.json` (also includes final artifact SHA256s).
- Fresh invocations intentionally use numeric suffixes and get distinct banners. Per-filename banner idempotency is covered by focused tests; no second expensive fresh run was misrepresented as “already present.” Parent D-01 resolved that request contradiction.

## Authorized AGY Verification

Trusted authorization relayed by parent: user answered **“Allow this one Antigravity verification run”** to **“May I rerun Antigravity (model gemini-3.8-flash-low) once to verify the backlink fix, sending the tracked REQ-632 report bundle and relevant repository source to that model?”**

Executed exactly one invocation at clean HEAD `69a2a998cfbf21391fd2a9c7c3cff4c7718da6d3`, with the same wrapper command/bundle above and `--backend agy --model gemini-3.8-flash-low`. A transparent scratch PATH wrapper called the installed `/Users/t2/.local/bin/agy`, preserving stdout, stderr and the reserved candidate before launcher cleanup. It changed no tracked source. Native probe and generation both resolved the explicitly requested model; generation emitted terminal `SUCCESS` for conversation `39bcf718-8719-4f38-8788-9543c0320b35`.

**Launcher result: exit2 / failure; state failed, actual_model gemini-3.8-flash-low, judge not-run, banner not-inserted, page/reply empty, outside_paths[].** Exact reason remains `new page must open its body with Original report: index.html linking to index.html`.

This candidate is now diagnosable: immediately after `<body>` it opens a container and sidebar headed `Report Navigation`, followed by a full table of contents. Only after that sidebar does `<main class="main-content">` begin with the correct `<p>Original report: <a href="index.html">index.html</a></p>`. The hyperlink is correct, but preceding navigation violates the prompt's required body-opening backlink. This is a real candidate ordering violation for the approved run; it does not retrospectively establish the cause of the earlier unretained batch candidate.

Evidence retained:
- `/private/tmp/req698-evidence/live-agy-approved.json` — typed launcher result.
- `/private/tmp/req698-evidence/agy-approved-probe.stdout` and `.stderr` — native zero-turn model resolution.
- `/private/tmp/req698-evidence/agy-approved-generation.stdout` and `.stderr` — full native generation stream (60,472 bytes stdout, empty stderr).
- `/private/tmp/req698-evidence/agy-approved-candidate.html` — exact model candidate, never published or amended.
- `/private/tmp/req698-evidence/agy-approved-before.json`, `agy-approved-original-before.html`, `agy-approved-invariants.json` — byte comparison and opening markup.

All eight pre-existing bundle files are byte-identical before/after the run, including the original with its two existing banners. Git status remains completely clean; no candidate, reply, asset or new banner was published. No new commit was needed; final HEAD remains69a2a998. No source changes and no second model retry were made. This completes the explicitly authorized single verification attempt.


## Remediation

The first independent review found a confirmed delayed-write race (F1). Although its policy mapping was Partial/report-only, the orchestrator chose to repair this current-request defect before finalization. Diagnosis confirmed: command.Wait reaped the backend leader before group ownership was checked, so a redirected background child survived and wrote after capture returned success. The overlay probe at `/private/tmp/req698-review/review_probe_test.go` failed in 0.66s. One remediation cycle is assigned with debugging guidance, the existing reaped-by-its-own-parent lesson, and the exiting-leader regression as closure evidence. No further external model calls are authorized or planned.

### First independent review record

### Review before remediation

**Overall: 80.625%** | 2026-10-15T12:58:20Z

| Dimension | Score |
|-----------|-------|
| Requirements | 87.5% |
| Code Quality | 85% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1: `report_second_opinion.go:582–589` reaps the backend leader before group cleanup; `owned_process_group_unix.go:50–54` then cannot establish group ownership and leaves redirected background writers alive. Independent temporary overlay probe returned success before a delayed file write; repair the runner's lifecycle before finalization. — impact-user-visible → report only
- F2: REQ live acceptance remains incomplete: targeted Codex and pinned Claude passed, actual-session proof is absent, Gemini failed exit 55, and authorized AGY retry correctly refused a backlink after navigation with judge not-run and no publication. — impact-user-visible → report only

**Minor findings:** None.
**Acceptance:** Partial — local and merged integration checked, live Codex/pinned-Claude success plus batch/AGY refusal paths observed; successful actual-session/full-batch acceptance incomplete; consumer deployment unassessed.
**Restatement sweep:** redefined report publication within existing bundles, original immutability's launcher-banner exception, second-opinion mode/default/exclusion semantics, backend/model-derived filename/reply suffixes and typed result evidence; sole-member wave-end sweep completed, no stale operative restatements found.
**Suggested testing:** 3 items.
**Follow-ups created:** None (2 findings report only).

*Reviewed by review-work action*



### Remediation handback


Remediation commit: `2abd7939ac52c9813616909e388cfedfb8b78b6a` on the retained builder branch, based on `69a2a998`. Exactly two files changed: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` and `report_second_opinion_test.go`. Parent owns the re-review, merged gate, lesson/prime update, lifecycle and release. No external model calls occurred in this remediation.

Read the complete `REQ-698-review.md` before production patching; reread `skills/do-work/crew-members/debugging.md`, the whole CLI prime and the exact `reaped-by-its-own-parent` satellite entry, including its `ESRCH` trap. Parent reserved D-12 for its lesson/prime update.

**D-11 — DECIDE & STATE:** Keep a small POSIX-shell supervisor as the isolated group owner until the existing `ownedprocess.TerminateGroup` cleanup completes. Backend executable and arguments remain positional (`"$@"`), never interpolated/evaluated as shell code. The supervisor waits for the backend, writes its exact exit status to a private inherited pipe, closes that pipe and waits on parent-owned stdin. The backend closes fd3 and receives `/dev/null` stdin, so it cannot inherit either completion/control channel. Parent status parsing is bounded to32bytes and rejects malformed/noncanonical statuses. Cleanup still proves ownership through the shared helper, runs before `Cmd.Wait`, retains cleanup errors, and then reaps the supervisor/drains captured output. No unchecked numeric process-group signal or shared ownedprocess change was introduced. Windows remains fail-closed before process start.

### Structured diagnosis and evidence
- Observation: independent review probe returned success, then the redirected child wrote `late-write`300ms later.
- Hypothesis: `Cmd.Wait` reaped the isolated backend leader before `TerminateGroup` could prove its PGID; its ESRCH fallback missed live descendants. Inherited stdout/stderr instead kept `Cmd.Wait` blocked while children could still write.
- Prediction: a redirected child will write after successful capture; an inherited-output child will delay capture until it closes the pipes.
- RED confirmed both: `remediation-red.log` reports `capture returned while a descendant could still write` in both variants and inherited-pipe elapsed919.212041ms. Test command: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run 'TestSecondOpinionCapture(DrainsExitedLeaderChildren|PreservesBackendStatusAndChannels)$' -count=1`, exit1.
- Root cause confirmed; one production approach was implemented. Initial strengthened cancellation fixture asserted no write even during cleanup and failed. Measurement showed its `late` file already existed at capture return (killing its waited-for sleep resumes the shell's next statement during teardown). The corrected assertion compares the snapshot at capture return with the later state; it verifies the promised post-return boundary. Parent explicitly agreed this distinction. The successful-leader cases continue to require that the delayed write never occurs at all.
- Focused GREEN: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run 'TestSecondOpinionCapture' -count=1 -v`, pass3.137s, `remediation-green.log`. Covers redirected/inherited-pipe children, cancellation with no post-return write, exact stdout/stderr and literal argv, truthful status0/7, closed backend stdin and closed status fd3.
- Reviewer's unchanged original probe GREEN: `go -C skills/do-work/tools/do-work-cli test -overlay /private/tmp/req698-evidence/remediation-review-overlay.json ./internal/toolboxcommands -run '^TestReviewSecondOpinionLeaderExitLeavesWriter$' -count=1`, pass2.047s, `remediation-review-green.log`. Overlay only points the retained reviewer test into this worktree.
- Final race packages: `go -C skills/do-work/tools/do-work-cli test -race ./internal/toolboxcommands ./internal/resultmodel`, pass toolbox15.804s/resultmodel cached, `remediation-race.log`. This includes native adapter/metadata contracts and all existing scoped command tests after the runner change.
- `go -C skills/do-work/tools/do-work-cli vet ./...`: exit0, `remediation-vet.log`.
- `GOOS=windows GOARCH=amd64 go -C skills/do-work/tools/do-work-cli test -c ./internal/toolboxcommands -o /private/tmp/req698-evidence/remediation-toolbox-windows.test.exe`: exit0, `remediation-windows.log`.
- gofmt and `git diff --check` clean; reviewed both changed files and119insertions/14deletions. Worktree clean after commit. Live artifact bytes untouched.

### Lifetime limits and lesson handoff
The supervisor is our child and is reaped by our `Cmd.Wait`. Backend descendants are terminated through the existing bottom-up shared group boundary. Descendants orphaned by a backend that exits itself can be proved stopped, but their eventual reap belongs to init; this patch does not claim to reap another process's child. A deliberate `setsid` escape remains outside group ownership and is not signalled, as the governing lesson already states. The reviewed in-group late-writer and inherited-pipe defects are covered by the new behavioral regressions.

Proposed parent lesson wording: **Symptom:** success returned before a redirected child wrote, while inherited-output children delayed completion. **Root cause:** backend completion was conflated with ownership-lifetime completion, and reaping the sole group leader erased PGID proof. **Fix:** a supervisor retains isolated ownership, reports backend status separately and stays alive until proven bottom-up cleanup; only then wait/reap and publish. **Discovery:** reviewer minimal delayed-write probe plus differential real-child RED/GREEN tests. Parent owns updating the existing generalized prime trap/satellite/index rather than adding duplicate lesson machinery.

## Qualification

**Result:** Pass for substantive implementation at `3b593d22..31e7fcdf`. The canonical request/path-bound qualify record is satisfied; P-A-U and diff hygiene passed. Read source, wiring, prompt/action/docs and live artifact diff: registered handler → trusted launcher template → native adapters → preservation checks → exclusive page/reply publication → shared judge → original banner. Code and deterministic tests cover all modes and both generation identity checks.

**Findings judged:** Three static unreferenced-file warnings are legitimate generated reply artifacts and a Go test file discovered by the test harness. Scope comparison compares literal strings: its four undeclared artifact paths are contained by the previously declared `index-*.html` / `*-reply.md` globs, not scope expansion. The six untouched glob entries are file-class allowances; no Gemini/AGY output was published and no extra assets were needed. No outside-scope implementation path was introduced. The final live AGY retry was separately authorized and verified prompt noncompliance; actual Claude-session execution remains unproven.

**Acceptance trace:** Router, guide and prompt preserve the four forms and full editorial constraints; handler/typed-result/launcher wiring is present. Baseline RED and later GREEN prove the new command and equivalent-backlink bug. Fake backends cover default selection, excluded/missing/failing batch members, actual identity mismatch, collisions, dirty paths, old-asset changes, rename sources, symlinks, missing backlinks, all judge states and per-file banners. Live Codex and pinned Claude output exists and opens; original diff is banners only. Review must retain session and batch limitations.

*Qualified by work action*


**Remediation qualification:** Canonical qualify remains satisfied on the cumulative range; no new project file or new scope error. The owner evidence commit between merges is do-work metadata, explicitly excluded from implementation judgment. The two-file runner change preserves a live supervisor until proven group teardown; reviewed argv isolation, closed channels, exact backend status, bounded status parsing, errors and final Wait. Existing group limitations (setsid escapees and init-owned orphan reaping) are stated. Reviewer delayed-write and inherited-pipe RED cases pass after the fix.

## Testing

**Tests run:** Canonical `bash _dev/tests/maintainer-verify.sh` through `run-timed-command`: exit 0 in 119 seconds on merged source. Board: 427 tests in 37 seconds, slowest file 17.86s. CLI: 937 tests in 56 seconds, slowest file 18.49s. Both satisfy the per-file <30s budget. The gate includes contract, shipped references, shell checks, Go vet and uncached tests. No failing canonical run required a retry.

**Focused merged probe:** `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands ./internal/resultmodel -run 'TestAIReport|TestSecondOpinion|TestHandlersRegister' -count=1` passed under the 30-second probe limit. `advance` returned request/path-bound satisfied `run-blocked-check` and `green-gate` records. Only judged unused-glob scope warnings remain. Evidence: `/private/tmp/req698-evidence/remediated-test-gates.json` and `remediated-maintainer.log`. Initial gate-call input errors ran no tests; corrected to the supported probe-file form.

**Red-green validation:** Initial fake-backend tests failed because the handler was absent, then passed after implementation. The equivalent-backlink test failed three valid examples before the parser fix and passed afterwards. See Builder Evidence for exact commands/logs; final builder race checks and vet passed after that fix. Captured session proof was adapted to explicit pinned Claude because the current chat is Codex and its exact model is unavailable; no session identity was invented.

**New tests added:** Native metadata, generation model switch, mode selection, unique filenames/replies, outside dirt/rename, preserved existing assets, symlinks, semantic backlink, judge states, banner placement/idempotency, serial batch, cancellation and required backend argv.

**Existing tests updated:** Handler registration/count; typed result serialization; source/installed launcher prompt discovery. These are intentional new-command wiring contracts.

**Live verification:** Codex gpt-6-astra and pinned Claude claude-opus-5-5 both returned judge pass and banner inserted, outside paths empty. Codex's disclosed source-map amendment separately passed wide/phone light/dark render checks. Parent browser interactions passed on both pages. Batch excluded actual caller Codex; Gemini failed exit55 before judge and AGY failed the original backlink validator before judge. That validator had a proven false-rejection case and is fixed. User explicitly authorized one AGY retry after two automatic approval-review rejections; the retry completed with native SUCCESS but was refused because navigation preceded its opening backlink; no page/reply/banner was published, judge not-run, outside paths empty. Native output and candidate were retained for this run. No all-backend live success or actual Claude-session success is claimed.

**Heavy verification plan:**
- Range: `3b593d22..31e7fcdf`. All six lanes selected because report HTML paths have uncertain classification.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html

*Verified by work action*


**Post-review closure verification:** Cumulative range `3b593d22..31e7fcdf`. Final canonical gate exit 0/119s and focused probe satisfied; no error findings. Runner RED cases for exited leaders with redirected and inherited output passed after the supervisor fix; the original independent reviewer overlay passed unchanged. Builder race tests passed in 15.804s, vet and Windows compile passed. Later evidence supersedes the pre-remediation green for runtime claims. Parent-owned lesson/prime/index scope is deliberately deferred to finalization and is an expected declared-not-touched warning.

## Lessons Learned

**What worked:**
- Probe installed native protocols before writing fake adapters, and verify actual identity again on the generation result. Isolated clean checkouts made strict outside-bundle proof meaningful.
- Independent review's minimal delayed-write probe exposed a lifecycle defect missed by cancellation tests; the unchanged reviewer probe passed after remediation.

**What didn't:**
- The assumed Codex JSON model header does not exist in this installation; the native text startup envelope supplies it. Attribute-order-sensitive backlink validation rejected equivalent HTML and needed a semantic regression.
- Waiting/reaping the backend leader before cleanup erased process-group ownership proof. A live supervisor now separates backend exit from group lifetime, retains proof through cleanup, then permits Wait/publication. Redirected-child and inherited-pipe RED/GREEN tests pin the boundary.

**Worth knowing:** A configured remote and local commit do not prove source links are published; use repository paths when uncertain. Keep model claims, launcher judge results and verifier edits distinct. Group cleanup proves in-group descendants stopped, not that init has reaped orphans; deliberate setsid escapees remain outside that group guarantee.

## Orientation

[MAP CHANGED] The toolbox can now request a model-specific second opinion on an existing report through a new CLI execution path. It preserves prior bundle files, reports native model and render evidence, and links only passing pages from the original. Two verified example pages are included; Gemini availability, actual-session execution and successful AGY output remain explicit live-proof limits. The five touched-area prime references were spot-checked for stale paths; none became stale. The existing CLI process-lifetime trap is being amended for the observed recurrence.


## Review

**Overall: 86.4375%** | 2026-10-15T13:13:51Z

| Dimension | Score |
|-----------|-------|
| Requirements | 93.75% |
| Code Quality | 96% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F2: Requested live acceptance remains incomplete: prior targeted Codex and pinned Claude passed, actual-session proof is absent, Gemini failed exit 55, and the authorized AGY retry correctly refused a backlink following navigation with judge not-run and no publication. No later model run established broader acceptance. — impact-user-visible → report only

**Minor findings:** None.
**Resolved findings:** F1 closed by `2abd7939`, merged at `31e7fcdf`: supervisor retains group ownership through descendant teardown; unchanged independent delayed-writer overlay passes on merged main (exit 0, 1.082 seconds), with inherited-pipe/status/cancellation regressions and final canonical gate passing.
**Acceptance:** Partial — local and merged integration pass after remediation; prior live Codex/pinned-Claude success and batch/AGY refusal evidence retained; actual-session/full-batch live success incomplete; consumer deployment unassessed.
**Restatement sweep:** redefined report publication within existing bundles, original immutability's launcher-banner exception, second-opinion modes/default/exclusion, backend/model-derived filename/reply suffixes and typed result evidence; sole-member wave-end sweep refreshed for owned-process completion semantics, no stale operative restatements found.
**Suggested testing:** (1) With separately authorized working credentials and known actual session identity, verify session invocation and remaining Gemini/AGY publication; no further external calls are authorized by this review. (2) After deployment, exercise the installed command in a clean consumer bundle to check trusted prompt discovery, relative links and banner-only changes. Consumer deployment remains unassessed. Both are report only.
**Follow-ups created:** None (1 finding report only).

*Reviewed by review-work action; re-reviewed after the single remediation cycle against cumulative range 3b593d22..31e7fcdf*


## Heavy Verification Plan

Base: `3b593d22993d38da80024346b4af12b80626dd09`
Target: `31e7fcdfd8fc7a934a5ca3f13202e66a9f8c3200`

- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html, ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html

## Timing

Observed 2026-10-15T11:43:46Z to 2026-10-15T13:10:54Z: 1h 27m 08s total, 1h 05m 24s attributed across 5 events, 21m 44s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 59m 15s | 2 |
| verification-gate | 6m 09s | 3 |

Slowest stage: builder-work / implementation, 50m 50s, outcome success.
Slowest command: verification-gate / baseline-maintainer, 2m 06s, exit 0, bash (2 argv tokens).

## Heavy Verification Result

Target and execution revision: `31e7fcdfd8fc7a934a5ca3f13202e66a9f8c3200`. All six selected lanes passed, each executed in this drain; none reused. The first browser attempt explicitly skipped because Chrome was not on PATH. Retried only that lane with `QUEUE_KANBAN_BROWSER=/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`; 39 browser tests passed.

- `queue-kanban-javascript`: exit 0, 8 seconds, disposition `executed` (fingerprint_mismatch); not skipped.
- `queue-kanban-browser`: exit 0, 89 seconds, disposition `executed` (fingerprint_uncertain); not skipped.
- `do-work-cli-integrations`: exit 0, 77 seconds, disposition `executed` (fingerprint_mismatch); not skipped.
- `staged-skills`: exit 0, 36 seconds, disposition `executed` (fingerprint_mismatch); not skipped.
- `updater`: exit 0, 68 seconds, disposition `executed` (fingerprint_mismatch); not skipped.
- `installer`: exit 0, 31 seconds, disposition `executed` (fingerprint_mismatch); not skipped.

Evidence: `/private/tmp/req698-evidence/heavy-result.json` and `heavy-browser-result.json`; stderr logs retain command output. Heavy lanes have no per-file 30-second budget.

## Completion Lessons Promotion

D-12 completed after the canonical finalization preview confirmed the exact archive postimage. Amended the existing reaped-by-its-own-parent prime trap, appended two linked satellite lessons and refreshed its exact family/size index. Canonical repository links are used because installed packages omit do-work/archive. The optional KB inbox handoff remains pending; no inbox file was written.
