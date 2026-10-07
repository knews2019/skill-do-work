---
id: REQ-642
title: '[impact-rule-change] Addendum: the wave-end decision reads the set-aside list, and a late set-aside is named in the Decision Brief'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: the maintainer session's capture source for the REQ-641 follow-up, 2026-10-07.*
