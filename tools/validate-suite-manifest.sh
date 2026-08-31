#!/usr/bin/env bash
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
  printf 'suite manifest: do-work-cli.sh is missing beside this launcher\n' >&2
  exit 2
fi

launcher_status=0
bash "$do_work_cli" --format text validate-manifest --root "$2" || launcher_status=$?
exit "$launcher_status"
