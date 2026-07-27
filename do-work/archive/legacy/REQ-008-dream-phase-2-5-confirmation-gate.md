---
id: REQ-008
title: "Code review: add Phase 2.5 confirmation gate to dream.md before destructive Phase 3"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:52:00Z
completed_at: 2026-05-29T18:54:00Z
commit: 624e825
route: B
review_generated: true
source: code-review
scope: actions/dream.md
domain: backend
---

# Code Review Fix: Add Phase 2.5 Confirmation Gate to dream.md

## What

`actions/dream.md` is destructive by design — Phase 3 deletes wiki pages, edits files in place, repoints inbound links, "prunes the untrue", and rewrites transcripts (L121-145, 162-198). The action philosophy says "destructive by design, so it never runs automatically" (line 3) and "Manual only / refuse if asked" (lines 21, 202).

**But the only consent point is the outer `do-work dream` invocation token** — the same single-token consent that would trigger any read-only action. The action jumps straight from Phase 1 ("find the .lock") to Phase 3 ("Prune the untrue") with no preview, no per-cluster confirmation, no dry-run mode.

A user who casually invokes `do-work dream` expecting "a consolidation pass" can lose pages they wanted to keep — especially on the first run when they don't yet know the worklist.

Compare the sibling destructive action `actions/stray-check.md` (lines 27-30, 74-80): requires both an explicit `fix` mode token **and** an in-step "Apply these fixes? [all / numbers / none]" ask-user prompt. That's the precedent to clone.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). This is the highest-impact security finding in the review.

## Requirements

- Add a **Phase 2.5: Preview & Confirm** after the worklist is built (Phase 2) and before Phase 3 begins:
  - Print the worklist summary: count and identity of merges, prunes, deletions, link-repointings.
  - Require an ask-user confirmation: "Apply these N fixes? [all / dry-run / specific clusters / none]"
  - Default to dry-run on ambiguous response.
- Add a `--dry-run` mode token (mirroring `stray-check`'s `report` mode) that lets users preview without writing.
- Add a Rule to `actions/dream.md`: "Phases 1–2 must produce a visible worklist before Phase 3 may begin. The scan phase makes zero writes."
- Update the Verification Checklist to lock the new invariant: "Phase 2.5 confirmation was presented before any Phase 3 writes" and "No Phase 3 writes occurred if user declined or chose dry-run".

## Acceptance

- `actions/dream.md` has a Phase 2.5 with explicit confirmation prose.
- `--dry-run` mode token is documented in the Modes/Input section.
- Verification Checklist includes the two new items above.
- A first-time invocation of `do-work dream` produces a visible worklist and waits for user confirmation before mutating anything.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Finding: `security.md` Finding 1 (Important)

---

## Triage

**Route: B** — Medium

**Reasoning:** Single-file edit but the gate has to fit cleanly between two existing phases without breaking the prose flow. Light exploration of dream.md's structure (Input, Steps 1-5, Rules, Common Rationalizations, Red Flags, Verification Checklist) before editing. Mirroring `stray-check`'s `report`-mode precedent (cited in the REQ) clarifies the wording.

**Planning:** Not required.

## Exploration

Read `actions/dream.md` in full. Confirmed:

- The four phases map to Steps 2 (Phase 1), 3 (Phase 2), 4 (Phase 3), 5 (Phase 4). The new gate fits as Step 3.5 = Phase 2.5.
- Phase 4 (Step 5) is the only step that releases the lock — a dry-run must still hit Step 5's lock-release path, but skip its writes.
- The `## Input` section was a single sentence; `--dry-run` needs to be added as a mode token to the documented args.
- `## Rules` already has a "Phases 1–2 are read-only" rule (L203 pre-edit). The new "must produce a visible worklist before Phase 3" rule layers on top.
- `## Verification Checklist` has a "worklist collected" item but no "preview was presented" item or "no writes when declined" item — both need to be added.
- Sibling action `stray-check.md` was cited in the REQ as precedent for `report`-mode + per-fix prompt. The same shape — `--dry-run` mode token + an in-step `Apply these N fixes? [all / dry-run / specific clusters / none]` prompt — fits here.

## Scope

**Files I will touch:**
- `actions/dream.md` (modify) — six insertion sites: Input section (mode token), Step 1 (peel `--dry-run` token), Step 3.5 (new section), Step 5 (dry-run short-circuit at the top), Rules (new "visible worklist" rule), Common Rationalizations (two new rows), Red Flags (two new rows), Verification Checklist (two new items).

**Files I will NOT touch:** `actions/stray-check.md` (cited as precedent, no change needed), SKILL.md (the existing `dream` routing handles `--dry-run` as `$ARGUMENTS` content — no router change required).

**Acceptance criteria (restated from REQ):**
- [x] `actions/dream.md` has a Phase 2.5 with explicit confirmation prose.
- [x] `--dry-run` mode token documented in the Input section.
- [x] Verification Checklist includes the two new items.
- [x] First-time invocation produces a visible worklist and waits for confirmation before any mutation.

## Implementation Summary

**Files changed:**
- `actions/dream.md` (modified) — added `--dry-run` mode token to `## Input`; added peel-out logic as Step 1 substep 1 (renumbered subsequent substeps); inserted new `### Step 3.5: Phase 2.5 — Preview & Confirm` section between Phase 2 (Step 3) and Phase 3 (Step 4) with the `Apply these N fixes? [all / dry-run / specific clusters / none]` prompt and the default-to-dry-run rule for ambiguous responses; added dry-run/declined short-circuit at the top of Step 5; added a new Rule (`Phases 1–2 must produce a visible worklist before Phase 3 may begin`); added two `## Common Rationalizations` rows guarding against gate-skipping; added two `## Red Flags`; added two `## Verification Checklist` items.

**What was done:** Closed the consent gap between the single-bit `do-work dream` invocation token and the destructive Phase 3 writes. The new Phase 2.5 step requires an explicit preview-then-confirm round before any mutation. `--dry-run` short-circuits the prompt and exits cleanly after Phase 4's lock release. The defaults bias toward safety: ambiguous responses become `dry-run`, never `all`.

## Qualification

Passed — 1 file modified per scope. The Phase 2.5 section sits between Phase 2 and Phase 3 in the prose order. The `--dry-run` token is documented in Input, peeled in Step 1, short-circuits Step 5 writes, and is honored by Step 6's summary tagging. The Rules / Common Rationalizations / Red Flags / Verification Checklist additions enforce the invariant at multiple layers. No debug artifacts, no scope drift.

## Testing

**Tests run:** Manual walkthrough of the new Phase 2.5 flow against each documented response (`all`, `dry-run`, `specific clusters`, `none`, ambiguous), plus the `--dry-run` mode-token shortcut.
**Result:** ✓ Every path either routes through Phase 3 with the chosen worklist subset or skips Phase 3 entirely, with Phase 4 always releasing the lock. ✓ The "ambiguous → dry-run" default closes the worst footgun. ✓ Step 1's peel-out keeps `--dry-run` from being treated as a path.

*Verified by work action*

## Review

**Overall: 92%** | 2026-05-29T18:54Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — gate prose, mode token, rule, and checklist items all landed.
**Suggested testing:** A future REQ could add an integration test (or at minimum a fixture-driven exercise) that proves the dry-run path makes zero writes on a sample memory dir. The current verification is documentation-level only.
**Follow-ups created:** None — suggested-testing is quality-of-life, not Important.

*Reviewed by work action (Route B self-review)*

## Lessons Learned

**What worked:** Cloning `stray-check`'s `report`-mode shape directly (the REQ named it as precedent) made the wording fall out — `all / dry-run / specific clusters / none` is a familiar pattern, not a new vocabulary. Adding the guardrails at four layers (Rule + Common Rationalizations + Red Flag + Verification Checklist) was deliberate redundancy: each layer catches a different kind of agent shortcut.
**What didn't:** Initial edit duplicated the substep number `2.` in Step 1 (the original had a list 1-2-3-4, the new "Peel out the mode token first" became a new 1, but the renumbering didn't cascade). Caught on re-read, fixed before commit. Lesson: when prepending an item to a numbered list, re-read the whole list.
**Worth knowing:** The dry-run path still releases the `.lock` (via Step 5 substep 5), but skips writes 1-4 and skips the `log.md` `[dream]` line. A dry-run is not a pass — that's why the log doesn't gain an entry.

