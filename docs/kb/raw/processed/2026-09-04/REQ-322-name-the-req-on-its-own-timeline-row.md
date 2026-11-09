---
source_type: req_lesson
req_id: REQ-322
req_path: do-work/archive/UR-065/REQ-322-name-the-req-on-its-own-timeline-row.md
date: 2026-08-23
domain: frontend
module: _dev/primes
tags: [frontend, name, timeline]
---

# Lessons from REQ-322: Name the REQ on its own timeline row

## What the REQ was about

Show each REQ's title in the row's label column, and put its detail in a tooltip at the
pointer instead of only in a readout at the foot of the panel.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

## What worked

- The monospace property. One measured advance for the whole render, and a
- guard that refuses a face it cannot describe — the review independently confirmed the 0.5px
- tolerance is right in both directions (Georgia and Arial miss it by three orders of
- magnitude). Building the refusal first meant the Unicode fix was a change of unit, not a

## What didn't work

- **I wrote a false claim into a decision record.** D-01 said a test asserted the column-width
- floor. No test did — both files restated `172` instead of reading `TIMELINE_LABEL_WIDTH`, so
- the number the whole REQ was about was pinned by nothing and reverting it passed everything.
- The prime records this exact class as REQ-265, one batch earlier. **A constant a decision
- turns on has to be read by the test, never restated beside it** — and a claim that a test
- exists is checkable in ten seconds, which is ten seconds I did not spend.
- **A verified assumption is only verified for what it sampled.** The guard proves the face is
- monospace using `i` and `M`, and I then applied its answer to arbitrary Unicode. 中 is 10px
- on the same 6.02px face. The guard was not wrong; its *scope* was narrower than the use.
- **Geometry is not legibility.** My render evidence measured widest-label pixels and overflow
- counts, both zero-defect, and quoted three labels that all happened to be non-review REQs.
- The sixteen newest REQs on this very board read `[impact-user-visib…` and named nothing —
- the exact failure the REQ exists to remove, on the first screen, invisible to every number I

## Worth knowing

- the label cell is the face's *Latin* advance. Anything that puts new text
- in that column — a different script, an emoji, a longer id — is measured in cells by
- `timelineLabelCellCount`, not characters, and non-ASCII deliberately over-counts so the error
- falls on the side of cutting early. And the `[impact-token] ` title convention and a
- nineteen-cell budget cannot both have the front of the string: the label strips the tag, the
- tooltip and table keep it.

## Back-reference

See `do-work/archive/UR-065/REQ-322-name-the-req-on-its-own-timeline-row.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1c42897`.
