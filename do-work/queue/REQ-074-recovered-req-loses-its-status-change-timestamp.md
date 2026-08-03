---
id: REQ-074
title: A recovered REQ loses the timestamp that says when it was reset
status: pending-answers
created_at: 2026-08-03T14:48:01Z
user_request: UR-013
addendum_to: REQ-071
domain: general
prime_files: []
tdd: false
depends_on: []
maintenance: false
review_generated: false
discovered_during: REQ-071
---

# A Recovered REQ Loses the Timestamp That Says When It Was Reset

## What

Two places in the skill reset a stuck REQ back to `pending`, and they disagree about whether to
record *when* the reset happened.

- When a **human** does it by hand, `actions/forensics.md` Check 1 tells them to stamp
  `status_changed_at` with the current time.
- When the **pipeline** does it automatically, `actions/work-reference.md`'s Crash Recovery substep 1
  resets the status and stamps nothing.

`status_changed_at` is the field the Kanban board reads to show how long a card has been sitting in
its current state. Without it, an automatically-recovered REQ falls back to its original creation
date, so the board shows a REQ that was reset two minutes ago as though it had been waiting since the
day it was written.

## Why This Is Worth Deciding Rather Than Just Fixing

The one-line fix is obvious (have Crash Recovery stamp the field too), which is exactly why it is
worth pausing: the field's own schema note says it is written **on any status flip that has no
dedicated timestamp of its own**, and lists the manual reset as one of the writers. By that rule
Crash Recovery has been out of compliance since the field was introduced — so the question is not
really "should we stamp it here" but "is the rule right, or is the omission deliberate?" Getting that
backwards writes a rule into a second place and makes the disagreement permanent instead of fixing it.

Nothing is broken today beyond a misleading age on one board card. There is no data loss and no
pipeline behavior depends on the field — the schema calls it display-only.

## Context

- `actions/work-reference.md` → **Crash Recovery (Step 1)**, substep 1 — resets `status`, removes
  `claimed_at` and `route`, stamps nothing.
- `actions/forensics.md` → Check 1 (Stuck Work), suggested remediation — the manual reset, which does
  stamp `status_changed_at`.
- `actions/work-reference.md` → **Request File Schema — Full Frontmatter**, the `status_changed_at`
  line — states the trigger condition ("any status flip that has no dedicated `*_at` stamp of its
  own") and names its writers illustratively, including the manual reset.
- `tools/queue-kanban/model.go` — the consumer: the board's state timer prefers `status_changed_at`
  over `created_at` and file mtime for pending-tier cards.

## Discovered During

REQ-071 (`do-work/archive/` — crash recovery respects a live claim). Found by that REQ's restatement
sweep while checking that the automatic and manual reset procedures agreed with each other. It is
pre-existing and was outside REQ-071's requirements, so it was filed rather than swept in.

## Open Questions

- [ ] Crash Recovery resets a REQ to `pending` without recording when — should it stamp
  `status_changed_at` like the manual reset does?
  Recommended: **Yes — stamp it in Crash Recovery substep 1.** One line in
  `actions/work-reference.md`; brings the automatic path in line with both the field's stated trigger
  condition and the manual path, and the board stops dating a just-recovered REQ from its creation day.
  Value: the board's "waiting for N days" figure becomes true for recovered REQs, and the two reset
  procedures stop contradicting each other.
  Risk: very low and reversible — the field is display-only by contract, nothing in the pipeline
  reads it, and no other REQ's data changes.
  Also: **(b) Leave Crash Recovery as it is and narrow the field's rule instead** — amend the
  `status_changed_at` schema line to say recovery resets are deliberately excluded, and say why. Pick
  this if the omission was intentional and the board should keep showing the original queue age.
  Also: **(c) Do nothing.** The disagreement stays, and the next person to notice it re-discovers it
  from scratch. Cheapest now, and the reason this REQ exists.

## Full Context

See `do-work/user-requests/UR-013/input.md` for the original batch input, and REQ-071 for the sweep
that surfaced this.
