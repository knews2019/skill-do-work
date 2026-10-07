# REQ-640 hand-back (mid-run messages steer a run in progress)

- Branch: `worktree-agent-REQ-640-mid-run-messages-steer-a-run-in-progress` (no name collision)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-640-mid-run-messages-steer-a-run-in-progress`
- Base: `31a6e35f`. Commits: `d7e2e186` (`[REQ-640] mid-run messages steer a run in progress`)

## File manifest

- `skills/do-work/actions/work.md` (modified): new `### Mid-Run Messages (any step)` between Step 3.5 and Mechanical Evidence-Gate Loop; Orchestrator Checklist line after the Step 3.5 line; Route C plan validation "Requirement coverage" also reads any `## Addendum` section.
- `skills/do-work/actions/review-work.md` (modified): Step 2 "What was requested" and Step 5 item 1 include any `## Addendum` section (Step 5 names `## Addendum (mid-run)`).
- `skills/do-work/actions/capture.md` (modified): Step 2 table `do-work/working/` row keeps NEVER modify and adds the pointer to `actions/work.md` → **Mid-Run Messages (any step)**.
- `skills/do-work/actions/work-reference.md` (modified): Decision Brief intro reads the run manifest's hand-back emphasis note and orders what it shows by it; guardrail table `manifest.md` row names that note. No heading changed.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read prime-action-files.md, lessons-action-files.md lines 50-61 (alternate-writer-contract-drift x6, restated-mechanism-unchecked), general.md, coding-guardrails.md, shared-principles.md, communication-style.md. Re-checked the mechanism claims in Go: `generatedRecoveryHeading` (requeststate/state_apply.go:981) does not list Addendum, so takeover keeps it; `advanceSections` (lifecycleadvance/advance_commands.go:356) refuses any duplicated section name, hence one section with entries inside. Approach: one condition-keyed subsection (branches keyed on what the message changes and whether the builder handed back), one checklist line, addendum added to every requirement-source reader, capture pointer in house arrow form, Decision Brief and manifest row made to name the emphasis note. Keep all work-reference headings.
- [x] **[APPLY]:** Edits exactly as planned in the four write-set files; the sweep added one phrase to work.md plan validation (D-05). Nothing outside the write set touched.
- [x] **[UNIFY]:** `git diff --stat 31a6e35f..d7e2e186`: capture.md 2 +-, review-work.md 4 ++--, work-reference.md 4 ++--, work.md 13 +++-; 4 files, 17 insertions, 6 deletions. `git diff --check` clean. `bash _dev/tests/shipped-package-reference-contract.sh` PASS (1s); `bash _dev/tests/contract-regressions.sh` PASS (21s; first run 24s). Cold-read of each changed section: the subsection states conditions, every citation resolves (Mid-Run Messages heading, Delegated integration bold label, Immutability Rule heading, Decision Brief heading, Outside-text containment bold label), no forward reference to text that does not exist. No debug artifacts.

## Decisions

- **D-01** DECIDE & STATE. No run manifest (coordinator note 1): keyed on "the run has a run directory". With one, the emphasis note goes in its `manifest.md`. Without one, the session that received the message renders the Decision Brief, holds the note itself and says so in one progress line; the text forbids creating a run directory or any file for this. Value: no new machinery. Risk: a crash in a serial run without a run directory loses the emphasis note; low cost, it only orders the hand-back.
- **D-02** DECIDE & STATE. Running integrator (coordinator note 2): written as part of the "builder already handed back" branch, which says it includes a REQ whose integrator is running: capture's in-flight path (new REQ with `addendum_to`), run by the coordinator in its next writing gap, never while an integrator runs, the words kept in session scratch outside the project root until then. Value: no third path, one-writer rule intact. Risk: the user's late change waits a loop; reversible prose.
- **D-03** DECIDE & STATE. Second message for one REQ (coordinator note 3): "A REQ carries one such section, with entries inside", each entry under its own timestamp line, with the reason (`advance` refuses a duplicated section name), stated as a condition, not a Go file citation. Value: no classifier refusal. Risk: none found.
- **D-04** DECIDE & STATE. Not-yet-handed-back branch under delegated integration also says the coordinator keeps the words in session scratch until its next writing gap, then puts them in the integrator brief. Value: a cold reader does not write the brief while another integrator runs. Risk: none.
- **D-05** DECIDE & STATE. Sweep: work.md Route C plan validation "Requirement coverage" now reads "What/Detailed Requirements and any `## Addendum` section". It restates the same requirement-source list as review-work Step 5 (family alternate-writer-contract-drift). Value: plan validation sees queued and mid-run addenda. Risk: none; one phrase.
- **D-06** DECIDE & STATE. review-work.md Step 2 line 55 changed the same way as Step 5 item 1 (it restates the source list). Value: the two review readers agree. Risk: none.
- **D-07** DECIDE & STATE. capture.md **Immutability Rule** left unchanged. It is capture's own rule and stays true for capture; the mid-run section is written by the work pipeline, which work-reference.md line 51 already exempts. The pointer sits in the Step 2 table, the point where capture meets a `working/` match. Value: delete-before-add, one pointer not two. Risk: a capture reader in a live run who stops at the Immutability Rule makes an addendum REQ instead; that path is still valid, only later.
- **D-08** DECIDE & STATE. Checklist line placed right after the Step 3.5 line, matching the subsection's position.

## Discovered Tasks

None open. The work-guide sentence found earlier is done, see Addendum.

## Restatement-grep accounting

Command (from worktree root): `grep -n "<pattern>" skills/do-work/actions/{work,work-reference,review-work,capture,clarify}.md skills/do-work/docs/work-guide.md skills/do-work/crew-members/background-agents.md`

- `mid-run\|mid run`: work.md 155 and 493 (escalated-answer rule, different rule, consistent); work.md 170 new; review-work 87 changed; work-reference 25, 165, 788 (unrelated state/blocked-flip); work-reference 470 already carries "any mid-run addendum text", consistent; clarify 99/259 unrelated; background-agents 20/36 unrelated; work-guide 72/101 consistent (D-discovered task).
- `## Addendum\|Addendum section\|addendum_to`: work.md 135 (addendum REQs at Triage, consistent), 170/171 new; capture 63, 123-125, 132-166 queued/in-flight capture paths, consistent and cited; work-reference 106, 670, 773, 779, 794 frontmatter/cascade, unrelated; review-work 55/87 changed, 347 follow-up template unrelated.
- `NEVER modify\|immutable\|Immutability Rule`: capture 61/63 kept (D-07), 124 changed, 125/164 unchanged; work-reference 51 consistent (work pipeline exempt).
- `integrator brief\|REQ-NNN-integrate.md\|mid-run addendum text`: work-reference 470/479 consistent; work-reference 430 unrelated (run artifacts settle order).
- `run manifest\|manifest.md`: work-reference 480 changed; 424/430 unrelated; background-agents 51/192 generic manifest, consistent (note is optional).
- `Decision Brief`: work-reference 846-848 changed (intro clause); work.md 536 Progress Reporting consistent; review-work 224, work-reference 493/694, clarify 38-42 consumers of the format, unaffected (the note is optional).
- `What/Detailed Requirements`: work.md 202 changed (D-05); review-work 55/87 changed.
- `paraphrase`: work.md 167 new; clarify 112 unrelated.

## Lessons read

`_dev/primes/lessons-action-files.md` lines 50-61: the six `alternate-writer-contract-drift` bullets and the REQ-639 `restated-mechanism-unchecked` bullet. Prime `_dev/primes/prime-action-files.md` in full.

## Proposed lesson entry (`_dev/primes/lessons-action-files.md`)

- [family: alternate-writer-contract-drift] [REQ-640: adding a new requirement source (`## Addendum (mid-run)`) is a reader-side contract change; the REQ named review Step 5 only, but review Step 2 and Route C plan validation restate the same "What/Detailed Requirements" source list and would have built and judged without the user's mid-run words. Grep the source-list phrase, not only the cited line](../../do-work/archive/REQ-640-mid-run-messages-steer-a-run-in-progress.md#lessons-learned)

## Proposed CHANGELOG entry

**Messages Sent During a Run Now Reach the Right REQ Without Stopping It.** `actions/work.md` gains a Mid-Run Messages rule. A question about progress is answered from the REQ files and run manifest. A change to a REQ still being built is written into that REQ as one `## Addendum (mid-run)` section in the user's own words, and review now checks the build against it. A change to work already handed back, queued work, or new work goes through capture. A note about what to show first goes into the run manifest, and the end-of-run Decision Brief reads it. Before this, such a message became a next-run addendum REQ and the builder never saw it.

## Integration seams

None. REQ-641 (wave-end consistency check) touches the same files next.

## Tests

- `bash _dev/tests/shipped-package-reference-contract.sh`: PASS, 1s
- `bash _dev/tests/contract-regressions.sh`: PASS, 21s (24s on the first run)
- `git diff --check`: clean

## Addendum

- New commit: `92f662d7` (`[REQ-640] work guide: you can keep talking while a run is in progress`). Branch commits are now `d7e2e186` and `92f662d7`.
- File manifest addition: `skills/do-work/docs/work-guide.md` (modified). One new paragraph after the coordinator sentences under **Building several REQs at once**. It says the user can keep talking during a run and the run never stops for it. A question is answered from disk. A change to work still being built reaches that REQ as an addendum in the user's own words. A change to work that already came back is queued as a follow-up REQ. No heading, no list.
- **D-09** DECIDE & STATE. The coordinator extended Scope and write_set with `skills/do-work/docs/work-guide.md` to close the discovered task. Value: users learn that steering mid-run works. Risk: none; two sentences, reversible.
- Test: `bash _dev/tests/shipped-package-reference-contract.sh` PASS, 1s. `git diff --check` clean.
