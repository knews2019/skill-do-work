---
id: UR-038
title: Stabilization plan v2 — ratchet, rubric home, scripts promotion
created_at: 2026-08-11T13:55:58Z
requests: [REQ-170, REQ-171]
word_count: 1
---

# Stabilization plan v2 — ratchet, rubric home, scripts promotion

## Summary

User approved the reworked stabilization plan discussed in-session ("how can we do it? ultrathink" → plan → "capture"). The plan was drafted against the then-pending UR-036/UR-037 batch as one new REQ plus four addenda to pending REQs — but the work loop processed the entire batch between plan approval and this capture (REQ-165..169 built, reviewed, archived; commits `a45d5c4`, `6538bdd`, `1a27c07`, `8703b66`, `063bb88`). The capture was therefore re-anchored against the shipped state: the pending-REQ addenda are moot, and the two deltas the shipped batch does not cover land as two queue REQs.

## Extracted Requests

| Item | Target | Change |
|---|---|---|
| REQ-170 (new) | work-reference, review-work, capture, coding-guardrails, validate-feedback | Finding-closure ratchet (test-or-deletion) + canonical earned-defense rubric home; validate-feedback's shipped restatement (REQ-169) condensed to application + citation |
| REQ-171 (new, `addendum_to: REQ-165`) | prescribed-shell primitives, `scripts/`, `_dev/tests/` | Promote multi-line primitives from prose (REQ-167's canonical guide) to shipped scripts with fixture-repo execution tests; lint harness (REQ-165) narrows to inline residue |

Dropped from the original plan shape: addenda to REQ-167/168/169 — their deltas either folded into REQ-170's scope (rubric single-home, audit citation) or are superseded by REQ-171 (script promotion subsumes the prose-narrowing).

## Batch Constraints

- The dividing line for script promotion is "does this block contain logic that can be wrong," not line count; one-liners and illustrative fragments stay inline.
- Every call site keeps a one-line intent statement so actions still work as standalone pasted prompts; canonical logic lives in the shipped script (floor respected: read/write files + run shell).
- The rubric and the ratchet are each stated **once**, pointed to elsewhere — every added rule is itself surface; the campaign must not become its own complexity.
- No metric machinery: "stable" = two consecutive review cycles where nothing above nitpick is accepted, read off validate-feedback's existing summary tables.
- Go-owned capabilities (atomic REQ reservation) get no shell twin — the script layer is for shell-portable primitives only.
- Existing ratchets survive the moves: `prescribed-shell-canonicalization.sh` re-targets script-as-home; validate-feedback's Surface-cost verdict contract stays enforced.

## Full Verbatim Input

capture

### Conversation context (same session, condensed)

The approved plan, in the user's flow: reviews keep returning 3–5 findings; the three levers ranked structure > ratchet > gate. Structure: promote multi-line prescribed shell (screenshot install block, manual next-REQ scan, curl download-and-rename, merge-safe first-parent diff) into per-package `scripts/` files with `_dev/tests/` fixture-repo tests — lint alone cannot catch the semantic trap classes (pipefail dead fallback, porcelain collapse, curl partials) that motivated the campaign; the session-start hook fix (REQ-166) is the proof of shape. Ratchet: a finding-origin REQ closes only with a named regression test or a deletion — canonical text in work-reference, enforced at the review-work gate, GREEN criterion sharpened at capture. Gate: the earned-defense rubric ("what incident earned this, and is the fix still cheaper than the surface it added?") lives once in coding-guardrails.md; validate-feedback, review-work, and the defensive-surface audit point at it.

---
*Captured: 2026-08-11T13:55:58Z (re-anchored against shipped batch state before commit, same day)*
