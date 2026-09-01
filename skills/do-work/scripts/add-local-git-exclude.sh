#!/usr/bin/env bash
# do-work-cli compatibility launcher: add-local-git-exclude
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -lt 1 || "$#" -gt 2 ]]; then
  exec bash "$script_directory/../tools/do-work-cli.sh" --format text add-local-git-exclude "$@"
fi
set -- --probe-path "$1" --pattern "${2:-**/${1#./}}"
exec bash "$script_directory/../tools/do-work-cli.sh" --format text add-local-git-exclude "$@"
