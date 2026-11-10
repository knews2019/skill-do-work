#!/usr/bin/env bash
# The logic of the `kanban-static` just recipe: build queue-kanban, generate the static board into
# <project-root>/build/queue-kanban-board, and add that folder to the local Git exclude file once.
# Usage: kanban-static.sh <project-root>
# The build and the generate stand on their own lines so `set -e` stops the script when either
# fails, before the exclude write. Inside an && list, set -e ignores every command but the last.
set -eu
project_root="${1:?usage: kanban-static.sh <project-root>}"
# Resolved first: bash 3.2 loses BASH_SOURCE after it calls a function imported from the environment.
queue_kanban_directory="$(dirname "${BASH_SOURCE[0]}")/queue-kanban"

cd "$queue_kanban_directory"
go build -o queue-kanban .
./queue-kanban generate --out "$project_root/build/queue-kanban-board" --repo-root "$project_root"
cd "$project_root"
if git rev-parse --git-dir >/dev/null 2>&1 && ! git check-ignore -q build/queue-kanban-board/index.html; then
  exclude_file="$(git rev-parse --git-path info/exclude)"
  mkdir -p "$(dirname "$exclude_file")"
  echo '/build/queue-kanban-board/' >> "$exclude_file"
  echo "kanban-static: added /build/queue-kanban-board/ to .git/info/exclude (local-only ignore)"
fi
