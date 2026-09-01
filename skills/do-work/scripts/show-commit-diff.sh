#!/usr/bin/env bash
# do-work-cli compatibility launcher: show-commit-diff
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 1 ]]; then
  set -- --commit "$1"
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text show-commit-diff "$@"
