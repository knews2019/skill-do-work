# REQ-641 hand-back (wave-end consistency check)

- Branch: worktree-agent-REQ-641-wave-end-consistency-check
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-641-wave-end-consistency-check
- Base: 2291812d. Commits: 39b21c30 (first build), 0b306f69 (coordinator answers applied). Review the range 2291812d..0b306f69.

## File manifest
- skills/do-work/actions/review-work.md (modified): Restatement Sweep label parenthetical widened; step 1 gains one wave-end condition (fires only when the orchestrator passes the last-integration fact and earlier members' REQ ids; reads their archived `**Restatement sweep:**` lines via work-reference.md **Folder Structure**; a missing line is named unread for the sweep, never re-derived); step 4 skips only when the whole trigger set is empty; Append to REQ File template gains the `**Restatement sweep:**` line (own elements only, plus `unread for the sweep: REQ-NNN`); Verification Checklist item points at that line.
- skills/do-work/actions/work.md (modified): Step 7 **Restatement sweep (MUST)** adds "or this REQ is its wave's last successful integration (passed below)"; **How to run it** has the orchestrator decide it from the run manifest (every other wave member finalized or set aside) and pass the fact plus the members' REQ ids; one clause says a run with no manifest has no wave.
- skills/do-work/actions/work-reference.md (modified): one sentence in the Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick". No heading changed.

## AI Execution State (P-A-U)
- [x] **[PLAN]:** Read the brief, REQ, prime-action-files.md, lessons-action-files.md bullets (lines 52-62), the delegated-integration paragraph and manifest row (work-reference.md 470, 480) and Step 7 of work.md. Approach: add the wave-end condition as a continuation of Restatement Sweep step 1 (keyed on "last successful integration of its wave", membership from the run manifest or integrator brief), widen step 4's skip to the whole trigger set, add a record step and template line holding only this diff's own redefinitions, repoint the checklist, align work.md Step 7's trigger restatement by citation, pass the manifest to the spawned reviewer, and add one sentence to work-reference.md line 456. Routing untouched.
- [x] **[APPLY]:** First build in 39b21c30. The coordinator then answered the four open choices. 0b306f69 applies them: the reviewer never scans the manifest, the orchestrator decides and passes the fact (work.md Step 7), missing lines become "unread for the sweep", the archive location cites **Folder Structure**, and the separate record step 5 was folded into the template line to keep the addition small.
- [x] **[UNIFY]:** `git diff --stat 2291812d..0b306f69`: review-work.md 8 (+5/-3), work-reference.md 2 (+1/-1), work.md 4 (+2/-2); 3 files, 8 insertions, 6 deletions. `git diff --check` clean. Cold read: review-work.md step 1 names the input (orchestrator-passed fact and ids), where reviews live (Folder Structure, a real `##` heading), what a missing line means, and the reason sentence; the template line says what it holds; work.md Step 7 says how the orchestrator decides; work-reference.md cites **Restatement Sweep**, a real bold label. No debug artifacts.

## Decisions
- **D-01 DECIDE & STATE (coordinator answer 1): who decides "last".** The orchestrator or integrator decides it from the run manifest, where every other row of the wave is finalized or set aside, and passes the fact plus the sibling REQ ids to the reviewer like the merge range. The reviewer never scans the manifest. Value: a set-aside member cannot skip the check, and the reviewer needs no manifest parsing. Risk: an orchestrator that forgets to pass it skips the check silently; the MUST paragraph in work.md Step 7 names the case.
- **D-02 DECIDE & STATE (answer 2): no manifest, no wave.** One clause in work.md Step 7: "a run with no manifest has no wave". Value: a condition, not a mode list. Risk: none.
- **D-03 DECIDE & STATE (answer 3): missing earlier line.** The member is named on this review's line as `unread for the sweep: REQ-NNN`; its elements are never re-derived from its diff, which would be a second review. Value: absence stays visible. Risk: waves straddling this release show unread members once.
- **D-04 DECIDE & STATE (answer 4): where earlier reviews are.** Archived REQ files, found by the passed ids, cited as `actions/work-reference.md` → **Folder Structure**, with no restated paths. Value: no second copy of the archive layout. Risk: none.
- **D-05 DECIDE & STATE: record-line content.** This diff's own elements only, never inherited ones, stated inside the template placeholder. Value: a later member reads it as what this member redefined. Risk: none.
- **D-06 DECIDE & STATE: label parenthetical widened.** "(when the diff, or an earlier member of its wave, redefines ...)" so the check's own label does not restate the old trigger. Value: sweep closed. Risk: a citation pinning the full label; the contract test passes and no full-label citation exists.
- **D-07 DECIDE & STATE: step 4 kept.** Required by the brief so a last-of-wave review whose own diff redefines nothing does not skip. Value: correctness. Risk: none.

## Discovered Tasks
None (the work-guide sentence was done in the Addendum).

## Lessons read
- `_dev/primes/lessons-action-files.md` (via prime-action-files.md Lessons list): family alternate-writer-contract-drift bullets (REQ-477, 498, 513, 461, 531, 566, 640) and REQ-639 restated-mechanism-unchecked. No new mechanism claim added: "Go checks section presence only" was verified by the orchestrator and is not restated in the shipped text.

## Proposed lesson entry (_dev/primes/lessons-action-files.md)
- [family: alternate-writer-contract-drift] [REQ-641: widening a trigger also stales the bold label that names the check, because the label's parenthetical restates the trigger; sweep the check's own heading and skip rule, not only other files, and record per-REQ facts a later reviewer must inherit as a fixed line in the durable review block, with the orchestrator, not the reviewer, deciding when to read them](../../do-work/archive/REQ-641-wave-end-consistency-check.md#lessons-learned)

## Proposed CHANGELOG entry
**The Last Review in a Parallel Wave Now Checks What Earlier REQs Redefined**
Each review now records, on one `Restatement sweep:` line, which contracts its own change redefined. When a REQ is the last one of a parallel wave to integrate, its review also sweeps every element the earlier members recorded, because their reviews ran before the later work merged and could not see a later REQ restating the old meaning. A git merge only catches line collisions; this catches meaning collisions. A set-aside member does not skip the check, and an earlier review with no line is named as unread for the sweep. Findings route as before.

## Integration seams
None. Last REQ of the UR-138 chain.

## Tests (from worktree root, at 0b306f69)
| Command | Result | Wall |
|---|---|---|
| bash _dev/tests/shipped-package-reference-contract.sh | exit 0, PASS | 1s |
| bash _dev/tests/contract-regressions.sh | exit 0, "Contract regression checks passed." | 20s |
(At 39b21c30: 1s and 23s, both pass.)

## Restatement grep accounting
Command (zsh, explicit files; the `grep` alias is ugrep, so `/usr/bin/grep` was used):
`/usr/bin/grep -n "<pattern>" skills/do-work/actions/review-work.md skills/do-work/actions/work.md skills/do-work/actions/work-reference.md skills/do-work/docs/review-work-guide.md skills/do-work/docs/work-guide.md skills/do-work/crew-members/background-agents.md`
- `estatement [Ss]weep`: review-work.md 130, 135, 405, 491 changed (the step-5 hit is gone). work.md 368 changed. work.md 205 (plan-time counterpart, cites by name, no trigger restated) stays. work-reference.md 456 changed.
- `redefin`: review-work.md 130, 134-138, 405, 491 changed or new; 134 (own-diff question) stays as the base trigger; 136 ("Sweep each redefined element") applies to the whole set, stays. work.md 368, work-reference.md 456 changed.
- `nothing redefined|sweep N/A`: old "sweep N/A" wording gone; all hits are the new text.
- `non-interference`: work-reference.md 456 changed; 462 (Auto-wave, points back at 456) stays; 112 (write_set schema comment) and work.md 550 (write_set advisory) are about the pick, stay.
- `line proximity|not meaning`: only work-reference.md 456, changed.
- `wave membership|last successful`: review-work.md 135, work.md 368, work-reference.md 456 new; work-reference.md 470 (integrator brief carries wave membership) stays and is cited by D-02.
- `run manifest`: review-work.md no longer names it (the orchestrator passes the fact); work.md 370 new; work.md 169, 173 (mid-run messages), 536, work-reference.md 848 unrelated, stay.
- `Append to REQ File`: review-work.md 491 (checklist) new, 381 heading unchanged, 226 unrelated; work.md 393 (defers to the template) stays.
- `report only`: routing unchanged everywhere; review-work.md 137 (step 3) is the route new findings use; docs/review-work-guide.md 55 and docs/work-guide.md 55 stay.
- background-agents.md: no hits for any pattern.

## Addendum
- Commit: a09d4e87 (after 0b306f69). Review range now 2291812d..a09d4e87.
- File manifest addition: skills/do-work/docs/work-guide.md (modified): one sentence in the "Building several REQs at once" paragraph, after "nowhere else.": a merge only catches edits to the same lines, so the review of the last REQ to land in a wave also checks whether a later REQ restated something an earlier one changed the meaning of, and reports what it finds.
- **D-09 DECIDE & STATE: scope extended by the coordinator from the discovered task.** The coordinator added skills/do-work/docs/work-guide.md to Scope and write_set, so the user guide now states the wave-end check. Value: users learn that meaning collisions are checked, not only line collisions. Risk: none; one sentence, no heading, no list.
- Test: `bash _dev/tests/shipped-package-reference-contract.sh` at a09d4e87: exit 0, PASS, 1s. `git diff --check` clean.
