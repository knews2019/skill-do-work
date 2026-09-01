#!/usr/bin/env bash
# do-work-cli compatibility launcher: record-commit-hash
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
record_arguments=()
if [[ "${1:-}" == "--verify" ]]; then
  record_arguments+=(--verify)
  shift
fi
if [[ "$#" -eq 2 ]]; then
  record_arguments+=(--request-path "$1" --implementation-hash "$2")
  set --
fi
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../do-work-cli.sh" --format text record-commit-hash "${record_arguments[@]}" "$@"
