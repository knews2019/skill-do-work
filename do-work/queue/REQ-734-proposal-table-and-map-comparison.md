---
id: REQ-734
title: 'Build the proposal as a supports-and-limits table and compare it with the map'
status: pending
created_at: 2026-11-14T23:21:51Z
user_request: UR-178
domain: frontend
prime_files: []
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-733]
related: [REQ-727, REQ-733]
batch: tinka-follow-ups
write_set: ["ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index-v1c04-*.html", "ai-reports/2026-11-14_1547_proposal-tinka-review-surface/v1c04-*-reply.md", "ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html"]
---
# Build the Proposal as a Supports-and-Limits Table and Compare It With the Map
## What
Build one more page of the same report: a table with one row per option and the columns "supports" and "limits". Then put the table page and the improved map page (REQ-733, rebuild the relation map page) side by side, so the maintainer can see which one tells the story better.
## Why
The maintainer's words: "basically compare the two, so I can see which one is better to tell the story". The assistant had doubted the map: the report is one decision with reasons for and against, and a table would show the same relations with no zooming.
## Context
- The report is `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/index.html`. Its first map page, `index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html`, holds the content as 28 boxes and 37 relations, each also as text and as a table row. The table page reuses that content.
- A new page is written with `ai-report second-opinion --editor-direction '<direction>' ai-reports/2026-11-14_1547_proposal-tinka-review-surface`, with `index.html` as the base. The action is `skills/do-work-toolbox/actions/ai-report.md`, Second-opinion form, Session sequence.
## Detailed Requirements
- One row per option, O1 to O4, with the decision above the table.
- Two main columns: what supports the option and what limits it. Each cell lists the kill attempts, seeds and limits that relate to the option, each with its code, one short line and its source.
- The relations that do not touch an option (for example a limit that weakens a seed) are shown too, in a place the builder chooses, and none of the 37 relations is lost.
- No zooming and no dragging. The page reads top to bottom on a wide screen and on a phone.
- The comparison: take one picture of each page's first screen at the same width, 1440 pixels, and one at phone width. The hand-back shows the two pictures side by side with both page paths, and says for each page how many of the relations can be read on the first screen without any click.
- The hand-back does not pick a winner. The maintainer picks.
## Constraints
- No file under `skills/` or `suite/` changes, so this is no release.
- Every existing file in the bundle keeps its bytes, except the iteration list that the second-opinion launcher rebuilds in `index.html`.
- The page passes `ai-report-judge` through the second-opinion finalize step.
- The pictures for the comparison are not committed into the bundle. They live in the run's own folder.
## Dependencies
REQ-733 (rebuild the relation map page) completes first, because the comparison needs the improved map.
## Red-Green Proof
**RED prompt/case:** Look in the report folder for a page that shows the options with what supports and what limits each one, in one table.
**Why RED now:** No such page exists. The relations can be read only on the map, or as 37 separate rows.
**GREEN when:** The table page exists and passes the render judge; every one of the 37 relations appears on it; the hand-back shows the map and the table side by side at 1440 pixels and at phone width, with both paths.
**Validation:** Inferred during capture, from the maintainer's own words quoted above.
## Builder Guidance
- Certainty: firm. The table form is quoted by the maintainer; the cell layout is the builder's.
- Scope cue from the maintainer: keep language clear, do not complicate it.
- Trap: second-opinion refuses to start when Git shows changes outside the bundle, and a claimed request file is such a change in the main tree. Run the prepare and finalize steps in this request's own worktree, where the tree is clean.
- The maintainer's choice between map and table decides what REQ-727 (ship a relation-map preset) becomes.
## Full Context
See `do-work/user-requests/UR-178/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: build this one too, a table with one row per option and columns supports and limits; compare the two*
