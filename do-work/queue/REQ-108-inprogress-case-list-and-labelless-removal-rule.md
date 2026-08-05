---
id: REQ-108
title: "Review fix: In-Progress Record still enumerates two recovery cases and owes no removal rule for a label-less entry"
status: pending
created_at: 2026-08-05T11:36:39Z
user_request: UR-018
addendum_to: REQ-104
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
review_generated: true
write_set: [actions/work-reference.md, actions/forensics.md, decisions/log.md, decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md]
related: [REQ-094, REQ-095, REQ-104]
batch: parallel-building
---

# Review Fix: In-Progress Record's Case List and the Label-Less Removal Rule

## What

REQ-104 dropped the label-less authorship heuristic in `actions/work-reference.md`'s Crash Recovery
ladder (a label-less checkpoint entry is now always report-only), but two consequences of that drop
were not carried into the rest of the contract:

1. `actions/work-reference.md:456` — In-Progress Record (Step 2)'s opening paragraph restates the
   classification as a **closed two-item** set: "one that is not — unnamed, or named under another
   checkout's label — is a foreign claim". Under the drop the correct non-own set is three (unnamed,
   foreign-label, label-less). The leading positive clause ("under this checkout's own writer label …
   recovers") still yields correct behavior, so nothing misbehaves today; the risk is a future editor
   re-deriving the case list from here. This is the Closed-Enumerations failure shape the skill's own
   conventions name.
2. The rewritten bullet routes reclaim of a genuinely-own pre-0.170.0 entry to "the takeover ladder
   below, or `actions/forensics.md` Check 1's manual reset" — but neither path has a rule for removing
   the label-less `## In Progress (interrupted)` entry. In-Progress Record's removal rule is scoped to
   "this checkout's **own** entry"; `actions/forensics.md:39` says to drop the "**own-label**" entry and
   leave "any entry under **another checkout's** `writer:` label untouched". A label-less entry is
   neither. Before the drop it was classified as own before recovery ran, so the own-entry rule reached
   it; it no longer does.

   Consequence: a reclaimed label-less REQ leaves a permanent checkpoint entry, and `actions/work.md`
   Step 10's session-start note forbids deleting `CHECKPOINT.md` while any "no label at all" entry
   remains — a phantom claim re-reported every session with no documented exit.

## Context

Found during review of REQ-104 (Important findings 1 and 2). Two Minor items travel with it:
`decisions/log.md:106` still records the edge as "filed rather than fixed in flight" while ADR-018's
Consequences now says resolved; and ADR-018's frontmatter `updated:` was not bumped for the
2026-08-05 body edits (`adr-001` sets the precedent that it is maintained on amendment).

## Detailed Requirements

- Widen `actions/work-reference.md:456`'s restatement to name all three non-own cases, or restate the
  condition rather than the list (per the Closed-Enumerations rule, prefer stating the condition).
- State who removes a label-less `## In Progress (interrupted)` entry when a human reclaims the REQ —
  in In-Progress Record (Step 2)'s removal rule and in `actions/forensics.md` Check 1's remediation.
  Both currently name only the two labeled cases.
- Update `decisions/log.md:106` to note the edge was closed by REQ-104, and bump ADR-018's `updated:`.
- Do not add liveness machinery; UR-018's ban is unchanged.
- Consider whether the reworded In-Progress Record paragraph needs a suite pin, or whether REQ-104's
  existing pair suffices.

---
*Source: REQ-104 review findings (Important 1 + 2, Minor 3 + 4)*
