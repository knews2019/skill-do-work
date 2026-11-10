# Builder Brief — REQ-717 (suite-managed justfile recipes call a script in one line)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The orchestrator merges, tests, reviews and releases.

- Worktree (the ONLY tree you write, plus the hand-back file below): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-717-justfile-recipes-call-scripts`
- Branch: `worktree-agent-REQ-717-justfile-recipes-call-scripts`. The coordinator creates the worktree from main HEAD right before dispatch. Do not create, move or remove a worktree yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its Triage, Exploration, Decisions D-01 to D-08 and Scope: read `do-work/working/REQ-717-suite-justfile-recipes-one-line.md` from the MAIN tree. Your worktree may hold an older copy.
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/2026-11-09-layout-queue/REQ-717-handback.md`
- Proof tools (main tree, run them by absolute path from inside your worktree; each checks the tree it runs in): `do-work/runs/2026-11-09-layout-queue/REQ-717-probe.sh` (GREEN probe, red at base), `REQ-717-live-proof.sh` (real servers in a throwaway consumer; green at base and after), `REQ-717-preflight-probe.sh` (baseline).

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-shell-commands.md` and `_dev/primes/prime-kanban-board.md`. From `_dev/primes/lessons-shell-commands.md` read the bullets tagged `check-reads-the-file-it-names` and `legacy-fixture-implementation-shape`, and the section "Every Flag on a Shipped Script Needs a Non-Test Caller".
3. `skills/do-work/tools/lessons-do-work-update.md` (the request's one required lesson).
4. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).

## What to build
Behavior-preserving move of recipe logic into two shipped scripts. The request's Decisions are settled; follow them.

- **Scripts** (D-01, D-02): `skills/do-work-board/tools/run-kanban.sh <project-root> <port>` and `skills/do-work-board/tools/kanban-static.sh <project-root>`. `#!/usr/bin/env bash`, a header comment naming the recipe it serves, mode 100755. Each finds `queue-kanban/` beside itself through `$(dirname "${BASH_SOURCE[0]}")`. No flags, no defaults.
- **Move the lines verbatim** (D-04) from `skills/do-work-board/justfile.template:6-8` and `:16-17`: every message, exit code, the 320-step wait and `serve --open`. Replace the `{{justfile_directory()}}` interpolations with the project-root argument and the `cd .../tools/queue-kanban` with the script-relative directory. Split the long one-liners into readable lines if you like, but the logic must not change. No `exec`, no new output.
- **Shell options** (D-03), the trap that would break both scripts silently:
  - `run-kanban.sh`: `set -u` only. Under `pipefail` and `-e`, `listener_pid="$(lsof -ti ... | head -n1)"` fails when nothing listens, and the script dies before serving.
  - `kanban-static.sh`: a failed build or generate must stop the script before the exclude write. `set -e` does not stop on a failure inside an `&&` list unless it is the last command. Use separate lines, or end the chain with `|| exit`, and say why in a comment.
- **Recipes** (D-05): the template's two bodies become one line each, `bash "{{justfile_directory()}}/.claude/skills/do-work-board/tools/<name>.sh" "{{justfile_directory()}}"` (run-kanban adds `"$port"`). In the root `justfile` use `bash skills/do-work-board/tools/<name>.sh "{{justfile_directory()}}"` (plus `"$port"`). Keep the recipe headers and comments unchanged. Touch no other recipe.
- **`_dev/tests/contracts/queue-kanban.sh`** (D-06): delete `extract_kanban_shutdown_line` and the root/template parity check. Build a fixture board folder in a temp dir: copy the real script into `<fixture>/tools/`, put a stub executable at `<fixture>/tools/queue-kanban/queue-kanban` that records its argv (and, for `generate --out DIR`, writes `DIR/index.html`). Stub `go`, `lsof`, `ps`, `kill`, `sleep` as functions and `export -f` them (`kill` is a builtin, so a PATH stub cannot replace it). The macOS `/bin/bash` is 3.2, so test with it. Cases, each with its own assertion message:
  1. Stale queue-kanban listener: the `lsof` stub reports pid 4242 until `kill` writes a marker, then nothing. Expect exit 0, the line `queue-kanban: stopping previous session on :8090 (pid 4242): ...`, `kill` called with 4242, and serve argv `serve --open --repo-root <root> --port 8090`. This case is new.
  2. A stuck listener is refused with exit 1, naming pid 4242 and the command. Keep the existing case.
  3. A foreign listener is refused with exit 1, and `kill` is never called. Keep the existing case.
  4. kanban-static, run twice in a temp `git init` project: `build/queue-kanban-board/index.html` exists, the exclude file holds exactly one `/build/queue-kanban-board/` line, and the announcement line appears on the first run only. Match it as a whole line with `grep -x`, because just echoes recipe text. This case is new.
  5. kanban-static with a failing `go` stub: non-zero exit and no exclude line. This case is new.
  Keep the REQ-633 stamp check at the end of the file as it is. Use here-strings for every quiet grep, never `| grep -q`.
- **`_dev/tests/staged-skills-contract.sh:928-991`** (D-06): replace the count of `.claude/skills/do-work-board/tools/queue-kanban` (4) with a count of `.claude/skills/do-work-board/tools/` (still 4: two queue-kanban recipes plus two script calls). Require both scripts with `require_file`. Read the four safety strings (`[ "$wait_count" -lt 320 ]`, `remaining_listener_pid=`, `listener_executable_name=`, `serve --open`) from `skills/do-work-board/tools/run-kanban.sh`, the file the template calls, and no longer from the template. Edit only this block. Lines 85-125 (`legacy_runtime_paths`, `retained_bootstrap_tool`) belong to REQ-716.
- **Go lock-in** (D-06): add `TestShippedBoardTemplateRecipesHoldNoLogic` to `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go`. It copies the shipped template (`filepath.Join("..", "..", "..", "..", "..", "do-work-board", "justfile.template")`, the same relative read `suiteinstall/install_transaction_test.go:489` uses) into a temp dir as `justfile`, calls `justfileLogicGaps`, and expects no gap. Its comment names the failure it pins: every consumer's `folder-structure-check` reported `FOLDER-STRUCTURE-JUST-LOGIC` for the suite's own recipes. The GREEN probe runs this exact test name.
- **`skills/do-work-board/SKILL.md`** (D-07): add one sentence to the Ownership boundary list naming `tools/run-kanban.sh` and `tools/kanban-static.sh` as the logic of the two managed recipes. Change no other prose.
- **Not yours**: `VERSION`, both `CHANGELOG.md` files and `skills/do-work/actions/version.md` (the finalizer writes the release); anything under `do-work/` except the hand-back file; the files in the Scope's "will NOT touch" list.

## Rules
- Two commits. Commit 1 holds the test changes only (contract, staged-skills block, Go test). Run them and record the failure output. Commit 2 holds the scripts, the two justfiles and SKILL.md, and makes the tests pass.
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason. Otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Subjects start with `[REQ-717]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not run `recover` or any `advance` form.
- Another builder (REQ-716, moving the five root installer scripts from `tools/` to `cmd/`) runs at the same time in its own worktree. Never touch its files. Never rename anything under `skills/do-work-board/tools/`.

## Verify before you hand back (all from your worktree)
Never run the repository gate `_dev/tests/maintainer-verify.sh`, because the integrator runs it. If a 30-second per-file budget fails once under load, rerun it once and record both runs.
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/toolboxcommands`, `go vet ./internal/toolboxcommands/`, `go test -count=1 ./internal/toolboxcommands/`.
2. `bash _dev/tests/contracts/queue-kanban.sh`, `bash _dev/tests/contract-regressions.sh`, `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` (about 30 s), `bash _dev/tests/prescribed-shell-canonicalization.sh`, `bash _dev/tests/quiet-grep-pipeline-audit.sh` (it scans tracked files, so run it after commit 2), and `shellcheck --severity=warning` on both new scripts and both changed test scripts.
3. `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/2026-11-09-layout-queue/REQ-717-probe.sh` prints `REQ-717 probe passed.`
4. `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/2026-11-09-layout-queue/REQ-717-live-proof.sh` prints all three PASS lines. The third must say `kanban-static calls skills/do-work-board/tools/kanban-static.sh`, not "still inline". It starts real servers on ports 18717 to 18757 and stops them on exit. Afterwards check that `lsof -ti tcp:18717 -sTCP:LISTEN` prints nothing.
5. `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text folder-structure-check` lists no `FOLDER-STRUCTURE-JUST-LOGIC` (the root-entry gaps remain; they are recorded exceptions or REQ-716's).
6. `just --justfile skills/do-work-board/justfile.template --summary` lists the same 41 names as before.

## Hand-back file content
Use these headings:
- `## Commits`: the hash and subject of each commit.
- `## Files`: every file created, modified or deleted, with the verb.
- `## Scope changes`
- `## Red-green evidence`: test names, the failing output before, and the passing run after.
- `## Verification`: each command above with its result.
- `## P-A-U`: three short lines, PLAN, APPLY and UNIFY. The UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions`: D-09 onward, each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams`: see below. Confirm or correct it.
- `## Proposed CHANGELOG entry`: a descriptive title and the line "what changes for an installed copy". For example, the run-kanban and kanban-static recipes now call `.claude/skills/do-work-board/tools/*.sh`, and `do-work update` rewrites the managed block. Behavior and messages are unchanged.
- `## Proposed lesson`: one bullet, or "none".
- `## Lessons read`

## Integration seams (REQ-716 builds in parallel; integration is serial, REQ-716 first)
- `_dev/tests/staged-skills-contract.sh`: you edit only the board-template block (about lines 928-991). REQ-716 edits `legacy_runtime_paths` and the `retained_bootstrap_tool` list (about lines 85-125). These are separate hunks and should merge cleanly.
- The root `justfile`: you edit only the `run-kanban` and `kanban-static` bodies. REQ-716 changes the root justfile only if a recipe names the root `tools/`, and today none does.
- `skills/do-work-board/tools/`: REQ-716 moves the root `tools/` only. If its sweep rewrites `tools/<script>` spellings, it must not match `skills/do-work-board/tools/`. Your new names (`run-kanban.sh`, `kanban-static.sh`) differ from its five scripts.
- `_dev/tests/heavy-lanes.json` and `heavy_maintainer_tree_test.go`: you change neither. The template change already selects the staged-skills, updater and installer heavy lanes for the integrator.

Your final message: under 2500 characters. State the branch tip hash, the probe and live-proof results, and anything the orchestrator must decide.
