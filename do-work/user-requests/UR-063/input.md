---
id: UR-063
title: Stop the pipeline generating work about its own prose
created_at: 2026-08-20T13:21:13Z
requests: [REQ-306, REQ-307]
word_count: 412
---

# Stop the Pipeline Generating Work About Its Own Prose

## Summary

A directed capture, not a review finding: the user analysed the queue's arrival rate rather than
its contents and concluded the queue is not stuck work but follow-up generation. The instruction
is to change what capture is allowed to mint, and to give the displaced class somewhere to go.

Measured by the user: **145 of 298 REQs ever written were spawned by a prior REQ**, and **four of
the five deepest descent chains terminate in prose reconciliation** — a stale count, a wrong
cross-reference number, a comment describing a superseded mechanism.

## Extracted Requests

| Request | Disposition |
|---|---|
| A prose-only discrepancy stops being its own REQ; the judgment moves to capture | REQ-306 |
| Create the standing prose-sweep REQ so REQ-306 has a destination on day one | REQ-307 |

## Batch Constraints

- **Exactly two REQs, by explicit instruction.** Anything else discovered while executing this
  batch goes in the hand-back message, not the queue. The restraint is the point: a UR about
  over-capture that spawned five REQs would refute itself.
- **The judgment already exists downstream and must not be re-derived.**
  `skills/do-work/actions/review-work.md:349` already routes an `impact-negligible` finding, or any
  set sharing one root cause, into a sweep rather than its own REQ. REQ-306 applies that same
  judgment earlier — at capture — to the class actually producing the volume. Cite it; do not
  restate it.
- **REQ-307 must exist before or with REQ-306**, or REQ-306 creates a rule with no destination.
- **REQ-272 and REQ-273 are the seed instances**, folded in and abandoned. Both premises were
  verified against the tree during this capture rather than taken from the report:
  - REQ-272: four citations say forensics **Check 11** (`repair-req-timestamps.sh:7`, `:22`,
    `:136`, and `work-reference.md:285`). Check 11 is *Unrecognized Status Values*; the
    future-stamp check is `### 12. Future-Dated Timestamps`. Confirmed wrong, all four.
  - REQ-273: `frontmatter_cli.go:34` correctly says "exactly three write surfaces", while
    `open_work.go:22` and `testing.go:42` both still say **two**. Confirmed, two stale comments.

## Full Verbatim Input

REQ A — a prose-only discrepancy stops being its own REQ. A stale count, a wrong
cross-reference number, a comment describing a superseded mechanism: these get appended to
one standing sweep REQ, or folded into the next commit that touches the file for another
reason. Reserve new-REQ creation for a change in behavior, in a checker's predicate, or in
a rule's stated condition. `skills/do-work/actions/review-work.md:349` already encodes this
judgment for impact-negligible findings — apply it at capture, in
`skills/do-work/actions/capture.md` and `capture-reference.md`, to the class actually
producing the volume. Measured: 145 of 298 REQs ever written were spawned by a prior REQ,
and four of the five deepest descent chains terminate in prose reconciliation.

REQ B — create that standing sweep REQ itself, `sweep: true`, with an empty `## Instances`
list and a stated cadence, so REQ A has somewhere to send findings on day one. Seed it with
REQ-272 (four citations point at forensics Check 11; the future-stamp check is Check 12) and
REQ-273 (two Go comments say the board tool has two write surfaces; it has three), then
abandon those two with reason "folded into the standing prose sweep REQ-NNN".

Note in REQ B that REQ-273 is evidence against its own rule: CLAUDE.md requires the
write-surface count sentence be amended in the same commit that adds a surface, and it
wasn't. That is the case for batching this class rather than trusting co-location.

---
*Captured: 2026-08-20T13:21:13Z*
