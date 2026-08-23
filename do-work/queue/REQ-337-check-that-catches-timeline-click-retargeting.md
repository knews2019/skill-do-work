---
id: REQ-337
title: "A check that can catch Timeline click retargeting"
status: pending
created_at: 2026-08-23T18:30:26Z
user_request: UR-067
addendum_to: REQ-324
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-336]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-336, REQ-338]
batch: timeline-click-regression
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# A Check That Can Catch Timeline Click Retargeting

## What

Give the Timeline probe lane a check that fails on the pre-REQ-336 behaviour (pointer capture on
pointerdown retargeting the synthesized click) and passes after REQ-336's fix. The existing
lock-in test passes against the broken build, so the lane currently cannot catch this class of
regression.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving`
(`timeline_browser_probe_test.go:626`) is REQ-324's lock-in test, and it passed while every real
mouse click in the Timeline was broken. Reason: it dispatches synthetic PointerEvents with
`pointerId: 1`, an id the engine does not know, so `setPointerCapture` throws and capture is
never established — the probe never exercises the code path that breaks real clicks. A lock-in
test that cannot catch the regression it pins is not doing its job.

## Detailed Requirements

- The check must fail on the broken behaviour (capture established on pointerdown, click
  retargeted to `#timeline-scroll`) and pass on the fixed behaviour (REQ-336).
- Either drive real input, or assert the structural property directly (where the capture call
  sits relative to the engage threshold `TIMELINE_PAN_THRESHOLD_PX`).
- **Mutation-test whichever you pick: a guard no mutation can break is dead code.** At minimum,
  reintroducing capture-on-pointerdown (reverting REQ-336's fix) must turn the check red.

## Builder Guidance

Certainty: Firm on the goal, open on the mechanism (real input vs structural assertion — the
user left the choice to the builder). Ordering was decided at capture and user-confirmed at
verify (2026-08-23): this REQ depends on REQ-336 so the committed suite never carries a red
check; prove the RED side by mutation (reintroduce the broken capture placement locally and
watch the check fail), not by committing against broken HEAD.
`tdd: false` because the deliverable IS the check — the mutation evidence is its proof, and the
work loop's test-first gate has no separate implementation to precede. Cautionary precedent from
REQ-333's UNIFY notes: its structural assertion first matched the `typeof` guard beside the
capture call instead of the call itself, and only mutation testing exposed that — expect the
same trap here.

## Red-Green Proof

**RED prompt/case:** Check out the pre-REQ-336 behaviour (or mutate the fix back to
capture-on-pointerdown) and run the new check: it must fail. On current HEAD before this REQ,
`TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving` passes against that same broken
behaviour — that gap is what this REQ closes.
**Why RED now:** The lane's only click/pan lock-in test never establishes pointer capture
(synthetic `pointerId` throws), so click retargeting is invisible to it.
**GREEN when:** The new check passes on post-REQ-336 HEAD, demonstrably fails when
capture-on-pointerdown is reintroduced (mutation evidence recorded in the REQ), and
`bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed (goal and both accepted mechanisms stated verbatim in the input)

## Prior Implementation

REQ-324 ("Give the timeline drag a movement threshold", archived, completed, commit 3486ab2)
added `TIMELINE_PAN_THRESHOLD_PX` (4px) so a press that does not move is still a click, and
pinned it with `TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving`
(`timeline_browser_probe_test.go:626`). REQ-333 (commit 36c4518) later noted in its UNIFY that a
synthetic `pointerId` cannot be captured in this lane, so it asserted the capture call
structurally and drove `lostpointercapture` directly — the workaround whose blind spot this REQ
closes.

## Dependencies

`depends_on: [REQ-336]` — the fix lands first; this check then pins it, proving RED by mutation.

## Notes

The input tagged this `[impact-high]`, which is not a canonical `impact:` value; normalize-and-warn
fired and the user confirmed `impact-user-visible` during capture (default tier, so the title
carries no tag).

---
*Source: UR-067 — see `do-work/user-requests/UR-067/input.md` for complete verbatim input.*
