---
id: REQ-197
title: Normalize completed-work presentation target IDs
status: pending
domain: general
created_at: 2026-08-15T16:34:08Z
user_request: UR-042
addendum_to: REQ-189
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
maintenance: true
---

# Review Fix: Normalize Completed-Work Presentation Target IDs

## What

Make every completed-work presentation ID path touched by UR-042 inherit `actions/work-reference.md` → **Target ID Resolution** before dispatch or archive lookup, so case-insensitive prefixes and numeric-value matching resolve canonical stored IDs such as `REQ-042` and `UR-011`.

This is a standalone user-visible input contract and cannot fold into a sweep: its fix is unrelated to output-directory publication and has one canonical resolver surface.

## Context

Found during review of REQ-189. The new shared reference says to find an exact archived folder or REQ match but does not first normalize equivalent input spellings such as `req-42`, `REQ-42`, and `REQ-042`.

Review of REQ-190 found the same root cause in `present-work` item-specific migration guidance: its dispatcher recognizes canonical-looking `UR-NNN`/`REQ-NNN` text but does not explicitly accept case-insensitive, numeric-equivalent forms before printing both replacement commands.

## Requirements

- Cite and apply the shared Target ID Resolution contract before UR or REQ archive lookup in the completed-work presentation reference.
- Apply the same token grammar to `present-work` item-specific migration dispatch while preserving the supplied token in the two printed replacement commands.
- Preserve the presentation action's archive-only search locations and terminal-success gates.
- Add or identify replayable contract assertions covering case-insensitive and zero-padding equivalents at both presentation entry paths.

## Red-Green Proof

**RED prompt/case:** Inspect the shared resolver and `present-work` migration dispatcher for a presentation request using `req-42` or `REQ-42` when canonical storage uses `REQ-042`; neither path currently applies the shared input-token grammar before lookup or routing.
**Why RED now:** Raw case-sensitive or zero-padding-sensitive matching can reject a valid equivalent ID token or print generic usage instead of migration guidance.
**GREEN when:** Both presentation entry paths cite Target ID Resolution, accept equivalent spellings, preserve their own lookup/write gates, and replayable assertions cover the shared grammar.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
