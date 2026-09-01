#!/usr/bin/env bash
# do-work-cli compatibility launcher: run-blocked-check
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 2 ]]; then
  set -- --probe-file "$1" --timeout-seconds "$2"
elif [[ "$#" -eq 1 ]]; then
  set -- --probe-file "$1"
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text run-blocked-check "$@"
