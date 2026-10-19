---
id: REQ-700
title: 'Second opinion from the session model runs as a visible in-session prompt, not a hidden shell subprocess'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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
