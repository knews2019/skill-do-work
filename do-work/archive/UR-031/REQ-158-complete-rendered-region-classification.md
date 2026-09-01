---
id: REQ-158
title: "Review fix: Complete rendered-region classification in shipped Markdown references"
status: completed
completed_at: 2026-08-10T09:48:42Z
commit: 47b71fd
claimed_at: 2026-08-10T09:24:04Z
route: C
status_changed_at: 2026-08-10T09:20:51Z
domain: general
created_at: 2026-08-09T18:38:08Z
user_request: UR-031
addendum_to: REQ-154
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-rendered-region-classification
write_set:
  - _dev/tests/shipped-package-reference-contract.sh
kb_status: promoted
kb_entry: REQ-158-review-fix-complete-rendered-region-clas.md
---

# Review Fix: Complete Rendered-Region Classification in Shipped Markdown References

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Confirmed the parser-only boundary, mapped all four reproduced Instances to exact production-helper fixtures, and retained the downstream publication-policy contracts.
- [x] **[APPLY]:** Captured focused RED for effective-column indentation and paragraph context, escaped-backtick delimiters, and escaped-link URL fallback; then made one offset-preserving rendered-region mask authoritative before target discovery.
- [x] **[UNIFY]:** Audited the one-file implementation, checked the same-root pattern, and passed focused/full/distribution, syntax, ShellCheck, changelog-identity, and diff-hygiene verification.

## What

Make one rendered-versus-ignored decision govern every target-discovery path in the shipped-package Markdown guard for the already-approved escaped-link, indented-code, and inline-code classes. Done means these classes cannot recur as either release-blocking false positives or broken-link false negatives; this does not add a documentation policy or broaden the downstream reference contract.

## Context

Found during review of REQ-154. Its ordinary fixtures and all full contracts pass, but bounded production-helper probes show that independent masking and bare-URL scans disagree at valid Markdown boundaries.

## Requirements

- Apply the ignored escaped-link decision to relative and bare first-party URL discovery alike.
- Recognize Markdown indentation by effective columns and block context so tab-expanded indented code is ignored without hiding a four-space continuation in an active paragraph.
- Do not treat backslash-escaped backticks as inline-code delimiters; preserve ordinary exact-run code-span behavior.
- Add exact positive/negative production-helper fixtures for every instance below and retain the current topology, raw/blob, path-containment, changelog-identity, and full-contract behavior.

## Instances

- [ ] `_dev/tests/shipped-package-reference-contract.sh`: an escaped link containing a first-party URL is ignored by inline extraction but re-added by the bare-URL fallback.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: one-to-three leading spaces followed by a tab can form four columns of indented code but is scanned as live Markdown.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: a four-space continuation inside an active paragraph is rendered Markdown but is masked as code.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: backslash-escaped backticks are treated as a code span and hide a published link between them.

## Open Questions

- [x] REQ-154 improved the release guard, but its parser can still reject valid documentation or miss a broken published link in the approved escaped-link, indented-code, and inline-code classes. This is a review of review-generated work, so your consent is required before another implementation task enters the work loop. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: C** - Complex

**Reasoning:** The implementation surface is narrow, but the bug joins four Markdown lexical and block-context cases across every target-discovery path. A full plan and exploration are needed to preserve the existing publication-policy contracts while making one rendered-region decision authoritative.

**Planning:** Required

## Plan

1. Add named RED cases to `run_parser_fixtures()` for escaped relative and first-party URL destinations, one-to-three spaces plus a tab, four-space paragraph continuations versus genuine indented code, and escaped backticks versus exact-run code spans.
2. Make one offset-preserving rendered/ignored mask authoritative for every discovery path, including escaped-link regions and the bare first-party URL fallback.
3. Replace literal indentation-prefix checks with effective-column calculation and minimal paragraph/block state so indented code is masked only where Markdown permits it.
4. Make inline-code delimiter scanning honor odd/even backslash parity while retaining exact-run closing behavior.
5. Run the focused guard, full contracts, distribution contracts, Bash syntax, warning-level ShellCheck, changelog identity, and diff hygiene.

**Root-cause hypothesis:** The current helper makes rendered-region decisions in separate passes: indentation and backtick masking are incomplete, escaped-link destinations are rejected only by structural extraction, and the later bare-URL scan re-examines the whole masked string. One shared classification must govern all target discovery without changing downstream publication policy.

**Exact implementation scope:** Modify `_dev/tests/shipped-package-reference-contract.sh` only. Existing regression wiring already invokes it; publication policy, Markdown content, installers/updaters, ADR-019, UR-031, and owner-only lifecycle files stay out of builder scope.

**Plan validation:** Every Requirement and all four Instances map to the fixture task plus their corresponding masking/classification task. No orphan tasks were identified. The plan has five ordered tasks, which is at the quality-warning threshold, but four are inseparable RED/GREEN phases within one production helper and the fifth is verification; splitting would weaken the single classification contract.

*Generated by Plan agent*

## Exploration

The defect is confined to `_dev/tests/shipped-package-reference-contract.sh`; `_dev/tests/contract-regressions.sh` already invokes the focused guard. The exact production seams are `strip_markdown_code()` for offset-preserving masking, `punctuation_is_escaped()` for odd/even parity, `inline_link_targets()` for structural targets, `markdown_targets()` for reference definitions and bare first-party URLs, and `run_parser_fixtures()` for exact production-helper assertions.

The current mask recognizes only a literal leading tab or four spaces and lacks paragraph context; its inline-code pass treats every backtick run as a delimiter; and escaped-link syntax is rejected only by structural extraction while the later bare-URL regex still sees the destination. Classification should be derived from original line/block state, then applied without changing source length or newline offsets before every target path consumes it. Named paired fixtures must protect escaped relative/URL targets, every one-to-three-spaces-plus-tab form, live four-space paragraph continuations, genuine indented code, escaped backticks, even-parity controls, and ordinary exact-run spans.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — unify rendered-region masking across all target discovery and add production-helper regression fixtures.

**Files I will NOT touch:** `_dev/tests/contract-regressions.sh` (already wired), publication/topology/raw-blob/path-containment policy, published Markdown, installers/updaters, ADR-019, UR-031 input, or owner-only lifecycle/version/changelog files.

**Acceptance criteria (restated from REQ):**
- [ ] Escaped relative and first-party URL links are ignored consistently, with live/even-parity controls still discovered.
- [ ] One-to-three spaces plus a tab are classified by effective columns, while a four-space continuation in an active paragraph remains rendered and genuine block-context indented code stays ignored.
- [ ] Odd-parity escaped backticks cannot delimit inline code; ordinary and even-parity exact-run code spans retain their current behavior.
- [ ] Exact positive/negative production-helper fixtures cover every listed Instance without changing source length or newline offsets.
- [ ] Existing topology, raw/blob, path containment, changelog identity, focused, full-contract, and distribution behavior remains green.

## Pre-Flight

**Git:** ⚠ Pre-existing `decisions/records/adr-019-four-skill-suite-contract.md` edit is user-owned; preserve it and exclude it from this REQ's staging and diff judgments. The approved REQ-158/159/160 queue edits and UR-031 input edit are lifecycle/user state, not builder scope.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh` passes (`launched: true`).
**Dependencies:** ✓ Bash and Python 3 are available; the focused parser contract is local and network-free.

*Checked by work action*

## Root Cause

The guard made non-rendered-region decisions in independent partial passes. Block masking recognized only literal tab or four-space prefixes and had no paragraph state, inline-code scanning treated escaped backticks as delimiters, and escaped-link rejection lived only in structural extraction while the bare first-party URL fallback rescanned the unclassified destination. Those paths therefore disagreed about the same source offsets.

## Decisions

- **D-01 — Keep one offset-preserving classification seam.** Extend `strip_markdown_code()` to mask every approved non-rendered region before structural links, reference definitions, or bare first-party URLs are discovered; leave downstream topology and URL policy unchanged.
- **D-02 — Track only required Markdown block state.** Expand leading spaces and tabs to effective columns and retain minimal paragraph-versus-block state so indented code cannot interrupt an active paragraph, without introducing a Markdown dependency or a broader documentation policy.
- **D-03 — Apply existing escape parity to code delimiters.** Reuse odd/even backslash classification for opening and closing backtick runs while preserving exact-run matching, rather than adding a second escape rule.

## Implementation Summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified) — centralizes offset-preserving rendered-region masking, adds effective-column and paragraph/block classification, makes code-span delimiters escape-aware, masks escaped inline-link and escaped reference-definition regions before every discovery pass, and adds exact paired production-helper fixtures.

**What was done:** The shipped Markdown guard now feeds structural links, reference definitions, and bare first-party URLs from the same rendered text while retaining the existing downstream publication policy and source-offset invariants.

## Qualification

Passed on re-qualification attempt 2 — one scoped implementation file verified, all four Requirements and Instances traced, P-A-U confirmed, and Scope/Implementation Summary sets match. Attempt 1 exposed escaped reference-definition URLs still reaching the bare fallback; exact opening/closing-escape RED controls were added, the shared length-preserving mask now covers them, and focused probes plus full contracts pass.

## Testing

**Tests run:** `bash _dev/tests/shipped-package-reference-contract.sh`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/staged-skills-contract.sh`; `bash _dev/tests/suite-manifest-contract.sh`; `bash -n _dev/tests/shipped-package-reference-contract.sh`; `shellcheck -S warning _dev/tests/shipped-package-reference-contract.sh`; `cmp -s CHANGELOG.md skills/do-work/CHANGELOG.md`; `git diff --check`
**Result:** ✓ All passing.

**Red-green validation:**
- Effective-column/block-context fixture: ✗ old helper returned the three hidden space-plus-tab targets and omitted the live four-space continuation → ✓ exact expected continuation only.
- Escaped-backtick fixture: ✗ old helper hid the live link between escaped delimiters → ✓ odd-parity escaped delimiters stay live while ordinary/even-parity exact-run spans stay hidden.
- Escaped-link URL fixture: ✗ old bare fallback re-added the hidden first-party URL → ✓ hidden inline-link regions and live even-parity controls classify consistently.
- Re-qualification escaped-reference fixture: ✗ opening/closing escaped definitions leaked both hidden first-party URLs → ✓ both stay hidden while ordinary and even-parity reference controls remain live.

**New tests added:** Exact named production-helper fixtures for effective indentation and paragraph continuation, escaped backticks, escaped inline-link regions, and escaped reference-definition regions; every case retains source-length and newline-offset assertions.

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/shipped-package-reference-contract.sh` (from REQ-154): extended the same production-helper contract; existing topology, raw/blob, containment, changelog-identity, and parser cases remain unchanged and green.

*Verified by work action*

## Review

**Overall: 63%** | 2026-08-10T09:46:11Z

| Dimension | Score |
|-----------|-------|
| Requirements | 50% |
| Code Quality | 72% |
| Test Adequacy | 68% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- Rendered-region classification still misses escaped closing-bracket/opening-parenthesis link shapes and four-space continuations inside bullet or ordered-list paragraphs, allowing hidden first-party URLs to re-enter or live links to disappear — gate: user-visible → rerouted pending-answers as REQ-161

**Minor findings:** 0 (report only)
**Acceptance:** Partial — declared contracts pass, but production-helper probes reproduce same-class URL leaks and list-paragraph false negatives.
**Suggested testing:** 2 items
**Follow-ups created:** REQ-161; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Exact production-helper probes made disagreement between structural extraction and the bare-URL fallback observable without changing downstream policy.
- Length-preserving masks plus paired live/hidden controls protected offsets and caught over-masking during re-qualification.

**What didn't:**
- Implementing only the reproduced opening-bracket and top-level paragraph forms left the same classification rule incomplete at other link delimiters and list-item paragraph contexts; REQ-161 records the remaining sweep for consent.

**Worth knowing:**
- A rendered-region guard needs delimiter-complete and block-context-complete tables, not one fixture per syntax family. Full live-corpus contracts can stay green while nearby production-helper variants still fail.

## Orientation

The shipped-reference guard now shares one offset-preserving Markdown classification across its target-discovery paths for effective indentation, top-level continuations, escaped backticks, escaped opening links, and escaped reference definitions. Review isolated the remaining delimiter/list-context facets in consent-gated REQ-161; publication policy and four-module topology are unchanged.
