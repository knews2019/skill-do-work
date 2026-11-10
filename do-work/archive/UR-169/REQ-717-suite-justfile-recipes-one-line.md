---
id: REQ-717
title: 'Suite-managed justfile recipes call a script in one line'
status: completed
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
write_set: [skills/do-work-board/tools/run-kanban.sh, skills/do-work-board/tools/kanban-static.sh, skills/do-work-board/justfile.template, justfile, skills/do-work-board/SKILL.md, _dev/tests/contracts/queue-kanban.sh, _dev/tests/staged-skills-contract.sh, skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go, _dev/tests/audit-lockins.sh]
dispatch_at: 2026-11-10T00:35:00Z
builder_handback_at: 2026-11-10T00:43:51Z
integration_at: 2026-11-10T00:45:15Z
review_at: 2026-11-10T00:52:58Z
kb_status: pending
commit: 36012976bbe2f21f3b347574a2be7ea80fb02374
heavy_verified_at: 2026-11-10T00:56:26Z
heavy_verified_revision: 36012976bbe2f21f3b347574a2be7ea80fb02374
completed_at: 2026-11-10T00:57:07Z
release_at: 2026-11-10T00:57:07Z
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
- [x] **[PLAN]:** Tests first (contract fixture with stub programs, staged block, Go lock-in), then two scripts holding the moved lines, then one-line recipes and the SKILL.md bullet.
- [x] **[APPLY]:** Built as planned, plus the early folder resolution (D-10) and the lock-in change (D-09).
- [x] **[UNIFY]:** git diff --stat 2f0ee277..b8b07228: 9 files, +239/-58; each checked: messages, exit codes, the 320-step wait and serve --open unchanged; both justfiles change only the two bodies; no debug output.
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
- D-09 (builder ESCALATE, coordinator decided O1): the REQ-551 lock-in in `_dev/tests/audit-lockins.sh` fails any shipped `.sh` that does not mention `do-work-cli`, and both new scripts failed it. The check now also accepts an invocation of `./queue-kanban `, the board's own Go program; a script with no Go call still fails. Value: the smallest change, and the same shell logic already shipped inside the template, so shipped shell does not grow. Risk: the lock-in matches a second string; reverting commit b8b07228 undoes it. The alternative, moving the stale-listener stop and the exclude write into queue-kanban's Go code, is recorded under Discovered Tasks.
- D-10 DECIDE & STATE: both scripts resolve `queue-kanban/` beside themselves once, at the top. macOS `/bin/bash` 3.2 unsets `BASH_SOURCE[0]` after a script calls a function imported from the environment, which is how the contract stubs `lsof` and `kill`.
- D-11 DECIDE & STATE: missing arguments stop with a usage message; the port uses the no-colon form so an empty port still reaches the original invalid-port message.
- D-12 DECIDE & STATE: the kanban-static exclude test is split over lines for readability; logic unchanged.
<!-- D-XX counter: last used D-12. Next decision: D-13. -->

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
- `_dev/tests/audit-lockins.sh` (modify, added at integration, D-09) — the shipped-shell delegation check also accepts a call to the board's Go program

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

## Discovered Tasks

- The delegation check in `_dev/tests/audit-lockins.sh` is a plain string match, so a `do-work-cli` mention in a comment satisfies it; it cannot tell a thin launcher from a script with logic. impact-negligible → report only
- D-09's alternative: move the stale-listener stop into `queue-kanban serve` and the exclude write into Go, so the two recipes call `queue-kanban` directly and no recipe script ships. A REQ-sized change that adds a write surface `prime-kanban-board.md` counts. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/run-kanban.sh` (created)
- `skills/do-work-board/tools/kanban-static.sh` (created)
- `skills/do-work-board/justfile.template` (modified)
- `justfile` (modified)
- `skills/do-work-board/SKILL.md` (modified)
- `_dev/tests/contracts/queue-kanban.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go` (modified)
- `_dev/tests/audit-lockins.sh` (modified)

**What was done:** The bodies of the `run-kanban` and `kanban-static` recipes moved verbatim into two scripts shipped with do-work-board, and each recipe in the template and in this repository's justfile is now one line that calls its script. The queue-kanban contract drives the scripts with stub `go` and `queue-kanban` programs (stale, stuck, foreign, static and failed-build cases) in place of the deleted shutdown-line parity check. The staged-skills contract reads the safety strings from the called script, a Go test checks that the shipped template holds no multi-line recipe, and the shell lock-in accepts a call to the board's Go program (D-09).

## Qualification

**Range:** `7cbbf980..36012976` (one merge of the REQ branch: fbdb1093 tests, 9b5bcac0 scripts and recipes, b8b07228 lock-in)
**Mechanical gates:** qualify ✓ satisfied, scope-drift ✓ satisfied (9 declared paths after the D-09 Scope addition, 9 touched).

**Requirement trace (orchestrator, read against the diff):**
- Template: `run-kanban` and `kanban-static` are each one line calling `.claude/skills/do-work-board/tools/<name>.sh` with the project root (and the port). Recipe headers and the 41 recipe names are unchanged.
- Root justfile: the same two bodies, calling the scripts by their repository path.
- Behavior: every message, exit code, the 320-step wait and `serve --open` moved verbatim. `run-kanban.sh` keeps just's `set -u` only, because the lsof pipelines exit 1 when nothing listens. `kanban-static.sh` puts build and generate on their own lines under `set -e`, so a failed build stops before the exclude write, as the two-line recipe did.
- Tests: the contract drives both scripts with stub programs; the deleted parity check compared two copies that no longer exist. The staged contract reads the safety strings from the script the template calls. A Go test runs the layout check over the shipped template.
- Lock-in: the REQ-551 delegation check accepts a call to `./queue-kanban` (D-09); a script with no Go call still fails.

**Live flow:** the builder's live proof started real servers in a throwaway consumer: a stale queue-kanban on its port was replaced, the static snapshot and one exclude line were written across two runs, and no listener was left.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/2026-11-09-layout-queue/REQ-717-probe.sh` (the Go template lock-in, the queue-kanban contract and the staged-skills board block), then `bash _dev/tests/maintainer-verify.sh` from a detached checkout at the merge revision. The builder and the reviewer each ran `REQ-717-live-proof.sh` (real servers in a throwaway consumer).
**Result:** ✓ Focused probe green. Repository gate exit 0 at the merge 36012976 (gate wall 162s). Live proof: three PASS lines both times, no listener left on ports 18717 to 18757. The contract passes under bash 5 and macOS `/bin/bash` 3.2.

**Red-green validation:**
- `TestShippedBoardTemplateRecipesHoldNoLogic`: ✗ at fbdb1093 (two `FOLDER-STRUCTURE-JUST-LOGIC` gaps, run-kanban 3 body lines, kanban-static 2) → ✓ at 9b5bcac0.
- `_dev/tests/contracts/queue-kanban.sh`: ✗ at fbdb1093 (`skills/do-work-board/tools/run-kanban.sh is missing`) → ✓ at 9b5bcac0. Two deliberate breaks were caught and reverted: an `&&` chain in kanban-static (case 5) and `pipefail` in run-kanban (case 1).

**New tests added:**
- `TestShippedBoardTemplateRecipesHoldNoLogic` in `folder_structure_test.go`
- Five fixture cases in `_dev/tests/contracts/queue-kanban.sh`: stale listener replaced, stuck listener refused, foreign listener refused, static run writes one exclude line across two runs, failed build writes no exclude line

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contracts/queue-kanban.sh` (from REQ-633): the root/template shutdown-line parity check is deleted, because one script replaces the two copies it compared; the REQ-633 stamp check is unchanged.
- `_dev/tests/audit-lockins.sh` (from REQ-551): the shipped-shell delegation check also accepts a call to `./queue-kanban` (D-09).

**Heavy verification plan:**
- Range: 7cbbf9802418bff72007e7732b3e0e0449627492..36012976bbe2f21f3b347574a2be7ea80fb02374
- All six lanes are selected because coverage is uncertain for the root `justfile`.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

*Verified by work action*

## Review

**Overall: 95%** | 2026-11-10T00:52:58Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. The two recipes are one-line script calls in both justfiles. The scripts keep the old messages, exit codes, the 320-step wait and `serve --open`. The installer replaces the whole managed block, so `do-work update` delivers the new recipes and the scripts together.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1: After the deleted root/template parity check, no gate test pins the root `justfile` `run-kanban`/`kanban-static` bodies (`justfile:6`, `justfile:14`). Only the run-local probe and live proof read them, and `folder-structure-check` is not in the gate. A wrong path or a return to inline logic in the root block would pass the gate. Fix: in `_dev/tests/contracts/queue-kanban.sh`, assert that each of the two root recipe bodies is exactly `bash skills/do-work-board/tools/<name>.sh "{{justfile_directory()}}"` (plus `"$port"`). — impact-negligible → report only
- M2: Both scripts use `project_root` after `cd "$queue_kanban_directory"` (`skills/do-work-board/tools/kanban-static.sh:12-15`, `run-kanban.sh:38`). A relative root therefore resolves against the tool folder. Reproduced in scratch: `bash skills/do-work-board/tools/kanban-static.sh .` wrote the snapshot to `tools/queue-kanban/build/`, and the exclude line still went to the project. Both recipes pass an absolute `{{justfile_directory()}}`, so only a direct caller can hit this. Fix: `project_root="$(cd "${1:?usage: ...}" && pwd)"` at the top of both scripts. — impact-negligible → report only
- M3: D-09 widened the REQ-551 delegation lock-in for every shipped `.sh`. Any shipped `.sh` under `skills/` or `tools/` that contains `./queue-kanban ` now passes, not only the two board recipe scripts (`_dev/tests/audit-lockins.sh:65`). Fix: accept `\./queue-kanban ` only when `$f` is under `skills/do-work-board/tools/`, and keep `do-work-cli` as the rule everywhere else. — impact-negligible → report only

**Nit findings:**
- N1: `skills/do-work-board/actions/board.md:58` says the git-optional root handling is "matching `justfile.template`". That git-optional exclude guard now lives in `tools/kanban-static.sh`, not in the template. Point the phrase at the script. — impact-negligible → report only
- N2: `kanban-static.sh:16-21` keeps an inline copy of the exclude write. It does not use the canonical Local Git ignore primitive (`skills/do-work/scripts/add-local-git-exclude.sh`, which `board.md` static mode already calls). This is correct for this REQ, because the REQ required the exact old message. It belongs with the D-09 alternative recorded under Discovered Tasks. — impact-negligible → report only

**Requirements checklist:**
- [x] `folder-structure-check` on this repository reports no `FOLDER-STRUCTURE-JUST-LOGIC`. Delivered: probe step 1 is green.
- [x] Each template recipe is one line that calls a script shipped with do-work-board. Delivered: `justfile.template:6`, `:14`, plus `TestShippedBoardTemplateRecipesHoldNoLogic`.
- [x] Same behavior and messages. Delivered: the diff shows the lines moved verbatim. D-03 runs `run-kanban.sh` under `set -u` only, and `kanban-static.sh` under `set -e` with build and generate on separate lines. Reviewer-noted differences (not regressions): bash `echo` replaces sh `echo`, and under `set -e` a failed `mkdir`/`git rev-parse` inside the exclude block now exits non-zero. The old recipe printed a false success line in that case.
- [x] `just run-kanban` replaces a stale instance and refuses a foreign listener and a listener that survives the wait. Delivered: contract cases 1-3, and the live proof replaced pid 99753 with 99907 on :18717.
- [x] `just kanban-static` writes the `build/` snapshot and exactly one exclude line. Delivered: contract cases 4-5 and the live proof.
- [x] The root justfile managed block follows, and the board contract tests pass. Delivered: `justfile:6`, `:14`, and the contract passes under bash 5 and `/bin/bash` 3.2.

**Specific judgments:**
- `dirname "${BASH_SOURCE[0]}"`: The root justfile has no `set working-directory` and no `[no-cd]`, so just runs the recipe in the justfile folder. The relative call `bash skills/do-work-board/tools/run-kanban.sh` then resolves to the right folder, because `run-kanban.sh` does not `cd` before using it. `just --dry-run` from `skills/do-work/` showed the same line. The installed template passes an absolute path. Scratch runs proved both forms reach `tools/queue-kanban` with the right `--repo-root`. D-10 is correct: under `/bin/bash` 3.2.57, `BASH_SOURCE[0]` was unset after a call to an `export -f` function.
- `do-work update`: `RunUpdate` runs the same install transaction. `buildJustCandidate` (`install_transaction.go:538-581`) extracts the template's single managed section and replaces the consumer's marked block wholesale through `managedsection.ReplaceSection`. Then `validateJustCandidate` checks the markers, all 41 recipe names and `just --list`. The scripts ship with the whole `skills/do-work-board` folder (`suite/modules.tsv`) in the same transaction, so a consumer cannot get the new recipe lines without the scripts. A consumer that keeps an old inline block still works, because `tools/queue-kanban/` did not move. `go test ./internal/suiteinstall/ ./internal/managedsection/` passed. That includes the shipped-template test at `install_transaction_test.go:488`.
- D-09: Sound as an interim decision. The amount of shipped shell did not grow, because the same logic shipped inline in the template before. The decision is recorded, coordinator-decided and in its own revertible commit (b8b07228), and the Go-first alternative is recorded under Discovered Tasks. Its weakness is breadth (M3), not direction.

**Acceptance:** Pass. Covered stages: implementation (diff walk, scratch runs of both scripts by relative and absolute path) and integration (merged tree: GREEN probe, contract under bash 5 and 3.2, `audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, the heavy staged-skills contract, shellcheck, and the Go tests for toolboxcommands, suiteinstall and managedsection). The live proof passed in a throwaway consumer, and afterwards nothing listened on 18717-18757.
**Restatement sweep:** redefined (1) where the `run-kanban` and `kanban-static` recipe logic lives (from inline recipe bodies to `skills/do-work-board/tools/*.sh`), and (2) the REQ-551 shipped-shell delegation lock-in (it now accepts `./queue-kanban `). Swept `skills/`, `README.md`, `_dev/`, `CLAUDE.md`, `suite/` and `docs/`. `board.md`, `help.md`, `board-guide.md`, `command-line-guide.md` and `prime-kanban-board.md` name recipes, not bodies, and agree. The parity helper `extract_kanban_shutdown_line` survives only in dated KB records (`docs/kb/wiki/sources/REQ-180-*`, `docs/kb/raw/processed/2026-09-01/REQ-180-*`), which are historical. The fixture recipe bodies in `update-script-behavior.sh:312-323` and the installer Go tests are synthetic inputs and are still valid. No other text restates the lock-in wording. One stale restatement was found: N1.
**Suggested testing:** 4 items
- Deployment: unassessed. Run `do-work update` from the prior release in a real consumer. Confirm that the managed block now calls `.claude/skills/do-work-board/tools/run-kanban.sh` and that the file exists. The live proof copied the package and did not run the installer.
- Live acceptance: unassessed. In an interactive terminal, press Ctrl-C on `just run-kanban`. Confirm that the server stops and the prompt returns, because the call now adds a bash process between just and queue-kanban.
- Linux consumer (dash as `sh`): run `just run-kanban` twice on one port and confirm that the stale replacement works.
- Run the heavy `_dev/tests/install-suite-behavior.sh` lane, which compares an installed justfile with the template. This review did not run it (heavy-only, needs maintainer permission).
**Follow-ups created:** None (5 findings report only)

### Domain Review
`crew-members/general.md` defines no review rubric. Its discipline rule held: the builder recorded the weakness of the lock-in's plain string match as a Discovered Task and did not patch it inline. The PLAN, APPLY and UNIFY boxes are all checked. Decisions D-01 to D-12 are in the REQ and the hand-back, and they cover every behavior-relevant choice in the diff.

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: 7cbbf9802418bff72007e7732b3e0e0449627492
- Target revision: 36012976bbe2f21f3b347574a2be7ea80fb02374
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for the root `justfile`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for the root `justfile`
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for the root `justfile`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for the root `justfile`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for the root `justfile`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for the root `justfile`

## Heavy Verification Result

- Target revision: 36012976bbe2f21f3b347574a2be7ea80fb02374
- Execution revision: 36012976bbe2f21f3b347574a2be7ea80fb02374 (detached checkout)
- queue-kanban-javascript: exit 0, 15s, executed
- queue-kanban-browser: exit 0, 93s, executed (real browser, not skipped)
- do-work-cli-integrations: exit 0, 80s, executed
- staged-skills: exit 0, 37s, executed
- updater: exit 0, 74s, executed
- installer: first run exit 1, 64s (`install-suite-behavior.sh`: "ordinary installation failed while discarding confirmed managed customizations", while the REQ-716 pre-dispatch agent ran installer tests at the same time); rerun alone exit 0, 44s, executed. Load-only: the failing case touches only `.claude/skills/do-work/SKILL.md`, not the justfile or the board scripts.

## Lessons Learned

**What worked:** Moving the recipe lines verbatim and proving parity with stub programs (`go`, `queue-kanban`, `lsof`, `kill`) in a fixture copy, so no test case can start the real board. Mutation checks (an `&&` chain, `pipefail`) showed that the new cases catch the two ways the scripts would most likely break.
**What didn't:** The first script draft read `BASH_SOURCE[0]` late; under macOS `/bin/bash` 3.2 it was unset after the script called a stubbed function, so only the test run died. The REQ-551 shell lock-in also failed both new scripts, because it accepts only a `do-work-cli` call; that needed a decision (D-09) the pre-dispatch did not foresee.
**Worth knowing:** `run-kanban.sh` must stay `set -u` only: its lsof pipelines exit 1 when nothing listens. `kanban-static.sh` relies on `set -e` with build and generate on their own lines, because `set -e` ignores every command of an `&&` list but the last. Both recipe scripts expect an absolute project root (the recipes pass `justfile_directory()`); a relative root writes the snapshot under the tool folder (review M2).

## Orientation

Now the suite's `run-kanban` and `kanban-static` recipes are one-line calls to `skills/do-work-board/tools/run-kanban.sh` and `kanban-static.sh`, so a consumer's justfile holds no recipe logic and `folder-structure-check` reports no `FOLDER-STRUCTURE-JUST-LOGIC` gap for the suite's own recipes. This lives in the board package's recipe launchers (`_dev/primes/prime-kanban-board.md`). A leaf change: the logic and messages did not change. Prime staleness check: `prime-kanban-board.md` and `prime-shell-commands.md` name no recipe body, so neither is stale.
