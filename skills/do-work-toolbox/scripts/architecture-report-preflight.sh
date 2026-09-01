#!/usr/bin/env bash
# do-work-cli compatibility launcher: architecture-report-preflight
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
repository_probe='.'
if [[ "${1:-}" == --scan && "$#" -eq 2 ]]; then
  repository_probe="$2"
elif [[ "${1:-}" == --publish && "$#" -eq 3 ]]; then
  repository_probe="$(dirname "$2")"
fi
repository_root="$(git -C "$repository_probe" rev-parse --show-toplevel 2>/dev/null || pwd -P)"
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../../do-work/tools/do-work-cli.sh" --repo-root "$repository_root" --format text architecture-report-preflight "$@"
