---
id: REQ-753
title: 'Show an actual preview for every gallery report'
status: claimed
created_at: 2026-11-19T12:02:21Z
user_request: UR-186
domain: general
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: false
maintenance: false
addendum_to: REQ-749
related: [REQ-747, REQ-748]
claimed_at: 2026-11-19T12:02:31Z
---
# Show an Actual Preview for Every Gallery Report

## What

Show a screenshot of the original page on each of the 55 gallery cards, including generated and older reports. Preserve existing report links and controls. Ship the prepared implementation as the 0.331.1 patch.

## Prior Implementation

REQ-749 (real reports and screenshot previews) delivered ten curated previews. The generated and older cards could still display without a preview. The prepared correction is commit `7ca9382ac3311273c8e9d08c6cab86833cdca135`.

## Constraints

Keep report source bytes and reader preferences intact. Verify the installed launcher and screenshot failure handling. Leave the first-round feedback hold and the coordinator's live viewer alone.

## Red-Green Proof

**RED prompt/case:** Open the gallery and expand Older reports. Inspect the five generated cards and 40 older cards.
**Why RED now:** Only the ten curated reports have screenshot covers.
**GREEN when:** Every card displays a decoded screenshot with its original report link, archive, starter and feedback controls still work, and capture failures report an error without publishing a bad screenshot.
**Validation:** Inferred during capture from the user's requirement and the prepared implementation.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Scope the integration of the prepared correction.
- [ ] **[APPLY]:** Retain the prepared implementation in main.
- [ ] **[UNIFY]:** Verify previews, controls, error handling and report hashes.

## Full Context

See `do-work/user-requests/UR-186/input.md` for the verbatim request.
