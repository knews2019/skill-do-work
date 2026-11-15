Assigned file: [index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html](index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html)

### Clarity changes
- The report is now one zoomable relation map (`#rev-1`): 23 boxes in three columns (what is known, what is proposed, what tests or limits it) and 28 relations of four kinds (supports, chooses, limits, answers), with a key. The first view shows ten main boxes and the words of all 12 relations between them.
- A click on a box opens its sources and relations in a side panel, shows its related boxes, and dims the rest a little. "Main boxes" and "Show everything" switch between the first view and all boxes.
- Every box is also given as text (`#rev-3`), and every relation as one table row (`#rev-4`). Without JavaScript the map shows at full size with the words on its lines.
- The page says that the lines are my reading of the report. Boxes without a report code get page codes (D, J, P1 to P5, L1 to L3, REQ-710).

### Corrections
- The report says REQ-710 "remains pending". It is completed and was released as 0.310.0 on 2026-11-05 (`do-work/archive/UR-163/REQ-710-add-optional-bounded-review-modes.md`, `CHANGELOG.md`).
- What shipped differs from the proposal: at most five total rounds or two stalled rounds, and an unresolved loop ends as `completed-with-issues`, not as a hold, which answers Q1 (`skills/do-work/actions/ai-evaluation-loop-reference.md:60` and `:69`).
- The report's line links pin `ccea405e`; the same rules now sit at `skills/do-work/actions/work.md:145-147`, `skills/do-work/actions/work-reference.md:365` and `skills/do-work-toolbox/actions/validate-feedback.md:95`.

### Open questions
- Q2 (which environments must support provider selection first) has no later answer in the repository.
- The upstream tstack files were not read again, because this pass may not use external services; their summaries come from the report.
- No one checked the map's lines against the report except me.

Variation slug: review loop map
