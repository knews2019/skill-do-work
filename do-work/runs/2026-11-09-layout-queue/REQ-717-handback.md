# REQ-717 hand-back (suite-managed justfile recipes call a script in one line)

Branch `worktree-agent-REQ-717-justfile-recipes-call-scripts`, tip `b8b07228`, base `2f0ee277`.

## Commits

- `fbdb1093` [REQ-717] tests: board recipe scripts drive the launcher contract
- `9b5bcac0` [REQ-717] run-kanban and kanban-static recipes call a board script in one line
- `b8b07228` [REQ-717] audit lock-in: board recipe scripts delegate to queue-kanban (scope change, see D-09)

## Files

- created `skills/do-work-board/tools/run-kanban.sh` (100755)
- created `skills/do-work-board/tools/kanban-static.sh` (100755)
- modified `skills/do-work-board/justfile.template` (run-kanban and kanban-static bodies only)
- modified `justfile` (run-kanban and kanban-static bodies only)
- modified `skills/do-work-board/SKILL.md` (one Ownership boundary bullet)
- modified `_dev/tests/contracts/queue-kanban.sh` (parity check and `extract_kanban_shutdown_line` deleted; five fixture cases; REQ-633 stamp check unchanged)
- modified `_dev/tests/staged-skills-contract.sh` (board-template block only: tools-path count, two `require_file`, safety strings read from `run-kanban.sh`)
- modified `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go` (`TestShippedBoardTemplateRecipesHoldNoLogic`)
- modified `_dev/tests/audit-lockins.sh` (outside the declared Scope, see Scope changes)

## Scope changes

- `_dev/tests/audit-lockins.sh`: the REQ-551 lock-in "Shipped shell delegating check" fails any `.sh` under `skills/` or `tools/` that does not contain `do-work-cli`. Both new scripts failed it (`FAIL: NON-DELEGATING: .../run-kanban.sh`, `.../kanban-static.sh`), which made `contract-regressions.sh` exit 1. The check now also accepts `./queue-kanban ` (an invocation of the board's own Go binary). A mutation file with no Go call still fails the check. Kept as its own commit so it can be dropped if D-09 goes the other way.

## Red-green evidence

Commit 1 (`fbdb1093`, tests only), before the scripts existed:
- `go test -run '^TestShippedBoardTemplateRecipesHoldNoLogic$'`: `FAIL ... shipped board template recipes hold logic: [{FOLDER-STRUCTURE-JUST-LOGIC justfile: run-kanban recipe has 3 body lines ...} {FOLDER-STRUCTURE-JUST-LOGIC justfile: kanban-static recipe has 2 body lines ...}]`
- `bash _dev/tests/contracts/queue-kanban.sh`: exit 1, `FAIL: skills/do-work-board/tools/run-kanban.sh is missing; the run-kanban and kanban-static recipes call it.`
- `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`: exit 1, two `required staged-suite file is missing` lines and four `board run-kanban script lost safety/launch contract` lines.

After commit 2/3: the Go test passes, the contract prints `queue-kanban contract probes passed.` under `bash` and `/bin/bash` (3.2.57), the staged contract prints `staged skills contract: PASS`.

Mutation checks on the contract (reverted after): an `&&` chain for build+generate in kanban-static.sh fails case 5 (`got exit 0`); `set -euo pipefail` in run-kanban.sh fails case 1 (`got exit 1`, no serve argv).

## Verification

1. `gofmt -l ./internal/toolboxcommands`: empty. `go vet`: ok. `go test -count=1 ./internal/toolboxcommands/`: ok (23-26 s).
2. `bash _dev/tests/contracts/queue-kanban.sh`: pass. `bash _dev/tests/contract-regressions.sh`: exit 0, no FAIL line (before commit 3 it failed on NON-DELEGATING only). `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`: PASS in 39 s. `bash _dev/tests/prescribed-shell-canonicalization.sh`: pass. `bash _dev/tests/quiet-grep-pipeline-audit.sh` after commit 2: pass (210 tracked shell files). `shellcheck --severity=warning` on both scripts, both changed test scripts and `audit-lockins.sh`: clean.
3. `REQ-717-probe.sh`: `REQ-717 probe passed.`
4. `REQ-717-live-proof.sh`: `PASS run-kanban: pid 87023 replaced pid 86241 on :18717`, `PASS kanban-static: snapshot written, one exclude line, idempotent`, `PASS root justfile: run-kanban validates the port; kanban-static calls skills/do-work-board/tools/kanban-static.sh`. Afterwards `lsof -ti tcp:<port> -sTCP:LISTEN` printed nothing for 18717, 18727, 18737, 18747 and 18757.
5. `folder-structure-check` on the worktree: 0 `FOLDER-STRUCTURE-JUST-LOGIC`; 5 `FOLDER-STRUCTURE-ROOT-ENTRY` remain (`VERSION`, `_dev/`, `skills/`, `suite/`, `tools/`).
6. `just --justfile skills/do-work-board/justfile.template --summary`: 41 names, identical to the base template's summary (`diff` empty).

## P-A-U

- PLAN: tests first (contract fixture with stub queue-kanban and exported function stubs, staged block, Go lock-in), then two scripts with the lines moved verbatim, then one-line recipes and the SKILL.md bullet.
- APPLY: built as planned; two additions: the tool folder is resolved at script top (bash 3.2 trap, see Proposed lesson), and the audit lock-in accepts `./queue-kanban ` (D-09).
- UNIFY: `git diff --stat 2f0ee277..HEAD` shows 9 files, +239/-58; checked each: scripts hold the base lines' messages, exit codes, 320-step wait and `serve --open`; both justfiles change only the two bodies; staged edit is inside the board-template block only (no change near `legacy_runtime_paths`); no debug output or scratch files in the diff.

## Decisions

- D-09 ESCALATE: shipped shell versus the thin-launcher rule. UR-081 and the REQ-551 lock-in say a shipped `.sh` is a thin launcher to Go; the two new scripts carry about 30 lines of real shell (lsof/ps/kill, exclude write). I widened the lock-in to accept the board's Go binary (commit 3) so the REQ ships as decided in D-01 to D-04.
  - O1 (done, recommended): keep commit 3. Value: smallest change; the same shell logic already shipped inside the template, so the amount of shipped shell does not grow. Risk: the lock-in now matches a second string; reversible by reverting one commit.
  - O2: move the stale-listener stop into `queue-kanban serve` and the exclude write into Go, so the recipes become one-line `queue-kanban` calls and no new `.sh` ships. Value: matches UR-081's direction fully. Risk: a new REQ-sized change to queue-kanban; the exclude write would be a fourth write surface, which `prime-kanban-board.md` counts and `contract-regressions.sh` pins.
- D-10 DECIDE & STATE: both scripts resolve `queue-kanban/` beside themselves once, at the top. Under bash 3.2, `BASH_SOURCE[0]` is unset after the script calls a function imported from the environment, which is how the contract stubs `lsof`/`kill`; production calls real commands, but resolving early costs nothing.
- D-11 DECIDE & STATE: missing arguments stop with `${1:?usage: ...}`; the port uses `${2?...}` (no colon) so an empty port still reaches the original invalid-port message.
- D-12 DECIDE & STATE: the kanban-static exclude `if` is split over lines for readability; logic unchanged.

## Discovered Tasks

- impact-low: `_dev/tests/audit-lockins.sh` delegation check is a plain string match, so a `do-work-cli` mention in a comment satisfies it; it cannot tell a thin launcher from a script with logic → report only

## Integration seams

Confirmed: the staged-skills edit is only inside the board-template block (around lines 928-995); REQ-716's `legacy_runtime_paths` / `retained_bootstrap_tool` hunks are untouched. The root `justfile` edit is the two bodies only. Correction: this branch also edits `_dev/tests/audit-lockins.sh` (one hunk, "Shipped shell delegating check"). Its `find "$repo_root/skills" "$repo_root/tools"` scan names the root `tools/` that REQ-716 deletes; if REQ-716 edits that line, the merge conflicts on the same hunk and the resolution keeps both changes. `heavy-lanes.json` and `heavy_maintainer_tree_test.go` are unchanged.

## Proposed CHANGELOG entry

**Board recipes call shipped scripts in one line** — the `run-kanban` and `kanban-static` recipes move their logic into `skills/do-work-board/tools/run-kanban.sh` and `tools/kanban-static.sh`, so `folder-structure-check` no longer reports `FOLDER-STRUCTURE-JUST-LOGIC` for the suite's own recipes.

What changes for an installed copy: the run-kanban and kanban-static recipes now call `.claude/skills/do-work-board/tools/*.sh`, and `do-work update` rewrites the managed block. Behavior and messages are unchanged; just now echoes the one-line script call instead of the inline shell.

## Proposed lesson

- [family: exported-function-stub-traps] REQ-717: macOS `/bin/bash` 3.2 unsets `BASH_SOURCE[0]` after a script calls a function imported with `export -f`, so a script under `set -u` that reads `BASH_SOURCE` late dies only under its stubbed test; resolve the script's own folder at the top.

## Lessons read

- `skills/do-work/tools/lessons-do-work-update.md` (required): whole file.
- `_dev/primes/lessons-shell-commands.md`: `check-reads-the-file-it-names` (REQ-555, applied: the staged contract now reads the safety strings from the called script), `legacy-fixture-implementation-shape` (REQ-420, applied: the contract asserts exit status, output lines, kill and serve argv, and the exclude file), and "Every Flag on a Shipped Script Needs a Non-Test Caller" (applied: no flags, positional arguments only).
- `_dev/primes/prime-shell-commands.md`, `_dev/primes/prime-kanban-board.md`.
