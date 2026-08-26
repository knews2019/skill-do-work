---
session_ended: 2026-08-26T15:05:00Z
last_completed: REQ-374
queue_state: [3 pending, 0 finished (awaiting archive), 3 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 1 in-progress]
reqs_processed_this_session: 1
session_depth: light
---

# Session Checkpoint

## Completed This Session

- REQ-374 (Route B, review 77% → Pass after remediation) put the implementation span on every
  Recently-Done card: `took 2h 40m`, `likely paused` past the four-hour read-time ceiling,
  `reversed stamps` where the bookkeeping is broken. Measured once in Go and shipped ready to
  draw, so the card and the Durations chart read one rule. Commit `5ad1d3d`.

Earlier sessions, kept for the trail:

- REQ-372 (Route A, review 100%) established one canonical two-path response when required files
  fall outside a declared write set.
- REQ-373 (Route A, review 100%) made project-harness membership the explicit boundary for
  `tdd: true` evidence and routed probe-only work to `tdd: false` plus repeatable proof.

## In Progress (interrupted)

- REQ-378 — claimed 2026-08-26T13:20:45Z — writer: vm:/home/user/skill-do-work

## Still Queued

From REQ-374's Discovered Tasks — each needs a yes/no before it can run. Run `do-work clarify`;
each flips to `pending` on a yes:

- REQ-375 `pending-answers` — restore the strict browser lane on current Chromium (`impact-rule-change`)
- REQ-376 `pending-answers` — raise the done line's faint text to readable contrast (`impact-user-visible`)
- REQ-377 `pending-answers` — stop preflight scratch showing up as untracked (`impact-negligible`)

From UR-075 and UR-076 — the ticket-id autocomplete program, dependency-ordered:

- REQ-379 `pending` — copy carries titles and a referenced-requests glossary (`depends_on: REQ-378`)
- REQ-380 `pending` — Cross-Reference Convention for newly authored ids (independent)
- REQ-381 `pending` — index cited ticket ids and let the filter box match them (`depends_on: REQ-379`)

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
