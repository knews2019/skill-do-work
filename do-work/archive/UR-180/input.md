---
id: UR-180
title: 'Phone layout for relation-map pages: readable titles after a click, and the map on the first screen'
created_at: 2026-11-15T12:45:58Z
requests: [REQ-737, REQ-738]
word_count: 25
---
## Summary
The maintainer wants the relation-map pages to work on a phone. Two phone limits were recorded as report-only when REQ-727 (ship a relation-map preset for ai-report second-opinion) was released as 0.324.0 on 2026-11-15. D-16: on a 390 px phone, one click zooms the map out to about 0.2 to 0.34, so box titles become about 4 to 7.5 px. D-17: the launcher prints the whole preset text above the page, about 514 px tall on desktop and about 1,302 px on a phone, so the map starts on the third phone screen.

The two limits live in two different places, so the capture makes two REQs. The preset text owns what a generated map does after a click. The launcher's Go code owns the direction block above every directed page, for every preset.

Context checked on disk at capture time (coordinator context, not the user's words):
- REQ-727 is archived at `do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md` (D-16 and D-17 at lines 197-198, findings F1 and F2 at lines 297-298). Release commit `d50775cd`, CHANGELOG `## 0.324.0` dated 2026-11-15, whose "Known limits" line names both limits.
- Integration report `do-work/runs/work-2026-11-15-012821/REQ-727-integration-report.md` line 8 and hand-back `do-work/runs/work-2026-11-15-012821/REQ-727-handback.md` lines 70-96 record the numbers. The review `do-work/runs/work-2026-11-15-012821/REQ-727-review.md` line 64 measured the map frame top at y = 853 (1440 px wide) and y = 2,069 (390 px wide).
- The direction block is written in `secondOpinionDecorate` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (line 1634 prints the `Direction:` paragraph).
- The shipped preset is `skills/do-work-toolbox/docs/second-opinion-presets/relation-map.md` (3,057 bytes, 20 checks).
- The test page is `ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html`.

## Extracted Requests

| REQ | Request |
| --- | --- |
| REQ-737 | Addendum to REQ-727: the relation-map preset gives a phone layout in which box titles stay readable after a click (preset text in `relation-map.md`). |
| REQ-738 | The launcher's direction block no longer pushes the page content down: on a 390 x 844 phone the relation map's top edge is inside the first screen. This block is on the pages of every preset, not only the map. |

## Batch Constraints
- The two REQs are independent. Neither depends on the other: they change different files, and each proves its own phone number. The user named no order.
- Both change shipped files under `skills/`, so each is a release (`_dev/primes/prime-releases.md`).
- Phone numbers are measured at 390 x 844 px, the size REQ-727's review used. Desktop is measured at 1440 x 1000 px.
- Numeric targets are approximate unless a REQ says strict.

## Full Verbatim Input
> ```
> relation-map preset needs a phone layout: readable titles after a click, and the launcher's preset text must not push the map below the first screen
> ```
