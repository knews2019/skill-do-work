---
id: REQ-700
title: 'Second opinion from the session model runs as a visible in-session prompt, not a hidden shell subprocess'
status: pending
created_at: 2026-10-19T14:43:26Z
user_request: UR-158
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
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
