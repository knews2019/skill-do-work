#!/usr/bin/env bash
# Shared append-only timing record for shell and Go test files.

duration_log_path="${DO_WORK_TEST_DURATION_LOG:-${repo_root:?}/do-work/test-durations.tsv}"
duration_run_id="${DO_WORK_TEST_RUN_ID:-standalone-$(date -u +%Y%m%dT%H%M%SZ)-$$}"
duration_other_gate_processes="${DO_WORK_TEST_OTHER_GATE_PROCESSES:-0}"
duration_log_header="run_id\tfile\tseconds\tother_gate_processes"
export DO_WORK_TEST_DURATION_LOG="$duration_log_path"
export DO_WORK_TEST_RUN_ID="$duration_run_id"
export DO_WORK_TEST_OTHER_GATE_PROCESSES="$duration_other_gate_processes"

if [ ! -e "$duration_log_path" ]; then
  duration_header_candidate="$duration_log_path.header.$$"
  printf '%b\n' "$duration_log_header" > "$duration_header_candidate"
  ln "$duration_header_candidate" "$duration_log_path" 2>/dev/null || true
  rm -f -- "$duration_header_candidate"
fi
if [ "$(sed -n '1p' "$duration_log_path" 2>/dev/null)" != "$(printf '%b' "$duration_log_header")" ]; then
  printf 'FAIL: test duration log has an invalid header: %s\n' "$duration_log_path" >&2
  exit 1
fi

record_test_duration() {
  local test_file="$1"
  local elapsed_seconds="$2"
  printf '%s\t%s\t%s\t%s\n' \
    "$duration_run_id" "$test_file" "$elapsed_seconds" "$duration_other_gate_processes" \
    >> "$duration_log_path"
}
