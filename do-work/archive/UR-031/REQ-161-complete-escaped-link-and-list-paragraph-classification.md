---
id: REQ-161
title: "Review fix: Complete escaped-link and list-paragraph classification"
status: completed
completed_at: 2026-08-10T11:31:29Z
commit: ad3f8bd
claimed_at: 2026-08-10T10:59:26Z
status_changed_at: 2026-08-10T09:53:36Z
domain: general
created_at: 2026-08-10T09:46:11Z
user_request: UR-031
addendum_to: REQ-158
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-rendered-region-classification
write_set:
  - _dev/tests/shipped-package-reference-contract.sh
kb_status: promoted
kb_entry: REQ-161-review-fix-complete-escaped-link-and-lis.md
---

# Review Fix: Complete Escaped-Link and List-Paragraph Classification

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Confirmed the one-file production-helper boundary, mapped the four Instances to four separately named fixtures, and preserved the downstream publication-policy and offset contracts.
- [x] **[APPLY]:** Added the fixtures first, captured focused RED for both escaped delimiter seams and both list contexts, then completed the shared mask and list paragraph-state classification with the smallest local changes.
- [x] **[UNIFY]:** Swept every escape-parity call, mask consumer, and list/indentation state branch in the same root and passed focused, full, distribution, syntax, ShellCheck, changelog-identity, scope, protected-file, and diff-hygiene checks.

## What

Complete the shared rendered-versus-ignored classification for the remaining escaped inline-link delimiter shapes and list-item paragraph continuations. Done means the approved Markdown class cannot recur through a different structural delimiter or block context; this remains a test-only release-guard correction and does not change the downstream publication policy.

## Context

Found during review of REQ-158. Its named regressions and full contracts pass, but production-helper probes show that closing-bracket/opening-parenthesis escapes still disagree with the bare-URL fallback and that list-item paragraph continuations are treated as indented code.

## Requirements

- Apply odd/even escape parity consistently at every inline-link structural delimiter so ignored escaped links cannot re-enter through bare first-party URL discovery.
- Recognize bullet and ordered-list paragraph continuation context before classifying four effective columns as indented code.
- Add exact positive/negative production-helper fixtures and retain all REQ-158, topology, raw/blob, path-containment, changelog-identity, focused, full-contract, and distribution behavior.

## Instances

- [ ] `_dev/tests/shipped-package-reference-contract.sh`: `[label\](first-party URL)` is rejected structurally but its URL is re-added by the bare fallback.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: `[label]\(first-party URL)` is rejected structurally but its URL is re-added by the bare fallback.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: a four-space continuation inside a bullet-list paragraph is rendered but masked as indented code.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: a four-space continuation inside an ordered-list paragraph is rendered but masked as indented code.

## Open Questions

- [x] REQ-158 fixes escaped opening brackets and top-level paragraph continuation, but escaped closing brackets or opening parentheses still let first-party URLs re-enter through the bare-URL scan, and four-space continuations inside bullet or ordered-list paragraphs are still hidden as code. This is another review of review-generated work, so the cascade-depth rule requires your consent before completing the already-approved rendered-region behavior. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: C** - Complex

**Reasoning:** The change is confined to a test-only release guard, but it completes a handwritten Markdown classifier across structural escape parity and list-block continuation state while preserving publication policy and prior parser contracts. Fresh planning and exploration are warranted.

**Planning:** Required

## Plan

1. **Capture RED in the production helper before changing classifier logic.** In `_dev/tests/shipped-package-reference-contract.sh`, add four separately named `run_parser_fixtures()` cases for odd-parity escaped closing brackets, odd-parity escaped opening parentheses, bullet-list paragraph continuations, and ordered-list paragraph continuations. Pair each with even-parity/live and genuine-code controls; retain exact-target, source-length, and newline-offset assertions. Run the focused guard before implementation and record RED.
2. **Complete inline-link escape classification at the shared mask.** Generalize `escaped_inline_link_end()` and its masking pass so a complete inline-link-shaped region is ignored when odd escape parity occurs at any link-forming delimiter—opening `[`, candidate closing `]`, or destination-opening `(`. Reuse `punctuation_is_escaped()`; preserve escaped label content, nesting, complete-region detection, and even-parity rendered controls. Leave the bare first-party URL policy unchanged.
3. **Preserve list-item paragraph state before indentation classification.** Refine `line_opens_paragraph()` and the existing `paragraph_active` transition so bullet and one-to-nine-digit ordered markers with same-line paragraph content establish an active list paragraph before the four-column indented-code branch. Blank lines, empty markers, fences, and blank-separated genuine indented code must remain non-rendered/resetting.
4. **Perform the same-root sweep after GREEN.** Audit every inline delimiter escape-parity call, every `strip_markdown_code()` consumer, and list-marker/effective-indentation branch in the same file. Retain prior opening-bracket, reference-definition, backtick, paragraph, and genuine-code controls; record any separate defect as a Discovered Task.
5. **Verify preservation and distribution.** Run the focused shipped-reference guard, full contract regressions, staged-skills and suite-manifest contracts, Bash syntax, warning-level ShellCheck, changelog identity, scoped diff audit, and `git diff --check`.

**Root-cause hypothesis:** The shared rendered-region mask is authoritative but incomplete at two seams. `escaped_inline_link_end()` is gated on an escaped opening bracket and only accepts an unescaped closing bracket followed by an unescaped parenthesis, so odd escapes at `]` or `(` are rejected structurally without masking the first-party URL from the fallback. Separately, `line_opens_paragraph()` treats every list marker as non-paragraph content, clearing `paragraph_active` before the following four-column continuation reaches indentation classification.

**Exact implementation scope:** Modify `_dev/tests/shipped-package-reference-contract.sh` only. Keep publication/target-resolution policy, suite/runtime/install/update behavior, shipped Markdown, manifest, skills/tools, protected ADR/UR files, and owner-only lifecycle/version/changelog files outside builder scope.

**Plan validation:** All three Requirements and all four Instances map to named RED fixtures, one minimal classifier correction per root-cause seam, the same-pattern sweep, and explicit verification. Five tasks reach the quality-warning threshold, but they are inseparable RED/GREEN/sweep/verification phases inside one production helper; no orphan task or uncovered requirement exists.

*Generated by Plan agent*

## Exploration

The defect remains confined to `_dev/tests/shipped-package-reference-contract.sh`. The production path is `strip_markdown_code()` → `markdown_targets()` → `inline_link_targets()` plus the bare first-party URL fallback, and `run_parser_fixtures()` already exercises these exact helpers with exact normalized target order, source-length preservation, and unchanged newline offsets.

Escape parity is centralized in `punctuation_is_escaped()`. The masking helper currently requires an escaped opening bracket and therefore cannot mask complete link-shaped regions whose candidate closing bracket or destination-opening parenthesis has odd parity; structural extraction rejects them, but the still-visible first-party URL reaches the fallback. List state is held only in `paragraph_active`; `line_opens_paragraph()` rejects every list marker even when same-line item text opens a paragraph, so the next four-column continuation is misclassified as indented code.

The minimal correction is to recognize a complete same-line link-shaped region and mask it when any structural delimiter has odd parity, while distinguishing escaped bracket content from the candidate delimiter and preserving nesting/even-parity forms. Supported nonempty bullet and one-to-nine-digit ordered markers should establish paragraph state; empty markers, blanks, fences, and blank-separated four-column code must not.

Repository-wide search found no second implementation of these helpers or the bare first-party URL scan. `_dev/tests/contract-regressions.sh` already invokes the focused guard. Existing opening-bracket, reference-definition, top-level paragraph, backtick, nesting, and destination-normalization fixtures remain preservation controls.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — complete inline-link delimiter-parity masking, preserve list-item paragraph state, and add four exact production-helper fixtures.

**Files I will NOT touch:** `_dev/tests/contract-regressions.sh` (already wired), shipped Markdown, publication/target-resolution policy, suite manifest, skills/tools/runtime/install/update behavior, fixture sidecars, protected ADR/UR inputs, or owner-only lifecycle/version/changelog files.

**Acceptance criteria (restated from REQ):**
- [ ] Odd-parity escapes at opening `[`, candidate closing `]`, or destination-opening `(` mask the complete link-shaped region before every target-discovery path; zero/even parity remains discoverable.
- [ ] Escaped label content, nested labels/parentheses, incomplete forms, escaped reference definitions, and existing destination normalization retain current behavior.
- [ ] Nonempty `-`, `+`, `*`, and one-to-nine-digit `[.)]` list items establish paragraph state so their four-column continuations remain rendered.
- [ ] Empty markers, blank lines, fences, and blank-separated four-column code remain non-paragraph/masked.
- [ ] Four named fixtures assert exact normalized targets plus unchanged source length and newline offsets.
- [ ] Focused, full-contract, distribution, syntax, warning-level ShellCheck, changelog-identity, scoped-diff, and diff-hygiene verification passes.

## Root Cause

The shared rendered-region mask was authoritative but incomplete at both reported seams. `escaped_inline_link_end()` entered only when the opening bracket was escaped and accepted only an unescaped closing bracket immediately followed by an unescaped opening parenthesis; an odd escape at either later structural delimiter therefore made structural extraction reject the link while leaving its first-party URL visible to the bare fallback. Separately, `line_opens_paragraph()` rejected every list-marker line, including markers with same-line paragraph text, so `paragraph_active` was cleared before the following four-column continuation reached indentation classification.

## Decisions

- **D-01 — Classify all inline-link delimiters at the shared mask.** Reuse `punctuation_is_escaped()` for the opening bracket, candidate closing bracket, and destination-opening parenthesis, and mask only a complete same-line link-shaped region when at least one delimiter has odd parity. Escaped label content, nested labels and destinations, incomplete forms, and zero/even parity remain unmasked for the existing discovery paths.
- **D-02 — Derive list paragraph state from item content.** Strip one supported bullet or one-to-nine-digit ordered marker inside `line_opens_paragraph()` and classify its nonempty same-line content with the existing block rules. Empty markers, comments, fences, nested block openers, blanks, and blank-separated indented code still clear or avoid paragraph state.
- **D-03 — Keep publication policy and offsets unchanged.** Confine the correction and four exact fixtures to the production helper; retain length/newline assertions and leave normalization, target resolution, topology, first-party URL, path-containment, and changelog-identity policy untouched.

## Implementation Summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified) — classifies odd escape parity at every inline-link structural delimiter in the shared mask, preserves list-item paragraph continuation state, and adds exact production-helper fixtures with parity, block-context, and offset controls.

**What was done:** The test-only shipped-reference guard now hides complete escaped link-shaped regions from every target-discovery path and keeps four-column continuations rendered inside nonempty bullet and ordered-list paragraphs, without changing publication policy or downstream target resolution.

## Qualification

Passed — one scoped implementation file is substantive and present in the diff; all three Requirements and four Instances trace to production-helper logic plus named fixtures; Scope and Implementation Summary match; exactly one checked PLAN/APPLY/UNIFY entry is present. The shared masked-text flow remains live through structural extraction and the bare first-party URL fallback, with no stub, policy change, debug artifact, or undeclared project-file touch.

## Testing

**Tests run:** `bash _dev/tests/shipped-package-reference-contract.sh`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/staged-skills-contract.sh`; `bash _dev/tests/suite-manifest-contract.sh`; `bash -n _dev/tests/shipped-package-reference-contract.sh`; `shellcheck -S warning _dev/tests/shipped-package-reference-contract.sh`; scope-drift check; changelog/version mirror checks; protected-file hashes; `git diff --check`
**Result:** ✓ All passing.

**Red-green validation:**
- Odd escaped closing bracket and opening parenthesis: ✗ each new fixture leaked exactly its hidden first-party URL while retaining live controls → ✓ the complete escaped link-shaped region is masked and even-parity/ordinary controls remain discoverable.
- Bullet and ordered list paragraph continuations: ✗ each fixture returned no live continuation targets → ✓ all approved marker forms return their four-column continuation targets while empty-marker, fence, blank, and genuine-code controls remain hidden.

**New tests added:** Four named production-helper fixture groups cover structural escape parity and bullet/ordered paragraph contexts with exact targets plus unchanged source-length/newline-offset assertions.

**Existing tests updated (cross-REQ impact):** `_dev/tests/shipped-package-reference-contract.sh` extends the REQ-154/REQ-158 parser contract; all earlier fixtures and downstream topology/publication/distribution checks remain unchanged and green.

*Verified by work action*

## Review

**Overall: 63%** | 2026-08-10T11:29:43Z

| Dimension | Score |
|-----------|-------|
| Requirements | 67% |
| Code Quality | 68% |
| Test Adequacy | 58% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- Inline-link masking false-masks escaped opening brackets used as label content and leaves even-parity opening-parenthesis relative targets undiscoverable; first-party fixtures conceal the latter through bare-URL fallback — gate: user-visible → consolidated into generation≥2 pending-answers sweep REQ-163
- List paragraph classification misses fence runs with attached info strings and nested tilde fences, allowing code links to be scanned as published Markdown — gate: user-visible → consolidated into generation≥2 pending-answers sweep REQ-163

**Minor findings:** 0 (report only)
**Acceptance:** Partial — the four named regressions and every declared suite pass, but actual-helper adversarial probes reproduce remaining same-class false negatives and false positives.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-163; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Exact production-helper RED cases closed the four approved defects without changing downstream publication or target-resolution policy.
- Pairing odd/even parity and rendered/code controls made the shared-mask intent directly testable while preserving offsets.

**What didn't:**
- First-party URL controls could pass through the fallback even when structural relative-link extraction was still broken, so they overstated parity completeness.
- Bare list-fence controls did not exercise attached info strings or escaped opening brackets nested inside live labels; the independent matrix exposed both blind spots.

**Worth knowing:**
- Parser completeness needs paired relative and first-party targets for every delimiter form, plus context-sensitive controls that distinguish a structural escape from the same punctuation inside label content. REQ-163 records the remaining bounded variants for consent.

## Orientation

The shipped-reference release guard now handles the four approved escaped-closing/opening-parenthesis and bullet/ordered continuation defects through its shared test-only Markdown classifier. Publication policy and runtime behavior are unchanged; review isolated the remaining relative-parity, label-content, and list-fence-info variants in consent-gated REQ-163.
