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

sh "$probe_script" &
probe_process_id=$!
probe_wait_ticks=0
probe_tick_limit=$((timeout_seconds * 10))
while kill -0 "$probe_process_id" 2>/dev/null && [ "$probe_wait_ticks" -lt "$probe_tick_limit" ]; do
  sleep 0.1
  probe_wait_ticks=$((probe_wait_ticks + 1))
done
if kill -0 "$probe_process_id" 2>/dev/null; then
  kill "$probe_process_id" 2>/dev/null || true
  sleep 0.1
  kill -9 "$probe_process_id" 2>/dev/null || true
  wait "$probe_process_id" 2>/dev/null || true
  exit 124
fi
wait "$probe_process_id"
