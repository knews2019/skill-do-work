---
id: REQ-293
title: Make the impact/effort lock-in checks pin the property instead of one spelling of it
status: pending
created_at: 2026-08-19T15:48:05Z
user_request: UR-060
addendum_to: REQ-289
domain: general
review_generated: true
sweep: true
sweep_key: impact-effort-lockin-checks-underpin
impact: impact-rule-change
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
depends_on: []
maintenance: false
related: [REQ-289]
write_set:
- _dev/tests/contract-regressions.sh
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Make the Impact/Effort Lock-In Checks Pin the Property Instead of One Spelling of It

## What

REQ-289's four lock-in checks all pass, and each one catches the exact defect it was written
against. Every one of them also pins a **spelling** — a literal verb, a markup shape, a partial file
set — rather than the property it claims to hold. Five instances, one root cause.

Done means the class cannot recur: each check holds its property against a re-drift written in
different words, not just against the phrasing that existed when it was authored.

## Why

This is the blind-spot lesson `_dev/primes/prime-kanban-board.md` already records twice — "asserting
a phrase is absent is not a guard: it passes when the whole string is replaced" (REQ-245), and a
fuzz's blind spots are the axes it holds constant (REQ-267). REQ-289's checks are the same shape one
layer up.

It matters now because REQ-290 depends on one of these properties (the `impact-user-visible`
default) and nothing pins it.

## Instances

- [ ] **F1 — Check A scans a strict subset of the files it claims.** Its roots are
  `skills/do-work/actions` and `skills/do-work-toolbox/actions` only. `skills/do-work/docs/` is
  excluded, and that is exactly where one of the fourteen defect sites REQ-289 fixed lived
  (`review-work-guide.md:55`, pre-change: "whose token **stamps** the follow-up's
  `effort_estimate`"). The check's own regex matches that line — only the directory list hid it.
  Also excluded: `crew-members/`, every `SKILL.md`, `do-work-knowledge/actions`,
  `do-work-board/actions`. The retired-vocabulary loop in the same file already includes
  `$core_root/docs`, so the two checks disagree about what "a shipped file" means.
- [ ] **F2 — Check A's proximity rule catches one verb out of six realistic re-drifts.** The pattern
  is `stamp\w*[^.]{0,80}effort_estimate` (either order). Two defects: the verb set is the single
  literal `stamp*`, so "derive", "set … from", "comes from", and "map … to" are invisible; and
  `[^.]` treats any period as a sentence end, so a file path between the verb and the field breaks
  the window — which is this repo's dominant citation style. Measured: 1 of 6 phrasings caught.
  The orchestrator's own mutation test used the word "stamping", the one verb the check greps, so
  that test was self-confirming and did not detect this.
- [ ] **F3 — the retired-vocabulary loop pins bold markup, not the token.** It greps
  `\*\*\[critical\]\*\*` and friends. `- [low] a style nit` and the backticked `` `[critical]` ``
  both pass clean — and the backticked form is one the tree actually carried before this REQ
  (`review-work.md:201`, "mirroring that section's `[critical]` exemption").
- [ ] **F4 — no check holds the `impact-user-visible` default.** Check D compares the token *set*
  in the schema line against `model.go`'s `canonicalValues` and never reads `defaultValue`; the
  schema line's "**Absent or unrecognized reads as `impact-user-visible`**" is matched by nothing.
  If that default ever flips to `impact-negligible`, REQ-290's `--skip-impact-negligible` inverts
  into "skip everything, including every REQ predating the field" and the suite stays green.
  **This is the highest-value instance** — it is the one property REQ-290 depends on.
- [ ] **F5 — the impact chip has no test coverage of any kind.** Neither `badge-impact` nor
  `badge-effort-estimate` appears in any test file. A cheap precedent is sitting in the same file
  the fix would touch: `generate_test.go:1226`
  (`TestGenerateInlinesWriteSetOverlapBadgeRenderPath`) pins a badge's render path by source token
  in about ten lines with no node runtime. REQ-289's own Discovered Task framed this as needing a
  JavaScript behavior probe; the cheaper precedent was next to it.

## Acceptance

- Each instance's property survives a re-drift written in different words — verify by mutation, and
  **use a different verb and a different markup shape than the one already in the file**, or the
  test repeats the self-confirming mistake above.
- Check A's scan roots and the retired-vocabulary loop's roots agree on what counts as a shipped file.
- A check fails when `defaultValue` for `impact` is changed to anything but `impact-user-visible`.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Full Context

Findings F1-F5 from REQ-289's review; F1, F2, and F3 independently reproduced by the orchestrator
before this REQ was created. See `do-work/user-requests/UR-060/input.md`.
