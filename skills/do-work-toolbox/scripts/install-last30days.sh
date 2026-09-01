#!/usr/bin/env bash
# do-work-cli compatibility launcher: install-last30days
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
launcher_arguments=(--format text)
if [[ -n "${DO_WORK_COMPATIBILITY_REPO_ROOT:-}" ]]; then
  launcher_arguments+=(--repo-root "$DO_WORK_COMPATIBILITY_REPO_ROOT")
fi
exec bash "$script_directory/../../do-work/tools/do-work-cli.sh" "${launcher_arguments[@]}" install-last30days "$@"
