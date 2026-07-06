# --- do-work board recipes (installed by `do-work install just-kanban`) ---

# Serve the do-work queue as a live Kanban board, replacing a stale instance on the port and opening your browser (Ctrl-C to stop; reload the page to refresh)
run-kanban $port="8090":
    case "$port" in ''|*[!0-9]*) echo "queue-kanban: invalid port '$port' - must be digits only (for a LAN-exposed host:port bind, run the queue-kanban serve command directly)" >&2; exit 1;; esac
    if command -v lsof >/dev/null 2>&1; then PID="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"; if [ -n "$PID" ]; then COMM="$(ps -p "$PID" -o comm= 2>/dev/null)"; COMM="${COMM##*/}"; if [ "$COMM" = "queue-kanban" ]; then kill "$PID" 2>/dev/null; i=0; while kill -0 "$PID" 2>/dev/null && [ "$i" -lt 20 ]; do sleep 0.1; i=$((i+1)); done; else echo "queue-kanban: port $port is already in use by another process ($COMM, pid $PID) - refusing to kill it. Stop it manually, or run 'just run-kanban <port>' with a different port." >&2; exit 1; fi; fi; fi
    cd tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban serve --open --repo-root "{{justfile_directory()}}" --port "$port"

# Shareable static snapshot → build/queue-kanban-board/index.html
kanban-static:
    cd tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban generate --out "{{justfile_directory()}}/build/queue-kanban-board" --repo-root "{{justfile_directory()}}"

# Column counts in the terminal, no browser
kanban-summary:
    cd tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban summary --repo-root "{{justfile_directory()}}"
