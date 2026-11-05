# REQ-710 implementation handback

Implemented the optional AI-Evaluation-Loop within all eight declared source paths. Ready for parent integration and final gates. No live cross-provider approval is claimed.

- Branch: `t3/req710-skill-review-modes`
- Full commit: `ad4ef86bb15d6d101a5abe02d4bb386668a4a186`
- Starting commit: `db0e8b1e984399da66c7bde24a50569459abee2c`
- Subject: `[REQ-710] Add optional AI-Evaluation-Loop modes`
- Checkout: `/Users/t2/.t3/worktrees/skill-do-work2/t3-d84d79ab`
- Source-only commit: eight paths, 149 insertions, 19 deletions. No deletions of files.
- Parent is sole integrator. No PR, release metadata, queue state, or upstream installation was created. This handback is the only main-tree write.

## Implementation Summary

- Modified `skills/do-work/SKILL.md`: discoverable run argument and companion routing.
- Modified `skills/do-work/actions/work.md`: parse/consume action-owned options before CLI advance; persist per-REQ policy before dispatch; select the companion at review; qualify ordinary score/failure fallbacks; retain heavy-test obligations and counters; validate the selected outcome before canonical finalization; expose unresolved results in reporting.
- Added `skills/do-work/actions/ai-evaluation-loop-reference.md`: parent-owned five-step policy for provider discovery and execution evidence, durable selection/counters/ledger, complete candidate identity, recoverable checkpoint, independent read-only rounds, five-total-round/two-stall decisions, advisory jester, and honest lifecycle result.
- Modified `skills/do-work/actions/review-work.md`: caller-owned evidence-only review with sanitized request snapshots, proposed writes returned to parent, material-finding approval rules distinct from ordinary scores, and explicit persisted selected outcome. Ordinary and standalone behavior retained.
- Modified `skills/do-work/actions/work-reference.md`: canonical pointer, recovery persistence, and unresolved-evaluation versus terminal-success reporting contract.
- Modified `skills/do-work/actions/help.md`: installed invocation and availability/stop guidance.
- Modified `skills/do-work/docs/work-guide.md`: usable Sol/council/jester examples, provider requirements, ledger/counters/checkpoint behavior, and explicit dependency semantics; qualified ordinary review restatements.
- Modified `skills/do-work/docs/review-work-guide.md`: selected-mode exception to scoring/writes, invocation, identities, stopping and unresolved outcomes, with links to canonical policy.

## Checks

All completed checks exited 0:

1. `bash _dev/tests/shipped-package-reference-contract.sh` — PASS. Also passed within the final contract suite, with the new companion staged and visible.
2. `bash _dev/tests/contract-regressions.sh` — PASS on the first implementation and again after the final edge-case edits. Final output: `Contract regression checks passed.` All reported individual fast files remained below 30 seconds. Final log: `/tmp/req710-contract-regressions.log` on this machine (temporary evidence, not a committed artifact).
3. `git diff --cached --check` — PASS before commit.
4. Reviewed every changed file and the full new companion against the authoritative acceptance criteria. Swept shipped Markdown for review-skip, score, one-remediation and completed-with-issues restatements. Remaining ordinary clauses are explicitly conditioned on ordinary review.

The aggregate's staged-skills/updater/installer heavy probes were explicitly skipped by its normal fast-tier behavior. Full `maintainer-verify.sh` and heavy lanes were not run; the parent owns those. No deterministic helper/schema was introduced, and no decorative prose-token tests were added.

Source inspection supporting the design (read-only, not live workflow execution): `requeststate/state_apply.go` uses `generatedRecoveryHeading` to strip specific generated sections; a separate top-level `AI-Evaluation-Loop` section is outside that set. `schemanormalization/schema_normalization.go` keeps `completed-with-issues` in `IsTerminalSuccess`, and dependency regression fixtures explicitly cover that status. No lifecycle semantics were changed.

## Simulated workflow walkthroughs

These are manual, repeatable instruction walkthroughs against the shipped prose, not automated execution of an LLM policy and not live model reviews. To repeat one, start with its stated REQ/body/candidate inputs, supply the described synthetic responses, and walk companion Steps 1–5 plus work Step 8. `C0`, `C1`, etc. denote distinct complete candidate manifests; they are symbolic IDs, not measured hashes. `SIM-*` provider/model/task records are synthetic fixtures, not actual harness evidence. An installed run must replace them with verified real execution metadata.

### S1 — Ordinary behavior and standalone review

Input: `do-work run REQ-042`, no saved loop section. Ordinary Acceptance Pass/Overall 80 follows the existing completed branch. Ordinary Partial/60 also follows its existing completed/report routing. Ordinary Fail gets exactly one remediation; continued failure chooses completed-with-issues. An ordinary reviewer failure still records the skipped review and continues. No loop record, counters or provider requirement is introduced. `do-work review-work REQ-042` still uses standalone review writes and commit behavior, even when archive context includes a prior loop result.

Result: ordinary and standalone branches are unchanged. Selecting a loop changes only that REQ's review contract.

### S2 — Selection reaches the action, not CLI advance

Input: `do-work run REQ-042 --review-mode sol-loop --jester`. Parse selection first; retain Sol/jester in action state, pass only the existing targeting/selector arguments to advance, and persist selection at committed claim before dispatch. Repeat with `--review-mode council`, `--fan-out 2`, or a UR target; existing selector semantics remain unchanged. `--jester` alone, unknown/missing mode values or duplicate mode selections fail before queue mutation. A fresh invocation omitting flags reloads a saved selection instead of downgrading it.

Result: no unknown option reaches the deterministic CLI and no selection is lost in its continuation.

### S3 — Clean approval, Sol and council

Input: C0 satisfies requirements; covering checks pass. Supply resolved/dispatched/returned synthetic identity records for one OpenAI Sol seat, e.g. `SIM-OpenAI/SIM-Sol/task-S1`. Reserve round 1. The reviewer and parent explicitly approve C0, no material findings remain, and identity is unchanged before finalization.

Result: AI-Evaluation-Loop approved at total 1; no second round. If no new verification beyond the baseline occurred, stalls is 1, but approval has priority. Repeat for council with three separate verified synthetic OpenAI/Anthropic/Google execution records and all three approving C0: approved at round 1. Two of three approvals is insufficient. A clean approval at round 5 is also permitted; round 6 is never dispatched.

### S4 — Five total rounds with progress and a remaining material defect

Input: accepted material finding F1 remains open throughout, with no dispute about its validity. Each pass has all required replies. Round 1 exercises a previously unassessed required check successfully. Between rounds 1/2, 2/3, 3/4 and 4/5, fix and verify separate accepted findings F2, F3, F4 and F5 without regressions. Candidates are C0 through C4. All reviewers keep requesting changes for F1.

Trace: `(total, stalls)` = `(1,0), (2,0), (3,0), (4,0), (5,0)`. F1 retains its stable ID and open evidence across all passes; F2–F5 keep fixed-and-verified records.

Result: stop at five total rounds, unresolved / not approved, F1 unmet. Retain best verified C4 with its open concern. Step 8 chooses completed-with-issues when existing gates permit. No sixth pass or counter reset after fixes. The resulting lifecycle still satisfies depends_on and UR completion; it is not consensus.

### S5 — Two stalled rounds despite rising scores

Input: C0 has accepted material F1. Round 1 returns no new checked improvement; round 2 merely repeats advice and raises a subjective score from 70 to 95. No behavior/check coverage improves and F1 remains open.

Trace: `(1,1)` then `(2,2)`.

Result: stop unresolved after round 2. The higher score cannot approve or reset stalls. No unnecessary fixes are invented to spend remaining rounds.

### S6 — Progress resets stalls only; recovery never resets totals

Input: selected council record; round 1 has no supported improvement `(1,1)`. Round 2 verifies an accepted fix `(2,0)`. Interrupt and recover/requeue, stripping ordinary Review/Testing but preserving the separate top-level loop section. Resume without options and with the same candidate/evidence.

Trace: restored mode is council and counters remain `(2,0)`. Round 3 stalls `(3,1)`; round 4 verifies another fix `(4,0)`; round 5 stalls `(5,1)` with F1 still open.

Result: stop at five total, not a fresh five-round allowance. If the crash instead occurs after reserving round 3, that round stays used; recover its existing task/results or report missing required evidence unresolved without re-dispatch. Conflicting mode selection or inconsistent/missing known loop evidence never initializes new zero counters.

### S7 — Stale approval, including dirty and untracked bytes

Input: parent and all required seats approve C0 at commit H. Change an unstaged implementation hunk while HEAD remains H; alternatively add a relevant untracked file. The complete manifest becomes C1. Recompute identity before finalization.

Result: C0 approvals cannot approve C1. Requalify/check C1 and use a fresh required round if budget permits; otherwise unresolved. The same rule covers integration edits. A log-only edit in a previously recorded non-behavior exclusion does not manufacture a source change. Unknown/omitted implementation bytes cannot be hidden under that exclusion.

### S8 — Recoverable best checkpoint and unrelated user changes

Input: after round 2, retain C1 containing REQ-owned `feature.txt` bytes `good` plus owned untracked `new-feature.txt`, its modes/manifest and check evidence. Record unrelated `notes.txt` bytes and an unrelated staged user hunk. Round 3's C2 regresses the feature. The best checkpoint remains C1, not the higher-scoring C2.

Replay restore: compare ownership/current bytes; restore only the owned C1 files/hunks, preserve `notes.txt` and the unrelated index entry, recompute identity and rerun the recorded checks. Total remains 3. Round 4 must obtain fresh parent plus required approvals for restored C1 before approval is possible. Repeat with total 5: restored C1 remains unresolved because no review budget remains. Repeat with a later user edit overlapping the owned hunk: refuse that restore, preserve both versions and report the conflict, never reset/clean the tree.

Result: best verified work is recoverable and named, without treating recovery or preservation as approval. This is a simulated ownership/restore walkthrough; no user files were actually restored during implementation.

### S9 — Missing/failed required reviewers and identity mismatch

Input variants: Sol discovery has no actual Sol; council lacks Google; a required task fails; a task produces no reply; a reply lacks execution identity; actual execution identity differs from the resolved explicit model. For post-dispatch cases, round 1 was already reserved.

Result: unresolved / not approved with the missing seat/failed task/identity evidence named. Pre-dispatch failure uses zero rounds; a dispatched failed round remains used. The parent cannot substitute itself or a jester. Step 8 and Error Handling cannot fall through to the ordinary skip-review or high-score completed branch.

### S10 — Absent and nonvoting jester

Input: council plus --jester; all required seats and parent approve C0. In separate runs, the optional Gemini cannot resolve, fails, is slow at the required-pass decision boundary, or the harness cannot run it nonblocking.

Result: record advisory absence and finish the required decision without waiting. Approval depends on required evidence only. Repeat with a timely advisory score of 1 and no supported defect: same approval. Repeat with score 100 while a required seat requests changes: still not approved. Advisory scores never enter council grades in a PR description. A missing required Google seat remains missing regardless of advisory availability.

### S11 — Supported advisory defect

Input: the harness fixture resolves `latest available Gemini` to `SIM-Gemini-latest` and returns matching synthetic invocation metadata. Its neutral brief is identical in requirements/candidate/evidence to other reviewers and contains no “jester,” advisory-role assignment or humor request. Its response identifies an input-handling defect F6 with a reproducible failing case against C0, despite a score of 100.

Result: parent reproduces F6 and adds it as a material accepted finding with origin/evidence. Approval is prevented by F6, not its score. Apply an in-scope fix to C1, rerun the failing check, then obtain all required approvals in the next budgeted round. Unverifiable allegations are declined with reasons; late advisory output is reported for validation, never silently rewriting the finished result. These model labels are simulation data, not claims about currently available Gemini releases.

### S12 — Persistent disagreement and scope control

Input: required reviewer disputes material F7; parent returns counterevidence. A second consecutive reply addressing that evidence remains in material disagreement. Alternative: F7 is material but its fix requires undeclared scope.

Result: preserve evidence/disposition instead of deleting the finding or averaging scores. Persistent disagreement stops unresolved. Out-of-scope material concerns cannot approve or authorize scope growth; only existing critical impact routing can create automatic follow-up work. Noncritical lines stay report-only.

## Decisions

- D1 — Kept all behavior in the skill/action prose. Candidate/checkpoint mechanisms use existing Git/file facilities; no runtime helper, lifecycle schema or status was needed.
- D2 — Stored selection/history in an independent top-level AI-Evaluation-Loop body section, because ordinary generated Review evidence is resettable. Reservations happen before dispatch so recovery cannot gain extra rounds.
- D3 — Bound approval to complete candidate bytes and explicit required verdicts, with a fresh pass after restoration. Checkpoint ownership rules preserve unrelated staged, unstaged and untracked user work.
- D4 — Reused completed-with-issues honestly: unresolved selected evaluation is not approved, while dependency and UR semantics remain unchanged. Selected outcomes cannot waive existing test/repository/heavy gates.
- D5 — Council uses one independent seat from each of OpenAI, Anthropic and Google; Sol resolves GPT Sol through the harness. Jester is a separate advisory invocation, even when council already has a Google seat. Neutral snapshots keep role assignments with the parent.

## Lessons and guidance read

Read the authoritative brief in full, repository CLAUDE.md, `_dev/primes/prime-action-files.md` and the complete `_dev/primes/lessons-action-files.md` satellite, plus `_dev/primes/prime-shell-commands.md` and `_dev/primes/prime-releases.md`. Read the general, coding-guardrails, shared-principles, communication-style and anti-slop crews.

Applied action lessons about downstream catchall closure, alternate writers/restatements, preserving extracted conditions, installed citation depth and checking the implementation before restating recovery semantics. No required-lessons entries were declared; no required guidance/lesson path was missing. Initial exploratory guessed CLI directories were absent; the actual recovery and status consumers were located under internal/requeststate and internal/schemanormalization. No source was changed there.

## Discovered Tasks

None requiring separate capture. No genuine scope expansion was needed.

## Limitations and parent next step

Live model discovery, cross-provider dispatch, actual returned identities, nonblocking advisory behavior and real consensus convergence are untested by instruction. The scenarios above demonstrate the intended workflow by simulation, not live approval. The implementation remains an agent-followed contract rather than a deterministic CLI enforcement engine.

Parent: integrate the source commit, consume this evidence, then run the final full gate and any selected heavy lanes under your ownership. Do not interpret the simulated approvals as approval of this implementation by external reviewers.

## Parent remediation validation

Independent review found that accepted fixes reused the first handback. Parent fixed the companion and work Step 6 (source 0fadcd9a, integrated 74b3de5a5bd90968441dfd5c47f8d73a0baf3653). These additional scenarios are instruction walkthroughs, not live model runs:

- S13 first accepted fix: initial attempt A has handback H-A and merge receipt M-A; accepted fix creates attempt B with absent H-B, updates active run row, and only then enters Step 6. H-A does not prevent a new dispatch. Original review base, total/stall counters and ledger survive.
- S14 recovery during B: restored active attempt points to existing invocation; recover it without new attempt or duplicate dispatch. Landed unconsumed H-B integrates once.
- S15 interruption after merge before receipt: source commit ancestry proves B is integrated; recover merge receipt and continue qualification/testing without empty re-merge. Unprovable integration ends unresolved.
- S16 another fix: consumed B remains in history; new C gets a distinct absent handback before dispatch. Fix attempts do not allocate or reset review rounds; only review dispatch consumes the next reserved round.

Parent traced these against work Step 6 and the fan-out merge/recovery rules. Remediation stays within the original eight-file scope.
