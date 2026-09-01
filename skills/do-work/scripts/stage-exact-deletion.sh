#!/usr/bin/env bash
# do-work-cli compatibility launcher: stage-exact-deletion
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 1 ]]; then
  set -- --path "$1"
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text stage-exact-deletion "$@"
