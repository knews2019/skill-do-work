# Builder Brief — REQ-736 (the ai-report judge line in the toolbox help shows the --page option)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases. This REQ changes one line of one shipped file.

- Operative name and branch: `worktree-agent-REQ-736-help-judge-page-option`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-736-help-judge-page-option`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its Triage and decision D-01: read `do-work/working/REQ-736-help-judge-line-shows-page-option.md` from the MAIN tree (your worktree may hold an older copy). Its UR: `do-work/user-requests/UR-179/input.md` (item 2).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-736-handback.md`
- Route A, no pre-flight. One sibling builder (REQ-735, scope of the moved standing rules) runs at the same time; it touches no file of yours. REQ-736 integrates first.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`.
2. `_dev/primes/prime-action-files.md` and `_dev/primes/prime-releases.md` (the REQ's prime files).
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add).
4. Required lessons: none added. Dropped only for budget: `_dev/primes/lessons-action-files.md` (8761 tokens). Read only its `[family: restated-mechanism-unchecked]` bullets (fixed-string grep): before shipped prose restates a command form, read the source that defines it.
5. The source of truth for the form: `skills/do-work-toolbox/actions/ai-report.md:36` (`judge <bundle-dir> [--page <file.html>]`) and `skills/do-work-toolbox/docs/ai-report-guide.md:80`.

## Decided (do not reopen)
- **D-01 — overflow shape:** the new usage `ai-report judge <dir> [--page <file.html>]` is longer than the 30-character usage column. Keep the description on the same line after two spaces, the shape line 20 (`ai-report --kind proposal|root-cause <target>  Decision brief…`) already uses. Do not realign or shorten any other line.

## What to build
One commit. In `skills/do-work-toolbox/actions/help.md` line 18 (inside the `text` fence), change

`  ai-report judge <dir>          Render-check a report bundle (widths, themes, links)`

to

`  ai-report judge <dir> [--page <file.html>]  Render-check a report bundle (widths, themes, links)`

You may extend the description to say the option checks one named page (for example "…, or one named page") if it stays short; record it as a D-02 DECIDE & STATE if you do. No other line of `help.md` changes.

**Not yours:** `VERSION`, both `CHANGELOG.md` files, the version line in `skills/do-work/actions/version.md` (the integrator writes the release); anything under `do-work/`; `ai-report.md` and `ai-report-guide.md` (already correct).

## Rules
- Write only `skills/do-work-toolbox/actions/help.md`. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add skills/do-work-toolbox/actions/help.md`). Never `git add -A` or `git commit -a`. Commit subject starts with `[REQ-736]`. End the commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-736-probe.sh` run with your worktree as the current directory prints `PASS: REQ-736 help.md judge line shows --page`. It fails at base by design (`FAIL: help.md judge line does not show [--page <file.html>]`); record both runs.
2. `git diff --stat HEAD~1` shows one file, one line changed (1 insertion, 1 deletion).
3. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/staged-skills-contract.sh` exit 0.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject.
- `## Files` — the one modified file.
- `## Scope changes` — "none" or the reason.
- `## Red-green evidence` — the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-02 onward (D-01 is taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — none expected; say so, or name any shared file you touched.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Toolbox Help Shows the ai-report judge --page Option") and the line "What changes for an installed copy": `do-work-toolbox help` lists `ai-report judge <dir> [--page <file.html>]`, so the single-page check added in 0.319.0 is visible from the menu. Help text only; no behaviour changes.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
| Shared file | REQ-736 edits | Sibling |
| --- | --- | --- |
| `skills/do-work-toolbox/actions/help.md` | line 18 only | none (REQ-735 edits only `crew-members/` files) |

Your final message: under 1500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
