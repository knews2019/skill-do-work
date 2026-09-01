---
title: "Lessons from REQ-194: Forward stray REQs through verify and forensics"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-194-forward-stray-reqs-through-verify-and-fo.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-194: Forward stray REQs through verify and forensics

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Make forensics surface every REQ file that the board already detects outside `queue/`, `working/`, and `archive/`, reusing the board's detector rather than describing another scan. Align the existing archived-UR/live-member probe with REQ-193 so legitimate review-generated follow-ups do not teach users to reopen a closed UR.

## Solution summary

**[MAP CHANGED]** Queue verification now shares the board walker's structured stray-request evidence and understands the closed-UR review-follow-up marker. Misplaced REQs stay invisible as cards but are no longer invisible to forensics, while legitimate same-UR review work does not trigger advice to reopen a closed folder.

## What worked

- Retaining the walker's structured stray records on the board lets warnings and verify share one detector without turning misplaced files into cards.
- A direct seam test with no warnings or filesystem evidence proves the intended data source more strongly than an end-to-end fixture alone.

## What didn't work

- The first integration tests allowed verify to reconstruct identical output from warning prose, so the forbidden coupling survived mutation.
- Forensics initially claimed every Go-backed probe had a manual equivalent, contradicting the deliberate no-second-scan boundary for strays.

## Back-reference

See `do-work/archive/UR-043/REQ-194-forward-stray-reqs-through-forensics.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ca34ef2`.
