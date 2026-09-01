#!/usr/bin/env bash
# do-work-cli compatibility launcher: associate-files
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
launcher_arguments=(--format text)
if [[ "${1:-}" == "--repo-root" ]]; then
  if [[ "$#" -lt 2 || -z "$2" ]]; then
    exit 2
  fi
  launcher_arguments=(--repo-root "$2" --format text)
  shift 2
fi
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../do-work-cli.sh" "${launcher_arguments[@]}" associate-files --paths-file /dev/stdin "$@"
