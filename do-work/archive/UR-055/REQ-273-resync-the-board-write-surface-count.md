---
id: REQ-273
title: Resync the board tool's write-surface count where it is restated
status: cancelled
created_at: 2026-08-18T23:13:36Z
status_changed_at: 2026-08-18T23:13:36Z
completed_at: 2026-08-20T13:21:13Z
user_request: UR-055
addendum_to: REQ-261
domain: general
review_generated: true
sweep: true
sweep_key: board-write-surface-count-restated
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work-board/tools/queue-kanban/timestamp.go
- skills/do-work-board/tools/queue-kanban/open_work.go
---

# Resync the Board Tool's Write-Surface Count Where It Is Restated

## What

The queue-kanban tool has **three** write surfaces (the board's Testing view, `next-version`, and `next-req`). Two Go comments still say two — they drifted unseen when `next-req` added the third:

- `skills/do-work-board/tools/queue-kanban/timestamp.go:24` — "the tool's **two-write-surface** contract is unchanged"
- `skills/do-work-board/tools/queue-kanban/open_work.go:22` — "Read-only, like every subcommand but the **two named write surfaces** (the board's testing view and `next-version`)" — which also names only two of the three

`frontmatter_cli.go:34` and `main.go:49` say three and are correct, so the codebase currently contradicts itself about a counted contract.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- Every in-code restatement of the write-surface count agrees with the actual count.
- **Sweep the primitive rather than fixing two lines.** These two were found by a reviewer sweeping for something else, so they are a sample: grep the whole repo for any restatement of this count — in Go comments, shipped markdown, primes, and tests — and verify each against the real set. Report what was found either way.
- **Ask whether the count should be restated at all.** `_dev/tests/contract-regressions.sh:2236` pins only the `CLAUDE.md` sentence, which is why these two drifted silently for a whole surface. The durable answer may be to stop repeating the number in code comments and point at the one canonical statement instead — deletion before addition, per `crew-members/maintenance.md`. That judgment is builder latitude; whichever way it goes, state it.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

Found by REQ-261's independent review while running its Restatement Sweep — not a finding against that REQ's diff, which redefines nothing about write surfaces, but a contradiction the reviewer hit while sweeping and declined to leave unreported.

This is the project's own **Closed Enumerations Go Stale** rule biting the file that enforces it: `CLAUDE.md` § Kanban Board Write Surfaces says "the rule is the count, not any list of subcommands" and requires the sentence to be amended in the same commit as a new surface. That happened in `CLAUDE.md` and in two Go files; it did not happen in these two, and only the `CLAUDE.md` sentence is mechanically pinned.

## Open Questions

None — the drift and the correct count were both verified against the code by the review.

## Cancelled

- **When:** 2026-08-20T13:21:13Z
- **Why:** folded into the standing prose sweep REQ-307
- **Decided by:** user, via `do-work abandon`

**Where the work went.** This REQ's finding is now an instance on REQ-307's `## Instances`
checklist, with its file:line citations intact and re-verified against the tree on 2026-08-20
rather than carried over on trust. Nothing is dropped; what changes is that it drains in a batch
with its own class instead of costing a dispatch, a review, a version bump and two commits on its
own. That is the whole point of UR-063, and this REQ is one of the two seed instances that gives
REQ-307 something to hold on its first day.

**This REQ is also the argument for its own folding.** `CLAUDE.md`'s board write-surface paragraph
ends "Adding a fourth write surface means amending this sentence in the same commit; that is the
co-location rule applied to itself." When the third surface landed, `CLAUDE.md` and
`frontmatter_cli.go` were both amended and `open_work.go` and `testing.go` were not — so
co-location held for the sites someone was looking at and missed the two they were not. Recorded on
REQ-307 as the case for batching this class rather than trusting co-location to be remembered.
