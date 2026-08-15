---
id: REQ-205
title: Make portfolio publication independent and exact
status: pending-answers
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

- [ ] Snapshot-first ordering is fixed, but durable immutability and directory-destination safety need another focused change. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.
