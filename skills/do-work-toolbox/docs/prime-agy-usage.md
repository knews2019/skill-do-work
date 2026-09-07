# Prime: AGY usage

AGY is the Antigravity CLI (`agy`). Read this before invoking it from any agent; project orchestration and quality gates remain owned by the consuming project.

## Read first

- `agy --help`, `agy models`, `agy changelog` — installed flags, current model catalog, version-specific behavior. Resolve the executable from PATH; do not copy another machine's binary path.
- [Headless guide](https://antigravity.google/docs/cli/headless/) — process I/O, output envelopes, schema output and conversation transport.
- [Boost](https://antigravity.google/docs/boost/) and [subagents](https://antigravity.google/docs/subagents/) — delegated work, context isolation and workspace options.

## Model and permissions

- Select the requested model from the live catalog and pin its complete slug with `--model`. `gemini-3.8-flash-high` was the source investigation's working example on 2026-09-07, not a permanent default. Avoid a conflicting `--effort` override.
- `agy --model <model-slug>` opens the interactive UI. `/model` opens its picker; `/model <model-slug>` saves a default. In CLI 1.1.27, trailing prompt text instead made a one-prompt switch. Prefer `--model` for a session override.
- Inspect effective selection without a model turn using `agy --model <model-slug> --output-format json -p '/model'`. Inspect only needed fields from `agy -p '/config' --output-format json`; configuration can contain account data.
- `--dangerously-skip-permissions` auto-approves tool requests. Use it when unattended tool execution is authorized; it does not select a model, submit work, detach the process, or remove hardcoded protection boundaries.

## Submit and verify work

Run from the target repository. Replace both placeholders before using this invocation:

```sh
agy --model '<model-slug>' --print-timeout 15m --output-format stream-json \
  -p '<concrete task with absolute repository and input paths, constraints, outputs and acceptance checks>' </dev/null
```

- **Include `-p` to submit once and exit.** Omitting it opens the UI. Choose an explicit `--print-timeout` appropriate to the task; the source CLI default was five minutes, and nested shell timeouts are separate.
- **Scope goes in the prompt as well as cwd.** The source investigation observed a vague file probe search outside the repository; naming the exact absolute target and repository boundary produced a direct file read. Inspect tool events for drift.
- Keep NDJSON stdout and stderr in separate files. `text` gives prose; `json` gives a terminal envelope; `stream-json` gives observable execution.
- Inspect `init` for cwd, explicit model and permission mode; `step_update` for actual tool work; the terminal `result` event for `result.status`. Require `SUCCESS`, process exit success and the task's artifact checks. Startup alone proves no completed work.
- For schema output, use the documented `--json-schema` option and inspect `structured_output` in the terminal envelope; verify the installed behavior before depending on it.

## Background work and continuation

- Allocate a unique run directory containing the task prompt, event stream, stderr, PID and conversation ID. Name one owner to monitor progress, check completion and report the outcome.
- For a programmatic launcher, use `subprocess.Popen` with an argument list, `cwd` set to the absolute repository root, `stdin=DEVNULL`, separate output files and `start_new_session=True`. Pass prompt text as an argument, never executable shell text. Launching is not completion: collect process status, terminal result and output checks.
- Resume only an inactive conversation with `agy --model <model-slug> --conversation <conversation-id> --print-timeout 15m -p '<continuation task>'`. Reapply the authorized permission mode as needed. Prefer the recorded ID over `--continue`, which can select an unrelated recent conversation; never attach a second writer to an active run.
- Persistent transport uses `--input-format stream-json --output-format stream-json`, with a line shaped as `{"event":"user","message":{"content":"..."}}` per turn. Wait for its `result` before the next prompt and close stdin when finished. Consult the headless guide before implementing this transport.

## Boost and delegation

- Enter `/boost <bounded task>` in AGY's interactive TUI for a difficult audit or fix. Google documents Boost on paid plans; check availability on the installed account. A headless file probe does not establish headless Boost support.
- Supply workers with absolute paths, scope, constraints and acceptance checks: subagents do not inherit the parent's conversation history. Keep one coordinator responsible for shared integration, builds and queue transitions.
- Inventory HEAD, tracked changes and relevant untracked files before delegation; verify each worker sees the intended candidate. Do not assume isolated worktrees received dirty or untracked state.
- A pinned parent model does not prove worker or nested subprocess model selection. Verify each dispatch before claiming every phase used the requested model. A DeepCoder mention alone does not prove Boost ran.

## Traps

- **Runtime names are not interchangeable.** Do not assume Claude Code's injected `Workflow` runtime or Gemini CLI syntax exists in AGY. Inspect available tools; a workflow file requiring injected globals is not automatically a standalone executable.
- **Conversation and pipeline IDs identify different state.** Preserve both when integrating with a project runner, and let the project's existing owner advance its gates.
- **Status-line arguments are shell commands.** Prose saved as a command caused `sh: include: command not found` in the source investigation. `/statusline reset` restores the built-in display; `/statusline off` disables it. See the [command reference](https://antigravity.google/docs/cli/commands/statusline).

## Verification boundary

Adapted from the operator-supplied WSD AGY usage prime dated 2026-09-07. That investigation reported CLI 1.1.27 help/catalog/config checks, explicit Gemini 3.8 Flash High selection, a scoped file read, streaming `SUCCESS` and exit 0. These are historical observations, not checks performed in a consumer installation. Headless, Boost and subagent docs were rechecked when packaging this reference on 2026-09-07; no AGY model task was launched. Boost execution, worker model routing, dirty-state transfer, schema output and persistent transport remain untested here.
