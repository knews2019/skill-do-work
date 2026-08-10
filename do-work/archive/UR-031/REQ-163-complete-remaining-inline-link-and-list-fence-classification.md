---
id: REQ-163
title: "Review fix: Complete remaining inline-link and list-fence classification"
status: completed
completed_at: 2026-08-10T12:43:05Z
commit: c9d1acd
status_changed_at: 2026-08-10T12:34:02Z
claimed_at: 2026-08-10T12:34:36Z
route: A
domain: general
created_at: 2026-08-10T11:30:40Z
user_request: UR-031
addendum_to: REQ-161
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-rendered-region-classification
write_set:
  - _dev/tests/shipped-package-reference-contract.sh
kb_status: pending
kb_entry:
---

# Review Fix: Complete Remaining Inline-Link and List-Fence Classification

## What

Finish the test-only rendered-region classifier so delimiter parity, label context, and list-item fences cannot disagree with target discovery. Done means relative and first-party targets classify consistently and fenced examples inside list items remain ignored, without changing publication, topology, raw/blob, path-containment, runtime, or documentation policy.

## Context

Found during independent review of REQ-161. Its four reproduced cases and declared suites pass, but actual production-helper probes found nearby variants that still hide live relative references or scan fenced code as published Markdown. Because REQ-161 is itself review-generated, the cascade-depth rule requires user consent before this non-critical continuation becomes claimable.

## Requirements

- Make structural extraction preserve zero/even-parity destination-opening-parenthesis forms for relative targets, not only first-party URLs recovered by the bare fallback.
- Keep escaped opening brackets used as label content from causing the containing live link to be masked.
- Recognize backtick and tilde list-item fences with attached info strings while preserving approved list-paragraph continuations.
- Add exact production-helper fixtures for every instance and retain all publication-policy, offset, focused, full-contract, and distribution behavior.

## Instances

- [ ] Even two/four-backslash destination-opening-parenthesis forms retain relative targets, not only first-party URLs recovered by fallback.
- [ ] An escaped opening bracket used as label content does not cause the containing live link to be masked.
- [ ] Bullet and one-to-nine-digit ordered list items recognize backtick/tilde fences with attached info strings, including nested tilde fences.

## Open Questions

- [x] REQ-161 fixes its four reproduced cases, but the remaining parity, label-content, and list-fence variants can still hide broken relative references or scan code as published Markdown. This is a non-critical follow-up to review-generated work, so the cascade-depth rule requires your approval instead of silently extending the queue. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: A** - Simple

**Reasoning:** The review already isolated three reproducible variants to one production test helper and supplied exact preservation boundaries. The fix can proceed directly as a focused bug correction with regression fixtures.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Loaded REQ-161, the bug-fix specification, and implementation guardrails; mapped the three remaining variants to exact production-helper fixtures and a one-file classifier correction.
- [x] **[APPLY]:** Added exact fixtures first and captured all three failures, then completed even-parity link discovery, live-label region skipping, and list-item fence state without changing target publication or resolution policy.
- [x] **[UNIFY]:** Reviewed the complete one-file diff and same-root helper surface; focused, full-contract, installer, staged-skills, suite-manifest, updater, Bash, ShellCheck, and diff-hygiene checks pass with no debug artifacts or unrelated project-file changes.

## Root Cause

The mask already preserved zero and even escape parity before a destination-opening parenthesis, but structural target extraction required the parenthesis immediately after the label, so relative targets behind two or four backslashes disappeared while first-party URLs were rescued by the unrelated bare-URL fallback. The mask also reconsidered every opening bracket after declining to mask a complete live link, allowing an escaped bracket inside that link's label to be treated as a new escaped-link start and mask the containing target. Finally, fenced-block state recognized only top-level fence markers; list-item markers with attached info strings opened paragraph state instead, so their code contents remained publishable and nested closers could not restore list-paragraph continuation state.

## Decisions

- **D-01 — Classify complete inline-link regions before masking.** Return both the region end and whether a structural delimiter is escaped, mask only the escaped form, and skip the full live region so escaped punctuation inside its label is not reconsidered independently.
- **D-02 — Resolve the actual destination opener after escape runs.** Structural extraction skips the backslash run after a label and accepts the parenthesis only when its parity is zero or even; odd parity remains masked by the shared region classifier.
- **D-03 — Track list-item fences at their content indentation.** Recognize bullet and one-to-nine-digit ordered fence openers with attached info strings, require compatible closers, and resume paragraph state after a list fence so approved indented continuations remain rendered. Root-level four-column list-shaped code remains masked.

## Implementation Summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified) — completes relative even-parity link extraction, prevents escaped label content from masking live links, tracks backtick and tilde list-item fences, and adds exact fixtures with boundary and offset controls.

**What was done:** The shipped-reference release guard now classifies all three remaining rendered-region variants consistently while preserving target normalization, publication policy, source length, newline offsets, list-paragraph continuations, and adjacent distribution contracts.

## Qualification

Passed — the single declared implementation file is present, substantive, and fully represented in the diff; all three Requirements and Instances trace to classifier logic plus exact fixtures; P-A-U contains exactly one completed entry per phase. The masked Markdown still flows through structural target extraction and the bare first-party fallback with no stub, policy change, debug artifact, or unrelated project-file touch. Mechanical qualification passed, and same-root search found no second implementation of these helpers.

## Testing

**Tests run:** `bash _dev/tests/shipped-package-reference-contract.sh`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/install-suite-behavior.sh`; `bash _dev/tests/staged-skills-contract.sh`; `bash _dev/tests/suite-manifest-contract.sh`; `bash _dev/tests/update-script-behavior.sh`; `bash -n _dev/tests/shipped-package-reference-contract.sh`; `shellcheck -S warning _dev/tests/shipped-package-reference-contract.sh`; `git diff --check`
**Result:** ✓ All passing.

**Red-green validation:**
- Zero/even destination-opening-parenthesis parity: ✗ the old helper found only the zero-backslash relative target and lost the two- and four-backslash forms → ✓ all three are structurally extracted while the odd-parity control remains hidden.
- Escaped opening brackets inside live labels: ✗ the old mask removed both containing live relative targets → ✓ both remain discoverable while a genuinely escaped outer link stays hidden.
- List-item fences with attached info strings: ✗ the old classifier published a fenced-code target and lost all five live continuation controls → ✓ bullet, ordered, nine-digit, and nested fences hide their code and preserve every post-fence/list-paragraph target.

**New tests added:** Three exact production-helper fixture groups cover relative parity, escaped label content, bullet/ordered/nested list fences, root-level indented-code preservation, normalized target order, source length, and newline offsets.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-10T12:42:01Z

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
**Acceptance:** Pass — production-helper execution proves relative parity, escaped label content, and bullet/ordered/nested list fences behave as specified while all adjacent contracts remain green.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Relative targets beside first-party controls exposed structural extraction gaps that the bare-URL fallback had previously hidden.
- Treating a complete live link as one region and tracking list fences at their content column fixed the context errors without altering downstream target policy.

**What didn't:**
- Earlier first-party-only parity fixtures could pass through the fallback while relative links remained broken, and scanning every bracket independently misread escaped label content as a new link.
- Bare list-fence controls did not exercise attached info strings, ordered-marker width, or nested fence indentation.

**Worth knowing:**
- Markdown classifier fixtures should pair relative and first-party targets for every delimiter form, and list-fence tests need the opener, matching content indentation, closer, and a live post-fence continuation in the same case.

## Orientation

The shipped-reference release guard now classifies relative delimiter parity, escaped bracket content inside live labels, and nested list-item fences consistently; publication and target-resolution policy remain unchanged.
