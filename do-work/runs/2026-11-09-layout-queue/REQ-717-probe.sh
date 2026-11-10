#!/usr/bin/env bash
# GREEN focused probe for REQ-717 (suite-managed justfile recipes call a script in one line).
# Expected to FAIL at the base revision: the two recipes still hold inline logic there.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

probe_fail() {
  printf 'REQ-717 probe FAIL: %s\n' "$1" >&2
  exit 1
}

# 1. Layout rule 5 on this repository's own justfile, through the real check.
folder_output="$(bash skills/do-work/tools/do-work-cli.sh --repo-root "$root" --format text folder-structure-check 2>&1 || true)"
grep -qF 'folder-structure-check:' <<<"$folder_output" || probe_fail "folder-structure-check did not run: $folder_output"
if grep -qF 'FOLDER-STRUCTURE-JUST-LOGIC' <<<"$folder_output"; then
  probe_fail "folder-structure-check still reports recipe logic: $folder_output"
fi

# 2. The shipped template (what every consumer receives) holds no multi-line recipe.
go_output="$(cd skills/do-work/tools/do-work-cli && go test -count=1 -v -run '^TestShippedBoardTemplateRecipesHoldNoLogic$' ./internal/toolboxcommands/ 2>&1)" \
  || probe_fail "template rule-5 test failed: $go_output"
grep -qF -- '--- PASS: TestShippedBoardTemplateRecipesHoldNoLogic' <<<"$go_output" \
  || probe_fail "template rule-5 test did not run: $go_output"

# 3. The two recipes call scripts shipped with do-work-board, and the scripts exist.
for board_script in run-kanban kanban-static; do
  [ -f "skills/do-work-board/tools/$board_script.sh" ] || probe_fail "missing skills/do-work-board/tools/$board_script.sh"
  grep -qF ".claude/skills/do-work-board/tools/$board_script.sh" skills/do-work-board/justfile.template \
    || probe_fail "template recipe does not call tools/$board_script.sh"
  grep -qF "skills/do-work-board/tools/$board_script.sh" justfile \
    || probe_fail "root justfile recipe does not call tools/$board_script.sh"
done

# 4. The board launcher contract drives the scripts (stale replace, stuck and foreign listener,
#    kanban-static snapshot and exclude line), not the retired inline recipe text.
contract_file='_dev/tests/contracts/queue-kanban.sh'
grep -qF 'tools/run-kanban.sh' "$contract_file" || probe_fail "$contract_file does not exercise tools/run-kanban.sh"
grep -qF 'tools/kanban-static.sh' "$contract_file" || probe_fail "$contract_file does not exercise tools/kanban-static.sh"
contract_output="$(bash "$contract_file" 2>&1)" || probe_fail "board launcher contract failed: $contract_output"
grep -qxF 'queue-kanban contract probes passed.' <<<"$contract_output" \
  || probe_fail "board launcher contract did not report a pass: $contract_output"

printf 'REQ-717 probe passed.\n'
