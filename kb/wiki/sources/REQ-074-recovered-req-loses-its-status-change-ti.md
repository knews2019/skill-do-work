---
title: "Lessons from REQ-074: A recovered REQ loses the timestamp that says when it was reset"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-074-recovered-req-loses-its-status-change-ti.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-074: A recovered REQ loses the timestamp that says when it was reset

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

Two procedures reset a stuck task-queue ticket back to pending, and they disagreed about recording *when*. The by-hand procedure stamped `status_changed_at`; the pipeline's automatic crash recovery stamped nothing. Since the board's state timer reads that field, an automatically-recovered ticket fell back to its creation date and displayed as though it had been waiting since the day it was written. The field's own schema note already said it is written on any status flip with no dedicated timestamp of its own — so the real question was whether the rule was right or the omission deliberate, not whether to add a line.

## Solution summary

Crash recovery's frontmatter-reset substep now stamps `status_changed_at` on both flip targets — `pending` and `pending-answers` — with the preserved-`blocked` case explicitly excluded, since its own `blocked_at` is intact and its status does not flip. The wording cites the Timestamp rule by section name rather than restating the format. A contract-suite assertion pins the prescription so a later edit cannot quietly drop it.

## What worked

- Stashing the source file to prove the new contract assertion actually fires. A grep assertion that has only ever been run against a tree where it passes is untested — this repo already shipped one that aborted the suite silently in exactly the case it existed to catch (REQ-073's first invariant check), so the stash/run/unstash cycle is cheap insurance.

## What didn't work

- Nothing failed, but the first instinct — "one line, no test needed" — was wrong for the same reason the REQ existed: the omission survived because nothing failed when it was absent.

## Worth knowing

- `status_changed_at` now has three writers with a shared rationale (the manual reset, the unblock flips, and crash recovery) and the common thread is that each one *removes* the stamp that would otherwise date the transition — `claimed_at` here, `blocked_at` for the unblock flips. That is the real reason the field exists, and it is a better trigger test than the field list: **if the flip discards its own `*_at`, it must stamp `status_changed_at`.**

## Back-reference

See `do-work/archive/UR-013/REQ-074-recovered-req-loses-its-status-change-timestamp.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `80e7b88`.
