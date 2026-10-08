# Hand-back — REQ-647 (Clarify tells the user when a released REQ stays earmarked and clears it on request)

- **Branch:** worktree-agent-REQ-647-clarify-reports-released-req-still-earmarked
- **Worktree:** /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-647-clarify-reports-released-req-still-earmarked
- **Base:** dd54712a
- **Commits:** fe2c1878 `[REQ-647] clarify reports a released REQ that stays earmarked and offers to clear it` (single commit)

## File manifest

- `skills/do-work/actions/clarify.md` (modified): Step 5.5 "Yes → unblock" gains the `assigned_to` condition with a two-option question and the commit clause; the closing sentence carries the `assigned_to` exception; Step 6 names every released REQ that stays earmarked; Verification Checklist gains one line.
- `skills/do-work/docs/work-guide.md` (modified): the Earmarking with `assigned_to` paragraph gains one closing sentence saying clarify flags a released-but-earmarked REQ and offers to clear the earmark. Every existing sentence kept.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the brief, REQ-647, crew members (general, coding-guardrails, shared-principles, communication-style), prime-action-files, prime-releases, and every `alternate-writer-contract-drift` lesson. Read clarify.md in full to find its commit discipline: Step 4/5 answers go through `answer --manifest … --commit`; the Step 5.5 `unblock` invocation carries no `--commit`; the Step 5.5 reclaim line and the Builder Was Right / Discarded fast paths state no commit step at all. So requirement 2 takes the brief's fallback: say clarify states no commit step for hand edits and point at `actions/commit.md` in one clause. Edits: (1) extend the "Yes → unblock" bullet in place after the "do not reproduce" sentence, (2) one sentence appended to Step 6, (3) one checklist line beside the existing blocked-REQ line, (4) one sentence appended to the work-guide Earmarking paragraph. No new bold labels or headings; cross-references in the prime's `actions/work-reference.md` → **section** spelling.
- [x] **[APPLY]:** Applied exactly the four planned edits to the two write-set files. Option wording: "Clear the earmark now" (hand edit, documented clear path cited to Request File Schema `assigned_to` and the Composed Exit Summary assigned-elsewhere row) and "Keep it, and run the REQ by name with `do-work run REQ-NNN`". Added "Never clear the field without that answer, and leave it untouched on keep." The closing sentence is the REQ's text verbatim.
- [x] **[UNIFY]:** `git diff --stat`: `clarify.md | 9 +++++++--`, `work-guide.md | 2 +-`, 2 files changed, 8 insertions(+), 3 deletions(-). `git diff --check` clean. `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 ("shipped package reference contract: PASS"). `bash _dev/tests/contract-regressions.sh` exit 0 ("Contract regression checks passed."). Files checked: clarify.md (Step 5.5 bullet reads cold, the transaction sentence and "do not reproduce the mutation free-form" are unchanged, Step 6 line, checklist line, all bold labels and headings unchanged); work-guide.md (paragraph intact, new sentence last). No debug artifacts.

## Red-green evidence

**RED (clarify.md at dd54712a):** `grep -n assigned_to skills/do-work/actions/clarify.md` returned nothing. The "Yes → unblock" bullet ended "The REQ re-enters the queue for the next `do-work run`." Step 6 said nothing about earmarks. A cold read of the RED case (blocked REQ with `assigned_to: 'user-interactive'`, user answers "Yes, unblock it") runs the transaction and tells the user the REQ is back in the queue, which is false for the default run.

**GREEN (clarify.md at fe2c1878):** the same cold read now hits "The transaction never touches `assigned_to`, so when the released REQ carries it, tell the user on the spot that the REQ is still earmarked for `<assigned_to, verbatim>` and that the default `do-work run` will skip it, then ask one question with two options", followed by the clear and keep options. The closing sentence reads "The REQ re-enters the queue for the next `do-work run` unless it carries `assigned_to`, in which case the default run skips it until it is named explicitly or the field is cleared." Step 6 reads "When a REQ unblocked in Step 5.5 still carries `assigned_to`, the summary names it with the session name verbatim and the run-by-name command `do-work run REQ-NNN`."

```
188: ... The transaction never touches `assigned_to`, so when the released REQ carries it, tell the user on the spot ...
189:   1. Clear the earmark now: remove the `assigned_to` line by hand edit, the documented clear path ...
192:   Never clear the field without that answer, and leave it untouched on keep. The REQ re-enters the queue for the next `do-work run` unless it carries `assigned_to`, ...
198: Summary of what was resolved ... When a REQ unblocked in Step 5.5 still carries `assigned_to`, the summary names it with the session name verbatim and the run-by-name command `do-work run REQ-NNN`.
275: - [ ] Every unblocked REQ still carrying `assigned_to` was named to the user at release and in the final report, and its field was cleared only on the user's answer.
```

## Decisions

- **D-01 (DECIDE & STATE) — commit discipline for the earmark hand edit (requirement 2).** Clarify has no stated commit discipline for hand edits: only Step 5's `answer` carries `--commit`; the Step 5.5 `unblock` invocation has no `--commit`, and the reclaim line and Builder Was Right / Discarded fast paths name no commit step. Per the brief's fallback, the new clause says so and routes the edit through `do-work commit` (`actions/commit.md`). Value: no invented commit step or transaction. Risk: low; the edit can sit uncommitted until the user commits, which is already true of the unblock itself.
- **D-02 (DECIDE & STATE) — checklist line (requirement 5).** Added. The checklist already has a per-condition line for blocked REQs ("`blocked` REQs the user confirmed satisfied flipped to `pending` …"), so one more line beside it fits the existing shape without ceremony. It also pins "cleared only on the user's answer", the never-silent constraint.
- **D-03 (DECIDE & STATE) — "keep" option names the claim side effect.** The keep option says `do-work run REQ-NNN` "claims it and clears the field", matching work-guide.md's existing sentence, so the user knows naming the REQ also drops the earmark.
- **D-04 (DECIDE & STATE) — placeholder spelling.** Used `<assigned_to, verbatim>` rather than `<session name, verbatim>` so it matches the exit-summary row 6 field cell ("assigned to <assigned_to, verbatim>") exactly.

## Discovered Tasks

- impact-low: `actions/clarify.md` Step 5.5 invokes `unblock` without `--commit` (the CLI accepts `--commit`, `internal/requeststate/state_commands.go`), while Step 5 commits its `answer` transaction; a clarify session that only unblocks leaves an uncommitted lifecycle change. → report only
- impact-low: `actions/clarify.md` "Builder Was Right / Discarded" section describes hand edits (status flip, archive, appended notes) while Step 5 says the `answer` command derives the disposition and "There is no hand-edit, helper, or manual Git fallback"; the two sections may disagree about who writes the fast-path result. → report only

## Lessons read

- Satellite: `_dev/primes/lessons-action-files.md`, every `[family: alternate-writer-contract-drift]` bullet (REQ-477, 498, 513, 461, 531, 566, 640, 641, 642, 644). Applied: swept the restated contract (the `assigned_to` clear path in work-reference schema line and exit-summary row 6, the work-guide Earmarking paragraph), and copied the placeholder from the row 6 source, not from the REQ text (REQ-644's lesson). `grep -rn "re-enters the queue" skills/` was not run beyond clarify; the REQ's verified facts state no `_dev/tests` pin.
- Prime traps: prime-action-files `## Traps` (alternate-writer-contract-drift, budgeted-context-routing); prime-releases `## Traps`.

## Proposed lesson bullet (integrator writes)

```
- [family: alternate-writer-contract-drift] [REQ-647: a transaction that clears part of a REQ's state can leave a sibling field that still steers the scheduler; the unblock prose said "re-enters the queue" while `assigned_to` kept the default run skipping it. When an action reports a state change, check every field the next reader filters on, not only the fields the transaction wrote](../../do-work/archive/UR-142/REQ-647-clarify-reports-released-req-still-earmarked.md#lessons-learned)
```

## Proposed CHANGELOG entry (integrator writes)

```
## Clarify Warns When a Released REQ Is Still Earmarked

When you unblock a REQ in `do-work clarify` and it also carries `assigned_to`, clarify now says so instead of claiming it is back in the queue. Before this, the default run skipped such REQs silently and they sat under Pending → Earmarked.

- Clarify names the session the REQ is still earmarked for and asks: clear the earmark now, or keep it and run the REQ by name.
- The clarify report lists every released REQ that stays earmarked, with the run-by-name command.
- The user guide's earmarking section explains this behavior.
```

## Integration seams

None. Write set was exactly the two files. REQ-648 owns `work-reference.md`; this change cites its Request File Schema `assigned_to` line and exit-summary row 6 without editing them.

## Test wall times

| Script | Exit | Wall time |
|---|---|---|
| `_dev/tests/shipped-package-reference-contract.sh` | 0 | 1 s |
| `_dev/tests/contract-regressions.sh` | 0 | 28 s |
