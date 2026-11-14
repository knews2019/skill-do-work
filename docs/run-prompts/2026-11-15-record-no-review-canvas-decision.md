# Capture prompt: record the decision to adopt no external review canvas

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.

## The request

Write one decision record: do-work adopts no external review canvas (Tinka) for ai-report bundles. The record exists so a later session does not propose it again without new facts.

### Context

- Tinka (github.com/bladnman/tinka) is a local board tool for trying designs with an AI agent. Two proposals tested it on 2026-11-14 as a place to review report pages: `ai-reports/2026-11-14_1452_proposal-tinka-review-surface` and `ai-reports/2026-11-14_1547_proposal-tinka-review-surface`.
- REQ-721 (archived under `do-work/archive/UR-172/`) shipped a canvas offer in every report plus a hand recipe as version 0.316.0. Version 0.316.1 fixed the canvas id. Version 0.317.0 removed both the same day. The three changelog entries are in `CHANGELOG.md`.
- What was observed:
  - A report page is about 5,000 pixels tall. As a picture tile it is a thin strip that cannot be read, and folded sections stay closed (38 to 46% of the text, measured in the second proposal).
  - Tinka answers only on 127.0.0.1, so a phone cannot open it.
  - A live HTML tile served from Tinka's own address can write feedback by itself. A test page did it.
  - A live tile served from a second local port cannot write (tested 2026-11-14: the browser blocked the read and one write, Tinka refused two writes with HTTP 403). But it is live only after two clicks and one page at a time, which a browser tab already does.
  - The marks Tinka adds (star, vote, note) move feedback from chat to a second place and add steps.
- What would reopen the question: Tinka shows several live pages side by side with an isolation that stops them from writing feedback, or do-work gets a project whose work is many rounds of screen designs.
- The page `index-v1c03-claude-claude-fable-5-2-custom-zoomable-relation-map.html` in the second proposal's bundle is the direction taken instead: a zoomable relation map inside the report.

### Files to change

- `docs/decisions/records/adr-023-<slug>.md` (new), in the format of `adr-022-one-broken-pipe-does-not-stop-the-factory.md`.
- `docs/decisions/log.md`, `docs/decisions/_master_index.md` and the fitting topic index under `docs/decisions/topics/` (probably `_index_pipeline-deliverables.md`), as the existing records are linked.

### What must not change

- Nothing under `skills/` or `suite/`. This is a maintainer record, so it is no release.
- No existing `ai-reports/` bundle.
- No new command and no new file outside `docs/decisions/`.

### Verification

RED first: `ls docs/decisions/records | grep -c 'review-canvas'` prints 0.

GREEN: the record exists with status accepted and names the two proposals, REQ-721 and the three versions as sources; the log and the topic index link it; `git diff --stat HEAD -- skills/ suite/` is empty; `bash _dev/tests/contract-regressions.sh` passes.
