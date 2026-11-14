---
id: REQ-726
title: '[impact-negligible] Record the decision to adopt no external review canvas'
status: claimed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: []
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-727, REQ-728, REQ-729, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["docs/decisions/records/adr-023-*.md", "docs/decisions/log.md", "docs/decisions/_master_index.md", "docs/decisions/topics/_index_pipeline-deliverables.md"]
claimed_at: 2026-11-14T23:28:12Z
---
# Record the Decision to Adopt No External Review Canvas
## What
Write one decision record: do-work adopts no external review canvas (Tinka) for ai-report bundles. The record exists so a later session does not propose it again without new facts.
## Why
Two proposals, one request and three versions went into Tinka on 2026-11-14, and the maintainer then abandoned it. Without a record, the next session can spend the same effort again.
## Context
- Tinka (github.com/bladnman/tinka) is a local board tool for trying designs with an AI agent. Two proposals tested it as a place to review report pages: `ai-reports/2026-11-14_1452_proposal-tinka-review-surface` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface`.
- REQ-721 (every AI report ends with a Tinka canvas offer), archived under `do-work/archive/UR-172/`, shipped a canvas offer in every report plus a hand recipe as version 0.316.0. Version 0.316.1 fixed the canvas id. Version 0.317.0 removed both the same day. The three entries are in `CHANGELOG.md`.
- The page `index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` in the second proposal's bundle is the direction taken instead: a zoomable relation map inside the report.
## Detailed Requirements
- The record states the decision: no external review canvas for ai-report bundles.
- It states what was observed. A report page is about 5,000 pixels tall; as a picture tile it is a thin strip that cannot be read, and folded sections stay closed (38 to 46% of the text, measured in the second proposal).
- Tinka answers only on 127.0.0.1, so a phone cannot open it.
- A live HTML tile served from Tinka's own address can write feedback by itself. A test page did it.
- A live tile served from a second local port cannot write (tested 2026-11-14: the browser blocked the read and one write, Tinka refused two writes with HTTP 403). But it is live only after two clicks and one page at a time, which a browser tab already does.
- The marks Tinka adds (star, vote, note) move feedback from chat to a second place and add steps.
- It states what would reopen the question: Tinka shows several live pages side by side with an isolation that stops them from writing feedback, or do-work gets a project whose work is many rounds of screen designs.
- Files: a new `docs/decisions/records/adr-023-<slug>.md` in the format of `adr-022-one-broken-pipe-does-not-stop-the-factory.md`; links from `docs/decisions/log.md`, `docs/decisions/_master_index.md` and the fitting topic index under `docs/decisions/topics/` (probably `_index_pipeline-deliverables.md`), as the existing records are linked.
## Constraints
- Nothing under `skills/` or `suite/` changes. This is a maintainer record, so it is no release.
- No existing `ai-reports/` bundle changes.
- No new command and no new file outside `docs/decisions/`.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** Search the decision records for this decision: `ls docs/decisions/records | grep -c 'review-canvas'`.
**Why RED now:** It prints 0. No record names the decision, the two proposals or the three versions.
**GREEN when:** The record exists with status accepted and names the two proposals, REQ-721 and the three versions as sources; the log and the topic index link it; `git diff --stat HEAD -- skills/ suite/` is empty; `bash _dev/tests/contract-regressions.sh` passes.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- Certainty: firm. The maintainer said "let's abandon tinka" on 2026-11-15.
- The record number is the next free one at build time; 023 is the number seen at capture.
- The slug in the file name is the builder's choice, two words minimum, and it contains `review-canvas` so the proof grep finds it.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-record-no-review-canvas-decision.md*
