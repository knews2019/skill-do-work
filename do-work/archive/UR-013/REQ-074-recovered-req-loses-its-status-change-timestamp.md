---
id: REQ-074
title: A recovered REQ loses the timestamp that says when it was reset
status: completed
claimed_at: 2026-08-03T15:52:32Z
completed_at: 2026-08-03T15:56:18Z
commit: 80e7b88
kb_status: promoted
kb_entry: REQ-074-recovered-req-loses-its-status-change-ti.md
route: A
created_at: 2026-08-03T14:48:01Z
status_changed_at: 2026-08-03T15:45:29Z
user_request: UR-013
addendum_to: REQ-071
domain: general
prime_files: []
tdd: false
depends_on: []
maintenance: false
review_generated: false
discovered_during: REQ-071
---

# A Recovered REQ Loses the Timestamp That Says When It Was Reset

## What

Two places in the skill reset a stuck REQ back to `pending`, and they disagree about whether to
record *when* the reset happened.

- When a **human** does it by hand, `actions/forensics.md` Check 1 tells them to stamp
  `status_changed_at` with the current time.
- When the **pipeline** does it automatically, `actions/work-reference.md`'s Crash Recovery substep 1
  resets the status and stamps nothing.

`status_changed_at` is the field the Kanban board reads to show how long a card has been sitting in
its current state. Without it, an automatically-recovered REQ falls back to its original creation
date, so the board shows a REQ that was reset two minutes ago as though it had been waiting since the
day it was written.

## Why This Is Worth Deciding Rather Than Just Fixing

The one-line fix is obvious (have Crash Recovery stamp the field too), which is exactly why it is
worth pausing: the field's own schema note says it is written **on any status flip that has no
dedicated timestamp of its own**, and lists the manual reset as one of the writers. By that rule
Crash Recovery has been out of compliance since the field was introduced — so the question is not
really "should we stamp it here" but "is the rule right, or is the omission deliberate?" Getting that
backwards writes a rule into a second place and makes the disagreement permanent instead of fixing it.

Nothing is broken today beyond a misleading age on one board card. There is no data loss and no
pipeline behavior depends on the field — the schema calls it display-only.

## Context

- `actions/work-reference.md` → **Crash Recovery (Step 1)**, substep 1 — resets `status`, removes
  `claimed_at` and `route`, stamps nothing.
- `actions/forensics.md` → Check 1 (Stuck Work), suggested remediation — the manual reset, which does
  stamp `status_changed_at`.
- `actions/work-reference.md` → **Request File Schema — Full Frontmatter**, the `status_changed_at`
  line — states the trigger condition ("any status flip that has no dedicated `*_at` stamp of its
  own") and names its writers illustratively, including the manual reset.
- `tools/queue-kanban/model.go` — the consumer: the board's state timer prefers `status_changed_at`
  over `created_at` and file mtime for pending-tier cards.

## Discovered During

REQ-071 (`do-work/archive/` — crash recovery respects a live claim). Found by that REQ's restatement
sweep while checking that the automatic and manual reset procedures agreed with each other. It is
pre-existing and was outside REQ-071's requirements, so it was filed rather than swept in.

## Open Questions

- [x] Crash Recovery resets a REQ to `pending` without recording when — should it stamp
  `status_changed_at` like the manual reset does? → **Confirmed: Yes — stamp it in Crash Recovery
  substep 1** (option (a); user, via `do-work clarify`, 2026-08-03)
  Recommended: **Yes — stamp it in Crash Recovery substep 1.** One line in
  `actions/work-reference.md`; brings the automatic path in line with both the field's stated trigger
  condition and the manual path, and the board stops dating a just-recovered REQ from its creation day.
  Value: the board's "waiting for N days" figure becomes true for recovered REQs, and the two reset
  procedures stop contradicting each other.
  Risk: very low and reversible — the field is display-only by contract, nothing in the pipeline
  reads it, and no other REQ's data changes.
  Also: **(b) Leave Crash Recovery as it is and narrow the field's rule instead** — amend the
  `status_changed_at` schema line to say recovery resets are deliberately excluded, and say why. Pick
  this if the omission was intentional and the board should keep showing the original queue age.
  Also: **(c) Do nothing.** The disagreement stays, and the next person to notice it re-discovers it
  from scratch. Cheapest now, and the reason this REQ exists.

## Answer Record

**[2026-08-03] Option (a) confirmed by the user via `do-work clarify`:** Crash Recovery substep 1
stamps `status_changed_at`. The field's stated trigger condition is right and the omission was not
deliberate — so the fix goes in the automatic path, not into a narrowed rule.

Two implementation points the answer settled, both raised during clarify:

1. **Stamp both flip outcomes, not just `pending`.** Substep 1 can land a recovered REQ in `pending`
   *or* `pending-answers` (the latter when an unresolved `- [ ]` survives in `## Open Questions`).
   Both are status flips with no dedicated `*_at` stamp, so both stamp. The third outcome —
   `status: blocked` preserved with its `blocked_by` — is **not** a flip and must **not** stamp.
2. **Substep 1 also removes `claimed_at`**, so with nothing stamped there is no trace at all of the
   recovery instant and the board falls all the way back to `created_at`. That is the concrete harm,
   and it makes the stamp the only surviving record of when recovery happened — the same reasoning
   the schema note already gives for the unblock flips that remove `blocked_at`.

Out of scope: no change to the schema's `status_changed_at` note (its trigger condition already
covers this — no second copy of the rule), no change to `actions/forensics.md` Check 1 (already
correct), and no change to `tools/queue-kanban/` (the field is display-only and the reader is right).

## Full Context

See `do-work/user-requests/UR-013/input.md` for the original batch input, and REQ-071 for the sweep
that surfaced this.

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact file and substep (`actions/work-reference.md` → Crash Recovery
substep 1), the user's clarify answer settled which of the three options to take, and the
`## Answer Record` pins the two implementation details (stamp both `pending` and `pending-answers`
flips; never stamp a preserved `blocked`). Prose-only change to one file.

**Planning:** Not required

## Plan

**Skipped** — Route A. Planning not required for a single-file prose change with a settled decision.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — Crash Recovery (Step 1) substep 1 now prescribes stamping
  `status_changed_at` on the reset, on both flip outcomes, with the preserved-`blocked` case
  explicitly excluded.
- `_dev/tests/contract-regressions.sh` (modified) — one assertion pinning the stamp inside the
  existing `crash_recovery_block` extraction (D-01).

**What was done:** Added the stamp instruction to Crash Recovery substep 1, immediately after the
`claimed_at`/`route` removal it depends on for its rationale. The wording cites the Timestamp rule by
its section name rather than restating the format, names both flip targets (`pending` and
`pending-answers`) so the `## Open Questions` branch above is covered too, states why the stamp is the
only surviving trace (this substep removes `claimed_at`), and carves out the preserved-`blocked`
exception, whose `blocked_at` is intact and whose status does not flip. No change to the
`status_changed_at` schema note, `actions/forensics.md`, or `tools/queue-kanban/` — all three were
already correct.

## Decisions

**D-01 — Added a contract-suite assertion pinning the stamp. DECIDE & STATE.**
The REQ asked for one line of prose and nothing more, so this is scope the builder added. Reasoning:
the defect being fixed is not "the stamp is missing" but "the stamp was missing for the field's entire
lifetime and nobody noticed" — the manual path stamped, the automatic one didn't, and only a
restatement sweep on an unrelated REQ surfaced it. A prose-only fix leaves the next editor free to drop
the sentence with nothing failing. The suite already extracts `crash_recovery_block` for REQ-071's six
assertions, so this is one `assert_block_contains` in an existing pattern — reversible, no new
machinery. Verified RED by stashing `actions/work-reference.md`: the assertion fails, naming REQ-074.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Contract regression checks passed

**Red-green validation:**
- `contract-regressions.sh` — Crash Recovery stamp assertion (D-01): ✗ FAILs with
  `actions/work-reference.md` stashed → ✓ passes with the fix restored. Run as an explicit stash /
  run / unstash cycle, not inferred.

**Restatement check:** `grep -rn "remove[sd]* \`claimed_at\`" actions/ docs/ crew-members/ tools/queue-kanban/*.md _dev/`
returns exactly three sites — this substep, `actions/forensics.md` Check 1 (the manual reset, already
stamping), and `actions/work.md`'s mid-run blocked flip (a different procedure that stamps its own
dedicated `blocked_at`). No fourth copy of the reset procedure needed updating.

The board-side reader (`tools/queue-kanban/model.go`) already prefers `status_changed_at` and has its
own tests (`state_timer_source_test.go`, `future_timestamp_test.go`) — unchanged here.

## Qualification

Passed — 2 files verified, 3 requirements traced (the stamp, both-flip coverage, blocked exclusion),
mechanical checks green via `tools/checks/qualify.sh`. Scope-drift comparison correctly skipped
(Route A, no `## Scope`).

## Review

**Overall: 95%** | 2026-08-03T15:55:45Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor

- **Minor — substep 1 is now ~380 words carrying five distinct rules** (flip target, `blocked`
  exception, `claimed_at`/`route` removal, the new stamp, the conditional `write_set` clear). It was
  already this file's densest paragraph and this addition does not improve it. Not worth a follow-up on
  its own; worth splitting the next time that substep is edited for another reason.

**Restatement sweep:** Run. The diff adds a *writer* of `status_changed_at`, not a new meaning — the
schema note states its trigger condition and names writers illustratively, so it needed no edit
(REQ's own out-of-scope list agrees). Consumers verified: `actions/forensics.md` Check 1 already
stamps; `actions/work.md`'s mid-run blocked flip is a different procedure with its own `blocked_at`;
`tools/queue-kanban/model.go` reads the field with unchanged semantics. No stale restatement found.

**Acceptance:** Pass — the contract assertion proves the prescription is present, and was shown to
fail without it.
**Suggested testing:** 1 item — the stamp is prose an agent must follow, so the only end-to-end proof
is a real crash-recovery run (kill a claimed REQ mid-build, re-run `do-work run`, confirm the
re-queued file carries a fresh `status_changed_at` and the board shows minutes rather than days).
Same class of untested-by-construction as REQ-073's two-builder test.
**Follow-ups created:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Stashing the source file to prove the new contract assertion actually fires. A grep
assertion that has only ever been run against a tree where it passes is untested — this repo already
shipped one that aborted the suite silently in exactly the case it existed to catch (REQ-073's first
invariant check), so the stash/run/unstash cycle is cheap insurance.

**What didn't:** Nothing failed, but the first instinct — "one line, no test needed" — was wrong for
the same reason the REQ existed: the omission survived because nothing failed when it was absent.

**Worth knowing:** `status_changed_at` now has three writers with a shared rationale (the manual reset,
the unblock flips, and crash recovery) and the common thread is that each one *removes* the stamp that
would otherwise date the transition — `claimed_at` here, `blocked_at` for the unblock flips. That is the
real reason the field exists, and it is a better trigger test than the field list: **if the flip
discards its own `*_at`, it must stamp `status_changed_at`.**

## Orientation

Crash recovery now dates its own resets. A REQ that the pipeline re-queues after a crash carries the
instant it was recovered, so the board's waiting-time figure counts from the reset instead of from the
day the request was written — matching what the by-hand reset in forensics already did. Lives in the
work pipeline's recovery path (`actions/work-reference.md`); no code, no schema change, and the board
reader was already correct. `prime_files` is empty for this REQ, so no prime staleness check applied.
