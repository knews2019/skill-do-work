#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fail_count=0
verification_tier="${DO_WORK_MAINTAINER_TIER:-fast}"
test_file_budget_seconds="${DO_WORK_TEST_FILE_BUDGET_SECONDS:-30}"
# Ratchet baseline: raise only when aggregate responsibility deliberately grows.
fast_contract_line_ceiling=77

if [ "$verification_tier" != fast ] && [ "$verification_tier" != heavy ]; then
  printf 'FAIL: unsupported maintainer tier: %s\n' "$verification_tier" >&2
  exit 2
fi

actual_contract_lines="$(wc -l < "${BASH_SOURCE[0]}" | tr -d ' ')"
if [ "$actual_contract_lines" -gt "$fast_contract_line_ceiling" ]; then
  printf 'FAIL: contract-regressions.sh grew to %s lines; ratchet ceiling is %s.\n' \
    "$actual_contract_lines" "$fast_contract_line_ceiling" >&2
  fail_count=$((fail_count + 1))
fi

# shellcheck source=_dev/tests/test-duration-log.sh
source "$repo_root/_dev/tests/test-duration-log.sh"

run_contract_file() {
  local contract_file="$1"
  local relative_file="${contract_file#"$repo_root/"}"
  local started_at
  local elapsed_seconds
  local contract_status=0
  local budget_label="<${test_file_budget_seconds}s"
  if [ "$verification_tier" = heavy ]; then
    budget_label='none (heavy)'
  fi
  started_at="$(date +%s)"
  (
    # shellcheck disable=SC1090
    source "$contract_file"
  ) || contract_status=$?
  elapsed_seconds=$(( $(date +%s) - started_at ))
  printf 'test-file duration: %s %ss (limit %s)\n' \
    "$relative_file" "$elapsed_seconds" "$budget_label"
  record_test_duration "$relative_file" "$elapsed_seconds"
  if [ "$contract_status" -ne 0 ]; then
    fail_count=$((fail_count + 1))
  fi
  if [ "$verification_tier" != heavy ] \
    && [ "$elapsed_seconds" -ge "$test_file_budget_seconds" ]; then
    printf 'FAIL: %s took %ss; each fast test file must finish under %ss\n' \
      "$relative_file" "$elapsed_seconds" "$test_file_budget_seconds" >&2
    fail_count=$((fail_count + 1))
  fi
}

for contract_file in \
  "$repo_root/_dev/tests/contracts/core-checks.sh" \
  "$repo_root/_dev/tests/contracts/queue-kanban.sh" \
  "$repo_root/_dev/tests/contracts/replace-text-section.sh" \
  "$repo_root/_dev/tests/contracts/recovery-set-aside.sh"
do
  if [ ! -f "$contract_file" ]; then
    printf 'FAIL: missing owner contract: %s\n' "${contract_file#"$repo_root/"}" >&2
    fail_count=$((fail_count + 1))
    continue
  fi
  run_contract_file "$contract_file"
done

# shellcheck source=_dev/tests/probe-batch.sh
source "$repo_root/_dev/tests/probe-batch.sh"
# shellcheck source=_dev/tests/contracts/probe-lanes.sh
source "$repo_root/_dev/tests/contracts/probe-lanes.sh"
collect_probes

[ "$fail_count" -eq 0 ] || exit 1
printf 'Contract regression checks passed.\n'
