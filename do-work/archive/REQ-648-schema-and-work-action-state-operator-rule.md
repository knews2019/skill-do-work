---
id: REQ-648
title: '[impact-rule-change] The assigned_to schema line and the work action say operator-gated work is blocked, never earmarked'
status: completed
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-08T15:54:08Z
route: A
created_at: 2026-10-08T15:50:17Z
user_request: UR-142
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-647, REQ-649]
batch: earmark-release-visibility
write_set: [skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md]
claimed_at: 2026-10-08T15:53:15Z
dispatch_at: 2026-10-08T15:56:02Z
builder_handback_at: 2026-10-08T16:00:34Z
integration_at: 2026-10-08T16:00:51Z
review_at: 2026-10-08T16:05:04Z
kb_status: pending
commit: 4c50ccdecf64fc63ec429ca252768252e32faf83
heavy_verified_at: 2026-10-08T16:06:42Z
heavy_verified_revision: 4c50ccdecf64fc63ec429ca252768252e32faf83
completed_at: 2026-10-08T16:07:19Z
release_at: 2026-10-08T16:07:19Z
---
# The assigned_to Schema Line and the Work Action Say Operator-Gated Work Is Blocked, Never Earmarked
## What
The rule "a session gets an earmark, the operator gets a block" lives only in `actions/capture.md` (Earmark assessment) and `docs/work-guide.md`. The `assigned_to` schema line in `actions/work-reference.md` and the run orchestrator's error-handling table in `actions/work.md` do not carry it, so a session that writes the field outside capture has no rule to follow. Add one sentence to each.
## Why
Consumer incident, Oct 7 22:58 UTC on 0.305.71: a run orchestrator, mid-run, hand-wrote `assigned_to: 'user-interactive'` on two pending REQs that needed the operator's production access. No action told it to. The field's own schema line is what a session reads first when it decides whether to write the field, and it says who seeds the field and that it is advisory, never what it is not for. `work.md`'s table covers a builder reporting a missing precondition, not an orchestrator that wants to park a queued REQ on the operator. Triage (validate-feedback, 2026-10-08, F4) accepted both sentences.
## Verified Facts (from triage)
- `grep -n operator skills/do-work/actions/work-reference.md skills/do-work/actions/work.md` finds nothing today.
- `skills/do-work/actions/work-reference.md:114` is the `assigned_to` schema line: "OPTIONAL advisory claim marker: the session this REQ is earmarked for … Seeded by capture when the user earmarks work, or written by a session claiming from another checkout. Not a lock and not a status …". It ends with the board placement sentence and the parser lock-step rule.
- `skills/do-work/actions/work.md:522` (Error Handling table): "Builder reports a missing external precondition | Apply the blocked-flip test (Step 8 → **Mid-run blocked flip**) …". That bold heading does not exist in `work.md`; the procedure is the Environment row of **Failure Classification (Step 8)** in `actions/work-reference.md` (line ~788, "First apply the blocked-flip test"), and `work.md` Step 8 substep 3 names it as "external-precondition blocked-flip judgments". `grep -rn "Mid-run blocked flip" _dev/tests/` finds nothing, so repointing is safe.
- The rule already stands at `actions/capture.md:107` (Earmark assessment), `actions/capture.md:108` (External-condition look-alikes), and `docs/work-guide.md:120`. `actions/capture-reference.md:41` and `actions/capture.md:143` restate earmark mechanics, not the rule. REQ-644 (0.305.78) wrote the capture and guide sentences and its restatement sweep found `work-reference.md:114` consistent then, because the rule was additive.
- The `assigned_to` line's parser lock-step rule (`../../do-work-board/tools/queue-kanban/model.go`) concerns the field's read semantics and board placement. Adding a sentence about who may write the field changes neither, so no parser change is needed.
## Detailed Requirements
1. `actions/work-reference.md:114`, append to the `assigned_to` schema comment, before or after the board-placement sentence as reads best: "**For another session or checkout only** — work that waits on the user as operator is `status: blocked` with `blocked_by` naming them (`actions/capture.md` Step 1's Earmark assessment), never an invented earmark." Keep every existing sentence of the line.
2. `actions/work.md` Error Handling table, next to the row at line ~522, add: `| Orchestrator wants to park a queued REQ on the operator (production access, credentials, a person at the keyboard) | Flip it to `status: blocked` with `blocked_by` naming the operator and `blocked_at`, the same mid-run flip as a missing precondition. Never write `assigned_to` for this: the earmark is for another session, it is invisible as "waiting on you", and the default run hides it from the user by skipping it. |`
3. Pointer fix: both rows (the existing "missing external precondition" row and the new one) point at the procedure that exists: `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row (the blocked-flip test). Replace the dangling "Step 8 → **Mid-run blocked flip**" with that pointer. Check `work-reference.md:165` and `:788`, which call it "the mid-run blocked flip" and "the mid-run blocked-flip procedure" in prose; those descriptions are fine as prose and stay.
4. Leave `actions/capture.md:107`, `capture.md:143`, `actions/capture-reference.md:41`, and `docs/work-guide.md:120` unchanged. Do not touch `model.go`.
5. Acceptance grep: `grep -n "operator" skills/do-work/actions/work-reference.md skills/do-work/actions/work.md` finds the schema sentence and the table row.
6. Run `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`; both green.
7. Release per `_dev/primes/prime-releases.md` (the integrator writes the CHANGELOG entry and bump).
## Constraints
- Prose only. No field, status, Go, routing, or parser change. `assigned_to` keeps its verbatim-read class and its advisory meaning; the default scan and `advance` select exactly as before.
- Keep the never-invent rule intact in spirit: the new sentences say when NOT to write `assigned_to`; they do not add a case for writing it.
- Read `_dev/primes/prime-action-files.md` before editing; cross-references use that prime's spelling.
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-647 or REQ-649 into it.
## Dependencies
None. REQ-647 (clarify) and REQ-649 (board tooltip) touch different files; the three may run in parallel.
## Builder Guidance
Certainty is high: both sentences are quoted in the feedback and accepted as worded. Latitude: where in the long schema line the new sentence sits, and the exact pointer wording for requirement 3.
## Red-Green Proof
**RED prompt/case:** A run orchestrator, mid-run, reads only `actions/work.md`'s Error Handling table and the `assigned_to` schema line in `actions/work-reference.md`, and wants to park a queued REQ that needs the operator's production access.
**Why RED now:** Neither place says the earmark is for another session only, so the orchestrator writes `assigned_to: 'user-interactive'`, the default run skips the REQ, and the board shows it under Pending → Earmarked instead of Needs input · Blocked.
**GREEN when:** The schema line says the field is for another session or checkout only and operator-wait work is `blocked` with `blocked_by`; the table row routes the park case to the mid-run blocked flip and forbids `assigned_to` for it; `grep -n operator` on both files finds both; the two rows point at a heading that exists.
**Validation:** User confirmed. The maintainer approved the triage (F4 accepted as worded) and said "capture and run".
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (6587 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` (REQ-644: a disambiguation added in one action stales sibling restatements) is exactly the gap this REQ closes for the schema line and the run action.
## Full Context
See `do-work/user-requests/UR-142/input.md` for complete verbatim input (the consumer incident timeline and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, UR-142, crew members (general, coding-guardrails, shared-principles, communication-style), prime-action-files (Cross-Referencing spelling: same-package `actions/<file>.md`, sections as `→ **Heading**`), prime-releases, and every `alternate-writer-contract-drift` lesson. Approach: one sentence inserted into the schema comment right after the "who writes it" sentence; one new table row directly after the missing-precondition row; replace the dangling pointer in the existing row with the real section pointer and reuse it in the new row. Confirmed the cited target exists: `actions/capture.md` § Step 1: Parse and Assess holds the **Earmark assessment** bullet (line 107), and `work-reference.md` holds `## Failure Classification (Step 8)` with the Environment row carrying the blocked-flip test. No test pins the old pointer text.
- [x] **[APPLY]:** Edited exactly the two write-set files with an anchored, count-asserted replacement. No other file touched.
- [x] **[UNIFY]:** `git diff --stat`: work-reference.md 2 +-, work.md 3 ++-, 2 files, 3 insertions, 2 deletions. `git diff --check` exit 0. `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 (PASS). `bash _dev/tests/contract-regressions.sh` exit 0 ("Contract regression checks passed."). Files checked: work-reference.md (line 114 keeps every prior sentence, bold labels unchanged, the new sentence adds no case for writing `assigned_to`); work.md (table keeps two columns, every existing row unchanged except the pointer text, the new row sits beside the missing-precondition row). No debug artifacts.
*Source: consumer feedback "the marker rule is stated only in capture" (Gap 2 / the ask §2), accepted by `do-work-toolbox validate-feedback` on 2026-10-08 as F4, plus the triage's adjacent observation about the dangling "Mid-run blocked flip" pointer.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the two files and lines (work-reference.md:114 schema line, work.md Error Handling table row ~522) and quotes both sentences; prose only, plus one pointer fix the triage located.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)

**What was done:** The `assigned_to` schema line now says the field is for another session or checkout only, and that operator-gated work is `status: blocked` with `blocked_by` naming the user. The `work.md` Error Handling table gains a row that routes an orchestrator's wish to park a queued REQ on the operator to the mid-run blocked flip and forbids writing `assigned_to`, and the existing missing-precondition row now points at `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row, instead of a heading that does not exist.

## Decisions

*(from the builder hand-back)*

- D-01: The schema sentence sits right after "Seeded by capture … or written by a session claiming from another checkout." and before "**Not a lock and not a status:**". It reads as a direct qualifier of who may write the field, and it leaves the closing board-placement and parser lock-step sentences intact as the line's tail.
- D-02: The pointer wording is `` `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row ``, matching the prime's same-package spelling and the existing `work.md` habit of `` `actions/work-reference.md` → **Heading** ``. The brief's parenthetical "(the blocked-flip test)" was left out of the inline pointer because the existing row already opens with "Apply the blocked-flip test", and the new row says "the same mid-run flip as a missing precondition".
- D-03: In the new row the pointer goes in parentheses after "the same mid-run flip as a missing precondition", so the REQ's quoted sentences stay verbatim around it.

## Discovered Tasks

*(from the builder hand-back)*

- impact-low: `skills/do-work/actions/work-reference.md:788` (Failure Classification, Environment row) says "see `actions/work.md`'s mid-run blocked-flip procedure", and `:165` says "Step 8's blocked-flip procedure". `work.md` has no such procedure; after this REQ, `work.md` points back at this very row, so the two citations now form a loop with no procedure text beyond the row itself. The brief and REQ kept these prose mentions on purpose. A later REQ could drop the "see work.md" clause at :788. → report only

**Integrator consolidation (Step 8):**
- The builder's item above is the same as review F1. Restamped with the house token: impact-negligible, because the rule text at `work-reference.md:788` is complete and the loop costs one jump → report only
- `_dev/primes/prime-action-files.md:21` names `pipeline.md` as a state-based action, and no package has that file any more (found by the Orientation staleness spot-check, pre-existing) — impact-negligible → report only

## Qualification

**Gate record:** `advance --diff-range 336cbb52..4c50ccde` reported the `qualify` gate `satisfied` on the merged range, with no findings (no debug artifacts, no output primitives, P-A-U boxes ticked).

**Requirement trace against `git diff 336cbb52..4c50ccde` (2 files, 3 insertions, 2 deletions):**
1. `work-reference.md:114`: the `assigned_to` line gains the quoted "**For another session or checkout only** — …" sentence verbatim, right after "Seeded by capture …", before "**Not a lock and not a status:**". Every prior sentence of the line, including the board-placement and parser lock-step sentences, is byte-identical. Met.
2. `work.md` Error Handling table: the new "Orchestrator wants to park a queued REQ on the operator …" row is present with the REQ's quoted text, plus a parenthetical pointer after "the same mid-run flip as a missing precondition" (D-03). The table keeps two columns. Met.
3. Pointer fix: both rows point at `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row. That heading exists at `work-reference.md:765` and its Environment row opens "First apply the blocked-flip test". `grep -rn "Mid-run blocked flip" skills/` prints nothing. The prose mentions at `work-reference.md:165` and `:788` are unchanged, as the REQ asked. Met.
4. `capture.md`, `capture-reference.md`, `docs/work-guide.md`, and `model.go` are absent from the diff. Met.
5. `grep -n operator` on both files hits `work-reference.md:114` and `work.md:523`. Met.
6. Contract checks: covered by the probe and repository gate in Testing.
7. Release: handled at finalization.

**Scope (Route A, no `## Scope`):** touched files equal `write_set` exactly (`skills/do-work/actions/work-reference.md`, `skills/do-work/actions/work.md`). No drift. The new sentences say when not to write `assigned_to` and add no case for writing it.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (repository gate, at merge `4c50ccde`, unpiped, wall 115s); focused probe `do-work/runs/work-2026-10-08-155307/REQ-648-probe.sh` run by `advance` (the two `operator` greps plus `bash _dev/tests/shipped-package-reference-contract.sh`).
**Result:** ✓ Repository gate exit 0 on the first run ("Maintainer verification passed."). `advance` recorded `green-gate` satisfied and the probe satisfied (`BLOCKED-PROBE-SUCCEEDED`, raw status 0).

**Regression evidence (non-behavioral prose change, `tdd: false`):** the Red-Green Proof is a read check, verified on the merged diff. Before: `grep -c operator` was 0 in both files, and `work.md` pointed at "Step 8 → **Mid-run blocked flip**", which does not exist. After: `grep -n operator` hits `work-reference.md:114` and `work.md:523`, and `grep -rn "Mid-run blocked flip" skills/` prints nothing.

**New tests added:** None.

**Heavy verification plan:**
- Range: 336cbb52dbdf2e91dea27b96f7e783ee40e9f238..4c50ccdecf64fc63ec429ca252768252e32faf83
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — `skills/do-work/actions/work-reference.md` and `skills/do-work/actions/work.md` matched subtree skills

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-08T16:05:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1: `skills/do-work/actions/work-reference.md:788` cites "`actions/work.md`'s mid-run blocked-flip procedure", which does not exist; `work.md:522-523` now point back at that row, so the citations loop (rule text is complete at `:788`; kept by DR3 on purpose; fix is deleting the "see work.md" clause) — impact-negligible → report only. Nit F2: `work.md:523` says the default run "hides" an earmarked REQ, while the schema line and `docs/work-guide.md:120` say it skips and reports; wording accepted as quoted by the maintainer — impact-negligible → report only
**Acceptance:** Pass — `grep -n operator` hits both files, target heading exists, reference contract PASS, `git diff --check` clean
**Restatement sweep:** redefined `assigned_to` audience (session or checkout only, operator work is `blocked`): `actions/capture.md:107`, `capture.md:108`, `docs/work-guide.md:120`, `work.md:37` agree; `actions/capture-reference.md:41`, `capture.md:143` state mechanics only, no conflict
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Putting the "not for this" sentence on the field's own schema line. A session that writes `assigned_to` outside capture reads that line first, so the rule now meets it there.
**What didn't:** The old `work.md` row pointed at "Step 8 → **Mid-run blocked flip**", a heading that never existed. Nothing tested the pointer, so it stayed broken until a triage read it.
**Worth knowing:** `work-reference.md:788` still sends readers to a "`work.md` mid-run blocked-flip procedure", and `work.md` now points back at that row. The rule text is complete at `:788`, so the loop costs one jump (review F1, report only).

## Orientation

The run action's error table and the `assigned_to` schema line now both say operator-gated work is blocked, never earmarked; lives in the work action's reference docs (`_dev/primes/prime-action-files.md` area). Staleness spot-check: this change made no prime stale. One pre-existing stale mention: `_dev/primes/prime-action-files.md:21` names a `pipeline.md` state-based action that no longer exists in any package (impact-negligible → report only).

## Heavy Verification Plan

- Base: 336cbb52dbdf2e91dea27b96f7e783ee40e9f238
- Target: 4c50ccdecf64fc63ec429ca252768252e32faf83
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — `skills/do-work/actions/work-reference.md` and `skills/do-work/actions/work.md` matched subtree skills

## Heavy Verification Result

- Target revision: 4c50ccdecf64fc63ec429ca252768252e32faf83
- Execution revision: 4c50ccdecf64fc63ec429ca252768252e32faf83 (detached checkout at the merge)
- staged-skills: exit 0, executed (no_prior_evidence), 43s wall

## Timing

Observed 2026-10-08T15:56:02Z to 2026-10-08T16:06:42Z: 10m 40s total, 9m 34s attributed across 5 events, 1m 06s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 4m 32s | 1 |
| verification-gate | 3m 19s | 2 |
| review | 1m 29s | 1 |
| handback-merge | 14s | 1 |

Slowest stage: builder-work / builder worktree build, 4m 32s, outcome success.
