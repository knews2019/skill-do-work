---
id: UR-053
title: Durations lane text overlap — audit, mockup alternatives, then Alternative 2
created_at: 2026-08-18T10:28:46Z
requests: [REQ-231]
word_count: 60
---

# Durations Lane Text Overlap — Audit, Mockup Alternatives, Then Alternative 2

## Full Verbatim Input

-request [Image #1] <- there is a lot of text overlapping, do an audit and propose sensible UIUX improvements, also generate several alternatives with mockups using the ai-report before implementing

*(after the audit + proposal report was delivered)*

ok, go with Alternative 2 — top-N extremes in the text band (recommended)

*(mid-implementation)*

ok, you started working, that's fine, but I also want to capture the intent via do-work capture request

## Screenshot Description

The attached screenshot shows the queue-kanban board's Durations view for the project `glw-game-find-the-difference`, generated 2026-08-18 09:54 UTC — 577 archived REQs with both stamps across 43 active days. Panel A's overflow lane (the `60+` strip) holds 51 samples. One label at the far left (`REQ-407 14h 15m`) is legible. In the dense right half, the first label row (`REQ-876 3h 48m`, a partly obliterated `REQ-88…`, `REQ-1177 16h 9m`) is overprinted by the overflow dots themselves; the second row (`REQ-881 2h 17m`, `REQ-897 1h 32m`, `+46 more over 60 min`) is clean. Panels B (median bars, `78 min` annotation, `45+` tick) and C (green count bars) are legible. This is the post-REQ-226 state: labels no longer collide with each other; marks overprint the labels.

## Clarifications Answered During Capture

1. **Which design should the fix implement?** → The maintainer chose **Alternative 2 — top-N extremes in the text band** from the five-option proposal report at `ai-reports/2026-08-18_1000_durations-lane-legibility-alternatives/` (commit 0f8e349): marks keep the lane's top strip, both label rows move below a divider with leader ticks (the captured REQ-231 spacing fix), **and** label selection changes from left-to-right first-fit to top-N-by-magnitude, so every drawn label is one of the band's longest spans.
2. **New REQ or existing?** → Existing: queued REQ-231 (`addendum_to: REQ-226`) already captured the spacing half of this defect. This UR is recorded as an addendum to it rather than as a duplicate REQ; the Alternative-2 selection rule is the scope it adds.

## Assets

- `assets/REQ-231-screenshot-1-durations-mark-over-label-overprint.png` — the submitted Durations view screenshot described above. Its ids and counts belong to another repository; it is evidence that the defect scales with sample count, not data to reproduce locally.
