#!/usr/bin/env bash
# do-work-cli compatibility launcher: associate-files
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
preamble_path="$script_directory/../do-work-cli-preamble.sh"
[ -f "$preamble_path" ] || { printf 'associate-files: do-work-cli-preamble.sh is missing beside this launcher\n' >&2; exit 2; }
# shellcheck source-path=SCRIPTDIR source=../do-work-cli-preamble.sh
. "$preamble_path"
if [[ "${1:-}" == "--repo-root" ]]; then
  if [[ "$#" -lt 2 || -z "$2" ]]; then
    exit 2
  fi
  # The requested repository root leads the preamble's shared argument prefix.
  launcher_arguments=(--repo-root "$2" "${launcher_arguments[@]}")
  shift 2
fi
# shellcheck disable=SC2154 # do_work_cli is set by the sourced preamble.
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$do_work_cli" "${launcher_arguments[@]}" associate-files --paths-file /dev/stdin "$@"
