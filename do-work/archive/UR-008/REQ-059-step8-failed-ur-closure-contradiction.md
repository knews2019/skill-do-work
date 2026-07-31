---
id: REQ-059
title: "work.md Step 8 archive table counts `failed` as finished for UR closure, contradicting work-reference.md's terminal-resolved set"
status: completed
claimed_at: 2026-07-29T15:08:48Z
commit: d414590
route: B
created_at: 2026-07-29T13:31:00Z
status_changed_at: 2026-07-29T15:25:13Z
user_request: UR-008
addendum_to: REQ-048
domain: general
prime_files: []
tdd: false
depends_on: []
write_set: [actions/work.md, actions/work-reference.md, docs/cleanup-guide.md, docs/forensics-guide.md]
maintenance: false
review_generated: false
---

# work.md Step 8 vs work-reference.md disagree on whether `failed` closes a UR

## What
Discovered during REQ-048. REQ-048 made cleanup Pass 1 treat `failed` as holding a UR open (per `actions/work-reference.md`'s terminal-resolved status set, which excludes `failed`). But `actions/work.md` Step 8's archive table (~line 583) still lists `failed` among the statuses that count as "finished" for UR closure (`completed`, `completed-with-issues`, `cancelled`, or `failed`). So after REQ-048, cleanup Pass 1 holds a UR open on a `failed` REQ while work Step 8 would close it — the two readers now agree on the key (`user_request`) and the locations but **not** on the `failed` status. This is the residual half of REQ-048's "make both readers agree" goal, in a file REQ-048 could not touch (`actions/work.md` was out of its write_set). Narrow in practice (Step 8's failure classification usually spawns a `pending` follow-up that holds the UR open anyway) but real when no follow-up is created (e.g. cycle detection skips it).

## Open Questions
- [x] Which reader is correct on `failed` — does a `failed` REQ **hold its UR open** (work-reference.md's terminal-resolved set) or **close it** (work.md Step 8's table)? → Confirmed: a `failed` REQ **holds its UR open**. Align work.md Step 8's table to work-reference.md — drop `failed` from the finished set and cite the terminal-resolved set, so a UR with an unresolved `failed` REQ stays open until a follow-up resolves it. (Do NOT change work-reference.md's terminal-resolved set — the user chose the smaller-blast-radius fix.)

## Full Context
Discovered task from REQ-048. See `do-work/archive/REQ-048-ur-closure-keying-consistency.md` → `## Discovered Tasks`.

## Implementation Summary

Step 8's archive table no longer enumerates its own "finished" set. The `user_request: UR-NNN` row now checks whether every REQ in the UR is **terminally resolved** and cites `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set as canonical, with the same "any status outside it holds the UR open, **`failed` included**" sentence that `actions/cleanup.md` Pass 1 and `actions/forensics.md` Check 4 already carry. The dead clause about failed REQs staying at archive root during consolidation went with it — under the corrected rule a UR can never reach all-resolved while a `failed` REQ is present — and `failed` now leads the row's list of statuses that hold the UR open. All three readers (work Step 8, cleanup Pass 1, forensics Check 4) now key UR closure on one canonical set from one file, so a UR with an unresolved `failed` REQ stays open until a follow-up resolves it or it is explicitly cancelled.

**Files changed:**
- `actions/work.md` — Step 8 archive table's `user_request` row: cite work-reference.md's terminal-resolved set instead of a locally-enumerated finished set that included `failed`.
- `actions/work-reference.md` — Terminal-resolved status set paragraph: marked its caller list illustrative rather than exhaustive and added `actions/forensics.md` Check 4 (the set's membership and semantics are unchanged).
- `docs/cleanup-guide.md` — Pass 1 description: closure predicate restated as terminally-resolved instead of "all REQs archived".
- `docs/forensics-guide.md` — Orphaned URs check row: same predicate correction.

## Discovered Tasks

- **[normal] Nothing in the skill can resolve a `failed` REQ, so a UR containing one now has no closure path at all.** With Step 8 aligned, all three closure readers (`actions/cleanup.md` Pass 1, `actions/forensics.md` Check 4, `actions/work.md` Step 8) hold a UR open on a `failed` REQ — but no `failed` → resolved transition exists anywhere in the skill. `actions/abandon.md` Step 2 (line ~36) refuses a `failed` target outright ("already terminal … Cancelling would erase the failure signal"), and `actions/work.md` Step 8's Failure Classification archives the original at `status: failed` and spawns a *separate* follow-up whose completion never touches the original. The gap pre-existed via cleanup Pass 1 (REQ-048) and forensics Check 4 (REQ-058); before this REQ, work Step 8 was the accidental release valve, so closing the last inconsistency turned a partial gap into a total one. The remedy (permit abandon on a failed REQ, add a `resolved_by:`/`superseded_by:` field, or have the follow-up's completion flip the original) lives in `actions/abandon.md` and `actions/work-reference.md` — outside this REQ's chosen blast radius, which the user deliberately limited to the smaller fix. Secondary consideration: UR folders archived by older skill versions may already hold still-`failed` REQs that consolidation-era rules never planned for, so any fix must say whether those are migrated or left alone. Queued as `do-work/queue/REQ-060-failed-req-resolution-path.md` (`status: pending-answers` — the mechanism is a user decision trading audit fidelity against queue hygiene).

## Review

Adversarial review workflow (4 Opus lenses -> 2 diverse refuters per Important+ finding; 12 agents): verdict FIX-THEN-PASS -> fixes applied -> PASS.

- Upheld, routed out of this commit by its own analysis: with Step 8 aligned, no mechanism anywhere resolves a `failed` REQ (abandon.md refuses; no failed->resolved transition exists), so a UR holding one has no closure path — pre-existing gap made total; recorded in `## Discovered Tasks` and queued as `do-work/queue/REQ-060-failed-req-resolution-path.md` (pending-answers; mechanism is a user decision).
- Adopted minors: work-reference.md's terminal-resolved caller list was a stale closed enumeration (now illustrative + forensics Check 4 added; set membership untouched); docs/cleanup-guide.md + docs/forensics-guide.md still stated "all REQs archived" as the closure rule (corrected).
- Skipped (recorded, deliberate): "third hardcoded copy of failed-outside-the-set" nit — the emphasis is intentionally mirrored across the three readers; cleanup-Pass-2 legacy-folder concern folded into REQ-060's considerations.
- 3 sibling findings killed by refutation (same systemic gap claimed as an in-commit defect; the recorded spec explicitly declines that scope).
