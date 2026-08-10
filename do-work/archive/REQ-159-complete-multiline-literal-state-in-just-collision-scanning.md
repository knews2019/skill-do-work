---
id: REQ-159
title: "Review fix: Complete multiline literal state in Just collision scanning"
status: completed
completed_at: 2026-08-10T10:15:38Z
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
write_set:
  - tools/replace-text-section.sh
  - skills/do-work/tools/replace-text-section.sh
  - _dev/tests/contract-regressions.sh
kb_status: pending
kb_entry:
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

## Exploration

The root and staged replacers are byte-identical. `just_multiline_string_state()` persists only `'''`/`"""`; ordinary quote state is recreated per physical line and triple backticks are grouped with ordinary backticks, while `just_definition_names()` skips a line only when a persistent delimiter was already active. `_dev/tests/contract-regressions.sh` already provides the production-helper pattern: Just-parseable before/after positives, `.before` byte snapshots, exact sorted diagnostic equality, and `cmp -s` pre-mutation controls.

Just 1.46.0 parse probes confirm the three captured forms and also accept physical newlines inside ordinary single-backtick commands. Because the Requirements enumerate ordinary single/double quotes plus triple backticks and explicitly describe ordinary backticks as one-line, that fourth form is a discovered boundary rather than permission to broaden this builder. The three-file scope is sufficient; installers consume the helper and remain verification-only. Main risks are longest-opener ordering, raw versus cooked closing rules, suffix scanning after close/reopen, and ensuring inactive comments/recipe bodies do not seed state.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `tools/replace-text-section.sh` (modify) — retain ordinary quote and triple-backtick lexical state across physical lines.
- `skills/do-work/tools/replace-text-section.sh` (modify) — byte-identical shipped helper correction.
- `_dev/tests/contract-regressions.sh` (modify) — Just-parseable positives and exact byte-preserving collision controls.

**Files I will NOT touch:** installers, updater, board template, reserved-name derivation, managed-marker handling, ADR-019, UR-031 input, REQ-158/160/161 surfaces, or owner-only lifecycle/version/changelog files.

**Acceptance criteria (restated from REQ):**
- [ ] Ordinary raw single-quoted and cooked double-quoted strings retain cross-line state with their actual closing and escape rules.
- [ ] Triple-backtick command literals retain cross-line state, matched before ordinary backticks, without comments, indented recipe bodies, or same-line ordinary backticks hiding real definitions.
- [ ] Just-parseable positive fixtures cover all three captured forms; real definitions immediately around them still produce exact sorted diagnostics and preserve target bytes.
- [ ] Existing triple-single/triple-double strings, managed-span exclusion, reserved-name derivation, deterministic reporting, mutation ordering, paired identities, and full contracts remain green.

## Discovered Tasks

- [normal] Just 1.46.0 accepts physical newlines inside ordinary single-backtick commands containing column-zero recipe/alias-shaped payload, despite this REQ's explicit one-line-backtick boundary. Capture a consent-gated follow-up so the broader multiline-literal completion claim can be reconciled without silently expanding REQ-159.

## Pre-Flight

**Git:** ✓ No pre-existing changes outside `do-work/`; previously protected ADR-019 and UR-031 input bytes are now committed by external workspace activity and remain hash-identical to the approved state. Live queue/claim changes remain lifecycle state, not builder scope.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh` passes (`launched: true`).
**Dependencies:** ✓ Bash, Python 3, Just 1.46.0, and ShellCheck are available; paired helpers begin byte-identical.

*Checked by work action*

## Root Cause

`just_multiline_string_state()` persisted only triple-single and triple-double delimiters. Ordinary quote state was recreated for every physical line, and triple backticks entered the one-line backtick branch, so subsequent column-zero payload reached `just_definition_name()` as if it were a top-level definition.

## Decisions

- **D-01 — Persist only the two captured ordinary quote forms.** Carry an unclosed raw single quote or cooked double quote to the next physical line, with raw literal closing and cooked contiguous-backslash parity; deliberately keep an unclosed ordinary single backtick line-local because that accepted Just form is already recorded as a separate Discovered Task.
- **D-02 — Treat only an exact triple-backtick run as the multiline command delimiter.** Match it before ordinary backticks, require a run boundary on both sides, retain it across physical lines, and leave the existing top-level/comment/indent gate unchanged so comments and recipe bodies cannot seed state.
- **D-03 — Prove false-positive removal separately from real-collision retention.** Use a Just-parseable positive containing all three captured forms, an exact four-name control bracketing each form while reserving `kanban-static` as literal-only payload, and a separate byte-preserving `kanban-static` control after inactive comment/body/one-line-backtick forms.

## Implementation Summary

**Files changed:**
- `tools/replace-text-section.sh` (modified) — persists ordinary raw single/cooked double quote and exact triple-backtick lexical state across physical lines.
- `skills/do-work/tools/replace-text-section.sh` (modified) — byte-identical shipped helper correction.
- `_dev/tests/contract-regressions.sh` (modified) — adds Just-parseable acceptance fixtures, exact sorted nearby-collision controls, inactive-form shielding controls, and pre-mutation byte checks.

**What was done:** The collision scanner now ignores reserved-looking payload inside the three captured multiline literal families while retaining exact real-definition detection, deterministic diagnostics, and atomic pre-mutation rejection around them.

## Qualification

Passed — all three scoped implementation files verified, all four Requirements and Instances traced, P-A-U confirmed, Scope/Implementation Summary sets match, and the paired helpers are byte-identical. The diff leaves definition grammar, managed-span removal, reserved derivation, sorting/diagnostics, and atomic ordering unchanged; contamination against REQ-158 is clean because no prior Markdown-guard file appears.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/install-suite-behavior.sh`; `bash _dev/tests/staged-skills-contract.sh`; `bash _dev/tests/suite-manifest-contract.sh`; paired helper `cmp -s`; `bash -n` and `shellcheck -S warning` on both helpers and the modified contract; `just --justfile skills/do-work-board/justfile.template --list`; changelog `cmp -s`; `git diff --check`
**Result:** ✓ All passing.

**Red-green validation:**
- Valid ordinary single/double multiline values and triple-backtick command literal: ✗ old scanner reported all five reserved-looking payload names → ✓ insertion succeeds and the resulting Justfile parses.
- Nearby real-definition control: ✗ old scanner added literal-only `kanban-static` to the diagnostic → ✓ exact four real names remain sorted and rejection preserves bytes.
- Inactive comment/body/one-line-backtick control: ✓ real `kanban-static` remained visible before and after the lexer change, with exact diagnostic and byte preservation.

**New tests added:** Just-parseable production-helper fixtures for raw single, cooked double, and exact triple-backtick multiline literals; exact nearby real-collision and inactive-form shielding controls.

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` (from REQ-156): extends the same collision-scanner contract while retaining triple-single/triple-double, managed-span, deterministic-reporting, and transaction-order assertions.

*Verified by work action*

## Review

**Overall: 84%** | 2026-08-10T10:14:57Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 94% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- Physical-newline ordinary single-backtick commands remain Just-parseable but their reserved-looking payload still false-collides, so the broad “every current Just multiline literal form” completion claim is partial — gate: user-visible → covered by the existing Discovered Task and rerouted pending-answers as REQ-162; no duplicate review follow-up created

**Minor findings:** 0 (report only)
**Acceptance:** Partial — all three explicit literal families and preservation controls pass, but the already-recorded ordinary single-backtick multiline boundary remains.
**Suggested testing:** 1 item
**Follow-ups created:** None by review; ordinary-single-backtick follow-up created from the existing Discovered Task as REQ-162; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Just-parseable positives paired with exact sorted diagnostics and byte snapshots proved both safe acceptance and unchanged rejection behavior through the production helper.
- Persisting delimiter identity while leaving the existing definition grammar untouched kept the paired change small and reviewable.

**What didn't:**
- The captured three-family boundary called ordinary backticks one-line, but Just 1.46.0 accepts them across physical lines too; implementing the broad “every current form” claim without probing adjacent grammar would have silently overstated completion.

**Worth knowing:**
- A safety scanner that approximates a lexer needs accepted-syntax probes for every delimiter family. Explicit Requirements should still bound the builder when a broad summary conflicts; record the adjacent accepted form for consent instead of expanding silently.

## Orientation

The managed-Just collision guard now accepts reserved-looking payload inside physical-newline ordinary single/double strings and exact triple-backtick commands while preserving real collision detection and atomic rejection. Just's additional ordinary single-backtick multiline form remains consent-gated in REQ-162.
