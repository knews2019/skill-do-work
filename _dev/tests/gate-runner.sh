#!/usr/bin/env bash
# Background scheduler for the maintainer gate. Runs `_dev/tests/maintainer-verify.sh`
# once per new HEAD and records a green result through `do-work-cli record-green-gate`,
# so a pipeline claim that finds HEAD already proven green skips its baseline run
# (`actions/work.md` Step 5.75 consults `check-green-gate`). The gate stops blocking
# every step and becomes evidence attached to a revision.
#
#   bash _dev/tests/gate-runner.sh            # watch HEAD forever, 20 s between polls
#   bash _dev/tests/gate-runner.sh --once     # run for the current HEAD and exit with its status
#   GATE_RUNNER_INTERVAL=60 bash _dev/tests/gate-runner.sh
#
# Logs land under $TMPDIR/do-work-gate-runs/<revision>.log. One line per run is printed:
#   gate <revision> green|red <seconds>s <log path> [recorded|not recorded: <reason>]
set -euo pipefail

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"
gate_script="$repo_root/_dev/tests/maintainer-verify.sh"
cli_launcher="$repo_root/skills/do-work/tools/do-work-cli.sh"
log_root="${TMPDIR:-/tmp}/do-work-gate-runs"
poll_interval="${GATE_RUNNER_INTERVAL:-20}"
run_once='no'
if [ "${1:-}" = '--once' ]; then
  run_once='yes'
elif [ "$#" -ne 0 ]; then
  printf 'Usage: %s [--once]\n' "$0" >&2
  exit 2
fi
mkdir -p "$log_root"

run_gate_for_head() {
  local revision="$1"
  local log_path="$log_root/$revision.log"
  local started_at
  local gate_status=0
  local record_note
  started_at="$(date +%s)"
  if ! (cd "$repo_root" && bash "$gate_script" > "$log_path" 2>&1); then
    gate_status=1
  fi
  if [ "$gate_status" -eq 0 ]; then
    if record_output="$(cd "$repo_root" && bash "$cli_launcher" --repo-root "$repo_root" --format json \
      record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh 2>&1)"; then
      record_note='recorded'
    else
      record_note="not recorded: $(printf '%s' "$record_output" | tr '\n' ' ' | cut -c1-160)"
    fi
    printf 'gate %s green %ss %s %s\n' "$revision" "$(( $(date +%s) - started_at ))" "$log_path" "$record_note"
  else
    printf 'gate %s red %ss %s\n' "$revision" "$(( $(date +%s) - started_at ))" "$log_path"
  fi
  return "$gate_status"
}

last_revision=''
while :; do
  head_revision="$(git -C "$repo_root" rev-parse HEAD)"
  if [ "$head_revision" != "$last_revision" ]; then
    last_revision="$head_revision"
    gate_result=0
    run_gate_for_head "$head_revision" || gate_result=$?
    if [ "$run_once" = 'yes' ]; then
      exit "$gate_result"
    fi
  fi
  sleep "$poll_interval"
done
