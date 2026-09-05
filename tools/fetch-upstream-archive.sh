#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
# Fetch the upstream suite archive by whichever route works, and publish it only whole.
#
# The implementation lives in the do-work-cli `fetch-archive` command; this file is the
# compatibility launcher that maps the three positional arguments onto its flags and keeps
# this script's conventional signal statuses. Requires Go 1.25.0 or newer, which the
# do-work-cli launcher enforces.
set -u

if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
  printf 'Usage: %s <archive-target-path> <upstream-tarball-url> [upstream-repo-url]\n' "$0" >&2
  exit 2
fi

# Cleanup stays on EXIT; a terminating signal exits with its conventional status and lets
# EXIT do the cleaning, so an interrupted fetch never resumes into a publish.
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
# The preamble ships beside these launchers and inside the staged package; probe both.
preamble_path="$script_dir/do-work-cli-preamble.sh"
[ -f "$preamble_path" ] || preamble_path="$script_dir/../skills/do-work/tools/do-work-cli-preamble.sh"
[ -f "$preamble_path" ] || { printf 'upstream archive could not be fetched. HTTP route: unavailable (do-work-cli-preamble.sh is missing beside this launcher). Git route: unavailable (do-work-cli-preamble.sh is missing beside this launcher).\nSet DO_WORK_UPSTREAM_URL to a reachable archive URL to route around a blocked host.\n' >&2; exit 1; }
# shellcheck source-path=SCRIPTDIR source=do-work-cli-preamble.sh
. "$preamble_path"
if [ -z "$do_work_cli" ]; then
  printf 'upstream archive could not be fetched. HTTP route: unavailable (do-work-cli.sh is missing beside this launcher). Git route: unavailable (do-work-cli.sh is missing beside this launcher).\n' >&2
  printf 'Set DO_WORK_UPSTREAM_URL to a reachable archive URL to route around a blocked host.\n' >&2
  exit 1
fi

cli_arguments=(--target "$1" --url "$2")
if [ "$#" -eq 3 ] && [ -n "$3" ]; then
  cli_arguments+=(--repo-url "$3")
fi
launcher_status=0
# shellcheck disable=SC2154 # launcher_arguments is set by the sourced preamble.
bash "$do_work_cli" "${launcher_arguments[@]}" fetch-archive "${cli_arguments[@]}" || launcher_status=$?
if [ "$launcher_status" -ne 0 ]; then
  exit 1
fi
exit 0
