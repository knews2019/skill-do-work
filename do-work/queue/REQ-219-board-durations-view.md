---
id: REQ-219
title: Durations view on the Kanban board
status: pending
created_at: 2026-08-17T17:17:22Z
user_request: UR-050
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
---

# Durations View on the Kanban Board

## What

Add a dedicated **Durations** view to the Go-served Kanban board showing how long archived REQs actually take and how often they run: three panels on one shared calendar axis — duration per REQ, median minutes per active day, and REQs completed per day. Data comes from the archive scan `queue-kanban` already performs, not from a new file.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

From the user, verbatim: *"that way I can identify if there is an overtime response time degradation, outliers, how often are these reqs executed."* Those three questions are the acceptance frame — a panel that does not answer one of them is not earning its place.

## Context

The chart already exists and is render-judged in light and dark at `ai-reports/2026-08-17_1401_UR-048-estimator-calibration-and-anomaly-surfacing/index.html` (section 8, commit `6e79932`). **That is the reference implementation, not a suggestion to re-derive.** It resolved several decisions the hard way; port them rather than rediscovering them:

- The three panels are separate on a shared x-axis precisely because a single overlaid plot would need two y-scales. Never a dual axis.
- Panel A needs an overflow lane above a scale break: REQ-064's 10h 55m span otherwise squashes the 90% of samples that sit under 35 minutes.
- Route A/B/C are **ordered difficulty tiers**, so they take a single-hue ordinal ramp, not categorical hues.
- An in-cloud trend line was tried first and discarded — a rolling window over completion order straddles session boundaries and produces spikes that read as drift, and a continuous stroke across a 27-day idle stretch asserts a climb that never happened.

## Detailed Requirements

**Data source — the archive scan, at build time.**

- Derive every sample from archived REQ frontmatter the board already parses: `id`, `route`, `claimed_at`, `completed_at`. No new parser fields, no new file reads.
- Include a REQ when its status normalizes to terminal success and **both** stamps parse. 195 samples qualify today, spanning 2026-05-29 to 2026-08-17.
- Record the wall span **raw and signed**. A negative span is data, not an error to swallow — it is the anomaly class REQ-213 taught this board to surface. The archive carries none today (REQ-091's reversed stamp was owner-repaired in commit `24915a8`), so this behaviour must be pinned by a fixture rather than by real data.
- Do **not** read `do-work/calibration-log.tsv`. It holds 8 rows, only grows going forward, and is append-only so a repaired stamp never propagates. Explicitly rejected during capture.

**Panel A — duration per REQ.**

- One mark per REQ, x by completion instant, y by wall minutes.
- Y ceiling 60 minutes with an **overflow lane above a scale break** for longer spans; today that is REQ-064 (10h 55m), REQ-043 (2h 29m), REQ-171 (1h 21m).
- Negative spans render **below the zero line** in the status-critical colour, labelled — never rounded up to zero. None exist in the archive today; keep the branch and cover it with a fixture, because the whole point of REQ-213/214 is that one can reappear.
- Direct-label only the marks that matter individually (the overflow marks and any negative span). Never a value on every point.
- Colour by route using the ordinal ramp below; unrouted legacy REQs take muted ink, not a step on the ramp.

**Panel B — median minutes per active day.**

- One column per active day. **Apply the calibration's documented read-time outlier rule** — spans over four hours are an assumed pause, negative spans are broken stamps, both excluded (`skills/do-work/actions/estimate-reference.md` → Calibration).
- This is the panel that answers "is it degrading". Applying the rule here is what stops one paused session inventing a five-hour day: 2026-07-31 has two samples, and including REQ-064's 655 minutes yields a median of 328.9 instead of the true 2.5.
- Panel A still shows those excluded spans, raw and labelled. State the split where the reader can see it.

**Panel C — REQs completed per day.**

- One column per active day, counting **all** samples including rule-excluded ones (it is a count, not a duration).
- 23 active days out of 81 calendar days today; peak 25 on 2026-08-15.

**Shared axis and chrome.**

- One linear calendar x-axis across all three panels. **Linear real time, not ordinal-by-active-day** — the idle gaps are the cadence finding, and compressing them destroys the answer to "how often are these executed".
- Drop a month tick that would collide with the first or last label.
- Gridlines and axes: solid hairlines, one step off the surface. Never dashed.

**Colour — validated, do not eyeball.**

- Single-hue blue ordinal ramp, light→dark by route weight, validated with the dataviz validator's `--ordinal` mode against the board's own surfaces.
- The invariant is **"heavier route sits further from the surface"**: on a light surface that means C darkest; on a dark surface it means **C brightest**. Reusing the light steps in dark mode was tried and failed — the darkest step measured 2.6:1 and 8px dots simply disappeared. Re-step per mode.
- Text wears text tokens, never a series colour.

**Interaction and accessibility.**

- Hover readout on all three panels via a nearest-mark layer, not pinpoint hit targets on 8px dots.
- A table view listing every sample — REQ id, route, completion stamp, duration, and a note on rows excluded from the day medians — so no value is reachable only by hovering.

## Constraints

- **Read-only. This view adds no write surface.** The board has exactly three (CLAUDE.md § Kanban Board Write Surfaces) and that sentence must not need amending in this REQ's commit. If the implementation finds itself wanting a fourth, stop and surface it rather than adding one.
- No new CDN or npm dependency. The board's web layer is vanilla JS and stays that way; the chart is hand-authored SVG.
- Board versioning is folded into the skill: normal `CHANGELOG.md` entry plus suite version bump, per `_dev/primes/prime-kanban-board.md`.
- The Kanban columns and the existing Testing view are untouched.
- Aggregation cost rides on a parse that already happens; do not introduce a second pass over the archive.

## Builder Guidance

**Certainty: firm** on the data source, the read-only constraint, the three panels, the outlier-rule split between panels A and B, and the per-mode ramp re-stepping — all four were decided with the user or learned from a failed render, and re-litigating them wastes the work already done.

**Latitude** on: the Go aggregation's file and function layout, how the view is wired into the existing view switcher, exact SVG geometry and margins, tick counts, and tooltip copy. The reference implementation is a design contract, not a code contract — it is HTML in a report, and the board is Go plus vanilla JS. Port the decisions, write idiomatic board code.

Read the reference implementation's section 8 before planning. It carries inline comments explaining *why* each of the firm decisions is what it is.

## Red-Green Proof

**RED prompt/case:** Build and open the board today (`queue-kanban` → the generated board). There is no Durations view, and the generated payload carries no per-REQ duration and no per-day aggregate. Nothing anywhere in the board answers "are REQs taking longer than they used to", "which ones are outliers", or "how often does this actually run" — the only duration-adjacent signal is the completion-anomaly strip, which fires on broken bookkeeping, not on slowness.

**Why RED now:** The board renders queue state and completion anomalies. Duration has never been aggregated or surfaced, so the whole question has to be answered by hand-mining the archive — which is exactly what UR-048's calibration had to do.

**GREEN when:**
1. A Durations view exists beside the board and Testing views, rendering all three panels over every qualifying archived REQ (195 today) on one shared calendar axis.
2. A Go test pins the aggregation against known values: with the read-time rule applied, `2026-07-31` yields a day median of **2.5 min** (REQ-064's 655-minute span excluded) and `2026-08-15` yields **19.6 min over 25 REQs**; without the rule, 2026-07-31 would be 328.9 — assert the rule is actually applied, not just present.
3. A fixture with a reversed span (completed before claimed) keeps its raw negative value through aggregation and is excluded from day medians rather than clamped to zero.
4. `go vet ./...` and `go test -count=1 ./...` green in the module; `queue-kanban verify` behaviour unchanged.
5. The three-write-surface sentence in `CLAUDE.md` is untouched by this REQ's commit.

**Validation:** User confirmed — three capture questions answered (archive scan as the data source; a new dedicated view; all three panels), each the recommended option.

## Full Context

See `do-work/user-requests/UR-050/input.md` for complete verbatim input.

---
*Source: UR-050 — port the report's duration chart to the board*
