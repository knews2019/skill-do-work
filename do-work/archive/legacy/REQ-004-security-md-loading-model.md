---
id: REQ-004
title: "Code review: resolve crew-members/security.md loading-model orphan status"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T19:05:00Z
completed_at: 2026-05-29T19:06:00Z
commit: 99a8ad3
route: B
review_generated: true
source: code-review
scope: CLAUDE.md, crew-members/security.md, actions/work.md, actions/code-review.md
---

# Code Review Fix: Resolve security.md Loading-Model Orphan Status

## What

`crew-members/security.md` is in the worst-of-both-worlds state: documented as a domain crew with JIT_CONTEXT triggers, but the loading model treats it as not-a-crew.

- CLAUDE.md Agent Rules section (L150-164) enumerates every crew-member with its trigger — `security.md` is **not listed**.
- `crew-members/security.md:3` JIT_CONTEXT claims it loads for "any REQ with domain: security".
- `actions/work.md` Step 6 (L507-518) never loads it.
- `actions/work.md` canonical domain enum (L194) only includes `frontend | backend | ui-design | general` — no `security` path exists.
- Only `actions/code-review.md:159` references it (conditionally: "If `crew-members/security.md` exists, load it…").

So the JIT_CONTEXT describes a load trigger that **doesn't exist in the codebase**.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`).

## Requirements

Pick one of two paths:

**Option A — Security is a real crew-member** (recommended; resolves both this finding and Architecture-F8):
- Expand `actions/work.md`'s canonical domain enum (L194) to include `security`.
- Add an entry to CLAUDE.md Agent Rules section: "`security.md` — loaded when REQ frontmatter `domain: security`, or when the REQ description references auth/crypto/input handling".
- Wire `actions/work.md` Step 6 (L511-518) to load `security.md` under the documented conditions, mirroring how `testing.md` loads on `tdd: true` or `domain: testing`.

**Option B — Security is a code-review-only checklist**:
- Remove the "any REQ with domain: security" claim from `crew-members/security.md:3` JIT_CONTEXT.
- Add a note to CLAUDE.md that `security.md` is a reference checklist loaded only by `code-review.md`, not a crew-member in the standard sense.
- Consider renaming to `references/security-owasp.md` or similar if it doesn't conceptually belong in `crew-members/`.

## Acceptance

- CLAUDE.md, `crew-members/security.md`, and `actions/work.md` agree on how/when `security.md` loads.
- A reader can find the answer to "when does security.md load?" in one place without contradiction.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Findings: `architecture.md` F2, `architecture.md` F8

---

## Triage

**Route: B** — Medium

**Reasoning:** Three small, coordinated edits across three files; the architectural decision (Path A vs B) was the only real choice and the user answered Path A. No planning needed.

**Planning:** Not required.

## Decisions

- **D-01:** Chose **Path A** (promote security to a real crew) over Path B (demote to code-review-only checklist), per the user's AskUserQuestion answer. Consistent with the just-shipped `crew-members/prompt-injection.md` as a sibling crew rule.

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — Schema Read Contract enum (L194) to add `security` + `testing`; Step 6 to add substep 4a; example frontmatter comment (L153) to match.
- `CLAUDE.md` (modify) — Agent Rules section to add `security.md` bullet next to `testing.md`.
- `crew-members/security.md` (modify) — JIT_CONTEXT rewritten to point at the new canonical source (work.md Step 6 substep 4a + CLAUDE.md Agent Rules).

**Files I will NOT touch:** `actions/code-review.md` (already conditionally loads security.md at L159 — no change needed; the JIT_CONTEXT now references it correctly).

**Acceptance criteria (restated from REQ):**
- [x] CLAUDE.md, `crew-members/security.md`, and `actions/work.md` agree on how/when security.md loads.
- [x] A reader can find the answer in one place (CLAUDE.md Agent Rules is now the canonical reference, work.md Step 6 has the exact load step, security.md JIT_CONTEXT points at both).

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified) — Schema Read Contract `domain` enum row now lists `frontend, backend, ui-design, general, security, testing` with two new normalization aliases (`sec` → `security`, `test` → `testing`). Example frontmatter comment at L153 updated to match. Step 6 gains a new substep 4a that loads `crew-members/security.md` on `domain: security` OR on a description-keyword match for auth/crypto/session/secrets/input-validation/OWASP — heuristic, lean-toward-loading.
- `CLAUDE.md` (modified) — new Agent Rules bullet for `security.md` immediately after `testing.md`, with the same OR-clause loading rule.
- `crew-members/security.md` (modified) — JIT_CONTEXT rewritten to point at the canonical loading reference (CLAUDE.md Agent Rules + work.md Step 6 substep 4a) and at the existing `code-review.md` conditional load.

**What was done:** Closed the loading-model orphan status. Three files now agree on when `security.md` loads: `domain: security` OR keyword-match on auth/crypto/etc. The `testing` domain is also formally added to the enum since the existing testing.md JIT trigger already used it. `code-review.md`'s conditional load at L159 unchanged (the JIT_CONTEXT now references it correctly).

## Qualification

Passed — 3 files modified per scope. The enum addition + Step 6 substep + CLAUDE.md bullet form a consistent contract. The security.md JIT_CONTEXT now points at the contract rather than re-describing it (single source of truth in CLAUDE.md / work.md). No drift to `actions/code-review.md` (its existing conditional load is already correct).

## Testing

**Tests run:** Manual cross-reference audit:
- `grep "domain.*security" actions/work.md CLAUDE.md crew-members/security.md` — all three files name the trigger consistently.
- `grep "Step 6.*security\|substep 4a" actions/work.md crew-members/security.md` — the load location is referenced from both.
- Existing `actions/code-review.md:159` reference still resolves to the file.

**Result:** ✓ All three files agree. ✓ No contradiction between JIT_CONTEXT and the canonical loading rule.

*Verified by work action*

## Review

**Overall: 95%** | 2026-05-29T19:06Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor
**Acceptance:** Pass — loading model is consistent across the three files; security.md is now a first-class crew rule.
**Follow-ups created:** None

*Reviewed by work action (Route B self-review)*

## Lessons Learned

**What worked:** Pinning the canonical loading reference to CLAUDE.md Agent Rules (and letting the crew file's JIT_CONTEXT just *point* at it) eliminated the divergence risk for next time. The previous JIT_CONTEXT was a near-duplicate of CLAUDE.md's missing entry — putting the answer in one place stops the two from drifting.
**Worth knowing:** `testing` was de facto a domain (testing.md loads on `domain: testing`) but wasn't formalized in the enum. Adding it now closes a smaller version of the same architectural inconsistency.

