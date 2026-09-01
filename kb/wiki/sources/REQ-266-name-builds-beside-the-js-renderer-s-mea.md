---
title: "Lessons from REQ-266: Name builds beside the JS renderer's measured face numbers"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-266-name-builds-beside-the-js-renderer-s-mea.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-266: Name builds beside the JS renderer's measured face numbers

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

`web/board-durations.js` presents measured face numbers (12.83 / 10.43 / 2.41 at `DURATIONS_LABEL_ROW_HEIGHT` and `DURATIONS_LABEL_TEXT_ASCENT`) as current fact with no build named — the same provenance gap REQ-252 closed in the Go files, on the JS surface its go/parser test cannot reach. Extend the rule: every browser-measured number in the JS comments names its build, and the mechanism that keeps it true (a JS-side check, or a stated review convention) is the builder's call, recorded either way.

## Solution summary

Removed the last undated measured number from the renderer's comments and replaced its stale test citation, then made the rule enforceable on the JS surface with a check rather than a convention.

## What worked

**What worked:** Doing the re-read the REQ asked for instead of trusting its own expected outcome. The REQ said twice that closing as no-longer-applicable was expected, and REQ-292 had just landed — it would have been easy and defensible to close it. One of the three instances had survived, along with a citation to a test deleted an hour earlier in the same run.

**What didn't:** The first mutation test of the new check passed when it should have failed, because the defect was reinserted into a comment block that already named two REQs. That is the check's real anchoring behaviour, not a mistake in the mutation — and finding it took writing a mutation that *should* fail and being surprised. A mutation test that passes is either a working guard or a broken check, and the only way to tell is to write the second mutation.

**Worth knowing:** The useful distinction here generalises past this file. A number in a comment is doing one of two jobs: asserting something true *now* about the environment, or citing evidence for a decision *then*. The first goes stale invisibly and needs dating; the second is already dated by the decision it supports. Demanding provenance for both makes the rule annoying enough to be ignored, which is roughly how the undated numbers got there.

## Back-reference

See `do-work/archive/UR-051/REQ-266-name-builds-beside-the-js-measured-face-numbers.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8fe9eb6`.
