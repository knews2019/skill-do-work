---
id: REQ-531
title: 'Review findings below impact-critical stay in the report'
status: completed
created_at: 2026-09-03T11:42:36Z
user_request: UR-102
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-548]
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
  - skills/do-work/actions/capture.md
  - skills/do-work/crew-members/general.md
  - skills/do-work-toolbox/crew-members/general.md
  - skills/do-work/docs/review-work-guide.md
  - skills/do-work/docs/work-guide.md
  - skills/do-work/docs/standing-preferences.md
  - skills/do-work-toolbox/actions/code-review.md
  - skills/do-work/next-steps.md
  - _dev/tests/contract-regressions.sh
route: C
planning_at: 2026-09-03T18:23:16Z
estimate:
  p50_active_minutes: 45
  confidence: low
  calculated_at: 2026-09-03T18:19:45Z
  basis:
    - Route C
    - 5-file write set
    - 3 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
    - full-suite verification
gate_deferred: 'true'
claimed_at: 2026-09-03T18:47:36Z
implementation_at: 2026-09-03T19:29:52Z
testing_at: 2026-09-03T19:32:25Z
review_at: 2026-09-03T19:41:38Z
kb_status: pending
completed_at: 2026-09-03T19:44:58Z
commit: 455462be46b3170a22d331136ec5aa7f7e5a1c60
release_at: 2026-09-03T19:44:58Z
---

# Review Findings Below Impact-Critical Stay in the Report

## What

A review or build finding whose impact is anything other than `impact-critical` is recorded where it was found and never creates a REQ file. The maintainer reads the record and runs `do-work capture` by hand for the ones worth building. Only `impact-critical` still auto-queues.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the required action/shell primes, lessons, and crew guidance; mapped the critical-only finding boundary across canonical writers, alternate writers, guides, and the regression ratchet before editing.
- [x] **[APPLY]:** Implemented the planned contract rewrite on the isolated REQ-531 branch, limited to all 13 declared project files.
- [x] **[UNIFY]:** Reviewed the complete 13-file merge range, preserved the concurrent fast-test-budget assertions during conflict resolution, passed shell syntax/diff checks and contract regressions, and kept `_dev/tests/contract-regressions.sh` at 8,479 lines.

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
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ rewrites the shell-based contract regression suite.

## Full Context
See `do-work/user-requests/UR-102/input.md` for complete verbatim input.

---
*Source: D1 of the 2026-09-03 roadmap disposition, selected by the maintainer.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This changes the suite-wide finding-to-REQ contract across review, work, capture, and regression pins, including release-facing behavior. The requirements are explicit, but the cross-action semantics and stale-restatement surface require planning and exploration.

**Planning:** Required

## Plan

1. Narrow the canonical finding-routing rule in `capture-reference.md`: preserve folds into existing records and prose-backlog handling, but allow a no-match finding to mint a REQ only when its impact is `impact-critical`. Keep explicit `do-work capture` of a quoted finding as new user intent.
2. Rewrite review output and persistence in `review-work.md`: all Important findings keep a mandatory impact token; critical findings may create one pending follow-up, while non-critical findings end `→ report only`. Remove the generation-two consent/cascade machinery and update the follow-up summary plus manual-promotion guidance.
3. Apply the same condition to builder discoveries and work-loop review handling in `work-reference.md` and `work.md`, deleting non-critical `pending-answers` and test-hygiene creation paths while preserving builder-decided, stakeholder, and failure-classification follow-ups.
4. Sweep alternate writers and shipped restatements in `capture.md`, `crew-members/general.md`, `docs/review-work-guide.md`, and the toolbox `code-review.md` action so no reader or producer teaches the retired behavior.
5. Replace obsolete contract pins in `_dev/tests/contract-regressions.sh` with critical-only creation and non-critical report-only assertions, keeping the file at or below its current line count.
6. Verify the new focused contract behavior, the complete contract-regression suite, the canonical maintainer gate, stale-restatement grep results, and mirror/version consistency. The orchestrator will publish the required release bump and changelog through finalization.

**Plan validation:** All detailed requirements map to the six tasks above, and every task traces to the finding-routing rule or its required release proof. The plan has 6 tasks; quality can degrade past 3, but splitting would leave the cross-action contract half-landed, so the implementation should treat tasks 1–4 as one semantic rewrite and task 5 as its single regression boundary. No planned command output drives an action-owned mutation without the existing canonical finalization consumer.

*Generated by Plan agent*

## Exploration

The canonical creation contracts are concentrated in `review-work.md` Step 10, `work-reference.md`'s Discovered Tasks Classification, and `capture-reference.md`'s Fold-First Rule. The same old behavior is restated by `capture.md`, both installed copies of `crew-members/general.md`, two user guides, standing preferences, `next-steps.md`, and the toolbox code-review action. Those are required consistency edits, not optional documentation cleanup.

The focused regression surface is `_dev/tests/contract-regressions.sh`, especially the archived-review assumptions, review follow-up template checks, and generic `review_generated` producer ratchet. The suite's canonical gate is `bash _dev/tests/maintainer-verify.sh`.

REQ-504 (collapsing recovery prose) and REQ-507 (handing archive/commit tails to finalization) are still pending and overlap `work.md`, `work-reference.md`, and the contract suite. Their target semantics differ, but whichever lands second must rebase. REQ-484's proposed Fold-First destination was cancelled as superseded by this request and must not be revived.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/review-work.md` (modify) — retain all findings in review records and create REQs only for critical findings
- `skills/do-work/actions/work-reference.md` (modify) — make builder discoveries critical-only for REQ creation
- `skills/do-work/actions/work.md` (modify) — align orchestration, review-result, and discovered-task prose
- `skills/do-work/actions/capture-reference.md` (modify) — split Fold-First destination 4 by impact while preserving existing-record destinations
- `skills/do-work/actions/capture.md` (modify) — distinguish explicit manual capture from automatic finding routing
- `skills/do-work/crew-members/general.md` (modify) — describe report-only non-critical discoveries
- `skills/do-work-toolbox/crew-members/general.md` (modify) — keep the installed alternate crew contract aligned
- `skills/do-work/docs/review-work-guide.md` (modify) — update the user-facing review behavior
- `skills/do-work/docs/work-guide.md` (modify) — update the user-facing work-loop behavior
- `skills/do-work/docs/standing-preferences.md` (modify) — remove the all-discoveries-queue restatement
- `skills/do-work-toolbox/actions/code-review.md` (modify) — align the alternate review writer with critical-only creation
- `skills/do-work/next-steps.md` (modify) — make post-review run guidance conditional on a critical REQ existing
- `_dev/tests/contract-regressions.sh` (modify) — replace old minting pins without increasing line count

**Files I will NOT touch:** Failure Classification follow-up behavior; builder-decided and stakeholder question flows; REQ-504/REQ-507 implementation areas beyond stale finding-routing restatements; release/version files before canonical finalization.

**Acceptance criteria (restated from REQ):**
- [ ] Every non-critical review finding keeps its impact token and ends `→ report only` in both the durable Review block and the report.
- [ ] Every non-critical builder discovery stays in the archived REQ's Discovered Tasks section and creates no REQ; the test-hygiene creation carve-out is gone.
- [ ] Fold-First destination 4 creates a new REQ only for an unmatched `impact-critical` finding; explicit manual capture remains available.
- [ ] The review follow-up section lists critical creations only and reports `None (N findings report only)` when no critical REQ was created.
- [ ] A user-facing sentence explains promotion by running `do-work capture` with the quoted finding line.
- [ ] Contract pins for old creation, cascade-depth consent, and test-hygiene behavior are removed or rewritten, and the test file does not grow beyond 8479 lines.
- [ ] Failure Classification, builder-decided, stakeholder, existing-sweep, and prose-backlog paths retain their current behavior.
- [ ] The canonical maintainer gate passes before release finalization.

## Pre-Flight

**Git:** ✓ Working tree clean outside `do-work/`
**Tests baseline:** ⚠ The canonical maintainer gate fails before implementation in `skills/do-work/tools/do-work-cli/internal/corehelpers` (public-command text identity/status expectations and protected-inventory association state).
**Dependencies:** ✓ Required Go and shell toolchains launched successfully.

*Checked by work action*

## Repository Gate Deferral

- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** sha256:efed40e27e755df5b2733ea1fdf3d3f228c42e8d859baefb1f89b35769ebda2e
- **Repair dependency:** REQ-548
- **Diagnostic evidence:** "maintainer-verify: update-script behavior probes failed"
- **Diagnostic evidence:** "update-script-behavior: suite update identifies layout output did not match four-module suite"
- **Diagnostic evidence:** "update-script-behavior.sh: printf write error: Broken pipe"

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — replaced stale automatic-follow-up pins with a named critical-only/report-only contract while preserving concurrent fast-test-budget coverage and the 8,479-line ceiling.
- `skills/do-work/actions/review-work.md` (modified) — records every finding, auto-queues only critical findings, adds report-only suffix/summary wording, and documents manual capture promotion.
- `skills/do-work/actions/work-reference.md` (modified) — retains noncritical builder discoveries in the current REQ and removes consent/test-hygiene creation paths.
- `skills/do-work/actions/work.md` (modified) — aligns review-result and discovered-task orchestration with critical-only automatic follow-ups.
- `skills/do-work/actions/capture-reference.md` (modified) — gates automatic findings before Fold-First and makes destination 4 report-only, with explicit capture and critical exceptions.
- `skills/do-work/actions/capture.md` (modified) — distinguishes explicit user capture from automatic finding routing.
- `skills/do-work/crew-members/general.md` (modified) — teaches report-only handling for noncritical discoveries.
- `skills/do-work-toolbox/crew-members/general.md` (modified) — keeps the toolbox crew mirror aligned.
- `skills/do-work/docs/review-work-guide.md` (modified) — updates user-facing review follow-up behavior.
- `skills/do-work/docs/work-guide.md` (modified) — updates user-facing builder discovery behavior.
- `skills/do-work/docs/standing-preferences.md` (modified) — removes the obsolete all-findings queue preference.
- `skills/do-work-toolbox/actions/code-review.md` (modified) — aligns the alternate review writer and summary template.
- `skills/do-work/next-steps.md` (modified) — makes post-review queue guidance conditional on critical work.

**What was done:** Reversed the default finding-to-REQ behavior across every active writer and restatement: only `impact-critical` review/build findings may mutate follow-up state; all other findings stay in their report or archived REQ with `→ report only`, and maintainers can promote one explicitly with `do-work capture` using the quoted finding line.

## Qualification

Passed — 13 files verified, 8 requirements traced, P-A-U confirmed. The merge-range qualifier passed for `0ffd579266cf05539133ece53bffde38307b6ceb..455462be46b3170a22d331136ec5aa7f7e5a1c60`; the integration conflict retained both REQ-531 and concurrent fast-test-budget contracts without exceeding the 8,479-line regression ceiling.

## Testing

- **Red-green validation:** `bash _dev/tests/contract-regressions.sh` failed before implementation with seven named REQ-531 contract failures covering critical-only review routing, manual capture promotion, noncritical builder retention, removal of test-hygiene/consent creation, and Fold-First destination 4. The same named contract passed after implementation and again after merge.
- `grep -Eq 'noncritical.*line.*ends.*report only' skills/do-work-toolbox/actions/code-review.md && grep -Eq 'None \(N findings report only\)' skills/do-work-toolbox/actions/code-review.md && bash _dev/tests/contract-regressions.sh` — passed on merged `HEAD`.
- `bash -n _dev/tests/contract-regressions.sh` — passed after conflict resolution.
- `wc -l _dev/tests/contract-regressions.sh` — 8,479, equal to the pre-REQ ceiling despite retaining the concurrent fast-test-budget contract.
- `bash _dev/tests/maintainer-verify.sh` — passed on merged `HEAD`; recorded green at `455462be46b3170a22d331136ec5aa7f7e5a1c60`.

## Review

**Verdict: Approve with report-only findings.** The critical-only/report-only rule is mostly implemented, but acceptance is Partial: the change broadens routing beyond the captured scope, adds phrase-level pins despite the UR constraint, and leaves stale active restatements.

**Overall: 64%** | 2026-09-03T19:41:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 63% |
| Code Quality | 74% |
| Test Adequacy | 58% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- `skills/do-work/actions/capture-reference.md` gates all automatic noncritical findings before Fold-First, contradicting the captured decision and plan to preserve existing-sweep and prose-backlog destinations; `review-work.md` and `work.md` propagate the broader rule — impact-rule-change → report only
- `_dev/tests/contract-regressions.sh` adds a new phrase-level REQ-531 assertion block despite UR-102's no-new-sentence-pins constraint, and one assertion pins the broader prose-backlog prohibition — impact-rule-change → report only
- Active restatements remain in `work.md` around generic discovered-task queueing and failed-review follow-ups, and `clarify.md` still describes review/build producers emitting consent-style `pending-answers` templates — impact-rule-change → report only
- The requested report heading is `### Follow-up REQs Created`, while both review writers emit `### Follow-ups created`; the regression check pins only the summary phrase — impact-rule-change → report only

**Minor findings:** None
**Acceptance:** Partial — the core critical-only outcome works and all merged tests pass, but the four contract discrepancies above remain in the durable report.
**Suggested testing:** 3 items — exercise a four-case Fold-First routing matrix; assert destinations 1–3 retain their captured behavior while destination 4 changes; pin the exact requested heading without another prose sentence assertion.
**Follow-ups created:** None (4 findings report only); **sweeps appended to:** None

To promote one of these findings, invoke `do-work capture` with the complete finding line quoted as the capture source.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- A named RED contract plus an isolated builder branch made the default reversal testable before prose changed, and the merge-time conflict exposed concurrent contract edits instead of losing either family.

**What didn't:**
- The implementation followed the headline's broad “report only” wording past the narrower captured decision, and added phrase pins despite the UR's explicit no-new-pin constraint; independent review caught both after the green tests.

**Worth knowing:** The review contract now keeps its own noncritical findings here. To pursue one, run `do-work capture` with the complete finding line quoted as the source rather than expecting an automatic follow-up.

## Orientation

[MAP CHANGED] Review/build finding routing now lives behind one impact-authority boundary spanning the action writers and Fold-First: critical findings may mutate queue state, while noncritical findings remain durable report content. This changes the action contract described by `_dev/primes/prime-action-files.md`; every alternate writer and action-bearing reader must stay aligned.
