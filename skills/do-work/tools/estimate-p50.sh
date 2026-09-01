#!/usr/bin/env bash
# do-work-cli compatibility launcher: estimate-p50
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
exec bash "$script_directory/do-work-cli.sh" --format text estimate-p50 "$@"
