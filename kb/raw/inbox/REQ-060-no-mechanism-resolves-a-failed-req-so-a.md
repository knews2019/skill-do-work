---
source_type: req_lesson
req_id: REQ-060
req_path: do-work/archive/legacy/REQ-060-failed-req-resolution-path.md
date: 2026-07-29
domain: general
module: actions
tags: [actions, mechanism, resolves, failed, containing]
---

# Lessons from REQ-060: No mechanism resolves a `failed` REQ, so a UR containing one can never close

## What the REQ was about

Every reader that decides whether a User Request (UR) can be closed and moved to `do-work/archive/` now gates on the **terminal-resolved status set** in `actions/work-reference.md`'s Schema Read Contract — `completed`, `completed-with-issues`, `cancelled`. `failed` is deliberately outside that set, so a single `failed` REQ holds its UR open. That part is intended.

## Solution summary

Implemented Option A — `do-work abandon` now resolves an already-archived `failed` REQ by cancelling it in place, so its UR can reach closure. In `actions/abandon.md`: the Step 1 **location** gate (which fires first because failed REQs always live in `archive/` root) now branches by status — archived-`failed`-at-root falls through to cancellable, archived-`failed`-inside-a-UR-folder is refused (constraint 4: closed UR folders are left untouched), and other archived statuses keep the "nothing to cancel" refusal; the Step 1 status-`failed` row (for the rare queue/working case) likewise became cancellable. Step 2 gained failed-specific consequence wording (per `clear-questions.md`); Step 3 re-stamps `completed_at` to the cancellation instant while retaining `error`/`error_type` and recording the original failure instant + prior status in a new `## Cancelled` **Previously** line; Step 5 gained an in-place branch that skips the move and the self-colliding collision guard for archived targets. Discoverability (the REQ's "user is never told the exit" complaint) added to When-to-Use, the blockquote, the bare-verb listing, an Output Format sample, two reframed Common Rationalizations rows, two Red Flags, two Verification Checklist items, and a Rules blast-radius bound. In `actions/work-reference.md`: the Terminal-resolved statement gained the canonical resolution-**rule** sentence, which the three closure readers cite by reference — their **predicate sentences were not edited**, keeping the predicate in lock-step; the rule explicitly permits a user-facing remedy *pointer* (not a forked definition), which is what `actions/cleanup.md` Pass 1's report line and `actions/forensics.md` Check 6 add. The schema now documents that `error`/`error_type` are retained on a failed→cancelled flip (drift guard so a maintenance pass won't strip the signal). `docs/cleanup-guide.md`'s stale claim corrected. Throughout, the messaging was fixed to the true rule surfaced in review: **a completed follow-up never flips the original out of `failed`; `do-work abandon` is the only transition, needed whether or not a follow-up ran.** Two adversarial-review rounds swept this framing to every instance: the stale "follow-up must happen / `cancelled` = no-follow-up-wanted" **gloss** was removed from all three closure readers' predicate sentences (`cleanup.md` Pass 1 step 2, `forensics.md` Check 4, `work.md` Step 8 — D-02 scope extension), replaced with an identical pointer that defers to the canonical home (keeps them in true lock-step and resolves a same-file contradiction in `cleanup.md`); `forensics.md` Check 6 and abandon's Common Rationalizations row were de-gated so abandon reads as required either way; the canonical statement's reader list was reworded from a closed enumeration to a trigger condition. Step 3 handles a legacy/hand-edited failure with missing `error_type`/`completed_at` by presence (never fabricating a value, single unambiguous treatment for an absent instant), Step 1 accepts a `failed` REQ at `archive/legacy/` as well as root and refuses a REQ duplicated across archive paths, and the no-ID listing surfaces both `root` and `legacy/`. No board/tooling change (verified: no new field, status value, or alias).

## What worked

- Adversarial multi-agent review with independent skeptic-verification earned its cost here — it caught a genuine self-contradiction (the ban vs. its own additions) and a mechanically-impossible claim (a follow-up "closing" a UR) that a single-pass review would likely have missed. Reading the *fully-edited* file as a cold reader (not just the diff) is what surfaced the same-file `cleanup:55`-vs-`cleanup:63` contradiction.

## What didn't work

- Spot-fixing the follow-up/cancel framing site-by-site instead of sweeping the primitive first. The correction reversed a semantic that turned out to be restated in ~9 places; round 1 hit 5, round 2 had to hit the other 4. This is exactly the CLAUDE.md "grep the same primitive across all actions before calling it fixed" lesson — the sweep should have been the *first* implementation move, not a remediation.

## Worth knowing

- The load-bearing, non-obvious fact this whole REQ turns on: **a completed follow-up REQ never flips its parent out of `failed`** — `do-work abandon` is the only transition, needed whether or not a follow-up ran. Any future edit near failure/closure semantics must preserve that. Also: abandon's Step 1 gate is ordered *location-first* (the "only in archive" branch decides the failed path before the status rows), because failure classification always sends failures to `archive/` root or `legacy/` — never assume a status row is reached for an archived target.

## Back-reference

See `do-work/archive/legacy/REQ-060-failed-req-resolution-path.md` for the full REQ — triage, implementation, review, and lessons. Commit `3a85811`.
