#!/usr/bin/env bash
# Fixture execution proofs for run-blocked-check.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# run-blocked-check: force the stock-Bash fallback, then prove timeout owns the
# isolated wrapper/descendant group without touching an unrelated group member.
fallback_bin="$fixture_root/fallback-bin"
mkdir -p "$fallback_bin"
ln -s "$(command -v bash)" "$fallback_bin/bash"
ln -s "$(command -v sh)" "$fallback_bin/sh"
ln -s "$(command -v sleep)" "$fallback_bin/sleep"
sleep 30 & unrelated_process_id=$!
background_process_ids="$background_process_ids $unrelated_process_id"
printf '%s\n' \
  'trap "" TERM' \
  'printf "%s\n" "$$" > "$BLOCKED_WRAPPER_PID_FILE"' \
  '(trap "" TERM; sleep 30) &' \
  'printf "%s\n" "$!" > "$BLOCKED_DESCENDANT_PID_FILE"' \
  'wait' > "$fixture_root/blocked-probe.sh"
PATH="$fallback_bin:$PATH" \
  BLOCKED_WRAPPER_PID_FILE="$fixture_root/blocked-wrapper.pid" \
  BLOCKED_DESCENDANT_PID_FILE="$fixture_root/blocked-descendant.pid" \
  "$core_scripts/run-blocked-check.sh" "$fixture_root/blocked-probe.sh" 1 >/dev/null 2>&1
blocked_status=$?
[ "$blocked_status" -eq 124 ] || fail_case "run-blocked-check portable-timeout case returned $blocked_status instead of 124"
blocked_wrapper_pid="$(cat "$fixture_root/blocked-wrapper.pid")"
blocked_descendant_pid="$(cat "$fixture_root/blocked-descendant.pid")"
background_process_ids="$background_process_ids $blocked_wrapper_pid $blocked_descendant_pid"
# kill -0 counts zombies as alive; a killed-but-unreaped descendant (reparent target
# reaps lazily, e.g. container PID 1) must read as dead, so check the process state.
process_runs_unreaped_excluded() {
  local process_state
  process_state="$(ps -o stat= -p "$1" 2>/dev/null | tr -d '[:space:]')" || return 1
  [ -n "$process_state" ] || return 1
  case "$process_state" in Z*) return 1 ;; esac
  return 0
}
process_runs_unreaped_excluded "$blocked_wrapper_pid" && fail_case 'run-blocked-check process-tree case left the wrapper alive'
process_runs_unreaped_excluded "$blocked_descendant_pid" && fail_case 'run-blocked-check process-tree case left the descendant alive'
kill -0 "$unrelated_process_id" 2>/dev/null || fail_case 'run-blocked-check process-tree cleanup killed an unrelated process in the test runner group'
: > "$fixture_root/test-runner-survived-timeout"
[ -e "$fixture_root/test-runner-survived-timeout" ] || fail_case 'run-blocked-check process-tree cleanup killed the test runner group'

# run-blocked-check: the fallback must preserve an ordinary probe status.
printf 'exit 23\n' > "$fixture_root/blocked-status-probe.sh"
PATH="$fallback_bin:$PATH" "$core_scripts/run-blocked-check.sh" "$fixture_root/blocked-status-probe.sh" 2 >/dev/null 2>&1
blocked_ordinary_status=$?
[ "$blocked_ordinary_status" -eq 23 ] || fail_case "run-blocked-check ordinary-status case returned $blocked_ordinary_status instead of 23"

prescribed_shell_finish
