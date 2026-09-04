---
id: UR-114
title: 'Delete the pending-heavy-testing status: held requests stay claimed and finalize on green'
created_at: 2026-09-04T22:52:00Z
requests: [REQ-570, REQ-571]
word_count: 211
---

# Delete the pending-heavy-testing Status: Held Requests Stay Claimed and Finalize on Green

## Summary

A green heavy-lane result currently returns a request to plain `pending`. The dependency rule does not count a `pending` request as landed source, so merged and verified requests wait in series for one review each, and the board shows the same word for "never started" and "merged, verified, awaiting review". The live selector at 2026-09-04 22:3x UTC showed it: REQ-504 selected with `resume_phase: review`, REQ-505, REQ-506 and REQ-507 excluded as `DEPENDENCIES-UNMET`, all four commits ancestors of `main`. The user asked why held requests are not simply left `claimed`, since they are being worked on, and approved that design. The heavy hold becomes a phase of a claimed request, marked by the `## Heavy Verification Plan` section and a `commit:` on `main`, exactly as `actions/work.md` line 93 already says phases are tracked. Review runs before the hold; a green drain finalizes in the same session; a red drain enters ordinary remediation. The status value and every reader of it are deleted. The finding and the proposal are recorded in `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html` section 03 (commits 8f462956 and eabc2984).

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-570 | Delete the `pending-heavy-testing` status from the core skill: held requests stay claimed, review precedes the hold, a green drain finalizes, recover routes a held claimed request to the drain, the readiness rule reads `claimed` plus `commit:` |
| REQ-571 | Remove the board's dead `pending-heavy-testing` reader case once the core skill no longer writes it |

## Batch Constraints

- Land after the advance chain (REQ-504 to REQ-507) has closed through the existing path, so no compatibility read of the old status is needed. REQ-570 depends on REQ-507 (finalization behind `advance`), which the drain will call.
- Delete before you add: no new status value, no new count in the checkpoint, no manual lane path. A held request counts as in progress; that is accurate.
- A skipped lane leaves the request claimed with its typed finding; the next drain retries; a persistent skip is named in the exit summary for a human, never a hand edit.
- One session owns a held request from claim to finalization. If the session dies during the hold, the next session's recover step routes the request to its own exhaustion drain by reading the plan section, never by clock age.
- Judgment stays prose; the CLI emits typed findings. The floor agent must still complete a run using only `advance` output plus the remaining prose.

## Full Verbatim Input

> ```
> how do you want to implement it, should we capture a req for it?
> 
> capture it and fix it
> 
> [The proposal the user approved with "capture it", drafted by the agent earlier in the same conversation:]
> Delete the pending-heavy-testing status. A held request stays claimed in do-work/working/; its phase is already marked by the Heavy Verification Plan section and a commit on main. Move Step 7 review and Step 7.5 lessons ahead of the hold, so a request is held only after fast tests, qualification and review pass. At queue exhaustion the same session runs the lanes, then finalizes each green request through Step 8 and 9 in the same turn and enters ordinary remediation for each red one. The readiness rule accepts claimed plus a nonblank commit in place of pending-heavy-testing plus a commit. Recover recognizes a claimed request that has a Heavy Verification Plan and no live writer and routes it to this session's drain instead of treating it as interrupted build work. Delete the heavy-testing answer mode and its refusal codes, resume_phase and matchingHeavyReviewPhase in the selector, clarify Step 2.5, and every remaining reader of the status. A skipped lane leaves the request claimed with its finding; the next drain retries. depends_on REQ-507. batch orchestrator-simplification. impact-rule-change, tdd, maintenance.
> ```

---
*Captured: 2026-09-04T22:52:00Z*
