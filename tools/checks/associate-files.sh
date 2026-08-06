#!/usr/bin/env bash
# associate-files.sh — mechanical form of the Step 3 association pass shared by
# actions/commit.md and actions/inspect.md: match candidate paths against the
# "## Implementation Summary" file lists of archived and in-flight REQs.
#
# Usage: tools/checks/associate-files.sh [--repo-root DIR] < candidate-paths
#        candidate paths arrive one per line on stdin (the "<path>" column of
#        tools/checks/uncommitted-inventory.sh output)
# Output: one TAB-separated row per candidate — "<owner>\t<path>"
#           REQ-NNN  the REQ whose Implementation Summary claims this path
#           -        unassociated; the caller groups it semantically instead
# Exit 0: rows emitted (associated or not). Exit 1: no candidates on stdin.
# Exit 2: usage error, or no do-work/ directory (the caller skips REQ tracing).
#
# Read-only. This script resolves ownership and prints it; grouping, commit
# messages, and verdicts stay with the caller.
set -uo pipefail

repository_root="."
while [ $# -gt 0 ]; do
  case "$1" in
    --repo-root) repository_root="${2:-}"; shift 2 ;;
    *) echo "usage: $0 [--repo-root DIR] < candidate-paths" >&2; exit 2 ;;
  esac
done
if [ ! -d "$repository_root" ]; then
  echo "usage: $0 [--repo-root DIR] < candidate-paths" >&2
  exit 2
fi
cd "$repository_root" || exit 2

if [ ! -d "do-work" ]; then
  echo "NO-DO-WORK-DIR: nothing to associate against" >&2
  exit 2
fi

mapfile -t candidate_paths
# Drop blank lines so a trailing newline on stdin does not become a candidate.
filtered_candidates=()
for candidate_path in "${candidate_paths[@]}"; do
  [ -n "$candidate_path" ] && filtered_candidates+=("$candidate_path")
done
if [ "${#filtered_candidates[@]}" -eq 0 ]; then
  exit 1
fi

# Terminal-success per actions/work-reference.md's Schema Read Contract. The
# aliases are the point: commit.md's Red Flags warn that testing only for the
# literal `completed` drops every remediated-with-issues REQ, and the contract
# additionally aliases done/finished/closed onto completed. Both prose copies
# had to restate this; now neither does.
is_terminal_success_status() {
  case "$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]')" in
    completed|completed-with-issues|done|finished|closed) return 0 ;;
  esac
  return 1
}

read_frontmatter_field() {
  awk -v wanted_field="$2" '
    NR==1 && /^---/ {inside_frontmatter=1; next}
    inside_frontmatter && /^---/ {exit}
    inside_frontmatter {
      field_separator = index($0, ":")
      if (field_separator == 0) next
      field_name = substr($0, 1, field_separator - 1)
      if (field_name != wanted_field) next
      field_value = substr($0, field_separator + 1)
      gsub(/^[ \t]+|[ \t]+$/, "", field_value)
      gsub(/^["'\'']|["'\'']$/, "", field_value)
      print field_value
      exit
    }
  ' "$1"
}

# Backticked paths from the "- `path` (verb)" bullets of one section, do-work/
# metadata excluded by contract. Same parse tools/checks/scope-drift.sh uses,
# so the two agree on what an Implementation Summary file list is.
extract_implementation_summary_paths() {
  awk '
    /^## Implementation Summary$/ {inside_section=1; next}
    inside_section && /^## / {inside_section=0}
    inside_section {print}
  ' "$1" | sed -n 's/^[[:space:]]*- `\([^`]*\)`.*/\1/p' | grep -v '^do-work/'
}

# owning_req_by_path / owning_stamp_by_path hold the current best claim per
# path. The stamp is the tie-break: when two REQs both list a path, the most
# recently completed one wins. An in-flight working/ REQ has no completed_at
# and is stamped empty, which sorts below every archived claim — deliberately,
# since a finished REQ is stronger evidence of ownership than a running one.
declare -A owning_req_by_path=()
declare -A owning_stamp_by_path=()

consider_request_file() {
  local request_file="$1" require_terminal_success="$2"
  [ -f "$request_file" ] || return 0

  local request_id request_status completed_stamp
  request_id="$(read_frontmatter_field "$request_file" id)"
  [ -n "$request_id" ] || return 0

  if [ "$require_terminal_success" = "yes" ]; then
    request_status="$(read_frontmatter_field "$request_file" status)"
    is_terminal_success_status "$request_status" || return 0
  fi
  completed_stamp="$(read_frontmatter_field "$request_file" completed_at)"

  local summary_path
  while IFS= read -r summary_path; do
    [ -n "$summary_path" ] || continue
    if [ -z "${owning_req_by_path[$summary_path]:-}" ] ||
       [[ "$completed_stamp" > "${owning_stamp_by_path[$summary_path]:-}" ]]; then
      owning_req_by_path["$summary_path"]="$request_id"
      owning_stamp_by_path["$summary_path"]="$completed_stamp"
    fi
  done < <(extract_implementation_summary_paths "$request_file")
}

while IFS= read -r archived_request_file; do
  consider_request_file "$archived_request_file" yes
done < <(find do-work/archive -type f -name 'REQ-*.md' 2>/dev/null | sort)

# In-flight REQs are considered regardless of status — they are claimed, not
# finished, so a terminal-success gate would exclude every one of them.
while IFS= read -r working_request_file; do
  consider_request_file "$working_request_file" no
done < <(find do-work/working -type f -name 'REQ-*.md' 2>/dev/null | sort)

for candidate_path in "${filtered_candidates[@]}"; do
  printf '%s\t%s\n' "${owning_req_by_path[$candidate_path]:--}" "$candidate_path"
done
exit 0
