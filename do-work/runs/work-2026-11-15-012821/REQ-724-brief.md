# Builder Brief — REQ-724 (move the maintainer's global standing rules into the skill's crew-member files)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-724-move-standing-rules-into-skill`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-724-move-standing-rules-into-skill`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with decisions D-01 to D-06, Exploration (file:line pointers, the six RED searches) and Scope: read `do-work/working/REQ-724-move-global-standing-rules-into-skill.md` from the MAIN tree (your worktree may hold an older copy).
- The source text of the rules: the verbatim block in `do-work/user-requests/UR-175/input.md` (main tree). Never read, edit or delete any file under `~/.claude/`.
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-724-handback.md`
- Up to ten builders run at the same time in this run. No sibling edits `skills/do-work/crew-members/`.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `clear-questions.md`, `anti-slop.md`.
2. `_dev/primes/prime-action-files.md` (§ Cross-Referencing for the cross-package path) and `_dev/primes/prime-releases.md`.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, no test that pins wording).
4. Required lesson: `_dev/primes/lessons-releases.md`. Dropped for budget: `_dev/primes/lessons-action-files.md`; its matching family `[family: budgeted-context-routing]` does not apply (see Exploration).

## Decided (do not reopen)
- **D-01:** G5 (prompts for another session or tool) is one bullet in `communication-style.md`, citing `../../do-work-toolbox/actions/folder-structure.md` for the run-prompts folder. Write no folder path of your own.
- **D-02:** the G4 bug-report line goes to `communication-style.md`. The numeric-target line goes to `coding-guardrails.md` § 1 Think Before Coding, the redesign line to § 3 Surgical Changes.
- **D-03:** the status-word rule is a clause in the new G1 bullet. Do not reword the shipped technical-term line (`communication-style.md:17`).
- **D-04:** "decide reversible trivia yourself" is a pointer to the DECIDE & STATE tier (`coding-guardrails.md` § The decide-vs-escalate gate) plus "report it in one line". Do not restate the tier.
- **D-05:** G7 is one paragraph at the end of `background-agents.md` § Match the Pattern to the Harness, before `## Manifest Format`.
- **D-06:** you never edit `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or any version mirror. You propose the entry; the integrator releases.

## What to build
One commit, prose only, about 25 shipped lines in total (approximate; fewer is better). Write every new line in the house style of `communication-style.md` itself: simple words, about 20 words per sentence, no semicolons or em-dash chains, no vendor, tool, host or machine name. New Core Patterns bullets go at the END of the `## Core Patterns` list (after line 23). The Aliases and Examples sections do not change, and no existing sentence is reworded.

The GREEN probe checks these exact phrases (case-insensitive), so each must appear once, in the file named:
1. `communication-style.md`, G1 bullet: explaining the state of work or a system gives the "big picture in one short paragraph" first, then the exact picture, and defines "every status word" (for example claimed, blocked, pending) in plain words the first time. The next-command close is already shipped (line 12); do not repeat it.
2. `communication-style.md`, G2 bullet: a decision for the user goes through the agent's ask-user tool when it has one, never a "prose list of questions"; load `crew-members/clear-questions.md` first. Reversible choices: D-04.
3. `clear-questions.md`, principle 5: add one sentence. Put the "recommended option first" and mark it, and give each option its value and its risk.
4. `communication-style.md`, G3 bullet: commands the user will paste are written for the user's "login shell", and a payload for a remote host is wrapped in that host's shell. Fill in every value, so a pasted command holds no placeholder. Logic longer than about 5 lines goes into a committed script with a "one-line launcher".
5. `communication-style.md`, G4 line: "A bug report means" fix the code, not write a memory note.
6. `coding-guardrails.md`: numeric targets in a request are approximate "unless the user says strict" (§ 1), and never remove existing UI "while redesigning" unless the user asked (§ 3).
7. `communication-style.md`, G5 bullet: write the full prompt to a file in the project's run-prompts folder (cite `../../do-work-toolbox/actions/folder-structure.md`), commit it by exact path, and reply with a paste block of 3 lines or fewer. The prompt must run unattended: it parks questions as pending answers, includes a "20-minute stall check", and leaves operator steps out.
8. `background-agents.md`, G7 paragraph: when the harness lets you choose, pick each agent's model by its work. Lookups, searches and file sweeps run on the "cheapest model tier" that can do them. Planning, review and build agents keep the session model. Never leave a pure lookup on the default.

If the same idea is already shipped in a file you read, do not add it, and say so under `## Decisions`. If a phrase above would read badly, you may change the wording around it, but keep the quoted phrase so the probe and the REQ's GREEN proof still find it.

**Not yours:** `VERSION`, both `CHANGELOG.md` files and version mirrors (D-06); anything under `do-work/`; `skills/do-work/actions/restart-with-parallel-handoff.md`; `skills/do-work-toolbox/`; any new file (no new crew-member file, no test file).

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-724]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. RED, at your base before editing: run the GREEN probe below and record that it fails (12 FAIL lines expected at base `1389e2be` or a later main HEAD). Also record the six RED searches listed in the REQ's Exploration (all empty).
2. GREEN: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh` with your worktree as the current directory prints `REQ-724 probe PASS`.
3. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` (about 25-45 s) pass.
4. `git diff --stat <base>..HEAD`: only the four crew-member files, about 25 added lines in total, no deleted lines except where an added sentence joins an existing paragraph.
5. Re-read each added line against `crew-members/communication-style.md` § Patterns to Avoid.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file modified, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the six RED searches empty at base and the GREEN probe failing at base, then the probe passing at your tip with its PASS lines.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-07 onward (D-01 to D-06 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams`
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Standing Rules for Status Explanations, User Decisions, Pasted Commands and Subagent Models Ship in the Crew Files") and the line "What changes for an installed copy": every session in an installed project now follows the new always-on rules in `communication-style.md` (big picture first, status words defined, decisions through the ask-user tool, paste-ready commands, prompts for another session), and question wording, implementation guardrails and fan-out model choice gain one rule each. This is a patch-or-minor bump per `skills/do-work/actions/work-reference.md` → Changelog Entry Procedure; say which and why.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
No sibling in this run edits `skills/do-work/crew-members/`. REQ-732 (allow `docs/run-prompts/` in the staged-skills contract) changes this repository's own run-prompts folder, but your G5 bullet cites `folder-structure.md` and names no path, so the two do not touch. List each file you changed with the heading you added text under.

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.

## Coordinator note (2026-11-15, the maintainer's answer to D-05)

The maintainer was asked where the model-per-agent rule lives and chose `skills/do-work/crew-members/background-agents.md`, as D-05 planned. Do not add a copy or a pointer line in `communication-style.md`. Record this in your hand-back under `## Decisions` as confirmed by the maintainer.
