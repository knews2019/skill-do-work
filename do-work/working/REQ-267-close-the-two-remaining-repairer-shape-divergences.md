---
id: REQ-267
title: Close the two remaining repairer shape divergences
status: claimed
claimed_at: 2026-08-18T22:59:48Z
route: B
created_at: 2026-08-18T21:03:15Z
status_changed_at: 2026-08-18T22:20:09Z
user_request: UR-056
addendum_to: REQ-255
domain: general
review_generated: true
sweep: true
sweep_key: repairer-detector-shape-parity
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Close the Two Remaining Repairer Shape Divergences

## What

REQ-255 closed six shape divergences between the timestamp repairer and the board's readers. Its independent review, fuzzing the whole value space, found two more — both pre-existing at REQ-255's branch point, neither among its declared instances, and neither corrupting.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] **An unterminated frontmatter block is repaired where the board sees only body text (reproduced by execution).** A file whose opening fence never closes is scanned to end-of-file by the extractor and its stamp lines are rewritten by the unattended hook, while the board's `splitFrontmatter` returns *no frontmatter* for exactly that shape. The script's own scope comment states the fence-bounded contract the code does not honour here. **Second face of the same root cause:** when such a file also ends with the defective stamp on its final line with no trailing newline, the changed-line guard expects four diff lines and sees two, so the repair is refused and the SessionStart hook exits 1 **every session, permanently**, with no self-heal. Refusing the fence-broken shape the way the read side does closes both at once.
- [ ] **Quoted stamps with padding inside the quotes are board-parseable but refused, and the refusal is undocumented (reproduced by execution).** A Go probe replicating the board's pipeline (YAML unquote, then trim, then parse) accepts `"2093-01-01 00:00:00 "` and would flag it future; the repairer refuses it byte-identical and the refusal falls into the header's catch-all "anything else unparseable", which is false for this shape. The header's parity rule claims the opposite family-wide, so the documentation is wrong rather than merely silent.

- [ ] **The clean-fixture comment still justifies the offset refusal with the reason REQ-257 repudiated (added from REQ-257's review, Important 3).** `_dev/tests/prescribed-shell-scripts-behavior.sh:1150` says the repairer must not touch "a numeric-offset value **it cannot compare without timezone arithmetic**" — verbatim the justification REQ-257's header now states is wrong ("the arithmetic is the risk, not the obstacle"). Same family as instance 2: a statement about what the repairer refuses that is wrong rather than merely silent. This file is already in this REQ's write set.

## Requirements

- Each instance is repaired to canonical form or refused byte-identical **with the refusal documented beside the existing refusal entries** — the never-half-rewrite rule from REQ-255 still governs.
- The permanent hook-failure loop is gone: no input shape may leave the SessionStart hook exiting nonzero every session with no path to self-heal.
- Lock-in cases per instance, through both scan scopes.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-255's independent review, findings I-1 and I-2 (gate: trivial each — never-corrupt holds, and both shapes are already malformed or exotic). Created `pending-answers` per the generation-≥2 cascade stop. The review also noted a consequence for a queued sibling: REQ-257's description called offset and fractional seconds "the last" board-detectable-but-unrepaired class, which instance 2 makes inaccurate — corrected in that REQ directly.

## Open Questions

- [x] REQ-255's review found two more shapes where the repairer and the board disagree — one repairs a file the board reads as having no frontmatter (and, in one variant, makes the session hook fail forever), the other refuses a shape the board accepts while the docs claim it is handled. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`, presented as the only known live defect in the queue that can wedge a session permanently. Nothing was put out of scope — both instances stand, and the never-half-rewrite rule from REQ-255 still governs the refusal path.
