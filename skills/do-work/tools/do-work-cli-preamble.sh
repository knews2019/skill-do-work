#!/usr/bin/env bash
# shellcheck disable=SC2034 # do_work_cli and launcher_arguments are read by the sourcing launcher.
# Shared do-work-cli launcher preamble. Sourced by every compatibility launcher, never run.
#
# Resolution is anchored on this file's own directory, not the sourcing launcher's, so the
# root copy and the byte-identical staged mirror each reach the do-work-cli.sh beside them:
# beside this preamble first, then the staged package one level up, which is where an
# extracted source archive keeps it. do_work_cli stays empty when neither exists, because
# each launcher owns a different missing-launcher message and exit status.

do_work_cli_preamble_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
do_work_cli=''
for cli_candidate in \
  "$do_work_cli_preamble_directory/do-work-cli.sh" \
  "$do_work_cli_preamble_directory/../skills/do-work/tools/do-work-cli.sh"; do
  if [ -f "$cli_candidate" ]; then
    do_work_cli="$cli_candidate"
    break
  fi
done
unset do_work_cli_preamble_directory cli_candidate

# The argument prefix every launcher passes ahead of its own command and flags.
launcher_arguments=(--format text)
