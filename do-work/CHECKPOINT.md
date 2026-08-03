---
session_ended: 2026-08-03T16:15:28Z
last_completed: REQ-076
queue_state: 0 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 3
session_depth: moderate
---

# Session Checkpoint

## Completed This Session

- REQ-074: Crash recovery stamps `status_changed_at` on the reset (Route A, 95%) — v0.166.1, commit `80e7b88`
- REQ-075: `write_set`'s display-only reason no longer rests on a falsified premise (Route B, 91%) — v0.166.2, commit `738e9fe`
- REQ-076: Go utility emits the canonical UTC timestamp (Route C, 92%) — v0.167.0, commit `3fb5938`

Queue is **empty**: no pending, no pending-answers, no blocked, nothing claimed. **UR-013 and UR-014
are both fully consolidated** into `do-work/archive/UR-013/` and `do-work/archive/UR-014/`, so
`do-work/user-requests/` is empty for the first time in this batch's life. `actions/version.md` and the
top `CHANGELOG.md` entry both read 0.167.0; `queue-kanban verify --repo-root .` exits 0.

The `## In Progress (interrupted)` section is deliberately **absent** — nothing was interrupted. A REQ
found in `do-work/working/` by a future session is therefore a foreign claim and must be left
byte-identical (`actions/work-reference.md` → Crash Recovery (Step 1)).

## Still Queued

Nothing. `do-work/working/baseline.json` is this run's pre-flight record, not a claim.

## Session Notes

- **REQ-076 shipped a subcommand usable immediately.** `queue-kanban now` prints the exact `*_at` stamp
  shape; the Timestamp rule prefers it when the binary is already built, with `date -u` still the floor.
  This checkpoint's `session_ended` was produced by it.
- **Deferred, carried into the next session:**
  - All three REQs carry `kb_status: pending`. `kb/` exists; the handoff was batched rather than
    interrupting the run three times, and was raised with the user at hand-back.
  - **REQ-073's live two-builder acceptance test still has not run** (carried from the previous
    checkpoint, now one batch older). Everything since has been built serially. Procedure is in REQ-073's
    `## Review` under Suggested additional testing.
  - **REQ-076's Windows fallback is reasoned, not tested.** `Get-Date -AsUTC` needs PowerShell 7+; 5.1
    would need `(Get-Date).ToUniversalTime().ToString(...)`. Nothing here can run either.
- **A sweep REQ's site inventory is a floor, not a list.** REQ-075 was captured naming five stale sites;
  the assertion written first found eight, and a second grep shape found three more — eleven total,
  including the `write_set` schema line itself and the browser tooltip. Write the check before trusting
  the prose inventory: on a sweep REQ the check *is* the inventory.
- **A retired premise leaves two fingerprints.** The strong form names what it claimed ("one REQ at a
  time"); the weak form names what it was called ("under the exclusive-session model"). Sweeping for one
  reads as clean. The weak form is more dangerous, because the thing it names is *still true* — only its
  relevance died — so the sentence survives inspection. Sweep for both.
- **Line-granularity greps are blind to wrapped comments; proximity windows false-positive on canonical
  text.** REQ-075 hit both: a line sweep passed over `model.go` and `board.js` entirely (their comments
  wrap), and a 3-line window then flagged the canonical Fan-Out section, where "integration runs one REQ
  at a time" is true. The working shape is per-class — line sweep for prose, file-level "must not name a
  builder count" negative for comment-carrying source.
- **Adding a *preferred* source to a rule does not make it used.** Ten of the eleven sites citing the
  Timestamp rule inline `date -u` as the command, so an agent following them never reaches the new
  preference order — and a Windows agent still gets the broken command. Recorded as a REQ-076 Minor
  finding; the fix serving both goals is subtraction (strip the inline command from the citing sites so
  one place says how), not another copy of the order.
- **`time.RFC3339` is right to parse with and wrong to emit with.** It writes a numeric offset for a
  non-UTC instant and can carry sub-second digits — neither of which the schema's `*_at` shape accepts.
  REQ-076 D-02.
