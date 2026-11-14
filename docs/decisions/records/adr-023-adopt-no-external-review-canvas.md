---
title: "ADR-023: Adopt No External Review Canvas"
type: architecture-decision-record
status: accepted
topic_cluster: pipeline-deliverables
decided: 2026-11-15
sources:
  - ai-reports/2026-11-14_1452_proposal-tinka-review-surface/
  - ai-reports/2026-11-14_1547_proposal-tinka-review-surface/
  - do-work/archive/UR-172/REQ-721-ai-report-closing-suggests-tinka-canvas.md
  - CHANGELOG.md (0.316.0 Every Report Offers a Tinka Canvas)
  - CHANGELOG.md (0.316.1 Tinka Canvas Opens in the Browser)
  - CHANGELOG.md (0.317.0 Tinka Canvas Offer Removed)
  - do-work/user-requests/UR-176/input.md
related: []
created: 2026-11-14
updated: 2026-11-14
confidence: high
---

# ADR-023: Adopt No External Review Canvas

Topic cluster: [[_index_pipeline-deliverables]] ([topic index](../topics/_index_pipeline-deliverables.md))

## Context

Tinka (github.com/bladnman/tinka) is a local board tool for trying designs with an AI agent. On 2026-11-14 two proposals tested it as a place to review ai-report pages: `ai-reports/2026-11-14_1452_proposal-tinka-review-surface` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface`.

REQ-721 (every AI report ends with a Tinka canvas offer) shipped that offer plus a hand recipe in the ai-report guide as version 0.316.0. Version 0.316.1 fixed the canvas id. Version 0.317.0 removed the offer and the recipe on the same day. On 2026-11-15 the maintainer said "let's abandon tinka".

What the two proposals and the pilot observed:

- A report page is about 5,000 pixels tall. As a picture tile it is a thin strip that cannot be read. Folded sections stay closed in a capture, and they hold 38 to 46% of the text (measured in the second proposal).
- Tinka answers only on 127.0.0.1, so a phone cannot open it.
- A live HTML tile served from Tinka's own address can write feedback by itself. A test page did it.
- A live tile served from a second local port cannot write. In a test on 2026-11-14 the browser blocked the read and one write, and Tinka refused two writes with HTTP 403. But such a tile is live only after two clicks, and only one page at a time. A browser tab already does that.
- The marks Tinka adds (star, vote, note) move feedback from chat to a second place and add steps.

## Decision

**do-work adopts no external review canvas for ai-report bundles.** Reports are reviewed as pages in a browser, and feedback stays in chat.

The direction taken instead is a zoomable relation map inside the report itself. The example is `index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` in the second proposal's bundle.

## When to Reopen

Reopen this question only when one of these becomes true:

1. Tinka shows several live pages side by side, with an isolation that stops them from writing feedback.
2. do-work gets a project whose work is many rounds of screen designs.

Without one of these, a new proposal for a review canvas repeats the work above.

## Consequences

Reports keep one review place: the page in a browser tab, with feedback in chat. No canvas recipe, offer or `.tinka/` folder ships with the skill. The cost is that a lineage of report pages cannot be seen side by side on one board.
