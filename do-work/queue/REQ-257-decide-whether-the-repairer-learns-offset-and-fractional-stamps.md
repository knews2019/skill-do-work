---
id: REQ-257
title: Decide whether the timestamp repairer learns offset and fractional stamps
status: pending-answers
created_at: 2026-08-18T17:49:24Z
user_request: UR-056
addendum_to: REQ-246
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-255]
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Decide Whether the Timestamp Repairer Learns Offset and Fractional Stamps

## What

REQ-246's repairer deliberately refuses stamps with numeric UTC offsets (`2093-01-01T00:00:00+02:00`) or fractional seconds — repairing them needs timezone arithmetic, and a wrong guess would rewrite a correct stamp (REQ-246 D-04, documented in the script header). The board and forensics still detect and warn on those shapes, so they remain a detection-without-repair residual. This asks whether that residual matters enough to implement offset arithmetic in `comparison_key_for`, or whether the documented refusal is the permanent answer.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Builder-discovered on REQ-246 (Discovered Tasks, first item), classified [normal]. Gated behind REQ-255 so the repairer's shape handling settles once, on top of the parity sweep, rather than twice.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-246: the repairer refuses offset/fractional stamps that the board still warns about — implementing offset arithmetic in `comparison_key_for` would close the last board-detectable-but-unrepaired timestamp class. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the documented refusal (D-04) is the permanent answer and the board's warning is disclosure enough.
