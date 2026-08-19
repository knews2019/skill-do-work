---
id: REQ-294
title: "Make capture's impact guard symmetric so the field is judged rather than defaulted"
status: pending
created_at: 2026-08-19T15:48:05Z
user_request: UR-060
addendum_to: REQ-289
domain: general
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
related: [REQ-289, REQ-290]
write_set:
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
---

# Make Capture's Impact Guard Symmetric So the Field Is Judged Rather Than Defaulted

## What

Three separate forces in the shipped capture path all push a REQ's `impact:` to
`impact-user-visible`, and the only stated anti-invention guard blocks the one value UR-060 exists
to surface. Make the guard symmetric, or stop the template pre-filling a verdict.

## Why

UR-060's whole point, in the user's words: "if the impact is TRIVIAL, I might not want to keep
implementing it. For example if it's a minor CSS change, or wording, that the client will never
see, then it might be a good time to stop."

That signal only exists if `impact-negligible` actually gets written when it is true. Today:

1. `capture-reference.md:24` ships an **uncommented** `impact: impact-user-visible` line in the
   Simple REQ template — a pre-filled verdict, not a prompt to judge.
2. The Schema Read Contract default is `impact-user-visible` (correct, and deliberately so — absence
   must never read as the stop signal).
3. `capture.md`'s new bullet closes with "Never invent `impact-negligible` for work you have not
   actually judged" — a one-directional guard against the only value the user asked for.

Every lazy path lands on `impact-user-visible`. The field goes decorative and REQ-290's
`--skip-impact-negligible` has nothing to filter.

## Detailed Requirements

- Make the anti-invention guard **symmetric**: neither value may be asserted without the judgment
  behind it. `impact-user-visible` asserted by default is exactly as wrong as `impact-negligible`
  invented, and it is the failure that is currently invisible.
- Decide the template line deliberately and say why in one clause: either comment it out (so the
  writer supplies a judged value) or keep it uncommented with a note that it is a placeholder to be
  replaced, never a default to accept. Prefer whichever needs fewer words to be unambiguous.
- Do **not** change the Schema Read Contract default. `impact-user-visible` is right for *absence*;
  this REQ is about *emission*, which is a different question.
- `maintenance: true` — the fix is most likely a narrowing or a deletion (an asymmetric guard is a
  bad instruction, not a missing one). Try subtraction before adding a rule.

## Acceptance

- A capture that has not judged impact cannot silently emit `impact-user-visible` any more easily
  than it can emit `impact-negligible`.
- The instruction reads the same in both directions.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [ ] Should the template line be commented out, or kept uncommented with a placeholder note?
  Recommended: comment it out — an uncommented value in a template is copied far more often than it
  is judged, and the contract default already covers absence safely.
  Also: keep it uncommented with an explicit "replace this, do not accept it" note, which keeps the
  field visible to a writer who would otherwise forget it exists.

## Full Context

Finding F6 from REQ-289's review. See `do-work/user-requests/UR-060/input.md`.
