#!/usr/bin/env bash
# Parallel launcher for the aggregate contract suite's behavioral sub-suites. Sourced by
# _dev/tests/contract-regressions.sh; expects `fail_count` in the caller's scope.
#
#   launch_probe <name> <failure message> <script>...   runs the scripts in order, in the background
#   collect_probes                                       waits, prints outputs in launch order, counts failures

probe_batch_root="$(mktemp -d "${TMPDIR:-/tmp}/contract-probe-batch.XXXXXX")"
probe_batch_names=()
probe_file_budget_seconds="${DO_WORK_TEST_FILE_BUDGET_SECONDS:-30}"

launch_probe() {
  local probe_name="$1"
  local probe_failure_message="$2"
  shift 2
  printf '%s\n' "$probe_failure_message" > "$probe_batch_root/$probe_name.message"
  # Job control on for the launch: a plain `&` in a non-interactive shell starts the job
  # with SIGINT and SIGQUIT ignored, and the update-script probe's interrupt case sends INT.
  set -m
  (
    probe_status=0
    for probe_script in "$@"; do
      probe_started_at="$(date +%s)"
      current_probe_status=0
      bash "$probe_script" || current_probe_status=$?
      probe_elapsed_seconds=$(( $(date +%s) - probe_started_at ))
      printf 'test-file duration: %s %ss (limit <%ss)\n' \
        "${probe_script##*/}" "$probe_elapsed_seconds" "$probe_file_budget_seconds"
      if [ "$probe_elapsed_seconds" -ge "$probe_file_budget_seconds" ]; then
        printf 'FAIL: %s took %ss; each test file must finish under %ss\n' \
          "$probe_script" "$probe_elapsed_seconds" "$probe_file_budget_seconds" >&2
        current_probe_status=1
      fi
      if [ "$current_probe_status" -ne 0 ]; then
        probe_status=1
        break
      fi
    done
    printf '%s\n' "$probe_status" > "$probe_batch_root/$probe_name.status"
  ) > "$probe_batch_root/$probe_name.out" 2>&1 &
  set +m
  probe_batch_names+=("$probe_name")
}
collect_probes() {
  local probe_name
  local probe_status
  wait
  for probe_name in ${probe_batch_names[@]+"${probe_batch_names[@]}"}; do
    cat "$probe_batch_root/$probe_name.out"
    probe_status="$(cat "$probe_batch_root/$probe_name.status" 2>/dev/null || printf 'none')"
    if [ "$probe_status" != '0' ]; then
      printf 'FAIL: %s\n' "$(cat "$probe_batch_root/$probe_name.message")" >&2
      fail_count=$((fail_count + 1))
    fi
  done
  rm -rf -- "$probe_batch_root"
}
