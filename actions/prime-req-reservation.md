# Prime: req-reservation

REQ reservations — the `reserved` status that allocates pending queue items to a *different* worktree or cloud session sharing this repo. Spans `actions/reserve.md` (the verb), the work loop's scan/claim rules, and the Kanban board's grayed-out rendering.

## Read first
- `actions/reserve.md` — reserve / release / list; the only writer of `reserved`
- `actions/work-reference.md` — Schema Read Contract (status enum) + exit-summary section 5
- `actions/work.md` — Step 1 special-statuses + stale check; Step 2 reservation-consuming claim
- `tools/queue-kanban/model.go` — `reserved` bucketing + `reservationStaleAfter`

## Traps
- **Reservation ≠ claim** — a reserved REQ stays in `do-work/queue/`; anything moved to `do-work/working/` gets force-re-queued by crash recovery at the next session start, silently destroying the allocation. Never "reserve" by claiming.
- **Targeted runs bypass reservations by design** — `do-work run REQ-NNN` claims a reserved REQ and clears `reserved_for`/`reserved_at`. That's the owning session's pickup path, not a bug; only the default full-queue scan honors reservations.
- **An unsynced reservation protects nothing** — checkouts only see it after the queue edit is committed and pushed; files are the sole channel between sessions.
- **The 24h staleness threshold lives in two places** — `actions/work.md` Step 1 (prose) and `tools/queue-kanban/model.go` `reservationStaleAfter` (code). Change them together. Staleness only ever *suggests* recategorizing (release / claim here / leave it); nothing auto-releases.
- **The status enum is closed** — adding a reservation-adjacent status means updating the Schema Read Contract row, work.md's special-statuses and `--wave` lists, the exit summary, cleanup/abandon/roadmap/forensics readers, and `model.go` in the same commit (see CLAUDE.md → Closed Enumerations Go Stale).

## Stakes
- `status: reserved` + `reserved_for`/`reserved_at` (Schema Read Contract, `actions/work-reference.md`)
  Req:   mark a pending REQ as allocated to a named other session, in the queue file itself, so every sibling checkout sees the same allocation after a git sync — while remaining invisible to crash recovery and the default claim scan.
  Value: one shared queue can feed parallel worktrees/cloud sessions without double-claiming; the label makes ownership auditable by a human.
  Risk:  if a reader treats `reserved` as unrecognized (or claims it by default), parallel sessions double-build the same REQ; if release stops restoring exactly `pending`, REQs leak out of the queue. Reversible per-REQ (edit frontmatter), but double-built work is not.
- 24h staleness (`actions/work.md` Step 1 ↔ `model.go` `reservationStaleAfter`)
  Req:   surface reservations whose owning session likely died, with a recategorize suggestion — never an auto-release.
  Value: dead sessions can't strand queue items forever; the user decides, so a slow-but-alive session is never robbed mid-build.
  Risk:  auto-releasing (or silently dropping the flag) reintroduces the exact double-claim race reservations exist to prevent.
