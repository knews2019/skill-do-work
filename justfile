# >>> do-work:recipes >>>
# --- board recipes (managed by the full-suite installer; run with `do-work-board board`) ---

# Serve the do-work queue as a live Kanban board, replacing a stale instance on the port and opening your browser (Ctrl-C to stop; reload the page to refresh)
run-kanban $port="8090":
    case "$port" in ''|*[!0-9]*) echo "queue-kanban: invalid port '$port' - must be digits only (for a LAN-exposed host:port bind, run the queue-kanban serve command directly)" >&2; exit 1;; esac
    if command -v lsof >/dev/null 2>&1; then listener_pid="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"; if [ -n "$listener_pid" ]; then listener_executable="$(lsof -a -p "$listener_pid" -d txt -Fn 2>/dev/null | sed -n 's/^n//p' | head -n1)"; listener_executable_name="${listener_executable##*/}"; listener_command="$(ps -p "$listener_pid" -o args= 2>/dev/null)"; case "$listener_executable_name" in *queue-kanban*) echo "queue-kanban: stopping previous session on :$port (pid $listener_pid): $listener_command"; kill "$listener_pid" 2>/dev/null; wait_count=0; while [ "$wait_count" -lt 320 ] && lsof -a -p "$listener_pid" -i tcp:"$port" -sTCP:LISTEN -t >/dev/null 2>&1; do sleep 0.1; wait_count=$((wait_count+1)); done;; *) echo "queue-kanban: port $port is already in use by another process ($listener_command, pid $listener_pid) - refusing to kill it. Stop it manually, or run 'just run-kanban <port>' with a different port." >&2; exit 1;; esac; fi; remaining_listener_pid="$(lsof -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null | head -n1)"; if [ -n "$remaining_listener_pid" ]; then remaining_listener_command="$(ps -p "$remaining_listener_pid" -o args= 2>/dev/null)"; echo "queue-kanban: port $port still has a listener after shutdown wait (pid $remaining_listener_pid): $remaining_listener_command - refusing to start." >&2; exit 1; fi; fi
    cd skills/do-work-board/tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban serve --open --repo-root "{{justfile_directory()}}" --port "$port"

# Fast terminal read of what's in flight — open count, every claimed REQ with its title, every needs-input/blocked REQ with its status (no browser, no server)
run-kanban-cli:
    cd skills/do-work-board/tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban open-work --repo-root "{{justfile_directory()}}"

# Shareable static snapshot → build/queue-kanban-board/index.html (locally git-excluded so it never dirties git status)
kanban-static:
    cd skills/do-work-board/tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban generate --out "{{justfile_directory()}}/build/queue-kanban-board" --repo-root "{{justfile_directory()}}"
    cd "{{justfile_directory()}}" && if git rev-parse --git-dir >/dev/null 2>&1 && ! git check-ignore -q build/queue-kanban-board/index.html; then exclude_file="$(git rev-parse --git-path info/exclude)"; mkdir -p "$(dirname "$exclude_file")"; echo '/build/queue-kanban-board/' >> "$exclude_file"; echo "kanban-static: added /build/queue-kanban-board/ to .git/info/exclude (local-only ignore)"; fi

# Column counts in the terminal, no browser
kanban-summary:
    cd skills/do-work-board/tools/queue-kanban && go build -o queue-kanban . && ./queue-kanban summary --repo-root "{{justfile_directory()}}"

# Update the project-local do-work suite without an agent (one reviewed archive, one confirmation, managed-path recovery)
run-do-work-update:
    project_root="{{justfile_directory()}}"; skill_root="$project_root/.claude/skills/do-work"; [ -f "$skill_root/SKILL.md" ] || skill_root="$project_root/skills/do-work"; bash "$skill_root/tools/do-work-update.sh" --project-root "$project_root"

# Initialize or fill gaps in a BKB scaffold through the canonical Go command
bkb-init kb="kb" *args:
    project_root="{{justfile_directory()}}"; skill_root="$project_root/.claude/skills/do-work"; [ -f "$skill_root/SKILL.md" ] || skill_root="$project_root/skills/do-work"; bash "$skill_root/tools/do-work-cli.sh" --repo-root "$project_root" bkb-init --kb {{quote(kb)}} {{args}}

# Read a deterministic BKB status snapshot
bkb-status kb="kb":
    project_root="{{justfile_directory()}}"; skill_root="$project_root/.claude/skills/do-work"; [ -f "$skill_root/SKILL.md" ] || skill_root="$project_root/skills/do-work"; bash "$skill_root/tools/do-work-cli.sh" --repo-root "$project_root" bkb-status --kb {{quote(kb)}}

# Run deterministic BKB structural lint without semantic edits
bkb-lint-structure kb="kb":
    project_root="{{justfile_directory()}}"; skill_root="$project_root/.claude/skills/do-work"; [ -f "$skill_root/SKILL.md" ] || skill_root="$project_root/skills/do-work"; bash "$skill_root/tools/do-work-cli.sh" --repo-root "$project_root" bkb-lint-structure --kb {{quote(kb)}}

# Run Dream's seven read-only deterministic scans
dream-scan path="memory":
    project_root="{{justfile_directory()}}"; skill_root="$project_root/.claude/skills/do-work"; [ -f "$skill_root/SKILL.md" ] || skill_root="$project_root/skills/do-work"; bash "$skill_root/tools/do-work-cli.sh" --repo-root "$project_root" dream-scan --path {{quote(path)}}
# <<< do-work:recipes <<<

# Repository-only verification; the script owns the command inventory.
maintainer-verify:
    bash "{{justfile_directory()}}/_dev/tests/maintainer-verify.sh"
