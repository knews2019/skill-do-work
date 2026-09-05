#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
# Validates the complete do-work suite in an extracted archive or staging tree.
#
# The implementation lives in the do-work-cli `validate-manifest` command; this file is the
# compatibility launcher that keeps the public argv working. Requires Go 1.25.0 or newer,
# which the do-work-cli launcher enforces.
set -euo pipefail

if [ "$#" -ne 2 ] || [ "$1" != '--root' ]; then
  printf 'suite manifest: usage: validate-suite-manifest.sh --root <archive-root>\n' >&2
  exit 2
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
# The preamble ships beside these launchers and inside the staged package; probe both.
preamble_path="$script_dir/do-work-cli-preamble.sh"
[ -f "$preamble_path" ] || preamble_path="$script_dir/../skills/do-work/tools/do-work-cli-preamble.sh"
# shellcheck source-path=SCRIPTDIR source=do-work-cli-preamble.sh
. "$preamble_path"
if [ -z "$do_work_cli" ]; then
  printf 'suite manifest: do-work-cli.sh is missing beside this launcher\n' >&2
  exit 2
fi

launcher_status=0
# shellcheck disable=SC2154 # launcher_arguments is set by the sourced preamble.
bash "$do_work_cli" "${launcher_arguments[@]}" validate-manifest --root "$2" || launcher_status=$?
exit "$launcher_status"
