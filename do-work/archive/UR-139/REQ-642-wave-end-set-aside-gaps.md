---
id: REQ-642
title: '[impact-rule-change] Addendum: the wave-end decision reads the set-aside list, and a late set-aside is named in the Decision Brief'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-07T20:37:26Z
created_at: 2026-10-07T20:34:36Z
user_request: UR-139
addendum_to: REQ-641
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
write_set: [skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md]
claimed_at: 2026-10-07T20:35:54Z
dispatch_at: 2026-10-07T20:38:36Z
builder_handback_at: 2026-10-07T20:40:17Z
integration_at: 2026-10-07T20:40:57Z
kb_status: pending
review_at: 2026-10-07T20:47:20Z
commit: 0c41de5528054ab5e8c070e85ba1dd65badfd3fc
heavy_verified_at: 2026-10-07T20:47:52Z
heavy_verified_revision: 0c41de5528054ab5e8c070e85ba1dd65badfd3fc
completed_at: 2026-10-07T20:48:12Z
release_at: 2026-10-07T20:48:12Z
---
# Addendum: Close the Set-Aside Gaps in the Wave-End Decision
## What
Close the two set-aside gaps that the review of REQ-641 (wave-end consistency check) left report-only as findings F2 and F6. This addendum extends REQ-641: the integrator brief carries the set-aside members, work.md Step 7 reads that list when it decides "last", and the case it still cannot cover (F2) is named in the Decision Brief instead of being built for.
## Why
The wave-end restatement sweep runs inside the review of the wave's last successful integration. work.md Step 7 says the orchestrator decides "last" from the archive plus the members it set aside.
- **F6.** Under delegated integration the integrator did not set anything aside. A set-aside REQ keeps its claim in `do-work/working/`, so to the integrator it looks like a member still being built. The integrator cannot tell that its REQ is the last successful integration, and the sweep is skipped with no record.
- **F2.** When a member is set aside after the wave's last review already ran, no review received the fact, so the sweep did not run for that wave.
## Prior Implementation
REQ-641 (wave-end consistency check) shipped in 0.305.76. Final commit `78037a3e` (re-merge of builder branch `worktree-agent-REQ-641-wave-end-consistency-check` after the review fixes in `f5f95c65`); builder commits `39b21c30`, `0b306f69`, `a09d4e87`, first merge `f7416a80`. Archived file: `do-work/archive/UR-138/REQ-641-wave-end-consistency-check.md`. Review: `do-work/runs/work-2026-10-07-191445/REQ-641-review.md`.

Key files and what each now says:
- `skills/do-work/actions/review-work.md` Step 6 **Restatement Sweep**: step 1 gains "At a wave end, the trigger set also inherits." When the orchestrator passes that this REQ is its wave's last successful integration, with the earlier members' REQ ids, the trigger set also includes every element those members recorded on the **Restatement sweep:** line of their archived `## Review`. A member with no such line is named "unread for the sweep". Step 4 skips only when the whole trigger set is empty. The Append to REQ File template gains the `**Restatement sweep:** redefined <elements> | nothing redefined` line, and the Verification Checklist item points at it.
- `skills/do-work/actions/work.md` Step 7: **Restatement sweep (MUST)** runs the sweep when this REQ "is its wave's last successful integration (passed below)". **How to run it** passes that fact and the earlier members' ids "when every other member of this REQ's wave is already archived or set aside", and says: "The run manifest's rows name the wave's members, but they do not track who is done: the archive shows which members finalized, and the orchestrator knows which it set aside. A run with no manifest has no wave." This is the wave-end clause this REQ changes.
- `skills/do-work/actions/work-reference.md`: the Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick" gains one sentence naming the wave-end sweep as the semantic half no merge provides; **Folder Structure** names the flat archive of an open UR.
- `skills/do-work/docs/work-guide.md`: one user-facing sentence that the last review in a wave also catches meaning collisions.

Pattern used: the orchestrator decides "last" and passes the fact; the reviewer never derives wave state itself.
## Detailed Requirements
Condition-keyed, no new mechanism:
1. `skills/do-work/actions/work-reference.md`, the "Delegated integration — the coordinator shape" paragraph, item (e) the integrator brief: the brief also names every wave member the coordinator has set aside, because a set-aside member keeps its claim and is otherwise indistinguishable from one still building. One sentence.
2. `skills/do-work/actions/work.md` Step 7, the wave-end clause: the orchestrator decides "last successful integration" from the archive plus the set-aside list it holds, which under delegated integration is the list the integrator brief carries. One clause.
3. F2, say it rather than build for it: in the same Step 7 clause, when a member is set aside after the wave's last review already ran, the wave-end sweep did not run for that wave. The orchestrator names that in the run's Decision Brief under HANDLED (`work-reference.md` → **Decision Brief (hand-back format)**) so the gap is visible, and no review is re-run. One sentence.
4. Sweep every restatement of the wave-end rule and change none unless it contradicts the new sentences: `review-work.md` Step 6 **Restatement Sweep**, the `work-reference.md` Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick", and `docs/work-guide.md`.
## Constraints
- Prose-only, two shipped files (`work-reference.md`, `work.md`). No new mechanism.
- Keep every existing heading unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings.
- Read `_dev/primes/prime-action-files.md` before editing.
## Builder Guidance
Certainty is high: the source names each insertion point and caps each change at one sentence or one clause. Latitude is limited to wording.

Relation to REQ-641: extending. Requirement 3 also narrows REQ-641's clause "a set-aside member does not skip the check": for a member set aside after the last review ran, the sweep is skipped, and the change makes that skip visible instead of closing it. REQ-641's discovered task suggested "passing the fact to the last reviewed member when a later one is set aside"; the source chose to name the gap and re-run no review.
## Red-Green Proof
**RED prompt/case:** Today a delegated integrator reading its brief has no way to tell a set-aside sibling from one still building, so the condition "every other member is already finalized or set aside" cannot be evaluated and the sweep is silently skipped.
**Why RED now:** The integrator brief carries the wave membership but not which members were set aside, and a set-aside REQ keeps its claim in `do-work/working/` like an in-flight one.
**GREEN when:** The brief carries the set-aside list, Step 7 reads it, and the F2 case is named in the Decision Brief.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read plus `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 6295 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: family `alternate-writer-contract-drift` (requirement 4 sweeps restatements of a changed rule), and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-139/input.md` for complete verbatim input. Source facts: REQ-641's review re-check (`do-work/runs/work-2026-10-07-191445/REQ-641-review.md`, F6) and its discovered task "(impact-rule-change, F2 and F6) Close the set-aside gaps in the wave-end decision". No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the crew rules, `_dev/primes/prime-action-files.md` and the alternate-writer-contract-drift bullets. Two in-place prose edits, no new mechanism: one clause in the integrator-brief list of work-reference.md **Delegated integration — the coordinator shape**, and in work.md Step 7 **How to run it** the archive-plus-set-aside-list decision plus one F2 sentence pointing at Decision Brief → HANDLED. Keep the decision out of **Restatement sweep (MUST)**, which says "passed below". Then grep every restatement. (from the builder hand-back)
- [x] **[APPLY]:** Both edits made as planned, one commit `22c2268e` on `worktree-agent-REQ-642-wave-end-set-aside-gaps`; only the two `write_set` files touched. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat`: work-reference.md 2 +-, work.md 2 +-. `git diff --check` clean. `shipped-package-reference-contract.sh` PASS 1.0s; `contract-regressions.sh` PASS 23.1s. Checked work-reference.md line 470 and work.md line 370 by cold read, and every restatement hit for set aside, last successful, wave membership, integrator brief, wave end, HANDLED and Decision Brief (accounting in the hand-back). (from the builder hand-back)
*Source: the maintainer session's capture source for the REQ-641 follow-up, 2026-10-07.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names both files, the exact paragraph and clause to change, and caps each change at one sentence or one clause. Its restatement sweep names its three targets, so no exploration is needed to find where the change goes. The original REQ-641 (wave-end consistency check) was read from `do-work/archive/UR-138/`; its Implementation Summary, Decisions D-10 to D-15 and the F2/F6 discovered task match this REQ's `## Prior Implementation`.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)

**What was done:** In work-reference.md **Delegated integration — the coordinator shape**, the integrator brief's content list gains every wave member the coordinator has set aside, with the reason in the same clause: a set-aside member keeps its claim in `do-work/working/`, so without the list the integrator cannot tell it from one still building. In work.md Step 7 **How to run it**, the orchestrator now decides "last successful integration" from the archive plus the set-aside list it holds, which under delegated integration is the list the integrator brief carries. One new sentence covers F2: a member set aside after the wave's last review already ran means the sweep did not run for that wave, and the orchestrator names that gap under HANDLED in the run's Decision Brief as its call not to re-run a review. Review F1 fix (`aeb3f98f`): the work-reference.md **Decision Brief (hand-back format)** HANDLED bullet also lists any run-level call the orchestrator states itself, such as that skipped sweep, and the block is no longer omitted while one exists. No heading or bold label changed, no integration seams. Builder commits `22c2268e` and `aeb3f98f`; first merge `03cdf347`, re-merge `0c41de55`; cumulative range `3f27039f..0c41de55`.

## Qualification

**Diff range:** 3f27039f..03cdf347 for the qualify record (builder commit 22c2268e, merge 03cdf347); cumulative range after the review fix 3f27039f..0c41de55 (fix aeb3f98f, re-merge 0c41de55).
**Gate records:** qualify satisfied. Route A has no `## Scope`, so there is no scope-drift record; the two changed files equal the captured `write_set`.
**Warnings judged:** none raised. Read the diff against each requirement: requirement 1 is the added clause in the integrator-brief list; requirement 2 is the replaced clause in work.md Step 7 **How to run it**; requirement 3 is the new sentence after "A run with no manifest has no wave."; requirement 4 is the builder's restatement accounting in the hand-back (no other file contradicts the new text). Each change is one clause or one sentence, as the REQ caps it.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at 0c41de55 from the detached checkout `.git/work-run-2026-10-07-191445/drain-head`, 20:44:14Z to 20:46:25Z, load under 4, no other gate running; focused probe `do-work/runs/work-2026-10-07-191445/REQ-642-probe.sh` (`shipped-package-reference-contract.sh`, `contract-regressions.sh`) run by the advance test gate.
**Result:** ✓ All passing. Gate exit 0, gate wall 131s, slowest files 20.95s and 20.74s under the 30s budget. Advance test gate: run-blocked-check and green-gate satisfied. Builder runs: shipped-package-reference-contract.sh 1.0s, contract-regressions.sh 23.1s (first commit) and 20.9s (F1 fix).

The change is non-behavioral prose, so there is no red-green validation. The `## Red-Green Proof` GREEN condition (the brief carries the set-aside list, Step 7 reads it, the F2 case is named in the Decision Brief) was checked by a prose read of the cumulative diff.

**Heavy verification plan:**
- Range: 3f27039f0af0f0544d1e74336aa15ab3f1a50f51..0c41de5528054ab5e8c070e85ba1dd65badfd3fc
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/work-reference.md and skills/do-work/actions/work.md matched subtree skills

*Verified by work action*

## Review

**Overall: 92%** | 2026-10-07T20:47:20Z

| Dimension | Score |
|-----------|-------|
| Requirements | 97% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Reviewed range: 3f27039f..0c41de55 (builder 22c2268e, fix aeb3f98f).

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 (fixed in aeb3f98f). `work-reference.md:871` HANDLED bullet defined HANDLED only as REQ `D-NN` decisions and could omit the block, dropping the `work.md:370` late-set-aside note; the list clause and the omit condition now cover orchestrator-stated run-level calls. — impact-rule-change → report only

**Minor findings:** F2. `work.md:370` "the set-aside list it holds" has no durable home in the non-delegated case (session context only; manifest does not track it), and an integrator setting its own REQ aside is not said to feed later briefs. — impact-rule-change → report only. F3 (nit). `work.md:370` **How to run it** paragraph carries three rules in one paragraph. — impact-negligible → report only
**Acceptance:** Pass — both edits and the F1 fix match UR-139's four requirements; reference-contract test re-run PASS at 0c41de55.
**Restatement sweep:** redefined the source of the "last successful integration" decision (archive + held set-aside list / integrator brief), the integrator brief's contents, the late-set-aside (F2) outcome (skipped, named under HANDLED, no re-run), and HANDLED's content and omit rule (adds orchestrator-stated run-level calls)
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

## Decisions

Builder decisions, in full in the hand-back `do-work/runs/work-2026-10-07-191445/REQ-642-handback.md` § Decisions and § Addendum:
- **D-01 (DECIDE & STATE), superseded by D-04:** the HANDLED bullet needed no edit, because the late set-aside note is a reversible call made without asking.
- **D-02 (DECIDE & STATE):** the F2 sentence says "no review was passed the wave-end fact, so the sweep did not run for that wave", making the cause explicit and keeping the REQ's trigger wording.
- **D-03 (DECIDE & STATE):** work.md cites **Delegated integration — the coordinator shape** so a Step 7 reader can find who writes the brief.
- **D-04 (DECIDE & STATE):** review F1 overturned D-01. The HANDLED bullet's omit rule made it exclusive, so the bullet now also lists run-level calls the orchestrator states itself and is not omitted while one exists (`aeb3f98f`).

Orchestrator decisions:
- **D-05 (DECIDE & STATE): Route A.** The REQ names both files and every insertion point and caps each change at one sentence or clause, so there was no Exploration, Scope or Pre-Flight. Value: the build took two minutes. Risk: no Scope comparison; the two files equal the captured `write_set`.
- **D-06 (DECIDE & STATE): review F1 fixed before the gate,** on the builder branch, re-merged with the same `<pre>` as `0c41de55`; the reviewer re-checked the fix delta and found nothing new. Requirement 4 asks to change a restatement that contradicts the new text, and work-reference.md is in the write set. Value: the F2 note cannot be dropped from a brief. Risk: none; one re-merge, one gate.
- **D-07 (DECIDE & STATE): review F2 and F3 stay report only.** F2 asks for a durable home for the non-delegated set-aside list, which is new mechanism the REQ rules out; in the serial loop the orchestrator holds its own set-aside findings. F3 is readability only.
- **D-08 (DECIDE & STATE): dispatch instant from the branch reflog.** The builder call returned only after the build finished, so the instant is the branch creation, 2026-10-07T20:38:36Z, the first trace of the accepted dispatch. The builder-work timing event was recorded right after the hand-back landed.
- **D-09 (DECIDE & STATE): no wave-end fact passed.** This REQ is a one-member wave in its run; the earlier manifest rows are the finished UR-138 wave. The review recorded its own **Restatement sweep:** line and inherited nothing.

## Discovered Tasks

- impact-rule-change: the orchestrator's set-aside list has no durable home outside delegated integration, so a context loss between REQs makes the "last successful integration" decision unprovable; and nothing says a coordinator copies an integrator's own set-aside into later briefs (review F2 plus the builder's second discovered task). → report only
- impact-negligible: `work.md` Step 7 **How to run it** now carries dispatch arguments, the wave-end decision and the late set-aside rule in one paragraph (review F3). → report only

## Lessons Learned

**What worked:** Launching the reviewer right after qualify put its Important finding (F1) ahead of the gate, so the fix cost one re-merge and one gate run. Routing the fix through the idle builder kept the range one cumulative range.
**What didn't:** The builder saw the HANDLED bullet conflict and judged it "descriptive, not exclusive" (D-01). It missed that the bullet's omit rule made it exclusive in practice, so the new note would vanish in exactly the run where it was the only entry.
**Worth knowing:** When a rule sends a new kind of entry into an existing output block, read that block's own content rule and its omit or empty condition. A block that is omitted when its usual source is empty drops the new entry.

## Orientation

Now a delegated integrator can tell which wave members were set aside, because its brief lists them, so the last successful integration still triggers the wave-end restatement check. When a member is set aside after the wave's last review already ran, the run's Decision Brief says under HANDLED that the check did not run, and no review is re-run. This lives in `actions/work.md` Step 7 **How to run it** and `actions/work-reference.md` → **Delegated integration — the coordinator shape** and **Decision Brief (hand-back format)**. `_dev/primes/prime-action-files.md` is not stale: no path it names moved or was removed.

## Heavy Verification Plan

- Base revision: 3f27039f0af0f0544d1e74336aa15ab3f1a50f51
- Target revision: 0c41de5528054ab5e8c070e85ba1dd65badfd3fc
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/work-reference.md and skills/do-work/actions/work.md matched subtree skills

## Heavy Verification Result

- Target revision: 0c41de5528054ab5e8c070e85ba1dd65badfd3fc
- Execution revision: 0c41de5528054ab5e8c070e85ba1dd65badfd3fc (detached checkout `.git/work-run-2026-10-07-191445/drain-head`, `QUEUE_KANBAN_BROWSER` set)
- staged-skills: exit 0, executed, 39s, no `HEAVY-RUN-LANE-SKIPPED` finding

## Timing

Observed 2026-10-07T20:38:36Z to 2026-10-07T20:40:41Z: 2m 05s total, 2m 05s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 2m 05s | 1 |

Slowest stage: builder-work / builder agent in worktree, 2m 05s, outcome success.
