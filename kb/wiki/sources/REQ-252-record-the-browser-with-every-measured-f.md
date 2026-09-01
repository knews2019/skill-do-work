---
title: "Lessons from REQ-252: Record the browser with every measured-face number in the Durations tests"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-01/REQ-252-record-the-browser-with-every-measured-f.md]
related:
  - page: concept-duration-estimation-and-breaks
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-252: Record the browser with every measured-face number in the Durations tests

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

The Durations suite now holds several constants that are a **browser's answer**, not arithmetic — and two independent measurements of the same face, taken on different Chromium builds, disagreed by enough to matter. Every such constant needs the build it came from recorded next to it, and the numbers that are now known to be per-browser need saying so.

## Solution summary

Every browser-measured constant in the Durations suite (7 `durationsMeasured*` constants across `durations_test.go` and `generate_test.go`) now names the Chromium build its number came from, written exactly as the source REQs recorded it, each with a fresh cross-check on this environment's Chromium 141.0.7390.37. A new `TestDurationsMeasuredConstantsNameTheirChromiumBuild` walks every Go file in the package (go/parser, vacuity-guarded) and fails any measured constant whose doc comment names no build — observed RED on all 7 before the edits. The Panel B clearance budget prose states it is per-browser and carries all three measurements (1.364 REQ-241 / 0.185 Chromium 146 / 1.111 fresh, post-REQ-248 geometry) plus the corrected model-space figure (0.49 units, pitch 14 breaks it). `durations.go` documents the hyphen-vs-U+2212 divergence (5.24 units on this build vs the REQ's recorded 1.73 — the delta itself is per-browser) and states the remainder reserve's over-shoot as deliberate with measured numbers; REQ-260's gofmt nit fixed in passing. **No constant's value changed** — two current-build measurements exceed their constants and that raise is a Discovered Task, per the REQ's no-fold rule.

## What worked

**What worked:** Enforcing a documentation convention with a real AST-walking test (vacuity-guarded, mutation-falsifiable) instead of trusting comments to stay honest — RED on all seven constants proved the gap before any edit. Recording provenance identifiers exactly as the archives state them, inconsistencies included, rather than inventing tidier versions.

**What didn't:** Discovered Tasks that live only in hand-back prose are one integration slip from evaporating — the review caught that neither had a durable artifact. Capture them in the REQ's own section at hand-back time.

**Worth knowing:** The 11px mark-label box measures LARGER than its recorded constant on Chromium 141 (12.9631 vs 12.84) — the raise is REQ-265; until it lands, the pitch-13 floor clears reality by 0.037 units, not the 0.16 the constant implies. Even the hyphen-vs-U+2212 delta is per-browser (1.73 recorded vs 5.24 measured here). The `durationsMeasured` prefix is a convention the test enforces comments for; a smuggled number under another name is review's job.

## Back-reference

See `do-work/archive/UR-051/REQ-252-record-the-browser-with-every-measured-face-number.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c752529`.
