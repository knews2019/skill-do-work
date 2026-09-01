---
source_type: req_lesson
req_id: REQ-314
req_path: do-work/archive/UR-064/REQ-314-judge-effort-estimate-on-review-minted-follow-ups.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, judge, effort, estimate, review]
---

# Lessons from REQ-314: Judge effort_estimate on review-minted follow-ups too

## What the REQ was about

REQ-308 made capture judge `effort_estimate` on every REQ it mints, by the same three-way contract
`impact:` already carried. The other writer of new REQs kept the weaker rule.
`actions/work-reference.md` → **Discovered Tasks Classification (Step 8)** and
`actions/review-work.md` Step 10 both tell that path to "write `effort-mechanical` only when you
have actually judged the fix small, and otherwise leave it absent to read as `effort-substantive`".

Half of that is right and stays: never invent `effort-mechanical`. The other half is permission not
to judge, and it is the rule capture just lost.

## Solution summary

Both automatic follow-up writers now judge size and emit the matching canonical
effort token, ask the user when size is genuinely unclear, and omit only when neither judging nor
asking was possible. The schema and board mirror name every writer, `work.md` delegates its
review-follow-up field list to the canonical Step 10 judgment, and REQ-308's existing semantic
contract independently verifies both field-anchored instructions and the emitted template.

## What worked

- Comparing every REQ writer with one canonical semantic property keeps capture, review follow-ups,
  and Discovered Tasks aligned without adding another competing contract.
- Mutation-testing each writer independently caught both one-sided and symmetric weakening, while a
  separate template assertion kept emission coverage at the actual output seam.

## What didn't work

- A section-wide keyword scan initially looked semantic but allowed unrelated Step 10 prose to
  satisfy the ask-user and never-copy-default legs. Fenced payload isolation alone was insufficient;
  the check also had to isolate the paragraph containing the actual `effort_estimate` directive.

## Back-reference

See `do-work/archive/UR-064/REQ-314-judge-effort-estimate-on-review-minted-follow-ups.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `328767f`.
