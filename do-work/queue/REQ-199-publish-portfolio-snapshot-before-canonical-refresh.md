---
id: REQ-199
title: Publish portfolio snapshot before canonical refresh
status: pending
domain: general
created_at: 2026-08-15T17:09:44Z
user_request: UR-042
addendum_to: REQ-190
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
---

# Review Fix: Publish Portfolio Snapshot Before Canonical Refresh

## What

Make the Yes/unavailable snapshot branch publish its exclusive, no-clobber snapshot from the retained bytes before atomically refreshing the canonical portfolio file. A snapshot publication failure must leave the prior canonical summary unchanged instead of partially completing the promised two-output branch.

This is a standalone user-visible publication-order contract and cannot fold into the existing image-generation or ID-normalization follow-ups: its fix applies only to the portfolio writer's canonical-plus-snapshot transaction.

## Context

Found during review of REQ-190. The new action resolves the snapshot name before writes but tells the agent to refresh the canonical file before publishing the exclusive snapshot, allowing a late collision/publication failure to leave only the canonical output.

## Requirements

- In the Yes/unavailable branch, exclusively publish the snapshot first from the retained byte sequence.
- Only after snapshot success, atomically refresh the canonical file from those same bytes.
- In the No branch, atomically refresh only the canonical file.
- Preserve collision suffixing, no-clobber snapshots, byte identity, artifact immutability, and no automatic cleanup.
- Add or identify replayable contract assertions for snapshot-publication failure ordering and both successful branches.

## Red-Green Proof

**RED prompt/case:** Inspect or replay the Yes/unavailable branch with a snapshot candidate that becomes occupied or cannot be exclusively published after path resolution; current ordering can refresh the canonical file before snapshot failure.
**Why RED now:** The branch promises canonical plus snapshot output but can partially complete while destroying the previous canonical state.
**GREEN when:** Snapshot publication happens first, failure leaves the prior canonical file unchanged, success atomically refreshes canonical from identical bytes, and the No branch still refreshes canonical only.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
