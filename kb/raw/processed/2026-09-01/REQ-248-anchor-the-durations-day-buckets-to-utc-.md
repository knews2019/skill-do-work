---
source_type: req_lesson
req_id: REQ-248
req_path: do-work/archive/UR-051/REQ-248-anchor-durations-day-buckets-to-utc-midnight.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, anchor, durations, buckets, midnight]
---

# Lessons from REQ-248: Anchor the Durations day buckets to UTC midnight so Panel B stays on canvas

## What the REQ was about

Panel B's bars are placed with `xOfEpoch`, which maps each day bucket's **midnight**, while `timeStart` is the **first completion instant**. The two disagree by however far into its first day the earliest sample falls, so the leftmost bar is drawn to the left of the plot area — and on a board with one or two active days the disagreement dominates the whole span and the panel renders off-canvas entirely.

## Solution summary

The Durations axis domain is anchored to whole UTC days — first completion floored to its UTC midnight, ending at the midnight after the last — with day buckets (Panel B/C bars, the slowest-day annotation, the hover day-nearest rule) centred on each day's noon, and the outermost bars clamped inside the plot for the >280-active-day case where the 4-unit minimum bar width exceeds the day slot. The Go label planner (`durationLabelTimeRange`) floors/ceils identically, so renderer and planner share one domain; a new mark-position agreement assertion fails if either side's domain ever drifts alone. Axis end label and aria-label keep naming the last *active* day. Verified live in Chromium 141.0.7390.37 headless at 1/2/14 active days: zero bars or annotations outside the plot area [54, 1182] (before: bar x=−5184.4 at one day).

## What worked

**What worked:** The mark-position agreement assertion is the class-closure this board's geometry work had been missing — it fails in *both* drift directions (JS-only revert, Go-only revert), so renderer and planner cannot silently become two definitions of the domain again. Sweeping day counts to 400 caught the hole behind the instance (the 4-unit minimum bar width overhanging past ~280 days), which the floor/ceil fix alone did not close. Driving the full real `renderDurationsView` over a DOM stub found what sliced-function tests hide.

**What didn't:** The captured write set excluded `durations.go` on a collision theory that was empty (REQ-252 is gated behind this REQ), and criterion 4 was unsatisfiable without it — the D-01 escalation cost a scope negotiation mid-build. When a REQ's own guidance says "check the other readers of this domain", the files encoding those readers belong in the write set from the start.

**Worth knowing:** Day buckets centre on noon, not midnight — a bar at its floored midnight straddles the previous slot. The axis end label deliberately names the last *active* day, not the domain's exclusive end (D-04, currently unpinned by any test). `durationLabelPlotX`'s zero-width-domain guard is now unreachable and its comment stale (review Minor 1).

## Back-reference

See `do-work/archive/UR-051/REQ-248-anchor-durations-day-buckets-to-utc-midnight.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1cb897f`.
