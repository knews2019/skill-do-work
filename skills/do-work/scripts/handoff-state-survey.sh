#!/usr/bin/env bash
# do-work-cli compatibility launcher: handoff-state-survey
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
exec bash "$script_directory/../tools/do-work-cli.sh" --format text handoff-state-survey "$@"
