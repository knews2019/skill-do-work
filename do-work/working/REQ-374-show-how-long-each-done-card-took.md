---
id: REQ-374
title: 'Show how long each done card took'
status: claimed
claimed_at: 2026-08-26T13:17:00Z
created_at: 2026-08-26T13:02:22Z
user_request: UR-074
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
route: B
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-26T13:22:00Z
  basis:
    - Route B
    - 4-file write set
    - 1 new files
    - 2 subsystems involved
    - 6 acceptance criteria
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board-core.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/*_test.go
---

# Show How Long Each Done Card Took

## What

A card in the Kanban board's Recently Done column states when the work finished (`done Aug 26, 12:47 UTC · 9min ago`) but never states how long it took. Add the implementation span to that card: the time from when the builder started the REQ to when it landed in Done with a completed status.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The board already measures this span — `durations.go` computes `completed_at − claimed_at` for the Durations view — but a reader scanning the Done column has to open another view to learn whether a card took twenty minutes or half a day.

## Context

- The span is `claimed_at` → the completion instant the card already renders. `claimedAt` is already in the client payload (`generate.go`), so this needs no new frontmatter field and no second walk of the archive.
- The done line is built in `web/board-cards.js` under `options.showCompleted`, using `makeInstantWithRelativeNode` from `web/board-core.js`.
- The four-hour ceiling is the board's existing read-time rule (`durations.go`'s `analysisOutlierCeiling`, stated once in `skills/do-work/actions/estimate-reference.md` § Calibration). Reuse that rule; do not restate the number as a second definition.

## Detailed Requirements

- Every Recently Done card whose status is `completed` or `completed-with-issues` and whose `claimed_at` parses shows the implementation span alongside the existing done stamp.
- A span at or under four hours renders plainly, in the same relative-time vocabulary the board already uses (`34min`, `2h`).
- A span over four hours renders with a marker saying the session was likely paused, so an overnight pause is never read as hours of work.
- A negative span renders as a broken-stamp warning rather than a number. The board already treats reversed stamps as an anomaly worth surfacing.
- A card with no parseable `claimed_at` shows no duration at all — the rest of the card is unchanged.
- Cancelled cards, which share the Recently Done column, show no duration. The user's request is scoped to completed status.

## Constraints

- Display only. No new frontmatter field, no new write surface, no change to `model.go`'s parsed status vocabulary.
- The four-hour ceiling and the reversed-stamp rule have one definition each already. Read them; do not copy the constant into the client.
- Board changes get a normal skill version bump and a root `CHANGELOG.md` entry; the tool has no independent changelog.

## Red-Green Proof

**RED prompt/case:** Run the board (`do-work-board board`) against an archive holding a REQ with `claimed_at: 2026-08-24T10:05:00Z` and `completed_at: 2026-08-24T12:45:00Z`. Its Recently Done card reads `done Aug 24, 12:45 UTC · <relative>` and says nothing about the 2h40m the work took.
**Why RED now:** `web/board-cards.js` renders only the completion instant on a done card. The span exists in the payload (`claimedAt`) and is computed for the Durations view, but no card reads it.
**GREEN when:** That same card also reads its span — `2h40m` — and three neighbours prove the marked cases: a REQ spanning 18 hours renders its span with a paused marker, a REQ whose `completed_at` precedes its `claimed_at` renders a broken-stamp warning instead of a number, and a REQ with no `claimed_at` renders the done line exactly as it does today with no duration text.
**Validation:** User confirmed (the odd-span handling was put to the user during capture and the mark-them option was chosen).

## Assets

Screenshot of the board's Recently Done column, supplied with the request. It could not be persisted as a file in the capture session; this is the record.

The column header reads `RECENTLY DONE` with a count badge of `18` and a `Copy all` button. Seven cards are visible, each a white rounded panel with a green left edge:

- `REQ-1684` · `completed` — "[impact-user-visible] Convert the two remaining apps/admin/ native confirms to CmsConfirmDialog" · chips `frontend`, `UR-388`, `ROUTE B` · footer `done Aug 26, 12:47 UTC · 9min ago`
- `REQ-1685` · `completed-with-issues` — "[impact-rule-change] Close the preview-to-apply drift window across the four destructive card-mutation flows" · chips `cms`, `UR-389`, `ROUTE C`, `impact-rule-change` · footer `done Aug 26, 12:23 UTC · 33min ago`
- `REQ-1683` · `completed` — "Align the admin Import/Export done feedback with the extended feedback contract" · chips `frontend`, `UR-388`, `ROUTE B` · a `NEEDS REQ-1682` dependency row · footer `done Aug 26, 00:16 UTC · 12h ago`
- `REQ-1682` · `completed` — "[impact-rule-change] Extend the prime-cms-ux feedback contract to admin browser surfaces" · footer `done Aug 25, 23:58 UTC · 12h ago`
- `REQ-1681` · `completed` — "Audit e2e/prime-bowser.md — its instructions predate the layout migration" · footer `done Aug 25, 23:47 UTC · 13h ago`
- `REQ-1680` · `completed` — "Addendum: Import card confirm becomes a CMS dialog with an Import-vs-Restore fork" · footer `done Aug 25, 23:39 UTC · 13h ago`
- `REQ-1679` · `completed` — "Admin can delete a card — any card, admin-only, mapped level assets deleted too" · footer `done Aug 25, 23:27 UTC · 13h ago`

The footer line on every card is the surface this REQ extends: it carries the completion instant and its relative companion, and nothing about elapsed implementation time.

## Full Context
See `do-work/user-requests/UR-074/input.md` for complete verbatim input.

---
*Source: "So when we show the recently done cards, please also show the duration, how long it took since it was started until it is finished to implement that card to make it delivered. By making it delivered, I mean it was moved to the Done column, completed status."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified and the target files are named, but the existing conventions still need discovery — where the four-hour ceiling is read from on the client side, which test file holds the board's card-rendering probes, and how the done line is styled.

**Planning:** Not required
