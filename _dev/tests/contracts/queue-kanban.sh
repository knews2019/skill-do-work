#!/usr/bin/env bash
# Runtime pins for the board launcher shutdown boundary.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
fail_count=0
extract_kanban_shutdown_line() {
  awk '
    /^run-kanban \$port=/ {
      getline
      getline
      sub(/^[[:space:]]*/, "")
      print
      exit
    }
  ' "$repo_root/$1"
}

root_kanban_shutdown_line="$(extract_kanban_shutdown_line justfile)"
installer_kanban_shutdown_line="$(extract_kanban_shutdown_line skills/do-work-board/justfile.template)"
if [ "$root_kanban_shutdown_line" != "$installer_kanban_shutdown_line" ]; then
  printf 'FAIL: Justfile and the board-owned template must carry one identical run-kanban shutdown line.\n' >&2
  fail_count=$((fail_count + 1))
fi
# Execute the canonical shutdown line with command seams. A queue-kanban PID that remains a
# listener throughout the bounded wait must make the recipe line fail before build+serve, and
# the diagnosis must name both the PID and full command. sleep is a no-op in the fixture, so all
# 320 iterations run without slowing the suite.
if stuck_listener_output="$(
  port=8090
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/tmp/queue-kanban\n'
    else
      printf '4242\n'
    fi
  }
  ps() { printf '/tmp/queue-kanban serve --port 8090\n'; }
  kill() { return 0; }
  sleep() { return 0; }
  eval "$root_kanban_shutdown_line" 2>&1
)"; then
  stuck_listener_status=0
else
  stuck_listener_status=$?
fi
if [ "$stuck_listener_status" -ne 1 ]; then
  printf 'FAIL: run-kanban shutdown must refuse startup when a listener remains after the bounded wait; got exit %s.\n' "$stuck_listener_status" >&2
  fail_count=$((fail_count + 1))
fi
if ! printf '%s\n' "$stuck_listener_output" | grep -qF 'pid 4242'; then
  printf 'FAIL: run-kanban stuck-listener refusal must name pid 4242.\n' >&2
  fail_count=$((fail_count + 1))
fi
if ! printf '%s\n' "$stuck_listener_output" | grep -qF '/tmp/queue-kanban serve --port 8090'; then
  printf 'FAIL: run-kanban stuck-listener refusal must name the listener command.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Preserve the older safety boundary: a foreign executable is refused immediately and never
# passed to kill. This behavior probe complements the executable-identity contract above.
foreign_kill_marker="$(mktemp)"
if foreign_listener_output="$(
  port=8090
  # The shutdown recipe reads this through eval, which ShellCheck cannot trace.
  : "$port"
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/usr/bin/python3\n'
    else
      printf '3131\n'
    fi
  }
  ps() { printf '/usr/bin/python3 -m http.server 8090\n'; }
  kill() { printf 'called\n' > "$foreign_kill_marker"; }
  sleep() { return 0; }
  eval "$root_kanban_shutdown_line" 2>&1
)"; then
  foreign_listener_status=0
else
  foreign_listener_status=$?
fi
# This fixture asserts the refusal status and kill boundary; its captured text is intentionally ignored.
: "$foreign_listener_output"
if [ "$foreign_listener_status" -ne 1 ] || [ -s "$foreign_kill_marker" ]; then
  printf 'FAIL: run-kanban must refuse a foreign listener without calling kill.\n' >&2
  fail_count=$((fail_count + 1))
fi
rm -f "$foreign_kill_marker"
[ "$fail_count" -eq 0 ] || exit 1
printf 'queue-kanban contract probes passed.\n'
