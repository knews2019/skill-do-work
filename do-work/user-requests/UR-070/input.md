---
id: UR-070
title: Needs Input · Blocked shows only actionable cards
created_at: 2026-08-24T14:03:59Z
requests: [REQ-366]
word_count: 434
---

# Needs Input · Blocked Shows Only Actionable Cards

## Summary

The NEEDS INPUT · BLOCKED column is the operator's inbox and must contain only cards the operator can act on right now. Today a `status: blocked` REQ with an unmet `depends_on` sits there even though its external condition cannot be satisfied until the dependency lands. Requested rule: such a ticket renders in PENDING with the existing waiting-on-deps treatment (plus its blocked badge) and is promoted to NEEDS INPUT · BLOCKED only when every dependency is completed. Applies everywhere the bucket is counted or rendered; presentation and counting only.

## Extracted Requests

| REQ | Title |
|---|---|
| REQ-366 | Keep dependency-gated blocked REQs out of Needs Input · Blocked |

## Batch Constraints

Single REQ — no batch-level constraints beyond the scope guards recorded in the REQ.

## Full Verbatim Input

First message (with screenshot):

> is it crazy that I don't want to see anything in that column that I can not act on?
> how can we make that happen?

Second message:

> more info:
>
> Board wish: NEEDS INPUT · BLOCKED must contain only cards I can act on right now
>
> I don't want to see anything in that column that I cannot act on. The column is my inbox: if a card is there, the deal should be "you read it, you do something, it leaves." Today that contract breaks for dependency-gated blocked REQs.
>
> Concrete case: an implementation REQ carries status: blocked with blocked_by: "owner explicitly approved the architecture decision report from REQ-A" and depends_on: [REQ-A], while REQ-A (the report) is still pending. The report does not exist yet, so "approve the report" is not something I can do — yet the card sits in NEEDS INPUT · BLOCKED, possibly for weeks, training me to ignore the column. The bucketing in tools/queue-kanban is status-only (pending-answers + blocked + the unrecognized-status catch-all go to NeedsInputOrBlocked); dependency readiness is never consulted, even though the PENDING column already computes exactly that for its ready/waiting split.
>
> Requested rule: a blocked ticket with unmet depends_on is waiting on its dependency first, not on me. Render it in the PENDING column with the existing waiting-on-deps treatment (NEEDS badge) plus its blocked badge, and promote it to NEEDS INPUT · BLOCKED only when every dependency is completed — the moment its external condition becomes the sole gate and therefore actionable. Apply the same rule everywhere the bucket is counted or rendered: the Board lens, the summary command's needs-input/blocked counter (such tickets count under pending → waiting on deps), the open-work digest line, and the timeline/calendar state coloring if they distinguish the bucket.
>
> Scope guards: pending-answers tickets keep their current placement (their questions are answerable regardless of deps — that's why answering them is called clarify, not unblock). Stakeholder-question REQs keep their current placement (getting answers from the named person is always actionable). The unrecognized-status catch-all stays in NEEDS INPUT so nothing ever goes invisible. blocked_check probe semantics are untouched — this is presentation and counting only, no frontmatter or scheduling change, so verify and work behavior are unaffected.
>
> Acceptance: a blocked REQ with an unmet dependency appears in PENDING/waiting with both badges; the same REQ appears in NEEDS INPUT · BLOCKED the moment its last dependency archives as completed; summary counts move accordingly; a blocked REQ with no depends_on (or all deps completed) behaves exactly as today. Invariant to state in the board's docs: NEEDS INPUT · BLOCKED contains only tickets the operator can act on right now.

Third message (routing instruction):

> the goal is to capture the intent via do-work capture-request

---
*Captured: 2026-08-24T14:03:59Z*
