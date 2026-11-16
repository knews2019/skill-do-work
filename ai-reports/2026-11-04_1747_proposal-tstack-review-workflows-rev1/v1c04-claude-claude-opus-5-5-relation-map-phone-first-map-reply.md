Assigned file: [index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html](index-v1c04-claude-claude-opus-5-5-relation-map-phone-first-map.html)

### Clarity changes
- The report is now one zoomable relation map (`#rev-1`): 23 boxes in three columns (what is known, what is proposed, what tests or limits it) and 28 relations of four kinds (supports, chooses, limits, answers), with a key. The first view shows ten main boxes and the words of all 12 relations between them.
- On a phone the map starts in the top half of the first screen: the intro paragraph moves below the map, the note on how I read the lines sits beside the key, and the first view puts the main boxes at the frame's top.
- A click on a box opens its sources and relations in the detail panel and shows its related boxes. On a phone, a click keeps every title at about 12 px, and each relation in the panel has a Show button that brings its box into view.
- Relation words keep clear of arrowheads and of the frame edge, because line ends count as taken space.
- Every box is also given as text (`#rev-3`), and every relation as one table row (`#rev-4`). Without JavaScript the map shows at full size with the words on its lines.

### Corrections
- The report says REQ-710 "remains pending". It is completed and was released as 0.310.0 on 2026-11-05 (`do-work/archive/UR-163/REQ-710-add-optional-bounded-review-modes.md`, `CHANGELOG.md`).
- What shipped differs from the proposal: at most five total rounds or two stalled rounds, and an unresolved loop ends as `completed-with-issues`, not as a hold, which answers Q1 (`skills/do-work/actions/ai-evaluation-loop-reference.md:60` and `:69`).
- The report's line links pin `ccea405e`; the same rules now sit at `skills/do-work/actions/work.md:145-147`, `skills/do-work/actions/work-reference.md:365` and `skills/do-work-toolbox/actions/validate-feedback.md:95`.

### Open questions
- Q2 (which environments must support provider selection first) has no later answer in the repository.
- The upstream tstack files were not read again, because this pass may not use external services; their summaries come from the report.
- No one checked the map's lines against the report except me.

Variation slug: phone first map
