---
id: REQ-026
title: next-steps.md — intent over enumeration
status: completed
created_at: 2026-07-27T07:34:50Z
claimed_at: 2026-07-27T07:40:27Z
completed_at: 2026-07-27T08:05:19Z
commit: b29ec34
user_request: UR-006
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-025, REQ-027, REQ-028, REQ-029, REQ-030, REQ-031]
batch: context-engineering-alignment
---

# next-steps.md — Intent Over Enumeration

## What

`next-steps.md` is 1,741 words of 40 hard-coded worked examples — one fenced "Next steps:" block per action. `SKILL.md` points at it after *every* action, so it is part of the always-read floor. Replace the enumeration with:

1. The **stated intent** — after an action completes, suggest 2–3 next actions inferred from the current queue/pipeline state, using fully-qualified `do-work <verb>` names the user can copy-paste.
2. The handful of **formatting constraints** that a model would not infer (block shape, command-name form, how many to show).
3. A short **table of the genuinely non-obvious cases** — the ones where the right suggestion is *not* inferable from the action that just ran: pipeline interrupted vs. completed, `reserve`, `capture`, `clarify` with pending answers.

Target: ≤ 500 words, down from 1,741.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/general.md`, `crew-members/coding-guardrails.md`, `crew-members/maintenance.md` (maintenance:true — delete-before-you-add applies). Read all 47 "After X" blocks in next-steps.md. Approach: classify each as obvious-from-context (delete, rely on the new intent statement + SKILL.md dispatch table for judgment) vs encodes-a-judgment-call (keep as a table row). Rewrite the file as an intent statement, format constraints, one anchor example block, and a non-obvious-cases table. Scope limited to `next-steps.md` only.
- [x] **[APPLY]:** Rewrote `next-steps.md` in place: 1,741 → 383 words. Kept 6 table rows (pipeline×2, reserve×1, capture-requests×1, clarify×2) out of 47 original blocks; deleted 41. SKILL.md untouched (its pointer at line 138 already just says "See `next-steps.md`" — no edit needed there).
- [x] **[UNIFY]:** `git diff --stat` shows only `next-steps.md` changed. `bash _dev/tests/contract-regressions.sh` passes clean (including the `ultracode|fable` absence check and the shipped-citation check, both of which scan `next-steps.md`). No debug artifacts. `wc -w next-steps.md` = 383 (≤500 target).

## Why (if provided)

Anthropic's context-engineering guidance for Claude 5 generation models replaces few-shot enumeration with interface design: a capable model given the intent plus the constraints produces better, more situation-appropriate output than one pattern-matching against 40 canned examples — and the examples cost real context on every single invocation. Most of the 40 blocks say the obvious thing (after `capture`, suggest `run`); only a few encode judgment the model couldn't reach on its own. Those few are exactly what the table preserves.

## Context

- `SKILL.md` line ~138: "After every action completes, suggest the next logical prompts the user might want to run. See `next-steps.md` for the full per-action reference..." That pointer stays; its target gets much smaller.
- The always-read floor is `SKILL.md` + `next-steps.md`. This REQ is the single largest available reduction to that floor without touching the router.
- Maintenance pass on the skill's own instructions — `crew-members/maintenance.md` (delete-before-you-add) loads via the `maintenance: true` marker.

## Detailed Requirements

- **Preserve the non-obvious cases.** Before deleting, read all 40 blocks and classify each: obvious-from-context (delete) vs. encodes-a-judgment-call (keep, as a table row). Record the classification count in the Implementation Summary. Known keepers from the audit: pipeline *interrupted* vs. *completed* (different suggestions), `reserve`, `capture`, `clarify` with pending answers. Do not treat that list as exhaustive — it is illustrative; the classification pass decides.
- **Keep the output shape.** Suggestions are rendered as a fenced block with aligned `do-work <verb>` on the left and a short gloss on the right. That format is a user-facing convention; state it as a constraint, don't drop it. One example block is acceptable to anchor the format — one, not forty.
- **Fully-qualified names.** Suggestions must name real, routable verbs. State this as a constraint plus the location of the authoritative verb list (`SKILL.md` routing table), rather than re-listing verbs — a hand-copied verb list is exactly the closed enumeration that goes stale.
- **Word budget:** ≤ 500 words. Report the final `wc -w`.
- **Do not grow `SKILL.md`** to compensate. If something must be said at routing time, this REQ is the wrong home for it.

## Constraints

- `bash _dev/tests/contract-regressions.sh` must pass clean — note it asserts `next-steps.md` contains no `ultracode|fable` mentions.
- Version bump + descriptive `CHANGELOG.md` entry.

## Dependencies

None. Independent of the other six REQs in the batch.

## Builder Guidance

**Certainty: Firm on the target, exploratory on the table's contents.** The classification pass is the real work — be honest about which cases are genuinely non-obvious rather than preserving blocks out of caution. A row that says what any competent reader would already suggest is a row to cut.

## Red-Green Proof

- **RED now:** `wc -w next-steps.md` ≈ 1,741, and the file contains ~40 fenced example blocks (`grep -c '^```' next-steps.md` returns a large count). The intent it encodes is never stated — it must be inferred from the examples.
- **GREEN when:** `wc -w next-steps.md` ≤ 500; the file opens with an explicit statement of intent; the non-obvious cases survive as table rows; at most one example block remains as a format anchor.
- **Validation:** word-count receipt (file, always-read floor, orchestrator load); the classification tally (blocks read / kept / deleted); `bash _dev/tests/contract-regressions.sh` clean. Spot-check by naming three actions and confirming the rewritten file gives an agent enough to produce correct suggestions for each.

## Open Questions

None.

## Full Context

See `do-work/user-requests/UR-006/input.md` for complete verbatim input.

---
*Source: "compare with the current skill, is there something that we need to update?" — resolved into the approved seven-REQ plan.*

Think carefully before answering.

## Triage

**Route: B** — outcome and target word count are fixed by the REQ (≤500 words, intent + constraints + one example + non-obvious table), the only real judgment is the classification pass over the 47 existing blocks. No design exploration needed; this is read-classify-rewrite against a single file.

## Scope

**Files I will touch:** `next-steps.md` only (per builder file scope — this REQ's declared scope matches).

**Files I will NOT touch:** `SKILL.md` (its pointer at line ~138 already just says "See `next-steps.md` for the full per-action reference" — no wording depends on the old enumeration's length or contents, so no edit needed there); `actions/version.md` and `CHANGELOG.md` (orchestrator-owned, explicitly off-limits to builders).

**Acceptance criteria (restated from REQ):**
- [x] File opens with an explicit statement of intent (infer 2-3 next actions from outcome + queue state, suggest as copy-pasteable `do-work <verb>` commands).
- [x] Format constraints stated (fenced `Next steps:` block, aligned columns, fully-qualified verbs pointing at `SKILL.md`'s Action Dispatch table as the authoritative list, 2-3 suggestions cap, closing `do-work help` reminder).
- [x] Non-obvious cases survive as table rows, not full worked blocks.
- [x] At most one example block remains, as a format anchor (exactly one, explicitly labeled "not a template to reuse verbatim").
- [x] `wc -w next-steps.md` ≤ 500 → actual 383.
- [x] `bash _dev/tests/contract-regressions.sh` passes clean.
- [x] SKILL.md not grown to compensate (untouched; word count unchanged at 2,557).

## Implementation Summary

**What was done:** Read all 47 "After X" blocks in the original `next-steps.md` (some action names covered more than one block, e.g. `pipeline` completed/interrupted, `clarify` pending/answered, `tidy-repo` plan-only/executed, several `bkb` subcommands, four `interview` states — 47 blocks across ~40 distinct actions per the REQ's estimate). Classified each:

- **Obvious-from-context (deleted, 41 blocks):** every block whose suggestion is the natural next step in the action's own pipeline and doesn't depend on a queue/domain condition a competent model wouldn't already infer from having just run the action — `work`, `verify-requests`, `review-work`, `validate-feedback`, `code-review`, `ui-review`, `ai-report`, `present-work`, `slop-check`, `forensics`, `roadmap`, `board`, `note`, `cleanup`, `stray-check`, `tidy-repo` (both states), `install`, `prime create`, `prime audit`, `quick-wins`, `scan-ideas`, `deep-explore`, `inspect`, `commit`, `abandon`, all six `bkb` blocks, `dream`, all three `prompts` blocks, all four `interview` blocks, `tutorial`, `version`/`recap`.
- **Encodes-a-judgment-call (kept, 6 rows covering 4 actions):** `pipeline` (interrupted vs. completed — genuinely different suggestions, and getting it backwards actively misleads), `reserve` (the non-obvious "sync via commit before another session can see it" reasoning), `capture-requests` (verify-requests-before-run ordering), `clarify` (pending vs. all-answered — suggesting `run` while questions are still open silently skips those REQs). These match all four of the REQ's named "known keepers" exactly; the classification pass did not surface additional non-obvious cases beyond them (see rationale below for `tidy-repo`, the one candidate I considered and rejected).
- Considered but rejected as a keeper: `tidy-repo` plan-only vs. executed. Tempting because a wrong suggestion there (commit before executing) would mislead the same way the kept rows do. Rejected because the state distinction is already carried in the action's own name/invocation ("plan-only" vs. an executed run) — an agent that just ran the action already knows which one happened, unlike `pipeline`/`clarify`/`reserve` where the *result* doesn't announce the state in its label.

Rewrote the file with: an opening **Intent** paragraph (infer suggestions from actual outcome + queue state, not a fixed action→suggestion table, with three concrete contrast examples baked into the sentence itself — empty queue, clean tree, `ui-design` domain); a **Format Constraints** list (fenced-block shape, alignment, fully-qualified-verbs-from-SKILL.md-Action-Dispatch-table constraint instead of a hand-copied verb list, 2-3 cap, closing `do-work help` line); one **Example** block explicitly labeled as a format anchor, not a template; and a **Non-Obvious Cases** table with the 6 surviving rows plus a lead-in sentence explaining what qualifies a row for the table (state-dependent + misleading-if-wrong) and what doesn't (ordinary pipeline-order inference).

**Files changed:**
- `next-steps.md` (modified) — 1,741 → 383 words. Structure changed from 47 fenced worked-example blocks (one per action/state) to: Intent statement, Format Constraints list, one anchor Example block, Non-Obvious Cases table (6 rows / 4 actions).

**Word-count receipt:**

| File | Before | After |
| --- | --- | --- |
| `next-steps.md` | 1,741 | 383 |
| `SKILL.md` | 2,557 | 2,557 (untouched) |
| Always-read floor (`SKILL.md` + `next-steps.md`) | 4,298 | 2,940 |
| Orchestrator load (`SKILL.md` + `actions/work.md` + `actions/work-reference.md` + `crew-members/general.md` + `crew-members/coding-guardrails.md` + `next-steps.md`) | 25,231 | 23,873 |

(`actions/work.md` 11,983w, `actions/work-reference.md` 7,426w, `crew-members/general.md` 678w, `crew-members/coding-guardrails.md` 846w — all unchanged, outside this REQ's scope.)

## Testing

- `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` (exercises the `ultracode|fable` absence check on `next-steps.md`, and the shipped-citation check that scans `next-steps.md` for dangling `CLAUDE.md`/`AGENTS.md` references — the rewritten file has neither).
- `wc -w next-steps.md` → 383 (target ≤ 500).
- `grep -c '^```' next-steps.md` → 2 (one open + one close fence for the single anchor example; RED baseline was 94, i.e. 47 blocks × 2 fences).
- `git diff --stat` → only `next-steps.md` changed (spot-checked; no edits leaked into `SKILL.md` or any other file).
- Spot-check (Red-Green Proof's "name three actions" validation), reasoning through what the rewritten file gives an agent with no other context beyond `SKILL.md`'s Action Dispatch table:
  - **`code-review`** (not in the table): the Intent line says to infer from outcome + queue state. A code review that found issues has follow-up REQs in the queue → `do-work run` to process them is inferable; `do-work capture-request:` for a finding not yet captured is inferable from the general "capture new requests" pattern any action can trigger. Correct.
  - **`board`** (not in the table): board just showed the queue as a Kanban view. `do-work run` to pick up a Ready card, `do-work roadmap` for the text-survey equivalent — both inferable from "what would a user naturally do after looking at their queue."
  - **`clarify`** (in the table, pending-answers row): the file explicitly tells the agent not to suggest `do-work run` while questions remain open, and to suggest `do-work clarify` again instead — this is exactly the case that would NOT be correctly inferred without the table row (a naive model would default to "clarify happened, now run the queue"), which is why it survived the classification pass.

## Lessons Learned

**What worked:** Reading all 47 blocks against one question — "would a competent agent, given a stated intent and format constraints but no example for THIS action, produce a suggestion the user would recognize as reasonable?" — separated the obvious from the non-obvious faster than trying to rate each block's individual merit. Almost every block restated "the natural next pipeline step," which is exactly what an intent statement replaces in one sentence instead of forty repetitions of it.
**What didn't:** Early draft kept `tidy-repo` plan-only/executed as a 7th/8th row by analogy to `pipeline`'s two-state pattern, before checking whether the distinguishing state was actually hidden from the agent the way `pipeline`'s is. It isn't — the action's own recent output announces "plan-only" or "executed" — so it was cut on a second pass. Worth re-checking each candidate row against "does the agent already know which state it's in from what just happened," not just "does this action have more than one possible follow-up."
**Worth knowing:** The four kept actions (`pipeline`, `reserve`, `capture-requests`, `clarify`) all share a common shape: the action's own recent transcript output does NOT self-announce which of two-plus states it left the queue in, so an agent working from intent alone could plausibly infer the wrong branch. That's a sharper test for "keep this row" than the REQ's own words ("not inferable from the action that just ran") might suggest in isolation — worth stating explicitly if this table ever needs to grow again, so a future editor doesn't over-add by loose analogy to an existing row's action name.
