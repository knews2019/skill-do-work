---
id: REQ-726
title: '[impact-negligible] Record the decision to adopt no external review canvas'
status: completed
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
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-14T23:30:08Z
  basis:
    - Route A
    - 4-file write set
    - 6 acceptance criteria
dispatch_at: 2026-11-14T23:32:56Z
builder_handback_at: 2026-11-15T02:47:34Z
integration_at: 2026-11-15T02:47:45Z
review_at: 2026-11-15T02:52:56Z
kb_status: pending
commit: 66e3128042eff0417b6feda7672ab520662eeda1
completed_at: 2026-11-15T02:53:07Z
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
- [x] **[PLAN]:** copy adr-022's shape (frontmatter, topic line, Context/Decision/Consequences), shorter; link per D-01 to D-03. (from the builder hand-back)
- [x] **[APPLY]:** wrote the record and three index edits; one commit by exact path. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat main...HEAD` shows _master_index.md, log.md, adr-023 record, _index_pipeline-deliverables.md only; nothing under skills/, suite/, ai-reports/, do-work/; no debug artifacts. (from the builder hand-back)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-record-no-review-canvas-decision.md*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names every file to write (one new ADR in the adr-022 format plus three index links) and states every fact the record must hold; nothing is unknown about location or pattern. The change is docs-only under the export-ignored `docs/decisions/`, so no code, no tests to write and no release.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

Pre-dispatch decisions, made without asking (unattended run). The builder continues at D-05.

- **D-01** (ESCALATE): link the record from `docs/decisions/topics/_index_pipeline-deliverables.md`, the topic the REQ names. That topic is titled "Retired Pipeline Deliverables" and holds only superseded ADRs, so the builder widens its title and opening sentence to cover the live ai-report deliverable too, lists the new record as accepted, and updates that cluster's line in `_master_index.md` ("2 superseded ADRs" becomes a count that includes the accepted one). Reasoning: the write_set and the REQ both point there, and ai-report bundles are what the pipeline delivers to a reviewer. Value: no new topic file and no write outside the write_set. Risk: low; a later reader who expects that topic to be history only sees one live record there; moving the link is a two-file edit.
- **D-02** (DECIDE & STATE): the record's `decided:` date is 2026-11-15, the day the maintainer said "let's abandon tinka"; `created:` and `updated:` are the build date. The removal release 0.317.0 (2026-11-14) is cited as a source, not as the decision date.
- **D-03** (DECIDE & STATE): the builder adds only this record to `_master_index.md` and does not backfill the older records (ADR-017 to ADR-022) that its cluster lines do not list; that gap is a Discovered Task.
- **D-04** (DECIDE & STATE): `docs/decisions/_progress.md` still says "Next ADR number: ADR-017". It is outside the write_set, so it stays unchanged; the builder reports it as a Discovered Task.

Builder decisions (from the builder hand-back):

- **D-05** (DECIDE & STATE): `created:`/`updated:` are 2026-11-14, the UTC build date (2026-11-14T23:33Z), per D-02; `decided:` stays 2026-11-15. The log section header uses the decision date 2026-11-15. If the integrator prefers local date, change three `updated:`/`created:` lines.
- **D-06** (DECIDE & STATE): added a `## When to Reopen` section to the record rather than burying the two reopen conditions in Consequences; adr-022 has extra sections (Alternatives, Instances) so the format allows it. `related: []` because no existing ADR covers ai-report review.
- **D-07** (DECIDE & STATE): the record quotes the maintainer's "let's abandon tinka" (from the REQ's Builder Guidance) as the decision's trigger.

<!-- D-XX counter: last used D-07. Next decision: D-08. -->

## Implementation Summary

**Files changed:**
- `docs/decisions/records/adr-023-adopt-no-external-review-canvas.md` (new)
- `docs/decisions/log.md` (modified)
- `docs/decisions/_master_index.md` (modified)
- `docs/decisions/topics/_index_pipeline-deliverables.md` (modified)

**What was done:** Wrote ADR-023, the accepted record that do-work adopts no external review canvas for ai-report bundles, with the observed facts, the relation-map direction and two reopen conditions. Linked it from the decision log, the master index and the widened Pipeline Deliverables topic index.

## Discovered Tasks

(from the builder hand-back)

- `docs/decisions/_master_index.md` cluster lines do not list ADR-017 to ADR-022 (D-03 gap); a backfill is needed. → report only
- `docs/decisions/_progress.md` still says "Next ADR number: ADR-017" (D-04); it should say ADR-024 after this lands. → report only

## Qualification

**Gate records:** `advance --diff-range 05217478..66e31280` returned phase `qualify`, gate record `qualify` (provenance `merged_range`) satisfied, no findings.

**Requirement trace against `git diff 05217478..66e31280` (4 files, +71/-6):**
- Decision stated: ADR-023 `## Decision` says do-work adopts no external review canvas for ai-report bundles; `status: accepted`.
- ~5,000 px page, unreadable strip tile, folded sections closed holding 38 to 46% of the text: first Context bullet.
- 127.0.0.1 only, no phone: second bullet.
- Same-address live tile writes feedback by itself: third bullet.
- Second-port tile blocked (browser blocked the read and one write, two HTTP 403 refusals), live only after two clicks and one page at a time: fourth bullet.
- Star/vote/note move feedback out of chat: fifth bullet.
- Reopen conditions (several isolated live pages side by side; a project of many screen-design rounds): `## When to Reopen`.
- Sources: both proposal bundles, REQ-721's archive path, CHANGELOG 0.316.0/0.316.1/0.317.0, UR-176 input. Checked on main: both `ai-reports/` folders and the relation-map page exist; the REQ-721 archive file and UR-176 input are tracked; the three CHANGELOG headings exist at the cited titles.
- Links: `[[adr-023-adopt-no-external-review-canvas]]` appears in `log.md`, `_master_index.md` and `_index_pipeline-deliverables.md` (topic index widened per D-01).
- No other `adr-023-*` exists on main. The record has no link into `_dev/run-prompts/` (moved by REQ-732), so nothing needed fixing.

**Scope:** Route A; write_set = the four `docs/decisions/` paths, and the diff touches exactly those. `git diff --stat 05217478..66e31280 -- skills/ suite/ ai-reports/ do-work/` is empty. No debug artifacts.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge `66e3128042eff0417b6feda7672ab520662eeda1` (detached drain-head checkout, `DO_WORK_FAST_STAGE_REUSE=off`, 1-minute load 2.89 at start); GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-726-probe.sh` through `advance --probe-file`.
**Result:** ✓ Gate exit 0, 131 s, "Maintainer verification passed." Probe exit 0; gate records `run-blocked-check` and `green-gate` satisfied.

**Red-green validation:** (from the builder hand-back)
- `REQ-726-probe.sh`: ✗ at base ("no docs/decisions/records/adr-*review-canvas*.md record", exit 1) → ✓ at tip ("PASS (adr-023-adopt-no-external-review-canvas)", exit 0)

**New tests added:**
- None (docs-only record; the run probe is run evidence, not a shipped test).

**Heavy verification plan:** none selected — `plan-heavy-verification` over 052174786428be10b5e768fcd51619a803a3ed6d..66e3128042eff0417b6feda7672ab520662eeda1 selected no lanes (only `docs/decisions/` changed, no uncovered paths).

## Review

**Overall: 97%** | 2026-11-15T02:52:56Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** F1 ADR-023 lacks adr-022's `## References` section with relative links; sources are frontmatter only — impact-negligible → report only. F2 (nit) the 5,000-pixel and second-port facts come from the UR-176 capture prompt, not the bundles; the record cites UR-176 as a source — impact-negligible → report only. F3 (nit) anti-bloat: 1 unnamed section (`## When to Reopen`, holds required content), 0 helpers, 0 options, 0 files, 0 tests — impact-negligible → report only.
**Acceptance:** Pass — the probe passes and `ls docs/decisions/records | grep -c review-canvas` prints 1; no change under skills/, suite/ or ai-reports/.
**Restatement sweep:** skipped (not the wave's last successful integration)
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

Independent reviewer report: `do-work/runs/work-2026-11-15-012821/REQ-726-review.md` (Route A quick scan). Verdict Approve at 97% (≥75% → continue); no fix taken, all three findings report only.

## Heavy Verification Plan

- Base: 052174786428be10b5e768fcd51619a803a3ed6d
- Target: 66e3128042eff0417b6feda7672ab520662eeda1
- Lanes: none selected (`plan-heavy-verification`: only `docs/decisions/` changed, no uncovered paths). No drain.

## Lessons Learned

**What worked:** A pre-dispatch decision list (D-01 to D-04) that settled the topic placement and dates let a docs-only builder finish without seams.
**What didn't:** Nothing beyond the hand-back, which proposed no lesson.
**Worth knowing:** `docs/decisions/_progress.md` and the master-index cluster lines lag behind the records (ADR-017 to ADR-022 unlisted); both are Discovered Tasks.

## Orientation

Decision records live in `docs/decisions/` (export-ignored, maintainer-only): one file per ADR under `records/`, linked from `log.md`, `_master_index.md` and one `topics/_index_*.md` cluster.

## Timing

Observed 2026-11-14T23:32:56Z to 2026-11-15T02:52:56Z: 3h 20m 00s total, 3h 19m 04s attributed across 4 events, 56s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 3h 14m 38s | 1 |
| verification-gate | 2m 23s | 1 |
| review | 1m 52s | 1 |
| handback-merge | 11s | 1 |

Slowest stage: builder-work / builder worktree build, 3h 14m 38s, outcome success.
