#!/usr/bin/env bash
# Audit do-work/archive/ for detectably wrong `*_at` stamps and repair them from
# git history. This is the stated mechanical-timestamp-repair exception to the
# archive immutability rule (actions/capture.md § Immutability Rule): a
# fabricated or timezone-shifted stamp in an archived REQ is wrong forever, the
# board warns on every render, and nothing else may correct it.
#
#   bash scripts/audit-archive-timestamps.sh [--fix] [project-root]
#
# Report-only by default: each defect prints as `would repair ...` and the run
# exits nonzero so a caller can gate on findings. `--fix` rewrites the stamps in
# the working tree; committing the repaired files follows the normal commit
# flow. NEVER wire this script into a hook — repairing the archive stays a
# conscious, deliberate invocation, which is exactly what keeps the exception
# narrower than the rule.
#
# WHAT COUNTS AS WRONG is the repairer's predicate, shared by sourcing
# scripts/repair-req-timestamps.sh rather than duplicating it: a top-level
# frontmatter `*_at` value beyond the 2-minute future skew allowance, or an
# impossible `created_at <= claimed_at <= completed_at` ordering. The value
# shapes recognized, the clamp, the audit-line format, and the
# verify-before-replace atomic-write guards are all the sourced file's —
# widening a shape there widens this audit in the same edit.
#
# WHERE THE REPLACEMENT COMES FROM: git only. The author time of the commit
# that introduced the stamp line, read with `git blame --line-porcelain`. File
# mtimes are never consulted here — a checkout resets them, so they carry no
# signal for committed archive content. A defect whose introducing commit
# cannot answer (an uncommitted or locally modified archive file, no git) is
# reported and left byte-identical, never invented.
#
# Exit 0 — every scanned archive file is clean, or --fix repaired every defect.
# Exit 1 — at least one defect was found and not fixed: a report-only finding,
#          an underivable replacement, or a tripped write guard. Guarded files
#          are byte-identical to how they were found.
set -uo pipefail

apply_fixes=0
project_root_argument=''
for script_argument in "$@"; do
  case "$script_argument" in
    --fix) apply_fixes=1 ;;
    --*)
      printf 'usage: audit-archive-timestamps.sh [--fix] [project-root]\n' >&2
      exit 2
      ;;
    *)
      if [ -n "$project_root_argument" ]; then
        printf 'usage: audit-archive-timestamps.sh [--fix] [project-root]\n' >&2
        exit 2
      fi
      project_root_argument="$script_argument"
      ;;
  esac
done

audit_project_root="${project_root_argument:-${CLAUDE_PROJECT_DIR:-.}}"
archive_root="$audit_project_root/do-work/archive"

# Never write through a link — same refusal as the sourced repairer, but voiced,
# because this run was asked for deliberately and silence would read as clean.
if [ -L "$audit_project_root/do-work" ] || [ -L "$archive_root" ]; then
  printf 'audit-archive-timestamps: refusing a symlinked do-work/ or archive/ path.\n' >&2
  exit 1
fi
if [ ! -d "$archive_root" ]; then
  printf 'do-work: archive audit clean (no archive directory).\n'
  exit 0
fi

# The shared machinery: detection, derivation, clamp, guards, and the
# repair_request_file worker, initialized against this project root. The
# switches must be set before the source — the sourced file reads them once.
# The counters are pre-seeded only so a lint run that does not follow the
# source still sees them assigned; the sourced file re-initializes all three.
# shellcheck disable=SC2034  # read by the sourced repairer, not by this file
timestamp_repair_apply_mode="$apply_fixes"
# shellcheck disable=SC2034  # read by the sourced repairer, not by this file
timestamp_repair_git_only=1
failure_count=0
repair_count=0
pending_repair_count=0
script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=skills/do-work/scripts/repair-req-timestamps.sh
source "$script_directory/repair-req-timestamps.sh" "$audit_project_root"

# The scan: archive only, at any depth (top-level legacy REQs, UR folders, and
# any layout a later schema adds — the pattern is the REQ filename, not a fixed
# directory list). find does not follow symlinked directories, and the worker
# itself refuses symlinked files. Sorted so the audit trail is deterministic.
scanned_file_count=0
while IFS= read -r -d '' archived_request_file; do
  scanned_file_count=$((scanned_file_count + 1))
  repair_request_file "$archived_request_file" "${archived_request_file#"$audit_project_root"/}"
done < <(find "$archive_root" -name 'REQ-*.md' -print0 | sort -z)

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi
if [ "$pending_repair_count" -gt 0 ]; then
  printf 'do-work: %s archived correction(s) pending — rerun with --fix to write them.\n' \
    "$pending_repair_count"
  exit 1
fi
if [ "$repair_count" -gt 0 ]; then
  printf 'do-work: repaired %s archived timestamp(s) — review and commit the correction(s) through the normal flow.\n' \
    "$repair_count"
  exit 0
fi
printf 'do-work: archive audit clean (%s file(s) scanned).\n' "$scanned_file_count"
exit 0
