#!/usr/bin/env bash
# archive-collision.sh — mechanical form of actions/work.md Step 2.0.
# Checks whether a REQ id is already archived before the work loop claims it.
#
# Usage: tools/checks/archive-collision.sh REQ-042
# Exit 0: no collision (safe to claim). Prints nothing.
# Exit 1: collision — prints every matching archive path, one per line.
# Exit 2: usage error.
#
# Run from the project root (the directory containing do-work/). If
# do-work/archive/ does not exist there is nothing to collide with — exit 0.
set -uo pipefail

request_id="${1:-}"
if [[ ! "$request_id" =~ ^REQ-[0-9]+$ ]]; then
  echo "usage: $0 REQ-NNN" >&2
  exit 2
fi

archive_root="do-work/archive"
[ -d "$archive_root" ] || exit 0

collision_paths="$(find "$archive_root" \( -name "${request_id}-*.md" -o -name "${request_id}.md" \) 2>/dev/null)"

if [ -n "$collision_paths" ]; then
  printf '%s\n' "$collision_paths"
  exit 1
fi
exit 0
