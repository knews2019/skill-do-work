#!/usr/bin/env bash
# do-work-cli compatibility launcher: uncommitted-inventory
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 1 ]]; then
  DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../do-work-cli.sh" --repo-root "$1" --format text uncommitted-inventory
fi
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../do-work-cli.sh" --format text uncommitted-inventory "$@"
