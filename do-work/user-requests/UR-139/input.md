---
id: UR-139
title: 'Addendum to REQ-641: close the set-aside gaps in the wave-end decision (review findings F2 and F6)'
created_at: 2026-10-07T20:34:36Z
requests: [REQ-642]
word_count: 447
---
# Addendum to REQ-641: Close the Set-Aside Gaps in the Wave-End Decision

## Summary
The maintainer's session asked for one follow-up REQ to REQ-641 (wave-end consistency check, released 0.305.76, archived in `do-work/archive/UR-138/`). It closes the two set-aside gaps REQ-641's review left report-only as findings F2 and F6. REQ-641 is archived and immutable, so this is a new UR and a new queued REQ with `addendum_to: REQ-641`. Prose-only, two shipped files.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-642 | [impact-rule-change] Addendum: the wave-end decision reads the set-aside list, and a late set-aside is named in the Decision Brief |

## Full Verbatim Input
> ```
> Follow-up to REQ-641 (wave-end consistency check, released 0.305.76, archived in do-work/archive/UR-138/): close the set-aside gaps its review left report-only as findings F2 and F6. One REQ, addendum_to REQ-641, prose-only, two shipped files.
> 
> The gap. The wave-end restatement sweep runs inside the review of the wave's last successful integration. work.md Step 7 says the orchestrator decides "last" from the archive plus the members it set aside. Under delegated integration (work-reference.md, Fan-Out Dispatch, "Delegated integration — the coordinator shape") the integrator did not set anything aside: a set-aside REQ keeps its claim in do-work/working/, so to the integrator it looks like a member still being built, the integrator cannot tell that its REQ is the last successful integration, and the sweep is skipped with no record (F6). Separately, when a member is set aside after the wave's last review already ran, no review received the fact, so the sweep did not run for that wave (F2).
> 
> Changes, condition-keyed, no new mechanism:
> 1. skills/do-work/actions/work-reference.md, the "Delegated integration — the coordinator shape" paragraph, item (e) the integrator brief: the brief also names every wave member the coordinator has set aside, because a set-aside member keeps its claim and is otherwise indistinguishable from one still building. One sentence.
> 2. skills/do-work/actions/work.md Step 7, the wave-end clause: the orchestrator decides "last successful integration" from the archive plus the set-aside list it holds, which under delegated integration is the list the integrator brief carries. One clause.
> 3. F2, say it rather than build for it: in the same Step 7 clause, when a member is set aside after the wave's last review already ran, the wave-end sweep did not run for that wave; the orchestrator names that in the run's Decision Brief under HANDLED (work-reference.md, Decision Brief) so the gap is visible, and no review is re-run. One sentence.
> 4. Sweep every restatement of the wave-end rule (review-work.md Step 6 Restatement Sweep, work-reference.md Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick", docs/work-guide.md) and change none unless it contradicts the new sentences. Keep every existing heading unchanged; _dev/tests/shipped-package-reference-contract.sh pins citation strings.
> 
> Red-Green: RED: today a delegated integrator reading its brief has no way to tell a set-aside sibling from one still building, so the condition "every other member is already finalized or set aside" cannot be evaluated and the sweep is silently skipped. GREEN: the brief carries the set-aside list, Step 7 reads it, and the F2 case is named in the Decision Brief. Validation: prose read plus bash _dev/tests/shipped-package-reference-contract.sh and bash _dev/tests/contract-regressions.sh.
> 
> Verified facts: REQ-641's review re-check (do-work/runs/work-2026-10-07-191445/REQ-641-review.md, F6) and its discovered task "(impact-rule-change, F2 and F6) Close the set-aside gaps in the wave-end decision". prime_files: _dev/primes/prime-action-files.md. impact-rule-change, effort-mechanical, tdd false.
> ```