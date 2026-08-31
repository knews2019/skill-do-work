#!/usr/bin/env bash
# Update a project-local four-module do-work suite.
#
# The transaction lives in the do-work-cli `update-suite` command; this file is the
# compatibility launcher that keeps the public argv working. The command fetches, validates
# and installs in one process, so the installer is no longer run as a subprocess and there is
# no cancellation status to thread through. Requires Go 1.25.0 or newer, which the
# do-work-cli launcher enforces.
set -euo pipefail

fail() {
  printf 'do-work update: %s\n' "$*" >&2
  exit 1
}

project_root=''
if [ "$#" -eq 2 ] && [ "$1" = '--project-root' ]; then
  project_root="$2"
else
  fail 'usage: do-work-update.sh --project-root <project-root>'
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
skill_root="$(cd "$script_dir/.." && pwd -P)"
do_work_cli="$script_dir/do-work-cli.sh"
[ -f "$do_work_cli" ] || fail 'do-work-cli.sh is missing beside the updater'

# --skill-root names the installed skill this launcher belongs to, so the command's
# refuse-a-shared-install guard measures the same tree the user is actually running.
launcher_status=0
bash "$do_work_cli" --repo-root "$project_root" --format text update-suite \
  --skill-root "$skill_root" || launcher_status=$?
exit "$launcher_status"
