---
id: REQ-162
title: "Review fix: Handle ordinary multiline backtick commands"
status: completed
completed_at: 2026-08-10T12:01:02Z
commit: aff7c9c
claimed_at: 2026-08-10T11:33:58Z
status_changed_at: 2026-08-10T10:57:17Z
route: C
domain: general
created_at: 2026-08-10T10:14:57Z
user_request: UR-031
addendum_to: REQ-159
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

# Review Fix: Handle Ordinary Multiline Backtick Commands

## What

Extend reserved-recipe collision scanning to retain physical-line state for ordinary single-backtick Just command literals. Done means the broader multiline-literal scanner accepts reserved-looking command payload without weakening real definition detection, exact diagnostics, or pre-mutation preservation.

## Context

Discovered during REQ-159 and independently confirmed in its review. Just 1.46.0 accepts physical newlines inside ordinary single-backtick commands, but REQ-159 deliberately kept them line-local because its explicit Requirements named three other literal families and described ordinary backticks as one-line.

## Requirements

- Retain ordinary single-backtick command state across physical lines using Just's actual closing behavior.
- Keep same-line backtick commands, comments, indented recipe bodies, triple-backtick commands, and every existing string family from hiding real definitions.
- Add Just-parseable positive fixtures plus exact byte-preserving real-collision controls and keep paired helpers/full contracts green.

## Instances

- [ ] `tools/replace-text-section.sh`: an ordinary backtick command spanning physical lines exposes column-zero reserved-looking payload to `just_definition_names()`.
- [ ] `skills/do-work/tools/replace-text-section.sh`: apply the correction byte-identically to the shipped helper.
- [ ] `_dev/tests/contract-regressions.sh`: add positive and exact negative production-helper fixtures for the accepted multiline form.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-159: Just accepts ordinary single-backtick commands across physical lines, and the collision scanner still rejects safe command payload that resembles a reserved recipe. REQ-159 explicitly treated ordinary backticks as one-line, so the cascade-depth rule requires your consent before extending that boundary. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: C** - Complex

**Reasoning:** The correction is bounded to a paired helper and its direct contract, but it extends a handwritten Just lexer with persistent ordinary-backtick state and must preserve exact closing behavior, real-definition detection, diagnostics, pre-mutation safety, and all previously repaired literal families.

**Planning:** Required

## Plan

1. Add a Just-parseable RED positive fixture for an ordinary single-backtick command spanning physical lines. Include column-zero `run-kanban:` and `alias kanban-summary := ignored` payload, and close/reopen on one physical line with `\`` followed by ` + \``. Just 1.46.0 confirms the next single backtick closes the command even when immediately preceded by `\`; ordinary backticks have no cooked-string escape parity.
2. Snapshot the positive fixture before invocation. Confirm the old helper rejects it with exactly `kanban-summary, run-kanban` and leaves bytes unchanged; after the fix, require success, exact expected insertion bytes, and successful `just --list` both before and after insertion.
3. Add a Just-parseable exact negative control with real `run-kanban`, `alias kanban-summary`, `@run-kanban-cli`, and `run-do-work-update` definitions around the multiline command, while `kanban-static` appears only inside command payload. Also place an unmatched ordinary backtick in a comment and indented recipe body plus a closed same-line backtick command before later real definitions. Require the exact sorted diagnostic `kanban-summary, run-do-work-update, run-kanban, run-kanban-cli` and `cmp -s` against the pre-invocation snapshot. This is RED because the current scanner incorrectly adds `kanban-static`.
4. Persist ordinary backtick state in `just_multiline_string_state()` across physical lines with the minimal existing-state extension: match exact triple backticks first as today; otherwise carry an unclosed ordinary backtick as `b"\`"`, close it at the next single backtick without escape processing, continue scanning the same physical-line suffix so close/reopen works, and keep the existing top-level/comment/indent gate unchanged.
5. Apply the production change byte-identically to both helpers. Do not change `just_definition_name()`, reserved-name derivation, managed-span removal, sorted diagnostic text, or transaction order: validate and scan collisions before constructing replacement bytes or calling `atomic_replace()`.
6. Preserve the existing contracts for closed same-line backticks, comments, indented recipe bodies, exact triple-backtick commands, ordinary raw single/cooked double strings, triple-single/triple-double strings, CRLF, same-line close/reopen, managed-span exclusion, and real definitions surrounding every literal family.

**Exact implementation scope:** `tools/replace-text-section.sh`, `skills/do-work/tools/replace-text-section.sh`, and `_dev/tests/contract-regressions.sh` only. Installers, updater, board template, marker/reserved-name policy, release files, ADRs, and UR-031 remain unchanged; the helper copies must finish byte-identical.

**Verification:** Capture RED with the production helper and exact diagnostics above, then run `bash _dev/tests/contract-regressions.sh`, installer/staged-skills/suite-manifest contracts, `cmp -s` on the paired helpers, `bash -n` and warning-level ShellCheck on the three scoped shell files, board-template `just --list`, changelog-mirror identity, and `git diff --check`.

**Discovered tasks:** Grep the paired scanner copies for the same state pattern. Any separate accepted Just syntax gap is reported in the REQ's top-level `## Discovered Tasks` lifecycle record for orchestrator classification; do not fix it inline or expand the three-file implementation scope. No current orphan task was found.

*Generated by Plan agent*

## Exploration

The defect is isolated to the paired helpers' `just_multiline_string_state()` return boundary. The scanner already recognizes byte `96` as an ordinary quote while scanning a physical line, but at end of line it persists only ordinary single/double quotes. An unclosed ordinary backtick is therefore discarded; on the next physical line, `just_definition_names()` sees no active delimiter and falsely classifies recipe- or alias-shaped command payload as a real definition.

Just 1.46.0 probes confirm that ordinary single-backtick commands may span physical lines; the next single backtick closes the command even when immediately preceded by `\`; suffix scanning must continue so a later backtick on the same line can reopen state; and comments or indented recipe bodies remain inactive opener contexts. The minimal correction is to persist byte `96` alongside `34` and `39`. Existing active-delimiter logic already supplies raw closing, suffix scanning, and longest-first exact triple-backtick matching.

Fresh production-helper reproduction shows a valid positive currently exits 1 with exactly `kanban-summary, run-kanban`, while a valid negative control adds literal-only `kanban-static` to the four real sorted collisions; both rejected targets remain byte-identical. Root and shipped helpers begin byte-identical. Similar-pattern search found the state function and persistence condition only in those two paired copies. The full contract-regressions baseline passes.

Main risks are accidentally adding cooked backslash parity, matching ordinary backticks before triple backticks, dropping suffix close/reopen scanning, letting comments or recipe bodies seed state, weakening exact diagnostics, or drifting the paired copies. `just_definition_name()`, reserved derivation, managed-span removal, sorting/error text, and mutation ordering need no change.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `tools/replace-text-section.sh` (modify) — persist an unclosed ordinary single-backtick command across physical lines using the existing raw active-delimiter path.
- `skills/do-work/tools/replace-text-section.sh` (modify) — apply the same correction byte-identically.
- `_dev/tests/contract-regressions.sh` (modify) — add production-helper positive and exact negative fixtures beside the existing multiline-literal controls.

**Files I will NOT touch:** `just_definition_name()`, installers, updater, board template, reserved-name policy, managed-marker handling, transaction ordering, ADRs, UR-031, other REQs, or owner-only release/lifecycle files.

**Acceptance criteria:**
- [ ] An ordinary single-backtick command retains state across physical lines and closes at the next literal backtick, including one immediately preceded by `\`.
- [ ] A close followed by ` + \`` on the same physical line closes and reopens state correctly.
- [ ] A Just-parseable positive containing column-zero `run-kanban:` and `alias kanban-summary := ignored` command payload succeeds through the production helper.
- [ ] The positive parses before and after insertion and matches exact expected insertion bytes.
- [ ] A Just-parseable negative control reports exactly `kanban-summary, run-do-work-update, run-kanban, run-kanban-cli`; literal-only `kanban-static` is absent.
- [ ] Rejection preserves the negative target byte-for-byte.
- [ ] Unmatched ordinary backticks in a comment or indented recipe body, plus a closed same-line command, cannot hide later real definitions.
- [ ] Existing triple-backtick commands and every existing ordinary/triple string fixture remain green.
- [ ] Managed-span exclusion, exact sorted reporting, collision-before-mutation ordering, and atomic preservation remain unchanged.
- [ ] Paired helpers finish byte-identical; focused/full, installer, staged-skills, suite-manifest, Just-parse, Bash/ShellCheck, changelog-identity, and diff-hygiene checks pass.

## Pre-Flight

**Git:** ✓ No pre-existing changes outside `do-work/`; the approved REQ-162 queue edit, claim move, and checkpoint entry are lifecycle state.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh` passes (`launched: true`).
**Dependencies:** ✓ Bash, Python 3, Just 1.46.0, and ShellCheck are available; paired helpers begin byte-identical.

*Checked by work action*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Loaded the required prior-literal repairs and implementation guardrails, isolated the remaining defect to ordinary backtick state persistence in the two paired helpers, and mapped positive, exact-negative, inactive-context, parse, byte, and adjacent-contract evidence to the declared three-file write set.
- [x] **[APPLY]:** Captured production-helper RED with exact diagnostics and byte preservation, added Just-parseable positive and negative fixtures, then persisted byte `96` through the existing raw active-delimiter path byte-identically in both helpers without changing definition grammar or mutation ordering.
- [x] **[UNIFY]:** Confirmed the focused behavior, full contract regressions, standalone installer/staged-skills/suite-manifest contracts, paired and changelog identities, Just parsing, Bash syntax, warning-level ShellCheck, same-root/similar-pattern/protected-file audits, and diff hygiene all pass.

## Root Cause

`just_multiline_string_state()` recognized ordinary backticks while scanning one physical line but returned ordinary state only for bytes `34` and `39`. An unclosed byte `96` was discarded at the line boundary, so the next column-zero command payload reached `just_definition_name()` as if it were a top-level definition.

## Decisions

- **D-01 — Reuse the raw active-delimiter path.** Persist byte `96` alongside the existing ordinary quote bytes; the established active-state loop already closes on the next literal backtick, ignores backslash parity, and rescans the same-line suffix for close/reopen behavior.
- **D-02 — Preserve longest-first and inactive-context behavior.** Leave exact triple-backtick matching ahead of ordinary backticks and keep the top-level/comment/indent gate unchanged, so comments, recipe bodies, and closed same-line commands cannot seed cross-line state.
- **D-03 — Prove acceptance and rejection separately.** Require exact appended bytes plus before/after Just parsing for the positive, and an exact four-name sorted diagnostic plus pre-mutation byte identity for real collisions surrounding literal-only `kanban-static` payload.

## Implementation Summary

**Files changed:**
- `tools/replace-text-section.sh` (modified) — persists ordinary single-backtick state across physical lines through the existing raw delimiter path.
- `skills/do-work/tools/replace-text-section.sh` (modified) — byte-identical shipped-helper correction.
- `_dev/tests/contract-regressions.sh` (modified) — adds Just-parseable acceptance and exact byte-preserving real-collision fixtures.

**What was done:** Reserved-looking recipe and alias payload inside valid ordinary multiline backtick commands is now ignored, while same-line commands, inactive contexts, triple-backticks, prior string families, exact real-collision reporting, and pre-mutation preservation remain intact.

## Qualification

Passed — all three scoped implementation files are present and substantive, Scope and Implementation Summary sets match, P-A-U contains exactly one checked entry per phase, and the paired helper delta is byte-identical. Every Requirement and Instance traces to the persistent byte-96 state change or its direct production-helper fixtures. The data path remains `just_multiline_string_state()` → `just_definition_names()` → exact reserved-name intersection before replacement construction and `atomic_replace()`; definition grammar, managed-span removal, diagnostic sorting, and mutation ordering are unchanged. Mechanical qualification and debug/wiring audits pass, with no contamination from REQ-161/163 surfaces.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/install-suite-behavior.sh`; `bash _dev/tests/staged-skills-contract.sh`; `bash _dev/tests/suite-manifest-contract.sh`; paired-helper `cmp -s`; `bash -n` and `shellcheck -S warning` on all three scoped shell files; board-template `just --list`; changelog-mirror `cmp -s`; `git diff --check`; protected-file SHA-256 verification
**Result:** ✓ All passing.

**Red-green validation:**
- Valid ordinary physical-line backtick command: ✗ old helper rejected it with exact false collisions `kanban-summary, run-kanban` and preserved bytes → ✓ corrected helper inserts the exact managed section and the result parses with Just.
- Nearby real-collision control: ✗ old helper incorrectly added literal-only `kanban-static` → ✓ corrected helper reports exactly `kanban-summary, run-do-work-update, run-kanban, run-kanban-cli` and preserves the target byte-for-byte.

**New tests added:** Just-parseable production-helper fixtures for raw backtick closing after `\`, same-line close/reopen, inactive comment/body and closed same-line controls, exact insertion bytes, exact sorted diagnostics, and rejection preservation.

**Existing tests updated (cross-REQ impact):** `_dev/tests/contract-regressions.sh` extends the REQ-156/159 multiline-literal collision contract while retaining every prior string/command family and transaction-safety assertion.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-10T11:59:32Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — both production helpers accept Just-parseable ordinary multiline backtick payload with exact bytes while retaining exact real-collision diagnostics and byte-preserving rejection.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Just-parseable positives plus exact stderr and byte snapshots made a one-condition lexer correction independently provable in both acceptance and rejection directions.
- Reusing the existing raw active-delimiter path preserved suffix close/reopen behavior and kept root/shipped helpers byte-identical.
- Replaying the production helper against all six Just multiline delimiters, then parsing the result with Just itself, exposed a real reserved recipe that the line-oriented collision scan had silently missed.

**What didn't:**
- The earlier multiline-literal repair treated ordinary backticks as line-local without checking Just's accepted physical-line form, leaving the final byte-96 state boundary uncovered.
- Literal state was completed without carrying recipe-header state with it: the opening line of a multiline parameter default has the recipe name but no colon, while the closing line has the colon and was skipped as literal content. The helper therefore accepted a duplicate reserved recipe in the no-Just path.

**Worth knowing:**
- Ordinary backticks close on the next literal backtick even after a backslash; they must not inherit cooked double-quote escape parity, and exact triple-backticks must remain the longest-first opener.
- A handwritten lexer must test each literal in every grammar position it can occupy, not only as assignment payload. For Just, variable values prove false-positive suppression; recipe-parameter defaults prove real-definition retention. The fix retains the pending top-level header until the literal closes, then evaluates the complete header while preserving `:=` assignment exclusion.

## Orientation

The reserved-recipe collision scanner now handles every multiline string/command delimiter in both assignment payloads and recipe-parameter defaults: safe literal content stays ignored, while real reserved recipes are rejected before a no-Just installation can create a duplicate definition.
