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

Make the shared completed-work presentation resolver inherit `actions/work-reference.md` → **Target ID Resolution** before archive lookup, so case-insensitive prefixes and numeric-value matching resolve canonical stored IDs such as `REQ-042` and `UR-011`.

This is a standalone user-visible input contract and cannot fold into a sweep: its fix is unrelated to output-directory publication and has one canonical resolver surface.

## Context

Found during review of REQ-189. The new shared reference says to find an exact archived folder or REQ match but does not first normalize equivalent input spellings such as `req-42`, `REQ-42`, and `REQ-042`.

## Requirements

- Cite and apply the shared Target ID Resolution contract before UR or REQ archive lookup.
- Preserve the presentation action's archive-only search locations and terminal-success gates.
- Add or identify a replayable contract assertion covering case-insensitive and zero-padding equivalents.

## Red-Green Proof

**RED prompt/case:** Inspect the shared resolver for a presentation request using `req-42` or `REQ-42` when the archived file stores `REQ-042`; no canonicalization step exists before exact lookup.
**Why RED now:** Raw case-sensitive or zero-padding-sensitive lookup can reject a valid equivalent ID token.
**GREEN when:** The shared resolver cites Target ID Resolution, canonicalizes the token before lookup, and a replayable assertion accepts equivalent spellings while preserving archive-only and success-status gates.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
