---
id: REQ-197
title: Normalize completed-work presentation target IDs
status: completed-with-issues
claimed_at: 2026-08-15T18:49:02Z
route: B
completed_at: 2026-08-15T19:20:09Z
commit: 89d068e
kb_status: promoted
kb_entry: REQ-197-normalize-completed-work-presentation-ta.md
domain: general
created_at: 2026-08-15T16:34:08Z
user_request: UR-042
addendum_to: REQ-189
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
maintenance: true
write_set: [skills/do-work-toolbox/actions/completed-work-presentation-reference.md, skills/do-work-toolbox/actions/present-work.md, _dev/tests/contract-regressions.sh]
---

# Review Fix: Normalize Completed-Work Presentation Target IDs

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Use exploration to identify the two presentation entry paths and the existing canonical resolver, then lock their inheritance at the current presentation contract seam.
- [x] **[APPLY]:** Added focused assertions first, observed semantic failures, then updated only the shared presentation resolver and migration-only portfolio dispatcher.
- [x] **[UNIFY]:** Reviewed the three-file diff, preserved archive/status and side-effect gates, ran the focused suite to GREEN, and passed `git diff --check`.

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

## Triage

**Route: B — Medium**

**Reasoning:** The defect and acceptance outcome are precise, but exploration is needed to locate every canonical token-grammar seam and decide whether behavior is prose-prescribed or executable before declaring the exact write boundary.

**Planning:** Not required — exploration-guided implementation.

## Exploration

- `skills/do-work/actions/work-reference.md` → **Target ID Resolution** is already the canonical source for case-insensitive prefixes, numeric-value matching, canonical stored IDs, and canonicalize-before-lookup. It should remain unchanged and be cited rather than copied.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` is the sole item resolver used by `ai-report` and `present-video`; its terminal-success resolver currently performs exact archive matching without first inheriting the token grammar.
- `skills/do-work-toolbox/actions/present-work.md` Step 1 is the separate migration-only entry path. It should recognize one item token through the shared grammar and print both replacements with the supplied spelling unchanged, without archive lookup, delegation, prompting, or writes.
- The replayable seam is the existing presentation-contract Python block in `_dev/tests/contract-regressions.sh`: first assert both citations/order and the canonical equivalent spellings, then apply the two action-doc changes.
- Preserve archive-only lookup, ambiguity rejection, terminal-success filtering, and the side-effect-free migration branch. Exclude the canonical resolver, schemas, executable tooling, consumer actions that already delegate, and portfolio publication.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (modify) — inherit the canonical token grammar before target lookup
- `skills/do-work-toolbox/actions/present-work.md` (modify) — apply the shared grammar to item-specific migration dispatch while preserving the supplied token
- `_dev/tests/contract-regressions.sh` (modify) — add replayable RED/GREEN assertions for both entry paths

**Files I will NOT touch:** `skills/do-work/actions/work-reference.md`; `ai-report.md`; `present-video.md`; archive/schema/tooling files; portfolio generation or publication mechanics; or unrelated presentation tests.

**Acceptance criteria:**
- [x] The shared presentation resolver cites and applies Target ID Resolution before any UR/REQ archive lookup.
- [x] The `present-work` item-dispatch path accepts case-insensitive, numeric-equivalent tokens and preserves the supplied spelling in both replacement commands.
- [x] Archive-only and terminal-success gates remain unchanged.
- [x] Focused RED/GREEN assertions pass; canonical verification remains the integration gate.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Made both completed-work presentation ID entry paths inherit the canonical Target ID Resolution grammar before lookup or migration dispatch, preserved caller-specific search/write behavior, and added replayable assertions for equivalent input spellings and supplied-token output.

## Testing

**Tests run:** baseline `bash _dev/tests/contract-regressions.sh`; test-only RED `bash _dev/tests/contract-regressions.sh`; GREEN `bash _dev/tests/contract-regressions.sh`; `git diff --check -- <three scoped files>`

**Result:** ✓ Focused suite passed after implementation.

**Red-green validation:** ✗ new assertions failed on both missing Target ID Resolution citations, equivalent-token acceptance, pre-lookup ordering, and supplied-token preservation → ✓ both action paths now cite and apply the shared grammar and every focused assertion passes.

**Remediation validation:** ✗ Attempt 1's tests accepted semantic negations and required duplicated caller examples → ✓ callers now retain only named inheritance plus caller-specific application, canonical grammar is checked once at its source, copied grammar is rejected, and both `Do not canonicalize` / `Do not recognize` mutations fail.

*Verified by work action*

## Qualification

Passed — all three scoped files are substantive and wired; Implementation Summary matches Scope; P-A-U is complete; the canonical grammar remains owned by `work-reference.md`; both callers preserve their own resolution/dispatch gates; shell syntax and `git diff --check` pass.

## Review — Attempt 1

**Overall: 50%** | 2026-08-15T19:06:02Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 70% |
| Test Adequacy | 40% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings:**
- Both callers restate the canonical grammar and the tests require those copies. Negating the application directives still passes, while removing local examples but retaining the source citation fails; the tests lock example prose rather than inheritance behavior. — gate: rule-change → remediation required before any generation-2 follow-up disposition.

**Minor findings:** 0
**Acceptance:** Fail — the intended inputs are described, but canonical inheritance and the Finding-Closure Ratchet are not proven.
**Suggested testing:** semantic-negation and source-seam-only mutations.

*Reviewed by review-work action*

## Remediation

**Attempt 1 diagnosis:** Confirmed that caller-local examples forked the canonical token grammar and example-presence tests were semantically blind.

**Single remediation:** Removed copied grammar and membership definitions, retained active named application before each caller-specific lookup/dispatch boundary, moved grammar assertions to the canonical source section, rejected copied caller grammar, and added semantic-negation mutation guards.

**Evidence:** The focused suite rejected both duplicated-grammar RED state and temporary `Do not canonicalize` / `Do not recognize` mutations, then passed after restoring the active directives. Scoped `git diff --check` passed.

## Qualification — Attempt 2

Passed — the same three-file scope remains exact; the canonical source is asserted once; caller tests now cover citation, active application, ordering, local consequences, copied-grammar rejection, supplied-token preservation, and negation mutations.

## Review — Attempt 2

**Overall: 50%** | 2026-08-15T19:19:21Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 95% |
| Test Adequacy | 50% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings:**
- The Target ID regression guard remains semantically and positionally porous: word-substring matching accepts “read without applying,” `present-work` has no pre-dispatch order assertion, and its copied-grammar rejection omits membership definitions. This leaves the captured GREEN and Finding-Closure Ratchet unproven. — gate: rule-change → remediation exhausted; rerouted pending-answers as REQ-203.

**Minor findings:** 0
**Acceptance:** Fail — current instructions are correct, but the required finding-closure regression proof remains incomplete after the single remediation attempt.
**Suggested testing:** 1 focused semantic/order mutation matrix.
**Follow-ups created:** REQ-203; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Exploration found the single canonical grammar and kept the product change limited to two callers plus one contract seam.
- Mutation testing exposed semantic blindness that ordinary GREEN runs could not reveal.

**What didn't:**
- Attempt 1 copied examples into callers and tested for their presence, creating the drift the reference contract forbids.
- The remediation removed the duplication but still used a substring-positive assertion; “read without applying” survived, so the single remediation attempt did not fully close the review finding.

**Worth knowing:** When an instruction caller inherits a canonical contract, test the source definition once and test each caller's active, ordered application with adversarial negations. Positive keyword presence alone is not evidence that the directive remains operative.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

Completed-work presentation callers now point to the shared Target ID Resolution contract before their archive lookup or migration dispatch, so equivalent case and zero-padding forms follow one canonical rule. The current prose is correct; REQ-203 records the remaining consent-gated test-hardening decision.
