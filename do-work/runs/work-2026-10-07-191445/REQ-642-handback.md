# REQ-642 hand-back (wave-end decision reads the set-aside list; late set-aside named in the Decision Brief)

- Branch: `worktree-agent-REQ-642-wave-end-set-aside-gaps` (no name collision)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-642-wave-end-set-aside-gaps`
- Base: d2d3c316. Commits: `22c2268e` (only commit).

## File manifest
- `skills/do-work/actions/work-reference.md` (modified): **Delegated integration — the coordinator shape**, last sentence: the integrator brief's list gains "every wave member the coordinator has set aside (a set-aside member keeps its claim in `do-work/working/`, so without this list the integrator cannot tell it from one still building)". One clause.
- `skills/do-work/actions/work.md` (modified): Step 7 **How to run it**: the wave-end clause now says the orchestrator decides "last successful integration" from the archive plus the set-aside list it holds, which under delegated integration is the list the integrator brief carries (cites the coordinator-shape label). One new sentence for F2: a member set aside after the wave's last review ran means no review was passed the wave-end fact, so the sweep did not run; the orchestrator names that gap under HANDLED in the run's Decision Brief as its call not to re-run a review, and re-runs none.

No heading or bold label changed. No other file touched.

## P-A-U
- **[PLAN]:** Read general.md, coding-guardrails.md, shared-principles.md, communication-style.md, prime-action-files.md, and the alternate-writer-contract-drift bullets (lessons-action-files.md lines 52-63). Two in-place prose edits, no new mechanism: (1) one clause in the integrator-brief list in work-reference.md coordinator shape; (2) replace the "orchestrator knows which it set aside" clause in work.md Step 7 How to run it with the archive-plus-set-aside-list decision naming the integrator brief as the list under delegation, and add one F2 sentence pointing at Decision Brief → HANDLED. Keep the decision out of the **Restatement sweep (MUST)** paragraph (it says "passed below"). Then grep every restatement.
- **[APPLY]:** Both edits made exactly as planned, via exact-string replacement asserting a single match. Committed as 22c2268e.
- **[UNIFY]:** `git diff --stat HEAD~1`: work-reference.md 2 +-, work.md 2 +- (2 files, 2 insertions, 2 deletions). `git diff --check`: clean. `bash _dev/tests/shipped-package-reference-contract.sh`: PASS, 1.0 s wall. `bash _dev/tests/contract-regressions.sh`: PASS ("Contract regression checks passed."), 23.1 s wall. Files checked: work-reference.md line 470 (cold read: who writes the set-aside list = coordinator, when = in the gap before each integrator, why = claim stays in working/); work.md line 370 (cold read: who decides = orchestrator, from what = archive + held set-aside list / integrator brief, F2 = HANDLED in Decision Brief, no review re-run); restatement sites below.

## Decisions
- **D-01 (DECIDE & STATE) HANDLED bullet needs no edit.** The HANDLED bullet (work-reference.md line 871) says it lists DECIDE & STATE decisions (reversible `D-NN` entries), and the template line is "decided <Y> because <Z> ← reversible calls made without asking". The new sentence phrases the F2 note as the orchestrator's call not to re-run a review, which is a reversible call made without asking, so it fits the existing rule. Bullet left unchanged. Slight residual: the bullet says the entries are read from each REQ's `## Decisions`, while this one comes from the orchestrator itself; I judged that descriptive, not exclusive, so no contradiction.
- **D-02 (DECIDE & STATE) F2 wording.** Used "no review was passed the wave-end fact, so the sweep did not run for that wave" to make the cause explicit for a cold reader, keeping the REQ's "after the wave's last review already ran" trigger verbatim.
- **D-03 (DECIDE & STATE) Cite the coordinator-shape label from work.md** so the reader of Step 7 can find who writes the brief; the label string is unchanged so the contract test stays green.

## Discovered Tasks
- (impact-low) `skills/do-work/CHANGELOG.md` line 19 (0.305.76 entry for REQ-641 wave-end consistency check) says the orchestrator decides "last" from the archive and the members it set aside; history, not a contradiction. → report only
- (impact-low) Under delegated integration an integrator may itself set its own REQ aside; the coordinator then learns it from the hand-back and must add it to later briefs. Current prose implies this (the coordinator writes the brief from what it holds) but does not say it. → report only

## Lessons read
`_dev/primes/lessons-action-files.md` alternate-writer-contract-drift bullets (REQ-477, 498, 513, 461, 531, 566, 640, 641) and the restated-mechanism-unchecked bullet (REQ-639); prime Traps.

## Proposed lesson
`[family: alternate-writer-contract-drift] [REQ-642: a rule that names who holds a state ("the orchestrator knows which it set aside") goes stale once that role is delegated; name where the delegate reads it (the integrator brief) and say what happens when the state changes after its only reader already ran, instead of building a re-run]`

## Proposed CHANGELOG entry
**The Last Review in a Wave Now Knows Which REQs Were Set Aside.** The integrator brief lists the wave members the coordinator set aside, so a delegated integrator can tell its REQ is the last successful one and run the wave-end check. If a member is set aside after the last review already ran, the run's Decision Brief says the check did not run for that wave under HANDLED, and no review is re-run.

## Integration seams
None.

## Restatement-grep accounting
Commands (from worktree root): `grep -rnE 'set aside|set-aside' skills/`, `grep -rnoE '.{0,60}(last successful).{0,60}' skills/`, `... (wave membership)`, `... (integrate\.md|integrator brief)`, `grep -rnE 'wave-end|wave end' skills/`, `grep -rnE 'HANDLED' skills/`, `grep -rnc 'Decision Brief' skills/`, `grep -rn 'skip the check\|does not skip' skills/`.
- set aside: Go code/tests (finalization set-aside records, unrelated vocabulary); run-with-recovery.md 33/49/55, commit.md 51, work-reference.md 310/401/514/516, work-guide.md 130 (recovery set-aside, not the wave-end rule); work.md 35 (isolation, unrelated); deep-explore-reference.md (unrelated sense); work.md 370 and work-reference.md 470 = my edits. No change.
- last successful: review-work.md 135 (reviewer is passed the fact, never decides "last"; consistent); work.md 368 ("passed below", consistent, not restated); work.md 370 (edited); work-reference.md 456 (names the sweep only; consistent). No change.
- wave membership: work-reference.md 470 only (edited).
- integrate.md / integrator brief: work.md 170 (addendum text in brief; consistent); work-reference.md 430 (run-dir artifact staging; consistent); 470 (edited); 479 guardrail row "per-integrator input" (points at the paragraph; consistent). No change.
- wave-end / wave end: review-work.md 135, 138, 405, 491 (reviewer side, trigger set and skip rule; no set-aside wording; consistent); work.md 370 (edited). No change.
- HANDLED / Decision Brief: work-reference.md 493 (reader list, illustrative), 864 template, 871 bullet (D-01); work.md 536 hand-back pointer; clarify.md, anti-slop.md, review-work.md mentions are generic. No change.
- "does not skip": only CHANGELOG history. docs/work-guide.md 134 user sentence: consistent, outside write boundary, no change.

## Test wall times
- shipped-package-reference-contract.sh: PASS, 1.0 s
- contract-regressions.sh: PASS, 23.1 s

## Addendum

Review fix for F1 (HANDLED bullet could drop the late set-aside note).

- Commit: `aeb3f98f` on `worktree-agent-REQ-642-wave-end-set-aside-gaps` (after 22c2268e).
- File: `skills/do-work/actions/work-reference.md` (modified), Decision Brief HANDLED bullet only; rest of the bullet byte-identical, no heading or label changed.
- Changed first sentence: "**HANDLED** lists the **DECIDE & STATE** decisions (reversible `D-NN` entries), plus any run-level call the orchestrator states itself, such as a wave-end sweep that did not run because a member was set aside after the wave's last review (`actions/work.md` Step 7), so the user can spot-check without being asked to ratify."
- Changed omit condition: "Omit the block when the sections were read and held no DECIDE & STATE entry and the orchestrator stated no such call;"
- Tests: `shipped-package-reference-contract.sh` PASS, 1.0 s wall; `contract-regressions.sh` PASS, 20.9 s wall; `git diff --check` clean.
- **D-04 (DECIDE & STATE):** D-01 (HANDLED bullet needs no edit) is superseded by review F1. The bullet's omit rule made it exclusive, so a run whose REQs recorded no decisions would omit HANDLED and lose the late set-aside note. The bullet now names orchestrator-stated run-level calls and keeps the block while one exists.
