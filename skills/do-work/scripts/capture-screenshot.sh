#!/usr/bin/env bash
# do-work-cli compatibility launcher: capture-screenshot
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 3 ]]; then
  set -- "$1" --source "$2" --destination "$3"
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text capture-screenshot "$@"
