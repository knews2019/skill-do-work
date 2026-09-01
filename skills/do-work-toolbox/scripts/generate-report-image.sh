#!/usr/bin/env bash
# do-work-cli compatibility launcher: generate-report-image
set -euo pipefail

script_directory="${BASH_SOURCE[0]%/*}"
[[ "$script_directory" != "${BASH_SOURCE[0]}" ]] || script_directory=.
script_directory="$(cd "$script_directory" && pwd -P)"
launcher_arguments=(--format text)
if [[ -n "${DO_WORK_COMPATIBILITY_REPO_ROOT:-}" ]]; then
  launcher_arguments+=(--repo-root "$DO_WORK_COMPATIBILITY_REPO_ROOT")
fi
exec /bin/bash "$script_directory/../../do-work/tools/do-work-cli.sh" "${launcher_arguments[@]}" generate-report-image "$@"
