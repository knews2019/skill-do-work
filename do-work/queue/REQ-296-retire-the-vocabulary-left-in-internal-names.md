---
id: REQ-296
title: Decide whether the retired vocabulary left in internal names should follow
status: pending-answers
created_at: 2026-08-19T15:48:05Z
user_request: UR-060
addendum_to: REQ-289
domain: general
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
depends_on: []
maintenance: false
related: [REQ-289]
write_set:
- skills/do-work-board/tools/queue-kanban/generate.go
- skills/do-work-board/tools/queue-kanban/timeline.go
- skills/do-work/tools/estimate-p50.sh
- skills/do-work/actions/estimate-reference.md
- _dev/tests/p50-estimator-determinism.sh
---

# Decide Whether the Retired Vocabulary Left in Internal Names Should Follow

## What

REQ-289 renamed `effort_estimate`'s values from `trivial`/`normal` to
`effort-mechanical`/`effort-substantive` everywhere they carry schema meaning. Two places still
spell the retired words on purpose, both named as explicit non-goals by the approved plan. This REQ
asks whether that should stay the answer.

## Why the plan left them

- **The board's timeline projection internals** — the JSON keys `trivialSamples` / `normalSamples` /
  `trivialMinutes` / `normalMinutes` in `generate.go`, and `TrivialMedianMinutes` /
  `NormalMedianMinutes` in `timeline.go`. These are internal payload names with no schema meaning;
  the user-facing labels they feed **were** renamed to mechanical/substantive. Renaming them costs
  `board-timeline.js` and `generate_test.go` for no user-visible gain.
- **`tools/estimate-p50.sh --trivial`** and its `- trivial short-circuit` basis string, pinned by
  `_dev/tests/p50-estimator-determinism.sh:77-80`. The flag names the estimator's *floor mode*, not
  the schema token — a distinction REQ-289 documented in `estimate-reference.md` rather than
  silently leaving ambiguous.

## The case for doing it anyway

An internal name that disagrees with the label it renders costs reader time at exactly the moment
someone is debugging why a projection looks wrong. UR-060's premise is that a token which means two
things is worth money to fix; these are the last two places in the codebase where the retired
spelling survives.

## The case against

Both divergences are now *stated* rather than silent, which is the condition the maintenance posture
cares about. Renaming the estimator flag changes a documented CLI surface and a pinned test for a
cosmetic gain, and `--trivial` is arguably the correct name for a floor mode regardless of what the
schema calls its values.

## Open Questions

- [ ] Should the board's timeline projection internals (`trivialSamples`, `normalMinutes`,
  `TrivialMedianMinutes`, and siblings) be renamed to match the labels they render?
  Recommended: No — they are internal payload names with no schema meaning, the rendered labels are
  already correct, and the rename costs two more files plus their tests for no user-visible gain.
  Also: Yes, rename them so no spelling of the retired vocabulary survives anywhere in the module.
- [ ] Should `tools/estimate-p50.sh --trivial` and its basis string be renamed?
  Recommended: No — the flag names the estimator's floor mode rather than the schema token, that
  distinction is now documented at `estimate-reference.md`, and renaming changes a CLI surface and a
  pinned determinism test for a cosmetic gain.
  Also: Yes, rename to `--mechanical` so a reader never has to be told the two `trivial`s are
  different things.

## Full Context

Discovered Tasks 2 and 3 from REQ-289, both recorded by the builder as deliberate non-goals with
reasons. The third discovered task (no JavaScript behavior probe covers the card chips) is already
queued as instance F5 of REQ-293 and is not repeated here.

See `do-work/user-requests/UR-060/input.md`.
