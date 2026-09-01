#!/usr/bin/env bash
# do-work-cli compatibility launcher: atomic-download
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 2 ]]; then
  set -- --source-url "$1" --target-path "$2"
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text atomic-download "$@"
