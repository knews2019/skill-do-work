#!/usr/bin/env bash
# Start or consume a secret-quarantined uncommitted-file inventory.
set -u

if [ "$#" -lt 1 ] || [ "$#" -gt 3 ]; then
  printf 'Usage: %s (start|associate) [repository-root] [quarantine-name]\n' "$0" >&2
  exit 2
fi

inventory_mode="$1"
repository_root="${2:-$(git rev-parse --show-toplevel 2>/dev/null)}"
quarantine_name="${3:-do-work-commit-secret-quarantine}"
if [ -z "$repository_root" ]; then
  printf 'Repository root could not be resolved.\n' >&2
  exit 2
fi
script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
checks_directory="$(cd "$script_directory/../tools/checks" && pwd)"
case "$quarantine_name" in
  ''|*/*) printf 'Quarantine name must be one Git-private basename.\n' >&2; exit 2 ;;
esac
quarantine_paths_file="$(git -C "$repository_root" rev-parse --git-path "$quarantine_name")" || exit 2
inventory_file="$(mktemp)" || exit 2
candidate_paths_file=""

cleanup_inventory() {
  rm -f "$inventory_file"
  if [ -n "$candidate_paths_file" ]; then
    rm -f "$candidate_paths_file"
  fi
}
trap cleanup_inventory EXIT HUP INT TERM

if "$checks_directory/uncommitted-inventory.sh" "$repository_root" > "$inventory_file"; then
  inventory_status=0
else
  inventory_status=$?
fi
case "$inventory_status" in
  0) ;;
  1)
    [ "$inventory_mode" != "start" ] || rm -f "$quarantine_paths_file"
    exit 1
    ;;
  *) exit 2 ;;
esac

case "$inventory_mode" in
  start)
    awk -F '\t' '$1 == "X" { sub(/^[^\t]*\t/, ""); print }' \
      "$inventory_file" > "$quarantine_paths_file" || exit 2
    cat "$inventory_file"
    ;;
  associate)
    test -f "$quarantine_paths_file" || {
      printf 'Protected inventory has not been started.\n' >&2
      exit 2
    }
    candidate_paths_file="$(mktemp)" || exit 2
    awk -F '\t' '$1 == "X" { sub(/^[^\t]*\t/, ""); print }' \
      "$inventory_file" >> "$quarantine_paths_file" || exit 2
    awk -F '\t' '
      FILENAME == ARGV[1] { excluded[$0] = 1; next }
      {
        tag = $1
        sub(/^[^\t]*\t/, "")
        if (tag != "X" && !($0 in excluded)) print
      }
    ' "$quarantine_paths_file" "$inventory_file" > "$candidate_paths_file" || exit 2
    "$checks_directory/associate-files.sh" --repo-root "$repository_root" < "$candidate_paths_file"
    ;;
  *)
    printf 'Unknown protected inventory mode: %s\n' "$inventory_mode" >&2
    exit 2
    ;;
esac
