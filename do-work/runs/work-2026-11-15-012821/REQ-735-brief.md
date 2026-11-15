# Builder Brief — REQ-735 (decide which moved standing rules load in ordinary sessions, and ship it)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-735-moved-rules-always-on-scope`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-735-moved-rules-always-on-scope`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with decisions D-01 to D-07 (inside `## Exploration`), the load-site evidence and Scope: read `do-work/working/REQ-735-rule-scope-for-moved-standing-rules.md` from the MAIN tree (your worktree may hold an older copy).
- Parent work: REQ-724 (moving the maintainer's global standing rules into the skill, released 0.322.0), archived at `do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md` (main tree). Read its `## Decisions` (D-05, D-14) and `## Lessons Learned`.
- Never read, edit or delete any file under `~/.claude/`. Deleting the maintainer's global file is the maintainer's own step after this ships.
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-735-handback.md`
- This pass has two builders: you and REQ-736 (one line in `skills/do-work-toolbox/actions/help.md`). REQ-736 integrates first, then you.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `background-agents.md` (§ Match the Pattern to the Harness, the last paragraph), and `skills/do-work/actions/commit.md` § Rules.
2. `_dev/primes/prime-action-files.md` (§ Traps, `[family: alternate-writer-contract-drift]`) and `_dev/primes/prime-releases.md`.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, no test that pins wording).
4. Required lesson: `_dev/primes/lessons-releases.md`. Dropped for budget: `_dev/primes/lessons-action-files.md`; read its two matching families with a fixed-string grep: `[family: alternate-writer-contract-drift]` and `[family: condition-preserving-prose-extraction]`.

## Decided (do not reopen)
- **D-01:** R1 moves. The line `- Read a number in the request (a size, a count, a time) as approximate unless the user says strict.` is deleted from § 1 of BOTH `coding-guardrails.md` copies (do-work and toolbox) and added byte-identical to `communication-style.md` as a plain bullet right after the bug-report bullet (`- A bug report means a request to fix the code, not to write a memory note.`). No pointer stays in `coding-guardrails.md`.
- **D-02:** R2 moves the same way. The line `- Never remove existing UI while redesigning unless the user asked for it.` is deleted from § 3 of BOTH `coding-guardrails.md` copies and added byte-identical right after the R1 bullet. No pointer.
- **D-03:** R3 stays. The model paragraph stays only in the three `background-agents.md` copies. No copy, no pointer and no mention of it in `communication-style.md` (REQ-724 D-05, confirmed by the maintainer).
- **D-04:** F3 folds in. (a) In `communication-style.md`, the first sentence of the ask-tool bullet becomes: "Use the agent's ask-user tool when it has one, never a prose list of questions, and load `crew-members/clear-questions.md` first, unless the action in use says how to ask." (b) In all three `background-agents.md` copies, "Planning, review and build agents keep the session model." becomes "Planning, review and build agents keep the session model unless the action names another model." Keep the paragraph's first two lines unchanged and rewrap only from the third line, so the phrase "cheapest model tier" stays on one line. The three paragraphs must stay byte-identical.
- **D-05:** F5 folds in as one new bullet in `skills/do-work/actions/commit.md` § Rules, right after the `git add -A` bullet: `- **Never revert, stash or discard a change you did not make.**` plus at most one short reason sentence (for example, another session or the user may own it). Do not edit the `git add -A` bullet and do not add `git commit -a` to it.
- **D-06:** you never edit `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or any version mirror. You propose the entry; the integrator releases.
- **D-07:** the probe's fixed phrases are in the REQ's Exploration ("RED searches"). Keep each quoted phrase exactly as written above.

## What to build
One commit, prose only, seven files, about 10 lines added and 11 removed. The pre-dispatch agent ran exactly these edits on a scratch clone: `git diff --stat` showed 7 files, +10/-11, and the GREEN probe passed. Do not reword any other shipped sentence. The Aliases and Examples sections of `communication-style.md` do not change. No machine, vendor, model or tool name in any added text.

**Not yours:** `VERSION`, both `CHANGELOG.md` files and version mirrors (D-06); anything under `do-work/`; `skills/do-work-toolbox/actions/tutorial.md` and `skills/do-work/actions/ai-evaluation-loop-reference.md` (cited, not changed); `crew-members/ui-design.md`; the toolbox `coding-guardrails.md` "Earned defense" drift (pre-existing, out of scope); any new file or test file.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-735]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Up to ten builders may run at the same time in this run. Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. RED, at your base before editing: run the GREEN probe below and record that it fails. At base `547d4cc5` it prints 7 FAIL lines (R1, R2, F3-ask, F5, the two moved lines, F3-model).
2. GREEN: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-735-probe.sh` with your worktree as the current directory prints `REQ-735 probe PASS`. The probe checks every package copy: the two `coding-guardrails.md` copies lose both lines, and the three model paragraphs stay identical (REQ-724 review F1 was exactly this miss).
3. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` (about 35 s) pass. This is the pre-flight baseline the pre-dispatch agent ran green at `547d4cc5`.
4. `git diff --stat <base>..HEAD`: exactly the seven Scope files.
5. Re-read each added or changed line against `crew-members/communication-style.md` § Patterns to Avoid.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file modified, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the GREEN probe failing at base (7 FAIL lines) and passing at your tip with its PASS lines.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-08 onward (D-01 to D-07 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams`
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Numeric-Target and Redesign Rules Apply in Every Session; Actions Can Override the Ask and Model Rules") and the line "What changes for an installed copy": every session in an installed project now reads numbers in a request as approximate and keeps existing UI during a redesign (these rules move from the implementation-only guardrails into the always-on communication rules); an action's own way of asking and its named models now win over the general ask-tool and model rules; the commit action gains a rule never to revert, stash or discard a change you did not make. Propose minor (impact-rule-change, every installed project's sessions change) per `skills/do-work/actions/work-reference.md` → Changelog Entry Procedure, and say why.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
None with REQ-736: it edits only `skills/do-work-toolbox/actions/help.md`. List each file you changed with the heading you edited under (for example `communication-style.md` § Core Patterns, `coding-guardrails.md` § 1 and § 3, `background-agents.md` § Match the Pattern to the Harness, `commit.md` § Rules).

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
