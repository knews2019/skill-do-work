---
id: UR-142
title: 'Clarify, the schema line, and the board stop hiding a released REQ that is still earmarked'
created_at: 2026-10-08T15:50:17Z
requests: [REQ-647, REQ-648, REQ-649]
word_count: 2569
---
# Three Accepted Findings: Clarify Names a Still-Earmarked Release, the Schema and Work Action State the Operator Rule, the Board Badge Names the Blocked Column

## Summary
The maintainer ran `do-work-toolbox validate-feedback` on a consumer-repo incident report (two REQs that needed the operator carried both `status: blocked` and `assigned_to`; `do-work clarify` released them, the `unblock` transaction kept `assigned_to`, both landed under Pending → Earmarked, and every default run skipped them without saying so), approved the triage, and said "capture and run". Accepted: the clarify prose fix (Option A: ask on release, clear by hand edit on request, name the REQ in the report, fix the false "re-enters the queue" sentence), one rule sentence each on the `assigned_to` schema line and in the work action's error-handling table, and one sentence on the board's Earmarked badge tooltip. Pushed back and out of scope: a typed `--clear-assigned-to` flag on `unblock` and a new `verify` probe (the hand edit is the documented clear path; the probe could not tell a deliberate keep from a forgotten one).

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-647 | Clarify tells the user when a released REQ stays earmarked and clears it on request |
| REQ-648 | [impact-rule-change] The assigned_to schema line and the work action say operator-gated work is blocked, never earmarked |
| REQ-649 | Board Earmarked badge tooltip names Needs input · Blocked as the home for operator-gated work |

## Batch Constraints
- No dependency between the three REQs: disjoint files (`actions/clarify.md` + `docs/work-guide.md`; `actions/work-reference.md` + `actions/work.md`; `queue-kanban/web/board-cards.js`), may run in parallel; integration and release stay serial.
- Each REQ is its own release per `_dev/primes/prime-releases.md`; do not fold one REQ's change into another.
- Prose and text only: no Go change, no new CLI flag, no new transaction, no `verify` probe. `TransitionUnblock` and the default scan stay as they are; `assigned_to` stays advisory and verbatim-read; clearing an earmark stays an explicit act.
- Surface-cost N/A for all three per the triage; the two flagged remedies (Option B, the verify probe) are not captured.
- Read `_dev/primes/prime-action-files.md` before REQ-647 and REQ-648; `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` before REQ-649.

## Full Verbatim Input
> ```
> capture and run
> 
> [Capture source: the maintainer ran `do-work validate-feedback` on a consumer-repo incident report (pasted, with a board screenshot showing two REQs under Pending → Earmarked with `ASSIGNED user-interactive` badges), approved the triage, then said "capture and run". The three accepted-finding payloads from the triage come first, then the original feedback verbatim, then the triage verdicts.]
> 
> --- Accepted-finding capture payloads (from the approved triage) ---
> 
> do-work capture-request: Clarify tells the user when a released REQ stays earmarked and clears it on request (Option A: clarify.md Step 4 question + hand-edit clear, Step 6 lists every released-but-earmarked REQ by verbatim session name, fix the false "re-enters the queue" sentence at clarify.md:188). Source: consumer incident Oct 8 08:56 UTC, validate-feedback F1/F2 accepted; Option B (--clear-assigned-to) and the verify probe pushed back because the hand edit is the documented clear path (work-reference.md:114, :511) and the probe would flag a deliberate keep.
> 
> do-work capture-request: State the earmark-vs-operator rule where the field is defined: append one sentence to the assigned_to schema line (work-reference.md:114) and add a symptom-table row to work.md beside line 522 for an orchestrator that wants to park a queued REQ on the operator. Acceptance: grep -n operator on both files finds both. Source: validate-feedback F4 accepted. Optional: repoint the row's "Step 8 → Mid-run blocked flip" to work-reference.md Failure Classification → Environment.
> 
> do-work capture-request: Board Earmarked badge tooltip (board-cards.js:239) names Needs input · Blocked as the home for operator-gated work; no group-header hint exists, add none. Source: validate-feedback F5 accepted.
> 
> --- Original feedback (verbatim) ---
> 
> ## Request
> 
> 0.305.78 taught capture to tell a session earmark (`assigned_to`) apart from work that needs the user at the
> keyboard (`status: blocked` + `blocked_by` naming them), and 0.305.79 gave earmarked REQs their own Pending
> group on the board. Both are right. Three gaps remain, and together they reproduced the exact confusion
> those releases were meant to end:
> 
> 1. `unblock` keeps `assigned_to`, and `clarify` says nothing about it, so a REQ the user just released lands
>    in Pending → Earmarked where every default run skips it. The user sees it as "pending" and waits.
> 2. The rule "sessions get an earmark, the operator gets a block" is stated only in `actions/capture.md` and
>    `docs/work-guide.md`. The field's own schema line and the run orchestrator's guidance do not carry it, so a
>    session that touches a REQ outside capture has no rule to follow and improvises.
> 3. The board's Earmarked badge says what the group does, not what it is not for.
> 
> Please close all three. Nothing here changes queue semantics: `assigned_to` stays advisory, `blocked` stays an
> external condition, and clearing an earmark stays an explicit act.
> 
> ## What happened, in one consumer repo (timeline)
> 
> The consumer ran 0.305.71 when it started and 0.305.79 by the end.
> 
> | When (UTC) | Suite | Event |
> | --- | --- | --- |
> | Oct 7, 22:58 | 0.305.71 | A run orchestrator, mid-run, hand-wrote `assigned_to: 'user-interactive'` on two pending REQs that need the operator's production and content-store access. No action told it to; 0.305.71's capture rule had no sentence separating the two markers, and no mid-run rule exists at any version. |
> | Oct 7, 23:15 | 0.305.71 | A second session saw the two REQs under Pending → Ready (pre-0.305.79 placement) while the run skipped them, and flipped both to `status: blocked` with `blocked_by` naming the operator. It kept `assigned_to`. |
> | Oct 8, 08:50 | 0.305.79 | `do-work update` landed 0.305.78 and 0.305.79. |
> | Oct 8, 08:56 | 0.305.79 | The user released both through `do-work clarify` ("yes, unblock"). The `unblock` transition removed `blocked_by` and `blocked_at` and left `assigned_to`. Both REQs returned to `pending` and appeared under Pending → Earmarked. |
> | Oct 8, 14:23 | 0.305.79 | The user asked why two REQs that need their input sit in Pending instead of Needs input · Blocked. A session re-blocked them by hand, which is the state 23:15 had already produced. |
> 
> The user did everything the skill asked. The skill handed back a REQ it will not run and did not say so.
> 
> ## Where the behaviour lives today
> 
> **Gap 1 — unblock keeps the earmark, clarify is silent.**
> 
> - `tools/do-work-cli/internal/requeststate/state_apply.go:604-621` (`TransitionUnblock`): sets `status: pending`,
>   stamps `status_changed_at`, deletes `blocked_by` and `blocked_at`, appends the Blocked history line. Never
>   reads `assigned_to`.
> - `state_apply.go:562-565`: `assigned_to` is deleted only on a claim with `Provenance == ProvenanceExplicit`.
>   Correct for claim; it is the only writer that clears the field.
> - `actions/clarify.md:188` ("Yes → unblock"): "The REQ re-enters the queue for the next `do-work run`." That
>   sentence is false when `assigned_to` is set: the next default run skips it and reports it under
>   "assigned to another session" (`actions/work-reference.md:511`).
> - The sibling probe already exists on the other side of the claim: `tools/queue-kanban/verify.go:908-923`
>   (`appendAssignedElsewhereFindings`) flags a REQ in `working/` that still carries `assigned_to`. There is no
>   equivalent for "released by clarify but still earmarked".
> 
> **Gap 2 — the marker rule is stated only in capture.**
> 
> - `actions/capture.md:107` (Earmark assessment) carries the rule: "`assigned_to` routes work between sessions;
>   when the user instead has to be present as operator … capture the REQ `blocked` … with `blocked_by` naming
>   that person … and `blocked_at`, which lands it in Needs input · Blocked for `do-work clarify` … to release."
> - `docs/work-guide.md:120` ends with the same rule.
> - `actions/work-reference.md:114`, the schema line for `assigned_to`, defines the field as "the session this
>   REQ is earmarked for … Seeded by capture when the user earmarks work, or written by a session claiming from
>   another checkout" and never says what it is not for. This is the line any session reads first when it
>   decides whether to write the field.
> - `actions/work.md` has the mid-run blocked flip (Step 8, and the row at `work.md:522`), which covers "a
>   builder reports a missing external precondition". It has no rule for an orchestrator that wants to park a
>   queued REQ on the operator without a builder having failed, which is the Oct 7 22:58 case.
> 
> **Gap 3 — the board badge.**
> 
> - `tools/queue-kanban/web/board-cards.js:234-239`: the Earmarked badge tooltip ends "The board only groups it
>   under Pending → Earmarked; it never reorders, blocks, or hides on this." A user reading the board learns
>   what the group does, not that a REQ waiting on them belongs one column to the right.
> 
> ## The ask
> 
> ### 1. `clarify` must not hand back a silently skipped REQ
> 
> Pick one of these two; the first is smaller.
> 
> **Option A — prose only.** In `actions/clarify.md` Step 4, after the "Yes → unblock" bullet, add: when the
> released REQ carries `assigned_to`, say so to the user on the spot ("REQ-NNN is still earmarked for `<name>`;
> the default run will skip it. Clear the earmark now, or keep it and run `do-work run REQ-NNN` by name?") and
> clear the field by hand edit on "clear". In the Step 5 summary, list every released REQ that stays earmarked
> with the verbatim session name. Fix the false sentence at `clarify.md:188`: "The REQ re-enters the queue for
> the next `do-work run` **unless it carries `assigned_to`, in which case the default run skips it until it is
> named explicitly or the field is cleared.**"
> 
> **Option B — typed.** Add `--clear-assigned-to` to the `unblock` transition in `state_apply.go` (delete the
> field in `TransitionUnblock` when set), surface it through the CLI, and have `clarify.md:188` pass it when
> the user answers "clear". Same prose in Step 4 and Step 5 as option A. This keeps the mutation inside the
> canonical transaction, which is the rule `clarify.md:188` already states ("do not reproduce the mutation
> free-form"), so option B is the one consistent with the skill's own contract.
> 
> Either way, add a `verify` probe beside `appendAssignedElsewhereFindings` (`verify.go:908`): a `pending` REQ
> whose Blocked history's last entry is "cleared by user via clarify" and which still carries `assigned_to`.
> Remedy text: "released by clarify but still earmarked for `<name>`; the default run skips it — clear the
> field or run it by name". Read-only, not fixable, same as its sibling.
> 
> ### 2. State the marker rule where the field is defined
> 
> One sentence each, so a session that reads only the schema or only the run action meets the rule.
> 
> - `actions/work-reference.md:114`, append to the `assigned_to` schema comment:
>   `**For another session or checkout only** — work that waits on the user as operator is \`status: blocked\`
>   with \`blocked_by\` naming them (the capture Earmark assessment), never an invented earmark.`
> - `actions/work.md`, next to the mid-run blocked flip (Step 8), add a row to the symptom table at
>   `work.md:522`:
>   `| Orchestrator wants to park a queued REQ on the operator (production access, credentials, a person at
>   the keyboard) | Flip it to \`status: blocked\` with \`blocked_by\` naming the operator and \`blocked_at\`,
>   the same mid-run flip as a missing precondition. Never write \`assigned_to\` for this: the earmark is for
>   another session, it is invisible as "waiting on you", and the default run hides it from the user by
>   skipping it. |`
> - `docs/work-guide.md:120` already has the rule; leave it.
> 
> ### 3. Board: say what Earmarked is not
> 
> `tools/queue-kanban/web/board-cards.js:239`, extend the tooltip's last sentence:
> `"… it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input ·
> Blocked (status: blocked), not here."`
> 
> If the board has a group-header hint for Pending → Earmarked, give it the same sentence.
> 
> ## Acceptance
> 
> - A REQ with `status: blocked` + `assigned_to`, released through `do-work clarify`: the user is told the REQ is
>   still earmarked and offered the clear; the Step 5 summary names it; `verify` reports it if the user keeps
>   the earmark and forgets.
> - `actions/work-reference.md:114` and `actions/work.md`'s symptom table each carry the one-sentence rule;
>   `grep -n "operator" actions/work-reference.md actions/work.md` finds both.
> - The Earmarked badge tooltip names Needs input · Blocked as the home for operator-gated work.
> - No change to what the default scan or `advance` selects, and `assigned_to` stays verbatim-read.
> 
> --- Triage (validate-feedback, 2026-10-08, approved by the maintainer) ---
> 
> F1 · `unblock` never clears `assigned_to`; `clarify` says nothing about it · Gap 1
> - Verdict: Accept
> - Evidence: `state_apply.go:604-621` (`TransitionUnblock`) sets `pending`, deletes `blocked_by`/`blocked_at`, appends the Blocked history line, never reads `assigned_to`. `state_apply.go:562-565` is the only writer that deletes it, and only on an explicit claim. `grep assigned_to clarify.md` → no matches. `clarify.md:188` says "The REQ re-enters the queue for the next `do-work run`", which `work-reference.md:511` (exit-summary row 6, assigned-elsewhere skip) contradicts when the field is set.
> - Surface-cost: N/A for the prose fix (corrects a false sentence, adds a question at the moment the state is created).
> 
> F2 · Fix shape: Option A (prose + hand edit) vs Option B (`--clear-assigned-to` on `unblock`) · The ask, §1
> - Verdict: Accept Option A; push back on Option B.
> - Evidence: The hand edit is already the documented clear path, not a free-form reproduction of a transaction: `work-reference.md:114` ("an assignment persists until an explicit run or a hand-edit clears it") and `work-reference.md:511` ("To drop an assignment without running it, remove the field by hand"). `clarify` already hand-edits REQ files in the same flow (Step 2 reclaim lines, Builder Was Right / Discarded sections). The "do not reproduce the mutation free-form" rule at `clarify.md:188` guards the unblock preimage check and stamps, which Option A leaves inside the CLI.
> - Surface-cost: Option B is Flagged. It adds a flag in `state_commands.go`, a `StateOptions` field, a plan-validation branch in `state_plan.go:170-181`, an apply branch, and a test, to delete one advisory field the schema already lets a hand edit remove. Incident is real (consumer Oct 8 08:56 UTC) but the cheaper remedy wins.
> - Remedy: In `clarify.md` Step 4 "Yes → unblock": after the transaction, if the REQ carries `assigned_to`, ask on the spot (still earmarked for `<name>`; clear now, or keep and run by name) and delete the field by hand on "clear". Replace the false closing sentence with the conditional one the feedback quotes. In Step 6 (the feedback says Step 5; this repo's report step is Step 6), list every released REQ that stays earmarked with the verbatim session name.
> 
> F3 · New `verify` probe: pending + last Blocked entry "cleared by user via clarify" + `assigned_to` · The ask, §1 tail
> - Verdict: Push back
> - Evidence: `verify.go:908-923` (`appendAssignedElsewhereFindings`) reads frontmatter only. The board model parses no body history: `model.go` reads `## Instances` for sweeps (`countSweepInstances`, line 928) and nothing else, so this probe needs a new `## Blocked` section parser. Once F2 ships, the state the probe would catch is one the user chose in clarify ("keep the earmark"), and the probe cannot tell a deliberate keep from a forgotten one: it would flag the user's own answer on every run, read-only and never fixable. The state is already visible twice: the Pending → Earmarked group (0.305.79) and the run's assigned-elsewhere exit row (`work-reference.md:511`).
> - Surface-cost: Flagged. A permanent probe plus a body parser, for a case F2 closes at the source and the board already shows.
> 
> F4 · The marker rule lives only in `capture.md:107` and `work-guide.md:120` · Gap 2
> - Verdict: Accept
> - Evidence: `grep -n operator skills/do-work/actions/work-reference.md skills/do-work/actions/work.md` → no matches. `work-reference.md:114` (the `assigned_to` schema line) says who seeds it and that it is advisory, never what it is not for. `work.md:522` covers only "builder reports a missing external precondition"; `grep -i park work.md` → nothing for an orchestrator parking a queued REQ on the operator, which is the Oct 7 22:58 UTC case. REQ-644's own restatement sweep listed `work-reference.md:114` as still agreeing, because the rule was additive then; this feedback is the first report that a session reads the schema line alone.
> - Surface-cost: N/A (one sentence each, keyed on the condition "work waits on the user as operator").
> - Remedy: Append the one-sentence rule to `work-reference.md:114` and add the symptom-table row to `work.md` beside line 522, both as the feedback words them. Leave `capture.md:107`, `capture-reference.md:41`, `capture.md:143`, and `work-guide.md:120` as they are.
> 
> F5 · Earmarked badge tooltip says what the group does, not what it is not for · Gap 3
> - Verdict: Accept (tooltip); group-header hint: nothing to extend
> - Evidence: `board-cards.js:234-239` ends "it never reorders, blocks, or hides on this." `makePendingGroup` (`board-cards.js:509-524`) renders a name and a count only; there is no header hint on any Pending group, so the conditional "if the board has a group-header hint" is false and adding one is new surface for a sentence the badge already carries.
> - Surface-cost: N/A (one sentence on an existing tooltip).
> - Remedy: Extend the tooltip's last sentence: "… it never reorders, blocks, or hides on this. A REQ waiting on the operator belongs under Needs input · Blocked (status: blocked), not here."
> 
> F6 · Acceptance criteria · The ask, Acceptance
> - Verdict: Accept with one line struck
> - Evidence: Bullets 1 (minus the `verify` clause), 2, 3, 4 follow from F2, F4, F5 and the constraints hold: `TransitionUnblock` and the scan are untouched, `assigned_to` stays verbatim-read. Bullet 1's "`verify` reports it if the user keeps the earmark and forgets" falls with F3.
> 
> Summary: Accept F1, F4, F5, F6; Accept narrowed F2 (Option A only); Push back F3 and Option B inside F2. Already done 0. Discuss 0.
> 
> Adjacent observation (not in the paste): `work.md:522` points at "Step 8 → Mid-run blocked flip", a heading that does not exist; the procedure is the Environment row of Failure Classification in `work-reference.md:788`. The new row from F4 should point where the existing row points, so this is a one-line pointer fix to take or leave.
> ```