#!/usr/bin/env bash
# Run one already-materialized blocked-condition probe with a portable timeout.
set -u

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
  printf 'Usage: %s <probe-script> [timeout-seconds]\n' "$0" >&2
  exit 2
fi

probe_script="$1"
timeout_seconds="${2:-30}"
case "$timeout_seconds" in
  ''|*[!0-9]*) printf 'Timeout must be a positive integer.\n' >&2; exit 2 ;;
esac
if [ "$timeout_seconds" -eq 0 ] || [ ! -r "$probe_script" ]; then
  printf 'Probe script must be readable and timeout must be positive.\n' >&2
  exit 2
fi

if command -v timeout >/dev/null 2>&1; then
  timeout "$timeout_seconds" sh "$probe_script"
  exit $?
fi
if command -v gtimeout >/dev/null 2>&1; then
  gtimeout "$timeout_seconds" sh "$probe_script"
  exit $?
fi

if [ -x /bin/ps ]; then
  process_status_command=/bin/ps
elif [ -x /usr/bin/ps ]; then
  process_status_command=/usr/bin/ps
else
  printf 'Portable timeout requires ps to verify process-group isolation.\n' >&2
  exit 125
fi

probe_process_id=""
probe_process_group_id=""
caller_process_group_id="$($process_status_command -o pgid= -p "$$" 2>/dev/null || true)"
caller_process_group_id="${caller_process_group_id//[[:space:]]/}"

probe_group_is_safe() {
  case "$probe_process_id:$probe_process_group_id:$caller_process_group_id" in
    *[!0-9:]*|:*|*:|*::*) return 1 ;;
  esac
  [ "$probe_process_group_id" = "$probe_process_id" ] \
    && [ "$probe_process_group_id" != "$caller_process_group_id" ]
}

terminate_probe_group() {
  probe_signal="${1:-TERM}"
  probe_grace_ticks=0
  probe_grace_limit=10

  probe_group_is_safe || return 0
  kill -"$probe_signal" -- "-$probe_process_group_id" 2>/dev/null || true
  while kill -0 -- "-$probe_process_group_id" 2>/dev/null \
    && [ "$probe_grace_ticks" -lt "$probe_grace_limit" ]; do
    sleep 0.1
    probe_grace_ticks=$((probe_grace_ticks + 1))
  done
  if kill -0 -- "-$probe_process_group_id" 2>/dev/null; then
    kill -KILL -- "-$probe_process_group_id" 2>/dev/null || true
  fi
  wait "$probe_process_id" 2>/dev/null || true
  probe_reap_ticks=0
  while kill -0 -- "-$probe_process_group_id" 2>/dev/null \
    && [ "$probe_reap_ticks" -lt "$probe_grace_limit" ]; do
    sleep 0.1
    probe_reap_ticks=$((probe_reap_ticks + 1))
  done
}

trap 'terminate_probe_group TERM; exit 129' HUP
trap 'terminate_probe_group TERM; exit 130' INT
trap 'terminate_probe_group TERM; exit 143' TERM

# Job control gives the background wrapper its own process group. The wrapper stops
# before executing the probe so the parent can verify isolation before user code runs.
set -m
sh -c 'kill -STOP "$$"; exec sh "$1"' blocked-check-wrapper "$probe_script" &
probe_process_id=$!
set +m

probe_launch_ticks=0
probe_launch_limit=100
probe_process_state=""
while kill -0 "$probe_process_id" 2>/dev/null \
  && [ "$probe_launch_ticks" -lt "$probe_launch_limit" ]; do
  probe_process_state="$($process_status_command -o stat= -p "$probe_process_id" 2>/dev/null || true)"
  probe_process_state="${probe_process_state//[[:space:]]/}"
  case "$probe_process_state" in
    T*) break ;;
  esac
  sleep 0.01
  probe_launch_ticks=$((probe_launch_ticks + 1))
done

probe_process_group_id="$($process_status_command -o pgid= -p "$probe_process_id" 2>/dev/null || true)"
probe_process_group_id="${probe_process_group_id//[[:space:]]/}"
case "$probe_process_state" in
  T*) ;;
  *)
    kill -KILL "$probe_process_id" 2>/dev/null || true
    wait "$probe_process_id" 2>/dev/null || true
    printf 'Portable timeout could not stop the probe before launch.\n' >&2
    exit 125
    ;;
esac
if ! probe_group_is_safe; then
  kill -KILL "$probe_process_id" 2>/dev/null || true
  wait "$probe_process_id" 2>/dev/null || true
  printf 'Portable timeout could not establish an isolated process group.\n' >&2
  exit 125
fi
if ! kill -CONT -- "-$probe_process_group_id" 2>/dev/null; then
  terminate_probe_group KILL
  printf 'Portable timeout could not start the isolated probe group.\n' >&2
  exit 125
fi

probe_wait_ticks=0
probe_tick_limit=$((timeout_seconds * 10))
while kill -0 "$probe_process_id" 2>/dev/null && [ "$probe_wait_ticks" -lt "$probe_tick_limit" ]; do
  sleep 0.1
  probe_wait_ticks=$((probe_wait_ticks + 1))
done
if kill -0 "$probe_process_id" 2>/dev/null; then
  terminate_probe_group TERM
  exit 124
fi
wait "$probe_process_id"
probe_status=$?
probe_process_id=""
probe_process_group_id=""
exit "$probe_status"
