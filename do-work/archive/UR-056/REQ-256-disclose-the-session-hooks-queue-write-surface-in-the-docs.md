---
id: REQ-256
title: Disclose the session hook's queue write surface in the docs
status: completed
created_at: 2026-08-18T17:48:08Z
claimed_at: 2026-08-18T20:08:45Z
completed_at: 2026-08-18T20:21:34Z
commit: fbc14e8
kb_status: promoted
kb_entry: REQ-256-disclose-the-session-hook-s-queue-write-.md
route: A
user_request: UR-056
addendum_to: REQ-246
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- README.md
- skills/do-work/actions/capture.md
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T20:08:45Z
  basis:
    - trivial short-circuit (effort_estimate would be trivial; gate-stamped normal for the disclosure class)
---

# Disclose the Session Hook's Queue Write Surface in the Docs

## What

REQ-246 made the SessionStart hook a *write* surface on consumer queue files — it mechanically repairs detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/` at session start. Two shipped texts still describe the hook as read-only-plus-banner; a consumer auditing "what writes to my repo at session start" is misled.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the prime and three crew files; planned the two exact sites; closed the class by condition — grepped all shipped `*.md` for session-start/SessionStart and confirmed the two REQ instances are the whole site class. (Transcribed from builder hand-back.)
- [x] **[APPLY]:** Exactly the two planned files, one edit each. (Transcribed from builder hand-back.)
- [x] **[UNIFY]:** 2 files, +3/−1, full diff reviewed; maintainer-verify exit 0 with the reference contract validating both citation forms. (Transcribed from builder hand-back.)

## Context

Instance 1 is REQ-246's review finding I3 (restatement sweep — the omission predates REQ-246 for the reservation cleanup, but a hook that edits files unattended is a different disclosure class than one that prints a banner). Instance 2 is the optional integration seam REQ-246's builder offered and deliberately did not write. Citations follow the literal cross-package path rule REQ-249 establishes.

## Instances

- [x] **`README.md` (~line 188):** "SessionStart hook that injects the installed version and pending REQ count" — add that the hook also reaps stale REQ-number reservations and mechanically repairs detectably wrong queue/working timestamps (one clause each; keep it one sentence if it fits).
- [x] **`skills/do-work/actions/capture.md`:** document `scripts/repair-req-timestamps.sh` (SessionStart hook) the way `cleanup-req-reservations.sh` is documented — one line stating that detectably wrong queue/working `*_at` stamps are mechanically corrected at session start.

## Requirements

- Both texts state the hook's write behavior; no behavior change anywhere.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

---

## Triage

**Route: A** - Simple

**Reasoning:** Two named doc sites with the sentence content specified; no behaviour change.

**Planning:** Not required

## Plan

**Planning not required** - Route A

*Skipped by work action*

---

## Implementation Summary

**What was done:** Both doc sites now disclose the SessionStart hook's queue write surface. README's hook bullet says session-start.sh also writes to consumer queue files — reaping stale reservation markers and mechanically repairing detectably wrong `*_at` stamps — citing both scripts root-relatively. capture.md gained one paragraph beside the Immutability Rule's timestamp-repair exception, framing the session-start repair as the same metadata-correction class with an explicit "never archive". A class sweep (grep for session-start/SessionStart over shipped markdown) confirmed the two instances were the whole site class.

**Files changed (2, +3/−1):**
- `README.md` (modified) — hook bullet extended with the write disclosure
- `skills/do-work/actions/capture.md` (modified) — session-start stamp-repair paragraph beside the archive-audit exception

*Integrated by orchestrator from builder hand-back; merge range `ec0ebfd..fbc14e8`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01:** placement beside the Immutability Rule's exception, not the reservations bullet — `working/` is declared immutable there, and that is exactly where an auditor would be misled.
- **D-02:** not framed as a third "Exception" — the existing "only exceptions" enumeration is scoped to archived content, which this script never touches; the metadata-not-content framing keeps that sentence true.
- **D-03:** README cites root-relative (it sits at repo root); capture.md cites package-locally, matching the cleanup script's existing citation.

## Qualification

Passed — 2 files in merge range `ec0ebfd..fbc14e8`, both instances traced and independently re-read in the merged tree, P-A-U audited (text-only diff). Reference contract validates both citation forms in the gate.

## Review

**Overall: 96%** | 2026-08-18T20:22:00Z (Route A quick scan, orchestrated inline)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

- [x] Both texts state the hook's write behavior — verified by reading the merged files, not the hand-back.
- [x] No behaviour change — +3/−1 prose only.
- [x] D-02's scoping judgment is right: the "only exceptions" sentence stays true because its enumeration covers archived content and the repairer never touches the archive.
- [x] The class claim (two instances = whole site class) rests on a real sweep with its condition stated; the knowledge package's memory hooks disclose their own writes already.

**Findings:** none Important, none Minor worth recording. **Acceptance: Pass.** No follow-ups.

*Reviewed by review-work action (orchestrated inline, Route A depth; merge range `ec0ebfd..fbc14e8`)*

## Lessons Learned

Route A, no surprises — skipped per the rule (the class sweep confirming the REQ's instance list was complete is already recorded in P-A-U).

## Orientation

Now the docs say what the SessionStart hook actually does to a consumer's repo: banner, reservation reaping, and mechanical stamp repair on queue/working. Leaf doc change; map unchanged.
