---
id: REQ-568
title: 'Show recently touched REQs on the board regardless of status'
status: pending
created_at: 2026-09-04T17:54:05Z
user_request: UR-112
domain: ui-design
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Show Recently Touched REQs on the Board Regardless of Status

## What

Give the Kanban board one surface that answers "what changed on the queue in the last N hours, and why", listing every REQ whose newest lifecycle stamp falls inside the selected window, newest first, with the stamp and the transition it records. Status must not filter it: a REQ that was claimed, held for heavy testing, deferred, blocked, completed, cancelled, or failed inside the window all belong on it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

On the board generated 2026-09-04 16:56 UTC, the only claimed card (REQ-505, moving selection and claim behind advance) was 20 minutes old and the newest Recently done card (REQ-485, canonicalizing reservation marker filenames) was two hours old. The maintainer read that as a gap and asked where the work went. Git history showed the loop never paused: REQ-567 (repairing shipped lesson links to archived UR paths), REQ-503 (adding the read-only advance lifecycle command), and REQ-504 (collapsing Step 10 and crash recovery prose into recovery) were each claimed, built, merged, and held as `pending-heavy-testing` between 14:57 and 16:39 UTC. Those three sit under Pending, Waiting, so no "recent" surface on the board showed them. Answering the question needed `git log`.

## Context

- The board has three time surfaces today, none of which answers the question: the Recently done column (terminal states only, 24h / 48h / 7d window), the Timeline view (wait and work spans per REQ, window buttons for last day / 7 / 30 / 90 / all days), and the Calendar view (one entry per REQ on its claim or resolve day). The `open-work` terminal digest excludes finished REQs on purpose.
- Tickets already carry several lifecycle stamps the board parses: `created_at`, `claimed_at`, `completed_at`, `blocked_at`, the REQ-448 phase milestones (`planning_at`, `exploration_at`, and the rest), and `testing_updated_at`. The `pending-heavy-testing` hold itself currently writes no stamp: REQ-503's frontmatter after the hold carries `claimed_at` and the phase milestones but nothing recording when it was held. The builder decides whether the newest-stamp rule is enough or whether the hold transition must also stamp a time; the REQ intent is that the hold event is visible on this surface either way.
- "Newest stamp on a ticket" is the proposed key, and the transition it belongs to is what the row should name in words (claimed, held for heavy testing, blocked, completed, and so on). Which control hosts the surface (a new toolbar tab, a strip like Verify findings, or reuse of the existing Recently done window buttons) is the builder's call; the maintainer's words were "a 'recently touched' window".
- Release commits (for example 0.275.3 at 15:54 UTC in the same gap) are not REQs and are out of scope for this surface unless the builder finds them free to include.
- Board versioning, parser lock-step, and build outputs follow `_dev/primes/prime-kanban-board.md`.

## Red-Green Proof
**RED prompt/case:** Generate the board against this repository at a state where REQ-503, REQ-504, and REQ-567 are `pending-heavy-testing` with claims inside the last two hours and REQ-505 is `claimed`, then look for any surface listing what changed in the last 2 hours.
**Why RED now:** Recently done shows only REQ-485 and older terminal REQs. The three held REQs appear only under Pending, Waiting, with no time ordering and no transition name. The Timeline draws their open bars but not the hold. Nothing on the board says "REQ-504 was held for heavy testing at 16:38".
**GREEN when:** One board surface lists REQ-505, REQ-504, REQ-503, REQ-567, and REQ-485 newest first for a "last 24h" window, each row naming its newest stamp and the transition it records, and a Go test on the aggregation pins that a `pending-heavy-testing` REQ with a hold inside the window is included and ordered by that stamp.
**Validation:** Inferred during capture from the maintainer's accepted proposal.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-kanban-board.md` — 4820 tokens; matches a queue-kanban view change, but the satellite is `slugged: partial`, so no targeted form is legal and the bare entry exceeds the 2000-token capture budget.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5744 tokens; matches queue-kanban model, UI, and timeline behavior, but the satellite is `slugged: partial`, so no targeted form is legal and the bare entry exceeds the 2000-token capture budget.

## Assets

None.

---
*Source: "capture that as a REQ" (accepting the proposal: a "recently touched" window keyed on the newest stamp on a ticket, `updated_at` or the hold time)*
