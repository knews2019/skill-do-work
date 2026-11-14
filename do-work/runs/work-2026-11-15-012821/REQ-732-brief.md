# Builder Brief — REQ-732 (allow docs/run-prompts/ in the staged-skills contract and move the run prompts there)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and writes the changelog.

- Operative name and branch: `worktree-agent-REQ-732-allow-docs-run-prompts`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-732-allow-docs-run-prompts`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its Triage, decisions D-01 to D-03 and the Red-Green Proof: read `do-work/working/REQ-732-allow-docs-run-prompts-folder.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-732-handback.md`
- Route A: no Exploration or Scope section. The change is small and fully named below.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-shell-commands.md` before you edit the `.sh` file, and `_dev/primes/prime-releases.md` (this REQ is no release; confirm it there).
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. Required lesson: the bullets carrying `[family: assertion-passes-when-feature-is-dead]` in `_dev/primes/lessons-shell-commands.md` (find them with a fixed-string grep). It is why the probe and your red-green evidence both carry a negative case: a scan that allows everything must not pass.

## Decided (do not reopen)
- **Maintainer, 2026-11-14:** run prompts live in `docs/run-prompts/`. Keeping `_dev/run-prompts/` with a `CLAUDE.md` line is out of scope.
- **D-01 — records keep their text.** Every remaining mention of `_dev/run-prompts/` sits in `do-work/` records (UR inputs, REQ `Source:` lines, the sibling REQs REQ-724 to REQ-731) or inside the moved prompts. Move the six prompts byte for byte so git records pure renames. Edit no `do-work/` path and no prompt text.
- **D-02 — the GREEN probe runs the contract's docs checks, not the whole contract** (the whole contract took 39 s under load). The probe cuts the segment from the line starting `for maintainer_docs_path in` to the line starting `board_files=(` out of the real script. Keep both marker lines as they are; if you must reword them, update nothing else and report it under Scope changes, because the probe then fails closed.
- **D-03 — scope stays at the request.** The `staged-skills` lane coverage in `_dev/tests/heavy-lanes.json` does not include root `docs/`, which is why the fast gate missed the 2026-11-14 break. Do not change it; list it under Discovered Tasks.

## What to build
One commit (the contract change and the move must land together, or the contract goes red in between).

1. `_dev/tests/staged-skills-contract.sh`:
   1. Line 75: add `docs/run-prompts` to the `maintainer_docs_path` list.
   2. Lines 104-107: add a `! -path "$repo_root/docs/run-prompts/*"` exclusion to the `find` beside the three `docs/` exclusions.
   3. Lines 71-74: the comment names the three folders; add run prompts in a few words (for example "and the maintainer's run prompts for other sessions (REQ-732)").
   Keep `docs` in `legacy_runtime_paths`: every other new file under root `docs/` must still fail.
2. `.gitattributes`:
   1. Add `/docs/run-prompts export-ignore` to the export-ignore block (lines 46-57), aligned like its neighbours.
   2. Add one short comment paragraph in the style of the `/docs/decisions` and `/docs/changelog` paragraphs: run prompts are dated prompts the maintainer writes for another session, never read at runtime.
   3. Lines 17-18: the sentence "asserts the /docs/kb, /docs/decisions and /docs/changelog lines exist" gains `/docs/run-prompts`.
3. Move the prompts: `git mv _dev/run-prompts/<each of the six .md files> docs/run-prompts/`, one path at a time or with the explicit list. Confirm `ls _dev/run-prompts` fails afterwards (git leaves no empty folder; remove it with `rmdir` if your tree still shows one).
4. Sweep: `git grep -n '_dev/run-prompts' -- . ':(exclude)do-work' ':(exclude)docs/run-prompts'` must print nothing. At base it already prints nothing outside the moved files, so expect no further edits.

**Not yours:** `VERSION`, `CHANGELOG.md`, anything under `do-work/` (even though the REQ's `write_set` names nothing there), `_dev/tests/heavy-lanes.json` (D-03), and everything under `skills/` and `suite/`.

## Rules
- Write only the files above. If the request's own requirements need another file, proceed and report it under Scope changes with the reason; otherwise stop and report.
- Out-of-scope finds go to Discovered Tasks in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`; the move with `git mv`). Never `git add -A` or `git commit -a`. The commit subject starts with `[REQ-732]`. End the commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Up to ten builders run at the same time in this run. Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. RED, before your commit: in the worktree, create `docs/run-prompts/probe.md`, run `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`, record the exit status and the FAIL line naming the probe file, then delete the probe file.
2. GREEN, after your commit, with `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` (about 40 s under load):
   1. With only the moved prompts present: exit 0 and `staged skills contract: PASS`.
   2. With an extra `docs/run-prompts/probe.md`: exit 0.
   3. With an extra `docs/other/probe.md` instead: non-zero, and the FAIL line names `docs/other/probe.md`.
   4. Delete both probe files and the `docs/other` folder; `git status --short` shows nothing.
3. `git check-attr export-ignore -- docs/run-prompts` prints `export-ignore: set`, and `git archive HEAD | tar -t` lists no `run-prompts` path.
4. `shellcheck _dev/tests/staged-skills-contract.sh` and `bash _dev/tests/quiet-grep-pipeline-audit.sh`.
5. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-732-probe.sh` run with your worktree as the current directory prints `REQ-732 probe PASS`. It fails at base by design (`the contract's docs checks do not name docs/run-prompts`). The pre-dispatch agent checked it on a simulated fix: it passes on the full fix, fails when the scan exclusion is missing, and fails when the scan allows every `docs/` file.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject.
- `## Files` — every file modified, with the verb; one line for the six-file move.
- `## Scope changes`
- `## Red-green evidence` — the contract failing with `docs/run-prompts/probe.md` at base, passing after your commit, still failing with `docs/other/probe.md`; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-04 onward (D-01 to D-03 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks` — at least D-03's lane coverage gap.
- `## Integration seams` — the line ranges you changed in the two shared files below, after your edit.
- `## Proposed CHANGELOG entry` — "none (no shipped file changes)": `_dev/`, `/docs/run-prompts` and `.gitattributes` are export-ignored and nothing under `skills/` or `suite/` changes. Confirm this against `.gitattributes` and `suite/modules.tsv` before you write it.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
All ten REQs of this run build at the same time from the same main HEAD; the integrator merges them one at a time.
| Shared file | REQ-732 edits | Sibling edits |
| --- | --- | --- |
| `_dev/tests/staged-skills-contract.sh` | the docs comment and loop at 71-81 and the `find` exclusions at 104-107 | none expected; a sibling that adds a shipped file list entry edits other lines |
| `.gitattributes` | the header sentence at 17-18, one new comment paragraph, one export-ignore line | none expected |
| `_dev/run-prompts/2026-11-15-move-global-rules-into-skill.md` | moved, text unchanged | REQ-724 (global rules into the skill) and REQ-725 (memory notes report) name this prompt as their source; they read it, they do not edit it |

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
