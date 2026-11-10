#!/usr/bin/env bash
# Pre-build baseline for REQ-717 (suite-managed justfile recipes call a script in one line).
# Every check here passes at the base revision and must still pass after the build:
#   - the board launcher contract (run-kanban stuck/foreign listener refusals),
#   - the folder-structure-check Go tests (the RED check this REQ clears),
#   - the live characterization of run-kanban and kanban-static in a throwaway consumer.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"

contract_output="$(bash _dev/tests/contracts/queue-kanban.sh 2>&1)" || {
  printf '%s\n' "$contract_output" >&2
  exit 1
}
grep -qxF 'queue-kanban contract probes passed.' <<<"$contract_output"

go_output="$(cd skills/do-work/tools/do-work-cli && go test -count=1 -v -run 'FolderStructure' ./internal/toolboxcommands/ 2>&1)" || {
  printf '%s\n' "$go_output" >&2
  exit 1
}
grep -qF -- '--- PASS: TestFolderStructure' <<<"$go_output"

bash do-work/runs/2026-11-09-layout-queue/REQ-717-live-proof.sh
printf 'REQ-717 pre-flight baseline passed.\n'
