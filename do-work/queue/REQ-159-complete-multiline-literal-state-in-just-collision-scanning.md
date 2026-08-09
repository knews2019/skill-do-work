---
id: REQ-159
title: "Review fix: Complete multiline literal state in Just collision scanning"
status: pending-answers
domain: general
created_at: 2026-08-09T19:21:41Z
user_request: UR-031
addendum_to: REQ-156
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: just-multiline-literal-state
---

# Review Fix: Complete Multiline Literal State in Just Collision Scanning

## What

Make the reserved-recipe collision scanner retain lexical state for every current Just multiline literal form that can contain column-zero recipe- or alias-shaped payload. Done means valid ordinary multiline single/double strings and triple-backtick command literals cannot recur as false collisions, while real definitions around every form remain detected exactly and pre-mutation preservation remains unchanged.

## Context

Found during review of REQ-156. That task correctly handles triple-single and triple-double strings, but Just 1.46 also accepts physical line breaks inside ordinary quoted strings and triple-backtick command literals; the scanner resets ordinary quote/backtick state at every line and rejects their payload.

## Requirements

- Retain cross-line state for ordinary single-quoted and cooked double-quoted Just strings, using their actual closing and escape rules.
- Retain cross-line state for triple-backtick indented command literals without letting ordinary one-line backticks, comments, or indented recipe bodies hide real definitions.
- Add Just-parseable positive fixtures for all three reproduced forms and exact, byte-preserving real-collision controls immediately around them.
- Keep the existing triple-single/triple-double behavior, reserved-name derivation, managed-span exclusion, deterministic reporting, installer ordering, paired identities, and full contracts passing.

## Instances

- [ ] `tools/replace-text-section.sh:116`: ordinary double-quoted strings may span lines, but their quote state is discarded at the line boundary.
- [ ] `tools/replace-text-section.sh:116`: ordinary single-quoted strings may span lines, but their quote state is discarded at the line boundary.
- [ ] `tools/replace-text-section.sh:162`: a triple-backtick opener is treated as one-line ordinary backticks, so multiline command payload is classified as definitions.
- [ ] `skills/do-work/tools/replace-text-section.sh:116`: apply the same correction byte-identically to the shipped helper copy.

## Open Questions

- [ ] REQ-156 fixed triple-quoted multiline strings, but other valid Just multiline literals can still make installation reject a safe custom Justfile. This is another review-generated follow-up, so the cascade-depth rule requires your consent before it enters the work loop. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
