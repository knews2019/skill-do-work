---
id: REQ-649
title: 'Board Earmarked badge tooltip names Needs input · Blocked as the home for operator-gated work'
status: claimed
created_at: 2026-10-08T15:50:17Z
user_request: UR-142
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-647, REQ-648]
batch: earmark-release-visibility
write_set: [skills/do-work-board/tools/queue-kanban/web/board-cards.js]
claimed_at: 2026-10-08T15:53:16Z
---
# Board Earmarked Badge Tooltip Names Needs Input · Blocked as the Home for Operator-Gated Work
## What
The `assigned` badge tooltip on a board card (`skills/do-work-board/tools/queue-kanban/web/board-cards.js:234-239`) ends "The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this." It says what the group does, not what it is not for. Add one sentence: a REQ waiting on the operator belongs under Needs input · Blocked (`status: blocked`), not here.
## Why
Consumer incident, Oct 8 14:23 UTC on 0.305.79: the user saw two REQs that needed their input under Pending → Earmarked and asked why they were not under Needs input · Blocked. The badge was the one place on the board that could have said so. Triage (validate-feedback, 2026-10-08, F5) accepted the tooltip sentence and found no group-header hint to extend.
## Verified Facts (from triage)
- `board-cards.js:229-241`: `assignedBadge.title` is built from four string literals ending `"Pending → Earmarked; it never reorders, blocks, or hides on this."`.
- `board-cards.js:509-524` (`makePendingGroup`): a Pending group header is a name and a count only; no group carries a hint or title attribute, so there is no header hint to extend and none is added.
- No Go or JavaScript test pins the tooltip text: `grep -rn "never reorders" skills/do-work-board/tools/queue-kanban` matches only `board-cards.js` (three other badges carry the shorter "Display only: the board never reorders, blocks, or hides on this." at lines 269, 299, 321; those are not the earmark badge and stay as they are).
- `generate_test.go` reads `web/board-cards.js` for class-name pins only; the badge class and text are unchanged by this REQ.
- 0.305.79 (REQ-643) last edited this tooltip to say the board groups on the field; this REQ only appends to it.
## Detailed Requirements
1. Extend the earmark badge tooltip's last sentence so the title reads "… The board only groups it under Pending → Earmarked; it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." Keep the badge text, class, placement, and every other sentence of the tooltip unchanged.
2. Add no group-header hint: `makePendingGroup` stays a name and a count.
3. No change to `model.go`, `generate.go`, the board data payload, `board.css`, or any other badge's tooltip.
4. Verify from the module root: `go build ./... && go vet ./... && go test ./... -count=1` in `skills/do-work-board/tools/queue-kanban` (the package tests read the web files); `gofmt -l .` prints nothing. The integrator's heavy lanes (`queue-kanban-javascript`, `queue-kanban-browser`) cover the rendered page.
5. Release per `_dev/primes/prime-releases.md` (the board is versioned with the skill; the integrator writes the CHANGELOG entry and bump).
## Constraints
- Text only. The field stays verbatim-read and display-only on the board: no column logic, no scheduling, no filter, no sort.
- Read `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` before editing. Never commit build outputs.
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-647 or REQ-648 into it.
## Dependencies
None. REQ-647 (clarify) and REQ-648 (schema line and work action) touch different modules; the three may run in parallel.
## Builder Guidance
Certainty is high: one appended string literal, quoted in the feedback and accepted as worded. No latitude beyond line wrapping of the string concatenation.
## Red-Green Proof
**RED prompt/case:** Build the board for a queue holding one `pending` REQ with `assigned_to: 'user-interactive'` and no `depends_on`; hover the card's `assigned` badge.
**Why RED now:** The title ends "it never reorders, blocks, or hides on this." and never says where operator-gated work belongs.
**GREEN when:** The title ends "… A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here." `grep -c "belongs under Needs input" skills/do-work-board/tools/queue-kanban/web/board-cards.js` prints 1, the package tests stay green, and the browser heavy lane renders the board.
**Validation:** User confirmed. The maintainer approved the triage (F5 accepted as worded) and said "capture and run".
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the board tool this REQ edits.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8110 tokens, over the budget; `slugged: partial`). Matching reason: its owning prime governs `web/board-cards.js`; REQ-643's entry there is this tooltip's own history.
## Full Context
See `do-work/user-requests/UR-142/input.md` for complete verbatim input (the consumer incident timeline, the board screenshot description, and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer feedback "the board's Earmarked badge says what the group does, not what it is not for" (Gap 3 / the ask §3), accepted by `do-work-toolbox validate-feedback` on 2026-10-08 as F5.*
