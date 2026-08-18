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
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` including its REQ-255, REQ-250, REQ-246 and REQ-243 lessons, four crew-member files, the whole repairer, `audit-archive-timestamps.sh`, the board's `parseTimestamp` and `detectFutureTimestampFields` seams, forensics Check 11, and the existing `# repair-req-timestamps:` case group. Plan: establish **empirically** what the read side actually does with each shape — the REQ's premise, which turned out to be wrong — then decide on merits, then either implement or pin. In the pin case the RED has to come from a mutation, because a lock-in for existing behaviour cannot go red on its own.
- [x] **[APPLY]:** Stayed inside the two declared files. Four commits, each individually green.
- [x] **[UNIFY]:** `git diff --stat 662788c..HEAD` → 2 files, +76/−7; both diffs read in full. `bash -n` and `shellcheck -S warning` clean on the repairer (exit 0). The suite file's shellcheck output carries only pre-existing info-level findings (SC2016/SC2086/SC2015/SC2329) in fixture-writing lines from earlier REQs — none on any added line, and the gate's own warning-level lane passes. Grepped the added lines for `TODO|FIXME|XXX|MUTATION|console.log|set -x|echo DEBUG` → none. `git status --porcelain` empty. The RED mutation was reverted with `git checkout --` and the file confirmed byte-identical to its committed state **before** the gate run.

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

---

## Implementation Summary

**Files changed:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** The offset/fractional refusal was **kept and made permanent**, and no arithmetic was added — `comparison_key_for` is byte-identical to its pre-REQ state. What changed is the reason and the enforcement. REQ-246 D-04 justified the refusal as *undecidable without timezone arithmetic*; that is wrong, since an RFC3339 offset denotes an exact instant. The header now states the real reason — the arithmetic is the **risk**, not the obstacle — with the concrete hazard spelled out: `2026-08-19T00:29:11+05:00` denotes `2026-08-18T19:29:11Z`, so a repairer that reads the wall clock and ignores the offset sees a value five hours later than the instant and erases a correct stamp as future-dated, unattended, from a SessionStart hook. Refusing can only fail to fix; repairing can destroy. The suite gained two case blocks — a six-shape refusal case and a read-side layout lock-step case that fails if the board's `parseTimestamp` layouts change underneath this decision — plus a numeric-offset fixture in the archive-parity block, taking the named case count from 64 to 66. `audit-archive-timestamps.sh` needed no edit; it sources the repairer and inherits the behaviour.

---

## Discovered Tasks

Transcribed by the orchestrator from `do-work/runs/work-2026-08-18-211613/REQ-257-handback.md` (a worktree builder cannot write this file — REQ-270).

- **[normal] The board's future warning prints a value that is not in the file.** For `created_at: 2093-01-01T00:00:00+02:00` the warning reads `created_at 2092-12-31T22:00:00Z`; for `2093-01-01T00:00:00.500Z` and `2093-01-01 00:00:00.5` it reads `2093-01-01T00:00:00Z`. The YAML layer types these as timestamps and `normalizeFrontmatterValue` re-formats them to RFC3339 UTC before the warning is built. A user who greps the file for the value the board named will not find it, and the warning's own remediation — "rewrite with the current UTC instant" — is aimed at a string that is not there. Observed by execution, not fixed: it lives in the board package, outside this REQ's write set. **This matters more now than before**, because this REQ's decision makes the board's warning the *only* signal for these shapes.
- **[low] `created_at: 2093-01-01 00:00:00Z` is repaired here but is not board-parseable** — a space separator with a `Z` matches no `parseTimestamp` layout, and the board emitted no warning for it. The repairer is *broader* than the read side in this one direction, which is benign (an unreadable value is replaced with a readable one). Recorded as a known asymmetry rather than proposed as a REQ.

