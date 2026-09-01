#!/usr/bin/env bash
# do-work-cli compatibility launcher: preflight
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../do-work-cli.sh" --format text preflight "$@"
