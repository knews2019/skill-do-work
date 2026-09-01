#!/usr/bin/env bash
# do-work-cli compatibility launcher: scope-drift
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 1 ]]; then
  set -- --request-path "$1"
fi
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../do-work-cli.sh" --format text scope-drift "$@"
