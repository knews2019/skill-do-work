# REQ-675 hand-back (core review names four delivery stages)

- Branch: worktree-agent-REQ-675-review-delivery-stages
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-675-review-delivery-stages
- Base: e313e870
- Commit: f12beed1 `[REQ-675] review names four delivery stages and marks unexercised ones unassessed` (one commit)

## File manifest
- skills/do-work/actions/review-work.md (modified): Step 7 gains a **Delivery stages** block (four one-sentence stage definitions, "which apply is a judgment, many REQs have no deployment stage"), and one paragraph after the Score block saying the result scores only exercised stages, the one-line summary names the stages covered, and Step 8 lists the rest. Step 8 gains a first category, **Unassessed delivery stages**, required whenever an applicable stage was not exercised, named by stage and marked unassessed.
- skills/do-work/docs/review-work-guide.md (modified): Phase 3 gains one paragraph saying the same in user words.

## P-A-U
- [PLAN]: Read the brief, REQ, UR-151, general.md, coding-guardrails.md, shared-principles.md, communication-style.md, anti-slop.md, prime-action-files.md, prime-releases.md, lessons-releases.md (whole) and the alternate-writer-contract-drift bullets of lessons-action-files.md. Plan: define the stages once in Step 7 between "What NOT to do" and "If you can't run the code"; put the "scores only exercised stages" rule after the Score block without touching the Pass/Partial/Fail/Untested bullets; put the "unassessed" entry rule in Step 8 as a category so the rule lives in one place; one paragraph in guide Phase 3.
- [APPLY]: Done as planned. +9 lines in review-work.md, +2 in the guide. No other file touched.
- [UNIFY]: `git diff --stat e313e870`: review-work.md +9, review-work-guide.md +2, 2 files, 11 insertions, 0 deletions. Checks from the worktree root:
  - REQ-675-probe.sh: exit 0, 1 s.
  - _dev/tests/shipped-package-reference-contract.sh: exit 0, 1 s.
  - _dev/tests/contract-regressions.sh: exit 1, 24 s. The only failure is the quiet-grep pipeline audit on `do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh` and `REQ-677-probe.sh` (run artifacts committed in e313e870, not in my diff). Unrelated to this change; see Discovered Tasks.
  - `git diff --check`: exit 0.
  - `git diff e313e870 -- skills/do-work/actions/review-work.md`: two hunks, @@ -185 and @@ -193, both inside Step 7 and Step 8 (Step 7 starts at line 160, Scoring Guidelines at 212 before the edit). Append to REQ File template (`**Acceptance:** [Pass/Partial/Fail/Untested] — [1-line summary]`), score table, Scoring Guidelines, Verdict mapping and Step 10 untouched.
  - Files checked: both changed files read in full around the edits; no em-dashes in the added prose.

## Restatement sweep
- skills/do-work/crew-members/shared-principles.md:14-15: "Acceptance is inferred from narrower tests → exercise end-to-end where applicable" and "Acceptance cannot be exercised → Record Untested". Still agree: Untested stays the result when no stage could be exercised, and the new text only scopes Pass/Partial/Fail to the exercised stages. No edit.
- skills/do-work/actions/work.md Step 7 (:366-368) and failure table (:518): route on Pass / Partial / Fail and the overall score. Unchanged values, unchanged routing. Still agree. No edit.
- grep over skills/ for `Untested`, `Suggested Additional Testing`, `delivery stage`, `live acceptance`: no other reader of the review's Acceptance result (toolbox inspect/code-review hits use their own scales).

## Behavioral exercise
Method: reading the edited Step 7 and Step 8 text against the scenario myself (no reviewer agent spawned). Scenario: scratch REQ "publish the updated rules page", diff correct, local tests pass, no access to the served site.
- GREEN (edited text): applicable stages are all four (it is a publish). Implementation and integration were exercised; deployment and live acceptance were not. Report lines produced:
  - `**Acceptance:** Pass — implementation and integration: the diff updates the rules page as asked and the local test suite passes.`
  - Suggested Additional Testing: `Deployment: unassessed. Confirm the published build reached the serving host.` and `Live acceptance: unassessed. Open the served rules page and confirm the updated rules show.`
- RED (base text e313e870): the code runs locally, so "If you can't run the code" does not fire and Untested does not apply; Pass reads "Feature works end-to-end as specified". Produced: `**Acceptance:** Pass — rules page updated, tests pass.` with no stage named, and Step 8 categories are all optional, so nothing forces a deployment or live line.
- Not exercised: a real reviewer agent run on a scratch REQ file. Reason: the change is prose only and the brief allows the reading method; a spawned run would add cost without testing more than the same text.

## Decisions
- D-01 DECIDE & STATE: the "unassessed" entry rule lives in Step 8 as the first category (required whenever such a stage exists), and Step 7 only points to it. One home for the rule, and it sits where the reviewer writes the entries.
- D-02 DECIDE & STATE: the Pass/Partial/Fail/Untested bullets are unchanged. The scoping sentence follows the Score block instead, so the score values every reader routes on keep their exact wording.
- D-03 DECIDE & STATE: used the latitude on the one-line summary: Step 7 tells the reviewer to name the stages covered there. The Append template itself is unchanged.
- D-04 DECIDE & STATE: the new Step 8 bullet uses a colon, not the em-dash its siblings use, per the no-em-dash rule. Small style mismatch inside the list, accepted.
- D-05 DECIDE & STATE: the existing Step 8 category "Integration testing" (third-party services, APIs) shares a word with the new "Integration" stage. Left as is: renaming it is outside the write boundary's intent. Noted for the reviewer.

## Discovered Tasks
- impact-important: `_dev/tests/contract-regressions.sh` fails on base e313e870 because `do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh` and `REQ-677-probe.sh` decide on a quiet grep fed from a pipeline (quiet-grep pipeline audit). These are orchestrator run artifacts; rewrite the two probes with a here-string (as REQ-675-probe.sh does) before the gate runs, or the canonical gate stays red for every member of this run → report only

## Lessons read
- _dev/primes/lessons-releases.md (whole): families canonical-link-outlives-its-target, manifest-ownership-vs-edit-content. Not triggered (no links or manifests touched).
- _dev/primes/lessons-action-files.md, only `[family: alternate-writer-contract-drift]` bullets (REQ-477, 498, 513, 461, 531, 566, 640, 641, 642, 644, 648, 647, 652). Applied as the restatement sweep above.

## Proposed CHANGELOG entry
**Review says which delivery stages it checked.** A review "Acceptance: Pass" from local runs could be read as proof that the served build works. The review now names the stages it covered and lists the rest as unassessed.
- `review-work.md` Step 7 defines four delivery stages once: implementation, integration, deployment and live acceptance. Which apply is a judgment.
- The Acceptance result scores only the stages the review exercised, and its one-line summary names them.
- Step 8 lists every applicable stage the review did not exercise as unassessed, by stage name.
- `docs/review-work-guide.md` Phase 3 says the same in user words. The persisted Review block, scoring and verdict mapping are unchanged.

## Proposed lesson bullet
none

## Integration seams
None found. REQ-676 and REQ-677 touch toolbox files only. REQ-678 (release-check) will cite these stage names: implementation, integration, deployment, live acceptance (headed **Delivery stages** in review-work.md Step 7).

## Test wall times
probe 1 s, shipped-package-reference-contract 1 s, contract-regressions 24 s.
