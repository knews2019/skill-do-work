---
id: REQ-109
title: "work.md session-start note still enumerates the recovery case list and calls a label-less entry a foreign claim"
status: pending-answers
created_at: 2026-08-05T11:44:27Z
user_request: UR-018
addendum_to: REQ-108
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: [actions/work.md]
related: [REQ-104, REQ-108]
batch: parallel-building
---

# work.md Session-Start Note: Recovery Case List Terminology

## What

Discovered during REQ-108 (`[low]`): `actions/work.md`'s Step 10 session-start note (the sentence
listing which `working/` REQs recovery may strip) carries the same closed-enumeration shape REQ-108
just removed from `actions/work-reference.md`'s In-Progress Record — "one that isn't (unlabeled, or
labeled for another checkout) is a foreign claim recovery must not strip." The set is complete and
the behavior is correct, but it calls the label-less case a *foreign claim*, whereas since REQ-104
the canonical term is a *claim of unknown origin* — and the very next sentence in the same note uses
the correct term. Same fix shape as REQ-108: state the condition, defer the list to Crash Recovery.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-108: `actions/work.md`'s session-start note restates the recovery case list with pre-REQ-104 terminology (calls a label-less entry a "foreign claim"; behavior is correct, wording predates the drop). Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---
*Source: REQ-108 builder, Discovered Tasks ([low])*
