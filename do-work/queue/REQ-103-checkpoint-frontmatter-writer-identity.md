---
id: REQ-103
title: "Checkpoint frontmatter carries no writer identity — session-start resume banner can report another checkout's session"
status: pending-answers
created_at: 2026-08-04T20:08:59Z
user_request: UR-018
addendum_to: REQ-094
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-094]
batch: parallel-building
---

# Checkpoint Frontmatter Writer Identity

## What

Discovered during REQ-094 (`[low]`): `do-work/CHECKPOINT.md`'s **frontmatter** (`session_ended`, `last_completed`, `queue_state`, `session_depth`) carries no writer identity, so on a synced tree the Step 10 session-start resume banner ("Resuming from previous session. Last completed: REQ-NNN...") can report another checkout's session as this one's. Harmless today — nothing classifies or strips on those fields — but it is the same ambiguity the In-Progress writer label just closed one section above it.

## Open Questions

- [ ] Should the checkpoint frontmatter gain the same `writer:` label (banner then says whose session it summarizes), or is the resume banner cosmetic enough to leave as-is?
  Recommended: leave as-is until the two-clone acceptance run (REQ-095) shows the banner actually misleading in practice — nothing acts on these fields, and every field added to the checkpoint is another thing Step 10's rewrite must preserve correctly.
  Also: add `writer:` to frontmatter now (one-line change, symmetric with the entry label), or drop the resume banner on synced checkouts entirely.

---
*Source: REQ-094 builder, Discovered Tasks ([low])*
