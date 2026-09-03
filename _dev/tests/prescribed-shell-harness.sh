#!/usr/bin/env bash
# Shared fixture preamble for the per-script proof files in prescribed-shell-cases/.
#
# Each case file sources this, runs one script's fixture blocks, and closes with
# prescribed_shell_finish. The runner (prescribed-shell-scripts-behavior.sh) executes every
# case file as its own process, so each gets an independent fixture root, trap, and failure
# tally — which is also what makes a single case file runnable on its own while iterating
# on the script it covers.
set -uo pipefail

if [ "${DO_WORK_MAINTAINER_TIER:-}" != heavy ]; then
  printf 'prescribed-shell behavior probes are heavy-only; run _dev/tests/maintainer-verify.sh --heavy after user permission.\n' >&2
  exit 2
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
# The file that sourced this one, so the closing line can count its cases at run time
# instead of carrying a remembered figure (REQ-234).
case_file_path="${BASH_SOURCE[1]}"
# shellcheck source=_dev/tests/fixture-repo.sh
source "$repo_root/_dev/tests/fixture-repo.sh"
# shellcheck source=_dev/tests/prescribed-shell-case-count.sh
source "$repo_root/_dev/tests/prescribed-shell-case-count.sh"
fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/prescribed-shell-scripts.XXXXXX")" || exit 1
background_process_ids=""
cleanup_prescribed_shell_fixture() {
  for background_process_id in $background_process_ids; do
    kill "$background_process_id" 2>/dev/null || true
    wait "$background_process_id" 2>/dev/null || true
  done
  chmod -R u+rwX "$fixture_root" 2>/dev/null || true
  rm -rf "$fixture_root"
}
trap cleanup_prescribed_shell_fixture EXIT
failure_count=0

fail_case() {
  printf 'FAIL: %s\n' "$1" >&2
  failure_count=$((failure_count + 1))
}

# The script roots the case files reach for. ShellCheck cannot see the sourcing files that
# consume them, so its unused warning here is structural rather than a real finding.
# shellcheck disable=SC2034
{
  core_scripts="$repo_root/skills/do-work/scripts"
  knowledge_scripts="$repo_root/skills/do-work-knowledge/scripts"
  toolbox_scripts="$repo_root/skills/do-work-toolbox/scripts"
}

# Every case file ends with this. Its status is that file's own failure tally, which is what
# the runner aggregates and what a standalone invocation reports.
prescribed_shell_finish() {
  local case_file_name
  local named_case_count
  local case_noun=cases
  local failure_noun=failures

  case_file_name="${case_file_path##*/}"
  # What counts as a case header — and therefore what this number means — is stated once,
  # beside the grep that applies it, in prescribed-shell-case-count.sh.
  named_case_count="$(count_named_case_headers "$case_file_path")"
  if [ "$named_case_count" -eq 1 ]; then
    case_noun=case
  fi
  if [ "$failure_count" -eq 1 ]; then
    failure_noun=failure
  fi
  printf '%s: %s %s, %s %s.\n' "${case_file_name%.sh}" \
    "$named_case_count" "$case_noun" "$failure_count" "$failure_noun"
  [ "$failure_count" -eq 0 ]
}
