# Claude implementation handback: compaction instruction retention

Committed live evidence is ready for serial integration. No capability blocker occurred for the tested Claude configuration. This handback does not mark the combined Claude plus Codex/T3 request complete.

## Operative checkout

- Branch: `t3/execute-claude-slice-work`
- Worktree: `/Users/t2/.t3/worktrees/skill-do-work2/t3-491e3e50`
- HEAD: `993c3befca2af4eeae8aad2fdba386ab7d16438a`
- Commit: `Record live Claude compaction and instruction retention evidence`
- Working tree: clean at handback.
- Primary report: `_dev/experiments/REQ-718/claude/experiment-report.md`.

## Observed result and acceptance coverage

Claude Code 2.1.322, model `claude-fable-5-1`, completed a single read-only Step 6 probe with safe mode, empty setting sources, no external MCP servers, Read as the only tool, and no transcript persistence. Exact argv/configuration scope is in `run-metadata.json`; prompts, local CLI help, input fingerprints and sanitized real stream are committed.

Baseline: eight actual successful Read calls loaded the do-work skill, action, implementation procedure and applicable general/coding/shared/communication/testing crews. Full tool-result contents match the fingerprinted source files. The observed function name `sumInvoiceAmounts` meets the existing descriptive naming rule.

Boundary: a harness `system/compact_boundary` event reports manual compaction, 36,066 pre tokens, 11,389 post tokens, and duration 41,513 ms. The command result reports success. Continuation made zero tool calls and generated `sumShipmentWeights`, also compliant with the observed function-naming rule. The emitted synthetic summary explicitly retains that rule. The preserved baseline output also supplies a naming example.

Claude-slice acceptance: baseline reads/rule, actual compaction, subsequent reads/compliance, repeatable prompts/version/configuration, evidence preservation, and distinctions/uncertainty are exercised. No missing-instruction or noncompliance failure was observed for the tested rule. Full verbatim file retention is not established. Lack of rereads is observed independently of available summary content. Compliant naming alone cannot establish causality.

## Exact file manifest

All 11 files are new; no modified/deleted source files or tests:

```
A	_dev/experiments/REQ-718/claude/baseline-prompt.txt
A	_dev/experiments/REQ-718/claude/cli-help.txt
A	_dev/experiments/REQ-718/claude/compact-prompt.txt
A	_dev/experiments/REQ-718/claude/continuation-prompt.txt
A	_dev/experiments/REQ-718/claude/experiment-report.md
A	_dev/experiments/REQ-718/claude/focused-checks.json
A	_dev/experiments/REQ-718/claude/input-fingerprints.json
A	_dev/experiments/REQ-718/claude/live-events.jsonl
A	_dev/experiments/REQ-718/claude/run-metadata.json
A	_dev/experiments/REQ-718/claude/run-probe.py
A	_dev/experiments/REQ-718/claude/verify-evidence.py
```

## Focused checks and actual outcomes

- `python3 _dev/experiments/REQ-718/claude/run-probe.py`: CLI exit 0, three successful turn results, observed compact boundary, 55.288 seconds, process reaped. Model-accounted session cost approximately $1.60, within the $3 cap.
- `python3 _dev/experiments/REQ-718/claude/verify-evidence.py`: exit 0. All eight Read results match complete source contents; fingerprints unchanged; manual boundary present; descriptive function names and retained summary rule present; no continuation reads; excluded metadata absent; runner Python syntax valid; CLI clean exit/reaping confirmed. Durable results: `focused-checks.json`.
- Initial ad hoc read matcher: failed due to Read's numbered trailing empty line. Corrected representation handling, then verified every content line. This is not recorded as a harness failure.
- `git diff --check` and `git diff --cached --check`: exit 0.
- Staged JSON/JSONL/Python syntax parsing and write-boundary assertion: passed; every committed path is under the assigned Claude directory. Reviewed every artifact, including all captured tool contents and summary.
- `git diff --cached --stat`: 11 new files, 705 inserted physical lines. Evidence includes long JSONL records.
- Canonical full gate: not run by worker, per brief. Integrator owns merged-state verification; pre-build gate was supplied as exit 0.

## Uncertainty and blockers

No blocker for manual Claude compaction in this configuration. Ordinary project/global customizations, automatic/repeated compaction, large-context pressure, skill automatic restoration, other models and other crew-rule compliance are **Untested**. The stream does not expose an exact effort level or the entire post-compaction model context; zero recorded thinking tokens is not proof of an effort setting. Global configuration/logs were not inspected. Actual Codex/T3 is **Untested by this worker** and remains the other slice's responsibility. This run supplies no basis for a categorical retention-failure claim or restoration fix.

## Decisions

- D-01 — DECIDE & STATE: use installed CLI safe mode with explicitly loaded repository files and read-only tools. This bounds custom hooks and external configuration while preserving a real harness compaction test. Conclusions are scoped accordingly.
- D-02 — DECIDE & STATE: use supported streamed `/compact` mid-phase rather than artificial filler. Require the harness boundary event before continuation. This gives direct compaction evidence with a bounded cost/time budget.
- D-03 — DECIDE & STATE: retain sanitized actual event records, removing account rate-limit events, credential expiration and socket paths. Preserve complete applicable file results, synthetic summary and generated behavior for independent review.

## Discovered Tasks

None. Broader experiments above are untested scope, not newly queued work.

## Lessons consulted

No required lessons were listed; the durable brief says none currently required. Read the implementation reference and applicable general, coding-guardrails, shared-principles, communication-style, testing, background-agents and anti-slop crews; also consulted prompt-injection guidance and shell/release primes. No touched utility or trap requires a lesson satellite for this artifact-only slice.

## Process cleanup and boundaries

PLAN/APPLY/UNIFY completed for this isolated slice; the working REQ was not edited. No delegates/threads, hooks, global config edits, pushes, PRs, lifecycle mutations, shipped-file changes or main-tree source writes. The CLI exposed only Read, reports zero spawned subagents, and was reaped after exit 0. No command process remains from the probe. Commit hooks were disabled for this single commit using command-local `core.hooksPath=/dev/null`; no configuration was persisted. Only the authorized main-tree progress log and this handback were written. Worktree and branch are preserved for integration.
