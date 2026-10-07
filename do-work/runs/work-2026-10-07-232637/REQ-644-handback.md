# Hand-back — REQ-644 (capture separates a session earmark from work that needs the user as operator)

- Branch: `worktree-agent-REQ-644-capture-earmark-vs-operator-blocked`
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-644-capture-earmark-vs-operator-blocked
- Commits: `b789ca74` ([REQ-644] capture separates a session earmark from work that needs the user as operator), on top of `531dd00e`.

## File manifest
- `skills/do-work/actions/capture.md` (modified) — Earmark assessment gains one sentence: `assigned_to` routes work between sessions; operator-present work is captured `blocked` per the External-condition assessment, with `blocked_by` naming the person in the user's words and `blocked_at`, landing in Needs Input · Blocked for `do-work clarify` to release. External-condition assessment: "the three look-alikes" → "every look-alike that is not blocked", and the earmark is added as the fourth (`assigned_to`, pointing back at the earmark assessment).
- `skills/do-work/docs/work-guide.md` (modified) — "Earmarking with `assigned_to`" opener changed from "leave this one for me" to "leave this one for that session"; one closing sentence: the field is for another session or checkout, work that waits on you as operator is captured `blocked` with `blocked_by` naming you and shows under Needs Input · Blocked until `do-work clarify` releases it.

## P-A-U
- **[PLAN]:** Read prime-action-files.md, general/coding-guardrails headings/shared-principles/communication-style crew, and the `alternate-writer-contract-drift` lesson bullets. Grep the restated earmark rule across `skills/` first. Add one sentence per site at capture.md Earmark assessment, capture.md External-condition look-alike list (replacing the closed count "three" with a condition-keyed phrase), and work-guide.md Earmarking paragraph (also fixing its ambiguous "leave this one for me" opener). No heading or bold label changes; no board Pending placement mentioned.
- **[APPLY]:** Three edits as planned in the two write_set files, committed as `b789ca74`. Bold labels **Earmark assessment**, **External-condition assessment**, **Earmarking with `assigned_to`.** unchanged. Never-invent rule untouched (the existing "Never invent or infer a session name" and the user's-words requirement for `blocked_by` both stand).
- **[UNIFY]:** `git diff HEAD~1 --stat`: capture.md 4 +/- 2 lines changed (2 ins, 2 del), work-guide.md (2 ins, 2 del); total 2 files, 4 insertions, 4 deletions. `git diff --check` clean. Checked: capture.md lines 107-108 (sentence placement, cross-reference spelling `actions/clarify.md`, Frontmatter Quoting still governs `blocked_by` via the External-condition bullet it points to); work-guide.md lines 113-119 (opener no longer reads as operator-present; new sentence true before and after REQ-643). No debug artifacts. `bash _dev/tests/shipped-package-reference-contract.sh` PASS (1s wall); `bash _dev/tests/contract-regressions.sh` PASS (21s wall).

## Red-Green (prose read)
- RED (before): "I need to be at the keyboard for this one, hold it for me" matched only the Earmark assessment, so capture seeded `assigned_to`.
- GREEN (after): the Earmark assessment now names that exact case and routes it to `blocked` via the External-condition assessment, with `blocked_by` (the person, user's words) and `blocked_at`. "leave this for cloud-alpha" still names a session, so it still seeds `assigned_to`, and the External-condition list now explicitly says a session reservation is `assigned_to`, not blocked.

## Restatement grep accounting
Commands (from worktree root):
1. `grep -rniE 'earmark|leave this|for the laptop' skills/`
2. `grep -rn 'assigned_to' skills/ --include='*.md'`
3. `grep -rn 'look-alike' skills/`

Hits:
- capture.md:107, :108 — changed (the two sites).
- capture.md:143 (addendum: "apply Step 1's earmark assessment") — defers to Step 1, so it inherits the new sentence; unchanged. See Discovered Tasks.
- work-guide.md:113, :116, :119 — the Earmarking paragraph; changed (113, 119); 116 is the YAML example.
- work-guide.md:99, :132 ("not earmarked for someone else/elsewhere" in the fan-out ready set) — describe the scan predicate, not when to use the field; unchanged, still true.
- work-guide.md:128 (verify flags a REQ built here that is still earmarked elsewhere) and :138 (`--skip-impact-negligible` still respects another session's earmark) — mechanics, not when to use the field; unchanged.
- capture-reference.md:41 (sweep conversion takes only unassigned REQs) — mechanics; unchanged.
- work-reference.md:77, :170, :259 — Frontmatter Quoting / stakeholder / verbatim-read contract mentions, mechanics only; unchanged.
- work-reference.md:114 (schema line) and :511 (Composed Exit Summary row 6) — field mechanics, out of boundary as the brief says; unchanged.
- work-reference.md:460, work.md:37, work.md:108 — scan/fan-out/skip predicates; mechanics; unchanged.
- Board Go/JS/CSS hits (generate.go:182, model.go:176, model_test.go, verify.go, verify_test.go, board-cards.js:220/234, board.css:1298) — display and verify mechanics; REQ-643 territory; unchanged.
- Look-alike grep: only capture.md:108.
No other file restates when to use the earmark; no out-of-boundary change needed.

## Decisions
- D-01 DECIDE & STATE: replaced "the three look-alikes" with "every look-alike that is not blocked" rather than "four". Value: condition-keyed, cannot go stale when another look-alike is added. Risk: none; the list still enumerates each.
- D-02 DECIDE & STATE: changed the work-guide opener "leave this one for me" to "leave this one for that session". Value: removes the operator-present reading the REQ names as the root ambiguity. Risk: low; the bold label is unchanged, so citation pins are unaffected (reference contract passes).
- D-03 DECIDE & STATE: the capture sentence says "the user's words" for `blocked_by` and points at the External-condition assessment for quoting and timestamp mechanics, instead of restating Frontmatter Quoting. Value: one home for the mechanics. Risk: none.

## Discovered Tasks
- impact-negligible: capture.md:143 (addendum path) only handles an addendum that names a session; an addendum saying the user must be present as operator would need the pending REQ flipped to `blocked` with `blocked_by`/`blocked_at`, and the addendum rule does not say so. Not observed in practice → report only

## Lessons read
- `_dev/primes/lessons-action-files.md` bullets tagged `[family: alternate-writer-contract-drift]` (REQ-477, 498, 513, 461, 531, 566, 640, 641, 642), and the prime's Traps line for that family.

## Proposed lesson bullet (orchestrator writes)
- [family: alternate-writer-contract-drift] [REQ-644: a disambiguation added to a capture assessment also stales any closed count in the sibling assessment that lists its look-alikes ("the three look-alikes"), and the user-guide example phrase ("leave this one for me") that seeded the ambiguity; key the list on its condition and fix the guide's example wording, not only add the sentence](../../do-work/archive/UR-140/REQ-644-capture-earmark-vs-operator-blocked.md#lessons-learned)

## Proposed CHANGELOG entry (orchestrator writes)
**Capture Tells a Session Earmark Apart From Work That Needs You at the Keyboard**
Capture now writes `assigned_to` only when another session or checkout will take the work. When you have to be present yourself, it captures the REQ as blocked on you, so it shows in Needs Input · Blocked and `do-work clarify` releases it.
- `actions/capture.md`: the Earmark assessment states the two cases; the External-condition assessment lists the earmark as a look-alike that is not blocked.
- `docs/work-guide.md`: the Earmarking paragraph says the same in one sentence.

## Integration seams
None. The text does not describe the board's Pending placement of an earmark, so it stays true before and after REQ-643 (Pending → Earmarked group).

## Test wall times
- shipped-package-reference-contract.sh: PASS, 1s
- contract-regressions.sh: PASS, 21s
