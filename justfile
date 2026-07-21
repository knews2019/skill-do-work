# --- do-work board recipes (installed by `do-work install just-kanban`) ---

# Serve the do-work queue as a live Kanban board, replacing a stale instance on the port and opening your browser (Ctrl-C to stop; reload the page to refresh)
run-kanban $port="8090":
    case "$port" in ''|*[!0-9]*) echo "queue-kanban: invalid port '$port' - must be digits only (for a LAN-exposed host:port bind, run the queue-kanban serve command directly)" >&2; exit 1;; esac
    if command -v lsof >/dev/null 2>&1; then listener_pid="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"; if [ -n "$listener_pid" ]; then listener_command="$(ps -p "$listener_pid" -o args= 2>/dev/null)"; case "$listener_command" in *queue-kanban*) echo "queue-kanban: stopping previous session on :$port (pid $listener_pid): $listener_command"; kill "$listener_pid" 2>/dev/null; wait_count=0; while kill -0 "$listener_pid" 2>/dev/null && [ "$wait_count" -lt 20 ]; do sleep 0.1; wait_count=$((wait_count+1)); done;; *) echo "queue-kanban: port $port is already in use by another process ($listener_command, pid $listener_pid) - refusing to kill it. Stop it manually, or run 'just run-kanban <port>' with a different port." >&2; exit 1;; esac; fi; fi
    cd tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban serve --open --repo-root "{{justfile_directory()}}" --port "$port"

# Shareable static snapshot → build/queue-kanban-board/index.html (locally git-excluded so it never dirties git status)
kanban-static:
    cd tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban generate --out "{{justfile_directory()}}/build/queue-kanban-board" --repo-root "{{justfile_directory()}}"
    cd "{{justfile_directory()}}" && if git rev-parse --git-dir >/dev/null 2>&1 && ! git check-ignore -q build/queue-kanban-board/index.html; then exclude_file="$(git rev-parse --git-path info/exclude)"; mkdir -p "$(dirname "$exclude_file")"; echo '/build/queue-kanban-board/' >> "$exclude_file"; echo "kanban-static: added /build/queue-kanban-board/ to .git/info/exclude (local-only ignore)"; fi

# Column counts in the terminal, no browser
kanban-summary:
    cd tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban summary --repo-root "{{justfile_directory()}}"
