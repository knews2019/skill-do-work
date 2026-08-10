---
id: REQ-159
title: "Review fix: Complete multiline literal state in Just collision scanning"
status: claimed
claimed_at: 2026-08-10T09:51:10Z
route: C
status_changed_at: 2026-08-10T09:20:51Z
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

- [x] REQ-156 fixed triple-quoted multiline strings, but other valid Just multiline literals can still make installation reject a safe custom Justfile. This is another review-generated follow-up, so the cascade-depth rule requires your consent before it enters the work loop. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: C** - Complex

**Reasoning:** The change is bounded to a paired helper and its contract, but it extends a handwritten Just lexer across three multiline literal families with different escape/closing rules while preserving exact pre-mutation collision behavior. A full plan and exploration are warranted.

**Planning:** Required

## Plan

1. Add Just-parseable RED fixtures for ordinary raw single-quoted and cooked double-quoted multiline values containing column-zero reserved-looking recipe/alias payload, including their actual raw and escape-parity closing rules.
2. Add a Just-parseable RED fixture for triple-backtick command literals plus controls proving comments, indented recipe bodies, and same-line ordinary backticks cannot start cross-line state.
3. Add real reserved definitions immediately around every new literal family and assert exact sorted diagnostics plus byte-preserving pre-mutation rejection.
4. Generalize the paired helpers' lexical state across physical lines, matching longest delimiters first, retaining raw/cooked closing semantics, and continuing suffix scans after closes without changing definition grammar or transaction order.
5. Run focused/full, installer, distribution, helper-identity, Just-parse, syntax/ShellCheck, changelog-identity, and diff checks.

**Root-cause hypothesis:** `just_multiline_string_state()` persists only triple-single/triple-double delimiters. Its ordinary quote state is recreated on every physical line and triple backticks fall through as ordinary tokens, so later column-zero literal payload reaches `just_definition_names()` as a real definition.

**Just semantics to preserve:** Ordinary single quotes are raw and close on the next literal quote; cooked double quotes close only outside an active backslash escape; triple backticks close on the next exact triple run with no escape processing. Active literals treat leading indentation and `#` as payload, while inactive comments/recipe bodies and same-line ordinary backticks must not create cross-line state.

**Exact implementation scope:** Modify `tools/replace-text-section.sh`, `skills/do-work/tools/replace-text-section.sh`, and `_dev/tests/contract-regressions.sh` only; the helper copies must finish byte-identical. Keep installers, updater, templates, reserved-name policy, managed markers, ADR-019, UR-031, and owner-only lifecycle/release files out of builder scope.

**Plan validation:** Every Requirement and Instance maps to a literal fixture, the shared scanner change, and preservation verification; no orphan task was found. Five ordered tasks reach the quality-warning threshold, but they are inseparable RED/control/GREEN/verification phases for one lexer boundary. Just 1.46 also appears to accept physical newlines in ordinary single-backtick commands despite the REQ's one-line wording; do not silently broaden this implementation—verify during exploration and record it as a discovered task if confirmed.

*Generated by Plan agent*
