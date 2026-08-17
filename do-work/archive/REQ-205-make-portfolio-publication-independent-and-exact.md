---
id: REQ-205
title: Make portfolio publication independent and exact
status: completed
completed_at: 2026-08-17T18:45:53Z
commit: 4a0b90e
claimed_at: 2026-08-17T18:41:41Z
status_changed_at: 2026-08-17T18:10:49Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-17T18:42:20Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - persistence changes
    - cross-route regression gates
write_set:
  - skills/do-work-toolbox/scripts/publish-portfolio-summary.sh
  - skills/do-work-toolbox/actions/present-work.md
  - _dev/tests/prescribed-shell-scripts-behavior.sh
domain: general
created_at: 2026-08-15T20:01:55Z
user_request: UR-042
addendum_to: REQ-199
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: portfolio-independent-exact-publication
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
---

# Review Fix: Make Portfolio Publication Independent and Exact

## What

Make snapshot and canonical outputs independent regular files at their exact requested paths while preserving snapshot-first ordering, exclusive no-clobber publication, atomic canonical refresh, and all existing failure guarantees.

## Context

REQ-199 closes canonical-first publication on ordinary paths, but review proved that its hard-linked outputs share later mutations and that `ln`/`mv` can treat directory operands as containers while returning success.

## Instances

- [ ] Output independence: later in-place changes to canonical must not change a published snapshot; tests must reject same-inode aliasing.
- [ ] Exact path semantics: snapshot-candidate directories must advance or fail closed, canonical directories must fail, and no hidden private file may be nested in either.

## Requirements

- Publish snapshot and canonical from the same verified bytes without sharing a mutable inode after success.
- Preserve snapshot-before-canonical ordering and prior-canonical preservation on snapshot failure.
- Treat an existing snapshot directory as occupied and safely suffix or fail closed.
- Reject an existing canonical directory before publication and leave it unchanged.
- Verify and report the exact regular-file paths actually published.
- Add behavior tests for canonical post-write mutation and both directory operand cases.

## Red-Green Proof

**RED prompt/case:** Publish canonical plus snapshot, overwrite canonical in place, and observe the snapshot changes; then pre-create snapshot/canonical targets as directories and observe private files nested while the helper reports success.
**Why RED now:** Same-inode equality is not durable immutability, and utility exit zero does not prove exact-path publication when the destination is a directory.
**GREEN when:** Snapshot bytes survive later canonical mutation, outputs are different inodes, snapshot-directory collision advances/fails safely, canonical-directory publication fails unchanged, no nested private file exists, and all prior branch/failure tests pass.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] Snapshot-first ordering is fixed, but durable immutability and directory-destination safety need another focused change. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Triage

**Route:** B — Explore then Build

**Reasoning:** Both instances are precisely stated, but the publisher's existing exclusivity primitive (`ln` as an atomic no-clobber create) and the exact test that *asserted* the defect had to be found before either could be changed safely.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route B.

## Exploration

**Key files:**
- `skills/do-work-toolbox/scripts/publish-portfolio-summary.sh` — one private verified copy, hard-linked to the snapshot, then renamed to canonical.
- `_dev/tests/prescribed-shell-scripts-behavior.sh` lines 227–295 — five existing named cases, one of which **asserted the defect**: `[ "$portfolio_candidate" -ef "$portfolio_canonical" ]` required both outputs to share an inode.
- `skills/do-work-toolbox/actions/present-work.md` lines 134–136 — the prose contract the script implements.

**Concerns found:**
- One inode, two names: `ln private snapshot` then `mv private canonical` leaves snapshot and canonical as links to the same file, so any later in-place write to canonical rewrites the "immutable" snapshot.
- `ln file dir` links *into* the directory and exits 0, so a snapshot candidate occupied by a directory published a hidden `.portfolio-summary.md.publishing.XXXXXX` inside it and reported success.
- `mv file dir` moves *into* the directory and exits 0, so a canonical path occupied by a directory did the same.

## Decisions

- **D-02: verify each publication after the fact rather than pre-checking the destination type.** A pre-check alone is a check-then-act window; the post-check catches both the pre-existing directory and one that appears during publication, using one mechanism instead of two. Reasoning: same shape REQ-204 landed for the ai-report batch an hour earlier, so the repo now states this rule once and applies it identically in both publishers. This is DECIDE & STATE — reversible, no user-visible surface.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/scripts/publish-portfolio-summary.sh` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**Acceptance criteria (restated from the REQ):**
1. Snapshot and canonical publish from the same verified bytes without sharing a mutable inode after success.
2. Snapshot-before-canonical ordering and prior-canonical preservation on snapshot failure are preserved.
3. An existing snapshot **directory** is treated as occupied — advance or fail closed.
4. An existing canonical **directory** is rejected before publication and left unchanged.
5. The exact regular-file paths actually published are verified and reported.
6. Behavior tests exist for canonical post-write mutation and both directory operand cases.

## Pre-Flight

- Working tree clean outside `do-work/`.
- Baseline `bash _dev/tests/prescribed-shell-scripts-behavior.sh` passing before the change (31 named cases after this REQ, 29 before).

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/scripts/publish-portfolio-summary.sh` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** Gave each output its own verified private copy and made both publication steps verify the path they actually wrote.

*Independence.* A new `allocate_private_copy` helper mints one private file per output and `cmp`-verifies it against the source, so snapshot and canonical carry identical bytes from two distinct inodes. After the snapshot's exclusive `ln` succeeds, its private name is unlinked, leaving the snapshot as the sole link to its own inode. A later in-place edit of canonical can no longer reach it.

*Exactness.* Both publications are verified after the operation, not only before it. The snapshot loop checks for its own private name nested under the candidate: present means `ln` linked into a directory, so it removes only its own nested link and advances to the next numeric suffix. Canonical is rejected up front when it is a directory, and the same nesting check runs after the rename to catch a directory that appeared in the window; either way the occupying directory is untouched and the run exits nonzero. A final `[ -f ]` gate confirms each reported path is a regular file before anything is printed.

*Contract.* `present-work.md` now states the independence guarantee and the directory-operand behavior, so the prose and the script cannot drift.

**Tests touched:** one existing assertion inverted (see below) plus two new named cases — snapshot candidate occupied by a directory, and canonical occupied by a directory. Named-case count 29 → 31.

## Qualification

Passed — 3 files verified, 6 requirements traced, no debug artifacts.

- All three declared files appear in the diff; nothing undeclared was touched.
- Substantive: the script's publication path was restructured, not reformatted; the tests gained two cases and a corrected assertion.
- Requirements traced: two private copies + unlink (1); snapshot loop still precedes the canonical branch and exits before it on failure (2); suffix advance on a directory candidate (3); canonical directory rejected and re-checked after the rename (4); `[ -f ]` gates before reporting (5); three new/changed assertions (6).
- Flowing: `allocate_private_copy` is called on both paths; `report_retained_snapshot` is called on every canonical failure branch.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh` (baseline, RED, GREEN); `bash _dev/tests/maintainer-verify.sh`

**Result:** ✓ prescribed-shell suite exit 0, 31 named cases; ✓ maintainer-verify exit 0, zero FAIL lines.

**Red-green validation:** ✗ RED — with the new assertions in place and the script reverted, the suite exits 1 with eight failures spanning all three defect families: `snapshot-success case aliased the snapshot to the canonical inode`, `snapshot-success case let a later canonical edit rewrite the snapshot`, `snapshot-directory case nested a private file inside the occupying directory`, `snapshot-directory case did not advance to the numeric suffix`, `snapshot-directory case reported the wrong published paths`, `canonical-directory case reported success`, `canonical-directory case did not leave the occupying directory unchanged`, `canonical-directory case leaked private bytes`. → ✓ GREEN — all eight pass, and the five pre-existing portfolio cases (canonical-only, snapshot success, occupied-file collision, `ln` failure, `mv` failure) pass unchanged.

**Existing tests updated (cross-REQ impact):** `_dev/tests/prescribed-shell-scripts-behavior.sh` — REQ-199's assertion `[ "$portfolio_candidate" -ef "$portfolio_canonical" ]` ("did not publish both paths from the same private bytes") is **inverted** to reject same-inode aliasing. REQ-199 was testing that both outputs came from one verified copy, which was the right intent; sharing an inode was the wrong proxy for it, and this REQ's `## What` names rejecting that aliasing as a requirement. Byte identity — REQ-199's actual guarantee — is still asserted by the two `cmp -s` checks on the line above, so nothing REQ-199 locked in was dropped. The change is traceable through the inline comment naming both REQs.

*Verified by work action*

## Review

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Six criteria, each with its own failing-then-passing assertion |
| Code Quality | 92% | Two private copies is one more moving part, earned by the guarantee it buys |
| Test Adequacy | 98% | The inverted assertion is the interesting one — it was locking in the defect |
| Scope | 100% | Three declared files; the sibling `install-last30days.sh` hit stayed with REQ-220 |
| Risk | Low | Publication is now stricter; the only new failure mode is refusing a directory that used to be silently nested into |
| Acceptance | Pass | Eight assertions fail on the old script and pass on the new one |

**Verdict: Approve** — the snapshot is now durably immutable, and both publications verify where they actually landed.

### Findings

**Minor:**
- Two `cp` operations of the same source instead of one is a small cost on large summaries. Copying twice is what buys independent inodes without a second read-verify strategy, and portfolio summaries are text; noted rather than optimized.

**Nit:**
- `report_retained_snapshot` prints to stderr on the canonical-directory path even in `--canonical-only` mode, where `snapshot_path` is empty and it prints nothing. Harmless, and guarding the call site would be more code than the guard inside it.

### Restatement Sweep

**Triggered** — the diff redefines what "published from the same bytes" guarantees. Swept `inode`, `same bytes`, and `byte-identical` across `skills/`, `_dev/`, and the guides. Results: `present-work.md` was the one consumer restating the contract and was updated in the same commit; `docs/present-work-guide.md`'s "byte-identical snapshot" is still exactly true; `_dev/primes/prime-shell-commands.md`'s REQ-199 lesson line ("not only command success or inode identity") already reads correctly against the new behavior. The only stale restatement was the test assertion itself, which this REQ inverts.

### Requirements Checklist

- [x] Same verified bytes without a shared mutable inode — delivered
- [x] Snapshot-before-canonical ordering and prior-canonical preservation on snapshot failure — delivered (pre-existing cases pass unchanged)
- [x] Snapshot directory treated as occupied; advance or fail closed — delivered (advances)
- [x] Canonical directory rejected, left unchanged — delivered
- [x] Exact regular-file paths verified and reported — delivered
- [x] Behavior tests for post-write mutation and both directory operands — delivered

### Acceptance Testing

**Result: Pass**
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh` — exit 0, 31 named cases.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines (covers the ShellCheck lint over the rewritten script).
- Finding-Closure Ratchet: the captured RED prompt is replayed exactly — publish both, overwrite canonical in place, observe the snapshot; then pre-create each target as a directory. All three observations are now assertions that fail against the pre-change script.

### Suggested Additional Testing

- A snapshot directory and a canonical directory in the *same* invocation is not replayed; each is covered alone. The paths are independent, but one combined case would prove the retained-snapshot message on the canonical-directory branch.
- Cross-filesystem behavior: `ln` from the canonical directory to the snapshot directory has always assumed one filesystem. Unchanged by this REQ, but worth an explicit case if snapshots ever move to a separate mount.

### Follow-up REQs Created

None — the one sibling file carrying this same class (`install-last30days.sh`) was already routed to REQ-220 by REQ-204's sweep, and reporting it twice would create a duplicate.

## Lessons Learned

**What worked:** Reading the existing test before the script. The assertion that had to be *inverted* was the fastest possible confirmation that the defect was designed in rather than overlooked — REQ-199 proved "one verified copy" by proving "one inode", and the proxy outlived the property.

**What didn't:** A pre-check on the destination type looked sufficient until it was written out; it leaves the same check-then-act window the REQ was complaining about. Verifying after the operation covers the pre-existing directory and the one that appears mid-flight with a single mechanism.

**Worth knowing:** `ln file dir` and `mv file dir` both succeed and both mean "put it inside", not "collide with". Neither exit status can distinguish a publication from a nesting, so any exclusive-create built on them needs a post-condition check on the exact path. And a hard link is not a copy: it makes two names for one mutable object, which is precisely wrong for anything called a snapshot.

## Orientation

The portfolio's timestamped snapshot is now durably immutable — editing the canonical summary afterwards no longer rewrites it — and both publication steps refuse to publish into a directory that occupies their destination. Lives in the portfolio publication path (`skills/do-work-toolbox/scripts/publish-portfolio-summary.sh`) with its contract stated in `present-work.md`. No new module or command; the guarantee behind an existing output got stronger, so the system's shape is unchanged.

Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` and `_dev/primes/prime-action-files.md` — all referenced paths resolve; the REQ-199 lesson line still reads true and needs no correction.
