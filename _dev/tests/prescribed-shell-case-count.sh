#!/usr/bin/env bash
# The one definition of what a prescribed-shell case header is. Sourced by both counters —
# the per-file line in prescribed-shell-harness.sh and the aggregate in
# prescribed-shell-scripts-behavior.sh — so the two figures cannot disagree with each other
# or with the files. This file runs nothing itself.

# A case header is a column-zero comment that opens with the name of the script the case
# file covers — the file's own basename — and reaches a colon before any period:
#
#   # <script-name>: <what it proves>
#   # <script-name><qualifier>: <what it proves>
#
count_named_case_headers() {
  local case_file_path="$1"
  local script_under_test="${case_file_path##*/}"
  script_under_test="${script_under_test%.sh}"

  # Keep the qualifier open-ended, but require a real boundary after the exact script
  # name. Stopping at a period keeps wrapped prose such as `# qualify.sh ...` out too.
  grep -cE "^# ${script_under_test}([^[:alnum:]_.:][^.:]*)?: " "$case_file_path"
}
