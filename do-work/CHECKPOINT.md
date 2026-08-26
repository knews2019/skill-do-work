---
session_ended: 2026-08-26T19:20:00Z
last_completed: REQ-379
queue_state: [3 pending, 0 finished (awaiting archive), 1 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress]
reqs_processed_this_session: 2
session_depth: deep
---

# Session Checkpoint

## Completed This Session

- REQ-378 (Route B, review Pass) put titles on every REQ/UR id the board's drawer renders, flagged
  dead references, and added a per-body glossary.
- REQ-379 (Route B, review Pass 91%) carried the same into the clipboard: bodies annotated, one
  referenced-requests appendix, frontmatter fence byte-exact.
- Five external findings verified and fixed across the two, all mutation-pinned. Released 0.239.0
  and 0.240.0; merged main at a290cd6, reconciling a second version collision.

## In Progress

None.

## Still Queued

REQ-380 and REQ-381 are ready with disjoint write sets. REQ-382 is gated behind REQ-381 because both
write `generate_test.go`. REQ-383 is at `pending-answers` — its Open Question asks whether to harden
the hand-rolled fence scanner or delete it by moving block-context detection to the Go side, and the
answer changes the order. See `do-work/RESTART-PROMPT.md`.

## Session Notes

- Current release: 0.237.0 on main; this branch bumps again when REQ-378 archives.
- **A numbering collision was reconciled here.** Two branches worked in parallel and both reserved
  from the same block: main's PR #168 took UR-074 and REQ-374, then REQ-375–377 for its follow-ups,
  while this branch had already reserved the same numbers for the ticket-id autocomplete program.
  This branch renumbered — UR-074→UR-075, UR-075→UR-076, REQ-374–377→REQ-378–381 — because main
  merged first. Reservation markers moved with them. The lesson is not "renumber carefully": the
  `.req-reservations` markers are per-checkout, so two checkouts that never fetch each other reserve
  the same ids by construction, and nothing detects it until the merge.
- UR-074 stays open in `do-work/user-requests/` — REQ-374 archived to `do-work/archive/` root
  because its three follow-ups are not terminally resolved.
- The canonical gate could not run in this container until three tools were installed:
  ShellCheck 0.11.0 (apt ships 0.9.0, below the floor — fetched the static release),
  `just`, and the Go 1.26.1 toolchain via `GOTOOLCHAIN`. Both sessions hit this independently, so a
  future session in a fresh container will hit it again. It passes unpiped, exit 0.
- The strict browser lane fails at HEAD for a reason that predates both sessions — that is REQ-375.
  Confirmed again from this branch by stashing the diff and re-running: identical failure, so it is
  not REQ-378's.
- A Chromium **is** available in this container at `/opt/pw-browsers/chromium` (Playwright
  chromium-1194). Passing it as `QUEUE_KANBAN_BROWSER` runs the browser lane for real instead of
  skipping it — worth doing before believing any rendering claim.
