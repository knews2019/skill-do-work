#!/usr/bin/env bash
# do-work-cli compatibility launcher: generate-report-image-batch
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
launcher_arguments=(--format text)
if [[ -n "${DO_WORK_COMPATIBILITY_REPO_ROOT:-}" ]]; then
  launcher_arguments+=(--repo-root "$DO_WORK_COMPATIBILITY_REPO_ROOT")
fi
compatibility_arguments=("$@")
report_index=0
if [[ "${compatibility_arguments[0]:-}" == --dry-run || "${compatibility_arguments[0]:-}" == --commit ]]; then
  report_index=1
fi
if [[ -n "${compatibility_arguments[$report_index]:-}" && "${compatibility_arguments[$report_index]}" != /* ]]; then
  compatibility_arguments[$report_index]="$PWD/${compatibility_arguments[$report_index]}"
fi
exec bash "$script_directory/../../do-work/tools/do-work-cli.sh" "${launcher_arguments[@]}" generate-report-image-batch "${compatibility_arguments[@]}"
