# Maintainability Audit

A measured, repeatable code-health baseline with deltas across runs. Metrics (size, complexity, churn, duplication) run under WATCH/FLAG bands calibrated to *your* repo — not generic thresholds — judgment reads only the hotspots, and every run writes a persistent report to `do-work/audits/` so the next run can prove things got better. Read-only outside that folder: the audit proposes, it never fixes.

## The loop

The audit is one step in a cycle you drive. One pass:

1. **Run it** — `do-work-toolbox maintainability-audit` (or `audit codebase`).
2. **Calibrate** — before measuring anything, the audit grounds itself in the repo and presents one bundled, editable proposal at the calibration gate: tool installs (user-local only, exact commands), WATCH/FLAG bands proposed against the repo's own measured distributions, and scope/excludes — plus at most three focused questions ("which area hurts most when you touch it?"). Edit any line, approve, and measurement starts. Repeat runs skip the full gate and ask exactly one question: "reuse last run's calibration, or recalibrate?"
3. **Read the report** — `do-work/audits/audit-<date>.md`: a metrics summary with a delta table against the prior run, at most a dozen root-cause finding classes, plus Pre-empted (candidates a documented decision or waiver already covers) and NOT-MEASURED (metrics whose tool was declined or absent — always stated, never estimated).
4. **Triage** — paste the report's `## Findings` section into `do-work-toolbox validate-feedback`. Findings are refutable claims with reproduce commands, and the validator verifies each one against the real code before anything reaches the queue.
5. **Capture and build** — accept findings through the validator's capture handoff; park anything you want to think about with `do-work-toolbox note`. Then `do-work run` builds the fixes.
6. **Re-audit** — after fixes land, run the audit again. The delta table must move; a flat table across runs means the loop is decorative.

## Lock-in limits

Each finding class proposes at most one lock-in limit: a single number or zero-hit assertion pinned at the current worst observed value ("no folder over today's max of 34"). Green on day one, it blocks regression immediately — and it only ever tightens as fixes lower the worst case. The audit only proposes; accepted limits flow through the validator and land as ordinary REQs in your own test suite or CI.

## Living with findings

A class you've decided to live with goes in `do-work/audits/waivers.md` — one line per waived class with the reason — not into another round of triage. The audit reads the waivers file every run and never re-flags a waived class.

## When a finding fights a decision

Findings that collide with a documented decision or recorded lesson land in Pre-empted, or explicitly challenge the decision — never emitted as if it didn't exist. If a push-back cites a decision you no longer agree with, change the decision doc — not the code. The audit respects whatever the record says, so keep the record honest.

## Usage

```
do-work-toolbox maintainability-audit                # full-repo audit
do-work-toolbox maintainability-audit src/ lib/      # propose scope narrowed to these dirs
do-work-toolbox maintainability-audit recalibrate    # force full calibration on a repeat run
do-work-toolbox audit codebase                       # same action
do-work-toolbox audit maintainability
```
