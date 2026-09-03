#!/usr/bin/env bash
# Entry point for the prescribed-shell behavior proofs. The cases themselves live one file
# per script under test in prescribed-shell-cases/, over the shared fixture preamble in
# prescribed-shell-harness.sh; this file owns dispatch and the closing count and nothing
# else. Callers (today: _dev/tests/staged-skills-contract.sh) consume only the exit status.
set -uo pipefail

if [ "${DO_WORK_MAINTAINER_TIER:-}" != heavy ]; then
  printf 'prescribed-shell behavior probes are heavy-only; run _dev/tests/maintainer-verify.sh --heavy after user permission.\n' >&2
  exit 2
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# shellcheck source=_dev/tests/prescribed-shell-case-count.sh
source "$repo_root/_dev/tests/prescribed-shell-case-count.sh"
case_directory="$repo_root/_dev/tests/prescribed-shell-cases"
case_file_count=0
failed_case_file_count=0
named_case_count=0

# One process per case file, so a file's fixture root, background processes, and failure
# tally are its own — and so one file can be run on its own while iterating on its script.
for case_file in "$case_directory"/*.sh; do
  [ -f "$case_file" ] || continue
  case_file_count=$((case_file_count + 1))
  case_header_count="$(count_named_case_headers "$case_file")"
  named_case_count=$((named_case_count + case_header_count))
  bash "$case_file" || failed_case_file_count=$((failed_case_file_count + 1))
done

if [ "$case_file_count" -eq 0 ]; then
  printf 'FAIL: no per-script case files found under %s\n' "$case_directory" >&2
  exit 1
fi
if [ "$failed_case_file_count" -gt 0 ]; then
  printf 'FAIL: %s of %s per-script case files reported failures.\n' \
    "$failed_case_file_count" "$case_file_count" >&2
  exit 1
fi
# One case is one fixture block, and every block opens with a header comment at column
# zero naming the script under test. The aggregate accumulated in the loop above applies
# prescribed-shell-case-count.sh's rule — the same rule each case file reports its own
# line with — to every file in turn, so the reported number and the files cannot disagree,
# and nothing here is a remembered figure.
printf 'Prescribed shell script behavior probes passed (%s named script cases across %s per-script files).\n' \
  "$named_case_count" "$case_file_count"
