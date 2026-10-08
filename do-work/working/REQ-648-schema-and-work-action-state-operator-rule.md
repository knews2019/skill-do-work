---
id: REQ-648
title: '[impact-rule-change] The assigned_to schema line and the work action say operator-gated work is blocked, never earmarked'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer feedback "the marker rule is stated only in capture" (Gap 2 / the ask §2), accepted by `do-work-toolbox validate-feedback` on 2026-10-08 as F4, plus the triage's adjacent observation about the dangling "Mid-run blocked flip" pointer.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the two files and lines (work-reference.md:114 schema line, work.md Error Handling table row ~522) and quotes both sentences; prose only, plus one pointer fix the triage located.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
