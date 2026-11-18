---
id: REQ-740
title: 'Browse report types and save new samples'
status: cancelled
created_at: 2026-11-18T15:04:36Z
user_request: UR-182
domain: frontend
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
maintenance: false
batch: report-sample-gallery
related: [REQ-741, REQ-742]
completed_at: 2026-11-18T15:28:06Z
---
# Browse report types and save new samples

## What
Build an interactive gallery with inspectable sample reports for every type the detailed AI report workflow supports. Let a reader add an example discovered elsewhere and find it again in the same gallery.

## Why
The reader should be able to understand the main finding in most reports.

## Detailed Requirements
Show what each type helps a reader understand. Include a way to inspect an actual example. Preserve enough of each addition to identify the source and the useful presentation idea on a future run. Cover the existing report kinds as well as the ways they can present evidence.

## Constraints
The screenshots illustrate reading patterns. They supply design references, and each generated report supplies its own facts.

## Red-Green Proof
**RED prompt/case:** Open the gallery, inspect its supported types, save the supplied inflation comparison, and reopen the gallery.
**Why RED now:** The existing style demonstration and the report catalog offer examples, but neither provides the requested ongoing collection and add flow.
**GREEN when:** The types have inspectable samples and the newly saved comparison is present and usable after reopening.
**Validation:** User confirmed the end-to-end inflation example during capture.

## Assets
- `do-work/user-requests/UR-182/assets/REQ-740-screenshot-1-daily-activity.png` — Daybook shows time across the day with a combined band chart and a second view with six aligned activity rows. Peak labels and field notes mark specific moments.
- `do-work/user-requests/UR-182/assets/REQ-740-screenshot-2-nobel-decades.png` — Nobel prizes by decade. Stacked columns lead to decade details, small charts compare categories, and hatching identifies a partial decade.
- `do-work/user-requests/UR-182/assets/REQ-740-screenshot-3-inflation-comparison.png` — “Inflation fell. Did prices?” pairs the annual inflation rate with the price level over the same months. A falling rate and rising price level answer different questions.
- `do-work/user-requests/UR-182/assets/REQ-740-screenshot-4-inflation-evidence.png` — The companion view labels the price-level change on the chart and shows a verdict with expandable evidence about the measure, units, dates and arithmetic.

## Builder Guidance
Use the existing report gallery demonstration, report action and second-opinion workflow as starting points. The choice of storage, matching technique and page design belongs to the builder.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8761 tokens) covers report routing and alternate writers. The index marks it partially slugged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read relevant rules and record the approach.
- [ ] **[APPLY]:** Implement the agreed behavior.
- [ ] **[UNIFY]:** Verify the diff and the acceptance example.

## Full Context
See `do-work/user-requests/UR-182/input.md` for the original request and the confirmed example.

## Cancelled

- **When:** 2026-11-18T15:28:06Z
- **Why:** Replaced on 2026-11-18 by a two-request rewrite the user chose: one request for the report gallery and one for the report writer that uses it.
- **Decided by:** user, via `do-work abandon`
