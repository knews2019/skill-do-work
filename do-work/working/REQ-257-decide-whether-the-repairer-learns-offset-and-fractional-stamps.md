---
id: REQ-257
title: Decide whether the timestamp repairer learns offset and fractional stamps
status: claimed
claimed_at: 2026-08-18T21:16:24Z
route: B
created_at: 2026-08-18T17:49:24Z
status_changed_at: 2026-08-18T20:55:14Z
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
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Decide Whether the Timestamp Repairer Learns Offset and Fractional Stamps

## What

REQ-246's repairer deliberately refuses stamps with numeric UTC offsets (`2093-01-01T00:00:00+02:00`) or fractional seconds — repairing them needs timezone arithmetic, and a wrong guess would rewrite a correct stamp (REQ-246 D-04, documented in the script header). The board and forensics still detect and warn on those shapes, so they remain a detection-without-repair residual. **Correction (REQ-255 review, I-2):** these are no longer the *only* such residual — a quoted stamp with padding inside the quotes is also board-parseable and refused here; that one is tracked in REQ-267. This asks whether that residual matters enough to implement offset arithmetic in `comparison_key_for`, or whether the documented refusal is the permanent answer.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Builder-discovered on REQ-246 (Discovered Tasks, first item), classified [normal]. Gated behind REQ-255 so the repairer's shape handling settles once, on top of the parity sweep, rather than twice.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-246: the repairer refuses offset/fractional stamps that the board still warns about — implementing offset arithmetic in `comparison_key_for` would close the last board-detectable-but-unrepaired timestamp class. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the documented refusal (D-04) is the permanent answer and the board's warning is disclosure enough.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (decide, then either implement offset/fractional arithmetic or make the refusal permanent and provable) but grounding that decision needs discovery: how `comparison_key_for` and `extract_timestamp_fields` currently recognize shapes after REQ-255's parity sweep, and what the read-side detectors accept.

**Planning:** Not required

---

## Exploration

`skills/do-work/scripts/repair-req-timestamps.sh` holds the whole shape surface in two functions the REQ names:

- `comparison_key_for` (line ~206) — turns a value token into a sortable key; ends by gating on `calendar_components_valid` (line ~179). This is the single place a shape is recognized or refused, and REQ-247's auditor sources this file, so widening it widens the auditor in the same edit.
- `extract_timestamp_fields` (line ~259) — whitespace-token extraction; REQ-255 taught it the quoted and unquoted space-separated spellings.
- The script header (line ~51) states the current refusal: *a numeric UTC offset or fractional seconds is not provably wrong without timezone arithmetic* (REQ-246 D-04).

Lock-ins live in `_dev/tests/prescribed-shell-scripts-behavior.sh` as a `# repair-req-timestamps:` comment-headed case group (lines ~1098–1500), including the two REQ-255 space-separated cases and the skew-constant lock-step case at ~1474 — the pattern any new case follows.

*Explored inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — the decision's implementation: either offset/fractional recognition in `comparison_key_for`, or the refusal made permanent and provable
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-ins for whichever answer wins

**Files I will NOT touch:** `skills/do-work/scripts/audit-archive-timestamps.sh` (it sources the repairer, so it inherits the change without an edit), the SessionStart hook, the board tool, `CHANGELOG.md`, `VERSION`, `skills/do-work/actions/version.md`.

**Acceptance criteria (restated from REQ):**
- [ ] The offset/fractional residual is decided either way, and the decision is stated where a reader of the script meets it
- [ ] Whichever way it goes, a lock-in pins it — a refusal that is only a comment is not pinned
- [ ] The one code body stays shared: nothing is duplicated into the auditor
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0
