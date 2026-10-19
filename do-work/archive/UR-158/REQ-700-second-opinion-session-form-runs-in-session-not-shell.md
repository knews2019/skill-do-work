---
id: REQ-700
title: 'Second opinion from the session model runs as a visible in-session prompt, not a hidden shell subprocess'
status: completed
created_at: 2026-10-19T14:43:26Z
user_request: UR-158
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-10-19T15:25:12Z
route: C
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md", "do-work/lessons-index.md"]
planning_at: 2026-10-19T15:39:22Z
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-10-19T15:34:40Z
  basis:
    - Route C
    - 6-file write set
    - 2 subsystems involved
    - 8 acceptance criteria
dispatch_at: 2026-10-19T16:34:20Z
builder_handback_at: 2026-10-19T16:49:40Z
integration_at: 2026-10-19T16:50:17Z
heavy_verified_at: 2026-10-19T17:14:27Z
heavy_verified_revision: c70b9b72deb421c93fd9da459c4020f27249f7a9
kb_status: pending
commit: c70b9b72deb421c93fd9da459c4020f27249f7a9
review_at: 2026-10-19T17:16:23Z
completed_at: 2026-10-19T17:16:59Z
release_at: 2026-10-19T17:16:59Z
---
# Second Opinion From The Session Model Runs As A Visible In-Session Prompt, Not A Hidden Shell Subprocess
## What
When `ai-report second-opinion` runs in its Session form (the calling agent is itself the backend, for example a Claude session passing `--backend claude --model claude-fable-5-1`), the clarity-editor pass must run where the user can watch it: as a normal prompt inside the agent session, or as a background subagent or workflow the agent starts. It must not run as a hidden shell subprocess of the agent's own CLI (`claude -p <prompt>` under the supervisor) that sits for minutes showing "No output available". The CLI keeps the mechanical parts it owns today: bundle resolution, pre-checks, staging filename and asset prefix reservation, preservation and outside-bundle Git checks, backlink check, render judge, banner insertion, reply file, and the typed `ai_report_second_opinion.runs` summary.
## Why
The user watched a 10-minute shell with no output and no way to see what the editor was reading or writing. Progress inside the session (tool calls, file reads, the page being written) is the feedback they asked for.
## Context
- Today `skills/do-work-toolbox/actions/ai-report.md` § Second-opinion form tells the agent to invoke `do-work-cli.sh --format json ai-report-second-opinion --backend <own-cli> --model <own-slug>`. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` then execs the backend executable with `-p <prompt> --output-format json` through `/bin/sh -c` supervisor, so a Claude session spawns a second headless Claude. The whole editorial pass is invisible until the process exits.
- The screenshot (see Assets) shows the agent's three status lines, then "Churned for 33s · done 5:26 PM · 1 shell still running", and the Shell details pane: Status running, Runtime 10m 11s, the full `ai-report-second-opinion --backend claude --model claude-fable-5-1 ai-reports/2026-10-19_1615_proposal-zoom-discovery-rev1` command, Output "No output available".
- The shipped prompt `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` already has placeholders (`{{BUNDLE_PATH}}`, `{{REPOSITORY_ROOT}}`, `{{OUTPUT_FILE}}`, `{{ASSET_PREFIX}}`, `{{BACKLINK_REQUIREMENT}}`) that the CLI fills before launching; the same filled prompt is what the agent would run itself.
- The follow-up screenshot (UR-158 asset `REQ-701-screenshot-2-…`) shows a second effect of the hidden-process design: the launcher reads the editor's model from the child process's reply envelope, so a subagent the editor spawned made the model look ambiguous and the finished page was deleted. REQ-701 (keep the finished page when the model check fails) owns that defect; in the in-session path there is no child envelope to probe, so the model is the session's declared one.
## Requirements
- The Session form, and a targeted `--backend <own-cli>` from inside that same agent, never launches a long-running backend subprocess of the agent's own CLI. The editorial work runs in the agent's own turns, or in a background subagent or workflow it dispatches, so progress is visible in the session.
- Everything the launcher guarantees today still holds for the in-session path: the original `index.html` and existing assets are never modified by the editor, new files use the reserved staging name and asset prefix, the published page is `index-<cli>-<model>[-N].html` paired with `<cli>-<model>[-N]-reply.md`, the banner is inserted only after preservation, outside-bundle Git status, and render judge `pass`, and the final reply lists page/reply paths, backend/model, judge verdict, and the 3 to 5 clarity changes, corrections and open questions.
- Foreign backends (codex, gemini, agy) have no in-session form, so their launcher path is unchanged by this REQ. Their lack of progress feedback is a separate request if wanted.
- Update the action prose (`actions/ai-report.md` § Second-opinion form), the user guide (`docs/ai-report-guide.md` § A Second Opinion), the CLI, its tests, and the changelog as a release per `_dev/primes/prime-releases.md`.
## Dependencies
None blocking. REQ-701 (keep the finished page when the model check fails) touches the same `secondOpinionRun` / `secondOpinionMetadata` code; build them in one pass or land REQ-701 first so the shared cleanup path is already fixed.
## Builder Guidance
Certainty: high on the outcome (no hidden shell for the session's own model, progress visible), open on the mechanism. The natural split is two CLI phases the agent calls around its own work: a prepare phase that resolves the bundle, runs the pre-checks, reserves the staging filename and asset prefix, and prints the filled prompt; and a finalize phase that runs preservation, outside-bundle Git status, backlink check, render judge, banner insertion, reply file, and the typed summary against the file the agent wrote. Model identity for the in-session path is the session's declared model, since there is no backend process to query; record it as such rather than claiming CLI verification. The builder decides the exact flags and names. Reuse `secondOpinionBundle`, `secondOpinionSnapshot`, `secondOpinionPreserved`, `secondOpinionOutside`, `secondOpinionNames`, `secondOpinionBacklink`, `secondOpinionBanner`, and `secondOpinionRender` rather than duplicating them. Delete the Claude-specific exec branch only if nothing else needs it after the split.
## Red-Green Proof
**RED prompt/case:** In a Claude Code session with an existing bundle, run `/do-work-toolbox ai-report second-opinion ai-reports/<bundle>/index.html`. A Go test in `report_second_opinion_test.go` with a fake `claude` executable on PATH that records every invocation asserts that the session-form path (own backend, own model) records zero editorial-pass invocations and instead hands back the filled prompt and reserved output path for the agent to act on.
**Why RED now:** The action prescribes a single CLI call, and `report_second_opinion.go` execs `claude -p <prompt>` for the session's own model; the fake executable records one long-running invocation and the session shows one silent shell for over 10 minutes.
**GREEN when:** The same invocation produces visible work in the session (the agent, or its background subagent or workflow, reads the bundle and repository and writes the sibling page in observable steps), the CLI then finalizes it with the same preservation, Git, judge and banner guarantees, the final reply lists page, reply, model, judge verdict and the editor's changes, and the fake executable records no editorial-pass invocation.
**Validation:** Inferred during capture from the screenshot and the user's words.
## Assets
- `do-work/user-requests/UR-158/assets/REQ-700-screenshot-1-second-opinion-shell-run.png` — Claude Code transcript: the slash command with the bundle path, three agent status lines ("I'll read the routed ai-report action first…", "Git is clean and the bundle exists. Running the session form: the CLI launches a clarity-editor pass under my own backend and model, then judges the result.", "The generation is running in the background. It reads the proposal report, checks its claims against the repo, writes a sibling page, and then runs the render judge. I'll report once it finishes."), then "Churned for 33s · done 5:26 PM · 1 shell still running". Below, the Shell details pane: Status running, Runtime 10m 11s, Command `.claude/skills/do-work/tools/do-work-cli.sh --repo-root . --format json ai-report-second-opinion --backend claude --model claude-fable-5-1 ai-reports/2026-10-19_1615_proposal-zoom-discovery-rev1 > /private/tmp/…`, Output "No output available".
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22085 tokens, over the 2000 budget; `slugged: partial` so no targeted entry is legal. Matched because the REQ changes native backend identity and background-worker behavior in do-work-cli internals (families native-backend-identity, background-worker-self-bound, interruptible-blocking-io).
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, over budget; `slugged: partial`. Matched because the REQ changes a prescribed command block inside a shipped action.
- `_dev/primes/lessons-action-files.md` — 8128 tokens, over budget; `slugged: partial`. Matched because the REQ changes an action's routing and an alternate artifact writer path.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ in full (What through Scope), the brief with its Coordinator note, UR-158's input, the archived REQ-701, the crew rules (general, coding-guardrails § 5, shared-principles, communication-style, testing, anti-slop), the three primes (action-files § Cross-Referencing and the alternate-writer-contract-drift trap; shell-commands § Closed Enumerations Go Stale; releases), and the `native-backend-identity` and `reaped-by-its-own-parent` bullets of `lessons-do-work-cli.md`. Followed the Plan's five tasks in order with the decided mechanism; no change to the mechanism. *(from the builder hand-back)*
- [x] **[APPLY]:** Task 1 RED prepare test; task 2 result model, parser, handler branch, `secondOpinionOutsideCheck`, `secondOpinionPrepare` (GREEN); task 3 refactor under green (digest-aware `secondOpinionPreserved`, `secondOpinionPublication` extraction, pre-flight probe rerun); task 4 RED finalize test, then `--finalize` parsing, handler bundle derivation with the supersession guard, `secondOpinionFinalize` (GREEN); task 5 prose. Scope held to the five write-boundary files. *(from the builder hand-back)*
- [x] **[UNIFY]:** `git diff 5f237933 --stat` shows exactly the five files (415 insertions, 65 deletions). Checks (worktree or module root): `git diff --check` clean before each commit; `gofmt -l ./internal` prints nothing; `go vet ./internal/...` exit 0; `go test -count=1 ./internal/toolboxcommands ./internal/resultmodel` both `ok` (17 s, 0.6 s); REQ-700 probe exit 0; pre-flight probe exit 0 (three runs); `action-shell-blocks.sh`, `prescribed-shell-canonicalization.sh`, `shipped-package-reference-contract.sh`, `audit-lockins.sh` all exit 0. Files checked by reading the full diff: the Go file (every hunk, the exec path's argv, metadata, capture, supervisor and probe calls unchanged), the test file, the result model, both prose files. No debug artifacts. *(from the builder hand-back)*
*Source: [Image #1] <- let's update the second opinion so we have a feedback, this should have just triggered a normal prompt inside the ai-coder that way I can see the progress, alternatively it's fine to start a background subagent, or a workflow, but not a shell like this*

---

## Triage

**Route: C** - Complex

**Reasoning:** The outcome is fixed (no hidden shell for the session's own model) but the Builder Guidance leaves the mechanism open, and the change crosses the CLI argument contract, the typed result model, the shipped prompt contract, the action prose, the user guide and the Go tests. The pre-dispatch agent decides the mechanism before the builder starts so the brief carries one decided design.

**Planning:** Required

<!-- D-XX counter: none used. Next decision: D-01. -->
*No `## Open Questions` section exists; Step 3.5 skipped.*

## Plan

**Decided mechanism (pre-dispatch, do not reopen):** two phases on the existing `ai-report-second-opinion` subcommand. `--prepare --backend <cli> --model <slug> [<bundle-dir>|latest]` resolves the bundle, reads the template, refuses outside-bundle Git dirt, snapshots the bundle, creates the dotfile staging page exactly as `secondOpinionRun` does today, writes a sidecar `<stem>.json` next to it (backend, model, prefix, and the before-snapshot as mode + sha256 digests, never bytes), fills the prompt with `{{OUTPUT_FILE}}` set to the staging page's absolute path, and prints one `runs[]` row with `state: prepared`, `staging_page`, `staging_reply` (`<stem>-reply.md`), `asset_prefix`, `prompt`, and `model_evidence: session-declared`. It never calls `exec.LookPath`, `secondOpinionCapture`, `secondOpinionArgv` or `secondOpinionMetadata`. `--finalize <staging-page-path>` (the exact `staging_page` value; no positional bundle, no backend flags) derives the bundle from the path, loads the sidecar, and runs: outside-bundle Git status; never-written stage (missing or 0 bytes) → failed and the stage, sidecar and reply stage are removed (the only deletion case); preservation against the sidecar digests with the stage, reply stage and sidecar as the allowed new names plus prefixed regular files; blank reply stage → failed, files kept; backlink; `secondOpinionNames`; publish page and reply (reply rewritten from the stage path forms to the published page name); remove the three staging leaves; render judge; banner only on `pass` with the identity and byte checks against a snapshot taken at finalize start. Result `mode` is `session` for both phases. The prompt template is unchanged. The `claude` exec branch stays for other agents that target `--backend claude`.

### Files to modify
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify): parser, handler branch, `secondOpinionPrepare`, `secondOpinionFinalize`, extracted `secondOpinionPublication`, generalized `secondOpinionPreserved`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify): two new tests plus a `secondOpinionFakeLogger` helper.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify): five `omitempty` string fields on `AIReportSecondOpinionRun`; the nil-slice normalizer is unchanged.
- `skills/do-work-toolbox/actions/ai-report.md` (modify): § Second-opinion form Session and Targeted bullets plus the model-evidence sentence.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify): one paragraph in § A Second Opinion.
- CHANGELOG, VERSION, mirror: integrator only; the builder proposes the entry text.

### Order of changes
1. **RED: prepare test.** Add `secondOpinionFakeLogger(t, "claude") string` (writes a `#!/bin/sh` script into `os.Getenv("PATH")` that appends `"$@"` to a log file and returns the log path) and `TestSecondOpinionSessionPrepareNeverExecsOwnBackend`. At base it fails at `secondOpinionRows` ("missing typed backend summary") because `secondOpinionParse` returns `unexpected argument "--prepare"`.
2. **GREEN: result model, parser, prepare.** `result_model.go`: add `StagingPage`, `StagingReply`, `AssetPrefix`, `PromptText`, `ModelEvidence` (JSON `staging_page`, `staging_reply`, `asset_prefix`, `prompt`, `model_evidence`, all omitempty) after `OutsidePaths`. `secondOpinionOptions` gains `prepare bool`, `finalize string`; `secondOpinionParse` adds `case argument == "--prepare"` and `--finalize`/`--finalize=` through `optionValue`; validation: `prepare` requires `backend` and `model` and forbids `all`/`exclude`; `finalize` forbids a positional bundle, `backend`, `model`, `all`, `exclude`, `prepare`. `secondOpinionUsage` becomes `ai-report-second-opinion [--backend codex|gemini|agy|claude [--model <slug>] | --all --exclude <own-cli> | --prepare --backend <cli> --model <slug>] [<bundle-dir>|latest] | --finalize <staging-page>`. `handleAIReportSecondOpinion`: when `prepare`, set `payload.InvocationMode = "session"`, build one row (`RequestedModel = ActualModel = model`, `ModelEvidence = "session-declared"`), call `secondOpinionPrepare(repository, bundle, template, &row)`; add `prepared` to the state exemption in the finding loop (`row.RunState != "ran" && != "excluded" && != "prepared"`) or prepare reports a failure. New `secondOpinionPrepare`: outside check (lift the `checkOutside` closure to `secondOpinionOutsideCheck(repository, bundle, row) error` shared by prepare, finalize and run), `secondOpinionSnapshot`, prefix, `os.CreateTemp(bundle, "."+prefix+"staging-*.html")`, sidecar `<stem>.json` from a `secondOpinionRecord{Backend, Model, Prefix string; Files map[string]secondOpinionDigest{Mode uint32; SHA256 string}}` built by `secondOpinionDigests(before)`, the same `strings.NewReplacer` with the absolute stage path for `{{OUTPUT_FILE}}`, the `{{` check, then the row fields and `RunState = "prepared"`. Set `row.ModelEvidence = "native"` in `secondOpinionRun` where `row.ActualModel` is assigned after generation. Test 1 goes GREEN.
3. **Refactor under green (no new test).** `secondOpinionFile` gains `mode fs.FileMode` and `digest string`; `secondOpinionSnapshot` fills them (sha256 of data for regular files, empty for directories). `secondOpinionPreserved(bundle, before, allowed []string, prefix)`: when `old.info != nil` keep today's mode, bytes and `os.SameFile` check; otherwise compare `mode` and `digest`; new paths pass if in `allowed` or a prefixed regular root file. Callers in `secondOpinionRun`: `nil, ""` for the probes, `[]string{stageName}, prefix` for generation. Extract the tail of `secondOpinionRun` from "read page" to banner publication into `secondOpinionPublication(repository, bundle, stageName string, page []byte, reply string, original secondOpinionFile, stageLeaves []string, <REQ-701 mismatch input>, row, result)`. It carries REQ-701's landed behavior (unverified-name publication, no deletion of a finished page); after both publishes succeed it removes `stageLeaves` through `os.OpenRoot(bundle)`, so REQ-701's post-publication stage removal moves into it and only the never-written cleanup stays in `secondOpinionRun`. Run the package tests; every test named below stays green.
4. **RED then GREEN: finalize.** Write `TestSecondOpinionSessionFinalizePublishesAgentWrittenPage` with four subtests (publishes; judge pass inserts banner; original modified keeps stage; never written cleans up). RED: `--finalize` unknown. Then the handler branch: `absolute := judgeAbsolutePath(repository, options.finalize)`; `bundle, err := secondOpinionBundle(repository, filepath.Dir(absolute))` and require the returned bundle to equal `filepath.Dir(absolute)` (the supersession chain-follow would otherwise jump to a newer revision); no template read. `secondOpinionFinalize(repository, bundle, stageName, row, result)`: validate `stageName` against `secondOpinionStagePattern` (`^\.[A-Za-z0-9][A-Za-z0-9._-]*-staging-[0-9]+\.html$`), load the sidecar (missing → "no prepare record for this staging page"), fill the row from it, then the decided order: outside check; stage missing or 0 bytes → "staging page was not written" plus removal of stage, sidecar and reply stage; `secondOpinionPreserved(bundle, recorded, []string{stage, sidecar, replyStage}, prefix)` where `recorded` is the sidecar turned into `secondOpinionFile{mode, digest}` entries; reply stage missing or blank → "no final editorial reply"; `after := secondOpinionSnapshot(bundle)`; `secondOpinionPublication(..., after["index.html"], []string{stage, replyStage, sidecar}, none, ...)`.
5. **Prose and changelog proposal.** `ai-report.md` § Second-opinion form: rewrite the Session bullet (and the Targeted bullet when the target is the agent's own CLI) to the prepare → do the prompt (own turns, or one background subagent or workflow told to write the FINAL REPLY to `staging_reply`) → finalize sequence; say `model_evidence: session-declared` is never CLI-verified; qualify "Native backend metadata resolves the model" to launcher-run backends. `ai-report-guide.md`: one plain-words paragraph. Propose a changelog entry with a descriptive title (the integrator bumps the version).

### Architectural decisions
- Same subcommand, one boolean flag and one valued flag; no new command, file, or result-level type.
- The sidecar holds digests, never bytes; `secondOpinionPreserved` compares digests only when `info == nil`, keeping the exec path's `os.SameFile` proof.
- One shared tail `secondOpinionPublication` serves exec and finalize; REQ-701's unverified branch lives there and finalize passes "none".
- Session rows record `model_evidence: session-declared`; exec rows record `native`.
- Deletion only when the stage was never written; every other failure keeps files on disk (REQ-701's rule).
- `secondOpinionArgv` and `secondOpinionMetadata` keep their `claude` branches.

### Testing approach
- Focused: `go test -count=1 ./internal/toolboxcommands -run 'TestSecondOpinionSessionPrepareNeverExecsOwnBackend|TestSecondOpinionSessionFinalizePublishesAgentWrittenPage'` from the module root `skills/do-work/tools/do-work-cli`.
- Package: `go test -count=1 ./internal/toolboxcommands ./internal/resultmodel`; `gofmt -l ./internal`; `go vet ./internal/...`.
- Must stay green: `TestSecondOpinionNativeModelCollisionAndSkippedJudge`, `TestSecondOpinionTargetNeverSubstitutes`, `TestSecondOpinionOutsideDirtAndOriginalMutationRefuse`, `TestSecondOpinionBatchIncludesMissingAndExcluded`, `TestSecondOpinionJudgeVerdictsGateOriginal`, `TestSecondOpinionAllNativeAdaptersAndSerialFailureContinuation`, `TestSecondOpinionPreservesAssetsAndRejectsUnprefixedSymlinkAndMissingBacklink`, `TestSecondOpinionFreshGenerationModelMismatchRefusesPublication` (as REQ-701 reshaped it), `TestSecondOpinionBackendCommandSafetyContract`, REQ-701's two new tests, `TestHandlersRegisterCanonicalToolboxCommands` (exactly 10 handlers, so flags, not a new command), and `resultmodel.TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`.

### Risks
- The nil-`info` branch of `secondOpinionPreserved`: a wrong default could skip the `os.SameFile` check for exec runs or dereference nil; the `original` and `asset.css` subtests catch a weakened check.
- Extracted tail ordering: `RunState = "ran"` is set after the judge and reset to `failed` on later errors; `JudgeVerdictsGateOriginal` (all four verdicts) and `NativeModelCollisionAndSkippedJudge` (reply rewrite, `-2` suffix) catch drift.
- Stage removal moving into the helper: `terminal selects once` counts bundle files and `FreshGenerationModelMismatch` checks on-disk state.
- `prepared` must be exempt from the finding loop, or prepare returns `failure`; test 1 asserts success.
- `secondOpinionBundle` follows supersession; finalize guards that the resolved bundle is the stage's own directory.

### Consumer field contract (per `runs[]` row)
- Identity: `backend`, `requested_model`, `actual_model` (equal for session rows).
- Provenance: `model_evidence` (`session-declared` | `native`), `asset_prefix`.
- State: `state` (`prepared` | `ran` | `failed` | `skipped` | `excluded`), `reason`, `outside_paths`.
- Prepare outputs: `staging_page` (absolute), `staging_reply` (absolute), `prompt` (filled text); `page` and `reply` empty; `judge_verdict: not-run`, `banner_state: not-inserted`.
- Finalize outputs: `page`, `reply` (repo-relative), `judge_verdict`, `judge_findings`, `banner_state` (`inserted` | `already-present` | `not-inserted`); staging fields empty. Result `mode` is `session` for both.

*Generated by Plan agent (constrained to the pre-dispatch agent's decided mechanism)*

### Plan validation (pre-dispatch)
- Requirement coverage: no hidden subprocess and visible work → tasks 1, 2, 5; every launcher guarantee kept in-session → tasks 3, 4; foreign backends unchanged → task 3 and the decisions; prose, guide, CLI, tests, changelog → task 5 plus the integrator's release. No uncovered requirement.
- Orphan tasks: none. `secondOpinionOutsideCheck` and `secondOpinionFakeLogger` are support for named tasks.
- Scope sanity: **Plan has 5 tasks — quality degrades past 3.** Not split: tasks 1 to 4 are one serialized Go slice in a file REQ-701 also edits, and task 5 is prose that must ship with it.
- Correction: the Plan agent's "then `_dev/tests/maintainer-verify.sh`" is struck. The builder never runs the repository gate; the integrator does.
- Consumer field contract: present above; the action prose reads exactly those fields.

## Exploration

Read directly by the pre-dispatch agent (every file below in full unless a range is given) plus one Explore agent sweep for restatements of the contract.

**Where the hidden shell comes from.** `handleAIReportSecondOpinion` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` looks the backend up with `exec.LookPath` and calls `secondOpinionRun`, which runs two subprocesses through `secondOpinionCapture` (a `/bin/sh -c` supervisor, 16-minute bound): a read-only identity probe (`secondOpinionArgv(..., true)`) and the generation (`secondOpinionArgv(..., false)`, for `claude`: `-p <prompt> --output-format json --model <slug> --permission-mode acceptEdits`). Nothing is printed until the process exits, which is the 10-minute silent shell in the screenshot. The agent's cwd is the repository; the subprocess cwd is the bundle, which is why `{{OUTPUT_FILE}}` is a bare basename today.

**What the launcher guarantees and where.** In `secondOpinionRun`: `secondOpinionOutside` (git status outside the bundle, before and after), `secondOpinionSnapshot` + `secondOpinionPreserved` (existing paths unchanged by mode, bytes and `os.SameFile`; new paths only the stage or prefixed regular root files), `secondOpinionBacklink`, `secondOpinionNames` (collision suffix), `rootedPublishFile` for page and reply (the reply's stage-path forms are rewritten to the published page name), `secondOpinionJudge` (a package var over `secondOpinionRender`, stubbed by tests), then the banner only after a pass with identity and byte checks on `index.html` and `secondOpinionBanner`. The stage leaf is `os.CreateTemp(bundle, "."+prefix+"staging-*.html")` (dotfile, digit suffix). The in-session design keeps every one of these by splitting the function at the generation call: everything before it is the prepare phase, everything after it is the finalize phase, and the before-snapshot crosses the process boundary as digests in a sidecar file.

**Model identity.** `secondOpinionMetadata` reads the model from the backend envelope (for `claude`, `modelUsage`). REQ-701 (keep the finished page when the model check fails; requested model wins when present) lands on main before this builder starts. In-session there is no envelope, so the row records the declared model with a new `model_evidence` field (`session-declared` versus `native`), never a verified claim.

**Result model.** `AIReportSecondOpinionRun` in `internal/resultmodel/result_model.go` (fields `backend`, `requested_model`, `actual_model`, `state`, `page`, `reply`, `judge_verdict`, `judge_findings`, `banner_state`, `reason`, `outside_paths`); the normalizer near line 904 only fills nil slices, and `TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence` round-trips without a golden string, so `omitempty` string fields are safe. `ExactTextOutput` is `json:"-"`, so under `--format json` the filled prompt must be a JSON field, not text output. `TestHandlersRegisterCanonicalToolboxCommands` pins exactly 10 handlers: flags, not a new command. No Go code outside `report_second_opinion.go` calls the `secondOpinion*` helpers. No `--prepare` or `--finalize` flag exists anywhere in the module; boolean flags are matched as literal strings in `secondOpinionParse` (`--all`) and valued flags go through `optionValue`.

**Prose sites.** `skills/do-work-toolbox/actions/ai-report.md` § Second-opinion form (the Session, Targeted and Batch bullets, the invoke paragraph, the guarantees paragraph, the finish paragraph) and `skills/do-work-toolbox/docs/ai-report-guide.md` § A Second Opinion (three paragraphs). The template `docs/ai-report-second-opinion-prompt.md` needs no change: its placeholders are filled the same way and the reply destination is stated by the action, not the template. Restatements elsewhere (Explore sweep): `actions/ai-report-reference.md:148` ("only the launcher may add a pass-gated banner") stays true because finalize is the launcher's phase and the agent never inserts the banner, and it links `ai-report.md#second-opinion-form`, so that heading must keep its text; `SKILL.md:3,24` are trigger words only; `help.md` has no second-opinion line; the 0.306.0 changelog entry is history and is not edited; `lessons-do-work-cli.md:5,7` stay true for the exec backends.

**Contract scripts that read the touched files** (each runs alone with `bash <path>`): `_dev/tests/action-shell-blocks.sh` (lints every shell fence in `skills/**/*.md` with the quiet-grep scanner, so a new `--prepare`/`--finalize` fence is checked), `_dev/tests/prescribed-shell-canonicalization.sh` (`ai-report.md` must keep its pointer to `../../do-work/docs/prescribed-shell-primitives.md`), `_dev/tests/shipped-package-reference-contract.sh` (links and heading anchors, including `#second-opinion-form`), `_dev/tests/audit-lockins.sh` (bans manual-fallback sentences in actions). No size, line-budget or description-trigger test covers `ai-report.md` or the guide.

**Required-lessons consult (re-run at claim):** the three matching satellites are still `slugged: partial` and over the 2000-token budget, so the capture-time `## Required Lessons — Dropped for Budget` section stands unchanged. The pre-dispatch agent read the two `native-backend-identity` and `reaped-by-its-own-parent` bullets in `lessons-do-work-cli.md` directly: they explain why the exec path verifies identity from the envelope and why `secondOpinionCapture` keeps its supervisor; both stay as they are.

*Generated by the pre-dispatch agent's direct reading plus an Explore agent sweep*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — prepare and finalize phases, parser, extracted publication tail, digest-aware preservation
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — the two session-form tests and the fake-logger helper
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — five omitempty string fields on the run row
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — § Second-opinion form: the in-session sequence
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — § A Second Opinion: one paragraph
- `CHANGELOG.md` (modify, integrator's release) — the entry
- `skills/do-work/CHANGELOG.md` (modify, integrator's release) — the byte-identical mirror
- `VERSION` (modify, integrator's release) — the version bump
- `skills/do-work/VERSION` (modify, integrator's release) — the version mirror
- `skills/do-work/actions/version.md` (modify, integrator's release) — the Current version line
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify, integrator, only if the hand-back proposes a lesson) — one lesson bullet
- `do-work/lessons-index.md` (modify, integrator, only with the lesson) — token refresh for that satellite

**Files I will NOT touch:** `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (placeholders unchanged), `skills/do-work-toolbox/actions/ai-report-reference.md` (its banner sentence stays true), `skills/do-work-toolbox/SKILL.md` and `actions/help.md` (trigger words unchanged), `skills/do-work/tools/do-work-cli.sh` (the template export stays), `commands.go` and `commands_test.go` (no new command), the `secondOpinionCapture` supervisor and the codex, gemini and agy adapters.

**Acceptance criteria (restated from REQ):**
- [ ] The Session form, and a targeted run whose backend is the agent's own CLI, never launch a subprocess of that CLI: a fake executable on PATH records zero invocations across prepare and finalize.
- [ ] Prepare prints the filled prompt, the reserved staging page path, the reserved reply path and the asset prefix, so the editorial work runs in the agent's own turns or in a background subagent or workflow it dispatches.
- [ ] Finalize keeps every launcher guarantee: original `index.html` and existing assets unchanged (digest check), new files only under the reserved staging names and asset prefix, outside-bundle Git status clean, backlink present, page published as `index-<cli>-<model>[-N].html` with `<cli>-<model>[-N]-reply.md`, render judge run, banner inserted only after preservation, Git status and judge `pass`.
- [ ] A finished staging page is never deleted by finalize; only a never-written (empty) stage is cleaned up together with its sidecar and reply stage.
- [ ] The in-session model is recorded as session-declared (`model_evidence: session-declared`), never as CLI-verified; exec rows say `native`.
- [ ] The final reply in the action lists page and reply paths, backend and model, judge verdict, and the 3 to 5 clarity changes, corrections and open questions.
- [ ] Foreign backends (codex, gemini, agy) keep the single-call launcher path; their existing tests stay green.
- [ ] `actions/ai-report.md` § Second-opinion form and `docs/ai-report-guide.md` § A Second Opinion describe the in-session sequence; the CLI, its tests and the changelog ship as one release per the releases prime.

## Pre-Flight

**Git:** ✓ clean at 5f237933 (main after REQ-701's release 0.306.3), the builder's base commit; nothing outside `do-work/` was dirty when the coordinator recorded pre-flight
**Tests baseline:** ✓ `bash do-work/runs/work-2026-10-19-152525/REQ-700-preflight-probe.sh` exit 0, launched (13 named `toolboxcommands` tests including the four REQ-701 tests and the handler-count pin, plus `TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`; `do-work/working/baseline.json`). Repository gate `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at 5f237933: run 1 exit 1 (board per-file budget at 1-minute load 10 plus a `TestBlockedProbeRunsFromSelectedRepositoryRoot` launch flake), run 2 exit 0 in 130 s under the one-retry rule; green-gate record `reported_direct_zero_exit` at revision 5f237933
**Dependencies:** ✓ Go toolchain and ShellCheck present (the gate and contract scripts ran)

*Checked by the coordinator before dispatch (2026-10-19T16:28:55Z to 16:33:34Z); recorded here by the integrator from the run's progress log, `baseline.json` and the green-gate record*

## Decisions

*(from the builder hand-back)*

- **D-01 (DECIDE & STATE): the sidecar is a dotfile inside the bundle**, `<stem>.json` beside the stage, mode 0600, holding backend, model, prefix and the before-snapshot as mode plus sha256 digests. An agent with repository write access could edit it, but the in-session agent is already the trusted editor, the exec path never reads a sidecar, and the finalize preservation still refuses any change to a pre-existing path the record names.
- **D-02 (DECIDE & STATE): the reply stage is agent-written, not pre-created.** Prepare returns `staging_reply` as `<stem>-reply.md` without creating it; finalize fails with "no final editorial reply" (files kept) when it is missing or blank. Pre-creating it would make an empty reply indistinguishable from an unwritten one.
- **D-03 (DECIDE & STATE): the extracted helper's REQ-701 input is `unverifiedModel string`**, the requested model whose identity native evidence did not confirm; empty means verified. The exec path passes `resolved` when unverified; finalize passes `""` ("none"). The seen list comes from `row.SeenModels`, so no second parameter.
- **D-04 (DECIDE & STATE): a prepare that fails after reserving its files removes them.** The deferred cleanup in `secondOpinionPrepare` removes the stage and sidecar when the row did not reach `prepared` (for example an unresolved placeholder). No agent has written to them at that point, so this is the never-written case, not a finished page.
- **D-05 (DECIDE & STATE): in the exec path, a reply publish failure after the page publish now keeps the staging page** and names it in `reason`, where before REQ-700 the deferred cleanup removed it (the old `stagePublished` flag was set after the page publish). The stage leaves are removed only after both publishes succeed, so the never-delete-a-finished-page rule now covers this path too. No test covered the old behavior; the kept bytes are harmless and visible.
- **D-06 (DECIDE & STATE): finalize derives the bundle from the staging page's parent and refuses a supersession jump** (`staging page must sit in a current report bundle`) so a page prepared in an older revision is never published into a newer bundle.

No ESCALATE decisions.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified) — parsing and validation of the prepare and finalize flags; the usage string names both phases; the staging-name pattern; the sidecar record and digest types with their builder; the snapshot entry gains mode and digest; the preservation check takes an allowed-names list and compares digests when no live file info exists; the outside-bundle check lifted out of the closure; the prepare phase; the finalize phase; the publication tail extracted from the exec run (carries REQ-701's unverified publication and removes the stage leaves after both publishes); exec rows record native model evidence; one package const deduplicates the backlink-requirement literal. The argv builder, the metadata reader, the capture supervisor and the adapters are untouched.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified) — a fake-logger backend helper, the prepare test (own backend never executed) and the finalize test with four subtests (publishes; judge pass inserts banner; original modified keeps stage; never written cleans up).
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — five optional string fields on the second-opinion run row (staging page, staging reply, asset prefix, prompt, model evidence); normalizer unchanged.
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — § Second-opinion form: the Session bullet rewritten to the prepare, do-the-prompt, finalize sequence with the five row fields; the Targeted bullet routes a target equal to the agent's own CLI through the same sequence; the native-metadata sentence qualified to launcher-run backends. Heading text unchanged; no new shell fence; no new citation.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — § A Second Opinion: three plain-words sentences on the in-session pass, what the CLI prepares and finalizes, and the declared model label; the native-metadata sentence qualified to launcher-run backends. REQ-701's unverified-name sentences in both files are kept.

**What was done:** The session form of the second-opinion command no longer execs the agent's own CLI as a hidden subprocess: the prepare phase reserves the staging page, writes a digest sidecar and returns the filled prompt; the agent does the editorial pass in its own turns or in a background subagent or workflow; the finalize phase re-checks outside-bundle Git status, preservation, backlink and reply, then publishes, judges and inserts the pass-gated banner through the publication tail now shared with the exec path. Session rows record the model as session-declared, exec rows as native; codex, gemini and agy keep the single-call launcher path.

## Discovered Tasks

*(from the builder hand-back)*

- Foreign-backend exec path (codex, gemini, agy) still shows no progress while the subprocess runs; the REQ says this is a separate request. impact-minor → report only
- `skills/do-work-toolbox/actions/help.md` has no second-opinion line, so `help` never surfaces the form or its two phases. impact-minor → report only
- `actual_model` still shows the probe model on a non-identity generation failure in the exec path (left by REQ-701's builder). impact-minor → report only
- Kept staging files (finished page after a refusal, in both the exec and session paths) are never auto-cleaned; a later prepare or finalize does not sweep them. impact-minor → report only
- The handler's finding text for a finalize that finds no sidecar reads `: no prepare record for this staging page` with an empty backend name, because the backend is only known from the sidecar. impact-negligible → report only

## Qualification

**Gate records (`advance --diff-range e11d3fca..4362eb61`):** `qualify` satisfied (no debug artifacts, no output primitives, all three P-A-U boxes ticked with the hand-back's text); `scope-drift` reports six `SCOPE-DECLARED-NOT-TOUCHED` warnings, all release files the Scope lists for the integrator on purpose (`CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`); they are touched by the finalization transaction, not by the builder. No undeclared touch: `git diff e11d3fca..4362eb61 --stat` lists exactly the five builder files (415 insertions, 65 deletions); nothing under `do-work/` is in the range (queue guard printed nothing before the merge). A first qualify pass refused with `QUALIFY-CLAIMED-PATH-MISSING` for every backticked identifier in the Implementation Summary; the section was rewritten with only file paths in backticks, no change to the claims.

**Requirement trace against the real files at 4362eb61:**
- *Session form and own-CLI target never launch a subprocess of the agent's CLI; work is visible in the session.* The handler's switch runs `secondOpinionFinalize` or `secondOpinionPrepare` before any `exec.LookPath`; prepare calls no capture, argv or metadata helper. `TestSecondOpinionSessionPrepareNeverExecsOwnBackend` and the `publishes` subtest assert the fake-logger backend on PATH was never invoked. The action's Session bullet now prescribes prepare, do the prompt in the agent's turns or one background subagent or workflow, finalize; the Targeted bullet sends a target equal to the agent's own CLI through that sequence.
- *Every launcher guarantee holds in-session.* Finalize: outside-bundle Git status (`secondOpinionOutsideCheck`, again inside the publication tail before the banner), preservation against the sidecar digests with the three staging leaves as the only allowed new names plus prefixed regular root files, backlink (inside `secondOpinionPublication` for the verified path), reserved staging name validated by the stage pattern, published names from `secondOpinionNames` (`index-<cli>-<model>[-N].html` with `<cli>-<model>[-N]-reply.md`), render judge, banner only after judge `pass`, outside check, inode and mode identity of `index.html` against the snapshot taken at finalize start. Subtests `publishes` (names, reply rewrite, leaves removed), `judge pass inserts banner` (original equals the banner-only expectation), `original modified keeps stage` (failed, page kept, no banner).
- *A finished page is never deleted by finalize.* The only removal is the stage-missing-or-empty branch (`never written cleans up`); every other refusal returns with the files on disk and the deferred append names the kept page.
- *Model recorded as session-declared, never CLI-verified; exec rows say native.* Prepare sets `model_evidence: session-declared` with actual equal to requested; finalize re-reads them from the sidecar; the exec path sets `native` after generation. Both prose files say the in-session label is the declared one.
- *Final reply lists page, reply, backend and model, judge verdict, and the clarity changes.* The action's step (e) reads the typed row as before; the shipped prompt (unchanged) still asks the editor for the 3 to 5 changes, corrections and open questions, and the reply stage is required non-blank.
- *Foreign backends keep the single-call path.* The exec path's probe, generation, argv, metadata and supervisor code is unchanged in the diff; its callers only pass `nil, ""` and `[]string{stageName}` to the generalized preservation and hand the tail to `secondOpinionPublication`. The pre-flight probe (13 named tests including the four REQ-701 tests) passed three times in the builder's worktree and runs again in the gate.
- *Prose, guide, CLI, tests and changelog as one release.* Action and guide sections rewritten; the changelog entry, version bump and mirrors are the finalization transaction below.

**Behavior change judged (D-05):** in the exec path a reply-publish failure after the page publish now keeps the staging page instead of removing it. Accepted: it extends REQ-701's never-delete-a-finished-page rule to one more branch, no test pinned the old removal, and the kept file is named in `reason`.

**Scope comparison:** Implementation Summary files equal the Scope's five builder files; the Scope's seven integrator-only paths are release and lesson files. Files the Scope said would not be touched (prompt template, reference action, SKILL.md, help.md, launcher, commands.go, supervisor, adapters) are absent from the diff.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` (the repository gate, unpiped, under `full-gate.lock`) at merge 4362eb61 and again at the re-merge c70b9b72; the GREEN probe `bash do-work/runs/work-2026-10-19-152525/REQ-700-probe.sh` through `advance` at both merges; the builder's package runs `go test -count=1 ./internal/toolboxcommands ./internal/resultmodel` (ok, 17 s and 0.6 s), `gofmt -l`, `go vet`, and the four contract scripts in the worktree; the reviewer's `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands` (ok, 12.9 s) and `./internal/resultmodel` at 4362eb61.
**Result:** ✓ All passing (gate: 427 queue-kanban tests and 942 do-work-cli tests plus every shell contract; probe: both session-form tests passed at both merges)

**Repository gate retry:** at 4362eb61 run 1 exited 1 after 116 s on the 30 s per-file budget only (`skills/do-work-board/tools/queue-kanban` `strict_behavior_regression_test.go` accumulated 30.32 s, a package this REQ does not touch; 1-minute load 1.89 at start), rerun exited 0 after 136 s (the same file 20.35 s); at c70b9b72 run 3 exited 0 after 168 s with no retry.

**Red-green validation:** (traced to `## Red-Green Proof`: RED = a fake `claude` on PATH records zero editorial-pass invocations for the session form, which the base cannot satisfy because it execs the backend; GREEN = the same invocation hands back the filled prompt and reserved output path, and finalize publishes with every guarantee)
- `TestSecondOpinionSessionPrepareNeverExecsOwnBackend`: ✗ before implementation (`unexpected argument "--prepare"` inside the usage refusal, no typed row) → ✓ after
- `TestSecondOpinionSessionFinalizePublishesAgentWrittenPage` (subtests publishes, judge pass inserts banner, original modified keeps stage, never written cleans up): ✗ before (`unexpected argument "--finalize"` in all four) → ✓ after
- The captured RED case is the slash command in a live session; the implemented tests use the CLI's typed invocation with a fake backend as the nearest harness-runnable equivalent, which is what the proof's own GREEN sentence names.

**New tests added:**
- `TestSecondOpinionSessionPrepareNeverExecsOwnBackend`, `TestSecondOpinionSessionFinalizePublishesAgentWrittenPage` (four subtests), helper `secondOpinionFakeLogger` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`

**Existing tests updated (cross-REQ impact):** none; REQ-701's four tests and the kept-staging assertion stay green by name inside the pre-flight probe.

**Heavy verification plan:**
- Range: e11d3fca5a99ae8c51de48f531e9df7f07354746..c70b9b72deb421c93fd9da459c4020f27249f7a9
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the three Go files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three Go files matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 96%** (95% at 4362eb61; 96% after the delta re-review of the M1/M2 fix at c70b9b72) | 2026-10-19T17:16:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** M1 `actions/ai-report.md:65` the FINAL-REPLY-to-`staging_reply` instruction is attached only to the subagent branch, so an own-turns pass can skip the reply file and finalize refuses (files kept) — impact-user-visible → report only; M2 `actions/ai-report.md:63` "delegates all execution and publication to the canonical CLI" is stale for the Session form — impact-negligible → report only; M3 `docs/ai-report-guide.md:109` "and the reply says so" has no code or action rule behind it — impact-negligible → report only; M4 `docs/ai-report-second-opinion-prompt.md` says the launcher publishes "under the verified backend/model name", untrue on the session path — impact-negligible → report only; M5 D-02 (reply missing, page kept) and D-06 (supersession guard) have no pinning test — impact-negligible → report only; M6 `lessons-do-work-cli.md:5` (inherited, REQ-701's M2) still describes probe-and-generation identity for every run; reconciled only by the proposed `session-is-the-backend` bullet — impact-negligible → report only. Nits: N1 `actions/ai-report.md:69` terminal-default sentence ignores `--finalize` — impact-negligible → report only; N2 finalize-without-sidecar finding text has an empty backend prefix — impact-negligible → report only; N3 never-written cleanup also removes an agent-written reply stage without saying so — impact-negligible → report only. Anti-bloat count: 1 (`secondOpinionBacklinkRequirement`, deduplicates the backlink literal); every other addition is on the brief's named list.
**Acceptance:** Pass — Implementation and Integration stages: `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands` ok 12.9 s and `go test -count=1 ./internal/resultmodel` ok at 4362eb61; Deployment and Live acceptance unassessed
**Restatement sweep:** redefined the `ai-report-second-opinion` argument contract (`--prepare`, `--finalize`, usage string), the run-row schema (`staging_page`, `staging_reply`, `asset_prefix`, `prompt`, `model_evidence`, state `prepared`), the meaning of `mode: session` (prepare and finalize phases, no exec of the agent's own CLI), the model-identity rule (session rows declared, exec rows native), the `{{OUTPUT_FILE}}` fill (absolute stage path in-session), and the staging-leaf cleanup rule (leaves removed after both publishes, D-05); wave end, inherited from REQ-701's recorded line: the `-unverified` filename contract, the row schema (`seen_models`, state `unverified`), the `AI-REPORT-SECOND-OPINION-<STATE>` code set, and the identity rule (mismatch publishes unverified, a subagent model is not ambiguity). Swept `skills/` and `_dev/` (no `docs/` at the root): agreeing — `actions/ai-report-reference.md:148` (finalize is the launcher's phase, so "only the launcher may add a pass-gated banner" holds; `#second-opinion-form` anchor kept), `docs/ai-report-guide.md:72` usage line, `SKILL.md:3,24` trigger words, `actions/help.md` (no second-opinion line), `tools/do-work-cli.sh:10` (prepare reads the same env var), `_dev/tests/do-work-cli-launcher-behavior.sh:52,252` (export only), `lessons-do-work-cli.md:7` (REQ-701 bullet still true), CHANGELOG 0.306.0 and 0.306.3 (history), `commands_test.go` 10-handler pin, `result_model_test.go` round trip, and the inherited `-unverified`/`seen_models`/`-UNVERIFIED` contract (code at `secondOpinionPublication` 654-659, action line 73, guide line 107, REQ-701's tests green); stale — M2 (`ai-report.md:63`), M4 (prompt template OUTPUT paragraph), M6 (`lessons-do-work-cli.md:5`), N1 (`ai-report.md:69`)
**Suggested testing:** 6 items
**Follow-ups created:** None (9 findings report only)

*Reviewed by review-work action (full report and the delta re-review in `do-work/runs/work-2026-10-19-152525/REQ-700-review.md`; M1 and M2 were fixed on the builder branch in 2ba229e4 and re-merged as c70b9b72, so the remaining findings are M3 to M6 and N1 to N3, all report only)*

## Heavy Verification Plan

- Base: e11d3fca5a99ae8c51de48f531e9df7f07354746
- Target: c70b9b72deb421c93fd9da459c4020f27249f7a9
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — `result_model.go`, `report_second_opinion.go`, `report_second_opinion_test.go` matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three Go files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: c70b9b72deb421c93fd9da459c4020f27249f7a9
- Execution revision: c70b9b72deb421c93fd9da459c4020f27249f7a9 (detached checkout `.git/work-run-work-2026-10-19-152525/drain-head-REQ-700`, removed afterwards; `QUEUE_KANBAN_BROWSER` set)
- do-work-cli-integrations: executed, exit 0, 82 s
- staged-skills: executed, exit 0, 50 s
- updater: executed, exit 0, 76 s
- installer: executed, exit 0, 33 s
- Drain wall 245 s, outcome success, no `HEAVY-RUN-LANE-SKIPPED`. (An earlier drain attempt at 4362eb61 never executed a lane: the lane argv was passed as one quoted token and the runner refused with `HEAVY-RUN-USAGE`; the review fix moved the target before a retry.)

## Lessons Learned

**What worked:** Splitting the launcher at the generation call and carrying the before-snapshot across the process boundary as digests in a sidecar kept every guarantee without a second copy of the checks; the pre-flight probe naming REQ-701's four tests by name made the shared file safe to refactor right after REQ-701's release; serializing the two builders on the same Go file avoided a merge conflict entirely.
**What didn't:** The first qualification pass refused with 35 `QUALIFY-CLAIMED-PATH-MISSING` findings because the Implementation Summary put Go identifiers in backticks and the qualifier reads every backticked token as a path; only file paths may be backticked there. Gate run 1 went red on the 30 s per-file budget in `queue-kanban` (30.32 s, an untouched package) and passed on the one retry. The first heavy drain executed nothing because the lane flags were passed as one shell word under zsh. Review found the one real gap in prose, not code: the own-turns branch of the action never said to write the final reply to `staging_reply`, which finalize requires.
**Worth knowing:** The sidecar is `<stem>.json` beside the dotfile staging page, mode 0600, digests only. Finalize derives the bundle from the staging page's parent and refuses a supersession jump (D-06). In the exec path a reply-publish failure after the page publish now keeps the staging page (D-05). Report-only follow-ups: the prompt template still says "verified backend/model name"; D-02 and D-06 have no pinning test; the never-written cleanup also removes an agent-written reply stage; `help.md` has no second-opinion line.

## Orientation

Now the session form of `ai-report second-opinion` runs its editorial pass inside the agent's own session (or a subagent or workflow it starts) between the CLI's `--prepare` and `--finalize` phases, instead of a hidden `claude -p` subprocess; lives in the do-work-cli toolbox commands (`report_second_opinion.go`, `result_model.go`) and the toolbox `ai-report` action and guide. [MAP CHANGED] the second-opinion CLI contract gains two phases and five row fields (`staging_page`, `staging_reply`, `asset_prefix`, `prompt`, `model_evidence`); the three touched primes (action-files, shell-commands, releases) reference no path this change moved, so none is stale.

## Timing

Observed 2026-10-19T15:32:15Z to 2026-10-19T17:16:24Z: 1h 44m 09s total, 58m 46s attributed across 8 events, 45m 23s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 16m 49s | 1 |
| builder-work | 15m 21s | 1 |
| verification-gate | 15m 04s | 3 |
| planning | 7m 16s | 1 |
| exploration-preflight | 4m 01s | 1 |
| handback-merge | 15s | 1 |

Slowest stage: review / independent review (Route C, wave-end restatement sweep) plus delta re-review, 16m 49s, outcome success.
