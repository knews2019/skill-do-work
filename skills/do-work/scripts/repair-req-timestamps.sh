#!/usr/bin/env bash
# do-work-cli compatibility launcher: repair-req-timestamps
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
if [[ "$#" -eq 1 ]]; then
  exec bash "$script_directory/../tools/do-work-cli.sh" --repo-root "$1" --format text repair-req-timestamps
fi
exec bash "$script_directory/../tools/do-work-cli.sh" --format text repair-req-timestamps "$@"
