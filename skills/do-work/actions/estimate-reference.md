# P50 Estimation — Reference

> **Companion file to `actions/work.md` (ensure-estimate step) and `actions/verify-requests.md` (post-repair recalculation).** Holds the signal-extraction guide, the `estimate:` frontmatter block template, the confidence rubric, and the presentation formats. Load it only at an estimation moment — never during ordinary claiming, building, or review. If it is already in context this session, reuse it. The `effort_estimate: effort-mechanical` short-circuit below deliberately skips this file entirely.

## What the Estimate Means

`p50_active_minutes` is the median estimate of **active agent wall-clock minutes** — roughly a **50% chance of completing within the estimated minutes** while do-work and its agents are actually working: planning, exploration, implementation, tests and builds, independent review, and ordinary remediation.

It **excludes**: time waiting for user input or approval, paused/suspended/stopped sessions, overnight or user-controlled gaps, queue wait time, and calendar completion dates. It is an **informational forecast, never a deadline or execution budget** — estimation must never block execution or require user clarification. If estimation fails for any reason, note it and proceed without an estimate.

**No P80 or other percentile fields exist, by design.** Do not add them.

## The `estimate:` Frontmatter Block

Optional and backwards-compatible — a REQ without it is fully valid, and every reader treats the field as absent-is-fine. Written once, then **frozen when execution begins**: never rewrite it with knowledge gained during implementation.

```yaml
estimate:
  p50_active_minutes: 75        # multiple of 5, never below 5
  confidence: medium            # low | medium | high — rubric below
  calculated_at: 2026-08-16T12:00:00Z   # current UTC instant (Timestamp rule, actions/work-reference.md)
  basis:                        # dominant sizing factors, echoed by the estimator
    - Route C
    - 12-file write set
    - browser evidence
```

## Extracting Signals from a REQ

The agent supplies judgment; the shipped script does the arithmetic — same normalized signals, same estimate, always. Read the REQ and map what you find onto the estimator's flags:

| REQ evidence | Flag |
|---|---|
| `route:` frontmatter (post-triage), or your route-equivalent read of the Decision Flow criteria when estimating pre-triage | `--route A\|B\|C` |
| `write_set:` length, or the Scope/requirements' implied file count | `--write-set N` |
| Files/assets/dependencies the REQ says to create | `--new-files N` |
| Distinct runtime subsystems the requirements touch | `--subsystems N` |
| Acceptance-criteria / checklist item count | `--acceptance N` |
| Depth of this REQ's `depends_on` chain (serialization cost) | `--deps-depth N` |
| Browser, responsive, visual, accessibility, or screenshot requirements | `--browser` |
| Persistence, migration, API, or schema changes | `--persistence` |
| Async lifecycle, teardown, race, or retry behavior | `--async-behavior` |
| Performance instrumentation | `--performance` |
| Lint, deploy, asset-integrity, or cross-route regression requirements | `--regression` |
| Full-suite (rather than focused) verification | `--full-suite` |

Independent review and ordinary remediation cost are folded into the route base — no flag. Missing signals never prevent estimation: omit the flag and the estimator still answers.

```bash
<skill-root>/tools/estimate-p50.sh --route C --write-set 12 --browser --persistence --full-suite
```

The output lines map directly onto the block: `p50_active_minutes`, `confidence`, and the `basis` list. Stamp `calculated_at` with the current UTC instant (Timestamp rule, `actions/work-reference.md`).

## Confidence Rubric (deterministic — computed by the script)

- **high** — mechanical-effort short-circuit, or Route A with a raw score ≤ 10 minutes.
- **low** — Route C with a write set ≥ 15 files, ≥ 3 subsystems, or a raw score ≥ 75 minutes: wide scope, wide error bars.
- **medium** — everything else.

## The Mechanical-Effort Short-Circuit

A REQ with `effort_estimate: effort-mechanical` (or obvious Route A indicators) gets the floor estimate without loading this file or extracting signals:

```bash
<skill-root>/tools/estimate-p50.sh --trivial
```

The `--trivial` flag names the estimator's own floor mode, not the schema token — it is the script's interface and does not change with the field's vocabulary. That keeps estimation overhead near zero exactly where the estimate is worth the least. `effort_estimate` itself stays the closed two-value triage chip, `effort-mechanical` | `effort-substantive` (`actions/work-reference.md` → Request File Schema); this short-circuit is the only bridge between the two fields.

## Multi-REQ Totals and Critical Path

For a selected set of more than one REQ, compute both figures with the estimator's graph mode — per-REQ minutes plus `depends_on` edges as `ID:MINUTES[:DEP,...]` triples:

```bash
<skill-root>/tools/estimate-p50.sh critical-path REQ-208:85 REQ-209:60:REQ-208 REQ-210:25:REQ-208
```

`total_estimated_effort_minutes` is the plain sum; `critical_path_minutes` is the longest path through the dependency graph — never the sum of parallel branches. A dependency id outside the set contributes zero (an archived dependency is already done). A **selected member without an estimate enters as a zero-minute vertex** (`REQ-NNN:0:deps`) so the edges through it survive — dropping it would break the transitive chain and understate the critical path. Build each member's dependency list by resolving the legacy `dependencies:` alias exactly as the work action's selection scan does (canonical `depends_on` wins when both are present). Present both figures, clearly labeled:

```
REQ-208  85 min
REQ-209  60 min  depends on REQ-208
REQ-210  25 min  depends on REQ-208

Total estimated effort: 170 active minutes
Estimated critical path: 145 active minutes
```

## Calibration

The scoring table is **calibrated to the archive's measured history** (2026-08-17): of 190 archived REQs carrying both `claimed_at` and `completed_at`, 188 were kept after excluding spans over 4 hours or negative (assumed user pauses / broken stamps). Route bases equal the measured per-route medians — A 4.7 → 5, B 9.2 → 10, C 21.4 → 20 minutes (n = 50/53/45) — which makes a signal-free estimate a true empirical P50; signal weights stretch heavy REQs toward the per-route p80 (A 8.7, B 17.8, C 37.5). Known bias, accepted: the corpus is mostly autonomous runs where wall ≈ active, and any un-filtered short pause inflates actuals slightly — so the table errs conservative.

**The recalibration input is `do-work/calibration-log.tsv`** — appended by the work action's archive step (one line per archived REQ that carried an estimate: `req_id`, `route`, `estimated_p50_minutes`, `wall_minutes`, `completed_at`). The log records raw wall spans; the outlier rule is applied when reading, never when writing.

Two rules keep future calibration honest:

- **`claimed_at` − `completed_at` must never be recorded as `actual_active_minutes`** — it is wall-clock time. It may be *analyzed* as a proxy under the outlier rule above (exclude spans > 4h or negative as assumed pauses), which is exactly how this table was fit and how the next re-fit should read the calibration log.
- **Recalibration changes future estimates only.** Frozen estimates on claimed or archived REQs are never rewritten to match a new table.
