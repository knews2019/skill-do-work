#!/usr/bin/env bash
# The logic of the `run-kanban` just recipe: validate the port, stop a stale queue-kanban listener
# on it, then build and serve the live board for the project.
# Usage: run-kanban.sh <project-root> <port>
# Only `set -u`, the same as just's `sh -cu`: under -e or pipefail the lsof pipelines below exit 1
# when nothing listens, and the script would die before it serves.
set -u
project_root="${1:?usage: run-kanban.sh <project-root> <port>}"
port="${2?usage: run-kanban.sh <project-root> <port>}"
# Resolved first: bash 3.2 loses BASH_SOURCE after it calls a function imported from the environment.
queue_kanban_directory="$(dirname "${BASH_SOURCE[0]}")/queue-kanban"

case "$port" in ''|*[!0-9]*) echo "queue-kanban: invalid port '$port' - must be digits only (for a LAN-exposed host:port bind, run the queue-kanban serve command directly)" >&2; exit 1;; esac
if command -v lsof >/dev/null 2>&1; then
  listener_pid="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"
  if [ -n "$listener_pid" ]; then
    listener_executable="$(lsof -a -p "$listener_pid" -d txt -Fn 2>/dev/null | sed -n 's/^n//p' | head -n1)"
    listener_executable_name="${listener_executable##*/}"
    listener_command="$(ps -p "$listener_pid" -o args= 2>/dev/null)"
    case "$listener_executable_name" in
      *queue-kanban*)
        echo "queue-kanban: stopping previous session on :$port (pid $listener_pid): $listener_command"
        kill "$listener_pid" 2>/dev/null
        wait_count=0
        while [ "$wait_count" -lt 320 ] && lsof -a -p "$listener_pid" -i tcp:"$port" -sTCP:LISTEN -t >/dev/null 2>&1; do sleep 0.1; wait_count=$((wait_count+1)); done;;
      *)
        echo "queue-kanban: port $port is already in use by another process ($listener_command, pid $listener_pid) - refusing to kill it. Stop it manually, or run 'just run-kanban <port>' with a different port." >&2
        exit 1;;
    esac
  fi
  remaining_listener_pid="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"
  if [ -n "$remaining_listener_pid" ]; then
    remaining_listener_command="$(ps -p "$remaining_listener_pid" -o args= 2>/dev/null)"
    echo "queue-kanban: port $port still has a listener after shutdown wait (pid $remaining_listener_pid): $remaining_listener_command - refusing to start." >&2
    exit 1
  fi
fi
cd "$queue_kanban_directory" && go build -o queue-kanban . && ./queue-kanban serve --open --repo-root "$project_root" --port "$port"
