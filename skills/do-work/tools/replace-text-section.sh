#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
# Atomically create, append, or replace do-work's managed text section.
#
# The implementation lives in the do-work-cli `replace-section` command; this file is the
# compatibility launcher that keeps the public argv working. Requires Go 1.25.0 or newer,
# which the do-work-cli launcher enforces.
set -euo pipefail

USAGE="usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--reject-recipe-collisions] [--begin-marker <line> --end-marker <line>]"

for launcher_argument in "$@"; do
  if [ "$launcher_argument" = '--help' ]; then
    printf '%s\n' "$USAGE"
    exit 0
  fi
done

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
do_work_cli=''
for cli_candidate in \
  "$script_dir/do-work-cli.sh" \
  "$script_dir/../skills/do-work/tools/do-work-cli.sh"; do
  if [ -f "$cli_candidate" ]; then
    do_work_cli="$cli_candidate"
    break
  fi
done
if [ -z "$do_work_cli" ]; then
  printf 'replace-text-section: do-work-cli.sh is missing beside this launcher\n' >&2
  exit 2
fi

# Invoked rather than exec'd so this launcher keeps its own signal statuses.
launcher_status=0
bash "$do_work_cli" --format text replace-section "$@" || launcher_status=$?
exit "$launcher_status"
