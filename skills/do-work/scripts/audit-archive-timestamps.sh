#!/usr/bin/env bash
# do-work-cli compatibility launcher: audit-archive-timestamps
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 1 && "$1" != --fix ]]; then
  exec bash "$script_directory/../tools/do-work-cli.sh" --repo-root "$1" --format text audit-archive-timestamps
fi
if [[ "$#" -eq 2 && "$1" == --fix ]]; then
  exec bash "$script_directory/../tools/do-work-cli.sh" --repo-root "$2" --format text audit-archive-timestamps --fix
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text audit-archive-timestamps "$@"
