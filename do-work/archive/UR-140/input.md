---
id: UR-140
title: 'Board shows an earmarked pending REQ as Ready, and capture does not separate a session earmark from work that needs the user present'
created_at: 2026-10-07T23:21:03Z
requests: [REQ-643, REQ-644]
word_count: 720
---
# Earmarked REQs Read as Ready on the Board; Capture Lacks the Session-vs-Operator Distinction

## Summary
The maintainer ran `do-work-toolbox validate-feedback` on their own board feedback and both findings were accepted, then asked to capture and run them. Finding 1: a `pending` REQ with a non-empty `assigned_to` and no unmet dependencies lands under Pending → Ready on the Kanban board (and in the "ready to work" CLI counts) although the run scan skips it. Finding 2: capture's Earmark assessment does not say when to use `assigned_to` (routing between sessions) versus `status: blocked` with `blocked_by` naming the person (the user must be present as operator). The triage chose the third Pending sub-group ("Earmarked") over routing to Needs Input · Blocked; the reasons are in the verbatim input and in REQ-643's body.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-643 | Board shows an earmarked pending REQ under Pending → Earmarked instead of Ready, and the ready counts agree |
| REQ-644 | [impact-rule-change] Capture distinguishes a session earmark from work that needs the user as operator |

## Batch Constraints
- The badge stays verbatim-read and its value is untouched; this is about placement only.
- `model.go`'s `AssignedTo` comment and the `assigned_to` schema line in `skills/do-work/actions/work-reference.md` (today: "no column logic") change in the same commit, as the lock-step rule requires.
- No dependency between the two REQs: REQ-644's sentence is true before and after REQ-643 ships. They touch disjoint files and may run in parallel.
- Read `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-releases.md` before REQ-643; `_dev/primes/prime-action-files.md` before REQ-644.

## Full Verbatim Input
> ```
> capture then, then run them
> 
> [Capture source: the two findings accepted by `do-work-toolbox validate-feedback` in this session, with the maintainer's original feedback verbatim, then the triage Evidence and Surface-cost for each.]
> 
> --- Original feedback (verbatim) ---
> 
> Board: a pending REQ earmarked for the user reads as "Ready", which misleads.
> 
> Today assigned_to is an advisory frontmatter field. The run scan skips an
> assigned REQ as a courtesy, but the Kanban board (do-work-board/tools/
> queue-kanban/model.go, bucketColumns) places the card by status and
> dependencies only, so an earmarked REQ with no unmet dependencies lands
> under Pending → Ready with a small "ASSIGNED <session>" badge. "Ready"
> tells the user the run will take it next, when the run will skip it. I
> went looking for my two earmarked REQs in Needs Input · Blocked, where
> "waiting on me" belongs, and found nothing.
> 
> Two changes, either or both:
> 
> 1. Board routing. A pending REQ with a non-empty assigned_to should not
>    count as Ready. Either put it in a third Pending sub-group
>    ("Earmarked N", after Ready and Waiting) or route it to Needs Input ·
>    Blocked with the badge text "assigned to <session>". Keep the badge
>    verbatim-read; this is about placement, not the value. The schema
>    comment on assigned_to in actions/work-reference.md says "no column
>    logic", so that line and model.go change in the same commit, as the
>    comment already requires.
> 
> 2. Capture guidance. When an earmark means "the user has to be present
>    as operator" rather than "another session will take it", capture
>    should prefer status: blocked with blocked_by naming the person and
>    blocked_at, which already lands in Needs Input · Blocked and is
>    released by clarify. Add one sentence to the Earmark assessment in
>    actions/capture.md distinguishing the two cases: assigned_to is for
>    routing between sessions; blocked-on-a-person is for work that needs
>    the user in the loop.
> 
> Repro: a queue with REQ-A (pending, no depends_on, assigned_to:
> 'user-interactive'). Expected: not under Ready. Actual: under Ready
> with the badge. Workaround used today: flipped both REQs to blocked
> with blocked_by 'the maintainer as operator, in an interactive session'.
> 
> --- Triage: Finding 1 (Earmarked pending REQ shows as Ready on the board), severity unrated, source: maintainer feedback ---
> 
> Verdict: Accept.
> Evidence: skills/do-work-board/tools/queue-kanban/model.go:1719-1726 buckets on status and unmet dependencies only; AssignedTo is never read there. The PendingReady comment at model.go:382 promises "actionable now", while the run scan excludes assigned REQs (skills/do-work/actions/work-reference.md:511, exit-summary section 6, and the auto-wave predicate at line 460). web/board.css:1301 states the original decision: "it is somebody's intention, not a state of the work." Origin: archived REQ-097 (the advisory field, display-only parse) in do-work/archive/UR-018/.
> Reasoning: the mislead is wider than the card. The same bucket feeds the CLI summary line "ready to work" (main.go:123) and the open-work digest (open_work.go:44), so three surfaces call a skipped REQ ready. The documented decision was "no scheduling, not a lock"; a Pending sub-group changes placement only and keeps that intact. The third sub-group ("Earmarked N", after Ready and Waiting) fits better than Needs Input · Blocked: that column is the operator-actionable inbox (board.md:92, model.go:1087), and a REQ reserved for another session is not actionable by the operator; routing it there would also break the status-based isNeedsInputOrBlockedStatus contract and the open-work counts. Frontend cost is one more makePendingGroup call (board-cards.js:497-498) plus a pendingEarmarked payload list in generate.go:146. Trap: TestAssignedToNeverAffectsColumnPlacement (model_test.go:1020) only asserts the two tickets share a status, so it stays green under either option and must be replaced by a test that pins the new placement. Also update board-guide.md:18 and the PendingReady struct comment, and bump the board version per the kanban prime.
> Surface-cost: N/A. Display correction, no guard, fallback, or warning added.
> 
> --- Triage: Finding 2 (Capture should distinguish "another session takes it" from "the user must be present"), severity unrated, source: maintainer feedback ---
> 
> Verdict: Accept.
> Evidence: skills/do-work/actions/capture.md:107 (Earmark assessment) never mentions blocked. Line 108 (External-condition assessment) lists three look-alikes (depends_on, pending-answers, Answerer:) but not the earmark. The release path exists: clarify.md:13 and Step 5.5 confirm human-confirmable blocked conditions, and work.md:97 keeps a blocked REQ out of the scan.
> Reasoning: both shapes already have correct homes; only the capture-time routing sentence is missing, and the user's workaround is exactly the documented blocked shape. Adding the earmark as a fourth look-alike on line 108 and one sentence on line 107 closes it without new machinery.
> Surface-cost: N/A. One disambiguation sentence in existing guidance.
> ```
