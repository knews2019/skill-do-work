---
id: REQ-178
title: Build the audit-metrics tool for mechanical audit measurement
status: pending
created_at: 2026-08-13T22:35:10Z
user_request: UR-040
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
related: [REQ-176, REQ-177]
batch: maintainability-audit
write_set: [skills/do-work-toolbox/tools/audit-metrics/]
maintenance: false
---

# Build the Audit-Metrics Tool for Mechanical Audit Measurement

## What

A small Go tool, `skills/do-work-toolbox/tools/audit-metrics/`, that produces the maintainability audit's deterministic numbers mechanically — inventory, distributions, band flags, churn — so the action pastes tool output instead of prescribing fragile find/wc/awk pipelines to an LLM. Script what can be scripted; judgment stays in prose.

## Why (if provided)

User (verbatim): "since you also have a go tool, consider building tools for the audit, that will also output mechanicanically some flagged folders, files, etc... for the MVP it does not have too be too complex, but basically whatever we can script would be good to have it as script not as LLM call, becuase those are cheaper and more robust". This is CLAUDE.md's "Programs beat prose for anything mechanical" applied to the audit itself.

## Detailed Requirements

MVP scope — what wc/find/git can answer robustly; CCN and duplication stay with the external tools (lizard/jscpd) and their NOT-MEASURED path:

1. **Inventory:** tracked-file counts and line/word totals by extension, honoring an exclude list (flag or config; defaults per REQ-176 requirement 14).
2. **Distributions:** per metric — file lines, file words, folder file-counts — median / p90 / p95 / max, plus top-N largest offenders.
3. **Band flags:** apply WATCH/FLAG thresholds passed as flags; output the flagged folders and files mechanically (path, value, band). Bands are inputs, never hardcoded — calibration happens in the action's conversation.
4. **Churn:** `git log --since=<window> --name-only` aggregation with: shallow-clone detection (report it, never silently truncate), exclude patterns for release-ceremony files, and current-path filtering so dead pre-rename paths don't rank. Top-N output.
5. **Hotspot join:** churn × size (size as the MVP complexity proxy), top-N.
6. **Output:** markdown tables suitable for pasting directly into the audit report; machine-readable TSV behind a flag if trivially cheap, else skip (YAGNI).
7. **Pattern match with queue-kanban:** vendored source, built on demand (`go build` then run), invoked by the action as an accelerator with the manual-fallback contract — if `go` is absent or the build fails, the action falls back to the manual commands in its reference file; the tool is never a dependency.
8. Go tests pin the contract: distribution math on a fixed fixture, band flagging edges (value == threshold is not flagged; > is), shallow-detection reporting, exclude-list honoring. Focused lock-in tests, not smoke slop.

## Constraints

- Read-only: the tool writes nothing; it prints. (The action owns the `do-work/audits/` report write.)
- Naming for reach per coding-guardrails § 5 — `audit-metrics` and its flags need findable, two-word names; single-word subcommands are exempt by design.

## Dependencies

REQ-176's action consumes this tool (its Phase 0/1 steps invoke it with fallback); build this first or in the same wave.

## Builder Guidance

Certainty: Firm on the MVP boundary (items 1–5, nothing more); Exploratory on CLI shape (single run vs subcommands — pick what queue-kanban's conventions suggest). Keep it deliberately small; complexity added here is complexity the audit pays forever.

## Red-Green Proof
**RED prompt/case:** `ls skills/do-work-toolbox/tools/audit-metrics/` fails; a failing Go test for the distribution math can be written first in the new package.
**Why RED now:** No mechanical measurement exists; the draft spec prescribes hand-run shell pipelines for every number.
**GREEN when:** `go test ./...` passes in the tool directory; running the built tool against this repo prints inventory, distributions, band flags (given thresholds), and churn tables matching hand-computed spot checks.
**Validation:** User requested mid-capture; MVP boundary inferred from their words ("does not have to be too complex").

## Full Context

See `do-work/user-requests/UR-040/input.md` for complete verbatim input.

---
*Source: UR-040 — user follow-up during capture*
