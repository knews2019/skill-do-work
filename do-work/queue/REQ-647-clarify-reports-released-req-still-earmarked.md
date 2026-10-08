---
id: REQ-647
title: 'Clarify tells the user when a released REQ stays earmarked and clears it on request'
status: pending
created_at: 2026-10-08T15:50:17Z
user_request: UR-142
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-648, REQ-649]
batch: earmark-release-visibility
write_set: [skills/do-work/actions/clarify.md, skills/do-work/docs/work-guide.md]
---
# Clarify Tells the User When a Released REQ Stays Earmarked and Clears It on Request
## What
`do-work clarify` releases a `blocked` REQ through the `unblock` transaction, which removes `blocked_by` and `blocked_at` and leaves `assigned_to` in place. `actions/clarify.md` Step 5.5 then tells the user "The REQ re-enters the queue for the next `do-work run`", which is false when the REQ is earmarked: the default run skips it and the board shows it under Pending → Earmarked. Make clarify say so at the moment of release, offer to clear the earmark (a hand edit of the field, the documented clear path), name every released-but-earmarked REQ in the report, and fix the false sentence.
## Why
Consumer incident, Oct 8 08:56 UTC on 0.305.79: two REQs that needed the operator carried both `status: blocked` and `assigned_to: 'user-interactive'`. The user released both through clarify ("yes, unblock"). They returned to `pending`, landed under Pending → Earmarked, and every default run skipped them without saying so. The user waited, then asked why REQs that needed their input sat under Pending. The skill handed back a REQ it would not run and did not say so. Triage (validate-feedback, 2026-10-08) accepted the prose fix (Option A) and pushed back on a typed `--clear-assigned-to` flag: the hand edit is already the documented clear path, so a CLI flag adds surface for one advisory field.
## Verified Facts (from triage)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:604-621` (`TransitionUnblock`): sets `status: pending`, stamps `status_changed_at`, deletes `blocked_by` and `blocked_at`, appends the Blocked history line; never reads `assigned_to`. The only writer that deletes the field is an explicit claim (`state_apply.go:562-565`). Both stay as they are.
- `skills/do-work/actions/clarify.md:188` ("Yes → unblock") ends "The REQ re-enters the queue for the next `do-work run`." `grep -n assigned_to skills/do-work/actions/clarify.md` finds nothing.
- `skills/do-work/actions/work-reference.md:114` (`assigned_to` schema line): "an assignment persists until an explicit run or a hand-edit clears it". `work-reference.md:511` (exit-summary row 6): "To drop an assignment without running it, remove the field by hand." The hand edit is the documented clear path, not a free-form reproduction of a transaction.
- Clarify already hand-edits queued REQ files in the same flow (Step 4 item 2 reclaim lines, the Builder Was Right / Discarded fast paths).
- `skills/do-work/docs/work-guide.md:120` (Earmarking with `assigned_to`): names the two ways to take an earmarked REQ (run it by name, or delete the field by hand) and says operator-wait work is captured `blocked` and released by clarify; it does not say what happens when a released REQ is also earmarked.
- No `_dev/tests` lock-in pins the sentence at `clarify.md:188` (grep for "re-enters the queue" across `_dev/tests/` finds nothing).
## Detailed Requirements
1. `actions/clarify.md` Step 5.5, the "Yes → unblock" bullet: keep the transaction sentence and the "do not reproduce the mutation free-form" rule unchanged. After the transaction, add the condition: when the released REQ carries `assigned_to`, tell the user on the spot that the REQ is still earmarked for `<session name, verbatim>` and that the default run will skip it, and ask one question with two concrete options (`crew-members/clear-questions.md` is already loaded by Step 3): clear the earmark now (remove the `assigned_to` line by hand edit, the documented clear path in `actions/work-reference.md` → Request File Schema `assigned_to` and the exit summary's assigned-elsewhere row), or keep it and run the REQ by name with `do-work run REQ-NNN`. Replace the closing sentence with: "The REQ re-enters the queue for the next `do-work run` unless it carries `assigned_to`, in which case the default run skips it until it is named explicitly or the field is cleared."
2. Say how the hand edit is committed: the same way clarify's other hand edits in the same invocation are committed today (read the file; do not invent a new commit step or a new transaction).
3. `actions/clarify.md` Step 6 (Report): the summary lists every released REQ that stays earmarked, with the verbatim session name and the run-by-name command. One sentence keyed on the condition, not a new subsection.
4. `docs/work-guide.md` Earmarking paragraph: one sentence saying that when clarify releases a blocked REQ that is also earmarked, it says so and offers to clear the earmark. Keep the paragraph's existing sentences.
5. Verification Checklist in `clarify.md`: one line that every unblocked REQ still carrying `assigned_to` was named to the user and in the report. Only if the checklist's existing shape takes it without ceremony; otherwise skip and say so in the hand-back.
6. Run `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`; both green.
7. Release per `_dev/primes/prime-releases.md` (the integrator writes the CHANGELOG entry and bump; the title says what shipped in plain words).
## Constraints
- Prose only. No Go change, no new CLI flag, no new transaction, no `verify` probe (both were pushed back at triage; see the UR). `TransitionUnblock` and the default scan stay as they are.
- `assigned_to` stays verbatim-read: the question quotes the session name as written, no normalization.
- Clearing stays an explicit act: never clear the field silently on unblock, and never clear it when the user picks "keep".
- Do not edit `actions/work-reference.md` (REQ-648 owns its `assigned_to` schema line; the exit-summary row 6 is unchanged) or `actions/capture.md`.
- Read `_dev/primes/prime-action-files.md` before editing; every cross-reference uses that prime's spelling (`actions/work-reference.md` → section name).
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-648 or REQ-649 into it.
## Dependencies
None. REQ-648 (the schema line and the work action state the operator rule) and REQ-649 (board tooltip) touch different files; the three may run in parallel.
## Builder Guidance
Certainty is high: the triage read the transaction and the action text, and the remedy is two conditions in one action file plus one sentence in the user guide. Latitude: exact wording of the question and the report line, and whether the checklist line is worth adding.
## Red-Green Proof
**RED prompt/case:** A queued REQ with `status: blocked`, `blocked_by: 'the operator'`, `blocked_at`, and `assigned_to: 'user-interactive'`. Run `do-work clarify` and answer "Yes, unblock it".
**Why RED now:** `clarify.md` Step 5.5 runs the unblock transaction and says the REQ re-enters the queue. Nothing reads `assigned_to`, so the user is not told the default run will skip it, is not offered the clear, and the Step 6 report does not name it.
**GREEN when:** A cold read of `clarify.md` Step 5.5 routes that case to a two-option question (clear the earmark by hand edit, or keep it and run by name), the closing sentence carries the `assigned_to` exception, and Step 6 names every released-but-earmarked REQ with its verbatim session name. `grep -n assigned_to skills/do-work/actions/clarify.md` finds the Step 5.5 condition and the Step 6 line. `shipped-package-reference-contract.sh` and `contract-regressions.sh` stay green.
**Validation:** User confirmed. The maintainer approved the triage's Option A remedy and said "capture and run".
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (6587 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` (REQ-644: a rule added to one action stales its restatements elsewhere) is the trap this REQ's work-guide sentence guards against.
## Full Context
See `do-work/user-requests/UR-142/input.md` for complete verbatim input (the consumer incident timeline and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer feedback "unblock keeps assigned_to, and clarify says nothing about it" (Gap 1 / the ask §1, Option A), accepted by `do-work-toolbox validate-feedback` on 2026-10-08 as F1 and F2; Option B and the verify probe pushed back.*
