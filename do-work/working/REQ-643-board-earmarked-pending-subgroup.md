---
id: REQ-643
title: 'Board shows an earmarked pending REQ under Pending → Earmarked instead of Ready, and the ready counts agree'
status: claimed
created_at: 2026-10-07T23:21:03Z
user_request: UR-140
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-644]
batch: earmark-placement
write_set: [skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/model_test.go, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/open_work.go, skills/do-work-board/tools/queue-kanban/main.go, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work-board/docs/board-guide.md, skills/do-work-board/actions/board.md, skills/do-work/actions/work-reference.md]
claimed_at: 2026-10-07T23:26:38Z
---
# Board Shows an Earmarked Pending REQ Under Pending → Earmarked Instead of Ready
## What
A `pending` REQ with a non-empty `assigned_to` and no unmet dependencies must not count as Ready anywhere the board reports readiness. It lands in a third Pending sub-group, "Earmarked", after Ready and Waiting, and the CLI summary and open-work digest stop counting it as ready.
## Why
"Ready" tells the user the run will take the REQ next. The run scan skips an assigned REQ as a courtesy (`skills/do-work/actions/work-reference.md` exit-summary section 6 and the auto-wave predicate), so the board and the pipeline disagree on the one thing the Ready label promises. The maintainer looked for two earmarked REQs under Needs Input · Blocked and found nothing.
## Verified Facts (from triage)
- `skills/do-work-board/tools/queue-kanban/model.go` `bucketColumns` (around line 1719) buckets by status and `UnmetDependencies` only; `AssignedTo` is never read there. The `PendingReady` comment (line 382) promises "actionable now".
- The same bucket feeds `main.go:123` ("ready to work") and `open_work.go:44` (the digest's `N ready`), so three surfaces call a skipped REQ ready.
- The documented decision (`model.go` `AssignedTo` comment, `work-reference.md` schema line "no column logic and no scheduling", `web/board.css` ".badge-assigned" comment "somebody's intention, not a state of the work") was about never making the field a gate or a scheduler. A sub-group gates nothing.
- Needs Input · Blocked is defined as the operator-actionable inbox (`skills/do-work-board/actions/board.md:92`, `isNeedsInputOrBlockedStatus` in `model.go`). A REQ reserved for another session is not actionable by this operator, and that bucket is status-keyed, so routing a `pending` ticket there would break its contract and the open-work counts.
- `TestAssignedToNeverAffectsColumnPlacement` (`model_test.go:1020`) only asserts that the assigned and unassigned tickets share a `Status`. It would stay green under either placement option, so it pins nothing about columns.
- The frontend already renders Pending sub-groups generically through `makePendingGroup` (`web/board-cards.js:497-515`); `fillPendingColumn` collapses to a flat list only when Waiting is empty.
## Detailed Requirements
1. `model.go`: add a `PendingEarmarked` bucket to `BoardColumns` for tickets with `Status == "pending"`, no unmet dependencies, and non-empty `AssignedTo`. Such a ticket leaves `PendingReady`; `Pending` stays the union (Ready, Waiting, Earmarked) and the priority sort applies to the new group as it does to the other two. A `pending` ticket with unmet dependencies stays in Waiting whether or not it is assigned (Waiting already says it is not actionable). Update the `BoardColumns` field comments and the `AssignedTo` doc comment (today: "The board never buckets, orders, or schedules on it") so they describe the new placement truthfully: the board groups on it, and still never orders across groups, schedules, or gates on it.
2. `skills/do-work/actions/work-reference.md`, the `assigned_to:` schema line: replace "with **no column logic and no scheduling**" with the true statement (the board shows an assigned pending REQ in the Pending column's Earmarked group instead of Ready; no scheduling), in the same commit as the `model.go` change. Keep every existing heading unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings.
3. `generate.go`: thread `pendingEarmarked` through the payload struct and its single copy site (the `PendingReady`/`PendingWaiting` pattern at lines 142-146 and 776).
4. `web/board-cards.js`: `fillPendingColumn` renders a third `makePendingGroup("Earmarked", …)` after Waiting when the group is non-empty; the flat-list shortcut applies only when both Waiting and Earmarked are empty. The badge keeps its current text and verbatim value; nothing about the badge changes.
5. `main.go` summary: "ready to work" counts `PendingReady` only; add an `earmarked` line beside "waiting on deps". `open_work.go`: the digest's pending breakdown reads `(N ready, N waiting, N earmarked)` and `countOpenWork` exposes the new count.
6. `model_test.go`: replace `TestAssignedToNeverAffectsColumnPlacement` with a test that writes a pending unassigned REQ and a pending assigned REQ with no dependencies through the normal board build and asserts the assigned one is in `PendingEarmarked` and not in `PendingReady`, the unassigned one is in `PendingReady`, and a pending assigned REQ with an unmet dependency is in `PendingWaiting`. Keep the existing verbatim-read and absent-reads-as-empty tests.
7. Docs: `skills/do-work-board/docs/board-guide.md:18` (Pending splits into Ready and Waiting) and `skills/do-work-board/actions/board.md:111-112` (digest and summary descriptions) name the third group.
8. Board version bump and changelog entry per `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-releases.md` (the release step of the run owns the mechanics).
## Constraints
- Placement only. The `assigned_to` value stays verbatim-read with no alias map and no case folding; the badge text and the drawer row are unchanged.
- No scheduling and no gating on `assigned_to` anywhere on the board. The one behavioral reader stays the work pipeline's scan.
- `model.go` and the `work-reference.md` schema line change in the same commit (lock-step rule, stated on both).
- `tdd: true`: the Go test in requirement 6 is written first and seen failing before the model change.
## Dependencies
None. REQ-644 (capture distinguishes a session earmark from operator-present work) is related and touches disjoint files; the two may run in parallel.
## Builder Guidance
Certainty is high on the decision (sub-group, not Needs Input routing) and on the surfaces to change; the triage listed them from the code. Latitude: the exact label text ("Earmarked") and the summary line wording. Do not add an `earmarked` filter, sort, or any new behavior beyond placement and counts.
## Red-Green Proof
**RED prompt/case:** A queue with one REQ: `status: pending`, no `depends_on`, `assigned_to: 'user-interactive'`. Build the board and run `queue-kanban summary`.
**Why RED now:** The card renders under Pending → Ready with the "assigned user-interactive" badge, and the summary says `ready to work : 1`, while `do-work run` would skip it.
**GREEN when:** The card renders under Pending → Earmarked (1), Ready shows 0, the summary says `ready to work : 0` and `earmarked : 1`, and the open-work digest reads `(0 ready, 0 waiting, 1 earmarked)`. The new Go test pins the placement.
**Validation:** User confirmed. The repro and expected/actual are the maintainer's own words in the verbatim input.
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 7636 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the index row governs queue-kanban model, parser, and UI changes, which this REQ makes; family `paired-predicate-drift` (two readers of one contract) is the shape of this bug.
- `_dev/primes/lessons-kanban-board.md` — 5912 tokens, `slugged: partial`; over budget. Matching reason: its owning prime governs the board tool this REQ edits.
## Full Context
See `do-work/user-requests/UR-140/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: maintainer feedback "Board: a pending REQ earmarked for the user reads as 'Ready', which misleads", change 1, accepted by validate-feedback.*
