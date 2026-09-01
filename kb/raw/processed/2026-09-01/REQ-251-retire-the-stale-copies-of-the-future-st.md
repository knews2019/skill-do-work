---
source_type: req_lesson
req_id: REQ-251
req_path: do-work/archive/UR-055/REQ-251-retire-the-stale-copies-of-the-future-stamp-message.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, retire, stale, copies, future]
---

# Lessons from REQ-251: Retire the stale copies of the future-stamp message

## What the REQ was about

REQ-245 rewrote the future-stamp diagnosis in five renderers. Two copies of the old wording survive outside its write set, both harmless to behaviour and both misleading to the next person who greps for the message.

## Solution summary

Retired the two surviving copies of the pre-REQ-245 future-stamp wording. The `verify_test.go` fixture literals (input-only; assertions unchanged) now derive their reason from production via the existing `reversedSpanAnomalyReason(t)` helper, so a future message change cannot strand another copy; `timestamp.go`'s `formatCanonicalTimestamp` comment now matches its test twin's claim ("one of the two corruptions…") instead of overstating ("the exact corruption…").

## What worked

**What worked:** Deriving fixture text from production (helper call) instead of re-pasting the current wording — the third copy would have stranded at the next message move exactly as the first two did.

**Worth knowing:** `reversedSpanAnomalyReason(t)` in `timestamp_test.go` is the canonical way for any test in the package to obtain the production reversed-span reason; new fixtures should call it, never paste.

## Back-reference

See `do-work/archive/UR-055/REQ-251-retire-the-stale-copies-of-the-future-stamp-message.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `96bb593`.
