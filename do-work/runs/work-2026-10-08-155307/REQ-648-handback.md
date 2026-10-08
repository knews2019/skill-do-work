# REQ-648 hand-back (the assigned_to schema line and the work action say operator-gated work is blocked, never earmarked)

- Branch: `worktree-agent-REQ-648-schema-and-work-action-state-operator-rule`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-648-schema-and-work-action-state-operator-rule`
- Base: `dd54712a`
- Commits: `4fc5967e` (single commit, `[REQ-648] assigned_to schema line and work action say operator-gated work is blocked, never earmarked`)

## File manifest

- `skills/do-work/actions/work-reference.md` (modified): the `assigned_to:` schema comment gains the sentence "**For another session or checkout only** — work that waits on the user as operator is `status: blocked` with `blocked_by` naming them (`actions/capture.md` Step 1's Earmark assessment), never an invented earmark." All existing sentences, including the parser lock-step sentence, are unchanged.
- `skills/do-work/actions/work.md` (modified): the Error Handling table gains the "Orchestrator wants to park a queued REQ on the operator" row, worded as the REQ quotes it, plus the pointer. The existing "Builder reports a missing external precondition" row now points at `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row, instead of the dangling "Step 8 → **Mid-run blocked flip**".

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the REQ, UR-142, crew members (general, coding-guardrails, shared-principles, communication-style), prime-action-files (Cross-Referencing spelling: same-package `actions/<file>.md`, sections as `→ **Heading**`), prime-releases, and every `alternate-writer-contract-drift` lesson. Approach: one sentence inserted into the schema comment right after the "who writes it" sentence; one new table row directly after the missing-precondition row; replace the dangling pointer in the existing row with the real section pointer and reuse it in the new row. Confirmed the cited target exists: `actions/capture.md` § Step 1: Parse and Assess holds the **Earmark assessment** bullet (line 107), and `work-reference.md` holds `## Failure Classification (Step 8)` with the Environment row carrying the blocked-flip test. No test pins the old pointer text.
- [x] **[APPLY]:** Edited exactly the two write-set files with an anchored, count-asserted replacement. No other file touched.
- [x] **[UNIFY]:** `git diff --stat`: work-reference.md 2 +-, work.md 3 ++-, 2 files, 3 insertions, 2 deletions. `git diff --check` exit 0. `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 (PASS). `bash _dev/tests/contract-regressions.sh` exit 0 ("Contract regression checks passed."). Files checked: work-reference.md (line 114 keeps every prior sentence, bold labels unchanged, the new sentence adds no case for writing `assigned_to`); work.md (table keeps two columns, every existing row unchanged except the pointer text, the new row sits beside the missing-precondition row). No debug artifacts.

## Red-green evidence

- RED (before the edit): `grep -c operator` returned 0 for both files. `work.md:522` pointed at "Step 8 → **Mid-run blocked flip**", a heading that does not exist. An orchestrator reading only these two places had no rule against writing `assigned_to` for operator-gated work.
- GREEN (after the edit): `grep -n operator` hits `work-reference.md:114` (the schema sentence) and `work.md:523` (the new table row). `grep -n "Mid-run blocked flip" skills/do-work/actions/work.md` prints nothing (grep exit 1). Both blocked-flip rows point at `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row, which exists.

## Decisions

- D-01: The schema sentence sits right after "Seeded by capture … or written by a session claiming from another checkout." and before "**Not a lock and not a status:**". It reads as a direct qualifier of who may write the field, and it leaves the closing board-placement and parser lock-step sentences intact as the line's tail.
- D-02: The pointer wording is `` `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row ``, matching the prime's same-package spelling and the existing `work.md` habit of `` `actions/work-reference.md` → **Heading** ``. The brief's parenthetical "(the blocked-flip test)" was left out of the inline pointer because the existing row already opens with "Apply the blocked-flip test", and the new row says "the same mid-run flip as a missing precondition".
- D-03: In the new row the pointer goes in parentheses after "the same mid-run flip as a missing precondition", so the REQ's quoted sentences stay verbatim around it.

## Discovered Tasks

- impact-low: `skills/do-work/actions/work-reference.md:788` (Failure Classification, Environment row) says "see `actions/work.md`'s mid-run blocked-flip procedure", and `:165` says "Step 8's blocked-flip procedure". `work.md` has no such procedure; after this REQ, `work.md` points back at this very row, so the two citations now form a loop with no procedure text beyond the row itself. The brief and REQ kept these prose mentions on purpose. A later REQ could drop the "see work.md" clause at :788. → report only

## Lessons read

- Satellite: `_dev/primes/lessons-action-files.md`, every bullet carrying `[family: alternate-writer-contract-drift]` (REQ-477, REQ-498, REQ-513, REQ-461, REQ-531, REQ-566, REQ-640, REQ-641, REQ-642, REQ-644). Applied: grepped the restated phrase ("blocked flip", "blocked-flip") across the action files rather than only the cited line, which surfaced the :788 loop above. The capture-side restatements (`capture.md:107/108/143`, `capture-reference.md:41`, `docs/work-guide.md:120`) were left untouched per requirement 4.

## Proposed lesson bullet (integrator writes it and refreshes the lessons-index row)

```
- [family: alternate-writer-contract-drift] [REQ-648: a rule stated only in the action that normally writes a field does not reach a session that writes it from elsewhere; put the "not for this" sentence on the field's own schema line, which is what an off-path writer reads first, and check that every pointer you add or keep lands on a heading that exists (the old "Step 8 → Mid-run blocked flip" pointed at nothing)](../../do-work/archive/UR-142/REQ-648-schema-and-work-action-state-operator-rule.md#lessons-learned)
```

## Proposed CHANGELOG entry (integrator writes it and bumps the version)

Title: **Operator-Gated Work Is Blocked, Not Earmarked, in the Run Action**

A run no longer has a gap where it could park work on you with a session earmark. In a consumer run, an orchestrator hand-wrote `assigned_to: 'user-interactive'` on two REQs that needed production access, so the default run skipped them and the board showed them as Earmarked instead of waiting on you.

- The `assigned_to` schema line in `actions/work-reference.md` now says the field is for another session or checkout only. Work that waits on you as operator is `status: blocked` with `blocked_by` naming you.
- The Error Handling table in `actions/work.md` has a new row for an orchestrator that wants to park a queued REQ on the operator: flip it to blocked, never write `assigned_to`.
- The missing-precondition row pointed at a heading that does not exist. Both rows now point at `actions/work-reference.md` → Failure Classification (Step 8), Environment row.

## Integration seams

None. The two files do not overlap REQ-647 (clarify) or REQ-649 (board tooltip) write sets. No parser, field, status, or routing change, so `model.go` stays in lock-step without edits.

## Wall times

| Check | Exit | Wall |
|---|---|---|
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 1s |
| `bash _dev/tests/contract-regressions.sh` | 0 | 27s |
| acceptance greps, `git diff --check` | as expected / 0 | under 1s |
