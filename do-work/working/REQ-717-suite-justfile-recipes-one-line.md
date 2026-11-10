---
id: REQ-717
title: 'Suite-managed justfile recipes call a script in one line'
status: claimed
created_at: 2026-11-09T21:47:36Z
user_request: UR-169
domain: general
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-shell-commands.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-11-09T23:48:22Z
  basis:
    - Route B
    - 8-file write set
    - 3 subsystems involved
    - 5 acceptance criteria
depends_on: [REQ-713]
required_lessons: [skills/do-work/tools/lessons-do-work-update.md]
claimed_at: 2026-11-09T23:39:43Z
route: B
write_set: [skills/do-work-board/tools/run-kanban.sh, skills/do-work-board/tools/kanban-static.sh, skills/do-work-board/justfile.template, justfile, skills/do-work-board/SKILL.md, _dev/tests/contracts/queue-kanban.sh, _dev/tests/staged-skills-contract.sh, skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go]
---
# Suite-managed justfile recipes call a script in one line
## What
Move the logic out of the suite-managed justfile recipes in `skills/do-work-board/justfile.template` so each managed recipe is one line that calls a script shipped with do-work-board. Today `run-kanban` has 3 body lines and `kanban-static` has 2, so every consumer repository gets this layout gap.
## Detailed Requirements
- `run-kanban`: move the port validation, the stale queue-kanban listener stop, and the build-and-serve into a shipped script (for example under `skills/do-work-board/tools/`). Keep the exact behavior and messages.
- `kanban-static`: move the build-and-generate and the `.git/info/exclude` write into a shipped script. Keep the exact behavior and messages.
- The installer rewrites the managed block in consumer justfiles, so also update this repository's root justfile.
- Depends on REQ-713 (folder-structure-check sees recipes with default parameter values); without it, the check cannot prove `run-kanban` is fixed.
- Layout rule this applies: a justfile recipe holds no logic and calls a script in one line (rule 5).
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository after REQ-713 lands.
**Why RED now:** It reports `FOLDER-STRUCTURE-JUST-LOGIC` for `justfile: kanban-static` and `justfile: run-kanban`.
**GREEN when:** No `FOLDER-STRUCTURE-JUST-LOGIC` gap is reported, `just run-kanban` still replaces a stale instance on its port, `just kanban-static` still writes the `build/` snapshot and the exclude line, and the board contract tests pass.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` — 6122 tokens, `slugged: partial`; matched because the run-kanban and kanban-static recipes drive the board.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, `slugged: partial`; matched because the recipe logic becomes shipped shell.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: Move the logic out of the suite-managed justfile recipes so each managed recipe is one line (layout rule 5: a justfile recipe holds no logic and calls a script in one line). In skills/do-work-board/justfile.template, run-kanban has 3 body lines (port validation, stop a stale queue-kanban listener, build and serve) and kanban-static has 2 (build and generate, then the .git/info/exclude write). Move that logic into scripts shipped with do-work-board (for example under skills/do-work-board/tools/), keep the exact behavior and messages, and make each recipe a single call. The installer rewrites the managed block in consumer justfiles, so also update this repo's root justfile. Depends on the REQ that fixes folder-structure-check for recipes with default parameter values. Verify: folder-structure-check reports no FOLDER-STRUCTURE-JUST-LOGIC on this repo, just run-kanban still replaces a stale instance on its port, just kanban-static still writes the build/ snapshot and the exclude line, and the board contract tests pass.*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and behavior-preserving: two recipes become one-line calls to two shipped scripts. The work is in finding what pins the old recipe text: a launcher contract evaluates the inline shutdown line, the staged-skills contract greps the template for safety strings and a path count, and the root justfile mirrors the template with repository-relative paths.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Required-lessons consult: `do-work/lessons-index.md` adds `skills/do-work/tools/lessons-do-work-update.md` (748 tokens; the managed Just section is an updater-managed path and the root justfile mirrors the template, a source/install parity case). Both captured candidates stay dropped: each is `slugged: partial` and over the 2000-token budget alone. Read for this exploration anyway, by family: `lessons-shell-commands.md` `check-reads-the-file-it-names` (REQ-555) and `legacy-fixture-implementation-shape` (REQ-420). `lessons-kanban-board.md` holds nothing on recipes or launchers.

**Where the logic lives today**
- `skills/do-work-board/justfile.template:5-8` `run-kanban $port="8090":` 3 body lines (port `case`, the `lsof` stale-listener stop with a 320-step wait, `cd .../tools/queue-kanban && go build && ./queue-kanban serve --open`). `:15-17` `kanban-static:` 2 lines (build + generate into `build/queue-kanban-board`, then the `.git/info/exclude` append guarded by `git check-ignore`).
- `justfile:5-8` and `:15-17` are the same recipes with repository-relative `cd skills/do-work-board/tools/queue-kanban` (this repository has no `.claude/skills` install). The rest of the root block differs from the template on purpose (core recipes fall back to `skills/do-work`), so the root block is hand-mirrored, not generated.
- RED confirmed at base: `folder-structure-check` on this repository reports `FOLDER-STRUCTURE-JUST-LOGIC` for `justfile: kanban-static` (2) and `justfile: run-kanban` (3). The same check over a temp repository whose justfile is the template reports the same two gaps and nothing else, so every consumer gets them.

**What pins the old recipe text (must change with it)**
- `_dev/tests/contracts/queue-kanban.sh:6-22` extracts the run-kanban shutdown line from both the root justfile and the template with awk, fails when they differ, then `eval`s it with function seams for `lsof`/`ps`/`kill`/`sleep` (`:23-88`: stuck listener refused naming pid and command; foreign listener refused without `kill`). Run by `_dev/tests/contract-regressions.sh:58`. No case covers the happy path (stale queue-kanban stopped, then serve) or kanban-static.
- `_dev/tests/staged-skills-contract.sh:928-991`: `:930-933` requires exactly 4 occurrences of `.claude/skills/do-work-board/tools/queue-kanban` in the template (drops to 2 after the change); `:983-990` greps the template for `[ "$wait_count" -lt 320 ]`, `remaining_listener_pid=`, `listener_executable_name=`, `serve --open` (these move into the script). Heavy-tier only; 32 s here.

**What does not pin bodies (no change needed)**
- Installer/updater copy the template verbatim and read recipe names only: `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go:539,612`; `install_transaction_test.go:488-500` (41 definitions); `_dev/tests/install-suite-behavior.sh:256-260` (`cmp` with the template, 41 recipes); `_dev/tests/update-script-behavior.sh:351-361` (greps `.claude/skills/do-work-board/tools/queue-kanban`, still matched by `run-kanban-cli` and `kanban-summary`; recipe-name parity); `publication_commands_test.go:106` runs a `do-work-*` recipe; `_dev/tests/contracts/replace-text-section.sh:73-79` reads reserved names.
- `suite/modules.tsv` ships the whole `skills/do-work-board` folder to `.claude/skills/do-work-board`; `.gitattributes` export-ignores nothing under it. A new file under `skills/do-work-board/tools/` ships with no manifest edit.
- `_dev/tests/prescribed-shell-canonicalization.sh:39-57` scans `skills/*/tools/*.sh` for a direct `curl`; the new scripts have none. The shellcheck lane lints every tracked `*.sh` at `--severity=warning`.

**Layout rule and placement**
- `folder_structure.go:58-83,195-231`: the check reads only the repository's root entries and its root justfile. It never looks inside `skills/<package>/`, so rule 5's `cmd/<name>.sh` is about a consumer's own root; the suite must not write into a consumer's `cmd/`. Core's `skills/do-work/tools/` already holds shell launchers (`do-work-cli.sh`, `do-work-update.sh`); `scripts/` holds primitives that actions prescribe. See D-01.
- `skills/do-work-board/SKILL.md:24-25` says the package owns the template bytes and that `tools/queue-kanban/` is the only board implementation.

**Shell traps found while probing**
- just echoes each recipe line to stderr before running it, so the old echoed `kanban-static` line contains the exclude message text. A check for the message must match a whole line (`grep -x`), or it passes on the echo.
- Under `set -e -o pipefail`, `listener_pid="$(lsof -ti ... | head -n1)"` fails when nothing listens (lsof exits 1), which would kill run-kanban before it serves. Under `set -e`, a failure in `cd && go build && ./queue-kanban generate` that is not the list's last command does not stop the script, so the exclude write would run after a failed build. See D-03.
- `queue-kanban serve --open` opens the browser through `open` (darwin) or `xdg-open` found on `PATH` (`skills/do-work-board/tools/queue-kanban/openbrowser.go:17-27`); a stub first on `PATH` keeps a live proof headless.

**Proof tools written for this REQ** (in `do-work/runs/2026-11-09-layout-queue/`)
- `REQ-717-live-proof.sh`: copies the board package (without its binary) into a throwaway Git consumer with the template as its justfile; runs `just run-kanban <port>` twice and checks the second run stops the first queue-kanban pid and serves in its place; runs `just kanban-static` twice and checks the snapshot, exactly one exclude line, and one announcement; checks the root justfile reaches the port validation. Every server it starts is stopped on exit. Passed at base in about 2 s (characterization).
- `REQ-717-preflight-probe.sh` (baseline, green at base) and `REQ-717-probe.sh` (GREEN probe, red at base).

*Generated by pre-dispatch agent*

## Decisions

- D-01 DECIDE & STATE: the scripts are `skills/do-work-board/tools/run-kanban.sh` and `skills/do-work-board/tools/kanban-static.sh`, named after their recipes. `folder-structure-check` never inspects a skill package's folders, so the `cmd/` rule does not apply inside the package, and a suite package cannot place files in a consumer's `cmd/`. Core's `tools/` already holds shell launchers, and these two launch the board tool.
- D-02 DECIDE & STATE: positional arguments only. `run-kanban.sh <project-root> <port>`, `kanban-static.sh <project-root>`. Each script finds `queue-kanban/` beside itself. The default port stays in the recipe header (`$port="8090"`), so no flag or default exists without a caller. The recipes call `bash "<path>"` like the core recipes, so the exec bit does not matter; commit the files as 100755 like every other shipped `.sh`.
- D-03 DECIDE & STATE: `run-kanban.sh` runs under `set -u` only, which matches just's `sh -cu` for the moved lines (the lsof pipelines exit 1 when nothing listens). `kanban-static.sh` must stop when the build or generate fails, before the exclude write: put build and generate on separate lines under `set -e`, or end the chain with `|| exit`. Either way a comment says why.
- D-04 DECIDE & STATE: messages, exit codes, the 320-step wait and `serve --open` move verbatim. No `exec`, no `@` quiet prefix, no new output. The one visible change is just's echo of the recipe line, which now shows the script call.
- D-05 DECIDE & STATE: template lines are `bash "{{justfile_directory()}}/.claude/skills/do-work-board/tools/<name>.sh" "{{justfile_directory()}}"` (plus `"$port"`). The root justfile uses `bash skills/do-work-board/tools/<name>.sh "{{justfile_directory()}}"` (plus `"$port"`), the repository-relative style its board recipes already use. Recipe headers stay unchanged (the contract and the REQ-713 test match `run-kanban $port=`).
- D-06 DECIDE & STATE: tests. Delete the root/template shutdown-line parity check, because one script replaces the two copies it compared. Drive the script from a fixture copy with stub `go` and a stub `queue-kanban`, so no case can start the real binary. Add the cases nothing covered: a stale queue-kanban listener is stopped and serve receives the port, kanban-static writes one exclude line across two runs, and a failed build writes no exclude line. Keep the stuck and foreign cases. The staged-skills contract reads the safety strings from the script the template calls (REQ-555: check the file you name). One Go test runs `justfileLogicGaps` over the shipped template, because the check on this repository reads only the root justfile and nothing else would catch a consumer-side gap.
- D-07 DECIDE & STATE: `skills/do-work-board/SKILL.md` gets one sentence in its ownership list naming the two recipe scripts, since it says `tools/queue-kanban/` is the only board implementation. board.md, help.md and board-guide.md name recipes, not bodies, and stay untouched.
- D-08 DECIDE & STATE: the live proof (real servers, real ports) stays out of the GREEN probe. The probe must be quick and must not depend on ports or start servers. The builder runs the live proof once and records the result, and the integrator may run it again.
<!-- D-XX counter: last used D-08. Next decision: D-09. -->

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/run-kanban.sh` (new) — port validation, stale-listener stop, build and serve
- `skills/do-work-board/tools/kanban-static.sh` (new) — build and generate, then the local exclude line
- `skills/do-work-board/justfile.template` (modify) — run-kanban and kanban-static become one-line script calls
- `justfile` (modify) — the same two recipes in this repository's managed block
- `skills/do-work-board/SKILL.md` (modify) — one ownership sentence naming the two recipe scripts
- `_dev/tests/contracts/queue-kanban.sh` (modify) — drive the scripts; stale, stuck, foreign, static and failed-build cases
- `_dev/tests/staged-skills-contract.sh` (modify) — board path count and safety strings read from the called script
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go` (modify) — the shipped template holds no multi-line recipe

**Files I will NOT touch:** `folder_structure.go` (the check is correct); `skills/do-work-board/actions/board.md`, `actions/help.md`, `docs/board-guide.md`; the installer and updater Go code and their tests; `_dev/tests/install-suite-behavior.sh`, `update-script-behavior.sh`, `contracts/replace-text-section.sh`, `contract-regressions.sh`, `heavy-lanes.json`, `fast-stages.json`; anything REQ-716 owns (root `tools/`, `cmd/`, the `legacy_runtime_paths` and `retained_bootstrap_tool` parts of the staged-skills contract, `heavy_maintainer_tree_test.go`, `README.md`); `VERSION`, both changelogs and `skills/do-work/actions/version.md` (the finalizer's).

**Acceptance criteria (restated from REQ):**
- [ ] `folder-structure-check` on this repository reports no `FOLDER-STRUCTURE-JUST-LOGIC` gap.
- [ ] In `skills/do-work-board/justfile.template`, `run-kanban` and `kanban-static` are each one line that calls a script shipped with do-work-board, with the same behavior and messages.
- [ ] `just run-kanban` still replaces a stale queue-kanban instance on its port, and still refuses a foreign listener and a listener that survives the wait.
- [ ] `just kanban-static` still writes the `build/` snapshot and the exclude line.
- [ ] This repository's root justfile managed block follows, and the board contract tests pass.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/` at 792f981c; only the claimed REQ files, `baseline.json` and the run directory are changed
**Tests baseline:** ✓ Focused pre-flight probe green (`bash do-work/runs/2026-11-09-layout-queue/REQ-717-preflight-probe.sh`); repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 129s at 792f981c, run from a detached checkout
**Dependencies:** ✓ None added

*Checked by work action*
