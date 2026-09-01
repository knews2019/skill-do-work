---
title: "Lessons from REQ-253: Decide the Timestamp rule's two uncovered stamp shapes"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-253-decide-the-timestamp-rule-s-two-uncovere.md]
related:
  - page: concept-timestamp-and-metadata-governance
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-253: Decide the Timestamp rule's two uncovered stamp shapes

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

Two clock-write shapes in shipped action files are governed by neither paragraph of the Timestamp rule. Both need a decision before they can be swept, and neither decision is the pipeline's to make.

## Solution summary

Both user-answered decisions implemented; the Timestamp rule did not grow a third paragraph. (1) The rule's date-only paragraph now names ui-review's report-header `**Date**` as a deliberately-UTC date-only consumer; the site cites the rule above the report fence and its placeholder reads `[today's UTC date]`. (2) One sentence in the same paragraph declares `## HH:MM UTC` time-of-day headings out of the rule's scope; the canonical statement lives in `memory-reference.md` § Daily-Log Entry Conventions, and all five write/template sites carry a uniform greppable marker ("time-of-day label, outside the Timestamp rule's scope"). The class was closed by grep over the corpus, not by the REQ's stale line list — the REQ's `memory.md:140` was a read-only checklist line; the real write sites (memory.md 50/84, memory-reference.md 46/93/135) are the marked set.

## What worked

**What worked:** Closing the class by condition ("writes or templates the heading") instead of the REQ's line list — the listed line turned out to be a read site, and the real write sites were elsewhere. Uniqueness-asserted scripted edits made six small text changes safe in one pass.

**Worth knowing:** The `## HH:MM UTC` shape is defined once in memory-reference.md § Daily-Log Entry Conventions; the site markers are pointers to it. The date-only paragraph now carries a tripped tripwire — it says "revisit if a second consumer appears" and ui-review is that second consumer (REQ-261 asks the question).

## Back-reference

See `do-work/archive/UR-055/REQ-253-decide-the-timestamp-rule-boundary-cases.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0d8d629`.
