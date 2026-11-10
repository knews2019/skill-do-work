## Review

**Overall: 95%** | REVIEW_AT

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
