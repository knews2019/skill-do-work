#!/usr/bin/env bash
# scope-drift.sh — mechanical form of the Step 5.5 review-time comparison in
# actions/work.md: set-difference between the files declared in "## Scope" and
# the files reported in "## Implementation Summary".
#
# Usage: tools/checks/scope-drift.sh <req-file>
# Exit 0: no drift. Exit 1: drift found (details on stdout).
# Exit 2: usage error or a section is missing (Route A REQs have no Scope
#         section — the caller skips the comparison, exactly as the prose says).
#
# Severity is judgment: the orchestrator decides Important vs Minor per
# actions/work.md Step 5.5. This script only computes the two lists.
set -uo pipefail

request_file="${1:-}"
if [ ! -f "$request_file" ]; then
  echo "usage: $0 <req-file>" >&2
  exit 2
fi

extract_section_paths() {
  local section_heading="$1"
  # Lines of the section up to the next ## heading; backtick-quoted paths from
  # "- `path` (verb)" bullets; do-work/ metadata excluded by contract.
  awk -v section="^## ${section_heading}$" '
    $0 ~ section {inside=1; next}
    inside && /^## / {inside=0}
    inside {print}
  ' "$request_file" | sed -n 's/^[[:space:]]*- `\([^`]*\)`.*/\1/p' | grep -v '^do-work/' | sort -u
}

declared_paths="$(extract_section_paths 'Scope')"
reported_paths="$(extract_section_paths 'Implementation Summary')"

if [ -z "$declared_paths" ]; then
  echo "SKIP: no '## Scope' file list found (Route A REQs have none — skip the comparison)"
  exit 2
fi
if [ -z "$reported_paths" ]; then
  echo "SKIP: no '## Implementation Summary' file list found — run this after Step 6.25"
  exit 2
fi

undeclared_touches="$(comm -13 <(printf '%s\n' "$declared_paths") <(printf '%s\n' "$reported_paths"))"
unused_declarations="$(comm -23 <(printf '%s\n' "$declared_paths") <(printf '%s\n' "$reported_paths"))"

drift_found=0
if [ -n "$undeclared_touches" ]; then
  echo "DRIFT: touched but never declared in ## Scope:"
  printf '  %s\n' $undeclared_touches
  drift_found=1
fi
if [ -n "$unused_declarations" ]; then
  echo "DRIFT: declared in ## Scope but never touched:"
  printf '  %s\n' $unused_declarations
  drift_found=1
fi
[ "$drift_found" -eq 0 ] && echo "OK: Implementation Summary matches the Scope declaration"
exit "$drift_found"
