---
id: REQ-531
title: 'Review findings below impact-critical stay in the report'
status: claimed
created_at: 2026-09-03T11:42:36Z
user_request: UR-102
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-532]
batch: close-the-tap
write_set:
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/capture-reference.md
  - _dev/tests/contract-regressions.sh
claimed_at: 2026-09-03T18:19:18Z
---

# Review Findings Below Impact-Critical Stay in the Report

## What

A review or build finding whose impact is anything other than `impact-critical` is recorded where it was found and never creates a REQ file. The maintainer reads the record and runs `do-work capture` by hand for the ones worth building. Only `impact-critical` still auto-queues.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"The entire development cycle seems to be a never ending story, now we have 46 reqs and they get more faster then being implemented." Measured 2026-09-03: 214 REQs created against 195 completed since 2026-08-20; 29 of 45 pending REQs carry `review_generated: true` or `addendum_to`; 0 of 45 carry `impact-negligible`, so the existing `--skip-impact-negligible` brake removes nothing. The 2026-08-20 treadmill report already concluded "the exit is a rule change, not a sprint".

## Context

Three flows mint REQs from findings today, and this REQ changes all three by condition, not by list:

- `skills/do-work/actions/review-work.md` Step 10 (Create Follow-up REQs): every Important finding becomes a `Review fix:` REQ, a sweep REQ, or a `pending-answers` REQ at generation 2 or deeper.
- `skills/do-work/actions/work-reference.md` → Discovered Tasks Classification (Step 8): every non-critical discovery becomes a `pending-answers` REQ, and the test-hygiene carve-out auto-queues some as `pending`.
- `skills/do-work/actions/capture-reference.md` → Fold-First Rule, destination 4 (no match → a new REQ), which the two flows above and `../../do-work-toolbox/actions/code-review.md` share.

The durable record already exists in each flow: the `## Review` section appended to the reviewed REQ (with the mandatory impact token per finding) and the builder's `## Discovered Tasks` section. Those sections become the finding's home.

Decisions taken at capture:

- D1 Scope is REQ creation from a finding or discovery. Appending an instance to an existing pending sweep (Fold-First destinations 1 and 2) and appending to `do-work/prose-backlog.md` (destination 3) are edits to existing records, not creation, and stay as they are.
- D2 `impact-critical` keeps its pierce: it still mints `status: pending` at any generation, with the prominent report line.
- D3 Failure Classification follow-ups (`work-reference.md` → Failure Classification (Step 8)) are out of scope. A failed REQ still gets its Intent/Spec/Code successor, or failed work dies silently.
- D4 The generation-2 cascade stop and its `pending-answers` consent path become unreachable for non-critical findings and are deleted rather than kept as dead prose.

## Detailed Requirements

- A non-critical review finding's line in the `## Review` section and in the report ends with `→ report only` instead of a REQ id. The token stays mandatory.
- A non-critical builder discovery stays in the archived REQ's `## Discovered Tasks` section with its token, and the work loop's Step 8 creates nothing for it. The test-hygiene carve-out is deleted with it.
- Fold-First destination 4 reads: no match and `impact-critical` → a new REQ; no match otherwise → the finding stays in the record that found it, named in the flow's report.
- The review report's `### Follow-up REQs Created` section lists critical REQs only and says `None (N findings report only)` otherwise.
- `do-work roadmap` or the review report tells the maintainer how to promote a report-only finding: `do-work capture` quoting the finding line. One sentence, no new action.
- Every sentence pin in `_dev/tests/contract-regressions.sh` that asserts the old minting behavior (about 30 lines match `Review fix`, `Create Follow-up REQs`, or `review_generated`) is deleted or rewritten to the new condition, without the file growing.
- The suite version bump and changelog entry name the rule change so a user reading only headings sees that reviews no longer mint REQs.

## Constraints

- Delete before you add: prefer removing the cascade-depth and consent machinery over adding a switch.
- No new `--flag` on `do-work run` or `do-work review-work`; the rule is unconditional.
- Overlap declared, not a dependency: REQ-504 and REQ-507 (UR-098) also rewrite `work.md` Step 8 prose; REQ-484 (UR-091) adds a Fold-First destination. Whichever lands second rebases.

## Red-Green Proof
**RED prompt/case:** Run `do-work review-work REQ-NNN` on an archived REQ whose review finds one Important `impact-user-visible` issue, then `ls do-work/queue/`.
**Why RED now:** A new `REQ-MMM-review-fix-*.md` appears in the queue with `review_generated: true`; Step 10 mandates it.
**GREEN when:** The same review leaves `do-work/queue/` unchanged, the finding line in the archived REQ's `## Review` section ends with `impact-user-visible → report only`, and an `impact-critical` finding in the same run still produces exactly one `status: pending` REQ.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes action routing and a status contract.

## Full Context
See `do-work/user-requests/UR-102/input.md` for complete verbatim input.

---
*Source: D1 of the 2026-09-03 roadmap disposition, selected by the maintainer.*
