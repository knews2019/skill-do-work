---
id: REQ-060
title: "No mechanism resolves a `failed` REQ, so a UR containing one can never close"
status: completed
created_at: 2026-07-29T15:19:48Z
status_changed_at: 2026-07-29T15:34:19Z
claimed_at: 2026-07-29T17:43:39Z
completed_at: 2026-07-29T18:42:07Z
commit: 3a85811
kb_status: pending
domain: general
prime_files: []
tdd: false
depends_on: []
write_set: [actions/abandon.md, actions/work-reference.md, actions/cleanup.md, actions/forensics.md, docs/cleanup-guide.md, actions/work.md]
maintenance: false
review_generated: true
---

# A `failed` REQ has no path to resolution, so its UR is stuck open forever

## What
Every reader that decides whether a User Request (UR) can be closed and moved to `do-work/archive/` now gates on the **terminal-resolved status set** in `actions/work-reference.md`'s Schema Read Contract — `completed`, `completed-with-issues`, `cancelled`. `failed` is deliberately outside that set, so a single `failed` REQ holds its UR open. That part is intended.

The gap is that **nothing in the skill ever moves a REQ out of `failed`.** There is no `failed` → resolved transition anywhere:

- `actions/abandon.md` Step 2 refuses a `failed` target outright ("already terminal; `do-work cleanup` will archive it. Cancelling would erase the failure signal"), so the one user-facing won't-do lever cannot be aimed at it.
- `actions/work.md` Step 8's failure classification archives the original at `status: failed` and spawns a *separate* follow-up REQ. The follow-up completing does not touch the original, which stays `failed` permanently.
- No other action writes a status onto an already-archived REQ.

So the UR waits on a condition that can only be met by hand-editing frontmatter. The user is never told that hand-editing is the exit, and `do-work forensics` will keep reporting the UR as open indefinitely without explaining why.

## Why
The three closure readers were aligned one at a time (`actions/cleanup.md` Pass 1 via REQ-048, `actions/forensics.md` Check 4 via REQ-058, `actions/work.md` Step 8 via REQ-059). Before that alignment, work Step 8 still counted `failed` as finished, so it acted as an accidental release valve — a UR whose last REQ failed would at least get closed by the work loop. Closing the last inconsistency removed the valve, which turns a partial gap into a total one. The fix belongs in the resolution mechanism, not in the closure predicate: the predicate is correct, and reverting it would re-introduce URs that close while real follow-up work is still queued.

Fixing this needs the three closure readers updated in lock-step with whatever mechanism is chosen — `actions/cleanup.md` Pass 1, `actions/forensics.md` Check 4, and `actions/work.md` Step 8's archive table all evaluate the same predicate and must not drift apart again.

## Open Questions
- [x] **How should a `failed` REQ become resolved, so its UR can close?** → Option A — let `do-work abandon` accept a `failed` REQ and cancel it. ("Resolved" here means: counted as done by the UR-closure check, the way `completed` and `cancelled` already are.) Options, each with its trade-off:
  - **A — Let `do-work abandon` accept a `failed` REQ.** Drop `actions/abandon.md`'s refusal and let the user cancel a failed REQ (flipping it to `cancelled`), either whenever a follow-up REQ already exists or whenever the user explicitly waives one. Smallest change and reuses a lever users already know; the cost is that the `failed` status is overwritten, so the board and archive no longer show that this REQ was attempted and failed — the reason survives only in the `## Cancelled` section.
  - **B — Add a `resolved_by:` (or `superseded_by:`) frontmatter field.** The REQ keeps `status: failed`, and a new field naming the follow-up REQ tells the closure readers to treat it as resolved once that follow-up is itself terminally resolved. Preserves the failure signal exactly and makes the link auditable; the cost is a new schema field that all three closure readers plus `tools/queue-kanban/model.go` must learn, and a decision about whether it is set automatically at follow-up completion or by the user.
  - **C — Have the follow-up's completion flip the original.** When a REQ with `addendum_to: REQ-NNN` completes, `actions/work.md` Step 8 also updates REQ-NNN out of `failed`. Fully automatic with no user action; the cost is that a completed addendum does not always mean the original's goal was met (an addendum may narrow scope), so this can close URs on work that was only partially recovered.

  This is your call rather than the builder's because it trades **audit fidelity** (keeping `failed` visible forever) against **queue hygiene** (URs that actually close), and which one matters more depends on how you use the archive.

  **Resolved via `do-work clarify` (2026-07-29T15:34:19Z): Option A — let `do-work abandon` accept a `failed` REQ, cancelling it like any other unwanted item.** Implementation must drop `actions/abandon.md` Step 2's refusal of `failed` targets and update its Step 2 language accordingly; no new frontmatter field is needed.

## Full Context
Discovered during REQ-059's adversarial review. See `do-work/archive/REQ-059-step8-failed-ur-closure-contradiction.md` → `## Discovered Tasks` (the REQ may still be in `do-work/working/` if that run has not archived yet).

**Secondary consideration for whoever implements this:** UR folders archived by older versions of the skill may already contain REQs still at `status: failed` — back then work Step 8 counted `failed` as finished and consolidated the UR anyway. Those legacy folders were never planned for by the current consolidation rules, so the chosen mechanism should either leave already-archived URs alone or state explicitly how they are migrated. Do not let a new closure rule re-open URs that are already sitting in `do-work/archive/`.

---

## Triage

**Route: B** - Medium

**Reasoning:** The mechanism is already decided (Option A, user-resolved via clarify): let `do-work abandon` accept a `failed` target and cancel it. The edit site in `actions/abandon.md` Step 2 is known, but the refusal and the "failed has no exit" claim are restated across the closure readers and schema docs (`actions/work-reference.md`, `actions/cleanup.md`, `actions/forensics.md`, possibly help/docs), so the full restatement-site list needs discovery before scoping — exactly the Closed-Enumerations/restatement trap this repo documents.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Pre-Flight

**Git:** ⚠ 1 uncommitted untracked file (`ai-reports/2026-07-29_1937_UR-007-008-deep-review-batch/index.html` — prior batch's ai-report, unrelated; explicit-file staging keeps it out of this REQ's commit)
**Tests baseline:** ✓ Passing (`bash _dev/tests/contract-regressions.sh`; baseline recorded in `do-work/working/baseline.json`)
**Dependencies:** ✓ N/A (markdown + shell; no package install surface)

*Checked by work action*

## Exploration

Four-angle sweep (`actions/abandon.md` deep-read, the three closure readers, a repo-wide restatement grep, and a board/tooling lock-step check). Key findings that reshape the naive "just delete the Step 2 refusal" plan:

- **The refusal is two gates in Step 1, and the *location* gate fires first.** `actions/abandon.md` Step 1 gates "**Only in archive** → nothing to cancel" (line 34) *before* the "**Status `failed`**" refusal (line 36). Per `actions/work-reference.md`'s Failure Classification, "failed REQs always go to archive root" — so every failed target hits the archive gate and is refused before the status gate is ever reached. Both gates must change; deleting line 36 alone does nothing.
- **The search surface already covers the archive** (`do-work/archive/**/REQ-NNN*.md` is in the Step 1 glob), so abandon does NOT need to widen where it looks — only stop short-circuiting on what it finds.
- **Step 5 ("Archive") is written entirely for queue-resident targets** ("Move each cancelled REQ file out of the queue"). An already-archived target has nothing to move; worse, its "if `archive/UR-NNN/` exists → move it there" bullet would migrate the REQ into a possibly-already-closed UR folder, and its collision guard (`archive/**/REQ-NNN*.md` exists → refuse) would fire against the target's *own* path. The in-archive case must be an explicit in-place no-op that skips the move and the collision guard.
- **`completed_at` collision:** a failed REQ already carries `completed_at` (stamped at failure). Decision recorded below (D-01): re-stamp to the cancellation instant (satisfies the existing STAMPING RULE with zero schema change, keeps the board's `detectCompletionAnomaly` clean) and record the original failure instant + prior status in `## Cancelled`.
- **No board/tooling change — verified, not assumed.** `tools/queue-kanban/model.go`: `failed` is already in `isNeedsInputOrBlockedStatus`, `cancelled` already in `isTerminalResolvedStatus`; the flip just moves an existing card between two existing buckets. No new field, status value, or alias ⇒ the CLAUDE.md lock-step rule does not fire. The board reads `archive/**` (`walk.go`), so an in-place flip is picked up with no move.
- **The closure predicate genuinely does not change.** All three readers (`actions/cleanup.md` Pass 1, `actions/forensics.md` Check 4, `actions/work.md` Step 8) cite the terminal-resolved set by reference (post REQ-048/058/059), so a semantics clarification at the canonical home (`actions/work-reference.md`'s Terminal-resolved set) propagates without touching them — which keeps them in lock-step. Stating the escape hatch once there, not in three copies, is the correct move under the repo's Closed-Enumerations discipline.
- **Stale "no exit" prose to correct:** `actions/abandon.md`'s Common Rationalizations row ("`failed` … holds the UR open"), `actions/work-reference.md`'s Terminal-resolved statement, `docs/cleanup-guide.md` ("holds its UR open until a follow-up resolves it"), and `actions/forensics.md` Check 6's remedy ("failure has no recovery path" → "create a follow-up REQ") + its Output Format sample.

*Generated by Explore agents (workflow)*

## Scope

**Files I will touch:**
- `actions/abandon.md` (modify) — the mechanism: fall-through both Step 1 gates for archived+`failed`; add extra-confirmation arm; Step 3 completed_at re-stamp + explicit error/error_type retention; `## Cancelled` prior-status line; Step 5 in-place branch (skip move + collision guard, never re-open a UR folder); When-to-Use trigger + narrowed exclusion; blockquote second entry point; bare-verb listing includes archived failed REQs; Common Rationalizations row; Red Flags; Verification Checklist
- `actions/work-reference.md` (modify) — Terminal-resolved statement (:199) gains the escape-hatch clause (canonical home; propagates by reference); `error: # Only if failed` comment (:139) widened; cancelled schema block (:141-143) notes error/error_type retention
- `actions/cleanup.md` (modify) — Pass 1 open-UR report line names `do-work abandon` when the blockers are failed REQs (the "user is never told the exit" gap)
- `actions/forensics.md` (modify) — Check 6 remedy text + its Output Format sample name `do-work abandon` as the second exit
- `docs/cleanup-guide.md` (modify) — Pass 1 stale "until a follow-up resolves it" claim
- `actions/work.md` (modify — **scope extension, D-02**) — Step 8 archive-table row (line 585) carries the same stale gloss as the other two closure readers; must be fixed in lock-step

**Files I will NOT touch:**
- The three closure readers' shared predicate sentence (`actions/cleanup.md` Pass 1 line 55, `actions/forensics.md` Check 4 line 72, `actions/work.md` Step 8 line 585) — predicate unchanged; escape hatch stated once at the canonical home keeps them in lock-step
- `tools/queue-kanban/*` — verified no field/status/alias change ⇒ lock-step rule does not fire
- `actions/help.md`, `SKILL.md`, `docs/forensics-guide.md` — status-agnostic / still accurate; word-budget discipline
- forensics Check 4 (a new "UR held open by failed REQ" finding) — feature beyond REQ scope; YAGNI

**Acceptance criteria (restated from REQ):**
- [ ] `do-work abandon REQ-NNN` accepts an already-archived `failed` REQ and flips it to `cancelled` (a terminally-resolved status, so its UR can close)
- [ ] The failure signal survives: `error`/`error_type` retained in frontmatter; `## Cancelled` records the prior `failed` status and original failure instant
- [ ] abandon handles an already-archived target correctly: both Step 1 gates fall through, the flip is in-place, no self-collision, no move into a UR folder
- [ ] Already-archived UR folders are never re-opened or migrated (constraint 4)
- [ ] The three closure readers stay in lock-step (predicate unchanged; escape hatch stated once at the canonical schema home)
- [ ] Stale "failed has no exit / follow-up is the only path" prose is corrected (abandon rationalization row, work-reference:199, docs/cleanup-guide, forensics Check 6 remedy + sample)
- [ ] No board/tooling change (verified against `model.go`)

## Decisions

- [~] **D-01: On the `failed` → `cancelled` flip, re-stamp `completed_at` to the cancellation instant rather than preserving the failure instant.** DECIDE & STATE (reversible, low-reach). Reasoning: the existing STAMPING RULE (`actions/work-reference.md`) already requires every flip to a terminal status to stamp `completed_at`, and `cancelled` is terminal — so re-stamping satisfies the rule with zero schema change and keeps the board's `detectCompletionAnomaly` clean, matching every other abandon path. The original failure instant is preserved in the `## Cancelled` body section (prior-status line), so the audit trail survives the re-stamp. Alternative (preserve the failure instant + write `status_changed_at`) would need a new work-reference carve-out for a terminal→terminal flip — more schema surface for no board benefit (YAGNI).

- [~] **D-02: Extend `write_set` to `actions/work.md` and edit its Step 8 archive-table row (line 585).** DECIDE & STATE (serial run — no co-dispatched REQ claims `work.md`, so the parallel-dispatch disjointness gate does not apply; orchestrator extends the set after confirming that). Surfaced by adversarial review: the correction reversed the "a follow-up must happen / `cancelled` = no-follow-up-wanted" framing, but that exact framing survived as a trailing **gloss** on all three closure readers' predicate sentences (`cleanup.md:55`, `forensics.md:72`, `work.md:585`). Leaving `work.md:585` stale would (a) reintroduce the REQ-059 bug on the pipeline's own archive path, (b) drift the three predicate readers apart, and (c) violate this REQ's own tightened ban at `work-reference.md:203`. The original Scope excluded these three under "predicate unchanged" — correct about the *predicate*, but the trailing gloss is a separate rule re-derivation the correction made false. Fix: delete the stale gloss from all three (keeping the predicate + the by-reference citation each already has), so they defer wholly to the canonical home — which is the lock-step outcome the REQ wanted. Value: eliminates a same-file contradiction and a reintroduced pipeline bug. Risk: editing `work.md` (large file) — bounded to one clause deletion in one table cell; reversible.

<!-- D-XX counter: last used D-02. Next decision: D-03. -->

## Implementation Summary

**Files changed:**
- `actions/abandon.md` (modified)
- `actions/work-reference.md` (modified)
- `actions/cleanup.md` (modified)
- `actions/forensics.md` (modified)
- `docs/cleanup-guide.md` (modified)
- `actions/work.md` (modified — D-02 scope extension)

**What was done:** Implemented Option A — `do-work abandon` now resolves an already-archived `failed` REQ by cancelling it in place, so its UR can reach closure. In `actions/abandon.md`: the Step 1 **location** gate (which fires first because failed REQs always live in `archive/` root) now branches by status — archived-`failed`-at-root falls through to cancellable, archived-`failed`-inside-a-UR-folder is refused (constraint 4: closed UR folders are left untouched), and other archived statuses keep the "nothing to cancel" refusal; the Step 1 status-`failed` row (for the rare queue/working case) likewise became cancellable. Step 2 gained failed-specific consequence wording (per `clear-questions.md`); Step 3 re-stamps `completed_at` to the cancellation instant while retaining `error`/`error_type` and recording the original failure instant + prior status in a new `## Cancelled` **Previously** line; Step 5 gained an in-place branch that skips the move and the self-colliding collision guard for archived targets. Discoverability (the REQ's "user is never told the exit" complaint) added to When-to-Use, the blockquote, the bare-verb listing, an Output Format sample, two reframed Common Rationalizations rows, two Red Flags, two Verification Checklist items, and a Rules blast-radius bound. In `actions/work-reference.md`: the Terminal-resolved statement gained the canonical resolution-**rule** sentence, which the three closure readers cite by reference — their **predicate sentences were not edited**, keeping the predicate in lock-step; the rule explicitly permits a user-facing remedy *pointer* (not a forked definition), which is what `actions/cleanup.md` Pass 1's report line and `actions/forensics.md` Check 6 add. The schema now documents that `error`/`error_type` are retained on a failed→cancelled flip (drift guard so a maintenance pass won't strip the signal). `docs/cleanup-guide.md`'s stale claim corrected. Throughout, the messaging was fixed to the true rule surfaced in review: **a completed follow-up never flips the original out of `failed`; `do-work abandon` is the only transition, needed whether or not a follow-up ran.** Two adversarial-review rounds swept this framing to every instance: the stale "follow-up must happen / `cancelled` = no-follow-up-wanted" **gloss** was removed from all three closure readers' predicate sentences (`cleanup.md` Pass 1 step 2, `forensics.md` Check 4, `work.md` Step 8 — D-02 scope extension), replaced with an identical pointer that defers to the canonical home (keeps them in true lock-step and resolves a same-file contradiction in `cleanup.md`); `forensics.md` Check 6 and abandon's Common Rationalizations row were de-gated so abandon reads as required either way; the canonical statement's reader list was reworded from a closed enumeration to a trigger condition. Step 3 handles a legacy/hand-edited failure with missing `error_type`/`completed_at` by presence (never fabricating a value, single unambiguous treatment for an absent instant), Step 1 accepts a `failed` REQ at `archive/legacy/` as well as root and refuses a REQ duplicated across archive paths, and the no-ID listing surfaces both `root` and `legacy/`. No board/tooling change (verified: no new field, status value, or alias).

## Qualification

Passed — `tools/checks/qualify.sh` OK (mechanical); `tools/checks/scope-drift.sh` OK (Implementation Summary matches Scope exactly). Judgment: 5 files verified as substantive prose modifications; every acceptance criterion traced to a change; no data-flow surface (instruction files). No debug artifacts.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Passed (exit 0) — matches the Step 5.75 baseline. Covers SKILL.md word budgets, shipped-file citation rules (no CLAUDE.md/AGENTS.md citations), and the Common-Rationalizations do-work-noun requirement (the reframed abandon rows contain REQ/UR/queue nouns).

Non-behavioral change (documentation/instruction files) — no red-green pair; regression evidence is the contract suite holding at baseline.

*Verified by work action*

## Review

Route-B review run as an adversarial multi-agent workflow (4 dimensions — requirements, correctness/logic-traps, restatement-sweep/lock-step, consistency/tooling — each Important+ finding verified by an independent skeptic), then a focused re-verification after remediation, then a final single-agent cold-read gate.

**Round 1 — Acceptance: Partial (all 4 dimensions).** Core mechanism traced correctly and the no-tooling-change claim verified against `tools/queue-kanban/model.go`. Confirmed defects, all remediated:
- **Self-contradiction:** the canonical sentence forbade closure readers from "restating the exit inline" while the same change added remedy-pointers to `cleanup.md` Pass 1 and `forensics.md` Check 6. → Scoped the ban to the predicate/set/rule; explicitly permit a remedy *pointer*.
- **Correctness (important):** `docs/cleanup-guide.md` and `cleanup.md`'s report line implied a completed follow-up closes the UR — mechanically impossible (only `abandon` flips the original out of `failed`), and `cleanup.md` said "retry the REQ" (no such op). → Rewritten to the true rule.
- **Robustness (minor):** Step 3 assumed the failed REQ always carries `completed_at`/`error_type`. → Made presence-conditional with an explicit no-fabrication rule.
- **Closed-enum gap / over-promise (minor):** Step 1 missed `archive/legacy/`; Check 6 over-promised for UR-nested failures abandon refuses. → Both scoped.

**Round 2 — re-verification.** All 6 round-1 findings confirmed resolved, but the re-review caught that the *same* corrected framing survived in four more places I'd fixed *around* rather than swept: the stale gloss on all three closure readers' predicate sentences (`cleanup:55`, `forensics:72`, **`work.md:585`** — the latter reintroducing the REQ-059 bug on the pipeline's own archive path), `forensics.md` Check 6's exit framing, and abandon's Common Rationalizations row. → Swept every instance (grep-verified gone), removed the gloss from all three readers in lock-step (D-02 scope extension to `work.md`), de-gated Check 6 and the rationalization row, and reworded the canonical reader list from a closed enumeration to a trigger condition.

**Final gate — cold read: clean on the governing rule.** Zero grep hits for the stale framing; the three closure-reader glosses byte-identical; all cross-references resolve; contract-regressions pass. Residual nits (pre-existing set-restatement in `cleanup:55`/`abandon:98`; `help.md`/`next-steps.md` not mentioning the failed-resolution use) are out of REQ-060's scope and left as-is per YAGNI — no follow-up REQs warranted.

**Acceptance: Pass.** Archived as `completed`.

## Lessons Learned

**What worked:** Adversarial multi-agent review with independent skeptic-verification earned its cost here — it caught a genuine self-contradiction (the ban vs. its own additions) and a mechanically-impossible claim (a follow-up "closing" a UR) that a single-pass review would likely have missed. Reading the *fully-edited* file as a cold reader (not just the diff) is what surfaced the same-file `cleanup:55`-vs-`cleanup:63` contradiction.

**What didn't:** Spot-fixing the follow-up/cancel framing site-by-site instead of sweeping the primitive first. The correction reversed a semantic that turned out to be restated in ~9 places; round 1 hit 5, round 2 had to hit the other 4. This is exactly the CLAUDE.md "grep the same primitive across all actions before calling it fixed" lesson — the sweep should have been the *first* implementation move, not a remediation.

**Worth knowing:** The load-bearing, non-obvious fact this whole REQ turns on: **a completed follow-up REQ never flips its parent out of `failed`** — `do-work abandon` is the only transition, needed whether or not a follow-up ran. Any future edit near failure/closure semantics must preserve that. Also: abandon's Step 1 gate is ordered *location-first* (the "only in archive" branch decides the failed path before the status rows), because failure classification always sends failures to `archive/` root or `legacy/` — never assume a status row is reached for an archived target.

## Orientation

`do-work abandon` now resolves a `failed` REQ (not just pending/blocked ones): it cancels an already-archived failure in place, flipping it to `cancelled` so its UR can finally close, while preserving the failure record. Lives in the abandon action; the closure predicate (`work-reference.md`'s Terminal-resolved set, read by cleanup/forensics/work) is unchanged — this fills the missing `failed`→resolved transition those readers were gating on. No new schema field, no board change. `[MAP CHANGED]` — a status that was a permanent dead-end is now user-resolvable, and the canonical failure-resolution rule now lives in one place (`work-reference.md`'s Terminal-resolved statement) with the three closure readers deferring to it by reference.
