# Claude Code instruction retention experiment

The tested Claude Code session completed real manual compaction, made no subsequent file reads, and followed the observed function-naming rule. Its compaction summary explicitly preserved that rule. This single run does not establish complete instruction retention or a general restoration failure.

## Harness and controls

Observed on 2026-11-11 at 01:20:24–01:21:20 UTC in the assigned isolated worktree. Executable: `/Users/t2/.local/bin/claude`. Installed version: **2.1.322 (Claude Code)**. The initialization and assistant events report **claude-fable-5-1**. No model override or fallback was supplied. Result accounting also reports auxiliary `claude-haiku-5-5` usage; this is not an observed delegated worker. Result events report zero spawned subagents.

`run-metadata.json` records the exact argv, cwd, start time, exit code, and turn outcomes. `cli-help.txt` preserves locally discovered controls. No external documentation or global logs/configuration were inspected. Scoped `.claude/settings.json`, `.claude/settings.local.json`, and `.mcp.json` were absent.

The probe used `--safe-mode --setting-sources '' --strict-mcp-config --tools Read --allowedTools Read --permission-mode dontAsk --no-session-persistence`, with streamed JSON input/output and a $3 budget cap. Safe mode disables user customizations, hooks, and CLAUDE.md discovery; phase files were explicitly loaded instead. The initialization still lists built-in plugins and skills, no MCP servers, only the Read tool, and no additional directories. Authentication was handled by the installed CLI; no credentials were printed or manually inspected. There was no explicit effort override. Initialization reports per-turn effort active, but not its exact selected level; recorded usage shows zero thinking tokens.

## Repeatable procedure

1. Run `python3 _dev/experiments/REQ-718/claude/run-probe.py` from this checkout. It replaces capture and metadata files in this experiment directory and emits timestamped progress to stderr, without depending on the original main-tree run directory.
2. Run `python3 _dev/experiments/REQ-718/claude/verify-evidence.py` to check the capture against the repository inputs and regenerate `focused-checks.json`.

The runner sends the exact three prompts stored in `baseline-prompt.txt`, `compact-prompt.txt`, and `continuation-prompt.txt` within one CLI process. It waits for each successful result and requires a real `system/compact_boundary` before sending the continuation. Limits are 150 seconds per turn and 360 seconds overall, with bounded termination and process reaping. Actual elapsed time was 55.288 seconds; process exit was 0. Subsequent runs may choose a different default model: inspect their initialization events before comparing results.

For an offline progress regression check, run `python3 _dev/experiments/REQ-718/claude/test-replay-progress.py`. It simulates an absent run directory by making filesystem opens fail and checks timestamped stderr output without launching the CLI. The evidence verifier in step 2 also runs offline.

## Observations and evidence

All row references address the sanitized `live-events.jsonl`, which preserves actual tool calls, tool results, generated code, and compaction events.

- **F1 — Actual baseline loads.** Rows 3–18 contain eight Read calls and their successful, full-content results: `SKILL.md`, `actions/work.md`, `actions/work-implementation-reference.md`, and crews `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`, all beneath `skills/do-work/`. The verifier matches every numbered result line to the original file. `input-fingerprints.json` stores SHA-256 fingerprints. This is a paused real Step 6 procedure under a read-only experimental boundary, not a complete queue implementation run. The sample domain is testing; tdd, maintenance and caveman are false; security/dependency triggers and required lessons are absent.
- **F2 — Baseline rule compliance.** `coding-guardrails.md` section 5 requires descriptive names with at least two words for functions. The baseline response and result at rows 19–20 contain `function sumInvoiceAmounts(invoices)`. The prompt asks for invoice summation without prescribing the function name or restating the naming rule.
- **F3 — Real compaction.** The second prompt is exactly `/compact`. Row 25 is a harness-generated `system/compact_boundary`, trigger `manual`, reporting 36,066 pre-compaction tokens, 11,389 post-compaction tokens, 24,677 cumulative dropped tokens and 41,513 ms duration. Row 28 is a successful result with `local_command: compact`. These events establish the live boundary independently of model self-report.
- **F4 — Visible retained material.** Row 26 is the synthetic continuation summary emitted by the harness. It explicitly retains Step 6 crew loading, phase constraints and “Naming for Reach (two-word minimum for anything with reach...)”. The boundary also identifies the baseline assistant response as a preserved message. The summary incorrectly describes that baseline response as still pending, while explaining that preserved recent messages can supersede its account. The continuation nevertheless answers the new shipment prompt.
- **F5 — Continuation.** Rows 29–31 contain initialization, one assistant text response, and a successful result. There are **zero tool calls**, including zero Read calls. The generated function is `sumShipmentWeights(shipments)`, with local accumulator `shipmentTotal`. Its function name satisfies the observed rule. The continuation prompt does not remind the model of the naming rule or request rereads.

The capture omits account rate-limit events, credential-expiration metadata, and local socket paths. No credentials or unrelated historical transcripts are included. Those omissions do not remove phase reads, compaction metadata, summary content, or output behavior.

## Interpretation and acceptance limits

The observed absence of rereads is distinct from missing instructions: the naming rule remains visibly available in the synthetic summary, and subsequent output follows it. This run contains no observed failure to follow that rule. It does not show that every original file remains available verbatim after compaction, nor does token-count reduction identify exactly which file contents were discarded or restored internally.

Compliant output alone cannot prove causal retention: descriptive names are common, the task wording contains domain terms, and the preserved baseline response provides an example. The visible summary supplies stronger evidence of rule availability, but this is still one simple task and one manual boundary with a small initial context. Automatic compaction under pressure, repeated compaction, ordinary project/global settings, skill auto-discovery, other crew rules, and other models are **Untested**. Actual Codex/T3 coverage belongs to the separate worker and is **Untested in this slice**.

All Claude-slice acceptance observations were exercised: actual baseline reads, real compaction, subsequent read behavior, harmless-rule behavior, exact controls/prompts, and scoped uncertainty. There is no capability blocker for this tested Claude configuration. No restoration hook or runtime fix is supported by these observations.

## Focused validation

`verify-evidence.py` exited 0: eight complete read-result matches, unchanged input fingerprints, real manual boundary, both observed descriptive function names, summary rule availability, zero continuation reads, sanitized metadata, runner syntax, and successful reaped CLI process. An initial ad hoc matcher failed because Read numbers the final empty line after a trailing newline; the verifier accounts for that representation and then matched all contents. This was a checking-script mismatch, not a harness read failure.

`git diff --check` is the whitespace check for the artifact change. No production code, shipped files, runtime hooks, dependency files or tests were changed. No full canonical gate was run by this worker; the integrator owns it.
