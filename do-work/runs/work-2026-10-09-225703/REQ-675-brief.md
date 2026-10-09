# Builder brief — REQ-675 (Core review names four delivery stages)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-675-review-delivery-stages
- Branch: worktree-agent-REQ-675-review-delivery-stages, created from main HEAD BASE_PENDING with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-675-review-delivery-stages.md. Read it fully: What, Why, Finding Provenance, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof, Required Lessons — Dropped for Budget, and the orchestrator's `## Triage`. The release requirement belongs to the integrator; you PROPOSE the changelog entry and a lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-151/input.md (the maintainer's brief and the validate-feedback triage that produced this REQ).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-09-225703/REQ-675-handback.md
- Route A, tdd: false, impact-user-visible, effort-mechanical, domain general. Prose only, two files.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md and anti-slop.md (load it before writing any prose that ships) (same directory). Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (template, earned sections, cross-reference spelling, conditions not lists, package justification) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release itself is the integrator's). Lessons: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file, small). The capture dropped `_dev/primes/lessons-action-files.md` for budget; read only its bullets carrying `[family: alternate-writer-contract-drift]`, because this change edits a rule other files restate.

## The change (decided at capture, do not reopen)
In `skills/do-work/actions/review-work.md` and `skills/do-work/docs/review-work-guide.md` of your worktree:
1. Step 7: Acceptance Testing (`review-work.md:160-195`): define the four delivery stages once, each in one plain sentence: implementation (the diff does what the REQ asks), integration (it works in the merged tree with the rest of the system, for example the test suite and adjacent flows), deployment (the built or packaged result reached its serving environment or consumer install), live acceptance (the consumer-visible behavior is right in the real environment). Say that which stages apply is a judgment and that many REQs have no deployment stage.
2. State that the Acceptance result scores only the stages the review exercised, and that every applicable stage the review did not exercise is listed in Step 8's Suggested Additional Testing as "unassessed", by stage name. The natural places are near the "If you can't run the code" paragraph (`:188`) and the Score block (`:190-194`) in Step 7, and the Step 8 categories (`:196-210`).
3. `review-work-guide.md` Phase 3 (`:26-34`, ends before `## Scoring` at `:36`): one line saying the same in user words.
Constraints: do not change the `## Review` Append to REQ File template (the `**Acceptance:** [Pass/Partial/Fail/Untested] — [1-line summary]` line stays as is; the one-line summary may name the stages covered, that is your latitude), the score table, the scoring formula and caps, the Verdict mapping, the impact tokens, Step 10 routing, or any status. No retrospective mode. Do not restate the targeted-counterexample rule (`:113`, `:174`, `:181`, `:478`). A few sentences in total.
Restatement Sweep (do it yourself before hand-back; the reviewer repeats it): `skills/do-work/crew-members/shared-principles.md:14-15` (the two Acceptance rows) and `skills/do-work/actions/work.md` Step 7 (`:366-368`, the Acceptance routing) and its failure table (`:518`) read the Acceptance result; confirm they still agree with your wording and say so in the hand-back. Do not edit them unless they now contradict (then stop and ESCALATE in Decisions instead of widening scope).

## Write boundary
Exactly `skills/do-work/actions/review-work.md` and `skills/do-work/docs/review-work-guide.md`. Anything else: stop and say so in the hand-back.

## Hard rules
- Every commit subject on your branch starts with `[REQ-675]` (for example `[REQ-675] add ...`). The board credits commits by that prefix only.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (which you write but never stage or commit). The integrator's queue guard refuses a branch that touches `do-work/`.
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, or `just do-work-update`.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or any version mirror: the release belongs to the integrator.
- Scratch fixtures go in `mktemp -d` directories outside both trees and are never committed. Clean up background servers you start.
- Retired-trigger trap: `_dev/tests/staged-skills-contract.sh` scans every file under `skills/` for `do-work <retired word>` (words in column 4 of `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, for example `scan`, `present`, `deliver`, `suggest`, `learn`, `prompt`, `note`, `inspect`). Always write `do-work-toolbox <action>`; never `do-work <toolbox action>`.

Integration seams: none expected. REQ-676 and REQ-677 (new toolbox actions) build in parallel but touch none of your files.

## Behavioral exercise (one-off, from the REQ's Red-Green Proof)
Scenario: a reviewer following your edited `review-work.md` reviews a scratch REQ "publish the updated rules page" whose diff is correct and whose local tests pass, with no access to the served site. Record it by reading the edited Step 7 and Step 8 text against that scenario and writing the report lines it yields: the Acceptance line (which stages it covers: implementation and integration) and the Suggested Additional Testing entries (deployment and live acceptance, each marked unassessed by name). If you can spawn a reviewer agent, you may have it follow the edited text on a scratch REQ file in a temp directory instead; either way, quote the lines produced and state which method you used. Also record the RED side: the same scenario against the base text produces "Acceptance: Pass" with no stage named.

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-09-225703/REQ-675-probe.sh` run from your worktree root (it reads relative paths, so it checks your tree): exit 0.
- `bash _dev/tests/shipped-package-reference-contract.sh` (about 1 s) and `bash _dev/tests/contract-regressions.sh` (about 19 s): exit 0.
- `git diff BASE_PENDING -- skills/do-work/actions/review-work.md`: every hunk is inside Step 7 or Step 8; none touches the Append to REQ File template, Scoring Guidelines, Verdict mapping or Step 10.
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (one commit preferred). Base: BASE_PENDING.
- File manifest: each file with (new)/(modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Behavioral exercise record: what you ran (fixtures, commands, tool used), what the action produced (paste the key output lines), and every part you could not exercise with the reason. Honest partials beat claimed passes.
- `## Decisions` (D-01 onwards; each DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellites and families).
- A proposed CHANGELOG entry: a descriptive title that says what shipped, in plain words, one or two sentences on why it matters, then specific bullets. The integrator writes it with the version.
- A proposed lesson bullet, naming the satellite you think it belongs in (`_dev/primes/lessons-action-files.md` or `_dev/primes/lessons-releases.md`), in that file's bullet shape (`- [family: <slug>] [REQ-675: <one-line lesson>](<relative archive link; the integrator fixes the path>)`). Say "none" if nothing was learned that a later builder would need.
- Integration seams (see the note above), test wall times.
