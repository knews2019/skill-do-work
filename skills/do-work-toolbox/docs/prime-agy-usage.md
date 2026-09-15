# Prime: AGY usage

AGY is the Antigravity CLI (`agy`). Read this before invoking it from any agent; project orchestration and quality gates remain owned by the consuming project.

## Read first

- `agy --help`, `agy models`, `agy changelog` — installed flags, current model catalog, version-specific behavior. Resolve the executable from PATH; do not copy another machine's binary path.
- [Headless guide](https://antigravity.google/docs/cli/headless/) — process I/O, output envelopes, schema output and conversation transport.
- [Boost](https://antigravity.google/docs/boost/), [subagents](https://antigravity.google/docs/subagents/) and [CLI background work](https://antigravity.google/docs/cli/subagents) — escalation, model tiers, context isolation and monitoring.

## Model and permissions

- **Default: the latest Gemini Flash release available in `agy models`, at Low effort for routine work or Medium for work requiring more reasoning.** Choose the newest Flash release first, then its effort variant; catalog order is not a ranking. Use Medium if that release has no Low variant. Do not select an older release just to obtain Low, invent a missing slug, or silently use High, Pro, Claude, GPT, or a saved default when Low/Medium is unavailable. Report an unavailable or ambiguous selection before submitting work. An explicit user model/effort request overrides this default.
- Pin the complete catalog slug with `--model` on every launch and continuation. Omit `--effort` when the slug already selects the desired effort. A bare `--effort` does not pin the model family; a conflicting override can resolve to another variant. The live catalog checked on 2026-09-15 offered `gemini-3.8-flash-low`, `gemini-3.8-flash-medium` and `gemini-3.8-flash-high`; these are dated examples, not permanent defaults.
- `agy --model <model-slug>` opens the interactive UI. `/model` opens its picker; `/model <model-slug>` saves a default. In CLI 1.1.27, trailing prompt text instead made a one-prompt switch. Prefer `--model` for a session override.
- Before submitting work, inspect effective selection without a model turn using `agy --model <model-slug> --output-format json -p '/model'`; `/effort` is also a read-only query. Require the resolved model and effort to match the intended selection. Record them with the run and recheck after a switch or continuation; an alias, deprecated model replacement, or effort override is not proof of the requested selection. Inspect only needed fields from `agy -p '/config' --output-format json`; configuration can contain account data.
- `--dangerously-skip-permissions` auto-approves tool requests. Use it when unattended tool execution is authorized; it does not select a model, submit work, detach the process, or remove hardcoded protection boundaries.

## Submit and verify work

Run from the target repository. Resolve `<model-slug>` using the default or an explicit override above, then replace both placeholders before using this invocation:

```sh
agy --model '<model-slug>' --print-timeout 15m --output-format stream-json \
  -p '<concrete task with absolute repository and input paths, constraints, outputs and acceptance checks>' </dev/null
```

- **Include `-p` to submit once and exit.** Omitting it opens the UI. Choose an explicit `--print-timeout` appropriate to the task; the source CLI default was five minutes, and nested shell timeouts are separate.
- **Scope goes in the prompt as well as cwd.** Name the absolute repository and input paths, current project phase, remaining deliverable, owned paths, known command interfaces and acceptance checks. For a handoff, distinguish already-landed work from the remaining assignment. The source investigation observed a vague file probe search outside the repository; naming the exact target and boundary produced a direct file read.
- Keep NDJSON stdout and stderr in separate files. `text` gives prose; `json` gives a terminal envelope; `stream-json` gives observable execution.
- Inspect `init` for cwd, explicit model and permission mode; `step_update` for actual tool work; the terminal `result` event for `result.status`. Treat AGY completion, underlying command outcomes and artifact correctness as separate checks. Require `SUCCESS` and process exit success, then independently verify the task's acceptance checks. Inspect command findings even after a successful AGY result; give each a disposition: fixed and rechecked, accepted with a stated reason, or unresolved with an owner and next action. Report accepted warnings and unresolved findings in the handback; unresolved acceptance failures mean the task is incomplete.
- Read stderr and any `denied_actions` alongside the result. The CLI changelog records that a print timeout can return partial output with exit 0 and a warning; exit 0 alone does not establish completion. A missing terminal result, truncation, or denied required action leaves work to reconcile and verify.
- For schema output, use the documented `--json-schema` option and inspect `structured_output` in the terminal envelope; verify the installed behavior before depending on it.

## Background work and continuation

- Allocate a unique run directory containing the task prompt, event stream, stderr, PID and conversation ID. Name one owner to monitor progress, check completion and report the outcome.
- Monitor tool events for scope drift and progress toward the remaining deliverable. Repeated reads without new evidence or progress call for reassessing and narrowing the assignment. Before interrupting, inspect fresh tool activity, worktree changes, new commits and the project's canonical phase; work may have landed since the last progress message. Record a deliberate interruption separately from an AGY failure.
- For a programmatic launcher, use `subprocess.Popen` with an argument list, `cwd` set to the absolute repository root, `stdin=DEVNULL`, separate output files and `start_new_session=True`. Pass prompt text as an argument, never executable shell text. Launching is not completion: collect process status, terminal result and output checks.
- After confirming the process is inactive, reconcile landed work through the project's existing recovery and phase owners. Resume only the remainder with `agy --model <model-slug> --conversation <conversation-id> --print-timeout 15m -p '<continuation task>'`. Include the reconciled state and remaining checks; reapply the authorized permission mode as needed. Prefer the recorded ID over `--continue`, which can select an unrelated recent conversation; never attach a second writer to an active run.
- Persistent transport uses `--input-format stream-json --output-format stream-json`, with a line shaped as `{"event":"user","message":{"content":"..."}}` per turn. Wait for its `result` before the next prompt and close stdin when finished. Consult the headless guide before implementing this transport.

## Escalation and additional capabilities

- **Low → Medium → High on the same latest Flash release.** Use Medium directly when the task needs cross-file reasoning. Escalate to High for a difficult diagnosis or failed acceptance check that better reasoning can address. First distinguish a reasoning failure from missing inputs, denied tools, quota exhaustion or a broken environment; higher effort does not repair those. Record the reason and exact new slug. Return to Low/Medium once the hard part is resolved; do not carry High into unrelated routine work.
- **Boost for a bounded problem that benefits from multiple investigations and independent verification.** Pin Flash High for the coordinator, then enter `/boost <bounded task>` in AGY's interactive TUI. High is a reasoning setting; Boost is a separate multi-agent workflow. Google documents Boost on paid plans. Check installed availability, supply the reproduction or failing check and desired result, and verify that the workflow actually dispatched. The documented CLI invocation is interactive; do not prescribe `-p '/boost ...'` as supported without checking the installed behavior. If unavailable in the execution surface, use a normal pinned High run with focused verification.
- **Focused subagents for separable work.** Use the installed agent tools for a bounded investigation or independent check while the coordinator has useful work to do. Custom subagent `model` values are tiers (`inherit`, `flash`, `pro`), not full CLI slugs. Prefer `inherit` from the selected Flash parent, then verify the effective worker model and effort. If exact selection cannot be established, use a separately pinned AGY process or keep the work in the coordinator. Do not infer that `flash` means the newest release at the desired effort.
- **Teamwork for an authorized project spanning multiple milestones.** [Teamwork](https://antigravity.google/docs/teamwork/) uses `/teamwork-preview <task>`, a scoping interview and a reviewable prompt before execution, on paid plans. Use it when sustained team coordination fits the assignment, not as an automatic retry for a small failure. Specify the intended repository/workspace explicitly and retain the consuming project's orchestration and release ownership.
- **Observe before escalating.** Use `/agents` for workers, `/tasks` for background commands and `/usage` for model quotas. `/usage` is also available as a read-only print query. Quota totals do not identify which model a particular run used. Consult the [CLI reference](https://antigravity.google/docs/cli/reference/) for installed controls; structured output and conversation continuation above help make verification reliable without increasing reasoning effort.

### Delegation boundaries

- Supply workers with absolute paths, scope, constraints and acceptance checks: subagents do not inherit the parent's conversation history. Keep one coordinator responsible for shared integration, builds and queue transitions.
- Inventory HEAD, tracked changes and relevant untracked files before delegation; verify each worker sees the intended candidate. Do not assume isolated worktrees received dirty or untracked state.
- Apply the Flash default and any explicit user constraints to workers and nested subprocesses too. Escalation authorizes the needed effort or workflow, not a silent switch to another model family. A pinned parent model does not prove worker selection, and Boost or Teamwork may manage their own routing. Record observable dispatch models and disclose any routing that cannot be verified; do not claim an all-Flash run without evidence. If the user requires every worker to be pinned, use individually verified dispatches instead of opaque routing. A DeepCoder mention alone does not prove Boost ran.

## Traps

- **Runtime names are not interchangeable.** Do not assume Claude Code's injected `Workflow` runtime or Gemini CLI syntax exists in AGY. Inspect available tools; a workflow file requiring injected globals is not automatically a standalone executable.
- **Conversation and pipeline IDs identify different state.** Preserve both when integrating with a project runner, and let the project's existing owner advance its gates.
- **Status-line arguments are shell commands.** Prose saved as a command caused `sh: include: command not found` in the source investigation. `/statusline reset` restores the built-in display; `/statusline off` disables it. See the [command reference](https://antigravity.google/docs/cli/commands/statusline).

## Verification boundary

On 2026-09-15, `agy --help`, the live model catalog and `agy changelog` (entries through 1.2.3) were checked. Read-only `/model` probes with each Gemini 3.8 Flash Low/Medium/High slug returned the matching `command.data.id` and `command.data.effort`, `SUCCESS`, exit 0 and zero model turns/tokens. Google's headless, Boost, subagent, Teamwork and CLI reference docs were rechecked. The default and escalation choices above are this guide's policy; these checks did not execute a model task, Boost or Teamwork, or verify worker routing.

Adapted from the operator-supplied WSD AGY usage prime dated 2026-09-07. That investigation reported CLI 1.1.27 help/catalog/config checks, explicit Gemini 3.8 Flash High selection, a scoped file read, streaming `SUCCESS` and exit 0. Headless, Boost and subagent docs were rechecked when packaging this reference on 2026-09-07; that packaging check launched no AGY model task.

On 2026-09-09, a repository-scoped headless REQ-624 run used `gemini-3.8-flash-high` with separate event/error logs and explicit conversation continuation. The coordinator intentionally interrupted the first run after finalization had already committed. The continuation returned `SUCCESS` and exit 0, but its cleanup command reported findings that independent verification corrected. Evidence preservation and stale references also needed correction. See the [execution evidence](https://github.com/knews2019/skill-do-work/blob/main/ai-reports/2026-09-09_1224_REQ-624-verification-review/evidence/agy-execution-summary.json) and [independent verification](https://github.com/knews2019/skill-do-work/blob/main/do-work/archive/UR-130/assets/req624-evidence/parent-verification.json). These observations motivate the guidance above; they do not establish a speed or cost improvement or verify a consumer installation. Boost execution, worker model routing, dirty-state transfer, schema output and persistent transport remain untested here.
