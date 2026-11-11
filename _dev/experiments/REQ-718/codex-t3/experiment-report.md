# Codex/T3 instruction retention probe

Real compaction was inaccessible through the controls exposed to this worker. Discovery and the implementation-phase baseline ran in the actual Codex/T3 worker. Post-compaction retention, rereads, and compliance are **Untested**. This report does not establish an instruction-retention failure or success.

## Runtime and boundary

The read-only T3 APIs report server **0.1.7**, provider instance **codex**, model **gpt-6.1-sol**, reasoning effort **high**, runtime mode **full-access**, and interaction mode **default**, on Darwin arm64. `runtime-evidence.json` preserves the scoped responses and exposed T3 tool inventory. The active developer context also specifies approval policy never and danger-full-access sandboxing. Those injected settings were not independently obtained from provider configuration.

The worker operates on branch `t3/execute-codex-brief` in `/Users/t2/.t3/worktrees/skill-do-work2/t3-d6e60aeb`, initially at `a867b1a08d38dfb5605602a6a0ac5408c0368eb9`. The T3 activity response independently associates this path and branch with the actual Codex subagent. The final committed SHA is supplied in the durable handback.

`/opt/homebrew/bin/codex --version` returned **codex-cli 0.167.0**. This identifies the installed PATH executable, not conclusively the binary backing the active T3 session. `command -v t3` returned no path, causing the combined discovery command to exit 1 before `t3 --version` ran. The T3 server version comes from its API instead.

`codex app-server daemon version` exited 1: “failed to connect to /Users/t2/.codex/app-server-control/app-server-control.sock”, followed by “No such file or directory (os error 2)”. This permitted version query did not inspect global logs or configuration. The active provider build, auto-compaction threshold, context size, compaction prompt, and merged provider configuration remain unverified. Repository `.codex` and `.t3` directories were absent; scoped discovery returned exit 2 for those missing paths. No global configuration was read or changed.

## Actual reads and harmless-rule baseline

The selected phase is **Step 6: Implementation**, under the durable brief. The worker read the full working request, including user answers F5–F8, and `work-implementation-reference.md`. It loaded general, coding-guardrails, shared-principles, communication-style, and testing crews. `domain: testing` requires testing; `tdd: false` does not require RED/GREEN. Background-agents was read for durable handback, and anti-slop for report prose. No required lessons were listed, and no applicable lesson satellite was required for these experiment artifacts.

The skill body `skills/do-work/SKILL.md` and `actions/work.md` were also read. These were explicit shell reads, not evidence of automatic skill activation. The lifecycle owner's metadata operations and spawning instructions do not apply to this isolated worker: the user explicitly prohibited delegates and REQ edits.

`transcript-evidence.json` contains exact excerpts from actual T3 activity items, with IDs, timestamps, positions, and truncation flags. Event 4 records the request/procedure read. Event 9 records the complete general/guardrails/shared-principles read, including the selected rule. Events 22 and 23 record subsequent complete crew and skill reads. Initial combined output was truncated (event 5); separate reads supplied the crew bodies afterward. Tool-output truncation is not a compaction boundary. `source-fingerprints.json` records the exact source contents by SHA-256 at this checkout.

The selected harmless rule is **Naming for Reach**, section 5 of coding-guardrails: new files and identifiers with reach need descriptive names with at least two words. The actual file-change event 28 creates `baseline-observation.md`, satisfying that narrow filename observation. Event 14 is a plain, specific status update consistent with communication-style. These examples show baseline compliance only. The brief and developer instructions also favor descriptive names and plain communication, so causation cannot be attributed solely to the loaded crew file. There was no blinded control.

## Compaction capability blocker

The available catalog contained **65 T3 tools**. None described a compaction operation. The relevant controls expose thread reads, model configuration, message delivery, interruption, and context transfer through fork/merge. None documents a way to compact the running provider session or reports the provider compaction configuration. The exact discovered names and descriptions are preserved in `runtime-evidence.json`.

CLI help discovery exposed app-server tooling but no command targeting compaction of this active T3 worker. The version query could not connect to its default daemon socket. No supported connection from the T3 thread identifier to a provider compaction endpoint was established. Starting a separate CLI agent would test a different harness. Sending `/compact` as ordinary thread text would not establish control semantics. Neither was attempted.

No manual compaction invocation was available or attempted. No runtime compaction event was observed in the captured activity prefix, positions 0–32. That does not prove the provider emits every such event into T3 activity. This is a scoped capability blocker, not a claim that Codex or T3 can never compact. Artificially filling the context was not attempted: no threshold was known and a large uncontrolled run would not supply a reliable, bounded test.

## Acceptance coverage and interpretation

| Check | Actual outcome |
|---|---|
| Real phase, procedure, crew reads and baseline | Observed in actual worker events; narrow descriptive filename compliance |
| Real compaction boundary | Untested: no supported active-session control established |
| Post-boundary file reads | Untested: no verified boundary |
| Post-boundary harmless-rule compliance | Untested: continuation not run |
| Exact harness/model/configuration | T3 server/model/options observed; installed CLI version observed; backing provider build and compaction settings unverified |
| Missing content vs absent rereads vs noncompliance | No loss established; no post-boundary reread or compliance classification possible |
| Claude coverage | Outside this worker's assigned slice; no Claude claim |

Before a verified boundary, repeated file reads only demonstrate explicit instruction loading. After a future boundary, absent rereads would not establish absent contents. A retained rule could survive as a summary without its source body. Noncompliance would require evidence the instruction was available at the decision point; otherwise the cause remains ambiguous. No restoration fix or hook is justified by these observations.

## Repeatability, checks, and cleanup

`probe-prompts.md` records actual discovery commands, the original dispatch evidence, the baseline action, and the unexecuted boundary/continuation procedure. Raw evidence is selected and excerpted rather than a full transcript export. It excludes reasoning items, unrelated skill documentation, environment preferences, and unrelated conversations. No credentials were found in the retained fields.

Focused validation checks JSON parsing, source fingerprints, referenced event evidence, artifact scope, and `git diff --check`. Exact outcomes are recorded in the handback after execution. These checks validate the evidence package, not retention. The integrator owns the canonical gate; the user reports its pre-build exit 0, and this worker did not rerun it.

Only experiment artifacts were added on this branch. No hooks, runtime instructions, queue files, global configuration, source in the main checkout, pushes, or PRs were changed. The sole main-checkout writes are the authorized progress log and handback. All commands completed, and no agent, daemon, background process, or server was started.
