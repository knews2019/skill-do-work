---
id: REQ-698
title: 'ai-report second-opinion: a second model rebuilds a report bundle as index-<cli>-<model-slug>.html (session, targeted and batch modes)'
status: claimed
created_at: 2026-10-15T11:02:38Z
user_request: UR-156
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "skills/do-work-toolbox/docs/prime-agy-usage.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-*.html", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/*-reply.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go", "skills/do-work/tools/do-work-cli.sh", "_dev/tests/do-work-cli-launcher-behavior.sh", "skills/do-work-toolbox/actions/ai-report-reference.md", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/gemini-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/agy-*", "ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-*"]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
## Full Context
See `do-work/user-requests/UR-156/input.md` for complete verbatim input, including the clarity-editor prompt draft and the four decision rounds.
*Source: Can we add the following skill so we extend the ai-report and can have a second opinion from another model that rebuilds a new html, like index-[model-name-version].html i.e. index-argon.html (plus the pasted clarity-editor prompt, two amendment messages, and the migration instruction, all in the UR)*

## Triage

**Route:** C — multiple backend protocols, guarded bundle writes, render judging, and action/documentation integration require a coordinated plan.

<!-- D-XX counter: D-04 used. Next decision: D-05. -->

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
