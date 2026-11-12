---
id: REQ-720
title: '[impact-negligible] Count work-reference.md sections and always-loaded crew files in the reading-volume measurement'
status: completed
route: A
heavy_verified_at: 2026-11-12T09:31:03Z
heavy_verified_revision: ca851acf3c95aafe136e29d101606d89445ec1dc
commit: 3a91903c
review_at: 2026-11-12T09:27:29Z
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-12T09:19:00Z
  basis:
    - Route A
    - 2-file write set
    - 6 acceptance criteria
created_at: 2026-11-12T09:05:50Z
user_request: UR-171
domain: general
prime_files: []
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
related: [REQ-719]
batch: phase-loading-follow-ups
write_set: ["_dev/tests/instruction-reading-volume.py", "_dev/reports/work-phase-loading.md"]
claimed_at: 2026-11-12T09:18:38Z
completed_at: 2026-11-12T09:31:15Z
---
# Count work-reference.md Sections and Always-Loaded Crew Files in the Reading-Volume Measurement
## What
Extend `_dev/tests/instruction-reading-volume.py` and the report `_dev/reports/work-phase-loading.md` so the three walkthroughs also count `work-reference.md`, by cited section, and the crew files that are always loaded. Measurement only: no shipped file changes.
## Why
The current totals leave out the largest instruction file. The report excludes `work-reference.md` (119,671 bytes), yet the selection phase file sends even a simple fix to two of its sections that total 52,760 bytes. The phase split saved 12,693 bytes, so the reported saving looks larger than it is against what an agent really reads.
## Context
- This request comes from Finding 5 of the validate-feedback triage of 2026-11-11 (source: pasted text, no severity given). The triage verdict was Accept, and its surface-cost result was N/A: this is a measurement fix, not a new defensive layer.
- `_dev/reports/work-phase-loading.md:18` states the current boundary: the router, `work.md`, the phase files and `review-work.md`.
- `skills/do-work/actions/work-selection-reference.md:53` points a simple fix at the Request File Schema section (35,527 bytes) and the Schema Read Contract section (17,233 bytes) of `skills/do-work/actions/work-reference.md`.
- The script compares a fixed baseline revision with the files on disk (`_dev/tests/instruction-reading-volume.py:11`), and the same file holds the preservation replay for the phase split. The builder edits that file.
- Fold-first scan on 2026-11-12: the queue held no pending or pending-answers REQ, so no candidate in any UR shares this root cause.
## Detailed Requirements
- Extend `_dev/tests/instruction-reading-volume.py` and the report `_dev/reports/work-phase-loading.md`.
- Count `work-reference.md` by cited section: the sections that each walkthrough's loaded files point to, not the whole file and not zero.
- Count the always-loaded crew files.
- Keep the same three walkthroughs: a simple mechanical fix, a complex multi-file change, and a review-only run.
- Record, for each walkthrough, the exact instruction files loaded and the total bytes read.
## Constraints
- Measurement only. No shipped file changes: nothing under `skills/` changes.
- Byte or token volume is used to spot unnecessary reading in unrelated phases. It is never proof of instruction quality and never a reason to delete a required rule.
## Red-Green Proof
**RED prompt/case:** Run `python3 _dev/tests/instruction-reading-volume.py` from the repository root.
**Why RED now:** The Simple fix row ends at 133281 bytes. No output line names `work-reference.md` or a crew file, so those bytes are not counted.
**GREEN when:** For each of the three walkthroughs, the script lists the cited `work-reference.md` sections and the always-loaded crew files with their bytes, and its totals include them. The table in `_dev/reports/work-phase-loading.md` shows the same totals the script prints. This script run, compared with the report table, is the regression check for the finding.
**Validation:** User confirmed on 2026-11-12: the proof is the script output and the report table together.
## Dependencies
None. REQ-719 (the prime byte budget) belongs to the same batch and does not depend on this request.
## Builder Guidance
The user is certain about what to count and where. How the script finds which `work-reference.md` sections a walkthrough cites, and which crew files count as always loaded for each walkthrough, are the builder's decisions; state both in the report so a reader can check them. The finding's verbatim claim also names checks, recovery rules and a `code-review` path; the remedy the user captured is the counting work above.
## Full Context
See `do-work/user-requests/UR-171/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Finding 5 of the 2026-11-11 phase-loading triage, as pasted by the user in UR-171: "Record the exact instruction files loaded, total bytes read, and every required check and recovery rule reached."*

## Triage

**Route: A** - Simple

**Reasoning:** Two named maintainer files with a fixed three-scenario measurement and script/report acceptance proof.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

Implementation approach: keep the pinned baseline and three scenarios, count reference sections outside fenced examples, deduplicate overlapping section reads, enumerate all read units on both sides, and regenerate the report table from script output. Claim-time lessons index consulted: action satellite is relevant but exceeds the shared budget and is partial; no mandatory lesson entry captured. Always-loaded crews read from the current pipeline instructions. No missing entries.

## Implementation Summary

**Files changed:**
- `_dev/tests/instruction-reading-volume.py` (modified)
- `_dev/reports/work-phase-loading.md` (modified)

**What was done:** Counts actual cited reference sections and crew files symmetrically at the pinned baseline and current candidate, prints exact per-unit bytes and unique file counts, excludes fenced example headings and deduplicates overlapping reads. Report totals match script output.

UNIFY evidence: both changed files inspected, git diff --check passed, script preservation replay passed, and script/report parity, schema inclusion, reference boundaries and review-only isolation checked. A first scratch assertion used the entry-only path list instead of the complete walkthrough; correcting that assertion confirmed the loaded-phase behavior.

## Decisions

- **D-01 — DECIDE & STATE:** Match named real reference headings, follow nested mentions to a fixed point, and subsume child sections under already-read parents. Conservative static accounting includes named conditional branches; external reference files and task-specific crews remain outside the stated boundary. No claim about total task context or runtime improvements.

## Qualification

Passed — substantive two-file diff a2fe2556..3a91903c matches the request and manifest. Advance reported QUALIFY-LIBRARY-OUTPUT for the new print calls: accepted by orchestrator and independent reviewer because this standalone script explicitly produces the requested human-readable report. These lines are contract output, not leftover instrumentation. No source change solely to appease the implicit-exit heuristic.

## Testing

**Tests run:** python3 _dev/tests/instruction-reading-volume.py; script/report table comparison; reference-boundary/schema-inclusion/review-isolation checks; bash _dev/tests/maintainer-verify.sh.
**Result:** All passed. Preservation replay passed. Canonical gate exited 0 in 131s; 427 board tests and 980 CLI tests passed. Request-bound focused/green gate records satisfied.
**Red-green validation:** Old simple-fix output omitted reference/crew files and reported 133,281 bytes. Updated output lists each read unit, counts 301,147 baseline bytes / 286,689 candidate bytes, and matches the report. Complex: 301,147 / 301,715; review-only: 103,725 / 103,717.
**Heavy verification plan:** a2fe2556..3a91903c; all six lanes selected because the maintainer test subtree changed.

## Review

**Overall: 99%**

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — implementation and merged integration; independent reviewer reproduced all report rows, checked fenced-heading exclusion, overlap deduplication, schema inclusion and the disclosed static boundary. No consumer deployment applies to this maintainer-only change.
**Restatement sweep:** redefined instruction measurement boundary and totals; script, report method, table and interpretation agree. Reviewed earlier REQ-719 prime budget element: create/audit/shrink/checklist/crew rules remain consistent. No stale restatement found.
**Suggested testing:** the six planner-selected heavy lanes at queue exhaustion.
**Follow-ups created:** None.

## Orientation

The instruction walkthrough report now exposes the schema/reference and crew reading it previously omitted, with reproducible file counts and byte totals for each scenario and revision.

## Heavy Verification Plan

- Base: a2fe25561f44fc860c14e5bac61b43dba5edf501
- Target: 3a91903ce62dc9c1164e8ed7c89497114803bcba
- Lanes: queue-kanban-javascript, queue-kanban-browser, do-work-cli-integrations, staged-skills, updater, installer.
- Each lane argv: env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane-id>
- Reason: _dev/tests/instruction-reading-volume.py matched the maintainer-test subtree.

## Timing

Observed 2026-11-12T09:21:51Z to 2026-11-12T09:24:02Z: 2m 11s total, 2m 12s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 2m 12s | 1 |

Slowest command: verification-gate / maintainer-verify, 2m 12s, exit 0, bash (2 argv tokens).

## Heavy Verification Result

Target: 3a91903c. Execution: ca851acf3c95aafe136e29d101606d89445ec1dc. All selected lanes executed and passed: JavaScript 11s, CLI integrations 65s, staged-skills 40s, updater 71s, installer 28s, browser 88s. Initial browser skip was resolved using installed Chrome via QUEUE_KANBAN_BROWSER; the separate typed result reports exit 0, not skipped. No lane result was reused.
