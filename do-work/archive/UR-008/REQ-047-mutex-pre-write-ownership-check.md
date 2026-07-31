---
id: REQ-047
title: "Orchestrator-lock mutex: re-verify ownership immediately before the write so an evicted slow owner cannot clobber its successor"
status: completed
route: A
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T13:14:44Z
completed_at: 2026-07-29T13:23:14Z
commit: 5b2c02f
user_request: UR-008
domain: general
prime_files: []
tdd: false
depends_on: []
related: [REQ-044]
batch: deep-review-followups
write_set: [actions/work-reference.md]
maintenance: false
---

# Mutex Pre-Write Ownership Check

## What

The serialized-lock mutex reclaims a mutex directory whose mtime is over one minute old as having "a verifiably dead owner" — but the same code block admits the critical section "can legitimately span a model round-trip," which can exceed 60 seconds. An evicted slow-but-live owner's in-flight `printf + mv` still lands, interleaving with the successor's read-modify-write: last writer wins, a claim is lost, and Crash Recovery re-queues a live REQ — the exact failure the mutex exists to prevent. (Confirmed against the current post-d839cf5 code: the owner-token check on *release* prevents cascade-deleting the successor's mutex, but nothing guards the lock-file *write* itself.)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/general.md`, `crew-members/coding-guardrails.md`, and CLAUDE.md's "Prescribed Shell Commands Must Surface What the Steps Consume". Approach: in the `Serialized lock updates` block, split the `printf … && mv … || rm` one-liner so a re-check of `$mutex_path/owner` against `$session_id` sits between staging the temp file and the publishing `mv`; on mismatch discard the temp file, publish nothing, leave the mutex to its new owner, and re-enter the block from the top. Reword the mtime comment to claim age rather than death and point at the new re-check as what makes eviction safe. No change to `-mmin -1`, no mtime refreshing.
- [x] **[APPLY]:** One file touched (`actions/work-reference.md`), one section (`## Concurrent-Orchestrator Lock Guard (Step 1)`): the mtime-reclaim comment, the staging/publish lines of the prescribed block, and the one paragraph below it that repeated the softened claim. No other action, crew, or tool file changed.
- [x] **[UNIFY]:** `git diff --stat` → `actions/work-reference.md | 40 ++++++--------` (32 insertions, 8 deletions); `git status --porcelain -uall` → `M actions/work-reference.md` only. Extracted the edited prescribed block verbatim to a scratch file (73 lines, `lock_path=` through the closing release `fi`): `bash -n` PASS and `sh -n` PASS (POSIX-clean). Ran the block's critical-section half against a scratch queue twice: (a) **eviction** — mutex owner token pre-set to `SESSION-B`, run as `SESSION-A` → exit 3, stderr eviction notice, `orchestrator-lock.json` still holds B's image byte-for-byte, zero `.tmp` files left, mutex still owned by `SESSION-B` (release correctly skipped); (b) **happy path** — owner token `SESSION-A`, run as `SESSION-A` → exit 0, new image published, zero `.tmp` files, mutex released. `_dev/tests/contract-regressions.sh` → "Contract regression checks passed." Grepped `mutex_path`/`lock_tmp` and `verifiably dead`/`dead owner` repo-wide (excluding `.git`, `do-work`, `build`): the mutex prescription exists in this file only, so the fix has no copy-paste siblings, and the softened claim survives nowhere else. No debug artifacts, no `echo`/`set -x` scaffolding, no scratch paths in the diff (all simulation files live in the session scratchpad).

## Requirements

1. In the prescribed mutex block (`actions/work-reference.md` → Serialized lock updates), immediately before the `mv` of the temp file onto `$lock_path`, re-verify `$(cat "$mutex_path/owner")` equals this session's `$session_id`. On mismatch: discard the temp file, write nothing, and re-enter the acquisition loop (the fresh read inside the new critical section re-judges everything). This shrinks the lost-update window from model-round-trip scale to milliseconds without touching the one-minute reclaim bound.
2. Soften the "verifiably dead owner" comment to what the mtime actually proves (age, not death) — the pre-write check is what makes eviction safe, so the comment should point at it.

## Constraints

- Do not shorten or lengthen the one-minute reclaim bound, and do not add mtime refreshing — the fixed-mtime property is load-bearing for reclaim of genuinely dead owners.
- Keep the block a single copy-pasteable prescription; the check must not reference shell state from an earlier block.

## Red-Green Proof
**RED prompt/case:** Walk the prescribed block with two sessions: A acquires, stalls >60s mid-critical-section (model round-trip); B reclaims at the mtime check, acquires, reads the lock; A's `printf + mv` then completes — A's stale image overwrites B's state (or vice versa), and nothing in the current block detects it.
**Why RED now:** Ownership is checked only at release (to protect the successor's *mutex*), not before the lock-file *write*.
**GREEN when:** The block re-checks the owner token immediately before `mv` and aborts-and-retries on mismatch; the walk-through above ends with A discarding its write; the comment no longer claims mtime verifies death.
**Validation:** User confirmed (approved capture of the adjudicated external finding)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** A (Direct to Builder)
**Reasoning:** A single-file (`actions/work-reference.md`), precisely-specified change to one named prescribed block (the serialized-lock mutex under "Serialized lock updates"): add an owner re-check immediately before the `mv`, plus soften one comment. No location discovery or pattern-matching needed — the REQ names the exact block and the exact edit. Route A.
**Complexity indicators:** Concurrency-critical prescribed shell — small in scope but unforgiving in detail (the re-check must abort-and-retry, discard the temp file, and not reference cross-block shell state), so the review verifies the shell logic and the retry re-entry, not just presence.
**Rigor:** Careful main-context review of the exact shell logic (owner re-check placement, temp-file discard, loop re-entry) + confirmation the one-minute bound and fixed-mtime property are untouched.

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Decisions

**Retry re-enters by re-running the block (`exit 3`), not by a shell `while` wrapper.** The obvious reading of "re-enter the acquisition loop" is to wrap the prescription in `while :; do … done` and `continue` on mismatch. That would be wrong here: the critical section is a *comment*, not shell — `$updated_lock_json` is assembled by the agent from the fresh in-mutex read across a model round-trip. A literal shell loop would therefore re-acquire the mutex and then `mv` the **same pre-eviction image**, publishing exactly the lost update this REQ exists to prevent, while looking correct. So the mismatch branch discards the temp file and exits with a distinct code 3 plus a stderr instruction, and the prose below the block states the retry contract: re-run from the top, re-derive `$updated_lock_json` inside the *new* critical section, never re-`mv` the discarded image. This also matches the block's existing convention for an invalidated premise ("write nothing, release the mutex, and re-enter the procedure at the branch the fresh state selects") — prose re-entry, not a shell loop.

**Exiting on mismatch deliberately skips the release.** An evicted session no longer owns the mutex, so there is nothing to release; the existing release-time owner check would no-op anyway, but exiting makes the intent explicit and cannot cascade into the successor's mutex.

**Staging failure is now reported instead of swallowed.** Splitting `printf … && mv … || rm -f` forced a decision about the failed-`printf` path: the old chain ran `rm -f` and exited 0, so a failed lock write reported success. Leaving it silent would also have let the new ownership check run on a path with no staged file. The staging step now reports to stderr and `exit 1`, mirroring the block's existing `mkdir`-failure convention.

**Amended the paragraph below the block, not just the comment.** That prose repeated the softened claim verbatim ("a mutex younger than a minute can only be removed out from under a live owner mid-write, which re-opens the lost-claim → Crash Recovery re-queue failure"). Same file, same section, same property this REQ changes — leaving it would have made the section self-contradicting. The replacement keeps the original conclusion (still no shorter attempt-count break) but on honest grounds: the re-check narrows the window to the instant between itself and the rename rather than closing it, so early eviction costs a discarded read-modify-write cycle and another draw at that window.

## Implementation Summary

**Files changed:**

- `actions/work-reference.md` (modified) — `## Concurrent-Orchestrator Lock Guard (Step 1)`: split the serialized-lock mutex block's `printf … && mv` one-liner to insert a pre-`mv` owner re-check (mismatch → discard temp file, `exit 3`, re-acquire), and softened the mtime-reclaim comment to claim age (not death).

**What was done:**

1. **Pre-`mv` ownership re-check added.** The prescribed `Serialized lock updates` block's one-line `printf … && mv … || rm -f` was split into: stage the temp file (with a reported `exit 1` on staging failure), then `if [ "$(cat "$mutex_path/owner" 2>/dev/null)" != "$session_id" ]` → `rm -f "$lock_tmp"`, stderr notice, `exit 3`; otherwise the original `mv -f "$lock_tmp" "$lock_path" || rm -f "$lock_tmp"` runs unchanged. Every variable the check reads (`$session_id`, `$mutex_path`, `$lock_path`, `$lock_tmp`) is already defined inside the same block, so no cross-block shell state was introduced. The block remains one copy-pasteable prescription and is `sh -n`-clean.
2. **The mtime comment now claims age, not death.** "a mutex older than a minute has a verifiably dead owner" became "proves the mutex's AGE — not that its owner is dead", followed by why reclaiming on age alone is nonetheless safe (the new pre-write re-check discards an evicted owner's stale image instead of publishing it). The "no attempt-count override on a younger mutex" conclusion is retained on the corrected rationale.
3. **Untouched, per the REQ's constraints:** `-mmin -1` (the one-minute reclaim bound), the fixed-mtime property (nothing was added that touches the mutex directory after creation — the owner file is still written exactly once), the `mkdir` acquisition loop and its environment-failure path, the release-time owner check, the temp-file-plus-rename atomicity guarantee, and the 45-minute heartbeat threshold. Serial and single-session runs are behaviorally identical: with no eviction, the re-check compares the token this session wrote moments earlier to its own id and passes.

**Red-Green:** the captured RED walk-through now ends with session A discarding its write. Simulated against a scratch queue: A staged its image, B held the owner token → A exited 3, published nothing, left no `.tmp`, and left B's mutex intact; with the token matching, the same code published and released as before.

## Review

**Acceptance: Pass — overall ~96%.** Main-context review of the exact shell logic + the RED/GREEN simulation the builder ran.

**Both requirements met:**
1. Pre-`mv` owner re-check added: stage temp → `if [ "$(cat "$mutex_path/owner" 2>/dev/null)" != "$session_id" ]` → discard temp, `exit 3`, re-acquire; else `mv` unchanged. Every variable read is defined within the same block — no cross-block shell state. `bash -n`/`sh -n` clean; simulated: eviction → exit 3, lock byte-unchanged, mutex left to successor, no `.tmp` leak; happy path → publish + release as before.
2. Comment softened to "age, not death," pointing at the re-check as what makes age-based eviction safe; the paragraph below (which repeated the old claim) reconciled too.

**Design correctness:** D-01's key insight — retry via `exit 3` + prose re-entry, NOT a shell `while` loop, because a loop would re-`mv` the *same pre-eviction image* and republish the lost update. Exit also correctly skips the release (the mutex is the successor's now). Staging failure now reported (`exit 1`) instead of silently swallowed — a sound in-scope improvement forced by splitting the one-liner.
**Constraints honored:** `-mmin -1` and the fixed-mtime property untouched; single copy-pasteable block; serial/single-session behaviorally identical. contract-regressions PASS; the mutex prescription has no copy-paste siblings (swept).

No Important/Critical findings. No follow-ups.

## Lessons Learned
**Worth knowing:** The mutex "critical section" is a *comment*, not shell — the agent assembles `$updated_lock_json` across a model round-trip. So "re-enter the loop" must be prose re-entry (re-derive from a fresh in-mutex read), never a literal `while`+`continue`, which would republish the pre-eviction image. The age-based reclaim + pre-write owner re-check narrows (does not close) the lost-update window; the residual is the instant between the re-check and the rename.

## Orientation
The serialized-lock mutex now re-verifies ownership immediately before publishing the lock file, so an owner evicted mid-critical-section (a slow/stalled live session past the one-minute reclaim age) discards its stale image (`exit 3`, re-acquire) instead of clobbering its successor's state; age-based reclaim is now honest about proving age, not death. Lives in `actions/work-reference.md` → Concurrent-Orchestrator Lock Guard (Serialized lock updates). No map change — hardens the lock mutex.
