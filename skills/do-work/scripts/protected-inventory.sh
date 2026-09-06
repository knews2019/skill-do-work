#!/usr/bin/env bash
# do-work-cli compatibility launcher: protected-inventory
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
global_arguments=()
command_arguments=()
while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --repo-root)
      if [[ "$#" -lt 2 || -z "$2" ]]; then
        exit 2
      fi
      global_arguments+=(--repo-root "$2")
      shift 2
      ;;
    --repo-root=*)
      global_arguments+=("$1")
      shift
      ;;
    --format)
      if [[ "$#" -lt 2 || -z "$2" ]]; then
        exit 2
      fi
      global_arguments+=(--format "$2")
      shift 2
      ;;
    --format=*)
      global_arguments+=("$1")
      shift
      ;;
    *)
      command_arguments+=("$1")
      shift
      ;;
  esac
done

DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../tools/do-work-cli.sh" "${global_arguments[@]}" --format text protected-inventory "${command_arguments[@]}"

