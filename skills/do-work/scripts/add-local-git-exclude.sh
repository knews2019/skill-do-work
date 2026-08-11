#!/usr/bin/env bash
# Add one worktree-safe pattern to Git's local exclude file.
set -u

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
  printf 'Usage: %s <path-to-probe> [exclude-pattern]\n' "$0" >&2
  exit 2
fi

probe_path="$1"
exclude_pattern="${2:-**/${probe_path#./}}"
case "$probe_path$exclude_pattern" in
  *"
"*) printf 'Paths and patterns must not contain newlines.\n' >&2; exit 2 ;;
esac

exclude_file="$(git rev-parse --git-path info/exclude 2>/dev/null)" || {
  printf 'Local Git exclude unavailable outside a repository.\n' >&2
  exit 0
}

if git check-ignore -q -- "$probe_path" 2>/dev/null; then
  printf 'present\t%s\n' "$exclude_pattern"
  exit 0
fi

printf '%s\n' "$exclude_pattern" >> "$exclude_file" || {
  printf 'Could not append to local Git exclude: %s\n' "$exclude_file" >&2
  exit 2
}
printf 'added\t%s\n' "$exclude_pattern"
