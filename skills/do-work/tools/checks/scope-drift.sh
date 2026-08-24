#!/usr/bin/env bash
# scope-drift.sh — mechanical form of the Step 5.5 review-time comparison in
# actions/work.md: set-difference between the files declared in "## Scope" and
# the files reported in "## Implementation Summary".
#
# Usage: tools/checks/scope-drift.sh <req-file>
# Exit 0: no drift. Exit 1: drift found, or a touch-list exists but cannot be
#         parsed (details on stdout).
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

path_parse_error_marker='__DO_WORK_UNMATCHED_PATH_BACKTICK__'

extract_section_paths() {
  local section_heading="$1"
  # Lines of the section up to the next ## heading; every closed backtick pair
  # from path-led bullets is a path. Requiring the bullet to lead with a path
  # keeps later code spans on prose-only bullets out without dropping valid
  # root-level filenames such as README.md or justfile.
  awk -v section="^## ${section_heading}$" -v parse_error="$path_parse_error_marker" '
    function emit_backticked_paths(line, part_count, part_index, parts) {
      part_count=split(line, parts, "`")
      if (part_count % 2 == 0) {
        print parse_error
        return
      }
      for (part_index=2; part_index<part_count; part_index+=2)
        if (parts[part_index] != "") print parts[part_index]
    }
    $0 ~ section {inside=1; next}
    inside && /^## / {inside=0}
    inside && /^[[:space:]]*-[[:space:]]*`/ {emit_backticked_paths($0)}
  ' "$request_file" | grep -v '^do-work/' | sort -u
}

# Declared work is ONLY the "Files I will touch" list — the template's
# "Files I will NOT touch" bullets document exclusions and must not count
# as declarations (they'd surface as false "declared but never touched" drift).
# Paths are accepted in BOTH shipped styles: inline on the header line
# ("**Files I will touch:** `a.js`, `b.js`") and as "- `path`" bullets below
# it. An awk that only took lines AFTER the header silently dropped inline
# lists, emptied declared_paths, and turned the whole check into a SKIP.
# The header may carry a trailing annotation before the colon ("**Files I
# will touch (all new):**") — recognize the header, not one literal spelling:
# both this extraction and the FAIL guard below use the same tolerant match
# ("Files I will touch" up to the first colon), so an annotated header either
# parses or fails loudly. A bare-literal match here degraded annotated
# headers to SKIP, silently disabling the comparison (REQ-179).
declared_paths="$(awk '
    /^## Scope$/ {inside_scope=1; next}
    inside_scope && /^## / {inside_scope=0}
    inside_scope {print}
  ' "$request_file" | awk -v parse_error="$path_parse_error_marker" '
    function emit_backticked_paths(line, part_count, part_index, parts) {
      part_count=split(line, parts, "`")
      if (part_count % 2 == 0) {
        print parse_error
        return
      }
      for (part_index=2; part_index<part_count; part_index+=2)
        if (parts[part_index] != "") print parts[part_index]
    }
    /\*\*Files I will touch[^:]*:/ {
      taking_list=1
      header_rest=$0
      sub(/.*\*\*Files I will touch[^:]*:/, "", header_rest)
      emit_backticked_paths(header_rest)
      next
    }
    taking_list && /^\*\*/ {taking_list=0}
    taking_list && /^[[:space:]]*-[[:space:]]*`/ {
      emit_backticked_paths($0)
    }
  ' | grep -v '^do-work/' | sort -u)"
reported_paths="$(extract_section_paths 'Implementation Summary')"

if grep -Fx "$path_parse_error_marker" <<<"$declared_paths" >/dev/null; then
  echo "FAIL: the ## Scope path list has an unmatched backtick — close every backticked path"
  exit 1
fi
if grep -Fx "$path_parse_error_marker" <<<"$reported_paths" >/dev/null; then
  echo "FAIL: the ## Implementation Summary path list has an unmatched backtick — close every backticked path"
  exit 1
fi

if [ -z "$declared_paths" ]; then
  # A touch-list header that exists but yields zero paths is a formatting
  # defect, not a Route A REQ — SKIPping it would silently disable the check.
  # grep without -q: -q quits on first match and can SIGPIPE the awk upstream,
  # which pipefail then reads as "no match" — the same trap qualify.sh had.
  if awk '/^## Scope$/ {inside_scope=1; next} inside_scope && /^## / {inside_scope=0} inside_scope {print}' "$request_file" \
      | grep '\*\*Files I will touch[^:]*:' >/dev/null; then
    echo "FAIL: a '**Files I will touch**' header is present in ## Scope but no backticked paths parse from it — fix the list formatting (backtick every path)"
    exit 1
  fi
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
  printf '%s\n' "$undeclared_touches" | sed 's/^/  /'
  drift_found=1
fi
if [ -n "$unused_declarations" ]; then
  echo "DRIFT: declared in ## Scope but never touched:"
  printf '%s\n' "$unused_declarations" | sed 's/^/  /'
  drift_found=1
fi
[ "$drift_found" -eq 0 ] && echo "OK: Implementation Summary matches the Scope declaration"
exit "$drift_found"
