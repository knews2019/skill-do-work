---
id: REQ-272
title: Point the future-stamp citations at forensics Check 12
status: cancelled
created_at: 2026-08-18T22:57:49Z
status_changed_at: 2026-08-18T22:57:49Z
completed_at: 2026-08-20T13:21:13Z
user_request: UR-056
addendum_to: REQ-257
domain: general
review_generated: true
sweep: true
sweep_key: forensics-check-number-citations
effort_estimate: trivial
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- skills/do-work/actions/work-reference.md
---

# Point the Future-Stamp Citations at Forensics Check 12

## What

Four shipped citations send a reader to `do-work forensics` **Check 11** for the future-dated-timestamp check. That check is **Check 12**. Check 11 is Unrecognized Status Vocabulary (`skills/do-work-toolbox/actions/forensics.md:143`); Future-Dated Timestamps is at `:156`, and `forensics.md:201` already cross-references it correctly as Check 12 — so the file contradicts itself and the wrong number is the majority spelling.

A reader following any of these pointers lands on the status-vocabulary check and finds nothing about timestamps.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `skills/do-work/scripts/repair-req-timestamps.sh:7`
- [ ] `skills/do-work/scripts/repair-req-timestamps.sh:22`
- [ ] `skills/do-work/scripts/repair-req-timestamps.sh:102`
- [ ] `skills/do-work/actions/work-reference.md:283`

## Requirements

- Every citation of the future-dated-timestamp forensics check names the number it actually has.
- **Sweep the primitive, do not fix four lines.** The four above came from one review reading one REQ's diff, so they are a sample: grep for every citation of a forensics check *by number* across the shipped tree and verify each against `forensics.md`'s actual headings. Any other mis-numbered citation is in scope; report what was found either way.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-257's independent review, Important finding 4 (gate: trivial). Pre-existing rather than introduced by that REQ — but three of the four sites are inside the header REQ-257 rewrote, and REQ-257's disclosure argument (*the board's warning and forensics still surface these shapes, so refusing loses nothing*) rests on the pointer being followable. A decision justified by a signal a reader cannot find is weaker than it looks.

Worth knowing while this sits in the queue: nothing mechanical checks a by-number citation of a section heading. If the sweep finds this class is more than four sites, the durable answer may be to stop citing checks by number at all — a name survives an insertion, a number does not.

## Open Questions

None — the review verified all four sites and the correct number against `forensics.md` directly.

## Cancelled

- **When:** 2026-08-20T13:21:13Z
- **Why:** folded into the standing prose sweep REQ-307
- **Decided by:** user, via `do-work abandon`

**Where the work went.** This REQ's finding is now an instance on REQ-307's `## Instances`
checklist, with its file:line citations intact and re-verified against the tree on 2026-08-20
rather than carried over on trust. Nothing is dropped; what changes is that it drains in a batch
with its own class instead of costing a dispatch, a review, a version bump and two commits on its
own. That is the whole point of UR-063, and this REQ is one of the two seed instances that gives
REQ-307 something to hold on its first day.
